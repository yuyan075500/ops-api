# 创建Grafana站点
Grafana支持使用OAuth2进行单点登录，站点配置如下：
![img.png](img/grafana.png)
注：Grafana的回调地址为：`<protocol>://<address>[:<port>]/login/generic_oauth`
# 修改Grafana配置
编辑`grafana.ini`配置文件，修改的配置如下：
```shell
[server]
# Grafana访问地址相关配置，一定要配置正确
protocol = http
http_port = 3000
domain = 192.168.200.21

[auth]
# 禁用登录页面，可以自动跳转至统一认证中心
disable_login_form = true
# 允许通过电子邮件查找用户（唯一标识）
oauth_allow_insecure_email_lookup = true

[auth.generic_oauth]
# 开启OAuth认证
enabled = true
# OAuth认证的名称
name = 信息化统一认证中心
# 允许用户注册
allow_sign_up = true
# 自动登录
auto_login = true
# client_id，从统一认证中心获取
client_id = tYOvydGyamAQUTcZ
# client_secret，从统一认证中心获取
client_secret = GXQzHtPSDAHIKRHuFHpSHarKQjDIIXmG
scopes = openid
auth_url = https://ops-test.50yc.cn/login
token_url = https://ops-test.50yc.cn/api/v1/oauth/token
api_url = https://ops-test.50yc.cn/api/v1/oauth/userinfo
```
* auth_url：前端登录地址
* token_url：获取Token接口
* api_url：获取用户信息接口