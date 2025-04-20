APP_NAME = go-auth-service
DOCKER_IMAGE = $(APP_NAME):latest
DOCKER_COMPOSE = docker-compose.yml

build:
	@echo "Building the application..."
	go build -o $(APP_NAME) ./main.go

run: build
	@echo "Running the application..."
	./$(APP_NAME)

clean:
	@echo "Cleaning the build files..."
	rm -f $(APP_NAME)

test:
	@echo "Running tests..."
	go test ./... -v

docker-build:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE) .

docker-compose-up:
	@echo "Starting Docker Compose..."
	docker-compose -f $(DOCKER_COMPOSE) up --build -d

docker-compose-down:
	@echo "Stopping Docker Compose..."
	docker-compose -f $(DOCKER_COMPOSE) down

logs:
	@echo "Tailing Docker Compose logs..."
	docker-compose -f $(DOCKER_COMPOSE) logs -f

docker-run: docker-build
	@echo "Running the application in a Docker container..."
	docker run -p 3000:3000 $(DOCKER_IMAGE)

help:
	@echo "Usage:"
	@echo "  make build                Build the Go application"
	@echo "  make run                  Run the Go application"
	@echo "  make clean                Clean the build files"
	@echo "  make test                 Run tests"
	@echo "  make docker-build         Build Docker image"
	@echo "  make docker-compose-up    Start Docker Compose"
	@echo "  make docker-compose-down  Stop Docker Compose"
	@echo "  make logs                 Tail Docker Compose logs"
	@echo "  make docker-run           Run the application in a Docker container"
	@echo "  make help                 Show this help message"
