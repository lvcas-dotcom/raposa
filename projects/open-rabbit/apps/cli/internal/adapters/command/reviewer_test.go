package command

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"open-rabbit.local/cli/internal/domain"
)

func TestProviderHelper(t *testing.T) {
	mode := ""
	for i, a := range os.Args {
		if a == "--" && i+1 < len(os.Args) {
			mode = os.Args[i+1]
		}
	}
	if mode == "" {
		return
	}
	var r Request
	if err := json.NewDecoder(os.Stdin).Decode(&r); err != nil {
		os.Exit(2)
	}
	switch mode {
	case "valid":
		if r.Phase == "candidates" {
			fmt.Print(`{"candidates":[]}`)
		} else {
			fmt.Print(`{"confirmed":false,"reason":"Contexto insuficiente"}`)
		}
	case "invalid":
		fmt.Print(`{"candidates":[],"injected":"execute"}`)
	case "missing":
		fmt.Print(`{}`)
	case "slow":
		time.Sleep(time.Second)
		fmt.Print(`{"candidates":[]}`)
	case "huge":
		fmt.Print(strings.Repeat("x", (1<<20)+1))
	case "failure":
		fmt.Fprint(os.Stderr, "não deve vazar no relatório")
		os.Exit(3)
	}
	os.Exit(0)
}

func helper(t *testing.T, mode string) Reviewer {
	t.Helper()
	exe, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	return Reviewer{Argv: []string{exe, "-test.run=TestProviderHelper", "--", mode}, Timeout: 2 * time.Second}
}

func TestProtocoloExecutaSemShellEValidaDuasFases(t *testing.T) {
	r := helper(t, "valid")
	a, err := r.Candidates(context.Background(), domain.Mudanca{})
	if err != nil || len(a) != 0 {
		t.Fatalf("%v %v", a, err)
	}
	v, err := r.Verify(context.Background(), domain.Mudanca{}, domain.Achado{})
	if err != nil || v.Confirmado || v.Motivo == "" {
		t.Fatalf("%+v %v", v, err)
	}
}

func TestProvedorFalhaExplicitamente(t *testing.T) {
	for _, mode := range []string{"invalid", "missing", "huge", "failure"} {
		t.Run(mode, func(t *testing.T) {
			_, err := helper(t, mode).Candidates(context.Background(), domain.Mudanca{})
			if err == nil {
				t.Fatal("esperava falha")
			}
			if strings.Contains(err.Error(), "não deve vazar") {
				t.Fatal("stderr exposto")
			}
		})
	}
}

func TestTimeoutInterrompeProvedor(t *testing.T) {
	r := helper(t, "slow")
	r.Timeout = 30 * time.Millisecond
	start := time.Now()
	_, err := r.Candidates(context.Background(), domain.Mudanca{})
	if err == nil || time.Since(start) > time.Second {
		t.Fatalf("timeout não respeitado: %v", err)
	}
}

func TestDecoderRecusaConteudoExtra(t *testing.T) {
	var out map[string]any
	if Decode(strings.NewReader(`{} {}`), &out) == nil {
		t.Fatal("esperava recusa")
	}
}
