# 项目部署
支持使用 [Docker Compose一键部署](#docker-compose部署) 和 [Kubernetes部署](#Kubernetes部署)，生产环境推荐使用 Kubernetes 部署。
## Docker Compose部署
如果你想快速拥有一个简易的环境用于测试、演示，且对性能、稳定性和安全性没有任何求的，那么推荐使用该部署方式。  
1. **部署环境准备**：你需要准备一台 Linux 服务器，并安装以下组件。
   * [x] Docker。
   * [x] Docker Compose。

    `Docker`和`Docker Compose`是部署毅必须准备的，其它组件在 `docker-compose.yaml` 配置清单中已指定。
2. **克隆项目**：
    ```shell
    git clone https://github.com/yuyan075500/ops-api.git
    或
    git clone https://gitee.com/yybluestorm/ops-api
    ```
3. **切换工作目录**：
    ```shell
    cd ops-api/deploy/docker-compose
    ```
4. **配置环境变量**：修改 `.env` 文件中环境变量。<br>
   此配置文件中主要指定了数据库、缓存、MinIO 等组件的初始化以及项目启动的系统版本，该步骤可以跳过，也可以按需要修改。
5. **修改项目配置**：修改`conf/config.yaml`文件中相关配置。<br>
   参考 [配置文件说明](#配置文件说明)，以下配置必须修改：
   * `externalUrl` 需要更改为 IDSphere 统一认证平台在浏览器实际的访问地址，否则导致单点功能等相关功能无法正常使用。
   * `oss.accessKey` 和 `oss.secretKey` 中指定的 `AK` 和 `SK` 需要在 Minio 启动完成后登录到后台手动创建。
   * `oss.endpoint` 配置的地址必须确保使用 IDSphere 统一认证平台的客户端电脑可以访问，如果实际的地址协议为 `HTTPS` 则需要将 `oss.ssl` 更改为 `true`。
6. **创建证书**。<br>
   参考 [创建项目证书](#项目证书)，将生成的新证书保存至`certs`目录中并覆盖目标文件，测试环境可以跳过此步骤。
7. **创建Minio数据目录**：需要手动创建Minio数据目录，并更改权限为`1001:1001`。
    ```shell
    mkdir -p data/minio
    chown -R 1001:1001 data/minio
    ```
8. **执行部署**：
    ```shell
    docker-compose up -d
    ```
9. **数据初始化**：将`deploy/data.sql`SQL中的数据导入到数据库中。

   > **注意**：如果使用的外部数据库，请确保数据库使用的字符集为`utf8mb4`，排序规则为`utf8mb4_general_ci`。如果你使用的默认数据库，`data.sql`文件已经打包进容器`ops-mysql`的`/root/data.sql`路径，可以直接导入。

   1. **系统登录**：部署完成后，系统会自动创建一个超级用户，此用户不受Casbin权限控制。用户名为：`admin`，密码为：`admin@123...`。
## Kubernetes部署
你需要自行准备以下相关资源：
* [x] [Kubernetes](https://kubernetes.io "Kubernetes") 运行环境。
* [x] [Helm](https://helm.sh "Helm") 客户端，确保能访问Kubernetes集群。
* [x] MySQL 8.0。
* [x] Redis 5.x。
* [x] MinIO或华为云OBS。
### 部署
1. **克隆项目**：将项目克隆到本地Helm客户端所在服务器。
    ```shell
    git clone https://github.com/yuyan075500/ops-api.git
    ```
2. **切换工作目录**：
    ```shell
    cd ops-api/deploy/kubernetes
    ```
3. **证书**：创建 [项目证书](#项目证书)，证书创建完成后使用新的证书替换`templates/configmap.yaml`文件中对应的内容。
4. **项目配置**：修改`templates/configmap.yaml`文件中`config.yaml`项的相关配置，请参考 [配置文件说明](#配置文件说明)。

   > **注意**：
   > * 必须修改`externalUrl`的值为实际的访问地址，否则导致单点功能登录无法使用。
   > * 另外MinIO的`accessKey`和`secretKey`需要在部署成功后登录进MinIO控制台手动创建，确保与`conf/config.yaml`配置文件中指定的值相同即可，默认值可以自行修改。
   > * Minio的`endpoint`项配置的地址必须确保使用该平台的客户端电脑可以访问，否则图片上传成功后将无法访问。

5. **部署**：
   ```shell
   helm install <自定义应用名> --namespace <名称空间> .
   ```

   > 说明：如果你使用Kubernetes之外的代理程序，那么你需要将`Service`类型修改为`NodePort`，并参考`templates/ingress.yaml`模板文件中的转发规则进行相关配置。如果有使用[钉钉](https://github.com/yuyan075500/ops-api/blob/main/deploy/dingtalk.md "钉钉配置")、[企业微信](https://github.com/yuyan075500/ops-api/blob/main/deploy/wechat.md "企业微信配置")或[飞书](https://github.com/yuyan075500/ops-api/blob/main/deploy/feishu.md "飞书配置")扫码认证，请按要求对前端项目进行单独构建打包，并修改`value.yaml`文件对中应前端的镜像配置。

6. **数据初始化**：将`deploy/data.sql`SQL中的数据导入到数据库中。

   > **注意**：请确保数据库使用的字符集为`utf8mb4`，排序规则为`utf8mb4_general_ci`。

7. **系统登录**：部署完成后，系统会自动创建一个超级用户，此用户不受Casbin权限控制。用户名为：`admin`，密码为：`admin@123...`。

   > 说明：如果需要高可用只需调整应用的副本数即可，数据库和中间件的高可用需要自行完成。

# 配置文件说明
```yaml
server: "0.0.0.0:8000"
externalUrl: ""
secret: "swfqezjzoqssvjck"
mysql:
  host: "127.0.0.1"
  port: 3306
  db: "ops"
  user: "root"
  password: ""
  maxIdleConns: 10
  maxOpenConns: 100
  maxLifeTime: 30
redis:
  host: "127.0.0.1:6379"
  password: ""
  db: 0
jwt:
  secret: "swfqezjzoqssvjck"
  expires: 6
mfa:
  enable: false
  issuer: "统一认证平台"
oss:
  endpoint: ""
  accessKey: ""
  secretKey: ""
  bucketName: ""
  ssl: true
ldap:
  host: ""
  bindUserDN: ""
  bindUserPassword: ""
  searchDN: ""
  userAttribute: ""
  maxPasswordAge: 90
sms:
  provider: ""
  url: ""
  appKey: ""
  appSecret: ""
  callbackUrl: ""
  resetPassword:
    sender: ""
    templateId: ""
    signature: ""
mail:
  smtpHost: ""
  smtpPort: 587
  from: ""
  password: ""
dingTalk:
   appKey: ""
   appSecret: ""
wechat:
   corpId: ""
   agentId: ""
   secret: ""
feishu:
   appId: ""
   appSecret: ""
swagger: true
```
* [x] server：后端服务监听的地址和端口，保持默认。
* [x] externalUrl：对台对外提供的访问地址，格式为：`http[s]://<address>[:<port>]`。
* [x] secret: `CAS3.0`票据签名字符串，生产环境请务必进行修改。
* [x] mysql：`MySQL`数据库相关配置。
* [x] redis：`Redis`相关配置。
* [x] jwt：`JWT`相关配置。
* [x] mfa：双因素认证相关配置，`issuer`为APP扫码后显示的名称。
* [x] oss：对象存储相关配置，支持MinIO和华为云OBS。
* [ ] ldap：参考 [LDAP配置](#LDAP配置)，配置完成后需要将用户同步到本地后，用户方可登录。
* [ ] sms：参考 [短信配置](#短信配置)。
* [ ] mail：邮件相关配置，目前系统中未使用。
* [ ] dingTalk：钉钉自建应用配置，如果不需要钉钉扫码登录，可以忽略，参考 [钉钉配置](https://github.com/yuyan075500/ops-api/blob/main/deploy/dingtalk.md "钉钉配置")。
* [ ] wechat：企业微信自建应用配置，如果不需要企业微信扫码登录，可以忽略，参考 [企业微信配置](https://github.com/yuyan075500/ops-api/blob/main/deploy/wechat.md "企业微信配置")。
* [ ] feishu：飞书自建应用配置，如果不需要飞书扫码登录，可以忽略，参考 [飞书配置](https://github.com/yuyan075500/ops-api/blob/main/deploy/feishu.md "飞书配置")。
* [x] swagger：Swagger接口，生产环境不建议关闭。

> **注意**： `externalUrl`地址一经固定，切忽随意更改，会影响SSO的相关功能，如果更改后SSO客户端无法登录，那么你需要重置进行客户端配置。

## LDAP配置
支持与Windows AD或OpenLDAP进行对接，实现用户认证，使配置说明如下：
* [x] host：服务器地址，格式为：`ldap[s]://<host>:<port>`。
* [x] bindUserDN：绑定的用户DN，格式为：`cn=admin,dc=example,dc=cn`。
* [x] bindUserPassword：绑定的用户密码。
* [x] searchDN：搜索用户的DN，格式为：`ou=IT,dc=example,dc=cn`，支持配置多个DN，之间使用`&`分割，如：`ou=IT,dc=example,dc=cn&ou=HR,dc=example,dc=cn`。
* [x] userAttribute：用户属性，如果是OpenLDAP则为`uid`，如果是Windows AD则为`sAMAccountName`。
* [ ] maxPasswordAge：密码最大有效期，该参数仅针对 `Windows AD`，需要与实际的域控用户密码有效期保持一致。

> 说明：如果需要更改Windows AD或OpenLDAP的用户密码，则需要绑定的用户有足够的权限。如果是Windows AD还要求使用`ldaps`协议进行连接，`ldaps`协议的默认端口为`636`。

## 短信配置
目前短信支持华为云和阿里云，具体配置如下所示：
* [x] provider：指定短信服务商，固定值，`aliyun`或`huawei`。
* [x] url：短信服务地址，不同服务商的配置不同，阿里云参考[短信服务接入点](https://help.aliyun.com/zh/sms/developer-reference/api-dysmsapi-2017-05-25-endpoint "阿里云短信服务接入点")，华为云参考[API请求地址](https://support.huaweicloud.com/api-msgsms/sms_05_0000.html#section1 "API请求地址")。
* [x] appKey: 华为云参考[开发数据准备](https://support.huaweicloud.com/devg-msgsms/sms_04_0006.html "开发数据准备")，阿里云参考[创建AccessKey](https://help.aliyun.com/zh/ram/user-guide/create-an-accesskey-pair "创建AccessKey")。
* [x] appSecret: 华为云参考[开发数据准备](https://support.huaweicloud.com/devg-msgsms/sms_04_0006.html "开发数据准备")，阿里云参考[创建AccessKey](https://help.aliyun.com/zh/ram/user-guide/create-an-accesskey-pair "创建AccessKey")。
* [x] callbackUrl：短信回调地址，用于接收短信发送状态，仅华为云需要配置，回调地址为`<externalUrl>/api/v1/sms/huawei/callback`，如果平台的访问地址为内网地址，则无法接收回调信息。
* [x] resetPassword.sender：重置密码短信通道号，仅华为云需要配置。
* [x] resetPassword.templateId：重置密码短信模板ID。
* [x] resetPassword.signature：重置密码短信签名名称。
### 短信模板
阿里云模板如下所示：
```
# 重置密码短信模板
您的验证码为：${code}，验证码在5分钟内有效，请勿泄漏他人！
```
华为云模板如下所示：
```
# 重置密码短信模板
您的验证码为：${1}，验证码在5分钟内有效，请勿泄漏他人！
```
# 项目证书
为确保重要信息不会泄露，在项目部署时建议生成一套全新的证书，推荐使用 [证书在线生成工具](https://www.qvdv.net/tools/qvdv-csrpfx.html "在线生成工具") 创建。建议将证书有效期设置为10年，证书生成完成后需要下载CRT证书文件、证书公钥和证书私钥并严格按以下名称命名：
* private.key：私钥
* public.key：公钥
* certificate.crt：证书  

你也可以使用`openssl`工具生成自签证书，参考命令如下所示：
```shell
# 生成私钥
openssl genpkey -algorithm RSA -out private.key -pkeyopt rsa_keygen_bits:2048 -outform PEM
# 创建证书
openssl req -new -x509 -key private.key -out certificate.crt -days 3650
# 从证书中提取公钥
openssl rsa -in private.key -pubout -out public.key
```
# IP地址库
记录用户登录信息中的源IP来源于离线库文件，该文件位于项目`config/GeoLite2-City.mmdb`目录，最后更新日志为`2024-11-8`，最新库文件可从官方获取并替换即可。
