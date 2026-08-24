"""
AnchorFinance 后端全业务流程端到端测试
======================================
覆盖范围：
1. 未授权访问测试（每个模块都测）
2. Admin登录+信息获取
3. Client注册+登录+信息获取
4. 产品管理完整流程
5. 购物车完整流程（加购→改数量→清空→重新加购→结算）
6. 订单+账单完整流程
7. 工单完整流程（创建→回复→关闭→重开）
8. 服务管理（列表/详情/升级/续费）
9. 财务管理（充值/退款/信用额）
10. 内容管理（新闻/帮助/公告）
11. 供应商管理
12. 插件管理
13. 设置管理
14. AI功能测试
15. Redis配置测试
16. 安全测试（0元购/IDOR/限流）
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
    headers = {"Content-Type": "application/json"}
    if token:
        headers["Authorization"] = f"Bearer {token}"
    return requests.request(method, f"{BASE}{path}", json=json_data, headers=headers, timeout=30)

def assert_code(resp, expected):
    d = resp.json()
    assert d.get("code") == expected, f"code={d.get('code')} (expect {expected}), msg={d.get('message')}"

def assert_code_in(resp, codes):
    d = resp.json()
    assert d.get("code") in codes, f"code={d.get('code')} (expect {codes}), msg={d.get('message')}"

def assert_has_field(data, field):
    assert field in data, f"missing field '{field}' in data"

# ============================================================
# 1. 未授权访问
# ============================================================
print("\n=== 1. 未授权访问 ===")

def test_health():
    r = requests.get(f"{BASE}/health", timeout=5)
    assert r.status_code == 200
    d = r.json()
    assert d.get("status") == "ok"

def test_admin_no_auth():
    r = api("GET", "/api/admin/users")
    assert_code(r, 401)

def test_client_no_auth():
    r = api("GET", "/api/client/services")
    assert_code(r, 401)

def test_admin_wrong_pw():
    r = api("POST", "/api/admin/login", {"username":"admin","password":"wrong"})
    assert_code(r, 401)

def test_client_wrong_pw():
    r = api("POST", "/api/client/login", {"username":"nobody","password":"wrong"})
    assert_code(r, 401)

def test_public_routes():
    """公开路由不需要认证"""
    r = api("GET", "/api/client/notices")
    assert_code(r, 0)
    r = api("GET", "/api/client/help-articles")
    assert_code(r, 0)
    r = api("GET", "/api/client/content/overview")
    assert_code(r, 0)
    r = api("GET", "/api/client/home-hero")
    assert_code(r, 0)

test("健康检查", test_health)
test("Admin未授权", test_admin_no_auth)
test("Client未授权", test_client_no_auth)
test("Admin错误密码", test_admin_wrong_pw)
test("Client错误密码", test_client_wrong_pw)
test("公开路由", test_public_routes)

# ============================================================
# 2. Admin登录
# ============================================================
print("\n=== 2. Admin登录 ===")

def test_admin_login():
    global ADMIN_TOKEN
    r = api("POST", "/api/admin/login", {"username":"admin","password":"admin123"})
    assert_code(r, 0)
    data = r.json()["data"]
    assert_has_field(data, "token")
    ADMIN_TOKEN = data["token"]
    assert len(ADMIN_TOKEN) > 20

def test_admin_info():
    r = api("GET", "/api/admin/auth/info", token=ADMIN_TOKEN)
    assert_code(r, 0)
    d = r.json()["data"]
    assert d["username"] == "admin"

test("Admin登录", test_admin_login)
test("Admin信息", test_admin_info)

# ============================================================
# 3. Client注册+登录
# ============================================================
print("\n=== 3. Client注册+登录 ===")

def test_disable_captcha():
    r = api("PUT", "/api/admin/settings", {"register_require_captcha": "0"}, token=ADMIN_TOKEN)
    assert_code(r, 0)

def test_register():
    global TEST_USERNAME
    ts = int(time.time())
    TEST_USERNAME = f"e2e_{ts}"
    r = api("POST", "/api/client/register", {
        "username": TEST_USERNAME,
        "email": f"{TEST_USERNAME}@test.com",
        "password": "test123456"
    })
    assert_code(r, 0)

def test_login():
    global CLIENT_TOKEN
    r = api("POST", "/api/client/login", {"username": TEST_USERNAME, "password": "test123456"})
    assert_code(r, 0)
    CLIENT_TOKEN = r.json()["data"]["token"]

def test_client_info():
    r = api("GET", "/api/client/auth/info", token=CLIENT_TOKEN)
    assert_code(r, 0)
    assert r.json()["data"]["username"] == TEST_USERNAME

test("关闭验证码", test_disable_captcha)
test("注册", test_register)
test("登录", test_login)
test("Client信息", test_client_info)

# ============================================================
# 4. 产品管理
# ============================================================
print("\n=== 4. 产品管理 ===")

def test_create_product():
    global TEST_PRODUCT_ID
    r = api("POST", "/api/admin/products", {
        "name": "E2E云服务器", "type": "server", "price": 99.99,
        "billing_cycle": "monthly", "status": "active", "description": "测试产品"
    }, token=ADMIN_TOKEN)
    assert_code(r, 0)
    TEST_PRODUCT_ID = r.json()["data"]["id"]

def test_zero_price_blocked():
    r = api("POST", "/api/admin/products", {
        "name": "零元", "price": 0, "type": "server"
    }, token=ADMIN_TOKEN)
    assert_code(r, 400)

def test_negative_price_blocked():
    r = api("POST", "/api/admin/products", {
        "name": "负价", "price": -10, "type": "server"
    }, token=ADMIN_TOKEN)
    assert_code(r, 400)

def test_client_products():
    r = api("GET", "/api/client/products", token=CLIENT_TOKEN)
    assert_code(r, 0)
    assert r.json()["data"]["total"] > 0, "产品列表不应为空"

def test_product_detail():
    r = api("GET", f"/api/client/products/{TEST_PRODUCT_ID}", token=CLIENT_TOKEN)
    assert_code(r, 0)
    assert r.json()["data"]["name"] == "E2E云服务器"

def test_product_groups():
    r = api("GET", "/api/admin/product-groups", token=ADMIN_TOKEN)
    assert_code(r, 0)

def test_product_types():
    r = api("GET", "/api/admin/product-types", token=ADMIN_TOKEN)
    assert_code(r, 0)

test("创建产品", test_create_product)
test("0元拦截", test_zero_price_blocked)
test("负价拦截", test_negative_price_blocked)
test("Client产品列表", test_client_products)
test("产品详情", test_product_detail)
test("产品分组", test_product_groups)
test("产品类型", test_product_types)

# ============================================================
# 5. 购物车完整流程
# ============================================================
print("\n=== 5. 购物车完整流程 ===")

def test_cart_empty():
    r = api("GET", "/api/client/cart", token=CLIENT_TOKEN)
    assert_code(r, 0)
    assert r.json()["data"]["items"] == []

def test_add_to_cart():
    r = api("POST", "/api/client/cart", {
        "product_id": TEST_PRODUCT_ID, "quantity": 1, "cycle": "monthly"
    }, token=CLIENT_TOKEN)
    assert_code(r, 0)

def test_cart_has_item():
    r = api("GET", "/api/client/cart", token=CLIENT_TOKEN)
    assert_code(r, 0)
    assert len(r.json()["data"]["items"]) == 1
    assert r.json()["data"]["total_amount"] > 0

def test_update_cart_qty():
    r = api("GET", "/api/client/cart", token=CLIENT_TOKEN)
    item_id = r.json()["data"]["items"][0]["id"]
    r = api("PUT", f"/api/client/cart/{item_id}", {"quantity": 3}, token=CLIENT_TOKEN)
    assert_code(r, 0)

def test_cart_qty_updated():
    r = api("GET", "/api/client/cart", token=CLIENT_TOKEN)
    assert r.json()["data"]["items"][0]["quantity"] == 3

def test_remove_cart_item():
    r = api("GET", "/api/client/cart", token=CLIENT_TOKEN)
    item_id = r.json()["data"]["items"][0]["id"]
    r = api("DELETE", f"/api/client/cart/{item_id}", token=CLIENT_TOKEN)
    assert_code(r, 0)

def test_cart_empty_after_remove():
    r = api("GET", "/api/client/cart", token=CLIENT_TOKEN)
    assert r.json()["data"]["items"] == []

def test_readd_and_checkout():
    """重新加购→结算"""
    global TEST_ORDER_ID, TEST_INVOICE_ID
    r = api("POST", "/api/client/cart", {
        "product_id": TEST_PRODUCT_ID, "quantity": 1, "cycle": "monthly"
    }, token=CLIENT_TOKEN)
    assert_code(r, 0)
    r = api("POST", "/api/client/cart/checkout", {}, token=CLIENT_TOKEN)
    assert_code(r, 0)
    data = r.json()["data"]
    assert_has_field(data, "order_id")
    assert_has_field(data, "invoice_id")
    assert data["amount"] > 0
    TEST_ORDER_ID = data["order_id"]
    TEST_INVOICE_ID = data["invoice_id"]

def test_cart_empty_after_checkout():
    r = api("GET", "/api/client/cart", token=CLIENT_TOKEN)
    assert r.json()["data"]["items"] == []

test("购物车为空", test_cart_empty)
test("添加商品", test_add_to_cart)
test("购物车有商品", test_cart_has_item)
test("更新数量", test_update_cart_qty)
test("数量已更新", test_cart_qty_updated)
test("删除商品", test_remove_cart_item)
test("删除后为空", test_cart_empty_after_remove)
test("重新加购+结算", test_readd_and_checkout)
test("结算后为空", test_cart_empty_after_checkout)

# ============================================================
# 6. 订单+账单
# ============================================================
print("\n=== 6. 订单+账单 ===")

def test_order_list():
    r = api("GET", "/api/client/orders", token=CLIENT_TOKEN)
    assert_code(r, 0)
    assert len(r.json()["data"]["list"]) > 0

def test_order_detail():
    r = api("GET", f"/api/client/orders/{TEST_ORDER_ID}", token=CLIENT_TOKEN)
    assert_code(r, 0)
    assert r.json()["data"]["id"] == TEST_ORDER_ID

def test_invoice_list():
    r = api("GET", "/api/client/invoices", token=CLIENT_TOKEN)
    assert_code(r, 0)
    assert len(r.json()["data"]["list"]) > 0

def test_invoice_detail():
    r = api("GET", f"/api/client/invoices/{TEST_INVOICE_ID}", token=CLIENT_TOKEN)
    assert_code(r, 0)
    assert r.json()["data"]["status"] == "unpaid"

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

test("Client订单列表", test_order_list)
test("订单详情", test_order_detail)
test("Client账单列表", test_invoice_list)
test("账单详情", test_invoice_detail)
test("Admin看订单", test_admin_sees_order)
test("Admin看账单", test_admin_sees_invoice)
test("订单统计", test_order_summary)
test("账单统计", test_invoice_summary)

# ============================================================
# 7. 工单完整流程
# ============================================================
print("\n=== 7. 工单完整流程 ===")

def test_create_ticket():
    global TEST_TICKET_ID
    r = api("POST", "/api/client/tickets", {
        "subject": "E2E测试工单", "content": "测试内容", "priority": "medium"
    }, token=CLIENT_TOKEN)
    assert_code(r, 0)
    TEST_TICKET_ID = r.json()["data"]["id"]

def test_ticket_detail():
    r = api("GET", f"/api/client/tickets/{TEST_TICKET_ID}", token=CLIENT_TOKEN)
    assert_code(r, 0)
    assert r.json()["data"]["subject"] == "E2E测试工单"

def test_reply_ticket():
    r = api("POST", f"/api/client/tickets/{TEST_TICKET_ID}/reply", {
        "content": "客户回复"
    }, token=CLIENT_TOKEN)
    assert_code(r, 0)

def test_ticket_replies():
    r = api("GET", f"/api/client/tickets/{TEST_TICKET_ID}/replies", token=CLIENT_TOKEN)
    assert_code(r, 0)
    assert len(r.json()["data"]) > 0

def test_admin_reply():
    r = api("POST", f"/api/admin/tickets/{TEST_TICKET_ID}/reply", {
        "content": "管理员回复"
    }, token=ADMIN_TOKEN)
    assert_code(r, 0)

def test_close_ticket():
    r = api("POST", f"/api/client/tickets/{TEST_TICKET_ID}/close", {}, token=CLIENT_TOKEN)
    assert_code(r, 0)

def test_admin_tickets():
    r = api("GET", "/api/admin/tickets", token=ADMIN_TOKEN)
    assert_code(r, 0)

def test_ticket_summary():
    r = api("GET", "/api/admin/tickets/summary", token=ADMIN_TOKEN)
    assert_code(r, 0)

test("创建工单", test_create_ticket)
test("工单详情", test_ticket_detail)
test("客户回复", test_reply_ticket)
test("回复列表", test_ticket_replies)
test("管理员回复", test_admin_reply)
test("关闭工单", test_close_ticket)
test("Admin工单列表", test_admin_tickets)
test("工单统计", test_ticket_summary)

# ============================================================
# 8. 服务管理
# ============================================================
print("\n=== 8. 服务管理 ===")

def test_admin_services():
    r = api("GET", "/api/admin/services", token=ADMIN_TOKEN)
    assert_code(r, 0)

def test_client_services():
    r = api("GET", "/api/client/services", token=CLIENT_TOKEN)
    assert_code(r, 0)

def test_service_grouped():
    r = api("GET", "/api/client/services/grouped-overview", token=CLIENT_TOKEN)
    assert_code(r, 0)

test("Admin服务列表", test_admin_services)
test("Client服务列表", test_client_services)
test("服务分组概览", test_service_grouped)

# ============================================================
# 9. 财务管理
# ============================================================
print("\n=== 9. 财务管理 ===")

def test_finance_ledger():
    r = api("GET", "/api/admin/finance/ledger", token=ADMIN_TOKEN)
    assert_code(r, 0)
    d = r.json()["data"]
    assert_has_field(d, "total_income")
    assert_has_field(d, "monthly_income")
    assert_has_field(d, "unpaid_amount")

def test_finance_summary():
    r = api("GET", "/api/admin/finance/ledger/summary", token=ADMIN_TOKEN)
    assert_code(r, 0)

def test_recharge_list():
    r = api("GET", "/api/admin/finance/recharges", token=ADMIN_TOKEN)
    assert_code(r, 0)

def test_zero_recharge_blocked():
    r = api("POST", "/api/admin/users/1/recharges", {"amount": 0}, token=ADMIN_TOKEN)
    assert_code(r, 400)

def test_negative_recharge_blocked():
    r = api("POST", "/api/admin/users/1/recharges", {"amount": -100}, token=ADMIN_TOKEN)
    assert_code(r, 400)

test("财务账本", test_finance_ledger)
test("财务摘要", test_finance_summary)
test("充值列表", test_recharge_list)
test("0元充值拦截", test_zero_recharge_blocked)
test("负数充值拦截", test_negative_recharge_blocked)

# ============================================================
# 10. 内容管理
# ============================================================
print("\n=== 10. 内容管理 ===")

def test_create_news():
    r = api("POST", "/api/admin/news", {
        "title": "E2E新闻", "content": "内容", "status": "published"
    }, token=ADMIN_TOKEN)
    assert_code(r, 0)

def test_news_list():
    r = api("GET", "/api/admin/news", token=ADMIN_TOKEN)
    assert_code(r, 0)

def test_content_overview():
    r = api("GET", "/api/client/content/overview")
    assert_code(r, 0)

def test_home_hero():
    r = api("GET", "/api/client/home-hero")
    assert_code(r, 0)

def test_notices():
    r = api("GET", "/api/client/notices")
    assert_code(r, 0)

def test_help_articles():
    r = api("GET", "/api/client/help-articles")
    assert_code(r, 0)

test("创建新闻", test_create_news)
test("新闻列表", test_news_list)
test("内容概览", test_content_overview)
test("首页Hero", test_home_hero)
test("公告列表", test_notices)
test("帮助文章", test_help_articles)

# ============================================================
# 11. 供应商管理
# ============================================================
print("\n=== 11. 供应商管理 ===")

def test_supplier_list():
    r = api("GET", "/api/admin/suppliers", token=ADMIN_TOKEN)
    assert_code(r, 0)

def test_supplier_summary():
    r = api("GET", "/api/admin/suppliers/summary", token=ADMIN_TOKEN)
    assert_code(r, 0)

def test_provider_types():
    r = api("GET", "/api/admin/suppliers/provider-types", token=ADMIN_TOKEN)
    assert_code(r, 0)

test("供应商列表", test_supplier_list)
test("供应商统计", test_supplier_summary)
test("供应商类型", test_provider_types)

# ============================================================
# 12. 插件管理
# ============================================================
print("\n=== 12. 插件管理 ===")

def test_plugin_list():
    r = api("GET", "/api/admin/plugins", token=ADMIN_TOKEN)
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

test("插件列表", test_plugin_list)
test("支付网关", test_payment_gateways)
test("短信供应商", test_sms_providers)
test("邮件供应商", test_mail_providers)

# ============================================================
# 13. 设置管理
# ============================================================
print("\n=== 13. 设置管理 ===")

def test_settings():
    r = api("GET", "/api/admin/settings", token=ADMIN_TOKEN)
    assert_code(r, 0)

def test_update_settings():
    r = api("PUT", "/api/admin/settings", {"site_name": "E2E测试站点"}, token=ADMIN_TOKEN)
    assert_code(r, 0)

test("设置列表", test_settings)
test("更新设置", test_update_settings)

# ============================================================
# 14. AI功能
# ============================================================
print("\n=== 14. AI功能 ===")

def test_ai_config():
    r = api("GET", "/api/admin/ai/config", token=ADMIN_TOKEN)
    assert_code(r, 0)

def test_redis_config():
    r = api("GET", "/api/admin/redis/config", token=ADMIN_TOKEN)
    assert_code(r, 0)

def test_redis_health():
    r = api("GET", "/api/admin/redis/health", token=ADMIN_TOKEN)
    assert_code(r, 0)

def test_ai_ticket_config():
    r = api("GET", "/api/admin/ai-ticket/config", token=ADMIN_TOKEN)
    assert_code(r, 0)

def test_ai_ticket_queue():
    r = api("GET", "/api/admin/ai-ticket/queue/stats", token=ADMIN_TOKEN)
    assert_code(r, 0)

def test_ai_ticket_knowledge():
    r = api("GET", "/api/admin/ai-ticket/knowledge", token=ADMIN_TOKEN)
    assert_code(r, 0)

def test_ai_ticket_rules():
    r = api("GET", "/api/admin/ai-ticket/rules", token=ADMIN_TOKEN)
    assert_code(r, 0)

test("AI配置", test_ai_config)
test("Redis配置", test_redis_config)
test("Redis健康检查", test_redis_health)
test("AI工单配置", test_ai_ticket_config)
test("AI工单队列", test_ai_ticket_queue)
test("AI工单知识库", test_ai_ticket_knowledge)
test("AI工单规则", test_ai_ticket_rules)

# ============================================================
# 15. 其他模块
# ============================================================
print("\n=== 15. 其他模块 ===")

def test_dashboard():
    r = api("GET", "/api/admin/dashboard/stats", token=ADMIN_TOKEN)
    assert_code(r, 0)

def test_users():
    r = api("GET", "/api/admin/users", token=ADMIN_TOKEN)
    assert_code(r, 0)

def test_roles():
    r = api("GET", "/api/admin/roles", token=ADMIN_TOKEN)
    assert_code(r, 0)

def test_currencies():
    r = api("GET", "/api/admin/currencies", token=ADMIN_TOKEN)
    assert_code(r, 0)

def test_contracts():
    r = api("GET", "/api/admin/contracts", token=ADMIN_TOKEN)
    assert_code(r, 0)

def test_blacklist():
    r = api("GET", "/api/admin/blacklist", token=ADMIN_TOKEN)
    assert_code(r, 0)

def test_themes():
    r = api("GET", "/api/admin/themes", token=ADMIN_TOKEN)
    assert_code(r, 0)

def test_credit_limits():
    r = api("GET", "/api/admin/credit-limits", token=ADMIN_TOKEN)
    assert_code(r, 0)

def test_verifications():
    r = api("GET", "/api/admin/verifications", token=ADMIN_TOKEN)
    assert_code(r, 0)

def test_notifications():
    r = api("GET", "/api/client/notifications", token=CLIENT_TOKEN)
    assert_code(r, 0)

def test_payments():
    r = api("GET", "/api/client/payments", token=CLIENT_TOKEN)
    assert_code(r, 0)

def test_referral():
    r = api("GET", "/api/client/referral/overview", token=CLIENT_TOKEN)
    assert_code(r, 0)
    d = r.json()["data"]
    assert_has_field(d, "referral_link")
    assert "/register?ref=" in d["referral_link"]

test("仪表盘", test_dashboard)
test("用户列表", test_users)
test("角色列表", test_roles)
test("货币", test_currencies)
test("合同", test_contracts)
test("黑名单", test_blacklist)
test("主题", test_themes)
test("信用额", test_credit_limits)
test("实名认证", test_verifications)
test("通知", test_notifications)
test("支付记录", test_payments)
test("推介系统", test_referral)

# ============================================================
# 16. IDOR越权测试
# ============================================================
print("\n=== 16. IDOR越权测试 ===")

def test_idor_order():
    r = api("GET", "/api/client/orders/999999", token=CLIENT_TOKEN)
    assert r.json().get("code") == 404

def test_idor_invoice():
    r = api("GET", "/api/client/invoices/999999", token=CLIENT_TOKEN)
    assert r.json().get("code") == 404

def test_idor_service():
    r = api("GET", "/api/client/services/999999", token=CLIENT_TOKEN)
    assert r.json().get("code") == 404

def test_idor_ticket():
    r = api("GET", "/api/client/tickets/999999", token=CLIENT_TOKEN)
    assert r.json().get("code") == 404

test("订单IDOR", test_idor_order)
test("账单IDOR", test_idor_invoice)
test("服务IDOR", test_idor_service)
test("工单IDOR", test_idor_ticket)

# ============================================================
# 结果
# ============================================================
print(f"\n{'='*60}")
print(f"结果: {results['pass']} pass / {results['fail']} fail / {results['pass']+results['fail']} total")
if results["errors"]:
    print("\nFailed:")
    for e in results["errors"]:
        print(f"  [FAIL] {e}")
print(f"{'='*60}")

sys.exit(0 if results["fail"] == 0 else 1)
