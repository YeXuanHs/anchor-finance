#!/usr/bin/env python3
"""Build complete API docs with query params, URL params, and response fields."""
import re, os, glob, json

HANDLER_DIR = os.path.join(os.path.dirname(__file__), '..', 'backend', 'internal', 'handler')
API_DIR = os.path.join(os.path.dirname(__file__), '..', 'backend', 'internal', 'api')
DOCS_DIR = os.path.dirname(__file__)

TYPE_MAP = {'string':'string','int':'int','int8':'int','int64':'int','uint':'int','float32':'number','float64':'number','bool':'bool','datatypes.JSON':'object','[]string':'array','time.Time':'string','*time.Time':'string'}
PD = {'account':'用户名/邮箱/手机号','password':'密码','old_password':'旧密码','new_password':'新密码','confirm_password':'确认密码','email':'邮箱','phone':'手机号','code':'验证码','username':'用户名','nickname':'昵称','access_token':'访问令牌','refresh_token':'刷新令牌','api_key':'API密钥','captcha_id':'验证码ID','answer':'验证码答案','product_id':'产品ID','billing_cycle':'计费周期','quantity':'数量','config':'配置选项(JSON)','domain':'域名','coupon_code':'优惠码','payment_method':'支付方式','payment_method_id':'支付方式ID','amount':'金额','type':'类型','target':'目标','open':'1=开启 0=关闭','status':'状态','keyword':'关键词','page':'页码','page_size':'每页数量','subject':'主题','content':'内容','department_id':'部门ID','priority':'优先级','remark':'备注','name':'名称','description':'描述','price':'价格','host_id':'主机ID','group_id':'分组ID','real_name':'真实姓名','id_card':'身份证号','title':'标题','sort_order':'排序','position':'位置','days':'天数','user_id':'用户ID','action':'操作','provider':'提供商','ip':'IP','os':'操作系统','cpu':'CPU','cpu_cores':'CPU核数','memory_mb':'内存(MB)','disk_size_gb':'磁盘(GB)','bandwidth_mbps':'带宽(Mbps)','datacenter_id':'数据中心ID','months':'月数','gateway':'网关','invoice_id':'账单ID','protocol':'协议','ext_port':'外部端口','int_port':'内部端口','int_ip':'内部IP','direction':'方向','port_range':'端口范围','source':'来源IP','default_action':'默认动作','cron_expr':'Cron表达式','command':'命令','timeout':'超时(秒)','comment':'备注','admin_note':'管理员备注','approve':'是否同意','method':'方式','reason':'原因','limit':'额度','start_date':'开始日期','end_date':'结束日期','contract_no':'合同号','bill_generation_day':'账单生成日','repayment_period':'还款期限(天)','is_active':'是否启用','is_published':'是否发布','language':'语言','parent_id':'上级ID','commission_rate':'佣金比例','template_id':'模板ID','channel':'渠道','days_before':'提前天数','owner_id':'所有者ID','expired_at':'过期时间','redirect_url':'回调URL','mark_read':'标记已读','transfer_code':'转移码','new_owner_id':'新所有者ID','client_id':'客户ID','server_id':'服务器ID','admin_id':'管理员ID','role_id':'角色ID','gateway_id':'网关ID','dept_id':'部门ID','ticket_id':'工单ID','sale_id':'销售ID','agent_id':'代理ID','service_id':'服务ID','module_id':'模块ID','plugin_id':'插件ID','file_id':'文件ID','msg_id':'消息ID','article_id':'文章ID','category_id':'分类ID','nav_id':'导航ID','banner_id':'轮播图ID','announcement_id':'公告ID','news_id':'新闻ID','link_id':'链接ID','suffix_id':'后缀ID','reply_id':'回复ID','ids':'ID数组','provider_id':'提供商ID','iso_id':'ISO镜像ID','key':'验证码key','digits':'验证码答案','host_id':'主机ID','images':'图片数组','user_ids':'用户ID数组','slug':'URL别名','image_url':'图片URL','link_url':'链接URL','button_text':'按钮文字','media_url':'媒体URL','cover_image':'封面图','keywords':'关键词','summary':'摘要','body':'HTML内容','tags':'标签','force':'是否强制','virtual_type':'虚拟化类型','price_monthly':'月价格','rack':'机架','rack_position':'机架位置','hostname':'主机名','mac':'MAC地址','ipv6':'IPv6地址','disk_type':'磁盘类型','traffic_gb':'流量(GB)','parent_server_id':'宿主机ID','oauth_code':'OAuth授权码','config_group_id':'配置组ID','option_id':'选项ID','tracker_id':'跟踪器ID','cancel_reason_id':'注销原因ID','target_email':'目标邮箱','target_phone':'目标手机号','is_sticky':'是否置顶','condition':'条件','data':'数据','email_suffix':'邮箱后缀','suffix':'后缀','content_type':'内容类型','work_start':'工作开始时间','work_end':'工作结束时间','holiday':'节假日','is_workday':'是否工作日','year':'年份','month':'月份','quarter':'季度','start_time':'开始时间','end_time':'结束时间','refund_amount':'退款金额','pay_amount':'支付金额','discount_amount':'折扣金额','tax_amount':'税额','subtotal':'小计','total':'总计','currency':'货币','exchange_rate':'汇率','balance_before':'变动前余额','balance_after':'变动后余额','before':'变动前','after':'变动后','old_status':'原状态','new_status':'新状态','operator':'操作人','operator_id':'操作人ID','target_type':'目标类型','target_id':'目标ID','ip_address':'IP地址','user_agent':'浏览器','request_url':'请求URL','request_method':'请求方法','response_code':'响应码','duration':'耗时(ms)','error_msg':'错误信息','stack_trace':'堆栈','file_name':'文件名','file_size':'文件大小','mime_type':'MIME类型','original_name':'原始文件名','storage_path':'存储路径','download_url':'下载URL','url':'URL','icon':'图标','is_default':'是否默认','is_system':'是否系统','max_uses':'最大使用次数','used_count':'已使用次数','min_amount':'最小金额','max_amount':'最大金额','discount_type':'折扣类型','discount_value':'折扣值','start_at':'开始时间','expire_at':'过期时间','scope':'范围','scope_id':'范围ID','ref_type':'关联类型','ref_id':'关联ID','ref_no':'关联编号','payment_time':'支付时间','payment_trade_no':'支付交易号','payment_amount':'支付金额','refund_time':'退款时间','refund_trade_no':'退款交易号','gateway_name':'网关名称','gateway_code':'网关代码','notify_url':'通知URL','return_url':'返回URL','sign_type':'签名类型','extra':'扩展参数','config_data':'配置数据','last_sync_at':'最后同步时间','sync_status':'同步状态','sync_message':'同步信息','upstream_id':'上游ID','upstream_product_id':'上游产品ID','upstream_host_id':'上游主机ID','client_area':'客户区域','admin_area':'管理区域','module_name':'模块名称','module_code':'模块代码','module_type':'模块类型','setting_key':'设置键','setting_value':'设置值','setting_group':'设置分组','setting_desc':'设置描述','setting_type':'设置类型','sort':'排序字段','order':'排序方式','start':'开始','end':'结束','format':'格式','export_type':'导出类型','batch_id':'批次ID','task_id':'任务ID','task_name':'任务名称','task_type':'任务类型','task_status':'任务状态','next_run_at':'下次执行时间','last_run_at':'上次执行时间','run_count':'执行次数','fail_count':'失败次数','success_count':'成功次数','total_count':'总数','message':'消息','result':'结果','error':'错误','warning':'警告','info':'信息','debug':'调试'}

