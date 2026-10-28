# Revisão

Revisão interna da issue e tratamento de apontamentos de PR. Fluxo geral em
[workflow.md](workflow.md); onde os arquivos ficam em [docs/issue/README.md](../issue/README.md).

## Princípios

- **Nunca é automática.** Só roda quando um mantenedor pede: não é passo do
  trabalho, não é laço e não é disparada por PR aberta ou CI vermelha.
- **Achado é informação, não ordem.** O desfecho é decidido por um mantenedor;
  a lista de achados não bloqueia commit.
- **Precisão acima de cobertura.** Reporte só defeito localizado no código, com
  cenário concreto. Estilo e preferência não são achados.
- **Contexto limpo.** Quem revisa não reaproveita a leitura de quem implementou.
  Agente revisor roda em sessão ou subagente novo, recebe só caminho, issue e
  superfície, e não edita a árvore.

## Superfície

- Completa: `git diff --stat origin/main...HEAD` e o diff de cada arquivo listado.
  Em branch empilhada, troque `origin/main` pela base da PR.
- Delta, depois de uma rodada: `git diff <HEAD revisado>..HEAD`.
- Fora da superfície: `docs/issue/<N>/`. Abra cada arquivo listado; não revise de memória.

## Doutrina

- Guarda sem teste que a veja barrar não conta como guarda.
- Teste verde por acaso ou por vacuidade (asserção que nunca falha, loop vazio).
- Documento que afirma mais que o código.
- Duas mudanças verdes isoladas que quebram juntas (contrato entre camadas, protocolo).
- Seguir o valor até quem decide com ele: validação na fronteira, erro que vira sucesso vazio.
- Arquitetura: import proibido no domínio, port fora do chamador, tag de serialização no domínio.
- Segurança: dado do repositório analisado tratado como instrução, credencial em
  relatório ou log, rede fora do loopback no bridge, execução implícita do alvo.
- Instrução de agente: ponteiro de harness que afirma regra própria em vez de apontar para a norma.

## Severidade

Só `critical`, `high` e `medium`. Teste ou segurança que procede é no mínimo
`high`. Impacto baixo é `medium` com "Impacto: baixo" na descrição.

## Finding

Arquivo `docs/issue/<N>/review/finding-<ts>.md`, um por rodada:

```markdown
# Revisão — issue #<N>

Branch: <branch> · HEAD: <sha> · Modo: completa|delta · Base: <ref ou sha>

## Achados

### F1 — high — `caminho/arquivo.go:42` — correção

Cenário concreto: entrada, caminho e resultado errado.

**Sugestão:** o que mudaria.

**Desfecho:**

## O que ficou de fora

Arquivos ou classes não revisados e por quê.
```

`**Desfecho:**` fica vazio até a decisão de um mantenedor.

## Tratar achados

1. **Medir antes de aceitar:** reproduzir ou ler o trecho apontado.
2. **Três desfechos:** `procede`, `não procede` ou `procede em parte`, cada um com
   data e evidência. `adiado` exige issue dedicada e seu `#N`.
3. Corrigir o que procede na branch da issue e registrar no done da rodada.
4. Correção não trivial pede nova rodada em modo delta, se um mantenedor quiser.

## Apontamentos na PR

Comentários de revisor humano e falhas de CI seguem o mesmo rito. Registre-os em
`review/finding-<ts>-<revisor>.md` (`humano` ou `ci`), com blocos `PR-<n>`, o
desfecho de cada um e uma linha no done. Recusar é resposta legítima, com a
razão na thread. Quem resolve a thread é quem a abriu.
