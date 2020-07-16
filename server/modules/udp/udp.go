package udp

import (
	"fmt"
	"net"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/lum8rjack/GoOut/server/modules/writefile"
)

type udpConf struct {
	port     int
	fileDir  string
	logFile  string
	buffer   int
	fnLength int
	filename string
	remoteip string
}

func NewUDP(lfile string, odir string, port int) udpConf {
	var udp udpConf

	udp.buffer = 1024
	udp.fnLength = 0
	udp.filename = ""
	udp.port = port
	udp.logFile = lfile
	udp.fileDir = odir
	udp.remoteip = ""

	return udp
}

func serve(pc net.PacketConn, addr *net.UDPAddr, buf []byte, udp *udpConf) {
	//pc.WriteTo(buf, addr)

	if string(buf[:5]) == "START" {
		// START,filename.txt
		//t := strings.TrimSuffix(string(buf), "\n")
		l := strings.Split(string(buf), ",")
		t := l[1]
		udp.filename = t
		udp.fnLength = len(udp.filename)
		udp.remoteip = addr.IP.String()
	} else if string(buf[:4]) == "STOP" {
		// STOP
		writefile.WriteLog(udp.logFile, "UDP from "+udp.remoteip+" - Wrote to "+udp.filename)
		udp.filename = ""
		udp.fnLength = 0
		udp.remoteip = ""
	} else if udp.fnLength != 0 && addr.IP.String() == udp.remoteip {
		writefile.WriteFile(path.Join(udp.fileDir, udp.filename), buf)
	} else {
		// do nothing
	}
}

func StartUDP(udp udpConf) {

	time.Sleep(1 * time.Second)

	s, err := net.ResolveUDPAddr("udp4", ":"+strconv.Itoa(udp.port))
	if err != nil {
		fmt.Println(err)
		return
	}

	l, err := net.ListenUDP("udp4", s)
	if err != nil {
		writefile.WriteLog(udp.logFile, "Error trying to start UDP server on port "+strconv.Itoa(udp.port))
		return
	}
	defer l.Close()

	writefile.WriteLog(udp.logFile, "Started UDP server on port "+strconv.Itoa(udp.port))

	for {
		buf := make([]byte, udp.buffer)
		//n, addr, err := l.ReadFrom(buf)
		n, addr, err := l.ReadFromUDP(buf)
		if err != nil {
			continue
		}

		go serve(l, addr, buf[:n], &udp)
	}
}
