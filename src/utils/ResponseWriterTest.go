/**
HTTP/2 的服务器推送功能测试，暂时没什么用
*/

package utils

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
)

var image []byte
var SAMPLE_PORT string = ":8080"

/* 准备将要推送的图片 */
func checkPushContent() {
	var err error
	if image, err = ioutil.ReadFile("./image.png"); err != nil {
		panic(err)
	}
}

/* 处理 / 的推送 */
func handlerHtml(w http.ResponseWriter, r *http.Request) {
	pusher, ok := w.(http.Pusher)
	if ok {
		fmt.Println("Push /image")
		pusher.Push("/image", nil)
	}
	w.Header().Add("Content-Type", "text/html")
	fmt.Fprintf(w, `<html><body><img src="/image"></body></html>`)
}

/* 处理 /image 的推送 */
func handlerImage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/png")
	w.Write(image)
}

func ResponseWriterTest() {
	checkPushContent()
	http.HandleFunc("/", handlerHtml)
	http.HandleFunc("/image", handlerImage)
	fmt.Println("监听端口", SAMPLE_PORT)

	server := &http.Server{
		Addr:    SAMPLE_PORT,
		Handler: http.DefaultServeMux,
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	go func() {
		for _ = range c {
			fmt.Println("关闭端口", SAMPLE_PORT)
			server.Close()
		}
	}()

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