EX = {'account':'admin','password':'Abc123456','old_password':'Abc123456','new_password':'NewPass789','confirm_password':'NewPass789','email':'user@example.com','phone':'13800138000','code':'123456','username':'testuser','nickname':'测试用户','access_token':'eyJ...','refresh_token':'eyJ...','api_key':'aB3$xY7!kL9m','captcha_id':'abc123','answer':'x5km','product_id':'1','billing_cycle':'monthly','quantity':'1','config':'{}','domain':'example.com','coupon_code':'SAVE20','payment_method':'alipay','payment_method_id':'1','amount':'99.00','type':'totp','target':'13800138000','open':'1','status':'active','keyword':'test','page':'1','page_size':'20','subject':'服务器无法访问','content':'服务器无法正常访问','department_id':'1','priority':'high','remark':'生产环境','name':'基础型云服务器','description':'产品描述','price':'99.00','host_id':'1','group_id':'1','real_name':'张三','id_card':'110101199001011234','title':'系统维护通知','sort_order':'0','position':'home','days':'30','user_id':'1','action':'start','provider':'github','ip':'192.168.1.1','os':'CentOS 7','cpu_cores':'4','memory_mb':'8192','disk_size_gb':'100','bandwidth_mbps':'10','datacenter_id':'1','virtual_type':'KVM','price_monthly':'99.00','force':'false','months':'1','gateway':'alipay','invoice_id':'1','protocol':'tcp','ext_port':'8080','int_port':'80','int_ip':'10.0.0.1','direction':'inbound','port_range':'80-443','source':'0.0.0.0/0','default_action':'allow','cron_expr':'0 */15 * * *','timeout':'300','approve':'true','method':'alipay','limit':'5000.00','start_date':'2026-01-01','end_date':'2026-12-31','contract_no':'HT2026001','bill_generation_day':'1','repayment_period':'30','is_active':'true','is_published':'true','language':'zh-CN','parent_id':'0','commission_rate':'10','template_id':'1','channel':'email','days_before':'7','owner_id':'1','expired_at':'2027-01-01','oauth_code':'auth_code_xxx','mark_read':'true','transfer_code':'TC123456','client_id':'1','server_id':'1','admin_id':'1','role_id':'1','gateway_id':'1','dept_id':'1','ticket_id':'1','sale_id':'1','agent_id':'1','service_id':'1','module_id':'1','plugin_id':'1','file_id':'1','msg_id':'1','article_id':'1','category_id':'1','nav_id':'1','banner_id':'1','announcement_id':'1','news_id':'1','link_id':'1','suffix_id':'1','reply_id':'1','config_group_id':'1','option_id':'1','digits':'1234','key':'captcha_key','iso_id':'1','hostname':'web01','mac':'00:11:22:33:44:55','ipv6':'2001:db8::1','tags':'生产','rack':'A01','rack_position':'U1','cpu':'Intel Xeon','reason':'额度不足','admin_note':'已处理','comment':'手动执行','command':'/usr/bin/backup.sh','condition':'{}','cover_image':'/uploads/cover.jpg','keywords':'云服务器','summary':'摘要','body':'<p>内容</p>','images':'[]','user_ids':'[1,2,3]','target_email':'new@example.com','target_phone':'13900139000','new_owner_id':'2','ids':'1,2,3','provider_id':'1','code':'alipay','data':'{}','slug':'test','image_url':'/img.jpg','link_url':'/products','button_text':'立即选购','media_url':'/video.mp4','is_sticky':'false','redirect_url':'/callback','page_size':'20','sort':'id','order':'desc'}

