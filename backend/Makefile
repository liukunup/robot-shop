.PHONY: init
init:
	go install github.com/google/wire/cmd/wire@latest
	go install github.com/golang/mock/mockgen@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/go-delve/delve/cmd/dlv@latest
	go install github.com/go-nunu/nunu@latest

.PHONY: bootstrap
bootstrap:
	cd ./deploy/docker-compose && docker compose up -d && cd ../../
	go run ./cmd/migration
	nunu run ./cmd/server

.PHONY: mock
mock:
	@# 生成 repository 的 mock
	@mkdir -p test/mocks/repository
	@for file in internal/repository/*.go; do \
		if [ -f "$$file" ]; then \
			mockgen -source=$$file -destination=test/mocks/repository/$$(basename $$file) || exit 1; \
			echo "Generated mock repository for $$file"; \
		fi \
	done
	@# 生成 service 的 mock
	@mkdir -p test/mocks/service
	@for file in internal/service/*.go; do \
		if [ -f "$$file" ] && [ "$$file" != "internal/service/service.go" ]; then \
			mockgen -source=$$file -destination=test/mocks/service/$$(basename $$file) || exit 1; \
			echo "Generated mock service for $$file"; \
		fi \
	done

.PHONY: test
test:
	go test -coverpkg=./internal/handler,./internal/service,./internal/repository -coverprofile=./coverage.out ./test/server/...
	go tool cover -html=./coverage.out -o coverage.html

.PHONY: build
build:
	go build -ldflags="-s -w" -o ./bin/server ./cmd/server

.PHONY: docker
docker:
	docker build -f deploy/build/Dockerfile --build-arg APP_RELATIVE_PATH=./cmd/task -t 1.1.1.1:5000/demo-task:v1 .
	docker run --rm -i 1.1.1.1:5000/demo-task:v1

.PHONY: swag
swag:
	swag init -g cmd/server/main.go -o docs --parseDependency
