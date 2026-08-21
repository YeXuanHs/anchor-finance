import paramiko
import sys
import json

sys.stdout.reconfigure(encoding='utf-8')

client = paramiko.SSHClient()
client.set_missing_host_key_policy(paramiko.AutoAddPolicy())
client.connect('45.207.210.235', username='root', password='iswlBSLY8118')

# 检查服务状态
print("=== 检查服务状态 ===")
stdin, stdout, stderr = client.exec_command("systemctl status anchor-finance | head -5")
print(stdout.read().decode('utf-8', errors='ignore'))

# 检查日志
print("\n=== 检查日志 ===")
stdin, stdout, stderr = client.exec_command("journalctl -u anchor-finance -n 5 --no-pager")
print(stdout.read().decode('utf-8', errors='ignore'))

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

print("\n=== 更新systemd配置 ===")
sftp = client.open_sftp()
with sftp.file('/etc/systemd/system/anchor-finance.service', 'w') as f:
    f.write(service_content)
sftp.close()

# 重启服务
print("=== 重启服务 ===")
stdin, stdout, stderr = client.exec_command("fuser -k 8080/tcp 2>/dev/null; systemctl daemon-reload; systemctl restart anchor-finance")
print(stdout.read().decode('utf-8', errors='ignore'))

import time
time.sleep(3)

# 测试登录
print("\n=== 测试登录 ===")
cmd = """curl -s -X POST http://localhost:8080/api/admin/login -H "Content-Type: application/json" -d '{"username":"admin","password":"admin123"}'"""
stdin, stdout, stderr = client.exec_command(cmd)
login_result = stdout.read().decode('utf-8', errors='ignore')
print(login_result)

try:
    token_data = json.loads(login_result)
    token = token_data.get('data', {}).get('token', '')
    
    if token:
        # 测试获取服务列表
        print("\n=== 测试获取服务列表 ===")
        list_cmd = f'''curl -s "http://localhost:8080/api/admin/services?page=1&page_size=10" \
            -H "Authorization: Bearer {token}"'''
        stdin, stdout, stderr = client.exec_command(list_cmd)
        print(stdout.read().decode('utf-8', errors='ignore'))
except Exception as e:
    print(f"错误: {e}")

client.close()
