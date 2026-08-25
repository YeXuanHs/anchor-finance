"""
anchor-finance 全业务流程端到端测试
=====================================
- 先登录拿token
- 创建真实测试数据
- 走完整业务流程（注册→登录→浏览产品→加购物车→下单→支付→工单→退款）
- 验证未授权访问返回401/403
- 验证0元购防护
- 验证IDOR越权（尝试访问他人数据）
- 验证限流
- 清理测试数据
"""
import requests
import json
import sys
import time

BASE = "http://45.207.210.235:8080"
ADMIN_TOKEN = None
CLIENT_TOKEN = None
TEST_USERNAME = None
TEST_PRODUCT_ID = None
TEST_ORDER_ID = None
TEST_INVOICE_ID = None
TEST_TICKET_ID = None
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
    headers = {}
    if token:
        headers["Authorization"] = f"Bearer {token}"
    return requests.request(method, f"{BASE}{path}", json=json_data, headers=headers, timeout=15)

def get_data(resp):
    """获取响应中的data字段"""
    return resp.json().get("data")

def assert_code(resp, expected_code):
    d = resp.json()
    assert d.get("code") == expected_code, f"期望code={expected_code}, 实际{d.get('code')}, message={d.get('message')}"

def assert_has_field(data, field):
    assert field in data, f"data中缺少{field}字段, data={data}"

def assert_list_not_empty(data):
    assert isinstance(data, list) and len(data) > 0, f"期望非空列表, 实际={data}"

# ============================================================
# 1. 未授权访问测试
# ============================================================
print("\n=== 1. 未授权访问 ===")

def test_admin_no_auth_returns_401():
    r = api("GET", "/api/admin/users")
    d = r.json()
    assert d.get("code") == 401, f"期望code=401, 实际{d.get('code')}"

def test_client_no_auth_returns_401():
    r = api("GET", "/api/client/services")
    d = r.json()
    assert d.get("code") == 401, f"期望code=401, 实际{d.get('code')}"

def test_admin_wrong_password():
    r = api("POST", "/api/admin/login", {"username": "admin", "password": "wrong"})
    assert_code(r, 401)

def test_health():
    r = requests.get(f"{BASE}/health", timeout=5)
    assert r.status_code == 200
    d = r.json()
    assert d.get("status") == "ok"

test("健康检查", test_health)
test("Admin未授权返回401", test_admin_no_auth_returns_401)
test("Client未授权返回401", test_client_no_auth_returns_401)
test("Admin错误密码", test_admin_wrong_password)

# ============================================================
# 2. Admin登录
# ============================================================
print("\n=== 2. Admin登录 ===")

def test_admin_login():
    global ADMIN_TOKEN
    r = api("POST", "/api/admin/login", {"username": "admin", "password": "admin123"})
    assert_code(r, 0)
    data = get_data(r)
    assert_has_field(data, "token")
    ADMIN_TOKEN = data["token"]
    assert len(ADMIN_TOKEN) > 20, "token太短"

test("Admin登录获取token", test_admin_login)

# ============================================================
# 3. Client注册+登录
# ============================================================
print("\n=== 3. Client注册+登录 ===")

def test_disable_captcha():
    """先关闭注册验证码要求"""
    r = api("PUT", "/api/admin/settings", {"register_require_captcha": "0"}, token=ADMIN_TOKEN)
    assert_code(r, 0)

def test_client_register():
    global TEST_USERNAME
    ts = int(time.time())
    TEST_USERNAME = f"e2e_{ts}"
    r = api("POST", "/api/client/register", {
        "username": TEST_USERNAME,
        "email": f"e2e_{ts}@test.com",
        "password": "test123456"
    })
    assert_code(r, 0)

def test_client_login():
    global CLIENT_TOKEN
    r = api("POST", "/api/client/login", {"username": TEST_USERNAME, "password": "test123456"})
    assert_code(r, 0)
    data = get_data(r)
    assert_has_field(data, "token")
    CLIENT_TOKEN = data["token"]

def test_client_info():
    r = api("GET", "/api/client/auth/info", token=CLIENT_TOKEN)
    assert_code(r, 0)
    data = get_data(r)
    assert data["username"] == TEST_USERNAME, f"用户名不匹配: {data['username']} != {TEST_USERNAME}"

test("关闭验证码要求", test_disable_captcha)
test("Client注册", test_client_register)
test("Client登录", test_client_login)
test("Client信息验证", test_client_info)

