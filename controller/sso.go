package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"net/http"
	"ops-api/middleware"
	"ops-api/service"
)

var SSO sso

type sso struct{}

// OAuthAuthorize 客户端授权
// @Summary 客户端授权
// @Description OAuth2.0认证相关接口
// @Tags OAuth2.0认证
// @Param Authorization header string true "Bearer 用户令牌"
// @Param authorize body service.OAuthAuthorize true "授权请求参数"
// @Success 200 {string} json "{"code": 0, "msg": 授权成功, "redirect_uri": redirect_uri}"
// @Router /api/v1/sso/oauth/authorize [post]
func (s *sso) OAuthAuthorize(c *gin.Context) {

	var data = &service.OAuthAuthorize{}

	// 请求参数绑定
	if err := c.ShouldBind(&data); err != nil {
		logger.Error("ERROR：" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 90400,
			"msg":  err.Error(),
		})
		return
	}

	// Token校验
	token := c.Request.Header.Get("Authorization")
	mc, err := middleware.ValidateJWT(token)
	if err != nil {
		logger.Error("ERROR：" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 90500,
			"msg":  err.Error(),
		})
		return
	}

	// 获取授权码
	callbackUrl, err := service.SSO.GetOAuthAuthorize(data, mc.ID)
	if err != nil {
		logger.Error("ERROR：" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 90500,
			"msg":  err.Error(),
		})
		return
	}

	// 返回客户端回调地址
	c.JSON(http.StatusOK, gin.H{
		"code":         0,
		"msg":          "授权成功",
		"redirect_uri": callbackUrl,
	})
}

// GetToken 客户端认证
// @Summary 客户端认证
// @Description OAuth2.0认证相关接口
// @Tags OAuth2.0认证
// @Param authorize body service.Token true "授权请求参数"
// @Success 200 {object} service.ResponseToken
// @Router /api/v1/oauth/token [post]
func (s *sso) GetToken(c *gin.Context) {

	var data = &service.Token{}

	// 请求参数绑定
	if err := c.Bind(&data); err != nil {
		logger.Error("ERROR：" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 90400,
			"msg":  err.Error(),
		})
		return
	}

	// 检查必需参数是否存在（要求客户端在获取Token时必须传入client_id和client_secret）
	if data.ClientId == "" && data.ClientSecret == "" {
		logger.Error("ERROR：" + "Missing required parameters")
		c.JSON(http.StatusOK, gin.H{
			"code": 90400,
			"msg":  "Missing required parameters",
		})
		return
	}

	token, err := service.SSO.GetToken(data)
	if err != nil {
		logger.Error("ERROR：" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 90500,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, token)
}

// GetUserInfo 获取用户信息
// @Summary 获取用户信息
// @Description OAuth2.0认证相关接口
// @Tags OAuth2.0认证
// @Param Authorization header string true "Bearer 用户令牌"
// @Success 200 {object} service.ResponseUserinfo
// @Router /api/v1/oauth/userinfo [get]
func (s *sso) GetUserInfo(c *gin.Context) {

	// 获取Token
	token := c.Request.Header.Get("Authorization")

	// 获取用户信息
	user, err := service.SSO.GetUserinfo(token)
	if err != nil {
		logger.Error("ERROR：" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 90500,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

// CASAuthorize 客户端授权
// @Summary 客户端授权
// @Description CAS3.0认证相关接口
// @Tags CAS3.0认证
// @Param Authorization header string true "Bearer 用户令牌"
// @Param authorize body service.CASAuthorize true "授权请求参数"
// @Success 200 {string} json "{"code": 0, "msg": 授权成功, "redirect_uri": redirect_uri}"
// @Router /api/v1/sso/cas/authorize [post]
func (s *sso) CASAuthorize(c *gin.Context) {

	var data = &service.CASAuthorize{}

	// 请求参数绑定
	if err := c.ShouldBind(&data); err != nil {
		logger.Error("ERROR：" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 90400,
			"msg":  err.Error(),
		})
		return
	}

	// Token校验
	token := c.Request.Header.Get("Authorization")
	mc, err := middleware.ValidateJWT(token)
	if err != nil {
		logger.Error("ERROR：" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 90500,
			"msg":  err.Error(),
		})
		return
	}

	// 获取票据
	callbackUrl, err := service.SSO.GetCASAuthorize(data, mc.ID, mc.Username)
	if err != nil {
		logger.Error("ERROR：" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 90500,
			"msg":  err.Error(),
		})
		return
	}

	// 返回客户端回调地址
	c.JSON(http.StatusOK, gin.H{
		"code":         0,
		"msg":          "授权成功",
		"redirect_uri": callbackUrl,
	})
}

