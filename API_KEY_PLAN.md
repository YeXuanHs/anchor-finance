# API密钥功能计划

## 一、背景

zjmf系统对接上游时，用`/v1/login_api`端点，传`account`（账号）+ `password`（API密钥）获取JWT。

锚点财务要兼容这个流程，同时锚点自己对接其他锚点也走同样的方式——**统一用账号+API密钥登录获取JWT，不管是zjmf对接还是锚点对接，流程完全一样。**

## 二、设计

### 2.1 API密钥 vs 登录密码

| | 登录密码 | API密钥 |
|---|---------|---------|
| 用途 | 人登录用（浏览器） | 程序对接用（API） |
| 生成方式 | 注册时设置 | 后台/前台生成 |
| 存储 | bcrypt哈希 | AES加密存储 |
| 格式 | 用户自定义 | 系统生成64位随机字符串 |
| 登录方式 | /api/client/login | /v1/login_api |

### 2.2 数据库

`users`表新增字段：
```sql
api_key VARCHAR(64) UNIQUE  -- API密钥（AES加密存储）
```

### 2.3 对接流程

**zjmf对接锚点 → 走zjmf兼容端点（/v1/）**
```
zjmf调 /v1/login_api → 获取JWT → 用JWT调 /v1/* 接口
```

**锚点对接锚点 → 走锚点自己的端点（/api/client/）**
```
锚点调 /api/client/auth/login-api → 获取JWT → 用JWT调 /api/client/* 接口
```

**锚点对接zjmf → 走zjmf的端点**
```
锚点调 zjmf的 /v1/login_api → 获取JWT → 用JWT调 zjmf的 /v1/* 接口
```

**三种场景都用同一套API密钥，只是登录端点不同。**

### 2.4 API开放条件（后台可配置，参考zjmf）

后台设置 → API开放条件：

| 配置项 | 可选值 | 说明 |
|--------|--------|------|
| API功能开关 | 开启/关闭 | 关闭后用户无法自助开通API，但管理员仍可操作 |
| 需要绑定手机 | 是/否 | 开启后用户必须绑定手机才能开通API |
| 需要实名认证 | 是/否 | 开启后用户必须实名认证才能开通API |

两个条件可以同时开（都要满足）或都不开（无条件开通）。

**用户自助开通流程：**
```
1. 用户进入"API管理"页面
2. 检查API功能是否开启（关闭则提示不可用）
3. 检查是否满足条件（绑定手机/实名认证）
4. 满足条件 → 点击"开通API" → 系统生成API密钥
5. 不满足条件 → 提示用户先完成绑定手机/实名认证
```

**管理员操作：**
```
1. 管理员进入用户详情页
2. 可以直接查看/重置/关闭用户的API密钥（无视条件）
3. 记录操作日志
```

### 2.5 API密钥管理

| 操作 | 说明 |
|------|------|
| 重置密钥 | 生成新密钥，旧密钥立即失效 |
| 关闭API | 停用API功能，密钥保留但不可用 |
| 自定义密钥 | ❌ 不支持，系统自动生成 |
| 多个密钥 | ❌ 不支持，一个用户只能一个 |

### 2.6 API端点

#### 用户自助操作（需JWT认证）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/client/api-key/status | 查看API状态（是否开通、密钥、条件） |
| POST | /api/client/api-key/enable | 自助开通API（检查条件） |
| POST | /api/client/api-key/reset | 重置API密钥 |
| POST | /api/client/api-key/disable | 关闭API功能 |

#### 管理员操作（需管理员JWT认证）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/admin/settings/api | 查看API开放条件配置 |
| PUT | /api/admin/settings/api | 修改API开放条件配置 |
| POST | /api/admin/users/:id/api-key/enable | 强制开通API（无视条件） |
| POST | /api/admin/users/:id/api-key/reset | 管理员重置用户API密钥 |
| POST | /api/admin/users/:id/api-key/disable | 管理员关闭用户API |

#### 登录获取JWT（公开端点）

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/client/auth/login-api | 账号+API密钥登录，返回JWT |
| POST | /v1/login_api | 账号+API密钥登录，返回JWT（兼容zjmf） |

### 2.7 对接场景

#### 场景A：zjmf对接锚点
```
zjmf后台 → 添加供应商 → 类型选"zjmf"
→ 填锚点网址 http://45.207.210.235:8080
→ 填账号（锚点用户名）+ API密钥（锚点生成的api_key）
→ zjmf调 /v1/login_api → 获取JWT
→ zjmf用JWT调 /v1/products、/v1/hosts/:id/module/status 等
```

#### 场景B：锚点对接锚点
```
锚点后台 → 添加供应商 → 类型选"anchor"
→ 填另一套锚点网址
→ 填账号 + API密钥
→ 调 /api/client/auth/login-api → 获取JWT
→ 用JWT调 /api/client/* 接口
```

#### 场景C：锚点对接zjmf
```
锚点后台 → 添加供应商 → 类型选"zjmf"
→ 填zjmf网址
→ 填账号 + API密钥（zjmf生成的）
→ 调 zjmf的 /v1/login_api → 获取JWT
→ 用JWT调 zjmf的接口
```

### 2.8 安全措施

- API密钥AES加密存储
- 查看密钥返回完整明文（不脱敏）
- 生成新密钥时旧密钥立即失效
- 无API密钥登录失败计数（不限制尝试次数）
- API密钥无有效期限制（用户可随时手动重置）

### 2.9 文件清单

| 文件 | 改动 |
|------|------|
| model/user.go | 新增api_key、api_enabled字段 |
| api/client/api_key.go | 新增：查看状态/自助开通/重置/关闭 |
| api/client/auth.go | 新增：LoginByAPIKey（/api/client/auth/login-api） |
| api/client/router.go | 注册新路由 |
| api/client/zjmf_compat.go | 改login_api用api_key验证 |
| api/admin/users.go | 新增：管理员强制开通/重置/关闭用户API |
| api/admin/settings.go | 新增：API开放条件配置 |
| cmd/server/main.go | 无需改动（AutoMigrate自动加字段） |

### 2.10 注意事项

- API密钥和登录密码是**两套独立凭证**，互不影响
- 用户可以没有API密钥（只用密码登录）
- 一个用户只能有一个API密钥（生成新的旧的失效）
- API密钥无有效期限制（用户可随时手动重置）
- 密钥不能自定义，系统自动生成
- 关闭API后密钥保留但不可用（重新开通恢复）
