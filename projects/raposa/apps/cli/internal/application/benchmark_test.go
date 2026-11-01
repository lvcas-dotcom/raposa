package application

import (
	"context"
	"errors"
	"strings"
	"testing"

	"raposa.local/cli/internal/domain"
)

type commitsStub struct {
	commits map[string]bool
	lines   map[string]int
}

func (c commitsStub) ResolveCommit(_ context.Context, ref string) (string, error) {
	if !c.commits[ref] {
		return "", errors.New("ausente")
	}
	return ref, nil
}

func (c commitsStub) CountLines(_ context.Context, commit, path string) (int, error) {
	n, ok := c.lines[commit+":"+path]
	if !ok {
		return 0, errors.New("ausente")
	}
	return n, nil
}

func datasetComUmDefeito(fim int) (domain.Dataset, string, string) {
	base, head := strings.Repeat("a", 40), strings.Repeat("b", 40)
	d := domain.Dataset{Casos: []domain.CasoBenchmark{{Base: base, Head: head, Defeitos: []domain.DefeitoConhecido{{Caminho: "a.go", Inicio: 1, Fim: fim}}}}}
	return d, base, head
}

func TestVerificaCommitsELinhasNoHead(t *testing.T) {
	d, base, head := datasetComUmDefeito(10)
	r := commitsStub{commits: map[string]bool{base: true, head: true}, lines: map[string]int{head + ":a.go": 10}}
	if err := VerificarNoRepositorio(context.Background(), r, d); err != nil {
		t.Fatal(err)
	}
}

func TestRecusaLinhaAlemDoArquivo(t *testing.T) {
	d, base, head := datasetComUmDefeito(11)
	r := commitsStub{commits: map[string]bool{base: true, head: true}, lines: map[string]int{head + ":a.go": 10}}
	err := VerificarNoRepositorio(context.Background(), r, d)
	if err == nil || !strings.Contains(err.Error(), "Casos[0].Defeitos[0].Fim") {
		t.Fatalf("esperava erro de linha, recebeu %v", err)
	}
}

func TestRecusaArquivoAusenteNoHead(t *testing.T) {
	d, base, head := datasetComUmDefeito(1)
	r := commitsStub{commits: map[string]bool{base: true, head: true}}
	err := VerificarNoRepositorio(context.Background(), r, d)
	if err == nil || !strings.Contains(err.Error(), "Casos[0].Defeitos[0].Caminho") {
		t.Fatalf("esperava erro de caminho, recebeu %v", err)
	}
}

func TestHeadAusenteNaoGeraRuidoDeLinhas(t *testing.T) {
	d, base, _ := datasetComUmDefeito(1)
	r := commitsStub{commits: map[string]bool{base: true}}
	err := VerificarNoRepositorio(context.Background(), r, d)
	if err == nil || !strings.Contains(err.Error(), "Casos[0].Head") || strings.Contains(err.Error(), "Defeitos") {
		t.Fatalf("esperava só o erro do head, recebeu %v", err)
	}
}
