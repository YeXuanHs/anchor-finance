#!/usr/bin/env python3
"""Enhance API docs: add inferred params for routes missing them."""
import re, os, glob, json, copy

HANDLER_DIR = os.path.join(os.path.dirname(__file__), '..', 'backend', 'internal', 'handler')
API_DIR = os.path.join(os.path.dirname(__file__), '..', 'backend', 'internal', 'api')
DOCS_DIR = os.path.dirname(__file__)

PD = {'id':'ID','page':'页码','page_size':'每页数量','keyword':'关键词','status':'状态','type':'类型','group_id':'分组ID','user_id':'用户ID','department_id':'部门ID','sort':'排序字段','order':'排序方式(asc/desc)','is_active':'是否启用','provider':'提供商','name':'名称','code':'代码','ids':'ID数组(逗号分隔)','host_id':'主机ID','days':'天数','start_date':'开始日期(YYYY-MM-DD)','end_date':'结束日期(YYYY-MM-DD)','limit':'每页数量','offset':'偏移量','ticket_id':'工单ID','invoice_id':'账单ID','order_id':'订单ID','category_id':'分类ID','admin_id':'管理员ID','role_id':'角色ID','server_id':'服务器ID','product_id':'产品ID','gateway_id':'网关ID','dept_id':'部门ID','sale_id':'销售ID','agent_id':'代理ID','service_id':'服务ID','module_id':'模块ID','plugin_id':'插件ID','file_id':'文件ID','msg_id':'消息ID','article_id':'文章ID','nav_id':'导航ID','banner_id':'轮播图ID','announcement_id':'公告ID','news_id':'新闻ID','link_id':'链接ID','suffix_id':'后缀ID','reply_id':'回复ID','config_group_id':'配置组ID','option_id':'选项ID','tracker_id':'跟踪器ID','cancel_reason_id':'注销原因ID','parent_id':'上级ID','template_id':'模板ID','owner_id':'所有者ID','new_owner_id':'新所有者ID','client_id':'客户ID','message':'消息','remark':'备注','admin_note':'管理员备注','description':'描述','content':'内容','title':'标题','subject':'主题','priority':'优先级(high/medium/low)','target_email':'目标邮箱','target_phone':'目标手机号','month':'月份','year':'年份','quarter':'季度','format':'格式(json/csv)','export_type':'导出类型','batch_id':'批次ID','task_id':'任务ID','review_status':'审核状态','mark_read':'是否标记已读','position':'位置(home)','language':'语言代码(zh-CN/en-US)','datacenter_id':'数据中心ID','parent_server_id':'宿主机ID','upstream_id':'上游ID','scope':'范围','confirm':'确认','force':'是否强制','comment':'备注','reason':'原因','action':'操作类型','ip':'IP地址','provider_id':'提供商ID','bill_generation_day':'账单生成日','repayment_period':'还款期限(天)','work_start':'工作开始时间','work_end':'工作结束时间','config':'配置选项(JSON)','billing_cycle':'计费周期','quantity':'数量','payment_method':'支付方式代码','payment_method_id':'支付方式ID','gateway':'支付网关代码','amount':'金额','coupon_code':'优惠码','domain':'域名','password':'密码','new_password':'新密码','username':'用户名','email':'邮箱','phone':'手机号','nickname':'昵称','real_name':'真实姓名','id_card':'身份证号','qq':'QQ号','company':'公司名','address':'地址','avatar':'头像URL','signature':'签名','api_key':'API密钥','api_password':'API密钥','token':'令牌','key':'验证码key','digits':'验证码答案','method':'方式','port':'端口','protocol':'协议','direction':'方向','source':'来源','default_action':'默认动作','comment':'备注','ext_port':'外部端口','int_port':'内部端口','int_ip':'内部IP','timeout':'超时(秒)','hostname':'主机名','cpu_cores':'CPU核数','memory_mb':'内存(MB)','disk_size_gb':'磁盘(GB)','bandwidth_mbps':'带宽(Mbps)','traffic_gb':'流量(GB)','virtual_type':'虚拟化类型','price_monthly':'月价格','tags':'标签','cron_expr':'Cron表达式','command':'命令','os':'操作系统','iso_id':'ISO镜像ID','images':'图片数组','cover_image':'封面图','keywords':'关键词','summary':'摘要','body':'HTML内容','slug':'URL别名','image_url':'图片URL','link_url':'链接URL','button_text':'按钮文字','media_url':'媒体URL','sort_order':'排序值','is_published':'是否发布','is_sticky':'是否置顶','commission_rate':'佣金比例','channel':'通知渠道','days_before':'提前天数','condition':'条件','data':'数据(JSON)','amount':'金额','approve':'是否同意','transfer_code':'转移码','redirect_url':'回调URL','oauth_code':'OAuth授权码','verify_code':'验证码','target':'目标(手机/邮箱)','open':'1=开启 0=关闭','captcha_id':'验证码ID','answer':'验证码答案','code':'验证码','captcha_output':'极验输出','lot_number':'流水号','pass_token':'验证令牌','gen_time':'生成时间','product_id':'产品ID','account':'用户名/邮箱/手机号','old_password':'旧密码','confirm_password':'确认密码','access_token':'访问令牌','refresh_token':'刷新令牌','search':'搜索关键词','scope_id':'范围ID','ref_type':'关联类型','ref_id':'关联ID','ref_no':'关联编号','setting_key':'设置键','setting_value':'设置值','setting_group':'设置分组','upload_type':'上传类型','file_name':'文件名','start':'开始时间','end':'结束时间','ip_address':'IP地址','is_workday':'是否工作day','expired_at':'过期时间','before':'变动前','after':'变动后','refund_amount':'退款金额','pay_amount':'支付金额','old_status':'原状态','new_status':'新状态','target_type':'目标类型','target_id':'目标ID','operator_id':'操作人ID','error_msg':'错误信息','batch_ids':'批次ID数组','host_ids':'主机ID数组','user_ids':'用户ID数组','client_ids':'客户ID数组'}

