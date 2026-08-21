#!/usr/bin/env python3
"""
运行数据库初始化脚本
"""

import paramiko
import sys

# 设置输出编码
sys.stdout.reconfigure(encoding='utf-8')

SERVER_HOST = "45.207.210.235"
SERVER_USER = "root"
SERVER_PASS = "iswlBSLY8118"
MYSQL_PASS = "*RhbY#m0IbS8SPaAteOI"

client = paramiko.SSHClient()
client.set_missing_host_key_policy(paramiko.AutoAddPolicy())
client.connect(SERVER_HOST, username=SERVER_USER, password=SERVER_PASS)

# 确保数据库存在
print("=== 确保数据库存在 ===")
stdin, stdout, stderr = client.exec_command(
    f"mysql -u root -p'{MYSQL_PASS}' -e 'CREATE DATABASE IF NOT EXISTS anchor_finance CHARACTER SET utf8mb4;'"
)
print(stdout.read().decode('utf-8', errors='ignore'))

# 运行seed脚本
print("\n=== 运行seed脚本 ===")
cmd = f"""
export DB_HOST=127.0.0.1
export DB_PORT=3306
export DB_USER=root
export DB_PASSWORD='{MYSQL_PASS}'
export DB_NAME=anchor_finance
cd /opt/anchor-finance/backend
/usr/local/go/bin/go run cmd/seed/main.go
"""

stdin, stdout, stderr = client.exec_command(cmd)
print(stdout.read().decode('utf-8', errors='ignore'))
err = stderr.read().decode('utf-8', errors='ignore')
if err:
    print("STDERR:", err)

client.close()
