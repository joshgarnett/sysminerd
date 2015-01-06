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

	SYSTEM_PROCESSOR_PERFORMANCE_INFORMATION struct {
		IdleTime       uint64
		KernelTime     uint64
		UserTime       uint64
		DpcTime        uint64
		InterruptTime  uint64
		InterruptCount uint32
	}

	SYSTEM_PERFORMANCE_INFORMATION struct {
		IdleProcessTime          uint64
		IoReadTransferCount      uint64
		IoWriteTransferCount     uint64
		IoOtherTransferCount     uint64
		IoReadOperationCount     uint32
		IoWriteOperationCount    uint32
		IoOtherOperationCount    uint32
		AvailablePages           uint32
		CommittedPages           uint32
		CommitLimit              uint32
		PeakCommitment           uint32
		PageFaultCount           uint32
		CopyOnWriteCount         uint32
		TransitionCount          uint32
		CacheTransitionCount     uint32
		DemandZeroCount          uint32
		PageReadCount            uint32
		PageReadIoCount          uint32
		CacheReadCount           uint32
		CacheIoCount             uint32
		DirtyPagesWriteCount     uint32
		DirtyWriteIoCount        uint32
		MappedPagesWriteCount    uint32
		MappedWriteIoCount       uint32
		PagedPoolPages           uint32
		NonPagedPoolPages        uint32
		PagedPoolAllocs          uint32
		PagedPoolFrees           uint32
		NonPagedPoolAllocs       uint32
		NonPagedPoolFrees        uint32
		FreeSystemPtes           uint32
		ResidentSystemCodePage   uint32
		TotalSystemDriverPages   uint32
		TotalSystemCodePages     uint32
		NonPagedPoolLookasideHit uint32
		PagedPoolLookasideHits   uint32
		AvailablePagedPoolPages  uint32
		ResidentSystemCachePage  uint32
		ResidentPagedPoolPage    uint32
		ResidentSystemDriverPage uint32
		CcFastReadNoWait         uint32
		CcFastReadWait           uint32
		CcFastReadResourceMiss   uint32
		CcFastReadNotPossible    uint32
		CcFastMdlReadNoWait      uint32
		CcFastMdlReadWait        uint32
		CcFastMdlReadResourceMis uint32
		CcFastMdlReadNotPossible uint32
		CcMapDataNoWait          uint32
		CcMapDataWait            uint32
		CcMapDataNoWaitMiss      uint32
		CcMapDataWaitMiss        uint32
		CcPinMappedDataCount     uint32
		CcPinReadNoWait          uint32
		CcPinReadWait            uint32
		CcPinReadNoWaitMiss      uint32
		CcPinReadWaitMiss        uint32
		CcCopyReadNoWait         uint32
		CcCopyReadWait           uint32
		CcCopyReadNoWaitMiss     uint32
		CcCopyReadWaitMiss       uint32
		CcMdlReadNoWait          uint32
		CcMdlReadWait            uint32
		CcMdlReadNoWaitMiss      uint32
		CcMdlReadWaitMiss        uint32
		CcReadAheadIos           uint32
		CcLazyWriteIos           uint32
		CcLazyWritePages         uint32
		CcDataFlushes            uint32
		CcDataPages              uint32
		ContextSwitches          uint32
		FirstLevelTbFills        uint32
		SecondLevelTbFills       uint32
		SystemCalls              uint32
		CcTotalDirtyPages        uint64
		CcDirtyPageThreshold     uint64
		ResidentAvailablePages   int64
	}

	SYSTEM_TIMEOFDAY_INFORMATION struct {
		BootTime          uint64
		CurrentTime       uint64
		TimeZoneBias      uint64
		CurrentTimeZoneId uint64
		Reserved1         uint64
	}
)

const (
	SystemBasicInformation                = 0
	SystemPerformanceInformation          = 2
	SystemTimeOfDayInformation            = 3
	SystemProcessInformation              = 5
	SystemProcessorPerformanceInformation = 8
	SystemInterruptInformation            = 23
	SystemExceptionInformation            = 33
	SystemRegistryQuotaInformation        = 37
	SystemLookasideInformation            = 45
	SystemPolicyInformation               = 134
)

var kernel32 = syscall.NewLazyDLL("kernel32.dll")
var ntDll = syscall.NewLazyDLL("Ntdll.dll")

var (
	_globalMemoryStatusEx     = kernel32.NewProc("GlobalMemoryStatusEx")
	_getDiskFreeSpaceEx       = kernel32.NewProc("GetDiskFreeSpaceExW")
	_getLogicalDrives         = kernel32.NewProc("GetLogicalDrives")
	_ntQuerySystemInformation = ntDll.NewProc("NtQuerySystemInformation")
)
