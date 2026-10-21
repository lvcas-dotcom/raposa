package application

import (
	"context"
	"errors"
	"testing"

	"open-rabbit.local/cli/internal/domain"
)

type repoStub struct{ m domain.Mudanca }

func (r repoStub) Snapshot(context.Context, string, string) (domain.Mudanca, error) { return r.m, nil }

type analyzerStub struct{ a []domain.Achado }

func (a analyzerStub) Analyze(context.Context, domain.Mudanca) ([]domain.Achado, []string, error) {
	return a.a, nil, nil
}

type reviewerStub struct {
	confirm bool
	err     error
	calls   int
}

func (r *reviewerStub) Candidates(context.Context, domain.Mudanca) ([]domain.Achado, error) {
	return nil, r.err
}
func (r *reviewerStub) Verify(context.Context, domain.Mudanca, domain.Achado) (Veredito, error) {
	r.calls++
	return Veredito{r.confirm, "proteção verificada"}, r.err
}

func fixture() (domain.Mudanca, domain.Achado) {
	m := domain.Mudanca{Arquivos: []domain.Arquivo{{Caminho: "a.go", Depois: "falha", Alteradas: []domain.Intervalo{{Inicio: 1, Fim: 1}}}}}
	a := domain.Achado{Regra: "r", Caminho: "a.go", Linha: 1, Titulo: "Falha", Severidade: "high", Confianca: .8, Condicao: "c", Impacto: "i", Esperado: "e", Contraevidencia: "v", Evidencias: []domain.Evidencia{{Caminho: "a.go", Linha: 1, Trecho: "falha"}}}
	return m, a
}

func TestPublicaSomenteCandidatosConfirmadosERastreaveis(t *testing.T) {
	for _, confirmed := range []bool{false, true} {
		m, a := fixture()
		bad := a
		bad.Linha = 50
		r := &reviewerStub{confirm: confirmed}
		s := Service{Repository: repoStub{m}, Analyzer: analyzerStub{[]domain.Achado{bad, a, a}}, Reviewer: r}
		out, err := s.Review(context.Background(), "base", "head")
		if err != nil {
			t.Fatal(err)
		}
		want := 0
		if confirmed {
			want = 1
		}
		if len(out.Achados) != want || len(out.Descartados) != 3-want {
			t.Fatalf("relatório incorreto: %+v", out)
		}
		if confirmed && r.calls != 1 {
			t.Fatal("candidato inválido ou duplicado chegou ao verificador")
		}
	}
}

func TestFalhaDoModeloNaoViraRelatorioLimpo(t *testing.T) {
	m, a := fixture()
	s := Service{Repository: repoStub{m}, Analyzer: analyzerStub{[]domain.Achado{a}}, Reviewer: &reviewerStub{err: errors.New("indisponível")}}
	if _, err := s.Review(context.Background(), "base", "head"); err == nil {
		t.Fatal("esperava falha")
	}
}
