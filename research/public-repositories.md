# Fontes públicas úteis para a Raposa

Consulta: 2026-10-21. Método: leitura de fontes públicas e comparação com o código
local, sem executar programas de terceiros ou modelos. Este recorte não é um
inventário exaustivo de todos os arquivos da organização.

## O que interessa

Os repositórios consultados oferecem demonstrações, regras e integrações. Não foi
identificado neles o motor completo de produção do CodeRabbit. Estudar um exemplo
publicado não comprova que o serviço comercial execute o mesmo pipeline ou modelo.

| Fonte | O que foi identificado | Aplicação na Raposa | Prioridade |
| --- | --- | --- | --- |
| [odsc-west-2025](https://github.com/coderabbitai/odsc-west-2025) | Revisor demonstrativo em Python; licença Apache-2.0 | Referência para buscar contexto além do diff | Alta |
| [ast-grep-essentials](https://github.com/coderabbitai/ast-grep-essentials) | Regras estáticas, utilitários e testes; Apache-2.0; arquivado na consulta | Candidatas a um adapter de análise, após validação por regra | Alta |
| [skills](https://github.com/coderabbitai/skills) | Integração de agentes com a CLI; MIT | Referência de fluxo e contratos de saída; a revisão é encaminhada à API do serviço | Média |
| [bitbucket](https://github.com/coderabbitai/bitbucket) | Cliente TypeScript para APIs do Bitbucket; Apache-2.0 | Referência de adapter Git se esse provedor entrar no escopo | Futura |
| [vscode-extension](https://github.com/coderabbitai/vscode-extension) | Repositório público de issues e discussões da extensão | Consultar problemas de integração; não tratá-lo como fonte do motor | Baixa |
| [coderabbit-pr-review](https://github.com/coderabbitai/coderabbit-pr-review) | Branch principal com README e PRs públicas de demonstração | Observar formato e comportamento dos comentários; não usar como gabarito | Apoio |

O estado de arquivamento e as licenças podem mudar. Conferir o arquivo de licença
na versão escolhida antes de incorporar qualquer código. Esta documentação não
incorpora código dessas fontes nem define a licença da Raposa.

## Demonstração de contexto além do diff

Versão examinada: `e1569721f2b851ea7b31d0c7d57615df4eed5da7`.
[Apresentação e escopo](https://github.com/coderabbitai/odsc-west-2025/blob/e1569721f2b851ea7b31d0c7d57615df4eed5da7/README.md).

O script usa AST Python para localizar símbolos alterados, formar relações de
chamada e coletar trechos de contexto. `one_hop_slice` seleciona arquivos de
impacto. O pacote distingue diff, código modificado e referências de callers e
callees; uma chamada ao modelo retorna JSON estruturado.
[Implementação observada](https://github.com/coderabbitai/odsc-west-2025/blob/e1569721f2b851ea7b31d0c7d57615df4eed5da7/review_demo.py).

Os limites do código lido são importantes: chamadas são associadas por `ast.Name`,
sem resolução completa de módulos ou métodos; nomes iguais podem ser confundidos.
A coleta de callers depende de alterações na linha da assinatura. O contexto
vem dos arquivos locais e de `HEAD~1`, sem o contrato de dois snapshots imutáveis.
Há exceções de leitura/parsing silenciadas. Não existe uma passagem independente
de refutação após a geração no fluxo observado.
[Pontos de coleta e execução](https://github.com/coderabbitai/odsc-west-2025/blob/e1569721f2b851ea7b31d0c7d57615df4eed5da7/review_demo.py).

Essas observações delimitam a demonstração; não são evidências de limitações do
produto comercial. Não foi medida a precisão do exemplo.

## Diferença em relação à Raposa atual

| Aspecto | Código local conferido | Trabalho necessário |
| --- | --- | --- |
| Identidade do código | Base/head resolvidos para commits | Preservar esse contrato ao indexar símbolos |
| Contexto auxiliar | Go do mesmo diretório, AGENTS ancestrais e README/go.mod raiz | Recuperar callers/callees entre pacotes |
| Símbolos | Não há índice de referências | Identidade por pacote, símbolo e posição, com ambiguidades explícitas |
| Verificação | Evidência literal, linha no diff e segunda passagem quando há provedor | Dar ao verificador o contexto ampliado e sua procedência |
| Ferramentas estáticas | Uma baseline AST Go | Normalizar regras externas por adapter e verificar sua utilidade |

Fontes locais: [snapshot Git](../projects/raposa/apps/cli/internal/adapters/git/repository.go),
[pipeline de revisão](../projects/raposa/apps/cli/internal/application/review.go)
e [baseline estática](../projects/raposa/apps/cli/internal/adapters/gostatic/analyzer.go).

## Proposta de aplicação

Esta seção é desenho proposto, ainda não implementado.

1. Mapear linhas alteradas para símbolos nos snapshots base e head. Começar com
   Go e preservar a distinção entre símbolo removido, renomeado e modificado.
2. Buscar um nível de callers e callees com identidade por pacote. Separar
   referências resolvidas, ambíguas e indisponíveis; correspondência só por nome
   não basta. Interfaces e despacho dinâmico precisam de limite declarado.
3. Incluir trechos e motivo de inclusão no contexto: contrato, consumidor ou teste.
   Resolver código do snapshot, sem executar builds ou scripts do alvo implicitamente.
4. Respeitar os budgets existentes. Se o contexto não couber, selecionar de forma
   determinística e registrar o que foi omitido, em vez de cortar silenciosamente.
5. Passar o mesmo pacote rastreável ao gerador e ao verificador. Contexto adicional
   não autoriza publicar achados preexistentes ou fora da superfície definida.

Exemplo de aceite proposto: alterar a assinatura de uma função em um pacote e
manter um consumidor incompatível em outro. A recuperação deve apresentar ambos,
com linhas e SHAs corretos; a revisão deve localizar o problema na mudança.
Também devem entrar no conjunto casos com consumidor atualizado, aliases,
métodos homônimos, mudança de comportamento sem assinatura nova e contexto truncado.

## Regras estáticas e integração

O [pacote de regras](https://github.com/coderabbitai/ast-grep-essentials) pode orientar
a seleção de classes e fixtures. A adoção deve fixar versão da ferramenta e das
regras, normalizar resultados no adapter e validar exemplos vulneráveis e seguros.
Status arquivado exige verificar manutenção e compatibilidade antes de depender dele.

As [integrações de agentes](https://github.com/coderabbitai/skills) ajudam a desenhar
contratos claros de execução concluída, omitida ou com erro. Esses estados devem
continuar distintos no produto; não é necessário copiar o formato da CLI comercial.

## Como decidir se valeu a pena

Comparar o pipeline atual com a recuperação por símbolos nos mesmos SHAs e com
o mesmo provedor/configuração. Registrar contexto recuperado, omissões, latência
e tokens quando disponíveis. Usar adjudicação humana para precisão, recall e
ruído, incluindo casos limpos e um conjunto reservado. Ver o
[protocolo de avaliação](../docs/evaluation.md).

Só promover a recuperação ampliada ou uma regra externa após observar ganho e
custos nos casos relevantes. A existência desses repositórios não estabelece
paridade de qualidade e não dispensa a medição da Raposa.
