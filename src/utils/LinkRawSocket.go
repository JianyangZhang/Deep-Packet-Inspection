package utils

/*
	不同操作系统下的syscall库不一样，windows版本无法编译，centOS测试通过
*/
/*
import (
	"fmt"
	"log"
	"net"
	"os"
	"syscall"

	"golang.org/x/net/ipv4"
)

func Receiver() {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_ICMP) // fd是socket的文件描述符
	if err != nil {
		panic(err)
	}
	f := os.NewFile(uintptr(fd), fmt.Sprintf("fd %d", fd)) // f为fd的文件指针

	for {
		buf := make([]byte, 1024)
		numRead, err := f.Read(buf)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("% x\n", buf[:numRead]) // x，将整数转换成十六进制表示，并将其格式化到指定位置
	}
}

func Sender() {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_RAW)
	if err != nil {
		panic(err)
	}
	addr := syscall.SockaddrInet4{
		Port: 0,
		Addr: [4]byte{127, 8, 8, 8},
	}
	packet := createPacket()
	err = syscall.Sendto(fd, packet, 0, &addr)
	if err != nil {
		log.Fatal("Sendto:", err)
	} else {
		fmt.Println("Raw Packet 已发送")
	}
}

func createPacket() []byte {
	icmp := []byte{
		8, // type: echo request
		0, // code: not used by echo request
		0, // checksum (16 bit), we fill in below
		0,
		0, // identifier (16 bit). zero allowed.
		0,
		0, // sequence number (16 bit). zero allowed.
		0,
		0xC0, // Optional data. ping puts time packet sent here
		0xDE,
	}
	cs := csum(icmp)
	icmp[2] = byte(cs)
	icmp[3] = byte(cs >> 8)

	h := ipv4.Header{
		Version:  4,
		Len:      20,
		TotalLen: 20 + 10, // 20 bytes for IP, 10 for ICMP
		TTL:      64,
		Protocol: 1, // ICMP
		Dst:      net.IPv4(127, 5, 5, 5),
		// ID, Src and Checksum will be set for us by the kernel
	}
	out, err := h.Marshal()
	if err != nil {
		log.Fatal(err)
	}

	return append(out, icmp...)
}

func csum(b []byte) uint16 {
	var s uint32
	for i := 0; i < len(b); i += 2 {
		s += uint32(b[i+1])<<8 | uint32(b[i])
	}
	// add back the carry
	s = s>>16 + s&0xffff
	s = s + s>>16
	return uint16(^s)
}
*/