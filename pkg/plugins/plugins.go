// Package plugins provides the plugin system for extending go-lisp functionality
package plugins

import (
	"fmt"
	"sort"
	"sync"

	"github.com/leinonen/lisp-interpreter/pkg/registry"
)

// Plugin represents a loadable plugin that can register functions
type Plugin interface {
	// Metadata
	Name() string
	Version() string
	Description() string
	Dependencies() []string

	// Lifecycle methods
	Initialize(registry registry.FunctionRegistry) error
	Shutdown() error

	// Function registration
	RegisterFunctions(registry registry.FunctionRegistry) error
}

// PluginInfo contains metadata about a loaded plugin
type PluginInfo struct {
	Name         string
	Version      string
	Description  string
	Dependencies []string
	Functions    []string
	Loaded       bool
}

// PluginManager manages plugin lifecycle
type PluginManager interface {
	LoadPlugin(plugin Plugin) error
	UnloadPlugin(name string) error
	ReloadPlugin(name string) error
	ListPlugins() []PluginInfo
	GetPlugin(name string) (Plugin, bool)
	IsLoaded(name string) bool
	GetDependencies(name string) []string
	GetDependents(name string) []string
}

// pluginManager implements PluginManager
type pluginManager struct {
	plugins         map[string]Plugin
	pluginFunctions map[string][]string // plugin name -> function names
	registry        registry.FunctionRegistry
	mutex           sync.RWMutex
}

// NewPluginManager creates a new plugin manager
func NewPluginManager(reg registry.FunctionRegistry) PluginManager {
	return &pluginManager{
		plugins:         make(map[string]Plugin),
		pluginFunctions: make(map[string][]string),
		registry:        reg,
	}
}

// LoadPlugin loads and initializes a plugin
func (pm *pluginManager) LoadPlugin(plugin Plugin) error {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	name := plugin.Name()
	if name == "" {
		return fmt.Errorf("plugin name cannot be empty")
	}

	if _, exists := pm.plugins[name]; exists {
		return fmt.Errorf("plugin %s already loaded", name)
	}

	// Check dependencies
	for _, dep := range plugin.Dependencies() {
		if !pm.isLoadedUnsafe(dep) {
			return fmt.Errorf("plugin %s requires dependency %s which is not loaded", name, dep)
		}
	}

	// Track functions before registration
	functionsBefore := pm.registry.List()

	// Initialize plugin
	if err := plugin.Initialize(pm.registry); err != nil {
		return fmt.Errorf("failed to initialize plugin %s: %v", name, err)
	}

	// Register functions
	if err := plugin.RegisterFunctions(pm.registry); err != nil {
		// Try to cleanup on failure
		plugin.Shutdown()
		return fmt.Errorf("failed to register functions for plugin %s: %v", name, err)
	}

	// Track which functions this plugin registered
	functionsAfter := pm.registry.List()
	newFunctions := difference(functionsAfter, functionsBefore)

	pm.plugins[name] = plugin
	pm.pluginFunctions[name] = newFunctions

	return nil
}

// UnloadPlugin unloads a plugin and its functions
func (pm *pluginManager) UnloadPlugin(name string) error {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	plugin, exists := pm.plugins[name]
	if !exists {
		return fmt.Errorf("plugin %s not loaded", name)
	}

	// Check if other plugins depend on this one
	dependents := pm.getDependentsUnsafe(name)
	if len(dependents) > 0 {
		return fmt.Errorf("cannot unload plugin %s: depended on by %v", name, dependents)
	}

	// Unregister functions
	if functions, exists := pm.pluginFunctions[name]; exists {
		for _, funcName := range functions {
			pm.registry.Unregister(funcName)
		}
		delete(pm.pluginFunctions, name)
	}

	// Shutdown plugin
	if err := plugin.Shutdown(); err != nil {
		// Log error but continue with unloading
		// In a real implementation, we'd use a proper logger
		fmt.Printf("Warning: error shutting down plugin %s: %v\n", name, err)
	}

	delete(pm.plugins, name)
	return nil
}

