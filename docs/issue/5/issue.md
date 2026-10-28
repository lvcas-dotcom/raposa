# #5 local — chore: renomear o projeto para Raposa

## Contexto

O mantenedor renomeou o projeto. O nome antigo colide com um revisor de PR por
IA já publicado e sugere vínculo com a marca usada como comparador. CodeFox foi
considerado e descartado pelo mesmo motivo. A organização `raposa-labs` foi
criada para hospedar o repositório.

## O que fazer

- Renomear módulo Go, diretórios do produto, binário, textos do CLI e do relatório.
- Atualizar scripts, workspace, moon, documentação viva e exemplos.
- Adotar `img/banner.png` e `img/logo.png` no README e remover as imagens antigas.
- Apontar o repositório para `raposa-labs/raposa`.

## Fora de escopo

Reescrever registros históricos (`docs/issue/1–4`, `research/evidence`), mudar
comportamento ou contratos, renomear o repositório remoto ou a pasta local,
commit ou push.

## Critério de aceite

- `scripts/check.sh` e `scripts/build.sh` passam; `bin/raposa --help` exibe Raposa.
- Nenhuma ocorrência do nome antigo fora do histórico preservado.
- Links relativos dos documentos alterados resolvem.

## Registro

Issue local desta execução. Ainda não foi criado um número correspondente no GitHub.
