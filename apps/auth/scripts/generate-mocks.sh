#!/bin/bash

set -euo pipefail

MOCK_DST_ROOT="${1:-test/mock}"
MOCK_TARGET_DIR="${2:-}"

if [ -z "$MOCK_TARGET_DIR" ]; then
  echo "❌ MOCK_TARGET_DIR is required."
  exit 1
fi

MODULE_PATH=$(go list -m)

echo "🔍 Generating mocks from '$MOCK_TARGET_DIR' into '$MOCK_DST_ROOT'"

find "$MOCK_TARGET_DIR" -type f -name '*.go' -not -path '*/mocks/*' | while read -r src; do
  pkg_dir=$(dirname "$src")
  rel_dir=${pkg_dir#./}                # ./ 제거
  import_path="${MODULE_PATH}/${rel_dir}"

  # 인터페이스 추출
  interfaces=$(grep -E 'type [A-Z][a-zA-Z0-9_]+ interface' "$src" | awk '{print $2}')
  for iface in $interfaces; do
    out_dir="${MOCK_DST_ROOT}/${rel_dir}"
    mkdir -p "$out_dir"
    out_file="${out_dir}/mock_${iface}.go"

    echo "⚙️  Generating $iface → $out_file"
    if ! mockgen "$import_path" "$iface" > "$out_file"; then
      echo "❌ Failed to generate mock for $iface in $src"
    fi
  done
done

echo "✅ Mock generation complete."
