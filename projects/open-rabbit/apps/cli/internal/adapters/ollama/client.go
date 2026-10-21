package ollama

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"open-rabbit.local/cli/internal/adapters/command"
)

type Client struct {
	HTTP     *http.Client
	Endpoint string
	Model    string
}

func (c Client) Exchange(ctx context.Context, request command.Request) (json.RawMessage, error) {
	if c.Model == "" {
		return nil, fmt.Errorf("informe o modelo escolhido com --model")
	}
	u, err := url.Parse(c.Endpoint)
	if err != nil || u.Scheme != "http" || u.User != nil || u.RawQuery != "" || u.Fragment != "" || (u.Hostname() != "127.0.0.1" && u.Hostname() != "localhost" && u.Hostname() != "::1") {
		return nil, fmt.Errorf("endpoint Ollama deve ser HTTP loopback sem credenciais")
	}
	if request.Version != "1" || (request.Phase != "candidates" && request.Phase != "verify") {
		return nil, fmt.Errorf("versão ou fase inválida")
	}
	data, err := json.Marshal(struct {
		Phase     string `json:"phase"`
		Change    any    `json:"change"`
		Candidate any    `json:"candidate,omitempty"`
	}{request.Phase, request.Change, request.Candidate})
	if err != nil {
		return nil, err
	}
	body, err := json.Marshal(map[string]any{"model": c.Model, "stream": false, "format": "json", "messages": []map[string]string{{"role": "system", "content": command.Instructions}, {"role": "user", "content": string(data)}}})
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.Endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	client := c.HTTP
	if client == nil {
		client = &http.Client{}
	}
	safeClient := *client
	safeClient.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
	response, err := safeClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Ollama indisponível: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Ollama retornou HTTP %d", response.StatusCode)
	}
	raw, err := io.ReadAll(io.LimitReader(response.Body, (2<<20)+1))
	if err != nil {
		return nil, err
	}
	if len(raw) > 2<<20 {
		return nil, fmt.Errorf("resposta Ollama excede 2 MiB")
	}
	var envelope struct {
		Done       bool   `json:"done"`
		DoneReason string `json:"done_reason"`
		Message    struct {
			Content string `json:"content"`
		} `json:"message"`
	}
	if err := json.Unmarshal(raw, &envelope); err != nil {
		return nil, fmt.Errorf("envelope Ollama inválido")
	}
	if !envelope.Done || envelope.DoneReason != "stop" || !json.Valid([]byte(envelope.Message.Content)) {
		return nil, fmt.Errorf("Ollama retornou resposta incompleta ou JSON inválido")
	}
	return json.RawMessage(envelope.Message.Content), nil
}
