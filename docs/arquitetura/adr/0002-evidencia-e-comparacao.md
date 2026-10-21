# ADR 0002 — Evidência antes de publicação e medição antes de paridade

Status: adotada em 2026-10-21.

## Decisão

Resolver referências para commits antes de ler arquivos. Preservar texto de base
e head, linhas novas do diff e limitações. Não executar o alvo. Não confiar em
instruções dentro dos arquivos analisados.

Todo candidato precisa de localização no diff e evidência literal do snapshot,
condição concreta, impacto, contrato esperado e contraevidência. A revisão
semântica passa por tentativa de refutação. Validação literal demonstra
rastreabilidade, não demonstra correção da conclusão do modelo.

Sem provedor, o relatório identifica a baseline restrita. Falha do provedor não
vira sucesso vazio. Severidade e confiança são dimensões separadas. Duplicação
é contada como ruído na avaliação. Precisão/recall indefinidos são `null`.

## Consequências

Conservadorismo reduz cobertura, especialmente com contexto incompleto.
Verificação pelo mesmo modelo pode repetir erros e não é validação independente.
A meta de aproximar o CodeRabbit depende de conjunto reservado, rótulos humanos,
mesmos SHAs e custos/latências registrados, não de testes unitários verdes.
