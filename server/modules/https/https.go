package https

import (
	b64 "encoding/base64"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"path"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/lum8rjack/GoOut/server/modules/writefile"
)

var hs httpsConf

type httpsConf struct {
	port       int
	get        string
	post       string
	fileDir    string
	logFile    string
	filename   string
	uploadsize int64
	cert       string
	key        string
}

func NewHTTPS(lfile string, odir string, port int, get string, post string, upsize int64, cert string, key string) httpsConf {
	var https httpsConf
	https.filename = ""
	https.port = port
	https.logFile = lfile
	https.fileDir = odir
	https.get = get
	https.post = post
	https.uploadsize = upsize
	https.cert = cert
	https.key = key

	return https
}

func serveSecurePage(w http.ResponseWriter, r *http.Request) {

	if (r.URL.Path == "/" || r.URL.Path == "/index.html") && r.Method == "GET" {
		http.ServeFile(w, r, "config/http/index.html")
	} else if r.URL.Path == "/"+hs.get && r.Method == "GET" {
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

		filename := filepath.Base(keys[0])

		uDec, _ := b64.URLEncoding.DecodeString(data)
		writefile.WriteFile(path.Join(hs.fileDir, filename), uDec)

	} else if r.URL.Path == "/"+hs.post && r.Method == "GET" {
		http.ServeFile(w, r, "config/http/post.html")
	} else if (r.URL.Path == "/upload" || r.URL.Path == "/"+hs.post) && r.Method == "POST" {
		// Parse our multipart form, 10 << 20 specifies a maximum
		// upload of 10 MB files.
		r.Body = http.MaxBytesReader(w, r.Body, hs.uploadsize<<20)
		r.ParseMultipartForm(hs.uploadsize << 20)

		file, handler, err := r.FormFile("filename")
		if err != nil {
			fmt.Fprintf(w, "File is too large to upload. Please try a different file.\n")
			return
		}
		defer file.Close()

		hs.filename = filepath.Base(handler.Filename)

		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Println(err)
		}
		// write this byte array to to Uploads directory
		writefile.WriteFile(path.Join(hs.fileDir, hs.filename), fileBytes)

		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			return
		}

		writefile.WriteLog(hs.logFile, "HTTPS from "+ip+" - Wrote to "+hs.filename)
		fmt.Fprintf(w, "Successfully uploaded:  %s\n", hs.filename)

	} else {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
}

func StartHTTPS(httpc httpsConf, wg *sync.WaitGroup) {

	defer wg.Done()

	hs = httpc

	if !writefile.IsValidFile(hs.key) || !writefile.IsValidFile(hs.cert) {
		writefile.WriteLog(hs.logFile, "HTTPS - Error: Certificate or key is not valid")
		return
	}
	time.Sleep(1 * time.Second)

	http.HandleFunc("/", serveSecurePage)

	writefile.WriteLog(hs.logFile, "Started HTTPS server on port "+strconv.Itoa(hs.port))

	if err := http.ListenAndServeTLS(":"+strconv.Itoa(hs.port), hs.cert, hs.key, nil); err != nil {
		return
	}
}
