import paramiko

client = paramiko.SSHClient()
client.set_missing_host_key_policy(paramiko.AutoAddPolicy())
client.connect('45.207.210.235', username='root', password='iswlBSLY8118')

# 检查服务状态
stdin, stdout, stderr = client.exec_command("systemctl status anchor-finance | head -10")
print("=== 服务状态 ===")
print(stdout.read().decode('utf-8', errors='ignore'))

# 检查日志
stdin, stdout, stderr = client.exec_command("journalctl -u anchor-finance -n 20 --no-pager")
print("\n=== 服务日志 ===")
print(stdout.read().decode('utf-8', errors='ignore'))

# 测试健康检查
stdin, stdout, stderr = client.exec_command("curl -s http://localhost:8080/health")
print("\n=== 健康检查 ===")
print(stdout.read().decode('utf-8', errors='ignore'))

# 测试登录
stdin, stdout, stderr = client.exec_command("""curl -s -X POST http://localhost:8080/api/admin/login -H "Content-Type: application/json" -d '{"username":"admin","password":"admin123"}'""")
print("\n=== 登录测试 ===")
print(stdout.read().decode('utf-8', errors='ignore'))

client.close()