EX = {'id':'1','page':'1','page_size':'20','keyword':'test','status':'active','type':'totp','group_id':'1','user_id':'1','department_id':'1','sort':'id','order':'desc','is_active':'true','name':'测试','code':'test','ids':'1,2,3','host_id':'1','days':'30','start_date':'2026-01-01','end_date':'2026-12-31','limit':'20','offset':'0','ticket_id':'1','invoice_id':'1','order_id':'1','category_id':'1','admin_id':'1','role_id':'1','server_id':'1','product_id':'1','gateway_id':'1','dept_id':'1','sale_id':'1','agent_id':'1','service_id':'1','module_id':'1','plugin_id':'1','file_id':'1','msg_id':'1','article_id':'1','nav_id':'1','banner_id':'1','announcement_id':'1','news_id':'1','link_id':'1','suffix_id':'1','reply_id':'1','config_group_id':'1','option_id':'1','tracker_id':'1','cancel_reason_id':'1','parent_id':'0','template_id':'1','owner_id':'1','new_owner_id':'2','client_id':'1','message':'操作成功','remark':'备注','admin_note':'管理员备注','description':'描述','content':'内容','title':'测试标题','subject':'工单主题','priority':'high','target_email':'new@example.com','target_phone':'13900139000','month':'1','year':'2026','quarter':'1','format':'json','export_type':'csv','batch_id':'1','task_id':'1','mark_read':'true','position':'home','language':'zh-CN','datacenter_id':'1','parent_server_id':'1','upstream_id':'1','scope':'all','confirm':'true','force':'false','comment':'手动执行','reason':'额度不足','action':'start','ip':'192.168.1.1','provider_id':'1','bill_generation_day':'1','repayment_period':'30','work_start':'09:00','work_end':'18:00','config':'{}','billing_cycle':'monthly','quantity':'1','payment_method':'alipay','payment_method_id':'1','gateway':'alipay','amount':'99.00','coupon_code':'SAVE20','domain':'example.com','password':'Abc123456','new_password':'NewPass789','username':'testuser','email':'user@example.com','phone':'13800138000','nickname':'测试用户','real_name':'张三','id_card':'110101199001011234','qq':'12345678','company':'测试公司','address':'北京市朝阳区','avatar':'/uploads/avatar.jpg','signature':'测试签名','api_key':'aB3$xY7!kL9m','api_password':'aB3$xY7!kL9m','token':'jwt_token','key':'captcha_key','digits':'1234','method':'alipay','port':'80','protocol':'tcp','direction':'inbound','source':'0.0.0.0/0','default_action':'allow','ext_port':'8080','int_port':'80','int_ip':'10.0.0.1','timeout':'300','hostname':'web01','cpu_cores':'4','memory_mb':'8192','disk_size_gb':'100','bandwidth_mbps':'10','traffic_gb':'1000','virtual_type':'KVM','price_monthly':'99.00','tags':'生产','cron_expr':'0 */15 * * *','command':'/usr/bin/backup.sh','os':'CentOS 7','iso_id':'1','images':'[]','cover_image':'/uploads/cover.jpg','keywords':'云服务器','summary':'摘要','body':'<p>内容</p>','slug':'test','image_url':'/img.jpg','link_url':'/products','button_text':'立即选购','media_url':'/video.mp4','sort_order':'0','is_published':'true','is_sticky':'false','commission_rate':'10','channel':'email','days_before':'7','condition':'{}','data':'{}','approve':'true','transfer_code':'TC123456','redirect_url':'/callback','oauth_code':'auth_code_xxx','verify_code':'123456','target':'13800138000','open':'1','captcha_id':'abc123','answer':'x5km','captcha_output':'xxx','lot_number':'lot_xxx','pass_token':'pass_xxx','gen_time':'1720000000','search':'keyword','scope_id':'1','ref_type':'order','ref_id':'1','ref_no':'ORD2026001','setting_key':'site_name','setting_value':'测试站','setting_group':'general','upload_type':'image','file_name':'test.jpg','start':'2026-01-01','end':'2026-12-31','expired_at':'2027-01-01','before':'100','after':'200','refund_amount':'50.00','pay_amount':'99.00','old_status':'pending','new_status':'active','target_type':'user','target_id':'1','operator_id':'1','error_msg':'错误','email':'user@example.com','provider':'github'}

