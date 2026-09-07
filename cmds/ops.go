package cmds

import "fmt"

// humanSize turns bytes into a human‑readable string.
func humanSize(b uint64) string {
	const (
		KB float64 = 1e3
		MB float64 = 1e6
		GB float64 = 1e9
		TB float64 = 1e12
		PB float64 = 1e15
	)

	switch {
	case b >= uint64(PB):
		return fmt.Sprintf("%.2f PB", float64(b)/PB)
	case b >= uint64(TB):
		return fmt.Sprintf("%.2f TB", float64(b)/TB)
	case b >= uint64(GB):
		return fmt.Sprintf("%.2f GB", float64(b)/GB)
	case b >= uint64(MB):
		return fmt.Sprintf("%.2f MB", float64(b)/MB)
	case b >= uint64(KB):
		return fmt.Sprintf("%.2f KB", float64(b)/KB)
	default:
		return fmt.Sprintf("%d B", b)
	}
}
