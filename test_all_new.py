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
print("=== 登录 ===")
cmd = """curl -s -X POST http://localhost:8080/api/admin/login -H "Content-Type: application/json" -d '{"username":"admin","password":"admin123"}'"""
stdin, stdout, stderr = client.exec_command(cmd)
login_result_raw = stdout.read().decode('utf-8', errors='ignore')
print(f"登录响应: {login_result_raw[:200]}")

try:
    login_result = json.loads(login_result_raw)
    token = login_result.get('data', {}).get('token', '')
    print(f"Token获取成功")
except:
    print("登录失败，无法获取token")
    client.close()
    exit(1)

# 测试所有新API
apis = [
    ("GET", "/api/admin/credit-limits", "信用额列表"),
    ("GET", "/api/admin/credit-limits/config", "信用额配置"),
    ("GET", "/api/admin/finance/recharges", "充值记录"),
    ("GET", "/api/admin/finance/recharges/summary", "充值统计"),
    ("GET", "/api/admin/media-files", "媒体文件列表"),
    ("GET", "/api/admin/log-cleanups/overview", "日志清理概览"),
    ("GET", "/api/admin/settings/email", "邮件配置"),
    ("GET", "/api/admin/settings/sms", "短信配置"),
    ("GET", "/api/admin/settings/register-login", "注册登录配置"),
    ("GET", "/api/admin/settings/captcha", "验证码配置"),
    ("GET", "/api/admin/send-message/search-params", "发送消息参数"),
    ("GET", "/api/admin/send-message/send-methods", "发送方式"),
    ("GET", "/api/admin/ticket-prereplies", "工单预回复列表"),
    ("GET", "/api/admin/ticket-prereply-categories", "预回复分类"),
    ("GET", "/api/admin/suppliers/provider-types", "供应商类型"),
    ("GET", "/api/admin/tickets/admin-users", "工单管理员列表"),
]

print("\n=== 测试所有新API ===")
success_count = 0
fail_count = 0

for method, path, desc in apis:
    cmd = f'curl -s "http://localhost:8080{path}" -H "Authorization: Bearer {token}"'
    stdin, stdout, stderr = client.exec_command(cmd)
    result = stdout.read().decode('utf-8', errors='ignore')
    try:
        data = json.loads(result)
        if data.get('code') == 0:
            print(f"  {desc}: OK")
            success_count += 1
        else:
            print(f"  {desc}: FAIL - {data.get('message')}")
            fail_count += 1
    except:
        print(f"  {desc}: FAIL - 解析错误")
        fail_count += 1

print(f"\n=== 测试结果 ===")
print(f"成功: {success_count}")
print(f"失败: {fail_count}")
print(f"总计: {success_count + fail_count}")

client.close()
