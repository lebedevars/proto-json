run:
	protoc --plugin protoc-gen-protojson --protojson_out=rpc example.proto
lints:
	golangci-lint run -c ./.linters.yml