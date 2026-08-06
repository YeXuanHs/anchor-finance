import paramiko
import sys

def ssh_exec(host, username, password, command):
    client = paramiko.SSHClient()
    client.set_missing_host_key_policy(paramiko.AutoAddPolicy())
    client.connect(host, username=username, password=password, timeout=30)
    
    stdin, stdout, stderr = client.exec_command(command, timeout=600)
    out = stdout.read().decode('utf-8', errors='replace')
    err = stderr.read().decode('utf-8', errors='replace')
    exit_code = stdout.channel.recv_exit_status()
    
    client.close()
    return out, err, exit_code

if __name__ == "__main__":
    host = "45.207.210.235"
    username = "root"
    password = "iswlBSLY8118"
    
    # Pull and build
    build_cmd = """
cd /opt/anchorfinance
git pull origin master
cd backend
export PATH=$PATH:/usr/local/go/bin
CGO_ENABLED=0 go build -ldflags="-s -w" -o /opt/anchorfinance/anchorfinance .
"""
    
    print("Pulling and building...")
    out, err, code = ssh_exec(host, username, password, build_cmd)
    if out:
        print(out)
    if err:
        print(err, file=sys.stderr)
    
    if code != 0:
        print(f"Build failed with exit code {code}")
        sys.exit(code)
    
    print("Build successful!")
