//Danny Radosevich
//COSC3750
//TCP connection example

// This file demonstrates a simple TCP server and TCP client in one program.
//
// Design summary:
// - The server listens on localhost:3070.
// - The client dials localhost:3070 from a goroutine.
// - The client reads input lines from stdin and sends them to the server.
// - The server reads each line, logs it, and echoes a response.
// - Typing "exit" triggers a coordinated, clean shutdown.

package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

// createConnection runs the client side of this demo.
//
// Parameters:
// - addr: server address to dial.
// - done: write-only signal channel; closed when client exits.
//
// Behavior:
// 1) Connect to the local server.
// 2) Start a background reader for server replies.
// 3) Loop reading stdin lines and writing them to the server.
// 4) Exit when stdin ends, write fails, or user types "exit".
func createConnection(addr *net.TCPAddr, done chan<- struct{}) {
	// Open outbound TCP connection to our server listener.
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		// If connect fails, signal done immediately so main can terminate.
		log.Println("Error connecting to", addr, ":", err)
		close(done)
		return
	}

	// Always notify main that client work has ended.
	defer close(done)

	// Ensure socket resources are released.
	defer conn.Close()

	// Final lifecycle message to make shutdown obvious in logs.
	defer log.Println("Closing connection to", addr)
	log.Println("Connected to", addr)

	// Background receive loop for server responses.
	// This enables full duplex behavior: receive while still sending input.
	go func() {
		reader := bufio.NewReader(conn)
		for {
			// Protocol is line-based, so read until newline.
			msg, err := reader.ReadString('\n')
			if err != nil {
				// Most likely normal connection closure during shutdown.
				return
			}
			// Prefix output to distinguish server messages from the prompt.
			fmt.Print("server> ", msg)
		}
	}()

	// Foreground input/send loop.
	r := bufio.NewReader(os.Stdin)
	for {
		// Prompt for user input.
		fmt.Print("# ")

		// Read one full line from stdin (includes trailing newline).
		text, err := r.ReadBytes('\n')
		if err != nil {
			// EOF is common when input is piped and stream ends.
			if !errors.Is(err, io.EOF) {
				log.Println("Error reading from stdin:", err)
			}
			return
		}

		// Send the line to the server.
		_, err = conn.Write(text)
		if err != nil {
			log.Println("Error writing to connection:", err)
			return
		}

		// Local command to stop the client.
		if strings.TrimSpace(string(text)) == "exit" {
			return
		}
	}
}

// handleConnection runs on the server side for one accepted client socket.
//
// For each incoming line:
// - "exit" => send goodbye and close this connection.
// - any other text => log and echo back.
func handleConnection(conn net.Conn) {
	// Always close the accepted connection when handler returns.
	defer conn.Close()

	// Buffered reader for line-based parsing on stream socket.
	reader := bufio.NewReader(conn)

	// Small delay retained from original demo code.
	time.Sleep(time.Second / 2)

	for {
		// Wait for one newline-terminated client message.
		msg, err := reader.ReadString('\n')
		if err != nil {
			log.Println("Error reading from connection:", err)
			return
		}

		// Normalize by trimming spaces/newline for comparisons and logs.
		trim_msg := strings.TrimSpace(msg)

		if trim_msg == "exit" {
			// Optional final acknowledgement before disconnecting.
			_, _ = conn.Write([]byte("goodbye\n"))
			log.Println("Closing connection to", conn.RemoteAddr())
			return
		}

		// Standard echo behavior for non-exit messages.
		log.Println("Received from", conn.RemoteAddr(), ":", trim_msg)
		_, _ = conn.Write([]byte("echo: " + trim_msg + "\n"))
	}
}

func main() {
	// Local server endpoint for both listening and dialing.
	local_addr := "localhost:3070"

	// Convert "host:port" string into a TCPAddr.
	tcp_addr, err := net.ResolveTCPAddr("tcp", local_addr)
	if err != nil {
		log.Fatalln("Error resolving address:", err)
	}

	// Start server listener socket.
	listener, err := net.ListenTCP("tcp", tcp_addr)
	if err != nil {
		log.Fatalln("Error starting TCP listener:", err)
	}

	// Safety cleanup in any exit path.
	defer listener.Close()
	log.Println("Listening on", local_addr)

	// done closes when the client goroutine exits.
	// This gives main a reliable shutdown trigger.
	done := make(chan struct{})

	// Launch the in-process client.
	go createConnection(tcp_addr, done)

	// When client finishes, close listener to unblock Accept().
	go func() {
		<-done
		_ = listener.Close()
	}()

	// Accept loop for incoming connections.
	for {
		conn, err := listener.Accept()
		if err != nil {
			// If client already ended, this accept error is expected.
			select {
			case <-done:
				log.Println("Shutting down server")
				return
			default:
			}

			// Listener closure can also show up as net.ErrClosed.
			if errors.Is(err, net.ErrClosed) {
				log.Println("Shutting down server")
				return
			}

			// Non-shutdown accept errors are logged and loop continues.
			log.Println("Error accepting connection:", err)
			continue
		}

		// Handle each accepted connection concurrently.
		log.Println("Accepted connection from", conn.RemoteAddr())
		go handleConnection(conn)
	}
}
