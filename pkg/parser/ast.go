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
