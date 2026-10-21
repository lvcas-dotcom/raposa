package domain

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestNucleoNaoImportaInfraestrutura(t *testing.T) {
	root := "../.."
	err := filepath.WalkDir(root, func(p string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(p, ".go") || strings.HasSuffix(p, "_test.go") {
			return nil
		}
		f, err := parser.ParseFile(token.NewFileSet(), p, nil, 0)
		if err != nil {
			return err
		}
		for _, violation := range violations(p, f) {
			t.Error(violation)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

func violations(path string, file *ast.File) []string {
	result := []string{}
	for _, imp := range file.Imports {
		value, _ := strconv.Unquote(imp.Path.Value)
		infra := strings.Contains(value, ".") || value == "os" || strings.HasPrefix(value, "os/") || value == "net" || strings.HasPrefix(value, "net/") || strings.HasPrefix(value, "database/") || value == "syscall" || value == "unsafe"
		if strings.Contains(path, "/domain/") && infra {
			result = append(result, "domain importa infraestrutura: "+path+": "+value)
		}
		if strings.Contains(path, "/application/") && strings.Contains(value, "/adapters/") {
			result = append(result, "application importa adapter: "+path)
		}
	}
	if strings.Contains(path, "/domain/") {
		ast.Inspect(file, func(node ast.Node) bool {
			if field, ok := node.(*ast.Field); ok && field.Tag != nil {
				result = append(result, "domain contém tag de infraestrutura: "+path)
			}
			return true
		})
	}
	return result
}

func TestGuardasDeArquiteturaRecusamViolacoes(t *testing.T) {
	tests := []struct{ name, path, source string }{
		{"http no domínio", "internal/domain/a.go", "package domain\nimport \"net/http\""},
		{"terceiro no domínio", "internal/domain/a.go", "package domain\nimport \"example.org/sdk\""},
		{"adapter no application", "internal/application/a.go", "package application\nimport \"open-rabbit.local/cli/internal/adapters/git\""},
		{"tag no domínio", "internal/domain/a.go", "package domain\ntype Entidade struct { Nome string `json:\"nome\"` }"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := parser.ParseFile(token.NewFileSet(), tt.path, tt.source, 0)
			if err != nil {
				t.Fatal(err)
			}
			if len(violations(tt.path, file)) == 0 {
				t.Fatal("guarda deveria reprovar o defeito injetado")
			}
		})
	}
}