def infer_params_from_path(path, method):
    """Infer required params from URL path like :id, :provider, etc."""
    params = []
    for m in re.finditer(r':(\w+)', path):
        pname = m.group(1)
        ptype = 'int' if 'id' in pname else 'string'
        desc = PD.get(pname, pname)
        example = EX.get(pname, '1')
        params.append({'n':pname,'t':ptype,'r':'必填','d':desc,'e':example,'src':'p'})
    return params

def infer_query_params(path, method):
    """Infer common query params based on route pattern."""
    params = []
    if method == 'GET':
        if 'list' in path or '/hosts' in path or '/orders' in path or '/invoices' in path or '/tickets' in path or '/users' in path or '/products' in path or '/logs' in path or '/records' in path or '/messages' in path or '/articles' in path or '/banners' in path or '/announcements' in path or '/news' in path or '/links' in path or '/certificates' in path or '/vouchers' in path:
            params.append({'n':'page','t':'int','r':'选填','d':'页码','e':'1','src':'q'})
            params.append({'n':'page_size','t':'int','r':'选填','d':'每页数量','e':'20','src':'q'})
        if 'search' in path or '/help/' in path:
            params.append({'n':'keyword','t':'string','r':'选填','d':'关键词','e':'test','src':'q'})
        if '/logs' in path or '/records' in path:
            params.append({'n':'start_date','t':'string','r':'选填','d':'开始日期','e':'2026-01-01','src':'q'})
            params.append({'n':'end_date','t':'string','r':'选填','d':'结束日期','e':'2026-12-31','src':'q'})
    return params

