/* 抓包 */
package utils

import (
	"fmt"
	"io"
	"log"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

const OFFLINE_PCAP_PATH string = "../samples/test_ethernet.pcap"
const ADAPTER_NAME string = "eth0"
const VPN_ADAPTER_NAME string = "\\Device\\NPF_{3687A031-BD34-4472-ACCA-29349877F279}"
const CENTOS_ADAPTER_NAME string = "ens33"
const WIN_LOOPBACK_ADAPTER_NAME string = "\\Device\\NPF_{EF48B70C-A89C-4E75-BBDD-B03A3ACCFC9E}"
const LINUX_LOOPBACK_ADAPTER_NAME string = "lo"
const FILTER_TCP_80 string = "tcp and port 80"
const FILTER_UDP_53 string = "udp and port 53"
const FILTER_ALL string = ""

/*
	抓取 TCP:80 数据包
	参数 deviceName: 抓取通过此设备的报文; packetChannel: 将抓取到的报文信息传入此channel; filter: 过滤器
*/
func GetLivePackets(deviceName string, packetChannel chan PacketInfo, filter string) {
	if handle, err := pcap.OpenLive(deviceName, 3000, true, pcap.BlockForever); err != nil { // 得到设备的handle
		panic(err)
	} else {
		defer handle.Close()
		err = handle.SetBPFFilter(filter) // 设置过滤器
		if err != nil {
			panic(err)
		}
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType()) // 得到数据源
		fmt.Println("开始抓包...")
		for packet := range packetSource.Packets() {
			if filter != "" {
				fmt.Println("\n抓取到一个", filter, "的数据包，等待解析")
			} else {
				fmt.Println("\n抓取到一个数据包，等待解析")
			}
			packetChannel <- decodePacket(packet)
		}
	}
}

/*
	解析一个数据包，将结果传入PacketInfo并返回
*/
func decodePacket(packet gopacket.Packet) PacketInfo {
	var rst PacketInfo
	var rstLayers []string
	for _, layer := range packet.Layers() {
		rstLayers = append(rstLayers, layer.LayerType().String())
	}

	ethernetLayer := packet.Layer(layers.LayerTypeEthernet)
	ipLayer := packet.Layer(layers.LayerTypeIPv4)
	tcpLayer := packet.Layer(layers.LayerTypeTCP)
	udpLayer := packet.Layer(layers.LayerTypeUDP)
	dnsLayer := packet.Layer(layers.LayerTypeDNS)
	applicationLayer := packet.ApplicationLayer()

	if ethernetLayer != nil {
		ethernetPacket, _ := ethernetLayer.(*layers.Ethernet)
		rst.SrcMac = ethernetPacket.SrcMAC.String()
		rst.DstMac = ethernetPacket.DstMAC.String()
		rst.EthType = ethernetPacket.EthernetType.String()
	}

	if ipLayer != nil {
		ip, _ := ipLayer.(*layers.IPv4)
		rst.SrcIP = ip.SrcIP.String()
		rst.DstIP = ip.DstIP.String()
		rst.Protocal = ip.Protocol.String()
	}

	if tcpLayer != nil {
		tcp, _ := tcpLayer.(*layers.TCP)
		rst.SrcPort = tcp.SrcPort.String()
		rst.DstPort = tcp.DstPort.String()
		rst.Sequence = tcp.Seq
	} else if udpLayer != nil {
		udp, _ := udpLayer.(*layers.UDP)
		rst.SrcPort = udp.SrcPort.String()
		rst.DstPort = udp.DstPort.String()
	}

	if dnsLayer != nil {
		dns, _ := dnsLayer.(*layers.DNS)
		rst.IsDNS = true
		if len(dns.Questions) > 0 {
			rst.Question = string(dns.Questions[0].Name)
		}
		if len(dns.Answers) > 0 {
			rst.Answer = dns.Answers[0].IP.String()
		}
	} else {
		rst.IsDNS = false
	}
	if applicationLayer != nil {
		rst.Payload = string(applicationLayer.Payload())
	}
	rst.Layers = rstLayers
	return rst
}

type PacketInfo struct {
	Layers   []string
	Protocal string
	EthType  string
	SrcMac   string
	DstMac   string
	SrcIP    string
	DstIP    string
	SrcPort  string
	DstPort  string
	Sequence uint32
	Payload  string
	IsDNS    bool
	Question string
	Answer   string
}

func (this PacketInfo) PrintAll() {
	for _, s := range this.Layers {
		fmt.Print("此报文包含", s, "层 ")
	}
	fmt.Println("\n报文协议:", this.Protocal)
	fmt.Println("IP类型:", this.EthType)
	fmt.Println("Mac:", this.SrcMac, "=>", this.DstMac)
	fmt.Println("IP:", this.SrcIP, "=>", this.DstIP)
	fmt.Println("端口:", this.SrcPort, "=>", this.DstPort)
	fmt.Println("报文序列:", fmt.Sprint(this.Sequence))
	// fmt.Println("Payload内容:", this.Payload)
	if this.IsDNS {
		fmt.Println("DNS查询:", this.Question)
		fmt.Println("DNS返回:", this.Answer)
	}
	fmt.Println("----------解析完成----------\n ")
}

