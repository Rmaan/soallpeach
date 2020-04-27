package main

import (
	"fmt"
	"github.com/valyala/fasthttp"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"sync/atomic"
	"time"
)

// Use atomic to update this. It's intentionally a pointer so you don't mistakenly think it's a
// normal integer.
var counter = new(int64)

func addHandler(ctx *fasthttp.RequestCtx) {
	if len(ctx.Path()) > 1 {
		countHandler(ctx)
		return
	}

	buf := ctx.Request.Body()
	n := len(buf)
	if n > 0 && buf[n-1] == '\n' {
		n--
	}

	// Simple []byte to int
	count := 0
	for i := 0; i < n; i++ {
		count = count*10 + int(buf[i]-'0')
	}

	atomic.AddInt64(counter, int64(count))
}

func countHandler(ctx *fasthttp.RequestCtx) {
	total := atomic.LoadInt64(counter)

	_, err := fmt.Fprintf(ctx, "%d", total)
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
	if false {
		f, err := os.Create("cpu-profile")
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close()
		runtime.SetCPUProfileRate(500)
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}

		go func() {
			time.Sleep(time.Second * 15)
			pprof.StopCPUProfile()
			f.Close()
		}()
	}

	port := envOrDefault("COUNTME_PORT", "8000")

	if err := fasthttp.ListenAndServe(":"+port, addHandler); err != nil {
		log.Fatalf("Error in ListenAndServe: %s", err)
	}
}
