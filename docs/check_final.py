import re
f = open(r'c:\Users\Administrator\Desktop\智简魔方\AnchorFinance\docs\index.html', 'r', encoding='utf-8')
html = f.read()
f.close()
total = len(re.findall(r"id:'", html))
empty_req = len(re.findall(r'reqParams:\[\]', html))
empty_res = len(re.findall(r'resParams:\[\]', html))
with_req = total - empty_req
with_res = total - empty_res
print(f"Total: {total}")
print(f"With req: {with_req} ({with_req*100//total}%)")
print(f"With res: {with_res} ({with_res*100//total}%)")
print(f"Empty req: {empty_req}")
print(f"Empty res: {empty_res}")
