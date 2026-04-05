.PHONY: run build tidy test clean

run:
	go run ./cmd/server

build:
	go build -o bin/finance-backend ./cmd/server

tidy:
	go mod tidy

test:
	go test ./... -v

clean:
	rm -f bin/finance-backend finance.db
