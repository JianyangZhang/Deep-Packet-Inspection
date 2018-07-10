package utils

import (
	"fmt"
	"log"
	"net"

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
	ethernetLayer := &layers.Ethernet{
		SrcMAC: net.HardwareAddr{src_mac[0], src_mac[1], src_mac[2], src_mac[3], src_mac[4], src_mac[5]},
		DstMAC: net.HardwareAddr{dst_mac[0], dst_mac[1], dst_mac[2], dst_mac[3], dst_mac[4], dst_mac[5]},
	}
	ipLayer := &layers.IPv4{
		SrcIP: net.IP{src_ip[0], src_ip[1], src_ip[2], src_ip[3]},
		DstIP: net.IP{dst_ip[0], dst_ip[1], dst_ip[2], dst_ip[3]},
	}
	tcpLayer := &layers.TCP{
		SrcPort: layers.TCPPort(src_port),
		DstPort: layers.TCPPort(dst_port),
	}
	gopacket.SerializeLayers(buffer, options,
		ethernetLayer,
		ipLayer,
		tcpLayer,
		gopacket.Payload(payload),
	)
	fmt.Println("=============> 报文创建成功")
	// 将报文转为byte数组
	return buffer.Bytes()

}
