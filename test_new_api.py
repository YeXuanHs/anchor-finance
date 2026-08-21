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

# 登录
cmd = """curl -s -X POST http://localhost:8080/api/admin/login -H "Content-Type: application/json" -d '{"username":"admin","password":"admin123"}'"""
stdin, stdout, stderr = client.exec_command(cmd)
login_result = json.loads(stdout.read().decode('utf-8', errors='ignore'))
token = login_result.get('data', {}).get('token', '')

# 测试工单回复API
print("=== 测试工单回复API ===")
cmd = f'curl -s "http://localhost:8080/api/admin/tickets/1/replies" -H "Authorization: Bearer {token}"'
stdin, stdout, stderr = client.exec_command(cmd)
result = stdout.read().decode('utf-8', errors='ignore')
print(result)

# 验证返回格式
data = json.loads(result)
if data.get('code') == 0:
    print("\nOK - Ticket replies API passed!")
else:
    print(f"\nFAIL - {data.get('message')}")

client.close()
