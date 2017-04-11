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
	go get github.com/mitchellh/mapstructure
	go get github.com/ghodss/yaml
	go get gopkg.in/r3labs/graph.v2
	go get github.com/r3labs/binary-prefix

dev-deps: deps
	go get golang.org/x/crypto/pbkdf2
	go get github.com/ernestio/crypto
	go get github.com/ernestio/crypto/aes
	go get github.com/tidwall/gjson
	go get github.com/smartystreets/goconvey/convey
	go get github.com/alecthomas/gometalinter
	go get github.com/stretchr/testify/suite
	gometalinter --install

clean:
	go clean
