"""
anchor-finance 全API流程测试
- 先登录拿token
- 创建测试数据
- 走完整业务流程
- 验证未授权访问返回401
- 验证购物车/订单/工单等完整流程
"""
import requests
import json
import sys

BASE = "http://45.207.210.235:8080"
ADMIN_TOKEN = None
CLIENT_TOKEN = None
TEST_USERNAME = None
TEST_USER_ID = None
TEST_PRODUCT_ID = None
TEST_ORDER_ID = None
TEST_INVOICE_ID = None
TEST_TICKET_ID = None
TEST_CART_ITEM_ID = None
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

def admin_api(method, path, json_data=None, token=None):
    headers = {}
    if token:
        headers["Authorization"] = f"Bearer {token}"
    r = requests.request(method, f"{BASE}{path}", json=json_data, headers=headers, timeout=10)
    return r

def client_api(method, path, json_data=None, token=None):
    headers = {}
    if token:
        headers["Authorization"] = f"Bearer {token}"
    r = requests.request(method, f"{BASE}{path}", json=json_data, headers=headers, timeout=10)
    return r

# ==================== 1. 未授权访问测试 ====================
print("\n=== 1. 未授权访问测试 ===")

def test_admin_no_auth():
    r = admin_api("GET", "/api/admin/users")
    assert r.status_code in [401, 403], f"期望401/403，实际{r.status_code}"
    d = r.json()
    assert d.get("code") in [401, 403], f"期望code 401/403，实际{d.get('code')}"

def test_client_no_auth():
    r = client_api("GET", "/api/client/services")
    assert r.status_code in [401, 403], f"期望401/403，实际{r.status_code}"
    d = r.json()
    assert d.get("code") in [401, 403], f"期望code 401/403，实际{d.get('code')}"

def test_admin_wrong_password():
    r = admin_api("POST", "/api/admin/login", {"username": "admin", "password": "wrongpassword"})
    d = r.json()
    assert d.get("code") == 401, f"期望401，实际{d.get('code')}"

def test_health():
    r = requests.get(f"{BASE}/health", timeout=5)
    assert r.status_code == 200
    d = r.json()
    assert d.get("status") == "ok"

test("健康检查", test_health)
test("Admin未授权访问", test_admin_no_auth)
test("Client未授权访问", test_client_no_auth)
test("Admin错误密码", test_admin_wrong_password)

# ==================== 2. Admin登录 ====================
print("\n=== 2. Admin登录 ===")

def test_admin_login():
    global ADMIN_TOKEN
    r = admin_api("POST", "/api/admin/login", {"username": "admin", "password": "admin123"})
    d = r.json()
    assert d.get("code") == 0, f"登录失败: {d}"
    assert "token" in d.get("data", {}), f"无token: {d}"
    ADMIN_TOKEN = d["data"]["token"]
    assert len(ADMIN_TOKEN) > 20, "token太短"

test("Admin登录", test_admin_login)

# ==================== 3. Admin信息 ====================
print("\n=== 3. Admin认证信息 ===")

def test_admin_info():
    r = admin_api("GET", "/api/admin/auth/info", token=ADMIN_TOKEN)
    d = r.json()
    assert d.get("code") == 0, f"获取信息失败: {d}"
    assert d["data"]["username"] == "admin"

def test_admin_profile():
    r = admin_api("GET", "/api/admin/auth/info", token=ADMIN_TOKEN)
    d = r.json()
    assert d.get("code") == 0

test("获取Admin信息", test_admin_info)
test("Admin Profile", test_admin_profile)

# ==================== 4. Client注册+登录 ====================
print("\n=== 4. Client注册+登录 ===")

def test_client_register():
    import time
    ts = int(time.time())
    global TEST_USERNAME
    TEST_USERNAME = f"e2e_user_{ts}"
    
    # 先关闭注册验证码要求
    admin_api("PUT", "/api/admin/settings", {
        "register_require_captcha": "0"
    }, token=ADMIN_TOKEN)
    
    r = client_api("POST", "/api/client/register", {
        "username": TEST_USERNAME,
        "email": f"e2e_{ts}@test.com",
        "password": "test123456"
    })
    d = r.json()
    assert d.get("code") == 0, f"注册失败: {d}"

def test_client_login():
    global CLIENT_TOKEN
    r = client_api("POST", "/api/client/login", {
        "username": TEST_USERNAME,
        "password": "test123456"
    })
    d = r.json()
    assert d.get("code") == 0, f"Client登录失败: {d}"
    assert "token" in d.get("data", {}), f"无token: {d}"
    CLIENT_TOKEN = d["data"]["token"]

