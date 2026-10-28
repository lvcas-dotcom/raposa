package domain

import (
	"crypto/sha256"
	"fmt"
	"strings"
)

type Intervalo struct {
	Inicio int
	Fim    int
}

type Arquivo struct {
	Caminho      string
	Anterior     string
	Antes        string
	Depois       string
	Alteradas    []Intervalo
	DiffNumerado string
}

type Contexto struct {
	Caminho  string
	Conteudo string
}

type Mudanca struct {
	Base       string
	Head       string
	Arquivos   []Arquivo
	Contexto   []Contexto
	Limitacoes []string
}

type Evidencia struct {
	Caminho string
	Linha   int
	Trecho  string
}

type Achado struct {
	ID              string
	Regra           string
	Caminho         string
	Linha           int
	Titulo          string
	Severidade      string
	Confianca       float64
	Condicao        string
	Impacto         string
	Esperado        string
	Evidencias      []Evidencia
	Contraevidencia string
}

func (a Achado) Fingerprint() string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(fmt.Sprintf("%s:%d:%s", a.Caminho, a.Linha, a.Regra))))[:16]
}

func (a Achado) Validate(m Mudanca) error {
	for _, field := range []string{a.Regra, a.Titulo, a.Condicao, a.Impacto, a.Esperado, a.Contraevidencia} {
		if strings.TrimSpace(field) == "" {
			return fmt.Errorf("achado sem justificativa completa")
		}
	}
	if a.Severidade != "critical" && a.Severidade != "high" && a.Severidade != "medium" && a.Severidade != "low" {
		return fmt.Errorf("severidade inválida")
	}
	if !(a.Confianca >= 0 && a.Confianca <= 1) {
		return fmt.Errorf("confiança inválida")
	}
	changed := false
	for _, f := range m.Arquivos {
		if f.Caminho != a.Caminho {
			continue
		}
		for _, r := range f.Alteradas {
			if a.Linha >= r.Inicio && a.Linha <= r.Fim {
				changed = true
			}
		}
	}
	if !changed {
		return fmt.Errorf("localização fora das linhas alteradas")
	}
	if len(a.Evidencias) == 0 {
		return fmt.Errorf("achado sem evidência")
	}
	anchored := false
	for _, e := range a.Evidencias {
		content, ok := m.Content(e.Caminho)
		if !ok || e.Linha < 1 || strings.TrimSpace(e.Trecho) == "" {
			return fmt.Errorf("evidência fora do contexto")
		}
		lines := strings.Split(content, "\n")
		if e.Linha > len(lines) || !strings.HasPrefix(strings.Join(lines[e.Linha-1:], "\n"), e.Trecho) {
			return fmt.Errorf("evidência não corresponde ao snapshot")
		}
		if e.Caminho == a.Caminho && e.Linha == a.Linha {
			anchored = true
		}
	}
	if !anchored {
		return fmt.Errorf("falta evidência na linha do achado")
	}
	return nil
}

func (m Mudanca) Content(path string) (string, bool) {
	for _, f := range m.Arquivos {
		if f.Caminho == path {
			return f.Depois, true
		}
	}
	for _, f := range m.Contexto {
		if f.Caminho == path {
			return f.Conteudo, true
		}
	}
	return "", false
}
