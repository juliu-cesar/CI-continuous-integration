# Continuous Integration - CI

A Integração Contínua é uma prática de desenvolvimento de software que consiste na integração de modificações no código de forma continua e automatizada, permitindo a identificação rápida de problemas de integração, erros humanos de verificação e garantindo maior segurança. Com essa abordagem temos uma estabilidade do código e uma entrega contínua e consistente de funcionalidades aos usuários finais.

Os principais processos que são interessantes de serem executados seriam:

- **Execução** de testes
- **Linter**: verificar se o código esta dentro do padrão de formatação utilizado no projeto.
- **Verificações de segurança**
- **Geração de artefatos prontos para o processo de deploy**: é possível gerar um zip ou uma imagem docker para ser enviado ao servidor e quando oportuno executar a nova versão.
- **Identificação da proxima versão a ser gerada no software**: ao seguir os padrões do SemVer e do conventional commits, podemos utilizar programas que leem as modificações feitas e geram uma nova versão automaticamente.
- **Geração de tags e releases**

Essas são algumas opções, mas o processo de integração continua é muito customizável e permite uma grande liberdade ao usuário decidir o que sera feito.

## Status check

No repositório do GitHub temos algumas configurações que obrigam a Pull Request a passar por diversos processos como code review, CI, etc. antes de poder ser *mergeada* na branch principal. Ou seja, enquanto o Status Check da PR não for positivo, as modificações não podem ser aplicadas.

## Ferramentas populares

Existem diversas opções de ferramentas para se trabalhar com CI, abaixo temos algumas:

- [Jenkins](https://www.jenkins.io)
- [GitHub Actions](https://github.com/features/actions)
- [Circle CI](https://circleci.com)
- [AWS CodeBuild](https://aws.amazon.com/pt/codebuild/)
- [Azure Devops](https://azure.microsoft.com/pt-br/products/devops)
- [Google Cloud Build](https://azure.microsoft.com/pt-br/products/devops)
- [GitLab Pipelines CI](https://docs.gitlab.com/ee/ci/)

Para este projeto, vamos utilizar o GitHub Actions.

## O que é o GitHub Actions

GitHub Actions é uma ferramenta de automação integrada ao GitHub que permite criar fluxos de trabalho (workflows) automatizados para diversas tarefas. Esses workflows podem ser configurados para serem acionados pelos principais eventos do GitHUb, como push de código, fechar uma pull request, criar uma issue, entre outros. É uma ferramenta bastante poderosa e que pode ser utilizada para diversos outras funções além da integração continua.

### O que é uma Actions

É a ação que sera executada em um dos **Steps** de um **Job** em um **Workflow**, podendo ser criada do zero ou utilizada uma pre existente. Ela pode ser desenvolvida em **Javascript** ou **Docker image**.

### Funcionamento do Workflow

É um processo automatizado configurável composto por um ou mais jobs.

- Workflow
    - É um conjunto de processos definidos pelo usuário. Como por exemplo, radar os testes, efetuar o build, criar um artefato, etc.
    - Pode haver mais de um workflow por repositório.
    - Esse workflow é definidos em um arquivo `.yml` dentro da pasta `.github/workflows`.
    - Possui um ou mais Jobs
    - Pode ser iniciado por eventos do GitHub ou por agendamento.

Suponhamos que o workflow foi definido e um evento ativou ele, como funciona esse fluxo:

1. **Evento** : digamos que um evento de **push** que ativou a action.
2. **Filtros** : vamos filtrar para que somente seja executado quando for efetuado um push para a branch **main**.
3. **Ambiente** : selecionamos o ambiente em que ele sera executado, por exemplo o **ubuntu**.
4. **Ações** : agora sera executado os passos (**steps**) definidos no workflow, digamos algo como ele preciso executar um *composer* do php e um npm run.

Abaixo temos um exemplo de como ficaria esse fluxo no arquivo **yml**.

```yml
name: fluxo n1

on: 
  push:
    branches: 
      - main
jobs:
  build:
    runs-on: ubuntu

    steps:
    - uses: action/run-composer
    - run: npm run prod
```

Note que no steps temos o `uses` e o `run`, sendo que o primeiro executa uma Action criada por algum desenvolvedor ou por voce mesmo, e o segundo executa um código dentro do ambiente selecionado.

### Marketplace

Dentro do [GitHub Marketplace](https://github.com/marketplace?type=actions) temos diversas Action ja prontas para serem executadas.