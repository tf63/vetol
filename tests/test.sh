#!/bin/bash

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

PASSED=0
FAILED=0

# Test function
run_test() {
  local test_name="$1"
  local args="$2"
  local expected="$3"

  output=$(cd "$PROJECT_ROOT" && eval "go run ./cmd/vetol $args" 2>&1 || true)

  if echo "$output" | grep -q "$expected"; then
    echo -e "${GREEN}✓${NC} $test_name"
    ((PASSED++))
  else
    echo -e "${RED}✗${NC} $test_name"
    echo "  Expected: $expected"
    echo "  Got: $output"
    ((FAILED++))
  fi
}

echo "Running vetol tests..."
echo

# Whitelist tests
echo "=== Whitelist Tests ==="
run_test "whitelist: allowed command (pwd)" \
  "check --config testdata/whitelist.json \"pwd\"" \
  "VALID"

run_test "whitelist: allowed multicommand (echo hello)" \
  "check --config testdata/whitelist.json 'echo hello'" \
  "VALID"

run_test "whitelist: denied single command (echo)" \
  "check --config testdata/whitelist.json \"echo\"" \
  "INVALID"

run_test "whitelist: denied command (ls)" \
  "check --config testdata/whitelist.json \"ls\"" \
  "INVALID"

echo
echo "=== Blacklist Tests ==="

run_test "blacklist: allowed command (ls)" \
  "check --config testdata/blacklist.json \"ls\"" \
  "VALID"

run_test "blacklist: allowed single command (echo)" \
  "check --config testdata/blacklist.json \"echo\"" \
  "VALID"

run_test "blacklist: denied command (pwd)" \
  "check --config testdata/blacklist.json \"pwd\"" \
  "INVALID"

run_test "blacklist: denied multicommand (echo hello)" \
  "check --config testdata/blacklist.json 'echo hello'" \
  "INVALID"

echo
echo "=== Complex Command Tests ==="

# AND operator (&&)
run_test "whitelist: && (pwd && cat)" \
  "check --config testdata/whitelist.json 'pwd && cat'" \
  "VALID"

run_test "whitelist: && with denied (pwd && echo)" \
  "check --config testdata/whitelist.json 'pwd && echo'" \
  "INVALID"

run_test "blacklist: && with denied (pwd && ls)" \
  "check --config testdata/blacklist.json 'pwd && ls'" \
  "INVALID"

# OR operator (||)
run_test "whitelist: || (pwd || cat)" \
  "check --config testdata/whitelist.json 'pwd || cat'" \
  "VALID"

run_test "whitelist: || with denied (pwd || echo)" \
  "check --config testdata/whitelist.json 'pwd || echo'" \
  "INVALID"

run_test "blacklist: || with denied (pwd || ls)" \
  "check --config testdata/blacklist.json 'pwd || ls'" \
  "INVALID"

# Semicolon (;)
run_test "whitelist: ; (pwd; cat)" \
  "check --config testdata/whitelist.json 'pwd; cat'" \
  "VALID"

run_test "whitelist: ; with denied (pwd; echo)" \
  "check --config testdata/whitelist.json 'pwd; echo'" \
  "INVALID"

# Pipe (|)
run_test "whitelist: | (pwd | grep)" \
  "check --config testdata/whitelist.json 'pwd | grep'" \
  "VALID"

run_test "whitelist: | with denied (echo | grep)" \
  "check --config testdata/whitelist.json 'echo | grep'" \
  "INVALID"

# Command substitution with $()
run_test "whitelist: \$() (echo \$(pwd))" \
  "check --config testdata/whitelist.json 'echo \$(pwd)'" \
  "INVALID"

run_test "blacklist: \$() with allowed (ls \$(pwd))" \
  "check --config testdata/blacklist.json 'ls \$(pwd)'" \
  "VALID"

# Backticks
run_test "whitelist: backticks (echo \\\`pwd\\\`)" \
  "check --config testdata/whitelist.json 'echo \`pwd\`'" \
  "INVALID"

run_test "blacklist: backticks (ls \\\`pwd\\\`)" \
  "check --config testdata/blacklist.json 'ls \`pwd\`'" \
  "VALID"

# Subshell with ()
run_test "whitelist: subshell (pwd; cat)" \
  "check --config testdata/whitelist.json '(pwd; cat)'" \
  "VALID"

run_test "whitelist: subshell with denied (pwd; echo)" \
  "check --config testdata/whitelist.json '(pwd; echo)'" \
  "INVALID"

# Redirect
run_test "whitelist: redirect (cat > /dev/null)" \
  "check --config testdata/whitelist.json 'cat > /dev/null'" \
  "VALID"

run_test "blacklist: redirect with denied (pwd > /dev/null)" \
  "check --config testdata/blacklist.json 'pwd > /dev/null'" \
  "INVALID"

# Complex pipeline and redirect
run_test "whitelist: complex pipeline (pwd | grep test > /dev/null)" \
  "check --config testdata/whitelist.json 'pwd | grep test > /dev/null'" \
  "VALID"

run_test "blacklist: complex pipeline (pwd | grep test > /dev/null)" \
  "check --config testdata/blacklist.json 'pwd | grep test > /dev/null'" \
  "INVALID"