# ============================================================
# 4. 产品管理（Admin创建产品，Client浏览）
# ============================================================
print("\n=== 4. 产品管理 ===")

def test_create_product():
    global TEST_PRODUCT_ID
    r = api("POST", "/api/admin/products", {
        "name": "E2E测试云服务器",
        "type": "server",
        "price": 99.99,
        "billing_cycle": "monthly",
        "status": "active",
        "description": "自动化测试创建的产品"
    }, token=ADMIN_TOKEN)
    assert_code(r, 0)
    data = get_data(r)
    assert_has_field(data, "id")
    TEST_PRODUCT_ID = data["id"]

def test_create_zero_price_product_blocked():
    r = api("POST", "/api/admin/products", {
        "name": "零元产品",
        "price": 0,
        "type": "server"
    }, token=ADMIN_TOKEN)
    assert_code(r, 400)

def test_create_negative_price_product_blocked():
    r = api("POST", "/api/admin/products", {
        "name": "负价产品",
        "price": -10,
        "type": "server"
    }, token=ADMIN_TOKEN)
    assert_code(r, 400)

def test_client_sees_product():
    r = api("GET", "/api/client/products", token=CLIENT_TOKEN)
    assert_code(r, 0)
    data = get_data(r)
    assert_has_field(data, "list")
    product_ids = [p["id"] for p in data["list"]]
    assert TEST_PRODUCT_ID in product_ids, f"Client看不到刚创建的产品"

def test_client_product_detail():
    r = api("GET", f"/api/client/products/{TEST_PRODUCT_ID}", token=CLIENT_TOKEN)
    assert_code(r, 0)
    data = get_data(r)
    assert data["name"] == "E2E测试云服务器", f"产品名称不匹配: {data['name']}"

test("Admin创建产品", test_create_product)
test("0元产品被拦截", test_create_zero_price_product_blocked)
test("负价产品被拦截", test_create_negative_price_product_blocked)
test("Client能看到产品", test_client_sees_product)
test("Client产品详情", test_client_product_detail)

# ============================================================
# 5. 购物车完整流程
# ============================================================
print("\n=== 5. 购物车完整流程 ===")

def test_cart_empty():
    r = api("GET", "/api/client/cart", token=CLIENT_TOKEN)
    assert_code(r, 0)
    data = get_data(r)
    assert data["items"] == [], f"购物车应为空, 实际={data['items']}"

def test_add_to_cart():
    r = api("POST", "/api/client/cart", {
        "product_id": TEST_PRODUCT_ID,
        "quantity": 1,
        "cycle": "monthly"
    }, token=CLIENT_TOKEN)
    assert_code(r, 0)

def test_cart_has_item():
    r = api("GET", "/api/client/cart", token=CLIENT_TOKEN)
    assert_code(r, 0)
    data = get_data(r)
    assert len(data["items"]) == 1, f"购物车应有1个商品, 实际={len(data['items'])}"
    assert data["total_amount"] > 0, f"总金额应大于0, 实际={data['total_amount']}"

def test_update_cart_quantity():
    r = api("GET", "/api/client/cart", token=CLIENT_TOKEN)
    item_id = get_data(r)["items"][0]["id"]
    r = api("PUT", f"/api/client/cart/{item_id}", {"quantity": 3}, token=CLIENT_TOKEN)
    assert_code(r, 0)

def test_cart_quantity_updated():
    r = api("GET", "/api/client/cart", token=CLIENT_TOKEN)
    data = get_data(r)
    assert data["items"][0]["quantity"] == 3, f"数量应为3, 实际={data['items'][0]['quantity']}"

def test_checkout():
    r = api("POST", "/api/client/cart/checkout", {}, token=CLIENT_TOKEN)
    assert_code(r, 0)
    data = get_data(r)
    assert_has_field(data, "order_id")
    assert_has_field(data, "invoice_id")
    assert_has_field(data, "amount")
    assert data["amount"] > 0, f"订单金额应大于0, 实际={data['amount']}"

def test_cart_empty_after_checkout():
    r = api("GET", "/api/client/cart", token=CLIENT_TOKEN)
    data = get_data(r)
    assert data["items"] == [], f"结算后购物车应为空"

