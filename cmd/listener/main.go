package main

import (
	_ "expvar"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
)

var eth1 string

func main() {
	port := os.Getenv("PORT")
	minPort, _ := strconv.Atoi(os.Getenv("MIN_PORT"))
	maxPort, _ := strconv.Atoi(os.Getenv("MAX_PORT"))

	iface, err := net.InterfaceByName("eth1")
	if err != nil {
		panic(err)
	}

	addrs, err := iface.Addrs()
	if err != nil {
		panic(err)
	}

	for _, addr := range addrs {
		switch ip := addr.(type) {
		case *net.IPNet:
			if ip.IP.DefaultMask() != nil {
				if ip.IP.To4() != nil {
					eth1 = (ip.IP).String()
				}
			}
		}
	}

	for p := minPort; p <= maxPort; p++ {
		go func(port int) {
			log.Printf("%s:%d", eth1, port)
			http.ListenAndServe(fmt.Sprintf("%s:%d", eth1, port), &handler{})
		}(p)
	}

	http.ListenAndServe(":"+port, &handler{})
}

type handler struct {
}

func (m *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hi"))
}
