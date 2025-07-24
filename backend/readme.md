kubectl port-forward mongo-0 27017:27017 --address 0.0.0.0

docker run --name black-shop-mongo -p 27017:27017 -e MONGO_INITDB_ROOT_USERNAME=omides248 -e MONGO_INITDB_ROOT_PASSWORD=123123 -d mongo
docker run --name black-shop-postgres -p 5432:5432 -e POSTGRES_USER=omides248 -e POSTGRES_PASSWORD=123123 -e POSTGRES_DB=identity_db -d postgres



go get github.com/grpc-ecosystem/grpc-gateway/v2
go get github.com/jackc/pgx/v5


git clone https://github.com/googleapis/googleapis.git third_party/googleapis

protoc --proto_path=. --proto_path=third_party/googleapis --go_out=. --go-grpc_out=. --grpc-gateway_out=. api/proto/v1/catalog_service.proto api/proto/v1/identity_service.proto


Catalog Service:
go ./cmd/catalog_service/main.go

Identity Service:
go ./cmd/identity_service/main.go