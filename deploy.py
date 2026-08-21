#!/usr/bin/env python3
"""
锚点财务部署脚本
使用paramiko连接服务器进行部署
"""

import paramiko
import os
import sys
import time

# 服务器配置
SERVER_HOST = "45.207.210.235"
SERVER_USER = "root"
SERVER_PASS = "iswlBSLY8118"
SERVER_PORT = 22

# 项目路径
LOCAL_PROJECT_DIR = os.path.dirname(os.path.abspath(__file__))
REMOTE_PROJECT_DIR = "/opt/anchor-finance"


def create_ssh_client():
    """创建SSH客户端"""
    client = paramiko.SSHClient()
    client.set_missing_host_key_policy(paramiko.AutoAddPolicy())
    client.connect(SERVER_HOST, port=SERVER_PORT, username=SERVER_USER, password=SERVER_PASS)
    return client


def execute_command(client, command, description=""):
    """执行远程命令"""
    if description:
        print(f"\n[执行] {description}")
    print(f"$ {command}")
    
    stdin, stdout, stderr = client.exec_command(command)
    exit_status = stdout.channel.recv_exit_status()
    
    output = stdout.read().decode('utf-8')
    error = stderr.read().decode('utf-8')
    
    if output:
        print(output)
    if error:
        print(f"[错误] {error}")
    
    return exit_status, output, error


def install_go(client):
    """安装Go"""
    print("\n[安装] Go 1.21...")
    commands = [
        "wget -q https://go.dev/dl/go1.21.6.linux-amd64.tar.gz -O /tmp/go.tar.gz",
        "rm -rf /usr/local/go && tar -C /usr/local -xzf /tmp/go.tar.gz",
        "echo 'export PATH=$PATH:/usr/local/go/bin' >> /etc/profile",
        "export PATH=$PATH:/usr/local/go/bin",
        "ln -sf /usr/local/go/bin/go /usr/local/bin/go",
        "go version",
    ]
    for cmd in commands:
        execute_command(client, cmd)


def install_nodejs(client):
    """安装Node.js"""
    print("\n[安装] Node.js 20...")
    commands = [
        "curl -fsSL https://deb.nodesource.com/setup_20.x | bash -",
        "apt-get install -y nodejs",
        "node --version",
        "npm --version",
    ]
    for cmd in commands:
        execute_command(client, cmd)


def install_php(client):
    """安装PHP"""
    print("\n[安装] PHP 8.1...")
    commands = [
        "apt-get install -y software-properties-common",
        "add-apt-repository -y ppa:ondrej/php",
        "apt-get update",
        "apt-get install -y php8.1 php8.1-cli php8.1-mbstring php8.1-xml php8.1-curl php8.1-mysql php8.1-zip php8.1-bcmath",
        "php --version",
    ]
    for cmd in commands:
        execute_command(client, cmd)


def install_mysql(client):
    """安装MySQL"""
    print("\n[安装] MySQL 8.0...")
    commands = [
        "apt-get install -y mysql-server",
        "systemctl start mysql",
        "systemctl enable mysql",
        "mysql --version",
    ]
    for cmd in commands:
        execute_command(client, cmd)


def install_redis(client):
    """安装Redis"""
    print("\n[安装] Redis...")
    commands = [
        "apt-get install -y redis-server",
        "systemctl start redis-server",
        "systemctl enable redis-server",
        "redis-server --version",
    ]
    for cmd in commands:
        execute_command(client, cmd)


def install_nginx(client):
    """安装Nginx"""
    print("\n[安装] Nginx...")
    commands = [
        "apt-get install -y nginx",
        "systemctl start nginx",
        "systemctl enable nginx",
        "nginx -v",
    ]
    for cmd in commands:
        execute_command(client, cmd)


def check_server_environment(client):
    """检查服务器环境，缺少则安装"""
    print("\n" + "="*60)
    print("检查服务器环境")
    print("="*60)
    
    # 更新包管理器
    execute_command(client, "apt-get update", "更新包管理器")
    
    # 检查并安装Go
    exit_status, _, _ = execute_command(client, "go version", "检查Go版本")
    if exit_status != 0:
        install_go(client)
    
    # 检查并安装Node.js
    exit_status, _, _ = execute_command(client, "node --version", "检查Node.js版本")
    if exit_status != 0:
        install_nodejs(client)
    
    # 检查并安装PHP
    exit_status, _, _ = execute_command(client, "php --version", "检查PHP版本")
    if exit_status != 0:
        install_php(client)
    
    # 检查并安装MySQL
    exit_status, _, _ = execute_command(client, "mysql --version", "检查MySQL版本")
    if exit_status != 0:
        install_mysql(client)
    
    # 检查并安装Redis
    exit_status, _, _ = execute_command(client, "redis-server --version", "检查Redis版本")
    if exit_status != 0:
        install_redis(client)
    
    # 检查并安装Nginx
    exit_status, _, _ = execute_command(client, "nginx -v", "检查Nginx版本")
    if exit_status != 0:
        install_nginx(client)
    
    return True


