# elk 学习笔记

## 0. 安装 filebeat + elk

### 0.1 Filebeat docker 安装

#### 背景

1. 依赖镜像：centos 7

2. 开源、免费使用

3. 参考

   > https://www.elastic.co/guide/en/beats/filebeat/6.5/running-on-docker.html

#### 安装命令

> ```sh
> docker pull docker.elastic.co/beats/filebeat:6.5.4
> ```

#### 运行命令

> ```sh
> docker run \
> docker.elastic.co/beats/filebeat:6.5.4 \
> setup -E setup.kibana.host=kibana:5601 \
> -E output.elasticsearch.hosts=["elasticsearch:9200"]
> 
> ```



### 0.2 ElasticSearch docker 安装

#### 背景

1. 参考

  > [1. 全文搜索引擎 Elasticsearch 入门教程](http://www.ruanyifeng.com/blog/2017/08/elasticsearch.html)
  >
  > [2. Install Elasticsearch with Docker](https://www.elastic.co/guide/en/elasticsearch/reference/current/docker.html)

#### 安装命令

> ```sh
> docker pull docker.elastic.co/elasticsearch/elasticsearch:5.6.9
> ```

#### 运行命令

> ```sh
> docker run -d -p 9200:9200 -p 9300:9300 -e "discovery.type=single-node" docker.elastic.co/elasticsearch/elasticsearch:5.6.9
> ```

#### 相关其他操作
```sh
docker ps 
docker exec -it beda5687cbe1 /bin/bash # 进入容器的centos命令控制台
./bin/elasticsearch-plugin remove x-pack # 卸载X-Pack

# 卸载后需要重启es docker容器： 先 stop 容器，再 start 容器
curl –user elastic:changeme -X GET 'http://10.7.80.194:9200'

# 安装 es head
docker pull mobz/elasticsearch-head:5 # 拉取镜像
docker run -p 9100:9100 mobz/elasticsearch-head:5 #运行容器

# 进入 es 容器修改信息
docker exec -it 6b6d4328ceed /bin/bash

vi elasticsearch/config/elasticsearch.yml

#添加如下配置然后重启elasticsearch
http.cors.enabled: true
http.cors.allow-origin: "*"

# 重启 es 容器
docker restart "es 容器的id"
http://localhost:9100

# 插入测试数据
curl -XPOST 'http://localhost:9200/indextest/tabletest/1?pretty' -d' {"name": "张三", "timestamp":1548760056 }'
```

### 0.3 Kibana docker 安装

#### 背景

1. 参考

   > [1. Docker 容器中运行 Kibana](https://www.elastic.co/guide/cn/kibana/current/docker.html)

#### 安装命令

> ```sh
> docker pull docker.elastic.co/kibana/kibana:6.0.0 [慢]
> 
> docker pull kibana:6.6.0 [快]
> 
> # remove x-pack 插件
> docker exec -it beda5678dd1 bash
> bin/kibana-plugin remove x-pack
> ```



## 1. 运行 elasticsearch 和 kibana

```sh
# 运行 elasticsearch 基于docker
docker run --name myesv3 -d -p 9200:9200 -p 9300:9300 -p 5601:5601 elasticsearch:6.4.1
# 运行 kibana 基于 docker
docker run -it -d -e ELASTICSEARCH_URL=http://127.0.0.1:9200 --name mykibanav2 --network=container:myesv1 kibana:6.4.1

```

> 在浏览器输入 `http://127.0.0.1:9200` 测试 es 安装是否成功
>
> 在浏览器输入 `http://127.0.0.1:5601` 测试 kibana 安装是否成功

## 2. kibana Dev tools 使用

### 2.1 插入数据

```
POST accounts/person/1
{
    "name": "zhangsan1",
    "gender": "male",
    "desc": "a boy test bla"
}

```

### 2.2 更新数据

```
POST accounts/person/1/_update
{
	"doc": {
        "desc": "greeting!!!"
	}
}

```

### 2.3 删除数据

```
DELETE accounts/person/1

```

### 2.4 查询的两种形式

1. Query String

    ` GET accounts/person/_search?q=zhangsan`

2. Query DSL

    ```
    GET accounts/person/_search
    {
        "query":{
            "term":{
                "name":"zhangsan"
            }
        }
    }
    
    ```

### 2.5 列出所有索引

`GET _cat/indices`

## 3. filebeat 使用

### 3.1 准备日志文件

`/Users/panlong/yoozoo/logs/access.log`

### 3.2 配置文件 `filebeat.yml`

> **input配置简介**
>
> `prospectors` 是个数组，可以有多个 通过 `-` 来配置多个数组元素
>
> `input_type` 输入类型，比如 log（日志文件） stdin（标准输入）
>
> `path` 路径，可配置多个
>
> **output配置简介**
>
> `console`
>
> ​    output.console:
>
> ​        pretty: true
>
> `elasticsearch`
>
> `logstash`
>
> `kafka`
>
> `redis`
>
> `file`

> demo 如下

```yml
filebeat.config:
  prospectors:
    path: ${path.config}/prospectors.d/*.yml
    reload.enabled: false
  modules:
    path: ${path.config}/modules.d/*.yml
    reload.enabled: false

processors:
- add_cloud_metadata:

output.elasticsearch:
  hosts: ['elasticsearch:9200']
  username: elastic
  password: changeme
  
```

### 3.3 测试

**step1复制一份 nginx.yml 文件，内容为如下**

```yml
filebeat.inputs:
- type: stdin
output.console:
  pretty: true
```



**step2 测试 filebeat 采集日志**

`head -n2 /Users/panlong/yoozoo/logs/access.log | /Users/panlong/yoozoo/opt/filebeat/filebeat -c nginx.yml`

## 4. etcd 安装使用

> [etcd 文档](https://etcd.readthedocs.io/en/latest/)

### 4.1 etcd 安装

> **step 1 Docker 安装**

###### Docker

etcd uses [`gcr.io/etcd-development/etcd`](https://gcr.io/etcd-development/etcd) as a primary container registry, and [`quay.io/coreos/etcd`](https://quay.io/coreos/etcd) as secondary.

```shell
rm -rf /tmp/etcd-data.tmp && mkdir -p /tmp/etcd-data.tmp && \
  docker rmi quay.io/coreos/etcd:v3.3.12 || true && \
  docker run \
  -p 2379:2379 \
  -p 2380:2380 \
  --mount type=bind,source=/tmp/etcd-data.tmp,destination=/etcd-data \
  --name etcd-gcr-v3.3.12 \
  gcr.io/etcd-development/etcd:v3.3.12 \
  /usr/local/bin/etcd \
  --name s1 \
  --data-dir /etcd-data \
  --listen-client-urls http://0.0.0.0:2379 \
  --advertise-client-urls http://0.0.0.0:2379 \
  --listen-peer-urls http://0.0.0.0:2380 \
  --initial-advertise-peer-urls http://0.0.0.0:2380 \
  --initial-cluster s1=http://0.0.0.0:2380 \
  --initial-cluster-token tkn \
  --initial-cluster-state new

docker exec etcd-gcr-v3.3.12 /bin/sh -c "/usr/local/bin/etcd --version"
docker exec etcd-gcr-v3.3.12 /bin/sh -c "ETCDCTL_API=3 /usr/local/bin/etcdctl version"
docker exec etcd-gcr-v3.3.12 /bin/sh -c "ETCDCTL_API=3 /usr/local/bin/etcdctl endpoint health"
docker exec etcd-gcr-v3.3.12 /bin/sh -c "ETCDCTL_API=3 /usr/local/bin/etcdctl put foo bar"
docker exec etcd-gcr-v3.3.12 /bin/sh -c "ETCDCTL_API=3 /usr/local/bin/etcdctl get foo"
```

### 4.2 etcd 运行、测试demo

## 5. confd 安装使用

### 5.1 confd 安装

```shell
# 下载
wget https://github.com/kelseyhightower/confd/releases/download/v0.16.0/confd-0.16.0-darwin-amd64
# 编译安装
mkdir -p ~/opt/confd/bin
mv confd-0.16.0-linux-amd64 ~/opt/confd/bin/confd
chmod +x ~/opt/confd/bin/confd
export PATH="$PATH:~/opt/confd/bin"

```



### 5.2 confd 使用

> **step 1 新建配置目录（Create the confdir）**

`mkdir -p ~/etc/confd/{conf.d,templates}`

> **step 2 新建模板配置文件（Create a template resource config）**
>
> 语法是 toml 配置规则，文件在 conf.d
>
> `touch ~/etc/confd/conf.d/myconfig.toml` 内容如下:

```toml
[template]
src = "myconfig.conf.tmpl"
dest = "/tmp/myconfig.conf"
keys = [
    "/myapp/database/url",
    "/myapp/database/user",
]

```

> **step 3 新建源模板文件（Create the source template）**
>
> 语法是 [Golang text templates](http://golang.org/pkg/text/template/#pkg-overview)
>
> `touch ~/etc/confd/templates/myconfig.conf.tmpl`

```toml
[template]
src = "myconfig.conf.tmpl"
dest = "/tmp/myconfig.conf"
keys = [
    "/myapp/database/url",
    "/myapp/database/user",
]
```

> **step 4 运行模板（Process the template）**
>
> **基于etcd的运行命令**
>
> `confd -onetime -backend etcd -node http://127.0.0.1:2379`

> **step 5 查看配置文件**
>
> `cat ~/tmp/myconfig.conf`
>
> *输出*

```


```



## 6. 证书相关

> 参考：https://www.jianshu.com/p/79c284e826fa
>
> mac 本机的 openssl 文件
>
> `vim /usr/local/etc/openssl/openssl.cnf`

### 6.1 建立 CA
1. 建目录
```
mkdir -p ./demoCA/{private,newcerts} && \
    touch ./demoCA/index.txt && \
    touch ./demoCA/serial && \
    echo 01 > ./demoCA/serial
```
2. 生成 CA 根密钥
`openssl genrsa -des3 -out ./demoCA/private/cakey.pem 2048 #可以去掉des，以后签证书不用输密码`

> 输入密码：123456

3. 生成 CA 证书请求

`openssl req -new -days 3650 -key ./demoCA/private/cakey.pem -out careq.pem`

>输出密码：Panlong123


4. 自签发 CA 根证书
`openssl ca -selfsign -in careq.pem -out ./demoCA/cacert.pem`

5. 以上合二为一(上面两步不做，直接可以做这步）
`openssl req -new -x509 -days 3650 -key ./demoCA/private/cakey.pem -out ./demoCA/cacert.pem`

> 到这里，我们已经有了自己的 CA 了，下面我们开始为用户颁发证书。


### 6.2 为用户颁发证书

1. 生成用户 RSA 密钥

`openssl genrsa -des3 -out userkey.pem 2048 # 4096`

> 输入密码：123456


2. 生成用户证书请求
```
# NO SAN 实际执行这个 输入密码：12345678
openssl req -new -days 365 -key userkey.pem -out userreq.pem # 
# with SAN
openssl req -new -days 365 -key userkey.pem -out userreq.pem -config /usr/local/etc/openssl/openssl.cnf
```

3. 使用 CA 签发证书
```
# NO SAN 实际执行这个 输入密码：12345678
openssl ca -in userreq.pem -out usercert.pem 
# with SAN
openssl ca -in userreq.pem -out usercert.pem -config /usr/local/etc/openssl/openssl.cnf -extensions v3_req
```
### 6.3 查看证书内容
`openssl x509 -in usercert.pem -text -noout`











