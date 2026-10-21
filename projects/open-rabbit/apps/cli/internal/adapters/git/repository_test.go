package git

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func gitcmd(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	cmd.Env = append(os.Environ(), "GIT_CONFIG_NOSYSTEM=1", "GIT_CONFIG_GLOBAL=/dev/null", "GIT_AUTHOR_NAME=Fixture", "GIT_AUTHOR_EMAIL=fixture@example.invalid", "GIT_COMMITTER_NAME=Fixture", "GIT_COMMITTER_EMAIL=fixture@example.invalid")
	b, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %s: %v", args, b, err)
	}
	return strings.TrimSpace(string(b))
}

func write(t *testing.T, dir, path, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, path), []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
}

func TestSnapshotIgnoraArvoreLocalEIncluiRenames(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "repo com espaços")
	if err := os.Mkdir(dir, 0700); err != nil {
		t.Fatal(err)
	}
	gitcmd(t, dir, "init", "--initial-branch=main")
	stable := strings.Repeat("// contexto estável para detecção de rename\n", 15)
	write(t, dir, "antes.go", "package demo\n\nfunc A() int { return 1 }\n"+stable)
	write(t, dir, "helper_test.go", "package demo\n")
	gitcmd(t, dir, "add", ".")
	gitcmd(t, dir, "-c", "core.hooksPath=/dev/null", "commit", "-m", "base")
	base := gitcmd(t, dir, "rev-parse", "HEAD")
	gitcmd(t, dir, "mv", "antes.go", "depois com espaço.go")
	write(t, dir, "depois com espaço.go", "package demo\n\nfunc A() int { return 2 }\n"+stable)
	gitcmd(t, dir, "add", ".")
	gitcmd(t, dir, "-c", "core.hooksPath=/dev/null", "commit", "-m", "head")
	head := gitcmd(t, dir, "rev-parse", "HEAD")
	write(t, dir, "depois com espaço.go", "não faz parte do commit")
	m, err := (Repository{dir}).Snapshot(context.Background(), base, head)
	if err != nil {
		t.Fatal(err)
	}
	if m.Base != base || m.Head != head {
		t.Fatal("SHAs alterados")
	}
	found := false
	for _, f := range m.Arquivos {
		if f.Caminho == "depois com espaço.go" {
			found = true
			if !strings.Contains(f.Depois, "return 2") || len(f.Alteradas) != 1 || f.Alteradas[0].Inicio != 3 {
				t.Fatalf("snapshot incorreto: %+v", f)
			}
		}
	}
	if !found || len(m.Contexto) != 1 {
		t.Fatalf("contexto incompleto: %+v", m)
	}
	if gitcmd(t, dir, "rev-parse", "HEAD") != head {
		t.Fatal("HEAD foi modificado")
	}
}

func TestRecusaReferenciaInvalidaSemInterpretarFlag(t *testing.T) {
	dir := t.TempDir()
	gitcmd(t, dir, "init")
	if _, err := (Repository{dir}).Snapshot(context.Background(), "--help", "HEAD"); err == nil {
		t.Fatal("esperava recusa")
	}
}

func TestDiffNaoIncluiLinhasDeContextoOuRemocao(t *testing.T) {
	r := intervals([]byte("@@ -2,3 +2,0 @@\n@@ -8 +5,2 @@\n@@ -12 +9 @@\n"))
	if len(r) != 2 || r[0].Inicio != 5 || r[0].Fim != 6 || r[1].Inicio != 9 {
		t.Fatalf("intervalos incorretos: %+v", r)
	}
}
