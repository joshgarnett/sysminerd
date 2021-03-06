// +build linux

package main

import (
	"fmt"
	"golang.org/x/sys/unix" //see https://godoc.org/golang.org/x/sys/unix
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
)

func (m *DiskspaceInputModule) GetMetrics() (*ModuleMetrics, error) {
	metrics := make([]Metric, 0, 50)

	filesystems, err := GetFileSystems(m)
	if err != nil {
		log.Printf("Error retrieving filesystems: %v", err)
		return nil, err
	}
	stats, err := GetFilesystemStats(filesystems)
	if err != nil {
		log.Printf("Error retrieving filesystem stats: %v", err)
		return nil, err
	}

	for _, stat := range stats {
		used := NewMetric(fmt.Sprintf("%s.used", stat.DeviceName), stat.Used)
		free := NewMetric(fmt.Sprintf("%s.free", stat.DeviceName), stat.Free)
		reserved := NewMetric(fmt.Sprintf("%s.reserved", stat.DeviceName), stat.Reserved)
		available := NewMetric(fmt.Sprintf("%s.available", stat.DeviceName), stat.Available)

		metrics = append(metrics, used)
		metrics = append(metrics, free)
		metrics = append(metrics, reserved)
		metrics = append(metrics, available)
	}

	return &ModuleMetrics{Module: m.Name(), Metrics: metrics}, nil
}

func GetFileSystems(m *DiskspaceInputModule) ([]Filesystem, error) {
	b, err := ioutil.ReadFile("/proc/mounts")
	if err != nil {
		return nil, err
	}
	content := string(b)
	lines := strings.Split(content, "\n")

	filesystems := make([]Filesystem, 0, len(lines))

	for _, line := range lines {
		fields := strings.Fields(line)

		if len(fields) < 3 {
			continue
		}

		device := fields[0]
		mount := fields[1]
		fsType := fields[2]

		if !m.CheckTypes.Contains(fsType) {
			continue
		}

		mountPrefix := strings.Split(mount, "/")[1]
		if mountPrefix == "proc" || mountPrefix == "dev" || mountPrefix == "sys" {
			continue
		}

		deviceSym, err := filepath.EvalSymlinks(device)
		if err == nil {
			device = deviceSym
		}

		stat := unix.Stat_t{}
		err = unix.Stat(mount, &stat)
		if err != nil {
			// filesystem likely not mounted
			continue
		}

		fs := Filesystem{
			Device: device,
			Mount:  mount,
			Type:   fsType,
		}

		filesystems = append(filesystems, fs)
	}

	return filesystems, nil
}

func GetFilesystemStats(filesystems []Filesystem) ([]FilesystemStats, error) {
	stats := make([]FilesystemStats, 0, len(filesystems))

	for _, fs := range filesystems {
		stat := unix.Statfs_t{}
		err := unix.Statfs(fs.Mount, &stat)
		if err != nil {
			continue
		}

		// change /dev/sda1 to sda1
		var name string
		if fs.Device[:1] == "/" {
			_, name = filepath.Split(fs.Device)
		}
		if name == "" {
			name = strings.Replace(fs.Mount[1:], "/", "_", -1)
		}

		used := float64(stat.Bsize) * float64(stat.Blocks-stat.Bfree)
		free := float64(stat.Bsize) * float64(stat.Bfree)
		reserved := float64(stat.Bsize) * float64(stat.Bfree-stat.Bavail)
		available := float64(stat.Bsize) * float64(stat.Bavail)

		stats = append(stats, FilesystemStats{DeviceName: name, Used: used, Free: free, Reserved: reserved, Available: available})
	}

	return stats, nil
}
