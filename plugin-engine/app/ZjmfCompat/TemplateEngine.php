<?php

namespace App\ZjmfCompat;

/**
 * 模板引擎兼容层
 * 支持zjmf的.tpl模板渲染
 */
class TemplateEngine
{
    // 模板目录
    private static $templateDir = '';

    // 编译目录
    private static $compileDir = '';

    // 模板变量
    private static $variables = [];

    /**
     * 设置模板目录
     * @param string $dir 目录路径
     */
    public static function setTemplateDir($dir)
    {
        self::$templateDir = rtrim($dir, '/');
    }

    /**
     * 设置编译目录
     * @param string $dir 目录路径
     */
    public static function setCompileDir($dir)
    {
        self::$compileDir = rtrim($dir, '/');
        if (!is_dir(self::$compileDir)) {
            mkdir(self::$compileDir, 0755, true);
        }
    }

    /**
     * 渲染模板
     * @param string $tplPath 模板路径
     * @param array $data 模板数据
     * @return string 渲染后的HTML
     */
    public static function render($tplPath, $data = [])
    {
        // 如果是相对路径，拼接模板目录
        if (!str_starts_with($tplPath, '/') && !str_starts_with($tplPath, ':')) {
            $tplPath = self::$templateDir . '/' . $tplPath;
        }

        // 检查模板文件是否存在
        if (!file_exists($tplPath)) {
            return "<!-- Template not found: {$tplPath} -->";
        }

        // 合并变量
        $variables = array_merge(self::$variables, $data);

        // 读取模板内容
        $content = file_get_contents($tplPath);

        // 简单的模板解析（兼容zjmf的.tpl格式）
        $content = self::parseTemplate($content, $variables);

        return $content;
    }

    /**
     * 解析模板内容
     * @param string $content 模板内容
     * @param array $variables 变量
     * @return string 解析后的内容
     */
    private static function parseTemplate($content, $variables)
    {
        // 解析 {$variable} 格式的变量
        $content = preg_replace_callback('/\{\$(\w+)\}/', function ($matches) use ($variables) {
            $key = $matches[1];
            return $variables[$key] ?? '';
        }, $content);

        // 解析 {if $condition} ... {else} ... {/if} 格式的条件
        $content = preg_replace_callback('/\{if\s+\$(\w+)\}\s*(.*?)\s*\{else\}\s*(.*?)\s*\{\/if\}/s', function ($matches) use ($variables) {
            $key = $matches[1];
            $trueContent = $matches[2];
            $falseContent = $matches[3];

            if (!empty($variables[$key])) {
                return $trueContent;
            }
            return $falseContent;
        }, $content);

        // 解析 {foreach $array as $item} ... {/foreach} 格式的循环
        $content = preg_replace_callback('/\{foreach\s+\$(\w+)\s+as\s+\$(\w+)\}\s*(.*?)\s*\{\/foreach\}/s', function ($matches) use ($variables) {
            $arrayKey = $matches[1];
            $itemKey = $matches[2];
            $template = $matches[3];

            $result = '';
            if (isset($variables[$arrayKey]) && is_array($variables[$arrayKey])) {
                foreach ($variables[$arrayKey] as $item) {
                    $itemContent = $template;
                    if (is_array($item)) {
                        foreach ($item as $k => $v) {
                            $itemContent = str_replace('{$' . $k . '}', $v, $itemContent);
                        }
                    }
                    $result .= $itemContent;
                }
            }

            return $result;
        }, $content);

        // 解析 {include file="xxx.tpl"} 格式的包含
        $content = preg_replace_callback('/\{include\s+file="([^"]+)"\}/', function ($matches) {
            $includePath = self::$templateDir . '/' . $matches[1];
            if (file_exists($includePath)) {
                return file_get_contents($includePath);
            }
            return "<!-- Include not found: {$matches[1]} -->";
        }, $content);

        // 解析 {$variable|default:'value'} 格式的默认值
        $content = preg_replace_callback('/\{\$(\w+)\|default:\'([^\']*)\'\}/', function ($matches) use ($variables) {
            $key = $matches[1];
            $default = $matches[2];
            return $variables[$key] ?? $default;
        }, $content);

        return $content;
    }

    /**
     * 设置模板变量
     * @param string $key 变量名
     * @param mixed $value 变量值
     */
    public static function assign($key, $value = null)
    {
        if (is_array($key)) {
            self::$variables = array_merge(self::$variables, $key);
        } else {
            self::$variables[$key] = $value;
        }
    }

    /**
     * 清除模板变量
     */
    public static function clearVariables()
    {
        self::$variables = [];
    }

    /**
     * 获取模板变量
     * @param string $key 变量名
     * @return mixed
     */
    public static function getVariable($key)
    {
        return self::$variables[$key] ?? null;
    }

    /**
     * 检查模板文件是否存在
     * @param string $tplPath 模板路径
     * @return bool
     */
    public static function exists($tplPath)
    {
        if (!str_starts_with($tplPath, '/') && !str_starts_with($tplPath, ':')) {
            $tplPath = self::$templateDir . '/' . $tplPath;
        }
        return file_exists($tplPath);
    }

    /**
     * 获取模板目录
     * @return string
     */
    public static function getTemplateDir()
    {
        return self::$templateDir;
    }

    /**
     * 获取编译目录
     * @return string
     */
    public static function getCompileDir()
    {
        return self::$compileDir;
    }
}