// ReloadPlugin unloads and reloads a plugin
func (pm *pluginManager) ReloadPlugin(name string) error {
	pm.mutex.RLock()
	plugin, exists := pm.plugins[name]
	pm.mutex.RUnlock()

	if !exists {
		return fmt.Errorf("plugin %s not loaded", name)
	}

	// Unload first
	if err := pm.UnloadPlugin(name); err != nil {
		return err
	}

	// Reload
	return pm.LoadPlugin(plugin)
}

// ListPlugins returns information about all plugins
func (pm *pluginManager) ListPlugins() []PluginInfo {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()

	plugins := make([]PluginInfo, 0, len(pm.plugins))
	for name, plugin := range pm.plugins {
		functions := pm.pluginFunctions[name]
		if functions == nil {
			functions = []string{}
		}

		plugins = append(plugins, PluginInfo{
			Name:         plugin.Name(),
			Version:      plugin.Version(),
			Description:  plugin.Description(),
			Dependencies: plugin.Dependencies(),
			Functions:    functions,
			Loaded:       true,
		})
	}

	// Sort by name
	sort.Slice(plugins, func(i, j int) bool {
		return plugins[i].Name < plugins[j].Name
	})

	return plugins
}

// GetPlugin retrieves a loaded plugin
func (pm *pluginManager) GetPlugin(name string) (Plugin, bool) {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()

	plugin, exists := pm.plugins[name]
	return plugin, exists
}

// IsLoaded checks if a plugin is loaded
func (pm *pluginManager) IsLoaded(name string) bool {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()

	return pm.isLoadedUnsafe(name)
}

// isLoadedUnsafe checks if a plugin is loaded (without locking)
func (pm *pluginManager) isLoadedUnsafe(name string) bool {
	_, exists := pm.plugins[name]
	return exists
}

// GetDependencies returns the dependencies of a plugin
func (pm *pluginManager) GetDependencies(name string) []string {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()

	if plugin, exists := pm.plugins[name]; exists {
		deps := plugin.Dependencies()
		if deps == nil {
			return []string{}
		}
		// Return a copy
		result := make([]string, len(deps))
		copy(result, deps)
		return result
	}
	return []string{}
}

// GetDependents returns plugins that depend on the given plugin
func (pm *pluginManager) GetDependents(name string) []string {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()

	return pm.getDependentsUnsafe(name)
}

// getDependentsUnsafe returns plugins that depend on the given plugin (without locking)
func (pm *pluginManager) getDependentsUnsafe(name string) []string {
	var dependents []string
	for _, plugin := range pm.plugins {
		for _, dep := range plugin.Dependencies() {
			if dep == name {
				dependents = append(dependents, plugin.Name())
				break
			}
		}
	}
	return dependents
}

// Helper function to find difference between two string slices
func difference(a, b []string) []string {
	mb := make(map[string]bool, len(b))
	for _, x := range b {
		mb[x] = true
	}

	var diff []string
	for _, x := range a {
		if !mb[x] {
			diff = append(diff, x)
		}
	}
	return diff
}

// BasePlugin provides a basic implementation of Plugin interface
type BasePlugin struct {
	name         string
	version      string
	description  string
	dependencies []string
}

// NewBasePlugin creates a new base plugin
func NewBasePlugin(name, version, description string, dependencies []string) *BasePlugin {
	if dependencies == nil {
		dependencies = []string{}
	}
	return &BasePlugin{
		name:         name,
		version:      version,
		description:  description,
		dependencies: dependencies,
	}
}

// Name returns the plugin name
func (bp *BasePlugin) Name() string {
	return bp.name
}

// Version returns the plugin version
func (bp *BasePlugin) Version() string {
	return bp.version
}

// Description returns the plugin description
func (bp *BasePlugin) Description() string {
	return bp.description
}

// Dependencies returns the plugin dependencies
func (bp *BasePlugin) Dependencies() []string {
	return bp.dependencies
}

// Initialize provides a default implementation (no-op)
func (bp *BasePlugin) Initialize(registry registry.FunctionRegistry) error {
	return nil
}

// Shutdown provides a default implementation (no-op)
func (bp *BasePlugin) Shutdown() error {
	return nil
}
