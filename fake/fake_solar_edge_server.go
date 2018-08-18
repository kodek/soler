package fake

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"fmt"

	"github.com/gorilla/mux"
)

func NewServer() *httptest.Server {
	r := mux.NewRouter()
	r.HandleFunc("/site/{site}/{method:[a-zA-z]+\\.json}", handleJson)
	ts := httptest.NewServer(r)
	return ts
}

func handleJson(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["method"]
	bytes, err := ioutil.ReadFile("fake/" + name)
	if err != nil {
		panic(err)
	}
	fmt.Fprint(w, string(bytes))
}
