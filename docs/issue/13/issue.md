# Issue #13 — tornar configurável o timeout do provedor

## Contexto

O prazo por chamada ao provedor é fixo: 120 s no CLI e 110 s no bridge Ollama. Em CPU, modelos de 7B podem estourar o prazo mesmo em mudanças pequenas.

## O que fazer

- Flag de prazo por chamada no `review` e no `ollama-provider`, com o bridge abaixo do prazo do CLI.
- Validar valores positivos e manter o prazo total existente.

## Fora de escopo

Retentativa automática, streaming.

## Critério de aceite

- Prazo efetivo registrado no relatório.
- Timeout continua encerrando com erro, sem relatório parcial.
- Testes do limite e da validação.

## Referências

Adapter de comando do provedor, `cmd/ollama-provider`.


## Registro

Issue [#13](https://github.com/lvcas-dotcom/raposa/issues/13) no GitHub.
