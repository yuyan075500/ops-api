package service

import (
	"context"
	"errors"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkauthen "github.com/larksuite/oapi-sdk-go/v3/service/authen/v1"
	"ops-api/config"
)

// FeishuClient 飞书SDK客户端
type FeishuClient struct {
	client *lark.Client
}

// NewFeishuClient 飞书客户端初始化，参考文档：https://github.com/larksuite/oapi-sdk-go
func NewFeishuClient() (*FeishuClient, error) {

	// 读取配置信息
	appId := config.Conf.Feishu.AppId         // 自建应用的AppId
	appSecret := config.Conf.Feishu.AppSecret // 自建应用的appSecret

	// 创建客户端
	client := lark.NewClient(appId, appSecret)

	return &FeishuClient{
		client: client,
	}, nil
}

// GetUserAccessToken 获取访问Token
func (c *FeishuClient) GetUserAccessToken(code string) (*larkauthen.CreateOidcAccessTokenResp, error) {

	// 创建请求对象
	req := larkauthen.NewCreateOidcAccessTokenReqBuilder().
		Body(larkauthen.NewCreateOidcAccessTokenReqBodyBuilder().
			GrantType(`authorization_code`).
			Code(code).
			Build()).
		Build()

	// 发起请求
	resp, err := c.client.Authen.OidcAccessToken.Create(context.Background(), req)
	// 处理错误
	if err != nil {
		return nil, err
	}
	// 服务端错误处理
	if !resp.Success() {
		return nil, errors.New(resp.Msg)
	}

	return resp, nil
}

// GetUserInfo 获取用户信息
func (c *FeishuClient) GetUserInfo(userAccessToken string) (*larkauthen.GetUserInfoResp, error) {

	// 发起请求
	resp, err := c.client.Authen.UserInfo.Get(context.Background(), larkcore.WithUserAccessToken(userAccessToken))

	// 处理错误
	if err != nil {
		return nil, err
	}
	// 服务端错误处理
	if !resp.Success() {
		return nil, errors.New(resp.Msg)
	}

	return resp, nil
}
