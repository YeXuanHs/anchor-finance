# 锚点财务 API清单

> 每完成一个功能必须更新此文件
> 禁止同一个功能有多个API入口

## 项目结构

```
backend/
├── cmd/server/main.go              # 入口
├── config/config.go                # 配置
├── internal/
│   ├── api/
│   │   ├── admin/                  # 管理后台API
│   │   │   ├── auth.go             # 认证模块
│   │   │   ├── users.go            # 用户管理
│   │   │   ├── orders.go           # 订单管理
│   │   │   ├── services.go         # 服务管理
│   │   │   ├── invoices.go         # 账单管理
│   │   │   ├── tickets.go          # 工单管理
│   │   │   ├── products.go         # 产品管理
│   │   │   ├── plugins.go          # 插件管理
│   │   │   └── settings.go         # 设置管理
│   │   └── client/                 # 用户前台API
│   │       ├── auth.go             # 认证模块
│   │       ├── user.go             # 用户信息
│   │       ├── services.go         # 我的服务
│   │       ├── orders.go           # 我的订单
│   │       ├── tickets.go          # 我的工单
│   │       └── finance.go          # 财务中心
│   ├── model/                      # 数据模型
│   ├── service/                    # 业务服务
│   ├── middleware/                  # 中间件
│   └── database/                   # 数据库连接
```

---

## 已实现功能清单

### 认证模块 (Admin)

| API | 文件 | 状态 | 说明 |
|-----|------|------|------|
| POST /api/admin/login | internal/api/admin/auth.go | ✅完成 | 管理员登录（带防暴力破解） |
| POST /api/admin/logout | internal/api/admin/auth.go | ✅完成 | 管理员登出 |
| GET /api/admin/auth/info | internal/api/admin/auth.go | ✅完成 | 获取管理员信息 |

### 认证模块 (Client)

| API | 文件 | 状态 | 说明 |
|-----|------|------|------|
| POST /api/client/login | internal/api/client/auth.go | ⏳待实现 | 用户登录 |
| POST /api/client/register | internal/api/client/auth.go | ⏳待实现 | 用户注册 |
| POST /api/client/password/reset | internal/api/client/auth.go | ⏳待实现 | 重置密码 |

### 用户管理 (Admin)

| API | 文件 | 状态 | 说明 |
|-----|------|------|------|
| GET /api/admin/users | internal/api/admin/users.go | ✅完成 | 用户列表（分页+搜索） |
| GET /api/admin/users/:id | internal/api/admin/users.go | ✅完成 | 用户详情 |
| POST /api/admin/users | internal/api/admin/users.go | ✅完成 | 创建用户 |
| PUT /api/admin/users/:id | internal/api/admin/users.go | ✅完成 | 更新用户 |
| DELETE /api/admin/users/:id | internal/api/admin/users.go | ✅完成 | 删除用户（软删除） |
| GET /api/admin/users/:id/orders | internal/api/admin/users.go | ✅完成 | 用户订单 |
| GET /api/admin/users/:id/invoices | internal/api/admin/users.go | ✅完成 | 用户账单 |
| GET /api/admin/users/:id/tickets | internal/api/admin/users.go | ✅完成 | 用户工单 |
| GET /api/admin/users/:id/services | internal/api/admin/users.go | ✅完成 | 用户服务 |

### 订单管理 (Admin)

| API | 文件 | 状态 | 说明 |
|-----|------|------|------|
| GET /api/admin/orders | internal/api/admin/orders.go | ✅完成 | 订单列表（分页+搜索+筛选） |
| GET /api/admin/orders/:id | internal/api/admin/orders.go | ✅完成 | 订单详情 |
| POST /api/admin/orders | internal/api/admin/orders.go | ✅完成 | 创建订单 |
| PUT /api/admin/orders/:id | internal/api/admin/orders.go | ✅完成 | 更新订单 |
| POST /api/admin/orders/:id/activate | internal/api/admin/orders.go | ✅完成 | 激活订单 |
| POST /api/admin/orders/:id/cancel | internal/api/admin/orders.go | ✅完成 | 取消订单 |

### 服务管理 (Admin)

| API | 文件 | 状态 | 说明 |
|-----|------|------|------|
| GET /api/admin/services | internal/api/admin/services.go | ✅完成 | 服务列表（分页+搜索+筛选） |
| GET /api/admin/services/:id | internal/api/admin/services.go | ✅完成 | 服务详情 |
| PUT /api/admin/services/:id | internal/api/admin/services.go | ✅完成 | 更新服务 |
| POST /api/admin/services/:id/suspend | internal/api/admin/services.go | ✅完成 | 暂停服务 |
| POST /api/admin/services/:id/unsuspend | internal/api/admin/services.go | ✅完成 | 取消暂停 |
| POST /api/admin/services/:id/terminate | internal/api/admin/services.go | ✅完成 | 终止服务 |

