package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/maogx8/secret-tunnel/common"
)

var (
	// KEY 加密的密钥
	KEY    = "123456"
	SERVER = "127.0.0.1:1025"
)

func handle(conn net.Conn) {
	defer conn.Close()
	fmt.Println("secret-tunnel-server: got a client from ", conn.RemoteAddr().String())
	proxy, err := net.Dial("tcp", SERVER)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer proxy.Close()

	clientWriter, err := common.NewSecretWriter(conn, KEY)
	if err != nil {
		fmt.Println(err)
		return
	}

	clientReader, err := common.NewSecretReader(conn, KEY)
	if err != nil {
		fmt.Println(err)
		return
	}

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

	defer listen.Close()

	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go handle(conn)
	}
}
