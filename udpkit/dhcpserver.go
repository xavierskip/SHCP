package main

import (
	"fmt"
	"net"
	"sync"
)

func main() {
	addr1, err := net.ResolveUDPAddr("udp", ":11000")
	panicIf(err)
	addr2, err := net.ResolveUDPAddr("udp", ":11001")
	panicIf(err)

	srv1 := server(addr1)
	srv2 := server(addr2)

	var wg sync.WaitGroup
	wg.Add(2)

	go read(srv1, &wg)
	go read(srv2, &wg)

	write(srv1, addr2, []byte("server1 -> server2"))
	write(srv2, addr1, []byte("server2 -> server1"))

	wg.Wait()

	fmt.Println("everything is ok")
}

func server(addr *net.UDPAddr) *net.UDPConn {
	conn, err := net.ListenUDP("udp", addr)
	panicIf(err)

	return conn
}

func read(conn *net.UDPConn, wg *sync.WaitGroup) {
	var buf [1024]byte
	_, _, err := conn.ReadFromUDP(buf[:])
	panicIf(err)

	wg.Done()
}

func write(conn *net.UDPConn, addr *net.UDPAddr, message []byte) {
	_, err := conn.WriteToUDP(message, addr)
	panicIf(err)
}

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}
