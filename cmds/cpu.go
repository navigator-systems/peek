package cmds

import (
	"fmt"
	"log"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/host"
)

func CPU() {

	cpu, err := cpu.Info()
	if err != nil {
		log.Fatalf("unable to get cpu info: %v", err)
	}

	host, err := host.Info()
	if err != nil {
		log.Fatalf("unable to get host info: %v", err)
	}

	fmt.Printf("OS: %s - Platform %s\n", host.OS, host.Platform)
	fmt.Printf("Kernel: %s - %s\n", host.KernelVersion, host.KernelArch)
	fmt.Printf("CPU: %s - Cores %v\n", cpu[0].ModelName, len(cpu))

}
