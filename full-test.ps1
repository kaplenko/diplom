$ErrorActionPreference = "Continue"
$BASE = "http://localhost:8080/api/v1"
$pass = 0; $fail = 0; $total = 0
$tmpFile = [System.IO.Path]::GetTempFileName()
$utf8 = New-Object System.Text.UTF8Encoding($false)

function T {
    param([string]$Name, [string]$Method, [string]$Url, [string]$Body, [string]$Token, [int]$Expect)
    $script:total++
    $curlArgs = @("-s", "-w", "`n%{http_code}", "-X", $Method, "-H", "Content-Type: application/json")
    if ($Token) { $curlArgs += @("-H", "Authorization: Bearer $Token") }
    if ($Body) {
        [System.IO.File]::WriteAllText($tmpFile, $Body, $utf8)
        $curlArgs += @("-d", "@$tmpFile")
    }
    $curlArgs += $Url
    $raw = (& curl.exe @curlArgs 2>&1) -join "`n"
    $lines = $raw.Trim().Split("`n")
    $statusStr = $lines[-1].Trim()
    $status = 0; [int]::TryParse($statusStr, [ref]$status) | Out-Null
    $body = if ($lines.Count -gt 1) { ($lines[0..($lines.Count-2)] -join "`n").Trim() } else { "" }
    if ($status -eq $Expect) {
        Write-Host "  PASS [$status]: $Name" -ForegroundColor Green
        $script:pass++
    } else {
        Write-Host "  FAIL [$status exp $Expect]: $Name" -ForegroundColor Red
        if ($body.Length -gt 0 -and $body.Length -lt 300) { Write-Host "    $body" -ForegroundColor DarkGray }
        $script:fail++
    }
    return $body
}

Write-Host "`n===== HEALTHCHECK =====" -ForegroundColor Cyan
T -Name "Swagger UI" -Method GET -Url "http://localhost:8080/swagger/index.html" -Expect 200 | Out-Null

# ====================== AUTH ======================
Write-Host "`n===== 1. AUTH =====" -ForegroundColor Cyan

$ts = [DateTimeOffset]::UtcNow.ToUnixTimeMilliseconds()
T -Name "Register new user" -Method POST -Url "$BASE/auth/register" -Body "{`"email`":`"test${ts}@test.com`",`"password`":`"secure123`",`"name`":`"Test User`"}" -Expect 201 | Out-Null
T -Name "Register duplicate (admin)" -Method POST -Url "$BASE/auth/register" -Body '{"email":"admin@example.com","password":"secure123","name":"Dupe"}' -Expect 409 | Out-Null
T -Name "Register empty body" -Method POST -Url "$BASE/auth/register" -Body '{}' -Expect 400 | Out-Null
T -Name "Register missing password" -Method POST -Url "$BASE/auth/register" -Body '{"email":"x@x.com","name":"X"}' -Expect 400 | Out-Null
T -Name "Register short password" -Method POST -Url "$BASE/auth/register" -Body '{"email":"s@x.com","password":"abc","name":"S"}' -Expect 400 | Out-Null
T -Name "Register invalid email" -Method POST -Url "$BASE/auth/register" -Body '{"email":"not-email","password":"secure123","name":"Bad"}' -Expect 400 | Out-Null
T -Name "Register no body" -Method POST -Url "$BASE/auth/register" -Expect 400 | Out-Null

$r = T -Name "Login admin" -Method POST -Url "$BASE/auth/login" -Body '{"email":"admin@example.com","password":"password123"}' -Expect 200
$tok = $r | ConvertFrom-Json; $ADMIN = $tok.access_token; $ADMIN_R = $tok.refresh_token

$r = T -Name "Login student" -Method POST -Url "$BASE/auth/login" -Body '{"email":"student@example.com","password":"password123"}' -Expect 200
$tok = $r | ConvertFrom-Json; $STUDENT = $tok.access_token

T -Name "Login wrong password" -Method POST -Url "$BASE/auth/login" -Body '{"email":"admin@example.com","password":"wrong"}' -Expect 401 | Out-Null
T -Name "Login nonexistent user" -Method POST -Url "$BASE/auth/login" -Body '{"email":"ghost@x.com","password":"x"}' -Expect 401 | Out-Null
T -Name "Login empty body" -Method POST -Url "$BASE/auth/login" -Body '{}' -Expect 400 | Out-Null
T -Name "Login malformed JSON" -Method POST -Url "$BASE/auth/login" -Body '{broken' -Expect 400 | Out-Null

