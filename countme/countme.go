package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync/atomic"
)

const bufferLength = 32

// Use atomic to update this. It's intentionally a pointer so you don't mistakenly think it's a
// normal integer.
var counter = new(int64)

func addHandler(w http.ResponseWriter, r *http.Request) {
	buf := [bufferLength]byte{}
	n, err := r.Body.Read(buf[:])
	if err != io.EOF {
		log.Printf("non EOF from body read: %v", err)
	}

	count, err := strconv.Atoi(string(buf[:n]))
	if err != nil {
		http.Error(w, "not an int", 400)
		return
	}

	atomic.AddInt64(counter, int64(count))
}

func countHandler(w http.ResponseWriter, r *http.Request) {
	total := atomic.LoadInt64(counter)

	s := strconv.FormatInt(total, 10)
	_, err := w.Write([]byte(s))
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
