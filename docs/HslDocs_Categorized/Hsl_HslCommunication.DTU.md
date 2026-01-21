# HslCommunication - HslCommunication.DTU

> 分類頁數: 24



---
## HslCommunication.DTU

[原文連結](http://api.hslcommunication.cn/html/bf8d39e1-bb7b-54e5-7d99-73cbdd1fc4af.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.DTU](../html/bf8d39e1-bb7b-54e5-7d99-73cbdd1fc4af.htm "HslCommunication.DTU")

[DTUServer 类](../html/2e9a3fda-3f3e-82a7-01af-95a940d16410.htm "DTUServer 类")

[DTUSettingType 类](../html/0a9eed32-45a3-d0eb-9104-fe18442a3351.htm "DTUSettingType 类")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslCommunication.DTU 命名空间 |

[缺少 "N:HslCommunication.DTU" 的 <summary> 文档]

![](../icons/SectionExpanded.png)类

|  | 类 | 说明 |
| --- | --- | --- |
| 公共类 | [DTUServer](2e9a3fda-3f3e-82a7-01af-95a940d16410.htm) | DTU的服务器信息，本服务器支持任意的hsl支持的网络对象，包括plc信息，modbus设备等等，通过DTU来连接， 然后支持多个连接对象。如果需要支持非hsl的注册报文，需要重写相关的方法  DTU server information, the server supports any network objects supported by hsl, including plc information, modbus devices, etc., connected through DTU, and then supports multiple connection objects. If you need to support non-HSL registration messages, you need to rewrite the relevant methods |
| 公共类 | [DTUSettingType](0a9eed32-45a3-d0eb-9104-fe18442a3351.htm) | DTU的类型设置器 |

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## DTUServer 类

[原文連結](http://api.hslcommunication.cn/html/2e9a3fda-3f3e-82a7-01af-95a940d16410.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.DTU](../html/bf8d39e1-bb7b-54e5-7d99-73cbdd1fc4af.htm "HslCommunication.DTU")

[DTUServer 类](../html/2e9a3fda-3f3e-82a7-01af-95a940d16410.htm "DTUServer 类")

[DTUServer 构造函数](../html/bcb571b0-02d7-c103-7cfa-98305625d416.htm "DTUServer 构造函数 ")

[DTUServer 属性](../html/2307930e-d020-7b41-ceb5-a1e1f1ccaf78.htm "DTUServer 属性")

[DTUServer 方法](../html/bd9fa841-c520-b11a-85b8-6f4ad77db379.htm "DTUServer 方法")

[DTUServer 事件](../html/38005bd6-daa6-8120-950f-ebf087418ac9.htm "DTUServer 事件")

[DTUServer 字段](../html/b75cd2c5-4594-0645-d74c-d0eba36b388c.htm "DTUServer 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DTUServer 类 |

DTU的服务器信息，本服务器支持任意的hsl支持的网络对象，包括plc信息，modbus设备等等，通过DTU来连接，
然后支持多个连接对象。如果需要支持非hsl的注册报文，需要重写相关的方法  
DTU server information, the server supports any network objects supported by hsl,
including plc information, modbus devices, etc., connected through DTU, and then supports multiple connection objects.
If you need to support non-HSL registration messages, you need to rewrite the relevant methods

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  [HslCommunication.Core.NetCommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)  
    [HslCommunication.Core.NetNetworkAlienClient](e29c23cc-75ae-02b5-2e55-00328c111f37.htm)  
      HslCommunication.DTUDTUServer

**命名空间：**
 [HslCommunication.DTU](bf8d39e1-bb7b-54e5-7d99-73cbdd1fc4af.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class DTUServer : NetworkAlienClient
```

```
Public Class DTUServer
	Inherits NetworkAlienClient
```

```
public ref class DTUServer : public NetworkAlienClient
```

```
type DTUServer =  
    class
        inherit NetworkAlienClient
    end
```

DTUServer 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [DTUServer(ListDTUSettingType)](493a5be9-6693-c4ba-b844-022368256626.htm) | 根据配置的列表信息来实例化相关的DTU服务器  Instantiate the relevant DTU server according to the configured list information |
| 公共方法 | [DTUServer(String, DeviceTcpNet)](02d9a717-15ad-7c11-6499-662cb3e7b303.htm) | 根据配置的列表信息来实例化相关的DTU服务器  Instantiate the relevant DTU server according to the configured list information |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [EnableIPv6](f15cbcb3-5dcf-a658-0a93-3f0921383c32.htm) | 获取或设置服务器是否支持IPv6的地址协议信息  Get or set whether the server supports IPv6 address protocol information (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |
| 公共属性 | [IsCheckPwd](984bb24d-5f9e-5b3e-0a97-2860778c1e26.htm) | 是否统一检查密码，如果每个会话需要自己检查密码，就需要设置为false  Whether to check the password uniformly, if each session needs to check the password by itself, it needs to be set to false (继承自 [NetworkAlienClient](e29c23cc-75ae-02b5-2e55-00328c111f37.htm)。) |
| 公共属性 | [IsResponseAck](f2bdfc03-4e28-57b4-d777-de774e274744.htm) | 在DTU设备发送了注册报文的时候，指示是否返回响应报文，用来通知DTU设备是否登录成功，默认为 True  When a DTU sends a registration packet, it indicates whether to return a response packet to notify the DTU whether the DTU is logged in, The default is True (继承自 [NetworkAlienClient](e29c23cc-75ae-02b5-2e55-00328c111f37.htm)。) |
| 公共属性 | [IsStarted](8d6786a8-d416-4f45-158c-de6be0634aca.htm) | 服务器引擎是否启动  Whether the server engine is started (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |
| 公共属性 | [IsUseSSL](f4daa1e8-f433-0ee2-8284-16c4754eba25.htm) | 获取当前的服务器是否使用了SSL证书功能 (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |
| 公共属性 | [Item](a3462ddd-95fb-0430-a9fb-1f06244ca42a.htm) | 根据DTU信息获取设备的连接对象  Obtain the connection object of the device according to the DTU information |
| 公共属性 | [LocalAddress](c0effb49-cfe6-3267-9aef-afdab27a41c3.htm) | 获取或设置服务器绑定的本地IP地址，默认为空，使用本地所有可用的ip地址  Obtain or set the local IP address bound to the server, which is empty by default, and use all available IP addresses locally (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |
| 公共属性 | [LogDebugMessage](8a6db355-1130-2af8-cc6a-aeccc1ac0b29.htm) | 记录一些调试日志的委托，将会进行输出调试文本。  The delegate that records some debug logs will output debug text. (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |
| 公共属性代码示例 | [LogNet](746b148d-1ee3-1ee7-ee5d-3938eb72b2df.htm) | 组件的日志工具，支持日志记录，只要实例化后，当前网络的基本信息，就以[DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm)等级进行输出  The component's logging tool supports logging. As long as the instantiation of the basic network information, the output will be output at [DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm) (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |
| 公共属性 | [Port](6f700f7a-bd7e-c4b3-1275-a30491469e8f.htm) | 获取或设置服务器的端口号，如果是设置，需要在服务器启动前设置完成，才能生效。  Gets or sets the port number of the server. If it is set, it needs to be set before the server starts to take effect. (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |
| 公共属性 | [SocketKeepAliveTime](7d77e3df-9bbe-1a6e-a051-c87d9516c3a2.htm) | 获取或设置客户端的Socket的心跳时间信息，这个是Socket底层自动实现的心跳包，不基于协议层实现。默认小于0，不开启心跳检测，如果需要开启，设置 60\_000 比较合适，单位毫秒  Get or set the heartbeat time information of the Socket of the client. This is the heartbeat packet automatically implemented by the bottom layer of the Socket, not based on the protocol layer. The default value is less than 0, and heartbeat detection is not enabled. If you need to enable it, it is more appropriate to set 60\_000, in milliseconds. (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |
| 公共属性 | [UseRegistrationPackage](8ad58263-7fbe-74f4-c963-5d60ede58a6b.htm) | 是否启用注册包的功能，如果不启用，就不会进行注册包的验证，有连接上来就自动覆盖本地的会话信息，不是很安全  Whether to enable the registration package function. If not enabled, the registration package verification will not be performed, and the local session information will be automatically overwritten when a connection is made, which is not very secure (继承自 [NetworkAlienClient](e29c23cc-75ae-02b5-2e55-00328c111f37.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [Dispose](a526b6f9-c988-d48c-5c7f-0d7fc3c69e77.htm) | (继承自 [NetworkAlienClient](e29c23cc-75ae-02b5-2e55-00328c111f37.htm)。) |
| 受保护的方法 | [Dispose(Boolean)](70f1c366-872f-04af-7560-1c9ef85be162.htm) | 释放当前的对象 (继承自 [NetworkAlienClient](e29c23cc-75ae-02b5-2e55-00328c111f37.htm)。) |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | [ExtraOnClose](6955ada4-316f-8bfb-eebe-d4d8da2c1271.htm) | 关闭的时候额外执行的功能代码 (重写 [CommunicationTcpServerExtraOnClose](2e788937-9ac2-73d5-2e82-6d8b58633902.htm).) |
| 受保护的方法 | [ExtraOnStart](2397069b-cc49-1ba5-41c6-ffe84e25219c.htm) | 服务器启动的时候额外执行的功能代码 (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | [GetDevices](dcba9753-9a7f-42b0-2200-e03576241c53.htm) | 获取所有的设备的信息，可以用来读写设备的数据信息  Get all device information, can be used to read and write device data information |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法 | [GetPipeSessions](8b0ec7b3-6291-9e63-6928-194ffdd03e61.htm) | 获取所有的会话信息，是否在线，上线的基本信息  Get all the session information, whether it is online, online basic information |
| 公共方法 | [GetTrustedClients](1a468d3b-7c50-7a36-f1f9-76a1d1488a8a.htm) | 获取受信任的客户端列表  Get a list of trusted clients (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |
| 公共方法 | GetType | (继承自 Object。) |
| 公共方法 | [IsClientOnline](9cdb7803-92ad-05ad-a793-f671f8f3124b.htm) | 检测当前的DTU是否在线 (重写 [NetworkAlienClientIsClientOnline(PipeDtuNet)](1b2b14cf-7c51-2c5a-754e-9190a072c9e1.htm).) |
| 受保护的方法 | [LogDebugMsg](3dc4b43f-2fdc-1d8e-ef40-66319b9f7055.htm) | 记录当前的日志信息 (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法 | [ServerClose](04a70d3a-d427-ed8d-53cf-3d3be56bdbd3.htm) | 关闭服务器的引擎  Shut down the server's engine (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |
| 公共方法 | [ServerStart](6aa1f94d-768c-4711-80e5-a529a5040f98.htm) | 指定端口号来启动服务器的引擎  Specify the port number to start the server's engine (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |
| 公共方法 | [ServerStart(Int32)](06e14b74-3deb-8e46-81c6-688f078375d9.htm) | 指定端口号来启动服务器的引擎  Specify the port number to start the server's engine (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |
| 公共方法 | [SetPassword](28933f98-e388-87d0-641c-235250b3eedb.htm) | 设置密码，需要传入长度为6的字节数组  To set the password, you need to pass in an array of bytes of length 6 (继承自 [NetworkAlienClient](e29c23cc-75ae-02b5-2e55-00328c111f37.htm)。) |
| 公共方法 | [SetSslPipeAction](5ecd4347-5dae-126d-f43c-76d4c7390705.htm) | 设置一个SSL的管道操作对象，在管道实例化之后，可以进行一些初始化的属性设置，例如自定义 SslProtocols 枚举 (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |
| 公共方法 | [SetTrustClients](4f6787a7-2c93-ff2a-bea2-d1ca09871957.htm) | 设置可信任的客户端列表，传入一个DTU的列表信息  Set up the list of trusted clients, passing in the list information for a DTU (继承自 [NetworkAlienClient](e29c23cc-75ae-02b5-2e55-00328c111f37.htm)。) |
| 公共方法 | [SetTrustedIpAddress](6204b3e7-af37-7c20-98db-08114431b5ac.htm) | 设置并启动受信任的客户端登录并读写，如果为null，将关闭对客户端的ip验证  Set and start the trusted client login and read and write, if it is null, the client's IP verification will be turned off (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |
| 受保护的方法 | [SocketAcceptExtraCheck](3ed24041-5b60-57a7-7629-f7c49baf8bf9.htm) | 当客户端的socket登录的时候额外检查的操作，并返回操作的结果信息。  The operation is additionally checked when the client's socket logs in, and the result information of the operation is returned. (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |
| 受保护的方法 | [ThreadPoolLogin](8ff83fc0-4b35-41ac-eeb1-2ab5d7a487b3.htm) | 当接收到了新的请求的时候执行的操作  An action performed when a new request is received (继承自 [NetworkAlienClient](e29c23cc-75ae-02b5-2e55-00328c111f37.htm)。) |
| 公共方法 | [ToString](523e19fd-a67a-db64-f7dd-144535ec6b0d.htm) | (重写 [NetworkAlienClientToString](99126697-d18d-bb3d-1899-c07f51ac2a96.htm).) |
| 公共方法 | [UseSSL(X509Certificate)](3577f66a-4337-6d6c-0f73-1f904f6bed38.htm) | 使用SSL通信，传递一个证书的对象 (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |
| 公共方法 | [UseSSL(String, String)](e6f37bac-c743-e05a-729a-0e4a64daf486.htm) | 使用SSL通信，传递一个证书的路径，以及证书的密码 (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)事件

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共事件 | [OnClientConnected](afbb5ac9-facf-52b4-2b7b-abd99733d1ca.htm) | 当有服务器连接上来的时候触发  Triggered when a server is connected (继承自 [NetworkAlienClient](e29c23cc-75ae-02b5-2e55-00328c111f37.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)字段

|  | 名称 | 说明 |
| --- | --- | --- |
| 受保护的字段 | [socketServer](f8c04282-bfb8-502e-730b-8d4127164ecb.htm) | 核心的socket服务器 (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)备注

针对异形客户端进行扩展信息

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.DTU 命名空间](bf8d39e1-bb7b-54e5-7d99-73cbdd1fc4af.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## DTUServer 构造函数 

[原文連結](http://api.hslcommunication.cn/html/bcb571b0-02d7-c103-7cfa-98305625d416.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.DTU](../html/bf8d39e1-bb7b-54e5-7d99-73cbdd1fc4af.htm "HslCommunication.DTU")

[DTUServer 类](../html/2e9a3fda-3f3e-82a7-01af-95a940d16410.htm "DTUServer 类")

[DTUServer 构造函数](../html/bcb571b0-02d7-c103-7cfa-98305625d416.htm "DTUServer 构造函数 ")

[DTUServer 构造函数 (List(DTUSettingType))](../html/493a5be9-6693-c4ba-b844-022368256626.htm "DTUServer 构造函数 (List(DTUSettingType))")

[DTUServer 构造函数 (String[], DeviceTcpNet[])](../html/02d9a717-15ad-7c11-6499-662cb3e7b303.htm "DTUServer 构造函数 (String[], DeviceTcpNet[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DTUServer 构造函数 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [DTUServer(ListDTUSettingType)](493a5be9-6693-c4ba-b844-022368256626.htm) | 根据配置的列表信息来实例化相关的DTU服务器  Instantiate the relevant DTU server according to the configured list information |
| 公共方法 | [DTUServer(String, DeviceTcpNet)](02d9a717-15ad-7c11-6499-662cb3e7b303.htm) | 根据配置的列表信息来实例化相关的DTU服务器  Instantiate the relevant DTU server according to the configured list information |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[DTUServer 类](2e9a3fda-3f3e-82a7-01af-95a940d16410.htm)

[HslCommunication.DTU 命名空间](bf8d39e1-bb7b-54e5-7d99-73cbdd1fc4af.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## DTUServer 构造函数 (List(DTUSettingType))

[原文連結](http://api.hslcommunication.cn/html/493a5be9-6693-c4ba-b844-022368256626.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.DTU](../html/bf8d39e1-bb7b-54e5-7d99-73cbdd1fc4af.htm "HslCommunication.DTU")

[DTUServer 类](../html/2e9a3fda-3f3e-82a7-01af-95a940d16410.htm "DTUServer 类")

[DTUServer 构造函数](../html/bcb571b0-02d7-c103-7cfa-98305625d416.htm "DTUServer 构造函数 ")

[DTUServer 构造函数 (List(DTUSettingType))](../html/493a5be9-6693-c4ba-b844-022368256626.htm "DTUServer 构造函数 (List(DTUSettingType))")

[DTUServer 构造函数 (String[], DeviceTcpNet[])](../html/02d9a717-15ad-7c11-6499-662cb3e7b303.htm "DTUServer 构造函数 (String[], DeviceTcpNet[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DTUServer 构造函数 (ListDTUSettingType) |

根据配置的列表信息来实例化相关的DTU服务器  
Instantiate the relevant DTU server according to the configured list information

**命名空间：**
 [HslCommunication.DTU](bf8d39e1-bb7b-54e5-7d99-73cbdd1fc4af.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public DTUServer(
	List<DTUSettingType> dTUSettings
)
```

```
Public Sub New ( 
	dTUSettings As List(Of DTUSettingType)
)
```

```
public:
DTUServer(
	List<DTUSettingType^>^ dTUSettings
)
```

```
new : 
        dTUSettings : List<DTUSettingType> -> DTUServer
```

#### 参数

dTUSettings
:   类型：System.Collections.GenericList[DTUSettingType](0a9eed32-45a3-d0eb-9104-fe18442a3351.htm)  
    DTU的配置信息

![](../icons/SectionExpanded.png)参见

#### 引用

[DTUServer 类](2e9a3fda-3f3e-82a7-01af-95a940d16410.htm)

[DTUServer 重载](bcb571b0-02d7-c103-7cfa-98305625d416.htm)

[HslCommunication.DTU 命名空间](bf8d39e1-bb7b-54e5-7d99-73cbdd1fc4af.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## DTUServer 构造函数 (String[], DeviceTcpNet[])

[原文連結](http://api.hslcommunication.cn/html/02d9a717-15ad-7c11-6499-662cb3e7b303.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.DTU](../html/bf8d39e1-bb7b-54e5-7d99-73cbdd1fc4af.htm "HslCommunication.DTU")

[DTUServer 类](../html/2e9a3fda-3f3e-82a7-01af-95a940d16410.htm "DTUServer 类")

[DTUServer 构造函数](../html/bcb571b0-02d7-c103-7cfa-98305625d416.htm "DTUServer 构造函数 ")

[DTUServer 构造函数 (List(DTUSettingType))](../html/493a5be9-6693-c4ba-b844-022368256626.htm "DTUServer 构造函数 (List(DTUSettingType))")

[DTUServer 构造函数 (String[], DeviceTcpNet[])](../html/02d9a717-15ad-7c11-6499-662cb3e7b303.htm "DTUServer 构造函数 (String[], DeviceTcpNet[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DTUServer 构造函数 (String, DeviceTcpNet) |

根据配置的列表信息来实例化相关的DTU服务器  
Instantiate the relevant DTU server according to the configured list information

**命名空间：**
 [HslCommunication.DTU](bf8d39e1-bb7b-54e5-7d99-73cbdd1fc4af.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public DTUServer(
	string[] dtuId,
	DeviceTcpNet[] networkDevices
)
```

```
Public Sub New ( 
	dtuId As String(),
	networkDevices As DeviceTcpNet()
)
```

```
public:
DTUServer(
	array<String^>^ dtuId, 
	array<DeviceTcpNet^>^ networkDevices
)
```

```
new : 
        dtuId : string[] * 
        networkDevices : DeviceTcpNet[] -> DTUServer
```

#### 参数

dtuId
:   类型：SystemString  
    Dtu信息

networkDevices
:   类型：[HslCommunication.Core.DeviceDeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)  
    设备信息

![](../icons/SectionExpanded.png)参见

#### 引用

[DTUServer 类](2e9a3fda-3f3e-82a7-01af-95a940d16410.htm)

[DTUServer 重载](bcb571b0-02d7-c103-7cfa-98305625d416.htm)

[HslCommunication.DTU 命名空间](bf8d39e1-bb7b-54e5-7d99-73cbdd1fc4af.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## DTUServer 属性

[原文連結](http://api.hslcommunication.cn/html/2307930e-d020-7b41-ceb5-a1e1f1ccaf78.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.DTU](../html/bf8d39e1-bb7b-54e5-7d99-73cbdd1fc4af.htm "HslCommunication.DTU")

[DTUServer 类](../html/2e9a3fda-3f3e-82a7-01af-95a940d16410.htm "DTUServer 类")

[DTUServer 属性](../html/2307930e-d020-7b41-ceb5-a1e1f1ccaf78.htm "DTUServer 属性")

[Item 属性](../html/a3462ddd-95fb-0430-a9fb-1f06244ca42a.htm "Item 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DTUServer 属性 |

[DTUServer](2e9a3fda-3f3e-82a7-01af-95a940d16410.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [EnableIPv6](f15cbcb3-5dcf-a658-0a93-3f0921383c32.htm) | 获取或设置服务器是否支持IPv6的地址协议信息  Get or set whether the server supports IPv6 address protocol information (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |
| 公共属性 | [IsCheckPwd](984bb24d-5f9e-5b3e-0a97-2860778c1e26.htm) | 是否统一检查密码，如果每个会话需要自己检查密码，就需要设置为false  Whether to check the password uniformly, if each session needs to check the password by itself, it needs to be set to false (继承自 [NetworkAlienClient](e29c23cc-75ae-02b5-2e55-00328c111f37.htm)。) |
| 公共属性 | [IsResponseAck](f2bdfc03-4e28-57b4-d777-de774e274744.htm) | 在DTU设备发送了注册报文的时候，指示是否返回响应报文，用来通知DTU设备是否登录成功，默认为 True  When a DTU sends a registration packet, it indicates whether to return a response packet to notify the DTU whether the DTU is logged in, The default is True (继承自 [NetworkAlienClient](e29c23cc-75ae-02b5-2e55-00328c111f37.htm)。) |
| 公共属性 | [IsStarted](8d6786a8-d416-4f45-158c-de6be0634aca.htm) | 服务器引擎是否启动  Whether the server engine is started (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |
| 公共属性 | [IsUseSSL](f4daa1e8-f433-0ee2-8284-16c4754eba25.htm) | 获取当前的服务器是否使用了SSL证书功能 (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |
| 公共属性 | [Item](a3462ddd-95fb-0430-a9fb-1f06244ca42a.htm) | 根据DTU信息获取设备的连接对象  Obtain the connection object of the device according to the DTU information |
| 公共属性 | [LocalAddress](c0effb49-cfe6-3267-9aef-afdab27a41c3.htm) | 获取或设置服务器绑定的本地IP地址，默认为空，使用本地所有可用的ip地址  Obtain or set the local IP address bound to the server, which is empty by default, and use all available IP addresses locally (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |
| 公共属性 | [LogDebugMessage](8a6db355-1130-2af8-cc6a-aeccc1ac0b29.htm) | 记录一些调试日志的委托，将会进行输出调试文本。  The delegate that records some debug logs will output debug text. (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |
| 公共属性代码示例 | [LogNet](746b148d-1ee3-1ee7-ee5d-3938eb72b2df.htm) | 组件的日志工具，支持日志记录，只要实例化后，当前网络的基本信息，就以[DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm)等级进行输出  The component's logging tool supports logging. As long as the instantiation of the basic network information, the output will be output at [DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm) (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |
| 公共属性 | [Port](6f700f7a-bd7e-c4b3-1275-a30491469e8f.htm) | 获取或设置服务器的端口号，如果是设置，需要在服务器启动前设置完成，才能生效。  Gets or sets the port number of the server. If it is set, it needs to be set before the server starts to take effect. (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |
| 公共属性 | [SocketKeepAliveTime](7d77e3df-9bbe-1a6e-a051-c87d9516c3a2.htm) | 获取或设置客户端的Socket的心跳时间信息，这个是Socket底层自动实现的心跳包，不基于协议层实现。默认小于0，不开启心跳检测，如果需要开启，设置 60\_000 比较合适，单位毫秒  Get or set the heartbeat time information of the Socket of the client. This is the heartbeat packet automatically implemented by the bottom layer of the Socket, not based on the protocol layer. The default value is less than 0, and heartbeat detection is not enabled. If you need to enable it, it is more appropriate to set 60\_000, in milliseconds. (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |
| 公共属性 | [UseRegistrationPackage](8ad58263-7fbe-74f4-c963-5d60ede58a6b.htm) | 是否启用注册包的功能，如果不启用，就不会进行注册包的验证，有连接上来就自动覆盖本地的会话信息，不是很安全  Whether to enable the registration package function. If not enabled, the registration package verification will not be performed, and the local session information will be automatically overwritten when a connection is made, which is not very secure (继承自 [NetworkAlienClient](e29c23cc-75ae-02b5-2e55-00328c111f37.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[DTUServer 类](2e9a3fda-3f3e-82a7-01af-95a940d16410.htm)

[HslCommunication.DTU 命名空间](bf8d39e1-bb7b-54e5-7d99-73cbdd1fc4af.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Item 属性 

[原文連結](http://api.hslcommunication.cn/html/a3462ddd-95fb-0430-a9fb-1f06244ca42a.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.DTU](../html/bf8d39e1-bb7b-54e5-7d99-73cbdd1fc4af.htm "HslCommunication.DTU")

[DTUServer 类](../html/2e9a3fda-3f3e-82a7-01af-95a940d16410.htm "DTUServer 类")

[DTUServer 属性](../html/2307930e-d020-7b41-ceb5-a1e1f1ccaf78.htm "DTUServer 属性")

[Item 属性](../html/a3462ddd-95fb-0430-a9fb-1f06244ca42a.htm "Item 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DTUServerItem 属性 |

根据DTU信息获取设备的连接对象  
Obtain the connection object of the device according to the DTU information

**命名空间：**
 [HslCommunication.DTU](bf8d39e1-bb7b-54e5-7d99-73cbdd1fc4af.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public DeviceCommunication this[
	string dtuId
] { get; }
```

```
Public ReadOnly Default Property Item ( 
	dtuId As String
) As DeviceCommunication
	Get
```

```
public:
property DeviceCommunication^ default[String^ dtuId] {
	DeviceCommunication^ get (String^ dtuId);
}
```

```
member Item : DeviceCommunication with get
```

#### 参数

dtuId
:   类型：SystemString  
    设备的id信息

#### 返回值

类型：[DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)  
设备的对象

![](../icons/SectionExpanded.png)参见

#### 引用

[DTUServer 类](2e9a3fda-3f3e-82a7-01af-95a940d16410.htm)

[HslCommunication.DTU 命名空间](bf8d39e1-bb7b-54e5-7d99-73cbdd1fc4af.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## DTUServer 方法

[原文連結](http://api.hslcommunication.cn/html/bd9fa841-c520-b11a-85b8-6f4ad77db379.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.DTU](../html/bf8d39e1-bb7b-54e5-7d99-73cbdd1fc4af.htm "HslCommunication.DTU")

[DTUServer 类](../html/2e9a3fda-3f3e-82a7-01af-95a940d16410.htm "DTUServer 类")

[DTUServer 方法](../html/bd9fa841-c520-b11a-85b8-6f4ad77db379.htm "DTUServer 方法")

[ExtraOnClose 方法](../html/6955ada4-316f-8bfb-eebe-d4d8da2c1271.htm "ExtraOnClose 方法 ")

[GetDevices 方法](../html/dcba9753-9a7f-42b0-2200-e03576241c53.htm "GetDevices 方法 ")

[GetPipeSessions 方法](../html/8b0ec7b3-6291-9e63-6928-194ffdd03e61.htm "GetPipeSessions 方法 ")

[IsClientOnline 方法](../html/9cdb7803-92ad-05ad-a793-f671f8f3124b.htm "IsClientOnline 方法 ")

[ToString 方法](../html/523e19fd-a67a-db64-f7dd-144535ec6b0d.htm "ToString 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DTUServer 方法 |

[DTUServer](2e9a3fda-3f3e-82a7-01af-95a940d16410.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [Dispose](a526b6f9-c988-d48c-5c7f-0d7fc3c69e77.htm) | (继承自 [NetworkAlienClient](e29c23cc-75ae-02b5-2e55-00328c111f37.htm)。) |
| 受保护的方法 | [Dispose(Boolean)](70f1c366-872f-04af-7560-1c9ef85be162.htm) | 释放当前的对象 (继承自 [NetworkAlienClient](e29c23cc-75ae-02b5-2e55-00328c111f37.htm)。) |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | [ExtraOnClose](6955ada4-316f-8bfb-eebe-d4d8da2c1271.htm) | 关闭的时候额外执行的功能代码 (重写 [CommunicationTcpServerExtraOnClose](2e788937-9ac2-73d5-2e82-6d8b58633902.htm).) |
| 受保护的方法 | [ExtraOnStart](2397069b-cc49-1ba5-41c6-ffe84e25219c.htm) | 服务器启动的时候额外执行的功能代码 (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | [GetDevices](dcba9753-9a7f-42b0-2200-e03576241c53.htm) | 获取所有的设备的信息，可以用来读写设备的数据信息  Get all device information, can be used to read and write device data information |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法 | [GetPipeSessions](8b0ec7b3-6291-9e63-6928-194ffdd03e61.htm) | 获取所有的会话信息，是否在线，上线的基本信息  Get all the session information, whether it is online, online basic information |
| 公共方法 | [GetTrustedClients](1a468d3b-7c50-7a36-f1f9-76a1d1488a8a.htm) | 获取受信任的客户端列表  Get a list of trusted clients (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |
| 公共方法 | GetType | (继承自 Object。) |
| 公共方法 | [IsClientOnline](9cdb7803-92ad-05ad-a793-f671f8f3124b.htm) | 检测当前的DTU是否在线 (重写 [NetworkAlienClientIsClientOnline(PipeDtuNet)](1b2b14cf-7c51-2c5a-754e-9190a072c9e1.htm).) |
| 受保护的方法 | [LogDebugMsg](3dc4b43f-2fdc-1d8e-ef40-66319b9f7055.htm) | 记录当前的日志信息 (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法 | [ServerClose](04a70d3a-d427-ed8d-53cf-3d3be56bdbd3.htm) | 关闭服务器的引擎  Shut down the server's engine (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |
| 公共方法 | [ServerStart](6aa1f94d-768c-4711-80e5-a529a5040f98.htm) | 指定端口号来启动服务器的引擎  Specify the port number to start the server's engine (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |
| 公共方法 | [ServerStart(Int32)](06e14b74-3deb-8e46-81c6-688f078375d9.htm) | 指定端口号来启动服务器的引擎  Specify the port number to start the server's engine (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |
| 公共方法 | [SetPassword](28933f98-e388-87d0-641c-235250b3eedb.htm) | 设置密码，需要传入长度为6的字节数组  To set the password, you need to pass in an array of bytes of length 6 (继承自 [NetworkAlienClient](e29c23cc-75ae-02b5-2e55-00328c111f37.htm)。) |
| 公共方法 | [SetSslPipeAction](5ecd4347-5dae-126d-f43c-76d4c7390705.htm) | 设置一个SSL的管道操作对象，在管道实例化之后，可以进行一些初始化的属性设置，例如自定义 SslProtocols 枚举 (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |
| 公共方法 | [SetTrustClients](4f6787a7-2c93-ff2a-bea2-d1ca09871957.htm) | 设置可信任的客户端列表，传入一个DTU的列表信息  Set up the list of trusted clients, passing in the list information for a DTU (继承自 [NetworkAlienClient](e29c23cc-75ae-02b5-2e55-00328c111f37.htm)。) |
| 公共方法 | [SetTrustedIpAddress](6204b3e7-af37-7c20-98db-08114431b5ac.htm) | 设置并启动受信任的客户端登录并读写，如果为null，将关闭对客户端的ip验证  Set and start the trusted client login and read and write, if it is null, the client's IP verification will be turned off (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |
| 受保护的方法 | [SocketAcceptExtraCheck](3ed24041-5b60-57a7-7629-f7c49baf8bf9.htm) | 当客户端的socket登录的时候额外检查的操作，并返回操作的结果信息。  The operation is additionally checked when the client's socket logs in, and the result information of the operation is returned. (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |
| 受保护的方法 | [ThreadPoolLogin](8ff83fc0-4b35-41ac-eeb1-2ab5d7a487b3.htm) | 当接收到了新的请求的时候执行的操作  An action performed when a new request is received (继承自 [NetworkAlienClient](e29c23cc-75ae-02b5-2e55-00328c111f37.htm)。) |
| 公共方法 | [ToString](523e19fd-a67a-db64-f7dd-144535ec6b0d.htm) | (重写 [NetworkAlienClientToString](99126697-d18d-bb3d-1899-c07f51ac2a96.htm).) |
| 公共方法 | [UseSSL(X509Certificate)](3577f66a-4337-6d6c-0f73-1f904f6bed38.htm) | 使用SSL通信，传递一个证书的对象 (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |
| 公共方法 | [UseSSL(String, String)](e6f37bac-c743-e05a-729a-0e4a64daf486.htm) | 使用SSL通信，传递一个证书的路径，以及证书的密码 (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[DTUServer 类](2e9a3fda-3f3e-82a7-01af-95a940d16410.htm)

[HslCommunication.DTU 命名空间](bf8d39e1-bb7b-54e5-7d99-73cbdd1fc4af.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ExtraOnClose 方法 

[原文連結](http://api.hslcommunication.cn/html/6955ada4-316f-8bfb-eebe-d4d8da2c1271.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.DTU](../html/bf8d39e1-bb7b-54e5-7d99-73cbdd1fc4af.htm "HslCommunication.DTU")

[DTUServer 类](../html/2e9a3fda-3f3e-82a7-01af-95a940d16410.htm "DTUServer 类")

[DTUServer 方法](../html/bd9fa841-c520-b11a-85b8-6f4ad77db379.htm "DTUServer 方法")

[ExtraOnClose 方法](../html/6955ada4-316f-8bfb-eebe-d4d8da2c1271.htm "ExtraOnClose 方法 ")

[GetDevices 方法](../html/dcba9753-9a7f-42b0-2200-e03576241c53.htm "GetDevices 方法 ")

[GetPipeSessions 方法](../html/8b0ec7b3-6291-9e63-6928-194ffdd03e61.htm "GetPipeSessions 方法 ")

[IsClientOnline 方法](../html/9cdb7803-92ad-05ad-a793-f671f8f3124b.htm "IsClientOnline 方法 ")

[ToString 方法](../html/523e19fd-a67a-db64-f7dd-144535ec6b0d.htm "ToString 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DTUServerExtraOnClose 方法 |

关闭的时候额外执行的功能代码

**命名空间：**
 [HslCommunication.DTU](bf8d39e1-bb7b-54e5-7d99-73cbdd1fc4af.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
protected override void ExtraOnClose()
```

```
Protected Overrides Sub ExtraOnClose
```

```
protected:
virtual void ExtraOnClose() override
```

```
abstract ExtraOnClose : unit -> unit 
override ExtraOnClose : unit -> unit
```

![](../icons/SectionExpanded.png)参见

#### 引用

[DTUServer 类](2e9a3fda-3f3e-82a7-01af-95a940d16410.htm)

[HslCommunication.DTU 命名空间](bf8d39e1-bb7b-54e5-7d99-73cbdd1fc4af.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GetDevices 方法 

[原文連結](http://api.hslcommunication.cn/html/dcba9753-9a7f-42b0-2200-e03576241c53.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.DTU](../html/bf8d39e1-bb7b-54e5-7d99-73cbdd1fc4af.htm "HslCommunication.DTU")

[DTUServer 类](../html/2e9a3fda-3f3e-82a7-01af-95a940d16410.htm "DTUServer 类")

[DTUServer 方法](../html/bd9fa841-c520-b11a-85b8-6f4ad77db379.htm "DTUServer 方法")

[ExtraOnClose 方法](../html/6955ada4-316f-8bfb-eebe-d4d8da2c1271.htm "ExtraOnClose 方法 ")

[GetDevices 方法](../html/dcba9753-9a7f-42b0-2200-e03576241c53.htm "GetDevices 方法 ")

[GetPipeSessions 方法](../html/8b0ec7b3-6291-9e63-6928-194ffdd03e61.htm "GetPipeSessions 方法 ")

[IsClientOnline 方法](../html/9cdb7803-92ad-05ad-a793-f671f8f3124b.htm "IsClientOnline 方法 ")

[ToString 方法](../html/523e19fd-a67a-db64-f7dd-144535ec6b0d.htm "ToString 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DTUServerGetDevices 方法 |

获取所有的设备的信息，可以用来读写设备的数据信息  
Get all device information, can be used to read and write device data information

**命名空间：**
 [HslCommunication.DTU](bf8d39e1-bb7b-54e5-7d99-73cbdd1fc4af.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public DeviceCommunication[] GetDevices()
```

```
Public Function GetDevices As DeviceCommunication()
```

```
public:
array<DeviceCommunication^>^ GetDevices()
```

```
member GetDevices : unit -> DeviceCommunication[] 
```

#### 返回值

类型：[DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)  
设备数组

![](../icons/SectionExpanded.png)参见

#### 引用

[DTUServer 类](2e9a3fda-3f3e-82a7-01af-95a940d16410.htm)

[HslCommunication.DTU 命名空间](bf8d39e1-bb7b-54e5-7d99-73cbdd1fc4af.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GetPipeSessions 方法 

[原文連結](http://api.hslcommunication.cn/html/8b0ec7b3-6291-9e63-6928-194ffdd03e61.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.DTU](../html/bf8d39e1-bb7b-54e5-7d99-73cbdd1fc4af.htm "HslCommunication.DTU")

[DTUServer 类](../html/2e9a3fda-3f3e-82a7-01af-95a940d16410.htm "DTUServer 类")

[DTUServer 方法](../html/bd9fa841-c520-b11a-85b8-6f4ad77db379.htm "DTUServer 方法")

[ExtraOnClose 方法](../html/6955ada4-316f-8bfb-eebe-d4d8da2c1271.htm "ExtraOnClose 方法 ")

[GetDevices 方法](../html/dcba9753-9a7f-42b0-2200-e03576241c53.htm "GetDevices 方法 ")

[GetPipeSessions 方法](../html/8b0ec7b3-6291-9e63-6928-194ffdd03e61.htm "GetPipeSessions 方法 ")

[IsClientOnline 方法](../html/9cdb7803-92ad-05ad-a793-f671f8f3124b.htm "IsClientOnline 方法 ")

[ToString 方法](../html/523e19fd-a67a-db64-f7dd-144535ec6b0d.htm "ToString 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DTUServerGetPipeSessions 方法 |

获取所有的会话信息，是否在线，上线的基本信息  
Get all the session information, whether it is online, online basic information

**命名空间：**
 [HslCommunication.DTU](bf8d39e1-bb7b-54e5-7d99-73cbdd1fc4af.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public PipeDtuNet[] GetPipeSessions()
```

```
Public Function GetPipeSessions As PipeDtuNet()
```

```
public:
array<PipeDtuNet^>^ GetPipeSessions()
```

```
member GetPipeSessions : unit -> PipeDtuNet[] 
```

#### 返回值

类型：[PipeDtuNet](97793361-a796-5f13-5029-1abdb2bff247.htm)  
会话列表

![](../icons/SectionExpanded.png)参见

#### 引用

[DTUServer 类](2e9a3fda-3f3e-82a7-01af-95a940d16410.htm)

[HslCommunication.DTU 命名空间](bf8d39e1-bb7b-54e5-7d99-73cbdd1fc4af.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## IsClientOnline 方法 

[原文連結](http://api.hslcommunication.cn/html/9cdb7803-92ad-05ad-a793-f671f8f3124b.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.DTU](../html/bf8d39e1-bb7b-54e5-7d99-73cbdd1fc4af.htm "HslCommunication.DTU")

[DTUServer 类](../html/2e9a3fda-3f3e-82a7-01af-95a940d16410.htm "DTUServer 类")

[DTUServer 方法](../html/bd9fa841-c520-b11a-85b8-6f4ad77db379.htm "DTUServer 方法")

[ExtraOnClose 方法](../html/6955ada4-316f-8bfb-eebe-d4d8da2c1271.htm "ExtraOnClose 方法 ")

[GetDevices 方法](../html/dcba9753-9a7f-42b0-2200-e03576241c53.htm "GetDevices 方法 ")

[GetPipeSessions 方法](../html/8b0ec7b3-6291-9e63-6928-194ffdd03e61.htm "GetPipeSessions 方法 ")

[IsClientOnline 方法](../html/9cdb7803-92ad-05ad-a793-f671f8f3124b.htm "IsClientOnline 方法 ")

[ToString 方法](../html/523e19fd-a67a-db64-f7dd-144535ec6b0d.htm "ToString 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DTUServerIsClientOnline 方法 |

检测当前的DTU是否在线

**命名空间：**
 [HslCommunication.DTU](bf8d39e1-bb7b-54e5-7d99-73cbdd1fc4af.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public override int IsClientOnline(
	PipeDtuNet pipe
)
```

```
Public Overrides Function IsClientOnline ( 
	pipe As PipeDtuNet
) As Integer
```

```
public:
virtual int IsClientOnline(
	PipeDtuNet^ pipe
) override
```

```
abstract IsClientOnline : 
        pipe : PipeDtuNet -> int 
override IsClientOnline : 
        pipe : PipeDtuNet -> int
```

#### 参数

pipe
:   类型：[HslCommunication.Core.PipePipeDtuNet](97793361-a796-5f13-5029-1abdb2bff247.htm)  
    当前的会话信息

#### 返回值

类型：Int32  
当前的会话是否在线

![](../icons/SectionExpanded.png)参见

#### 引用

[DTUServer 类](2e9a3fda-3f3e-82a7-01af-95a940d16410.htm)

[HslCommunication.DTU 命名空间](bf8d39e1-bb7b-54e5-7d99-73cbdd1fc4af.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ToString 方法 

[原文連結](http://api.hslcommunication.cn/html/523e19fd-a67a-db64-f7dd-144535ec6b0d.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.DTU](../html/bf8d39e1-bb7b-54e5-7d99-73cbdd1fc4af.htm "HslCommunication.DTU")

[DTUServer 类](../html/2e9a3fda-3f3e-82a7-01af-95a940d16410.htm "DTUServer 类")

[DTUServer 方法](../html/bd9fa841-c520-b11a-85b8-6f4ad77db379.htm "DTUServer 方法")

[ExtraOnClose 方法](../html/6955ada4-316f-8bfb-eebe-d4d8da2c1271.htm "ExtraOnClose 方法 ")

[GetDevices 方法](../html/dcba9753-9a7f-42b0-2200-e03576241c53.htm "GetDevices 方法 ")

[GetPipeSessions 方法](../html/8b0ec7b3-6291-9e63-6928-194ffdd03e61.htm "GetPipeSessions 方法 ")

[IsClientOnline 方法](../html/9cdb7803-92ad-05ad-a793-f671f8f3124b.htm "IsClientOnline 方法 ")

[ToString 方法](../html/523e19fd-a67a-db64-f7dd-144535ec6b0d.htm "ToString 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DTUServerToString 方法 |

[缺少 "M:HslCommunication.DTU.DTUServer.ToString" 的 <summary> 文档]

**命名空间：**
 [HslCommunication.DTU](bf8d39e1-bb7b-54e5-7d99-73cbdd1fc4af.htm)  
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

[缺少 "M:HslCommunication.DTU.DTUServer.ToString" 的 <returns> 文档]

![](../icons/SectionExpanded.png)参见

#### 引用

[DTUServer 类](2e9a3fda-3f3e-82a7-01af-95a940d16410.htm)

[HslCommunication.DTU 命名空间](bf8d39e1-bb7b-54e5-7d99-73cbdd1fc4af.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## DTUServer 事件

[原文連結](http://api.hslcommunication.cn/html/38005bd6-daa6-8120-950f-ebf087418ac9.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.DTU](../html/bf8d39e1-bb7b-54e5-7d99-73cbdd1fc4af.htm "HslCommunication.DTU")

[DTUServer 类](../html/2e9a3fda-3f3e-82a7-01af-95a940d16410.htm "DTUServer 类")

[DTUServer 构造函数](../html/bcb571b0-02d7-c103-7cfa-98305625d416.htm "DTUServer 构造函数 ")

[DTUServer 属性](../html/2307930e-d020-7b41-ceb5-a1e1f1ccaf78.htm "DTUServer 属性")

[DTUServer 方法](../html/bd9fa841-c520-b11a-85b8-6f4ad77db379.htm "DTUServer 方法")

[DTUServer 事件](../html/38005bd6-daa6-8120-950f-ebf087418ac9.htm "DTUServer 事件")

[DTUServer 字段](../html/b75cd2c5-4594-0645-d74c-d0eba36b388c.htm "DTUServer 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DTUServer 事件 |

[DTUServer](2e9a3fda-3f3e-82a7-01af-95a940d16410.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)事件

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共事件 | [OnClientConnected](afbb5ac9-facf-52b4-2b7b-abd99733d1ca.htm) | 当有服务器连接上来的时候触发  Triggered when a server is connected (继承自 [NetworkAlienClient](e29c23cc-75ae-02b5-2e55-00328c111f37.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[DTUServer 类](2e9a3fda-3f3e-82a7-01af-95a940d16410.htm)

[HslCommunication.DTU 命名空间](bf8d39e1-bb7b-54e5-7d99-73cbdd1fc4af.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## DTUServer 字段

[原文連結](http://api.hslcommunication.cn/html/b75cd2c5-4594-0645-d74c-d0eba36b388c.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.DTU](../html/bf8d39e1-bb7b-54e5-7d99-73cbdd1fc4af.htm "HslCommunication.DTU")

[DTUServer 类](../html/2e9a3fda-3f3e-82a7-01af-95a940d16410.htm "DTUServer 类")

[DTUServer 构造函数](../html/bcb571b0-02d7-c103-7cfa-98305625d416.htm "DTUServer 构造函数 ")

[DTUServer 属性](../html/2307930e-d020-7b41-ceb5-a1e1f1ccaf78.htm "DTUServer 属性")

[DTUServer 方法](../html/bd9fa841-c520-b11a-85b8-6f4ad77db379.htm "DTUServer 方法")

[DTUServer 事件](../html/38005bd6-daa6-8120-950f-ebf087418ac9.htm "DTUServer 事件")

[DTUServer 字段](../html/b75cd2c5-4594-0645-d74c-d0eba36b388c.htm "DTUServer 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DTUServer 字段 |

[DTUServer](2e9a3fda-3f3e-82a7-01af-95a940d16410.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)字段

|  | 名称 | 说明 |
| --- | --- | --- |
| 受保护的字段 | [socketServer](f8c04282-bfb8-502e-730b-8d4127164ecb.htm) | 核心的socket服务器 (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[DTUServer 类](2e9a3fda-3f3e-82a7-01af-95a940d16410.htm)

[HslCommunication.DTU 命名空间](bf8d39e1-bb7b-54e5-7d99-73cbdd1fc4af.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## DTUSettingType 类

[原文連結](http://api.hslcommunication.cn/html/0a9eed32-45a3-d0eb-9104-fe18442a3351.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.DTU](../html/bf8d39e1-bb7b-54e5-7d99-73cbdd1fc4af.htm "HslCommunication.DTU")

[DTUSettingType 类](../html/0a9eed32-45a3-d0eb-9104-fe18442a3351.htm "DTUSettingType 类")

[DTUSettingType 构造函数](../html/d7f0e42f-2cd7-bf60-ef0c-bee1af02ce82.htm "DTUSettingType 构造函数 ")

[DTUSettingType 属性](../html/ad5183a9-6f24-d32d-8e77-2971a7a39409.htm "DTUSettingType 属性")

[DTUSettingType 方法](../html/740e1d50-2e55-5f2b-f263-068383a8dbfd.htm "DTUSettingType 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DTUSettingType 类 |

DTU的类型设置器

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  HslCommunication.DTUDTUSettingType

**命名空间：**
 [HslCommunication.DTU](bf8d39e1-bb7b-54e5-7d99-73cbdd1fc4af.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class DTUSettingType
```

```
Public Class DTUSettingType
```

```
public ref class DTUSettingType
```

```
type DTUSettingType =  class end
```

DTUSettingType 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [DTUSettingType](d7f0e42f-2cd7-bf60-ef0c-bee1af02ce82.htm) | 初始化 DTUSettingType 类的一个新实例 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [DtuId](18d55052-f440-e5c5-140d-3351be3149ba.htm) | 设备的唯一ID信息 |
| 公共属性 | [DtuType](3c155da9-4a86-8ebb-74f2-b5bdddfd5152.htm) | 当前的设备的类型 |
| 公共属性 | [JsonParameter](8c0300dd-5899-ae59-203c-0f85f346058b.htm) | 额外的参数都存放在json里面 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | [GetClient](6b84fa80-15c5-8223-789f-173b8ef65c4d.htm) | 根据类型，获取连接对象 |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法 | [ToString](6974f5c4-cea4-7956-15e3-11a810c9757f.htm) | (重写 ObjectToString.) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.DTU 命名空间](bf8d39e1-bb7b-54e5-7d99-73cbdd1fc4af.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## DTUSettingType 构造函数 

[原文連結](http://api.hslcommunication.cn/html/d7f0e42f-2cd7-bf60-ef0c-bee1af02ce82.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.DTU](../html/bf8d39e1-bb7b-54e5-7d99-73cbdd1fc4af.htm "HslCommunication.DTU")

[DTUSettingType 类](../html/0a9eed32-45a3-d0eb-9104-fe18442a3351.htm "DTUSettingType 类")

[DTUSettingType 构造函数](../html/d7f0e42f-2cd7-bf60-ef0c-bee1af02ce82.htm "DTUSettingType 构造函数 ")

[DTUSettingType 属性](../html/ad5183a9-6f24-d32d-8e77-2971a7a39409.htm "DTUSettingType 属性")

[DTUSettingType 方法](../html/740e1d50-2e55-5f2b-f263-068383a8dbfd.htm "DTUSettingType 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DTUSettingType 构造函数 |

初始化 [DTUSettingType](0a9eed32-45a3-d0eb-9104-fe18442a3351.htm) 类的一个新实例

**命名空间：**
 [HslCommunication.DTU](bf8d39e1-bb7b-54e5-7d99-73cbdd1fc4af.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public DTUSettingType()
```

```
Public Sub New
```

```
public:
DTUSettingType()
```

```
new : unit -> DTUSettingType
```

![](../icons/SectionExpanded.png)参见

#### 引用

[DTUSettingType 类](0a9eed32-45a3-d0eb-9104-fe18442a3351.htm)

[HslCommunication.DTU 命名空间](bf8d39e1-bb7b-54e5-7d99-73cbdd1fc4af.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## DTUSettingType 属性

[原文連結](http://api.hslcommunication.cn/html/ad5183a9-6f24-d32d-8e77-2971a7a39409.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.DTU](../html/bf8d39e1-bb7b-54e5-7d99-73cbdd1fc4af.htm "HslCommunication.DTU")

[DTUSettingType 类](../html/0a9eed32-45a3-d0eb-9104-fe18442a3351.htm "DTUSettingType 类")

[DTUSettingType 属性](../html/ad5183a9-6f24-d32d-8e77-2971a7a39409.htm "DTUSettingType 属性")

[DtuId 属性](../html/18d55052-f440-e5c5-140d-3351be3149ba.htm "DtuId 属性 ")

[DtuType 属性](../html/3c155da9-4a86-8ebb-74f2-b5bdddfd5152.htm "DtuType 属性 ")

[JsonParameter 属性](../html/8c0300dd-5899-ae59-203c-0f85f346058b.htm "JsonParameter 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DTUSettingType 属性 |

[DTUSettingType](0a9eed32-45a3-d0eb-9104-fe18442a3351.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [DtuId](18d55052-f440-e5c5-140d-3351be3149ba.htm) | 设备的唯一ID信息 |
| 公共属性 | [DtuType](3c155da9-4a86-8ebb-74f2-b5bdddfd5152.htm) | 当前的设备的类型 |
| 公共属性 | [JsonParameter](8c0300dd-5899-ae59-203c-0f85f346058b.htm) | 额外的参数都存放在json里面 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[DTUSettingType 类](0a9eed32-45a3-d0eb-9104-fe18442a3351.htm)

[HslCommunication.DTU 命名空间](bf8d39e1-bb7b-54e5-7d99-73cbdd1fc4af.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## DtuId 属性 

[原文連結](http://api.hslcommunication.cn/html/18d55052-f440-e5c5-140d-3351be3149ba.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.DTU](../html/bf8d39e1-bb7b-54e5-7d99-73cbdd1fc4af.htm "HslCommunication.DTU")

[DTUSettingType 类](../html/0a9eed32-45a3-d0eb-9104-fe18442a3351.htm "DTUSettingType 类")

[DTUSettingType 属性](../html/ad5183a9-6f24-d32d-8e77-2971a7a39409.htm "DTUSettingType 属性")

[DtuId 属性](../html/18d55052-f440-e5c5-140d-3351be3149ba.htm "DtuId 属性 ")

[DtuType 属性](../html/3c155da9-4a86-8ebb-74f2-b5bdddfd5152.htm "DtuType 属性 ")

[JsonParameter 属性](../html/8c0300dd-5899-ae59-203c-0f85f346058b.htm "JsonParameter 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DTUSettingTypeDtuId 属性 |

设备的唯一ID信息

**命名空间：**
 [HslCommunication.DTU](bf8d39e1-bb7b-54e5-7d99-73cbdd1fc4af.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public string DtuId { get; set; }
```

```
Public Property DtuId As String
	Get
	Set
```

```
public:
property String^ DtuId {
	String^ get ();
	void set (String^ value);
}
```

```
member DtuId : string with get, set
```

#### 属性值

类型：String

![](../icons/SectionExpanded.png)参见

#### 引用

[DTUSettingType 类](0a9eed32-45a3-d0eb-9104-fe18442a3351.htm)

[HslCommunication.DTU 命名空间](bf8d39e1-bb7b-54e5-7d99-73cbdd1fc4af.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## DtuType 属性 

[原文連結](http://api.hslcommunication.cn/html/3c155da9-4a86-8ebb-74f2-b5bdddfd5152.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.DTU](../html/bf8d39e1-bb7b-54e5-7d99-73cbdd1fc4af.htm "HslCommunication.DTU")

[DTUSettingType 类](../html/0a9eed32-45a3-d0eb-9104-fe18442a3351.htm "DTUSettingType 类")

[DTUSettingType 属性](../html/ad5183a9-6f24-d32d-8e77-2971a7a39409.htm "DTUSettingType 属性")

[DtuId 属性](../html/18d55052-f440-e5c5-140d-3351be3149ba.htm "DtuId 属性 ")

[DtuType 属性](../html/3c155da9-4a86-8ebb-74f2-b5bdddfd5152.htm "DtuType 属性 ")

[JsonParameter 属性](../html/8c0300dd-5899-ae59-203c-0f85f346058b.htm "JsonParameter 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DTUSettingTypeDtuType 属性 |

当前的设备的类型

**命名空间：**
 [HslCommunication.DTU](bf8d39e1-bb7b-54e5-7d99-73cbdd1fc4af.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public string DtuType { get; set; }
```

```
Public Property DtuType As String
	Get
	Set
```

```
public:
property String^ DtuType {
	String^ get ();
	void set (String^ value);
}
```

```
member DtuType : string with get, set
```

#### 属性值

类型：String

![](../icons/SectionExpanded.png)参见

#### 引用

[DTUSettingType 类](0a9eed32-45a3-d0eb-9104-fe18442a3351.htm)

[HslCommunication.DTU 命名空间](bf8d39e1-bb7b-54e5-7d99-73cbdd1fc4af.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## JsonParameter 属性 

[原文連結](http://api.hslcommunication.cn/html/8c0300dd-5899-ae59-203c-0f85f346058b.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.DTU](../html/bf8d39e1-bb7b-54e5-7d99-73cbdd1fc4af.htm "HslCommunication.DTU")

[DTUSettingType 类](../html/0a9eed32-45a3-d0eb-9104-fe18442a3351.htm "DTUSettingType 类")

[DTUSettingType 属性](../html/ad5183a9-6f24-d32d-8e77-2971a7a39409.htm "DTUSettingType 属性")

[DtuId 属性](../html/18d55052-f440-e5c5-140d-3351be3149ba.htm "DtuId 属性 ")

[DtuType 属性](../html/3c155da9-4a86-8ebb-74f2-b5bdddfd5152.htm "DtuType 属性 ")

[JsonParameter 属性](../html/8c0300dd-5899-ae59-203c-0f85f346058b.htm "JsonParameter 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DTUSettingTypeJsonParameter 属性 |

额外的参数都存放在json里面

**命名空间：**
 [HslCommunication.DTU](bf8d39e1-bb7b-54e5-7d99-73cbdd1fc4af.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public string JsonParameter { get; set; }
```

```
Public Property JsonParameter As String
	Get
	Set
```

```
public:
property String^ JsonParameter {
	String^ get ();
	void set (String^ value);
}
```

```
member JsonParameter : string with get, set
```

#### 属性值

类型：String

![](../icons/SectionExpanded.png)参见

#### 引用

[DTUSettingType 类](0a9eed32-45a3-d0eb-9104-fe18442a3351.htm)

[HslCommunication.DTU 命名空间](bf8d39e1-bb7b-54e5-7d99-73cbdd1fc4af.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## DTUSettingType 方法

[原文連結](http://api.hslcommunication.cn/html/740e1d50-2e55-5f2b-f263-068383a8dbfd.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.DTU](../html/bf8d39e1-bb7b-54e5-7d99-73cbdd1fc4af.htm "HslCommunication.DTU")

[DTUSettingType 类](../html/0a9eed32-45a3-d0eb-9104-fe18442a3351.htm "DTUSettingType 类")

[DTUSettingType 方法](../html/740e1d50-2e55-5f2b-f263-068383a8dbfd.htm "DTUSettingType 方法")

[GetClient 方法](../html/6b84fa80-15c5-8223-789f-173b8ef65c4d.htm "GetClient 方法 ")

[ToString 方法](../html/6974f5c4-cea4-7956-15e3-11a810c9757f.htm "ToString 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DTUSettingType 方法 |

[DTUSettingType](0a9eed32-45a3-d0eb-9104-fe18442a3351.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | [GetClient](6b84fa80-15c5-8223-789f-173b8ef65c4d.htm) | 根据类型，获取连接对象 |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法 | [ToString](6974f5c4-cea4-7956-15e3-11a810c9757f.htm) | (重写 ObjectToString.) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[DTUSettingType 类](0a9eed32-45a3-d0eb-9104-fe18442a3351.htm)

[HslCommunication.DTU 命名空间](bf8d39e1-bb7b-54e5-7d99-73cbdd1fc4af.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GetClient 方法 

[原文連結](http://api.hslcommunication.cn/html/6b84fa80-15c5-8223-789f-173b8ef65c4d.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.DTU](../html/bf8d39e1-bb7b-54e5-7d99-73cbdd1fc4af.htm "HslCommunication.DTU")

[DTUSettingType 类](../html/0a9eed32-45a3-d0eb-9104-fe18442a3351.htm "DTUSettingType 类")

[DTUSettingType 方法](../html/740e1d50-2e55-5f2b-f263-068383a8dbfd.htm "DTUSettingType 方法")

[GetClient 方法](../html/6b84fa80-15c5-8223-789f-173b8ef65c4d.htm "GetClient 方法 ")

[ToString 方法](../html/6974f5c4-cea4-7956-15e3-11a810c9757f.htm "ToString 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DTUSettingTypeGetClient 方法 |

根据类型，获取连接对象

**命名空间：**
 [HslCommunication.DTU](bf8d39e1-bb7b-54e5-7d99-73cbdd1fc4af.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public virtual DeviceTcpNet GetClient()
```

```
Public Overridable Function GetClient As DeviceTcpNet
```

```
public:
virtual DeviceTcpNet^ GetClient()
```

```
abstract GetClient : unit -> DeviceTcpNet 
override GetClient : unit -> DeviceTcpNet
```

#### 返回值

类型：[DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)  
获取设备的连接对象

![](../icons/SectionExpanded.png)参见

#### 引用

[DTUSettingType 类](0a9eed32-45a3-d0eb-9104-fe18442a3351.htm)

[HslCommunication.DTU 命名空间](bf8d39e1-bb7b-54e5-7d99-73cbdd1fc4af.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ToString 方法 

[原文連結](http://api.hslcommunication.cn/html/6974f5c4-cea4-7956-15e3-11a810c9757f.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.DTU](../html/bf8d39e1-bb7b-54e5-7d99-73cbdd1fc4af.htm "HslCommunication.DTU")

[DTUSettingType 类](../html/0a9eed32-45a3-d0eb-9104-fe18442a3351.htm "DTUSettingType 类")

[DTUSettingType 方法](../html/740e1d50-2e55-5f2b-f263-068383a8dbfd.htm "DTUSettingType 方法")

[GetClient 方法](../html/6b84fa80-15c5-8223-789f-173b8ef65c4d.htm "GetClient 方法 ")

[ToString 方法](../html/6974f5c4-cea4-7956-15e3-11a810c9757f.htm "ToString 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DTUSettingTypeToString 方法 |

[缺少 "M:HslCommunication.DTU.DTUSettingType.ToString" 的 <summary> 文档]

**命名空间：**
 [HslCommunication.DTU](bf8d39e1-bb7b-54e5-7d99-73cbdd1fc4af.htm)  
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

[缺少 "M:HslCommunication.DTU.DTUSettingType.ToString" 的 <returns> 文档]

![](../icons/SectionExpanded.png)参见

#### 引用

[DTUSettingType 类](0a9eed32-45a3-d0eb-9104-fe18442a3351.htm)

[HslCommunication.DTU 命名空间](bf8d39e1-bb7b-54e5-7d99-73cbdd1fc4af.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)