#!/usr/bin/env python3
"""
运行数据库初始化脚本
"""

import paramiko

SERVER_HOST = "45.207.210.235"
SERVER_USER = "root"
SERVER_PASS = "iswlBSLY8118"

client = paramiko.SSHClient()
client.set_missing_host_key_policy(paramiko.AutoAddPolicy())
client.connect(SERVER_HOST, username=SERVER_USER, password=SERVER_PASS)

# 先检查MySQL状态
print("=== 检查MySQL ===")
stdin, stdout, stderr = client.exec_command("systemctl status mysql | head -5")
print(stdout.read().decode())

# 尝试创建数据库
print("\n=== 创建数据库 ===")
stdin, stdout, stderr = client.exec_command(
    "mysql -u root -e 'CREATE DATABASE IF NOT EXISTS anchor_finance CHARACTER SET utf8mb4;' 2>&1 || echo '需要密码'"
)
print(stdout.read().decode())

# 尝试用不同密码
print("\n=== 测试MySQL连接 ===")
for password in ['', 'root', 'mysql', 'iswlBSLY8118']:
    cmd = f"mysql -u root -p'{password}' -e 'SELECT 1;' 2>&1 && echo '密码正确: {password}' || echo '密码错误: {password}'"
    stdin, stdout, stderr = client.exec_command(cmd)
    output = stdout.read().decode()
    print(output)

client.close()
