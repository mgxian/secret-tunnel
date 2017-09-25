package common

import (
	"bufio"
	"crypto/md5"
	"crypto/rc4"
	"errors"
	"fmt"
	"io"
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
func NewsecretReader(r io.Reader, key string) (*secretReader, error) {
	md5Byte := md5.New().Sum([]byte(key))
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
func NewsecretWriter(w io.Writer, key string) (*secretWriter, error) {
	md5Byte := md5.New().Sum([]byte(key))
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
