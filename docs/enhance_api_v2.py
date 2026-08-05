#!/usr/bin/env python3
"""Extract model structs to fill response fields for handlers returning model objects."""
import re, os, glob, json

HANDLER_DIR = os.path.join(os.path.dirname(__file__), '..', 'backend', 'internal', 'handler')
MODEL_DIR = os.path.join(os.path.dirname(__file__), '..', 'backend', 'internal', 'model')
SERVICE_DIR = os.path.join(os.path.dirname(__file__), '..', 'backend', 'internal', 'service')
API_DIR = os.path.join(os.path.dirname(__file__), '..', 'backend', 'internal', 'api')
DOCS_DIR = os.path.dirname(__file__)

# Variable name -> likely model type mapping
VAR_TO_MODEL = {
    'user': 'User', 'users': '[]User', 'host': 'Host', 'hosts': '[]Host',
    'order': 'Order', 'orders': '[]Order', 'invoice': 'Invoice', 'invoices': '[]Invoice',
    'ticket': 'Ticket', 'tickets': '[]Ticket', 'product': 'Product', 'products': '[]Product',
    'host': 'Host', 'hosts': '[]Host', 'host': 'Host',
    'banner': 'Banner', 'banners': '[]Banner', 'announcement': 'Announcement',
    'news': 'News', 'article': 'Article', 'articles': '[]Article',
    'group': 'Group', 'groups': '[]Group', 'category': 'Category',
    'provider': 'Provider', 'providers': '[]Provider', 'payment_method': 'PaymentMethod',
    'gateway': 'Gateway', 'gateways': '[]Gateway', 'config': 'SystemConfig',
    'cfg': 'SystemConfig', 'setting': 'SystemConfig', 'settings': '[]SystemConfig',
    'admin': 'Admin', 'admins': '[]Admin', 'role': 'Role', 'roles': '[]Role',
    'server': 'Server', 'servers': '[]Server', 'module': 'ServerModule',
    'log': 'Log', 'logs': '[]Log', 'record': 'Record', 'records': '[]Record',
    'message': 'Message', 'messages': '[]Message', 'contact': 'Contact',
    'contacts': '[]Contact', 'note': 'Note', 'notes': '[]Note',
    'nav': 'Nav', 'menus': '[]Nav', 'menu': 'Nav', 'link': 'FriendlyLink',
    'links': '[]FriendlyLink', 'promo_code': 'PromoCode', 'promo_codes': '[]PromoCode',
    'suffix': 'EmailSuffix', 'suffixes': '[]EmailSuffix', 'dept': 'Department',
    'depts': '[]Department', 'department': 'Department', 'departments': '[]Department',
    'cert': 'Certification', 'certification': 'Certification', 'plugin': 'Plugin',
    'plugins': '[]Plugin', 'template': 'Template', 'templates': '[]Template',
    'datacenter': 'Datacenter', 'datacenters': '[]Datacenter', 'iso': 'IsoImage',
    'isos': '[]IsoImage', 'snapshot': 'Snapshot', 'snapshots': '[]Snapshot',
    'backup': 'Backups', 'backups': '[]Backup', 'session': 'ChatSession',
    'sessions': '[]ChatSession', 'listing': 'MarketListing', 'listings': '[]MarketListing',
    'balance': 'Balance', 'transaction': 'Transaction', 'transactions': '[]Transaction',
    'ssl': 'SSL', 'ssl_cert': 'SSL', 'contract': 'Contract',
    'contracts': '[]Contract', 'sale': 'SaleRecord', 'sales': '[]SaleRecord',
    'agent': 'Agent', 'agents': '[]Agent', 'kb': 'KBArticle', 'kbs': '[]KBArticle',
    'cancel_reason': 'CancelReason', 'cancel_reasons': '[]CancelReason',
    'tracker': 'ClientTracker', 'trackers': '[]ClientTracker',
    'care': 'ClientCare', 'vouchers': '[]Voucher', 'voucher': 'Voucher',
    'voucher_type': 'VoucherType', 'voucher_types': '[]VoucherType',
    'divert': 'ProductDivert', 'diverts': '[]ProductDivert',
    'task': 'CronTask', 'tasks': '[]CronTask', 'transfer': 'Transfer',
    'account': 'Account', 'oauth_account': 'OAuthAccount',
    'note': 'TicketNote', 'notes': '[]TicketNote',
    'info': 'Info', 'result': 'Result', 'stats': 'Stats',
    'options': '[]ConfigOption', 'option': 'ConfigOption',
    'option_group': 'OptionGroup', 'option_groups': '[]OptionGroup',
    'item': 'Item', 'items': '[]Item', 'list': '[]Item',
    'ssl_list': '[]SSL', 'error': 'Error',
}

# Model fields by model name
MODEL_FIELDS = {}

