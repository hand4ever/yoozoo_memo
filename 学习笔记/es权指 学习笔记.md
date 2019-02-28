# 《elasticsearch 权威指南》 学习笔记
## 1. 基本概念入门
### 1.1 base 概念
> **索引（indexing）** 存储数据的行为叫做索引
> **类型（type）** 
> **文档（document）**
> **列（field）**

```
DB: Databases -> Tables -> Rows -> Columns
ES: Indices -> Types -> Documents -> Fields
```
> ES 集群可以包含多个索引（indices）（数据库），每个索引可以包含多个类型（types）（表），每个类型可包含多个文档（documents）（行），然后每个文档可包含多个字段（fields）（列）

### 1.2 入门实操
> **case 1 插入数据**
> **Q:需要进行的操作**
> * 为每个员工的文档（docuemnt）建立索引，每个文档包含了相应员工的所有信息。
> * 每个文档的类型为 `employee`。
> * `employee` 类型归属于索引 `megacorp`。
> * `megacorp` 索引存储在 Elasticsearch 集群中。
> 
> **这里插入 3 个文档**
> 
> **A:命令如下：**

```json
PUT /megacorp/employee/1
{
    "first_name": "John",
    "last_name": "Smith",
    "age": 25,
    "about": "I love to go rock climbing",
    "interests": ["sports", "music"]    
}
```

> **case 2 检索数据**
> 
> `GET /megacorp/employee/1`
> 从结果看出原始的 JSON 文档包含在 _source 字段中。
> 
> *我们通过 HTTP 方法 GET 来检索文档，同样的，可以用 DELETE 来删除文档，使用 HEAD 来检查某文档是否存在，如果只是想更新已存在的文档，我们只需要再 PUT 一次*
> 
> *更新操作如下*

```json
POST /megacorp/employee/3/_update
{
  "doc": {
    "age": 32
  }
}
```
> 
> **case 3 简单搜索**
> 
> `GET /megacorp/employee/_search`
> 搜索全部员工
> 
> *响应内容的 hits 数组中包含了我们所有的 3 个文档，默认情况下搜索会返回前 10 个结果*
> 
> `GET /megacorp/employee/_search?q=last_name:Smith`
> 搜索 last_name 中包含 “Smith” 的员工，上面是 `query string` 式搜索，也可以用 DSL 语句来查询。

```json
GET /megacorp/employee/_search
{
    "query": {
        "match": {
            "last_name": "Smith"
        }
    }
}
```
> **case 4 更复杂的搜索**

```json
GET /megacorp/employee/_search
{
  "query": {
    "bool": {
      "must": {
          "match": {
            "last_name": "Smith"
          }
      },
      "filter": {
        "range": {
          "age": {"gt":30}
        }
      }
    }
  }
}
```

> 以上为 搜索 last_name 为 Smith 的员工，并且是年龄大于 30 的员工。这里主要是用了一个过滤器（filter）（书中是es5的语法，这里写的是es6的语法）

> **case 5 更高级的搜索：全文搜索**

```json
GET /megacorp/employee/_search
{
  "query": {
    "match": {
      "about": "rock climbing"
    }
  }
}
```

> 以上为 搜索 about 字段中的 *rock climbing* ，结果匹配到了 2 个文档。从结果看出来，返回的两个文档，其中包含的 _score 字段分数不同，匹配度高的分数高，部分匹配的分数低。由此看出：es 进行全文搜索的时候，首先返回相关性最大的结果，这也是与传统关系型数据库中记录只有匹配和不匹配概念的不同之处。

> **case 6 短语搜索**
> 上例中，匹配的是单独的单词，但有的时候，就是需要匹配到确切的单词序列或者短语（phrases），可以用下面的语法：

```json
GET  /megacorp/employee/_search
{
  "query": {
    "match_phrase": {
      "about": "rock climbing"
    }
  }
}
```

> **case 7 聚合（aggregations）【需求：允许管理者在职员目录中分析】**

```json
GET /megacorp/employee/_search
{
  "size": 0,
  "aggs": {
    "aaaaaaaa": {
      "terms": {
        "field": "age"
      }
    }
  }
}
```

> 以上的 field 如果换成 interests，则会出现错误提示 `Fielddata is disabled on text fields by default. Set fielddata=true on [interests] i`，网上搜到的解决办法是在field后面的字段加上 .keyword，此处为 interests.keyword

```json
GET /megacorp/employee/_search
{   
    ...
    "terms": {
        "field": "interests.keyword"
    }
    ...
}
```

**case 8 聚合允许分级汇总，例如：统计每种兴趣下的职员的平均年龄**

```json
GET /megacorp/employee/_search
{
  "size": 1,
  "aggs": {
    "aaaaaaaa": {
      "terms": {
        "field": "interests.keyword"
      },
      "aggs": {
        "bbbbb": {
          "avg": {
            "field": "age"
          }
        }
      }
    }
  }
}
```

### 1.3 supersdk 订单搜索的一个例子
> 

```json
GET msl-201902/_search
{
  "query": {
		"bool": {
			"must": [{
				"query_string": {
					"query": "1410310230 AND 36028895455670619 AND 15501473451000009236",
					"analyze_wildcard": true
				}
			}, {
				"range": {
					"@timestamp": {
						"gte": 1551163469189,
						"lte": 1551249869189,
						"format": "epoch_millis"
					}
				}
			}],
			"must_not": []
		}
	},
	"size": 1,
  "aggs": {
    "errorcodeaggs": {
      "terms": {
        "field": "order_id.keyword"
      }
    }
  }
}
```

## 2. 集群
## 3. 数据


