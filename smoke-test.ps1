$ErrorActionPreference = "Continue"
$BASE = "http://localhost:8080/api/v1"

function Invoke-Api {
    param([string]$Method, [string]$Url, [string]$Body, [string]$Token)
    $headers = @{ "Content-Type" = "application/json" }
    if ($Token) { $headers["Authorization"] = "Bearer $Token" }
    $params = @{ Method = $Method; Uri = $Url; Headers = $headers; ContentType = "application/json" }
    if ($Body) { $params["Body"] = [System.Text.Encoding]::UTF8.GetBytes($Body) }
    try {
        $resp = Invoke-RestMethod @params
        return $resp
    } catch {
        Write-Host ("  ERROR: " + $_.Exception.Message) -ForegroundColor Red
        if ($_.ErrorDetails.Message) { Write-Host ("  BODY: " + $_.ErrorDetails.Message) -ForegroundColor Red }
        return $null
    }
}

Write-Host ""
Write-Host "========== 1. AUTH: Login as admin ==========" -ForegroundColor Cyan
$body = '{"email":"admin@example.com","password":"password123"}'
$admin = Invoke-Api -Method POST -Url "$BASE/auth/login" -Body $body
if ($admin) {
    Write-Host ("  OK: access_token = " + $admin.access_token.Substring(0,30) + "...") -ForegroundColor Green
    $ADMIN_TOKEN = $admin.access_token
} else { Write-Host "  FAIL" -ForegroundColor Red; exit 1 }

Write-Host ""
Write-Host "========== 2. AUTH: Login as student ==========" -ForegroundColor Cyan
$body = '{"email":"student@example.com","password":"password123"}'
$student = Invoke-Api -Method POST -Url "$BASE/auth/login" -Body $body
if ($student) {
    Write-Host ("  OK: access_token = " + $student.access_token.Substring(0,30) + "...") -ForegroundColor Green
    $STUDENT_TOKEN = $student.access_token
} else { Write-Host "  FAIL" -ForegroundColor Red; exit 1 }

Write-Host ""
Write-Host "========== 3. LIST COURSES ==========" -ForegroundColor Cyan
$courses = Invoke-Api -Method GET -Url "$BASE/courses" -Token $STUDENT_TOKEN
if ($courses) {
    Write-Host ("  OK: total = " + $courses.total + ", first = '" + $courses.data[0].title + "'") -ForegroundColor Green
} else { Write-Host "  FAIL" -ForegroundColor Red }

Write-Host ""
Write-Host "========== 4. GET COURSE 1 ==========" -ForegroundColor Cyan
$course = Invoke-Api -Method GET -Url "$BASE/courses/1" -Token $STUDENT_TOKEN
if ($course) {
    Write-Host ("  OK: '" + $course.title + "'") -ForegroundColor Green
} else { Write-Host "  FAIL" -ForegroundColor Red }

Write-Host ""
Write-Host "========== 5. LIST LESSONS (course 1) ==========" -ForegroundColor Cyan
$lessons = Invoke-Api -Method GET -Url "$BASE/courses/1/lessons" -Token $STUDENT_TOKEN
if ($lessons) {
    Write-Host ("  OK: total = " + $lessons.total) -ForegroundColor Green
} else { Write-Host "  FAIL" -ForegroundColor Red }

Write-Host ""
Write-Host "========== 6. LIST TASKS (lesson 1) ==========" -ForegroundColor Cyan
$tasks = Invoke-Api -Method GET -Url "$BASE/lessons/1/tasks" -Token $STUDENT_TOKEN
if ($tasks) {
    Write-Host ("  OK: total = " + $tasks.total) -ForegroundColor Green
    if ($tasks.data.Count -gt 0) {
        Write-Host ("  First task: '" + $tasks.data[0].title + "'") -ForegroundColor Green
    }
} else { Write-Host "  FAIL" -ForegroundColor Red }

Write-Host ""
Write-Host "========== 7. GET TASK 1 (with test_cases) ==========" -ForegroundColor Cyan
$task = Invoke-Api -Method GET -Url "$BASE/tasks/1" -Token $STUDENT_TOKEN
if ($task) {
    Write-Host ("  OK: '" + $task.title + "', test_cases count = " + $task.test_cases.Count) -ForegroundColor Green
    Write-Host ("  test_cases[0]: input='" + $task.test_cases[0].input + "' expected='" + $task.test_cases[0].expected + "'") -ForegroundColor Gray
} else { Write-Host "  FAIL" -ForegroundColor Red }

