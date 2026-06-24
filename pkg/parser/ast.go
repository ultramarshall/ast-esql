package parser

import (
	"encoding/json"
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
)

type ASTNode struct {
	Type     NodeType    `json:"type"`
	Value    interface{} `json:"value,omitempty"`
	Children []ASTNode   `json:"children,omitempty"`
	Line     int         `json:"line"`
	Column   int         `json:"column"`
	Token    string      `json:"-"`
}

func NewASTNode(nodeType NodeType, token string, line, column int) ASTNode {
	return ASTNode{
		Type:     nodeType,
		Children: []ASTNode{},
		Line:     line,
		Column:   column,
		Token:    token,
	}
}

func (n *ASTNode) AddChild(child ASTNode) {
	n.Children = append(n.Children, child)
}

func (n ASTNode) ToJSON() ([]byte, error) {
	return json.MarshalIndent(n, "", "  ")
}

type Program struct {
	Statements []ASTNode `json:"statements"`
}

func NewProgram() Program {
	return Program{
		Statements: []ASTNode{},
	}
}

func (p *Program) AddStatement(stmt ASTNode) {
	p.Statements = append(p.Statements, stmt)
}

func (p Program) ToJSON() ([]byte, error) {
	return json.MarshalIndent(p, "", "  ")
}
