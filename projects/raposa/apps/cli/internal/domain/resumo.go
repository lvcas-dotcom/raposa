package domain

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

const (
	OrigemModelo      = "modelo"
	LimiteResumo      = 1500
	LimiteLinhaResumo = 240
)

type ResumoArquivo struct {
	Caminho string
	Linha   string
	Ancora  *Intervalo `json:",omitempty"`
}

type Resumo struct {
	Origem   string
	Resumo   string
	Arquivos []ResumoArquivo
}

func (r Resumo) Validate(m Mudanca) error {
	if strings.TrimSpace(r.Resumo) == "" {
		return fmt.Errorf("resumo vazio")
	}
	if utf8.RuneCountInString(r.Resumo) > LimiteResumo {
		return fmt.Errorf("resumo excede %d caracteres", LimiteResumo)
	}
	files := map[string]Arquivo{}
	for _, f := range m.Arquivos {
		files[f.Caminho] = f
	}
	seen := map[string]bool{}
	for _, e := range r.Arquivos {
		f, ok := files[e.Caminho]
		if !ok {
			return fmt.Errorf("resumo cita arquivo fora da mudança")
		}
		if seen[e.Caminho] {
			return fmt.Errorf("resumo repete arquivo")
		}
		seen[e.Caminho] = true
		if strings.TrimSpace(e.Linha) == "" || strings.ContainsAny(e.Linha, "\r\n") {
			return fmt.Errorf("linha de arquivo vazia ou com quebra de linha")
		}
		if utf8.RuneCountInString(e.Linha) > LimiteLinhaResumo {
			return fmt.Errorf("linha de arquivo excede %d caracteres", LimiteLinhaResumo)
		}
		if err := validaAncora(e.Ancora, f.Alteradas); err != nil {
			return err
		}
	}
	if len(seen) != len(files) {
		return fmt.Errorf("resumo não cobre todos os arquivos alterados")
	}
	return nil
}

func validaAncora(a *Intervalo, alteradas []Intervalo) error {
	if len(alteradas) == 0 {
		if a != nil {
			return fmt.Errorf("âncora em arquivo sem linhas alteradas")
		}
		return nil
	}
	if a == nil {
		return fmt.Errorf("arquivo alterado sem âncora")
	}
	for _, r := range alteradas {
		if a.Inicio >= r.Inicio && a.Fim <= r.Fim && a.Inicio <= a.Fim {
			return nil
		}
	}
	return fmt.Errorf("âncora fora das linhas alteradas")
}
