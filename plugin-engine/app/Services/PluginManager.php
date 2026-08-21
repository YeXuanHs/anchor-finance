<?php

namespace App\Services;

use App\ZjmfCompat\DatabaseCompat;
use App\ZjmfCompat\HookCompat;

/**
 * 插件管理器
 * 负责插件的扫描、加载、安装、卸载等操作
 */
class PluginManager
{
    // 插件目录
    private $pluginDir;

    // 已加载的插件
    private $loadedPlugins = [];

    /**
     * 构造函数
     */
    public function __construct()
    {
        $this->pluginDir = base_path('plugins');
    }

    /**
     * 扫描插件目录
     * @return array 插件列表
     */
    public function scan()
    {
        $plugins = [];
        $domains = ['gateways', 'sms', 'mail', 'oauth', 'servers', 'captcha', 'certification', 'addons'];

        foreach ($domains as $domain) {
            $domainDir = $this->pluginDir . '/' . $domain;
            if (!is_dir($domainDir)) {
                continue;
            }

            $items = scandir($domainDir);
            foreach ($items as $item) {
                if ($item === '.' || $item === '..') {
                    continue;
                }

                $pluginDir = $domainDir . '/' . $item;
                if (!is_dir($pluginDir)) {
                    continue;
                }

                $manifestFile = $pluginDir . '/manifest.json';
                if (file_exists($manifestFile)) {
                    $manifest = json_decode(file_get_contents($manifestFile), true);
                    if ($manifest) {
                        $manifest['domain'] = $domain;
                        $manifest['directory'] = $pluginDir;
                        $plugins[] = $manifest;
                    }
                }
            }
        }

        return $plugins;
    }

    /**
     * 安装插件
     * @param string $domain 插件域
     * @param string $name 插件名称
     * @return bool
     */
    public function install($domain, $name)
    {
        $pluginDir = $this->pluginDir . '/' . $domain . '/' . $name;
        $manifestFile = $pluginDir . '/manifest.json';

        if (!file_exists($manifestFile)) {
            return false;
        }

        $manifest = json_decode(file_get_contents($manifestFile), true);
        if (!$manifest) {
            return false;
        }

        // 检查是否已安装
        $exists = DatabaseCompat::table('plugin')
            ->where('name', $name)
            ->where('domain', $domain)
            ->exists();

        if ($exists) {
            return false;
        }

        // 读取默认配置
        $configFile = $pluginDir . '/config.php';
        $config = [];
        if (file_exists($configFile)) {
            $config = include $configFile;
        }

        // 保存到数据库
        DatabaseCompat::table('plugin')->insert([
            'name' => $name,
            'title' => $manifest['name'] ?? $name,
            'domain' => $domain,
            'description' => $manifest['description'] ?? '',
            'version' => $manifest['version'] ?? '1.0.0',
            'author' => $manifest['author'] ?? '',
            'status' => 'inactive',
            'config' => json_encode($config, JSON_UNESCAPED_UNICODE),
            'created_at' => now(),
            'updated_at' => now(),
        ]);

        // 执行安装钩子
        $this->callPluginMethod($domain, $name, 'install');

        return true;
    }

    /**
     * 卸载插件
     * @param string $domain 插件域
     * @param string $name 插件名称
     * @return bool
     */
    public function uninstall($domain, $name)
    {
        // 执行卸载钩子
        $this->callPluginMethod($domain, $name, 'uninstall');

        // 从数据库删除
        DatabaseCompat::table('plugin')
            ->where('name', $name)
            ->where('domain', $domain)
            ->delete();

        // 删除Hook绑定
        DatabaseCompat::table('hook_plugin')
            ->where('plugin_key', $domain . '/' . $name)
            ->delete();

        return true;
    }

    /**
     * 启用插件
     * @param string $domain 插件域
     * @param string $name 插件名称
     * @return bool
     */
    public function enable($domain, $name)
    {
        DatabaseCompat::table('plugin')
            ->where('name', $name)
            ->where('domain', $domain)
            ->update(['status' => 'active', 'updated_at' => now()]);

        // 加载插件的Hook
        $this->loadPluginHooks($domain, $name);

        return true;
    }

