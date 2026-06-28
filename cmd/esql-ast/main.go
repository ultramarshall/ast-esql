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
		refactorCmd = flag.String("refactor", "", "Refactoring operation: suggest, dead-code, rename")
		refactorOld = flag.String("old", "", "Old name (for rename operations)")
		refactorNew = flag.String("new", "", "New name (for rename operations)")
		dryRun      = flag.Bool("dry-run", false, "Preview changes without applying")
		apply       = flag.Bool("apply", false, "Apply refactoring changes to file")

		explain    = flag.Bool("explain", false, "Explain the code in natural language")
		search     = flag.String("search", "", "Search type: procedure, function, variable, call, unused")
		searchName = flag.String("search-name", "", "Name to search for")
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

	// SEARCH
	if *search != "" {
		an := analyzer.NewAnalyzer()
		analysisResult := an.Analyze(program)
		var searchResult refactor.SearchResult

		switch *search {
		case "procedure":
			if *searchName == "" {
				result += "\n❌ Please provide -search-name for procedure search\n"
			} else {
				searchResult = refactor.SearchProcedure(program, *searchName)
			}
		case "function":
			if *searchName == "" {
				result += "\n❌ Please provide -search-name for function search\n"
			} else {
				searchResult = refactor.SearchFunction(program, *searchName)
			}
		case "variable":
			if *searchName == "" {
				result += "\n❌ Please provide -search-name for variable search\n"
			} else {
				searchResult = refactor.SearchVariable(program, *searchName)
			}
		case "call":
			if *searchName == "" {
				result += "\n❌ Please provide -search-name for call search\n"
			} else {
				searchResult = refactor.SearchCall(program, *searchName)
			}
		case "unused":
			searchResult = refactor.SearchUnused(program, analysisResult)
		default:
			result += fmt.Sprintf("\n❌ Unknown search type: %s\n", *search)
			result += "Available search types:\n"
			result += "  procedure   - Search for procedure\n"
			result += "  function    - Search for function\n"
			result += "  variable    - Search for variable\n"
			result += "  call        - Search for CALL statements\n"
			result += "  unused      - Find all unused code\n"
		}

		if searchResult.TotalCount > 0 || searchResult.Message != "" {
			result += refactor.FormatSearchResult(searchResult)
		}
	}
	// REFACTORING
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

		case "rename":
			if *refactorOld == "" || *refactorNew == "" {
				result += "\n❌ Please provide both -old and -new names\n"
				result += "Usage: esql-ast -f file.esql -refactor rename -old <oldName> -new <newName>\n"
			} else {
				// Try variable rename first, then procedure, then function
				renameResult := refactor.RenameVariable(program, *refactorOld, *refactorNew, *dryRun)
				if !renameResult.Success {
					renameResult = refactor.RenameProcedure(program, *refactorOld, *refactorNew, *dryRun)
				}
				if !renameResult.Success {
					renameResult = refactor.RenameFunction(program, *refactorOld, *refactorNew, *dryRun)
				}
				result += refactor.FormatRenameResult(renameResult, *dryRun)

				// Apply changes if -apply is set and not dry-run
				if *apply && !*dryRun && renameResult.Success {
					newContent := refactor.ApplyRenameChanges(input, renameResult)
					if *output != "" {
						err := os.WriteFile(*output, []byte(newContent), 0644)
						if err != nil {
							result += fmt.Sprintf("\n❌ Error writing to output file: %v\n", err)
						} else {
							result += fmt.Sprintf("\n✅ Changes saved to: %s\n", *output)
						}
					} else if *file != "" {
						// Backup original
						backupFile := *file + ".bak"
						err := os.WriteFile(backupFile, []byte(input), 0644)
						if err != nil {
							result += fmt.Sprintf("\n⚠️ Could not create backup: %v\n", err)
						} else {
							result += fmt.Sprintf("\n📁 Backup saved to: %s\n", backupFile)
						}

						// Write changes
						err = os.WriteFile(*file, []byte(newContent), 0644)
						if err != nil {
							result += fmt.Sprintf("\n❌ Error writing to file: %v\n", err)
						} else {
							result += fmt.Sprintf("\n✅ File updated: %s\n", *file)
						}
					}
				} else if *apply && renameResult.Success {
					result += "\n💡 Dry-run mode: changes not applied. Remove -dry-run to apply.\n"
				}
			}

		default:
			result += fmt.Sprintf("\n❌ Unknown refactor operation: %s\n", *refactorCmd)
			result += "Available operations:\n"
			result += "  suggest     - Show refactoring suggestions\n"
			result += "  dead-code   - Show dead code analysis\n"
			result += "  rename      - Rename variable/procedure/function\n"
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

	if *explain {
		an := analyzer.NewAnalyzer()
		analysisResult := an.Analyze(program)
		explanationResult := refactor.Explain(program, analysisResult)
		result += refactor.FormatExplanation(explanationResult)
	}

	if *output != "" && *refactorCmd != "rename" {
		err := os.WriteFile(*output, []byte(result), 0644)
		if err != nil {
			fmt.Printf("Error writing output: %v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Println(result)
	}
}
