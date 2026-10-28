# Registro da issue

Cada issue tem `docs/issue/<N>/`, com o número do GitHub. É o planejamento e o
histórico da execução; a descrição da issue continua sendo a fonte de verdade
do escopo. Fluxo em [workflow.md](../processo/workflow.md).

## Layout

| Arquivo | Quem escreve | Quando |
| --- | --- | --- |
| `issue.md` | responsável | entrada; cópia da descrição com a seção Registro |
| `notes/notes-<ts>.md` | responsável | antes de implementar |
| `review/finding-<ts>.md` | revisor | a cada rodada de revisão interna |
| `review/finding-<ts>-<revisor>.md` | responsável | a cada rodada de apontamentos na PR |
| `done/done-<ts>.md` | responsável | antes do commit; atualizado se o trabalho mudar |

- `<ts>` é medido com `date -u +%Y%m%dT%H%M%SZ`, nunca escrito à mão.
- Saída de ferramenta (JSON, log, relatório gerado, screenshot) não entra:
  vai para `reports/`, que é ignorado pelo Git. O registro cita o comando e o resultado.
- Nunca entra token, senha, chave, cookie, `.env`, DSN com credencial ou dado
  pessoal. Achado de segurança cita `arquivo:linha` e a classe, nunca o valor.
- Registros fechados não são reescritos. Correção posterior entra como nota datada.

## Numeração

- `0/`: bootstrap do projeto, anterior ao GitHub (criado como `1/`).
- `2/`–`5/`: registros locais espelhados no GitHub como #2–#5.
- Demais: número da issue no GitHub. A #1 é o piloto de revisão semântica.

## Modelos

### notes

```markdown
# Notas — issue #<N> — <assunto>

## O que a análise mediu

Fatos verificados, com `arquivo:linha` ou comando.

## Decisões

- [usuário] Decisão respondida por mantenedor.
- Decisão de engenharia e o motivo.

## Guardrails

O que não pode acontecer nesta execução.

## Perguntas pendentes

## O que fica fora
```

### finding

Formato em [revisao.md](../processo/revisao.md#finding).

### done

```markdown
# Entrega — issue #<N> — <assunto>

## Resultado

O que mudou, por arquivo ou capacidade.

## Verificação

| Comando | Resultado |
| --- | --- |
| `sh scripts/check.sh` | exit 0 |

## Revisão

Rodadas e desfechos, ou "não solicitada".

## Desvios

Diferenças em relação à issue ou às notes, com motivo.

## Limites e estado

O que não foi validado, o que ficou de fora e pendências de mantenedor.
```
