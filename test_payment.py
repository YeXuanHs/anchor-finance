import paramiko
import json

client = paramiko.SSHClient()
client.set_missing_host_key_policy(paramiko.AutoAddPolicy())
client.connect('45.207.210.235', username='root', password='iswlBSLY8118')

# 用户登录
cmd = """curl -s -X POST http://localhost:8080/api/client/login -H "Content-Type: application/json" -d '{"username":"testuser","password":"123456"}'"""
stdin, stdout, stderr = client.exec_command(cmd)
login_result = json.loads(stdout.read().decode('utf-8', errors='ignore'))
token = login_result.get('data', {}).get('token', '')

# 测试支付记录API
apis = [
    ("GET", "/api/client/payments", "支付记录列表"),
    ("GET", "/api/client/payments/summary", "支付统计"),
]

print("=== 测试支付记录API ===")
for method, path, desc in apis:
    cmd = f'curl -s "http://localhost:8080{path}" -H "Authorization: Bearer {token}"'
    stdin, stdout, stderr = client.exec_command(cmd)
    result = stdout.read().decode('utf-8', errors='ignore')
    try:
        data = json.loads(result)
        if data.get('code') == 0:
            print(f"  {desc}: OK")
        else:
            print(f"  {desc}: FAIL - {data.get('message')}")
    except:
        print(f"  {desc}: FAIL")

client.close()
