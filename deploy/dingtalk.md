# 配置钉钉应用
1. **登录钉钉开放平台**：https://open.dingtalk.com ，并进入开发者后台。
2. **创建应用**：参考 [官方文档](https://open.dingtalk.com/document/orgapp/create-an-application "官方文档")。
3. **应用配置**：进入应用详情页，单击这【开发配置】 > 【安全设置】，填写重定向 URL（回调域名）。回调域名为该平台的登录地址，为`http[s]://<address>[:<port>]/login`，当用户扫码后浏览器默认会跳转至该地址。
4. **应用授权**：进入应用详情页，单击【开发配置】 > 【权限管理】，在权限搜索框中输入权限名称并申请权限。需要授与该应用`Contact.User.mobile`和`Contact.User.Read`权限，参考 [官方文档](https://open.dingtalk.com/document/orgapp/tutorial-obtaining-user-personal-information#c4647d84328mg "官方文档")
5. **发布应用**：参考 [官方文档](https://open.dingtalk.com/document/orgapp/publish-dingtalk-application "官方文档")。
# 后端应用配置
# 前端应用配置