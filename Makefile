# grpc and gateways
SERVICE_NAME=trueauth
RPCS_DIR=../rpcs
PROTO_DIR=./internal/proto

GO_RPC_DIR=$(RPCS_DIR)/$(SERVICE_NAME)/go

# open api swagger documentations
STATIK_OUT=./docs
SWAGGER_OUT=./docs/swagger

# TEST: database configs 
DB_MIGRATIONS=./migration
DB_URL=postgres://postgres:testpassword@localhost:5432/postgres?sslmode=disable

# Generate protobuf and gRPC code, Swagger documentation, and static assets
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

# Clean up go.sum and tidy dependencies
tidy:
	rm -f ./go.sum
	go get github.com/sirjager/rpcs@latest
	go mod tidy

# Run tests
test:
	go clean -testcache
	go test -v -cover -short ./... 

# Build the project
build:
	golint ./...
	go build -o ./dist/main ./cmd/main.go

# Linting and formatting
lint:
	golint ./...

# Run the project
run:
	go run ./cmd/main.go

# Generate database documentation
dbdocs:
	dbdocs build ./internal/db/db.dbml;

# Generate database schema
dbschema:
	dbml2sql --postgres -o ./internal/db/schema.sql ./internal/db/db.dbml

# Create a new database migration
migratenew:
	migrate create -ext sql -dir $(DB_MIGRATIONS) -seq $(filter-out $@,$(MAKECMDGOALS))

# Run database migrations (up)
migrateup:
	migrate -source file://$(DB_MIGRATIONS) -database $(DB_URL) -verbose up

# Rollback database migrations (down)
migratedown:
	migrate -source file://$(DB_MIGRATIONS) -database $(DB_URL) -verbose down --all

# Generate SQL code using SQLC
sqlc:
	sqlc generate

release:
	goreleaser --snapshot --clean

image:
	docker build --pull --rm -f "Dockerfile" -t gomicro:latest "."

# Declare phony targets to prevent conflicts with file names
.PHONY: proto tidy test run dbdocs dbschema migratenew migrateup migratedown sqlc release