def test_client_info():
    r = client_api("GET", "/api/client/auth/info", token=CLIENT_TOKEN)
    d = r.json()
    assert d.get("code") == 0, f"获取Client信息失败: {d}"
    assert d["data"]["username"] == TEST_USERNAME

test("Client注册", test_client_register)
test("Client登录", test_client_login)
test("Client信息", test_client_info)

# ==================== 5. Admin用户管理 ====================
print("\n=== 5. Admin用户管理 ===")

def test_user_list():
    global TEST_USER_ID
    r = admin_api("GET", "/api/admin/users", token=ADMIN_TOKEN)
    d = r.json()
    assert d.get("code") == 0, f"用户列表失败: {d}"
    assert "list" in d.get("data", {}), f"无list: {d}"
    assert "total" in d.get("data", {}), f"无total: {d}"
    if len(d["data"]["list"]) > 0:
        TEST_USER_ID = d["data"]["list"][0]["id"]

def test_user_detail():
    if not TEST_USER_ID:
        return
    r = admin_api("GET", f"/api/admin/users/{TEST_USER_ID}", token=ADMIN_TOKEN)
    d = r.json()
    assert d.get("code") == 0, f"用户详情失败: {d}"
    assert d["data"]["id"] == TEST_USER_ID

def test_user_create():
    r = admin_api("POST", "/api/admin/users", {
        "username": "admin_created_user",
        "email": "admin_created@example.com",
        "password": "admin123"
    }, token=ADMIN_TOKEN)
    d = r.json()
    assert d.get("code") in [0, 400], f"创建用户异常: {d}"

def test_user_balance_logs():
    if not TEST_USER_ID:
        return
    r = admin_api("GET", f"/api/admin/users/{TEST_USER_ID}/balance-logs", token=ADMIN_TOKEN)
    d = r.json()
    assert d.get("code") == 0, f"余额日志失败: {d}"

def test_user_operation_logs():
    if not TEST_USER_ID:
        return
    r = admin_api("GET", f"/api/admin/users/{TEST_USER_ID}/operation-logs", token=ADMIN_TOKEN)
    d = r.json()
    assert d.get("code") == 0, f"操作日志失败: {d}"

test("用户列表", test_user_list)
test("用户详情", test_user_detail)
test("创建用户", test_user_create)
test("用户余额日志", test_user_balance_logs)
test("用户操作日志", test_user_operation_logs)

# ==================== 6. 产品管理 ====================
print("\n=== 6. 产品管理 ===")

def test_product_list():
    global TEST_PRODUCT_ID
    r = admin_api("GET", "/api/admin/products", token=ADMIN_TOKEN)
    d = r.json()
    assert d.get("code") == 0, f"产品列表失败: {d}"
    assert "list" in d.get("data", {})
    if len(d["data"]["list"]) > 0:
        TEST_PRODUCT_ID = d["data"]["list"][0]["id"]

def test_product_create():
    global TEST_PRODUCT_ID
    r = admin_api("POST", "/api/admin/products", {
        "name": "E2E测试产品",
        "type": "server",
        "price": 99.99,
        "billing_cycle": "monthly",
        "status": "active"
    }, token=ADMIN_TOKEN)
    d = r.json()
    assert d.get("code") == 0, f"创建产品失败: {d}"
    assert "id" in d.get("data", {}), f"无id: {d}"
    TEST_PRODUCT_ID = d["data"]["id"]

def test_product_detail():
    if not TEST_PRODUCT_ID:
        return
    r = admin_api("GET", f"/api/admin/products/{TEST_PRODUCT_ID}", token=ADMIN_TOKEN)
    d = r.json()
    assert d.get("code") == 0, f"产品详情失败: {d}"
    assert d["data"]["name"] == "E2E测试产品"

def test_product_zero_price():
    r = admin_api("POST", "/api/admin/products", {
        "name": "零元产品",
        "price": 0,
        "type": "server"
    }, token=ADMIN_TOKEN)
    d = r.json()
    assert d.get("code") == 400, f"0元购防护失败，应返回400，实际{d.get('code')}"

def test_product_groups():
    r = admin_api("GET", "/api/admin/product-groups", token=ADMIN_TOKEN)
    d = r.json()
    assert d.get("code") == 0, f"产品分组失败: {d}"

def test_product_types():
    r = admin_api("GET", "/api/admin/product-types", token=ADMIN_TOKEN)
    d = r.json()
    assert d.get("code") == 0, f"产品类型失败: {d}"

test("产品列表", test_product_list)
test("创建产品", test_product_create)
test("产品详情", test_product_detail)
test("0元购防护", test_product_zero_price)
test("产品分组", test_product_groups)
test("产品类型", test_product_types)

