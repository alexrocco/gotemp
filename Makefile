lint:
	golangci-lint run --enable-all ./...
build-arm:
	env CC=arm-linux-gnueabihf-gcc GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=1 go build -v -o gotemp cmd/gotemp/main.go 