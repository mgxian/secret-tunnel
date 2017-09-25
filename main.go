package main

import (
	"bufio"
	"crypto/md5"
	"crypto/rc4"
	"errors"
	"fmt"
	"io"
	"net"
)

var (
	// KEY 加密的密钥
	KEY = "123456"
)

type secretReader struct {
	r      *bufio.Reader
	cipher *rc4.Cipher
}

type secretWriter struct {
	w      *bufio.Writer
	cipher *rc4.Cipher
}

// NewsecretReader 创建加密读
func NewsecretReader(r io.Reader) (*secretReader, error) {
	md5Byte := md5.New().Sum([]byte(KEY))
	cipher, err := rc4.NewCipher(md5Byte)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("create cipher fail")
	}

	return &secretReader{
		r:      bufio.NewReader(r),
		cipher: cipher,
	}, nil
}

// NewsecretWriter 创建加密写
func NewsecretWriter(w io.Writer) (*secretWriter, error) {
	md5Byte := md5.New().Sum([]byte(KEY))
	cipher, err := rc4.NewCipher(md5Byte)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("create cipher fail")
	}
	return &secretWriter{
		w:      bufio.NewWriter(w),
		cipher: cipher,
	}, nil
}

func (r *secretReader) Read(b []byte) (n int, err error) {
	n, err = r.r.Read(b)
	if err != nil {
		fmt.Println(err)
		return
	}
	r.cipher.XORKeyStream(b, b)
	return
}

func (w *secretWriter) Write(p []byte) (n int, err error) {
	w.cipher.XORKeyStream(p, p)
	n, err = w.w.Write(p)
	if err != nil {
		fmt.Println(err)
		return
	}
	w.w.Flush()
	return
}

func handle(conn net.Conn) {
	reader, err := NewsecretReader(conn)
	if err != nil {
		fmt.Println(err)
		return
	}

	p := make([]byte, 1024)
	n, err := reader.Read(p)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(p[:n]))
	return
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