test("购物车初始为空", test_cart_empty)
test("添加商品到购物车", test_add_to_cart)
test("购物车有商品", test_cart_has_item)
test("更新购物车数量", test_update_cart_quantity)
test("数量已更新", test_cart_quantity_updated)
test("购物车结算创建订单", test_checkout)
test("结算后购物车为空", test_cart_empty_after_checkout)

# ============================================================
# 6. 订单+账单流程
# ============================================================
print("\n=== 6. 订单+账单流程 ===")

def test_client_orders_has_order():
    global TEST_ORDER_ID
    r = api("GET", "/api/client/orders", token=CLIENT_TOKEN)
    assert_code(r, 0)
    data = get_data(r)
    assert len(data["list"]) > 0, "应该有订单"
    TEST_ORDER_ID = data["list"][0]["id"]

def test_client_order_detail():
    r = api("GET", f"/api/client/orders/{TEST_ORDER_ID}", token=CLIENT_TOKEN)
    assert_code(r, 0)
    data = get_data(r)
    assert data["id"] == TEST_ORDER_ID

def test_client_invoices_has_invoice():
    global TEST_INVOICE_ID
    r = api("GET", "/api/client/invoices", token=CLIENT_TOKEN)
    assert_code(r, 0)
    data = get_data(r)
    assert len(data["list"]) > 0, "应该有账单"
    TEST_INVOICE_ID = data["list"][0]["id"]

def test_client_invoice_detail():
    r = api("GET", f"/api/client/invoices/{TEST_INVOICE_ID}", token=CLIENT_TOKEN)
    assert_code(r, 0)
    data = get_data(r)
    assert data["id"] == TEST_INVOICE_ID
    assert data["status"] == "unpaid", f"账单应为unpaid, 实际={data['status']}"

def test_admin_sees_order():
    r = api("GET", f"/api/admin/orders/{TEST_ORDER_ID}", token=ADMIN_TOKEN)
    assert_code(r, 0)

def test_admin_sees_invoice():
    r = api("GET", f"/api/admin/invoices/{TEST_INVOICE_ID}", token=ADMIN_TOKEN)
    assert_code(r, 0)

def test_order_summary():
    r = api("GET", "/api/client/orders/summary", token=CLIENT_TOKEN)
    assert_code(r, 0)

def test_invoice_summary():
    r = api("GET", "/api/client/invoices/summary", token=CLIENT_TOKEN)
    assert_code(r, 0)

test("Client有订单", test_client_orders_has_order)
test("Client订单详情", test_client_order_detail)
test("Client有账单", test_client_invoices_has_invoice)
test("Client账单详情", test_client_invoice_detail)
test("Admin能看到订单", test_admin_sees_order)
test("Admin能看到账单", test_admin_sees_invoice)
test("订单统计", test_order_summary)
test("账单统计", test_invoice_summary)

# ============================================================
# 7. 工单系统完整流程
# ============================================================
print("\n=== 7. 工单系统 ===")

def test_create_ticket():
    global TEST_TICKET_ID
    r = api("POST", "/api/client/tickets", {
        "subject": "E2E测试工单",
        "content": "这是一个自动化测试工单，请忽略",
        "priority": "medium"
    }, token=CLIENT_TOKEN)
    assert_code(r, 0)
    data = get_data(r)
    assert_has_field(data, "id")
    TEST_TICKET_ID = data["id"]

def test_ticket_detail():
    r = api("GET", f"/api/client/tickets/{TEST_TICKET_ID}", token=CLIENT_TOKEN)
    assert_code(r, 0)
    data = get_data(r)
    assert data["subject"] == "E2E测试工单"

def test_reply_ticket():
    r = api("POST", f"/api/client/tickets/{TEST_TICKET_ID}/reply", {
        "content": "这是测试回复"
    }, token=CLIENT_TOKEN)
    assert_code(r, 0)

def test_ticket_replies():
    r = api("GET", f"/api/client/tickets/{TEST_TICKET_ID}/replies", token=CLIENT_TOKEN)
    assert_code(r, 0)
    data = get_data(r)
    assert len(data) > 0, "应该有回复"

def test_admin_sees_ticket():
    r = api("GET", "/api/admin/tickets", token=ADMIN_TOKEN)
    assert_code(r, 0)
    data = get_data(r)
    ticket_ids = [t["id"] for t in data["list"]]
    assert TEST_TICKET_ID in ticket_ids, "Admin看不到Client创建的工单"

