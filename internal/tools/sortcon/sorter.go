package sortcon

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"sort"
	"strconv"
)

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