T -Name "Refresh token" -Method POST -Url "$BASE/auth/refresh" -Body "{`"refresh_token`":`"$ADMIN_R`"}" -Expect 200 | Out-Null
T -Name "Refresh with access token" -Method POST -Url "$BASE/auth/refresh" -Body "{`"refresh_token`":`"$ADMIN`"}" -Expect 401 | Out-Null
T -Name "Refresh with garbage" -Method POST -Url "$BASE/auth/refresh" -Body '{"refresh_token":"not.a.real.token"}' -Expect 401 | Out-Null
T -Name "Refresh empty body" -Method POST -Url "$BASE/auth/refresh" -Body '{}' -Expect 400 | Out-Null

# ====================== AUTHORIZATION ======================
Write-Host "`n===== 2. AUTHORIZATION =====" -ForegroundColor Cyan

T -Name "No token -> 401" -Method GET -Url "$BASE/users/me" -Expect 401 | Out-Null
T -Name "Garbage token -> 401" -Method GET -Url "$BASE/users/me" -Token "not-a-jwt" -Expect 401 | Out-Null
T -Name "Student -> list users (admin)" -Method GET -Url "$BASE/users" -Token $STUDENT -Expect 403 | Out-Null
T -Name "Student -> create course (admin)" -Method POST -Url "$BASE/courses" -Token $STUDENT -Body '{"title":"H","description":"x"}' -Expect 403 | Out-Null
T -Name "Student -> delete user (admin)" -Method DELETE -Url "$BASE/users/1" -Token $STUDENT -Expect 403 | Out-Null
T -Name "Student -> delete course (admin)" -Method DELETE -Url "$BASE/courses/1" -Token $STUDENT -Expect 403 | Out-Null
T -Name "Student -> create lesson (admin)" -Method POST -Url "$BASE/courses/1/lessons" -Token $STUDENT -Body '{"title":"H","content":"x","order_index":1}' -Expect 403 | Out-Null
T -Name "Student -> create task (admin)" -Method POST -Url "$BASE/lessons/1/tasks" -Token $STUDENT -Body '{"title":"T","description":"D","test_cases":[{"input":"","expected":""}],"difficulty":"easy"}' -Expect 403 | Out-Null

# ====================== USERS ======================
Write-Host "`n===== 3. USERS =====" -ForegroundColor Cyan

T -Name "Get me (student)" -Method GET -Url "$BASE/users/me" -Token $STUDENT -Expect 200 | Out-Null
T -Name "Update my name" -Method PUT -Url "$BASE/users/me" -Token $STUDENT -Body '{"name":"Renamed Student"}' -Expect 200 | Out-Null
T -Name "Update name empty (no-op)" -Method PUT -Url "$BASE/users/me" -Token $STUDENT -Body '{}' -Expect 200 | Out-Null
T -Name "Admin list users" -Method GET -Url "$BASE/users" -Token $ADMIN -Expect 200 | Out-Null
T -Name "Admin search users" -Method GET -Url "$BASE/users?search=admin" -Token $ADMIN -Expect 200 | Out-Null

# ====================== COURSES ======================
Write-Host "`n===== 4. COURSES =====" -ForegroundColor Cyan

$r = T -Name "List courses" -Method GET -Url "$BASE/courses" -Token $STUDENT -Expect 200
$parsed = $r | ConvertFrom-Json; Write-Host "    $($parsed.total) courses" -ForegroundColor DarkGray
T -Name "Get course 1" -Method GET -Url "$BASE/courses/1" -Token $STUDENT -Expect 200 | Out-Null
T -Name "Get course 999" -Method GET -Url "$BASE/courses/999" -Token $STUDENT -Expect 404 | Out-Null
T -Name "Get course abc" -Method GET -Url "$BASE/courses/abc" -Token $STUDENT -Expect 400 | Out-Null
T -Name "Admin create course" -Method POST -Url "$BASE/courses" -Token $ADMIN -Body '{"title":"Test Course","description":"Desc"}' -Expect 201 | Out-Null
T -Name "Create missing title" -Method POST -Url "$BASE/courses" -Token $ADMIN -Body '{"description":"no title"}' -Expect 400 | Out-Null
T -Name "Create empty body" -Method POST -Url "$BASE/courses" -Token $ADMIN -Body '{}' -Expect 400 | Out-Null
T -Name "Update course 3" -Method PUT -Url "$BASE/courses/3" -Token $ADMIN -Body '{"title":"Updated"}' -Expect 200 | Out-Null
T -Name "Update course 999" -Method PUT -Url "$BASE/courses/999" -Token $ADMIN -Body '{"title":"x"}' -Expect 404 | Out-Null
T -Name "Search courses" -Method GET -Url "$BASE/courses?search=Go" -Token $STUDENT -Expect 200 | Out-Null

