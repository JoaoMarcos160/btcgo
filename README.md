# BTC GO v0.6.0
[![instalação do Go no Windows](https://img.youtube.com/vi/679Zc7ZQLtI/0.jpg)](https://www.youtube.com/watch?v=679Zc7ZQLtI)

## Requisitos
  -  [Go][install-go]
  -  [Git][install-git]
  -  Terminal

# Instruções para rodar o projeto no Windows.

 * Clona o repo e brota na pasta:
```bash
git clone [git@github.com:lmajowka/btcgo.git](https://github.com/lmajowka/btcgo.git) && cd btcgo
```
 
 * Instala as parada:
```bash
go mod tidy
```

 * Faz o build do projeto:
```bash
go build -o btcgo.exe ./cmd/main.go
```

 * Executa o que foi compilado:
```bash
btcgo
```


# Instruções para rodar o projeto no Linux / MacOS.

 * Clona o repo e brota na pasta:
```bash
  git clone [git@github.com:lmajowka/btcgo.git](https://github.com/lmajowka/btcgo.git) && cd btcgo
```
 
 * Instala as parada:
```bash
 go mod tidy
```

 * Faz o build do projeto:
```bash
go build -o btcgo ./cmd/main.go
```

 * Executa o que foi compilado:
```bash
./btcgo
```

# Instruções para rodar o projeto em container.

## Requisitos
  -  [Docker][install-docker]
  -  [Docker-compose][install-docker-compose]
  -  [Git][install-git]

## Execução da parada

#### Use a imagem desse repositório já construída
 * Rode com o seguinte comando alterando as variáveis de ambiente como quiser:
```bash
docker run -e CPU_COUNT=0 -e RANGE_NUMBER=68 -e MODS=2 -e USE_DB=1 -e NUM_RECS=10000 -e START_MODE=2 -e START_PERCENT=75 -d ghcr.io/joaomarcos160/btcgo:main
```

#### Contrua a imagem Docker locamente
 * Clona o repo:
```bash
git clone git@github.com:lmajowka/btcgo.git && cd btcgo
```
 * Build do Dockerfile:
```bash
docker buildx build --no-cache -t btcgo .
```
 * Executa a imagem contruída no passo anterior:
```bash
docker run --rm -it --name btcgo btcgo
```
Este container será deletado se a aplicação parar, para executar novamente basta executar o último comando acima.

[install-go]: https://go.dev/doc/install
[install-git]: https://git-scm.com/download/win
[install-docker]: https://www.docker.com/get-started/
[install-docker-compose]: https://docs.docker.com/compose/install/
