package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math/rand"
	"net"
	"time"
)

var magicCookie = [4]byte{99, 130, 83, 99}

// DHCPv4 packet
type DHCPv4 struct {
	op      byte
	htype   byte
	hlen    byte
	hops    byte
	xid     TransactionID
	secs    [2]byte
	flags   [2]byte
	ciaddr  [4]byte
	yiaddr  [4]byte
	siaddr  [4]byte
	giaddr  [4]byte
	chaddr  net.HardwareAddr
	sname   [64]byte
	file    [128]byte
	options Options
}

// GenerateTransactionID make xid
func GenerateTransactionID() TransactionID {
	var xid TransactionID
	rand.Seed(time.Now().UnixNano())
	token := make([]byte, 4)
	rand.Read(token)
	// fmt.Println("Token:", token)
	copy(xid[:], token)
	return xid
}

// AddOptionsToBytes make DHCP optiones to bytes
func (p *DHCPv4) AddOptionsToBytes() []byte {
	var out []byte
	for _, opt := range p.options {
		// out = append(out, []byte{byte(code), byte(len(value))}...)
		out = append(out, opt.OptBytes()...)
	}
	return out
}

// ToBytes make DHCP packet to bytes
func (p *DHCPv4) ToBytes() []byte {
	out := make([]byte, DHCPMINLEN)
	out[0] = p.op
	out[1] = p.htype
	out[2] = p.hlen
	out[3] = p.hops
	copy(out[4:8], p.xid[:])
	copy(out[8:10], p.secs[:])
	copy(out[10:12], p.flags[:])
	copy(out[12:16], p.ciaddr[:])
	copy(out[16:20], p.yiaddr[:])
	copy(out[20:24], p.siaddr[:])
	copy(out[24:28], p.giaddr[:])
	copy(out[28:44], p.chaddr[:])
	copy(out[44:108], p.sname[:])
	copy(out[108:236], p.file[:])
	// copy(out[236:240], magicCookie[:])
	out = append(out, magicCookie[:]...)
	out = append(out, p.AddOptionsToBytes()...)
	out = append(out, 255)
	return out
}

// FromBytes let Bytes to type DHCPv4
func FromBytes(b []byte) (*DHCPv4, bool) {
	var p DHCPv4
	p.op = b[0]
	p.htype = b[1]
	p.hlen = b[2]
	p.hops = b[3]
	copy(p.xid[:], b[4:8])
	copy(p.secs[:], b[8:10])
	copy(p.flags[:], b[10:12])
	copy(p.ciaddr[:], b[12:16])
	copy(p.yiaddr[:], b[16:20])
	copy(p.siaddr[:], b[20:24])
	copy(p.giaddr[:], b[24:28])
	p.chaddr = b[28:34] //mac only need 6 bytes
	copy(p.sname[:], b[44:108])
	copy(p.file[:], b[108:236])
	// skip the magic code

	if opts, ok := makeoptions(b[240:]); ok {
		p.options = opts
		return &p, true
	}
	return &p, false
}

// makeoptions use Buffer to read
// maybe can only use UDPConn as buffer to read
// https://v2ex.com/t/727922#r_9814982
// https://juejin.cn/post/6844903721269198855
func makeoptions(b []byte) (Options, bool) {
	var options Options = make(Options)
	if b[0] == 0 {
		return options, false
	}

	buf := make([]byte, len(b))
	reader := bytes.NewBuffer(b)

	var l int
	for {
		l = 1
		// read 1
		n, err := io.ReadFull(reader, buf[:l])
		if n != l {
			fmt.Println("read OptionCode error:", err)
			return options, false
		}
		code := OptionCode(buf[0])

		if code == End {
			break
		}

		// read 1
		n, err = io.ReadFull(reader, buf[:l])
		if n != l {
			fmt.Println("read OptionLength error:", err)
			return options, false
		}
		length := int(buf[0])
		// read length

		n, err = io.ReadFull(reader, buf[:length])
		if n != length {
			fmt.Println("value error:", err)
			return options, false
		}
		value := make([]byte, length)
		copy(value, buf[:length])
		// getOption
		if int(code) == 53 {
			options[code] = MessageType(value[0])
		} else {
			options[code] = DHCPOption{code, value}
		}
	}
	return options, true
}

func getMacAddr(name string) net.HardwareAddr {
	ifa, _ := net.InterfaceByName(name)
	return ifa.HardwareAddr
}

