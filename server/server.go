package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/lum8rjack/GoOut/server/modules/http"
	"github.com/lum8rjack/GoOut/server/modules/https"
	"github.com/lum8rjack/GoOut/server/modules/loadconfig"
	"github.com/lum8rjack/GoOut/server/modules/tcp"
	"github.com/lum8rjack/GoOut/server/modules/udp"
	"github.com/lum8rjack/GoOut/server/modules/writefile"
)

var (
	version = "1.0"
	wg      sync.WaitGroup
	config  loadconfig.Configuration
)

func startServers() {

	if config.TCP.Enabled {
		wg.Add(1)
		tcpc := tcp.NewTCP(config.Logging.LogFile, config.Logging.UploadDir, config.TCP.Port)
		go tcp.StartTCP(tcpc, &wg)
	}

	if config.UDP.Enabled {
		wg.Add(1)
		udpc := udp.NewUDP(config.Logging.LogFile, config.Logging.UploadDir, config.UDP.Port)
		go udp.StartUDP(udpc, &wg)
	}

	if config.HTTP.Enabled {
		wg.Add(1)
		httpc := http.NewHTTP(config.Logging.LogFile, config.Logging.UploadDir, config.HTTP.Port, config.HTTP.Get, config.HTTP.Post, config.HTTP.Uploadsize)
		go http.StartHTTP(httpc, &wg)
	}

	if config.HTTPS.Enabled {
		wg.Add(1)
		httpsc := https.NewHTTPS(config.Logging.LogFile, config.Logging.UploadDir, config.HTTPS.Port, config.HTTPS.Get, config.HTTPS.Post, config.HTTPS.Uploadsize, config.HTTPS.Certificate, config.HTTPS.Key)
		go https.StartHTTPS(httpsc, &wg)
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
		writefile.WriteLog(config.Logging.LogFile, "Stopped GoOutServer")
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
	config, err = loadconfig.LoadConf(configFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	time.Sleep(1 * time.Second)
	writefile.WriteLog(config.Logging.LogFile, "Starting GoOutServer")
	startServers()

	wg.Wait()
}
