package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"os"
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
				conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", target, port))
				if err != nil {
					fmt.Print(err)
				}

				fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
				status, err := bufio.NewReader(conn).ReadString('\n')
				if err != nil {
					fmt.Print(err)
				}
				fmt.Printf("worker-goroutine.%d status=%s", worker, status)
				time.Sleep(1000 * time.Millisecond)
			}
		}(i)
	}

	for {
	}
}
