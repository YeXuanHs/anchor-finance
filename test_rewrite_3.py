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

# 服务已外部启动，等待就绪
print("=== 等待服务 ===")
import time
for _ in range(5):
    time.sleep(2)
    if json.loads(run("curl -s http://127.0.0.1:8080/health")) or True:
        break

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

# ===== configurable-options =====
print("\n=== 全局可配置项 ===")
d = test("创建-可配置项", f'curl -s -X POST {BASE}/api/admin/configurable-options {AUTH} -H "Content-Type: application/json" -d \'{{"name":"内存","type":"select","options":"[2,4,8]"}}\'')
if d:
    oid = d["data"]["id"]
    test("列表-可配置项", f'curl -s "{BASE}/api/admin/configurable-options?page=1&page_size=10" {AUTH}')
    test("更新-可配置项", f'curl -s -X PUT {BASE}/api/admin/configurable-options/{oid} {AUTH} -H "Content-Type: application/json" -d \'{{"name":"内存GB"}}\'')
    test("删除-可配置项", f'curl -s -X DELETE {BASE}/api/admin/configurable-options/{oid} {AUTH}')

# ===== oauth-providers =====
print("\n=== 第三方登录 ===")
d = test("创建-OAuth", f'curl -s -X POST {BASE}/api/admin/oauth-providers {AUTH} -H "Content-Type: application/json" -d \'{{"name":"微信","code":"wechat","icon":"wx.png"}}\'')
if d:
    pid = d["data"]["id"]
    test("列表-OAuth", f'curl -s "{BASE}/api/admin/oauth-providers?page=1&page_size=10" {AUTH}')
    test("重复code拦截", f'curl -s -X POST {BASE}/api/admin/oauth-providers {AUTH} -H "Content-Type: application/json" -d \'{{"name":"微信2","code":"wechat"}}\'', 400)
    test("更新-OAuth", f'curl -s -X PUT {BASE}/api/admin/oauth-providers/{pid} {AUTH} -H "Content-Type: application/json" -d \'{{"name":"微信登录"}}\'')
    test("删除-OAuth", f'curl -s -X DELETE {BASE}/api/admin/oauth-providers/{pid} {AUTH}')

# ===== custom-template-fields =====
print("\n=== 官网自定义字段 ===")
d = test("创建-模板字段", f'curl -s -X POST {BASE}/api/admin/custom-template-fields {AUTH} -H "Content-Type: application/json" -d \'{{"page":"home","name":"标题","key":"title","value":"欢迎"}}\'')
if d:
    fid = d["data"]["id"]
    test("列表-模板字段", f'curl -s "{BASE}/api/admin/custom-template-fields?page=home&page_size=10" {AUTH}')
    test("重复字段拦截", f'curl -s -X POST {BASE}/api/admin/custom-template-fields {AUTH} -H "Content-Type: application/json" -d \'{{"page":"home","name":"x","key":"title"}}\'', 400)
    test("更新-模板字段", f'curl -s -X PUT {BASE}/api/admin/custom-template-fields/{fid} {AUTH} -H "Content-Type: application/json" -d \'{{"value":"新欢迎"}}\'')
    test("删除-模板字段", f'curl -s -X DELETE {BASE}/api/admin/custom-template-fields/{fid} {AUTH}')

# ===== traffic-packages =====
print("\n=== 流量包 ===")
d = test("创建-流量包", f'curl -s -X POST {BASE}/api/admin/traffic-packages {AUTH} -H "Content-Type: application/json" -d \'{{"name":"1GB包","volume":1024,"price":10}}\'')
if d:
    tid = d["data"]["id"]
    test("列表-流量包", f'curl -s "{BASE}/api/admin/traffic-packages?page=1&page_size=10" {AUTH}')
    test("0元购防护-负价", f'curl -s -X POST {BASE}/api/admin/traffic-packages {AUTH} -H "Content-Type: application/json" -d \'{{"name":"bad","volume":1,"price":-5}}\'', 400)
    test("更新-流量包", f'curl -s -X PUT {BASE}/api/admin/traffic-packages/{tid} {AUTH} -H "Content-Type: application/json" -d \'{{"price":15}}\'')
    test("删除-流量包", f'curl -s -X DELETE {BASE}/api/admin/traffic-packages/{tid} {AUTH}')
test("列表-流量日志", f'curl -s "{BASE}/api/admin/traffic-logs?page=1&page_size=10" {AUTH}')

print(f"\n=== 最终: 通过 {passed} / 失败 {failed} ===")
client.close()
