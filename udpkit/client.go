package main

import (
	"fmt"
	"net"
)

func main() {
	// message := os.Args[1]
	ip := net.ParseIP("127.0.0.1")
	srcAddr := &net.UDPAddr{IP: net.IPv4zero, Port: 68}
	dstAddr := &net.UDPAddr{IP: ip, Port: 67}
	// conn, err := net.ListenUDP("udp", srcAddr)
	conn, err := net.DialUDP("udp", srcAddr, dstAddr)
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()
	message := []byte{
		0x01, 0x01, 0x06, 0x00, 0xc8, 0xee, 0xba, 0x66, 0x00, 0x00, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x4c, 0x56, 0x9d, 0x57,
		0xc1, 0xec, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x63, 0x82, 0x53, 0x63, // magiccode
		// 0x45, 0x10, 0xbb, 0x41, 0x10, 0x02, 0x03, 0x40, 0x55, 0x60, 0x70, 0x90, 0x63, 0x82, 0x53, 0x63, // magiccode
		0x35, 0x01, 0x05, 0x1e, 0x01, 0xa1,
		0x35, 0x01, 0x06, 0x1e, 0x01, 0x21,
		0x35, 0x01, 0x07, 0x1e, 0x01, 0x31,
		0xff,
	}

	conn.Write(message)
	return

	// data := make([]byte, 124)
	// n, _, err := conn.ReadFrom(data)
	// fmt.Printf("read %s from <%s>\n", data[:n], conn.RemoteAddr())

	// n, err := conn.WriteToUDP([]byte("hello"), dstAddr)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// data := make([]byte, 1024)
	// n, _, err = conn.ReadFrom(data)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Printf("read %s from <%s>\n", data[:n], conn.RemoteAddr())
	// b := make([]byte, 1)
	// os.Stdin.Read(b)
}