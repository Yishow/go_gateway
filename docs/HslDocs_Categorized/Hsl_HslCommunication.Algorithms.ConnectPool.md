# HslCommunication - HslCommunication.Algorithms.ConnectPool

> 分類頁數: 30



---
## HslCommunication.Algorithms.ConnectPool

[原文連結](http://api.hslcommunication.cn/html/9ec7e0cd-f203-a971-db9f-ea1f53b0b4b3.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Algorithms.ConnectPool](../html/9ec7e0cd-f203-a971-db9f-ea1f53b0b4b3.htm "HslCommunication.Algorithms.ConnectPool")

[ConnectPool(TConnector) 类](../html/8697bcd0-0fa1-9274-3029-3ee6a893fd56.htm "ConnectPool(TConnector) 类")

[DeviceCommConnector 类](../html/9714842a-ab7e-6320-7092-b2a8ca7d7c0e.htm "DeviceCommConnector 类")

[IConnector 接口](../html/ede7f47f-b60e-e320-8bc2-389883e77083.htm "IConnector 接口")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslCommunication.Algorithms.ConnectPool 命名空间 |

[缺少 "N:HslCommunication.Algorithms.ConnectPool" 的 <summary> 文档]

![](../icons/SectionExpanded.png)类

|  | 类 | 说明 |
| --- | --- | --- |
| 公共类代码示例 | [ConnectPoolTConnector](8697bcd0-0fa1-9274-3029-3ee6a893fd56.htm) | 一个连接池管理器，负责维护多个可用的连接，并且自动清理，扩容，用于快速读写服务器或是PLC时使用。  A connection pool manager is responsible for maintaining multiple available connections, and automatically cleans up, expands, and is used to quickly read and write servers or PLCs. |
| 公共类 | [DeviceCommConnector](9714842a-ab7e-6320-7092-b2a8ca7d7c0e.htm) | 基于设备通信的连接池信息  Connection pool information based on device communication |

![](../icons/SectionExpanded.png)接口

|  | 接口 | 说明 |
| --- | --- | --- |
| 公共接口代码示例 | [IConnector](ede7f47f-b60e-e320-8bc2-389883e77083.htm) | 连接池的接口，连接池的管理对象必须实现此接口 |

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ConnectPool(TConnector) 构造函数 

[原文連結](http://api.hslcommunication.cn/html/29586bec-621f-6470-c2ca-3808221c70f2.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Algorithms.ConnectPool](../html/9ec7e0cd-f203-a971-db9f-ea1f53b0b4b3.htm "HslCommunication.Algorithms.ConnectPool")

[ConnectPool(TConnector) 类](../html/8697bcd0-0fa1-9274-3029-3ee6a893fd56.htm "ConnectPool(TConnector) 类")

[ConnectPool(TConnector) 构造函数](../html/29586bec-621f-6470-c2ca-3808221c70f2.htm "ConnectPool(TConnector) 构造函数 ")

[ConnectPool(TConnector) 属性](../html/73798dd3-7274-847c-6aec-dae2533dacef.htm "ConnectPool(TConnector) 属性")

[ConnectPool(TConnector) 方法](../html/c22ef39f-c4e9-7389-66d5-37459c7522f6.htm "ConnectPool(TConnector) 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ConnectPoolTConnector 构造函数 |

实例化一个连接池对象，需要指定如果创建新实例的方法  
To instantiate a connection pool object, you need to specify how to create a new instance

**命名空间：**
 [HslCommunication.Algorithms.ConnectPool](9ec7e0cd-f203-a971-db9f-ea1f53b0b4b3.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public ConnectPool(
	Func<TConnector> createConnector
)
```

```
Public Sub New ( 
	createConnector As Func(Of TConnector)
)
```

```
public:
ConnectPool(
	Func<TConnector>^ createConnector
)
```

```
new : 
        createConnector : Func<'TConnector> -> ConnectPool
```

#### 参数

createConnector
:   类型：SystemFunc[TConnector](8697bcd0-0fa1-9274-3029-3ee6a893fd56.htm)  
    创建连接对象的委托

![](../icons/SectionExpanded.png)参见

#### 引用

[ConnectPoolTConnector 类](8697bcd0-0fa1-9274-3029-3ee6a893fd56.htm)

[HslCommunication.Algorithms.ConnectPool 命名空间](9ec7e0cd-f203-a971-db9f-ea1f53b0b4b3.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ConnectPool(TConnector) 属性

[原文連結](http://api.hslcommunication.cn/html/73798dd3-7274-847c-6aec-dae2533dacef.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Algorithms.ConnectPool](../html/9ec7e0cd-f203-a971-db9f-ea1f53b0b4b3.htm "HslCommunication.Algorithms.ConnectPool")

[ConnectPool(TConnector) 类](../html/8697bcd0-0fa1-9274-3029-3ee6a893fd56.htm "ConnectPool(TConnector) 类")

[ConnectPool(TConnector) 属性](../html/73798dd3-7274-847c-6aec-dae2533dacef.htm "ConnectPool(TConnector) 属性")

[ConectionExpireTime 属性](../html/9f233c68-58fa-91e5-8a28-09b6cd0c10d5.htm "ConectionExpireTime 属性 ")

[MaxConnector 属性](../html/9dc834aa-57d2-ee27-248f-8118aded699a.htm "MaxConnector 属性 ")

[UseConnectorMax 属性](../html/feef2662-bb80-a0f0-9db8-7f8fd2c34197.htm "UseConnectorMax 属性 ")

[UsedConnector 属性](../html/f96d08a2-465a-8076-82d7-222c0b3d225f.htm "UsedConnector 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ConnectPoolTConnector 属性 |

[ConnectPoolTConnector](8697bcd0-0fa1-9274-3029-3ee6a893fd56.htm) 泛型公开以下成员。

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [ConectionExpireTime](9f233c68-58fa-91e5-8a28-09b6cd0c10d5.htm) | 获取或设置当前连接过期的时间，单位秒，默认30秒，也就是说，当前的连接在设置的时间段内未被使用，就进行释放连接，减少内存消耗。  Get or set the expiration time of the current connection, in seconds, the default is 30 seconds, that is, if the current connection is not used within the set time period, the connection will be released to reduce memory consumption. |
| 公共属性 | [MaxConnector](9dc834aa-57d2-ee27-248f-8118aded699a.htm) | 获取或设置最大的连接数，当实际的连接数超过最大的连接数的时候，就会进行阻塞，直到有新的连接对象为止。  Get or set the maximum number of connections. When the actual number of connections exceeds the maximum number of connections, it will block until there is a new connection object. |
| 公共属性 | [UseConnectorMax](feef2662-bb80-a0f0-9db8-7f8fd2c34197.htm) | 当前已经使用的连接数的峰值，可以用来衡量当前系统的适用的连接池上限。  The current peak value of the number of connections used can be used to measure the upper limit of the applicable connection pool of the current system. |
| 公共属性 | [UsedConnector](f96d08a2-465a-8076-82d7-222c0b3d225f.htm) | 当前已经使用的连接数，会根据使用的频繁程度进行动态的变化。  The number of currently used connections will dynamically change according to the frequency of use. |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[ConnectPoolTConnector 类](8697bcd0-0fa1-9274-3029-3ee6a893fd56.htm)

[HslCommunication.Algorithms.ConnectPool 命名空间](9ec7e0cd-f203-a971-db9f-ea1f53b0b4b3.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ConectionExpireTime 属性 

[原文連結](http://api.hslcommunication.cn/html/9f233c68-58fa-91e5-8a28-09b6cd0c10d5.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Algorithms.ConnectPool](../html/9ec7e0cd-f203-a971-db9f-ea1f53b0b4b3.htm "HslCommunication.Algorithms.ConnectPool")

[ConnectPool(TConnector) 类](../html/8697bcd0-0fa1-9274-3029-3ee6a893fd56.htm "ConnectPool(TConnector) 类")

[ConnectPool(TConnector) 属性](../html/73798dd3-7274-847c-6aec-dae2533dacef.htm "ConnectPool(TConnector) 属性")

[ConectionExpireTime 属性](../html/9f233c68-58fa-91e5-8a28-09b6cd0c10d5.htm "ConectionExpireTime 属性 ")

[MaxConnector 属性](../html/9dc834aa-57d2-ee27-248f-8118aded699a.htm "MaxConnector 属性 ")

[UseConnectorMax 属性](../html/feef2662-bb80-a0f0-9db8-7f8fd2c34197.htm "UseConnectorMax 属性 ")

[UsedConnector 属性](../html/f96d08a2-465a-8076-82d7-222c0b3d225f.htm "UsedConnector 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ConnectPoolTConnectorConectionExpireTime 属性 |

获取或设置当前连接过期的时间，单位秒，默认30秒，也就是说，当前的连接在设置的时间段内未被使用，就进行释放连接，减少内存消耗。  
Get or set the expiration time of the current connection, in seconds, the default is 30 seconds, that is,
if the current connection is not used within the set time period, the connection will be released to reduce memory consumption.

**命名空间：**
 [HslCommunication.Algorithms.ConnectPool](9ec7e0cd-f203-a971-db9f-ea1f53b0b4b3.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public int ConectionExpireTime { get; set; }
```

```
Public Property ConectionExpireTime As Integer
	Get
	Set
```

```
public:
property int ConectionExpireTime {
	int get ();
	void set (int value);
}
```

```
member ConectionExpireTime : int with get, set
```

#### 属性值

类型：Int32

![](../icons/SectionExpanded.png)参见

#### 引用

[ConnectPoolTConnector 类](8697bcd0-0fa1-9274-3029-3ee6a893fd56.htm)

[HslCommunication.Algorithms.ConnectPool 命名空间](9ec7e0cd-f203-a971-db9f-ea1f53b0b4b3.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## MaxConnector 属性 

[原文連結](http://api.hslcommunication.cn/html/9dc834aa-57d2-ee27-248f-8118aded699a.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Algorithms.ConnectPool](../html/9ec7e0cd-f203-a971-db9f-ea1f53b0b4b3.htm "HslCommunication.Algorithms.ConnectPool")

[ConnectPool(TConnector) 类](../html/8697bcd0-0fa1-9274-3029-3ee6a893fd56.htm "ConnectPool(TConnector) 类")

[ConnectPool(TConnector) 属性](../html/73798dd3-7274-847c-6aec-dae2533dacef.htm "ConnectPool(TConnector) 属性")

[ConectionExpireTime 属性](../html/9f233c68-58fa-91e5-8a28-09b6cd0c10d5.htm "ConectionExpireTime 属性 ")

[MaxConnector 属性](../html/9dc834aa-57d2-ee27-248f-8118aded699a.htm "MaxConnector 属性 ")

[UseConnectorMax 属性](../html/feef2662-bb80-a0f0-9db8-7f8fd2c34197.htm "UseConnectorMax 属性 ")

[UsedConnector 属性](../html/f96d08a2-465a-8076-82d7-222c0b3d225f.htm "UsedConnector 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ConnectPoolTConnectorMaxConnector 属性 |

获取或设置最大的连接数，当实际的连接数超过最大的连接数的时候，就会进行阻塞，直到有新的连接对象为止。  
Get or set the maximum number of connections. When the actual number of connections exceeds the maximum number of connections,
it will block until there is a new connection object.

**命名空间：**
 [HslCommunication.Algorithms.ConnectPool](9ec7e0cd-f203-a971-db9f-ea1f53b0b4b3.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public int MaxConnector { get; set; }
```

```
Public Property MaxConnector As Integer
	Get
	Set
```

```
public:
property int MaxConnector {
	int get ();
	void set (int value);
}
```

```
member MaxConnector : int with get, set
```

#### 属性值

类型：Int32

![](../icons/SectionExpanded.png)参见

#### 引用

[ConnectPoolTConnector 类](8697bcd0-0fa1-9274-3029-3ee6a893fd56.htm)

[HslCommunication.Algorithms.ConnectPool 命名空间](9ec7e0cd-f203-a971-db9f-ea1f53b0b4b3.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## UseConnectorMax 属性 

[原文連結](http://api.hslcommunication.cn/html/feef2662-bb80-a0f0-9db8-7f8fd2c34197.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Algorithms.ConnectPool](../html/9ec7e0cd-f203-a971-db9f-ea1f53b0b4b3.htm "HslCommunication.Algorithms.ConnectPool")

[ConnectPool(TConnector) 类](../html/8697bcd0-0fa1-9274-3029-3ee6a893fd56.htm "ConnectPool(TConnector) 类")

[ConnectPool(TConnector) 属性](../html/73798dd3-7274-847c-6aec-dae2533dacef.htm "ConnectPool(TConnector) 属性")

[ConectionExpireTime 属性](../html/9f233c68-58fa-91e5-8a28-09b6cd0c10d5.htm "ConectionExpireTime 属性 ")

[MaxConnector 属性](../html/9dc834aa-57d2-ee27-248f-8118aded699a.htm "MaxConnector 属性 ")

[UseConnectorMax 属性](../html/feef2662-bb80-a0f0-9db8-7f8fd2c34197.htm "UseConnectorMax 属性 ")

[UsedConnector 属性](../html/f96d08a2-465a-8076-82d7-222c0b3d225f.htm "UsedConnector 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ConnectPoolTConnectorUseConnectorMax 属性 |

当前已经使用的连接数的峰值，可以用来衡量当前系统的适用的连接池上限。  
The current peak value of the number of connections used can be used to measure the upper limit of the applicable connection pool of the current system.

**命名空间：**
 [HslCommunication.Algorithms.ConnectPool](9ec7e0cd-f203-a971-db9f-ea1f53b0b4b3.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public int UseConnectorMax { get; }
```

```
Public ReadOnly Property UseConnectorMax As Integer
	Get
```

```
public:
property int UseConnectorMax {
	int get ();
}
```

```
member UseConnectorMax : int with get
```

#### 属性值

类型：Int32

![](../icons/SectionExpanded.png)参见

#### 引用

[ConnectPoolTConnector 类](8697bcd0-0fa1-9274-3029-3ee6a893fd56.htm)

[HslCommunication.Algorithms.ConnectPool 命名空间](9ec7e0cd-f203-a971-db9f-ea1f53b0b4b3.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## UsedConnector 属性 

[原文連結](http://api.hslcommunication.cn/html/f96d08a2-465a-8076-82d7-222c0b3d225f.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Algorithms.ConnectPool](../html/9ec7e0cd-f203-a971-db9f-ea1f53b0b4b3.htm "HslCommunication.Algorithms.ConnectPool")

[ConnectPool(TConnector) 类](../html/8697bcd0-0fa1-9274-3029-3ee6a893fd56.htm "ConnectPool(TConnector) 类")

[ConnectPool(TConnector) 属性](../html/73798dd3-7274-847c-6aec-dae2533dacef.htm "ConnectPool(TConnector) 属性")

[ConectionExpireTime 属性](../html/9f233c68-58fa-91e5-8a28-09b6cd0c10d5.htm "ConectionExpireTime 属性 ")

[MaxConnector 属性](../html/9dc834aa-57d2-ee27-248f-8118aded699a.htm "MaxConnector 属性 ")

[UseConnectorMax 属性](../html/feef2662-bb80-a0f0-9db8-7f8fd2c34197.htm "UseConnectorMax 属性 ")

[UsedConnector 属性](../html/f96d08a2-465a-8076-82d7-222c0b3d225f.htm "UsedConnector 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ConnectPoolTConnectorUsedConnector 属性 |

当前已经使用的连接数，会根据使用的频繁程度进行动态的变化。  
The number of currently used connections will dynamically change according to the frequency of use.

**命名空间：**
 [HslCommunication.Algorithms.ConnectPool](9ec7e0cd-f203-a971-db9f-ea1f53b0b4b3.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public int UsedConnector { get; }
```

```
Public ReadOnly Property UsedConnector As Integer
	Get
```

```
public:
property int UsedConnector {
	int get ();
}
```

```
member UsedConnector : int with get
```

#### 属性值

类型：Int32

![](../icons/SectionExpanded.png)参见

#### 引用

[ConnectPoolTConnector 类](8697bcd0-0fa1-9274-3029-3ee6a893fd56.htm)

[HslCommunication.Algorithms.ConnectPool 命名空间](9ec7e0cd-f203-a971-db9f-ea1f53b0b4b3.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ConnectPool(TConnector) 方法

[原文連結](http://api.hslcommunication.cn/html/c22ef39f-c4e9-7389-66d5-37459c7522f6.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Algorithms.ConnectPool](../html/9ec7e0cd-f203-a971-db9f-ea1f53b0b4b3.htm "HslCommunication.Algorithms.ConnectPool")

[ConnectPool(TConnector) 类](../html/8697bcd0-0fa1-9274-3029-3ee6a893fd56.htm "ConnectPool(TConnector) 类")

[ConnectPool(TConnector) 方法](../html/c22ef39f-c4e9-7389-66d5-37459c7522f6.htm "ConnectPool(TConnector) 方法")

[GetAvailableConnector 方法](../html/4db5960c-88c8-8abc-28e5-f37df26e66ed.htm "GetAvailableConnector 方法 ")

[ResetAllConnector 方法](../html/e183a5a2-bbf4-3557-34c6-44d16a8bb413.htm "ResetAllConnector 方法 ")

[ReturnConnector 方法](../html/278f482c-b77b-294c-6ef4-5d598378ff78.htm "ReturnConnector 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ConnectPoolTConnector 方法 |

[ConnectPoolTConnector](8697bcd0-0fa1-9274-3029-3ee6a893fd56.htm) 泛型公开以下成员。

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | [GetAvailableConnector](4db5960c-88c8-8abc-28e5-f37df26e66ed.htm) | 获取一个可用的连接对象，如果已经达到上限，就进行阻塞等待。当使用完连接对象的时候，需要调用[ReturnConnector(TConnector)](278f482c-b77b-294c-6ef4-5d598378ff78.htm)方法归还连接对象。  Get an available connection object, if the upper limit has been reached, block waiting. When the connection object is used up, you need to call the [ReturnConnector(TConnector)](278f482c-b77b-294c-6ef4-5d598378ff78.htm) method to return the connection object. |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法 | [ResetAllConnector](e183a5a2-bbf4-3557-34c6-44d16a8bb413.htm) | 将目前连接中的所有对象进行关闭，然后移除队列。  Close all objects in the current connection, and then remove the queue. |
| 公共方法 | [ReturnConnector](278f482c-b77b-294c-6ef4-5d598378ff78.htm) | 使用完之后需要通知连接池的管理器，本方法调用之前需要获取到连接对象信息。  After using it, you need to notify the manager of the connection pool, and you need to get the connection object information before calling this method. |
| 公共方法 | ToString | (继承自 Object。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[ConnectPoolTConnector 类](8697bcd0-0fa1-9274-3029-3ee6a893fd56.htm)

[HslCommunication.Algorithms.ConnectPool 命名空间](9ec7e0cd-f203-a971-db9f-ea1f53b0b4b3.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GetAvailableConnector 方法 

[原文連結](http://api.hslcommunication.cn/html/4db5960c-88c8-8abc-28e5-f37df26e66ed.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Algorithms.ConnectPool](../html/9ec7e0cd-f203-a971-db9f-ea1f53b0b4b3.htm "HslCommunication.Algorithms.ConnectPool")

[ConnectPool(TConnector) 类](../html/8697bcd0-0fa1-9274-3029-3ee6a893fd56.htm "ConnectPool(TConnector) 类")

[ConnectPool(TConnector) 方法](../html/c22ef39f-c4e9-7389-66d5-37459c7522f6.htm "ConnectPool(TConnector) 方法")

[GetAvailableConnector 方法](../html/4db5960c-88c8-8abc-28e5-f37df26e66ed.htm "GetAvailableConnector 方法 ")

[ResetAllConnector 方法](../html/e183a5a2-bbf4-3557-34c6-44d16a8bb413.htm "ResetAllConnector 方法 ")

[ReturnConnector 方法](../html/278f482c-b77b-294c-6ef4-5d598378ff78.htm "ReturnConnector 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ConnectPoolTConnectorGetAvailableConnector 方法 |

获取一个可用的连接对象，如果已经达到上限，就进行阻塞等待。当使用完连接对象的时候，需要调用[ReturnConnector(TConnector)](278f482c-b77b-294c-6ef4-5d598378ff78.htm)方法归还连接对象。  
Get an available connection object, if the upper limit has been reached, block waiting. When the connection object is used up,
you need to call the [ReturnConnector(TConnector)](278f482c-b77b-294c-6ef4-5d598378ff78.htm) method to return the connection object.

**命名空间：**
 [HslCommunication.Algorithms.ConnectPool](9ec7e0cd-f203-a971-db9f-ea1f53b0b4b3.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public TConnector GetAvailableConnector()
```

```
Public Function GetAvailableConnector As TConnector
```

```
public:
TConnector GetAvailableConnector()
```

```
member GetAvailableConnector : unit -> 'TConnector 
```

#### 返回值

类型：[TConnector](8697bcd0-0fa1-9274-3029-3ee6a893fd56.htm)  
可用的连接对象

![](../icons/SectionExpanded.png)参见

#### 引用

[ConnectPoolTConnector 类](8697bcd0-0fa1-9274-3029-3ee6a893fd56.htm)

[HslCommunication.Algorithms.ConnectPool 命名空间](9ec7e0cd-f203-a971-db9f-ea1f53b0b4b3.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ResetAllConnector 方法 

[原文連結](http://api.hslcommunication.cn/html/e183a5a2-bbf4-3557-34c6-44d16a8bb413.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Algorithms.ConnectPool](../html/9ec7e0cd-f203-a971-db9f-ea1f53b0b4b3.htm "HslCommunication.Algorithms.ConnectPool")

[ConnectPool(TConnector) 类](../html/8697bcd0-0fa1-9274-3029-3ee6a893fd56.htm "ConnectPool(TConnector) 类")

[ConnectPool(TConnector) 方法](../html/c22ef39f-c4e9-7389-66d5-37459c7522f6.htm "ConnectPool(TConnector) 方法")

[GetAvailableConnector 方法](../html/4db5960c-88c8-8abc-28e5-f37df26e66ed.htm "GetAvailableConnector 方法 ")

[ResetAllConnector 方法](../html/e183a5a2-bbf4-3557-34c6-44d16a8bb413.htm "ResetAllConnector 方法 ")

[ReturnConnector 方法](../html/278f482c-b77b-294c-6ef4-5d598378ff78.htm "ReturnConnector 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ConnectPoolTConnectorResetAllConnector 方法 |

将目前连接中的所有对象进行关闭，然后移除队列。  
Close all objects in the current connection, and then remove the queue.

**命名空间：**
 [HslCommunication.Algorithms.ConnectPool](9ec7e0cd-f203-a971-db9f-ea1f53b0b4b3.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public void ResetAllConnector()
```

```
Public Sub ResetAllConnector
```

```
public:
void ResetAllConnector()
```

```
member ResetAllConnector : unit -> unit 
```

![](../icons/SectionExpanded.png)参见

#### 引用

[ConnectPoolTConnector 类](8697bcd0-0fa1-9274-3029-3ee6a893fd56.htm)

[HslCommunication.Algorithms.ConnectPool 命名空间](9ec7e0cd-f203-a971-db9f-ea1f53b0b4b3.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReturnConnector 方法 

[原文連結](http://api.hslcommunication.cn/html/278f482c-b77b-294c-6ef4-5d598378ff78.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Algorithms.ConnectPool](../html/9ec7e0cd-f203-a971-db9f-ea1f53b0b4b3.htm "HslCommunication.Algorithms.ConnectPool")

[ConnectPool(TConnector) 类](../html/8697bcd0-0fa1-9274-3029-3ee6a893fd56.htm "ConnectPool(TConnector) 类")

[ConnectPool(TConnector) 方法](../html/c22ef39f-c4e9-7389-66d5-37459c7522f6.htm "ConnectPool(TConnector) 方法")

[GetAvailableConnector 方法](../html/4db5960c-88c8-8abc-28e5-f37df26e66ed.htm "GetAvailableConnector 方法 ")

[ResetAllConnector 方法](../html/e183a5a2-bbf4-3557-34c6-44d16a8bb413.htm "ResetAllConnector 方法 ")

[ReturnConnector 方法](../html/278f482c-b77b-294c-6ef4-5d598378ff78.htm "ReturnConnector 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ConnectPoolTConnectorReturnConnector 方法 |

使用完之后需要通知连接池的管理器，本方法调用之前需要获取到连接对象信息。  
After using it, you need to notify the manager of the connection pool, and you need to get the connection object information before calling this method.

**命名空间：**
 [HslCommunication.Algorithms.ConnectPool](9ec7e0cd-f203-a971-db9f-ea1f53b0b4b3.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public void ReturnConnector(
	TConnector connector
)
```

```
Public Sub ReturnConnector ( 
	connector As TConnector
)
```

```
public:
void ReturnConnector(
	TConnector connector
)
```

```
member ReturnConnector : 
        connector : 'TConnector -> unit 
```

#### 参数

connector
:   类型：[TConnector](8697bcd0-0fa1-9274-3029-3ee6a893fd56.htm)  
    连接对象

![](../icons/SectionExpanded.png)参见

#### 引用

[ConnectPoolTConnector 类](8697bcd0-0fa1-9274-3029-3ee6a893fd56.htm)

[HslCommunication.Algorithms.ConnectPool 命名空间](9ec7e0cd-f203-a971-db9f-ea1f53b0b4b3.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## DeviceCommConnector 类

[原文連結](http://api.hslcommunication.cn/html/9714842a-ab7e-6320-7092-b2a8ca7d7c0e.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Algorithms.ConnectPool](../html/9ec7e0cd-f203-a971-db9f-ea1f53b0b4b3.htm "HslCommunication.Algorithms.ConnectPool")

[DeviceCommConnector 类](../html/9714842a-ab7e-6320-7092-b2a8ca7d7c0e.htm "DeviceCommConnector 类")

[DeviceCommConnector 构造函数](../html/c560d0aa-1721-5972-ecb4-1faea1f81cd3.htm "DeviceCommConnector 构造函数 ")

[DeviceCommConnector 属性](../html/2e6c7eaf-68cc-a4fa-7b7a-93a312942dbb.htm "DeviceCommConnector 属性")

[DeviceCommConnector 方法](../html/51e31065-2f63-2447-7fdf-0401049c2eb2.htm "DeviceCommConnector 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DeviceCommConnector 类 |

基于设备通信的连接池信息  
Connection pool information based on device communication

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  HslCommunication.Algorithms.ConnectPoolDeviceCommConnector

**命名空间：**
 [HslCommunication.Algorithms.ConnectPool](9ec7e0cd-f203-a971-db9f-ea1f53b0b4b3.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class DeviceCommConnector : IConnector
```

```
Public Class DeviceCommConnector
	Implements IConnector
```

```
public ref class DeviceCommConnector : IConnector
```

```
type DeviceCommConnector =  
    class
        interface IConnector
    end
```

DeviceCommConnector 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [DeviceCommConnector](c560d0aa-1721-5972-ecb4-1faea1f81cd3.htm) | 使用指定的设备来实例化一个对象  Use the specified device to instantiate an object |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [Device](c9a30751-3e4e-c7cb-c450-989fbc18b541.htm) | 获取当前实际的设备通信对象  Gets the actual device communication object at present |
| 公共属性 | [GuidToken](c4e3d30a-1b3c-eb47-6640-71f577f13c79.htm) | 唯一的GUID码 |
| 公共属性 | [IsConnectUsing](a0d8bb95-6409-13fe-c29b-ba54ee830713.htm) | 指示当前的连接是否在使用用 |
| 公共属性 | [LastUseTime](1dda7ec4-aae2-0bd0-3f29-656e3c2a210f.htm) | 最新一次使用的时间 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [Close](2964530b-9774-cd60-85a7-78add7a65360.htm) | 关闭并释放 |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法 | [Open](61e4ba6d-9f81-3370-9ff6-8926fd8dac5d.htm) | 打开连接 |
| 公共方法 | ToString | (继承自 Object。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.Algorithms.ConnectPool 命名空间](9ec7e0cd-f203-a971-db9f-ea1f53b0b4b3.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## DeviceCommConnector 构造函数 

[原文連結](http://api.hslcommunication.cn/html/c560d0aa-1721-5972-ecb4-1faea1f81cd3.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Algorithms.ConnectPool](../html/9ec7e0cd-f203-a971-db9f-ea1f53b0b4b3.htm "HslCommunication.Algorithms.ConnectPool")

[DeviceCommConnector 类](../html/9714842a-ab7e-6320-7092-b2a8ca7d7c0e.htm "DeviceCommConnector 类")

[DeviceCommConnector 构造函数](../html/c560d0aa-1721-5972-ecb4-1faea1f81cd3.htm "DeviceCommConnector 构造函数 ")

[DeviceCommConnector 属性](../html/2e6c7eaf-68cc-a4fa-7b7a-93a312942dbb.htm "DeviceCommConnector 属性")

[DeviceCommConnector 方法](../html/51e31065-2f63-2447-7fdf-0401049c2eb2.htm "DeviceCommConnector 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DeviceCommConnector 构造函数 |

使用指定的设备来实例化一个对象  
Use the specified device to instantiate an object

**命名空间：**
 [HslCommunication.Algorithms.ConnectPool](9ec7e0cd-f203-a971-db9f-ea1f53b0b4b3.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public DeviceCommConnector(
	DeviceCommunication device
)
```

```
Public Sub New ( 
	device As DeviceCommunication
)
```

```
public:
DeviceCommConnector(
	DeviceCommunication^ device
)
```

```
new : 
        device : DeviceCommunication -> DeviceCommConnector
```

#### 参数

device
:   类型：[HslCommunication.Core.DeviceDeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)  
    指定的设备通信

![](../icons/SectionExpanded.png)参见

#### 引用

[DeviceCommConnector 类](9714842a-ab7e-6320-7092-b2a8ca7d7c0e.htm)

[HslCommunication.Algorithms.ConnectPool 命名空间](9ec7e0cd-f203-a971-db9f-ea1f53b0b4b3.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## DeviceCommConnector 属性

[原文連結](http://api.hslcommunication.cn/html/2e6c7eaf-68cc-a4fa-7b7a-93a312942dbb.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Algorithms.ConnectPool](../html/9ec7e0cd-f203-a971-db9f-ea1f53b0b4b3.htm "HslCommunication.Algorithms.ConnectPool")

[DeviceCommConnector 类](../html/9714842a-ab7e-6320-7092-b2a8ca7d7c0e.htm "DeviceCommConnector 类")

[DeviceCommConnector 属性](../html/2e6c7eaf-68cc-a4fa-7b7a-93a312942dbb.htm "DeviceCommConnector 属性")

[Device 属性](../html/c9a30751-3e4e-c7cb-c450-989fbc18b541.htm "Device 属性 ")

[GuidToken 属性](../html/c4e3d30a-1b3c-eb47-6640-71f577f13c79.htm "GuidToken 属性 ")

[IsConnectUsing 属性](../html/a0d8bb95-6409-13fe-c29b-ba54ee830713.htm "IsConnectUsing 属性 ")

[LastUseTime 属性](../html/1dda7ec4-aae2-0bd0-3f29-656e3c2a210f.htm "LastUseTime 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DeviceCommConnector 属性 |

[DeviceCommConnector](9714842a-ab7e-6320-7092-b2a8ca7d7c0e.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [Device](c9a30751-3e4e-c7cb-c450-989fbc18b541.htm) | 获取当前实际的设备通信对象  Gets the actual device communication object at present |
| 公共属性 | [GuidToken](c4e3d30a-1b3c-eb47-6640-71f577f13c79.htm) | 唯一的GUID码 |
| 公共属性 | [IsConnectUsing](a0d8bb95-6409-13fe-c29b-ba54ee830713.htm) | 指示当前的连接是否在使用用 |
| 公共属性 | [LastUseTime](1dda7ec4-aae2-0bd0-3f29-656e3c2a210f.htm) | 最新一次使用的时间 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[DeviceCommConnector 类](9714842a-ab7e-6320-7092-b2a8ca7d7c0e.htm)

[HslCommunication.Algorithms.ConnectPool 命名空间](9ec7e0cd-f203-a971-db9f-ea1f53b0b4b3.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Device 属性 

[原文連結](http://api.hslcommunication.cn/html/c9a30751-3e4e-c7cb-c450-989fbc18b541.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Algorithms.ConnectPool](../html/9ec7e0cd-f203-a971-db9f-ea1f53b0b4b3.htm "HslCommunication.Algorithms.ConnectPool")

[DeviceCommConnector 类](../html/9714842a-ab7e-6320-7092-b2a8ca7d7c0e.htm "DeviceCommConnector 类")

[DeviceCommConnector 属性](../html/2e6c7eaf-68cc-a4fa-7b7a-93a312942dbb.htm "DeviceCommConnector 属性")

[Device 属性](../html/c9a30751-3e4e-c7cb-c450-989fbc18b541.htm "Device 属性 ")

[GuidToken 属性](../html/c4e3d30a-1b3c-eb47-6640-71f577f13c79.htm "GuidToken 属性 ")

[IsConnectUsing 属性](../html/a0d8bb95-6409-13fe-c29b-ba54ee830713.htm "IsConnectUsing 属性 ")

[LastUseTime 属性](../html/1dda7ec4-aae2-0bd0-3f29-656e3c2a210f.htm "LastUseTime 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DeviceCommConnectorDevice 属性 |

获取当前实际的设备通信对象  
Gets the actual device communication object at present

**命名空间：**
 [HslCommunication.Algorithms.ConnectPool](9ec7e0cd-f203-a971-db9f-ea1f53b0b4b3.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public DeviceCommunication Device { get; }
```

```
Public ReadOnly Property Device As DeviceCommunication
	Get
```

```
public:
property DeviceCommunication^ Device {
	DeviceCommunication^ get ();
}
```

```
member Device : DeviceCommunication with get
```

#### 属性值

类型：[DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[DeviceCommConnector 类](9714842a-ab7e-6320-7092-b2a8ca7d7c0e.htm)

[HslCommunication.Algorithms.ConnectPool 命名空间](9ec7e0cd-f203-a971-db9f-ea1f53b0b4b3.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GuidToken 属性 

[原文連結](http://api.hslcommunication.cn/html/c4e3d30a-1b3c-eb47-6640-71f577f13c79.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Algorithms.ConnectPool](../html/9ec7e0cd-f203-a971-db9f-ea1f53b0b4b3.htm "HslCommunication.Algorithms.ConnectPool")

[DeviceCommConnector 类](../html/9714842a-ab7e-6320-7092-b2a8ca7d7c0e.htm "DeviceCommConnector 类")

[DeviceCommConnector 属性](../html/2e6c7eaf-68cc-a4fa-7b7a-93a312942dbb.htm "DeviceCommConnector 属性")

[Device 属性](../html/c9a30751-3e4e-c7cb-c450-989fbc18b541.htm "Device 属性 ")

[GuidToken 属性](../html/c4e3d30a-1b3c-eb47-6640-71f577f13c79.htm "GuidToken 属性 ")

[IsConnectUsing 属性](../html/a0d8bb95-6409-13fe-c29b-ba54ee830713.htm "IsConnectUsing 属性 ")

[LastUseTime 属性](../html/1dda7ec4-aae2-0bd0-3f29-656e3c2a210f.htm "LastUseTime 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DeviceCommConnectorGuidToken 属性 |

唯一的GUID码

**命名空间：**
 [HslCommunication.Algorithms.ConnectPool](9ec7e0cd-f203-a971-db9f-ea1f53b0b4b3.htm)  
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

[DeviceCommConnector 类](9714842a-ab7e-6320-7092-b2a8ca7d7c0e.htm)

[HslCommunication.Algorithms.ConnectPool 命名空间](9ec7e0cd-f203-a971-db9f-ea1f53b0b4b3.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## IsConnectUsing 属性 

[原文連結](http://api.hslcommunication.cn/html/a0d8bb95-6409-13fe-c29b-ba54ee830713.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Algorithms.ConnectPool](../html/9ec7e0cd-f203-a971-db9f-ea1f53b0b4b3.htm "HslCommunication.Algorithms.ConnectPool")

[DeviceCommConnector 类](../html/9714842a-ab7e-6320-7092-b2a8ca7d7c0e.htm "DeviceCommConnector 类")

[DeviceCommConnector 属性](../html/2e6c7eaf-68cc-a4fa-7b7a-93a312942dbb.htm "DeviceCommConnector 属性")

[Device 属性](../html/c9a30751-3e4e-c7cb-c450-989fbc18b541.htm "Device 属性 ")

[GuidToken 属性](../html/c4e3d30a-1b3c-eb47-6640-71f577f13c79.htm "GuidToken 属性 ")

[IsConnectUsing 属性](../html/a0d8bb95-6409-13fe-c29b-ba54ee830713.htm "IsConnectUsing 属性 ")

[LastUseTime 属性](../html/1dda7ec4-aae2-0bd0-3f29-656e3c2a210f.htm "LastUseTime 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DeviceCommConnectorIsConnectUsing 属性 |

指示当前的连接是否在使用用

**命名空间：**
 [HslCommunication.Algorithms.ConnectPool](9ec7e0cd-f203-a971-db9f-ea1f53b0b4b3.htm)  
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

[DeviceCommConnector 类](9714842a-ab7e-6320-7092-b2a8ca7d7c0e.htm)

[HslCommunication.Algorithms.ConnectPool 命名空间](9ec7e0cd-f203-a971-db9f-ea1f53b0b4b3.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## LastUseTime 属性 

[原文連結](http://api.hslcommunication.cn/html/1dda7ec4-aae2-0bd0-3f29-656e3c2a210f.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Algorithms.ConnectPool](../html/9ec7e0cd-f203-a971-db9f-ea1f53b0b4b3.htm "HslCommunication.Algorithms.ConnectPool")

[DeviceCommConnector 类](../html/9714842a-ab7e-6320-7092-b2a8ca7d7c0e.htm "DeviceCommConnector 类")

[DeviceCommConnector 属性](../html/2e6c7eaf-68cc-a4fa-7b7a-93a312942dbb.htm "DeviceCommConnector 属性")

[Device 属性](../html/c9a30751-3e4e-c7cb-c450-989fbc18b541.htm "Device 属性 ")

[GuidToken 属性](../html/c4e3d30a-1b3c-eb47-6640-71f577f13c79.htm "GuidToken 属性 ")

[IsConnectUsing 属性](../html/a0d8bb95-6409-13fe-c29b-ba54ee830713.htm "IsConnectUsing 属性 ")

[LastUseTime 属性](../html/1dda7ec4-aae2-0bd0-3f29-656e3c2a210f.htm "LastUseTime 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DeviceCommConnectorLastUseTime 属性 |

最新一次使用的时间

**命名空间：**
 [HslCommunication.Algorithms.ConnectPool](9ec7e0cd-f203-a971-db9f-ea1f53b0b4b3.htm)  
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

[DeviceCommConnector 类](9714842a-ab7e-6320-7092-b2a8ca7d7c0e.htm)

[HslCommunication.Algorithms.ConnectPool 命名空间](9ec7e0cd-f203-a971-db9f-ea1f53b0b4b3.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## DeviceCommConnector 方法

[原文連結](http://api.hslcommunication.cn/html/51e31065-2f63-2447-7fdf-0401049c2eb2.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Algorithms.ConnectPool](../html/9ec7e0cd-f203-a971-db9f-ea1f53b0b4b3.htm "HslCommunication.Algorithms.ConnectPool")

[DeviceCommConnector 类](../html/9714842a-ab7e-6320-7092-b2a8ca7d7c0e.htm "DeviceCommConnector 类")

[DeviceCommConnector 方法](../html/51e31065-2f63-2447-7fdf-0401049c2eb2.htm "DeviceCommConnector 方法")

[Close 方法](../html/2964530b-9774-cd60-85a7-78add7a65360.htm "Close 方法 ")

[Open 方法](../html/61e4ba6d-9f81-3370-9ff6-8926fd8dac5d.htm "Open 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DeviceCommConnector 方法 |

[DeviceCommConnector](9714842a-ab7e-6320-7092-b2a8ca7d7c0e.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [Close](2964530b-9774-cd60-85a7-78add7a65360.htm) | 关闭并释放 |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法 | [Open](61e4ba6d-9f81-3370-9ff6-8926fd8dac5d.htm) | 打开连接 |
| 公共方法 | ToString | (继承自 Object。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[DeviceCommConnector 类](9714842a-ab7e-6320-7092-b2a8ca7d7c0e.htm)

[HslCommunication.Algorithms.ConnectPool 命名空间](9ec7e0cd-f203-a971-db9f-ea1f53b0b4b3.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Close 方法 

[原文連結](http://api.hslcommunication.cn/html/2964530b-9774-cd60-85a7-78add7a65360.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Algorithms.ConnectPool](../html/9ec7e0cd-f203-a971-db9f-ea1f53b0b4b3.htm "HslCommunication.Algorithms.ConnectPool")

[DeviceCommConnector 类](../html/9714842a-ab7e-6320-7092-b2a8ca7d7c0e.htm "DeviceCommConnector 类")

[DeviceCommConnector 方法](../html/51e31065-2f63-2447-7fdf-0401049c2eb2.htm "DeviceCommConnector 方法")

[Close 方法](../html/2964530b-9774-cd60-85a7-78add7a65360.htm "Close 方法 ")

[Open 方法](../html/61e4ba6d-9f81-3370-9ff6-8926fd8dac5d.htm "Open 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DeviceCommConnectorClose 方法 |

关闭并释放

**命名空间：**
 [HslCommunication.Algorithms.ConnectPool](9ec7e0cd-f203-a971-db9f-ea1f53b0b4b3.htm)  
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

[DeviceCommConnector 类](9714842a-ab7e-6320-7092-b2a8ca7d7c0e.htm)

[HslCommunication.Algorithms.ConnectPool 命名空间](9ec7e0cd-f203-a971-db9f-ea1f53b0b4b3.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Open 方法 

[原文連結](http://api.hslcommunication.cn/html/61e4ba6d-9f81-3370-9ff6-8926fd8dac5d.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Algorithms.ConnectPool](../html/9ec7e0cd-f203-a971-db9f-ea1f53b0b4b3.htm "HslCommunication.Algorithms.ConnectPool")

[DeviceCommConnector 类](../html/9714842a-ab7e-6320-7092-b2a8ca7d7c0e.htm "DeviceCommConnector 类")

[DeviceCommConnector 方法](../html/51e31065-2f63-2447-7fdf-0401049c2eb2.htm "DeviceCommConnector 方法")

[Close 方法](../html/2964530b-9774-cd60-85a7-78add7a65360.htm "Close 方法 ")

[Open 方法](../html/61e4ba6d-9f81-3370-9ff6-8926fd8dac5d.htm "Open 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DeviceCommConnectorOpen 方法 |

打开连接

**命名空间：**
 [HslCommunication.Algorithms.ConnectPool](9ec7e0cd-f203-a971-db9f-ea1f53b0b4b3.htm)  
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

[DeviceCommConnector 类](9714842a-ab7e-6320-7092-b2a8ca7d7c0e.htm)

[HslCommunication.Algorithms.ConnectPool 命名空间](9ec7e0cd-f203-a971-db9f-ea1f53b0b4b3.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## IConnector 接口

[原文連結](http://api.hslcommunication.cn/html/ede7f47f-b60e-e320-8bc2-389883e77083.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Algorithms.ConnectPool](../html/9ec7e0cd-f203-a971-db9f-ea1f53b0b4b3.htm "HslCommunication.Algorithms.ConnectPool")

[IConnector 接口](../html/ede7f47f-b60e-e320-8bc2-389883e77083.htm "IConnector 接口")

[IConnector 属性](../html/e916cee1-6786-375f-3534-f4bdc8cb689b.htm "IConnector 属性")

[IConnector 方法](../html/3a534303-8ef9-c7d2-9f66-46835568ff9d.htm "IConnector 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IConnector 接口 |

连接池的接口，连接池的管理对象必须实现此接口

**命名空间：**
 [HslCommunication.Algorithms.ConnectPool](9ec7e0cd-f203-a971-db9f-ea1f53b0b4b3.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public interface IConnector
```

```
Public Interface IConnector
```

```
public interface class IConnector
```

```
type IConnector =  interface end
```

IConnector 类型公开以下成员。

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [GuidToken](5d2d5ad0-74d8-0b6d-29a0-5cc6b0ead492.htm) | 唯一的GUID码 |
| 公共属性 | [IsConnectUsing](60329f21-7e63-6365-78c8-c4aad8efaa85.htm) | 指示当前的连接是否在使用用 |
| 公共属性 | [LastUseTime](25b06530-9b51-19db-ba38-a91998928c48.htm) | 最新一次使用的时间 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [Close](be1b3bb2-108e-0069-994c-c7a14953cc4e.htm) | 关闭并释放 |
| 公共方法 | [Open](08379c6c-8548-7e37-08d4-6a8b71bfd509.htm) | 打开连接 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)备注

为了使用完整的连接池功能，需要先实现本接口，然后配合[ConnectPoolTConnector](8697bcd0-0fa1-9274-3029-3ee6a893fd56.htm)来使用

![](../icons/SectionExpanded.png)示例

下面举例实现一个modbus的连接池对象

IConnector示例

[复制](# "复制")

```
/// <summary>
/// 此处示例实现一个modbus-tcp连接对象，事实上这里可以实现任何的连接对象，PLC的，数据库的，redis的等等操作
/// </summary>
public class ModbusConnector : IConnector
{

    private ModbusTcpNet modbusTcp = null;

    public ModbusConnector( string ipAddress, int port )
    {
        modbusTcp = new ModbusTcpNet( ipAddress, port, 0x01 );   // 默认站号1
    }


    public ModbusTcpNet ModbusTcp
    {
        get { return modbusTcp; }
    } 

    public bool IsConnectUsing { get; set; }


    public string GuidToken { get; set; }


    public DateTime LastUseTime { get; set; }


    public void Close( )
    {
        modbusTcp.ConnectClose( );
    }


    public void Open( )
    {

    }
}
```

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.Algorithms.ConnectPool 命名空间](9ec7e0cd-f203-a971-db9f-ea1f53b0b4b3.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## IConnector 属性

[原文連結](http://api.hslcommunication.cn/html/e916cee1-6786-375f-3534-f4bdc8cb689b.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Algorithms.ConnectPool](../html/9ec7e0cd-f203-a971-db9f-ea1f53b0b4b3.htm "HslCommunication.Algorithms.ConnectPool")

[IConnector 接口](../html/ede7f47f-b60e-e320-8bc2-389883e77083.htm "IConnector 接口")

[IConnector 属性](../html/e916cee1-6786-375f-3534-f4bdc8cb689b.htm "IConnector 属性")

[GuidToken 属性](../html/5d2d5ad0-74d8-0b6d-29a0-5cc6b0ead492.htm "GuidToken 属性 ")

[IsConnectUsing 属性](../html/60329f21-7e63-6365-78c8-c4aad8efaa85.htm "IsConnectUsing 属性 ")

[LastUseTime 属性](../html/25b06530-9b51-19db-ba38-a91998928c48.htm "LastUseTime 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IConnector 属性 |

[IConnector](ede7f47f-b60e-e320-8bc2-389883e77083.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [GuidToken](5d2d5ad0-74d8-0b6d-29a0-5cc6b0ead492.htm) | 唯一的GUID码 |
| 公共属性 | [IsConnectUsing](60329f21-7e63-6365-78c8-c4aad8efaa85.htm) | 指示当前的连接是否在使用用 |
| 公共属性 | [LastUseTime](25b06530-9b51-19db-ba38-a91998928c48.htm) | 最新一次使用的时间 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[IConnector 接口](ede7f47f-b60e-e320-8bc2-389883e77083.htm)

[HslCommunication.Algorithms.ConnectPool 命名空间](9ec7e0cd-f203-a971-db9f-ea1f53b0b4b3.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GuidToken 属性 

[原文連結](http://api.hslcommunication.cn/html/5d2d5ad0-74d8-0b6d-29a0-5cc6b0ead492.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Algorithms.ConnectPool](../html/9ec7e0cd-f203-a971-db9f-ea1f53b0b4b3.htm "HslCommunication.Algorithms.ConnectPool")

[IConnector 接口](../html/ede7f47f-b60e-e320-8bc2-389883e77083.htm "IConnector 接口")

[IConnector 属性](../html/e916cee1-6786-375f-3534-f4bdc8cb689b.htm "IConnector 属性")

[GuidToken 属性](../html/5d2d5ad0-74d8-0b6d-29a0-5cc6b0ead492.htm "GuidToken 属性 ")

[IsConnectUsing 属性](../html/60329f21-7e63-6365-78c8-c4aad8efaa85.htm "IsConnectUsing 属性 ")

[LastUseTime 属性](../html/25b06530-9b51-19db-ba38-a91998928c48.htm "LastUseTime 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IConnectorGuidToken 属性 |

唯一的GUID码

**命名空间：**
 [HslCommunication.Algorithms.ConnectPool](9ec7e0cd-f203-a971-db9f-ea1f53b0b4b3.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
string GuidToken { get; set; }
```

```
Property GuidToken As String
	Get
	Set
```

```
property String^ GuidToken {
	String^ get ();
	void set (String^ value);
}
```

```
abstract GuidToken : string with get, set
```

#### 属性值

类型：String

![](../icons/SectionExpanded.png)参见

#### 引用

[IConnector 接口](ede7f47f-b60e-e320-8bc2-389883e77083.htm)

[HslCommunication.Algorithms.ConnectPool 命名空间](9ec7e0cd-f203-a971-db9f-ea1f53b0b4b3.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## IsConnectUsing 属性 

[原文連結](http://api.hslcommunication.cn/html/60329f21-7e63-6365-78c8-c4aad8efaa85.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Algorithms.ConnectPool](../html/9ec7e0cd-f203-a971-db9f-ea1f53b0b4b3.htm "HslCommunication.Algorithms.ConnectPool")

[IConnector 接口](../html/ede7f47f-b60e-e320-8bc2-389883e77083.htm "IConnector 接口")

[IConnector 属性](../html/e916cee1-6786-375f-3534-f4bdc8cb689b.htm "IConnector 属性")

[GuidToken 属性](../html/5d2d5ad0-74d8-0b6d-29a0-5cc6b0ead492.htm "GuidToken 属性 ")

[IsConnectUsing 属性](../html/60329f21-7e63-6365-78c8-c4aad8efaa85.htm "IsConnectUsing 属性 ")

[LastUseTime 属性](../html/25b06530-9b51-19db-ba38-a91998928c48.htm "LastUseTime 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IConnectorIsConnectUsing 属性 |

指示当前的连接是否在使用用

**命名空间：**
 [HslCommunication.Algorithms.ConnectPool](9ec7e0cd-f203-a971-db9f-ea1f53b0b4b3.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
bool IsConnectUsing { get; set; }
```

```
Property IsConnectUsing As Boolean
	Get
	Set
```

```
property bool IsConnectUsing {
	bool get ();
	void set (bool value);
}
```

```
abstract IsConnectUsing : bool with get, set
```

#### 属性值

类型：Boolean

![](../icons/SectionExpanded.png)参见

#### 引用

[IConnector 接口](ede7f47f-b60e-e320-8bc2-389883e77083.htm)

[HslCommunication.Algorithms.ConnectPool 命名空间](9ec7e0cd-f203-a971-db9f-ea1f53b0b4b3.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## LastUseTime 属性 

[原文連結](http://api.hslcommunication.cn/html/25b06530-9b51-19db-ba38-a91998928c48.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Algorithms.ConnectPool](../html/9ec7e0cd-f203-a971-db9f-ea1f53b0b4b3.htm "HslCommunication.Algorithms.ConnectPool")

[IConnector 接口](../html/ede7f47f-b60e-e320-8bc2-389883e77083.htm "IConnector 接口")

[IConnector 属性](../html/e916cee1-6786-375f-3534-f4bdc8cb689b.htm "IConnector 属性")

[GuidToken 属性](../html/5d2d5ad0-74d8-0b6d-29a0-5cc6b0ead492.htm "GuidToken 属性 ")

[IsConnectUsing 属性](../html/60329f21-7e63-6365-78c8-c4aad8efaa85.htm "IsConnectUsing 属性 ")

[LastUseTime 属性](../html/25b06530-9b51-19db-ba38-a91998928c48.htm "LastUseTime 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IConnectorLastUseTime 属性 |

最新一次使用的时间

**命名空间：**
 [HslCommunication.Algorithms.ConnectPool](9ec7e0cd-f203-a971-db9f-ea1f53b0b4b3.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
DateTime LastUseTime { get; set; }
```

```
Property LastUseTime As DateTime
	Get
	Set
```

```
property DateTime LastUseTime {
	DateTime get ();
	void set (DateTime value);
}
```

```
abstract LastUseTime : DateTime with get, set
```

#### 属性值

类型：DateTime

![](../icons/SectionExpanded.png)参见

#### 引用

[IConnector 接口](ede7f47f-b60e-e320-8bc2-389883e77083.htm)

[HslCommunication.Algorithms.ConnectPool 命名空间](9ec7e0cd-f203-a971-db9f-ea1f53b0b4b3.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## IConnector 方法

[原文連結](http://api.hslcommunication.cn/html/3a534303-8ef9-c7d2-9f66-46835568ff9d.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Algorithms.ConnectPool](../html/9ec7e0cd-f203-a971-db9f-ea1f53b0b4b3.htm "HslCommunication.Algorithms.ConnectPool")

[IConnector 接口](../html/ede7f47f-b60e-e320-8bc2-389883e77083.htm "IConnector 接口")

[IConnector 方法](../html/3a534303-8ef9-c7d2-9f66-46835568ff9d.htm "IConnector 方法")

[Close 方法](../html/be1b3bb2-108e-0069-994c-c7a14953cc4e.htm "Close 方法 ")

[Open 方法](../html/08379c6c-8548-7e37-08d4-6a8b71bfd509.htm "Open 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IConnector 方法 |

[IConnector](ede7f47f-b60e-e320-8bc2-389883e77083.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [Close](be1b3bb2-108e-0069-994c-c7a14953cc4e.htm) | 关闭并释放 |
| 公共方法 | [Open](08379c6c-8548-7e37-08d4-6a8b71bfd509.htm) | 打开连接 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[IConnector 接口](ede7f47f-b60e-e320-8bc2-389883e77083.htm)

[HslCommunication.Algorithms.ConnectPool 命名空间](9ec7e0cd-f203-a971-db9f-ea1f53b0b4b3.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Close 方法 

[原文連結](http://api.hslcommunication.cn/html/be1b3bb2-108e-0069-994c-c7a14953cc4e.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Algorithms.ConnectPool](../html/9ec7e0cd-f203-a971-db9f-ea1f53b0b4b3.htm "HslCommunication.Algorithms.ConnectPool")

[IConnector 接口](../html/ede7f47f-b60e-e320-8bc2-389883e77083.htm "IConnector 接口")

[IConnector 方法](../html/3a534303-8ef9-c7d2-9f66-46835568ff9d.htm "IConnector 方法")

[Close 方法](../html/be1b3bb2-108e-0069-994c-c7a14953cc4e.htm "Close 方法 ")

[Open 方法](../html/08379c6c-8548-7e37-08d4-6a8b71bfd509.htm "Open 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IConnectorClose 方法 |

关闭并释放

**命名空间：**
 [HslCommunication.Algorithms.ConnectPool](9ec7e0cd-f203-a971-db9f-ea1f53b0b4b3.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
void Close()
```

```
Sub Close
```

```
void Close()
```

```
abstract Close : unit -> unit 
```

![](../icons/SectionExpanded.png)参见

#### 引用

[IConnector 接口](ede7f47f-b60e-e320-8bc2-389883e77083.htm)

[HslCommunication.Algorithms.ConnectPool 命名空间](9ec7e0cd-f203-a971-db9f-ea1f53b0b4b3.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Open 方法 

[原文連結](http://api.hslcommunication.cn/html/08379c6c-8548-7e37-08d4-6a8b71bfd509.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Algorithms.ConnectPool](../html/9ec7e0cd-f203-a971-db9f-ea1f53b0b4b3.htm "HslCommunication.Algorithms.ConnectPool")

[IConnector 接口](../html/ede7f47f-b60e-e320-8bc2-389883e77083.htm "IConnector 接口")

[IConnector 方法](../html/3a534303-8ef9-c7d2-9f66-46835568ff9d.htm "IConnector 方法")

[Close 方法](../html/be1b3bb2-108e-0069-994c-c7a14953cc4e.htm "Close 方法 ")

[Open 方法](../html/08379c6c-8548-7e37-08d4-6a8b71bfd509.htm "Open 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IConnectorOpen 方法 |

打开连接

**命名空间：**
 [HslCommunication.Algorithms.ConnectPool](9ec7e0cd-f203-a971-db9f-ea1f53b0b4b3.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
void Open()
```

```
Sub Open
```

```
void Open()
```

```
abstract Open : unit -> unit 
```

![](../icons/SectionExpanded.png)参见

#### 引用

[IConnector 接口](ede7f47f-b60e-e320-8bc2-389883e77083.htm)

[HslCommunication.Algorithms.ConnectPool 命名空间](9ec7e0cd-f203-a971-db9f-ea1f53b0b4b3.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)