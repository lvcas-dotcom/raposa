package git

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"open-rabbit.local/cli/internal/domain"
)

const maxOutput = 4 << 20
const maxBlob = 128 << 10
const maxContext = 2 << 20

type Repository struct{ Path string }

type boundedBuffer struct{ bytes.Buffer }

func (b *boundedBuffer) Write(p []byte) (int, error) {
	if b.Len()+len(p) > maxOutput {
		return 0, fmt.Errorf("saída Git excede limite")
	}
	return b.Buffer.Write(p)
}

func (r Repository) run(ctx context.Context, args ...string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()
	absolute, err := filepath.Abs(r.Path)
	if err != nil {
		return nil, err
	}
	argv := append([]string{"--no-pager", "--literal-pathspecs", "-C", absolute, "-c", "core.quotePath=false", "-c", "core.fsmonitor=false"}, args...)
	cmd := exec.CommandContext(ctx, "git", argv...)
	for _, v := range os.Environ() {
		if !strings.HasPrefix(v, "GIT_") {
			cmd.Env = append(cmd.Env, v)
		}
	}
	cmd.Env = append(cmd.Env, "GIT_CONFIG_NOSYSTEM=1", "GIT_CONFIG_GLOBAL=/dev/null", "GIT_NO_REPLACE_OBJECTS=1", "GIT_TERMINAL_PROMPT=0", "LC_ALL=C")
	var out boundedBuffer
	cmd.Stdout = &out
	cmd.WaitDelay = time.Second
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("falha no Git %s: %w", args[0], err)
	}
	return out.Bytes(), nil
}

func (r Repository) resolve(ctx context.Context, ref string) (string, error) {
	if ref == "" {
		return "", fmt.Errorf("referência vazia")
	}
	b, err := r.run(ctx, "rev-parse", "--verify", "--end-of-options", ref+"^{commit}")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(b)), nil
}

func (r Repository) blob(ctx context.Context, sha, p string) (string, bool, error) {
	size, err := r.run(ctx, "cat-file", "-s", sha+":"+p)
	if err != nil {
		return "", false, err
	}
	n, err := strconv.Atoi(strings.TrimSpace(string(size)))
	if err != nil {
		return "", false, err
	}
	if n > maxBlob {
		return "", false, nil
	}
	b, err := r.run(ctx, "show", sha+":"+p)
	if err != nil {
		return "", false, err
	}
	if bytes.ContainsRune(b, 0) || !utf8.Valid(b) {
		return "", false, nil
	}
	return string(b), true, nil
}

var hunk = regexp.MustCompile(`(?m)^@@ -[0-9]+(?:,[0-9]+)? \+([0-9]+)(?:,([0-9]+))? @@`)

func intervals(diff []byte) []domain.Intervalo {
	ranges := []domain.Intervalo{}
	for _, m := range hunk.FindAllSubmatch(diff, -1) {
		start, _ := strconv.Atoi(string(m[1]))
		count := 1
		if len(m[2]) > 0 {
			count, _ = strconv.Atoi(string(m[2]))
		}
		if count > 0 {
			ranges = append(ranges, domain.Intervalo{Inicio: start, Fim: start + count - 1})
		}
	}
	return ranges
}

