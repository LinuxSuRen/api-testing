# Make SAN certs

Create a SAN certificate here, following these steps

#### 1. 生成私钥

```shell
openssl genrsa -out server.key 2048
```

#### 2. 生成证书（会提示即可，不必填写）

```shell
openssl req -new -x509 -key server.key -out server.crt -days 36500
国家名字
Country Name (2 letter code) [AU]:CN
省份全名
State or Province Name (full name) [Some-State]:GuangDong
城市名
Locality Name (eg, city) []:Meizhou
组织名
Organization Name (eg, company) [Internet Widgits Pty Ltd]:Xuexiangban
组织单位名
Organizational Unit Name (eg, section) []:go
服务器or用户的名字
Common Name (e.g. server FQDN or YOUR name) []:kuangstudy
邮箱地址
Email Address []:24736743@qq.com
```

#### 3.生成csr

```shell
openssl req -new -key server.key -out server.csr
```

#### 4.配置openssl.cfg

```shell
1) 查找openssl在服务器的安装目录并且找到openssl.cnf并将其所在的目录
2) 编辑[ CA_default ] ，打开 copy_extensions = copy
3) 找到[ req ]，打开 req_extensions = v3.req # The extensions to add to a
certificate request
4) 找到[ v3_req ]，添加 subjectAltName = @alt_names
5) 找到最底部的 [ alt_names ]，修改名字和
DNS.1 = *.your.com
```

#### 5.生成本地私钥test.key

```shell
openssl genpkey -algorithm RSA -out test.key
```

#### 6.根据私钥生成csr请求文件test.csr

```shell
openssl req -new -nodes -key test.key -out test.csr -days 3650 -subj "/C=cn/O=myorg/OU=mycomp/CN=myname" -config ./openssl.cnf -extensions v3.req
```

#### 7.生成ca证书 pem

```shell
openssl x509 -req -days 365 -in test.csr -out test.pem -CA server.crt -CAkey server.key -CAcreateserial -extfile ./openssl.cnf -extensions v3.req
```

