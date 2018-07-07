package tests

import (
	"fmt"
	"log"
	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

/* 推送报文 */
func TryPush() {
	if packet, dstIPaddr, err := initPacket(); err != nil {
		panic(err)
	} else {
		// 推送
		fmt.Println("\n报文创建成功, 开始推送...\n ")
		connection, err := net.ListenPacket("ip4:tcp", "0.0.0.1")
		if err != nil {
			panic(err)
		}
		_, err = connection.WriteTo(packet.Bytes(), &dstIPaddr)
		if err != nil {
			panic(err)
		}
		log.Print("发送成功!\n")
	}
}

/* 创建报文 */
func initPacket() (gopacket.SerializeBuffer, net.IPAddr, error) {
	srcIP := net.ParseIP("127.0.0.1")
	dstIP := net.ParseIP("127.0.0.1")
	dstIPaddr := net.IPAddr{
		IP: dstIP,
	}
	ipLayer := layers.IPv4{
		SrcIP:    srcIP,
		DstIP:    dstIP,
		Protocol: layers.IPProtocolTCP,
	}
	tcpLayer := layers.TCP{
		SrcPort: layers.TCPPort(80),
		DstPort: layers.TCPPort(80),
		SYN:     true,
	}
	tcpLayer.SetNetworkLayerForChecksum(&ipLayer)
	packet := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}
	err := gopacket.SerializeLayers(packet, opts, &ipLayer, &tcpLayer)
	return packet, dstIPaddr, err
}
