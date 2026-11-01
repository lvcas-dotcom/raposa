package domain

import (
	"errors"
	"fmt"
	"path"
	"strings"
)

const VersaoDataset = "1"

var classesDefeito = map[string]bool{"autenticacao": true, "fluxo-de-erro": true, "contrato": true, "concorrencia": true, "dados": true, "seguranca": true, "logica": true}
var origensRotulo = map[string]bool{"humano": true, "szz": true, "sintetico": true}
var particoes = map[string]bool{"desenvolvimento": true, "reservado": true}

type DefeitoConhecido struct {
	ID        string
	Caminho   string
	Inicio    int
	Fim       int
	Classe    string
	Descricao string
	Origem    string
}

type CasoBenchmark struct {
	ID        string
	Linguagem string
	Base      string
	Head      string
	Particao  string
	Limpo     bool
	Defeitos  []DefeitoConhecido
}

type Dataset struct {
	Versao      string
	Repositorio string
	Licenca     string
	Casos       []CasoBenchmark
}

type ResumoDataset struct {
	Casos      int
	Defeitos   int
	Limpos     int
	Reservados int
}

type problemas []error

func (p *problemas) add(campo, format string, args ...any) {
	*p = append(*p, fmt.Errorf("%s: %s", campo, fmt.Sprintf(format, args...)))
}

func ValidarDataset(d Dataset) (ResumoDataset, error) {
	var p problemas
	if d.Versao != VersaoDataset {
		p.add("Versao", "versão %q não suportada; use %q", d.Versao, VersaoDataset)
	}
	if !strings.HasPrefix(d.Repositorio, "https://") || strings.ContainsAny(d.Repositorio, " \t\n") {
		p.add("Repositorio", "informe a URL https do repositório público")
	}
	if strings.TrimSpace(d.Licenca) == "" {
		p.add("Licenca", "informe a licença do repositório de origem")
	}
	if len(d.Casos) == 0 {
		p.add("Casos", "dataset sem casos")
	}
	r := ResumoDataset{Casos: len(d.Casos)}
	casos := map[string]bool{}
	for i, c := range d.Casos {
		campo := fmt.Sprintf("Casos[%d]", i)
		if c.ID == "" || casos[c.ID] {
			p.add(campo+".ID", "ID vazio ou repetido %q", c.ID)
		}
		casos[c.ID] = true
		if strings.TrimSpace(c.Linguagem) == "" {
			p.add(campo+".Linguagem", "informe a linguagem")
		}
		if !sha(c.Base) {
			p.add(campo+".Base", "SHA completo obrigatório")
		}
		if !sha(c.Head) {
			p.add(campo+".Head", "SHA completo obrigatório")
		}
		if c.Base == c.Head {
			p.add(campo+".Head", "head igual à base")
		}
		if !particoes[c.Particao] {
			p.add(campo+".Particao", "partição %q inválida; use desenvolvimento ou reservado", c.Particao)
		}
		if c.Limpo && len(c.Defeitos) > 0 {
			p.add(campo+".Limpo", "caso limpo não pode ter defeitos")
		}
		if !c.Limpo && len(c.Defeitos) == 0 {
			p.add(campo+".Defeitos", "caso sem defeitos precisa de Limpo: true")
		}
		if c.Limpo {
			r.Limpos++
		}
		if c.Particao == "reservado" {
			r.Reservados++
		}
		r.Defeitos += len(c.Defeitos)
		defeitos := map[string]bool{}
		for j, def := range c.Defeitos {
			campoDefeito := fmt.Sprintf("%s.Defeitos[%d]", campo, j)
			if def.ID == "" || defeitos[def.ID] {
				p.add(campoDefeito+".ID", "ID vazio ou repetido %q", def.ID)
			}
			defeitos[def.ID] = true
			if !caminhoRelativo(def.Caminho) {
				p.add(campoDefeito+".Caminho", "caminho relativo à raiz do repositório obrigatório, recebido %q", def.Caminho)
			}
			if def.Inicio < 1 || def.Fim < def.Inicio {
				p.add(campoDefeito+".Fim", "intervalo de linhas inválido %d–%d", def.Inicio, def.Fim)
			}
			if !classesDefeito[def.Classe] {
				p.add(campoDefeito+".Classe", "classe desconhecida %q", def.Classe)
			}
			if strings.TrimSpace(def.Descricao) == "" {
				p.add(campoDefeito+".Descricao", "descreva o defeito")
			}
			if !origensRotulo[def.Origem] {
				p.add(campoDefeito+".Origem", "origem %q inválida; use humano, szz ou sintetico", def.Origem)
			}
		}
	}
	return r, errors.Join(p...)
}

func caminhoRelativo(p string) bool {
	return p != "" && path.Clean(p) == p && !path.IsAbs(p) && p != ".." && !strings.HasPrefix(p, "../") && !strings.Contains(p, `\`)
}
