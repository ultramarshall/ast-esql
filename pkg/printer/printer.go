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
	sb.WriteString("Program [line: 1, col: 1]\n")

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

	// Determine display name
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
	}

	sb.WriteString(fmt.Sprintf("%s%s [line: %d, col: %d]\n", indent, displayName, node.Line, node.Column))

	// Print children
	for _, child := range node.Children {
		sb.WriteString(p.printNode(child, level+1))
	}

	return sb.String()
}
