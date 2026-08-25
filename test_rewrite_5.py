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

# 服务已外部启动
BASE = "http://127.0.0.1:8080"
print("=== 等待服务 ===")
for _ in range(8):
    time.sleep(2)
    login = run(f'curl -s -X POST {BASE}/api/admin/login -H "Content-Type: application/json" -d \'{{"username":"admin","password":"admin123"}}\'')
    if "code" in login:
        break
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

# ===== 分类日志 =====
print("\n=== 分类日志 ===")
test("短信日志", f'curl -s "{BASE}/api/admin/logs/sms?page=1&page_size=10" {AUTH}')
test("邮件日志", f'curl -s "{BASE}/api/admin/logs/email?page=1&page_size=10" {AUTH}')
test("API日志", f'curl -s "{BASE}/api/admin/logs/api?page=1&page_size=10" {AUTH}')
test("定时日志", f'curl -s "{BASE}/api/admin/logs/cron?page=1&page_size=10" {AUTH}')
test("管理员登录日志", f'curl -s "{BASE}/api/admin/logs/admin-login?page=1&page_size=10" {AUTH}')
test("站内信日志", f'curl -s "{BASE}/api/admin/logs/notification?page=1&page_size=10" {AUTH}')

# ===== 客户分组 =====
print("\n=== 客户分组 ===")
d = test("创建分组", f'curl -s -X POST {BASE}/api/admin/customer-groups {AUTH} -H "Content-Type: application/json" -d \'{{"name":"VIP","discount":90}}\'')
if d:
    gid = d["data"]["id"]
    test("分组列表", f'curl -s "{BASE}/api/admin/customer-groups?page=1&page_size=10" {AUTH}')
    test("折扣越界拦截", f'curl -s -X POST {BASE}/api/admin/customer-groups {AUTH} -H "Content-Type: application/json" -d \'{{"name":"bad","discount":150}}\'', 400)
    test("更新分组", f'curl -s -X PUT {BASE}/api/admin/customer-groups/{gid} {AUTH} -H "Content-Type: application/json" -d \'{{"discount":85}}\'')
    test("删除分组", f'curl -s -X DELETE {BASE}/api/admin/customer-groups/{gid} {AUTH}')

# ===== catalogs =====
print("\n=== CPU/实例规格目录 ===")
test("CPU目录", f'curl -s {BASE}/api/admin/cpu-model-catalog {AUTH}')
test("实例规格目录", f'curl -s {BASE}/api/admin/instance-spec-catalog {AUTH}')

# ===== 定时任务 =====
print("\n=== 定时任务 ===")
test("cron任务列表", f'curl -s {BASE}/api/admin/cron-tasks {AUTH}')
test("调度概览", f'curl -s {BASE}/api/admin/schedules/overview {AUTH}')
test("调度运行列表", f'curl -s "{BASE}/api/admin/schedule-runs?page=1&page_size=10" {AUTH}')

# ===== 插件域 =====
print("\n=== 插件域 ===")
test("支付网关", f'curl -s {BASE}/api/admin/payment-gateways {AUTH}')
test("短信提供商", f'curl -s {BASE}/api/admin/sms-providers {AUTH}')
test("邮件提供商", f'curl -s {BASE}/api/admin/mail-providers {AUTH}')
test("认证提供商", f'curl -s {BASE}/api/admin/certification-providers {AUTH}')
test("服务器模块", f'curl -s {BASE}/api/admin/server-modules {AUTH}')

print(f"\n=== 最终: 通过 {passed} / 失败 {failed} ===")
client.close()
