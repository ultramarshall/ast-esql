# Dokumentasi Struktur & Kode Project

Dokumen ini dihasilkan secara otomatis untuk memetakan seluruh struktur folder dan isi kode di dalam project ini.

## Struktur Project (Tree)

```text
.
├── cmd
│   └── esql-ast
│       └── main.go
├── DOC.md
├── esql-ast
├── examples
│   ├── sample.esql
│   ├── test_between.esql
│   ├── test_call_graph.esql
│   ├── test_call_graph.esql.bak
│   ├── test_case.esql
│   ├── test_case_nested_if.esql
│   ├── test_case_searched_only.esql
│   ├── test_case_simple.esql
│   ├── test_case_simple_only.esql
│   ├── test_cast.esql
│   ├── test_coalesce_nullif.esql
│   ├── test_in.esql
│   ├── test_is_null.esql
│   ├── test_like.esql
│   ├── test_nested_cast.esql
│   ├── test_param.esql
│   └── test_propagate.esql
├── generate_doc.sh
├── go.mod
├── internal
│   └── token
│       └── token.go
├── Makefile
├── pkg
│   ├── analyzer
│   │   └── analyzer.go
│   ├── generator
│   │   └── generator.go
│   ├── parser
│   │   ├── ast.go
│   │   ├── lexer.go
│   │   ├── parser_expr.go
│   │   ├── parser.go
│   │   ├── parser_primary.go
│   │   ├── parser_stmt_call.go
│   │   ├── parser_stmt_control.go
│   │   ├── parser_stmt_create.go
│   │   ├── parser_stmt_declare.go
│   │   ├── parser_stmt.go
│   │   ├── parser_stmt_if.go
│   │   ├── parser_stmt_loop.go
│   │   ├── parser_stmt_set.go
│   │   └── parser_utils.go
│   ├── printer
│   │   └── printer.go
│   └── refactor
│       └── refactor.go
├── scripts
│   ├── baseline.sh
│   ├── diff.sh
│   └── test.sh
└── tests

14 directories, 45 files
```

## Isi Kode Berdasarkan File

### File: `scripts/diff.sh`

```text
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
```

---

### File: `scripts/test.sh`

```text
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
```

---

### File: `scripts/baseline.sh`

```text
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
```

---

### File: `go.mod`

```text
module esql-ast-tool

go 1.22

```

---

### File: `examples/test_between.esql`

```text
CREATE COMPUTE MODULE TestBetween
    DECLARE score INTEGER;
    DECLARE result INTEGER;
    DECLARE grade STRING;
    
    SET score = 85;
    
    -- Basic BETWEEN
    IF score BETWEEN 80 AND 100 THEN
        SET result = 1;
    END IF;
    
    -- NOT BETWEEN
    IF score NOT BETWEEN 0 AND 50 THEN
        SET result = 2;
    END IF;
    
    -- BETWEEN in CASE
    SET grade = CASE 
        WHEN score BETWEEN 90 AND 100 THEN 'A'
        WHEN score BETWEEN 80 AND 89 THEN 'B'
        WHEN score BETWEEN 70 AND 79 THEN 'C'
        ELSE 'D'
    END;
    
    -- Nested BETWEEN with expressions
    IF (score + 5) BETWEEN 80 AND 100 THEN
        SET result = 3;
    END IF;
END MODULE;
```

---

### File: `examples/test_like.esql`

```text
CREATE COMPUTE MODULE TestLike
    DECLARE name STRING;
    DECLARE result INTEGER;
    DECLARE status STRING;
    
    SET name = 'John Doe';
    
    -- Basic LIKE
    IF name LIKE 'John%' THEN
        SET result = 1;
    END IF;
    
    -- NOT LIKE
    IF name NOT LIKE '%Smith%' THEN
        SET result = 2;
    END IF;
    
    -- LIKE in CASE
    SET status = CASE 
        WHEN name LIKE 'A%' THEN 'Starts with A'
        WHEN name LIKE '%Doe' THEN 'Ends with Doe'
        ELSE 'Other'
    END;
    
    -- LIKE with underscore
    IF name LIKE 'J__n%' THEN
        SET result = 3;
    END IF;
END MODULE;
```

---

### File: `examples/test_param.esql`

```text
CREATE COMPUTE MODULE TestParameterModes

    -- Procedure dengan berbagai mode parameter
    CREATE PROCEDURE ProcessData(
        IN p_input INTEGER,
        OUT p_output CHARACTER,
        INOUT p_accumulator FLOAT
    )
    BEGIN
        SET p_output = 'Processed';
        SET p_accumulator = p_accumulator + p_input;
    END;

    -- Procedure dengan default IN (tanpa mode)
    CREATE PROCEDURE SimpleProc(
        p_name CHARACTER,
        p_value INTEGER
    )
    BEGIN
        SET p_value = p_value * 2;
    END;

    -- Function (semua parameter default IN)
    CREATE FUNCTION Calculate(
        p_a INTEGER,
        p_b INTEGER
    ) RETURNS INTEGER
    BEGIN
        RETURN p_a + p_b;
    END;

    DECLARE inputVal INTEGER;
    DECLARE outputVal CHARACTER;
    DECLARE acc FLOAT;

    SET inputVal = 10;
    SET acc = 0.0;

    CALL ProcessData(inputVal, outputVal, acc);
    SET inputVal = Calculate(5, 7);

END MODULE;
```

---

### File: `examples/test_coalesce_nullif.esql`

```text
CREATE COMPUTE MODULE TestCoalesceNullIf
    DECLARE var1 INTEGER;
    DECLARE var2 INTEGER;
    DECLARE var3 INTEGER;
    DECLARE result INTEGER;
    DECLARE resultStr STRING;
    
    -- Test data
    SET var1 = NULL;
    SET var2 = 10;
    SET var3 = 20;
    
    -- COALESCE - returns first non-NULL
    SET result = COALESCE(var1, var2, var3, 0);
    
    -- COALESCE with all NULL
    SET result = COALESCE(NULL, NULL, NULL, 0);
    
    -- COALESCE with strings
    DECLARE str1 STRING;
    DECLARE str2 STRING;
    SET str1 = NULL;
    SET str2 = 'World';
    SET resultStr = COALESCE(str1, str2, 'Default');
    
    -- COALESCE in IF condition
    IF COALESCE(var1, var2, 0) > 5 THEN
        SET result = 100;
    END IF;
    
    -- COALESCE in CASE
    SET result = CASE 
        WHEN COALESCE(var1, var2, 0) = 10 THEN 1
        ELSE 0
    END;
    
    -- NULLIF with equal values
    SET result = NULLIF(10, 10);
    
    -- NULLIF with different values
    SET result = NULLIF(10, 20);
    
    -- NULLIF with NULL as first arg
    SET result = NULLIF(NULL, 10);
    
    -- NULLIF with NULL as second arg
    SET result = NULLIF(10, NULL);
    
    -- NULLIF with both NULL
    SET result = NULLIF(NULL, NULL);
    
    -- NULLIF in IF condition
    IF NULLIF(var2, 10) IS NULL THEN
        SET result = 200;
    END IF;
    
    -- NULLIF in CASE
    SET result = CASE 
        WHEN NULLIF(var2, 10) IS NULL THEN 1
        ELSE 0
    END;
    
    -- COALESCE with NULLIF
    SET result = COALESCE(NULLIF(var2, 10), NULLIF(var3, 20), 999);
    
    -- Nested NULLIF
    SET result = NULLIF(NULLIF(10, 10), 10);
    
END MODULE;
```

---

### File: `examples/test_in.esql`

```text
CREATE COMPUTE MODULE TestIn
    DECLARE score INTEGER;
    DECLARE result INTEGER;
    DECLARE status STRING;
    
    SET score = 85;
    
    -- Basic IN
    IF score IN (80, 90, 100) THEN
        SET result = 1;
    END IF;
    
    -- NOT IN
    IF score NOT IN (1, 2, 3) THEN
        SET result = 2;
    END IF;
    
    -- IN with strings
    DECLARE name STRING;
    SET name = 'John';
    IF name IN ('John', 'Jane', 'Bob') THEN
        SET result = 3;
    END IF;
    
    -- IN in CASE
    SET status = CASE 
        WHEN score IN (90, 100) THEN 'Excellent'
        WHEN score IN (70, 80) THEN 'Good'
        ELSE 'Average'
    END;
    
    -- Nested IN with expressions
    IF (score + 5) IN (85, 90, 95) THEN
        SET result = 4;
    END IF;
END MODULE;
```

---

### File: `examples/test_case.esql`

```text
CREATE COMPUTE MODULE TestCase
    DECLARE score INTEGER;
    DECLARE grade STRING;
    DECLARE status STRING;
    DECLARE result INTEGER;
    
    SET score = 85;
    
    -- Simple CASE
    SET grade = CASE score
        WHEN 90 THEN 'A'
        WHEN 80 THEN 'B'
        WHEN 70 THEN 'C'
        ELSE 'D'
    END;
    
    -- Searched CASE
    SET status = CASE 
        WHEN score >= 90 THEN 'Excellent'
        WHEN score >= 80 THEN 'Good'
        WHEN score >= 70 THEN 'Average'
        ELSE 'Poor'
    END;
    
    -- CASE in IF condition
    IF CASE WHEN score > 80 THEN 1 ELSE 0 END = 1 THEN
        SET result = 100;
    END IF;
    
    -- Nested CASE
    SET result = CASE 
        WHEN CASE WHEN score > 80 THEN 1 ELSE 0 END = 1 THEN 100
        ELSE 0
    END;
END MODULE;

```

---

### File: `examples/sample.esql`

```text
CREATE COMPUTE MODULE TestModule
    DECLARE myVar INTEGER;
    SET myVar = 42;
    IF myVar > 0 THEN
        SET Environment.Variables.Status = 'OK';
    END IF;
END MODULE;

```

---

### File: `examples/test_case_simple.esql`

```text
CREATE COMPUTE MODULE TestCaseSimple
    DECLARE score INTEGER;
    DECLARE grade STRING;
    SET score = 85;
    SET grade = CASE score
        WHEN 90 THEN 'A'
        WHEN 80 THEN 'B'
        ELSE 'C'
    END;
END MODULE;

```

---

### File: `examples/test_case_searched_only.esql`

```text
CREATE COMPUTE MODULE TestCaseSearched
    DECLARE score INTEGER;
    DECLARE status STRING;
    SET score = 85;
    SET status = CASE 
        WHEN score >= 90 THEN 'Excellent'
        WHEN score >= 80 THEN 'Good'
        WHEN score >= 70 THEN 'Average'
        ELSE 'Poor'
    END;
END MODULE;

```

---

### File: `examples/test_cast.esql`

```text
CREATE COMPUTE MODULE TestCast
    DECLARE intVar INTEGER;
    DECLARE strVar STRING;
    DECLARE result INTEGER;
    
    SET strVar = '123';
    SET intVar = CAST(strVar AS INTEGER);
    SET result = CAST('456' AS INTEGER);
    SET strVar = CAST(789 AS STRING);
    
    IF CAST(strVar AS INTEGER) > 100 THEN
        SET result = 1;
    END IF;
END MODULE;

```

---

### File: `examples/test_nested_cast.esql`

```text
CREATE COMPUTE MODULE TestNestedCast
    DECLARE strVar STRING;
    DECLARE intVar INTEGER;
    DECLARE result INTEGER;
    
    SET strVar = '123';
    SET intVar = CAST(CAST(strVar AS STRING) AS INTEGER);
    SET result = CAST(CAST('456' AS INTEGER) AS INTEGER);
    SET strVar = CAST(CAST(789 AS STRING) AS STRING);
    
    IF CAST(CAST(strVar AS STRING) AS INTEGER) > 100 THEN
        SET result = 1;
    END IF;
END MODULE;

```

---

### File: `examples/test_case_simple_only.esql`

```text
CREATE COMPUTE MODULE TestCaseSimple
    DECLARE score INTEGER;
    DECLARE grade STRING;
    SET score = 85;
    SET grade = CASE score
        WHEN 90 THEN 'A'
        WHEN 80 THEN 'B'
        ELSE 'C'
    END;
END MODULE;

```

---

### File: `examples/test_call_graph.esql`

```text
CREATE COMPUTE MODULE TestCallGraph
    DECLARE studentScore INTEGER;
    DECLARE result INTEGER;
    
    CREATE PROCEDURE ProcA()
    BEGIN
        CALL ProcB();
        SET studentScore = 10;
    END;
    
    CREATE PROCEDURE ProcB()
    BEGIN
        CALL ProcC();
        SET studentScore = 20;
    END;
    
    CREATE PROCEDURE ProcC()
    BEGIN
        SET studentScore = 30;
    END;
    
    CREATE FUNCTION FuncA() RETURNS INTEGER
    BEGIN
        RETURN studentScore + 10;
    END;
    
    -- Main
    CALL ProcA();
    SET result = FuncA();
END MODULE;

```

---

### File: `examples/test_propagate.esql`

```text
CREATE COMPUTE MODULE TestPropagate
    DECLARE status STRING;
    DECLARE result INTEGER;
    
    -- Basic PROPAGATE
    PROPAGATE status;
    
    -- PROPAGATE with expression
    PROPAGATE result + 10;
    
    -- PROPAGATE in IF
    IF status IS NOT NULL THEN
        PROPAGATE status;
    END IF;
    
    -- PROPAGATE with multiple expressions
    PROPAGATE status, result;
END MODULE;

```

---

### File: `examples/test_case_nested_if.esql`

```text
CREATE COMPUTE MODULE TestCaseNestedIf
    DECLARE score INTEGER;
    DECLARE result INTEGER;
    SET score = 85;
    IF CASE WHEN score > 80 THEN 1 ELSE 0 END = 1 THEN
        SET result = 100;
    END IF;
END MODULE;

```

---

### File: `examples/test_is_null.esql`

```text
CREATE COMPUTE MODULE TestIsNull
    DECLARE var1 INTEGER;
    DECLARE var2 STRING;
    DECLARE result INTEGER;
    
    SET var1 = NULL;
    SET var2 = 'Hello';
    
    -- Test IS NULL
    IF var1 IS NULL THEN
        SET result = 1;
    END IF;
    
    -- Test IS NOT NULL
    IF var2 IS NOT NULL THEN
        SET result = 2;
    END IF;
    
    -- Test dalam CASE
    SET result = CASE 
        WHEN var1 IS NULL THEN 100
        WHEN var2 IS NOT NULL THEN 200
        ELSE 0
    END;
END MODULE;
```

---

### File: `generate_doc.sh`

```text
#!/bin/bash

OUTPUT_FILE="DOC.md"

echo "📝 Menyusun dokumentasi ke $OUTPUT_FILE..."

# Tulis Header Utama
cat << 'EOF' > "$OUTPUT_FILE"
# Dokumentasi Struktur & Kode Project

Dokumen ini dihasilkan secara otomatis untuk memetakan seluruh struktur folder dan isi kode di dalam project ini.

## Struktur Project (Tree)

```text
EOF

# Jalankan perintah tree jika ada, jika tidak pakai find untuk simulasi tree sederhana
if command -v tree &> /dev/null; then
    tree -I "node_modules|.git|vendor|.next|dist" >> "$OUTPUT_FILE"
else
    find . -not -path '*/.*' -not -path './vendor*' -not -path './node_modules*' | sed -e 's/^[^\/]*\//⎹  /' -e 's/\/[^\/]*$/⎹__/' >> "$OUTPUT_FILE"
fi

echo '```' >> "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"
echo "## Isi Kode Berdasarkan File" >> "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"

# Mencari file kode yang valid
find . -type f \( -name "*.go" -o -name "go.mod" -o -name "*.js" -o -name "*.ts" -o -name "*.py" -o -name "*.php" -o -name "*.json" -o -name "*.html" -o -name "*.css" -o -name "*.esql" -o -name "*.sh" \) \
-not -path "*/.*" \
-not -path "./vendor/*" \
-not -path "./node_modules/*" \
-not -path "./dist/*" \
-not -path "./.next/*" | while read -r file; do
    
    # Hapus `./` di depan nama file agar rapi
    clean_path=$(echo "$file" | sed 's|^\./||')
    
    echo "### File: \`$clean_path\`" >> "$OUTPUT_FILE"
    echo "" >> "$OUTPUT_FILE"
    
    # Ambil ekstensi file
    ext="${file##*.}"
    
    # Tentukan syntax highlighting menggunakan kutip tunggal agar aman
    if [[ "$file" == *"go.mod" ]]; then
        echo '```text' >> "$OUTPUT_FILE"
    elif [[ "$ext" =~ ^(go|js|ts|py|php|json|html|css)$ ]]; then
        echo '```'"$ext" >> "$OUTPUT_FILE"
    else
        echo '```text' >> "$OUTPUT_FILE"
    fi
    
    # Masukkan isi file
    cat "$file" >> "$OUTPUT_FILE"
    
    echo "" >> "$OUTPUT_FILE"
    echo '```' >> "$OUTPUT_FILE"
    echo "" >> "$OUTPUT_FILE"
    echo "---" >> "$OUTPUT_FILE"
    echo "" >> "$OUTPUT_FILE"
done

echo "✅ Berhasil! File $OUTPUT_FILE telah dibuat."
```

---

### File: `pkg/parser/parser_stmt_set.go`

```go
package parser

import (
	"fmt"

	"esql-ast-tool/internal/token"
)

