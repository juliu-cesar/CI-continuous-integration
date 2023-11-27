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

Lembrando que o GitHub disponibiliza 2000 minutos (no momento da criação deste documento) de maquina por mes para repositório públicos. Logo é interessante não ficar rodando diversas actions quando não for necessario.

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
    branches: [ main ]
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

## Criando uma Action

Antes de trabalharmos com Actions e PRs, vamos criar um exemplo apenas para ver o funcionamento de uma action. Primeiramente adicionamos um programa muito simples em Go para efetuar uma soma de dois números, e então criamos um outro para testar essa função, e serão eles que iremos executar dentro o ubuntu na Action. Apos isso, na pasta `.github/workflows` vamos adicionar o arquivo `ci-test-go.yaml`.

```yaml
name: continuous-integration-go

on: [push]
jobs: 
  check-application:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-go@v2
        with:
          go-version: 1.15
      - run: go test
      - run: go run math.go
```

No arquivo acima definimos que sempre que houver um push, sera executado o job `check-application`, que por sua vez ira rodar uma maquina ubuntu. Dentro dessa maquina copiamos os código para dentro dela com o `checkout@v2`, e criamos um setup para o go com o `setup-go@v2` passando a versão 1.15. Por fim as ultimas duas linhas executam os testes em go.

Com isso temos a Action configurada e sempre que houver um push ela sera executada, bastando ir na aba `Actions` do repositório para verificar todos os passos que foram executados.

## Adicionando Status Check para as PRs

Para que a verificação da Action tenha alguma validade precisamos aplica-la a uma PR, onde somente sera liberado o merge caso o status check tenha sido aprovado. Para fazer tal configuração vamos em `settings`, `branches` e `add rule`. As regras que vamos aplicar são:

- `Require status checks to pass before merging` : isso garante que só sera possível efetuar um merge caso o status check for aprovado.
  - `Require branches to be up to date before merging` : a branch precisa ser a mais atual para poder ser feito o merge.
  - `Status checks that are required` : para esta opção vamos passar o nome do job que foi definido na Action, no caso o `check-application`. Caso o GitHub não sugerir ela, basta pesquisar no box de texto acima desta opção.

Agora a Action precisa ser ativa apenas quando for criado uma nova PR, logo vamos modificar o arquivo *yaml*, ou criar um novo.

```yaml
name: ci-prs

on: 
  pull_request:
    branches: 
      - develop

jobs: 
  check-application:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-go@v2
        with:
          go-version: 1.15
      - run: go test
      - run: go run math.go
```

Alteramos o evento que ativa a action, porem anteriormente utilizamos o `[push]` para informar o tipo do evento, mas dessa vez passamos o evento como um campo, com isso podemos informar em qual branch esse evento sera ativo, que no caso é somente na `develop`.

## Múltiplas ambientes com Strategy Matrix

Outra parte importante sobre a fase de testes da Action, é que pode ser necessario testar o software em múltiplos ambientes e versões de uma mesma linguagem, e para essa tarefa temos o `Strategy` do Actions. Nele podemos passar uma matriz especificando as **versões** e até mesmo o **sistema** onde deve ser executado. Levando isso em consideração vamos modificar o arquivo *yaml* da seguinte forma:

```yaml
name: ci-multi-version

on: 
  pull_request:
    branches: 
      - develop

jobs: 
  check-application:
    strategy:
      matrix:
        version: [1.14, 1.15]
        os: [ubuntu-latest, windows-latest]
    runs-on: ${{ matrix.os }}
    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-go@v2
        with:
          go-version: ${{ matrix.version }}
      - run: go test
      - run: go run math.go
```

Dentro do Job criamos a estrategia `strategy` e passamos uma matriz `matrix`. Dentro dessa matriz poderíamos passar apenas as versões, ou apenas os sistemas operacionais, mas vamos utilizar os dois, dessa forma ao final sera executado 4 jobs. Para indicar a Action qual sistema utilizar, é preciso ir em `runs-on` e passar o campo da matriz responsável pelos sistemas, no caso `${{ matrix.os }}`. O mesmo vale para as versões `${{ matrix.version }}`.

