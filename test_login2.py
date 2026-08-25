import paramiko
import json
import time
import sys

sys.stdout.reconfigure(encoding='utf-8')

client = paramiko.SSHClient()
client.set_missing_host_key_policy(paramiko.AutoAddPolicy())
client.connect('45.207.210.235', username='root', password='iswlBSLY8118')

def run(cmd):
    stdin, stdout, stderr = client.exec_command(cmd, timeout=15)
    return stdout.read().decode('utf-8', errors='replace').strip()

# 杀残留进程 + 重启
print("=== 彻底重启 ===")
print(run("fuser -k 8080/tcp 2>&1; sleep 3"))
print(run("systemctl restart anchor-finance"))
time.sleep(5)
print("状态:", run("systemctl is-active anchor-finance"))
print("端口:", run("ss -tlnp | grep 8080 || echo 无"))
print("二进制:", run("ls -la /opt/anchor-finance/server/anchor-finance"))

# 登录（用最简单的登录试试）
print("\n=== 登录尝试 ===")
for pw in ["admin123", "123456", "admin"]:
    login = run(f'curl -s -X POST http://127.0.0.1:8080/api/admin/login -H "Content-Type: application/json" -d \'{{"username":"admin","password":"{pw}"}}\'')
    print(f"  pw={pw}: {login[:100]}")

client.close()
