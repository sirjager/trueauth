# grpc and gateways
SERVICE_NAME=trueauth
RPCS_DIR=../rpcs
PROTO_DIR=./proto

GO_RPC_DIR=$(RPCS_DIR)/$(SERVICE_NAME)/go
JS_RPC_DIR=$(RPCS_DIR)/$(SERVICE_NAME)/js
PY_RPC_DIR=$(RPCS_DIR)/$(SERVICE_NAME)/py

# open api swagger documentations
STATIK_OUT=./docs
SWAGGER_OUT=./docs/swagger

# TEST: database configs 
DB_MIGRATIONS=./db/migration
DB_URL=postgres://postgres:password@localhost:5432/testdb?sslmode=disable

proto:
	- rm -rf $(GO_RPC_DIR) $(JS_RPC_DIR) $(PY_RPC_DIR)
	- mkdir -p $(GO_RPC_DIR) $(JS_RPC_DIR) $(PY_RPC_DIR)
	- rm -f $(SWAGGER_OUT)/*.swagger.json
	- rm -rf $(STATIK_OUT)/statik

	protoc \
	--proto_path=$(PROTO_DIR) --go_out=$(GO_RPC_DIR) --go_opt=paths=source_relative	\
	--go-grpc_out=$(GO_RPC_DIR) --go-grpc_opt=paths=source_relative	\
	--grpc-gateway_out=$(GO_RPC_DIR) --grpc-gateway_opt=paths=source_relative \
	--openapiv2_out=$(SWAGGER_OUT) --openapiv2_opt=allow_merge=true,merge_file_name=$(SERVICE_NAME) \
	--grpc-web_out=import_style=commonjs+dts,mode=grpcwebtext:$(JS_RPC_DIR) \
	$(PROTO_DIR)/*.proto
	statik -src=$(SWAGGER_OUT) -dest=$(STATIK_OUT)
	- python -m grpc_tools.protoc -I$(PROTO_DIR) \
	--python_out=$(PY_RPC_DIR) --grpc_python_out=$(PY_RPC_DIR) \
	$(PROTO_DIR)/*.proto


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
	dbdocs build ./db/db.dbml

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