// CASServiceValidate 票据校验
// @Summary 票据校验
// @Description CAS3.0认证相关接口
// @Tags CAS3.0认证
// @Param authorize body service.CASServiceValidate true "授权请求参数"
// @Produce xml
// @Success 200
// @Router /p3/serviceValidate [get]
func (s *sso) CASServiceValidate(c *gin.Context) {

	var data = &service.CASServiceValidate{}

	// 请求参数绑定
	if err := c.ShouldBind(&data); err != nil {
		logger.Error("ERROR：" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 90400,
			"msg":  err.Error(),
		})
		return
	}

	// 获取票据
	response, err := service.SSO.ServiceValidate(data)
	if err != nil {
		logger.Error("ERROR：" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 90500,
			"msg":  err.Error(),
		})
		return
	}

	// 设置响应头为XML格式
	c.Header("Content-Type", "application/xml")

	// 返回客户端回调地址（使用c.XML返回）
	c.XML(http.StatusOK, response)
}

// GetIdPMetadata 获取元数据
// @Summary 获取元数据
// @Description SAML2认证相关接口
// @Tags SAML2认证
// @Produce xml
// @Success 200
// @Router /api/v1/sso/saml/metadata [get]
func (s *sso) GetIdPMetadata(c *gin.Context) {

	// 获取票据
	response, err := service.SSO.GetIdPMetadata()
	if err != nil {
		logger.Error("ERROR：" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 90500,
			"msg":  err.Error(),
		})
		return
	}

	// 设置响应头为XML格式
	c.Header("Content-Type", "application/xml")

	// 返回客户端回调地址（使用c.Data返回，因为response非结构体，如果使用c.XML返回，在最外层会嵌套<string></string>）
	c.Data(http.StatusOK, response, []byte(response))
}

// SPAuthorize SP授权
// @Summary SP授权
// @Description SAML2认证相关接口
// @Tags SAML2认证
// @Success 200
// @Router /api/v1/sso/saml/authorize [post]
func (s *sso) SPAuthorize(c *gin.Context) {
	var data = &service.SAMLRequest{}

	// 请求参数绑定
	if err := c.ShouldBind(&data); err != nil {
		logger.Error("ERROR：" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 90400,
			"msg":  err.Error(),
		})
		return
	}

	// Token校验
	token := c.Request.Header.Get("Authorization")
	mc, err := middleware.ValidateJWT(token)
	if err != nil {
		logger.Error("ERROR：" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 90500,
			"msg":  err.Error(),
		})
		return
	}

	// authnRequest校验
	if err := service.SSO.SPAuthorize(data, mc.ID); err != nil {
		logger.Error("ERROR：" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 90500,
			"msg":  err.Error(),
		})
		return
	}
}

// ParseSPMetadata SP Metadata解析
// @Summary SP Metadata解析
// @Description SAML2认证相关接口
// @Tags SAML2认证
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param url body service.ParseSPMetadata true "授权请求参数"
// @Success 200 {string} json "{"code": 0, "msg": "解析成功", "data": nil}"
// @Router /api/v1/sso/saml/metadata [post]
func (s *site) ParseSPMetadata(c *gin.Context) {
	var data = &service.ParseSPMetadata{}

	if err := c.ShouldBind(&data); err != nil {
		logger.Error("ERROR：" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 90400,
			"msg":  err.Error(),
		})
		return
	}

	// 获取SP Metadata信息
	metadataInfo, err := service.SSO.ParseSPMetadata(data.SPMetadataURL)
	if err != nil {
		logger.Error("ERROR：" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 90500,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "解析成功",
		"data": metadataInfo,
	})
}
