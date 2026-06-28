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

	// Group by type
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
