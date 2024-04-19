# 目录说明
* config：全局配置，如监听地址、数据库连接配置等。
* controller：定义路由规则以及接口同业入参和响应。
* service：处理接口的业务逻辑。
* dao：数据库操作。
* model：定义数据库表信息。
* db：数据库、缓存、对象存储客户端初始化。
* middleware：中间件层，添加全局逻辑处理，如跨域、JWT认证等。
* utils：工具目录，定义常用工具，如Token解析，文件操作等。
# 项目依赖
* MySQL
* Redis
* Minio
# 功能概览
* 用户管理
* RBAC权限管理
* 双因素认证（支持Google Authenticator）
* 单点登录（CAS 3.0、OAuth 2.0、SAML 2）
* 钉钉扫码登录（需要配置钉钉应用）