    /**
     * 禁用插件
     * @param string $domain 插件域
     * @param string $name 插件名称
     * @return bool
     */
    public function disable($domain, $name)
    {
        DatabaseCompat::table('plugin')
            ->where('name', $name)
            ->where('domain', $domain)
            ->update(['status' => 'inactive', 'updated_at' => now()]);

        return true;
    }

    /**
     * 获取插件配置
     * @param string $domain 插件域
     * @param string $name 插件名称
     * @return array
     */
    public function getConfig($domain, $name)
    {
        $plugin = DatabaseCompat::table('plugin')
            ->where('name', $name)
            ->where('domain', $domain)
            ->first();

        if (!$plugin) {
            return [];
        }

        return json_decode($plugin->config ?? '{}', true);
    }

    /**
     * 更新插件配置
     * @param string $domain 插件域
     * @param string $name 插件名称
     * @param array $config 配置
     * @return bool
     */
    public function updateConfig($domain, $name, $config)
    {
        DatabaseCompat::table('plugin')
            ->where('name', $name)
            ->where('domain', $domain)
            ->update([
                'config' => json_encode($config, JSON_UNESCAPED_UNICODE),
                'updated_at' => now(),
            ]);

        return true;
    }

    /**
     * 加载插件的Hook
     * @param string $domain 插件域
     * @param string $name 插件名称
     */
    private function loadPluginHooks($domain, $name)
    {
        $hooksFile = $this->pluginDir . '/' . $domain . '/' . $name . '/hooks.php';

        if (file_exists($hooksFile)) {
            $hooks = include $hooksFile;
            if (is_array($hooks)) {
                foreach ($hooks as $tag => $config) {
                    if (isset($config['class']) && class_exists($config['class'])) {
                        $handler = [new $config['class'], 'handle'];
                        HookCompat::add($tag, $handler, $config['priority'] ?? 10);
                    }
                }
            }
        }
    }

    /**
     * 调用插件方法
     * @param string $domain 插件域
     * @param string $name 插件名称
     * @param string $method 方法名
     * @param array $params 参数
     * @return mixed
     */
    private function callPluginMethod($domain, $name, $method, $params = [])
    {
        $pluginClass = $this->getPluginClass($domain, $name);

        if ($pluginClass && method_exists($pluginClass, $method)) {
            $instance = new $pluginClass();
            return call_user_func_array([$instance, $method], $params);
        }

        return null;
    }

    /**
     * 获取插件类名
     * @param string $domain 插件域
     * @param string $name 插件名称
     * @return string|null
     */
    private function getPluginClass($domain, $name)
    {
        $manifestFile = $this->pluginDir . '/' . $domain . '/' . $name . '/manifest.json';

        if (!file_exists($manifestFile)) {
            return null;
        }

        $manifest = json_decode(file_get_contents($manifestFile), true);
        return $manifest['entry_class'] ?? null;
    }

    /**
     * 健康检查
     * @param string $domain 插件域
     * @param string $name 插件名称
     * @return array
     */
    public function healthCheck($domain, $name)
    {
        $pluginDir = $this->pluginDir . '/' . $domain . '/' . $name;
        $manifestFile = $pluginDir . '/manifest.json';

        $status = 'healthy';
        $issues = [];

        // 检查manifest.json
        if (!file_exists($manifestFile)) {
            $status = 'error';
            $issues[] = 'manifest.json not found';
        }

        // 检查入口类
        $manifest = json_decode(file_get_contents($manifestFile), true);
        if (isset($manifest['entry_class'])) {
            if (!class_exists($manifest['entry_class'])) {
                $status = 'error';
                $issues[] = 'Entry class not found: ' . $manifest['entry_class'];
            }
        }

        return [
            'status' => $status,
            'issues' => $issues,
            'directory' => $pluginDir,
        ];
    }
}