# ==================== 7. 购物车完整流程 ====================
print("\n=== 7. 购物车完整流程 ===")

def test_cart_empty():
    r = client_api("GET", "/api/client/cart", token=CLIENT_TOKEN)
    d = r.json()
    assert d.get("code") == 0, f"获取购物车失败: {d}"
    assert d["data"]["items"] == [], f"购物车应为空: {d}"

def test_cart_add():
    global TEST_CART_ITEM_ID
    if not TEST_PRODUCT_ID:
        return
    r = client_api("POST", "/api/client/cart", {
        "product_id": TEST_PRODUCT_ID,
        "quantity": 2,
        "cycle": "monthly"
    }, token=CLIENT_TOKEN)
    d = r.json()
    assert d.get("code") == 0, f"添加购物车失败: {d}"

def test_cart_not_empty():
    r = client_api("GET", "/api/client/cart", token=CLIENT_TOKEN)
    d = r.json()
    assert d.get("code") == 0
    assert len(d["data"]["items"]) > 0, f"购物车应有商品: {d}"
    TEST_CART_ITEM_ID = d["data"]["items"][0]["id"]

def test_cart_update():
    if not TEST_CART_ITEM_ID:
        return
    r = client_api("PUT", f"/api/client/cart/{TEST_CART_ITEM_ID}", {
        "quantity": 3
    }, token=CLIENT_TOKEN)
    d = r.json()
    assert d.get("code") == 0, f"更新购物车失败: {d}"

def test_cart_checkout():
    if not TEST_PRODUCT_ID:
        return
    r = client_api("POST", "/api/client/cart/checkout", {}, token=CLIENT_TOKEN)
    d = r.json()
    assert d.get("code") in [0, 400], f"结算异常: {d}"
    if d.get("code") == 0:
        assert "order_id" in d.get("data", {}), f"无order_id: {d}"
        assert "invoice_id" in d.get("data", {}), f"无invoice_id: {d}"

test("购物车为空", test_cart_empty)
test("添加购物车", test_cart_add)
test("购物车非空", test_cart_not_empty)
test("更新购物车", test_cart_update)
test("购物车结算", test_cart_checkout)

# ==================== 8. 订单管理 ====================
print("\n=== 8. 订单管理 ===")

def test_order_list():
    global TEST_ORDER_ID
    r = admin_api("GET", "/api/admin/orders", token=ADMIN_TOKEN)
    d = r.json()
    assert d.get("code") == 0, f"订单列表失败: {d}"
    assert "list" in d.get("data", {})
    if len(d["data"]["list"]) > 0:
        TEST_ORDER_ID = d["data"]["list"][0]["id"]

def test_order_detail():
    if not TEST_ORDER_ID:
        return
    r = admin_api("GET", f"/api/admin/orders/{TEST_ORDER_ID}", token=ADMIN_TOKEN)
    d = r.json()
    assert d.get("code") == 0, f"订单详情失败: {d}"

def test_client_order_list():
    r = client_api("GET", "/api/client/orders", token=CLIENT_TOKEN)
    d = r.json()
    assert d.get("code") == 0, f"Client订单列表失败: {d}"

def test_order_summary():
    r = client_api("GET", "/api/client/orders/summary", token=CLIENT_TOKEN)
    d = r.json()
    assert d.get("code") == 0, f"订单统计失败: {d}"

test("Admin订单列表", test_order_list)
test("Admin订单详情", test_order_detail)
test("Client订单列表", test_client_order_list)
test("订单统计", test_order_summary)

# ==================== 9. 账单管理 ====================
print("\n=== 9. 账单管理 ===")

def test_invoice_list():
    global TEST_INVOICE_ID
    r = admin_api("GET", "/api/admin/invoices", token=ADMIN_TOKEN)
    d = r.json()
    assert d.get("code") == 0, f"账单列表失败: {d}"
    if len(d["data"]["list"]) > 0:
        TEST_INVOICE_ID = d["data"]["list"][0]["id"]

def test_invoice_detail():
    if not TEST_INVOICE_ID:
        return
    r = admin_api("GET", f"/api/admin/invoices/{TEST_INVOICE_ID}", token=ADMIN_TOKEN)
    d = r.json()
    assert d.get("code") == 0, f"账单详情失败: {d}"

def test_client_invoices():
    r = client_api("GET", "/api/client/invoices", token=CLIENT_TOKEN)
    d = r.json()
    assert d.get("code") == 0, f"Client账单列表失败: {d}"

