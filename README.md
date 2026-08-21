# 锚点财务 (AnchorFinance)

一个兼容智简魔方(zjmf)接口的现代化财务管理系统。

## 项目结构

```
anchor-finance/
├── backend/                # Go后端 (Gin + GORM)
│   ├── cmd/server/        # 主入口
│   ├── config/            # 配置
│   ├── internal/          # 内部包
│   │   ├── api/           # API处理器
│   │   ├── model/         # 数据模型
│   │   ├── service/       # 业务服务
│   │   ├── middleware/    # 中间件
│   │   └── database/      # 数据库连接
│   └── pkg/               # 公共包
│
├── plugin-engine/          # PHP插件引擎 (Laravel)
│   ├── app/               # 应用代码
│   │   ├── ZjmfCompat/    # zjmf兼容层
│   │   ├── Http/          # HTTP控制器
│   │   └── Services/      # 服务
│   ├── plugins/           # 插件目录
│   └── zjmf_compat/       # zjmf函数兼容
│
├── admin-frontend/         # Admin后台 (Vue 3 + Arco Design Pro)
│   └── src/
│       ├── api/           # API接口
│       ├── components/    # 组件
│       ├── layouts/       # 布局
│       ├── views/         # 页面
│       ├── router/        # 路由
│       └── store/         # 状态管理
│
├── client-frontend/        # 用户前台 (Vue 3 + Element Plus)
│   └── src/
│       ├── api/           # API接口
│       ├── components/    # 组件
│       ├── layouts/       # 布局
│       ├── views/         # 页面
│       ├── router/        # 路由
│       └── store/         # 状态管理
│
└── .github/workflows/     # CI/CD配置
```

## 技术栈

- **后端**: Go 1.21+ / Gin / GORM
- **插件引擎**: PHP 8.1+ / Laravel 10
- **Admin后台**: Vue 3.5 / Arco Design Pro / Pinia
- **用户前台**: Vue 3.5 / Element Plus / Pinia
- **数据库**: MySQL 8.0
- **缓存**: Redis 7

## 快速开始

### 1. 克隆项目

```bash
git clone https://github.com/YeXuanHs/anchor-finance.git
cd anchor-finance
```

### 2. 配置环境

复制环境变量示例文件：

```bash
cp backend/.env.example backend/.env
```

编辑 `.env` 文件，配置数据库和Redis连接信息。

### 3. 启动后端

```bash
cd backend
go mod tidy
go run cmd/server/main.go
```

后端将在 http://localhost:8080 启动。

### 4. 启动Admin前端

```bash
cd admin-frontend
npm install
npm run dev
```

Admin前端将在 http://localhost:3000 启动。

### 5. 启动用户前台

```bash
cd client-frontend
npm install
npm run dev
```

用户前台将在 http://localhost:3001 启动。

### 6. 启动PHP插件引擎

```bash
cd plugin-engine
composer install
php artisan serve --port=9000
```

插件引擎将在 http://localhost:9000 启动。

## 默认账号

- **Admin后台**: admin / admin123
- **用户前台**: 需要先注册

## API文档

### Admin API

- 基础路径: `/api/admin/`
- 认证方式: Bearer Token (JWT)

### Client API

- 基础路径: `/api/client/`
- 认证方式: Bearer Token (JWT)

## 插件系统

### 支持的插件域

| 域 | 目录 | 说明 |
|----|------|------|
| payment | gateways/ | 支付网关 |
| sms | sms/ | 短信服务 |
| mail | mail/ | 邮件服务 |
| oauth | oauth/ | 第三方登录 |
| servers | servers/ | 上游接口/自动化接口 |

### zjmf插件兼容

zjmf插件可以直接放入 `plugin-engine/plugins/` 目录，无需修改即可使用。

## 部署

### 测试服务器

- **地址**: 45.207.210.235
- **后台**: http://45.207.210.235:8080/SB111111/

### GitHub Actions

项目使用GitHub Actions自动构建和部署：

1. 推送代码到 `main` 分支
2. GitHub Actions自动构建Go二进制和前端资源
3. 自动部署到服务器

## 开发规范

- **禁止创建stub/TODO文件**: 写一个文件就写完整
- **每个文件必须有实际逻辑**: 不允许只写函数签名+TODO注释
- **SSH连接**: 使用paramiko写Python脚本，不用MCP SSH工具

## 安全特性

- 参数化查询防止SQL注入
- 价格服务端计算防止0元购
- JWT认证与权限控制
- 频率限制防止暴力破解
- 安全审计日志

## 许可证

MIT License
