package cmds

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

func CPU() {
	// 1. Get Operating System and Core Count using the runtime package
	osName := runtime.GOOS
	cores := runtime.NumCPU()

	var cpuModel string

	cpuModel = getLinuxCPU()

	fmt.Println("OS: ", osName)
	fmt.Printf("CPU: %s - %d Cores\n", cpuModel, cores)
}

// Linux implementation: Parsers /proc/cpuinfo
func getLinuxCPU() string {
	data, err := os.ReadFile("/proc/cpuinfo")
	if err != nil {
		return "Unknown Linux CPU"
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		// Look for the "model name" line
		if strings.HasPrefix(line, "model name") || strings.HasPrefix(line, "Processor") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				return strings.TrimSpace(parts[1])
			}
		}
	}
	return "Unknown Linux CPU"
}
