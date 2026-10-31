# #26 — docs: ADR do benchmark aberto como foco do projeto

## Contexto

Revisores de código com IA estão se multiplicando, mas não existe uma forma aberta, neutra e reproduzível de medir a qualidade deles. O núcleo da Raposa (SHAs imutáveis, evidência literal, adjudicação e `evaluate`) já foi desenhado para medir. Vale registrar o benchmark como foco do projeto e orientar o roadmap a partir dele.

## O que fazer

- ADR 0003: benchmark aberto de code review como foco principal, com o revisor da Raposa como implementação de referência medida nele.
- Atualizar `docs/roadmap.md` com a nova sequência e o vínculo das issues.
- Revisar na ADR 0002 a meta de comparação, para que o benchmark seja a régua comum a todos os revisores.

## Fora de escopo

Mudanças de código, README em inglês (issue própria) e escolha de licença (#17).

## Critério de aceite

- ADR com contexto, decisão, alternativas consideradas e consequências.
- O roadmap cita as issues do benchmark em ordem de dependência.

## Referências

`docs/roadmap.md`, `docs/arquitetura/adr/0002-evidencia-e-comparacao.md`, `docs/evaluation.md`.


## Registro

Issue [#26](https://github.com/lvcas-dotcom/raposa/issues/26), branch `chore/26-adr-benchmark`.
