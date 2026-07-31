package server

import (
	"fmt"
	"sync"
)

// registry 模块注册表
var (
	registry   = make(map[string]ModuleConstructor)
	registryMu sync.RWMutex
)

// RegisterModule 注册服务器模块
func RegisterModule(name string, constructor ModuleConstructor) {
	registryMu.Lock()
	defer registryMu.Unlock()
	registry[name] = constructor
}

// CreateModule 创建服务器模块实例
func CreateModule(name string, config map[string]interface{}) (ServerModule, error) {
	registryMu.RLock()
	constructor, ok := registry[name]
	registryMu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("server module %q not registered", name)
	}

	return constructor(config)
}

// ListModules 返回所有已注册的模块名称
func ListModules() []string {
	registryMu.RLock()
	defer registryMu.RUnlock()

	names := make([]string, 0, len(registry))
	for name := range registry {
		names = append(names, name)
	}
	return names
}

// GetModuleInfo 获取模块信息
func GetModuleInfo(name string) (map[string]interface{}, error) {
	registryMu.RLock()
	constructor, ok := registry[name]
	registryMu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("server module %q not registered", name)
	}

	// 创建临时实例获取信息
	module, err := constructor(nil)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"name":    module.Name(),
		"title":   module.Title(),
		"options": module.GetConfigOptions(),
	}, nil
}
