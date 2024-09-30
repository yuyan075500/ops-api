# 项目部署
项目支持使用 [Docker Compose一键部署](#docker-compose部署) 和 [Kubernetes部署](#Kubernetes部署)，生产环境推荐使用Kubernetes部署。
## Docker Compose部署
如果你想快速拥有一个简易的环境用于测试、演示，对性能、稳定性以及安全性没有任何求的，那么推荐使用该部署方式。  
1. **部署环境准备**：你需要准备一台Linux服务器，并安装以下组件。
   * [x] Docker。
   * [x] Docker Compose。
   * [ ] MySQL 8.0。
   * [ ] Redis 5.x。
   * [ ] MinIO。

    `Docker`和`Docker Compose`是部署环境必须的，其它的都可以使用`docker-compose.yaml`指定的，也可以使用独立的。
2. **克隆项目**：
    ```shell
    git clone https://github.com/yuyan075500/ops-api.git
    ```
3. **切换工作目录**：
    ```shell
    cd ops-api/deploy/docker-compose
    ```
4. **环境变量配置**：修改`.env`文件中环境变量，如果你使用`docker-compose.yaml`指定的`MySQL`、`Redis`、`MinIO`，则可以跳过此步骤。

   > **注意**：如果有使用[钉钉](https://github.com/yuyan075500/ops-api/blob/main/deploy/dingtalk.md "钉钉配置")、[企业微信](https://github.com/yuyan075500/ops-api/blob/main/deploy/wechat.md "企业微信配置")或[飞书](https://github.com/yuyan075500/ops-api/blob/main/deploy/feishu.md "飞书配置")扫码认证，请按要求对前端项目进行单独构建打包，并修改`.env`文件对中应前端的镜像配置。

5. **项目配置**：修改`conf/config.yaml`文件中相关配置，请参考 [配置文件说明](#配置文件说明)。

   > **注意**：MinIO的`accessKey`和`secretKey`需要在部署成功后登录进MinIO控制台手动创建，确保与配置文件中指定的值相同即可。

6. **证书**：[创建项目证书](#项目证书)，将生成的新证书保存至`certs`目录中并覆盖目标文件。如果是测试环境你也可以跳过此步骤使用项目自带的证书。
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

   > **注意**：如果使用的外部数据库，请确保数据库使用的字符集为`utf8mb4`，排序规则为`utf8mb4_general_ci`。

10. **系统登录**：部署完成后，系统会自动创建一个超级用户，此用户不受Casbin权限控制。用户名为：`admin`，密码为：`admin@123...`。
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
sms:
  url: "https://smsapi.cn-north-4.myhuaweicloud.com:443/sms/batchSendDiffSms/v1"
  appKey: ""
  appSecret: ""
  callbackUrl: "<externalUrl>/api/v1/sms/callback"
  verificationCode:
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
* [ ] sms：参考 [短信配置](#LDAP配置)。
* [ ] mail：邮件相关配置，目前系统中未使用。
* [ ] dingTalk：钉钉自建应用配置，如果不需要钉钉扫码登录，可以忽略，参考 [钉钉配置](https://github.com/yuyan075500/ops-api/blob/main/deploy/dingtalk.md "钉钉配置")。
* [ ] wechat：企业微信自建应用配置，如果不需要企业微信扫码登录，可以忽略，参考 [企业微信配置](https://github.com/yuyan075500/ops-api/blob/main/deploy/wechat.md "企业微信配置")。
* [ ] feishu：飞书自建应用配置，如果不需要飞书扫码登录，可以忽略，参考 [飞书配置](https://github.com/yuyan075500/ops-api/blob/main/deploy/feishu.md "飞书配置")。
* [x] swagger：Swagger接口，生产环境不建议关闭。

> **注意**： `externalUrl`地址一经固定，切忽随意更改，会影响SSO的相关功能，如果更改后SSO客户端无法登录，那么你需要重置进行客户端配置。

## LDAP配置
平台用户支持与Windows AD或OpenLDAP进行对接，实现用户认证，使配置说明如下：
* [x] host：服务器地址，格式为：`ldap[s]://<host>:<port>`。
* [x] bindUserDN：绑定的用户DN，格式为：`cn=admin,dc=example,dc=cn`。
* [x] bindUserPassword：绑定的用户密码。
* [x] searchDN：搜索用户的DN，格式为：`ou=IT,dc=example,dc=cn`，支持配置多个DN，之间使用`&`分割，如：`ou=IT,dc=example,dc=cn&ou=HR,dc=example,dc=cn`。
* [x] userAttribute：用户属性，如果是OpenLDAP则为`uid`，如果是Windows AD则为`sAMAccountName`。

> 说明：如果需要更改Windows AD或OpenLDAP的用户密码，则需要绑定的用户有足够的权限，Windows AD还要求使用`ldaps`协议进行连接。

## 短信模板
目前仅支持华为云短信服务（MSGSMS），需要在华为云开通短信服务。短信将用于用户自助密码修改，不使用该功能则可以忽略，短信模板如下所示：
```
${1}您好，您的校验码为：${2}，校验码在${3}分钟内有效，保管好校验码，请勿泄漏！
```
短信模板三个变量，分别代表用户名、校验码和校验码有效时间。
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
记录用户登录信息中的源IP来源于离线库文件，该文件位于项目`config/GeoLite2-City.mmdb`目录，最后更新日志为`2024-07-23`，最新库文件可从官方获取并替换即可。
