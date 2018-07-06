package tests

import (
	"fmt"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

var OFFLINE_PCAP_PATH string = "../samples/test_ethernet.pcap"
var VPN_ADAPTER_PATH string = "\\Device\\NPF_{3687A031-BD34-4472-ACCA-29349877F279}"

/*
	本地pcap数据包测试
*/
func OfflinePacketsTest(path string) {
	fmt.Println("\n --- 本地pcap数据包测试开始 ---\n ")
	if offline_source, err := pcap.OpenOffline(path); err != nil { // 得到原始数据源
		panic(err)
	} else {
		packetSource := gopacket.NewPacketSource(offline_source, offline_source.LinkType()) // 原始数据源 => gopacket数据源
		for packet := range packetSource.Packets() {                                        // gopacket数据源 => Packet channel
			handlePacket(packet)
		}
	}
}

/*
	实时数据包测试
	参数amount int: 将要分析的数据包的数量上限
*/
func LivePacketsTest(path string, amount int) {
	fmt.Println("\n --- 实时数据包测试开始 ---\n ")
	count := 0
	if live_source, err := pcap.OpenLive(path, 1600, true, pcap.BlockForever); err != nil {
		panic(err)
	} else {
		packetSource := gopacket.NewPacketSource(live_source, live_source.LinkType())
		for packet := range packetSource.Packets() {
			handlePacket(packet)
			count++
			if count >= amount {
				break
			}
		}
	}
}

/*
	简单分析一个包的信息：是否为TCP协议、来源端口、目标端口、每一层的名称
	返回值bool：true为TCP包; false为非TCP包
*/
func handlePacket(packet gopacket.Packet) bool {
	var isTCP bool = true
	if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil { // 判断是否为TCP协议
		fmt.Println("TCP数据包")
		tcp, _ := tcpLayer.(*layers.TCP) // 得到来源端口与目标端口
		fmt.Printf("从端口 %d 发送至端口 %d\n", tcp.SrcPort, tcp.DstPort)

	} else {
		isTCP = false
		fmt.Println("非TCP数据包")
	}
	// 打印此包每一层的名称
	for _, layer := range packet.Layers() {
		fmt.Println("此数据包有:", layer.LayerType(), "层")
	}
	fmt.Println()
	return isTCP
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
