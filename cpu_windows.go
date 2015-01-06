// +build windows

package main

import (
	"fmt"
	"runtime"
	"unsafe"
)

func (m *CPUInputModule) GetMetrics() ([]Metric, error) {

	cpuNum := runtime.NumCPU()

	var perfSize SYSTEM_PROCESSOR_PERFORMANCE_INFORMATION

	performanceInformation := make([]SYSTEM_PROCESSOR_PERFORMANCE_INFORMATION, cpuNum)
	sizeofPerformanceInformation := uint64(unsafe.Sizeof(perfSize)) * uint64(cpuNum)

	r1, _, err := _ntQuerySystemInformation.Call(SystemProcessorPerformanceInformation, uintptr(unsafe.Pointer(&performanceInformation[0])), uintptr(sizeofPerformanceInformation), uintptr(0))

	if !(r1 >= 0) {
		return nil, err
	}

	var SysTime uint64
	var userTime uint64
	var idleTime uint64
	var intTime uint64

	var interruptCount uint32

	content := string("")

	for i := 0; i < cpuNum; i++ {
		SysTime += (performanceInformation[i].KernelTime - performanceInformation[i].IdleTime) * 100 / 10000000
		idleTime += performanceInformation[i].IdleTime * 100 / 10000000
		userTime += performanceInformation[i].UserTime * 100 / 10000000
		intTime += performanceInformation[i].InterruptTime * 100 / 10000000
	}

	content += fmt.Sprintf("cpu %d %d %d %d %d %d %d %d \n", userTime, 0, SysTime, idleTime, 0, intTime, 0, 0)

	SysTime = 0
	userTime = 0
	idleTime = 0
	intTime = 0

	for i := 0; i < cpuNum; i++ {
		interruptCount += performanceInformation[i].InterruptCount

		SysTime = (performanceInformation[i].KernelTime - performanceInformation[i].IdleTime) * 100 / 10000000
		idleTime = performanceInformation[i].IdleTime * 100 / 10000000
		userTime = performanceInformation[i].UserTime * 100 / 10000000
		intTime = performanceInformation[i].InterruptTime * 100 / 10000000

		content += fmt.Sprintf("cpu%d %d %d %d %d %d %d %d %d \n", i, userTime, 0, SysTime, idleTime, 0, intTime, 0, 0)
	}

	spi := SYSTEM_PERFORMANCE_INFORMATION{}

	sizeof_spi := unsafe.Sizeof(spi)

	r1, _, err = _ntQuerySystemInformation.Call(SystemPerformanceInformation, uintptr(unsafe.Pointer(&spi)), uintptr(sizeof_spi), uintptr(0))
	if !(r1 >= 0) {

		return nil, err
	}

	var timeSize SYSTEM_TIMEOFDAY_INFORMATION
	sizeofTimeInformation := uint64(unsafe.Sizeof(timeSize))

	r1, _, err = _ntQuerySystemInformation.Call(SystemTimeOfDayInformation, uintptr(unsafe.Pointer(&timeSize)), uintptr(sizeofTimeInformation), uintptr(0))

	if !(r1 >= 0) {
		return nil, err
	}

	var pagesIn uint32 = spi.PageReadCount
	var pagesOut uint32 = spi.PageReadIoCount + spi.MappedWriteIoCount
	var swapIn uint32 = spi.PageReadCount
	var swapout uint32 = spi.PageReadIoCount
	var contextSwitches uint32 = spi.ContextSwitches
	var bootTime int64 = m.GetBootTime(timeSize.BootTime)

	content += fmt.Sprintf("page %d %d \n", pagesIn, pagesOut)
	content += fmt.Sprintf("swap %d %d \n", swapIn, swapout)
	content += fmt.Sprintf("intr %d \n", interruptCount)
	content += fmt.Sprintf("ctxt %d \n", contextSwitches)
	content += fmt.Sprintf("btime %d \n", bootTime)

	return m.ParseProcStat(content)
}

func (m *CPUInputModule) GetBootTime(inTime uint64) int64 {

	x := inTime

	if inTime == 0 {
		return 0
	}

	x -= 0x19db1ded53e8000
	x /= 10000000

	return int64(x)
}
