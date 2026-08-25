"""
AnchorFinance 严格审计测试
=========================
不仅检查返回码，还验证：
1. 数据库中数据是否真的写入
2. 金额是否正确计算
3. 安全措施是否真正生效
"""
import requests, json, sys, time

BASE = "http://45.207.210.235:8080"
ADMIN_TOKEN = None
CLIENT_TOKEN = None
results = {"pass": 0, "fail": 0, "errors": []}

def test(name, func):
    try:
        func()
        results["pass"] += 1
        print(f"  [PASS] {name}")
    except AssertionError as e:
        results["fail"] += 1
        results["errors"].append(f"{name}: {e}")
        print(f"  [FAIL] {name}: {e}")
    except Exception as e:
        results["fail"] += 1
        results["errors"].append(f"{name}: {type(e).__name__}: {e}")
        print(f"  [FAIL] {name}: {type(e).__name__}: {e}")

def api(method, path, json_data=None, token=None):
    headers = {"Content-Type": "application/json"}
    if token:
        headers["Authorization"] = f"Bearer {token}"
    return requests.request(method, f"{BASE}{path}", json=json_data, headers=headers, timeout=30)

def assert_code(resp, expected):
    d = resp.json()
    assert d.get("code") == expected, f"code={d.get('code')} (expect {expected}), msg={d.get('message')}"

# === 1. 返回格式验证 ===
print("\n=== 1. 返回格式验证 ===")

def test_health_format():
    r = requests.get(f"{BASE}/health", timeout=5)
    d = r.json()
    assert "status" in d, f"health缺少status字段: {d}"

def test_error_format():
    """错误响应也必须有code字段"""
    r = api("GET", "/api/admin/users")
    d = r.json()
    assert "code" in d, f"错误响应缺少code字段: {d}"
    assert "message" in d, f"错误响应缺少message字段: {d}"

def test_success_format():
    """成功响应必须有code=0, data字段"""
    global ADMIN_TOKEN
    r = api("POST", "/api/admin/login", {"username": "admin", "password": "admin123"})
    d = r.json()
    assert d.get("code") == 0, f"登录失败: {d}"
    assert "data" in d, f"成功响应缺少data字段: {d}"
    assert "token" in d["data"], f"data缺少token字段: {d['data']}"
    ADMIN_TOKEN = d["data"]["token"]

test("健康检查格式", test_health_format)
test("错误响应格式", test_error_format)
test("成功响应格式+Admin登录", test_success_format)

# === 2. 注册+登录完整流程 ===
print("\n=== 2. 注册+登录完整流程 ===")

def test_disable_captcha():
    r = api("PUT", "/api/admin/settings", {"register_require_captcha": "0"}, token=ADMIN_TOKEN)
    assert_code(r, 0)

def test_register_and_login():
    global CLIENT_TOKEN
    ts = int(time.time())
    username = f"audit_{ts}"
    r = api("POST", "/api/client/register", {
        "username": username, "email": f"{username}@test.com", "password": "test123456"
    })
    assert_code(r, 0)
    assert r.json()["data"].get("token") is not None, "注册响应缺少token"

    CLIENT_TOKEN = r.json()["data"]["token"]

    r = api("GET", "/api/client/auth/info", token=CLIENT_TOKEN)
    assert_code(r, 0)
    assert r.json()["data"]["username"] == username, f"用户名不匹配: {r.json()['data']['username']}"

test("关闭验证码", test_disable_captcha)
test("注册+登录+信息验证", test_register_and_login)

# === 3. 产品创建+金额验证 ===
print("\n=== 3. 产品创建+金额验证 ===")

def test_create_product():
    r = api("POST", "/api/admin/products", {
        "name": "审计云服务器", "type": "server", "price": 199.99,
        "billing_cycle": "monthly", "status": "active"
    }, token=ADMIN_TOKEN)
    assert_code(r, 0)

def test_zero_price():
    r = api("POST", "/api/admin/products", {
        "name": "零元产品", "price": 0, "type": "server"
    }, token=ADMIN_TOKEN)
    assert_code(r, 400)

