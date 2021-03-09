# Elasticsearch

### 介绍

Elasticsearch（ES）是一个基于Lucene构建的开源、分布式、RESTful接口的全文搜索引擎。Elasticsearch还是一个分布式文档数据库，其中每个字段均可被索引，而且每个字段的数据均可被搜索，ES能够横向扩展至数以百计的服务器存储以及处理PB级的数据。可以在极短的时间内存储、搜索和分析大量的数据。通常作为具有复杂搜索场景情况下的核心发动机。

### Elasticsearch能做什么

1. 当你经营一家网上商店，你可以让你的客户搜索你卖的商品。在这种情况下，你可以使用ElasticSearch来存储你的整个产品目录和库存信息，为客户提供精准搜索，可以为客户推荐相关商品。
2. 当你想收集日志或者交易数据的时候，需要分析和挖掘这些数据，寻找趋势，进行统计，总结，或发现异常。在这种情况下，你可以使用Logstash或者其他工具来进行收集数据，当这引起数据存储到ElasticsSearch中。你可以搜索和汇总这些数据，找到任何你感兴趣的信息。
3. 对于程序员来说，比较有名的案例是GitHub，GitHub的搜索是基于ElasticSearch构建的，在github.com/search页面，你可以搜索项目、用户、issue、pull request，还有代码。共有40~50个索引库，分别用于索引网站需要跟踪的各种数据。虽然只索引项目的主分支（master），但这个数据量依然巨大，包括20亿个索引文档，30TB的索引文件。

### Elasticsearch基本概念

#### Near Realtime(NRT) 几乎实时

Elasticsearch是一个几乎实时的搜索平台。意思是，从索引一个文档到这个文档可被搜索只需要一点点的延迟，这个时间一般为毫秒级。

#### Cluster 集群

群集是一个或多个节点（服务器）的集合， 这些节点共同保存整个数据，并在所有节点上提供联合索引和搜索功能。一个集群由一个唯一集群ID确定，并指定一个集群名（默认为“elasticsearch”）。该集群名非常重要，因为节点可以通过这个集群名加入群集，一个节点只能是群集的一部分。

确保在不同的环境中不要使用相同的群集名称，否则可能会导致连接错误的群集节点。例如，你可以使用logging-dev、logging-stage、logging-prod分别为开发、阶段产品、生产集群做记录。

#### Node节点

节点是单个服务器实例，它是群集的一部分，可以存储数据，并参与群集的索引和搜索功能。就像一个集群，节点的名称默认为一个随机的通用唯一标识符（UUID），确定在启动时分配给该节点。如果不希望默认，可以定义任何节点名。这个名字对管理很重要，目的是要确定你的网络服务器对应于你的ElasticSearch群集节点。

我们可以通过群集名配置节点以连接特定的群集。默认情况下，每个节点设置加入名为“elasticSearch”的集群。这意味着如果你启动多个节点在网络上，假设他们能发现彼此都会自动形成和加入一个名为“elasticsearch”的集群。

在单个群集中，你可以拥有尽可能多的节点。此外，如果“elasticsearch”在同一个网络中，没有其他节点正在运行，从单个节点的默认情况下会形成一个新的单节点名为”elasticsearch”的集群。

#### Index索引

索引是具有相似特性的文档集合。例如，可以为客户数据提供索引，为产品目录建立另一个索引，以及为订单数据建立另一个索引。索引由名称（必须全部为小写）标识，该名称用于在对其中的文档执行索引、搜索、更新和删除操作时引用索引。在单个群集中，你可以定义尽可能多的索引。

#### Type类型

在索引中，可以定义一个或多个类型。类型是索引的逻辑类别/分区，其语义完全取决于你。一般来说，类型定义为具有公共字段集的文档。例如，假设你运行一个博客平台，并将所有数据存储在一个索引中。在这个索引中，你可以为用户数据定义一种类型，为博客数据定义另一种类型，以及为注释数据定义另一类型。

#### Document文档

文档是可以被索引的信息的基本单位。例如，你可以为单个客户提供一个文档，单个产品提供另一个文档，以及单个订单提供另一个文档。本文件的表示形式为JSON（JavaScript Object Notation）格式，这是一种非常普遍的互联网数据交换格式。

