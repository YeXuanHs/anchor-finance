<?php

namespace App\Services\Plugin;

class HookRegistry
{
    private static array $hooks = [];
    private static array $pluginHooks = [];

    /**
     * Register a hook handler
     */
    public static function register(string $hook, string $class, int $priority = 10): void
    {
        self::$hooks[$hook][] = [
            'class' => $class,
            'priority' => $priority,
        ];
        // Sort by priority
        usort(self::$hooks[$hook], fn($a, $b) => $a['priority'] - $b['priority']);
    }

    /**
     * Load hooks from a plugin's hooks.php
     */
    public static function loadPluginHooks(string $pluginDir, string $pluginKey): void
    {
        $hooksFile = $pluginDir . '/hooks.php';
        if (!file_exists($hooksFile)) {
            return;
        }

        $hooks = include $hooksFile;
        if (!is_array($hooks)) {
            return;
        }

        foreach ($hooks as $hook => $config) {
            $class = $config['class'] ?? $config;
            $priority = $config['priority'] ?? 10;
            self::register($hook, $class, $priority);
            self::$pluginHooks[$pluginKey][] = $hook;
        }
    }

    /**
     * Trigger a hook and collect results
     */
    public static function trigger(string $hook, array $params = []): array
    {
        $results = [];
        $handlers = self::$hooks[$hook] ?? [];

        foreach ($handlers as $handler) {
            $class = $handler['class'];
            if (!class_exists($class)) {
                continue;
            }

            try {
                $instance = new $class();
                if (method_exists($instance, 'handle')) {
                    $result = $instance->handle($params);
                    $results[] = [
                        'handler' => $class,
                        'result' => $result,
                        'status' => 'success',
                    ];
                }
            } catch (\Throwable $e) {
                $results[] = [
                    'handler' => $class,
                    'error' => $e->getMessage(),
                    'status' => 'failed',
                ];
            }
        }

        return $results;
    }

    /**
     * Get all registered hooks
     */
    public static function getHooks(): array
    {
        return self::$hooks;
    }

    /**
     * Get hooks for a specific plugin
     */
    public static function getPluginHooks(string $pluginKey): array
    {
        return self::$pluginHooks[$pluginKey] ?? [];
    }

    /**
     * Clear all hooks (for testing)
     */
    public static function clear(): void
    {
        self::$hooks = [];
        self::$pluginHooks = [];
    }
}
