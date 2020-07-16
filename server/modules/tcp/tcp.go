package tcp

import (
	"fmt"
	"net"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/lum8rjack/GoOut/server/modules/writefile"
)

type tcpConf struct {
	port     int
	fileDir  string
	logFile  string
	buffer   int
	fnLength int
	filename string
}

func NewTCP(lfile string, odir string, port int) tcpConf {
	var tcp tcpConf

	tcp.buffer = 1024
	tcp.fnLength = 0
	tcp.filename = ""
	tcp.port = port
	tcp.logFile = lfile
	tcp.fileDir = odir

	return tcp
}

func handleConnection(c net.Conn, tcp *tcpConf) {
	tcp.filename = ""
	tcp.fnLength = 0

	start := make([]byte, 64)
	_, err := c.Read(start)
	if err != nil {
		fmt.Println(err) //Prints EOF
		return
	}

	if string(start[:5]) == "START" {
		// START,filename.txt,N
		t := strings.Split(string(start), "\n")
		tcp.filename = strings.Split(t[0], ",")[1]
		tcp.fnLength = len(tcp.filename)

		buf := make([]byte, tcp.buffer)
		rm := c.RemoteAddr().String()

		writefile.WriteLog(tcp.logFile, "TCP from "+strings.Split(rm, ":")[0]+" - Wrote to "+tcp.filename)
		for {
			n, err := c.Read(buf)
			if err != nil {
				//fmt.Println(err) //Prints EOF
				return
			}

			if tcp.fnLength != 0 {
				writefile.WriteFile(path.Join(tcp.fileDir, tcp.filename), buf[:n])
			}

		}
	} else {
		c.Close()
	}
}

func StartTCP(tcp tcpConf) {

	time.Sleep(1 * time.Second)

	l, err := net.Listen("tcp4", ":"+strconv.Itoa(tcp.port))
	if err != nil {
		writefile.WriteLog(tcp.logFile, "Error trying to start TCP server on port "+strconv.Itoa(tcp.port))
		return
	}
	defer l.Close()

	writefile.WriteLog(tcp.logFile, "Started TCP server on port "+strconv.Itoa(tcp.port))

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		defer c.Close()
		go handleConnection(c, &tcp)

	}
}
