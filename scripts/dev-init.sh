#!/usr/bin/env bash
# dev-init.sh — Run `worky init` interactively using the local codebase.
# After init completes, automatically applies a go.mod replace directive
# so the generated workshop uses the local library instead of the published one.
#
# Usage:
#   ./scripts/dev-init.sh              # interactive, workshop created in SCAFFOLD_DIR
#   ./scripts/dev-init.sh my-workshop  # pass args directly to worky init
#   SCAFFOLD_DIR=/tmp/my-dir ./scripts/dev-init.sh

set -euo pipefail

WORKY_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
SCAFFOLD_DIR="${SCAFFOLD_DIR:-/tmp/worky-scaffold}"
WORKY_BIN="$WORKY_ROOT/bin/worky"

# Build the CLI from local source
echo "Building worky CLI..."
go build -o "$WORKY_BIN" "$WORKY_ROOT/cmd/worky"
echo ""

mkdir -p "$SCAFFOLD_DIR"
cd "$SCAFFOLD_DIR"

# Snapshot existing dirs before running init
BEFORE=$(ls -d */ 2>/dev/null | sort || true)

# Run init interactively, capturing stdout with tee so we can parse it later.
# stdin is unaffected by the pipe, so prompts/input still work normally.
TMPOUT=$(mktemp)
"$WORKY_BIN" init "$@" 2>&1 | tee "$TMPOUT"

# Parse the created directory from worky's output: 'Workshop "..." created in slug/'
SLUG=$(grep -oE 'created in [^ ]+/' "$TMPOUT" | awk '{print $3}' | tr -d '/' || true)
rm -f "$TMPOUT"

# Fall back to diffing directory listing if parsing failed
if [ -z "$SLUG" ]; then
    AFTER=$(ls -d */ 2>/dev/null | sort || true)
    SLUG=$(comm -13 <(echo "$BEFORE") <(echo "$AFTER") | head -1 | tr -d '/')
fi

if [ -z "$SLUG" ] || [ ! -f "$SCAFFOLD_DIR/$SLUG/go.mod" ]; then
    echo ""
    echo "Could not find the generated workshop directory. Skipping replace directive."
    exit 0
fi

WORKSHOP_DIR="$SCAFFOLD_DIR/$SLUG"

echo ""
echo "Applying local replace directive (local worky → $WORKY_ROOT)..."
cd "$WORKSHOP_DIR"
go mod edit -replace "github.com/davideimola/worky=$WORKY_ROOT"
go mod tidy

echo ""
echo "Workshop ready at: $WORKSHOP_DIR"
echo ""
echo "  cd $WORKSHOP_DIR && go run . serve --open"
echo ""
