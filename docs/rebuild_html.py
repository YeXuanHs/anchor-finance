import os

base = os.path.dirname(os.path.abspath(__file__))
html_path = os.path.join(base, 'index.html')

with open(html_path, 'r', encoding='utf-8') as f:
    lines = f.readlines()

# Find the structure:
# Line 31: <script src="api_data.js"></script>
# Line 32: <script>
# Line 33: // ========== API Data loaded from api_data.js ==========
# Line 34+: OLD API_DATA content...
# Some line: };
# Next line: // ========== Mobile Menu ==========  (or similar)
# ...
# Last lines: </script></body></html>

# Find the start of old data (line 34, 0-indexed 33)
start_remove = 33  # "  public: { title: ..."

# Find the end - look for "};\n" followed by empty line and "// =========="
end_remove = -1
for i in range(start_remove, len(lines)):
    if lines[i].strip() == '};':
        # Check if next non-empty line is a comment
        for j in range(i+1, min(i+3, len(lines))):
            if lines[j].strip().startswith('// =========='):
                end_remove = i + 1  # include the "};\n"
                break
        if end_remove > 0:
            break

if end_remove < 0:
    # Fallback: just find "};"
    for i in range(start_remove, len(lines)):
        if lines[i].strip() == '};':
            end_remove = i + 1
            break

print(f"Removing lines {start_remove+1} to {end_remove} (0-indexed {start_remove} to {end_remove-1})")
print(f"Total lines before: {len(lines)}")

# Build new file
new_lines = lines[:start_remove] + ['\n'] + lines[end_remove:]
print(f"Total lines after: {len(new_lines)}")

with open(html_path, 'w', encoding='utf-8') as f:
    f.writelines(new_lines)

print("Done! Verifying...")

# Verify the structure
with open(html_path, 'r', encoding='utf-8') as f:
    content = f.read()

if 'api_data.js' in content and 'var API_DATA' not in content and 'function render()' in content:
    print("OK: Uses external api_data.js, no inline API_DATA, render function present")
else:
    print("WARNING: Structure may be incorrect")
    if 'var API_DATA' in content:
        print("  - Still has inline API_DATA")
    if 'function render()' not in content:
        print("  - Missing render function")
