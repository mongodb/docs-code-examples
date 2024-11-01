#! /bin/bash

PROJECT=$(git rev-parse --show-toplevel)
GO_EXAMPLES=$PROJECT/go/examples
GENERATED_EXAMPLES=$PROJECT/generated/go

# Bluehawk Go examples
echo "Bluehawking Go examples"
npx bluehawk snip $GO_EXAMPLES -o $GENERATED_EXAMPLES
