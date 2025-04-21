#!/usr/bin/env bash

# Copy project files to generated directory using Bluehawk
# Copied files are then pushed to Arch Center artifact repo via Copier App
set -euo pipefail

# Set input and output directories
PROJECT=$(git rev-parse --show-toplevel)
INPUT_DIR="$PROJECT/usage-examples/go/atlas-sdk-go/"
OUTPUT_DIR="$PROJECT/generated-usage-examples/go/atlas-sdk-go/project-copy/"

# Set ignored internal files and directories
IGNORE=(
  "README.md"
  "bluehawk/"
  "tests/"
  ".env"
  "*.gz"
  "*.log"
)
IGNORE_ARGS=()
for path in "${IGNORE[@]}"; do
  IGNORE_ARGS+=("--ignore=$path")
done

# Run Bluehawk copy command, passing all ignore args and "copy" state
bluehawk copy \
  "${IGNORE_ARGS[@]}" \
  --state copy \
  -o "$OUTPUT_DIR" \
  "$INPUT_DIR"
