#!/bin/bash

INPUT=$(cat)
COMMAND=$(echo "$INPUT" | jq -r '.tool_input.command // empty')

VETOL_CONFIG_FILE="./vetol.json"

if [ ! -f "$VETOL_CONFIG_FILE" ]; then
  jq -n --arg cfg "$VETOL_CONFIG_FILE" '{
    hookSpecificOutput: {
      hookEventName: "PreToolUse",
      permissionDecision: "deny",
      permissionDecisionReason: ("vetol config not found: " + $cfg)
    }
  }'
  exit 2
fi

if vetol --config "$VETOL_CONFIG_FILE" "$COMMAND" >/dev/null 2>&1; then
  exit 0
fi

jq -n --arg cmd "$COMMAND" '{
  hookSpecificOutput: {
    hookEventName: "PreToolUse",
    permissionDecision: "deny",
    permissionDecisionReason: ("Denied by vetol: " + $cmd)
  }
}'
exit 2
