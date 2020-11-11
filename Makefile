build:
	go build -o . ./cmd/protoc-gen-gojrpcmodel
run:
	protoc --plugin protoc-gen-gojrpcmodel --gojrpcmodel_out=. rpc/*.proto
lints:
	golangci-lint run --fix -c ./.linters.yml