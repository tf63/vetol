#!/bin/bash

INPUT=$(cat)
COMMAND=$(echo "$INPUT" | jq -r '.tool_input.command // empty')

BLACKLIST=(
  # File deletion
  '^rm( .*)?$'
  '^unlink( .*)?$'
  '^shred( .*)?$'

  # Privilege escalation
  '^sudo( .*)?$'
  '^su( .*)?$'

  # System shutdown/reboot
  '^reboot( .*)?$'
  '^shutdown( .*)?$'
  '^halt( .*)?$'
  '^poweroff( .*)?$'

  # Dangerous disk operations
  '^dd( .*)?$'
  '^mkfs(\..*)?( .*)?$'
  '^fdisk( .*)?$'
  '^parted( .*)?$'

  # Git history destruction
  '^git reset --hard( .*)?$'
  '^git clean -f.*$'
  '^git push --force( .*)?$'
  '^git push -f( .*)?$'

  # Terraform
  '^terraform destroy( .*)?$'

  # Kubernetes
  '^kubectl delete( .*)?$'
  '^kubectl drain( .*)?$'
  '^helm uninstall( .*)?$'

  # Docker cleanup
  '^docker system prune( .*)?$'
  '^docker volume rm( .*)?$'
  '^docker rm -f( .*)?$'
  '^docker compose down -v( .*)?$'

  # Download and execute
  '^curl .*\|.*$'
  '^wget .*\|.*$'
)

ALLOWED=true

for pattern in "${BLACKLIST[@]}"; do
  echo "TEST: $pattern" >&2
  if echo "$COMMAND" | grep -Eq "$pattern"; then
    ALLOWED=false
    break
  fi
done

if [ "$ALLOWED" = false ]; then
  jq -n '{
    hookSpecificOutput: {
      hookEventName: "PreToolUse",
      permissionDecision: "deny",
      permissionDecisionReason: "Command not allowed by blacklist policy"
    }
  }'
  exit 2
else
  exit 0
fi
