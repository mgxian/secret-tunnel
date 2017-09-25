package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"sync"
)

var (
	// KEY 加密的密钥
	KEY    = "123456"
	SERVER = "127.0.0.1:1025"
)

func handle(conn net.Conn) {
	fmt.Println("server: got a client")
	proxy, err := net.Dial("tcp", SERVER)
	if err != nil {
		fmt.Println(err)
		return
	}

	//clientWriter, err := common.NewsecretWriter(conn, KEY)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	clientWriter := bufio.NewWriter(conn)

	// clientReader, err := common.NewsecretReader(conn, KEY)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	clientReader := bufio.NewReader(conn)

	defer proxy.Close()
	defer conn.Close()

	proxyReader := bufio.NewReader(proxy)
	proxyWriter := bufio.NewWriter(proxy)

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		io.Copy(proxyWriter, clientReader)
	}()

	go func() {
		defer wg.Done()
		io.Copy(clientWriter, proxyReader)
	}()

	wg.Wait()
}

func main() {
	listen, err := net.Listen("tcp", ":1026")
	if err != nil {
		fmt.Println(err)
		panic("listen error")
	}

	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go handle(conn)
	}
}
