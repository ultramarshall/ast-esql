#!/bin/bash
# scripts/baseline.sh - Generate all baselines

BINARY="./esql-ast"
BASELINE_DIR="tests/baseline"
EXAMPLES_DIR="examples"

mkdir -p "$BASELINE_DIR"

echo "Generating baselines..."

for f in "$EXAMPLES_DIR"/*.esql; do
    name=$(basename "$f" .esql)
    echo "  $name..."
    "$BINARY" -f "$f" -pretty > "$BASELINE_DIR/$name.pretty.txt" 2>&1
    "$BINARY" -f "$f" -json > "$BASELINE_DIR/$name.json.txt" 2>&1
    "$BINARY" -f "$f" -generate > "$BASELINE_DIR/$name.generate.txt" 2>&1
    "$BINARY" -f "$f" -analyze > "$BASELINE_DIR/$name.analyze.txt" 2>&1
done

echo "✅ Baselines generated!"