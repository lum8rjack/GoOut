package socket

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"
)

const SOCKETBUFFER = 1024

// Buffer is smaller because once 768 bytes are b64 encoded it will be 1024 bytes
//const SOCKETB64BUFFER = 768

func fillString(retunString string, toLength int) string {
	for {
		lengtString := len(retunString)
		if lengtString < toLength {
			retunString = retunString + ":"
			continue
		}
		break
	}
	return retunString
}

func SocketRun(data []string, file string) {

	if len(data) != 2 {
		flag.Usage()
		return
	}

	method := data[0]
	connect := data[1]

	if method == "" || connect == "" || file == "" {
		flag.Usage()
		return
	}

	method = strings.ToLower(method)

	if method != "tcp" && method != "udp" {
		fmt.Println("Must specify TCP or UDP for the method.")
		return
	}

	// Converting host and port to host:port
	//CONNECT := net.JoinHostPort(host, port)

	// Get information on the file
	info, err := os.Stat(file)
	if os.IsNotExist(err) {
		fmt.Println("File does not exist")
		return
	}

	FILESIZE := info.Size()

	// Make sure the user didnt provide a directory
	if info.IsDir() {
		fmt.Println("Do not provie a directory")
		return
	}

	// Connect to the remote host based on the method (tcp/udp) provided
	c, err := net.Dial(method, connect)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer c.Close()

	f, err := os.Open(file)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	fmt.Printf("Sending file: %v (%d bytes)\n", file, FILESIZE)

	sendBuffer := make([]byte, SOCKETBUFFER)
	start := []byte(fillString("START,"+file+",N", 64))
	stop := []byte("STOP," + file + "\n")
	s := 0
	if method == "udp" {
		// Milliseconds
		s = 100
	}
	// Send initial data about the file
	c.Write(start)

	for {
		n, err := f.Read(sendBuffer)
		if err != nil && err != io.EOF {
			break
		}
		if n == 0 {
			break
		}

		time.Sleep(time.Duration(s) * time.Millisecond)
		c.Write(sendBuffer[:n])

	}
	// Send final data to tell server the file is done if using UDP
	if method == "udp" {
		time.Sleep(time.Duration(s) * time.Millisecond)
		c.Write(stop)
	}
	fmt.Println("File sent!")
}
