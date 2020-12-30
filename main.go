package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

const (
	appname = "SDHCP"
	version = "0.1"
)

var configfile string

// Config is read from a json file
var Config map[string]interface{}

// http://www.networksorcery.com/enp/protocol/bootp/option001.htm
// var defaultSubnetMask net.IPMask = net.IPv4Mask(255, 255, 255, 0) // Mask 4 byte.
var defaultRouter net.IP = net.ParseIP("192.168.20.1")[12:] // net.ipv4 has IPv6len
var defaultSubnetMask net.IPMask = defaultRouter.DefaultMask()
var defaultDNS net.IP = defaultRouter
var defaultServerIdentifier net.IP = defaultRouter
var defaultLeaseTime uint32 = 60
var defaultDatabase string = "./network.db"
var defaultTable string = "192.168.20.1/24"

// use macaddr to get ip from sqlite database
func matchAddress(mac net.HardwareAddr) (net.IP, bool) {
	var ip string
	queryaddr := fmt.Sprintf("%s", mac)
	// db field `mac` like 60-ab-67-f8-d5-6c
	queryaddr = strings.ReplaceAll(queryaddr, ":", "-")

	// DB must exist!
	//https://stackoverflow.com/questions/12518876/how-to-check-if-a-file-exists-in-go/12518877
	if _, err := os.Stat(defaultDatabase); os.IsNotExist(err) {
		logErr(err)
	}

	db, err := sql.Open("sqlite3", defaultDatabase)
	logErr(err)
	defer db.Close()

	sql := fmt.Sprintf("SELECT ip FROM `%s` WHERE mac=?;", defaultTable)
	stmt, err := db.Prepare(sql)
	logErr(err)
	defer stmt.Close()

	err = stmt.QueryRow(queryaddr).Scan(&ip)
	if err != nil {
		return nil, false
	}
	return net.ParseIP(ip), true
}

// just send DHCP BOOTREPLY in broadcast
func udpSend(socket *net.UDPConn, b []byte) {
	dstaddr := &net.UDPAddr{IP: net.IPv4bcast, Port: 68}
	_, err := socket.WriteToUDP(b, dstaddr)
	logErr(err)
}

func listenServer() {
	port := 67
	addr := &net.UDPAddr{IP: nil, Port: port} // IP: net.IPv4zero
	socket, err := net.ListenUDP("udp", addr)
	logErr(err)
	log.Println("listening....", port)
	defer socket.Close()

	for {
		data := make([]byte, DHCPMAXLEN)
		read, remoteAddr, err := socket.ReadFromUDP(data)
		if err != nil {
			log.Println("ReadFromUDP:", err)
		}
		log.Printf("Read %d From %v ", read, remoteAddr)
		if read < 240 {
			continue // Packet too small to be DHCP
		}

		go process(socket, data[:read]) // main loop will not panic
	}
}

func process(socket *net.UDPConn, data []byte) {
	p, ok := FromBytes(data)
	if !ok {
		log.Println("incomplete packet!")
		return
	}
	messagetype, ok := p.options[OptionDHCPMessageType]
	if !ok {
		log.Println("NO DHCPMessageType!")
		return
	}
	log.Printf("xid: %s, %v,\n", p.xid, messagetype)

	ip, ok := matchAddress(p.chaddr)
	if ok {
		log.Printf("MATCH: %s %s\n", p.chaddr, ip)
	} else {
		log.Printf("NO FOUND: %s\n", p.chaddr)
		return
	}
	switch messagetype {
	case Discover:
		log.Println("Discover -> Offer")
		offer := NewOffer(p, ip)
		udpSend(socket, offer.ToBytes())
	case Request:
		log.Println("Request -> Ack")
		ack := NewAck(p, ip)
		udpSend(socket, ack.ToBytes())
	case Decline:
		log.Println("Some one already user thie IP")
	case Release:
		log.Println("Release IP")
	case Inform:
		log.Println("Just note it")
	default:
		log.Println("? unknow messagetype", messagetype)
	}
}

func main() {
	flag.StringVar(&configfile, "c", "config.json", "json config file")
	v := flag.Bool("v", false, "version")
	flag.Parse()

	if *v { // show version
		fmt.Println(appname, version)
		os.Exit(0)
	}
	// read json and set Config
	dat, err := ioutil.ReadFile(configfile)
	logErr(err)
	err = json.Unmarshal(dat, &Config)
	logErr(err)
	defaultRouter = net.ParseIP(Config["route"].(string))[12:]
	m := net.ParseIP(Config["mask"].(string))[12:]
	defaultSubnetMask = net.IPv4Mask(m[0], m[1], m[2], m[3])
	if val, ok := Config["dns"]; ok {
		defaultDNS = net.ParseIP(val.(string))[12:]
	}
	defaultServerIdentifier = defaultRouter
	defaultLeaseTime = uint32(Config["leasetime"].(float64))
	if val, ok := Config["database"]; ok {
		defaultDatabase = val.(string)
	}
	if val, ok := Config["table"]; ok {
		defaultTable = val.(string)
	}

	// testsend()

	// testlisten()

	listenServer()

}
