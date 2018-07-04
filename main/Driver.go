package main

import (
	"fmt"
	"time"
)

func main() {
	c := make(chan bool, 8)
	go func() {
		fmt.Println("Go Go Go")
		c <- true
		c <- true
		c <- true
		c <- true
		close(c)
	}()
	for v := range c {
		time.Sleep(time.Second)
		fmt.Println(v)
	}
}
