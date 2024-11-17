run:
	go run cmd/weather/main.go

test:
	go test ./... -v -cover

generate:
	go generate ./...

pprof-test:
	go test -coverprofile=cover.out ./...
	go tool cover -html=cover.out -o cover.html