package git

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"open-rabbit.local/cli/internal/domain"
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
			if !strings.Contains(f.DiffNumerado, "     3 + func A() int { return 2 }\n") || !strings.Contains(f.DiffNumerado, "       - func A() int { return 1 }\n") {
				t.Fatalf("diff numerado incorreto:\n%s", f.DiffNumerado)
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

func TestDiffNumeradoMarcaSomenteLinhasAdicionadasComoAlteradas(t *testing.T) {
	diff := "diff --git a/a.go b/a.go\n--- a/a.go\n+++ b/a.go\n" +
		"@@ -2,3 +2,2 @@ func A() {\n contexto\n-removida\n contexto2\n" +
		"@@ -8,2 +7,4 @@\n mantida\n-antiga\n+nova1\n+nova2\n+nova3\n\\ No newline at end of file\n"
	ranges, numbered := parseHunks([]byte(diff))
	if len(ranges) != 1 || ranges[0].Inicio != 8 || ranges[0].Fim != 10 {
		t.Fatalf("intervalos incorretos: %+v", ranges)
	}
	expected := "@@ -2,3 +2,2 @@ func A() {\n" +
		"     2   contexto\n" +
		"       - removida\n" +
		"     3   contexto2\n" +
		"@@ -8,2 +7,4 @@\n" +
		"     7   mantida\n" +
		"       - antiga\n" +
		"     8 + nova1\n" +
		"     9 + nova2\n" +
		"    10 + nova3\n"
	if numbered != expected {
		t.Fatalf("diff numerado incorreto:\n%s", numbered)
	}
}

func TestDiffNumeradoSeparaBlocosAdicionadosNaoContiguos(t *testing.T) {
	ranges, _ := parseHunks([]byte("@@ -1,3 +1,5 @@\n+a\n x\n+b\n+c\n y\n"))
	if len(ranges) != 2 || ranges[0] != (domain.Intervalo{Inicio: 1, Fim: 1}) || ranges[1] != (domain.Intervalo{Inicio: 3, Fim: 4}) {
		t.Fatalf("intervalos incorretos: %+v", ranges)
	}
}

func intervalosSemContexto(t *testing.T, dir, base, head, caminho string) []domain.Intervalo {
	t.Helper()
	saida := gitcmd(t, dir, "diff", "--no-color", "--unified=0", base, head, "--", caminho)
	cabecalho := regexp.MustCompile(`^@@ -\d+(?:,\d+)? \+(\d+)(?:,(\d+))? @@`)
	intervalos := []domain.Intervalo{}
	for _, linha := range strings.Split(saida, "\n") {
		m := cabecalho.FindStringSubmatch(linha)
		if m == nil {
			continue
		}
		inicio, _ := strconv.Atoi(m[1])
		quantidade := 1
		if m[2] != "" {
			quantidade, _ = strconv.Atoi(m[2])
		}
		if quantidade > 0 {
			intervalos = append(intervalos, domain.Intervalo{Inicio: inicio, Fim: inicio + quantidade - 1})
		}
	}
	return intervalos
}

func TestAlteradasEquivaleAoDiffSemContextoEmCasosReais(t *testing.T) {
	longo := "l1\nl2\nl3\nl4\nl5\nl6\nl7\nl8\nl9\nl10\nl11\nl12\n"
	casos := []struct {
		nome, antes, depois string
		quantidade          int
	}{
		{"modificacao simples", "a\nb\nc\n", "a\nB\nc\n", 1},
		{"remocao pura", "a\nb\nc\n", "a\nc\n", 0},
		{"dois blocos separados por contexto", longo, strings.Replace(strings.Replace(longo, "l2\n", "X\n", 1), "l11\n", "Y\n", 1), 2},
		{"remocoes e adicoes intercaladas", "a\nb\nc\nd\ne\n", "a\nB\nD\ne\n", 1},
		{"arquivo sem newline final", "a\nb", "a\nb\nc", 1},
		{"adicao no fim com newline", "a\n", "a\nb\nc\n", 1},
	}
	for _, c := range casos {
		t.Run(c.nome, func(t *testing.T) {
			dir := t.TempDir()
			gitcmd(t, dir, "init", "--initial-branch=main")
			write(t, dir, "f.txt", c.antes)
			write(t, dir, "README.md", "x\n")
			gitcmd(t, dir, "add", ".")
			gitcmd(t, dir, "-c", "core.hooksPath=/dev/null", "commit", "-m", "base")
			base := gitcmd(t, dir, "rev-parse", "HEAD")
			write(t, dir, "f.txt", c.depois)
			gitcmd(t, dir, "add", ".")
			gitcmd(t, dir, "-c", "core.hooksPath=/dev/null", "commit", "-m", "head")
			head := gitcmd(t, dir, "rev-parse", "HEAD")
			m, err := (Repository{dir}).Snapshot(context.Background(), base, head)
			if err != nil {
				t.Fatal(err)
			}
			esperado := intervalosSemContexto(t, dir, base, head, "f.txt")
			if len(esperado) != c.quantidade {
				t.Fatalf("premissa do caso falhou: %+v", esperado)
			}
			var obtido []domain.Intervalo
			for _, f := range m.Arquivos {
				if f.Caminho == "f.txt" {
					obtido = f.Alteradas
				}
			}
			if !reflect.DeepEqual(append([]domain.Intervalo{}, obtido...), esperado) {
				t.Fatalf("Alteradas = %+v, esperado %+v", obtido, esperado)
			}
		})
	}
}

func TestArquivoNovoTemTodasAsLinhasAlteradas(t *testing.T) {
	dir := t.TempDir()
	gitcmd(t, dir, "init", "--initial-branch=main")
	write(t, dir, "README.md", "x\n")
	gitcmd(t, dir, "add", ".")
	gitcmd(t, dir, "-c", "core.hooksPath=/dev/null", "commit", "-m", "base")
	base := gitcmd(t, dir, "rev-parse", "HEAD")
	write(t, dir, "novo.txt", "um\ndois\ntres\n")
	gitcmd(t, dir, "add", ".")
	gitcmd(t, dir, "-c", "core.hooksPath=/dev/null", "commit", "-m", "head")
	head := gitcmd(t, dir, "rev-parse", "HEAD")
	m, err := (Repository{dir}).Snapshot(context.Background(), base, head)
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range m.Arquivos {
		if f.Caminho == "novo.txt" {
			if !reflect.DeepEqual(f.Alteradas, []domain.Intervalo{{Inicio: 1, Fim: 3}}) || !strings.Contains(f.DiffNumerado, "     3 + tres\n") {
				t.Fatalf("arquivo novo incorreto: %+v", f)
			}
			return
		}
	}
	t.Fatal("novo.txt ausente")
}

func TestConteudoComPrefixoDeCabecalhoEhAdicaoERemocao(t *testing.T) {
	diff := "diff --git a/a.md b/a.md\n--- a/a.md\n+++ b/a.md\n@@ -1,2 +1,2 @@\n----\n-- x\n+++\n+++ y\n"
	ranges, numbered := parseHunks([]byte(diff))
	esperado := "@@ -1,2 +1,2 @@\n       - ---\n       - - x\n     1 + ++\n     2 + ++ y\n"
	if numbered != esperado {
		t.Fatalf("diff numerado incorreto:\n%s", numbered)
	}
	if !reflect.DeepEqual(ranges, []domain.Intervalo{{Inicio: 1, Fim: 2}}) {
		t.Fatalf("intervalos incorretos: %+v", ranges)
	}
}

func TestAchadoCopiadoDoDiffNumeradoPassaNaValidacao(t *testing.T) {
	dir := t.TempDir()
	gitcmd(t, dir, "init", "--initial-branch=main")
	write(t, dir, "a.go", "package a\n\nfunc A() int {\n\treturn 1\n}\n")
	gitcmd(t, dir, "add", ".")
	gitcmd(t, dir, "-c", "core.hooksPath=/dev/null", "commit", "-m", "base")
	base := gitcmd(t, dir, "rev-parse", "HEAD")
	write(t, dir, "a.go", "package a\n\nfunc A() int {\n\tx := 0\n\treturn 1 / x\n}\n")
	gitcmd(t, dir, "add", ".")
	gitcmd(t, dir, "-c", "core.hooksPath=/dev/null", "commit", "-m", "head")
	head := gitcmd(t, dir, "rev-parse", "HEAD")
	m, err := (Repository{dir}).Snapshot(context.Background(), base, head)
	if err != nil {
		t.Fatal(err)
	}
	var numerado string
	for _, f := range m.Arquivos {
		if f.Caminho == "a.go" {
			numerado = f.DiffNumerado
		}
	}
	linhaAlvo := regexp.MustCompile(`(?m)^\s*(\d+) \+ (\t.*1 / x)$`).FindStringSubmatch(numerado)
	if linhaAlvo == nil {
		t.Fatalf("linha adicionada não encontrada:\n%s", numerado)
	}
	linha, _ := strconv.Atoi(linhaAlvo[1])
	achado := domain.Achado{Regra: "R", Caminho: "a.go", Linha: linha, Titulo: "t", Severidade: "high", Confianca: 0.9,
		Condicao: "c", Impacto: "i", Esperado: "e", Contraevidencia: "ce",
		Evidencias: []domain.Evidencia{{Caminho: "a.go", Linha: linha, Trecho: linhaAlvo[2]}}}
	if err := achado.Validate(m); err != nil {
		t.Fatalf("achado copiado do diff numerado rejeitado (linha %d): %v", linha, err)
	}
}
