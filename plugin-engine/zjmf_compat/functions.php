<?php
/**
 * zjmf全局函数兼容层
 * 让zjmf插件可以零修改运行在AnchorFinance插件引擎中
 */

use App\ZjmfCompat\DatabaseCompat;
use App\ZjmfCompat\HookCompat;
use App\ZjmfCompat\CmfCompat;
use App\ZjmfCompat\TemplateEngine;

// ============================================================================
// 数据库查询
// ============================================================================

/**
 * 兼容 zjmf 的 db('table') 写法
 * @param string $table 表名
 * @return \Illuminate\Database\Query\Builder
 */
function db($table) {
    return DatabaseCompat::table($table);
}

/**
 * 兼容 ThinkPHP 的 Db::name() 写法
 */
class Db {
    public static function name($table) {
        return DatabaseCompat::table($table);
    }

    public static function table($table) {
        return DatabaseCompat::table($table);
    }

    public static function query($sql, $bindings = []) {
        return DatabaseCompat::query($sql, $bindings);
    }

    public static function execute($sql, $bindings = []) {
        return DatabaseCompat::execute($sql, $bindings);
    }
}

// ============================================================================
// Hook注册和触发
// ============================================================================

/**
 * 注册Hook
 * @param string $tag Hook名称
 * @param callable $fun 处理函数
 */
function hook_add($tag, $fun) {
    HookCompat::add($tag, $fun);
}

/**
 * 触发Hook（返回所有结果）
 * @param string $hook Hook名称
 * @param mixed $params 参数
 * @return array 所有handler的返回结果
 */
function hook($hook, $params = null) {
    return HookCompat::trigger($hook, $params);
}

/**
 * 触发单个Hook（返回第一个结果）
 * @param string $hook Hook名称
 * @param mixed $params 参数
 * @return mixed 第一个handler的返回结果
 */
function hook_one($hook, $params = null) {
    return HookCompat::triggerOne($hook, $params);
}

/**
 * 获取系统Hook列表
 * @return array
 */
function getSystemHook() {
    return HookCompat::getSystemHooks();
}

// ============================================================================
// 系统配置
// ============================================================================

/**
 * 获取系统配置
 * @param string $key 配置键名（可选）
 * @return mixed
 */
function configuration($key = '') {
    return CmfCompat::configuration($key);
}

/**
 * 更新系统配置
 * @param string $key 配置键名
 * @param mixed $value 配置值
 * @return bool
 */
function updateConfiguration($key, $value) {
    return CmfCompat::updateConfiguration($key, $value);
}

// ============================================================================
// 加密解密
// ============================================================================

/**
 * 加密数据
 * @param string $data 要加密的数据
 * @param string $key 加密密钥
 * @return string 加密后的数据
 */
function cmf_encrypt($data, $key) {
    return CmfCompat::encrypt($data, $key);
}

/**
 * 解密数据
 * @param string $data 要解密的数据
 * @param string $key 解密密钥
 * @return string 解密后的数据
 */
function cmf_decrypt($data, $key) {
    return CmfCompat::decrypt($data, $key);
}

// ============================================================================
// URL生成
// ============================================================================

/**
 * 生成插件URL
 * @param string $url URL路径
 * @return string 完整URL
 */
function shd_addon_url($url) {
    return CmfCompat::addonUrl($url);
}

/**
 * 获取后台地址
 * @return string
 */
function adminAddress() {
    return CmfCompat::adminAddress();
}

/**
 * 生成站点URL
 * @param string $path 路径
 * @return string
 */
function cmf_url($path = '') {
    return CmfCompat::siteUrl($path);
}

// ============================================================================
// 用户信息
// ============================================================================

/**
 * 获取当前管理员ID
 * @return int|null
 */
function cmf_get_current_admin_id() {
    return CmfCompat::getCurrentAdminId();
}

/**
 * 获取当前用户信息
 * @return array|null
 */
function cmf_get_current_user() {
    return CmfCompat::getCurrentUser();
}

/**
 * 获取当前管理员信息
 * @return array|null
 */
function cmf_get_current_admin() {
    return CmfCompat::getCurrentAdmin();
}

// ============================================================================
// 插件相关
// ============================================================================

/**
 * 获取插件类
 * @param string $name 插件名称
 * @param string $type 插件类型
 * @return string|null 类名
 */
function cmf_get_plugin_class_shd($name, $type) {
    return CmfCompat::getPluginClass($name, $type);
}

/**
 * 名称解析
 * @param string $name 名称
 * @param int $type 类型（0:下划线转驼峰, 1:驼峰转下划线）
 * @return string
 */
function cmf_parse_name($name, $type = 0) {
    return CmfCompat::parseName($name, $type);
}

// ============================================================================
// 模板渲染
// ============================================================================

/**
 * 渲染模板
 * @param string $tpl 模板路径
 * @param array $data 模板数据
 * @return string 渲染后的HTML
 */
function template($tpl, $data = []) {
    return TemplateEngine::render($tpl, $data);
}

// ============================================================================
// 语言
// ============================================================================

/**
 * 获取语言文本
 * @param string $key 语言键名
 * @param array $params 替换参数
 * @return string
 */
function lang($key, $params = []) {
    return CmfCompat::lang($key, $params);
}

// ============================================================================
// 系统状态
// ============================================================================

/**
 * 系统是否已安装
 * @return bool
 */
function cmf_is_installed() {
    return true; // AnchorFinance始终返回true
}

/**
 * 用户是否已登录
 * @return bool
 */
function cmf_is_login() {
    return CmfCompat::isLogin();
}

/**
 * 管理员是否已登录
 * @return bool
 */
function cmf_is_admin_login() {
    return CmfCompat::isAdminLogin();
}

// ============================================================================
// 其他工具函数
// ============================================================================

/**
 * 获取客户端IP
 * @return string
 */
function cmf_get_ip() {
    return CmfCompat::getClientIp();
}

/**
 * 获取用户代理
 * @return string
 */
function cmf_get_user_agent() {
    return CmfCompat::getUserAgent();
}

/**
 * 生成随机字符串
 * @param int $length 长度
 * @return string
 */
function cmf_random_string($length = 16) {
    return CmfCompat::randomString($length);
}

/**
 * 发送HTTP请求
 * @param string $url URL
 * @param array $options 选项
 * @return array
 */
function cmf_http_request($url, $options = []) {
    return CmfCompat::httpRequest($url, $options);
}

/**
 * 记录日志
 * @param string $type 日志类型
 * @param string $content 日志内容
 * @param int $userid 用户ID
 */
function cmf_log($type, $content, $userid = 0) {
    CmfCompat::log($type, $content, $userid);
}

/**
 * 获取插件配置
 * @param string $pluginName 插件名称
 * @param string $key 配置键名
 * @return mixed
 */
function get_plugin_config($pluginName, $key = '') {
    return CmfCompat::getPluginConfig($pluginName, $key);
}

/**
 * 保存插件配置
 * @param string $pluginName 插件名称
 * @param array $config 配置
 * @return bool
 */
function save_plugin_config($pluginName, $config) {
    return CmfCompat::savePluginConfig($pluginName, $config);
}

/**
 * 实名认证检查后的取消暂停
 * @param int $uid 用户ID
 */
function unsuspendAfterCertify($uid) {
    // AnchorFinance兼容实现：认证后取消暂停
    // 实际逻辑在Go后端实现
    CmfCompat::unsuspendAfterCertify($uid);
}
