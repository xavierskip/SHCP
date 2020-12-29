package main

import (
	"fmt"
	"net"
)

// do like a dhcp client
func testdhcpClient(message []byte) *DHCPv4 {
	// https://colobu.com/2016/10/19/Go-UDP-Programming/
	srcAddr := &net.UDPAddr{IP: net.IPv4zero, Port: 68}
	dstAddr := &net.UDPAddr{IP: net.IPv4bcast, Port: 67}

	conn, err := net.ListenUDP("udp", srcAddr)
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()

	n, err := conn.WriteToUDP(message, dstAddr)
	if err != nil {
		fmt.Println(err)
	}
	data := make([]byte, 1024)

	n, raddr, err := conn.ReadFrom(data)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("client read %d\n%s\nfrom %v\n", n, data, raddr)
	// b := make([]byte, 1)
	// os.Stdin.Read(b)

	p, ok := FromBytes(data)
	if !ok {
		return p
	}
	fmt.Println("xid", p.xid, "yiaddr", p.yiaddr, "chaddr", p.chaddr)
	fmt.Println("MessageType", p.options[53])
	// for k, v := range p.options {
	// 	fmt.Println(k, v)
	// }
	return p
}

func testsend() {
	p := NewRequest()
	testdhcpClient(p.ToBytes())
	fmt.Println("send Request")

	// dhcp Discovery
	// discovery := NewDiscovery()
	// offer := dhcpClient(discovery.ToBytes())
	// fmt.Println("client send Discovery")

	// // dhcp request
	// req := NewRequest()
	// req.xid = offer.xid
	// req.chaddr = offer.chaddr
	// req.options[OptionRequestedIPAddress] = offer.yiaddr[:]
	// req.options[OptionServerIdentifier] = offer.options[OptionServerIdentifier]
	// udpSend(req.ToBytes())
}

// https://github.com/aler9/howto-udp-broadcast-golang
func testlisten() {
	pc, err := net.ListenPacket("udp4", ":67")
	if err != nil {
		panic(err)
	}
	defer pc.Close()

	buf := make([]byte, 1024)
	n, addr, err := pc.ReadFrom(buf)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s sent this: %s\n", addr, buf[:n])
}
