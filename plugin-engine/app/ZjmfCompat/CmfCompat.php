<?php

namespace App\ZjmfCompat;

use Illuminate\Support\Facades\Request;
use Illuminate\Support\Facades\Cache;
use Illuminate\Support\Facades\Log;

/**
 * CMF兼容层
 * 兼容zjmf的cmf_xxx()系列函数
 */
class CmfCompat
{
    // 配置缓存
    private static $configCache = [];

    // 当前用户信息缓存
    private static $currentUser = null;
    private static $currentAdmin = null;

    /**
     * 获取系统配置
     * @param string $key 配置键名（可选）
     * @return mixed
     */
    public static function configuration($key = '')
    {
        if (empty($key)) {
            return self::$configCache;
        }

        // 先从缓存获取
        if (isset(self::$configCache[$key])) {
            return self::$configCache[$key];
        }

        // 从数据库获取
        try {
            $config = DatabaseCompat::table('configuration')
                ->where('key', $key)
                ->first();

            if ($config) {
                $value = $config->value;
                // 尝试JSON解码
                $decoded = json_decode($value, true);
                if (json_last_error() === JSON_ERROR_NONE) {
                    $value = $decoded;
                }
                self::$configCache[$key] = $value;
                return $value;
            }
        } catch (\Exception $e) {
            // 静默失败
        }

        return null;
    }

    /**
     * 更新系统配置
     * @param string $key 配置键名
     * @param mixed $value 配置值
     * @return bool
     */
    public static function updateConfiguration($key, $value)
    {
        try {
            if (is_array($value) || is_object($value)) {
                $value = json_encode($value, JSON_UNESCAPED_UNICODE);
            }

            $exists = DatabaseCompat::table('configuration')
                ->where('key', $key)
                ->exists();

            if ($exists) {
                DatabaseCompat::table('configuration')
                    ->where('key', $key)
                    ->update(['value' => $value, 'updated_at' => now()]);
            } else {
                DatabaseCompat::table('configuration')->insert([
                    'key' => $key,
                    'value' => $value,
                    'created_at' => now(),
                    'updated_at' => now(),
                ]);
            }

            // 更新缓存
            self::$configCache[$key] = is_string($value) ? json_decode($value, true) ?? $value : $value;

            return true;
        } catch (\Exception $e) {
            return false;
        }
    }

    /**
     * 加密数据（AES-256-CBC）
     * @param string $data 要加密的数据
     * @param string $key 加密密钥
     * @return string 加密后的数据（Base64编码）
     */
    public static function encrypt($data, $key)
    {
        $key = substr(hash('sha256', $key), 0, 32);
        $iv = openssl_random_pseudo_bytes(16);
        $encrypted = openssl_encrypt($data, 'AES-256-CBC', $key, 0, $iv);
        return base64_encode($iv . '::' . $encrypted);
    }

    /**
     * 解密数据
     * @param string $data 要解密的数据（Base64编码）
     * @param string $key 解密密钥
     * @return string 解密后的数据
     */
    public static function decrypt($data, $key)
    {
        $key = substr(hash('sha256', $key), 0, 32);
        $data = base64_decode($data);
        list($iv, $encrypted) = explode('::', $data, 2);
        return openssl_decrypt($encrypted, 'AES-256-CBC', $key, 0, $iv);
    }

    /**
     * 生成插件URL
     * @param string $url URL路径
     * @return string 完整URL
     */
    public static function addonUrl($url)
    {
        $baseUrl = config('app.url', '');
        return rtrim($baseUrl, '/') . '/plugin-engine/' . ltrim($url, '/');
    }

    /**
     * 获取后台地址
     * @return string
     */
    public static function adminAddress()
    {
        return config('app.url', '') . '/admin';
    }

    /**
     * 生成站点URL
     * @param string $path 路径
     * @return string
     */
    public static function siteUrl($path = '')
    {
        $baseUrl = config('app.url', '');
        return rtrim($baseUrl, '/') . '/' . ltrim($path, '/');
    }

    /**
     * 获取当前管理员ID
     * @return int|null
     */
    public static function getCurrentAdminId()
    {
        $admin = self::getCurrentAdmin();
        return $admin ? $admin['id'] : null;
    }

    /**
     * 获取当前用户信息
     * @return array|null
     */
    public static function getCurrentUser()
    {
        if (self::$currentUser !== null) {
            return self::$currentUser;
        }

        try {
            $userId = request()->header('X-User-ID');
            if ($userId) {
                $user = DatabaseCompat::table('clients')
                    ->where('id', $userId)
                    ->first();
                if ($user) {
                    self::$currentUser = (array) $user;
                    return self::$currentUser;
                }
            }
        } catch (\Exception $e) {
            // 静默失败
        }

        return null;
    }

    /**
     * 获取当前管理员信息
     * @return array|null
     */
    public static function getCurrentAdmin()
    {
        if (self::$currentAdmin !== null) {
            return self::$currentAdmin;
        }

        try {
            $adminId = request()->header('X-Admin-ID');
            if ($adminId) {
                $admin = DatabaseCompat::table('admins')
                    ->where('id', $adminId)
                    ->first();
                if ($admin) {
                    self::$currentAdmin = (array) $admin;
                    return self::$currentAdmin;
                }
            }
        } catch (\Exception $e) {
            // 静默失败
        }

        return null;
    }

