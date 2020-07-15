package modules

// Reference
// https://medium.com/hackervalleystudio/hacking-with-go-packet-crafting-and-manipulation-in-golang-pt-1-f31cdb066e3a

import (
	b64 "encoding/base64"
	"fmt"
	"io"
	"os"
	"time"
	//"github.com/google/gopacket"
	//"github.com/google/gopacket/layers"
	//"github.com/google/gopacket/pcap"
)

const ICMPBUFFER = 100

func sendPacket(data string) {
	fmt.Println(data)

	//rawBytes := []byte(data)

	//var buffer gopacket.SerializeBuffer
	/*buffer := gopacket.NewSerializeBuffer()
	var options gopacket.SerializeOptions

	gopacket.SerializeLayers(buffer, options,
		&layers.Ethernet{},
		&layers.IPv4{},
		&layers.TCP{},
		gopacket.Payload(rawBytes),
	)
	outgoingPacket := buffer.Bytes()*/
}

func sendICMP(url string, path string) {

	f, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	sendBuffer := make([]byte, ICMPBUFFER)
	rsent := 0

	for {
		n, err := f.Read(sendBuffer)
		if err != nil && err != io.EOF {
			break
		}
		if n == 0 {
			break
		}
		sEnc := b64.URLEncoding.EncodeToString(sendBuffer[:n])
		sendPacket(sEnc)
		rsent++
		//time.Sleep(2 * time.Second)
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Printf("Number of ICMP packets sent: %v\n", rsent)

}

func IcmpRun(data []string, file string) {

	host := data[0]

	if host == "" || file == "" {
		return
	}

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

	fmt.Printf("Sending file: %v (%d bytes) to %v\n", file, FILESIZE, host)
	sendICMP(host, file)

}
