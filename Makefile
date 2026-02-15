PROTO_SERVICE=?
PROTO_CLIENT=?
PROTO_VERSION=v1

GEN_DIR=${PROTO_SERVICE}/api
CLIENT_DIR=${PROTO_CLIENT}/api

proto:
	mkdir -p ${GEN_DIR}

	protoc --proto_path=contracts \
	--go_out=${GEN_DIR} --go_opt=paths=source_relative \
	--go-grpc_out=${GEN_DIR} --go-grpc_opt=paths=source_relative \
	contracts/${PROTO_SERVICE}/${PROTO_VERSION}/${PROTO_SERVICE}.proto

protoclient:
	mkdir -p ${CLIENT_DIR}

	protoc --proto_path=contracts \
	--go_out=${CLIENT_DIR} --go_opt=paths=source_relative \
	--go-grpc_out=${CLIENT_DIR} --go-grpc_opt=paths=source_relative \
	contracts/${PROTO_SERVICE}/${PROTO_VERSION}/${PROTO_SERVICE}.proto