#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat <<EOF
Usage: $(basename "$0") <command> [flags]

Commands:
  snip   Extract code examples from Bluehawk snippets
  copy   Copy sanitized project files

Example (non‐interactive):
  $(basename "$0") copy --ignore="tests/*.go" --rename='{"old.md":"new.md"}'
EOF
  exit 1
}

# defaults
CMD=""
PROJECT=$(git rev-parse --show-toplevel)
INPUT_DIR="$PROJECT/usage-examples/go/atlas-sdk-go"
OUTPUT_DIR="$PROJECT/generated-usage-examples/go/atlas-sdk-go/"
STATE=""
IGNORE_PATTERNS=(
  "internal_*.*"
  "scripts/"
  "*_test.go"
  ".env"
  "*.gz"
  "*.log"
  "./logs"
)

# ─── Interactive mode ────────────────────────────────────────────────────────
if [[ $# -eq 0 ]]; then
  echo "=== Run Bluehawk ==="

  # 1) pick snip or copy
  while true; do
    read -rp "Enter command (snip/copy): " CMD
    [[ "$CMD" == "snip" || "$CMD" == "copy" ]] && break
    echo "enter 'snip' or 'copy'"
  done

  STATE=$([[ "$CMD" == "snip" ]] && echo "" || echo "copy")
  OUTPUT_DIR=$([[ "$CMD" == "snip" ]] && echo "$OUTPUT_DIR" || echo "$OUTPUT_DIR/project-copy")

  IGNORE_ARGS=()
  for pattern in "${IGNORE_PATTERNS[@]}"; do
    IGNORE_ARGS+=(--ignore="$pattern")
  done

# RENAME_ARGS=()
#if [[ "$CMD" != "snip" ]]; then
#  for rule in "${RENAME_PATTERNS[@]}"; do
#    RENAME_ARGS+=(--rename="$rule")
#  done
#else
#  RENAME_ARGS=()
#fi

  # call bluehawk with all the args
  bluehawk "$CMD" \
    --state="$STATE" \
    -o "$OUTPUT_DIR" \
    "${IGNORE_ARGS[@]}" \
    "$INPUT_DIR"
else
  usage
fi
