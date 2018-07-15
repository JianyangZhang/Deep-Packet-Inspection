package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	"./utils"
)

func main() {
	// LinkRawSocketTest()
	// DNSServerTest()
	PacketTransferTest()
}

/*------------------------------ 测试 ------------------------------*/

/*
// LinkRawSocket.go 测试
func LinkRawSocketTest() {
	go func() {
		utils.Sender()

	}()
	utils.Receiver()
}
*/

// DNSServer.go 测试
func DNSServerTest() {
	_port := flag.Int("p", 3535, "服务器端口")
	flag.Parse()
	port := *_port
	utils.StartDNSServer(port)
}

// PacketCapter.go PacketSender.go 测试
func PacketTransferTest() {
	_protocal := flag.String("p", "tcp", "协议类型")
	_device := flag.String("d", utils.VPN_ADAPTER_NAME, "网络适配器名称")
	_filter := flag.String("f", "", "报文过滤")
	flag.Parse()
	protocal := *_protocal
	device := *_device
	filter := *_filter
	if filter != "" {
		fmt.Println("已设置过滤器:", filter)
	}
	// utils.GetDevices()
	pChan := make(chan utils.PacketInfo, 100) // 数据包channel

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
		utils.GetLivePackets(device, pChan, filter)
	}()

	// 推送一个包
	go func() {
		time.Sleep(time.Second)
		// 创建一个报文
		// 顺序 src_mac string, dst_mac string, src_ip string, dst_ip string, src_port int, dst_port int, payload []byte
		newPacket := utils.CreatePacket(
			protocal,
			"11:22:33:44:55:66",
			"66:55:44:33:22:11",
			"192.166.6.6",
			"192.188.8.8",
			60,
			80,
			[]byte{5, 8, 10})
		// 发送报文
		utils.SendPacket(device, newPacket)
	}()

	// 打印抓到的包
	for p := range pChan {
		p.PrintAll()
	}
}
