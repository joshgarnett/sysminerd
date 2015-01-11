package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

type Process struct {
	Pid    int64
	Name   string
	State  string
	Fields map[string]float64
}

// See http://man7.org/linux/man-pages/man5/proc.5.html
var processStatFields = [...]string{
	"pid",                   // The process ID.
	"comm",                  // The filename of the executable, in parentheses.
	"state",                 // character representing process state
	"ppid",                  // The PID of the parent of this process.
	"pgrp",                  // The process group ID of the process.
	"session",               // The session ID of the process.
	"tty_nr",                // The controlling terminal of the process.
	"tpgid",                 // The ID of the foreground process group.
	"flags",                 // The kernel flags word of the process.
	"minflt",                // The number of minor faults the process has made
	"cminflt",               // The number of minor faults that the process's waited-for children have made.
	"majflt",                // The number of major faults the process has made
	"cmajflt",               // The number of major faults that the process's waited-for children have made.
	"utime",                 // Amount of time that this process has been scheduled in user mode
	"stime",                 // Amount of time that this process has been scheduled in kernel mode
	"cutime",                // Amount of time that this process's waited-for children have been scheduled in user mode
	"cstime",                // Amount of time that this process's waited-for children have been scheduled in kernel mode
	"priority",              // process priority
	"nice",                  // The nice value
	"num_threads",           // Number of threads in this process
	"itrealvalue",           // The time in jiffies before the next SIGALRM is sent to the process due to an interval timer.
	"starttime",             // The time the process started after system boot
	"vsize",                 // Virtual memory size in bytes
	"rss",                   // Resident Set Size: number of pages the process has in real memory
	"rsslim",                // Current soft limit in bytes on the rss of the process
	"startcode",             // The address above which program text can run.
	"endcode",               // The address below which program text can run.
	"startstack",            // The address of the start (i.e., bottom) of the stack.
	"kstkesp",               // The current value of ESP (stack pointer)
	"kstkeip",               // The current EIP (instruction pointer)
	"signal",                // The bitmap of pending signals, displayed as a decimal number
	"blocked",               // The bitmap of blocked signals, displayed as a decimal number
	"sigignore",             // The bitmap of ignored signals, displayed as a decimal number
	"sigcatch",              // The bitmap of caught signals, displayed as a decimal number
	"wchan",                 // This is the "channel" in which the process is waiting
	"nswap",                 // Number of pages swapped
	"cnswap",                // Cumulative nswap for child processes
	"exit_signal",           // Signal to be sent to parent when we die
	"processor",             // CPU number last executed on
	"rt_priority",           // Real-time scheduling priority
	"policy",                // Scheduling policy
	"delayacct_blkio_ticks", // Aggregated block I/O delays
	"guest_time",            // Guest time of the process
	"cguest_time",           // Guest time of the process's children
}

func GetProcessStats(pid int64) (*Process, error) {
	path := fmt.Sprintf("/proc/%d/stat", pid)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	content := string(data)
	fields := strings.Fields(content)

	if len(fields) < 4 {
		return nil, errors.New("Invalid stat file")
	}

	process := Process{}
	process.Pid = pid
	process.Name = fields[1]
	process.State = fields[2]
	process.Fields = make(map[string]float64)
	for i, field := range fields {
		if i < 3 || i >= len(processStatFields) {
			continue
		}

		value, err := strconv.ParseFloat(field, 64)
		if err != nil {
			log.Printf("Error parsing %s as float64: %v", field, err)
		} else {
			process.Fields[processStatFields[i]] = value
		}
	}

	// convert state to full string name
	switch strings.ToUpper(process.State) {
	case "R":
		process.State = "running"
	case "S":
		process.State = "sleeping"
	case "D":
		process.State = "blocked"
	case "Z":
		process.State = "zombies"
	case "T":
		process.State = "stopped"
	case "W":
		process.State = "paging"
	case "X":
		process.State = "dead"
	case "K":
		process.State = "wake_kill"
	case "P":
		process.State = "parked"
	default:
		process.State = "unknown"
	}

	return &process, nil
}
