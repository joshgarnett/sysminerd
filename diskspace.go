package main

import (
	"log"
)

type Filesystem struct {
	Device string
	Mount  string
	Type   string
}

type FilesystemStats struct {
	DeviceName string
	Used       float64
	Free       float64
	Reserved   float64
	Available  float64
}

const DiskspaceModuleName = "diskspace"

var defaultCheckTypes = []string{"ext2", "ext3", "ext4", "xfs", "glusterfs", "nfs", "ntfs", "hfs", "fat32", "fat16", "btrfs"}

type DiskspaceInputModule struct {
	CheckTypes StringSet
}

func (m *DiskspaceInputModule) Name() string {
	return DiskspaceModuleName
}

func (m *DiskspaceInputModule) Init(config *Config, moduleConfig *ModuleConfig) error {
	m.CheckTypes = StringSet{}

	types, err := moduleConfig.SettingsStringArray("filesystems")
	if err != nil {
		log.Printf("filesystems not set, using default settings: %v", err)
		types = defaultCheckTypes
	}

	log.Printf("types: %v", types)
	m.CheckTypes.AddAll(types)

	return nil
}

func (m *DiskspaceInputModule) TearDown() error {
	return nil
}