def setup_project_structure(client):
    """创建项目目录结构"""
    print("\n" + "="*60)
    print("创建项目目录结构")
    print("="*60)
    
    commands = [
        f"mkdir -p {REMOTE_PROJECT_DIR}",
        f"mkdir -p {REMOTE_PROJECT_DIR}/server",
        f"mkdir -p {REMOTE_PROJECT_DIR}/admin",
        f"mkdir -p {REMOTE_PROJECT_DIR}/client",
        f"mkdir -p {REMOTE_PROJECT_DIR}/plugin-engine",
        f"mkdir -p {REMOTE_PROJECT_DIR}/logs",
    ]
    
    for cmd in commands:
        execute_command(client, cmd)


def upload_files(client):
    """上传项目文件"""
    print("\n" + "="*60)
    print("上传项目文件")
    print("="*60)
    
    sftp = client.open_sftp()
    
    # 上传后端代码
    backend_dir = os.path.join(LOCAL_PROJECT_DIR, "backend")
    if os.path.exists(backend_dir):
        print("[上传] 后端代码...")
        upload_directory(sftp, backend_dir, f"{REMOTE_PROJECT_DIR}/backend")
    
    # 上传前端代码
    admin_dir = os.path.join(LOCAL_PROJECT_DIR, "admin-frontend")
    if os.path.exists(admin_dir):
        print("[上传] Admin前端代码...")
        upload_directory(sftp, admin_dir, f"{REMOTE_PROJECT_DIR}/admin-frontend")
    
    client_dir = os.path.join(LOCAL_PROJECT_DIR, "client-frontend")
    if os.path.exists(client_dir):
        print("[上传] Client前端代码...")
        upload_directory(sftp, client_dir, f"{REMOTE_PROJECT_DIR}/client-frontend")
    
    # 上传插件引擎
    plugin_dir = os.path.join(LOCAL_PROJECT_DIR, "plugin-engine")
    if os.path.exists(plugin_dir):
        print("[上传] 插件引擎代码...")
        upload_directory(sftp, plugin_dir, f"{REMOTE_PROJECT_DIR}/plugin-engine")
    
    sftp.close()


def upload_directory(sftp, local_dir, remote_dir):
    """递归上传目录"""
    try:
        sftp.mkdir(remote_dir)
    except:
        pass
    
    for item in os.listdir(local_dir):
        local_path = os.path.join(local_dir, item)
        remote_path = f"{remote_dir}/{item}"
        
        if os.path.isfile(local_path):
            print(f"  上传: {local_path} -> {remote_path}")
            sftp.put(local_path, remote_path)
        elif os.path.isdir(local_path):
            upload_directory(sftp, local_path, remote_path)


def build_backend(client):
    """构建后端"""
    print("\n" + "="*60)
    print("构建Go后端")
    print("="*60)
    
    commands = [
        f"cd {REMOTE_PROJECT_DIR}/backend && go mod tidy",
        f"cd {REMOTE_PROJECT_DIR}/backend && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o {REMOTE_PROJECT_DIR}/server/anchor-finance ./cmd/server/",
    ]
    
    for cmd in commands:
        exit_status, _, error = execute_command(client, cmd)
        if exit_status != 0:
            print(f"[错误] 构建失败: {error}")
            return False
    
    return True


def build_frontend(client):
    """构建前端"""
    print("\n" + "="*60)
    print("构建前端")
    print("="*60)
    
    # 构建Admin前端
    commands = [
        f"cd {REMOTE_PROJECT_DIR}/admin-frontend && npm install",
        f"cd {REMOTE_PROJECT_DIR}/admin-frontend && npm run build",
    ]
    
    for cmd in commands:
        exit_status, _, error = execute_command(client, cmd)
        if exit_status != 0:
            print(f"[错误] Admin前端构建失败: {error}")
            return False
    
    # 构建Client前端
    commands = [
        f"cd {REMOTE_PROJECT_DIR}/client-frontend && npm install",
        f"cd {REMOTE_PROJECT_DIR}/client-frontend && npm run build",
    ]
    
    for cmd in commands:
        exit_status, _, error = execute_command(client, cmd)
        if exit_status != 0:
            print(f"[错误] Client前端构建失败: {error}")
            return False
    
    return True


