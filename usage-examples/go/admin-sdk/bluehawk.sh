#! /bin/bash

PROJECT=$(git rev-parse --show-toplevel)
GO_SDK_EXAMPLES=$PROJECT/usage-examples/go/admin-sdk/config:$PROJECT/usage-examples/go/admin-sdk/scripts:$PROJECT/usage-examples/go/admin-sdk/internal:$PROJECT/usage-examples/go/admin-sdk/types
GENERATED_EXAMPLES=$PROJECT/generated-examples/go/admin-sdk

# Bluehawk Admin Go SDK examples
echo "Bluehawking Admin Go SDK examples"
npx bluehawk snip $GO_SDK_EXAMPLES -o $GENERATED_EXAMPLES
