kubectl port-forward mongo-0 27017:27017 --address 0.0.0.0

docker run --name black-shop-mongo -p 27017:27017 -e MONGO_INITDB_ROOT_USERNAME=omides248 -e MONGO_INITDB_ROOT_PASSWORD=123123 -d mongo
docker run --name black-shop-postgres -p 5432:5432 -e POSTGRES_USER=omides248 -e POSTGRES_PASSWORD=123123 -e POSTGRES_DB=identity_db -d postgres



go get github.com/grpc-ecosystem/grpc-gateway/v2
go get github.com/jackc/pgx/v5


git clone https://github.com/googleapis/googleapis.git third_party/googleapis

protoc --proto_path=. --proto_path=third_party/googleapis --go_out=. --go-grpc_out=. --grpc-gateway_out=. api/proto/v1/catalog_service.proto api/proto/v1/identity_service.proto api/proto/v1/order_service.proto

protoc --proto_path=. --proto_path=third_party/googleapis --go_out=. --go-grpc_out=. --grpc-gateway_out=. $(find api/proto/v1 -name "*.proto")

protoc --proto_path=. --proto_path=third_party/googleapis --go_out=. --go-grpc_out=. --grpc-gateway_out=. api/proto/v1/*.proto




protoc --proto_path=api/proto --proto_path=../third_party/googleapis --go_out=api/proto --go-grpc_out=api/proto --grpc-gateway_out=api/proto api/proto/catalog/v1/*.proto

protoc \
--proto_path=api/proto \
--proto_path=../third_party/googleapis \
--go_out=api/proto \
--go-grpc_out=api/proto \
--grpc-gateway_out=api/proto \
$(find api/proto -name "*.proto")



Catalog Service:
go ./cmd/catalog_service/main.go

Identity Service:
go ./cmd/identity_service/main.go

Order Service:
go ./cmd/identity_service/main.go



go run github.com/99designs/gqlgen generate


make up
make down





# Structure
```text
D:.
│   .gitignore
│   go.work
│   Makefile
│   readme.md
│   tools.go
│
├───api_gateway
│   │   go.mod
│   │   go.sum
│   │   main.go
│   │
│   └───internal
│       └───delivery
│
├───catalog
│   │   go.mod
│   │   go.sum
│   │   main.go
│   │
│   ├───api
│   │   └───v1
│   │           catalog_service.proto
│   │           ...
│   │
│   └───internal
│       ├───adapters
│       ├───application
│       ├───delivery
│       └───domain
│
├───identity
│   │   go.mod
│   │   go.sum
│   │   main.go
│   │
│   ├───api
│   │   └───v1
│   │           identity_service.proto
│   │
│   └───internal
│       ├───adapters
│       ├───application
│       ├───delivery
│       └───domain
│
├───order
│   │   go.mod
│   │   go.sum
│   │   main.go
│   │
│   ├───api
│   │   └───v1
│   │           order_service.proto
│   │
│   └───internal
│       ├───adapters
│       ├───application
│       ├───delivery
│       └───domain
│
├───pkg
│   ├───auth
│   │   │   go.mod
│   │   │   jwt.go
│   │
│   ├───config
│   │   │   go.mod
│   │   │   config.go
│   │
│   └───logger
│       │   go.mod
│       │   logger.go
│
├───build
├───configs
├───deployments
├───scripts
└───test
```
