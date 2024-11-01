# 阿里云单点登录
支持的单点登录方式：SAML2
## 配置方法
1. **获取阿里云元数据**：登录阿里云，进入【RAM访问控制】，按下图所示，依次点击进入【用户SSO】管理。
![img.png](img/aliyun-metadata.jpg)
请保存SAML服务提供商元数据URL，后续在平台注册站点时使用。
2. **创建身份提供商**：接上步，点击SSO登录设置右边的【编辑】按钮，打开SSO登录设置，如下图所示：
![img.png](img/aliyun-sso-config.jpg)
配置说明：
   * SSO功能状态：开启。
   * 元数据文件：这里需要上传IDP元数据文件，IDP的元数据文件可以访问平台`http[s]://<address>[:<port>]/api/v1/sso/saml/metadata`获取。
   * 辅助域名：关闭。
3. **获取登录地址**：点击【概览】即可获取，如下图所示：
![img.png](img/aliyun-login-url.jpg)
请保存登录地址，后续在平台注册站点时使用。
4. **获取账号域名**：点击【设置】-【账号域名】获取，如下所示：
![img.png](img/aliyun-domain.jpg)
请保存域名，后续在平台注册站点时使用。
5. **创建RAM用户**：需要确保用户名和平台中用户的`username`保持一致，所有需要登录的账号都需要创建一个与之对应的RAM实体用户。
6. **站点注册**：登录到平台，点击【资产管理】-【站点管理】-【新增】将华为云站点信息注册到平台，配置如下所示：
![img.png](img/aliyun-site.jpg)
配置说明：
   * 站点名称：指定一个名称，便于用户区分。
   * 登录地址：填写从第3步中获取的登录地址。
   * SSO认证：启用。
   * 认证类型：选择`SAML2`。
   * 站点描述：描述信息。
   * SP Metadata URL：填写从第1步中获取的地址，点击【获取】可以自动从阿里云元数据中加载`SP EntityID`和`SP 证书`相关信息。
7. **站点修改**：登录到平台`MySQL`数据库，在`site`表中找到刚注册的站点信息，将字段`domain_id`的值修改为从第4步中获取的域名。
8. **登录测试**：在浏览器打开从第3步中获取的登录地址，然后点击【使用企业账号登录】即可，如下图所示：
![img.png](img/aliyun-login.jpg)