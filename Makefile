
test:
	go test -v -cover ./...

server:
	go run main.go

ampqServer:
	go run consumer/main.go

.PHONY:  test server consServer