# 详细设计

## 各部分介绍

### 1. 配置服务器

此服务器包含 CA 认证中心和一个配置中心

####  CA 认证中心

提供CA认证和证书服务，服务器上存储根证书，扶着签发所有用于 filebeat 和 logstash 通信的证书

1. 初始化 CA 认证中心
```bash
# 1. 创建根密钥
# umask 077，确保创建的证书权限正确
umask 077;openssl genrsa -out /etc/pki/CA/private/cakey.pem 2048
# 2. 生成根证书
umask 077;openssl req -new -x509 -key /etc/pki/CA/private/cakey.pem -out /etc/pki/CA/cacert.pem
# 3. 初始化目录
touch /etc/pki/CA/newcerts
touch /etc/pki/CA/index.txt
touch /etc/pki/CA/serial
echo "01" > /etc/pki/CA/serial
```

2. 为采集器颁发证书

```bash
(umask 077; openssl genrsa -out ori/logstash.key 1024  -extfile ori/logstash.conf)
# 生成 pkcs8 格式证书
openssl pkcs8 -topk8 -inform PEM -in ori/logstash.key -outform PEM -nocrypt > ori/logstash.pem

openssl req -new -key ori/logstash.pem -out ori/logstash.csr

openssl ca -in ori/logstash.csr -out ori/logstash.cst -days 3650



(umask 077; openssl genrsa -out ori/filebeat.key 1024  -extfile ori/filebeat.conf)
# 生成 pkcs8 格式证书
openssl pkcs8 -topk8 -inform PEM -in ori/filebeat.key -outform PEM -nocrypt > ori/filebeat.pem

openssl req -new -key ori/filebeat.pem -out ori/filebeat.csr

openssl ca -in ori/filebeat.csr -out ori/filebeat.cst -days 3650
```

#### 配置中心

1. 创建应用
提供应用的基础信息
提供项目注册功能和配置功能，项目接入日志系统，需在此处填写基本信息，配置中心通过CA 认证中心和用户填写的信息生成证书

2. 提供应用信息相关接口
可以通过私密的链接获取对应的应用的信息、证书、配置等等

3. 自动同步配置信息到 logstash

4. 配置检测功能，本机上应该有同版本的 logstash，检测配置文件的正确性

### 2. 项目服务器

项目服务器上选装采集工具，配置的工具信息是需要用到配置中心生成的证书和项目信息

1. 提供 shell 脚本的安装方式
```bash
curl 'http://log.youzu.com' | bash
```
2. 自行安装
```bash
# 安装 filebeat
> yum install filebeat
# 下载证书文件
> curl "http://log.youzu.com/cert.zip" 
# 修改配置文件
> vim /etc/filebeat/filebeat.yml
# 启动filebeat
```
3. 后续提供 docker 安装的方式


### 3. logstash 日志过滤服务器

接收项目服务器上报的信息，跟据在配置中心配置的信息，对项目的日志进行分类、过滤、格式化，最终存入 ES 集群

1. 接收来自采集端的日志信息
2. 区分来源，打上标记，进行各自的过滤规则
3. 创建索引，存入 ES
4. 定时同步过滤规则

### 4. ES 集群

存储所有项目格式化后的日志信息
