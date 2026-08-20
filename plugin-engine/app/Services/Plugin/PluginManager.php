<?php

namespace App\Services\Plugin;

use PDO;

class PluginManager
{
    private PDO $db;
    private PluginScanner $scanner;
    private PluginFileLoader $loader;

    public function __construct(PDO $db, PluginScanner $scanner, PluginFileLoader $loader)
    {
        $this->db = $db;
        $this->scanner = $scanner;
        $this->loader = $loader;
    }

    /**
     * List plugins with status from database
     */
    public function list(string $domain = ''): array
    {
        $scanned = $this->scanner->scan($domain);
        $installed = $this->getInstalledPlugins();

        $result = [];
        foreach ($scanned as $manifest) {
            $slug = $manifest['_slug'];
            $d = $manifest['_domain'];
            $key = "{$d}:{$slug}";

            $dbRecord = $installed[$key] ?? null;

            $result[] = [
                'id' => $dbRecord['id'] ?? null,
                'domain' => $d,
                'slug' => $slug,
                'plugin_key' => $manifest['plugin_key'] ?? $key,
                'name' => $manifest['name'] ?? $slug,
                'version' => $manifest['version'] ?? '1.0.0',
                'description' => $manifest['description'] ?? '',
                'provider_class' => $manifest['provider_class'] ?? null,
                'entry_class' => $manifest['entry_class'] ?? null,
                'config_schema' => $manifest['config_schema'] ?? [],
                'capabilities' => $manifest['capabilities'] ?? [],
                'is_installed' => $dbRecord !== null,
                'is_enabled' => ($dbRecord['status'] ?? 0) == 2,
                'status' => $dbRecord['status'] ?? 0,
                'installed_at' => $dbRecord['installed_at'] ?? null,
                'config' => $dbRecord ? $this->decryptConfig($dbRecord['config_json'] ?? '{}') : [],
                'manifest_hash' => $manifest['_hash'] ?? null,
                'path' => $manifest['_path'],
            ];
        }

        return $result;
    }

    /**
     * Install a plugin
     */
    public function install(string $domain, string $slug): array
    {
        $plugins = $this->list($domain);
        $plugin = null;
        foreach ($plugins as $p) {
            if ($p['slug'] === $slug && $p['domain'] === $domain) {
                $plugin = $p;
                break;
            }
        }

        if (!$plugin) {
            throw new \RuntimeException("Plugin not found: {$domain}/{$slug}");
        }

        if ($plugin['is_installed']) {
            return $plugin;
        }

        $pluginKey = $plugin['plugin_key'];
        $now = date('Y-m-d H:i:s');

        $stmt = $this->db->prepare(
            "INSERT INTO integration_plugins (domain, slug, plugin_key, name, version, provider_class, entry_class, capabilities_json, config_schema_json, status, installed_at, manifest_hash, created_at, updated_at)
             VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, 1, ?, ?, ?, ?)"
        );
        $stmt->execute([
            $domain,
            $slug,
            $pluginKey,
            $plugin['name'],
            $plugin['version'],
            $plugin['provider_class'],
            $plugin['entry_class'],
            json_encode($plugin['capabilities']),
            json_encode($plugin['config_schema']),
            $now,
            $plugin['manifest_hash'],
            $now,
            $now,
        ]);

        $id = $this->db->lastInsertId();

        // Create default config
        $defaultConfig = [];
        foreach ($plugin['config_schema'] ?? [] as $field) {
            if (isset($field['default'])) {
                $defaultConfig[$field['key']] = $field['default'];
            }
        }
        $configStmt = $this->db->prepare(
            "INSERT INTO integration_plugin_configs (plugin_id, config_json, created_at, updated_at) VALUES (?, ?, ?, ?)"
        );
        $configStmt->execute([$id, json_encode($defaultConfig), $now, $now]);

        // Run install hook
        $this->loader->load($plugin['path']);
        HookRegistry::loadPluginHooks($plugin['path'], $pluginKey);
        HookRegistry::trigger('plugin_installed', ['plugin' => $plugin]);

