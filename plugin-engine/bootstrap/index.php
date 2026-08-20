<?php
/**
 * AnchorFinance Plugin Engine - HTTP API Entry Point
 *
 * Called by Go backend via internal HTTP requests.
 * All requests are JSON.
 */

require_once __DIR__ . '/../vendor/autoload.php';

use App\Http\Controllers\Internal\PluginController;
use App\Services\Plugin\PluginFileLoader;
use App\Services\Plugin\PluginManager;
use App\Services\Plugin\PluginScanner;

// Load .env if exists
$envFile = __DIR__ . '/../.env';
if (file_exists($envFile)) {
    $dotenv = Dotenv\Dotenv::createImmutable(dirname($envFile));
    $dotenv->safeLoad();
}

// Database connection
$dbHost = $_ENV['DB_HOST'] ?? 'localhost';
$dbPort = $_ENV['DB_PORT'] ?? '3306';
$dbName = $_ENV['DB_NAME'] ?? 'anchorfinance';
$dbUser = $_ENV['DB_USER'] ?? 'root';
$dbPass = $_ENV['DB_PASS'] ?? '';

try {
    $pdo = new PDO(
        "mysql:host={$dbHost};port={$dbPort};dbname={$dbName};charset=utf8mb4",
        $dbUser,
        $dbPass,
        [PDO::ATTR_ERRMODE => PDO::ERRMODE_EXCEPTION]
    );
} catch (PDOException $e) {
    http_response_code(500);
    echo json_encode(['error' => 'Database connection failed: ' . $e->getMessage()]);
    exit;
}

// Initialize services
$scanner = new PluginScanner();
$loader = new PluginFileLoader();
$manager = new PluginManager($pdo, $scanner, $loader);
$controller = new PluginController($manager, $scanner);

// Parse request
$method = $_SERVER['REQUEST_METHOD'];
$uri = $_SERVER['REQUEST_URI'];
$path = parse_url($uri, PHP_URL_PATH);
$path = preg_replace('#^/api/plugin-engine#', '', $path);
$path = rtrim($path, '/');

// Parse JSON body
$body = [];
if (in_array($method, ['POST', 'PUT', 'PATCH'])) {
    $input = file_get_contents('php://input');
    $body = json_decode($input, true) ?? [];
}

// Route matching
header('Content-Type: application/json');

try {
    $result = null;

    // GET /domains
    if ($method === 'GET' && $path === '/domains') {
        $result = $controller->domains();
    }
    // GET /plugins
    elseif ($method === 'GET' && $path === '/plugins') {
        $result = $controller->list($_GET);
    }
    // POST /plugins/install
    elseif ($method === 'POST' && $path === '/plugins/install') {
        $result = $controller->install($body);
    }
    // POST /plugins/scan
    elseif ($method === 'POST' && $path === '/plugins/scan') {
        $result = $controller->scan($body);
    }
    // GET /plugins/config/domain/{domain}
    elseif ($method === 'GET' && preg_match('#^/plugins/config/domain/(.+)$#', $path, $m)) {
        $result = $controller->getActiveConfig(['domain' => $m[1]]);
    }
    // POST /plugins/{id}/enable
    elseif ($method === 'POST' && preg_match('#^/plugins/(\d+)/enable$#', $path, $m)) {
        $result = $controller->enable(['id' => $m[1]]);
    }
    // POST /plugins/{id}/disable
    elseif ($method === 'POST' && preg_match('#^/plugins/(\d+)/disable$#', $path, $m)) {
        $result = $controller->disable(['id' => $m[1]]);
    }
    // DELETE /plugins/{id}
    elseif ($method === 'DELETE' && preg_match('#^/plugins/(\d+)$#', $path, $m)) {
        $force = $_GET['force'] ?? 'false';
        $result = $controller->uninstall(['id' => $m[1], 'force' => $force]);
    }
    // GET /plugins/{id}/config
    elseif ($method === 'GET' && preg_match('#^/plugins/(\d+)/config$#', $path, $m)) {
        $result = $controller->getConfig(['id' => $m[1]]);
    }
    // PUT /plugins/{id}/config
    elseif ($method === 'PUT' && preg_match('#^/plugins/(\d+)/config$#', $path, $m)) {
        $result = $controller->updateConfig(['id' => $m[1]], $body);
    }
    // POST /hooks/trigger
    elseif ($method === 'POST' && $path === '/hooks/trigger') {
        $result = $controller->triggerHook($body);
    }
    // GET /health
    elseif ($method === 'GET' && $path === '/health') {
        $result = ['status' => 'ok', 'version' => '1.0.0'];
    }
    else {
        http_response_code(404);
        $result = ['error' => 'Not found', 'path' => $path];
    }

    // Handle error responses
    if (isset($result['error']) && isset($result['code'])) {
        http_response_code($result['code']);
    }

    echo json_encode($result, JSON_UNESCAPED_UNICODE);

} catch (Throwable $e) {
    http_response_code(500);
    echo json_encode([
        'error' => $e->getMessage(),
        'file' => $e->getFile(),
        'line' => $e->getLine(),
    ]);
}
