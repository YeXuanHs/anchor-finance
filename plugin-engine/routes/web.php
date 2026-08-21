<?php

use Illuminate\Support\Facades\Route;
use App\Http\Controllers\HookController;
use App\Http\Controllers\PluginController;

/*
|--------------------------------------------------------------------------
| Web Routes
|--------------------------------------------------------------------------
|
| 插件引擎路由
| 所有请求都通过Go后端代理，不直接暴露给前端
|
*/

// 健康检查
Route::get('/health', function () {
    return response()->json([
        'status' => 'ok',
        'service' => 'AnchorFinance Plugin Engine',
    ]);
});

// 内部API（Go后端调用）
Route::prefix('/internal')->group(function () {
    // Hook相关
    Route::post('/hook/trigger', [HookController::class, 'trigger']);
    Route::get('/hook/list', [HookController::class, 'list']);
    Route::post('/hook/register', [HookController::class, 'register']);

    // 插件管理
    Route::get('/plugins', [PluginController::class, 'index']);
    Route::post('/plugins/install', [PluginController::class, 'install']);
    Route::post('/plugins/uninstall/{id}', [PluginController::class, 'uninstall']);
    Route::post('/plugins/enable/{id}', [PluginController::class, 'enable']);
    Route::post('/plugins/disable/{id}', [PluginController::class, 'disable']);
    Route::get('/plugins/{id}/config', [PluginController::class, 'getConfig']);
    Route::put('/plugins/{id}/config', [PluginController::class, 'updateConfig']);
    Route::post('/plugins/{id}/health', [PluginController::class, 'healthCheck']);
    Route::post('/plugins/scan', [PluginController::class, 'scan']);
});

// 客户端area路由（servers模块前台界面）
Route::prefix('/clientarea')->group(function () {
    Route::get('/{module}', function ($module) {
        // 加载servers模块的前台模板
        $serviceId = request('service_id');
        $pluginPath = base_path('plugins/servers/' . $module);

        if (!is_dir($pluginPath)) {
            return response('Module not found', 404);
        }

        // 加载模块的控制器
        $controllerFile = $pluginPath . '/' . ucfirst($module) . 'Controller.php';
        if (file_exists($controllerFile)) {
            require_once $controllerFile;
            $controllerClass = 'Plugins\\Servers\\' . ucfirst($module) . '\\' . ucfirst($module) . 'Controller';

            if (class_exists($controllerClass)) {
                $controller = new $controllerClass();
                if (method_exists($controller, 'clientArea')) {
                    return $controller->clientArea($serviceId);
                }
            }
        }

        // 如果没有控制器，尝试渲染模板
        $templateFile = $pluginPath . '/templates/clientarea.html';
        if (file_exists($templateFile)) {
            return response(file_get_contents($templateFile));
        }

        return response('Module template not found', 404);
    });
});

// 支付回调路由
Route::prefix('/payment')->group(function () {
    Route::get('/{gateway}/callback', function ($gateway) {
        $pluginPath = base_path('plugins/gateways/' . $gateway);

        if (!is_dir($pluginPath)) {
            return response('Gateway not found', 404);
        }

        // 加载网关的回调处理
        $callbackFile = $pluginPath . '/callback.php';
        if (file_exists($callbackFile)) {
            require_once $callbackFile;
        }

        return response('OK');
    });

    Route::post('/{gateway}/notify', function ($gateway) {
        $pluginPath = base_path('plugins/gateways/' . $gateway);

        if (!is_dir($pluginPath)) {
            return response('Gateway not found', 404);
        }

        // 加载网关的通知处理
        $notifyFile = $pluginPath . '/notify.php';
        if (file_exists($notifyFile)) {
            require_once $notifyFile;
        }

        return response('OK');
    });
});