Write-Host ""
Write-Host "========== 8. SUBMIT SOLUTION (async) ==========" -ForegroundColor Cyan
$gocode = "package main`nimport `"fmt`"`nfunc main() {`n    var a, b int`n    fmt.Scan(&a, &b)`n    fmt.Println(a + b)`n}"
$submitBody = (@{ code = $gocode } | ConvertTo-Json -Compress)
$sub = Invoke-Api -Method POST -Url "$BASE/tasks/1/submissions" -Body $submitBody -Token $STUDENT_TOKEN
if ($sub) {
    Write-Host ("  OK: id=" + $sub.id + ", status='" + $sub.status + "'") -ForegroundColor Green
    $SUB_ID = $sub.id
} else { Write-Host "  FAIL" -ForegroundColor Red; exit 1 }

Write-Host ""
Write-Host "========== 9. POLL SUBMISSION (waiting for result) ==========" -ForegroundColor Cyan
$maxAttempts = 20
$attempt = 0
$finalSub = $null
while ($attempt -lt $maxAttempts) {
    Start-Sleep -Seconds 3
    $attempt++
    Write-Host ("  Poll attempt " + $attempt + "...") -ForegroundColor Gray
    $polled = Invoke-Api -Method GET -Url "$BASE/submissions/$SUB_ID" -Token $STUDENT_TOKEN
    if ($polled -and $polled.status -ne "pending") {
        $finalSub = $polled
        break
    }
}
if ($finalSub) {
    Write-Host ("  OK: status='" + $finalSub.status + "', score=" + $finalSub.score) -ForegroundColor Green
    Write-Host ("  result=" + $finalSub.result) -ForegroundColor Gray
} else {
    Write-Host "  TIMEOUT: submission still pending after $maxAttempts attempts" -ForegroundColor Red
}

Write-Host ""
Write-Host "========== 10. SUBMIT FAILING SOLUTION ==========" -ForegroundColor Cyan
$badcode = "package main`nimport `"fmt`"`nfunc main() {`n    fmt.Println(`"wrong answer`")`n}"
$badBody = (@{ code = $badcode } | ConvertTo-Json -Compress)
$badSub = Invoke-Api -Method POST -Url "$BASE/tasks/1/submissions" -Body $badBody -Token $STUDENT_TOKEN
if ($badSub) {
    Write-Host ("  OK: id=" + $badSub.id + ", status='" + $badSub.status + "'") -ForegroundColor Green
    $BAD_SUB_ID = $badSub.id
} else { Write-Host "  FAIL" -ForegroundColor Red; exit 1 }

Write-Host ""
Write-Host "========== 11. POLL FAILING SUBMISSION ==========" -ForegroundColor Cyan
$attempt = 0
$finalBad = $null
while ($attempt -lt $maxAttempts) {
    Start-Sleep -Seconds 3
    $attempt++
    Write-Host ("  Poll attempt " + $attempt + "...") -ForegroundColor Gray
    $polled = Invoke-Api -Method GET -Url "$BASE/submissions/$BAD_SUB_ID" -Token $STUDENT_TOKEN
    if ($polled -and $polled.status -ne "pending") {
        $finalBad = $polled
        break
    }
}
if ($finalBad) {
    Write-Host ("  OK: status='" + $finalBad.status + "', score=" + $finalBad.score) -ForegroundColor Green
} else {
    Write-Host "  TIMEOUT: submission still pending" -ForegroundColor Red
}

Write-Host ""
Write-Host "========== 12. LIST SUBMISSIONS (task 1) ==========" -ForegroundColor Cyan
$subList = Invoke-Api -Method GET -Url "$BASE/tasks/1/submissions" -Token $STUDENT_TOKEN
if ($subList) {
    Write-Host ("  OK: total = " + $subList.total) -ForegroundColor Green
    foreach ($s in $subList.data) {
        Write-Host ("    id=" + $s.id + " status=" + $s.status + " score=" + $s.score) -ForegroundColor Gray
    }
} else { Write-Host "  FAIL" -ForegroundColor Red }

Write-Host ""
Write-Host "========== 13. PROGRESS ==========" -ForegroundColor Cyan
$prog = Invoke-Api -Method GET -Url "$BASE/courses/1/progress" -Token $STUDENT_TOKEN
if ($prog) {
    Write-Host ("  OK: " + ($prog | ConvertTo-Json -Compress)) -ForegroundColor Green
} else { Write-Host "  FAIL (may be expected if no completed lessons)" -ForegroundColor Yellow }

Write-Host ""
Write-Host "========== ALL TESTS COMPLETE ==========" -ForegroundColor Cyan
