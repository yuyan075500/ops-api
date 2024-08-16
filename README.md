# 目录说明
* config：全局配置，如监听地址、数据库连接配置等。
* controller：定义路由规则以及接口同业入参和响应。
* service：处理接口的业务逻辑。
* dao：数据库操作。
* model：定义数据库表信息。
* db：数据库、缓存、对象存储客户端初始化。
* middleware：中间件层，添加全局逻辑处理，如跨域、JWT认证等。
* utils：工具目录，定义常用工具，如Token解析，文件操作等。
# 项目依赖（必须）
* MySQL
* Redis：主要用于存放黑名单（注销用户）Token以及校验码缓存
* Minio：用于存储用户头像
# 功能概览
## 基础功能
* RBAC权限管理（基于CasBin实现）
* 统一站点管理（SSO认证）
## 认证相关
* 双因素认证（支持Google Authenticator、华为云、阿里云）
* 单点登录（支持CAS 3.0、OAuth 2.0和SAML2）
* 钉钉扫码登录（需要配置钉钉应用）
* AD认证
## 其它
* 短信验证码（仅支持华为云）
* Swagger接口文档
* 用户密码自助更改
* 前端水印
# 后端返回状态码说明
* 0：请求成功
* 90400：请求参数错误
* 90401：认证失败
* 90403：拒绝访问
* 90404：访问的对象不存在
* 90500：其它错误
* 90514：Token过期或无效
# 项目部署
在部署前需要确保项目依赖必须项已全部准备完成，如：MySQL、Redis、Minio。
## 项目配置文件
项目配置文件路径为`config/config.yaml`，如果没有则创建，配置说明如下：
```yaml
server: "0.0.0.0:8000"
mysql:
  host: "127.0.0.1"
  port: 3306
  db: "ops"
  user: "root"
  password: ""
  maxIdleConns: 10
  maxOpenConns: 100
  maxLifeTime: 30
redis:
  host: "127.0.0.1:6379"
  password: ""
  db: 0
jwt:
  secret: "swfqezjzoqssvjck"
  expires: 6
mfa:
  enable: false
  issuer: "统一认证平台"
oss:
  endpoint: ""
  accessKey: ""
  secretKey: ""
  bucketName: ""
  ssl: true
ldap:
  host: ""
  bindUserDN: ""
  bindUserPassword: ""
  searchDN: ""
sms:
  url: "https://smsapi.cn-north-4.myhuaweicloud.com:443/sms/batchSendDiffSms/v1"
  appKey: ""
  appSecret: ""
  callbackUrl: "https://ops-test.50yc.cn/api/v1/sms/callback"
  verificationCode:
    sender: ""
    templateId: ""
    signature: ""
mail:
  smtpHost: ""
  smtpPort: 587
  from: ""
  password: ""
swagger: true
```
* server：服务监听的地址和端口，必须
* mysql：MySQL数据库相关配置，必须
* redis：Redis相关配置，必须
* jwt：JWT相关配置，注意：secret请务必更改，必须
* mfa：双因素认证相关配置，issuer为APP扫码后显示的名称，必须
* oss：Minio对象存储相关配置，主要存储用户头像和资产图片，必须
* ldap：LDAP相关配置，用于AD域用户登录，可选
* sms：短信相关配置，仅支持华为云，用户自主重置密码需要使用，可选
* mail：邮件相关配置，用于用户自助密码重置，可选
* swagger：是否开启Swagger，生产环境请忽开启，必须
## 导入初始化数据
初始化数据SQL文件位于`deploy/data.sql`。
## 更新IP地址库文件
地址库文件用于分析用户登录城市，文件位于`config/GeoLite2-City.mmdb`，本地址库截止更新日志为2024-07-23，如果有需要可从官方获取最新文件替换即可。
## 创建管理员账号
管理员账号需要项目运行后创建，具体操作步骤如下：
## 项目其它配置
### 短信
短信功能用于用户自助密码重置，目前仅支持华为云，使用的模板格式如下：
```
${1}您好，您的校验码为：${2}，校验码在${3}分钟内有效，保管好校验码，请勿泄漏！
```
需要确保模板中包含三个变量，分别代表用户名、校验码和校验码有效时间，其它文字可以自定义。
### 密钥
密钥用于数据库敏感字段加密解密，其存放路径为：
```shell
config/certs/
```
项目部署时需要新生成相关密钥、证书，以确保重要信息不会泄露，可以使用[在线生成工具](https://www.qvdv.net/tools/qvdv-csrpfx.html, "在线生成工具")。建议证书有效期设置为10年，不设置私钥密码（不支持），证书生成完成后需要下载CRT证书、公钥和私钥并按以下名称命名：
* private.key：私钥
* public.key：公钥
* certificate.crt：证书