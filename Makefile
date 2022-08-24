build: tools-install generate
	go mod vendor

engine:
	CGO_ENABLED=0 GOOS=darwin go build -a -installsuffix cgo -o ${BINARY} /tests2 cmd/*.go

test:
	go test -v -cover -covermode=atomic ./...

docker:
	docker build -t order-test-service .

run:
	go run cmd/main.go