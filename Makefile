# grpc and gateways
SERVICE_NAME=trueauth
RPCS_DIR=../rpcs
PROTO_DIR=./proto

GO_RPC_DIR=$(RPCS_DIR)/$(SERVICE_NAME)/go

# open api swagger documentations
STATIK_OUT=./docs
SWAGGER_OUT=./docs/swagger

# TEST: database configs 
DB_MIGRATIONS=./db/migration
DB_URL=postgres://postgres:FW3F6ojfGN6IbZfJ@localhost:5432/testdb?sslmode=disable

proto:
	- rm -rf $(GO_RPC_DIR)
	- mkdir -p $(GO_RPC_DIR) 
	- rm -f $(SWAGGER_OUT)/*.swagger.json
	- rm -rf $(STATIK_OUT)/statik

	protoc \
	--proto_path=$(PROTO_DIR) --go_out=$(GO_RPC_DIR) --go_opt=paths=source_relative	\
	--go-grpc_out=$(GO_RPC_DIR) --go-grpc_opt=paths=source_relative	\
	--grpc-gateway_out=$(GO_RPC_DIR) --grpc-gateway_opt=paths=source_relative \
	--openapiv2_out=$(SWAGGER_OUT) --openapiv2_opt=allow_merge=true,merge_file_name=$(SERVICE_NAME) \
	$(PROTO_DIR)/*.proto
	statik -src=$(SWAGGER_OUT) -dest=$(STATIK_OUT)


tidy:
	rm -f ./go.sum
	rm -rf ./vendor
	go get github.com/sirjager/rpcs@latest
	go mod tidy
	# go mod vendor

test:
	go clean -testcache
	go test -v -cover -short ./... 

build:
	golint ./...
	go build -o ./dist/main ./cmd/main.go

run:
	go run ./cmd/main.go

dbdocs:
	dbdocs build ./db/db.dbml;

dbschema:
	dbml2sql --postgres -o ./db/schema.sql ./db/db.dbml


migratenew:
	migrate create -ext sql -dir $(DB_MIGRATIONS) -seq $(filter-out $@,$(MAKECMDGOALS))

migrateup:
	migrate -source file://$(DB_MIGRATIONS) -database $(DB_URL) -verbose up

migratedown:
	migrate -source file://$(DB_MIGRATIONS) -database $(DB_URL) -verbose down --all

sqlc:
	sqlc generate


gen:
	- make dbschema
	- make sqlc

.PHONY: proto tidy test run dbdocs dbschema migratenew migrateup migratedown sqlc gen

