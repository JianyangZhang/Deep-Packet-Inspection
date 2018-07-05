package main

import (
	"fmt"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

var SAMPLE_PCAP_PATH string = "../samples/test_ethernet.pcap"

func main() {
	fmt.Println("-DPI Test Begin-")
	fmt.Println()

	if psource, err := pcap.OpenOffline(SAMPLE_PCAP_PATH); err != nil { // 得到原始数据源
		panic(err)
	} else {
		packetSource := gopacket.NewPacketSource(psource, psource.LinkType()) // 原始数据源 => gopacket数据源
		for packet := range packetSource.Packets() {                          // gopacket数据源 => Packet channel
			handlePacket(packet)
		}
	}

	fmt.Println("-DPI Test End-")
}

/* 简单分析一个包的信息：是否为TCP协议、来源端口、目标端口、每一层的名称 */
func handlePacket(packet gopacket.Packet) {
	if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil { // 判断是否为TCP协议
		fmt.Println("这是一个TCP的数据包")
		tcp, _ := tcpLayer.(*layers.TCP) // 得到来源端口与目标端口
		fmt.Printf("从端口 %d 发送至端口 %d\n", tcp.SrcPort, tcp.DstPort)
	} else {
		fmt.Println("这不是一个TCP的数据包")
	}
	// 打印此包每一层的名称
	for _, layer := range packet.Layers() {
		fmt.Println("此数据包有:", layer.LayerType(), "层")
	}
	fmt.Println()
}