# ====================== LESSONS ======================
Write-Host "`n===== 5. LESSONS =====" -ForegroundColor Cyan

T -Name "List lessons (course 1)" -Method GET -Url "$BASE/courses/1/lessons" -Token $STUDENT -Expect 200 | Out-Null
T -Name "Get lesson 1" -Method GET -Url "$BASE/lessons/1" -Token $STUDENT -Expect 200 | Out-Null
T -Name "Get lesson 999" -Method GET -Url "$BASE/lessons/999" -Token $STUDENT -Expect 404 | Out-Null
T -Name "Get lesson abc" -Method GET -Url "$BASE/lessons/abc" -Token $STUDENT -Expect 400 | Out-Null
T -Name "Create lesson" -Method POST -Url "$BASE/courses/1/lessons" -Token $ADMIN -Body '{"title":"Test Lesson","content":"Content","order_index":10}' -Expect 201 | Out-Null
T -Name "Create lesson in course 999" -Method POST -Url "$BASE/courses/999/lessons" -Token $ADMIN -Body '{"title":"G","content":"x","order_index":1}' -Expect 404 | Out-Null
T -Name "Create lesson missing title" -Method POST -Url "$BASE/courses/1/lessons" -Token $ADMIN -Body '{"content":"x","order_index":1}' -Expect 400 | Out-Null
T -Name "Update lesson 1" -Method PUT -Url "$BASE/lessons/1" -Token $ADMIN -Body '{"title":"Updated Hello World"}' -Expect 200 | Out-Null
T -Name "Update lesson 999" -Method PUT -Url "$BASE/lessons/999" -Token $ADMIN -Body '{"title":"x"}' -Expect 404 | Out-Null

# ====================== TASKS ======================
Write-Host "`n===== 6. TASKS =====" -ForegroundColor Cyan

T -Name "List tasks (lesson 1)" -Method GET -Url "$BASE/lessons/1/tasks" -Token $STUDENT -Expect 200 | Out-Null
T -Name "Get task 1" -Method GET -Url "$BASE/tasks/1" -Token $STUDENT -Expect 200 | Out-Null
T -Name "Get task 999" -Method GET -Url "$BASE/tasks/999" -Token $STUDENT -Expect 404 | Out-Null
T -Name "Get task abc" -Method GET -Url "$BASE/tasks/abc" -Token $STUDENT -Expect 400 | Out-Null

T -Name "Create task" -Method POST -Url "$BASE/lessons/1/tasks" -Token $ADMIN `
    -Body '{"title":"New Task","description":"desc","initial_code":"pkg","test_cases":[{"input":"","expected":"ok"}],"difficulty":"easy"}' -Expect 201 | Out-Null
T -Name "Update task with test_cases" -Method PUT -Url "$BASE/tasks/1" -Token $ADMIN `
    -Body '{"title":"Updated","test_cases":[{"input":"","expected":"ok"}]}' -Expect 200 | Out-Null

T -Name "Create task in lesson 999" -Method POST -Url "$BASE/lessons/999/tasks" -Token $ADMIN -Body '{"title":"T","description":"D","test_cases":[{"input":"","expected":""}],"difficulty":"easy"}' -Expect 404 | Out-Null
T -Name "Create task no test_cases" -Method POST -Url "$BASE/lessons/1/tasks" -Token $ADMIN -Body '{"title":"T","description":"D","difficulty":"easy"}' -Expect 400 | Out-Null
T -Name "Create task bad difficulty" -Method POST -Url "$BASE/lessons/1/tasks" -Token $ADMIN -Body '{"title":"T","description":"D","test_cases":[{"input":"","expected":""}],"difficulty":"insane"}' -Expect 400 | Out-Null

# Update without test_cases should work (only updates title)
T -Name "Update task title only" -Method PUT -Url "$BASE/tasks/1" -Token $ADMIN -Body '{"title":"Updated Task"}' -Expect 200 | Out-Null
T -Name "Update task 999" -Method PUT -Url "$BASE/tasks/999" -Token $ADMIN -Body '{"title":"x"}' -Expect 404 | Out-Null

