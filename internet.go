package main

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"time"
)

func getNetworkStats() (int64, int64, error) {
	// Run ifstat to get current network statistics (received and sent bytes)
	cmd := exec.Command("ifstat", "-i", "eth0", "1", "1") // Adjust interface name ("eth0") as needed
	output, err := cmd.CombinedOutput()
	if err != nil {
		return 0, 0, err
	}

	// Process output
	re := regexp.MustCompile(`\s+(\d+)\s+(\d+)`)
	matches := re.FindStringSubmatch(string(output))
	if len(matches) < 3 {
		return 0, 0, fmt.Errorf("failed to parse network stats")
	}

	// Convert matched bytes to integers
	bytesReceived, err := strconv.ParseInt(matches[1], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to convert received bytes: %v", err)
	}

	bytesSent, err := strconv.ParseInt(matches[2], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to convert sent bytes: %v", err)
	}

	return bytesReceived, bytesSent, nil
}

func main() {
	var totalReceived int64
	var totalSent int64

	// Start tracking network usage over time
	for {
		// Get the current network stats
		received, sent, err := getNetworkStats()
		if err != nil {
			fmt.Printf("Error getting network stats: %v\n", err)
			return
		}

		// Calculate the difference in received and sent bytes since last check
		totalReceived += received
		totalSent += sent

		// Print the accumulated usage
		fmt.Printf("Total data received: %d bytes\n", totalReceived)
		fmt.Printf("Total data sent: %d bytes\n", totalSent)

		// Wait for 1 minute before next update
		time.Sleep(1 * time.Minute)
	}
}
