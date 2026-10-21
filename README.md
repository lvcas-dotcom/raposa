<p align="center">
  <img src="img/banner-open-rabbit-wide.png" alt="Banner do OpenRabbit" width="100%">
</p>

# <img src="img/icon-open-rabbit-v2.png" alt="" width="40" height="40"> Open Rabbit

Revisor de código local e laboratório de avaliação, inspirado nas capacidades
públicas do CodeRabbit. Objetivo: aproximar qualidade por evidência e medição.
O estado atual é uma primeira implementação funcional, sem equivalência demonstrada.
Projeto pessoal open source independente, preparado para hospedagem no GitHub.

## Começar

Requisitos: Git, Go 1.27.1; moon 2.3.4 é opcional para executar as mesmas tasks.
As versões ficam em `.prototools`. Não há dependências Go externas.

```sh
cd open-rabbit
proto install
sh scripts/check.sh
sh scripts/build.sh
./bin/open-rabbit --help
```

Recebe dois commits de um repositório existente. Não faz checkout e não executa
os builds, testes ou hooks do alvo. O relatório identifica os SHAs completos.
Use as referências desejadas do seu repositório:

```sh
./bin/open-rabbit review \
  --repo '/caminho/do/repositorio' \
  --base main --head minha-branch \
  --format md --output /tmp/review.md

./bin/open-rabbit prepare \
  --repo '/caminho/do/repositorio' \
  --base main --head minha-branch \
  --output /tmp/review-context.json
```

Sem provedor, `review` executa **somente uma regra estática Go**: identifica
`tls.Config` com `InsecureSkipVerify: true` introduzido no diff. É uma baseline
restrita, não uma revisão de lógica geral. Confiança é estimativa não calibrada.
Literais com callbacks customizados são omitidos pela baseline; tipos, fluxo e
atribuições posteriores ainda exigem investigação semântica.

## Revisão com modelo

O operador escolhe modelo e provedor. Não existe modelo padrão, download
automático ou uso implícito de credenciais. O port permite um processo externo
com JSON em stdin/stdout. Há um bridge para Ollama local, com duas passagens:
candidatos e tentativa de refutação.

1. Tenha Ollama rodando e um modelo escolhido disponível nele.
2. Copie `examples/provider.example.json` para `provider.local.json`.
3. Substitua o caminho absoluto e o modelo no arquivo.
4. Execute:

```sh
./bin/open-rabbit review \
  --repo '/caminho/do/repositorio' \
  --base main --head minha-branch \
  --provider-config provider.local.json \
  --format md --output /tmp/review-semantic.md
```

Um provedor de API também pode implementar o mesmo [protocolo](docs/provider-protocol.md).
`argv` é configuração confiável do operador: o processo roda com suas permissões.
O diretório temporário não é uma sandbox. O bridge Ollama aceita somente loopback,
recusa redirects e não solicita tools ao modelo.

A geração/verificação foi testada com processos de fixture. O bridge HTTP foi
testado com transporte controlado; nenhum modelo real nem API paga foi usado
nesta entrega. Falha do provedor encerra com erro e não produz relatório novo.

## O que está pronto

- CLI `prepare`, `review`, `evaluate`, JSON e Markdown.
- [Exemplo executado da baseline](examples/report.example.md) em fixture Git sintética.
- SHAs imutáveis, diff com renames, leitura de blobs e contexto auxiliar no head.
- Evidência literal e linha no diff verificadas; candidatos descartados rastreáveis.
- Segunda passagem por candidato e deduplicação por localização/classe.
- Limites de memória/saída/contexto e timeouts; cobertura omitida declarada.
- Avaliação por rótulos humanos: precisão, recall, duplicação e falsos negativos.
- Testes de domínio, arquitetura, Git real, protocolo, CLI e bridge HTTP.

Contexto auxiliar: Go do mesmo diretório, `AGENTS.md` ancestral e README/go.mod
raiz. Até 200 arquivos alterados, 128 KiB por blob, 40 arquivos auxiliares e
2 MiB de texto total. Binários, remoções e recortes omitidos aparecem nas limitações.
O prazo total padrão é dez minutos; cada chamada ao provedor tem até dois minutos.

Saída zero significa execução concluída, mesmo com achados. **Zero achados não
significa aprovação**, sobretudo com contexto parcial ou baseline estática.
O JSON contém trechos do repositório; mantenha os relatórios no ambiente autorizado.

## Pesquisa e qualidade

- [Evidências do REA e comportamento público](research/findings.md)
- [Comparação de capacidades e próximos marcos](docs/roadmap.md)
- [Protocolo de avaliação](docs/evaluation.md)
- [Decisões de arquitetura](docs/arquitetura/adr/0001-nucleo-go-e-adapters.md)
- [Metodologia de desenvolvimento](docs/processo/workflow.md)
- [Issue local #1](docs/issue/1/issue.md)

O projeto usa product-first, dependency rule, ports no chamador, ADRs, issues,
notas, validação e fecho. A CI está em `.github/workflows/ci.yml`.
Repositório: [lvcas-dotcom/open-rabbit](https://github.com/lvcas-dotcom/open-rabbit).
A licença de distribuição ainda será definida pelo mantenedor.
