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
    
    if len(sys.argv) < 2:
        print("Usage: python ssh.py <command>")
        sys.exit(1)
    
    command = " ".join(sys.argv[1:])
    out, err, code = ssh_exec(host, username, password, command)
    
    if out:
        print(out)
    if err:
        print(err, file=sys.stderr)
    sys.exit(code)
