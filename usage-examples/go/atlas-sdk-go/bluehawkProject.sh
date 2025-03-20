#! /bin/bash

# NOTE: This project serves as the source repo for a customer-facing `atlas-architecture-center-go-sdk` repo, so we
# need to run Bluehawk on multiple directories that need to be copied to the target repo.

PROJECT=$(git rev-parse --show-toplevel)
GO_SDK_EXAMPLES=$PROJECT/usage-examples/go/atlas-sdk-go/
GENERATED_EXAMPLES=$PROJECT/generated-usage-examples/go/atlas-sdk-go/

# Define an array of source directories/files
SOURCE_PATHS=(
  "$GO_SDK_EXAMPLES/cmd/get_logs/"
  "$GO_SDK_EXAMPLES/cmd/get_metrics/"
  "$GO_SDK_EXAMPLES/internal"
  "$GO_SDK_EXAMPLES/internal/auth"
  "$GO_SDK_EXAMPLES"
)

# Define an array of corresponding target directories
TARGET_PATHS=(
  "$GENERATED_EXAMPLES/cmd/get_logs/"
  "$GENERATED_EXAMPLES/cmd/get_metrics/"
  "$GENERATED_EXAMPLES/internal"
  "$GENERATED_EXAMPLES/internal/auth"
  "$GENERATED_EXAMPLES"
)

# Run Bluehawk snip on each source path and output to the corresponding target path
for i in "${!SOURCE_PATHS[@]}"; do
    # Create the target directory if it doesn't exist
    mkdir -p "${TARGET_PATHS[$i]}"

  echo "Running Bluehawk snip on ${SOURCE_PATHS[$i]}..."
  echo "Extracting snippets to the ${TARGET_PATHS[$i]} directory"
  npx bluehawk snip "${SOURCE_PATHS[$i]}" -o "${TARGET_PATHS[$i]}"
done

# Copy project files to generated examples directory
mkdir -p "$GENERATED_EXAMPLES/configs"
cp "$GO_SDK_EXAMPLES/configs/config.json" "$GENERATED_EXAMPLES/configs"
cp "$GO_SDK_EXAMPLES/go.*" "$GENERATED_EXAMPLES"
cp "$GO_SDK_EXAMPLES/README.md" "$GENERATED_EXAMPLES"

# Clean up any .gz log files
find "$PROJECT" -name "*.gz" -type f -delete
