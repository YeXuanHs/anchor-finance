import paramiko
import json
import time
import sys

sys.stdout.reconfigure(encoding='utf-8')
client = paramiko.SSHClient()
client.set_missing_host_key_policy(paramiko.AutoAddPolicy())
client.connect('45.207.210.235', username='root', password='iswlBSLY8118')

def run(cmd):
    i, o, e = client.exec_command(cmd)
    return o.read().decode('utf-8', errors='replace').strip()

BASE = "http://127.0.0.1:8080"
# 等待服务
for _ in range(10):
    time.sleep(2)
    login = run(f'curl -s -X POST {BASE}/api/admin/login -H "Content-Type: application/json" -d \'{{"username":"admin","password":"admin123"}}\'')
    if '"code"' in login:
        break
try:
    token = json.loads(login)["data"]["token"]
except:
    print("❌ 登录失败:", login)
    sys.exit(1)

AUTH = f'-H "Authorization: Bearer {token}"'
passed = 0
failed = 0

def test(name, cmd, expect_code=0):
    global passed, failed
    resp = run(cmd)
    try:
        d = json.loads(resp)
        if d.get("code") == expect_code:
            print(f"  ✅ {name}")
            passed += 1
            return d
        else:
            print(f"  ❌ {name}: code={d.get('code')} msg={d.get('message')}")
            failed += 1
            return None
    except Exception:
        print(f"  ❌ {name}: 无法解析 {resp[:80]}")
        failed += 1
        return None

print("=== 核心功能回归（验证没被破坏） ===")
test("登录成功", f'curl -s -X POST {BASE}/api/admin/login -H "Content-Type: application/json" -d \'{{"username":"admin","password":"admin123"}}\'')
test("用户列表", f'curl -s "{BASE}/api/admin/users?page=1&page_size=5" {AUTH}')
test("订单列表", f'curl -s "{BASE}/api/admin/orders?page=1&page_size=5" {AUTH}')
test("产品列表", f'curl -s "{BASE}/api/admin/products?page=1&page_size=5" {AUTH}')
test("工单列表", f'curl -s "{BASE}/api/admin/tickets?page=1&page_size=5" {AUTH}')
test("仪表盘", f'curl -s "{BASE}/api/admin/dashboard/stats" {AUTH}')

print("\n=== 安全防护验证 ===")
test("0元购防护-产品负价拦截", f'curl -s -X POST {BASE}/api/admin/products {AUTH} -H "Content-Type: application/json" -d \'{{"name":"x","price":-100}}\'', 400)
test("0元购防护-订单负价拦截", f'curl -s -X POST {BASE}/api/admin/orders {AUTH} -H "Content-Type: application/json" -d \'{{"user_id":1,"product_id":1,"product_name":"t","amount":-50}}\'', 400)

print("\n=== 信用额配置 ===")
test("获取信用额配置", f'curl -s {BASE}/api/admin/credit-limits/config {AUTH}')
test("保存信用额配置", f'curl -s -X POST {BASE}/api/admin/credit-limits/config {AUTH} -H "Content-Type: application/json" -d \'{{"enabled":true,"default_amount":100,"max_amount":5000}}\'')
test("最大额度负值拦截", f'curl -s -X POST {BASE}/api/admin/credit-limits/config {AUTH} -H "Content-Type: application/json" -d \'{{"max_amount":-1}}\'', 400)

print(f"\n=== 最终: 通过 {passed} / 失败 {failed} ===")
client.close()