在索引/类型中，你可以存储尽可能多的文档。请注意，尽管文档物理驻留在索引中，文档实际上必须索引或分配到索引中的类型。

#### Shards & Replicas分片与副本

索引可以存储大量的数据，这些数据可能超过单个节点的硬件限制。例如，十亿个文件占用磁盘空间1TB的单指标可能不适合对单个节点的磁盘或可能太慢服务仅从单个节点的搜索请求。

为了解决这一问题，Elasticsearch提供细分你的指标分成多个块称为分片的能力。当你创建一个索引，你可以简单地定义你想要的分片数量。每个分片本身是一个全功能的、独立的“指数”，可以托管在集群中的任何节点。

**Shards分片的重要性主要体现在以下两个特征：**

1. 分片允许你水平拆分或缩放内容的大小
2. 分片允许你分配和并行操作的碎片（可能在多个节点上）从而提高性能/吞吐量 这个机制中的碎片是分布式的以及其文件汇总到搜索请求是完全由ElasticSearch管理，对用户来说是透明的。

在同一个集群网络或云环境上，故障是任何时候都会出现的，拥有一个故障转移机制以防分片和节点因为某些原因离线或消失是非常有用的，并且被强烈推荐。为此，Elasticsearch允许你创建一个或多个拷贝，你的索引分片进入所谓的副本或称作复制品的分片，简称Replicas。

**Replicas的重要性主要体现在以下两个特征：**

1. 副本为分片或节点失败提供了高可用性。为此，需要注意的是，一个副本的分片不会分配在同一个节点作为原始的或主分片，副本是从主分片那里复制过来的。
2. 副本允许用户扩展你的搜索量或吞吐量，因为搜索可以在所有副本上并行执行。

#### ES基本概念与关系型数据库的比较

|                     ES概念                     |    关系型数据库    |
| :--------------------------------------------: | :----------------: |
|           Index（索引）支持全文检索            | Database（数据库） |
|                  Type（类型）                  |    Table（表）     |
| Document（文档），不同文档可以有不同的字段集合 |   Row（数据行）    |
|                 Field（字段）                  |  Column（数据列）  |
|                Mapping（映射）                 |   Schema（模式）   |



# go操作es

也可以用net/http 包做put get post 发送json操作

https://github.com/olivere/elastic

