# #4 local — feat: diff numerado no protocolo e janela de contexto do Ollama

## Contexto

O provedor recebe o conteúdo do head sem numeração e precisa contar a linha do
achado e copiar o trecho literal. Erros de contagem fazem o pipeline descartar o
candidato. O bridge Ollama não define a janela de contexto nem desativa o
truncamento do prompt; uma mudança grande pode ser revisada só em parte.

## O que fazer

- Acrescentar ao contrato v1, de forma aditiva, os hunks do diff com o número
  da linha no head, derivados do mesmo parse que define as linhas alteradas.
- Orientar o provedor a ancorar `Linha` e `Trecho` nesse campo.
- Exigir a janela de contexto no bridge Ollama, desativar truncamento e shift e
  recusar respostas que esgotaram a janela.
- Atualizar protocolo, exemplo e README.

## Fora de escopo

Protocolo v2, ampliação do contexto auxiliar, execução com modelo real,
medição de precisão, commit ou push.

## Critério de aceite

- Provedores v1 existentes continuam aceitos; o novo campo não altera respostas.
- Linhas alteradas permanecem iguais às da implementação anterior nos testes.
- Prompt acima da janela encerra com erro acionável, sem relatório parcial.
- `scripts/check.sh` passa; limites e o que não foi validado registrados no done.

## Registro

Issue local desta execução. Ainda não foi criado um número correspondente no GitHub.
