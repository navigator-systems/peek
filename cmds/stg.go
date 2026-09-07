package cmds

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/shirou/gopsutil/v4/disk"
)

func STG() {

	partitions, err := disk.Partitions(false)

	if err != nil {
		log.Fatalf("unable to read partitions: %v", err)
	}

	// TabWriter for a clean table layout
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "MOUNT\tTYPE\tTOTAL\tUSED\tFREE")

	for _, p := range partitions {
		usage, err := disk.Usage(p.Mountpoint)
		if err != nil {
			// Skip mount points we can't read (e.g., /proc, /sys)
			continue
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			p.Mountpoint,
			p.Fstype,
			humanSize(usage.Total),
			humanSize(usage.Used),
			humanSize(usage.Free))
	}
	w.Flush()

	// Summary: total/used/free across all partitions
	total, used, free := uint64(0), uint64(0), uint64(0)
	for _, p := range partitions {
		usage, err := disk.Usage(p.Mountpoint)
		if err != nil {
			continue
		}
		total += usage.Total
		used += usage.Used
		free += usage.Free
	}
	fmt.Println("\nTOTAL STORAGE OVERVIEW")
	fmt.Printf("Total: %s  Used: %s  Free: %s\n",
		humanSize(total), humanSize(used), humanSize(free))
}
