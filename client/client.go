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
	SERVER = "47.52.79.208:1026"
)

func handle(conn net.Conn) {
	fmt.Println("client: got a client")
	remote, err := net.Dial("tcp", SERVER)
	if err != nil {
		fmt.Println(err)
		return
	}

	remoteWriter, err := common.NewSecretWriter(remote, KEY)
	if err != nil {
		fmt.Println(err)
		return
	}

	remoteReader, err := common.NewSecretReader(remote, KEY)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer remote.Close()
	defer conn.Close()

	clientReader := bufio.NewReader(conn)
	clientWriter := bufio.NewWriter(conn)

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		io.Copy(remoteWriter, clientReader)
		conn.Close()
	}()

	go func() {
		defer wg.Done()
		io.Copy(clientWriter, remoteReader)
		remote.Close()
	}()

	wg.Wait()
}

func main() {
	listen, err := net.Listen("tcp", ":1027")
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
