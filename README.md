# 锚点财务 (AnchorFinance)

新一代财务管理系统，使用 Go + Vue 3 构建。

## 技术栈

### 后端
- **语言**: Go 1.21+
- **框架**: Gin
- **数据库**: PostgreSQL 14+
- **ORM**: GORM
- **缓存**: Redis 7+
- **认证**: JWT

### 前端
- **框架**: Vue 3 + TypeScript
- **构建**: Vite
- **UI**: Naive UI
- **状态**: Pinia
- **图表**: ECharts (管理后台)

## 项目结构

```
AnchorFinance/
├── backend/                # Go 后端
│   ├── cmd/server/        # 入口
│   ├── internal/          # 内部包
│   │   ├── api/          # API路由
│   │   ├── config/       # 配置
│   │   ├── handler/      # 处理器
│   │   ├── service/      # 业务逻辑
│   │   ├── model/        # 数据模型
│   │   └── job/          # 定时任务
│   ├── pkg/              # 公共包
│   ├── configs/          # 配置文件
│   └── install/          # 安装程序
├── frontend/
│   ├── web/              # 用户前台
│   └── admin/            # 管理后台
└── Makefile
```

## 快速开始

### 1. 安装依赖
```bash
make install
```

### 2. 配置
```bash
cd backend
cp configs/config.example.yaml configs/config.yaml
# 编辑 config.yaml 配置数据库和Redis
```

### 3. 启动后端
```bash
make dev
```

### 4. 启动前端
```bash
# 用户前台
make frontend

# 管理后台
make admin
```

### 5. 访问
- 用户前台: http://localhost:5173
- 管理后台: http://localhost:5174/admin
- API: http://localhost:8080/api/

## 安装程序

首次访问时会自动跳转到安装页面，按提示完成：
1. 环境检测
2. 数据库配置
3. 管理员设置
4. 完成安装

## API 接口

### V1 (兼容智简魔方)
- `POST /api/v1/user/login` - 登录
- `POST /api/v1/user/register` - 注册
- `GET /api/v1/products` - 产品列表
- `GET /api/v1/user/info` - 用户信息

### V2 (原生接口)
- `POST /api/v2/auth/login` - 登录
- `POST /api/v2/auth/register` - 注册
- `GET /api/v2/products` - 产品列表
- `GET /api/v2/orders` - 订单列表
- `GET /api/v2/invoices` - 账单列表
- `GET /api/v2/tickets` - 工单列表

## 默认账户

安装后首个管理员账户由安装程序创建。

## 许可证

MIT License
