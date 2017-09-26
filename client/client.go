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
	SERVER = "127.0.0.1:1026"
)

func handle(conn net.Conn) {
	defer conn.Close()
	fmt.Println("secret-tunnel-client: got a client from ", conn.RemoteAddr().String())
	remote, err := net.Dial("tcp", SERVER)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer remote.Close()

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

	clientReader := bufio.NewReader(conn)
	clientWriter := bufio.NewWriter(conn)

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		io.Copy(remoteWriter, clientReader)
	}()

	go func() {
		defer wg.Done()
		io.Copy(clientWriter, remoteReader)
	}()

	wg.Wait()
}

func main() {
	listen, err := net.Listen("tcp", ":1027")
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
