package main

import (
	"net"
	"runtime"
)

var listenChannel chan net.Conn = make(chan net.Conn, 10)

func StartCrossDomain() error {
	var err error
	ln, err := net.Listen("tcp", ":8430")
	if err != nil {
		return err
	}
	var conn net.Conn
	for {
		conn, err = ln.Accept()
		if err != nil {
			continue
		}
		go handleClient(conn)
		runtime.Gosched()
	}
	return nil
}

func handleClient(conn net.Conn) {
	defer conn.Close()
	for {
		runtime.Gosched()
		Buffer := make([]byte, 200)
		numbytes, err := conn.Read(Buffer)
		if err != nil {
			runtime.Goexit()
		}
		if numbytes == 0 {
			runtime.Gosched()
			continue
		}
		conn.Write([]byte("<?xml version=\"1.0\"?><cross-domain-policy><site-control permitted-cross-domain-policies=\"all\"/><allow-access-from domain=\"*\" secure=\"true\" to-ports=\"*\"/></cross-domain-policy>"))
		runtime.Goexit()
	}
}

func main() {
	var err = StartCrossDomain()
	if err != nil {
		panic(err)
	}
}
