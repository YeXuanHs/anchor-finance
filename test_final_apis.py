import paramiko
import json

client = paramiko.SSHClient()
client.set_missing_host_key_policy(paramiko.AutoAddPolicy())
client.connect('45.207.210.235', username='root', password='iswlBSLY8118')

# 更新systemd配置
service_content = """[Unit]
Description=AnchorFinance Server
After=network.target mysql.service

[Service]
Type=simple
User=root
WorkingDirectory=/opt/anchor-finance/server
ExecStart=/opt/anchor-finance/server/anchor-finance
Restart=always
RestartSec=5
Environment=DB_HOST=127.0.0.1
Environment=DB_PORT=3306
Environment=DB_USER=root
Environment=DB_PASSWORD='*RhbY#m0IbS8SPaAteOI'
Environment=DB_NAME=anchor_finance
Environment=SERVER_PORT=8080
Environment=SERVER_MODE=release
Environment=JWT_SECRET=anchor-finance-secret-key-2024

[Install]
WantedBy=multi-user.target
"""

sftp = client.open_sftp()
with sftp.file('/etc/systemd/system/anchor-finance.service', 'w') as f:
    f.write(service_content)
sftp.close()

stdin, stdout, stderr = client.exec_command("fuser -k 8080/tcp 2>/dev/null; systemctl daemon-reload; systemctl restart anchor-finance")
stdout.read()

import time
time.sleep(3)

# 管理员登录
print("=== 管理员登录 ===")
cmd = """curl -s -X POST http://localhost:8080/api/admin/login -H "Content-Type: application/json" -d '{"username":"admin","password":"admin123"}'"""
stdin, stdout, stderr = client.exec_command(cmd)
admin_login = json.loads(stdout.read().decode('utf-8', errors='ignore'))
admin_token = admin_login.get('data', {}).get('token', '')
print(f"OK")

# 测试所有管理后台API
print("\n=== 测试管理后台API ===")
admin_apis = [
    ("GET", "/api/admin/auth/info", "管理员信息"),
    ("GET", "/api/admin/dashboard/stats", "仪表盘统计"),
    ("GET", "/api/admin/users", "用户列表"),
    ("GET", "/api/admin/orders", "订单列表"),
    ("GET", "/api/admin/services", "服务列表"),
    ("GET", "/api/admin/invoices", "账单列表"),
    ("GET", "/api/admin/tickets", "工单列表"),
    ("GET", "/api/admin/tickets/summary", "工单统计"),
    ("GET", "/api/admin/products", "产品列表"),
    ("GET", "/api/admin/products/summary", "产品统计"),
    ("GET", "/api/admin/product-groups", "产品分组"),
    ("GET", "/api/admin/plugins", "插件列表"),
    ("GET", "/api/admin/settings", "系统设置"),
    ("GET", "/api/admin/menus", "菜单列表"),
    ("GET", "/api/admin/admins", "管理员列表"),
    ("GET", "/api/admin/roles", "角色列表"),
    ("GET", "/api/admin/permissions", "权限列表"),
    ("GET", "/api/admin/currencies", "货币列表"),
    ("GET", "/api/admin/suppliers", "供应商列表"),
    ("GET", "/api/admin/promo-codes", "优惠码列表"),
    ("GET", "/api/admin/verifications", "认证列表"),
    ("GET", "/api/admin/verifications/summary", "认证统计"),
    ("GET", "/api/admin/news", "新闻列表"),
    ("GET", "/api/admin/news-categories", "新闻分类"),
    ("GET", "/api/admin/knowledge/categories", "知识库分类"),
    ("GET", "/api/admin/knowledge/articles", "知识库文章"),
    ("GET", "/api/admin/downloads", "下载列表"),
    ("GET", "/api/admin/downloads/categories", "下载分类"),
    ("GET", "/api/admin/system-logs", "系统日志"),
    ("GET", "/api/admin/operation-logs", "操作日志"),
    ("GET", "/api/admin/login-logs", "登录日志"),
    ("GET", "/api/admin/database/status", "数据库状态"),
    ("GET", "/api/admin/referral/overview", "推介概览"),
    ("GET", "/api/admin/referral/rewards", "推介奖励"),
    ("GET", "/api/admin/referral-withdrawals", "提现列表"),
]

admin_success = 0
for method, path, desc in admin_apis:
    cmd = f'curl -s "http://localhost:8080{path}" -H "Authorization: Bearer {admin_token}"'
    stdin, stdout, stderr = client.exec_command(cmd)
    result = stdout.read().decode('utf-8', errors='ignore')
    try:
        data = json.loads(result)
        if data.get('code') == 0:
            print(f"  {desc}: OK")
            admin_success += 1
        else:
            print(f"  {desc}: FAIL - {data.get('message')}")
    except:
        print(f"  {desc}: FAIL")

# 测试用户前台API
print("\n=== 测试用户前台API ===")

# 用户注册
cmd = """curl -s -X POST http://localhost:8080/api/client/register -H "Content-Type: application/json" -d '{"username":"testuser2","email":"test2@client.com","password":"123456"}'"""
stdin, stdout, stderr = client.exec_command(cmd)
stdout.read()

# 用户登录
cmd = """curl -s -X POST http://localhost:8080/api/client/login -H "Content-Type: application/json" -d '{"username":"testuser2","password":"123456"}'"""
stdin, stdout, stderr = client.exec_command(cmd)
client_login = json.loads(stdout.read().decode('utf-8', errors='ignore'))
client_token = client_login.get('data', {}).get('token', '')

client_apis = [
    ("GET", "/api/client/auth/info", "用户信息"),
    ("GET", "/api/client/services", "服务列表"),
    ("GET", "/api/client/orders", "订单列表"),
    ("GET", "/api/client/tickets", "工单列表"),
    ("GET", "/api/client/invoices", "账单列表"),
    ("GET", "/api/client/balance-logs", "余额日志"),
    ("GET", "/api/client/recharge/gateways", "充值网关"),
    ("GET", "/api/client/coupons", "优惠券"),
    ("GET", "/api/client/notifications", "通知"),
    ("GET", "/api/client/notices", "公告"),
    ("GET", "/api/client/help-articles", "帮助文章"),
]

client_success = 0
for method, path, desc in client_apis:
    cmd = f'curl -s "http://localhost:8080{path}" -H "Authorization: Bearer {client_token}"'
    stdin, stdout, stderr = client.exec_command(cmd)
    result = stdout.read().decode('utf-8', errors='ignore')
    try:
        data = json.loads(result)
        if data.get('code') == 0:
            print(f"  {desc}: OK")
            client_success += 1
        else:
            print(f"  {desc}: FAIL - {data.get('message')}")
    except:
        print(f"  {desc}: FAIL")

print(f"\n=== 测试结果 ===")
print(f"Admin API: {admin_success}/{len(admin_apis)} 通过")
print(f"Client API: {client_success}/{len(client_apis)} 通过")
print(f"总计: {admin_success + client_success}/{len(admin_apis) + len(client_apis)} 通过")

client.close()
