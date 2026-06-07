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
  "grep -r /tmp"
  "ls -la"
  "ls -la --color=auto"
  "ls -la /tmp"
  "pwd && echo hello"
  "pwd || echo hello"
  "pwd; echo hello"
  "(pwd; echo hello)"
  "pwd && echo hello || grep -r test"
  "echo hello > a.txt"
  "echo hello > /dev/null 2>&1"
)

ALLOWLIST_DENY=(
  "cat"
  "cat README.md"
  "echo"
  "echo \"hello world\""
  "grep"
  "grep -l"
  "grep -l -r"
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
  "grep -l -r"
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
  "ls -la"
  "ls \$(pwd)"
  "ls \`pwd\`"
  "pwd && ls"
  "pwd || ls"
  "pwd > /dev/null"
  "pwd | grep test > /dev/null"
  "pwd > /dev/null 2>&1"

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
# CLI Arguments Tests
# ============================================================

echo
echo "=== Command Line Arguments Tests ==="

run_test "CLI: allowlist mode with pwd" \
  "ALLOW" \
  -m allowlist -r pwd "pwd"

run_test "CLI: denylist mode with pwd" \
  "DENY" \
  -m denylist -r pwd "pwd"

run_test "CLI: multiple rules (comma-separated)" \
  "ALLOW" \
  -m allowlist -r "pwd,cat" "pwd"

run_test "CLI: multiple rules with denied command" \
  "DENY" \
  -m allowlist -r "pwd,cat" "ls"

run_test "CLI: multiple rules in denylist" \
  "ALLOW" \
  -m denylist -r "rm,dd" "ls"

# ============================================================
# Error Cases Tests
# ============================================================

echo
echo "=== Error Cases Tests ==="

run_test "CLI error: missing --rules" \
  "ERROR: either --config or both --mode and --rules are required" \
  -m allowlist "pwd"

run_test "CLI error: missing --mode" \
  "ERROR: either --config or both --mode and --rules are required" \
  -r pwd "pwd"

run_test "CLI error: missing both --mode and --rules" \
  "ERROR: either --config or both --mode and --rules are required" \
  "pwd"

run_test "CLI error: missing command string" \
  "ERROR: command string is required" \
  -m allowlist -r pwd

run_test "CLI error: invalid mode" \
  "ERROR: invalid mode [mode invalid]" \
  -m invalid -r pwd "pwd"

run_test "CLI error: mixing --config with --mode" \
  "ERROR: cannot mix --config with --mode/--rules" \
  --config testdata/allowlist.json -m allowlist "pwd"

run_test "CLI error: mixing --config with --rules" \
  "ERROR: cannot mix --config with --mode/--rules" \
  --config testdata/allowlist.json -r pwd "pwd"

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
