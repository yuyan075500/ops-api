# Jenkins 单点登录
支持的单点登录方式：OIDC
## 配置方法
1. **站点注册**：登录到平台，点击【资产管理】-【站点管理】-【新增】将Confluence站点信息注册到平台，配置如下所示：
   ![img.png](img/jira-site.jpg)
    * 站点名称：指定一个名称，便于用户区分。
    * 登录地址：Confluence的登录地址。
    * SSO认证：启用。
    * 认证类型：选择`OAuth2`。
    * 站点描述：描述信息。
    * 回调地址：单点登录的回调地址，务必填写正确，默认为：`<protocol>://<address>[:<port>]/plugins/servlet/oauth/callback`。
2. **应用安装**：登录到Jira并点击右上角齿轮进入【管理应用】，在应用商店搜索`OIDC SSO`，找到如下图所示的插件并安装。
   ![img.png](img/confluence-marketplace.jpg)
   应用安装完成后按提示申请1个30天免费适用License，并激活。
   > **提示**：该应用过期后，无法使用单点登录，可以注册一个Atlassian账号，每月申请一个30天的免费License使用。
3. **应用配置**：登录到Confluence并点击右上角齿轮进入【管理应用】，点击左侧的【miniOrange OAuth/OIDC SSO】进入OIDC配置。
   点击【Add New App】按钮创建一个身份提供商，如下图所示：
   ![img.png](img/confluence-config1.jpg)
   选择【Custom OIDC】创建一个自定义OIDC应用，如下图所示：
   ![img.png](img/confluence-config2.jpg)
   点击【Import Details】，将OIDC的配置信息导入到下面的表单中，如下图所示：
   ![img.png](img/confluence-config3.jpg)
   在打开的表单中，Well-Known Endpoint地址为平台的：`<protocol>://<address>[:<port>]/.well-known/openid-configuration`。导入成功后还需要填写应用的`Client Id`和`Client Secret`，这两项配置从平台的站点详情中获取。
   配置完成后点击【Save】保存即可。
   > **提示**：
   > * 在此页面显示了Jira正确的`Callback URL`，可以将此地址复制到平台的站点配置中覆盖之前的回调地址，以确保平台配置正确。
   > * 表单中`Custom App Name`可以自定义，用于显示在在Jira的登录页。
4. **用户配置**：接上步，点击左侧的【User Profile】，需要进行用户属于配置，按如下图所示：
   ![img.png](img/confluence-config5.jpg)
    * Username Attribute：`username`
    * Email Attribute：`email`
    * Full Name Attribute：`name`
5. **高级设置**：接上步，点击左侧的【Advanced Settings】找到`Send Parameters in Token Endpoint`，将值更改为`Http Body`，如下图所示：
   ![img.png](img/confluence-config4.jpg)