# Jira 单点登录
Jira 支持的单点登录方式：OAuth2。
## 配置方法
1. **站点注册**：登录到 IDSphere 统一认证平台，点击【资产管理】-【站点管理】-【新增】将 Jira 站点信息注册到 IDSphere 统一认证平台，配置如下所示：<br><br>
![img.png](img/jira-site.jpg)<br><br>
   * 站点名称：指定一个名称，便于用户区分。
   * 登录地址：Jira 的登录地址。
   * SSO认证：启用。
   * 认证类型：选择 `OAuth2`。
   * 站点描述：描述信息。
   * 回调地址：Jira 的回调地址，默认为：`http[s]://<address>[:<port>]/plugins/servlet/oauth/callback`。<br><br>
2. **应用安装**：登录到 Jira 并点击右上角齿轮进入【管理应用】，在应用商店搜索 `OIDC SSO`，找到如下图所示的插件并安装。<br><br>
![img.png](img/jira-marketplace.jpg)<br><br>
应用安装完成后根据提示申请1个30天免费适用的 License，并激活。Licence 过期后，无法使用单点登录，可以注册一个 Atlassian 账号，每月都可以申请一个30天的免费 License。<br><br>
3. **应用配置**：登录到 Jira 并点击右上角齿轮进入【管理应用】，点击左侧的【miniOrange OAuth/OIDC SSO】进入 OIDC 配置，如下图所示：<br><br>
![img.png](img/jira-app-config.jpg)<br><br>
接下来点击【Add New App】按钮创建一个身份提供商，如下图所示：<br><br>
![img.png](img/jira-app-config1.jpg)<br><br>
选择【Custom OAuth】创建一个自定义 OAuth 应用，如下图所示：<br><br>
![img.png](img/jira-app-config2.jpg)<br><br>
将 OAuth 的配置信息填到下面的表单中，如下图所示：<br><br>
![img.png](img/jira-app-config3.jpg)<br><br>
   * Custom App Name：可以自定义，用于显示在 Jira 的登录页。
   * Client Id：在 IDSphere 统一认证平台的站点详情中获取。
   * Client Secret：在 IDSphere 统一认证平台的站点详情中获取。
   * Scope：`openid`。
   * Authorization Endpoint：`<externalUrl>/login`。
   * Access Token Endpoint：`<externalUrl>/api/v1/sso/oauth/token`。
   * User Info Endpoint：`<externalUrl>/api/v1/sso/oauth/userinfo`。<br><br>
   配置完成后点击【Save】保存。 <br><br>
   **注意**：在此页面显示了 Jira 正确的 `Callback URL`，可以将此地址复制到 IDSphere 统一认证平台的 Jira 站点配置中作为回调地址。<br><br>
4. **用户配置**：点击左侧的【User Profile】，需要进行用户属性配置，如下图所示：<br><br>
![img.png](img/jira-app-config5.jpg)<br><br>
   * Username Attribute：`username`。
   * Email Attribute：`email`。
   * Full Name Attribute：`name`。<br><br>
   **注意**：建议将 `User Profile Mapping` 选项开启，以便自动更新用户信息到已存在的用户。<br><br>
6. **高级设置**：接上步，点击左侧的【Advanced Settings】找到 `Send Parameters in Token Endpoint`，将值更改为 `Http Body`，如下图所示：<br><br>
![img.png](img/jira-app-config4.jpg)