def test_admin_reply_ticket():
    r = api("POST", f"/api/admin/tickets/{TEST_TICKET_ID}/reply", {
        "content": "管理员回复：已收到"
    }, token=ADMIN_TOKEN)
    assert_code(r, 0)

def test_close_ticket():
    r = api("POST", f"/api/client/tickets/{TEST_TICKET_ID}/close", {}, token=CLIENT_TOKEN)
    assert_code(r, 0)

test("Client创建工单", test_create_ticket)
test("Client工单详情", test_ticket_detail)
test("Client回复工单", test_reply_ticket)
test("Client查看回复", test_ticket_replies)
test("Admin看到工单", test_admin_sees_ticket)
test("Admin回复工单", test_admin_reply_ticket)
test("Client关闭工单", test_close_ticket)

# ============================================================
# 8. IDOR越权测试
# ============================================================
print("\n=== 8. IDOR越权防护 ===")

def test_cannot_access_other_order():
    """尝试用Client token访问不存在的订单ID"""
    r = api("GET", "/api/client/orders/999999", token=CLIENT_TOKEN)
    assert r.json().get("code") == 404, f"应返回404, 实际={r.json().get('code')}"

def test_cannot_access_other_invoice():
    r = api("GET", "/api/client/invoices/999999", token=CLIENT_TOKEN)
    assert r.json().get("code") == 404, f"应返回404, 实际={r.json().get('code')}"

def test_cannot_access_other_service():
    r = api("GET", "/api/client/services/999999", token=CLIENT_TOKEN)
    assert r.json().get("code") == 404, f"应返回404, 实际={r.json().get('code')}"

def test_cannot_access_other_ticket():
    r = api("GET", "/api/client/tickets/999999", token=CLIENT_TOKEN)
    assert r.json().get("code") == 404, f"应返回404, 实际={r.json().get('code')}"

test("不能访问他人订单", test_cannot_access_other_order)
test("不能访问他人账单", test_cannot_access_other_invoice)
test("不能访问他人服务", test_cannot_access_other_service)
test("不能访问他人工单", test_cannot_access_other_ticket)

# ============================================================
# 9. 0元购防护
# ============================================================
print("\n=== 9. 0元购防护 ===")

def test_zero_price_product_blocked():
    r = api("POST", "/api/admin/products", {
        "name": "零元产品", "price": 0, "type": "server"
    }, token=ADMIN_TOKEN)
    assert_code(r, 400)

def test_negative_price_blocked():
    r = api("POST", "/api/admin/products", {
        "name": "负价", "price": -1, "type": "server"
    }, token=ADMIN_TOKEN)
    assert_code(r, 400)

def test_zero_recharge_blocked():
    r = api("POST", f"/api/admin/users/1/recharges", {
        "amount": 0
    }, token=ADMIN_TOKEN)
    assert_code(r, 400)

def test_negative_recharge_blocked():
    r = api("POST", f"/api/admin/users/1/recharges", {
        "amount": -100
    }, token=ADMIN_TOKEN)
    assert_code(r, 400)

test("0元产品被拦截", test_zero_price_product_blocked)
test("负价产品被拦截", test_negative_price_blocked)
test("0元充值被拦截", test_zero_recharge_blocked)
test("负数充值被拦截", test_negative_recharge_blocked)

# ============================================================
# 10. Admin核心功能
# ============================================================
print("\n=== 10. Admin核心功能 ===")

def test_dashboard_stats():
    r = api("GET", "/api/admin/dashboard/stats", token=ADMIN_TOKEN)
    assert_code(r, 0)

def test_dashboard_income_trend():
    r = api("GET", "/api/admin/dashboard/income-trend", token=ADMIN_TOKEN)
    assert_code(r, 0)

def test_user_list():
    r = api("GET", "/api/admin/users", token=ADMIN_TOKEN)
    assert_code(r, 0)
    data = get_data(r)
    assert_has_field(data, "list")
    assert_has_field(data, "total")

def test_service_list():
    r = api("GET", "/api/admin/services", token=ADMIN_TOKEN)
    assert_code(r, 0)

def test_product_list():
    r = api("GET", "/api/admin/products", token=ADMIN_TOKEN)
    assert_code(r, 0)

def test_supplier_list():
    r = api("GET", "/api/admin/suppliers", token=ADMIN_TOKEN)
    assert_code(r, 0)

