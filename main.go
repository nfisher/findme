package main

import (
	"flag"
	"fmt"
	"html"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"regexp"
	"strconv"

	"github.com/hashicorp/consul/api"
)

func main() {
	var listenAddr string

	flag.StringVar(&listenAddr, "listen", "127.0.0.1:0", "Listening address. Port 0 will delegate port selection to kernel.")
	flag.Parse()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)

	conn, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatal(err)
	}

	addr := conn.Addr().String()
	var ipMatch = regexp.MustCompile("^[^:]*")
	var portMatch = regexp.MustCompile("[^:]*$")
	ip := ipMatch.FindString(addr)
	port, _ := strconv.Atoi(portMatch.FindString(addr))

	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for _ = range c {
			cleanup(client, ip, port)
			os.Exit(0)
		}
	}()

	registerCheck(client, port)
	register(client, ip, port)
	go check(client, port)

	server := &http.Server{Handler: handler{}}
	log.Printf("Bound to %v.", addr)
	err = server.Serve(conn)

	log.Fatal(err)
}

type handler struct{}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}
