/* 创建报文 发送报文*/
package utils

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

/* 将byte array格式的报文推送给指定的设备 */
func SendPacket(deviceName string, data []byte) {
	if handle, err := pcap.OpenLive(deviceName, 1024, false, 60*time.Second); err != nil { // 得到设备的handle
		log.Fatal(err)
	} else {
		// 发送报文
		fmt.Println("=============> 准备发送报文")
		err = handle.WritePacketData(data)
		if err != nil {
			fmt.Println("=============> 报文发送失败")
			log.Fatal(err)
		} else {
			fmt.Println("=============> 报文发送成功")
		}
	}
}

/* 创建byte array格式的完整报文 */
func CreatePacket(protocal string, src_mac string, dst_mac string, src_ip string, dst_ip string, src_port int, dst_port int, payload []byte) []byte {

	// 准备报文内容
	buffer := gopacket.NewSerializeBuffer()
	options := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}
	tcpLayer := &layers.TCP{
		SrcPort: layers.TCPPort(src_port),
		DstPort: layers.TCPPort(dst_port),
	}
	udpLayer := &layers.UDP{
		SrcPort: layers.UDPPort(src_port),
		DstPort: layers.UDPPort(dst_port),
	}
	ipLayer := &layers.IPv4{
		Version:    4,
		IHL:        5,
		TOS:        0,
		Id:         0,
		Flags:      0,
		FragOffset: 0,
		TTL:        255,
		SrcIP:      net.ParseIP(src_ip),
		DstIP:      net.ParseIP(dst_ip),
	}

	src_mac_ready, err0 := net.ParseMAC(src_mac)
	dst_mac_ready, err1 := net.ParseMAC(dst_mac)
	if err0 != nil {
		panic(err0)
	}
	if err1 != nil {
		panic(err1)
	}
	ethernetLayer := &layers.Ethernet{
		SrcMAC:       src_mac_ready,
		DstMAC:       dst_mac_ready,
		EthernetType: layers.EthernetTypeIPv4,
	}

	if "udp" == strings.ToLower(protocal) {
		ipLayer.Protocol = layers.IPProtocolUDP
		udpLayer.SetNetworkLayerForChecksum(ipLayer)
		gopacket.SerializeLayers(buffer, options,
			ethernetLayer,
			ipLayer,
			udpLayer,
			gopacket.Payload(payload),
		)
		fmt.Println("=============> UDP报文创建成功")
	} else if "tcp" == strings.ToLower(protocal) {
		ipLayer.Protocol = layers.IPProtocolTCP
		tcpLayer.SetNetworkLayerForChecksum(ipLayer)
		gopacket.SerializeLayers(buffer, options,
			ethernetLayer,
			ipLayer,
			tcpLayer,
			gopacket.Payload(payload),
		)
		fmt.Println("=============> TCP报文创建成功")
	} else {
		panic("protocal must be 'tcp' or 'udp'")
	}

	// 将报文转为byte数组
	return buffer.Bytes()
}