[接口文档](https://olivere.github.io/elastic/)

https://pkg.go.dev/github.com/olivere/elastic/v7?readme=expanded#section-readme

**版本对应**

| Elasticsearch version | Elastic version | Package URL                                                  | Remarks                               |
| --------------------- | --------------- | ------------------------------------------------------------ | ------------------------------------- |
| 7.x                   | 7.0             | [`github.com/olivere/elastic/v7`](https://github.com/olivere/elastic) ([source](https://github.com/olivere/elastic/tree/release-branch.v7) [doc](http://godoc.org/github.com/olivere/elastic)) | Use Go modules.                       |
| 6.x                   | 6.0             | [`github.com/olivere/elastic`](https://github.com/olivere/elastic) ([source](https://github.com/olivere/elastic/tree/release-branch.v6) [doc](http://godoc.org/github.com/olivere/elastic)) | Use a dependency manager (see below). |
| 5.x                   | 5.0             | [`gopkg.in/olivere/elastic.v5`](https://gopkg.in/olivere/elastic.v5) ([source](https://github.com/olivere/elastic/tree/release-branch.v5) [doc](http://godoc.org/gopkg.in/olivere/elastic.v5)) | Actively maintained.                  |
| 2.x                   | 3.0             | [`gopkg.in/olivere/elastic.v3`](https://gopkg.in/olivere/elastic.v3) ([source](https://github.com/olivere/elastic/tree/release-branch.v3) [doc](http://godoc.org/gopkg.in/olivere/elastic.v3)) | Deprecated. Please update.            |
| 1.x                   | 2.0             | [`gopkg.in/olivere/elastic.v2`](https://gopkg.in/olivere/elastic.v2) ([source](https://github.com/olivere/elastic/tree/release-branch.v2) [doc](http://godoc.org/gopkg.in/olivere/elastic.v2)) | Deprecated. Please update.            |
| 0.9-1.3               | 1.0             | [`gopkg.in/olivere/elastic.v1`](https://gopkg.in/olivere/elastic.v1) ([source](https://github.com/olivere/elastic/tree/release-branch.v1) [doc](http://godoc.org/gopkg.in/olivere/elastic.v1)) | Deprecated. Please update.            |

****

**daemon**

```go
package main

import (
	"context"
	"fmt"

	"github.com/olivere/elastic/v7"
)

// Elasticsearch demo

type Person struct {
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Married bool   `json:"married"`
}

func main() {
	// 初始化es客户端,   elastic.SetSniff(false)忽略内弯该地址链接失败  https://blog.csdn.net/m1126m/article/details/108751132
	//client, err := elastic.NewClient(elastic.SetURL("http://82.156.98.236:9200"))
	client, err := elastic.NewClient(elastic.SetURL("http://82.156.98.236:9200"),elastic.SetSniff(false))
	if err != nil {
		// Handle error
		panic(err)
	}

	fmt.Println("connect to es success")
	// 创建user1 index，并推送数据p1
	p1 := Person{Name: "rion", Age: 22, Married: false}
	put1, err := client.Index().
		Index("user1").
		BodyJson(p1).
		Do(context.Background())
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("Indexed user %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type)
}
```

****

## 初始化客户端

```go
client, err := elastic.NewClient(elastic.SetURL("http://82.156.98.236:9200"))

// 链接失败处理
elastic.NewClient(elastic.SetURL("http://82.156.98.236:9200"),elastic.SetSniff(false))
```



遇到错误

```go
es client connect failed:no active connection found: no Elasticsearch node available
```

解决方法(两种)

[参考](https://blog.csdn.net/m1126m/article/details/108751132)

1. 改变golang代码初始化client时的参数.  client, err := elastic.NewClient(elastic.SetSniff(false),elastic.SetURL(host…))  新增参数 elastic.SetSniff(false), 用于关闭 Sniff

2. 调整es node 的 publish_address 配置，新增 network.publish_host: 127.0.0.1, 整体配置文件如下:

   ```shell
   [root@e33229400bf1 config]# cat elasticsearch.yml
   	cluster.name: "docker-cluster"
   	network.host: 0.0.0.0
   	network.publish_host: 127.0.0.1 # 新增配置项
   ```

   

## 判断索引是否存在

**IndexExists**

```go
func (c *Client) IndexExists(indices ...string) *IndicesExistsService
```

其中

	* 可以传入多个参数[]string{"indx1","indx2",...}...或者"idx1","idx2",...

返回值源码**indices_exists.go**：

```go
func (c *Client) IndexExists(indices ...string) *IndicesExistsService {
	return NewIndicesExistsService(c).Index(indices)
}

// IndicesExistsService checks if an index or indices exist or not.
//
// See https://www.elastic.co/guide/en/elasticsearch/reference/7.0/indices-exists.html
// for details.
type IndicesExistsService struct {
	client *Client

	pretty     *bool       // pretty format the returned JSON response
	human      *bool       // return human readable values for statistics
	errorTrace *bool       // include the stack trace of returned errors
	filterPath []string    // list of filters used to reduce the response
	headers    http.Header // custom request-level HTTP headers

	index             []string
	ignoreUnavailable *bool
	allowNoIndices    *bool
	expandWildcards   string
	local             *bool
}

func (s *IndicesExistsService) Do(ctx context.Context) (bool, error){}
```

示例：

```go
package main

import (
	"context"
	"fmt"

	"github.com/olivere/elastic/v7"
)

func main() {
	// 初始化es客户端
	client, err := elastic.NewClient(elastic.SetURL("http://82.156.98.236:9200"),elastic.SetSniff(false))
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Println("connect to es success")

	// 判断user1 索引是否存在
	exists, err := client.IndexExists("user1").Do(context.Background())
	if err != nil {
		// Handle error
	}
	if !exists {
		fmt.Println("Index does not exist yet")
	}
}
```

## 获取版本

1. 通过ping

   ```go
   package main
   
   import (
   	"context"
   	"fmt"
   
   	"github.com/olivere/elastic/v7"
   )
   
   func main() {
   	// 初始化es客户端
   	client, err := elastic.NewClient(elastic.SetURL("http://82.156.98.236:9200"),elastic.SetSniff(false))
   	if err != nil {
   		// Handle error
   		panic(err)
   	}
   	fmt.Println("connect to es success")
   
   	// 判断user1 索引是否存在
   	info, code, err := client.Ping("http://82.156.98.236:9200").Do(context.Background())
   	if err != nil {
   		// Handle error
   		panic(err)
   	}
   	fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)
   }
   ```

   输出：

   ```go
   connect to es success
   Elasticsearch returned with code 200 and version 7.11.1
   ```

2. 通过ElasticsearchVersion

   ```go
   package main
   
   import (
   	"fmt"
   
   	"github.com/olivere/elastic/v7"
   )
   
   func main() {
   	// 初始化es客户端
   	client, err := elastic.NewClient(elastic.SetURL("http://82.156.98.236:9200"),elastic.SetSniff(false))
   	if err != nil {
   		// Handle error
   		panic(err)
   	}
   	fmt.Println("connect to es success")
   
   	// 判断user1 索引是否存在
   	esversion, err := client.ElasticsearchVersion("http://82.156.98.236:9200")
   	if err != nil {
   		// Handle error
   		panic(err)
   	}
   	fmt.Printf("Elasticsearch version %s\n", esversion)
   }
   ```

   输出：

   ```go
   connect to es success
   Elasticsearch version 7.11.1
   ```

## 创建index

### 创建指定maping索引

```go
package main

import (
	"context"
	"fmt"

	"github.com/olivere/elastic/v7"
)

func main() {
	// 初始化es客户端
	client, err := elastic.NewClient(elastic.SetURL("http://82.156.98.236:9200"),elastic.SetSniff(false))
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Println("connect to es success")

	// 判断user2，不存在则创建
	indexName := "user2"
	exists, err := client.IndexExists(indexName).Do(context.Background())
	if err != nil {
		// Handle error
		panic(err)
	}
	if !exists {
		// Create a new index.
		mapping := `
{
	"settings":{
		"number_of_shards":1,
		"number_of_replicas":0
	},
	"mappings":{
		"properties":{
			"user":{
				"type":"keyword"
			},
			"message":{
				"type":"text",
				"store": true,
				"fielddata": true
			},
			"retweets":{
				"type":"long"
			},
			"tags":{
				"type":"keyword"
			},
			"location":{
				"type":"geo_point"
			},
			"suggest_field":{
				"type":"completion"
			}
		}
	}
}
`
		createIndex, err := client.CreateIndex(indexName).Body(mapping).Do(context.Background())
		if err != nil {
			// Handle error
			panic(err)
		}
		if !createIndex.Acknowledged {
			// Not acknowledged
		}
	}
}
```

确认mapping

```shell
[root@VM-0-14-centos soft]# curl http://82.156.98.236:9200/user2/_mapping?pretty
{
  "user2" : {
    "mappings" : {
      "properties" : {
        "location" : {
          "type" : "geo_point"
        },
        "message" : {
          "type" : "text",
          "store" : true,
          "fielddata" : true
        },
        "retweets" : {
          "type" : "long"
        },
        "suggest_field" : {
          "type" : "completion",
          "analyzer" : "simple",
          "preserve_separators" : true,
          "preserve_position_increments" : true,
          "max_input_length" : 50
        },
        "tags" : {
          "type" : "keyword"
        },
        "user" : {
          "type" : "keyword"
        }
      }
    }
  }
}
```

### 自动创建mapping

```go
package main

import (
	"context"
	"fmt"

	"github.com/olivere/elastic/v7"
)

func main() {
	// 初始化es客户端
	client, err := elastic.NewClient(elastic.SetURL("http://82.156.98.236:9200"),elastic.SetSniff(false))
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Println("connect to es success")

	type Person struct {
		User    string `json:"user"`
		Message     string    `json:"message"`
		Retweets int   `json:"retweets"`
	}

	// 创建 user index,并推送数据
	indexName := "user3"
	tweet1 := Person{User: "olivere", Message: "Take Five", Retweets: 0}
	put1, err := client.Index().
		Index(indexName).
		Id("1"). // 注释id会随机生成，否则插入只更新id为1的数据
		BodyJson(tweet1).
		Do(context.Background())
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("Indexed tweet %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type)
}
```

验证结果：

```shell
[root@VM-0-14-centos soft]# curl -XGET http://82.156.98.236:9200/user3/_search?pretty -H 'Content-Type:application/json' -d '
{
    "query":{
        "match_all":{}
    }
}'


{
  "took" : 2,
  "timed_out" : false,
  "_shards" : {
    "total" : 1,
    "successful" : 1,
    "skipped" : 0,
    "failed" : 0
  },
  "hits" : {
    "total" : {
      "value" : 1,
      "relation" : "eq"
    },
    "max_score" : 1.0,
    "hits" : [
      {
        "_index" : "user3",
        "_type" : "_doc",
        "_id" : "1",
        "_score" : 1.0,
        "_source" : {
          "user" : "olivere",
          "message" : "Take Five",
          "retweets" : 0
        }
      }
    ]
  }
}
```



## 判断指定id是否有数据

```go
package main

import (
	"context"
	"fmt"

	"github.com/olivere/elastic/v7"
)

func main() {
	// 初始化es客户端
	client, err := elastic.NewClient(elastic.SetURL("http://82.156.45.46:9200"),elastic.SetSniff(false))
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Println("connect to es success")

	type Person struct {
		User    string `json:"user"`
		Message     string    `json:"message"`
		Retweets int   `json:"retweets"`
	}

	// 获取user3 index下指定id数据是否存在
	get1, err := client.Get().
		Index("user3").
		Id("1").
		Do(context.Background())
	if err != nil {
		switch {
		case elastic.IsNotFound(err):
			panic(fmt.Sprintf("Document not found: %v", err))
		case elastic.IsTimeout(err):
			panic(fmt.Sprintf("Timeout retrieving document: %v", err))
		case elastic.IsConnErr(err):
			panic(fmt.Sprintf("Connection problem: %v", err))
		default:
			// Some other kind of error
			panic(err)
		}
	}
	fmt.Printf("Got document %s in version %d from index %s, type %s\n", get1.Id, get1.Version, get1.Index, get1.Type)
}
```

输出：

```go
connect to es success
Got document 1 in version 294711632 from index user3, type _doc
```

## 刷新确保文档可搜索

```go
// 初始化es客户端
client, err := elastic.NewClient(elastic.SetURL("http://82.156.45.46:9200"),elastic.SetSniff(false))
if err != nil {
    // Handle error
    panic(err)
}
fmt.Println("connect to es success")

type Person struct {
    User    string `json:"user"`
    Message     string    `json:"message"`
    Retweets int   `json:"retweets"`
}

// 刷新以确保文档可搜索
_, err = client.Refresh().Index("user3").Do(context.Background())
if err != nil {
    panic(err)
}
```

## 查询数据

```go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic/v7"
	"reflect"
)

func main() {
	// 初始化es客户端
	client, err := elastic.NewClient(elastic.SetURL("http://82.156.45.46:9200"),elastic.SetSniff(false))
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Println("connect to es success")

	type Person struct {
		User    string `json:"user"`
		Message     string    `json:"message"`
		Retweets int   `json:"retweets"`
	}

	// Search with a term query,term精确匹配费keyword类型字段会报错
	termQuery := elastic.NewTermQuery("user", "olivere")
	searchResult, err := client.Search().
		Index("user2").          // search in index "twitter"
		Query(termQuery).        // specify the query
		Sort("user", true).      // sort by "user" field, ascending
		From(0).Size(10).        // take documents 0-9
		Pretty(true).            // pretty print request and response JSON
		Do(context.Background()) // execute
	if err != nil {
		// Handle error
		panic(err)
	}
	// searchResult is of type SearchResult and returns hits, suggestions,
	// and all kinds of other information from Elasticsearch.
	fmt.Printf("Query took %d milliseconds\n", searchResult.TookInMillis)

	// Each is a convenience function that iterates over hits in a search result.
	// It makes sure you don't need to check for nil values in the response.
	// However, it ignores errors in serialization. If you want full control
	// over iterating the hits, see below.
	// 遍历查询到的数据
	var ttyp Person
	for _, item := range searchResult.Each(reflect.TypeOf(ttyp)) {
		t := item.(Person)
		fmt.Printf("Person by %s: %s\n", t.User, t.Message)
	}
	// TotalHits is another convenience function that works even when something goes wrong.
	fmt.Printf("Found a total of %d user3\n", searchResult.TotalHits())

	// Here's how you iterate through results with full control over each step.
	if searchResult.TotalHits() > 0 {
		fmt.Printf("Found a total of %d user3\n", searchResult.TotalHits())

		// Iterate through results
		for _, hit := range searchResult.Hits.Hits {
			// hit.Index contains the name of the index

			// Deserialize hit.Source into a Tweet (could also be just a map[string]interface{}).
			var t Person
			err := json.Unmarshal(hit.Source, &t)
			if err != nil {
				// Deserialization failed
			}

			// Work with tweet
			fmt.Printf("Person by %s: %s\n", t.User, t.Message)
		}
	} else {
		// No hits
		fmt.Print("Found no Person\n")
	}
}
```



// https://studygolang.com/articles/29230?fr=sidebar

```go
////搜索
func query() {
    var res *elastic.SearchResult
    var err error
    //取所有
    res, err = client.Search("megacorp").Type("employee").Do(context.Background())
    printEmployee(res, err)

    //字段相等
    q := elastic.NewQueryStringQuery("last_name:Smith")
    res, err = client.Search("megacorp").Type("employee").Query(q).Do(context.Background())
    if err != nil {
        println(err.Error())
    }
    printEmployee(res, err)



    //条件查询
    //年龄大于30岁的
    boolQ := elastic.NewBoolQuery()
    boolQ.Must(elastic.NewMatchQuery("last_name", "smith"))
    boolQ.Filter(elastic.NewRangeQuery("age").Gt(30))
    res, err = client.Search("megacorp").Type("employee").Query(q).Do(context.Background())
    printEmployee(res, err)

    //短语搜索 搜索about字段中有 rock climbing
    matchPhraseQuery := elastic.NewMatchPhraseQuery("about", "rock climbing")
    res, err = client.Search("megacorp").Type("employee").Query(matchPhraseQuery).Do(context.Background())
    printEmployee(res, err)

    //分析 interests
    aggs := elastic.NewTermsAggregation().Field("interests")
    res, err = client.Search("megacorp").Type("employee").Aggregation("all_interests", aggs).Do(context.Background())
    printEmployee(res, err)

}
//
////简单分页
func list(size,page int) {
    if size < 0 || page < 1 {
        fmt.Printf("param error")
        return
    }
    res,err := client.Search("megacorp").
        Type("employee").
        Size(size).
        From((page-1)*size).
        Do(context.Background())
    printEmployee(res, err)

}
//
//打印查询到的Employee
func printEmployee(res *elastic.SearchResult, err error) {
    if err != nil {
        print(err.Error())
        return
    }
    var typ Employee
    for _, item := range res.Each(reflect.TypeOf(typ)) { //从搜索结果中取数据的方法
        t := item.(Employee)
        fmt.Printf("%#v\n", t)
    }
}

func main() {
    create()
    delete()
    update()
    gets()
    query()
    list(2,1)
}
```



## 更新数据

```go
package main

import (
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
)

func main() {
	// 初始化es客户端
	client, err := elastic.NewClient(elastic.SetURL("http://82.156.45.46:9200"),elastic.SetSniff(false))
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Println("connect to es success")

	// Update a tweet by the update API of Elasticsearch.
	// We just increment the number of retweets.
	//script := elastic.NewScript("ctx._source.retweets += params.num").Param("num", 22)
	script2 := elastic.NewScript("ctx._source.message = params.message").Param("message", "mewmessage")
	update, err := client.Update().Index("user2").Id("1").
		Script(script2).
		Upsert(map[string]interface{}{"retweets": 0}).
		Do(context.Background())
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("New version of tweet id %q is now %d", update.Id, update.Version)
}
```

## 删除索引

```go
package main

import (
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
)

func main() {
	// 初始化es客户端
	client, err := elastic.NewClient(elastic.SetURL("http://82.156.45.46:9200"),elastic.SetSniff(false))
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Println("connect to es success")

	// Delete an index.
	indexName := "user3"
	deleteIndex, err := client.DeleteIndex(indexName).Do(context.Background())
	if err != nil {
		// Handle error
		panic(err)
	}
	if !deleteIndex.Acknowledged {
		// Not acknowledged
	}
}
```




