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
	for i, stmt := range program.Statements {
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

	_ = strings.Repeat(g.indent, level) // avoid unused variable

	switch node.Type {
	case parser.CreateNode:
		return g.generateCreate(node, level)

	case parser.ModuleNode:
		return g.generateModule(node, level)

	case parser.DeclareNode:
		return g.generateDeclare(node, level)

	case parser.SetNode:
		return g.generateSet(node, level)

	case parser.IfNode:
		return g.generateIf(node, level)

	case parser.BlockNode:
		return g.generateBlock(node, level)

	case parser.IdentifierNode:
		if name, ok := node.Value.(string); ok && name != "" {
			return name
		}
		return node.Token

	case parser.FieldReferenceNode:
		return g.generateFieldReference(node)

	case parser.LiteralNode:
		if str, ok := node.Value.(string); ok {
			return "'" + str + "'"
		}
		if num, ok := node.Value.(float64); ok {
			return fmt.Sprintf("%v", num)
		}
		return node.Token

	case parser.ComparisonNode, parser.BinaryOpNode:
		return g.generateBinaryOp(node, level)

	case parser.CastNode:
		return g.generateCast(node, level)

	case parser.CaseNode:
		return g.generateCase(node, level)

	case parser.WhenNode:
		return g.generateWhen(node, false)

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
		if len(node.Children) >= 3 {
			expr := g.generateNode(node.Children[0], 0)
			lower := g.generateNode(node.Children[1], 0)
			upper := g.generateNode(node.Children[2], 0)
			return expr + " BETWEEN " + lower + " AND " + upper
		}
		return "BETWEEN"

	default:
		var sb strings.Builder
		if node.Token != "" {
			sb.WriteString(node.Token)
		}
		for _, child := range node.Children {
			sb.WriteString(" " + g.generateNode(child, level))
		}
		return sb.String()
	}
}

func (g *Generator) generateCreate(node parser.ASTNode, level int) string {
	var sb strings.Builder
	sb.WriteString("CREATE COMPUTE MODULE")
	for _, child := range node.Children {
		if child.Type == parser.ModuleNode {
			sb.WriteString(" " + g.generateModule(child, level))
			break
		}
	}
	return sb.String()
}

func (g *Generator) generateModule(node parser.ASTNode, level int) string {
	var sb strings.Builder
	moduleName := "UnnamedModule"
	for _, child := range node.Children {
		if child.Type == parser.IdentifierNode {
			if name, ok := child.Value.(string); ok && name != "" {
				moduleName = name
				break
			}
		}
	}

	sb.WriteString(moduleName + "\n")
	for _, child := range node.Children {
		if child.Type != parser.IdentifierNode {
			sb.WriteString(g.generateNode(child, level+1))
		}
	}
	sb.WriteString("END MODULE;\n")
	return sb.String()
}

func (g *Generator) generateDeclare(node parser.ASTNode, level int) string {
	var sb strings.Builder
	sb.WriteString(strings.Repeat(g.indent, level) + "DECLARE ")

	name := ""
	typ := ""
	for _, child := range node.Children {
		if child.Type == parser.IdentifierNode {
			if name == "" {
				if v, ok := child.Value.(string); ok {
					name = v
				} else {
					name = child.Token
				}
			} else {
				typ = child.Token
			}
		}
	}
	sb.WriteString(name + " " + typ + ";\n")
	return sb.String()
}

func (g *Generator) generateSet(node parser.ASTNode, level int) string {
	var sb strings.Builder
	sb.WriteString(strings.Repeat(g.indent, level) + "SET ")

	var target, value string
	for _, child := range node.Children {
		if child.Type == parser.BlockNode {
			if child.Token == "target" && len(child.Children) > 0 {
				target = g.generateNode(child.Children[0], 0)
			} else if child.Token == "value" && len(child.Children) > 0 {
				value = g.generateNode(child.Children[0], 0)
			}
		}
	}

	sb.WriteString(target + " = " + value + ";\n")
	return sb.String()
}

