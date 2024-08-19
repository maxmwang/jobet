gen: gen-proto gen-sqlc gen-mocks

.PHONY: gen-sqlc
gen-sqlc:
	sqlc generate

.PHONY: gen-proto
gen-proto:
	rm -r ./internal/proto || true # ignore doesn't exist error
	protoc --go_out=. --go_opt=paths=import \
		--go-grpc_out=. --go-grpc_opt=paths=import \
		api/prober.proto
	protoc --go_out=. --go_opt=paths=import \
		api/zeromq.proto

.PHONY: gen-mocks
gen-mocks:
	mockgen -source=internal/scrape/scraper.go -destination=internal/mocks/mock_scraper.go -package=mocks Scraper
	mockgen -source=internal/worker/publisher.go -destination=internal/mocks/mock_publisher.go -package=mocks Publisher

.PHONY: build-tsunami
build-tsunami:
	GOOS=linux GOARCH=amd64 go build -o build ./...

local-deploy: build-tsunami
	scp build/* maxmwang@tsunami.ocf.berkeley.edu:~/jobet/build
