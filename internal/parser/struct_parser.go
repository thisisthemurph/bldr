package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
)

type Field struct {
	Name string
	Type string
}

type StructDetail struct {
	Name        string
	PackageName string
	PackageDir  string
	Import      string
	Fields      []Field
}

func ParseStruct(goModule, filePath, structName string) (*StructDetail, error) {
	src, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	fileset := token.NewFileSet()
	f, err := parser.ParseFile(fileset, filePath, src, parser.AllErrors)
	if err != nil {
		return nil, err
	}

	for _, declaration := range f.Decls {
		genDecl, ok := declaration.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}

		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok || typeSpec.Name.Name != structName {
				continue
			}

			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				return nil, fmt.Errorf("%s is not a struct", structName)
			}

			var fields []Field
			for _, field := range structType.Fields.List {
				for _, name := range field.Names {
					fields = append(fields, Field{
						Name: name.Name,
						Type: fmt.Sprintf("%s", fieldTypeAsString(field.Type)),
					})
				}
			}

			return &StructDetail{
				Name:        structName,
				PackageName: f.Name.Name,
				PackageDir:  filepath.Dir(filePath),
				Import:      fmt.Sprintf("%s/%s", goModule, filepath.Dir(filePath)),
				Fields:      fields,
			}, nil
		}
	}

	return nil, fmt.Errorf("struct %s not found in file %s", structName, filePath)
}

func fieldTypeAsString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		return fmt.Sprintf("%s.%s", fieldTypeAsString(t.X), t.Sel.Name)
	case *ast.StarExpr:
		return "*" + fieldTypeAsString(t.X)
	case *ast.ArrayType:
		return "[]" + fieldTypeAsString(t.Elt)
	default:
		return fmt.Sprintf("%T", expr)
	}
}