        $plugin['id'] = $id;
        $plugin['is_installed'] = true;
        $plugin['status'] = 1;
        $plugin['installed_at'] = $now;
        return $plugin;
    }

    /**
     * Enable a plugin
     */
    public function enable(int $id): void
    {
        $plugin = $this->getById($id);
        if (!$plugin) {
            throw new \RuntimeException("Plugin not found: {$id}");
        }

        $now = date('Y-m-d H:i:s');
        $stmt = $this->db->prepare("UPDATE integration_plugins SET status = 2, updated_at = ? WHERE id = ?");
        $stmt->execute([$now, $id]);

        // Load plugin files and hooks
        $this->loader->load($plugin['_path']);
        HookRegistry::loadPluginHooks($plugin['_path'], $plugin['plugin_key']);
        HookRegistry::trigger('plugin_enabled', ['plugin_id' => $id]);
    }

    /**
     * Disable a plugin
     */
    public function disable(int $id): void
    {
        $now = date('Y-m-d H:i:s');
        $stmt = $this->db->prepare("UPDATE integration_plugins SET status = 1, updated_at = ? WHERE id = ?");
        $stmt->execute([$now, $id]);

        HookRegistry::trigger('plugin_disabled', ['plugin_id' => $id]);
    }

    /**
     * Uninstall a plugin
     */
    public function uninstall(int $id, bool $force = false): void
    {
        $plugin = $this->getById($id);
        if (!$plugin) {
            throw new \RuntimeException("Plugin not found: {$id}");
        }

        // Delete config
        $this->db->prepare("DELETE FROM integration_plugin_configs WHERE plugin_id = ?")->execute([$id]);
        // Delete plugin record
        $this->db->prepare("DELETE FROM integration_plugins WHERE id = ?")->execute([$id]);

        HookRegistry::trigger('plugin_uninstalled', ['plugin_id' => $id, 'force' => $force]);
    }

    /**
     * Update plugin config
     */
    public function updateConfig(int $id, array $config): void
    {
        $now = date('Y-m-d H:i:s');
        $stmt = $this->db->prepare(
            "INSERT INTO integration_plugin_configs (plugin_id, config_json, created_at, updated_at)
             VALUES (?, ?, ?, ?)
             ON DUPLICATE KEY UPDATE config_json = VALUES(config_json), updated_at = VALUES(updated_at)"
        );
        $stmt->execute([$id, json_encode($config), $now, $now]);

        HookRegistry::trigger('plugin_config_updated', ['plugin_id' => $id, 'config' => $config]);
    }

    /**
     * Get plugin config
     */
    public function getConfig(int $id): array
    {
        $stmt = $this->db->prepare("SELECT config_json FROM integration_plugin_configs WHERE plugin_id = ?");
        $stmt->execute([$id]);
        $row = $stmt->fetch(PDO::FETCH_ASSOC);
        return $row ? json_decode($row['config_json'], true) ?? [] : [];
    }

    /**
     * Get enabled plugins for a domain
     */
    public function getEnabledForDomain(string $domain): array
    {
        $stmt = $this->db->prepare(
            "SELECT p.*, c.config_json FROM integration_plugins p
             LEFT JOIN integration_plugin_configs c ON c.plugin_id = p.id
             WHERE p.domain = ? AND p.status = 2 AND p.deleted_at IS NULL"
        );
        $stmt->execute([$domain]);
        return $stmt->fetchAll(PDO::FETCH_ASSOC);
    }

    /**
     * Get plugin config by domain (for Go backend to call)
     */
    public function getActiveConfigForDomain(string $domain): ?array
    {
        $plugins = $this->getEnabledForDomain($domain);
        if (empty($plugins)) {
            return null;
        }

        $plugin = $plugins[0];
        return [
            'plugin_key' => $plugin['plugin_key'],
            'name' => $plugin['name'],
            'config' => $this->decryptConfig($plugin['config_json'] ?? '{}'),
        ];
    }

    private function getById(int $id): ?array
    {
        $stmt = $this->db->prepare("SELECT * FROM integration_plugins WHERE id = ?");
        $stmt->execute([$id]);
        $row = $stmt->fetch(PDO::FETCH_ASSOC);
        if (!$row) return null;

        // Get the path
        $plugins = $this->list($row['domain']);
        foreach ($plugins as $p) {
            if ($p['slug'] === $row['slug']) {
                return array_merge($p, $row);
            }
        }
        return null;
    }

    private function getInstalledPlugins(): array
    {
        $stmt = $this->db->query("SELECT * FROM integration_plugins WHERE deleted_at IS NULL");
        $rows = $stmt->fetchAll(PDO::FETCH_ASSOC);

        $result = [];
        foreach ($rows as $row) {
            $key = "{$row['domain']}:{$row['slug']}";
            $result[$key] = $row;
        }
        return $result;
    }

    private function decryptConfig(string $json): array
    {
        return json_decode($json, true) ?? [];
    }
}
