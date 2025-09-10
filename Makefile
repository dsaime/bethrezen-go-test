.PHONY: test vet lint check run build genmock compose-up compose-build

# Get number of CPU cores minus 1 for parallel execution
CORES := $(shell echo $$(( $$(nproc) - 1 )))

test:
	echo "Running go test..."
	go test ./...

vet:
	echo "Running go vet..."
	go vet ./...

lint:
	echo "Running golangci-lint with $(CORES) workers..."
# The -j parameter for golangci-lint will use all available CPU cores minus one (to avoid overloading your system)
	#go run github.com/golangci/golangci-lint/cmd/golangci-lint@latest run -v -j $(CORES)
	golangci-lint run -v -j $(CORES)

# Combined target to run both vet and lint
check: vet lint

run:
	echo "Running npc main package..."
	go run ./cmd/newsapi/main.go

build:
	echo "Running newsapi main package..."
	CGO_ENABLED=0 GOOS=linux go build -o bin/newsapi cmd/newsapi/main.go

genmock:
	go run github.com/vektra/mockery/v3@v3.5.4

compose-up:
	docker compose -f infra/docker-compose.yaml -p infra up -d dsaime.test.newsapi

compose-build:
	docker compose -f infra/docker-compose.yaml build dsaime.test.newsapi