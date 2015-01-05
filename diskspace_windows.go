// +build windows

package main

import (
	"fmt"
	"strconv"
	"syscall"
	"time"
	"unsafe"
)

func (m *DiskspaceInputModule) GetMetrics() ([]Metric, error) {

	metrics := make([]Metric, 0, 50)
	now := time.Now()

	// we grab all the drives
	drives := m.GetLogicalDrives()

	for _, drive := range drives {

		lpFreeBytesAvailable := int64(0)
		lpTotalNumberOfBytes := int64(0)
		lpTotalNumberOfFreeBytes := int64(0)

		_, _, _ = _getDiskFreeSpaceEx.Call(uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(drive))),
			uintptr(unsafe.Pointer(&lpFreeBytesAvailable)),
			uintptr(unsafe.Pointer(&lpTotalNumberOfBytes)),
			uintptr(unsafe.Pointer(&lpTotalNumberOfFreeBytes)))

		usedBytes := lpTotalNumberOfBytes - lpTotalNumberOfFreeBytes

		used := Metric{
			module:    m.Name(),
			name:      fmt.Sprintf("%s.used", drive[:1]),
			value:     float64(usedBytes),
			timestamp: now,
		}
		free := Metric{
			module:    m.Name(),
			name:      fmt.Sprintf("%s.free", drive[:1]),
			value:     float64(lpTotalNumberOfFreeBytes),
			timestamp: now,
		}
		available := Metric{
			module:    m.Name(),
			name:      fmt.Sprintf("%s.available", drive[:1]),
			value:     float64(lpFreeBytesAvailable),
			timestamp: now,
		}
		metrics = append(metrics, used)
		metrics = append(metrics, free)
		metrics = append(metrics, available)
	}

	return metrics, nil
}

func (m *DiskspaceInputModule) GetLogicalDrives() []string {

	_getLogicalDrives := kernel32.NewProc("GetLogicalDrives")

	n, _, _ := _getLogicalDrives.Call()
	s := strconv.FormatInt(int64(n), 2)

	var drives_all = []string{"A:", "B:", "C:", "D:", "E:", "F:", "G:", "H:", "I:", "J:", "K:", "L:", "M:", "N:", "O:", "P：", "Q：", "R：", "S：", "T：", "U：", "V：", "W：", "X：", "Y：", "Z："}
	temp := drives_all[0:len(s)]

	var d []string
	for i, v := range s {

		if v == 49 {
			l := len(s) - i - 1
			d = append(d, temp[l])
		}
	}

	var drives []string
	for i, v := range d {
		drives = append(drives[i:], append([]string{v}, drives[:i]...)...)
	}
	return drives

}
