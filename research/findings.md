# Investigação — REA e CodeRabbit

Data: 2026-10-21. Escopo: documentação e artefatos públicos, sem conta autenticada.
Esta investigação não recupera o backend do CodeRabbit.

## Evidência obtida

REA 6.7.0 instalado em `~/.local/share/rea-tool`, com lockfile e integrity
verificáveis; `~/.local/bin/rea` aponta para a instalação. Scripts de instalação
foram desativados. Nenhum registro MCP global foi alterado. A CLI foi usada
diretamente nesta sessão.

Fontes: [repositório REA](https://github.com/morluto/rea),
[contrato CLI](https://github.com/morluto/rea/blob/main/docs/cli.md),
[observação CDP](https://github.com/morluto/rea/blob/main/docs/browser-observation.md).

| Execução | Resultado preservado |
| --- | --- |
| `rea --version` | `6.7.0` |
| `list-browser-targets` | página pública de arquitetura identificada no Chrome temporário |
| `inspect-web-page` | [captura](evidence/rea-page.json), 54 scripts da origem selecionada |
| `analyze-web-bundle` | [falha registrada](evidence/rea-bundle-failure.json), OOM/exit 134 |
| `analyze-javascript-application` | [amostra estática](evidence/rea-static-sample.json), exit 0 |

A captura de página é passiva: não inclui tráfego anterior à conexão, corpos,
cookies ou valores de storage. Outras origens foram filtradas. Não corresponde
a uma sessão de revisão nem observa raciocínio ou chamadas do backend.

A amostra tem 3.494 bytes e SHA256
`6cb1c3f716b086e8c1538912e685da51dcb65b4e7624fd8d7d2d03bbc95a98d1`.
URL original:
[chunk público de documentação](https://docs.coderabbit.ai/mintlify-assets/_next/static/chunks/ed45ca9fc4fad2f8.js?dpl=dpl_8AUWT8nFRmT7D3RgoEGEPVymfFQP).
REA visitou 544 nós AST, sem falha de parsing; o grafo semântico contém 295 nós,
426 relações e 105 incógnitas. A cobertura semântica geral é `unknown`.
Não foram encontrados módulos de revisão: é um chunk do site de documentação,
e não o motor do produto. Grafo recuperado não prova execução em runtime.

## O que as fontes sustentam

A [arquitetura publicada](https://docs.coderabbit.ai/overview/architecture)
descreve clone isolado do repositório, ferramentas estáticas, exploração de
contexto, funções especializadas de revisão/verificação e memória de feedback.
Isso sustenta a escolha de um pipeline com contexto e refutação. Não revela
prompts, pesos, algoritmos de recuperação, seleção de modelos ou limiares internos.

Também foram consultadas as
[instruções públicas da integração](https://github.com/coderabbitai/skills/blob/main/skills/code-review/SKILL.md).
Elas especificam saídas próprias e distinguem execução omitida de revisão concluída.
São documentação externa analisada como dado; não foram instaladas ou executadas.

## Um comportamento público concreto

Na [PR 92 de demonstração](https://github.com/coderabbitai/coderabbit-pr-review/pull/92),
a API retornou cinco comentários de `coderabbitai[bot]` ligados ao head
`0349d730a9f34013cfb1c59bbd61915f396f15e9`. Foram preservados somente metadados e
paráfrases em [public-pr-92.json](evidence/public-pr-92.json).

A amostra mostra localização por arquivo/linha, níveis de severidade,
explicação de impacto, sugestões e uma cadeia de análise exposta. Inclui
apontamentos de runtime além do propósito de demonstrar analisadores estáticos.
Não permite reconstruir o processo interno. Nenhum dos cinco apontamentos foi
adjudicado, nem a amostra foi usada para medir qualidade do Open Rabbit.

## Hipóteses de produto a testar

1. Mais contexto relevante deve aumentar a detecção de falhas entre arquivos.
2. Refutação antes de publicar deve reduzir ruído sem perder excessivamente recall.
3. Ferramentas estáticas podem aumentar cobertura em classes bem definidas.

São hipóteses de engenharia, não conclusões experimentais desta execução.
O [plano de avaliação](../docs/evaluation.md) define como medi-las.

## Reprodução

Use Chrome com perfil temporário e CDP loopback, sem sessão pessoal. Descubra o
ID atual com `rea list-browser-targets http://127.0.0.1:9231 --json`. O ID muda
entre sessões. As capturas originais preservam parâmetros e Evidence IDs.

```sh
rea inspect-web-page http://127.0.0.1:9231 TARGET_ID \
  --allowed-origins https://docs.coderabbit.ai \
  --observation-ms 1000 --json > /tmp/rea-page.json

rea analyze-javascript-application /tmp/open-rabbit-docs-sample \
  --source-path script.js \
  --max-semantic-nodes 2000 --max-semantic-nodes-per-file 2000 \
  --max-heap-mb 256 --analysis-timeout-ms 30000 \
  --json > /tmp/rea-static-sample.json
```

Baixe o chunk da URL registrada e verifique o hash antes da análise. Conteúdo
público pode mudar. O arquivo `evidence/manifest.json` identifica os artefatos
preservados. A análise completa que falhou não deve ser contada como concluída.
