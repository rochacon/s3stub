// s3stub is a simple server that will write a local file
// on PUT and retrieve files with GET
package main

import (
	"crypto/sha256"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"hash"
	"io"
	"log"
	"net/http"
	"os"
	"path"
)

var BasePath string

type ReadHasher struct {
	h hash.Hash
	r io.Reader
}

func (rh *ReadHasher) Read(buf []byte) (n int, err error) {
	n, err = rh.r.Read(buf)
	if err != nil {
		return n, err
	}

	rh.h.Write(buf[:n])

	return
}

func (rh *ReadHasher) Sum(b []byte) []byte {
	return rh.h.Sum(b)
}

func main() {
	flag.StringVar(&BasePath, "r", "", "The root path of the server")
	flag.Parse()
	if BasePath == "" {
		fmt.Println("s3stub:")
		flag.PrintDefaults()
		return
	}

	r := mux.NewRouter()
	r.HandleFunc("/{path:.+}", download).Methods("GET")
	r.HandleFunc("/{path:.+}", upload).Methods("PUT")

	http.Handle("/", r)

	log.Println("Listening on :8000")
	log.Println("BasePath:", BasePath)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func download(w http.ResponseWriter, r *http.Request) {
	filename := path.Join(BasePath, r.URL.Path)

	fp, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			http.Error(w, err.Error(), 404)
		} else {
			http.Error(w, err.Error(), 500)
		}
		return
	}
	defer fp.Close()

	io.Copy(w, fp)
}

func upload(w http.ResponseWriter, r *http.Request) {
	filename := path.Join(BasePath, r.URL.Path)

	os.MkdirAll(path.Dir(filename), 0700)

	fp, err := os.Create(filename)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer fp.Close()

	bodyReader := &ReadHasher{h: sha256.New(), r: r.Body}
	io.Copy(fp, bodyReader)
	fmt.Fprintf(w, "%x", bodyReader.Sum(nil))
}
