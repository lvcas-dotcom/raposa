# Capacidades e próximos marcos

Comparação funcional baseada na documentação pública, sem comparação de acurácia.

| Capacidade | Open Rabbit hoje | Próximo critério de aceite |
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
   e parâmetros. Executar casos controlados e confirmar que falhas permanecem visíveis.
2. Coletar 30–50 mudanças reais com casos limpos e defeitos confirmados. Reservar
   parte sem usar no desenvolvimento de prompts e recuperação.
3. Medir baseline, modelo sem refutação e pipeline completo. Investigar erros por
   classe antes de adicionar complexidade.
4. Ampliar contexto e ferramentas conforme os falsos negativos observados.
5. Integrar provedores Git e memória depois que houver utilidade mensurável.

Cada marco vira issue com aceite verificável, notas antes de implementar e done
com evidências. A issue #1 entrega o laboratório inicial; não representa conclusão
do produto completo ou da meta de qualidade.
