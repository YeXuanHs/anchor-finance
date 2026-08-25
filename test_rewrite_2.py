import paramiko
import json
import sys

sys.stdout.reconfigure(encoding='utf-8')
client = paramiko.SSHClient()
client.set_missing_host_key_policy(paramiko.AutoAddPolicy())
client.connect('45.207.210.235', username='root', password='iswlBSLY8118')

def run(cmd):
    i, o, e = client.exec_command(cmd)
    return o.read().decode('utf-8', errors='replace').strip()

BASE = "http://127.0.0.1:8080"
login = run(f'curl -s -X POST {BASE}/api/admin/login -H "Content-Type: application/json" -d \'{{"username":"admin","password":"admin123"}}\'')
token = json.loads(login)["data"]["token"]
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
        print(f"  ❌ {name}: 无法解析 {resp[:100]}")
        failed += 1
        return None

# 需要一个测试用户
print("\n=== 准备测试用户 ===")
d = test("创建测试用户", f'curl -s -X POST {BASE}/api/admin/users {AUTH} -H "Content-Type: application/json" -d \'{{"username":"testuser","email":"test@test.com","password":"test123","phone":"13800138000"}}\'')
uid = d["data"]["id"] if d else 1

# ===== 合同 =====
print("\n=== 合同 ===")
d = test("合同-创建", f'curl -s -X POST {BASE}/api/admin/contracts {AUTH} -H "Content-Type: application/json" -d \'{{"user_id":{uid},"title":"测试合同","content":"条款...","type":"service"}}\'')
if d:
    cid = d["data"]["id"]
    test("合同-列表", f'curl -s "{BASE}/api/admin/contracts?page=1&page_size=10" {AUTH}')
    test("合同-详情", f'curl -s "{BASE}/api/admin/contracts/{cid}" {AUTH}')
    test("合同-更新", f'curl -s -X PUT {BASE}/api/admin/contracts/{cid} {AUTH} -H "Content-Type: application/json" -d \'{{"title":"改名"}}\'')
    test("合同-签署pipeline校验", f'curl -s -X POST {BASE}/api/admin/contracts/{cid}/sign {AUTH}', 400)  # draft不能签
    test("合同-取消", f'curl -s -X POST {BASE}/api/admin/contracts/{cid}/cancel {AUTH}')
    test("合同-删除", f'curl -s -X DELETE {BASE}/api/admin/contracts/{cid} {AUTH}')

print("\n=== 合同模板 ===")
d = test("合同模板-创建", f'curl -s -X POST {BASE}/api/admin/contract-templates {AUTH} -H "Content-Type: application/json" -d \'{{"name":"模板A","content":"hello username"}}\'')
if d:
    tid = d["data"]["id"]
    test("合同模板-列表", f'curl -s {BASE}/api/admin/contract-templates {AUTH}')
    test("合同模板-更新", f'curl -s -X PUT {BASE}/api/admin/contract-templates/{tid} {AUTH} -H "Content-Type: application/json" -d \'{{"name":"改名A"}}\'')
    test("合同模板-删除", f'curl -s -X DELETE {BASE}/api/admin/contract-templates/{tid} {AUTH}')

# ===== 营销推送 =====
print("\n=== 营销推送 ===")
d = test("营销-创建", f'curl -s -X POST {BASE}/api/admin/marketing/pushes {AUTH} -H "Content-Type: application/json" -d \'{{"title":"活动","content":"大促","type":"sms","target_type":"all"}}\'')
if d:
    mid = d["data"]["id"]
    test("营销-列表", f'curl -s "{BASE}/api/admin/marketing/pushes?page=1&page_size=10" {AUTH}')
    # 发送走PHP引擎，可能离线但接口要通
    test("营销-发送", f'curl -s -X POST {BASE}/api/admin/marketing/pushes/{mid}/send {AUTH}')
    test("营销-删除", f'curl -s -X DELETE {BASE}/api/admin/marketing/pushes/{mid} {AUTH}')

# ===== 取消请求 =====
print("\n=== 取消请求 ===")
# 直接查库生成一个取消请求测试
test("取消请求-列表", f'curl -s "{BASE}/api/admin/cancel-requests?page=1&page_size=10" {AUTH}')

# ===== 销售统计 =====
print("\n=== 销售统计 ===")
test("销售-统计", f'curl -s "{BASE}/api/admin/sales/statistics" {AUTH}')
test("销售-记录", f'curl -s "{BASE}/api/admin/sales/records?page=1&page_size=10" {AUTH}')

print(f"\n=== 最终: 通过 {passed} / 失败 {failed} ===")
client.close()