# ====================== PAGINATION ======================
Write-Host "`n===== 7. PAGINATION =====" -ForegroundColor Cyan

T -Name "page=0 (normalize)" -Method GET -Url "$BASE/courses?page=0&page_size=10" -Token $STUDENT -Expect 200 | Out-Null
T -Name "page=-5 (normalize)" -Method GET -Url "$BASE/courses?page=-5&page_size=10" -Token $STUDENT -Expect 200 | Out-Null
T -Name "page_size=9999 (cap)" -Method GET -Url "$BASE/courses?page=1&page_size=9999" -Token $STUDENT -Expect 200 | Out-Null
T -Name "page_size=0 (default)" -Method GET -Url "$BASE/courses?page=1&page_size=0" -Token $STUDENT -Expect 200 | Out-Null
T -Name "page=abc (default)" -Method GET -Url "$BASE/courses?page=abc" -Token $STUDENT -Expect 200 | Out-Null
T -Name "page=99999 (empty)" -Method GET -Url "$BASE/courses?page=99999" -Token $STUDENT -Expect 200 | Out-Null

# ====================== SUBMISSIONS ======================
Write-Host "`n===== 8. SUBMISSIONS =====" -ForegroundColor Cyan

# Use proper JSON files for code submissions
$codeCorrect = '{"code":"package main\nimport \"fmt\"\nfunc main() { fmt.Println(\"Hello, World!\") }"}'
$r = T -Name "Submit correct (task 1)" -Method POST -Url "$BASE/tasks/1/submissions" -Token $STUDENT -Body $codeCorrect -Expect 201
if ($r -and $r -ne "") { $sub1 = ($r | ConvertFrom-Json).id; Write-Host "    id=$sub1" -ForegroundColor DarkGray } else { $sub1 = $null }

$codeWrong = '{"code":"package main\nimport \"fmt\"\nfunc main() { fmt.Println(\"wrong\") }"}'
$r = T -Name "Submit wrong (task 1)" -Method POST -Url "$BASE/tasks/1/submissions" -Token $STUDENT -Body $codeWrong -Expect 201
if ($r -and $r -ne "") { $sub2 = ($r | ConvertFrom-Json).id; Write-Host "    id=$sub2" -ForegroundColor DarkGray } else { $sub2 = $null }

$codeCompile = '{"code":"package main\nfunc main() { badcode }"}'
$r = T -Name "Submit compile err (task 1)" -Method POST -Url "$BASE/tasks/1/submissions" -Token $STUDENT -Body $codeCompile -Expect 201
if ($r -and $r -ne "") { $sub3 = ($r | ConvertFrom-Json).id; Write-Host "    id=$sub3" -ForegroundColor DarkGray } else { $sub3 = $null }

T -Name "Submit empty code" -Method POST -Url "$BASE/tasks/1/submissions" -Token $STUDENT -Body '{"code":""}' -Expect 400 | Out-Null
T -Name "Submit no body" -Method POST -Url "$BASE/tasks/1/submissions" -Token $STUDENT -Body '{}' -Expect 400 | Out-Null
T -Name "Submit to task 999" -Method POST -Url "$BASE/tasks/999/submissions" -Token $STUDENT -Body '{"code":"x"}' -Expect 404 | Out-Null
T -Name "Submit to task abc" -Method POST -Url "$BASE/tasks/abc/submissions" -Token $STUDENT -Body '{"code":"x"}' -Expect 400 | Out-Null

Write-Host "  Waiting 25s for runner..." -ForegroundColor Yellow
Start-Sleep -Seconds 25

if ($sub1) {
    $r = T -Name "Poll correct sub $sub1" -Method GET -Url "$BASE/submissions/$sub1" -Token $STUDENT -Expect 200
    $p = $r | ConvertFrom-Json; Write-Host "    status=$($p.status) score=$($p.score)" -ForegroundColor DarkGray
}
if ($sub2) {
    $r = T -Name "Poll wrong sub $sub2" -Method GET -Url "$BASE/submissions/$sub2" -Token $STUDENT -Expect 200
    $p = $r | ConvertFrom-Json; Write-Host "    status=$($p.status) score=$($p.score)" -ForegroundColor DarkGray
}
if ($sub3) {
    $r = T -Name "Poll compile-err sub $sub3" -Method GET -Url "$BASE/submissions/$sub3" -Token $STUDENT -Expect 200
    $p = $r | ConvertFrom-Json; Write-Host "    status=$($p.status) score=$($p.score)" -ForegroundColor DarkGray
}

