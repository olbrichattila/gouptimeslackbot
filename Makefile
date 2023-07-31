build:
	go build -o ./upclient .
run:
	go run .
run-background: build
	./upclient > /dev/null 2>&1 &
run-test:
	go test
	