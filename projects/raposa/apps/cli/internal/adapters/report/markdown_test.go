package report

import (
	"bytes"
	"strings"
	"testing"

	"raposa.local/cli/internal/application"
	"raposa.local/cli/internal/domain"
)

func TestResumoApareceAntesDosAchadosEEscapado(t *testing.T) {
	r := application.Relatorio{
		Achados: []domain.Achado{{Titulo: "Falha"}},
		Resumo: &domain.Resumo{Origem: domain.OrigemModelo, Resumo: "<script>x</script>", Arquivos: []domain.ResumoArquivo{
			{Caminho: "a.go", Linha: "ajusta", Ancora: &domain.Intervalo{Inicio: 2, Fim: 4}},
			{Caminho: "b.go", Linha: "removido"},
		}},
	}
	var out bytes.Buffer
	if err := Markdown(&out, r); err != nil {
		t.Fatal(err)
	}
	s := out.String()
	i, j := strings.Index(s, "## Resumo da mudança"), strings.Index(s, "## [")
	if i < 0 || j < 0 || i > j {
		t.Fatalf("seção ausente ou fora de ordem:\n%s", s)
	}
	for _, want := range []string{"Gerado pelo modelo, não verificado", "&lt;script&gt;", "- a.go:2-4: ajusta", "- b.go: removido"} {
		if !strings.Contains(s, want) {
			t.Fatalf("faltou %q:\n%s", want, s)
		}
	}
}

func TestRelatorioSemResumoNaoTemSecao(t *testing.T) {
	var out bytes.Buffer
	if err := Markdown(&out, application.Relatorio{Limitacoes: []string{"Resumo da mudança indisponível: x"}}); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(out.String(), "## Resumo da mudança") || !strings.Contains(out.String(), "indisponível") {
		t.Fatal("ausência não declarada corretamente")
	}
}

func TestRelatorioEscapaHTMLDeDadosNaoConfiaveis(t *testing.T) {
	r := application.Relatorio{Achados: []domain.Achado{{Titulo: "<script>alert(1)</script>", Caminho: "[link](url)", Evidencias: []domain.Evidencia{{Trecho: "```\n<script>"}}}}}
	var out bytes.Buffer
	if err := Markdown(&out, r); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(out.String(), "## [Warning]") || strings.Contains(out.String(), "## [] <script>") || strings.Contains(out.String(), "[link](url)") {
		t.Fatal("conteúdo não escapado")
	}
	if !strings.Contains(out.String(), "&lt;script&gt;") || !strings.Contains(out.String(), "    <script>") {
		t.Fatal("evidência deveria estar como bloco de código")
	}
}
