package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/lum8rjack/GoOut/server/modules"
)

var (
	version = "0.1"
	wg      sync.WaitGroup
	config  modules.Configuration
)

func startServers() {

	if config.TCP.Enabled {
		wg.Add(1)
		tcp := modules.NewTCP(config.Logging.LogFile, config.Logging.UploadDir, config.TCP.Port)
		go modules.StartTCP(tcp)
	}

	if config.UDP.Enabled {
		wg.Add(1)
		udp := modules.NewUDP(config.Logging.LogFile, config.Logging.UploadDir, config.UDP.Port)
		go modules.StartUDP(udp)
	}

	if config.HTTP.Enabled {
		wg.Add(1)
		http := modules.NewHTTP(config.Logging.LogFile, config.Logging.UploadDir, config.HTTP.Port, config.HTTP.Get, config.HTTP.Post, config.HTTP.Uploadsize)
		go modules.StartHTTP(http)
	}

	if config.HTTPS.Enabled {
		wg.Add(1)
		https := modules.NewHTTPS(config.Logging.LogFile, config.Logging.UploadDir, config.HTTPS.Port, config.HTTPS.Get, config.HTTPS.Post, config.HTTPS.Uploadsize, config.HTTPS.Certificate, config.HTTPS.Key)
		go modules.StartHTTPS(https)
	}

}

func printLogo() {
	fmt.Println("GoOutServer v" + version)
}

func usage() {
	printLogo()
	filename := os.Args[0]
	fmt.Printf("Usage: %s [options] \n\n", filename)
	flag.PrintDefaults()
	os.Exit(0)
}

func SetupCloseHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		modules.WriteLog(config.Logging.LogFile, "Stopped GoOutServer")
		os.Exit(0)
	}()
}

func main() {

	SetupCloseHandler()

	flag.Usage = usage

	configFile := flag.String("c", "config/config.json", "Configuration file to use")
	flag.Parse()

	if *configFile == "" {
		flag.Usage()
		return
	}

	var err error
	config, err = modules.LoadConf(configFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	time.Sleep(1 * time.Second)
	modules.WriteLog(config.Logging.LogFile, "Starting GoOutServer")
	startServers()

	wg.Wait()
}
