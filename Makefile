test:
	$(MAKE) clean local-mongo
	go test -v ./...

build:
	go build -o bin/seeder cmd/seeder/main.go
	go build -o bin/client cmd/client/main.go

run-seeder:
	$(MAKE) clean local-mongo
	go run cmd/seeder/main.go --mongo.uri mongodb://localhost:27017 --check.interval 10s

run-client:
	go run cmd/client/main.go --port $(PORT) --server.address http://127.0.0.1:5000

local-mongo:
	docker run -d --network host --name seeder-mongo mongo:latest

local-seeder:
	$(MAKE) clean local-mongo
	docker build -t seeder-image .
	docker run --network host --name seeder-server seeder-image:latest seeder --mongo.uri mongodb://localhost:27017 --check.interval 10s

local-client:
	docker build -t seeder-image .
	docker run --network host --name seeder-client seeder-image:latest client --port $(PORT) --server.address http://127.0.0.1:5000

clean:
	rm -rf bin
	docker rm -f seeder-server seeder-client seeder-mongo
