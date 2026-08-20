<?php

namespace App\Services\Plugin;

use RecursiveDirectoryIterator;
use RecursiveIteratorIterator;
use SplFileInfo;

class PluginScanner
{
    private string $pluginsDir;

    public function __construct(?string $pluginsDir = null)
    {
        $this->pluginsDir = $pluginsDir ?? dirname(__DIR__, 3) . '/plugins';
    }

    /**
     * Scan all plugin directories and return manifests
     */
    public function scan(string $domain = ''): array
    {
        $plugins = [];
        $domains = $domain ? [$domain] : $this->getDomains();

        foreach ($domains as $d) {
            $domainDir = $this->pluginsDir . '/' . $d;
            if (!is_dir($domainDir)) {
                continue;
            }

            $dirs = array_filter(glob($domainDir . '/*'), 'is_dir');
            foreach ($dirs as $dir) {
                $manifest = $this->loadManifest($dir, $d);
                if ($manifest) {
                    $plugins[] = $manifest;
                }
            }
        }

        return $plugins;
    }

    /**
     * Get available domains
     */
    public function getDomains(): array
    {
        $domains = [];
        $items = scandir($this->pluginsDir);
        foreach ($items as $item) {
            if ($item === '.' || $item === '..') continue;
            if (is_dir($this->pluginsDir . '/' . $item)) {
                $domains[] = $item;
            }
        }
        return $domains;
    }

    /**
     * Load manifest.json from a plugin directory
     */
    public function loadManifest(string $pluginDir, string $domain): ?array
    {
        $manifestFile = $pluginDir . '/manifest.json';
        if (!file_exists($manifestFile)) {
            // Try to generate manifest from config.php
            return $this->generateManifestFromConfig($pluginDir, $domain);
        }

        $manifest = json_decode(file_get_contents($manifestFile), true);
        if (!$manifest) {
            return null;
        }

        $manifest['_path'] = $pluginDir;
        $manifest['_domain'] = $domain;
        $manifest['_slug'] = basename($pluginDir);
        $manifest['_hash'] = md5_file($manifestFile);

        return $manifest;
    }

    /**
     * Generate manifest from zjmf-style config.php
     */
    private function generateManifestFromConfig(string $pluginDir, string $domain): ?array
    {
        $configFile = $pluginDir . '/config.php';
        if (!file_exists($configFile)) {
            return null;
        }

        $config = include $configFile;
        if (!is_array($config)) {
            return null;
        }

        $slug = basename($pluginDir);
        return [
            'name' => $config['title'] ?? $config['name'] ?? $slug,
            'slug' => $slug,
            'version' => $config['version'] ?? '1.0.0',
            'domain' => $domain,
            'description' => $config['description'] ?? '',
            'provider_class' => $config['provider_class'] ?? "Plugins\\{$domain}\\{$slug}\\Plugin",
            'entry_class' => $config['entry_class'] ?? $config['provider_class'] ?? null,
            'config_schema' => $this->configToSchema($config['config'] ?? []),
            '_path' => $pluginDir,
            '_domain' => $domain,
            '_slug' => $slug,
        ];
    }

    /**
     * Convert zjmf config array to schema format
     */
    private function configToSchema(array $config): array
    {
        $schema = [];
        foreach ($config as $key => $item) {
            if (is_array($item) && isset($item['type'])) {
                $schema[] = [
                    'key' => $key,
                    'label' => $item['title'] ?? $key,
                    'type' => $this->mapConfigType($item['type']),
                    'default' => $item['value'] ?? null,
                    'description' => $item['description'] ?? '',
                    'required' => ($item['required'] ?? false) ? true : false,
                ];
            }
        }
        return $schema;
    }

    private function mapConfigType(string $type): string
    {
        return match($type) {
            'text', 'string' => 'text',
            'textarea', 'longtext' => 'textarea',
            'select', 'dropdown' => 'select',
            'switch', 'boolean', 'bool' => 'switch',
            'number', 'int', 'integer' => 'number',
            'password', 'secret' => 'password',
            default => 'text',
        };
    }
}
