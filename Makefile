install:
	go install -v

build:
	go build -v ./...

lint:
	golint ./... && go vet ./...

test:
	go test -v ./...

cover:
	go test -v ./... --cover

deps: dev-deps
	go get -u github.com/nats-io/nats
	go get -u github.com/julienschmidt/httprouter
	go get -u github.com/lib/pq
	go get -u github.com/r3labs/binary-prefix
	go get -u github.com/ernestio/ernest-config-client

dev-deps:
	go get -u github.com/golang/lint/golint
	go get -u github.com/smartystreets/goconvey/convey

clean:
	go clean
