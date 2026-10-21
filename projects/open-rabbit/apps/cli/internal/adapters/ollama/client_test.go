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
	client := Client{Endpoint: "http://127.0.0.1:11434/api/chat", Model: "modelo-do-operador", HTTP: &http.Client{Transport: transport(func(r *http.Request) (*http.Response, error) {
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if r.Method != "POST" || body["stream"] != false || body["format"] != "json" || body["model"] != "modelo-do-operador" || body["tools"] != nil {
			t.Fatalf("request incorreto: %v", body)
		}
		return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(`{"done":true,"done_reason":"stop","message":{"content":"{\"candidates\":[]}"}}`)), Header: http.Header{}}, nil
	})}}
	out, err := client.Exchange(context.Background(), command.Request{Version: "1", Phase: "candidates"})
	if err != nil || string(out) != `{"candidates":[]}` {
		t.Fatalf("%s %v", out, err)
	}
}

func TestOllamaRecusaRespostaTruncadaEStatusInvalido(t *testing.T) {
	for _, body := range []string{`{"done":false}`, `{"done":true,"done_reason":"length","message":{"content":"{}"}}`, `{"done":true,"done_reason":"stop","message":{"content":"texto"}}`} {
		client := Client{Endpoint: "http://localhost:11434/api/chat", Model: "modelo", HTTP: &http.Client{Transport: transport(func(*http.Request) (*http.Response, error) {
			return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(body))}, nil
		})}}
		if _, err := client.Exchange(context.Background(), command.Request{Version: "1", Phase: "verify"}); err == nil {
			t.Fatal("esperava recusa")
		}
	}
}

func TestOllamaNaoEnviaCodigoParaHostRemoto(t *testing.T) {
	client := Client{Endpoint: "https://example.com/api/chat", Model: "modelo"}
	if _, err := client.Exchange(context.Background(), command.Request{Version: "1", Phase: "verify"}); err == nil {
		t.Fatal("esperava recusa")
	}
}
