.PHONY: gen-proto gen-sqlc gen-mocks

gen: gen-proto gen-sqlc gen-mocks

gen-sqlc:
	sqlc generate

gen-proto:
	rm -r ./internal/proto || true # ignore doesn't exist error
	protoc --go_out=. --go_opt=paths=import \
		--go-grpc_out=. --go-grpc_opt=paths=import \
		api/prober.proto
	protoc --go_out=. --go_opt=paths=import \
		api/zeromq.proto

gen-mocks:
	mockgen -source=internal/scrape/scraper.go -destination=internal/mocks/mock_scraper.go -package=mocks Scraper
	mockgen -source=internal/worker/publisher.go -destination=internal/mocks/mock_publisher.go -package=mocks Publisher
