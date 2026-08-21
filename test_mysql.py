import paramiko

client = paramiko.SSHClient()
client.set_missing_host_key_policy(paramiko.AutoAddPolicy())
client.connect('45.207.210.235', username='root', password='iswlBSLY8118')

# 测试MySQL连接
stdin, stdout, stderr = client.exec_command("mysql -u root -p'iswlBSLY8118' -e 'SELECT 1;' 2>&1")
print(stdout.read().decode())

client.close()
