package utils

import (
	"bytes"
	"compress/flate"
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

// AuthnRequest SAMLRequest数据绑定结构体
type AuthnRequest struct {
	XMLName                     xml.Name     `xml:"urn:oasis:names:tc:SAML:2.0:protocol AuthnRequest"`
	AssertionConsumerServiceURL string       `xml:"AssertionConsumerServiceURL,attr"`
	Destination                 string       `xml:"Destination,attr"`
	ID                          string       `xml:"ID,attr"`
	IssueInstant                string       `xml:"IssueInstant,attr"`
	ProtocolBinding             string       `xml:"ProtocolBinding,attr"`
	Version                     string       `xml:"Version,attr"`
	Issuer                      Issuer       `xml:"urn:oasis:names:tc:SAML:2.0:assertion Issuer"`
	NameIDPolicy                NameIDPolicy `xml:"NameIDPolicy"`
	OriginalString              string       // 保存原始字符串,用于签名验证
}
type Issuer struct {
	Value string `xml:",chardata"`
}
type NameIDPolicy struct {
	AllowCreate     string `xml:"AllowCreate,attr"`
	Format          string `xml:"Format,attr"`
	SPNameQualifier string `xml:"SPNameQualifier,attr"`
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

// LoadIdpKey 获取IDP私钥
//func LoadIdpKey() (interface{}, error) {
//
//	// 读取证书
//	buf, err := os.ReadFile("config/certs/private.key")
//	if err != nil {
//		return nil, err
//	}
//
//	// 解码PEM格式证书
//	block, _ := pem.Decode(buf)
//	if block == nil {
//		return nil, errors.New("证书解码失败")
//	}
//
//	// 解析证书
//	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
//	if err != nil {
//		return nil, err
//	}
//
//	return privateKey, nil
//}

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

// Decompress 解压缩
func Decompress(in []byte) ([]byte, error) {
	buf := new(bytes.Buffer)
	decompressor := flate.NewReader(bytes.NewReader(in))
	if _, err := io.Copy(buf, decompressor); err != nil {
		return nil, err
	}
	if err := decompressor.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// ParseSAMLRequest SP SAMLRequest请求解析
func ParseSAMLRequest(samlRequest string) (data *AuthnRequest, err error) {

	var authnRequest AuthnRequest

	// Base64解码
	compressedXML, err := base64.StdEncoding.DecodeString(samlRequest)
	if err != nil {
		return nil, err
	}

	// 解压缩
	bXML, err := Decompress(compressedXML)
	if err != nil {
		return nil, err
	}

	// 数据保存
	if err := xml.Unmarshal(bXML, &authnRequest); err != nil {
		return nil, err
	}

	// 保留原始字符串用于签名验证
	authnRequest.OriginalString = string(bXML)

	return &authnRequest, nil
}
