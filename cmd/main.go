package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/shirou/gopsutil/v3/net"
)

// Constants for unit conversions using binary (1024-based) prefixes
const (
	KB = 1024.0
	MB = KB * 1024
	GB = MB * 1024
)

// NetStats represents comprehensive network statistics for a specific network interface.
type NetStats struct {
	Interface  string `json:"interface"`
	SentSpeed  Speed  `json:"sentSpeed"`
	RecvSpeed  Speed  `json:"recvSpeed"`
	TotalSent  Usage  `json:"totalSent"`
	TotalRecv  Usage  `json:"totalRecv"`
	TotalUsage Usage  `json:"totalUsage"`
}

// Speed describes network transfer speed with a numerical value and its unit.
type Speed struct {
	Value float64 `json:"value"`
	Unit  string  `json:"unit"`
}

// Usage represents network data transfer amount with a numerical value and its unit.
type Usage struct {
	Value float64 `json:"value"`
	Unit  string  `json:"unit"`
}

// NetworkMonitor manages the collection and processing of network interface statistics.
type NetworkMonitor struct {
	interfaceName   string         // Name of the network interface being monitored
	refreshInterval int            // Time between statistical updates in seconds
	precision       int            // Number of decimal places for rounding numerical values
	format          string         // Output format ("json" or "table")
	interrupt       chan os.Signal // Channel to handle interrupt signals
	stats           NetStats       // Most recent network statistics
	mu              sync.RWMutex   // Mutex for thread-safe access to stats
}

// NewNetworkMonitor creates and initializes a new NetworkMonitor instance.
func NewNetworkMonitor(iface string, interval, precision int, format string) *NetworkMonitor {
	return &NetworkMonitor{
		interfaceName:   iface,
		refreshInterval: interval,
		precision:       precision,
		format:          format,
		interrupt:       make(chan os.Signal, 1),
	}
}

// round calculates a floating-point number rounded to a specified number of decimal places.
func round(value float64, precision int) float64 {
	multiplier := math.Pow(10, float64(precision))
	return math.Round(value*multiplier) / multiplier
}

// calculateSpeed determines the most appropriate unit for network transfer speed (B/s, KB/s, MB/s, GB/s).
func calculateSpeed(bytes uint64, interval int, precision int) Speed {
	speed := float64(bytes) / float64(interval)

	switch {
	case speed >= GB:
		return Speed{
			Value: round(speed/GB, precision),
			Unit:  "GB/s",
		}
	case speed >= MB:
		return Speed{
			Value: round(speed/MB, precision),
			Unit:  "MB/s",
		}
	case speed >= KB:
		return Speed{
			Value: round(speed/KB, precision),
			Unit:  "KB/s",
		}
	default:
		return Speed{
			Value: round(speed, precision),
			Unit:  "B/s",
		}
	}
}

// calculateUsage determines the most appropriate unit for network data transfer (B, KB, MB, GB).
func calculateUsage(bytes uint64, precision int) Usage {
	usage := float64(bytes)

	switch {
	case usage >= GB:
		return Usage{
			Value: round(usage/GB, precision),
			Unit:  "GB",
		}
	case usage >= MB:
		return Usage{
			Value: round(usage/MB, precision),
			Unit:  "MB",
		}
	case usage >= KB:
		return Usage{
			Value: round(usage/KB, precision),
			Unit:  "KB",
		}
	default:
		return Usage{
			Value: round(usage, precision),
			Unit:  "B",
		}
	}
}

// getInterfaceIOCounters retrieves network I/O statistics for a specific network interface.
func getInterfaceIOCounters(ifaceName string) (net.IOCountersStat, error) {
	netIO, err := net.IOCounters(true)
	if err != nil {
		return net.IOCountersStat{}, err
	}

	for _, io := range netIO {
		if io.Name == ifaceName {
			return io, nil
		}
	}

	return net.IOCountersStat{}, fmt.Errorf("interface not found: %s", ifaceName)
}