def go_type_to_js(t):
    t = t.strip('*')
    if t.startswith('[]'): return 'array'
    if t.startswith('map['): return 'object'
    return TYPE_MAP.get(t, 'string')

def extract_structs(content):
    structs = {}
    struct_map = {}
    for m in re.finditer(r'type\s+(\w+)\s+struct\s*\{([^}]*(?:\{[^}]*\}[^}]*)*)\}', content, re.DOTALL):
        name, body = m.group(1), m.group(2)
        fields = []
        for fm in re.finditer(r'(\w+)\s+(\S+)\s+`json:"([^"]+)"[^`]*`', body):
            fname = fm.group(3).split(',')[0]
            if fname in ('-', ''): continue
            ftype = go_type_to_js(fm.group(2))
            is_req = 'required' in fm.group(0)
            fields.append({'n':fname,'t':ftype,'r':'必填' if is_req else '选填','d':PD.get(fname,fname),'e':EX.get(fname,'')})
        if fields:
            structs[name] = fields
            struct_map[name] = fields
    func_positions = []
    for m in re.finditer(r'func\s+\([^)]+\)\s+(\w+)\([^)]*\*gin\.Context[^)]*\)\s*\{', content):
        func_positions.append((m.end(), m.group(1)))
    for i, (pos, fn) in enumerate(func_positions):
        if fn in structs: continue
        end_pos = func_positions[i+1][0] if i+1 < len(func_positions) else min(pos+5000, len(content))
        chunk = content[pos:end_pos]
        inline = re.search(r'var\s+req\s+struct\s*\{([^}]+)\}', chunk)
        if inline:
            fields = []
            for fm in re.finditer(r'(\w+)\s+(\S+)\s+`json:"([^"]+)"[^`]*`', inline.group(1)):
                fname = fm.group(3).split(',')[0]
                if fname in ('-', ''): continue
                ftype = go_type_to_js(fm.group(2))
                is_req = 'required' in fm.group(0)
                fields.append({'n':fname,'t':ftype,'r':'必填' if is_req else '选填','d':PD.get(fname,fname),'e':EX.get(fname,'')})
            if fields: structs[fn] = fields
            continue
        named = re.search(r'var\s+req\s+(\w+)', chunk)
        if named and named.group(1) in struct_map:
            structs[fn] = struct_map[named.group(1)]
    return structs

