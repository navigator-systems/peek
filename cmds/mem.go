package cmds

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// MemoryInfo stores raw kilobytes parsed from /proc/meminfo
type MemoryInfo struct {
	Total        uint64
	Free         uint64
	Available    uint64
	Buffers      uint64
	Cached       uint64
	SReclaimable uint64 // Reclaimable kernel allocations (part of the actual cache)
	SwapTotal    uint64
	SwapFree     uint64
}

func MEM() {
	fmt.Println("--- Linux Memory Fetch ---")

	mem, err := parseMemInfo()
	if err != nil {
		fmt.Printf("Error reading memory: %v\n", err)
		return
	}

	// --- 1. Math Formulas (Matching modern 'free' command logic) ---
	// Total Cache = Cached file pages + Reclaimable slab memory
	totalCache := mem.Cached + mem.SReclaimable

	// Used Memory = Total - Free - Buffers - Total Cache
	usedMem := mem.Total - mem.Free - mem.Buffers - totalCache

	// Swap Used = SwapTotal - SwapFree
	usedSwap := mem.SwapTotal - mem.SwapFree

	// --- 2. Print System Info (Neofetch Style) ---
	fmt.Printf("Memory (Used):    %s / %s\n", formatKiB(usedMem), formatKiB(mem.Total))
	fmt.Printf("Memory (Avail):   %s\n", formatKiB(mem.Available))
	fmt.Printf("Memory (Free):    %s\n", formatKiB(mem.Free))
	fmt.Printf("Buff/Cache:       %s (Disk buffers: %s | Cache: %s)\n", formatKiB(mem.Buffers+totalCache), formatKiB(mem.Buffers), formatKiB(totalCache))

	if mem.SwapTotal > 0 {
		fmt.Printf("Swap Space:       %s / %s\n", formatKiB(usedSwap), formatKiB(mem.SwapTotal))
	}
}

// Parses /proc/meminfo line by line
func parseMemInfo() (*MemoryInfo, error) {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	mem := &MemoryInfo{}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		// Remove the trailing colon from key (e.g., "MemTotal:")
		key := strings.TrimSuffix(fields[0], ":")
		value, err := strconv.ParseUint(fields[1], 10, 64)
		if err != nil {
			continue
		}

		switch key {
		case "MemTotal":
			mem.Total = value
		case "MemFree":
			mem.Free = value
		case "MemAvailable":
			mem.Available = value
		case "Buffers":
			mem.Buffers = value
		case "Cached":
			mem.Cached = value
		case "SReclaimable":
			mem.SReclaimable = value
		case "SwapTotal":
			mem.SwapTotal = value
		case "SwapFree":
			mem.SwapFree = value
		}
	}

	return mem, scanner.Err()
}

// Helper function to format raw KiB into MiB or GiB strings
func formatKiB(kib uint64) string {
	const mib = 1024
	const gib = 1024 * 1024

	if kib >= gib {
		return fmt.Sprintf("%.2f GiB", float64(kib)/float64(gib))
	}
	return fmt.Sprintf("%.1f MiB", float64(kib)/float64(mib))
}