def build_enhanced_item(r):
    """Build JS item with enhanced params."""
    rid = r['id']
    title = r['title']
    url = r['url']
    method = r['method']
    
    req_p = r.get('req_params', [])
    resp_p = r.get('res_params', [])
    
    # Enhance req params: add path params if missing
    existing_names = {p['n'] for p in req_p}
    
    # Add path params from URL
    for m in re.finditer(r':(\w+)', url):
        pname = m.group(1)
        if pname not in existing_names:
            ptype = 'int' if 'id' in pname else 'string'
            req_p.append({'n':pname,'t':ptype,'r':'必填','d':PD.get(pname,pname),'e':EX.get(pname,'1')})
            existing_names.add(pname)
    
    # Add common query params for GET list endpoints
    if method == 'GET':
        path_lower = url.lower()
        is_list = any(x in path_lower for x in ['/list', '/logs', '/records', '/messages', '/articles', '/banners', '/announcements', '/news', '/links', '/certificates', '/vouchers', '/hosts', '/orders', '/invoices', '/tickets', '/users', '/products', '/clients', '/download', '/notes'])
        if is_list:
            if 'page' not in existing_names:
                req_p.append({'n':'page','t':'int','r':'选填','d':'页码','e':'1'})
                existing_names.add('page')
            if 'page_size' not in existing_names:
                req_p.append({'n':'page_size','t':'int','r':'选填','d':'每页数量','e':'20'})
                existing_names.add('page_size')
        if '/search' in path_lower and 'keyword' not in existing_names:
            req_p.append({'n':'keyword','t':'string','r':'选填','d':'关键词','e':'test'})
        if ('/log' in path_lower or '/record' in path_lower) and 'start_date' not in existing_names:
            req_p.append({'n':'start_date','t':'string','r':'选填','d':'开始日期','e':'2026-01-01'})
            req_p.append({'n':'end_date','t':'string','r':'选填','d':'结束日期','e':'2026-12-31'})
    
    # Build request example
    body_params = {p['n']:p['e'] for p in req_p if p.get('src') not in ('q','p') and p.get('e')}
    query_params = {p['n']:p['e'] for p in req_p if p.get('src') == 'q' and p.get('e')}
    
    actual_url = url
    for p in req_p:
        if p.get('src') == 'p' and p.get('e'):
            actual_url = actual_url.replace(':'+p['n'], p['e'])
    # Also replace any remaining :params with 1
    actual_url = re.sub(r':(\w+)', '1', actual_url)
    
    if method in ('POST','PUT','PATCH') and body_params:
        body = json.dumps(body_params, ensure_ascii=False)
        req_ex = f"{method} {actual_url}"
        if query_params: req_ex += '?' + '&'.join([f"{k}={v}" for k,v in query_params.items()])
        req_ex += f"\\nAuthorization: JWT eyJ...\\n\\n{body}"
    elif query_params:
        qs = '&'.join([f"{k}={v}" for k,v in query_params.items()])
        req_ex = f"{method} {actual_url}?{qs}\\nAuthorization: JWT eyJ..."
    else:
        need_auth = any(x in url for x in ['/user/','/hosts','/orders','/invoices','/tickets','/cart','/balance','/credit','/contacts','/messages','/certification','/marketplace','/product-diverts','/affiliate','/contracts','/vouchers','/upgrades','/oauth','/ssl','/ai-shopping','/v10cloud','/multi-renew','/api-logs','/login-logs','/admin'])
        req_ex = f"{method} {actual_url}\\nAuthorization: JWT eyJ..." if need_auth else f"{method} {actual_url}"
    
    # Build response example
    if resp_p:
        data = {}
        for p in resp_p:
            k = p['n'].replace('data.','')
            if '.' not in k and k not in ('code','message'):
                t = p.get('t','object')
                if t == 'array': data[k] = '...'
                elif t == 'int': data[k] = 0
                elif t == 'string': data[k] = ''
                elif t == 'number': data[k] = 0.0
                elif t == 'bool': data[k] = true
                else: data[k] = '{}'
        res_ex = json.dumps({'code':0,'data':data,'message':'success'}, ensure_ascii=False)
    else:
        res_ex = '{"code":0,"data":{},"message":"success"}'
    
    # Format req params
    req_s = ','.join([f"{{n:'{p['n']}',t:'{p['t']}',r:'{p['r']}',d:'{p['d']}',e:'{p.get('e','')}'}}" for p in req_p])
    
    # Format resp params
    res_parts = []
    for p in resp_p:
        pt = p.get('t', 'object')
        res_parts.append("{n:'" + p['n'] + "',t:'" + pt + "',r:'',d:'" + p['d'] + "',e:''}")
    res_s = ','.join(res_parts)
    
    return (f"        {{id:'{rid}',title:'{title}',method:'{method}',url:'{url}',"
            f"desc:'{title}',reqParams:[{req_s}],resParams:[{res_s}],"
            f"reqExample:'{req_ex}',resExample:'{res_ex}'}}")