def test_invoice_summary():
    r = client_api("GET", "/api/client/invoices/summary", token=CLIENT_TOKEN)
    d = r.json()
    assert d.get("code") == 0, f"账单统计失败: {d}"

test("Admin账单列表", test_invoice_list)
test("Admin账单详情", test_invoice_detail)
test("Client账单列表", test_client_invoices)
test("账单统计", test_invoice_summary)

# ==================== 10. 工单系统 ====================
print("\n=== 10. 工单系统 ===")

def test_ticket_create():
    global TEST_TICKET_ID
    r = client_api("POST", "/api/client/tickets", {
        "subject": "E2E测试工单",
        "content": "这是一个自动化测试工单",
        "priority": "medium"
    }, token=CLIENT_TOKEN)
    d = r.json()
    assert d.get("code") in [0, 400], f"创建工单异常: {d}"
    if d.get("code") == 0:
        TEST_TICKET_ID = d.get("data", {}).get("id")

def test_ticket_list():
    global TEST_TICKET_ID
    r = client_api("GET", "/api/client/tickets", token=CLIENT_TOKEN)
    d = r.json()
    assert d.get("code") == 0, f"工单列表失败: {d}"
    if not TEST_TICKET_ID and len(d["data"]["list"]) > 0:
        TEST_TICKET_ID = d["data"]["list"][0]["id"]

def test_ticket_detail():
    if not TEST_TICKET_ID:
        return
    r = client_api("GET", f"/api/client/tickets/{TEST_TICKET_ID}", token=CLIENT_TOKEN)
    d = r.json()
    assert d.get("code") == 0, f"工单详情失败: {d}"

def test_ticket_reply():
    if not TEST_TICKET_ID:
        return
    r = client_api("POST", f"/api/client/tickets/{TEST_TICKET_ID}/reply", {
        "content": "这是测试回复"
    }, token=CLIENT_TOKEN)
    d = r.json()
    assert d.get("code") == 0, f"工单回复失败: {d}"

def test_ticket_replies():
    if not TEST_TICKET_ID:
        return
    r = client_api("GET", f"/api/client/tickets/{TEST_TICKET_ID}/replies", token=CLIENT_TOKEN)
    d = r.json()
    assert d.get("code") == 0, f"工单回复列表失败: {d}"

def test_admin_tickets():
    r = admin_api("GET", "/api/admin/tickets", token=ADMIN_TOKEN)
    d = r.json()
    assert d.get("code") == 0, f"Admin工单列表失败: {d}"

def test_ticket_summary():
    r = admin_api("GET", "/api/admin/tickets/summary", token=ADMIN_TOKEN)
    d = r.json()
    assert d.get("code") == 0, f"工单统计失败: {d}"

test("创建工单", test_ticket_create)
test("Client工单列表", test_ticket_list)
test("工单详情", test_ticket_detail)
test("工单回复", test_ticket_reply)
test("工单回复列表", test_ticket_replies)
test("Admin工单列表", test_admin_tickets)
test("工单统计", test_ticket_summary)

# ==================== 11. 仪表盘 ====================
print("\n=== 11. 仪表盘 ===")

def test_dashboard_stats():
    r = admin_api("GET", "/api/admin/dashboard/stats", token=ADMIN_TOKEN)
    d = r.json()
    assert d.get("code") == 0, f"仪表盘统计失败: {d}"

def test_dashboard_income():
    r = admin_api("GET", "/api/admin/dashboard/income-trend", token=ADMIN_TOKEN)
    d = r.json()
    assert d.get("code") == 0, f"收入趋势失败: {d}"

def test_dashboard_revenue():
    r = admin_api("GET", "/api/admin/dashboard/monthly-revenue", token=ADMIN_TOKEN)
    d = r.json()
    assert d.get("code") == 0, f"月收入失败: {d}"

test("仪表盘统计", test_dashboard_stats)
test("收入趋势", test_dashboard_income)
test("月收入", test_dashboard_revenue)

# ==================== 12. 服务管理 ====================
print("\n=== 12. 服务管理 ===")

def test_service_list():
    r = admin_api("GET", "/api/admin/services", token=ADMIN_TOKEN)
    d = r.json()
    assert d.get("code") == 0, f"服务列表失败: {d}"

def test_client_services():
    r = client_api("GET", "/api/client/services", token=CLIENT_TOKEN)
    d = r.json()
    assert d.get("code") == 0, f"Client服务列表失败: {d}"

def test_service_grouped():
    r = client_api("GET", "/api/client/services/grouped-overview", token=CLIENT_TOKEN)
    d = r.json()
    assert d.get("code") == 0, f"服务分组概览失败: {d}"

