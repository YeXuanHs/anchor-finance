<?php
namespace app\admin\lib;

/**
 * zjmf插件基类兼容（让zjmf插件零修改运行）
 * 参考 zjmfv376/app/admin/lib/Plugin.php
 */
abstract class Plugin
{
    public $info = [];
    public $hasAdmin = 0;
    private $name = "";
    private $pluginPath = "";

    public function __construct()
    {
        $this->name = $this->getName();
    }

    /**
     * 安装插件
     */
    public function install()
    {
        return true;
    }

    /**
     * 卸载插件
     */
    public function uninstall()
    {
        return true;
    }

    /**
     * 获取插件配置
     */
    public function config()
    {
        return [];
    }

    /**
     * 获取插件名称
     */
    public function getName()
    {
        if ($this->name) {
            return $this->name;
        }
        $class = get_class($this);
        $parts = explode("\\", $class);
        return end($parts);
    }

    /**
     * 获取插件路径
     */
    public function getPluginPath()
    {
        return $this->pluginPath;
    }

    /**
     * 获取插件配置文件内容
     */
    public function getConfig()
    {
        $configFile = $this->pluginPath . "config.php";
        if (file_exists($configFile)) {
            return include $configFile;
        }
        return [];
    }

    /**
     * 保存插件配置
     */
    public function saveConfig($config)
    {
        $configFile = $this->pluginPath . "config.php";
        $content = "<?php\nreturn " . var_export($config, true) . ";\n";
        return file_put_contents($configFile, $content) !== false;
    }
}