## CI com Docker

Uma funcionalidade muito importante com as Actions, é que podemos construir imagens docker dentro dela. Onde setamos um ambiente que rode Docker e então efetuamos o build da imagem, e até mesmo o push para o Docker Hub. 

Primeiramente vamos construir uma imagem docker do Go, onde ela apenas ira fazer o build do arquivo `math.go` e executar ele. Dentro da raiz do projeto vamos adicionar o `Dockerfile`:

```dockerfile
FROM golang:1.19

WORKDIR /app

COPY . .

RUN go build -o math

CMD [ "./math" ]
```

Dentro do Marketplace temos uma action pronta que faz o build e o push de imagens docker, [Build and push Docker images](https://github.com/marketplace/actions/build-and-push-docker-images). Na pagina da Action temos uma pequena documentação com alguns códigos de exemplos que vamos utilizar. No primeiro momento iremos **apenas efetuar o build** da imagem Docker.

```yaml
# Código omitido ...
jobs: 
  check-application-and-build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-go@v2
        with:
          go-version: 1.15
      - run: go test
      - run: go run math.go

      - name: Set up QEMU
        uses: docker/setup-qemu-action@v3

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v3

      - name: Build and push
        id: docker_build
        uses: docker/build-push-action@v5
        with:
          push: false
          tags: juliucesar/ci-action-golang:latest
```

Começamos definindo um ambiente para rodar o Docker com a action **QEMU** que permite rodar o docker em diversas arquiteturas. Em seguida executamos outra action **Docker Buildx**, que permite executar o build da imagem nos passos seguintes. Por fim rodamos o build imagem, passando o nome da imagem como `ci-action-golang`, este sera o nome da imagem na maquina temporária da action, ela sera importante quando formos fazer o push da imagem.

### Push da imagem

Para efetuar o push no Docker Hub é preciso estar logado na maquina, e para isso vamos utilizar uma outra action que também consta na documentação vista no passo anterior.

```yaml
jobs: 
  check-application-and-build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-go@v2
        with:
          go-version: 1.15
      - run: go test
      - run: go run math.go

      - name: Set up QEMU
        uses: docker/setup-qemu-action@v3

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v3

      - name: Login to Docker Hub
        uses: docker/login-action@v3
        with:
          username: ${{ secrets.DOCKERHUB_USERNAME }}
          password: ${{ secrets.DOCKERHUB_TOKEN }}

      - name: Build and push
        id: docker_build
        uses: docker/build-push-action@v5
        with:
          push: true
          tags: juliucesar/ci-action-golang:latest
```

Note que no campo de usuário e senha, utilizamos uma **secret** do github, e isso é necessario por questões de segurança, uma vez que quem obter essas informações terão acesso a canta do Docker. Vejamos os passos para configurar essas secrets:

- Na aba `Settings` do repositório, vamos em `Secrets and variables`, `Actions` e por fim em `New repository secret`, 

Nesta janela vamos configurar as duas variáveis, a primeira sendo a `DOCKERHUB_USERNAME` onde informaremos o usuário, e a segunda a `DOCKERHUB_TOKEN` onde informaremos o token do Docker. Para entender mais sobre o token de acesso e como cria-lo, temos a documentação do Docker [Create and manage access tokens](https://docs.docker.com/security/for-developers/access-tokens/). Mas para resumir como criar o token:

- Efetue o login no Docker Hub, selecione o nome do usuário no canto superior direito e clique na opção `Account Settings`, depois em `Security` e por fim `New Access Token`.

Apos essas configurações, sempre que for criada uma nova PR, sera feito o push de uma imagem docker para o Docker Hub. 

Lembrando que este é apenas um exemplo, e podemos alterar o evento que ativa essa action, para por exemplo, somente fazer o push quando a for feito o merge na branch main.
