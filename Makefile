.PHONY: generate

generate:
	mkdir -p generated
	protoc --proto_path=api api/*.proto --go_out=generated --go-grpc_out=generated
