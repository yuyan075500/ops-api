package middleware

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/gin-gonic/gin"
	"ops-api/global"
	"strings"
)

// CasBinInit 权限初始化
func CasBinInit() {

	// 初始化CasBin适配器
	adapter, err := gormadapter.NewAdapterByDB(global.MySQLClient)
	if err != nil {
		fmt.Println("初始化CasBin适配器失败：" + err.Error())
		return
	}

	// 初始化CasBin执行器
	enforcer, err := casbin.NewEnforcer("config/rbac_model.conf", adapter)
	if err != nil {
		fmt.Println("初始化CasBin执行器失败：" + err.Error())
		return
	}

	// 加载规则
	err = enforcer.LoadPolicy()
	if err != nil {
		fmt.Println("加载CasBin规则失败：" + err.Error())
		return
	}

	global.CasBinServer = enforcer
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
			"/login",
			"/health",
			"/swagger/",
		}
		for _, item := range ignorePath {
			if strings.HasPrefix(item, path) {
				c.Next()
				return
			}
		}

		// 检查用户权限
		ok, err := global.CasBinServer.Enforce(username, path, method)
		if err != nil {
			c.JSON(200, gin.H{
				"code": 500,
				"msg":  "检查权限失败：" + err.Error(),
			})
			c.Abort()
			return
		} else if !ok {
			c.JSON(403, gin.H{
				"code": 90403,
				"msg":  "您无权访问该资源，请联系管理员",
			})
			c.Abort()
			return
		} else {
			c.Next()
		}
	}
}
