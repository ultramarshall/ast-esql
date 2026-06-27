#!/bin/bash
# scripts/test.sh - Simple test runner

BINARY="./esql-ast"
BASELINE_DIR="tests/baseline"
OUTPUT_DIR="tests/output"
EXAMPLES_DIR="examples"

mkdir -p "$BASELINE_DIR" "$OUTPUT_DIR"

echo "Running tests..."
echo ""

for f in "$EXAMPLES_DIR"/*.esql; do
    name=$(basename "$f" .esql)
    echo "Testing: $name"
    
    passed=0
    total=0
    
    for mode in pretty json generate analyze; do
        total=$((total + 1))
        "$BINARY" -f "$f" -"$mode" > "$OUTPUT_DIR/$name.$mode.txt" 2>&1
        
        if [ -f "$BASELINE_DIR/$name.$mode.txt" ]; then
            if diff -q "$BASELINE_DIR/$name.$mode.txt" "$OUTPUT_DIR/$name.$mode.txt" > /dev/null 2>&1; then
                echo "  ✅ $mode"
                passed=$((passed + 1))
            else
                echo "  ❌ $mode"
                mkdir -p tests/diff
                diff -u "$BASELINE_DIR/$name.$mode.txt" "$OUTPUT_DIR/$name.$mode.txt" > "tests/diff/$name.$mode.diff" 2>&1
            fi
        else
            echo "  ⚠️  No baseline for $mode (run: make baseline-$name)"
        fi
    done
    
    if [ $passed -eq $total ]; then
        echo "  ✅ $name PASSED ($passed/$total)"
    else
        echo "  ❌ $name FAILED ($passed/$total)"
    fi
    echo ""
done