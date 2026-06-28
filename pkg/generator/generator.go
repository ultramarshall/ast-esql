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
