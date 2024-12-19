package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
)

// Encrypt 字符串加密
func Encrypt(str string) (string, error) {
	file, err := ReadFile("/data/certs/public.key")
	if err != nil {
		file, err = ReadFile("config/certs/public.pem")
		if err != nil {
			return "", err
		}
	}

	// 解析公钥数据
	block, _ := pem.Decode(file)

	// 解析PEM格式的公钥
	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", err
	}

	// 根据公钥加密
	encryptedData, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey.(*rsa.PublicKey), []byte(str))
	return base64.RawURLEncoding.EncodeToString(encryptedData), nil
}
