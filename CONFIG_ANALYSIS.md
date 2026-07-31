# 魔方财务配置项分析

## 一、配置系统架构

### 魔方财务配置存储
- 配置表：`configuration` (MySQL)
- 字段：`setting` (配置键), `value` (配置值)
- 访问方式：`configuration("key")` 获取配置

### 锚点财务当前实现
- 基础配置：`config.yaml` (数据库连接、JWT等)
- 业务配置：`system_settings` 表
- 验证码配置：`captcha_configs` 表

---

## 二、配置项分类及影响范围

### 1. 验证码配置 (已实现 ✅)

| 配置项 | 说明 | 影响范围 |
|--------|------|----------|
| `is_captcha` | 验证码总开关 | 所有验证码场景 |
| `captcha_length` | 验证码长度(4/5/6) | 图形验证码生成 |
| `captcha_combination` | 验证码类型(数字/字母/混合) | 图形验证码生成 |
| `captcha_configuration` | 验证码高级配置(JSON) | 验证码样式、干扰线等 |
| `allow_register_email_captcha` | 邮件注册显示验证码 | 注册页面 |
| `allow_register_phone_captcha` | 手机注册显示验证码 | 注册页面 |
| `allow_login_phone_captcha` | 手机登录显示验证码 | 登录页面 |
| `allow_login_email_captcha` | 邮件登录显示验证码 | 登录页面 |
| `allow_login_code_captcha` | 验证码登录显示验证码 | 登录页面 |
| `allow_login_id_captcha` | ID登录显示验证码 | 登录页面 |
| `allow_login_admin_captcha` | 后台登录显示验证码 | 后台登录页面 |
| `allow_phone_forgetpwd_captcha` | 手机找回密码显示验证码 | 找回密码页面 |
| `allow_email_forgetpwd_captcha` | 邮件找回密码显示验证码 | 找回密码页面 |
| `allow_resetpwd_captcha` | 重置密码显示验证码 | 重置密码页面 |
| `allow_setpwd_captcha` | 设置密码显示验证码 | 设置密码页面 |
| `allow_phone_bind_captcha` | 手机绑定显示验证码 | 绑定页面 |
| `allow_email_bind_captcha` | 邮件绑定显示验证码 | 绑定页面 |
| `allow_cancel_sms_captcha` | 取消短信提醒显示验证码 | 设置页面 |
| `allow_cancel_email_captcha` | 取消邮件提醒显示验证码 | 设置页面 |

### 2. 登录/注册配置 (部分实现 ⚠️)

| 配置项 | 说明 | 影响范围 | 锚点状态 |
|--------|------|----------|----------|
| `allow_phone` | 允许手机登录 | 登录页面 | ❌ 缺失 |
| `allow_email` | 允许邮箱登录 | 登录页面 | ❌ 缺失 |
| `allow_register_phone` | 允许手机注册 | 注册页面 | ❌ 缺失 |
| `allow_register_email` | 允许邮箱注册 | 注册页面 | ❌ 缺失 |
| `allow_register_wechat` | 允许微信注册 | 注册页面 | ❌ 缺失 |
| `allow_login_phone` | 允许手机登录 | 登录页面 | ❌ 缺失 |
| `allow_login_email` | 允许邮箱登录 | 登录页面 | ❌ 缺失 |
| `allow_login_wechat` | 允许微信登录 | 登录页面 | ❌ 缺失 |
| `login_error_max_num` | 登录错误次数限制 | 登录逻辑 | ❌ 缺失 |
| `login_error_switch` | 登录错误限制开关 | 登录逻辑 | ❌ 缺失 |
| `clients_profoptional` | 用户资料可选字段 | 注册/个人资料 | ❌ 缺失 |
| `login_register_custom_require` | 自定义注册字段(JSON) | 注册页面 | ❌ 缺失 |

### 3. 安全配置 (部分实现 ⚠️)

| 配置项 | 说明 | 影响范围 | 锚点状态 |
|--------|------|----------|----------|
| `required_pwstrength` | 密码强度要求 | 注册/修改密码 | ❌ 缺失 |
| `invalid_logins_banlength` | 登录失败封禁时长 | 登录逻辑 | ❌ 缺失 |
| `home_ip_check` | 前台登录IP检测 | 登录安全 | ❌ 缺失 |
| `admin_ip_check` | 后台登录IP检测 | 后台登录安全 | ❌ 缺失 |
| `second_verify_action_home` | 前台二次验证操作 | 安全验证 | ❌ 缺失 |
| `second_verify_action_admin` | 后台二次验证操作 | 安全验证 | ❌ 缺失 |
| `second_verify_action_home_type` | 前台二次验证类型 | 安全验证 | ❌ 缺失 |

