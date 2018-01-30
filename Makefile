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
	go get -u github.com/golang/dep/cmd/dep
	dep ensure -v

dev-deps: deps
	go get github.com/alecthomas/gometalinter
	gometalinter --install

clean:
	go clean
