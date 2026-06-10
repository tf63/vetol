#!/bin/bash

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

PASSED=0
FAILED=0

run_test() {
  local test_name="$1"
  local expected="$2"
  shift 2

  local output

  output=$(
    cd "$PROJECT_ROOT" &&
    go run ./cmd/vetol "$@" 2>&1 || true
  )

  if echo "$output" | grep -qxF "$expected"; then
    echo -e "${GREEN}✓${NC} $test_name"
    ((PASSED++))
  else
    echo -e "${RED}✗${NC} $test_name"
    ((FAILED++))
  fi
}

echo "Running vetol tests..."
echo

# ============================================================
# Test data: Allowlist mode
# "rules": ["pwd", "echo hello", "grep -r", "ls -la"]
# ============================================================
ALLOWLIST_ALLOW=(
  "pwd"
  "pwd -P"
  "echo hello"
  "echo hello world"
  "echo hello \$(pwd)"
  "echo hello \`pwd\`"
  "grep -r"
  "grep -r -l"
  "grep -l -r"
  "grep -rl"
  "grep -r /tmp"
  "ls -la"
  "ls -la --color=auto"
  "ls -la /tmp"
  "pwd && echo hello"
  "pwd || echo hello"
  "pwd; echo hello"
  "(pwd; echo hello)"
  "pwd && echo hello || grep -r test"
  "echo hello"
  "echo hello test"
)

ALLOWLIST_DENY=(
  "cat"
  "cat README.md"
  "echo"
  "echo \"hello world\""
  "grep"
  "grep -l"
  "grep /tmp"
  "ls"
  "ls -l"
  "ls -a"
  "ls --color=auto"
  "ls /tmp"
  "pwd && echo"
  "pwd || echo"
  "pwd; echo"
  "(pwd; echo)"
  "echo && pwd"
  "echo \$(pwd)"
  "echo \`pwd\`"
  "pwd && cat || echo"
)

# ============================================================
# Test data: Denylist mode
# "rules": ["pwd", "echo hello", "grep -r", "ls -la"]
# ============================================================
DENYLIST_ALLOW=(
  "ls"
  "echo"
  "cat"
  "grep"
  "grep -l"
  "ls -l"
  "ls -a"
  "ls --color=auto"
  "ls && echo"
  "ls || echo"
  "ls; echo"
  "(ls; echo)"
  "echo && ls"
  "ls && cat || echo"
)

DENYLIST_DENY=(
  "pwd"
  "pwd -P"
  "echo hello"
  "echo hello world"
  "echo hello \$(pwd)"
  "grep -r"
  "grep -r -a"
  "grep -l -r"
  "ls -la"
  "ls \$(pwd)"
  "ls \`pwd\`"
  "pwd && ls"
  "pwd || ls"
  "pwd test"
  "pwd | grep test"
  "pwd arg1 arg2"

)

# ============================================================
# Run Allowlist tests
# ============================================================

echo "=== Allowlist Mode Tests ==="

for cmd in "${ALLOWLIST_ALLOW[@]}"; do
  run_test "allowlist: allowed ($cmd)" \
    "ALLOW" \
    --config testdata/allowlist.json "$cmd"
done

for cmd in "${ALLOWLIST_DENY[@]}"; do
  run_test "allowlist: denied ($cmd)" \
    "DENY" \
    --config testdata/allowlist.json "$cmd"
done

# ============================================================
# Run Denylist tests
# ============================================================

echo
echo "=== Denylist Mode Tests ==="

for cmd in "${DENYLIST_ALLOW[@]}"; do
  run_test "denylist: allowed ($cmd)" \
    "ALLOW" \
    --config testdata/denylist.json "$cmd"
done

for cmd in "${DENYLIST_DENY[@]}"; do
  run_test "denylist: denied ($cmd)" \
    "DENY" \
    --config testdata/denylist.json "$cmd"
done

# ============================================================
# Allowlist with Exclude Tests
# ============================================================

echo
echo "=== Allowlist with Exclude Tests ==="

# Test cases with exclude flag
run_test "allowlist with exclude: grep -r (allowed, no exclude match)" \
  "ALLOW" \
  --config testdata/allowlist_with_exclude.json "grep -r pattern"

run_test "allowlist with exclude: grep -r -q (denied, exclude match)" \
  "DENY" \
  --config testdata/allowlist_with_exclude.json "grep -r -q pattern"

