#!/bin/bash
# scripts/diff.sh - Show differences

BINARY="./esql-ast"
BASELINE_DIR="tests/baseline"
OUTPUT_DIR="tests/output"
EXAMPLES_DIR="examples"

mkdir -p "$OUTPUT_DIR"

for f in "$EXAMPLES_DIR"/*.esql; do
    name=$(basename "$f" .esql)
    echo "=== $name ==="
    
    for mode in pretty json generate analyze; do
        if [ -f "$BASELINE_DIR/$name.$mode.txt" ]; then
            "$BINARY" -f "$f" -"$mode" > "$OUTPUT_DIR/$name.$mode.txt" 2>&1
            diff -u "$BASELINE_DIR/$name.$mode.txt" "$OUTPUT_DIR/$name.$mode.txt" || echo "  No changes for $mode"
        else
            echo "  No baseline for $mode"
        fi
    done
    echo ""
done