def main():
    # Read existing api_data_complete.js
    js_path = os.path.join(DOCS_DIR, 'api_data_complete.js')
    with open(js_path, 'r', encoding='utf-8') as f:
        content = f.read()
    
    # Parse existing data
    # This is complex, so let's re-read the routes from router files and enhance
    
    all_structs = {}
    all_query_params = {}
    all_responses = {}
    
    from build_api_v3 import extract_structs, extract_query_params, extract_responses, extract_routes, infer_title, path_to_id, NAMES
    
    for fpath in sorted(glob.glob(os.path.join(HANDLER_DIR, '*.go'))):
        with open(fpath, 'r', encoding='utf-8') as f:
            fc = f.read()
        all_structs.update(extract_structs(fc))
        all_query_params.update(extract_query_params(fc))
        all_responses.update(extract_responses(fc))
    
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
                
                # Merge all params
                seen = set()
                combined = []
                for p in body_p + query_p:
                    if p['n'] not in seen:
                        seen.add(p['n'])
                        combined.append(p)
                
                # Add inferred path params
                for m in re.finditer(r':(\w+)', r['path']):
                    pname = m.group(1)
                    if pname not in seen:
                        ptype = 'int' if 'id' in pname else 'string'
                        combined.append({'n':pname,'t':ptype,'r':'必填','d':PD.get(pname,pname),'e':EX.get(pname,'1')})
                        seen.add(pname)
                
                # Add inferred query params for GET list endpoints
                if r['method'] == 'GET':
                    path_lower = r['path'].lower()
                    is_list = any(x in path_lower for x in ['/list', '/logs', '/records', '/messages', '/articles', '/banners', '/announcements', '/news', '/links', '/certificates', '/vouchers'])
                    if is_list:
                        if 'page' not in seen:
                            combined.append({'n':'page','t':'int','r':'选填','d':'页码','e':'1'})
                            seen.add('page')
                        if 'page_size' not in seen:
                            combined.append({'n':'page_size','t':'int','r':'选填','d':'每页数量','e':'20'})
                            seen.add('page_size')
                
                r['req_fields'] = combined
                r['resp_fields'] = resp_f
            all_routes.extend(routes)
    
    # Build items
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
    with open(js_path, 'w', encoding='utf-8') as f:
        f.write(js)
    
    total = len(all_routes)
    wr = sum(1 for r in all_routes if r.get('req_fields'))
    wrs = sum(1 for r in all_routes if r.get('resp_fields'))
    print(f"Total: {total}, Req: {wr} ({wr*100//total}%), Resp: {wrs} ({wrs*100//total}%)")
    print(f"Pub: {len(pub)}, Usr: {len(usr)}, Adm: {len(adm)}")

if __name__ == '__main__':
    main()
