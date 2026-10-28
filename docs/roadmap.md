# Capacidades e próximos marcos

Comparação funcional baseada na documentação pública, sem comparação de acurácia.

| Capacidade | Raposa hoje | Próximo critério de aceite |
| --- | --- | --- |
| Contexto de repositório | snapshots e recorte auxiliar | importar contratos, callers e testes por símbolos, com evidência de recuperação |
| Análise estática | uma regra AST Go | adapters para analisadores versionados com resultados normalizados |
| Revisão semântica | port, protocolo e bridge local | escolher provedor e medir execução com modelo real |
| Verificação | segunda passagem e evidência literal | ablação com/sem refutação e adjudicação humana |
| Memória | não implementada | preferências explícitas, escopo de repo, origem e revogação |
| Execução isolada | sem executar alvo | sandbox real antes de permitir builds/testes não confiáveis |
| Integração Git | somente local | read-only GitHub, atualização incremental e idempotência |
| Publicação de comentários | não implementada | preview humano e autorização para publicar |
| Qualidade comparável | não demonstrada | conjunto reservado e comparativo nos mesmos SHAs |

## Sequência sugerida

1. Escolher provedor/modelo e stack inicial de avaliação. Registrar custo, versões
   e parâmetros. Executar casos controlados e confirmar que falhas permanecem visíveis (#1).
2. Coletar 30–50 mudanças reais com casos limpos e defeitos confirmados. Reservar
   parte sem usar no desenvolvimento de prompts e recuperação (#12).
3. Medir baseline, modelo sem refutação e pipeline completo. Investigar erros por
   classe antes de adicionar complexidade.
4. Ampliar contexto e ferramentas conforme os falsos negativos observados (#7, #9).
5. Integrar provedores Git e memória depois que houver utilidade mensurável (#10, #11).

Cada marco é uma issue no [board](https://github.com/users/lvcas-dotcom/projects/6),
com aceite verificável, notas antes de implementar e done com evidências. O
laboratório inicial está em `docs/issue/0/`; não representa conclusão do produto
completo ou da meta de qualidade.

## Referências que orientam o próximo recorte

[Fontes públicas úteis](../research/public-repositories.md) registra o exemplo de
contexto além do diff, regras estáticas e integrações, separando fatos e propostas.
O recorte prioritário é recuperar um nível de callers/callees por símbolo em Go,
com procedência e limites explícitos. Seu aceite inclui consumidores entre pacotes,
casos seguros, nomes ambíguos e budgets de contexto. Não está implementado.

A adoção depende da comparação com a recuperação atual nos mesmos SHAs, mantendo
provedor/configuração iguais e adjudicação humana. Ferramentas externas entram por
adapter e por classe validada; nenhuma fonte pública substitui essa medição.
