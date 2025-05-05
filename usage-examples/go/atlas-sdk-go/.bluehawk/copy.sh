#!/usr/bin/env bash

# Copy project files to generated directory using Bluehawk
# Copied files are then pushed to Arch Center artifact repo via Copier App
set -euo pipefail

# Set input and output directories
PROJECT=$(git rev-parse --show-toplevel)
INPUT_DIR="$PROJECT/usage-examples/go/atlas-sdk-go/"
OUTPUT_DIR="$PROJECT/generated-usage-examples/go/atlas-sdk-go/project-copy/"

# Set directories and files to ignore
IGNORE=(
  "README.md"
  "bluehawk/"
  "tests/"
  ".*"
  "*.gz"
  "*.log"
)
IGNORE_ARGS=()
for path in "${IGNORE[@]}"; do
  IGNORE_ARGS+=("--ignore=$path")
done

# Set directories and files to rename
RENAME=(
  "REPO_README.md:README.md"
)
RENAME_ARGS=()
for path in "${RENAME[@]}"; do
  IFS=":" read -r src dst <<< "$path" # Split the path into source and destination
  RENAME_ARGS+=("--rename=$src:$dst") # Add the rename argument to the array
done

# Run Bluehawk copy command, passing all ignore args and "copy" state
.bluehawk copy \
  "${IGNORE_ARGS[@]}" \
  --state copy \
  -o "$OUTPUT_DIR" \
  "$INPUT_DIR"
