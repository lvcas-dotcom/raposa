package ollama

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"open-rabbit.local/cli/internal/adapters/command"
)

type transport func(*http.Request) (*http.Response, error)

func (t transport) RoundTrip(r *http.Request) (*http.Response, error) { return t(r) }

func TestOllamaEnviaProtocoloSemFerramentasNemStreaming(t *testing.T) {
	client := Client{Endpoint: "http://127.0.0.1:11434/api/chat", Model: "modelo-do-operador", ContextLength: 8192, HTTP: &http.Client{Transport: transport(func(r *http.Request) (*http.Response, error) {
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if r.Method != "POST" || body["stream"] != false || body["format"] != "json" || body["model"] != "modelo-do-operador" || body["tools"] != nil || body["truncate"] != false || body["shift"] != false {
			t.Fatalf("request incorreto: %v", body)
		}
		if options, ok := body["options"].(map[string]any); !ok || options["num_ctx"] != float64(8192) {
			t.Fatalf("janela de contexto ausente: %v", body["options"])
		}
		return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(`{"done":true,"done_reason":"stop","message":{"content":"{\"candidates\":[]}"}}`)), Header: http.Header{}}, nil
	})}}
	out, err := client.Exchange(context.Background(), command.Request{Version: "1", Phase: "candidates"})
	if err != nil || string(out) != `{"candidates":[]}` {
		t.Fatalf("%s %v", out, err)
	}
}

func TestOllamaRecusaRespostaTruncadaEStatusInvalido(t *testing.T) {
	for _, body := range []string{`{"done":false}`, `{"done":true,"done_reason":"length","message":{"content":"{}"}}`, `{"done":true,"done_reason":"stop","message":{"content":"texto"}}`, `{"done":true,"done_reason":"stop","prompt_eval_count":4000,"eval_count":96,"message":{"content":"{}"}}`} {
		client := Client{Endpoint: "http://localhost:11434/api/chat", Model: "modelo", ContextLength: 4096, HTTP: &http.Client{Transport: transport(func(*http.Request) (*http.Response, error) {
			return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(body))}, nil
		})}}
		if _, err := client.Exchange(context.Background(), command.Request{Version: "1", Phase: "verify"}); err == nil {
			t.Fatal("esperava recusa")
		}
	}
}

func TestOllamaNaoEnviaCodigoParaHostRemoto(t *testing.T) {
	client := Client{Endpoint: "https://example.com/api/chat", Model: "modelo", ContextLength: 4096}
	if _, err := client.Exchange(context.Background(), command.Request{Version: "1", Phase: "verify"}); err == nil {
		t.Fatal("esperava recusa")
	}
}

func TestOllamaExigeJanelaDeContextoAntesDeEnviar(t *testing.T) {
	client := Client{Endpoint: "http://127.0.0.1:11434/api/chat", Model: "modelo", HTTP: &http.Client{Transport: transport(func(*http.Request) (*http.Response, error) {
		t.Fatal("não deveria chamar o Ollama")
		return nil, nil
	})}}
	if _, err := client.Exchange(context.Background(), command.Request{Version: "1", Phase: "candidates"}); err == nil || !strings.Contains(err.Error(), "--num-ctx") {
		t.Fatalf("esperava recusa por janela ausente: %v", err)
	}
}

func TestOllamaExpoeMotivoDoPromptAcimaDaJanela(t *testing.T) {
	client := Client{Endpoint: "http://127.0.0.1:11434/api/chat", Model: "modelo", ContextLength: 4096, HTTP: &http.Client{Transport: transport(func(*http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: 400, Body: io.NopCloser(strings.NewReader(`{"error":"the prompt is longer than the context length"}`))}, nil
	})}}
	_, err := client.Exchange(context.Background(), command.Request{Version: "1", Phase: "candidates"})
	if err == nil || !strings.Contains(err.Error(), "HTTP 400: the prompt is longer than the context length") {
		t.Fatalf("esperava motivo do Ollama: %v", err)
	}
}