func (p *Parser) parseSet() ASTNode {
	debugPrint("  [parseSet] START: token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)

	node := NewASTNode(SetNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
	p.nextToken()

	// Parse target
	var target ASTNode
	switch p.curToken.Type {
	case token.IDENTIFIER:
		target = p.parseIdentifier()
	case token.ENVIRONMENT:
		target = p.parseFieldReferenceFromKeyword("Environment")
	case token.FIELD:
		target = p.parseField()
	default:
		target = p.parseExpression()
	}

	if target.Type != "" {
		targetWrapper := NewASTNode(BlockNode, "target", target.Span.Start.Line, target.Span.Start.Column)
		targetWrapper.AddChild(target)
		targetWrapper.Span = target.Span
		node.AddChild(targetWrapper)
	}

	// Parse '='
	if p.curToken.Type == token.ASSIGN || p.curToken.Type == token.EQ {
		p.nextToken()

		// ✅ Cek apakah ada expression setelah '='
		if p.curToken.Type == token.SEMICOLON || p.curToken.Type == token.EOF {
			p.errors = append(p.errors,
				fmt.Sprintf("expected expression after '=' in SET statement at line %d",
					p.curToken.Line))
			// Consume semicolon or EOF to avoid infinite loop
			if p.curToken.Type == token.SEMICOLON {
				p.nextToken()
			}
			return node
		}

		// ✅ Cek juga jika token adalah ';' setelah whitespace
		// (ini sudah di-handle di atas)

		value := p.parseAssignmentRHS()
		if value.Type != "" {
			valueWrapper := NewASTNode(BlockNode, "value", value.Span.Start.Line, value.Span.Start.Column)
			valueWrapper.AddChild(value)
			valueWrapper.Span = value.Span
			node.AddChild(valueWrapper)
			node.Span.End = value.Span.End
		} else {
			// ✅ Kalau parseAssignmentRHS return kosong, tambahkan error
			p.errors = append(p.errors,
				fmt.Sprintf("invalid expression after '=' in SET statement at line %d",
					p.curToken.Line))
			return node
		}
	} else {
		p.errors = append(p.errors,
			fmt.Sprintf("expected '=' in SET statement, got %s '%s' at line %d",
				p.curToken.Type, p.curToken.Literal, p.curToken.Line))
		p.nextToken()
	}

	// Consume semicolon
	if p.curToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	debugPrint("  [parseSet] END: token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)

	return node
}

func (p *Parser) parseAssignmentRHS() ASTNode {
	// Parse sebagai additive dulu
	left := p.parseAdditive()

	debugPrint("  [parseAssignmentRHS] after parseAdditive: token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)

	// Jika setelah parseAdditive masih ada operator comparison selain '='
	// Misalnya: >, <, >=, <=, !=
	if p.curToken.Type == token.GT || p.curToken.Type == token.LT ||
		p.curToken.Type == token.GTE || p.curToken.Type == token.LTE ||
		p.curToken.Type == token.NOT_EQ {
		// Ini adalah comparison, parse sebagai comparison
		debugPrint("  [parseAssignmentRHS] found comparison operator, parsing as comparison\n")
		return p.parseComparisonSuffix(left)
	}

	// ✅ Jika token saat ini adalah ';', berarti selesai
	if p.curToken.Type == token.SEMICOLON {
		debugPrint("  [parseAssignmentRHS] found ';', returning left\n")
		return left
	}

	// ✅ Jika token saat ini adalah '=', abaikan (sudah di-consume di parseSet)
	if p.curToken.Type == token.EQ || p.curToken.Type == token.ASSIGN {
		debugPrint("  [parseAssignmentRHS] found '=', but in SET context, skipping...\n")
		// Consume token '=' agar tidak menyebabkan infinite loop
		p.nextToken()
		return left
	}

	debugPrint("  [parseAssignmentRHS] returning left (no more operators)\n")
	return left
}

```

---

### File: `pkg/parser/ast.go`

```go
package parser

import (
	"encoding/json"
	"fmt"
)

type NodeType string

const (
	// Statements
	ProgramNode     NodeType = "Program"
	ModuleNode      NodeType = "Module"
	FunctionNode    NodeType = "Function"
	ProcedureNode   NodeType = "Procedure"
	DeclareNode     NodeType = "Declare"
	SetNode         NodeType = "Set"
	IfNode          NodeType = "If"
	ElseNode        NodeType = "Else"
	ElseIfNode      NodeType = "ElseIf"
	WhileNode       NodeType = "While"
	ForNode         NodeType = "For"
	CaseNode        NodeType = "Case"
	WhenNode        NodeType = "When"
	OtherwiseNode   NodeType = "Otherwise"
	ReturnNode      NodeType = "Return"
	ThrowNode       NodeType = "Throw"
	CreateNode      NodeType = "Create"
	EnvironmentNode NodeType = "Environment"
	DatabaseNode    NodeType = "Database"
	PassthruNode    NodeType = "Passthru"
	MoveNode        NodeType = "Move"
	PropagateNode   NodeType = "Propagate"
	ContinueNode    NodeType = "Continue"
	BreakNode       NodeType = "Break"
	LabelNode       NodeType = "Label"
	BlockNode       NodeType = "Block"
	CallNode        NodeType = "Call"

	// Expressions
	BinaryOpNode       NodeType = "BinaryOp"
	UnaryOpNode        NodeType = "UnaryOp"
	ComparisonNode     NodeType = "Comparison"
	FunctionCallNode   NodeType = "FunctionCall"
	FieldReferenceNode NodeType = "FieldReference"
	ArrayIndexNode     NodeType = "ArrayIndex"
	LiteralNode        NodeType = "Literal"
	IdentifierNode     NodeType = "Identifier"
	CastNode           NodeType = "Cast"
	IsNullNode         NodeType = "IsNull"
	IsNotNullNode      NodeType = "IsNotNull"
	BetweenNode        NodeType = "Between"
	ParenthesizedNode  NodeType = "Parenthesized"
	LikeNode           NodeType = "Like"
	InNode             NodeType = "In"
	CoalesceNode       NodeType = "Coalesce"
	NullIfNode         NodeType = "NullIf"
	ParameterNode      NodeType = "Parameter"
	ModeNode           NodeType = "Mode"
	ParameterNameNode  NodeType = "ParamName"
	ParameterTypeNode  NodeType = "ParamType"
	ReturnTypeNode     NodeType = "ReturnType"
)

// Position represents a position in source code
type Position struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

// Span represents a range in source code (start to end)
type Span struct {
	Start Position `json:"start"`
	End   Position `json:"end"`
}

// ASTNode represents a node in the Abstract Syntax Tree
type ASTNode struct {
	Type     NodeType    `json:"type"`
	Value    interface{} `json:"value,omitempty"`
	Children []ASTNode   `json:"children,omitempty"`
	Span     Span        `json:"span"`
	Token    string      `json:"-"`
	Not      bool        `json:"not,omitempty"`
	Target   string      `json:"target,omitempty"`
}

// NewASTNode creates a new AST node with start position
func NewASTNode(nodeType NodeType, token string, line, column int) ASTNode {
	return ASTNode{
		Type:     nodeType,
		Children: []ASTNode{},
		Span: Span{
			Start: Position{Line: line, Column: column},
			End:   Position{Line: line, Column: column + len(token)},
		},
		Token: token,
		Not:   false,
	}
}

// AddChild adds a child node and extends span
func (n *ASTNode) AddChild(child ASTNode) {
	n.Children = append(n.Children, child)
	n.ExtendSpan(child)
}

// ExtendSpan extends the node's span to include the child's span
func (n *ASTNode) ExtendSpan(child ASTNode) {
	if child.Span.End.Line > n.Span.End.Line ||
		(child.Span.End.Line == n.Span.End.Line && child.Span.End.Column > n.Span.End.Column) {
		n.Span.End = child.Span.End
	}
}

// SetEnd sets the end position of the node
func (n *ASTNode) SetEnd(line, column int) {
	n.Span.End = Position{Line: line, Column: column}
}

// ToJSON returns JSON representation of the node
func (p Program) ToJSON() ([]byte, error) {
	if p.Root.Type != "" {
		return json.MarshalIndent(p.Root, "", "  ")
	}
	return json.MarshalIndent(p.Statements, "", "  ")
}

// Program represents a complete program
type Program struct {
	Statements []ASTNode `json:"statements"`
	Root       ASTNode   `json:"root,omitempty"`
}

func NewProgram() Program {
	return Program{
		Statements: []ASTNode{},
		Root:       ASTNode{},
	}
}

func (p *Program) AddStatement(stmt ASTNode) {
	p.Statements = append(p.Statements, stmt)
}

// String returns a string representation of Position
func (p Position) String() string {
	return fmt.Sprintf("%d:%d", p.Line, p.Column)
}

// String returns a string representation of Span
func (s Span) String() string {
	return fmt.Sprintf("[%s - %s]", s.Start.String(), s.End.String())
}

```

---

### File: `pkg/parser/parser_stmt_create.go`

```go
package parser

import (
	"esql-ast-tool/internal/token"
	"fmt"
)

// ============================================
// CREATE COMPUTE MODULE
// ============================================

func (p *Parser) parseCreate() ASTNode {
	p.nextToken() // consume CREATE

	if p.curToken.Type == token.COMPUTE {
		p.nextToken()
		if p.curToken.Type == token.MODULE {
			p.nextToken()
			// Langsung return module node, tanpa dibungkus CreateNode
			return p.parseComputeModule()
		}
	}

	// Fallback
	return ASTNode{}
}

func (p *Parser) parseComputeModule() ASTNode {
	moduleNode := NewASTNode(ModuleNode, "COMPUTE MODULE", p.curToken.Line, p.curToken.Column)
	moduleNode.Value = "COMPUTE"

	// Parse module name
	if p.curToken.Type == token.IDENTIFIER {
		nameNode := NewASTNode(IdentifierNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
		nameNode.Value = p.curToken.Literal
		moduleNode.AddChild(nameNode)
		p.nextToken()
	}

	// Parse body - LANGSUNG parse statement, JANGAN lewat parseCreate lagi
	for p.curToken.Type != token.END && p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt.Type != "" {
			moduleNode.AddChild(stmt)
		}
		if p.curToken.Type == token.SEMICOLON {
			p.nextToken()
		}
	}

	// Consume END MODULE
	if p.curToken.Type == token.END {
		endLine := p.curToken.Line
		endCol := p.curToken.Column + 3
		p.nextToken()
		if p.curToken.Type == token.MODULE {
			endCol = p.curToken.Column + len(p.curToken.Literal)
			p.nextToken()
		}
		moduleNode.Span.End = Position{Line: endLine, Column: endCol}
	}

	return moduleNode
}

// ============================================
// MODULE Statement
// ============================================

func (p *Parser) parseModuleStatement() ASTNode {
	node := NewASTNode(ModuleNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
	p.nextToken()

	// Parse module name
	if p.curToken.Type == token.IDENTIFIER {
		nameNode := NewASTNode(IdentifierNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
		nameNode.Value = p.curToken.Literal
		node.AddChild(nameNode)
		p.nextToken()
	}

	// Parse body
	for p.curToken.Type != token.END && p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt.Type != "" {
			node.AddChild(stmt)
		}
		if p.curToken.Type == token.SEMICOLON {
			p.nextToken()
		}
	}

	if p.curToken.Type == token.END {
		p.nextToken()
		if p.curToken.Type == token.MODULE {
			p.nextToken()
		}
	}

	if p.curToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return node
}

// ============================================
// FUNCTION Statement
// ============================================

func (p *Parser) parseFunctionStatement() ASTNode {
	node := NewASTNode(FunctionNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
	p.nextToken()

	// Parse function name
	if p.curToken.Type == token.IDENTIFIER {
		nameNode := NewASTNode(IdentifierNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
		nameNode.Value = p.curToken.Literal
		node.AddChild(nameNode)
		p.nextToken()
	} else {
		p.errors = append(p.errors, fmt.Sprintf("expected function name, got %s", p.curToken.Type))
		return node
	}

	// Parse parameters (harus ada LPAREN)
	if p.curToken.Type == token.LPAREN {
		p.parseFunctionParameters(&node)
	} else {
		p.errors = append(p.errors, fmt.Sprintf("expected '(' after function name, got %s", p.curToken.Type))
		return node
	}

	// Parse return type
	if p.curToken.Type == token.RETURNS {
		p.parseFunctionReturnType(&node)
	} else {
		p.errors = append(p.errors, fmt.Sprintf("expected RETURNS, got %s", p.curToken.Type))
		return node
	}

	// Parse function body
	if p.curToken.Type == token.BEGIN {
		p.parseFunctionBody(&node)
	} else {
		p.errors = append(p.errors, fmt.Sprintf("expected BEGIN, got %s", p.curToken.Type))
		return node
	}

	if p.curToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return node
}

func (p *Parser) parseFunctionBody(node *ASTNode) {
	p.nextToken()
	for p.curToken.Type != token.END && p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt.Type != "" {
			node.AddChild(stmt)
		}
		if p.curToken.Type != token.END && p.curToken.Type != token.EOF {
			p.nextToken()
		}
	}
	if p.curToken.Type == token.END {
		p.nextToken()
	}
}

func (p *Parser) parseFunctionParameters(node *ASTNode) {
	p.nextToken() // consume '('

	for p.curToken.Type != token.RPAREN && p.curToken.Type != token.EOF {
		// Simpan posisi awal parameter
		paramStartLine := p.curToken.Line
		paramStartCol := p.curToken.Column

		// 1. Parse mode
		var mode string
		var modeStartLine, modeStartCol int
		if p.curToken.Type == token.IN || p.curToken.Type == token.OUT || p.curToken.Type == token.INOUT {
			mode = p.curToken.Literal
			modeStartLine = p.curToken.Line
			modeStartCol = p.curToken.Column
			p.nextToken()
		} else {
			mode = "IN"
			// Jika tidak ada mode, posisi awal adalah nama parameter
			modeStartLine = p.curToken.Line
			modeStartCol = p.curToken.Column
		}

		// 2. Parse parameter name
		if p.curToken.Type != token.IDENTIFIER {
			p.errors = append(p.errors, fmt.Sprintf("expected parameter name, got %s", p.curToken.Type))
			break
		}
		nameStartLine := p.curToken.Line
		nameStartCol := p.curToken.Column
		paramName := p.curToken.Literal
		p.nextToken()

		// 3. Parse parameter type
		if p.curToken.Type != token.IDENTIFIER {
			p.errors = append(p.errors, fmt.Sprintf("expected parameter type, got %s", p.curToken.Type))
			break
		}
		typeStartLine := p.curToken.Line
		typeStartCol := p.curToken.Column
		paramType := p.curToken.Literal
		p.nextToken()

		// 4. Create Parameter node
		paramNode := NewASTNode(ParameterNode, "parameter", paramStartLine, paramStartCol)

		// Mode node
		modeNode := NewASTNode(LiteralNode, mode, modeStartLine, modeStartCol)
		modeNode.Value = mode
		modeNode.Span = Span{
			Start: Position{Line: modeStartLine, Column: modeStartCol},
			End:   Position{Line: modeStartLine, Column: modeStartCol + len(mode)},
		}
		paramNode.AddChild(modeNode)

		// Name node
		nameNode := NewASTNode(IdentifierNode, paramName, nameStartLine, nameStartCol)
		nameNode.Value = paramName
		nameNode.Span = Span{
			Start: Position{Line: nameStartLine, Column: nameStartCol},
			End:   Position{Line: nameStartLine, Column: nameStartCol + len(paramName)},
		}
		paramNode.AddChild(nameNode)

		// Type node
		typeNode := NewASTNode(IdentifierNode, paramType, typeStartLine, typeStartCol)
		typeNode.Value = paramType
		typeNode.Span = Span{
			Start: Position{Line: typeStartLine, Column: typeStartCol},
			End:   Position{Line: typeStartLine, Column: typeStartCol + len(paramType)},
		}
		paramNode.AddChild(typeNode)

		// Span parameter: dari awal parameter sampai akhir type
		paramNode.Span = Span{
			Start: Position{Line: paramStartLine, Column: paramStartCol},
			End:   Position{Line: typeStartLine, Column: typeStartCol + len(paramType)},
		}
		node.AddChild(paramNode)

		if p.curToken.Type == token.COMMA {
			p.nextToken()
		}
	}
	if p.curToken.Type == token.RPAREN {
		p.nextToken()
	}
}

func (p *Parser) parseFunctionReturnType(node *ASTNode) {
	p.nextToken() // consume RETURNS
	if p.curToken.Type == token.IDENTIFIER {
		returnType := NewASTNode(ReturnTypeNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
		returnType.Value = p.curToken.Literal
		node.AddChild(returnType)
		p.nextToken()
	}
}

func (p *Parser) parseProcedureBody(node *ASTNode) {
	p.nextToken() // consume BEGIN
	for p.curToken.Type != token.END && p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt.Type != "" {
			node.AddChild(stmt)
		}
		// parseStatement sudah memajukan token ke setelah statement,
		// tidak perlu p.nextToken() lagi di sini.
	}
	if p.curToken.Type == token.END {
		p.nextToken()
	}
}

// ============================================
// PROCEDURE Statement
// ============================================

func (p *Parser) parseProcedureStatement() ASTNode {
	node := NewASTNode(ProcedureNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
	p.nextToken()

	// Parse procedure name
	if p.curToken.Type == token.IDENTIFIER {
		nameNode := NewASTNode(IdentifierNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
		nameNode.Value = p.curToken.Literal
		node.AddChild(nameNode)
		p.nextToken()
	} else {
		p.errors = append(p.errors, fmt.Sprintf("expected procedure name, got %s", p.curToken.Type))
		return node
	}

	// Parse parameters
	if p.curToken.Type == token.LPAREN {
		p.parseProcedureParameters(&node)
	} else {
		p.errors = append(p.errors, fmt.Sprintf("expected '(' after procedure name, got %s", p.curToken.Type))
		return node
	}

	// Parse procedure body
	if p.curToken.Type == token.BEGIN {
		p.parseProcedureBody(&node)
	} else {
		p.errors = append(p.errors, fmt.Sprintf("expected BEGIN, got %s", p.curToken.Type))
		return node
	}

	if p.curToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return node
}

func (p *Parser) parseProcedureParameters(node *ASTNode) {
	p.nextToken() // consume '('

	for p.curToken.Type != token.RPAREN && p.curToken.Type != token.EOF {
		paramStartLine := p.curToken.Line
		paramStartCol := p.curToken.Column

		var mode string
		var modeStartLine, modeStartCol int
		if p.curToken.Type == token.IN || p.curToken.Type == token.OUT || p.curToken.Type == token.INOUT {
			mode = p.curToken.Literal
			modeStartLine = p.curToken.Line
			modeStartCol = p.curToken.Column
			p.nextToken()
		} else {
			mode = "IN"
			modeStartLine = p.curToken.Line
			modeStartCol = p.curToken.Column
		}

		if p.curToken.Type != token.IDENTIFIER {
			p.errors = append(p.errors, fmt.Sprintf("expected parameter name, got %s", p.curToken.Type))
			break
		}
		nameStartLine := p.curToken.Line
		nameStartCol := p.curToken.Column
		paramName := p.curToken.Literal
		p.nextToken()

		if p.curToken.Type != token.IDENTIFIER {
			p.errors = append(p.errors, fmt.Sprintf("expected parameter type, got %s", p.curToken.Type))
			break
		}
		typeStartLine := p.curToken.Line
		typeStartCol := p.curToken.Column
		paramType := p.curToken.Literal
		p.nextToken()

		paramNode := NewASTNode(ParameterNode, "parameter", paramStartLine, paramStartCol)

		modeNode := NewASTNode(LiteralNode, mode, modeStartLine, modeStartCol)
		modeNode.Value = mode
		modeNode.Span = Span{
			Start: Position{Line: modeStartLine, Column: modeStartCol},
			End:   Position{Line: modeStartLine, Column: modeStartCol + len(mode)},
		}
		paramNode.AddChild(modeNode)

		nameNode := NewASTNode(IdentifierNode, paramName, nameStartLine, nameStartCol)
		nameNode.Value = paramName
		nameNode.Span = Span{
			Start: Position{Line: nameStartLine, Column: nameStartCol},
			End:   Position{Line: nameStartLine, Column: nameStartCol + len(paramName)},
		}
		paramNode.AddChild(nameNode)

		typeNode := NewASTNode(IdentifierNode, paramType, typeStartLine, typeStartCol)
		typeNode.Value = paramType
		typeNode.Span = Span{
			Start: Position{Line: typeStartLine, Column: typeStartCol},
			End:   Position{Line: typeStartLine, Column: typeStartCol + len(paramType)},
		}
		paramNode.AddChild(typeNode)

		paramNode.Span = Span{
			Start: Position{Line: paramStartLine, Column: paramStartCol},
			End:   Position{Line: typeStartLine, Column: typeStartCol + len(paramType)},
		}
		node.AddChild(paramNode)

		if p.curToken.Type == token.COMMA {
			p.nextToken()
		}
	}
	if p.curToken.Type == token.RPAREN {
		p.nextToken()
	}
}

```

---

### File: `pkg/parser/parser_primary.go`

```go
package parser

import (
	"esql-ast-tool/internal/token"
	"fmt"
	"strconv"
)

func (p *Parser) parsePrimary() ASTNode {
	debugPrint("    [parsePrimary] START: token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)

	switch p.curToken.Type {
	case token.IDENTIFIER:
		result := p.parseIdentifier()
		debugPrint("    [parsePrimary] after IDENTIFIER: token=%s, literal='%s'\n",
			p.curToken.Type, p.curToken.Literal)
		return result
	case token.NUMBER:
		result := p.parseNumber()
		debugPrint("    [parsePrimary] after NUMBER: token=%s, literal='%s'\n",
			p.curToken.Type, p.curToken.Literal)
		return result
	case token.STRING:
		result := p.parseString()
		debugPrint("    [parsePrimary] after STRING: token=%s, literal='%s'\n",
			p.curToken.Type, p.curToken.Literal)
		return result
	case token.LPAREN:
		result := p.parseGroupedExpression()
		debugPrint("    [parsePrimary] after LPAREN: token=%s, literal='%s'\n",
			p.curToken.Type, p.curToken.Literal)
		return result
	case token.CASE:
		result := p.parseCase()
		debugPrint("    [parsePrimary] after CASE: token=%s, literal='%s'\n",
			p.curToken.Type, p.curToken.Literal)
		return result
	case token.CAST:
		result := p.parseCast()
		debugPrint("    [parsePrimary] after CAST: token=%s, literal='%s'\n",
			p.curToken.Type, p.curToken.Literal)
		return result
	case token.COALESCE:
		result := p.parseCoalesce()
		debugPrint("    [parsePrimary] after COALESCE: token=%s, literal='%s'\n",
			p.curToken.Type, p.curToken.Literal)
		return result
	case token.NULLIF:
		result := p.parseNullIf()
		debugPrint("    [parsePrimary] after NULLIF: token=%s, literal='%s'\n",
			p.curToken.Type, p.curToken.Literal)
		return result
	case token.DOT:
		debugPrint("    [parsePrimary] WARNING: DOT without identifier, skipping...\n")
		p.nextToken()
		return ASTNode{}
	default:
		// Kalau nemu token yang gak dikenal, consume biar gak loop
		debugPrint("    [parsePrimary] UNKNOWN: token=%s, literal='%s', consuming...\n",
			p.curToken.Type, p.curToken.Literal)
		p.nextToken() // ← INI PENTING! consume token biar gak loop
		return ASTNode{}
	}
}

func (p *Parser) parseIdentifier() ASTNode {
	debugPrint("      [parseIdentifier] token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)

	node := NewASTNode(IdentifierNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
	node.Value = p.curToken.Literal
	p.nextToken()

	if p.curToken.Type == token.LPAREN {
		return p.parseFunctionCall(node)
	}

	if p.curToken.Type == token.DOT {
		return p.parseFieldReference(node)
	}

	debugPrint("      [parseIdentifier] END: token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)

	return node
}

func (p *Parser) parseFunctionCall(name ASTNode) ASTNode {
	funcName := name.Token
	if name.Value != nil {
		if str, ok := name.Value.(string); ok {
			funcName = str
		}
	}
	node := NewASTNode(FunctionCallNode, funcName, name.Span.Start.Line, name.Span.Start.Column)
	node.Value = funcName
	node.Span.Start = name.Span.Start

	// ✅ Tambahkan name sebagai child pertama
	nameNode := NewASTNode(IdentifierNode, funcName, name.Span.Start.Line, name.Span.Start.Column)
	nameNode.Value = funcName
	node.AddChild(nameNode)

	p.nextToken() // consume '('

	if p.curToken.Type != token.RPAREN {
		arg := p.parseExpression()
		if arg.Type != "" {
			node.AddChild(arg)
		}

		for p.curToken.Type == token.COMMA {
			p.nextToken()
			arg = p.parseExpression()
			if arg.Type != "" {
				node.AddChild(arg)
			}
		}
	}

	if p.curToken.Type == token.RPAREN {
		endLine := p.curToken.Line
		endCol := p.curToken.Column + 1
		p.nextToken()
		node.Span.End = Position{Line: endLine, Column: endCol}
	}

	return node
}

func (p *Parser) parseFieldReference(base ASTNode) ASTNode {
	fieldNode := NewASTNode(FieldReferenceNode, "field", base.Span.Start.Line, base.Span.Start.Column)
	fieldNode.Span.Start = base.Span.Start
	fieldNode.AddChild(base)

	if base.Value != nil {
		fieldNode.Value = base.Value
	} else {
		fieldNode.Value = base.Token
	}

	for p.curToken.Type == token.DOT {
		p.nextToken()
		if p.curToken.Type == token.IDENTIFIER {
			identNode := NewASTNode(IdentifierNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
			identNode.Value = p.curToken.Literal

			newFieldNode := NewASTNode(FieldReferenceNode, "field", fieldNode.Span.Start.Line, fieldNode.Span.Start.Column)
			newFieldNode.AddChild(fieldNode)
			newFieldNode.AddChild(identNode)

			if fieldNode.Value != nil {
				newFieldNode.Value = fieldNode.Value.(string) + "." + p.curToken.Literal
			} else {
				newFieldNode.Value = p.curToken.Literal
			}

			newFieldNode.Span.Start = fieldNode.Span.Start
			newFieldNode.Span.End = Position{Line: p.curToken.Line, Column: p.curToken.Column + len(p.curToken.Literal)}

			fieldNode = newFieldNode
			p.nextToken()
		}
	}

	return fieldNode
}

func (p *Parser) parseNumber() ASTNode {
	val, _ := strconv.ParseFloat(p.curToken.Literal, 64)
	node := NewASTNode(LiteralNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
	node.Value = val
	p.nextToken()
	return node
}

func (p *Parser) parseString() ASTNode {
	node := NewASTNode(LiteralNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
	node.Value = p.curToken.Literal
	p.nextToken()
	return node
}

func (p *Parser) parseGroupedExpression() ASTNode {
	startLine := p.curToken.Line
	startCol := p.curToken.Column
	p.nextToken()

	expr := p.parseExpression()
	if p.curToken.Type == token.RPAREN {
		endLine := p.curToken.Line
		endCol := p.curToken.Column + 1
		p.nextToken()

		// Buat Parenthesized node untuk menyimpan tanda kurung
		parenNode := NewASTNode(ParenthesizedNode, "()", startLine, startCol)
		parenNode.AddChild(expr)
		parenNode.Span = Span{
			Start: Position{Line: startLine, Column: startCol},
			End:   Position{Line: endLine, Column: endCol},
		}
		return parenNode
	}
	return expr
}

func (p *Parser) parseCast() ASTNode {
	debugPrint("    [parseCast] START: token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)

	node := NewASTNode(CastNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
	node.Value = "CAST"
	p.nextToken()

	debugPrint("    [parseCast] after CAST: token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)

	if p.curToken.Type != token.LPAREN {
		p.errors = append(p.errors,
			fmt.Sprintf("expected '(' after CAST, got %s at line %d",
				p.curToken.Type, p.curToken.Line))
		return node
	}
	p.nextToken()

	debugPrint("    [parseCast] after '(': token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)

	expr := p.parseExpression()
	if expr.Type != "" {
		node.AddChild(expr)
	}

	debugPrint("    [parseCast] after expression: token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)

	if p.curToken.Type != token.AS {
		p.errors = append(p.errors,
			fmt.Sprintf("expected 'AS' in CAST expression, got %s at line %d",
				p.curToken.Type, p.curToken.Line))
		return node
	}
	p.nextToken()

	debugPrint("    [parseCast] after AS: token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)

	if p.curToken.Type != token.IDENTIFIER {
		p.errors = append(p.errors,
			fmt.Sprintf("expected type name after AS in CAST, got %s at line %d",
				p.curToken.Type, p.curToken.Line))
		return node
	}

	typeNode := NewASTNode(IdentifierNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
	typeNode.Value = p.curToken.Literal
	node.AddChild(typeNode)
	p.nextToken()

	debugPrint("    [parseCast] after type: token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)

	if p.curToken.Type != token.RPAREN {
		p.errors = append(p.errors,
			fmt.Sprintf("expected ')' in CAST expression, got %s at line %d",
				p.curToken.Type, p.curToken.Line))
		return node
	}
	endLine := p.curToken.Line
	endCol := p.curToken.Column + 1
	p.nextToken()

	node.Span.End = Position{Line: endLine, Column: endCol}
	debugPrint("    [parseCast] END: token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)

	return node
}

func (p *Parser) parseCase() ASTNode {
	debugPrint("[parseCase] START: token=%s, literal='%s', line=%d\n",
		p.curToken.Type, p.curToken.Literal, p.curToken.Line)

	node := NewASTNode(CaseNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
	p.nextToken()

	debugPrint("[parseCase] after CASE: token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)

	var isSimpleCase bool
	var caseExpr ASTNode

	if p.curToken.Type != token.WHEN {
		isSimpleCase = true
		debugPrint("[parseCase] Simple CASE, parsing expression\n")
		caseExpr = p.parseExpression()
		if caseExpr.Type != "" {
			node.AddChild(caseExpr)
		}
		debugPrint("[parseCase] after expression: token=%s, literal='%s'\n",
			p.curToken.Type, p.curToken.Literal)
	} else {
		debugPrint("[parseCase] Searched CASE\n")
	}

	whenCount := 0
	for p.curToken.Type == token.WHEN {
		whenCount++
		debugPrint("[parseCase] Parsing WHEN #%d at line %d\n", whenCount, p.curToken.Line)
		whenNode := p.parseWhen(isSimpleCase)
		if whenNode.Type != "" {
			node.AddChild(whenNode)
		}
		debugPrint("[parseCase] after WHEN #%d: token=%s, literal='%s'\n",
			whenCount, p.curToken.Type, p.curToken.Literal)
	}

	if p.curToken.Type == token.ELSE {
		debugPrint("[parseCase] Parsing ELSE\n")
		p.nextToken()
		elseExpr := p.parseExpression()
		if elseExpr.Type != "" {
			elseNode := NewASTNode(BlockNode, "else", elseExpr.Span.Start.Line, elseExpr.Span.Start.Column)
			elseNode.AddChild(elseExpr)
			elseNode.Span = elseExpr.Span
			node.AddChild(elseNode)
		}
		debugPrint("[parseCase] after ELSE: token=%s, literal='%s'\n",
			p.curToken.Type, p.curToken.Literal)
	}

	if p.curToken.Type == token.END {
		debugPrint("[parseCase] Found END, consuming it\n")
		endLine := p.curToken.Line
		endCol := p.curToken.Column + 3
		p.nextToken()
		node.Span.End = Position{Line: endLine, Column: endCol}
	} else {
		p.errors = append(p.errors,
			fmt.Sprintf("expected END in CASE expression, got %s at line %d",
				p.curToken.Type, p.curToken.Line))
	}

	debugPrint("[parseCase] END: token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)

	return node
}

func (p *Parser) parseWhen(isSimpleCase bool) ASTNode {
	debugPrint("  [parseWhen] START: token=%s, literal='%s', isSimple=%v\n",
		p.curToken.Type, p.curToken.Literal, isSimpleCase)

	node := NewASTNode(WhenNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
	p.nextToken()

	debugPrint("  [parseWhen] after WHEN: token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)

	var condition ASTNode

	if isSimpleCase {
		debugPrint("  [parseWhen] Simple CASE: parsing value\n")
		condition = p.parseExpression()
		if condition.Type != "" {
			node.AddChild(condition)
		}
		debugPrint("  [parseWhen] after value: token=%s, literal='%s'\n",
			p.curToken.Type, p.curToken.Literal)
	} else {
		debugPrint("  [parseWhen] Searched CASE: parsing condition\n")
		condition = p.parseExpression()
		if condition.Type != "" {
			node.AddChild(condition)
		}
		debugPrint("  [parseWhen] after condition: token=%s, literal='%s'\n",
			p.curToken.Type, p.curToken.Literal)
	}

	if p.curToken.Type == token.THEN {
		debugPrint("  [parseWhen] Found THEN\n")
		p.nextToken()
	} else {
		p.errors = append(p.errors,
			fmt.Sprintf("expected THEN in CASE WHEN, got %s at line %d",
				p.curToken.Type, p.curToken.Line))
		return node
	}

	debugPrint("  [parseWhen] after THEN: token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)

	result := p.parseExpression()
	if result.Type != "" {
		node.AddChild(result)
		node.Span.End = result.Span.End
	}

	debugPrint("  [parseWhen] END: token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)

	return node
}

func (p *Parser) parseCoalesce() ASTNode {
	debugPrint("    [parseCoalesce] START: token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)

	node := NewASTNode(CoalesceNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
	node.Value = "COALESCE"
	startLine := p.curToken.Line
	startCol := p.curToken.Column
	p.nextToken() // consume COALESCE

	// Expect '('
	if p.curToken.Type != token.LPAREN {
		p.errors = append(p.errors,
			fmt.Sprintf("expected '(' after COALESCE at line %d", p.curToken.Line))
		return node
	}
	p.nextToken() // consume '('

	// Parse arguments (at least 1)
	var args []ASTNode
	for p.curToken.Type != token.RPAREN && p.curToken.Type != token.EOF {
		arg := p.parseExpression()
		if arg.Type != "" {
			args = append(args, arg)
		}
		if p.curToken.Type == token.COMMA {
			p.nextToken()
		}
	}

	if len(args) < 1 {
		p.errors = append(p.errors,
			fmt.Sprintf("COALESCE requires at least 1 argument at line %d", p.curToken.Line))
		return node
	}

	if p.curToken.Type != token.RPAREN {
		p.errors = append(p.errors,
			fmt.Sprintf("expected ')' in COALESCE expression at line %d", p.curToken.Line))
		return node
	}
	endLine := p.curToken.Line
	endCol := p.curToken.Column + 1
	p.nextToken() // consume ')'

	// Add all arguments as children
	for _, arg := range args {
		node.AddChild(arg)
	}

	// Span dari COALESCE sampai ')'
	node.Span.Start = Position{Line: startLine, Column: startCol}
	node.Span.End = Position{Line: endLine, Column: endCol}

	debugPrint("    [parseCoalesce] END: returning COALESCE node with %d args\n", len(args))
	return node
}

// parseNullIf menangani NULLIF(expr1, expr2)
func (p *Parser) parseNullIf() ASTNode {
	debugPrint("    [parseNullIf] START: token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)

	node := NewASTNode(NullIfNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
	node.Value = "NULLIF"
	startLine := p.curToken.Line
	startCol := p.curToken.Column
	p.nextToken() // consume NULLIF

	// Expect '('
	if p.curToken.Type != token.LPAREN {
		p.errors = append(p.errors,
			fmt.Sprintf("expected '(' after NULLIF at line %d", p.curToken.Line))
		return node
	}
	p.nextToken() // consume '('

	// Parse first argument
	arg1 := p.parseExpression()
	if arg1.Type == "" {
		p.errors = append(p.errors,
			fmt.Sprintf("expected expression in NULLIF at line %d", p.curToken.Line))
		return node
	}
	node.AddChild(arg1)

	// Expect ','
	if p.curToken.Type != token.COMMA {
		p.errors = append(p.errors,
			fmt.Sprintf("expected ',' in NULLIF expression at line %d", p.curToken.Line))
		return node
	}
	p.nextToken() // consume ','

	// Parse second argument
	arg2 := p.parseExpression()
	if arg2.Type == "" {
		p.errors = append(p.errors,
			fmt.Sprintf("expected expression in NULLIF at line %d", p.curToken.Line))
		return node
	}
	node.AddChild(arg2)

	if p.curToken.Type != token.RPAREN {
		p.errors = append(p.errors,
			fmt.Sprintf("expected ')' in NULLIF expression at line %d", p.curToken.Line))
		return node
	}
	endLine := p.curToken.Line
	endCol := p.curToken.Column + 1
	p.nextToken() // consume ')'

	// Span dari NULLIF sampai ')'
	node.Span.Start = Position{Line: startLine, Column: startCol}
	node.Span.End = Position{Line: endLine, Column: endCol}

	debugPrint("    [parseNullIf] END: returning NULLIF node\n")
	return node
}

```

---

### File: `pkg/parser/parser_expr.go`

```go
package parser

import (
	"esql-ast-tool/internal/token"
	"fmt"
)

// Helper: combine span dari left dan right
func combineSpan(left, right ASTNode) Span {
	return Span{
		Start: left.Span.Start,
		End:   right.Span.End,
	}
}

// Helper: combine span dari multiple nodes
func combineSpans(nodes ...ASTNode) Span {
	if len(nodes) == 0 {
		return Span{}
	}
	span := nodes[0].Span
	for _, n := range nodes[1:] {
		if n.Span.End.Line > span.End.Line ||
			(n.Span.End.Line == span.End.Line && n.Span.End.Column > span.End.Column) {
			span.End = n.Span.End
		}
	}
	return span
}

func (p *Parser) parseExpression() ASTNode {
	debugPrint("  [parseExpression] token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)

	// STOP - jangan parse jika token bukan bagian dari expression
	if p.curToken.Type == token.END || p.curToken.Type == token.THEN ||
		p.curToken.Type == token.ELSE || p.curToken.Type == token.WHEN ||
		p.curToken.Type == token.EOF ||
		p.curToken.Type == token.RETURNS ||
		p.curToken.Type == token.BEGIN ||
		p.curToken.Type == token.MODULE ||
		p.curToken.Type == token.FUNCTION ||
		p.curToken.Type == token.PROCEDURE {
		debugPrint("  [parseExpression] STOP: token is %s, consuming...\n", p.curToken.Type)
		p.nextToken() // ← CONSUME token biar gak loop!
		return ASTNode{}
	}

	left := p.parseLogicalOr()

	if p.inSet {
		debugPrint("  [parseExpression] inSet=true, skipping '=' comparison\n")
		return left
	}

	if p.curToken.Type == token.EQ || p.curToken.Type == token.ASSIGN ||
		p.curToken.Type == token.NOT_EQ ||
		p.curToken.Type == token.LT || p.curToken.Type == token.GT ||
		p.curToken.Type == token.LTE || p.curToken.Type == token.GTE {
		debugPrint("  [parseExpression] FOUND OPERATOR: %s\n", p.curToken.Literal)
		tok := p.curToken
		p.nextToken()
		right := p.parseAdditive()
		if right.Type != "" {
			compNode := NewASTNode(ComparisonNode, tok.Literal, tok.Line, tok.Column)
			compNode.AddChild(left)
			compNode.AddChild(right)
			compNode.Span = combineSpan(left, right)
			debugPrint("  [parseExpression] returning comparison node\n")
			return compNode
		}
	}

	return left
}

func (p *Parser) parseLogicalOr() ASTNode {
	debugPrint("    [parseLogicalOr] token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)
	node := p.parseLogicalAnd()

	for p.curToken.Type == token.OR {
		tok := p.curToken
		p.nextToken()
		right := p.parseLogicalAnd()
		if right.Type != "" {
			binOp := NewASTNode(BinaryOpNode, tok.Literal, tok.Line, tok.Column)
			binOp.AddChild(node)
			binOp.AddChild(right)
			binOp.Span = combineSpan(node, right)
			node = binOp
		}
	}

	return node
}

func (p *Parser) parseLogicalAnd() ASTNode {
	debugPrint("    [parseLogicalAnd] token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)
	node := p.parseComparison()

	for p.curToken.Type == token.AND {
		tok := p.curToken
		p.nextToken()
		right := p.parseComparison()
		if right.Type != "" {
			binOp := NewASTNode(BinaryOpNode, tok.Literal, tok.Line, tok.Column)
			binOp.AddChild(node)
			binOp.AddChild(right)
			binOp.Span = combineSpan(node, right)
			node = binOp
		}
	}

	return node
}

func (p *Parser) parseComparison() ASTNode {
	debugPrint("    [parseComparison] START: token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)

	node := p.parseAdditive()

	debugPrint("    [parseComparison] after parseAdditive: token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)

	// Cek operator-operator yang mungkin muncul setelah node
	return p.parseComparisonSuffix(node)
}

// parseComparisonSuffix menangani operator setelah node kiri
func (p *Parser) parseComparisonSuffix(left ASTNode) ASTNode {
	switch p.curToken.Type {
	case token.ISNULL, token.NOTNULL:
		return p.parseIsNull(left)

	case token.NOT:
		return p.parseNotOperator(left)

	case token.BETWEEN:
		return p.parseBetween(left, false)

	case token.LIKE:
		return p.parseLike(left, false)

	case token.IN: // ← Tambahkan ini
		return p.parseIn(left, false)

	case token.EQ, token.NOT_EQ, token.LT, token.GT, token.LTE, token.GTE:
		return p.parseComparisonOperator(left)

	default:
		return left
	}
}

// parseIsNull menangani IS NULL / IS NOT NULL
func (p *Parser) parseIsNull(left ASTNode) ASTNode {
	debugPrint("    [parseIsNull] found IS NULL/NOT NULL: %s\n", p.curToken.Literal)
	tok := p.curToken
	p.nextToken()

	var nullNode ASTNode
	if tok.Type == token.ISNULL {
		nullNode = NewASTNode(IsNullNode, "IS NULL", tok.Line, tok.Column)
	} else {
		nullNode = NewASTNode(IsNotNullNode, "IS NOT NULL", tok.Line, tok.Column)
	}
	nullNode.AddChild(left)
	nullNode.Span = combineSpan(left, nullNode)
	debugPrint("    [parseIsNull] returning IS NULL/NOT NULL node\n")
	return nullNode
}

// parseNotOperator menangani NOT (termasuk NOT BETWEEN, NOT LIKE)
func (p *Parser) parseNotOperator(left ASTNode) ASTNode {
	debugPrint("    [parseNotOperator] found NOT, checking next token\n")
	tok := p.curToken
	pos := p.position
	p.nextToken() // consume NOT

	switch p.curToken.Type {
	case token.BETWEEN:
		debugPrint("    [parseNotOperator] found NOT BETWEEN\n")
		return p.parseBetween(left, true)

	case token.LIKE:
		debugPrint("    [parseNotOperator] found NOT LIKE\n")
		return p.parseLike(left, true)

	case token.IN: // ← Tambahkan ini
		debugPrint("    [parseNotOperator] found NOT IN\n")
		return p.parseIn(left, true)

	default:
		// Not NOT BETWEEN/LIKE/IN, treat as unary NOT
		debugPrint("    [parseNotOperator] NOT followed by %s, treating as unary NOT\n", p.curToken.Type)
		p.position = pos
		p.curToken = p.tokens[p.position]

		tok = p.curToken
		p.nextToken()
		right := p.parseComparison()
		if right.Type != "" {
			unaryNode := NewASTNode(UnaryOpNode, tok.Literal, tok.Line, tok.Column)
			unaryNode.AddChild(right)
			unaryNode.Span = combineSpan(
				NewASTNode(IdentifierNode, tok.Literal, tok.Line, tok.Column),
				right,
			)
			return unaryNode
		}
		return left
	}
}

// parseBetween menangani BETWEEN / NOT BETWEEN
func (p *Parser) parseBetween(left ASTNode, isNot bool) ASTNode {
	debugPrint("    [parseBetween] parsing BETWEEN (not=%v)\n", isNot)
	tok := p.curToken
	if isNot {
		p.nextToken() // consume BETWEEN (sudah di-consume di parseNotOperator)
	} else {
		p.nextToken() // consume BETWEEN
	}

	lower := p.parseAdditive()
	if lower.Type == "" {
		p.errors = append(p.errors,
			fmt.Sprintf("expected lower bound in BETWEEN expression at line %d", tok.Line))
		return left
	}

	if p.curToken.Type != token.AND {
		p.errors = append(p.errors,
			fmt.Sprintf("expected AND in BETWEEN expression, got %s at line %d",
				p.curToken.Type, p.curToken.Line))
		return left
	}
	p.nextToken()

	upper := p.parseAdditive()
	if upper.Type == "" {
		p.errors = append(p.errors,
			fmt.Sprintf("expected upper bound in BETWEEN expression at line %d", tok.Line))
		return left
	}

	betweenNode := NewASTNode(BetweenNode, tok.Literal, tok.Line, tok.Column)
	betweenNode.Not = isNot
	betweenNode.AddChild(left)
	betweenNode.AddChild(lower)
	betweenNode.AddChild(upper)
	betweenNode.Span = combineSpans(left, lower, upper)

	debugPrint("    [parseBetween] returning BETWEEN node (not=%v)\n", isNot)
	return betweenNode
}

// parseLike menangani LIKE / NOT LIKE
func (p *Parser) parseLike(left ASTNode, isNot bool) ASTNode {
	debugPrint("    [parseLike] parsing LIKE (not=%v)\n", isNot)
	tok := p.curToken
	if isNot {
		p.nextToken() // consume LIKE (sudah di-consume di parseNotOperator)
	} else {
		p.nextToken() // consume LIKE
	}

	pattern := p.parseAdditive()
	if pattern.Type == "" {
		p.errors = append(p.errors,
			fmt.Sprintf("expected pattern in LIKE expression at line %d", tok.Line))
		return left
	}

	likeNode := NewASTNode(LikeNode, tok.Literal, tok.Line, tok.Column)
	likeNode.Not = isNot
	likeNode.AddChild(left)
	likeNode.AddChild(pattern)
	likeNode.Span = combineSpans(left, pattern)

	debugPrint("    [parseLike] returning LIKE node (not=%v)\n", isNot)
	return likeNode
}

// parseComparisonOperator menangani operator comparison biasa (=, <, >, dll)
func (p *Parser) parseComparisonOperator(left ASTNode) ASTNode {
	debugPrint("    [parseComparisonOperator] found operator: %s\n", p.curToken.Literal)
	tok := p.curToken
	p.nextToken()
	right := p.parseAdditive()
	if right.Type != "" {
		compNode := NewASTNode(ComparisonNode, tok.Literal, tok.Line, tok.Column)
		compNode.AddChild(left)
		compNode.AddChild(right)
		compNode.Span = combineSpan(left, right)
		debugPrint("    [parseComparisonOperator] returning comparison node\n")
		return compNode
	}
	return left
}

func (p *Parser) parseAdditive() ASTNode {
	debugPrint("    [parseAdditive] START: token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)

	node := p.parseMultiplicative()

	for p.curToken.Type == token.PLUS || p.curToken.Type == token.MINUS {
		tok := p.curToken
		p.nextToken()
		right := p.parseMultiplicative()
		binOp := NewASTNode(BinaryOpNode, tok.Literal, tok.Line, tok.Column)
		binOp.AddChild(node)
		binOp.AddChild(right)
		binOp.Span = combineSpan(node, right)
		node = binOp
	}

	debugPrint("    [parseAdditive] END: token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)

	return node
}

func (p *Parser) parseMultiplicative() ASTNode {
	debugPrint("    [parseMultiplicative] START: token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)

	node := p.parseUnary()

	for p.curToken.Type == token.ASTERISK || p.curToken.Type == token.SLASH ||
		p.curToken.Type == token.MODULO {
		tok := p.curToken
		p.nextToken()
		right := p.parseUnary()
		binOp := NewASTNode(BinaryOpNode, tok.Literal, tok.Line, tok.Column)
		binOp.AddChild(node)
		binOp.AddChild(right)
		binOp.Span = combineSpan(node, right)
		node = binOp
	}

	debugPrint("    [parseMultiplicative] END: token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)

	return node
}

func (p *Parser) parseUnary() ASTNode {
	debugPrint("    [parseUnary] START: token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)

	if p.curToken.Type == token.MINUS {
		tok := p.curToken
		p.nextToken()
		operand := p.parsePrimary()
		unaryNode := NewASTNode(UnaryOpNode, tok.Literal, tok.Line, tok.Column)
		unaryNode.AddChild(operand)
		unaryNode.Span = combineSpan(
			NewASTNode(IdentifierNode, tok.Literal, tok.Line, tok.Column),
			operand,
		)
		debugPrint("    [parseUnary] END (unary minus): token=%s, literal='%s'\n",
			p.curToken.Type, p.curToken.Literal)
		return unaryNode
	}

	// Note: NOT is now handled in parseComparison, but keep this for safety
	if p.curToken.Type == token.NOT {
		tok := p.curToken
		p.nextToken()
		operand := p.parseComparison()
		unaryNode := NewASTNode(UnaryOpNode, tok.Literal, tok.Line, tok.Column)
		unaryNode.AddChild(operand)
		unaryNode.Span = combineSpan(
			NewASTNode(IdentifierNode, tok.Literal, tok.Line, tok.Column),
			operand,
		)
		debugPrint("    [parseUnary] END (unary not): token=%s, literal='%s'\n",
			p.curToken.Type, p.curToken.Literal)
		return unaryNode
	}

	result := p.parsePrimary()
	debugPrint("    [parseUnary] END (primary): token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)
	return result
}

func (p *Parser) parseIn(left ASTNode, isNot bool) ASTNode {
	debugPrint("    [parseIn] parsing IN (not=%v)\n", isNot)
	tok := p.curToken
	if isNot {
		p.nextToken() // consume IN (sudah di-consume di parseNotOperator)
	} else {
		p.nextToken() // consume IN
	}

	// Expect '('
	if p.curToken.Type != token.LPAREN {
		p.errors = append(p.errors,
			fmt.Sprintf("expected '(' after IN expression at line %d", tok.Line))
		return left
	}
	p.nextToken() // consume '('

	// Parse list of values
	var values []ASTNode
	for p.curToken.Type != token.RPAREN && p.curToken.Type != token.EOF {
		val := p.parseExpression()
		if val.Type != "" {
			values = append(values, val)
		}
		if p.curToken.Type == token.COMMA {
			p.nextToken()
		}
	}

	if p.curToken.Type != token.RPAREN {
		p.errors = append(p.errors,
			fmt.Sprintf("expected ')' in IN expression at line %d", tok.Line))
		return left
	}
	endLine := p.curToken.Line
	endCol := p.curToken.Column + 1
	p.nextToken() // consume ')'

	// Buat InNode
	inNode := NewASTNode(InNode, tok.Literal, tok.Line, tok.Column)
	inNode.Not = isNot
	inNode.AddChild(left)

	// Tambahkan semua values sebagai children
	for _, val := range values {
		inNode.AddChild(val)
	}

	// Span dari left sampai akhir ')'
	if len(values) > 0 {
		inNode.Span = combineSpans(append([]ASTNode{left}, values...)...)
		inNode.Span.End = Position{Line: endLine, Column: endCol}
	} else {
		inNode.Span = combineSpan(left, inNode)
	}

	debugPrint("    [parseIn] returning IN node (not=%v, values=%d)\n", isNot, len(values))
	return inNode
}

```

---

### File: `pkg/parser/parser.go`

```go
package parser

import (
	"esql-ast-tool/internal/token"
	"fmt"
)

var DebugMode = false

func debugPrint(format string, args ...interface{}) {
	if DebugMode {
		fmt.Printf(format, args...)
	}
}

type Parser struct {
	l         *Lexer
	tokens    []token.Token
	position  int
	curToken  token.Token
	peekToken token.Token
	errors    []string
	inSet     bool
}

func NewParser(input string) *Parser {
	l := NewLexer(input)
	tokens := l.Tokenize()
	p := &Parser{
		l:        l,
		tokens:   tokens,
		position: 0,
		errors:   []string{},
	}

	if len(tokens) > 0 {
		p.curToken = tokens[0]
		if len(tokens) > 1 {
			p.peekToken = tokens[1]
		}
	}

	return p
}

func (p *Parser) nextToken() {
	p.position++
	if p.position < len(p.tokens) {
		p.curToken = p.tokens[p.position]
		if p.position+1 < len(p.tokens) {
			p.peekToken = p.tokens[p.position+1]
		} else {
			p.peekToken = token.Token{Type: token.EOF}
		}
	} else {
		p.curToken = token.Token{Type: token.EOF}
		p.peekToken = token.Token{Type: token.EOF}
	}
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) GetCurToken() token.Token {
	return p.curToken
}

func (p *Parser) GetPeekToken() token.Token {
	return p.peekToken
}

func (p *Parser) GetNextToken() {
	p.nextToken()
}

func (p *Parser) GetPosition() int {
	return p.position
}

func (p *Parser) DebugTokens() {
	for i, tok := range p.tokens {
		fmt.Printf("[%d] Type: %s, Literal: '%s', Line: %d, Col: %d\n",
			i, tok.Type, tok.Literal, tok.Line, tok.Column)
	}
}

func (p *Parser) GetTokens() []token.Token {
	return p.tokens
}

func (p *Parser) DebugPrintTokens() {
	fmt.Println("=== TOKENS ===")
	for i, tok := range p.tokens {
		fmt.Printf("[%d] Type: %s, Literal: '%s', Line: %d, Col: %d\n",
			i, tok.Type, tok.Literal, tok.Line, tok.Column)
	}
}

// ParseProgram - entry point
func (p *Parser) ParseProgram() Program {
	program := NewProgram()
	rootNode := NewASTNode(ProgramNode, "PROGRAM", 1, 1)

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt.Type != "" {
			rootNode.AddChild(stmt)
			program.AddStatement(stmt) // <-- tambahkan ini
		}

		if p.curToken.Type == token.SEMICOLON {
			p.nextToken()
		} else if p.curToken.Type != token.EOF && p.curToken.Type != token.END {
			p.nextToken()
		}
	}

	program.Root = rootNode
	return program
}

```

---

### File: `pkg/parser/parser_stmt_loop.go`

```go
package parser

import (
	"esql-ast-tool/internal/token"
)

func (p *Parser) parseWhile() ASTNode {
	node := NewASTNode(WhileNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
	p.nextToken()

	condition := p.parseExpression()
	if condition.Type != "" {
		node.AddChild(condition)
	}

	if p.curToken.Type == token.DO {
		p.nextToken()
	}

	bodyNode := NewASTNode(BlockNode, "body", p.curToken.Line, p.curToken.Column)
	for p.curToken.Type != token.END && p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt.Type != "" {
			bodyNode.AddChild(stmt)
		}
		// parseStatement sudah memajukan token, tidak perlu p.nextToken() lagi
	}
	node.AddChild(bodyNode)

	if p.curToken.Type == token.END {
		p.nextToken()
		if p.curToken.Type == token.WHILE {
			p.nextToken()
		}
	}

	if p.curToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return node
}

func (p *Parser) parseFor() ASTNode {
	node := NewASTNode(ForNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
	p.nextToken()

	if p.curToken.Type == token.IDENTIFIER {
		varNode := NewASTNode(IdentifierNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
		varNode.Value = varNode.Token
		node.AddChild(varNode)
		p.nextToken()
	}

	if p.curToken.Type == token.AS {
		p.nextToken()
	}

	for p.curToken.Type != token.DO && p.curToken.Type != token.EOF {
		expr := p.parseExpression()
		if expr.Type != "" {
			node.AddChild(expr)
		}
	}

	if p.curToken.Type == token.DO {
		p.nextToken()
	}

	bodyNode := NewASTNode(BlockNode, "body", p.curToken.Line, p.curToken.Column)
	for p.curToken.Type != token.END && p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt.Type != "" {
			bodyNode.AddChild(stmt)
		}
		if p.curToken.Type != token.END && p.curToken.Type != token.EOF {
			p.nextToken()
		}
	}
	node.AddChild(bodyNode)

	if p.curToken.Type == token.END {
		p.nextToken()
		if p.curToken.Type == token.FOR {
			p.nextToken()
		}
	}

	if p.curToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return node
}

```

---

### File: `pkg/parser/parser_stmt_if.go`

```go
package parser

import (
	"fmt"

	"esql-ast-tool/internal/token"
)

func (p *Parser) parseIf() ASTNode {
	debugPrint("[parseIf] START: token=%s, literal='%s', line=%d\n",
		p.curToken.Type, p.curToken.Literal, p.curToken.Line)

	node := NewASTNode(IfNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
	p.nextToken() // consume IF

	debugPrint("[parseIf] after IF: token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)

	// Parse condition - parseExpression will handle NOT
	cond := p.parseExpression()
	if cond.Type != "" {
		debugPrint("[parseIf] condition parsed: type=%s\n", cond.Type)
		node.AddChild(cond)
	} else {
		debugPrint("[parseIf] WARNING: empty condition\n")
	}

	debugPrint("[parseIf] after condition: token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)

	// Expect THEN
	if p.curToken.Type == token.THEN {
		debugPrint("[parseIf] Found THEN, consuming it\n")
		p.nextToken() // consume THEN
	} else {
		debugPrint("[parseIf] ERROR: expected THEN, got %s\n", p.curToken.Type)
		p.errors = append(p.errors,
			fmt.Sprintf("expected THEN after IF condition, got %s at line %d",
				p.curToken.Type, p.curToken.Line))
		if p.curToken.Type != token.EOF {
			p.nextToken()
		}
		return node
	}

	debugPrint("[parseIf] after THEN: token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)

	// Parse THEN block
	thenBlock := NewASTNode(BlockNode, "then", p.curToken.Line, p.curToken.Column)
	for p.curToken.Type != token.END && p.curToken.Type != token.ELSE && p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt.Type != "" {
			thenBlock.AddChild(stmt)
		}
		if p.curToken.Type == token.SEMICOLON {
			p.nextToken()
		}
	}
	node.AddChild(thenBlock)

	// Parse ELSE if ada
	if p.curToken.Type == token.ELSE {
		debugPrint("[parseIf] Found ELSE\n")
		p.nextToken()
		elseBlock := NewASTNode(BlockNode, "else", p.curToken.Line, p.curToken.Column)
		for p.curToken.Type != token.END && p.curToken.Type != token.EOF {
			stmt := p.parseStatement()
			if stmt.Type != "" {
				elseBlock.AddChild(stmt)
			}
			if p.curToken.Type == token.SEMICOLON {
				p.nextToken()
			}
		}
		node.AddChild(elseBlock)
	}

	// Konsumsi END IF
	if p.curToken.Type == token.END {
		debugPrint("[parseIf] Found END\n")
		endLine := p.curToken.Line
		endCol := p.curToken.Column
		p.nextToken()
		if p.curToken.Type == token.IF {
			debugPrint("[parseIf] Found IF after END\n")
			endCol = p.curToken.Column + len(p.curToken.Literal)
			p.nextToken()
		}
		// Update end span
		node.Span.End = Position{Line: endLine, Column: endCol}
	}

	if p.curToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	debugPrint("[parseIf] END: token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)

	return node
}

```

---

### File: `pkg/parser/parser_stmt_call.go`

```go
package parser

import (
	"esql-ast-tool/internal/token"
)

func (p *Parser) parseCall() ASTNode {
	node := NewASTNode(CallNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
	p.nextToken()

	// Parse procedure name
	if p.curToken.Type == token.IDENTIFIER {
		nameNode := NewASTNode(IdentifierNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
		node.AddChild(nameNode)
		p.nextToken()
	}

	// Parse arguments
	if p.curToken.Type == token.LPAREN {
		p.nextToken()
		for p.curToken.Type != token.RPAREN && p.curToken.Type != token.EOF {
			arg := p.parseExpression()
			if arg.Type != "" {
				node.AddChild(arg)
			}
			if p.curToken.Type == token.COMMA {
				p.nextToken()
			}
		}
		if p.curToken.Type == token.RPAREN {
			p.nextToken()
		}
	}

	if p.curToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return node
}

```

---

### File: `pkg/parser/parser_stmt_declare.go`

```go
package parser

import (
	"esql-ast-tool/internal/token"
)

func (p *Parser) parseDeclare() ASTNode {
	node := NewASTNode(DeclareNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
	p.nextToken()

	// Parse variable name
	if p.curToken.Type != token.IDENTIFIER {
		p.errors = append(p.errors, "expected identifier after DECLARE")
		return node
	}

	nameNode := NewASTNode(IdentifierNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
	nameNode.Value = p.curToken.Literal
	node.AddChild(nameNode)
	p.nextToken()

	// Parse type (INTEGER, STRING, etc.)
	if p.curToken.Type == token.IDENTIFIER {
		typeNode := NewASTNode(IdentifierNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
		typeNode.Value = p.curToken.Literal
		node.AddChild(typeNode)
		p.nextToken()
	}

	// Optional DEFAULT
	if p.curToken.Type == token.DEFAULT {
		p.nextToken()
		expr := p.parseExpression()
		if expr.Type != "" {
			node.AddChild(expr)
		}
	}

	// Consume semicolon
	if p.curToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return node
}

```

---

### File: `pkg/parser/parser_stmt.go`

```go
package parser

import (
	"esql-ast-tool/internal/token"
)

func (p *Parser) parseStatement() ASTNode {
	switch p.curToken.Type {
	case token.CREATE:
		return p.parseCreate()
	case token.DECLARE:
		return p.parseDeclare()
	case token.SET:
		return p.parseSet()
	case token.IF:
		return p.parseIf()
	case token.WHILE:
		return p.parseWhile()
	case token.FOR:
		return p.parseFor()
	case token.RETURN:
		return p.parseReturn()
	case token.THROW:
		return p.parseThrow()
	case token.PROPAGATE:
		return p.parsePropagate()
	case token.MOVE:
		return p.parseMove()
	case token.CONTINUE:
		return p.parseContinue()
	case token.BREAK:
		return p.parseBreak()
	case token.LABEL:
		return p.parseLabel()
	case token.MODULE:
		return p.parseModuleStatement()
	case token.FUNCTION:
		return p.parseFunctionStatement()
	case token.PROCEDURE:
		return p.parseProcedureStatement()
	case token.CALL:
		return p.parseCall()
	case token.END:
		p.nextToken()
		return ASTNode{}
	case token.WHEN, token.ELSE, token.THEN:
		p.nextToken()
		return ASTNode{}
	default:
		expr := p.parseExpression()
		if p.curToken.Type == token.SEMICOLON {
			p.nextToken()
		}
		return expr
	}
}

```

---

### File: `pkg/parser/parser_stmt_control.go`

```go
package parser

import (
	"esql-ast-tool/internal/token"
)

// ============================================
// RETURN
// ============================================

func (p *Parser) parseReturn() ASTNode {
	node := NewASTNode(ReturnNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
	p.nextToken()

	if p.curToken.Type != token.SEMICOLON {
		expr := p.parseExpression()
		if expr.Type != "" {
			node.AddChild(expr)
		}
	}

	if p.curToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return node
}

// ============================================
// THROW
// ============================================

func (p *Parser) parseThrow() ASTNode {
	node := NewASTNode(ThrowNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
	p.nextToken()

	if p.curToken.Type != token.SEMICOLON {
		expr := p.parseExpression()
		if expr.Type != "" {
			node.AddChild(expr)
		}
	}

	if p.curToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return node
}

// ============================================
// PROPAGATE
// ============================================

func (p *Parser) parsePropagate() ASTNode {
	node := NewASTNode(PropagateNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
	p.nextToken() // consume PROPAGATE

	// Parse expressions sampai semicolon
	// Bisa multiple expressions dipisah koma
	for p.curToken.Type != token.SEMICOLON && p.curToken.Type != token.EOF {
		expr := p.parseExpression()
		if expr.Type != "" {
			node.AddChild(expr)
		}
		// Jika ketemu koma, lanjut ke expression berikutnya
		if p.curToken.Type == token.COMMA {
			p.nextToken()
		}
	}

	if p.curToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return node
}

// ============================================
// MOVE
// ============================================

func (p *Parser) parseMove() ASTNode {
	node := NewASTNode(MoveNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
	p.nextToken()

	target := p.parseExpression()
	if target.Type != "" {
		node.AddChild(target)
	}

	if p.curToken.Type == token.TO {
		p.nextToken()
		source := p.parseExpression()
		if source.Type != "" {
			node.AddChild(source)
		}
	}

	if p.curToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return node
}

// ============================================
// CONTINUE
// ============================================

func (p *Parser) parseContinue() ASTNode {
	node := NewASTNode(ContinueNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
	if p.curToken.Type == token.IDENTIFIER {
		node.Value = p.curToken.Literal
		p.nextToken()
	}

	if p.curToken.Type == token.SEMICOLON {
		p.nextToken()
	}
	return node
}

// ============================================
// BREAK
// ============================================

func (p *Parser) parseBreak() ASTNode {
	node := NewASTNode(BreakNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
	p.nextToken()

	if p.curToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return node
}

// ============================================
// LABEL
// ============================================

func (p *Parser) parseLabel() ASTNode {
	node := NewASTNode(LabelNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
	p.nextToken()

	if p.curToken.Type == token.IDENTIFIER {
		node.Value = p.curToken.Literal
		p.nextToken()
	}

	if p.curToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return node
}

```

---

### File: `pkg/parser/lexer.go`

```go
package parser

import (
	"strings"

	"esql-ast-tool/internal/token"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
	line         int
	column       int
}

func NewLexer(input string) *Lexer {
	l := &Lexer{
		input:  input,
		line:   1,
		column: 0,
	}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
	l.column++
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func (l *Lexer) peekIdentifier() string {
	pos := l.position
	oldPos := l.position
	oldReadPos := l.readPosition
	oldCh := l.ch

	for isLetter(l.ch) || isDigit(l.ch) || l.ch == '_' {
		l.readChar()
	}
	ident := l.input[pos:l.position]

	l.position = oldPos
	l.readPosition = oldReadPos
	l.ch = oldCh

	return ident
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		if l.ch == '\n' {
			l.line++
			l.column = 0
		}
		l.readChar()
	}
}

func (l *Lexer) skipComments() {
	if l.ch == '-' && l.peekChar() == '-' {
		for l.ch != '\n' && l.ch != 0 {
			l.readChar()
		}
		l.skipWhitespace()
	} else if l.ch == '/' && l.peekChar() == '*' {
		l.readChar()
		l.readChar()
		for !(l.ch == '*' && l.peekChar() == '/') && l.ch != 0 {
			if l.ch == '\n' {
				l.line++
				l.column = 0
			}
			l.readChar()
		}
		if l.ch != 0 {
			l.readChar()
			l.readChar()
		}
		l.skipWhitespace()
	}
}

func (l *Lexer) NextToken() token.Token {
	l.skipWhitespace()
	l.skipComments()

	var tok token.Token
	tok.Line = l.line
	tok.Column = l.column
	tok.Pos = l.position

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok.Type = token.EQ
			tok.Literal = string(ch) + string(l.ch)
		} else {
			tok.Type = token.ASSIGN
			tok.Literal = string(l.ch)
		}
	case '+':
		tok.Type = token.PLUS
		tok.Literal = string(l.ch)
	case '-':
		tok.Type = token.MINUS
		tok.Literal = string(l.ch)
	case '*':
		tok.Type = token.ASTERISK
		tok.Literal = string(l.ch)
	case '/':
		tok.Type = token.SLASH
		tok.Literal = string(l.ch)
	case '%':
		tok.Type = token.MODULO
		tok.Literal = string(l.ch)
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok.Type = token.NOT_EQ
			tok.Literal = string(ch) + string(l.ch)
		} else {
			tok.Type = token.ILLEGAL
			tok.Literal = string(l.ch)
		}
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok.Type = token.LTE
			tok.Literal = string(ch) + string(l.ch)
		} else {
			tok.Type = token.LT
			tok.Literal = string(l.ch)
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok.Type = token.GTE
			tok.Literal = string(ch) + string(l.ch)
		} else {
			tok.Type = token.GT
			tok.Literal = string(l.ch)
		}
	case ',':
		tok.Type = token.COMMA
		tok.Literal = string(l.ch)
	case ';':
		tok.Type = token.SEMICOLON
		tok.Literal = string(l.ch)
	case '(':
		tok.Type = token.LPAREN
		tok.Literal = string(l.ch)
	case ')':
		tok.Type = token.RPAREN
		tok.Literal = string(l.ch)
	case '{':
		tok.Type = token.LBRACE
		tok.Literal = string(l.ch)
	case '}':
		tok.Type = token.RBRACE
		tok.Literal = string(l.ch)
	case '[':
		tok.Type = token.LBRACKET
		tok.Literal = string(l.ch)
	case ']':
		tok.Type = token.RBRACKET
		tok.Literal = string(l.ch)
	case '.':
		tok.Type = token.DOT
		tok.Literal = string(l.ch)
	case 0:
		tok.Type = token.EOF
		tok.Literal = ""
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(strings.ToUpper(tok.Literal))

			// Handle IS NULL / IS NOT NULL
			if tok.Type == token.IS {
				l.skipWhitespace()
				if strings.ToUpper(l.peekIdentifier()) == "NULL" {
					l.readIdentifier()
					tok.Type = token.ISNULL
					tok.Literal = "IS NULL"
				} else if strings.ToUpper(l.peekIdentifier()) == "NOT" {
					l.readIdentifier()
					l.skipWhitespace()
					if strings.ToUpper(l.peekIdentifier()) == "NULL" {
						l.readIdentifier()
						tok.Type = token.NOTNULL
						tok.Literal = "IS NOT NULL"
					}
				}
				return tok
			}

			// Handle NOT BETWEEN - tapi jangan gabung jadi satu token
			// Biarkan NOT dan BETWEEN sebagai token terpisah
			if tok.Type == token.NOT {
				// Cek apakah next token adalah BETWEEN
				l.skipWhitespace()
				if strings.ToUpper(l.peekIdentifier()) == "BETWEEN" {
					// Return NOT, next token akan BETWEEN
					return tok
				}
			}

			return tok
		} else if isDigit(l.ch) {
			tok.Type = token.NUMBER
			tok.Literal = l.readNumber()
			return tok
		} else if l.ch == '\'' {
			tok.Type = token.STRING
			tok.Literal = l.readString()
			return tok
		} else {
			tok.Type = token.ILLEGAL
			tok.Literal = string(l.ch)
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) || l.ch == '_' {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) || l.ch == '.' {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readString() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '\'' || l.ch == 0 {
			break
		}
	}
	str := l.input[position:l.position]
	l.readChar()
	return str
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) Tokenize() []token.Token {
	var tokens []token.Token
	for {
		tok := l.NextToken()
		if tok.Type == token.EOF {
			break
		}
		tokens = append(tokens, tok)
	}
	return tokens
}

```

---

### File: `pkg/parser/parser_utils.go`

```go
package parser

import (
	"esql-ast-tool/internal/token"
)

func (p *Parser) parseFieldReferenceFromKeyword(keyword string) ASTNode {
	debugPrint("  [parseFieldReferenceFromKeyword] keyword=%s\n", keyword)

	baseNode := NewASTNode(IdentifierNode, keyword, p.curToken.Line, p.curToken.Column)
	baseNode.Value = keyword
	p.nextToken()

	return p.parseFieldReference(baseNode)
}

func (p *Parser) parseField() ASTNode {
	debugPrint("  [parseField] START: token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)

	node := NewASTNode(FieldReferenceNode, "FIELD", p.curToken.Line, p.curToken.Column)
	p.nextToken()

	if p.curToken.Type == token.IDENTIFIER || p.curToken.Type == token.ENVIRONMENT {
		base := p.parseIdentifier()
		node.AddChild(base)

		for p.curToken.Type == token.DOT {
			node = p.parseFieldReference(node)
		}
	}

	debugPrint("  [parseField] END: token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)

	return node
}

func (p *Parser) parseExpressionStatement() ASTNode {
	expr := p.parseExpression()
	if p.curToken.Type == token.SEMICOLON {
		p.nextToken()
	}
	return expr
}

func (p *Parser) parseComparisonFromNode(left ASTNode) ASTNode {
	node := left

	if p.curToken.Type == token.EQ || p.curToken.Type == token.NOT_EQ ||
		p.curToken.Type == token.LT || p.curToken.Type == token.GT ||
		p.curToken.Type == token.LTE || p.curToken.Type == token.GTE {
		tok := p.curToken
		p.nextToken()
		right := p.parseAdditive()
		if right.Type != "" {
			compNode := NewASTNode(ComparisonNode, tok.Literal, tok.Line, tok.Column)
			compNode.AddChild(node)
			compNode.AddChild(right)
			return compNode
		}
	}

	return node
}

```

---

### File: `pkg/refactor/refactor.go`

```go
package refactor

import (
	"fmt"
	"sort"
	"strings"

	"esql-ast-tool/pkg/analyzer"
	"esql-ast-tool/pkg/parser"
)

type Suggestion struct {
	Type     string // "dead_code", "code_smell", "improvement"
	Severity string // "high", "medium", "low"
	Message  string
	Line     int
	Details  []string
}

type RefactorResult struct {
	Suggestions []Suggestion
	Stats       map[string]int
}

type RenameResult struct {
	OldName     string
	NewName     string
	Occurrences int
	Locations   []RenameLocation
	Success     bool
	Message     string
}

type RenameLocation struct {
	Line    int
	Column  int
	OldText string
	NewText string
	Context string // "DECLARE", "SET", "CALL", etc.
}

// ============================================
// SUGGEST
// ============================================

func Suggest(program parser.Program, analysisResult analyzer.AnalysisResult) RefactorResult {
	var suggestions []Suggestion
	stats := make(map[string]int)

	// 1. Detect Dead Code - Unused Procedures
	for name, info := range analysisResult.Procedures {
		if _, ok := analysisResult.ReverseCallGraph[name]; !ok {
			suggestions = append(suggestions, Suggestion{
				Type:     "dead_code",
				Severity: "high",
				Message:  fmt.Sprintf("Procedure '%s' is never called", name),
				Line:     info.Line,
				Details:  []string{"Remove this procedure or add a CALL statement"},
			})
			stats["dead_procedures"]++
		}
	}

	// 2. Detect Dead Code - Unused Functions
	for name, info := range analysisResult.Functions {
		if _, ok := analysisResult.ReverseCallGraph[name]; !ok && info.ReturnType != "BUILTIN" {
			suggestions = append(suggestions, Suggestion{
				Type:     "dead_code",
				Severity: "high",
				Message:  fmt.Sprintf("Function '%s' is never called", name),
				Line:     info.Line,
				Details:  []string{"Remove this function or add a call"},
			})
			stats["dead_functions"]++
		}
	}

	// 3. Detect Unused Variables
	for name, info := range analysisResult.Variables {
		if !isUsed(analysisResult.UsedVariables, name) {
			suggestions = append(suggestions, Suggestion{
				Type:     "dead_code",
				Severity: "medium",
				Message:  fmt.Sprintf("Variable '%s' is declared but never used", name),
				Line:     info.Line,
				Details:  []string{"Remove this declaration or use it"},
			})
			stats["unused_variables"]++
		}
	}

	// 4. Code Smells - Variables used in many places
	for name, info := range analysisResult.Variables {
		if usages, ok := analysisResult.UsageMap[name]; ok && len(usages) > 5 {
			suggestions = append(suggestions, Suggestion{
				Type:     "code_smell",
				Severity: "low",
				Message:  fmt.Sprintf("Variable '%s' is used in %d places", name, len(usages)),
				Line:     info.Line,
				Details:  []string{fmt.Sprintf("Used at lines: %v", getLines(usages))},
			})
			stats["high_usage_variables"]++
		}
	}

	// 5. Improvements - Single-call procedures
	for name, info := range analysisResult.Procedures {
		if callers, ok := analysisResult.CallGraph[name]; ok && len(callers) == 1 {
			suggestions = append(suggestions, Suggestion{
				Type:     "improvement",
				Severity: "low",
				Message:  fmt.Sprintf("Procedure '%s' only calls one other procedure", name),
				Line:     info.Line,
				Details:  []string{fmt.Sprintf("Calls: %v", callers)},
			})
			stats["single_call_procedures"]++
		}
	}

	return RefactorResult{
		Suggestions: suggestions,
		Stats:       stats,
	}
}

// ============================================
// RENAME
// ============================================

func RenameVariable(program parser.Program, oldName, newName string, dryRun bool) RenameResult {
	var locations []RenameLocation
	occurrences := 0
	declareLines := make(map[int]bool)

	var search func(node parser.ASTNode)
	search = func(node parser.ASTNode) {
		// DECLARE
		if node.Type == parser.DeclareNode {
			if len(node.Children) > 0 && node.Children[0].Type == parser.IdentifierNode {
				name := getNodeName(node.Children[0])
				if name == oldName {
					declareLines[node.Span.Start.Line] = true
					occurrences++
					locations = append(locations, RenameLocation{
						Line:    node.Span.Start.Line,
						Column:  node.Span.Start.Column,
						OldText: oldName,
						NewText: newName,
						Context: "DECLARE",
					})
				}
			}
		}

		// Identifier usage
		if node.Type == parser.IdentifierNode {
			name := getNodeName(node)
			if name == oldName && !declareLines[node.Span.Start.Line] {
				occurrences++
				locations = append(locations, RenameLocation{
					Line:    node.Span.Start.Line,
					Column:  node.Span.Start.Column,
					OldText: oldName,
					NewText: newName,
					Context: "USAGE",
				})
			}
		}

		for _, child := range node.Children {
			search(child)
		}
	}

	for _, stmt := range program.Statements {
		search(stmt)
	}

	if occurrences == 0 {
		return RenameResult{
			OldName:     oldName,
			NewName:     newName,
			Occurrences: 0,
			Locations:   locations,
			Success:     false,
			Message:     fmt.Sprintf("Variable '%s' not found", oldName),
		}
	}

	return RenameResult{
		OldName:     oldName,
		NewName:     newName,
		Occurrences: occurrences,
		Locations:   locations,
		Success:     true,
		Message:     fmt.Sprintf("Renamed '%s' to '%s' in %d locations", oldName, newName, occurrences),
	}
}

func RenameProcedure(program parser.Program, oldName, newName string, dryRun bool) RenameResult {
	var locations []RenameLocation
	occurrences := 0

	var search func(node parser.ASTNode)
	search = func(node parser.ASTNode) {
		// PROCEDURE definition
		if node.Type == parser.ProcedureNode {
			if len(node.Children) > 0 {
				name := getNodeName(node.Children[0])
				if name == oldName {
					occurrences++
					locations = append(locations, RenameLocation{
						Line:    node.Span.Start.Line,
						Column:  node.Span.Start.Column,
						OldText: oldName,
						NewText: newName,
						Context: "PROCEDURE",
					})
				}
			}
		}

		// CALL statement
		if node.Type == parser.CallNode {
			if len(node.Children) > 0 {
				name := getNodeName(node.Children[0])
				if name == oldName {
					occurrences++
					locations = append(locations, RenameLocation{
						Line:    node.Span.Start.Line,
						Column:  node.Span.Start.Column,
						OldText: oldName,
						NewText: newName,
						Context: "CALL",
					})
				}
			}
		}

		for _, child := range node.Children {
			search(child)
		}
	}

	for _, stmt := range program.Statements {
		search(stmt)
	}

	if occurrences == 0 {
		return RenameResult{
			OldName:     oldName,
			NewName:     newName,
			Occurrences: 0,
			Locations:   locations,
			Success:     false,
			Message:     fmt.Sprintf("Procedure '%s' not found", oldName),
		}
	}

	return RenameResult{
		OldName:     oldName,
		NewName:     newName,
		Occurrences: occurrences,
		Locations:   locations,
		Success:     true,
		Message:     fmt.Sprintf("Renamed procedure '%s' to '%s' in %d locations", oldName, newName, occurrences),
	}
}

func RenameFunction(program parser.Program, oldName, newName string, dryRun bool) RenameResult {
	var locations []RenameLocation
	occurrences := 0

	var search func(node parser.ASTNode)
	search = func(node parser.ASTNode) {
		// FUNCTION definition
		if node.Type == parser.FunctionNode {
			if len(node.Children) > 0 {
				name := getNodeName(node.Children[0])
				if name == oldName {
					occurrences++
					locations = append(locations, RenameLocation{
						Line:    node.Span.Start.Line,
						Column:  node.Span.Start.Column,
						OldText: oldName,
						NewText: newName,
						Context: "FUNCTION",
					})
				}
			}
		}

		// FUNCTION CALL
		if node.Type == parser.FunctionCallNode {
			name := getNodeName(node)
			if name == oldName {
				occurrences++
				locations = append(locations, RenameLocation{
					Line:    node.Span.Start.Line,
					Column:  node.Span.Start.Column,
					OldText: oldName,
					NewText: newName,
					Context: "FUNCTION_CALL",
				})
			}
		}

		for _, child := range node.Children {
			search(child)
		}
	}

	for _, stmt := range program.Statements {
		search(stmt)
	}

	if occurrences == 0 {
		return RenameResult{
			OldName:     oldName,
			NewName:     newName,
			Occurrences: 0,
			Locations:   locations,
			Success:     false,
			Message:     fmt.Sprintf("Function '%s' not found", oldName),
		}
	}

	return RenameResult{
		OldName:     oldName,
		NewName:     newName,
		Occurrences: occurrences,
		Locations:   locations,
		Success:     true,
		Message:     fmt.Sprintf("Renamed function '%s' to '%s' in %d locations", oldName, newName, occurrences),
	}
}

// ============================================
// SEARCH & REPLACE HELPERS
// ============================================

func searchAndReplace(nodes []parser.ASTNode, oldName, newName string, locations *[]RenameLocation, occurrences *int, targetType string) {
	for _, node := range nodes {
		// Procedure definition
		if node.Type == parser.ProcedureNode && targetType == "procedure" {
			if len(node.Children) > 0 && node.Children[0].Type == parser.IdentifierNode {
				var name string
				if val, ok := node.Children[0].Value.(string); ok && val != "" {
					name = val
				} else if node.Children[0].Token != "" {
					name = node.Children[0].Token
				}
				if name == oldName {
					*occurrences++
					*locations = append(*locations, RenameLocation{
						Line:    node.Span.Start.Line,
						Column:  node.Span.Start.Column,
						OldText: oldName,
						NewText: newName,
						Context: "PROCEDURE",
					})
				}
			}
		}

		// Function definition
		if node.Type == parser.FunctionNode && targetType == "function" {
			if len(node.Children) > 0 && node.Children[0].Type == parser.IdentifierNode {
				var name string
				if val, ok := node.Children[0].Value.(string); ok && val != "" {
					name = val
				} else if node.Children[0].Token != "" {
					name = node.Children[0].Token
				}
				if name == oldName {
					*occurrences++
					*locations = append(*locations, RenameLocation{
						Line:    node.Span.Start.Line,
						Column:  node.Span.Start.Column,
						OldText: oldName,
						NewText: newName,
						Context: "FUNCTION",
					})
				}
			}
		}

		// Variable declaration
		if node.Type == parser.DeclareNode && targetType == "variable" {
			if len(node.Children) > 0 && node.Children[0].Type == parser.IdentifierNode {
				var name string
				if val, ok := node.Children[0].Value.(string); ok && val != "" {
					name = val
				} else if node.Children[0].Token != "" {
					name = node.Children[0].Token
				}
				if name == oldName {
					*occurrences++
					*locations = append(*locations, RenameLocation{
						Line:    node.Span.Start.Line,
						Column:  node.Span.Start.Column,
						OldText: oldName,
						NewText: newName,
						Context: "DECLARE",
					})
				}
			}
		}

		// CALL statement (untuk procedure rename)
		if node.Type == parser.CallNode && targetType == "procedure" {
			// Cek child pertama (nama procedure)
			if len(node.Children) > 0 && node.Children[0].Type == parser.IdentifierNode {
				var name string
				if val, ok := node.Children[0].Value.(string); ok && val != "" {
					name = val
				} else if node.Children[0].Token != "" {
					name = node.Children[0].Token
				}
				if name == oldName {
					*occurrences++
					*locations = append(*locations, RenameLocation{
						Line:    node.Span.Start.Line,
						Column:  node.Span.Start.Column,
						OldText: oldName,
						NewText: newName,
						Context: "CALL",
					})
				}
			}
		}

		// Function call (untuk function rename)
		if node.Type == parser.FunctionCallNode && targetType == "function" {
			var name string
			if val, ok := node.Value.(string); ok && val != "" {
				name = val
			}
			if name == "" && len(node.Children) > 0 {
				if node.Children[0].Type == parser.IdentifierNode {
					if val, ok := node.Children[0].Value.(string); ok && val != "" {
						name = val
					} else if node.Children[0].Token != "" {
						name = node.Children[0].Token
					}
				}
			}
			if name == oldName {
				*occurrences++
				*locations = append(*locations, RenameLocation{
					Line:    node.Span.Start.Line,
					Column:  node.Span.Start.Column,
					OldText: oldName,
					NewText: newName,
					Context: "FUNCTION_CALL",
				})
			}
		}

		// Identifier usage (untuk variable rename)
		if node.Type == parser.IdentifierNode && targetType == "variable" {
			if val, ok := node.Value.(string); ok && val == oldName {
				// Skip jika ini adalah deklarasi (sudah di-handle di atas)
				// Kita akan handle di level atas
			}
		}

		// Recurse into children
		for _, child := range node.Children {
			searchAndReplace([]parser.ASTNode{child}, oldName, newName, locations, occurrences, targetType)
		}
	}
}

func getContext(node parser.ASTNode) string {
	// Try to determine context from parent
	// For now, return simple context
	return "USAGE"
}

// ============================================
// FORMAT OUTPUT
// ============================================

func FormatRenameResult(result RenameResult, dryRun bool) string {
	var sb strings.Builder

	if !result.Success {
		sb.WriteString(fmt.Sprintf("\n❌ %s\n", result.Message))
		return sb.String()
	}

	if dryRun {
		sb.WriteString("\n🔍 Dry Run - Preview changes:\n")
		sb.WriteString(strings.Repeat("-", 40) + "\n\n")
	} else {
		sb.WriteString("\n✅ " + result.Message + "\n")
		sb.WriteString(strings.Repeat("-", 40) + "\n\n")
	}

	sb.WriteString(fmt.Sprintf("📝 Changes made (%d occurrences):\n", result.Occurrences))
	for _, loc := range result.Locations {
		sb.WriteString(fmt.Sprintf("  Line %d: %s → %s (%s)\n",
			loc.Line, loc.OldText, loc.NewText, loc.Context))
	}

	if dryRun {
		sb.WriteString(fmt.Sprintf("\n📊 %d changes will be made\n", result.Occurrences))
		sb.WriteString("❌ No files were modified (dry-run mode)\n")
	} else {
		sb.WriteString(fmt.Sprintf("\n📊 %d changes applied\n", result.Occurrences))
	}

	return sb.String()
}

// ============================================
// EXISTING HELPERS
// ============================================

func isUsed(usedVars []string, name string) bool {
	for _, v := range usedVars {
		if v == name {
			return true
		}
	}
	return false
}

func getLines(usages []analyzer.UsageInfo) []int {
	var lines []int
	for _, u := range usages {
		lines = append(lines, u.Line)
	}
	return lines
}

func FormatSuggestions(result RefactorResult) string {
	var sb strings.Builder

	if len(result.Suggestions) == 0 {
		sb.WriteString("✅ No refactoring suggestions found. Code looks clean!\n")
		return sb.String()
	}

	sb.WriteString("\n📊 Refactoring Suggestions\n")
	sb.WriteString(strings.Repeat("=", 50) + "\n\n")

	var deadCode, codeSmells, improvements []Suggestion
	for _, s := range result.Suggestions {
		switch s.Type {
		case "dead_code":
			deadCode = append(deadCode, s)
		case "code_smell":
			codeSmells = append(codeSmells, s)
		case "improvement":
			improvements = append(improvements, s)
		}
	}

	if len(deadCode) > 0 {
		sb.WriteString("🔴 Dead Code Detected:\n")
		for _, s := range deadCode {
			sb.WriteString(fmt.Sprintf("  - %s (line %d)\n", s.Message, s.Line))
			for _, d := range s.Details {
				sb.WriteString(fmt.Sprintf("    → %s\n", d))
			}
		}
		sb.WriteString("\n")
	}

	if len(codeSmells) > 0 {
		sb.WriteString("🟡 Code Smells:\n")
		for _, s := range codeSmells {
			sb.WriteString(fmt.Sprintf("  - %s (line %d)\n", s.Message, s.Line))
			for _, d := range s.Details {
				sb.WriteString(fmt.Sprintf("    → %s\n", d))
			}
		}
		sb.WriteString("\n")
	}

	if len(improvements) > 0 {
		sb.WriteString("🟢 Improvements:\n")
		for _, s := range improvements {
			sb.WriteString(fmt.Sprintf("  - %s (line %d)\n", s.Message, s.Line))
			for _, d := range s.Details {
				sb.WriteString(fmt.Sprintf("    → %s\n", d))
			}
		}
		sb.WriteString("\n")
	}

	sb.WriteString("📈 Statistics:\n")
	for key, value := range result.Stats {
		sb.WriteString(fmt.Sprintf("  - %s: %d\n", key, value))
	}
	sb.WriteString(fmt.Sprintf("  - Total suggestions: %d\n", len(result.Suggestions)))

	return sb.String()
}

func FormatDeadCode(result RefactorResult) string {
	var sb strings.Builder

	var deadCode []Suggestion
	for _, s := range result.Suggestions {
		if s.Type == "dead_code" {
			deadCode = append(deadCode, s)
		}
	}

	if len(deadCode) == 0 {
		sb.WriteString("✅ No dead code detected. Code looks clean!\n")
		return sb.String()
	}

	sb.WriteString("\n🗑️ Dead Code Analysis\n")
	sb.WriteString(strings.Repeat("=", 40) + "\n\n")

	for _, s := range deadCode {
		emoji := "🔴"
		if s.Severity == "medium" {
			emoji = "🟡"
		} else if s.Severity == "low" {
			emoji = "🟢"
		}
		sb.WriteString(fmt.Sprintf("%s %s (line %d)\n", emoji, s.Message, s.Line))
		for _, d := range s.Details {
			sb.WriteString(fmt.Sprintf("    → %s\n", d))
		}
		sb.WriteString("\n")
	}

	sb.WriteString(fmt.Sprintf("📊 Total dead code: %d items\n", len(deadCode)))

	return sb.String()
}

func ApplyRenameChanges(originalContent string, result RenameResult) string {
	lines := strings.Split(originalContent, "\n")

	// Create a map of line -> replacements
	replacements := make(map[int][]string)
	for _, loc := range result.Locations {
		// Replace old name with new name on that line
		// Note: This is simplified; for production, use more precise replacement
		oldLine := lines[loc.Line-1]
		newLine := strings.ReplaceAll(oldLine, loc.OldText, loc.NewText)
		if oldLine != newLine {
			replacements[loc.Line-1] = append(replacements[loc.Line-1], newLine)
		}
	}

	// Apply replacements
	for lineNum, newLines := range replacements {
		if len(newLines) > 0 {
			// Use the last replacement (most complete)
			lines[lineNum] = newLines[len(newLines)-1]
		}
	}

	return strings.Join(lines, "\n")
}

// ============================================
// Helper Functions
// ============================================

// appendUnique appends item to slice if not already present.
func appendUnique(slice []string, item string) []string {
	for _, s := range slice {
		if s == item {
			return slice
		}
	}
	return append(slice, item)
}

// ============================================
// EXPLAIN - Natural Language Explanation
// ============================================

type ExplanationResult struct {
	ModuleName  string
	Variables   []VariableInfo
	Procedures  []ProcedureInfo
	Functions   []FunctionInfo
	CallFlow    []string
	Summary     string
	Warnings    []string
	Suggestions []string
}

type VariableInfo struct {
	Name string
	Type string
	Line int
}

type ProcedureInfo struct {
	Name     string
	Line     int
	Calls    []string
	IsCalled bool
}

type FunctionInfo struct {
	Name       string
	Line       int
	ReturnType string
	IsCalled   bool
}

func Explain(program parser.Program, analysisResult analyzer.AnalysisResult) ExplanationResult {
	var result ExplanationResult

	// 1. Extract module name
	for _, stmt := range program.Statements {
		// Cek ModuleNode langsung
		if stmt.Type == parser.ModuleNode {
			if len(stmt.Children) > 0 && stmt.Children[0].Type == parser.IdentifierNode {
				if name, ok := stmt.Children[0].Value.(string); ok && name != "" {
					result.ModuleName = name
				} else if stmt.Children[0].Token != "" {
					result.ModuleName = stmt.Children[0].Token
				}
			}
		}
		// Cek CreateNode -> ModuleNode
		if stmt.Type == parser.CreateNode {
			for _, child := range stmt.Children {
				if child.Type == parser.ModuleNode {
					if len(child.Children) > 0 && child.Children[0].Type == parser.IdentifierNode {
						if name, ok := child.Children[0].Value.(string); ok && name != "" {
							result.ModuleName = name
						} else if child.Children[0].Token != "" {
							result.ModuleName = child.Children[0].Token
						}
					}
				}
			}
		}
	}
	if result.ModuleName == "" {
		result.ModuleName = "Unnamed"
	}

	// 2. Variables
	for name, info := range analysisResult.Variables {
		result.Variables = append(result.Variables, VariableInfo{
			Name: name,
			Type: info.Type,
			Line: info.Line,
		})
	}
	sort.Slice(result.Variables, func(i, j int) bool {
		return result.Variables[i].Name < result.Variables[j].Name
	})

	// ============================================
	// 3. Manual scan for ALL CALLs and FunctionCalls
	// ============================================
	callGraph, reverseCallGraph := BuildCallGraph(program)
	mergedCallGraph := callGraph
	mergedReverseCallGraph := reverseCallGraph

	var scanCalls func(node parser.ASTNode, inProcedure bool, currentProc string)
	scanCalls = func(node parser.ASTNode, inProcedure bool, currentProc string) {
		// Handle CallNode
		if node.Type == parser.CallNode {
			var callee string
			if len(node.Children) > 0 && node.Children[0].Type == parser.IdentifierNode {
				if v, ok := node.Children[0].Value.(string); ok {
					callee = v
				} else if node.Children[0].Token != "" {
					callee = node.Children[0].Token
				}
			}
			if callee != "" {
				caller := "MAIN"
				if inProcedure && currentProc != "" {
					caller = currentProc
				}
				callGraph[caller] = appendUnique(callGraph[caller], callee)
				reverseCallGraph[callee] = appendUnique(reverseCallGraph[callee], caller)
			}
		}

		// Handle FunctionCallNode (e.g., FuncA())
		if node.Type == parser.FunctionCallNode {
			var callee string
			if v, ok := node.Value.(string); ok {
				callee = v
			}
			if callee != "" {
				caller := "MAIN"
				if inProcedure && currentProc != "" {
					caller = currentProc
				}
				callGraph[caller] = appendUnique(callGraph[caller], callee)
				reverseCallGraph[callee] = appendUnique(reverseCallGraph[callee], caller)
			}
		}

		// Handle ProcedureNode: enter procedure scope
		if node.Type == parser.ProcedureNode {
			if len(node.Children) > 0 && node.Children[0].Type == parser.IdentifierNode {
				if name, ok := node.Children[0].Value.(string); ok {
					for _, child := range node.Children {
						scanCalls(child, true, name)
					}
					return
				}
			}
		}

		// Recurse into children
		for _, child := range node.Children {
			scanCalls(child, inProcedure, currentProc)
		}
	}

	for _, stmt := range program.Statements {
		scanCalls(stmt, false, "")
	}

	// 4. Procedures
	for name, info := range analysisResult.Procedures {
		proc := ProcedureInfo{
			Name:     name,
			Line:     info.Line,
			IsCalled: false,
			Calls:    []string{},
		}
		if callers, ok := mergedReverseCallGraph[name]; ok && len(callers) > 0 {
			for _, caller := range callers {
				if caller == "MAIN" || caller != "" {
					proc.IsCalled = true
					break
				}
			}
		}
		if callees, ok := mergedCallGraph[name]; ok {
			proc.Calls = callees
		}
		result.Procedures = append(result.Procedures, proc)
	}
	sort.Slice(result.Procedures, func(i, j int) bool {
		return result.Procedures[i].Name < result.Procedures[j].Name
	})

	// 5. Functions
	for name, info := range analysisResult.Functions {
		returnType := info.ReturnType

		// Jika returnType masih UNKNOWN atau kosong, cari dari AST
		if returnType == "" || returnType == "UNKNOWN" {
			var search func(node parser.ASTNode)
			search = func(node parser.ASTNode) {
				if node.Type == parser.FunctionNode {
					if len(node.Children) > 0 && node.Children[0].Type == parser.IdentifierNode {
						var funcName string
						if val, ok := node.Children[0].Value.(string); ok && val != "" {
							funcName = val
						} else if node.Children[0].Token != "" {
							funcName = node.Children[0].Token
						}
						if funcName == name {
							for _, child := range node.Children {
								if child.Type == parser.ReturnTypeNode {
									if val, ok := child.Value.(string); ok && val != "" {
										returnType = val
									} else if child.Token != "" {
										returnType = child.Token
									}
									// Coba dari children jika masih kosong
									if returnType == "" && len(child.Children) > 0 {
										if val, ok := child.Children[0].Value.(string); ok && val != "" {
											returnType = val
										} else if child.Children[0].Token != "" {
											returnType = child.Children[0].Token
										}
									}
								}
							}
						}
					}
				}
				for _, child := range node.Children {
					search(child)
				}
			}
			for _, stmt := range program.Statements {
				search(stmt)
			}
		}

		// ✅ PAKAI returnType yang sudah dicari, BUKAN info.ReturnType
		funcInfo := FunctionInfo{
			Name:       name,
			Line:       info.Line,
			ReturnType: returnType, // ← Pakai returnType, bukan info.ReturnType
			IsCalled:   false,
		}
		if callers, ok := mergedReverseCallGraph[name]; ok && len(callers) > 0 {
			for _, caller := range callers {
				if caller == "MAIN" || caller != "" {
					funcInfo.IsCalled = true
					break
				}
			}
		}
		result.Functions = append(result.Functions, funcInfo)
	}
	sort.Slice(result.Functions, func(i, j int) bool {
		return result.Functions[i].Name < result.Functions[j].Name
	})

	// 6. Call Flow
	if len(mergedCallGraph) > 0 {
		var callers []string
		for caller := range mergedCallGraph {
			callers = append(callers, caller)
		}
		sort.Strings(callers)
		for _, caller := range callers {
			callees := mergedCallGraph[caller]
			sort.Strings(callees)
			for _, callee := range callees {
				if caller == "MAIN" {
					result.CallFlow = append(result.CallFlow, fmt.Sprintf("(main) → %s", callee))
				} else {
					result.CallFlow = append(result.CallFlow, fmt.Sprintf("%s → %s", caller, callee))
				}
			}
		}
	}

	// 7. Summary
	var parts []string

	if len(result.Variables) > 0 {
		word := "variables"
		if len(result.Variables) == 1 {
			word = "variable"
		}
		parts = append(parts, fmt.Sprintf("%d %s", len(result.Variables), word))
	}

	if len(result.Procedures) > 0 {
		word := "procedures"
		if len(result.Procedures) == 1 {
			word = "procedure"
		}
		parts = append(parts, fmt.Sprintf("%d %s", len(result.Procedures), word))
	}

	if len(result.Functions) > 0 {
		word := "functions"
		if len(result.Functions) == 1 {
			word = "function"
		}
		parts = append(parts, fmt.Sprintf("%d %s", len(result.Functions), word))
	}

	if len(parts) == 0 {
		result.Summary = fmt.Sprintf("Module '%s' is empty.", result.ModuleName)
	} else if len(parts) == 1 {
		result.Summary = fmt.Sprintf("Module '%s' contains %s.", result.ModuleName, parts[0])
	} else if len(parts) == 2 {
		result.Summary = fmt.Sprintf("Module '%s' contains %s and %s.", result.ModuleName, parts[0], parts[1])
	} else {
		result.Summary = fmt.Sprintf("Module '%s' contains %s, %s, and %s.", result.ModuleName, parts[0], parts[1], parts[2])
	}

	// 8. Warnings
	for _, proc := range result.Procedures {
		if !proc.IsCalled {
			result.Warnings = append(result.Warnings, fmt.Sprintf("Procedure '%s' is never called (line %d)", proc.Name, proc.Line))
		}
	}
	for _, fn := range result.Functions {
		if !fn.IsCalled && fn.ReturnType != "BUILTIN" {
			result.Warnings = append(result.Warnings, fmt.Sprintf("Function '%s' is never called (line %d)", fn.Name, fn.Line))
		}
	}
	sort.Strings(result.Warnings)

	// 9. Suggestions
	processed := make(map[string]bool)
	for _, proc := range result.Procedures {
		if processed[proc.Name] {
			continue
		}
		processed[proc.Name] = true

		if !proc.IsCalled {
			result.Suggestions = append(result.Suggestions, fmt.Sprintf("Consider removing or using procedure '%s'", proc.Name))
		}
		if calls, ok := mergedCallGraph[proc.Name]; ok && len(calls) == 1 {
			result.Suggestions = append(result.Suggestions,
				fmt.Sprintf("Procedure '%s' only calls one other procedure (%s), consider inlining",
					proc.Name, calls[0]))
		}
	}
	sort.Strings(result.Suggestions)

	return result
}

// FormatExplanation returns a human-readable string from ExplanationResult.
func FormatExplanation(result ExplanationResult) string {
	var sb strings.Builder

	sb.WriteString("\n📖 Code Explanation\n")
	sb.WriteString(strings.Repeat("=", 50) + "\n\n")

	sb.WriteString(fmt.Sprintf("📦 Module: %s\n\n", result.ModuleName))

	if len(result.Variables) > 0 {
		sb.WriteString("📊 Variables:\n")
		for _, v := range result.Variables {
			sb.WriteString(fmt.Sprintf("  - %s: %s (line %d)\n", v.Name, v.Type, v.Line))
		}
		sb.WriteString("\n")
	}

	if len(result.Procedures) > 0 {
		sb.WriteString("🔧 Procedures:\n")
		for _, p := range result.Procedures {
			called := "❌ unused"
			if p.IsCalled {
				called = "✅ used"
			}
			sb.WriteString(fmt.Sprintf("  - %s (line %d) [%s]\n", p.Name, p.Line, called))
			if len(p.Calls) > 0 {
				sb.WriteString(fmt.Sprintf("    → Calls: %v\n", p.Calls))
			}
		}
		sb.WriteString("\n")
	}

	if len(result.Functions) > 0 {
		sb.WriteString("⚡ Functions:\n")
		for _, f := range result.Functions {
			called := "❌ unused"
			if f.IsCalled {
				called = "✅ used"
			}
			sb.WriteString(fmt.Sprintf("  - %s: %s (line %d) [%s]\n", f.Name, f.ReturnType, f.Line, called))
		}
		sb.WriteString("\n")
	}

	if len(result.CallFlow) > 0 {
		sb.WriteString("🔄 Call Flow:\n")
		for _, flow := range result.CallFlow {
			sb.WriteString(fmt.Sprintf("  %s\n", flow))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("📝 Summary:\n")
	sb.WriteString(fmt.Sprintf("  %s\n\n", result.Summary))

	if len(result.Warnings) > 0 {
		sb.WriteString("⚠️ Warnings:\n")
		for _, w := range result.Warnings {
			sb.WriteString(fmt.Sprintf("  - %s\n", w))
		}
		sb.WriteString("\n")
	}

	if len(result.Suggestions) > 0 {
		sb.WriteString("💡 Suggestions:\n")
		for _, s := range result.Suggestions {
			sb.WriteString(fmt.Sprintf("  - %s\n", s))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

// ============================================
// SEARCH FUNCTIONS
// ============================================

type SearchResult struct {
	Query      string        `json:"query"`
	Type       string        `json:"type"`
	Matches    []SearchMatch `json:"matches"`
	TotalCount int           `json:"totalCount"`
	Message    string        `json:"message"`
}

type SearchMatch struct {
	Name     string `json:"name"`
	Line     int    `json:"line"`
	Column   int    `json:"column"`
	Context  string `json:"context"`
	FullText string `json:"fullText"`
}

// SearchProcedure finds all occurrences of a procedure (definition and calls)
func SearchProcedure(program parser.Program, name string) SearchResult {
	var matches []SearchMatch

	var search func(node parser.ASTNode)
	search = func(node parser.ASTNode) {
		// Check procedure definition
		if node.Type == parser.ProcedureNode {
			if len(node.Children) > 0 && node.Children[0].Type == parser.IdentifierNode {
				if val, ok := node.Children[0].Value.(string); ok && val == name {
					matches = append(matches, SearchMatch{
						Name:     val,
						Line:     node.Span.Start.Line,
						Column:   node.Span.Start.Column,
						Context:  "PROCEDURE",
						FullText: fmt.Sprintf("CREATE PROCEDURE %s()", val),
					})
				}
			}
		}

		// ✅ Check CALL statements
		if node.Type == parser.CallNode {
			var calleeName string

			// Coba dari children pertama (nama procedure)
			if len(node.Children) > 0 && node.Children[0].Type == parser.IdentifierNode {
				// Coba dari Value
				if val, ok := node.Children[0].Value.(string); ok {
					calleeName = val
				}
				// Coba dari Token (karena kadang Value kosong)
				if calleeName == "" && node.Children[0].Token != "" {
					calleeName = node.Children[0].Token
				}
			}

			if calleeName == name {
				matches = append(matches, SearchMatch{
					Name:     calleeName,
					Line:     node.Span.Start.Line,
					Column:   node.Span.Start.Column,
					Context:  "CALL",
					FullText: fmt.Sprintf("CALL %s()", calleeName),
				})
			}
		}

		for _, child := range node.Children {
			search(child)
		}
	}

	for _, stmt := range program.Statements {
		search(stmt)
	}

	// Deduplicate
	unique := make(map[string]SearchMatch)
	for _, m := range matches {
		key := fmt.Sprintf("%d:%d", m.Line, m.Column)
		unique[key] = m
	}
	var deduped []SearchMatch
	for _, m := range unique {
		deduped = append(deduped, m)
	}
	sort.Slice(deduped, func(i, j int) bool {
		return deduped[i].Line < deduped[j].Line
	})

	if len(deduped) == 0 {
		return SearchResult{
			Query:      name,
			Type:       "procedure",
			Matches:    deduped,
			TotalCount: 0,
			Message:    fmt.Sprintf("Procedure '%s' not found", name),
		}
	}

	return SearchResult{
		Query:      name,
		Type:       "procedure",
		Matches:    deduped,
		TotalCount: len(deduped),
		Message:    fmt.Sprintf("Found %d occurrence(s) of procedure '%s'", len(deduped), name),
	}
}

// SearchFunction finds all occurrences of a function (definition and calls)
func SearchFunction(program parser.Program, name string) SearchResult {
	var matches []SearchMatch

	var search func(node parser.ASTNode)
	search = func(node parser.ASTNode) {
		// Check function definition
		if node.Type == parser.FunctionNode {
			if len(node.Children) > 0 && node.Children[0].Type == parser.IdentifierNode {
				if val, ok := node.Children[0].Value.(string); ok && val == name {
					returnType := "UNKNOWN"
					if len(node.Children) > 1 && node.Children[1].Type == parser.IdentifierNode {
						if v, ok := node.Children[1].Value.(string); ok {
							returnType = v
						}
					}
					matches = append(matches, SearchMatch{
						Name:     val,
						Line:     node.Span.Start.Line,
						Column:   node.Span.Start.Column,
						Context:  "FUNCTION",
						FullText: fmt.Sprintf("CREATE FUNCTION %s() RETURNS %s", val, returnType),
					})
				}
			}
		}

		// Check function calls
		if node.Type == parser.FunctionCallNode {
			if val, ok := node.Value.(string); ok && val == name {
				matches = append(matches, SearchMatch{
					Name:     val,
					Line:     node.Span.Start.Line,
					Column:   node.Span.Start.Column,
					Context:  "FUNCTION_CALL",
					FullText: fmt.Sprintf("%s()", val),
				})
			}
		}

		for _, child := range node.Children {
			search(child)
		}
	}

	for _, stmt := range program.Statements {
		search(stmt)
	}

	if len(matches) == 0 {
		return SearchResult{
			Query:      name,
			Type:       "function",
			Matches:    matches,
			TotalCount: 0,
			Message:    fmt.Sprintf("Function '%s' not found", name),
		}
	}

	return SearchResult{
		Query:      name,
		Type:       "function",
		Matches:    matches,
		TotalCount: len(matches),
		Message:    fmt.Sprintf("Found %d occurrence(s) of function '%s'", len(matches), name),
	}
}

// SearchVariable finds all occurrences of a variable (declaration and usage)
func SearchVariable(program parser.Program, name string) SearchResult {
	var matches []SearchMatch
	var declareLines map[int]bool = make(map[int]bool) // Track lines yang sudah di-declare

	var search func(node parser.ASTNode)
	search = func(node parser.ASTNode) {
		// Check variable declaration
		if node.Type == parser.DeclareNode {
			if len(node.Children) > 0 && node.Children[0].Type == parser.IdentifierNode {
				if val, ok := node.Children[0].Value.(string); ok && val == name {
					varType := "UNKNOWN"
					if len(node.Children) > 1 && node.Children[1].Type == parser.IdentifierNode {
						if v, ok := node.Children[1].Value.(string); ok {
							varType = v
						}
					}
					matches = append(matches, SearchMatch{
						Name:     val,
						Line:     node.Span.Start.Line,
						Column:   node.Span.Start.Column,
						Context:  "DECLARE",
						FullText: fmt.Sprintf("DECLARE %s %s", val, varType),
					})
					declareLines[node.Span.Start.Line] = true // Mark as declared
				}
			}
		}

		// Check identifier usage - skip if it's the declaration itself
		if node.Type == parser.IdentifierNode {
			if val, ok := node.Value.(string); ok && val == name {
				// Skip if this is the declaration line (already handled)
				if !declareLines[node.Span.Start.Line] {
					matches = append(matches, SearchMatch{
						Name:     val,
						Line:     node.Span.Start.Line,
						Column:   node.Span.Start.Column,
						Context:  "USAGE",
						FullText: val,
					})
				}
			}
		}

		for _, child := range node.Children {
			search(child)
		}
	}

	for _, stmt := range program.Statements {
		search(stmt)
	}

	// Remove duplicates (same line and column)
	unique := make(map[string]SearchMatch)
	for _, m := range matches {
		key := fmt.Sprintf("%d:%d", m.Line, m.Column)
		unique[key] = m
	}
	var deduped []SearchMatch
	for _, m := range unique {
		deduped = append(deduped, m)
	}
	// Sort by line
	sort.Slice(deduped, func(i, j int) bool {
		return deduped[i].Line < deduped[j].Line
	})

	if len(deduped) == 0 {
		return SearchResult{
			Query:      name,
			Type:       "variable",
			Matches:    deduped,
			TotalCount: 0,
			Message:    fmt.Sprintf("Variable '%s' not found", name),
		}
	}

	return SearchResult{
		Query:      name,
		Type:       "variable",
		Matches:    deduped,
		TotalCount: len(deduped),
		Message:    fmt.Sprintf("Found %d occurrence(s) of variable '%s'", len(deduped), name),
	}
}

// SearchCall finds all CALL statements to a specific procedure
func SearchCall(program parser.Program, name string) SearchResult {
	var matches []SearchMatch

	var search func(node parser.ASTNode)
	search = func(node parser.ASTNode) {
		if node.Type == parser.CallNode {
			var callee string

			// Coba dari children pertama
			if len(node.Children) > 0 && node.Children[0].Type == parser.IdentifierNode {
				if val, ok := node.Children[0].Value.(string); ok {
					callee = val
				}
				if callee == "" && node.Children[0].Token != "" {
					callee = node.Children[0].Token
				}
			}

			if callee == name {
				matches = append(matches, SearchMatch{
					Name:     callee,
					Line:     node.Span.Start.Line,
					Column:   node.Span.Start.Column,
					Context:  "CALL",
					FullText: fmt.Sprintf("CALL %s()", callee),
				})
			}
		}
		for _, child := range node.Children {
			search(child)
		}
	}

	for _, stmt := range program.Statements {
		search(stmt)
	}

	if len(matches) == 0 {
		return SearchResult{
			Query:      name,
			Type:       "call",
			Matches:    matches,
			TotalCount: 0,
			Message:    fmt.Sprintf("CALL to '%s' not found", name),
		}
	}

	return SearchResult{
		Query:      name,
		Type:       "call",
		Matches:    matches,
		TotalCount: len(matches),
		Message:    fmt.Sprintf("Found %d CALL(s) to '%s'", len(matches), name),
	}
}

// SearchUnused finds all unused code (procedures, functions, variables)
func SearchUnused(program parser.Program, analysisResult analyzer.AnalysisResult) SearchResult {
	var matches []SearchMatch

	// ✅ Hapus baris ini karena tidak terpakai
	// reverseCallGraph := analysisResult.ReverseCallGraph

	callGraph := analysisResult.CallGraph

	// Build set of all procedures that are called (directly or indirectly) from MAIN
	calledFromMain := make(map[string]bool)

	// Start with procedures called directly from MAIN
	var queue []string
	if mainCalls, ok := callGraph["MAIN"]; ok {
		for _, callee := range mainCalls {
			if !calledFromMain[callee] {
				calledFromMain[callee] = true
				queue = append(queue, callee)
			}
		}
	}

	// BFS: mark all procedures called by procedures that are called from MAIN
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		if callees, ok := callGraph[current]; ok {
			for _, callee := range callees {
				if !calledFromMain[callee] {
					calledFromMain[callee] = true
					queue = append(queue, callee)
				}
			}
		}
	}

	// Check unused procedures
	for name, info := range analysisResult.Procedures {
		if !calledFromMain[name] {
			matches = append(matches, SearchMatch{
				Name:     name,
				Line:     info.Line,
				Column:   0,
				Context:  "PROCEDURE (unused)",
				FullText: fmt.Sprintf("CREATE PROCEDURE %s()", name),
			})
		}
	}

	// Check unused functions
	for name, info := range analysisResult.Functions {
		if !calledFromMain[name] && info.ReturnType != "BUILTIN" {
			matches = append(matches, SearchMatch{
				Name:     name,
				Line:     info.Line,
				Column:   0,
				Context:  "FUNCTION (unused)",
				FullText: fmt.Sprintf("CREATE FUNCTION %s()", name),
			})
		}
	}

	// Check unused variables
	usedVars := make(map[string]bool)
	for _, v := range analysisResult.UsedVariables {
		usedVars[v] = true
	}
	for name, info := range analysisResult.Variables {
		if !usedVars[name] {
			matches = append(matches, SearchMatch{
				Name:     name,
				Line:     info.Line,
				Column:   0,
				Context:  "VARIABLE (unused)",
				FullText: fmt.Sprintf("DECLARE %s %s", name, info.Type),
			})
		}
	}

	// Sort by line
	sort.Slice(matches, func(i, j int) bool {
		return matches[i].Line < matches[j].Line
	})

	if len(matches) == 0 {
		return SearchResult{
			Type:       "unused",
			Matches:    matches,
			TotalCount: 0,
			Message:    "No unused code found",
		}
	}

	return SearchResult{
		Type:       "unused",
		Matches:    matches,
		TotalCount: len(matches),
		Message:    fmt.Sprintf("Found %d unused item(s)", len(matches)),
	}
}

// FormatSearchResult formats a SearchResult for human-readable output
func FormatSearchResult(result SearchResult) string {
	var sb strings.Builder

	if result.TotalCount == 0 {
		sb.WriteString(fmt.Sprintf("\n🔍 Search Result: %s\n", result.Message))
		return sb.String()
	}

	sb.WriteString(fmt.Sprintf("\n🔍 Search Result: %s\n", result.Message))
	sb.WriteString(strings.Repeat("=", 50) + "\n\n")

	for _, match := range result.Matches {
		sb.WriteString(fmt.Sprintf("  Line %d: %s (%s)\n", match.Line, match.FullText, match.Context))
	}

	sb.WriteString(fmt.Sprintf("\n📊 Total: %d match(es)\n", result.TotalCount))
	return sb.String()
}

func BuildCallGraph(program parser.Program) (map[string][]string, map[string][]string) {
	callGraph := make(map[string][]string)
	reverseCallGraph := make(map[string][]string)

	var scan func(node parser.ASTNode, inProcedure bool, currentProc string)
	scan = func(node parser.ASTNode, inProcedure bool, currentProc string) {
		// Handle CallNode
		if node.Type == parser.CallNode {
			var callee string
			if len(node.Children) > 0 && node.Children[0].Type == parser.IdentifierNode {
				if val, ok := node.Children[0].Value.(string); ok && val != "" {
					callee = val
				} else if node.Children[0].Token != "" {
					callee = node.Children[0].Token
				}
			}
			if callee == "" {
				if val, ok := node.Value.(string); ok && val != "" {
					callee = val
				}
			}
			if callee != "" {
				caller := "MAIN"
				if inProcedure && currentProc != "" {
					caller = currentProc
				}
				callGraph[caller] = appendUnique(callGraph[caller], callee)
				reverseCallGraph[callee] = appendUnique(reverseCallGraph[callee], caller)
			}
		}

		// Handle FunctionCallNode
		if node.Type == parser.FunctionCallNode {
			var callee string
			if val, ok := node.Value.(string); ok && val != "" {
				callee = val
			}
			if callee == "" && len(node.Children) > 0 {
				if node.Children[0].Type == parser.IdentifierNode {
					if val, ok := node.Children[0].Value.(string); ok && val != "" {
						callee = val
					} else if node.Children[0].Token != "" {
						callee = node.Children[0].Token
					}
				}
			}
			if callee != "" {
				caller := "MAIN"
				if inProcedure && currentProc != "" {
					caller = currentProc
				}
				callGraph[caller] = appendUnique(callGraph[caller], callee)
				reverseCallGraph[callee] = appendUnique(reverseCallGraph[callee], caller)
			}
		}

		// Track procedure entry - skip identifier (nama procedure)
		if node.Type == parser.ProcedureNode {
			if len(node.Children) > 0 && node.Children[0].Type == parser.IdentifierNode {
				name := ""
				if val, ok := node.Children[0].Value.(string); ok && val != "" {
					name = val
				} else if node.Children[0].Token != "" {
					name = node.Children[0].Token
				}
				if name != "" {
					// ✅ SKIP child pertama (IdentifierNode - nama procedure)
					// Scan dari child index 1 (body statements)
					for i := 1; i < len(node.Children); i++ {
						scan(node.Children[i], true, name)
					}
					return
				}
			}
		}

		// Track function entry - skip identifier (nama function)
		if node.Type == parser.FunctionNode {
			if len(node.Children) > 0 && node.Children[0].Type == parser.IdentifierNode {
				name := ""
				if val, ok := node.Children[0].Value.(string); ok && val != "" {
					name = val
				} else if node.Children[0].Token != "" {
					name = node.Children[0].Token
				}
				if name != "" {
					// ✅ SKIP child pertama (IdentifierNode - nama function)
					for i := 1; i < len(node.Children); i++ {
						scan(node.Children[i], true, name)
					}
					return
				}
			}
		}

		// Recurse into children
		for _, child := range node.Children {
			scan(child, inProcedure, currentProc)
		}
	}

	for _, stmt := range program.Statements {
		scan(stmt, false, "")
	}

	return callGraph, reverseCallGraph
}

func getNodeName(node parser.ASTNode) string {
	if node.Type == parser.IdentifierNode {
		if val, ok := node.Value.(string); ok && val != "" {
			return val
		}
		return node.Token
	}
	if node.Type == parser.FunctionCallNode {
		if val, ok := node.Value.(string); ok && val != "" {
			return val
		}
		if len(node.Children) > 0 {
			return getNodeName(node.Children[0])
		}
	}
	return ""
}

```

---

### File: `pkg/generator/generator.go`

```go
package generator

import (
	"fmt"
	"strings"

	"esql-ast-tool/pkg/parser"
)

type Generator struct {
	indent string
}

func NewGenerator() *Generator {
	return &Generator{indent: "    "}
}

func (g *Generator) Generate(program parser.Program) string {
	var sb strings.Builder
	var stmts []parser.ASTNode
	if len(program.Statements) > 0 {
		stmts = program.Statements
	} else if program.Root.Type != "" {
		stmts = program.Root.Children
	} else {
		return ""
	}
	for i, stmt := range stmts {
		if i > 0 {
			sb.WriteString("\n")
		}
		sb.WriteString(g.generateNode(stmt, 0))
	}
	return sb.String()
}

func (g *Generator) generateNode(node parser.ASTNode, level int) string {
	if node.Type == "" {
		return ""
	}

	switch node.Type {
	case parser.ProgramNode:
		return g.generateChildren(node.Children, 0)

	case parser.ModuleNode:
		return g.generateModule(node, level)

	case parser.CreateNode:
		for _, child := range node.Children {
			if child.Type == parser.ModuleNode {
				return g.generateModule(child, level)
			}
		}
		return ""

	case parser.ProcedureNode:
		return g.generateProcedure(node, level)

	case parser.FunctionNode:
		return g.generateFunction(node, level)

	case parser.DeclareNode:
		return g.generateDeclare(node, level)

	case parser.SetNode:
		return g.generateSet(node, level)

	case parser.IfNode:
		return g.generateIf(node, level)

	case parser.WhileNode:
		return g.generateWhile(node, level)

	case parser.ForNode:
		return g.generateFor(node, level)

	case parser.ReturnNode:
		return g.generateReturn(node, level)

	case parser.ThrowNode:
		return g.generateThrow(node, level)

	case parser.CallNode:
		return g.generateCall(node, level)

	case parser.BlockNode:
		return g.generateBlock(node, level)

	case parser.IdentifierNode:
		return g.getIdentifierName(node)

	case parser.LiteralNode:
		return g.formatLiteral(node)

	case parser.FieldReferenceNode:
		return g.generateFieldReference(node)

	case parser.ComparisonNode, parser.BinaryOpNode:
		return g.generateBinaryOp(node)

	case parser.UnaryOpNode:
		return g.generateUnaryOp(node)

	case parser.CastNode:
		return g.generateCast(node)

	case parser.CaseNode:
		return g.generateCase(node)

	case parser.WhenNode:
		return g.generateWhen(node)

	case parser.IsNullNode:
		if len(node.Children) > 0 {
			return g.generateNode(node.Children[0], 0) + " IS NULL"
		}
		return "IS NULL"

	case parser.IsNotNullNode:
		if len(node.Children) > 0 {
			return g.generateNode(node.Children[0], 0) + " IS NOT NULL"
		}
		return "IS NOT NULL"

	case parser.BetweenNode:
		return g.generateBetween(node)

	case parser.LikeNode:
		return g.generateLike(node)

	case parser.InNode:
		return g.generateIn(node)

	case parser.CoalesceNode:
		return g.generateCoalesce(node)

	case parser.NullIfNode:
		return g.generateNullIf(node)

	case parser.ParameterNode:
		return g.generateParameter(node)

	case parser.ReturnTypeNode:
		return g.getIdentifierName(node)

	case parser.ParenthesizedNode:
		if len(node.Children) > 0 {
			return "(" + g.generateNode(node.Children[0], 0) + ")"
		}
		return "()"
	case parser.PropagateNode:
		return g.generatePropagate(node, level)

	case parser.MoveNode:
		return g.generateMove(node, level)
	case parser.ContinueNode:
		return g.generateContinue(node, level)
	case parser.BreakNode:
		return g.generateBreak(node, level)
	case parser.LabelNode:
		return g.generateLabel(node, level)
	case parser.FunctionCallNode:
		return g.generateFunctionCall(node)

	default:
		return g.generateChildren(node.Children, level)
	}
}

// ============================================
// HELPER: Get Node Value (tanpa quote untuk mode)
// ============================================

func (g *Generator) getNodeValue(node parser.ASTNode) string {
	switch node.Type {
	case parser.LiteralNode:
		if val, ok := node.Value.(string); ok {
			// Mode parameter tidak boleh di-quote
			if val == "IN" || val == "OUT" || val == "INOUT" {
				return val
			}
			return val
		}
		if num, ok := node.Value.(float64); ok {
			if num == float64(int(num)) {
				return fmt.Sprintf("%d", int(num))
			}
			return fmt.Sprintf("%v", num)
		}
		return node.Token

	case parser.IdentifierNode:
		if val, ok := node.Value.(string); ok && val != "" {
			return val
		}
		return node.Token

	case parser.FieldReferenceNode:
		if val, ok := node.Value.(string); ok {
			return val
		}
		var parts []string
		for _, child := range node.Children {
			parts = append(parts, g.getNodeValue(child))
		}
		return strings.Join(parts, ".")

	default:
		if val, ok := node.Value.(string); ok && val != "" {
			return val
		}
		return node.Token
	}
}

// ============================================
// MODULE
// ============================================

func (g *Generator) generateModule(node parser.ASTNode, level int) string {
	var sb strings.Builder
	name := g.getModuleName(node)
	sb.WriteString("CREATE COMPUTE MODULE " + name + "\n")
	for _, child := range node.Children {
		if child.Type != parser.IdentifierNode {
			sb.WriteString(g.generateNode(child, level+1))
		}
	}
	sb.WriteString("END MODULE;\n")
	return sb.String()
}

func (g *Generator) getModuleName(node parser.ASTNode) string {
	for _, child := range node.Children {
		if child.Type == parser.IdentifierNode {
			return g.getIdentifierName(child)
		}
	}
	return "UnnamedModule"
}

// ============================================
// PROCEDURE
// ============================================

func (g *Generator) generateProcedure(node parser.ASTNode, level int) string {
	var sb strings.Builder
	indent := strings.Repeat(g.indent, level)
	name := g.getProcedureNameFromNode(node)
	params := g.getParameterList(node)
	sb.WriteString(indent + "CREATE PROCEDURE " + name + "(" + params + ")\n")
	sb.WriteString(indent + "BEGIN\n")
	for _, child := range node.Children {
		if child.Type != parser.IdentifierNode && child.Type != parser.ParameterNode {
			sb.WriteString(g.generateNode(child, level+1))
		}
	}
	sb.WriteString(indent + "END;\n")
	return sb.String()
}

func (g *Generator) getProcedureNameFromNode(node parser.ASTNode) string {
	for _, child := range node.Children {
		if child.Type == parser.IdentifierNode {
			return g.getIdentifierName(child)
		}
	}
	return "UnnamedProc"
}

func (g *Generator) getParameterList(node parser.ASTNode) string {
	var parts []string
	for _, child := range node.Children {
		if child.Type == parser.ParameterNode {
			parts = append(parts, g.generateParameter(child))
		}
	}
	return strings.Join(parts, ", ")
}

func (g *Generator) generateParameter(node parser.ASTNode) string {
	if len(node.Children) < 3 {
		return ""
	}
	// Gunakan getNodeValue, bukan generateNode (agar tidak di-quote)
	mode := g.getNodeValue(node.Children[0])
	name := g.getNodeValue(node.Children[1])
	typ := g.getNodeValue(node.Children[2])
	return mode + " " + name + " " + typ
}

// ============================================
// FUNCTION
// ============================================

func (g *Generator) generateFunction(node parser.ASTNode, level int) string {
	var sb strings.Builder
	indent := strings.Repeat(g.indent, level)
	name := g.getFunctionNameFromNode(node)
	params := g.getParameterList(node)
	returnType := g.getReturnType(node)
	sb.WriteString(indent + "CREATE FUNCTION " + name + "(" + params + ") RETURNS " + returnType + "\n")
	sb.WriteString(indent + "BEGIN\n")
	for _, child := range node.Children {
		if child.Type != parser.IdentifierNode && child.Type != parser.ParameterNode && child.Type != parser.ReturnTypeNode {
			sb.WriteString(g.generateNode(child, level+1))
		}
	}
	sb.WriteString(indent + "END;\n")
	return sb.String()
}

func (g *Generator) getFunctionNameFromNode(node parser.ASTNode) string {
	for _, child := range node.Children {
		if child.Type == parser.IdentifierNode {
			return g.getIdentifierName(child)
		}
	}
	return "UnnamedFunc"
}

func (g *Generator) getReturnType(node parser.ASTNode) string {
	for _, child := range node.Children {
		if child.Type == parser.ReturnTypeNode {
			// Ambil dari value
			if val, ok := child.Value.(string); ok && val != "" {
				return val
			}
			// Atau dari token
			if child.Token != "" {
				return child.Token
			}
			// Atau dari children
			if len(child.Children) > 0 {
				return g.getNodeValue(child.Children[0])
			}
		}
	}
	return "INTEGER" // Default
}

// ============================================
// DECLARE
// ============================================

func (g *Generator) generateDeclare(node parser.ASTNode, level int) string {
	var name, typ string
	for _, child := range node.Children {
		if child.Type == parser.IdentifierNode {
			if name == "" {
				name = g.getIdentifierName(child)
			} else {
				typ = g.getIdentifierName(child)
			}
		}
	}
	if name == "" || typ == "" {
		return ""
	}
	return strings.Repeat(g.indent, level) + "DECLARE " + name + " " + typ + ";\n"
}

// ============================================
// SET
// ============================================

func (g *Generator) generateSet(node parser.ASTNode, level int) string {
	var target, value string
	for _, child := range node.Children {
		if child.Type == parser.BlockNode {
			if child.Token == "target" && len(child.Children) > 0 {
				target = g.generateNode(child.Children[0], 0)
			} else if child.Token == "value" && len(child.Children) > 0 {
				// Cek apakah value adalah function call
				if child.Children[0].Type == parser.FunctionCallNode {
					value = g.generateFunctionCall(child.Children[0])
				} else {
					value = g.generateNode(child.Children[0], 0)
				}
			}
		}
	}
	if target == "" || value == "" {
		return ""
	}
	return strings.Repeat(g.indent, level) + "SET " + target + " = " + value + ";\n"
}

// ============================================
// IF
// ============================================

func (g *Generator) generateIf(node parser.ASTNode, level int) string {
	var sb strings.Builder
	indent := strings.Repeat(g.indent, level)
	if len(node.Children) == 0 {
		return ""
	}
	cond := g.generateNode(node.Children[0], 0)
	sb.WriteString(indent + "IF " + cond + " THEN\n")
	if len(node.Children) > 1 {
		thenBlock := node.Children[1]
		for _, stmt := range thenBlock.Children {
			sb.WriteString(g.generateNode(stmt, level+1))
		}
	}
	if len(node.Children) > 2 {
		elseBlock := node.Children[2]
		sb.WriteString(indent + "ELSE\n")
		for _, stmt := range elseBlock.Children {
			sb.WriteString(g.generateNode(stmt, level+1))
		}
	}
	sb.WriteString(indent + "END IF;\n")
	return sb.String()
}

// ============================================
// WHILE
// ============================================

func (g *Generator) generateWhile(node parser.ASTNode, level int) string {
	var sb strings.Builder
	indent := strings.Repeat(g.indent, level)
	if len(node.Children) == 0 {
		return ""
	}
	cond := g.generateNode(node.Children[0], 0)
	sb.WriteString(indent + "WHILE " + cond + " DO\n")
	if len(node.Children) > 1 {
		body := node.Children[1]
		for _, stmt := range body.Children {
			sb.WriteString(g.generateNode(stmt, level+1))
		}
	}
	sb.WriteString(indent + "END WHILE;\n")
	return sb.String()
}

// ============================================
// FOR
// ============================================

func (g *Generator) generateFor(node parser.ASTNode, level int) string {
	var sb strings.Builder
	indent := strings.Repeat(g.indent, level)
	if len(node.Children) == 0 {
		return ""
	}
	var varName string
	var expr string
	for i, child := range node.Children {
		if i == 0 && child.Type == parser.IdentifierNode {
			varName = g.getIdentifierName(child)
		} else if child.Type != parser.BlockNode {
			expr = g.generateNode(child, 0)
		}
	}
	if varName != "" {
		sb.WriteString(indent + "FOR " + varName + " AS " + expr + " DO\n")
	} else {
		sb.WriteString(indent + "FOR " + expr + " DO\n")
	}
	for _, child := range node.Children {
		if child.Type == parser.BlockNode {
			for _, stmt := range child.Children {
				sb.WriteString(g.generateNode(stmt, level+1))
			}
		}
	}
	sb.WriteString(indent + "END FOR;\n")
	return sb.String()
}

// ============================================
// RETURN
// ============================================

func (g *Generator) generateReturn(node parser.ASTNode, level int) string {
	if len(node.Children) > 0 {
		expr := g.generateNode(node.Children[0], 0)
		return strings.Repeat(g.indent, level) + "RETURN " + expr + ";\n"
	}
	return strings.Repeat(g.indent, level) + "RETURN;\n"
}

// ============================================
// THROW
// ============================================

func (g *Generator) generateThrow(node parser.ASTNode, level int) string {
	if len(node.Children) > 0 {
		expr := g.generateNode(node.Children[0], 0)
		return strings.Repeat(g.indent, level) + "THROW " + expr + ";\n"
	}
	return strings.Repeat(g.indent, level) + "THROW;\n"
}

// ============================================
// CALL
// ============================================

func (g *Generator) generateCall(node parser.ASTNode, level int) string {
	var name string
	var args []string
	for i, child := range node.Children {
		if i == 0 && child.Type == parser.IdentifierNode {
			name = g.getIdentifierName(child)
		} else {
			args = append(args, g.generateNode(child, 0))
		}
	}
	if name == "" {
		return ""
	}
	return strings.Repeat(g.indent, level) + "CALL " + name + "(" + strings.Join(args, ", ") + ");\n"
}

// ============================================
// BLOCK
// ============================================

func (g *Generator) generateBlock(node parser.ASTNode, level int) string {
	var sb strings.Builder
	for _, child := range node.Children {
		sb.WriteString(g.generateNode(child, level))
	}
	return sb.String()
}

// ============================================
// EXPRESSIONS
// ============================================

func (g *Generator) generateBinaryOp(node parser.ASTNode) string {
	if len(node.Children) != 2 {
		return node.Token
	}
	left := g.generateNode(node.Children[0], 0)
	right := g.generateNode(node.Children[1], 0)
	return left + " " + node.Token + " " + right
}

func (g *Generator) generateUnaryOp(node parser.ASTNode) string {
	if len(node.Children) == 0 {
		return node.Token
	}
	operand := g.generateNode(node.Children[0], 0)
	return node.Token + " " + operand
}

func (g *Generator) generateCast(node parser.ASTNode) string {
	if len(node.Children) < 2 {
		return "CAST"
	}
	expr := g.generateNode(node.Children[0], 0)
	typ := g.generateNode(node.Children[1], 0)
	return "CAST(" + expr + " AS " + typ + ")"
}

func (g *Generator) generateCase(node parser.ASTNode) string {
	var sb strings.Builder

	// Cek apakah ini simple CASE atau searched CASE
	isSimpleCase := false
	var caseExpr string
	var whenNodes []parser.ASTNode
	var elseNode parser.ASTNode

	for i, child := range node.Children {
		if i == 0 && child.Type != parser.WhenNode {
			// Simple CASE: ada expression setelah CASE
			isSimpleCase = true
			caseExpr = g.generateNode(child, 0)
		} else if child.Type == parser.WhenNode {
			whenNodes = append(whenNodes, child)
		} else if child.Type == parser.BlockNode && child.Token == "else" {
			elseNode = child
		}
	}

	// Tulis CASE
	if isSimpleCase {
		sb.WriteString("CASE " + caseExpr)
	} else {
		sb.WriteString("CASE")
	}

	// Tulis WHEN clauses (masing-masing di baris baru dengan indent)
	for _, when := range whenNodes {
		sb.WriteString("\n        " + g.generateWhen(when))
	}

	// Tulis ELSE jika ada
	if elseNode.Type != "" && len(elseNode.Children) > 0 {
		elseExpr := g.generateNode(elseNode.Children[0], 0)
		sb.WriteString("\n        ELSE " + elseExpr)
	}

	// Tulis END
	sb.WriteString("\n    END")

	return sb.String()
}

func (g *Generator) generateWhen(node parser.ASTNode) string {
	if len(node.Children) < 2 {
		return "WHEN"
	}
	cond := g.generateNode(node.Children[0], 0)
	result := g.generateNode(node.Children[1], 0)
	return "WHEN " + cond + " THEN " + result
}

func (g *Generator) generateBetween(node parser.ASTNode) string {
	if len(node.Children) < 3 {
		return "BETWEEN"
	}
	expr := g.generateNode(node.Children[0], 0)
	low := g.generateNode(node.Children[1], 0)
	high := g.generateNode(node.Children[2], 0)
	if node.Not {
		return expr + " NOT BETWEEN " + low + " AND " + high
	}
	return expr + " BETWEEN " + low + " AND " + high
}

func (g *Generator) generateLike(node parser.ASTNode) string {
	if len(node.Children) < 2 {
		return "LIKE"
	}
	expr := g.generateNode(node.Children[0], 0)
	pattern := g.generateNode(node.Children[1], 0)
	if node.Not {
		return expr + " NOT LIKE " + pattern
	}
	return expr + " LIKE " + pattern
}

func (g *Generator) generateIn(node parser.ASTNode) string {
	if len(node.Children) < 2 {
		return "IN"
	}
	expr := g.generateNode(node.Children[0], 0)
	var values []string
	for i := 1; i < len(node.Children); i++ {
		values = append(values, g.generateNode(node.Children[i], 0))
	}
	if node.Not {
		return expr + " NOT IN (" + strings.Join(values, ", ") + ")"
	}
	return expr + " IN (" + strings.Join(values, ", ") + ")"
}

func (g *Generator) generateCoalesce(node parser.ASTNode) string {
	var args []string
	for _, child := range node.Children {
		args = append(args, g.generateNode(child, 0))
	}
	return "COALESCE(" + strings.Join(args, ", ") + ")"
}

func (g *Generator) generateNullIf(node parser.ASTNode) string {
	if len(node.Children) < 2 {
		return "NULLIF()"
	}
	arg1 := g.generateNode(node.Children[0], 0)
	arg2 := g.generateNode(node.Children[1], 0)
	return "NULLIF(" + arg1 + ", " + arg2 + ")"
}

func (g *Generator) generateFunctionCall(node parser.ASTNode) string {
	var funcName string
	var args []string

	// Ambil nama function
	if val, ok := node.Value.(string); ok {
		funcName = val
	} else if len(node.Children) > 0 {
		funcName = g.getIdentifierName(node.Children[0])
	}

	// Ambil arguments (skip first child = function name)
	for i, child := range node.Children {
		if i == 0 {
			continue // skip function name identifier
		}
		// Generate each argument
		arg := g.generateNode(child, 0)
		if arg != "" {
			args = append(args, arg)
		}
	}

	if funcName == "" {
		return "()"
	}
	return funcName + "(" + strings.Join(args, ", ") + ")"
}

// ============================================
// FIELD REFERENCE
// ============================================

func (g *Generator) generateFieldReference(node parser.ASTNode) string {
	if node.Value != nil {
		if path, ok := node.Value.(string); ok {
			return path
		}
	}
	var parts []string
	for _, child := range node.Children {
		parts = append(parts, g.generateNode(child, 0))
	}
	return strings.Join(parts, ".")
}

// ============================================
// HELPERS
// ============================================

func (g *Generator) getIdentifierName(node parser.ASTNode) string {
	if node.Type == parser.IdentifierNode {
		if val, ok := node.Value.(string); ok && val != "" {
			return val
		}
		return node.Token
	}
	return ""
}

func (g *Generator) formatLiteral(node parser.ASTNode) string {
	if val, ok := node.Value.(string); ok {
		// Mode parameter tidak boleh di-quote
		if val == "IN" || val == "OUT" || val == "INOUT" {
			return val
		}
		if strings.HasPrefix(val, "'") && strings.HasSuffix(val, "'") {
			return val
		}
		return "'" + val + "'"
	}
	if num, ok := node.Value.(float64); ok {
		if num == float64(int(num)) {
			return fmt.Sprintf("%d", int(num))
		}
		return fmt.Sprintf("%v", num)
	}
	return node.Token
}

func (g *Generator) generateChildren(children []parser.ASTNode, level int) string {
	var sb strings.Builder
	for _, child := range children {
		sb.WriteString(g.generateNode(child, level))
	}
	return sb.String()
}

func (g *Generator) generatePropagate(node parser.ASTNode, level int) string {
	indent := strings.Repeat(g.indent, level)

	if len(node.Children) == 0 {
		return indent + "PROPAGATE;\n"
	}

	var exprs []string
	for _, child := range node.Children {
		exprs = append(exprs, g.generateNode(child, 0))
	}

	return indent + "PROPAGATE " + strings.Join(exprs, ", ") + ";\n"
}

// ============================================
// MOVE
// ============================================

func (g *Generator) generateMove(node parser.ASTNode, level int) string {
	indent := strings.Repeat(g.indent, level)
	if len(node.Children) >= 2 {
		target := g.generateNode(node.Children[0], 0)
		source := g.generateNode(node.Children[1], 0)
		return indent + "MOVE " + target + " TO " + source + ";\n"
	}
	return indent + "MOVE;\n"
}

// ============================================
// CONTINUE
// ============================================

func (g *Generator) generateContinue(node parser.ASTNode, level int) string {
	indent := strings.Repeat(g.indent, level)
	if node.Value != nil {
		return indent + "CONTINUE " + node.Value.(string) + ";\n"
	}
	return indent + "CONTINUE;\n"
}

// ============================================
// BREAK
// ============================================

func (g *Generator) generateBreak(node parser.ASTNode, level int) string {
	return strings.Repeat(g.indent, level) + "BREAK;\n"
}

// ============================================
// LABEL
// ============================================

func (g *Generator) generateLabel(node parser.ASTNode, level int) string {
	indent := strings.Repeat(g.indent, level)
	if node.Value != nil {
		return indent + "LABEL " + node.Value.(string) + ";\n"
	}
	return indent + "LABEL;\n"
}

```

---

### File: `pkg/analyzer/analyzer.go`

```go
package analyzer

import (
	"esql-ast-tool/pkg/parser"
	"sort"
	"strconv"
)

// ============================================
// Struct Definitions
// ============================================

type VariableInfo struct {
	Type string `json:"type"`
	Line int    `json:"line"`
}

type FunctionInfo struct {
	Parameters []string `json:"parameters"`
	ReturnType string   `json:"returnType"`
	Line       int      `json:"line"`
}

type ProcedureInfo struct {
	Parameters []string `json:"parameters"`
	Line       int      `json:"line"`
}

type UsageInfo struct {
	Name     string `json:"name"`
	Location string `json:"location"` // "line:col"
	Context  string `json:"context"`  // "DECLARE", "SET", "IF", "CALL", etc.
	Line     int    `json:"line"`
	Column   int    `json:"column"`
}

type ModuleInfo struct {
	Name       string   `json:"name"`
	Line       int      `json:"line"`
	Procedures []string `json:"procedures"`
	Functions  []string `json:"functions"`
	Variables  []string `json:"variables"`
}

// ============================================
// Analysis Result
// ============================================

type AnalysisResult struct {
	Variables        map[string]VariableInfo  `json:"variables"`
	Functions        map[string]FunctionInfo  `json:"functions"`
	Procedures       map[string]ProcedureInfo `json:"procedures"`
	UsedVariables    []string                 `json:"usedVariables"`
	DefinedVariables []string                 `json:"definedVariables"`
	Issues           []string                 `json:"issues"`

	// Relational Info
	CallGraph        map[string][]string    `json:"callGraph"`        // Caller -> Callees
	ReverseCallGraph map[string][]string    `json:"reverseCallGraph"` // Callee -> Callers
	UsageMap         map[string][]UsageInfo `json:"usageMap"`         // Name -> Usage locations
	ImpactMap        map[string][]string    `json:"impactMap"`        // Change X -> Affects Y
	ModuleInfo       ModuleInfo             `json:"moduleInfo"`
}

// ============================================
// Analyzer
// ============================================

type Analyzer struct {
	variables        map[string]VariableInfo
	functions        map[string]FunctionInfo
	procedures       map[string]ProcedureInfo
	usedVariables    map[string]bool
	definedVariables map[string]bool
	issues           []string

	// Relational tracking
	callGraph        map[string][]string
	reverseCallGraph map[string][]string
	usageMap         map[string][]UsageInfo
	currentScope     string

	// Module info - gunakan map untuk cegah duplikasi
	moduleName       string
	moduleLine       int
	moduleProcedures map[string]bool // ← Ubah ke map
	moduleFunctions  map[string]bool // ← Ubah ke map
	moduleVariables  map[string]bool // ← Ubah ke map
}

func NewAnalyzer() *Analyzer {
	return &Analyzer{
		variables:        make(map[string]VariableInfo),
		functions:        make(map[string]FunctionInfo),
		procedures:       make(map[string]ProcedureInfo),
		usedVariables:    make(map[string]bool),
		definedVariables: make(map[string]bool),
		issues:           []string{},
		callGraph:        make(map[string][]string),
		reverseCallGraph: make(map[string][]string),
		usageMap:         make(map[string][]UsageInfo),
		moduleProcedures: make(map[string]bool), // ← Ubah
		moduleFunctions:  make(map[string]bool), // ← Ubah
		moduleVariables:  make(map[string]bool), // ← Ubah
	}
}

// ============================================
// Main Analysis
// ============================================

func (a *Analyzer) Analyze(program parser.Program) AnalysisResult {
	// Hanya analyze sekali
	for _, stmt := range program.Statements {
		a.analyzeNode(stmt)
	}

	// Konversi map ke slice untuk output
	var procedures []string
	for name := range a.moduleProcedures {
		procedures = append(procedures, name)
	}
	sort.Strings(procedures)

	var functions []string
	for name := range a.moduleFunctions {
		functions = append(functions, name)
	}
	sort.Strings(functions)

	var variables []string
	for name := range a.moduleVariables {
		variables = append(variables, name)
	}
	sort.Strings(variables)

	impactMap := a.buildImpactMap()

	return AnalysisResult{
		Variables:        a.sortVariables(),
		Functions:        a.sortFunctions(),
		Procedures:       a.sortProcedures(),
		UsedVariables:    a.sortUsedVariables(),
		DefinedVariables: a.sortDefinedVariables(),
		Issues:           a.issues,
		CallGraph:        a.sortCallGraph(),
		ReverseCallGraph: a.sortReverseCallGraph(),
		UsageMap:         a.usageMap,
		ImpactMap:        impactMap,
		ModuleInfo: ModuleInfo{
			Name:       a.moduleName,
			Line:       a.moduleLine,
			Procedures: procedures,
			Functions:  functions,
			Variables:  variables,
		},
	}
}

// ============================================
// Sorting Helpers
// ============================================

func (a *Analyzer) sortVariables() map[string]VariableInfo {
	var names []string
	for name := range a.variables {
		names = append(names, name)
	}
	sort.Strings(names)
	result := make(map[string]VariableInfo)
	for _, name := range names {
		result[name] = a.variables[name]
	}
	return result
}

func (a *Analyzer) sortFunctions() map[string]FunctionInfo {
	var names []string
	for name := range a.functions {
		names = append(names, name)
	}
	sort.Strings(names)
	result := make(map[string]FunctionInfo)
	for _, name := range names {
		result[name] = a.functions[name]
	}
	return result
}

func (a *Analyzer) sortProcedures() map[string]ProcedureInfo {
	var names []string
	for name := range a.procedures {
		names = append(names, name)
	}
	sort.Strings(names)
	result := make(map[string]ProcedureInfo)
	for _, name := range names {
		result[name] = a.procedures[name]
	}
	return result
}

func (a *Analyzer) sortUsedVariables() []string {
	var vars []string
	for v := range a.usedVariables {
		vars = append(vars, v)
	}
	sort.Strings(vars)
	return vars
}

func (a *Analyzer) sortDefinedVariables() []string {
	var vars []string
	for v := range a.definedVariables {
		vars = append(vars, v)
	}
	sort.Strings(vars)
	return vars
}

func (a *Analyzer) sortCallGraph() map[string][]string {
	result := make(map[string][]string)
	var keys []string
	for k := range a.callGraph {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		vals := a.callGraph[k]
		sort.Strings(vals)
		result[k] = vals
	}
	return result
}

func (a *Analyzer) sortReverseCallGraph() map[string][]string {
	result := make(map[string][]string)
	var keys []string
	for k := range a.reverseCallGraph {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		vals := a.reverseCallGraph[k]
		sort.Strings(vals)
		result[k] = vals
	}
	return result
}

// ============================================
// Impact Analysis
// ============================================

func (a *Analyzer) buildImpactMap() map[string][]string {
	impact := make(map[string][]string)

	// Untuk setiap variable, cari di mana dia digunakan
	for varName := range a.variables {
		if usages, ok := a.usageMap[varName]; ok {
			// Gunakan map untuk deduplicate
			seen := make(map[string]bool)
			var affected []string
			for _, u := range usages {
				key := u.Context + " at line " + strconv.Itoa(u.Line)
				if !seen[key] {
					seen[key] = true
					affected = append(affected, key)
				}
			}
			if len(affected) > 0 {
				// SORT affected
				sort.Strings(affected)
				impact[varName] = affected
			}
		}
	}

	// Untuk setiap procedure/function, cari siapa yang memanggilnya
	for name := range a.procedures {
		if callers, ok := a.reverseCallGraph[name]; ok {
			seen := make(map[string]bool)
			var unique []string
			for _, caller := range callers {
				if !seen[caller] {
					seen[caller] = true
					unique = append(unique, caller)
				}
			}
			sort.Strings(unique)
			impact[name] = unique
		}
	}
	for name := range a.functions {
		if callers, ok := a.reverseCallGraph[name]; ok {
			seen := make(map[string]bool)
			var unique []string
			for _, caller := range callers {
				if !seen[caller] {
					seen[caller] = true
					unique = append(unique, caller)
				}
			}
			sort.Strings(unique)
			impact[name] = unique
		}
	}

	return impact
}

// ============================================
// Node Analysis
// ============================================

func (a *Analyzer) analyzeNode(node parser.ASTNode) {
	if node.Type == "" {
		return
	}

	switch node.Type {
	case parser.ModuleNode:
		if len(node.Children) > 0 && node.Children[0].Type == parser.IdentifierNode {
			if name, ok := node.Children[0].Value.(string); ok {
				a.moduleName = name
				a.moduleLine = node.Span.Start.Line
			}
		}
		for _, child := range node.Children {
			a.analyzeNode(child)
		}

	case parser.DeclareNode:
		if len(node.Children) > 0 && node.Children[0].Type == parser.IdentifierNode {
			if name, ok := node.Children[0].Value.(string); ok {
				varType := "UNKNOWN"
				if len(node.Children) > 1 && node.Children[1].Type == parser.IdentifierNode {
					if v, ok := node.Children[1].Value.(string); ok {
						varType = v
					}
				}
				a.variables[name] = VariableInfo{
					Type: varType,
					Line: node.Span.Start.Line,
				}
				a.definedVariables[name] = true
				a.moduleVariables[name] = true // ← Pakai map

				a.usageMap[name] = append(a.usageMap[name], UsageInfo{
					Name:     name,
					Location: formatLocation(node.Span.Start.Line, node.Span.Start.Column),
					Context:  "DECLARE",
					Line:     node.Span.Start.Line,
					Column:   node.Span.Start.Column,
				})
			}
		}
	case parser.SetNode:
		if len(node.Children) > 0 {
			if node.Children[0].Type == parser.BlockNode && len(node.Children[0].Children) > 0 {
				a.analyzeNode(node.Children[0].Children[0])
			}
		}
		if len(node.Children) > 1 {
			if node.Children[1].Type == parser.BlockNode && len(node.Children[1].Children) > 0 {
				a.analyzeNode(node.Children[1].Children[0])
			}
		}

	case parser.IfNode:
		if len(node.Children) > 0 {
			a.analyzeNode(node.Children[0])
		}
		if len(node.Children) > 1 {
			for _, child := range node.Children[1].Children {
				a.analyzeNode(child)
			}
		}
		if len(node.Children) > 2 {
			for _, child := range node.Children[2].Children {
				a.analyzeNode(child)
			}
		}

	case parser.FunctionNode:
		if len(node.Children) > 0 && node.Children[0].Type == parser.IdentifierNode {
			if name, ok := node.Children[0].Value.(string); ok && name != "" {
				returnType := "UNKNOWN"
				// Cari ReturnTypeNode di children
				for _, child := range node.Children {
					if child.Type == parser.ReturnTypeNode {
						if val, ok := child.Value.(string); ok && val != "" {
							returnType = val
						} else if child.Token != "" {
							returnType = child.Token
						} else if len(child.Children) > 0 {
							if val, ok := child.Children[0].Value.(string); ok && val != "" {
								returnType = val
							} else if child.Children[0].Token != "" {
								returnType = child.Children[0].Token
							}
						}
					}
				}
				a.functions[name] = FunctionInfo{
					Parameters: []string{},
					ReturnType: returnType,
					Line:       node.Span.Start.Line,
				}
				a.moduleFunctions[name] = true
				a.currentScope = name
			}
		}
		for _, child := range node.Children {
			if child.Type != parser.IdentifierNode {
				a.analyzeNode(child)
			}
		}
		a.currentScope = ""

	case parser.ProcedureNode:
		if len(node.Children) > 0 && node.Children[0].Type == parser.IdentifierNode {
			if name, ok := node.Children[0].Value.(string); ok {
				a.procedures[name] = ProcedureInfo{
					Parameters: []string{},
					Line:       node.Span.Start.Line,
				}
				a.moduleProcedures[name] = true // ← Pakai map
				a.currentScope = name
			}
		}
		for _, child := range node.Children {
			if child.Type != parser.IdentifierNode {
				a.analyzeNode(child)
			}
		}
		a.currentScope = ""

	case parser.CallNode:
		var callee string
		if len(node.Children) > 0 && node.Children[0].Type == parser.IdentifierNode {
			if val, ok := node.Children[0].Value.(string); ok && val != "" {
				callee = val
			}
			if callee == "" && node.Children[0].Token != "" {
				callee = node.Children[0].Token
			}
		}
		if callee != "" {
			// ✅ HANYA jika currentScope tidak kosong
			if a.currentScope != "" {
				a.callGraph[a.currentScope] = appendUnique(a.callGraph[a.currentScope], callee)
				a.reverseCallGraph[callee] = appendUnique(a.reverseCallGraph[callee], a.currentScope)
			}
			// ❌ JANGAN tambahkan ke MAIN secara otomatis
			// Hanya CALL di root yang boleh masuk MAIN
		}
		for _, child := range node.Children {
			a.analyzeNode(child)
		}

	case parser.IdentifierNode:
		if name, ok := node.Value.(string); ok {
			a.usedVariables[name] = true
			context := "USAGE"
			if a.currentScope != "" {
				context = "USAGE in " + a.currentScope
			}

			// Cek apakah sudah ada entry dengan line/column yang sama
			existing := false
			for _, u := range a.usageMap[name] {
				if u.Line == node.Span.Start.Line && u.Column == node.Span.Start.Column {
					existing = true
					break
				}
			}
			if !existing {
				a.usageMap[name] = append(a.usageMap[name], UsageInfo{
					Name:     name,
					Location: formatLocation(node.Span.Start.Line, node.Span.Start.Column),
					Context:  context,
					Line:     node.Span.Start.Line,
					Column:   node.Span.Start.Column,
				})
			}
		}

	case parser.FunctionCallNode:
		var funcName string
		if val, ok := node.Value.(string); ok && val != "" {
			funcName = val
		}
		if funcName == "" && len(node.Children) > 0 {
			if node.Children[0].Type == parser.IdentifierNode {
				if val, ok := node.Children[0].Value.(string); ok && val != "" {
					funcName = val
				} else if node.Children[0].Token != "" {
					funcName = node.Children[0].Token
				}
			}
		}
		if funcName != "" {
			if _, exists := a.functions[funcName]; !exists {
				a.functions[funcName] = FunctionInfo{
					Parameters: []string{},
					ReturnType: "BUILTIN",
					Line:       node.Span.Start.Line,
				}
			}
			// ✅ HANYA jika currentScope tidak kosong
			if a.currentScope != "" {
				a.callGraph[a.currentScope] = appendUnique(a.callGraph[a.currentScope], funcName)
				a.reverseCallGraph[funcName] = appendUnique(a.reverseCallGraph[funcName], a.currentScope)
			}
			// ❌ JANGAN tambahkan ke MAIN secara otomatis
		}
		for _, child := range node.Children {
			a.analyzeNode(child)
		}

	case parser.CastNode, parser.CaseNode, parser.WhenNode,
		parser.IsNullNode, parser.IsNotNullNode, parser.BetweenNode,
		parser.LikeNode, parser.InNode, parser.CoalesceNode, parser.NullIfNode:
		for _, child := range node.Children {
			a.analyzeNode(child)
		}
	}

	for _, child := range node.Children {
		a.analyzeNode(child)
	}
}

// ============================================
// Helper Functions
// ============================================

func appendUnique(slice []string, item string) []string {
	for _, s := range slice {
		if s == item {
			return slice
		}
	}
	return append(slice, item)
}

func formatLocation(line, column int) string {
	return strconv.Itoa(line) + ":" + strconv.Itoa(column)
}

```

---

### File: `pkg/printer/printer.go`

```go
package printer

import (
	"fmt"
	"strings"

	"esql-ast-tool/pkg/parser"
)

type Printer struct {
	indent string
}

func NewPrinter() *Printer {
	return &Printer{indent: "  "}
}

func (p *Printer) PrintProgram(program parser.Program) string {
	var sb strings.Builder

	// If root exists, print it
	if program.Root.Type != "" {
		sb.WriteString(p.printNode(program.Root, 0))
	} else {
		// Fallback to old behavior
		sb.WriteString(fmt.Sprintf("Program %s\n", programSpan(program)))
		for _, stmt := range program.Statements {
			sb.WriteString(p.printNode(stmt, 1))
		}
	}

	return sb.String()
}

func (p *Printer) printNode(node parser.ASTNode, level int) string {
	if node.Type == "" {
		return ""
	}

	indent := strings.Repeat(p.indent, level)
	var sb strings.Builder

	// Special handling untuk node yang butuh format khusus
	switch node.Type {
	case parser.IfNode:
		spanStr := node.Span.String()
		sb.WriteString(fmt.Sprintf("%sIf %s\n", indent, spanStr))
		for _, child := range node.Children {
			sb.WriteString(p.printNode(child, level+1))
		}
		return sb.String()

	case parser.ParameterNode:
		if len(node.Children) >= 3 {
			mode := p.getNodeValue(node.Children[0])
			name := p.getNodeValue(node.Children[1])
			typ := p.getNodeValue(node.Children[2])
			return fmt.Sprintf("%sParameter [%s]:\n%s  Mode: %s\n%s  Name: %s\n%s  Type: %s\n",
				indent, node.Span.String(),
				indent, mode,
				indent, name,
				indent, typ)
		}
		return fmt.Sprintf("%sParameter [%s]\n", indent, node.Span.String())

	case parser.DeclareNode:
		var name, typ string
		for _, child := range node.Children {
			if child.Type == parser.IdentifierNode {
				if name == "" {
					name = p.getNodeValue(child)
				} else {
					typ = p.getNodeValue(child)
				}
			}
		}
		return fmt.Sprintf("%sDeclare [%s]: %s %s\n",
			indent, node.Span.String(), name, typ)

	case parser.SetNode:
		var target, value string
		for _, child := range node.Children {
			if child.Type == parser.BlockNode {
				if child.Token == "target" && len(child.Children) > 0 {
					target = p.getNodeValue(child.Children[0])
				} else if child.Token == "value" && len(child.Children) > 0 {
					// Special handling untuk FunctionCall
					if child.Children[0].Type == parser.FunctionCallNode {
						value = p.getFunctionCallString(child.Children[0])
					} else {
						value = p.getNodeValue(child.Children[0])
					}
				}
			}
		}
		if value == "" {
			value = "?"
		}
		return fmt.Sprintf("%sSet [%s]: %s = %s\n",
			indent, node.Span.String(), target, value)

	case parser.CallNode:
		var callee string
		var args []string
		for i, child := range node.Children {
			if i == 0 {
				callee = p.getNodeValue(child)
			} else {
				args = append(args, p.getNodeValue(child))
			}
		}
		if len(args) > 0 {
			return fmt.Sprintf("%sCall [%s]: %s(%s)\n",
				indent, node.Span.String(), callee, strings.Join(args, ", "))
		}
		return fmt.Sprintf("%sCall [%s]: %s\n",
			indent, node.Span.String(), callee)

	case parser.FunctionCallNode:
		funcName := p.getNodeValue(node)
		var args []string
		// Skip the first child (function name identifier)
		for i, child := range node.Children {
			if i == 0 {
				continue
			}
			if child.Type != "" {
				val := p.getNodeValue(child)
				if val != "" && val != "?" {
					args = append(args, val)
				}
			}
		}
		if len(args) > 0 {
			return fmt.Sprintf("%sFunctionCall [%s]: %s(%s)\n",
				indent, node.Span.String(), funcName, strings.Join(args, ", "))
		}
		return fmt.Sprintf("%sFunctionCall [%s]: %s()\n",
			indent, node.Span.String(), funcName)
	case parser.ReturnTypeNode:
		return fmt.Sprintf("%sReturnType: %s [%s]\n",
			indent, p.getNodeValue(node), node.Span.String())

	case parser.ReturnNode:
		var val string
		if len(node.Children) > 0 {
			val = p.getNodeValue(node.Children[0])
		}
		if val != "" {
			return fmt.Sprintf("%sReturn [%s]: %s\n",
				indent, node.Span.String(), val)
		}
		return fmt.Sprintf("%sReturn [%s]\n", indent, node.Span.String())
	case parser.PropagateNode:
		if len(node.Children) > 0 {
			var exprs []string
			for _, child := range node.Children {
				exprs = append(exprs, p.getNodeValue(child))
			}
			return fmt.Sprintf("%sPropagate [%s]: %s\n",
				indent, node.Span.String(), strings.Join(exprs, ", "))
		}
		return fmt.Sprintf("%sPropagate [%s]\n", indent, node.Span.String())

	}
	// Get display name untuk node lainnya
	displayName := p.getDisplayName(node)
	spanStr := node.Span.String()
	sb.WriteString(fmt.Sprintf("%s%s %s\n", indent, displayName, spanStr))

	// Print children with proper indentation
	for _, child := range node.Children {
		sb.WriteString(p.printNode(child, level+1))
	}

	return sb.String()
}

// getDisplayName returns the display name for a node
func (p *Printer) getDisplayName(node parser.ASTNode) string {
	switch node.Type {
	case parser.CreateNode:
		return "CreateModule"
	case parser.ModuleNode:
		if node.Value != nil {
			if val, ok := node.Value.(string); ok && val == "COMPUTE" {
				return "ComputeModule"
			}
		}
		return "Module"
	case parser.DeclareNode:
		return "Declare"
	case parser.SetNode:
		return "Set"
	case parser.BlockNode:
		switch node.Token {
		case "condition":
			return "Condition"
		case "then":
			return "Then"
		case "else":
			return "Else"
		case "target":
			return "Target"
		case "value":
			return "Value"
		default:
			return "Block"
		}
	case parser.ComparisonNode:
		if node.Token != "" {
			return fmt.Sprintf("Comparison (%s)", node.Token)
		}
		return "Comparison"
	case parser.FieldReferenceNode:
		if node.Value != nil {
			if path, ok := node.Value.(string); ok {
				return fmt.Sprintf("FieldReference (%s)", path)
			}
		}
		return "FieldReference"
	case parser.IdentifierNode:
		if name, ok := node.Value.(string); ok && name != "" {
			return fmt.Sprintf("Identifier: %s", name)
		}
		return fmt.Sprintf("Identifier: %s", node.Token)
	case parser.LiteralNode:
		if valStr, ok := node.Value.(string); ok {
			return fmt.Sprintf("Literal: '%s'", valStr)
		}
		if num, ok := node.Value.(float64); ok {
			return fmt.Sprintf("Literal: %v", num)
		}
		return "Literal"
	case parser.CastNode:
		return "Cast"
	case parser.CaseNode:
		return "Case"
	case parser.WhenNode:
		return "When"
	case parser.IsNullNode:
		return "IsNull"
	case parser.IsNotNullNode:
		return "IsNotNull"
	case parser.BetweenNode:
		if node.Not {
			return "Between (NOT)"
		}
		return "Between"
	case parser.UnaryOpNode:
		return fmt.Sprintf("UnaryOp (%s)", node.Token)
	case parser.BinaryOpNode:
		return fmt.Sprintf("BinaryOp (%s)", node.Token)
	case parser.ParenthesizedNode:
		return "Parenthesized"
	case parser.LikeNode:
		if node.Not {
			return "Like (NOT)"
		}
		return "Like"
	case parser.InNode:
		if node.Not {
			return "In (NOT)"
		}
		return "In"
	case parser.CoalesceNode:
		return "Coalesce"
	case parser.NullIfNode:
		return "NullIf"
	case parser.FunctionCallNode:
		return "FunctionCall"
	case parser.ProcedureNode:
		return "Procedure"
	case parser.FunctionNode:
		return "Function"
	case parser.CallNode:
		return "Call"
	case parser.ReturnNode:
		return "Return"
	case parser.ThrowNode:
		return "Throw"
	case parser.PropagateNode:
		return "Propagate"
	case parser.MoveNode:
		return "Move"
	case parser.ContinueNode:
		return "Continue"
	case parser.BreakNode:
		return "Break"
	case parser.LabelNode:
		return "Label"
	case parser.WhileNode:
		return "While"
	case parser.ForNode:
		return "For"
	case parser.EnvironmentNode:
		return "Environment"
	case parser.DatabaseNode:
		return "Database"
	case parser.PassthruNode:
		return "Passthru"
	case parser.OtherwiseNode:
		return "Otherwise"
	case parser.ElseIfNode:
		return "ElseIf"
	case parser.ElseNode:
		return "Else"
	case parser.ReturnTypeNode:
		return "ReturnType"
	default:
		if node.Token != "" {
			return string(node.Type) + " (" + node.Token + ")"
		}
		return string(node.Type)
	}
}

// getNodeValue returns the string value of a node
func (p *Printer) getNodeValue(node parser.ASTNode) string {
	switch node.Type {
	case parser.LiteralNode:
		if val, ok := node.Value.(string); ok {
			// Jangan tambahkan quotes untuk MODE (IN, OUT, INOUT)
			if val == "IN" || val == "OUT" || val == "INOUT" {
				return val
			}
			// Cek apakah string sudah punya quotes
			if strings.HasPrefix(val, "'") && strings.HasSuffix(val, "'") {
				return val
			}
			return "'" + val + "'"
		}
		if num, ok := node.Value.(float64); ok {
			// Format number tanpa desimal jika integer
			if num == float64(int(num)) {
				return fmt.Sprintf("%d", int(num))
			}
			return fmt.Sprintf("%v", num)
		}
		return node.Token

	case parser.IdentifierNode:
		if val, ok := node.Value.(string); ok && val != "" {
			return val
		}
		return node.Token

	case parser.FieldReferenceNode:
		if val, ok := node.Value.(string); ok {
			return val
		}
		// Build from children
		var parts []string
		for _, child := range node.Children {
			parts = append(parts, p.getNodeValue(child))
		}
		return strings.Join(parts, ".")

	case parser.FunctionCallNode:
		if val, ok := node.Value.(string); ok {
			return val
		}
		if len(node.Children) > 0 {
			return p.getNodeValue(node.Children[0])
		}
		return node.Token

	case parser.BinaryOpNode, parser.ComparisonNode:
		if len(node.Children) >= 2 {
			left := p.getNodeValue(node.Children[0])
			right := p.getNodeValue(node.Children[1])
			return left + " " + node.Token + " " + right
		}
		return node.Token

	case parser.CastNode:
		if len(node.Children) >= 2 {
			expr := p.getNodeValue(node.Children[0])
			typ := p.getNodeValue(node.Children[1])
			return "CAST(" + expr + " AS " + typ + ")"
		}
		return "CAST"

	case parser.CaseNode:
		return "CASE"

	case parser.WhenNode:
		return "WHEN"

	case parser.IsNullNode:
		return "IS NULL"

	case parser.IsNotNullNode:
		return "IS NOT NULL"

	case parser.BetweenNode:
		if len(node.Children) >= 3 {
			expr := p.getNodeValue(node.Children[0])
			lower := p.getNodeValue(node.Children[1])
			upper := p.getNodeValue(node.Children[2])
			if node.Not {
				return expr + " NOT BETWEEN " + lower + " AND " + upper
			}
			return expr + " BETWEEN " + lower + " AND " + upper
		}
		return "BETWEEN"

	case parser.LikeNode:
		if len(node.Children) >= 2 {
			expr := p.getNodeValue(node.Children[0])
			pattern := p.getNodeValue(node.Children[1])
			if node.Not {
				return expr + " NOT LIKE " + pattern
			}
			return expr + " LIKE " + pattern
		}
		return "LIKE"

	case parser.InNode:
		if len(node.Children) >= 2 {
			expr := p.getNodeValue(node.Children[0])
			var values []string
			for i := 1; i < len(node.Children); i++ {
				values = append(values, p.getNodeValue(node.Children[i]))
			}
			if node.Not {
				return expr + " NOT IN (" + strings.Join(values, ", ") + ")"
			}
			return expr + " IN (" + strings.Join(values, ", ") + ")"
		}
		return "IN"

	case parser.CoalesceNode:
		var args []string
		for _, child := range node.Children {
			args = append(args, p.getNodeValue(child))
		}
		return "COALESCE(" + strings.Join(args, ", ") + ")"

	case parser.NullIfNode:
		if len(node.Children) >= 2 {
			arg1 := p.getNodeValue(node.Children[0])
			arg2 := p.getNodeValue(node.Children[1])
			return "NULLIF(" + arg1 + ", " + arg2 + ")"
		}
		return "NULLIF()"

	case parser.ReturnTypeNode:
		if val, ok := node.Value.(string); ok && val != "" {
			return val
		}
		return node.Token

	default:
		if node.Value != nil {
			if str, ok := node.Value.(string); ok && str != "" {
				return str
			}
		}
		if node.Token != "" {
			return node.Token
		}
		if len(node.Children) > 0 {
			return p.getNodeValue(node.Children[0])
		}
		return "?"
	}
}

func programSpan(program parser.Program) string {
	if len(program.Statements) == 0 {
		return "[1:1 - 1:1]"
	}
	start := program.Statements[0].Span.Start
	end := program.Statements[len(program.Statements)-1].Span.End
	return fmt.Sprintf("[%d:%d - %d:%d]", start.Line, start.Column, end.Line, end.Column)
}

func (p *Printer) getFunctionCallString(node parser.ASTNode) string {
	if node.Type != parser.FunctionCallNode {
		return p.getNodeValue(node)
	}

	funcName := p.getNodeValue(node)
	var args []string
	// Skip the first child (function name identifier)
	for i, child := range node.Children {
		if i == 0 {
			continue
		}
		if child.Type != "" {
			val := p.getNodeValue(child)
			if val != "" && val != "?" {
				args = append(args, val)
			}
		}
	}
	if len(args) > 0 {
		return funcName + "(" + strings.Join(args, ", ") + ")"
	}
	return funcName + "()"
}

```

---

### File: `cmd/esql-ast/main.go`

```go
package main

import (
	"flag"
	"fmt"
	"os"
	"sort"

	"esql-ast-tool/pkg/analyzer"
	"esql-ast-tool/pkg/generator"
	"esql-ast-tool/pkg/parser"
	"esql-ast-tool/pkg/printer"
	"esql-ast-tool/pkg/refactor"
)

func main() {
	var (
		file     = flag.String("f", "", "ESQL file to parse")
		code     = flag.String("c", "", "ESQL code string to parse")
		jsonOut  = flag.Bool("json", false, "Output AST as JSON")
		pretty   = flag.Bool("pretty", false, "Pretty print AST")
		analyze  = flag.Bool("analyze", false, "Perform analysis")
		validate = flag.Bool("validate", false, "Validate AST")
		generate = flag.Bool("generate", false, "Generate ESQL code from AST")
		output   = flag.String("o", "", "Output file")
		debug    = flag.Bool("debug", false, "Enable debug output")

		// Refactoring flags
		refactorCmd = flag.String("refactor", "", "Refactoring operation: suggest, dead-code, rename")
		refactorOld = flag.String("old", "", "Old name (for rename operations)")
		refactorNew = flag.String("new", "", "New name (for rename operations)")
		dryRun      = flag.Bool("dry-run", false, "Preview changes without applying")
		apply       = flag.Bool("apply", false, "Apply refactoring changes to file")

		explain    = flag.Bool("explain", false, "Explain the code in natural language")
		search     = flag.String("search", "", "Search type: procedure, function, variable, call, unused")
		searchName = flag.String("search-name", "", "Name to search for")
	)
	flag.Parse()

	// Enable debug mode if flag is set
	if *debug {
		parser.DebugMode = true
	}

	if *file == "" && *code == "" {
		fmt.Println("Usage: esql-ast -f <file> or -c <code>")
		flag.PrintDefaults()
		os.Exit(1)
	}

	var input string
	if *file != "" {
		data, err := os.ReadFile(*file)
		if err != nil {
			fmt.Printf("Error reading file: %v\n", err)
			os.Exit(1)
		}
		input = string(data)
	} else {
		input = *code
	}

	p := parser.NewParser(input)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		fmt.Println("Parse errors:")
		for _, err := range p.Errors() {
			fmt.Printf("  %s\n", err)
		}
		os.Exit(1)
	}

	var result string

	if *jsonOut {
		jsonData, err := program.ToJSON()
		if err != nil {
			fmt.Printf("Error marshaling JSON: %v\n", err)
			os.Exit(1)
		}
		result = string(jsonData)
	} else if *pretty {
		pr := printer.NewPrinter()
		result = pr.PrintProgram(program)
	} else if *generate {
		gen := generator.NewGenerator()
		result = gen.Generate(program)
	} else {
		result = fmt.Sprintf("Program has %d statements\n", len(program.Statements))
	}

	if *analyze {
		an := analyzer.NewAnalyzer()
		analysisResult := an.Analyze(program)

		result += "\n=== Analysis Results ===\n"
		result += fmt.Sprintf("Variables defined: %d\n", len(analysisResult.DefinedVariables))
		result += fmt.Sprintf("Variables used: %d\n", len(analysisResult.UsedVariables))

		if len(analysisResult.Variables) > 0 {
			result += "\nVariables:\n"
			var varNames []string
			for name := range analysisResult.Variables {
				varNames = append(varNames, name)
			}
			sort.Strings(varNames)
			for _, name := range varNames {
				info := analysisResult.Variables[name]
				result += fmt.Sprintf("  %s: %s (line %d)\n", name, info.Type, info.Line)
			}
		}

		if len(analysisResult.Functions) > 0 {
			result += "\nFunctions:\n"
			var funcNames []string
			for name := range analysisResult.Functions {
				funcNames = append(funcNames, name)
			}
			sort.Strings(funcNames)
			for _, name := range funcNames {
				info := analysisResult.Functions[name]
				result += fmt.Sprintf("  %s (line %d)\n", name, info.Line)
			}
		}

		if len(analysisResult.Procedures) > 0 {
			result += "\nProcedures:\n"
			var procNames []string
			for name := range analysisResult.Procedures {
				procNames = append(procNames, name)
			}
			sort.Strings(procNames)
			for _, name := range procNames {
				info := analysisResult.Procedures[name]
				result += fmt.Sprintf("  %s (line %d)\n", name, info.Line)
			}
		}

		if len(analysisResult.CallGraph) > 0 {
			result += "\n=== Call Graph (Caller -> Callees) ===\n"
			var callers []string
			for caller := range analysisResult.CallGraph {
				callers = append(callers, caller)
			}
			sort.Strings(callers)
			for _, caller := range callers {
				callees := analysisResult.CallGraph[caller]
				result += fmt.Sprintf("  %s -> %v\n", caller, callees)
			}
		}

		if len(analysisResult.ReverseCallGraph) > 0 {
			result += "\n=== Reverse Call Graph (Callee -> Callers) ===\n"
			var callees []string
			for callee := range analysisResult.ReverseCallGraph {
				callees = append(callees, callee)
			}
			sort.Strings(callees)
			for _, callee := range callees {
				callers := analysisResult.ReverseCallGraph[callee]
				result += fmt.Sprintf("  %s <- %v\n", callee, callers)
			}
		}

		if len(analysisResult.ImpactMap) > 0 {
			result += "\n=== Impact Analysis (Change X -> Affects Y) ===\n"
			var keys []string
			for name := range analysisResult.ImpactMap {
				keys = append(keys, name)
			}
			sort.Strings(keys)
			for _, name := range keys {
				affected := analysisResult.ImpactMap[name]
				result += fmt.Sprintf("  %s -> %v\n", name, affected)
			}
		}

		if analysisResult.ModuleInfo.Name != "" {
			result += "\n=== Module Info ===\n"
			result += fmt.Sprintf("  Name: %s\n", analysisResult.ModuleInfo.Name)
			result += fmt.Sprintf("  Line: %d\n", analysisResult.ModuleInfo.Line)
			if len(analysisResult.ModuleInfo.Procedures) > 0 {
				result += fmt.Sprintf("  Procedures: %v\n", analysisResult.ModuleInfo.Procedures)
			}
			if len(analysisResult.ModuleInfo.Functions) > 0 {
				result += fmt.Sprintf("  Functions: %v\n", analysisResult.ModuleInfo.Functions)
			}
			if len(analysisResult.ModuleInfo.Variables) > 0 {
				result += fmt.Sprintf("  Variables: %v\n", analysisResult.ModuleInfo.Variables)
			}
		}
	}

	// SEARCH
	if *search != "" {
		an := analyzer.NewAnalyzer()
		analysisResult := an.Analyze(program)
		var searchResult refactor.SearchResult

		switch *search {
		case "procedure":
			if *searchName == "" {
				result += "\n❌ Please provide -search-name for procedure search\n"
			} else {
				searchResult = refactor.SearchProcedure(program, *searchName)
			}
		case "function":
			if *searchName == "" {
				result += "\n❌ Please provide -search-name for function search\n"
			} else {
				searchResult = refactor.SearchFunction(program, *searchName)
			}
		case "variable":
			if *searchName == "" {
				result += "\n❌ Please provide -search-name for variable search\n"
			} else {
				searchResult = refactor.SearchVariable(program, *searchName)
			}
		case "call":
			if *searchName == "" {
				result += "\n❌ Please provide -search-name for call search\n"
			} else {
				searchResult = refactor.SearchCall(program, *searchName)
			}
		case "unused":
			searchResult = refactor.SearchUnused(program, analysisResult)
		default:
			result += fmt.Sprintf("\n❌ Unknown search type: %s\n", *search)
			result += "Available search types:\n"
			result += "  procedure   - Search for procedure\n"
			result += "  function    - Search for function\n"
			result += "  variable    - Search for variable\n"
			result += "  call        - Search for CALL statements\n"
			result += "  unused      - Find all unused code\n"
		}

		if searchResult.TotalCount > 0 || searchResult.Message != "" {
			result += refactor.FormatSearchResult(searchResult)
		}
	}
	// REFACTORING
	if *refactorCmd != "" {
		an := analyzer.NewAnalyzer()
		analysisResult := an.Analyze(program)

		switch *refactorCmd {
		case "suggest":
			refactorResult := refactor.Suggest(program, analysisResult)
			result += refactor.FormatSuggestions(refactorResult)

		case "dead-code":
			refactorResult := refactor.Suggest(program, analysisResult)
			result += refactor.FormatDeadCode(refactorResult)

		case "rename":
			if *refactorOld == "" || *refactorNew == "" {
				result += "\n❌ Please provide both -old and -new names\n"
				result += "Usage: esql-ast -f file.esql -refactor rename -old <oldName> -new <newName>\n"
			} else {
				// Try variable rename first, then procedure, then function
				renameResult := refactor.RenameVariable(program, *refactorOld, *refactorNew, *dryRun)
				if !renameResult.Success {
					renameResult = refactor.RenameProcedure(program, *refactorOld, *refactorNew, *dryRun)
				}
				if !renameResult.Success {
					renameResult = refactor.RenameFunction(program, *refactorOld, *refactorNew, *dryRun)
				}
				result += refactor.FormatRenameResult(renameResult, *dryRun)

				// Apply changes if -apply is set and not dry-run
				if *apply && !*dryRun && renameResult.Success {
					newContent := refactor.ApplyRenameChanges(input, renameResult)
					if *output != "" {
						err := os.WriteFile(*output, []byte(newContent), 0644)
						if err != nil {
							result += fmt.Sprintf("\n❌ Error writing to output file: %v\n", err)
						} else {
							result += fmt.Sprintf("\n✅ Changes saved to: %s\n", *output)
						}
					} else if *file != "" {
						// Backup original
						backupFile := *file + ".bak"
						err := os.WriteFile(backupFile, []byte(input), 0644)
						if err != nil {
							result += fmt.Sprintf("\n⚠️ Could not create backup: %v\n", err)
						} else {
							result += fmt.Sprintf("\n📁 Backup saved to: %s\n", backupFile)
						}

						// Write changes
						err = os.WriteFile(*file, []byte(newContent), 0644)
						if err != nil {
							result += fmt.Sprintf("\n❌ Error writing to file: %v\n", err)
						} else {
							result += fmt.Sprintf("\n✅ File updated: %s\n", *file)
						}
					}
				} else if *apply && renameResult.Success {
					result += "\n💡 Dry-run mode: changes not applied. Remove -dry-run to apply.\n"
				}
			}

		default:
			result += fmt.Sprintf("\n❌ Unknown refactor operation: %s\n", *refactorCmd)
			result += "Available operations:\n"
			result += "  suggest     - Show refactoring suggestions\n"
			result += "  dead-code   - Show dead code analysis\n"
			result += "  rename      - Rename variable/procedure/function\n"
		}
	}

	if *validate {
		an := analyzer.NewAnalyzer()
		analysisResult := an.Analyze(program)

		result += "\n=== Validation Results ===\n"
		if len(analysisResult.Issues) > 0 {
			for _, issue := range analysisResult.Issues {
				result += fmt.Sprintf("  %s\n", issue)
			}
		} else {
			result += "  No issues found\n"
		}
	}

	if *explain {
		an := analyzer.NewAnalyzer()
		analysisResult := an.Analyze(program)
		explanationResult := refactor.Explain(program, analysisResult)
		result += refactor.FormatExplanation(explanationResult)
	}

	if *output != "" && *refactorCmd != "rename" {
		err := os.WriteFile(*output, []byte(result), 0644)
		if err != nil {
			fmt.Printf("Error writing output: %v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Println(result)
	}
}

```

---

### File: `internal/token/token.go`

```go
package token

const (
	// Keywords and identifiers
	IDENTIFIER = "IDENTIFIER"
	NUMBER     = "NUMBER"
	STRING     = "STRING"
	EOF        = "EOF"
	ILLEGAL    = "ILLEGAL"

	// Operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	ASTERISK = "*"
	SLASH    = "/"
	MODULO   = "%"
	EQ       = "=="
	NOT_EQ   = "!="
	LT       = "<"
	GT       = ">"
	LTE      = "<="
	GTE      = ">="

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"
	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"
	LBRACKET  = "["
	RBRACKET  = "]"
	DOT       = "."

	// Keywords
	CREATE      = "CREATE"
	DECLARE     = "DECLARE"
	SET         = "SET"
	IF          = "IF"
	ELSE        = "ELSE"
	ELSEIF      = "ELSEIF"
	WHILE       = "WHILE"
	FOR         = "FOR"
	RETURN      = "RETURN"
	THROW       = "THROW"
	PROPAGATE   = "PROPAGATE"
	MOVE        = "MOVE"
	CONTINUE    = "CONTINUE"
	BREAK       = "BREAK"
	LABEL       = "LABEL"
	MODULE      = "MODULE"
	FUNCTION    = "FUNCTION"
	PROCEDURE   = "PROCEDURE"
	CALL        = "CALL"
	BEGIN       = "BEGIN"
	END         = "END"
	THEN        = "THEN"
	DO          = "DO"
	AS          = "AS"
	RETURNS     = "RETURNS"
	DEFAULT     = "DEFAULT"
	COMPUTE     = "COMPUTE"
	FIELD       = "FIELD"
	ENVIRONMENT = "ENVIRONMENT"
	DATABASE    = "DATABASE"
	PASSTHRU    = "PASSTHRU"
	AND         = "AND"
	OR          = "OR"
	NOT         = "NOT"
	TO          = "TO"
	CAST        = "CAST"
	CASE        = "CASE"
	WHEN        = "WHEN"
	IS          = "IS"
	ISNULL      = "ISNULL"
	NOTNULL     = "NOTNULL"
	BETWEEN     = "BETWEEN"
	LIKE        = "LIKE"
	COALESCE    = "COALESCE"
	NULLIF      = "NULLIF"
	IN          = "IN"
	OUT         = "OUT"
	INOUT       = "INOUT"
	TERMINAL    = "TERMINAL"
	DELETE      = "DELETE"
	NONE        = "NONE"
)

type Token struct {
	Type    string
	Literal string
	Line    int
	Column  int
	Pos     int
}

func LookupIdent(ident string) string {
	keywords := map[string]string{
		"CREATE":      CREATE,
		"DECLARE":     DECLARE,
		"SET":         SET,
		"IF":          IF,
		"ELSE":        ELSE,
		"ELSEIF":      ELSEIF,
		"WHILE":       WHILE,
		"FOR":         FOR,
		"RETURN":      RETURN,
		"THROW":       THROW,
		"PROPAGATE":   PROPAGATE,
		"MOVE":        MOVE,
		"CONTINUE":    CONTINUE,
		"BREAK":       BREAK,
		"LABEL":       LABEL,
		"MODULE":      MODULE,
		"FUNCTION":    FUNCTION,
		"PROCEDURE":   PROCEDURE,
		"CALL":        CALL,
		"BEGIN":       BEGIN,
		"END":         END,
		"THEN":        THEN,
		"DO":          DO,
		"AS":          AS,
		"RETURNS":     RETURNS,
		"DEFAULT":     DEFAULT,
		"COMPUTE":     COMPUTE,
		"FIELD":       FIELD,
		"ENVIRONMENT": ENVIRONMENT,
		"DATABASE":    DATABASE,
		"PASSTHRU":    PASSTHRU,
		"AND":         AND,
		"OR":          OR,
		"NOT":         NOT,
		"TO":          TO,
		"CAST":        CAST,
		"CASE":        CASE,
		"WHEN":        WHEN,
		"IS":          IS,
		"NULL":        IDENTIFIER,
		"BETWEEN":     BETWEEN,
		"LIKE":        LIKE,
		"IN":          IN,
		"OUT":         OUT,
		"INOUT":       INOUT,
		"COALESCE":    COALESCE,
		"NULLIF":      NULLIF,
		"TERMINAL":    TERMINAL,
		"DELETE":      DELETE,
		"NONE":        NONE,
	}
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENTIFIER
}

```

---

