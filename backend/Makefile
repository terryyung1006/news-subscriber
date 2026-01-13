.PHONY: build run gen-proto install-proto update-proto generate-invites db-up db-migrate db-down db-clean db-reset db-seed

# Build commands
build:
	go build -o server ./src/main.go

run: build
	./server

# Proto commands
gen-proto:
	mkdir -p api/proto/v1
	protoc --go_out=. --go_opt=module=news-subscriber-core \
	--go-grpc_out=. --go-grpc_opt=module=news-subscriber-core \
	-I proto-repo/proto \
	proto-repo/proto/*.proto

install-proto:
	git submodule update --init --recursive

update-proto: install-proto gen-proto

# Utility commands
generate-invites:
	go run tools/generate-invites/main.go

compose-up:
	docker-compose up -d

db-migrate:
	@go run tools/migrate/main.go -action=migrate

db-reset:
	@go run tools/migrate/main.go -action=reset

db-seed:
	@go run tools/migrate/main.go -action=seed

db-clean:
	docker-compose down -v
	@echo "✓ Removed containers and volumes"

# Convenience: reset database completely (down, up, reset, seed)
db-fresh: db-clean compose-up
	@go run tools/migrate/main.go -action=fresh

setup-dev: compose-up update-proto db-migrate run
