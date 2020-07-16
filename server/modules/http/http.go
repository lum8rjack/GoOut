package http

import (
	b64 "encoding/base64"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"path"
	"strconv"
	"time"

	"github.com/lum8rjack/GoOut/server/modules/writefile"
)

var hc httpConf

type httpConf struct {
	port       int
	get        string
	post       string
	fileDir    string
	logFile    string
	filename   string
	uploadsize int64
}

func NewHTTP(lfile string, odir string, port int, get string, post string, upsize int64) httpConf {
	var http httpConf
	http.port = port
	http.logFile = lfile
	http.fileDir = odir
	http.get = get
	http.post = post
	http.uploadsize = upsize

	return http
}

func servePage(w http.ResponseWriter, r *http.Request) {
	if (r.URL.Path == "/" || r.URL.Path == "/index.html") && r.Method == "GET" {
		http.ServeFile(w, r, "config/http/index.html")
	} else if r.URL.Path == "/"+hc.get && r.Method == "GET" {
		keys, ok := r.URL.Query()["d"]

		if !ok || len(keys[0]) < 1 {
			fmt.Fprintf(w, "Missing data\n")
			return
		}

		data := keys[0]

		keys, ok = r.URL.Query()["f"]

		if !ok || len(keys[0]) < 1 {
			fmt.Fprintf(w, "Missing filename\n")
			return
		}

		filename := keys[0]

		//fmt.Fprintf(w, "Filename = %s\n", filename)
		//fmt.Fprintf(w, "Data = %s\n", data)
		uDec, _ := b64.URLEncoding.DecodeString(data)
		writefile.WriteFile(path.Join(hc.fileDir, filename), uDec)
		//http.ServeFile(w, r, "config/http/index.html")

	} else if r.URL.Path == "/"+hc.post && r.Method == "GET" {
		http.ServeFile(w, r, "config/http/post.html")
	} else if r.URL.Path == "/upload" && r.Method == "POST" {
		// Parse our multipart form, 10 << 20 specifies a maximum
		// upload of 10 MB files.
		r.Body = http.MaxBytesReader(w, r.Body, hc.uploadsize<<20)
		r.ParseMultipartForm(hc.uploadsize << 20)

		file, handler, err := r.FormFile("filename")
		if err != nil {
			fmt.Fprintf(w, "File is too large to upload. Please try a different file.\n")
			return
		}
		defer file.Close()

		//s := strconv.FormatInt(handler.Size, 10)

		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Println(err)
		}
		// write this byte array to to Uploads directory
		writefile.WriteFile(path.Join(hc.fileDir, handler.Filename), fileBytes)

		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			return
		}

		writefile.WriteLog(hc.logFile, "HTTP from "+ip+" - Wrote to "+handler.Filename)
		fmt.Fprintf(w, "Successfully uploaded:  %s\n", handler.Filename)

	} else {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
}

func StartHTTP(httpc httpConf) {

	hc = httpc
	time.Sleep(1 * time.Second)

	http.HandleFunc("/", servePage)

	writefile.WriteLog(hc.logFile, "Started HTTP server on port "+strconv.Itoa(hc.port))

	if err := http.ListenAndServe(":"+strconv.Itoa(hc.port), nil); err != nil {
		return
	}
}
