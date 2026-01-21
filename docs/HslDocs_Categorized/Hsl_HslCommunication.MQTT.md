# HslCommunication - HslCommunication.MQTT

> 分類頁數: 30



---
## HslCommunication.MQTT

[原文連結](http://api.hslcommunication.cn/html/68b3b27c-e8db-3b64-4a05-22be97dc0b2b.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.MQTT](../html/68b3b27c-e8db-3b64-4a05-22be97dc0b2b.htm "HslCommunication.MQTT")

[IMqttSyncConnector 类](../html/50cf47b3-ed1b-4799-b23a-3e0dfb8bbdeb.htm "IMqttSyncConnector 类")

[MqttApplicationMessage 类](../html/120454c9-8979-5ed7-c8c5-f68f2a971dd9.htm "MqttApplicationMessage 类")

[MqttClient 类](../html/8626b3c4-5983-7308-94ba-65e3dd480b58.htm "MqttClient 类")

[MqttClient.MqttMessageReceiveDelegate 委托](../html/dcb32bdc-7192-4e9b-3b4a-2a9411ef1756.htm "MqttClient.MqttMessageReceiveDelegate 委托")

[MqttClient.OnClientConnectedDelegate 委托](../html/89da0b36-b26a-237a-9759-6e321cf281db.htm "MqttClient.OnClientConnectedDelegate 委托")

[MqttClientApplicationMessage 类](../html/83342fdd-249e-ae3e-9be3-7d93fb6f9c1a.htm "MqttClientApplicationMessage 类")

[MqttConnectionOptions 类](../html/ba13f335-1fd8-1492-9d41-9f649c890f04.htm "MqttConnectionOptions 类")

[MqttControlMessage 类](../html/c803871a-cf62-d39d-6fa2-44111b61d29b.htm "MqttControlMessage 类")

[MqttCredential 类](../html/6aa84648-38e8-5c9e-887e-707d8c6a8c30.htm "MqttCredential 类")

[MqttHelper 类](../html/ddd32b50-8b26-e4cf-bb5a-1a05845c7fd5.htm "MqttHelper 类")

[MqttPublishMessage 类](../html/90929e44-65c0-951f-9b09-d45133bcdd3b.htm "MqttPublishMessage 类")

[MqttQualityOfServiceLevel 枚举](../html/09eab11b-6f45-6662-193e-835e88f804c4.htm "MqttQualityOfServiceLevel 枚举")

[MqttRpcApiInfo 类](../html/38f1fe3d-727e-fd82-17f3-9aee65a3019a.htm "MqttRpcApiInfo 类")

[MqttRpcDevice 类](../html/fb20799e-d772-079d-97fb-2185b6c362c3.htm "MqttRpcDevice 类")

[MqttServer 类](../html/41e2aff2-d1ad-d10c-bad0-61b644d686fb.htm "MqttServer 类")

[MqttServer.ClientVerificationDelegate 委托](../html/78e973f6-1322-6aab-fcab-a3c369de62da.htm "MqttServer.ClientVerificationDelegate 委托")

[MqttServer.FileChangedDelegate 委托](../html/dce0b285-6f53-91fe-0298-d1575d5829ce.htm "MqttServer.FileChangedDelegate 委托")

[MqttServer.FileOperateVerificationDelegate 委托](../html/4cc60da7-76db-e712-5476-3eb9b819ef8a.htm "MqttServer.FileOperateVerificationDelegate 委托")

[MqttServer.OnClientApplicationMessageReceiveDelegate 委托](../html/79f3cb71-4050-7330-4bdf-a8814fd11399.htm "MqttServer.OnClientApplicationMessageReceiveDelegate 委托")

[MqttServer.OnClientConnectedDelegate 委托](../html/157668d8-d2d0-3f59-9846-05ae5d15663f.htm "MqttServer.OnClientConnectedDelegate 委托")

[MqttSession 类](../html/55ded54f-b600-a5b3-9019-480cbd6bca00.htm "MqttSession 类")

[MqttSessionInfo 类](../html/8931ec91-35e9-6312-e56c-87be05d9e1f0.htm "MqttSessionInfo 类")

[MqttSubscribeMessage 类](../html/975fb0c4-0b12-5c3e-e0d0-f5786a1d3e9d.htm "MqttSubscribeMessage 类")

[MqttSyncClient 类](../html/ce9fedd4-ea75-c8b5-62fd-6aa23f74eb8e.htm "MqttSyncClient 类")

[MqttSyncClientPool 类](../html/08fe9311-16ed-1fe2-bc9f-ba3beda123fe.htm "MqttSyncClientPool 类")

[SubscribeTopic 类](../html/33ff31d8-87be-dc7c-50b7-ed428567eb0c.htm "SubscribeTopic 类")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslCommunication.MQTT 命名空间 |

[缺少 "N:HslCommunication.MQTT" 的 <summary> 文档]

![](../icons/SectionExpanded.png)类

|  | 类 | 说明 |
| --- | --- | --- |
| 公共类 | [IMqttSyncConnector](50cf47b3-ed1b-4799-b23a-3e0dfb8bbdeb.htm) | 关于MqttSyncClient实现的接口[IConnector](ede7f47f-b60e-e320-8bc2-389883e77083.htm)，从而实现了数据连接池的操作信息 |
| 公共类 | [MqttApplicationMessage](120454c9-8979-5ed7-c8c5-f68f2a971dd9.htm) | Mqtt的一次完整消息内容，包含主题，负载数据，消息等级。  Mqtt's complete message content, including subject, payload data, message level. |
| 公共类代码示例 | [MqttClient](8626b3c4-5983-7308-94ba-65e3dd480b58.htm) | Mqtt协议的客户端实现，支持订阅消息，发布消息，详细的使用例子参考api文档  The client implementation of the Mqtt protocol supports subscription messages and publishing messages. For detailed usage examples, refer to the api documentation. |
| 公共类 | [MqttClientApplicationMessage](83342fdd-249e-ae3e-9be3-7d93fb6f9c1a.htm) | 来自客户端的一次消息的内容，当前类主要是在MQTT的服务端进行使用  The content of a message from the client. The current class is mainly used on the MQTT server |
| 公共类 | [MqttConnectionOptions](ba13f335-1fd8-1492-9d41-9f649c890f04.htm) | 连接MQTT服务器的一些参数信息，适用[MqttClient](8626b3c4-5983-7308-94ba-65e3dd480b58.htm)消息发布订阅客户端以及[MqttSyncClient](ce9fedd4-ea75-c8b5-62fd-6aa23f74eb8e.htm)同步请求客户端。  Some parameter information for connecting to the MQTT server is applicable to the [MqttClient](8626b3c4-5983-7308-94ba-65e3dd480b58.htm) message publishing and subscription client and the [MqttSyncClient](ce9fedd4-ea75-c8b5-62fd-6aa23f74eb8e.htm) synchronization request client. |
| 公共类 | [MqttControlMessage](c803871a-cf62-d39d-6fa2-44111b61d29b.htm) | 定义了Mqtt的相关的控制报文的信息 |
| 公共类 | [MqttCredential](6aa84648-38e8-5c9e-887e-707d8c6a8c30.htm) | Mqtt协议的验证对象，包含用户名和密码  Authentication object of Mqtt protocol, including username and password |
| 公共类 | [MqttHelper](ddd32b50-8b26-e4cf-bb5a-1a05845c7fd5.htm) | Mqtt协议的辅助类，提供了一些协议相关的基础方法，方便客户端和服务器端一起调用。  The auxiliary class of the Mqtt protocol provides some protocol-related basic methods for the client and server to call together. |
| 公共类 | [MqttPublishMessage](90929e44-65c0-951f-9b09-d45133bcdd3b.htm) | Mqtt发送的消息封装对象，是对 [MqttApplicationMessage](120454c9-8979-5ed7-c8c5-f68f2a971dd9.htm) 对象的封装，添加了序号，还有是否重发的信息  The message encapsulation object sent by Mqtt is an encapsulation of the [MqttApplicationMessage](120454c9-8979-5ed7-c8c5-f68f2a971dd9.htm) object, with the serial number added, and whether to retransmit |
| 公共类 | [MqttRpcApiInfo](38f1fe3d-727e-fd82-17f3-9aee65a3019a.htm) | Mqtt的同步网络服务的单Api信息描述类  Single Api information description class of Mqtt's synchronous network service |
| 公共类代码示例 | [MqttRpcDevice](fb20799e-d772-079d-97fb-2185b6c362c3.htm) | 基于MRPC实现的远程设备访问的接口，实现了和基础PLC一样的访问功能，适用的设备为 [MqttServer](41e2aff2-d1ad-d10c-bad0-61b644d686fb.htm) 将PLC实际的通信对象注册为 RPC 接口服务。  The interface for remote device access based on MRPC implements the same access function as the basic PLC. The applicable device is [MqttServer](41e2aff2-d1ad-d10c-bad0-61b644d686fb.htm) to register the actual communication object of the PLC as an RPC interface service. |
| 公共类代码示例 | [MqttServer](41e2aff2-d1ad-d10c-bad0-61b644d686fb.htm) | 一个Mqtt的服务器类对象，本服务器支持发布订阅操作，支持从服务器强制推送数据，支持往指定的客户端推送，支持基于一问一答的远程过程调用（RPC）的数据交互，支持文件上传下载。根据这些功能从而定制化出满足各个场景的服务器，详细的使用说明可以参见代码api文档示例。  An Mqtt server class object. This server supports publish and subscribe operations, supports forced push data from the server, supports push to designated clients, supports data interaction based on one-question-one-answer remote procedure calls (RPC), and supports file upload and download . According to these functions, the server can be customized to meet various scenarios. For detailed instructions, please refer to the code api document example. |
| 公共类 | [MqttSession](55ded54f-b600-a5b3-9019-480cbd6bca00.htm) | Mqtt的会话信息，包含了一些基本的信息内容，客户端的IP地址及端口，Client ID，用户名，活动时间，是否允许发布数据等等  Mqtt's session information includes some basic information content, the client's IP address and port, Client ID, user name, activity time, whether it is allowed to publish data, etc. |
| 公共类 | [MqttSessionInfo](8931ec91-35e9-6312-e56c-87be05d9e1f0.htm) | 用于客户端获取服务器会话状态监视数据的类 |
| 公共类 | [MqttSubscribeMessage](975fb0c4-0b12-5c3e-e0d0-f5786a1d3e9d.htm) | 订阅的消息类，用于客户端向服务器请求订阅的信息  Subscribed message class, used by the client to request subscription information from the server |
| 公共类代码示例 | [MqttSyncClient](ce9fedd4-ea75-c8b5-62fd-6aa23f74eb8e.htm) | 基于MQTT协议的同步访问的客户端程序，支持以同步的方式访问服务器的数据信息，并及时的反馈结果，当服务器启动文件功能时，也支持文件的上传，下载，删除操作等。  The client program based on MQTT protocol for synchronous access supports synchronous access to the server's data information and timely feedback of results, When the server starts the file function, it also supports file upload, download, and delete operations. |
| 公共类 | [MqttSyncClientPool](08fe9311-16ed-1fe2-bc9f-ba3beda123fe.htm) | **[商业授权]** MqttSyncClient客户端的连接池类对象，用于共享当前的连接池，合理的动态调整连接对象，然后进行高效通信的操作，默认连接数无限大。 **[Authorization]** The connection pool class object of the MqttSyncClient is used to share the current connection pool, reasonably dynamically adjust the connection object, and then perform efficient communication operations, The default number of connections is unlimited |
| 公共类 | [SubscribeTopic](33ff31d8-87be-dc7c-50b7-ed428567eb0c.htm) | 订阅的主题信息  Subscribed topic information |

![](../icons/SectionExpanded.png)委托

|  | 委托 | 说明 |
| --- | --- | --- |
| 公共委托 | [MqttClientMqttMessageReceiveDelegate](dcb32bdc-7192-4e9b-3b4a-2a9411ef1756.htm) | 当接收到Mqtt订阅的信息的时候触发  Triggered when receiving Mqtt subscription information |
| 公共委托 | [MqttClientOnClientConnectedDelegate](89da0b36-b26a-237a-9759-6e321cf281db.htm) | 连接服务器成功的委托  Connection server successfully delegated |
| 公共委托 | [MqttServerClientVerificationDelegate](78e973f6-1322-6aab-fcab-a3c369de62da.htm) | 验证的委托 |
| 公共委托 | [MqttServerFileChangedDelegate](dce0b285-6f53-91fe-0298-d1575d5829ce.htm) | 文件变化的委托信息 |
| 公共委托 | [MqttServerFileOperateVerificationDelegate](4cc60da7-76db-e712-5476-3eb9b819ef8a.htm) | 当客户端进行文件操作时，校验客户端合法性的委托，操作码具体查看的常量值  When client performing file operations, verify the legitimacy of the client, and check the constant value of for the operation code. |
| 公共委托 | [MqttServerOnClientApplicationMessageReceiveDelegate](79f3cb71-4050-7330-4bdf-a8814fd11399.htm) | Mqtt的消息收到委托 |
| 公共委托 | [MqttServerOnClientConnectedDelegate](157668d8-d2d0-3f59-9846-05ae5d15663f.htm) | 当前mqtt客户端连接上服务器的事件委托 |

![](../icons/SectionExpanded.png)枚举

|  | 枚举 | 说明 |
| --- | --- | --- |
| 公共枚举 | [MqttQualityOfServiceLevel](09eab11b-6f45-6662-193e-835e88f804c4.htm) | Mqtt消息的质量等级  Mqtt message quality level |

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## IMqttSyncConnector 类

[原文連結](http://api.hslcommunication.cn/html/50cf47b3-ed1b-4799-b23a-3e0dfb8bbdeb.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.MQTT](../html/68b3b27c-e8db-3b64-4a05-22be97dc0b2b.htm "HslCommunication.MQTT")

[IMqttSyncConnector 类](../html/50cf47b3-ed1b-4799-b23a-3e0dfb8bbdeb.htm "IMqttSyncConnector 类")

[IMqttSyncConnector 构造函数](../html/f72816f6-4e43-a231-cd22-77dc3e108807.htm "IMqttSyncConnector 构造函数 ")

[IMqttSyncConnector 属性](../html/48a83527-5004-4f9a-a200-d9e6cc278e96.htm "IMqttSyncConnector 属性")

[IMqttSyncConnector 方法](../html/d32e2f94-c196-13a3-8d22-6e422071d34f.htm "IMqttSyncConnector 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IMqttSyncConnector 类 |

关于MqttSyncClient实现的接口[IConnector](ede7f47f-b60e-e320-8bc2-389883e77083.htm)，从而实现了数据连接池的操作信息

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  HslCommunication.MQTTIMqttSyncConnector

**命名空间：**
 [HslCommunication.MQTT](68b3b27c-e8db-3b64-4a05-22be97dc0b2b.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class IMqttSyncConnector : IConnector
```

```
Public Class IMqttSyncConnector
	Implements IConnector
```

```
public ref class IMqttSyncConnector : IConnector
```

```
type IMqttSyncConnector =  
    class
        interface IConnector
    end
```

IMqttSyncConnector 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [IMqttSyncConnector](4cfbecf4-1ca6-cc48-72d5-bdaa0340dbe3.htm) | 实例化一个默认的对象  Instantiate a default object |
| 公共方法 | [IMqttSyncConnector(MqttConnectionOptions)](8c52b2ef-8f82-bfc4-0756-7c332c74c5bb.htm) | 根据连接的MQTT参数，实例化一个默认的对象  According to the connected MQTT parameters, instantiate a default object |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [GuidToken](3c4e204e-4e01-1340-dcec-6bb49c5d5704.htm) | 唯一的GUID码 |
| 公共属性 | [IsConnectUsing](6cdab0f5-39cc-f749-40f0-8ef4e49c6d3e.htm) | 指示当前的连接是否在使用用 |
| 公共属性 | [LastUseTime](89d7da75-e752-0113-686b-a64e6ca7cbd6.htm) | 最新一次使用的时间 |
| 公共属性 | [SyncClient](98afddbf-9ddd-ad0f-e6fe-60615a0905b4.htm) | MQTT的连接对象 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [Close](38470d5d-9093-1866-545c-b9da54ece336.htm) | 关闭并释放 |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法 | [Open](9cffd7fe-9fd5-7364-aa10-d62c44bcbc84.htm) | 打开连接 |
| 公共方法 | ToString | (继承自 Object。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.MQTT 命名空间](68b3b27c-e8db-3b64-4a05-22be97dc0b2b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## IMqttSyncConnector 构造函数 

[原文連結](http://api.hslcommunication.cn/html/f72816f6-4e43-a231-cd22-77dc3e108807.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.MQTT](../html/68b3b27c-e8db-3b64-4a05-22be97dc0b2b.htm "HslCommunication.MQTT")

[IMqttSyncConnector 类](../html/50cf47b3-ed1b-4799-b23a-3e0dfb8bbdeb.htm "IMqttSyncConnector 类")

[IMqttSyncConnector 构造函数](../html/f72816f6-4e43-a231-cd22-77dc3e108807.htm "IMqttSyncConnector 构造函数 ")

[IMqttSyncConnector 构造函数](../html/4cfbecf4-1ca6-cc48-72d5-bdaa0340dbe3.htm "IMqttSyncConnector 构造函数 ")

[IMqttSyncConnector 构造函数 (MqttConnectionOptions)](../html/8c52b2ef-8f82-bfc4-0756-7c332c74c5bb.htm "IMqttSyncConnector 构造函数 (MqttConnectionOptions)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IMqttSyncConnector 构造函数 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [IMqttSyncConnector](4cfbecf4-1ca6-cc48-72d5-bdaa0340dbe3.htm) | 实例化一个默认的对象  Instantiate a default object |
| 公共方法 | [IMqttSyncConnector(MqttConnectionOptions)](8c52b2ef-8f82-bfc4-0756-7c332c74c5bb.htm) | 根据连接的MQTT参数，实例化一个默认的对象  According to the connected MQTT parameters, instantiate a default object |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[IMqttSyncConnector 类](50cf47b3-ed1b-4799-b23a-3e0dfb8bbdeb.htm)

[HslCommunication.MQTT 命名空间](68b3b27c-e8db-3b64-4a05-22be97dc0b2b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## IMqttSyncConnector 构造函数 

[原文連結](http://api.hslcommunication.cn/html/4cfbecf4-1ca6-cc48-72d5-bdaa0340dbe3.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.MQTT](../html/68b3b27c-e8db-3b64-4a05-22be97dc0b2b.htm "HslCommunication.MQTT")

[IMqttSyncConnector 类](../html/50cf47b3-ed1b-4799-b23a-3e0dfb8bbdeb.htm "IMqttSyncConnector 类")

[IMqttSyncConnector 构造函数](../html/f72816f6-4e43-a231-cd22-77dc3e108807.htm "IMqttSyncConnector 构造函数 ")

[IMqttSyncConnector 构造函数](../html/4cfbecf4-1ca6-cc48-72d5-bdaa0340dbe3.htm "IMqttSyncConnector 构造函数 ")

[IMqttSyncConnector 构造函数 (MqttConnectionOptions)](../html/8c52b2ef-8f82-bfc4-0756-7c332c74c5bb.htm "IMqttSyncConnector 构造函数 (MqttConnectionOptions)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IMqttSyncConnector 构造函数 |

实例化一个默认的对象  
Instantiate a default object

**命名空间：**
 [HslCommunication.MQTT](68b3b27c-e8db-3b64-4a05-22be97dc0b2b.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public IMqttSyncConnector()
```

```
Public Sub New
```

```
public:
IMqttSyncConnector()
```

```
new : unit -> IMqttSyncConnector
```

![](../icons/SectionExpanded.png)参见

#### 引用

[IMqttSyncConnector 类](50cf47b3-ed1b-4799-b23a-3e0dfb8bbdeb.htm)

[IMqttSyncConnector 重载](f72816f6-4e43-a231-cd22-77dc3e108807.htm)

[HslCommunication.MQTT 命名空间](68b3b27c-e8db-3b64-4a05-22be97dc0b2b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## IMqttSyncConnector 构造函数 (MqttConnectionOptions)

[原文連結](http://api.hslcommunication.cn/html/8c52b2ef-8f82-bfc4-0756-7c332c74c5bb.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.MQTT](../html/68b3b27c-e8db-3b64-4a05-22be97dc0b2b.htm "HslCommunication.MQTT")

[IMqttSyncConnector 类](../html/50cf47b3-ed1b-4799-b23a-3e0dfb8bbdeb.htm "IMqttSyncConnector 类")

[IMqttSyncConnector 构造函数](../html/f72816f6-4e43-a231-cd22-77dc3e108807.htm "IMqttSyncConnector 构造函数 ")

[IMqttSyncConnector 构造函数](../html/4cfbecf4-1ca6-cc48-72d5-bdaa0340dbe3.htm "IMqttSyncConnector 构造函数 ")

[IMqttSyncConnector 构造函数 (MqttConnectionOptions)](../html/8c52b2ef-8f82-bfc4-0756-7c332c74c5bb.htm "IMqttSyncConnector 构造函数 (MqttConnectionOptions)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IMqttSyncConnector 构造函数 (MqttConnectionOptions) |

根据连接的MQTT参数，实例化一个默认的对象  
According to the connected MQTT parameters, instantiate a default object

**命名空间：**
 [HslCommunication.MQTT](68b3b27c-e8db-3b64-4a05-22be97dc0b2b.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public IMqttSyncConnector(
	MqttConnectionOptions options
)
```

```
Public Sub New ( 
	options As MqttConnectionOptions
)
```

```
public:
IMqttSyncConnector(
	MqttConnectionOptions^ options
)
```

```
new : 
        options : MqttConnectionOptions -> IMqttSyncConnector
```

#### 参数

options
:   类型：[HslCommunication.MQTTMqttConnectionOptions](ba13f335-1fd8-1492-9d41-9f649c890f04.htm)  
    连接的参数信息

![](../icons/SectionExpanded.png)参见

#### 引用

[IMqttSyncConnector 类](50cf47b3-ed1b-4799-b23a-3e0dfb8bbdeb.htm)

[IMqttSyncConnector 重载](f72816f6-4e43-a231-cd22-77dc3e108807.htm)

[HslCommunication.MQTT 命名空间](68b3b27c-e8db-3b64-4a05-22be97dc0b2b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## IMqttSyncConnector 属性

[原文連結](http://api.hslcommunication.cn/html/48a83527-5004-4f9a-a200-d9e6cc278e96.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.MQTT](../html/68b3b27c-e8db-3b64-4a05-22be97dc0b2b.htm "HslCommunication.MQTT")

[IMqttSyncConnector 类](../html/50cf47b3-ed1b-4799-b23a-3e0dfb8bbdeb.htm "IMqttSyncConnector 类")

[IMqttSyncConnector 属性](../html/48a83527-5004-4f9a-a200-d9e6cc278e96.htm "IMqttSyncConnector 属性")

[GuidToken 属性](../html/3c4e204e-4e01-1340-dcec-6bb49c5d5704.htm "GuidToken 属性 ")

[IsConnectUsing 属性](../html/6cdab0f5-39cc-f749-40f0-8ef4e49c6d3e.htm "IsConnectUsing 属性 ")

[LastUseTime 属性](../html/89d7da75-e752-0113-686b-a64e6ca7cbd6.htm "LastUseTime 属性 ")

[SyncClient 属性](../html/98afddbf-9ddd-ad0f-e6fe-60615a0905b4.htm "SyncClient 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IMqttSyncConnector 属性 |

[IMqttSyncConnector](50cf47b3-ed1b-4799-b23a-3e0dfb8bbdeb.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [GuidToken](3c4e204e-4e01-1340-dcec-6bb49c5d5704.htm) | 唯一的GUID码 |
| 公共属性 | [IsConnectUsing](6cdab0f5-39cc-f749-40f0-8ef4e49c6d3e.htm) | 指示当前的连接是否在使用用 |
| 公共属性 | [LastUseTime](89d7da75-e752-0113-686b-a64e6ca7cbd6.htm) | 最新一次使用的时间 |
| 公共属性 | [SyncClient](98afddbf-9ddd-ad0f-e6fe-60615a0905b4.htm) | MQTT的连接对象 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[IMqttSyncConnector 类](50cf47b3-ed1b-4799-b23a-3e0dfb8bbdeb.htm)

[HslCommunication.MQTT 命名空间](68b3b27c-e8db-3b64-4a05-22be97dc0b2b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GuidToken 属性 

[原文連結](http://api.hslcommunication.cn/html/3c4e204e-4e01-1340-dcec-6bb49c5d5704.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.MQTT](../html/68b3b27c-e8db-3b64-4a05-22be97dc0b2b.htm "HslCommunication.MQTT")

[IMqttSyncConnector 类](../html/50cf47b3-ed1b-4799-b23a-3e0dfb8bbdeb.htm "IMqttSyncConnector 类")

[IMqttSyncConnector 属性](../html/48a83527-5004-4f9a-a200-d9e6cc278e96.htm "IMqttSyncConnector 属性")

[GuidToken 属性](../html/3c4e204e-4e01-1340-dcec-6bb49c5d5704.htm "GuidToken 属性 ")

[IsConnectUsing 属性](../html/6cdab0f5-39cc-f749-40f0-8ef4e49c6d3e.htm "IsConnectUsing 属性 ")

[LastUseTime 属性](../html/89d7da75-e752-0113-686b-a64e6ca7cbd6.htm "LastUseTime 属性 ")

[SyncClient 属性](../html/98afddbf-9ddd-ad0f-e6fe-60615a0905b4.htm "SyncClient 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IMqttSyncConnectorGuidToken 属性 |

唯一的GUID码

**命名空间：**
 [HslCommunication.MQTT](68b3b27c-e8db-3b64-4a05-22be97dc0b2b.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public string GuidToken { get; set; }
```

```
Public Property GuidToken As String
	Get
	Set
```

```
public:
virtual property String^ GuidToken {
	String^ get () sealed;
	void set (String^ value) sealed;
}
```

```
abstract GuidToken : string with get, set
override GuidToken : string with get, set
```

#### 属性值

类型：String

#### 实现

[IConnectorGuidToken](5d2d5ad0-74d8-0b6d-29a0-5cc6b0ead492.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[IMqttSyncConnector 类](50cf47b3-ed1b-4799-b23a-3e0dfb8bbdeb.htm)

[HslCommunication.MQTT 命名空间](68b3b27c-e8db-3b64-4a05-22be97dc0b2b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## IsConnectUsing 属性 

[原文連結](http://api.hslcommunication.cn/html/6cdab0f5-39cc-f749-40f0-8ef4e49c6d3e.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.MQTT](../html/68b3b27c-e8db-3b64-4a05-22be97dc0b2b.htm "HslCommunication.MQTT")

[IMqttSyncConnector 类](../html/50cf47b3-ed1b-4799-b23a-3e0dfb8bbdeb.htm "IMqttSyncConnector 类")

[IMqttSyncConnector 属性](../html/48a83527-5004-4f9a-a200-d9e6cc278e96.htm "IMqttSyncConnector 属性")

[GuidToken 属性](../html/3c4e204e-4e01-1340-dcec-6bb49c5d5704.htm "GuidToken 属性 ")

[IsConnectUsing 属性](../html/6cdab0f5-39cc-f749-40f0-8ef4e49c6d3e.htm "IsConnectUsing 属性 ")

[LastUseTime 属性](../html/89d7da75-e752-0113-686b-a64e6ca7cbd6.htm "LastUseTime 属性 ")

[SyncClient 属性](../html/98afddbf-9ddd-ad0f-e6fe-60615a0905b4.htm "SyncClient 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IMqttSyncConnectorIsConnectUsing 属性 |

指示当前的连接是否在使用用

**命名空间：**
 [HslCommunication.MQTT](68b3b27c-e8db-3b64-4a05-22be97dc0b2b.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public bool IsConnectUsing { get; set; }
```

```
Public Property IsConnectUsing As Boolean
	Get
	Set
```

```
public:
virtual property bool IsConnectUsing {
	bool get () sealed;
	void set (bool value) sealed;
}
```

```
abstract IsConnectUsing : bool with get, set
override IsConnectUsing : bool with get, set
```

#### 属性值

类型：Boolean

#### 实现

[IConnectorIsConnectUsing](60329f21-7e63-6365-78c8-c4aad8efaa85.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[IMqttSyncConnector 类](50cf47b3-ed1b-4799-b23a-3e0dfb8bbdeb.htm)

[HslCommunication.MQTT 命名空间](68b3b27c-e8db-3b64-4a05-22be97dc0b2b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## LastUseTime 属性 

[原文連結](http://api.hslcommunication.cn/html/89d7da75-e752-0113-686b-a64e6ca7cbd6.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.MQTT](../html/68b3b27c-e8db-3b64-4a05-22be97dc0b2b.htm "HslCommunication.MQTT")

[IMqttSyncConnector 类](../html/50cf47b3-ed1b-4799-b23a-3e0dfb8bbdeb.htm "IMqttSyncConnector 类")

[IMqttSyncConnector 属性](../html/48a83527-5004-4f9a-a200-d9e6cc278e96.htm "IMqttSyncConnector 属性")

[GuidToken 属性](../html/3c4e204e-4e01-1340-dcec-6bb49c5d5704.htm "GuidToken 属性 ")

[IsConnectUsing 属性](../html/6cdab0f5-39cc-f749-40f0-8ef4e49c6d3e.htm "IsConnectUsing 属性 ")

[LastUseTime 属性](../html/89d7da75-e752-0113-686b-a64e6ca7cbd6.htm "LastUseTime 属性 ")

[SyncClient 属性](../html/98afddbf-9ddd-ad0f-e6fe-60615a0905b4.htm "SyncClient 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IMqttSyncConnectorLastUseTime 属性 |

最新一次使用的时间

**命名空间：**
 [HslCommunication.MQTT](68b3b27c-e8db-3b64-4a05-22be97dc0b2b.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public DateTime LastUseTime { get; set; }
```

```
Public Property LastUseTime As DateTime
	Get
	Set
```

```
public:
virtual property DateTime LastUseTime {
	DateTime get () sealed;
	void set (DateTime value) sealed;
}
```

```
abstract LastUseTime : DateTime with get, set
override LastUseTime : DateTime with get, set
```

#### 属性值

类型：DateTime

#### 实现

[IConnectorLastUseTime](25b06530-9b51-19db-ba38-a91998928c48.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[IMqttSyncConnector 类](50cf47b3-ed1b-4799-b23a-3e0dfb8bbdeb.htm)

[HslCommunication.MQTT 命名空间](68b3b27c-e8db-3b64-4a05-22be97dc0b2b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## SyncClient 属性 

[原文連結](http://api.hslcommunication.cn/html/98afddbf-9ddd-ad0f-e6fe-60615a0905b4.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.MQTT](../html/68b3b27c-e8db-3b64-4a05-22be97dc0b2b.htm "HslCommunication.MQTT")

[IMqttSyncConnector 类](../html/50cf47b3-ed1b-4799-b23a-3e0dfb8bbdeb.htm "IMqttSyncConnector 类")

[IMqttSyncConnector 属性](../html/48a83527-5004-4f9a-a200-d9e6cc278e96.htm "IMqttSyncConnector 属性")

[GuidToken 属性](../html/3c4e204e-4e01-1340-dcec-6bb49c5d5704.htm "GuidToken 属性 ")

[IsConnectUsing 属性](../html/6cdab0f5-39cc-f749-40f0-8ef4e49c6d3e.htm "IsConnectUsing 属性 ")

[LastUseTime 属性](../html/89d7da75-e752-0113-686b-a64e6ca7cbd6.htm "LastUseTime 属性 ")

[SyncClient 属性](../html/98afddbf-9ddd-ad0f-e6fe-60615a0905b4.htm "SyncClient 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IMqttSyncConnectorSyncClient 属性 |

MQTT的连接对象

**命名空间：**
 [HslCommunication.MQTT](68b3b27c-e8db-3b64-4a05-22be97dc0b2b.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public MqttSyncClient SyncClient { get; set; }
```

```
Public Property SyncClient As MqttSyncClient
	Get
	Set
```

```
public:
property MqttSyncClient^ SyncClient {
	MqttSyncClient^ get ();
	void set (MqttSyncClient^ value);
}
```

```
member SyncClient : MqttSyncClient with get, set
```

#### 属性值

类型：[MqttSyncClient](ce9fedd4-ea75-c8b5-62fd-6aa23f74eb8e.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[IMqttSyncConnector 类](50cf47b3-ed1b-4799-b23a-3e0dfb8bbdeb.htm)

[HslCommunication.MQTT 命名空间](68b3b27c-e8db-3b64-4a05-22be97dc0b2b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## IMqttSyncConnector 方法

[原文連結](http://api.hslcommunication.cn/html/d32e2f94-c196-13a3-8d22-6e422071d34f.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.MQTT](../html/68b3b27c-e8db-3b64-4a05-22be97dc0b2b.htm "HslCommunication.MQTT")

[IMqttSyncConnector 类](../html/50cf47b3-ed1b-4799-b23a-3e0dfb8bbdeb.htm "IMqttSyncConnector 类")

[IMqttSyncConnector 方法](../html/d32e2f94-c196-13a3-8d22-6e422071d34f.htm "IMqttSyncConnector 方法")

[Close 方法](../html/38470d5d-9093-1866-545c-b9da54ece336.htm "Close 方法 ")

[Open 方法](../html/9cffd7fe-9fd5-7364-aa10-d62c44bcbc84.htm "Open 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IMqttSyncConnector 方法 |

[IMqttSyncConnector](50cf47b3-ed1b-4799-b23a-3e0dfb8bbdeb.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [Close](38470d5d-9093-1866-545c-b9da54ece336.htm) | 关闭并释放 |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法 | [Open](9cffd7fe-9fd5-7364-aa10-d62c44bcbc84.htm) | 打开连接 |
| 公共方法 | ToString | (继承自 Object。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[IMqttSyncConnector 类](50cf47b3-ed1b-4799-b23a-3e0dfb8bbdeb.htm)

[HslCommunication.MQTT 命名空间](68b3b27c-e8db-3b64-4a05-22be97dc0b2b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Close 方法 

[原文連結](http://api.hslcommunication.cn/html/38470d5d-9093-1866-545c-b9da54ece336.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.MQTT](../html/68b3b27c-e8db-3b64-4a05-22be97dc0b2b.htm "HslCommunication.MQTT")

[IMqttSyncConnector 类](../html/50cf47b3-ed1b-4799-b23a-3e0dfb8bbdeb.htm "IMqttSyncConnector 类")

[IMqttSyncConnector 方法](../html/d32e2f94-c196-13a3-8d22-6e422071d34f.htm "IMqttSyncConnector 方法")

[Close 方法](../html/38470d5d-9093-1866-545c-b9da54ece336.htm "Close 方法 ")

[Open 方法](../html/9cffd7fe-9fd5-7364-aa10-d62c44bcbc84.htm "Open 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IMqttSyncConnectorClose 方法 |

关闭并释放

**命名空间：**
 [HslCommunication.MQTT](68b3b27c-e8db-3b64-4a05-22be97dc0b2b.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public void Close()
```

```
Public Sub Close
```

```
public:
virtual void Close() sealed
```

```
abstract Close : unit -> unit 
override Close : unit -> unit
```

#### 实现

[IConnectorClose](be1b3bb2-108e-0069-994c-c7a14953cc4e.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[IMqttSyncConnector 类](50cf47b3-ed1b-4799-b23a-3e0dfb8bbdeb.htm)

[HslCommunication.MQTT 命名空间](68b3b27c-e8db-3b64-4a05-22be97dc0b2b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Open 方法 

[原文連結](http://api.hslcommunication.cn/html/9cffd7fe-9fd5-7364-aa10-d62c44bcbc84.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.MQTT](../html/68b3b27c-e8db-3b64-4a05-22be97dc0b2b.htm "HslCommunication.MQTT")

[IMqttSyncConnector 类](../html/50cf47b3-ed1b-4799-b23a-3e0dfb8bbdeb.htm "IMqttSyncConnector 类")

[IMqttSyncConnector 方法](../html/d32e2f94-c196-13a3-8d22-6e422071d34f.htm "IMqttSyncConnector 方法")

[Close 方法](../html/38470d5d-9093-1866-545c-b9da54ece336.htm "Close 方法 ")

[Open 方法](../html/9cffd7fe-9fd5-7364-aa10-d62c44bcbc84.htm "Open 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IMqttSyncConnectorOpen 方法 |

打开连接

**命名空间：**
 [HslCommunication.MQTT](68b3b27c-e8db-3b64-4a05-22be97dc0b2b.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public void Open()
```

```
Public Sub Open
```

```
public:
virtual void Open() sealed
```

```
abstract Open : unit -> unit 
override Open : unit -> unit
```

#### 实现

[IConnectorOpen](08379c6c-8548-7e37-08d4-6a8b71bfd509.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[IMqttSyncConnector 类](50cf47b3-ed1b-4799-b23a-3e0dfb8bbdeb.htm)

[HslCommunication.MQTT 命名空间](68b3b27c-e8db-3b64-4a05-22be97dc0b2b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## MqttApplicationMessage 类

[原文連結](http://api.hslcommunication.cn/html/120454c9-8979-5ed7-c8c5-f68f2a971dd9.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.MQTT](../html/68b3b27c-e8db-3b64-4a05-22be97dc0b2b.htm "HslCommunication.MQTT")

[MqttApplicationMessage 类](../html/120454c9-8979-5ed7-c8c5-f68f2a971dd9.htm "MqttApplicationMessage 类")

[MqttApplicationMessage 构造函数](../html/9eb36127-3e46-6e94-f5c8-2ef12a620f9a.htm "MqttApplicationMessage 构造函数 ")

[MqttApplicationMessage 属性](../html/39a876f7-8e1b-cfd3-78d9-d5bc06f824c0.htm "MqttApplicationMessage 属性")

[MqttApplicationMessage 方法](../html/bc8f7ff9-24ce-98e4-32fd-4aebea6d55d5.htm "MqttApplicationMessage 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MqttApplicationMessage 类 |

Mqtt的一次完整消息内容，包含主题，负载数据，消息等级。  
Mqtt's complete message content, including subject, payload data, message level.

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  HslCommunication.MQTTMqttApplicationMessage  
    [HslCommunication.MQTTMqttClientApplicationMessage](83342fdd-249e-ae3e-9be3-7d93fb6f9c1a.htm)

**命名空间：**
 [HslCommunication.MQTT](68b3b27c-e8db-3b64-4a05-22be97dc0b2b.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class MqttApplicationMessage
```

```
Public Class MqttApplicationMessage
```

```
public ref class MqttApplicationMessage
```

```
type MqttApplicationMessage =  class end
```

MqttApplicationMessage 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [MqttApplicationMessage](9eb36127-3e46-6e94-f5c8-2ef12a620f9a.htm) | 初始化 MqttApplicationMessage 类的一个新实例 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [Payload](f9cdfff9-4ef5-1a47-50a9-c172362551c6.htm) | 有效载荷包含将被发布的应用消息。数据的内容和格式是应用特定的。  The payload contains application messages to be published. The content and format of the data is application specific. |
| 公共属性 | [QualityOfServiceLevel](d4e9ed1b-415b-e146-d416-d117d1b21624.htm) | 这个字段表示应用消息分发的服务质量等级保证。分为，最多一次，最少一次，正好一次，只发不推送。  This field indicates the quality of service level guarantee for application message distribution. Divided into, at most once, at least once, exactly once |
| 公共属性 | [Retain](f8c619d3-a883-5e37-4e7d-499ca4e5a73d.htm) | 该消息是否在服务器端进行保留，详细的说明参照文档的备注  Whether the message is retained on the server. For details, refer to the remarks of the document. |
| 公共属性 | [Topic](3fef75f8-da83-3b25-3617-d6f93f61a44a.htm) | 主题名（Topic Name）用于识别有效载荷数据应该被发布到哪一个信息通道。  The Topic Name is used to identify which information channel the payload data should be published to. |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法 | [ToString](fab1a5b9-f28d-c8fa-5e7d-215aeb1ded62.htm) | (重写 ObjectToString.) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.MQTT 命名空间](68b3b27c-e8db-3b64-4a05-22be97dc0b2b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## MqttApplicationMessage 构造函数 

[原文連結](http://api.hslcommunication.cn/html/9eb36127-3e46-6e94-f5c8-2ef12a620f9a.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.MQTT](../html/68b3b27c-e8db-3b64-4a05-22be97dc0b2b.htm "HslCommunication.MQTT")

[MqttApplicationMessage 类](../html/120454c9-8979-5ed7-c8c5-f68f2a971dd9.htm "MqttApplicationMessage 类")

[MqttApplicationMessage 构造函数](../html/9eb36127-3e46-6e94-f5c8-2ef12a620f9a.htm "MqttApplicationMessage 构造函数 ")

[MqttApplicationMessage 属性](../html/39a876f7-8e1b-cfd3-78d9-d5bc06f824c0.htm "MqttApplicationMessage 属性")

[MqttApplicationMessage 方法](../html/bc8f7ff9-24ce-98e4-32fd-4aebea6d55d5.htm "MqttApplicationMessage 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MqttApplicationMessage 构造函数 |

初始化 [MqttApplicationMessage](120454c9-8979-5ed7-c8c5-f68f2a971dd9.htm) 类的一个新实例

**命名空间：**
 [HslCommunication.MQTT](68b3b27c-e8db-3b64-4a05-22be97dc0b2b.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public MqttApplicationMessage()
```

```
Public Sub New
```

```
public:
MqttApplicationMessage()
```

```
new : unit -> MqttApplicationMessage
```

![](../icons/SectionExpanded.png)参见

#### 引用

[MqttApplicationMessage 类](120454c9-8979-5ed7-c8c5-f68f2a971dd9.htm)

[HslCommunication.MQTT 命名空间](68b3b27c-e8db-3b64-4a05-22be97dc0b2b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## MqttApplicationMessage 属性

[原文連結](http://api.hslcommunication.cn/html/39a876f7-8e1b-cfd3-78d9-d5bc06f824c0.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.MQTT](../html/68b3b27c-e8db-3b64-4a05-22be97dc0b2b.htm "HslCommunication.MQTT")

[MqttApplicationMessage 类](../html/120454c9-8979-5ed7-c8c5-f68f2a971dd9.htm "MqttApplicationMessage 类")

[MqttApplicationMessage 属性](../html/39a876f7-8e1b-cfd3-78d9-d5bc06f824c0.htm "MqttApplicationMessage 属性")

[Payload 属性](../html/f9cdfff9-4ef5-1a47-50a9-c172362551c6.htm "Payload 属性 ")

[QualityOfServiceLevel 属性](../html/d4e9ed1b-415b-e146-d416-d117d1b21624.htm "QualityOfServiceLevel 属性 ")

[Retain 属性](../html/f8c619d3-a883-5e37-4e7d-499ca4e5a73d.htm "Retain 属性 ")

[Topic 属性](../html/3fef75f8-da83-3b25-3617-d6f93f61a44a.htm "Topic 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MqttApplicationMessage 属性 |

[MqttApplicationMessage](120454c9-8979-5ed7-c8c5-f68f2a971dd9.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [Payload](f9cdfff9-4ef5-1a47-50a9-c172362551c6.htm) | 有效载荷包含将被发布的应用消息。数据的内容和格式是应用特定的。  The payload contains application messages to be published. The content and format of the data is application specific. |
| 公共属性 | [QualityOfServiceLevel](d4e9ed1b-415b-e146-d416-d117d1b21624.htm) | 这个字段表示应用消息分发的服务质量等级保证。分为，最多一次，最少一次，正好一次，只发不推送。  This field indicates the quality of service level guarantee for application message distribution. Divided into, at most once, at least once, exactly once |
| 公共属性 | [Retain](f8c619d3-a883-5e37-4e7d-499ca4e5a73d.htm) | 该消息是否在服务器端进行保留，详细的说明参照文档的备注  Whether the message is retained on the server. For details, refer to the remarks of the document. |
| 公共属性 | [Topic](3fef75f8-da83-3b25-3617-d6f93f61a44a.htm) | 主题名（Topic Name）用于识别有效载荷数据应该被发布到哪一个信息通道。  The Topic Name is used to identify which information channel the payload data should be published to. |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[MqttApplicationMessage 类](120454c9-8979-5ed7-c8c5-f68f2a971dd9.htm)

[HslCommunication.MQTT 命名空间](68b3b27c-e8db-3b64-4a05-22be97dc0b2b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Payload 属性 

[原文連結](http://api.hslcommunication.cn/html/f9cdfff9-4ef5-1a47-50a9-c172362551c6.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.MQTT](../html/68b3b27c-e8db-3b64-4a05-22be97dc0b2b.htm "HslCommunication.MQTT")

[MqttApplicationMessage 类](../html/120454c9-8979-5ed7-c8c5-f68f2a971dd9.htm "MqttApplicationMessage 类")

[MqttApplicationMessage 属性](../html/39a876f7-8e1b-cfd3-78d9-d5bc06f824c0.htm "MqttApplicationMessage 属性")

[Payload 属性](../html/f9cdfff9-4ef5-1a47-50a9-c172362551c6.htm "Payload 属性 ")

[QualityOfServiceLevel 属性](../html/d4e9ed1b-415b-e146-d416-d117d1b21624.htm "QualityOfServiceLevel 属性 ")

[Retain 属性](../html/f8c619d3-a883-5e37-4e7d-499ca4e5a73d.htm "Retain 属性 ")

[Topic 属性](../html/3fef75f8-da83-3b25-3617-d6f93f61a44a.htm "Topic 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MqttApplicationMessagePayload 属性 |

有效载荷包含将被发布的应用消息。数据的内容和格式是应用特定的。  
The payload contains application messages to be published. The content and format of the data is application specific.

**命名空间：**
 [HslCommunication.MQTT](68b3b27c-e8db-3b64-4a05-22be97dc0b2b.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public byte[] Payload { get; set; }
```

```
Public Property Payload As Byte()
	Get
	Set
```

```
public:
property array<unsigned char>^ Payload {
	array<unsigned char>^ get ();
	void set (array<unsigned char>^ value);
}
```

```
member Payload : byte[] with get, set
```

#### 属性值

类型：Byte

![](../icons/SectionExpanded.png)参见

#### 引用

[MqttApplicationMessage 类](120454c9-8979-5ed7-c8c5-f68f2a971dd9.htm)

[HslCommunication.MQTT 命名空间](68b3b27c-e8db-3b64-4a05-22be97dc0b2b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## QualityOfServiceLevel 属性 

[原文連結](http://api.hslcommunication.cn/html/d4e9ed1b-415b-e146-d416-d117d1b21624.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.MQTT](../html/68b3b27c-e8db-3b64-4a05-22be97dc0b2b.htm "HslCommunication.MQTT")

[MqttApplicationMessage 类](../html/120454c9-8979-5ed7-c8c5-f68f2a971dd9.htm "MqttApplicationMessage 类")

[MqttApplicationMessage 属性](../html/39a876f7-8e1b-cfd3-78d9-d5bc06f824c0.htm "MqttApplicationMessage 属性")

[Payload 属性](../html/f9cdfff9-4ef5-1a47-50a9-c172362551c6.htm "Payload 属性 ")

[QualityOfServiceLevel 属性](../html/d4e9ed1b-415b-e146-d416-d117d1b21624.htm "QualityOfServiceLevel 属性 ")

[Retain 属性](../html/f8c619d3-a883-5e37-4e7d-499ca4e5a73d.htm "Retain 属性 ")

[Topic 属性](../html/3fef75f8-da83-3b25-3617-d6f93f61a44a.htm "Topic 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MqttApplicationMessageQualityOfServiceLevel 属性 |

这个字段表示应用消息分发的服务质量等级保证。分为，最多一次，最少一次，正好一次，只发不推送。  
This field indicates the quality of service level guarantee for application message distribution. Divided into, at most once, at least once, exactly once

**命名空间：**
 [HslCommunication.MQTT](68b3b27c-e8db-3b64-4a05-22be97dc0b2b.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public MqttQualityOfServiceLevel QualityOfServiceLevel { get; set; }
```

```
Public Property QualityOfServiceLevel As MqttQualityOfServiceLevel
	Get
	Set
```

```
public:
property MqttQualityOfServiceLevel QualityOfServiceLevel {
	MqttQualityOfServiceLevel get ();
	void set (MqttQualityOfServiceLevel value);
}
```

```
member QualityOfServiceLevel : MqttQualityOfServiceLevel with get, set
```

#### 属性值

类型：[MqttQualityOfServiceLevel](09eab11b-6f45-6662-193e-835e88f804c4.htm)

![](../icons/SectionExpanded.png)备注

在实际的开发中的情况下，最多一次是最省性能的，正好一次是最消耗性能的，如果应有场景为推送实时的数据，那么，最多一次的性能是最高的

![](../icons/SectionExpanded.png)参见

#### 引用

[MqttApplicationMessage 类](120454c9-8979-5ed7-c8c5-f68f2a971dd9.htm)

[HslCommunication.MQTT 命名空间](68b3b27c-e8db-3b64-4a05-22be97dc0b2b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Retain 属性 

[原文連結](http://api.hslcommunication.cn/html/f8c619d3-a883-5e37-4e7d-499ca4e5a73d.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.MQTT](../html/68b3b27c-e8db-3b64-4a05-22be97dc0b2b.htm "HslCommunication.MQTT")

[MqttApplicationMessage 类](../html/120454c9-8979-5ed7-c8c5-f68f2a971dd9.htm "MqttApplicationMessage 类")

[MqttApplicationMessage 属性](../html/39a876f7-8e1b-cfd3-78d9-d5bc06f824c0.htm "MqttApplicationMessage 属性")

[Payload 属性](../html/f9cdfff9-4ef5-1a47-50a9-c172362551c6.htm "Payload 属性 ")

[QualityOfServiceLevel 属性](../html/d4e9ed1b-415b-e146-d416-d117d1b21624.htm "QualityOfServiceLevel 属性 ")

[Retain 属性](../html/f8c619d3-a883-5e37-4e7d-499ca4e5a73d.htm "Retain 属性 ")

[Topic 属性](../html/3fef75f8-da83-3b25-3617-d6f93f61a44a.htm "Topic 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MqttApplicationMessageRetain 属性 |

该消息是否在服务器端进行保留，详细的说明参照文档的备注  
Whether the message is retained on the server. For details, refer to the remarks of the document.

**命名空间：**
 [HslCommunication.MQTT](68b3b27c-e8db-3b64-4a05-22be97dc0b2b.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public bool Retain { get; set; }
```

```
Public Property Retain As Boolean
	Get
	Set
```

```
public:
property bool Retain {
	bool get ();
	void set (bool value);
}
```

```
member Retain : bool with get, set
```

#### 属性值

类型：Boolean

![](../icons/SectionExpanded.png)备注

如果客户端发给服务端的PUBLISH报文的保留（RETAIN）标志被设置为1，服务端必须存储这个应用消息和它的服务质量等级（QoS），
以便它可以被分发给未来的主题名匹配的订阅者 [MQTT-3.3.1-5]。一个新的订阅建立时，对每个匹配的主题名
，如果存在最近保留的消息，它必须被发送给这个订阅者 [MQTT-3.3.1-6]。如果服务端收到一条保留（RETAIN）标志为1的QoS 0消息，
它必须丢弃之前为那个主题保留的任何消息。它应该将这个新的QoS 0消息当作那个主题的新保留消息，但是任何时候都可以选择丢弃它 — 如果这种情况发生了，
那个主题将没有保留消息 [MQTT-3.3.1-7]

![](../icons/SectionExpanded.png)参见

#### 引用

[MqttApplicationMessage 类](120454c9-8979-5ed7-c8c5-f68f2a971dd9.htm)

[HslCommunication.MQTT 命名空间](68b3b27c-e8db-3b64-4a05-22be97dc0b2b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Topic 属性 

[原文連結](http://api.hslcommunication.cn/html/3fef75f8-da83-3b25-3617-d6f93f61a44a.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.MQTT](../html/68b3b27c-e8db-3b64-4a05-22be97dc0b2b.htm "HslCommunication.MQTT")

[MqttApplicationMessage 类](../html/120454c9-8979-5ed7-c8c5-f68f2a971dd9.htm "MqttApplicationMessage 类")

[MqttApplicationMessage 属性](../html/39a876f7-8e1b-cfd3-78d9-d5bc06f824c0.htm "MqttApplicationMessage 属性")

[Payload 属性](../html/f9cdfff9-4ef5-1a47-50a9-c172362551c6.htm "Payload 属性 ")

[QualityOfServiceLevel 属性](../html/d4e9ed1b-415b-e146-d416-d117d1b21624.htm "QualityOfServiceLevel 属性 ")

[Retain 属性](../html/f8c619d3-a883-5e37-4e7d-499ca4e5a73d.htm "Retain 属性 ")

[Topic 属性](../html/3fef75f8-da83-3b25-3617-d6f93f61a44a.htm "Topic 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MqttApplicationMessageTopic 属性 |

主题名（Topic Name）用于识别有效载荷数据应该被发布到哪一个信息通道。  
The Topic Name is used to identify which information channel the payload data should be published to.

**命名空间：**
 [HslCommunication.MQTT](68b3b27c-e8db-3b64-4a05-22be97dc0b2b.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public string Topic { get; set; }
```

```
Public Property Topic As String
	Get
	Set
```

```
public:
property String^ Topic {
	String^ get ();
	void set (String^ value);
}
```

```
member Topic : string with get, set
```

#### 属性值

类型：String

![](../icons/SectionExpanded.png)备注

UTF-8编码字符串中的字符数据必须是按照Unicode规范 [Unicode] 定义的和在RFC3629 [RFC3629] 中重申的有效的UTF-8格式。特别需要指出的是，
这些数据不能包含字符码在U+D800和U+DFFF之间的数据。如果服务端或客户端收到了一个包含无效UTF-8字符的控制报文，它必须关闭网络连接 [MQTT-1.5.3-1].
PUBLISH报文中的主题名不能包含通配符 [MQTT-3.3.2-2]。

![](../icons/SectionExpanded.png)参见

#### 引用

[MqttApplicationMessage 类](120454c9-8979-5ed7-c8c5-f68f2a971dd9.htm)

[HslCommunication.MQTT 命名空间](68b3b27c-e8db-3b64-4a05-22be97dc0b2b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## MqttApplicationMessage 方法

[原文連結](http://api.hslcommunication.cn/html/bc8f7ff9-24ce-98e4-32fd-4aebea6d55d5.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.MQTT](../html/68b3b27c-e8db-3b64-4a05-22be97dc0b2b.htm "HslCommunication.MQTT")

[MqttApplicationMessage 类](../html/120454c9-8979-5ed7-c8c5-f68f2a971dd9.htm "MqttApplicationMessage 类")

[MqttApplicationMessage 方法](../html/bc8f7ff9-24ce-98e4-32fd-4aebea6d55d5.htm "MqttApplicationMessage 方法")

[ToString 方法](../html/fab1a5b9-f28d-c8fa-5e7d-215aeb1ded62.htm "ToString 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MqttApplicationMessage 方法 |

[MqttApplicationMessage](120454c9-8979-5ed7-c8c5-f68f2a971dd9.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法 | [ToString](fab1a5b9-f28d-c8fa-5e7d-215aeb1ded62.htm) | (重写 ObjectToString.) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[MqttApplicationMessage 类](120454c9-8979-5ed7-c8c5-f68f2a971dd9.htm)

[HslCommunication.MQTT 命名空间](68b3b27c-e8db-3b64-4a05-22be97dc0b2b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ToString 方法 

[原文連結](http://api.hslcommunication.cn/html/fab1a5b9-f28d-c8fa-5e7d-215aeb1ded62.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.MQTT](../html/68b3b27c-e8db-3b64-4a05-22be97dc0b2b.htm "HslCommunication.MQTT")

[MqttApplicationMessage 类](../html/120454c9-8979-5ed7-c8c5-f68f2a971dd9.htm "MqttApplicationMessage 类")

[MqttApplicationMessage 方法](../html/bc8f7ff9-24ce-98e4-32fd-4aebea6d55d5.htm "MqttApplicationMessage 方法")

[ToString 方法](../html/fab1a5b9-f28d-c8fa-5e7d-215aeb1ded62.htm "ToString 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MqttApplicationMessageToString 方法 |

[缺少 "M:HslCommunication.MQTT.MqttApplicationMessage.ToString" 的 <summary> 文档]

**命名空间：**
 [HslCommunication.MQTT](68b3b27c-e8db-3b64-4a05-22be97dc0b2b.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public override string ToString()
```

```
Public Overrides Function ToString As String
```

```
public:
virtual String^ ToString() override
```

```
abstract ToString : unit -> string 
override ToString : unit -> string
```

#### 返回值

类型：String  

[缺少 "M:HslCommunication.MQTT.MqttApplicationMessage.ToString" 的 <returns> 文档]

![](../icons/SectionExpanded.png)参见

#### 引用

[MqttApplicationMessage 类](120454c9-8979-5ed7-c8c5-f68f2a971dd9.htm)

[HslCommunication.MQTT 命名空间](68b3b27c-e8db-3b64-4a05-22be97dc0b2b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## MqttClient 类

[原文連結](http://api.hslcommunication.cn/html/8626b3c4-5983-7308-94ba-65e3dd480b58.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.MQTT](../html/68b3b27c-e8db-3b64-4a05-22be97dc0b2b.htm "HslCommunication.MQTT")

[MqttClient 类](../html/8626b3c4-5983-7308-94ba-65e3dd480b58.htm "MqttClient 类")

[MqttClient 构造函数](../html/4817faf9-a999-ef74-9d36-d9efe7831b1b.htm "MqttClient 构造函数 ")

[MqttClient 属性](../html/05c7efc3-bf81-a5ea-a5a0-b67d68df3e10.htm "MqttClient 属性")

[MqttClient 方法](../html/5745a0d8-2dd3-b35a-76f5-d0b6dc6fb012.htm "MqttClient 方法")

[MqttClient 事件](../html/22de311a-df03-7e25-27c9-2aa41f8c46df.htm "MqttClient 事件")

[MqttClient 字段](../html/384cf7bb-b210-9388-1d20-ef614282bfe0.htm "MqttClient 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MqttClient 类 |

Mqtt协议的客户端实现，支持订阅消息，发布消息，详细的使用例子参考api文档  
The client implementation of the Mqtt protocol supports subscription messages and publishing messages. For detailed usage examples, refer to the api documentation.

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  [HslCommunication.Core.NetNetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)  
    [HslCommunication.Core.NetNetworkXBase](ab528a6a-e713-5202-59e0-90dd6a66de5f.htm)  
      HslCommunication.MQTTMqttClient

**命名空间：**
 [HslCommunication.MQTT](68b3b27c-e8db-3b64-4a05-22be97dc0b2b.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class MqttClient : NetworkXBase, IDisposable
```

```
Public Class MqttClient
	Inherits NetworkXBase
	Implements IDisposable
```

```
public ref class MqttClient : public NetworkXBase, 
	IDisposable
```

```
type MqttClient =  
    class
        inherit NetworkXBase
        interface IDisposable
    end
```

MqttClient 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [MqttClient](4817faf9-a999-ef74-9d36-d9efe7831b1b.htm) | 实例化一个默认的对象 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [ConnectionOptions](59421305-6be2-389d-e9e6-093321835fac.htm) | 获取当前的连接配置参数信息  Get current connection configuration parameter information |
| 公共属性 | [IsConnected](91d587d0-b6d8-96e7-0e8d-b73360d6820f.htm) | 获取或设置当前的服务器连接是否成功，定时获取本属性可用于实时更新连接状态信息。  Get or set whether the current server connection is successful or not. This property can be obtained regularly and can be used to update the connection status information in real time. |
| 公共属性代码示例 | [LogNet](46dfe698-46c6-73a2-5458-8fd568d264ce.htm) | 组件的日志工具，支持日志记录，只要实例化后，当前网络的基本信息，就以[DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm)等级进行输出  The component's logging tool supports logging. As long as the instantiation of the basic network information, the output will be output at [DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm) (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 公共属性 | [SubcribeTopics](8b884727-7781-3342-cd59-b5d4af0f2fac.htm) | 获取当前的客户端对象已经订阅的所有的Topic信息  Get all Topic information that the current client object has subscribed to |
| 公共属性 | [Tag](7d226192-f98d-e251-5fd2-3a664d0933de.htm) | 获取或设置当前客户端关联的自定义的对象内容  Gets or sets the custom object content associated with the current client |
| 公共属性代码示例 | [Token](80db00a7-aaa5-164f-e87c-af5400dc80fa.htm) | 网络类的身份令牌，在hsl协议的模式下会有效，在和设备进行通信的时候是无效的  Network-type identity tokens will be valid in the hsl protocol mode and will not be valid when communicating with the device (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 公共属性 | [UseTimerCheckDropped](3b46c470-9ee6-79b9-f84d-f40a746ef99b.htm) | 获取或设置是否启动定时器去检测当前客户端是否超时掉线。默认为 True  Get or set whether to start the timer to detect whether the current client timeout and disconnection. Default is True |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 受保护的方法 | [CheckRemoteToken](1d231984-101b-0fff-0daf-1d18581a4241.htm) | 检查当前的头子节信息的令牌是否是正确的，仅用于某些特殊的协议实现  Check whether the token of the current header subsection information is correct, only for some special protocol implementations (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 公共方法 | [ConnectClose](dde2bdc2-8bf6-9d78-d898-20f78621b055.htm) | 关闭Mqtt服务器的连接。  Close the connection to the Mqtt server. |
| 公共方法 | [ConnectCloseAsync](f33c8fa6-24a6-b352-43c9-b481fa3597b1.htm) | 关闭Mqtt服务器的连接。  Close the connection to the Mqtt server. |
| 公共方法 | [ConnectServer](c4a4b5c2-2c46-fe79-f347-0f312dc7bc21.htm) | 连接服务器，如果连接失败，请稍候重试。  Connect to the server. If the connection fails, try again later. |
| 公共方法 | [ConnectServerAsync](04751e70-b5b1-4364-b599-662925b0aacf.htm) | 连接服务器，如果连接失败，请稍候重试。  Connect to the server. If the connection fails, try again later. |
| 受保护的方法代码示例 | [CreateSocketAndConnect(String, Int32)](a153ae27-92e3-4fe7-b2e7-c8d207bebbb0.htm) | 创建一个新的socket对象并连接到远程的地址，默认超时时间为10秒钟，需要指定ip地址以及端口号信息  Create a new socket object and connect to the remote address. The default timeout is 10 seconds. You need to specify the IP address and port number. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法代码示例 | [CreateSocketAndConnect(IPEndPoint, Int32, IPEndPoint)](7ce6288f-c6db-e968-d676-f851d2b05c77.htm) | 创建一个新的socket对象并连接到远程的地址，需要指定远程终结点，超时时间（单位是毫秒），如果需要绑定本地的IP或是端口，传入 local对象  To create a new socket object and connect to the remote address, you need to specify the remote endpoint, the timeout period (in milliseconds), if you need to bind the local IP or port, pass in the local object (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法代码示例 | [CreateSocketAndConnect(String, Int32, Int32)](d0e14a5f-1e97-e59d-078f-4fbf0412b865.htm) | 创建一个新的socket对象并连接到远程的地址，需要指定ip地址以及端口号信息，还有超时时间，单位是毫秒  To create a new socket object and connect to a remote address, you need to specify the IP address and port number information, and the timeout period in milliseconds (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法代码示例 | [CreateSocketAndConnectAsync(String, Int32)](0227bc60-1c7f-3c1d-0cd6-e49bb904d6df.htm) | 创建一个新的socket对象并连接到远程的地址，默认超时时间为10秒钟，需要指定ip地址以及端口号信息  Create a new socket object and connect to the remote address. The default timeout is 10 seconds. You need to specify the IP address and port number. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法代码示例 | [CreateSocketAndConnectAsync(IPEndPoint, Int32, IPEndPoint)](ccec07bf-cfa4-b302-210e-441821e095a4.htm) | 创建一个新的socket对象并连接到远程的地址，需要指定远程终结点，超时时间（单位是毫秒），如果需要绑定本地的IP或是端口，传入 local对象  To create a new socket object and connect to the remote address, you need to specify the remote endpoint, the timeout period (in milliseconds), if you need to bind the local IP or port, pass in the local object (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法代码示例 | [CreateSocketAndConnectAsync(String, Int32, Int32)](7ec25883-60e2-c7ac-2f90-0ea1ae0c4945.htm) | 创建一个新的socket对象并连接到远程的地址，需要指定ip地址以及端口号信息，还有超时时间，单位是毫秒  To create a new socket object and connect to a remote address, you need to specify the IP address and port number information, and the timeout period in milliseconds (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [DeleteFileByName](59afa4a4-3f3c-9831-26ee-c7c13a213bf7.htm) | 删除一个指定的文件，如果文件不存在，直接返回 True，如果文件存在则直接删除，删除成功返回 True，如果发生了异常，返回False  Delete a specified file, if the file does not exist, return True directly, if the file exists, delete it directly, if the deletion is successful, return True, if an exception occurs, return  False (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 公共方法 | [Dispose](deed427b-fd4f-771a-f4d5-7c0d6392f1ec.htm) | 释放被 MqttClient 使用的所有资源 |
| 受保护的方法 | [Dispose(Boolean)](973e0523-d8d4-b870-2e38-d9e2b9b12989.htm) | 释放当前的对象 |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法 | [GetSubscribeTopic](f32877cd-5f40-92dc-9203-f2ac9ef0964e.htm) | 获取已经订阅的主题信息，方便针对不同的界面订阅不同的主题。  Get subscribed topic information, which is convenient for subscribing to different topics for different interfaces. |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法代码示例 | [PublishMessage](8175840c-06e5-f82e-3079-ce613f81e2cc.htm) | 发布一个MQTT协议的消息到服务器。该消息包含主题，负载数据，消息等级，是否保留信息。  Publish an MQTT protocol message to the server. The message contains the subject, payload data, message level, and whether to retain information. |
| 公共方法代码示例 | [PublishMessageAsync](ea62a552-79aa-2470-7db4-e249268edce1.htm) | 发布一个MQTT协议的消息到服务器。该消息包含主题，负载数据，消息等级，是否保留信息。  Publish an MQTT protocol message to the server. The message contains the subject, payload data, message level, and whether to retain information. |
| 受保护的方法 | [Receive(SslStream, Int32, Int32, ActionInt64, Int64)](50ad65ee-3ff1-6a08-31f3-a29090797796.htm) | 接收固定长度的字节数组，允许指定超时时间，默认为60秒，当length大于0时，接收固定长度的数据内容，当length小于0时，接收不大于2048长度的随机数据信息  Receiving a fixed-length byte array, allowing a specified timeout time. The default is 60 seconds. When length is greater than 0, fixed-length data content is received. When length is less than 0, random data information of a length not greater than 2048 is received. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [Receive(Socket, Int32, Int32, ActionInt64, Int64)](28c887ec-7a68-e90a-f531-b2298ded707b.htm) | 接收固定长度的字节数组，允许指定超时时间，默认为60秒，当length大于0时，接收固定长度的数据内容，当length小于0时，接收不大于2048长度的随机数据信息  Receiving a fixed-length byte array, allowing a specified timeout time. The default is 60 seconds. When length is greater than 0, fixed-length data content is received. When length is less than 0, random data information of a length not greater than 2048 is received. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [Receive(SslStream, Byte, Int32, Int32, Int32, ActionInt64, Int64)](d6fbd69f-3aa1-9f84-139a-04003a9ce7c0.htm) | 接收固定长度的字节数组，允许指定超时时间，默认为60秒，当length大于0时，接收固定长度的数据内容，当length小于0时，buffer长度的缓存数据  Receiving a fixed-length byte array, allowing a specified timeout time. The default is 60 seconds. When length is greater than 0, fixed-length data content is received. When length is less than 0, random data information of a length not greater than 2048 is received. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [Receive(Socket, Byte, Int32, Int32, Int32, ActionInt64, Int64)](403209a0-350c-ede4-a17c-8ab8ad5ddc5b.htm) | 接收固定长度的字节数组，允许指定超时时间，默认为60秒，当length大于0时，接收固定长度的数据内容，当length小于0时，buffer长度的缓存数据  Receiving a fixed-length byte array, allowing a specified timeout time. The default is 60 seconds. When length is greater than 0, fixed-length data content is received. When length is less than 0, random data information of a length not greater than 2048 is received. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [ReceiveAndCheckBytes](1d627617-63a7-7078-86e8-cfb10c3ad500.htm) | [自校验] 接收一条完整的同步数据，包含头子节和内容字节，基础的数据，如果结果异常，则结束通讯  [Self-checking] Receive a complete synchronization data, including header subsection and content bytes, basic data, if the result is abnormal, the communication ends (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [ReceiveAndCheckBytesAsync](a9e13d42-804a-afee-4c86-d59ad03ec469.htm) | [自校验] 接收一条完整的同步数据，包含头子节和内容字节，基础的数据，如果结果异常，则结束通讯  [Self-checking] Receive a complete synchronization data, including header subsection and content bytes, basic data, if the result is abnormal, the communication ends (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [ReceiveAsync(SslStream, Int32, Int32, ActionInt64, Int64)](548ecb57-8a13-a71d-0a29-a71c282f94d8.htm) | 接收固定长度的字节数组，允许指定超时时间，默认为60秒，当length大于0时，接收固定长度的数据内容，当length小于0时，接收不大于2048长度的随机数据信息  Receiving a fixed-length byte array, allowing a specified timeout time. The default is 60 seconds. When length is greater than 0, fixed-length data content is received. When length is less than 0, random data information of a length not greater than 2048 is received. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [ReceiveAsync(Socket, Int32, Int32, ActionInt64, Int64)](6edaea53-6855-bfcd-33e6-5ca3df268636.htm) | 接收固定长度的字节数组，允许指定超时时间，默认为60秒，当length大于0时，接收固定长度的数据内容，当length小于0时，接收不大于2048长度的随机数据信息  Receiving a fixed-length byte array, allowing a specified timeout time. The default is 60 seconds. When length is greater than 0, fixed-length data content is received. When length is less than 0, random data information of a length not greater than 2048 is received. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [ReceiveAsync(SslStream, Byte, Int32, Int32, Int32, ActionInt64, Int64)](6104ebb8-3044-1bc6-3972-add79588e90c.htm) | 接收固定长度的字节数组，允许指定超时时间，默认为60秒，当length大于0时，接收固定长度的数据内容，当length小于0时，buffer长度的缓存数据  Receiving a fixed-length byte array, allowing a specified timeout time. The default is 60 seconds. When length is greater than 0, fixed-length data content is received. When length is less than 0, random data information of a length not greater than 2048 is received. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [ReceiveAsync(Socket, Byte, Int32, Int32, Int32, ActionInt64, Int64)](5d183a4a-19ba-46b2-9338-a4a88f2a7c70.htm) | 接收固定长度的字节数组，允许指定超时时间，默认为60秒，当length大于0时，接收固定长度的数据内容，当length小于0时，buffer长度的缓存数据  Receiving a fixed-length byte array, allowing a specified timeout time. The default is 60 seconds. When length is greater than 0, fixed-length data content is received. When length is less than 0, random data information of a length not greater than 2048 is received. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [ReceiveByMessage](33d5e9f5-67f0-c2c9-4ad3-c59ebfdde4d7.htm) | 接收一条完整的 数据内容，需要指定超时时间，单位为毫秒。   Receive a complete data content, Need to specify a timeout period in milliseconds (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [ReceiveByMessageAsync](e75c8b31-a6df-7f85-cd9f-d5e6a095a983.htm) | 接收一条完整的 数据内容，需要指定超时时间，单位为毫秒。   Receive a complete data content, Need to specify a timeout period in milliseconds (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [ReceiveBytesContentFromSocket](190b4c27-554b-713d-0f49-7942516e96f0.htm) | [自校验] 从网络中接收一串字节数据，如果结果异常，则结束通讯  [Self-checking] Receive a string of byte data from the network. If the result is abnormal, the communication ends. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [ReceiveBytesContentFromSocketAsync](2e3e6ee6-4d7e-a5c9-f642-13969c4c3b78.htm) | [自校验] 从网络中接收一串字节数据，如果结果异常，则结束通讯  [Self-checking] Receive a string of byte data from the network. If the result is abnormal, the communication ends. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [ReceiveCommandLineFromSocket(Socket, Byte, Int32)](f6fbd8fd-5f1f-f4c2-2ab4-3b7f1166dd4c.htm) | 接收一行命令数据，需要自己指定这个结束符，默认超时时间为60秒，也即是60000，单位是毫秒  To receive a line of command data, you need to specify the terminator yourself. The default timeout is 60 seconds, which is 60,000, in milliseconds. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [ReceiveCommandLineFromSocket(Socket, Byte, Byte, Int32)](9fb845a9-a46c-2285-2076-33bcaf62c16a.htm) | 接收一行命令数据，需要自己指定这个结束符，默认超时时间为60秒，也即是60000，单位是毫秒  To receive a line of command data, you need to specify the terminator yourself. The default timeout is 60 seconds, which is 60,000, in milliseconds. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [ReceiveCommandLineFromSocketAsync(Socket, Byte, Int32)](6f21495a-3b08-501e-3219-30f00b831a36.htm) | 接收一行命令数据，需要自己指定这个结束符，默认超时时间为60秒，也即是60000，单位是毫秒  To receive a line of command data, you need to specify the terminator yourself. The default timeout is 60 seconds, which is 60,000, in milliseconds. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [ReceiveCommandLineFromSocketAsync(Socket, Byte, Byte, Int32)](bcc92f0b-de8e-c78f-4aba-fff2cdd34c82.htm) | 接收一行命令数据，需要自己指定这个结束符，默认超时时间为60秒，也即是60000，单位是毫秒  To receive a line of command data, you need to specify the terminator yourself. The default timeout is 60 seconds, which is 60,000, in milliseconds. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [ReceiveFileFromSocket(Socket, Stream, ActionInt64, Int64)](a227f0b6-3050-e5d1-79c8-88ea9719bd06.htm) | [自校验] 从网络中接收一个文件，写入数据流，如果结果异常，则结束通讯，参数顺序文件名，文件大小，文件标识，上传人  [Self-checking] Receive a file from the network. If the result is abnormal, the communication ends. (继承自 [NetworkXBase](ab528a6a-e713-5202-59e0-90dd6a66de5f.htm)。) |
| 受保护的方法 | [ReceiveFileFromSocket(Socket, String, ActionInt64, Int64)](6bcd5a10-d79c-649e-fb65-cd78fecf99b7.htm) | [自校验] 从网络中接收一个文件，如果结果异常，则结束通讯  [Self-checking] Receive a file from the network. If the result is abnormal, the communication ends. (继承自 [NetworkXBase](ab528a6a-e713-5202-59e0-90dd6a66de5f.htm)。) |
| 受保护的方法 | [ReceiveFileFromSocketAsync(Socket, Stream, ActionInt64, Int64)](e3b415d2-58a9-7ab6-0097-0bdc736ca0dc.htm) | [自校验] 从网络中接收一个文件，写入数据流，如果结果异常，则结束通讯，参数顺序文件名，文件大小，文件标识，上传人  [Self-checking] Receive a file from the network. If the result is abnormal, the communication ends. (继承自 [NetworkXBase](ab528a6a-e713-5202-59e0-90dd6a66de5f.htm)。) |
| 受保护的方法 | [ReceiveFileFromSocketAsync(Socket, String, ActionInt64, Int64)](ab9b6255-580b-464b-447a-5de92b6f3dd9.htm) | [自校验] 从网络中接收一个文件，如果结果异常，则结束通讯  [Self-checking] Receive a file from the network. If the result is abnormal, the communication ends. (继承自 [NetworkXBase](ab528a6a-e713-5202-59e0-90dd6a66de5f.htm)。) |
| 受保护的方法 | [ReceiveFileHeadFromSocket](37b7d2d2-faeb-2263-627d-6a74bcbce47d.htm) | [自校验] 从套接字中接收文件头信息  [Self-checking] Receive file header information from socket (继承自 [NetworkXBase](ab528a6a-e713-5202-59e0-90dd6a66de5f.htm)。) |
| 受保护的方法 | [ReceiveFileHeadFromSocketAsync](9f10bf85-76d0-c721-60f1-1b42425a52af.htm) | [自校验] 从套接字中接收文件头信息  [Self-checking] Receive file header information from socket (继承自 [NetworkXBase](ab528a6a-e713-5202-59e0-90dd6a66de5f.htm)。) |
| 受保护的方法 | [ReceiveHslMessage](9cd1bc29-c99a-5426-1fb1-047268154fe3.htm) | 接收一条hsl协议的数据信息，自动解析，解压，解码操作，获取最后的实际的数据，接收结果依次为暗号，用户码，负载数据  Receive a piece of hsl protocol data information, automatically parse, decompress, and decode operations to obtain the last actual data. The result is a opCode, user code, and payload data in order. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [ReceiveHslMessageAsync](63e3bd3f-efd6-fc4b-9361-5caab03c47d9.htm) | 接收一条hsl协议的数据信息，自动解析，解压，解码操作，获取最后的实际的数据，接收结果依次为暗号，用户码，负载数据  Receive a piece of hsl protocol data information, automatically parse, decompress, and decode operations to obtain the last actual data. The result is a opCode, user code, and payload data in order. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [ReceiveMqttFile](a529cfa6-f57f-bede-b266-4f8942c4934a.htm) | 使用MQTT协议从网络接收字节数组，然后写入文件或流中，支持进度报告  Use MQTT protocol to receive byte array from the network, and then write it to file or stream, support progress report (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [ReceiveMqttFileAsync](2b725eb6-1051-c1c8-5e31-9787204d9ffa.htm) | 使用MQTT协议从网络接收字节数组，然后写入文件或流中，支持进度报告  Use MQTT protocol to receive byte array from the network, and then write it to file or stream, support progress report (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [ReceiveMqttMessage(SslStream, Int32, ActionInt64, Int64)](c268dc7d-d3d1-3806-8e2c-8de3a7736cc2.htm) | 接收一条完整的MQTT协议的报文信息，包含控制码和负载数据  Receive a message of a completed MQTT protocol, including control code and payload data (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [ReceiveMqttMessage(Socket, Int32, ActionInt64, Int64)](81c27af3-012e-665f-129c-7aae05c77021.htm) | 接收一条完整的MQTT协议的报文信息，包含控制码和负载数据  Receive a message of a completed MQTT protocol, including control code and payload data (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [ReceiveMqttMessageAsync(SslStream, Int32, ActionInt64, Int64)](5fa070b0-e679-7201-df00-bf7c34811100.htm) | 接收一条完整的MQTT协议的报文信息，包含控制码和负载数据  Receive a message of a completed MQTT protocol, including control code and payload data (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [ReceiveMqttMessageAsync(Socket, Int32, ActionInt64, Int64)](f6ec22af-2c46-7ace-716c-a198e9a70a82.htm) | 接收一条完整的MQTT协议的报文信息，包含控制码和负载数据  Receive a message of a completed MQTT protocol, including control code and payload data (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [ReceiveMqttStream](dcd504a7-4231-7f98-2abc-dd6aed3ed1e2.htm) | 使用MQTT协议从socket接收指定长度的字节数组，然后全部写入到流中，可以指定进度报告  Use the MQTT protocol to receive a byte array of specified length from the socket, and then write all of them to the stream, and you can specify a progress report (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [ReceiveMqttStreamAsync](8e562525-0bdf-6592-82ea-c27cb6aa339b.htm) | 使用MQTT协议从socket接收指定长度的字节数组，然后全部写入到流中，可以指定进度报告  Use the MQTT protocol to receive a byte array of specified length from the socket, and then write all of them to the stream, and you can specify a progress report (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [ReceiveRedisCommand](896bd785-c06b-c7f6-475f-cc78ca50b201.htm) | 从网络接收一条完整的redis报文的消息  Receive a complete redis message from the network (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [ReceiveRedisCommandAsync](7da184d3-e92d-fb42-9c44-a3c00f990157.htm) | 从网络接收一条完整的redis报文的消息  Receive a complete redis message from the network (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [ReceiveRedisCommandString](3299524f-73b1-476e-18f3-cddd811a5b2e.htm) | 接收一行基于redis协议的字符串的信息，需要指定固定的长度  Receive a line of information based on the redis protocol string, you need to specify a fixed length (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [ReceiveRedisCommandStringAsync](9bdf24f8-847b-5177-d45d-2524e58774d9.htm) | 接收一行基于redis协议的字符串的信息，需要指定固定的长度  Receive a line of information based on the redis protocol string, you need to specify a fixed length (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [ReceiveStringArrayContentFromSocket](d32e83aa-c62b-0c81-ebe6-1482fe8a5bfa.htm) | [自校验] 从网络中接收一个字符串数组，如果结果异常，则结束通讯  [Self-check] Receive an array of strings from the network. If the result is abnormal, the communication ends. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [ReceiveStringArrayContentFromSocketAsync](ca8a25f7-200d-b95c-8b31-d1374c77c280.htm) | [自校验] 从网络中接收一个字符串数组，如果结果异常，则结束通讯  [Self-check] Receive an array of strings from the network. If the result is abnormal, the communication ends. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [ReceiveStringContentFromSocket](ad3b2915-8cfb-ff45-f24f-1eb255bd289e.htm) | [自校验] 从网络中接收一个字符串数据，如果结果异常，则结束通讯  [Self-checking] Receive a string of data from the network. If the result is abnormal, the communication ends. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [ReceiveStringContentFromSocketAsync](3cd827b7-1fca-6184-0067-85c016d3eb15.htm) | [自校验] 从网络中接收一个字符串数据，如果结果异常，则结束通讯  [Self-checking] Receive a string of data from the network. If the result is abnormal, the communication ends. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [Send(SslStream, Byte)](8fad1450-4dc9-4505-2873-09fd128d0f09.htm) | 发送消息给套接字，直到完成的时候返回，经过测试，本方法是线程安全的。  Send a message to the socket until it returns when completed. After testing, this method is thread-safe. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [Send(Socket, Byte)](d98c31b7-b055-011c-0546-93434bd77af5.htm) | 发送消息给套接字，直到完成的时候返回，经过测试，本方法是线程安全的。  Send a message to the socket until it returns when completed. After testing, this method is thread-safe. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [Send(SslStream, Byte, Int32, Int32)](b701a27a-0c44-1cf4-7bdf-73518ea1f2e6.htm) | 发送消息给套接字，直到完成的时候返回，经过测试，本方法是线程安全的。  Send a message to the socket until it returns when completed. After testing, this method is thread-safe. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [Send(Socket, Byte, Int32, Int32)](09e7c105-bfea-0186-03fe-981d8124c208.htm) | 发送消息给套接字，直到完成的时候返回，经过测试，本方法是线程安全的。  Send a message to the socket until it returns when completed. After testing, this method is thread-safe. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [SendAccountAndCheckReceive](86806e1d-1b59-4333-9c64-65fc445a64e1.htm) | [自校验] 直接发送字符串数组并确认对方接收完成数据，如果结果异常，则结束通讯  [Self-check] Send string array directly and confirm that the other party has received the completed data. If the result is abnormal, the communication ends. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [SendAccountAndCheckReceiveAsync](cb6aa97e-a3f7-e64a-e68a-636756d28207.htm) | [自校验] 直接发送字符串数组并确认对方接收完成数据，如果结果异常，则结束通讯  [Self-check] Send string array directly and confirm that the other party has received the completed data. If the result is abnormal, the communication ends. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [SendAsync(SslStream, Byte)](9ab7bc43-5a28-5e39-ec77-811d091507a5.htm) | 发送消息给套接字，直到完成的时候返回，经过测试，本方法是线程安全的。  Send a message to the socket until it returns when completed. After testing, this method is thread-safe. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [SendAsync(Socket, Byte)](dec1ad4d-89b9-fc85-1426-8e61acb444aa.htm) | 发送消息给套接字，直到完成的时候返回，经过测试，本方法是线程安全的。  Send a message to the socket until it returns when completed. After testing, this method is thread-safe. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [SendAsync(SslStream, Byte, Int32, Int32)](4ce38aad-0f1a-da9b-aa13-a755fa236ac7.htm) | 发送消息给套接字，直到完成的时候返回，经过测试，本方法是线程安全的。  Send a message to the socket until it returns when completed. After testing, this method is thread-safe. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [SendAsync(Socket, Byte, Int32, Int32)](e18a57a9-9caf-193f-9971-fbd41bba9e61.htm) | 发送消息给套接字，直到完成的时候返回，经过测试，本方法是线程安全的。  Send a message to the socket until it returns when completed. After testing, this method is thread-safe. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [SendBaseAndCheckReceive](3ab1c2a3-2a33-b6f8-b648-98ca72b4f139.htm) | [自校验] 发送字节数据并确认对方接收完成数据，如果结果异常，则结束通讯  [Self-check] Send the byte data and confirm that the other party has received the completed data. If the result is abnormal, the communication ends. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [SendBaseAndCheckReceiveAsync](b6344eac-286c-c80f-ed85-56c5784065a5.htm) | [自校验] 发送字节数据并确认对方接收完成数据，如果结果异常，则结束通讯  [Self-check] Send the byte data and confirm that the other party has received the completed data. If the result is abnormal, the communication ends. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [SendBytesAndCheckReceive](ddbf44f9-7dff-964d-3e8c-d4198887183f.htm) | [自校验] 发送字节数据并确认对方接收完成数据，如果结果异常，则结束通讯  [Self-check] Send the byte data and confirm that the other party has received the completed data. If the result is abnormal, the communication ends. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [SendBytesAndCheckReceiveAsync](038563fa-a4e0-2103-4f93-f52f1ee2ded3.htm) | [自校验] 发送字节数据并确认对方接收完成数据，如果结果异常，则结束通讯  [Self-check] Send the byte data and confirm that the other party has received the completed data. If the result is abnormal, the communication ends. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [SendFileAndCheckReceive(Socket, Stream, String, String, String, ActionInt64, Int64)](602977f6-83fa-486f-7b4e-7913594bb847.htm) | [自校验] 将流数据发送至套接字，具体发送细节将在继承类中实现，如果结果异常，则结束通讯  [Self-checking] Send stream data to the socket. The specific sending details will be implemented in the inherited class. If the result is abnormal, the communication will be terminated (继承自 [NetworkXBase](ab528a6a-e713-5202-59e0-90dd6a66de5f.htm)。) |
| 受保护的方法 | [SendFileAndCheckReceive(Socket, String, String, String, String, ActionInt64, Int64)](cad4c115-2f7b-5301-36e4-66b6da8d864f.htm) | [自校验] 将文件数据发送至套接字，具体发送细节将在继承类中实现，如果结果异常，则结束通讯  [Self-checking] Send the file data to the socket. The specific sending details will be implemented in the inherited class. If the result is abnormal, the communication will end (继承自 [NetworkXBase](ab528a6a-e713-5202-59e0-90dd6a66de5f.htm)。) |
| 受保护的方法 | [SendFileAndCheckReceiveAsync(Socket, Stream, String, String, String, ActionInt64, Int64)](8ad77d7e-8710-5bb5-e39a-9e613957cf72.htm) | [自校验] 将流数据发送至套接字，具体发送细节将在继承类中实现，如果结果异常，则结束通讯  [Self-checking] Send stream data to the socket. The specific sending details will be implemented in the inherited class. If the result is abnormal, the communication will be terminated (继承自 [NetworkXBase](ab528a6a-e713-5202-59e0-90dd6a66de5f.htm)。) |
| 受保护的方法 | [SendFileAndCheckReceiveAsync(Socket, String, String, String, String, ActionInt64, Int64)](9eff5d17-5f46-6900-b36b-f06e76f3605e.htm) | [自校验] 将文件数据发送至套接字，具体发送细节将在继承类中实现，如果结果异常，则结束通讯  [Self-checking] Send the file data to the socket. The specific sending details will be implemented in the inherited class. If the result is abnormal, the communication will end (继承自 [NetworkXBase](ab528a6a-e713-5202-59e0-90dd6a66de5f.htm)。) |
| 受保护的方法 | [SendFileStreamToSocket](73df4306-17c8-3649-8214-b0704cd009af.htm) | [自校验] 将文件数据发送至套接字，如果结果异常，则结束通讯  [Self-check] Send the file data to the socket. If the result is abnormal, the communication is ended. (继承自 [NetworkXBase](ab528a6a-e713-5202-59e0-90dd6a66de5f.htm)。) |
| 受保护的方法 | [SendFileStreamToSocketAsync](41b6ebdb-22c3-fe61-113b-a60f1674750c.htm) | [自校验] 将文件数据发送至套接字，如果结果异常，则结束通讯  [Self-check] Send the file data to the socket. If the result is abnormal, the communication is ended. (继承自 [NetworkXBase](ab528a6a-e713-5202-59e0-90dd6a66de5f.htm)。) |
| 受保护的方法 | [SendMqttFile(Socket, Stream, String, String, ActionInt64, Int64, AesCryptography, HslCancelToken)](46ad6f88-4010-ff7d-12d7-196516140b31.htm) | 使用MQTT协议将一个数据流发送到网络上去，需要保存的文件名，可选指定文件描述信息，进度报告  Use the MQTT protocol to send a data stream to the network, the file name that needs to be saved, optional file description information, progress report (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [SendMqttFile(Socket, String, String, String, ActionInt64, Int64, AesCryptography, HslCancelToken)](29af3d39-9ebe-f5a6-b203-69243b1ec55e.htm) | 使用MQTT协议将一个文件发送到网络上去，需要指定文件名，保存的文件名，可选指定文件描述信息，进度报告  To send a file to the network using the MQTT protocol, you need to specify the file name, the saved file name, optionally specify the file description information, and the progress report (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [SendMqttFileAsync(Socket, Stream, String, String, ActionInt64, Int64, AesCryptography, HslCancelToken)](38401fac-1a74-bbfb-2847-27fa314ef1e7.htm) | 使用MQTT协议将一个数据流发送到网络上去，需要保存的文件名，可选指定文件描述信息，进度报告  Use the MQTT protocol to send a data stream to the network, the file name that needs to be saved, optional file description information, progress report (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [SendMqttFileAsync(Socket, String, String, String, ActionInt64, Int64, AesCryptography, HslCancelToken)](92b852f5-9536-9588-4919-6faf4e6b835a.htm) | 使用MQTT协议将一个文件发送到网络上去，需要指定文件名，保存的文件名，可选指定文件描述信息，进度报告  To send a file to the network using the MQTT protocol, you need to specify the file name, the saved file name, optionally specify the file description information, and the progress report (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [SendMqttStream](12db1776-8fb3-7d64-4331-2df07128c59a.htm) | 使用MQTT协议将流中的数据读取到字节数组，然后都写入到socket里面，可以指定进度报告，主要用于将文件发送到网络。  Use the MQTT protocol to read the data in the stream into a byte array, and then write them all into the socket. You can specify a progress report, which is mainly used to send files to the network. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [SendMqttStreamAsync](8493fd7e-4ba3-cc9c-efbf-18747cae084d.htm) | 使用MQTT协议将流中的数据读取到字节数组，然后都写入到socket里面，可以指定进度报告，主要用于将文件发送到网络。  Use the MQTT protocol to read the data in the stream into a byte array, and then write them all into the socket. You can specify a progress report, which is mainly used to send files to the network. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [SendStreamToSocket](f6220302-c53a-4f8a-ab69-7827a9779519.htm) | 发送一个流的所有数据到指定的网络套接字，需要指定发送的数据长度，支持按照百分比的进度报告  Send all the data of a stream to the specified network socket. You need to specify the length of the data to be sent. It supports the progress report in percentage. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [SendStreamToSocketAsync](e4e82bab-973b-99b8-a3da-098c009d244f.htm) | 发送一个流的所有数据到指定的网络套接字，需要指定发送的数据长度，支持按照百分比的进度报告  Send all the data of a stream to the specified network socket. You need to specify the length of the data to be sent. It supports the progress report in percentage. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [SendStringAndCheckReceive(Socket, Int32, String)](4aa4b2d4-98fd-33ce-9a08-4a4f9529f950.htm) | [自校验] 直接发送字符串数据并确认对方接收完成数据，如果结果异常，则结束通讯  [Self-checking] Send string data directly and confirm that the other party has received the completed data. If the result is abnormal, the communication ends. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [SendStringAndCheckReceive(Socket, Int32, String)](b7bdf13e-1a1d-cbe5-f511-e6a5e0be4c69.htm) | [自校验] 直接发送字符串数组并确认对方接收完成数据，如果结果异常，则结束通讯  [Self-check] Send string array directly and confirm that the other party has received the completed data. If the result is abnormal, the communication ends. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [SendStringAndCheckReceiveAsync(Socket, Int32, String)](41c7589d-8116-66a0-5311-685001239af6.htm) | [自校验] 直接发送字符串数据并确认对方接收完成数据，如果结果异常，则结束通讯  [Self-checking] Send string data directly and confirm that the other party has received the completed data. If the result is abnormal, the communication ends. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [SendStringAndCheckReceiveAsync(Socket, Int32, String)](814826a1-27c4-ca81-424c-ae50df6eb0cb.htm) | [自校验] 直接发送字符串数组并确认对方接收完成数据，如果结果异常，则结束通讯  [Self-check] Send string array directly and confirm that the other party has received the completed data. If the result is abnormal, the communication ends. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 公共方法代码示例 | [SubscribeMessage(String)](983ebe00-28da-7d06-f7b0-92c0acae3000.htm) | 从服务器订阅一个或多个主题信息  Subscribe to one or more topics from the server |
| 公共方法代码示例 | [SubscribeMessage(String)](82824f68-6335-bcf0-58a6-b10915502d46.htm) | 从服务器订阅一个或多个主题信息  Subscribe to one or more topics from the server |
| 公共方法 | [SubscribeMessage(MqttSubscribeMessage)](0957c4dd-38f2-567b-b013-6757ca59f107.htm) | 向服务器订阅一个主题消息，可以指定订阅的主题数组，订阅的质量等级，还有消息标识符  To subscribe to a topic message from the server, you can specify the subscribed topic array, the subscription quality level, and the message identifier |
| 公共方法代码示例 | [SubscribeMessageAsync(String)](a7833b6a-d4e3-6a8a-aae9-d8fa3a47d78b.htm) | 从服务器订阅一个或多个主题信息  Subscribe to one or more topics from the server |
| 公共方法代码示例 | [SubscribeMessageAsync(String)](3fff899a-b9d0-c3c1-9a01-e64dfb1a0a82.htm) | 从服务器订阅一个或多个主题信息  Subscribe to one or more topics from the server |
| 公共方法 | [ToString](07418f5b-741a-e8c3-1f2b-977f55554e5b.htm) | (重写 [NetworkXBaseToString](08bf006b-af6e-7c58-6935-8e8c387e96ab.htm).) |
| 公共方法代码示例 | [UnSubscribeMessage(String)](62cbd093-a214-fea7-7687-cac520b1b43a.htm) | 取消订阅指定的主题信息，取消之后，就不再接收当前主题的数据，除非服务器强制推送  Unsubscribe from the specified topic information. After cancellation, the data of the current topic will no longer be received unless the server forces push |
| 公共方法代码示例 | [UnSubscribeMessage(String)](ca8a3e8e-7d28-3092-6218-ab521dd806d2.htm) | 取消订阅多个主题信息，取消之后，当前的订阅数据就不在接收到，除非服务器强制推送。  Unsubscribe from multiple topic information. After cancellation, the current subscription data will not be received unless the server forces it to push it. |
| 公共方法代码示例 | [UnSubscribeMessageAsync(String)](cf59d92a-4e0a-647a-df99-3bb3517983d8.htm) | 取消订阅指定的主题信息，取消之后，就不再接收当前主题的数据，除非服务器强制推送  Unsubscribe from the specified topic information. After cancellation, the data of the current topic will no longer be received unless the server forces push |
| 公共方法代码示例 | [UnSubscribeMessageAsync(String)](36409683-6754-83b8-e610-1acd93d0c3f6.htm) | 取消订阅多个主题信息，取消之后，当前的订阅数据就不在接收到，除非服务器强制推送。  Unsubscribe from multiple topic information. After cancellation, the current subscription data will not be received unless the server forces it to push it. |
| 受保护的方法 | [WriteStreamFromSocket](97bf95ea-51fe-8f2e-5da6-1c83766aaa78.htm) | 从套接字中接收所有的数据然后写入到指定的流当中去，需要指定数据的长度，支持按照百分比进行进度报告  Receives all data from the socket and writes it to the specified stream. The length of the data needs to be specified, and progress reporting is supported in percentage. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [WriteStreamFromSocketAsync](b4cf7f6d-1e0c-a266-4c49-683fdb14a922.htm) | 从套接字中接收所有的数据然后写入到指定的流当中去，需要指定数据的长度，支持按照百分比进行进度报告  Receives all data from the socket and writes it to the specified stream. The length of the data needs to be specified, and progress reporting is supported in percentage. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)事件

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共事件 | [OnClientConnected](7cad095d-7edb-9867-b8f3-253274bd1325.htm) | 当客户端连接成功触发事件，就算是重新连接服务器后，也是会触发的  The event is triggered when the client is connected successfully, even after reconnecting to the server. |
| 公共事件 | [OnMqttMessageReceived](6d7f1a92-1ed6-509c-4ad0-83bd46fc1d74.htm) | 当接收到Mqtt订阅的信息的时候触发 |
| 公共事件 | [OnNetworkError](401ab742-c9d0-2d08-2b79-e0b72abd3a5c.htm) | 当网络发生异常的时候触发的事件，用户应该在事件里进行重连服务器 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)字段

|  | 名称 | 说明 |
| --- | --- | --- |
| 受保护的字段 | [CoreSocket](6831da9f-21e0-8967-f17c-641b330e0d80.htm) | 对客户端而言是的通讯用的套接字，对服务器来说是用于侦听的套接字  A communication socket for the client, or a listening socket for the server (继承自 [NetworkXBase](ab528a6a-e713-5202-59e0-90dd6a66de5f.htm)。) |
| 受保护的字段 | [fileCacheSize](2328fb89-cfda-7fa9-88d0-98f1ae3c7c54.htm) | 文件传输的时候的缓存大小，直接影响传输的速度，值越大，传输速度越快，越占内存，默认为100K大小  The size of the cache during file transfer directly affects the speed of the transfer. The larger the value, the faster the transfer speed and the more memory it takes. The default size is 100K. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)备注

这是一个MQTT的客户端实现，参照MQTT协议的3.1.1版本设计实现的。服务器可以是其他的组件提供的，其他的可以参考示例  
This is an MQTT client implementation, designed and implemented with reference to version 3.1.1 of the MQTT protocol. The server can be provided by other components.

![](../icons/SectionExpanded.png)示例

简单的实例化

[复制](# "复制")

```
// 简单的实例化例子
MqttClient mqttClient = new MqttClient( new MqttConnectionOptions( )
{
    ClientId = "ABC",                     // 客户端的唯一的ID信息
    IpAddress = "127.0.0.1",              // 服务器的地址
} );
```

带用户名密码的实例化

[复制](# "复制")

```
// 如果有密码的情况
MqttClient mqttClient = new MqttClient( new MqttConnectionOptions( )
{
    ClientId = "ABC",                                            // 客户端的唯一的ID信息
    IpAddress = "127.0.0.1",                                     // 服务器的地址
    Credentials = new MqttCredential( "admin", "123456" )        // 设置了用户名和密码
} );
```

如果使用证书的情况

[复制](# "复制")

```
// 如果有证书的情况
MqttClient mqttClient = new MqttClient( new MqttConnectionOptions( )
{
    ClientId = "ABC",                                            // 客户端的唯一的ID信息
    IpAddress = "127.0.0.1",                                     // 服务器的地址
    Credentials = new MqttCredential( "admin", "123456" ),       // 设置了用户名和密码

    CertificateFile = @"D:\emqx\etc\certs\cacert.pem",           // 指定证书的路径，这里选择的是emqx的测试证书
} );
```

简单的加密操作

[复制](# "复制")

```
// 如果连接的hslcommunication自身的服务器，还支持直接使用通信加密技术，配置非常的简单，这样就可以防止通讯的过程中，被抓包发现账户名，密码，通信数据等信息
MqttClient mqttClient = new MqttClient( new MqttConnectionOptions( )
{
    ClientId = "ABC",                                            // 客户端的唯一的ID信息
    IpAddress = "127.0.0.1",                                     // 服务器的地址
    Credentials = new MqttCredential( "admin", "123456" ),       // 设置了用户名和密码

    UseRSAProvider = true
} );
```

连接示例

[复制](# "复制")

```
// 连接服务器
HslCommunication.OperateResult connect =  mqttClient.ConnectServer( );
if (connect.IsSuccess)
{
    // 连接成功
}
else
{
    // 连接失败，过会就需要重新连接了
}

// 重点说明。如果过会网络不行了，断开了，内部会自动连接服务器的，你只管publish就可以了
```

发布示例

[复制](# "复制")

```
// 发布示例
HslCommunication.OperateResult connect = mqttClient.PublishMessage( new MqttApplicationMessage( )
{
    Topic = "A",                                                           // 主题
    QualityOfServiceLevel = MqttQualityOfServiceLevel.AtMostOnce,          // 如果是实时数据，适合用这个
    Payload = Encoding.UTF8.GetBytes("Test data")                          // 发布的数据
} );
if (connect.IsSuccess)
{
    // 发布成功
}
else
{
    // 发布失败
}

// 重点说明。如果过会网络不行了，断开了，内部会自动连接服务器的，你只管publish就可以了
```

订阅示例

[复制](# "复制")

```
// 订阅示例
mqttClient.OnMqttMessageReceived += ( MqttClient client, MqttApplicationMessage message ) =>
{
    Console.WriteLine( "Time:" + DateTime.Now.ToString( ) );
    Console.WriteLine( "Topic:" + message.Topic );
    Console.WriteLine( "Payload:" + Encoding.UTF8.GetString( message.Payload ) );
};

// 然后添加订阅
HslCommunication.OperateResult sub = mqttClient.SubscribeMessage( "A" );
if (sub.IsSuccess)
{
    // 订阅成功
}
else
{
    // 订阅失败
}
```

网络重连示例

[复制](# "复制")

```
// 网络失败的情况，需要自己来手动控制重连，实例化后进行事件绑定
mqttClient.OnNetworkError += ( object sender, EventArgs e ) =>
{
    // 当网络异常的时候触发，可以在此处重连服务器
    if (sender is MqttClient client)
    {
        // 开始重连服务器，直到连接成功为止
        client.LogNet?.WriteInfo( "网络异常，准备10秒后重新连接。" );
        while (true)
        {
            // 每隔10秒重连
            System.Threading.Thread.Sleep( 10_000 );
            client.LogNet?.WriteInfo( "准备重新连接服务器..." );
            HslCommunication.OperateResult connect = client.ConnectServer( );
            if (connect.IsSuccess)
            {
                // 连接成功后，可以在下方break之前进行订阅，或是数据初始化操作
                client.LogNet?.WriteInfo( "连接服务器成功！" );
                break;
            }
            client.LogNet?.WriteInfo( "连接失败，准备10秒后重新连接。" );
        }
    }
};
```

当我们在一个多窗体的客户端里使用了MqttClient类，可能很多界面都需要订阅主题，显示一些实时数据信息。只由主窗体来订阅再把数据传递给子窗体却不是很容易操作。
所以在hsl里提供了更加便捷的操作方法。方便在每个子窗体界面中，订阅，显示，取消订阅操作。核心代码如下：

子窗体订阅操作

[复制](# "复制")

```
// 如果有证书的情况
MqttClient mqttClient = new MqttClient( new MqttConnectionOptions( )
{
    ClientId = "ABC",                                            // 客户端的唯一的ID信息
    IpAddress = "127.0.0.1",                                     // 服务器的地址
    Credentials = new MqttCredential( "admin", "123456" ),       // 设置了用户名和密码

    CertificateFile = @"D:\emqx\etc\certs\cacert.pem",           // 指定证书的路径，这里选择的是emqx的测试证书
} );
```

以下的例子是DEMO程序的一个例子代码，也可以作为参考

DEMO子窗体

[复制](# "复制")

```
namespace HslCommunicationDemo.MQTT
{
    public partial class FormMqttSubscribe : System.Windows.Forms.Form
    {
        // 将 mqttClient 传递给子窗体
        public FormMqttSubscribe( MqttClient mqttClient )
        {
            InitializeComponent( );

            this.mqttClient = mqttClient;
        }

        private void FormMqttSubscribe_Load( object sender, EventArgs e )
        {
            button8.Enabled = false;
        }


        MqttClient mqttClient = null;

        private void button7_Click( object sender, EventArgs e )
        {
            // 子窗体订阅操作，在子窗体订阅的情况下，一般来说每个子窗体不同的topic主题
            OperateResult send = mqttClient.SubscribeMessage( new string[] { textBox5.Text } );

            if (!send.IsSuccess) DemoUtils.ShowMessage( "SubscribeMessage Failed:" + send.Message );
            else
            {
                // 获取订阅的信息，绑定本类的触发事件
                SubscribeTopic subscribeTopic = mqttClient.GetSubscribeTopic( textBox5.Text );
                if (subscribeTopic != null) subscribeTopic.OnMqttMessageReceived += FormMqttSubscribe_OnMqttMessageReceived;
                button7.Enabled = false;
                button8.Enabled = true;
            }
        }

        private void FormMqttSubscribe_OnMqttMessageReceived( MqttClient client, MqttApplicationMessage message )
        {
            // 订阅触发，这里举例是显示出来
            try
            {
                string topic = message.Topic;
                byte[] payload = message.Payload;
                Invoke( new Action( ( ) =>
                {
                    receiveCount++;
                    label10.Text = "Receive Count: " + receiveCount;
                    string msg = string.Empty;
                    if (radioButton_binary.Checked)
                    {
                        msg = payload.ToHexString( ' ' );
                    }
                    else if (radioButton_text.Checked)
                    {
                        msg = Encoding.UTF8.GetString( payload );
                    }
                    else if (radioButton_xml.Checked)
                    {
                        try
                        {
                            msg = XElement.Parse( Encoding.UTF8.GetString( payload ) ).ToString( );
                        }
                        catch
                        {
                            msg = Encoding.UTF8.GetString( payload );
                        }
                    }
                    else if (radioButton_json.Checked)
                    {
                        try
                        {
                            msg = Newtonsoft.Json.Linq.JObject.Parse( Encoding.UTF8.GetString( payload ) ).ToString( );
                        }
                        catch
                        {
                            msg = Encoding.UTF8.GetString( payload );
                        }
                    }

                    if (checkBox_long_message_hide.Checked)
                    {
                        if (msg?.Length > 200)
                        {
                            msg = msg.Substring( 200 );
                        }
                    }
                    if (radioButton2.Checked)
                        textBox8.AppendText( DateTime.Now.ToString( "yyyy-MM-dd HH:mm:ss.fff" ) + $" Topic[{topic}]: " + msg + Environment.NewLine );
                    else if (radioButton1.Checked)
                        textBox8.Text = DateTime.Now.ToString( "yyyy-MM-dd HH:mm:ss.fff" ) + $" Topic[{topic}]: " + msg;
                } ) );
            }
            catch
            {

            }
        }

        private long receiveCount = 0;


        private void button8_Click( object sender, EventArgs e )
        {
            // 取消订阅操作，需要先解绑注册事件
            SubscribeTopic subscribeTopic = mqttClient.GetSubscribeTopic( textBox5.Text );
            if (subscribeTopic != null)
            {
                subscribeTopic.OnMqttMessageReceived -= FormMqttSubscribe_OnMqttMessageReceived;

                // 取消订阅，为了不影响其他的界面的订阅信息（可能其他子界面也有订阅相同的主题），这里先判断订阅计数，减到 0 才真正的取消订阅
                // 如果你的主题只有一个窗体使用到的话，那么这里就不需要判断，直接取消订阅即可
                if (subscribeTopic.RemoveSubscribeTick( ) <= 0)
                {
                    OperateResult send = mqttClient.UnSubscribeMessage( textBox5.Text );

                    if (!send.IsSuccess) DemoUtils.ShowMessage( "UnSubscribeMessage Failed:" + send.Message );
                    else
                    {
                        button7.Enabled = true;
                        button8.Enabled = false;
                    }
                }
                else
                {
                    button7.Enabled = true;
                    button8.Enabled = false;
                }
            }
        }

        private void button4_Click( object sender, EventArgs e )
        {
            // 清空文本的操作
            textBox8.Clear( );
        }
    }
}
```

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.MQTT 命名空间](68b3b27c-e8db-3b64-4a05-22be97dc0b2b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## MqttClient 构造函数 

[原文連結](http://api.hslcommunication.cn/html/4817faf9-a999-ef74-9d36-d9efe7831b1b.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.MQTT](../html/68b3b27c-e8db-3b64-4a05-22be97dc0b2b.htm "HslCommunication.MQTT")

[MqttClient 类](../html/8626b3c4-5983-7308-94ba-65e3dd480b58.htm "MqttClient 类")

[MqttClient 构造函数](../html/4817faf9-a999-ef74-9d36-d9efe7831b1b.htm "MqttClient 构造函数 ")

[MqttClient 属性](../html/05c7efc3-bf81-a5ea-a5a0-b67d68df3e10.htm "MqttClient 属性")

[MqttClient 方法](../html/5745a0d8-2dd3-b35a-76f5-d0b6dc6fb012.htm "MqttClient 方法")

[MqttClient 事件](../html/22de311a-df03-7e25-27c9-2aa41f8c46df.htm "MqttClient 事件")

[MqttClient 字段](../html/384cf7bb-b210-9388-1d20-ef614282bfe0.htm "MqttClient 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MqttClient 构造函数 |

实例化一个默认的对象

**命名空间：**
 [HslCommunication.MQTT](68b3b27c-e8db-3b64-4a05-22be97dc0b2b.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public MqttClient(
	MqttConnectionOptions options
)
```

```
Public Sub New ( 
	options As MqttConnectionOptions
)
```

```
public:
MqttClient(
	MqttConnectionOptions^ options
)
```

```
new : 
        options : MqttConnectionOptions -> MqttClient
```

#### 参数

options
:   类型：[HslCommunication.MQTTMqttConnectionOptions](ba13f335-1fd8-1492-9d41-9f649c890f04.htm)  
    配置信息

![](../icons/SectionExpanded.png)参见

#### 引用

[MqttClient 类](8626b3c4-5983-7308-94ba-65e3dd480b58.htm)

[HslCommunication.MQTT 命名空间](68b3b27c-e8db-3b64-4a05-22be97dc0b2b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## MqttClient 属性

[原文連結](http://api.hslcommunication.cn/html/05c7efc3-bf81-a5ea-a5a0-b67d68df3e10.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.MQTT](../html/68b3b27c-e8db-3b64-4a05-22be97dc0b2b.htm "HslCommunication.MQTT")

[MqttClient 类](../html/8626b3c4-5983-7308-94ba-65e3dd480b58.htm "MqttClient 类")

[MqttClient 属性](../html/05c7efc3-bf81-a5ea-a5a0-b67d68df3e10.htm "MqttClient 属性")

[ConnectionOptions 属性](../html/59421305-6be2-389d-e9e6-093321835fac.htm "ConnectionOptions 属性 ")

[IsConnected 属性](../html/91d587d0-b6d8-96e7-0e8d-b73360d6820f.htm "IsConnected 属性 ")

[SubcribeTopics 属性](../html/8b884727-7781-3342-cd59-b5d4af0f2fac.htm "SubcribeTopics 属性 ")

[Tag 属性](../html/7d226192-f98d-e251-5fd2-3a664d0933de.htm "Tag 属性 ")

[UseTimerCheckDropped 属性](../html/3b46c470-9ee6-79b9-f84d-f40a746ef99b.htm "UseTimerCheckDropped 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MqttClient 属性 |

[MqttClient](8626b3c4-5983-7308-94ba-65e3dd480b58.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [ConnectionOptions](59421305-6be2-389d-e9e6-093321835fac.htm) | 获取当前的连接配置参数信息  Get current connection configuration parameter information |
| 公共属性 | [IsConnected](91d587d0-b6d8-96e7-0e8d-b73360d6820f.htm) | 获取或设置当前的服务器连接是否成功，定时获取本属性可用于实时更新连接状态信息。  Get or set whether the current server connection is successful or not. This property can be obtained regularly and can be used to update the connection status information in real time. |
| 公共属性代码示例 | [LogNet](46dfe698-46c6-73a2-5458-8fd568d264ce.htm) | 组件的日志工具，支持日志记录，只要实例化后，当前网络的基本信息，就以[DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm)等级进行输出  The component's logging tool supports logging. As long as the instantiation of the basic network information, the output will be output at [DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm) (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 公共属性 | [SubcribeTopics](8b884727-7781-3342-cd59-b5d4af0f2fac.htm) | 获取当前的客户端对象已经订阅的所有的Topic信息  Get all Topic information that the current client object has subscribed to |
| 公共属性 | [Tag](7d226192-f98d-e251-5fd2-3a664d0933de.htm) | 获取或设置当前客户端关联的自定义的对象内容  Gets or sets the custom object content associated with the current client |
| 公共属性代码示例 | [Token](80db00a7-aaa5-164f-e87c-af5400dc80fa.htm) | 网络类的身份令牌，在hsl协议的模式下会有效，在和设备进行通信的时候是无效的  Network-type identity tokens will be valid in the hsl protocol mode and will not be valid when communicating with the device (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 公共属性 | [UseTimerCheckDropped](3b46c470-9ee6-79b9-f84d-f40a746ef99b.htm) | 获取或设置是否启动定时器去检测当前客户端是否超时掉线。默认为 True  Get or set whether to start the timer to detect whether the current client timeout and disconnection. Default is True |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[MqttClient 类](8626b3c4-5983-7308-94ba-65e3dd480b58.htm)

[HslCommunication.MQTT 命名空间](68b3b27c-e8db-3b64-4a05-22be97dc0b2b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ConnectionOptions 属性 

[原文連結](http://api.hslcommunication.cn/html/59421305-6be2-389d-e9e6-093321835fac.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.MQTT](../html/68b3b27c-e8db-3b64-4a05-22be97dc0b2b.htm "HslCommunication.MQTT")

[MqttClient 类](../html/8626b3c4-5983-7308-94ba-65e3dd480b58.htm "MqttClient 类")

[MqttClient 属性](../html/05c7efc3-bf81-a5ea-a5a0-b67d68df3e10.htm "MqttClient 属性")

[ConnectionOptions 属性](../html/59421305-6be2-389d-e9e6-093321835fac.htm "ConnectionOptions 属性 ")

[IsConnected 属性](../html/91d587d0-b6d8-96e7-0e8d-b73360d6820f.htm "IsConnected 属性 ")

[SubcribeTopics 属性](../html/8b884727-7781-3342-cd59-b5d4af0f2fac.htm "SubcribeTopics 属性 ")

[Tag 属性](../html/7d226192-f98d-e251-5fd2-3a664d0933de.htm "Tag 属性 ")

[UseTimerCheckDropped 属性](../html/3b46c470-9ee6-79b9-f84d-f40a746ef99b.htm "UseTimerCheckDropped 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MqttClientConnectionOptions 属性 |

获取当前的连接配置参数信息  
Get current connection configuration parameter information

**命名空间：**
 [HslCommunication.MQTT](68b3b27c-e8db-3b64-4a05-22be97dc0b2b.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public MqttConnectionOptions ConnectionOptions { get; }
```

```
Public ReadOnly Property ConnectionOptions As MqttConnectionOptions
	Get
```

```
public:
property MqttConnectionOptions^ ConnectionOptions {
	MqttConnectionOptions^ get ();
}
```

```
member ConnectionOptions : MqttConnectionOptions with get
```

#### 属性值

类型：[MqttConnectionOptions](ba13f335-1fd8-1492-9d41-9f649c890f04.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[MqttClient 类](8626b3c4-5983-7308-94ba-65e3dd480b58.htm)

[HslCommunication.MQTT 命名空间](68b3b27c-e8db-3b64-4a05-22be97dc0b2b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## IsConnected 属性 

[原文連結](http://api.hslcommunication.cn/html/91d587d0-b6d8-96e7-0e8d-b73360d6820f.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.MQTT](../html/68b3b27c-e8db-3b64-4a05-22be97dc0b2b.htm "HslCommunication.MQTT")

[MqttClient 类](../html/8626b3c4-5983-7308-94ba-65e3dd480b58.htm "MqttClient 类")

[MqttClient 属性](../html/05c7efc3-bf81-a5ea-a5a0-b67d68df3e10.htm "MqttClient 属性")

[ConnectionOptions 属性](../html/59421305-6be2-389d-e9e6-093321835fac.htm "ConnectionOptions 属性 ")

[IsConnected 属性](../html/91d587d0-b6d8-96e7-0e8d-b73360d6820f.htm "IsConnected 属性 ")

[SubcribeTopics 属性](../html/8b884727-7781-3342-cd59-b5d4af0f2fac.htm "SubcribeTopics 属性 ")

[Tag 属性](../html/7d226192-f98d-e251-5fd2-3a664d0933de.htm "Tag 属性 ")

[UseTimerCheckDropped 属性](../html/3b46c470-9ee6-79b9-f84d-f40a746ef99b.htm "UseTimerCheckDropped 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MqttClientIsConnected 属性 |

获取或设置当前的服务器连接是否成功，定时获取本属性可用于实时更新连接状态信息。  
Get or set whether the current server connection is successful or not.
This property can be obtained regularly and can be used to update the connection status information in real time.

**命名空间：**
 [HslCommunication.MQTT](68b3b27c-e8db-3b64-4a05-22be97dc0b2b.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public bool IsConnected { get; }
```

```
Public ReadOnly Property IsConnected As Boolean
	Get
```

```
public:
property bool IsConnected {
	bool get ();
}
```

```
member IsConnected : bool with get
```

#### 属性值

类型：Boolean

![](../icons/SectionExpanded.png)参见

#### 引用

[MqttClient 类](8626b3c4-5983-7308-94ba-65e3dd480b58.htm)

[HslCommunication.MQTT 命名空间](68b3b27c-e8db-3b64-4a05-22be97dc0b2b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## SubcribeTopics 属性 

[原文連結](http://api.hslcommunication.cn/html/8b884727-7781-3342-cd59-b5d4af0f2fac.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.MQTT](../html/68b3b27c-e8db-3b64-4a05-22be97dc0b2b.htm "HslCommunication.MQTT")

[MqttClient 类](../html/8626b3c4-5983-7308-94ba-65e3dd480b58.htm "MqttClient 类")

[MqttClient 属性](../html/05c7efc3-bf81-a5ea-a5a0-b67d68df3e10.htm "MqttClient 属性")

[ConnectionOptions 属性](../html/59421305-6be2-389d-e9e6-093321835fac.htm "ConnectionOptions 属性 ")

[IsConnected 属性](../html/91d587d0-b6d8-96e7-0e8d-b73360d6820f.htm "IsConnected 属性 ")

[SubcribeTopics 属性](../html/8b884727-7781-3342-cd59-b5d4af0f2fac.htm "SubcribeTopics 属性 ")

[Tag 属性](../html/7d226192-f98d-e251-5fd2-3a664d0933de.htm "Tag 属性 ")

[UseTimerCheckDropped 属性](../html/3b46c470-9ee6-79b9-f84d-f40a746ef99b.htm "UseTimerCheckDropped 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MqttClientSubcribeTopics 属性 |

获取当前的客户端对象已经订阅的所有的Topic信息  
Get all Topic information that the current client object has subscribed to

**命名空间：**
 [HslCommunication.MQTT](68b3b27c-e8db-3b64-4a05-22be97dc0b2b.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public string[] SubcribeTopics { get; }
```

```
Public ReadOnly Property SubcribeTopics As String()
	Get
```

```
public:
property array<String^>^ SubcribeTopics {
	array<String^>^ get ();
}
```

```
member SubcribeTopics : string[] with get
```

#### 属性值

类型：String

![](../icons/SectionExpanded.png)参见

#### 引用

[MqttClient 类](8626b3c4-5983-7308-94ba-65e3dd480b58.htm)

[HslCommunication.MQTT 命名空间](68b3b27c-e8db-3b64-4a05-22be97dc0b2b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Tag 属性 

[原文連結](http://api.hslcommunication.cn/html/7d226192-f98d-e251-5fd2-3a664d0933de.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.MQTT](../html/68b3b27c-e8db-3b64-4a05-22be97dc0b2b.htm "HslCommunication.MQTT")

[MqttClient 类](../html/8626b3c4-5983-7308-94ba-65e3dd480b58.htm "MqttClient 类")

[MqttClient 属性](../html/05c7efc3-bf81-a5ea-a5a0-b67d68df3e10.htm "MqttClient 属性")

[ConnectionOptions 属性](../html/59421305-6be2-389d-e9e6-093321835fac.htm "ConnectionOptions 属性 ")

[IsConnected 属性](../html/91d587d0-b6d8-96e7-0e8d-b73360d6820f.htm "IsConnected 属性 ")

[SubcribeTopics 属性](../html/8b884727-7781-3342-cd59-b5d4af0f2fac.htm "SubcribeTopics 属性 ")

[Tag 属性](../html/7d226192-f98d-e251-5fd2-3a664d0933de.htm "Tag 属性 ")

[UseTimerCheckDropped 属性](../html/3b46c470-9ee6-79b9-f84d-f40a746ef99b.htm "UseTimerCheckDropped 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MqttClientTag 属性 |

获取或设置当前客户端关联的自定义的对象内容  
Gets or sets the custom object content associated with the current client

**命名空间：**
 [HslCommunication.MQTT](68b3b27c-e8db-3b64-4a05-22be97dc0b2b.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public Object Tag { get; set; }
```

```
Public Property Tag As Object
	Get
	Set
```

```
public:
property Object^ Tag {
	Object^ get ();
	void set (Object^ value);
}
```

```
member Tag : Object with get, set
```

#### 属性值

类型：Object

![](../icons/SectionExpanded.png)参见

#### 引用

[MqttClient 类](8626b3c4-5983-7308-94ba-65e3dd480b58.htm)

[HslCommunication.MQTT 命名空间](68b3b27c-e8db-3b64-4a05-22be97dc0b2b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## UseTimerCheckDropped 属性 

[原文連結](http://api.hslcommunication.cn/html/3b46c470-9ee6-79b9-f84d-f40a746ef99b.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.MQTT](../html/68b3b27c-e8db-3b64-4a05-22be97dc0b2b.htm "HslCommunication.MQTT")

[MqttClient 类](../html/8626b3c4-5983-7308-94ba-65e3dd480b58.htm "MqttClient 类")

[MqttClient 属性](../html/05c7efc3-bf81-a5ea-a5a0-b67d68df3e10.htm "MqttClient 属性")

[ConnectionOptions 属性](../html/59421305-6be2-389d-e9e6-093321835fac.htm "ConnectionOptions 属性 ")

[IsConnected 属性](../html/91d587d0-b6d8-96e7-0e8d-b73360d6820f.htm "IsConnected 属性 ")

[SubcribeTopics 属性](../html/8b884727-7781-3342-cd59-b5d4af0f2fac.htm "SubcribeTopics 属性 ")

[Tag 属性](../html/7d226192-f98d-e251-5fd2-3a664d0933de.htm "Tag 属性 ")

[UseTimerCheckDropped 属性](../html/3b46c470-9ee6-79b9-f84d-f40a746ef99b.htm "UseTimerCheckDropped 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MqttClientUseTimerCheckDropped 属性 |

获取或设置是否启动定时器去检测当前客户端是否超时掉线。默认为 True  
Get or set whether to start the timer to detect whether the current client timeout and disconnection. Default is True

**命名空间：**
 [HslCommunication.MQTT](68b3b27c-e8db-3b64-4a05-22be97dc0b2b.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public bool UseTimerCheckDropped { get; set; }
```

```
Public Property UseTimerCheckDropped As Boolean
	Get
	Set
```

```
public:
property bool UseTimerCheckDropped {
	bool get ();
	void set (bool value);
}
```

```
member UseTimerCheckDropped : bool with get, set
```

#### 属性值

类型：Boolean

![](../icons/SectionExpanded.png)参见

#### 引用

[MqttClient 类](8626b3c4-5983-7308-94ba-65e3dd480b58.htm)

[HslCommunication.MQTT 命名空间](68b3b27c-e8db-3b64-4a05-22be97dc0b2b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)