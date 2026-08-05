import re, os

html_path = os.path.join(os.path.dirname(os.path.abspath(__file__)), 'index.html')
print(f"Reading: {html_path}")

with open(html_path, 'r', encoding='utf-8') as f:
    content = f.read()

print(f"Original size: {len(content)} bytes")

# First fix: remove double single-quotes from previous bad run
# Pattern: reqExample:'...'...'' -> reqExample:'...'
# Be careful: only fix double quotes at end of example strings
content = re.sub(
    r"(reqExample|resExample):'((?:[^'\\]|\\.)*)''",
    lambda m: f"{m.group(1)}:'{m.group(2)}'",
    content,
    flags=re.DOTALL
)

# Second fix: escape literal newlines in remaining strings
def fix_newlines(m):
    prefix = m.group(1)
    body = m.group(2)
    body = body.replace('\n', '\\n').replace('\r', '')
    return f"{prefix}:'{body}'"

content = re.sub(
    r"(reqExample|resExample):'((?:[^'\\]|\\.)*)'",
    fix_newlines,
    content,
    flags=re.DOTALL
)

print(f"New size: {len(content)} bytes")

with open(html_path, 'w', encoding='utf-8') as f:
    f.write(content)

# Verify
test = content[content.find('reqExample:'):content.find('reqExample:')+200]
print(f"Sample: {test[:100]}...")
print("Done!")
