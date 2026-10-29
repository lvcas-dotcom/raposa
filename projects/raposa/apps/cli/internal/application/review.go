package application

import (
	"context"
	"fmt"
	"sort"
	"time"

	"raposa.local/cli/internal/domain"
)

type Repository interface {
	Snapshot(context.Context, string, string) (domain.Mudanca, error)
}

type Analyzer interface {
	Analyze(context.Context, domain.Mudanca) ([]domain.Achado, []string, error)
}

type Reviewer interface {
	Candidates(context.Context, domain.Mudanca) ([]domain.Achado, error)
	Verify(context.Context, domain.Mudanca, domain.Achado) (Veredito, error)
}

type Veredito struct {
	Confirmado bool
	Motivo     string
}

type Descarte struct {
	Achado domain.Achado
	Motivo string
}

type Relatorio struct {
	Versao          string
	Mudanca         domain.Mudanca
	Achados         []domain.Achado
	Descartados     []Descarte
	Modo            string
	DuracaoMS       int64
	PrazoProvedorMS int64 `json:",omitempty"`
	Limitacoes      []string
}

type Service struct {
	Repository Repository
	Analyzer   Analyzer
	Reviewer   Reviewer
}

func (s Service) Review(ctx context.Context, base, head string) (Relatorio, error) {
	start := time.Now()
	m, err := s.Repository.Snapshot(ctx, base, head)
	if err != nil {
		return Relatorio{}, err
	}
	r := Relatorio{Versao: "1", Mudanca: m, Modo: "static", Achados: []domain.Achado{}, Descartados: []Descarte{}, Limitacoes: append([]string{}, m.Limitacoes...)}
	candidates, limits, err := s.Analyzer.Analyze(ctx, m)
	if err != nil {
		return Relatorio{}, fmt.Errorf("análise estática: %w", err)
	}
	r.Limitacoes = append(r.Limitacoes, limits...)
	if s.Reviewer != nil {
		r.Modo = "static+semantic"
		semantic, err := s.Reviewer.Candidates(ctx, m)
		if err != nil {
			return Relatorio{}, fmt.Errorf("geração de candidatos: %w", err)
		}
		candidates = append(candidates, semantic...)
	} else {
		r.Limitacoes = append(r.Limitacoes, "Provedor semântico não configurado; apenas baseline estática Go, sem revisão geral de lógica.")
	}
	seen := map[string]bool{}
	for _, a := range candidates {
		if err := a.Validate(m); err != nil {
			r.Descartados = append(r.Descartados, Descarte{a, err.Error()})
			continue
		}
		a.ID = a.Fingerprint()
		if seen[a.ID] {
			r.Descartados = append(r.Descartados, Descarte{a, "candidato duplicado"})
			continue
		}
		if s.Reviewer != nil {
			v, err := s.Reviewer.Verify(ctx, m, a)
			if err != nil {
				return Relatorio{}, fmt.Errorf("verificação do candidato: %w", err)
			}
			if !v.Confirmado {
				r.Descartados = append(r.Descartados, Descarte{a, v.Motivo})
				continue
			}
			a.Contraevidencia = v.Motivo
		}
		seen[a.ID] = true
		r.Achados = append(r.Achados, a)
	}
	priority := map[string]int{"critical": 0, "high": 1, "medium": 2, "low": 3}
	sort.SliceStable(r.Achados, func(i, j int) bool {
		a, b := r.Achados[i], r.Achados[j]
		if a.Severidade != b.Severidade {
			return priority[a.Severidade] < priority[b.Severidade]
		}
		if a.Caminho != b.Caminho {
			return a.Caminho < b.Caminho
		}
		return a.Linha < b.Linha
	})
	r.DuracaoMS = time.Since(start).Milliseconds()
	return r, nil
}
