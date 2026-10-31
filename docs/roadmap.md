# Roadmap

Foco: benchmark aberto de code review, com o revisor da Raposa como implementação
de referência ([ADR 0003](arquitetura/adr/0003-benchmark-aberto.md)).

## Benchmark

| Etapa | Issue | Depende de |
| --- | --- | --- |
| Decisão e roadmap | #26 | — |
| Formato versionado do caso e validação | #28 | #26 |
| Mineração de candidatos por SZZ | #29 | #28 |
| Amostra adjudicada por humanos e conjunto reservado | #12 | #28 |
| Runner nos mesmos SHAs | #30 | #28 |
| Piloto com modelo real no runner | #1 | #30 |
| Casamento achado→defeito e adjudicação cega | #31 | #30 |
| Termos de uso dos revisores | #36 | — |
| Adapters de revisores externos | #32 | #30, #36 |
| Métricas por classe, intervalos, custo e latência | #33 | #31 |
| Leaderboard estático | #34 | #33 |

Em paralelo: README em inglês (#27), releases e GitHub Action (#35), linguagens
da v0 e base multi-linguagem (#37), licença (#17).

## Revisor de referência

Comparação funcional baseada na documentação pública, sem comparação de acurácia.
Cada capacidade só é promovida após ganho medido no benchmark.

| Capacidade | Raposa hoje | Próximo critério de aceite |
| --- | --- | --- |
| Contexto de repositório | snapshots e recorte auxiliar | importar contratos, callers e testes por símbolos, com evidência de recuperação (#7) |
| Análise estática | uma regra AST Go | adapters para analisadores versionados com resultados normalizados (#9) |
| Revisão semântica | port, protocolo e bridge local | escolher provedor e medir execução com modelo real (#1) |
| Verificação | segunda passagem e evidência literal | ablação com/sem refutação e adjudicação humana |
| Memória | não implementada | preferências explícitas, escopo de repo, origem e revogação (#11) |
| Execução isolada | sem executar alvo | sandbox real antes de permitir builds/testes não confiáveis |
| Integração Git | somente local | read-only GitHub, atualização incremental e idempotência (#10) |
| Publicação de comentários | não implementada | preview humano e autorização para publicar (#10) |
| Qualidade comparável | não demonstrada | resultado no benchmark com conjunto reservado |

Ampliar contexto e ferramentas conforme os falsos negativos observados no
benchmark. Integrar provedores Git e memória depois que houver utilidade mensurável.

Cada marco é uma issue no [board](https://github.com/users/lvcas-dotcom/projects/6),
com aceite verificável, notas antes de implementar e done com evidências. O
laboratório inicial está em `docs/issue/0/`; não representa conclusão do produto
completo ou da meta de qualidade.

## Referências que orientam o próximo recorte

[Fontes públicas úteis](../research/public-repositories.md) registra o exemplo de
contexto além do diff, regras estáticas e integrações, separando fatos e propostas.
O recorte prioritário do revisor é recuperar um nível de callers/callees por
símbolo em Go, com procedência e limites explícitos. Seu aceite inclui consumidores
entre pacotes, casos seguros, nomes ambíguos e budgets de contexto. Não está implementado.

A adoção depende da comparação com a recuperação atual nos mesmos SHAs, mantendo
provedor/configuração iguais e adjudicação humana. Ferramentas externas entram por
adapter e por classe validada; nenhuma fonte pública substitui essa medição.