test("Admin服务列表", test_service_list)
test("Client服务列表", test_client_services)
test("服务分组概览", test_service_grouped)

# ==================== 13. 设置管理 ====================
print("\n=== 13. 设置管理 ===")

def test_settings():
    r = admin_api("GET", "/api/admin/settings", token=ADMIN_TOKEN)
    d = r.json()
    assert d.get("code") == 0, f"设置列表失败: {d}"

def test_settings_group():
    r = admin_api("GET", "/api/admin/settings/general", token=ADMIN_TOKEN)
    d = r.json()
    assert d.get("code") == 0, f"通用设置失败: {d}"

def test_settings_update():
    r = admin_api("PUT", "/api/admin/settings", {
        "settings": {"site_name": "E2E测试站点"}
    }, token=ADMIN_TOKEN)
    d = r.json()
    assert d.get("code") == 0, f"更新设置失败: {d}"

test("设置列表", test_settings)
test("通用设置", test_settings_group)
test("更新设置", test_settings_update)

# ==================== 14. 内容管理 ====================
print("\n=== 14. 内容管理 ===")

def test_news_list():
    r = admin_api("GET", "/api/admin/news", token=ADMIN_TOKEN)
    d = r.json()
    assert d.get("code") == 0, f"新闻列表失败: {d}"

def test_news_create():
    r = admin_api("POST", "/api/admin/news", {
        "title": "E2E测试新闻",
        "content": "测试内容",
        "status": "published"
    }, token=ADMIN_TOKEN)
    d = r.json()
    assert d.get("code") in [0, 400], f"创建新闻异常: {d}"

def test_content_overview():
    r = client_api("GET", "/api/client/content/overview")
    d = r.json()
    assert d.get("code") == 0, f"内容概览失败: {d}"

def test_home_hero():
    r = client_api("GET", "/api/client/home-hero")
    d = r.json()
    assert d.get("code") == 0, f"首页Hero失败: {d}"

def test_notices():
    r = client_api("GET", "/api/client/notices")
    d = r.json()
    assert d.get("code") == 0, f"公告列表失败: {d}"

def test_help_articles():
    r = client_api("GET", "/api/client/help-articles")
    d = r.json()
    assert d.get("code") == 0, f"帮助文章失败: {d}"

test("新闻列表", test_news_list)
test("创建新闻", test_news_create)
test("内容概览", test_content_overview)
test("首页Hero", test_home_hero)
test("公告列表", test_notices)
test("帮助文章", test_help_articles)

# ==================== 15. 插件管理 ====================
print("\n=== 15. 插件管理 ===")

def test_plugin_list():
    r = admin_api("GET", "/api/admin/plugins", token=ADMIN_TOKEN)
    d = r.json()
    assert d.get("code") == 0, f"插件列表失败: {d}"

def test_payment_gateways():
    r = admin_api("GET", "/api/admin/payment-gateways", token=ADMIN_TOKEN)
    d = r.json()
    assert d.get("code") == 0, f"支付网关失败: {d}"

def test_sms_providers():
    r = admin_api("GET", "/api/admin/sms-providers", token=ADMIN_TOKEN)
    d = r.json()
    assert d.get("code") == 0, f"短信供应商失败: {d}"

def test_mail_providers():
    r = admin_api("GET", "/api/admin/mail-providers", token=ADMIN_TOKEN)
    d = r.json()
    assert d.get("code") == 0, f"邮件供应商失败: {d}"

test("插件列表", test_plugin_list)
test("支付网关", test_payment_gateways)
test("短信供应商", test_sms_providers)
test("邮件供应商", test_mail_providers)

# ==================== 16. 财务报表 ====================
print("\n=== 16. 财务报表 ===")

def test_finance_ledger():
    r = admin_api("GET", "/api/admin/finance/ledger", token=ADMIN_TOKEN)
    d = r.json()
    assert d.get("code") == 0, f"财务账本失败: {d}"

def test_finance_summary():
    r = admin_api("GET", "/api/admin/finance/ledger/summary", token=ADMIN_TOKEN)
    d = r.json()
    assert d.get("code") == 0, f"财务摘要失败: {d}"

def test_recharge_list():
    r = admin_api("GET", "/api/admin/finance/recharges", token=ADMIN_TOKEN)
    d = r.json()
    assert d.get("code") == 0, f"充值列表失败: {d}"

test("财务账本", test_finance_ledger)
test("财务摘要", test_finance_summary)
test("充值列表", test_recharge_list)

# ==================== 17. 供应商管理 ====================
print("\n=== 17. 供应商管理 ===")