// printTable prints the network statistics in a tabular format to the console.
func printTable(stats NetStats, precision int) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Interface", "Sent Speed", "Recv Speed", "Total Sent", "Total Recv", "Total Usage"})

	table.Append([]string{
		stats.Interface,
		fmt.Sprintf("%.*f %s", precision, stats.SentSpeed.Value, stats.SentSpeed.Unit),
		fmt.Sprintf("%.*f %s", precision, stats.RecvSpeed.Value, stats.RecvSpeed.Unit),
		fmt.Sprintf("%.*f %s", precision, stats.TotalSent.Value, stats.TotalSent.Unit),
		fmt.Sprintf("%.*f %s", precision, stats.TotalRecv.Value, stats.TotalRecv.Unit),
		fmt.Sprintf("%.*f %s", precision, stats.TotalUsage.Value, stats.TotalUsage.Unit),
	})

	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetBorder(true)
	table.SetRowLine(true)

	table.Render()
}

// printJSON prints the network statistics in JSON format to the console.
func printJSON(stats NetStats) {
	jsonData, err := json.Marshal(stats)
	if err != nil {
		log.Printf("Error marshaling to JSON: %v", err)
	}
	fmt.Println(string(jsonData))
}

// collectStats continuously gathers and processes network statistics.
func (nm *NetworkMonitor) collectStats() error {
	initialNetIO, err := getInterfaceIOCounters(nm.interfaceName)
	if err != nil {
		return fmt.Errorf("error getting initial network stats: %v", err)
	}

	totalSentStart := initialNetIO.BytesSent
	totalRecvStart := initialNetIO.BytesRecv
	prevNetIO := initialNetIO

	ticker := time.NewTicker(time.Duration(nm.refreshInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			currentNetIO, err := getInterfaceIOCounters(nm.interfaceName)
			if err != nil {
				log.Printf("Error getting network stats: %v", err)
				continue
			}

			sentBytes := currentNetIO.BytesSent - prevNetIO.BytesSent
			recvBytes := currentNetIO.BytesRecv - prevNetIO.BytesRecv

			totalSent := currentNetIO.BytesSent - totalSentStart
			totalRecv := currentNetIO.BytesRecv - totalRecvStart

			stats := NetStats{
				Interface:  nm.interfaceName,
				SentSpeed:  calculateSpeed(sentBytes, nm.refreshInterval, nm.precision),
				RecvSpeed:  calculateSpeed(recvBytes, nm.refreshInterval, nm.precision),
				TotalSent:  calculateUsage(totalSent, nm.precision),
				TotalRecv:  calculateUsage(totalRecv, nm.precision),
				TotalUsage: calculateUsage(totalSent+totalRecv, nm.precision),
			}

			nm.mu.Lock()
			nm.stats = stats
			nm.mu.Unlock()

			if nm.format == "table" {
				printTable(stats, nm.precision)
			} else {
				printJSON(stats)
			}

			prevNetIO = currentNetIO

		case <-nm.interrupt:
			return nil
		}
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	interfaceName := flag.String("i", "", "Network interface to monitor (required)")
	refreshInterval := flag.Int("t", 1, "Refresh interval in seconds")
	precision := flag.Int("p", 2, "Precision for rounding numbers")
	format := flag.String("f", "table", "Output format: json or table")
	flag.Parse()

	if *interfaceName == "" {
		flag.Usage()
		fmt.Print("\n")
		log.Fatal("Error: the -i (interface) flag is required.\n" +
			"Usage: ./zag-netStats -i <interface_name> -t <interval> -p <precision> -f <format>")
	}

	if *precision < 0 || *precision > 6 {
		log.Fatal("Precision must be between 0 and 6 decimal places")
	}

	if *refreshInterval <= 0 || *refreshInterval > 3600 {
		log.Fatal("Refresh interval must be between 1 and 3600 seconds")
	}

	if *format != "json" && *format != "table" {
		log.Fatal("Invalid output format. Allowed values: json, table")
	}

	monitor := NewNetworkMonitor(*interfaceName, *refreshInterval, *precision, *format)

	signal.Notify(monitor.interrupt, os.Interrupt, syscall.SIGTERM)

	if err := monitor.collectStats(); err != nil {
		log.Fatalf("Network monitoring error: %v", err)
	}
}
