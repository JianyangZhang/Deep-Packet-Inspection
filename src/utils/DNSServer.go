package utils

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"

	dnss "github.com/miekg/dns"
)

var domainsToAddresses = map[string]string{
/*
	"baidu.com.": "1.2.3.4",
	"zhihu.com.": "104.198.14.52",
*/
}

func (this *handler) serveDNS(w dnss.ResponseWriter, r *dnss.Msg) {
	msg := dnss.Msg{}
	msg.SetReply(r)
	switch r.Question[0].Qtype {
	case dnss.TypeA:
		msg.Authoritative = true
		domain := msg.Question[0].Name
		/*
			address, ok := domainsToAddresses[domain]
			if ok {
				msg.Answer = append(msg.Answer, &dnss.A{
					Hdr: dnss.RR_Header{Name: domain, Rrtype: dnss.TypeA, Class: dnss.ClassINET, Ttl: 60},
					A:   net.ParseIP(address),
				})
			}
		*/
		msg.Answer = append(msg.Answer, &dnss.A{
			Hdr: dnss.RR_Header{Name: domain, Rrtype: dnss.TypeA, Class: dnss.ClassINET, Ttl: 60},
			A:   net.ParseIP("192.168.8.8"),
		})
	}
	w.WriteMsg(&msg)
}

type handler struct{}

/* 参数port: DNS服务器的端口号 (35被真实DNS服务占用) */
func StartDNSServer(port int) {
	srv := &dnss.Server{Addr: ":" + strconv.Itoa(port), Net: "udp"}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	go func() {
		for _ = range c {
			fmt.Println("DNS服务器关闭...")
			srv.Shutdown()
		}
	}()

	srv.Handler = &handler{}
	srv.NotifyStartedFunc = func() {
		fmt.Println("DNS服务器启动成功...")
		fmt.Println("测试执行命令: dig @localhost -p", port, "baidu.com")
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("DNS服务器启动失败... %s\n", err.Error())
	}
}
