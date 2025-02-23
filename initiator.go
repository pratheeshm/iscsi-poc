package main

import (
	"fmt"
	"log"
	"net"
)

const (
	targetAddress = "localhost:3260"
	blockNumber   = 256 // Block number to write to and read from
)

// initiator is the entry point of the iSCSI initiator program. It establishes a TCP connection
// to the specified iSCSI target address, writes data to block #blockNumber, and then reads data from block #blockNumber.
// The program logs any errors encountered during these operations and prints the results to the console.
func initiator() {
	conn, err := net.Dial("tcp", targetAddress)
	if err != nil {
		log.Fatalf("Failed to connect to target: %v", err)
	}
	defer conn.Close()

	fmt.Printf("Connected to iSCSI Target at %s\n", targetAddress)

	// Write data to block #blockNumber
	writeData := make([]byte, headerBytes+blockSize)
	writeData[0] = 'W'
	writeData[1] = byte(blockNumber >> 24)
	writeData[2] = byte(blockNumber >> 16)
	writeData[3] = byte(blockNumber >> 8)
	writeData[4] = byte(blockNumber & 0xFF)
	copy(writeData[headerBytes:], []byte("Hello, iSCSI Target!"))

	// Ensure the writeData buffer is exactly 516 bytes
	if len(writeData) < headerBytes+blockSize {
		writeData = append(writeData, make([]byte, headerBytes+blockSize-len(writeData))...)
	}

	fmt.Printf("Sending write request: %v\n", writeData)

	_, err = conn.Write(writeData)
	if err != nil {
		log.Fatalf("Failed to write to target: %v", err)
	}
	fmt.Printf("Wrote data to block %d\n", blockNumber)

	// Read data from block #blockNumber
	readData := make([]byte, headerBytes)
	readData[0] = 'R'
	readData[1] = byte(blockNumber >> 24)
	readData[2] = byte(blockNumber >> 16)
	readData[3] = byte(blockNumber >> 8)
	readData[4] = byte(blockNumber & 0xFF)
	fmt.Printf("Sending read request: %v\n", readData)

	_, err = conn.Write(readData)
	if err != nil {
		log.Fatalf("Failed to send read request to target: %v", err)
	}

	buf := make([]byte, blockSize)
	n, err := conn.Read(buf)
	if err != nil {
		log.Fatalf("Failed to read from target: %v", err)
	}
	fmt.Printf("Read data from block %d: %s\n", blockNumber, string(buf[:n]))
}
