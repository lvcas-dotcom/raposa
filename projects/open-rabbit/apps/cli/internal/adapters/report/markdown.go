package report

import (
	"fmt"
	"io"
	"strings"

	"open-rabbit.local/cli/internal/application"
)

func text(v string) string {
	r := strings.NewReplacer("<", "&lt;", ">", "&gt;", "`", "&#96;", "[", "&#91;", "]", "&#93;")
	return r.Replace(v)
}

func Markdown(w io.Writer, r application.Relatorio) error {
	var b strings.Builder
	fmt.Fprintf(&b, "# Open Rabbit\n\nBase: `%s`\n\nHead: `%s`\n\nModo: %s · duração: %d ms\n\n%d achado(s), %d candidato(s) descartado(s). Resultado informativo; ausência de achados não representa aprovação.\n", r.Mudanca.Base, r.Mudanca.Head, r.Modo, r.DuracaoMS, len(r.Achados), len(r.Descartados))
	for _, a := range r.Achados {
		fmt.Fprintf(&b, "\n## [%s] %s\n\n%s:%d · confiança %.2f · ID %s\n\nCondição: %s\n\nImpacto: %s\n\nEsperado: %s\n\nVerificação: %s\n", a.Severidade, text(a.Titulo), text(a.Caminho), a.Linha, a.Confianca, a.ID, text(a.Condicao), text(a.Impacto), text(a.Esperado), text(a.Contraevidencia))
		for _, e := range a.Evidencias {
			fmt.Fprintf(&b, "\nEvidência em %s:%d:\n\n", text(e.Caminho), e.Linha)
			for _, line := range strings.Split(e.Trecho, "\n") {
				fmt.Fprintf(&b, "    %s\n", line)
			}
		}
	}
	if len(r.Descartados) > 0 {
		fmt.Fprintln(&b, "\n## Candidatos descartados")
		for _, d := range r.Descartados {
			fmt.Fprintf(&b, "\n- %s:%d: %s\n", text(d.Achado.Caminho), d.Achado.Linha, text(d.Motivo))
		}
	}
	fmt.Fprintln(&b, "\n## Limitações")
	for _, l := range r.Limitacoes {
		fmt.Fprintf(&b, "\n- %s\n", text(l))
	}
	_, err := io.WriteString(w, b.String())
	return err
}
