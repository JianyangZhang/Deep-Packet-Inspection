package tests

import (
	"fmt"
	"log"
	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

const PROTOCOL_IPV4_TCP string = "ip4:tcp"

/* 测试 */
func TryPush() {
	if packet, err := InitPacket("127.0.0.1", 80, "127.0.0.1", 80); err != nil {
		panic(err)
	} else {
		PushPacket(packet, "127.0.0.1")
	}
}

/*
	推送报文
*/
func PushPacket(packet gopacket.SerializeBuffer, ip_dst string) {
	// 推送
	dstIPaddr := net.IPAddr{
		IP: net.ParseIP(ip_dst),
	}
	fmt.Println("\n报文创建成功, 开始推送...\n ")
	connection, err := net.ListenPacket(PROTOCOL_IPV4_TCP, "127.0.0.1")
	if err != nil {
		panic(err)
	}
	_, err = connection.WriteTo(packet.Bytes(), &dstIPaddr) // 为什么要传入指针
	if err != nil {
		panic(err)
	}
	log.Print("发送成功!\n")
}

/*
	创建报文
*/
func InitPacket(ip_src string, port_src int, ip_dst string, port_dst int) (gopacket.SerializeBuffer, error) {
	srcIP := net.ParseIP(ip_src)
	dstIP := net.ParseIP(ip_dst)
	ipLayer := layers.IPv4{
		SrcIP:    srcIP,
		DstIP:    dstIP,
		Protocol: layers.IPProtocolTCP,
	}
	tcpLayer := layers.TCP{
		SrcPort: layers.TCPPort(port_src),
		DstPort: layers.TCPPort(port_dst),
		SYN:     true,
	}
	tcpLayer.SetNetworkLayerForChecksum(&ipLayer)
	packet := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}
	err := gopacket.SerializeLayers(packet, opts, &ipLayer, &tcpLayer)
	return packet, err
}
