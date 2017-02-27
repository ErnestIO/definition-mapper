install:
	go install -v

build:
	go build -v ./...

lint:
	gometalinter --config .linter.conf

test:
	go test -v ./...

cover:
	go test -v ./... --cover

deps:
	go get github.com/nats-io/nats
	go get github.com/ernestio/ernest-config-client
	go get github.com/ernestio/libmapper

dev-deps: deps
	go get github.com/smartystreets/goconvey/convey
	go get github.com/alecthomas/gometalinter
	gometalinter --install

clean:
	go clean