### 4. 公司/网站信息配置 (部分实现 ⚠️)

| 配置项 | 说明 | 影响范围 | 锚点状态 |
|--------|------|----------|----------|
| `company_name` | 公司名称 | 全站显示 | ✅ 已有 |
| `company_email` | 公司邮箱 | 联系信息 | ✅ 已有 |
| `domain` | 网站域名 | 系统链接 | ✅ 已有 |
| `logo_url` | Logo地址 | 网站头部 | ✅ 已有 |
| `logo_url_home` | 前台Logo地址 | 前台头部 | ❌ 缺失 |
| `logo_url_bill` | 账单Logo地址 | 账单/PDF | ❌ 缺失 |
| `www_logo` | 官网LOGO | 官网 | ❌ 缺失 |
| `system_url` | 系统链接 | 邮件/通知 | ✅ 已有 |
| `main_phone` | 联系电话 | 联系信息 | ✅ 已有 |
| `main_address` | 公司地址 | 联系信息 | ✅ 已有 |
| `record_no` | 备案号 | 网站底部 | ✅ 已有 |
| `map` | 坐标 | 地图显示 | ❌ 缺失 |
| `company_profile` | 公司简介 | 关于页面 | ❌ 缺失 |
| `seo_keywords` | SEO关键词 | 网站SEO | ❌ 缺失 |
| `seo_desc` | SEO描述 | 网站SEO | ❌ 缺失 |
| `server_clause_url` | 服务条款地址 | 注册页面 | ❌ 缺失 |
| `privacy_clause_url` | 隐私条款地址 | 注册页面 | ❌ 缺失 |
| `cancellation_time` | 注销时间 | 用户注销 | ❌ 缺失 |

### 5. 界面/模板配置 (缺失 ❌)

| 配置项 | 说明 | 影响范围 | 锚点状态 |
|--------|------|----------|----------|
| `themes_templates` | 主题模板 | 整站样式 | ❌ 缺失 |
| `clientarea_default_themes` | 用户中心默认主题 | 用户中心样式 | ❌ 缺失 |
| `order_page_style` | 订单页面样式 | 购物车样式 | ❌ 缺失 |
| `language` | 默认语言 | 整站语言 | ❌ 缺失 |
| `allow_user_language` | 允许用户切换语言 | 语言切换 | ❌ 缺失 |
| `date_format` | 后台日期格式 | 后台显示 | ❌ 缺失 |
| `client_date_format` | 用户日期格式 | 前台显示 | ❌ 缺失 |
| `default_country` | 默认国家 | 注册/地址 | ❌ 缺失 |
| `header` | 自定义头部 | 网站头部 | ❌ 缺失 |
| `footer` | 自定义底部 | 网站底部 | ❌ 缺失 |
| `login_header` | 登录页头部 | 登录页面 | ❌ 缺失 |
| `login_footer` | 登录页底部 | 登录页面 | ❌ 缺失 |
| `login_header_footer` | 登录页显示头底部 | 登录页面 | ❌ 缺失 |
| `web_widgets` | 网页挂件 | 网站显示 | ❌ 缺失 |
| `cart_product_description` | 购物车产品说明 | 购物车页面 | ❌ 缺失 |
| `per_page_limit` | 每页显示条数 | 分页显示 | ❌ 缺失 |

### 6. 功能开关配置 (缺失 ❌)

| 配置项 | 说明 | 影响范围 | 锚点状态 |
|--------|------|----------|----------|
| `main_tenance_mode` | 维护模式 | 整站访问 | ❌ 缺失 |
| `main_tenance_mode_url` | 维护模式重定向 | 维护模式 | ❌ 缺失 |
| `main_tenance_mode_message` | 维护模式信息 | 维护模式 | ❌ 缺失 |
| `show_cancel` | 显示取消按钮 | 订单操作 | ❌ 缺失 |
| `display_errors` | 显示错误 | 错误处理 | ❌ 缺失 |
| `sql_error_reporting` | SQL错误报告 | 错误处理 | ❌ 缺失 |
| `hooks_debug_mode` | 钩子调试模式 | 开发调试 | ❌ 缺失 |
| `dl_incl_product` | 下载包含产品 | 下载功能 | ❌ 缺失 |
| `nologin_send_ticket` | 未登录发送工单 | 工单功能 | ❌ 缺失 |
| `evaluate_ticket` | 工单评价 | 工单功能 | ❌ 缺失 |
| `ticket_reply_order` | 工单回复排序 | 工单功能 | ❌ 缺失 |

