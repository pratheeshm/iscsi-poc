// Package main implements a simple iSCSI target server.
// The server listens on a specified port and handles read and write operations
// from connected clients. The storage is simulated using an in-memory byte slice.

package main

import (
	"fmt"
	"log"
	"net"
	"sync"
)

const (
	port           = ":3260" // Port on which the iSCSI target server listens
	numBlocks      = 1024    // Number of blocks in the storage
	operationBytes = 1       // Number of bytes used to represent the operation code (R/W)
)

var (
	storage = make([]byte, blockSize*numBlocks) // In-memory storage
	mu      sync.Mutex                          // Mutex to synchronize access to storage
)

// target function starts the iSCSI target server and listens for incoming connections.
func target() {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", port, err)
	}
	defer listener.Close()

	fmt.Printf("iSCSI Target listening on port %s\n", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		go handleConnection(conn)
	}
}

// handleConnection handles an individual client connection.
func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Printf("Client connected: %s\n", conn.RemoteAddr().String())

	buf := make([]byte, blockSize+headerBytes) // Adjust buffer size to accommodate the operation code and data
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err.Error() != "EOF" {
				log.Printf("Failed to read from connection: %v", err)
			} else {
				fmt.Printf("Client disconnected: %s\n", conn.RemoteAddr().String())
			}
			return
		}
		// Ignore empty messages
		if n < 1 {
			continue
		}

		fmt.Printf("Received data: %v\n", buf[:n])

		switch buf[0] {
		case 'R': // Read operation
			handleRead(conn, buf[1:n])
		case 'W': // Write operation
			handleWrite(buf[1:n])
		default:
			log.Printf("Unknown operation: %c", buf[0])
		}
	}
}

// handleRead handles a read operation from the client.
func handleRead(conn net.Conn, data []byte) {
	fmt.Printf("Received read request: %v\n", data)
	blockIndex := (int(data[0])<<24 | int(data[1])<<16 | int(data[2])<<8 | int(data[3]))
	if blockIndex < 0 || blockIndex >= numBlocks {
		log.Printf("Invalid block index: %d", blockIndex)
		return
	}
	fmt.Printf("Received read request for block %d\n", blockIndex)
	mu.Lock()
	defer mu.Unlock()

	// Send the data from the specified block to the client
	conn.Write(storage[blockIndex*blockSize : (blockIndex+1)*blockSize])
	fmt.Printf("Sent data from block %d\n", blockIndex)
}

// handleWrite handles a write operation from the client.
func handleWrite(data []byte) {
	blockIndex := (int(data[0])<<24 | int(data[1])<<16 | int(data[2])<<8 | int(data[3]))
	if blockIndex < 0 || blockIndex >= numBlocks {
		log.Printf("Invalid block index: %d", blockIndex)
		return
	}
	fmt.Printf("Received write request for block %d\n", blockIndex)

	mu.Lock()
	defer mu.Unlock()

	copy(storage[blockIndex*blockSize:], data[headerBytes-operationBytes:])
	fmt.Printf("Wrote data to block %d\n", blockIndex)
}
