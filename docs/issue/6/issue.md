# #6 — chore: adotar fluxo de issues, board e templates

## Contexto

O projeto passa a ter dois mantenedores. O fluxo atual cabe em `docs/processo/workflow.md` (cinco passos) e em `AGENTS.md`, com issues locais em `docs/issue/<N>`. Não há board, labels de tipo nem templates; um colaborador novo não tem um guia único para pegar, executar e entregar uma issue.

## O que fazer

- Labels `type:feature`, `type:melhoria`, `type:bug`, `type:chore`, `type:tech-debt`, `type:spike`; o tipo define o prefixo da branch.
- Board do projeto com Status `Todo → Doing → Review → Done`; o status vive no board, não em label.
- Espelhar os registros locais #2–#5 como issues e criar o backlog a partir do roadmap.
- `docs/processo/workflow.md`: quatro fases (entrada, trabalho, revisão sob demanda, fecho), quando parar e perguntar, pós-merge e o que nunca acontece.
- `docs/processo/revisao.md`: princípios, superfície, doutrina, severidade, formato do finding e tratamento com três desfechos.
- `docs/issue/README.md`: layout do registro (`notes/`, `review/`, `done/`), `<ts>` medido e modelos.
- Templates de issue e PR em `.github/`; `CONTRIBUTING.md` como guia do colaborador; `AGENTS.md` apontando para os documentos sem duplicar a norma.
- Mover o registro do bootstrap de `docs/issue/1` para `docs/issue/0`, já que a #1 do GitHub é o piloto.

## Fora de escopo

Comandos de agente (#18), automações de CI para o board, mudanças de código ou de comportamento, fechar issues.

## Critério de aceite

- Issues #1–#18 com label de tipo e presentes no board com o Status correto.
- Documentos navegáveis, links relativos resolvidos e templates YAML válidos.
- Um colaborador consegue, só com `CONTRIBUTING.md` e `docs/processo/`, ir do board a uma PR com registro completo.

## Referências

`AGENTS.md`, `docs/processo/workflow.md`, `CONTRIBUTING.md`, `.github/pull_request_template.md`.

## Registro

Issue [#6](https://github.com/lvcas-dotcom/raposa/issues/6), branch `chore/6-fluxo-issues-board`.
