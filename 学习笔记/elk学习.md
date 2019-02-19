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

## 3. filebeat 使用

### 3.1 准备日志文件

`/Users/panlong/yoozoo/logs/access.log`

