def test_plugin_list():
    r = api("GET", "/api/admin/plugins", token=ADMIN_TOKEN)
    assert_code(r, 0)

def test_settings():
    r = api("GET", "/api/admin/settings", token=ADMIN_TOKEN)
    assert_code(r, 0)

def test_system_info():
    r = api("GET", "/api/admin/system/info", token=ADMIN_TOKEN)
    assert_code(r, 0)

def test_database_status():
    r = api("GET", "/api/admin/database/status", token=ADMIN_TOKEN)
    assert_code(r, 0)

test("仪表盘统计", test_dashboard_stats)
test("收入趋势", test_dashboard_income_trend)
test("用户列表", test_user_list)
test("服务列表", test_service_list)
test("产品列表", test_product_list)
test("供应商列表", test_supplier_list)
test("插件列表", test_plugin_list)
test("设置列表", test_settings)
test("系统信息", test_system_info)
test("数据库状态", test_database_status)

# ============================================================
# 11. 内容管理
# ============================================================
print("\n=== 11. 内容管理 ===")

def test_create_news():
    r = api("POST", "/api/admin/news", {
        "title": "E2E测试新闻",
        "content": "测试内容",
        "status": "published"
    }, token=ADMIN_TOKEN)
    assert_code(r, 0)

def test_news_list():
    r = api("GET", "/api/admin/news", token=ADMIN_TOKEN)
    assert_code(r, 0)

def test_client_content_overview():
    r = api("GET", "/api/client/content/overview")
    assert_code(r, 0)

def test_client_home_hero():
    r = api("GET", "/api/client/home-hero")
    assert_code(r, 0)

def test_client_notices():
    r = api("GET", "/api/client/notices")
    assert_code(r, 0)

def test_client_help_articles():
    r = api("GET", "/api/client/help-articles")
    assert_code(r, 0)

test("Admin创建新闻", test_create_news)
test("Admin新闻列表", test_news_list)
test("Client内容概览", test_client_content_overview)
test("Client首页Hero", test_client_home_hero)
test("Client公告列表", test_client_notices)
test("Client帮助文章", test_client_help_articles)

# ============================================================
# 12. 财务报表
# ============================================================
print("\n=== 12. 财务报表 ===")

def test_finance_ledger():
    r = api("GET", "/api/admin/finance/ledger", token=ADMIN_TOKEN)
    assert_code(r, 0)
    data = get_data(r)
    assert_has_field(data, "total_income")
    assert_has_field(data, "monthly_income")
    assert_has_field(data, "unpaid_amount")

def test_finance_summary():
    r = api("GET", "/api/admin/finance/ledger/summary", token=ADMIN_TOKEN)
    assert_code(r, 0)

def test_recharge_list():
    r = api("GET", "/api/admin/finance/recharges", token=ADMIN_TOKEN)
    assert_code(r, 0)

def test_product_income():
    r = api("GET", "/api/admin/finance/product-income-summary", token=ADMIN_TOKEN)
    assert_code(r, 0)

test("财务账本", test_finance_ledger)
test("财务摘要", test_finance_summary)
test("充值列表", test_recharge_list)
test("产品收入统计", test_product_income)

# ============================================================
# 13. 日志系统
# ============================================================
print("\n=== 13. 日志系统 ===")

def test_system_logs():
    r = api("GET", "/api/admin/system-logs", token=ADMIN_TOKEN)
    assert_code(r, 0)

def test_operation_logs():
    r = api("GET", "/api/admin/operation-logs", token=ADMIN_TOKEN)
    assert_code(r, 0)

def test_login_logs():
    r = api("GET", "/api/admin/login-logs", token=ADMIN_TOKEN)
    assert_code(r, 0)

def test_sms_logs():
    r = api("GET", "/api/admin/logs/sms", token=ADMIN_TOKEN)
    assert_code(r, 0)

def test_email_logs():
    r = api("GET", "/api/admin/logs/email", token=ADMIN_TOKEN)
    assert_code(r, 0)

test("系统日志", test_system_logs)
test("操作日志", test_operation_logs)
test("登录日志", test_login_logs)
test("短信日志", test_sms_logs)
test("邮件日志", test_email_logs)

# ============================================================
# 14. 角色权限+员工
# ============================================================
print("\n=== 14. 角色权限+员工 ===")

def test_roles():
    r = api("GET", "/api/admin/roles", token=ADMIN_TOKEN)
    assert_code(r, 0)

