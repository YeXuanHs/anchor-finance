<?php

namespace App\ZjmfCompat;

/**
 * Hook兼容层
 * 兼容zjmf的Hook注册和触发机制
 */
class HookCompat
{
    // 已注册的Hook处理器
    private static $hooks = [];

    // 系统内置Hook列表
    private static $systemHooks = [];

    // Hook执行结果缓存
    private static $results = [];

    /**
     * 注册Hook处理器
     * @param string $tag Hook名称
     * @param callable $handler 处理函数
     * @param int $priority 优先级（越小越先执行）
     */
    public static function add($tag, $handler, $priority = 10)
    {
        if (!isset(self::$hooks[$tag])) {
            self::$hooks[$tag] = [];
        }

        self::$hooks[$tag][] = [
            'handler' => $handler,
            'priority' => $priority,
        ];

        // 按优先级排序
        usort(self::$hooks[$tag], function ($a, $b) {
            return $a['priority'] - $b['priority'];
        });
    }

    /**
     * 触发Hook（返回所有结果）
     * @param string $tag Hook名称
     * @param mixed $params 参数
     * @return array 所有handler的返回结果
     */
    public static function trigger($tag, $params = null)
    {
        $results = [];

        if (!isset(self::$hooks[$tag])) {
            return $results;
        }

        foreach (self::$hooks[$tag] as $hook) {
            try {
                $result = call_user_func($hook['handler'], $params);
                $results[] = [
                    'handler' => $hook['handler'],
                    'result' => $result,
                    'status' => 'success',
                ];
            } catch (\Exception $e) {
                $results[] = [
                    'handler' => $hook['handler'],
                    'result' => null,
                    'status' => 'error',
                    'error' => $e->getMessage(),
                ];
            }
        }

        // 缓存结果
        self::$results[$tag] = $results;

        return $results;
    }

    /**
     * 触发单个Hook（返回第一个结果）
     * @param string $tag Hook名称
     * @param mixed $params 参数
     * @return mixed 第一个handler的返回结果
     */
    public static function triggerOne($tag, $params = null)
    {
        if (!isset(self::$hooks[$tag]) || empty(self::$hooks[$tag])) {
            return null;
        }

        $hook = self::$hooks[$tag][0];

        try {
            return call_user_func($hook['handler'], $params);
        } catch (\Exception $e) {
            return null;
        }
    }

    /**
     * 获取系统Hook列表
     * @return array
     */
    public static function getSystemHooks()
    {
        return self::$systemHooks;
    }

    /**
     * 设置系统Hook列表
     * @param array $hooks Hook列表
     */
    public static function setSystemHooks($hooks)
    {
        self::$systemHooks = $hooks;
    }

    /**
     * 获取指定Hook的所有处理器
     * @param string $tag Hook名称
     * @return array
     */
    public static function getHandlers($tag)
    {
        return self::$hooks[$tag] ?? [];
    }

    /**
     * 获取所有已注册的Hook
     * @return array
     */
    public static function getAllHooks()
    {
        return self::$hooks;
    }

    /**
     * 检查Hook是否有处理器
     * @param string $tag Hook名称
     * @return bool
     */
    public static function hasHandlers($tag)
    {
        return !empty(self::$hooks[$tag]);
    }

    /**
     * 移除Hook处理器
     * @param string $tag Hook名称
     * @param callable $handler 处理函数（可选，不传则移除所有）
     */
    public static function remove($tag, $handler = null)
    {
        if (!isset(self::$hooks[$tag])) {
            return;
        }

        if ($handler === null) {
            unset(self::$hooks[$tag]);
            return;
        }

        self::$hooks[$tag] = array_filter(self::$hooks[$tag], function ($hook) use ($handler) {
            return $hook['handler'] !== $handler;
        });

        if (empty(self::$hooks[$tag])) {
            unset(self::$hooks[$tag]);
        }
    }

    /**
     * 获取Hook执行结果
     * @param string $tag Hook名称
     * @return array|null
     */
    public static function getResults($tag)
    {
        return self::$results[$tag] ?? null;
    }

    /**
     * 清除结果缓存
     */
    public static function clearResults()
    {
        self::$results = [];
    }

    /**
     * 从数据库加载Hook绑定
     * 用于从hook_bindings表加载已注册的Hook
     */
    public static function loadFromDatabase()
    {
        try {
            $bindings = DatabaseCompat::table('hook_plugin')
                ->where('status', 'active')
                ->get();

            foreach ($bindings as $binding) {
                // 加载插件的hooks.php文件
                $pluginPath = base_path('plugins/' . $binding->plugin_key . '/hooks.php');
                if (file_exists($pluginPath)) {
                    $hooks = include $pluginPath;
                    if (is_array($hooks)) {
                        foreach ($hooks as $tag => $config) {
                            if (isset($config['class']) && class_exists($config['class'])) {
                                $handler = [new $config['class'], 'handle'];
                                self::add($tag, $handler, $config['priority'] ?? 10);
                            }
                        }
                    }
                }
            }
        } catch (\Exception $e) {
            // 静默失败，不影响系统运行
        }
    }
}
