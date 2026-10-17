package cmds

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/shirou/gopsutil/v4/process"
)

type procInfo struct {
	Name    string   `json:"name"`
	User    string   `json:"user"`
	Status  []string `json:"status"`
	PID     int32    `json:"pid"`
	Command string   `json:"command"`
	CPU     float64  `json:"cpu"`
	Mem     float32  `json:"mem"`
}

type procDetails struct {
	Name    string   `json:"name"`
	User    string   `json:"user"`
	Status  []string `json:"status"`
	PID     int32    `json:"pid"`
	Command string   `json:"command"`
}

func PROC(action string, pidFlag string, jsonOutput bool) error {
	proc, err := process.Processes()
	if err != nil {
		return fmt.Errorf("unable to get processes info: %w", err)
	}

	if action == "pid" {
		pid, err := strconv.ParseInt(pidFlag, 10, 32)
		if err != nil || pid <= 0 {
			return fmt.Errorf("invalid process ID %q", pidFlag)
		}
		for _, p := range proc {
			if p.Pid == int32(pid) {
				info, err := processDetails(p)
				if err != nil {
					return fmt.Errorf("unable to get process %d info: %w", pid, err)
				}
				if jsonOutput {
					return printProcJSON(info)
				}
				printProcDetails(info)
				return nil
			}
		}
		return fmt.Errorf("process %d not found", pid)
	}

	// Keep the same process handles so the second CPU reading uses the first
	// reading as its baseline.
	sampled := make([]*process.Process, 0, len(proc))
	for _, p := range proc {
		if _, err := p.CPUPercent(); err == nil {
			sampled = append(sampled, p)
		}
	}

	time.Sleep(time.Second)
	infoProcs := make([]procInfo, 0, len(sampled))
	for _, p := range sampled {
		info, err := sampleProcess(p)
		if err != nil {
			continue
		}
		infoProcs = append(infoProcs, info)
	}

	switch action {
	case "cpu":
		sort.Slice(infoProcs, func(i, j int) bool {
			return infoProcs[i].CPU > infoProcs[j].CPU
		})
		infoProcs = topProcesses(infoProcs)
	case "mem":
		sort.Slice(infoProcs, func(i, j int) bool {
			return infoProcs[i].Mem > infoProcs[j].Mem
		})
		infoProcs = topProcesses(infoProcs)
	}

	if jsonOutput {
		return printProcJSON(infoProcs)
	}
	printNice(infoProcs)
	return nil
}

func sampleProcess(p *process.Process) (procInfo, error) {
	name, err := p.Name()
	if err != nil {
		return procInfo{}, err
	}
	user, err := p.Username()
	if err != nil {
		return procInfo{}, err
	}
	status, err := p.Status()
	if err != nil {
		return procInfo{}, err
	}
	cmd, err := p.Cmdline()
	if err != nil {
		return procInfo{}, err
	}
	cpu, err := p.CPUPercent()
	if err != nil {
		return procInfo{}, err
	}
	mem, err := p.MemoryPercent()
	if err != nil {
		return procInfo{}, err
	}

	return procInfo{
		Name:    name,
		User:    user,
		Status:  status,
		PID:     p.Pid,
		Command: cmd,
		CPU:     cpu,
		Mem:     mem,
	}, nil
}

func printProcJSON(info any) error {
	encoded, err := json.Marshal(info)
	if err != nil {
		return fmt.Errorf("unable to encode process info as JSON: %w", err)
	}
	fmt.Println(string(encoded))
	return nil
}

func processDetails(p *process.Process) (procDetails, error) {
	name, err := p.Name()
	if err != nil {
		return procDetails{}, err
	}
	user, err := p.Username()
	if err != nil {
		return procDetails{}, err
	}
	status, err := p.Status()
	if err != nil {
		return procDetails{}, err
	}
	cmd, err := p.Cmdline()
	if err != nil {
		return procDetails{}, err
	}
	return procDetails{
		Name:    name,
		User:    user,
		Status:  status,
		PID:     p.Pid,
		Command: cmd,
	}, nil
}

func printProcDetails(p procDetails) {
	fmt.Println("Process")
	fmt.Printf("Name:    %s\n", p.Name)
	fmt.Printf("PID:     %d\n", p.PID)
	fmt.Printf("User:    %s\n", p.User)
	fmt.Printf("Status:  %v\n", p.Status)
	fmt.Printf("Command: %s\n", p.Command)
}

func topProcesses(info []procInfo) []procInfo {
	if len(info) > 10 {
		return info[:10]
	}
	return info
}

func printNice(info []procInfo) {
	fmt.Printf("%-20s %-6s %-6s %-6s %-20s\n", "Name", "PID", "CPU%", "Mem%", "Command")
	for _, i := range info {
		fmt.Printf("%-20s %-6d %-6.3f %-6.3f %s\n", i.Name, i.PID, i.CPU, i.Mem, i.Command)
	}
}
