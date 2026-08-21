import paramiko
import sys

sys.stdout.reconfigure(encoding='utf-8')

client = paramiko.SSHClient()
client.set_missing_host_key_policy(paramiko.AutoAddPolicy())
client.connect('45.207.210.235', username='root', password='iswlBSLY8118')

# 更新systemd配置（密码加引号）
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

print("=== 更新systemd配置 ===")
sftp = client.open_sftp()
with sftp.file('/etc/systemd/system/anchor-finance.service', 'w') as f:
    f.write(service_content)
sftp.close()

# 杀掉占用端口的进程并重启
print("=== 重启服务 ===")
stdin, stdout, stderr = client.exec_command("fuser -k 8080/tcp 2>/dev/null; systemctl daemon-reload; systemctl restart anchor-finance")
print(stdout.read().decode('utf-8', errors='ignore'))

import time
time.sleep(3)

# 测试健康检查
print("\n=== 测试健康检查 ===")
stdin, stdout, stderr = client.exec_command("curl -s http://localhost:8080/health")
print(stdout.read().decode('utf-8', errors='ignore'))

# 测试登录
print("\n=== 测试登录 ===")
cmd = """curl -s -X POST http://localhost:8080/api/admin/login -H "Content-Type: application/json" -d '{"username":"admin","password":"admin123"}'"""
stdin, stdout, stderr = client.exec_command(cmd)
login_result = stdout.read().decode('utf-8', errors='ignore')
print(login_result)

# 提取token
import json
try:
    token_data = json.loads(login_result)
    token = token_data.get('data', {}).get('token', '')
    
    if token:
        # 测试创建订单
        print("\n=== 测试创建订单 ===")
        create_cmd = f'''curl -s -X POST http://localhost:8080/api/admin/orders \
            -H "Content-Type: application/json" \
            -H "Authorization: Bearer {token}" \
            -d '{{"user_id":1,"product_id":1,"product_name":"云服务器 4核8G","amount":199.00}}' '''
        stdin, stdout, stderr = client.exec_command(create_cmd)
        print(stdout.read().decode('utf-8', errors='ignore'))
        
        # 测试获取订单列表
        print("\n=== 测试获取订单列表 ===")
        list_cmd = f'''curl -s "http://localhost:8080/api/admin/orders?page=1&page_size=10" \
            -H "Authorization: Bearer {token}"'''
        stdin, stdout, stderr = client.exec_command(list_cmd)
        print(stdout.read().decode('utf-8', errors='ignore'))
except Exception as e:
    print(f"解析失败: {e}")

client.close()
