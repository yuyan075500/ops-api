package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

var publicKey []byte

func readPublicKeyFile(file string) {
	publicKey, _ = ReadFile(file)
}

// Encrypt 字符串加密
func Encrypt(str []byte) ([]byte, error) {
	readPublicKeyFile("config/certs/public.cert")
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, errors.New("公钥错误")
	}

	// 解析公钥
	publicInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	// 类型断言
	pub := publicInterface.(*rsa.PublicKey)

	// 加密明明文
	return rsa.EncryptPKCS1v15(rand.Reader, pub, str)
}
