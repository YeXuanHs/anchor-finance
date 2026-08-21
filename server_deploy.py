#!/usr/bin/env python3
"""
服务器端部署脚本
在服务器上执行：拉取代码、构建、重启服务
使用paramiko连接服务器
"""

import paramiko
import sys

# 服务器配置
SERVER_HOST = "45.207.210.235"
SERVER_USER = "root"
SERVER_PASS = "iswlBSLY8118"
SERVER_PORT = 22

# 项目配置
REMOTE_DIR = "/opt/anchor-finance"
GITHUB_REPO = "https://github.com/YeXuanHs/anchor-finance.git"


def create_ssh_client():
    """创建SSH客户端"""
    client = paramiko.SSHClient()
    client.set_missing_host_key_policy(paramiko.AutoAddPolicy())
    client.connect(SERVER_HOST, port=SERVER_PORT, username=SERVER_USER, password=SERVER_PASS)
    return client


def execute_command(client, command, description="", check_error=True):
    """执行远程命令"""
    if description:
        print(f"\n[{'='*50}]")
        print(f"[执行] {description}")
        print(f"[命令] {command}")
    
    stdin, stdout, stderr = client.exec_command(command)
    exit_status = stdout.channel.recv_exit_status()
    
    output = stdout.read().decode('utf-8', errors='ignore')
    error = stderr.read().decode('utf-8', errors='ignore')
    
    if output:
        print(output)
    
    if exit_status != 0 and check_error:
        print(f"[错误] 命令执行失败，退出码: {exit_status}")
        if error:
            print(f"[错误输出] {error}")
        return False, output, error
    
    return True, output, error


def pull_latest_code(client):
    """拉取最新代码"""
    # 检查目录是否存在
    success, _, _ = execute_command(
        client,
        f"test -d {REMOTE_DIR}/.git && echo 'exists' || echo 'not_exists'",
        "检查项目目录"
    )
    
    # 如果目录不存在，克隆仓库
    success, output, _ = execute_command(
        client,
        f"test -d {REMOTE_DIR}/.git && echo 'exists' || echo 'not_exists'",
        check_error=False
    )
    
    if 'not_exists' in output:
        print("[克隆] 项目目录不存在或无.git，清理后克隆仓库...")
        # 清理目录
        execute_command(
            client,
            f"rm -rf {REMOTE_DIR}",
            "清理旧目录",
            check_error=False
        )
        # 克隆
        success, _, _ = execute_command(
            client,
            f"git clone {GITHUB_REPO} {REMOTE_DIR}",
            "克隆GitHub仓库"
        )
        if not success:
            return False
    else:
        print("[更新] 拉取最新代码...")
        # 先重置本地修改，再拉取
        success, _, _ = execute_command(
            client,
            f"cd {REMOTE_DIR} && git reset --hard origin/main && git pull origin main --rebase",
            "拉取最新代码"
        )
        if not success:
            # 如果还是失败，强制重置
            success, _, _ = execute_command(
                client,
                f"cd {REMOTE_DIR} && git fetch origin && git reset --hard origin/main",
                "强制重置到远程版本"
            )
        if not success:
            return False
    
    return True


def check_and_install_go(client):
    """检查并安装Go"""
    success, output, _ = execute_command(
        client,
        "go version 2>/dev/null || echo 'not_installed'",
        "检查Go版本",
        check_error=False
    )
    
    if 'not_installed' in output or 'command not found' in output:
        print("[安装] Go未安装，正在安装...")
        commands = [
            "wget -q https://go.dev/dl/go1.21.6.linux-amd64.tar.gz -O /tmp/go.tar.gz",
            "rm -rf /usr/local/go && tar -C /usr/local -xzf /tmp/go.tar.gz",
            "echo 'export PATH=$PATH:/usr/local/go/bin' >> /etc/profile",
            "ln -sf /usr/local/go/bin/go /usr/local/bin/go",
            "go version",
        ]
        for cmd in commands:
            success, _, _ = execute_command(client, cmd, f"安装Go: {cmd[:50]}...")
            if not success:
                return False
    
    return True


