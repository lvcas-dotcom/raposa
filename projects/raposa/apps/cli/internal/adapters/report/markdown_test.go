package report

import (
	"bytes"
	"strings"
	"testing"

	"raposa.local/cli/internal/application"
	"raposa.local/cli/internal/domain"
)

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

func TestPrazoDoProvedorSoApareceQuandoConfigurado(t *testing.T) {
	var static, semantic bytes.Buffer
	if err := Markdown(&static, application.Relatorio{Modo: "static"}); err != nil {
		t.Fatal(err)
	}
	if err := Markdown(&semantic, application.Relatorio{Modo: "static+semantic", PrazoProvedorMS: 90000}); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(static.String(), "prazo por chamada") || !strings.Contains(semantic.String(), "· prazo por chamada: 90000 ms\n") {
		t.Fatalf("prazo do provedor incorreto:\n%s\n%s", static.String(), semantic.String())
	}
}
