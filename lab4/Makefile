PROTOC = protoc
PROTOC_GEN_GO = protoc-gen-go
PROTOC_GEN_GRPC = protoc-gen-go-grpc

PROTO_FILES = ./api/internal_api.proto ./api/external_api.proto
OUTPUT_DIR = ./src/generated

all: generate

generate: $(PROTO_FILES)
	$(PROTOC) --go_out=$(OUTPUT_DIR) --go-grpc_out=$(OUTPUT_DIR) $(PROTO_FILES)
	@echo 'All proto files generated successfully!'

clean:
	rm -rf $(OUTPUT_DIR)/*

.PHONY: all generate clean
