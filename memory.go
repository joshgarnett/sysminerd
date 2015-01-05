package main

const MemoryModuleName = "memory"

type MemoryInputModule struct{}

func (m *MemoryInputModule) Name() string {
	return MemoryModuleName
}

func (m *MemoryInputModule) Init(config *Config, moduleConfig *ModuleConfig) error {
	return nil
}

func (m *MemoryInputModule) TearDown() error {
	return nil
}
