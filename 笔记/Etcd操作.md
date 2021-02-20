https://segmentfault.com/a/1190000020868242?utm_source=tag-newest

https://segmentfault.com/a/1190000020787391

# etcd

https://www.infoq.cn/article/etcd-interpretation-application-scenario-implement-principle/

https://etcd.io/docs/current/learning/api/

https://www.liwenzhou.com/posts/Go/go_etcd/

[etcd](https://etcd.io/)是使用Go语言开发的一个开源的、高可用的分布式key-value存储系统，可以用于配置共享和服务的注册和发现。

类似项目有zookeeper和consul。

etcd具有以下特点：

- 完全复制：集群中的每个节点都可以使用完整的存档
- 高可用性：Etcd可用于避免硬件的单点故障或网络问题
- 一致性：每次读取都会返回跨多主机的最新写入
- 简单：包括一个定义良好、面向用户的API（gRPC）
- 安全：实现了带有可选的客户端证书身份验证的自动化TLS
- 快速：每秒10000次写入的基准速度
- 可靠：使用[raft](https://www.cnblogs.com/xybaby/p/10124083.html)算法实现了强一致、高可用的服务存储目录

## etcd应用场景

### 服务发现

服务发现要解决的也是分布式系统中最常见的问题之一，即在同一个分布式集群中的进程或服务，要如何才能找到对方并建立连接。本质上来说，服务发现就是想要了解集群中是否有进程在监听 udp 或 tcp 端口，并且通过名字就可以查找和连接。

![1613550626202](Etcd操作.assets\1613550626202.png)

### 配置中心

将一些配置信息放到 etcd 上进行集中管理。

这类场景的使用方式通常是这样：应用在启动的时候主动从 etcd 获取一次配置信息，同时，在 etcd 节点上注册一个 Watcher 并等待，以后每次配置有更新的时候，etcd 都会实时通知订阅者，以此达到获取最新配置信息的目的。

### 分布式锁

因为 etcd 使用 Raft 算法保持了数据的强一致性，某次操作存储到集群中的值必然是全局一致的，所以很容易实现分布式锁。锁服务有两种使用方式，一是保持独占，二是控制时序。

- **保持独占即所有获取锁的用户最终只有一个可以得到**。etcd 为此提供了一套实现分布式锁原子操作 CAS（`CompareAndSwap`）的 API。通过设置`prevExist`值，可以保证在多个节点同时去创建某个目录时，只有一个成功。而创建成功的用户就可以认为是获得了锁。
- 控制时序，即所有想要获得锁的用户都会被安排执行，但是**获得锁的顺序也是全局唯一的，同时决定了执行顺序**。etcd 为此也提供了一套 API（自动创建有序键），对一个目录建值时指定为`POST`动作，这样 etcd 会自动在目录下生成一个当前最大的值为键，存储这个新的值（客户端编号）。同时还可以使用 API 按顺序列出所有当前目录下的键值。此时这些键的值就是客户端的时序，而这些键中存储的值可以是代表客户端的编号。

![1613550643169](Etcd操作.assets\1613550643169.png)



## etcd架构

![1613631172756](Etcd操作.assets\1613631172756.png)

从 etcd 的架构图中我们可以看到，etcd 主要分为四个部分。

- HTTP Server： 用于处理用户发送的 API 请求以及其它 etcd 节点的同步与心跳信息请求。
- Store：用于处理 etcd 支持的各类功能的事务，包括数据索引、节点状态变更、监控与反馈、事件处理与执行等等，是 etcd 对用户提供的大多数 API 功能的具体实现。
- Raft：Raft 强一致性算法的具体实现，是 etcd 的核心。
- WAL：Write Ahead Log（预写式日志），是 etcd 的数据存储方式。除了在内存中存有所有数据的状态以及节点的索引以外，etcd 就通过 WAL 进行持久化存储。WAL 中，所有的数据提交前都会事先记录日志。Snapshot 是为了防止数据过多而进行的状态快照；Entry 表示存储的具体日志内容。

通常，一个用户的请求发送过来，会经由 HTTP Server 转发给 Store 进行具体的事务处理，如果涉及到节点的修改，则交给 Raft 模块进行状态的变更、日志的记录，然后再同步给别的 etcd 节点以确认数据提交，最后进行数据的提交，再次同步。

## 数据存储

etcd 的存储分为内存存储和持久化（硬盘）存储两部分，内存中的存储除了顺序化的记录下所有用户对节点数据变更的记录外，还会对用户数据进行索引、建堆等方便查询的操作。而持久化则使用预写式日志（WAL：Write Ahead Log）进行记录存储。

在 WAL 的体系中，所有的数据在提交之前都会进行日志记录。在 etcd 的持久化存储目录中，有两个子目录。一个是 WAL，存储着所有事务的变化记录；另一个则是 snapshot，用于存储某一个时刻 etcd 所有目录的数据。通过 WAL 和 snapshot 相结合的方式，etcd 可以有效的进行数据存储和节点故障恢复等操作。

既然有了 WAL 实时存储了所有的变更，为什么还需要 snapshot 呢？随着使用量的增加，WAL 存储的数据会暴增，为了防止磁盘很快就爆满，etcd 默认每 10000 条记录做一次 snapshot，经过 snapshot 以后的 WAL 文件就可以删除。而通过 API 可以查询的历史 etcd 操作默认为 1000 条。

首次启动时，etcd 会把启动的配置信息存储到`data-dir`参数指定的数据目录中。配置信息包括本地节点的 ID、集群 ID 和初始时集群信息。用户需要避免 etcd 从一个过期的数据目录中重新启动，因为使用过期的数据目录启动的节点会与集群中的其他节点产生不一致（如：之前已经记录并同意 Leader 节点存储某个信息，重启后又向 Leader 节点申请这个信息）。所以，为了最大化集群的安全性，一旦有任何数据损坏或丢失的可能性，你就应该把这个节点从集群中移除，然后加入一个不带数据目录的新节点。

## 为什么用 etcd 而不用ZooKeeper？

etcd 实现的这些功能，ZooKeeper都能实现。那么为什么要用 etcd 而非直接使用ZooKeeper呢？

### 为什么不选择ZooKeeper？

1. 部署维护复杂，其使用的`Paxos`强一致性算法复杂难懂。官方只提供了`Java`和`C`两种语言的接口。
2. 使用`Java`编写引入大量的依赖。运维人员维护起来比较麻烦。
3. 最近几年发展缓慢，不如`etcd`和`consul`等后起之秀。

### 为什么选择etcd？

1. 简单。使用 Go 语言编写部署简单；支持HTTP/JSON API,使用简单；使用 Raft 算法保证强一致性让用户易于理解。
2. etcd 默认数据一更新就进行持久化。
3. etcd 支持 SSL 客户端安全认证。

最后，etcd 作为一个年轻的项目，正在高速迭代和开发中，这既是一个优点，也是一个缺点。优点是它的未来具有无限的可能性，缺点是无法得到大项目长时间使用的检验。然而，目前 `CoreOS`、`Kubernetes`和`CloudFoundry`等知名项目均在生产环境中使用了`etcd`，所以总的来说，etcd值得你去尝试。

## etcd配置指标

https://etcd.io/docs/current/op-guide/configuration/

https://www.cnblogs.com/cbkj-xd/p/11934599.html

原文地址：[Configuration flags](https://github.com/etcd-io/etcd/blob/master/Documentation/op-guide/configuration.md)
etcd通过配置文件，多命令行参数和环境变量进行配置，

可重用的配置文件是YAML文件，其名称和值由一个或多个下面描述的命令行标志组成。为了使用此文件，请将文件路径指定为`--config-file`标志或`ETCD_CONFIG_FILE`环境变量的值。如果需要的话[配置文件示例](https://github.com/etcd-io/etcd/blob/master/etcd.conf.yml.sample)可以作为入口点创建新的配置文件。

在命令行上设置的选项优先于环境中的选项。 如果提供了配置文件，则其他命令行标志和环境变量将被忽略。例如，`etcd --config-file etcd.conf.yml.sample --data-dir /tmp`将会忽略`--data-dir`参数。

参数`--my-flag`的环境变量的格式为`ETCD_MY_FLAG`.它适用于所有参数。

客户端请求[官方的etcd端口](http://www.iana.org/assignments/service-names-port-numbers/service-names-port-numbers.txt)为2379,2380是节点通信端口。可以将etcd端口设置为接受TLS流量，非TLS流量，或同时接受TLS和非TLS流量。

要在Linux启动时使用自定义设置自动启动etcd，强烈建议使用[systemd](https://www.cnblogs.com/cbkj-xd/p/freedesktop.org/wiki/Software/systemd/)单元。

### 成员标记

------

**--name**

- 人类可读的该成员的名字
- 默认值："default"
- 环境变量：ETCD_NAME
- 该值被该节点吃的`--initial-cluster`参数引用(例如 `default=http://localhost:2380`).如果使用[静态引导程序](https://www.cnblogs.com/cbkj-xd/p/11934599.html)，则需要与标志中使用的键匹配。当使用发现服务时，每一个成员需要有唯一的名字。`Hostname`或者`machine-id`是好的选择。

**--data-dir**

- 数据目录的路径
- 默认值："${name}.etcd"  ，当前目录下的xxx.etcd目录，可以指定绝对路劲
- 环境变量：ETCD_DATA_DIR

**--wal-dir**

- 专用的wal目录的路径。如果这个参数被设置，etcd将会写WAL文件到walDir而不是dataDir，允许使用专用磁盘，并有助于避免日志记录和其他IO操作之间的io竞争。
- 默认值：""
- 环境变量：ETCD_WAL_DIR

**--snapshot-count**

- 触发一个快照到磁盘的已提交交易的数量
- 默认值："100000"
- 环境变量：ETCD_SNAPSHOP_COUNT

**--heartbeat-interval**

- 心跳间隔(毫秒为单位)
- 默认值:"100"
- 环境变量：ETCD_HEARTBEAT_INTERVAL

**--election-timeout**

- 选举超时时间(毫秒为单位)，从[文档/tuning.md](https://github.com/etcd-io/etcd/blob/master/Documentation/tuning.md)发现更多细节
- 默认值："1000"
- 环境变量：ETCD_ELECTION_TIMEOUT

**--listen-peer-urls**

- 监听在对等节点流量上的URL列表，该参数告诉etcd在指定的协议://IP:port组合上接受来自其对等方的传入请求。协议可以是http或者https。或者，使用`unix://`或者`unixs://`到unix sockets。如果将0.0.0.0作为IP，etcd将监听在所有的接口上的给定端口。如果给定了Ip和端口，etcd将监听指定的接口和端口。可以使用多个URL指定要监听的地址和端口的数量。 etcd将响应来自任何列出的地址和端口的请求。
- 默认值："[http://localhost:2380](http://localhost:2380/)"
- 环境变量:ETCD_LISTEN_PEER_URLS
- 示例："[http://10.0.0.1:2380](http://10.0.0.1:2380/)"
- 无效的示例："[http://example.com:2380](http://example.com:2380/)"(绑定的域名是无效的)

**--listen-client-urls**

- 监听在客户端流量上的URL列表，该参数告诉etcd在指定的协议://IP:port组合上接受来自客户端的传入请求。协议可以是http或者https。或者，使用`unix://`或者`unixs://`到unix sockets。如果将0.0.0.0作为IP，etcd将监听在所有的接口上的给定端口。如果给定了Ip和端口，etcd将监听指定的接口和端口。可以使用多个URL指定要监听的地址和端口的数量。 etcd将响应来自任何列出的地址和端口的请求。
- 默认值："[http://localhost:2379](http://localhost:2379/)"
- 环境变量:ETCD_LISTEN_CLIENT_URLS
- 示例："[http://10.0.0.1:2379](http://10.0.0.1:2379/)"
- 无效的示例："[http://example.com:2379](http://example.com:2379/)"(绑定的域名是无效的)

**--max-snapshots**

- 保留的快照文件最大数量（0为无限）
- 默认值：5
- 环境变量：ETCD_MAX_SNAPSHOTS
- Windows用户的默认设置是无限制的，建议手动设置到5（或出于安全性的考虑）。

**--max-wals**

- 保留的wal文件最大数量（0为无限）
- 默认值：5
- 环境变量：ETCD_MAX_WALS
- Windows用户的默认设置是无限制的，建议手动设置到5（或出于安全性的考虑）。

**--cors**

- 以逗号分隔的CORS来源白名单（跨来源资源共享）。
- 默认值：""
- 环境变量：ETCD_CORS

**--quota-backent-bytes**

- 后端大小超过给定配额时引发警报（0默认为低空间配额）。
- 默认值：0
- 环境变量：ETCD_QUOTA_BACKEND_BYTES

**--backend-batch-limit**

- BackendBatchLimit是提交后端事务之前的最大数量的操作。
- 默认值：0
- 环境变量：ETCD_BACKEND_BATCH_LIMIT

**--backend-bbolt-freelist-type**

- etcd后端（bboltdb）使用的自由列表类型（支持数组和映射的类型）。
- 默认值：map
- 环境变量：ETCD_BACKEND_BBOLT_FREELIST_TYPE

**--backend-batch-interval**

- BackendBatchInterval是提交后端事务之前的最长时间。
- 默认值：0
- 环境变量：ETCD_BACKEND_BATCH_INTERVAL

**--max-txn-ops**

- 交易中允许的最大操作数。
- 默认值：128
- 环境变量：ETCD_MAX_TXN_OPS

**--max-request-bytes**

- 服务器将接受的最大客户端请求大小（以字节为单位）。
- 默认值：1572864
- 环境变量：ETCD_MAX_REQUEST_BYTES

**--grpc-keepalive-min-time**

- 客户端在ping服务器之前应等待的最小持续时间间隔。
- 默认值：5s
- 环境变量：ETCD_GRPC_KEEPALIVE_MIN_TIME

**--grpc-keepalive-interval**

- 服务器到客户端ping的频率持续时间，以检查连接是否有效（0禁用）。
- 默认值：2h
- 环境变量：ETCD_GRPC_KEEPALIVE_INTERVAL

**--grpc-keepalive-timeout**

- 关闭无响应的连接之前的额外等待时间（0禁用）。
- 默认值：20s
- 环境变量：ETCD_GRPC_KEEPALIVE_TIMEOUT

### 集群参数

------

`--initial-advertise-peer-urls`,`--initial-cluster`,`--initial-cluster-state`,和`--initial-cluster-token`参数用于启动([静态启动](https://www.cnblogs.com/cbkj-xd/p/11934599.html),[发现服务启动](https://www.cnblogs.com/cbkj-xd/p/11934599.html)或者[运行时重新配置](https://newonexd.github.io/2019/11/23/blog/etcd/运行时重新配置/))一个新成员，当重启已经存在的成员时将忽略。
前缀为`--discovery`的参数在使用[发现服务](https://newonexd.github.io/2019/11/23/blog/etcd/gRPC命名与发现/)时需要被设置。

**--initial-advertise-peer-urls**

- 此成员的对等URL的列表，以通告到集群的其余部分。 这些地址用于在集群周围传送etcd数据。 所有集群成员必须至少有一个路由。 这些URL可以包含域名。
- 默认值："[http://localhost:2380](http://localhost:2380/)"
- 环境变量：ETCD_INITIAL_ADVERTISE_PEER_URLS
- 示例："[http://example.com:2380](http://example.com:2380/), [http://10.0.0.1:2380](http://10.0.0.1:2380/)"

**--initial-cluster**

- 启动集群的初始化配置
- 默认值："default=[http://localhost:2380](http://localhost:2380/)"
- 环境变量：ETCD_INITIAL_CLUSTER
- 关键是所提供的每个节点的`--name`参数的值。 默认值使用`default`作为密钥，因为这是`--name`参数的默认值。

**--initial-cluster-state**

- 初始群集状态（“新”或“现有”）。 对于在初始静态或DNS引导过程中存在的所有成员，将其设置为`new`。 如果此选项设置为`existing`，则etcd将尝试加入现存集群。 如果设置了错误的值，etcd将尝试启动，但会安全地失败。
- 默认值："new:
- 环境变量：ETCD_INITIAL_CLUSTER_STATE

**--initial-cluster-token**

- 引导期间etcd群集的初始集群令牌。
- 默认值："etcd-cluster"
- 环境变量：ETCD_INITIAL_CLUSTER_TOKEN

**--advertise-client-urls**

- 此成员的客户端URL的列表，这些URL广播给集群的其余部分。 这些URL可以包含域名。
- 默认值：[http://localhost:2379](http://localhost:2379/)
- 环境变量：ETCD_ADVERTISE_CLIENT_URLS
- 示例："[http://example.com:2379](http://example.com:2379/), [http://10.0.0.1:2379](http://10.0.0.1:2379/)"
- 如果从集群成员中发布诸如http://localhost:2379之类的URL并使用etcd的代理功能，请小心。这将导致循环，因为代理将向其自身转发请求，直到其资源（内存，文件描述符）最终耗尽为止。

**--discovery**

- 发现URL用于引导启动集群
- 默认值：""
- 环境变量：ETCD_DISCOVERY

**--discovery-srv**

- 用于引导集群的DNS srv域。
- 默认值：""
- 环境变量：ETCD_DISCOVERY_SRV

**--discovery-srv-name**

- 使用DNS引导时查询的DNS srv名称的后缀。
- 默认值：""
- 环境变量：ETCD_DISCOVERY_SRV_NAME

**--discovery-fallback**

- 发现服务失败时的预期行为(“退出”或“代理”)。“代理”仅支持v2 API。
- 默认值： "proxy"
- 环境变量：ETCD_DISCOVERY_FALLBACK

**--discovery-proxy**

- HTTP代理，用于发现服务的流量。
- 默认值：""
- 环境变量：ETCD_DISCOVERY_PROXY

**--strict-reconfig-check**

- 拒绝可能导致quorum丢失的重新配置请求。
- 默认值：true
- 环境变量：ETCD_STRICT_RECONFIG_CHECK

**--auto-compaction-retention**

- mvcc密钥值存储的自动压缩保留时间（小时）。 0表示禁用自动压缩。
- 默认值：0
- 环境变量：ETCD_AUTO_COMPACTION_RETENTION

**--auto-compaction-mode**

- 解释“自动压缩保留”之一：“定期”，“修订”。 基于期限的保留的“定期”，如果未提供时间单位（例如“ 5m”），则默认为小时。 “修订”用于基于修订号的保留。
- 默认值：periodic
- 环境变量：ETCD_AUTO_COMPACTION_MODE

**--enable-v2**

- 接受etcd V2客户端请求
- 默认值：false
- 环境变量：ETCD_ENABLE_V2

### 代理参数

------

--proxy前缀标志将etcd配置为以代理模式运行。 “代理”仅支持v2 API。

**--proxy**

- 代理模式设置(”off","readonly"或者"on")
- 默认值："off"
- 环境变量：ETCD_PROXY

**--proxy-failure-wait**

- 在重新考虑端点请求之前，端点将保持故障状态的时间（以毫秒为单位）。
- 默认值：5000
- 环境变量：ETCD_PROXY_FAILURE_WAIT

**--proxy-refresh-interval**

- 节点刷新间隔的时间（以毫秒为单位）。
- 默认值：30000
- 环境变量：ETCD_PROXY_REFRESH_INTERVAL

**--proxy-dial-timeout**

- 拨号超时的时间（以毫秒为单位），或0以禁用超时
- 默认值：1000
- 环境变量：ETCD_PROXY_DIAL_TIMEOUT

**--proxy-write-timeout**

- 写入超时的时间（以毫秒为单位）或禁用超时的时间为0。
- 默认值：5000
- 环境变量：ETCD_PROXY_WRITE_TIMEOUT

**--proxy-read-timeout**

- 读取超时的时间（以毫秒为单位），或者为0以禁用超时。
- 如果使用Watch，请勿更改此值，因为会使用较长的轮询请求。
- 默认值：0
- 环境变量：ETCD_PROXY_READ_TIMEOUT

### 安全参数

------

安全参数有助于[构建一个安全的etcd集群](https://www.cnblogs.com/cbkj-xd/p/11934599.html)
**--ca-file**
**DEPRECATED**

- 客户端服务器TLS CA文件的路径。 `--ca-file ca.crt`可以替换为`--trusted-ca-file ca.crt --client-cert-auth`，而etcd将执行相同的操作。
- 默认值：""
- 环境变量：ETCD_CA_FILE

**--cert-file**

- 客户端服务器TLS证书文件的路径
- 默认值：""
- 环境变量：ETCD_CERT_FILE

**--key-file**

- 客户端服务器TLS秘钥文件的路径
- 默认值：""
- 环境变量：ETCD_KEY_FILE

**--client-cert-auth**

- 开启客户端证书认证
- 默认值：false
- 环境变量：ETCD_CLIENT_CERT_AUTH
- CN 权限认证不支持gRPC-网关

**--client-crl-file**

- 客户端被撤销的TLS证书文件的路径
- 默认值：""
- 环境变量：ETCD_CLIENT_CERT_ALLOWED_HOSTNAME

**--client-cert-allowed-hostname**

- 允许客户端证书身份验证的TLS名称。
- 默认值：""
- 环境变量：ETCD_CLIENT_CERT_ALLOWED_HOSTNAME

**--trusted-ca-file**

- 客户端服务器受信任的TLS CA证书文件的路径
- 默认值：""
- 环境变量：ETCD_TRUSTED_CA_FILE

**--auto-tls**

- 客户端TLS使用自动生成的证书
- 默认值：false
- 环境变量：ETCD_AUTO_TLS

**--peer-ca-file**
**已淘汰**

- 节点TLS CA文件的路径.`--peer-ca-file`可以替换为`--peer-trusted-ca-file ca.crt --peer-client-cert-auth`，而etcd将执行相同的操作。
- 默认值：”“
- 环境变量：ETCD_PEER_CA_FILE

**--peer-cert-file**

- 对等服务器TLS证书文件的路径。 这是对等节点通信证书，在服务器和客户端都可以使用。
- 默认值：""
- 环境变量：ETCD_PEER_CERT_FILE

**--peer-key-file**

- 对等服务器TLS秘钥文件的路径。 这是对等节点通信秘钥，在服务器和客户端都可以使用。
- 默认值：""
- 环境变量：ETCD_PEER_KEY_FILE

**--peer-client-cert-auth**

- 启动节点客户端证书认证
- 默认值：false
- 环境变量：ETCD_PEER_CLIENT_CERT_AUTH

**--peer-crl-file**

- 节点被撤销的TLS证书文件的路径
- 默认值：""
- 环境变量：ETCD_PEER_CRL_FILE

**--peer-trusted-ca-file**

- 节点受信任的TLS CA证书文件的路径
- 默认值：""
- 环境变量：ETCD_PEER_TRUSTED_CA_FILE

**--peer-auto-tls**

- 节点使用自动生成的证书
- 默认值：false
- 环境变量：ETCD_PEER_AUTO_TLS

**--peer-cert-allowed-cn**

- 允许使用CommonName进行对等身份验证。
- 默认值：""
- 环境变量：ETCD_PEER_CERT_ALLOWED_CN

**--peer-cert-allowed-hostname**

- 允许的TLS证书名称用于对等身份验证。
- 默认值：""
- 环境变量：ETCD_PEER_CERT_ALLOWED_HOSTNAME

**--cipher-suites**

- 以逗号分隔的服务器/客户端和对等方之间受支持的TLS密码套件列表。
- 默认值：""
- 环境变量：ETCD_CIPHER_SUITES

### 日志参数

------

**--logger**

**v3.4可以使用，警告：`--logger=capnslog`在v3.5被抛弃使用**

- 指定“ zap”用于结构化日志记录或“ capnslog”。
- 默认值：capnslog
- 环境变量：ETCD_LOGGER

**--log-outputs**

- 指定“ stdout”或“ stderr”以跳过日志记录，即使在systemd或逗号分隔的输出目标列表下运行时也是如此。
- 默认值：defalut
- 环境变量：ETCD_LOG_OUTPUTS
- `default`在zap logger迁移期间对v3.4使用`stderr`配置

**--log-level**
**v3.4可以使用**

- 配置日志等级，仅支持`debug,info,warn,error,panic,fatal`
- 默认值：info
- 环境变量：ETCD_LOG_LEVEL
- `default`使用`info`.

**--debug**
**警告：在v3.5被抛弃使用**

- 将所有子程序包的默认日志级别降为DEBUG。
- 默认值：false(所有的包使用INFO)
- 环境变量：ETCD_DEBUG

**--log-package-levels**
**警告：在v3.5被抛弃使用**

- 将各个etcd子软件包设置为特定的日志级别。 一个例子是`etcdserver = WARNING，security = DEBUG`
- 默认值：""(所有的包使用INFO)
- 环境变量：ETCD_LOG_PACKAGE_LEVELS

### 风险参数

------

使用不安全标志时请小心，因为它将破坏共识协议提供的保证。 例如，如果群集中的其他成员仍然存在，可能会`panic`。 使用这些标志时，请遵循说明。
**--force-new-cluster**

- 强制创建一个新的单成员群集。 它提交配置更改，以强制删除群集中的所有现有成员并添加自身，但是强烈建议不要这样做。 请查看[灾难恢复文档](https://www.cnblogs.com/cbkj-xd/p/11934599.html)以了解首选的v3恢复过程。
- 默认值：false
- 环境变量：ETCD_FORCE_NEW_CLUSTER

### 杂项参数

------

**--version**

- 打印版本并退出
- 默认值：false

**--config-file**

- 从文件加载服务器配置。 请注意，如果提供了配置文件，则其他命令行标志和环境变量将被忽略。
- 默认值：""
- 示例：[配置文件示例](https://github.com/etcd-io/etcd/blob/master/etcd.conf.yml.sample)
- 环境变量：ETCD_CONFIG_FILE

### 分析参数

------

**--enable-pprof**

- 通过HTTP服务器启用运行时分析数据。地址位于客户端`URL+“/debug/pprof/”`
- 默认值：false
- 环境变量：ETCD_ENABLE_PPROF

**--metrics**

- 设置导出指标的详细程度，specify 'extensive' to include server side grpc histogram metrics.
- 默认值：basic
- 环境变量：ETCD_METRICS

**--listen-metrics-urls**

- 可以响应`/metrics`和`/health`端点的其他URL列表
- 默认值：""
- 环境变量：ETCD_LISTEN_METRICS_URLS

### 权限参数

------

**--auth-token**

- 指定令牌类型和特定于令牌的选项，特别是对于JWT,格式为`type,var1=val1,var2=val2,...`,可能的类型是`simple`或者`jwt`.对于具体的签名方法jwt可能的变量为`sign-method`（可能的值为`'ES256', 'ES384', 'ES512', 'HS256', 'HS384', 'HS512', 'RS256', 'RS384', 'RS512', 'PS256', 'PS384','PS512'`）
- 对于非对称算法（“ RS”，“ PS”，“ ES”），公钥是可选的，因为私钥包含足够的信息来签名和验证令牌。`pub-key`用于指定用于验证jwt的公钥的路径,`priv-key`用于指定用于对jwt进行签名的私钥的路径，`ttl`用于指定jwt令牌的TTL。
- JWT的示例选项：`-auth-token jwt，pub-key=app.rsa.pub，privkey=app.rsasign-method = RS512，ttl = 10m`
- 默认值："simple"
- 环境变量：ETCD_AUTH_TOKEN

**--bcrypt-cost**

- 指定用于哈希认证密码的bcrypt算法的成本/强度。 有效值在4到31之间。
- 默认值：10
- 环境变量：(不支持)

### 实验参数

------

**--experimental-corrupt-check-time**

- 群集损坏检查通过之间的时间间隔
- 默认值：0s
- 环境变量：ETCD_EXPERIMENTAL_CORRUPT_CHECK_TIME

**--experimental-compaction-batch-limit**

- 设置每个压缩批处理中删除的最大修订。
- 默认值：1000
- 环境变量：ETCD_EXPERIMENTAL_COMPACTION_BATCH_LIMIT

**--experimental-peer-skip-client-san-verification**

- 跳过客户端证书中对等连接的SAN字段验证。 这可能是有帮助的，例如 如果群集成员在NAT后面的不同网络中运行。在这种情况下，请确保使用基于私有证书颁发机构的对等证书.`--peer-cert-file, --peer-key-file, --peer-trusted-ca-file`
- 默认值：false
- 环境变量：ETCD_EXPERIMENTAL_PEER_SKIP_CLIENT_SAN_VERIFICATION

## etcdctl

https://etcd.io/docs/current/demo/

https://etcd.io/docs/current/learning/api/

### put

```bash
etcdctl --endpoints=172.21.0.14:2379 put /key1 v1
```

### get

帮助信息

```bash
etcdctl --endpoints=172.21.0.14:2379 get -h
```

具体命令

```shell
# 普通get单个key
etcdctl --endpoints=172.21.0.14:2379 get /key1
# get多个key（范围操作，顾头不顾尾）
etcdctl --endpoints=172.21.0.14:2379 get /key1 /key3    #得到 key1 key2
# 获取某个前缀的所有键值对，通过 --prefix 可以指定前缀,--limit=xxx 返回长度
etcdctl --endpoints=172.21.0.14:2379 get --prefix /key  #得到key1 key2 key3
# 读取键过往版本的值
etcdctl --endpoints=$ENDPOINTS get --prefix --rev=3 foo #  访问第3个版本的key
# 读取大于等于指定键的 byte 值的键
a = 123
b = 456
z = 789
etcdctl -endpoints=$ENDPOINTS get --from-key b   # 返回b z

# 列出/下所有key
etcdctl --endpoints=172.21.0.14:2379 get / --prefix --keys-only
```

### del

返回影响行数

```shell
# 删除单个key
etcdctl --endpoints=172.21.0.14:2379 del /key5
# 删除范围key
etcdctl --endpoints=172.21.0.14:2379 del /key3 /key5
# 匹配删除多个键值对
etcdctl --endpoints=172.21.0.14:2379 del --prefix /key
# 删除大于等于键 b 的 byte 值的键的命令：
etcdctl del --from-key b   # 返回 2 删除了两个键
```

删除并返回k，v键值对

```shell
etcdctl --endpoints=172.21.0.14:2379 del --prev-kv /key2
```

### watch

 watch 监测一个键值的变化，一旦键值发生更新，就会输出最新的值并退出。 

* 消息发布与订阅

  应用在启动的时候主动从 etcd 获取一次配置信息，同时，在 etcd 节点上注册一个 Watcher 并等待，以后每次配置有更新的时候，etcd 都会实时通知订阅者，以此达到获取最新配置信息的目的。 

* 分布式通知与协调

  分布式通知与协调，与消息发布和订阅有些相似。都用到了 etcd 中的 Watcher 机制，通过注册与异步通知机制，实现分布式环境下不同系统之间的通知与协调，从而对数据变更做到实时处理。实现方式通常是这样：不同系统都在 etcd 上对同一个目录进行注册，同时设置 Watcher 观测该目录的变化（如果对子目录的变化也有需要，可以设置递归模式），当某个系统更新了 etcd 的目录，那么设置了 Watcher 的系统就会收到通知，并作出相应处理

* 集群监控与 Leader 竞选

  前面几个场景已经提到 Watcher 机制，当某个节点消失或有变动时，Watcher 会第一时间发现并告知用户。

```shell
# watch单个key
etcdctl --endpoints=82.156.98.236:2379 watch /key2/k2
	### 在另一终端 etcdctl --endpoints=82.156.98.236:2379 put /key2/k2 v22
	### 返回key变化结果
	PUT
    /key2/k2
    v22
    PUT
    /key2/k2
    v222
```

从 foo to foo9 范围内键的命令： 

```shell
$ etcdctl watch foo foo9
# 在另外一个终端: etcdctl put foo bar
PUT
foo
bar
# 在另外一个终端: etcdctl put foo1 bar1
PUT
foo1
bar1
```

以16进制格式在键 foo 上进行观察的命令： 

```shell
$ etcdctl watch foo --hex
# 在另外一个终端: etcdctl put foo bar
PUT
\x66\x6f\x6f          # 键
\x62\x61\x72          # 值
```

观察多个键 foo 和 zoo 的命令： 

```shell
$ etcdctl watch -i
$ watch foo
$ watch zoo
# 在另外一个终端: etcdctl put foo bar
PUT
foo
bar
# 在另外一个终端: etcdctl put zoo val
PUT
zoo
val
```

观察历史改动

```shell
# 从修订版本 2 开始观察键 `foo` 的改动
$ etcdctl watch --rev=2 foo
PUT
foo
bar
PUT
foo
bar_new
```

从上一次历史修改开始观察： 

```shell
# 在键 `foo` 上观察变更并返回被修改的值和上个修订版本的值
$ etcdctl watch --prev-kv foo
# 在另外一个终端: etcdctl put foo bar_latest
PUT
foo         # 键
bar_new     # 在修改前键foo的上一个值
foo         # 键
bar_latest  # 修改后键foo的值
```

### ttl 租约

* 授予租约

应用可以为 etcd 集群里面的键授予租约。当键被附加到租约时，它的存活时间被绑定到租约的存活时间，而租约的存活时间相应的被 time-to-live (TTL)管理。在租约授予时每个租约的最小TTL值由应用指定。租约的实际 TTL 值是不低于最小 TTL，由 etcd 集群选择。一旦租约的 TTL 到期，租约就过期并且所有附带的键都将被删除。

```shell
# 授予租约，TTL为100秒
$ etcdctl lease grant 100
lease 694d71ddacfda227 granted with TTL(10s)

# 附加键 foo 到租约694d71ddacfda227
$ etcdctl put --lease=694d71ddacfda227 foo10 bar
OK
```

建议时间设置久一点，否则来不及操作会出现如下的错误：

* 撤销租约

应用通过租约 id 可以撤销租约。撤销租约将删除所有它附带的 key。
 假设我们完成了下列的操作：

```shell
$ etcdctl lease revoke 694d71ddacfda227
lease 694d71ddacfda227 revoked

$ etcdctl get foo10
```

* 刷新租期

应用程序可以通过刷新其TTL来保持租约活着，因此不会过期。

```shell
$ etcdctl lease keep-alive 694d71ddacfda227
lease 694d71ddacfda227 keepalived with TTL(100)
lease 694d71ddacfda227 keepalived with TTL(100)
...
```

* 查询租期

应用程序可能想要了解租赁信息，以便它们可以续订或检查租赁是否仍然存在或已过期。应用程序也可能想知道特定租约所附的 key。

假设我们完成了以下一系列操作：

```ruby
$ etcdctl lease grant 300
lease 694d71ddacfda22c granted with TTL(300s)

$ etcdctl put --lease=694d71ddacfda22c foo10 bar
OK
```

获取有关租赁信息以及哪些 key 使用了租赁信息：

```dart
$ etcdctl lease timetolive 694d71ddacfda22c
lease 694d71ddacfda22c granted with TTL(300s), remaining(282s)

$ etcdctl lease timetolive --keys 694d71ddacfda22c
lease 694d71ddacfda22c granted with TTL(300s), remaining(220s), attached keys([foo10])
```

## ssl

http://play.etcd.io/install

## etcd单机

```bash
etcd --data-dir=opt.etcd --name local-1 \
	--initial-advertise-peer-urls http://172.21.0.14:2380 --listen-peer-urls http://172.21.0.14:2380 \
	--advertise-client-urls http://172.21.0.14:2379 --listen-client-urls http://172.21.0.14:2379
```

## etcd集群

https://etcd.io/docs/current/op-guide/clustering/

etcd 作为一个高可用键值存储系统，天生就是为集群化而设计的。由于 Raft 算法在做决策时需要多数节点的投票，所以 etcd 一般部署集群推荐奇数个节点，推荐的数量为 3、5 或者 7 个节点构成一个集群。

### 搭建一个3节点集群示例：

在每个etcd节点指定集群成员，为了区分不同的集群最好同时配置一个独一无二的token。

下面是提前定义好的集群信息，其中`n1`、`n2`和`n3`表示3个不同的etcd节点。

```bash
TOKEN=token-01
CLUSTER_STATE=new
CLUSTER=n1=http://10.240.0.17:2380,n2=http://10.240.0.18:2380,n3=http://10.240.0.19:2380
```

在`n1`这台机器上执行以下命令来启动etcd：

```bash
etcd --data-dir=data.etcd --name n1 \
	--initial-advertise-peer-urls http://10.240.0.17:2380 --listen-peer-urls http://10.240.0.17:2380 \
	--advertise-client-urls http://10.240.0.17:2379 --listen-client-urls http://10.240.0.17:2379 \
	--initial-cluster ${CLUSTER} \
	--initial-cluster-state ${CLUSTER_STATE} --initial-cluster-token ${TOKEN}
```

在`n2`这台机器上执行以下命令启动etcd：

```bash
etcd --data-dir=data.etcd --name n2 \
	--initial-advertise-peer-urls http://10.240.0.18:2380 --listen-peer-urls http://10.240.0.18:2380 \
	--advertise-client-urls http://10.240.0.18:2379 --listen-client-urls http://10.240.0.18:2379 \
	--initial-cluster ${CLUSTER} \
	--initial-cluster-state ${CLUSTER_STATE} --initial-cluster-token ${TOKEN}
```

在`n3`这台机器上执行以下命令启动etcd：

```bash
etcd --data-dir=data.etcd --name n3 \
	--initial-advertise-peer-urls http://10.240.0.19:2380 --listen-peer-urls http://10.240.0.19:2380 \
	--advertise-client-urls http://10.240.0.19:2379 --listen-client-urls http://10.240.0.19:2379 \
	--initial-cluster ${CLUSTER} \
	--initial-cluster-state ${CLUSTER_STATE} --initial-cluster-token ${TOKEN}
```

etcd 官网提供了一个可以公网访问的 etcd 存储地址。你可以通过如下命令得到 etcd 服务的目录，并把它作为`-discovery`参数使用。

```bash
curl https://discovery.etcd.io/new?size=3
https://discovery.etcd.io/a81b5818e67a6ea83e9d4daea5ecbc92

# grab this token
TOKEN=token-01
CLUSTER_STATE=new
DISCOVERY=https://discovery.etcd.io/a81b5818e67a6ea83e9d4daea5ecbc92


etcd --data-dir=data.etcd --name n1 \
	--initial-advertise-peer-urls http://10.240.0.17:2380 --listen-peer-urls http://10.240.0.17:2380 \
	--advertise-client-urls http://10.240.0.17:2379 --listen-client-urls http://10.240.0.17:2379 \
	--discovery ${DISCOVERY} \
	--initial-cluster-state ${CLUSTER_STATE} --initial-cluster-token ${TOKEN}


etcd --data-dir=data.etcd --name n2 \
	--initial-advertise-peer-urls http://10.240.0.18:2380 --listen-peer-urls http://10.240.0.18:2380 \
	--advertise-client-urls http://10.240.0.18:2379 --listen-client-urls http://10.240.0.18:2379 \
	--discovery ${DISCOVERY} \
	--initial-cluster-state ${CLUSTER_STATE} --initial-cluster-token ${TOKEN}


etcd --data-dir=data.etcd --name n3 \
	--initial-advertise-peer-urls http://10.240.0.19:2380 --listen-peer-urls http://10.240.0.19:2380 \
	--advertise-client-urls http://10.240.0.19:2379 --listen-client-urls http:/10.240.0.19:2379 \
	--discovery ${DISCOVERY} \
	--initial-cluster-state ${CLUSTER_STATE} --initial-cluster-token ${TOKEN}
```

到此etcd集群就搭建起来了，可以使用`etcdctl`来连接etcd。

```bash
export ETCDCTL_API=3
HOST_1=10.240.0.17
HOST_2=10.240.0.18
HOST_3=10.240.0.19
ENDPOINTS=$HOST_1:2379,$HOST_2:2379,$HOST_3:2379

etcdctl --endpoints=$ENDPOINTS member list
```





# go操作etcd

https://pkg.go.dev/go.etcd.io/etcd/clientv3/concurrency

https://github.com/etcd-io/etcd/tree/master/client/v3

https://www.cnblogs.com/sunlong88/p/11295424.html

## 安装

```bash
go get go.etcd.io/etcd/clientv3
```

报错

```go
undefined: resolver.BuildOption

undefined: resolver.ResolveNowOption
使用其他版本
google.golang.org/grpc v1.26.0 // indirect
```



## put/get

 `put`命令用来设置键值对数据，`get`命令用来根据key获取值。 

```go
package main

import (
	"context"
	"fmt"
	"time"

	"go.etcd.io/etcd/clientv3"
)

// etcd client put/get demo
// use etcd/clientv3

func main() {
	// 创建一个etcd客户端，初始化配置
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"82.156.98.236:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		// handle error!
		fmt.Printf("connect to etcd failed, err:%v\n", err)
		return
	}
	fmt.Println("connect to etcd success")
	defer cli.Close()
	// pub操作
	// 创建context，1秒超时
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	// 写入数据
	_, err = cli.Put(ctx, "/test/k1", "v1")
	// 关闭context
	cancel()
	if err != nil {
		fmt.Printf("put to etcd failed, err:%v\n", err)
		return
	}

	// get操作
	// 创建context，1秒超时
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	// get数据
	resp, err := cli.Get(ctx, "/test/k1")
	// 关闭context
	cancel()
	if err != nil {
		fmt.Printf("get from etcd failed, err:%v\n", err)
		return
	}
	// 读取键值信息
	for _, ev := range resp.Kvs {
		fmt.Printf("%s:%s\n", ev.Key, ev.Value)
	}
}
```

## watch操作

 `watch`用来获取未来更改的通知。 

```go
package main

import (
	"context"
	"fmt"
	"time"

	"go.etcd.io/etcd/clientv3"
)

// watch demo
func main() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"82.156.98.236:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		fmt.Printf("connect to etcd failed, err:%v\n", err)
		return
	}
	fmt.Println("connect to etcd success")
	defer cli.Close()
	// watch 键/test/k1 的变化
	rch := cli.Watch(context.Background(), "/test/k1") // <-chan WatchResponse
	for wresp := range rch {
		for _, ev := range wresp.Events {
			fmt.Printf("Type: %s Key:%s Value:%s\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
		}
	}
}
```

结果：

```go
connect to etcd success
Type: PUT Key:/test/k1 Value:111
Type: PUT Key:/test/k1 Value:222
Type: PUT Key:/test/k1 Value:333
```

## lease租约



```go
package main

import (
	"fmt"
	"time"
)

// etcd lease

import (
	"context"
	"log"

	"go.etcd.io/etcd/clientv3"
)

func main() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"82.156.98.236:2379"},
		DialTimeout: time.Second * 5,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("connect to etcd success.")
	defer cli.Close()

	// 创建一个30秒的租约
	resp, err := cli.Grant(context.TODO(), 30)
	if err != nil {
		log.Fatal(err)
	}
	// 5秒钟之后, /test/k2 这个key就会被移除
	_, err = cli.Put(context.TODO(), "/test/k2", "v2", clientv3.WithLease(resp.ID))
	if err != nil {
		log.Fatal(err)
	}
}
```

使用.count 判断租约过期

```go
package main

import (
	"fmt"
	"time"
)

// etcd lease

import (
	"context"
	"log"

	"go.etcd.io/etcd/clientv3"
)

func main() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"82.156.98.236:2379"},
		DialTimeout: time.Second * 5,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("connect to etcd success.")
	defer cli.Close()

	// 创建一个30秒的租约
	resp, err := cli.Grant(context.TODO(), 30)
	if err != nil {
		log.Fatal(err)
	}
	// 5秒钟之后, /test/k2 这个key就会被移除
	putResp, err := cli.Put(context.TODO(), "/test/k2", "v2", clientv3.WithLease(resp.ID))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("写入成功",putResp.Header.Revision)
	//定时的看一下/test/k2过期了没有
	//var getResp *clientv3.GetResponse
	//var kv clientv3.KV
	for{
		getResp,err := cli.Get(context.TODO(),"/test/k2")
		if err != nil{
			fmt.Println(err)
			return
		}
		if getResp.Count == 0{
			fmt.Println("kv过期了")
			break
		}
		fmt.Println("还没过期：",getResp.Kvs)
		time.Sleep(time.Second*2)
	}
}
```



## 续租 keepAlive

```go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.etcd.io/etcd/clientv3"
)

// etcd keepAlive

func main() {
	// 初始化一个客户端
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"82.156.98.236:2379"},
		DialTimeout: time.Second * 5,
	})
	if err != nil {
		log.Fatal(err)
	}
	// 非连列成功
	fmt.Println("connect to etcd success.")
	defer cli.Close()

	// 创建5秒的租约
	resp, err := cli.Grant(context.TODO(), 10)
	if err != nil {
		log.Fatal(err)
	}
	// 绑定租约创建key
	_, err = cli.Put(context.TODO(), "/test/k1", "v1", clientv3.WithLease(resp.ID))
	if err != nil {
		log.Fatal(err)
	}

	// 按最大租约持续 续租
	ch, kaerr := cli.KeepAlive(context.TODO(), resp.ID)
	if kaerr != nil {
		log.Fatal(kaerr)
	}
	// 查看剩余时间  一直是最大租约时间
	for {
		ka := <-ch
		fmt.Println("ttl:", ka.TTL)
	}
}
```

## 基于etcd实现分布式锁

`go.etcd.io/etcd/clientv3/concurrency`在etcd之上实现并发操作，如分布式锁、屏障和选举。

导入该包：

```go
import "go.etcd.io/etcd/clientv3/concurrency"
```

基于etcd实现的分布式锁示例：

`go.etcd.io/etcd/clientv3/concurrency`在etcd之上实现并发操作，如分布式锁、屏障和选举。

导入该包：

```go
import "github.com/coreos/etcd/clientv3/concurrency"
```

基于etcd实现的分布式锁示例：

输出：

```bash
acquired lock for s1
released lock for s1
acquired lock for s2
```



```go
import (
    "github.com/coreos/etcd/clientv3"
    "github.com/coreos/etcd/clientv3/concurrency"
)
```



```go
package main

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/concurrency"
	"log"
	"time"
)

func main() {

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"82.156.98.236:2379"},
		DialTimeout: time.Second * 5,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	// 创建两个单独的会话用来演示锁竞争
	s1, err := concurrency.NewSession(cli)
	if err != nil {
		log.Fatal(err)
	}
	defer s1.Close()
	m1 := concurrency.NewMutex(s1, "/my-lock/")

	s2, err := concurrency.NewSession(cli)
	if err != nil {
		log.Fatal(err)
	}
	defer s2.Close()
	m2 := concurrency.NewMutex(s2, "/my-lock/")

	// 会话s1获取锁
	if err := m1.Lock(context.TODO()); err != nil {
		log.Fatal(err)
	}
	fmt.Println("acquired lock for s1")

	m2Locked := make(chan struct{})
	go func() {
		defer close(m2Locked)
		// 等待直到会话s1释放了/my-lock/的锁
		if err := m2.Lock(context.TODO()); err != nil {
			log.Fatal(err)
		}
	}()

	if err := m1.Unlock(context.TODO()); err != nil {
		log.Fatal(err)
	}
	fmt.Println("released lock for s1")

	<-m2Locked
	fmt.Println("acquired lock for s2")
}
```



