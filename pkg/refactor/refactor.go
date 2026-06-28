package refactor

import (
	"fmt"
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

	// Traverse AST and find all occurrences of the variable
	searchAndReplace(program.Statements, oldName, newName, &locations, &occurrences, "variable")

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

	searchAndReplace(program.Statements, oldName, newName, &locations, &occurrences, "procedure")

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

	searchAndReplace(program.Statements, oldName, newName, &locations, &occurrences, "function")

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
		switch node.Type {
		case parser.IdentifierNode:
			if val, ok := node.Value.(string); ok && val == oldName {
				*occurrences++
				context := getContext(node)
				*locations = append(*locations, RenameLocation{
					Line:    node.Span.Start.Line,
					Column:  node.Span.Start.Column,
					OldText: oldName,
					NewText: newName,
					Context: context,
				})
			}

		case parser.DeclareNode:
			// Check if this declares the variable
			if targetType == "variable" && len(node.Children) > 0 {
				if node.Children[0].Type == parser.IdentifierNode {
					if val, ok := node.Children[0].Value.(string); ok && val == oldName {
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

		case parser.FunctionNode:
			if targetType == "function" && len(node.Children) > 0 {
				if node.Children[0].Type == parser.IdentifierNode {
					if val, ok := node.Children[0].Value.(string); ok && val == oldName {
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

		case parser.ProcedureNode:
			if targetType == "procedure" && len(node.Children) > 0 {
				if node.Children[0].Type == parser.IdentifierNode {
					if val, ok := node.Children[0].Value.(string); ok && val == oldName {
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

		case parser.CallNode:
			if targetType == "procedure" && len(node.Children) > 0 {
				if node.Children[0].Type == parser.IdentifierNode {
					if val, ok := node.Children[0].Value.(string); ok && val == oldName {
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

		case parser.FunctionCallNode:
			if targetType == "function" {
				if val, ok := node.Value.(string); ok && val == oldName {
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
		}

		// Recursively search children
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
