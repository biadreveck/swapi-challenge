# swapi-challenge

Este projeto é a implementação de uma API REST para cadastro de planetas do mundo de Star Wars.

### Funcionalidades
- Adicionar um planeta (com nome, clima e terreno)
- Listar planetas
- Buscar planeta por nome
- Buscar planeta por ID
- Remover planeta

### Ferramentas
- Linguagem: [Go](https://golang.org/ "Go")
- Banco de dados: [MongoDB](https://www.mongodb.com/ "MongoDB")
- Criação de mocks: [mockery](https://github.com/vektra/mockery "mockery")

### Dependências
- Gin (go get github.com/gin-gonic/gin)
- MongoDB Go Driver (go get go.mongodb.org/mongo-driver/mongo)
- Viper (go get github.com/spf13/viper)
- Testify (go get github.com/stretchr/testify)

### Arquivo de configuração: *config/config.yml*
Exemplo:
```yaml
database:
  host: localhost:5000  
  dbname: dbname
  connectionTimeout: 10s
  commandTimeout: 30s
  user: dbUser
  password: dbPass

swapi:
  baseUrl: https://swapi.dev/api

server:
  address: :8080
```
- **database**: configurações do banco de dados MongoDB
	- **host**: endereço e porta para acesso ao banco de dados
	- **dbname**: nome do banco de dados
	- **connectionTimeout**: limite de tempo de conexão com o banco de dados
	- **commandTimeout**: limite de tempo dos requisições ao banco de dados
	- **user**: usuário para acesso ao banco de dados [opcional]
	- **password**: senha para acesso ao banco de dados [opcional]
- **swapi**: configurações da SWAPI (API de Star Wars)
	- **baseUrl**: endereço base para acesso à API
- **server**: configurações do servidor da API
	- **address**: endereço e porta de acesso à API

#### Adicionar um planeta (com nome, clima e terreno)

> Método: POST
Endpoint: /v1/planets

- **Campos do corpo**:
	- **name**: endereço e porta para acesso ao banco de dados [obrigatório]
	- **climate**: nome do banco de dados
	- **terrain**: limite de tempo de conexão com o banco de dados

##### Exemplo requisição:
> POST /v1/planets
```json
{
	"name": "Kamino",
	"climate": "temperate",
	"terrain": "ocean"
}
```

##### Exemplo resposta:
```json
{
    "data": {
        "id": "5f300f1713bd94e33937a4d0",
        "name": "Kamino",
        "climate": "temperate",
        "terrain": "ocean",
        "apparitions": 1
    }
}
```

#### Listar planetas

> Método: GET
Endpoint: /v1/planets

##### Exemplo requisição:
> GET /v1/planets

##### Exemplo resposta:
```json 
{
    "data": [
        {
            "id": "5f300c776a9701275f757eca",
            "name": "Tatooine",
            "climate": "arid",
            "terrain": "desert",
            "apparitions": 5
        },
        {
            "id": "5f300ef113bd94e33937a4cf",
            "name": "Alderaan",
            "climate": "temperate",
            "terrain": "grasslands, mountains",
            "apparitions": 2
        },
        {
            "id": "5f300f1713bd94e33937a4d0",
            "name": "Kamino",
            "climate": "temperate",
            "terrain": "ocean",
            "apparitions": 1
        },
        {
            "id": "5f3026ce8ab5e691db4aa03c",
            "name": "Yavin IV",
            "climate": "temperate, tropical",
            "terrain": "jungle, rainforests",
            "apparitions": 1
        }
    ]
}
```

#### Buscar planeta por nome

> Método: GET
Endpoint: /v1/planets?name={nome do planeta}

##### Exemplo requisição:
> GET /v1/planets?name=Yavin+IV

##### Exemplo resposta:
```json 
{
    "data": {
        "id": "5f3026ce8ab5e691db4aa03c",
        "name": "Yavin IV",
        "climate": "temperate, tropical",
        "terrain": "jungle, rainforests",
        "apparitions": 1
    }
}
```

#### Buscar planeta por ID

> Método: GET
Endpoint: /v1/planets/{id do planeta}

##### Exemplo requisição:
> GET /v1/planets/5f300ef113bd94e33937a4cf

##### Exemplo resposta:
- **200 OK**
```json 
{
    "data": {
        "id": "5f300ef113bd94e33937a4cf",
        "name": "Alderaan",
        "climate": "temperate",
        "terrain": "grasslands, mountains",
        "apparitions": 2
    }
}
```

#### Remover planeta

> Método: DELETE
Endpoint: /v1/planets/{id do planeta}

##### Exemplo requisição:
> DELETE /v1/planets/5f300ef113bd94e33937a4cf

##### Exemplo resposta:
- **204 No Content**

------------

#### Usando localmente:
Para rodar a aplicação localmente é necessário executar os seguintes passos:
1. Instalar as ferramentas abaixo na máquina local:
	- Go v1.14.6+
	- MongoDB v4.4.0+
2. Clonar esse repositório em qualquer diretório
3. Alterar o arquivo *config/config.yml* com as configurações desejadas
4. No diretório clonado, rodar a aplicação usando: **go run main.go**
5. Se desejar, executar os testes com o comando: **go test ./...**