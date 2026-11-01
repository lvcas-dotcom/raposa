# Formato do dataset do benchmark

Versão `1`. Um arquivo JSON por repositório de origem. Decisão em
[ADR 0003](../arquitetura/adr/0003-benchmark-aberto.md).

```sh
./bin/raposa benchmark validate --input dataset.json
./bin/raposa benchmark validate --input dataset.json --repo /caminho/do/clone
```

Sem `--repo`, valida apenas a estrutura. Com `--repo`, confere também que base e
head existem no clone e que cada defeito cabe nas linhas do arquivo no head.
Nenhum código do repositório é executado. Todos os problemas aparecem de uma vez,
com o campo apontado (`Casos[2].Defeitos[0].Classe: …`). Exemplo sintético em
[`examples/benchmark.example.json`](../../examples/benchmark.example.json).

## Campos

Nomes em pt-BR, iguais aos do domínio. Campo desconhecido é recusado.

| Campo | Regra |
| --- | --- |
| `Versao` | `"1"` |
| `Repositorio` | URL `https://` do repositório público |
| `Licenca` | licença do repositório de origem (SPDX recomendado) |
| `Casos[]` | pelo menos um |

Caso:

| Campo | Regra |
| --- | --- |
| `ID` | único no arquivo |
| `Linguagem` | linguagem principal da mudança |
| `Base`, `Head` | SHAs completos (40 ou 64 hex), diferentes entre si |
| `Particao` | `desenvolvimento` ou `reservado`; o reservado não entra em ajuste de nenhum revisor |
| `Limpo` | `true` exige `Defeitos` vazio; `false` exige pelo menos um defeito |
| `Defeitos[]` | defeitos conhecidos introduzidos entre base e head |

Defeito:

| Campo | Regra |
| --- | --- |
| `ID` | único no caso |
| `Caminho` | relativo à raiz, normalizado, com `/` |
| `Inicio`, `Fim` | linhas no head, `1 ≤ Inicio ≤ Fim` |
| `Classe` | `autenticacao`, `fluxo-de-erro`, `contrato`, `concorrencia`, `dados`, `seguranca`, `logica` |
| `Descricao` | gatilho e impacto em uma ou duas frases |
| `Origem` | `humano` (adjudicado), `szz` (minerado) ou `sintetico` (construído) |

## Limites

A validação não confirma que o defeito existe nem que a classe está correta: isso
é adjudicação. Rótulos `szz` têm ruído e dependem de amostra adjudicada (#12).
