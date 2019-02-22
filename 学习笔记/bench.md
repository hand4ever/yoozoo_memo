# filebeat 到 logstash 压测

> 压测两种情况
>
> 1. 含证书情况下吞吐量
> 2. 不含证书情况下吞吐量

## 1. 准备工作

> * 安装好 filebeat、logstash，并分别启动之
>
> * filebeat 从 stdin 中获取输入，然后传输到 logstash 拆解，最后传输到 es 里，整个步骤完成

 ## 2. 测试项目

> * 连续测试 10 万数据，查看耗时
> * 连续测试 1 分钟，查看导量

## 3. 详细设计

### 3.0 elasticsearch(6.4.1) 、 kibana(6.4.1) 准备

> 使用 `docker` 安装，这里不做赘述
>
> ```shell
> # 执行命令如下：
> docker ps --format "{{.Image}}: {{.Command}}" --no-trunc
> # 输入结果如下：
> kibana:6.4.1: "/usr/local/bin/kibana-docker"
> elasticsearch:6.4.1: "/usr/local/bin/docker-entrypoint.sh eswrapper"
> 
> ```
>
> 

### 3.1 filebeat(6.4.1) 准备

> **step 1 filebeat 安装后位置**
>
> `/Users/panlong/yoozoo/opt/filebeat`
>
> **step 2 准备**
>
> **配置文件 fb_bench.yml**
>
> ```yaml
> filebeat.inputs:
> - type: stdin
> output.console:
>   pretty: true
> ```
>
> **源日志文件 access.log**
>
> `/Users/panlong/yoozoo/logs/access.log`



### 3.2 logstash(6.4.1) 准备

> **step 1 logstash 安装后位置**
>
> `/Users/panlong/yoozoo/opt/logstash`
>
> **step 2 准备**
>
> **配置文件 ls_bench.yml**
>
> ```
> input {
>   stdin {}
> }
> 
> filter {
>     grok {
>         match => { "message" => '%{IPORHOST:remote_addr} - %{DATA:remote_user} \[%{HTTPDATE:time}\] "%{WORD:request_action} %{DATA:request} HTTP/%{NUMBER:http_version}" %{NUMBER:response_code} %{NUMBER:bytes} "%{DATA:referrer}" "%{DATA:agent}" "%{DATA:xforward}"' }
>     }
> }
> 
> output {
>   stdout { codec => rubydebug }
> }
> 
> ```
>
> 
>
> **step 3 运行测试命令**
>
> `head -n 2 /Users/panlong/yoozoo/logs/access.log | bin/logstash -f config/ls_bench_test.yml`

### 3.3 filebeat 到 logstash 配置连通

> **step 0 openssl 配置**

```shell
# 编辑 openssl 文件
vim /usr/local/etc/openssl/openssl.cnf
## 在 v3_ca 下面添加 subjectAltName
[ v3_ca ]
subjectAltName = IP:127.0.0.1

# filebeat 端生成证书和私钥
openssl req -subj '/CN=127.0.0.1/' -x509 -days $((100 * 365)) -batch -nodes -newkey rsa:2048 -keyout pki/tls/private/filebeat.key -out pki/tls/certs/filebeat.crt

# logstash 端生成证书和私钥
openssl req -subj '/CN=127.0.0.1/' -x509 -days $((100 * 365)) -batch -nodes -newkey rsa:2048 -keyout pki/tls/private/logstash.key -out pki/tls/certs/logstash.crt

# 启动lgstash
bin/logstash -f config/ls_bench_cert_test.conf

# 测试 filebeat
./filebeat -c fb_bench_cert_test.yml test output

# 启动 filebeat 【-e -v 是为了看更详细的输出信息，可选】
./filebeat -c fb_bench_cert_test.yml -e -v


```



> **step 1 filebeat 配置**

```
demo

```

> **step 2 logstash 配置**

```
demo

```

> **step 4 启动 filebeat**

`tail -f /Users/panlong/yoozoo/logs/access.log | ./filebeat -c fb_bench.yml`

> **step 5 测试文件编写如下

```php


```

