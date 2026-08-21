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
print(f'Token: {token[:50]}...')

# 测试几个API
apis = [
    '/api/admin/auth/info',
    '/api/admin/dashboard/stats',
    '/api/admin/users?page=1&page_size=5',
]

for api in apis:
    cmd = f'curl -s "http://localhost:8080{api}" -H "Authorization: Bearer {token}"'
    stdin, stdout, stderr = client.exec_command(cmd)
    result = stdout.read().decode('utf-8', errors='ignore')
    print(f'\n{api}:')
    print(result[:200])

client.close()
