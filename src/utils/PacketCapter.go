/*

 */
package utils

import (
	"fmt"
	"io"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

const OFFLINE_PCAP_PATH string = "../samples/test_ethernet.pcap"
const VPN_ADAPTER_PATH string = "\\Device\\NPF_{3687A031-BD34-4472-ACCA-29349877F279}"

var eth layers.Ethernet
var ip4 layers.IPv4
var ip6 layers.IPv6
var tcp layers.TCP
var udp layers.UDP
var payload gopacket.Payload
var parser *gopacket.DecodingLayerParser = gopacket.NewDecodingLayerParser(layers.LayerTypeEthernet, &eth, &ip4, &ip6, &tcp, &udp, &payload)

/*
	本地pcap数据包测试
	mode: 0打印数据包的基本信息; 1完整打印包的信息
*/
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
					handlePacket(packetData)
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

/*
	实时数据包测试
	amount int: 将要分析的数据包的数量上限
	mode: 0打印数据包的基本信息; 1完整打印包的信息
*/
func LivePacketsTest(path string, amount int, mode int) {
	fmt.Println("\n --- 实时数据包测试开始 ---\n ")
	count := 0
	if live_source, err := pcap.OpenLive(path, 1600, true, pcap.BlockForever); err != nil {
		panic(err)
	} else {
		switch mode {
		case 0:
			for {
				if packetData, _, err := live_source.ReadPacketData(); err != nil {
					panic(err)
				} else {
					handlePacket(packetData)
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

/*
	简单分析一个包的信息：每一层的名称、来源、目标
*/
func handlePacket(packetData []byte) {
	decodedLayers := make([]gopacket.LayerType, 0, 10)

	if err := parser.DecodeLayers(packetData, &decodedLayers); err != nil {
		panic(err)
	}
	fmt.Println("得到一个数据包")
	for _, ltype := range decodedLayers {
		fmt.Println(ltype, "层:")
		switch ltype {
		case layers.LayerTypeEthernet:
			fmt.Println("Mac", eth.SrcMAC, "=>", eth.DstMAC)
		case layers.LayerTypeIPv4:
			fmt.Println("IP", ip4.SrcIP, "=>", ip4.DstIP)
		case layers.LayerTypeIPv6:
			fmt.Println("IP", ip6.SrcIP, "=>", ip6.DstIP)
		case layers.LayerTypeTCP:
			fmt.Println("端口", tcp.SrcPort, "=>", tcp.DstPort)
		case layers.LayerTypeUDP:
			fmt.Println("端口", udp.SrcPort, "=>", udp.DstPort)
		default:
		}
	}
	fmt.Println()
}

/*
	打印所有可用设备信息
*/
func GetDevices() {
	fmt.Println("\n --- 获取所有可用设备信息 ---\n ")
	devices, err := pcap.FindAllDevs()
	if err != nil {
		panic(err)
	}

	for _, device := range devices {
		fmt.Println("设备名称: ", device.Name)
		fmt.Println("设备描述: ", device.Description)
		fmt.Println("设备地址: ")
		for _, address := range device.Addresses {
			fmt.Println("- IP 地址: ", address.IP)
			fmt.Println("- 子网掩码: ", address.Netmask)
		}
		fmt.Println()
	}
}
