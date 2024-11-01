#! /bin/bash

PROJECT=$(git rev-parse --show-toplevel)
JAVA_EXAMPLES=$PROJECT/java/src/main/java
GENERATED_EXAMPLES=$PROJECT/generated/java

echo "Bluehawking Java examples"
npx bluehawk snip $JAVA_EXAMPLES -o $GENERATED_EXAMPLES