def extract_query_params(content):
    func_params = {}
    func_positions = []
    for m in re.finditer(r'func\s+\([^)]+\)\s+(\w+)\([^)]*\*gin\.Context[^)]*\)\s*\{', content):
        func_positions.append((m.end(), m.group(1)))
    for i, (pos, fn) in enumerate(func_positions):
        end_pos = func_positions[i+1][0] if i+1 < len(func_positions) else min(pos+10000, len(content))
        chunk = content[pos:end_pos]
        params = []
        for qm in re.finditer(r'c\.(?:Query|DefaultQuery)\(\s*"(\w+)"\s*(?:,\s*"([^"]*)")?\s*\)', chunk):
            qname = qm.group(1)
            default = qm.group(2) or ''
            ptype = 'int' if any(x in qname for x in ['page','limit','size','count','num','id','port','days','month','year']) else 'string'
            params.append({'n':qname,'t':ptype,'r':'选填','d':PD.get(qname,qname),'e':EX.get(qname,default or qname),'src':'q'})
        for pm in re.finditer(r'c\.Param\(\s*"(\w+)"\s*\)', chunk):
            pname = pm.group(1)
            ptype = 'int' if 'id' in pname else 'string'
            params.append({'n':pname,'t':ptype,'r':'必填','d':PD.get(pname,pname),'e':EX.get(pname,'1'),'src':'p'})
        if params:
            func_params[fn] = params
    return func_params

def extract_responses(content):
    func_responses = {}
    func_positions = []
    for m in re.finditer(r'func\s+\([^)]+\)\s+(\w+)\([^)]*\*gin\.Context[^)]*\)\s*\{', content):
        func_positions.append((m.end(), m.group(1)))
    for i, (pos, fn) in enumerate(func_positions):
        end_pos = func_positions[i+1][0] if i+1 < len(func_positions) else min(pos+10000, len(content))
        chunk = content[pos:end_pos]
        resp_fields = []
        for rm in re.finditer(r'response\.Success\(\s*c\s*,\s*gin\.H\{([^}]+)\}', chunk):
            for fm in re.finditer(r'"(\w+)":', rm.group(1)):
                resp_fields.append({'n':f'data.{fm.group(1)}','t':'object','d':fm.group(1)})
        for rm in re.finditer(r'response\.Success\(\s*c\s*,\s*(\w+)\s*\)', chunk):
            vn = rm.group(1)
            if vn not in ('c', 'gin', 'nil'):
                resp_fields.append({'n':'data','t':'object','d':'返回对象'})
        if re.search(r'response\.SuccessWithList\(', chunk):
            resp_fields.append({'n':'data.list','t':'array','d':'列表'})
            resp_fields.append({'n':'data.total','t':'int','d':'总数'})
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

