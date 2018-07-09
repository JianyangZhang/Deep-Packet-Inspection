package utils

import (
	"fmt"
	"log"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

/* 将byte array格式的报文推送给指定的设备 */
func SendPacket(deviceName string, data []byte) {
	if handle, err := pcap.OpenLive(deviceName, 3000, true, pcap.BlockForever); err != nil { // 得到设备的handle
		panic(err)
	} else {
		// 发送报文
		err = handle.WritePacketData(data)
		if err != nil {
			log.Fatal(err)
		} else {
			fmt.Println("=============> 报文已推送")
		}
	}
}

/* 创建byte array格式的完整报文 */
func CreatePacket(src_mac []byte, dst_mac []byte, src_ip []byte, dst_ip []byte, src_port int, dst_port int, payload []byte) []byte {

	// 准备报文内容
	buffer := gopacket.NewSerializeBuffer()
	options := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}
	/*
			ethernetLayer := &layers.Ethernet{
				SrcMAC: src_mac,
				DstMAC: dst_mac,
		    }
	*/
	ipLayer := &layers.IPv4{
		SrcIP: src_ip,
		DstIP: dst_ip,
	}
	tcpLayer := &layers.TCP{
		SrcPort: layers.TCPPort(src_port),
		DstPort: layers.TCPPort(dst_port),
	}
	gopacket.SerializeLayers(buffer, options,
		// ethernetLayer,
		ipLayer,
		tcpLayer,
		gopacket.Payload(payload),
	)
	fmt.Println("=============> 报文创建成功")
	// 将报文转为byte数组
	return buffer.Bytes()

}
