<?php
/**
 * AnchorFinance Plugin Engine
 * 原生PHP插件引擎 - 兼容zjmf插件
 */

// 错误处理
error_reporting(E_ALL);
ini_set('display_errors', '0');

// 数据库连接
$db_host = getenv('DB_HOST') ?: '127.0.0.1';
$db_port = getenv('DB_PORT') ?: '3306';
$db_user = getenv('DB_USER') ?: 'root';
$db_pass = getenv('DB_PASS') ?: '';
$db_name = getenv('DB_NAME') ?: 'anchor_finance';

try {
    $pdo = new PDO(
        "mysql:host={$db_host};port={$db_port};dbname={$db_name};charset=utf8mb4",
        $db_user,
        $db_pass,
        [PDO::ATTR_ERRMODE => PDO::ERRMODE_EXCEPTION]
    );
} catch (PDOException $e) {
    json_response(500, '数据库连接失败');
    exit;
}

// 插件目录
$plugins_dir = __DIR__ . '/plugins';

// 引入zjmf兼容层（让zjmf插件可以零修改运行）
require_once __DIR__ . '/zjmf_compat/functions.php';
require_once __DIR__ . '/app/ZjmfCompat/DatabaseCompat.php';
require_once __DIR__ . '/app/ZjmfCompat/HookCompat.php';
require_once __DIR__ . '/app/ZjmfCompat/CmfCompat.php';
require_once __DIR__ . '/app/ZjmfCompat/TemplateEngine.php';

// 初始化DatabaseCompat的PDO连接
\App\ZjmfCompat\DatabaseCompat::setConnection($pdo);

// 加载已启用插件的hooks.php
loadEnabledPluginHooks($pdo, $plugins_dir);

// 路由
$request_uri = $_SERVER['REQUEST_URI'] ?? '/';
$method = $_SERVER['REQUEST_METHOD'];
$path = parse_url($request_uri, PHP_URL_PATH);
$path = preg_replace('#^/index\.php#', '', $path);
if ($path !== '/' && substr($path, -1) === '/') {
    $path = rtrim($path, '/');
}

// JSON响应函数
function json_response($code, $message, $data = null) {
    header('Content-Type: application/json; charset=utf-8');
    echo json_encode([
        'code' => $code,
        'message' => $message,
        'data' => $data
    ], JSON_UNESCAPED_UNICODE);
}

// 获取请求体
function get_json_body() {
    $input = file_get_contents('php://input');
    return json_decode($input, true) ?: [];
}

// 加载已启用插件的hooks.php
function loadEnabledPluginHooks($pdo, $plugins_dir) {
    $stmt = $pdo->query("SELECT slug, domain FROM plugins WHERE status = 'active'");
    $plugins = $stmt->fetchAll(PDO::FETCH_ASSOC);

    foreach ($plugins as $plugin) {
        $hooks_file = $plugins_dir . '/' . $plugin['domain'] . '/' . $plugin['slug'] . '/hooks.php';
        if (file_exists($hooks_file)) {
            $hooks = include $hooks_file;
            if (is_array($hooks)) {
                foreach ($hooks as $tag => $config) {
                    if (isset($config['class']) && class_exists($config['class'])) {
                        $handler = [new $config['class'](), 'handle'];
                        \App\ZjmfCompat\HookCompat::add($tag, $handler, $config['priority'] ?? 10);
                    }
                }
            }
        }
    }
}

// 加载单个插件的hooks.php
function loadSinglePluginHooks($pdo, $plugins_dir, $domain, $slug) {
    $hooks_file = $plugins_dir . '/' . $domain . '/' . $slug . '/hooks.php';
    if (!file_exists($hooks_file)) return;

    $hooks = include $hooks_file;
    if (!is_array($hooks)) return;

    foreach ($hooks as $tag => $config) {
        if (isset($config['class']) && class_exists($config['class'])) {
            $handler = [new $config['class'](), 'handle'];
            \App\ZjmfCompat\HookCompat::add($tag, $handler, $config['priority'] ?? 10);
        }
    }
}

// 调用插件方法（install/uninstall等）
function callPluginMethod($plugins_dir, $domain, $slug, $method) {
    $manifest_file = $plugins_dir . '/' . $domain . '/' . $slug . '/manifest.json';
    if (!file_exists($manifest_file)) return;

    $manifest = json_decode(file_get_contents($manifest_file), true);
    $entry_class = $manifest['entry_class'] ?? '';
    if (empty($entry_class)) return;

    // 尝试加载入口类
    $class_file = $plugins_dir . '/' . $domain . '/' . $slug . '/src/' . basename(str_replace('\\', '/', $entry_class)) . '.php';
    if (!file_exists($class_file)) return;

    require_once $class_file;
    if (!class_exists($entry_class)) return;

    $instance = new $entry_class();
    if (method_exists($instance, $method)) {
        $instance->$method();
    }
}

