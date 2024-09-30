# 配置钉钉应用
1. **登录飞机开放平台**：https://open.feishu.cn ，并进入开发者后台。
2. **创建自建应用**：参考 [官方文档](https://open.feishu.cn/document/home/introduction-to-custom-app-development/self-built-application-development-process#a0a7f6b0 "官方文档")。
3. **应用配置**：进入应用详情页，单击【开发配置】 > 【安全设置】，填写重定向 URL（回调域名）。回调域名为该平台的登录地址，为`http[s]://<address>[:<port>]/login`，当用户扫码后浏览器默认会跳转至该地址。
4. **应用授权**：进入应用详情页，单击【开发配置】 > 【权限管理】，在权限搜索框中输入权限名称并申请权限。需要授与该应用`contact:user.email:readonly`和`contact:user.phone:readonly`权限。)
5. **发布应用**：参考 [官方文档](https://open.feishu.cn/document/home/introduction-to-custom-app-development/self-built-application-development-process#baf09c7d "官方文档")。
# 后端应用配置
需要在配置中添加钉钉应用的相关配置。
```yaml
feishu:
  appId: ""
  appSecret: ""
```
* [x] appId：在飞书开放平台，应用详情页，左侧的【凭证与基础信息】中获取。
* [x] appSecret：在飞书开放平台，应用详情页，左侧的【凭证与基础信息】中获取。
# 前端应用配置
1. **配置修改**：修改配置文件`.env.production`中飞书相关的配置：
    ```js
    VUE_APP_FEISHU_CLIENT_ID = ''
    ```
* [x] VUE_APP_FEISHU_CLIENT_ID：和后端应用中配置的`appId`一致。
2. **代码编译**：执行下面的命令对前端项目进行编译打包。
    ```shell
    npm install
    npm run build:prod
    ```
3. **镜像打包**：前端项目在部署时需要将编译好的前端静态文件打包到容器镜像中，推荐使用项目默认的`Dockerfile`进行打包。
# 创建本地用户
飞书扫码登录需要事先在本地创建对应的用户，或者从LDAP、Windows AD中同步用户到本地。确保用户姓名、电话和邮箱与飞书用户一致且本地用户状态正常，则可以登录成功。