package cmds

import (
	"fmt"
	"log"
	"sort"

	"github.com/shirou/gopsutil/v4/process"
)

type procInfo struct {
	Name    string
	PID     int32
	Command string
	CPU     float64
	Mem     float32
}

func PROC(action string) {

	proc, err := process.Processes()
	if err != nil {
		log.Fatalf("unable to get processes info: %v", err)
	}

	var infoProcs []procInfo

	for i := range proc {
		procID := proc[i].Pid
		p, _ := process.NewProcess(procID)
		name, _ := p.Name()
		cmd, _ := p.Cmdline()
		cpu, _ := p.CPUPercent()
		mem, _ := p.MemoryPercent()
		infoProcs = append(infoProcs, procInfo{
			Name:    name,
			PID:     procID,
			Command: cmd,
			CPU:     cpu,
			Mem:     mem,
		})
	}

	switch action {
	case "cpu":
		ProcByCPU(infoProcs)
	case "mem":
		ProcByMem(infoProcs)
	default:
		printNice(infoProcs)
	}

}

func ProcByCPU(info []procInfo) {
	var infoProcs []procInfo

	// Ordenar de mayor a menor CPU
	sort.Slice(info, func(i, j int) bool {
		return info[i].CPU > info[j].CPU
	})

	// Mostrar los primeros 10
	for i, p := range info {
		if i == 10 {
			break
		}

		infoProcs = append(infoProcs, procInfo{
			Name:    p.Name,
			PID:     p.PID,
			Command: p.Command,
			CPU:     p.CPU,
			Mem:     p.Mem,
		})
	}

	printNice(infoProcs)
}

func ProcByMem(info []procInfo) {
	var infoProcs []procInfo

	// Ordenar de mayor a menor CPU
	sort.Slice(info, func(i, j int) bool {
		return info[i].Mem > info[j].Mem
	})

	// Mostrar los primeros 10
	for i, p := range info {
		if i == 10 {
			break
		}

		infoProcs = append(infoProcs, procInfo{
			Name:    p.Name,
			PID:     p.PID,
			Command: p.Command,
			CPU:     p.CPU,
			Mem:     p.Mem,
		})
	}
	printNice(infoProcs)

}

func printNice(info []procInfo) {
	fmt.Printf("%-20s %-6s %-6s %-6s %-20s\n", "Name", "PID", "CPU", "Mem", "Command")
	for _, i := range info {
		fmt.Printf("%-20s %-6v %.3f %.3f %-20s\n", i.Name, i.PID, i.CPU, i.Mem, i.Command)
	}
}
