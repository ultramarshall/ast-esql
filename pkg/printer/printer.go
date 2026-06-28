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
