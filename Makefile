test:
	go test -v ./...
build:
	go build -o bin/seeder cmd/seeder/main.go
	go build -o bin/client cmd/client/main.go
run-server:
	go run cmd/seeder/main.go
run-client:
	go run cmd/client/main.go --port 7687
