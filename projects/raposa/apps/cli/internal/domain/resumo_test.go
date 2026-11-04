package domain

import (
	"strings"
	"testing"
)

func resumoFixture() (Mudanca, Resumo) {
	m := Mudanca{Arquivos: []Arquivo{
		{Caminho: "a.go", Alteradas: []Intervalo{{Inicio: 3, Fim: 8}}},
		{Caminho: "b.go"},
	}}
	r := Resumo{Resumo: "Muda A e remove B", Arquivos: []ResumoArquivo{
		{Caminho: "a.go", Linha: "ajusta A", Ancora: &Intervalo{Inicio: 4, Fim: 6}},
		{Caminho: "b.go", Linha: "removido"},
	}}
	return m, r
}

func TestResumoValido(t *testing.T) {
	m, r := resumoFixture()
	if err := r.Validate(m); err != nil {
		t.Fatal(err)
	}
}

func TestResumoRecusadoPorRegra(t *testing.T) {
	cases := map[string]func(*Resumo){
		"vazio":                func(r *Resumo) { r.Resumo = " " },
		"acima do limite":      func(r *Resumo) { r.Resumo = strings.Repeat("é", LimiteResumo+1) },
		"linha longa":          func(r *Resumo) { r.Arquivos[0].Linha = strings.Repeat("x", LimiteLinhaResumo+1) },
		"linha vazia":          func(r *Resumo) { r.Arquivos[0].Linha = "" },
		"quebra de linha":      func(r *Resumo) { r.Arquivos[0].Linha = "a\nb" },
		"caminho desconhecido": func(r *Resumo) { r.Arquivos[0].Caminho = "z.go" },
		"caminho repetido":     func(r *Resumo) { r.Arquivos[1] = r.Arquivos[0] },
		"sem cobertura":        func(r *Resumo) { r.Arquivos = r.Arquivos[:1] },
		"âncora fora":          func(r *Resumo) { r.Arquivos[0].Ancora = &Intervalo{Inicio: 1, Fim: 4} },
		"âncora invertida":     func(r *Resumo) { r.Arquivos[0].Ancora = &Intervalo{Inicio: 6, Fim: 4} },
		"sem âncora":           func(r *Resumo) { r.Arquivos[0].Ancora = nil },
		"âncora sem alteradas": func(r *Resumo) { r.Arquivos[1].Ancora = &Intervalo{Inicio: 1, Fim: 1} },
	}
	for name, mutate := range cases {
		t.Run(name, func(t *testing.T) {
			m, r := resumoFixture()
			mutate(&r)
			if r.Validate(m) == nil {
				t.Fatal("esperava recusa")
			}
		})
	}
}

func TestResumoNaoServeComoEvidencia(t *testing.T) {
	m := Mudanca{Arquivos: []Arquivo{{Caminho: "a.go", Depois: "x", Alteradas: []Intervalo{{Inicio: 1, Fim: 1}}}}}
	a := Achado{Regra: "r", Caminho: "a.go", Linha: 1, Titulo: "t", Severidade: "high", Confianca: .5, Condicao: "c", Impacto: "i", Esperado: "e", Contraevidencia: "v",
		Evidencias: []Evidencia{{Caminho: "resumo", Linha: 1, Trecho: "Muda A"}}}
	if a.Validate(m) == nil {
		t.Fatal("evidência fora do snapshot deveria ser recusada")
	}
}
