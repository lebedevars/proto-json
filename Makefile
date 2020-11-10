run:
	go build . && protoc --plugin protoc-gen-protojson --protojson_out=rpc *.proto
lints:
	golangci-lint run --fix -c ./.linters.yml