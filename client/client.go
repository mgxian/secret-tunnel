package client

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/maogx8/secret-tunnel/common"
)

func handle(conn net.Conn, key, server string) {
	defer conn.Close()
	fmt.Println("secret-tunnel-client: got a client from ", conn.RemoteAddr().String())
	remote, err := net.Dial("tcp", server)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer remote.Close()

	remoteWriter, err := common.NewSecretWriter(remote, key)
	if err != nil {
		fmt.Println(err)
		return
	}

	remoteReader, err := common.NewSecretReader(remote, key)
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
		remote.Close()
		fmt.Println("client say bye bye")
	}()

	go func() {
		defer wg.Done()
		io.Copy(clientWriter, remoteReader)
		conn.Close()
		fmt.Println("server say bye bye")
	}()

	wg.Wait()
}

// Run 启动客户端
func Run(listenAddress, key, server string) {
	listen, err := net.Listen("tcp", listenAddress)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("secret-tunnel-client listening on", listenAddress)

	defer listen.Close()

	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go handle(conn, key, server)
	}
}
