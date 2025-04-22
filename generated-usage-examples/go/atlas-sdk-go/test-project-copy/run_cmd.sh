#!/bin/bash

# Check if an action is provided
if [ -z "$1" ]; then
  echo "Usage: $0 <action>"
  exit 1
fi

ACTION=$1
MAIN_FILE="cmd/$ACTION/main.go"

# Check if the main.go file exists for the given action
if [ ! -f "$MAIN_FILE" ]; then
  echo "Error: $MAIN_FILE does not exist."
  exit 1
fi

# Run the main.go file
go run "$MAIN_FILE"
