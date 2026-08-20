<?php

namespace App\Services\Plugin;

class PluginFileLoader
{
    private array $loaded = [];

    /**
     * Load all PHP files from a plugin directory
     */
    public function load(string $pluginDir): void
    {
        if (isset($this->loaded[$pluginDir])) {
            return;
        }

        // Load lib, src, controller directories
        foreach (['lib', 'src', 'controller', 'Controllers', 'Services'] as $dir) {
            $this->requirePhpFilesIn($pluginDir . '/' . $dir);
        }

        // Load vendor autoload if exists
        $vendorAutoload = $pluginDir . '/vendor/autoload.php';
        if (is_file($vendorAutoload)) {
            require_once $vendorAutoload;
        }

        // Load root PHP files (except config.php)
        foreach (glob($pluginDir . '/*.php') ?: [] as $file) {
            if (basename($file) === 'config.php') continue;
            require_once $file;
        }

        $this->loaded[$pluginDir] = true;
    }

    /**
     * Recursively load PHP files from a directory
     */
    private function requirePhpFilesIn(string $path): void
    {
        if (!is_dir($path)) {
            return;
        }

        $iterator = new \RecursiveIteratorIterator(
            new \RecursiveDirectoryIterator($path, \FilesystemIterator::SKIP_DOTS)
        );

        foreach ($iterator as $file) {
            if ($file->isFile() && strtolower($file->getExtension()) === 'php') {
                require_once $file->getPathname();
            }
        }
    }

    /**
     * Check if a plugin directory has been loaded
     */
    public function isLoaded(string $pluginDir): bool
    {
        return isset($this->loaded[$pluginDir]);
    }
}
