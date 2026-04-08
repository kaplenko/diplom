-- Seed data for the diplom application.
-- Run from the project root: go run ./cmd/seed
--
-- Password for all seed users: password123
-- The hash placeholder {{PASSWORD_HASH}} is replaced at runtime by the Go seed script.

BEGIN;

-- Wipe existing data so the seed is idempotent
TRUNCATE TABLE progress, submissions, tasks, lessons, courses, users
    RESTART IDENTITY CASCADE;

-- ============================================================
-- Users
-- ============================================================
INSERT INTO users (email, password_hash, name, role) VALUES
    ('admin@example.com',   '{{PASSWORD_HASH}}', 'Admin User',   'admin'),
    ('student@example.com', '{{PASSWORD_HASH}}', 'Student User', 'student');

-- ============================================================
-- Courses  (created_by = 1 → admin)
-- ============================================================
INSERT INTO courses (title, description, created_by) VALUES
    ('Introduction to Go',        'Learn the basics of the Go programming language, from Hello World to control flow.', 1),
    ('Web Development with Gin',  'Build REST APIs and web applications using the Gin framework.',                       1);

-- ============================================================
-- Lessons — "Introduction to Go" (course_id = 1)
-- ============================================================
INSERT INTO lessons (course_id, title, content, order_index) VALUES
    (1, 'Hello World',          'In this lesson you will write your first Go program and learn about package main and fmt.Println.', 1),
    (1, 'Variables and Types',  'Go is a statically typed language. Learn how to declare variables and use basic types.',            2),
    (1, 'Control Flow',         'Learn about if/else statements, for loops, and switch in Go.',                                      3);

-- ============================================================
-- Lessons — "Web Development with Gin" (course_id = 2)
-- ============================================================
INSERT INTO lessons (course_id, title, content, order_index) VALUES
    (2, 'Setting Up Gin',        'Install the Gin package and create your first HTTP server with a single route.',         1),
    (2, 'Routing and Handlers',  'Define route groups, path parameters, and write handler functions that return JSON.',     2);

-- ============================================================
-- Tasks — Hello World (lesson_id = 1)
-- ============================================================
INSERT INTO tasks (lesson_id, title, description, initial_code, test_cases, difficulty) VALUES
(1,
 'Print Hello World',
 'Write a program that prints "Hello, World!" to standard output.',
 E'package main\n\nfunc main() {\n\t// your code here\n}',
 '[{"input": "", "expected": "Hello, World!"}]',
 'easy');

-- ============================================================
-- Tasks — Variables and Types (lesson_id = 2)
-- ============================================================
INSERT INTO tasks (lesson_id, title, description, initial_code, test_cases, difficulty) VALUES
(2,
 'Sum Two Numbers',
 'Write a function that accepts two integers and returns their sum.',
 E'package main\n\nfunc sum(a, b int) int {\n\treturn 0\n}',
 '[{"input": "2 3", "expected": "5"}, {"input": "-1 1", "expected": "0"}]',
 'easy'),
(2,
 'String Length',
 'Write a function that returns the length of the given string.',
 E'package main\n\nfunc strLen(s string) int {\n\treturn 0\n}',
 '[{"input": "hello", "expected": "5"}, {"input": "", "expected": "0"}]',
 'easy');

-- ============================================================
-- Tasks — Control Flow (lesson_id = 3)
-- ============================================================
INSERT INTO tasks (lesson_id, title, description, initial_code, test_cases, difficulty) VALUES
(3,
 'FizzBuzz',
 'Return "Fizz" if divisible by 3, "Buzz" if by 5, "FizzBuzz" if by both, otherwise the number as a string.',
 E'package main\n\nfunc fizzBuzz(n int) string {\n\treturn ""\n}',
 '[{"input": "3", "expected": "Fizz"}, {"input": "5", "expected": "Buzz"}, {"input": "15", "expected": "FizzBuzz"}, {"input": "7", "expected": "7"}]',
 'medium');

-- ============================================================
-- Tasks — Setting Up Gin (lesson_id = 4)
-- ============================================================
INSERT INTO tasks (lesson_id, title, description, initial_code, test_cases, difficulty) VALUES
(4,
 'Create a Ping Endpoint',
 'Create a GET /ping endpoint that responds with {"message": "pong"}.',
 E'package main\n\nimport "github.com/gin-gonic/gin"\n\nfunc main() {\n\tr := gin.Default()\n\t// add route here\n\tr.Run()\n}',
 '[{"input": "GET /ping", "expected": "{\"message\":\"pong\"}"}]',
 'easy');

-- ============================================================
-- Tasks — Routing and Handlers (lesson_id = 5)
-- ============================================================
INSERT INTO tasks (lesson_id, title, description, initial_code, test_cases, difficulty) VALUES
(5,
 'Path Parameters',
 'Create a GET /users/:id endpoint that returns the id from the URL path as JSON.',
 E'package main\n\nimport "github.com/gin-gonic/gin"\n\nfunc main() {\n\tr := gin.Default()\n\t// add route here\n\tr.Run()\n}',
 '[{"input": "GET /users/42", "expected": "{\"id\":\"42\"}"}]',
 'medium'),
(5,
 'Query Parameters',
 'Create a GET /search endpoint that reads a "q" query parameter and echoes it back as JSON.',
 E'package main\n\nimport "github.com/gin-gonic/gin"\n\nfunc main() {\n\tr := gin.Default()\n\t// add route here\n\tr.Run()\n}',
 '[{"input": "GET /search?q=golang", "expected": "{\"query\":\"golang\"}"}]',
 'medium');

-- ============================================================
-- Progress
--
-- Schema: one row per (user, lesson). The course_id is denormalised
-- for fast per-course queries.
--
-- Seed data reference:
--   Course 1 "Introduction to Go"      → lessons 1, 2, 3  (3 total)
--   Course 2 "Web Development with Gin" → lessons 4, 5     (2 total)
--   User 1 = admin,  User 2 = student
-- ============================================================

-- Admin (user_id = 1): completed all of Course 1 (3/3), lesson 4 of Course 2 (1/2)
INSERT INTO progress (user_id, course_id, lesson_id, completed, completed_at) VALUES
    (1, 1, 1, TRUE,  '2026-03-01 10:00:00+00'),
    (1, 1, 2, TRUE,  '2026-03-02 14:30:00+00'),
    (1, 1, 3, TRUE,  '2026-03-03 09:15:00+00'),
    (1, 2, 4, TRUE,  '2026-03-10 11:00:00+00'),
    (1, 2, 5, FALSE, NULL);

-- Student (user_id = 2): completed 2/3 of Course 1, started Course 2 (1/2)
INSERT INTO progress (user_id, course_id, lesson_id, completed, completed_at) VALUES
    (2, 1, 1, TRUE,  '2026-03-15 16:00:00+00'),
    (2, 1, 2, TRUE,  '2026-03-16 18:45:00+00'),
    (2, 1, 3, FALSE, NULL),
    (2, 2, 4, TRUE,  '2026-03-20 12:30:00+00'),
    (2, 2, 5, FALSE, NULL);

COMMIT;
