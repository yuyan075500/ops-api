package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
)

var publicKey []byte

func readPublicKeyFile(file string) {
	publicKey, _ = ReadFile(file)
}

// Encrypt 字符串加密
func Encrypt(str string) (string, error) {
	readPublicKeyFile("config/certs/public.pem")

	// 解析公钥数据
	block, _ := pem.Decode(publicKey)

	// 解析PEM格式的公钥
	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		fmt.Println(err)
		return "", errors.New("公钥解析错误")
	}

	// 根据公钥加密
	encryptedData, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey.(*rsa.PublicKey), []byte(str))
	return base64.StdEncoding.EncodeToString(encryptedData), nil
}
