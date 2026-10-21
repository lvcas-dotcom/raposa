# Avaliação de qualidade

Meta: descobrir em quais classes o Open Rabbit aproxima ou supera um comparador,
e onde perde. CodeRabbit é um comparador, não um gabarito infalível.

## Protocolo

1. Fixar repo, base/head completos e política de diff. Não misturar revisão do
   working tree com revisão de commit. Registrar parâmetros, modelo, versão,
   tokens/custo disponíveis, duração e execução concluída/omitida/falha.
2. Selecionar 30–50 mudanças da stack escolhida: auth, fluxo de erro, contratos,
   concorrência, dados e casos sem defeito. Separar conjunto reservado por mudança,
   evitando variantes da mesma falha nos conjuntos de ajuste e teste.
3. Executar os revisores nos mesmos SHAs. Preservar resultados brutos no ambiente
   autorizado. Uma revisão histórica só entra se for possível identificar o
   snapshot realmente analisado.
4. Embaralhar os achados e ocultar o revisor da pessoa que adjudica. Pedir gatilho,
   localização, impacto e evidência. Medir defeitos que ambos deixaram passar com
   inspeção humana independente; a união de comentários não basta para recall.
5. Registrar um identificador estável por defeito e associar achados a ele.
   Duplicatas contam como ruído adicional, não como novos verdadeiros positivos.
6. Calcular métricas, intervalos de incerteza e erros por classe. Repetir apenas
   quando houver hipótese ou mudança que justifique nova medição.

## Calculadora implementada

```sh
./bin/open-rabbit evaluate \
  --input examples/evaluation.example.json \
  --output /tmp/metrics-demo.json
```

O exemplo é sintético, com hashes ilustrativos; seus números não medem produto.
O comando recebe a adjudicação humana. Ele não decide se um comentário é correto,
não verifica a existência dos SHAs e não executa os revisores.

Contrato de entrada: lista de `Caso` com `ID`, `Base`, `Head`, `Defeitos` e
`Revisoes`. Cada revisão tem `Revisor` e `Achados`; cada julgamento tem `ID`,
`Defeito` e `Valido`. Um achado inválido usa `Defeito` vazio. Defeito válido deve
estar no conjunto conhecido do caso. O mesmo conjunto de revisores é exigido em
todos os casos; cabe ao adjudicador atestar que cada resultado usou os SHAs do caso.

Saída: verdadeiros positivos, falsos positivos, falsos negativos, precisão e recall
microagregados por revisor. Divisão por zero vira `null`. Métrica sem rótulos
confiáveis não é evidência. Intervalos de confiança, custo, estabilidade entre
execuções e análise por severidade ainda precisam ser calculados no experimento.

## Estado da avaliação

- Pipeline, validações e cálculo: testes automatizados com fixtures.
- Modelo real, CodeRabbit em execução e conjunto adjudicado: pendentes.
- PR pública observada na pesquisa: não adjudicada, fora das métricas.
- Nenhum percentual de equivalência foi estimado.
