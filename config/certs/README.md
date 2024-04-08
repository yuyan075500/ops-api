# 使用Linux系统生成密钥对的方法如下
1. 生成一个自签名的私钥和公钥对。
```shell
openssl genrsa -out private.key 2048
openssl rsa -in private.key -pubout -out public.key
```
2. 将私钥和公钥转换为`PEM`格式。
```shell
openssl rsa -in private.key -outform PEM -out private.pem
openssl rsa -pubin -in public.key -outform PEM -out public.pem
```