def extract_routes(fpath):
    with open(fpath, 'r', encoding='utf-8') as f:
        content = f.read()
    routes = []
    base = '/api/admin' if os.sep + 'admin' + os.sep in fpath else '/api/v1'
    prefixes = {'r': base, 'public': base, 'user': base, 'admin': base}
    for m in re.finditer(r'(\w+)\s*:?=\s*\w+\.Group\("([^"]*)"', content):
        prefixes[m.group(1)] = base + m.group(2)
    for m in re.finditer(r'(\w+)\.(GET|POST|PUT|DELETE|PATCH)\(\s*"([^"]+)"\s*,\s*(\w+\.(\w+))', content):
        vn, method, path, handler, mn = m.group(1), m.group(2), m.group(3), m.group(4), m.group(5)
        if not path.startswith('/'): path = '/' + path
        full_path = prefixes.get(vn, base) + path
        routes.append({'group':vn,'method':method,'path':full_path,'method_name':mn})
    for m in re.finditer(r'(\w+)\.(GET|POST|PUT|DELETE|PATCH)\(\s*"([^"]+)"\s*,\s*func\s*\(', content):
        vn, method, path = m.group(1), m.group(2), m.group(3)
        if not path.startswith('/'): path = '/' + path
        full_path = prefixes.get(vn, base) + path
        func_start = m.end()
        func_body = content[func_start:func_start+3000]
        fields = []
        inline = re.search(r'var\s+req\s+struct\s*\{([^}]+)\}', func_body)
        if inline:
            for fm in re.finditer(r'(\w+)\s+(\S+)\s+`json:"([^"]+)"[^`]*`', inline.group(1)):
                fname = fm.group(3).split(',')[0]
                if fname in ('-', ''): continue
                ftype = go_type_to_js(fm.group(2))
                is_req = 'required' in fm.group(0)
                fields.append({'n':fname,'t':ftype,'r':'必填' if is_req else '选填','d':PD.get(fname,fname),'e':EX.get(fname,'')})
        qp = []
        for qm in re.finditer(r'c\.(?:Query|DefaultQuery)\(\s*"(\w+)"\s*(?:,\s*"([^"]*)")?\s*\)', func_body):
            qname = qm.group(1)
            default = qm.group(2) or ''
            qp.append({'n':qname,'t':'string','r':'选填','d':PD.get(qname,qname),'e':EX.get(qname,default),'src':'q'})
        for pm in re.finditer(r'c\.Param\(\s*"(\w+)"\s*\)', func_body):
            pname = pm.group(1)
            qp.append({'n':pname,'t':'string','r':'必填','d':PD.get(pname,pname),'e':EX.get(pname,'1'),'src':'p'})
        resp_fields = []
        for rm in re.finditer(r'response\.Success\(\s*c\s*,\s*gin\.H\{([^}]+)\}', func_body):
            for fm in re.finditer(r'"(\w+)":', rm.group(1)):
                resp_fields.append({'n':f'data.{fm.group(1)}','t':'object','d':fm.group(1)})
        if re.search(r'response\.SuccessMsg\(', func_body) and not resp_fields:
            resp_fields.append({'n':'message','t':'string','d':'提示消息'})
        routes.append({'group':vn,'method':method,'path':full_path,'method_name':'inline','req_fields':fields,'query_params':qp,'resp_fields':resp_fields})
    return routes

