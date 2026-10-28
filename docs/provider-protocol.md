# Protocolo do provedor v1

`--provider-config` aponta para um JSON do operador com `argv`, lista de strings.
O primeiro item deve ser executável absoluto. Nenhuma expansão de shell ocorre.
Não obtenha essa configuração do repositório analisado.

Cada invocação recebe um documento JSON em stdin, contendo:

| Campo | Conteúdo |
| --- | --- |
| `version` | `"1"` |
| `phase` | `candidates` ou `verify` |
| `instructions` | regras confiáveis do pipeline |
| `change` | `Base`, `Head`, `Arquivos`, `Contexto`, `Limitacoes` |
| `candidate` | achado a refutar, só em verify |

Arquivos contêm `Caminho`, `Anterior`, `Antes`, `Depois`, `Alteradas`
(`Inicio`/`Fim`, inclusivos no head) e `DiffNumerado`. Contexto contém
`Caminho`/`Conteudo`. Esses nomes são o contrato JSON v1, preservados pelos
tipos de domínio.

`DiffNumerado` é aditivo: repete os hunks do diff com três linhas de contexto.
Cada linha traz o número no head com seis colunas, o marcador (`+` adicionada,
`-` removida, espaço para contexto), um espaço e o texto literal. Linhas
removidas não têm número. `Alteradas` é derivada das linhas `+` desse mesmo diff.

```text
@@ -2,2 +2,2 @@ func A() int {
     2   	x := 1
       - 	return x
     3 + 	return x + 1
```

`Linha` do achado usa o número exibido; `Trecho` é o texto após o marcador.

Resposta de candidates:

```json
{
  "candidates": [
    {
      "Regra": "classe-estavel",
      "Caminho": "arquivo.go",
      "Linha": 12,
      "Titulo": "Descrição objetiva do defeito",
      "Severidade": "high",
      "Confianca": 0.8,
      "Condicao": "Entrada e caminho que disparam a falha",
      "Impacto": "Efeito observável",
      "Esperado": "Contrato exigido",
      "Evidencias": [{"Caminho": "arquivo.go", "Linha": 12, "Trecho": "texto literal desde o início da linha"}],
      "Contraevidencia": "Proteções procuradas e resultado"
    }
  ]
}
```

Máximo de 20 candidatos. `ID` é atribuído pelo pipeline. Todos os campos de
justificativa são obrigatórios; confiança deve estar em 0..1. Severidades:
`critical`, `high`, `medium`, `low`. Não são mapeadas automaticamente da escala
do CodeRabbit. Evidência precisa incluir a própria linha do achado.

Resposta de verify:

```json
{"confirmed": false, "reason": "Proteção presente, hipótese incorreta ou contexto insuficiente, com justificativa"}
```

Somente um JSON, sem markdown ou logs no stdout. Campos desconhecidos, resposta
extra, conteúdo incompleto, exit code não zero, timeout e saída acima de 1 MiB
encerram a execução com erro. O stderr do processo não é colocado no relatório.

O comando é confiável e executa com permissões do operador. O diretório
temporário não restringe filesystem ou rede. Não use um agente com acesso amplo
a ferramentas se a intenção for apenas consultar um modelo. O bridge fornecido
faz uma chamada de chat ao Ollama loopback, sem ferramentas ou redirects.
Ele exige `--num-ctx` e envia `options.num_ctx`, `truncate:false` e `shift:false`:
prompt acima da janela falha no Ollama em vez de ser cortado. Também recusa a
resposta quando tokens de prompt e gerados alcançam a janela.
Consulte o [contrato Ollama](https://docs.ollama.com/api/chat).
