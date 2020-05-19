build:
	@echo ">> building binaries - darwin_amd64"
	@GOOS=darwin GOARCH=amd64 go build -o build/amqpsniff-darwin-amd64 amqpsniff.go && gzip -f build/amqpsniff-darwin-amd64
	@echo ">> building binaries - linux_amd64"
	@GOOS=linux GOARCH=amd64 go build -o build/amqpsniff-linux-amd64 amqpsniff.go && gzip -f build/amqpsniff-linux-amd64
	@echo ">> building binaries - windows_amd64"
	@GOOS=windows GOARCH=amd64 go build -o build/amqpsniff-windows-amd64.exe amqpsniff.go \
		&& zip -9 build/amqpsniff-windows-amd64.zip build/amqpsniff-windows-amd64.exe && rm build/amqpsniff-windows-amd64.exe
	@echo ">> building binaries - alpine_amd64"
	@docker run --rm -v $(PWD):/app -w /app golang:1.14-alpine go build -o build/amqpsniff-alpine-amd64 \
		&& gzip -f build/amqpsniff-alpine-amd64


.PHONY: build