def build_backend(client):
    """构建Go后端"""
    print("\n[构建] Go后端...")
    
    # 设置GOTOOLCHAIN=local避免下载新版本
    commands = [
        f"cd {REMOTE_DIR}/backend && export GOTOOLCHAIN=local && /usr/local/go/bin/go mod tidy",
        f"cd {REMOTE_DIR}/backend && export GOTOOLCHAIN=local && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 /usr/local/go/bin/go build -o {REMOTE_DIR}/server/anchor-finance ./cmd/server/",
    ]
    
    for cmd in commands:
        success, _, _ = execute_command(client, cmd, f"构建: {cmd.split('&&')[-1].strip()[:50]}...")
        if not success:
            return False
    
    return True


def setup_systemd(client):
    """配置systemd服务"""
    print("\n[配置] systemd服务...")
    
    service_content = """[Unit]
Description=AnchorFinance Server
After=network.target mysql.service

[Service]
Type=simple
User=root
WorkingDirectory=/opt/anchor-finance/server
ExecStart=/opt/anchor-finance/server/anchor-finance
Restart=always
RestartSec=5
Environment=DB_HOST=127.0.0.1
Environment=DB_PORT=3306
Environment=DB_USER=root
Environment=DB_PASSWORD=iswlBSLY8118
Environment=DB_NAME=anchor_finance
Environment=SERVER_PORT=8080
Environment=SERVER_MODE=release
Environment=JWT_SECRET=anchor-finance-secret-key-2024

[Install]
WantedBy=multi-user.target
"""
    
    # 写入service文件
    sftp = client.open_sftp()
    with sftp.file('/etc/systemd/system/anchor-finance.service', 'w') as f:
        f.write(service_content)
    sftp.close()
    
    commands = [
        "systemctl daemon-reload",
        "systemctl enable anchor-finance",
    ]
    
    for cmd in commands:
        execute_command(client, cmd, cmd)
    
    return True


def restart_service(client):
    """重启服务"""
    print("\n[重启] anchor-finance服务...")
    
    commands = [
        "systemctl restart anchor-finance",
        "sleep 2",
        "systemctl status anchor-finance --no-pager",
    ]
    
    for cmd in commands:
        success, _, _ = execute_command(client, cmd, cmd, check_error=False)
    
    return True


def check_service_status(client):
    """检查服务状态"""
    print("\n[检查] 服务状态...")
    
    # 检查进程
    execute_command(client, "ps aux | grep anchor-finance | grep -v grep", "检查进程", check_error=False)
    
    # 检查端口
    execute_command(client, "netstat -tlnp | grep 8080", "检查端口8080", check_error=False)
    
    # 测试API
    execute_command(client, "curl -s http://localhost:8080/health || echo 'API未响应'", "测试健康检查API", check_error=False)
    
    return True


def main():
    """主函数"""
    print("="*60)
    print("锚点财务 - 服务器部署脚本")
    print("="*60)
    
    try:
        # 连接服务器
        print("\n[连接] 正在连接服务器...")
        client = create_ssh_client()
        print("[成功] 已连接到服务器")
        
        # 拉取最新代码
        if not pull_latest_code(client):
            print("\n[失败] 拉取代码失败")
            client.close()
            return 1
        
        # 检查并安装Go
        if not check_and_install_go(client):
            print("\n[失败] Go安装失败")
            client.close()
            return 1
        
        # 构建后端
        if not build_backend(client):
            print("\n[失败] 后端构建失败")
            client.close()
            return 1
        
        # 配置systemd
        setup_systemd(client)
        
        # 重启服务
        restart_service(client)
        
        # 检查状态
        check_service_status(client)
        
        print("\n" + "="*60)
        print("[完成] 部署完成！")
        print(f"[地址] http://{SERVER_HOST}:8080")
        print(f"[健康检查] http://{SERVER_HOST}:8080/health")
        print("="*60)
        
        client.close()
        return 0
        
    except Exception as e:
        print(f"\n[错误] {str(e)}")
        import traceback
        traceback.print_exc()
        return 1


if __name__ == "__main__":
    sys.exit(main())
