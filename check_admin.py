import paramiko

client = paramiko.SSHClient()
client.set_missing_host_key_policy(paramiko.AutoAddPolicy())
client.connect('45.207.210.235', username='root', password='iswlBSLY8118')

# 查询管理员信息
cmd = """
mysql -u root -p'*RhbY#m0IbS8SPaAteOI' anchor_finance -e 'SELECT id, username, email, status, login_fail_count, locked_until FROM admins;'
"""
stdin, stdout, stderr = client.exec_command(cmd)
print("管理员账号:")
print(stdout.read().decode('utf-8', errors='ignore'))

client.close()
