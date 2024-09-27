# 配置企业微信内部应用
1. **登录企业微信管理后台**：https://work.weixin.qq.com。
2. **创建自建应用**：参考 [官方文档](https://open.work.weixin.qq.com/help2/pc/16892?person_id=1%3Freplykey%3D10aea9b3c7ab01d8948c254e43b2ww "官方文档")。
3. **开启自建应用网页授权登录**：参考[官方文档](https://developer.work.weixin.qq.com/document/path/98151#%E5%BC%80%E5%90%AF%E7%BD%91%E9%A1%B5%E6%8E%88%E6%9D%83%E7%99%BB%E5%BD%95 "官方文档")，回调域为该平台访问域。
4. **配置应用可信域名**：
5. **配置应用可信IP**：
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
# 创建本地用户
企业微信扫码登录需要事先在本地创建对应的用户，或者从LDAP、Windows AD中同步用户到本地。默认只要用户姓名和电话与钉钉用户一致处且本地用户状态正常，则可以登录成功。
# 本地用户与企业微信用户关联