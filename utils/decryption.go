package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

var privateKey []byte

func readPrivateKeyFile(file string) {
	privateKey, _ = ReadFile(file)
}

// Decrypt 字符串解密
func Decrypt(cipherText []byte) ([]byte, error) {
	readPrivateKeyFile("config/certs/private.key")
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("公钥错误")
	}

	// 解析私钥
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	// 解密密文
	return rsa.DecryptPKCS1v15(rand.Reader, priv, cipherText)
}
