package main

import (
	"flag"
	"fmt"
	"os"
	"sort"

	"esql-ast-tool/pkg/analyzer"
	"esql-ast-tool/pkg/generator"
	"esql-ast-tool/pkg/parser"
	"esql-ast-tool/pkg/printer"
	"esql-ast-tool/pkg/refactor"
)

func main() {
	var (
		file     = flag.String("f", "", "ESQL file to parse")
		code     = flag.String("c", "", "ESQL code string to parse")
		jsonOut  = flag.Bool("json", false, "Output AST as JSON")
		pretty   = flag.Bool("pretty", false, "Pretty print AST")
		analyze  = flag.Bool("analyze", false, "Perform analysis")
		validate = flag.Bool("validate", false, "Validate AST")
		generate = flag.Bool("generate", false, "Generate ESQL code from AST")
		output   = flag.String("o", "", "Output file")
		debug    = flag.Bool("debug", false, "Enable debug output")

		// Refactoring flags
		refactorCmd = flag.String("refactor", "", "Refactoring operation: suggest, dead-code")
	)
	flag.Parse()

	// Enable debug mode if flag is set
	if *debug {
		parser.DebugMode = true
	}

	if *file == "" && *code == "" {
		fmt.Println("Usage: esql-ast -f <file> or -c <code>")
		flag.PrintDefaults()
		os.Exit(1)
	}

	var input string
	if *file != "" {
		data, err := os.ReadFile(*file)
		if err != nil {
			fmt.Printf("Error reading file: %v\n", err)
			os.Exit(1)
		}
		input = string(data)
	} else {
		input = *code
	}

	p := parser.NewParser(input)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		fmt.Println("Parse errors:")
		for _, err := range p.Errors() {
			fmt.Printf("  %s\n", err)
		}
		os.Exit(1)
	}

	var result string

	if *jsonOut {
		jsonData, err := program.ToJSON()
		if err != nil {
			fmt.Printf("Error marshaling JSON: %v\n", err)
			os.Exit(1)
		}
		result = string(jsonData)
	} else if *pretty {
		pr := printer.NewPrinter()
		result = pr.PrintProgram(program)
	} else if *generate {
		gen := generator.NewGenerator()
		result = gen.Generate(program)
	} else {
		result = fmt.Sprintf("Program has %d statements\n", len(program.Statements))
	}

	if *analyze {
		an := analyzer.NewAnalyzer()
		analysisResult := an.Analyze(program)

		result += "\n=== Analysis Results ===\n"
		result += fmt.Sprintf("Variables defined: %d\n", len(analysisResult.DefinedVariables))
		result += fmt.Sprintf("Variables used: %d\n", len(analysisResult.UsedVariables))

		if len(analysisResult.Variables) > 0 {
			result += "\nVariables:\n"
			var varNames []string
			for name := range analysisResult.Variables {
				varNames = append(varNames, name)
			}
			sort.Strings(varNames)
			for _, name := range varNames {
				info := analysisResult.Variables[name]
				result += fmt.Sprintf("  %s: %s (line %d)\n", name, info.Type, info.Line)
			}
		}

		if len(analysisResult.Functions) > 0 {
			result += "\nFunctions:\n"
			var funcNames []string
			for name := range analysisResult.Functions {
				funcNames = append(funcNames, name)
			}
			sort.Strings(funcNames)
			for _, name := range funcNames {
				info := analysisResult.Functions[name]
				result += fmt.Sprintf("  %s (line %d)\n", name, info.Line)
			}
		}

		if len(analysisResult.Procedures) > 0 {
			result += "\nProcedures:\n"
			var procNames []string
			for name := range analysisResult.Procedures {
				procNames = append(procNames, name)
			}
			sort.Strings(procNames)
			for _, name := range procNames {
				info := analysisResult.Procedures[name]
				result += fmt.Sprintf("  %s (line %d)\n", name, info.Line)
			}
		}

		if len(analysisResult.CallGraph) > 0 {
			result += "\n=== Call Graph (Caller -> Callees) ===\n"
			var callers []string
			for caller := range analysisResult.CallGraph {
				callers = append(callers, caller)
			}
			sort.Strings(callers)
			for _, caller := range callers {
				callees := analysisResult.CallGraph[caller]
				result += fmt.Sprintf("  %s -> %v\n", caller, callees)
			}
		}

		if len(analysisResult.ReverseCallGraph) > 0 {
			result += "\n=== Reverse Call Graph (Callee -> Callers) ===\n"
			var callees []string
			for callee := range analysisResult.ReverseCallGraph {
				callees = append(callees, callee)
			}
			sort.Strings(callees)
			for _, callee := range callees {
				callers := analysisResult.ReverseCallGraph[callee]
				result += fmt.Sprintf("  %s <- %v\n", callee, callers)
			}
		}

		if len(analysisResult.ImpactMap) > 0 {
			result += "\n=== Impact Analysis (Change X -> Affects Y) ===\n"
			var keys []string
			for name := range analysisResult.ImpactMap {
				keys = append(keys, name)
			}
			sort.Strings(keys)
			for _, name := range keys {
				affected := analysisResult.ImpactMap[name]
				result += fmt.Sprintf("  %s -> %v\n", name, affected)
			}
		}

		if analysisResult.ModuleInfo.Name != "" {
			result += "\n=== Module Info ===\n"
			result += fmt.Sprintf("  Name: %s\n", analysisResult.ModuleInfo.Name)
			result += fmt.Sprintf("  Line: %d\n", analysisResult.ModuleInfo.Line)
			if len(analysisResult.ModuleInfo.Procedures) > 0 {
				result += fmt.Sprintf("  Procedures: %v\n", analysisResult.ModuleInfo.Procedures)
			}
			if len(analysisResult.ModuleInfo.Functions) > 0 {
				result += fmt.Sprintf("  Functions: %v\n", analysisResult.ModuleInfo.Functions)
			}
			if len(analysisResult.ModuleInfo.Variables) > 0 {
				result += fmt.Sprintf("  Variables: %v\n", analysisResult.ModuleInfo.Variables)
			}
		}
	}

	// ============================================
	// REFACTORING
	// ============================================
	if *refactorCmd != "" {
		an := analyzer.NewAnalyzer()
		analysisResult := an.Analyze(program)

		switch *refactorCmd {
		case "suggest":
			refactorResult := refactor.Suggest(program, analysisResult)
			result += refactor.FormatSuggestions(refactorResult)

		case "dead-code":
			refactorResult := refactor.Suggest(program, analysisResult)
			result += refactor.FormatDeadCode(refactorResult)

		default:
			result += fmt.Sprintf("\n❌ Unknown refactor operation: %s\n", *refactorCmd)
			result += "Available operations:\n"
			result += "  suggest     - Show refactoring suggestions\n"
			result += "  dead-code   - Show dead code analysis\n"
		}
	}

	if *validate {
		an := analyzer.NewAnalyzer()
		analysisResult := an.Analyze(program)

		result += "\n=== Validation Results ===\n"
		if len(analysisResult.Issues) > 0 {
			for _, issue := range analysisResult.Issues {
				result += fmt.Sprintf("  %s\n", issue)
			}
		} else {
			result += "  No issues found\n"
		}
	}

	if *output != "" {
		err := os.WriteFile(*output, []byte(result), 0644)
		if err != nil {
			fmt.Printf("Error writing output: %v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Println(result)
	}
}
