package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/lum8rjack/client/GoOut/modules"
)

func printLogo() {
	fmt.Println("GoOut v1.0")
}

func usage() {
	printLogo()
	filename := os.Args[0]
	fmt.Printf("Usage: %s [options] \n\n", filename)
	flag.PrintDefaults()
	fmt.Println("\nEach module is different but should be created similar to the following:")
	fmt.Println("\tHTTP\t[module,method,host]")
	fmt.Println("\t\tEx. http,get,http://example.com:8080/status")
	fmt.Println("\t\tEx. http,post,http://192.168.1.1/upload")
	fmt.Println("\n\tSocket\t[module,method,host:port]")
	fmt.Println("\t\tEx. socket,tcp,192.168.1.1:8080")
	fmt.Println("\t\tEx. socket,udp,192.168.1.1:8888")
}

func isValidFile(file string) bool {

	info, err := os.Stat(file)

	if os.IsNotExist(err) {
		fmt.Println("File does not exist")
		return false
	} else if info.IsDir() {
		fmt.Println("Do not provie a directory")
		return false
	} else {
		return true
	}
}

func main() {
	flag.Usage = usage

	module := flag.String("module", "", "Specify module [HTTP|ICMP|SOCKET]")
	file := flag.String("file", "", "File you want to send")
	flag.Parse()

	if *module == "" || *file == "" {
		flag.Usage()
		return
	}

	*module = strings.ToLower(*module)

	s := strings.Split(*module, ",")

	if len(s) < 2 {
		flag.Usage()
		return
	}

	mod := s[0]
	data := s[1:]

	// Make sure the file is a valid file
	if !isValidFile(*file) {
		return
	}

	// Make sure a valid module is specified
	if mod == "http" {
		modules.HttpRun(data, *file)
	} else if mod == "icmp" {
		fmt.Println("Exfil via ICMP")
		modules.IcmpRun(data, *file)
	} else if mod == "socket" {
		modules.SocketRun(data, *file)
	} else {
		fmt.Println("Invalid module specified")
		flag.Usage()
		return
	}
}
