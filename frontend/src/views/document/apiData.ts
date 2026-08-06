// API Documentation Data
// This file is auto-generated from docs/index.html
export interface ApiParam {
  n: string  // name
  t: string  // type
  r: string  // required
  d: string  // description
  e: string  // example
}

export interface ApiItem {
  id: string
  title: string
  method: string
  url: string
  desc: string
  reqParams: ApiParam[]
  resParams: ApiParam[]
  reqExample: string
  resExample: string
}

export interface ApiGroup {
  name: string
  items: ApiItem[]
}

export interface ApiPage {
  title: string
  groups: ApiGroup[]
}

export interface ApiDocData {
  public: ApiPage
  user: ApiPage
  admin: ApiPage
}

export const API_DATA: ApiDocData = {
  public: { title: "公共接口", groups: [
    { name: '认证', items: [
      { id: 'post-auth-login', title: '登录', method: 'POST', url: '/api/v1/auth/login', desc: '用户登录', reqParams: [{ n: 'account', t: 'string', r: '必填', d: '用户名/邮箱/手机号', e: 'admin' }, { n: 'password', t: 'string', r: '必填', d: '密码', e: 'Abc123456' }, { n: 'captcha_id', t: 'string', r: '选填', d: '验证码ID', e: 'abc123' }, { n: 'answer', t: 'string', r: '选填', d: '验证码答案', e: 'x5km' }], resParams: [{ n: 'data.access_token', t: 'string', r: '', d: '访问令牌', e: '' }, { n: 'data.refresh_token', t: 'string', r: '', d: '刷新令牌', e: '' }, { n: 'data.user', t: 'object', r: '', d: '用户信息', e: '' }], reqExample: 'POST /api/v1/auth/login\nContent-Type: application/json\n\n{"account":"admin","password":"Abc123456"}', resExample: '{"code":0,"data":{"access_token":"eyJ...","refresh_token":"eyJ...","user":{"id":1,"username":"admin"}},"message":"success"}' },
      { id: 'post-auth-sms-login', title: '短信登录', method: 'POST', url: '/api/v1/auth/sms-login', desc: '短信验证码登录', reqParams: [{ n: 'phone', t: 'string', r: '必填', d: '手机号', e: '13800138000' }, { n: 'code', t: 'string', r: '必填', d: '验证码', e: '123456' }], resParams: [{ n: 'data.access_token', t: 'string', r: '', d: '访问令牌', e: '' }, { n: 'data.refresh_token', t: 'string', r: '', d: '刷新令牌', e: '' }, { n: 'data.user', t: 'object', r: '', d: '用户信息', e: '' }], reqExample: 'POST /api/v1/auth/sms-login\nContent-Type: application/json\n\n{"phone":"13800138000","code":"123456"}', resExample: '{"code":0,"data":{"access_token":"eyJ...","refresh_token":"eyJ...","user":{"id":1}},"message":"success"}' },
      { id: 'post-auth-register', title: '注册', method: 'POST', url: '/api/v1/auth/register', desc: '用户注册', reqParams: [{ n: 'username', t: 'string', r: '必填', d: '用户名', e: 'testuser' }, { n: 'password', t: 'string', r: '必填', d: '密码', e: 'Abc123456' }, { n: 'email', t: 'string', r: '选填', d: '邮箱', e: 'user@example.com' }, { n: 'phone', t: 'string', r: '选填', d: '手机号', e: '13800138000' }], resParams: [{ n: 'data.access_token', t: 'string', r: '', d: '访问令牌', e: '' }, { n: 'data.user', t: 'object', r: '', d: '用户信息', e: '' }], reqExample: 'POST /api/v1/auth/register\nContent-Type: application/json\n\n{"username":"testuser","password":"Abc123456","email":"user@example.com"}', resExample: '{"code":0,"data":{"access_token":"eyJ...","user":{"id":2}},"message":"success"}' },
      { id: 'post-auth-refresh', title: '刷新令牌', method: 'POST', url: '/api/v1/auth/refresh', desc: '刷新JWT令牌', reqParams: [{ n: 'refresh_token', t: 'string', r: '必填', d: '刷新令牌', e: 'eyJ...' }], resParams: [{ n: 'data.access_token', t: 'string', r: '', d: '新访问令牌', e: '' }, { n: 'data.refresh_token', t: 'string', r: '', d: '新刷新令牌', e: '' }], reqExample: 'POST /api/v1/auth/refresh\nContent-Type: application/json\n\n{"refresh_token":"eyJ..."}', resExample: '{"code":0,"data":{"access_token":"eyJ...","refresh_token":"eyJ..."},"message":"success"}' },
      { id: 'post-auth-logout', title: '退出登录', method: 'POST', url: '/api/v1/auth/logout', desc: '退出登录', reqParams: [], resParams: [{ n: 'message', t: 'string', r: '', d: '提示消息', e: '' }], reqExample: 'POST /api/v1/auth/logout\nAuthorization: JWT eyJ...', resExample: '{"code":0,"data":null,"message":"success"}' },
      { id: 'post-auth-access-token', title: 'API密钥登录', method: 'POST', url: '/api/v1/auth/access-token', desc: '使用API密钥登录', reqParams: [{ n: 'api_key', t: 'string', r: '必填', d: 'API密钥', e: 'aB3$xY7!kL9m' }], resParams: [{ n: 'data.access_token', t: 'string', r: '', d: '访问令牌', e: '' }, { n: 'data.user', t: 'object', r: '', d: '用户信息', e: '' }], reqExample: 'POST /api/v1/auth/access-token\nContent-Type: application/json\n\n{"api_key":"aB3$xY7!kL9m"}', resExample: '{"code":0,"data":{"access_token":"eyJ...","user":{"id":1}},"message":"success"}' },
    ]},
    { name: '验证码', items: [
      { id: 'get-captcha-generate', title: '获取验证码', method: 'GET', url: '/api/v1/captcha/generate', desc: '获取图形验证码', reqParams: [], resParams: [{ n: 'data.captcha_id', t: 'string', r: '', d: '验证码ID', e: '' }, { n: 'data.image', t: 'string', r: '', d: '验证码图片Base64', e: '' }], reqExample: 'GET /api/v1/captcha/generate', resExample: '{"code":0,"data":{"captcha_id":"abc123","image":"data:image/png;base64,..."},"message":"success"}' },
      { id: 'post-captcha-verify', title: '验证码校验', method: 'POST', url: '/api/v1/captcha/verify', desc: '校验图形验证码', reqParams: [{ n: 'captcha_id', t: 'string', r: '必填', d: '验证码ID', e: 'abc123' }, { n: 'answer', t: 'string', r: '必填', d: '验证码答案', e: 'x5km' }], resParams: [{ n: 'data.success', t: 'bool', r: '', d: '验证结果', e: '' }], reqExample: 'POST /api/v1/captcha/verify\nContent-Type: application/json\n\n{"captcha_id":"abc123","answer":"x5km"}', resExample: '{"code":0,"data":{"success":true},"message":"success"}' },
      { id: 'get-payment-methods', title: '获取支付方式', method: 'GET', url: '/api/v1/payment-methods', desc: '获取可用支付方式列表', reqParams: [], resParams: [{ n: 'data', t: 'array', r: '', d: '支付方式列表', e: '' }], reqExample: 'GET /api/v1/payment-methods', resExample: '{"code":0,"data":[{"id":1,"name":"支付宝","code":"alipay"}],"message":"success"}' },
    ]},
    { name: '公共内容', items: [
      { id: 'get-products', title: '商品列表', method: 'GET', url: '/api/v1/products', desc: '获取商品列表', reqParams: [{ n: 'group_id', t: 'int', r: '选填', d: '分组ID', e: '1' }, { n: 'keyword', t: 'string', r: '选填', d: '关键词', e: '云服务器' }], resParams: [{ n: 'data', t: 'array', r: '', d: '商品列表', e: '' }], reqExample: 'GET /api/v1/products?group_id=1', resExample: '{"code":0,"data":[{"id":1,"name":"基础型云服务器","price":"99.00"}],"message":"success"}' },
      { id: 'get-product-detail', title: '商品详情', method: 'GET', url: '/api/v1/products/:id', desc: '获取商品详情', reqParams: [{ n: 'id', t: 'int', r: '必填', d: '商品ID', e: '1' }], resParams: [{ n: 'data', t: 'object', r: '', d: '商品详情', e: '' }], reqExample: 'GET /api/v1/products/1', resExample: '{"code":0,"data":{"id":1,"name":"基础型云服务器","price":"99.00"},"message":"success"}' },
      { id: 'get-announcements', title: '公告列表', method: 'GET', url: '/api/v1/announcements', desc: '获取公告列表', reqParams: [{ n: 'page', t: 'int', r: '选填', d: '页码', e: '1' }, { n: 'page_size', t: 'int', r: '选填', d: '每页数量', e: '10' }], resParams: [{ n: 'data.list', t: 'array', r: '', d: '公告列表', e: '' }, { n: 'data.total', t: 'int', r: '', d: '总数', e: '' }], reqExample: 'GET /api/v1/announcements?page=1', resExample: '{"code":0,"data":{"list":[{"id":1,"title":"系统维护通知"}],"total":1},"message":"success"}' },
      { id: 'get-news', title: '新闻中心', method: 'GET', url: '/api/v1/news', desc: '获取新闻中心列表', reqParams: [{ n: 'page', t: 'int', r: '选填', d: '页码', e: '1' }], resParams: [{ n: 'data.list', t: 'array', r: '', d: '新闻列表', e: '' }], reqExample: 'GET /api/v1/news?page=1', resExample: '{"code":0,"data":{"list":[{"id":1,"title":"..."}],"total":1},"message":"success"}' },
      { id: 'get-help', title: '帮助中心', method: 'GET', url: '/api/v1/help', desc: '获取帮助中心', reqParams: [{ n: 'keyword', t: 'string', r: '选填', d: '搜索关键词', e: '域名' }], resParams: [{ n: 'data.categories', t: 'array', r: '', d: '分类列表', e: '' }, { n: 'data.articles', t: 'array', r: '', d: '文章列表', e: '' }], reqExample: 'GET /api/v1/help?keyword=域名', resExample: '{"code":0,"data":{"categories":[],"articles":[]},"message":"success"}' },
      { id: 'get-homepage-banners', title: '轮播图', method: 'GET', url: '/api/v1/homepage/banners', desc: '获取首页轮播图', reqParams: [], resParams: [{ n: 'data', t: 'array', r: '', d: '轮播图列表', e: '' }], reqExample: 'GET /api/v1/homepage/banners', resExample: '{"code":0,"data":[{"id":1,"title":"高性能云服务器","image_url":"/carousel/1.jpg"}],"message":"success"}' },
      { id: 'get-homepage-nav', title: '导航菜单', method: 'GET', url: '/api/v1/homepage/nav/:position', desc: '获取顶部/底部导航', reqParams: [{ n: 'position', t: 'string', r: '必填', d: '位置(top/bottom)', e: 'top' }], resParams: [{ n: 'data', t: 'array', r: '', d: '菜单列表', e: '' }], reqExample: 'GET /api/v1/homepage/nav/top', resExample: '{"code":0,"data":[{"id":1,"title":"首页","url":"/"}],"message":"success"}' },
      { id: 'get-settings-public', title: '公开设置', method: 'GET', url: '/api/v1/settings/public', desc: '获取公开站点设置', reqParams: [], resParams: [{ n: 'data.site_name', t: 'string', r: '', d: '站点名称', e: '' }, { n: 'data.site_logo', t: 'string', r: '', d: '站点Logo', e: '' }], reqExample: 'GET /api/v1/settings/public', resExample: '{"code":0,"data":{"site_name":"锚点财务","site_logo":"/logo.png"},"message":"success"}' },
    ]},
  ]},
  user: { title: "用户接口", groups: [
    { name: '个人资料', items: [
      { id: 'get-user-profile', title: '获取个人资料', method: 'GET', url: '/api/v1/user/profile', desc: '获取当前用户资料', reqParams: [], resParams: [{ n: 'data', t: 'object', r: '', d: '用户资料', e: '' }], reqExample: 'GET /api/v1/user/profile\nAuthorization: JWT eyJ...', resExample: '{"code":0,"data":{"id":1,"username":"admin","email":"user@example.com"},"message":"success"}' },
      { id: 'put-user-profile', title: '修改个人资料', method: 'PUT', url: '/api/v1/user/profile', desc: '修改个人资料', reqParams: [{ n: 'nickname', t: 'string', r: '选填', d: '昵称', e: '新昵称' }, { n: 'avatar', t: 'string', r: '选填', d: '头像URL', e: '/uploads/avatar.jpg' }, { n: 'qq', t: 'string', r: '选填', d: 'QQ号', e: '12345678' }], resParams: [{ n: 'message', t: 'string', r: '', d: '提示消息', e: '' }], reqExample: 'PUT /api/v1/user/profile\nAuthorization: JWT eyJ...\n\n{"nickname":"新昵称"}', resExample: '{"code":0,"data":null,"message":"修改成功"}' },
      { id: 'put-user-password', title: '修改密码', method: 'PUT', url: '/api/v1/user/password', desc: '修改登录密码', reqParams: [{ n: 'old_password', t: 'string', r: '必填', d: '旧密码', e: 'Abc123456' }, { n: 'new_password', t: 'string', r: '必填', d: '新密码', e: 'NewPass789' }], resParams: [{ n: 'message', t: 'string', r: '', d: '提示消息', e: '' }], reqExample: 'PUT /api/v1/user/password\nAuthorization: JWT eyJ...\n\n{"old_password":"Abc123456","new_password":"NewPass789"}', resExample: '{"code":0,"data":null,"message":"密码修改成功"}' },
    ]},
    { name: '产品/主机', items: [
      { id: 'get-user-hosts', title: '产品列表', method: 'GET', url: '/api/v1/user/hosts', desc: '获取已购产品列表', reqParams: [{ n: 'page', t: 'int', r: '选填', d: '页码', e: '1' }, { n: 'status', t: 'string', r: '选填', d: '状态(active/suspended)', e: 'active' }], resParams: [{ n: 'data.list', t: 'array', r: '', d: '产品列表', e: '' }, { n: 'data.total', t: 'int', r: '', d: '总数', e: '' }], reqExample: 'GET /api/v1/user/hosts?page=1&status=active\nAuthorization: JWT eyJ...', resExample: '{"code":0,"data":{"list":[{"id":1,"product_name":"云服务器","hostname":"web01","status":"active"}],"total":1},"message":"success"}' },
      { id: 'get-user-host-detail', title: '产品详情', method: 'GET', url: '/api/v1/user/hosts/:id', desc: '获取产品详情', reqParams: [{ n: 'id', t: 'int', r: '必填', d: '主机ID', e: '1' }], resParams: [{ n: 'data', t: 'object', r: '', d: '产品详情', e: '' }], reqExample: 'GET /api/v1/user/hosts/1\nAuthorization: JWT eyJ...', resExample: '{"code":0,"data":{"id":1,"hostname":"web01","ip":"192.168.1.1","os":"CentOS 7"},"message":"success"}' },
      { id: 'post-user-host-start', title: '开机', method: 'POST', url: '/api/v1/user/hosts/:id/start', desc: '开机', reqParams: [{ n: 'id', t: 'int', r: '必填', d: '主机ID', e: '1' }], resParams: [{ n: 'message', t: 'string', r: '', d: '提示消息', e: '' }], reqExample: 'POST /api/v1/user/hosts/1/start\nAuthorization: JWT eyJ...', resExample: '{"code":0,"data":null,"message":"开机成功"}' },
      { id: 'post-user-host-stop', title: '关机', method: 'POST', url: '/api/v1/user/hosts/:id/stop', desc: '关机', reqParams: [{ n: 'id', t: 'int', r: '必填', d: '主机ID', e: '1' }], resParams: [{ n: 'message', t: 'string', r: '', d: '提示消息', e: '' }], reqExample: 'POST /api/v1/user/hosts/1/stop\nAuthorization: JWT eyJ...', resExample: '{"code":0,"data":null,"message":"关机成功"}' },
      { id: 'post-user-host-restart', title: '重启', method: 'POST', url: '/api/v1/user/hosts/:id/restart', desc: '重启', reqParams: [{ n: 'id', t: 'int', r: '必填', d: '主机ID', e: '1' }], resParams: [{ n: 'message', t: 'string', r: '', d: '提示消息', e: '' }], reqExample: 'POST /api/v1/user/hosts/1/restart\nAuthorization: JWT eyJ...', resExample: '{"code":0,"data":null,"message":"重启成功"}' },
      { id: 'post-user-host-reinstall', title: '重装系统', method: 'POST', url: '/api/v1/user/hosts/:id/reinstall', desc: '重装操作系统', reqParams: [{ n: 'id', t: 'int', r: '必填', d: '主机ID', e: '1' }, { n: 'os', t: 'string', r: '必填', d: '操作系统', e: 'CentOS 7' }], resParams: [{ n: 'data.password', t: 'string', r: '', d: 'root密码', e: '' }], reqExample: 'POST /api/v1/user/hosts/1/reinstall\nAuthorization: JWT eyJ...\n\n{"os":"CentOS 7"}', resExample: '{"code":0,"data":{"password":"auto123"},"message":"重装中"}' },
      { id: 'get-user-host-snapshots', title: '快照列表', method: 'GET', url: '/api/v1/user/hosts/:id/snapshots', desc: '获取快照列表', reqParams: [{ n: 'id', t: 'int', r: '必填', d: '主机ID', e: '1' }], resParams: [{ n: 'data', t: 'array', r: '', d: '快照列表', e: '' }], reqExample: 'GET /api/v1/user/hosts/1/snapshots\nAuthorization: JWT eyJ...', resExample: '{"code":0,"data":[{"id":1,"name":"snapshot_20260101"}],"message":"success"}' },
    ]},
    { name: '购物车', items: [
      { id: 'get-cart', title: '获取购物车', method: 'GET', url: '/api/v1/cart', desc: '获取购物车信息', reqParams: [], resParams: [{ n: 'data.items', t: 'array', r: '', d: '购物车商品', e: '' }, { n: 'data.total', t: 'number', r: '', d: '总价', e: '' }], reqExample: 'GET /api/v1/cart\nAuthorization: JWT eyJ...', resExample: '{"code":0,"data":{"items":[],"total":0},"message":"success"}' },
      { id: 'post-cart-add', title: '添加至购物车', method: 'POST', url: '/api/v1/cart/add', desc: '添加商品至购物车', reqParams: [{ n: 'product_id', t: 'int', r: '必填', d: '商品ID', e: '1' }, { n: 'billing_cycle', t: 'string', r: '必填', d: '计费周期', e: 'monthly' }], resParams: [{ n: 'message', t: 'string', r: '', d: '提示消息', e: '' }], reqExample: 'POST /api/v1/cart/add\nAuthorization: JWT eyJ...\n\n{"product_id":1,"billing_cycle":"monthly"}', resExample: '{"code":0,"data":null,"message":"添加成功"}' },
      { id: 'post-cart-checkout', title: '结算', method: 'POST', url: '/api/v1/cart/checkout', desc: '购物车结算', reqParams: [{ n: 'payment_method_id', t: 'int', r: '必填', d: '支付方式ID', e: '1' }], resParams: [{ n: 'data.invoice_id', t: 'int', r: '', d: '账单ID', e: '' }, { n: 'data.total', t: 'number', r: '', d: '应付金额', e: '' }], reqExample: 'POST /api/v1/cart/checkout\nAuthorization: JWT eyJ...\n\n{"payment_method_id":1}', resExample: '{"code":0,"data":{"invoice_id":1,"total":99.00},"message":"success"}' },
    ]},
    { name: '账单', items: [
      { id: 'get-user-invoices', title: '账单列表', method: 'GET', url: '/api/v1/user/invoices', desc: '获取账单列表', reqParams: [{ n: 'page', t: 'int', r: '选填', d: '页码', e: '1' }, { n: 'status', t: 'string', r: '选填', d: '状态(unpaid/paid)', e: 'unpaid' }], resParams: [{ n: 'data.list', t: 'array', r: '', d: '账单列表', e: '' }], reqExample: 'GET /api/v1/user/invoices?page=1&status=unpaid\nAuthorization: JWT eyJ...', resExample: '{"code":0,"data":{"list":[{"id":1,"status":"unpaid","total":99.00}],"total":1},"message":"success"}' },
      { id: 'post-user-invoice-pay', title: '支付账单', method: 'POST', url: '/api/v1/user/invoices/:id/pay', desc: '发起支付', reqParams: [{ n: 'id', t: 'int', r: '必填', d: '账单ID', e: '1' }, { n: 'payment_method_id', t: 'int', r: '必填', d: '支付方式ID', e: '1' }], resParams: [{ n: 'data.pay_url', t: 'string', r: '', d: '支付链接', e: '' }], reqExample: 'POST /api/v1/user/invoices/1/pay\nAuthorization: JWT eyJ...\n\n{"payment_method_id":1}', resExample: '{"code":0,"data":{"pay_url":"https://pay.example.com/..."},"message":"success"}' },
    ]},
    { name: '工单', items: [
      { id: 'get-user-tickets', title: '工单列表', method: 'GET', url: '/api/v1/user/tickets', desc: '获取工单列表', reqParams: [{ n: 'page', t: 'int', r: '选填', d: '页码', e: '1' }, { n: 'status', t: 'string', r: '选填', d: '状态(open/answered/closed)', e: 'open' }], resParams: [{ n: 'data.list', t: 'array', r: '', d: '工单列表', e: '' }], reqExample: 'GET /api/v1/user/tickets?page=1\nAuthorization: JWT eyJ...', resExample: '{"code":0,"data":{"list":[{"id":1,"subject":"服务器无法访问","status":"open"}],"total":1},"message":"success"}' },
      { id: 'post-user-ticket-create', title: '提交工单', method: 'POST', url: '/api/v1/user/tickets', desc: '提交新工单', reqParams: [{ n: 'subject', t: 'string', r: '必填', d: '主题', e: '服务器无法访问' }, { n: 'content', t: 'string', r: '必填', d: '内容', e: '我的服务器无法正常访问' }, { n: 'department_id', t: 'int', r: '必填', d: '部门ID', e: '1' }], resParams: [{ n: 'data.id', t: 'int', r: '', d: '工单ID', e: '' }], reqExample: 'POST /api/v1/user/tickets\nAuthorization: JWT eyJ...\n\n{"subject":"服务器无法访问","content":"...","department_id":1}', resExample: '{"code":0,"data":{"id":1},"message":"工单提交成功"}' },
      { id: 'post-user-ticket-reply', title: '回复工单', method: 'POST', url: '/api/v1/user/tickets/:id/reply', desc: '回复工单', reqParams: [{ n: 'id', t: 'int', r: '必填', d: '工单ID', e: '1' }, { n: 'content', t: 'string', r: '必填', d: '回复内容', e: '已检查' }], resParams: [{ n: 'message', t: 'string', r: '', d: '提示消息', e: '' }], reqExample: 'POST /api/v1/user/tickets/1/reply\nAuthorization: JWT eyJ...\n\n{"content":"已检查"}', resExample: '{"code":0,"data":null,"message":"回复成功"}' },
    ]},
    { name: '消息中心', items: [
      { id: 'get-user-messages', title: '消息列表', method: 'GET', url: '/api/v1/user/messages', desc: '获取消息列表', reqParams: [{ n: 'page', t: 'int', r: '选填', d: '页码', e: '1' }], resParams: [{ n: 'data.list', t: 'array', r: '', d: '消息列表', e: '' }, { n: 'data.unread_count', t: 'int', r: '', d: '未读数量', e: '' }], reqExample: 'GET /api/v1/user/messages?page=1\nAuthorization: JWT eyJ...', resExample: '{"code":0,"data":{"list":[{"id":1,"title":"系统通知","is_read":false}],"unread_count":1},"message":"success"}' },
    ]},
  ]},
  admin: { title: "管理接口", groups: [
    { name: '管理员认证', items: [
      { id: 'post-admin-login', title: '管理员登录', method: 'POST', url: '/api/admin/auth/login', desc: '管理员登录', reqParams: [{ n: 'username', t: 'string', r: '必填', d: '管理员用户名', e: 'admin' }, { n: 'password', t: 'string', r: '必填', d: '密码', e: 'Abc123456' }], resParams: [{ n: 'data.access_token', t: 'string', r: '', d: '访问令牌', e: '' }, { n: 'data.admin', t: 'object', r: '', d: '管理员信息', e: '' }], reqExample: 'POST /api/admin/auth/login\nContent-Type: application/json\n\n{"username":"admin","password":"Abc123456"}', resExample: '{"code":0,"data":{"access_token":"eyJ...","admin":{"id":1,"username":"admin"}},"message":"success"}' },
    ]},
    { name: '仪表盘', items: [
      { id: 'get-admin-dashboard', title: '仪表盘概览', method: 'GET', url: '/api/admin/dashboard', desc: '获取仪表盘统计数据', reqParams: [], resParams: [{ n: 'data.user_count', t: 'int', r: '', d: '用户总数', e: '' }, { n: 'data.host_count', t: 'int', r: '', d: '主机总数', e: '' }, { n: 'data.income', t: 'number', r: '', d: '本月收入', e: '' }], reqExample: 'GET /api/admin/dashboard\nAuthorization: JWT eyJ...', resExample: '{"code":0,"data":{"user_count":100,"host_count":50,"income":50000},"message":"success"}' },
    ]},
    { name: '用户管理', items: [
      { id: 'get-admin-users', title: '用户列表', method: 'GET', url: '/api/admin/users', desc: '获取用户列表', reqParams: [{ n: 'page', t: 'int', r: '选填', d: '页码', e: '1' }, { n: 'keyword', t: 'string', r: '选填', d: '搜索', e: 'test' }], resParams: [{ n: 'data.list', t: 'array', r: '', d: '用户列表', e: '' }], reqExample: 'GET /api/admin/users?page=1&keyword=test\nAuthorization: JWT eyJ...', resExample: '{"code":0,"data":{"list":[{"id":1,"username":"test","email":"test@example.com"}],"total":1},"message":"success"}' },
      { id: 'get-admin-user-detail', title: '用户详情', method: 'GET', url: '/api/admin/users/:id', desc: '获取用户详情(不脱敏)', reqParams: [{ n: 'id', t: 'int', r: '必填', d: '用户ID', e: '1' }], resParams: [{ n: 'data', t: 'object', r: '', d: '用户完整信息', e: '' }], reqExample: 'GET /api/admin/users/1\nAuthorization: JWT eyJ...', resExample: '{"code":0,"data":{"id":1,"username":"test","phone":"13800138000","balance":500},"message":"success"}' },
    ]},
    { name: '产品管理', items: [
      { id: 'get-admin-products', title: '产品列表', method: 'GET', url: '/api/admin/products', desc: '获取产品列表', reqParams: [{ n: 'page', t: 'int', r: '选填', d: '页码', e: '1' }, { n: 'group_id', t: 'int', r: '选填', d: '分组ID', e: '1' }], resParams: [{ n: 'data.list', t: 'array', r: '', d: '产品列表', e: '' }], reqExample: 'GET /api/admin/products?page=1\nAuthorization: JWT eyJ...', resExample: '{"code":0,"data":{"list":[{"id":1,"name":"基础型云服务器","price":"99.00"}],"total":1},"message":"success"}' },
    ]},
    { name: '工单管理', items: [
      { id: 'get-admin-tickets', title: '工单列表', method: 'GET', url: '/api/admin/tickets', desc: '获取工单列表', reqParams: [{ n: 'page', t: 'int', r: '选填', d: '页码', e: '1' }, { n: 'status', t: 'string', r: '选填', d: '状态', e: 'open' }], resParams: [{ n: 'data.list', t: 'array', r: '', d: '工单列表', e: '' }], reqExample: 'GET /api/admin/tickets?page=1&status=open\nAuthorization: JWT eyJ...', resExample: '{"code":0,"data":{"list":[{"id":1,"subject":"服务器无法访问","status":"open"}],"total":1},"message":"success"}' },
    ]},
    { name: '系统设置', items: [
      { id: 'get-admin-settings-site', title: '站点设置', method: 'GET', url: '/api/admin/settings/site', desc: '获取站点设置', reqParams: [], resParams: [{ n: 'data', t: 'object', r: '', d: '站点设置', e: '' }], reqExample: 'GET /api/admin/settings/site\nAuthorization: JWT eyJ...', resExample: '{"code":0,"data":{"site_name":"锚点财务","contact_phone":"400-xxx-xxxx"},"message":"success"}' },
      { id: 'put-admin-settings-site', title: '保存站点设置', method: 'PUT', url: '/api/admin/settings/site', desc: '保存站点设置', reqParams: [{ n: 'site_name', t: 'string', r: '选填', d: '站点名称', e: '锚点财务' }, { n: 'contact_phone', t: 'string', r: '选填', d: '联系电话', e: '400-xxx-xxxx' }], resParams: [{ n: 'message', t: 'string', r: '', d: '提示消息', e: '' }], reqExample: 'PUT /api/admin/settings/site\nAuthorization: JWT eyJ...\n\n{"site_name":"锚点财务"}', resExample: '{"code":0,"data":null,"message":"保存成功"}' },
    ]},
  ]},
}
