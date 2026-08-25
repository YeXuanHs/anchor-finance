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
for _ in range(8):
    time.sleep(2)
    login = run(f'curl -s -X POST {BASE}/api/admin/login -H "Content-Type: application/json" -d \'{{"username":"admin","password":"admin123"}}\'')
    if '"code"' in login:
        break

try:
    token = json.loads(login)["data"]["token"]
except Exception as e:
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

print("=== 用户管理 ===")
test("用户列表", f'curl -s "{BASE}/api/admin/users?page=1&page_size=5" {AUTH}')
test("用户余额日志", f'curl -s "{BASE}/api/admin/users/1/balance-logs?page=1&page_size=5" {AUTH}')

print("\n=== 服务/订单/账单 ===")
test("服务列表", f'curl -s "{BASE}/api/admin/services?page=1&page_size=5" {AUTH}')
test("订单列表", f'curl -s "{BASE}/api/admin/orders?page=1&page_size=5" {AUTH}')
test("账单列表", f'curl -s "{BASE}/api/admin/invoices?page=1&page_size=5" {AUTH}')
test("交易流水", f'curl -s "{BASE}/api/admin/transactions?page=1&page_size=5" {AUTH}')

print("\n=== 工单/产品/供应商 ===")
test("工单列表", f'curl -s "{BASE}/api/admin/tickets?page=1&page_size=5" {AUTH}')
test("产品列表", f'curl -s "{BASE}/api/admin/products?page=1&page_size=5" {AUTH}')
test("供应商列表", f'curl -s "{BASE}/api/admin/suppliers?page=1&page_size=5" {AUTH}')

print("\n=== 新增模块回归 ===")
test("插件列表", f'curl -s "{BASE}/api/admin/plugins" {AUTH}')
test("任务队列概览", f'curl -s "{BASE}/api/admin/task-queue/overview" {AUTH}')
test("主题列表", f'curl -s "{BASE}/api/admin/themes?page=1&page_size=5" {AUTH}')
test("财务账本", f'curl -s "{BASE}/api/admin/finance/ledger" {AUTH}')
test("分类日志-短信", f'curl -s "{BASE}/api/admin/logs/sms?page=1&page_size=5" {AUTH}')
test("支付网关插件", f'curl -s "{BASE}/api/admin/payment-gateways" {AUTH}')

# 登录日志（管理员登录应该被记录了）
test("管理员登录日志", f'curl -s "{BASE}/api/admin/logs/admin-login?page=1&page_size=5" {AUTH}')

print(f"\n=== 最终: 通过 {passed} / 失败 {failed} ===")
client.close()
