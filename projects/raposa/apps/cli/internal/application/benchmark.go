package application

import (
	"context"
	"errors"
	"fmt"

	"raposa.local/cli/internal/domain"
)

type Commits interface {
	ResolveCommit(ctx context.Context, ref string) (string, error)
	CountLines(ctx context.Context, commit, path string) (int, error)
}

func VerificarNoRepositorio(ctx context.Context, r Commits, d domain.Dataset) error {
	var problems []error
	for i, c := range d.Casos {
		campo := fmt.Sprintf("Casos[%d]", i)
		headOK := true
		for _, ref := range []struct{ nome, sha string }{{"Base", c.Base}, {"Head", c.Head}} {
			got, err := r.ResolveCommit(ctx, ref.sha)
			if ctx.Err() != nil {
				return ctx.Err()
			}
			if err != nil || got != ref.sha {
				problems = append(problems, fmt.Errorf("%s.%s: commit %s não encontrado no repositório", campo, ref.nome, ref.sha))
				if ref.nome == "Head" {
					headOK = false
				}
			}
		}
		if !headOK {
			continue
		}
		for j, def := range c.Defeitos {
			n, err := r.CountLines(ctx, c.Head, def.Caminho)
			if ctx.Err() != nil {
				return ctx.Err()
			}
			if err != nil {
				problems = append(problems, fmt.Errorf("%s.Defeitos[%d].Caminho: arquivo %q ausente no head: %v", campo, j, def.Caminho, err))
			} else if def.Fim > n {
				problems = append(problems, fmt.Errorf("%s.Defeitos[%d].Fim: linha %d além do fim do arquivo (%d linhas)", campo, j, def.Fim, n))
			}
		}
	}
	return errors.Join(problems...)
}