// ==================== 路由分发 ====================

// 健康检查
if ($path === '/health') {
    json_response(0, 'ok', [
        'status' => 'ok',
        'service' => 'AnchorFinance Plugin Engine',
        'version' => '1.0.0',
        'php' => PHP_VERSION,
    ]);
    exit;
}

// 内部API
if (strpos($path, '/internal/') === 0) {
    $sub_path = ltrim(substr($path, 9), '/');

    // ===== Hook触发 =====
    if ($sub_path === 'hook/trigger' && $method === 'POST') {
        $body = get_json_body();
        $hook = $body['hook'] ?? '';
        $params = $body['params'] ?? [];

        if (empty($hook)) {
            json_response(400, 'Hook name is required');
            exit;
        }

        // 通过兼容层触发Hook
        $results = \App\ZjmfCompat\HookCompat::trigger($hook, $params);

        json_response(0, 'success', [
            'hook' => $hook,
            'results' => $results,
        ]);
        exit;
    }

    // ===== Hook列表 =====
    if ($sub_path === 'hook/list' && $method === 'POST') {
        $hooks = \App\ZjmfCompat\HookCompat::getAllHooks();
        $hook_list = [];
        foreach ($hooks as $tag => $handlers) {
            $hook_list[] = [
                'tag' => $tag,
                'handler_count' => count($handlers),
            ];
        }

        json_response(0, 'success', $hook_list);
        exit;
    }

    // ===== Hook注册 =====
    if ($sub_path === 'hook/register' && $method === 'POST') {
        $body = get_json_body();
        $hook_name = $body['tag'] ?? '';
        $plugin_key = $body['plugin_key'] ?? '';
        $handler_class = $body['handler_class'] ?? '';
        $priority = $body['priority'] ?? 10;

        if (empty($hook_name) || empty($handler_class)) {
            json_response(400, 'tag and handler_class are required');
            exit;
        }

        // 保存到数据库
        $stmt = $pdo->prepare("INSERT INTO hook_bindings (hook_name, plugin_key, handler_class, priority, created_at, updated_at) VALUES (?, ?, ?, ?, NOW(), NOW())");
        $stmt->execute([$hook_name, $plugin_key, $handler_class, $priority]);

        // 如果类已存在，立即注册到内存
        if (class_exists($handler_class)) {
            $handler = [new $handler_class(), 'handle'];
            \App\ZjmfCompat\HookCompat::add($hook_name, $handler, $priority);
        }

        json_response(0, 'Hook registered successfully');
        exit;
    }

    // ===== 插件列表（每次请求都读目录，合并数据库状态）=====
    if ($sub_path === 'plugins' && $method === 'GET') {
        $domain_filter = $_GET['domain'] ?? '';
        $domains = ['payment', 'sms', 'mail', 'oauth', 'servers', 'captcha', 'certification', 'addons'];
        $result = [];

        foreach ($domains as $domain) {
            if ($domain_filter && $domain !== $domain_filter) continue;

            $domain_dir = $plugins_dir . '/' . $domain;
            if (!is_dir($domain_dir)) continue;

            $items = scandir($domain_dir);
            foreach ($items as $item) {
                if ($item === '.' || $item === '..') continue;

                $plugin_dir = $domain_dir . '/' . $item;
                if (!is_dir($plugin_dir)) continue;

                $manifest_file = $plugin_dir . '/manifest.json';
                if (!file_exists($manifest_file)) continue;

                $manifest = json_decode(file_get_contents($manifest_file), true);
                if (!$manifest) continue;

                $slug = $manifest['slug'] ?? $item;

                // 查数据库获取安装/启用状态
                $stmt = $pdo->prepare("SELECT id, status, config FROM plugins WHERE slug = ? AND domain = ?");
                $stmt->execute([$slug, $domain]);
                $db_plugin = $stmt->fetch(PDO::FETCH_ASSOC);

                $result[] = [
                    'id' => $db_plugin['id'] ?? null,
                    'name' => $manifest['name'] ?? $item,
                    'slug' => $slug,
                    'domain' => $domain,
                    'description' => $manifest['description'] ?? '',
                    'version' => $manifest['version'] ?? '1.0.0',
                    'author' => $manifest['author'] ?? '',
                    'status' => $db_plugin['status'] ?? 'not_installed',
                    'installed' => $db_plugin ? true : false,
                    'config_schema' => $manifest['config_schema'] ?? [],
                ];
            }
        }

        json_response(0, 'success', [
            'list' => $result,
            'total' => count($result),
        ]);
        exit;
    }

    // ===== 安装插件（数据库记录，不改文件）=====
    if ($sub_path === 'plugins/install' && $method === 'POST') {
        $body = get_json_body();
        $domain = $body['domain'] ?? '';
        $slug = $body['slug'] ?? '';

        if (empty($domain) || empty($slug)) {
            json_response(400, 'domain and slug are required');
            exit;
        }

        // 检查文件是否存在
        $manifest_file = $plugins_dir . '/' . $domain . '/' . $slug . '/manifest.json';
        if (!file_exists($manifest_file)) {
            json_response(404, '插件文件不存在');
            exit;
        }

        $manifest = json_decode(file_get_contents($manifest_file), true);

        // 检查是否已安装
        $stmt = $pdo->prepare("SELECT COUNT(*) FROM plugins WHERE slug = ? AND domain = ?");
        $stmt->execute([$slug, $domain]);
        if ($stmt->fetchColumn() > 0) {
            json_response(400, '插件已安装');
            exit;
        }

        // 读取默认配置
        $config_file = $plugins_dir . '/' . $domain . '/' . $slug . '/config.php';
        $config = '{}';
        if (file_exists($config_file)) {
            $config_data = include $config_file;
            $config = json_encode($config_data, JSON_UNESCAPED_UNICODE);
        }

        $stmt = $pdo->prepare("INSERT INTO plugins (name, slug, domain, description, version, config, status, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, 'inactive', NOW(), NOW())");
        $stmt->execute([
            $manifest['name'] ?? $slug,
            $slug,
            $domain,
            $manifest['description'] ?? '',
            $manifest['version'] ?? '1.0.0',
            $config,
        ]);

        // 调用插件的install方法
        callPluginMethod($plugins_dir, $domain, $slug, 'install');

        json_response(0, '安装成功', ['id' => (int)$pdo->lastInsertId()]);
        exit;
    }

    // ===== 卸载插件 =====
    if (preg_match('#^plugins/(\d+)/uninstall$#', $sub_path, $m) && $method === 'POST') {
        $id = $m[1];
        $stmt = $pdo->prepare("SELECT * FROM plugins WHERE id = ?");
        $stmt->execute([$id]);
        $plugin = $stmt->fetch(PDO::FETCH_ASSOC);

        if (!$plugin) {
            json_response(404, '插件不存在');
            exit;
        }

        // 调用插件的uninstall方法
        callPluginMethod($plugins_dir, $plugin['domain'], $plugin['slug'], 'uninstall');

        // 删除数据库记录
        $stmt = $pdo->prepare("DELETE FROM plugins WHERE id = ?");
        $stmt->execute([$id]);

        // 删除Hook绑定
        $stmt = $pdo->prepare("DELETE FROM hook_bindings WHERE plugin_key = ?");
        $stmt->execute([$plugin['domain'] . '/' . $plugin['slug']]);

        json_response(0, '卸载成功');
        exit;
    }

    // ===== 启用插件 =====
    if (preg_match('#^plugins/(\d+)/enable$#', $sub_path, $m) && $method === 'POST') {
        $id = $m[1];
        $stmt = $pdo->prepare("SELECT * FROM plugins WHERE id = ?");
        $stmt->execute([$id]);
        $plugin = $stmt->fetch(PDO::FETCH_ASSOC);

        if (!$plugin) {
            json_response(404, '插件不存在');
            exit;
        }

        $stmt = $pdo->prepare("UPDATE plugins SET status = 'active', updated_at = NOW() WHERE id = ?");
        $stmt->execute([$id]);

        // 加载插件的hooks
        loadSinglePluginHooks($pdo, $plugins_dir, $plugin['domain'], $plugin['slug']);

        json_response(0, '启用成功');
        exit;
    }

    // ===== 禁用插件 =====
    if (preg_match('#^plugins/(\d+)/disable$#', $sub_path, $m) && $method === 'POST') {
        $id = $m[1];
        $stmt = $pdo->prepare("SELECT * FROM plugins WHERE id = ?");
        $stmt->execute([$id]);
        $plugin = $stmt->fetch(PDO::FETCH_ASSOC);

        if (!$plugin) {
            json_response(404, '插件不存在');
            exit;
        }

        $stmt = $pdo->prepare("UPDATE plugins SET status = 'inactive', updated_at = NOW() WHERE id = ?");
        $stmt->execute([$id]);

        json_response(0, '禁用成功');
        exit;
    }

    // ===== 插件配置Schema =====
    if (preg_match('#^plugins/(\d+)/schema$#', $sub_path, $m) && $method === 'GET') {
        $id = $m[1];
        $stmt = $pdo->prepare("SELECT * FROM plugins WHERE id = ?");
        $stmt->execute([$id]);
        $plugin = $stmt->fetch(PDO::FETCH_ASSOC);

        if (!$plugin) {
            json_response(404, '插件不存在');
            exit;
        }

        $manifest_path = $plugins_dir . '/' . $plugin['domain'] . '/' . $plugin['slug'] . '/manifest.json';
        $config_schema = null;
        if (file_exists($manifest_path)) {
            $manifest = json_decode(file_get_contents($manifest_path), true);
            $config_schema = $manifest['config_schema'] ?? null;
        }

        json_response(0, 'success', [
            'plugin_id' => (int)$plugin['id'],
            'name' => $plugin['name'],
            'slug' => $plugin['slug'],
            'domain' => $plugin['domain'],
            'config_schema' => $config_schema,
        ]);
        exit;
    }

    // ===== 邮件发送 =====
    if ($sub_path === 'mail/send' && $method === 'POST') {
        $body = get_json_body();
        $to = $body['to'] ?? '';
        $subject = $body['subject'] ?? '';
        $content = $body['content'] ?? '';

        if (empty($to) || empty($subject)) {
            json_response(400, 'to and subject are required');
            exit;
        }

        // 查找已启用的邮件插件
        $stmt = $pdo->prepare("SELECT * FROM plugins WHERE domain = 'mail' AND status = 'active' LIMIT 1");
        $stmt->execute();
        $mail_plugin = $stmt->fetch(PDO::FETCH_ASSOC);

        if (!$mail_plugin) {
            json_response(500, '没有可用的邮件插件');
            exit;
        }

        // TODO: 调用邮件插件的发送方法

        json_response(0, '邮件发送成功');
        exit;
    }

    // ===== 短信发送 =====
    if ($sub_path === 'sms/send' && $method === 'POST') {
        $body = get_json_body();
        $phone = $body['phone'] ?? '';
        $content = $body['content'] ?? '';

        if (empty($phone) || empty($content)) {
            json_response(400, 'phone and content are required');
            exit;
        }

        // 查找已启用的短信插件
        $stmt = $pdo->prepare("SELECT * FROM plugins WHERE domain = 'sms' AND status = 'active' LIMIT 1");
        $stmt->execute();
        $sms_plugin = $stmt->fetch(PDO::FETCH_ASSOC);

        if (!$sms_plugin) {
            json_response(500, '没有可用的短信插件');
            exit;
        }

        // TODO: 调用短信插件的发送方法

        json_response(0, '短信发送成功');
        exit;
    }

    json_response(404, 'API not found');
    exit;
}