def test_supplier_list():
    r = admin_api("GET", "/api/admin/suppliers", token=ADMIN_TOKEN)
    d = r.json()
    assert d.get("code") == 0, f"供应商列表失败: {d}"

def test_supplier_summary():
    r = admin_api("GET", "/api/admin/suppliers/summary", token=ADMIN_TOKEN)
    d = r.json()
    assert d.get("code") == 0, f"供应商统计失败: {d}"

test("供应商列表", test_supplier_list)
test("供应商统计", test_supplier_summary)

# ==================== 18. 通知系统 ====================
print("\n=== 18. 通知系统 ===")

def test_notifications():
    r = client_api("GET", "/api/client/notifications", token=CLIENT_TOKEN)
    d = r.json()
    assert d.get("code") == 0, f"通知列表失败: {d}"

def test_unread_count():
    r = client_api("GET", "/api/client/notifications/unread-count", token=CLIENT_TOKEN)
    d = r.json()
    assert d.get("code") == 0, f"未读数失败: {d}"

test("通知列表", test_notifications)
test("未读数", test_unread_count)

# ==================== 19. 支付记录 ====================
print("\n=== 19. 支付记录 ===")

def test_payment_list():
    r = client_api("GET", "/api/client/payments", token=CLIENT_TOKEN)
    d = r.json()
    assert d.get("code") == 0, f"支付列表失败: {d}"

def test_payment_summary():
    r = client_api("GET", "/api/client/payments/summary", token=CLIENT_TOKEN)
    d = r.json()
    assert d.get("code") == 0, f"支付统计失败: {d}"

test("支付列表", test_payment_list)
test("支付统计", test_payment_summary)

# ==================== 20. 日志系统 ====================
print("\n=== 20. 日志系统 ===")

def test_system_logs():
    r = admin_api("GET", "/api/admin/system-logs", token=ADMIN_TOKEN)
    d = r.json()
    assert d.get("code") == 0, f"系统日志失败: {d}"

def test_operation_logs():
    r = admin_api("GET", "/api/admin/operation-logs", token=ADMIN_TOKEN)
    d = r.json()
    assert d.get("code") == 0, f"操作日志失败: {d}"

def test_login_logs():
    r = admin_api("GET", "/api/admin/login-logs", token=ADMIN_TOKEN)
    d = r.json()
    assert d.get("code") == 0, f"登录日志失败: {d}"

def test_sms_logs():
    r = admin_api("GET", "/api/admin/logs/sms", token=ADMIN_TOKEN)
    d = r.json()
    assert d.get("code") == 0, f"短信日志失败: {d}"

def test_email_logs():
    r = admin_api("GET", "/api/admin/logs/email", token=ADMIN_TOKEN)
    d = r.json()
    assert d.get("code") == 0, f"邮件日志失败: {d}"

test("系统日志", test_system_logs)
test("操作日志", test_operation_logs)
test("登录日志", test_login_logs)
test("短信日志", test_sms_logs)
test("邮件日志", test_email_logs)

# ==================== 21. 推介系统 ====================
print("\n=== 21. 推介系统 ===")

def test_referral_overview():
    r = client_api("GET", "/api/client/referral/overview", token=CLIENT_TOKEN)
    d = r.json()
    assert d.get("code") == 0, f"推介概览失败: {d}"

def test_referral_rewards():
    r = client_api("GET", "/api/client/referral/rewards", token=CLIENT_TOKEN)
    d = r.json()
    assert d.get("code") == 0, f"推介奖励失败: {d}"

test("推介概览", test_referral_overview)
test("推介奖励", test_referral_rewards)

# ==================== 22. 实名认证 ====================
print("\n=== 22. 实名认证 ===")

def test_verification_status():
    r = client_api("GET", "/api/client/verification/status", token=CLIENT_TOKEN)
    d = r.json()
    assert d.get("code") == 0, f"认证状态失败: {d}"

def test_admin_verifications():
    r = admin_api("GET", "/api/admin/verifications", token=ADMIN_TOKEN)
    d = r.json()
    assert d.get("code") == 0, f"认证列表失败: {d}"

def test_verification_summary():
    r = admin_api("GET", "/api/admin/verifications/summary", token=ADMIN_TOKEN)
    d = r.json()
    assert d.get("code") == 0, f"认证统计失败: {d}"

test("认证状态", test_verification_status)
test("Admin认证列表", test_admin_verifications)
test("认证统计", test_verification_summary)

# ==================== 23. 货币管理 ====================
print("\n=== 23. 货币管理 ===")

