package domain

import (
	"strings"
	"testing"
)

func datasetValido() Dataset {
	return Dataset{Versao: "1", Repositorio: "https://example.invalid/r", Licenca: "MIT", Casos: []CasoBenchmark{
		{ID: "c1", Linguagem: "go", Base: strings.Repeat("a", 40), Head: strings.Repeat("b", 40), Particao: "desenvolvimento", Defeitos: []DefeitoConhecido{{ID: "d1", Caminho: "pkg/a.go", Inicio: 3, Fim: 4, Classe: "contrato", Descricao: "retorno incompatível", Origem: "humano"}}},
		{ID: "c2", Linguagem: "go", Base: strings.Repeat("c", 40), Head: strings.Repeat("d", 40), Particao: "reservado", Limpo: true},
	}}
}

func TestDatasetValidoGeraResumo(t *testing.T) {
	r, err := ValidarDataset(datasetValido())
	if err != nil {
		t.Fatal(err)
	}
	if r != (ResumoDataset{Casos: 2, Defeitos: 1, Limpos: 1, Reservados: 1}) {
		t.Fatalf("resumo incorreto: %+v", r)
	}
}

func TestDatasetInvalidoApontaCampo(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*Dataset)
		campo  string
	}{
		{"versão", func(d *Dataset) { d.Versao = "2" }, "Versao"},
		{"repositório sem https", func(d *Dataset) { d.Repositorio = "git@host:r" }, "Repositorio"},
		{"licença", func(d *Dataset) { d.Licenca = " " }, "Licenca"},
		{"sem casos", func(d *Dataset) { d.Casos = nil }, "Casos:"},
		{"caso repetido", func(d *Dataset) { d.Casos[1].ID = "c1" }, "Casos[1].ID"},
		{"linguagem", func(d *Dataset) { d.Casos[0].Linguagem = "" }, "Casos[0].Linguagem"},
		{"SHA curto", func(d *Dataset) { d.Casos[0].Base = "abc123" }, "Casos[0].Base"},
		{"head igual à base", func(d *Dataset) { d.Casos[0].Head = d.Casos[0].Base }, "Casos[0].Head"},
		{"partição", func(d *Dataset) { d.Casos[0].Particao = "teste" }, "Casos[0].Particao"},
		{"limpo com defeito", func(d *Dataset) { d.Casos[0].Limpo = true }, "Casos[0].Limpo"},
		{"sem defeito e não limpo", func(d *Dataset) { d.Casos[1].Limpo = false }, "Casos[1].Defeitos"},
		{"defeito repetido", func(d *Dataset) { d.Casos[0].Defeitos = append(d.Casos[0].Defeitos, d.Casos[0].Defeitos[0]) }, "Casos[0].Defeitos[1].ID"},
		{"caminho absoluto", func(d *Dataset) { d.Casos[0].Defeitos[0].Caminho = "/etc/a.go" }, "Defeitos[0].Caminho"},
		{"caminho escapa da raiz", func(d *Dataset) { d.Casos[0].Defeitos[0].Caminho = "../a.go" }, "Defeitos[0].Caminho"},
		{"caminho não normalizado", func(d *Dataset) { d.Casos[0].Defeitos[0].Caminho = "pkg//a.go" }, "Defeitos[0].Caminho"},
		{"intervalo invertido", func(d *Dataset) { d.Casos[0].Defeitos[0].Fim = 2 }, "Defeitos[0].Fim"},
		{"linha zero", func(d *Dataset) { d.Casos[0].Defeitos[0].Inicio = 0 }, "Defeitos[0].Fim"},
		{"classe", func(d *Dataset) { d.Casos[0].Defeitos[0].Classe = "estilo" }, "Defeitos[0].Classe"},
		{"descrição", func(d *Dataset) { d.Casos[0].Defeitos[0].Descricao = "" }, "Defeitos[0].Descricao"},
		{"origem", func(d *Dataset) { d.Casos[0].Defeitos[0].Origem = "llm" }, "Defeitos[0].Origem"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := datasetValido()
			tt.mutate(&d)
			_, err := ValidarDataset(d)
			if err == nil || !strings.Contains(err.Error(), tt.campo) {
				t.Fatalf("esperava erro em %s, recebeu %v", tt.campo, err)
			}
		})
	}
}

func TestDatasetInvalidoReportaTodosOsProblemas(t *testing.T) {
	d := datasetValido()
	d.Licenca = ""
	d.Casos[0].Defeitos[0].Classe = "estilo"
	d.Casos[1].Particao = ""
	_, err := ValidarDataset(d)
	if err == nil {
		t.Fatal("esperava recusa")
	}
	for _, campo := range []string{"Licenca", "Casos[0].Defeitos[0].Classe", "Casos[1].Particao"} {
		if !strings.Contains(err.Error(), campo) {
			t.Errorf("faltou %s em %q", campo, err)
		}
	}
}
