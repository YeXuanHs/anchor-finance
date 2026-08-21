import paramiko
import sys
import json

sys.stdout.reconfigure(encoding='utf-8')

client = paramiko.SSHClient()
client.set_missing_host_key_policy(paramiko.AutoAddPolicy())
client.connect('45.207.210.235', username='root', password='iswlBSLY8118')

# 更新systemd配置并重启
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

# 登录获取token
print("=== 登录 ===")
cmd = """curl -s -X POST http://localhost:8080/api/admin/login -H "Content-Type: application/json" -d '{"username":"admin","password":"admin123"}'"""
stdin, stdout, stderr = client.exec_command(cmd)
login_result = json.loads(stdout.read().decode('utf-8', errors='ignore'))
token = login_result.get('data', {}).get('token', '')
print(f"✅ 登录成功")

# 测试所有API
apis = [
    ("GET", "/api/admin/auth/info", "获取管理员信息"),
    ("GET", "/api/admin/dashboard/stats", "仪表盘统计"),
    ("GET", "/api/admin/users?page=1&page_size=5", "用户列表"),
    ("GET", "/api/admin/orders?page=1&page_size=5", "订单列表"),
    ("GET", "/api/admin/services?page=1&page_size=5", "服务列表"),
    ("GET", "/api/admin/invoices?page=1&page_size=5", "账单列表"),
    ("GET", "/api/admin/tickets?page=1&page_size=5", "工单列表"),
    ("GET", "/api/admin/tickets/summary", "工单统计"),
    ("GET", "/api/admin/products?page=1&page_size=5", "产品列表"),
    ("GET", "/api/admin/products/summary", "产品统计"),
    ("GET", "/api/admin/product-groups", "产品分组"),
    ("GET", "/api/admin/suppliers?page=1&page_size=5", "供应商列表"),
    ("GET", "/api/admin/plugins", "插件列表"),
    ("GET", "/api/admin/settings", "系统设置"),
    ("GET", "/api/admin/menus", "菜单列表"),
    ("GET", "/api/admin/admins?page=1&page_size=5", "管理员列表"),
    ("GET", "/api/admin/roles", "角色列表"),
    ("GET", "/api/admin/system-logs?page=1&page_size=5", "系统日志"),
    ("GET", "/api/admin/news?page=1&page_size=5", "新闻列表"),
    ("GET", "/api/admin/currencies", "货币列表"),
    ("GET", "/api/admin/promo-codes?page=1&page_size=5", "优惠码列表"),
    ("GET", "/api/admin/verifications?page=1&page_size=5", "实名认证列表"),
]

print("\n=== 测试所有API ===")
success_count = 0
fail_count = 0

for method, path, desc in apis:
    cmd = f'''curl -s "{path}" -H "Authorization: Bearer {token}"'''
    stdin, stdout, stderr = client.exec_command(cmd)
    result = stdout.read().decode('utf-8', errors='ignore')
    
    try:
        data = json.loads(result)
        if data.get('code') == 0:
            print(f"✅ {desc}")
            success_count += 1
        else:
            print(f"❌ {desc}: {data.get('message', '未知错误')}")
            fail_count += 1
    except:
        print(f"❌ {desc}: 解析失败")
        fail_count += 1

print(f"\n=== 测试结果 ===")
print(f"成功: {success_count}")
print(f"失败: {fail_count}")
print(f"总计: {success_count + fail_count}")

client.close()