T -Name "Get sub 999999" -Method GET -Url "$BASE/submissions/999999" -Token $STUDENT -Expect 404 | Out-Null
T -Name "Get sub abc" -Method GET -Url "$BASE/submissions/abc" -Token $STUDENT -Expect 400 | Out-Null
T -Name "List subs (task 1)" -Method GET -Url "$BASE/tasks/1/submissions" -Token $STUDENT -Expect 200 | Out-Null
T -Name "Admin list all subs" -Method GET -Url "$BASE/submissions" -Token $ADMIN -Expect 200 | Out-Null

# ====================== PROGRESS ======================
Write-Host "`n===== 9. PROGRESS =====" -ForegroundColor Cyan

$r = T -Name "Course 1 progress" -Method GET -Url "$BASE/courses/1/progress" -Token $STUDENT -Expect 200
$p = $r | ConvertFrom-Json; Write-Host "    completed=$($p.completed_count)/$($p.total_lessons) ($($p.percentage)%)" -ForegroundColor DarkGray
T -Name "All progress" -Method GET -Url "$BASE/progress" -Token $STUDENT -Expect 200 | Out-Null
T -Name "Lesson 1 progress" -Method GET -Url "$BASE/lessons/1/progress" -Token $STUDENT -Expect 200 | Out-Null
T -Name "Course 1 lessons progress" -Method GET -Url "$BASE/courses/1/lessons/progress" -Token $STUDENT -Expect 200 | Out-Null
T -Name "Progress course 999" -Method GET -Url "$BASE/courses/999/progress" -Token $STUDENT -Expect 404 | Out-Null
T -Name "Progress lesson 999" -Method GET -Url "$BASE/lessons/999/progress" -Token $STUDENT -Expect 404 | Out-Null
T -Name "Progress course abc" -Method GET -Url "$BASE/courses/abc/progress" -Token $STUDENT -Expect 400 | Out-Null

# ====================== UNEXPECTED ======================
Write-Host "`n===== 10. UNEXPECTED INPUTS =====" -ForegroundColor Cyan

T -Name "SQL injection in search" -Method GET -Url "$BASE/courses?search=';DROP+TABLE+courses;--" -Token $STUDENT -Expect 200 | Out-Null
T -Name "XSS in title" -Method POST -Url "$BASE/courses" -Token $ADMIN -Body '{"title":"<script>alert(1)</script>","description":"xss"}' -Expect 201 | Out-Null
T -Name "Unicode in name" -Method PUT -Url "$BASE/users/me" -Token $STUDENT -Body '{"name":"Test \u2764"}' -Expect 200 | Out-Null
T -Name "Negative ID (-1)" -Method GET -Url "$BASE/courses/-1" -Token $STUDENT -Expect 404 | Out-Null
T -Name "Zero ID (0)" -Method GET -Url "$BASE/courses/0" -Token $STUDENT -Expect 404 | Out-Null
T -Name "MAX INT64 ID" -Method GET -Url "$BASE/courses/9223372036854775807" -Token $STUDENT -Expect 404 | Out-Null
T -Name "Overflow INT64" -Method GET -Url "$BASE/courses/99999999999999999999" -Token $STUDENT -Expect 400 | Out-Null

# ====================== DELETE ======================
Write-Host "`n===== 11. DELETE =====" -ForegroundColor Cyan

T -Name "Delete nonexistent task" -Method DELETE -Url "$BASE/tasks/9999" -Token $ADMIN -Expect 404 | Out-Null
T -Name "Delete lesson 7" -Method DELETE -Url "$BASE/lessons/7" -Token $ADMIN -Expect 204 | Out-Null
T -Name "Delete nonexistent lesson" -Method DELETE -Url "$BASE/lessons/9999" -Token $ADMIN -Expect 404 | Out-Null
T -Name "Delete latest course" -Method DELETE -Url "$BASE/courses/5" -Token $ADMIN -Expect 204 | Out-Null
T -Name "Delete nonexistent course" -Method DELETE -Url "$BASE/courses/9999" -Token $ADMIN -Expect 404 | Out-Null

# ====================== RESULTS ======================
Write-Host "`n========================================" -ForegroundColor Cyan
Write-Host "  TOTAL=$total  PASS=$pass  FAIL=$fail" -ForegroundColor $(if ($fail -eq 0) {"Green"} else {"Yellow"})
Write-Host "========================================`n" -ForegroundColor Cyan

Remove-Item $tmpFile -ErrorAction SilentlyContinue
