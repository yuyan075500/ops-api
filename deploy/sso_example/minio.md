# Minio 单点登录
支持的单点登录方式：OIDC
## 配置方法
1. **站点注册**：登录到平台，点击【资产管理】-【站点管理】-【新增】将Minio站点信息注册到平台，配置如下所示：
![img.png](img/minio-site.jpg)
配置说明：
    * 站点名称：指定一个名称，便于用户区分。
    * 登录地址：Minio控制台的登录地址。
    * SSO认证：启用。
    * 认证类型：选择`OAuth2`。
    * 站点描述：描述信息。
    * 回调地址：单点登录的回调地址，务必填写正确，默认为：`http[s]://<address>[:<port>]/oauth_callback`。
2. **Minio OIDC配置**：登录进Minio控制台，点击左侧【Identity】-【OpenID】-【Create Configuration】创建一个OIDC的提供商，如下图所示：
![img.png](img/minio-config1.jpg)
配置说明：
    * Name：指定一个OIDC提供商的名称，建议是英文，经测试中文显示会有问题。
    * Config URL：指定OIDC提供商（平台）的配置地址，默认`http[s]://<address>[:<port>]/.well-known/openid-configuration`。
    * Client ID：在平台站点详情中获取。
    * Client Secret：在平台站点详情中获取。  
    填写完成后点击【Save】并重启Minio即可生效。  

   > 说明：默认授OIDC单点登录的用户的`policy`为`readwrite`。