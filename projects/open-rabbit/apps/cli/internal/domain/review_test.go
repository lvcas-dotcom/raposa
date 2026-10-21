package domain

import (
	"math"
	"testing"
)

func amostra() (Mudanca, Achado) {
	m := Mudanca{Arquivos: []Arquivo{{Caminho: "a.go", Depois: "package a\nproblema\nfim", Alteradas: []Intervalo{{2, 2}}}}}
	a := Achado{Regra: "regra", Caminho: "a.go", Linha: 2, Titulo: "Falha", Severidade: "high", Confianca: .9, Condicao: "gatilho", Impacto: "efeito", Esperado: "contrato", Contraevidencia: "busca registrada", Evidencias: []Evidencia{{Caminho: "a.go", Linha: 2, Trecho: "problema"}}}
	return m, a
}

func TestAceitaEvidenciaDoSnapshot(t *testing.T) {
	m, a := amostra()
	if err := a.Validate(m); err != nil {
		t.Fatal(err)
	}
}

func TestRecusaAchadosSemRastreabilidade(t *testing.T) {
	tests := []struct {
		name   string
		change func(*Achado)
	}{
		{"linha fora do diff", func(a *Achado) { a.Linha = 3 }},
		{"evidência inventada", func(a *Achado) { a.Evidencias[0].Trecho = "outro" }},
		{"evidência deslocada", func(a *Achado) { a.Evidencias[0].Linha = 1 }},
		{"evidência fora do contexto", func(a *Achado) { a.Evidencias[0].Caminho = "segredo" }},
		{"sem evidência", func(a *Achado) { a.Evidencias = nil }},
		{"evidência só com espaços", func(a *Achado) { a.Evidencias[0].Trecho = " " }},
		{"justificativa só com espaços", func(a *Achado) { a.Impacto = " " }},
		{"sem contraevidência", func(a *Achado) { a.Contraevidencia = "" }},
		{"confiança NaN", func(a *Achado) { a.Confianca = math.NaN() }},
		{"severidade inventada", func(a *Achado) { a.Severidade = "extreme" }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, a := amostra()
			tt.change(&a)
			if a.Validate(m) == nil {
				t.Fatal("esperava recusa")
			}
		})
	}
}

func TestEvidenciaPodeConterMultiplasLinhas(t *testing.T) {
	m, a := amostra()
	a.Evidencias[0].Trecho = "problema\nfim"
	if err := a.Validate(m); err != nil {
		t.Fatal(err)
	}
}