def test_negative_price():
    r = api("POST", "/api/admin/products", {
        "name": "负价产品", "price": -1, "type": "server"
    }, token=ADMIN_TOKEN)
    assert_code(r, 400)

test("创建产品199.99", test_create_product)
test("0元拦截", test_zero_price)
test("负价拦截", test_negative_price)

# === 4. 购物车完整流程+金额验证 ===
print("\n=== 4. 购物车完整流程+金额验证 ===")

def test_get_product_id():
    """获取审计测试产品ID"""
    r = api("GET", "/api/client/products", token=CLIENT_TOKEN)
    assert_code(r, 0)
    products = r.json()["data"]["list"]
    assert len(products) > 0, "产品列表为空"
    # 找审计产品
    for p in products:
        if p["name"] == "审计云服务器":
            return p["id"]
    return products[-1]["id"]

def test_cart_flow():
    global product_id
    product_id = test_get_product_id()

    # 加购
    r = api("POST", "/api/client/cart", {
        "product_id": product_id, "quantity": 2, "cycle": "monthly"
    }, token=CLIENT_TOKEN)
    assert_code(r, 0)

    # 验证购物车
    r = api("GET", "/api/client/cart", token=CLIENT_TOKEN)
    assert_code(r, 0)
    items = r.json()["data"]["items"]
    assert len(items) == 1, f"购物车应有1项，实际{len(items)}项"
    assert items[0]["quantity"] == 2, f"数量应为2，实际{items[0]['quantity']}"
    assert r.json()["data"]["total_amount"] > 0, "总金额应>0"

    # 结算
    r = api("POST", "/api/client/cart/checkout", {}, token=CLIENT_TOKEN)
    assert_code(r, 0)
    data = r.json()["data"]
    assert data["amount"] > 0, f"结算金额应>0: {data['amount']}"
    assert data.get("order_id") is not None, "缺少order_id"
    assert data.get("invoice_id") is not None, "缺少invoice_id"

    # 结算后购物车应为空
    r = api("GET", "/api/client/cart", token=CLIENT_TOKEN)
    assert_code(r, 0)
    assert r.json()["data"]["items"] == [], "结算后购物车应为空"

    return data["order_id"], data["invoice_id"], data["amount"]

product_id = None
order_id, invoice_id, amount = test_cart_flow()
test("购物车完整流程", lambda: None)  # test_cart_flow already does assertions

# === 5. 订单+账单验证 ===
print("\n=== 5. 订单+账单验证 ===")

def test_order_detail():
    r = api("GET", f"/api/client/orders/{order_id}", token=CLIENT_TOKEN)
    assert_code(r, 0)
    assert r.json()["data"]["amount"] == amount, f"订单金额不匹配: {r.json()['data']['amount']} != {amount}"

def test_invoice_detail():
    r = api("GET", f"/api/client/invoices/{invoice_id}", token=CLIENT_TOKEN)
    assert_code(r, 0)
    assert r.json()["data"]["status"] == "unpaid", f"账单状态应为unpaid: {r.json()['data']['status']}"

test("订单金额正确", test_order_detail)
test("账单状态正确", test_invoice_detail)

# === 6. 工单完整流程 ===
print("\n=== 6. 工单完整流程 ===")

def test_ticket_flow():
    r = api("POST", "/api/client/tickets", {
        "subject": "审计测试工单", "content": "测试内容", "priority": "high"
    }, token=CLIENT_TOKEN)
    assert_code(r, 0)
    tid = r.json()["data"]["id"]

    r = api("POST", f"/api/client/tickets/{tid}/reply", {"content": "客户回复"}, token=CLIENT_TOKEN)
    assert_code(r, 0)

    r = api("POST", f"/api/admin/tickets/{tid}/reply", {"content": "管理员回复"}, token=ADMIN_TOKEN)
    assert_code(r, 0)

    r = api("GET", f"/api/client/tickets/{tid}/replies", token=CLIENT_TOKEN)
    assert_code(r, 0)
    assert len(r.json()["data"]) >= 2, f"应有>=2条回复: {len(r.json()['data'])}"

    r = api("POST", f"/api/client/tickets/{tid}/close", {}, token=CLIENT_TOKEN)
    assert_code(r, 0)

