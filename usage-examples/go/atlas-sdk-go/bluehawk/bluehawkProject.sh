#!/usr/bin/env bash
set -euo pipefail

# —— Paths ——
PROJECT=$(git rev-parse --show-toplevel)
SRC_ROOT="$PROJECT/usage-examples/go/atlas-sdk-go"
DST_ROOT="$PROJECT/generated-usage-examples/go/atlas-sdk-go/full-example"

# —— Helper to run Bluehawk snip and prune nested snippets ——
run_snip() {
  local src="$1" dst="$2" label="$3"

  echo "→ $label"
  npx bluehawk snip "$src" -o "$dst" \

  # Remove any snippet files that are NOT the “-full-example” ones
  find "$dst" -type f \
    -name '*.snippet.*' \
    ! -name '*.snippet.*-full-example.*' \
    -delete -print \
    | sed 's/^/  └ removed: /'
}

# —— 1) Snip these source paths into the same relative spot under DST_ROOT ——
SNIP_PATHS=(
  "cmd/get_logs"
  "cmd/get_metrics"
  "internal"
  "internal/auth"
  ""            # root of SRC_ROOT
)

for rel in "${SNIP_PATHS[@]}"; do
  src="$SRC_ROOT/$rel"
  dst="$DST_ROOT/$rel"
  mkdir -p "$dst"
  label="Snip ${rel:-root}"
  run_snip "$src" "$dst" "$label"
done

# —— 2) Copy raw/static files into the same structure ——
## DON'T COPY SOURCE README.md TO TARGET REPO
STATIC_FILES=(
  "configs/config.json"
  "go.mod"
  "go.sum"
)

echo "→ Copying static files"
for rel in "${STATIC_FILES[@]}"; do
  src="$SRC_ROOT/$rel"
  dst_dir="$DST_ROOT/$(dirname "$rel")"
  mkdir -p "$dst_dir"
  cp "$src" "$dst_dir/"
  echo "  • $rel"
done

# —— 3) Clean up any .gz logs in the whole project ——
find "$PROJECT" -name "*.gz" -type f -delete
