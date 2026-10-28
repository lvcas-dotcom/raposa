> Demonstração sintética da baseline estática, sem modelo ou comparação de qualidade.

# Raposa

Base: `26c03bb82b389239f6af5a21e8459d8a825eb9f2`
Head: `61541090d4b1774863b997e962eab691ffe6f0cc`
Modo: static · duração: 14 ms

1 achado(s), 0 candidato(s) descartado(s). Resultado informativo; ausência de achados não representa aprovação.

## [high] Verificação padrão do certificado TLS desativada

client.go:4 · confiança 0.85 · ID b85ed8dbf170d720

Condição: Se esta configuração for usada em uma conexão sem verificação customizada equivalente, certificados não confiáveis serão aceitos.

Impacto: A conexão pode aceitar um servidor impostor; a explorabilidade depende do uso desta configuração.

Esperado: Manter a validação padrão ou verificar explicitamente a cadeia e o hostname em callback próprio.

Verificação: AST identifica crypto/tls.Config e literal true no diff, sem callback no mesmo literal; fluxo e atribuições posteriores em tls ainda exigem investigação.

Evidência em client.go:4:

     InsecureSkipVerify: true,

## Limitações

- Contexto auxiliar: arquivos Go do mesmo diretório, AGENTS ancestrais e README/go.mod raiz; sem grafo completo de consumidores ou execução de testes.

- Baseline estática: apenas tls.Config com InsecureSkipVerify literal true; sem resolução de tipos, fluxo ou prova de explorabilidade.

- Provedor semântico não configurado; apenas baseline estática Go, sem revisão geral de lógica.
