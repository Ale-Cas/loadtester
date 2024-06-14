install: 
	go install
	export GOROOT=$HOME/go
	export PATH=$PATH:$GOROOT/bin

test:
	go test -v ./...

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

fmt:
	gofmt -w .