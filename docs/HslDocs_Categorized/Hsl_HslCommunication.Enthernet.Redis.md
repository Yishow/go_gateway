# HslCommunication - HslCommunication.Enthernet.Redis

> 分類頁數: 30



---
## HslCommunication.Enthernet.Redis

[原文連結](http://api.hslcommunication.cn/html/cc44ce00-3781-c690-b233-24000870bf42.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Enthernet.Redis](../html/cc44ce00-3781-c690-b233-24000870bf42.htm "HslCommunication.Enthernet.Redis")

[IRedisConnector 类](../html/e1de67c7-4f90-73e3-2b89-9c49b48edf71.htm "IRedisConnector 类")

[RedisClient 类](../html/e33d9088-6df0-32e9-33f1-426e944c3b7f.htm "RedisClient 类")

[RedisClient.RedisMessageReceiveDelegate 委托](../html/76cad875-7768-57e1-fa76-495febf93410.htm "RedisClient.RedisMessageReceiveDelegate 委托")

[RedisClientPool 类](../html/116fac7a-cbca-c7b8-0608-1d7c0009d11c.htm "RedisClientPool 类")

[RedisHelper 类](../html/ff485d12-efb9-ba99-f502-24b87fb4c8be.htm "RedisHelper 类")

[RedisSubscribe 类](../html/9a7bbb54-a4ea-6051-cf31-9deeedb15d3f.htm "RedisSubscribe 类")

[RedisSubscribe.RedisMessageReceiveDelegate 委托](../html/970d78e2-5cd8-b8bf-084b-df92bc8e9c5f.htm "RedisSubscribe.RedisMessageReceiveDelegate 委托")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslCommunication.Enthernet.Redis 命名空间 |

[缺少 "N:HslCommunication.Enthernet.Redis" 的 <summary> 文档]

![](../icons/SectionExpanded.png)类

|  | 类 | 说明 |
| --- | --- | --- |
| 公共类 | [IRedisConnector](e1de67c7-4f90-73e3-2b89-9c49b48edf71.htm) | 关于Redis实现的接口[IConnector](ede7f47f-b60e-e320-8bc2-389883e77083.htm)，从而实现了数据连接池的操作信息 |
| 公共类代码示例 | [RedisClient](e33d9088-6df0-32e9-33f1-426e944c3b7f.htm) | 这是一个redis的客户端类，支持读取，写入，发布订阅，但是不支持订阅，如果需要订阅，请使用另一个类[RedisSubscribe](9a7bbb54-a4ea-6051-cf31-9deeedb15d3f.htm) |
| 公共类 | [RedisClientPool](116fac7a-cbca-c7b8-0608-1d7c0009d11c.htm) | **[商业授权]** Redis客户端的连接池类对象，用于共享当前的连接池，合理的动态调整连接对象，然后进行高效通信的操作，默认连接数无限大。 **[Authorization]** The connection pool class object of the Redis client is used to share the current connection pool, reasonably dynamically adjust the connection object, and then perform efficient communication operations, The default number of connections is unlimited |
| 公共类 | [RedisHelper](ff485d12-efb9-ba99-f502-24b87fb4c8be.htm) | 提供了redis辅助类的一些方法 |
| 公共类 | [RedisSubscribe](9a7bbb54-a4ea-6051-cf31-9deeedb15d3f.htm) | Redis协议的订阅操作，一个对象订阅一个或是多个频道的信息，当发生网络异常的时候，内部会进行自动重连，并恢复之前的订阅信息。  In the subscription operation of the Redis protocol, an object subscribes to the information of one or more channels. When a network abnormality occurs, the internal will automatically reconnect and restore the previous subscription information. |

![](../icons/SectionExpanded.png)委托

|  | 委托 | 说明 |
| --- | --- | --- |
| 公共委托 | [RedisClientRedisMessageReceiveDelegate](76cad875-7768-57e1-fa76-495febf93410.htm) | 当接收到Redis订阅的信息的时候触发  Triggered when receiving Redis subscription information |
| 公共委托 | [RedisSubscribeRedisMessageReceiveDelegate](970d78e2-5cd8-b8bf-084b-df92bc8e9c5f.htm) | 当接收到Redis订阅的信息的时候触发  Triggered when receiving Redis subscription information |

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## IRedisConnector 类

[原文連結](http://api.hslcommunication.cn/html/e1de67c7-4f90-73e3-2b89-9c49b48edf71.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Enthernet.Redis](../html/cc44ce00-3781-c690-b233-24000870bf42.htm "HslCommunication.Enthernet.Redis")

[IRedisConnector 类](../html/e1de67c7-4f90-73e3-2b89-9c49b48edf71.htm "IRedisConnector 类")

[IRedisConnector 构造函数](../html/245cc13d-51e0-a1d8-05d0-c16c27f19ddc.htm "IRedisConnector 构造函数 ")

[IRedisConnector 属性](../html/8c0eeaf3-5588-f642-27af-c72afe9902be.htm "IRedisConnector 属性")

[IRedisConnector 方法](../html/f477e470-60ae-13e1-aa3b-87cc8b7002fa.htm "IRedisConnector 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IRedisConnector 类 |

关于Redis实现的接口[IConnector](ede7f47f-b60e-e320-8bc2-389883e77083.htm)，从而实现了数据连接池的操作信息

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  HslCommunication.Enthernet.RedisIRedisConnector

**命名空间：**
 [HslCommunication.Enthernet.Redis](cc44ce00-3781-c690-b233-24000870bf42.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class IRedisConnector : IConnector
```

```
Public Class IRedisConnector
	Implements IConnector
```

```
public ref class IRedisConnector : IConnector
```

```
type IRedisConnector =  
    class
        interface IConnector
    end
```

IRedisConnector 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [IRedisConnector](245cc13d-51e0-a1d8-05d0-c16c27f19ddc.htm) | 初始化 IRedisConnector 类的一个新实例 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [GuidToken](30b26e60-232f-44a8-2ebc-698983f48a9c.htm) | 唯一的GUID码 |
| 公共属性 | [IsConnectUsing](5f2162de-1fc7-1100-055a-a6c2e6baf5a2.htm) | 指示当前的连接是否在使用用 |
| 公共属性 | [LastUseTime](c7216bb1-d39c-f706-7882-ec4a94e13ca6.htm) | 最新一次使用的时间 |
| 公共属性 | [Redis](11b0612b-67af-bd55-b73d-d2361720cf6f.htm) | Redis的连接对象 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [Close](446fd186-a993-884b-125b-5fdbaefd95bb.htm) | 关闭并释放 |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法 | [Open](e690ef36-e1a0-ca3b-9ab5-f4127c15742b.htm) | 打开连接 |
| 公共方法 | ToString | (继承自 Object。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.Enthernet.Redis 命名空间](cc44ce00-3781-c690-b233-24000870bf42.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## IRedisConnector 构造函数 

[原文連結](http://api.hslcommunication.cn/html/245cc13d-51e0-a1d8-05d0-c16c27f19ddc.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Enthernet.Redis](../html/cc44ce00-3781-c690-b233-24000870bf42.htm "HslCommunication.Enthernet.Redis")

[IRedisConnector 类](../html/e1de67c7-4f90-73e3-2b89-9c49b48edf71.htm "IRedisConnector 类")

[IRedisConnector 构造函数](../html/245cc13d-51e0-a1d8-05d0-c16c27f19ddc.htm "IRedisConnector 构造函数 ")

[IRedisConnector 属性](../html/8c0eeaf3-5588-f642-27af-c72afe9902be.htm "IRedisConnector 属性")

[IRedisConnector 方法](../html/f477e470-60ae-13e1-aa3b-87cc8b7002fa.htm "IRedisConnector 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IRedisConnector 构造函数 |

初始化 [IRedisConnector](e1de67c7-4f90-73e3-2b89-9c49b48edf71.htm) 类的一个新实例

**命名空间：**
 [HslCommunication.Enthernet.Redis](cc44ce00-3781-c690-b233-24000870bf42.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public IRedisConnector()
```

```
Public Sub New
```

```
public:
IRedisConnector()
```

```
new : unit -> IRedisConnector
```

![](../icons/SectionExpanded.png)参见

#### 引用

[IRedisConnector 类](e1de67c7-4f90-73e3-2b89-9c49b48edf71.htm)

[HslCommunication.Enthernet.Redis 命名空间](cc44ce00-3781-c690-b233-24000870bf42.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## IRedisConnector 属性

[原文連結](http://api.hslcommunication.cn/html/8c0eeaf3-5588-f642-27af-c72afe9902be.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Enthernet.Redis](../html/cc44ce00-3781-c690-b233-24000870bf42.htm "HslCommunication.Enthernet.Redis")

[IRedisConnector 类](../html/e1de67c7-4f90-73e3-2b89-9c49b48edf71.htm "IRedisConnector 类")

[IRedisConnector 属性](../html/8c0eeaf3-5588-f642-27af-c72afe9902be.htm "IRedisConnector 属性")

[GuidToken 属性](../html/30b26e60-232f-44a8-2ebc-698983f48a9c.htm "GuidToken 属性 ")

[IsConnectUsing 属性](../html/5f2162de-1fc7-1100-055a-a6c2e6baf5a2.htm "IsConnectUsing 属性 ")

[LastUseTime 属性](../html/c7216bb1-d39c-f706-7882-ec4a94e13ca6.htm "LastUseTime 属性 ")

[Redis 属性](../html/11b0612b-67af-bd55-b73d-d2361720cf6f.htm "Redis 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IRedisConnector 属性 |

[IRedisConnector](e1de67c7-4f90-73e3-2b89-9c49b48edf71.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [GuidToken](30b26e60-232f-44a8-2ebc-698983f48a9c.htm) | 唯一的GUID码 |
| 公共属性 | [IsConnectUsing](5f2162de-1fc7-1100-055a-a6c2e6baf5a2.htm) | 指示当前的连接是否在使用用 |
| 公共属性 | [LastUseTime](c7216bb1-d39c-f706-7882-ec4a94e13ca6.htm) | 最新一次使用的时间 |
| 公共属性 | [Redis](11b0612b-67af-bd55-b73d-d2361720cf6f.htm) | Redis的连接对象 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[IRedisConnector 类](e1de67c7-4f90-73e3-2b89-9c49b48edf71.htm)

[HslCommunication.Enthernet.Redis 命名空间](cc44ce00-3781-c690-b233-24000870bf42.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GuidToken 属性 

[原文連結](http://api.hslcommunication.cn/html/30b26e60-232f-44a8-2ebc-698983f48a9c.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Enthernet.Redis](../html/cc44ce00-3781-c690-b233-24000870bf42.htm "HslCommunication.Enthernet.Redis")

[IRedisConnector 类](../html/e1de67c7-4f90-73e3-2b89-9c49b48edf71.htm "IRedisConnector 类")

[IRedisConnector 属性](../html/8c0eeaf3-5588-f642-27af-c72afe9902be.htm "IRedisConnector 属性")

[GuidToken 属性](../html/30b26e60-232f-44a8-2ebc-698983f48a9c.htm "GuidToken 属性 ")

[IsConnectUsing 属性](../html/5f2162de-1fc7-1100-055a-a6c2e6baf5a2.htm "IsConnectUsing 属性 ")

[LastUseTime 属性](../html/c7216bb1-d39c-f706-7882-ec4a94e13ca6.htm "LastUseTime 属性 ")

[Redis 属性](../html/11b0612b-67af-bd55-b73d-d2361720cf6f.htm "Redis 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IRedisConnectorGuidToken 属性 |

唯一的GUID码

**命名空间：**
 [HslCommunication.Enthernet.Redis](cc44ce00-3781-c690-b233-24000870bf42.htm)  
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

[IRedisConnector 类](e1de67c7-4f90-73e3-2b89-9c49b48edf71.htm)

[HslCommunication.Enthernet.Redis 命名空间](cc44ce00-3781-c690-b233-24000870bf42.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## IsConnectUsing 属性 

[原文連結](http://api.hslcommunication.cn/html/5f2162de-1fc7-1100-055a-a6c2e6baf5a2.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Enthernet.Redis](../html/cc44ce00-3781-c690-b233-24000870bf42.htm "HslCommunication.Enthernet.Redis")

[IRedisConnector 类](../html/e1de67c7-4f90-73e3-2b89-9c49b48edf71.htm "IRedisConnector 类")

[IRedisConnector 属性](../html/8c0eeaf3-5588-f642-27af-c72afe9902be.htm "IRedisConnector 属性")

[GuidToken 属性](../html/30b26e60-232f-44a8-2ebc-698983f48a9c.htm "GuidToken 属性 ")

[IsConnectUsing 属性](../html/5f2162de-1fc7-1100-055a-a6c2e6baf5a2.htm "IsConnectUsing 属性 ")

[LastUseTime 属性](../html/c7216bb1-d39c-f706-7882-ec4a94e13ca6.htm "LastUseTime 属性 ")

[Redis 属性](../html/11b0612b-67af-bd55-b73d-d2361720cf6f.htm "Redis 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IRedisConnectorIsConnectUsing 属性 |

指示当前的连接是否在使用用

**命名空间：**
 [HslCommunication.Enthernet.Redis](cc44ce00-3781-c690-b233-24000870bf42.htm)  
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

[IRedisConnector 类](e1de67c7-4f90-73e3-2b89-9c49b48edf71.htm)

[HslCommunication.Enthernet.Redis 命名空间](cc44ce00-3781-c690-b233-24000870bf42.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## LastUseTime 属性 

[原文連結](http://api.hslcommunication.cn/html/c7216bb1-d39c-f706-7882-ec4a94e13ca6.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Enthernet.Redis](../html/cc44ce00-3781-c690-b233-24000870bf42.htm "HslCommunication.Enthernet.Redis")

[IRedisConnector 类](../html/e1de67c7-4f90-73e3-2b89-9c49b48edf71.htm "IRedisConnector 类")

[IRedisConnector 属性](../html/8c0eeaf3-5588-f642-27af-c72afe9902be.htm "IRedisConnector 属性")

[GuidToken 属性](../html/30b26e60-232f-44a8-2ebc-698983f48a9c.htm "GuidToken 属性 ")

[IsConnectUsing 属性](../html/5f2162de-1fc7-1100-055a-a6c2e6baf5a2.htm "IsConnectUsing 属性 ")

[LastUseTime 属性](../html/c7216bb1-d39c-f706-7882-ec4a94e13ca6.htm "LastUseTime 属性 ")

[Redis 属性](../html/11b0612b-67af-bd55-b73d-d2361720cf6f.htm "Redis 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IRedisConnectorLastUseTime 属性 |

最新一次使用的时间

**命名空间：**
 [HslCommunication.Enthernet.Redis](cc44ce00-3781-c690-b233-24000870bf42.htm)  
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

[IRedisConnector 类](e1de67c7-4f90-73e3-2b89-9c49b48edf71.htm)

[HslCommunication.Enthernet.Redis 命名空间](cc44ce00-3781-c690-b233-24000870bf42.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Redis 属性 

[原文連結](http://api.hslcommunication.cn/html/11b0612b-67af-bd55-b73d-d2361720cf6f.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Enthernet.Redis](../html/cc44ce00-3781-c690-b233-24000870bf42.htm "HslCommunication.Enthernet.Redis")

[IRedisConnector 类](../html/e1de67c7-4f90-73e3-2b89-9c49b48edf71.htm "IRedisConnector 类")

[IRedisConnector 属性](../html/8c0eeaf3-5588-f642-27af-c72afe9902be.htm "IRedisConnector 属性")

[GuidToken 属性](../html/30b26e60-232f-44a8-2ebc-698983f48a9c.htm "GuidToken 属性 ")

[IsConnectUsing 属性](../html/5f2162de-1fc7-1100-055a-a6c2e6baf5a2.htm "IsConnectUsing 属性 ")

[LastUseTime 属性](../html/c7216bb1-d39c-f706-7882-ec4a94e13ca6.htm "LastUseTime 属性 ")

[Redis 属性](../html/11b0612b-67af-bd55-b73d-d2361720cf6f.htm "Redis 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IRedisConnectorRedis 属性 |

Redis的连接对象

**命名空间：**
 [HslCommunication.Enthernet.Redis](cc44ce00-3781-c690-b233-24000870bf42.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public RedisClient Redis { get; set; }
```

```
Public Property Redis As RedisClient
	Get
	Set
```

```
public:
property RedisClient^ Redis {
	RedisClient^ get ();
	void set (RedisClient^ value);
}
```

```
member Redis : RedisClient with get, set
```

#### 属性值

类型：[RedisClient](e33d9088-6df0-32e9-33f1-426e944c3b7f.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[IRedisConnector 类](e1de67c7-4f90-73e3-2b89-9c49b48edf71.htm)

[HslCommunication.Enthernet.Redis 命名空间](cc44ce00-3781-c690-b233-24000870bf42.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## IRedisConnector 方法

[原文連結](http://api.hslcommunication.cn/html/f477e470-60ae-13e1-aa3b-87cc8b7002fa.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Enthernet.Redis](../html/cc44ce00-3781-c690-b233-24000870bf42.htm "HslCommunication.Enthernet.Redis")

[IRedisConnector 类](../html/e1de67c7-4f90-73e3-2b89-9c49b48edf71.htm "IRedisConnector 类")

[IRedisConnector 方法](../html/f477e470-60ae-13e1-aa3b-87cc8b7002fa.htm "IRedisConnector 方法")

[Close 方法](../html/446fd186-a993-884b-125b-5fdbaefd95bb.htm "Close 方法 ")

[Open 方法](../html/e690ef36-e1a0-ca3b-9ab5-f4127c15742b.htm "Open 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IRedisConnector 方法 |

[IRedisConnector](e1de67c7-4f90-73e3-2b89-9c49b48edf71.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [Close](446fd186-a993-884b-125b-5fdbaefd95bb.htm) | 关闭并释放 |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法 | [Open](e690ef36-e1a0-ca3b-9ab5-f4127c15742b.htm) | 打开连接 |
| 公共方法 | ToString | (继承自 Object。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[IRedisConnector 类](e1de67c7-4f90-73e3-2b89-9c49b48edf71.htm)

[HslCommunication.Enthernet.Redis 命名空间](cc44ce00-3781-c690-b233-24000870bf42.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Close 方法 

[原文連結](http://api.hslcommunication.cn/html/446fd186-a993-884b-125b-5fdbaefd95bb.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Enthernet.Redis](../html/cc44ce00-3781-c690-b233-24000870bf42.htm "HslCommunication.Enthernet.Redis")

[IRedisConnector 类](../html/e1de67c7-4f90-73e3-2b89-9c49b48edf71.htm "IRedisConnector 类")

[IRedisConnector 方法](../html/f477e470-60ae-13e1-aa3b-87cc8b7002fa.htm "IRedisConnector 方法")

[Close 方法](../html/446fd186-a993-884b-125b-5fdbaefd95bb.htm "Close 方法 ")

[Open 方法](../html/e690ef36-e1a0-ca3b-9ab5-f4127c15742b.htm "Open 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IRedisConnectorClose 方法 |

关闭并释放

**命名空间：**
 [HslCommunication.Enthernet.Redis](cc44ce00-3781-c690-b233-24000870bf42.htm)  
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

[IRedisConnector 类](e1de67c7-4f90-73e3-2b89-9c49b48edf71.htm)

[HslCommunication.Enthernet.Redis 命名空间](cc44ce00-3781-c690-b233-24000870bf42.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Open 方法 

[原文連結](http://api.hslcommunication.cn/html/e690ef36-e1a0-ca3b-9ab5-f4127c15742b.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Enthernet.Redis](../html/cc44ce00-3781-c690-b233-24000870bf42.htm "HslCommunication.Enthernet.Redis")

[IRedisConnector 类](../html/e1de67c7-4f90-73e3-2b89-9c49b48edf71.htm "IRedisConnector 类")

[IRedisConnector 方法](../html/f477e470-60ae-13e1-aa3b-87cc8b7002fa.htm "IRedisConnector 方法")

[Close 方法](../html/446fd186-a993-884b-125b-5fdbaefd95bb.htm "Close 方法 ")

[Open 方法](../html/e690ef36-e1a0-ca3b-9ab5-f4127c15742b.htm "Open 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IRedisConnectorOpen 方法 |

打开连接

**命名空间：**
 [HslCommunication.Enthernet.Redis](cc44ce00-3781-c690-b233-24000870bf42.htm)  
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

[IRedisConnector 类](e1de67c7-4f90-73e3-2b89-9c49b48edf71.htm)

[HslCommunication.Enthernet.Redis 命名空间](cc44ce00-3781-c690-b233-24000870bf42.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## RedisClient 类

[原文連結](http://api.hslcommunication.cn/html/e33d9088-6df0-32e9-33f1-426e944c3b7f.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Enthernet.Redis](../html/cc44ce00-3781-c690-b233-24000870bf42.htm "HslCommunication.Enthernet.Redis")

[RedisClient 类](../html/e33d9088-6df0-32e9-33f1-426e944c3b7f.htm "RedisClient 类")

[RedisClient 构造函数](../html/873ec4ce-423e-d519-ac3e-894ca64609ba.htm "RedisClient 构造函数 ")

[RedisClient 属性](../html/5708cfe8-172e-431e-d36e-02d5e23288e1.htm "RedisClient 属性")

[RedisClient 方法](../html/e41be4a5-1490-ffbd-f241-938f2d6c3f07.htm "RedisClient 方法")

[RedisClient 事件](../html/931bae45-0bc1-0713-e984-2d589dd7c843.htm "RedisClient 事件")

[RedisClient 字段](../html/4a43acdf-f8bd-333b-da81-99a1d696e961.htm "RedisClient 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| RedisClient 类 |

这是一个redis的客户端类，支持读取，写入，发布订阅，但是不支持订阅，如果需要订阅，请使用另一个类[RedisSubscribe](9a7bbb54-a4ea-6051-cf31-9deeedb15d3f.htm)

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  [HslCommunication.Core.NetNetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)  
    [HslCommunication.Core.NetNetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)  
      HslCommunication.Enthernet.RedisRedisClient

**命名空间：**
 [HslCommunication.Enthernet.Redis](cc44ce00-3781-c690-b233-24000870bf42.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class RedisClient : NetworkDoubleBase
```

```
Public Class RedisClient
	Inherits NetworkDoubleBase
```

```
public ref class RedisClient : public NetworkDoubleBase
```

```
type RedisClient =  
    class
        inherit NetworkDoubleBase
    end
```

RedisClient 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [RedisClient(String)](5407e3c7-7bdb-d5bc-5556-e51e73c318d0.htm) | 实例化一个客户端对象，需要手动指定Ip地址和端口 |
| 公共方法 | [RedisClient(String, Int32, String)](3ba84c15-5c80-601c-6d24-fc66a6d1eba1.htm) | 实例化一个客户端的对象，用于和服务器通信 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [AlienSession](5e30a9ab-7102-79a2-c3dc-450163b93107.htm) | 当前的异形连接对象，如果设置了异形连接的话，仅用于异形模式的情况使用  The current alien connection object, if alien connection is set, is only used in the case of alien mode (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共属性代码示例 | [ByteTransform](720ae40b-1a12-b12f-7845-11c3bf80de8a.htm) | 当前的数据变换机制，当你需要从字节数据转换类型数据的时候需要。  The current data transformation mechanism is required when you need to convert type data from byte data. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共属性 | [ConnectionId](902454e3-c22d-b58e-171c-c59aa43ab52b.htm) | 当前连接的唯一ID号，默认为长度20的guid码加随机数组成，方便列表管理，也可以自己指定  The unique ID number of the current connection. The default is a 20-digit guid code plus a random number. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共属性代码示例 | [ConnectTimeOut](8fc841ea-c377-ce38-79d9-0252e0ffb28c.htm) | 获取或设置连接的超时时间，单位是毫秒   Gets or sets the timeout for the connection, in milliseconds (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共属性代码示例 | [IpAddress](bea866c8-a585-e1c1-8d93-517eb2505454.htm) | 获取或是设置远程服务器的IP地址，如果是本机测试，那么需要设置为127.0.0.1   Get or set the IP address of the remote server. If it is a local test, then it needs to be set to 127.0.0.1 (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共属性 | [LocalBinding](4c347174-2cb7-c637-f92a-1b7502844749.htm) | 获取或设置绑定的本地的IP地址和端口号信息，如果端口设置为0，代表任何可用的端口  Get or set the bound local IP address and port number information, if the port is set to 0, it means any available port (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共属性代码示例 | [LogNet](46dfe698-46c6-73a2-5458-8fd568d264ce.htm) | 组件的日志工具，支持日志记录，只要实例化后，当前网络的基本信息，就以[DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm)等级进行输出  The component's logging tool supports logging. As long as the instantiation of the basic network information, the output will be output at [DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm) (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 公共属性代码示例 | [Port](45e6a4b4-1190-6d1a-e1e6-d8c010e15fe4.htm) | 获取或设置服务器的端口号，具体的值需要取决于对方的配置  Gets or sets the port number of the server. The specific value depends on the configuration of the other party. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共属性代码示例 | [ReceiveTimeOut](995db577-ab5c-6c9c-6f6f-253e26600472.htm) | 获取或设置接收服务器反馈的时间，如果为负数，则不接收反馈   Gets or sets the time to receive server feedback, and if it is a negative number, does not receive feedback (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共属性 | [SendBeforeHex](59a68d6e-a784-a757-40cf-138d6b05099b.htm) | 获取或设置在发送通信报文前追加发送的字节信息，HEX格式，通常用于lora组网时，需要携带 00 00 00 02 四个字节的站地址功能。  Obtain or set the byte information sent before sending communication packets, HEX format, usually used for LORA networking, you need to carry 00 00 00 02 four-byte station address function. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共属性 | [SleepTime](d275ea69-962f-70f9-b486-cb3dc3382b74.htm) | 获取或设置在正式接收对方返回数据前的时候，需要休息的时间，当设置为0的时候，不需要休息。  Get or set the time required to rest before officially receiving the data from the other party. When it is set to 0, no rest is required. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共属性 | [SocketKeepAliveTime](7f4b5ca3-e4db-76e4-ef2b-4c93bf9447ac.htm) | 获取或设置客户端的Socket的心跳时间信息，这个是Socket底层自动实现的心跳包，不基于协议层实现。默认小于0，不开启心跳检测，如果需要开启，设置 60\_000 比较合适，单位毫秒  Get or set the heartbeat time information of the Socket of the client. This is the heartbeat packet automatically implemented by the bottom layer of the Socket, not based on the protocol layer. The default value is less than 0, and heartbeat detection is not enabled. If you need to enable it, it is more appropriate to set 60\_000, in milliseconds. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共属性代码示例 | [Token](80db00a7-aaa5-164f-e87c-af5400dc80fa.htm) | 网络类的身份令牌，在hsl协议的模式下会有效，在和设备进行通信的时候是无效的  Network-type identity tokens will be valid in the hsl protocol mode and will not be valid when communicating with the device (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的属性 | [UseServerActivePush](0e77fcc5-1744-cedf-a8e1-958031151826.htm) | 获取或设置当前的连接是否激活从服务器主动推送的功能 (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 受保护的方法 | [AccountCertificate](e9ef8320-5daf-4629-dcc6-4076e7944a9a.htm) | 认证账号，根据已经设置的用户名和密码，进行发送服务器进行账号认证。  Authentication account, according to the user name and password that have been set, sending server for account authentication. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 受保护的方法 | [AccountCertificateAsync](2dc6eb5a-79c8-1e51-7aa8-c3d29769033a.htm) | 认证账号，根据已经设置的用户名和密码，进行发送服务器进行账号认证。  Authentication account, according to the user name and password that have been set, sending server for account authentication. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | [AppendKey](f2521b8c-0fb1-6b34-4adc-ca6804c7195f.htm) | 如果 key 已经存在并且是一个字符串， APPEND 命令将 value 追加到 key 原来的值的末尾。 如果 key 不存在， APPEND 就简单地将给定 key 设为 value ，就像执行 SET key value 一样。 返回追加 value 之后， key 中字符串的长度。 |
| 公共方法 | [AppendKeyAsync](c39b4c6a-b7f6-68fe-dfc9-7cce1b653e24.htm) | 如果 key 已经存在并且是一个字符串， APPEND 命令将 value 追加到 key 原来的值的末尾。 如果 key 不存在， APPEND 就简单地将给定 key 设为 value ，就像执行 SET key value 一样。 返回追加 value 之后， key 中字符串的长度。 |
| 公共方法 | [ChangePassword](78a37d4d-bdee-7a09-d0ce-dc245627b1d1.htm) | 修改Redis的密码信息，如果不需要密码，则传入空字符串即可 |
| 公共方法 | [ChangePasswordAsync](58f6f736-093d-6e4d-3ad3-19590fa2a3a3.htm) | 修改Redis的密码信息，如果不需要密码，则传入空字符串即可 |
| 受保护的方法 | [CheckReceiveDataComplete](ea6ac6fd-de29-39d0-1098-6d90d5a1fb5d.htm) | 检查当前从网口接收的数据是否是完整的，如果是完整的，则需要返回 True，表示数据接收立即完成，默认返回 True  Check whether the data currently received from the network port is complete, and if it is complete, you need to return True, indicating that the data reception is completed immediately, and the default value is True (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 受保护的方法 | [CheckRemoteToken](1d231984-101b-0fff-0daf-1d18581a4241.htm) | 检查当前的头子节信息的令牌是否是正确的，仅用于某些特殊的协议实现  Check whether the token of the current header subsection information is correct, only for some special protocol implementations (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 公共方法代码示例 | [ConnectClose](e235581f-1c77-3b52-6f49-a6d1c79e559a.htm) | 手动断开与远程服务器的连接，如果当前是长连接模式，那么就会切换到短连接模式  Manually disconnect from the remote server, if it is currently in long connection mode, it will switch to short connection mode (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法代码示例 | [ConnectCloseAsync](ecec49bf-ad46-e74b-69ce-211464eaea6d.htm) | 手动断开与远程服务器的连接，如果当前是长连接模式，那么就会切换到短连接模式  Manually disconnect from the remote server, if it is currently in long connection mode, it will switch to short connection mode (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法代码示例 | [ConnectServer](ea436d31-7950-42df-a9c0-3c749c47e31d.htm) | 尝试连接远程的服务器，如果连接成功，就切换短连接模式到长连接模式，后面的每次请求都共享一个通道，使得通讯速度更快速  Try to connect to a remote server. If the connection is successful, switch the short connection mode to the long connection mode. Each subsequent request will share a channel, making the communication speed faster. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法代码示例 | [ConnectServer(AlienSession)](a5e9ec5a-548c-3d66-2509-fafa5ec21c3b.htm) | 使用指定的套接字创建异形客户端，在异形客户端的模式下，网络通道需要被动创建。  Use the specified socket to create the alien client. In the alien client mode, the network channel needs to be created passively. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | [ConnectServer(MqttClient, String, String)](3ae46111-33b7-c7dc-02d8-f689c1084b5a.htm) | 使用一个MQTT中转服务器来连接设备对象，并进行相关的读取操作 (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法代码示例 | [ConnectServerAsync](7eeb4d56-e784-e48e-b49a-8498a8fe8384.htm) | 尝试连接远程的服务器，如果连接成功，就切换短连接模式到长连接模式，后面的每次请求都共享一个通道，使得通讯速度更快速  Try to connect to a remote server. If the connection is successful, switch the short connection mode to the long connection mode. Each subsequent request will share a channel, making the communication speed faster. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 受保护的方法代码示例 | [CreateSocketAndConnect(String, Int32)](a153ae27-92e3-4fe7-b2e7-c8d207bebbb0.htm) | 创建一个新的socket对象并连接到远程的地址，默认超时时间为10秒钟，需要指定ip地址以及端口号信息  Create a new socket object and connect to the remote address. The default timeout is 10 seconds. You need to specify the IP address and port number. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法代码示例 | [CreateSocketAndConnect(IPEndPoint, Int32, IPEndPoint)](7ce6288f-c6db-e968-d676-f851d2b05c77.htm) | 创建一个新的socket对象并连接到远程的地址，需要指定远程终结点，超时时间（单位是毫秒），如果需要绑定本地的IP或是端口，传入 local对象  To create a new socket object and connect to the remote address, you need to specify the remote endpoint, the timeout period (in milliseconds), if you need to bind the local IP or port, pass in the local object (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法代码示例 | [CreateSocketAndConnect(String, Int32, Int32)](d0e14a5f-1e97-e59d-078f-4fbf0412b865.htm) | 创建一个新的socket对象并连接到远程的地址，需要指定ip地址以及端口号信息，还有超时时间，单位是毫秒  To create a new socket object and connect to a remote address, you need to specify the IP address and port number information, and the timeout period in milliseconds (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法代码示例 | [CreateSocketAndConnectAsync(String, Int32)](0227bc60-1c7f-3c1d-0cd6-e49bb904d6df.htm) | 创建一个新的socket对象并连接到远程的地址，默认超时时间为10秒钟，需要指定ip地址以及端口号信息  Create a new socket object and connect to the remote address. The default timeout is 10 seconds. You need to specify the IP address and port number. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法代码示例 | [CreateSocketAndConnectAsync(IPEndPoint, Int32, IPEndPoint)](ccec07bf-cfa4-b302-210e-441821e095a4.htm) | 创建一个新的socket对象并连接到远程的地址，需要指定远程终结点，超时时间（单位是毫秒），如果需要绑定本地的IP或是端口，传入 local对象  To create a new socket object and connect to the remote address, you need to specify the remote endpoint, the timeout period (in milliseconds), if you need to bind the local IP or port, pass in the local object (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法代码示例 | [CreateSocketAndConnectAsync(String, Int32, Int32)](7ec25883-60e2-c7ac-2f90-0ea1ae0c4945.htm) | 创建一个新的socket对象并连接到远程的地址，需要指定ip地址以及端口号信息，还有超时时间，单位是毫秒  To create a new socket object and connect to a remote address, you need to specify the IP address and port number information, and the timeout period in milliseconds (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 公共方法 | [DBSize](5814129a-9a73-25ac-e726-fe1da141eb28.htm) | 返回当前数据库的 key 的数量。 |
| 公共方法 | [DBSizeAsync](243a1f5f-a474-2e20-1298-06eff125f8b5.htm) | 返回当前数据库的 key 的数量。 |
| 受保护的方法 | [DecideWhetherQAMessage](77091dac-82d7-8c67-d691-31ae5b2aaee4.htm) | 决定当前的消息是否是应答机制的消息内容，需要在客户端进行重写实现，如果是应答机制，返回 True, 否则返回 False  To determine whether the current message is the message content of the response mechanism, it needs to be rewritten on the client side. If it is the response mechanism, return True, otherwise return False (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | [DecrementKey(String)](0f2a8a0d-9e7f-c6ab-2081-bba6606c37a9.htm) | 将 key 中储存的数字值减一。如果 key 不存在，那么 key 的值会先被初始化为 0 ，然后再执行 DECR 操作。 如果值包含错误的类型，或字符串类型的值不能表示为数字，那么返回一个错误。 本操作的值限制在 64 位(bit)有符号数字表示之内。 返回执行 DECR 命令之后 key 的值。 |
| 公共方法 | [DecrementKey(String, Int64)](0bb5b243-e3e3-064c-1c76-219ad3e2768a.htm) | 将 key 所储存的值减去减量 decrement 。如果 key 不存在，那么 key 的值会先被初始化为 0 ，然后再执行 DECR 操作。 如果值包含错误的类型，或字符串类型的值不能表示为数字，那么返回一个错误。 本操作的值限制在 64 位(bit)有符号数字表示之内。 返回减去 decrement 之后， key 的值。 |
| 公共方法 | [DecrementKeyAsync(String)](d0ade959-e1b8-cd2c-6a58-7104767acb23.htm) | 将 key 中储存的数字值减一。如果 key 不存在，那么 key 的值会先被初始化为 0 ，然后再执行 DECR 操作。 如果值包含错误的类型，或字符串类型的值不能表示为数字，那么返回一个错误。 本操作的值限制在 64 位(bit)有符号数字表示之内。 返回执行 DECR 命令之后 key 的值。 |
| 公共方法 | [DecrementKeyAsync(String, Int64)](743ce5f1-4a4e-314d-424d-a739f9e2384a.htm) | 将 key 所储存的值减去减量 decrement 。如果 key 不存在，那么 key 的值会先被初始化为 0 ，然后再执行 DECR 操作。 如果值包含错误的类型，或字符串类型的值不能表示为数字，那么返回一个错误。 本操作的值限制在 64 位(bit)有符号数字表示之内。 返回减去 decrement 之后， key 的值。 |
| 受保护的方法 | [DeleteFileByName](59afa4a4-3f3c-9831-26ee-c7c13a213bf7.htm) | 删除一个指定的文件，如果文件不存在，直接返回 True，如果文件存在则直接删除，删除成功返回 True，如果发生了异常，返回False  Delete a specified file, if the file does not exist, return True directly, if the file exists, delete it directly, if the deletion is successful, return True, if an exception occurs, return  False (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 公共方法 | [DeleteHashKey(String, String)](3f085ae2-ff3d-3c08-98af-5796aba0d975.htm) | 删除哈希表 key 中的一个或多个指定域，不存在的域将被忽略。 |
| 公共方法 | [DeleteHashKey(String, String)](eeed0572-14a6-8f5f-9e59-0df1636a2cb0.htm) | 删除哈希表 key 中的一个或多个指定域，不存在的域将被忽略。返回被成功移除的域的数量，不包括被忽略的域。 |
| 公共方法 | [DeleteHashKeyAsync(String, String)](b973444b-9bfd-30ce-6e59-006588402029.htm) | 删除哈希表 key 中的一个或多个指定域，不存在的域将被忽略。 |
| 公共方法 | [DeleteHashKeyAsync(String, String)](b1cdcd69-7004-177f-29bd-414a6aff71e2.htm) | 删除哈希表 key 中的一个或多个指定域，不存在的域将被忽略。返回被成功移除的域的数量，不包括被忽略的域。 |
| 公共方法 | [DeleteKey(String)](6eab9b1d-4556-b9c8-e7fa-2a5c13bd9940.htm) | 删除给定的一个或多个 key 。不存在的 key 会被忽略。 |
| 公共方法 | [DeleteKey(String)](d0fa44e9-f8c1-e302-d4f0-cd57917d2bf0.htm) | 删除给定的一个或多个 key 。不存在的 key 会被忽略。 |
| 公共方法 | [DeleteKeyAsync(String)](b760fe40-12aa-c0d2-5638-acb37ac81383.htm) | 删除给定的一个或多个 key 。不存在的 key 会被忽略。 |
| 公共方法 | [DeleteKeyAsync(String)](a157ea04-65d0-a685-6fbe-ff377c7d7e66.htm) | 删除给定的一个或多个 key 。不存在的 key 会被忽略。 |
| 公共方法 | [Dispose](8a86a200-fbe0-3b45-0fed-c551405c9b59.htm) | 释放当前的资源，如果调用了本方法，那么该对象再使用的时候，需要重新实例化。  Release the current resource. If this method is called, the object needs to be instantiated again when it is used again. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 受保护的方法 | [Dispose(Boolean)](32f03275-9ae0-525c-84bc-6e2515d11c12.htm) | 释放当前的资源，并自动关闭长连接，如果设置了的话 (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | Equals | (继承自 Object。) |
| 公共方法 | [ExistsHashKey](5b0bbeb9-4ab9-1f83-9e09-5f7227ac28e6.htm) | 查看哈希表 key 中，给定域 field 是否存在。如果哈希表含有给定域，返回 1 。 如果哈希表不含有给定域，或 key 不存在，返回 0 。 |
| 公共方法 | [ExistsHashKeyAsync](a80d408e-d621-a502-7829-6c7ea2a081f9.htm) | 查看哈希表 key 中，给定域 field 是否存在。如果哈希表含有给定域，返回 1 。 如果哈希表不含有给定域，或 key 不存在，返回 0 。 |
| 公共方法 | [ExistsKey](2b1b5b08-13ed-88c3-445d-d3613b48bc95.htm) | 检查给定 key 是否存在。若 key 存在，返回 1 ，否则返回 0 。 |
| 公共方法 | [ExistsKeyAsync](0d4357ed-5de3-1444-fce2-0e9be619d20f.htm) | 检查给定 key 是否存在。若 key 存在，返回 1 ，否则返回 0 。 |
| 公共方法 | [ExpireKey](0272044f-cc55-6b6a-42f0-6a0d89f37fe7.htm) | 为给定 key 设置生存时间，当 key 过期时(生存时间为 0 )，它会被自动删除。设置成功返回 1 。当 key 不存在或者不能为 key 设置生存时间时，返回 0 。 |
| 公共方法 | [ExpireKeyAsync](53c18b2d-2169-0519-602e-f65302f490e1.htm) | 为给定 key 设置生存时间，当 key 过期时(生存时间为 0 )，它会被自动删除。设置成功返回 1 。当 key 不存在或者不能为 key 设置生存时间时，返回 0 。 |
| 受保护的方法 | [ExtraAfterReadFromCoreServer](9f8357e9-dff4-1002-9368-0744cdedff29.htm) | 和服务器交互完成的时候调用的方法，可以根据读写结果进行一些额外的操作，具体的操作需要根据实际的需求来重写实现  The method called when the interaction with the server is completed can perform some additional operations based on the read and write results. The specific operations need to be rewritten according to actual needs. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 受保护的方法代码示例 | [ExtraOnDisconnect](49d9b263-02d8-e822-6a33-49f93737b98b.htm) | 根据实际的协议选择是否重写本方法，有些协议在断开连接之前，需要发送一些报文来关闭当前的网络通道  Select whether to rewrite this method according to the actual protocol. Some protocols need to send some packets to close the current network channel before disconnecting. (重写 [NetworkDoubleBaseExtraOnDisconnect(Socket)](b8569c75-5e55-cd55-5251-207465a8ba39.htm).) |
| 受保护的方法代码示例 | [ExtraOnDisconnectAsync](f681350d-2aad-40e0-0824-30fbd3fe57ff.htm) | 根据实际的协议选择是否重写本方法，有些协议在断开连接之前，需要发送一些报文来关闭当前的网络通道  Select whether to rewrite this method according to the actual protocol. Some protocols need to send some packets to close the current network channel before disconnecting. (重写 [NetworkDoubleBaseExtraOnDisconnectAsync(Socket)](86aa00d8-b2d6-ada9-5926-2b7a2237fe6b.htm).) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | [FlushDB](38c28c3f-63bf-f1a0-67bc-6547799257f9.htm) | 清空当前的数据库的key信息 |
| 公共方法 | [FlushDBAsync](40e3ccb2-9fe7-c0d7-88b9-95266a13c336.htm) | 清空当前的数据库的key信息 |
| 受保护的方法 | [GetAvailableSocket](e86a6fc1-eb91-8de9-5ff8-799b0ea3786b.htm) | 获取本次操作的可用的网络通道，如果是短连接，就重新生成一个新的网络通道，如果是长连接，就复用当前的网络通道。  Obtain the available network channels for this operation. If it is a short connection, a new network channel is regenerated. If it is a long connection, the current network channel is reused. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 受保护的方法 | [GetAvailableSocketAsync](40cf9ce7-a5e6-5431-2e4b-56e3b7c98346.htm) | 获取本次操作的可用的网络通道，如果是短连接，就重新生成一个新的网络通道，如果是长连接，就复用当前的网络通道。  Obtain the available network channels for this operation. If it is a short connection, a new network channel is regenerated. If it is a long connection, the current network channel is reused. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法 | [GetListLength](36330604-2557-4fc1-bdd7-fc3aff05e997.htm) | 返回列表 key 的长度。如果 key 不存在，则 key 被解释为一个空列表，返回 0 .如果 key 不是列表类型，返回一个错误。 |
| 公共方法 | [GetListLengthAsync](0e09431d-45b3-b539-a618-a147ae377f88.htm) | 返回列表 key 的长度。如果 key 不存在，则 key 被解释为一个空列表，返回 0 .如果 key 不是列表类型，返回一个错误。 |
| 受保护的方法 | [GetNewNetMessage](481a40e0-8539-5fc8-3e2b-c25bcc952697.htm) | 获取一个新的消息对象的方法，需要在继承类里面进行重写  The method to get a new message object needs to be overridden in the inheritance class (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | [GetPipeSocket](19904680-372e-e0a6-1b27-b4180e759668.htm) | 获取当前用于通信的管道信息  Get the current pipe information used for communication (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | GetType | (继承自 Object。) |
| 公共方法 | [IncrementHashKey(String, String, Int64)](81e32c12-a45b-ebe3-5e74-6357c50256a8.htm) | 为哈希表 key 中的域 field 的值加上增量 increment 。增量也可以为负数，相当于对给定域进行减法操作。 如果 key 不存在，一个新的哈希表被创建并执行 HINCRBY 命令。返回执行 HINCRBY 命令之后，哈希表 key 中域 field 的值。 |
| 公共方法 | [IncrementHashKey(String, String, Single)](421d1ddd-aea7-bd21-76ed-b437d990aab6.htm) | 为哈希表 key 中的域 field 的值加上增量 increment 。增量也可以为负数，相当于对给定域进行减法操作。 如果 key 不存在，一个新的哈希表被创建并执行 HINCRBY 命令。返回执行 HINCRBY 命令之后，哈希表 key 中域 field 的值。 |
| 公共方法 | [IncrementHashKeyAsync(String, String, Int64)](9d381989-878b-7e2f-519a-86c11af596c9.htm) | 为哈希表 key 中的域 field 的值加上增量 increment 。增量也可以为负数，相当于对给定域进行减法操作。 如果 key 不存在，一个新的哈希表被创建并执行 HINCRBY 命令。返回执行 HINCRBY 命令之后，哈希表 key 中域 field 的值。 |
| 公共方法 | [IncrementHashKeyAsync(String, String, Single)](adb31272-47d3-1aee-b224-447cdda8e9aa.htm) | 为哈希表 key 中的域 field 的值加上增量 increment 。增量也可以为负数，相当于对给定域进行减法操作。 如果 key 不存在，一个新的哈希表被创建并执行 HINCRBY 命令。返回执行 HINCRBY 命令之后，哈希表 key 中域 field 的值。 |
| 公共方法 | [IncrementKey(String)](86472b84-a7e5-63fa-7adf-769b1e5df83e.htm) | 将 key 中储存的数字值增一。如果 key 不存在，那么 key 的值会先被初始化为 0 ，然后再执行 INCR 操作。 如果值包含错误的类型，或字符串类型的值不能表示为数字，那么返回一个错误。 返回执行 INCR 命令之后 key 的值。 |
| 公共方法 | [IncrementKey(String, Int64)](91c1b3c1-c826-ae22-41ec-fd210c7fc5b7.htm) | 将 key 所储存的值加上增量 increment 。如果 key 不存在，那么 key 的值会先被初始化为 0 ，然后再执行 INCR 操作。 如果值包含错误的类型，或字符串类型的值不能表示为数字，那么返回一个错误。 |
| 公共方法 | [IncrementKey(String, Single)](438bd452-d082-3631-34b0-586d8fd2224a.htm) | 将 key 所储存的值加上增量 increment 。如果 key 不存在，那么 key 的值会先被初始化为 0 ，然后再执行 INCRBYFLOAT 操作。 如果命令执行成功，那么 key 的值会被更新为（执行加法之后的）新值，并且新值会以字符串的形式返回给调用者 |
| 公共方法 | [IncrementKeyAsync(String)](0134664a-49ce-5aad-bf5e-eea10440d453.htm) | 将 key 中储存的数字值增一。如果 key 不存在，那么 key 的值会先被初始化为 0 ，然后再执行 INCR 操作。 如果值包含错误的类型，或字符串类型的值不能表示为数字，那么返回一个错误。 返回执行 INCR 命令之后 key 的值。 |
| 公共方法 | [IncrementKeyAsync(String, Int64)](83d89534-28ac-6edd-be62-cbd721b9b058.htm) | 将 key 所储存的值加上增量 increment 。如果 key 不存在，那么 key 的值会先被初始化为 0 ，然后再执行 INCR 操作。 如果值包含错误的类型，或字符串类型的值不能表示为数字，那么返回一个错误。 |
| 公共方法 | [IncrementKeyAsync(String, Single)](5c2ec4f8-9dd6-a174-032f-26bed895e1e7.htm) | 将 key 所储存的值加上增量 increment 。如果 key 不存在，那么 key 的值会先被初始化为 0 ，然后再执行 INCRBYFLOAT 操作。 如果命令执行成功，那么 key 的值会被更新为（执行加法之后的）新值，并且新值会以字符串的形式返回给调用者 |
| 受保护的方法代码示例 | [InitializationOnConnect](72414fc8-b43a-3377-acdc-c59ae0628e11.htm) | 根据实际的协议选择是否重写本方法，有些协议在创建连接之后，需要进行一些初始化的信号握手，才能最终建立网络通道。  Whether to rewrite this method is based on the actual protocol. Some protocols require some initial signal handshake to establish a network channel after the connection is created. (重写 [NetworkDoubleBaseInitializationOnConnect(Socket)](40d5eb26-9bca-11cf-42a1-55af6ea5aebd.htm).) |
| 受保护的方法代码示例 | [InitializationOnConnectAsync](9156d090-6b1d-7381-c396-81491b4b6756.htm) | 根据实际的协议选择是否重写本方法，有些协议在创建连接之后，需要进行一些初始化的信号握手，才能最终建立网络通道。  Whether to rewrite this method is based on the actual protocol. Some protocols require some initial signal handshake to establish a network channel after the connection is created. (重写 [NetworkDoubleBaseInitializationOnConnectAsync(Socket)](14cac21b-4001-da7d-09a6-941eba868e14.htm).) |
| 公共方法 | [IpAddressPing](d9611e18-4033-de95-444e-b91ae0a81709.htm) | 对当前设备的IP地址进行PING的操作，返回PING的结果，正常来说，返回Success  PING the IP address of the current device and return the PING result. Normally, it returns Success (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | [ListInsertAfter](11b6eb28-b2ee-57a5-c724-a99df9101c71.htm) | 将值 value 插入到列表 key 当中，位于值 pivot 之后。 当 pivot 不存在于列表 key 时，不执行任何操作。 当 key 不存在时， key 被视为空列表，不执行任何操作。 如果 key 不是列表类型，返回一个错误。 |
| 公共方法 | [ListInsertAfterAsync](9defd0d6-5681-4c26-0973-14ecca4fc821.htm) | 将值 value 插入到列表 key 当中，位于值 pivot 之后。 当 pivot 不存在于列表 key 时，不执行任何操作。 当 key 不存在时， key 被视为空列表，不执行任何操作。 如果 key 不是列表类型，返回一个错误。 |
| 公共方法 | [ListInsertBefore](b19eed6b-2778-4b91-6f6b-e75cd9f1088e.htm) | 将值 value 插入到列表 key 当中，位于值 pivot 之前。 当 pivot 不存在于列表 key 时，不执行任何操作。 当 key 不存在时， key 被视为空列表，不执行任何操作。 如果 key 不是列表类型，返回一个错误。 |
| 公共方法 | [ListInsertBeforeAsync](0bf5390c-af9f-ecd6-1148-34a20abee73d.htm) | 将值 value 插入到列表 key 当中，位于值 pivot 之前。 当 pivot 不存在于列表 key 时，不执行任何操作。 当 key 不存在时， key 被视为空列表，不执行任何操作。 如果 key 不是列表类型，返回一个错误。 |
| 公共方法 | [ListLeftPop](8865ed9c-69dc-2f28-7f53-42e1c03d152b.htm) | 移除并返回列表 key 的头元素。列表的头元素。当 key 不存在时，返回 nil 。 |
| 公共方法 | [ListLeftPopAsync](f4d7f464-73f7-ff5d-0834-e18b5a91c994.htm) | 移除并返回列表 key 的头元素。列表的头元素。当 key 不存在时，返回 nil 。 |
| 公共方法 | [ListLeftPush(String, String)](74a8557b-7b50-f4fe-ba71-807683204706.htm) | 将一个或多个值 value 插入到列表 key 的表头，如果 key 不存在，一个空列表会被创建并执行 LPUSH 操作。当 key 存在但不是列表类型时，返回一个错误。返回执行 LPUSH 命令后，列表的长度。 |
| 公共方法 | [ListLeftPush(String, String)](dca179dc-794a-9936-b7a4-66338fefe209.htm) | 将一个或多个值 value 插入到列表 key 的表头，如果 key 不存在，一个空列表会被创建并执行 LPUSH 操作。当 key 存在但不是列表类型时，返回一个错误。返回执行 LPUSH 命令后，列表的长度。 |
| 公共方法 | [ListLeftPushAsync(String, String)](64d628e5-166c-ab96-0e90-b7d59d79e4cb.htm) | 将一个或多个值 value 插入到列表 key 的表头，如果 key 不存在，一个空列表会被创建并执行 LPUSH 操作。当 key 存在但不是列表类型时，返回一个错误。返回执行 LPUSH 命令后，列表的长度。 |
| 公共方法 | [ListLeftPushAsync(String, String)](6404ecc4-4c22-cb1a-11b4-4d8f4cdafe52.htm) | 将一个或多个值 value 插入到列表 key 的表头，如果 key 不存在，一个空列表会被创建并执行 LPUSH 操作。当 key 存在但不是列表类型时，返回一个错误。返回执行 LPUSH 命令后，列表的长度。 |
| 公共方法 | [ListLeftPushX](33093159-9864-a0dd-6a09-8bf134942031.htm) | 将值 value 插入到列表 key 的表头，当且仅当 key 存在并且是一个列表。和 LPUSH 命令相反，当 key 不存在时， LPUSHX 命令什么也不做。 返回LPUSHX 命令执行之后，表的长度。 |
| 公共方法 | [ListLeftPushXAsync](6bd2ab31-ff58-71f0-b68f-f36d7f53d845.htm) | 将值 value 插入到列表 key 的表头，当且仅当 key 存在并且是一个列表。和 LPUSH 命令相反，当 key 不存在时， LPUSHX 命令什么也不做。 返回LPUSHX 命令执行之后，表的长度。 |
| 公共方法 | [ListRange](efe06eb0-b945-fe65-5e81-21888cdda04f.htm) | 返回列表 key 中指定区间内的元素，区间以偏移量 start 和 stop 指定。 下标(index)参数 start 和 stop 都以 0 为底，也就是说，以 0 表示列表的第一个元素，以 1 表示列表的第二个元素，以此类推。 你也可以使用负数下标，以 -1 表示列表的最后一个元素， -2 表示列表的倒数第二个元素，以此类推。 返回一个列表，包含指定区间内的元素。 |
| 公共方法 | [ListRangeAsync](32e344f9-cdac-fc0e-2c50-bf2dd54b4f98.htm) | 返回列表 key 中指定区间内的元素，区间以偏移量 start 和 stop 指定。 下标(index)参数 start 和 stop 都以 0 为底，也就是说，以 0 表示列表的第一个元素，以 1 表示列表的第二个元素，以此类推。 你也可以使用负数下标，以 -1 表示列表的最后一个元素， -2 表示列表的倒数第二个元素，以此类推。 返回一个列表，包含指定区间内的元素。 |
| 公共方法 | [ListRemoveElementMatch](0ab40e89-d265-0b48-05de-5c1a7b41d326.htm) | 根据参数 count 的值，移除列表中与参数 value 相等的元素。count 的值可以是以下几种： count > 0 : 从表头开始向表尾搜索，移除与 value 相等的元素，数量为 count 。 count < 0 : 从表尾开始向表头搜索，移除与 value 相等的元素，数量为 count 的绝对值。 count = 0 : 移除表中所有与 value 相等的值。 返回被移除的数量。 |
| 公共方法 | [ListRemoveElementMatchAsync](5610de76-3198-5b5a-8f3d-b8fc937d7c83.htm) | 根据参数 count 的值，移除列表中与参数 value 相等的元素。count 的值可以是以下几种： count > 0 : 从表头开始向表尾搜索，移除与 value 相等的元素，数量为 count 。 count < 0 : 从表尾开始向表头搜索，移除与 value 相等的元素，数量为 count 的绝对值。 count = 0 : 移除表中所有与 value 相等的值。 返回被移除的数量。 |
| 公共方法 | [ListRightPop](97855b86-6a07-76e2-b887-1e830be705de.htm) | 移除并返回列表 key 的尾元素。当 key 不存在时，返回 nil 。 |
| 公共方法 | [ListRightPopAsync](fde8cee2-f393-a463-d561-74568545edc7.htm) | 移除并返回列表 key 的尾元素。当 key 不存在时，返回 nil 。 |
| 公共方法 | [ListRightPopLeftPush](e0b58cc4-4d80-64d6-729c-0dcdaafc4701.htm) | 命令 RPOPLPUSH 在一个原子时间内，执行以下两个动作：  1. 将列表 source 中的最后一个元素( 尾元素)弹出，并返回给客户端。  2. 将 source 弹出的元素插入到列表 destination ，作为 destination 列表的的头元素。   举个例子，你有两个列表 source 和 destination ， source 列表有元素 a, b, c ， destination 列表有元素 x, y, z ，执行 RPOPLPUSH source destination 之后， source 列表包含元素 a, b ， destination 列表包含元素 c, x, y, z ，并且元素 c 会被返回给客户端。 如果 source 不存在，值 nil 被返回，并且不执行其他动作。 如果 source 和 destination 相同，则列表中的表尾元素被移动到表头，并返回该元素，可以把这种特殊情况视作列表的旋转( rotation)操作。 |
| 公共方法 | [ListRightPopLeftPushAsync](41069df2-0dad-6f01-895f-a91086a92fa9.htm) | 命令 RPOPLPUSH 在一个原子时间内，执行以下两个动作：  1. 将列表 source 中的最后一个元素( 尾元素)弹出，并返回给客户端。  2. 将 source 弹出的元素插入到列表 destination ，作为 destination 列表的的头元素。   举个例子，你有两个列表 source 和 destination ， source 列表有元素 a, b, c ， destination 列表有元素 x, y, z ，执行 RPOPLPUSH source destination 之后， source 列表包含元素 a, b ， destination 列表包含元素 c, x, y, z ，并且元素 c 会被返回给客户端。 如果 source 不存在，值 nil 被返回，并且不执行其他动作。 如果 source 和 destination 相同，则列表中的表尾元素被移动到表头，并返回该元素，可以把这种特殊情况视作列表的旋转( rotation)操作。 |
| 公共方法 | [ListRightPush(String, String)](bec290e4-5da3-397c-5c53-848990279e78.htm) | 将一个或多个值 value 插入到列表 key 的表尾(最右边)。 如果 key 不存在，一个空列表会被创建并执行 RPUSH 操作。当 key 存在但不是列表类型时，返回一个错误。 |
| 公共方法 | [ListRightPush(String, String)](af196915-e868-386c-8d4a-0cc302f17a1a.htm) | 将一个或多个值 value 插入到列表 key 的表尾(最右边)。 如果有多个 value 值，那么各个 value 值按从左到右的顺序依次插入到表尾：比如对一个空列表 mylist 执行 RPUSH mylist a b c ，得出的结果列表为 a b c ， 如果 key 不存在，一个空列表会被创建并执行 RPUSH 操作。当 key 存在但不是列表类型时，返回一个错误。 返回执行 RPUSH 操作后，表的长度。 |
| 公共方法 | [ListRightPushAsync(String, String)](b30a73f7-59c3-e2f9-923e-19c6e62445d3.htm) | 将一个或多个值 value 插入到列表 key 的表尾(最右边)。 如果 key 不存在，一个空列表会被创建并执行 RPUSH 操作。当 key 存在但不是列表类型时，返回一个错误。 |
| 公共方法 | [ListRightPushAsync(String, String)](adf383fc-794b-3145-2320-e661aba47f95.htm) | 将一个或多个值 value 插入到列表 key 的表尾(最右边)。 如果有多个 value 值，那么各个 value 值按从左到右的顺序依次插入到表尾：比如对一个空列表 mylist 执行 RPUSH mylist a b c ，得出的结果列表为 a b c ， 如果 key 不存在，一个空列表会被创建并执行 RPUSH 操作。当 key 存在但不是列表类型时，返回一个错误。 返回执行 RPUSH 操作后，表的长度。 |
| 公共方法 | [ListRightPushX](e6ccfd1a-25ee-5eff-d56c-0ff9eb017fbc.htm) | 将值 value 插入到列表 key 的表尾，当且仅当 key 存在并且是一个列表。 和 RPUSH 命令相反，当 key 不存在时， RPUSHX 命令什么也不做。 |
| 公共方法 | [ListRightPushXAsync](3252db2a-56cf-1aa2-7716-f9eb2596ec62.htm) | 将值 value 插入到列表 key 的表尾，当且仅当 key 存在并且是一个列表。 和 RPUSH 命令相反，当 key 不存在时， RPUSHX 命令什么也不做。 |
| 公共方法 | [ListSet](5f2e7b8b-b0e7-87e6-dac7-81b572cf3fa3.htm) | 设置数组的某一个索引的数据信息，当 index 参数超出范围，或对一个空列表( key 不存在)进行 LSET 时，返回一个错误。 |
| 公共方法 | [ListSetAsync](ca019ec2-eb32-41f8-a460-89c70f90ace9.htm) | 设置数组的某一个索引的数据信息，当 index 参数超出范围，或对一个空列表( key 不存在)进行 LSET 时，返回一个错误。 |
| 公共方法 | [ListTrim](acf21364-1692-f9da-d337-c460b12a012d.htm) | 对一个列表进行修剪(trim)，就是说，让列表只保留指定区间内的元素，不在指定区间之内的元素都将被删除。 举个例子，执行命令 LTRIM list 0 2 ，表示只保留列表 list 的前三个元素，其余元素全部删除。 下标( index)参数 start 和 stop 都以 0 为底，也就是说，以 0 表示列表的第一个元素，以 1 表示列表的第二个元素，以此类推。 你也可以使用负数下标，以 -1 表示列表的最后一个元素， -2 表示列表的倒数第二个元素，以此类推。 当 key 不是列表类型时，返回一个错误。 |
| 公共方法 | [ListTrimAsync](9e86492d-205c-8427-9052-50ea3d8c1c64.htm) | 对一个列表进行修剪(trim)，就是说，让列表只保留指定区间内的元素，不在指定区间之内的元素都将被删除。 举个例子，执行命令 LTRIM list 0 2 ，表示只保留列表 list 的前三个元素，其余元素全部删除。 下标( index)参数 start 和 stop 都以 0 为底，也就是说，以 0 表示列表的第一个元素，以 1 表示列表的第二个元素，以此类推。 你也可以使用负数下标，以 -1 表示列表的最后一个元素， -2 表示列表的倒数第二个元素，以此类推。 当 key 不是列表类型时，返回一个错误。 |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法 | [MoveKey](5aa50043-85b9-a245-606a-de2b69164218.htm) | 将当前数据库的 key 移动到给定的数据库 db 当中。 如果当前数据库(源数据库)和给定数据库(目标数据库)有相同名字的给定 key ，或者 key 不存在于当前数据库，那么 MOVE 没有任何效果。 因此，也可以利用这一特性，将 MOVE 当作锁(locking)原语(primitive)。 |
| 公共方法 | [MoveKeyAsync](a93a74bf-65bd-733b-e39a-06d705252d99.htm) | 将当前数据库的 key 移动到给定的数据库 db 当中。 如果当前数据库(源数据库)和给定数据库(目标数据库)有相同名字的给定 key ，或者 key 不存在于当前数据库，那么 MOVE 没有任何效果。 因此，也可以利用这一特性，将 MOVE 当作锁(locking)原语(primitive)。 |
| 公共方法 | [OperateLongNumberFromServer](d2cdcd44-a6fb-7aa7-8c9c-08ff43107f65.htm) | 向服务器请求指令，并返回long数字的结果对象 |
| 公共方法 | [OperateLongNumberFromServerAsync](0005132d-e8c6-d51b-69dc-c9f8feaf8a9d.htm) | 向服务器请求指令，并返回long数字的结果对象 |
| 公共方法 | [OperateNumberFromServer](69aecc05-6901-f53a-da37-ebf85bae7053.htm) | 向服务器请求指定，并返回数字的结果对象 |
| 公共方法 | [OperateNumberFromServerAsync](b53bd6fa-d215-55f7-bbac-2ee171ed7498.htm) | 向服务器请求指定，并返回数字的结果对象 |
| 公共方法 | [OperateStatusFromServer](d460c3bb-e484-3fcb-e206-d20dc7201d63.htm) | 向服务器请求指令，并返回状态的结果对象，通常用于写入的判断，或是请求类型的判断 |
| 公共方法 | [OperateStatusFromServerAsync](e1a241a3-e200-747f-bb8b-5b7578a08999.htm) | 向服务器请求指令，并返回状态的结果对象，通常用于写入的判断，或是请求类型的判断 |
| 公共方法 | [OperateStringFromServer](822914af-cf62-1ccc-d22e-2df7dfac7294.htm) | 向服务器请求指令，并返回字符串的结果对象 |
| 公共方法 | [OperateStringFromServerAsync](eb31e52d-5387-eb14-291e-22f0e64666a4.htm) | 向服务器请求指令，并返回字符串的结果对象 |
| 公共方法 | [OperateStringsFromServer](c85c1a0e-83f9-69c9-a9b9-8a947ca54583.htm) | 向服务器请求指令，并返回字符串数组的结果对象 |
| 公共方法 | [OperateStringsFromServerAsync](2db50b51-32b8-1e92-cb9a-f2c847437e50.htm) | 向服务器请求指令，并返回字符串数组的结果对象 |
| 公共方法 | [PackCommandWithHeader](f96881d6-e59f-9e16-07e3-1188eee1a934.htm) | 对当前的命令进行打包处理，通常是携带命令头内容，标记当前的命令的长度信息，需要进行重写，否则默认不打包  The current command is packaged, usually carrying the content of the command header, marking the length of the current command, and it needs to be rewritten, otherwise it is not packaged by default (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | [PersistKey](c68a840b-ce64-3303-9339-a10f88bcbb77.htm) | 移除给定 key 的生存时间，将这个 key 从『易失的』(带生存时间 key )转换成『持久的』(一个不带生存时间、永不过期的 key )。 当生存时间移除成功时，返回 1 . 如果 key 不存在或 key 没有设置生存时间，返回 0 。 |
| 公共方法 | [PersistKeyAsync](a5d7ee67-9041-6455-7fc1-c6e96a442799.htm) | 移除给定 key 的生存时间，将这个 key 从『易失的』(带生存时间 key )转换成『持久的』(一个不带生存时间、永不过期的 key )。 当生存时间移除成功时，返回 1 . 如果 key 不存在或 key 没有设置生存时间，返回 0 。 |
| 公共方法 | [Ping](9cdd4e99-90a2-0ae8-141c-26f0f4a4d698.htm) | 向服务器进行PING的操作，服务器会返回PONG操作 |
| 公共方法 | [PingAsync](d6c8e625-edb2-7ff2-a54d-e852d49c73c8.htm) | 向服务器进行PING的操作，服务器会返回PONG操作 |
| 公共方法 | [Publish](74f20db5-aaaa-af2e-b4e6-5d192196b222.htm) | 将信息 message 发送到指定的频道 channel，返回接收到信息 message 的订阅者数量。 |
| 公共方法 | [PublishAsync](61428af8-94a2-b2b0-419b-5c0d3df4bee5.htm) | 将信息 message 发送到指定的频道 channel，返回接收到信息 message 的订阅者数量。 |
| 公共方法代码示例 | [ReadT](08365a89-e2d7-2aa5-3fd5-ab1a5d60c440.htm) | 从设备里读取支持Hsl特性的数据内容， 该特性为[HslRedisKeyAttribute](0b3a62cc-e2c8-dfb7-89de-7da47983c3c0.htm)，[HslRedisListItemAttribute](2bbc9c13-e0d2-6942-b529-1ef6ad142370.htm)， [HslRedisListAttribute](fd07c099-2c4b-18b4-45dc-3f869fdc50c3.htm)，[HslRedisHashFieldAttribute](13118df4-a00f-12bb-887d-e78f8b46236b.htm) 详细参考代码示例的操作说明。 |
| 公共方法 | [ReadAllKeys](53c48696-f991-f734-12b0-1ad306a84c6b.htm) | 查找所有符合给定模式 pattern 的 key 。 \* 匹配数据库中所有 key。 h?llo 匹配 hello ， hallo 和 hxllo 等。 h[ae]llo 匹配 hello 和 hallo ，但不匹配 hillo 。 |
| 公共方法 | [ReadAllKeysAsync](84c1be8c-3286-0410-0bd1-bb0185ff197e.htm) | 查找所有符合给定模式 pattern 的 key 。 \* 匹配数据库中所有 key。 h?llo 匹配 hello ， hallo 和 hxllo 等。 h[ae]llo 匹配 hello 和 hallo ，但不匹配 hillo 。 |
| 公共方法 | [ReadAndWriteKey](a6908f0e-1991-5210-f6c6-1e8b06369e56.htm) | 将给定 key 的值设为 value ，并返回 key 的旧值(old value)。当 key 存在但不是字符串类型时，返回一个错误。 |
| 公共方法 | [ReadAndWriteKeyAsync](6686726f-93ba-6435-76b0-b501c31e637e.htm) | 将给定 key 的值设为 value ，并返回 key 的旧值(old value)。当 key 存在但不是字符串类型时，返回一个错误。 |
| 公共方法代码示例 | [ReadAsyncT](43a59914-1502-1353-0bf3-20281ac92ca9.htm) | 从设备里读取支持Hsl特性的数据内容， 该特性为[HslRedisKeyAttribute](0b3a62cc-e2c8-dfb7-89de-7da47983c3c0.htm)，[HslRedisListItemAttribute](2bbc9c13-e0d2-6942-b529-1ef6ad142370.htm)， [HslRedisListAttribute](fd07c099-2c4b-18b4-45dc-3f869fdc50c3.htm)，[HslRedisHashFieldAttribute](13118df4-a00f-12bb-887d-e78f8b46236b.htm) 详细参考代码示例的操作说明。 |
| 公共方法 | [ReadCustomer](a4b7f43e-34ba-8aac-a705-a9bc8d4a5298.htm) | 自定义的指令交互方法，该指令用空格分割，举例：LTRIM AAAAA 0 999 就是收缩列表，GET AAA 就是获取键值，需要对返回的数据进行二次分析 |
| 公共方法 | [ReadCustomerAsync](bafea5ed-4293-ae60-c8e7-bb3d3481db44.htm) | 自定义的指令交互方法，该指令用空格分割，举例：LTRIM AAAAA 0 999 就是收缩列表，GET AAA 就是获取键值，需要对返回的数据进行二次分析 |
| 公共方法代码示例 | [ReadFromCoreServer(Byte)](b3e14431-d773-0d1a-f453-a3ad60648c29.htm) | 将数据发送到当前的网络通道中，并从网络通道中接收一个[INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm)指定的完整的报文，网络通道将根据[GetAvailableSocket](e86a6fc1-eb91-8de9-5ff8-799b0ea3786b.htm)方法自动获取，本方法是线程安全的。  Send data to the current network channel and receive a complete message specified by [INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm) from the network channel. The network channel will be automatically obtained according to the [GetAvailableSocket](e86a6fc1-eb91-8de9-5ff8-799b0ea3786b.htm) method This method is thread-safe. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | [ReadFromCoreServer(IEnumerableByte)](580dae22-fad5-5816-d58f-43068e89e419.htm) | 将多个数据报文按顺序发到设备，并从设备接收返回的数据内容，然后拼接成一个Byte[]信息，需要重写[UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm)方法才能返回正确的结果。  Send multiple data packets to the device in sequence, and receive the returned data content from the device, and then splicing them into a Byte[] message, you need to rewrite [UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm) method to return the correct result. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServer(Byte, Boolean, Boolean)](1a33d974-805d-f149-3665-eb8d8a99ea84.htm) | 将数据发送到当前的网络通道中，并从网络通道中接收一个[INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm)指定的完整的报文，网络通道将根据[GetAvailableSocket](e86a6fc1-eb91-8de9-5ff8-799b0ea3786b.htm)方法自动获取，本方法是线程安全的。  Send data to the current network channel and receive a complete message specified by [INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm) from the network channel. The network channel will be automatically obtained according to the [GetAvailableSocket](e86a6fc1-eb91-8de9-5ff8-799b0ea3786b.htm) method This method is thread-safe. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServer(Socket, Byte, Boolean, Boolean)](b79f1191-25bc-6cfc-3c43-bd25b03bbbba.htm) | 将数据报文发送指定的网络通道上，根据当前指定的[INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm)类型，返回一条完整的数据指令  Sends a data message to the specified network channel, and returns a complete data command according to the currently specified [INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm) type (重写 [NetworkDoubleBaseReadFromCoreServer(Socket, Byte, Boolean, Boolean)](3b72b289-b201-7548-0981-ac049c9a4135.htm).) |
| 公共方法代码示例 | [ReadFromCoreServerAsync(Byte)](ec51a18e-0f5d-dbe1-c229-d62e8760ac97.htm) | 将数据发送到当前的网络通道中，并从网络通道中接收一个[INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm)指定的完整的报文，网络通道将根据[GetAvailableSocket](e86a6fc1-eb91-8de9-5ff8-799b0ea3786b.htm)方法自动获取，本方法是线程安全的。  Send data to the current network channel and receive a complete message specified by [INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm) from the network channel. The network channel will be automatically obtained according to the [GetAvailableSocket](e86a6fc1-eb91-8de9-5ff8-799b0ea3786b.htm) method This method is thread-safe. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | [ReadFromCoreServerAsync(IEnumerableByte)](0fce97fd-8441-6b54-0877-b8dfacc738e3.htm) | 将多个数据报文按顺序发到设备，并从设备接收返回的数据内容，然后拼接成一个Byte[]信息，需要重写[UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm)方法才能返回正确的结果。  Send multiple data packets to the device in sequence, and receive the returned data content from the device, and then splicing them into a Byte[] message, you need to rewrite [UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm) method to return the correct result. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServerAsync(Byte, Boolean, Boolean)](7525b379-04c2-8dac-b451-a8b5413d6747.htm) | 将数据发送到当前的网络通道中，并从网络通道中接收一个[INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm)指定的完整的报文，网络通道将根据[GetAvailableSocket](e86a6fc1-eb91-8de9-5ff8-799b0ea3786b.htm)方法自动获取，本方法是线程安全的。  Send data to the current network channel and receive a complete message specified by [INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm) from the network channel. The network channel will be automatically obtained according to the [GetAvailableSocket](e86a6fc1-eb91-8de9-5ff8-799b0ea3786b.htm) method This method is thread-safe. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServerAsync(Socket, Byte, Boolean, Boolean)](6bde9f59-8b15-abba-9638-3e949612874d.htm) | 将数据报文发送指定的网络通道上，根据当前指定的[INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm)类型，返回一条完整的数据指令  Sends a data message to the specified network channel, and returns a complete data command according to the currently specified [INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm) type (重写 [NetworkDoubleBaseReadFromCoreServerAsync(Socket, Byte, Boolean, Boolean)](07db4a00-b548-8bd6-5de6-12deb91f5311.htm).) |
| 公共方法 | [ReadHashKey(String, String)](b0d43de4-1f62-fac9-970b-7e48f8a4f1fc.htm) | 返回哈希表 key 中给定域 field 的值。当给定域不存在或是给定 key 不存在时，返回 nil |
| 公共方法 | [ReadHashKey(String, String)](e7041162-0c8b-4345-d679-f92f7315d454.htm) | 返回哈希表 key 中，一个或多个给定域的值。如果给定的域不存在于哈希表，那么返回一个 nil 值。 因为不存在的 key 被当作一个空哈希表来处理，所以对一个不存在的 key 进行 HMGET 操作将返回一个只带有 nil 值的表。 |
| 公共方法 | [ReadHashKeyAll](63828fbe-fbdf-772a-010b-1425b0584572.htm) | 返回哈希表 key 中，所有的域和值。在返回值里，紧跟每个域名(field name)之后是域的值(value)，所以返回值的长度是哈希表大小的两倍。 |
| 公共方法 | [ReadHashKeyAllAsync](31c54864-f7bc-b3ce-7e77-15f0c1dcc540.htm) | 返回哈希表 key 中，所有的域和值。在返回值里，紧跟每个域名(field name)之后是域的值(value)，所以返回值的长度是哈希表大小的两倍。 |
| 公共方法 | [ReadHashKeyAsync(String, String)](406f96a9-706e-3ba2-ca7d-80238a8156c7.htm) | 返回哈希表 key 中给定域 field 的值。当给定域不存在或是给定 key 不存在时，返回 nil |
| 公共方法 | [ReadHashKeyAsync(String, String)](bc7ab677-80ae-4b7e-6609-96c10916faff.htm) | 返回哈希表 key 中，一个或多个给定域的值。如果给定的域不存在于哈希表，那么返回一个 nil 值。 因为不存在的 key 被当作一个空哈希表来处理，所以对一个不存在的 key 进行 HMGET 操作将返回一个只带有 nil 值的表。 |
| 公共方法 | [ReadHashKeyLength](32d483da-49e0-a205-bdc4-a659711be1c0.htm) | 返回哈希表 key 中域的数量。 |
| 公共方法 | [ReadHashKeyLengthAsync](5d0d4570-1f16-4277-cff7-a40607a530ab.htm) | 返回哈希表 key 中域的数量。 |
| 公共方法 | [ReadHashKeys](420286f8-5876-ee1b-b328-0b6eba7cefac.htm) | 返回哈希表 key 中的所有域。当 key 不存在时，返回一个空表。 |
| 公共方法 | [ReadHashKeysAsync](bbc8310c-029d-1ea5-b0e3-239b29ee8743.htm) | 返回哈希表 key 中的所有域。当 key 不存在时，返回一个空表。 |
| 公共方法 | [ReadHashValues](1d8b8e06-3d34-f4c6-68d0-20c14b309de1.htm) | 返回哈希表 key 中所有域的值。当 key 不存在时，返回一个空表。 |
| 公共方法 | [ReadHashValuesAsync](403e7344-c665-def6-70d9-81e5b22a740b.htm) | 返回哈希表 key 中所有域的值。当 key 不存在时，返回一个空表。 |
| 公共方法 | [ReadKey(String)](ea47037b-d2bf-4b7b-e28d-555e2e0f373a.htm) | 返回 key 所关联的字符串值。如果 key 不存在那么返回特殊值 nil 。 假如 key 储存的值不是字符串类型，返回一个错误，因为 GET 只能用于处理字符串值。 |
| 公共方法 | [ReadKey(String)](ac5af648-f537-7c6c-a064-281d1b07052b.htm) | 返回所有(一个或多个)给定 key 的值。 如果给定的 key 里面，有某个 key 不存在，那么这个 key 返回特殊值 null 。因此，该命令永不失败。 |
| 公共方法 | [ReadKeyAsync(String)](1159cbb5-7384-a9f8-a1dc-3a697345a482.htm) | 返回 key 所关联的字符串值。如果 key 不存在那么返回特殊值 nil 。 假如 key 储存的值不是字符串类型，返回一个错误，因为 GET 只能用于处理字符串值。 |
| 公共方法 | [ReadKeyAsync(String)](ae8c05c5-8f24-dc12-0dae-2083c4ed4c04.htm) | 返回所有(一个或多个)给定 key 的值。 如果给定的 key 里面，有某个 key 不存在，那么这个 key 返回特殊值 null 。因此，该命令永不失败。 |
| 公共方法 | [ReadKeyLength](7746357b-7e1b-d150-d019-722b752f5eeb.htm) | 返回 key 所储存的字符串值的长度。当 key 储存的不是字符串值时，返回一个错误。返回符串值的长度。当 key 不存在时，返回 0 。 |
| 公共方法 | [ReadKeyLengthAsync](7dd1b305-77b1-920c-14c0-2fe38e970791.htm) | 返回 key 所储存的字符串值的长度。当 key 储存的不是字符串值时，返回一个错误。返回符串值的长度。当 key 不存在时，返回 0 。 |
| 公共方法 | [ReadKeyRange](2e3e9582-cb24-44e6-d78b-7c672acdd4c8.htm) | 返回 key 中字符串值的子字符串，字符串的截取范围由 start 和 end 两个偏移量决定(包括 start 和 end 在内)。 负数偏移量表示从字符串最后开始计数， -1 表示最后一个字符， -2 表示倒数第二个，以此类推。 返回截取得出的子字符串。 |
| 公共方法 | [ReadKeyRangeAsync](407fdd5b-8b57-75d2-63be-93294311cb30.htm) | 返回 key 中字符串值的子字符串，字符串的截取范围由 start 和 end 两个偏移量决定(包括 start 和 end 在内)。 负数偏移量表示从字符串最后开始计数， -1 表示最后一个字符， -2 表示倒数第二个，以此类推。 返回截取得出的子字符串。 |
| 公共方法 | [ReadKeyTTL](679079c3-36e5-1618-3409-18459a74a0c4.htm) | 以秒为单位，返回给定 key 的剩余生存时间(TTL, time to live)。 |
| 公共方法 | [ReadKeyTTLAsync](24887926-0fd5-b65f-d070-1d439a5273ef.htm) | 以秒为单位，返回给定 key 的剩余生存时间(TTL, time to live)。 |
| 公共方法 | [ReadKeyType](1f702770-7804-880a-038e-c002ed31c42f.htm) | 返回 key 所储存的值的类型。none (key不存在)，string (字符串)，list (列表)，set (集合)，zset (有序集)，hash (哈希表) |
| 公共方法 | [ReadKeyTypeAsync](91c104ff-a646-76de-9bed-5e64237a57fc.htm) | 返回 key 所储存的值的类型。none (key不存在)，string (字符串)，list (列表)，set (集合)，zset (有序集)，hash (哈希表) |
| 公共方法 | [ReadListByIndex](43c7dae8-4c01-d2c9-b38b-0e7bef68b581.htm) | 返回列表 key 中，下标为 index 的元素。下标(index)参数 start 和 stop 都以 0 为底，也就是说，以 0 表示列表的第一个元素，以 1 表示列表的第二个元素，以此类推。 你也可以使用负数下标，以 -1 表示列表的最后一个元素， -2 表示列表的倒数第二个元素，以此类推。如果 key 不是列表类型，返回一个错误。 |
| 公共方法 | [ReadListByIndexAsync](1e1163ad-8161-9390-cff2-7919c8145c92.htm) | 返回列表 key 中，下标为 index 的元素。下标(index)参数 start 和 stop 都以 0 为底，也就是说，以 0 表示列表的第一个元素，以 1 表示列表的第二个元素，以此类推。 你也可以使用负数下标，以 -1 表示列表的最后一个元素， -2 表示列表的倒数第二个元素，以此类推。如果 key 不是列表类型，返回一个错误。 |
| 公共方法 | [ReadRandomKey](0d358a8f-b9e0-c659-c218-f051aaedcbfe.htm) | 从当前数据库中随机返回(不删除)一个 key 。 当数据库不为空时，返回一个 key 。 当数据库为空时，返回 nil 。 |
| 公共方法 | [ReadRandomKeyAsync](b97f91d2-2ab6-a39b-86f7-71238cd44fac.htm) | 从当前数据库中随机返回(不删除)一个 key 。 当数据库不为空时，返回一个 key 。 当数据库为空时，返回 nil 。 |
| 公共方法 | [ReadServerTime](c6c5f1fa-4639-0a32-8a42-48d56985e8e6.htm) | 获取服务器的时间戳信息，可用于本地时间的数据同步问题 |
| 公共方法 | [ReadServerTimeAsync](53022c3b-766e-c31b-1b55-548ceb3b2f42.htm) | 获取服务器的时间戳信息，可用于本地时间的数据同步问题 |
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
| 公共方法 | [RenameKey](ecf533af-233c-16b8-f88d-cb232e115013.htm) | 将 key 改名为 newkey 。 当 key 和 newkey 相同，或者 key 不存在时，返回一个错误。 当 newkey 已经存在时， RENAME 命令将覆盖旧值。 |
| 公共方法 | [RenameKeyAsync](20073ebe-0c16-2f57-a96b-79246e01129d.htm) | 将 key 改名为 newkey 。 当 key 和 newkey 相同，或者 key 不存在时，返回一个错误。 当 newkey 已经存在时， RENAME 命令将覆盖旧值。 |
| 公共方法 | [Save](261a8daa-9890-95bc-f8e4-02b78eb86277.htm) | SAVE 命令执行一个同步保存操作，将当前 Redis 实例的所有数据快照(snapshot)以 RDB 文件的形式保存到硬盘。 |
| 公共方法 | [SaveAsync](88c2043f-e9e6-9d95-4a2e-455452cf8adf.htm) | 在后台异步(Asynchronously)保存当前数据库的数据到磁盘。 BGSAVE 命令执行之后立即返回 OK ，然后 Redis fork 出一个新子进程，原来的 Redis 进程(父进程)继续处理客户端请求，而子进程则负责将数据保存到磁盘，然后退出。 |
| 公共方法 | [SelectDB](584a38b8-96e0-736f-7a24-981a8e920574.htm) | 切换到指定的数据库，数据库索引号 index 用数字值指定，以 0 作为起始索引值。默认使用 0 号数据库。 |
| 公共方法 | [SelectDBAsync](462ed998-9938-a261-bde4-579b4b9d3dc1.htm) | 切换到指定的数据库，数据库索引号 index 用数字值指定，以 0 作为起始索引值。默认使用 0 号数据库。 |
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
| 公共方法 | [SetAdd(String, String)](061c607b-58ff-1585-f1b0-be2f6950dc8d.htm) | 将一个member 元素加入到集合 key 当中，已经存在于集合的 member 元素将被忽略。假如 key 不存在，则创建一个只包含 member 元素作成员的集合。当 key 不是集合类型时，返回一个错误。 |
| 公共方法 | [SetAdd(String, String)](687ec76f-50a9-5ecb-7732-47a4944373d4.htm) | 将一个或多个 member 元素加入到集合 key 当中，已经存在于集合的 member 元素将被忽略。假如 key 不存在，则创建一个只包含 member 元素作成员的集合。当 key 不是集合类型时，返回一个错误。 |
| 公共方法 | [SetAddAsync(String, String)](4fc5a399-e408-8cf9-11fe-10cb0e718cb6.htm) | 将一个member 元素加入到集合 key 当中，已经存在于集合的 member 元素将被忽略。假如 key 不存在，则创建一个只包含 member 元素作成员的集合。当 key 不是集合类型时，返回一个错误。 |
| 公共方法 | [SetAddAsync(String, String)](ee87b3df-3919-93c3-7a7e-9dfda289d0f5.htm) | 将一个或多个 member 元素加入到集合 key 当中，已经存在于集合的 member 元素将被忽略。假如 key 不存在，则创建一个只包含 member 元素作成员的集合。当 key 不是集合类型时，返回一个错误。 |
| 公共方法 | [SetCard](cc0a3d1d-019b-f0ba-90a5-1bab767d6f05.htm) | 返回集合 key 的基数(集合中元素的数量)。当 key 不存在时，返回 0 。 |
| 公共方法 | [SetCardAsync](0d0665a9-a174-d154-f95b-2bf10d885bc3.htm) | 返回集合 key 的基数(集合中元素的数量)。当 key 不存在时，返回 0 。 |
| 公共方法 | [SetDiff(String, String)](e9ecf346-db86-b5f4-4f91-dcf5cdbefc1e.htm) | 返回一个集合的全部成员，该集合是所有给定集合之间的差集。 |
| 公共方法 | [SetDiff(String, String)](93f80f0e-351d-b56a-96ec-7a2af40349fe.htm) | 返回一个集合的全部成员，该集合是所有给定集合之间的差集。 |
| 公共方法 | [SetDiffAsync(String, String)](8f89285e-ce7c-802a-a473-9347afccba05.htm) | 返回一个集合的全部成员，该集合是所有给定集合之间的差集。 |
| 公共方法 | [SetDiffAsync(String, String)](2c1c06c0-6089-b734-a5f2-b8c7a27ff88d.htm) | 返回一个集合的全部成员，该集合是所有给定集合之间的差集。 |
| 公共方法 | [SetDiffStore(String, String, String)](5dcd4b1b-48b7-52a8-e0d5-7a4969b9aef5.htm) | 这个命令的作用和 SDIFF 类似，但它将结果保存到 destination 集合，而不是简单地返回结果集。如果 destination 集合已经存在，则将其覆盖。destination 可以是 key 本身。 |
| 公共方法 | [SetDiffStore(String, String, String)](7f4dcdf7-214e-a9c5-b93f-6fc71c2d8016.htm) | 这个命令的作用和 SDIFF 类似，但它将结果保存到 destination 集合，而不是简单地返回结果集。如果 destination 集合已经存在，则将其覆盖。destination 可以是 key 本身。 |
| 公共方法 | [SetDiffStoreAsync(String, String, String)](948e092b-b776-43f6-bf01-a432d1fe3f68.htm) | 这个命令的作用和 SDIFF 类似，但它将结果保存到 destination 集合，而不是简单地返回结果集。如果 destination 集合已经存在，则将其覆盖。destination 可以是 key 本身。 |
| 公共方法 | [SetDiffStoreAsync(String, String, String)](42a69dfd-f064-4337-8713-d89339c62906.htm) | 这个命令的作用和 SDIFF 类似，但它将结果保存到 destination 集合，而不是简单地返回结果集。如果 destination 集合已经存在，则将其覆盖。destination 可以是 key 本身。 |
| 公共方法 | [SetInter(String, String)](f627bd19-f348-9704-0581-a8d442ee60b4.htm) | 返回一个集合的全部成员，该集合是所有给定集合的交集。不存在的 key 被视为空集。当给定集合当中有一个空集时，结果也为空集(根据集合运算定律)。 |
| 公共方法 | [SetInter(String, String)](8b28cdae-05cc-3259-840f-c91cafa6d8bc.htm) | 返回一个集合的全部成员，该集合是所有给定集合的交集。不存在的 key 被视为空集。当给定集合当中有一个空集时，结果也为空集(根据集合运算定律)。 |
| 公共方法 | [SetInterAsync(String, String)](4035a5be-411b-638b-6559-b56e4442a2c4.htm) | 返回一个集合的全部成员，该集合是所有给定集合的交集。不存在的 key 被视为空集。当给定集合当中有一个空集时，结果也为空集(根据集合运算定律)。 |
| 公共方法 | [SetInterAsync(String, String)](4f9ab180-ff1f-bdf7-2c21-f14f093ac5dc.htm) | 返回一个集合的全部成员，该集合是所有给定集合的交集。不存在的 key 被视为空集。当给定集合当中有一个空集时，结果也为空集(根据集合运算定律)。 |
| 公共方法 | [SetInterStore(String, String, String)](1b4b8824-96fa-2ac7-b54b-ee3b2ece1c69.htm) | 这个命令类似于 SINTER 命令，但它将结果保存到 destination 集合，而不是简单地返回结果集。如果 destination 集合已经存在，则将其覆盖。destination 可以是 key 本身。 |
| 公共方法 | [SetInterStore(String, String, String)](1b14a4a9-4149-5e72-3d06-f4eab94d8d1f.htm) | 这个命令类似于 SINTER 命令，但它将结果保存到 destination 集合，而不是简单地返回结果集。如果 destination 集合已经存在，则将其覆盖。destination 可以是 key 本身。 |
| 公共方法 | [SetInterStoreAsync(String, String, String)](8d29dcd5-7a8a-39ed-5c0f-eb0b8ffa135f.htm) | 这个命令类似于 SINTER 命令，但它将结果保存到 destination 集合，而不是简单地返回结果集。如果 destination 集合已经存在，则将其覆盖。destination 可以是 key 本身。 |
| 公共方法 | [SetInterStoreAsync(String, String, String)](01f68462-704d-7d6c-c7f3-fa4c0edceb36.htm) | 这个命令类似于 SINTER 命令，但它将结果保存到 destination 集合，而不是简单地返回结果集。如果 destination 集合已经存在，则将其覆盖。destination 可以是 key 本身。 |
| 公共方法 | [SetIsMember](1c6d465b-353a-34f2-8f46-9d4037e58911.htm) | 判断 member 元素是否集合 key 的成员。如果 member 元素是集合的成员，返回 1 。如果 member 元素不是集合的成员，或 key 不存在，返回 0 。 |
| 公共方法 | [SetIsMemberAsync](21419054-3273-2aab-d486-633ab8861039.htm) | 判断 member 元素是否集合 key 的成员。如果 member 元素是集合的成员，返回 1 。如果 member 元素不是集合的成员，或 key 不存在，返回 0 。 |
| 公共方法 | [SetLoginAccount](725afc25-6e3f-bd7d-25ed-580270c8fc56.htm) | 设置当前的登录的账户名和密码信息，并启用账户验证的功能，账户名为空时设置不生效  Set the current login account name and password information, and enable the account verification function. The account name setting will not take effect when it is empty (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | [SetMembers](7a5b277f-1fe1-141f-346b-3722d97ac7ca.htm) | 返回集合 key 中的所有成员。不存在的 key 被视为空集合。 |
| 公共方法 | [SetMembersAsync](ce86155f-4be2-3836-3594-1398d78183ad.htm) | 返回集合 key 中的所有成员。不存在的 key 被视为空集合。 |
| 公共方法 | [SetMove](8b76c71a-c755-708b-e1e7-61eda8b9456c.htm) | 将 member 元素从 source 集合移动到 destination 集合。如果 source 集合不存在或不包含指定的 member 元素，则 SMOVE 命令不执行任何操作，仅返回 0 。 否则， member 元素从 source 集合中被移除，并添加到 destination 集合中去。当 destination 集合已经包含 member 元素时， SMOVE 命令只是简单地将 source 集合中的 member 元素删除。 当 source 或 destination 不是集合类型时，返回一个错误。 |
| 公共方法 | [SetMoveAsync](dc01893c-b09e-e134-49d3-76d0dbc308b9.htm) | 将 member 元素从 source 集合移动到 destination 集合。如果 source 集合不存在或不包含指定的 member 元素，则 SMOVE 命令不执行任何操作，仅返回 0 。 否则， member 元素从 source 集合中被移除，并添加到 destination 集合中去。当 destination 集合已经包含 member 元素时， SMOVE 命令只是简单地将 source 集合中的 member 元素删除。 当 source 或 destination 不是集合类型时，返回一个错误。 |
| 公共方法代码示例 | [SetPersistentConnection](d8e8df9f-cbe9-6863-54e9-aa3b38906b26.htm) | 在读取数据之前可以调用本方法将客户端设置为长连接模式，相当于跳过了ConnectServer的结果验证，对异形客户端无效，当第一次进行通信时再进行创建连接请求。  Before reading the data, you can call this method to set the client to the long connection mode, which is equivalent to skipping the result verification of ConnectServer, and it is invalid for the alien client. When the first communication is performed, the connection creation request is performed. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | [SetPipeSocket](b6a7a243-2a6d-c4ea-45b7-bf980450a2bd.htm) | 设置一个新的网络管道，一般来说不需要调用本方法，当多个网口设备共用一个网络连接时才需要使用本方法进行设置共享的管道。  To set up a new network channel, generally speaking, you do not need to call this method. This method is only needed to set up a shared channel when multiple network port devices share a network connection. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | [SetPop](8375117d-d7b5-1f33-cb59-80cc632b736d.htm) | 移除并返回集合中的一个随机元素。如果只想获取一个随机元素，但不想该元素从集合中被移除的话，可以使用 SRANDMEMBER 命令。 |
| 公共方法 | [SetPopAsync](66039533-6b1e-ad66-e89b-ef8f7d46c47b.htm) | 移除并返回集合中的一个随机元素。如果只想获取一个随机元素，但不想该元素从集合中被移除的话，可以使用 SRANDMEMBER 命令。 |
| 公共方法 | [SetRandomMember(String)](a812ee61-36b4-1e8e-7545-1c304eefcf3b.htm) | 那么返回集合中的一个随机元素。 |
| 公共方法 | [SetRandomMember(String, Int32)](fdd8ab0f-92da-8e38-2367-3cc6955d4770.htm) | 返回集合中的多个随机元素。  如果 count 为正数，且小于集合基数，那么命令返回一个包含 count 个元素的数组，数组中的元素各不相同。如果 count 大于等于集合基数，那么返回整个集合。  如果 count 为负数，那么命令返回一个数组，数组中的元素可能会重复出现多次，而数组的长度为 count 的绝对值。 |
| 公共方法 | [SetRandomMemberAsync(String)](ddc2b940-2954-4e8b-9a50-4f00dd729dca.htm) | 那么返回集合中的一个随机元素。 |
| 公共方法 | [SetRandomMemberAsync(String, Int32)](363a263f-5e96-9e89-f0d2-1a59458b56c9.htm) | 返回集合中的多个随机元素。  如果 count 为正数，且小于集合基数，那么命令返回一个包含 count 个元素的数组，数组中的元素各不相同。如果 count 大于等于集合基数，那么返回整个集合。  如果 count 为负数，那么命令返回一个数组，数组中的元素可能会重复出现多次，而数组的长度为 count 的绝对值。 |
| 公共方法 | [SetRemove(String, String)](915dac15-42ce-1db0-1975-6e2619ede6b1.htm) | 移除集合 key 中的一个元素，不存在的 member 元素会被忽略。 |
| 公共方法 | [SetRemove(String, String)](6f77ad2b-b35f-c91c-2cb6-b3e39cba9f3a.htm) | 移除集合 key 中的一个或多个 member 元素，不存在的 member 元素会被忽略。 |
| 公共方法 | [SetRemoveAsync(String, String)](9e56cae7-1944-e5fa-90c1-508e1c87499e.htm) | 移除集合 key 中的一个元素，不存在的 member 元素会被忽略。 |
| 公共方法 | [SetRemoveAsync(String, String)](3023a9ef-1c62-e0f9-c913-d3b15a77821e.htm) | 移除集合 key 中的一个或多个 member 元素，不存在的 member 元素会被忽略。 |
| 公共方法 | [SetUnion(String, String)](86a1df0d-4454-f293-fb36-9ddb123a92e8.htm) | 返回一个集合的全部成员，该集合是所有给定集合的并集。不存在的 key 被视为空集。 |
| 公共方法 | [SetUnion(String, String)](2d459ccd-11fe-cc5b-d3c5-526e2cd8268e.htm) | 返回一个或多个集合的全部成员，该集合是所有给定集合的并集。不存在的 key 被视为空集。 |
| 公共方法 | [SetUnionAsync(String, String)](1a1439cb-d4a0-9c89-b4e5-534f4e16eee6.htm) | 返回一个集合的全部成员，该集合是所有给定集合的并集。不存在的 key 被视为空集。 |
| 公共方法 | [SetUnionAsync(String, String)](3c66c82e-0b60-c57d-5f74-ebd1c408b83a.htm) | 返回一个或多个集合的全部成员，该集合是所有给定集合的并集。不存在的 key 被视为空集。 |
| 公共方法 | [SetUnionStore(String, String, String)](d3c6d04a-2dfb-f129-5104-56b9602661cf.htm) | 这个命令类似于 SUNION 命令，但它将结果保存到 destination 集合，而不是简单地返回结果集。如果 destination 已经存在，则将其覆盖。destination 可以是 key 本身。 |
| 公共方法 | [SetUnionStore(String, String, String)](683dea58-da66-fda5-20b5-35a393be979b.htm) | 这个命令类似于 SUNION 命令，但它将结果保存到 destination 集合，而不是简单地返回结果集。如果 destination 已经存在，则将其覆盖。destination 可以是 key 本身。 |
| 公共方法 | [SetUnionStoreAsync(String, String, String)](0b81d3ed-fd9c-4374-b670-f6a7d33ea7fc.htm) | 这个命令类似于 SUNION 命令，但它将结果保存到 destination 集合，而不是简单地返回结果集。如果 destination 已经存在，则将其覆盖。destination 可以是 key 本身。 |
| 公共方法 | [SetUnionStoreAsync(String, String, String)](121a97ed-2d7c-ef0b-4660-4efbdcb060f1.htm) | 这个命令类似于 SUNION 命令，但它将结果保存到 destination 集合，而不是简单地返回结果集。如果 destination 已经存在，则将其覆盖。destination 可以是 key 本身。 |
| 公共方法 | [SubscribeMessage(String)](5a08d3e8-17f8-273a-77ea-fa4f8c682f29.htm) | 从Redis服务器订阅一个或多个主题信息  Subscribe to one or more topics from the redis server |
| 公共方法 | [SubscribeMessage(String)](63b458d5-24ae-2600-ee11-1e0eff778662.htm) | 从Redis服务器订阅一个或多个主题信息  Subscribe to one or more topics from the redis server |
| 公共方法 | [ToString](6ff40e11-12f4-8759-45a0-594ca3e4fcdd.htm) | (重写 [NetworkDoubleBaseToString](85a8532f-8fc5-886c-9d76-70844e0a7d64.htm).) |
| 公共方法 | [UnpackResponseContent](32bad9a0-b59d-c171-ce85-74055edfc597.htm) | 根据对方返回的报文命令，对命令进行基本的拆包，例如各种Modbus协议拆包为统一的核心报文，还支持对报文的验证  According to the message command returned by the other party, the command is basically unpacked, for example, various Modbus protocols are unpacked into a unified core message, and the verification of the message is also supported (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | [UnSubscribeMessage(String)](dfdd1fbd-c51a-8456-ff3e-913a51b69b6d.htm) | 取消订阅一个或多个主题信息，取消之后，当前的订阅数据就不在接收到。  Unsubscribe from multiple topic information. After cancellation, the current subscription data will not be received. |
| 公共方法 | [UnSubscribeMessage(String)](e24cb10d-a565-8e43-8d83-4c6e083884d2.htm) | 取消订阅一个或多个主题信息，取消之后，当前的订阅数据就不在接收到。  Unsubscribe from multiple topic information. After cancellation, the current subscription data will not be received. |
| 公共方法代码示例 | [WriteT](70a08f37-0f86-76f1-55a5-df1ddd5a7d16.htm) | 从设备里写入支持Hsl特性的数据内容， 该特性为[HslRedisKeyAttribute](0b3a62cc-e2c8-dfb7-89de-7da47983c3c0.htm) ，[HslRedisHashFieldAttribute](13118df4-a00f-12bb-887d-e78f8b46236b.htm) 需要注意的是写入并不支持[HslRedisListAttribute](fd07c099-2c4b-18b4-45dc-3f869fdc50c3.htm)，[HslRedisListItemAttribute](2bbc9c13-e0d2-6942-b529-1ef6ad142370.htm)特性，详细参考代码示例的操作说明。 |
| 公共方法 | [WriteAndPublishKey](69029e5e-a78a-5c83-0300-b6dd057f9308.htm) | 将字符串值 value 关联到 key 。并发布一个订阅的频道数据，都成功时，才返回成功 |
| 公共方法 | [WriteAndPublishKeyAsync](d44abf4b-02bf-cfa5-cbc8-0b36edab4521.htm) | 将字符串值 value 关联到 key 。并发布一个订阅的频道数据，都成功时，才返回成功 |
| 公共方法代码示例 | [WriteAsyncT](9d54e352-2bd0-2655-3e14-49a0213e6040.htm) | 从设备里写入支持Hsl特性的数据内容， 该特性为[HslRedisKeyAttribute](0b3a62cc-e2c8-dfb7-89de-7da47983c3c0.htm) ，[HslRedisHashFieldAttribute](13118df4-a00f-12bb-887d-e78f8b46236b.htm) 需要注意的是写入并不支持[HslRedisListAttribute](fd07c099-2c4b-18b4-45dc-3f869fdc50c3.htm)，[HslRedisListItemAttribute](2bbc9c13-e0d2-6942-b529-1ef6ad142370.htm)特性，详细参考代码示例的操作说明。 |
| 公共方法 | [WriteExpireKey](816ca8c0-5e3a-4609-acc3-6f227ac2dfde.htm) | 将值 value 关联到 key ，并将 key 的生存时间设为 seconds (以秒为单位)。如果 key 已经存在， SETEX 命令将覆写旧值。 |
| 公共方法 | [WriteExpireKeyAsync](f1dce55a-5681-32f6-2c7c-cc664465490c.htm) | 将值 value 关联到 key ，并将 key 的生存时间设为 seconds (以秒为单位)。如果 key 已经存在， SETEX 命令将覆写旧值。 |
| 公共方法 | [WriteHashKey(String, String, String)](09ed7629-7800-0eff-6ad2-b9f6cbe9100e.htm) | 将哈希表 key 中的域 field 的值设为 value 。 如果 key 不存在，一个新的哈希表被创建并进行 HSET 操作。 如果域 field 已经存在于哈希表中，旧值将被覆盖。 如果 field 是哈希表中的一个新建域，并且值设置成功，返回 1 。 如果哈希表中域 field 已经存在且旧值已被新值覆盖，返回 0 。 |
| 公共方法 | [WriteHashKey(String, String, String)](85db4be6-4ef2-6062-2a2c-607a2c8255b6.htm) | 同时将多个 field-value (域-值)对设置到哈希表 key 中。 此命令会覆盖哈希表中已存在的域。 如果 key 不存在，一个空哈希表被创建并执行 HMSET 操作。 |
| 公共方法 | [WriteHashKeyAsync(String, String, String)](5801856f-9b52-fded-8714-80c52b90d2a8.htm) | 将哈希表 key 中的域 field 的值设为 value 。 如果 key 不存在，一个新的哈希表被创建并进行 HSET 操作。 如果域 field 已经存在于哈希表中，旧值将被覆盖。 如果 field 是哈希表中的一个新建域，并且值设置成功，返回 1 。 如果哈希表中域 field 已经存在且旧值已被新值覆盖，返回 0 。 |
| 公共方法 | [WriteHashKeyAsync(String, String, String)](f638b9b1-6813-8e8d-6b75-1c6784cec658.htm) | 同时将多个 field-value (域-值)对设置到哈希表 key 中。 此命令会覆盖哈希表中已存在的域。 如果 key 不存在，一个空哈希表被创建并执行 HMSET 操作。 |
| 公共方法 | [WriteHashKeyNx](819ddbb0-74d4-4667-6e55-372fac166c68.htm) | 将哈希表 key 中的域 field 的值设置为 value ，当且仅当域 field 不存在。若域 field 已经存在，该操作无效。 设置成功，返回 1 。如果给定域已经存在且没有操作被执行，返回 0 。 |
| 公共方法 | [WriteHashKeyNxAsync](b4571eb9-be29-fbb8-691b-f3a93d9245ec.htm) | 将哈希表 key 中的域 field 的值设置为 value ，当且仅当域 field 不存在。若域 field 已经存在，该操作无效。 设置成功，返回 1 。如果给定域已经存在且没有操作被执行，返回 0 。 |
| 公共方法 | [WriteKey(String, String)](eff07a08-c650-1d95-1d9d-d5e3736c439d.htm) | 将字符串值 value 关联到 key 。 如果 key 已经持有其他值， SET 就覆写旧值，无视类型。 对于某个原本带有生存时间（TTL）的键来说， 当 SET 命令成功在这个键上执行时，这个键原有的 TTL 将被清除。 |
| 公共方法 | [WriteKey(String, String)](97158b0a-0ca6-749b-ee3e-80557e13c619.htm) | 同时设置一个或多个 key-value 对。 如果某个给定 key 已经存在，那么 MSET 会用新值覆盖原来的旧值，如果这不是你所希望的效果，请考虑使用 MSETNX 命令：它只会在所有给定 key 都不存在的情况下进行设置操作。 |
| 公共方法 | [WriteKeyAsync(String, String)](50c841fd-8905-fe84-46ff-cfd3a0b8e4e5.htm) | 将字符串值 value 关联到 key 。 如果 key 已经持有其他值， SET 就覆写旧值，无视类型。 对于某个原本带有生存时间（TTL）的键来说， 当 SET 命令成功在这个键上执行时，这个键原有的 TTL 将被清除。 |
| 公共方法 | [WriteKeyAsync(String, String)](7f526544-53c5-6d54-5300-1c100dcb5c8e.htm) | 同时设置一个或多个 key-value 对。 如果某个给定 key 已经存在，那么 MSET 会用新值覆盖原来的旧值，如果这不是你所希望的效果，请考虑使用 MSETNX 命令：它只会在所有给定 key 都不存在的情况下进行设置操作。 |
| 公共方法 | [WriteKeyIfNotExists](9d48312b-2ce1-077f-bf47-3e6b5a328c45.htm) | 将 key 的值设为 value ，当且仅当 key 不存在。若给定的 key 已经存在，则 SETNX 不做任何动作。设置成功，返回 1 。设置失败，返回 0 。 |
| 公共方法 | [WriteKeyIfNotExistsAsync](e6844e26-ba17-a0a8-0ff2-bb1a1e8bce8e.htm) | 将 key 的值设为 value ，当且仅当 key 不存在。若给定的 key 已经存在，则 SETNX 不做任何动作。设置成功，返回 1 。设置失败，返回 0 。 |
| 公共方法 | [WriteKeyRange](c1f98cc9-f791-4a35-bf67-ea29e8930d32.htm) | 用 value 参数覆写(overwrite)给定 key 所储存的字符串值，从偏移量 offset 开始。不存在的 key 当作空白字符串处理。返回被 SETRANGE 修改之后，字符串的长度。 |
| 公共方法 | [WriteKeyRangeAsync](cdfaffec-518b-2754-5a10-3c5b68ee1767.htm) | 用 value 参数覆写(overwrite)给定 key 所储存的字符串值，从偏移量 offset 开始。不存在的 key 当作空白字符串处理。返回被 SETRANGE 修改之后，字符串的长度。 |
| 受保护的方法 | [WriteStreamFromSocket](97bf95ea-51fe-8f2e-5da6-1c83766aaa78.htm) | 从套接字中接收所有的数据然后写入到指定的流当中去，需要指定数据的长度，支持按照百分比进行进度报告  Receives all data from the socket and writes it to the specified stream. The length of the data needs to be specified, and progress reporting is supported in percentage. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [WriteStreamFromSocketAsync](b4cf7f6d-1e0c-a266-4c49-683fdb14a922.htm) | 从套接字中接收所有的数据然后写入到指定的流当中去，需要指定数据的长度，支持按照百分比进行进度报告  Receives all data from the socket and writes it to the specified stream. The length of the data needs to be specified, and progress reporting is supported in percentage. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 公共方法 | [ZSetAdd(String, String, Double)](593c2f94-650b-1b61-d5ec-09b5231295c0.htm) | 将一个 member 元素及其 score 值加入到有序集 key 当中。如果某个 member 已经是有序集的成员，那么更新这个 member 的 score 值，并通过重新插入这个 member 元素，来保证该 member 在正确的位置上。 score 值可以是整数值或双精度浮点数。  如果 key 不存在，则创建一个空的有序集并执行 ZADD 操作。当 key 存在但不是有序集类型时，返回一个错误。 |
| 公共方法 | [ZSetAdd(String, String, Double)](edbcdc0f-4537-5be6-9ddc-7c16d9773ee6.htm) | 将一个或多个 member 元素及其 score 值加入到有序集 key 当中。如果某个 member 已经是有序集的成员，那么更新这个 member 的 score 值，并通过重新插入这个 member 元素，来保证该 member 在正确的位置上。 score 值可以是整数值或双精度浮点数。  如果 key 不存在，则创建一个空的有序集并执行 ZADD 操作。当 key 存在但不是有序集类型时，返回一个错误。 |
| 公共方法 | [ZSetAddAsync(String, String, Double)](2e2ef743-be44-5c57-e604-a8c96b584693.htm) | 将一个 member 元素及其 score 值加入到有序集 key 当中。如果某个 member 已经是有序集的成员，那么更新这个 member 的 score 值，并通过重新插入这个 member 元素，来保证该 member 在正确的位置上。 score 值可以是整数值或双精度浮点数。  如果 key 不存在，则创建一个空的有序集并执行 ZADD 操作。当 key 存在但不是有序集类型时，返回一个错误。 |
| 公共方法 | [ZSetAddAsync(String, String, Double)](79d1d709-cc24-d2be-5e4d-51cb3395290a.htm) | 将一个或多个 member 元素及其 score 值加入到有序集 key 当中。如果某个 member 已经是有序集的成员，那么更新这个 member 的 score 值，并通过重新插入这个 member 元素，来保证该 member 在正确的位置上。 score 值可以是整数值或双精度浮点数。  如果 key 不存在，则创建一个空的有序集并执行 ZADD 操作。当 key 存在但不是有序集类型时，返回一个错误。 |
| 公共方法 | [ZSetCard](aa928601-c5f6-63fb-6bd8-e2e383796043.htm) | 返回有序集 key 的基数。 |
| 公共方法 | [ZSetCardAsync](b760ff15-4c7e-510f-1927-e17e6310e142.htm) | 返回有序集 key 的基数。 |
| 公共方法 | [ZSetCount](9872ddc0-8b8d-d9d6-4d08-64001fdefd29.htm) | 返回有序集 key 中， score 值在 min 和 max 之间(默认包括 score 值等于 min 或 max )的成员的数量。 |
| 公共方法 | [ZSetCountAsync](b294e1ee-27c7-7326-eb4b-8669040aecf0.htm) | 返回有序集 key 中， score 值在 min 和 max 之间(默认包括 score 值等于 min 或 max )的成员的数量。 |
| 公共方法 | [ZSetIncreaseBy](c84e0c14-15cd-451a-dec4-2a112a73772e.htm) | 为有序集 key 的成员 member 的 score 值加上增量 increment 。可以通过传递一个负数值 increment ，让 score 减去相应的值，比如 ZINCRBY key -5 member ，就是让 member 的 score 值减去 5 。 当 key 不存在，或 member 不是 key 的成员时， ZINCRBY key increment member 等同于 ZADD key increment member 。当 key 不是有序集类型时，返回一个错误。 |
| 公共方法 | [ZSetIncreaseByAsync](9805ac6b-aa05-2152-7dbb-7f8077ea462b.htm) | 为有序集 key 的成员 member 的 score 值加上增量 increment 。可以通过传递一个负数值 increment ，让 score 减去相应的值，比如 ZINCRBY key -5 member ，就是让 member 的 score 值减去 5 。 当 key 不存在，或 member 不是 key 的成员时， ZINCRBY key increment member 等同于 ZADD key increment member 。当 key 不是有序集类型时，返回一个错误。 |
| 公共方法 | [ZSetRange](759e2523-e2ad-0dc1-a904-354031a8ae50.htm) | 返回有序集 key 中，指定区间内的成员。其中成员的位置按 score 值递增(从小到大)来排序。具有相同 score 值的成员按字典序来排列。 下标参数 start 和 stop 都以 0 为底，也就是说，以 0 表示有序集第一个成员，以 1 表示有序集第二个成员，以此类推。你也可以使用负数下标，以 -1 表示最后一个成员， -2 表示倒数第二个成员，以此类推。 |
| 公共方法 | [ZSetRangeAsync](7724eb62-2653-65a7-779d-391d630c2faf.htm) | 返回有序集 key 中，指定区间内的成员。其中成员的位置按 score 值递增(从小到大)来排序。具有相同 score 值的成员按字典序来排列。 下标参数 start 和 stop 都以 0 为底，也就是说，以 0 表示有序集第一个成员，以 1 表示有序集第二个成员，以此类推。你也可以使用负数下标，以 -1 表示最后一个成员， -2 表示倒数第二个成员，以此类推。 |
| 公共方法 | [ZSetRangeByScore](acab27d0-6b6c-1f25-00ad-fafd4ac068fa.htm) | 返回有序集 key 中，所有 score 值介于 min 和 max 之间(包括等于 min 或 max )的成员。有序集成员按 score 值递增(从小到大)次序排列。 min 和 max 可以是 -inf 和 +inf ，这样一来，你就可以在不知道有序集的最低和最高 score 值的情况下，使用 ZRANGEBYSCORE 这类命令。 默认情况下，区间的取值使用闭区间 (小于等于或大于等于)，你也可以通过给参数前增加 "(" 符号来使用可选的开区间 (小于或大于)。"(5"代表不包含5 |
| 公共方法 | [ZSetRangeByScoreAsync](d1e0aa3c-d5f9-caeb-4017-699f956e661d.htm) | 返回有序集 key 中，所有 score 值介于 min 和 max 之间(包括等于 min 或 max )的成员。有序集成员按 score 值递增(从小到大)次序排列。 min 和 max 可以是 -inf 和 +inf ，这样一来，你就可以在不知道有序集的最低和最高 score 值的情况下，使用 ZRANGEBYSCORE 这类命令。 默认情况下，区间的取值使用闭区间 (小于等于或大于等于)，你也可以通过给参数前增加 "(" 符号来使用可选的开区间 (小于或大于)。"(5"代表不包含5 |
| 公共方法 | [ZSetRank](7957411f-32c1-4e07-a4f2-9d1948c63a2b.htm) | 返回有序集 key 中成员 member 的排名。其中有序集成员按 score 值递增(从小到大)顺序排列。排名以 0 为底，也就是说， score 值最小的成员排名为 0 。 |
| 公共方法 | [ZSetRankAsync](1181fd8c-02aa-9d53-cbeb-35d542bed241.htm) | 返回有序集 key 中成员 member 的排名。其中有序集成员按 score 值递增(从小到大)顺序排列。排名以 0 为底，也就是说， score 值最小的成员排名为 0 。 |
| 公共方法 | [ZSetRemove(String, String)](bde5c871-2a34-579a-a9f8-2e685047cf7b.htm) | 移除有序集 key 中的指定成员，不存在的成员将被忽略。当 key 存在但不是有序集类型时，返回一个错误。 |
| 公共方法 | [ZSetRemove(String, String)](66350924-7cee-4c1f-a584-4e480cd24dab.htm) | 移除有序集 key 中的一个或多个成员，不存在的成员将被忽略。当 key 存在但不是有序集类型时，返回一个错误。 |
| 公共方法 | [ZSetRemoveAsync(String, String)](4e9ef9ce-44d3-82be-bfcc-a2cf7743e6e0.htm) | 移除有序集 key 中的指定成员，不存在的成员将被忽略。当 key 存在但不是有序集类型时，返回一个错误。 |
| 公共方法 | [ZSetRemoveAsync(String, String)](f2528ee0-4f06-dd02-ab3f-b1232e123cd0.htm) | 移除有序集 key 中的一个或多个成员，不存在的成员将被忽略。当 key 存在但不是有序集类型时，返回一个错误。 |
| 公共方法 | [ZSetRemoveRangeByRank](cb81dfc2-6e30-0aca-8bcb-348cdc32bc35.htm) | 移除有序集 key 中，指定排名(rank)区间内的所有成员。区间分别以下标参数 start 和 stop 指出，包含 start 和 stop 在内。 下标参数 start 和 stop 都以 0 为底，也就是说，以 0 表示有序集第一个成员，以 1 表示有序集第二个成员，以此类推。你也可以使用负数下标，以 -1 表示最后一个成员， -2 表示倒数第二个成员，以此类推。 |
| 公共方法 | [ZSetRemoveRangeByRankAsync](9bde4b5f-b9bd-e01b-feb4-82c9b532dee4.htm) | 移除有序集 key 中，指定排名(rank)区间内的所有成员。区间分别以下标参数 start 和 stop 指出，包含 start 和 stop 在内。 下标参数 start 和 stop 都以 0 为底，也就是说，以 0 表示有序集第一个成员，以 1 表示有序集第二个成员，以此类推。你也可以使用负数下标，以 -1 表示最后一个成员， -2 表示倒数第二个成员，以此类推。 |
| 公共方法 | [ZSetRemoveRangeByScore](ac7be899-fe4e-3f12-d60f-81a038a927c7.htm) | 移除有序集 key 中，所有 score 值介于 min 和 max 之间(包括等于 min 或 max )的成员。 min 和 max 可以是 -inf 和 +inf ，这样一来，你就可以在不知道有序集的最低和最高 score 值的情况下，使用 ZRANGEBYSCORE 这类命令。 默认情况下，区间的取值使用闭区间 (小于等于或大于等于)，你也可以通过给参数前增加 "(" 符号来使用可选的开区间 (小于或大于)。例如"(5"代表不包括5 |
| 公共方法 | [ZSetRemoveRangeByScoreAsync](6ccee6bf-8635-14bf-2e7f-d1ac351c70c9.htm) | 移除有序集 key 中，所有 score 值介于 min 和 max 之间(包括等于 min 或 max )的成员。 min 和 max 可以是 -inf 和 +inf ，这样一来，你就可以在不知道有序集的最低和最高 score 值的情况下，使用 ZRANGEBYSCORE 这类命令。 默认情况下，区间的取值使用闭区间 (小于等于或大于等于)，你也可以通过给参数前增加 "(" 符号来使用可选的开区间 (小于或大于)。例如"(5"代表不包括5 |
| 公共方法 | [ZSetReverseRange](ef6defc7-62a2-c8a4-4025-32f6b3a82ae4.htm) | 返回有序集 key 中，指定区间内的成员。其中成员的位置按 score 值递减(从大到小)来排列。具有相同 score 值的成员按字典序来排列。 下标参数 start 和 stop 都以 0 为底，也就是说，以 0 表示有序集第一个成员，以 1 表示有序集第二个成员，以此类推。你也可以使用负数下标，以 -1 表示最后一个成员， -2 表示倒数第二个成员，以此类推。 |
| 公共方法 | [ZSetReverseRangeAsync](b40d37b7-1118-6aec-c7da-85ec54b457d0.htm) | 返回有序集 key 中，指定区间内的成员。其中成员的位置按 score 值递减(从大到小)来排列。具有相同 score 值的成员按字典序来排列。 下标参数 start 和 stop 都以 0 为底，也就是说，以 0 表示有序集第一个成员，以 1 表示有序集第二个成员，以此类推。你也可以使用负数下标，以 -1 表示最后一个成员， -2 表示倒数第二个成员，以此类推。 |
| 公共方法 | [ZSetReverseRangeByScore](3c1a26a8-f8e5-a835-7051-94b4d9c851ab.htm) | 返回有序集 key 中，所有 score 值介于 min 和 max 之间(包括等于 min 或 max )的成员。序集成员按 score 值递减(从大到小)的次序排列。 min 和 max 可以是 -inf 和 +inf ，这样一来，你就可以在不知道有序集的最低和最高 score 值的情况下，使用 ZRANGEBYSCORE 这类命令。 默认情况下，区间的取值使用闭区间 (小于等于或大于等于)，你也可以通过给参数前增加 ( 符号来使用可选的开区间 (小于或大于)。(5代表不包含5 |
| 公共方法 | [ZSetReverseRangeByScoreAsync](3a3cfbff-cee7-e1f3-d9c1-e7528db4e998.htm) | 返回有序集 key 中，所有 score 值介于 min 和 max 之间(包括等于 min 或 max )的成员。序集成员按 score 值递减(从大到小)的次序排列。 min 和 max 可以是 -inf 和 +inf ，这样一来，你就可以在不知道有序集的最低和最高 score 值的情况下，使用 ZRANGEBYSCORE 这类命令。 默认情况下，区间的取值使用闭区间 (小于等于或大于等于)，你也可以通过给参数前增加 ( 符号来使用可选的开区间 (小于或大于)。(5代表不包含5 |
| 公共方法 | [ZSetReverseRank](d524195b-6b02-a8e8-330a-3454a6ed21cd.htm) | 返回有序集 key 中成员 member 的排名。其中有序集成员按 score 值递减(从大到小)排序。排名以 0 为底，也就是说，score 值最大的成员排名为 0 。 |
| 公共方法 | [ZSetReverseRankAsync](c010781c-2851-7158-f953-1779929e3242.htm) | 返回有序集 key 中成员 member 的排名。其中有序集成员按 score 值递减(从大到小)排序。排名以 0 为底，也就是说，score 值最大的成员排名为 0 。 |
| 公共方法 | [ZSetScore](a98ee7e4-f361-f403-975b-a197e0ca0268.htm) | 返回有序集 key 中，成员 member 的 score 值。如果 member 元素不是有序集 key 的成员，或 key 不存在，返回 nil 。 |
| 公共方法 | [ZSetScoreAsync](fd24c433-944e-2ecb-4bcd-266bfc6825d4.htm) | 返回有序集 key 中，成员 member 的 score 值。如果 member 元素不是有序集 key 的成员，或 key 不存在，返回 nil 。 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)事件

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共事件 | [OnRedisMessageReceived](358f1cc5-4e05-6b87-8610-610ba86c5eb2.htm) | 当接收到Redis订阅的信息的时候触发 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)字段

|  | 名称 | 说明 |
| --- | --- | --- |
| 受保护的字段 | [fileCacheSize](2328fb89-cfda-7fa9-88d0-98f1ae3c7c54.htm) | 文件传输的时候的缓存大小，直接影响传输的速度，值越大，传输速度越快，越占内存，默认为100K大小  The size of the cache during file transfer directly affects the speed of the transfer. The larger the value, the faster the transfer speed and the more memory it takes. The default size is 100K. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的字段 | [isPersistentConn](17c281db-b635-a496-ef0e-8e91770c94be.htm) | 是否是长连接的状态  Whether it is a long connection state (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 受保护的字段 | [isUseAccountCertificate](c3f9fdea-29af-2f88-8f6a-f4651318a8a6.htm) | 是否使用账号登录，这个账户登录的功能是HSL组件创建的服务器特有的功能。  Whether to log in using an account. The function of this account login is a server-specific function created by the  HSL  component. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 受保护的字段 | [LogMsgFormatBinary](a00be7ba-6c1e-fd46-c68c-e004433341e5.htm) | 设置日志记录报文是否二进制，如果为False，那就使用ASCII码  Set whether the log message is binary, if it is False, then use ASCII code (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 受保护的字段 | [pipeSocket](d4efd9b0-ad04-ce88-5dde-843794e4e3fa.htm) | 当前的网络的管道信息 (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)备注

本类库的API指令的参考及注释来源：http://doc.redisfans.com/index.html

![](../icons/SectionExpanded.png)示例

基本的操作如下所示，举例了几个比较常见的指令，更多的需要参考api接口描述

基本操作代码

[复制](# "复制")

```
RedisClient redis = new RedisClient( "127.0.0.1", 6379, string.Empty );
// 如果连接服务器
OperateResult connect = redis.ConnectServer( );
if (connect.IsSuccess)
{
    Console.WriteLine( "connect success" );
}
else
{
    Console.WriteLine( "connect failed" );
    return;
}

// 连接服务器之后，默认选择的是db0，如果需要选择db1
redis.SelectDB( 1 );

// 好了，接下来是基本的读写操作
redis.WriteKey( "A", "test1" );
OperateResult<string> readA = redis.ReadKey( "A" );

// 支持批量读写操作
redis.WriteKey( new string[] { "A", "B", "C" }, new string[] { "test1", "test2", "test3" } );
OperateResult<string[]> readB = redis.ReadKey( new string[] { "A", "B", "C" } );

// 支持写入有效期的关键字，过两秒后读就不见了
redis.WriteExpireKey( "A", "test1", 2 );

// 自增自减操作
redis.WriteKey( "A", "1" );
OperateResult<long> result = redis.IncrementKey( "A" ); // A就变成2了，并且返回2
result = redis.DecrementKey( "A" ); // A就变成1了，并且返回1

// 列表的操作
redis.ListLeftPush( "List", "3" );
redis.ListLeftPush( "List", "2" );
redis.ListLeftPush( "List", "1" );  // 这样就生成了一个List的列表 [1,2,3]
var lists = redis.ListRange( "List", 0, -1 ); // 读取到了List的所有数据
var listIndex1 = redis.ReadListByIndex( "List", 1 ); // 读取某一个元素
var listLength = redis.GetListLength( "List" ).Content;  // 获取数组的长度
redis.ListTrim( "List", 0, 1 );   // 只保留2个元素

// 哈希键的操作
redis.WriteHashKey( "Hash", "A", "test1" );
redis.WriteHashKey( "Hash", "B", "test2" );
redis.WriteHashKey( "Hash", "C", "test3" );
redis.ReadHashKey( "Hash", "A" );   // 读取单个的值，返回test1
redis.ReadHashKeys( "Hash" );   // 读取所有域，["A","B","C"]
redis.ReadHashValues( "Hash" );  // 读取所有的值，["test1", "test2", "text3"]

// 其他的操作详细参考类的api文档信息

// 关闭
redis.ConnectClose( );
```

如下是基于特性的操作，有必要说明以下：

基础的使用

[复制](# "复制")

```
RedisClient redis = new RedisClient( "127.0.0.1", 6379, string.Empty );
redis.ConnectServer( );

// 上面是连接，此处不关心，只关心读数情况
// The above is the connection, here does not care, only care about the reading
OperateResult<string> readA = redis.ReadKey( "A" );
OperateResult<string> readB = redis.ReadKey( "B" );
OperateResult<string> readC = redis.ReadKey( "C" );
OperateResult<string> readD = redis.ReadKey( "D" );

if (readA.IsSuccess && readB.IsSuccess && readC.IsSuccess && readD.IsSuccess)
{
    // do somethong
}


// 当然你也可以这么干，加快读取的速度，并且代码更加简洁
// Of course, you can also do this, speed up the reading speed, and the code is more concise
OperateResult<string[]> readAD = redis.ReadKey( new string[] { "A", "B", "C", "D" } );
if (readAD.IsSuccess)
{
    // do something
    // A = readAD.Content[0]
    // B = readAD.Content[0]
    // C = readAD.Content[0]
    // D = readAD.Content[0]
}

// 当你需要读取加转换代码的时候，就没有那么方便了
// When you need to read plus conversion code, it is not so convenient
OperateResult<string[]> readAD2 = redis.ReadKey( new string[] { "A", "B", "C", "D" } );
if (readAD.IsSuccess)
{
    // 我们假设读取的所有数据需要转换int类型数据
    // We assume that all data read needs to be converted to int data
    int[] buffer = readAD2.Content.Select( m => int.Parse( m ) ).ToArray( );

    // do something
}

// 当你还要读取其他的信息的时候，就不得不多写点代码了
// When you need to read other information, you have to write more code
OperateResult<string> readListA = redis.ReadListByIndex( "List", 0 );
OperateResult<string> readHashA1 = redis.ReadHashKey( "HashA", "A1" );
OperateResult<string> readHashA2 = redis.ReadHashKey( "HashA", "A2" );
// ...
// and so on

redis.ConnectClose( );
```

总的来说，当读取的数据种类比较多的时候，读取的关键字比较多的时候，处理起来就比较的麻烦，此处推荐一个全新的写法，为了更好的对比，我们假设实现一种需求

同等代码

[复制](# "复制")

```
RedisClient redis = new RedisClient( "127.0.0.1", 6379, string.Empty );
redis.ConnectServer( );

// 上面是连接，此处不关心，只关心读数情况
// The above is the connection, here does not care, only care about the reading

// 我们来看看下面的两种等效写法，就可以发现此处的新写法是多么的便捷
// Let's take a look at the two equivalents below to see how convenient the new writing is here

// 假设下面是我们需要读取的，并且转化数据的，读取的数据在redis是真实存在的
// Suppose the following is what we need to read and transform the data. The read data is real in redis
OperateResult<string> readA = redis.ReadKey( "A" );
OperateResult<string> readB = redis.ReadKey( "B" );
OperateResult<string> readC = redis.ReadKey( "C" );
OperateResult<string> readD = redis.ReadKey( "D" );
OperateResult<string> readE = redis.ReadListByIndex( "E", 0 );
OperateResult<string[]> readF = redis.ListRange( "F", 0, -1 );
OperateResult<string[]> readG = redis.ReadHashKey( "G", new string[] { "G1", "G2", "G4" } );
OperateResult<string> readH = redis.ReadHashKey( "H", "H1" );

if (readA.IsSuccess && readB.IsSuccess && readC.IsSuccess && readD.IsSuccess &&
    readE.IsSuccess && readF.IsSuccess && readG.IsSuccess && readH.IsSuccess)
{
    // 进行相关的数据转换
    // Perform related data conversion
    string A = readA.Content;
    int B = int.Parse( readB.Content );
    ushort C = ushort.Parse( readC.Content );
    double D = double.Parse( readD.Content );
    float E = float.Parse( readE.Content );
    float[] F = readF.Content.Select( m => float.Parse( m ) ).ToArray( );
    int G1 = int.Parse( readG.Content[0] );
    int G2 = int.Parse( readG.Content[1] );
    int G3 = int.Parse( readG.Content[2] );
    string H = readH.Content;

    // do something
}

// 现在我们紧紧只需要一行代码即可完成任务，前提定义一个数据对象，还可以形成数据复用
// Now we need only one line of code to complete the task, provided that a data object is defined and data reuse can be formed
OperateResult<MyClass> read = redis.Read<MyClass>( );
if (read.IsSuccess)
{
    // do something
}

redis.ConnectClose( );
```

为此我们只需要实现一个特性类即可。代码如下：(注意，实际是很灵活的，类型都是自动转换的)

数据类

[复制](# "复制")

```
public class MyClass
{
    [HslCommunication.Reflection.HslRedisKey( "A" )]
    public string A { get; set; }

    [HslCommunication.Reflection.HslRedisKey( "B" )]
    public int B { get; set; }

    [HslCommunication.Reflection.HslRedisKey( "C" )]
    public ushort C { get; set; }

    [HslCommunication.Reflection.HslRedisKey( "D" )]
    public double D { get; set; }

    [HslCommunication.Reflection.HslRedisListItem( "E", 0 )]
    public float E { get; set; }

    [HslCommunication.Reflection.HslRedisList( "F", 0, -1 )]
    public float[] F { get; set; }

    [HslCommunication.Reflection.HslRedisHashField( "G", "G1" )]
    public int G1 { get; set; }

    [HslCommunication.Reflection.HslRedisHashField( "G", "G2" )]
    public int G2 { get; set; }

    [HslCommunication.Reflection.HslRedisHashField( "G", "G3" )]
    public int G3 { get; set; }

    [HslCommunication.Reflection.HslRedisHashField( "H", "H1" )]
    public string H { get; set; }
}
```

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.Enthernet.Redis 命名空间](cc44ce00-3781-c690-b233-24000870bf42.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## RedisClient 构造函数 

[原文連結](http://api.hslcommunication.cn/html/873ec4ce-423e-d519-ac3e-894ca64609ba.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Enthernet.Redis](../html/cc44ce00-3781-c690-b233-24000870bf42.htm "HslCommunication.Enthernet.Redis")

[RedisClient 类](../html/e33d9088-6df0-32e9-33f1-426e944c3b7f.htm "RedisClient 类")

[RedisClient 构造函数](../html/873ec4ce-423e-d519-ac3e-894ca64609ba.htm "RedisClient 构造函数 ")

[RedisClient 构造函数 (String)](../html/5407e3c7-7bdb-d5bc-5556-e51e73c318d0.htm "RedisClient 构造函数 (String)")

[RedisClient 构造函数 (String, Int32, String)](../html/3ba84c15-5c80-601c-6d24-fc66a6d1eba1.htm "RedisClient 构造函数 (String, Int32, String)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| RedisClient 构造函数 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [RedisClient(String)](5407e3c7-7bdb-d5bc-5556-e51e73c318d0.htm) | 实例化一个客户端对象，需要手动指定Ip地址和端口 |
| 公共方法 | [RedisClient(String, Int32, String)](3ba84c15-5c80-601c-6d24-fc66a6d1eba1.htm) | 实例化一个客户端的对象，用于和服务器通信 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[RedisClient 类](e33d9088-6df0-32e9-33f1-426e944c3b7f.htm)

[HslCommunication.Enthernet.Redis 命名空间](cc44ce00-3781-c690-b233-24000870bf42.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## RedisClient 构造函数 (String)

[原文連結](http://api.hslcommunication.cn/html/5407e3c7-7bdb-d5bc-5556-e51e73c318d0.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Enthernet.Redis](../html/cc44ce00-3781-c690-b233-24000870bf42.htm "HslCommunication.Enthernet.Redis")

[RedisClient 类](../html/e33d9088-6df0-32e9-33f1-426e944c3b7f.htm "RedisClient 类")

[RedisClient 构造函数](../html/873ec4ce-423e-d519-ac3e-894ca64609ba.htm "RedisClient 构造函数 ")

[RedisClient 构造函数 (String)](../html/5407e3c7-7bdb-d5bc-5556-e51e73c318d0.htm "RedisClient 构造函数 (String)")

[RedisClient 构造函数 (String, Int32, String)](../html/3ba84c15-5c80-601c-6d24-fc66a6d1eba1.htm "RedisClient 构造函数 (String, Int32, String)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| RedisClient 构造函数 (String) |

实例化一个客户端对象，需要手动指定Ip地址和端口

**命名空间：**
 [HslCommunication.Enthernet.Redis](cc44ce00-3781-c690-b233-24000870bf42.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public RedisClient(
	string password
)
```

```
Public Sub New ( 
	password As String
)
```

```
public:
RedisClient(
	String^ password
)
```

```
new : 
        password : string -> RedisClient
```

#### 参数

password
:   类型：SystemString  
    密码，如果服务器没有设置，密码设置为null

![](../icons/SectionExpanded.png)参见

#### 引用

[RedisClient 类](e33d9088-6df0-32e9-33f1-426e944c3b7f.htm)

[RedisClient 重载](873ec4ce-423e-d519-ac3e-894ca64609ba.htm)

[HslCommunication.Enthernet.Redis 命名空间](cc44ce00-3781-c690-b233-24000870bf42.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## RedisClient 构造函数 (String, Int32, String)

[原文連結](http://api.hslcommunication.cn/html/3ba84c15-5c80-601c-6d24-fc66a6d1eba1.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Enthernet.Redis](../html/cc44ce00-3781-c690-b233-24000870bf42.htm "HslCommunication.Enthernet.Redis")

[RedisClient 类](../html/e33d9088-6df0-32e9-33f1-426e944c3b7f.htm "RedisClient 类")

[RedisClient 构造函数](../html/873ec4ce-423e-d519-ac3e-894ca64609ba.htm "RedisClient 构造函数 ")

[RedisClient 构造函数 (String)](../html/5407e3c7-7bdb-d5bc-5556-e51e73c318d0.htm "RedisClient 构造函数 (String)")

[RedisClient 构造函数 (String, Int32, String)](../html/3ba84c15-5c80-601c-6d24-fc66a6d1eba1.htm "RedisClient 构造函数 (String, Int32, String)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| RedisClient 构造函数 (String, Int32, String) |

实例化一个客户端的对象，用于和服务器通信

**命名空间：**
 [HslCommunication.Enthernet.Redis](cc44ce00-3781-c690-b233-24000870bf42.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public RedisClient(
	string ipAddress,
	int port,
	string password
)
```

```
Public Sub New ( 
	ipAddress As String,
	port As Integer,
	password As String
)
```

```
public:
RedisClient(
	String^ ipAddress, 
	int port, 
	String^ password
)
```

```
new : 
        ipAddress : string * 
        port : int * 
        password : string -> RedisClient
```

#### 参数

ipAddress
:   类型：SystemString  
    服务器的ip地址

port
:   类型：SystemInt32  
    服务器的端口号

password
:   类型：SystemString  
    密码，如果服务器没有设置，密码设置为null

![](../icons/SectionExpanded.png)参见

#### 引用

[RedisClient 类](e33d9088-6df0-32e9-33f1-426e944c3b7f.htm)

[RedisClient 重载](873ec4ce-423e-d519-ac3e-894ca64609ba.htm)

[HslCommunication.Enthernet.Redis 命名空间](cc44ce00-3781-c690-b233-24000870bf42.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## RedisClient 属性

[原文連結](http://api.hslcommunication.cn/html/5708cfe8-172e-431e-d36e-02d5e23288e1.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Enthernet.Redis](../html/cc44ce00-3781-c690-b233-24000870bf42.htm "HslCommunication.Enthernet.Redis")

[RedisClient 类](../html/e33d9088-6df0-32e9-33f1-426e944c3b7f.htm "RedisClient 类")

[RedisClient 构造函数](../html/873ec4ce-423e-d519-ac3e-894ca64609ba.htm "RedisClient 构造函数 ")

[RedisClient 属性](../html/5708cfe8-172e-431e-d36e-02d5e23288e1.htm "RedisClient 属性")

[RedisClient 方法](../html/e41be4a5-1490-ffbd-f241-938f2d6c3f07.htm "RedisClient 方法")

[RedisClient 事件](../html/931bae45-0bc1-0713-e984-2d589dd7c843.htm "RedisClient 事件")

[RedisClient 字段](../html/4a43acdf-f8bd-333b-da81-99a1d696e961.htm "RedisClient 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| RedisClient 属性 |

[RedisClient](e33d9088-6df0-32e9-33f1-426e944c3b7f.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [AlienSession](5e30a9ab-7102-79a2-c3dc-450163b93107.htm) | 当前的异形连接对象，如果设置了异形连接的话，仅用于异形模式的情况使用  The current alien connection object, if alien connection is set, is only used in the case of alien mode (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共属性代码示例 | [ByteTransform](720ae40b-1a12-b12f-7845-11c3bf80de8a.htm) | 当前的数据变换机制，当你需要从字节数据转换类型数据的时候需要。  The current data transformation mechanism is required when you need to convert type data from byte data. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共属性 | [ConnectionId](902454e3-c22d-b58e-171c-c59aa43ab52b.htm) | 当前连接的唯一ID号，默认为长度20的guid码加随机数组成，方便列表管理，也可以自己指定  The unique ID number of the current connection. The default is a 20-digit guid code plus a random number. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共属性代码示例 | [ConnectTimeOut](8fc841ea-c377-ce38-79d9-0252e0ffb28c.htm) | 获取或设置连接的超时时间，单位是毫秒   Gets or sets the timeout for the connection, in milliseconds (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共属性代码示例 | [IpAddress](bea866c8-a585-e1c1-8d93-517eb2505454.htm) | 获取或是设置远程服务器的IP地址，如果是本机测试，那么需要设置为127.0.0.1   Get or set the IP address of the remote server. If it is a local test, then it needs to be set to 127.0.0.1 (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共属性 | [LocalBinding](4c347174-2cb7-c637-f92a-1b7502844749.htm) | 获取或设置绑定的本地的IP地址和端口号信息，如果端口设置为0，代表任何可用的端口  Get or set the bound local IP address and port number information, if the port is set to 0, it means any available port (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共属性代码示例 | [LogNet](46dfe698-46c6-73a2-5458-8fd568d264ce.htm) | 组件的日志工具，支持日志记录，只要实例化后，当前网络的基本信息，就以[DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm)等级进行输出  The component's logging tool supports logging. As long as the instantiation of the basic network information, the output will be output at [DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm) (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 公共属性代码示例 | [Port](45e6a4b4-1190-6d1a-e1e6-d8c010e15fe4.htm) | 获取或设置服务器的端口号，具体的值需要取决于对方的配置  Gets or sets the port number of the server. The specific value depends on the configuration of the other party. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共属性代码示例 | [ReceiveTimeOut](995db577-ab5c-6c9c-6f6f-253e26600472.htm) | 获取或设置接收服务器反馈的时间，如果为负数，则不接收反馈   Gets or sets the time to receive server feedback, and if it is a negative number, does not receive feedback (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共属性 | [SendBeforeHex](59a68d6e-a784-a757-40cf-138d6b05099b.htm) | 获取或设置在发送通信报文前追加发送的字节信息，HEX格式，通常用于lora组网时，需要携带 00 00 00 02 四个字节的站地址功能。  Obtain or set the byte information sent before sending communication packets, HEX format, usually used for LORA networking, you need to carry 00 00 00 02 four-byte station address function. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共属性 | [SleepTime](d275ea69-962f-70f9-b486-cb3dc3382b74.htm) | 获取或设置在正式接收对方返回数据前的时候，需要休息的时间，当设置为0的时候，不需要休息。  Get or set the time required to rest before officially receiving the data from the other party. When it is set to 0, no rest is required. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共属性 | [SocketKeepAliveTime](7f4b5ca3-e4db-76e4-ef2b-4c93bf9447ac.htm) | 获取或设置客户端的Socket的心跳时间信息，这个是Socket底层自动实现的心跳包，不基于协议层实现。默认小于0，不开启心跳检测，如果需要开启，设置 60\_000 比较合适，单位毫秒  Get or set the heartbeat time information of the Socket of the client. This is the heartbeat packet automatically implemented by the bottom layer of the Socket, not based on the protocol layer. The default value is less than 0, and heartbeat detection is not enabled. If you need to enable it, it is more appropriate to set 60\_000, in milliseconds. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共属性代码示例 | [Token](80db00a7-aaa5-164f-e87c-af5400dc80fa.htm) | 网络类的身份令牌，在hsl协议的模式下会有效，在和设备进行通信的时候是无效的  Network-type identity tokens will be valid in the hsl protocol mode and will not be valid when communicating with the device (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的属性 | [UseServerActivePush](0e77fcc5-1744-cedf-a8e1-958031151826.htm) | 获取或设置当前的连接是否激活从服务器主动推送的功能 (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[RedisClient 类](e33d9088-6df0-32e9-33f1-426e944c3b7f.htm)

[HslCommunication.Enthernet.Redis 命名空间](cc44ce00-3781-c690-b233-24000870bf42.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## RedisClient 方法

[原文連結](http://api.hslcommunication.cn/html/e41be4a5-1490-ffbd-f241-938f2d6c3f07.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Enthernet.Redis](../html/cc44ce00-3781-c690-b233-24000870bf42.htm "HslCommunication.Enthernet.Redis")

[RedisClient 类](../html/e33d9088-6df0-32e9-33f1-426e944c3b7f.htm "RedisClient 类")

[RedisClient 方法](../html/e41be4a5-1490-ffbd-f241-938f2d6c3f07.htm "RedisClient 方法")

[AppendKey 方法](../html/f2521b8c-0fb1-6b34-4adc-ca6804c7195f.htm "AppendKey 方法 ")

[AppendKeyAsync 方法](../html/c39b4c6a-b7f6-68fe-dfc9-7cce1b653e24.htm "AppendKeyAsync 方法 ")

[ChangePassword 方法](../html/78a37d4d-bdee-7a09-d0ce-dc245627b1d1.htm "ChangePassword 方法 ")

[ChangePasswordAsync 方法](../html/58f6f736-093d-6e4d-3ad3-19590fa2a3a3.htm "ChangePasswordAsync 方法 ")

[DBSize 方法](../html/5814129a-9a73-25ac-e726-fe1da141eb28.htm "DBSize 方法 ")

[DBSizeAsync 方法](../html/243a1f5f-a474-2e20-1298-06eff125f8b5.htm "DBSizeAsync 方法 ")

[DecrementKey 方法](../html/e9fa6c94-43ca-8291-79ca-6e04ca9289e7.htm "DecrementKey 方法 ")

[DecrementKeyAsync 方法](../html/6837e351-b8ef-a90e-c5e2-73b31621f837.htm "DecrementKeyAsync 方法 ")

[DeleteHashKey 方法](../html/19e3ebc3-e6fc-71cf-e80d-cb73e3011170.htm "DeleteHashKey 方法 ")

[DeleteHashKeyAsync 方法](../html/0dcc0d8a-1ed1-5c81-fb0e-e1c93be9d579.htm "DeleteHashKeyAsync 方法 ")

[DeleteKey 方法](../html/e2f0f305-d8bc-d515-3a62-cc91a3c1de76.htm "DeleteKey 方法 ")

[DeleteKeyAsync 方法](../html/0d8b5b0c-5979-9140-4ccf-992f4eaf1553.htm "DeleteKeyAsync 方法 ")

[ExistsHashKey 方法](../html/5b0bbeb9-4ab9-1f83-9e09-5f7227ac28e6.htm "ExistsHashKey 方法 ")

[ExistsHashKeyAsync 方法](../html/a80d408e-d621-a502-7829-6c7ea2a081f9.htm "ExistsHashKeyAsync 方法 ")

[ExistsKey 方法](../html/2b1b5b08-13ed-88c3-445d-d3613b48bc95.htm "ExistsKey 方法 ")

[ExistsKeyAsync 方法](../html/0d4357ed-5de3-1444-fce2-0e9be619d20f.htm "ExistsKeyAsync 方法 ")

[ExpireKey 方法](../html/0272044f-cc55-6b6a-42f0-6a0d89f37fe7.htm "ExpireKey 方法 ")

[ExpireKeyAsync 方法](../html/53c18b2d-2169-0519-602e-f65302f490e1.htm "ExpireKeyAsync 方法 ")

[ExtraOnDisconnect 方法](../html/49d9b263-02d8-e822-6a33-49f93737b98b.htm "ExtraOnDisconnect 方法 ")

[ExtraOnDisconnectAsync 方法](../html/f681350d-2aad-40e0-0824-30fbd3fe57ff.htm "ExtraOnDisconnectAsync 方法 ")

[FlushDB 方法](../html/38c28c3f-63bf-f1a0-67bc-6547799257f9.htm "FlushDB 方法 ")

[FlushDBAsync 方法](../html/40e3ccb2-9fe7-c0d7-88b9-95266a13c336.htm "FlushDBAsync 方法 ")

[GetListLength 方法](../html/36330604-2557-4fc1-bdd7-fc3aff05e997.htm "GetListLength 方法 ")

[GetListLengthAsync 方法](../html/0e09431d-45b3-b539-a618-a147ae377f88.htm "GetListLengthAsync 方法 ")

[IncrementHashKey 方法](../html/6404f402-9882-5453-7df7-275f61748021.htm "IncrementHashKey 方法 ")

[IncrementHashKeyAsync 方法](../html/7a69e1a1-8689-3e23-8aa3-cdad33e7e3f8.htm "IncrementHashKeyAsync 方法 ")

[IncrementKey 方法](../html/04d46420-bedd-9016-075c-3bb1a6623583.htm "IncrementKey 方法 ")

[IncrementKeyAsync 方法](../html/95f88e06-114c-b5b3-ee2f-2e9f97f27b8c.htm "IncrementKeyAsync 方法 ")

[InitializationOnConnect 方法](../html/72414fc8-b43a-3377-acdc-c59ae0628e11.htm "InitializationOnConnect 方法 ")

[InitializationOnConnectAsync 方法](../html/9156d090-6b1d-7381-c396-81491b4b6756.htm "InitializationOnConnectAsync 方法 ")

[ListInsertAfter 方法](../html/11b6eb28-b2ee-57a5-c724-a99df9101c71.htm "ListInsertAfter 方法 ")

[ListInsertAfterAsync 方法](../html/9defd0d6-5681-4c26-0973-14ecca4fc821.htm "ListInsertAfterAsync 方法 ")

[ListInsertBefore 方法](../html/b19eed6b-2778-4b91-6f6b-e75cd9f1088e.htm "ListInsertBefore 方法 ")

[ListInsertBeforeAsync 方法](../html/0bf5390c-af9f-ecd6-1148-34a20abee73d.htm "ListInsertBeforeAsync 方法 ")

[ListLeftPop 方法](../html/8865ed9c-69dc-2f28-7f53-42e1c03d152b.htm "ListLeftPop 方法 ")

[ListLeftPopAsync 方法](../html/f4d7f464-73f7-ff5d-0834-e18b5a91c994.htm "ListLeftPopAsync 方法 ")

[ListLeftPush 方法](../html/232a2dbd-e2c6-5fbf-af96-5269a10df4c9.htm "ListLeftPush 方法 ")

[ListLeftPushAsync 方法](../html/cb91133d-ce36-f4f7-a9d3-81dbda7d6617.htm "ListLeftPushAsync 方法 ")

[ListLeftPushX 方法](../html/33093159-9864-a0dd-6a09-8bf134942031.htm "ListLeftPushX 方法 ")

[ListLeftPushXAsync 方法](../html/6bd2ab31-ff58-71f0-b68f-f36d7f53d845.htm "ListLeftPushXAsync 方法 ")

[ListRange 方法](../html/efe06eb0-b945-fe65-5e81-21888cdda04f.htm "ListRange 方法 ")

[ListRangeAsync 方法](../html/32e344f9-cdac-fc0e-2c50-bf2dd54b4f98.htm "ListRangeAsync 方法 ")

[ListRemoveElementMatch 方法](../html/0ab40e89-d265-0b48-05de-5c1a7b41d326.htm "ListRemoveElementMatch 方法 ")

[ListRemoveElementMatchAsync 方法](../html/5610de76-3198-5b5a-8f3d-b8fc937d7c83.htm "ListRemoveElementMatchAsync 方法 ")

[ListRightPop 方法](../html/97855b86-6a07-76e2-b887-1e830be705de.htm "ListRightPop 方法 ")

[ListRightPopAsync 方法](../html/fde8cee2-f393-a463-d561-74568545edc7.htm "ListRightPopAsync 方法 ")

[ListRightPopLeftPush 方法](../html/e0b58cc4-4d80-64d6-729c-0dcdaafc4701.htm "ListRightPopLeftPush 方法 ")

[ListRightPopLeftPushAsync 方法](../html/41069df2-0dad-6f01-895f-a91086a92fa9.htm "ListRightPopLeftPushAsync 方法 ")

[ListRightPush 方法](../html/d71ca156-74d3-536e-6925-315d948c0ace.htm "ListRightPush 方法 ")

[ListRightPushAsync 方法](../html/53144654-da3d-b367-b39c-0dbc68275e59.htm "ListRightPushAsync 方法 ")

[ListRightPushX 方法](../html/e6ccfd1a-25ee-5eff-d56c-0ff9eb017fbc.htm "ListRightPushX 方法 ")

[ListRightPushXAsync 方法](../html/3252db2a-56cf-1aa2-7716-f9eb2596ec62.htm "ListRightPushXAsync 方法 ")

[ListSet 方法](../html/5f2e7b8b-b0e7-87e6-dac7-81b572cf3fa3.htm "ListSet 方法 ")

[ListSetAsync 方法](../html/ca019ec2-eb32-41f8-a460-89c70f90ace9.htm "ListSetAsync 方法 ")

[ListTrim 方法](../html/acf21364-1692-f9da-d337-c460b12a012d.htm "ListTrim 方法 ")

[ListTrimAsync 方法](../html/9e86492d-205c-8427-9052-50ea3d8c1c64.htm "ListTrimAsync 方法 ")

[MoveKey 方法](../html/5aa50043-85b9-a245-606a-de2b69164218.htm "MoveKey 方法 ")

[MoveKeyAsync 方法](../html/a93a74bf-65bd-733b-e39a-06d705252d99.htm "MoveKeyAsync 方法 ")

[OperateLongNumberFromServer 方法](../html/d2cdcd44-a6fb-7aa7-8c9c-08ff43107f65.htm "OperateLongNumberFromServer 方法 ")

[OperateLongNumberFromServerAsync 方法](../html/0005132d-e8c6-d51b-69dc-c9f8feaf8a9d.htm "OperateLongNumberFromServerAsync 方法 ")

[OperateNumberFromServer 方法](../html/69aecc05-6901-f53a-da37-ebf85bae7053.htm "OperateNumberFromServer 方法 ")

[OperateNumberFromServerAsync 方法](../html/b53bd6fa-d215-55f7-bbac-2ee171ed7498.htm "OperateNumberFromServerAsync 方法 ")

[OperateStatusFromServer 方法](../html/d460c3bb-e484-3fcb-e206-d20dc7201d63.htm "OperateStatusFromServer 方法 ")

[OperateStatusFromServerAsync 方法](../html/e1a241a3-e200-747f-bb8b-5b7578a08999.htm "OperateStatusFromServerAsync 方法 ")

[OperateStringFromServer 方法](../html/822914af-cf62-1ccc-d22e-2df7dfac7294.htm "OperateStringFromServer 方法 ")

[OperateStringFromServerAsync 方法](../html/eb31e52d-5387-eb14-291e-22f0e64666a4.htm "OperateStringFromServerAsync 方法 ")

[OperateStringsFromServer 方法](../html/c85c1a0e-83f9-69c9-a9b9-8a947ca54583.htm "OperateStringsFromServer 方法 ")

[OperateStringsFromServerAsync 方法](../html/2db50b51-32b8-1e92-cb9a-f2c847437e50.htm "OperateStringsFromServerAsync 方法 ")

[PersistKey 方法](../html/c68a840b-ce64-3303-9339-a10f88bcbb77.htm "PersistKey 方法 ")

[PersistKeyAsync 方法](../html/a5d7ee67-9041-6455-7fc1-c6e96a442799.htm "PersistKeyAsync 方法 ")

[Ping 方法](../html/9cdd4e99-90a2-0ae8-141c-26f0f4a4d698.htm "Ping 方法 ")

[PingAsync 方法](../html/d6c8e625-edb2-7ff2-a54d-e852d49c73c8.htm "PingAsync 方法 ")

[Publish 方法](../html/74f20db5-aaaa-af2e-b4e6-5d192196b222.htm "Publish 方法 ")

[PublishAsync 方法](../html/61428af8-94a2-b2b0-419b-5c0d3df4bee5.htm "PublishAsync 方法 ")

[Read(T) 方法](../html/08365a89-e2d7-2aa5-3fd5-ab1a5d60c440.htm "Read(T) 方法 ")

[ReadAllKeys 方法](../html/53c48696-f991-f734-12b0-1ad306a84c6b.htm "ReadAllKeys 方法 ")

[ReadAllKeysAsync 方法](../html/84c1be8c-3286-0410-0bd1-bb0185ff197e.htm "ReadAllKeysAsync 方法 ")

[ReadAndWriteKey 方法](../html/a6908f0e-1991-5210-f6c6-1e8b06369e56.htm "ReadAndWriteKey 方法 ")

[ReadAndWriteKeyAsync 方法](../html/6686726f-93ba-6435-76b0-b501c31e637e.htm "ReadAndWriteKeyAsync 方法 ")

[ReadAsync(T) 方法](../html/43a59914-1502-1353-0bf3-20281ac92ca9.htm "ReadAsync(T) 方法 ")

[ReadCustomer 方法](../html/a4b7f43e-34ba-8aac-a705-a9bc8d4a5298.htm "ReadCustomer 方法 ")

[ReadCustomerAsync 方法](../html/bafea5ed-4293-ae60-c8e7-bb3d3481db44.htm "ReadCustomerAsync 方法 ")

[ReadFromCoreServer 方法](../html/a5985ec7-e909-5a40-b2b9-7160363a1156.htm "ReadFromCoreServer 方法 ")

[ReadFromCoreServerAsync 方法](../html/984a89de-20da-e7ed-e471-3e889955aa31.htm "ReadFromCoreServerAsync 方法 ")

[ReadHashKey 方法](../html/6e8fc531-b2d6-2b1e-4885-1b937745652f.htm "ReadHashKey 方法 ")

[ReadHashKeyAll 方法](../html/63828fbe-fbdf-772a-010b-1425b0584572.htm "ReadHashKeyAll 方法 ")

[ReadHashKeyAllAsync 方法](../html/31c54864-f7bc-b3ce-7e77-15f0c1dcc540.htm "ReadHashKeyAllAsync 方法 ")

[ReadHashKeyAsync 方法](../html/36c995b4-42a7-e264-8a03-3390e24376ec.htm "ReadHashKeyAsync 方法 ")

[ReadHashKeyLength 方法](../html/32d483da-49e0-a205-bdc4-a659711be1c0.htm "ReadHashKeyLength 方法 ")

[ReadHashKeyLengthAsync 方法](../html/5d0d4570-1f16-4277-cff7-a40607a530ab.htm "ReadHashKeyLengthAsync 方法 ")

[ReadHashKeys 方法](../html/420286f8-5876-ee1b-b328-0b6eba7cefac.htm "ReadHashKeys 方法 ")

[ReadHashKeysAsync 方法](../html/bbc8310c-029d-1ea5-b0e3-239b29ee8743.htm "ReadHashKeysAsync 方法 ")

[ReadHashValues 方法](../html/1d8b8e06-3d34-f4c6-68d0-20c14b309de1.htm "ReadHashValues 方法 ")

[ReadHashValuesAsync 方法](../html/403e7344-c665-def6-70d9-81e5b22a740b.htm "ReadHashValuesAsync 方法 ")

[ReadKey 方法](../html/93fa3dab-d22c-d028-00c3-e81391085e6f.htm "ReadKey 方法 ")

[ReadKeyAsync 方法](../html/48cee995-f908-c3a2-ab70-829d789f58a5.htm "ReadKeyAsync 方法 ")

[ReadKeyLength 方法](../html/7746357b-7e1b-d150-d019-722b752f5eeb.htm "ReadKeyLength 方法 ")

[ReadKeyLengthAsync 方法](../html/7dd1b305-77b1-920c-14c0-2fe38e970791.htm "ReadKeyLengthAsync 方法 ")

[ReadKeyRange 方法](../html/2e3e9582-cb24-44e6-d78b-7c672acdd4c8.htm "ReadKeyRange 方法 ")

[ReadKeyRangeAsync 方法](../html/407fdd5b-8b57-75d2-63be-93294311cb30.htm "ReadKeyRangeAsync 方法 ")

[ReadKeyTTL 方法](../html/679079c3-36e5-1618-3409-18459a74a0c4.htm "ReadKeyTTL 方法 ")

[ReadKeyTTLAsync 方法](../html/24887926-0fd5-b65f-d070-1d439a5273ef.htm "ReadKeyTTLAsync 方法 ")

[ReadKeyType 方法](../html/1f702770-7804-880a-038e-c002ed31c42f.htm "ReadKeyType 方法 ")

[ReadKeyTypeAsync 方法](../html/91c104ff-a646-76de-9bed-5e64237a57fc.htm "ReadKeyTypeAsync 方法 ")

[ReadListByIndex 方法](../html/43c7dae8-4c01-d2c9-b38b-0e7bef68b581.htm "ReadListByIndex 方法 ")

[ReadListByIndexAsync 方法](../html/1e1163ad-8161-9390-cff2-7919c8145c92.htm "ReadListByIndexAsync 方法 ")

[ReadRandomKey 方法](../html/0d358a8f-b9e0-c659-c218-f051aaedcbfe.htm "ReadRandomKey 方法 ")

[ReadRandomKeyAsync 方法](../html/b97f91d2-2ab6-a39b-86f7-71238cd44fac.htm "ReadRandomKeyAsync 方法 ")

[ReadServerTime 方法](../html/c6c5f1fa-4639-0a32-8a42-48d56985e8e6.htm "ReadServerTime 方法 ")

[ReadServerTimeAsync 方法](../html/53022c3b-766e-c31b-1b55-548ceb3b2f42.htm "ReadServerTimeAsync 方法 ")

[RenameKey 方法](../html/ecf533af-233c-16b8-f88d-cb232e115013.htm "RenameKey 方法 ")

[RenameKeyAsync 方法](../html/20073ebe-0c16-2f57-a96b-79246e01129d.htm "RenameKeyAsync 方法 ")

[Save 方法](../html/261a8daa-9890-95bc-f8e4-02b78eb86277.htm "Save 方法 ")

[SaveAsync 方法](../html/88c2043f-e9e6-9d95-4a2e-455452cf8adf.htm "SaveAsync 方法 ")

[SelectDB 方法](../html/584a38b8-96e0-736f-7a24-981a8e920574.htm "SelectDB 方法 ")

[SelectDBAsync 方法](../html/462ed998-9938-a261-bde4-579b4b9d3dc1.htm "SelectDBAsync 方法 ")

[SetAdd 方法](../html/6f52901b-4d7f-70ca-1bf1-2e7146604dba.htm "SetAdd 方法 ")

[SetAddAsync 方法](../html/a59f7c04-0b20-35a4-a431-53299cb829f0.htm "SetAddAsync 方法 ")

[SetCard 方法](../html/cc0a3d1d-019b-f0ba-90a5-1bab767d6f05.htm "SetCard 方法 ")

[SetCardAsync 方法](../html/0d0665a9-a174-d154-f95b-2bf10d885bc3.htm "SetCardAsync 方法 ")

[SetDiff 方法](../html/bcc43604-dc32-a27f-53a4-a89ab29916be.htm "SetDiff 方法 ")

[SetDiffAsync 方法](../html/b70eae65-3537-3ab7-4cf9-ae6c314884c9.htm "SetDiffAsync 方法 ")

[SetDiffStore 方法](../html/7c356901-c174-c81b-50b9-ce8d60258d1e.htm "SetDiffStore 方法 ")

[SetDiffStoreAsync 方法](../html/697ac488-c2e1-8c5a-5e5d-845f676c85fd.htm "SetDiffStoreAsync 方法 ")

[SetInter 方法](../html/fa3a123d-447d-ce8b-5670-e47956d427d1.htm "SetInter 方法 ")

[SetInterAsync 方法](../html/622854f7-5e07-092e-42eb-06cc63580467.htm "SetInterAsync 方法 ")

[SetInterStore 方法](../html/996ff5a7-0d0b-5b17-546b-eb006fa5773f.htm "SetInterStore 方法 ")

[SetInterStoreAsync 方法](../html/115a5c5d-2c7b-4c9a-daec-92477fb3044a.htm "SetInterStoreAsync 方法 ")

[SetIsMember 方法](../html/1c6d465b-353a-34f2-8f46-9d4037e58911.htm "SetIsMember 方法 ")

[SetIsMemberAsync 方法](../html/21419054-3273-2aab-d486-633ab8861039.htm "SetIsMemberAsync 方法 ")

[SetMembers 方法](../html/7a5b277f-1fe1-141f-346b-3722d97ac7ca.htm "SetMembers 方法 ")

[SetMembersAsync 方法](../html/ce86155f-4be2-3836-3594-1398d78183ad.htm "SetMembersAsync 方法 ")

[SetMove 方法](../html/8b76c71a-c755-708b-e1e7-61eda8b9456c.htm "SetMove 方法 ")

[SetMoveAsync 方法](../html/dc01893c-b09e-e134-49d3-76d0dbc308b9.htm "SetMoveAsync 方法 ")

[SetPop 方法](../html/8375117d-d7b5-1f33-cb59-80cc632b736d.htm "SetPop 方法 ")

[SetPopAsync 方法](../html/66039533-6b1e-ad66-e89b-ef8f7d46c47b.htm "SetPopAsync 方法 ")

[SetRandomMember 方法](../html/2769a6d4-3b4c-75b2-c57b-abe19bb2f8e9.htm "SetRandomMember 方法 ")

[SetRandomMemberAsync 方法](../html/fb4fda18-00e6-2545-d037-6585a9487334.htm "SetRandomMemberAsync 方法 ")

[SetRemove 方法](../html/a0ba101e-935e-aca6-da5e-bb6187b42f02.htm "SetRemove 方法 ")

[SetRemoveAsync 方法](../html/f80e5db7-218b-c4bf-0c57-ec2555f6ec4e.htm "SetRemoveAsync 方法 ")

[SetUnion 方法](../html/c255e92b-1a91-f3e4-49cd-79f39d7a1079.htm "SetUnion 方法 ")

[SetUnionAsync 方法](../html/184ffefd-f6f5-b5d8-e1d6-45ae6d2d653f.htm "SetUnionAsync 方法 ")

[SetUnionStore 方法](../html/56394825-4bd1-6422-4f9b-1520ed05f927.htm "SetUnionStore 方法 ")

[SetUnionStoreAsync 方法](../html/99689f9d-ebcc-cc16-19e3-4c46f3016084.htm "SetUnionStoreAsync 方法 ")

[SubscribeMessage 方法](../html/bbaad51e-d35a-fb1f-4bb5-c001d47693be.htm "SubscribeMessage 方法 ")

[ToString 方法](../html/6ff40e11-12f4-8759-45a0-594ca3e4fcdd.htm "ToString 方法 ")

[UnSubscribeMessage 方法](../html/592327f8-7e50-4d48-c229-ca17a039b6a3.htm "UnSubscribeMessage 方法 ")

[Write(T) 方法](../html/70a08f37-0f86-76f1-55a5-df1ddd5a7d16.htm "Write(T) 方法 ")

[WriteAndPublishKey 方法](../html/69029e5e-a78a-5c83-0300-b6dd057f9308.htm "WriteAndPublishKey 方法 ")

[WriteAndPublishKeyAsync 方法](../html/d44abf4b-02bf-cfa5-cbc8-0b36edab4521.htm "WriteAndPublishKeyAsync 方法 ")

[WriteAsync(T) 方法](../html/9d54e352-2bd0-2655-3e14-49a0213e6040.htm "WriteAsync(T) 方法 ")

[WriteExpireKey 方法](../html/816ca8c0-5e3a-4609-acc3-6f227ac2dfde.htm "WriteExpireKey 方法 ")

[WriteExpireKeyAsync 方法](../html/f1dce55a-5681-32f6-2c7c-cc664465490c.htm "WriteExpireKeyAsync 方法 ")

[WriteHashKey 方法](../html/a7c80b1f-ec11-c8f4-dae2-0fa0acada685.htm "WriteHashKey 方法 ")

[WriteHashKeyAsync 方法](../html/20122681-134d-f220-683c-11c6c9201ac9.htm "WriteHashKeyAsync 方法 ")

[WriteHashKeyNx 方法](../html/819ddbb0-74d4-4667-6e55-372fac166c68.htm "WriteHashKeyNx 方法 ")

[WriteHashKeyNxAsync 方法](../html/b4571eb9-be29-fbb8-691b-f3a93d9245ec.htm "WriteHashKeyNxAsync 方法 ")

[WriteKey 方法](../html/12919427-9ab8-d765-27cc-0c3ee2c481d2.htm "WriteKey 方法 ")

[WriteKeyAsync 方法](../html/b136bcbb-3813-f660-54f1-8c5ccefb5ddf.htm "WriteKeyAsync 方法 ")

[WriteKeyIfNotExists 方法](../html/9d48312b-2ce1-077f-bf47-3e6b5a328c45.htm "WriteKeyIfNotExists 方法 ")

[WriteKeyIfNotExistsAsync 方法](../html/e6844e26-ba17-a0a8-0ff2-bb1a1e8bce8e.htm "WriteKeyIfNotExistsAsync 方法 ")

[WriteKeyRange 方法](../html/c1f98cc9-f791-4a35-bf67-ea29e8930d32.htm "WriteKeyRange 方法 ")

[WriteKeyRangeAsync 方法](../html/cdfaffec-518b-2754-5a10-3c5b68ee1767.htm "WriteKeyRangeAsync 方法 ")

[ZSetAdd 方法](../html/9c15cbbe-df53-5280-c9e4-b9d020222e0c.htm "ZSetAdd 方法 ")

[ZSetAddAsync 方法](../html/5318d2e6-d635-ed62-18e5-6573a530ab3e.htm "ZSetAddAsync 方法 ")

[ZSetCard 方法](../html/aa928601-c5f6-63fb-6bd8-e2e383796043.htm "ZSetCard 方法 ")

[ZSetCardAsync 方法](../html/b760ff15-4c7e-510f-1927-e17e6310e142.htm "ZSetCardAsync 方法 ")

[ZSetCount 方法](../html/9872ddc0-8b8d-d9d6-4d08-64001fdefd29.htm "ZSetCount 方法 ")

[ZSetCountAsync 方法](../html/b294e1ee-27c7-7326-eb4b-8669040aecf0.htm "ZSetCountAsync 方法 ")

[ZSetIncreaseBy 方法](../html/c84e0c14-15cd-451a-dec4-2a112a73772e.htm "ZSetIncreaseBy 方法 ")

[ZSetIncreaseByAsync 方法](../html/9805ac6b-aa05-2152-7dbb-7f8077ea462b.htm "ZSetIncreaseByAsync 方法 ")

[ZSetRange 方法](../html/759e2523-e2ad-0dc1-a904-354031a8ae50.htm "ZSetRange 方法 ")

[ZSetRangeAsync 方法](../html/7724eb62-2653-65a7-779d-391d630c2faf.htm "ZSetRangeAsync 方法 ")

[ZSetRangeByScore 方法](../html/acab27d0-6b6c-1f25-00ad-fafd4ac068fa.htm "ZSetRangeByScore 方法 ")

[ZSetRangeByScoreAsync 方法](../html/d1e0aa3c-d5f9-caeb-4017-699f956e661d.htm "ZSetRangeByScoreAsync 方法 ")

[ZSetRank 方法](../html/7957411f-32c1-4e07-a4f2-9d1948c63a2b.htm "ZSetRank 方法 ")

[ZSetRankAsync 方法](../html/1181fd8c-02aa-9d53-cbeb-35d542bed241.htm "ZSetRankAsync 方法 ")

[ZSetRemove 方法](../html/a76625ea-7a7b-161d-a039-02f4f5a8e997.htm "ZSetRemove 方法 ")

[ZSetRemoveAsync 方法](../html/2f1b3de6-ad74-e2b7-3350-79d9afe13c9d.htm "ZSetRemoveAsync 方法 ")

[ZSetRemoveRangeByRank 方法](../html/cb81dfc2-6e30-0aca-8bcb-348cdc32bc35.htm "ZSetRemoveRangeByRank 方法 ")

[ZSetRemoveRangeByRankAsync 方法](../html/9bde4b5f-b9bd-e01b-feb4-82c9b532dee4.htm "ZSetRemoveRangeByRankAsync 方法 ")

[ZSetRemoveRangeByScore 方法](../html/ac7be899-fe4e-3f12-d60f-81a038a927c7.htm "ZSetRemoveRangeByScore 方法 ")

[ZSetRemoveRangeByScoreAsync 方法](../html/6ccee6bf-8635-14bf-2e7f-d1ac351c70c9.htm "ZSetRemoveRangeByScoreAsync 方法 ")

[ZSetReverseRange 方法](../html/ef6defc7-62a2-c8a4-4025-32f6b3a82ae4.htm "ZSetReverseRange 方法 ")

[ZSetReverseRangeAsync 方法](../html/b40d37b7-1118-6aec-c7da-85ec54b457d0.htm "ZSetReverseRangeAsync 方法 ")

[ZSetReverseRangeByScore 方法](../html/3c1a26a8-f8e5-a835-7051-94b4d9c851ab.htm "ZSetReverseRangeByScore 方法 ")

[ZSetReverseRangeByScoreAsync 方法](../html/3a3cfbff-cee7-e1f3-d9c1-e7528db4e998.htm "ZSetReverseRangeByScoreAsync 方法 ")

[ZSetReverseRank 方法](../html/d524195b-6b02-a8e8-330a-3454a6ed21cd.htm "ZSetReverseRank 方法 ")

[ZSetReverseRankAsync 方法](../html/c010781c-2851-7158-f953-1779929e3242.htm "ZSetReverseRankAsync 方法 ")

[ZSetScore 方法](../html/a98ee7e4-f361-f403-975b-a197e0ca0268.htm "ZSetScore 方法 ")

[ZSetScoreAsync 方法](../html/fd24c433-944e-2ecb-4bcd-266bfc6825d4.htm "ZSetScoreAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| RedisClient 方法 |

[RedisClient](e33d9088-6df0-32e9-33f1-426e944c3b7f.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 受保护的方法 | [AccountCertificate](e9ef8320-5daf-4629-dcc6-4076e7944a9a.htm) | 认证账号，根据已经设置的用户名和密码，进行发送服务器进行账号认证。  Authentication account, according to the user name and password that have been set, sending server for account authentication. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 受保护的方法 | [AccountCertificateAsync](2dc6eb5a-79c8-1e51-7aa8-c3d29769033a.htm) | 认证账号，根据已经设置的用户名和密码，进行发送服务器进行账号认证。  Authentication account, according to the user name and password that have been set, sending server for account authentication. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | [AppendKey](f2521b8c-0fb1-6b34-4adc-ca6804c7195f.htm) | 如果 key 已经存在并且是一个字符串， APPEND 命令将 value 追加到 key 原来的值的末尾。 如果 key 不存在， APPEND 就简单地将给定 key 设为 value ，就像执行 SET key value 一样。 返回追加 value 之后， key 中字符串的长度。 |
| 公共方法 | [AppendKeyAsync](c39b4c6a-b7f6-68fe-dfc9-7cce1b653e24.htm) | 如果 key 已经存在并且是一个字符串， APPEND 命令将 value 追加到 key 原来的值的末尾。 如果 key 不存在， APPEND 就简单地将给定 key 设为 value ，就像执行 SET key value 一样。 返回追加 value 之后， key 中字符串的长度。 |
| 公共方法 | [ChangePassword](78a37d4d-bdee-7a09-d0ce-dc245627b1d1.htm) | 修改Redis的密码信息，如果不需要密码，则传入空字符串即可 |
| 公共方法 | [ChangePasswordAsync](58f6f736-093d-6e4d-3ad3-19590fa2a3a3.htm) | 修改Redis的密码信息，如果不需要密码，则传入空字符串即可 |
| 受保护的方法 | [CheckReceiveDataComplete](ea6ac6fd-de29-39d0-1098-6d90d5a1fb5d.htm) | 检查当前从网口接收的数据是否是完整的，如果是完整的，则需要返回 True，表示数据接收立即完成，默认返回 True  Check whether the data currently received from the network port is complete, and if it is complete, you need to return True, indicating that the data reception is completed immediately, and the default value is True (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 受保护的方法 | [CheckRemoteToken](1d231984-101b-0fff-0daf-1d18581a4241.htm) | 检查当前的头子节信息的令牌是否是正确的，仅用于某些特殊的协议实现  Check whether the token of the current header subsection information is correct, only for some special protocol implementations (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 公共方法代码示例 | [ConnectClose](e235581f-1c77-3b52-6f49-a6d1c79e559a.htm) | 手动断开与远程服务器的连接，如果当前是长连接模式，那么就会切换到短连接模式  Manually disconnect from the remote server, if it is currently in long connection mode, it will switch to short connection mode (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法代码示例 | [ConnectCloseAsync](ecec49bf-ad46-e74b-69ce-211464eaea6d.htm) | 手动断开与远程服务器的连接，如果当前是长连接模式，那么就会切换到短连接模式  Manually disconnect from the remote server, if it is currently in long connection mode, it will switch to short connection mode (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法代码示例 | [ConnectServer](ea436d31-7950-42df-a9c0-3c749c47e31d.htm) | 尝试连接远程的服务器，如果连接成功，就切换短连接模式到长连接模式，后面的每次请求都共享一个通道，使得通讯速度更快速  Try to connect to a remote server. If the connection is successful, switch the short connection mode to the long connection mode. Each subsequent request will share a channel, making the communication speed faster. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法代码示例 | [ConnectServer(AlienSession)](a5e9ec5a-548c-3d66-2509-fafa5ec21c3b.htm) | 使用指定的套接字创建异形客户端，在异形客户端的模式下，网络通道需要被动创建。  Use the specified socket to create the alien client. In the alien client mode, the network channel needs to be created passively. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | [ConnectServer(MqttClient, String, String)](3ae46111-33b7-c7dc-02d8-f689c1084b5a.htm) | 使用一个MQTT中转服务器来连接设备对象，并进行相关的读取操作 (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法代码示例 | [ConnectServerAsync](7eeb4d56-e784-e48e-b49a-8498a8fe8384.htm) | 尝试连接远程的服务器，如果连接成功，就切换短连接模式到长连接模式，后面的每次请求都共享一个通道，使得通讯速度更快速  Try to connect to a remote server. If the connection is successful, switch the short connection mode to the long connection mode. Each subsequent request will share a channel, making the communication speed faster. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 受保护的方法代码示例 | [CreateSocketAndConnect(String, Int32)](a153ae27-92e3-4fe7-b2e7-c8d207bebbb0.htm) | 创建一个新的socket对象并连接到远程的地址，默认超时时间为10秒钟，需要指定ip地址以及端口号信息  Create a new socket object and connect to the remote address. The default timeout is 10 seconds. You need to specify the IP address and port number. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法代码示例 | [CreateSocketAndConnect(IPEndPoint, Int32, IPEndPoint)](7ce6288f-c6db-e968-d676-f851d2b05c77.htm) | 创建一个新的socket对象并连接到远程的地址，需要指定远程终结点，超时时间（单位是毫秒），如果需要绑定本地的IP或是端口，传入 local对象  To create a new socket object and connect to the remote address, you need to specify the remote endpoint, the timeout period (in milliseconds), if you need to bind the local IP or port, pass in the local object (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法代码示例 | [CreateSocketAndConnect(String, Int32, Int32)](d0e14a5f-1e97-e59d-078f-4fbf0412b865.htm) | 创建一个新的socket对象并连接到远程的地址，需要指定ip地址以及端口号信息，还有超时时间，单位是毫秒  To create a new socket object and connect to a remote address, you need to specify the IP address and port number information, and the timeout period in milliseconds (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法代码示例 | [CreateSocketAndConnectAsync(String, Int32)](0227bc60-1c7f-3c1d-0cd6-e49bb904d6df.htm) | 创建一个新的socket对象并连接到远程的地址，默认超时时间为10秒钟，需要指定ip地址以及端口号信息  Create a new socket object and connect to the remote address. The default timeout is 10 seconds. You need to specify the IP address and port number. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法代码示例 | [CreateSocketAndConnectAsync(IPEndPoint, Int32, IPEndPoint)](ccec07bf-cfa4-b302-210e-441821e095a4.htm) | 创建一个新的socket对象并连接到远程的地址，需要指定远程终结点，超时时间（单位是毫秒），如果需要绑定本地的IP或是端口，传入 local对象  To create a new socket object and connect to the remote address, you need to specify the remote endpoint, the timeout period (in milliseconds), if you need to bind the local IP or port, pass in the local object (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法代码示例 | [CreateSocketAndConnectAsync(String, Int32, Int32)](7ec25883-60e2-c7ac-2f90-0ea1ae0c4945.htm) | 创建一个新的socket对象并连接到远程的地址，需要指定ip地址以及端口号信息，还有超时时间，单位是毫秒  To create a new socket object and connect to a remote address, you need to specify the IP address and port number information, and the timeout period in milliseconds (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 公共方法 | [DBSize](5814129a-9a73-25ac-e726-fe1da141eb28.htm) | 返回当前数据库的 key 的数量。 |
| 公共方法 | [DBSizeAsync](243a1f5f-a474-2e20-1298-06eff125f8b5.htm) | 返回当前数据库的 key 的数量。 |
| 受保护的方法 | [DecideWhetherQAMessage](77091dac-82d7-8c67-d691-31ae5b2aaee4.htm) | 决定当前的消息是否是应答机制的消息内容，需要在客户端进行重写实现，如果是应答机制，返回 True, 否则返回 False  To determine whether the current message is the message content of the response mechanism, it needs to be rewritten on the client side. If it is the response mechanism, return True, otherwise return False (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | [DecrementKey(String)](0f2a8a0d-9e7f-c6ab-2081-bba6606c37a9.htm) | 将 key 中储存的数字值减一。如果 key 不存在，那么 key 的值会先被初始化为 0 ，然后再执行 DECR 操作。 如果值包含错误的类型，或字符串类型的值不能表示为数字，那么返回一个错误。 本操作的值限制在 64 位(bit)有符号数字表示之内。 返回执行 DECR 命令之后 key 的值。 |
| 公共方法 | [DecrementKey(String, Int64)](0bb5b243-e3e3-064c-1c76-219ad3e2768a.htm) | 将 key 所储存的值减去减量 decrement 。如果 key 不存在，那么 key 的值会先被初始化为 0 ，然后再执行 DECR 操作。 如果值包含错误的类型，或字符串类型的值不能表示为数字，那么返回一个错误。 本操作的值限制在 64 位(bit)有符号数字表示之内。 返回减去 decrement 之后， key 的值。 |
| 公共方法 | [DecrementKeyAsync(String)](d0ade959-e1b8-cd2c-6a58-7104767acb23.htm) | 将 key 中储存的数字值减一。如果 key 不存在，那么 key 的值会先被初始化为 0 ，然后再执行 DECR 操作。 如果值包含错误的类型，或字符串类型的值不能表示为数字，那么返回一个错误。 本操作的值限制在 64 位(bit)有符号数字表示之内。 返回执行 DECR 命令之后 key 的值。 |
| 公共方法 | [DecrementKeyAsync(String, Int64)](743ce5f1-4a4e-314d-424d-a739f9e2384a.htm) | 将 key 所储存的值减去减量 decrement 。如果 key 不存在，那么 key 的值会先被初始化为 0 ，然后再执行 DECR 操作。 如果值包含错误的类型，或字符串类型的值不能表示为数字，那么返回一个错误。 本操作的值限制在 64 位(bit)有符号数字表示之内。 返回减去 decrement 之后， key 的值。 |
| 受保护的方法 | [DeleteFileByName](59afa4a4-3f3c-9831-26ee-c7c13a213bf7.htm) | 删除一个指定的文件，如果文件不存在，直接返回 True，如果文件存在则直接删除，删除成功返回 True，如果发生了异常，返回False  Delete a specified file, if the file does not exist, return True directly, if the file exists, delete it directly, if the deletion is successful, return True, if an exception occurs, return  False (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 公共方法 | [DeleteHashKey(String, String)](3f085ae2-ff3d-3c08-98af-5796aba0d975.htm) | 删除哈希表 key 中的一个或多个指定域，不存在的域将被忽略。 |
| 公共方法 | [DeleteHashKey(String, String)](eeed0572-14a6-8f5f-9e59-0df1636a2cb0.htm) | 删除哈希表 key 中的一个或多个指定域，不存在的域将被忽略。返回被成功移除的域的数量，不包括被忽略的域。 |
| 公共方法 | [DeleteHashKeyAsync(String, String)](b973444b-9bfd-30ce-6e59-006588402029.htm) | 删除哈希表 key 中的一个或多个指定域，不存在的域将被忽略。 |
| 公共方法 | [DeleteHashKeyAsync(String, String)](b1cdcd69-7004-177f-29bd-414a6aff71e2.htm) | 删除哈希表 key 中的一个或多个指定域，不存在的域将被忽略。返回被成功移除的域的数量，不包括被忽略的域。 |
| 公共方法 | [DeleteKey(String)](6eab9b1d-4556-b9c8-e7fa-2a5c13bd9940.htm) | 删除给定的一个或多个 key 。不存在的 key 会被忽略。 |
| 公共方法 | [DeleteKey(String)](d0fa44e9-f8c1-e302-d4f0-cd57917d2bf0.htm) | 删除给定的一个或多个 key 。不存在的 key 会被忽略。 |
| 公共方法 | [DeleteKeyAsync(String)](b760fe40-12aa-c0d2-5638-acb37ac81383.htm) | 删除给定的一个或多个 key 。不存在的 key 会被忽略。 |
| 公共方法 | [DeleteKeyAsync(String)](a157ea04-65d0-a685-6fbe-ff377c7d7e66.htm) | 删除给定的一个或多个 key 。不存在的 key 会被忽略。 |
| 公共方法 | [Dispose](8a86a200-fbe0-3b45-0fed-c551405c9b59.htm) | 释放当前的资源，如果调用了本方法，那么该对象再使用的时候，需要重新实例化。  Release the current resource. If this method is called, the object needs to be instantiated again when it is used again. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 受保护的方法 | [Dispose(Boolean)](32f03275-9ae0-525c-84bc-6e2515d11c12.htm) | 释放当前的资源，并自动关闭长连接，如果设置了的话 (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | Equals | (继承自 Object。) |
| 公共方法 | [ExistsHashKey](5b0bbeb9-4ab9-1f83-9e09-5f7227ac28e6.htm) | 查看哈希表 key 中，给定域 field 是否存在。如果哈希表含有给定域，返回 1 。 如果哈希表不含有给定域，或 key 不存在，返回 0 。 |
| 公共方法 | [ExistsHashKeyAsync](a80d408e-d621-a502-7829-6c7ea2a081f9.htm) | 查看哈希表 key 中，给定域 field 是否存在。如果哈希表含有给定域，返回 1 。 如果哈希表不含有给定域，或 key 不存在，返回 0 。 |
| 公共方法 | [ExistsKey](2b1b5b08-13ed-88c3-445d-d3613b48bc95.htm) | 检查给定 key 是否存在。若 key 存在，返回 1 ，否则返回 0 。 |
| 公共方法 | [ExistsKeyAsync](0d4357ed-5de3-1444-fce2-0e9be619d20f.htm) | 检查给定 key 是否存在。若 key 存在，返回 1 ，否则返回 0 。 |
| 公共方法 | [ExpireKey](0272044f-cc55-6b6a-42f0-6a0d89f37fe7.htm) | 为给定 key 设置生存时间，当 key 过期时(生存时间为 0 )，它会被自动删除。设置成功返回 1 。当 key 不存在或者不能为 key 设置生存时间时，返回 0 。 |
| 公共方法 | [ExpireKeyAsync](53c18b2d-2169-0519-602e-f65302f490e1.htm) | 为给定 key 设置生存时间，当 key 过期时(生存时间为 0 )，它会被自动删除。设置成功返回 1 。当 key 不存在或者不能为 key 设置生存时间时，返回 0 。 |
| 受保护的方法 | [ExtraAfterReadFromCoreServer](9f8357e9-dff4-1002-9368-0744cdedff29.htm) | 和服务器交互完成的时候调用的方法，可以根据读写结果进行一些额外的操作，具体的操作需要根据实际的需求来重写实现  The method called when the interaction with the server is completed can perform some additional operations based on the read and write results. The specific operations need to be rewritten according to actual needs. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 受保护的方法代码示例 | [ExtraOnDisconnect](49d9b263-02d8-e822-6a33-49f93737b98b.htm) | 根据实际的协议选择是否重写本方法，有些协议在断开连接之前，需要发送一些报文来关闭当前的网络通道  Select whether to rewrite this method according to the actual protocol. Some protocols need to send some packets to close the current network channel before disconnecting. (重写 [NetworkDoubleBaseExtraOnDisconnect(Socket)](b8569c75-5e55-cd55-5251-207465a8ba39.htm).) |
| 受保护的方法代码示例 | [ExtraOnDisconnectAsync](f681350d-2aad-40e0-0824-30fbd3fe57ff.htm) | 根据实际的协议选择是否重写本方法，有些协议在断开连接之前，需要发送一些报文来关闭当前的网络通道  Select whether to rewrite this method according to the actual protocol. Some protocols need to send some packets to close the current network channel before disconnecting. (重写 [NetworkDoubleBaseExtraOnDisconnectAsync(Socket)](86aa00d8-b2d6-ada9-5926-2b7a2237fe6b.htm).) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | [FlushDB](38c28c3f-63bf-f1a0-67bc-6547799257f9.htm) | 清空当前的数据库的key信息 |
| 公共方法 | [FlushDBAsync](40e3ccb2-9fe7-c0d7-88b9-95266a13c336.htm) | 清空当前的数据库的key信息 |
| 受保护的方法 | [GetAvailableSocket](e86a6fc1-eb91-8de9-5ff8-799b0ea3786b.htm) | 获取本次操作的可用的网络通道，如果是短连接，就重新生成一个新的网络通道，如果是长连接，就复用当前的网络通道。  Obtain the available network channels for this operation. If it is a short connection, a new network channel is regenerated. If it is a long connection, the current network channel is reused. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 受保护的方法 | [GetAvailableSocketAsync](40cf9ce7-a5e6-5431-2e4b-56e3b7c98346.htm) | 获取本次操作的可用的网络通道，如果是短连接，就重新生成一个新的网络通道，如果是长连接，就复用当前的网络通道。  Obtain the available network channels for this operation. If it is a short connection, a new network channel is regenerated. If it is a long connection, the current network channel is reused. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法 | [GetListLength](36330604-2557-4fc1-bdd7-fc3aff05e997.htm) | 返回列表 key 的长度。如果 key 不存在，则 key 被解释为一个空列表，返回 0 .如果 key 不是列表类型，返回一个错误。 |
| 公共方法 | [GetListLengthAsync](0e09431d-45b3-b539-a618-a147ae377f88.htm) | 返回列表 key 的长度。如果 key 不存在，则 key 被解释为一个空列表，返回 0 .如果 key 不是列表类型，返回一个错误。 |
| 受保护的方法 | [GetNewNetMessage](481a40e0-8539-5fc8-3e2b-c25bcc952697.htm) | 获取一个新的消息对象的方法，需要在继承类里面进行重写  The method to get a new message object needs to be overridden in the inheritance class (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | [GetPipeSocket](19904680-372e-e0a6-1b27-b4180e759668.htm) | 获取当前用于通信的管道信息  Get the current pipe information used for communication (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | GetType | (继承自 Object。) |
| 公共方法 | [IncrementHashKey(String, String, Int64)](81e32c12-a45b-ebe3-5e74-6357c50256a8.htm) | 为哈希表 key 中的域 field 的值加上增量 increment 。增量也可以为负数，相当于对给定域进行减法操作。 如果 key 不存在，一个新的哈希表被创建并执行 HINCRBY 命令。返回执行 HINCRBY 命令之后，哈希表 key 中域 field 的值。 |
| 公共方法 | [IncrementHashKey(String, String, Single)](421d1ddd-aea7-bd21-76ed-b437d990aab6.htm) | 为哈希表 key 中的域 field 的值加上增量 increment 。增量也可以为负数，相当于对给定域进行减法操作。 如果 key 不存在，一个新的哈希表被创建并执行 HINCRBY 命令。返回执行 HINCRBY 命令之后，哈希表 key 中域 field 的值。 |
| 公共方法 | [IncrementHashKeyAsync(String, String, Int64)](9d381989-878b-7e2f-519a-86c11af596c9.htm) | 为哈希表 key 中的域 field 的值加上增量 increment 。增量也可以为负数，相当于对给定域进行减法操作。 如果 key 不存在，一个新的哈希表被创建并执行 HINCRBY 命令。返回执行 HINCRBY 命令之后，哈希表 key 中域 field 的值。 |
| 公共方法 | [IncrementHashKeyAsync(String, String, Single)](adb31272-47d3-1aee-b224-447cdda8e9aa.htm) | 为哈希表 key 中的域 field 的值加上增量 increment 。增量也可以为负数，相当于对给定域进行减法操作。 如果 key 不存在，一个新的哈希表被创建并执行 HINCRBY 命令。返回执行 HINCRBY 命令之后，哈希表 key 中域 field 的值。 |
| 公共方法 | [IncrementKey(String)](86472b84-a7e5-63fa-7adf-769b1e5df83e.htm) | 将 key 中储存的数字值增一。如果 key 不存在，那么 key 的值会先被初始化为 0 ，然后再执行 INCR 操作。 如果值包含错误的类型，或字符串类型的值不能表示为数字，那么返回一个错误。 返回执行 INCR 命令之后 key 的值。 |
| 公共方法 | [IncrementKey(String, Int64)](91c1b3c1-c826-ae22-41ec-fd210c7fc5b7.htm) | 将 key 所储存的值加上增量 increment 。如果 key 不存在，那么 key 的值会先被初始化为 0 ，然后再执行 INCR 操作。 如果值包含错误的类型，或字符串类型的值不能表示为数字，那么返回一个错误。 |
| 公共方法 | [IncrementKey(String, Single)](438bd452-d082-3631-34b0-586d8fd2224a.htm) | 将 key 所储存的值加上增量 increment 。如果 key 不存在，那么 key 的值会先被初始化为 0 ，然后再执行 INCRBYFLOAT 操作。 如果命令执行成功，那么 key 的值会被更新为（执行加法之后的）新值，并且新值会以字符串的形式返回给调用者 |
| 公共方法 | [IncrementKeyAsync(String)](0134664a-49ce-5aad-bf5e-eea10440d453.htm) | 将 key 中储存的数字值增一。如果 key 不存在，那么 key 的值会先被初始化为 0 ，然后再执行 INCR 操作。 如果值包含错误的类型，或字符串类型的值不能表示为数字，那么返回一个错误。 返回执行 INCR 命令之后 key 的值。 |
| 公共方法 | [IncrementKeyAsync(String, Int64)](83d89534-28ac-6edd-be62-cbd721b9b058.htm) | 将 key 所储存的值加上增量 increment 。如果 key 不存在，那么 key 的值会先被初始化为 0 ，然后再执行 INCR 操作。 如果值包含错误的类型，或字符串类型的值不能表示为数字，那么返回一个错误。 |
| 公共方法 | [IncrementKeyAsync(String, Single)](5c2ec4f8-9dd6-a174-032f-26bed895e1e7.htm) | 将 key 所储存的值加上增量 increment 。如果 key 不存在，那么 key 的值会先被初始化为 0 ，然后再执行 INCRBYFLOAT 操作。 如果命令执行成功，那么 key 的值会被更新为（执行加法之后的）新值，并且新值会以字符串的形式返回给调用者 |
| 受保护的方法代码示例 | [InitializationOnConnect](72414fc8-b43a-3377-acdc-c59ae0628e11.htm) | 根据实际的协议选择是否重写本方法，有些协议在创建连接之后，需要进行一些初始化的信号握手，才能最终建立网络通道。  Whether to rewrite this method is based on the actual protocol. Some protocols require some initial signal handshake to establish a network channel after the connection is created. (重写 [NetworkDoubleBaseInitializationOnConnect(Socket)](40d5eb26-9bca-11cf-42a1-55af6ea5aebd.htm).) |
| 受保护的方法代码示例 | [InitializationOnConnectAsync](9156d090-6b1d-7381-c396-81491b4b6756.htm) | 根据实际的协议选择是否重写本方法，有些协议在创建连接之后，需要进行一些初始化的信号握手，才能最终建立网络通道。  Whether to rewrite this method is based on the actual protocol. Some protocols require some initial signal handshake to establish a network channel after the connection is created. (重写 [NetworkDoubleBaseInitializationOnConnectAsync(Socket)](14cac21b-4001-da7d-09a6-941eba868e14.htm).) |
| 公共方法 | [IpAddressPing](d9611e18-4033-de95-444e-b91ae0a81709.htm) | 对当前设备的IP地址进行PING的操作，返回PING的结果，正常来说，返回Success  PING the IP address of the current device and return the PING result. Normally, it returns Success (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | [ListInsertAfter](11b6eb28-b2ee-57a5-c724-a99df9101c71.htm) | 将值 value 插入到列表 key 当中，位于值 pivot 之后。 当 pivot 不存在于列表 key 时，不执行任何操作。 当 key 不存在时， key 被视为空列表，不执行任何操作。 如果 key 不是列表类型，返回一个错误。 |
| 公共方法 | [ListInsertAfterAsync](9defd0d6-5681-4c26-0973-14ecca4fc821.htm) | 将值 value 插入到列表 key 当中，位于值 pivot 之后。 当 pivot 不存在于列表 key 时，不执行任何操作。 当 key 不存在时， key 被视为空列表，不执行任何操作。 如果 key 不是列表类型，返回一个错误。 |
| 公共方法 | [ListInsertBefore](b19eed6b-2778-4b91-6f6b-e75cd9f1088e.htm) | 将值 value 插入到列表 key 当中，位于值 pivot 之前。 当 pivot 不存在于列表 key 时，不执行任何操作。 当 key 不存在时， key 被视为空列表，不执行任何操作。 如果 key 不是列表类型，返回一个错误。 |
| 公共方法 | [ListInsertBeforeAsync](0bf5390c-af9f-ecd6-1148-34a20abee73d.htm) | 将值 value 插入到列表 key 当中，位于值 pivot 之前。 当 pivot 不存在于列表 key 时，不执行任何操作。 当 key 不存在时， key 被视为空列表，不执行任何操作。 如果 key 不是列表类型，返回一个错误。 |
| 公共方法 | [ListLeftPop](8865ed9c-69dc-2f28-7f53-42e1c03d152b.htm) | 移除并返回列表 key 的头元素。列表的头元素。当 key 不存在时，返回 nil 。 |
| 公共方法 | [ListLeftPopAsync](f4d7f464-73f7-ff5d-0834-e18b5a91c994.htm) | 移除并返回列表 key 的头元素。列表的头元素。当 key 不存在时，返回 nil 。 |
| 公共方法 | [ListLeftPush(String, String)](74a8557b-7b50-f4fe-ba71-807683204706.htm) | 将一个或多个值 value 插入到列表 key 的表头，如果 key 不存在，一个空列表会被创建并执行 LPUSH 操作。当 key 存在但不是列表类型时，返回一个错误。返回执行 LPUSH 命令后，列表的长度。 |
| 公共方法 | [ListLeftPush(String, String)](dca179dc-794a-9936-b7a4-66338fefe209.htm) | 将一个或多个值 value 插入到列表 key 的表头，如果 key 不存在，一个空列表会被创建并执行 LPUSH 操作。当 key 存在但不是列表类型时，返回一个错误。返回执行 LPUSH 命令后，列表的长度。 |
| 公共方法 | [ListLeftPushAsync(String, String)](64d628e5-166c-ab96-0e90-b7d59d79e4cb.htm) | 将一个或多个值 value 插入到列表 key 的表头，如果 key 不存在，一个空列表会被创建并执行 LPUSH 操作。当 key 存在但不是列表类型时，返回一个错误。返回执行 LPUSH 命令后，列表的长度。 |
| 公共方法 | [ListLeftPushAsync(String, String)](6404ecc4-4c22-cb1a-11b4-4d8f4cdafe52.htm) | 将一个或多个值 value 插入到列表 key 的表头，如果 key 不存在，一个空列表会被创建并执行 LPUSH 操作。当 key 存在但不是列表类型时，返回一个错误。返回执行 LPUSH 命令后，列表的长度。 |
| 公共方法 | [ListLeftPushX](33093159-9864-a0dd-6a09-8bf134942031.htm) | 将值 value 插入到列表 key 的表头，当且仅当 key 存在并且是一个列表。和 LPUSH 命令相反，当 key 不存在时， LPUSHX 命令什么也不做。 返回LPUSHX 命令执行之后，表的长度。 |
| 公共方法 | [ListLeftPushXAsync](6bd2ab31-ff58-71f0-b68f-f36d7f53d845.htm) | 将值 value 插入到列表 key 的表头，当且仅当 key 存在并且是一个列表。和 LPUSH 命令相反，当 key 不存在时， LPUSHX 命令什么也不做。 返回LPUSHX 命令执行之后，表的长度。 |
| 公共方法 | [ListRange](efe06eb0-b945-fe65-5e81-21888cdda04f.htm) | 返回列表 key 中指定区间内的元素，区间以偏移量 start 和 stop 指定。 下标(index)参数 start 和 stop 都以 0 为底，也就是说，以 0 表示列表的第一个元素，以 1 表示列表的第二个元素，以此类推。 你也可以使用负数下标，以 -1 表示列表的最后一个元素， -2 表示列表的倒数第二个元素，以此类推。 返回一个列表，包含指定区间内的元素。 |
| 公共方法 | [ListRangeAsync](32e344f9-cdac-fc0e-2c50-bf2dd54b4f98.htm) | 返回列表 key 中指定区间内的元素，区间以偏移量 start 和 stop 指定。 下标(index)参数 start 和 stop 都以 0 为底，也就是说，以 0 表示列表的第一个元素，以 1 表示列表的第二个元素，以此类推。 你也可以使用负数下标，以 -1 表示列表的最后一个元素， -2 表示列表的倒数第二个元素，以此类推。 返回一个列表，包含指定区间内的元素。 |
| 公共方法 | [ListRemoveElementMatch](0ab40e89-d265-0b48-05de-5c1a7b41d326.htm) | 根据参数 count 的值，移除列表中与参数 value 相等的元素。count 的值可以是以下几种： count > 0 : 从表头开始向表尾搜索，移除与 value 相等的元素，数量为 count 。 count < 0 : 从表尾开始向表头搜索，移除与 value 相等的元素，数量为 count 的绝对值。 count = 0 : 移除表中所有与 value 相等的值。 返回被移除的数量。 |
| 公共方法 | [ListRemoveElementMatchAsync](5610de76-3198-5b5a-8f3d-b8fc937d7c83.htm) | 根据参数 count 的值，移除列表中与参数 value 相等的元素。count 的值可以是以下几种： count > 0 : 从表头开始向表尾搜索，移除与 value 相等的元素，数量为 count 。 count < 0 : 从表尾开始向表头搜索，移除与 value 相等的元素，数量为 count 的绝对值。 count = 0 : 移除表中所有与 value 相等的值。 返回被移除的数量。 |
| 公共方法 | [ListRightPop](97855b86-6a07-76e2-b887-1e830be705de.htm) | 移除并返回列表 key 的尾元素。当 key 不存在时，返回 nil 。 |
| 公共方法 | [ListRightPopAsync](fde8cee2-f393-a463-d561-74568545edc7.htm) | 移除并返回列表 key 的尾元素。当 key 不存在时，返回 nil 。 |
| 公共方法 | [ListRightPopLeftPush](e0b58cc4-4d80-64d6-729c-0dcdaafc4701.htm) | 命令 RPOPLPUSH 在一个原子时间内，执行以下两个动作：  1. 将列表 source 中的最后一个元素( 尾元素)弹出，并返回给客户端。  2. 将 source 弹出的元素插入到列表 destination ，作为 destination 列表的的头元素。   举个例子，你有两个列表 source 和 destination ， source 列表有元素 a, b, c ， destination 列表有元素 x, y, z ，执行 RPOPLPUSH source destination 之后， source 列表包含元素 a, b ， destination 列表包含元素 c, x, y, z ，并且元素 c 会被返回给客户端。 如果 source 不存在，值 nil 被返回，并且不执行其他动作。 如果 source 和 destination 相同，则列表中的表尾元素被移动到表头，并返回该元素，可以把这种特殊情况视作列表的旋转( rotation)操作。 |
| 公共方法 | [ListRightPopLeftPushAsync](41069df2-0dad-6f01-895f-a91086a92fa9.htm) | 命令 RPOPLPUSH 在一个原子时间内，执行以下两个动作：  1. 将列表 source 中的最后一个元素( 尾元素)弹出，并返回给客户端。  2. 将 source 弹出的元素插入到列表 destination ，作为 destination 列表的的头元素。   举个例子，你有两个列表 source 和 destination ， source 列表有元素 a, b, c ， destination 列表有元素 x, y, z ，执行 RPOPLPUSH source destination 之后， source 列表包含元素 a, b ， destination 列表包含元素 c, x, y, z ，并且元素 c 会被返回给客户端。 如果 source 不存在，值 nil 被返回，并且不执行其他动作。 如果 source 和 destination 相同，则列表中的表尾元素被移动到表头，并返回该元素，可以把这种特殊情况视作列表的旋转( rotation)操作。 |
| 公共方法 | [ListRightPush(String, String)](bec290e4-5da3-397c-5c53-848990279e78.htm) | 将一个或多个值 value 插入到列表 key 的表尾(最右边)。 如果 key 不存在，一个空列表会被创建并执行 RPUSH 操作。当 key 存在但不是列表类型时，返回一个错误。 |
| 公共方法 | [ListRightPush(String, String)](af196915-e868-386c-8d4a-0cc302f17a1a.htm) | 将一个或多个值 value 插入到列表 key 的表尾(最右边)。 如果有多个 value 值，那么各个 value 值按从左到右的顺序依次插入到表尾：比如对一个空列表 mylist 执行 RPUSH mylist a b c ，得出的结果列表为 a b c ， 如果 key 不存在，一个空列表会被创建并执行 RPUSH 操作。当 key 存在但不是列表类型时，返回一个错误。 返回执行 RPUSH 操作后，表的长度。 |
| 公共方法 | [ListRightPushAsync(String, String)](b30a73f7-59c3-e2f9-923e-19c6e62445d3.htm) | 将一个或多个值 value 插入到列表 key 的表尾(最右边)。 如果 key 不存在，一个空列表会被创建并执行 RPUSH 操作。当 key 存在但不是列表类型时，返回一个错误。 |
| 公共方法 | [ListRightPushAsync(String, String)](adf383fc-794b-3145-2320-e661aba47f95.htm) | 将一个或多个值 value 插入到列表 key 的表尾(最右边)。 如果有多个 value 值，那么各个 value 值按从左到右的顺序依次插入到表尾：比如对一个空列表 mylist 执行 RPUSH mylist a b c ，得出的结果列表为 a b c ， 如果 key 不存在，一个空列表会被创建并执行 RPUSH 操作。当 key 存在但不是列表类型时，返回一个错误。 返回执行 RPUSH 操作后，表的长度。 |
| 公共方法 | [ListRightPushX](e6ccfd1a-25ee-5eff-d56c-0ff9eb017fbc.htm) | 将值 value 插入到列表 key 的表尾，当且仅当 key 存在并且是一个列表。 和 RPUSH 命令相反，当 key 不存在时， RPUSHX 命令什么也不做。 |
| 公共方法 | [ListRightPushXAsync](3252db2a-56cf-1aa2-7716-f9eb2596ec62.htm) | 将值 value 插入到列表 key 的表尾，当且仅当 key 存在并且是一个列表。 和 RPUSH 命令相反，当 key 不存在时， RPUSHX 命令什么也不做。 |
| 公共方法 | [ListSet](5f2e7b8b-b0e7-87e6-dac7-81b572cf3fa3.htm) | 设置数组的某一个索引的数据信息，当 index 参数超出范围，或对一个空列表( key 不存在)进行 LSET 时，返回一个错误。 |
| 公共方法 | [ListSetAsync](ca019ec2-eb32-41f8-a460-89c70f90ace9.htm) | 设置数组的某一个索引的数据信息，当 index 参数超出范围，或对一个空列表( key 不存在)进行 LSET 时，返回一个错误。 |
| 公共方法 | [ListTrim](acf21364-1692-f9da-d337-c460b12a012d.htm) | 对一个列表进行修剪(trim)，就是说，让列表只保留指定区间内的元素，不在指定区间之内的元素都将被删除。 举个例子，执行命令 LTRIM list 0 2 ，表示只保留列表 list 的前三个元素，其余元素全部删除。 下标( index)参数 start 和 stop 都以 0 为底，也就是说，以 0 表示列表的第一个元素，以 1 表示列表的第二个元素，以此类推。 你也可以使用负数下标，以 -1 表示列表的最后一个元素， -2 表示列表的倒数第二个元素，以此类推。 当 key 不是列表类型时，返回一个错误。 |
| 公共方法 | [ListTrimAsync](9e86492d-205c-8427-9052-50ea3d8c1c64.htm) | 对一个列表进行修剪(trim)，就是说，让列表只保留指定区间内的元素，不在指定区间之内的元素都将被删除。 举个例子，执行命令 LTRIM list 0 2 ，表示只保留列表 list 的前三个元素，其余元素全部删除。 下标( index)参数 start 和 stop 都以 0 为底，也就是说，以 0 表示列表的第一个元素，以 1 表示列表的第二个元素，以此类推。 你也可以使用负数下标，以 -1 表示列表的最后一个元素， -2 表示列表的倒数第二个元素，以此类推。 当 key 不是列表类型时，返回一个错误。 |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法 | [MoveKey](5aa50043-85b9-a245-606a-de2b69164218.htm) | 将当前数据库的 key 移动到给定的数据库 db 当中。 如果当前数据库(源数据库)和给定数据库(目标数据库)有相同名字的给定 key ，或者 key 不存在于当前数据库，那么 MOVE 没有任何效果。 因此，也可以利用这一特性，将 MOVE 当作锁(locking)原语(primitive)。 |
| 公共方法 | [MoveKeyAsync](a93a74bf-65bd-733b-e39a-06d705252d99.htm) | 将当前数据库的 key 移动到给定的数据库 db 当中。 如果当前数据库(源数据库)和给定数据库(目标数据库)有相同名字的给定 key ，或者 key 不存在于当前数据库，那么 MOVE 没有任何效果。 因此，也可以利用这一特性，将 MOVE 当作锁(locking)原语(primitive)。 |
| 公共方法 | [OperateLongNumberFromServer](d2cdcd44-a6fb-7aa7-8c9c-08ff43107f65.htm) | 向服务器请求指令，并返回long数字的结果对象 |
| 公共方法 | [OperateLongNumberFromServerAsync](0005132d-e8c6-d51b-69dc-c9f8feaf8a9d.htm) | 向服务器请求指令，并返回long数字的结果对象 |
| 公共方法 | [OperateNumberFromServer](69aecc05-6901-f53a-da37-ebf85bae7053.htm) | 向服务器请求指定，并返回数字的结果对象 |
| 公共方法 | [OperateNumberFromServerAsync](b53bd6fa-d215-55f7-bbac-2ee171ed7498.htm) | 向服务器请求指定，并返回数字的结果对象 |
| 公共方法 | [OperateStatusFromServer](d460c3bb-e484-3fcb-e206-d20dc7201d63.htm) | 向服务器请求指令，并返回状态的结果对象，通常用于写入的判断，或是请求类型的判断 |
| 公共方法 | [OperateStatusFromServerAsync](e1a241a3-e200-747f-bb8b-5b7578a08999.htm) | 向服务器请求指令，并返回状态的结果对象，通常用于写入的判断，或是请求类型的判断 |
| 公共方法 | [OperateStringFromServer](822914af-cf62-1ccc-d22e-2df7dfac7294.htm) | 向服务器请求指令，并返回字符串的结果对象 |
| 公共方法 | [OperateStringFromServerAsync](eb31e52d-5387-eb14-291e-22f0e64666a4.htm) | 向服务器请求指令，并返回字符串的结果对象 |
| 公共方法 | [OperateStringsFromServer](c85c1a0e-83f9-69c9-a9b9-8a947ca54583.htm) | 向服务器请求指令，并返回字符串数组的结果对象 |
| 公共方法 | [OperateStringsFromServerAsync](2db50b51-32b8-1e92-cb9a-f2c847437e50.htm) | 向服务器请求指令，并返回字符串数组的结果对象 |
| 公共方法 | [PackCommandWithHeader](f96881d6-e59f-9e16-07e3-1188eee1a934.htm) | 对当前的命令进行打包处理，通常是携带命令头内容，标记当前的命令的长度信息，需要进行重写，否则默认不打包  The current command is packaged, usually carrying the content of the command header, marking the length of the current command, and it needs to be rewritten, otherwise it is not packaged by default (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | [PersistKey](c68a840b-ce64-3303-9339-a10f88bcbb77.htm) | 移除给定 key 的生存时间，将这个 key 从『易失的』(带生存时间 key )转换成『持久的』(一个不带生存时间、永不过期的 key )。 当生存时间移除成功时，返回 1 . 如果 key 不存在或 key 没有设置生存时间，返回 0 。 |
| 公共方法 | [PersistKeyAsync](a5d7ee67-9041-6455-7fc1-c6e96a442799.htm) | 移除给定 key 的生存时间，将这个 key 从『易失的』(带生存时间 key )转换成『持久的』(一个不带生存时间、永不过期的 key )。 当生存时间移除成功时，返回 1 . 如果 key 不存在或 key 没有设置生存时间，返回 0 。 |
| 公共方法 | [Ping](9cdd4e99-90a2-0ae8-141c-26f0f4a4d698.htm) | 向服务器进行PING的操作，服务器会返回PONG操作 |
| 公共方法 | [PingAsync](d6c8e625-edb2-7ff2-a54d-e852d49c73c8.htm) | 向服务器进行PING的操作，服务器会返回PONG操作 |
| 公共方法 | [Publish](74f20db5-aaaa-af2e-b4e6-5d192196b222.htm) | 将信息 message 发送到指定的频道 channel，返回接收到信息 message 的订阅者数量。 |
| 公共方法 | [PublishAsync](61428af8-94a2-b2b0-419b-5c0d3df4bee5.htm) | 将信息 message 发送到指定的频道 channel，返回接收到信息 message 的订阅者数量。 |
| 公共方法代码示例 | [ReadT](08365a89-e2d7-2aa5-3fd5-ab1a5d60c440.htm) | 从设备里读取支持Hsl特性的数据内容， 该特性为[HslRedisKeyAttribute](0b3a62cc-e2c8-dfb7-89de-7da47983c3c0.htm)，[HslRedisListItemAttribute](2bbc9c13-e0d2-6942-b529-1ef6ad142370.htm)， [HslRedisListAttribute](fd07c099-2c4b-18b4-45dc-3f869fdc50c3.htm)，[HslRedisHashFieldAttribute](13118df4-a00f-12bb-887d-e78f8b46236b.htm) 详细参考代码示例的操作说明。 |
| 公共方法 | [ReadAllKeys](53c48696-f991-f734-12b0-1ad306a84c6b.htm) | 查找所有符合给定模式 pattern 的 key 。 \* 匹配数据库中所有 key。 h?llo 匹配 hello ， hallo 和 hxllo 等。 h[ae]llo 匹配 hello 和 hallo ，但不匹配 hillo 。 |
| 公共方法 | [ReadAllKeysAsync](84c1be8c-3286-0410-0bd1-bb0185ff197e.htm) | 查找所有符合给定模式 pattern 的 key 。 \* 匹配数据库中所有 key。 h?llo 匹配 hello ， hallo 和 hxllo 等。 h[ae]llo 匹配 hello 和 hallo ，但不匹配 hillo 。 |
| 公共方法 | [ReadAndWriteKey](a6908f0e-1991-5210-f6c6-1e8b06369e56.htm) | 将给定 key 的值设为 value ，并返回 key 的旧值(old value)。当 key 存在但不是字符串类型时，返回一个错误。 |
| 公共方法 | [ReadAndWriteKeyAsync](6686726f-93ba-6435-76b0-b501c31e637e.htm) | 将给定 key 的值设为 value ，并返回 key 的旧值(old value)。当 key 存在但不是字符串类型时，返回一个错误。 |
| 公共方法代码示例 | [ReadAsyncT](43a59914-1502-1353-0bf3-20281ac92ca9.htm) | 从设备里读取支持Hsl特性的数据内容， 该特性为[HslRedisKeyAttribute](0b3a62cc-e2c8-dfb7-89de-7da47983c3c0.htm)，[HslRedisListItemAttribute](2bbc9c13-e0d2-6942-b529-1ef6ad142370.htm)， [HslRedisListAttribute](fd07c099-2c4b-18b4-45dc-3f869fdc50c3.htm)，[HslRedisHashFieldAttribute](13118df4-a00f-12bb-887d-e78f8b46236b.htm) 详细参考代码示例的操作说明。 |
| 公共方法 | [ReadCustomer](a4b7f43e-34ba-8aac-a705-a9bc8d4a5298.htm) | 自定义的指令交互方法，该指令用空格分割，举例：LTRIM AAAAA 0 999 就是收缩列表，GET AAA 就是获取键值，需要对返回的数据进行二次分析 |
| 公共方法 | [ReadCustomerAsync](bafea5ed-4293-ae60-c8e7-bb3d3481db44.htm) | 自定义的指令交互方法，该指令用空格分割，举例：LTRIM AAAAA 0 999 就是收缩列表，GET AAA 就是获取键值，需要对返回的数据进行二次分析 |
| 公共方法代码示例 | [ReadFromCoreServer(Byte)](b3e14431-d773-0d1a-f453-a3ad60648c29.htm) | 将数据发送到当前的网络通道中，并从网络通道中接收一个[INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm)指定的完整的报文，网络通道将根据[GetAvailableSocket](e86a6fc1-eb91-8de9-5ff8-799b0ea3786b.htm)方法自动获取，本方法是线程安全的。  Send data to the current network channel and receive a complete message specified by [INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm) from the network channel. The network channel will be automatically obtained according to the [GetAvailableSocket](e86a6fc1-eb91-8de9-5ff8-799b0ea3786b.htm) method This method is thread-safe. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | [ReadFromCoreServer(IEnumerableByte)](580dae22-fad5-5816-d58f-43068e89e419.htm) | 将多个数据报文按顺序发到设备，并从设备接收返回的数据内容，然后拼接成一个Byte[]信息，需要重写[UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm)方法才能返回正确的结果。  Send multiple data packets to the device in sequence, and receive the returned data content from the device, and then splicing them into a Byte[] message, you need to rewrite [UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm) method to return the correct result. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServer(Byte, Boolean, Boolean)](1a33d974-805d-f149-3665-eb8d8a99ea84.htm) | 将数据发送到当前的网络通道中，并从网络通道中接收一个[INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm)指定的完整的报文，网络通道将根据[GetAvailableSocket](e86a6fc1-eb91-8de9-5ff8-799b0ea3786b.htm)方法自动获取，本方法是线程安全的。  Send data to the current network channel and receive a complete message specified by [INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm) from the network channel. The network channel will be automatically obtained according to the [GetAvailableSocket](e86a6fc1-eb91-8de9-5ff8-799b0ea3786b.htm) method This method is thread-safe. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServer(Socket, Byte, Boolean, Boolean)](b79f1191-25bc-6cfc-3c43-bd25b03bbbba.htm) | 将数据报文发送指定的网络通道上，根据当前指定的[INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm)类型，返回一条完整的数据指令  Sends a data message to the specified network channel, and returns a complete data command according to the currently specified [INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm) type (重写 [NetworkDoubleBaseReadFromCoreServer(Socket, Byte, Boolean, Boolean)](3b72b289-b201-7548-0981-ac049c9a4135.htm).) |
| 公共方法代码示例 | [ReadFromCoreServerAsync(Byte)](ec51a18e-0f5d-dbe1-c229-d62e8760ac97.htm) | 将数据发送到当前的网络通道中，并从网络通道中接收一个[INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm)指定的完整的报文，网络通道将根据[GetAvailableSocket](e86a6fc1-eb91-8de9-5ff8-799b0ea3786b.htm)方法自动获取，本方法是线程安全的。  Send data to the current network channel and receive a complete message specified by [INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm) from the network channel. The network channel will be automatically obtained according to the [GetAvailableSocket](e86a6fc1-eb91-8de9-5ff8-799b0ea3786b.htm) method This method is thread-safe. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | [ReadFromCoreServerAsync(IEnumerableByte)](0fce97fd-8441-6b54-0877-b8dfacc738e3.htm) | 将多个数据报文按顺序发到设备，并从设备接收返回的数据内容，然后拼接成一个Byte[]信息，需要重写[UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm)方法才能返回正确的结果。  Send multiple data packets to the device in sequence, and receive the returned data content from the device, and then splicing them into a Byte[] message, you need to rewrite [UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm) method to return the correct result. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServerAsync(Byte, Boolean, Boolean)](7525b379-04c2-8dac-b451-a8b5413d6747.htm) | 将数据发送到当前的网络通道中，并从网络通道中接收一个[INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm)指定的完整的报文，网络通道将根据[GetAvailableSocket](e86a6fc1-eb91-8de9-5ff8-799b0ea3786b.htm)方法自动获取，本方法是线程安全的。  Send data to the current network channel and receive a complete message specified by [INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm) from the network channel. The network channel will be automatically obtained according to the [GetAvailableSocket](e86a6fc1-eb91-8de9-5ff8-799b0ea3786b.htm) method This method is thread-safe. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServerAsync(Socket, Byte, Boolean, Boolean)](6bde9f59-8b15-abba-9638-3e949612874d.htm) | 将数据报文发送指定的网络通道上，根据当前指定的[INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm)类型，返回一条完整的数据指令  Sends a data message to the specified network channel, and returns a complete data command according to the currently specified [INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm) type (重写 [NetworkDoubleBaseReadFromCoreServerAsync(Socket, Byte, Boolean, Boolean)](07db4a00-b548-8bd6-5de6-12deb91f5311.htm).) |
| 公共方法 | [ReadHashKey(String, String)](b0d43de4-1f62-fac9-970b-7e48f8a4f1fc.htm) | 返回哈希表 key 中给定域 field 的值。当给定域不存在或是给定 key 不存在时，返回 nil |
| 公共方法 | [ReadHashKey(String, String)](e7041162-0c8b-4345-d679-f92f7315d454.htm) | 返回哈希表 key 中，一个或多个给定域的值。如果给定的域不存在于哈希表，那么返回一个 nil 值。 因为不存在的 key 被当作一个空哈希表来处理，所以对一个不存在的 key 进行 HMGET 操作将返回一个只带有 nil 值的表。 |
| 公共方法 | [ReadHashKeyAll](63828fbe-fbdf-772a-010b-1425b0584572.htm) | 返回哈希表 key 中，所有的域和值。在返回值里，紧跟每个域名(field name)之后是域的值(value)，所以返回值的长度是哈希表大小的两倍。 |
| 公共方法 | [ReadHashKeyAllAsync](31c54864-f7bc-b3ce-7e77-15f0c1dcc540.htm) | 返回哈希表 key 中，所有的域和值。在返回值里，紧跟每个域名(field name)之后是域的值(value)，所以返回值的长度是哈希表大小的两倍。 |
| 公共方法 | [ReadHashKeyAsync(String, String)](406f96a9-706e-3ba2-ca7d-80238a8156c7.htm) | 返回哈希表 key 中给定域 field 的值。当给定域不存在或是给定 key 不存在时，返回 nil |
| 公共方法 | [ReadHashKeyAsync(String, String)](bc7ab677-80ae-4b7e-6609-96c10916faff.htm) | 返回哈希表 key 中，一个或多个给定域的值。如果给定的域不存在于哈希表，那么返回一个 nil 值。 因为不存在的 key 被当作一个空哈希表来处理，所以对一个不存在的 key 进行 HMGET 操作将返回一个只带有 nil 值的表。 |
| 公共方法 | [ReadHashKeyLength](32d483da-49e0-a205-bdc4-a659711be1c0.htm) | 返回哈希表 key 中域的数量。 |
| 公共方法 | [ReadHashKeyLengthAsync](5d0d4570-1f16-4277-cff7-a40607a530ab.htm) | 返回哈希表 key 中域的数量。 |
| 公共方法 | [ReadHashKeys](420286f8-5876-ee1b-b328-0b6eba7cefac.htm) | 返回哈希表 key 中的所有域。当 key 不存在时，返回一个空表。 |
| 公共方法 | [ReadHashKeysAsync](bbc8310c-029d-1ea5-b0e3-239b29ee8743.htm) | 返回哈希表 key 中的所有域。当 key 不存在时，返回一个空表。 |
| 公共方法 | [ReadHashValues](1d8b8e06-3d34-f4c6-68d0-20c14b309de1.htm) | 返回哈希表 key 中所有域的值。当 key 不存在时，返回一个空表。 |
| 公共方法 | [ReadHashValuesAsync](403e7344-c665-def6-70d9-81e5b22a740b.htm) | 返回哈希表 key 中所有域的值。当 key 不存在时，返回一个空表。 |
| 公共方法 | [ReadKey(String)](ea47037b-d2bf-4b7b-e28d-555e2e0f373a.htm) | 返回 key 所关联的字符串值。如果 key 不存在那么返回特殊值 nil 。 假如 key 储存的值不是字符串类型，返回一个错误，因为 GET 只能用于处理字符串值。 |
| 公共方法 | [ReadKey(String)](ac5af648-f537-7c6c-a064-281d1b07052b.htm) | 返回所有(一个或多个)给定 key 的值。 如果给定的 key 里面，有某个 key 不存在，那么这个 key 返回特殊值 null 。因此，该命令永不失败。 |
| 公共方法 | [ReadKeyAsync(String)](1159cbb5-7384-a9f8-a1dc-3a697345a482.htm) | 返回 key 所关联的字符串值。如果 key 不存在那么返回特殊值 nil 。 假如 key 储存的值不是字符串类型，返回一个错误，因为 GET 只能用于处理字符串值。 |
| 公共方法 | [ReadKeyAsync(String)](ae8c05c5-8f24-dc12-0dae-2083c4ed4c04.htm) | 返回所有(一个或多个)给定 key 的值。 如果给定的 key 里面，有某个 key 不存在，那么这个 key 返回特殊值 null 。因此，该命令永不失败。 |
| 公共方法 | [ReadKeyLength](7746357b-7e1b-d150-d019-722b752f5eeb.htm) | 返回 key 所储存的字符串值的长度。当 key 储存的不是字符串值时，返回一个错误。返回符串值的长度。当 key 不存在时，返回 0 。 |
| 公共方法 | [ReadKeyLengthAsync](7dd1b305-77b1-920c-14c0-2fe38e970791.htm) | 返回 key 所储存的字符串值的长度。当 key 储存的不是字符串值时，返回一个错误。返回符串值的长度。当 key 不存在时，返回 0 。 |
| 公共方法 | [ReadKeyRange](2e3e9582-cb24-44e6-d78b-7c672acdd4c8.htm) | 返回 key 中字符串值的子字符串，字符串的截取范围由 start 和 end 两个偏移量决定(包括 start 和 end 在内)。 负数偏移量表示从字符串最后开始计数， -1 表示最后一个字符， -2 表示倒数第二个，以此类推。 返回截取得出的子字符串。 |
| 公共方法 | [ReadKeyRangeAsync](407fdd5b-8b57-75d2-63be-93294311cb30.htm) | 返回 key 中字符串值的子字符串，字符串的截取范围由 start 和 end 两个偏移量决定(包括 start 和 end 在内)。 负数偏移量表示从字符串最后开始计数， -1 表示最后一个字符， -2 表示倒数第二个，以此类推。 返回截取得出的子字符串。 |
| 公共方法 | [ReadKeyTTL](679079c3-36e5-1618-3409-18459a74a0c4.htm) | 以秒为单位，返回给定 key 的剩余生存时间(TTL, time to live)。 |
| 公共方法 | [ReadKeyTTLAsync](24887926-0fd5-b65f-d070-1d439a5273ef.htm) | 以秒为单位，返回给定 key 的剩余生存时间(TTL, time to live)。 |
| 公共方法 | [ReadKeyType](1f702770-7804-880a-038e-c002ed31c42f.htm) | 返回 key 所储存的值的类型。none (key不存在)，string (字符串)，list (列表)，set (集合)，zset (有序集)，hash (哈希表) |
| 公共方法 | [ReadKeyTypeAsync](91c104ff-a646-76de-9bed-5e64237a57fc.htm) | 返回 key 所储存的值的类型。none (key不存在)，string (字符串)，list (列表)，set (集合)，zset (有序集)，hash (哈希表) |
| 公共方法 | [ReadListByIndex](43c7dae8-4c01-d2c9-b38b-0e7bef68b581.htm) | 返回列表 key 中，下标为 index 的元素。下标(index)参数 start 和 stop 都以 0 为底，也就是说，以 0 表示列表的第一个元素，以 1 表示列表的第二个元素，以此类推。 你也可以使用负数下标，以 -1 表示列表的最后一个元素， -2 表示列表的倒数第二个元素，以此类推。如果 key 不是列表类型，返回一个错误。 |
| 公共方法 | [ReadListByIndexAsync](1e1163ad-8161-9390-cff2-7919c8145c92.htm) | 返回列表 key 中，下标为 index 的元素。下标(index)参数 start 和 stop 都以 0 为底，也就是说，以 0 表示列表的第一个元素，以 1 表示列表的第二个元素，以此类推。 你也可以使用负数下标，以 -1 表示列表的最后一个元素， -2 表示列表的倒数第二个元素，以此类推。如果 key 不是列表类型，返回一个错误。 |
| 公共方法 | [ReadRandomKey](0d358a8f-b9e0-c659-c218-f051aaedcbfe.htm) | 从当前数据库中随机返回(不删除)一个 key 。 当数据库不为空时，返回一个 key 。 当数据库为空时，返回 nil 。 |
| 公共方法 | [ReadRandomKeyAsync](b97f91d2-2ab6-a39b-86f7-71238cd44fac.htm) | 从当前数据库中随机返回(不删除)一个 key 。 当数据库不为空时，返回一个 key 。 当数据库为空时，返回 nil 。 |
| 公共方法 | [ReadServerTime](c6c5f1fa-4639-0a32-8a42-48d56985e8e6.htm) | 获取服务器的时间戳信息，可用于本地时间的数据同步问题 |
| 公共方法 | [ReadServerTimeAsync](53022c3b-766e-c31b-1b55-548ceb3b2f42.htm) | 获取服务器的时间戳信息，可用于本地时间的数据同步问题 |
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
| 公共方法 | [RenameKey](ecf533af-233c-16b8-f88d-cb232e115013.htm) | 将 key 改名为 newkey 。 当 key 和 newkey 相同，或者 key 不存在时，返回一个错误。 当 newkey 已经存在时， RENAME 命令将覆盖旧值。 |
| 公共方法 | [RenameKeyAsync](20073ebe-0c16-2f57-a96b-79246e01129d.htm) | 将 key 改名为 newkey 。 当 key 和 newkey 相同，或者 key 不存在时，返回一个错误。 当 newkey 已经存在时， RENAME 命令将覆盖旧值。 |
| 公共方法 | [Save](261a8daa-9890-95bc-f8e4-02b78eb86277.htm) | SAVE 命令执行一个同步保存操作，将当前 Redis 实例的所有数据快照(snapshot)以 RDB 文件的形式保存到硬盘。 |
| 公共方法 | [SaveAsync](88c2043f-e9e6-9d95-4a2e-455452cf8adf.htm) | 在后台异步(Asynchronously)保存当前数据库的数据到磁盘。 BGSAVE 命令执行之后立即返回 OK ，然后 Redis fork 出一个新子进程，原来的 Redis 进程(父进程)继续处理客户端请求，而子进程则负责将数据保存到磁盘，然后退出。 |
| 公共方法 | [SelectDB](584a38b8-96e0-736f-7a24-981a8e920574.htm) | 切换到指定的数据库，数据库索引号 index 用数字值指定，以 0 作为起始索引值。默认使用 0 号数据库。 |
| 公共方法 | [SelectDBAsync](462ed998-9938-a261-bde4-579b4b9d3dc1.htm) | 切换到指定的数据库，数据库索引号 index 用数字值指定，以 0 作为起始索引值。默认使用 0 号数据库。 |
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
| 公共方法 | [SetAdd(String, String)](061c607b-58ff-1585-f1b0-be2f6950dc8d.htm) | 将一个member 元素加入到集合 key 当中，已经存在于集合的 member 元素将被忽略。假如 key 不存在，则创建一个只包含 member 元素作成员的集合。当 key 不是集合类型时，返回一个错误。 |
| 公共方法 | [SetAdd(String, String)](687ec76f-50a9-5ecb-7732-47a4944373d4.htm) | 将一个或多个 member 元素加入到集合 key 当中，已经存在于集合的 member 元素将被忽略。假如 key 不存在，则创建一个只包含 member 元素作成员的集合。当 key 不是集合类型时，返回一个错误。 |
| 公共方法 | [SetAddAsync(String, String)](4fc5a399-e408-8cf9-11fe-10cb0e718cb6.htm) | 将一个member 元素加入到集合 key 当中，已经存在于集合的 member 元素将被忽略。假如 key 不存在，则创建一个只包含 member 元素作成员的集合。当 key 不是集合类型时，返回一个错误。 |
| 公共方法 | [SetAddAsync(String, String)](ee87b3df-3919-93c3-7a7e-9dfda289d0f5.htm) | 将一个或多个 member 元素加入到集合 key 当中，已经存在于集合的 member 元素将被忽略。假如 key 不存在，则创建一个只包含 member 元素作成员的集合。当 key 不是集合类型时，返回一个错误。 |
| 公共方法 | [SetCard](cc0a3d1d-019b-f0ba-90a5-1bab767d6f05.htm) | 返回集合 key 的基数(集合中元素的数量)。当 key 不存在时，返回 0 。 |
| 公共方法 | [SetCardAsync](0d0665a9-a174-d154-f95b-2bf10d885bc3.htm) | 返回集合 key 的基数(集合中元素的数量)。当 key 不存在时，返回 0 。 |
| 公共方法 | [SetDiff(String, String)](e9ecf346-db86-b5f4-4f91-dcf5cdbefc1e.htm) | 返回一个集合的全部成员，该集合是所有给定集合之间的差集。 |
| 公共方法 | [SetDiff(String, String)](93f80f0e-351d-b56a-96ec-7a2af40349fe.htm) | 返回一个集合的全部成员，该集合是所有给定集合之间的差集。 |
| 公共方法 | [SetDiffAsync(String, String)](8f89285e-ce7c-802a-a473-9347afccba05.htm) | 返回一个集合的全部成员，该集合是所有给定集合之间的差集。 |
| 公共方法 | [SetDiffAsync(String, String)](2c1c06c0-6089-b734-a5f2-b8c7a27ff88d.htm) | 返回一个集合的全部成员，该集合是所有给定集合之间的差集。 |
| 公共方法 | [SetDiffStore(String, String, String)](5dcd4b1b-48b7-52a8-e0d5-7a4969b9aef5.htm) | 这个命令的作用和 SDIFF 类似，但它将结果保存到 destination 集合，而不是简单地返回结果集。如果 destination 集合已经存在，则将其覆盖。destination 可以是 key 本身。 |
| 公共方法 | [SetDiffStore(String, String, String)](7f4dcdf7-214e-a9c5-b93f-6fc71c2d8016.htm) | 这个命令的作用和 SDIFF 类似，但它将结果保存到 destination 集合，而不是简单地返回结果集。如果 destination 集合已经存在，则将其覆盖。destination 可以是 key 本身。 |
| 公共方法 | [SetDiffStoreAsync(String, String, String)](948e092b-b776-43f6-bf01-a432d1fe3f68.htm) | 这个命令的作用和 SDIFF 类似，但它将结果保存到 destination 集合，而不是简单地返回结果集。如果 destination 集合已经存在，则将其覆盖。destination 可以是 key 本身。 |
| 公共方法 | [SetDiffStoreAsync(String, String, String)](42a69dfd-f064-4337-8713-d89339c62906.htm) | 这个命令的作用和 SDIFF 类似，但它将结果保存到 destination 集合，而不是简单地返回结果集。如果 destination 集合已经存在，则将其覆盖。destination 可以是 key 本身。 |
| 公共方法 | [SetInter(String, String)](f627bd19-f348-9704-0581-a8d442ee60b4.htm) | 返回一个集合的全部成员，该集合是所有给定集合的交集。不存在的 key 被视为空集。当给定集合当中有一个空集时，结果也为空集(根据集合运算定律)。 |
| 公共方法 | [SetInter(String, String)](8b28cdae-05cc-3259-840f-c91cafa6d8bc.htm) | 返回一个集合的全部成员，该集合是所有给定集合的交集。不存在的 key 被视为空集。当给定集合当中有一个空集时，结果也为空集(根据集合运算定律)。 |
| 公共方法 | [SetInterAsync(String, String)](4035a5be-411b-638b-6559-b56e4442a2c4.htm) | 返回一个集合的全部成员，该集合是所有给定集合的交集。不存在的 key 被视为空集。当给定集合当中有一个空集时，结果也为空集(根据集合运算定律)。 |
| 公共方法 | [SetInterAsync(String, String)](4f9ab180-ff1f-bdf7-2c21-f14f093ac5dc.htm) | 返回一个集合的全部成员，该集合是所有给定集合的交集。不存在的 key 被视为空集。当给定集合当中有一个空集时，结果也为空集(根据集合运算定律)。 |
| 公共方法 | [SetInterStore(String, String, String)](1b4b8824-96fa-2ac7-b54b-ee3b2ece1c69.htm) | 这个命令类似于 SINTER 命令，但它将结果保存到 destination 集合，而不是简单地返回结果集。如果 destination 集合已经存在，则将其覆盖。destination 可以是 key 本身。 |
| 公共方法 | [SetInterStore(String, String, String)](1b14a4a9-4149-5e72-3d06-f4eab94d8d1f.htm) | 这个命令类似于 SINTER 命令，但它将结果保存到 destination 集合，而不是简单地返回结果集。如果 destination 集合已经存在，则将其覆盖。destination 可以是 key 本身。 |
| 公共方法 | [SetInterStoreAsync(String, String, String)](8d29dcd5-7a8a-39ed-5c0f-eb0b8ffa135f.htm) | 这个命令类似于 SINTER 命令，但它将结果保存到 destination 集合，而不是简单地返回结果集。如果 destination 集合已经存在，则将其覆盖。destination 可以是 key 本身。 |
| 公共方法 | [SetInterStoreAsync(String, String, String)](01f68462-704d-7d6c-c7f3-fa4c0edceb36.htm) | 这个命令类似于 SINTER 命令，但它将结果保存到 destination 集合，而不是简单地返回结果集。如果 destination 集合已经存在，则将其覆盖。destination 可以是 key 本身。 |
| 公共方法 | [SetIsMember](1c6d465b-353a-34f2-8f46-9d4037e58911.htm) | 判断 member 元素是否集合 key 的成员。如果 member 元素是集合的成员，返回 1 。如果 member 元素不是集合的成员，或 key 不存在，返回 0 。 |
| 公共方法 | [SetIsMemberAsync](21419054-3273-2aab-d486-633ab8861039.htm) | 判断 member 元素是否集合 key 的成员。如果 member 元素是集合的成员，返回 1 。如果 member 元素不是集合的成员，或 key 不存在，返回 0 。 |
| 公共方法 | [SetLoginAccount](725afc25-6e3f-bd7d-25ed-580270c8fc56.htm) | 设置当前的登录的账户名和密码信息，并启用账户验证的功能，账户名为空时设置不生效  Set the current login account name and password information, and enable the account verification function. The account name setting will not take effect when it is empty (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | [SetMembers](7a5b277f-1fe1-141f-346b-3722d97ac7ca.htm) | 返回集合 key 中的所有成员。不存在的 key 被视为空集合。 |
| 公共方法 | [SetMembersAsync](ce86155f-4be2-3836-3594-1398d78183ad.htm) | 返回集合 key 中的所有成员。不存在的 key 被视为空集合。 |
| 公共方法 | [SetMove](8b76c71a-c755-708b-e1e7-61eda8b9456c.htm) | 将 member 元素从 source 集合移动到 destination 集合。如果 source 集合不存在或不包含指定的 member 元素，则 SMOVE 命令不执行任何操作，仅返回 0 。 否则， member 元素从 source 集合中被移除，并添加到 destination 集合中去。当 destination 集合已经包含 member 元素时， SMOVE 命令只是简单地将 source 集合中的 member 元素删除。 当 source 或 destination 不是集合类型时，返回一个错误。 |
| 公共方法 | [SetMoveAsync](dc01893c-b09e-e134-49d3-76d0dbc308b9.htm) | 将 member 元素从 source 集合移动到 destination 集合。如果 source 集合不存在或不包含指定的 member 元素，则 SMOVE 命令不执行任何操作，仅返回 0 。 否则， member 元素从 source 集合中被移除，并添加到 destination 集合中去。当 destination 集合已经包含 member 元素时， SMOVE 命令只是简单地将 source 集合中的 member 元素删除。 当 source 或 destination 不是集合类型时，返回一个错误。 |
| 公共方法代码示例 | [SetPersistentConnection](d8e8df9f-cbe9-6863-54e9-aa3b38906b26.htm) | 在读取数据之前可以调用本方法将客户端设置为长连接模式，相当于跳过了ConnectServer的结果验证，对异形客户端无效，当第一次进行通信时再进行创建连接请求。  Before reading the data, you can call this method to set the client to the long connection mode, which is equivalent to skipping the result verification of ConnectServer, and it is invalid for the alien client. When the first communication is performed, the connection creation request is performed. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | [SetPipeSocket](b6a7a243-2a6d-c4ea-45b7-bf980450a2bd.htm) | 设置一个新的网络管道，一般来说不需要调用本方法，当多个网口设备共用一个网络连接时才需要使用本方法进行设置共享的管道。  To set up a new network channel, generally speaking, you do not need to call this method. This method is only needed to set up a shared channel when multiple network port devices share a network connection. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | [SetPop](8375117d-d7b5-1f33-cb59-80cc632b736d.htm) | 移除并返回集合中的一个随机元素。如果只想获取一个随机元素，但不想该元素从集合中被移除的话，可以使用 SRANDMEMBER 命令。 |
| 公共方法 | [SetPopAsync](66039533-6b1e-ad66-e89b-ef8f7d46c47b.htm) | 移除并返回集合中的一个随机元素。如果只想获取一个随机元素，但不想该元素从集合中被移除的话，可以使用 SRANDMEMBER 命令。 |
| 公共方法 | [SetRandomMember(String)](a812ee61-36b4-1e8e-7545-1c304eefcf3b.htm) | 那么返回集合中的一个随机元素。 |
| 公共方法 | [SetRandomMember(String, Int32)](fdd8ab0f-92da-8e38-2367-3cc6955d4770.htm) | 返回集合中的多个随机元素。  如果 count 为正数，且小于集合基数，那么命令返回一个包含 count 个元素的数组，数组中的元素各不相同。如果 count 大于等于集合基数，那么返回整个集合。  如果 count 为负数，那么命令返回一个数组，数组中的元素可能会重复出现多次，而数组的长度为 count 的绝对值。 |
| 公共方法 | [SetRandomMemberAsync(String)](ddc2b940-2954-4e8b-9a50-4f00dd729dca.htm) | 那么返回集合中的一个随机元素。 |
| 公共方法 | [SetRandomMemberAsync(String, Int32)](363a263f-5e96-9e89-f0d2-1a59458b56c9.htm) | 返回集合中的多个随机元素。  如果 count 为正数，且小于集合基数，那么命令返回一个包含 count 个元素的数组，数组中的元素各不相同。如果 count 大于等于集合基数，那么返回整个集合。  如果 count 为负数，那么命令返回一个数组，数组中的元素可能会重复出现多次，而数组的长度为 count 的绝对值。 |
| 公共方法 | [SetRemove(String, String)](915dac15-42ce-1db0-1975-6e2619ede6b1.htm) | 移除集合 key 中的一个元素，不存在的 member 元素会被忽略。 |
| 公共方法 | [SetRemove(String, String)](6f77ad2b-b35f-c91c-2cb6-b3e39cba9f3a.htm) | 移除集合 key 中的一个或多个 member 元素，不存在的 member 元素会被忽略。 |
| 公共方法 | [SetRemoveAsync(String, String)](9e56cae7-1944-e5fa-90c1-508e1c87499e.htm) | 移除集合 key 中的一个元素，不存在的 member 元素会被忽略。 |
| 公共方法 | [SetRemoveAsync(String, String)](3023a9ef-1c62-e0f9-c913-d3b15a77821e.htm) | 移除集合 key 中的一个或多个 member 元素，不存在的 member 元素会被忽略。 |
| 公共方法 | [SetUnion(String, String)](86a1df0d-4454-f293-fb36-9ddb123a92e8.htm) | 返回一个集合的全部成员，该集合是所有给定集合的并集。不存在的 key 被视为空集。 |
| 公共方法 | [SetUnion(String, String)](2d459ccd-11fe-cc5b-d3c5-526e2cd8268e.htm) | 返回一个或多个集合的全部成员，该集合是所有给定集合的并集。不存在的 key 被视为空集。 |
| 公共方法 | [SetUnionAsync(String, String)](1a1439cb-d4a0-9c89-b4e5-534f4e16eee6.htm) | 返回一个集合的全部成员，该集合是所有给定集合的并集。不存在的 key 被视为空集。 |
| 公共方法 | [SetUnionAsync(String, String)](3c66c82e-0b60-c57d-5f74-ebd1c408b83a.htm) | 返回一个或多个集合的全部成员，该集合是所有给定集合的并集。不存在的 key 被视为空集。 |
| 公共方法 | [SetUnionStore(String, String, String)](d3c6d04a-2dfb-f129-5104-56b9602661cf.htm) | 这个命令类似于 SUNION 命令，但它将结果保存到 destination 集合，而不是简单地返回结果集。如果 destination 已经存在，则将其覆盖。destination 可以是 key 本身。 |
| 公共方法 | [SetUnionStore(String, String, String)](683dea58-da66-fda5-20b5-35a393be979b.htm) | 这个命令类似于 SUNION 命令，但它将结果保存到 destination 集合，而不是简单地返回结果集。如果 destination 已经存在，则将其覆盖。destination 可以是 key 本身。 |
| 公共方法 | [SetUnionStoreAsync(String, String, String)](0b81d3ed-fd9c-4374-b670-f6a7d33ea7fc.htm) | 这个命令类似于 SUNION 命令，但它将结果保存到 destination 集合，而不是简单地返回结果集。如果 destination 已经存在，则将其覆盖。destination 可以是 key 本身。 |
| 公共方法 | [SetUnionStoreAsync(String, String, String)](121a97ed-2d7c-ef0b-4660-4efbdcb060f1.htm) | 这个命令类似于 SUNION 命令，但它将结果保存到 destination 集合，而不是简单地返回结果集。如果 destination 已经存在，则将其覆盖。destination 可以是 key 本身。 |
| 公共方法 | [SubscribeMessage(String)](5a08d3e8-17f8-273a-77ea-fa4f8c682f29.htm) | 从Redis服务器订阅一个或多个主题信息  Subscribe to one or more topics from the redis server |
| 公共方法 | [SubscribeMessage(String)](63b458d5-24ae-2600-ee11-1e0eff778662.htm) | 从Redis服务器订阅一个或多个主题信息  Subscribe to one or more topics from the redis server |
| 公共方法 | [ToString](6ff40e11-12f4-8759-45a0-594ca3e4fcdd.htm) | (重写 [NetworkDoubleBaseToString](85a8532f-8fc5-886c-9d76-70844e0a7d64.htm).) |
| 公共方法 | [UnpackResponseContent](32bad9a0-b59d-c171-ce85-74055edfc597.htm) | 根据对方返回的报文命令，对命令进行基本的拆包，例如各种Modbus协议拆包为统一的核心报文，还支持对报文的验证  According to the message command returned by the other party, the command is basically unpacked, for example, various Modbus protocols are unpacked into a unified core message, and the verification of the message is also supported (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | [UnSubscribeMessage(String)](dfdd1fbd-c51a-8456-ff3e-913a51b69b6d.htm) | 取消订阅一个或多个主题信息，取消之后，当前的订阅数据就不在接收到。  Unsubscribe from multiple topic information. After cancellation, the current subscription data will not be received. |
| 公共方法 | [UnSubscribeMessage(String)](e24cb10d-a565-8e43-8d83-4c6e083884d2.htm) | 取消订阅一个或多个主题信息，取消之后，当前的订阅数据就不在接收到。  Unsubscribe from multiple topic information. After cancellation, the current subscription data will not be received. |
| 公共方法代码示例 | [WriteT](70a08f37-0f86-76f1-55a5-df1ddd5a7d16.htm) | 从设备里写入支持Hsl特性的数据内容， 该特性为[HslRedisKeyAttribute](0b3a62cc-e2c8-dfb7-89de-7da47983c3c0.htm) ，[HslRedisHashFieldAttribute](13118df4-a00f-12bb-887d-e78f8b46236b.htm) 需要注意的是写入并不支持[HslRedisListAttribute](fd07c099-2c4b-18b4-45dc-3f869fdc50c3.htm)，[HslRedisListItemAttribute](2bbc9c13-e0d2-6942-b529-1ef6ad142370.htm)特性，详细参考代码示例的操作说明。 |
| 公共方法 | [WriteAndPublishKey](69029e5e-a78a-5c83-0300-b6dd057f9308.htm) | 将字符串值 value 关联到 key 。并发布一个订阅的频道数据，都成功时，才返回成功 |
| 公共方法 | [WriteAndPublishKeyAsync](d44abf4b-02bf-cfa5-cbc8-0b36edab4521.htm) | 将字符串值 value 关联到 key 。并发布一个订阅的频道数据，都成功时，才返回成功 |
| 公共方法代码示例 | [WriteAsyncT](9d54e352-2bd0-2655-3e14-49a0213e6040.htm) | 从设备里写入支持Hsl特性的数据内容， 该特性为[HslRedisKeyAttribute](0b3a62cc-e2c8-dfb7-89de-7da47983c3c0.htm) ，[HslRedisHashFieldAttribute](13118df4-a00f-12bb-887d-e78f8b46236b.htm) 需要注意的是写入并不支持[HslRedisListAttribute](fd07c099-2c4b-18b4-45dc-3f869fdc50c3.htm)，[HslRedisListItemAttribute](2bbc9c13-e0d2-6942-b529-1ef6ad142370.htm)特性，详细参考代码示例的操作说明。 |
| 公共方法 | [WriteExpireKey](816ca8c0-5e3a-4609-acc3-6f227ac2dfde.htm) | 将值 value 关联到 key ，并将 key 的生存时间设为 seconds (以秒为单位)。如果 key 已经存在， SETEX 命令将覆写旧值。 |
| 公共方法 | [WriteExpireKeyAsync](f1dce55a-5681-32f6-2c7c-cc664465490c.htm) | 将值 value 关联到 key ，并将 key 的生存时间设为 seconds (以秒为单位)。如果 key 已经存在， SETEX 命令将覆写旧值。 |
| 公共方法 | [WriteHashKey(String, String, String)](09ed7629-7800-0eff-6ad2-b9f6cbe9100e.htm) | 将哈希表 key 中的域 field 的值设为 value 。 如果 key 不存在，一个新的哈希表被创建并进行 HSET 操作。 如果域 field 已经存在于哈希表中，旧值将被覆盖。 如果 field 是哈希表中的一个新建域，并且值设置成功，返回 1 。 如果哈希表中域 field 已经存在且旧值已被新值覆盖，返回 0 。 |
| 公共方法 | [WriteHashKey(String, String, String)](85db4be6-4ef2-6062-2a2c-607a2c8255b6.htm) | 同时将多个 field-value (域-值)对设置到哈希表 key 中。 此命令会覆盖哈希表中已存在的域。 如果 key 不存在，一个空哈希表被创建并执行 HMSET 操作。 |
| 公共方法 | [WriteHashKeyAsync(String, String, String)](5801856f-9b52-fded-8714-80c52b90d2a8.htm) | 将哈希表 key 中的域 field 的值设为 value 。 如果 key 不存在，一个新的哈希表被创建并进行 HSET 操作。 如果域 field 已经存在于哈希表中，旧值将被覆盖。 如果 field 是哈希表中的一个新建域，并且值设置成功，返回 1 。 如果哈希表中域 field 已经存在且旧值已被新值覆盖，返回 0 。 |
| 公共方法 | [WriteHashKeyAsync(String, String, String)](f638b9b1-6813-8e8d-6b75-1c6784cec658.htm) | 同时将多个 field-value (域-值)对设置到哈希表 key 中。 此命令会覆盖哈希表中已存在的域。 如果 key 不存在，一个空哈希表被创建并执行 HMSET 操作。 |
| 公共方法 | [WriteHashKeyNx](819ddbb0-74d4-4667-6e55-372fac166c68.htm) | 将哈希表 key 中的域 field 的值设置为 value ，当且仅当域 field 不存在。若域 field 已经存在，该操作无效。 设置成功，返回 1 。如果给定域已经存在且没有操作被执行，返回 0 。 |
| 公共方法 | [WriteHashKeyNxAsync](b4571eb9-be29-fbb8-691b-f3a93d9245ec.htm) | 将哈希表 key 中的域 field 的值设置为 value ，当且仅当域 field 不存在。若域 field 已经存在，该操作无效。 设置成功，返回 1 。如果给定域已经存在且没有操作被执行，返回 0 。 |
| 公共方法 | [WriteKey(String, String)](eff07a08-c650-1d95-1d9d-d5e3736c439d.htm) | 将字符串值 value 关联到 key 。 如果 key 已经持有其他值， SET 就覆写旧值，无视类型。 对于某个原本带有生存时间（TTL）的键来说， 当 SET 命令成功在这个键上执行时，这个键原有的 TTL 将被清除。 |
| 公共方法 | [WriteKey(String, String)](97158b0a-0ca6-749b-ee3e-80557e13c619.htm) | 同时设置一个或多个 key-value 对。 如果某个给定 key 已经存在，那么 MSET 会用新值覆盖原来的旧值，如果这不是你所希望的效果，请考虑使用 MSETNX 命令：它只会在所有给定 key 都不存在的情况下进行设置操作。 |
| 公共方法 | [WriteKeyAsync(String, String)](50c841fd-8905-fe84-46ff-cfd3a0b8e4e5.htm) | 将字符串值 value 关联到 key 。 如果 key 已经持有其他值， SET 就覆写旧值，无视类型。 对于某个原本带有生存时间（TTL）的键来说， 当 SET 命令成功在这个键上执行时，这个键原有的 TTL 将被清除。 |
| 公共方法 | [WriteKeyAsync(String, String)](7f526544-53c5-6d54-5300-1c100dcb5c8e.htm) | 同时设置一个或多个 key-value 对。 如果某个给定 key 已经存在，那么 MSET 会用新值覆盖原来的旧值，如果这不是你所希望的效果，请考虑使用 MSETNX 命令：它只会在所有给定 key 都不存在的情况下进行设置操作。 |
| 公共方法 | [WriteKeyIfNotExists](9d48312b-2ce1-077f-bf47-3e6b5a328c45.htm) | 将 key 的值设为 value ，当且仅当 key 不存在。若给定的 key 已经存在，则 SETNX 不做任何动作。设置成功，返回 1 。设置失败，返回 0 。 |
| 公共方法 | [WriteKeyIfNotExistsAsync](e6844e26-ba17-a0a8-0ff2-bb1a1e8bce8e.htm) | 将 key 的值设为 value ，当且仅当 key 不存在。若给定的 key 已经存在，则 SETNX 不做任何动作。设置成功，返回 1 。设置失败，返回 0 。 |
| 公共方法 | [WriteKeyRange](c1f98cc9-f791-4a35-bf67-ea29e8930d32.htm) | 用 value 参数覆写(overwrite)给定 key 所储存的字符串值，从偏移量 offset 开始。不存在的 key 当作空白字符串处理。返回被 SETRANGE 修改之后，字符串的长度。 |
| 公共方法 | [WriteKeyRangeAsync](cdfaffec-518b-2754-5a10-3c5b68ee1767.htm) | 用 value 参数覆写(overwrite)给定 key 所储存的字符串值，从偏移量 offset 开始。不存在的 key 当作空白字符串处理。返回被 SETRANGE 修改之后，字符串的长度。 |
| 受保护的方法 | [WriteStreamFromSocket](97bf95ea-51fe-8f2e-5da6-1c83766aaa78.htm) | 从套接字中接收所有的数据然后写入到指定的流当中去，需要指定数据的长度，支持按照百分比进行进度报告  Receives all data from the socket and writes it to the specified stream. The length of the data needs to be specified, and progress reporting is supported in percentage. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [WriteStreamFromSocketAsync](b4cf7f6d-1e0c-a266-4c49-683fdb14a922.htm) | 从套接字中接收所有的数据然后写入到指定的流当中去，需要指定数据的长度，支持按照百分比进行进度报告  Receives all data from the socket and writes it to the specified stream. The length of the data needs to be specified, and progress reporting is supported in percentage. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 公共方法 | [ZSetAdd(String, String, Double)](593c2f94-650b-1b61-d5ec-09b5231295c0.htm) | 将一个 member 元素及其 score 值加入到有序集 key 当中。如果某个 member 已经是有序集的成员，那么更新这个 member 的 score 值，并通过重新插入这个 member 元素，来保证该 member 在正确的位置上。 score 值可以是整数值或双精度浮点数。  如果 key 不存在，则创建一个空的有序集并执行 ZADD 操作。当 key 存在但不是有序集类型时，返回一个错误。 |
| 公共方法 | [ZSetAdd(String, String, Double)](edbcdc0f-4537-5be6-9ddc-7c16d9773ee6.htm) | 将一个或多个 member 元素及其 score 值加入到有序集 key 当中。如果某个 member 已经是有序集的成员，那么更新这个 member 的 score 值，并通过重新插入这个 member 元素，来保证该 member 在正确的位置上。 score 值可以是整数值或双精度浮点数。  如果 key 不存在，则创建一个空的有序集并执行 ZADD 操作。当 key 存在但不是有序集类型时，返回一个错误。 |
| 公共方法 | [ZSetAddAsync(String, String, Double)](2e2ef743-be44-5c57-e604-a8c96b584693.htm) | 将一个 member 元素及其 score 值加入到有序集 key 当中。如果某个 member 已经是有序集的成员，那么更新这个 member 的 score 值，并通过重新插入这个 member 元素，来保证该 member 在正确的位置上。 score 值可以是整数值或双精度浮点数。  如果 key 不存在，则创建一个空的有序集并执行 ZADD 操作。当 key 存在但不是有序集类型时，返回一个错误。 |
| 公共方法 | [ZSetAddAsync(String, String, Double)](79d1d709-cc24-d2be-5e4d-51cb3395290a.htm) | 将一个或多个 member 元素及其 score 值加入到有序集 key 当中。如果某个 member 已经是有序集的成员，那么更新这个 member 的 score 值，并通过重新插入这个 member 元素，来保证该 member 在正确的位置上。 score 值可以是整数值或双精度浮点数。  如果 key 不存在，则创建一个空的有序集并执行 ZADD 操作。当 key 存在但不是有序集类型时，返回一个错误。 |
| 公共方法 | [ZSetCard](aa928601-c5f6-63fb-6bd8-e2e383796043.htm) | 返回有序集 key 的基数。 |
| 公共方法 | [ZSetCardAsync](b760ff15-4c7e-510f-1927-e17e6310e142.htm) | 返回有序集 key 的基数。 |
| 公共方法 | [ZSetCount](9872ddc0-8b8d-d9d6-4d08-64001fdefd29.htm) | 返回有序集 key 中， score 值在 min 和 max 之间(默认包括 score 值等于 min 或 max )的成员的数量。 |
| 公共方法 | [ZSetCountAsync](b294e1ee-27c7-7326-eb4b-8669040aecf0.htm) | 返回有序集 key 中， score 值在 min 和 max 之间(默认包括 score 值等于 min 或 max )的成员的数量。 |
| 公共方法 | [ZSetIncreaseBy](c84e0c14-15cd-451a-dec4-2a112a73772e.htm) | 为有序集 key 的成员 member 的 score 值加上增量 increment 。可以通过传递一个负数值 increment ，让 score 减去相应的值，比如 ZINCRBY key -5 member ，就是让 member 的 score 值减去 5 。 当 key 不存在，或 member 不是 key 的成员时， ZINCRBY key increment member 等同于 ZADD key increment member 。当 key 不是有序集类型时，返回一个错误。 |
| 公共方法 | [ZSetIncreaseByAsync](9805ac6b-aa05-2152-7dbb-7f8077ea462b.htm) | 为有序集 key 的成员 member 的 score 值加上增量 increment 。可以通过传递一个负数值 increment ，让 score 减去相应的值，比如 ZINCRBY key -5 member ，就是让 member 的 score 值减去 5 。 当 key 不存在，或 member 不是 key 的成员时， ZINCRBY key increment member 等同于 ZADD key increment member 。当 key 不是有序集类型时，返回一个错误。 |
| 公共方法 | [ZSetRange](759e2523-e2ad-0dc1-a904-354031a8ae50.htm) | 返回有序集 key 中，指定区间内的成员。其中成员的位置按 score 值递增(从小到大)来排序。具有相同 score 值的成员按字典序来排列。 下标参数 start 和 stop 都以 0 为底，也就是说，以 0 表示有序集第一个成员，以 1 表示有序集第二个成员，以此类推。你也可以使用负数下标，以 -1 表示最后一个成员， -2 表示倒数第二个成员，以此类推。 |
| 公共方法 | [ZSetRangeAsync](7724eb62-2653-65a7-779d-391d630c2faf.htm) | 返回有序集 key 中，指定区间内的成员。其中成员的位置按 score 值递增(从小到大)来排序。具有相同 score 值的成员按字典序来排列。 下标参数 start 和 stop 都以 0 为底，也就是说，以 0 表示有序集第一个成员，以 1 表示有序集第二个成员，以此类推。你也可以使用负数下标，以 -1 表示最后一个成员， -2 表示倒数第二个成员，以此类推。 |
| 公共方法 | [ZSetRangeByScore](acab27d0-6b6c-1f25-00ad-fafd4ac068fa.htm) | 返回有序集 key 中，所有 score 值介于 min 和 max 之间(包括等于 min 或 max )的成员。有序集成员按 score 值递增(从小到大)次序排列。 min 和 max 可以是 -inf 和 +inf ，这样一来，你就可以在不知道有序集的最低和最高 score 值的情况下，使用 ZRANGEBYSCORE 这类命令。 默认情况下，区间的取值使用闭区间 (小于等于或大于等于)，你也可以通过给参数前增加 "(" 符号来使用可选的开区间 (小于或大于)。"(5"代表不包含5 |
| 公共方法 | [ZSetRangeByScoreAsync](d1e0aa3c-d5f9-caeb-4017-699f956e661d.htm) | 返回有序集 key 中，所有 score 值介于 min 和 max 之间(包括等于 min 或 max )的成员。有序集成员按 score 值递增(从小到大)次序排列。 min 和 max 可以是 -inf 和 +inf ，这样一来，你就可以在不知道有序集的最低和最高 score 值的情况下，使用 ZRANGEBYSCORE 这类命令。 默认情况下，区间的取值使用闭区间 (小于等于或大于等于)，你也可以通过给参数前增加 "(" 符号来使用可选的开区间 (小于或大于)。"(5"代表不包含5 |
| 公共方法 | [ZSetRank](7957411f-32c1-4e07-a4f2-9d1948c63a2b.htm) | 返回有序集 key 中成员 member 的排名。其中有序集成员按 score 值递增(从小到大)顺序排列。排名以 0 为底，也就是说， score 值最小的成员排名为 0 。 |
| 公共方法 | [ZSetRankAsync](1181fd8c-02aa-9d53-cbeb-35d542bed241.htm) | 返回有序集 key 中成员 member 的排名。其中有序集成员按 score 值递增(从小到大)顺序排列。排名以 0 为底，也就是说， score 值最小的成员排名为 0 。 |
| 公共方法 | [ZSetRemove(String, String)](bde5c871-2a34-579a-a9f8-2e685047cf7b.htm) | 移除有序集 key 中的指定成员，不存在的成员将被忽略。当 key 存在但不是有序集类型时，返回一个错误。 |
| 公共方法 | [ZSetRemove(String, String)](66350924-7cee-4c1f-a584-4e480cd24dab.htm) | 移除有序集 key 中的一个或多个成员，不存在的成员将被忽略。当 key 存在但不是有序集类型时，返回一个错误。 |
| 公共方法 | [ZSetRemoveAsync(String, String)](4e9ef9ce-44d3-82be-bfcc-a2cf7743e6e0.htm) | 移除有序集 key 中的指定成员，不存在的成员将被忽略。当 key 存在但不是有序集类型时，返回一个错误。 |
| 公共方法 | [ZSetRemoveAsync(String, String)](f2528ee0-4f06-dd02-ab3f-b1232e123cd0.htm) | 移除有序集 key 中的一个或多个成员，不存在的成员将被忽略。当 key 存在但不是有序集类型时，返回一个错误。 |
| 公共方法 | [ZSetRemoveRangeByRank](cb81dfc2-6e30-0aca-8bcb-348cdc32bc35.htm) | 移除有序集 key 中，指定排名(rank)区间内的所有成员。区间分别以下标参数 start 和 stop 指出，包含 start 和 stop 在内。 下标参数 start 和 stop 都以 0 为底，也就是说，以 0 表示有序集第一个成员，以 1 表示有序集第二个成员，以此类推。你也可以使用负数下标，以 -1 表示最后一个成员， -2 表示倒数第二个成员，以此类推。 |
| 公共方法 | [ZSetRemoveRangeByRankAsync](9bde4b5f-b9bd-e01b-feb4-82c9b532dee4.htm) | 移除有序集 key 中，指定排名(rank)区间内的所有成员。区间分别以下标参数 start 和 stop 指出，包含 start 和 stop 在内。 下标参数 start 和 stop 都以 0 为底，也就是说，以 0 表示有序集第一个成员，以 1 表示有序集第二个成员，以此类推。你也可以使用负数下标，以 -1 表示最后一个成员， -2 表示倒数第二个成员，以此类推。 |
| 公共方法 | [ZSetRemoveRangeByScore](ac7be899-fe4e-3f12-d60f-81a038a927c7.htm) | 移除有序集 key 中，所有 score 值介于 min 和 max 之间(包括等于 min 或 max )的成员。 min 和 max 可以是 -inf 和 +inf ，这样一来，你就可以在不知道有序集的最低和最高 score 值的情况下，使用 ZRANGEBYSCORE 这类命令。 默认情况下，区间的取值使用闭区间 (小于等于或大于等于)，你也可以通过给参数前增加 "(" 符号来使用可选的开区间 (小于或大于)。例如"(5"代表不包括5 |
| 公共方法 | [ZSetRemoveRangeByScoreAsync](6ccee6bf-8635-14bf-2e7f-d1ac351c70c9.htm) | 移除有序集 key 中，所有 score 值介于 min 和 max 之间(包括等于 min 或 max )的成员。 min 和 max 可以是 -inf 和 +inf ，这样一来，你就可以在不知道有序集的最低和最高 score 值的情况下，使用 ZRANGEBYSCORE 这类命令。 默认情况下，区间的取值使用闭区间 (小于等于或大于等于)，你也可以通过给参数前增加 "(" 符号来使用可选的开区间 (小于或大于)。例如"(5"代表不包括5 |
| 公共方法 | [ZSetReverseRange](ef6defc7-62a2-c8a4-4025-32f6b3a82ae4.htm) | 返回有序集 key 中，指定区间内的成员。其中成员的位置按 score 值递减(从大到小)来排列。具有相同 score 值的成员按字典序来排列。 下标参数 start 和 stop 都以 0 为底，也就是说，以 0 表示有序集第一个成员，以 1 表示有序集第二个成员，以此类推。你也可以使用负数下标，以 -1 表示最后一个成员， -2 表示倒数第二个成员，以此类推。 |
| 公共方法 | [ZSetReverseRangeAsync](b40d37b7-1118-6aec-c7da-85ec54b457d0.htm) | 返回有序集 key 中，指定区间内的成员。其中成员的位置按 score 值递减(从大到小)来排列。具有相同 score 值的成员按字典序来排列。 下标参数 start 和 stop 都以 0 为底，也就是说，以 0 表示有序集第一个成员，以 1 表示有序集第二个成员，以此类推。你也可以使用负数下标，以 -1 表示最后一个成员， -2 表示倒数第二个成员，以此类推。 |
| 公共方法 | [ZSetReverseRangeByScore](3c1a26a8-f8e5-a835-7051-94b4d9c851ab.htm) | 返回有序集 key 中，所有 score 值介于 min 和 max 之间(包括等于 min 或 max )的成员。序集成员按 score 值递减(从大到小)的次序排列。 min 和 max 可以是 -inf 和 +inf ，这样一来，你就可以在不知道有序集的最低和最高 score 值的情况下，使用 ZRANGEBYSCORE 这类命令。 默认情况下，区间的取值使用闭区间 (小于等于或大于等于)，你也可以通过给参数前增加 ( 符号来使用可选的开区间 (小于或大于)。(5代表不包含5 |
| 公共方法 | [ZSetReverseRangeByScoreAsync](3a3cfbff-cee7-e1f3-d9c1-e7528db4e998.htm) | 返回有序集 key 中，所有 score 值介于 min 和 max 之间(包括等于 min 或 max )的成员。序集成员按 score 值递减(从大到小)的次序排列。 min 和 max 可以是 -inf 和 +inf ，这样一来，你就可以在不知道有序集的最低和最高 score 值的情况下，使用 ZRANGEBYSCORE 这类命令。 默认情况下，区间的取值使用闭区间 (小于等于或大于等于)，你也可以通过给参数前增加 ( 符号来使用可选的开区间 (小于或大于)。(5代表不包含5 |
| 公共方法 | [ZSetReverseRank](d524195b-6b02-a8e8-330a-3454a6ed21cd.htm) | 返回有序集 key 中成员 member 的排名。其中有序集成员按 score 值递减(从大到小)排序。排名以 0 为底，也就是说，score 值最大的成员排名为 0 。 |
| 公共方法 | [ZSetReverseRankAsync](c010781c-2851-7158-f953-1779929e3242.htm) | 返回有序集 key 中成员 member 的排名。其中有序集成员按 score 值递减(从大到小)排序。排名以 0 为底，也就是说，score 值最大的成员排名为 0 。 |
| 公共方法 | [ZSetScore](a98ee7e4-f361-f403-975b-a197e0ca0268.htm) | 返回有序集 key 中，成员 member 的 score 值。如果 member 元素不是有序集 key 的成员，或 key 不存在，返回 nil 。 |
| 公共方法 | [ZSetScoreAsync](fd24c433-944e-2ecb-4bcd-266bfc6825d4.htm) | 返回有序集 key 中，成员 member 的 score 值。如果 member 元素不是有序集 key 的成员，或 key 不存在，返回 nil 。 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[RedisClient 类](e33d9088-6df0-32e9-33f1-426e944c3b7f.htm)

[HslCommunication.Enthernet.Redis 命名空间](cc44ce00-3781-c690-b233-24000870bf42.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## AppendKey 方法 

[原文連結](http://api.hslcommunication.cn/html/f2521b8c-0fb1-6b34-4adc-ca6804c7195f.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Enthernet.Redis](../html/cc44ce00-3781-c690-b233-24000870bf42.htm "HslCommunication.Enthernet.Redis")

[RedisClient 类](../html/e33d9088-6df0-32e9-33f1-426e944c3b7f.htm "RedisClient 类")

[RedisClient 方法](../html/e41be4a5-1490-ffbd-f241-938f2d6c3f07.htm "RedisClient 方法")

[AppendKey 方法](../html/f2521b8c-0fb1-6b34-4adc-ca6804c7195f.htm "AppendKey 方法 ")

[AppendKeyAsync 方法](../html/c39b4c6a-b7f6-68fe-dfc9-7cce1b653e24.htm "AppendKeyAsync 方法 ")

[ChangePassword 方法](../html/78a37d4d-bdee-7a09-d0ce-dc245627b1d1.htm "ChangePassword 方法 ")

[ChangePasswordAsync 方法](../html/58f6f736-093d-6e4d-3ad3-19590fa2a3a3.htm "ChangePasswordAsync 方法 ")

[DBSize 方法](../html/5814129a-9a73-25ac-e726-fe1da141eb28.htm "DBSize 方法 ")

[DBSizeAsync 方法](../html/243a1f5f-a474-2e20-1298-06eff125f8b5.htm "DBSizeAsync 方法 ")

[DecrementKey 方法](../html/e9fa6c94-43ca-8291-79ca-6e04ca9289e7.htm "DecrementKey 方法 ")

[DecrementKeyAsync 方法](../html/6837e351-b8ef-a90e-c5e2-73b31621f837.htm "DecrementKeyAsync 方法 ")

[DeleteHashKey 方法](../html/19e3ebc3-e6fc-71cf-e80d-cb73e3011170.htm "DeleteHashKey 方法 ")

[DeleteHashKeyAsync 方法](../html/0dcc0d8a-1ed1-5c81-fb0e-e1c93be9d579.htm "DeleteHashKeyAsync 方法 ")

[DeleteKey 方法](../html/e2f0f305-d8bc-d515-3a62-cc91a3c1de76.htm "DeleteKey 方法 ")

[DeleteKeyAsync 方法](../html/0d8b5b0c-5979-9140-4ccf-992f4eaf1553.htm "DeleteKeyAsync 方法 ")

[ExistsHashKey 方法](../html/5b0bbeb9-4ab9-1f83-9e09-5f7227ac28e6.htm "ExistsHashKey 方法 ")

[ExistsHashKeyAsync 方法](../html/a80d408e-d621-a502-7829-6c7ea2a081f9.htm "ExistsHashKeyAsync 方法 ")

[ExistsKey 方法](../html/2b1b5b08-13ed-88c3-445d-d3613b48bc95.htm "ExistsKey 方法 ")

[ExistsKeyAsync 方法](../html/0d4357ed-5de3-1444-fce2-0e9be619d20f.htm "ExistsKeyAsync 方法 ")

[ExpireKey 方法](../html/0272044f-cc55-6b6a-42f0-6a0d89f37fe7.htm "ExpireKey 方法 ")

[ExpireKeyAsync 方法](../html/53c18b2d-2169-0519-602e-f65302f490e1.htm "ExpireKeyAsync 方法 ")

[ExtraOnDisconnect 方法](../html/49d9b263-02d8-e822-6a33-49f93737b98b.htm "ExtraOnDisconnect 方法 ")

[ExtraOnDisconnectAsync 方法](../html/f681350d-2aad-40e0-0824-30fbd3fe57ff.htm "ExtraOnDisconnectAsync 方法 ")

[FlushDB 方法](../html/38c28c3f-63bf-f1a0-67bc-6547799257f9.htm "FlushDB 方法 ")

[FlushDBAsync 方法](../html/40e3ccb2-9fe7-c0d7-88b9-95266a13c336.htm "FlushDBAsync 方法 ")

[GetListLength 方法](../html/36330604-2557-4fc1-bdd7-fc3aff05e997.htm "GetListLength 方法 ")

[GetListLengthAsync 方法](../html/0e09431d-45b3-b539-a618-a147ae377f88.htm "GetListLengthAsync 方法 ")

[IncrementHashKey 方法](../html/6404f402-9882-5453-7df7-275f61748021.htm "IncrementHashKey 方法 ")

[IncrementHashKeyAsync 方法](../html/7a69e1a1-8689-3e23-8aa3-cdad33e7e3f8.htm "IncrementHashKeyAsync 方法 ")

[IncrementKey 方法](../html/04d46420-bedd-9016-075c-3bb1a6623583.htm "IncrementKey 方法 ")

[IncrementKeyAsync 方法](../html/95f88e06-114c-b5b3-ee2f-2e9f97f27b8c.htm "IncrementKeyAsync 方法 ")

[InitializationOnConnect 方法](../html/72414fc8-b43a-3377-acdc-c59ae0628e11.htm "InitializationOnConnect 方法 ")

[InitializationOnConnectAsync 方法](../html/9156d090-6b1d-7381-c396-81491b4b6756.htm "InitializationOnConnectAsync 方法 ")

[ListInsertAfter 方法](../html/11b6eb28-b2ee-57a5-c724-a99df9101c71.htm "ListInsertAfter 方法 ")

[ListInsertAfterAsync 方法](../html/9defd0d6-5681-4c26-0973-14ecca4fc821.htm "ListInsertAfterAsync 方法 ")

[ListInsertBefore 方法](../html/b19eed6b-2778-4b91-6f6b-e75cd9f1088e.htm "ListInsertBefore 方法 ")

[ListInsertBeforeAsync 方法](../html/0bf5390c-af9f-ecd6-1148-34a20abee73d.htm "ListInsertBeforeAsync 方法 ")

[ListLeftPop 方法](../html/8865ed9c-69dc-2f28-7f53-42e1c03d152b.htm "ListLeftPop 方法 ")

[ListLeftPopAsync 方法](../html/f4d7f464-73f7-ff5d-0834-e18b5a91c994.htm "ListLeftPopAsync 方法 ")

[ListLeftPush 方法](../html/232a2dbd-e2c6-5fbf-af96-5269a10df4c9.htm "ListLeftPush 方法 ")

[ListLeftPushAsync 方法](../html/cb91133d-ce36-f4f7-a9d3-81dbda7d6617.htm "ListLeftPushAsync 方法 ")

[ListLeftPushX 方法](../html/33093159-9864-a0dd-6a09-8bf134942031.htm "ListLeftPushX 方法 ")

[ListLeftPushXAsync 方法](../html/6bd2ab31-ff58-71f0-b68f-f36d7f53d845.htm "ListLeftPushXAsync 方法 ")

[ListRange 方法](../html/efe06eb0-b945-fe65-5e81-21888cdda04f.htm "ListRange 方法 ")

[ListRangeAsync 方法](../html/32e344f9-cdac-fc0e-2c50-bf2dd54b4f98.htm "ListRangeAsync 方法 ")

[ListRemoveElementMatch 方法](../html/0ab40e89-d265-0b48-05de-5c1a7b41d326.htm "ListRemoveElementMatch 方法 ")

[ListRemoveElementMatchAsync 方法](../html/5610de76-3198-5b5a-8f3d-b8fc937d7c83.htm "ListRemoveElementMatchAsync 方法 ")

[ListRightPop 方法](../html/97855b86-6a07-76e2-b887-1e830be705de.htm "ListRightPop 方法 ")

[ListRightPopAsync 方法](../html/fde8cee2-f393-a463-d561-74568545edc7.htm "ListRightPopAsync 方法 ")

[ListRightPopLeftPush 方法](../html/e0b58cc4-4d80-64d6-729c-0dcdaafc4701.htm "ListRightPopLeftPush 方法 ")

[ListRightPopLeftPushAsync 方法](../html/41069df2-0dad-6f01-895f-a91086a92fa9.htm "ListRightPopLeftPushAsync 方法 ")

[ListRightPush 方法](../html/d71ca156-74d3-536e-6925-315d948c0ace.htm "ListRightPush 方法 ")

[ListRightPushAsync 方法](../html/53144654-da3d-b367-b39c-0dbc68275e59.htm "ListRightPushAsync 方法 ")

[ListRightPushX 方法](../html/e6ccfd1a-25ee-5eff-d56c-0ff9eb017fbc.htm "ListRightPushX 方法 ")

[ListRightPushXAsync 方法](../html/3252db2a-56cf-1aa2-7716-f9eb2596ec62.htm "ListRightPushXAsync 方法 ")

[ListSet 方法](../html/5f2e7b8b-b0e7-87e6-dac7-81b572cf3fa3.htm "ListSet 方法 ")

[ListSetAsync 方法](../html/ca019ec2-eb32-41f8-a460-89c70f90ace9.htm "ListSetAsync 方法 ")

[ListTrim 方法](../html/acf21364-1692-f9da-d337-c460b12a012d.htm "ListTrim 方法 ")

[ListTrimAsync 方法](../html/9e86492d-205c-8427-9052-50ea3d8c1c64.htm "ListTrimAsync 方法 ")

[MoveKey 方法](../html/5aa50043-85b9-a245-606a-de2b69164218.htm "MoveKey 方法 ")

[MoveKeyAsync 方法](../html/a93a74bf-65bd-733b-e39a-06d705252d99.htm "MoveKeyAsync 方法 ")

[OperateLongNumberFromServer 方法](../html/d2cdcd44-a6fb-7aa7-8c9c-08ff43107f65.htm "OperateLongNumberFromServer 方法 ")

[OperateLongNumberFromServerAsync 方法](../html/0005132d-e8c6-d51b-69dc-c9f8feaf8a9d.htm "OperateLongNumberFromServerAsync 方法 ")

[OperateNumberFromServer 方法](../html/69aecc05-6901-f53a-da37-ebf85bae7053.htm "OperateNumberFromServer 方法 ")

[OperateNumberFromServerAsync 方法](../html/b53bd6fa-d215-55f7-bbac-2ee171ed7498.htm "OperateNumberFromServerAsync 方法 ")

[OperateStatusFromServer 方法](../html/d460c3bb-e484-3fcb-e206-d20dc7201d63.htm "OperateStatusFromServer 方法 ")

[OperateStatusFromServerAsync 方法](../html/e1a241a3-e200-747f-bb8b-5b7578a08999.htm "OperateStatusFromServerAsync 方法 ")

[OperateStringFromServer 方法](../html/822914af-cf62-1ccc-d22e-2df7dfac7294.htm "OperateStringFromServer 方法 ")

[OperateStringFromServerAsync 方法](../html/eb31e52d-5387-eb14-291e-22f0e64666a4.htm "OperateStringFromServerAsync 方法 ")

[OperateStringsFromServer 方法](../html/c85c1a0e-83f9-69c9-a9b9-8a947ca54583.htm "OperateStringsFromServer 方法 ")

[OperateStringsFromServerAsync 方法](../html/2db50b51-32b8-1e92-cb9a-f2c847437e50.htm "OperateStringsFromServerAsync 方法 ")

[PersistKey 方法](../html/c68a840b-ce64-3303-9339-a10f88bcbb77.htm "PersistKey 方法 ")

[PersistKeyAsync 方法](../html/a5d7ee67-9041-6455-7fc1-c6e96a442799.htm "PersistKeyAsync 方法 ")

[Ping 方法](../html/9cdd4e99-90a2-0ae8-141c-26f0f4a4d698.htm "Ping 方法 ")

[PingAsync 方法](../html/d6c8e625-edb2-7ff2-a54d-e852d49c73c8.htm "PingAsync 方法 ")

[Publish 方法](../html/74f20db5-aaaa-af2e-b4e6-5d192196b222.htm "Publish 方法 ")

[PublishAsync 方法](../html/61428af8-94a2-b2b0-419b-5c0d3df4bee5.htm "PublishAsync 方法 ")

[Read(T) 方法](../html/08365a89-e2d7-2aa5-3fd5-ab1a5d60c440.htm "Read(T) 方法 ")

[ReadAllKeys 方法](../html/53c48696-f991-f734-12b0-1ad306a84c6b.htm "ReadAllKeys 方法 ")

[ReadAllKeysAsync 方法](../html/84c1be8c-3286-0410-0bd1-bb0185ff197e.htm "ReadAllKeysAsync 方法 ")

[ReadAndWriteKey 方法](../html/a6908f0e-1991-5210-f6c6-1e8b06369e56.htm "ReadAndWriteKey 方法 ")

[ReadAndWriteKeyAsync 方法](../html/6686726f-93ba-6435-76b0-b501c31e637e.htm "ReadAndWriteKeyAsync 方法 ")

[ReadAsync(T) 方法](../html/43a59914-1502-1353-0bf3-20281ac92ca9.htm "ReadAsync(T) 方法 ")

[ReadCustomer 方法](../html/a4b7f43e-34ba-8aac-a705-a9bc8d4a5298.htm "ReadCustomer 方法 ")

[ReadCustomerAsync 方法](../html/bafea5ed-4293-ae60-c8e7-bb3d3481db44.htm "ReadCustomerAsync 方法 ")

[ReadFromCoreServer 方法](../html/a5985ec7-e909-5a40-b2b9-7160363a1156.htm "ReadFromCoreServer 方法 ")

[ReadFromCoreServerAsync 方法](../html/984a89de-20da-e7ed-e471-3e889955aa31.htm "ReadFromCoreServerAsync 方法 ")

[ReadHashKey 方法](../html/6e8fc531-b2d6-2b1e-4885-1b937745652f.htm "ReadHashKey 方法 ")

[ReadHashKeyAll 方法](../html/63828fbe-fbdf-772a-010b-1425b0584572.htm "ReadHashKeyAll 方法 ")

[ReadHashKeyAllAsync 方法](../html/31c54864-f7bc-b3ce-7e77-15f0c1dcc540.htm "ReadHashKeyAllAsync 方法 ")

[ReadHashKeyAsync 方法](../html/36c995b4-42a7-e264-8a03-3390e24376ec.htm "ReadHashKeyAsync 方法 ")

[ReadHashKeyLength 方法](../html/32d483da-49e0-a205-bdc4-a659711be1c0.htm "ReadHashKeyLength 方法 ")

[ReadHashKeyLengthAsync 方法](../html/5d0d4570-1f16-4277-cff7-a40607a530ab.htm "ReadHashKeyLengthAsync 方法 ")

[ReadHashKeys 方法](../html/420286f8-5876-ee1b-b328-0b6eba7cefac.htm "ReadHashKeys 方法 ")

[ReadHashKeysAsync 方法](../html/bbc8310c-029d-1ea5-b0e3-239b29ee8743.htm "ReadHashKeysAsync 方法 ")

[ReadHashValues 方法](../html/1d8b8e06-3d34-f4c6-68d0-20c14b309de1.htm "ReadHashValues 方法 ")

[ReadHashValuesAsync 方法](../html/403e7344-c665-def6-70d9-81e5b22a740b.htm "ReadHashValuesAsync 方法 ")

[ReadKey 方法](../html/93fa3dab-d22c-d028-00c3-e81391085e6f.htm "ReadKey 方法 ")

[ReadKeyAsync 方法](../html/48cee995-f908-c3a2-ab70-829d789f58a5.htm "ReadKeyAsync 方法 ")

[ReadKeyLength 方法](../html/7746357b-7e1b-d150-d019-722b752f5eeb.htm "ReadKeyLength 方法 ")

[ReadKeyLengthAsync 方法](../html/7dd1b305-77b1-920c-14c0-2fe38e970791.htm "ReadKeyLengthAsync 方法 ")

[ReadKeyRange 方法](../html/2e3e9582-cb24-44e6-d78b-7c672acdd4c8.htm "ReadKeyRange 方法 ")

[ReadKeyRangeAsync 方法](../html/407fdd5b-8b57-75d2-63be-93294311cb30.htm "ReadKeyRangeAsync 方法 ")

[ReadKeyTTL 方法](../html/679079c3-36e5-1618-3409-18459a74a0c4.htm "ReadKeyTTL 方法 ")

[ReadKeyTTLAsync 方法](../html/24887926-0fd5-b65f-d070-1d439a5273ef.htm "ReadKeyTTLAsync 方法 ")

[ReadKeyType 方法](../html/1f702770-7804-880a-038e-c002ed31c42f.htm "ReadKeyType 方法 ")

[ReadKeyTypeAsync 方法](../html/91c104ff-a646-76de-9bed-5e64237a57fc.htm "ReadKeyTypeAsync 方法 ")

[ReadListByIndex 方法](../html/43c7dae8-4c01-d2c9-b38b-0e7bef68b581.htm "ReadListByIndex 方法 ")

[ReadListByIndexAsync 方法](../html/1e1163ad-8161-9390-cff2-7919c8145c92.htm "ReadListByIndexAsync 方法 ")

[ReadRandomKey 方法](../html/0d358a8f-b9e0-c659-c218-f051aaedcbfe.htm "ReadRandomKey 方法 ")

[ReadRandomKeyAsync 方法](../html/b97f91d2-2ab6-a39b-86f7-71238cd44fac.htm "ReadRandomKeyAsync 方法 ")

[ReadServerTime 方法](../html/c6c5f1fa-4639-0a32-8a42-48d56985e8e6.htm "ReadServerTime 方法 ")

[ReadServerTimeAsync 方法](../html/53022c3b-766e-c31b-1b55-548ceb3b2f42.htm "ReadServerTimeAsync 方法 ")

[RenameKey 方法](../html/ecf533af-233c-16b8-f88d-cb232e115013.htm "RenameKey 方法 ")

[RenameKeyAsync 方法](../html/20073ebe-0c16-2f57-a96b-79246e01129d.htm "RenameKeyAsync 方法 ")

[Save 方法](../html/261a8daa-9890-95bc-f8e4-02b78eb86277.htm "Save 方法 ")

[SaveAsync 方法](../html/88c2043f-e9e6-9d95-4a2e-455452cf8adf.htm "SaveAsync 方法 ")

[SelectDB 方法](../html/584a38b8-96e0-736f-7a24-981a8e920574.htm "SelectDB 方法 ")

[SelectDBAsync 方法](../html/462ed998-9938-a261-bde4-579b4b9d3dc1.htm "SelectDBAsync 方法 ")

[SetAdd 方法](../html/6f52901b-4d7f-70ca-1bf1-2e7146604dba.htm "SetAdd 方法 ")

[SetAddAsync 方法](../html/a59f7c04-0b20-35a4-a431-53299cb829f0.htm "SetAddAsync 方法 ")

[SetCard 方法](../html/cc0a3d1d-019b-f0ba-90a5-1bab767d6f05.htm "SetCard 方法 ")

[SetCardAsync 方法](../html/0d0665a9-a174-d154-f95b-2bf10d885bc3.htm "SetCardAsync 方法 ")

[SetDiff 方法](../html/bcc43604-dc32-a27f-53a4-a89ab29916be.htm "SetDiff 方法 ")

[SetDiffAsync 方法](../html/b70eae65-3537-3ab7-4cf9-ae6c314884c9.htm "SetDiffAsync 方法 ")

[SetDiffStore 方法](../html/7c356901-c174-c81b-50b9-ce8d60258d1e.htm "SetDiffStore 方法 ")

[SetDiffStoreAsync 方法](../html/697ac488-c2e1-8c5a-5e5d-845f676c85fd.htm "SetDiffStoreAsync 方法 ")

[SetInter 方法](../html/fa3a123d-447d-ce8b-5670-e47956d427d1.htm "SetInter 方法 ")

[SetInterAsync 方法](../html/622854f7-5e07-092e-42eb-06cc63580467.htm "SetInterAsync 方法 ")

[SetInterStore 方法](../html/996ff5a7-0d0b-5b17-546b-eb006fa5773f.htm "SetInterStore 方法 ")

[SetInterStoreAsync 方法](../html/115a5c5d-2c7b-4c9a-daec-92477fb3044a.htm "SetInterStoreAsync 方法 ")

[SetIsMember 方法](../html/1c6d465b-353a-34f2-8f46-9d4037e58911.htm "SetIsMember 方法 ")

[SetIsMemberAsync 方法](../html/21419054-3273-2aab-d486-633ab8861039.htm "SetIsMemberAsync 方法 ")

[SetMembers 方法](../html/7a5b277f-1fe1-141f-346b-3722d97ac7ca.htm "SetMembers 方法 ")

[SetMembersAsync 方法](../html/ce86155f-4be2-3836-3594-1398d78183ad.htm "SetMembersAsync 方法 ")

[SetMove 方法](../html/8b76c71a-c755-708b-e1e7-61eda8b9456c.htm "SetMove 方法 ")

[SetMoveAsync 方法](../html/dc01893c-b09e-e134-49d3-76d0dbc308b9.htm "SetMoveAsync 方法 ")

[SetPop 方法](../html/8375117d-d7b5-1f33-cb59-80cc632b736d.htm "SetPop 方法 ")

[SetPopAsync 方法](../html/66039533-6b1e-ad66-e89b-ef8f7d46c47b.htm "SetPopAsync 方法 ")

[SetRandomMember 方法](../html/2769a6d4-3b4c-75b2-c57b-abe19bb2f8e9.htm "SetRandomMember 方法 ")

[SetRandomMemberAsync 方法](../html/fb4fda18-00e6-2545-d037-6585a9487334.htm "SetRandomMemberAsync 方法 ")

[SetRemove 方法](../html/a0ba101e-935e-aca6-da5e-bb6187b42f02.htm "SetRemove 方法 ")

[SetRemoveAsync 方法](../html/f80e5db7-218b-c4bf-0c57-ec2555f6ec4e.htm "SetRemoveAsync 方法 ")

[SetUnion 方法](../html/c255e92b-1a91-f3e4-49cd-79f39d7a1079.htm "SetUnion 方法 ")

[SetUnionAsync 方法](../html/184ffefd-f6f5-b5d8-e1d6-45ae6d2d653f.htm "SetUnionAsync 方法 ")

[SetUnionStore 方法](../html/56394825-4bd1-6422-4f9b-1520ed05f927.htm "SetUnionStore 方法 ")

[SetUnionStoreAsync 方法](../html/99689f9d-ebcc-cc16-19e3-4c46f3016084.htm "SetUnionStoreAsync 方法 ")

[SubscribeMessage 方法](../html/bbaad51e-d35a-fb1f-4bb5-c001d47693be.htm "SubscribeMessage 方法 ")

[ToString 方法](../html/6ff40e11-12f4-8759-45a0-594ca3e4fcdd.htm "ToString 方法 ")

[UnSubscribeMessage 方法](../html/592327f8-7e50-4d48-c229-ca17a039b6a3.htm "UnSubscribeMessage 方法 ")

[Write(T) 方法](../html/70a08f37-0f86-76f1-55a5-df1ddd5a7d16.htm "Write(T) 方法 ")

[WriteAndPublishKey 方法](../html/69029e5e-a78a-5c83-0300-b6dd057f9308.htm "WriteAndPublishKey 方法 ")

[WriteAndPublishKeyAsync 方法](../html/d44abf4b-02bf-cfa5-cbc8-0b36edab4521.htm "WriteAndPublishKeyAsync 方法 ")

[WriteAsync(T) 方法](../html/9d54e352-2bd0-2655-3e14-49a0213e6040.htm "WriteAsync(T) 方法 ")

[WriteExpireKey 方法](../html/816ca8c0-5e3a-4609-acc3-6f227ac2dfde.htm "WriteExpireKey 方法 ")

[WriteExpireKeyAsync 方法](../html/f1dce55a-5681-32f6-2c7c-cc664465490c.htm "WriteExpireKeyAsync 方法 ")

[WriteHashKey 方法](../html/a7c80b1f-ec11-c8f4-dae2-0fa0acada685.htm "WriteHashKey 方法 ")

[WriteHashKeyAsync 方法](../html/20122681-134d-f220-683c-11c6c9201ac9.htm "WriteHashKeyAsync 方法 ")

[WriteHashKeyNx 方法](../html/819ddbb0-74d4-4667-6e55-372fac166c68.htm "WriteHashKeyNx 方法 ")

[WriteHashKeyNxAsync 方法](../html/b4571eb9-be29-fbb8-691b-f3a93d9245ec.htm "WriteHashKeyNxAsync 方法 ")

[WriteKey 方法](../html/12919427-9ab8-d765-27cc-0c3ee2c481d2.htm "WriteKey 方法 ")

[WriteKeyAsync 方法](../html/b136bcbb-3813-f660-54f1-8c5ccefb5ddf.htm "WriteKeyAsync 方法 ")

[WriteKeyIfNotExists 方法](../html/9d48312b-2ce1-077f-bf47-3e6b5a328c45.htm "WriteKeyIfNotExists 方法 ")

[WriteKeyIfNotExistsAsync 方法](../html/e6844e26-ba17-a0a8-0ff2-bb1a1e8bce8e.htm "WriteKeyIfNotExistsAsync 方法 ")

[WriteKeyRange 方法](../html/c1f98cc9-f791-4a35-bf67-ea29e8930d32.htm "WriteKeyRange 方法 ")

[WriteKeyRangeAsync 方法](../html/cdfaffec-518b-2754-5a10-3c5b68ee1767.htm "WriteKeyRangeAsync 方法 ")

[ZSetAdd 方法](../html/9c15cbbe-df53-5280-c9e4-b9d020222e0c.htm "ZSetAdd 方法 ")

[ZSetAddAsync 方法](../html/5318d2e6-d635-ed62-18e5-6573a530ab3e.htm "ZSetAddAsync 方法 ")

[ZSetCard 方法](../html/aa928601-c5f6-63fb-6bd8-e2e383796043.htm "ZSetCard 方法 ")

[ZSetCardAsync 方法](../html/b760ff15-4c7e-510f-1927-e17e6310e142.htm "ZSetCardAsync 方法 ")

[ZSetCount 方法](../html/9872ddc0-8b8d-d9d6-4d08-64001fdefd29.htm "ZSetCount 方法 ")

[ZSetCountAsync 方法](../html/b294e1ee-27c7-7326-eb4b-8669040aecf0.htm "ZSetCountAsync 方法 ")

[ZSetIncreaseBy 方法](../html/c84e0c14-15cd-451a-dec4-2a112a73772e.htm "ZSetIncreaseBy 方法 ")

[ZSetIncreaseByAsync 方法](../html/9805ac6b-aa05-2152-7dbb-7f8077ea462b.htm "ZSetIncreaseByAsync 方法 ")

[ZSetRange 方法](../html/759e2523-e2ad-0dc1-a904-354031a8ae50.htm "ZSetRange 方法 ")

[ZSetRangeAsync 方法](../html/7724eb62-2653-65a7-779d-391d630c2faf.htm "ZSetRangeAsync 方法 ")

[ZSetRangeByScore 方法](../html/acab27d0-6b6c-1f25-00ad-fafd4ac068fa.htm "ZSetRangeByScore 方法 ")

[ZSetRangeByScoreAsync 方法](../html/d1e0aa3c-d5f9-caeb-4017-699f956e661d.htm "ZSetRangeByScoreAsync 方法 ")

[ZSetRank 方法](../html/7957411f-32c1-4e07-a4f2-9d1948c63a2b.htm "ZSetRank 方法 ")

[ZSetRankAsync 方法](../html/1181fd8c-02aa-9d53-cbeb-35d542bed241.htm "ZSetRankAsync 方法 ")

[ZSetRemove 方法](../html/a76625ea-7a7b-161d-a039-02f4f5a8e997.htm "ZSetRemove 方法 ")

[ZSetRemoveAsync 方法](../html/2f1b3de6-ad74-e2b7-3350-79d9afe13c9d.htm "ZSetRemoveAsync 方法 ")

[ZSetRemoveRangeByRank 方法](../html/cb81dfc2-6e30-0aca-8bcb-348cdc32bc35.htm "ZSetRemoveRangeByRank 方法 ")

[ZSetRemoveRangeByRankAsync 方法](../html/9bde4b5f-b9bd-e01b-feb4-82c9b532dee4.htm "ZSetRemoveRangeByRankAsync 方法 ")

[ZSetRemoveRangeByScore 方法](../html/ac7be899-fe4e-3f12-d60f-81a038a927c7.htm "ZSetRemoveRangeByScore 方法 ")

[ZSetRemoveRangeByScoreAsync 方法](../html/6ccee6bf-8635-14bf-2e7f-d1ac351c70c9.htm "ZSetRemoveRangeByScoreAsync 方法 ")

[ZSetReverseRange 方法](../html/ef6defc7-62a2-c8a4-4025-32f6b3a82ae4.htm "ZSetReverseRange 方法 ")

[ZSetReverseRangeAsync 方法](../html/b40d37b7-1118-6aec-c7da-85ec54b457d0.htm "ZSetReverseRangeAsync 方法 ")

[ZSetReverseRangeByScore 方法](../html/3c1a26a8-f8e5-a835-7051-94b4d9c851ab.htm "ZSetReverseRangeByScore 方法 ")

[ZSetReverseRangeByScoreAsync 方法](../html/3a3cfbff-cee7-e1f3-d9c1-e7528db4e998.htm "ZSetReverseRangeByScoreAsync 方法 ")

[ZSetReverseRank 方法](../html/d524195b-6b02-a8e8-330a-3454a6ed21cd.htm "ZSetReverseRank 方法 ")

[ZSetReverseRankAsync 方法](../html/c010781c-2851-7158-f953-1779929e3242.htm "ZSetReverseRankAsync 方法 ")

[ZSetScore 方法](../html/a98ee7e4-f361-f403-975b-a197e0ca0268.htm "ZSetScore 方法 ")

[ZSetScoreAsync 方法](../html/fd24c433-944e-2ecb-4bcd-266bfc6825d4.htm "ZSetScoreAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| RedisClientAppendKey 方法 |

如果 key 已经存在并且是一个字符串， APPEND 命令将 value 追加到 key 原来的值的末尾。
如果 key 不存在， APPEND 就简单地将给定 key 设为 value ，就像执行 SET key value 一样。
返回追加 value 之后， key 中字符串的长度。

**命名空间：**
 [HslCommunication.Enthernet.Redis](cc44ce00-3781-c690-b233-24000870bf42.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public OperateResult<int> AppendKey(
	string key,
	string value
)
```

```
Public Function AppendKey ( 
	key As String,
	value As String
) As OperateResult(Of Integer)
```

```
public:
OperateResult<int>^ AppendKey(
	String^ key, 
	String^ value
)
```

```
member AppendKey : 
        key : string * 
        value : string -> OperateResult<int> 
```

#### 参数

key
:   类型：SystemString  
    关键字

value
:   类型：SystemString  
    数值

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Int32  
追加 value 之后， key 中字符串的长度。

![](../icons/SectionExpanded.png)参见

#### 引用

[RedisClient 类](e33d9088-6df0-32e9-33f1-426e944c3b7f.htm)

[HslCommunication.Enthernet.Redis 命名空间](cc44ce00-3781-c690-b233-24000870bf42.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## AppendKeyAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/c39b4c6a-b7f6-68fe-dfc9-7cce1b653e24.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Enthernet.Redis](../html/cc44ce00-3781-c690-b233-24000870bf42.htm "HslCommunication.Enthernet.Redis")

[RedisClient 类](../html/e33d9088-6df0-32e9-33f1-426e944c3b7f.htm "RedisClient 类")

[RedisClient 方法](../html/e41be4a5-1490-ffbd-f241-938f2d6c3f07.htm "RedisClient 方法")

[AppendKey 方法](../html/f2521b8c-0fb1-6b34-4adc-ca6804c7195f.htm "AppendKey 方法 ")

[AppendKeyAsync 方法](../html/c39b4c6a-b7f6-68fe-dfc9-7cce1b653e24.htm "AppendKeyAsync 方法 ")

[ChangePassword 方法](../html/78a37d4d-bdee-7a09-d0ce-dc245627b1d1.htm "ChangePassword 方法 ")

[ChangePasswordAsync 方法](../html/58f6f736-093d-6e4d-3ad3-19590fa2a3a3.htm "ChangePasswordAsync 方法 ")

[DBSize 方法](../html/5814129a-9a73-25ac-e726-fe1da141eb28.htm "DBSize 方法 ")

[DBSizeAsync 方法](../html/243a1f5f-a474-2e20-1298-06eff125f8b5.htm "DBSizeAsync 方法 ")

[DecrementKey 方法](../html/e9fa6c94-43ca-8291-79ca-6e04ca9289e7.htm "DecrementKey 方法 ")

[DecrementKeyAsync 方法](../html/6837e351-b8ef-a90e-c5e2-73b31621f837.htm "DecrementKeyAsync 方法 ")

[DeleteHashKey 方法](../html/19e3ebc3-e6fc-71cf-e80d-cb73e3011170.htm "DeleteHashKey 方法 ")

[DeleteHashKeyAsync 方法](../html/0dcc0d8a-1ed1-5c81-fb0e-e1c93be9d579.htm "DeleteHashKeyAsync 方法 ")

[DeleteKey 方法](../html/e2f0f305-d8bc-d515-3a62-cc91a3c1de76.htm "DeleteKey 方法 ")

[DeleteKeyAsync 方法](../html/0d8b5b0c-5979-9140-4ccf-992f4eaf1553.htm "DeleteKeyAsync 方法 ")

[ExistsHashKey 方法](../html/5b0bbeb9-4ab9-1f83-9e09-5f7227ac28e6.htm "ExistsHashKey 方法 ")

[ExistsHashKeyAsync 方法](../html/a80d408e-d621-a502-7829-6c7ea2a081f9.htm "ExistsHashKeyAsync 方法 ")

[ExistsKey 方法](../html/2b1b5b08-13ed-88c3-445d-d3613b48bc95.htm "ExistsKey 方法 ")

[ExistsKeyAsync 方法](../html/0d4357ed-5de3-1444-fce2-0e9be619d20f.htm "ExistsKeyAsync 方法 ")

[ExpireKey 方法](../html/0272044f-cc55-6b6a-42f0-6a0d89f37fe7.htm "ExpireKey 方法 ")

[ExpireKeyAsync 方法](../html/53c18b2d-2169-0519-602e-f65302f490e1.htm "ExpireKeyAsync 方法 ")

[ExtraOnDisconnect 方法](../html/49d9b263-02d8-e822-6a33-49f93737b98b.htm "ExtraOnDisconnect 方法 ")

[ExtraOnDisconnectAsync 方法](../html/f681350d-2aad-40e0-0824-30fbd3fe57ff.htm "ExtraOnDisconnectAsync 方法 ")

[FlushDB 方法](../html/38c28c3f-63bf-f1a0-67bc-6547799257f9.htm "FlushDB 方法 ")

[FlushDBAsync 方法](../html/40e3ccb2-9fe7-c0d7-88b9-95266a13c336.htm "FlushDBAsync 方法 ")

[GetListLength 方法](../html/36330604-2557-4fc1-bdd7-fc3aff05e997.htm "GetListLength 方法 ")

[GetListLengthAsync 方法](../html/0e09431d-45b3-b539-a618-a147ae377f88.htm "GetListLengthAsync 方法 ")

[IncrementHashKey 方法](../html/6404f402-9882-5453-7df7-275f61748021.htm "IncrementHashKey 方法 ")

[IncrementHashKeyAsync 方法](../html/7a69e1a1-8689-3e23-8aa3-cdad33e7e3f8.htm "IncrementHashKeyAsync 方法 ")

[IncrementKey 方法](../html/04d46420-bedd-9016-075c-3bb1a6623583.htm "IncrementKey 方法 ")

[IncrementKeyAsync 方法](../html/95f88e06-114c-b5b3-ee2f-2e9f97f27b8c.htm "IncrementKeyAsync 方法 ")

[InitializationOnConnect 方法](../html/72414fc8-b43a-3377-acdc-c59ae0628e11.htm "InitializationOnConnect 方法 ")

[InitializationOnConnectAsync 方法](../html/9156d090-6b1d-7381-c396-81491b4b6756.htm "InitializationOnConnectAsync 方法 ")

[ListInsertAfter 方法](../html/11b6eb28-b2ee-57a5-c724-a99df9101c71.htm "ListInsertAfter 方法 ")

[ListInsertAfterAsync 方法](../html/9defd0d6-5681-4c26-0973-14ecca4fc821.htm "ListInsertAfterAsync 方法 ")

[ListInsertBefore 方法](../html/b19eed6b-2778-4b91-6f6b-e75cd9f1088e.htm "ListInsertBefore 方法 ")

[ListInsertBeforeAsync 方法](../html/0bf5390c-af9f-ecd6-1148-34a20abee73d.htm "ListInsertBeforeAsync 方法 ")

[ListLeftPop 方法](../html/8865ed9c-69dc-2f28-7f53-42e1c03d152b.htm "ListLeftPop 方法 ")

[ListLeftPopAsync 方法](../html/f4d7f464-73f7-ff5d-0834-e18b5a91c994.htm "ListLeftPopAsync 方法 ")

[ListLeftPush 方法](../html/232a2dbd-e2c6-5fbf-af96-5269a10df4c9.htm "ListLeftPush 方法 ")

[ListLeftPushAsync 方法](../html/cb91133d-ce36-f4f7-a9d3-81dbda7d6617.htm "ListLeftPushAsync 方法 ")

[ListLeftPushX 方法](../html/33093159-9864-a0dd-6a09-8bf134942031.htm "ListLeftPushX 方法 ")

[ListLeftPushXAsync 方法](../html/6bd2ab31-ff58-71f0-b68f-f36d7f53d845.htm "ListLeftPushXAsync 方法 ")

[ListRange 方法](../html/efe06eb0-b945-fe65-5e81-21888cdda04f.htm "ListRange 方法 ")

[ListRangeAsync 方法](../html/32e344f9-cdac-fc0e-2c50-bf2dd54b4f98.htm "ListRangeAsync 方法 ")

[ListRemoveElementMatch 方法](../html/0ab40e89-d265-0b48-05de-5c1a7b41d326.htm "ListRemoveElementMatch 方法 ")

[ListRemoveElementMatchAsync 方法](../html/5610de76-3198-5b5a-8f3d-b8fc937d7c83.htm "ListRemoveElementMatchAsync 方法 ")

[ListRightPop 方法](../html/97855b86-6a07-76e2-b887-1e830be705de.htm "ListRightPop 方法 ")

[ListRightPopAsync 方法](../html/fde8cee2-f393-a463-d561-74568545edc7.htm "ListRightPopAsync 方法 ")

[ListRightPopLeftPush 方法](../html/e0b58cc4-4d80-64d6-729c-0dcdaafc4701.htm "ListRightPopLeftPush 方法 ")

[ListRightPopLeftPushAsync 方法](../html/41069df2-0dad-6f01-895f-a91086a92fa9.htm "ListRightPopLeftPushAsync 方法 ")

[ListRightPush 方法](../html/d71ca156-74d3-536e-6925-315d948c0ace.htm "ListRightPush 方法 ")

[ListRightPushAsync 方法](../html/53144654-da3d-b367-b39c-0dbc68275e59.htm "ListRightPushAsync 方法 ")

[ListRightPushX 方法](../html/e6ccfd1a-25ee-5eff-d56c-0ff9eb017fbc.htm "ListRightPushX 方法 ")

[ListRightPushXAsync 方法](../html/3252db2a-56cf-1aa2-7716-f9eb2596ec62.htm "ListRightPushXAsync 方法 ")

[ListSet 方法](../html/5f2e7b8b-b0e7-87e6-dac7-81b572cf3fa3.htm "ListSet 方法 ")

[ListSetAsync 方法](../html/ca019ec2-eb32-41f8-a460-89c70f90ace9.htm "ListSetAsync 方法 ")

[ListTrim 方法](../html/acf21364-1692-f9da-d337-c460b12a012d.htm "ListTrim 方法 ")

[ListTrimAsync 方法](../html/9e86492d-205c-8427-9052-50ea3d8c1c64.htm "ListTrimAsync 方法 ")

[MoveKey 方法](../html/5aa50043-85b9-a245-606a-de2b69164218.htm "MoveKey 方法 ")

[MoveKeyAsync 方法](../html/a93a74bf-65bd-733b-e39a-06d705252d99.htm "MoveKeyAsync 方法 ")

[OperateLongNumberFromServer 方法](../html/d2cdcd44-a6fb-7aa7-8c9c-08ff43107f65.htm "OperateLongNumberFromServer 方法 ")

[OperateLongNumberFromServerAsync 方法](../html/0005132d-e8c6-d51b-69dc-c9f8feaf8a9d.htm "OperateLongNumberFromServerAsync 方法 ")

[OperateNumberFromServer 方法](../html/69aecc05-6901-f53a-da37-ebf85bae7053.htm "OperateNumberFromServer 方法 ")

[OperateNumberFromServerAsync 方法](../html/b53bd6fa-d215-55f7-bbac-2ee171ed7498.htm "OperateNumberFromServerAsync 方法 ")

[OperateStatusFromServer 方法](../html/d460c3bb-e484-3fcb-e206-d20dc7201d63.htm "OperateStatusFromServer 方法 ")

[OperateStatusFromServerAsync 方法](../html/e1a241a3-e200-747f-bb8b-5b7578a08999.htm "OperateStatusFromServerAsync 方法 ")

[OperateStringFromServer 方法](../html/822914af-cf62-1ccc-d22e-2df7dfac7294.htm "OperateStringFromServer 方法 ")

[OperateStringFromServerAsync 方法](../html/eb31e52d-5387-eb14-291e-22f0e64666a4.htm "OperateStringFromServerAsync 方法 ")

[OperateStringsFromServer 方法](../html/c85c1a0e-83f9-69c9-a9b9-8a947ca54583.htm "OperateStringsFromServer 方法 ")

[OperateStringsFromServerAsync 方法](../html/2db50b51-32b8-1e92-cb9a-f2c847437e50.htm "OperateStringsFromServerAsync 方法 ")

[PersistKey 方法](../html/c68a840b-ce64-3303-9339-a10f88bcbb77.htm "PersistKey 方法 ")

[PersistKeyAsync 方法](../html/a5d7ee67-9041-6455-7fc1-c6e96a442799.htm "PersistKeyAsync 方法 ")

[Ping 方法](../html/9cdd4e99-90a2-0ae8-141c-26f0f4a4d698.htm "Ping 方法 ")

[PingAsync 方法](../html/d6c8e625-edb2-7ff2-a54d-e852d49c73c8.htm "PingAsync 方法 ")

[Publish 方法](../html/74f20db5-aaaa-af2e-b4e6-5d192196b222.htm "Publish 方法 ")

[PublishAsync 方法](../html/61428af8-94a2-b2b0-419b-5c0d3df4bee5.htm "PublishAsync 方法 ")

[Read(T) 方法](../html/08365a89-e2d7-2aa5-3fd5-ab1a5d60c440.htm "Read(T) 方法 ")

[ReadAllKeys 方法](../html/53c48696-f991-f734-12b0-1ad306a84c6b.htm "ReadAllKeys 方法 ")

[ReadAllKeysAsync 方法](../html/84c1be8c-3286-0410-0bd1-bb0185ff197e.htm "ReadAllKeysAsync 方法 ")

[ReadAndWriteKey 方法](../html/a6908f0e-1991-5210-f6c6-1e8b06369e56.htm "ReadAndWriteKey 方法 ")

[ReadAndWriteKeyAsync 方法](../html/6686726f-93ba-6435-76b0-b501c31e637e.htm "ReadAndWriteKeyAsync 方法 ")

[ReadAsync(T) 方法](../html/43a59914-1502-1353-0bf3-20281ac92ca9.htm "ReadAsync(T) 方法 ")

[ReadCustomer 方法](../html/a4b7f43e-34ba-8aac-a705-a9bc8d4a5298.htm "ReadCustomer 方法 ")

[ReadCustomerAsync 方法](../html/bafea5ed-4293-ae60-c8e7-bb3d3481db44.htm "ReadCustomerAsync 方法 ")

[ReadFromCoreServer 方法](../html/a5985ec7-e909-5a40-b2b9-7160363a1156.htm "ReadFromCoreServer 方法 ")

[ReadFromCoreServerAsync 方法](../html/984a89de-20da-e7ed-e471-3e889955aa31.htm "ReadFromCoreServerAsync 方法 ")

[ReadHashKey 方法](../html/6e8fc531-b2d6-2b1e-4885-1b937745652f.htm "ReadHashKey 方法 ")

[ReadHashKeyAll 方法](../html/63828fbe-fbdf-772a-010b-1425b0584572.htm "ReadHashKeyAll 方法 ")

[ReadHashKeyAllAsync 方法](../html/31c54864-f7bc-b3ce-7e77-15f0c1dcc540.htm "ReadHashKeyAllAsync 方法 ")

[ReadHashKeyAsync 方法](../html/36c995b4-42a7-e264-8a03-3390e24376ec.htm "ReadHashKeyAsync 方法 ")

[ReadHashKeyLength 方法](../html/32d483da-49e0-a205-bdc4-a659711be1c0.htm "ReadHashKeyLength 方法 ")

[ReadHashKeyLengthAsync 方法](../html/5d0d4570-1f16-4277-cff7-a40607a530ab.htm "ReadHashKeyLengthAsync 方法 ")

[ReadHashKeys 方法](../html/420286f8-5876-ee1b-b328-0b6eba7cefac.htm "ReadHashKeys 方法 ")

[ReadHashKeysAsync 方法](../html/bbc8310c-029d-1ea5-b0e3-239b29ee8743.htm "ReadHashKeysAsync 方法 ")

[ReadHashValues 方法](../html/1d8b8e06-3d34-f4c6-68d0-20c14b309de1.htm "ReadHashValues 方法 ")

[ReadHashValuesAsync 方法](../html/403e7344-c665-def6-70d9-81e5b22a740b.htm "ReadHashValuesAsync 方法 ")

[ReadKey 方法](../html/93fa3dab-d22c-d028-00c3-e81391085e6f.htm "ReadKey 方法 ")

[ReadKeyAsync 方法](../html/48cee995-f908-c3a2-ab70-829d789f58a5.htm "ReadKeyAsync 方法 ")

[ReadKeyLength 方法](../html/7746357b-7e1b-d150-d019-722b752f5eeb.htm "ReadKeyLength 方法 ")

[ReadKeyLengthAsync 方法](../html/7dd1b305-77b1-920c-14c0-2fe38e970791.htm "ReadKeyLengthAsync 方法 ")

[ReadKeyRange 方法](../html/2e3e9582-cb24-44e6-d78b-7c672acdd4c8.htm "ReadKeyRange 方法 ")

[ReadKeyRangeAsync 方法](../html/407fdd5b-8b57-75d2-63be-93294311cb30.htm "ReadKeyRangeAsync 方法 ")

[ReadKeyTTL 方法](../html/679079c3-36e5-1618-3409-18459a74a0c4.htm "ReadKeyTTL 方法 ")

[ReadKeyTTLAsync 方法](../html/24887926-0fd5-b65f-d070-1d439a5273ef.htm "ReadKeyTTLAsync 方法 ")

[ReadKeyType 方法](../html/1f702770-7804-880a-038e-c002ed31c42f.htm "ReadKeyType 方法 ")

[ReadKeyTypeAsync 方法](../html/91c104ff-a646-76de-9bed-5e64237a57fc.htm "ReadKeyTypeAsync 方法 ")

[ReadListByIndex 方法](../html/43c7dae8-4c01-d2c9-b38b-0e7bef68b581.htm "ReadListByIndex 方法 ")

[ReadListByIndexAsync 方法](../html/1e1163ad-8161-9390-cff2-7919c8145c92.htm "ReadListByIndexAsync 方法 ")

[ReadRandomKey 方法](../html/0d358a8f-b9e0-c659-c218-f051aaedcbfe.htm "ReadRandomKey 方法 ")

[ReadRandomKeyAsync 方法](../html/b97f91d2-2ab6-a39b-86f7-71238cd44fac.htm "ReadRandomKeyAsync 方法 ")

[ReadServerTime 方法](../html/c6c5f1fa-4639-0a32-8a42-48d56985e8e6.htm "ReadServerTime 方法 ")

[ReadServerTimeAsync 方法](../html/53022c3b-766e-c31b-1b55-548ceb3b2f42.htm "ReadServerTimeAsync 方法 ")

[RenameKey 方法](../html/ecf533af-233c-16b8-f88d-cb232e115013.htm "RenameKey 方法 ")

[RenameKeyAsync 方法](../html/20073ebe-0c16-2f57-a96b-79246e01129d.htm "RenameKeyAsync 方法 ")

[Save 方法](../html/261a8daa-9890-95bc-f8e4-02b78eb86277.htm "Save 方法 ")

[SaveAsync 方法](../html/88c2043f-e9e6-9d95-4a2e-455452cf8adf.htm "SaveAsync 方法 ")

[SelectDB 方法](../html/584a38b8-96e0-736f-7a24-981a8e920574.htm "SelectDB 方法 ")

[SelectDBAsync 方法](../html/462ed998-9938-a261-bde4-579b4b9d3dc1.htm "SelectDBAsync 方法 ")

[SetAdd 方法](../html/6f52901b-4d7f-70ca-1bf1-2e7146604dba.htm "SetAdd 方法 ")

[SetAddAsync 方法](../html/a59f7c04-0b20-35a4-a431-53299cb829f0.htm "SetAddAsync 方法 ")

[SetCard 方法](../html/cc0a3d1d-019b-f0ba-90a5-1bab767d6f05.htm "SetCard 方法 ")

[SetCardAsync 方法](../html/0d0665a9-a174-d154-f95b-2bf10d885bc3.htm "SetCardAsync 方法 ")

[SetDiff 方法](../html/bcc43604-dc32-a27f-53a4-a89ab29916be.htm "SetDiff 方法 ")

[SetDiffAsync 方法](../html/b70eae65-3537-3ab7-4cf9-ae6c314884c9.htm "SetDiffAsync 方法 ")

[SetDiffStore 方法](../html/7c356901-c174-c81b-50b9-ce8d60258d1e.htm "SetDiffStore 方法 ")

[SetDiffStoreAsync 方法](../html/697ac488-c2e1-8c5a-5e5d-845f676c85fd.htm "SetDiffStoreAsync 方法 ")

[SetInter 方法](../html/fa3a123d-447d-ce8b-5670-e47956d427d1.htm "SetInter 方法 ")

[SetInterAsync 方法](../html/622854f7-5e07-092e-42eb-06cc63580467.htm "SetInterAsync 方法 ")

[SetInterStore 方法](../html/996ff5a7-0d0b-5b17-546b-eb006fa5773f.htm "SetInterStore 方法 ")

[SetInterStoreAsync 方法](../html/115a5c5d-2c7b-4c9a-daec-92477fb3044a.htm "SetInterStoreAsync 方法 ")

[SetIsMember 方法](../html/1c6d465b-353a-34f2-8f46-9d4037e58911.htm "SetIsMember 方法 ")

[SetIsMemberAsync 方法](../html/21419054-3273-2aab-d486-633ab8861039.htm "SetIsMemberAsync 方法 ")

[SetMembers 方法](../html/7a5b277f-1fe1-141f-346b-3722d97ac7ca.htm "SetMembers 方法 ")

[SetMembersAsync 方法](../html/ce86155f-4be2-3836-3594-1398d78183ad.htm "SetMembersAsync 方法 ")

[SetMove 方法](../html/8b76c71a-c755-708b-e1e7-61eda8b9456c.htm "SetMove 方法 ")

[SetMoveAsync 方法](../html/dc01893c-b09e-e134-49d3-76d0dbc308b9.htm "SetMoveAsync 方法 ")

[SetPop 方法](../html/8375117d-d7b5-1f33-cb59-80cc632b736d.htm "SetPop 方法 ")

[SetPopAsync 方法](../html/66039533-6b1e-ad66-e89b-ef8f7d46c47b.htm "SetPopAsync 方法 ")

[SetRandomMember 方法](../html/2769a6d4-3b4c-75b2-c57b-abe19bb2f8e9.htm "SetRandomMember 方法 ")

[SetRandomMemberAsync 方法](../html/fb4fda18-00e6-2545-d037-6585a9487334.htm "SetRandomMemberAsync 方法 ")

[SetRemove 方法](../html/a0ba101e-935e-aca6-da5e-bb6187b42f02.htm "SetRemove 方法 ")

[SetRemoveAsync 方法](../html/f80e5db7-218b-c4bf-0c57-ec2555f6ec4e.htm "SetRemoveAsync 方法 ")

[SetUnion 方法](../html/c255e92b-1a91-f3e4-49cd-79f39d7a1079.htm "SetUnion 方法 ")

[SetUnionAsync 方法](../html/184ffefd-f6f5-b5d8-e1d6-45ae6d2d653f.htm "SetUnionAsync 方法 ")

[SetUnionStore 方法](../html/56394825-4bd1-6422-4f9b-1520ed05f927.htm "SetUnionStore 方法 ")

[SetUnionStoreAsync 方法](../html/99689f9d-ebcc-cc16-19e3-4c46f3016084.htm "SetUnionStoreAsync 方法 ")

[SubscribeMessage 方法](../html/bbaad51e-d35a-fb1f-4bb5-c001d47693be.htm "SubscribeMessage 方法 ")

[ToString 方法](../html/6ff40e11-12f4-8759-45a0-594ca3e4fcdd.htm "ToString 方法 ")

[UnSubscribeMessage 方法](../html/592327f8-7e50-4d48-c229-ca17a039b6a3.htm "UnSubscribeMessage 方法 ")

[Write(T) 方法](../html/70a08f37-0f86-76f1-55a5-df1ddd5a7d16.htm "Write(T) 方法 ")

[WriteAndPublishKey 方法](../html/69029e5e-a78a-5c83-0300-b6dd057f9308.htm "WriteAndPublishKey 方法 ")

[WriteAndPublishKeyAsync 方法](../html/d44abf4b-02bf-cfa5-cbc8-0b36edab4521.htm "WriteAndPublishKeyAsync 方法 ")

[WriteAsync(T) 方法](../html/9d54e352-2bd0-2655-3e14-49a0213e6040.htm "WriteAsync(T) 方法 ")

[WriteExpireKey 方法](../html/816ca8c0-5e3a-4609-acc3-6f227ac2dfde.htm "WriteExpireKey 方法 ")

[WriteExpireKeyAsync 方法](../html/f1dce55a-5681-32f6-2c7c-cc664465490c.htm "WriteExpireKeyAsync 方法 ")

[WriteHashKey 方法](../html/a7c80b1f-ec11-c8f4-dae2-0fa0acada685.htm "WriteHashKey 方法 ")

[WriteHashKeyAsync 方法](../html/20122681-134d-f220-683c-11c6c9201ac9.htm "WriteHashKeyAsync 方法 ")

[WriteHashKeyNx 方法](../html/819ddbb0-74d4-4667-6e55-372fac166c68.htm "WriteHashKeyNx 方法 ")

[WriteHashKeyNxAsync 方法](../html/b4571eb9-be29-fbb8-691b-f3a93d9245ec.htm "WriteHashKeyNxAsync 方法 ")

[WriteKey 方法](../html/12919427-9ab8-d765-27cc-0c3ee2c481d2.htm "WriteKey 方法 ")

[WriteKeyAsync 方法](../html/b136bcbb-3813-f660-54f1-8c5ccefb5ddf.htm "WriteKeyAsync 方法 ")

[WriteKeyIfNotExists 方法](../html/9d48312b-2ce1-077f-bf47-3e6b5a328c45.htm "WriteKeyIfNotExists 方法 ")

[WriteKeyIfNotExistsAsync 方法](../html/e6844e26-ba17-a0a8-0ff2-bb1a1e8bce8e.htm "WriteKeyIfNotExistsAsync 方法 ")

[WriteKeyRange 方法](../html/c1f98cc9-f791-4a35-bf67-ea29e8930d32.htm "WriteKeyRange 方法 ")

[WriteKeyRangeAsync 方法](../html/cdfaffec-518b-2754-5a10-3c5b68ee1767.htm "WriteKeyRangeAsync 方法 ")

[ZSetAdd 方法](../html/9c15cbbe-df53-5280-c9e4-b9d020222e0c.htm "ZSetAdd 方法 ")

[ZSetAddAsync 方法](../html/5318d2e6-d635-ed62-18e5-6573a530ab3e.htm "ZSetAddAsync 方法 ")

[ZSetCard 方法](../html/aa928601-c5f6-63fb-6bd8-e2e383796043.htm "ZSetCard 方法 ")

[ZSetCardAsync 方法](../html/b760ff15-4c7e-510f-1927-e17e6310e142.htm "ZSetCardAsync 方法 ")

[ZSetCount 方法](../html/9872ddc0-8b8d-d9d6-4d08-64001fdefd29.htm "ZSetCount 方法 ")

[ZSetCountAsync 方法](../html/b294e1ee-27c7-7326-eb4b-8669040aecf0.htm "ZSetCountAsync 方法 ")

[ZSetIncreaseBy 方法](../html/c84e0c14-15cd-451a-dec4-2a112a73772e.htm "ZSetIncreaseBy 方法 ")

[ZSetIncreaseByAsync 方法](../html/9805ac6b-aa05-2152-7dbb-7f8077ea462b.htm "ZSetIncreaseByAsync 方法 ")

[ZSetRange 方法](../html/759e2523-e2ad-0dc1-a904-354031a8ae50.htm "ZSetRange 方法 ")

[ZSetRangeAsync 方法](../html/7724eb62-2653-65a7-779d-391d630c2faf.htm "ZSetRangeAsync 方法 ")

[ZSetRangeByScore 方法](../html/acab27d0-6b6c-1f25-00ad-fafd4ac068fa.htm "ZSetRangeByScore 方法 ")

[ZSetRangeByScoreAsync 方法](../html/d1e0aa3c-d5f9-caeb-4017-699f956e661d.htm "ZSetRangeByScoreAsync 方法 ")

[ZSetRank 方法](../html/7957411f-32c1-4e07-a4f2-9d1948c63a2b.htm "ZSetRank 方法 ")

[ZSetRankAsync 方法](../html/1181fd8c-02aa-9d53-cbeb-35d542bed241.htm "ZSetRankAsync 方法 ")

[ZSetRemove 方法](../html/a76625ea-7a7b-161d-a039-02f4f5a8e997.htm "ZSetRemove 方法 ")

[ZSetRemoveAsync 方法](../html/2f1b3de6-ad74-e2b7-3350-79d9afe13c9d.htm "ZSetRemoveAsync 方法 ")

[ZSetRemoveRangeByRank 方法](../html/cb81dfc2-6e30-0aca-8bcb-348cdc32bc35.htm "ZSetRemoveRangeByRank 方法 ")

[ZSetRemoveRangeByRankAsync 方法](../html/9bde4b5f-b9bd-e01b-feb4-82c9b532dee4.htm "ZSetRemoveRangeByRankAsync 方法 ")

[ZSetRemoveRangeByScore 方法](../html/ac7be899-fe4e-3f12-d60f-81a038a927c7.htm "ZSetRemoveRangeByScore 方法 ")

[ZSetRemoveRangeByScoreAsync 方法](../html/6ccee6bf-8635-14bf-2e7f-d1ac351c70c9.htm "ZSetRemoveRangeByScoreAsync 方法 ")

[ZSetReverseRange 方法](../html/ef6defc7-62a2-c8a4-4025-32f6b3a82ae4.htm "ZSetReverseRange 方法 ")

[ZSetReverseRangeAsync 方法](../html/b40d37b7-1118-6aec-c7da-85ec54b457d0.htm "ZSetReverseRangeAsync 方法 ")

[ZSetReverseRangeByScore 方法](../html/3c1a26a8-f8e5-a835-7051-94b4d9c851ab.htm "ZSetReverseRangeByScore 方法 ")

[ZSetReverseRangeByScoreAsync 方法](../html/3a3cfbff-cee7-e1f3-d9c1-e7528db4e998.htm "ZSetReverseRangeByScoreAsync 方法 ")

[ZSetReverseRank 方法](../html/d524195b-6b02-a8e8-330a-3454a6ed21cd.htm "ZSetReverseRank 方法 ")

[ZSetReverseRankAsync 方法](../html/c010781c-2851-7158-f953-1779929e3242.htm "ZSetReverseRankAsync 方法 ")

[ZSetScore 方法](../html/a98ee7e4-f361-f403-975b-a197e0ca0268.htm "ZSetScore 方法 ")

[ZSetScoreAsync 方法](../html/fd24c433-944e-2ecb-4bcd-266bfc6825d4.htm "ZSetScoreAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| RedisClientAppendKeyAsync 方法 |

如果 key 已经存在并且是一个字符串， APPEND 命令将 value 追加到 key 原来的值的末尾。
如果 key 不存在， APPEND 就简单地将给定 key 设为 value ，就像执行 SET key value 一样。
返回追加 value 之后， key 中字符串的长度。

**命名空间：**
 [HslCommunication.Enthernet.Redis](cc44ce00-3781-c690-b233-24000870bf42.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public Task<OperateResult<int>> AppendKeyAsync(
	string key,
	string value
)
```

```
Public Function AppendKeyAsync ( 
	key As String,
	value As String
) As Task(Of OperateResult(Of Integer))
```

```
public:
Task<OperateResult<int>^>^ AppendKeyAsync(
	String^ key, 
	String^ value
)
```

```
member AppendKeyAsync : 
        key : string * 
        value : string -> Task<OperateResult<int>> 
```

#### 参数

key
:   类型：SystemString  
    关键字

value
:   类型：SystemString  
    数值

#### 返回值

类型：Task[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Int32  
追加 value 之后， key 中字符串的长度。

![](../icons/SectionExpanded.png)参见

#### 引用

[RedisClient 类](e33d9088-6df0-32e9-33f1-426e944c3b7f.htm)

[HslCommunication.Enthernet.Redis 命名空间](cc44ce00-3781-c690-b233-24000870bf42.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ChangePassword 方法 

[原文連結](http://api.hslcommunication.cn/html/78a37d4d-bdee-7a09-d0ce-dc245627b1d1.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Enthernet.Redis](../html/cc44ce00-3781-c690-b233-24000870bf42.htm "HslCommunication.Enthernet.Redis")

[RedisClient 类](../html/e33d9088-6df0-32e9-33f1-426e944c3b7f.htm "RedisClient 类")

[RedisClient 方法](../html/e41be4a5-1490-ffbd-f241-938f2d6c3f07.htm "RedisClient 方法")

[AppendKey 方法](../html/f2521b8c-0fb1-6b34-4adc-ca6804c7195f.htm "AppendKey 方法 ")

[AppendKeyAsync 方法](../html/c39b4c6a-b7f6-68fe-dfc9-7cce1b653e24.htm "AppendKeyAsync 方法 ")

[ChangePassword 方法](../html/78a37d4d-bdee-7a09-d0ce-dc245627b1d1.htm "ChangePassword 方法 ")

[ChangePasswordAsync 方法](../html/58f6f736-093d-6e4d-3ad3-19590fa2a3a3.htm "ChangePasswordAsync 方法 ")

[DBSize 方法](../html/5814129a-9a73-25ac-e726-fe1da141eb28.htm "DBSize 方法 ")

[DBSizeAsync 方法](../html/243a1f5f-a474-2e20-1298-06eff125f8b5.htm "DBSizeAsync 方法 ")

[DecrementKey 方法](../html/e9fa6c94-43ca-8291-79ca-6e04ca9289e7.htm "DecrementKey 方法 ")

[DecrementKeyAsync 方法](../html/6837e351-b8ef-a90e-c5e2-73b31621f837.htm "DecrementKeyAsync 方法 ")

[DeleteHashKey 方法](../html/19e3ebc3-e6fc-71cf-e80d-cb73e3011170.htm "DeleteHashKey 方法 ")

[DeleteHashKeyAsync 方法](../html/0dcc0d8a-1ed1-5c81-fb0e-e1c93be9d579.htm "DeleteHashKeyAsync 方法 ")

[DeleteKey 方法](../html/e2f0f305-d8bc-d515-3a62-cc91a3c1de76.htm "DeleteKey 方法 ")

[DeleteKeyAsync 方法](../html/0d8b5b0c-5979-9140-4ccf-992f4eaf1553.htm "DeleteKeyAsync 方法 ")

[ExistsHashKey 方法](../html/5b0bbeb9-4ab9-1f83-9e09-5f7227ac28e6.htm "ExistsHashKey 方法 ")

[ExistsHashKeyAsync 方法](../html/a80d408e-d621-a502-7829-6c7ea2a081f9.htm "ExistsHashKeyAsync 方法 ")

[ExistsKey 方法](../html/2b1b5b08-13ed-88c3-445d-d3613b48bc95.htm "ExistsKey 方法 ")

[ExistsKeyAsync 方法](../html/0d4357ed-5de3-1444-fce2-0e9be619d20f.htm "ExistsKeyAsync 方法 ")

[ExpireKey 方法](../html/0272044f-cc55-6b6a-42f0-6a0d89f37fe7.htm "ExpireKey 方法 ")

[ExpireKeyAsync 方法](../html/53c18b2d-2169-0519-602e-f65302f490e1.htm "ExpireKeyAsync 方法 ")

[ExtraOnDisconnect 方法](../html/49d9b263-02d8-e822-6a33-49f93737b98b.htm "ExtraOnDisconnect 方法 ")

[ExtraOnDisconnectAsync 方法](../html/f681350d-2aad-40e0-0824-30fbd3fe57ff.htm "ExtraOnDisconnectAsync 方法 ")

[FlushDB 方法](../html/38c28c3f-63bf-f1a0-67bc-6547799257f9.htm "FlushDB 方法 ")

[FlushDBAsync 方法](../html/40e3ccb2-9fe7-c0d7-88b9-95266a13c336.htm "FlushDBAsync 方法 ")

[GetListLength 方法](../html/36330604-2557-4fc1-bdd7-fc3aff05e997.htm "GetListLength 方法 ")

[GetListLengthAsync 方法](../html/0e09431d-45b3-b539-a618-a147ae377f88.htm "GetListLengthAsync 方法 ")

[IncrementHashKey 方法](../html/6404f402-9882-5453-7df7-275f61748021.htm "IncrementHashKey 方法 ")

[IncrementHashKeyAsync 方法](../html/7a69e1a1-8689-3e23-8aa3-cdad33e7e3f8.htm "IncrementHashKeyAsync 方法 ")

[IncrementKey 方法](../html/04d46420-bedd-9016-075c-3bb1a6623583.htm "IncrementKey 方法 ")

[IncrementKeyAsync 方法](../html/95f88e06-114c-b5b3-ee2f-2e9f97f27b8c.htm "IncrementKeyAsync 方法 ")

[InitializationOnConnect 方法](../html/72414fc8-b43a-3377-acdc-c59ae0628e11.htm "InitializationOnConnect 方法 ")

[InitializationOnConnectAsync 方法](../html/9156d090-6b1d-7381-c396-81491b4b6756.htm "InitializationOnConnectAsync 方法 ")

[ListInsertAfter 方法](../html/11b6eb28-b2ee-57a5-c724-a99df9101c71.htm "ListInsertAfter 方法 ")

[ListInsertAfterAsync 方法](../html/9defd0d6-5681-4c26-0973-14ecca4fc821.htm "ListInsertAfterAsync 方法 ")

[ListInsertBefore 方法](../html/b19eed6b-2778-4b91-6f6b-e75cd9f1088e.htm "ListInsertBefore 方法 ")

[ListInsertBeforeAsync 方法](../html/0bf5390c-af9f-ecd6-1148-34a20abee73d.htm "ListInsertBeforeAsync 方法 ")

[ListLeftPop 方法](../html/8865ed9c-69dc-2f28-7f53-42e1c03d152b.htm "ListLeftPop 方法 ")

[ListLeftPopAsync 方法](../html/f4d7f464-73f7-ff5d-0834-e18b5a91c994.htm "ListLeftPopAsync 方法 ")

[ListLeftPush 方法](../html/232a2dbd-e2c6-5fbf-af96-5269a10df4c9.htm "ListLeftPush 方法 ")

[ListLeftPushAsync 方法](../html/cb91133d-ce36-f4f7-a9d3-81dbda7d6617.htm "ListLeftPushAsync 方法 ")

[ListLeftPushX 方法](../html/33093159-9864-a0dd-6a09-8bf134942031.htm "ListLeftPushX 方法 ")

[ListLeftPushXAsync 方法](../html/6bd2ab31-ff58-71f0-b68f-f36d7f53d845.htm "ListLeftPushXAsync 方法 ")

[ListRange 方法](../html/efe06eb0-b945-fe65-5e81-21888cdda04f.htm "ListRange 方法 ")

[ListRangeAsync 方法](../html/32e344f9-cdac-fc0e-2c50-bf2dd54b4f98.htm "ListRangeAsync 方法 ")

[ListRemoveElementMatch 方法](../html/0ab40e89-d265-0b48-05de-5c1a7b41d326.htm "ListRemoveElementMatch 方法 ")

[ListRemoveElementMatchAsync 方法](../html/5610de76-3198-5b5a-8f3d-b8fc937d7c83.htm "ListRemoveElementMatchAsync 方法 ")

[ListRightPop 方法](../html/97855b86-6a07-76e2-b887-1e830be705de.htm "ListRightPop 方法 ")

[ListRightPopAsync 方法](../html/fde8cee2-f393-a463-d561-74568545edc7.htm "ListRightPopAsync 方法 ")

[ListRightPopLeftPush 方法](../html/e0b58cc4-4d80-64d6-729c-0dcdaafc4701.htm "ListRightPopLeftPush 方法 ")

[ListRightPopLeftPushAsync 方法](../html/41069df2-0dad-6f01-895f-a91086a92fa9.htm "ListRightPopLeftPushAsync 方法 ")

[ListRightPush 方法](../html/d71ca156-74d3-536e-6925-315d948c0ace.htm "ListRightPush 方法 ")

[ListRightPushAsync 方法](../html/53144654-da3d-b367-b39c-0dbc68275e59.htm "ListRightPushAsync 方法 ")

[ListRightPushX 方法](../html/e6ccfd1a-25ee-5eff-d56c-0ff9eb017fbc.htm "ListRightPushX 方法 ")

[ListRightPushXAsync 方法](../html/3252db2a-56cf-1aa2-7716-f9eb2596ec62.htm "ListRightPushXAsync 方法 ")

[ListSet 方法](../html/5f2e7b8b-b0e7-87e6-dac7-81b572cf3fa3.htm "ListSet 方法 ")

[ListSetAsync 方法](../html/ca019ec2-eb32-41f8-a460-89c70f90ace9.htm "ListSetAsync 方法 ")

[ListTrim 方法](../html/acf21364-1692-f9da-d337-c460b12a012d.htm "ListTrim 方法 ")

[ListTrimAsync 方法](../html/9e86492d-205c-8427-9052-50ea3d8c1c64.htm "ListTrimAsync 方法 ")

[MoveKey 方法](../html/5aa50043-85b9-a245-606a-de2b69164218.htm "MoveKey 方法 ")

[MoveKeyAsync 方法](../html/a93a74bf-65bd-733b-e39a-06d705252d99.htm "MoveKeyAsync 方法 ")

[OperateLongNumberFromServer 方法](../html/d2cdcd44-a6fb-7aa7-8c9c-08ff43107f65.htm "OperateLongNumberFromServer 方法 ")

[OperateLongNumberFromServerAsync 方法](../html/0005132d-e8c6-d51b-69dc-c9f8feaf8a9d.htm "OperateLongNumberFromServerAsync 方法 ")

[OperateNumberFromServer 方法](../html/69aecc05-6901-f53a-da37-ebf85bae7053.htm "OperateNumberFromServer 方法 ")

[OperateNumberFromServerAsync 方法](../html/b53bd6fa-d215-55f7-bbac-2ee171ed7498.htm "OperateNumberFromServerAsync 方法 ")

[OperateStatusFromServer 方法](../html/d460c3bb-e484-3fcb-e206-d20dc7201d63.htm "OperateStatusFromServer 方法 ")

[OperateStatusFromServerAsync 方法](../html/e1a241a3-e200-747f-bb8b-5b7578a08999.htm "OperateStatusFromServerAsync 方法 ")

[OperateStringFromServer 方法](../html/822914af-cf62-1ccc-d22e-2df7dfac7294.htm "OperateStringFromServer 方法 ")

[OperateStringFromServerAsync 方法](../html/eb31e52d-5387-eb14-291e-22f0e64666a4.htm "OperateStringFromServerAsync 方法 ")

[OperateStringsFromServer 方法](../html/c85c1a0e-83f9-69c9-a9b9-8a947ca54583.htm "OperateStringsFromServer 方法 ")

[OperateStringsFromServerAsync 方法](../html/2db50b51-32b8-1e92-cb9a-f2c847437e50.htm "OperateStringsFromServerAsync 方法 ")

[PersistKey 方法](../html/c68a840b-ce64-3303-9339-a10f88bcbb77.htm "PersistKey 方法 ")

[PersistKeyAsync 方法](../html/a5d7ee67-9041-6455-7fc1-c6e96a442799.htm "PersistKeyAsync 方法 ")

[Ping 方法](../html/9cdd4e99-90a2-0ae8-141c-26f0f4a4d698.htm "Ping 方法 ")

[PingAsync 方法](../html/d6c8e625-edb2-7ff2-a54d-e852d49c73c8.htm "PingAsync 方法 ")

[Publish 方法](../html/74f20db5-aaaa-af2e-b4e6-5d192196b222.htm "Publish 方法 ")

[PublishAsync 方法](../html/61428af8-94a2-b2b0-419b-5c0d3df4bee5.htm "PublishAsync 方法 ")

[Read(T) 方法](../html/08365a89-e2d7-2aa5-3fd5-ab1a5d60c440.htm "Read(T) 方法 ")

[ReadAllKeys 方法](../html/53c48696-f991-f734-12b0-1ad306a84c6b.htm "ReadAllKeys 方法 ")

[ReadAllKeysAsync 方法](../html/84c1be8c-3286-0410-0bd1-bb0185ff197e.htm "ReadAllKeysAsync 方法 ")

[ReadAndWriteKey 方法](../html/a6908f0e-1991-5210-f6c6-1e8b06369e56.htm "ReadAndWriteKey 方法 ")

[ReadAndWriteKeyAsync 方法](../html/6686726f-93ba-6435-76b0-b501c31e637e.htm "ReadAndWriteKeyAsync 方法 ")

[ReadAsync(T) 方法](../html/43a59914-1502-1353-0bf3-20281ac92ca9.htm "ReadAsync(T) 方法 ")

[ReadCustomer 方法](../html/a4b7f43e-34ba-8aac-a705-a9bc8d4a5298.htm "ReadCustomer 方法 ")

[ReadCustomerAsync 方法](../html/bafea5ed-4293-ae60-c8e7-bb3d3481db44.htm "ReadCustomerAsync 方法 ")

[ReadFromCoreServer 方法](../html/a5985ec7-e909-5a40-b2b9-7160363a1156.htm "ReadFromCoreServer 方法 ")

[ReadFromCoreServerAsync 方法](../html/984a89de-20da-e7ed-e471-3e889955aa31.htm "ReadFromCoreServerAsync 方法 ")

[ReadHashKey 方法](../html/6e8fc531-b2d6-2b1e-4885-1b937745652f.htm "ReadHashKey 方法 ")

[ReadHashKeyAll 方法](../html/63828fbe-fbdf-772a-010b-1425b0584572.htm "ReadHashKeyAll 方法 ")

[ReadHashKeyAllAsync 方法](../html/31c54864-f7bc-b3ce-7e77-15f0c1dcc540.htm "ReadHashKeyAllAsync 方法 ")

[ReadHashKeyAsync 方法](../html/36c995b4-42a7-e264-8a03-3390e24376ec.htm "ReadHashKeyAsync 方法 ")

[ReadHashKeyLength 方法](../html/32d483da-49e0-a205-bdc4-a659711be1c0.htm "ReadHashKeyLength 方法 ")

[ReadHashKeyLengthAsync 方法](../html/5d0d4570-1f16-4277-cff7-a40607a530ab.htm "ReadHashKeyLengthAsync 方法 ")

[ReadHashKeys 方法](../html/420286f8-5876-ee1b-b328-0b6eba7cefac.htm "ReadHashKeys 方法 ")

[ReadHashKeysAsync 方法](../html/bbc8310c-029d-1ea5-b0e3-239b29ee8743.htm "ReadHashKeysAsync 方法 ")

[ReadHashValues 方法](../html/1d8b8e06-3d34-f4c6-68d0-20c14b309de1.htm "ReadHashValues 方法 ")

[ReadHashValuesAsync 方法](../html/403e7344-c665-def6-70d9-81e5b22a740b.htm "ReadHashValuesAsync 方法 ")

[ReadKey 方法](../html/93fa3dab-d22c-d028-00c3-e81391085e6f.htm "ReadKey 方法 ")

[ReadKeyAsync 方法](../html/48cee995-f908-c3a2-ab70-829d789f58a5.htm "ReadKeyAsync 方法 ")

[ReadKeyLength 方法](../html/7746357b-7e1b-d150-d019-722b752f5eeb.htm "ReadKeyLength 方法 ")

[ReadKeyLengthAsync 方法](../html/7dd1b305-77b1-920c-14c0-2fe38e970791.htm "ReadKeyLengthAsync 方法 ")

[ReadKeyRange 方法](../html/2e3e9582-cb24-44e6-d78b-7c672acdd4c8.htm "ReadKeyRange 方法 ")

[ReadKeyRangeAsync 方法](../html/407fdd5b-8b57-75d2-63be-93294311cb30.htm "ReadKeyRangeAsync 方法 ")

[ReadKeyTTL 方法](../html/679079c3-36e5-1618-3409-18459a74a0c4.htm "ReadKeyTTL 方法 ")

[ReadKeyTTLAsync 方法](../html/24887926-0fd5-b65f-d070-1d439a5273ef.htm "ReadKeyTTLAsync 方法 ")

[ReadKeyType 方法](../html/1f702770-7804-880a-038e-c002ed31c42f.htm "ReadKeyType 方法 ")

[ReadKeyTypeAsync 方法](../html/91c104ff-a646-76de-9bed-5e64237a57fc.htm "ReadKeyTypeAsync 方法 ")

[ReadListByIndex 方法](../html/43c7dae8-4c01-d2c9-b38b-0e7bef68b581.htm "ReadListByIndex 方法 ")

[ReadListByIndexAsync 方法](../html/1e1163ad-8161-9390-cff2-7919c8145c92.htm "ReadListByIndexAsync 方法 ")

[ReadRandomKey 方法](../html/0d358a8f-b9e0-c659-c218-f051aaedcbfe.htm "ReadRandomKey 方法 ")

[ReadRandomKeyAsync 方法](../html/b97f91d2-2ab6-a39b-86f7-71238cd44fac.htm "ReadRandomKeyAsync 方法 ")

[ReadServerTime 方法](../html/c6c5f1fa-4639-0a32-8a42-48d56985e8e6.htm "ReadServerTime 方法 ")

[ReadServerTimeAsync 方法](../html/53022c3b-766e-c31b-1b55-548ceb3b2f42.htm "ReadServerTimeAsync 方法 ")

[RenameKey 方法](../html/ecf533af-233c-16b8-f88d-cb232e115013.htm "RenameKey 方法 ")

[RenameKeyAsync 方法](../html/20073ebe-0c16-2f57-a96b-79246e01129d.htm "RenameKeyAsync 方法 ")

[Save 方法](../html/261a8daa-9890-95bc-f8e4-02b78eb86277.htm "Save 方法 ")

[SaveAsync 方法](../html/88c2043f-e9e6-9d95-4a2e-455452cf8adf.htm "SaveAsync 方法 ")

[SelectDB 方法](../html/584a38b8-96e0-736f-7a24-981a8e920574.htm "SelectDB 方法 ")

[SelectDBAsync 方法](../html/462ed998-9938-a261-bde4-579b4b9d3dc1.htm "SelectDBAsync 方法 ")

[SetAdd 方法](../html/6f52901b-4d7f-70ca-1bf1-2e7146604dba.htm "SetAdd 方法 ")

[SetAddAsync 方法](../html/a59f7c04-0b20-35a4-a431-53299cb829f0.htm "SetAddAsync 方法 ")

[SetCard 方法](../html/cc0a3d1d-019b-f0ba-90a5-1bab767d6f05.htm "SetCard 方法 ")

[SetCardAsync 方法](../html/0d0665a9-a174-d154-f95b-2bf10d885bc3.htm "SetCardAsync 方法 ")

[SetDiff 方法](../html/bcc43604-dc32-a27f-53a4-a89ab29916be.htm "SetDiff 方法 ")

[SetDiffAsync 方法](../html/b70eae65-3537-3ab7-4cf9-ae6c314884c9.htm "SetDiffAsync 方法 ")

[SetDiffStore 方法](../html/7c356901-c174-c81b-50b9-ce8d60258d1e.htm "SetDiffStore 方法 ")

[SetDiffStoreAsync 方法](../html/697ac488-c2e1-8c5a-5e5d-845f676c85fd.htm "SetDiffStoreAsync 方法 ")

[SetInter 方法](../html/fa3a123d-447d-ce8b-5670-e47956d427d1.htm "SetInter 方法 ")

[SetInterAsync 方法](../html/622854f7-5e07-092e-42eb-06cc63580467.htm "SetInterAsync 方法 ")

[SetInterStore 方法](../html/996ff5a7-0d0b-5b17-546b-eb006fa5773f.htm "SetInterStore 方法 ")

[SetInterStoreAsync 方法](../html/115a5c5d-2c7b-4c9a-daec-92477fb3044a.htm "SetInterStoreAsync 方法 ")

[SetIsMember 方法](../html/1c6d465b-353a-34f2-8f46-9d4037e58911.htm "SetIsMember 方法 ")

[SetIsMemberAsync 方法](../html/21419054-3273-2aab-d486-633ab8861039.htm "SetIsMemberAsync 方法 ")

[SetMembers 方法](../html/7a5b277f-1fe1-141f-346b-3722d97ac7ca.htm "SetMembers 方法 ")

[SetMembersAsync 方法](../html/ce86155f-4be2-3836-3594-1398d78183ad.htm "SetMembersAsync 方法 ")

[SetMove 方法](../html/8b76c71a-c755-708b-e1e7-61eda8b9456c.htm "SetMove 方法 ")

[SetMoveAsync 方法](../html/dc01893c-b09e-e134-49d3-76d0dbc308b9.htm "SetMoveAsync 方法 ")

[SetPop 方法](../html/8375117d-d7b5-1f33-cb59-80cc632b736d.htm "SetPop 方法 ")

[SetPopAsync 方法](../html/66039533-6b1e-ad66-e89b-ef8f7d46c47b.htm "SetPopAsync 方法 ")

[SetRandomMember 方法](../html/2769a6d4-3b4c-75b2-c57b-abe19bb2f8e9.htm "SetRandomMember 方法 ")

[SetRandomMemberAsync 方法](../html/fb4fda18-00e6-2545-d037-6585a9487334.htm "SetRandomMemberAsync 方法 ")

[SetRemove 方法](../html/a0ba101e-935e-aca6-da5e-bb6187b42f02.htm "SetRemove 方法 ")

[SetRemoveAsync 方法](../html/f80e5db7-218b-c4bf-0c57-ec2555f6ec4e.htm "SetRemoveAsync 方法 ")

[SetUnion 方法](../html/c255e92b-1a91-f3e4-49cd-79f39d7a1079.htm "SetUnion 方法 ")

[SetUnionAsync 方法](../html/184ffefd-f6f5-b5d8-e1d6-45ae6d2d653f.htm "SetUnionAsync 方法 ")

[SetUnionStore 方法](../html/56394825-4bd1-6422-4f9b-1520ed05f927.htm "SetUnionStore 方法 ")

[SetUnionStoreAsync 方法](../html/99689f9d-ebcc-cc16-19e3-4c46f3016084.htm "SetUnionStoreAsync 方法 ")

[SubscribeMessage 方法](../html/bbaad51e-d35a-fb1f-4bb5-c001d47693be.htm "SubscribeMessage 方法 ")

[ToString 方法](../html/6ff40e11-12f4-8759-45a0-594ca3e4fcdd.htm "ToString 方法 ")

[UnSubscribeMessage 方法](../html/592327f8-7e50-4d48-c229-ca17a039b6a3.htm "UnSubscribeMessage 方法 ")

[Write(T) 方法](../html/70a08f37-0f86-76f1-55a5-df1ddd5a7d16.htm "Write(T) 方法 ")

[WriteAndPublishKey 方法](../html/69029e5e-a78a-5c83-0300-b6dd057f9308.htm "WriteAndPublishKey 方法 ")

[WriteAndPublishKeyAsync 方法](../html/d44abf4b-02bf-cfa5-cbc8-0b36edab4521.htm "WriteAndPublishKeyAsync 方法 ")

[WriteAsync(T) 方法](../html/9d54e352-2bd0-2655-3e14-49a0213e6040.htm "WriteAsync(T) 方法 ")

[WriteExpireKey 方法](../html/816ca8c0-5e3a-4609-acc3-6f227ac2dfde.htm "WriteExpireKey 方法 ")

[WriteExpireKeyAsync 方法](../html/f1dce55a-5681-32f6-2c7c-cc664465490c.htm "WriteExpireKeyAsync 方法 ")

[WriteHashKey 方法](../html/a7c80b1f-ec11-c8f4-dae2-0fa0acada685.htm "WriteHashKey 方法 ")

[WriteHashKeyAsync 方法](../html/20122681-134d-f220-683c-11c6c9201ac9.htm "WriteHashKeyAsync 方法 ")

[WriteHashKeyNx 方法](../html/819ddbb0-74d4-4667-6e55-372fac166c68.htm "WriteHashKeyNx 方法 ")

[WriteHashKeyNxAsync 方法](../html/b4571eb9-be29-fbb8-691b-f3a93d9245ec.htm "WriteHashKeyNxAsync 方法 ")

[WriteKey 方法](../html/12919427-9ab8-d765-27cc-0c3ee2c481d2.htm "WriteKey 方法 ")

[WriteKeyAsync 方法](../html/b136bcbb-3813-f660-54f1-8c5ccefb5ddf.htm "WriteKeyAsync 方法 ")

[WriteKeyIfNotExists 方法](../html/9d48312b-2ce1-077f-bf47-3e6b5a328c45.htm "WriteKeyIfNotExists 方法 ")

[WriteKeyIfNotExistsAsync 方法](../html/e6844e26-ba17-a0a8-0ff2-bb1a1e8bce8e.htm "WriteKeyIfNotExistsAsync 方法 ")

[WriteKeyRange 方法](../html/c1f98cc9-f791-4a35-bf67-ea29e8930d32.htm "WriteKeyRange 方法 ")

[WriteKeyRangeAsync 方法](../html/cdfaffec-518b-2754-5a10-3c5b68ee1767.htm "WriteKeyRangeAsync 方法 ")

[ZSetAdd 方法](../html/9c15cbbe-df53-5280-c9e4-b9d020222e0c.htm "ZSetAdd 方法 ")

[ZSetAddAsync 方法](../html/5318d2e6-d635-ed62-18e5-6573a530ab3e.htm "ZSetAddAsync 方法 ")

[ZSetCard 方法](../html/aa928601-c5f6-63fb-6bd8-e2e383796043.htm "ZSetCard 方法 ")

[ZSetCardAsync 方法](../html/b760ff15-4c7e-510f-1927-e17e6310e142.htm "ZSetCardAsync 方法 ")

[ZSetCount 方法](../html/9872ddc0-8b8d-d9d6-4d08-64001fdefd29.htm "ZSetCount 方法 ")

[ZSetCountAsync 方法](../html/b294e1ee-27c7-7326-eb4b-8669040aecf0.htm "ZSetCountAsync 方法 ")

[ZSetIncreaseBy 方法](../html/c84e0c14-15cd-451a-dec4-2a112a73772e.htm "ZSetIncreaseBy 方法 ")

[ZSetIncreaseByAsync 方法](../html/9805ac6b-aa05-2152-7dbb-7f8077ea462b.htm "ZSetIncreaseByAsync 方法 ")

[ZSetRange 方法](../html/759e2523-e2ad-0dc1-a904-354031a8ae50.htm "ZSetRange 方法 ")

[ZSetRangeAsync 方法](../html/7724eb62-2653-65a7-779d-391d630c2faf.htm "ZSetRangeAsync 方法 ")

[ZSetRangeByScore 方法](../html/acab27d0-6b6c-1f25-00ad-fafd4ac068fa.htm "ZSetRangeByScore 方法 ")

[ZSetRangeByScoreAsync 方法](../html/d1e0aa3c-d5f9-caeb-4017-699f956e661d.htm "ZSetRangeByScoreAsync 方法 ")

[ZSetRank 方法](../html/7957411f-32c1-4e07-a4f2-9d1948c63a2b.htm "ZSetRank 方法 ")

[ZSetRankAsync 方法](../html/1181fd8c-02aa-9d53-cbeb-35d542bed241.htm "ZSetRankAsync 方法 ")

[ZSetRemove 方法](../html/a76625ea-7a7b-161d-a039-02f4f5a8e997.htm "ZSetRemove 方法 ")

[ZSetRemoveAsync 方法](../html/2f1b3de6-ad74-e2b7-3350-79d9afe13c9d.htm "ZSetRemoveAsync 方法 ")

[ZSetRemoveRangeByRank 方法](../html/cb81dfc2-6e30-0aca-8bcb-348cdc32bc35.htm "ZSetRemoveRangeByRank 方法 ")

[ZSetRemoveRangeByRankAsync 方法](../html/9bde4b5f-b9bd-e01b-feb4-82c9b532dee4.htm "ZSetRemoveRangeByRankAsync 方法 ")

[ZSetRemoveRangeByScore 方法](../html/ac7be899-fe4e-3f12-d60f-81a038a927c7.htm "ZSetRemoveRangeByScore 方法 ")

[ZSetRemoveRangeByScoreAsync 方法](../html/6ccee6bf-8635-14bf-2e7f-d1ac351c70c9.htm "ZSetRemoveRangeByScoreAsync 方法 ")

[ZSetReverseRange 方法](../html/ef6defc7-62a2-c8a4-4025-32f6b3a82ae4.htm "ZSetReverseRange 方法 ")

[ZSetReverseRangeAsync 方法](../html/b40d37b7-1118-6aec-c7da-85ec54b457d0.htm "ZSetReverseRangeAsync 方法 ")

[ZSetReverseRangeByScore 方法](../html/3c1a26a8-f8e5-a835-7051-94b4d9c851ab.htm "ZSetReverseRangeByScore 方法 ")

[ZSetReverseRangeByScoreAsync 方法](../html/3a3cfbff-cee7-e1f3-d9c1-e7528db4e998.htm "ZSetReverseRangeByScoreAsync 方法 ")

[ZSetReverseRank 方法](../html/d524195b-6b02-a8e8-330a-3454a6ed21cd.htm "ZSetReverseRank 方法 ")

[ZSetReverseRankAsync 方法](../html/c010781c-2851-7158-f953-1779929e3242.htm "ZSetReverseRankAsync 方法 ")

[ZSetScore 方法](../html/a98ee7e4-f361-f403-975b-a197e0ca0268.htm "ZSetScore 方法 ")

[ZSetScoreAsync 方法](../html/fd24c433-944e-2ecb-4bcd-266bfc6825d4.htm "ZSetScoreAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| RedisClientChangePassword 方法 |

修改Redis的密码信息，如果不需要密码，则传入空字符串即可

**命名空间：**
 [HslCommunication.Enthernet.Redis](cc44ce00-3781-c690-b233-24000870bf42.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public OperateResult ChangePassword(
	string password
)
```

```
Public Function ChangePassword ( 
	password As String
) As OperateResult
```

```
public:
OperateResult^ ChangePassword(
	String^ password
)
```

```
member ChangePassword : 
        password : string -> OperateResult 
```

#### 参数

password
:   类型：SystemString  
    密码信息

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否更新了密码信息

![](../icons/SectionExpanded.png)参见

#### 引用

[RedisClient 类](e33d9088-6df0-32e9-33f1-426e944c3b7f.htm)

[HslCommunication.Enthernet.Redis 命名空间](cc44ce00-3781-c690-b233-24000870bf42.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ChangePasswordAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/58f6f736-093d-6e4d-3ad3-19590fa2a3a3.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Enthernet.Redis](../html/cc44ce00-3781-c690-b233-24000870bf42.htm "HslCommunication.Enthernet.Redis")

[RedisClient 类](../html/e33d9088-6df0-32e9-33f1-426e944c3b7f.htm "RedisClient 类")

[RedisClient 方法](../html/e41be4a5-1490-ffbd-f241-938f2d6c3f07.htm "RedisClient 方法")

[AppendKey 方法](../html/f2521b8c-0fb1-6b34-4adc-ca6804c7195f.htm "AppendKey 方法 ")

[AppendKeyAsync 方法](../html/c39b4c6a-b7f6-68fe-dfc9-7cce1b653e24.htm "AppendKeyAsync 方法 ")

[ChangePassword 方法](../html/78a37d4d-bdee-7a09-d0ce-dc245627b1d1.htm "ChangePassword 方法 ")

[ChangePasswordAsync 方法](../html/58f6f736-093d-6e4d-3ad3-19590fa2a3a3.htm "ChangePasswordAsync 方法 ")

[DBSize 方法](../html/5814129a-9a73-25ac-e726-fe1da141eb28.htm "DBSize 方法 ")

[DBSizeAsync 方法](../html/243a1f5f-a474-2e20-1298-06eff125f8b5.htm "DBSizeAsync 方法 ")

[DecrementKey 方法](../html/e9fa6c94-43ca-8291-79ca-6e04ca9289e7.htm "DecrementKey 方法 ")

[DecrementKeyAsync 方法](../html/6837e351-b8ef-a90e-c5e2-73b31621f837.htm "DecrementKeyAsync 方法 ")

[DeleteHashKey 方法](../html/19e3ebc3-e6fc-71cf-e80d-cb73e3011170.htm "DeleteHashKey 方法 ")

[DeleteHashKeyAsync 方法](../html/0dcc0d8a-1ed1-5c81-fb0e-e1c93be9d579.htm "DeleteHashKeyAsync 方法 ")

[DeleteKey 方法](../html/e2f0f305-d8bc-d515-3a62-cc91a3c1de76.htm "DeleteKey 方法 ")

[DeleteKeyAsync 方法](../html/0d8b5b0c-5979-9140-4ccf-992f4eaf1553.htm "DeleteKeyAsync 方法 ")

[ExistsHashKey 方法](../html/5b0bbeb9-4ab9-1f83-9e09-5f7227ac28e6.htm "ExistsHashKey 方法 ")

[ExistsHashKeyAsync 方法](../html/a80d408e-d621-a502-7829-6c7ea2a081f9.htm "ExistsHashKeyAsync 方法 ")

[ExistsKey 方法](../html/2b1b5b08-13ed-88c3-445d-d3613b48bc95.htm "ExistsKey 方法 ")

[ExistsKeyAsync 方法](../html/0d4357ed-5de3-1444-fce2-0e9be619d20f.htm "ExistsKeyAsync 方法 ")

[ExpireKey 方法](../html/0272044f-cc55-6b6a-42f0-6a0d89f37fe7.htm "ExpireKey 方法 ")

[ExpireKeyAsync 方法](../html/53c18b2d-2169-0519-602e-f65302f490e1.htm "ExpireKeyAsync 方法 ")

[ExtraOnDisconnect 方法](../html/49d9b263-02d8-e822-6a33-49f93737b98b.htm "ExtraOnDisconnect 方法 ")

[ExtraOnDisconnectAsync 方法](../html/f681350d-2aad-40e0-0824-30fbd3fe57ff.htm "ExtraOnDisconnectAsync 方法 ")

[FlushDB 方法](../html/38c28c3f-63bf-f1a0-67bc-6547799257f9.htm "FlushDB 方法 ")

[FlushDBAsync 方法](../html/40e3ccb2-9fe7-c0d7-88b9-95266a13c336.htm "FlushDBAsync 方法 ")

[GetListLength 方法](../html/36330604-2557-4fc1-bdd7-fc3aff05e997.htm "GetListLength 方法 ")

[GetListLengthAsync 方法](../html/0e09431d-45b3-b539-a618-a147ae377f88.htm "GetListLengthAsync 方法 ")

[IncrementHashKey 方法](../html/6404f402-9882-5453-7df7-275f61748021.htm "IncrementHashKey 方法 ")

[IncrementHashKeyAsync 方法](../html/7a69e1a1-8689-3e23-8aa3-cdad33e7e3f8.htm "IncrementHashKeyAsync 方法 ")

[IncrementKey 方法](../html/04d46420-bedd-9016-075c-3bb1a6623583.htm "IncrementKey 方法 ")

[IncrementKeyAsync 方法](../html/95f88e06-114c-b5b3-ee2f-2e9f97f27b8c.htm "IncrementKeyAsync 方法 ")

[InitializationOnConnect 方法](../html/72414fc8-b43a-3377-acdc-c59ae0628e11.htm "InitializationOnConnect 方法 ")

[InitializationOnConnectAsync 方法](../html/9156d090-6b1d-7381-c396-81491b4b6756.htm "InitializationOnConnectAsync 方法 ")

[ListInsertAfter 方法](../html/11b6eb28-b2ee-57a5-c724-a99df9101c71.htm "ListInsertAfter 方法 ")

[ListInsertAfterAsync 方法](../html/9defd0d6-5681-4c26-0973-14ecca4fc821.htm "ListInsertAfterAsync 方法 ")

[ListInsertBefore 方法](../html/b19eed6b-2778-4b91-6f6b-e75cd9f1088e.htm "ListInsertBefore 方法 ")

[ListInsertBeforeAsync 方法](../html/0bf5390c-af9f-ecd6-1148-34a20abee73d.htm "ListInsertBeforeAsync 方法 ")

[ListLeftPop 方法](../html/8865ed9c-69dc-2f28-7f53-42e1c03d152b.htm "ListLeftPop 方法 ")

[ListLeftPopAsync 方法](../html/f4d7f464-73f7-ff5d-0834-e18b5a91c994.htm "ListLeftPopAsync 方法 ")

[ListLeftPush 方法](../html/232a2dbd-e2c6-5fbf-af96-5269a10df4c9.htm "ListLeftPush 方法 ")

[ListLeftPushAsync 方法](../html/cb91133d-ce36-f4f7-a9d3-81dbda7d6617.htm "ListLeftPushAsync 方法 ")

[ListLeftPushX 方法](../html/33093159-9864-a0dd-6a09-8bf134942031.htm "ListLeftPushX 方法 ")

[ListLeftPushXAsync 方法](../html/6bd2ab31-ff58-71f0-b68f-f36d7f53d845.htm "ListLeftPushXAsync 方法 ")

[ListRange 方法](../html/efe06eb0-b945-fe65-5e81-21888cdda04f.htm "ListRange 方法 ")

[ListRangeAsync 方法](../html/32e344f9-cdac-fc0e-2c50-bf2dd54b4f98.htm "ListRangeAsync 方法 ")

[ListRemoveElementMatch 方法](../html/0ab40e89-d265-0b48-05de-5c1a7b41d326.htm "ListRemoveElementMatch 方法 ")

[ListRemoveElementMatchAsync 方法](../html/5610de76-3198-5b5a-8f3d-b8fc937d7c83.htm "ListRemoveElementMatchAsync 方法 ")

[ListRightPop 方法](../html/97855b86-6a07-76e2-b887-1e830be705de.htm "ListRightPop 方法 ")

[ListRightPopAsync 方法](../html/fde8cee2-f393-a463-d561-74568545edc7.htm "ListRightPopAsync 方法 ")

[ListRightPopLeftPush 方法](../html/e0b58cc4-4d80-64d6-729c-0dcdaafc4701.htm "ListRightPopLeftPush 方法 ")

[ListRightPopLeftPushAsync 方法](../html/41069df2-0dad-6f01-895f-a91086a92fa9.htm "ListRightPopLeftPushAsync 方法 ")

[ListRightPush 方法](../html/d71ca156-74d3-536e-6925-315d948c0ace.htm "ListRightPush 方法 ")

[ListRightPushAsync 方法](../html/53144654-da3d-b367-b39c-0dbc68275e59.htm "ListRightPushAsync 方法 ")

[ListRightPushX 方法](../html/e6ccfd1a-25ee-5eff-d56c-0ff9eb017fbc.htm "ListRightPushX 方法 ")

[ListRightPushXAsync 方法](../html/3252db2a-56cf-1aa2-7716-f9eb2596ec62.htm "ListRightPushXAsync 方法 ")

[ListSet 方法](../html/5f2e7b8b-b0e7-87e6-dac7-81b572cf3fa3.htm "ListSet 方法 ")

[ListSetAsync 方法](../html/ca019ec2-eb32-41f8-a460-89c70f90ace9.htm "ListSetAsync 方法 ")

[ListTrim 方法](../html/acf21364-1692-f9da-d337-c460b12a012d.htm "ListTrim 方法 ")

[ListTrimAsync 方法](../html/9e86492d-205c-8427-9052-50ea3d8c1c64.htm "ListTrimAsync 方法 ")

[MoveKey 方法](../html/5aa50043-85b9-a245-606a-de2b69164218.htm "MoveKey 方法 ")

[MoveKeyAsync 方法](../html/a93a74bf-65bd-733b-e39a-06d705252d99.htm "MoveKeyAsync 方法 ")

[OperateLongNumberFromServer 方法](../html/d2cdcd44-a6fb-7aa7-8c9c-08ff43107f65.htm "OperateLongNumberFromServer 方法 ")

[OperateLongNumberFromServerAsync 方法](../html/0005132d-e8c6-d51b-69dc-c9f8feaf8a9d.htm "OperateLongNumberFromServerAsync 方法 ")

[OperateNumberFromServer 方法](../html/69aecc05-6901-f53a-da37-ebf85bae7053.htm "OperateNumberFromServer 方法 ")

[OperateNumberFromServerAsync 方法](../html/b53bd6fa-d215-55f7-bbac-2ee171ed7498.htm "OperateNumberFromServerAsync 方法 ")

[OperateStatusFromServer 方法](../html/d460c3bb-e484-3fcb-e206-d20dc7201d63.htm "OperateStatusFromServer 方法 ")

[OperateStatusFromServerAsync 方法](../html/e1a241a3-e200-747f-bb8b-5b7578a08999.htm "OperateStatusFromServerAsync 方法 ")

[OperateStringFromServer 方法](../html/822914af-cf62-1ccc-d22e-2df7dfac7294.htm "OperateStringFromServer 方法 ")

[OperateStringFromServerAsync 方法](../html/eb31e52d-5387-eb14-291e-22f0e64666a4.htm "OperateStringFromServerAsync 方法 ")

[OperateStringsFromServer 方法](../html/c85c1a0e-83f9-69c9-a9b9-8a947ca54583.htm "OperateStringsFromServer 方法 ")

[OperateStringsFromServerAsync 方法](../html/2db50b51-32b8-1e92-cb9a-f2c847437e50.htm "OperateStringsFromServerAsync 方法 ")

[PersistKey 方法](../html/c68a840b-ce64-3303-9339-a10f88bcbb77.htm "PersistKey 方法 ")

[PersistKeyAsync 方法](../html/a5d7ee67-9041-6455-7fc1-c6e96a442799.htm "PersistKeyAsync 方法 ")

[Ping 方法](../html/9cdd4e99-90a2-0ae8-141c-26f0f4a4d698.htm "Ping 方法 ")

[PingAsync 方法](../html/d6c8e625-edb2-7ff2-a54d-e852d49c73c8.htm "PingAsync 方法 ")

[Publish 方法](../html/74f20db5-aaaa-af2e-b4e6-5d192196b222.htm "Publish 方法 ")

[PublishAsync 方法](../html/61428af8-94a2-b2b0-419b-5c0d3df4bee5.htm "PublishAsync 方法 ")

[Read(T) 方法](../html/08365a89-e2d7-2aa5-3fd5-ab1a5d60c440.htm "Read(T) 方法 ")

[ReadAllKeys 方法](../html/53c48696-f991-f734-12b0-1ad306a84c6b.htm "ReadAllKeys 方法 ")

[ReadAllKeysAsync 方法](../html/84c1be8c-3286-0410-0bd1-bb0185ff197e.htm "ReadAllKeysAsync 方法 ")

[ReadAndWriteKey 方法](../html/a6908f0e-1991-5210-f6c6-1e8b06369e56.htm "ReadAndWriteKey 方法 ")

[ReadAndWriteKeyAsync 方法](../html/6686726f-93ba-6435-76b0-b501c31e637e.htm "ReadAndWriteKeyAsync 方法 ")

[ReadAsync(T) 方法](../html/43a59914-1502-1353-0bf3-20281ac92ca9.htm "ReadAsync(T) 方法 ")

[ReadCustomer 方法](../html/a4b7f43e-34ba-8aac-a705-a9bc8d4a5298.htm "ReadCustomer 方法 ")

[ReadCustomerAsync 方法](../html/bafea5ed-4293-ae60-c8e7-bb3d3481db44.htm "ReadCustomerAsync 方法 ")

[ReadFromCoreServer 方法](../html/a5985ec7-e909-5a40-b2b9-7160363a1156.htm "ReadFromCoreServer 方法 ")

[ReadFromCoreServerAsync 方法](../html/984a89de-20da-e7ed-e471-3e889955aa31.htm "ReadFromCoreServerAsync 方法 ")

[ReadHashKey 方法](../html/6e8fc531-b2d6-2b1e-4885-1b937745652f.htm "ReadHashKey 方法 ")

[ReadHashKeyAll 方法](../html/63828fbe-fbdf-772a-010b-1425b0584572.htm "ReadHashKeyAll 方法 ")

[ReadHashKeyAllAsync 方法](../html/31c54864-f7bc-b3ce-7e77-15f0c1dcc540.htm "ReadHashKeyAllAsync 方法 ")

[ReadHashKeyAsync 方法](../html/36c995b4-42a7-e264-8a03-3390e24376ec.htm "ReadHashKeyAsync 方法 ")

[ReadHashKeyLength 方法](../html/32d483da-49e0-a205-bdc4-a659711be1c0.htm "ReadHashKeyLength 方法 ")

[ReadHashKeyLengthAsync 方法](../html/5d0d4570-1f16-4277-cff7-a40607a530ab.htm "ReadHashKeyLengthAsync 方法 ")

[ReadHashKeys 方法](../html/420286f8-5876-ee1b-b328-0b6eba7cefac.htm "ReadHashKeys 方法 ")

[ReadHashKeysAsync 方法](../html/bbc8310c-029d-1ea5-b0e3-239b29ee8743.htm "ReadHashKeysAsync 方法 ")

[ReadHashValues 方法](../html/1d8b8e06-3d34-f4c6-68d0-20c14b309de1.htm "ReadHashValues 方法 ")

[ReadHashValuesAsync 方法](../html/403e7344-c665-def6-70d9-81e5b22a740b.htm "ReadHashValuesAsync 方法 ")

[ReadKey 方法](../html/93fa3dab-d22c-d028-00c3-e81391085e6f.htm "ReadKey 方法 ")

[ReadKeyAsync 方法](../html/48cee995-f908-c3a2-ab70-829d789f58a5.htm "ReadKeyAsync 方法 ")

[ReadKeyLength 方法](../html/7746357b-7e1b-d150-d019-722b752f5eeb.htm "ReadKeyLength 方法 ")

[ReadKeyLengthAsync 方法](../html/7dd1b305-77b1-920c-14c0-2fe38e970791.htm "ReadKeyLengthAsync 方法 ")

[ReadKeyRange 方法](../html/2e3e9582-cb24-44e6-d78b-7c672acdd4c8.htm "ReadKeyRange 方法 ")

[ReadKeyRangeAsync 方法](../html/407fdd5b-8b57-75d2-63be-93294311cb30.htm "ReadKeyRangeAsync 方法 ")

[ReadKeyTTL 方法](../html/679079c3-36e5-1618-3409-18459a74a0c4.htm "ReadKeyTTL 方法 ")

[ReadKeyTTLAsync 方法](../html/24887926-0fd5-b65f-d070-1d439a5273ef.htm "ReadKeyTTLAsync 方法 ")

[ReadKeyType 方法](../html/1f702770-7804-880a-038e-c002ed31c42f.htm "ReadKeyType 方法 ")

[ReadKeyTypeAsync 方法](../html/91c104ff-a646-76de-9bed-5e64237a57fc.htm "ReadKeyTypeAsync 方法 ")

[ReadListByIndex 方法](../html/43c7dae8-4c01-d2c9-b38b-0e7bef68b581.htm "ReadListByIndex 方法 ")

[ReadListByIndexAsync 方法](../html/1e1163ad-8161-9390-cff2-7919c8145c92.htm "ReadListByIndexAsync 方法 ")

[ReadRandomKey 方法](../html/0d358a8f-b9e0-c659-c218-f051aaedcbfe.htm "ReadRandomKey 方法 ")

[ReadRandomKeyAsync 方法](../html/b97f91d2-2ab6-a39b-86f7-71238cd44fac.htm "ReadRandomKeyAsync 方法 ")

[ReadServerTime 方法](../html/c6c5f1fa-4639-0a32-8a42-48d56985e8e6.htm "ReadServerTime 方法 ")

[ReadServerTimeAsync 方法](../html/53022c3b-766e-c31b-1b55-548ceb3b2f42.htm "ReadServerTimeAsync 方法 ")

[RenameKey 方法](../html/ecf533af-233c-16b8-f88d-cb232e115013.htm "RenameKey 方法 ")

[RenameKeyAsync 方法](../html/20073ebe-0c16-2f57-a96b-79246e01129d.htm "RenameKeyAsync 方法 ")

[Save 方法](../html/261a8daa-9890-95bc-f8e4-02b78eb86277.htm "Save 方法 ")

[SaveAsync 方法](../html/88c2043f-e9e6-9d95-4a2e-455452cf8adf.htm "SaveAsync 方法 ")

[SelectDB 方法](../html/584a38b8-96e0-736f-7a24-981a8e920574.htm "SelectDB 方法 ")

[SelectDBAsync 方法](../html/462ed998-9938-a261-bde4-579b4b9d3dc1.htm "SelectDBAsync 方法 ")

[SetAdd 方法](../html/6f52901b-4d7f-70ca-1bf1-2e7146604dba.htm "SetAdd 方法 ")

[SetAddAsync 方法](../html/a59f7c04-0b20-35a4-a431-53299cb829f0.htm "SetAddAsync 方法 ")

[SetCard 方法](../html/cc0a3d1d-019b-f0ba-90a5-1bab767d6f05.htm "SetCard 方法 ")

[SetCardAsync 方法](../html/0d0665a9-a174-d154-f95b-2bf10d885bc3.htm "SetCardAsync 方法 ")

[SetDiff 方法](../html/bcc43604-dc32-a27f-53a4-a89ab29916be.htm "SetDiff 方法 ")

[SetDiffAsync 方法](../html/b70eae65-3537-3ab7-4cf9-ae6c314884c9.htm "SetDiffAsync 方法 ")

[SetDiffStore 方法](../html/7c356901-c174-c81b-50b9-ce8d60258d1e.htm "SetDiffStore 方法 ")

[SetDiffStoreAsync 方法](../html/697ac488-c2e1-8c5a-5e5d-845f676c85fd.htm "SetDiffStoreAsync 方法 ")

[SetInter 方法](../html/fa3a123d-447d-ce8b-5670-e47956d427d1.htm "SetInter 方法 ")

[SetInterAsync 方法](../html/622854f7-5e07-092e-42eb-06cc63580467.htm "SetInterAsync 方法 ")

[SetInterStore 方法](../html/996ff5a7-0d0b-5b17-546b-eb006fa5773f.htm "SetInterStore 方法 ")

[SetInterStoreAsync 方法](../html/115a5c5d-2c7b-4c9a-daec-92477fb3044a.htm "SetInterStoreAsync 方法 ")

[SetIsMember 方法](../html/1c6d465b-353a-34f2-8f46-9d4037e58911.htm "SetIsMember 方法 ")

[SetIsMemberAsync 方法](../html/21419054-3273-2aab-d486-633ab8861039.htm "SetIsMemberAsync 方法 ")

[SetMembers 方法](../html/7a5b277f-1fe1-141f-346b-3722d97ac7ca.htm "SetMembers 方法 ")

[SetMembersAsync 方法](../html/ce86155f-4be2-3836-3594-1398d78183ad.htm "SetMembersAsync 方法 ")

[SetMove 方法](../html/8b76c71a-c755-708b-e1e7-61eda8b9456c.htm "SetMove 方法 ")

[SetMoveAsync 方法](../html/dc01893c-b09e-e134-49d3-76d0dbc308b9.htm "SetMoveAsync 方法 ")

[SetPop 方法](../html/8375117d-d7b5-1f33-cb59-80cc632b736d.htm "SetPop 方法 ")

[SetPopAsync 方法](../html/66039533-6b1e-ad66-e89b-ef8f7d46c47b.htm "SetPopAsync 方法 ")

[SetRandomMember 方法](../html/2769a6d4-3b4c-75b2-c57b-abe19bb2f8e9.htm "SetRandomMember 方法 ")

[SetRandomMemberAsync 方法](../html/fb4fda18-00e6-2545-d037-6585a9487334.htm "SetRandomMemberAsync 方法 ")

[SetRemove 方法](../html/a0ba101e-935e-aca6-da5e-bb6187b42f02.htm "SetRemove 方法 ")

[SetRemoveAsync 方法](../html/f80e5db7-218b-c4bf-0c57-ec2555f6ec4e.htm "SetRemoveAsync 方法 ")

[SetUnion 方法](../html/c255e92b-1a91-f3e4-49cd-79f39d7a1079.htm "SetUnion 方法 ")

[SetUnionAsync 方法](../html/184ffefd-f6f5-b5d8-e1d6-45ae6d2d653f.htm "SetUnionAsync 方法 ")

[SetUnionStore 方法](../html/56394825-4bd1-6422-4f9b-1520ed05f927.htm "SetUnionStore 方法 ")

[SetUnionStoreAsync 方法](../html/99689f9d-ebcc-cc16-19e3-4c46f3016084.htm "SetUnionStoreAsync 方法 ")

[SubscribeMessage 方法](../html/bbaad51e-d35a-fb1f-4bb5-c001d47693be.htm "SubscribeMessage 方法 ")

[ToString 方法](../html/6ff40e11-12f4-8759-45a0-594ca3e4fcdd.htm "ToString 方法 ")

[UnSubscribeMessage 方法](../html/592327f8-7e50-4d48-c229-ca17a039b6a3.htm "UnSubscribeMessage 方法 ")

[Write(T) 方法](../html/70a08f37-0f86-76f1-55a5-df1ddd5a7d16.htm "Write(T) 方法 ")

[WriteAndPublishKey 方法](../html/69029e5e-a78a-5c83-0300-b6dd057f9308.htm "WriteAndPublishKey 方法 ")

[WriteAndPublishKeyAsync 方法](../html/d44abf4b-02bf-cfa5-cbc8-0b36edab4521.htm "WriteAndPublishKeyAsync 方法 ")

[WriteAsync(T) 方法](../html/9d54e352-2bd0-2655-3e14-49a0213e6040.htm "WriteAsync(T) 方法 ")

[WriteExpireKey 方法](../html/816ca8c0-5e3a-4609-acc3-6f227ac2dfde.htm "WriteExpireKey 方法 ")

[WriteExpireKeyAsync 方法](../html/f1dce55a-5681-32f6-2c7c-cc664465490c.htm "WriteExpireKeyAsync 方法 ")

[WriteHashKey 方法](../html/a7c80b1f-ec11-c8f4-dae2-0fa0acada685.htm "WriteHashKey 方法 ")

[WriteHashKeyAsync 方法](../html/20122681-134d-f220-683c-11c6c9201ac9.htm "WriteHashKeyAsync 方法 ")

[WriteHashKeyNx 方法](../html/819ddbb0-74d4-4667-6e55-372fac166c68.htm "WriteHashKeyNx 方法 ")

[WriteHashKeyNxAsync 方法](../html/b4571eb9-be29-fbb8-691b-f3a93d9245ec.htm "WriteHashKeyNxAsync 方法 ")

[WriteKey 方法](../html/12919427-9ab8-d765-27cc-0c3ee2c481d2.htm "WriteKey 方法 ")

[WriteKeyAsync 方法](../html/b136bcbb-3813-f660-54f1-8c5ccefb5ddf.htm "WriteKeyAsync 方法 ")

[WriteKeyIfNotExists 方法](../html/9d48312b-2ce1-077f-bf47-3e6b5a328c45.htm "WriteKeyIfNotExists 方法 ")

[WriteKeyIfNotExistsAsync 方法](../html/e6844e26-ba17-a0a8-0ff2-bb1a1e8bce8e.htm "WriteKeyIfNotExistsAsync 方法 ")

[WriteKeyRange 方法](../html/c1f98cc9-f791-4a35-bf67-ea29e8930d32.htm "WriteKeyRange 方法 ")

[WriteKeyRangeAsync 方法](../html/cdfaffec-518b-2754-5a10-3c5b68ee1767.htm "WriteKeyRangeAsync 方法 ")

[ZSetAdd 方法](../html/9c15cbbe-df53-5280-c9e4-b9d020222e0c.htm "ZSetAdd 方法 ")

[ZSetAddAsync 方法](../html/5318d2e6-d635-ed62-18e5-6573a530ab3e.htm "ZSetAddAsync 方法 ")

[ZSetCard 方法](../html/aa928601-c5f6-63fb-6bd8-e2e383796043.htm "ZSetCard 方法 ")

[ZSetCardAsync 方法](../html/b760ff15-4c7e-510f-1927-e17e6310e142.htm "ZSetCardAsync 方法 ")

[ZSetCount 方法](../html/9872ddc0-8b8d-d9d6-4d08-64001fdefd29.htm "ZSetCount 方法 ")

[ZSetCountAsync 方法](../html/b294e1ee-27c7-7326-eb4b-8669040aecf0.htm "ZSetCountAsync 方法 ")

[ZSetIncreaseBy 方法](../html/c84e0c14-15cd-451a-dec4-2a112a73772e.htm "ZSetIncreaseBy 方法 ")

[ZSetIncreaseByAsync 方法](../html/9805ac6b-aa05-2152-7dbb-7f8077ea462b.htm "ZSetIncreaseByAsync 方法 ")

[ZSetRange 方法](../html/759e2523-e2ad-0dc1-a904-354031a8ae50.htm "ZSetRange 方法 ")

[ZSetRangeAsync 方法](../html/7724eb62-2653-65a7-779d-391d630c2faf.htm "ZSetRangeAsync 方法 ")

[ZSetRangeByScore 方法](../html/acab27d0-6b6c-1f25-00ad-fafd4ac068fa.htm "ZSetRangeByScore 方法 ")

[ZSetRangeByScoreAsync 方法](../html/d1e0aa3c-d5f9-caeb-4017-699f956e661d.htm "ZSetRangeByScoreAsync 方法 ")

[ZSetRank 方法](../html/7957411f-32c1-4e07-a4f2-9d1948c63a2b.htm "ZSetRank 方法 ")

[ZSetRankAsync 方法](../html/1181fd8c-02aa-9d53-cbeb-35d542bed241.htm "ZSetRankAsync 方法 ")

[ZSetRemove 方法](../html/a76625ea-7a7b-161d-a039-02f4f5a8e997.htm "ZSetRemove 方法 ")

[ZSetRemoveAsync 方法](../html/2f1b3de6-ad74-e2b7-3350-79d9afe13c9d.htm "ZSetRemoveAsync 方法 ")

[ZSetRemoveRangeByRank 方法](../html/cb81dfc2-6e30-0aca-8bcb-348cdc32bc35.htm "ZSetRemoveRangeByRank 方法 ")

[ZSetRemoveRangeByRankAsync 方法](../html/9bde4b5f-b9bd-e01b-feb4-82c9b532dee4.htm "ZSetRemoveRangeByRankAsync 方法 ")

[ZSetRemoveRangeByScore 方法](../html/ac7be899-fe4e-3f12-d60f-81a038a927c7.htm "ZSetRemoveRangeByScore 方法 ")

[ZSetRemoveRangeByScoreAsync 方法](../html/6ccee6bf-8635-14bf-2e7f-d1ac351c70c9.htm "ZSetRemoveRangeByScoreAsync 方法 ")

[ZSetReverseRange 方法](../html/ef6defc7-62a2-c8a4-4025-32f6b3a82ae4.htm "ZSetReverseRange 方法 ")

[ZSetReverseRangeAsync 方法](../html/b40d37b7-1118-6aec-c7da-85ec54b457d0.htm "ZSetReverseRangeAsync 方法 ")

[ZSetReverseRangeByScore 方法](../html/3c1a26a8-f8e5-a835-7051-94b4d9c851ab.htm "ZSetReverseRangeByScore 方法 ")

[ZSetReverseRangeByScoreAsync 方法](../html/3a3cfbff-cee7-e1f3-d9c1-e7528db4e998.htm "ZSetReverseRangeByScoreAsync 方法 ")

[ZSetReverseRank 方法](../html/d524195b-6b02-a8e8-330a-3454a6ed21cd.htm "ZSetReverseRank 方法 ")

[ZSetReverseRankAsync 方法](../html/c010781c-2851-7158-f953-1779929e3242.htm "ZSetReverseRankAsync 方法 ")

[ZSetScore 方法](../html/a98ee7e4-f361-f403-975b-a197e0ca0268.htm "ZSetScore 方法 ")

[ZSetScoreAsync 方法](../html/fd24c433-944e-2ecb-4bcd-266bfc6825d4.htm "ZSetScoreAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| RedisClientChangePasswordAsync 方法 |

修改Redis的密码信息，如果不需要密码，则传入空字符串即可

**命名空间：**
 [HslCommunication.Enthernet.Redis](cc44ce00-3781-c690-b233-24000870bf42.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public Task<OperateResult> ChangePasswordAsync(
	string password
)
```

```
Public Function ChangePasswordAsync ( 
	password As String
) As Task(Of OperateResult)
```

```
public:
Task<OperateResult^>^ ChangePasswordAsync(
	String^ password
)
```

```
member ChangePasswordAsync : 
        password : string -> Task<OperateResult> 
```

#### 参数

password
:   类型：SystemString  
    密码信息

#### 返回值

类型：Task[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否更新了密码信息

![](../icons/SectionExpanded.png)参见

#### 引用

[RedisClient 类](e33d9088-6df0-32e9-33f1-426e944c3b7f.htm)

[HslCommunication.Enthernet.Redis 命名空间](cc44ce00-3781-c690-b233-24000870bf42.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## DBSize 方法 

[原文連結](http://api.hslcommunication.cn/html/5814129a-9a73-25ac-e726-fe1da141eb28.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Enthernet.Redis](../html/cc44ce00-3781-c690-b233-24000870bf42.htm "HslCommunication.Enthernet.Redis")

[RedisClient 类](../html/e33d9088-6df0-32e9-33f1-426e944c3b7f.htm "RedisClient 类")

[RedisClient 方法](../html/e41be4a5-1490-ffbd-f241-938f2d6c3f07.htm "RedisClient 方法")

[AppendKey 方法](../html/f2521b8c-0fb1-6b34-4adc-ca6804c7195f.htm "AppendKey 方法 ")

[AppendKeyAsync 方法](../html/c39b4c6a-b7f6-68fe-dfc9-7cce1b653e24.htm "AppendKeyAsync 方法 ")

[ChangePassword 方法](../html/78a37d4d-bdee-7a09-d0ce-dc245627b1d1.htm "ChangePassword 方法 ")

[ChangePasswordAsync 方法](../html/58f6f736-093d-6e4d-3ad3-19590fa2a3a3.htm "ChangePasswordAsync 方法 ")

[DBSize 方法](../html/5814129a-9a73-25ac-e726-fe1da141eb28.htm "DBSize 方法 ")

[DBSizeAsync 方法](../html/243a1f5f-a474-2e20-1298-06eff125f8b5.htm "DBSizeAsync 方法 ")

[DecrementKey 方法](../html/e9fa6c94-43ca-8291-79ca-6e04ca9289e7.htm "DecrementKey 方法 ")

[DecrementKeyAsync 方法](../html/6837e351-b8ef-a90e-c5e2-73b31621f837.htm "DecrementKeyAsync 方法 ")

[DeleteHashKey 方法](../html/19e3ebc3-e6fc-71cf-e80d-cb73e3011170.htm "DeleteHashKey 方法 ")

[DeleteHashKeyAsync 方法](../html/0dcc0d8a-1ed1-5c81-fb0e-e1c93be9d579.htm "DeleteHashKeyAsync 方法 ")

[DeleteKey 方法](../html/e2f0f305-d8bc-d515-3a62-cc91a3c1de76.htm "DeleteKey 方法 ")

[DeleteKeyAsync 方法](../html/0d8b5b0c-5979-9140-4ccf-992f4eaf1553.htm "DeleteKeyAsync 方法 ")

[ExistsHashKey 方法](../html/5b0bbeb9-4ab9-1f83-9e09-5f7227ac28e6.htm "ExistsHashKey 方法 ")

[ExistsHashKeyAsync 方法](../html/a80d408e-d621-a502-7829-6c7ea2a081f9.htm "ExistsHashKeyAsync 方法 ")

[ExistsKey 方法](../html/2b1b5b08-13ed-88c3-445d-d3613b48bc95.htm "ExistsKey 方法 ")

[ExistsKeyAsync 方法](../html/0d4357ed-5de3-1444-fce2-0e9be619d20f.htm "ExistsKeyAsync 方法 ")

[ExpireKey 方法](../html/0272044f-cc55-6b6a-42f0-6a0d89f37fe7.htm "ExpireKey 方法 ")

[ExpireKeyAsync 方法](../html/53c18b2d-2169-0519-602e-f65302f490e1.htm "ExpireKeyAsync 方法 ")

[ExtraOnDisconnect 方法](../html/49d9b263-02d8-e822-6a33-49f93737b98b.htm "ExtraOnDisconnect 方法 ")

[ExtraOnDisconnectAsync 方法](../html/f681350d-2aad-40e0-0824-30fbd3fe57ff.htm "ExtraOnDisconnectAsync 方法 ")

[FlushDB 方法](../html/38c28c3f-63bf-f1a0-67bc-6547799257f9.htm "FlushDB 方法 ")

[FlushDBAsync 方法](../html/40e3ccb2-9fe7-c0d7-88b9-95266a13c336.htm "FlushDBAsync 方法 ")

[GetListLength 方法](../html/36330604-2557-4fc1-bdd7-fc3aff05e997.htm "GetListLength 方法 ")

[GetListLengthAsync 方法](../html/0e09431d-45b3-b539-a618-a147ae377f88.htm "GetListLengthAsync 方法 ")

[IncrementHashKey 方法](../html/6404f402-9882-5453-7df7-275f61748021.htm "IncrementHashKey 方法 ")

[IncrementHashKeyAsync 方法](../html/7a69e1a1-8689-3e23-8aa3-cdad33e7e3f8.htm "IncrementHashKeyAsync 方法 ")

[IncrementKey 方法](../html/04d46420-bedd-9016-075c-3bb1a6623583.htm "IncrementKey 方法 ")

[IncrementKeyAsync 方法](../html/95f88e06-114c-b5b3-ee2f-2e9f97f27b8c.htm "IncrementKeyAsync 方法 ")

[InitializationOnConnect 方法](../html/72414fc8-b43a-3377-acdc-c59ae0628e11.htm "InitializationOnConnect 方法 ")

[InitializationOnConnectAsync 方法](../html/9156d090-6b1d-7381-c396-81491b4b6756.htm "InitializationOnConnectAsync 方法 ")

[ListInsertAfter 方法](../html/11b6eb28-b2ee-57a5-c724-a99df9101c71.htm "ListInsertAfter 方法 ")

[ListInsertAfterAsync 方法](../html/9defd0d6-5681-4c26-0973-14ecca4fc821.htm "ListInsertAfterAsync 方法 ")

[ListInsertBefore 方法](../html/b19eed6b-2778-4b91-6f6b-e75cd9f1088e.htm "ListInsertBefore 方法 ")

[ListInsertBeforeAsync 方法](../html/0bf5390c-af9f-ecd6-1148-34a20abee73d.htm "ListInsertBeforeAsync 方法 ")

[ListLeftPop 方法](../html/8865ed9c-69dc-2f28-7f53-42e1c03d152b.htm "ListLeftPop 方法 ")

[ListLeftPopAsync 方法](../html/f4d7f464-73f7-ff5d-0834-e18b5a91c994.htm "ListLeftPopAsync 方法 ")

[ListLeftPush 方法](../html/232a2dbd-e2c6-5fbf-af96-5269a10df4c9.htm "ListLeftPush 方法 ")

[ListLeftPushAsync 方法](../html/cb91133d-ce36-f4f7-a9d3-81dbda7d6617.htm "ListLeftPushAsync 方法 ")

[ListLeftPushX 方法](../html/33093159-9864-a0dd-6a09-8bf134942031.htm "ListLeftPushX 方法 ")

[ListLeftPushXAsync 方法](../html/6bd2ab31-ff58-71f0-b68f-f36d7f53d845.htm "ListLeftPushXAsync 方法 ")

[ListRange 方法](../html/efe06eb0-b945-fe65-5e81-21888cdda04f.htm "ListRange 方法 ")

[ListRangeAsync 方法](../html/32e344f9-cdac-fc0e-2c50-bf2dd54b4f98.htm "ListRangeAsync 方法 ")

[ListRemoveElementMatch 方法](../html/0ab40e89-d265-0b48-05de-5c1a7b41d326.htm "ListRemoveElementMatch 方法 ")

[ListRemoveElementMatchAsync 方法](../html/5610de76-3198-5b5a-8f3d-b8fc937d7c83.htm "ListRemoveElementMatchAsync 方法 ")

[ListRightPop 方法](../html/97855b86-6a07-76e2-b887-1e830be705de.htm "ListRightPop 方法 ")

[ListRightPopAsync 方法](../html/fde8cee2-f393-a463-d561-74568545edc7.htm "ListRightPopAsync 方法 ")

[ListRightPopLeftPush 方法](../html/e0b58cc4-4d80-64d6-729c-0dcdaafc4701.htm "ListRightPopLeftPush 方法 ")

[ListRightPopLeftPushAsync 方法](../html/41069df2-0dad-6f01-895f-a91086a92fa9.htm "ListRightPopLeftPushAsync 方法 ")

[ListRightPush 方法](../html/d71ca156-74d3-536e-6925-315d948c0ace.htm "ListRightPush 方法 ")

[ListRightPushAsync 方法](../html/53144654-da3d-b367-b39c-0dbc68275e59.htm "ListRightPushAsync 方法 ")

[ListRightPushX 方法](../html/e6ccfd1a-25ee-5eff-d56c-0ff9eb017fbc.htm "ListRightPushX 方法 ")

[ListRightPushXAsync 方法](../html/3252db2a-56cf-1aa2-7716-f9eb2596ec62.htm "ListRightPushXAsync 方法 ")

[ListSet 方法](../html/5f2e7b8b-b0e7-87e6-dac7-81b572cf3fa3.htm "ListSet 方法 ")

[ListSetAsync 方法](../html/ca019ec2-eb32-41f8-a460-89c70f90ace9.htm "ListSetAsync 方法 ")

[ListTrim 方法](../html/acf21364-1692-f9da-d337-c460b12a012d.htm "ListTrim 方法 ")

[ListTrimAsync 方法](../html/9e86492d-205c-8427-9052-50ea3d8c1c64.htm "ListTrimAsync 方法 ")

[MoveKey 方法](../html/5aa50043-85b9-a245-606a-de2b69164218.htm "MoveKey 方法 ")

[MoveKeyAsync 方法](../html/a93a74bf-65bd-733b-e39a-06d705252d99.htm "MoveKeyAsync 方法 ")

[OperateLongNumberFromServer 方法](../html/d2cdcd44-a6fb-7aa7-8c9c-08ff43107f65.htm "OperateLongNumberFromServer 方法 ")

[OperateLongNumberFromServerAsync 方法](../html/0005132d-e8c6-d51b-69dc-c9f8feaf8a9d.htm "OperateLongNumberFromServerAsync 方法 ")

[OperateNumberFromServer 方法](../html/69aecc05-6901-f53a-da37-ebf85bae7053.htm "OperateNumberFromServer 方法 ")

[OperateNumberFromServerAsync 方法](../html/b53bd6fa-d215-55f7-bbac-2ee171ed7498.htm "OperateNumberFromServerAsync 方法 ")

[OperateStatusFromServer 方法](../html/d460c3bb-e484-3fcb-e206-d20dc7201d63.htm "OperateStatusFromServer 方法 ")

[OperateStatusFromServerAsync 方法](../html/e1a241a3-e200-747f-bb8b-5b7578a08999.htm "OperateStatusFromServerAsync 方法 ")

[OperateStringFromServer 方法](../html/822914af-cf62-1ccc-d22e-2df7dfac7294.htm "OperateStringFromServer 方法 ")

[OperateStringFromServerAsync 方法](../html/eb31e52d-5387-eb14-291e-22f0e64666a4.htm "OperateStringFromServerAsync 方法 ")

[OperateStringsFromServer 方法](../html/c85c1a0e-83f9-69c9-a9b9-8a947ca54583.htm "OperateStringsFromServer 方法 ")

[OperateStringsFromServerAsync 方法](../html/2db50b51-32b8-1e92-cb9a-f2c847437e50.htm "OperateStringsFromServerAsync 方法 ")

[PersistKey 方法](../html/c68a840b-ce64-3303-9339-a10f88bcbb77.htm "PersistKey 方法 ")

[PersistKeyAsync 方法](../html/a5d7ee67-9041-6455-7fc1-c6e96a442799.htm "PersistKeyAsync 方法 ")

[Ping 方法](../html/9cdd4e99-90a2-0ae8-141c-26f0f4a4d698.htm "Ping 方法 ")

[PingAsync 方法](../html/d6c8e625-edb2-7ff2-a54d-e852d49c73c8.htm "PingAsync 方法 ")

[Publish 方法](../html/74f20db5-aaaa-af2e-b4e6-5d192196b222.htm "Publish 方法 ")

[PublishAsync 方法](../html/61428af8-94a2-b2b0-419b-5c0d3df4bee5.htm "PublishAsync 方法 ")

[Read(T) 方法](../html/08365a89-e2d7-2aa5-3fd5-ab1a5d60c440.htm "Read(T) 方法 ")

[ReadAllKeys 方法](../html/53c48696-f991-f734-12b0-1ad306a84c6b.htm "ReadAllKeys 方法 ")

[ReadAllKeysAsync 方法](../html/84c1be8c-3286-0410-0bd1-bb0185ff197e.htm "ReadAllKeysAsync 方法 ")

[ReadAndWriteKey 方法](../html/a6908f0e-1991-5210-f6c6-1e8b06369e56.htm "ReadAndWriteKey 方法 ")

[ReadAndWriteKeyAsync 方法](../html/6686726f-93ba-6435-76b0-b501c31e637e.htm "ReadAndWriteKeyAsync 方法 ")

[ReadAsync(T) 方法](../html/43a59914-1502-1353-0bf3-20281ac92ca9.htm "ReadAsync(T) 方法 ")

[ReadCustomer 方法](../html/a4b7f43e-34ba-8aac-a705-a9bc8d4a5298.htm "ReadCustomer 方法 ")

[ReadCustomerAsync 方法](../html/bafea5ed-4293-ae60-c8e7-bb3d3481db44.htm "ReadCustomerAsync 方法 ")

[ReadFromCoreServer 方法](../html/a5985ec7-e909-5a40-b2b9-7160363a1156.htm "ReadFromCoreServer 方法 ")

[ReadFromCoreServerAsync 方法](../html/984a89de-20da-e7ed-e471-3e889955aa31.htm "ReadFromCoreServerAsync 方法 ")

[ReadHashKey 方法](../html/6e8fc531-b2d6-2b1e-4885-1b937745652f.htm "ReadHashKey 方法 ")

[ReadHashKeyAll 方法](../html/63828fbe-fbdf-772a-010b-1425b0584572.htm "ReadHashKeyAll 方法 ")

[ReadHashKeyAllAsync 方法](../html/31c54864-f7bc-b3ce-7e77-15f0c1dcc540.htm "ReadHashKeyAllAsync 方法 ")

[ReadHashKeyAsync 方法](../html/36c995b4-42a7-e264-8a03-3390e24376ec.htm "ReadHashKeyAsync 方法 ")

[ReadHashKeyLength 方法](../html/32d483da-49e0-a205-bdc4-a659711be1c0.htm "ReadHashKeyLength 方法 ")

[ReadHashKeyLengthAsync 方法](../html/5d0d4570-1f16-4277-cff7-a40607a530ab.htm "ReadHashKeyLengthAsync 方法 ")

[ReadHashKeys 方法](../html/420286f8-5876-ee1b-b328-0b6eba7cefac.htm "ReadHashKeys 方法 ")

[ReadHashKeysAsync 方法](../html/bbc8310c-029d-1ea5-b0e3-239b29ee8743.htm "ReadHashKeysAsync 方法 ")

[ReadHashValues 方法](../html/1d8b8e06-3d34-f4c6-68d0-20c14b309de1.htm "ReadHashValues 方法 ")

[ReadHashValuesAsync 方法](../html/403e7344-c665-def6-70d9-81e5b22a740b.htm "ReadHashValuesAsync 方法 ")

[ReadKey 方法](../html/93fa3dab-d22c-d028-00c3-e81391085e6f.htm "ReadKey 方法 ")

[ReadKeyAsync 方法](../html/48cee995-f908-c3a2-ab70-829d789f58a5.htm "ReadKeyAsync 方法 ")

[ReadKeyLength 方法](../html/7746357b-7e1b-d150-d019-722b752f5eeb.htm "ReadKeyLength 方法 ")

[ReadKeyLengthAsync 方法](../html/7dd1b305-77b1-920c-14c0-2fe38e970791.htm "ReadKeyLengthAsync 方法 ")

[ReadKeyRange 方法](../html/2e3e9582-cb24-44e6-d78b-7c672acdd4c8.htm "ReadKeyRange 方法 ")

[ReadKeyRangeAsync 方法](../html/407fdd5b-8b57-75d2-63be-93294311cb30.htm "ReadKeyRangeAsync 方法 ")

[ReadKeyTTL 方法](../html/679079c3-36e5-1618-3409-18459a74a0c4.htm "ReadKeyTTL 方法 ")

[ReadKeyTTLAsync 方法](../html/24887926-0fd5-b65f-d070-1d439a5273ef.htm "ReadKeyTTLAsync 方法 ")

[ReadKeyType 方法](../html/1f702770-7804-880a-038e-c002ed31c42f.htm "ReadKeyType 方法 ")

[ReadKeyTypeAsync 方法](../html/91c104ff-a646-76de-9bed-5e64237a57fc.htm "ReadKeyTypeAsync 方法 ")

[ReadListByIndex 方法](../html/43c7dae8-4c01-d2c9-b38b-0e7bef68b581.htm "ReadListByIndex 方法 ")

[ReadListByIndexAsync 方法](../html/1e1163ad-8161-9390-cff2-7919c8145c92.htm "ReadListByIndexAsync 方法 ")

[ReadRandomKey 方法](../html/0d358a8f-b9e0-c659-c218-f051aaedcbfe.htm "ReadRandomKey 方法 ")

[ReadRandomKeyAsync 方法](../html/b97f91d2-2ab6-a39b-86f7-71238cd44fac.htm "ReadRandomKeyAsync 方法 ")

[ReadServerTime 方法](../html/c6c5f1fa-4639-0a32-8a42-48d56985e8e6.htm "ReadServerTime 方法 ")

[ReadServerTimeAsync 方法](../html/53022c3b-766e-c31b-1b55-548ceb3b2f42.htm "ReadServerTimeAsync 方法 ")

[RenameKey 方法](../html/ecf533af-233c-16b8-f88d-cb232e115013.htm "RenameKey 方法 ")

[RenameKeyAsync 方法](../html/20073ebe-0c16-2f57-a96b-79246e01129d.htm "RenameKeyAsync 方法 ")

[Save 方法](../html/261a8daa-9890-95bc-f8e4-02b78eb86277.htm "Save 方法 ")

[SaveAsync 方法](../html/88c2043f-e9e6-9d95-4a2e-455452cf8adf.htm "SaveAsync 方法 ")

[SelectDB 方法](../html/584a38b8-96e0-736f-7a24-981a8e920574.htm "SelectDB 方法 ")

[SelectDBAsync 方法](../html/462ed998-9938-a261-bde4-579b4b9d3dc1.htm "SelectDBAsync 方法 ")

[SetAdd 方法](../html/6f52901b-4d7f-70ca-1bf1-2e7146604dba.htm "SetAdd 方法 ")

[SetAddAsync 方法](../html/a59f7c04-0b20-35a4-a431-53299cb829f0.htm "SetAddAsync 方法 ")

[SetCard 方法](../html/cc0a3d1d-019b-f0ba-90a5-1bab767d6f05.htm "SetCard 方法 ")

[SetCardAsync 方法](../html/0d0665a9-a174-d154-f95b-2bf10d885bc3.htm "SetCardAsync 方法 ")

[SetDiff 方法](../html/bcc43604-dc32-a27f-53a4-a89ab29916be.htm "SetDiff 方法 ")

[SetDiffAsync 方法](../html/b70eae65-3537-3ab7-4cf9-ae6c314884c9.htm "SetDiffAsync 方法 ")

[SetDiffStore 方法](../html/7c356901-c174-c81b-50b9-ce8d60258d1e.htm "SetDiffStore 方法 ")

[SetDiffStoreAsync 方法](../html/697ac488-c2e1-8c5a-5e5d-845f676c85fd.htm "SetDiffStoreAsync 方法 ")

[SetInter 方法](../html/fa3a123d-447d-ce8b-5670-e47956d427d1.htm "SetInter 方法 ")

[SetInterAsync 方法](../html/622854f7-5e07-092e-42eb-06cc63580467.htm "SetInterAsync 方法 ")

[SetInterStore 方法](../html/996ff5a7-0d0b-5b17-546b-eb006fa5773f.htm "SetInterStore 方法 ")

[SetInterStoreAsync 方法](../html/115a5c5d-2c7b-4c9a-daec-92477fb3044a.htm "SetInterStoreAsync 方法 ")

[SetIsMember 方法](../html/1c6d465b-353a-34f2-8f46-9d4037e58911.htm "SetIsMember 方法 ")

[SetIsMemberAsync 方法](../html/21419054-3273-2aab-d486-633ab8861039.htm "SetIsMemberAsync 方法 ")

[SetMembers 方法](../html/7a5b277f-1fe1-141f-346b-3722d97ac7ca.htm "SetMembers 方法 ")

[SetMembersAsync 方法](../html/ce86155f-4be2-3836-3594-1398d78183ad.htm "SetMembersAsync 方法 ")

[SetMove 方法](../html/8b76c71a-c755-708b-e1e7-61eda8b9456c.htm "SetMove 方法 ")

[SetMoveAsync 方法](../html/dc01893c-b09e-e134-49d3-76d0dbc308b9.htm "SetMoveAsync 方法 ")

[SetPop 方法](../html/8375117d-d7b5-1f33-cb59-80cc632b736d.htm "SetPop 方法 ")

[SetPopAsync 方法](../html/66039533-6b1e-ad66-e89b-ef8f7d46c47b.htm "SetPopAsync 方法 ")

[SetRandomMember 方法](../html/2769a6d4-3b4c-75b2-c57b-abe19bb2f8e9.htm "SetRandomMember 方法 ")

[SetRandomMemberAsync 方法](../html/fb4fda18-00e6-2545-d037-6585a9487334.htm "SetRandomMemberAsync 方法 ")

[SetRemove 方法](../html/a0ba101e-935e-aca6-da5e-bb6187b42f02.htm "SetRemove 方法 ")

[SetRemoveAsync 方法](../html/f80e5db7-218b-c4bf-0c57-ec2555f6ec4e.htm "SetRemoveAsync 方法 ")

[SetUnion 方法](../html/c255e92b-1a91-f3e4-49cd-79f39d7a1079.htm "SetUnion 方法 ")

[SetUnionAsync 方法](../html/184ffefd-f6f5-b5d8-e1d6-45ae6d2d653f.htm "SetUnionAsync 方法 ")

[SetUnionStore 方法](../html/56394825-4bd1-6422-4f9b-1520ed05f927.htm "SetUnionStore 方法 ")

[SetUnionStoreAsync 方法](../html/99689f9d-ebcc-cc16-19e3-4c46f3016084.htm "SetUnionStoreAsync 方法 ")

[SubscribeMessage 方法](../html/bbaad51e-d35a-fb1f-4bb5-c001d47693be.htm "SubscribeMessage 方法 ")

[ToString 方法](../html/6ff40e11-12f4-8759-45a0-594ca3e4fcdd.htm "ToString 方法 ")

[UnSubscribeMessage 方法](../html/592327f8-7e50-4d48-c229-ca17a039b6a3.htm "UnSubscribeMessage 方法 ")

[Write(T) 方法](../html/70a08f37-0f86-76f1-55a5-df1ddd5a7d16.htm "Write(T) 方法 ")

[WriteAndPublishKey 方法](../html/69029e5e-a78a-5c83-0300-b6dd057f9308.htm "WriteAndPublishKey 方法 ")

[WriteAndPublishKeyAsync 方法](../html/d44abf4b-02bf-cfa5-cbc8-0b36edab4521.htm "WriteAndPublishKeyAsync 方法 ")

[WriteAsync(T) 方法](../html/9d54e352-2bd0-2655-3e14-49a0213e6040.htm "WriteAsync(T) 方法 ")

[WriteExpireKey 方法](../html/816ca8c0-5e3a-4609-acc3-6f227ac2dfde.htm "WriteExpireKey 方法 ")

[WriteExpireKeyAsync 方法](../html/f1dce55a-5681-32f6-2c7c-cc664465490c.htm "WriteExpireKeyAsync 方法 ")

[WriteHashKey 方法](../html/a7c80b1f-ec11-c8f4-dae2-0fa0acada685.htm "WriteHashKey 方法 ")

[WriteHashKeyAsync 方法](../html/20122681-134d-f220-683c-11c6c9201ac9.htm "WriteHashKeyAsync 方法 ")

[WriteHashKeyNx 方法](../html/819ddbb0-74d4-4667-6e55-372fac166c68.htm "WriteHashKeyNx 方法 ")

[WriteHashKeyNxAsync 方法](../html/b4571eb9-be29-fbb8-691b-f3a93d9245ec.htm "WriteHashKeyNxAsync 方法 ")

[WriteKey 方法](../html/12919427-9ab8-d765-27cc-0c3ee2c481d2.htm "WriteKey 方法 ")

[WriteKeyAsync 方法](../html/b136bcbb-3813-f660-54f1-8c5ccefb5ddf.htm "WriteKeyAsync 方法 ")

[WriteKeyIfNotExists 方法](../html/9d48312b-2ce1-077f-bf47-3e6b5a328c45.htm "WriteKeyIfNotExists 方法 ")

[WriteKeyIfNotExistsAsync 方法](../html/e6844e26-ba17-a0a8-0ff2-bb1a1e8bce8e.htm "WriteKeyIfNotExistsAsync 方法 ")

[WriteKeyRange 方法](../html/c1f98cc9-f791-4a35-bf67-ea29e8930d32.htm "WriteKeyRange 方法 ")

[WriteKeyRangeAsync 方法](../html/cdfaffec-518b-2754-5a10-3c5b68ee1767.htm "WriteKeyRangeAsync 方法 ")

[ZSetAdd 方法](../html/9c15cbbe-df53-5280-c9e4-b9d020222e0c.htm "ZSetAdd 方法 ")

[ZSetAddAsync 方法](../html/5318d2e6-d635-ed62-18e5-6573a530ab3e.htm "ZSetAddAsync 方法 ")

[ZSetCard 方法](../html/aa928601-c5f6-63fb-6bd8-e2e383796043.htm "ZSetCard 方法 ")

[ZSetCardAsync 方法](../html/b760ff15-4c7e-510f-1927-e17e6310e142.htm "ZSetCardAsync 方法 ")

[ZSetCount 方法](../html/9872ddc0-8b8d-d9d6-4d08-64001fdefd29.htm "ZSetCount 方法 ")

[ZSetCountAsync 方法](../html/b294e1ee-27c7-7326-eb4b-8669040aecf0.htm "ZSetCountAsync 方法 ")

[ZSetIncreaseBy 方法](../html/c84e0c14-15cd-451a-dec4-2a112a73772e.htm "ZSetIncreaseBy 方法 ")

[ZSetIncreaseByAsync 方法](../html/9805ac6b-aa05-2152-7dbb-7f8077ea462b.htm "ZSetIncreaseByAsync 方法 ")

[ZSetRange 方法](../html/759e2523-e2ad-0dc1-a904-354031a8ae50.htm "ZSetRange 方法 ")

[ZSetRangeAsync 方法](../html/7724eb62-2653-65a7-779d-391d630c2faf.htm "ZSetRangeAsync 方法 ")

[ZSetRangeByScore 方法](../html/acab27d0-6b6c-1f25-00ad-fafd4ac068fa.htm "ZSetRangeByScore 方法 ")

[ZSetRangeByScoreAsync 方法](../html/d1e0aa3c-d5f9-caeb-4017-699f956e661d.htm "ZSetRangeByScoreAsync 方法 ")

[ZSetRank 方法](../html/7957411f-32c1-4e07-a4f2-9d1948c63a2b.htm "ZSetRank 方法 ")

[ZSetRankAsync 方法](../html/1181fd8c-02aa-9d53-cbeb-35d542bed241.htm "ZSetRankAsync 方法 ")

[ZSetRemove 方法](../html/a76625ea-7a7b-161d-a039-02f4f5a8e997.htm "ZSetRemove 方法 ")

[ZSetRemoveAsync 方法](../html/2f1b3de6-ad74-e2b7-3350-79d9afe13c9d.htm "ZSetRemoveAsync 方法 ")

[ZSetRemoveRangeByRank 方法](../html/cb81dfc2-6e30-0aca-8bcb-348cdc32bc35.htm "ZSetRemoveRangeByRank 方法 ")

[ZSetRemoveRangeByRankAsync 方法](../html/9bde4b5f-b9bd-e01b-feb4-82c9b532dee4.htm "ZSetRemoveRangeByRankAsync 方法 ")

[ZSetRemoveRangeByScore 方法](../html/ac7be899-fe4e-3f12-d60f-81a038a927c7.htm "ZSetRemoveRangeByScore 方法 ")

[ZSetRemoveRangeByScoreAsync 方法](../html/6ccee6bf-8635-14bf-2e7f-d1ac351c70c9.htm "ZSetRemoveRangeByScoreAsync 方法 ")

[ZSetReverseRange 方法](../html/ef6defc7-62a2-c8a4-4025-32f6b3a82ae4.htm "ZSetReverseRange 方法 ")

[ZSetReverseRangeAsync 方法](../html/b40d37b7-1118-6aec-c7da-85ec54b457d0.htm "ZSetReverseRangeAsync 方法 ")

[ZSetReverseRangeByScore 方法](../html/3c1a26a8-f8e5-a835-7051-94b4d9c851ab.htm "ZSetReverseRangeByScore 方法 ")

[ZSetReverseRangeByScoreAsync 方法](../html/3a3cfbff-cee7-e1f3-d9c1-e7528db4e998.htm "ZSetReverseRangeByScoreAsync 方法 ")

[ZSetReverseRank 方法](../html/d524195b-6b02-a8e8-330a-3454a6ed21cd.htm "ZSetReverseRank 方法 ")

[ZSetReverseRankAsync 方法](../html/c010781c-2851-7158-f953-1779929e3242.htm "ZSetReverseRankAsync 方法 ")

[ZSetScore 方法](../html/a98ee7e4-f361-f403-975b-a197e0ca0268.htm "ZSetScore 方法 ")

[ZSetScoreAsync 方法](../html/fd24c433-944e-2ecb-4bcd-266bfc6825d4.htm "ZSetScoreAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| RedisClientDBSize 方法 |

返回当前数据库的 key 的数量。

**命名空间：**
 [HslCommunication.Enthernet.Redis](cc44ce00-3781-c690-b233-24000870bf42.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public OperateResult<long> DBSize()
```

```
Public Function DBSize As OperateResult(Of Long)
```

```
public:
OperateResult<long long>^ DBSize()
```

```
member DBSize : unit -> OperateResult<int64> 
```

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Int64  
当前数据库的 key 的数量。

![](../icons/SectionExpanded.png)参见

#### 引用

[RedisClient 类](e33d9088-6df0-32e9-33f1-426e944c3b7f.htm)

[HslCommunication.Enthernet.Redis 命名空间](cc44ce00-3781-c690-b233-24000870bf42.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## DBSizeAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/243a1f5f-a474-2e20-1298-06eff125f8b5.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Enthernet.Redis](../html/cc44ce00-3781-c690-b233-24000870bf42.htm "HslCommunication.Enthernet.Redis")

[RedisClient 类](../html/e33d9088-6df0-32e9-33f1-426e944c3b7f.htm "RedisClient 类")

[RedisClient 方法](../html/e41be4a5-1490-ffbd-f241-938f2d6c3f07.htm "RedisClient 方法")

[AppendKey 方法](../html/f2521b8c-0fb1-6b34-4adc-ca6804c7195f.htm "AppendKey 方法 ")

[AppendKeyAsync 方法](../html/c39b4c6a-b7f6-68fe-dfc9-7cce1b653e24.htm "AppendKeyAsync 方法 ")

[ChangePassword 方法](../html/78a37d4d-bdee-7a09-d0ce-dc245627b1d1.htm "ChangePassword 方法 ")

[ChangePasswordAsync 方法](../html/58f6f736-093d-6e4d-3ad3-19590fa2a3a3.htm "ChangePasswordAsync 方法 ")

[DBSize 方法](../html/5814129a-9a73-25ac-e726-fe1da141eb28.htm "DBSize 方法 ")

[DBSizeAsync 方法](../html/243a1f5f-a474-2e20-1298-06eff125f8b5.htm "DBSizeAsync 方法 ")

[DecrementKey 方法](../html/e9fa6c94-43ca-8291-79ca-6e04ca9289e7.htm "DecrementKey 方法 ")

[DecrementKeyAsync 方法](../html/6837e351-b8ef-a90e-c5e2-73b31621f837.htm "DecrementKeyAsync 方法 ")

[DeleteHashKey 方法](../html/19e3ebc3-e6fc-71cf-e80d-cb73e3011170.htm "DeleteHashKey 方法 ")

[DeleteHashKeyAsync 方法](../html/0dcc0d8a-1ed1-5c81-fb0e-e1c93be9d579.htm "DeleteHashKeyAsync 方法 ")

[DeleteKey 方法](../html/e2f0f305-d8bc-d515-3a62-cc91a3c1de76.htm "DeleteKey 方法 ")

[DeleteKeyAsync 方法](../html/0d8b5b0c-5979-9140-4ccf-992f4eaf1553.htm "DeleteKeyAsync 方法 ")

[ExistsHashKey 方法](../html/5b0bbeb9-4ab9-1f83-9e09-5f7227ac28e6.htm "ExistsHashKey 方法 ")

[ExistsHashKeyAsync 方法](../html/a80d408e-d621-a502-7829-6c7ea2a081f9.htm "ExistsHashKeyAsync 方法 ")

[ExistsKey 方法](../html/2b1b5b08-13ed-88c3-445d-d3613b48bc95.htm "ExistsKey 方法 ")

[ExistsKeyAsync 方法](../html/0d4357ed-5de3-1444-fce2-0e9be619d20f.htm "ExistsKeyAsync 方法 ")

[ExpireKey 方法](../html/0272044f-cc55-6b6a-42f0-6a0d89f37fe7.htm "ExpireKey 方法 ")

[ExpireKeyAsync 方法](../html/53c18b2d-2169-0519-602e-f65302f490e1.htm "ExpireKeyAsync 方法 ")

[ExtraOnDisconnect 方法](../html/49d9b263-02d8-e822-6a33-49f93737b98b.htm "ExtraOnDisconnect 方法 ")

[ExtraOnDisconnectAsync 方法](../html/f681350d-2aad-40e0-0824-30fbd3fe57ff.htm "ExtraOnDisconnectAsync 方法 ")

[FlushDB 方法](../html/38c28c3f-63bf-f1a0-67bc-6547799257f9.htm "FlushDB 方法 ")

[FlushDBAsync 方法](../html/40e3ccb2-9fe7-c0d7-88b9-95266a13c336.htm "FlushDBAsync 方法 ")

[GetListLength 方法](../html/36330604-2557-4fc1-bdd7-fc3aff05e997.htm "GetListLength 方法 ")

[GetListLengthAsync 方法](../html/0e09431d-45b3-b539-a618-a147ae377f88.htm "GetListLengthAsync 方法 ")

[IncrementHashKey 方法](../html/6404f402-9882-5453-7df7-275f61748021.htm "IncrementHashKey 方法 ")

[IncrementHashKeyAsync 方法](../html/7a69e1a1-8689-3e23-8aa3-cdad33e7e3f8.htm "IncrementHashKeyAsync 方法 ")

[IncrementKey 方法](../html/04d46420-bedd-9016-075c-3bb1a6623583.htm "IncrementKey 方法 ")

[IncrementKeyAsync 方法](../html/95f88e06-114c-b5b3-ee2f-2e9f97f27b8c.htm "IncrementKeyAsync 方法 ")

[InitializationOnConnect 方法](../html/72414fc8-b43a-3377-acdc-c59ae0628e11.htm "InitializationOnConnect 方法 ")

[InitializationOnConnectAsync 方法](../html/9156d090-6b1d-7381-c396-81491b4b6756.htm "InitializationOnConnectAsync 方法 ")

[ListInsertAfter 方法](../html/11b6eb28-b2ee-57a5-c724-a99df9101c71.htm "ListInsertAfter 方法 ")

[ListInsertAfterAsync 方法](../html/9defd0d6-5681-4c26-0973-14ecca4fc821.htm "ListInsertAfterAsync 方法 ")

[ListInsertBefore 方法](../html/b19eed6b-2778-4b91-6f6b-e75cd9f1088e.htm "ListInsertBefore 方法 ")

[ListInsertBeforeAsync 方法](../html/0bf5390c-af9f-ecd6-1148-34a20abee73d.htm "ListInsertBeforeAsync 方法 ")

[ListLeftPop 方法](../html/8865ed9c-69dc-2f28-7f53-42e1c03d152b.htm "ListLeftPop 方法 ")

[ListLeftPopAsync 方法](../html/f4d7f464-73f7-ff5d-0834-e18b5a91c994.htm "ListLeftPopAsync 方法 ")

[ListLeftPush 方法](../html/232a2dbd-e2c6-5fbf-af96-5269a10df4c9.htm "ListLeftPush 方法 ")

[ListLeftPushAsync 方法](../html/cb91133d-ce36-f4f7-a9d3-81dbda7d6617.htm "ListLeftPushAsync 方法 ")

[ListLeftPushX 方法](../html/33093159-9864-a0dd-6a09-8bf134942031.htm "ListLeftPushX 方法 ")

[ListLeftPushXAsync 方法](../html/6bd2ab31-ff58-71f0-b68f-f36d7f53d845.htm "ListLeftPushXAsync 方法 ")

[ListRange 方法](../html/efe06eb0-b945-fe65-5e81-21888cdda04f.htm "ListRange 方法 ")

[ListRangeAsync 方法](../html/32e344f9-cdac-fc0e-2c50-bf2dd54b4f98.htm "ListRangeAsync 方法 ")

[ListRemoveElementMatch 方法](../html/0ab40e89-d265-0b48-05de-5c1a7b41d326.htm "ListRemoveElementMatch 方法 ")

[ListRemoveElementMatchAsync 方法](../html/5610de76-3198-5b5a-8f3d-b8fc937d7c83.htm "ListRemoveElementMatchAsync 方法 ")

[ListRightPop 方法](../html/97855b86-6a07-76e2-b887-1e830be705de.htm "ListRightPop 方法 ")

[ListRightPopAsync 方法](../html/fde8cee2-f393-a463-d561-74568545edc7.htm "ListRightPopAsync 方法 ")

[ListRightPopLeftPush 方法](../html/e0b58cc4-4d80-64d6-729c-0dcdaafc4701.htm "ListRightPopLeftPush 方法 ")

[ListRightPopLeftPushAsync 方法](../html/41069df2-0dad-6f01-895f-a91086a92fa9.htm "ListRightPopLeftPushAsync 方法 ")

[ListRightPush 方法](../html/d71ca156-74d3-536e-6925-315d948c0ace.htm "ListRightPush 方法 ")

[ListRightPushAsync 方法](../html/53144654-da3d-b367-b39c-0dbc68275e59.htm "ListRightPushAsync 方法 ")

[ListRightPushX 方法](../html/e6ccfd1a-25ee-5eff-d56c-0ff9eb017fbc.htm "ListRightPushX 方法 ")

[ListRightPushXAsync 方法](../html/3252db2a-56cf-1aa2-7716-f9eb2596ec62.htm "ListRightPushXAsync 方法 ")

[ListSet 方法](../html/5f2e7b8b-b0e7-87e6-dac7-81b572cf3fa3.htm "ListSet 方法 ")

[ListSetAsync 方法](../html/ca019ec2-eb32-41f8-a460-89c70f90ace9.htm "ListSetAsync 方法 ")

[ListTrim 方法](../html/acf21364-1692-f9da-d337-c460b12a012d.htm "ListTrim 方法 ")

[ListTrimAsync 方法](../html/9e86492d-205c-8427-9052-50ea3d8c1c64.htm "ListTrimAsync 方法 ")

[MoveKey 方法](../html/5aa50043-85b9-a245-606a-de2b69164218.htm "MoveKey 方法 ")

[MoveKeyAsync 方法](../html/a93a74bf-65bd-733b-e39a-06d705252d99.htm "MoveKeyAsync 方法 ")

[OperateLongNumberFromServer 方法](../html/d2cdcd44-a6fb-7aa7-8c9c-08ff43107f65.htm "OperateLongNumberFromServer 方法 ")

[OperateLongNumberFromServerAsync 方法](../html/0005132d-e8c6-d51b-69dc-c9f8feaf8a9d.htm "OperateLongNumberFromServerAsync 方法 ")

[OperateNumberFromServer 方法](../html/69aecc05-6901-f53a-da37-ebf85bae7053.htm "OperateNumberFromServer 方法 ")

[OperateNumberFromServerAsync 方法](../html/b53bd6fa-d215-55f7-bbac-2ee171ed7498.htm "OperateNumberFromServerAsync 方法 ")

[OperateStatusFromServer 方法](../html/d460c3bb-e484-3fcb-e206-d20dc7201d63.htm "OperateStatusFromServer 方法 ")

[OperateStatusFromServerAsync 方法](../html/e1a241a3-e200-747f-bb8b-5b7578a08999.htm "OperateStatusFromServerAsync 方法 ")

[OperateStringFromServer 方法](../html/822914af-cf62-1ccc-d22e-2df7dfac7294.htm "OperateStringFromServer 方法 ")

[OperateStringFromServerAsync 方法](../html/eb31e52d-5387-eb14-291e-22f0e64666a4.htm "OperateStringFromServerAsync 方法 ")

[OperateStringsFromServer 方法](../html/c85c1a0e-83f9-69c9-a9b9-8a947ca54583.htm "OperateStringsFromServer 方法 ")

[OperateStringsFromServerAsync 方法](../html/2db50b51-32b8-1e92-cb9a-f2c847437e50.htm "OperateStringsFromServerAsync 方法 ")

[PersistKey 方法](../html/c68a840b-ce64-3303-9339-a10f88bcbb77.htm "PersistKey 方法 ")

[PersistKeyAsync 方法](../html/a5d7ee67-9041-6455-7fc1-c6e96a442799.htm "PersistKeyAsync 方法 ")

[Ping 方法](../html/9cdd4e99-90a2-0ae8-141c-26f0f4a4d698.htm "Ping 方法 ")

[PingAsync 方法](../html/d6c8e625-edb2-7ff2-a54d-e852d49c73c8.htm "PingAsync 方法 ")

[Publish 方法](../html/74f20db5-aaaa-af2e-b4e6-5d192196b222.htm "Publish 方法 ")

[PublishAsync 方法](../html/61428af8-94a2-b2b0-419b-5c0d3df4bee5.htm "PublishAsync 方法 ")

[Read(T) 方法](../html/08365a89-e2d7-2aa5-3fd5-ab1a5d60c440.htm "Read(T) 方法 ")

[ReadAllKeys 方法](../html/53c48696-f991-f734-12b0-1ad306a84c6b.htm "ReadAllKeys 方法 ")

[ReadAllKeysAsync 方法](../html/84c1be8c-3286-0410-0bd1-bb0185ff197e.htm "ReadAllKeysAsync 方法 ")

[ReadAndWriteKey 方法](../html/a6908f0e-1991-5210-f6c6-1e8b06369e56.htm "ReadAndWriteKey 方法 ")

[ReadAndWriteKeyAsync 方法](../html/6686726f-93ba-6435-76b0-b501c31e637e.htm "ReadAndWriteKeyAsync 方法 ")

[ReadAsync(T) 方法](../html/43a59914-1502-1353-0bf3-20281ac92ca9.htm "ReadAsync(T) 方法 ")

[ReadCustomer 方法](../html/a4b7f43e-34ba-8aac-a705-a9bc8d4a5298.htm "ReadCustomer 方法 ")

[ReadCustomerAsync 方法](../html/bafea5ed-4293-ae60-c8e7-bb3d3481db44.htm "ReadCustomerAsync 方法 ")

[ReadFromCoreServer 方法](../html/a5985ec7-e909-5a40-b2b9-7160363a1156.htm "ReadFromCoreServer 方法 ")

[ReadFromCoreServerAsync 方法](../html/984a89de-20da-e7ed-e471-3e889955aa31.htm "ReadFromCoreServerAsync 方法 ")

[ReadHashKey 方法](../html/6e8fc531-b2d6-2b1e-4885-1b937745652f.htm "ReadHashKey 方法 ")

[ReadHashKeyAll 方法](../html/63828fbe-fbdf-772a-010b-1425b0584572.htm "ReadHashKeyAll 方法 ")

[ReadHashKeyAllAsync 方法](../html/31c54864-f7bc-b3ce-7e77-15f0c1dcc540.htm "ReadHashKeyAllAsync 方法 ")

[ReadHashKeyAsync 方法](../html/36c995b4-42a7-e264-8a03-3390e24376ec.htm "ReadHashKeyAsync 方法 ")

[ReadHashKeyLength 方法](../html/32d483da-49e0-a205-bdc4-a659711be1c0.htm "ReadHashKeyLength 方法 ")

[ReadHashKeyLengthAsync 方法](../html/5d0d4570-1f16-4277-cff7-a40607a530ab.htm "ReadHashKeyLengthAsync 方法 ")

[ReadHashKeys 方法](../html/420286f8-5876-ee1b-b328-0b6eba7cefac.htm "ReadHashKeys 方法 ")

[ReadHashKeysAsync 方法](../html/bbc8310c-029d-1ea5-b0e3-239b29ee8743.htm "ReadHashKeysAsync 方法 ")

[ReadHashValues 方法](../html/1d8b8e06-3d34-f4c6-68d0-20c14b309de1.htm "ReadHashValues 方法 ")

[ReadHashValuesAsync 方法](../html/403e7344-c665-def6-70d9-81e5b22a740b.htm "ReadHashValuesAsync 方法 ")

[ReadKey 方法](../html/93fa3dab-d22c-d028-00c3-e81391085e6f.htm "ReadKey 方法 ")

[ReadKeyAsync 方法](../html/48cee995-f908-c3a2-ab70-829d789f58a5.htm "ReadKeyAsync 方法 ")

[ReadKeyLength 方法](../html/7746357b-7e1b-d150-d019-722b752f5eeb.htm "ReadKeyLength 方法 ")

[ReadKeyLengthAsync 方法](../html/7dd1b305-77b1-920c-14c0-2fe38e970791.htm "ReadKeyLengthAsync 方法 ")

[ReadKeyRange 方法](../html/2e3e9582-cb24-44e6-d78b-7c672acdd4c8.htm "ReadKeyRange 方法 ")

[ReadKeyRangeAsync 方法](../html/407fdd5b-8b57-75d2-63be-93294311cb30.htm "ReadKeyRangeAsync 方法 ")

[ReadKeyTTL 方法](../html/679079c3-36e5-1618-3409-18459a74a0c4.htm "ReadKeyTTL 方法 ")

[ReadKeyTTLAsync 方法](../html/24887926-0fd5-b65f-d070-1d439a5273ef.htm "ReadKeyTTLAsync 方法 ")

[ReadKeyType 方法](../html/1f702770-7804-880a-038e-c002ed31c42f.htm "ReadKeyType 方法 ")

[ReadKeyTypeAsync 方法](../html/91c104ff-a646-76de-9bed-5e64237a57fc.htm "ReadKeyTypeAsync 方法 ")

[ReadListByIndex 方法](../html/43c7dae8-4c01-d2c9-b38b-0e7bef68b581.htm "ReadListByIndex 方法 ")

[ReadListByIndexAsync 方法](../html/1e1163ad-8161-9390-cff2-7919c8145c92.htm "ReadListByIndexAsync 方法 ")

[ReadRandomKey 方法](../html/0d358a8f-b9e0-c659-c218-f051aaedcbfe.htm "ReadRandomKey 方法 ")

[ReadRandomKeyAsync 方法](../html/b97f91d2-2ab6-a39b-86f7-71238cd44fac.htm "ReadRandomKeyAsync 方法 ")

[ReadServerTime 方法](../html/c6c5f1fa-4639-0a32-8a42-48d56985e8e6.htm "ReadServerTime 方法 ")

[ReadServerTimeAsync 方法](../html/53022c3b-766e-c31b-1b55-548ceb3b2f42.htm "ReadServerTimeAsync 方法 ")

[RenameKey 方法](../html/ecf533af-233c-16b8-f88d-cb232e115013.htm "RenameKey 方法 ")

[RenameKeyAsync 方法](../html/20073ebe-0c16-2f57-a96b-79246e01129d.htm "RenameKeyAsync 方法 ")

[Save 方法](../html/261a8daa-9890-95bc-f8e4-02b78eb86277.htm "Save 方法 ")

[SaveAsync 方法](../html/88c2043f-e9e6-9d95-4a2e-455452cf8adf.htm "SaveAsync 方法 ")

[SelectDB 方法](../html/584a38b8-96e0-736f-7a24-981a8e920574.htm "SelectDB 方法 ")

[SelectDBAsync 方法](../html/462ed998-9938-a261-bde4-579b4b9d3dc1.htm "SelectDBAsync 方法 ")

[SetAdd 方法](../html/6f52901b-4d7f-70ca-1bf1-2e7146604dba.htm "SetAdd 方法 ")

[SetAddAsync 方法](../html/a59f7c04-0b20-35a4-a431-53299cb829f0.htm "SetAddAsync 方法 ")

[SetCard 方法](../html/cc0a3d1d-019b-f0ba-90a5-1bab767d6f05.htm "SetCard 方法 ")

[SetCardAsync 方法](../html/0d0665a9-a174-d154-f95b-2bf10d885bc3.htm "SetCardAsync 方法 ")

[SetDiff 方法](../html/bcc43604-dc32-a27f-53a4-a89ab29916be.htm "SetDiff 方法 ")

[SetDiffAsync 方法](../html/b70eae65-3537-3ab7-4cf9-ae6c314884c9.htm "SetDiffAsync 方法 ")

[SetDiffStore 方法](../html/7c356901-c174-c81b-50b9-ce8d60258d1e.htm "SetDiffStore 方法 ")

[SetDiffStoreAsync 方法](../html/697ac488-c2e1-8c5a-5e5d-845f676c85fd.htm "SetDiffStoreAsync 方法 ")

[SetInter 方法](../html/fa3a123d-447d-ce8b-5670-e47956d427d1.htm "SetInter 方法 ")

[SetInterAsync 方法](../html/622854f7-5e07-092e-42eb-06cc63580467.htm "SetInterAsync 方法 ")

[SetInterStore 方法](../html/996ff5a7-0d0b-5b17-546b-eb006fa5773f.htm "SetInterStore 方法 ")

[SetInterStoreAsync 方法](../html/115a5c5d-2c7b-4c9a-daec-92477fb3044a.htm "SetInterStoreAsync 方法 ")

[SetIsMember 方法](../html/1c6d465b-353a-34f2-8f46-9d4037e58911.htm "SetIsMember 方法 ")

[SetIsMemberAsync 方法](../html/21419054-3273-2aab-d486-633ab8861039.htm "SetIsMemberAsync 方法 ")

[SetMembers 方法](../html/7a5b277f-1fe1-141f-346b-3722d97ac7ca.htm "SetMembers 方法 ")

[SetMembersAsync 方法](../html/ce86155f-4be2-3836-3594-1398d78183ad.htm "SetMembersAsync 方法 ")

[SetMove 方法](../html/8b76c71a-c755-708b-e1e7-61eda8b9456c.htm "SetMove 方法 ")

[SetMoveAsync 方法](../html/dc01893c-b09e-e134-49d3-76d0dbc308b9.htm "SetMoveAsync 方法 ")

[SetPop 方法](../html/8375117d-d7b5-1f33-cb59-80cc632b736d.htm "SetPop 方法 ")

[SetPopAsync 方法](../html/66039533-6b1e-ad66-e89b-ef8f7d46c47b.htm "SetPopAsync 方法 ")

[SetRandomMember 方法](../html/2769a6d4-3b4c-75b2-c57b-abe19bb2f8e9.htm "SetRandomMember 方法 ")

[SetRandomMemberAsync 方法](../html/fb4fda18-00e6-2545-d037-6585a9487334.htm "SetRandomMemberAsync 方法 ")

[SetRemove 方法](../html/a0ba101e-935e-aca6-da5e-bb6187b42f02.htm "SetRemove 方法 ")

[SetRemoveAsync 方法](../html/f80e5db7-218b-c4bf-0c57-ec2555f6ec4e.htm "SetRemoveAsync 方法 ")

[SetUnion 方法](../html/c255e92b-1a91-f3e4-49cd-79f39d7a1079.htm "SetUnion 方法 ")

[SetUnionAsync 方法](../html/184ffefd-f6f5-b5d8-e1d6-45ae6d2d653f.htm "SetUnionAsync 方法 ")

[SetUnionStore 方法](../html/56394825-4bd1-6422-4f9b-1520ed05f927.htm "SetUnionStore 方法 ")

[SetUnionStoreAsync 方法](../html/99689f9d-ebcc-cc16-19e3-4c46f3016084.htm "SetUnionStoreAsync 方法 ")

[SubscribeMessage 方法](../html/bbaad51e-d35a-fb1f-4bb5-c001d47693be.htm "SubscribeMessage 方法 ")

[ToString 方法](../html/6ff40e11-12f4-8759-45a0-594ca3e4fcdd.htm "ToString 方法 ")

[UnSubscribeMessage 方法](../html/592327f8-7e50-4d48-c229-ca17a039b6a3.htm "UnSubscribeMessage 方法 ")

[Write(T) 方法](../html/70a08f37-0f86-76f1-55a5-df1ddd5a7d16.htm "Write(T) 方法 ")

[WriteAndPublishKey 方法](../html/69029e5e-a78a-5c83-0300-b6dd057f9308.htm "WriteAndPublishKey 方法 ")

[WriteAndPublishKeyAsync 方法](../html/d44abf4b-02bf-cfa5-cbc8-0b36edab4521.htm "WriteAndPublishKeyAsync 方法 ")

[WriteAsync(T) 方法](../html/9d54e352-2bd0-2655-3e14-49a0213e6040.htm "WriteAsync(T) 方法 ")

[WriteExpireKey 方法](../html/816ca8c0-5e3a-4609-acc3-6f227ac2dfde.htm "WriteExpireKey 方法 ")

[WriteExpireKeyAsync 方法](../html/f1dce55a-5681-32f6-2c7c-cc664465490c.htm "WriteExpireKeyAsync 方法 ")

[WriteHashKey 方法](../html/a7c80b1f-ec11-c8f4-dae2-0fa0acada685.htm "WriteHashKey 方法 ")

[WriteHashKeyAsync 方法](../html/20122681-134d-f220-683c-11c6c9201ac9.htm "WriteHashKeyAsync 方法 ")

[WriteHashKeyNx 方法](../html/819ddbb0-74d4-4667-6e55-372fac166c68.htm "WriteHashKeyNx 方法 ")

[WriteHashKeyNxAsync 方法](../html/b4571eb9-be29-fbb8-691b-f3a93d9245ec.htm "WriteHashKeyNxAsync 方法 ")

[WriteKey 方法](../html/12919427-9ab8-d765-27cc-0c3ee2c481d2.htm "WriteKey 方法 ")

[WriteKeyAsync 方法](../html/b136bcbb-3813-f660-54f1-8c5ccefb5ddf.htm "WriteKeyAsync 方法 ")

[WriteKeyIfNotExists 方法](../html/9d48312b-2ce1-077f-bf47-3e6b5a328c45.htm "WriteKeyIfNotExists 方法 ")

[WriteKeyIfNotExistsAsync 方法](../html/e6844e26-ba17-a0a8-0ff2-bb1a1e8bce8e.htm "WriteKeyIfNotExistsAsync 方法 ")

[WriteKeyRange 方法](../html/c1f98cc9-f791-4a35-bf67-ea29e8930d32.htm "WriteKeyRange 方法 ")

[WriteKeyRangeAsync 方法](../html/cdfaffec-518b-2754-5a10-3c5b68ee1767.htm "WriteKeyRangeAsync 方法 ")

[ZSetAdd 方法](../html/9c15cbbe-df53-5280-c9e4-b9d020222e0c.htm "ZSetAdd 方法 ")

[ZSetAddAsync 方法](../html/5318d2e6-d635-ed62-18e5-6573a530ab3e.htm "ZSetAddAsync 方法 ")

[ZSetCard 方法](../html/aa928601-c5f6-63fb-6bd8-e2e383796043.htm "ZSetCard 方法 ")

[ZSetCardAsync 方法](../html/b760ff15-4c7e-510f-1927-e17e6310e142.htm "ZSetCardAsync 方法 ")

[ZSetCount 方法](../html/9872ddc0-8b8d-d9d6-4d08-64001fdefd29.htm "ZSetCount 方法 ")

[ZSetCountAsync 方法](../html/b294e1ee-27c7-7326-eb4b-8669040aecf0.htm "ZSetCountAsync 方法 ")

[ZSetIncreaseBy 方法](../html/c84e0c14-15cd-451a-dec4-2a112a73772e.htm "ZSetIncreaseBy 方法 ")

[ZSetIncreaseByAsync 方法](../html/9805ac6b-aa05-2152-7dbb-7f8077ea462b.htm "ZSetIncreaseByAsync 方法 ")

[ZSetRange 方法](../html/759e2523-e2ad-0dc1-a904-354031a8ae50.htm "ZSetRange 方法 ")

[ZSetRangeAsync 方法](../html/7724eb62-2653-65a7-779d-391d630c2faf.htm "ZSetRangeAsync 方法 ")

[ZSetRangeByScore 方法](../html/acab27d0-6b6c-1f25-00ad-fafd4ac068fa.htm "ZSetRangeByScore 方法 ")

[ZSetRangeByScoreAsync 方法](../html/d1e0aa3c-d5f9-caeb-4017-699f956e661d.htm "ZSetRangeByScoreAsync 方法 ")

[ZSetRank 方法](../html/7957411f-32c1-4e07-a4f2-9d1948c63a2b.htm "ZSetRank 方法 ")

[ZSetRankAsync 方法](../html/1181fd8c-02aa-9d53-cbeb-35d542bed241.htm "ZSetRankAsync 方法 ")

[ZSetRemove 方法](../html/a76625ea-7a7b-161d-a039-02f4f5a8e997.htm "ZSetRemove 方法 ")

[ZSetRemoveAsync 方法](../html/2f1b3de6-ad74-e2b7-3350-79d9afe13c9d.htm "ZSetRemoveAsync 方法 ")

[ZSetRemoveRangeByRank 方法](../html/cb81dfc2-6e30-0aca-8bcb-348cdc32bc35.htm "ZSetRemoveRangeByRank 方法 ")

[ZSetRemoveRangeByRankAsync 方法](../html/9bde4b5f-b9bd-e01b-feb4-82c9b532dee4.htm "ZSetRemoveRangeByRankAsync 方法 ")

[ZSetRemoveRangeByScore 方法](../html/ac7be899-fe4e-3f12-d60f-81a038a927c7.htm "ZSetRemoveRangeByScore 方法 ")

[ZSetRemoveRangeByScoreAsync 方法](../html/6ccee6bf-8635-14bf-2e7f-d1ac351c70c9.htm "ZSetRemoveRangeByScoreAsync 方法 ")

[ZSetReverseRange 方法](../html/ef6defc7-62a2-c8a4-4025-32f6b3a82ae4.htm "ZSetReverseRange 方法 ")

[ZSetReverseRangeAsync 方法](../html/b40d37b7-1118-6aec-c7da-85ec54b457d0.htm "ZSetReverseRangeAsync 方法 ")

[ZSetReverseRangeByScore 方法](../html/3c1a26a8-f8e5-a835-7051-94b4d9c851ab.htm "ZSetReverseRangeByScore 方法 ")

[ZSetReverseRangeByScoreAsync 方法](../html/3a3cfbff-cee7-e1f3-d9c1-e7528db4e998.htm "ZSetReverseRangeByScoreAsync 方法 ")

[ZSetReverseRank 方法](../html/d524195b-6b02-a8e8-330a-3454a6ed21cd.htm "ZSetReverseRank 方法 ")

[ZSetReverseRankAsync 方法](../html/c010781c-2851-7158-f953-1779929e3242.htm "ZSetReverseRankAsync 方法 ")

[ZSetScore 方法](../html/a98ee7e4-f361-f403-975b-a197e0ca0268.htm "ZSetScore 方法 ")

[ZSetScoreAsync 方法](../html/fd24c433-944e-2ecb-4bcd-266bfc6825d4.htm "ZSetScoreAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| RedisClientDBSizeAsync 方法 |

返回当前数据库的 key 的数量。

**命名空间：**
 [HslCommunication.Enthernet.Redis](cc44ce00-3781-c690-b233-24000870bf42.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public Task<OperateResult<long>> DBSizeAsync()
```

```
Public Function DBSizeAsync As Task(Of OperateResult(Of Long))
```

```
public:
Task<OperateResult<long long>^>^ DBSizeAsync()
```

```
member DBSizeAsync : unit -> Task<OperateResult<int64>> 
```

#### 返回值

类型：Task[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Int64  
当前数据库的 key 的数量。

![](../icons/SectionExpanded.png)参见

#### 引用

[RedisClient 类](e33d9088-6df0-32e9-33f1-426e944c3b7f.htm)

[HslCommunication.Enthernet.Redis 命名空间](cc44ce00-3781-c690-b233-24000870bf42.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## DecrementKey 方法 

[原文連結](http://api.hslcommunication.cn/html/e9fa6c94-43ca-8291-79ca-6e04ca9289e7.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Enthernet.Redis](../html/cc44ce00-3781-c690-b233-24000870bf42.htm "HslCommunication.Enthernet.Redis")

[RedisClient 类](../html/e33d9088-6df0-32e9-33f1-426e944c3b7f.htm "RedisClient 类")

[RedisClient 方法](../html/e41be4a5-1490-ffbd-f241-938f2d6c3f07.htm "RedisClient 方法")

[DecrementKey 方法](../html/e9fa6c94-43ca-8291-79ca-6e04ca9289e7.htm "DecrementKey 方法 ")

[DecrementKey 方法 (String)](../html/0f2a8a0d-9e7f-c6ab-2081-bba6606c37a9.htm "DecrementKey 方法 (String)")

[DecrementKey 方法 (String, Int64)](../html/0bb5b243-e3e3-064c-1c76-219ad3e2768a.htm "DecrementKey 方法 (String, Int64)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| RedisClientDecrementKey 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [DecrementKey(String)](0f2a8a0d-9e7f-c6ab-2081-bba6606c37a9.htm) | 将 key 中储存的数字值减一。如果 key 不存在，那么 key 的值会先被初始化为 0 ，然后再执行 DECR 操作。 如果值包含错误的类型，或字符串类型的值不能表示为数字，那么返回一个错误。 本操作的值限制在 64 位(bit)有符号数字表示之内。 返回执行 DECR 命令之后 key 的值。 |
| 公共方法 | [DecrementKey(String, Int64)](0bb5b243-e3e3-064c-1c76-219ad3e2768a.htm) | 将 key 所储存的值减去减量 decrement 。如果 key 不存在，那么 key 的值会先被初始化为 0 ，然后再执行 DECR 操作。 如果值包含错误的类型，或字符串类型的值不能表示为数字，那么返回一个错误。 本操作的值限制在 64 位(bit)有符号数字表示之内。 返回减去 decrement 之后， key 的值。 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[RedisClient 类](e33d9088-6df0-32e9-33f1-426e944c3b7f.htm)

[HslCommunication.Enthernet.Redis 命名空间](cc44ce00-3781-c690-b233-24000870bf42.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## DecrementKey 方法 (String)

[原文連結](http://api.hslcommunication.cn/html/0f2a8a0d-9e7f-c6ab-2081-bba6606c37a9.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Enthernet.Redis](../html/cc44ce00-3781-c690-b233-24000870bf42.htm "HslCommunication.Enthernet.Redis")

[RedisClient 类](../html/e33d9088-6df0-32e9-33f1-426e944c3b7f.htm "RedisClient 类")

[RedisClient 方法](../html/e41be4a5-1490-ffbd-f241-938f2d6c3f07.htm "RedisClient 方法")

[DecrementKey 方法](../html/e9fa6c94-43ca-8291-79ca-6e04ca9289e7.htm "DecrementKey 方法 ")

[DecrementKey 方法 (String)](../html/0f2a8a0d-9e7f-c6ab-2081-bba6606c37a9.htm "DecrementKey 方法 (String)")

[DecrementKey 方法 (String, Int64)](../html/0bb5b243-e3e3-064c-1c76-219ad3e2768a.htm "DecrementKey 方法 (String, Int64)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| RedisClientDecrementKey 方法 (String) |

将 key 中储存的数字值减一。如果 key 不存在，那么 key 的值会先被初始化为 0 ，然后再执行 DECR 操作。
如果值包含错误的类型，或字符串类型的值不能表示为数字，那么返回一个错误。
本操作的值限制在 64 位(bit)有符号数字表示之内。
返回执行 DECR 命令之后 key 的值。

**命名空间：**
 [HslCommunication.Enthernet.Redis](cc44ce00-3781-c690-b233-24000870bf42.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public OperateResult<long> DecrementKey(
	string key
)
```

```
Public Function DecrementKey ( 
	key As String
) As OperateResult(Of Long)
```

```
public:
OperateResult<long long>^ DecrementKey(
	String^ key
)
```

```
member DecrementKey : 
        key : string -> OperateResult<int64> 
```

#### 参数

key
:   类型：SystemString  
    关键字

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Int64  
执行 DECR 命令之后 key 的值。

![](../icons/SectionExpanded.png)参见

#### 引用

[RedisClient 类](e33d9088-6df0-32e9-33f1-426e944c3b7f.htm)

[DecrementKey 重载](e9fa6c94-43ca-8291-79ca-6e04ca9289e7.htm)

[HslCommunication.Enthernet.Redis 命名空间](cc44ce00-3781-c690-b233-24000870bf42.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## DecrementKey 方法 (String, Int64)

[原文連結](http://api.hslcommunication.cn/html/0bb5b243-e3e3-064c-1c76-219ad3e2768a.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Enthernet.Redis](../html/cc44ce00-3781-c690-b233-24000870bf42.htm "HslCommunication.Enthernet.Redis")

[RedisClient 类](../html/e33d9088-6df0-32e9-33f1-426e944c3b7f.htm "RedisClient 类")

[RedisClient 方法](../html/e41be4a5-1490-ffbd-f241-938f2d6c3f07.htm "RedisClient 方法")

[DecrementKey 方法](../html/e9fa6c94-43ca-8291-79ca-6e04ca9289e7.htm "DecrementKey 方法 ")

[DecrementKey 方法 (String)](../html/0f2a8a0d-9e7f-c6ab-2081-bba6606c37a9.htm "DecrementKey 方法 (String)")

[DecrementKey 方法 (String, Int64)](../html/0bb5b243-e3e3-064c-1c76-219ad3e2768a.htm "DecrementKey 方法 (String, Int64)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| RedisClientDecrementKey 方法 (String, Int64) |

将 key 所储存的值减去减量 decrement 。如果 key 不存在，那么 key 的值会先被初始化为 0 ，然后再执行 DECR 操作。
如果值包含错误的类型，或字符串类型的值不能表示为数字，那么返回一个错误。
本操作的值限制在 64 位(bit)有符号数字表示之内。
返回减去 decrement 之后， key 的值。

**命名空间：**
 [HslCommunication.Enthernet.Redis](cc44ce00-3781-c690-b233-24000870bf42.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public OperateResult<long> DecrementKey(
	string key,
	long value
)
```

```
Public Function DecrementKey ( 
	key As String,
	value As Long
) As OperateResult(Of Long)
```

```
public:
OperateResult<long long>^ DecrementKey(
	String^ key, 
	long long value
)
```

```
member DecrementKey : 
        key : string * 
        value : int64 -> OperateResult<int64> 
```

#### 参数

key
:   类型：SystemString  
    关键字

value
:   类型：SystemInt64  
    操作的值

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Int64  
返回减去 decrement 之后， key 的值。

![](../icons/SectionExpanded.png)参见

#### 引用

[RedisClient 类](e33d9088-6df0-32e9-33f1-426e944c3b7f.htm)

[DecrementKey 重载](e9fa6c94-43ca-8291-79ca-6e04ca9289e7.htm)

[HslCommunication.Enthernet.Redis 命名空间](cc44ce00-3781-c690-b233-24000870bf42.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## DecrementKeyAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/6837e351-b8ef-a90e-c5e2-73b31621f837.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Enthernet.Redis](../html/cc44ce00-3781-c690-b233-24000870bf42.htm "HslCommunication.Enthernet.Redis")

[RedisClient 类](../html/e33d9088-6df0-32e9-33f1-426e944c3b7f.htm "RedisClient 类")

[RedisClient 方法](../html/e41be4a5-1490-ffbd-f241-938f2d6c3f07.htm "RedisClient 方法")

[DecrementKeyAsync 方法](../html/6837e351-b8ef-a90e-c5e2-73b31621f837.htm "DecrementKeyAsync 方法 ")

[DecrementKeyAsync 方法 (String)](../html/d0ade959-e1b8-cd2c-6a58-7104767acb23.htm "DecrementKeyAsync 方法 (String)")

[DecrementKeyAsync 方法 (String, Int64)](../html/743ce5f1-4a4e-314d-424d-a739f9e2384a.htm "DecrementKeyAsync 方法 (String, Int64)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| RedisClientDecrementKeyAsync 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [DecrementKeyAsync(String)](d0ade959-e1b8-cd2c-6a58-7104767acb23.htm) | 将 key 中储存的数字值减一。如果 key 不存在，那么 key 的值会先被初始化为 0 ，然后再执行 DECR 操作。 如果值包含错误的类型，或字符串类型的值不能表示为数字，那么返回一个错误。 本操作的值限制在 64 位(bit)有符号数字表示之内。 返回执行 DECR 命令之后 key 的值。 |
| 公共方法 | [DecrementKeyAsync(String, Int64)](743ce5f1-4a4e-314d-424d-a739f9e2384a.htm) | 将 key 所储存的值减去减量 decrement 。如果 key 不存在，那么 key 的值会先被初始化为 0 ，然后再执行 DECR 操作。 如果值包含错误的类型，或字符串类型的值不能表示为数字，那么返回一个错误。 本操作的值限制在 64 位(bit)有符号数字表示之内。 返回减去 decrement 之后， key 的值。 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[RedisClient 类](e33d9088-6df0-32e9-33f1-426e944c3b7f.htm)

[HslCommunication.Enthernet.Redis 命名空间](cc44ce00-3781-c690-b233-24000870bf42.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## DecrementKeyAsync 方法 (String)

[原文連結](http://api.hslcommunication.cn/html/d0ade959-e1b8-cd2c-6a58-7104767acb23.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Enthernet.Redis](../html/cc44ce00-3781-c690-b233-24000870bf42.htm "HslCommunication.Enthernet.Redis")

[RedisClient 类](../html/e33d9088-6df0-32e9-33f1-426e944c3b7f.htm "RedisClient 类")

[RedisClient 方法](../html/e41be4a5-1490-ffbd-f241-938f2d6c3f07.htm "RedisClient 方法")

[DecrementKeyAsync 方法](../html/6837e351-b8ef-a90e-c5e2-73b31621f837.htm "DecrementKeyAsync 方法 ")

[DecrementKeyAsync 方法 (String)](../html/d0ade959-e1b8-cd2c-6a58-7104767acb23.htm "DecrementKeyAsync 方法 (String)")

[DecrementKeyAsync 方法 (String, Int64)](../html/743ce5f1-4a4e-314d-424d-a739f9e2384a.htm "DecrementKeyAsync 方法 (String, Int64)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| RedisClientDecrementKeyAsync 方法 (String) |

将 key 中储存的数字值减一。如果 key 不存在，那么 key 的值会先被初始化为 0 ，然后再执行 DECR 操作。
如果值包含错误的类型，或字符串类型的值不能表示为数字，那么返回一个错误。
本操作的值限制在 64 位(bit)有符号数字表示之内。
返回执行 DECR 命令之后 key 的值。

**命名空间：**
 [HslCommunication.Enthernet.Redis](cc44ce00-3781-c690-b233-24000870bf42.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public Task<OperateResult<long>> DecrementKeyAsync(
	string key
)
```

```
Public Function DecrementKeyAsync ( 
	key As String
) As Task(Of OperateResult(Of Long))
```

```
public:
Task<OperateResult<long long>^>^ DecrementKeyAsync(
	String^ key
)
```

```
member DecrementKeyAsync : 
        key : string -> Task<OperateResult<int64>> 
```

#### 参数

key
:   类型：SystemString  
    关键字

#### 返回值

类型：Task[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Int64  
执行 DECR 命令之后 key 的值。

![](../icons/SectionExpanded.png)参见

#### 引用

[RedisClient 类](e33d9088-6df0-32e9-33f1-426e944c3b7f.htm)

[DecrementKeyAsync 重载](6837e351-b8ef-a90e-c5e2-73b31621f837.htm)

[HslCommunication.Enthernet.Redis 命名空间](cc44ce00-3781-c690-b233-24000870bf42.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## DecrementKeyAsync 方法 (String, Int64)

[原文連結](http://api.hslcommunication.cn/html/743ce5f1-4a4e-314d-424d-a739f9e2384a.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Enthernet.Redis](../html/cc44ce00-3781-c690-b233-24000870bf42.htm "HslCommunication.Enthernet.Redis")

[RedisClient 类](../html/e33d9088-6df0-32e9-33f1-426e944c3b7f.htm "RedisClient 类")

[RedisClient 方法](../html/e41be4a5-1490-ffbd-f241-938f2d6c3f07.htm "RedisClient 方法")

[DecrementKeyAsync 方法](../html/6837e351-b8ef-a90e-c5e2-73b31621f837.htm "DecrementKeyAsync 方法 ")

[DecrementKeyAsync 方法 (String)](../html/d0ade959-e1b8-cd2c-6a58-7104767acb23.htm "DecrementKeyAsync 方法 (String)")

[DecrementKeyAsync 方法 (String, Int64)](../html/743ce5f1-4a4e-314d-424d-a739f9e2384a.htm "DecrementKeyAsync 方法 (String, Int64)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| RedisClientDecrementKeyAsync 方法 (String, Int64) |

将 key 所储存的值减去减量 decrement 。如果 key 不存在，那么 key 的值会先被初始化为 0 ，然后再执行 DECR 操作。
如果值包含错误的类型，或字符串类型的值不能表示为数字，那么返回一个错误。
本操作的值限制在 64 位(bit)有符号数字表示之内。
返回减去 decrement 之后， key 的值。

**命名空间：**
 [HslCommunication.Enthernet.Redis](cc44ce00-3781-c690-b233-24000870bf42.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public Task<OperateResult<long>> DecrementKeyAsync(
	string key,
	long value
)
```

```
Public Function DecrementKeyAsync ( 
	key As String,
	value As Long
) As Task(Of OperateResult(Of Long))
```

```
public:
Task<OperateResult<long long>^>^ DecrementKeyAsync(
	String^ key, 
	long long value
)
```

```
member DecrementKeyAsync : 
        key : string * 
        value : int64 -> Task<OperateResult<int64>> 
```

#### 参数

key
:   类型：SystemString  
    关键字

value
:   类型：SystemInt64  
    操作的值

#### 返回值

类型：Task[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Int64  
返回减去 decrement 之后， key 的值。

![](../icons/SectionExpanded.png)参见

#### 引用

[RedisClient 类](e33d9088-6df0-32e9-33f1-426e944c3b7f.htm)

[DecrementKeyAsync 重载](6837e351-b8ef-a90e-c5e2-73b31621f837.htm)

[HslCommunication.Enthernet.Redis 命名空间](cc44ce00-3781-c690-b233-24000870bf42.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## DeleteHashKey 方法 

[原文連結](http://api.hslcommunication.cn/html/19e3ebc3-e6fc-71cf-e80d-cb73e3011170.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Enthernet.Redis](../html/cc44ce00-3781-c690-b233-24000870bf42.htm "HslCommunication.Enthernet.Redis")

[RedisClient 类](../html/e33d9088-6df0-32e9-33f1-426e944c3b7f.htm "RedisClient 类")

[RedisClient 方法](../html/e41be4a5-1490-ffbd-f241-938f2d6c3f07.htm "RedisClient 方法")

[DeleteHashKey 方法](../html/19e3ebc3-e6fc-71cf-e80d-cb73e3011170.htm "DeleteHashKey 方法 ")

[DeleteHashKey 方法 (String, String)](../html/3f085ae2-ff3d-3c08-98af-5796aba0d975.htm "DeleteHashKey 方法 (String, String)")

[DeleteHashKey 方法 (String, String[])](../html/eeed0572-14a6-8f5f-9e59-0df1636a2cb0.htm "DeleteHashKey 方法 (String, String[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| RedisClientDeleteHashKey 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [DeleteHashKey(String, String)](3f085ae2-ff3d-3c08-98af-5796aba0d975.htm) | 删除哈希表 key 中的一个或多个指定域，不存在的域将被忽略。 |
| 公共方法 | [DeleteHashKey(String, String)](eeed0572-14a6-8f5f-9e59-0df1636a2cb0.htm) | 删除哈希表 key 中的一个或多个指定域，不存在的域将被忽略。返回被成功移除的域的数量，不包括被忽略的域。 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[RedisClient 类](e33d9088-6df0-32e9-33f1-426e944c3b7f.htm)

[HslCommunication.Enthernet.Redis 命名空间](cc44ce00-3781-c690-b233-24000870bf42.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)