def extract_model_fields(content, filename):
    """Extract struct fields from model files."""
    for m in re.finditer(r'type\s+(\w+)\s+struct\s*\{([^}]*(?:\{[^}]*\}[^}]*)*)\}', content, re.DOTALL):
        name = m.group(1)
        body = m.group(2)
        fields = []
        for fm in re.finditer(r'(\w+)\s+(\S+)\s+`json:"([^"]+)"[^`]*`', body):
            fname = fm.group(3).split(',')[0]
            if fname in ('-', ''): continue
            ftype = fm.group(2).strip('*')
            if ftype.startswith('[]'): ftype = 'array'
            elif ftype.startswith('map['): ftype = 'object'
            elif ftype in ('string','int','int64','uint','float32','float64','bool','time.Time','datatypes.JSON'): pass
            else: ftype = 'object'
            fields.append({'n': fname, 't': ftype})
        if fields:
            MODEL_FIELDS[name] = fields

# Common response field mappings for gin.H returns
COMMON_GIN_H = {
    'token': ('string', 'JWT令牌'),
    'access_token': ('string', '访问令牌'),
    'refresh_token': ('string', '刷新令牌'),
    'expires_at': ('int', '过期时间戳'),
    'user': ('object', '用户信息'),
    'admin': ('object', '管理员信息'),
    'host': ('object', '主机信息'),
    'order': ('object', '订单信息'),
    'invoice': ('object', '账单信息'),
    'ticket': ('object', '工单信息'),
    'product': ('object', '产品信息'),
    'provider': ('object', '提供商信息'),
    'providers': ('array', '提供商列表'),
    'accounts': ('array', '已绑定账号列表'),
    'list': ('array', '列表'),
    'total': ('int', '总数'),
    'page': ('int', '页码'),
    'page_size': ('int', '每页数量'),
    'message': ('string', '提示消息'),
    'is_active': ('bool', '是否启用'),
    'code': ('string', '代码/编号'),
    'name': ('string', '名称'),
    'id': ('int', 'ID'),
    'key': ('string', '密钥'),
    'password': ('string', '密码'),
    'url': ('string', 'URL'),
    'status': ('string', '状态'),
    'amount': ('number', '金额'),
    'count': ('int', '数量'),
    'imported': ('int', '导入数量'),
    'deleted': ('int', '删除数量'),
    'updated': ('int', '更新数量'),
    'results': ('array', '结果列表'),
    'data': ('object', '数据'),
    'config': ('object', '配置'),
    'settings': ('object', '设置'),
    'options': ('array', '选项列表'),
    'groups': ('array', '分组列表'),
    'menu': ('object', '菜单'),
    'menus': ('array', '菜单列表'),
    'nav': ('object', '导航'),
    'navs': ('array', '导航列表'),
    'banner': ('object', '轮播图'),
    'banners': ('array', '轮播图列表'),
    'announcement': ('object', '公告'),
    'announcements': ('array', '公告列表'),
    'news': ('object', '新闻'),
    'news_list': ('array', '新闻列表'),
    'article': ('object', '文章'),
    'articles': ('array', '文章列表'),
    'categories': ('array', '分类列表'),
    'category': ('object', '分类'),
    'links': ('array', '链接列表'),
    'link': ('object', '链接'),
    'servers': ('array', '服务器列表'),
    'server': ('object', '服务器'),
    'modules': ('array', '模块列表'),
    'module': ('object', '模块'),
    'plugin': ('object', '插件'),
    'plugins': ('array', '插件列表'),
    'template': ('object', '模板'),
    'templates': ('array', '模板列表'),
    'datacenter': ('object', '数据中心'),
    'datacenters': ('array', '数据中心列表'),
    'iso': ('object', 'ISO镜像'),
    'isos': ('array', 'ISO镜像列表'),
    'snapshot': ('object', '快照'),
    'snapshots': ('array', '快照列表'),
    'backup': ('object', '备份'),
    'backups': ('array', '备份列表'),
    'session': ('object', '会话'),
    'sessions': ('array', '会话列表'),
    'listing': ('object', '挂售信息'),
    'listings': ('array', '挂售列表'),
    'balance': ('number', '余额'),
    'credit': ('number', '信用额度'),
    'transaction': ('object', '交易记录'),
    'transactions': ('array', '交易记录列表'),
    'ssl': ('object', 'SSL证书'),
    'ssl_list': ('array', 'SSL证书列表'),
    'contract': ('object', '合同'),
    'contracts': ('array', '合同列表'),
    'sale': ('object', '销售记录'),
    'sales': ('array', '销售记录列表'),
    'agent': ('object', '代理'),
    'agents': ('array', '代理列表'),
    'kb': ('object', '知识库文章'),
    'kbs': ('array', '知识库文章列表'),
    'cancel_reason': ('object', '注销原因'),
    'cancel_reasons': ('array', '注销原因列表'),
    'tracker': ('object', '客户跟踪'),
    'trackers': ('array', '客户跟踪列表'),
    'care': ('object', '客户关怀'),
    'cares': ('array', '客户关怀列表'),
    'voucher': ('object', '发票'),
    'vouchers': ('array', '发票列表'),
    'voucher_type': ('object', '发票抬头'),
    'voucher_types': ('array', '发票抬头列表'),
    'divert': ('object', '产品转移'),
    'diverts': ('array', '产品转移列表'),
    'task': ('object', '定时任务'),
    'tasks': ('array', '定时任务列表'),
    'transfer': ('object', '转移信息'),
    'account': ('object', '账号信息'),
    'oauth_account': ('object', 'OAuth账号'),
    'note': ('object', '备注'),
    'notes': ('array', '备注列表'),
    'info': ('object', '信息'),
    'result': ('object', '结果'),
    'stats': ('object', '统计'),
    'config_options': ('array', '配置选项列表'),
    'option_group': ('object', '选项组'),
    'option_groups': ('array', '选项组列表'),
    'item': ('object', '条目'),
    'items': ('array', '条目列表'),
    'error': ('string', '错误信息'),
    'login_url': ('string', 'OAuth登录URL'),
    'redirect_url': ('string', '重定向URL'),
    'oauth_code': ('string', 'OAuth授权码'),
    'provider_name': ('string', '提供商名称'),
    'provider_code': ('string', '提供商代码'),
    'avatar_url': ('string', '头像URL'),
    'nickname': ('string', '昵称'),
    'bindding': ('object', '绑定信息'),
    'ok': ('bool', '是否成功'),
    'success': ('bool', '是否成功'),
    'applied': ('bool', '是否已申请'),
    'available': ('bool', '是否可用'),
    'captcha_id': ('string', '验证码ID'),
    'captcha_type': ('string', '验证码类型'),
    'image': ('string', '图片Base64'),
    'expired_at': ('string', '过期时间'),
    'ticket_count': ('int', '工单数量'),
    'order_count': ('int', '订单数量'),
    'host_count': ('int', '主机数量'),
    'user_count': ('int', '用户数量'),
    'income': ('number', '收入'),
    'expenses': ('number', '支出'),
    'commission': ('number', '佣金'),
    'referral_code': ('string', '推荐码'),
    'referral_link': ('string', '推荐链接'),
    'ip': ('string', 'IP地址'),
    'ipv6': ('string', 'IPv6地址'),
    'mac': ('string', 'MAC地址'),
    'hostname': ('string', '主机名'),
    'os': ('string', '操作系统'),
    'cpu': ('string', 'CPU型号'),
    'cpu_cores': ('int', 'CPU核数'),
    'memory_mb': ('int', '内存(MB)'),
    'disk_size_gb': ('int', '磁盘(GB)'),
    'bandwidth_mbps': ('int', '带宽(Mbps)'),
    'traffic_gb': ('int', '流量(GB)'),
    'virtual_type': ('string', '虚拟化类型'),
    'datacenter_id': ('int', '数据中心ID'),
    'product_name': ('string', '产品名称'),
    'group_name': ('string', '分组名称'),
    'domain': ('string', '域名'),
    'due_date': ('string', '到期日期'),
    'created_at': ('string', '创建时间'),
    'updated_at': ('string', '更新时间'),
    'login_url': ('string', '登录URL'),
    'logout_url': ('string', '登出URL'),
    'callback_url': ('string', '回调URL'),
    'department': ('object', '部门'),
    'departments': ('array', '部门列表'),
    'manager': ('object', '管理者'),
    'managers': ('array', '管理者列表'),
    'member': ('object', '成员'),
    'members': ('array', '成员列表'),
    'template': ('object', '模板'),
    'templates': ('array', '模板列表'),
    'language': ('string', '语言'),
    'languages': ('array', '语言列表'),
    'translations': ('object', '翻译'),
    'keys': ('array', '语言键列表'),
    'groups': ('array', '分组列表'),
    'group': ('object', '分组'),
    'logs': ('array', '日志列表'),
    'log': ('object', '日志'),
    'records': ('array', '记录列表'),
    'record': ('object', '记录'),
    'attachments': ('array', '附件列表'),
    'attachment': ('object', '附件'),
    'file': ('object', '文件'),
    'files': ('array', '文件列表'),
    'image': ('string', '图片'),
    'images': ('array', '图片列表'),
    'message': ('string', '消息'),
    'messages': ('array', '消息列表'),
    'unread_count': ('int', '未读数量'),
    'contact': ('object', '联系人'),
    'contacts': ('array', '联系人列表'),
    'promo_code': ('object', '优惠码'),
    'promo_codes': ('array', '优惠码列表'),
    'suffix': ('object', '邮箱后缀'),
    'suffixes': ('array', '邮箱后缀列表'),
    'nav_group': ('object', '导航分组'),
    'nav_groups': ('array', '导航分组列表'),
    'friendly_link': ('object', '友情链接'),
    'friendly_links': ('array', '友情链接列表'),
    'role': ('object', '角色'),
    'roles': ('array', '角色列表'),
    'permission': ('object', '权限'),
    'permissions': ('array', '权限列表'),
    'admin': ('object', '管理员'),
    'admins': ('array', '管理员列表'),
    'login_log': ('object', '登录日志'),
    'login_logs': ('array', '登录日志列表'),
    'operation_log': ('object', '操作日志'),
    'operation_logs': ('array', '操作日志列表'),
    'system_log': ('object', '系统日志'),
    'system_logs': ('array', '系统日志列表'),
    'upgrade': ('object', '升级'),
    'upgrades': ('array', '升级列表'),
    'sign': ('object', '签署信息'),
    'cancel_request': ('object', '注销申请'),
    'cancel_requests': ('array', '注销申请列表'),
    'blacklist': ('object', '黑名单'),
    'blacklists': ('array', '黑名单列表'),
    'client_group': ('object', '客户分组'),
    'client_groups': ('array', '客户分组列表'),
    'client_care': ('object', '客户关怀'),
    'client_cares': ('array', '客户关怀列表'),
    'client_track': ('object', '客户跟踪'),
    'client_tracks': ('array', '客户跟踪列表'),
    'email_suffix': ('object', '邮箱后缀'),
    'email_suffixes': ('array', '邮箱后缀列表'),
    'email_template': ('object', '邮件模板'),
    'email_templates': ('array', '邮件模板列表'),
    'sms_template': ('object', '短信模板'),
    'sms_templates': ('array', '短信模板列表'),
    'download': ('object', '下载'),
    'downloads': ('array', '下载列表'),
    'download_cate': ('object', '下载分类'),
    'download_cates': ('array', '下载分类列表'),
    'payment_method': ('object', '支付方式'),
    'payment_methods': ('array', '支付方式列表'),
    'gateway': ('object', '网关'),
    'gateways': ('array', '网关列表'),
    'oauth_provider': ('object', 'OAuth提供商'),
    'oauth_providers': ('array', 'OAuth提供商列表'),
    'oauth_account': ('object', 'OAuth账号'),
    'oauth_accounts': ('array', 'OAuth账号列表'),
    'kb': ('object', '知识库'),
    'kbs': ('array', '知识库列表'),
    'kb_category': ('object', '知识库分类'),
    'kb_categories': ('array', '知识库分类列表'),
    'notification': ('object', '通知'),
    'notifications': ('array', '通知列表'),
    'task': ('object', '任务'),
    'tasks': ('array', '任务列表'),
    'cron': ('object', '定时任务'),
    'crons': ('array', '定时任务列表'),
    'health': ('object', '健康状态'),
    'version': ('string', '版本号'),
    'uptime': ('string', '运行时间'),
    'memory': ('object', '内存信息'),
    'cpu_usage': ('number', 'CPU使用率'),
    'disk_usage': ('object', '磁盘使用情况'),
    'network': ('object', '网络信息'),
    'marketplace': ('object', '交易市场'),
    'marketplace_config': ('object', '交易市场配置'),
    'chat_session': ('object', '聊天会话'),
    'chat_sessions': ('array', '聊天会话列表'),
    'chat_message': ('object', '聊天消息'),
    'chat_messages': ('array', '聊天消息列表'),
    'market_listing': ('object', '市场挂售'),
    'market_listings': ('array', '市场挂售列表'),
    'market_order': ('object', '市场订单'),
    'market_orders': ('array', '市场订单列表'),
    'market_chat': ('object', '市场聊天'),
    'market_chats': ('array', '市场聊天列表'),
    'market_transaction': ('object', '市场交易'),
    'market_transactions': ('array', '市场交易列表'),
    'v10cloud': ('object', 'V10云'),
    'v10clouds': ('array', 'V10云列表'),
    'ai_session': ('object', 'AI会话'),
    'ai_sessions': ('array', 'AI会话列表'),
    'shopping': ('object', '购物助手'),
    'shopping_session': ('object', '购物会话'),
    'acfp': ('object', 'ACFP模块'),
    'acfps': ('array', 'ACFP模块列表'),
    'server_module': ('object', '服务器模块'),
    'server_modules': ('array', '服务器模块列表'),
    'server_group': ('object', '服务器分组'),
    'server_groups': ('array', '服务器分组列表'),
    'server_option': ('object', '服务器选项'),
    'server_options': ('array', '服务器选项列表'),
    'client_service': ('object', '客户服务'),
    'client_services': ('array', '客户服务列表'),
    'product_group': ('object', '产品分组'),
    'product_groups': ('array', '产品分组列表'),
    'product_divert': ('object', '产品转移'),
    'product_diverts': ('array', '产品转移列表'),
    'product_account': ('object', '产品账号'),
    'product_accounts': ('array', '产品账号列表'),
    'config_option': ('object', '配置选项'),
    'config_options': ('array', '配置选项列表'),
    'config_option_group': ('object', '配置选项组'),
    'config_option_groups': ('array', '配置选项组列表'),
    'affiliate': ('object', '代理推荐'),
    'affiliates': ('array', '代理推荐列表'),
    'affiliate_ladder': ('object', '代理阶梯'),
    'affiliate_ladders': ('array', '代理阶梯列表'),
    'referral': ('object', '推荐'),
    'referrals': ('array', '推荐列表'),
    'commission': ('object', '佣金'),
    'commissions': ('array', '佣金列表'),
    'earnings': ('object', '收益'),
    'earning': ('object', '收益'),
    'sale_record': ('object', '销售记录'),
    'sale_records': ('array', '销售记录列表'),
    'sale_user': ('object', '销售用户'),
    'sale_users': ('array', '销售用户列表'),
    'invoice_item': ('object', '账单条目'),
    'invoice_items': ('array', '账单条目列表'),
    'payment': ('object', '支付记录'),
    'payments': ('array', '支付记录列表'),
    'refund': ('object', '退款'),
    'refunds': ('array', '退款列表'),
    'credit_limit': ('object', '信用额度'),
    'credit_limits': ('array', '信用额度列表'),
    'credit_log': ('object', '信用额度日志'),
    'credit_logs': ('array', '信用额度日志列表'),
    'balance_log': ('object', '余额日志'),
    'balance_logs': ('array', '余额日志列表'),
    'recharge': ('object', '充值'),
    'recharges': ('array', '充值列表'),
    'withdraw': ('object', '提现'),
    'withdraws': ('array', '提现列表'),
    'withdraw_method': ('object', '提现方式'),
    'withdraw_methods': ('array', '提现方式列表'),
    'contract': ('object', '合同'),
    'contracts': ('array', '合同列表'),
    'sign': ('object', '签署'),
    'signs': ('array', '签署列表'),
    'ssl_cert': ('object', 'SSL证书'),
    'ssl_certs': ('array', 'SSL证书列表'),
    'ssl_order': ('object', 'SSL订单'),
    'ssl_orders': ('array', 'SSL订单列表'),
    'captcha': ('object', '验证码'),
    'captchas': ('array', '验证码列表'),
    'geetest': ('object', '极验验证'),
    'geetest_config': ('object', '极验配置'),
    'chat': ('object', '聊天'),
    'chats': ('array', '聊天列表'),
    'chat_message': ('object', '聊天消息'),
    'chat_messages': ('array', '聊天消息列表'),
    'shopping': ('object', '购物助手'),
    'shopping_session': ('object', '购物会话'),
    'shopping_sessions': ('array', '购物会话列表'),
    'shopping_message': ('object', '购物消息'),
    'shopping_messages': ('array', '购物消息列表'),
    'v10cloud': ('object', 'V10云'),
    'v10clouds': ('array', 'V10云列表'),
    'ai_session': ('object', 'AI会话'),
    'ai_sessions': ('array', 'AI会话列表'),
    'acfp': ('object', 'ACFP模块'),
    'acfps': ('array', 'ACFP模块列表'),
    'server_module': ('object', '服务器模块'),
    'server_modules': ('array', '服务器模块列表'),
    'server_group': ('object', '服务器分组'),
    'server_groups': ('array', '服务器分组列表'),
    'server_option': ('object', '服务器选项'),
    'server_options': ('array', '服务器选项列表'),
    'client_service': ('object', '客户服务'),
    'client_services': ('array', '客户服务列表'),
    'product_group': ('object', '产品分组'),
    'product_groups': ('array', '产品分组列表'),
    'product_divert': ('object', '产品转移'),
    'product_diverts': ('array', '产品转移列表'),
    'product_account': ('object', '产品账号'),
    'product_accounts': ('array', '产品账号列表'),
    'config_option': ('object', '配置选项'),
    'config_options': ('array', '配置选项列表'),
    'config_option_group': ('object', '配置选项组'),
    'config_option_groups': ('array', '配置选项组列表'),
    'affiliate': ('object', '代理推荐'),
    'affiliates': ('array', '代理推荐列表'),
    'affiliate_ladder': ('object', '代理阶梯'),
    'affiliate_ladders': ('array', '代理阶梯列表'),
    'referral': ('object', '推荐'),
    'referrals': ('array', '推荐列表'),
    'commission': ('object', '佣金'),
    'commissions': ('array', '佣金列表'),
    'earnings': ('object', '收益'),
    'earning': ('object', '收益'),
    'sale_record': ('object', '销售记录'),
    'sale_records': ('array', '销售记录列表'),
    'sale_user': ('object', '销售用户'),
    'sale_users': ('array', '销售用户列表'),
    'invoice_item': ('object', '账单条目'),
    'invoice_items': ('array', '账单条目列表'),
    'payment': ('object', '支付记录'),
    'payments': ('array', '支付记录列表'),
    'refund': ('object', '退款'),
    'refunds': ('array', '退款列表'),
    'credit_limit': ('object', '信用额度'),
    'credit_limits': ('array', '信用额度列表'),
    'credit_log': ('object', '信用额度日志'),
    'credit_logs': ('array', '信用额度日志列表'),
    'balance_log': ('object', '余额日志'),
    'balance_logs': ('array', '余额日志列表'),
    'recharge': ('object', '充值'),
    'recharges': ('array', '充值列表'),
    'withdraw': ('object', '提现'),
    'withdraws': ('array', '提现列表'),
    'withdraw_method': ('object', '提现方式'),
    'withdraw_methods': ('array', '提现方式列表'),
    'contract': ('object', '合同'),
    'contracts': ('array', '合同列表'),
    'sign': ('object', '签署'),
    'signs': ('array', '签署列表'),
    'ssl_cert': ('object', 'SSL证书'),
    'ssl_certs': ('array', 'SSL证书列表'),
    'ssl_order': ('object', 'SSL订单'),
    'ssl_orders': ('array', 'SSL订单列表'),
    'captcha': ('object', '验证码'),
    'captchas': ('array', '验证码列表'),
    'geetest': ('object', '极验验证'),
    'geetest_config': ('object', '极验配置'),
    'chat': ('object', '聊天'),
    'chats': ('array', '聊天列表'),
    'chat_message': ('object', '聊天消息'),
    'chat_messages': ('array', '聊天消息列表'),
    'shopping': ('object', '购物助手'),
    'shopping_session': ('object', '购物会话'),
    'shopping_sessions': ('array', '购物会话列表'),
    'shopping_message': ('object', '购物消息'),
    'shopping_messages': ('array', '购物消息列表'),
    'v10cloud': ('object', 'V10云'),
    'v10clouds': ('array', 'V10云列表'),
    'ai_session': ('object', 'AI会话'),
    'ai_sessions': ('array', 'AI会话列表'),
    'acfp': ('object', 'ACFP模块'),
    'acfps': ('array', 'ACFP模块列表'),
    'server_module': ('object', '服务器模块'),
    'server_modules': ('array', '服务器模块列表'),
    'server_group': ('object', '服务器分组'),
    'server_groups': ('array', '服务器分组列表'),
    'server_option': ('object', '服务器选项'),
    'server_options': ('array', '服务器选项列表'),
    'client_service': ('object', '客户服务'),
    'client_services': ('array', '客户服务列表'),
    'product_group': ('object', '产品分组'),
    'product_groups': ('array', '产品分组列表'),
    'product_divert': ('object', '产品转移'),
    'product_diverts': ('array', '产品转移列表'),
    'product_account': ('object', '产品账号'),
    'product_accounts': ('array', '产品账号列表'),
    'config_option': ('object', '配置选项'),
    'config_options': ('array', '配置选项列表'),
    'config_option_group': ('object', '配置选项组'),
    'config_option_groups': ('array', '配置选项组列表'),
    'affiliate': ('object', '代理推荐'),
    'affiliates': ('array', '代理推荐列表'),
    'affiliate_ladder': ('object', '代理阶梯'),
    'affiliate_ladders': ('array', '代理阶梯列表'),
    'referral': ('object', '推荐'),
    'referrals': ('array', '推荐列表'),
    'commission': ('object', '佣金'),
    'commissions': ('array', '佣金列表'),
    'earnings': ('object', '收益'),
    'earning': ('object', '收益'),
    'sale_record': ('object', '销售记录'),
    'sale_records': ('array', '销售记录列表'),
    'sale_user': ('object', '销售用户'),
    'sale_users': ('array', '销售用户列表'),
    'invoice_item': ('object', '账单条目'),
    'invoice_items': ('array', '账单条目列表'),
    'payment': ('object', '支付记录'),
    'payments': ('array', '支付记录列表'),
    'refund': ('object', '退款'),
    'refunds': ('array', '退款列表'),
    'credit_limit': ('object', '信用额度'),
    'credit_limits': ('array', '信用额度列表'),
    'credit_log': ('object', '信用额度日志'),
    'credit_logs': ('array', '信用额度日志列表'),
    'balance_log': ('object', '余额日志'),
    'balance_logs': ('array', '余额日志列表'),
    'recharge': ('object', '充值'),
    'recharges': ('array', '充值列表'),
    'withdraw': ('object', '提现'),
    'withdraws': ('array', '提现列表'),
    'withdraw_method': ('object', '提现方式'),
    'withdraw_methods': ('array', '提现方式列表'),
    'contract': ('object', '合同'),
    'contracts': ('array', '合同列表'),
    'sign': ('object', '签署'),
    'signs': ('array', '签署列表'),
    'ssl_cert': ('object', 'SSL证书'),
    'ssl_certs': ('array', 'SSL证书列表'),
    'ssl_order': ('object', 'SSL订单'),
    'ssl_orders': ('array', 'SSL订单列表'),
    'captcha': ('object', '验证码'),
    'captchas': ('array', '验证码列表'),
    'geetest': ('object', '极验验证'),
    'geetest_config': ('object', '极验配置'),
    'chat': ('object', '聊天'),
    'chats': ('array', '聊天列表'),
    'chat_message': ('object', '聊天消息'),
    'chat_messages': ('array', '聊天消息列表'),
    'shopping': ('object', '购物助手'),
    'shopping_session': ('object', '购物会话'),
    'shopping_sessions': ('array', '购物会话列表'),
    'shopping_message': ('object', '购物消息'),
    'shopping_messages': ('array', '购物消息列表'),
    'v10cloud': ('object', 'V10云'),
    'v10clouds': ('array', 'V10云列表'),
    'ai_session': ('object', 'AI会话'),
    'ai_sessions': ('array', 'AI会话列表'),
    'acfp': ('object', 'ACFP模块'),
    'acfps': ('array', 'ACFP模块列表'),
}

