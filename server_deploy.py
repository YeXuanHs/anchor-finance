#!/usr/bin/env python3
"""
服务器部署脚本 - 服务器端直接下载GitHub Release，不做任何编译
"""
import paramiko
import time

SERVER_HOST = "45.207.210.235"
SERVER_USER = "root"
SERVER_PASS = "iswlBSLY8118"
REMOTE_DIR = "/opt/anchor-finance"


def run(client, cmd):
    stdin, stdout, stderr = client.exec_command(cmd, timeout=180)
    exit_status = stdout.channel.recv_exit_status()
    out = stdout.read().decode('utf-8', errors='replace').strip()
    err = stderr.read().decode('utf-8', errors='replace').strip()
    if out:
        print(out)
    if err and exit_status != 0:
        print(f"[stderr] {err[:500]}")
    return out


def main():
    client = paramiko.SSHClient()
    client.set_missing_host_key_policy(paramiko.AutoAddPolicy())
    client.connect(SERVER_HOST, username=SERVER_USER, password=SERVER_PASS)

    # 1. 先停服务，避免 Text file busy
    print("\n=== 停服务 ===")
    run(client, "systemctl stop anchor-finance")
    run(client, "pkill -f 'anchor-finance$' 2>/dev/null || true")
    time.sleep(2)

    # 2. 拉取代码（备份.env防止被覆盖）
    print("\n=== 拉取代码 ===")
    run(client, f"cp {REMOTE_DIR}/server/.env /tmp/anchor-finance.env.bak 2>/dev/null; true")
    run(client, f"cd {REMOTE_DIR} && git fetch origin && git reset --hard origin/main")
    run(client, f"cp /tmp/anchor-finance.env.bak {REMOTE_DIR}/server/.env 2>/dev/null; true")

    # 3. 下载Release二进制
    print("\n=== 下载Release===")
    run(client, (
        f"cd {REMOTE_DIR} && "
        "URL=$(curl -s https://api.github.com/repos/YeXuanHs/anchor-finance/releases/latest "
        "| python3 -c 'import sys,json; d=json.load(sys.stdin); "
        "print([a[\"browser_download_url\"] for a in d.get(\"assets\",[]) "
        "if \"anchor-finance-server\" in a[\"name\"]][0])') "
        "&& curl -L -o server/anchor-finance \"$URL\" "
        "&& chmod +x server/anchor-finance "
        "&& ls -la server/anchor-finance"
    ))

    # 4. 重启（先彻底释放端口，必须daemon-reload加载新环境变量）
    print("\n=== 启动服务 ===")
    run(client, "systemctl stop anchor-finance 2>/dev/null; true")
    run(client, "fuser -k 8080/tcp 2>/dev/null; true")
    run(client, "sleep 3")
    run(client, "systemctl daemon-reload")
    run(client, "systemctl start anchor-finance")
    time.sleep(4)
    status = run(client, "systemctl is-active anchor-finance")
    print(f"服务状态: {status}")
    if status.strip() == "active":
        run(client, "curl -s http://127.0.0.1:8080/health", )
    else:
        run(client, "journalctl -u anchor-finance --no-pager -n 20")

    client.close()
    print("\n=== 部署完成（未编译） ===")


if __name__ == "__main__":
    main()
