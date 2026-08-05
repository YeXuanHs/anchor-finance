#!/usr/bin/env python3
"""Replace API_DATA in index.html with data from api_data_complete.js"""
import os, re

DOCS_DIR = os.path.dirname(__file__)
html_path = os.path.join(DOCS_DIR, 'index.html')
js_path = os.path.join(DOCS_DIR, 'api_data_complete.js')

with open(html_path, 'r', encoding='utf-8') as f:
    html = f.read()

with open(js_path, 'r', encoding='utf-8') as f:
    new_data = f.read()

# Replace the var API_DATA = {...}; block
pattern = r'var API_DATA\s*=\s*\{.*?\};'
html = re.sub(pattern, new_data, html, flags=re.DOTALL)

with open(html_path, 'w', encoding='utf-8') as f:
    f.write(html)

print("API_DATA replaced successfully!")
print(f"HTML size: {len(html)} bytes")
