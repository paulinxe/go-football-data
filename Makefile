test:
	go test ./... 2>&1

check:
	export GOFLAGS="-buildvcs=false" && \
	golangci-lint run