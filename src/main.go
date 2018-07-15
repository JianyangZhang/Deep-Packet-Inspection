package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"./utils"
)

func main() {
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
	utils.StartDNSServer(3535)
}

// PacketCapter.go PacketSender.go 测试
func PacketTransferTest() {
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
		utils.GetLivePackets(utils.LINUX_LOOPBACK_ADAPTER_NAME, pChan, utils.FILTER_ALL)
	}()

	// 推送一个包
	go func() {
		time.Sleep(time.Second)
		// 创建一个报文
		// 顺序 src_mac string, dst_mac string, src_ip string, dst_ip string, src_port int, dst_port int, payload []byte
		newPacket := utils.CreateUDPPacket(
			"11:22:33:44:55:66",
			"66:55:44:33:22:11",
			"192.166.6.6",
			"192.188.8.8",
			60,
			80,
			[]byte{5, 8, 10})
		// 发送报文
		utils.SendPacket(utils.LINUX_LOOPBACK_ADAPTER_NAME, newPacket)
	}()

	// 打印抓到的包
	for p := range pChan {
		p.PrintAll()
	}
}
