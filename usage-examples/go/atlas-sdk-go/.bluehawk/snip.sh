#!/usr/bin/env bash
set -euo pipefail

PROJECT=$(git rev-parse --show-toplevel)
GO_SDK_EXAMPLES="$PROJECT/usage-examples/go/atlas-sdk-go/"
GENERATED_EXAMPLES="$PROJECT/generated-usage-examples/go/atlas-sdk-go/usage-examples/"

# ——— helper: run .bluehawk and only show the key lines ———
run_snip() {
  local extra_flags="$1"
  local label="$2"
  local dest="$GENERATED_EXAMPLES"

  echo "→ $label"
  npx bluehawk snip "$GO_SDK_EXAMPLES" -o "$GENERATED_EXAMPLES"
}
# ——— 1) global snippets ———
run_snip "" "Global snippets"

# ——— 2) state‑specific snippets ———
read -r -p "Do you have any state tags to enter? [y/N]: " resp
if [[ $resp =~ ^[Yy]$ ]]; then
  read -r -a STATES -p "Enter one or more state tags (space‑separated): "
  for s in "${STATES[@]}"; do
    run_snip "--state $s" "State: $s"
  done
else
  echo "No state tags — skipping."
fi

# ——— 3) copy non‑snippable files ———
# list paths *relative* to $GO_SDK_EXAMPLES
readonly COPY_FILES=(
  "configs/config.json"
  # add any other JSON (or other) files here…
  # DON'T COPY SOURCE README.md
)

echo "→ Copying raw files (${#COPY_FILES[@]})"
for rel in "${COPY_FILES[@]}"; do
  src="$GO_SDK_EXAMPLES/$rel"
  dest_dir="$GENERATED_EXAMPLES/"
  mkdir -p "$dest_dir"
  base=$(basename "$rel")
  cp "$src" "$dest_dir/snippet.$base"
  echo "  • $rel → snippet.$base"
done

# ——— 4) cleanup ———
  find "$GENERATED_EXAMPLES/" -type f \
    -name '*.snippet.*-full-example.*' \
    -delete -print \
    | sed 's/^/  └ removed: /'

find "$PROJECT" -name "*.gz" -delete
