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

## SonarQube

A proxima ferramenta que vamos trabalhar sera a [SonarQube](https://www.sonarsource.com/products/sonarqube/), que tem como a principal função avaliar a qualidade do código. Ela faz isso analisando todas as linhas de código do software para identificar problemas de segurança, bugs, vulnerabilidades, duplicações, cobertura de testes, entre outros aspectos, proporcionando uma visão detalhada sobre a saúde do código. Vejamos alguns pontos:

- **Análise de Qualidade de Código** : o SonarQube avalia a qualidade do código-fonte, identificando problemas como complexidade, vulnerabilidades, duplicações, bugs, más práticas de codificação, entre outros.

- **Feedback Rápido** : integra-se com pipelines de CI para fornecer feedback imediato aos desenvolvedores sobre a qualidade do código, permitindo a correção rápida de problemas.

- **Padronização e Boas Práticas** : ajuda a aplicar padrões de codificação e boas práticas para garantir um código mais limpo, legível e sustentável.

- **Melhoria Contínua** : facilita a identificação de áreas que precisam de melhorias, permitindo que a equipe foque em resolver problemas e aprimorar a qualidade do software.

- **Integração com CI/CD** : ao se integrar com pipelines de Integração Contínua (CI), o SonarQube automatiza a análise do código, possibilitando a detecção precoce de problemas e garantindo que o código submetido esteja alinhado com os padrões de qualidade definidos.

Outro ponto positivo é que os aspectos que definem o que é um código de qualidade são customizáveis, podendo ser alterar individualmente para cada projeto de acordo com a demanda.

O SonarQube pode ser instalado localmente em uma máquina pessoal ou em um servidor dedicado, sendo distribuído como um servidor web que oferece uma interface de usuário para a configuração, execução e análise das verificações de código.

### Rodando o SonarQube localmente

Na [documentação](https://docs.sonarsource.com/sonarqube/latest/try-out-sonarqube/#installing-a-local-instance-of-sonarqube) temos uma imagem Docker para rodar o SonarQube localmente, sem a necessidade de instalar na maquina, com o ponto negativo de que sempre que o container for reiniciado, os dados das analises e os projetos cadastrados serão perdidos. Lembrando que esta imagem é apenas para testar o SonarQube, em produção utilizaremos outros métodos.

```bash
docker run -d --name sonarqube -e SONAR_ES_BOOTSTRAP_CHECKS_DISABLE=true -p 9000:9000 sonarqube:latest
```

Sera disponibilizado a porta 9000 para acessar o SonarQube, com o login e senha sendo `admin`.

### Conceitos principais do SonarQube

Ao abrir a plataforma, podemos notar algumas opções na barra de navegação, e posteriormente ao criar um projeto, algumas dessas opções serão repetidas, criando assim **configurações globais** e **configurações locais**, sendo as que vamos explorar a seguir da primeira opção.

<img src="https://github.com/juliu-cesar/CI-continuous-integration/assets/121033909/6eaf530e-af18-4f88-9e12-4c0a7b6b4abe" height="50" />

As principais configurações que precisamos compreender são:

- **Rules** : as regras são as diretrizes que definem os padrões de codificação, elas determinarão o que é certo e o que é errado. São categorizadas em severidade como **Hight** > **Medium** > **Low**, e em tipos como **Bug**, **Vulnerability**, **Code Smell** e **Security Hotspot**. Elas são aplicadas durante a análise estática e cada linguagem suportada pelo SonarQube possui um conjunto específico de regras, podendo ser ativadas, desativadas e personalizadas conforme necessário para se adequar ao contexto do projeto.

<img src="https://github.com/juliu-cesar/CI-continuous-integration/assets/121033909/e153f36c-49ca-4da5-b735-bf9f7859ac6c" height="70" />

- **Quality Profiles** : os perfis de qualidade são conjuntos de regras configurados por linguagem de programação. Eles agrupam as regras que serão aplicadas durante a análise do código-fonte. O SonarQube vem com um perfil padrão para cada linguagem chamado de `Sonar way`, mas é possível criar perfis personalizados selecionando as regras desejadas, e controlando as métricas e o nível de rigor das verificações. Com isso atendendo melhor as necessidades específicas do projeto.

<img src="https://github.com/juliu-cesar/CI-continuous-integration/assets/121033909/7dd41502-abff-4514-b8ae-9a2d7f95058b" height="130" />

- **Quality Gates** : os Quality Gates ou em tradução livre *portões de qualidade*, são um conjunto de condições que determinam se um projeto pode ser considerado aceitável ou não. Eles são usados para definir critérios de qualidade que devem ser atendidos antes de permitir a integração do código e caso não, o código pode ser bloqueado até que os problemas sejam resolvidos. Como mencionado acima, estamos analisando as configurações globais, porem podemos criar quality gates personalizado por projeto. Algumas das métricas utilizadas são: 
  - Cobertura de testes
  - Número de bugs críticos
  - Vulnerabilidades 
  - Linhas duplicadas
  - Avaliação de manutenção
  - Confiabilidade de código

<img src="https://github.com/juliu-cesar/CI-continuous-integration/assets/121033909/e1d94959-cae8-4f2e-b4a1-16c5f10a095f" height="300" />

### Criando um projeto com SonarQube

Para testar o SonarQube localmente vamos criar um projeto simples em Go na pasta `SonarQube/go`, onde nele adicionamos o arquivo principal e o de teste. Após isso é preciso criar um projeto na plataforma do SonarQube, então vamos navegar para a aba `Projects` e clicar em `Create a local project`. Neles vamos inserir um **nome** (ci-go-project), uma **key** (ci-go-project) e o nome da **branch principal** (develop). Na seção seguinte vamos utilizar as configurações globais e pronto, projeto criado.

O proximo passo é criar o token para acessar esse projeto online, com o projeto selecionado vamos na aba `Projects` e clicamos na opção `Locally`, então escolhemos um **nome para o token** (ci-go-project-token) e o tempo para expirar a chave. Após isso escolhemos como vamos rodar a analise do projeto, onde o SonarQube ja possui integração com o `Maven`, `Gradle` e `.NET`, porem como estamos criando um projeto go vamos utilizar a opção `Other`, e em seguida o sistema operacional `Linux`. Para rodar a analise precisamos instalar o [sonar-scanner](#instalando-o-sonar-scanner) que sera responsável tanto pela analise como por enviar as informações para a plataforma onde esta o projeto online. Após escolher o sistema sera exibido um comando para rodar o scanner:

```bash
sonar-scanner \
  -Dsonar.projectKey=ci-go-project \
  -Dsonar.sources=. \
  -Dsonar.host.url=http://localhost:9000 \
  -Dsonar.token=sqp_71e7db520fb0f5c2692e875442df92c8ac8ee1c5
```

- projectKey : a key do projeto.
- sources : o caminho onde estão os arquivos que serão analisados. Passamos apenas um `.` pois estamos considerando que o comando foi executado com o shell estando na pasta do projeto.
- host.url : local onde se encontra a plataforma do SonarQube, por padrão ela é disponibilizada na porta 9000.
- token : o token que acabamos de gerar.

Ao executar o comando podemos notar que dentro do projeto é adicionado uma pasta chamada do scanner, e na plataforma do SonarQube indo em `Overview` agora temos algumas informações do projeto, e muito provavelmente ele passou no teste por possuir pouquíssimas linhas. Outro detalhe é que na aba `Code` podemos ver todo nosso código que foi enviado, e também as funções que foram cobertas com testes ou não.

Apesar de funcionar com o comando acima, o mais indicado é criar um arquivo com as propriedades do projeto, o que ira facilitar até mesmo para informar o status de cobertura do código.

### Cobertura de código com SonarQube

Algo muito importante sobre o SonarQube é que ele não roda os testes para saber quais funções estão cobertas ou não, o que ele faz é pegar essas informações de um arquivo que deve ser disponibilizado e enviar para a plataforma. Cada linguagem possui uma forma ou formas para executar esse processo, sendo necessario procurar na [documentação](https://docs.sonarsource.com/sonarqube/9.9/analyzing-source-code/test-coverage/overview/) os tipos de arquivos aceitos. Para o projeto go temos o seguinte comando:

```bash
go test -coverprofile=coverage.out
```

Com a opção `coverprofile` o go ira criar um arquivo (coverage.out) com o resultado dos testes.

Agora o proximo passo é criar o arquivo `sonar-project.properties`, onde iremos informar as propriedades necessários para o projeto, vejamos como configura-lo:

```properties
sonar.projectKey=ci-go-project
sonar.sources=.
sonar.host.url=http://localhost:9000
sonar.token=sqp_9b9851ceba256bd6304d51766e274e0dd9eb5431

sonar.tests=.
sonar.test.inclusions=**/*_test.go
sonar.exclusions=**/*_test.go
sonar.go.coverage.reportPaths=coverage.out
```

No começo de arquivo temos as quatro configurações vistas anteriormente, porem ao final do arquivo temos outras quatro:

- tests : a pasta onde se encontra os testes, e assim como em *sources* eles estão na mesma pasta onde se encontra o arquivo.
- test.inclusions : quais os arquivos responsáveis pelos testes, no caso todos que terminem com "_test.go".
- exclusions : quais os arquivos que dever ser excluídos de cobertura, no caso todos os arquivos de teste.
- go.coverage.reportPaths : a localização do arquivo com o resultado da cobertura de testes, como estão na mesma pasta informamos apenas o nome.

Apos configurar o arquivo de propriedades, basta rodar o comando `sonar-scanner` na mesma pasta para executar a analise do código, com isso a informação de `Coverage` na plataforma deve ser atualizada. Porem mesmo que a cobertura esteja abaixo do padrão exigido pelo SonarQube, o programa provavelmente passou, isso acontece por que existem poucas linhas de códigos, então a analise não o bloqueia. Para testar o funcionamento completo da analise, vamos adicionar mais códigos ao arquivo `sum.go`, com isso a cobertura deve ficar ainda menor e o programa se barrado.

## Instalando o Sonar Scanner

Abaixo temos o paço a paço para instalar o **sonar-scanner** no Ubuntu, mas para acessar a documentação original segue o [link](https://docs.sonarsource.com/sonarqube/10.3/analyzing-source-code/scanners/sonarscanner/).

1. Instalar os pacotes `wget`, `unzip` e `nodejs`

```bash
apt-get update
apt-get install unzip wget nodejs
```

2. Baixar e extrair o software

```bash
wget https://binaries.sonarsource.com/Distribution/sonar-scanner-cli/sonar-scanner-cli-5.0.1.3006-linux.zip
unzip sonar-scanner-cli-5.0.1.3006-linux.zip
mv sonar-scanner-cli-5.0.1.3006-linux.zip /opt/sonar-scanner
```

> O link pode ser obtido no site do paragrafo acima, passando o mouse sobre a opção `Linux 64-bit`.

3. Exportar a variável PATH no arquivo do shell (.bash, .zshrc, etc)

```bash
export PATH="$PATH:/opt/sonar-scanner/bin"
```