/*
	获取所有可用设备信息
*/
func GetDevices() []Device {
	fmt.Println("\n --- 获取所有可用设备信息 ---\n ")
	devices, err := pcap.FindAllDevs()
	if err != nil {
		panic(err)
	}
	rst := make([]Device, 10)
	for _, device := range devices {
		var one Device = Device{name: "", description: "", ip: "", netmask: ""}
		fmt.Println("设备名称: ", device.Name)
		one.name = device.Name
		fmt.Println("设备描述: ", device.Description)
		one.description = device.Description
		fmt.Println("设备地址: ")
		for _, address := range device.Addresses {
			fmt.Println("- IP 地址: ", address.IP)
			one.ip = string(address.IP)
			fmt.Println("- 子网掩码: ", address.Netmask)
			one.netmask = string(address.Netmask)
		}
		fmt.Println()
		rst = append(rst, one)
	}
	return rst
}

type Device struct {
	name        string
	description string
	ip          string
	netmask     string
}

//----------------------------------------- 测试 --------------------------------------------

var eth layers.Ethernet
var ip4 layers.IPv4
var ip6 layers.IPv6
var tcp layers.TCP
var udp layers.UDP
var dns layers.DNS
var payload gopacket.Payload
var parser *gopacket.DecodingLayerParser = gopacket.NewDecodingLayerParser(layers.LayerTypeEthernet, &eth, &ip4, &ip6, &tcp, &udp, &payload, &dns)

//	本地pcap数据包测试
//	参数 path: .pcap文件路径; mode: 0数据包基本信息, 1数据包完整信息
func OfflinePacketsTest(path string, mode int) {
	fmt.Println("\n --- 本地pcap数据包测试开始 ---\n ")
	if offline_source, err := pcap.OpenOffline(path); err != nil { // 得到原始数据源
		panic(err)
	} else {
		switch mode {
		case 0:
			for {
				if packetData, _, err := offline_source.ReadPacketData(); err != nil {
					if err == io.EOF {
						break
					} else {
						panic(err)
					}
				} else {
					handlePacket(packetData, true)
				}
			}
			break
		case 1:
			packetSource := gopacket.NewPacketSource(offline_source, offline_source.LinkType())
			for packet := range packetSource.Packets() {
				fmt.Println(packet)
			}
			break
		}
	}
}

//	实时数据包测试
//	参数 deviceName: 抓取通过此设备的报文; amount int: 将要分析的数据包的数量上限; mode: 0数据包基本信息, 1数据包完整信息
func LivePacketsTest(deviceName string, amount int, mode int) {
	fmt.Println("\n --- 实时数据包测试开始 ---\n ")
	count := 0
	if live_source, err := pcap.OpenLive(deviceName, 3000, true, pcap.BlockForever); err != nil {
		panic(err)
	} else {
		switch mode {
		case 0:
			for {
				if packetData, _, err := live_source.ReadPacketData(); err != nil {
					panic(err)
				} else {
					handlePacket(packetData, true)
					count++
					if count >= amount {
						break
					}
				}
			}
			break
		case 1:
			packetSource := gopacket.NewPacketSource(live_source, live_source.LinkType())
			for packet := range packetSource.Packets() {
				fmt.Println(packet)
				count++
				if count >= amount {
					break
				}
			}
			break
		}
	}
}

// 简单分析一个包的信息：每一层的名称、来源、目标
// 参数 packetData: 数据包原始数据; flag: 是否打印数据包信息
func handlePacket(packetData []byte, print bool) PacketInfo {
	var rst PacketInfo
	var rstLayers []string
	decodedLayers := make([]gopacket.LayerType, 0, 10)
	if err := parser.DecodeLayers(packetData, &decodedLayers); err != nil {
		log.Fatal(err)
	}
	for _, ltype := range decodedLayers {
		if print {
			fmt.Println(ltype, "层:")
			rstLayers = append(rstLayers, ltype.String())
		}
		switch ltype {
		case layers.LayerTypeEthernet:
			if print {
				fmt.Println("Mac", eth.SrcMAC, "=>", eth.DstMAC)
			}
			rst.SrcMac = eth.SrcMAC.String()
			rst.DstMac = eth.DstMAC.String()
		case layers.LayerTypeIPv4:
			if print {
				fmt.Println("IP", ip4.SrcIP, "=>", ip4.DstIP)
			}
			rst.SrcIP = ip4.SrcIP.String()
			rst.DstIP = ip4.DstIP.String()
		case layers.LayerTypeIPv6:
			if print {
				fmt.Println("IP", ip6.SrcIP, "=>", ip6.DstIP)
			}
			rst.SrcIP = ip6.SrcIP.String()
			rst.DstIP = ip6.DstIP.String()
		case layers.LayerTypeTCP:
			if print {
				fmt.Println("端口", tcp.SrcPort, "=>", tcp.DstPort)
			}
			rst.SrcPort = tcp.SrcPort.String()
			rst.DstPort = tcp.DstPort.String()
		case layers.LayerTypeUDP:
			if print {
				fmt.Println("端口", udp.SrcPort, "=>", udp.DstPort)
			}
			rst.SrcPort = udp.SrcPort.String()
			rst.DstPort = udp.DstPort.String()
		default:
		}
	}
	fmt.Println()
	rst.Layers = rstLayers
	return rst
}
