package utils

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

/*------------------------------ 测试 ------------------------------*/

// LinkRawSocket.go 测试
func LinkRawSocketTest() {
	go func() {
		Sender()

	}()
	Receiver()
}

// DNSServer.go 测试
func DNSServerTest() {
	_port := flag.Int("p", 3535, "服务器端口")
	flag.Parse()
	port := *_port
	StartDNSServer(port)
}

// PacketCapter.go PacketSender.go 测试
func PacketTransferTest() {
	_protocal := flag.String("p", "tcp", "协议类型")
	_device := flag.String("d", VPN_ADAPTER_NAME, "网络适配器名称")
	_filter := flag.String("f", "", "报文过滤")
	flag.Parse()
	protocal := *_protocal
	device := *_device
	filter := *_filter
	if filter != "" {
		fmt.Println("已设置过滤器:", filter)
	}
	// utils.GetDevices()
	pChan := make(chan PacketInfo, 100) // 数据包channel

	// ctrl-c 关闭channel, 结束程序
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, os.Kill)
	go func() {
		for _ = range sigChan {
			fmt.Println("Ctrl-C 停止抓取...")
			close(pChan)
		}
	}()

	// 开始抓包
	go func() {
		GetLivePackets(device, pChan, filter)
	}()

	// 推送一个包
	go func() {
		time.Sleep(time.Second)
		// 创建一个报文
		// 顺序 src_mac string, dst_mac string, src_ip string, dst_ip string, src_port int, dst_port int, payload []byte
		newPacket := CreatePacket(
			protocal,
			"11:22:33:44:55:66",
			"66:55:44:33:22:11",
			"192.166.6.6",
			"192.188.8.8",
			60,
			80,
			[]byte{5, 8, 10})
		// 发送报文
		SendPacket(device, newPacket)
	}()

	// 打印抓到的包
	for p := range pChan {
		p.PrintAll()
	}
}

/*
	测试PacketSender.go中DNS包的创建
*/
func SendDNSPacketTest() {
	_device := flag.String("d", VPN_ADAPTER_NAME, "网络适配器名称")
	flag.Parse()
	device := *_device

	// 抓一个dns包
	dns_request_packet := GetOneLivePacket(device, "udp and port 53")
	DecodePacket(dns_request_packet).PrintAll()

	// 创建dns响应报文
	dns_response_packet := CreatePacket(
		"udp",
		"11:22:33:44:55:66",
		"66:55:44:33:22:11",
		"192.166.6.6",
		"192.188.8.8",
		80,
		53,
		CreateDNSResponse(dns_request_packet, "192.222.2.2"))

	// 发送报文
	SendPacket(device, dns_response_packet)

	// 解析检查发送的报文
	// (因为 “自己抓取自己发送的包” 在 PacketTransferTest() 中已经测试成功，这里发出的"假"DNS包，只测试能否被解析，而不去测试能否被抓取)
	mypacket := gopacket.NewPacket(dns_response_packet, layers.LayerTypeEthernet, gopacket.Default)
	DecodePacket(mypacket).PrintAll()
}