### 7. 代理/推广配置 (缺失 ❌)

| 配置项 | 说明 | 影响范围 | 锚点状态 |
|--------|------|----------|----------|
| `affiliate_enabled` | 代理功能开关 | 代理系统 | ❌ 缺失 |
| `affiliate_percent` | 代理佣金比例 | 代理结算 | ❌ 缺失 |
| `affiliate_bonusde_posit` | 代理奖金 | 代理结算 | ❌ 缺失 |
| `affiliate_payout` | 代理提现 | 代理提现 | ❌ 缺失 |
| `affiliate_cookie` | 代理Cookie | 代理追踪 | ❌ 缺失 |
| `affiliate_withdraw` | 代理提现设置 | 代理提现 | ❌ 缺失 |
| `affiliate_is_authentication` | 代理认证 | 代理认证 | ❌ 缺失 |
| `affiliate_delay_commission` | 延迟佣金 | 代理结算 | ❌ 缺失 |
| `affiliate_is_reorder` | 代理重复订单 | 代理统计 | ❌ 缺失 |
| `affiliate_reorder` | 代理重复订单设置 | 代理统计 | ❌ 缺失 |
| `affiliate_is_renew` | 代理续费 | 代理统计 | ❌ 缺失 |
| `affiliate_renew` | 代理续费设置 | 代理统计 | ❌ 缺失 |
| `aff_report` | 代理报告 | 代理报表 | ❌ 缺失 |

### 8. 充值/财务配置 (缺失 ❌)

| 配置项 | 说明 | 影响范围 | 锚点状态 |
|--------|------|----------|----------|
| `addfunds_enabled` | 充值功能开关 | 充值功能 | ❌ 缺失 |
| `addfunds_minimum` | 单笔最小金额 | 充值限制 | ❌ 缺失 |
| `addfunds_maximum` | 单笔最大金额 | 充值限制 | ❌ 缺失 |
| `addfunds_maximum_balance` | 最高余额 | 充值限制 | ❌ 缺失 |
| `addfunds_require_order` | 充值需要订单 | 充值流程 | ❌ 缺失 |
| `no_auto_apply_credit` | 不自动使用信用额 | 支付逻辑 | ❌ 缺失 |
| `credit_on_downgrade` | 降级时信用额 | 升降级 | ❌ 缺失 |
| `credit_limit` | 信用额开关 | 信用额功能 | ❌ 缺失 |

### 9. 发票配置 (缺失 ❌)

| 配置项 | 说明 | 影响范围 | 锚点状态 |
|--------|------|----------|----------|
| `in_circulation_create` | 循环订单创建发票 | 发票生成 | ❌ 缺失 |
| `in_pdf` | PDF发票 | 发票下载 | ❌ 缺失 |
| `in_save_user_info` | 保存用户信息 | 发票信息 | ❌ 缺失 |
| `in_batch_pay` | 批量支付 | 支付功能 | ❌ 缺失 |
| `in_select_payment` | 选择支付方式 | 支付流程 | ❌ 缺失 |
| `in_unpaid_tick` | 未支付订单 | 订单管理 | ❌ 缺失 |
| `in_continuous_pay_num` | 连续支付数量 | 支付流程 | ❌ 缺失 |
| `in_continuous_pay_num_type` | 连续支付类型 | 支付流程 | ❌ 缺失 |
| `in_overdue_fine` | 逾期罚款 | 账单管理 | ❌ 缺失 |
| `in_overdue_fine_min` | 逾期罚款最小值 | 账单管理 | ❌ 缺失 |
| `invoice_payto` | 付款条文 | 发票显示 | ❌ 缺失 |

### 10. 邮件/短信配置 (部分实现 ⚠️)

