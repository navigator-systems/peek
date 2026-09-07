package cmds

import (
	"fmt"
	"log"

	"github.com/shirou/gopsutil/v4/mem"
)

func MEM() {
	v, err := mem.VirtualMemory()
	if err != nil {
		log.Fatalf("could not get virtual memory: %v", err)
	}

	total := humanSize(v.Total)

	free := humanSize(v.Free)
	used := humanSize(v.Used)

	usedPer := float32(v.UsedPercent)

	fmt.Printf("Total: %s, Free:%s, Used:%s, UsedPercent:%.3f%%\n", total, free, used, usedPer)
}

// Helper function to format raw KiB into MiB or GiB strings
func formatKiB(kb uint64) string {

	var mb, gb float64
	mb = float64(kb) / 1000
	gb = mb / 1000 / 1000

	if float64(kb) >= gb {
		return fmt.Sprintf("%.2f Gb", gb)
	}
	return fmt.Sprintf("%.2f Mb", mb)
}
