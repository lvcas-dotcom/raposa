# Contribuir

Abra uma issue antes de iniciar uma mudança, descrevendo o problema, o desenho
proposto e um aceite verificável. Durante o bootstrap, as issues estão em
`docs/issue/`; depois da publicação, o espaço principal será GitHub Issues.

Use branch dedicada `<tipo>/<numero>-<slug>`. Código do produto fica em
`projects/raposa/`. Preserve `adapters → application → domain`, ports no
chamador e domínio sem dependências de infraestrutura. Ferramental em inglês;
domínio, documentação, erros e cenários de teste em português.

Execute `sh scripts/check.sh` e `sh scripts/build.sh`. Os testes devem terminar,
e qualquer alteração de comportamento deve ter evidência proporcional ao risco.
Não use resultados de fixtures para afirmar precisão de modelos.

Envie uma pull request com problema, resultado e validação, referenciando a issue
com `Refs #N`. Mantenha commits focados e sem trailers de coautoria automática.
Publicar, fazer merge e fechar issues são decisões do mantenedor.

Registre mudanças de arquitetura em ADRs. Propostas de novos modelos ou
provedores devem incluir protocolo, limites de dados e como serão avaliadas.
Discussões e arquivos do projeto devem se limitar à Raposa e fontes públicas.
