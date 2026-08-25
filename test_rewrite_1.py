import paramiko
import json
import time
import sys

sys.stdout.reconfigure(encoding='utf-8')

client = paramiko.SSHClient()
client.set_missing_host_key_policy(paramiko.AutoAddPolicy())
client.connect('45.207.210.235', username='root', password='iswlBSLY8118')

def run(cmd):
    stdin, stdout, stderr = client.exec_command(cmd, timeout=15)
    return stdout.read().decode('utf-8', errors='replace').strip()

# 等待服务（不重启，假设外部已重启）
print("等待服务就绪...")
for i in range(5):
    time.sleep(2)
    status = run("systemctl is-active anchor-finance")
    port = run("ss -tlnp | grep 8080 | head -1")
    if status == "active" and port:
        print(f"  服务就绪: status=active port=OK")
        break

BASE = "http://127.0.0.1:8080"
login = run(f'curl -s -X POST {BASE}/api/admin/login -H "Content-Type: application/json" -d \'{{"username":"admin","password":"admin123"}}\'')
print("登录:", login[:150])
token = json.loads(login)["data"]["token"]
AUTH = f'-H "Authorization: Bearer {token}"'
passed = 0
failed = 0

def test(name, cmd, expect_code=0):
    global passed, failed
    resp = run(cmd)
    try:
        data = json.loads(resp)
        if data.get("code") == expect_code:
            print(f"  ✅ {name}")
            passed += 1
            return data
        else:
            print(f"  ❌ {name}: code={data.get('code')} msg={data.get('message')}")
            failed += 1
            return None
    except Exception:
        print(f"  ❌ {name}: 无法解析 {resp[:100]}")
        failed += 1
        return None

print("\n=== 黑名单 ===")
d = test("黑名单-创建", f'curl -s -X POST {BASE}/api/admin/blacklist {AUTH} -H "Content-Type: application/json" -d \'{{"type":"ip","value":"5.6.7.8","reason":"t"}}\'')
if d: test("黑名单-列表", f'curl -s "{BASE}/api/admin/blacklist?page=1&page_size=10" {AUTH}'); test("黑名单-删除", f'curl -s -X DELETE {BASE}/api/admin/blacklist/{d["data"]["id"]} {AUTH}')

print("\n=== 邮件模板 ===")
d = test("邮件模板-创建", f'curl -s -X POST {BASE}/api/admin/email-templates {AUTH} -H "Content-Type: application/json" -d \'{{"name":"t1","subject":"s","content":"h"}}\'')
if d:
    eid = d["data"]["id"]
    test("邮件模板-列表", f'curl -s {BASE}/api/admin/email-templates {AUTH}')
    test("邮件模板-详情", f'curl -s {BASE}/api/admin/email-templates/{eid} {AUTH}')
    test("邮件模板-更新", f'curl -s -X PUT {BASE}/api/admin/email-templates/{eid} {AUTH} -H "Content-Type: application/json" -d \'{{"subject":"new"}}\'')
    test("邮件模板-删除", f'curl -s -X DELETE {BASE}/api/admin/email-templates/{eid} {AUTH}')

print("\n=== 短信模板 ===")
d = test("短信模板-创建", f'curl -s -X POST {BASE}/api/admin/sms-templates {AUTH} -H "Content-Type: application/json" -d \'{{"name":"v","code":"vc2","content":"c"}}\'')
if d:
    sid = d["data"]["id"]
    test("短信模板-重复code拦截", f'curl -s -X POST {BASE}/api/admin/sms-templates {AUTH} -H "Content-Type: application/json" -d \'{{"name":"dup","code":"vc2","content":"x"}}\'', 400)
    test("短信模板-更新", f'curl -s -X PUT {BASE}/api/admin/sms-templates/{sid} {AUTH} -H "Content-Type: application/json" -d \'{{"name":"v2"}}\'')
    test("短信模板-删除", f'curl -s -X DELETE {BASE}/api/admin/sms-templates/{sid} {AUTH}')

print("\n=== 友情链接 ===")
d = test("友情链接-创建", f'curl -s -X POST {BASE}/api/admin/friendly-links {AUTH} -H "Content-Type: application/json" -d \'{{"name":"link","url":"https://ex.com"}}\'')
if d:
    fid = d["data"]["id"]
    test("友情链接-列表", f'curl -s "{BASE}/api/admin/friendly-links?page=1&page_size=10" {AUTH}')
    test("友情链接-删除", f'curl -s -X DELETE {BASE}/api/admin/friendly-links/{fid} {AUTH}')

print(f"\n=== 最终: 通过 {passed} / 失败 {failed} ===")
client.close()
