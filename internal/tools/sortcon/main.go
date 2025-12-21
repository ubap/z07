package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"sort"
	"strconv"
)

func main() {
	// GOFILE is automatically set by 'go generate'
	fileName := os.Getenv("GOFILE")
	if fileName == "" {
		if len(os.Args) < 2 {
			fmt.Println("Usage: Run via 'go generate' or provide a file path as an argument.")
			os.Exit(1)
		}
		fileName = os.Args[1]
	}

	// 1. Read the file
	src, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	// 2. Process and Sort
	out, err := SortSource(src)
	if err != nil {
		fmt.Printf("Error processing source: %v\n", err)
		os.Exit(1)
	}

	// 3. Write back to the same file
	err = os.WriteFile(fileName, out, 0644)
	if err != nil {
		fmt.Printf("Error writing file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully sorted constants in %s\n", fileName)
}

// SortSource parses Go source code, sorts constant blocks, and returns formatted code.
func SortSource(src []byte) ([]byte, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "", src, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	for _, decl := range node.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		// We only care about 'const' blocks
		if !ok || genDecl.Tok != token.CONST {
			continue
		}

		// Sort the specs within the block
		sort.SliceStable(genDecl.Specs, func(i, j int) bool {
			a := genDecl.Specs[i].(*ast.ValueSpec)
			b := genDecl.Specs[j].(*ast.ValueSpec)
			return getValue(a) < getValue(b)
		})

		// FIX: Reset positions to avoid the "extra whitespace" bug.
		// Setting NamePos to token.NoPos tells go/format to use default spacing.
		for _, spec := range genDecl.Specs {
			if s, ok := spec.(*ast.ValueSpec); ok {
				for _, name := range s.Names {
					name.NamePos = token.NoPos
				}
			}
		}
	}

	// Format the AST back into a byte slice
	var buf bytes.Buffer
	if err := format.Node(&buf, fset, node); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// getValue extracts the numerical value of a constant for sorting.
func getValue(v *ast.ValueSpec) uint64 {
	if len(v.Values) == 0 {
		return 0
	}

	lit, ok := v.Values[0].(*ast.BasicLit)
	if !ok {
		return 0
	}

	valStr := lit.Value
	// Handle hex (0x0A), decimal (10), or octal
	val, err := strconv.ParseUint(valStr, 0, 64)
	if err != nil {
		return 0
	}
	return val
}
