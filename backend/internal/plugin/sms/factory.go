package sms

import (
	"fmt"
	"sync"
)

// PluginConstructor is a function that creates an SmsSender from a config map.
type PluginConstructor func(config map[string]interface{}) (SmsSender, error)

// registry holds all registered SMS plugin constructors.
var (
	registry   = make(map[string]PluginConstructor)
	registryMu sync.RWMutex
)

// RegisterPlugin registers an SMS plugin constructor under the given name.
// This is typically called from init() in each plugin file.
func RegisterPlugin(name string, constructor PluginConstructor) {
	registryMu.Lock()
	defer registryMu.Unlock()
	registry[name] = constructor
}

// CreatePlugin creates an SMS sender instance by provider name and config.
func CreatePlugin(provider string, config map[string]interface{}) (SmsSender, error) {
	registryMu.RLock()
	constructor, ok := registry[provider]
	registryMu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("sms plugin %q not registered", provider)
	}

	return constructor(config)
}

// ListPlugins returns the names of all registered SMS plugins.
func ListPlugins() []string {
	registryMu.RLock()
	defer registryMu.RUnlock()

	names := make([]string, 0, len(registry))
	for name := range registry {
		names = append(names, name)
	}
	return names
}