def test_currencies():
    r = admin_api("GET", "/api/admin/currencies", token=ADMIN_TOKEN)
    d = r.json()
    assert d.get("code") == 0, f"货币列表失败: {d}"

test("货币列表", test_currencies)

# ==================== 24. 角色权限 ====================
print("\n=== 24. 角色权限 ===")

def test_roles():
    r = admin_api("GET", "/api/admin/roles", token=ADMIN_TOKEN)
    d = r.json()
    assert d.get("code") == 0, f"角色列表失败: {d}"

def test_permissions():
    r = admin_api("GET", "/api/admin/permissions", token=ADMIN_TOKEN)
    d = r.json()
    assert d.get("code") == 0, f"权限列表失败: {d}"

def test_staff():
    r = admin_api("GET", "/api/admin/staff", token=ADMIN_TOKEN)
    d = r.json()
    assert d.get("code") == 0, f"员工列表失败: {d}"

test("角色列表", test_roles)
test("权限列表", test_permissions)
test("员工列表", test_staff)

# ==================== 25. 定时任务 ====================
print("\n=== 25. 定时任务 ===")

def test_cron_tasks():
    r = admin_api("GET", "/api/admin/cron-tasks", token=ADMIN_TOKEN)
    d = r.json()
    assert d.get("code") == 0, f"定时任务失败: {d}"

def test_schedule_overview():
    r = admin_api("GET", "/api/admin/schedules/overview", token=ADMIN_TOKEN)
    d = r.json()
    assert d.get("code") == 0, f"调度概览失败: {d}"

def test_schedule_runs():
    r = admin_api("GET", "/api/admin/schedule-runs", token=ADMIN_TOKEN)
    d = r.json()
    assert d.get("code") == 0, f"调度运行失败: {d}"

test("定时任务", test_cron_tasks)
test("调度概览", test_schedule_overview)
test("调度运行", test_schedule_runs)

# ==================== 26. 客户端产品浏览 ====================
print("\n=== 26. 客户端产品浏览 ===")

def test_client_products():
    r = client_api("GET", "/api/client/products", token=CLIENT_TOKEN)
    d = r.json()
    assert d.get("code") == 0, f"产品列表失败: {d}"

def test_client_product_categories():
    r = client_api("GET", "/api/client/products/categories", token=CLIENT_TOKEN)
    d = r.json()
    assert d.get("code") == 0, f"产品分类失败: {d}"

test("Client产品列表", test_client_products)
test("Client产品分类", test_client_product_categories)

# ==================== 27. 黑名单 ====================
print("\n=== 27. 黑名单 ===")

def test_blacklist():
    r = admin_api("GET", "/api/admin/blacklist", token=ADMIN_TOKEN)
    d = r.json()
    assert d.get("code") == 0, f"黑名单失败: {d}"

test("黑名单", test_blacklist)

# ==================== 28. 主题模板 ====================
print("\n=== 28. 主题模板 ===")

def test_themes():
    r = admin_api("GET", "/api/admin/themes", token=ADMIN_TOKEN)
    d = r.json()
    assert d.get("code") == 0, f"主题列表失败: {d}"

def test_home_hero_admin():
    r = admin_api("GET", "/api/admin/home-hero", token=ADMIN_TOKEN)
    d = r.json()
    assert d.get("code") == 0, f"Admin首页Hero失败: {d}"

test("主题列表", test_themes)
test("Admin首页Hero", test_home_hero_admin)

# ==================== 29. 信用额管理 ====================
print("\n=== 29. 信用额管理 ===")

def test_credit_limits():
    r = admin_api("GET", "/api/admin/credit-limits", token=ADMIN_TOKEN)
    d = r.json()
    assert d.get("code") == 0, f"信用额列表失败: {d}"

def test_credit_config():
    r = admin_api("GET", "/api/admin/credit-limits/config", token=ADMIN_TOKEN)
    d = r.json()
    assert d.get("code") == 0, f"信用额配置失败: {d}"

test("信用额列表", test_credit_limits)
test("信用额配置", test_credit_config)

# ==================== 30. 优惠券 ====================
print("\n=== 30. 优惠券 ===")

def test_coupons():
    r = admin_api("GET", "/api/admin/coupons", token=ADMIN_TOKEN)
    d = r.json()
    assert d.get("code") == 0, f"优惠券列表失败: {d}"

def test_coupon_summary():
    r = admin_api("GET", "/api/admin/coupons/summary", token=ADMIN_TOKEN)
    d = r.json()
    assert d.get("code") == 0, f"优惠券统计失败: {d}"

def test_client_coupons():
    r = client_api("GET", "/api/client/coupons", token=CLIENT_TOKEN)
    d = r.json()
    assert d.get("code") == 0, f"Client优惠券失败: {d}"

