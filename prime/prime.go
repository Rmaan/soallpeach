package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
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
	in := os.Stdin
	if len(os.Args) > 1 {
		in, err = os.Open(os.Args[1])
		if err != nil {
			return fmt.Errorf("opening file: %w", err)
		}
		defer in.Close()
	}
	reader := bufio.NewReader(in)

	out := os.Stdout
	if len(os.Args) > 2 {
		out, err = os.Create(os.Args[2])
		if err != nil {
			return fmt.Errorf("create file: %w", err)
		}
		defer out.Close()
	}
	writer := bufio.NewWriter(out)
	defer writer.Flush()

	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			if line == "" {
				break
			}
			// Seems input doesn't have a new line at end of file 🤷‍♂️
			// Simulate an extra \n
			err = nil
			line = line + "\n"
		}
		if err != nil {
			return fmt.Errorf("read string: %w", err)
		}
		line = line[:len(line)-1]

		number, err := strconv.Atoi(line)
		if err != nil {
			return fmt.Errorf("non-integer value provided: %v, %w", line, err)
		}

		result := []byte("0\n")
		if isPrime(number) {
			result = []byte("1\n")
		}

		_, err = writer.Write(result)
		if err != nil {
			return fmt.Errorf("fprintln: %w", err)
		}
	}
	return nil
}

func main() {
	start := time.Now()
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
