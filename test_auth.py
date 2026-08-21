import paramiko
import sys
import json

sys.stdout.reconfigure(encoding='utf-8')

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

# 重启服务
stdin, stdout, stderr = client.exec_command("fuser -k 8080/tcp 2>/dev/null; systemctl daemon-reload; systemctl restart anchor-finance")
stdout.read()

import time
time.sleep(3)

# 登录
print("=== 登录 ===")
cmd = """curl -s -X POST http://localhost:8080/api/admin/login -H "Content-Type: application/json" -d '{"username":"admin","password":"admin123"}'"""
stdin, stdout, stderr = client.exec_command(cmd)
login_result = json.loads(stdout.read().decode('utf-8', errors='ignore'))
token = login_result.get('data', {}).get('token', '')
print(f"Token: {token[:50]}...")

if token:
    # 测试获取管理员信息
    print("\n=== 测试获取管理员信息 ===")
    cmd = f'''curl -s "http://localhost:8080/api/admin/auth/info" -H "Authorization: Bearer {token}"'''
    stdin, stdout, stderr = client.exec_command(cmd)
    print(stdout.read().decode('utf-8', errors='ignore'))
    
    # 测试更新个人资料
    print("\n=== 测试更新个人资料 ===")
    cmd = f'''curl -s -X PUT http://localhost:8080/api/admin/auth/profile \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer {token}" \
        -d '{{"real_name":"超级管理员","email":"admin@anchor.com"}}' '''
    stdin, stdout, stderr = client.exec_command(cmd)
    print(stdout.read().decode('utf-8', errors='ignore'))
    
    # 测试修改密码（用旧密码）
    print("\n=== 测试修改密码 ===")
    cmd = f'''curl -s -X PUT http://localhost:8080/api/admin/auth/password \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer {token}" \
        -d '{{"old_password":"admin123","new_password":"admin123456"}}' '''
    stdin, stdout, stderr = client.exec_command(cmd)
    print(stdout.read().decode('utf-8', errors='ignore'))
    
    # 测试用新密码登录
    print("\n=== 测试用新密码登录 ===")
    cmd = """curl -s -X POST http://localhost:8080/api/admin/login -H "Content-Type: application/json" -d '{"username":"admin","password":"admin123456"}'"""
    stdin, stdout, stderr = client.exec_command(cmd)
    print(stdout.read().decode('utf-8', errors='ignore'))
    
    # 改回原密码
    print("\n=== 改回原密码 ===")
    cmd = """curl -s -X POST http://localhost:8080/api/admin/login -H "Content-Type: application/json" -d '{"username":"admin","password":"admin123456"}'"""
    stdin, stdout, stderr = client.exec_command(cmd)
    new_login = json.loads(stdout.read().decode('utf-8', errors='ignore'))
    new_token = new_login.get('data', {}).get('token', '')
    
    if new_token:
        cmd = f'''curl -s -X PUT http://localhost:8080/api/admin/auth/password \
            -H "Content-Type: application/json" \
            -H "Authorization: Bearer {new_token}" \
            -d '{{"old_password":"admin123456","new_password":"admin123"}}' '''
        stdin, stdout, stderr = client.exec_command(cmd)
        print(stdout.read().decode('utf-8', errors='ignore'))

client.close()
