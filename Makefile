# Makefile for Cloud Storage Service

.PHONY: help build run test clean migrate docker-up docker-down lint format

# 默认目标
help:
	@echo "可用命令:"
	@echo "  make build      - 构建应用程序"
	@echo "  make run        - 运行应用程序"
	@echo "  make test       - 运行测试"
	@echo "  make clean      - 清理构建文件"
	@echo "  make migrate    - 运行数据库迁移"
	@echo "  make docker-up  - 启动Docker容器（全部）"
	@echo "  make docker-backend    - 启动后端服务"
	@echo "  make docker-frontend   - 启动前端服务"
	@echo "  make docker-full      - 启动完整服务（包括Nginx）"
	@echo "  make docker-down - 停止Docker容器"
	@echo "  make lint       - 运行代码检查"
	@echo "  make format     - 格式化代码"
	@echo "  make dev        - 开发模式运行"

# 构建应用程序
build:
	@echo "构建应用程序..."
	go build -o cloud-storage ./cmd/server

# 运行应用程序
run: build
	@echo "启动应用程序..."
	./cloud-storage

# 开发模式运行
dev:
	@echo "开发模式启动..."
	go run ./cmd/server

# 运行测试
test:
	@echo "运行测试..."
	go test ./... -v

# 清理构建文件
clean:
	@echo "清理构建文件..."
	rm -f cloud-storage
	rm -rf dist/
	rm -rf coverage.out

# 数据库迁移
migrate:
	@echo "运行数据库迁移..."
	go run ./cmd/migrate

# 启动Docker容器（全部）
docker-up:
	@echo "启动Docker容器（全部）..."
	docker-compose --profile backend --profile frontend up -d

# 启动后端服务
docker-backend:
	@echo "启动后端服务..."
	docker-compose --profile backend up -d

# 启动前端服务
docker-frontend:
	@echo "启动前端服务..."
	docker-compose --profile frontend up -d

# 启动完整服务（包括Nginx）
docker-full:
	@echo "启动完整服务（包括Nginx代理）..."
	docker-compose --profile backend --profile frontend --profile full up -d

# 停止Docker容器
docker-down:
	@echo "停止Docker容器..."
	docker-compose down

# 清理Docker容器和卷
docker-clean:
	@echo "清理Docker容器和卷..."
	docker-compose down -v
	docker system prune -f

# 运行代码检查
lint:
	@echo "运行代码检查..."
	golangci-lint run ./...

# 格式化代码
format:
	@echo "格式化代码..."
	go fmt ./...

# 生成Swagger文档
swagger:
	@echo "生成Swagger文档..."
	swag init -g ./cmd/server/main.go -o ./api/docs

# 安装开发依赖
deps:
	@echo "安装开发依赖..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/swaggo/swag/cmd/swag@latest

# 创建.env文件
env:
	@echo "创建.env文件..."
	cp .env.example .env
	@echo "请编辑.env文件以配置环境变量"

# 数据库备份
db-backup:
	@echo "备份数据库..."
	docker-compose exec postgres pg_dump -U postgres cloud_storage > backup_$(shell date +%Y%m%d_%H%M%S).sql

# 数据库恢复
db-restore:
	@echo "恢复数据库..."
	@read -p "输入备份文件路径: " file; \
	docker-compose exec -T postgres psql -U postgres cloud_storage < $$file

# 查看日志
logs:
	@echo "查看应用程序日志..."
	docker-compose logs -f app

# 查看数据库日志
logs-db:
	@echo "查看数据库日志..."
	docker-compose logs -f postgres

# 查看Redis日志
logs-redis:
	@echo "查看Redis日志..."
	docker-compose logs -f redis

# 进入数据库容器
db-shell:
	@echo "进入数据库容器..."
	docker-compose exec postgres psql -U postgres cloud_storage

# 进入Redis容器
redis-shell:
	@echo "进入Redis容器..."
	docker-compose exec redis redis-cli -a redispassword

# 进入应用程序容器
app-shell:
	@echo "进入应用程序容器..."
	docker-compose exec app sh

# 运行压力测试
bench:
	@echo "运行压力测试..."
	ab -n 1000 -c 100 http://localhost:8080/health

# 代码覆盖率
coverage:
	@echo "生成代码覆盖率报告..."
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "覆盖率报告已生成: coverage.html"

# 依赖更新
update-deps:
	@echo "更新依赖..."
	go get -u ./...
	go mod tidy

# 安全检查
security:
	@echo "运行安全检查..."
	gosec ./...
	nancy audit

# 构建Docker镜像
docker-build:
	@echo "构建Docker镜像..."
	docker build -t cloud-storage:latest .

# 推送Docker镜像
docker-push:
	@echo "推送Docker镜像..."
	docker tag cloud-storage:latest your-registry/cloud-storage:latest
	docker push your-registry/cloud-storage:latest

# 重新构建Web前端
docker-rebuild-web:
	@echo "重新构建Web前端..."
	docker-compose build --no-cache web
	docker-compose up -d web

# 初始化项目
init: env deps
	@echo "项目初始化完成"
	@echo "请运行 'make docker-up' 启动服务"
	@echo "或运行 'make dev' 在开发模式运行"

# 默认目标
.DEFAULT_GOAL := help