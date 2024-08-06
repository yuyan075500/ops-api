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

// Authorize 客户端授权
// @Summary 客户端授权
// @Description OAuth2.0认证相关接口
// @Tags OAuth2.0认证
// @Param Authorization header string true "Bearer 用户令牌"
// @Param authorize body service.Authorize true "授权请求参数"
// @Success 200 {string} json "{"code": 0, "msg": 授权成功, "redirect_uri": redirect_uri}"
// @Router /api/v1/oauth/authorize [post]
func (s *sso) Authorize(c *gin.Context) {

	var data = &service.Authorize{}

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
	callbackUrl, err := service.SSO.GetAuthorize(data, mc.ID)
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
