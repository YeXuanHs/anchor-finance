<?php
/**
 * PHP内置服务器路由文件
 * 所有请求都转发到index.php
 */

$uri = $_SERVER['REQUEST_URI'];
$path = parse_url($uri, PHP_URL_PATH);

// 静态文件直接返回
if ($path !== '/index.php' && file_exists(__DIR__ . $path)) {
    return false;
}

// 所有其他请求转发到index.php
require_once __DIR__ . '/index.php';
