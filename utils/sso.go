package utils

import (
	"bytes"
	"compress/flate"
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"encoding/xml"
	"errors"
	"io"
	"net/http"
	"os"
)

// EntityDescriptor SP Metadata中的数据绑定结构体
type EntityDescriptor struct {
	XMLName         xml.Name        `xml:"EntityDescriptor"`
	EntityID        string          `xml:"entityID,attr"`
	SPSSODescriptor SPSSODescriptor `xml:"SPSSODescriptor"`
}
type SPSSODescriptor struct {
	KeyDescriptors []KeyDescriptor `xml:"KeyDescriptor"`
}
type KeyDescriptor struct {
	Use     string  `xml:"use,attr"`
	KeyInfo KeyInfo `xml:"KeyInfo"`
}
type KeyInfo struct {
	X509Data X509Data `xml:"X509Data"`
}
type X509Data struct {
	X509Certificate string `xml:"X509Certificate"`
}

// ParseSPMetadata SP Metadata数据解析
func ParseSPMetadata(metadataUrl string) (*EntityDescriptor, error) {

	var entityDescriptor = &EntityDescriptor{}

	// 请求SP Metadata地址
	resp, err := http.Get(metadataUrl)
	if err != nil {
		return nil, err
	}

	// 获取请求到的数据
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 将Metadata数据绑定到结构体
	err = xml.Unmarshal(data, &entityDescriptor)
	if err != nil {
		return nil, err
	}

	return entityDescriptor, nil
}

// VerifySignature SP签名验证
// 请求参数：
//   - samlRequest：原始的SAMLRequest
//   - certificate：base64编辑的证书
//   - signature：签名
//   - sigAlg：签名算法
func VerifySignature(samlRequest bytes.Buffer, certificate, signature, sigAlg string) error {

	// 证书解码（证书在SP Metadata中一般都是base64编码的，需要先解码）
	certBytes, err := base64.StdEncoding.DecodeString(certificate)
	if err != nil {
		return err
	}

	// 证书解析
	cert, err := x509.ParseCertificate(certBytes)
	if err != nil {
		return err
	}

	// 签名解码（签名在SP Metadata中一般都是base64编码的，需要先解码）
	signatureBytes, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return err
	}

	// 选择签名算法
	var h crypto.Hash
	switch sigAlg {
	case "http://www.w3.org/2001/04/xmldsig-more#rsa-sha1":
		h = crypto.SHA1
	case "http://www.w3.org/2001/04/xmldsig-more#rsa-sha256":
		h = crypto.SHA256
	case "http://www.w3.org/2001/04/xmldsig-more#rsa-sha384":
		h = crypto.SHA384
	case "http://www.w3.org/2001/04/xmldsig-more#rsa-sha512":
		h = crypto.SHA512
	default:
		return errors.New("不支持的签名算法")
	}

	// 创建哈希对象
	hashFunc := h.New()
	hashFunc.Write(samlRequest.Bytes())
	hashed := hashFunc.Sum(nil)

	// 验证签名
	rsaPublicKey, ok := cert.PublicKey.(*rsa.PublicKey)
	if !ok {
		return errors.New("证书不包含RAS公钥")
	}

	if err = rsa.VerifyPKCS1v15(rsaPublicKey, h, hashed, signatureBytes); err != nil {
		return errors.New("签名验证失败")
	}

	return nil
}

// LoadIdpCertificate 获取IDP证书
func LoadIdpCertificate() (*x509.Certificate, error) {

	// 读取证书
	certData, err := os.ReadFile("config/certs/certificate.crt")
	if err != nil {
		return nil, err
	}

	// 解码PEM格式证书
	block, _ := pem.Decode(certData)
	if block == nil {
		return nil, errors.New("证书解码失败")
	}

	// 解析证书
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}

	return cert, nil
}

// LoadIdpPrivateKey 获取IDP私钥
func LoadIdpPrivateKey() (string, error) {

	// 读取私钥文件
	privateKeyBytes, err := os.ReadFile("config/certs/private.key")
	if err != nil {
		return "", err
	}

	// 解码PEM格式证书
	block, _ := pem.Decode(privateKeyBytes)
	if block == nil {
		return "", errors.New("私钥解码失败")
	}

	// 解析私钥
	privateKeyInterface, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}

	// 私钥转换
	privateKey, ok := privateKeyInterface.(*rsa.PrivateKey)
	if !ok {
		return "", err
	}

	// 将私钥编码为PEM格式的字符串
	privateKeyPem := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	return string(privateKeyPem), nil
}

// ParseSAMLRequest SP SAMLRequest请求解析
func ParseSAMLRequest(samlRequest string) (data bytes.Buffer, err error) {

	var decompressed bytes.Buffer

	// Base64解码
	decoded, err := base64.StdEncoding.DecodeString(samlRequest)
	if err != nil {
		return decompressed, err
	}

	// DEFLATE解压缩，并读取里面的数据
	reader := flate.NewReader(bytes.NewReader(decoded))

	_, err = io.Copy(&decompressed, reader)
	if err != nil {
		return decompressed, err
	}

	return decompressed, nil
}
