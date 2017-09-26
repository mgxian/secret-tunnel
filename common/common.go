package common

import (
	"crypto/md5"
	"crypto/rc4"
	"errors"
	"fmt"
	"io"
)

// SecretReader 加密读
type SecretReader struct {
	r      io.Reader
	cipher *rc4.Cipher
}

// SecretWriter 加密写
type SecretWriter struct {
	w      io.Writer
	cipher *rc4.Cipher
}

// NewSecretReader 创建加密读
func NewSecretReader(r io.Reader, key string) (*SecretReader, error) {
	md5Byte := md5.New().Sum([]byte(key))
	cipher, err := rc4.NewCipher(md5Byte)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("create cipher fail")
	}

	return &SecretReader{
		r:      r,
		cipher: cipher,
	}, nil
}

// NewSecretWriter 创建加密写
func NewSecretWriter(w io.Writer, key string) (*SecretWriter, error) {
	md5Byte := md5.New().Sum([]byte(key))
	cipher, err := rc4.NewCipher(md5Byte)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("create cipher fail")
	}
	return &SecretWriter{
		w:      w,
		cipher: cipher,
	}, nil
}

func (r *SecretReader) Read(b []byte) (n int, err error) {
	n, err = r.r.Read(b)
	if err != nil {
		fmt.Println(err)
		return
	}
	r.cipher.XORKeyStream(b, b[:n])
	return
}

func (w *SecretWriter) Write(p []byte) (n int, err error) {
	w.cipher.XORKeyStream(p, p)
	n, err = w.w.Write(p)
	if err != nil {
		fmt.Println(err)
		return
	}
	return
}
