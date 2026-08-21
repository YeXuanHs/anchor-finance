import paramiko
import time

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

# 杀掉占用端口的进程
print("=== 杀掉占用端口的进程 ===")
stdin, stdout, stderr = client.exec_command("fuser -k 8080/tcp 2>/dev/null")
print(stdout.read().decode('utf-8', errors='ignore'))

# 重启服务
print("=== 重启服务 ===")
stdin, stdout, stderr = client.exec_command("systemctl daemon-reload && systemctl restart anchor-finance")
print(stdout.read().decode('utf-8', errors='ignore'))

time.sleep(3)

# 检查状态
print("\n=== 服务状态 ===")
stdin, stdout, stderr = client.exec_command("systemctl status anchor-finance | head -8")
print(stdout.read().decode('utf-8', errors='ignore'))

# 测试健康检查
print("\n=== 健康检查 ===")
stdin, stdout, stderr = client.exec_command("curl -s http://localhost:8080/health")
print(stdout.read().decode('utf-8', errors='ignore'))

# 测试登录
print("\n=== 登录测试 ===")
stdin, stdout, stderr = client.exec_command("""curl -s -X POST http://localhost:8080/api/admin/login -H "Content-Type: application/json" -d '{"username":"admin","password":"admin123"}'""")
print(stdout.read().decode('utf-8', errors='ignore'))

client.close()
