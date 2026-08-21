import paramiko
import sys

sys.stdout.reconfigure(encoding='utf-8')

client = paramiko.SSHClient()
client.set_missing_host_key_policy(paramiko.AutoAddPolicy())
client.connect('45.207.210.235', username='root', password='iswlBSLY8118')

# 杀掉占用8080端口的进程
print("=== 杀掉占用8080端口的进程 ===")
stdin, stdout, stderr = client.exec_command("fuser -k 8080/tcp 2>/dev/null; sleep 1; systemctl restart anchor-finance")
print(stdout.read().decode('utf-8', errors='ignore'))

import time
time.sleep(3)

# 检查状态
print("\n=== 服务状态 ===")
stdin, stdout, stderr = client.exec_command("systemctl status anchor-finance | head -8")
print(stdout.read().decode('utf-8', errors='ignore'))

# 测试登录
print("\n=== 测试登录 ===")
cmd = """curl -s -X POST http://localhost:8080/api/admin/login -H "Content-Type: application/json" -d '{"username":"admin","password":"admin123"}'"""
stdin, stdout, stderr = client.exec_command(cmd)
print(stdout.read().decode('utf-8', errors='ignore'))

client.close()
