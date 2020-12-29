package main

import (
	"fmt"
)

// profit from https://github.com/insomniacslk/dhcp/blob/master/dhcpv4/types.go
// values from http://www.networksorcery.com/enp/protocol/dhcp.htm and
// http://www.networksorcery.com/enp/protocol/bootp/options.htm

// TransactionID a random number chosen by the
// client, used by the client and server to associate
type TransactionID [4]byte

// String prints a hex transaction ID.
func (xid TransactionID) String() string {
	return fmt.Sprintf("0x%x", xid[:])
}

// OptionCode just DHCP Option Code
type OptionCode byte

// OptionValue just DHCP Option Value
type OptionValue []byte

// DHCPOption just a DHCP Option
type DHCPOption struct {
	code  OptionCode
	value OptionValue
}

// OptBytes return code Length and value
func (o DHCPOption) OptBytes() []byte {
	// value length limit ???
	return BytesCombine([]byte{byte(o.code)}, []byte{byte(len(o.value))}, o.value)
}

func (o DHCPOption) String() string {
	if s, ok := optionCodeToString[o.code]; ok {
		return s
	}
	return fmt.Sprintf("unknown option (%d)", uint8(o.code))
}

// Option one of Options
type Option interface {
	OptBytes() []byte
	String() string
}

// Options made by Option
type Options map[OptionCode]Option

// DHCP package len
// > 240 is options
// 4 byte for magic code
// where told you maxlen?
const (
	DHCPMAXLEN = 576
	DHCPMINLEN = 236
)

// OpCodes
const (
	BOOTREQUEST byte = 1 // From Client
	BOOTREPLY   byte = 2 // From Server
)

// MessageType represents the possible DHCP message types - DISCOVER, OFFER, etc
type MessageType byte

// DHCP message types
const (
	Discover MessageType = 1 // Broadcast Packet From Client - Can I have an IP?
	Offer    MessageType = 2 // Broadcast From Server - Here's an IP
	Request  MessageType = 3 // Broadcast From Client - I'll take that IP (Also start for renewals)
	Decline  MessageType = 4 // Broadcast From Client - Sorry I can't use that IP
	ACK      MessageType = 5 // From Server, Yes you can have that IP
	NAK      MessageType = 6 // From Server, No you cannot have that IP
	Release  MessageType = 7 // From Client, I don't need that IP anymore
	Inform   MessageType = 8 // From Client, I have this IP and there's nothing you can do about it
)

var messageTypeToString = map[MessageType]string{
	Discover: "DISCOVER",
	Offer:    "OFFER",
	Request:  "REQUEST",
	Decline:  "DECLINE",
	ACK:      "ACK",
	NAK:      "NAK",
	Release:  "RELEASE",
	Inform:   "INFORM",
}

// OptBytes returns the serialized version of this option described by RFC 2132,
// Section 9.6.
func (m MessageType) OptBytes() []byte {
	b := []byte{53, 1}
	b = append(b, byte(m))
	return b
}

// String prints a human-readable message type name.
func (m MessageType) String() string {
	if s, ok := messageTypeToString[m]; ok {
		return fmt.Sprintf("DHCP %s", s)
	}
	return fmt.Sprintf("unknown MessageType (%d)", byte(m))
}

