import os

html_path = os.path.join(os.path.dirname(os.path.abspath(__file__)), 'index.html')
with open(html_path, 'rb') as f:
    content = f.read()

# Find the area around line 36
lines = content.split(b'\n')
print(f"Total lines: {len(lines)}")

# Check line 36 (0-indexed: 35)
line36 = lines[35]
print(f"Line 36 length: {len(line36)}")
print(f"Line 36 first 20 bytes: {line36[:20]}")
print(f"Line 36 hex: {line36[:20].hex()}")

# Check for BOM or special chars at start
for i in range(min(10, len(line36))):
    b = line36[i]
    if b < 32 or b > 127:
        print(f"  Special byte at pos {i}: 0x{b:02X}")

# Also check the area around the first reqExample
idx = content.find(b"reqExample:")
if idx >= 0:
    chunk = content[idx:idx+200]
    print(f"\nFirst reqExample area ({len(chunk)} bytes):")
    print(chunk)
