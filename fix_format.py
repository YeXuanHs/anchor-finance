import re

# 精确修复：找到 "message": "纯字符串", 后面直接是 \n\t}) 的位置，加 "data": nil
# 规则：只处理纯字符串（无+拼接、无err.Error），且 message 后面直接是 \n\t})

count = 0
files_fixed = set()

for filepath in [
    'backend/internal/api/client/auth.go',
    'backend/internal/api/client/invoices.go',
    'backend/internal/api/client/orders.go',
    'backend/internal/api/client/services.go',
    'backend/internal/api/client/services_manage.go',
    'backend/internal/api/client/tickets.go',
    'backend/internal/api/client/payments.go',
    'backend/internal/api/client/referral.go',
    'backend/internal/api/client/finance.go',
    'backend/internal/api/client/verification.go',
    'backend/internal/api/admin/admins.go',
    'backend/internal/api/admin/auth.go',
    'backend/internal/api/admin/coupons.go',
    'backend/internal/api/admin/content_manage.go',
    'backend/internal/api/admin/contracts.go',
    'backend/internal/api/admin/currencies.go',
    'backend/internal/api/admin/custom_fields.go',
    'backend/internal/api/admin/credit_limits.go',
    'backend/internal/api/admin/invoices.go',
    'backend/internal/api/admin/customers.go',
    'backend/internal/api/admin/orders.go',
    'backend/internal/api/admin/products.go',
    'backend/internal/api/admin/services.go',
    'backend/internal/api/admin/plugins.go',
    'backend/internal/api/admin/settings.go',
    'backend/internal/api/admin/staff.go',
    'backend/internal/api/admin/suppliers.go',
    'backend/internal/api/admin/dashboard.go',
    'backend/internal/api/admin/log_cleanup.go',
]:
    try:
        with open(filepath, 'r', encoding='utf-8') as f:
            content = f.read()
    except FileNotFoundError:
        continue

    lines = content.split('\n')
    new_lines = []
    i = 0
    changed = False

    while i < len(lines):
        line = lines[i]
        # 匹配: "message": "纯字符串", （无+拼接、无err.Error）
        match = re.match(r'^(\s*)"message"\s*:\s*"([^"]+)"(,?)\s*$', line)
        if match and '+' not in line and 'err.Error' not in line:
            indent = match.group(1)
            # 检查下一行是否是 }) 或 }
            if i + 1 < len(lines):
                next_line = lines[i+1].strip()
                if next_line == '})' or next_line == '}':
                    # 这是错误响应的末尾，需要在 message 前面的 gin.H 里加 "data": nil
                    # 找到前面的 "code" 行来确认是错误响应
                    code_line = None
                    for j in range(i-1, max(0, i-5), -1):
                        if '"code"' in lines[j] and any(str(n) in lines[j] for n in [400, 401, 403, 404, 500, 502]):
                            code_line = lines[j]
                            break

                    if code_line and '"data"' not in code_line:
                        # 是错误响应且无data，加 "data": nil
                        new_lines.append(line)
                        new_lines.append(f'{indent}"data": nil,')
                        count += 1
                        changed = True
                        i += 1
                        continue

        new_lines.append(line)
        i += 1

    if changed:
        with open(filepath, 'w', encoding='utf-8') as f:
            f.write('\n'.join(new_lines))
        files_fixed.add(filepath)

print(f"修复了 {count} 处，涉及 {len(files_fixed)} 个文件")
for f in sorted(files_fixed):
    print(f"  - {f}")
