package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"open-rabbit.local/cli/internal/adapters/command"
	"open-rabbit.local/cli/internal/application"
	"open-rabbit.local/cli/internal/domain"
)

func commandGit(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	cmd.Env = append(os.Environ(), "GIT_CONFIG_GLOBAL=/dev/null", "GIT_CONFIG_NOSYSTEM=1", "GIT_AUTHOR_NAME=Fixture", "GIT_AUTHOR_EMAIL=fixture@example.invalid", "GIT_COMMITTER_NAME=Fixture", "GIT_COMMITTER_EMAIL=fixture@example.invalid")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("%v: %s", err, out)
	}
	return strings.TrimSpace(string(out))
}

func TestSemanticHelper(t *testing.T) {
	if len(os.Args) == 0 || os.Args[len(os.Args)-1] != "semantic-provider" {
		return
	}
	var r command.Request
	if err := json.NewDecoder(os.Stdin).Decode(&r); err != nil {
		os.Exit(2)
	}
	if r.Phase == "verify" {
		fmt.Print(`{"confirmed":true,"reason":"Fixture de protocolo: literal true sem callback no snapshot. Isto não é avaliação de modelo."}`)
		os.Exit(0)
	}
	a := domain.Achado{Regra: "go/tls-insecure-skip-verify", Caminho: "client.go", Linha: 4, Titulo: "Verificação TLS desativada", Severidade: "high", Confianca: .9, Condicao: "conexão com certificado inválido", Impacto: "aceitação de servidor inválido", Esperado: "validar certificado", Contraevidencia: "sem callback no contexto", Evidencias: []domain.Evidencia{{Caminho: "client.go", Linha: 4, Trecho: " InsecureSkipVerify: true,"}}}
	json.NewEncoder(os.Stdout).Encode(map[string]any{"candidates": []domain.Achado{a}})
	os.Exit(0)
}

func TestCLIRevisaGitRealEExportaJSONMarkdownEContexto(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "repo com espaço")
	if err := os.Mkdir(dir, 0700); err != nil {
		t.Fatal(err)
	}
	commandGit(t, dir, "init", "--initial-branch=main")
	file := filepath.Join(dir, "client.go")
	before := "package client\nimport \"crypto/tls\"\nvar Config = tls.Config{\n InsecureSkipVerify: false,\n}\n"
	if err := os.WriteFile(file, []byte(before), 0600); err != nil {
		t.Fatal(err)
	}
	commandGit(t, dir, "add", ".")
	commandGit(t, dir, "-c", "core.hooksPath=/dev/null", "commit", "-m", "base")
	base := commandGit(t, dir, "rev-parse", "HEAD")
	if err := os.WriteFile(file, []byte(strings.Replace(before, "false", "true", 1)), 0600); err != nil {
		t.Fatal(err)
	}
	commandGit(t, dir, "add", ".")
	commandGit(t, dir, "-c", "core.hooksPath=/dev/null", "commit", "-m", "head")
	head := commandGit(t, dir, "rev-parse", "HEAD")
	exe, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	config := filepath.Join(t.TempDir(), "provider.local.json")
	raw, _ := json.Marshal(map[string]any{"argv": []string{exe, "-test.run=TestSemanticHelper", "--", "semantic-provider"}})
	if err := os.WriteFile(config, raw, 0600); err != nil {
		t.Fatal(err)
	}
	var out, stderr bytes.Buffer
	args := []string{"review", "--repo", dir, "--base", base, "--head", head, "--provider-config", config}
	if err := run(context.Background(), args, &out, &stderr); err != nil {
		t.Fatal(err)
	}
	var report application.Relatorio
	if err := json.Unmarshal(out.Bytes(), &report); err != nil {
		t.Fatal(err)
	}
	if len(report.Achados) != 1 || len(report.Descartados) != 1 || report.Modo != "static+semantic" {
		t.Fatalf("relatório incorreto: %+v", report)
	}
	if !strings.Contains(report.Achados[0].Contraevidencia, "Fixture") {
		t.Fatal("verificador não executado")
	}
	md := filepath.Join(t.TempDir(), "relatório.md")
	if err := run(context.Background(), append(args, "--format", "md", "--output", md), &out, &stderr); err != nil {
		t.Fatal(err)
	}
	content, err := os.ReadFile(md)
	if err != nil || !bytes.Contains(content, []byte("client.go:4")) {
		t.Fatalf("Markdown ausente: %v", err)
	}
	out.Reset()
	if err := run(context.Background(), []string{"prepare", "--repo", dir, "--base", base, "--head", head}, &out, &stderr); err != nil {
		t.Fatal(err)
	}
	var packet command.Request
	if err := json.Unmarshal(out.Bytes(), &packet); err != nil {
		t.Fatal(err)
	}
	if packet.Change.Head != head || packet.Phase != "candidates" {
		t.Fatal("pacote inválido")
	}
	if err := run(context.Background(), []string{"review", "--repo", dir, "--base", "inexistente", "--head", head, "--output", md}, &out, &stderr); err == nil {
		t.Fatal("esperava falha")
	}
	after, err := os.ReadFile(md)
	if err != nil || !bytes.Equal(content, after) {
		t.Fatal("falha sobrescreveu relatório anterior")
	}
}
