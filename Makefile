.PHONY: gen-proto gen-sqlc gen-mocks

gen: gen-proto gen-sqlc gen-mocks

gen-sqlc:
	sqlc generate

gen-proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		internal/proto/jobet.proto
	protoc --go_out=. --go_opt=paths=source_relative \
		api/zeromq.proto

gen-mocks:
	mockgen -source=internal/scraper/scraper.go -destination=internal/scraper/mock_scraper.go -package=scraper Scraper
