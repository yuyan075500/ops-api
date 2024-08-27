# SSO介绍
该平台提供了SSO单点登录功能，支持CAS3.0、OAuth2.0、SAML2和OIDC协议对接，支持多种客户端。本文档中提供了[已通过测试的客户端列表](#已测试通过的客户端)，以及各个客户端的配置说明。
# CAS3.0客户端配置
# OAuth2.0客户端配置
# SAML2客户端配置
# OIDC客户端配置
# Nginx代理鉴权
对于一些客户端可以在没有账号密码的情况下进行访问，如：Kibana、Consul Server UI等，为了实现这类客户端的认证，可以使用Nginx对这类客户端进行代理，跳转至本平台进行认证。
<br>
因为是基于Cookie实现，对于这类客户端也有一定使用，要求如下：
1. 客户端的所在域必须与SSO平台所在的二级域一致，假如平台的访问域名为：`test.ops.cn`，则客户端的所在域必须为`xxx.ops.cn`。
2. 不支持非域名访问的客户端。如：`127.0.0.1`、`localhost`、`192.168.1.10`等。 
3. 不支持非HTTPS应用。
## 认证规则
假设有A、B、C三个客户端应用都使用Nginx进行代理鉴权，鉴权规则如下：
* 如果用户未登录平台，当访问A、B、C三个客户端应用中的任何一个都将跳转至登录界面。
* 如果用户已经登录平台，当访问A、B、C三个客户端应用中的任何一个都能直接访问应用。
## Nginx配置
在开始配置前请确保Nginx支持auth_request模块，可以使用命令`nginx -V`查看，具体配置如下：
```nginx
server {
	listen 80;

	location / {
		auth_request /auth;
		error_page 401 500 = @error401;
		# 下面是被代理应用的相关配置
	}
	
	location = /auth {
		internal;
		proxy_pass_request_body off;
		proxy_set_header Content-Length "";
		proxy_set_header X-Original-URI $request_uri;
		proxy_set_header Cookie $http_cookie;
		proxy_pass https://<平台域名>/api/v1/sso/cookie/auth;
	}

	# 认证失败后的处理
	location @error401 {
		# 跳转至登录页
		return 302 https://<平台域名>/login?nginx_redirect_uri=$scheme://$host$request_uri;
	}
}
```
在上面的示例中对`/`根路径进么了鉴权，如果需要对其它路径进行鉴权可以添加其它location，并配置`auth_request`和`error_page`。
## Kubernetes Ingress配置
```yaml
kind: Ingress
apiVersion: networking.k8s.io/v1
metadata:
  name: my-app
  annotations:
    nginx.ingress.kubernetes.io/auth-url: https://<平台域名>/api/v1/sso/cookie/auth
    nginx.ingress.kubernetes.io/server-snippet: |
      error_page 401 500 = @login;
      proxy_set_header Cookie $http_cookie;
      location @login {
        return 302 https://<平台域名>/login?nginx_redirect_uri=$scheme://$host$request_uri;
      }
spec:
  ingressClassName: nginx
  tls:
    - hosts:
        - my-app.test.cn
      secretName: test.cn
  rules:
    - host: my-app.test.cn
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: my-app
                port:
                  number: 80
```
# 已测试通过的客户端
针对CAS3.0、OAuth2.0、SAML2、OIDC协议，目前已测试的SSO客户端如下：
| 客户端名称    | 协议名称     | 参考文档                                                                                                       |
| :---        |    :----    |          ---                                                                                                 |
| Grafana     | OAuth2.0    | [参考文档](https://github.com/yuyan075500/ops-api/blob/main/deploy/sso_example/grafana.md "参考文档")           |
| Jenkins     | CAS3.0      | [参考文档](https://github.com/yuyan075500/ops-api/blob/main/deploy/sso_example/jenkins.md "参考文档")           |
| Zabbix      | SAML2       | [参考文档](https://github.com/yuyan075500/ops-api/blob/main/deploy/sso_example/zabbix.md "参考文档")            |
| 华为云       | SAML2       | [参考文档](https://github.com/yuyan075500/ops-api/blob/main/deploy/sso_example/huawei_cloud.md "参考文档")      |
| JumpServer  | OAuth2.0    | [参考文档](https://github.com/yuyan075500/ops-api/blob/main/deploy/sso_example/jumpserver.md "参考文档")        |
| Jira        | OAuth2.0    | [参考文档](https://github.com/yuyan075500/ops-api/blob/main/deploy/sso_example/jira.md "参考文档")              |
| Confluence  | OAuth2.0    | [参考文档](https://github.com/yuyan075500/ops-api/blob/main/deploy/sso_example/confluence.md "参考文档")        |
| KubePi      | SAML2       | 未完成      |
| 阿里云       | SAML2       | [参考文档](https://github.com/yuyan075500/ops-api/blob/main/deploy/sso_example/aliyun.md "参考文档")            |
| 腾讯云       | SAML2       | [参考文档](https://github.com/yuyan075500/ops-api/blob/main/deploy/sso_example/tencent.md "参考文档")           |
| GitLab      | CAS3.0      | 未完成      |