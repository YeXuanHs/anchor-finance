.PHONY: build run test clean dev frontend frontend-build

# 后端
build:
	cd backend && go build -o bin/anchorfinance .

run: build
	cd backend && ./bin/anchorfinance

dev:
	cd backend && go run .

test:
	cd backend && go test ./...

clean:
	rm -rf backend/bin/

# 前端 - 用户前台
frontend:
	cd frontend/web && npm install && npm run dev

frontend-build:
	cd frontend/web && npm install && npm run build

# 前端 - 管理后台
admin:
	cd frontend/admin && npm install && npm run dev

admin-build:
	cd frontend/admin && npm install && npm run build

# 全部构建
all: build frontend-build admin-build

# 安装依赖
install:
	cd backend && go mod tidy
	cd frontend/web && npm install
	cd frontend/admin && npm install
