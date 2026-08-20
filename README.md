# 锚点财务 (AnchorFinance)

新一代财务管理系统，兼容智简魔方(zjmf)，支持上游对接、产品管理、工单系统、AI 客服等功能。

## 环境要求

| 项目 | 要求 |
|------|------|
| 操作系统 | Linux (推荐 Ubuntu 20.04+ / CentOS 7+) |
| Go | 1.21+ |
| Node.js | 18+ |
| MySQL | 5.7+ / 8.0+ |
| Redis | 可选（可在后台开启） |

## 安装

### 一键安装

```bash
git clone https://github.com/YeXuanHs/anchor-finance.git
cd anchor-finance
bash install.sh
```

安装脚本会自动完成：
1. 检测编译环境（Go、Node.js、MySQL）
2. 询问数据库、站点名称、管理员账号等配置
3. 编译后端和前端
4. 创建数据库并导入初始数据
5. 将程序安装到 `/opt/anchorfinance`
6. 可选安装 systemd 服务

安装完成后会显示管理后台地址和登录信息。

### 手动编译

如果只需要编译不安装：

```bash
bash build.sh
```

编译产物会输出到 `/opt/anchorfinance`。

## 目录结构

安装完成后，程序位于 `/opt/anchorfinance`：

```
/opt/anchorfinance/
├── anchorfinance        # 主程序
├── .env                 # 数据库连接配置
├── init.sql             # 数据库初始化脚本
├── frontend/            # 用户前台页面
├── admin/               # 管理后台页面
└── README.md
```

## 使用

### 启动

```bash
cd /opt/anchorfinance
./anchorfinance
```

如果安装了 systemd 服务：

```bash
systemctl start anchorfinance
```

### 访问

- **用户前台**: `http://你的域名:端口`
- **管理后台**: `http://你的域名:端口/后台路径`

后台路径在安装时设置，默认随机生成。

### 停止

```bash
systemctl stop anchorfinance
# 或直接 Ctrl+C
```

### 查看日志

```bash
journalctl -u anchorfinance -f
```

## 配置

### 数据库连接

数据库信息存储在 `/opt/anchorfinance/.env` 文件中：

```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASS=你的密码
DB_NAME=anchorfinance
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
```

修改后需重启服务。

### 其他配置

站点名称、Logo、支付方式、邮件短信等配置均在 **管理后台 → 系统设置** 中修改，无需编辑文件。

### Redis（可选）

Redis 默认关闭。在管理后台 → 系统设置 → 高级设置 中开启并配置连接信息。

## 功能特性

- **产品管理**: 云服务器、独立服务器、虚拟主机等多种产品类型
- **订单系统**: 完整的下单、支付、续费流程
- **工单系统**: 多部门工单、AI 自动回复、知识库
- **财务管理**: 账单、发票、优惠码、信用额度
- **上游对接**: 支持智简魔方、WHMCS、V10 等上游系统
- **产品对接**: 从上游拉取产品列表，自动同步库存和价格
- **OAuth 登录**: 支持 24 个第三方登录平台
- **多语言**: 简体中文、英文、繁体中文
- **安全防护**: 登录限制、验证码、XSS 防护、审计日志
- **AI 客服**: 工单自动回复、购物助手、工具调用

## 从智简魔方迁移

如果从智简魔方(zjmf)迁移到锚点财务：

```bash
# 1. 备份 zjmf 数据库
mysqldump -u root -p zjmf_db > zjmf_backup.sql

# 2. 安装锚点财务
bash install.sh

# 3. 执行迁移脚本
mysql -u root -p anchorfinance_db < /opt/anchorfinance/migrate_from_zjmf.sql
```

迁移脚本会自动转换表结构和数据，不会删除原始数据。

## 常见问题

**Q: 忘记管理员密码怎么办？**

```bash
mysql -u root -p anchorfinance -e "UPDATE users SET password='新密码的bcrypt哈希' WHERE username='admin';"
```

**Q: 如何修改后台路径？**

在管理后台 → 系统设置 → 安全设置 中修改，修改后需重启服务。

**Q: 如何更新版本？**

```bash
cd anchor-finance
git pull
bash build.sh
systemctl restart anchorfinance
```

**Q: 如何备份数据库？**

```bash
mysqldump -u root -p anchorfinance > backup_$(date +%Y%m%d).sql
```

## 开发者

- **系统**: 锚点财务
- **团队**: 锚点云计算
- **开发者**: @Sora
- **QQ**: 2338795574
- **Telegram**: @soranb666
- **交流群**: [1009624286](https://qm.qq.com/q/m3i0A7bwga)

> 本系统由锚点云计算团队成员 @Sora 一人开发

## 许可证

MIT License
