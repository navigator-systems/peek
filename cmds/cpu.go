package cmds

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/host"
)

type CPUStruct struct {
	OS       string `json:"os"`
	Platform string `json:"platform"`
	Kernel   string `json:"kernel"`
	Arch     string `json:"arch"`
	Model    string `json:"model"`
	Cores    int    `json:"cores"`
}

func CPU(jsonOutput bool) {

	cpu, err := cpu.Info()
	if err != nil {
		log.Fatalf("unable to get cpu info: %v", err)
	}

	host, err := host.Info()
	if err != nil {
		log.Fatalf("unable to get host info: %v", err)
	}

	if jsonOutput {
		cpuInfo := CPUStruct{
			OS:       host.OS,
			Platform: host.Platform,
			Kernel:   host.KernelVersion,
			Arch:     host.KernelArch,
			Model:    cpu[0].ModelName,
			Cores:    len(cpu),
		}
		encoded, err := json.Marshal(cpuInfo)
		if err != nil {
			log.Fatalf("unable to encode cpu info as json: %v", err)
		}
		fmt.Println(string(encoded))
		return
	}

	fmt.Printf("OS: %s - Platform %s\n", host.OS, host.Platform)
	fmt.Printf("Kernel: %s - %s\n", host.KernelVersion, host.KernelArch)
	fmt.Printf("CPU: %s - Cores %v\n", cpu[0].ModelName, len(cpu))
}
