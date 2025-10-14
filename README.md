# Clean Architecture - Order Management

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)
[![MySQL](https://img.shields.io/badge/MySQL-8.0-orange.svg)](https://www.mysql.com)
[![Docker](https://img.shields.io/badge/Docker-Supported-blue.svg)](https://docker.com)
[![gRPC](https://img.shields.io/badge/gRPC-Supported-green.svg)](https://grpc.io)
[![GraphQL](https://img.shields.io/badge/GraphQL-Supported-pink.svg)](https://graphql.org)

### Sistema de Gerenciamento de Pedidos com Clean Architecture

<img width="1920" height="1080" alt="clean arc" src="https://github.com/user-attachments/assets/b2cf052e-9e03-44f9-9c6e-1ba5cfcf9a86" />

### 💬 Sobre: 
API completa de gerenciamento de orders implementada em Go seguindo os princípios da Clean Architecture, oferecendo três interfaces de comunicação distintas:
```plaintext
- CRUD completo de Orders
- REST API (JSON/HTTP)
- gRPC (Protocol Buffers)
- GraphQL (Query Language)
- Persistência MySQL
- Docker Compose para infraestrutura
- Testes unitários com 81% de cobertura
```

### ✨ Executando:

#### Via Docker Compose (Recomendado):
```shell
❯ docker-compose up --build
```

#### Via binário local:
```shell
❯ go run cmd/server/main.go
```

#### Parâmetros de configuração:
```plaintext
DB_HOST         Endereço do MySQL (padrão: localhost)
DB_PORT         Porta do MySQL (padrão: 3306)
DB_USER         Usuário do banco (padrão: root)
DB_PASSWORD     Senha do banco (padrão: root)
DB_NAME         Nome do banco (padrão: orders)
REST_PORT       Porta REST API (padrão: 8080)
GRPC_PORT       Porta gRPC (padrão: 50051)
GRAPHQL_PORT    Porta GraphQL (padrão: 8081)
```

#### Exemplo de uso:
```shell
❯ curl -X POST http://localhost:8080/order -d '{"price":100,"tax":10}'
❯ curl http://localhost:8080/order
❯ curl http://localhost:8080/order/{id}
```

## 🏗️ Arquitetura

```
clean-architecture/
├── cmd/
│   └── server/
│       └── main.go              # Entry point da aplicação
├── internal/
│   ├── entity/
│   │   ├── order.go             # Entidade de domínio
│   │   └── interface.go         # Contratos do repositório
│   ├── usecase/
│   │   ├── create_order.go      # Caso de uso: Criar
│   │   ├── list_orders.go       # Caso de uso: Listar
│   │   ├── get_order_by_id.go   # Caso de uso: Buscar por ID
│   │   ├── update_order.go      # Caso de uso: Atualizar
│   │   ├── delete_order.go      # Caso de uso: Deletar
│   │   └── test/                # Testes unitários (81% coverage)
│   └── infra/
│       ├── database/
│       │   └── order_repository.go    # Implementação MySQL
│       ├── grpc/
│       │   ├── proto/                 # Definições protobuf
│       │   └── service/               # Serviços gRPC
│       ├── graphql/
│       │   └── graph/                 # Resolvers GraphQL
│       └── web/
│           └── webserver/             # Handlers REST
├── sql/
│   └── migrations/
│       └── 001_create_orders_table.sql
├── docker-compose.yml           # Orquestração Docker
├── Dockerfile                   # Multi-stage build
├── api.http                     # Exemplos de requisições
└── README.md                    # Esta documentação
```

## 🐳 Docker

### Build da Imagem
```bash
docker build -t clean-architecture .
```

## ⭐ APIs Disponíveis

### REST API (Port 8080)
```bash
POST   /order          # Criar order
GET    /order          # Listar orders
GET    /order/{id}     # Buscar por ID
PUT    /order/{id}     # Atualizar order
DELETE /order/{id}     # Deletar order
```

### GraphQL (Port 8081)
```graphql
# Queries
query {
  listOrders { id price tax final_price }
  getOrder(id: "uuid") { id price tax final_price }
}

# Mutations
mutation {
  createOrder(price: 100.0, tax: 10.0) { id price tax final_price }
  updateOrder(id: "uuid", price: 200.0, tax: 20.0) { id price tax final_price }
  deleteOrder(id: "uuid")
}
```

**GraphiQL**: http://localhost:8081/graphql

### gRPC (Port 50051)
```bash
# Listar serviços
grpcurl -plaintext localhost:50051 list

# Criar order
grpcurl -plaintext -d '{"price": 100, "tax": 10}' \
  localhost:50051 pb.OrderService/CreateOrder

# Listar orders
grpcurl -plaintext -d '{}' \
  localhost:50051 pb.OrderService/ListOrders
```

## 🚀 Instalação

    - Clonar repositório 
    $ git clone https://github.com/Luuan11/clean-architecture-go.git

    - Instalar dependências
    $ go mod tidy

    - Compilar projeto
    $ go build -o server ./cmd/server

    - Executar
    $ ./server

---

Made with 💜 by [Luan Fernando](https://www.linkedin.com/in/luan-fernando/).
