package test

import (
	"bufio"
	"log"
	"net"
	"os"
	"sia/backend/lib"
	"strconv"
	"time"
)

// sendData sends data over the UDP connection with a specific dataType.
func sendData(conn *net.UDPConn, dataType int, data []byte) {
	packet := append([]byte{byte(dataType)}, data...)
	_, err := conn.Write(packet)
	if err != nil {
		log.Printf("Failed to send data: %v", err)
	}
}

// RunTestUDPClient reads data from a file line by line and sends it over UDP at a specified frequency.
func RunTestUDPClient(filePath string, frequencyHz int) {
	lib.Print(lib.UDP_SERVICE, "Starting UDP client")

	serverAddr := net.UDPAddr{
		Port: int(lib.UDP_PORT),
		IP:   net.ParseIP(lib.UDP_ADDR),
	}
	conn, err := net.DialUDP("udp", nil, &serverAddr)
	if err != nil {
		lib.Print(lib.UDP_SERVICE, "Failed to connect to server")
		return
	}
	defer conn.Close()

	// Calculate the interval between sends based on the frequency
	interval := time.Second / time.Duration(frequencyHz)

	// Open the file for streaming
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// Read each line as data
		line := scanner.Text()
		value, err := strconv.Atoi(line)
		if err != nil {
			log.Printf("Failed to convert line to integer: %v", err)
			continue
		}

		// Send data
		sendData(conn, 0, []byte{byte(value)})
		// lib.Print(lib.UDP_SERVICE, "Sent EKG sensor data:", value)

		// Wait for the next interval
		time.Sleep(interval)
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading file: %v", err)
	}
}