func (g *Generator) generateIf(node parser.ASTNode, level int) string {
	var sb strings.Builder
	indentStr := strings.Repeat(g.indent, level)

	sb.WriteString(indentStr + "IF ")

	// Generate condition (child 0)
	if len(node.Children) > 0 {
		cond := node.Children[0]
		sb.WriteString(g.generateNode(cond, 0))
	}

	sb.WriteString(" THEN\n")

	// Generate then block (child 1)
	if len(node.Children) > 1 {
		thenBlock := node.Children[1]
		for _, stmt := range thenBlock.Children {
			sb.WriteString(g.generateNode(stmt, level+1))
		}
	}

	// Generate else block if exists (child 2)
	if len(node.Children) > 2 {
		elseBlock := node.Children[2]
		sb.WriteString(indentStr + "ELSE\n")
		for _, stmt := range elseBlock.Children {
			sb.WriteString(g.generateNode(stmt, level+1))
		}
	}

	sb.WriteString(indentStr + "END IF;\n")
	return sb.String()
}

func (g *Generator) generateBlock(node parser.ASTNode, level int) string {
	var sb strings.Builder

	// Untuk BlockNode dengan token "else" di CASE expression
	if node.Token == "else" && len(node.Children) > 0 {
		sb.WriteString("ELSE " + g.generateNode(node.Children[0], 0))
		return sb.String()
	}

	// Untuk BlockNode biasa (then, target, value, condition)
	for _, child := range node.Children {
		sb.WriteString(g.generateNode(child, level))
	}
	return sb.String()
}

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

func (g *Generator) generateBinaryOp(node parser.ASTNode, level int) string {
	if len(node.Children) != 2 {
		return node.Token
	}
	left := g.generateNode(node.Children[0], 0)
	right := g.generateNode(node.Children[1], 0)
	return left + " " + node.Token + " " + right
}

func (g *Generator) generateCast(node parser.ASTNode, level int) string {
	var sb strings.Builder

	sb.WriteString("CAST(")

	// Generate expression
	if len(node.Children) > 0 {
		sb.WriteString(g.generateNode(node.Children[0], 0))
	}

	sb.WriteString(" AS ")

	// Generate type
	if len(node.Children) > 1 {
		sb.WriteString(g.generateNode(node.Children[1], 0))
	}

	sb.WriteString(")")

	return sb.String()
}

func (g *Generator) generateCase(node parser.ASTNode, level int) string {
	var sb strings.Builder

	sb.WriteString("CASE")

	// Cek apakah ini simple CASE (child 0 adalah expression, bukan WHEN)
	if len(node.Children) > 0 && node.Children[0].Type != parser.WhenNode {
		// Simple CASE: CASE expression
		sb.WriteString(" " + g.generateNode(node.Children[0], 0))

		// Generate WHEN clauses (mulai dari index 1)
		for i := 1; i < len(node.Children); i++ {
			child := node.Children[i]
			if child.Type == parser.WhenNode {
				sb.WriteString(" " + g.generateWhen(child, true))
			} else if child.Type == parser.BlockNode && child.Token == "else" {
				sb.WriteString(" " + g.generateNode(child, 0))
			}
		}
	} else {
		// Searched CASE: CASE WHEN condition THEN result ...
		for _, child := range node.Children {
			if child.Type == parser.WhenNode {
				sb.WriteString(" " + g.generateWhen(child, false))
			} else if child.Type == parser.BlockNode && child.Token == "else" {
				sb.WriteString(" " + g.generateNode(child, 0))
			}
		}
	}

	sb.WriteString(" END")
	return sb.String()
}

func (g *Generator) generateWhen(node parser.ASTNode, isSimpleCase bool) string {
	var sb strings.Builder

	sb.WriteString("WHEN ")

	if len(node.Children) >= 2 {
		if isSimpleCase {
			// Simple CASE: WHEN value THEN result
			sb.WriteString(g.generateNode(node.Children[0], 0))
			sb.WriteString(" THEN ")
			sb.WriteString(g.generateNode(node.Children[1], 0))
		} else {
			// Searched CASE: WHEN condition THEN result
			sb.WriteString(g.generateNode(node.Children[0], 0))
			sb.WriteString(" THEN ")
			sb.WriteString(g.generateNode(node.Children[1], 0))
		}
	}

	return sb.String()
}
