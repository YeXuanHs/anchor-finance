<?php
/**
 * AnchorFinance Plugin Engine
 * PHP插件引擎入口文件
 * 兼容zjmf插件，通过HTTP与Go后端通信
 */

use Illuminate\Http\Request;

define('LARAVEL_START', microtime(true));

// 确定当前路径
if ($maintenance = __DIR__.'/../storage/framework/maintenance.php') {
    require $maintenance;
}

// 注册自动加载器
require __DIR__.'/../vendor/autoload.php';

// 引入zjmf兼容层函数
require_once __DIR__.'/../zjmf_compat/functions.php';

// 引导Laravel
$app = require_once __DIR__.'/../bootstrap/app.php';

$kernel = $app->make(Illuminate\Contracts\Http\Kernel::class);

$response = $kernel->handle(
    $request = Request::capture()
);

$response->send();

$kernel->terminate($request, $response);
