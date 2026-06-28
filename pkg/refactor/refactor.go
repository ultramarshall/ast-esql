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
