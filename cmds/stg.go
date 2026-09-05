package cmds

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/shirou/gopsutil/v3/disk"
)

// humanSize turns bytes into a human‑readable string.
func humanSize(b uint64) string {
	const (
		_           = iota
		KiB float64 = 1 << (10 * iota)
		MiB
		GiB
		TiB
		PiB
	)

	switch {
	case b >= uint64(PiB):
		return fmt.Sprintf("%.1f PiB", float64(b)/PiB)
	case b >= uint64(TiB):
		return fmt.Sprintf("%.1f TiB", float64(b)/TiB)
	case b >= uint64(GiB):
		return fmt.Sprintf("%.1f GiB", float64(b)/GiB)
	case b >= uint64(MiB):
		return fmt.Sprintf("%.1f MiB", float64(b)/MiB)
	case b >= uint64(KiB):
		return fmt.Sprintf("%.1f KiB", float64(b)/KiB)
	default:
		return fmt.Sprintf("%d B", b)
	}
}

func STG() {
	// Fetch all mount points (all devices, even USBs, network shares, etc.)
	partitions, err := disk.Partitions(true)
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
