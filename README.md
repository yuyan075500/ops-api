# 项目介绍
该项目采用前后端分离式开发，[前端项目](https://github.com/yuyan075500/ops-web "前端项目")基于[Vue Admin Template](https://github.com/PanJiaChen/vue-admin-template "Vue Admin Template")进行二次开发，后端项目使用Golang、Gin、Gorm、CasBin进行开发。项目将主要解决项目运维过程中的**统一用户管理**和**统一系统认证**，以提高工作效能。
# 目录说明
* config：全局配置。
* controller：路由规则和业务接口的入参与响应。
* service：接口的业务处理逻辑。
* dao：数据库操作。
* model：数据库模型。
* db：数据库、缓存、对象存储客户端初始化。
* middleware：中间件层，全局逻辑处理，如跨域、JWT认证、权限校验等。
* utils：常用工具，如Token解析，文件操作等。
# 后端返回状态码说明
* 0：请求成功
* 90400：请求参数错误
* 90401：认证失败
* 90403：拒绝访问
* 90404：资源不存在
* 90500：其它错误
* 90514：Token过期或无效
# 功能概览
## 认证相关
* **SSO单点登录**：支持基于CAS 3.0、OAuth 2.0和SAML2的客户端单点登录，可以参考[单点登录配置指南](https://github.com/yuyan075500/ops-api/blob/main/deploy/sso.md "配置指南")和[已测试客户端列表](https://github.com/yuyan075500/ops-api/blob/main/deploy/sso.md#%E5%B7%B2%E6%B5%8B%E8%AF%95%E9%80%9A%E8%BF%87%E7%9A%84%E5%AE%A2%E6%88%B7%E7%AB%AF "客户端列表")。
* **用户认证**：支持~~钉钉扫码登录~~、AD用户认证和本地账号密码认证。
* **双因素**：支持Google Authenticator、阿里云APP和华为云APP。
## 其它
* 支持Swagger接口文档：访问地址：`/swagger/index.html`
* 支持用户密码自助更改：访问地址：`/reset_password`
# 项目部署
参考[Docker Compose部署](https://github.com/yuyan075500/ops-api/blob/main/deploy/deploy.md#docker-compose%E9%83%A8%E7%BD%B2 "docker-compose部署")和[Kubernetes部署](https://github.com/yuyan075500/ops-api/blob/main/deploy/deploy.md#kubernetes%E9%83%A8%E7%BD%B2 "Kubernetes部署")。
# 项目交流
如果你对此项目感兴趣，欢迎扫描下方二维码加入微信交流群  
<img src="deploy/sso_example/img/wechat.png" alt="img" width="200" height="200"/>