package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
)

const LargestCachedNumber = 1 << 17

// If n'th bit is off, n is a prime number.
// This array is 1-based in a sense that bitmap[1] shows 1 is prime or not
// It's a bitmap just for fun, I don't think it adds a lot of value/performance.
var bitmap [LargestCachedNumber/8 + 2]byte
var primes = make([]int32, 0, 10000)

func fillCache() {
	setNotPrime(0)
	setNotPrime(1)
	for i := 2; i <= LargestCachedNumber; i++ {
		if !isPrimeCached(i) {
			continue
		}
		primes = append(primes, int32(i))
		for j := 2 * i; j <= LargestCachedNumber; j += i {
			setNotPrime(j)
		}
	}
}

func setNotPrime(x int) {
	bitmap[x/8] |= 1 << (x % 8)
}

func isPrimeCached(x int) bool {
	return (bitmap[x/8]>>(x%8))&1 == 0
}

func isPrime(x int) bool {
	if x <= LargestCachedNumber {
		return isPrimeCached(x)
	}
	if x > LargestCachedNumber*LargestCachedNumber {
		panic(fmt.Sprintf("Too big for me bro! I only support up to %v", LargestCachedNumber*LargestCachedNumber))
	}
	// Iterate on prime numbers up to sqrt(x)
	for _, prime := range primes {
		prime := int(prime)
		if prime*prime > x {
			break
		}
		if x%prime == 0 {
			return false
		}
	}
	return true
}

func readInput() error {
	var err error

	// IO by far is the biggest bottleneck of this program. Golang doesn't buffer anything by default.
	const bufferSize = 1000000 // It seems Docker IO is bad!
	in := os.Stdin
	if len(os.Args) > 1 {
		in, err = os.Open(os.Args[1])
		if err != nil {
			return fmt.Errorf("opening file: %w", err)
		}
		defer in.Close()
	}

	out := os.Stdout
	if len(os.Args) > 2 {
		out, err = os.Create(os.Args[2])
		if err != nil {
			return fmt.Errorf("create file: %w", err)
		}
		defer out.Close()
	}

	// We use a hand written version of bufio.Writer to reduce overhead. Mainly memory copying.
	writeBuf := make([]byte, 0, bufferSize)
	defer func() {
		_, _ = out.Write(writeBuf)
	}()

	// Also a hand written version for reader. There was no big bottleneck in profiler, this reader
	// was the biggest.
	readBuf := make([]byte, 0, bufferSize)
	readCursor := 0

	for {
		// We can use scanner or ReadyByte here but they will eventually allocate a string
		// and become our bottleneck
		number := 0
		for {
			if readCursor >= len(readBuf) {
				n, err := in.Read(readBuf[:cap(readBuf)])
				if err == io.EOF {
					// n should be zero by definition of file.Read
					if n != 0 {
						panic("n != 0")
					}
					if number == 0 {
						return nil
					}
					break
				}
				if err != nil {
					return fmt.Errorf("read error: %w", err)
				}
				if n == 0 {
					// Theoretically this can happen but I think not possible on Unix,
					// correct implementation would be looping until we read something.
					panic("n == 0!")
				}
				readBuf = readBuf[:n]
				readCursor = 0
			}
			singleByte := readBuf[readCursor]
			readCursor++

			if singleByte == '\n' {
				break
			}
			// Budget string to int conversion
			number = number * 10 + int(singleByte) - '0'
		}

		if isPrime(number) {
			writeBuf = append(writeBuf, '1', '\n')
		} else {
			writeBuf = append(writeBuf, '0', '\n')
		}

		if len(writeBuf) + 2 >= bufferSize {
			_, _ = out.Write(writeBuf)
			writeBuf = writeBuf[:0]
		}
	}
}

var cpuprofile = flag.String("cpuprofile", "", "")

func main() {
	start := time.Now()

	flag.Parse()
	os.Args = append([]string{os.Args[0]}, flag.Args()...)

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}

		runtime.SetCPUProfileRate(5000)
		err = pprof.StartCPUProfile(f)
		if err != nil {
			log.Fatalf("profile failed: %v", err)
		}
		defer pprof.StopCPUProfile()
	}

	fillCache()
	duration := time.Since(start)
	log.Printf("%v primes found in %v", len(primes), duration)

	// Some tests
	for _, i := range []int{1, 2, 3, 4, 5, 6, 7, 2147483647} {
		log.Printf("%v %v", i, isPrime(i))
	}

	err := readInput()
	if err != nil {
		log.Fatalf("Fatal: %v", err)
	}
	duration = time.Since(start)
	log.Printf("Done %v", duration)
}
