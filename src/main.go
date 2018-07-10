package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"./utils"
)

func main() {
	go func() {
		utils.Receiver()
	}()
	utils.Sender()

/*--------------- 测试 ---------------*/
func PacketPushTest() {
	utils.GetDevices()

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
		utils.GetLivePackets(utils.CENTOS_ADAPTER_NAME, pChan, utils.TCP_80)
	}()

	// 推送一个包
	go func() {
		time.Sleep(time.Second)
		// 创建一个报文
		// 顺序 src_mac []byte, dst_mac []byte, src_ip []byte, dst_ip []byte, src_port int, dst_port int, payload []byte
		newPacket := utils.CreatePacket(
			[]byte{0x00, 0xff, 0x36, 0x87, 0xa0, 0x31},
			[]byte{0x00, 0xff, 0x37, 0x87, 0xa0, 0x31},
			[]byte{192, 168, 8, 8},
			[]byte{42, 236, 9, 26},
			999,
			80,
			[]byte{5, 8, 10})
		// 发送报文
		utils.SendPacket(utils.CENTOS_ADAPTER_NAME, newPacket)
	}()

	// 打印抓到的包
	for p := range pChan {
		p.PrintAll()
	}
}
