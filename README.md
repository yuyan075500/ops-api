# 项目介绍
仅需一次认证，即可访问所有授权访问的应用系统，可以为企业人员提供便捷、高效的访问体验。
## 架构设计
项目采用前后端分离架构设计，项目地址如下：
| 项目   | 项目地址 |
|:------|:-----|
| 前端   | https://github.com/yuyan075500/ops-web    |                                                                                                              |
| 后端   | https://github.com/yuyan075500/ops-api    |

如果你无法访问GitHub，可访问Gitee获取项目源代码：

| 项目   | 项目地址 |
|:------|:-----|
| 前端   | https://gitee.com/yybluestorm/ops-web    |                                                                                                              |
| 后端   | https://gitee.com/yybluestorm/ops-api    |
## 后端目录说明
* config：全局配置。
* controller：路由规则配置和接口的入参与响应。
* service：接口的处理逻辑。
* dao：数据库操作。
* model：数据库模型定义。
* db：数据库、缓存等客户端初始化。
* middleware：中间件层，作用于全局，如跨域、JWT认证、权限校验等。
* utils：工具层，如Token解析，文件操作等。
## 后端Code状态码说明
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
* **用户认证**：同时支持 [钉钉扫码登录](https://github.com/yuyan075500/ops-api/blob/main/deploy/dingtalk.md "扫码配置")、[企业微信扫码登录](https://github.com/yuyan075500/ops-api/blob/main/deploy/wechat.md "企业微信配置")、[飞书扫码登录](https://github.com/yuyan075500/ops-api/blob/main/deploy/feishu.md "飞书扫码配置")、[OpenLDAP认证、Windows AD认证](https://github.com/yuyan075500/ops-api/blob/main/deploy/deploy.md#ldap%E9%85%8D%E7%BD%AE "LDAP配置") 和本地账号认证。前端登录页面支持个性化配置，只显示某个平台，如钉钉、企业微信、飞书等，参考 [前端配置指南](https://github.com/yuyan075500/ops-web "前端配置")，按要需求修改前端项目配置文件并打包新的镜像即可。
* **双因素**：支持使用Google Authenticator、阿里云APP和华为云APP扫描获取动态验证码。

    <br>
    <img src="deploy/sso_example/img/login-1.gif" alt="img" width="350" height="200"/>
    <img src="deploy/sso_example/img/login-mfa.gif" alt="img" width="350" height="200"/>
    <br>

### 用户登录策略
✅支持，🟡敬请期待，❌不支持

| 用户来源       | 用户登录 | 账号同步 | 用户密码修改 | 用户信息修改（电话、邮箱） | 双因素认证 | 单点登录 | [NGINX鉴权](https://github.com/yuyan075500/ops-api/blob/main/deploy/sso.md#nginx%E4%BB%A3%E7%90%86%E9%89%B4%E6%9D%83 "NGINX鉴权") |
|:-----------|:-----|:-----|:-------|:--------------|:------|:-----|:------------------------------------------------------------------------------------------------------------------------------|
| 本地         | ✅    | ✅    | ✅      | ✅             | ✅     | ✅    | ✅                                                                                                                             |
| Windows AD | ✅    | ✅    | ✅      | 🟡            | ✅     | ✅    | ✅                                                                                                                             |
| OpenLDAP   | ✅    | ✅    | ✅      | 🟡            | ✅     | ✅    | ✅                                                                                                                             |
| 钉钉         | ✅    | ❌    | ❌      | ❌             | ❌     | ✅    | 🟡                                                                                                                            | 
| 企业微信       | ✅    | ❌    | ❌      | ❌             | ❌     | ✅    | 🟡                                                                                                                            | 
| 飞书         | ✅    | ❌    | ❌      | ❌             | ❌     | ✅    | 🟡                                                                                                                            | 
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
* 支持Swagger接口文档：访问地址：`/swagger/index.html`，无需要登录。
* 支持用户密码自助更改：访问地址：`/reset_password`，无需要登录。
* 支持企业网站导航：访问地址：`/sites```，无需要登录。
* 支持企业账号密码管理，登录后位于左侧【资产管理】-【账号管理】。
# 项目部署
参考 [Docker Compose部署](https://github.com/yuyan075500/ops-api/blob/main/deploy/deploy.md#docker-compose%E9%83%A8%E7%BD%B2 "docker-compose部署") 和 [Kubernetes部署](https://github.com/yuyan075500/ops-api/blob/main/deploy/deploy.md#kubernetes%E9%83%A8%E7%BD%B2 "Kubernetes部署")。
# 开发环境搭建
参考 [开发环境搭建](https://github.com/yuyan075500/ops-api/blob/main/deploy/dev.md "开发环境搭建")。
# 项目交流
如果你对此项目感兴趣，可添加作者联系方式
WeChat：270142877。  
Email：270142877@qq.com。  
<br>
联系时请注名来意。