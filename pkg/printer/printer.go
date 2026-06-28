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
	sb.WriteString(fmt.Sprintf("Program %s\n", programSpan(program)))

	for _, stmt := range program.Statements {
		sb.WriteString(p.printNode(stmt, 1))
	}

	return sb.String()
}

func (p *Printer) printNode(node parser.ASTNode, level int) string {
	if node.Type == "" {
		return ""
	}

	indent := strings.Repeat(p.indent, level)
	var sb strings.Builder

	displayName := string(node.Type)
	switch node.Type {
	case parser.CreateNode:
		displayName = "CreateModule"
	case parser.ModuleNode:
		if node.Value != nil {
			if val, ok := node.Value.(string); ok && val == "COMPUTE" {
				displayName = "ComputeModule"
			} else {
				displayName = "Module"
			}
		} else {
			displayName = "Module"
		}
	case parser.DeclareNode:
		displayName = "Declare"
	case parser.SetNode:
		displayName = "Set"
	case parser.IfNode:
		displayName = "If"
		// Print span
		spanStr := node.Span.String()
		sb.WriteString(fmt.Sprintf("%s%s %s\n", indent, displayName, spanStr))
		// Print children
		for _, child := range node.Children {
			sb.WriteString(p.printNode(child, level+1))
		}
		return sb.String()

	case parser.BlockNode:
		switch node.Token {
		case "condition":
			displayName = "Condition"
		case "then":
			displayName = "Then"
		case "else":
			displayName = "Else"
		case "target":
			displayName = "Target"
		case "value":
			displayName = "Value"
		default:
			displayName = "Block"
		}
	case parser.ComparisonNode:
		if node.Token != "" {
			displayName = fmt.Sprintf("Comparison (%s)", node.Token)
		} else {
			displayName = "Comparison"
		}

	case parser.FieldReferenceNode:
		if node.Value != nil {
			if path, ok := node.Value.(string); ok {
				displayName = fmt.Sprintf("FieldReference (%s)", path)
			} else {
				displayName = "FieldReference"
			}
		} else {
			displayName = "FieldReference"
		}

	case parser.IdentifierNode:
		if name, ok := node.Value.(string); ok && name != "" && name != "error" {
			displayName = fmt.Sprintf("Identifier: %s", name)
		} else {
			displayName = fmt.Sprintf("Identifier: %s", node.Token)
		}

	case parser.LiteralNode:
		if valStr, ok := node.Value.(string); ok {
			displayName = fmt.Sprintf("Literal: '%s'", valStr)
		} else if num, ok := node.Value.(float64); ok {
			displayName = fmt.Sprintf("Literal: %v", num)
		} else {
			displayName = "Literal"
		}

	case parser.CastNode:
		displayName = "Cast"
	case parser.CaseNode:
		displayName = "Case"
	case parser.WhenNode:
		displayName = "When"
	case parser.IsNullNode:
		displayName = "IsNull"
	case parser.IsNotNullNode:
		displayName = "IsNotNull"
	case parser.BetweenNode:
		if node.Not {
			displayName = "Between (NOT)"
		} else {
			displayName = "Between"
		}
	case parser.UnaryOpNode:
		displayName = fmt.Sprintf("UnaryOp (%s)", node.Token)
	case parser.BinaryOpNode:
		displayName = fmt.Sprintf("BinaryOp (%s)", node.Token)
	case parser.ParenthesizedNode:
		displayName = "Parenthesized"

	case parser.LikeNode:
		if node.Not {
			displayName = "Like (NOT)"
		} else {
			displayName = "Like"
		}

	case parser.InNode:
		if node.Not {
			displayName = "In (NOT)"
		} else {
			displayName = "In"
		}

	}

	spanStr := node.Span.String()
	sb.WriteString(fmt.Sprintf("%s%s %s\n", indent, displayName, spanStr))

	for _, child := range node.Children {
		sb.WriteString(p.printNode(child, level+1))
	}

	return sb.String()
}

func programSpan(program parser.Program) string {
	if len(program.Statements) == 0 {
		return "[1:1 - 1:1]"
	}
	start := program.Statements[0].Span.Start
	end := program.Statements[len(program.Statements)-1].Span.End
	return fmt.Sprintf("[%d:%d - %d:%d]", start.Line, start.Column, end.Line, end.Column)
}
