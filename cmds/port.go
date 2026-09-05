package cmds

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"

	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
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
	/*	for _, c := range conns {
			if c.Status != "LISTEN" {
				continue
			}
			// Skip connections that don't have a PID (e.g. kernel sockets).
			if c.Pid == 0 {
				continue
			}
			infos = append(infos, connectionInfo{
				LocalAddr: fmt.Sprintf("%s:%s", c.LocalAddr, c.LocalPort),
				Port:      c.LocalPort,
				Proto:     strings.ToUpper(c.Type),
				PID:       c.Pid,
			})
		}
	*/
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

	// Sort by port for readability.
	sort.Slice(infos, func(i, j int) bool { return infos[i].Port < infos[j].Port })

	// Resolve the process name for each PID (cached to avoid duplicate look‑ups).
	pidToName := make(map[int32]string)
	for _, i := range infos {
		if name, ok := pidToName[i.PID]; ok {
			i.ProcName = name
		} else {
			p, err := process.NewProcess(i.PID)
			if err != nil {
				i.ProcName = "???"
			} else {
				name, _ := p.Name()
				i.ProcName = name
			}
			pidToName[i.PID] = i.ProcName
		}
	}

	// Print a nice table header.
	fmt.Printf("%-20s %-6s %-6s %-15s\n", "Local Address", "Proto", "PID", "Process")
	for _, i := range infos {
		fmt.Printf("%-20s %-6s %-6d %-15s\n", i.LocalAddr, i.Proto, i.PID, i.ProcName)
	}

	if len(infos) == 0 {
		if portFlag != "" {
			fmt.Fprintf(os.Stderr, "No listening socket found for port %s\n", portFlag)
		} else {
			fmt.Println("No listening sockets found.")
		}
	}
}
