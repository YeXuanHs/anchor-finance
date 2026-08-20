<?php

namespace App\Http\Controllers\Internal;

use App\Services\Plugin\HookRegistry;
use App\Services\Plugin\PluginManager;
use App\Services\Plugin\PluginScanner;

class PluginController
{
    private PluginManager $manager;
    private PluginScanner $scanner;

    public function __construct(PluginManager $manager, PluginScanner $scanner)
    {
        $this->manager = $manager;
        $this->scanner = $scanner;
    }

    /**
     * GET /plugins?domain=xxx
     */
    public function list(array $params): array
    {
        $domain = $params['domain'] ?? '';
        return ['list' => $this->manager->list($domain)];
    }

    /**
     * POST /plugins/install
     */
    public function install(array $body): array
    {
        $domain = $body['domain'] ?? '';
        $slug = $body['slug'] ?? '';
        if (!$domain || !$slug) {
            return ['error' => 'domain and slug required', 'code' => 400];
        }
        return $this->manager->install($domain, $slug);
    }

    /**
     * POST /plugins/{id}/enable
     */
    public function enable(array $params): array
    {
        $id = (int)($params['id'] ?? 0);
        $this->manager->enable($id);
        return ['message' => 'enabled'];
    }

    /**
     * POST /plugins/{id}/disable
     */
    public function disable(array $params): array
    {
        $id = (int)($params['id'] ?? 0);
        $this->manager->disable($id);
        return ['message' => 'disabled'];
    }

    /**
     * DELETE /plugins/{id}
     */
    public function uninstall(array $params): array
    {
        $id = (int)($params['id'] ?? 0);
        $force = ($params['force'] ?? 'false') === 'true';
        $this->manager->uninstall($id, $force);
        return ['message' => 'uninstalled'];
    }

    /**
     * GET /plugins/{id}/config
     */
    public function getConfig(array $params): array
    {
        $id = (int)($params['id'] ?? 0);
        return ['config' => $this->manager->getConfig($id)];
    }

    /**
     * PUT /plugins/{id}/config
     */
    public function updateConfig(array $params, array $body): array
    {
        $id = (int)($params['id'] ?? 0);
        $this->manager->updateConfig($id, $body);
        return ['message' => 'config updated'];
    }

    /**
     * POST /plugins/scan
     */
    public function scan(array $params): array
    {
        $domain = $params['domain'] ?? '';
        $plugins = $this->scanner->scan($domain);
        return ['list' => $plugins];
    }

    /**
     * POST /hooks/trigger - Trigger a hook
     */
    public function triggerHook(array $body): array
    {
        $hook = $body['hook'] ?? '';
        $hookParams = $body['params'] ?? [];

        if (!$hook) {
            return ['error' => 'hook name required', 'code' => 400];
        }

        $results = HookRegistry::trigger($hook, $hookParams);
        return ['results' => $results];
    }

    /**
     * GET /plugins/config/domain/{domain} - Get active config for a domain (used by Go backend)
     */
    public function getActiveConfig(array $params): array
    {
        $domain = $params['domain'] ?? '';
        $config = $this->manager->getActiveConfigForDomain($domain);
        return ['config' => $config];
    }

    /**
     * GET /domains - List available plugin domains
     */
    public function domains(): array
    {
        return ['domains' => $this->scanner->getDomains()];
    }
}