| 配置项 | 说明 | 影响范围 | 锚点状态 |
|--------|------|----------|----------|
| `type` | 邮件类型 | 邮件发送 | ✅ 已有 |
| `charset` | 字符集 | 邮件编码 | ✅ 已有 |
| `port` | 端口 | 邮件连接 | ✅ 已有 |
| `host` | 主机 | 邮件连接 | ✅ 已有 |
| `username` | 用户名 | 邮件认证 | ✅ 已有 |
| `password` | 密码 | 邮件认证 | ✅ 已有 |
| `fromname` | 发件人名称 | 邮件显示 | ✅ 已有 |
| `systememail` | 系统邮箱 | 邮件发送 | ✅ 已有 |
| `sendmsgtimes` | 每天短信发送次数 | 短信限制 | ❌ 缺失 |
| `sendmsgphone` | 每天短信发送手机个数 | 短信限制 | ❌ 缺失 |
| `deletelogtime` | 删除日志天数 | 日志清理 | ❌ 缺失 |

### 11. OAuth/第三方登录配置 (部分实现 ⚠️)

| 配置项 | 说明 | 影响范围 | 锚点状态 |
|--------|------|----------|----------|
| `wechat_login_appid` | 微信登录AppID | 微信登录 | ✅ 已有 |
| `wechat_login_secret` | 微信登录Secret | 微信登录 | ✅ 已有 |
| `allow_email_register_code` | 邮箱注册发送验证码 | 注册流程 | ❌ 缺失 |

### 12. 用户ID配置 (缺失 ❌)

| 配置项 | 说明 | 影响范围 | 锚点状态 |
|--------|------|----------|----------|
| `allow_custom_clients_id` | 自定义用户ID | 用户管理 | ❌ 缺失 |
| `custom_clients_id_start` | 自定义ID起始值 | 用户管理 | ❌ 缺失 |

### 13. Cookie配置 (缺失 ❌)

| 配置项 | 说明 | 影响范围 | 锚点状态 |
|--------|------|----------|----------|
| `cookie_clientarea_nmae` | 用户中心Cookie名称 | 登录状态 | ❌ 缺失 |

### 14. 认证/授权配置 (缺失 ❌)

| 配置项 | 说明 | 影响范围 | 锚点状态 |
|--------|------|----------|----------|
| `zjmf_authorize` | 魔方授权 | 系统授权 | ❌ 不需要 |
| `system_license` | 系统许可证 | 系统授权 | ❌ 不需要 |

### 15. 支付宝认证配置 (缺失 ❌)

| 配置项 | 说明 | 影响范围 | 锚点状态 |
|--------|------|----------|----------|
| `app_id` | 支付宝AppID | 支付宝认证 | ❌ 缺失 |
| `private_key` | 支付宝私钥 | 支付宝认证 | ❌ 缺失 |
| `public_key` | 支付宝公钥 | 支付宝认证 | ❌ 缺失 |
| `biz_code` | 支付宝业务码 | 支付宝认证 | ❌ 缺失 |

### 16. 短信服务配置 (缺失 ❌)

| 配置项 | 说明 | 影响范围 | 锚点状态 |
|--------|------|----------|----------|
| `accesskeyid` | 阿里云AccessKeyID | 短信发送 | ❌ 缺失 |
| `accesskeysecret` | 阿里云AccessKeySecret | 短信发送 | ❌ 缺失 |
| `code` | 短信模板码 | 短信发送 | ❌ 缺失 |
| `signature` | 短信签名 | 短信发送 | ❌ 缺失 |
| `submail_appid` | SubMail AppID | 短信发送 | ❌ 缺失 |
| `submail_appkey` | SubMail AppKey | 短信发送 | ❌ 缺失 |

### 17. 后台配置 (缺失 ❌)

| 配置项 | 说明 | 影响范围 | 锚点状态 |
|--------|------|----------|----------|
| `admin_default_theme` | 后台默认主题 | 后台样式 | ❌ 缺失 |
| `admin_application` | 后台应用名称 | 后台路径 | ❌ 缺失 |

---

## 三、配置项影响范围映射

### 登录页面影响的配置项
```
登录页面
├── 验证码相关
│   ├── is_captcha (总开关)
│   ├── allow_login_phone_captcha (手机登录验证码)
│   ├── allow_login_email_captcha (邮箱登录验证码)
│   ├── allow_login_code_captcha (验证码登录验证码)
│   ├── allow_login_id_captcha (ID登录验证码)
│   └── allow_login_admin_captcha (后台登录验证码)
├── 登录方式
│   ├── allow_phone (手机登录)
│   ├── allow_email (邮箱登录)
│   └── allow_wechat (微信登录)
├── 安全限制
│   ├── login_error_max_num (错误次数限制)
│   ├── login_error_switch (错误限制开关)
│   ├── invalid_logins_banlength (封禁时长)
│   └── home_ip_check (IP检测)
└── 界面显示
    ├── login_header (登录页头部)
    ├── login_footer (登录页底部)
    └── login_header_footer (显示头底部)
```