// DHCP Options
const (
	End                          OptionCode = 255
	Pad                          OptionCode = 0
	OptionSubnetMask             OptionCode = 1
	OptionTimeOffset             OptionCode = 2
	OptionRouter                 OptionCode = 3
	OptionTimeServer             OptionCode = 4
	OptionNameServer             OptionCode = 5
	OptionDomainNameServer       OptionCode = 6
	OptionLogServer              OptionCode = 7
	OptionCookieServer           OptionCode = 8
	OptionLPRServer              OptionCode = 9
	OptionImpressServer          OptionCode = 10
	OptionResourceLocationServer OptionCode = 11
	OptionHostName               OptionCode = 12
	OptionBootFileSize           OptionCode = 13
	OptionMeritDumpFile          OptionCode = 14
	OptionDomainName             OptionCode = 15
	OptionSwapServer             OptionCode = 16
	OptionRootPath               OptionCode = 17
	OptionExtensionsPath         OptionCode = 18

	// IP Layer Parameters per Host
	OptionIPForwardingEnableDisable          OptionCode = 19
	OptionNonLocalSourceRoutingEnableDisable OptionCode = 20
	OptionPolicyFilter                       OptionCode = 21
	OptionMaximumDatagramReassemblySize      OptionCode = 22
	OptionDefaultIPTimeToLive                OptionCode = 23
	OptionPathMTUAgingTimeout                OptionCode = 24
	OptionPathMTUPlateauTable                OptionCode = 25

	// IP Layer Parameters per Interface
	OptionInterfaceMTU              OptionCode = 26
	OptionAllSubnetsAreLocal        OptionCode = 27
	OptionBroadcastAddress          OptionCode = 28
	OptionPerformMaskDiscovery      OptionCode = 29
	OptionMaskSupplier              OptionCode = 30
	OptionPerformRouterDiscovery    OptionCode = 31
	OptionRouterSolicitationAddress OptionCode = 32
	OptionStaticRoute               OptionCode = 33

	// Link Layer Parameters per Interface
	OptionTrailerEncapsulation  OptionCode = 34
	OptionARPCacheTimeout       OptionCode = 35
	OptionEthernetEncapsulation OptionCode = 36

	// TCP Parameters
	OptionTCPDefaultTTL        OptionCode = 37
	OptionTCPKeepaliveInterval OptionCode = 38
	OptionTCPKeepaliveGarbage  OptionCode = 39

	// Application and Service Parameters
	OptionNetworkInformationServiceDomain            OptionCode = 40
	OptionNetworkInformationServers                  OptionCode = 41
	OptionNetworkTimeProtocolServers                 OptionCode = 42
	OptionVendorSpecificInformation                  OptionCode = 43
	OptionNetBIOSOverTCPIPNameServer                 OptionCode = 44
	OptionNetBIOSOverTCPIPDatagramDistributionServer OptionCode = 45
	OptionNetBIOSOverTCPIPNodeType                   OptionCode = 46
	OptionNetBIOSOverTCPIPScope                      OptionCode = 47
	OptionXWindowSystemFontServer                    OptionCode = 48
	OptionXWindowSystemDisplayManager                OptionCode = 49
	OptionNetworkInformationServicePlusDomain        OptionCode = 64
	OptionNetworkInformationServicePlusServers       OptionCode = 65
	OptionMobileIPHomeAgent                          OptionCode = 68
	OptionSimpleMailTransportProtocol                OptionCode = 69
	OptionPostOfficeProtocolServer                   OptionCode = 70
	OptionNetworkNewsTransportProtocol               OptionCode = 71
	OptionDefaultWorldWideWebServer                  OptionCode = 72
	OptionDefaultFingerServer                        OptionCode = 73
	OptionDefaultInternetRelayChatServer             OptionCode = 74
	OptionStreetTalkServer                           OptionCode = 75
	OptionStreetTalkDirectoryAssistance              OptionCode = 76

	OptionRelayAgentInformation OptionCode = 82

	// DHCP Extensions
	OptionRequestedIPAddress     OptionCode = 50
	OptionIPAddressLeaseTime     OptionCode = 51
	OptionOverload               OptionCode = 52
	OptionDHCPMessageType        OptionCode = 53
	OptionServerIdentifier       OptionCode = 54
	OptionParameterRequestList   OptionCode = 55
	OptionMessage                OptionCode = 56
	OptionMaximumDHCPMessageSize OptionCode = 57
	OptionRenewalTimeValue       OptionCode = 58
	OptionRebindingTimeValue     OptionCode = 59
	OptionVendorClassIdentifier  OptionCode = 60
	OptionClientIdentifier       OptionCode = 61

	OptionTFTPServerName OptionCode = 66
	OptionBootFileName   OptionCode = 67

	OptionUserClass OptionCode = 77

	OptionClientArchitecture OptionCode = 93

	OptionTZPOSIXString    OptionCode = 100
	OptionTZDatabaseString OptionCode = 101

	OptionDomainSearch OptionCode = 119

	OptionClasslessRouteFormat OptionCode = 121

	// From RFC3942 - Options Used by PXELINUX
	OptionPxelinuxMagic      OptionCode = 208
	OptionPxelinuxConfigfile OptionCode = 209
	OptionPxelinuxPathprefix OptionCode = 210
	OptionPxelinuxReboottime OptionCode = 211
)

