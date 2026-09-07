package cmds

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/shirou/gopsutil/v4/net"
	"github.com/shirou/gopsutil/v4/process"
)

// connectionInfo holds the data we’ll print for a single listening socket.
type connectionInfo struct {
	LocalAddr string // e.g. "0.0.0.0:8080"
	Port      uint32
	Proto     string // "tcp", "udp", ...
	PID       int32
	ProcName  string
}

func PORT(portFlag string) {

	// Gather all listening connections.
	conns, err := net.Connections("inet")
	if err != nil {
		log.Fatalf("could not fetch connections: %v", err)
	}

	var infos []connectionInfo

	for _, c := range conns {
		if c.Status != "LISTEN" {
			continue
		}
		if c.Pid == 0 {
			continue
		}

		// Convert uint32 type to string (1 = stream/tcp, 2 = dgram/udp, etc.)
		var proto string
		switch c.Type {
		case 1:
			proto = "TCP"
		case 2:
			proto = "UDP"
		case 3:
			proto = "UNIX"
		default:
			proto = fmt.Sprintf("UNKNOWN(%d)", c.Type)
		}

		infos = append(infos, connectionInfo{
			LocalAddr: fmt.Sprintf("%s:%d", c.Laddr.IP, c.Laddr.Port),
			Port:      c.Laddr.Port,
			Proto:     proto,
			PID:       c.Pid,
		})
	}

	// If the user supplied a port filter, prune the list.
	if portFlag != "" {
		p, err := strconv.ParseUint(portFlag, 10, 32)
		if err != nil {
			log.Fatalf("invalid port: %v", err)
		}
		var filtered []connectionInfo
		for _, i := range infos {
			if uint32(p) == i.Port {
				filtered = append(filtered, i)
			}
		}
		infos = filtered
	}

	for i := range infos {
		PID := infos[i].PID
		p, err := process.NewProcess(PID)
		if err != nil {
			infos[i].ProcName = "???"
		}
		name, _ := p.Name()
		infos[i].ProcName = name
	}

	if len(infos) == 0 {
		if portFlag != "" {
			fmt.Fprintf(os.Stderr, "No listening socket found for port %s\n", portFlag)
			return
		} else {
			fmt.Println("No listening sockets found.")
			return
		}
	}

	fmt.Printf("%-20s %-6s %-6s %-15s\n", "Local Address", "Proto", "PID", "Process")
	for _, i := range infos {
		fmt.Printf("%-20s %-6s %-6d %-15s\n", i.LocalAddr, i.Proto, i.PID, i.ProcName)
	}
}
