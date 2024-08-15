package middleware

import (
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"net/http"
	"ops-api/global"
	"ops-api/model"
	"strings"
)

// CasBinInit 权限初始化
func CasBinInit() error {

	// 初始化CasBin适配器
	adapter, err := gormadapter.NewAdapterByDBWithCustomTable(global.MySQLClient, &model.CasbinRule{}, "casbin_rules")
	if err != nil {
		return err
	}

	// 初始化CasBin执行器
	enforcer, err := casbin.NewEnforcer("config/rbac_model.conf", adapter)
	if err != nil {
		return err
	}

	// 加载规则
	err = enforcer.LoadPolicy()
	if err != nil {
		return err
	}

	global.CasBinServer = enforcer

	return nil
}

// PermissionCheck 用户权限检查
func PermissionCheck() gin.HandlerFunc {
	return func(c *gin.Context) {

		// 用户名
		username, _ := c.Get("username")

		// 请求路径
		path := c.Request.URL.Path

		// 请求访问
		method := c.Request.Method

		// 排除不需要权限验证的接口，支持前缀匹配
		ignorePath := []string{
			"/api/auth/login",                 // 登录接口
			"/api/auth/logout",                // 注销接口
			"/health",                         // 预留健身检查接口
			"/api/v1/user/info",               // 用户登录成功后获取用户信息接口
			"/api/v1/user/avatarUpload",       // 用户头像上传接口
			"/swagger/",                       // swagger接口
			"/api/v1/sms/callback",            // 短信回调接口
			"/api/v1/sms/reset_password_code", // 获取重置密码验证码
			"/api/v1/reset_password",          // 密码自助重置接口
			"/api/v1/user/mfa_qrcode",         // 获取MFA二维码
			"/api/v1/user/mfa_auth",           // MFA认证
			"/api/v1/site/logoUpload",         // 站点图片上传
			"/api/v1/site/guide",              // 获取导航站点信息
			"/api/v1/sso/oauth/authorize",     // OAuth2.0客户端授权
			"/api/v1/sso/cas/authorize",       // CAS3.0客户端授权
			"/api/v1/sso/token",               // OAuth2.0客户端认证
			"/api/v1/sso/userinfo",            // OAuth2.0客户端获取用户信息
			"/p3/serviceValidate",             // CAS3.0票据校验
			"/api/v1/sso/saml/metadata",       // SAML2元数据
			"/api/v1/sso/saml/authorize",      // SP授权
			"/api/v1/sso/saml/metadata",       // SP元数据解析
		}
		for _, item := range ignorePath {
			if strings.HasPrefix(path, item) {
				c.Next()
				return
			}
		}

		// 检查用户权限
		ok, err := global.CasBinServer.Enforce(username, path, method)
		if err != nil {
			logger.Error("ERROR：", err.Error())
			c.JSON(http.StatusOK, gin.H{
				"code": 90500,
				"msg":  err.Error(),
			})
			c.Abort()
			return
		} else if !ok {
			c.JSON(http.StatusOK, gin.H{
				"code": 90403,
				"msg":  "该资源您无权访问",
			})
			c.Abort()
			return
		} else {
			c.Next()
		}
	}
}