### 账单管理 (Admin)

| API | 文件 | 状态 | 说明 |
|-----|------|------|------|
| GET /api/admin/invoices | internal/api/admin/invoices.go | ⏳待实现 | 账单列表 |
| GET /api/admin/invoices/:id | internal/api/admin/invoices.go | ⏳待实现 | 账单详情 |
| POST /api/admin/invoices/:id/cancel | internal/api/admin/invoices.go | ⏳待实现 | 取消账单 |
| GET /api/admin/transactions | internal/api/admin/invoices.go | ⏳待实现 | 交易流水 |

### 工单管理 (Admin)

| API | 文件 | 状态 | 说明 |
|-----|------|------|------|
| GET /api/admin/tickets | internal/api/admin/tickets.go | ⏳待实现 | 工单列表 |
| GET /api/admin/tickets/:id | internal/api/admin/tickets.go | ⏳待实现 | 工单详情 |
| POST /api/admin/tickets/:id/reply | internal/api/admin/tickets.go | ⏳待实现 | 回复工单 |
| POST /api/admin/tickets/:id/close | internal/api/admin/tickets.go | ⏳待实现 | 关闭工单 |
| GET /api/admin/ticket-departments | internal/api/admin/tickets.go | ⏳待实现 | 工单部门 |
| GET /api/admin/ticket-statuses | internal/api/admin/tickets.go | ⏳待实现 | 工单状态 |

### 产品管理 (Admin)

| API | 文件 | 状态 | 说明 |
|-----|------|------|------|
| GET /api/admin/products | internal/api/admin/products.go | ⏳待实现 | 产品列表 |
| GET /api/admin/products/:id | internal/api/admin/products.go | ⏳待实现 | 产品详情 |
| POST /api/admin/products | internal/api/admin/products.go | ⏳待实现 | 创建产品 |
| PUT /api/admin/products/:id | internal/api/admin/products.go | ⏳待实现 | 更新产品 |
| DELETE /api/admin/products/:id | internal/api/admin/products.go | ⏳待实现 | 删除产品 |
| GET /api/admin/product-groups | internal/api/admin/products.go | ⏳待实现 | 产品分组 |

### 插件管理 (Admin)

| API | 文件 | 状态 | 说明 |
|-----|------|------|------|
| GET /api/admin/plugins | internal/api/admin/plugins.go | ⏳待实现 | 插件列表 |
| POST /api/admin/plugins/install | internal/api/admin/plugins.go | ⏳待实现 | 安装插件 |
| POST /api/admin/plugins/scan | internal/api/admin/plugins.go | ⏳待实现 | 扫描插件 |
| POST /api/admin/plugins/:id/enable | internal/api/admin/plugins.go | ⏳待实现 | 启用插件 |
| POST /api/admin/plugins/:id/disable | internal/api/admin/plugins.go | ⏳待实现 | 禁用插件 |
| DELETE /api/admin/plugins/:id | internal/api/admin/plugins.go | ⏳待实现 | 卸载插件 |
| GET /api/admin/plugins/:id/config | internal/api/admin/plugins.go | ⏳待实现 | 获取配置 |
| PUT /api/admin/plugins/:id/config | internal/api/admin/plugins.go | ⏳待实现 | 更新配置 |

### 设置管理 (Admin)

| API | 文件 | 状态 | 说明 |
|-----|------|------|------|
| GET /api/admin/settings | internal/api/admin/settings.go | ⏳待实现 | 获取设置 |
| PUT /api/admin/settings | internal/api/admin/settings.go | ⏳待实现 | 更新设置 |
| GET /api/admin/settings/:group | internal/api/admin/settings.go | ⏳待实现 | 获取分组设置 |
| GET /api/admin/menus | internal/api/admin/settings.go | ⏳待实现 | 菜单列表 |
| POST /api/admin/menus | internal/api/admin/settings.go | ⏳待实现 | 创建菜单 |
| PUT /api/admin/menus/:id | internal/api/admin/settings.go | ⏳待实现 | 更新菜单 |
| DELETE /api/admin/menus/:id | internal/api/admin/settings.go | ⏳待实现 | 删除菜单 |

### 用户前台

| API | 文件 | 状态 | 说明 |
|-----|------|------|------|
| POST /api/client/login | internal/api/client/auth.go | ⏳待实现 | 用户登录 |
| POST /api/client/register | internal/api/client/auth.go | ⏳待实现 | 用户注册 |
| GET /api/client/user/info | internal/api/client/user.go | ⏳待实现 | 获取用户信息 |
| PUT /api/client/user/info | internal/api/client/user.go | ⏳待实现 | 更新用户信息 |
| GET /api/client/services | internal/api/client/services.go | ⏳待实现 | 我的服务列表 |
| GET /api/client/orders | internal/api/client/orders.go | ⏳待实现 | 我的订单列表 |
| GET /api/client/tickets | internal/api/client/tickets.go | ⏳待实现 | 我的工单列表 |
| GET /api/client/invoices | internal/api/client/finance.go | ⏳待实现 | 我的账单列表 |

