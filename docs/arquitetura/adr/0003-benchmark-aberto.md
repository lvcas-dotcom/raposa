# ADR 0003 — Benchmark aberto de code review como foco

Status: adotada em 2026-10-31. Complementa a ADR 0002.

## Contexto

Revisores de código com IA se multiplicam, mas não há forma aberta, neutra e
reproduzível de medir a qualidade deles. Resultados publicados costumam vir de
quem vende o revisor, com casos e critérios fora do alcance de terceiros.

A Raposa já resolve referências para commits imutáveis, exige evidência literal,
separa execução concluída, omitida e com falha, e calcula métricas a partir de
adjudicação humana (`evaluate`). Esses são os requisitos de um benchmark.

## Decisão

O foco do projeto é um benchmark aberto de code review: dataset versionado de
mudanças reais com defeitos conhecidos e casos limpos, runner que executa
revisores por adapter nos mesmos SHAs, métricas com incerteza e resultados públicos.

O revisor da Raposa continua como implementação de referência. Sua evolução
(contexto por símbolos, ferramentas estáticas, refutação) é guiada e medida pelo
benchmark, nas mesmas condições dos demais revisores.

Rótulos têm origem declarada: `humano`, `szz` ou `sintetico`. Rótulo minerado
precisa de uma amostra adjudicada que estime seu erro. O conjunto reservado não
entra no ajuste de prompts, recuperação ou regras de nenhum revisor da Raposa.

Revisores de terceiros entram no ranking público somente quando seus termos de
uso permitem execução automatizada e publicação dos resultados.

## Alternativas consideradas

- **Revisor completo como produto principal:** compete por recursos com produtos
  mantidos por equipes dedicadas e não resolve a ausência de medição comparável.
- **Motor de contexto como servidor de ferramentas para agentes:** útil, mas com
  alternativas maduras; permanece como componente interno do revisor.

## Consequências

O valor do projeto passa a depender do dataset: volume, diversidade de classes e
linguagens, qualidade dos rótulos e facilidade de contribuição.

Mineração por SZZ escala a coleta, mas produz rótulos com ruído; o ranking
precisa exibir tamanho da amostra, intervalos e origem dos rótulos.

Cada execução de revisor tem custo de modelo ou de plano. Resultados brutos,
versões e parâmetros ficam registrados para que terceiros reproduzam.

Sequência e issues em [roadmap](../../roadmap.md).
