![brcepgopher](docs/img/gopher.png)

# brcep 

[![build status](https://img.shields.io/travis/leogregianin/brcep/master.svg?style=flat-square)](https://travis-ci.org/leogregianin/brcep) [![Go Report Card](https://goreportcard.com/badge/github.com/leogregianin/brcep)](https://goreportcard.com/report/github.com/leogregianin/brcep) [![github closed issues](https://img.shields.io/github/issues-closed-raw/leogregianin/brcep.svg?style=flat-square)](https://github.com/leogregianin/brcep/issues?q=is%3Aissue+is%3Aclosed) [![codecov](https://codecov.io/gh/leogregianin/brcep/branch/master/graph/badge.svg)](https://codecov.io/gh/leogregianin/brcep)

API para acesso a informações dos CEPs do Brasil. A ideia central é não ficar dependente de uma API específica, e sim, ter a facilidade de acessar a __brcep__ e ela se encarrega em consultar diversas fontes e lhe devolver as informações do CEP de forma rápida e fácil.

O projeto __brcep__ faz consultas às APIs [ViaCEP](http://viacep.com.br) e [CEPAberto](http://cepaberto.com).

![brcep](docs/img/brcep.png)

### Sidecar Pattern

A idea desse projeto é que você utilize a imagem do Docker como um [sidecar](https://dzone.com/articles/sidecar-design-pattern-in-your-microservices-ecosy-1) para sua aplicação atual. Esse projeto não é uma biblioteca pra consumir APIs, mais sim um server que deve rodar ao lado (dai o sidecar) da sua aplicação atual, e quando você precisar fazer a requisição de um cep, você fará a requisição para o endpoint do sidecar e não diretamente para uma API. Dessa forma você ganha a vantagem de um middleware que fará o uso correto de multiplas APIs. 

Considere o docker-compose abaixo para entender melhor:

```yaml
version: '2.1'

services:
  myapp:
    container_name: myapp
    build:
      context: .
      dockerfile: Dockerfile
    ports:
      - "3000:3000"
  brcep:
    image: brcep:latest
    ports:
      - "8000:8000"
    links:
      - myapp
    container_name: brcep
    environment:
      - PORT=8000
```

A idea é que sua aplicação rode na porta 3000 e o brcep rode na porta 8000, considerando os multiplos [exemplos](./docs/exemplos.md) que temos definido, substituindo a URL por http://brcep/78048000/json. Sendo assim a sua aplicação agora faz o consumo de forma transparente de diversas APIs através do brcep. 


Tópicos
=================

  * [Exemplo API](#exemplo-api)
  	* [Acesso](#acesso)
  	* [Response](#response)
  * [Execução](#execução)
  	* [Configuração do ambiente](#configurao-do-ambiente)
  	* [Executar com Docker](#executar-com-docker)
  	* [Executando localmente](#executando-localmente)
  	* [Executando testes](#executando-testes)
  * [Exemplos](./docs/exemplos.md)
	* [Curl](./docs/exemplos.md#curl)
	* [Javascript](./docs/exemplos.md#javascript)
	* [Python](./docs/exemplos.md#python)
	* [Golang](./docs/exemplos.md#golang)
	* [Ruby](./docs/exemplos.md#ruby)
	* [PHP](./docs/exemplos.md#php)
	* [Java](./docs/exemplos.md#java)
	* [C#](./docs/exemplos.md#c-sharp)
	* [Delphi](./docs/exemplos.md#delphi)
  * [Licença de uso](#licença-de-uso)


## Exemplo API

### Acesso

Para facilitar a visualização do que esperar deste projeto, a versão atual está disponível para visualização dos dados em [https://brcep-qnlohrjtbl.now.sh/78048000/json](https://brcep-qnlohrjtbl.now.sh/78048000/json).

### Response

```json
{
  "cep": "78048000",
  "endereco": "Avenida Miguel Sutil, de 5799/5800 a 7887/7888",
  "bairro": "Consil",
  "complemento": "",
  "cidade": "Cuiabá",
  "uf": "MT",
  "latitude": "-15.5786867",
  "longitude": "-56.0952081",
  "ddd": "",
  "unidade": "",
  "ibge": "5103403"
}
```

* O campo "CEP" retorna somente números.
* Os campos "complemento", "latitude" e "longitude" podem retornar em branco dependendo da API consultada.
* Os demais campos sempre retornarão valores.

## Execução

### Configuração do ambiente

* A API do CEPAberto necessita do token de autorização e a API do ViaCEP não necessita de token.
* Renomear o arquivo .env.example para .env e incluir o seu token de acesso da API CEPAberto.com

### Executar com Docker

Utilizando Docker (imagem `golang:alpine`) com o comando abaixo a imagem será compilada e executada na porta `8000`. 

```sh
$ make run.docker

___.
\_ |_________   ____  ____ ______
| __ \_  __ \_/ ___\/ __ \\____ \
| \_\ \  | \/\  \__\  ___/|  |_> >
|___  /__|    \___  >___  >   __/
    \/            \/    \/|__|
http://github.com/leogregianin/brcep

starting server on 8000
```

Para visualizar os dados acesse [http://localhost:8000/78048000/json](http://localhost:8000/78048000/json).

### Executando localmente

Considerando que você tem o Golang 1.13 instalado localmente, o comando abaixo irá efetuar o download das dependencias e compilar um binário para execução local.

```sh
$ make run.local
```

Você pode escolher outras arquiteturas para compilar o binário caso precise fazer o deploy do binário para outros sistemas:

```bash
$ make build.local
$ make build.linux.armv8
$ make build.linux.armv7
$ make build.linux
$ make build.osx
$ make build.windows
```

Os comandos acima geran um binário na pasta `bin`.

### Executando testes

```sh
$ make test
```

## Licença de uso

[MIT License](LICENSE)

Copyright (c) 2017 Leonardo Gregianin

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