---

## 数据库表清单

| 表名 | Model文件 | 状态 | 说明 |
|------|----------|------|------|
| users | internal/model/user.go | ✅完成 | 用户表 |
| admins | internal/model/admin.go | ✅完成 | 管理员表（含防暴力破解字段） |
| roles | internal/model/admin.go | ✅完成 | 角色表 |
| orders | internal/model/order.go | ✅完成 | 订单表 |
| services | internal/model/service.go | ✅完成 | 服务表 |
| invoices | internal/model/invoice.go | ✅完成 | 账单表 |
| tickets | internal/model/ticket.go | ✅完成 | 工单表 |
| products | internal/model/product.go | ⏳待实现 | 产品表 |
| plugins | internal/model/plugin.go | ⏳待实现 | 插件表 |
| settings | internal/model/setting.go | ⏳待实现 | 设置表 |
| menus | internal/model/menu.go | ⏳待实现 | 菜单表 |

---

## 更新日志

### 2024-01-21
- ✅ 初始化项目结构
- ✅ 完成认证模块（Admin登录/登出/获取信息）
- ✅ 完成User/Admin/Role模型
- ✅ Go后端部署成功，健康检查API正常
- ✅ 用户管理CRUD完整实现
- ✅ 管理员登录防暴力破解（5次失败冻结6小时）
- ✅ 订单管理API（6个接口）
- ✅ 服务管理API（6个接口）
- ✅ 账单管理API（4个接口）
- ✅ 工单管理API（6个接口）
- ✅ 产品管理API（6个接口）
- ✅ 设置管理API（7个接口）
- ✅ 管理员管理API（5个接口）
- ✅ 仪表盘API（3个接口）
- ✅ 所有API测试通过

### 实名认证
| API | 文件 | 状态 | 说明 |
|-----|------|------|------|
| GET /api/admin/verifications | internal/api/admin/verifications.go | ✅完成 | 认证列表 |
| GET /api/admin/verifications/summary | internal/api/admin/verifications.go | ✅完成 | 认证统计 |
| GET /api/admin/verifications/:id | internal/api/admin/verifications.go | ✅完成 | 认证详情 |
| POST /api/admin/verifications/:id/approve | internal/api/admin/verifications.go | ✅完成 | 批准认证 |
| POST /api/admin/verifications/:id/reject | internal/api/admin/verifications.go | ✅完成 | 拒绝认证 |

### 供应商管理
| API | 文件 | 状态 | 说明 |
|-----|------|------|------|
| GET /api/admin/suppliers | internal/api/admin/suppliers.go | ✅完成 | 供应商列表 |
| GET /api/admin/suppliers/summary | internal/api/admin/suppliers.go | ✅完成 | 供应商统计 |
| GET /api/admin/suppliers/:id | internal/api/admin/suppliers.go | ✅完成 | 供应商详情 |
| POST /api/admin/suppliers | internal/api/admin/suppliers.go | ✅完成 | 创建供应商 |
| PUT /api/admin/suppliers/:id | internal/api/admin/suppliers.go | ✅完成 | 更新供应商 |
| DELETE /api/admin/suppliers/:id | internal/api/admin/suppliers.go | ✅完成 | 删除供应商 |
| GET /api/admin/suppliers/:id/products | internal/api/admin/suppliers.go | ✅完成 | 供应商产品 |

### 插件管理
| API | 文件 | 状态 | 说明 |
|-----|------|------|------|
| GET /api/admin/plugins | internal/api/admin/plugins.go | ✅完成 | 插件列表 |
| GET /api/admin/plugins/:id | internal/api/admin/plugins.go | ✅完成 | 插件详情 |
| POST /api/admin/plugins/:id/enable | internal/api/admin/plugins.go | ✅完成 | 启用插件 |
| POST /api/admin/plugins/:id/disable | internal/api/admin/plugins.go | ✅完成 | 禁用插件 |
| DELETE /api/admin/plugins/:id | internal/api/admin/plugins.go | ✅完成 | 卸载插件 |
| GET /api/admin/plugins/:id/config | internal/api/admin/plugins.go | ✅完成 | 获取配置 |
| PUT /api/admin/plugins/:id/config | internal/api/admin/plugins.go | ✅完成 | 更新配置 |

### 待完成
- ⏳ 优惠码管理API
- ⏳ 前端开发（所有API完成后）

---

## 状态说明

- ✅完成 - 代码完整，已测试
- ⏳待实现 - 未开始开发
- 🔧开发中 - 正在开发
- ❌废弃 - 已废弃，不再使用
