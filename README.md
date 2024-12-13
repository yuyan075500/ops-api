# IDSphere 统一认证平台项目介绍
仅需一次认证，即可访问所有授权访问的应用系统，为企业办公人员提供高效、便捷的访问体验。
## 架构设计
项目采用前后端分离架构设计，项目地址如下：

| 项目  | 项目地址                                   |
|:----|:---------------------------------------|
| 前端  | https://github.com/yuyan075500/ops-web |                                                                                                              |
| 后端  | https://github.com/yuyan075500/ops-api |

如果你无法访问`GitHub`，可访问`Gitee`获取项目源代码，地址如下：

| 项目  | 项目地址                                  |
|:----|:--------------------------------------|
| 前端  | https://gitee.com/yybluestorm/ops-web |                                                                                                              |
| 后端  | https://gitee.com/yybluestorm/ops-api |
## 目录说明
* config：全局配置。
* controller：路由规则配置和接口的入参与响应。
* service：接口处理逻辑。
* dao：数据库操作。
* model：数据库模型定义。
* db：数据库、缓存以及文件存储客户端初始化。
* middleware：全局中间件层，如跨域、JWT认证、权限校验等。
* utils：全局工具层，如Token解析、文件操作、字符串操作以及加解密等。
## Code状态码说明
* 0：请求成功。
* 90400：请求参数错误。
* 90401：认证失败。
* 90403：拒绝访问。
* 90404：访问的对象或资源不存在。
* 90514：Token过期或无效。
* 90500：其它服务器错误。
# 项目功能介绍
## 认证相关
* **SSO单点登录**：支持与使用 `CAS 3.0`、`OAuth 2.0`、`OIDC`和`SAML2` 协议的客户端进行对接，对接方法可以参考 [客户端配置指南](https://github.com/yuyan075500/ops-api/blob/main/deploy/sso.md "配置指南") 和 [已测试客户端列表](https://github.com/yuyan075500/ops-api/blob/main/deploy/sso.md#%E5%B7%B2%E6%B5%8B%E8%AF%95%E9%80%9A%E8%BF%87%E7%9A%84%E5%AE%A2%E6%88%B7%E7%AB%AF "客户端列表")。
* **用户认证**：支持 [钉钉扫码登录](https://github.com/yuyan075500/ops-api/blob/main/deploy/dingtalk.md "扫码配置")、[企业微信扫码登录](https://github.com/yuyan075500/ops-api/blob/main/deploy/wechat.md "企业微信配置")、[飞书扫码登录](https://github.com/yuyan075500/ops-api/blob/main/deploy/feishu.md "飞书扫码配置")、[OpenLDAP 账号密码认证、Windows AD 账号密码认证](https://github.com/yuyan075500/ops-api/blob/main/deploy/deploy.md#ldap%E9%85%8D%E7%BD%AE "LDAP配置") 和本地用户密码认证方式登录。另外前端登录页面支持个性化配置，隐藏或显示必要的登录选项，可以参考 [前端配置指南](https://github.com/yuyan075500/ops-web "前端配置")。
* **双因素认证**：支持使用 Google Authenticator、阿里云和华为云手机 APP 进行双因素认证，双因素认证仅在使用账号密码认证时生效。

    <br>
    <img src="deploy/sso_example/img/login-1.gif" alt="img" width="350" height="200"/>
    <img src="deploy/sso_example/img/login-mfa.gif" alt="img" width="350" height="200"/>
    <br>

### 用户登录策略
✅支持，🟡敬请期待，❌不支持

| 登录方法          | 本地登录 | 双因素认证 | SSO 登录 | [NGINX鉴权](https://github.com/yuyan075500/ops-api/blob/main/deploy/sso.md#nginx%E4%BB%A3%E7%90%86%E9%89%B4%E6%9D%83 "NGINX鉴权") |
|:--------------|:-----|:------|:-------|:------------------------------------------------------------------------------------------------------------------------------|
| 本地账号          | ✅    | ✅     | ✅      | ✅                                                                                                                             |
| Windows AD 账号 | ✅    | ✅     | ✅      | ✅                                                                                                                             |
| OpenLDAP 账号   | ✅    | ✅     | ✅      | ✅                                                                                                                             |
| 钉钉扫码          | ✅    | ❌     | ✅      | 🟡                                                                                                                            | 
| 企业微信扫码        | ✅    | ❌     | ✅      | 🟡                                                                                                                            | 
| 飞书扫码          | ✅    | ❌     | ✅      | 🟡                                                                                                                            | 

关于 Windows AD 或 OpenLDAP 配置可以参考 [配置指南](https://github.com/yuyan075500/ops-api/blob/main/deploy/deploy.md#ldap%E9%85%8D%E7%BD%AE "LDAP配置")，关于同步和登录策略可以参考 [注意事项](https://github.com/yuyan075500/ops-api/blob/main/deploy/ldap.md "注意事项")。

**注意：使用 Windows AD 或 OpenLDAP 登录需要确保在 IDSphere 统一认证平台中存在对应的用户，否则无法登录。**
## 企业级账号管理
## 域名及证书管理
## 其它
* 支持`Swagger`接口文档：部署成功后访问地址为：`/swagger/index.html`，无需要登录。
* 支持用户密码自助更改：部署成功后访问地址：`/reset_password`，无需要登录。
* 支持企业网站导航：部署成功后访问地址：`/sites`，无需要登录。
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