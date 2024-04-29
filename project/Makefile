# This file help with some shortcuts to start, build, take it down our application. 
# For mac, you need to install make -> brew install make

GRPC_SERVER_BINARY=grpcServer
POD_BINA_BINARY=podBINA
POD_BITP_BINARY=podBITP
POD_MBTC_BINARY=podMBTC
POD_FOXB_BINARY=podFOXB
POD_RIPI_BINARY=podRIPI

## PODS
###########

## start: starts all pods
start-all: start-BINA start-BITP start-FOXB start-MBTC start-RIPI
	@echo "Starting pods"
	@echo "Done!"

## stop: stop the front end
stop-all:
	@echo "Stopping all pods..."
	@-pkill -SIGTERM -f "./${POD_BINA_BINARY}"
	@-pkill -SIGTERM -f "./${POD_BITP_BINARY}"
	@-pkill -SIGTERM -f "./${POD_MBTC_BINARY}"
	@-pkill -SIGTERM -f "./${POD_FOXB_BINARY}"
	@-pkill -SIGTERM -f "./${POD_RIPI_BINARY}"
	@echo "Stopped all pods!"

## start: start pod BINA
start-BINA: 
	@echo "Starting pod BINA"
	cd ../pods/bin && ./${POD_BINA_BINARY} &
	@echo "Done!"

## start: start pod MBTC
start-BITP: 
	@echo "Starting pod BITP"
	cd ../pods/bin && ./${POD_BITP_BINARY} &
	@echo "Done!"

## start: start pod MBTC
start-MBTC: 
	@echo "Starting pod MBTC"
	cd ../pods/bin && ./${POD_MBTC_BINARY} &
	@echo "Done!"

## start: start pod FOXB
start-FOXB: 
	@echo "Starting pod FOXB"
	cd ../pods/bin && ./${POD_FOXB_BINARY} &
	@echo "Done!"

## start: start pod RIPI
start-RIPI: 
	@echo "Starting pod RIPI"
	cd ../pods/bin && ./${POD_RIPI_BINARY} &
	@echo "Done!"


## build_all: build all pods
build-all: build-pod-BINA build-pod-BITP build-pod-FOXB build-pod-MBTC build-pod-RIPI
	@echo "Building all pods..."
	@echo "Done!"

## build-pod-BINA: builds the pod binary as a macos executable
build-pod-BINA:
	@echo "Building pod BINA binary..."
	cd ../pods && env GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -o ./bin/${POD_BINA_BINARY} ./cmd/client-BINA
	@echo "Done!"

## build-pod-BITP: builds the pod binary as a macos executable
build-pod-BITP:
	@echo "Building pod BITP binary..."
	cd ../pods && env GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -o ./bin/${POD_BITP_BINARY} ./cmd/client-BITP
	@echo "Done!"

## build-pod-MBTC: builds the pod binary as a macos executable
build-pod-MBTC:
	@echo "Building pod MBTC binary..."
	cd ../pods && env GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -o ./bin/${POD_MBTC_BINARY} ./cmd/client-MBTC
	@echo "Done!"

## build-pod-FOXB: builds the pod binary as a macos executable
build-pod-FOXB:
	@echo "Building pod FOXB binary..."
	cd ../pods && env GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -o ./bin/${POD_FOXB_BINARY} ./cmd/client-FOXB
	@echo "Done!"

## build-pod-RIPI: builds the pod binary as a macos executable
build-pod-RIPI:
	@echo "Building pod RIPI binary..."
	cd ../pods && env GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -o ./bin/${POD_RIPI_BINARY} ./cmd/client-RIPI
	@echo "Done!"


## GRPC SERVER
###################

## build_server: build gRPC server
build-server:
	@echo "Building gRPC server..."
	cd ../grpc-server && env GOOS=darwin GOARCH=arm64 CGO_ENABLED=1 go build -o ./bin/${GRPC_SERVER_BINARY} ./cmd/api
	@echo "Done!"

## start: start gRPC server
start-server: 
	@echo "Starting gRPC server"
	cd ../grpc-server/bin && ./${GRPC_SERVER_BINARY}
	@echo "Done!"

## stop: stop the gRPC server
stop-server:
	@echo "Stopping gRPC server..."
	@-pkill -SIGTERM -f "./${GRPC_SERVER_BINARY}"
	@echo "Stopped gRPC server!"

## DOCKER-BUILD
###############

## up: starts all containers in the background without forcing build
# up:
# 	@echo "Starting Docker images..."
# 	docker-compose up -d
# 	@echo "Docker images started!"

# ## up_build: stops docker-compose (if running), builds all projects and starts docker compose
# up_build: build_auth build_broker build_listener build_logger build_mail
# 	@echo "Stopping docker images (if running...)"
# 	docker-compose down
# 	@echo "Building (when required) and starting docker images..."
# 	docker-compose up --build -d
# 	@echo "Docker images built and started!"

# ## down: stop docker compose
# down:
# 	@echo "Stopping docker compose..."
# 	docker-compose down
# 	@echo "Done!"
