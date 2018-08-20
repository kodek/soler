package soler

import (
	"crypto/md5"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strconv"
	"time"

	"io/ioutil"

	"bitbucket.org/kodek64/soler/greenbutton"
	"github.com/golang/glog"
)

const uploadTemplate = `
<html>
<head>
       <title>Upload file</title>
</head>
<body>
<form enctype="multipart/form-data" action="/upload" method="post">
    <input type="file" name="uploadfile" />
    <input type="hidden" name="token" value="{{.}}"/>
    <input type="submit" value="upload" />
</form>
</body>
</html>`

func (h *GreenButtonHandler) handleGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		panic("Should be a GET request")
	}
	crutime := time.Now().Unix()
	hash := md5.New()
	io.WriteString(hash, strconv.FormatInt(crutime, 10))
	token := fmt.Sprintf("%x", hash.Sum(nil))

	t, err := template.New("upload").Parse(uploadTemplate)
	if err != nil {
		writeError(w, err)
		return
	}
	t.Execute(w, token)
}

func writeError(w http.ResponseWriter, err error) {
	glog.Warning("Upload handler error: ", err)
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprint(w, err.Error())
}

func (h *GreenButtonHandler) handlePost(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("uploadfile")
	if err != nil {
		writeError(w, err)
		return
	}
	defer file.Close()

	glog.Info("Received file with header: ", handler.Header)

	uploadedBytes, err := ioutil.ReadAll(file)
	if err != nil {
		writeError(w, err)
		return
	}

	err = h.processUploadedBytes(uploadedBytes)

	if err != nil {
		writeError(w, err)
		return
	}
	fmt.Fprint(w, "Done!")
}

type GreenButtonHandler struct {
	Db *Database
}

func (h *GreenButtonHandler) processUploadedBytes(b []byte) error {
	dataPoints, err := greenbutton.Read(string(b))
	if err != nil {
		return err
	}
	return h.Db.AddConsumptionPoints(dataPoints)
}

// https://astaxie.gitbooks.io/build-web-application-with-golang/en/04.5.html
func (h *GreenButtonHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	glog.Info("Got green button upload request")
	if r.Method == "GET" {
		h.handleGet(w, r)
	} else {
		h.handlePost(w, r)
	}
}
