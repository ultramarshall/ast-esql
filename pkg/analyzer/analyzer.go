package analyzer

import (
	"esql-ast-tool/pkg/parser"
)

type AnalysisResult struct {
	Variables        map[string]VariableInfo  `json:"variables"`
	Functions        map[string]FunctionInfo  `json:"functions"`
	Procedures       map[string]ProcedureInfo `json:"procedures"`
	UsedVariables    []string                 `json:"usedVariables"`
	DefinedVariables []string                 `json:"definedVariables"`
	Issues           []string                 `json:"issues"`
}

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

type Analyzer struct {
	variables        map[string]VariableInfo
	functions        map[string]FunctionInfo
	procedures       map[string]ProcedureInfo
	usedVariables    map[string]bool
	definedVariables map[string]bool
	issues           []string
}

func NewAnalyzer() *Analyzer {
	return &Analyzer{
		variables:        make(map[string]VariableInfo),
		functions:        make(map[string]FunctionInfo),
		procedures:       make(map[string]ProcedureInfo),
		usedVariables:    make(map[string]bool),
		definedVariables: make(map[string]bool),
		issues:           []string{},
	}
}

func (a *Analyzer) Analyze(program parser.Program) AnalysisResult {
	for _, stmt := range program.Statements {
		a.analyzeNode(stmt)
	}

	usedVars := make([]string, 0, len(a.usedVariables))
	for v := range a.usedVariables {
		usedVars = append(usedVars, v)
	}

	definedVars := make([]string, 0, len(a.definedVariables))
	for v := range a.definedVariables {
		definedVars = append(definedVars, v)
	}

	return AnalysisResult{
		Variables:        a.variables,
		Functions:        a.functions,
		Procedures:       a.procedures,
		UsedVariables:    usedVars,
		DefinedVariables: definedVars,
		Issues:           a.issues,
	}
}

func (a *Analyzer) analyzeNode(node parser.ASTNode) {
	if node.Type == "" {
		return
	}

	switch node.Type {
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
					Line: node.Line,
				}
				a.definedVariables[name] = true
			}
		}

	case parser.SetNode:
		if len(node.Children) > 0 {
			a.analyzeNode(node.Children[0])
		}
		if len(node.Children) > 1 {
			a.analyzeNode(node.Children[1])
		}

	case parser.FunctionNode:
		if len(node.Children) > 0 && node.Children[0].Type == parser.IdentifierNode {
			if name, ok := node.Children[0].Value.(string); ok {
				a.functions[name] = FunctionInfo{
					Parameters: []string{},
					ReturnType: "UNKNOWN",
					Line:       node.Line,
				}
			}
		}

	case parser.ProcedureNode:
		if len(node.Children) > 0 && node.Children[0].Type == parser.IdentifierNode {
			if name, ok := node.Children[0].Value.(string); ok {
				a.procedures[name] = ProcedureInfo{
					Parameters: []string{},
					Line:       node.Line,
				}
			}
		}

	case parser.IdentifierNode:
		if name, ok := node.Value.(string); ok {
			a.usedVariables[name] = true
		}

	case parser.FunctionCallNode:
		if name, ok := node.Value.(string); ok {
			if _, exists := a.functions[name]; !exists {
				a.functions[name] = FunctionInfo{
					Parameters: []string{},
					ReturnType: "BUILTIN",
					Line:       node.Line,
				}
			}
		}
	case parser.CastNode:
		// Cast doesn't define variables, just analyze children
		for _, child := range node.Children {
			a.analyzeNode(child)
		}

	case parser.CaseNode:
		// Analyze semua children (WHEN clauses dan ELSE)
		for _, child := range node.Children {
			a.analyzeNode(child)
		}
	case parser.WhenNode:
		// Analyze condition dan result
		for _, child := range node.Children {
			a.analyzeNode(child)
		}

	}

	for _, child := range node.Children {
		a.analyzeNode(child)
	}
}
