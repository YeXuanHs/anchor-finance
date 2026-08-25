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

# 部署最新Release
print("=== 部署 ===")
run("cd /opt/anchor-finance && URL=$(curl -s https://api.github.com/repos/YeXuanHs/anchor-finance/releases/latest | python3 -c 'import sys,json; d=json.load(sys.stdin); print([a[\"browser_download_url\"] for a in d.get(\"assets\",[]) if \"anchor-finance-server\" in a[\"name\"]][0])') && curl -L -o server/anchor-finance-new \"$URL\" && mv server/anchor-finance-new server/anchor-finance && chmod +x server/anchor-finance")
run("fuser -k 8080/tcp 2>/dev/null")
time.sleep(2)
run("systemctl restart anchor-finance")
time.sleep(5)
print("status:", run("systemctl is-active anchor-finance"))

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
        print(f"  ❌ {name}: 无法解析 {resp[:80]}")
        failed += 1
        return None

# ===== 任务队列 =====
print("\n=== 任务队列 ===")
test("概览", f'curl -s {BASE}/api/admin/task-queue/overview {AUTH}')
test("列表", f'curl -s "{BASE}/api/admin/task-queue?page=1&page_size=10" {AUTH}')

# ===== 二次验证配置 =====
print("\n=== 二次验证配置 ===")
test("获取", f'curl -s {BASE}/api/admin/two-factor-config {AUTH}')
test("更新", f'curl -s -X PUT {BASE}/api/admin/two-factor-config {AUTH} -H "Content-Type: application/json" -d \'{{"enabled":"true"}}\'')

# ===== 销售设置 =====
print("\n=== 销售设置 ===")
test("销售配置-获取", f'curl -s {BASE}/api/admin/sales-config {AUTH}')
d = test("销售分组-创建", f'curl -s -X POST {BASE}/api/admin/sales-groups {AUTH} -H "Content-Type: application/json" -d \'{{"name":"华东组","commission":10}}\'')
if d:
    sgid = d["data"]["id"]
    test("销售分组-列表", f'curl -s "{BASE}/api/admin/sales-groups?page=1&page_size=10" {AUTH}')
    test("销售分组-佣金越界拦截", f'curl -s -X POST {BASE}/api/admin/sales-groups {AUTH} -H "Content-Type: application/json" -d \'{{"name":"bad","commission":150}}\'', 400)
    test("销售分组-更新", f'curl -s -X PUT {BASE}/api/admin/sales-groups/{sgid} {AUTH} -H "Content-Type: application/json" -d \'{{"commission":15}}\'')
    test("销售分组-删除", f'curl -s -X DELETE {BASE}/api/admin/sales-groups/{sgid} {AUTH}')

# ===== 主题 =====
print("\n=== 主题模板 ===")
d = test("主题-创建", f'curl -s -X POST {BASE}/api/admin/themes {AUTH} -H "Content-Type: application/json" -d \'{{"name":"默认主题","code":"default"}}\'')
if d:
    thid = d["data"]["id"]
    test("主题-列表", f'curl -s "{BASE}/api/admin/themes?page=1&page_size=10" {AUTH}')
    test("主题-设为默认", f'curl -s -X POST {BASE}/api/admin/themes/{thid}/set-default {AUTH}')
    test("主题-激活详情", f'curl -s {BASE}/api/admin/themes/active {AUTH}')
    test("主题-重复code拦截", f'curl -s -X POST {BASE}/api/admin/themes {AUTH} -H "Content-Type: application/json" -d \'{{"name":"x","code":"default"}}\'', 400)
    # 删除默认主题应被拦截
    test("主题-删除默认被拦", f'curl -s -X DELETE {BASE}/api/admin/themes/{thid} {AUTH}', 400)
    # 建非默认再删
    d2 = test("主题-创建2", f'curl -s -X POST {BASE}/api/admin/themes {AUTH} -H "Content-Type: application/json" -d \'{{"name":"第二","code":"two"}}\'')
    if d2:
        test("主题-删除2", f'curl -s -X DELETE {BASE}/api/admin/themes/{d2["data"]["id"]} {AUTH}')

# ===== 工单传递规则 =====
print("\n=== 工单传递规则 ===")
d = test("规则-创建", f'curl -s -X POST {BASE}/api/admin/ticket-rules {AUTH} -H "Content-Type: application/json" -d \'{{"name":"高优先级","condition_type":"priority","condition_value":"high","action_type":"notify","action_value":"admin"}}\'')
if d:
    rid = d["data"]["id"]
    test("规则-列表", f'curl -s "{BASE}/api/admin/ticket-rules?page=1&page_size=10" {AUTH}')
    test("规则-更新", f'curl -s -X PUT {BASE}/api/admin/ticket-rules/{rid} {AUTH} -H "Content-Type: application/json" -d \'{{"name":"改"}}\'')
    test("规则-删除", f'curl -s -X DELETE {BASE}/api/admin/ticket-rules/{rid} {AUTH}')

# ===== 商品订购/财务配置 =====
print("\n=== 商品订购/财务配置 ===")
test("订单配置-获取", f'curl -s {BASE}/api/admin/order-config {AUTH}')
test("订单配置-更新", f'curl -s -X PUT {BASE}/api/admin/order-config {AUTH} -H "Content-Type: application/json" -d \'{{"allow_purchase":"1"}}\'')
test("财务配置-获取", f'curl -s {BASE}/api/admin/finance-config {AUTH}')
test("财务配置-更新", f'curl -s -X PUT {BASE}/api/admin/finance-config {AUTH} -H "Content-Type: application/json" -d \'{{"currency":"CNY"}}\'')

print(f"\n=== 最终: 通过 {passed} / 失败 {failed} ===")
client.close()