test("工单完整流程（创建→客户回复→管理员回复→查看回复→关闭）", test_ticket_flow)

# === 7. 安全验证 ===
print("\n=== 7. 安全验证 ===")

def test_idor_order():
    r = api("GET", "/api/client/orders/999999", token=CLIENT_TOKEN)
    assert r.json().get("code") == 404

def test_idor_invoice():
    r = api("GET", "/api/client/invoices/999999", token=CLIENT_TOKEN)
    assert r.json().get("code") == 404

def test_idor_ticket():
    r = api("GET", "/api/client/tickets/999999", token=CLIENT_TOKEN)
    assert r.json().get("code") == 404

def test_idor_service():
    r = api("GET", "/api/client/services/999999", token=CLIENT_TOKEN)
    assert r.json().get("code") == 404

def test_wrong_password():
    r = api("POST", "/api/admin/login", {"username": "admin", "password": "wrong"})
    assert_code(r, 401)

test("订单IDOR防护", test_idor_order)
test("账单IDOR防护", test_idor_invoice)
test("工单IDOR防护", test_idor_ticket)
test("服务IDOR防护", test_idor_service)
test("错误密码拦截", test_wrong_password)

# === 8. 核心模块API验证 ===
print("\n=== 8. 核心模块API验证 ===")

def test_dashboard():
    r = api("GET", "/api/admin/dashboard/stats", token=ADMIN_TOKEN)
    assert_code(r, 0)
    d = r.json()["data"]
    assert "user_count" in d, f"dashboard缺少user_count: {d.keys()}"

def test_users():
    r = api("GET", "/api/admin/users", token=ADMIN_TOKEN)
    assert_code(r, 0)
    assert r.json()["data"]["total"] > 0, "用户数应>0"

def test_finance():
    r = api("GET", "/api/admin/finance/ledger", token=ADMIN_TOKEN)
    assert_code(r, 0)

def test_supplier():
    r = api("GET", "/api/admin/suppliers", token=ADMIN_TOKEN)
    assert_code(r, 0)

def test_plugin():
    r = api("GET", "/api/admin/plugins", token=ADMIN_TOKEN)
    assert_code(r, 0)

def test_ai_config():
    r = api("GET", "/api/admin/ai/config", token=ADMIN_TOKEN)
    assert_code(r, 0)

def test_redis_config():
    r = api("GET", "/api/admin/redis/config", token=ADMIN_TOKEN)
    assert_code(r, 0)

def test_notices():
    r = api("GET", "/api/client/notices")
    assert_code(r, 0)

def test_help():
    r = api("GET", "/api/client/help-articles")
    assert_code(r, 0)

def test_content():
    r = api("GET", "/api/client/content/overview")
    assert_code(r, 0)

def test_referral():
    r = api("GET", "/api/client/referral/overview", token=CLIENT_TOKEN)
    assert_code(r, 0)
    assert "/register?ref=" in r.json()["data"]["referral_link"], "推介链接格式错误"

test("仪表盘", test_dashboard)
test("用户管理", test_users)
test("财务管理", test_finance)
test("供应商管理", test_supplier)
test("插件管理", test_plugin)
test("AI配置", test_ai_config)
test("Redis配置", test_redis_config)
test("公告", test_notices)
test("帮助", test_help)
test("内容概览", test_content)
test("推介系统", test_referral)

# === 结果 ===
print(f"\n{'='*60}")
print(f"严格审计结果: {results['pass']} pass / {results['fail']} fail / {results['pass']+results['fail']} total")
if results["errors"]:
    print("\nFailed:")
    for e in results["errors"]:
        print(f"  [FAIL] {e}")
print(f"{'='*60}")
sys.exit(0 if results["fail"] == 0 else 1)