run_test "allowlist with exclude: grep -rq (denied, exclude match)" \
  "DENY" \
  --config testdata/allowlist_with_exclude.json "grep -rq pattern"

run_test "allowlist with exclude: grep -qr (denied, exclude match)" \
  "DENY" \
  --config testdata/allowlist_with_exclude.json "grep -qr pattern"

run_test "allowlist with exclude: ls -la (allowed, no exclude match)" \
  "ALLOW" \
  --config testdata/allowlist_with_exclude.json "ls -la"

run_test "allowlist with exclude: ls -la -i (denied, exclude match)" \
  "DENY" \
  --config testdata/allowlist_with_exclude.json "ls -la -i"

run_test "allowlist with exclude: echo hello (allowed)" \
  "ALLOW" \
  --config testdata/allowlist_with_exclude.json "echo hello"

run_test "allowlist with exclude: pwd (allowed)" \
  "ALLOW" \
  --config testdata/allowlist_with_exclude.json "pwd"

# ============================================================
# Blacklist with Exclude Tests
# ============================================================

echo
echo "=== Blacklist with Exclude Tests ==="

run_test "denylist with exclude: pwd (denied, rule matches)" \
  "DENY" \
  --config testdata/denylist_with_exclude.json "pwd"

run_test "denylist with exclude: pwd -P (allowed, exclude prevents match)" \
  "ALLOW" \
  --config testdata/denylist_with_exclude.json "pwd -P"

run_test "denylist with exclude: grep -r (denied, rule matches)" \
  "DENY" \
  --config testdata/denylist_with_exclude.json "grep -r pattern"

run_test "denylist with exclude: grep -r -q (allowed, exclude prevents match)" \
  "ALLOW" \
  --config testdata/denylist_with_exclude.json "grep -r -q pattern"

run_test "denylist with exclude: grep -rq (allowed, exclude prevents match)" \
  "ALLOW" \
  --config testdata/denylist_with_exclude.json "grep -rq pattern"

run_test "denylist with exclude: grep -qr (allowed, exclude prevents match)" \
  "ALLOW" \
  --config testdata/denylist_with_exclude.json "grep -qr pattern"

run_test "denylist with exclude: ls -la (denied, rule matches)" \
  "DENY" \
  --config testdata/denylist_with_exclude.json "ls -la"

run_test "denylist with exclude: ls -la -i (allowed, exclude prevents match)" \
  "ALLOW" \
  --config testdata/denylist_with_exclude.json "ls -la -i"

run_test "denylist with exclude: echo hello (denied, rule matches)" \
  "DENY" \
  --config testdata/denylist_with_exclude.json "echo hello"

run_test "denylist with exclude: cat (allowed, no match)" \
  "ALLOW" \
  --config testdata/denylist_with_exclude.json "cat file.txt"

# ============================================================
# Long Flag with Value Tests
# ============================================================

echo
echo "=== Long Flag with Value Tests ==="

run_test "long flag: grep -r --color (rule has both -r and --color)" \
  "ALLOW" \
  --config testdata/long_flag.json "grep -r --color=auto pattern"

run_test "long flag: grep -r --color=never (rule matches --color prefix)" \
  "ALLOW" \
  --config testdata/long_flag.json "grep -r --color=never pattern"

run_test "long flag: grep --color=auto pattern (missing required -r)" \
  "DENY" \
  --config testdata/long_flag.json "grep --color=auto pattern"

run_test "long flag: ls --color=auto (rule has --color)" \
  "ALLOW" \
  --config testdata/long_flag.json "ls --color=auto"

run_test "long flag: ls --color (exact match)" \
  "ALLOW" \
  --config testdata/long_flag.json "ls --color"

run_test "long flag: ls (missing --color)" \
  "DENY" \
  --config testdata/long_flag.json "ls"

# ============================================================
# Error Cases Tests
# ============================================================

echo
echo "=== Error Cases Tests ==="

run_test "CLI error: missing --config" \
  "ERROR: --config is required" \
  "pwd"

run_test "CLI error: missing command string" \
  "ERROR: command string is required" \
  --config testdata/allowlist.json

# ============================================================
# Summary
# ============================================================

echo
echo "=== Summary ==="
echo -e "Passed: ${GREEN}$PASSED${NC}"
echo -e "Failed: ${RED}$FAILED${NC}"

if [ $FAILED -eq 0 ]; then
  exit 0
else
  exit 1
fi
