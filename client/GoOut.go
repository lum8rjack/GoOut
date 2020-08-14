package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/lum8rjack/GoOut/client/modules/http"
	"github.com/lum8rjack/GoOut/client/modules/icmp"
	"github.com/lum8rjack/GoOut/client/modules/socket"
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
	fmt.Println("\tWEB\tUploads file via http(s) POST or multiple GET requests [module,method,host]")
	fmt.Println("\t\tEx. web,get,http://192.168.1.1:8080/status")
	fmt.Println("\t\tEx. web,post,https://192.168.1.1/upload")
	fmt.Println("\n\tSocket\tUploads file via raw TCP or UDP packets [module,method,host:port]")
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

	module := flag.String("module", "", "Specify module [WEB|ICMP|SOCKET]")
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
	if mod == "web" {
		http.HttpRun(data, *file)
	} else if mod == "icmp" {
		fmt.Println("Exfil via ICMP")
		icmp.IcmpRun(data, *file)
	} else if mod == "socket" {
		socket.SocketRun(data, *file)
	} else {
		fmt.Println("Invalid module specified")
		flag.Usage()
		return
	}
}
