#!/bin/bash

set -euo pipefail

INPUT=$(cat)

ISOLATION=$(
  echo "$INPUT" | jq -r '.tool_input.isolation // ""'
)

if [ "$ISOLATION" = "worktree" ]; then
  echo "ERROR: isolation=worktree is forbidden"
  exit 2
fi

exit 0
