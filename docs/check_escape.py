import os

html_path = os.path.join(os.path.dirname(os.path.abspath(__file__)), 'index.html')
with open(html_path, 'rb') as f:
    content = f.read()

# Find first reqExample and check what's between the quotes
idx = content.find(b"reqExample:'")
if idx < 0:
    print("NOT FOUND")
else:
    start = idx + len(b"reqExample:'")
    # Find the closing quote - look for ',res or '}
    end = content.find(b"',", start)
    if end < 0:
        end = content.find(b"'}", start)
    chunk = content[start:end]
    print(f"reqExample length: {len(chunk)} bytes")
    print(f"Hex dump of first 100 bytes:")
    print(chunk[:100].hex(' '))
    print(f"\nAs text:")
    print(chunk[:200])
    
    # Check for literal backslash-n vs actual newline
    if b'\n' in chunk:
        print("\nWARNING: Contains actual newline bytes!")
    if b'\\n' in chunk:
        print("\nOK: Contains \\n escape sequences")