/* Notes
A DHCP server always returns its own address in the 'server identifier' option.
DHCP defines a new 'client identifier' option that is used to pass an explicit client identifier to a DHCP server.
via:
https://github.com/krolaw/dhcp4/blob/master/packet.go
https://github.com/insomniacslk/dhcp/blob/master/dhcpv4/options.go
*/

var optionCodeToString = map[OptionCode]string{
	Pad:                    "Pad",
	OptionSubnetMask:       "Subnet Mask",
	OptionTimeOffset:       "Time Offset",
	OptionRouter:           "Router",
	OptionTimeServer:       "Time Server",
	OptionNameServer:       "Name Server",
	OptionDomainNameServer: "Domain Name Server",
	OptionLogServer:        "Log Server",
	// OptionQuoteServer:                                "Quote Server",
	OptionLPRServer:              "LPR Server",
	OptionImpressServer:          "Impress Server",
	OptionResourceLocationServer: "Resource Location Server",
	OptionHostName:               "Host Name",
	OptionBootFileSize:           "Boot File Size",
	OptionMeritDumpFile:          "Merit Dump File",
	OptionDomainName:             "Domain Name",
	OptionSwapServer:             "Swap Server",
	OptionRootPath:               "Root Path",
	OptionExtensionsPath:         "Extensions Path",
	// OptionIPForwarding:                               "IP Forwarding enable/disable",
	// OptionNonLocalSourceRouting:                      "Non-local Source Routing enable/disable",
	OptionPolicyFilter: "Policy Filter",
	// OptionMaximumDatagramAssemblySize:                "Maximum Datagram Reassembly Size",
	// OptionDefaultIPTTL:                               "Default IP Time-to-live",
	OptionPathMTUAgingTimeout:       "Path MTU Aging Timeout",
	OptionPathMTUPlateauTable:       "Path MTU Plateau Table",
	OptionInterfaceMTU:              "Interface MTU",
	OptionAllSubnetsAreLocal:        "All Subnets Are Local",
	OptionBroadcastAddress:          "Broadcast Address",
	OptionPerformMaskDiscovery:      "Perform Mask Discovery",
	OptionMaskSupplier:              "Mask Supplier",
	OptionPerformRouterDiscovery:    "Perform Router Discovery",
	OptionRouterSolicitationAddress: "Router Solicitation Address",
	// OptionStaticRoutingTable:                         "Static Routing Table",
	OptionTrailerEncapsulation: "Trailer Encapsulation",
	// OptionArpCacheTimeout:                            "ARP Cache Timeout",
	OptionEthernetEncapsulation: "Ethernet Encapsulation",
	// OptionDefaulTCPTTL:                               "Default TCP TTL",
	OptionTCPKeepaliveInterval:            "TCP Keepalive Interval",
	OptionTCPKeepaliveGarbage:             "TCP Keepalive Garbage",
	OptionNetworkInformationServiceDomain: "Network Information Service Domain",
	OptionNetworkInformationServers:       "Network Information Servers",
	// OptionNTPServers:                                 "NTP Servers",
	OptionVendorSpecificInformation:                  "Vendor Specific Information",
	OptionNetBIOSOverTCPIPNameServer:                 "NetBIOS over TCP/IP Name Server",
	OptionNetBIOSOverTCPIPDatagramDistributionServer: "NetBIOS over TCP/IP Datagram Distribution Server",
	OptionNetBIOSOverTCPIPNodeType:                   "NetBIOS over TCP/IP Node Type",
	OptionNetBIOSOverTCPIPScope:                      "NetBIOS over TCP/IP Scope",
	OptionXWindowSystemFontServer:                    "X Window System Font Server",
	// OptionXWindowSystemDisplayManger:                 "X Window System Display Manager",
	OptionRequestedIPAddress: "Requested IP Address",
	OptionIPAddressLeaseTime: "IP Addresses Lease Time",
	// OptionOptionOverload:                       "Option Overload",
	OptionDHCPMessageType:        "DHCP Message Type",
	OptionServerIdentifier:       "Server Identifier",
	OptionParameterRequestList:   "Parameter Request List",
	OptionMessage:                "Message",
	OptionMaximumDHCPMessageSize: "Maximum DHCP Message Size",
	// OptionRenewTimeValue:                       "Renew Time Value",
	OptionRebindingTimeValue: "Rebinding Time Value",
	// OptionClassIdentifier:                      "Class Identifier",
	OptionClientIdentifier: "Client identifier",
	// OptionNetWareIPDomainName:                  "NetWare/IP Domain Name",
	// OptionNetWareIPInformation:                 "NetWare/IP Information",
	OptionNetworkInformationServicePlusDomain:  "Network Information Service+ Domain",
	OptionNetworkInformationServicePlusServers: "Network Information Service+ Servers",
	OptionTFTPServerName:                       "TFTP Server Name",
	// OptionBootfileName:                         "Bootfile Name",
	OptionMobileIPHomeAgent: "Mobile IP Home Agent",
	// OptionSimpleMailTransportProtocolServer:   "SMTP Server",
	OptionPostOfficeProtocolServer: "POP Server",
	// OptionNetworkNewsTransportProtocolServer:  "NNTP Server",
	OptionDefaultWorldWideWebServer:      "Default WWW Server",
	OptionDefaultFingerServer:            "Default Finger Server",
	OptionDefaultInternetRelayChatServer: "Default IRC Server",
	OptionStreetTalkServer:               "StreetTalk Server",
	// OptionStreetTalkDirectoryAssistanceServer: "StreetTalk Directory Assistance Server",
	// OptionUserClassInformation:                "User Class Information",
	// OptionSLPDirectoryAgent:                   "SLP DIrectory Agent",
	// OptionSLPServiceScope:                     "SLP Service Scope",
	// OptionRapidCommit:                         "Rapid Commit",
	// OptionFQDN:                                "FQDN",
	OptionRelayAgentInformation: "Relay Agent Information",
	// OptionInternetStorageNameService:          "Internet Storage Name Service",
	// Option 84 returned in RFC 3679
	// OptionNDSServers:                       "NDS Servers",
	// OptionNDSTreeName:                      "NDS Tree Name",
	// OptionNDSContext:                       "NDS Context",
	// OptionBCMCSControllerDomainNameList:    "BCMCS Controller Domain Name List",
	// OptionBCMCSControllerIPv4AddressList:   "BCMCS Controller IPv4 Address List",
	// OptionAuthentication:                   "Authentication",
	// OptionClientLastTransactionTime:        "Client Last Transaction Time",
	// OptionAssociatedIP:                     "Associated IP",
	// OptionClientSystemArchitectureType:     "Client System Architecture Type",
	// OptionClientNetworkInterfaceIdentifier: "Client Network Interface Identifier",
	// OptionLDAP: "LDAP",
	// Option 96 returned in RFC 3679
	// OptionClientMachineIdentifier:     "Client Machine Identifier",
	// OptionOpenGroupUserAuthentication: "OpenGroup's User Authentication",
	// OptionGeoConfCivic:                "GEOCONF_CIVIC",
	// OptionIEEE10031TZString:           "IEEE 1003.1 TZ String",
	// OptionReferenceToTZDatabase:       "Reference to the TZ Database",
	// Options 102-111 returned in RFC 3679
	// OptionNetInfoParentServerAddress: "NetInfo Parent Server Address",
	// OptionNetInfoParentServerTag:     "NetInfo Parent Server Tag",
	// OptionURL:                        "URL",
	// Option 115 returned in RFC 3679
	// OptionAutoConfigure:                   "Auto-Configure",
	// OptionNameServiceSearch:               "Name Service Search",
	// OptionSubnetSelection:                 "Subnet Selection",
	// OptionDNSDomainSearchList:             "DNS Domain Search List",
	// OptionSIPServers:                      "SIP Servers",
	// OptionClasslessStaticRoute:            "Classless Static Route",
	// OptionCCC:                             "CCC, CableLabs Client Configuration",
	// OptionGeoConf:                         "GeoConf",
	// OptionVendorIdentifyingVendorClass:    "Vendor-Identifying Vendor Class",
	// OptionVendorIdentifyingVendorSpecific: "Vendor-Identifying Vendor-Specific",
	// Options 126-127 returned in RFC 3679
	// OptionTFTPServerIPAddress:                   "TFTP Server IP Address",
	// OptionCallServerIPAddress:                   "Call Server IP Address",
	// OptionDiscriminationString:                  "Discrimination String",
	// OptionRemoteStatisticsServerIPAddress:       "RemoteStatistics Server IP Address",
	// Option8021PVLANID:                           "802.1P VLAN ID",
	// Option8021QL2Priority:                       "802.1Q L2 Priority",
	// OptionDiffservCodePoint:                     "Diffserv Code Point",
	// OptionHTTPProxyForPhoneSpecificApplications: "HTTP Proxy for phone-specific applications",
	// OptionPANAAuthenticationAgent:               "PANA Authentication Agent",
	// OptionLoSTServer:                            "LoST Server",
	// OptionCAPWAPAccessControllerAddresses:       "CAPWAP Access Controller Addresses",
	// OptionOPTIONIPv4AddressMoS:             "OPTION-IPv4_Address-MoS",
	// OptionOPTIONIPv4FQDNMoS:                "OPTION-IPv4_FQDN-MoS",
	// OptionSIPUAConfigurationServiceDomains: "SIP UA Configuration Service Domains",
	// OptionOPTIONIPv4AddressANDSF:           "OPTION-IPv4_Address-ANDSF",
	// OptionOPTIONIPv6AddressANDSF:           "OPTION-IPv6_Address-ANDSF",
	// // Options 144-149 returned in RFC 3679
	// OptionTFTPServerAddress: "TFTP Server Address",
	// OptionStatusCode:        "Status Code",
	// OptionBaseTime:          "Base Time",
	// OptionStartTimeOfState:  "Start Time of State",
	// OptionQueryStartTime:    "Query Start Time",
	// OptionQueryEndTime: "Query End Time",
	// OptionDHCPState:    "DHCP Staet",
	// OptionDataSource:   "Data Source",
	// // Options 158-174 returned in RFC 3679
	// OptionEtherboot:                        "Etherboot",
	// OptionIPTelephone:                      "IP Telephone",
	// OptionEtherbootPacketCableAndCableHome: "Etherboot / PacketCable and CableHome",
	// // Options 178-207 returned in RFC 3679
	// OptionPXELinuxMagicString:  "PXELinux Magic String",
	// OptionPXELinuxConfigFile:   "PXELinux Config File",
	// OptionPXELinuxPathPrefix:   "PXELinux Path Prefix",
	// // OptionPXELinuxRebootTime:   "PXELinux Reboot Time",
	// OptionOPTION6RD:            "OPTION_6RD",
	// OptionOPTIONv4AccessDomain: "OPTION_V4_ACCESS_DOMAIN",
	// // Options 214-219 returned in RFC 3679
	// OptionSubnetAllocation:        "Subnet Allocation",
	// OptionVirtualSubnetAllocation: "Virtual Subnet Selection",
	// Options 222-223 returned in RFC 3679
	// Options 224-254 are reserved for private use

}
