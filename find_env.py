import paramiko

client = paramiko.SSHClient()
client.set_missing_host_key_policy(paramiko.AutoAddPolicy())
client.connect('45.207.210.235', username='root', password='iswlBSLY8118')

# 读取旧版锚点财务的env
stdin, stdout, stderr = client.exec_command("cat /opt/anchorfinance/.env")
print("旧版.env内容:")
print(stdout.read().decode())

client.close()
