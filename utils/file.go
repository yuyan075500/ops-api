package utils

import (
	"crypto/x509"
	"encoding/pem"
	"os"
)

// ReadFile 读取文件
func ReadFile(file string) ([]byte, error) {
	if f, err := os.Open(file); err != nil {
		return nil, err
	} else {
		content := make([]byte, 4096)
		if n, err := f.Read(content); err != nil {
			return nil, err
		} else {
			return content[:n], err
		}
	}
}

// ReadFileString 读取文件
func ReadFileString(file string) (string, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// LoadPublicKey 读取公钥
func LoadPublicKey() (interface{}, error) {
	// 读取公钥文件
	pubKeyPEM, err := os.ReadFile("config/certs/public.key")
	if err != nil {
		return nil, err
	}

	// 解析PEM块
	block, _ := pem.Decode(pubKeyPEM)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, err
	}

	// 解析公钥
	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return pubKey, nil
}
