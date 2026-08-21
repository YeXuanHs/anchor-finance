<?php

namespace App\Http\Controllers;

use App\Services\PluginManager;
use Illuminate\Http\Request;
use Illuminate\Http\JsonResponse;

/**
 * 插件控制器
 * 处理插件的CRUD操作
 */
class PluginController extends Controller
{
    private $pluginManager;

    public function __construct(PluginManager $pluginManager)
    {
        $this->pluginManager = $pluginManager;
    }

    /**
     * 获取插件列表
     * @param Request $request
     * @return JsonResponse
     */
    public function index(Request $request): JsonResponse
    {
        $domain = $request->input('domain');

        // 从数据库获取已安装的插件
        $query = \App\ZjmfCompat\DatabaseCompat::table('plugin');

        if ($domain) {
            $query->where('domain', $domain);
        }

        $plugins = $query->get();

        // 扫描文件系统中的插件
        $scannedPlugins = $this->pluginManager->scan();

        // 合并数据
        $result = [];
        foreach ($scannedPlugins as $plugin) {
            $installed = $plugins->firstWhere('name', $plugin['slug'] ?? '');
            $result[] = [
                'id' => $installed->id ?? null,
                'name' => $plugin['name'] ?? '',
                'slug' => $plugin['slug'] ?? '',
                'domain' => $plugin['domain'] ?? '',
                'description' => $plugin['description'] ?? '',
                'version' => $plugin['version'] ?? '',
                'author' => $plugin['author'] ?? '',
                'status' => $installed->status ?? 'not_installed',
                'config_schema' => $plugin['config_schema'] ?? [],
            ];
        }

        return response()->json([
            'code' => 0,
            'message' => 'success',
            'data' => [
                'list' => $result,
                'total' => count($result),
            ],
        ]);
    }

    /**
     * 安装插件
     * @param Request $request
     * @return JsonResponse
     */
    public function install(Request $request): JsonResponse
    {
        $domain = $request->input('domain');
        $name = $request->input('name');

        if (empty($domain) || empty($name)) {
            return response()->json([
                'code' => 400,
                'message' => 'Domain and name are required',
            ], 400);
        }

        $result = $this->pluginManager->install($domain, $name);

        if ($result) {
            return response()->json([
                'code' => 0,
                'message' => 'Plugin installed successfully',
            ]);
        }

        return response()->json([
            'code' => 500,
            'message' => 'Failed to install plugin',
        ], 500);
    }

    /**
     * 卸载插件
     * @param int $id
     * @return JsonResponse
     */
    public function uninstall(int $id): JsonResponse
    {
        $plugin = \App\ZjmfCompat\DatabaseCompat::table('plugin')->where('id', $id)->first();

        if (!$plugin) {
            return response()->json([
                'code' => 404,
                'message' => 'Plugin not found',
            ], 404);
        }

        $result = $this->pluginManager->uninstall($plugin->domain, $plugin->name);

        if ($result) {
            return response()->json([
                'code' => 0,
                'message' => 'Plugin uninstalled successfully',
            ]);
        }

        return response()->json([
            'code' => 500,
            'message' => 'Failed to uninstall plugin',
        ], 500);
    }

    /**
     * 启用插件
     * @param int $id
     * @return JsonResponse
     */
    public function enable(int $id): JsonResponse
    {
        $plugin = \App\ZjmfCompat\DatabaseCompat::table('plugin')->where('id', $id)->first();

        if (!$plugin) {
            return response()->json([
                'code' => 404,
                'message' => 'Plugin not found',
            ], 404);
        }

        $result = $this->pluginManager->enable($plugin->domain, $plugin->name);

        if ($result) {
            return response()->json([
                'code' => 0,
                'message' => 'Plugin enabled successfully',
            ]);
        }

        return response()->json([
            'code' => 500,
            'message' => 'Failed to enable plugin',
        ], 500);
    }

    /**
     * 禁用插件
     * @param int $id
     * @return JsonResponse
     */
    public function disable(int $id): JsonResponse
    {
        $plugin = \App\ZjmfCompat\DatabaseCompat::table('plugin')->where('id', $id)->first();

        if (!$plugin) {
            return response()->json([
                'code' => 404,
                'message' => 'Plugin not found',
            ], 404);
        }

        $result = $this->pluginManager->disable($plugin->domain, $plugin->name);

        if ($result) {
            return response()->json([
                'code' => 0,
                'message' => 'Plugin disabled successfully',
            ]);
        }

        return response()->json([
            'code' => 500,
            'message' => 'Failed to disable plugin',
        ], 500);
    }

    /**
     * 获取插件配置
     * @param int $id
     * @return JsonResponse
     */
    public function getConfig(int $id): JsonResponse
    {
        $plugin = \App\ZjmfCompat\DatabaseCompat::table('plugin')->where('id', $id)->first();

        if (!$plugin) {
            return response()->json([
                'code' => 404,
                'message' => 'Plugin not found',
            ], 404);
        }

        $config = $this->pluginManager->getConfig($plugin->domain, $plugin->name);

        return response()->json([
            'code' => 0,
            'message' => 'success',
            'data' => $config,
        ]);
    }

    /**
     * 更新插件配置
     * @param int $id
     * @param Request $request
     * @return JsonResponse
     */
    public function updateConfig(int $id, Request $request): JsonResponse
    {
        $plugin = \App\ZjmfCompat\DatabaseCompat::table('plugin')->where('id', $id)->first();

        if (!$plugin) {
            return response()->json([
                'code' => 404,
                'message' => 'Plugin not found',
            ], 404);
        }

        $config = $request->input('config', []);
        $result = $this->pluginManager->updateConfig($plugin->domain, $plugin->name, $config);

        if ($result) {
            return response()->json([
                'code' => 0,
                'message' => 'Plugin config updated successfully',
            ]);
        }

        return response()->json([
            'code' => 500,
            'message' => 'Failed to update plugin config',
        ], 500);
    }

    /**
     * 健康检查
     * @param int $id
     * @return JsonResponse
     */
    public function healthCheck(int $id): JsonResponse
    {
        $plugin = \App\ZjmfCompat\DatabaseCompat::table('plugin')->where('id', $id)->first();

        if (!$plugin) {
            return response()->json([
                'code' => 404,
                'message' => 'Plugin not found',
            ], 404);
        }

        $result = $this->pluginManager->healthCheck($plugin->domain, $plugin->name);

        return response()->json([
            'code' => 0,
            'message' => 'success',
            'data' => $result,
        ]);
    }

    /**
     * 扫描插件
     * @return JsonResponse
     */
    public function scan(): JsonResponse
    {
        $plugins = $this->pluginManager->scan();

        return response()->json([
            'code' => 0,
            'message' => 'success',
            'data' => [
                'list' => $plugins,
                'total' => count($plugins),
            ],
        ]);
    }
}
