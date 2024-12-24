APP_NAME = go-auth-service
DOCKER_IMAGE = $(APP_NAME):latest
DOCKER_COMPOSE = docker-compose.yml

.PHONY: all
all: build

.PHONY: build
build:
    @echo "Building the application..."
    go build -o $(APP_NAME) ./cmd/main.go

.PHONY: run
run: build
    @echo "Running the application..."
    ./$(APP_NAME)

.PHONY: clean
clean:
    @echo "Cleaning the build files..."
    rm -f $(APP_NAME)

.PHONY: test
test:
    @echo "Running tests..."
    go test ./... -v

.PHONY: docker-build
docker-build:
    @echo "Building Docker image..."
    docker build -t $(DOCKER_IMAGE) .

.PHONY: docker-compose-up
docker-compose-up:
    @echo "Starting Docker Compose..."
    docker-compose -f $(DOCKER_COMPOSE) up --build

.PHONY: docker-compose-down
docker-compose-down:
    @echo "Stopping Docker Compose..."
    docker-compose -f $(DOCKER_COMPOSE) down

.PHONY: logs
logs:
    @echo "Tailing Docker Compose logs..."
    docker-compose -f $(DOCKER_COMPOSE) logs -f

.PHONY: docker-run
docker-run: docker-build
    @echo "Running the application in a Docker container..."
    docker run -p 3000:3000 $(DOCKER_IMAGE)

.PHONY: help
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