NAMES = {'auth':'认证','login':'登录','register':'注册','refresh':'刷新','logout':'退出','access-token':'API登录','sms-login':'短信登录','password':'密码','verify-code':'验证重置码','reset':'重置','captcha':'验证码','config':'配置','generate':'生成','verify':'验证','check':'检查','geetest':'极验','products':'产品','product':'产品','hot':'热门','groups':'分组','pricing':'价格','users':'用户','user':'用户','payment-methods':'支付方式','system':'系统','lang':'语言','settings':'设置','public':'公开','menus':'菜单','nav':'导航','top':'顶部','bottom':'底部','homepage':'首页','base-info':'基础信息','downloads':'下载中心','announcements':'公告','news':'新闻','categories':'分类','help':'帮助中心','articles':'文章','search':'搜索','related':'相关','feedback':'反馈','languages':'语言列表','translations':'翻译','promo-codes':'优惠码','validate':'验证','oauth':'OAuth','provider':'提供商','callback':'回调','banners':'轮播图','contact':'联系','profile':'资料','change-password':'改密码','bind-phone':'绑手机','bind-email':'绑邮箱','2fa':'二步验证','enable':'开启','disable':'关闭','api':'API','summary':'摘要','open':'开关','providers':'提供商列表','bind':'绑定','unbind':'解绑','accounts':'已绑定账号','tastes':'偏好','login-logs':'登录日志','logs':'日志','hosts':'主机列表','host':'主机','operations':'可用操作','billing':'账单','download':'下载','remark':'备注','cart':'购物车','add':'添加','checkout':'结算','batch-delete':'批量删除','orders':'订单列表','order':'订单','pay':'支付','cancel':'取消','preview':'预览','invoices':'账单列表','invoice':'账单','combine':'合并','tickets':'工单列表','ticket':'工单','reply':'回复','close':'关闭','attachments':'附件','upload':'上传','balance':'余额','recharge':'充值','withdraw':'提现','credit':'信用额度','apply':'申请','repay':'还款','bills':'账单','contacts':'联系人','default':'默认','messages':'消息','unread-count':'未读数','read':'已读','read-all':'全部已读','certification':'实名认证','status':'状态','submit':'提交','enterprise':'企业认证','upgrades':'升级','available':'可用升级','contracts':'合同','sign':'签署','affiliate':'代理推荐','info':'信息','records':'记录','vouchers':'发票','voucher-types':'发票抬头','voucher-posts':'收件地址','multi-renew':'批量续费','execute':'执行','api-logs':'API日志','product-groups':'产品分组','product-diverts':'产品转移','received':'收到的转移','accept':'接受','reject':'拒绝','transfer-code':'转移码','regenerate-code':'重新生成转移码','system-logs':'系统日志','export':'导出','sms':'短信','send':'发送','ssl':'SSL证书','certificates':'证书','install':'安装','renew':'续费','ai-shopping':'AI购物助手','session':'会话','config-options':'配置选项','v10cloud':'V10云','transactions':'交易记录','earnings':'收益','mine':'我的挂售','buyer':'买家订单','seller':'卖家订单','complete':'完成','chat-sessions':'聊天会话','dashboard':'仪表盘','stats':'统计','income-trend':'收入趋势','global-search':'全局搜索','duplicate':'复制','stock':'库存','batch-update':'批量更新','sort':'排序','discounts':'折扣','res':'资源','select-type':'选择类型','first-groups':'一级分组','create':'创建','delete':'删除','update':'更新','manage':'管理','sale':'销售','activate':'激活','change-status':'变更状态','clients':'客户列表','enhanced':'增强版','refund':'退款','items':'条目','assign':'指派','merge':'合并','transfer':'转交','statistics':'统计','friendly-links':'友情链接','tree':'树形','lang-keys':'语言键','nav-groups':'导航分组','log-records':'操作日志','rules':'规则','sale-records':'销售记录','sale-users':'销售用户','admin-list':'管理员列表','timetypes':'时间类型','batch':'批量','retry-failed':'重试失败','templates':'模板','send-batch':'批量发送','batches':'批次','ticket-depts':'工单部门','members':'成员','managers':'管理者','servers':'服务器','options':'选项','server-list':'服务器列表','test-link':'测试连接','certifi':'实名认证配置','advanced-options':'高级选项','links':'链接','batch-sync':'批量同步','community':'社区','posts':'帖子','comments':'评论','cron-url':'URL定时任务','person':'个人','company':'企业','black-list':'黑名单','cancel-requests':'注销申请','authorize':'授权','group':'用户组','review':'审核','product-accounts':'产品账号','login-as':'代登录','record-log':'记录日志','operation-logs':'操作日志','provisions':'自动开通','test':'测试','suspend':'暂停','terminate':'终止','unsuspend':'解除暂停','usage':'用量','traffic':'流量','refresh':'刷新','detail':'详情','rescue':'救援','snapshots':'快照','backups':'备份','restore':'恢复','power':'电源','os-list':'系统列表','plugins':'插件','server-modules':'服务器模块','server-groups':'服务器分组','client-groups':'客户分组','client-services':'客户服务','delete-host':'删除主机','batch-renew':'批量续费','process':'处理','adjust-balance':'调整余额','reset-password':'重置密码','add-note':'添加备注','get-notes':'获取备注','client-logs':'客户日志','host-renew':'主机续费','upgrade-config':'升级配置','user-manage':'用户管理','filter':'筛选','ban':'封禁','unban':'解封','boot':'开机','shutdown':'关机','reboot':'重启','config-general':'通用配置','config-email':'邮件配置','config-sms':'短信配置','config-payment':'支付配置','config-security':'安全配置','config-affiliate':'代理配置','config-captcha':'验证码配置','config-register':'注册配置','config-login':'登录配置','config-api':'API配置','config-2fa':'二步验证配置','config-invoice':'发票配置','config-recharge':'充值配置','config-safe':'安全配置','config-local':'本地配置','config-buy-product':'购买产品配置','config-language':'语言配置','config-header-footer':'页头页脚配置','config-login-page':'登录页配置','config-nav-groups':'导航分组配置','config-certifi':'实名认证配置','config-smtp-test':'SMTP测试','config-sms-test':'短信测试','config-debug-mode':'调试模式','config-second-verify':'二次验证配置','email':'邮件','email-templates':'邮件模板','notification':'通知','kb':'知识库','oauth-provider':'OAuth提供商','aggregate-login':'聚合登录','account':'账号','client-care':'客户关怀','client-group':'客户分组','client-track':'客户跟踪','cancel-reason':'注销原因','email-suffix':'邮箱后缀','sms-template':'短信模板','download-cate':'下载分类','marketplace':'交易市场','marketplace-config':'交易市场配置','setting':'设置','system-log':'系统日志','client-service':'客户服务','notices':'公告','maintenance':'维护','solutions':'解决方案','datacenters':'数据中心','colocation':'托管','advantages':'优势','ai':'AI','payment':'支付','sub':'子分类','data':'数据','id':'ID','page':'页码','page_size':'每页数量','keyword':'关键词','code':'代码','ids':'ID数组','host_id':'主机ID','days':'天数','enhanced':'增强版'}

