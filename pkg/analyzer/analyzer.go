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
			if name, ok := node.Children[0].Value.(string); ok {
				a.functions[name] = FunctionInfo{
					Parameters: []string{},
					ReturnType: "UNKNOWN",
					Line:       node.Span.Start.Line,
				}
				a.moduleFunctions[name] = true // ← Pakai map
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
		if len(node.Children) > 0 && node.Children[0].Type == parser.IdentifierNode {
			if callee, ok := node.Children[0].Value.(string); ok {
				caller := a.currentScope
				if caller != "" {
					a.callGraph[caller] = appendUnique(a.callGraph[caller], callee)
					a.reverseCallGraph[callee] = appendUnique(a.reverseCallGraph[callee], caller)
				}
			}
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
		if name, ok := node.Value.(string); ok {
			if _, exists := a.functions[name]; !exists {
				a.functions[name] = FunctionInfo{
					Parameters: []string{},
					ReturnType: "BUILTIN",
					Line:       node.Span.Start.Line,
				}
			}
			caller := a.currentScope
			if caller != "" {
				a.callGraph[caller] = appendUnique(a.callGraph[caller], name)
				a.reverseCallGraph[name] = appendUnique(a.reverseCallGraph[name], caller)
			}
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
