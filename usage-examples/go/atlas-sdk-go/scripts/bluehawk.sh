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
INPUT_DIR="$PROJECT/usage-examples/go/atlas-sdk-go" # source directory
STATE=""
IGNORE_PATTERNS=(
  "internal_*.*" # for INTERNAL_README.md
  "scripts/"
  ".idea"
  "*_test.go" # we're not including test files in artifact repo
  ".env"
  "*.gz"
  "*.log"
  "./logs" # for generated logs directory
  # NOTE: DO NOT add pattern for ".gitignore"; we are including it in the artifact repo
)
RENAME_PATTERNS=()
VERBOSE=false

# Process command-line args
if [[ $# -gt 0 ]]; then
  CMD="$1"
  shift

  if [[ "$CMD" == "snip" ]]; then
    OUTPUT_DIR="$PROJECT/generated-usage-examples/go/atlas-sdk-go/"
    # No default STATE for snip
  elif [[ "$CMD" == "copy" ]]; then
    OUTPUT_DIR="$PROJECT/generated-usage-examples/go/atlas-sdk-go/project-copy"
    # Default STATE for copy
    STATE="copy"
  else
    usage
  fi

  # Process additional flags
  while [[ $# -gt 0 ]]; do
    case "$1" in
      --ignore=*)
        IGNORE_PATTERNS+=("${1#*=}")
        shift
        ;;
      --rename=*)
        if [[ "$CMD" == "copy" ]]; then
          RENAME_PATTERNS+=("${1#*=}")
        else
          echo "Warning: --rename is not supported for snip command, ignoring"
        fi
        shift
        ;;
      --state=*)
        STATE="${1#*=}"
        shift
        ;;
      --verbose)
        VERBOSE=true
        shift
        ;;
      *)
        echo "Unknown option: $1"
        usage
        ;;
    esac
  done
else
  # Interactive mode
  echo "=== Run Bluehawk ==="

  # 1) pick snip or copy
  while true; do
    read -rp "Enter command (snip/copy): " CMD
    [[ "$CMD" == "snip" || "$CMD" == "copy" ]] && break
    echo "enter 'snip' or 'copy'"
  done

  # Ask for verbose logging
  read -rp "Enable verbose output? (y/N): " verbose_response
  # Convert to lowercase using tr instead of bash parameter expansion
  verbose_response=$(echo "$verbose_response" | tr '[:upper:]' '[:lower:]')
  if [[ "$verbose_response" == "y" ]]; then
    VERBOSE=true
  fi
fi

# Set up the command and its arguments
if [[ "$CMD" == "snip" ]]; then
  OUTPUT_DIR="$PROJECT/generated-usage-examples/go/atlas-sdk-go/"
  # No default STATE for snip if not already set
elif [[ "$CMD" == "copy" ]]; then
  OUTPUT_DIR="$PROJECT/generated-usage-examples/go/atlas-sdk-go/project-copy"
  # Default STATE for copy if not already set
  STATE=${STATE:-"copy"}
else
  usage
fi

# Prepare ignore arguments
IGNORE_ARGS=()
for pattern in "${IGNORE_PATTERNS[@]}"; do
  IGNORE_ARGS+=(--ignore="$pattern")
done

# Prepare rename arguments (only for copy)
RENAME_ARGS=()
if [[ "$CMD" == "copy" ]] && [[ ${#RENAME_PATTERNS[@]} -gt 0 ]]; then
  for rule in "${RENAME_PATTERNS[@]}"; do
    RENAME_ARGS+=(--rename="$rule")
  done
fi

# Check for errors first
echo "Checking for Bluehawk parsing errors..."
if ! check_output=$(bluehawk check "${IGNORE_ARGS[@]}" "$INPUT_DIR" 2>&1); then
  echo "Bluehawk check failed. Errors found:"
  echo "$check_output" | grep -A 1 "bluehawk errors"
  exit 1
fi

echo "Validation successful! No errors found."

# Build the command - DO NOT use --quiet flag as it prevents file generation
CMD_ARGS=(bluehawk "$CMD" -o "$OUTPUT_DIR" "${IGNORE_ARGS[@]}")

# Add state argument if set
if [[ -n "$STATE" ]]; then
  CMD_ARGS+=(--state="$STATE")
fi

# Add rename arguments for copy command
if [[ "$CMD" == "copy" ]] && [[ ${#RENAME_ARGS[@]} -gt 0 ]]; then
  CMD_ARGS+=("${RENAME_ARGS[@]}")
fi

# Add input directory
CMD_ARGS+=("$INPUT_DIR")

# Execute the command
echo "Running: ${CMD_ARGS[0]} ${CMD_ARGS[1]} [options] ${CMD_ARGS[${#CMD_ARGS[@]}-1]}"
if [[ "$VERBOSE" == true ]]; then
  "${CMD_ARGS[@]}"
else
  output=$("${CMD_ARGS[@]}" 2>&1)
  cmd_status=$?

  # Display summary
  summary=$(echo "$output" | grep -A 3 "Processed [0-9]* files:")
  if [[ -n "$summary" ]]; then
    echo -e "\n$summary"
  else
    echo -e "\nNo summary available"
  fi

  # Filter and display only important information
  written_count=$(echo "$output" | grep -c "wrote text file" || true)
  echo -e "\nSuccessfully wrote the following $written_count files:"
  files_written=$(echo "$output" | grep "wrote text file" | sed 's/^wrote text file based on.*-> /  /')
  if [[ -n "$files_written" ]]; then
  echo "$files_written"
  else
  echo "  No files written"
  fi

  exit ${cmd_status:-0}
fi
