package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/lyeslabs/mcpgen/internal/generator"
)

func main() {

	// Define command-line flags
	inputFile := flag.String("input", "", "Path to the OpenAPI specification file (JSON or YAML)")
	outputDir := flag.String("output", "", "Path to the output MCP server directory")

	validation := flag.Bool("validation", false, "Enable OpenAPI validation")
	packageName := flag.String("package", "mcpgen", "Generated package name")
	includes := flag.String("includes", "", "Comma-separated list of includes for the generated code")

	// Parse command-line flags
	flag.Parse()

	// Validate required flags
	if *inputFile == "" {
		fmt.Println("Error: input file is required")
		flag.Usage()
		os.Exit(1)
	}

	if *outputDir == "" {
		fmt.Println("Error: output file is required")
		flag.Usage()
		os.Exit(1)
	}

	// Create the output directory if it doesn't exist
	if *outputDir != "" && *outputDir != "." {
		err := os.MkdirAll(*outputDir, 0755)
		if err != nil {
			fmt.Printf("Error creating output directory '%s': %v\n", *outputDir, err)
			os.Exit(1)
		}
	}

	generator, err := generator.NewGenerator(*inputFile, *validation, *packageName, *outputDir)
	if err != nil {
		fmt.Printf("Error creating generator: %v\n", err)
		os.Exit(1)
	}

	// Generate the HTTP CLIENT
	if *includes != "" {
		err = generator.GenerateHTTPClient(strings.Split(*includes, ","))
		if err != nil {
			fmt.Printf("Error generating HTTP client: %v\n", err)
			os.Exit(1)
		}
	}

	// Generate the MCP configuration
	err = generator.GenerateMCP()
	if err != nil {
		fmt.Printf("Error generating MCP configuration: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully converted OpenAPI specification to MCP configuration: %s\n", *outputDir)
}