def infer_title(path, method):
    parts = path.strip('/').split('/')
    if len(parts) >= 3 and parts[0] == 'api' and parts[1] in ('v1', 'admin', 'public'): parts = parts[2:]
    if not parts: return '根路径'
    readable = [NAMES.get(re.sub(r'^:', '', p), re.sub(r'^:', '', p)) for p in parts]
    last = readable[-1]
    parent = '/'.join(readable[:-1]) if len(readable) > 1 else ''
    if method == 'GET':
        if parts[-1] in ('list','index') or 'list' in parts[-1].lower(): return f"获取{parent}列表" if parent else "列表"
        elif re.match(r'^:\w+$|^id$', parts[-1].lower()): return f"获取{parent}详情" if parent else "详情"
        elif parts[-1] == 'search': return f"搜索{parent}" if parent else "搜索"
        elif parts[-1] == 'hot': return f"获取{parent}热门" if parent else "热门"
        elif parts[-1] == 'active': return f"获取{parent}有效列表" if parent else "有效列表"
        else: return f"获取{last}"
    elif method == 'POST': return last
    elif method == 'PUT': return f"更新{last}"
    elif method == 'DELETE': return f"删除{last}"
    return f"{method} {'/'.join(readable)}"

def path_to_id(path, method):
    p = re.sub(r'^/api/(v1|admin|public)/?', '', path).strip('/')
    p = re.sub(r':(\w+)', r'\1', p).replace('/', '-')
    return f"{method.lower()}-{p}" if p else f"{method.lower()}-root"

def main():
    all_structs = {}
    all_query_params = {}
    all_responses = {}
    for fpath in sorted(glob.glob(os.path.join(HANDLER_DIR, '*.go'))):
        with open(fpath, 'r', encoding='utf-8') as f:
            content = f.read()
        all_structs.update(extract_structs(content))
        all_query_params.update(extract_query_params(content))
        all_responses.update(extract_responses(content))
    print(f"Structs: {len(all_structs)}, Query: {len(all_query_params)}, Resp: {len(all_responses)}")

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
        rid = path_to_id(r['path'], r['method'])
        title = infer_title(r['path'], r['method'])
        url, method = r['path'], r['method']
        req_p = r.get('req_fields', [])
        req_s = ','.join([f"{{n:'{p['n']}',t:'{p['t']}',r:'{p['r']}',d:'{p['d']}',e:'{p['e']}'}}" for p in req_p])
        resp_p = r.get('resp_fields', [{'n':'data','t':'object','d':'返回数据'}])
        res_parts = []
        for p in resp_p:
            pt = p.get('t', 'object')
            res_parts.append("{n:'" + p['n'] + "',t:'" + pt + "',r:'',d:'" + p['d'] + "',e:''}")
        res_s = ','.join(res_parts)

        body_p = {p['n']:p['e'] for p in req_p if p.get('src') not in ('q','p') and p['e']}
        query_p = {p['n']:p['e'] for p in req_p if p.get('src')=='q' and p['e']}
        path_p = {p['n']:p['e'] for p in req_p if p.get('src')=='p' and p['e']}

        actual_url = url
        for pn, pv in path_p.items():
            actual_url = actual_url.replace(f':{pn}', pv)

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
                if '.' not in k: data[k] = '...' if p.get('t')=='array' else '{}'
            res_ex = json.dumps({'code':0,'data':data}, ensure_ascii=False)
        else:
            res_ex = '{"code":0,"data":{}}'

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
    print(f"Total: {total}, Req: {wr}, Resp: {wrs}")
    print(f"Pub: {len(pub)}, Usr: {len(usr)}, Adm: {len(adm)}")

if __name__ == '__main__':
    main()
