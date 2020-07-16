package https

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"path"
	"strconv"
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

		filename := keys[0]

		fmt.Fprintf(w, "Filename = %s\n", filename)
		fmt.Fprintf(w, "Data = %s\n", data)

	} else if r.URL.Path == "/"+hs.post && r.Method == "GET" {
		http.ServeFile(w, r, "config/http/post.html")
	} else if r.URL.Path == "/upload" && r.Method == "POST" {
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

		hs.filename = handler.Filename
		//s := strconv.FormatInt(handler.Size, 10)

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

func StartHTTPS(httpc httpsConf) {

	hs = httpc
	time.Sleep(1 * time.Second)

	http.HandleFunc("/", serveSecurePage)

	writefile.WriteLog(hs.logFile, "Started HTTPS server on port "+strconv.Itoa(hs.port))

	if err := http.ListenAndServe(":"+strconv.Itoa(hs.port), nil); err != nil {
		return
	}
}
