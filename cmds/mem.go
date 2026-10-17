package cmds

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/shirou/gopsutil/v4/mem"
)

type memStruct struct {
	Total       string  `json:"total"`
	Free        string  `json:"free"`
	Used        string  `json:"used"`
	UsedPercent float32 `json:"used_percent"`
}

func MEM(jsonOutput bool) {
	v, err := mem.VirtualMemory()
	if err != nil {
		log.Fatalf("could not get virtual memory: %v", err)
	}

	total := humanSize(v.Total)

	free := humanSize(v.Free)
	used := humanSize(v.Used)

	usedPer := float32(v.UsedPercent)
	if jsonOutput {
		memInfo := memStruct{
			Total:       total,
			Free:        free,
			Used:        used,
			UsedPercent: usedPer,
		}
		encoded, err := json.Marshal(memInfo)
		if err != nil {
			log.Fatalf("unable to encode memory info as JSON: %v", err)
		}
		fmt.Println(string(encoded))
		return
	} else {
		fmt.Printf("Total: %s, Free:%s, Used:%s, UsedPercent:%.3f%%\n", total, free, used, usedPer)
	}

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
