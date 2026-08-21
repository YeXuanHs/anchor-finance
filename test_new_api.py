import paramiko
import json

client = paramiko.SSHClient()
client.set_missing_host_key_policy(paramiko.AutoAddPolicy())
client.connect('45.207.210.235', username='root', password='iswlBSLY8118')

# 登录
cmd = """curl -s -X POST http://localhost:8080/api/admin/login -H "Content-Type: application/json" -d '{"username":"admin","password":"admin123"}'"""
stdin, stdout, stderr = client.exec_command(cmd)
login_result = json.loads(stdout.read().decode('utf-8', errors='ignore'))
token = login_result.get('data', {}).get('token', '')

# 测试用户操作日志API
print("=== 测试用户操作日志API ===")
cmd = f'curl -s "http://localhost:8080/api/admin/users/1/operation-logs?page=1&page_size=10" -H "Authorization: Bearer {token}"'
stdin, stdout, stderr = client.exec_command(cmd)
result = stdout.read().decode('utf-8', errors='ignore')
print(result)

# 验证返回格式
data = json.loads(result)
if data.get('code') == 0:
    print("\n✅ 用户操作日志API测试通过！")
else:
    print(f"\n❌ 测试失败: {data.get('message')}")

client.close()
