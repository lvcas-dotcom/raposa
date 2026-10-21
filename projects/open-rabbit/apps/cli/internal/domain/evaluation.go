package domain

import "fmt"

type Julgamento struct {
	ID      string
	Defeito string
	Valido  bool
}

type RevisaoJulgada struct {
	Revisor string
	Achados []Julgamento
}

type Caso struct {
	ID       string
	Base     string
	Head     string
	Defeitos []string
	Revisoes []RevisaoJulgada
}

type Metrica struct {
	Revisor              string
	Casos                int
	VerdadeirosPositivos int
	FalsosPositivos      int
	FalsosNegativos      int
	Precisao             *float64
	Recall               *float64
}

func ratio(a, b int) *float64 {
	if b == 0 {
		return nil
	}
	v := float64(a) / float64(b)
	return &v
}

func Avaliar(cases []Caso) ([]Metrica, error) {
	if len(cases) == 0 {
		return nil, fmt.Errorf("avaliação sem casos")
	}
	metrics := []Metrica{}
	indexes := map[string]int{}
	caseIDs := map[string]bool{}
	for _, c := range cases {
		if c.ID == "" || caseIDs[c.ID] || !sha(c.Base) || !sha(c.Head) {
			return nil, fmt.Errorf("caso sem identidade ou SHAs completos")
		}
		caseIDs[c.ID] = true
		defects := map[string]bool{}
		for _, d := range c.Defeitos {
			if d == "" || defects[d] {
				return nil, fmt.Errorf("defeito duplicado ou vazio")
			}
			defects[d] = true
		}
		if len(c.Revisoes) == 0 {
			return nil, fmt.Errorf("caso sem revisões adjudicadas")
		}
		reviewers := map[string]bool{}
		for _, review := range c.Revisoes {
			if review.Revisor == "" || reviewers[review.Revisor] {
				return nil, fmt.Errorf("revisor duplicado ou vazio no caso")
			}
			reviewers[review.Revisor] = true
			i, ok := indexes[review.Revisor]
			if !ok {
				i = len(metrics)
				indexes[review.Revisor] = i
				metrics = append(metrics, Metrica{Revisor: review.Revisor})
			}
			m := &metrics[i]
			m.Casos++
			found := map[string]bool{}
			ids := map[string]bool{}
			for _, a := range review.Achados {
				if a.ID == "" || ids[a.ID] {
					return nil, fmt.Errorf("achado sem ID único")
				}
				ids[a.ID] = true
				if a.Valido {
					if !defects[a.Defeito] {
						return nil, fmt.Errorf("achado válido sem defeito conhecido")
					}
					if found[a.Defeito] {
						m.FalsosPositivos++
					} else {
						m.VerdadeirosPositivos++
						found[a.Defeito] = true
					}
				} else {
					if a.Defeito != "" {
						return nil, fmt.Errorf("achado inválido não pode mapear defeito")
					}
					m.FalsosPositivos++
				}
			}
			m.FalsosNegativos += len(defects) - len(found)
		}
	}
	for i := range metrics {
		m := &metrics[i]
		if m.Casos != len(cases) {
			return nil, fmt.Errorf("comparação exige os mesmos casos para todos os revisores")
		}
		m.Precisao = ratio(m.VerdadeirosPositivos, m.VerdadeirosPositivos+m.FalsosPositivos)
		m.Recall = ratio(m.VerdadeirosPositivos, m.VerdadeirosPositivos+m.FalsosNegativos)
	}
	return metrics, nil
}

func sha(v string) bool {
	if len(v) != 40 && len(v) != 64 {
		return false
	}
	for _, c := range v {
		if !(c >= '0' && c <= '9' || c >= 'a' && c <= 'f') {
			return false
		}
	}
	return true
}
