package server

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/maogx8/secret-tunnel/common"
)

func handle(conn net.Conn, key, remote string) {
	defer conn.Close()
	fmt.Println("secret-tunnel-server: got a client from ", conn.RemoteAddr().String())
	proxy, err := net.Dial("tcp", remote)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer proxy.Close()

	clientWriter, err := common.NewSecretWriter(conn, key)
	if err != nil {
		fmt.Println(err)
		return
	}

	clientReader, err := common.NewSecretReader(conn, key)
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
		proxy.Close()
		fmt.Println("client say bye bye")
	}()

	go func() {
		defer wg.Done()
		io.Copy(clientWriter, proxyReader)
		conn.Close()
		fmt.Println("proxy say bye bye")
	}()

	wg.Wait()
}

// Run start secret-tunnel server
func Run(listenAddress, key, remote string) {
	listen, err := net.Listen("tcp", listenAddress)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("secret-tunnel-server listening on", listenAddress)
	defer listen.Close()

	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go handle(conn, key, remote)
	}
}
