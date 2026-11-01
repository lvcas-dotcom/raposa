package main

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"raposa.local/cli/internal/domain"
)

func TestBenchmarkValidateConfereDatasetComGitReal(t *testing.T) {
	dir := t.TempDir()
	commandGit(t, dir, "init", "--initial-branch=main")
	file := filepath.Join(dir, "load.go")
	if err := os.WriteFile(file, []byte("package load\n"), 0600); err != nil {
		t.Fatal(err)
	}
	commandGit(t, dir, "add", ".")
	commandGit(t, dir, "-c", "core.hooksPath=/dev/null", "commit", "-m", "base")
	base := commandGit(t, dir, "rev-parse", "HEAD")
	if err := os.WriteFile(file, []byte("package load\n\nfunc Load() {\n}\n"), 0600); err != nil {
		t.Fatal(err)
	}
	commandGit(t, dir, "-c", "core.hooksPath=/dev/null", "commit", "-am", "head")
	head := commandGit(t, dir, "rev-parse", "HEAD")

	write := func(fim int) string {
		d := domain.Dataset{Versao: "1", Repositorio: "https://example.invalid/load", Licenca: "MIT", Casos: []domain.CasoBenchmark{{
			ID: "c1", Linguagem: "go", Base: base, Head: head, Particao: "desenvolvimento",
			Defeitos: []domain.DefeitoConhecido{{ID: "d1", Caminho: "load.go", Inicio: 3, Fim: fim, Classe: "logica", Descricao: "corpo vazio", Origem: "sintetico"}},
		}}}
		b, err := json.Marshal(d)
		if err != nil {
			t.Fatal(err)
		}
		p := filepath.Join(t.TempDir(), "dataset.json")
		if err := os.WriteFile(p, b, 0600); err != nil {
			t.Fatal(err)
		}
		return p
	}

	var out bytes.Buffer
	if err := run(context.Background(), []string{"benchmark", "validate", "--input", write(4), "--repo", dir}, &out, &out); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), `"VerificadoNoRepositorio": true`) || !strings.Contains(out.String(), `"Defeitos": 1`) {
		t.Fatalf("resumo inesperado: %s", out.String())
	}

	err := run(context.Background(), []string{"benchmark", "validate", "--input", write(5), "--repo", dir}, &out, &out)
	if err == nil || !strings.Contains(err.Error(), "Casos[0].Defeitos[0].Fim") {
		t.Fatalf("esperava linha além do arquivo, recebeu %v", err)
	}
}

func TestBenchmarkValidateAceitaExemploDoRepositorio(t *testing.T) {
	var out bytes.Buffer
	example := filepath.Join("..", "..", "..", "..", "..", "..", "examples", "benchmark.example.json")
	if err := run(context.Background(), []string{"benchmark", "validate", "--input", example}, &out, &out); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), `"VerificadoNoRepositorio": false`) {
		t.Fatalf("sem --repo não deve verificar: %s", out.String())
	}
}

func TestBenchmarkSemSubcomandoValidoFalha(t *testing.T) {
	var out bytes.Buffer
	if err := run(context.Background(), []string{"benchmark", "run"}, &out, &out); err == nil {
		t.Fatal("esperava recusa")
	}
}
