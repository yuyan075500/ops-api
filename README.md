# 项目介绍
该项目主要提供**统一用户管理**和**统一系统认证**服务，采用前后端分离的架构模式。后端项目基于Gin + Gorm + Casbin实现，[前端项目](https://github.com/yuyan075500/ops-web "前端项目") 基于 [Vue Admin Template](https://github.com/PanJiaChen/vue-admin-template "Vue Admin Template") 进行二次开发。
# 目录说明
* config：全局配置。
* controller：路由规则配置和接口的入参与响应。
* service：接口的处理逻辑。
* dao：数据库操作。
* model：数据库模型定义。
* db：数据库、缓存等客户端初始化。
* middleware：中间件层，作用于全局，如跨域、JWT认证、权限校验等。
* utils：工具层，如Token解析，文件操作等。
# 后端Code状态码说明
* 0：请求成功。
* 90400：请求参数错误。
* 90401：认证失败。
* 90403：拒绝访问。
* 90404：访问的对象或资源不存在。
* 90500：其它错误。
* 90514：Token过期或无效。
# 功能概览
## 认证相关
* **SSO单点登录**：支持`CAS 3.0`、`OAuth 2.0`和`SAML2`协议，可以参考 [客户端配置指南](https://github.com/yuyan075500/ops-api/blob/main/deploy/sso.md "配置指南") 和 [已测试客户端列表](https://github.com/yuyan075500/ops-api/blob/main/deploy/sso.md#%E5%B7%B2%E6%B5%8B%E8%AF%95%E9%80%9A%E8%BF%87%E7%9A%84%E5%AE%A2%E6%88%B7%E7%AB%AF "客户端列表")。
* **用户认证**：同时支持钉钉扫码登录、~~企业微信扫码登录~~、~~飞书扫码登录~~、OpenLDAP认证、Windows AD认证和本地账号认证。
* **双因素**：支持使用Google Authenticator、阿里云APP和华为云APP扫描获取动态验证码。
### 用户登录策略
| 用户来源       | 用户登录   | 账号同步  | 用户密码修改 | 用户信息修改（电话、邮箱） | 双因素认证 | 单点登录 | [NGINX鉴权](https://github.com/yuyan075500/ops-api/blob/main/deploy/sso.md#nginx%E4%BB%A3%E7%90%86%E9%89%B4%E6%9D%83 "NGINX鉴权") |
|:-----------|:-------|:------|:-------|:--------------|:------|:-----|:------------------------------------------------------------------------------------------------------------------------------|
| 本地         | ✅ 支持   | ✅ 支持  | ✅ 支持   | ✅ 支持          | ✅ 支持  | ✅ 支持 | ✅ 支持                                                                                                                          |
| Windows AD | ✅ 支持   | ✅ 支持  | ✅ 支持   | 🟡 待支持        | ✅ 支持  | ✅ 支持 | ✅ 支持                                                                                                                          |
| OpenLDAP   | ✅ 支持   | ✅ 支持  | ✅ 支持   | 🟡 待支持        | ✅ 支持  | ✅ 支持 | ✅ 支持                                                                                                                          |
| 钉钉         | 🟡 待支持 | ❌ 不支持 | ❌ 不支持  | ❌ 不支持         | ❌ 不支持 | ✅ 支持 | 🟡 待支持                                                                                                                        | 
| 企业微信       | 🟡 待支持 | ❌ 不支持 | ❌ 不支持  | ❌ 不支持         | ❌ 不支持 | ✅ 支持 | 🟡 待支持                                                                                                                        | 
| 飞书         | 🟡 待支持 | ❌ 不支持 | ❌ 不支持  | ❌ 不支持         | ❌ 不支持 | ✅ 支持 | 🟡 待支持                                                                                                                        | 
### 账号同步规则
无论使用哪一种用户认证方式，都需要确保本地系统中用户存在，所以当配置好Windows AD或OpenLDAP后，需要登录平台点击【用户管理】-【分组管理】-【LDAP账号同步】执行一次用户同步，用户的同步规则如下：
1. 如果本地系统中没有，LDAP中有，则创建。
2. 如果本地系统有，LDAP中有，则更新（仅更新用户来源为LDAP且`username`相同，更新手机号和邮箱字段）。

同步用户时，用户字段映射规则如下：
```shell
# OpenLDAP
{
	"name": "uid",
	"username": "cn",
	"email": "mail",
	"phone_number": "telephoneNumber"
}

# Windows AD
{
	"name": "sAMAccountName",
	"username": "cn",
	"email": "mail",
	"phone_number": "telephoneNumber",
	"is_active": "userAccountControl"
}
```
## 其它
* 支持Swagger接口文档：访问地址：`/swagger/index.html`。
* 支持用户密码自助更改：访问地址：`/reset_password`。
# 项目部署
参考 [Docker Compose部署](https://github.com/yuyan075500/ops-api/blob/main/deploy/deploy.md#docker-compose%E9%83%A8%E7%BD%B2 "docker-compose部署") 和 [Kubernetes部署](https://github.com/yuyan075500/ops-api/blob/main/deploy/deploy.md#kubernetes%E9%83%A8%E7%BD%B2 "Kubernetes部署")。
# 项目交流
如果你对此项目感兴趣，欢迎扫描下方二维码加入微信交流群。  
<br>
<img src="deploy/sso_example/img/wechat.png" alt="img" width="200" height="200"/>
## 作者联系方式
WeChat：270142877。  
Email：270142877@qq.com。  
<br>
联系时请注名来意。