def extract_model_fields_from_dir():
    """Extract all model struct fields from model directory."""
    for fpath in sorted(glob.glob(os.path.join(MODEL_DIR, '*.go'))):
        with open(fpath, 'r', encoding='utf-8') as f:
            content = f.read()
        extract_model_fields(content, os.path.basename(fpath))

def extract_responses_enhanced(content):
    """Enhanced response extraction that also traces variables."""
    func_responses = {}
    func_positions = []
    for m in re.finditer(r'func\s+\([^)]+\)\s+(\w+)\([^)]*\*gin\.Context[^)]*\)\s*\{', content):
        func_positions.append((m.end(), m.group(1)))
    
    for i, (pos, fn) in enumerate(func_positions):
        end_pos = func_positions[i+1][0] if i+1 < len(func_positions) else min(pos+10000, len(content))
        chunk = content[pos:end_pos]
        resp_fields = []
        
        # Pattern 1: response.Success(c, gin.H{...})
        for rm in re.finditer(r'response\.Success\(\s*c\s*,\s*gin\.H\{([^}]+)\}', chunk):
            body = rm.group(1)
            for fm in re.finditer(r'"(\w+)":', body):
                fname = fm.group(1)
                t, d = COMMON_GIN_H.get(fname, ('object', fname))
                resp_fields.append({'n':f'data.{fname}','t':t,'d':d})
        
        # Pattern 2: response.Success(c, variable) - trace variable
        for rm in re.finditer(r'response\.Success\(\s*c\s*,\s*(\w+)\s*\)', chunk):
            vn = rm.group(1)
            if vn in ('c', 'gin', 'nil'): continue
            
            # Check if variable name maps to known model
            if vn in VAR_TO_MODEL:
                model_name = VAR_TO_MODEL[vn]
                if model_name in MODEL_FIELDS:
                    for f in MODEL_FIELDS[model_name][:15]:  # Limit to 15 fields
                        resp_fields.append({'n':f'data.{f["n"]}','t':f['t'],'d':f['n']})
                    continue
            
            # Check COMMON_GIN_H
            if vn in COMMON_GIN_H:
                t, d = COMMON_GIN_H[vn]
                resp_fields.append({'n':f'data.{vn}','t':t,'d':d})
            else:
                resp_fields.append({'n':'data','t':'object','d':'返回数据'})
        
        # Pattern 3: response.SuccessWithList
        if re.search(r'response\.SuccessWithList\(', chunk):
            resp_fields.append({'n':'data.list','t':'array','d':'列表'})
            resp_fields.append({'n':'data.total','t':'int','d':'总数'})
        
        # Pattern 4: response.SuccessWithPagination
        if re.search(r'response\.SuccessWithPagination\(', chunk):
            resp_fields.append({'n':'data.list','t':'array','d':'列表'})
            resp_fields.append({'n':'data.total','t':'int','d':'总数'})
            resp_fields.append({'n':'data.page','t':'int','d':'页码'})
            resp_fields.append({'n':'data.page_size','t':'int','d':'每页数量'})
        
        if re.search(r'response\.SuccessMsg\(', chunk) and not resp_fields:
            resp_fields.append({'n':'message','t':'string','d':'提示消息'})
        
        if resp_fields:
            func_responses[fn] = resp_fields
    
    return func_responses

