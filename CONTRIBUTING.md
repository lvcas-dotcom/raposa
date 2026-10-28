# Contribuir

Guia de quem trabalha na Raposa. A norma do processo está em
[docs/processo/workflow.md](docs/processo/workflow.md); aqui fica o caminho prático.

## Primeiro acesso

Requisitos: Git e as versões de `.prototools` (Go; moon é opcional).

```sh
git clone https://github.com/lvcas-dotcom/raposa.git
cd raposa
proto install
sh scripts/check.sh
sh scripts/build.sh
./bin/raposa --help
```

Leia `AGENTS.md` (arquitetura e idioma), o workflow e as ADRs em
`docs/arquitetura/adr/` antes da primeira issue.

## Pegar uma issue

1. No [board](https://github.com/users/lvcas-dotcom/projects/6), escolha uma issue em
   **Todo** priorizada pelos mantenedores. Sem issue, abra uma pelo template e combine o escopo.
2. Atribua-se e mova para **Doing**.
3. Crie a branch pelo label de tipo: `feat/`, `fix/`, `chore/` ou `spike/` seguido de `<N>-<slug>`,
   a partir de `origin/main`.
4. Escreva `docs/issue/<N>/notes/notes-<ts>.md` antes de implementar
   (modelos em [docs/issue/README.md](docs/issue/README.md)).

## Implementar

- Código do produto em `projects/raposa/`. Preserve `adapters → application → domain`,
  ports no chamador e domínio sem dependências de infraestrutura.
- Ferramental em inglês; domínio, documentação, erros e cenários de teste em português.
- O revisor nunca executa código, build, testes ou hooks do repositório analisado.
- Mudança estrutural precisa de ADR. Novo modelo ou provedor descreve protocolo,
  limites de dados e como será avaliado.
- Resultado de fixture não mede precisão de modelo.

## Validar

`sh scripts/check.sh` (gofmt, vet, testes com `-race`, build) e `sh scripts/build.sh`
depois da última edição, os mesmos da CI. Registre comando e código de saída no done.

## Entregar

1. `docs/issue/<N>/done/done-<ts>.md` antes do commit.
2. Commits em pt-BR, Conventional Commits, com `Refs #<N>`; nunca `Closes #N`;
   sem trailers de coautoria.
3. PR pelo template, com validação e limites declarados. Mova a issue para **Review**.
4. Revisão interna e apontamentos da PR seguem [docs/processo/revisao.md](docs/processo/revisao.md).

## Quem decide

| Ação | Quem |
| --- | --- |
| Escopo, prioridade e trade-off de produto | mantenedores |
| Decisão de engenharia dentro da issue | responsável, registrada nas notes |
| Desfecho de achado de revisão | mantenedor |
| Merge, fechar issue, publicar release | mantenedor, que não seja o autor da PR sem acordo |

Discussões e arquivos do projeto se limitam à Raposa e a fontes públicas: nada de
código, nomes ou links de projetos privados ou corporativos.
