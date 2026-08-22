import re

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
    'backend/internal/api/admin/content_manage.go',
    'backend/internal/api/admin/coupons.go',
    'backend/internal/api/admin/custom_fields.go',
    'backend/internal/api/admin/log_cleanup.go',
    'backend/internal/api/admin/staff.go',
    'backend/internal/api/admin/orders.go',
    'backend/internal/api/admin/products.go',
    'backend/internal/api/admin/services.go',
    'backend/internal/api/admin/plugins.go',
    'backend/internal/api/admin/settings.go',
    'backend/internal/api/admin/suppliers.go',
    'backend/internal/api/admin/dashboard.go',
    'backend/internal/api/admin/invoices.go',
    'backend/internal/api/admin/customers.go',
    'backend/internal/api/admin/configurable_options.go',
    'backend/internal/api/admin/custom_template_fields.go',
    'backend/internal/api/admin/currencies.go',
    'backend/internal/api/admin/credit_limits.go',
    'backend/internal/api/admin/email_templates.go',
    'backend/internal/api/admin/friendly_links.go',
    'backend/internal/api/admin/home_hero.go',
    'backend/internal/api/admin/logs_detail.go',
    'backend/internal/api/admin/member_levels.go',
    'backend/internal/api/admin/notifications.go',
    'backend/internal/api/admin/oauth_providers.go',
    'backend/internal/api/admin/product_types.go',
    'backend/internal/api/admin/referrals.go',
    'backend/internal/api/admin/ticket_prereplies.go',
    'backend/internal/api/admin/traffic_packages.go',
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
        # 匹配含 err.Error() 的 message 行
        match = re.match(r'^(\s*)"message"\s*:\s*".*".*$', line)
        if match and ('err.Error' in line or '" + ' in line):
            indent = match.group(1)
            # 检查下一行是否是 }) 或 }
            if i + 1 < len(lines):
                next_line = lines[i+1].strip()
                if next_line == '})' or next_line == '}':
                    # 检查这是否是错误响应（前面有 "code": 400/500 等）
                    has_error_code = False
                    for j in range(i-1, max(0, i-5), -1):
                        if '"code"' in lines[j] and any(str(n) in lines[j] for n in [400, 401, 403, 404, 500, 502]):
                            has_error_code = True
                            break

                    # 确保当前行末尾有逗号
                    if has_error_code and not line.rstrip().endswith(','):
                        line = line.rstrip() + ','

                    new_lines.append(line)
                    if has_error_code:
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
