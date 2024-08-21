package middleware

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/wonderivan/logger"
	"net/http"
	"ops-api/config"
	"ops-api/global"
	"os"
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
		mc, err := ValidateJWT(token)
		if err != nil {
			logger.Error("ERROR：", err)
			c.JSON(http.StatusOK, gin.H{
				"code": 90514,
				"msg":  err.Error(),
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

// ValidateJWT 校验Token
func ValidateJWT(token string) (mc *UserClaims, err error) {
	// 如果Token为空，则表示未认证
	if token == "" {
		return nil, errors.New("未认证")
	}

	// Token校验
	parts := strings.SplitN(token, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		return nil, errors.New("token无效")
	}

	// Token解析
	mc, err = ParseToken(parts[1])
	if err != nil {
		return nil, err
	}

	// 判断Token是否已注销
	val, err := global.RedisClient.Exists(parts[1]).Result()
	if err != nil {
		return nil, err
	}
	if val == 1 {
		return nil, errors.New("token无效")
	}

	return mc, nil
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
	//token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 返回Token字符串
	//return token.SignedString([]byte(config.Conf.JWT.Secret))

	// 使用RS256签名算法生成Token
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	// 读取私钥
	privateKeyData, err := os.ReadFile("config/certs/private.key")
	if err != nil {
		return "", err
	}

	// 解析私钥
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyData)
	if err != nil {
		return "", err
	}

	// 返回Token字符串（使用密钥签名）
	return token.SignedString(privateKey)
}

// ParseToken 解析Token
func ParseToken(tokenString string) (*UserClaims, error) {
	var mc = new(UserClaims)

	// 读取公钥
	publicKeyData, err := os.ReadFile("config/certs/public.key")
	if err != nil {
		return nil, err
	}

	// 解析公钥
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyData)
	if err != nil {
		return nil, err
	}

	token, err := jwt.ParseWithClaims(tokenString, mc, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New(fmt.Sprintf("unexpected signing method: %v", token.Header["alg"]))
		}
		return publicKey, nil
	})
	if err != nil {
		return nil, err
	}

	// 对token对象中的Claim进行类型断言，校验Token
	if token.Valid {
		return mc, nil
	}

	return nil, errors.New("token无效")
}
