package gostatic

import (
	"context"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path"
	"strconv"
	"strings"

	"open-rabbit.local/cli/internal/domain"
)

type Analyzer struct{}

func (Analyzer) Analyze(ctx context.Context, m domain.Mudanca) ([]domain.Achado, []string, error) {
	results := []domain.Achado{}
	limits := []string{"Baseline estática: apenas tls.Config com InsecureSkipVerify literal true; sem resolução de tipos, fluxo ou prova de explorabilidade."}
	for _, f := range m.Arquivos {
		if err := ctx.Err(); err != nil {
			return nil, nil, err
		}
		if path.Ext(f.Caminho) != ".go" {
			limits = append(limits, "Sem analisador estático para: "+f.Caminho)
			continue
		}
		if strings.HasSuffix(f.Caminho, "_test.go") {
			limits = append(limits, "Baseline TLS não analisa fixtures de teste: "+f.Caminho)
			continue
		}
		fs := token.NewFileSet()
		file, err := parser.ParseFile(fs, f.Caminho, f.Depois, 0)
		if err != nil {
			limits = append(limits, "Go não pôde ser parseado: "+f.Caminho)
			continue
		}
		aliases := map[string]bool{}
		for _, imp := range file.Imports {
			v, _ := strconv.Unquote(imp.Path.Value)
			if v != "crypto/tls" {
				continue
			}
			name := "tls"
			if imp.Name != nil {
				name = imp.Name.Name
			}
			aliases[name] = true
		}
		ast.Inspect(file, func(n ast.Node) bool {
			lit, ok := n.(*ast.CompositeLit)
			if !ok {
				return true
			}
			typ, ok := lit.Type.(*ast.SelectorExpr)
			if !ok || typ.Sel.Name != "Config" {
				return true
			}
			pkg, ok := typ.X.(*ast.Ident)
			if !ok || !aliases[pkg.Name] {
				return true
			}
			if hasCustomVerification(lit) {
				limits = append(limits, "Baseline omitiu configuração com callback TLS customizado: "+f.Caminho)
				return true
			}
			for _, elt := range lit.Elts {
				kv, ok := elt.(*ast.KeyValueExpr)
				if !ok {
					continue
				}
				key, ok := kv.Key.(*ast.Ident)
				if !ok || key.Name != "InsecureSkipVerify" {
					continue
				}
				value, ok := kv.Value.(*ast.Ident)
				if !ok || value.Name != "true" {
					continue
				}
				line := fs.Position(kv.Pos()).Line
				changed := false
				for _, r := range f.Alteradas {
					if line >= r.Inicio && line <= r.Fim {
						changed = true
					}
				}
				if !changed {
					continue
				}
				lines := strings.Split(f.Depois, "\n")
				results = append(results, domain.Achado{
					Regra:           "go/tls-insecure-skip-verify",
					Caminho:         f.Caminho,
					Linha:           line,
					Titulo:          "Verificação padrão do certificado TLS desativada",
					Severidade:      "high",
					Confianca:       0.85,
					Condicao:        "Se esta configuração for usada em uma conexão sem verificação customizada equivalente, certificados não confiáveis serão aceitos.",
					Impacto:         "A conexão pode aceitar um servidor impostor; a explorabilidade depende do uso desta configuração.",
					Esperado:        "Manter a validação padrão ou verificar explicitamente a cadeia e o hostname em callback próprio.",
					Evidencias:      []domain.Evidencia{{Caminho: f.Caminho, Linha: line, Trecho: lines[line-1]}},
					Contraevidencia: fmt.Sprintf("AST identifica crypto/tls.Config e literal true no diff, sem callback no mesmo literal; fluxo e atribuições posteriores em %s ainda exigem investigação.", pkg.Name),
				})
			}
			return true
		})
	}
	return results, limits, nil
}

func hasCustomVerification(lit *ast.CompositeLit) bool {
	for _, elt := range lit.Elts {
		kv, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			continue
		}
		key, ok := kv.Key.(*ast.Ident)
		if !ok || (key.Name != "VerifyConnection" && key.Name != "VerifyPeerCertificate") {
			continue
		}
		value, ok := kv.Value.(*ast.Ident)
		if !ok || value.Name != "nil" {
			return true
		}
	}
	return false
}
