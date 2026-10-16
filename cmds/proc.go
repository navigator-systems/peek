package cmds

import (
	"fmt"
	"log"
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

func PROC(action string, pidFlag string) {

	proc, err := process.Processes()
	if err != nil {
		log.Fatalf("unable to get processes info: %v", err)
	}

	for _, p := range proc {
		_, err := p.CPUPercent()
		if err != nil {
			continue
		}
	}

	time.Sleep(1 * time.Second)
	var infoProcs []procInfo

	for i := range proc {
		procID := proc[i].Pid
		p, _ := process.NewProcess(procID)
		name, _ := p.Name()
		user, _ := p.Username()
		status, _ := p.Status()
		cmd, _ := p.Cmdline()
		cpu, _ := p.CPUPercent()
		mem, _ := p.MemoryPercent()
		infoProcs = append(infoProcs, procInfo{
			Name:    name,
			User:    user,
			Status:  status,
			PID:     procID,
			Command: cmd,
			CPU:     cpu,
			Mem:     mem,
		})
	}

	switch action {
	case "cpu":
		procByCPU(infoProcs)
	case "mem":
		procByMem(infoProcs)
	case "pid":
		infoProc(infoProcs, pidFlag)
	default:
		printNice(infoProcs)
	}

}

func infoProc(info []procInfo, pid string) {
	PID, err := strconv.ParseInt(pid, 10, 32)
	PID32 := int32(PID)
	if err != nil {
		log.Fatalln("Error reading value", pid)
	}
	for _, p := range info {
		if p.PID == PID32 {
			fmt.Println("Process")
			fmt.Printf("Name:    %s\n", p.Name)
			fmt.Printf("PID:     %v\n", PID32)
			fmt.Printf("User:    %s\n", p.User)
			fmt.Printf("Status:  %v\n", p.Status)
			fmt.Printf("Command: %s\n", p.Command)

		}
	}

}

func procByCPU(info []procInfo) {
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

func procByMem(info []procInfo) {
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
	fmt.Printf("%-20s %-6s %-6s %-6s %-20s\n", "Name", "PID", "CPU%", "Mem%", "Command")
	for _, i := range info {
		fmt.Printf("%-20s %-6v %.3f %.3f %-20s\n", i.Name, i.PID, i.CPU, i.Mem, i.Command)
	}
}
