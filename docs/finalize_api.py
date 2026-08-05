#!/usr/bin/env python3
"""Final pass: fill remaining empty params with inferred defaults."""
import re, json, os

DOCS_DIR = os.path.dirname(__file__)
js_path = os.path.join(DOCS_DIR, 'api_data_complete.js')

with open(js_path, 'r', encoding='utf-8') as f:
    content = f.read()

# Parse all items and fix them
def fix_item(match):
    """Fix a single API item."""
    item = match.group(0)
    
    # Extract fields
    id_m = re.search(r"id:'([^']*)'", item)
    title_m = re.search(r"title:'([^']*)'", item)
    method_m = re.search(r"method:'([^']*)'", item)
    url_m = re.search(r"url:'([^']*)'", item)
    
    if not all([id_m, title_m, method_m, url_m]):
        return item
    
    rid = id_m.group(1)
    title = title_m.group(1)
    method = method_m.group(1)
    url = url_m.group(1)
    
    # Check if reqParams is empty
    req_empty = "reqParams:[]" in item
    res_empty = "resParams:[]" in item
    
    if not req_empty and not res_empty:
        return item  # Already has params
    
    # Build req params if empty
    req_additions = []
    if req_empty:
        # Path params from URL
        for m in re.finditer(r':(\w+)', url):
            pname = m.group(1)
            ptype = 'int' if 'id' in pname else 'string'
            req_additions.append(f"{{n:'{pname}',t:'{ptype}',r:'必填',d:'{pname}',e:'1'}}")
        
        # Common query params for list endpoints
        if method == 'GET' and ('/list' in url or '/logs' in url or '/records' in url):
            req_additions.append("{n:'page',t:'int',r:'选填',d:'页码',e:'1'}")
            req_additions.append("{n:'page_size',t:'int',r:'选填',d:'每页数量',e:'20'}")
        
        # Common params based on URL patterns
        if '/login' in url and method == 'POST':
            if not any('account' in a for a in req_additions):
                req_additions.append("{n:'account',t:'string',r:'必填',d:'用户名/邮箱/手机号',e:'admin'}")
                req_additions.append("{n:'password',t:'string',r:'必填',d:'密码',e:'Abc123456'}")
        elif '/register' in url and method == 'POST':
            if not any('username' in a for a in req_additions):
                req_additions.append("{n:'username',t:'string',r:'必填',d:'用户名',e:'testuser'}")
                req_additions.append("{n:'password',t:'string',r:'必填',d:'密码',e:'Abc123456'}")
        elif '/refresh' in url and method == 'POST':
            req_additions.append("{n:'refresh_token',t:'string',r:'必填',d:'刷新令牌',e:'eyJ...'}")
        elif '/validate' in url and method == 'POST':
            if 'geetest' in url:
                req_additions.append("{n:'lot_number',t:'string',r:'必填',d:'流水号',e:'lot_xxx'}")
                req_additions.append("{n:'captcha_output',t:'string',r:'必填',d:'验证输出',e:'xxx'}")
                req_additions.append("{n:'pass_token',t:'string',r:'必填',d:'验证令牌',e:'pass_xxx'}")
                req_additions.append("{n:'gen_time',t:'string',r:'必填',d:'生成时间',e:'1720000000'}")
        elif '/verify' in url and method == 'POST':
            if 'captcha' in url:
                req_additions.append("{n:'captcha_id',t:'string',r:'必填',d:'验证码ID',e:'abc123'}")
                req_additions.append("{n:'answer',t:'string',r:'必填',d:'验证码答案',e:'x5km'}")
        elif '/reply' in url and method == 'POST':
            req_additions.append("{n:'content',t:'string',r:'必填',d:'回复内容',e:'回复内容'}")
        elif '/close' in url and method == 'POST':
            req_additions.append("{n:'reason',t:'string',r:'选填',d:'关闭原因',e:'已解决'}")
        elif '/pay' in url and method == 'POST':
            req_additions.append("{n:'payment_method_id',t:'int',r:'必填',d:'支付方式ID',e:'1'}")
        elif '/cancel' in url and method == 'POST':
            req_additions.append("{n:'reason',t:'string',r:'选填',d:'取消原因',e:'不需要了'}")
        elif '/upload' in url and method == 'POST':
            req_additions.append("{n:'file',t:'file',r:'必填',d:'文件',e:'file.bin'}")
        elif '/read' in url and method == 'POST':
            if 'all' not in url:
                req_additions.append("{n:'id',t:'int',r:'必填',d:'ID',e:'1'}")
        elif method == 'POST' and not req_additions:
            # Generic POST
            if ':id' in url:
                for m in re.finditer(r':(\w+)', url):
                    req_additions.append(f"{{n:'{m.group(1)}',t:'int',r:'必填',d:'{m.group(1)}',e:'1'}}")
        elif method == 'DELETE' and not req_additions:
            if ':id' in url:
                for m in re.finditer(r':(\w+)', url):
                    req_additions.append(f"{{n:'{m.group(1)}',t:'int',r:'必填',d:'{m.group(1)}',e:'1'}}")
    
    # Build res params if empty
    res_additions = []
    if res_empty:
        # Auth responses
        if '/login' in url or '/register' in url or '/refresh' in url or '/access-token' in url:
            res_additions.append("{n:'data.access_token',t:'string',r:'',d:'访问令牌',e:''}")
            res_additions.append("{n:'data.refresh_token',t:'string',r:'',d:'刷新令牌',e:''}")
            res_additions.append("{n:'data.user',t:'object',r:'',d:'用户信息',e:''}")
        elif '/oauth' in url and '/login' in url:
            res_additions.append("{n:'data.login_url',t:'string',r:'',d:'OAuth登录URL',e:''}")
            res_additions.append("{n:'data.state',t:'string',r:'',d:'状态码',e:''}")
        elif '/oauth' in url and '/callback' in url:
            res_additions.append("{n:'data.access_token',t:'string',r:'',d:'访问令牌',e:''}")
            res_additions.append("{n:'data.user',t:'object',r:'',d:'用户信息',e:''}")
        elif '/oauth' in url and '/providers' in url:
            res_additions.append("{n:'data.providers',t:'array',r:'',d:'提供商列表',e:''}")
        elif '/oauth' in url and '/accounts' in url:
            res_additions.append("{n:'data.accounts',t:'array',r:'',d:'已绑定账号列表',e:''}")
        elif '/captcha' in url and '/generate' in url:
            res_additions.append("{n:'data.captcha_id',t:'string',r:'',d:'验证码ID',e:''}")
            res_additions.append("{n:'data.image',t:'string',r:'',d:'验证码图片Base64',e:''}")
        elif '/captcha' in url and '/config' in url:
            res_additions.append("{n:'data.captcha_type',t:'string',r:'',d:'验证码类型',e:''}")
            res_additions.append("{n:'data.enabled',t:'bool',r:'',d:'是否启用',e:''}")
        elif '/geetest' in url and '/register' in url:
            res_additions.append("{n:'data.captcha_id',t:'string',r:'',d:'验证ID',e:''}")
            res_additions.append("{n:'data.lot_number',t:'string',r:'',d:'流水号',e:''}")
            res_additions.append("{n:'data.gt',t:'string',r:'',d:'极验GT',e:''}")
            res_additions.append("{n:'data.challenge',t:'string',r:'',d:'挑战码',e:''}")
        elif '/verify' in url:
            res_additions.append("{n:'data.success',t:'bool',r:'',d:'验证结果',e:''}")
            res_additions.append("{n:'message',t:'string',r:'',d:'提示消息',e:''}")
        elif '/check' in url:
            res_additions.append("{n:'data.success',t:'bool',r:'',d:'检查结果',e:''}")
            res_additions.append("{n:'message',t:'string',r:'',d:'提示消息',e:''}")
        elif '/validate' in url:
            res_additions.append("{n:'data.success',t:'bool',r:'',d:'验证结果',e:''}")
            res_additions.append("{n:'message',t:'string',r:'',d:'提示消息',e:''}")
        elif '/close' in url:
            res_additions.append("{n:'message',t:'string',r:'',d:'提示消息',e:''}")
        elif '/pay' in url:
            res_additions.append("{n:'data.pay_url',t:'string',r:'',d:'支付URL',e:''}")
            res_additions.append("{n:'data.trade_no',t:'string',r:'',d:'交易号',e:''}")
        elif '/cancel' in url:
            res_additions.append("{n:'message',t:'string',r:'',d:'提示消息',e:''}")
        elif '/read' in url:
            res_additions.append("{n:'message',t:'string',r:'',d:'提示消息',e:''}")
        elif '/upload' in url:
            res_additions.append("{n:'data.url',t:'string',r:'',d:'文件URL',e:''}")
            res_additions.append("{n:'data.file_id',t:'int',r:'',d:'文件ID',e:''}")
        elif '/reply' in url:
            res_additions.append("{n:'data.id',t:'int',r:'',d:'回复ID',e:''}")
            res_additions.append("{n:'message',t:'string',r:'',d:'提示消息',e:''}")
        elif '/open' in url:
            res_additions.append("{n:'data.is_active',t:'bool',r:'',d:'是否启用',e:''}")
            res_additions.append("{n:'message',t:'string',r:'',d:'提示消息',e:''}")
        elif '/disable' in url:
            res_additions.append("{n:'data.is_active',t:'bool',r:'',d:'是否启用',e:''}")
            res_additions.append("{n:'message',t:'string',r:'',d:'提示消息',e:''}")
        elif '/enable' in url:
            res_additions.append("{n:'data.is_active',t:'bool',r:'',d:'是否启用',e:''}")
            res_additions.append("{n:'message',t:'string',r:'',d:'提示消息',e:''}")
        elif '/bind' in url:
            res_additions.append("{n:'message',t:'string',r:'',d:'提示消息',e:''}")
        elif '/unbind' in url:
            res_additions.append("{n:'message',t:'string',r:'',d:'提示消息',e:''}")
        elif method == 'GET':
            # GET endpoints typically return data
            if ':id' in url:
                # Detail endpoint
                parts = url.rstrip('/').split('/')
                model_name = parts[-2] if parts[-1].startswith(':') else parts[-1]
                res_additions.append(f"{{n:'data',t:'object',r:'',d:'{model_name}详情',e:''}}")
            else:
                # List or other GET
                if '/list' in url or '/logs' in url or '/records' in url:
                    res_additions.append("{n:'data.list',t:'array',r:'',d:'列表',e:''}")
                    res_additions.append("{n:'data.total',t:'int',r:'',d:'总数',e:''}")
                else:
                    res_additions.append("{n:'data',t:'object',r:'',d:'返回数据',e:''}")
        elif method in ('POST', 'PUT', 'PATCH'):
            res_additions.append("{n:'message',t:'string',r:'',d:'提示消息',e:''}")
        elif method == 'DELETE':
            res_additions.append("{n:'message',t:'string',r:'',d:'提示消息',e:''}")
        else:
            res_additions.append("{n:'data',t:'object',r:'',d:'返回数据',e:''}")
    
    # Rebuild item
    if req_additions and req_empty:
        item = item.replace('reqParams:[]', 'reqParams:[' + ','.join(req_additions) + ']')
        # Also update reqExample
        if method in ('POST', 'PUT', 'PATCH') and req_additions:
            body_params = []
            query_params = []
            for ra in req_additions:
                nm = re.search(r"n:'(\w+)'", ra)
                em = re.search(r"e:'([^']*)'", ra)
                src_m = re.search(r"d:'([^']*)'", ra)
                if nm and em:
                    pname = nm.group(1)
                    pval = em.group(1)
                    desc = src_m.group(1) if src_m else pname
                    if any(x in desc for x in ['选填']) and pname in ('page','page_size','keyword','status','type','sort','order'):
                        query_params.append(f'{pname}={pval}')
                    else:
                        body_params.append(f'"{pname}": "{pval}"')
            actual_url = url
            for m in re.finditer(r':(\w+)', url):
                actual_url = actual_url.replace(':'+m.group(1), '1')
            if body_params:
                new_ex = f"{method} {actual_url}"
                if query_params: new_ex += '?' + '&'.join(query_params)
                new_ex += f"\\nAuthorization: JWT eyJ...\\n\\n" + '{' + ', '.join(body_params) + '}'
                item = re.sub(r"reqExample:'[^']*'", f"reqExample:'{new_ex}'", item)
    
    if res_additions and res_empty:
        item = item.replace('resParams:[]', 'resParams:[' + ','.join(res_additions) + ']')
        # Update resExample
        data = {}
        for ra in res_additions:
            nm = re.search(r"n:'([^']*)'", ra)
            tm = re.search(r"t:'([^']*)'", ra)
            if nm:
                k = nm.group(1).replace('data.', '')
                t = tm.group(1) if tm else 'object'
                if '.' not in k:
                    if t == 'array': data[k] = '...'
                    elif t == 'int': data[k] = 0
                    elif t == 'string': data[k] = ''
                    elif t == 'number': data[k] = 0.0
                    elif t == 'bool': data[k] = True
                    else: data[k] = '{}'
        if data:
            new_res = json.dumps({'code':0,'data':data,'message':'success'}, ensure_ascii=False)
            item = re.sub(r"resExample:'[^']*'", f"resExample:'{new_res}'", item)
        else:
            item = re.sub(r"resExample:'[^']*'", "resExample:'{\"code\":0,\"data\":{},\"message\":\"success\"}'", item)
    
    return item

# Apply fixes
new_content = re.sub(
    r"\{id:'[^']*',title:'[^']*',method:'[^']*',url:'[^']*',desc:'[^']*',reqParams:\[.*?\],resParams:\[.*?\],reqExample:'[^']*',resExample:'[^']*'\}",
    fix_item,
    content,
    flags=re.DOTALL
)

with open(js_path, 'w', encoding='utf-8') as f:
    f.write(new_content)

# Count results
total = len(re.findall(r"id:'", new_content))
with_req = len(re.findall(r"reqParams:\[(?!\])", new_content))
with_res = len(re.findall(r"resParams:\[(?!\])", new_content))
print(f"Total: {total}, With req: {with_req} ({with_req*100//total}%), With res: {with_res} ({with_res*100//total}%)")

# Replace in index.html
from replace_api_data import main as replace_main
replace_main()
print("Done!")
