# OpenLDAP 注意事项
OpenLDAP 账号需要基于 `inetOrgPerson` 和 `shadowAccount` 两个类进行配置，参考配置如下：
```shell
objectClass: inetOrgPerson
objectClass: shadowAccount
uid: lisi
givenName: 李
sn: 四
displayName: 李四
shadowMin: 0
shadowWarning: 0
shadowInactive: 0
cn: 李四
shadowExpire: 99999
shadowMax: 90
mail: lisi@idsphere.cn
mobile: 132****8888
userPassword: {SHA512}5r**********is2EgNEu61tQ==
shadowLastChange: 20068
```
* uid：账号唯一标识，建议使用用户名，如：`lisi`。
* givenName：姓，如：`李`。
* sn：名，如：`四`。
* displayName：显示姓名，如：`李四`。
* shadowMin：从上次修改密码后，多久可再次修改密码，0 表示不限制。
* shadowWarning：密码过期前多久开始提示，0 表示不提示。
* shadowInactive：密码过期后还可以登录的天数，0 表示过期后不允许登录。
* cn：通常与`givenName`和`sn`组合而成，如：`李四`。
* shadowExpire：账号过期时间，该值表示距离 1970-01-01 的天数，其中 99999 表示永不过期。
* shadowMax：密码过期最大天数，如：90。
* mail：邮箱地址，如：`lisi@idsphere.cn`。
* mobile：手机号，如：`132****8888`。
* userPassword：用户密码，使用`IDSphere`平台修改用户密码，默认采用`SHA1`加密。
* shadowLastChange：最后一次更改密码的时间，该值表示距离 1970-01-01 的天数。<br><br>

将 OpenLDAP 用户同步到 IDSphere 统一认证平台，用户映射规则如下：
```shell
{
	"name": "cn",
	"username": "uid",
	"email": "mail",
	"phone_number": "mobile"
	"password_expired_at": "shadowMax"
	"is_active": "shadowExpire"
}
```
# Windows AD 注意事项
Windows AD 用户使用 `sAMAccountName` 作为用户名，相关字段映射规则如下：
```shell
{
	"name": "cn",
	"username": "sAMAccountName",
	"email": "mail",
	"phone_number": "mobile",
	"is_active": "userAccountControl"
}
```
由于用户密码过期时间受域控的策略影响，通过 API 无法获取，所以需要在项目配置文件中手动指定。
# 登录策略
| OpenLDAP 或 Windows AD | IDSphere 统一认证平台 | 是否可以登录 |
|-----------------------|-----------------|--------|
| 账号状态禁用                | 账号状态为禁用         | 否      |
| 账号状态禁用                | 账号状态为启用         | 否      |
| 账号状态启用                | 账号状态为禁用         | 否      |
| 账号状态启用                | 账号状态为启用         | 是      |
| 账号密码过期                | 账号密码过期          | 否      |
| 账号密码过期                | 账号密码未过期         | 否      |
| 账号密码未过期               | 账号密码未过期         | 否      |
| 账号密码未过期               | 账号密码未过期         | 是      |

只有当 OpenLDAP 或 Windows AD 中的账号状态为启用和账号密码未过期，并且需要本地对应的账号同时满足，才允许登录。
# 同步策略
用户同步判断依据为 `username`，具体地同步规则如下：

| OpenLDAP 或 Windows AD | IDSphere 统一认证平台 | 动作     |
|-----------------------|-----------------|--------|
| 有                     | 没               | 创建     |
| 有                     | 有               | 更新     |
| 没有                    | 有               | 不作任何动作 |