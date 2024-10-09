# 配置企业微信内部应用
1. **登录企业微信管理后台**：https://work.weixin.qq.com 。
2. **创建自建应用**：参考 [官方文档](https://open.work.weixin.qq.com/help2/pc/16892?person_id=1%3Freplykey%3D10aea9b3c7ab01d8948c254e43b2ww "官方文档")。
3. **配置应用登录授权**：进入【应用详情】 > 【企业微信授权登录】中配置。参考[官方文档](https://developer.work.weixin.qq.com/document/path/98151#%E5%BC%80%E5%90%AF%E7%BD%91%E9%A1%B5%E6%8E%88%E6%9D%83%E7%99%BB%E5%BD%95 "官方文档")，回调域为该平台访问域。
4. **设置应用可信域名**：进入【应用详情】 > 【网页授权及JS-SDK】中配置。参考[配置说明](https://open.work.weixin.qq.com/help2/pc/21316 "配置说明")和[配置介绍](https://developer.work.weixin.qq.com/document/path/98152#%E5%8F%82%E6%95%B0%E8%AF%B4%E6%98%8E "配置介绍")。
5. **设置应用可信IP**：进入【应用详情】 > 【企业可信IP】中配置。
# 后端应用配置
需要在配置中添加企业微信应用的相关配置。
```yaml
wechat:
  corpId: ""
  agentId: ""
  secret: ""
```
* [x] corpId：企业微信管理后台，【我的企业】中获取。
* [x] agentId：企业微信管理后台，应用详情中获取。
* [x] secret：企业微信管理后台，应用详情中获取。
# 前端应用配置
1. **配置修改**：修改配置文件`.env.production`中飞书相关的配置：
    ```js
    VUE_APP_WECHAT_APP_ID = ''
    VUE_APP_WECHAT_AGENT_ID = ''
    ```
* [x] VUE_APP_WECHAT_APP_ID：和后端应用中配置的`corpId`一致。
* [x] VUE_APP_WECHAT_AGENT_ID：和后端应用中配置的`agentId`一致。
2. **代码编译**：执行下面的命令对前端项目进行编译打包。
    ```shell
    npm install
    npm run build:prod
    ```
3. **镜像打包**：前端项目在部署时需要将编译好的前端静态文件打包到容器镜像中，推荐使用项目默认的`Dockerfile`进行打包。
# 创建本地用户
企业微信扫码登录需要事先在本地创建对应的用户，或者从LDAP、Windows AD中同步用户到本地。
# 本地用户与企业微信用户关联
在本平台的用户管理页面中，点击用户右边的【编辑】进行绑定。
![img.png](sso_example/img/ww-bind.png)