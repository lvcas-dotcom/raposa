# Desenvolvimento

Entrada → trabalho → revisão sob demanda → fecho. Projeto pessoal, com issues,
pull requests e CI no GitHub.

1. Definir propósito, restrições e aceite. Ler/criar issue; criar branch dedicada.
   Worktree externa é recomendada em repositórios com mais de uma mudança simultânea.
2. Registrar análise, decisões e guardrails em `docs/issue/<N>/notes/notes-<UTC>.md`.
3. Implementar e executar checks que terminam. Preservar evidências e limites em done.
4. Revisão formal só quando solicitada; achado é informação, não ordem de execução.
5. Criar `docs/issue/<N>/done/done-<UTC>.md` antes de eventual commit.
6. Commit/push/PR dependem de pedido explícito; referências usam `Refs #N`.
   Fechar issue e fazer merge são ações do usuário.
   Não adicionar trailers `Co-authored-by`.

Adaptação de bootstrap: a issue inicial foi registrada localmente antes de existir
remoto. Preserve esse registro histórico; novas issues usam o GitHub e seus números
reais. O primeiro envio para `main` de um repositório vazio é permitido quando o
mantenedor autoriza a publicação inicial; as mudanças seguintes usam branches dedicadas.