def setup_database(client):
    """初始化数据库"""
    print("\n" + "="*60)
    print("初始化数据库")
    print("="*60)
    
    # 创建数据库和用户
    sql_commands = """
CREATE DATABASE IF NOT EXISTS anchor_finance CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER IF NOT EXISTS 'anchor'@'localhost' IDENTIFIED BY 'anchor123';
GRANT ALL PRIVILEGES ON anchor_finance.* TO 'anchor'@'localhost';
FLUSH PRIVILEGES;
"""
    
    # 写入SQL文件
    sftp = client.open_sftp()
    with sftp.file(f"{REMOTE_PROJECT_DIR}/init.sql", 'w') as f:
        f.write(sql_commands)
    sftp.close()
    
    # 执行SQL
    execute_command(client, f"mysql -u root -p{SERVER_PASS} < {REMOTE_PROJECT_DIR}/init.sql", "创建数据库")


def create_systemd_service(client):
    """创建systemd服务"""
    print("\n" + "="*60)
    print("创建systemd服务")
    print("="*60)
    
    service_content = f"""[Unit]
Description=AnchorFinance Server
After=network.target mysql.service redis.service

[Service]
Type=simple
User=root
WorkingDirectory={REMOTE_PROJECT_DIR}/server
ExecStart={REMOTE_PROJECT_DIR}/server/anchor-finance
Restart=always
RestartSec=5
Environment=DB_HOST=127.0.0.1
Environment=DB_PORT=3306
Environment=DB_USER=anchor
Environment=DB_PASSWORD=anchor123
Environment=DB_NAME=anchor_finance
Environment=REDIS_HOST=127.0.0.1
Environment=REDIS_PORT=6379
Environment=SERVER_PORT=8080
Environment=SERVER_MODE=release
Environment=JWT_SECRET=anchor-finance-secret-key-2024

[Install]
WantedBy=multi-user.target
"""
    
    sftp = client.open_sftp()
    with sftp.file('/etc/systemd/system/anchor-finance.service', 'w') as f:
        f.write(service_content)
    sftp.close()
    
    commands = [
        "systemctl daemon-reload",
        "systemctl enable anchor-finance",
    ]
    
    for cmd in commands:
        execute_command(client, cmd)


def setup_nginx(client):
    """配置Nginx"""
    print("\n" + "="*60)
    print("配置Nginx")
    print("="*60)
    
    nginx_config = f"""server {{
    listen 80;
    server_name _;
    
    # Admin前端
    location /admin {{
        alias {REMOTE_PROJECT_DIR}/admin-frontend/dist;
        try_files $uri $uri/ /admin/index.html;
    }}
    
    # 用户前台
    location / {{
        root {REMOTE_PROJECT_DIR}/client-frontend/dist;
        try_files $uri $uri/ /index.html;
    }}
    
    # API代理
    location /api/ {{
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }}
    
    # 插件引擎代理
    location /plugin-engine/ {{
        proxy_pass http://127.0.0.1:9000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }}
}}
"""
    
    sftp = client.open_sftp()
    with sftp.file('/etc/nginx/sites-available/anchor-finance', 'w') as f:
        f.write(nginx_config)
    sftp.close()
    
    commands = [
        "ln -sf /etc/nginx/sites-available/anchor-finance /etc/nginx/sites-enabled/",
        "nginx -t",
        "systemctl reload nginx",
    ]
    
    for cmd in commands:
        execute_command(client, cmd)


def start_services(client):
    """启动服务"""
    print("\n" + "="*60)
    print("启动服务")
    print("="*60)
    
    commands = [
        "systemctl restart anchor-finance",
        "systemctl status anchor-finance",
    ]
    
    for cmd in commands:
        execute_command(client, cmd)


def main():
    """主函数"""
    print("="*60)
    print("锚点财务部署脚本")
    print("="*60)
    
    try:
        # 创建SSH连接
        print("\n连接服务器...")
        client = create_ssh_client()
        print("[成功] 已连接到服务器")
        
        # 检查环境
        check_server_environment(client)
        
        # 创建目录结构
        setup_project_structure(client)
        
        # 上传文件
        upload_files(client)
        
        # 构建后端
        if not build_backend(client):
            print("\n[失败] 后端构建失败")
            client.close()
            return
        
        # 构建前端
        if not build_frontend(client):
            print("\n[失败] 前端构建失败")
            client.close()
            return
        
        # 初始化数据库
        setup_database(client)
        
        # 创建systemd服务
        create_systemd_service(client)
        
        # 配置Nginx
        setup_nginx(client)
        
        # 启动服务
        start_services(client)
        
        print("\n" + "="*60)
        print("部署完成！")
        print("="*60)
        print(f"后台地址: http://{SERVER_HOST}/admin")
        print(f"前台地址: http://{SERVER_HOST}/")
        print(f"API地址: http://{SERVER_HOST}:8080")
        
        client.close()
        
    except Exception as e:
        print(f"\n[错误] {str(e)}")
        import traceback
        traceback.print_exc()


if __name__ == "__main__":
    main()
