package domain

import (
	"strings"
	"testing"
)

func TestAvaliaCoberturaRuidoEDuplicacao(t *testing.T) {
	c := Caso{ID: "caso", Base: strings.Repeat("a", 40), Head: strings.Repeat("b", 40), Defeitos: []string{"d1", "d2"}, Revisoes: []RevisaoJulgada{{Revisor: "raposa", Achados: []Julgamento{{ID: "a1", Defeito: "d1", Valido: true}, {ID: "a2", Defeito: "d1", Valido: true}, {ID: "a3"}}}}}
	metrics, err := Avaliar([]Caso{c})
	if err != nil {
		t.Fatal(err)
	}
	m := metrics[0]
	if m.VerdadeirosPositivos != 1 || m.FalsosPositivos != 2 || m.FalsosNegativos != 1 || *m.Recall != .5 || *m.Precisao != 1.0/3 {
		t.Fatalf("métrica incorreta: %+v", m)
	}
}

func TestNaoInventaPrecisaoSemAchados(t *testing.T) {
	c := Caso{ID: "vazio", Base: strings.Repeat("a", 40), Head: strings.Repeat("b", 40), Revisoes: []RevisaoJulgada{{Revisor: "r"}}}
	m, err := Avaliar([]Caso{c})
	if err != nil {
		t.Fatal(err)
	}
	if m[0].Precisao != nil || m[0].Recall != nil {
		t.Fatal("métrica indefinida deve ser null")
	}
}

func TestRecusaComparacaoComConjuntosDiferentes(t *testing.T) {
	a := Caso{ID: "a", Base: strings.Repeat("a", 40), Head: strings.Repeat("b", 40), Revisoes: []RevisaoJulgada{{Revisor: "r1"}, {Revisor: "r2"}}}
	b := a
	b.ID = "b"
	b.Revisoes = []RevisaoJulgada{{Revisor: "r1"}}
	if _, err := Avaliar([]Caso{a, b}); err == nil {
		t.Fatal("esperava recusa")
	}
}
