# Fluxo de desenvolvimento

Norma de processo da Raposa, válida para pessoas e agentes. `AGENTS.md` aponta
para cá; regras de arquitetura e idioma ficam nele. Revisão em
[revisao.md](revisao.md); registro da issue em [docs/issue/README.md](../issue/README.md).

```
entrada  →  trabalho   →  revisão sob demanda  →  fecho
issue       implementa    finding em review/      done
branch      valida        desfechos               commit Refs #N
notes                                             PR
```

As fases não são passos com laço obrigatório. Só há duas ordens fixas: notes
antes de implementar e done antes do commit.

## Papéis

- **Responsável:** quem está atribuído à issue no board. Executa entrada, trabalho e fecho.
- **Mantenedores:** priorizam o board, revisam PRs, fazem merge e fecham issues.
  Quem abriu a PR não faz o merge da própria PR sem acordo com o outro mantenedor.
- Agente atua como responsável quando um mantenedor pede, dentro dos mesmos limites.

## Quando parar e perguntar

- **Pare** só no que é decisão de mantenedor: escopo, prioridade, trade-off de
  produto, alternativas igualmente defensáveis ou algo que a issue não responde
  e muda o resultado. Reporte as alternativas; não implemente por suposição.
- **Não pare** para pedir licença de continuar, confirmar o que a issue já diz ou
  narrar progresso. Decisão de engenharia é do responsável: decida e registre nas notes.
- **Gates, não perguntas:** o primeiro push da branch e a abertura da PR acontecem
  quando o responsável humano decide, ou quando um mantenedor pede ao agente.

## Fase 1 — Entrada

1. **Issue antes de qualquer edição**, sem exceção por tamanho.
   - Sem issue: brainstorm com propósito, restrições e critério de aceite; com o
     aceite de um mantenedor, abra pelo template (`.github/ISSUE_TEMPLATE/issue.yml`).
   - Seções: Contexto, O que fazer, Fora de escopo, Critério de aceite, Referências.
     Em pt-BR. Por CLI, envie o corpo por arquivo (`gh issue create --body-file`).
   - Ordem de verdade: descrição da issue > comentários > notes.
2. **Label de tipo** define o prefixo da branch:

   | Label | Branch |
   | --- | --- |
   | `type:feature`, `type:melhoria` | `feat/<N>-<slug>` |
   | `type:bug` | `fix/<N>-<slug>` |
   | `type:chore`, `type:tech-debt` | `chore/<N>-<slug>` |
   | `type:spike` | `spike/<N>-<slug>` |

3. **Branch** a partir de `origin/main` atualizado:
   `git fetch origin && git switch -c <tipo>/<N>-<slug> origin/main`.
   Trabalho que depende de outra branch ainda não mergeada parte dela e declara a
   dependência na issue e na PR.
4. **Worktree** opcional, fora do repositório, quando houver mais de uma issue em
   andamento: `git worktree add ../raposa-wt-<N> <branch>`.
5. **Board:** atribua-se à issue e mova para **Doing**.
6. **Notes** em `docs/issue/<N>/notes/notes-<ts>.md`, **antes** de implementar:
   o que a análise mediu, decisões (`[usuário]` nas respondidas por mantenedor),
   guardrails, perguntas pendentes e o que fica fora. Modelo no registro da issue.

## Fase 2 — Trabalho

- Regras de `AGENTS.md`: `adapters → application → domain`, ports no chamador,
  domínio só com stdlib, ferramental em inglês e domínio/docs/erros em pt-BR.
- Não executar código, build, testes ou hooks do repositório analisado pelo revisor.
- **Validação com prova:** comando e código de saída nas notes ou no done.
  Verde sem comando é registro falso. Runner em modo watch não conta.
- **Mínimo igual à CI:** `sh scripts/check.sh` e `sh scripts/build.sh`, depois da
  última edição. O que não rodou entra no registro com o motivo.
- Mudança estrutural vira ADR em `docs/arquitetura/adr/`.
- Ao concluir, **sugira** revisão. Revisão só roda quando um mantenedor pede.

## Fase 3 — Revisão sob demanda

Ver [revisao.md](revisao.md). Achado é informação; quem decide o desfecho é um
mantenedor. Achado válido que não cabe na issue vira issue nova, com o número
citado no finding.

## Fase 4 — Fecho

1. **Done** em `docs/issue/<N>/done/done-<ts>.md`, antes do commit: o que mudou,
   validação com códigos de saída, revisão, desvios e o que ficou de fora.
2. **Varredura de segredos** no registro e no diff: nenhum token, senha, chave,
   cookie, `.env`, DSN com credencial ou dado pessoal.
3. **Commit** em pt-BR, Conventional Commits, com `Refs #<N>` no corpo. Nunca
   `Closes`, `Fixes` ou `Resolves`: a issue é fechada por um mantenedor.
   Sem trailers de coautoria nem atribuição a ferramentas.
4. **Push e PR** pelo template, com `Refs #<N>`, validação e `## Limites declarados`.
   Base `main`, ou a branch de que a issue depende. Mova a issue para **Review**.
5. Correções pedidas na PR entram na mesma branch; agrupe-as num push por rodada.

## Pós-merge

Feito por um mantenedor, depois do merge:

1. Fechar a issue e mover para **Done**.
2. `git switch main && git pull --ff-only`.
3. `git worktree remove ../raposa-wt-<N>`, se houver, e `git branch -d <branch>`
   (nunca `-D`); `git remote prune origin`.
4. PR empilhada sobre a branch mergeada: conferir se o GitHub reapontou a base para `main`.

## O que nunca acontece

- Fechar issue pelo commit ou pela PR.
- Merge, push em `main`, apagar branch remota ou `--force` fora da branch da issue.
- Push ou PR por agente sem pedido de mantenedor.
- Segredo, dado pessoal ou saída bruta de ferramenta no registro.
- Revisão disparada sem pedido.
- Achado válido que fica só na conversa.
- Declarar paridade ou precisão sem avaliação nos mesmos SHAs e adjudicação humana.

## Bootstrap

As issues anteriores ao fluxo no GitHub foram registradas localmente. O bootstrap
está em `docs/issue/0/`; `docs/issue/2–5` foram espelhadas como #2–#5 com os
mesmos números. Novas issues usam apenas o GitHub.
