# Raposa

Fonte normativa para agentes deste projeto pessoal open source, com destino GitHub.

- Leia a issue completa antes de editar. Sem remoto, issue em `docs/issue/<N>/issue.md`.
  Brainstorm e aceite precedem a issue quando o escopo ainda não foi definido.
- Branch `<tipo>/<N>-<slug>`; notas antes da implementação; done com comandos e limites.
- Um agente por padrão. Modelo e esforço são escolhas do usuário.
- Código de produto em `projects/raposa/`; sem libs compartilhadas com um consumidor.
- Dependências: adapters → application → domain. Ports pertencem ao application.
  Domain usa apenas stdlib, sem infraestrutura nem tags de serialização.
- Ferramental em inglês; domínio, cenários, docs e erros em pt-BR.
- ADR para decisão estrutural. Não duplicar normas em ponteiros de harness.
- Validar o comportamento proporcional ao risco. Não executar código do alvo implicitamente.
- Revisões formais sob demanda. Validação da implementação não equivale a revisão independente.
- Não declarar paridade ou precisão sem avaliação nos mesmos SHAs e adjudicação humana.
- Commit, push, PR e publicações externas só quando pedidos. Nunca fechar issue ou fazer merge.
- Commits sem trailers de coautoria. Não incluir identidades, links ou arquivos corporativos.
- Documentação e metadados descrevem apenas este projeto e suas fontes públicas.

Procedimento: `docs/processo/workflow.md`.