# AND/OR combination
run_test "whitelist: && || (pwd && cat || echo)" \
  "check --config testdata/whitelist.json 'pwd && cat || echo'" \
  "INVALID"

run_test "blacklist: && || (ls && cat || echo)" \
  "check --config testdata/blacklist.json 'ls && cat || echo'" \
  "VALID"

# Multiple redirects
run_test "whitelist: multiple redirects (pwd > /dev/null 2>&1)" \
  "check --config testdata/whitelist.json 'pwd > /dev/null 2>&1'" \
  "VALID"

run_test "blacklist: multiple redirects (pwd > /dev/null 2>&1)" \
  "check --config testdata/blacklist.json 'pwd > /dev/null 2>&1'" \
  "INVALID"

echo
echo "=== Command with Options Tests ==="

run_test "whitelist: allowed option (grep -r)" \
  "check --config testdata/whitelist.json \"grep -r\"" \
  "VALID"

run_test "whitelist: allowed option (ls -la)" \
  "check --config testdata/whitelist.json \"ls -la\"" \
  "VALID"

run_test "whitelist: denied option (grep -l)" \
  "check --config testdata/whitelist.json \"grep -l\"" \
  "INVALID"

run_test "whitelist: denied command (grep)" \
  "check --config testdata/whitelist.json \"grep\"" \
  "INVALID"

run_test "whitelist: denied option (ls -l)" \
  "check --config testdata/whitelist.json \"ls -l\"" \
  "INVALID"

run_test "whitelist: allowed command without option (pwd)" \
  "check --config testdata/whitelist.json \"pwd\"" \
  "VALID"

run_test "whitelist: allowed command (echo hello)" \
  "check --config testdata/whitelist.json \"echo hello\"" \
  "VALID"

run_test "whitelist: denied command (echo test)" \
  "check --config testdata/whitelist.json \"echo test\"" \
  "INVALID"

run_test "blacklist: denied option (grep -r)" \
  "check --config testdata/blacklist.json \"grep -r\"" \
  "INVALID"

run_test "blacklist: denied option (ls -la)" \
  "check --config testdata/blacklist.json \"ls -la\"" \
  "INVALID"

run_test "blacklist: allowed option (grep -l)" \
  "check --config testdata/blacklist.json \"grep -l\"" \
  "VALID"

run_test "blacklist: allowed command (grep)" \
  "check --config testdata/blacklist.json \"grep\"" \
  "VALID"

run_test "blacklist: allowed option (ls -l)" \
  "check --config testdata/blacklist.json \"ls -l\"" \
  "VALID"

run_test "blacklist: denied command (pwd)" \
  "check --config testdata/blacklist.json \"pwd\"" \
  "INVALID"

run_test "blacklist: denied command (echo hello)" \
  "check --config testdata/blacklist.json \"echo hello\"" \
  "INVALID"

run_test "blacklist: allowed command (cat)" \
  "check --config testdata/blacklist.json \"cat\"" \
  "VALID"

echo
echo "=== Command Line Arguments Tests ==="

# --mode and --rules
run_test "CLI: whitelist mode with pwd" \
  "check -m whitelist -r pwd \"pwd\"" \
  "VALID"

run_test "CLI: blacklist mode with pwd" \
  "check -m blacklist -r pwd \"pwd\"" \
  "INVALID"

run_test "CLI: long flags --mode --rules" \
  "check --mode whitelist --rules pwd \"pwd\"" \
  "VALID"

# Multiple rules (comma-separated)
run_test "CLI: multiple rules (comma-separated)" \
  "check -m whitelist -r \"pwd,cat\" \"pwd\"" \
  "VALID"

run_test "CLI: multiple rules with denied command" \
  "check -m whitelist -r \"pwd,cat\" \"ls\"" \
  "INVALID"

run_test "CLI: multiple rules in blacklist" \
  "check -m blacklist -r \"rm,dd\" \"ls\"" \
  "VALID"

# Error cases - missing required options
run_test "CLI: error - missing --rules" \
  "check -m whitelist \"pwd\"" \
  "required"

run_test "CLI: error - missing --mode" \
  "check -r pwd \"pwd\"" \
  "required"

run_test "CLI: error - missing both --mode and --rules" \
  "check \"pwd\"" \
  "required"

# Error cases - conflicting options
run_test "CLI: error - mixing --config with --mode" \
  "check --config testdata/whitelist.json -m whitelist \"pwd\"" \
  "cannot mix"

run_test "CLI: error - mixing --config with --rules" \
  "check --config testdata/whitelist.json -r pwd \"pwd\"" \
  "cannot mix"

# Error cases - missing command string
run_test "CLI: error - missing command string" \
  "check -m whitelist -r pwd" \
  "required"

run_test "CLI: error - no subcommand" \
  "check" \
  "Usage"

run_test "CLI: error - invalid subcommand" \
  "validate -m whitelist -r pwd \"pwd\"" \
  "unknown subcommand"

# Error cases - invalid mode
run_test "CLI: error - invalid mode" \
  "check -m invalid -r pwd \"pwd\"" \
  "invalid mode"

echo
echo "=== Summary ==="
echo -e "Passed: ${GREEN}$PASSED${NC}"
echo -e "Failed: ${RED}$FAILED${NC}"

if [ $FAILED -eq 0 ]; then
  exit 0
else
  exit 1
fi
