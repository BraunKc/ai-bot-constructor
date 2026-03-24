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

dotenvs:
	cp ./auth-gateway/.env.example ./auth-gateway/.env
	cp ./auth-service/.env.example ./auth-service/.env
	cp ./orchestrator-gateway/.env.example ./orchestrator-gateway/.env
	cp ./orchestrator-service/.env.example ./orchestrator-service/.env
	cp ./executor-service/.env.example ./executor-service/.env
	echo "u need to get open router api key (https://openrouter.ai/workspaces/default/keys) and paste it to OPEN_ROUTER_TOKEN at ./executor-service/.env"

docker-builds:
	docker build -t ai-auth-gateway ./auth-gateway
	docker build -t ai-auth-service ./auth-service
	docker build -t ai-orchestrator-gateway ./orchestrator-gateway
	docker build -t ai-orchestrator-service ./orchestrator-service
	docker build -t ai-executor-service ./executor-service