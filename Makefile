.PHONY: gen-proto gen-sqlc

gen: gen-proto gen-sqlc

gen-sqlc:
	sqlc generate

gen-proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		internal/proto/jobet.proto
