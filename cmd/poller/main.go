package main

import (
	"fmt"
	"io/ioutil"
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
				fmt.Printf("worker-goroutine.%d port=%d", worker, port)

				fmt.Printf("worker-goroutine.%d at=Dialing", worker)
				conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", target, port), 3*time.Second)
				if err != nil {
					fmt.Print(err)
				}
				fmt.Printf("worker-goroutine.%d at=FinishDialing", worker)

				fmt.Printf("worker-goroutine.%d at=GET", worker)
				fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
				all, err := ioutil.ReadAll(conn)
				fmt.Printf("worker-goroutine.%d at=FinishGET", worker)
				if err != nil {
					fmt.Print(err)
				}
				fmt.Printf("worker-goroutine.%d target=%s port=%d status=%s", worker, target, port, all)
				fmt.Printf("worker-goroutine.%d at=Sleeping", worker)
				conn.Close()
				time.Sleep(2000 * time.Millisecond)
				fmt.Printf("worker-goroutine.%d at=FinishSleeping", worker)
			}
		}(i)
	}

	for {
	}
}
