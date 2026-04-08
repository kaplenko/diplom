#!/usr/bin/env bash
set -euo pipefail

BASE_URL="${1:-http://localhost:8080}"
API="${BASE_URL}/api/v1"

TIMESTAMP=$(date +%s)
EMAIL="testuser_${TIMESTAMP}@example.com"
PASSWORD="TestPass123"
NAME="Smoke Test User"

echo "=== Smoke Test ==="
echo "Target: ${API}"
echo ""

# 1. Register
echo ">> POST /auth/register"
REG_RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "${API}/auth/register" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"${EMAIL}\",\"password\":\"${PASSWORD}\",\"name\":\"${NAME}\"}")

REG_BODY=$(echo "$REG_RESPONSE" | head -n -1)
REG_STATUS=$(echo "$REG_RESPONSE" | tail -n 1)

echo "   Status: ${REG_STATUS}"
echo "   Body:   ${REG_BODY}"

if [ "$REG_STATUS" -ne 201 ]; then
  echo "FAIL: Expected 201, got ${REG_STATUS}"
  exit 1
fi
echo "   OK"
echo ""

# 2. Login
echo ">> POST /auth/login"
LOGIN_RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "${API}/auth/login" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"${EMAIL}\",\"password\":\"${PASSWORD}\"}")

LOGIN_BODY=$(echo "$LOGIN_RESPONSE" | head -n -1)
LOGIN_STATUS=$(echo "$LOGIN_RESPONSE" | tail -n 1)

echo "   Status: ${LOGIN_STATUS}"
echo "   Body:   ${LOGIN_BODY}"

if [ "$LOGIN_STATUS" -ne 200 ]; then
  echo "FAIL: Expected 200, got ${LOGIN_STATUS}"
  exit 1
fi

ACCESS_TOKEN=$(echo "$LOGIN_BODY" | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4)

if [ -z "$ACCESS_TOKEN" ]; then
  echo "FAIL: Could not extract access_token"
  exit 1
fi
echo "   Token:  ${ACCESS_TOKEN:0:20}..."
echo "   OK"
echo ""

# 3. Get current user (protected route)
echo ">> GET /users/me"
ME_RESPONSE=$(curl -s -w "\n%{http_code}" -X GET "${API}/users/me" \
  -H "Authorization: Bearer ${ACCESS_TOKEN}")

ME_BODY=$(echo "$ME_RESPONSE" | head -n -1)
ME_STATUS=$(echo "$ME_RESPONSE" | tail -n 1)

echo "   Status: ${ME_STATUS}"
echo "   Body:   ${ME_BODY}"

if [ "$ME_STATUS" -ne 200 ]; then
  echo "FAIL: Expected 200, got ${ME_STATUS}"
  exit 1
fi
echo "   OK"
echo ""

# 4. List courses (protected route)
echo ">> GET /courses"
COURSES_RESPONSE=$(curl -s -w "\n%{http_code}" -X GET "${API}/courses" \
  -H "Authorization: Bearer ${ACCESS_TOKEN}")

COURSES_BODY=$(echo "$COURSES_RESPONSE" | head -n -1)
COURSES_STATUS=$(echo "$COURSES_RESPONSE" | tail -n 1)

echo "   Status: ${COURSES_STATUS}"
echo "   Body:   ${COURSES_BODY}"

if [ "$COURSES_STATUS" -ne 200 ]; then
  echo "FAIL: Expected 200, got ${COURSES_STATUS}"
  exit 1
fi
echo "   OK"
echo ""

echo "=== All smoke tests passed ==="