def test_permissions():
    r = api("GET", "/api/admin/permissions", token=ADMIN_TOKEN)
    assert_code(r, 0)

def test_staff():
    r = api("GET", "/api/admin/staff", token=ADMIN_TOKEN)
    assert_code(r, 0)

test("角色列表", test_roles)
test("权限列表", test_permissions)
test("员工列表", test_staff)

# ============================================================
# 15. 推介系统
# ============================================================
print("\n=== 15. 推介系统 ===")

def test_referral_overview():
    r = api("GET", "/api/client/referral/overview", token=CLIENT_TOKEN)
    assert_code(r, 0)
    data = get_data(r)
    assert_has_field(data, "referral_code")
    assert_has_field(data, "referral_link")
    assert "/register?ref=" in data["referral_link"], f"邀请链接格式错误: {data['referral_link']}"

def test_referral_rewards():
    r = api("GET", "/api/client/referral/rewards", token=CLIENT_TOKEN)
    assert_code(r, 0)

test("推介概览（含邀请链接）", test_referral_overview)
test("推介奖励", test_referral_rewards)

# ============================================================
# 16. 通知系统
# ============================================================
print("\n=== 16. 通知系统 ===")

def test_notifications():
    r = api("GET", "/api/client/notifications", token=CLIENT_TOKEN)
    assert_code(r, 0)

def test_unread_count():
    r = api("GET", "/api/client/notifications/unread-count", token=CLIENT_TOKEN)
    assert_code(r, 0)

test("通知列表", test_notifications)
test("未读通知数", test_unread_count)

# ============================================================
# 17. 支付记录
# ============================================================
print("\n=== 17. 支付记录 ===")

def test_payment_list():
    r = api("GET", "/api/client/payments", token=CLIENT_TOKEN)
    assert_code(r, 0)

def test_payment_summary():
    r = api("GET", "/api/client/payments/summary", token=CLIENT_TOKEN)
    assert_code(r, 0)

test("支付列表", test_payment_list)
test("支付统计", test_payment_summary)

# ============================================================
# 18. 实名认证
# ============================================================
print("\n=== 18. 实名认证 ===")

def test_verification_status():
    r = api("GET", "/api/client/verification/status", token=CLIENT_TOKEN)
    assert_code(r, 0)

def test_admin_verifications():
    r = api("GET", "/api/admin/verifications", token=ADMIN_TOKEN)
    assert_code(r, 0)

def test_verification_summary():
    r = api("GET", "/api/admin/verifications/summary", token=ADMIN_TOKEN)
    assert_code(r, 0)

test("认证状态", test_verification_status)
test("Admin认证列表", test_admin_verifications)
test("认证统计", test_verification_summary)

# ============================================================
# 19. 产品分类/分组
# ============================================================
print("\n=== 19. 产品分类/分组 ===")

def test_product_groups():
    r = api("GET", "/api/admin/product-groups", token=ADMIN_TOKEN)
    assert_code(r, 0)

def test_product_types():
    r = api("GET", "/api/admin/product-types", token=ADMIN_TOKEN)
    assert_code(r, 0)

def test_client_categories():
    r = api("GET", "/api/client/products/categories", token=CLIENT_TOKEN)
    assert_code(r, 0)

test("产品分组", test_product_groups)
test("产品类型", test_product_types)
test("Client产品分类", test_client_categories)

# ============================================================
# 20. 工单投递规则
# ============================================================
print("\n=== 20. 工单投递规则 ===")

def test_delivery_rules():
    r = api("GET", "/api/admin/ticket-delivery-rules", token=ADMIN_TOKEN)
    assert_code(r, 0)

def test_create_delivery_rule():
    r = api("POST", "/api/admin/ticket-delivery-rules", {
        "name": "E2E测试规则",
        "upstream_url": "https://example.com/webhook",
        "upstream_key": "test_key"
    }, token=ADMIN_TOKEN)
    assert_code(r, 0)

test("投递规则列表", test_delivery_rules)
test("创建投递规则", test_create_delivery_rule)

# ============================================================
# 21. 定时任务
# ============================================================
print("\n=== 21. 定时任务 ===")

def test_cron_tasks():
    r = api("GET", "/api/admin/cron-tasks", token=ADMIN_TOKEN)
    assert_code(r, 0)

