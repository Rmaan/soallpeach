package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
)

const bufferLength = 32

// Use atomic to update this. It's intentionally a pointer so you don't mistakenly think it's a
// normal integer.
var counter = new(int64)

var bufPool = sync.Pool{New: newBuf}
type bufType [bufferLength]byte

func newBuf() interface{} {
	return bufType{}
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	buf := bufPool.Get().(bufType)
	defer bufPool.Put(buf)

	n, err := r.Body.Read(buf[:])
	if err != io.EOF {
		log.Printf("non EOF from body read: %v", err)
		w.WriteHeader(500)
		return
	}

	if n > 0 && buf[n-1] == '\n' {
		n--
	}
	count, err := strconv.Atoi(string(buf[:n]))
	if err != nil {
		log.Printf("Error parsing int: %v", err)
		http.Error(w, "Not an int", 400)
		return
	}

	atomic.AddInt64(counter, int64(count))
}

func countHandler(w http.ResponseWriter, r *http.Request) {
	total := atomic.LoadInt64(counter)

	_, err := fmt.Fprintf(w, "%d", total)
	if err != nil {
		log.Printf("err in write %v", err)
		return
	}
}

func envOrDefault(key, def string) string {
	val := os.Getenv(key)
	if val == "" {
		return def
	}
	return val
}

func main() {
	port := envOrDefault("COUNTME_PORT", "8000")

	http.HandleFunc("/count", countHandler)
	http.HandleFunc("/", addHandler)
	log.Fatal(http.ListenAndServe(":" + port, nil))
}
