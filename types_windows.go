// +build windows

package main

import (
	"syscall"
)

type (
	dword     uint32
	dwordLong uint64

	memorystatusex struct {
		dwLength                dword
		dwMemoryLoad            dword
		ullTotalPhys            dwordLong
		ullAvailPhys            dwordLong
		ullTotalPageFile        dwordLong
		ullAvailPageFile        dwordLong
		ullTotalVirtual         dwordLong
		ullAvailVirtual         dwordLong
		ullAvailExtendedVirtual dwordLong
	}
)

var kernel32 = syscall.NewLazyDLL("kernel32.dll")

var (
	_globalMemoryStatusEx = kernel32.NewProc("GlobalMemoryStatusEx")
	_getDiskFreeSpaceEx   = kernel32.NewProc("GetDiskFreeSpaceExW")
	_getLogicalDrives     = kernel32.NewProc("GetLogicalDrives")
)
