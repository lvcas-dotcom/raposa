package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"raposa.local/cli/internal/adapters/command"
	"raposa.local/cli/internal/adapters/ollama"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "ollama-provider:", err)
		os.Exit(1)
	}
}

func run() error {
	model := flag.String("model", "", "modelo escolhido pelo operador")
	endpoint := flag.String("endpoint", "http://127.0.0.1:11434/api/chat", "endpoint local Ollama")
	contextLength := flag.Int("num-ctx", 0, "janela de contexto em tokens configurada para o modelo")
	flag.Parse()
	if flag.NArg() != 0 {
		return fmt.Errorf("argumento inesperado")
	}
	raw, err := io.ReadAll(io.LimitReader(os.Stdin, (4<<20)+1))
	if err != nil {
		return err
	}
	if len(raw) > 4<<20 {
		return fmt.Errorf("request excede 4 MiB")
	}
	var request command.Request
	if err := command.Decode(bytes.NewReader(raw), &request); err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 110*time.Second)
	defer cancel()
	response, err := (ollama.Client{Model: *model, Endpoint: *endpoint, ContextLength: *contextLength}).Exchange(ctx, request)
	if err != nil {
		return err
	}
	_, err = os.Stdout.Write(append(response, '\n'))
	return err
}
