BINARY=api

DOCKER_IMAGE_NAME=arizalsaputro/api

build: ## Build the binary
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags "-s" -o ${BINARY} .

clean: ## Clean up build artifacts
	go clean

docker-build: ## Build Docker image
	docker build -t ${DOCKER_IMAGE_NAME} .

docker-push: ## Push Docker image to registry
	docker push ${DOCKER_IMAGE_NAME}