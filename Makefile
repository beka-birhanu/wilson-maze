build:
	@go build -o ./bin/willson-maze ./main.go

test:
	go test -v ./...

run: build
	@./bin/willson-maze