def main():
    # Load model fields
    extract_model_fields_from_dir()
    print(f"Models: {len(MODEL_FIELDS)}")
    
    # Re-extract from handlers with enhanced logic
    all_structs = {}
    all_query_params = {}
    all_responses = {}
    
    from build_api_v3 import extract_structs, extract_query_params as extract_qp_basic, extract_routes, infer_title, path_to_id, NAMES, PD as PD3, EX as EX3
    
    for fpath in sorted(glob.glob(os.path.join(HANDLER_DIR, '*.go'))):
        with open(fpath, 'r', encoding='utf-8') as f:
            content = f.read()
        all_structs.update(extract_structs(content))
        all_query_params.update(extract_qp_basic(content))
        # Use enhanced response extraction
        resp = extract_responses_enhanced(content)
        all_responses.update(resp)
    
    print(f"Structs: {len(all_structs)}, Query: {len(all_query_params)}, Resp: {len(all_responses)}")
    
    # Build routes with enhanced data
    all_routes = []
    for subdir in ['v1', 'admin']:
        rpath = os.path.join(API_DIR, subdir, 'router.go')
        if os.path.exists(rpath):
            routes = extract_routes(rpath)
            for r in routes:
                mn = r.get('method_name', '')
                body_p = r.get('req_fields', all_structs.get(mn, []))
                query_p = r.get('query_params', all_query_params.get(mn, []))
                resp_f = r.get('resp_fields', all_responses.get(mn, []))
                
                seen = set()
                combined = []
                for p in body_p + query_p:
                    if p['n'] not in seen:
                        seen.add(p['n'])
                        combined.append(p)
                
                # Add path params
                for m in re.finditer(r':(\w+)', r['path']):
                    pname = m.group(1)
                    if pname not in seen:
                        ptype = 'int' if 'id' in pname else 'string'
                        from build_api_v3 import PD, EX
                        combined.append({'n':pname,'t':ptype,'r':'必填','d':PD.get(pname,pname),'e':EX.get(pname,'1')})
                        seen.add(pname)
                
                # Add common query params
                if r['method'] == 'GET':
                    path_lower = r['path'].lower()
                    is_list = any(x in path_lower for x in ['/list', '/logs', '/records', '/messages'])
                    if is_list:
                        if 'page' not in seen:
                            combined.append({'n':'page','t':'int','r':'选填','d':'页码','e':'1'})
                        if 'page_size' not in seen:
                            combined.append({'n':'page_size','t':'int','r':'选填','d':'每页数量','e':'20'})
                
                r['req_fields'] = combined
                r['resp_fields'] = resp_f
            all_routes.extend(routes)
    
    pub = [r for r in all_routes if r['path'].startswith('/api/v1') and r['group'] != 'user']
    usr = [r for r in all_routes if r['group'] == 'user']
    adm = [r for r in all_routes if r['path'].startswith('/api/admin')]
    
    def grp(routes):
        mods = {}
        for r in routes:
            p = re.sub(r'^/api/(v1|admin|public)/?', '', r['path']).strip('/')
            mod = p.split('/')[0] if p else 'other'
            if mod not in mods: mods[mod] = []
            mods[mod].append(r)
        return mods
    
    def item(r):
        from build_api_v3 import PD, EX
        rid = path_to_id(r['path'], r['method'])
        title = infer_title(r['path'], r['method'])
        url, method = r['path'], r['method']
        req_p = r.get('req_fields', [])
        req_s = ','.join([f"{{n:'{p['n']}',t:'{p['t']}',r:'{p['r']}',d:'{p['d']}',e:'{p.get('e','')}'}}" for p in req_p])
        
        resp_p = r.get('resp_fields', [{'n':'data','t':'object','d':'返回数据'}])
        res_parts = []
        for p in resp_p:
            pt = p.get('t', 'object')
            res_parts.append("{n:'" + p['n'] + "',t:'" + pt + "',r:'',d:'" + p['d'] + "',e:''}")
        res_s = ','.join(res_parts)
        
        body_p = {p['n']:p['e'] for p in req_p if p.get('src') not in ('q','p') and p.get('e')}
        query_p = {p['n']:p['e'] for p in req_p if p.get('src')=='q' and p.get('e')}
        
        actual_url = url
        for m in re.finditer(r':(\w+)', url):
            pname = m.group(1)
            example = EX.get(pname, '1')
            actual_url = actual_url.replace(f':{pname}', example)
        
        if method in ('POST','PUT','PATCH') and body_p:
            body = json.dumps(body_p, ensure_ascii=False)
            req_ex = f"{method} {actual_url}"
            if query_p: req_ex += '?' + '&'.join([f"{k}={v}" for k,v in query_p.items()])
            req_ex += f"\\nAuthorization: JWT eyJ...\\n\\n{body}"
        elif query_p:
            qs = '&'.join([f"{k}={v}" for k,v in query_p.items()])
            req_ex = f"{method} {actual_url}?{qs}\\nAuthorization: JWT eyJ..."
        else:
            auth = any(x in url for x in ['/user/','/hosts','/orders','/invoices','/tickets','/cart','/balance','/credit','/contacts','/messages','/certification','/marketplace','/product-diverts','/affiliate','/contracts','/vouchers','/upgrades','/oauth','/ssl','/ai-shopping','/v10cloud','/multi-renew','/api-logs','/login-logs','/admin'])
            req_ex = f"{method} {actual_url}\\nAuthorization: JWT eyJ..." if auth else f"{method} {actual_url}"
        
        if resp_p and len(resp_p) > 1:
            data = {}
            for p in resp_p:
                k = p['n'].replace('data.','')
                if '.' not in k:
                    t = p.get('t','object')
                    if t == 'array': data[k] = '...'
                    elif t == 'int': data[k] = 0
                    elif t == 'string': data[k] = ''
                    elif t == 'number': data[k] = 0.0
                    elif t == 'bool': data[k] = True
                    else: data[k] = '{}'
            res_ex = json.dumps({'code':0,'data':data,'message':'success'}, ensure_ascii=False)
        else:
            res_ex = '{"code":0,"data":{},"message":"success"}'
        
        return (f"        {{id:'{rid}',title:'{title}',method:'{method}',url:'{url}',"
                f"desc:'{title}',reqParams:[{req_s}],resParams:[{res_s}],"
                f"reqExample:'{req_ex}',resExample:'{res_ex}'}}")
    
    out = ['var API_DATA = {']
    for sn, sr in [('public',pub),('user',usr),('admin',adm)]:
        t = {'public':'公共接口','user':'用户接口','admin':'管理接口'}[sn]
        out.append(f'  {sn}: {{ title: "{t}", groups: [')
        for mn, mr in grp(sr).items():
            out.append(f"      {{ name:'{mn}', items:[")
            for r in mr: out.append(item(r) + ',')
            out.append("      ]},")
        out.append('  ]},')
    out.append('};')
    
    js = '\n'.join(out)
    with open(os.path.join(DOCS_DIR, 'api_data_complete.js'), 'w', encoding='utf-8') as f:
        f.write(js)
    
    total = len(all_routes)
    wr = sum(1 for r in all_routes if r.get('req_fields'))
    wrs = sum(1 for r in all_routes if r.get('resp_fields'))
    print(f"Total: {total}, Req: {wr} ({wr*100//total}%), Resp: {wrs} ({wrs*100//total}%)")

if __name__ == '__main__':
    main()
