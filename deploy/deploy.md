# 项目部署
项目支持使用[Docker Compose一键部署](#docker-compose部署)和[Kubernetes部署](#Kubernetes部署)。
## Docker Compose部署
如果你想快速拥有一个简易的环境用于测试、演示，且对性能、稳定性以及安全性没有任何求的，那么推荐使用该部署方式。  
1. **准备部署环境**：你需要准备一台Linux服务器，并安装以下相关组件。
* [x] Docker。
* [x] Docker Compose。
* [ ] MySQL 8.0。
* [ ] Redis 5.x。
* [ ] MinIO。  
`Docker`和`Docker Compose`是环境必须的，其它的都可以使用配置文件自带的，也可以使用独立的`MySQL`、`Redis`、`MinIO`的。
> 说明：如果使用了独立的`MySQL`、`Redis`和`MinIO`，在执行部署的时候也会部署自带的版本。如果你不想部署自带的版本，删除`docker-compose.yaml`文件中相关的配置即可。`
2. **克隆项目**：将项目克隆到服务器中。
    ```shell
    git clone https://github.com/yuyan075500/ops-api.git
    ```
3. **进入部署目录**：切换至`Docker Compose`部署目录。
    ```shell
    cd ops-api/deploy/docker-compose
    ```
4. **修改环境变量**：修改`.env`文件中的相关配置，如果你使用了独立的`MySQL`、`Redis`、`MinIO`，那么可以跳过此步骤。
5. **修改项目配置**：修改`conf/config.yaml`配置文件，如果使用了独立的`MySQL`、`Redis`、`MinIO`，请确保配置文件中的相关连接信息正确，参考[配置说明](#配置文件说明)。
> 注意：配置文件中`secret`项需要更改成随机的字符串，用于CAS3.0票据签名。`OSS`的`accessKey`和`secretKey`可以先随机生成，在第Minio部署完成后登录后创建即可。
6. **创建证书**：创建[项目证书](#项目证书)，将生成的新证书保存至`certs`目录中并覆盖目标文件。如果是测试环境你也可以跳过此步骤使用项目自带的证书，但在生产环境中不推荐如此使用。
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
10. **系统登录**：部署完成后，系统会自动创建一个超级用户，此用户不受Casbin权限控制。用户名为：`admin`，密码为：`admin@123...`。
## Kubernetes部署（生环境环境推荐）
在Kubernetes中部署，需要用到Helm，请确保已安成[Helm安装](https://helm.sh/docs/intro/install/#from-the-binary-releases "Helm安装")。
### 运行环境准备
在Kubernetes中部署需要独立准备额外的资源，包含：
* [x] MySQL 8.0。
* [x] Redis 5.x。
* [x] MinIO。
### 部署
1. **克隆项目**：将项目克隆到服务器中。
    ```shell
    git clone https://github.com/yuyan075500/ops-api.git
    ```
2. **进入部署目录**：切换至`Docker Compose`部署目录。
    ```shell
    cd ops-api/deploy/kubernetes
    ```
3. **创建证书**：创建[项目证书](#项目证书)，证书创建完成后使用新的证书替换`templates/configmap.yaml`文件中对应的配置项。
4. **修改项目配置**：修改`templates/configmap.yaml`文件中`config.yaml`的相关配置，请参考[配置说明](#配置文件说明)。
5. **部署**：如果你需要同步创建`ingress`资源，那么需要在执行`helm`命令部署前修改`values.yaml`文件中的对应的配置项，**推荐同步创建**。如果你使用Kubernetes之外的代理程序，那么你需要将`Service`类型修改为`NodePort`，并参考`templates/ingress.yaml`模板文件中的转发规则进行相关配置。
   ```shell
   helm install <自定义应用名> --namespace <名称空间> .
   ```
7. **数据初始化**：将`deploy/data.sql`SQL中的数据导入到数据库中。
8. **系统登录**：部署完成后，系统会自动创建一个超级用户，此用户不受Casbin权限控制。用户名为：`admin`，密码为：`admin@123...`。
> 说明：应用的高可用只需调整应用的副本数即可，数据库和中间件的高可用需要自行完成。
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
swagger: true
```
* [x] server：服务监听的地址和端口。
* [x] externalUrl：对外提供的访问地址，格式为：`<protocol>://<address>[:<port>]`。
* [x] secret: `CAS3.0`票据签名字符串。
* [x] mysql：`MySQL`数据库相关配置。
* [x] redis：`Redis`相关配置。
* [x] jwt：`JWT`相关配置。
* [x] mfa：双因素认证相关配置，`issuer`为手机APP扫码后显示的名称。
* [x] oss：`Minio`对象存储相关配置。
* [ ] ldap：`AD`相关配置。
* [ ] sms：短信相关配置，目前仅支持华为云，需要在华为云开通短信服务，并配置[短信模板](#短信模板)。
* [ ] mail：邮件相关配置。
* [x] swagger：Swagger接口，如果是生产环境不建议开启。
    > 注意： `externalUrl`地址一经固定，切忽随意更改，如果有使用SSO的相关功能，那么客户端可能会受此影响无法登录，你需要重置进行配置。
# 项目证书
为确保重要信息不会泄露，在项目部署时建议生成一套全新的证书，推荐使用[证书在线生成工具](https://www.qvdv.net/tools/qvdv-csrpfx.html "在线生成工具")创建。建议将证书有效期设置为10年，证书生成完成后需要下载CRT证书文件、证书公钥和证书私钥并严格按以下名称命名：
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
# 短信模板
目前短信用于用户自助密码修改，不使用该功能则可以忽略，短信模板如下所示：
```
${1}您好，您的校验码为：${2}，校验码在${3}分钟内有效，保管好校验码，请勿泄漏！
```
短信模板需要包含三个变量，分别代表用户名、校验码和校验码有效时间，其它文字可以自定义。
# IP地址库
记录用户登录信息中的源IP来源于离线库文件，该文件位于项目`config/GeoLite2-City.mmdb`目录，最后更新日志为`2024-07-23`，最新库文件可从官方获取并替换即可。
