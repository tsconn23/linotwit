.PHONY: build clean test run

GO=CGO_ENABLED=0 go

APPLICATION=cmd/agent

.PHONY: $(APPLICATION)

build: $(APPLICATION)

cmd/agent:
	$(GO) build -o $@ ./cmd/agent

clean:
	rm -f $(APPLICATION)

test:
	go test -coverprofile=coverage.out ./...
	go vet ./...
	gofmt -l .
	[ "`gofmt -l .`" = "" ]

run:
	cd bin && ./launch.sh