    /**
     * 获取插件类名
     * @param string $name 插件名称
     * @param string $type 插件类型
     * @return string|null 类名
     */
    public static function getPluginClass($name, $type)
    {
        $namespace = 'Plugins\\' . ucfirst($type) . '\\' . ucfirst($name);
        $className = $namespace . '\\' . ucfirst($name) . 'Plugin';

        if (class_exists($className)) {
            return $className;
        }

        return null;
    }

    /**
     * 名称解析
     * @param string $name 名称
     * @param int $type 类型（0:下划线转驼峰, 1:驼峰转下划线）
     * @return string
     */
    public static function parseName($name, $type = 0)
    {
        if ($type === 0) {
            // 下划线转驼峰
            return str_replace(' ', '', ucwords(str_replace('_', ' ', $name)));
        } else {
            // 驼峰转下划线
            return strtolower(preg_replace('/(?<!^)[A-Z]/', '_$0', $name));
        }
    }

    /**
     * 获取语言文本
     * @param string $key 语言键名
     * @param array $params 替换参数
     * @return string
     */
    public static function lang($key, $params = [])
    {
        // 简单实现：返回键名
        // 后续可以加载语言文件
        $text = $key;

        if (!empty($params)) {
            foreach ($params as $k => $v) {
                $text = str_replace('{' . $k . '}', $v, $text);
            }
        }

        return $text;
    }

    /**
     * 用户是否已登录
     * @return bool
     */
    public static function isLogin()
    {
        return self::getCurrentUser() !== null;
    }

    /**
     * 管理员是否已登录
     * @return bool
     */
    public static function isAdminLogin()
    {
        return self::getCurrentAdmin() !== null;
    }

    /**
     * 获取客户端IP
     * @return string
     */
    public static function getClientIp()
    {
        return Request::ip();
    }

    /**
     * 获取用户代理
     * @return string
     */
    public static function getUserAgent()
    {
        return Request::userAgent();
    }

    /**
     * 生成随机字符串
     * @param int $length 长度
     * @return string
     */
    public static function randomString($length = 16)
    {
        return bin2hex(random_bytes($length / 2));
    }

    /**
     * 发送HTTP请求
     * @param string $url URL
     * @param array $options 选项
     * @return array
     */
    public static function httpRequest($url, $options = [])
    {
        try {
            $client = new \GuzzleHttp\Client();
            $method = $options['method'] ?? 'GET';
            $response = $client->request($method, $url, $options);

            return [
                'status' => $response->getStatusCode(),
                'body' => $response->getBody()->getContents(),
                'headers' => $response->getHeaders(),
            ];
        } catch (\Exception $e) {
            return [
                'status' => 0,
                'body' => '',
                'error' => $e->getMessage(),
            ];
        }
    }

    /**
     * 记录日志
     * @param string $type 日志类型
     * @param string $content 日志内容
     * @param int $userid 用户ID
     */
    public static function log($type, $content, $userid = 0)
    {
        try {
            DatabaseCompat::table('system_logs')->insert([
                'type' => $type,
                'content' => $content,
                'user_id' => $userid,
                'ip' => self::getClientIp(),
                'created_at' => now(),
            ]);
        } catch (\Exception $e) {
            // 静默失败
        }
    }

    /**
     * 获取插件配置
     * @param string $pluginName 插件名称
     * @param string $key 配置键名
     * @return mixed
     */
    public static function getPluginConfig($pluginName, $key = '')
    {
        try {
            $plugin = DatabaseCompat::table('plugin')
                ->where('name', $pluginName)
                ->first();

            if (!$plugin) {
                return null;
            }

            $config = json_decode($plugin->config ?? '{}', true);

            if (empty($key)) {
                return $config;
            }

            return $config[$key] ?? null;
        } catch (\Exception $e) {
            return null;
        }
    }

    /**
     * 保存插件配置
     * @param string $pluginName 插件名称
     * @param array $config 配置
     * @return bool
     */
    public static function savePluginConfig($pluginName, $config)
    {
        try {
            $configJson = json_encode($config, JSON_UNESCAPED_UNICODE);

            $exists = DatabaseCompat::table('plugin')
                ->where('name', $pluginName)
                ->exists();

            if ($exists) {
                DatabaseCompat::table('plugin')
                    ->where('name', $pluginName)
                    ->update(['config' => $configJson, 'updated_at' => now()]);
            } else {
                DatabaseCompat::table('plugin')->insert([
                    'name' => $pluginName,
                    'config' => $configJson,
                    'status' => 'active',
                    'created_at' => now(),
                    'updated_at' => now(),
                ]);
            }

            return true;
        } catch (\Exception $e) {
            return false;
        }
    }

    /**
     * 认证后取消暂停
     * @param int $uid 用户ID
     */
    public static function unsuspendAfterCertify($uid)
    {
        try {
            DatabaseCompat::table('host')
                ->where('userid', $uid)
                ->where('status', 'Suspended')
                ->update(['status' => 'Active']);
        } catch (\Exception $e) {
            // 静默失败
        }
    }
}
