#! /bin/bash

PROJECT=$(git rev-parse --show-toplevel)
JS_EXAMPLES=$PROJECT/usage-examples/javascript/driver/examples
GENERATED_EXAMPLES=$PROJECT/generated-usage-examples/javascript/driver

echo "Bluehawking JavaScript examples"
npx bluehawk snip $JS_EXAMPLES -o $GENERATED_EXAMPLES