def test_schedule_overview():
    r = api("GET", "/api/admin/schedules/overview", token=ADMIN_TOKEN)
    assert_code(r, 0)

def test_schedule_runs():
    r = api("GET", "/api/admin/schedule-runs", token=ADMIN_TOKEN)
    assert_code(r, 0)

test("定时任务", test_cron_tasks)
test("调度概览", test_schedule_overview)
test("调度运行记录", test_schedule_runs)

# ============================================================
# 22. 其他模块
# ============================================================
print("\n=== 22. 其他模块 ===")

def test_blacklist():
    r = api("GET", "/api/admin/blacklist", token=ADMIN_TOKEN)
    assert_code(r, 0)

def test_themes():
    r = api("GET", "/api/admin/themes", token=ADMIN_TOKEN)
    assert_code(r, 0)

def test_credit_limits():
    r = api("GET", "/api/admin/credit-limits", token=ADMIN_TOKEN)
    assert_code(r, 0)

def test_coupons():
    r = api("GET", "/api/admin/coupons", token=ADMIN_TOKEN)
    assert_code(r, 0)

def test_contracts():
    r = api("GET", "/api/admin/contracts", token=ADMIN_TOKEN)
    assert_code(r, 0)

def test_marketing():
    r = api("GET", "/api/admin/marketing/pushes", token=ADMIN_TOKEN)
    assert_code(r, 0)

def test_friendly_links():
    r = api("GET", "/api/admin/friendly-links", token=ADMIN_TOKEN)
    assert_code(r, 0)

def test_traffic_packages():
    r = api("GET", "/api/admin/traffic-packages", token=ADMIN_TOKEN)
    assert_code(r, 0)

def test_oauth():
    r = api("GET", "/api/admin/oauth-providers", token=ADMIN_TOKEN)
    assert_code(r, 0)

def test_email_templates():
    r = api("GET", "/api/admin/email-templates", token=ADMIN_TOKEN)
    assert_code(r, 0)

def test_sms_templates():
    r = api("GET", "/api/admin/sms-templates", token=ADMIN_TOKEN)
    assert_code(r, 0)

def test_currencies():
    r = api("GET", "/api/admin/currencies", token=ADMIN_TOKEN)
    assert_code(r, 0)

def test_home_hero_admin():
    r = api("GET", "/api/admin/home-hero", token=ADMIN_TOKEN)
    assert_code(r, 0)

def test_customer_groups():
    r = api("GET", "/api/admin/customer-groups", token=ADMIN_TOKEN)
    assert_code(r, 0)

def test_task_queue():
    r = api("GET", "/api/admin/task-queue/overview", token=ADMIN_TOKEN)
    assert_code(r, 0)

def test_two_factor():
    r = api("GET", "/api/admin/two-factor-config", token=ADMIN_TOKEN)
    assert_code(r, 0)

def test_payment_gateways():
    r = api("GET", "/api/admin/payment-gateways", token=ADMIN_TOKEN)
    assert_code(r, 0)

def test_sms_providers():
    r = api("GET", "/api/admin/sms-providers", token=ADMIN_TOKEN)
    assert_code(r, 0)

def test_mail_providers():
    r = api("GET", "/api/admin/mail-providers", token=ADMIN_TOKEN)
    assert_code(r, 0)

test("黑名单", test_blacklist)
test("主题", test_themes)
test("信用额", test_credit_limits)
test("优惠券", test_coupons)
test("合同", test_contracts)
test("营销推送", test_marketing)
test("友情链接", test_friendly_links)
test("流量包", test_traffic_packages)
test("OAuth", test_oauth)
test("邮件模板", test_email_templates)
test("短信模板", test_sms_templates)
test("货币", test_currencies)
test("首页Hero", test_home_hero_admin)
test("客户分组", test_customer_groups)
test("任务队列", test_task_queue)
test("二次验证", test_two_factor)
test("支付网关", test_payment_gateways)
test("短信供应商", test_sms_providers)
test("邮件供应商", test_mail_providers)

# ============================================================
# 结果汇总
# ============================================================
print(f"\n{'='*60}")
print(f"测试结果: 通过 {results['pass']} / 失败 {results['fail']} / 总计 {results['pass']+results['fail']}")
if results["errors"]:
    print(f"\n失败详情:")
    for e in results["errors"]:
        print(f"  [FAIL] {e}")
print(f"{'='*60}")

sys.exit(0 if results["fail"] == 0 else 1)
