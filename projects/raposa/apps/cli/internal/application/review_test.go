package application

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"raposa.local/cli/internal/domain"
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

type summarizerStub struct {
	*reviewerStub
	resumo domain.Resumo
	err    error
	calls  int
}

func (s *summarizerStub) Summarize(context.Context, domain.Mudanca) (domain.Resumo, error) {
	s.calls++
	return s.resumo, s.err
}

func resumoDe(m domain.Mudanca) domain.Resumo {
	return domain.Resumo{Origem: "forjada", Resumo: "Muda a.go", Arquivos: []domain.ResumoArquivo{{Caminho: "a.go", Linha: "ajusta", Ancora: &domain.Intervalo{Inicio: 1, Fim: 1}}}}
}

func TestResumoNaoAlteraAchadosEDeclaraAusencia(t *testing.T) {
	m, a := fixture()
	run := func(r Reviewer) Relatorio {
		out, err := Service{Repository: repoStub{m}, Analyzer: analyzerStub{[]domain.Achado{a}}, Reviewer: r}.Review(context.Background(), "b", "h")
		if err != nil {
			t.Fatal(err)
		}
		return out
	}
	sem := run(&reviewerStub{confirm: true})
	com := run(&summarizerStub{reviewerStub: &reviewerStub{confirm: true}, resumo: resumoDe(m)})
	if len(sem.Achados) != 1 || len(com.Achados) != 1 || sem.Achados[0].ID != com.Achados[0].ID || len(sem.Descartados) != len(com.Descartados) {
		t.Fatal("resumo alterou os achados")
	}
	if sem.Resumo != nil || !hasLimit(sem, "provedor sem a fase summary") {
		t.Fatalf("ausência não declarada: %+v", sem.Limitacoes)
	}
	if com.Resumo == nil || com.Resumo.Origem != domain.OrigemModelo {
		t.Fatalf("resumo válido ausente ou sem origem fixada: %+v", com.Resumo)
	}
}

func TestResumoInvalidoOuComFalhaEDescartado(t *testing.T) {
	m, a := fixture()
	ruim := resumoDe(m)
	ruim.Arquivos[0].Ancora = &domain.Intervalo{Inicio: 9, Fim: 9}
	for name, s := range map[string]*summarizerStub{
		"invalido": {reviewerStub: &reviewerStub{confirm: true}, resumo: ruim},
		"falha":    {reviewerStub: &reviewerStub{confirm: true}, err: errors.New("sem suporte")},
	} {
		t.Run(name, func(t *testing.T) {
			out, err := Service{Repository: repoStub{m}, Analyzer: analyzerStub{[]domain.Achado{a}}, Reviewer: s}.Review(context.Background(), "b", "h")
			if err != nil || len(out.Achados) != 1 || out.Resumo != nil || !hasLimit(out, "Resumo da mudança indisponível") {
				t.Fatalf("%v %+v", err, out)
			}
		})
	}
}

func TestSemProvedorNaoChamaResumo(t *testing.T) {
	m, a := fixture()
	out, err := Service{Repository: repoStub{m}, Analyzer: analyzerStub{[]domain.Achado{a}}}.Review(context.Background(), "b", "h")
	if err != nil || out.Resumo != nil || hasLimit(out, "Resumo da mudança") {
		t.Fatalf("%v %+v", err, out.Limitacoes)
	}
}

func hasLimit(r Relatorio, part string) bool {
	for _, l := range r.Limitacoes {
		if strings.Contains(l, part) {
			return true
		}
	}
	return false
}

func TestResumoNoJSONSoComOrigemEOmitidoQuandoAusente(t *testing.T) {
	sem, _ := json.Marshal(Relatorio{})
	if strings.Contains(string(sem), `"Resumo"`) {
		t.Fatalf("campo deveria estar ausente: %s", sem)
	}
	com, _ := json.Marshal(Relatorio{Resumo: &domain.Resumo{Origem: domain.OrigemModelo, Resumo: "x"}})
	if !strings.Contains(string(com), `"Origem":"modelo"`) {
		t.Fatalf("origem ausente: %s", com)
	}
}
