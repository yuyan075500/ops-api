package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/wonderivan/logger"
	"net/http"
	"ops-api/config"
	"ops-api/global"
	"strings"
	"time"
)

// Login 保存不需要验证的URL结构体信息
type Login struct {
	paths []string
}

// UserClaims 保存需要保存到JWT中的信息结构体
type UserClaims struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func LoginBuilder() *Login {
	return &Login{}
}

// IgnorePaths 保存不需要认证的URL到结构体
func (l *Login) IgnorePaths(path string) *Login {
	l.paths = append(l.paths, path)
	return l
}

func (l *Login) Build() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 不需要认证的路径，支持前缀匹配
		for _, path := range l.paths {
			if c.Request.URL.Path == path || strings.HasPrefix(c.Request.URL.Path, path) {
				return
			}
		}

		// 获取Token
		token := c.Request.Header.Get("Authorization")

		// 未认证
		if token == "" {
			c.JSON(http.StatusOK, gin.H{
				"code": 90514,
				"msg":  "未认证",
			})
			c.Abort()
			return
		}

		// Token校验
		parts := strings.SplitN(token, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusOK, gin.H{
				"code": 90514,
				"msg":  "Token无效",
			})
			c.Abort()
			return
		}

		// Token解析
		mc, err := ParseToken(parts[1])
		if err != nil {
			logger.Error("ERROR：", err)
			c.JSON(http.StatusOK, gin.H{
				"code": 90514,
				"msg":  err.Error(),
			})
			c.Abort()
			return
		}

		// 判断Token是否已注销
		val, err := global.RedisClient.Exists(parts[1]).Result()
		if err != nil {
			logger.Error("ERROR：", err)
			c.JSON(http.StatusOK, gin.H{
				"code": 90514,
				"msg":  err.Error(),
			})
			c.Abort()
			return
		}
		if val == 1 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 90514,
				"msg":  "token无效",
			})
			c.Abort()
			return
		}

		// 将当前请求的用户信息保存到请求的上下文c
		c.Set("id", mc.ID)
		c.Set("name", mc.Name)
		c.Set("username", mc.Username)
		// 后续的处理函数可以用过c.Get("username")来获取当前请求的用户信息
		c.Next()
	}
}

// GenerateJWT 生成Token
func GenerateJWT(id uint, name, username string) (string, error) {
	claims := UserClaims{
		id,
		name,
		username,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(config.Conf.JWT.Expires) * time.Hour)), // 过期时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),                                                         // 签发时间
			NotBefore: jwt.NewNumericDate(time.Now()),                                                         // 生效时间
		},
	}

	// 使用HS256签名算法生成Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 返回Token字符串
	return token.SignedString([]byte(config.Conf.JWT.Secret))
}

// ParseToken 解析Token
func ParseToken(tokenString string) (*UserClaims, error) {
	var mc = new(UserClaims)
	token, err := jwt.ParseWithClaims(tokenString, mc, func(token *jwt.Token) (i interface{}, err error) {
		return []byte(config.Conf.JWT.Secret), nil
	})
	if err != nil {
		return nil, err
	}

	// 对token对象中的Claim进行类型断言， 校验Token
	if token.Valid {
		return mc, nil
	}

	return nil, errors.New("token无效")
}
