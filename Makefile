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
	go get github.com/nats-io/nats
	go get github.com/julienschmidt/httprouter
	go get github.com/lib/pq
	go get github.com/r3labs/binary-prefix
	go get github.com/ernestio/ernest-config-client

dev-deps:
	go get github.com/golang/lint/golint
	go get github.com/smartystreets/goconvey/convey

clean:
	go clean