// 客户端area（servers模块前台 - iframe嵌入）
if (strpos($path, '/clientarea/') === 0) {
    $module = substr($path, 12);
    $service_id = $_GET['service_id'] ?? '';

    $plugin_path = $plugins_dir . '/servers/' . $module;
    if (!is_dir($plugin_path)) {
        http_response_code(404);
        echo 'Module not found';
        exit;
    }

    // 尝试渲染模板
    $template_file = $plugin_path . '/templates/clientarea.html';
    if (file_exists($template_file)) {
        header('Content-Type: text/html; charset=utf-8');
        echo file_get_contents($template_file);
        exit;
    }

    http_response_code(404);
    echo 'Module template not found';
    exit;
}

// 支付回调
if (strpos($path, '/payment/') === 0) {
    $parts = explode('/', trim($path, '/'));
    $gateway = $parts[1] ?? '';

    $plugin_path = $plugins_dir . '/payment/' . $gateway;
    if (!is_dir($plugin_path)) {
        http_response_code(404);
        echo 'Gateway not found';
        exit;
    }

    $callback_file = $plugin_path . '/callback.php';
    if (file_exists($callback_file)) {
        require_once $callback_file;
        exit;
    }

    http_response_code(404);
    echo 'Callback not found';
    exit;
}

// 默认404
json_response(404, 'API not found');
