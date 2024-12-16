# Zabbix 单点登录
Zabbix 支持的单点登录方式：SAML2
## 配置方法
1. **创建密钥和证书**：可以使用 [在线生成工具](https://www.qvdv.net/tools/qvdv-csrpfx.html "在线生成工具")。建议证书有效期设置为10年，不设置密码，生成完成后需要下载 CRT 证书和私钥并按以下名称命名：<br><br>
   * sp.key：私钥。
   * sp.crt：证书。<br><br>
2. **获取 IDP 证书**：IDP 的证书的存放路径为项目的 `config/certs/certificate.crt`，需要将此证书下载并保存为 `idp.crt`。<br><br>
3. **上传密钥和证书**：将 `sp.key`、`sp.crt`、`idp.crt` 上传到 Zabbix 站点部署的 `ui/conf/certs/` 目录下，除非 `zabbix.conf.php` 中提供了自定义路径，否则 Zabbix 默认在 `ui/conf/certs/` 路径中查找文件。<br><br>
4. **Zabbix 单点登录配置**：登录到 Zabbix，进入【认证】配置界面，如下图所示：<br><br>
![img.png](img/zabbix-config.jpg)<br><br>
   * 启用 SAML 身份验证：选中复选框以启用 SAML 身份验证。
   * IDP 实体 ID：IDP 的唯一标识符，此处为 `http[s]://<address>[:<port>]`。
   * SSO 服务 URL：用户登录时被重定向到的 URL，此处为：`http[s]://<address>[:<port>]/login`。
   * Username attribute：固定值 `username`。
   * SP entity ID：通常为 Zabbix 的访问地址，如：`http[s]://<address>[:<port>]`。
   * SP name ID format：固定值 `urn:oasis:names:tc:SAML:2.0:attrname-format:unspecified`。<br><br>
其它选项按图示配置即可，也可以参考 [官方文档](https://www.zabbix.com/documentation/6.0/zh/manual/web_interface/frontend_sections/administration/authentication#advanced-settings "官方文档") 进行其它配置。<br><br>
5. **站点注册**：登录到 IDSphere 统一认证平台，点击【资产管理】-【站点管理】-【新增】将 Zabbix 站点信息注册到 IDSphere 统一认证平台，配置如下所示：<br><br>
![img.png](img/zabbix-config.jpg)<br><br>
   * 站点名称：指定一个名称，便于用户区分。
   * 登录地址：Zabbix 的登录地址。
   * SSO 认证：启用。
   * 认证类型：选择 `SAML2`。
   * 站点描述：描述信息。
   * SP EntityID：Zabbix的SP EntityID，与 Zabbix【认证】配置界面中的保持一致。
   * SP 证书：将 `sp.crt` 文件内容粘贴到此处。
