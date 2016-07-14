package main

import (
	"expvar"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
)

var eth1 string
var counts = expvar.NewMap("counters")

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
		go func(po int) {
			log.Printf("%s:%d", eth1, po)
			http.ListenAndServe(fmt.Sprintf("%s:%d", eth1, po), &handler{})
		}(p)
	}

	mux := http.NewServeMux()
	mux.Handle("/hi", &handler{})
	mux.Handle("/debug/vars", http.DefaultServeMux)

	http.ListenAndServe(":"+port, mux)
}

type handler struct {
}

func (m *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	counts.Add("hits", 1)
	w.Write([]byte("Hi"))
}
