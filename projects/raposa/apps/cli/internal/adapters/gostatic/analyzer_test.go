package gostatic

import (
	"context"
	"strings"
	"testing"

	"raposa.local/cli/internal/domain"
)

func TestIdentificaTLSComAliasESomenteNoDiff(t *testing.T) {
	source := "package demo\nimport secure \"crypto/tls\"\nvar c = secure.Config{\n InsecureSkipVerify: true,\n}\n"
	for _, line := range []int{1, 4} {
		m := domain.Mudanca{Arquivos: []domain.Arquivo{{Caminho: "a.go", Depois: source, Alteradas: []domain.Intervalo{{Inicio: line, Fim: line}}}}}
		a, _, err := (Analyzer{}).Analyze(context.Background(), m)
		if err != nil {
			t.Fatal(err)
		}
		if line == 1 && len(a) != 0 {
			t.Fatal("achado preexistente reportado")
		}
		if line == 4 {
			if len(a) != 1 {
				t.Fatal("TLS não identificado")
			}
			if err := a[0].Validate(m); err != nil {
				t.Fatal(err)
			}
		}
	}
}

func TestBaselineOmiteCallbackCustomizadoSemPresumirSuaCorrecao(t *testing.T) {
	for _, callback := range []string{"verifyPeer", "nil"} {
		t.Run(callback, func(t *testing.T) {
			source := "package demo\nimport \"crypto/tls\"\nvar c = tls.Config{\n InsecureSkipVerify: true,\n VerifyConnection: " + callback + ",\n}\n"
			m := domain.Mudanca{Arquivos: []domain.Arquivo{{Caminho: "a.go", Depois: source, Alteradas: []domain.Intervalo{{Inicio: 4, Fim: 5}}}}}
			a, limits, err := (Analyzer{}).Analyze(context.Background(), m)
			if err != nil {
				t.Fatal(err)
			}
			if callback == "verifyPeer" && (len(a) != 0 || !strings.Contains(strings.Join(limits, " "), "callback TLS")) {
				t.Fatal("callback deveria exigir investigação semântica")
			}
			if callback == "nil" && len(a) != 1 {
				t.Fatal("nil não implementa verificação customizada")
			}
		})
	}
}

func TestNaoConfundeNomeDeCampoDeOutroTipo(t *testing.T) {
	m := domain.Mudanca{Arquivos: []domain.Arquivo{{Caminho: "a.go", Depois: "package a\nvar a = Outro{InsecureSkipVerify: true}\n", Alteradas: []domain.Intervalo{{Inicio: 2, Fim: 2}}}}}
	a, _, err := (Analyzer{}).Analyze(context.Background(), m)
	if err != nil || len(a) != 0 {
		t.Fatalf("%v %v", a, err)
	}
}
