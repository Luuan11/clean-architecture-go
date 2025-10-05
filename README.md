# Clean Architecture - Order Management System

Sistema de gerenciamento de pedidos (orders) implementado com Clean Architecture em Go, oferecendo três interfaces de comunicação: REST API, gRPC e GraphQL.

## 🏗️ Arquitetura

Este projeto segue os princípios da Clean Architecture:

```
├── cmd/
│   └── server/           # Aplicação principal
├── internal/
│   ├── entity/           # Entidades de domínio
│   ├── usecase/          # Casos de uso (regras de negócio)
│   └── infra/            # Infraestrutura
│       ├── database/     # Repositórios
│       ├── grpc/         # Servidor gRPC
│       ├── graphql/      # Servidor GraphQL
│       └── web/          # Servidor REST
└── sql/
    └── migrations/       # Migrações do banco de dados
```

## 🚀 Tecnologias Utilizadas

- **Go 1.21+**
- **MySQL 8.0**
- **gRPC** - Comunicação entre serviços
- **GraphQL** - API query language
- **REST API** - HTTP/JSON
- **Docker & Docker Compose** - Containerização
- **Protocol Buffers** - Serialização de dados

## 📋 Pré-requisitos

- [Go 1.21+](https://golang.org/dl/)
- [Docker](https://www.docker.com/get-started)
- [Docker Compose](https://docs.docker.com/compose/install/)
- [Protocol Buffers Compiler (protoc)](https://grpc.io/docs/protoc-installation/)

## 🔧 Instalação e Configuração

### 1. Clone o repositório

```bash
git clone <repository-url>
cd clean-architecture
```

### 2. Configurar variáveis de ambiente

O arquivo `.env` já está configurado com valores padrão:

```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=root
DB_NAME=orders

REST_PORT=8080
GRPC_PORT=50051
GRAPHQL_PORT=8081
```

### 3. Gerar código Protocol Buffers

**No Linux/Mac:**
```bash
chmod +x scripts/generate-proto.sh
./scripts/generate-proto.sh
```

**No Windows:**
```bash
scripts\generate-proto.bat
```

**Ou usando Make:**
```bash
make proto
```

### 4. Instalar dependências

```bash
go mod download
go mod tidy
```

**Ou usando Make:**
```bash
make deps
```

## 🐳 Executando com Docker

A forma mais fácil de executar o projeto é usando Docker Compose:

```bash
docker-compose up -d
```

Isso irá:
1. Criar o banco de dados MySQL
2. Executar as migrações
3. Iniciar a aplicação

Para visualizar os logs:
```bash
docker-compose logs -f
```

Para parar os containers:
```bash
docker-compose down
```

**Ou usando Make:**
```bash
make docker-up      # Subir containers
make docker-logs    # Ver logs
make docker-down    # Parar containers
```

## 💻 Executando Localmente (sem Docker)

### 1. Inicie o MySQL

Certifique-se de ter o MySQL rodando localmente ou use apenas o container do MySQL:

```bash
docker-compose up -d mysql
```

### 2. Execute as migrações

```bash
mysql -h localhost -u root -proot orders < sql/migrations/001_create_orders_table.sql
```

### 3. Execute a aplicação

```bash
go run cmd/server/main.go
```

**Ou usando Make:**
```bash
make run
```

## 📡 Portas dos Serviços

| Serviço  | Porta | Endpoint                        |
|----------|-------|---------------------------------|
| REST API | 8080  | http://localhost:8080/order     |
| gRPC     | 50051 | localhost:50051                 |
| GraphQL  | 8081  | http://localhost:8081/graphql   |

## 📚 Documentação das APIs

### REST API

#### Criar Order (POST)
```bash
POST http://localhost:8080/order
Content-Type: application/json

{
  "price": 100.50,
  "tax": 10.05
}
```

#### Listar Orders (GET)
```bash
GET http://localhost:8080/order
```

### GraphQL

Acesse o GraphiQL em: http://localhost:8081/graphql

#### Criar Order (Mutation)
```graphql
mutation {
  createOrder(price: 150.75, tax: 15.08) {
    id
    price
    tax
    final_price
  }
}
```

#### Listar Orders (Query)
```graphql
{
  listOrders {
    id
    price
    tax
    final_price
  }
}
```

### gRPC

Para testar o gRPC, você pode usar [grpcurl](https://github.com/fullstorydev/grpcurl):

#### Listar serviços disponíveis
```bash
grpcurl -plaintext localhost:50051 list
```

#### Criar Order
```bash
grpcurl -plaintext -d '{"price": 200.00, "tax": 20.00}' localhost:50051 pb.OrderService/CreateOrder
```

#### Listar Orders
```bash
grpcurl -plaintext -d '{}' localhost:50051 pb.OrderService/ListOrders
```

## 🧪 Testando com api.http

O projeto inclui um arquivo `api.http` com exemplos de requisições. Se você usar o VS Code com a extensão REST Client, pode executar as requisições diretamente do editor.

## 🗂️ Estrutura do Banco de Dados

### Tabela: orders

| Campo       | Tipo           | Descrição              |
|-------------|----------------|------------------------|
| id          | VARCHAR(36)    | UUID (Primary Key)     |
| price       | DECIMAL(10,2)  | Preço base             |
| tax         | DECIMAL(10,2)  | Taxa/Imposto           |
| final_price | DECIMAL(10,2)  | Preço final calculado  |
| created_at  | TIMESTAMP      | Data de criação        |
| updated_at  | TIMESTAMP      | Data de atualização    |

## 🛠️ Comandos Úteis (Make)

```bash
make help         # Mostra todos os comandos disponíveis
make proto        # Gera código protobuf
make deps         # Instala dependências
make build        # Compila a aplicação
make run          # Executa a aplicação
make docker-up    # Sobe containers Docker
make docker-down  # Para containers Docker
make docker-logs  # Mostra logs
make test         # Executa testes
make clean        # Remove arquivos gerados
```

## 🔍 Verificando se está funcionando

Após iniciar a aplicação, você deve ver mensagens como:

```
Successfully connected to database
REST server listening on port 8080
gRPC server listening on port 50051
GraphQL server listening on port 8081
GraphiQL available at http://localhost:8081/graphql
```

Teste rapidamente com:

```bash
# REST API
curl -X POST http://localhost:8080/order \
  -H "Content-Type: application/json" \
  -d '{"price": 100, "tax": 10}'

curl http://localhost:8080/order
```

## 📝 Notas

- O banco de dados é criado automaticamente quando você executa `docker-compose up`
- As migrações são executadas automaticamente na inicialização do container MySQL
- O GraphiQL (interface web interativa) está disponível em http://localhost:8081/graphql
- Para ambientes de produção, lembre-se de configurar senhas fortes e variáveis de ambiente apropriadas

## 🤝 Contribuindo

Sinta-se à vontade para contribuir com o projeto!

## 📄 Licença

Este projeto é livre para uso educacional e comercial.
