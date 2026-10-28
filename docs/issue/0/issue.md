# #1 — feat: estabelecer investigação reproduzível e revisor local

> Registro do bootstrap, anterior ao fluxo no GitHub. Foi criado como `docs/issue/1/`
> e movido para `docs/issue/0/` porque a #1 do GitHub é o piloto de revisão semântica.
> Caminhos citados dentro deste registro são os da época.

## Contexto

O usuário solicitou instalar e usar o REA para estudar o comportamento observável
do CodeRabbit e criar o Open Rabbit com arquitetura limpa e processo rastreável.
O nome e o destino foram definidos na conversa em 2026-10-21.
É um projeto pessoal independente, com destino GitHub e sem vínculo corporativo.

## O que fazer

- Instalar o REA com versão fixa em diretório independente.
- Registrar fontes, capturas, limitações e hipóteses sobre o CodeRabbit.
- Implementar CLI Go: preparar contexto de SHAs imutáveis, analisar alterações,
  integrar um provedor de revisão substituível, verificar achados e exportar relatório.
- Entregar baseline estática para Go e protocolo para revisão semântica por comando.
- Adotar product-first, ports no application, domínio puro, ADRs, moon, registros e CI.
- Fornecer avaliação por adjudicação humana e instruções reproduzíveis.
- Aplicar no README o banner e o logo fornecidos em `img/`, com links relativos
  compatíveis com o GitHub e descrição acessível.

## Fora de escopo

Backend proprietário, promessa de equivalência, contas autenticadas, publicação
de comentários, PR, interface web e execução de código do repositório analisado.
Integrações SaaS, memória e múltiplos repositórios são futuras. A autorização
posterior do mantenedor inclui configurar o remoto, criar commits e fazer o push inicial.

## Critério de aceite

- REA instalado e executado; evidência preservada com versão e hash.
- CLI executa em fixture Git real e aceita caminhos com espaços.
- Base/head são resolvidos para commits; alterações locais não contaminam o contexto.
- Achados fora do diff, evidência inválida e respostas malformadas são recusados.
- Falhas ou cobertura parcial ficam explícitas; zero achados não significa aprovação.
- Testes de comportamento, arquitetura, race, vet e build passam.
- Há instruções de uso, matriz de funcionalidades e avaliação sem alegar paridade.
- O README começa com o banner e apresenta o logo junto ao nome do projeto.
- Publicar a base inicial no GitHub pessoal com commits focados e sem coautoria.

## Referências

- https://github.com/morluto/rea
- https://docs.coderabbit.ai/overview/architecture
- https://docs.github.com/en/actions/tutorials/build-and-test-code/go

## Estado

Implementação inicial validada localmente; fechamento da issue fica com o mantenedor.
Destino definido: https://github.com/lvcas-dotcom/open-rabbit.
O mantenedor autorizou montar os commits sem coautoria e publicar a base inicial.
Licença e alteração de visibilidade ainda não foram definidas.
