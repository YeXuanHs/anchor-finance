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
cmd = """curl -s -X POST http://localhost:8080/api/admin/login -H "Content-Type: application/json" -d '{"username":"admin","password":"admin123"}'"""
stdin, stdout, stderr = client.exec_command(cmd)
login_result = json.loads(stdout.read().decode('utf-8', errors='ignore'))
token = login_result.get('data', {}).get('token', '')

# 测试会员等级API
print("=== 测试会员等级API ===")

# 获取等级列表
cmd = f'curl -s "http://localhost:8080/api/admin/member-levels" -H "Authorization: Bearer {token}"'
stdin, stdout, stderr = client.exec_command(cmd)
result = stdout.read().decode('utf-8', errors='ignore')
data = json.loads(result)
if data.get('code') == 0:
    print("  等级列表: OK")
else:
    print(f"  等级列表: FAIL - {data.get('message')}")

# 创建等级
cmd = f'''curl -s -X POST "http://localhost:8080/api/admin/member-levels" -H "Authorization: Bearer {token}" -H "Content-Type: application/json" -d '{{"name":"VIP1","description":"VIP1会员","discount":95,"min_points":100}}' '''
stdin, stdout, stderr = client.exec_command(cmd)
result = stdout.read().decode('utf-8', errors='ignore')
data = json.loads(result)
if data.get('code') == 0:
    print("  创建等级: OK")
else:
    print(f"  创建等级: FAIL - {data.get('message')}")

client.close()
