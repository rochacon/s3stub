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
	"path/filepath"
	"strings"
)

var Root string

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
	var bind string

	flag.StringVar(&bind, "b", "127.0.0.1:8000", "The address to bind to")
	flag.StringVar(&Root, "r", "", "The root path of the server")
	flag.Parse()

	if Root == "" {
		fmt.Println("s3stub:")
		flag.PrintDefaults()
		return
	}

	r := mux.NewRouter()
	r.HandleFunc("/", list).Methods("GET")
	r.HandleFunc("/{path:.+}", download).Methods("GET")
	r.HandleFunc("/{path:.+}", exists).Methods("HEAD")
	r.HandleFunc("/{path:.+}", upload).Methods("PUT")
	r.HandleFunc("/{path:.+}", delete).Methods("DELETE")

	http.Handle("/", r)

	log.Println("Listening on:", bind)
	log.Println("Root:", Root)
	log.Fatal(http.ListenAndServe(bind, nil))
}

func download(w http.ResponseWriter, r *http.Request) {
	log.Println("GET", r.URL.Path)
	filename := path.Join(Root, r.URL.Path)

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

func exists(w http.ResponseWriter, r *http.Request) {
	log.Println("HEAD", r.URL.Path)
	filename := path.Join(Root, r.URL.Path)

	fp, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			w.WriteHeader(404)
		} else {
			log.Println("HEAD", r.URL.Path, err.Error())
			w.WriteHeader(500)
		}
		return
	}
	defer fp.Close()

	w.WriteHeader(204)
}

func list(w http.ResponseWriter, r *http.Request) {
	log.Println("GET", r.URL.Path)

	printpath := func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		fmt.Fprintf(w, "%s\n", strings.Replace(path, Root, "", 1))
		return nil
	}

	filepath.Walk(Root, printpath)
}

func upload(w http.ResponseWriter, r *http.Request) {
	log.Println("PUT", r.URL.Path)
	filename := path.Join(Root, r.URL.Path)

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

func delete(w http.ResponseWriter, r *http.Request) {
	log.Println("DELETE", r.URL.Path)
	filename := path.Join(Root, r.URL.Path)

	err := os.Remove(filename)
	if err != nil {
		if os.IsNotExist(err) {
			http.Error(w, err.Error(), 404)
		} else {
			fmt.Println("wat")
			http.Error(w, err.Error(), 500)
		}
		return
	}

	w.WriteHeader(204)
}