test("优惠券列表", test_coupons)
test("优惠券统计", test_coupon_summary)
test("Client优惠券", test_client_coupons)

# ==================== 31. 合同管理 ====================
print("\n=== 31. 合同管理 ===")

def test_contracts():
    r = admin_api("GET", "/api/admin/contracts", token=ADMIN_TOKEN)
    d = r.json()
    assert d.get("code") == 0, f"合同列表失败: {d}"

test("合同列表", test_contracts)

# ==================== 32. 营销推送 ====================
print("\n=== 32. 营销推送 ===")

def test_marketing():
    r = admin_api("GET", "/api/admin/marketing/pushes", token=ADMIN_TOKEN)
    d = r.json()
    assert d.get("code") == 0, f"营销列表失败: {d}"

test("营销列表", test_marketing)

# ==================== 33. 系统信息 ====================
print("\n=== 33. 系统信息 ===")

def test_system_info():
    r = admin_api("GET", "/api/admin/system/info", token=ADMIN_TOKEN)
    d = r.json()
    assert d.get("code") == 0, f"系统信息失败: {d}"

def test_system_modules():
    r = admin_api("GET", "/api/admin/system/modules", token=ADMIN_TOKEN)
    d = r.json()
    assert d.get("code") == 0, f"系统模块失败: {d}"

test("系统信息", test_system_info)
test("系统模块", test_system_modules)

# ==================== 34. 数据库管理 ====================
print("\n=== 34. 数据库管理 ===")

def test_database_status():
    r = admin_api("GET", "/api/admin/database/status", token=ADMIN_TOKEN)
    d = r.json()
    assert d.get("code") == 0, f"数据库状态失败: {d}"

test("数据库状态", test_database_status)

# ==================== 35. 第三方登录 ====================
print("\n=== 35. 第三方登录 ===")

def test_oauth_providers():
    r = admin_api("GET", "/api/admin/oauth-providers", token=ADMIN_TOKEN)
    d = r.json()
    assert d.get("code") == 0, f"OAuth供应商失败: {d}"

test("OAuth供应商", test_oauth_providers)

# ==================== 36. 邮件/短信模板 ====================
print("\n=== 36. 邮件/短信模板 ===")

def test_email_templates():
    r = admin_api("GET", "/api/admin/email-templates", token=ADMIN_TOKEN)
    d = r.json()
    assert d.get("code") == 0, f"邮件模板失败: {d}"

def test_sms_templates():
    r = admin_api("GET", "/api/admin/sms-templates", token=ADMIN_TOKEN)
    d = r.json()
    assert d.get("code") == 0, f"短信模板失败: {d}"

test("邮件模板", test_email_templates)
test("短信模板", test_sms_templates)

# ==================== 37. 友情链接 ====================
print("\n=== 37. 友情链接 ===")

def test_friendly_links():
    r = admin_api("GET", "/api/admin/friendly-links", token=ADMIN_TOKEN)
    d = r.json()
    assert d.get("code") == 0, f"友情链接失败: {d}"

test("友情链接", test_friendly_links)

# ==================== 38. 流量包 ====================
print("\n=== 38. 流量包 ===")

def test_traffic_packages():
    r = admin_api("GET", "/api/admin/traffic-packages", token=ADMIN_TOKEN)
    d = r.json()
    assert d.get("code") == 0, f"流量包失败: {d}"

test("流量包", test_traffic_packages)

# ==================== 39. 任务队列 ====================
print("\n=== 39. 任务队列 ===")

def test_task_queue():
    r = admin_api("GET", "/api/admin/task-queue/overview", token=ADMIN_TOKEN)
    d = r.json()
    assert d.get("code") == 0, f"任务队列失败: {d}"

test("任务队列", test_task_queue)

# ==================== 40. 二次验证 ====================
print("\n=== 40. 二次验证 ===")

def test_two_factor():
    r = admin_api("GET", "/api/admin/two-factor-config", token=ADMIN_TOKEN)
    d = r.json()
    assert d.get("code") == 0, f"二次验证失败: {d}"

test("二次验证", test_two_factor)

# ==================== 结果汇总 ====================
print(f"\n{'='*50}")
print(f"测试结果: 通过 {results['pass']} / 失败 {results['fail']} / 总计 {results['pass']+results['fail']}")
if results["errors"]:
    print(f"\nFailed details:")
    for e in results["errors"]:
        print(f"  [FAIL] {e}")
print(f"{'='*50}")

sys.exit(0 if results["fail"] == 0 else 1)
