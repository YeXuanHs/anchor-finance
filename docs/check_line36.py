import os

html_path = os.path.join(os.path.dirname(os.path.abspath(__file__)), 'index.html')
with open(html_path, 'r', encoding='utf-8') as f:
    lines = f.readlines()

for i in range(32, 42):
    line = lines[i]
    issues = []
    for j, ch in enumerate(line):
        if ord(ch) < 32 and ch not in ('\n', '\r', '\t'):
            issues.append(f"col {j}: U+{ord(ch):04X}")
    if issues:
        print(f"Line {i+1}: BAD - {', '.join(issues)}")
    else:
        print(f"Line {i+1}: OK")
