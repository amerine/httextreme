package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"os"
	"runtime"
	"strconv"
	"time"
)

func main() {
	rand.Seed(time.Now().Unix())

	minPort, err := strconv.Atoi(os.Getenv("MIN_PORT"))
	if err != nil {
		panic(err)
	}

	maxPort, err := strconv.Atoi(os.Getenv("MAX_PORT"))
	if err != nil {
		panic(err)
	}
	concurrency, err := strconv.Atoi(os.Getenv("CONCURRENCY"))

	target := os.Getenv("TARGET")

	for i := 0; i < concurrency; i++ {
		go func(worker int) {
			for {
				port := minPort
				if (maxPort - minPort) > 1 {
					port = rand.Intn(maxPort-minPort) + minPort
				}

				conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", target, port), 3*time.Second)
				if err != nil {
					fmt.Println("Dial Failed", err, target, port)
				}

				fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
				_, err = ioutil.ReadAll(conn)
				if err != nil {
					fmt.Println("Read All Failed", err, target, port)
				}

				// conn.Close()
				time.Sleep(2000 * time.Millisecond)
			}
		}(i)
	}

	for {
		fmt.Printf("Number of goroutines: %d", runtime.NumGoroutine())
		time.Sleep(time.Duration((rand.Int31n(1000) + 1)) * time.Millisecond)
	}
}