### 注册页面影响的配置项
```
注册页面
├── 验证码相关
│   ├── is_captcha (总开关)
│   ├── allow_register_phone_captcha (手机注册验证码)
│   └── allow_register_email_captcha (邮箱注册验证码)
├── 注册方式
│   ├── allow_register_phone (手机注册)
│   ├── allow_register_email (邮箱注册)
│   └── allow_register_wechat (微信注册)
├── 用户资料
│   ├── clients_profoptional (可选字段)
│   └── login_register_custom_require (自定义字段)
├── 安全要求
│   ├── required_pwstrength (密码强度)
│   └── allow_email_register_code (邮箱验证码)
└── 法律条款
    ├── server_clause_url (服务条款)
    └── privacy_clause_url (隐私条款)
```

### 用户中心影响的配置项
```
用户中心
├── 界面主题
│   ├── clientarea_default_themes (默认主题)
│   └── themes_templates (可用主题)
├── 显示设置
│   ├── client_date_format (日期格式)
│   ├── allow_user_language (语言切换)
│   ├── language (默认语言)
│   └── per_page_limit (每页条数)
├── 功能开关
│   ├── credit_limit (信用额开关)
│   ├── addfunds_enabled (充值开关)
│   └── affiliate_enabled (代理开关)
└── 安全设置
    ├── allow_phone_bind_captcha (手机绑定验证码)
    ├── allow_email_bind_captcha (邮箱绑定验证码)
    ├── allow_cancel_sms_captcha (取消短信验证码)
    └── allow_cancel_email_captcha (取消邮箱验证码)
```

### 购物车/订单影响的配置项
```
购物车/订单
├── 样式配置
│   ├── order_page_style (订单页面样式)
│   └── cart_product_description (产品说明)
├── 支付配置
│   ├── in_select_payment (选择支付方式)
│   ├── in_batch_pay (批量支付)
│   ├── in_continuous_pay_num (连续支付)
│   └── in_continuous_pay_num_type (支付类型)
├── 发票配置
│   ├── in_circulation_create (循环发票)
│   ├── in_pdf (PDF发票)
│   ├── in_save_user_info (保存用户信息)
│   └── invoice_payto (付款条文)
└── 逾期配置
    ├── in_overdue_fine (逾期罚款)
    └── in_overdue_fine_min (罚款最小值)
```

### 全站影响的配置项
```
全站影响
├── 公司信息
│   ├── company_name (公司名称)
│   ├── company_email (公司邮箱)
│   ├── logo_url (Logo)
│   ├── logo_url_home (前台Logo)
│   ├── logo_url_bill (账单Logo)
│   ├── www_logo (官网Logo)
│   ├── main_phone (电话)
│   ├── main_address (地址)
│   └── record_no (备案号)
├── SEO设置
│   ├── seo_keywords (关键词)
│   └── seo_desc (描述)
├── 维护模式
│   ├── main_tenance_mode (开关)
│   ├── main_tenance_mode_url (重定向)
│   └── main_tenance_mode_message (提示信息)
├── 自定义内容
│   ├── header (头部)
│   ├── footer (底部)
│   └── web_widgets (挂件)
└── 系统设置
    ├── system_url (系统链接)
    ├── domain (域名)
    └── display_errors (错误显示)
```

---

## 四、锚点财务需要补充的配置项

### 优先级 P0 (核心功能必须)
1. 登录/注册方式开关 (allow_phone, allow_email, allow_register_*)
2. 安全配置 (密码强度、登录限制、IP检测)
3. 维护模式
4. 公司Logo多场景 (logo_url_home, logo_url_bill)

### 优先级 P1 (重要功能)
1. 界面/模板配置
2. 代理/推广配置
3. 充值/财务配置
4. 发票配置

### 优先级 P2 (完善功能)
1. 自定义内容 (header, footer, web_widgets)
2. SEO配置
3. 日志配置
4. 短信服务配置
