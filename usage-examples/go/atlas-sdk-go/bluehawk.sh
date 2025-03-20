#! /bin/bash

PROJECT=$(git rev-parse --show-toplevel)
# This project requires Bluehawking multiple directories so the generated code can be copied to the user-facing `atlas-architecture-center-go-sdk` repo
GO_SDK_EXAMPLES=$PROJECT/usage-examples/go/atlas-sdk-go/cmd
GENERATED_EXAMPLES=$PROJECT/generated-usage-examples/go/sdk

# Bluehawk Admin Go SDK examples
echo "Running Bluehawk snip on $GO_SDK_EXAMPLES"
echo "Extracting snippets to the $GENERATED_EXAMPLES directory"
npx bluehawk snip "$GO_SDK_EXAMPLES" -o "$GENERATED_EXAMPLES"

# Bluehawk with state tags
echo "Running Bluehawk snip on state tags"
npx bluehawk snip "$GO_SDK_EXAMPLES" -o "$GENERATED_EXAMPLES" --state dev
npx bluehawk snip "$GO_SDK_EXAMPLES" -o "$GENERATED_EXAMPLES" --state prod

# Clean up any .gz log files
find "$PROJECT" -name "*.gz" -type f -delete
