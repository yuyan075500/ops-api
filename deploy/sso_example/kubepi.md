# Kubepi 单点登录
支持的单点登录方式：OIDC
## 配置方法
1. **站点注册**：登录到平台，点击【资产管理】-【站点管理】-【新增】将Kubepi站点信息注册到平台，配置如下所示：
![img.png](img/kubepi-site.jpg)
    * 站点名称：指定一个名称，便于用户区分。
    * 登录地址：Kubepi控制台的登录地址。
    * SSO认证：启用。
    * 认证类型：选择`OAuth2`。
    * 站点描述：描述信息。
    * 回调地址：单点登录的回调地址，务必填写正确，默认为：`<protocol>://<address>[:<port>]/kubepi/api/v1/sso/callback`。
2. **Kubepi OIDC配置**：登录进Kubepi控制台，点击左侧【用户管理】-【SSO】如下图所示：
![img.png](img/kubepi-config.jpg)
    * 协议：选择`OpenID Connect`。
    * 接口地址：指定OIDC提供商（平台）的配置地址，默认`<protocol>://<address>[:<port>]`，Kubepi会默认在此路径后面加上`/.well-known/openid-configuration`
    * 客户端ID：在平台站点详情中获取。
    * 客户端密钥：在平台站点详情中获取。  
    填写完成后点击【Save】并重启Minio即可生效。
