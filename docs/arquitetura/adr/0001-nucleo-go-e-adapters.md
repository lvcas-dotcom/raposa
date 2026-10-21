# ADR 0001 — Núcleo Go com integrações substituíveis

Status: adotada na implementação inicial em 2026-10-21; revisável antes da expansão.

## Contexto

O projeto adota product-first e `adapters → application → domain`, aproveitando
práticas conhecidas pelo mantenedor. A linguagem não determina a qualidade da
revisão. TypeScript também atende orquestração e integrações; Python favorece
experimentação. Distribuição simples, stdlib e controle de subprocessos
fundamentam a escolha inicial de Go.

## Decisão

Go 1.27.1 no núcleo e CLI. Ports `Repository`, `Analyzer`, `Reviewer` pertencem
ao application. Tipos de domínio não carregam tags de infraestrutura. Não há
pacote compartilhado antes de dois consumidores. Um módulo de produto em
`projects/open-rabbit/apps/cli` e workspace Go não vazio.

Provider por processo tem protocolo versionado e pode ser escrito em qualquer
linguagem. Um bridge Ollama demonstra a integração sem escolher modelo. REA é
instrumento de investigação independente do runtime do revisor.

## Consequências

Binários simples e baixa dependência inicial. Go é adequado para Git, subprocessos,
HTTP e controle de prazo. O ecossistema de avaliação e SDKs pode justificar um
adapter Python/TypeScript posteriormente; não há motivo para impor Go a ele.
Sem grafo de dependências completo nesta versão. Não há execução isolada de builds
nem afirmação de sandbox para o processo do provedor.

Referências: `os/exec` na [stdlib](https://pkg.go.dev/os/exec),
[SDKs oficiais MCP](https://modelcontextprotocol.io/docs/sdk),
[contexto registrado nas notas](../../issue/1/notes/notes-20261021T140601Z.md).
