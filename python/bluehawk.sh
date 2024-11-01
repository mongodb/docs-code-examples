#! /bin/bash

PROJECT=$(git rev-parse --show-toplevel)
PYTHON_EXAMPLES=$PROJECT/python/examples
GENERATED_EXAMPLES=$PROJECT/generated/python

# Bluehawk asymmetric examples
echo "Bluehawking Python examples"
npx bluehawk snip $PYTHON_EXAMPLES --ignore build -o $GENERATED_EXAMPLES
