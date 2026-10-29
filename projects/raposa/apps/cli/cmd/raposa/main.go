package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"raposa.local/cli/internal/adapters/command"
	gitadapter "raposa.local/cli/internal/adapters/git"
	"raposa.local/cli/internal/adapters/gostatic"
	"raposa.local/cli/internal/adapters/report"
	"raposa.local/cli/internal/application"
	"raposa.local/cli/internal/domain"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	if err := run(ctx, os.Args[1:], os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, "raposa:", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string, out, stderr io.Writer) error {
	if len(args) == 0 || args[0] == "help" || args[0] == "--help" {
		_, err := fmt.Fprintln(out, "Raposa\n\n  prepare --repo PATH --base REF --head REF [--output FILE]\n  review  --repo PATH --base REF --head REF [--provider-config FILE] [--provider-timeout DUR] [--format json|md] [--output FILE]\n  evaluate --input FILE [--output FILE]\n\nSem provider-config: baseline estática Go. Nenhum resultado representa aprovação automática.")
		return err
	}
	if args[0] != "prepare" && args[0] != "review" && args[0] != "evaluate" {
		return fmt.Errorf("comando desconhecido")
	}
	f := flag.NewFlagSet(args[0], flag.ContinueOnError)
	f.SetOutput(stderr)
	repo := f.String("repo", ".", "caminho do repositório Git")
	base := f.String("base", "", "commit base")
	head := f.String("head", "", "commit head")
	provider := f.String("provider-config", "", "JSON com argv do comando confiável")
	format := f.String("format", "json", "json ou md")
	output := f.String("output", "", "arquivo de saída; default stdout")
	input := f.String("input", "", "avaliação adjudicada por humano")
	timeout := f.Duration("timeout", 10*time.Minute, "prazo total")
	providerTimeout := f.Duration("provider-timeout", 120*time.Second, "prazo por chamada ao provedor; no máximo igual a --timeout")
	if err := f.Parse(args[1:]); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return err
	}
	if f.NArg() != 0 {
		return fmt.Errorf("argumento posicional inesperado")
	}
	if *format != "json" && *format != "md" {
		return fmt.Errorf("formato inválido")
	}
	if *timeout <= 0 {
		return fmt.Errorf("timeout deve ser positivo")
	}
	if *providerTimeout <= 0 {
		return fmt.Errorf("provider-timeout deve ser positivo")
	}
	if *provider != "" && *providerTimeout > *timeout {
		return fmt.Errorf("provider-timeout não pode exceder timeout")
	}
	ctx, cancel := context.WithTimeout(ctx, *timeout)
	defer cancel()
	var value any
	var result *application.Relatorio
	switch args[0] {
	case "evaluate":
		if *input == "" {
			return fmt.Errorf("informe --input")
		}
		var cases []domain.Caso
		if err := readJSON(*input, &cases); err != nil {
			return err
		}
		metrics, err := domain.Avaliar(cases)
		if err != nil {
			return err
		}
		value = metrics
	case "prepare", "review":
		if *base == "" || *head == "" {
			return fmt.Errorf("informe --base e --head")
		}
		r := gitadapter.Repository{Path: *repo}
		if args[0] == "prepare" {
			m, err := r.Snapshot(ctx, *base, *head)
			if err != nil {
				return err
			}
			value = command.Request{Version: "1", Phase: "candidates", Instructions: command.Instructions, Change: m}
		} else {
			s := application.Service{Repository: r, Analyzer: gostatic.Analyzer{}}
			if *provider != "" {
				var config struct {
					Argv []string `json:"argv"`
				}
				if err := readJSON(*provider, &config); err != nil {
					return err
				}
				if len(config.Argv) == 0 || !filepath.IsAbs(config.Argv[0]) {
					return fmt.Errorf("argv exige executável absoluto")
				}
				s.Reviewer = command.Reviewer{Argv: config.Argv, Timeout: *providerTimeout}
			}
			r, err := s.Review(ctx, *base, *head)
			if err != nil {
				return err
			}
			if *provider != "" {
				r.PrazoProvedorMS = providerTimeout.Milliseconds()
			}
			result = &r
			value = r
		}
	}
	if *format == "md" && result == nil {
		return fmt.Errorf("formato md disponível apenas para review")
	}
	write := func(w io.Writer) error {
		if *format == "md" {
			return report.Markdown(w, *result)
		}
		e := json.NewEncoder(w)
		e.SetIndent("", "  ")
		return e.Encode(value)
	}
	if *output == "" {
		return write(out)
	}
	return atomicWrite(*output, write)
}

func readJSON(path string, v any) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		return err
	}
	if info.Size() > 4<<20 {
		return fmt.Errorf("JSON excede 4 MiB")
	}
	return command.Decode(io.LimitReader(f, 4<<20), v)
}

func atomicWrite(path string, write func(io.Writer) error) error {
	f, err := os.CreateTemp(filepath.Dir(path), ".raposa-*")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())
	if err := write(f); err != nil {
		f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return os.Rename(f.Name(), path)
}