func (r Repository) Snapshot(ctx context.Context, base, head string) (domain.Mudanca, error) {
	b, err := r.resolve(ctx, base)
	if err != nil {
		return domain.Mudanca{}, err
	}
	h, err := r.resolve(ctx, head)
	if err != nil {
		return domain.Mudanca{}, err
	}
	m := domain.Mudanca{Base: b, Head: h, Arquivos: []domain.Arquivo{}, Contexto: []domain.Contexto{}, Limitacoes: []string{}}
	status, err := r.run(ctx, "diff", "--no-ext-diff", "--no-textconv", "--name-status", "-z", "--find-renames", b, h, "--")
	if err != nil {
		return m, err
	}
	parts := strings.Split(strings.TrimSuffix(string(status), "\x00"), "\x00")
	if len(status) == 0 {
		return m, nil
	}
	used := map[string]bool{}
	total := 0
	for i := 0; i < len(parts); {
		if len(m.Arquivos) >= 200 {
			return m, fmt.Errorf("mais de 200 arquivos alterados; divida a revisão")
		}
		kind := parts[i]
		i++
		if i >= len(parts) {
			return m, fmt.Errorf("status Git incompleto")
		}
		p := parts[i]
		i++
		old := p
		if strings.HasPrefix(kind, "R") || strings.HasPrefix(kind, "C") {
			if i >= len(parts) {
				return m, fmt.Errorf("rename Git incompleto")
			}
			p = parts[i]
			i++
		}
		used[p] = true
		if kind == "D" {
			m.Limitacoes = append(m.Limitacoes, "Arquivo removido sem linha nova para comentário: "+p)
			continue
		}
		if kind == "T" {
			m.Limitacoes = append(m.Limitacoes, "Mudança de tipo ignorada: "+p)
			continue
		}
		after, ok, err := r.blob(ctx, h, p)
		if err != nil {
			return m, err
		}
		if !ok {
			m.Limitacoes = append(m.Limitacoes, "Arquivo binário, não UTF-8 ou maior que 128 KiB: "+p)
			continue
		}
		before := ""
		if kind != "A" {
			var prior bool
			before, prior, err = r.blob(ctx, b, old)
			if err != nil {
				return m, err
			}
			if !prior {
				m.Limitacoes = append(m.Limitacoes, "Base indisponível por limite: "+old)
			}
		}
		total += len(before) + len(after)
		if total > maxContext {
			return m, fmt.Errorf("snapshot alterado excede 2 MiB; divida a revisão")
		}
		diff, err := r.run(ctx, "diff", "--no-ext-diff", "--no-textconv", "--no-color", "--unified=0", b, h, "--", old, p)
		if err != nil {
			return m, err
		}
		m.Arquivos = append(m.Arquivos, domain.Arquivo{Caminho: p, Anterior: old, Antes: before, Depois: after, Alteradas: intervals(diff)})
	}
	tree, err := r.run(ctx, "ls-tree", "-r", "-z", "--name-only", h)
	if err != nil {
		return m, err
	}
	paths := strings.Split(strings.TrimSuffix(string(tree), "\x00"), "\x00")
	for _, p := range paths {
		if used[p] {
			continue
		}
		selectPath := p == "AGENTS.md" || p == "README.md" || p == "go.mod"
		for _, f := range m.Arquivos {
			if path.Ext(p) == ".go" && path.Dir(p) == path.Dir(f.Caminho) {
				selectPath = true
			}
			if path.Base(p) == "AGENTS.md" && strings.HasPrefix(f.Caminho, path.Dir(p)+"/") {
				selectPath = true
			}
		}
		if !selectPath {
			continue
		}
		if len(m.Contexto) >= 40 {
			m.Limitacoes = append(m.Limitacoes, "Contexto auxiliar limitado a 40 arquivos.")
			break
		}
		content, ok, err := r.blob(ctx, h, p)
		if err != nil {
			return m, err
		}
		if !ok {
			m.Limitacoes = append(m.Limitacoes, "Contexto omitido pelo limite: "+p)
			continue
		}
		if total+len(content) > maxContext {
			m.Limitacoes = append(m.Limitacoes, "Contexto auxiliar limitado ao orçamento de 2 MiB.")
			break
		}
		total += len(content)
		m.Contexto = append(m.Contexto, domain.Contexto{Caminho: p, Conteudo: content})
	}
	m.Limitacoes = append(m.Limitacoes, "Contexto auxiliar: arquivos Go do mesmo diretório, AGENTS ancestrais e README/go.mod raiz; sem grafo completo de consumidores ou execução de testes.")
	return m, nil
}
