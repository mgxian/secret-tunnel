package main

import (
	"flag"

	"github.com/maogx8/secret-tunnel/client"
	"github.com/maogx8/secret-tunnel/server"
)

// Parameter 存储命令行参数的struct
type Parameter struct {
	host       string
	port       string
	remote     string
	server     string
	key        string
	serverMode bool
}

// NewParameter 创建存储命令行参数的struct
func NewParameter(host, port, remote, server, key string, serverMode bool) (*Parameter, error) {
	return &Parameter{
		host:       host,
		port:       port,
		remote:     remote,
		server:     server,
		key:        key,
		serverMode: serverMode,
	}, nil
}

func handleParameter() (*Parameter, error) {
	key := flag.String("k", "will", "encrypt string")
	host := flag.String("h", "", "listen address defualt '0.0.0.0'")
	port := flag.String("p", "1026", "listen port")
	remote := flag.String("r", "127.0.0.1:1025", "remote address for proxy")
	server := flag.String("s", "127.0.0.1:1026", "secret-tunnel-server address")
	serverMode := flag.Bool("server", false, "secret-tunnel start on server mode default true")

	flag.Parse()

	return NewParameter(*host, *port, *remote, *server, *key, *serverMode)
}

func main() {
	p, _ := handleParameter()

	listenAddress := p.host + ":" + p.port
	if p.serverMode {
		server.Run(listenAddress, p.key, p.remote)
	}

	client.Run(listenAddress, p.key, p.server)
}
