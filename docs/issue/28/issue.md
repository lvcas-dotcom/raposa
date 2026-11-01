# #28 — feat: formato versionado do caso de benchmark e validação

## Contexto

O `evaluate` recebe casos já adjudicados, sem um formato de dataset versionado. O benchmark precisa de casos autocontidos e reproduzíveis, que qualquer pessoa consiga contribuir e validar.

## O que fazer

- Formato versionado do caso: repositório público, base/head completos, licença da origem, linguagem, defeitos (id estável, arquivo, intervalo de linhas no head, classe, descrição) e origem do rótulo (`humano`, `szz`, `sintetico`).
- Marcar casos limpos (sem defeito) e o particionamento `desenvolvimento`/`reservado`.
- Comando `raposa benchmark validate`, que verifica o esquema e, com `--repo`, a existência dos SHAs e das linhas.
- Taxonomia inicial de classes alinhada a `docs/evaluation.md`.

## Fora de escopo

Mineração automática de casos, execução dos revisores e casamento de achados (issues próprias).

## Critério de aceite

- Exemplo em `examples/` validado pelo comando com exit 0; caso inválido falha com erro em pt-BR apontando o campo.
- Testes de domínio para as regras do esquema.
- Documentação do formato em `docs/benchmark/`.

## Referências

`projects/raposa/apps/cli/internal/domain/evaluation.go`, `docs/evaluation.md`, #12.


## Registro

Issue [#28](https://github.com/lvcas-dotcom/raposa/issues/28), branch `feat/28-formato-caso-benchmark`, a partir de `chore/26-adr-benchmark` (PR #38).