// NewDiscovery used by client
func NewDiscovery() *DHCPv4 {
	xid := GenerateTransactionID()
	chaddr := []byte(getMacAddr("en0")) // just for test

	p := DHCPv4{
		op:      BOOTREQUEST,
		htype:   1,
		hlen:    6,
		hops:    0, // optionally used by relay agents
		xid:     xid,
		secs:    [2]byte{0, 0},
		flags:   [2]byte{0x80, 0}, // 0x80 broadcast 0 unicast
		chaddr:  chaddr,
		options: make(Options),
	}
	p.options[OptionDHCPMessageType] = MessageType(Discover)
	p.options[OptionHostName] = DHCPOption{OptionHostName, []byte(appname)}
	return &p
}

// NewRequest used by client
func NewRequest() *DHCPv4 {
	xid := GenerateTransactionID()
	chaddr := []byte(getMacAddr("en0")) // just for test

	p := DHCPv4{
		op:      BOOTREQUEST,
		htype:   1,
		hlen:    6,
		hops:    0, // optionally used by relay agents
		xid:     xid,
		secs:    [2]byte{0, 0},
		flags:   [2]byte{0x80, 0}, // 0x80 broadcast 0 unicast
		chaddr:  chaddr,
		options: make(Options),
	}
	p.options[OptionDHCPMessageType] = MessageType(Request) // ???
	// p.options[OptionRequestedIPAddress] = net.IP{10, 69, 33, 110}
	// p.options[OptionServerIdentifier] = net.IP{10, 69, 33, 1}
	p.options[OptionHostName] = DHCPOption{OptionHostName, []byte(appname)}
	return &p
}

// NewOffer is DHCP Offer
func NewOffer(discover *DHCPv4, ip net.IP) *DHCPv4 {
	yiaddr := [4]byte{}
	copy(yiaddr[:], ip[12:]) // IPv4 same length as IPv6

	second := make([]byte, 4)
	binary.BigEndian.PutUint32(second, defaultLeaseTime)

	p := DHCPv4{
		op:      BOOTREPLY, //1 = BOOTREQUEST, 2 = BOOTREPLY
		htype:   1,
		hlen:    6,
		xid:     discover.xid,
		flags:   [2]byte{0x80, 0}, // 0x80 broadcast 0 unicast
		yiaddr:  yiaddr,
		giaddr:  discover.giaddr,
		chaddr:  discover.chaddr,
		options: make(Options),
	}
	p.options[OptionDHCPMessageType] = MessageType(Offer)
	p.options[OptionSubnetMask] = DHCPOption{OptionSubnetMask, []byte(defaultSubnetMask)}
	p.options[OptionRouter] = DHCPOption{OptionRouter, []byte(defaultRouter)}
	p.options[OptionDomainNameServer] = DHCPOption{OptionDomainNameServer, []byte(defaultRouter)}
	p.options[OptionIPAddressLeaseTime] = DHCPOption{OptionIPAddressLeaseTime, second}
	p.options[OptionServerIdentifier] = DHCPOption{OptionServerIdentifier, []byte(defaultServerIdentifier)}
	p.options[OptionHostName] = DHCPOption{OptionHostName, []byte(appname)}
	return &p
}

// NewAck is DHCP Ack
func NewAck(request *DHCPv4, ip net.IP) *DHCPv4 {
	yiaddr := [4]byte{}
	copy(yiaddr[:], ip[12:])

	second := make([]byte, 4)
	binary.BigEndian.PutUint32(second, defaultLeaseTime)

	p := DHCPv4{
		op:      BOOTREPLY, //1 = BOOTREQUEST, 2 = BOOTREPLY
		htype:   1,
		hlen:    6,
		xid:     request.xid,
		flags:   [2]byte{0x80, 0}, // 0x80 broadcast 0 unicast
		ciaddr:  request.ciaddr,
		yiaddr:  yiaddr,
		giaddr:  request.giaddr,
		chaddr:  request.chaddr,
		options: make(Options),
	}
	p.options[OptionDHCPMessageType] = MessageType(ACK)
	p.options[OptionSubnetMask] = DHCPOption{OptionSubnetMask, []byte(defaultSubnetMask)}
	p.options[OptionRouter] = DHCPOption{OptionRouter, []byte(defaultRouter)}
	p.options[OptionDomainNameServer] = DHCPOption{OptionDomainNameServer, []byte(defaultRouter)}
	p.options[OptionIPAddressLeaseTime] = DHCPOption{OptionIPAddressLeaseTime, second}
	p.options[OptionServerIdentifier] = DHCPOption{OptionServerIdentifier, []byte(defaultServerIdentifier)}
	p.options[OptionHostName] = DHCPOption{OptionHostName, []byte(appname)}
	return &p
}
