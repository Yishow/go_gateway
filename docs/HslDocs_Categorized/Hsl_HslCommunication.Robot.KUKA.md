# HslCommunication - HslCommunication.Robot.KUKA

> 分類頁數: 30



---
## HslCommunication.Robot.KUKA

[原文連結](http://api.hslcommunication.cn/html/185b17e5-ca64-0bf5-9521-90aab446dd88.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.KUKA](../html/185b17e5-ca64-0bf5-9521-90aab446dd88.htm "HslCommunication.Robot.KUKA")

[KukaAvarProxyNet 类](../html/24a1f5b9-8b90-7fe9-b45f-6e625494f1d5.htm "KukaAvarProxyNet 类")

[KukaTcpNet 类](../html/12f49624-a008-ac45-d356-6af762d5cc4c.htm "KukaTcpNet 类")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslCommunication.Robot.KUKA 命名空间 |

[缺少 "N:HslCommunication.Robot.KUKA" 的 <summary> 文档]

![](../icons/SectionExpanded.png)类

|  | 类 | 说明 |
| --- | --- | --- |
| 公共类 | [KukaAvarProxyNet](24a1f5b9-8b90-7fe9-b45f-6e625494f1d5.htm) | Kuka机器人的数据交互类，通讯支持的条件为KUKA 的 KRC4 控制器中运行KUKAVARPROXY 这个第三方软件，端口通常为7000  The data interaction class of Kuka robot is supported by the third-party software KUKAVARPROXY running in the KRC4 controller of Kuka. The port is usually 7000 |
| 公共类 | [KukaTcpNet](12f49624-a008-ac45-d356-6af762d5cc4c.htm) | Kuka机器人的数据交互类，通讯支持的条件为KUKA 的 TCP通讯 |

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## KukaAvarProxyNet 类

[原文連結](http://api.hslcommunication.cn/html/24a1f5b9-8b90-7fe9-b45f-6e625494f1d5.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.KUKA](../html/185b17e5-ca64-0bf5-9521-90aab446dd88.htm "HslCommunication.Robot.KUKA")

[KukaAvarProxyNet 类](../html/24a1f5b9-8b90-7fe9-b45f-6e625494f1d5.htm "KukaAvarProxyNet 类")

[KukaAvarProxyNet 构造函数](../html/5f800939-b86d-c3da-6054-bb7345036f9e.htm "KukaAvarProxyNet 构造函数 ")

[KukaAvarProxyNet 属性](../html/dc604acd-e69c-c07d-0762-7f25045a24b8.htm "KukaAvarProxyNet 属性")

[KukaAvarProxyNet 方法](../html/820392c9-645f-c6f1-db0a-1ff5b1a1a86a.htm "KukaAvarProxyNet 方法")

[KukaAvarProxyNet 字段](../html/5e7871de-3ec3-a3cb-e1ea-e82fd2737fab.htm "KukaAvarProxyNet 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| KukaAvarProxyNet 类 |

Kuka机器人的数据交互类，通讯支持的条件为KUKA 的 KRC4 控制器中运行KUKAVARPROXY 这个第三方软件，端口通常为7000  
The data interaction class of Kuka robot is supported by the third-party software KUKAVARPROXY running in the KRC4 controller of Kuka. The port is usually 7000

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  [HslCommunication.Core.NetNetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)  
    [HslCommunication.Core.NetNetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)  
      HslCommunication.Robot.KUKAKukaAvarProxyNet

**命名空间：**
 [HslCommunication.Robot.KUKA](185b17e5-ca64-0bf5-9521-90aab446dd88.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class KukaAvarProxyNet : NetworkDoubleBase, 
	IRobotNet
```

```
Public Class KukaAvarProxyNet
	Inherits NetworkDoubleBase
	Implements IRobotNet
```

```
public ref class KukaAvarProxyNet : public NetworkDoubleBase, 
	IRobotNet
```

```
type KukaAvarProxyNet =  
    class
        inherit NetworkDoubleBase
        interface IRobotNet
    end
```

KukaAvarProxyNet 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [KukaAvarProxyNet](f7c931d7-d126-6854-fd91-ec8056183187.htm) | 实例化一个默认的对象  Instantiate a default object |
| 公共方法 | [KukaAvarProxyNet(String, Int32)](59e6c77e-ce47-8cce-d98a-5069a859a4e4.htm) | 实例化一个默认的Kuka机器人对象，并指定IP地址和端口号，端口号通常为7000  Instantiate a default Kuka robot object and specify the IP address and port number, usually 7000 |

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
| 受保护的方法 | [DecideWhetherQAMessage](77091dac-82d7-8c67-d691-31ae5b2aaee4.htm) | 决定当前的消息是否是应答机制的消息内容，需要在客户端进行重写实现，如果是应答机制，返回 True, 否则返回 False  To determine whether the current message is the message content of the response mechanism, it needs to be rewritten on the client side. If it is the response mechanism, return True, otherwise return False (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 受保护的方法 | [DeleteFileByName](59afa4a4-3f3c-9831-26ee-c7c13a213bf7.htm) | 删除一个指定的文件，如果文件不存在，直接返回 True，如果文件存在则直接删除，删除成功返回 True，如果发生了异常，返回False  Delete a specified file, if the file does not exist, return True directly, if the file exists, delete it directly, if the deletion is successful, return True, if an exception occurs, return  False (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 公共方法 | [Dispose](8a86a200-fbe0-3b45-0fed-c551405c9b59.htm) | 释放当前的资源，如果调用了本方法，那么该对象再使用的时候，需要重新实例化。  Release the current resource. If this method is called, the object needs to be instantiated again when it is used again. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 受保护的方法 | [Dispose(Boolean)](32f03275-9ae0-525c-84bc-6e2515d11c12.htm) | 释放当前的资源，并自动关闭长连接，如果设置了的话 (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | [ExtraAfterReadFromCoreServer](9f8357e9-dff4-1002-9368-0744cdedff29.htm) | 和服务器交互完成的时候调用的方法，可以根据读写结果进行一些额外的操作，具体的操作需要根据实际的需求来重写实现  The method called when the interaction with the server is completed can perform some additional operations based on the read and write results. The specific operations need to be rewritten according to actual needs. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 受保护的方法代码示例 | [ExtraOnDisconnect](b8569c75-5e55-cd55-5251-207465a8ba39.htm) | 根据实际的协议选择是否重写本方法，有些协议在断开连接之前，需要发送一些报文来关闭当前的网络通道  Select whether to rewrite this method according to the actual protocol. Some protocols need to send some packets to close the current network channel before disconnecting. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 受保护的方法代码示例 | [ExtraOnDisconnectAsync](86aa00d8-b2d6-ada9-5926-2b7a2237fe6b.htm) | 根据实际的协议选择是否重写本方法，有些协议在断开连接之前，需要发送一些报文来关闭当前的网络通道  Select whether to rewrite this method according to the actual protocol. Some protocols need to send some packets to close the current network channel before disconnecting. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 受保护的方法 | [GetAvailableSocket](e86a6fc1-eb91-8de9-5ff8-799b0ea3786b.htm) | 获取本次操作的可用的网络通道，如果是短连接，就重新生成一个新的网络通道，如果是长连接，就复用当前的网络通道。  Obtain the available network channels for this operation. If it is a short connection, a new network channel is regenerated. If it is a long connection, the current network channel is reused. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 受保护的方法 | [GetAvailableSocketAsync](40cf9ce7-a5e6-5431-2e4b-56e3b7c98346.htm) | 获取本次操作的可用的网络通道，如果是短连接，就重新生成一个新的网络通道，如果是长连接，就复用当前的网络通道。  Obtain the available network channels for this operation. If it is a short connection, a new network channel is regenerated. If it is a long connection, the current network channel is reused. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 受保护的方法 | [GetNewNetMessage](108e6067-f615-e3e2-b417-65c87bba7d1b.htm) | 获取一个新的消息对象的方法，需要在继承类里面进行重写  The method to get a new message object needs to be overridden in the inheritance class (重写 [NetworkDoubleBaseGetNewNetMessage](481a40e0-8539-5fc8-3e2b-c25bcc952697.htm).) |
| 公共方法 | [GetPipeSocket](19904680-372e-e0a6-1b27-b4180e759668.htm) | 获取当前用于通信的管道信息  Get the current pipe information used for communication (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法代码示例 | [InitializationOnConnect](40d5eb26-9bca-11cf-42a1-55af6ea5aebd.htm) | 根据实际的协议选择是否重写本方法，有些协议在创建连接之后，需要进行一些初始化的信号握手，才能最终建立网络通道。  Whether to rewrite this method is based on the actual protocol. Some protocols require some initial signal handshake to establish a network channel after the connection is created. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 受保护的方法代码示例 | [InitializationOnConnectAsync](14cac21b-4001-da7d-09a6-941eba868e14.htm) | 根据实际的协议选择是否重写本方法，有些协议在创建连接之后，需要进行一些初始化的信号握手，才能最终建立网络通道。  Whether to rewrite this method is based on the actual protocol. Some protocols require some initial signal handshake to establish a network channel after the connection is created. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | [IpAddressPing](d9611e18-4033-de95-444e-b91ae0a81709.htm) | 对当前设备的IP地址进行PING的操作，返回PING的结果，正常来说，返回Success  PING the IP address of the current device and return the PING result. Normally, it returns Success (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法 | [PackCommandWithHeader](f96881d6-e59f-9e16-07e3-1188eee1a934.htm) | 对当前的命令进行打包处理，通常是携带命令头内容，标记当前的命令的长度信息，需要进行重写，否则默认不打包  The current command is packaged, usually carrying the content of the command header, marking the length of the current command, and it needs to be rewritten, otherwise it is not packaged by default (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | [Read](4e54eb4f-af7d-2c60-ec92-adf6d70537dd.htm) | 读取Kuka机器人的数据内容，根据输入的变量名称来读取  Read the data content of the Kuka robot according to the input variable name |
| 公共方法 | [ReadAsync](e35f4cca-c738-ec11-a428-2f17bc766822.htm) | 读取Kuka机器人的数据内容，根据输入的变量名称来读取  Read the data content of the Kuka robot according to the input variable name |
| 公共方法代码示例 | [ReadFromCoreServer(Byte)](b3e14431-d773-0d1a-f453-a3ad60648c29.htm) | 将数据发送到当前的网络通道中，并从网络通道中接收一个[INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm)指定的完整的报文，网络通道将根据[GetAvailableSocket](e86a6fc1-eb91-8de9-5ff8-799b0ea3786b.htm)方法自动获取，本方法是线程安全的。  Send data to the current network channel and receive a complete message specified by [INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm) from the network channel. The network channel will be automatically obtained according to the [GetAvailableSocket](e86a6fc1-eb91-8de9-5ff8-799b0ea3786b.htm) method This method is thread-safe. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | [ReadFromCoreServer(IEnumerableByte)](580dae22-fad5-5816-d58f-43068e89e419.htm) | 将多个数据报文按顺序发到设备，并从设备接收返回的数据内容，然后拼接成一个Byte[]信息，需要重写[UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm)方法才能返回正确的结果。  Send multiple data packets to the device in sequence, and receive the returned data content from the device, and then splicing them into a Byte[] message, you need to rewrite [UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm) method to return the correct result. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServer(Byte, Boolean, Boolean)](1a33d974-805d-f149-3665-eb8d8a99ea84.htm) | 将数据发送到当前的网络通道中，并从网络通道中接收一个[INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm)指定的完整的报文，网络通道将根据[GetAvailableSocket](e86a6fc1-eb91-8de9-5ff8-799b0ea3786b.htm)方法自动获取，本方法是线程安全的。  Send data to the current network channel and receive a complete message specified by [INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm) from the network channel. The network channel will be automatically obtained according to the [GetAvailableSocket](e86a6fc1-eb91-8de9-5ff8-799b0ea3786b.htm) method This method is thread-safe. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServer(Socket, Byte, Boolean, Boolean)](3b72b289-b201-7548-0981-ac049c9a4135.htm) | 将数据报文发送指定的网络通道上，根据当前指定的[INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm)类型，返回一条完整的数据指令  Sends a data message to the specified network channel, and returns a complete data command according to the currently specified [INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm) type (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServerAsync(Byte)](ec51a18e-0f5d-dbe1-c229-d62e8760ac97.htm) | 将数据发送到当前的网络通道中，并从网络通道中接收一个[INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm)指定的完整的报文，网络通道将根据[GetAvailableSocket](e86a6fc1-eb91-8de9-5ff8-799b0ea3786b.htm)方法自动获取，本方法是线程安全的。  Send data to the current network channel and receive a complete message specified by [INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm) from the network channel. The network channel will be automatically obtained according to the [GetAvailableSocket](e86a6fc1-eb91-8de9-5ff8-799b0ea3786b.htm) method This method is thread-safe. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | [ReadFromCoreServerAsync(IEnumerableByte)](0fce97fd-8441-6b54-0877-b8dfacc738e3.htm) | 将多个数据报文按顺序发到设备，并从设备接收返回的数据内容，然后拼接成一个Byte[]信息，需要重写[UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm)方法才能返回正确的结果。  Send multiple data packets to the device in sequence, and receive the returned data content from the device, and then splicing them into a Byte[] message, you need to rewrite [UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm) method to return the correct result. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServerAsync(Byte, Boolean, Boolean)](7525b379-04c2-8dac-b451-a8b5413d6747.htm) | 将数据发送到当前的网络通道中，并从网络通道中接收一个[INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm)指定的完整的报文，网络通道将根据[GetAvailableSocket](e86a6fc1-eb91-8de9-5ff8-799b0ea3786b.htm)方法自动获取，本方法是线程安全的。  Send data to the current network channel and receive a complete message specified by [INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm) from the network channel. The network channel will be automatically obtained according to the [GetAvailableSocket](e86a6fc1-eb91-8de9-5ff8-799b0ea3786b.htm) method This method is thread-safe. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServerAsync(Socket, Byte, Boolean, Boolean)](07db4a00-b548-8bd6-5de6-12deb91f5311.htm) | 将数据报文发送指定的网络通道上，根据当前指定的[INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm)类型，返回一条完整的数据指令  Sends a data message to the specified network channel, and returns a complete data command according to the currently specified [INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm) type (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | [ReadString](1bd22882-b614-1913-93d6-50652f750ae1.htm) | 读取Kuka机器人的所有的数据信息，返回字符串信息，解码方式为ANSI，需要指定变量名称  Read all the data information of the Kuka robot, return the string information, decode by ANSI, need to specify the variable name |
| 公共方法 | [ReadStringAsync](9380eee6-0b90-94d3-52f2-73bc520962ec.htm) | 读取Kuka机器人的所有的数据信息，返回字符串信息，解码方式为ANSI，需要指定变量名称  Read all the data information of the Kuka robot, return the string information, decode by ANSI, need to specify the variable name |
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
| 公共方法 | [SetLoginAccount](725afc25-6e3f-bd7d-25ed-580270c8fc56.htm) | 设置当前的登录的账户名和密码信息，并启用账户验证的功能，账户名为空时设置不生效  Set the current login account name and password information, and enable the account verification function. The account name setting will not take effect when it is empty (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法代码示例 | [SetPersistentConnection](d8e8df9f-cbe9-6863-54e9-aa3b38906b26.htm) | 在读取数据之前可以调用本方法将客户端设置为长连接模式，相当于跳过了ConnectServer的结果验证，对异形客户端无效，当第一次进行通信时再进行创建连接请求。  Before reading the data, you can call this method to set the client to the long connection mode, which is equivalent to skipping the result verification of ConnectServer, and it is invalid for the alien client. When the first communication is performed, the connection creation request is performed. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | [SetPipeSocket](b6a7a243-2a6d-c4ea-45b7-bf980450a2bd.htm) | 设置一个新的网络管道，一般来说不需要调用本方法，当多个网口设备共用一个网络连接时才需要使用本方法进行设置共享的管道。  To set up a new network channel, generally speaking, you do not need to call this method. This method is only needed to set up a shared channel when multiple network port devices share a network connection. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | [ToString](097819cc-7222-4fe5-0a5c-5d565906c647.htm) | (重写 [NetworkDoubleBaseToString](85a8532f-8fc5-886c-9d76-70844e0a7d64.htm).) |
| 公共方法 | [UnpackResponseContent](32bad9a0-b59d-c171-ce85-74055edfc597.htm) | 根据对方返回的报文命令，对命令进行基本的拆包，例如各种Modbus协议拆包为统一的核心报文，还支持对报文的验证  According to the message command returned by the other party, the command is basically unpacked, for example, various Modbus protocols are unpacked into a unified core message, and the verification of the message is also supported (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | [Write(String, Byte)](dcdb2813-429f-eb16-53f7-793cc06914fc.htm) | 根据Kuka机器人的变量名称，写入原始的数据内容  Write the original data content according to the variable name of the Kuka robot |
| 公共方法 | [Write(String, String)](3a9995fa-bb16-db1c-fb1f-094f31b4f57b.htm) | 根据Kuka机器人的变量名称，写入ANSI编码的字符串数据信息  Writes ansi-encoded string data information based on the variable name of the Kuka robot |
| 公共方法 | [WriteAsync(String, Byte)](b4cd9554-b881-5e87-194c-39abd678bd95.htm) | 根据Kuka机器人的变量名称，写入原始的数据内容  Write the original data content according to the variable name of the Kuka robot |
| 公共方法 | [WriteAsync(String, String)](ab50cf71-4c25-d98b-23b9-298fdd41dd8a.htm) | 根据Kuka机器人的变量名称，写入ANSI编码的字符串数据信息  Writes ansi-encoded string data information based on the variable name of the Kuka robot |
| 受保护的方法 | [WriteStreamFromSocket](97bf95ea-51fe-8f2e-5da6-1c83766aaa78.htm) | 从套接字中接收所有的数据然后写入到指定的流当中去，需要指定数据的长度，支持按照百分比进行进度报告  Receives all data from the socket and writes it to the specified stream. The length of the data needs to be specified, and progress reporting is supported in percentage. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [WriteStreamFromSocketAsync](b4cf7f6d-1e0c-a266-4c49-683fdb14a922.htm) | 从套接字中接收所有的数据然后写入到指定的流当中去，需要指定数据的长度，支持按照百分比进行进度报告  Receives all data from the socket and writes it to the specified stream. The length of the data needs to be specified, and progress reporting is supported in percentage. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |

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

非常感谢 昆山-LT 网友的测试和意见反馈。  
其中KUKAVARPROXY 这个第三方软件在来源地址：
https://github.com/ImtsSrl/KUKAVARPROXY

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.Robot.KUKA 命名空间](185b17e5-ca64-0bf5-9521-90aab446dd88.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## KukaAvarProxyNet 构造函数 

[原文連結](http://api.hslcommunication.cn/html/5f800939-b86d-c3da-6054-bb7345036f9e.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.KUKA](../html/185b17e5-ca64-0bf5-9521-90aab446dd88.htm "HslCommunication.Robot.KUKA")

[KukaAvarProxyNet 类](../html/24a1f5b9-8b90-7fe9-b45f-6e625494f1d5.htm "KukaAvarProxyNet 类")

[KukaAvarProxyNet 构造函数](../html/5f800939-b86d-c3da-6054-bb7345036f9e.htm "KukaAvarProxyNet 构造函数 ")

[KukaAvarProxyNet 构造函数](../html/f7c931d7-d126-6854-fd91-ec8056183187.htm "KukaAvarProxyNet 构造函数 ")

[KukaAvarProxyNet 构造函数 (String, Int32)](../html/59e6c77e-ce47-8cce-d98a-5069a859a4e4.htm "KukaAvarProxyNet 构造函数 (String, Int32)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| KukaAvarProxyNet 构造函数 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [KukaAvarProxyNet](f7c931d7-d126-6854-fd91-ec8056183187.htm) | 实例化一个默认的对象  Instantiate a default object |
| 公共方法 | [KukaAvarProxyNet(String, Int32)](59e6c77e-ce47-8cce-d98a-5069a859a4e4.htm) | 实例化一个默认的Kuka机器人对象，并指定IP地址和端口号，端口号通常为7000  Instantiate a default Kuka robot object and specify the IP address and port number, usually 7000 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[KukaAvarProxyNet 类](24a1f5b9-8b90-7fe9-b45f-6e625494f1d5.htm)

[HslCommunication.Robot.KUKA 命名空间](185b17e5-ca64-0bf5-9521-90aab446dd88.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## KukaAvarProxyNet 构造函数 

[原文連結](http://api.hslcommunication.cn/html/f7c931d7-d126-6854-fd91-ec8056183187.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.KUKA](../html/185b17e5-ca64-0bf5-9521-90aab446dd88.htm "HslCommunication.Robot.KUKA")

[KukaAvarProxyNet 类](../html/24a1f5b9-8b90-7fe9-b45f-6e625494f1d5.htm "KukaAvarProxyNet 类")

[KukaAvarProxyNet 构造函数](../html/5f800939-b86d-c3da-6054-bb7345036f9e.htm "KukaAvarProxyNet 构造函数 ")

[KukaAvarProxyNet 构造函数](../html/f7c931d7-d126-6854-fd91-ec8056183187.htm "KukaAvarProxyNet 构造函数 ")

[KukaAvarProxyNet 构造函数 (String, Int32)](../html/59e6c77e-ce47-8cce-d98a-5069a859a4e4.htm "KukaAvarProxyNet 构造函数 (String, Int32)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| KukaAvarProxyNet 构造函数 |

实例化一个默认的对象  
Instantiate a default object

**命名空间：**
 [HslCommunication.Robot.KUKA](185b17e5-ca64-0bf5-9521-90aab446dd88.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public KukaAvarProxyNet()
```

```
Public Sub New
```

```
public:
KukaAvarProxyNet()
```

```
new : unit -> KukaAvarProxyNet
```

![](../icons/SectionExpanded.png)参见

#### 引用

[KukaAvarProxyNet 类](24a1f5b9-8b90-7fe9-b45f-6e625494f1d5.htm)

[KukaAvarProxyNet 重载](5f800939-b86d-c3da-6054-bb7345036f9e.htm)

[HslCommunication.Robot.KUKA 命名空间](185b17e5-ca64-0bf5-9521-90aab446dd88.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## KukaAvarProxyNet 属性

[原文連結](http://api.hslcommunication.cn/html/dc604acd-e69c-c07d-0762-7f25045a24b8.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.KUKA](../html/185b17e5-ca64-0bf5-9521-90aab446dd88.htm "HslCommunication.Robot.KUKA")

[KukaAvarProxyNet 类](../html/24a1f5b9-8b90-7fe9-b45f-6e625494f1d5.htm "KukaAvarProxyNet 类")

[KukaAvarProxyNet 构造函数](../html/5f800939-b86d-c3da-6054-bb7345036f9e.htm "KukaAvarProxyNet 构造函数 ")

[KukaAvarProxyNet 属性](../html/dc604acd-e69c-c07d-0762-7f25045a24b8.htm "KukaAvarProxyNet 属性")

[KukaAvarProxyNet 方法](../html/820392c9-645f-c6f1-db0a-1ff5b1a1a86a.htm "KukaAvarProxyNet 方法")

[KukaAvarProxyNet 字段](../html/5e7871de-3ec3-a3cb-e1ea-e82fd2737fab.htm "KukaAvarProxyNet 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| KukaAvarProxyNet 属性 |

[KukaAvarProxyNet](24a1f5b9-8b90-7fe9-b45f-6e625494f1d5.htm) 类型公开以下成员。

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

[KukaAvarProxyNet 类](24a1f5b9-8b90-7fe9-b45f-6e625494f1d5.htm)

[HslCommunication.Robot.KUKA 命名空间](185b17e5-ca64-0bf5-9521-90aab446dd88.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GetNewNetMessage 方法 

[原文連結](http://api.hslcommunication.cn/html/108e6067-f615-e3e2-b417-65c87bba7d1b.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.KUKA](../html/185b17e5-ca64-0bf5-9521-90aab446dd88.htm "HslCommunication.Robot.KUKA")

[KukaAvarProxyNet 类](../html/24a1f5b9-8b90-7fe9-b45f-6e625494f1d5.htm "KukaAvarProxyNet 类")

[KukaAvarProxyNet 方法](../html/820392c9-645f-c6f1-db0a-1ff5b1a1a86a.htm "KukaAvarProxyNet 方法")

[GetNewNetMessage 方法](../html/108e6067-f615-e3e2-b417-65c87bba7d1b.htm "GetNewNetMessage 方法 ")

[Read 方法](../html/4e54eb4f-af7d-2c60-ec92-adf6d70537dd.htm "Read 方法 ")

[ReadAsync 方法](../html/e35f4cca-c738-ec11-a428-2f17bc766822.htm "ReadAsync 方法 ")

[ReadString 方法](../html/1bd22882-b614-1913-93d6-50652f750ae1.htm "ReadString 方法 ")

[ReadStringAsync 方法](../html/9380eee6-0b90-94d3-52f2-73bc520962ec.htm "ReadStringAsync 方法 ")

[ToString 方法](../html/097819cc-7222-4fe5-0a5c-5d565906c647.htm "ToString 方法 ")

[Write 方法](../html/b3f29076-b464-4a14-4093-81d8a68bb603.htm "Write 方法 ")

[WriteAsync 方法](../html/4c251649-c7dc-46be-6996-d9eca2f77047.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| KukaAvarProxyNetGetNewNetMessage 方法 |

获取一个新的消息对象的方法，需要在继承类里面进行重写  
The method to get a new message object needs to be overridden in the inheritance class

**命名空间：**
 [HslCommunication.Robot.KUKA](185b17e5-ca64-0bf5-9521-90aab446dd88.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
protected override INetMessage GetNewNetMessage()
```

```
Protected Overrides Function GetNewNetMessage As INetMessage
```

```
protected:
virtual INetMessage^ GetNewNetMessage() override
```

```
abstract GetNewNetMessage : unit -> INetMessage 
override GetNewNetMessage : unit -> INetMessage
```

#### 返回值

类型：[INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm)  
消息类对象

![](../icons/SectionExpanded.png)参见

#### 引用

[KukaAvarProxyNet 类](24a1f5b9-8b90-7fe9-b45f-6e625494f1d5.htm)

[HslCommunication.Robot.KUKA 命名空间](185b17e5-ca64-0bf5-9521-90aab446dd88.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Read 方法 

[原文連結](http://api.hslcommunication.cn/html/4e54eb4f-af7d-2c60-ec92-adf6d70537dd.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.KUKA](../html/185b17e5-ca64-0bf5-9521-90aab446dd88.htm "HslCommunication.Robot.KUKA")

[KukaAvarProxyNet 类](../html/24a1f5b9-8b90-7fe9-b45f-6e625494f1d5.htm "KukaAvarProxyNet 类")

[KukaAvarProxyNet 方法](../html/820392c9-645f-c6f1-db0a-1ff5b1a1a86a.htm "KukaAvarProxyNet 方法")

[GetNewNetMessage 方法](../html/108e6067-f615-e3e2-b417-65c87bba7d1b.htm "GetNewNetMessage 方法 ")

[Read 方法](../html/4e54eb4f-af7d-2c60-ec92-adf6d70537dd.htm "Read 方法 ")

[ReadAsync 方法](../html/e35f4cca-c738-ec11-a428-2f17bc766822.htm "ReadAsync 方法 ")

[ReadString 方法](../html/1bd22882-b614-1913-93d6-50652f750ae1.htm "ReadString 方法 ")

[ReadStringAsync 方法](../html/9380eee6-0b90-94d3-52f2-73bc520962ec.htm "ReadStringAsync 方法 ")

[ToString 方法](../html/097819cc-7222-4fe5-0a5c-5d565906c647.htm "ToString 方法 ")

[Write 方法](../html/b3f29076-b464-4a14-4093-81d8a68bb603.htm "Write 方法 ")

[WriteAsync 方法](../html/4c251649-c7dc-46be-6996-d9eca2f77047.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| KukaAvarProxyNetRead 方法 |

读取Kuka机器人的数据内容，根据输入的变量名称来读取  
Read the data content of the Kuka robot according to the input variable name

**命名空间：**
 [HslCommunication.Robot.KUKA](185b17e5-ca64-0bf5-9521-90aab446dd88.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public OperateResult<byte[]> Read(
	string address
)
```

```
Public Function Read ( 
	address As String
) As OperateResult(Of Byte())
```

```
public:
virtual OperateResult<array<unsigned char>^>^ Read(
	String^ address
) sealed
```

```
abstract Read : 
        address : string -> OperateResult<byte[]> 
override Read : 
        address : string -> OperateResult<byte[]>
```

#### 参数

address
:   类型：SystemString  
    地址数据

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
带有成功标识的byte[]数组

#### 实现

[IRobotNetRead(String)](6a9e3fe2-bbb6-8753-8d55-3c5d41fe640a.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[KukaAvarProxyNet 类](24a1f5b9-8b90-7fe9-b45f-6e625494f1d5.htm)

[HslCommunication.Robot.KUKA 命名空间](185b17e5-ca64-0bf5-9521-90aab446dd88.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/e35f4cca-c738-ec11-a428-2f17bc766822.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.KUKA](../html/185b17e5-ca64-0bf5-9521-90aab446dd88.htm "HslCommunication.Robot.KUKA")

[KukaAvarProxyNet 类](../html/24a1f5b9-8b90-7fe9-b45f-6e625494f1d5.htm "KukaAvarProxyNet 类")

[KukaAvarProxyNet 方法](../html/820392c9-645f-c6f1-db0a-1ff5b1a1a86a.htm "KukaAvarProxyNet 方法")

[GetNewNetMessage 方法](../html/108e6067-f615-e3e2-b417-65c87bba7d1b.htm "GetNewNetMessage 方法 ")

[Read 方法](../html/4e54eb4f-af7d-2c60-ec92-adf6d70537dd.htm "Read 方法 ")

[ReadAsync 方法](../html/e35f4cca-c738-ec11-a428-2f17bc766822.htm "ReadAsync 方法 ")

[ReadString 方法](../html/1bd22882-b614-1913-93d6-50652f750ae1.htm "ReadString 方法 ")

[ReadStringAsync 方法](../html/9380eee6-0b90-94d3-52f2-73bc520962ec.htm "ReadStringAsync 方法 ")

[ToString 方法](../html/097819cc-7222-4fe5-0a5c-5d565906c647.htm "ToString 方法 ")

[Write 方法](../html/b3f29076-b464-4a14-4093-81d8a68bb603.htm "Write 方法 ")

[WriteAsync 方法](../html/4c251649-c7dc-46be-6996-d9eca2f77047.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| KukaAvarProxyNetReadAsync 方法 |

读取Kuka机器人的数据内容，根据输入的变量名称来读取  
Read the data content of the Kuka robot according to the input variable name

**命名空间：**
 [HslCommunication.Robot.KUKA](185b17e5-ca64-0bf5-9521-90aab446dd88.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public Task<OperateResult<byte[]>> ReadAsync(
	string address
)
```

```
Public Function ReadAsync ( 
	address As String
) As Task(Of OperateResult(Of Byte()))
```

```
public:
virtual Task<OperateResult<array<unsigned char>^>^>^ ReadAsync(
	String^ address
) sealed
```

```
abstract ReadAsync : 
        address : string -> Task<OperateResult<byte[]>> 
override ReadAsync : 
        address : string -> Task<OperateResult<byte[]>>
```

#### 参数

address
:   类型：SystemString  
    地址数据

#### 返回值

类型：Task[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
带有成功标识的byte[]数组

#### 实现

[IRobotNetReadAsync(String)](5a0ad618-3b88-4c9a-0582-e90ec8650628.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[KukaAvarProxyNet 类](24a1f5b9-8b90-7fe9-b45f-6e625494f1d5.htm)

[HslCommunication.Robot.KUKA 命名空间](185b17e5-ca64-0bf5-9521-90aab446dd88.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadString 方法 

[原文連結](http://api.hslcommunication.cn/html/1bd22882-b614-1913-93d6-50652f750ae1.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.KUKA](../html/185b17e5-ca64-0bf5-9521-90aab446dd88.htm "HslCommunication.Robot.KUKA")

[KukaAvarProxyNet 类](../html/24a1f5b9-8b90-7fe9-b45f-6e625494f1d5.htm "KukaAvarProxyNet 类")

[KukaAvarProxyNet 方法](../html/820392c9-645f-c6f1-db0a-1ff5b1a1a86a.htm "KukaAvarProxyNet 方法")

[GetNewNetMessage 方法](../html/108e6067-f615-e3e2-b417-65c87bba7d1b.htm "GetNewNetMessage 方法 ")

[Read 方法](../html/4e54eb4f-af7d-2c60-ec92-adf6d70537dd.htm "Read 方法 ")

[ReadAsync 方法](../html/e35f4cca-c738-ec11-a428-2f17bc766822.htm "ReadAsync 方法 ")

[ReadString 方法](../html/1bd22882-b614-1913-93d6-50652f750ae1.htm "ReadString 方法 ")

[ReadStringAsync 方法](../html/9380eee6-0b90-94d3-52f2-73bc520962ec.htm "ReadStringAsync 方法 ")

[ToString 方法](../html/097819cc-7222-4fe5-0a5c-5d565906c647.htm "ToString 方法 ")

[Write 方法](../html/b3f29076-b464-4a14-4093-81d8a68bb603.htm "Write 方法 ")

[WriteAsync 方法](../html/4c251649-c7dc-46be-6996-d9eca2f77047.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| KukaAvarProxyNetReadString 方法 |

读取Kuka机器人的所有的数据信息，返回字符串信息，解码方式为ANSI，需要指定变量名称  
Read all the data information of the Kuka robot, return the string information, decode by ANSI, need to specify the variable name

**命名空间：**
 [HslCommunication.Robot.KUKA](185b17e5-ca64-0bf5-9521-90aab446dd88.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public OperateResult<string> ReadString(
	string address
)
```

```
Public Function ReadString ( 
	address As String
) As OperateResult(Of String)
```

```
public:
virtual OperateResult<String^>^ ReadString(
	String^ address
) sealed
```

```
abstract ReadString : 
        address : string -> OperateResult<string> 
override ReadString : 
        address : string -> OperateResult<string>
```

#### 参数

address
:   类型：SystemString  
    地址信息

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)String  
带有成功标识的字符串数据

#### 实现

[IRobotNetReadString(String)](c3f26090-6d4e-8bad-be21-35d5371c1d6f.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[KukaAvarProxyNet 类](24a1f5b9-8b90-7fe9-b45f-6e625494f1d5.htm)

[HslCommunication.Robot.KUKA 命名空间](185b17e5-ca64-0bf5-9521-90aab446dd88.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadStringAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/9380eee6-0b90-94d3-52f2-73bc520962ec.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.KUKA](../html/185b17e5-ca64-0bf5-9521-90aab446dd88.htm "HslCommunication.Robot.KUKA")

[KukaAvarProxyNet 类](../html/24a1f5b9-8b90-7fe9-b45f-6e625494f1d5.htm "KukaAvarProxyNet 类")

[KukaAvarProxyNet 方法](../html/820392c9-645f-c6f1-db0a-1ff5b1a1a86a.htm "KukaAvarProxyNet 方法")

[GetNewNetMessage 方法](../html/108e6067-f615-e3e2-b417-65c87bba7d1b.htm "GetNewNetMessage 方法 ")

[Read 方法](../html/4e54eb4f-af7d-2c60-ec92-adf6d70537dd.htm "Read 方法 ")

[ReadAsync 方法](../html/e35f4cca-c738-ec11-a428-2f17bc766822.htm "ReadAsync 方法 ")

[ReadString 方法](../html/1bd22882-b614-1913-93d6-50652f750ae1.htm "ReadString 方法 ")

[ReadStringAsync 方法](../html/9380eee6-0b90-94d3-52f2-73bc520962ec.htm "ReadStringAsync 方法 ")

[ToString 方法](../html/097819cc-7222-4fe5-0a5c-5d565906c647.htm "ToString 方法 ")

[Write 方法](../html/b3f29076-b464-4a14-4093-81d8a68bb603.htm "Write 方法 ")

[WriteAsync 方法](../html/4c251649-c7dc-46be-6996-d9eca2f77047.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| KukaAvarProxyNetReadStringAsync 方法 |

读取Kuka机器人的所有的数据信息，返回字符串信息，解码方式为ANSI，需要指定变量名称  
Read all the data information of the Kuka robot, return the string information, decode by ANSI, need to specify the variable name

**命名空间：**
 [HslCommunication.Robot.KUKA](185b17e5-ca64-0bf5-9521-90aab446dd88.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public Task<OperateResult<string>> ReadStringAsync(
	string address
)
```

```
Public Function ReadStringAsync ( 
	address As String
) As Task(Of OperateResult(Of String))
```

```
public:
virtual Task<OperateResult<String^>^>^ ReadStringAsync(
	String^ address
) sealed
```

```
abstract ReadStringAsync : 
        address : string -> Task<OperateResult<string>> 
override ReadStringAsync : 
        address : string -> Task<OperateResult<string>>
```

#### 参数

address
:   类型：SystemString  
    地址信息

#### 返回值

类型：Task[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)String  
带有成功标识的字符串数据

#### 实现

[IRobotNetReadStringAsync(String)](860414d5-bc23-62a1-7ae4-ac578deaa375.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[KukaAvarProxyNet 类](24a1f5b9-8b90-7fe9-b45f-6e625494f1d5.htm)

[HslCommunication.Robot.KUKA 命名空间](185b17e5-ca64-0bf5-9521-90aab446dd88.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ToString 方法 

[原文連結](http://api.hslcommunication.cn/html/097819cc-7222-4fe5-0a5c-5d565906c647.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.KUKA](../html/185b17e5-ca64-0bf5-9521-90aab446dd88.htm "HslCommunication.Robot.KUKA")

[KukaAvarProxyNet 类](../html/24a1f5b9-8b90-7fe9-b45f-6e625494f1d5.htm "KukaAvarProxyNet 类")

[KukaAvarProxyNet 方法](../html/820392c9-645f-c6f1-db0a-1ff5b1a1a86a.htm "KukaAvarProxyNet 方法")

[GetNewNetMessage 方法](../html/108e6067-f615-e3e2-b417-65c87bba7d1b.htm "GetNewNetMessage 方法 ")

[Read 方法](../html/4e54eb4f-af7d-2c60-ec92-adf6d70537dd.htm "Read 方法 ")

[ReadAsync 方法](../html/e35f4cca-c738-ec11-a428-2f17bc766822.htm "ReadAsync 方法 ")

[ReadString 方法](../html/1bd22882-b614-1913-93d6-50652f750ae1.htm "ReadString 方法 ")

[ReadStringAsync 方法](../html/9380eee6-0b90-94d3-52f2-73bc520962ec.htm "ReadStringAsync 方法 ")

[ToString 方法](../html/097819cc-7222-4fe5-0a5c-5d565906c647.htm "ToString 方法 ")

[Write 方法](../html/b3f29076-b464-4a14-4093-81d8a68bb603.htm "Write 方法 ")

[WriteAsync 方法](../html/4c251649-c7dc-46be-6996-d9eca2f77047.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| KukaAvarProxyNetToString 方法 |

[缺少 "M:HslCommunication.Robot.KUKA.KukaAvarProxyNet.ToString" 的 <summary> 文档]

**命名空间：**
 [HslCommunication.Robot.KUKA](185b17e5-ca64-0bf5-9521-90aab446dd88.htm)  
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

[缺少 "M:HslCommunication.Robot.KUKA.KukaAvarProxyNet.ToString" 的 <returns> 文档]

![](../icons/SectionExpanded.png)参见

#### 引用

[KukaAvarProxyNet 类](24a1f5b9-8b90-7fe9-b45f-6e625494f1d5.htm)

[HslCommunication.Robot.KUKA 命名空间](185b17e5-ca64-0bf5-9521-90aab446dd88.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Write 方法 

[原文連結](http://api.hslcommunication.cn/html/b3f29076-b464-4a14-4093-81d8a68bb603.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.KUKA](../html/185b17e5-ca64-0bf5-9521-90aab446dd88.htm "HslCommunication.Robot.KUKA")

[KukaAvarProxyNet 类](../html/24a1f5b9-8b90-7fe9-b45f-6e625494f1d5.htm "KukaAvarProxyNet 类")

[KukaAvarProxyNet 方法](../html/820392c9-645f-c6f1-db0a-1ff5b1a1a86a.htm "KukaAvarProxyNet 方法")

[Write 方法](../html/b3f29076-b464-4a14-4093-81d8a68bb603.htm "Write 方法 ")

[Write 方法 (String, Byte[])](../html/dcdb2813-429f-eb16-53f7-793cc06914fc.htm "Write 方法 (String, Byte[])")

[Write 方法 (String, String)](../html/3a9995fa-bb16-db1c-fb1f-094f31b4f57b.htm "Write 方法 (String, String)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| KukaAvarProxyNetWrite 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [Write(String, Byte)](dcdb2813-429f-eb16-53f7-793cc06914fc.htm) | 根据Kuka机器人的变量名称，写入原始的数据内容  Write the original data content according to the variable name of the Kuka robot |
| 公共方法 | [Write(String, String)](3a9995fa-bb16-db1c-fb1f-094f31b4f57b.htm) | 根据Kuka机器人的变量名称，写入ANSI编码的字符串数据信息  Writes ansi-encoded string data information based on the variable name of the Kuka robot |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[KukaAvarProxyNet 类](24a1f5b9-8b90-7fe9-b45f-6e625494f1d5.htm)

[HslCommunication.Robot.KUKA 命名空间](185b17e5-ca64-0bf5-9521-90aab446dd88.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Write 方法 (String, Byte[])

[原文連結](http://api.hslcommunication.cn/html/dcdb2813-429f-eb16-53f7-793cc06914fc.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.KUKA](../html/185b17e5-ca64-0bf5-9521-90aab446dd88.htm "HslCommunication.Robot.KUKA")

[KukaAvarProxyNet 类](../html/24a1f5b9-8b90-7fe9-b45f-6e625494f1d5.htm "KukaAvarProxyNet 类")

[KukaAvarProxyNet 方法](../html/820392c9-645f-c6f1-db0a-1ff5b1a1a86a.htm "KukaAvarProxyNet 方法")

[Write 方法](../html/b3f29076-b464-4a14-4093-81d8a68bb603.htm "Write 方法 ")

[Write 方法 (String, Byte[])](../html/dcdb2813-429f-eb16-53f7-793cc06914fc.htm "Write 方法 (String, Byte[])")

[Write 方法 (String, String)](../html/3a9995fa-bb16-db1c-fb1f-094f31b4f57b.htm "Write 方法 (String, String)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| KukaAvarProxyNetWrite 方法 (String, Byte) |

根据Kuka机器人的变量名称，写入原始的数据内容  
Write the original data content according to the variable name of the Kuka robot

**命名空间：**
 [HslCommunication.Robot.KUKA](185b17e5-ca64-0bf5-9521-90aab446dd88.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public OperateResult Write(
	string address,
	byte[] value
)
```

```
Public Function Write ( 
	address As String,
	value As Byte()
) As OperateResult
```

```
public:
virtual OperateResult^ Write(
	String^ address, 
	array<unsigned char>^ value
) sealed
```

```
abstract Write : 
        address : string * 
        value : byte[] -> OperateResult 
override Write : 
        address : string * 
        value : byte[] -> OperateResult
```

#### 参数

address
:   类型：SystemString  
    变量名称

value
:   类型：SystemByte  
    原始的字节数据信息

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否成功的写入

#### 实现

[IRobotNetWrite(String, Byte)](4abedd7b-e673-48ed-05b2-ba17d7709ea6.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[KukaAvarProxyNet 类](24a1f5b9-8b90-7fe9-b45f-6e625494f1d5.htm)

[Write 重载](b3f29076-b464-4a14-4093-81d8a68bb603.htm)

[HslCommunication.Robot.KUKA 命名空间](185b17e5-ca64-0bf5-9521-90aab446dd88.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Write 方法 (String, String)

[原文連結](http://api.hslcommunication.cn/html/3a9995fa-bb16-db1c-fb1f-094f31b4f57b.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.KUKA](../html/185b17e5-ca64-0bf5-9521-90aab446dd88.htm "HslCommunication.Robot.KUKA")

[KukaAvarProxyNet 类](../html/24a1f5b9-8b90-7fe9-b45f-6e625494f1d5.htm "KukaAvarProxyNet 类")

[KukaAvarProxyNet 方法](../html/820392c9-645f-c6f1-db0a-1ff5b1a1a86a.htm "KukaAvarProxyNet 方法")

[Write 方法](../html/b3f29076-b464-4a14-4093-81d8a68bb603.htm "Write 方法 ")

[Write 方法 (String, Byte[])](../html/dcdb2813-429f-eb16-53f7-793cc06914fc.htm "Write 方法 (String, Byte[])")

[Write 方法 (String, String)](../html/3a9995fa-bb16-db1c-fb1f-094f31b4f57b.htm "Write 方法 (String, String)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| KukaAvarProxyNetWrite 方法 (String, String) |

根据Kuka机器人的变量名称，写入ANSI编码的字符串数据信息  
Writes ansi-encoded string data information based on the variable name of the Kuka robot

**命名空间：**
 [HslCommunication.Robot.KUKA](185b17e5-ca64-0bf5-9521-90aab446dd88.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public OperateResult Write(
	string address,
	string value
)
```

```
Public Function Write ( 
	address As String,
	value As String
) As OperateResult
```

```
public:
virtual OperateResult^ Write(
	String^ address, 
	String^ value
) sealed
```

```
abstract Write : 
        address : string * 
        value : string -> OperateResult 
override Write : 
        address : string * 
        value : string -> OperateResult
```

#### 参数

address
:   类型：SystemString  
    变量名称

value
:   类型：SystemString  
    ANSI编码的字符串

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否成功的写入

#### 实现

[IRobotNetWrite(String, String)](a03312cc-d8cd-0bda-d849-0ebb57e7223a.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[KukaAvarProxyNet 类](24a1f5b9-8b90-7fe9-b45f-6e625494f1d5.htm)

[Write 重载](b3f29076-b464-4a14-4093-81d8a68bb603.htm)

[HslCommunication.Robot.KUKA 命名空间](185b17e5-ca64-0bf5-9521-90aab446dd88.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## WriteAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/4c251649-c7dc-46be-6996-d9eca2f77047.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.KUKA](../html/185b17e5-ca64-0bf5-9521-90aab446dd88.htm "HslCommunication.Robot.KUKA")

[KukaAvarProxyNet 类](../html/24a1f5b9-8b90-7fe9-b45f-6e625494f1d5.htm "KukaAvarProxyNet 类")

[KukaAvarProxyNet 方法](../html/820392c9-645f-c6f1-db0a-1ff5b1a1a86a.htm "KukaAvarProxyNet 方法")

[WriteAsync 方法](../html/4c251649-c7dc-46be-6996-d9eca2f77047.htm "WriteAsync 方法 ")

[WriteAsync 方法 (String, Byte[])](../html/b4cd9554-b881-5e87-194c-39abd678bd95.htm "WriteAsync 方法 (String, Byte[])")

[WriteAsync 方法 (String, String)](../html/ab50cf71-4c25-d98b-23b9-298fdd41dd8a.htm "WriteAsync 方法 (String, String)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| KukaAvarProxyNetWriteAsync 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [WriteAsync(String, Byte)](b4cd9554-b881-5e87-194c-39abd678bd95.htm) | 根据Kuka机器人的变量名称，写入原始的数据内容  Write the original data content according to the variable name of the Kuka robot |
| 公共方法 | [WriteAsync(String, String)](ab50cf71-4c25-d98b-23b9-298fdd41dd8a.htm) | 根据Kuka机器人的变量名称，写入ANSI编码的字符串数据信息  Writes ansi-encoded string data information based on the variable name of the Kuka robot |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[KukaAvarProxyNet 类](24a1f5b9-8b90-7fe9-b45f-6e625494f1d5.htm)

[HslCommunication.Robot.KUKA 命名空间](185b17e5-ca64-0bf5-9521-90aab446dd88.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## WriteAsync 方法 (String, Byte[])

[原文連結](http://api.hslcommunication.cn/html/b4cd9554-b881-5e87-194c-39abd678bd95.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.KUKA](../html/185b17e5-ca64-0bf5-9521-90aab446dd88.htm "HslCommunication.Robot.KUKA")

[KukaAvarProxyNet 类](../html/24a1f5b9-8b90-7fe9-b45f-6e625494f1d5.htm "KukaAvarProxyNet 类")

[KukaAvarProxyNet 方法](../html/820392c9-645f-c6f1-db0a-1ff5b1a1a86a.htm "KukaAvarProxyNet 方法")

[WriteAsync 方法](../html/4c251649-c7dc-46be-6996-d9eca2f77047.htm "WriteAsync 方法 ")

[WriteAsync 方法 (String, Byte[])](../html/b4cd9554-b881-5e87-194c-39abd678bd95.htm "WriteAsync 方法 (String, Byte[])")

[WriteAsync 方法 (String, String)](../html/ab50cf71-4c25-d98b-23b9-298fdd41dd8a.htm "WriteAsync 方法 (String, String)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| KukaAvarProxyNetWriteAsync 方法 (String, Byte) |

根据Kuka机器人的变量名称，写入原始的数据内容  
Write the original data content according to the variable name of the Kuka robot

**命名空间：**
 [HslCommunication.Robot.KUKA](185b17e5-ca64-0bf5-9521-90aab446dd88.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public Task<OperateResult> WriteAsync(
	string address,
	byte[] value
)
```

```
Public Function WriteAsync ( 
	address As String,
	value As Byte()
) As Task(Of OperateResult)
```

```
public:
virtual Task<OperateResult^>^ WriteAsync(
	String^ address, 
	array<unsigned char>^ value
) sealed
```

```
abstract WriteAsync : 
        address : string * 
        value : byte[] -> Task<OperateResult> 
override WriteAsync : 
        address : string * 
        value : byte[] -> Task<OperateResult>
```

#### 参数

address
:   类型：SystemString  
    变量名称

value
:   类型：SystemByte  
    原始的字节数据信息

#### 返回值

类型：Task[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否成功的写入

#### 实现

[IRobotNetWriteAsync(String, Byte)](17f6ba6b-5f84-0d7b-7f35-f0ed21ed0d22.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[KukaAvarProxyNet 类](24a1f5b9-8b90-7fe9-b45f-6e625494f1d5.htm)

[WriteAsync 重载](4c251649-c7dc-46be-6996-d9eca2f77047.htm)

[HslCommunication.Robot.KUKA 命名空间](185b17e5-ca64-0bf5-9521-90aab446dd88.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## WriteAsync 方法 (String, String)

[原文連結](http://api.hslcommunication.cn/html/ab50cf71-4c25-d98b-23b9-298fdd41dd8a.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.KUKA](../html/185b17e5-ca64-0bf5-9521-90aab446dd88.htm "HslCommunication.Robot.KUKA")

[KukaAvarProxyNet 类](../html/24a1f5b9-8b90-7fe9-b45f-6e625494f1d5.htm "KukaAvarProxyNet 类")

[KukaAvarProxyNet 方法](../html/820392c9-645f-c6f1-db0a-1ff5b1a1a86a.htm "KukaAvarProxyNet 方法")

[WriteAsync 方法](../html/4c251649-c7dc-46be-6996-d9eca2f77047.htm "WriteAsync 方法 ")

[WriteAsync 方法 (String, Byte[])](../html/b4cd9554-b881-5e87-194c-39abd678bd95.htm "WriteAsync 方法 (String, Byte[])")

[WriteAsync 方法 (String, String)](../html/ab50cf71-4c25-d98b-23b9-298fdd41dd8a.htm "WriteAsync 方法 (String, String)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| KukaAvarProxyNetWriteAsync 方法 (String, String) |

根据Kuka机器人的变量名称，写入ANSI编码的字符串数据信息  
Writes ansi-encoded string data information based on the variable name of the Kuka robot

**命名空间：**
 [HslCommunication.Robot.KUKA](185b17e5-ca64-0bf5-9521-90aab446dd88.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public Task<OperateResult> WriteAsync(
	string address,
	string value
)
```

```
Public Function WriteAsync ( 
	address As String,
	value As String
) As Task(Of OperateResult)
```

```
public:
virtual Task<OperateResult^>^ WriteAsync(
	String^ address, 
	String^ value
) sealed
```

```
abstract WriteAsync : 
        address : string * 
        value : string -> Task<OperateResult> 
override WriteAsync : 
        address : string * 
        value : string -> Task<OperateResult>
```

#### 参数

address
:   类型：SystemString  
    变量名称

value
:   类型：SystemString  
    ANSI编码的字符串

#### 返回值

类型：Task[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否成功的写入

#### 实现

[IRobotNetWriteAsync(String, String)](bfe4625b-e39b-3e1b-5269-26f0af34ae17.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[KukaAvarProxyNet 类](24a1f5b9-8b90-7fe9-b45f-6e625494f1d5.htm)

[WriteAsync 重载](4c251649-c7dc-46be-6996-d9eca2f77047.htm)

[HslCommunication.Robot.KUKA 命名空间](185b17e5-ca64-0bf5-9521-90aab446dd88.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## KukaAvarProxyNet 字段

[原文連結](http://api.hslcommunication.cn/html/5e7871de-3ec3-a3cb-e1ea-e82fd2737fab.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.KUKA](../html/185b17e5-ca64-0bf5-9521-90aab446dd88.htm "HslCommunication.Robot.KUKA")

[KukaAvarProxyNet 类](../html/24a1f5b9-8b90-7fe9-b45f-6e625494f1d5.htm "KukaAvarProxyNet 类")

[KukaAvarProxyNet 构造函数](../html/5f800939-b86d-c3da-6054-bb7345036f9e.htm "KukaAvarProxyNet 构造函数 ")

[KukaAvarProxyNet 属性](../html/dc604acd-e69c-c07d-0762-7f25045a24b8.htm "KukaAvarProxyNet 属性")

[KukaAvarProxyNet 方法](../html/820392c9-645f-c6f1-db0a-1ff5b1a1a86a.htm "KukaAvarProxyNet 方法")

[KukaAvarProxyNet 字段](../html/5e7871de-3ec3-a3cb-e1ea-e82fd2737fab.htm "KukaAvarProxyNet 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| KukaAvarProxyNet 字段 |

[KukaAvarProxyNet](24a1f5b9-8b90-7fe9-b45f-6e625494f1d5.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)字段

|  | 名称 | 说明 |
| --- | --- | --- |
| 受保护的字段 | [fileCacheSize](2328fb89-cfda-7fa9-88d0-98f1ae3c7c54.htm) | 文件传输的时候的缓存大小，直接影响传输的速度，值越大，传输速度越快，越占内存，默认为100K大小  The size of the cache during file transfer directly affects the speed of the transfer. The larger the value, the faster the transfer speed and the more memory it takes. The default size is 100K. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的字段 | [isPersistentConn](17c281db-b635-a496-ef0e-8e91770c94be.htm) | 是否是长连接的状态  Whether it is a long connection state (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 受保护的字段 | [isUseAccountCertificate](c3f9fdea-29af-2f88-8f6a-f4651318a8a6.htm) | 是否使用账号登录，这个账户登录的功能是HSL组件创建的服务器特有的功能。  Whether to log in using an account. The function of this account login is a server-specific function created by the  HSL  component. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 受保护的字段 | [LogMsgFormatBinary](a00be7ba-6c1e-fd46-c68c-e004433341e5.htm) | 设置日志记录报文是否二进制，如果为False，那就使用ASCII码  Set whether the log message is binary, if it is False, then use ASCII code (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 受保护的字段 | [pipeSocket](d4efd9b0-ad04-ce88-5dde-843794e4e3fa.htm) | 当前的网络的管道信息 (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[KukaAvarProxyNet 类](24a1f5b9-8b90-7fe9-b45f-6e625494f1d5.htm)

[HslCommunication.Robot.KUKA 命名空间](185b17e5-ca64-0bf5-9521-90aab446dd88.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## KukaTcpNet 类

[原文連結](http://api.hslcommunication.cn/html/12f49624-a008-ac45-d356-6af762d5cc4c.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.KUKA](../html/185b17e5-ca64-0bf5-9521-90aab446dd88.htm "HslCommunication.Robot.KUKA")

[KukaTcpNet 类](../html/12f49624-a008-ac45-d356-6af762d5cc4c.htm "KukaTcpNet 类")

[KukaTcpNet 构造函数](../html/cee3d8db-13ef-cb45-3b84-777bfe79c535.htm "KukaTcpNet 构造函数 ")

[KukaTcpNet 属性](../html/6bd4fdff-252d-96e2-d206-4600b799da8b.htm "KukaTcpNet 属性")

[KukaTcpNet 方法](../html/b0df9dbe-319e-b403-b1a2-b18b345d0728.htm "KukaTcpNet 方法")

[KukaTcpNet 字段](../html/36242a04-fc58-ecd5-a572-4db7621d6d78.htm "KukaTcpNet 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| KukaTcpNet 类 |

Kuka机器人的数据交互类，通讯支持的条件为KUKA 的 TCP通讯

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  [HslCommunication.Core.NetNetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)  
    [HslCommunication.Core.NetNetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)  
      HslCommunication.Robot.KUKAKukaTcpNet

**命名空间：**
 [HslCommunication.Robot.KUKA](185b17e5-ca64-0bf5-9521-90aab446dd88.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class KukaTcpNet : NetworkDoubleBase, IRobotNet
```

```
Public Class KukaTcpNet
	Inherits NetworkDoubleBase
	Implements IRobotNet
```

```
public ref class KukaTcpNet : public NetworkDoubleBase, 
	IRobotNet
```

```
type KukaTcpNet =  
    class
        inherit NetworkDoubleBase
        interface IRobotNet
    end
```

KukaTcpNet 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [KukaTcpNet](15e39d65-1299-147c-dba7-8e60ec37c0fa.htm) | 实例化一个默认的对象  Instantiate a default object |
| 公共方法 | [KukaTcpNet(String, Int32)](95666452-5792-d70d-28ca-1f333b2f089c.htm) | 实例化一个默认的Kuka机器人对象，并指定IP地址和端口号，端口号通常为9999  Instantiate a default Kuka robot object and specify the IP address and port number, usually 9999 |

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
| 公共方法静态成员 | [BuildReadCommands(String)](000b663f-b787-0b88-4487-ba8853851296.htm) | 构建读取变量的报文命令 |
| 公共方法静态成员 | [BuildReadCommands(String)](a3e56815-ccd1-0a07-ab03-09391f01a81a.htm) | 构建读取变量的报文命令 |
| 公共方法静态成员 | [BuildWriteCommands(String, String)](d709da01-1bc3-2179-4886-412a92997137.htm) | 构建写入变量的报文命令 |
| 公共方法静态成员 | [BuildWriteCommands(String, String)](bca4cdb6-0084-3dc5-30e1-f0f1f3c3fbaf.htm) | 构建写入变量的报文命令 |
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
| 受保护的方法 | [DecideWhetherQAMessage](77091dac-82d7-8c67-d691-31ae5b2aaee4.htm) | 决定当前的消息是否是应答机制的消息内容，需要在客户端进行重写实现，如果是应答机制，返回 True, 否则返回 False  To determine whether the current message is the message content of the response mechanism, it needs to be rewritten on the client side. If it is the response mechanism, return True, otherwise return False (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 受保护的方法 | [DeleteFileByName](59afa4a4-3f3c-9831-26ee-c7c13a213bf7.htm) | 删除一个指定的文件，如果文件不存在，直接返回 True，如果文件存在则直接删除，删除成功返回 True，如果发生了异常，返回False  Delete a specified file, if the file does not exist, return True directly, if the file exists, delete it directly, if the deletion is successful, return True, if an exception occurs, return  False (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 公共方法 | [Dispose](8a86a200-fbe0-3b45-0fed-c551405c9b59.htm) | 释放当前的资源，如果调用了本方法，那么该对象再使用的时候，需要重新实例化。  Release the current resource. If this method is called, the object needs to be instantiated again when it is used again. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 受保护的方法 | [Dispose(Boolean)](32f03275-9ae0-525c-84bc-6e2515d11c12.htm) | 释放当前的资源，并自动关闭长连接，如果设置了的话 (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | [ExtraAfterReadFromCoreServer](9f8357e9-dff4-1002-9368-0744cdedff29.htm) | 和服务器交互完成的时候调用的方法，可以根据读写结果进行一些额外的操作，具体的操作需要根据实际的需求来重写实现  The method called when the interaction with the server is completed can perform some additional operations based on the read and write results. The specific operations need to be rewritten according to actual needs. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 受保护的方法代码示例 | [ExtraOnDisconnect](b8569c75-5e55-cd55-5251-207465a8ba39.htm) | 根据实际的协议选择是否重写本方法，有些协议在断开连接之前，需要发送一些报文来关闭当前的网络通道  Select whether to rewrite this method according to the actual protocol. Some protocols need to send some packets to close the current network channel before disconnecting. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 受保护的方法代码示例 | [ExtraOnDisconnectAsync](86aa00d8-b2d6-ada9-5926-2b7a2237fe6b.htm) | 根据实际的协议选择是否重写本方法，有些协议在断开连接之前，需要发送一些报文来关闭当前的网络通道  Select whether to rewrite this method according to the actual protocol. Some protocols need to send some packets to close the current network channel before disconnecting. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 受保护的方法 | [GetAvailableSocket](e86a6fc1-eb91-8de9-5ff8-799b0ea3786b.htm) | 获取本次操作的可用的网络通道，如果是短连接，就重新生成一个新的网络通道，如果是长连接，就复用当前的网络通道。  Obtain the available network channels for this operation. If it is a short connection, a new network channel is regenerated. If it is a long connection, the current network channel is reused. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 受保护的方法 | [GetAvailableSocketAsync](40cf9ce7-a5e6-5431-2e4b-56e3b7c98346.htm) | 获取本次操作的可用的网络通道，如果是短连接，就重新生成一个新的网络通道，如果是长连接，就复用当前的网络通道。  Obtain the available network channels for this operation. If it is a short connection, a new network channel is regenerated. If it is a long connection, the current network channel is reused. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 受保护的方法 | [GetNewNetMessage](481a40e0-8539-5fc8-3e2b-c25bcc952697.htm) | 获取一个新的消息对象的方法，需要在继承类里面进行重写  The method to get a new message object needs to be overridden in the inheritance class (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | [GetPipeSocket](19904680-372e-e0a6-1b27-b4180e759668.htm) | 获取当前用于通信的管道信息  Get the current pipe information used for communication (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法代码示例 | [InitializationOnConnect](40d5eb26-9bca-11cf-42a1-55af6ea5aebd.htm) | 根据实际的协议选择是否重写本方法，有些协议在创建连接之后，需要进行一些初始化的信号握手，才能最终建立网络通道。  Whether to rewrite this method is based on the actual protocol. Some protocols require some initial signal handshake to establish a network channel after the connection is created. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 受保护的方法代码示例 | [InitializationOnConnectAsync](14cac21b-4001-da7d-09a6-941eba868e14.htm) | 根据实际的协议选择是否重写本方法，有些协议在创建连接之后，需要进行一些初始化的信号握手，才能最终建立网络通道。  Whether to rewrite this method is based on the actual protocol. Some protocols require some initial signal handshake to establish a network channel after the connection is created. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | [IpAddressPing](d9611e18-4033-de95-444e-b91ae0a81709.htm) | 对当前设备的IP地址进行PING的操作，返回PING的结果，正常来说，返回Success  PING the IP address of the current device and return the PING result. Normally, it returns Success (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法 | [PackCommandWithHeader](f96881d6-e59f-9e16-07e3-1188eee1a934.htm) | 对当前的命令进行打包处理，通常是携带命令头内容，标记当前的命令的长度信息，需要进行重写，否则默认不打包  The current command is packaged, usually carrying the content of the command header, marking the length of the current command, and it needs to be rewritten, otherwise it is not packaged by default (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | [Read](b5f857a6-6bf8-d640-07f7-2bc51e53fc47.htm) | 读取Kuka机器人的数据内容，根据输入的变量名称来读取  Read the data content of the Kuka robot according to the input variable name |
| 公共方法 | [ReadAsync](c82bc37f-fa1a-ebeb-50fc-3054ed60b086.htm) | 读取Kuka机器人的数据内容，根据输入的变量名称来读取  Read the data content of the Kuka robot according to the input variable name |
| 公共方法代码示例 | [ReadFromCoreServer(Byte)](b3e14431-d773-0d1a-f453-a3ad60648c29.htm) | 将数据发送到当前的网络通道中，并从网络通道中接收一个[INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm)指定的完整的报文，网络通道将根据[GetAvailableSocket](e86a6fc1-eb91-8de9-5ff8-799b0ea3786b.htm)方法自动获取，本方法是线程安全的。  Send data to the current network channel and receive a complete message specified by [INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm) from the network channel. The network channel will be automatically obtained according to the [GetAvailableSocket](e86a6fc1-eb91-8de9-5ff8-799b0ea3786b.htm) method This method is thread-safe. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | [ReadFromCoreServer(IEnumerableByte)](580dae22-fad5-5816-d58f-43068e89e419.htm) | 将多个数据报文按顺序发到设备，并从设备接收返回的数据内容，然后拼接成一个Byte[]信息，需要重写[UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm)方法才能返回正确的结果。  Send multiple data packets to the device in sequence, and receive the returned data content from the device, and then splicing them into a Byte[] message, you need to rewrite [UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm) method to return the correct result. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServer(Byte, Boolean, Boolean)](1a33d974-805d-f149-3665-eb8d8a99ea84.htm) | 将数据发送到当前的网络通道中，并从网络通道中接收一个[INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm)指定的完整的报文，网络通道将根据[GetAvailableSocket](e86a6fc1-eb91-8de9-5ff8-799b0ea3786b.htm)方法自动获取，本方法是线程安全的。  Send data to the current network channel and receive a complete message specified by [INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm) from the network channel. The network channel will be automatically obtained according to the [GetAvailableSocket](e86a6fc1-eb91-8de9-5ff8-799b0ea3786b.htm) method This method is thread-safe. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServer(Socket, Byte, Boolean, Boolean)](df34b731-3b68-3e23-5e5e-4ab4c4316ff8.htm) | 将数据报文发送指定的网络通道上，根据当前指定的[INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm)类型，返回一条完整的数据指令  Sends a data message to the specified network channel, and returns a complete data command according to the currently specified [INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm) type (重写 [NetworkDoubleBaseReadFromCoreServer(Socket, Byte, Boolean, Boolean)](3b72b289-b201-7548-0981-ac049c9a4135.htm).) |
| 公共方法代码示例 | [ReadFromCoreServerAsync(Byte)](ec51a18e-0f5d-dbe1-c229-d62e8760ac97.htm) | 将数据发送到当前的网络通道中，并从网络通道中接收一个[INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm)指定的完整的报文，网络通道将根据[GetAvailableSocket](e86a6fc1-eb91-8de9-5ff8-799b0ea3786b.htm)方法自动获取，本方法是线程安全的。  Send data to the current network channel and receive a complete message specified by [INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm) from the network channel. The network channel will be automatically obtained according to the [GetAvailableSocket](e86a6fc1-eb91-8de9-5ff8-799b0ea3786b.htm) method This method is thread-safe. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | [ReadFromCoreServerAsync(IEnumerableByte)](0fce97fd-8441-6b54-0877-b8dfacc738e3.htm) | 将多个数据报文按顺序发到设备，并从设备接收返回的数据内容，然后拼接成一个Byte[]信息，需要重写[UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm)方法才能返回正确的结果。  Send multiple data packets to the device in sequence, and receive the returned data content from the device, and then splicing them into a Byte[] message, you need to rewrite [UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm) method to return the correct result. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServerAsync(Byte, Boolean, Boolean)](7525b379-04c2-8dac-b451-a8b5413d6747.htm) | 将数据发送到当前的网络通道中，并从网络通道中接收一个[INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm)指定的完整的报文，网络通道将根据[GetAvailableSocket](e86a6fc1-eb91-8de9-5ff8-799b0ea3786b.htm)方法自动获取，本方法是线程安全的。  Send data to the current network channel and receive a complete message specified by [INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm) from the network channel. The network channel will be automatically obtained according to the [GetAvailableSocket](e86a6fc1-eb91-8de9-5ff8-799b0ea3786b.htm) method This method is thread-safe. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServerAsync(Socket, Byte, Boolean, Boolean)](25f70345-396d-b26b-c8d6-41ca0ca0e62a.htm) | 将数据报文发送指定的网络通道上，根据当前指定的[INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm)类型，返回一条完整的数据指令  Sends a data message to the specified network channel, and returns a complete data command according to the currently specified [INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm) type (重写 [NetworkDoubleBaseReadFromCoreServerAsync(Socket, Byte, Boolean, Boolean)](07db4a00-b548-8bd6-5de6-12deb91f5311.htm).) |
| 公共方法 | [ReadString](c8680435-6f44-adcc-7bd7-f5002035f45c.htm) | 读取Kuka机器人的所有的数据信息，返回字符串信息，解码方式为UTF8，需要指定变量名称  Read all the data information of the Kuka robot, return the string information, decode by ANSI, need to specify the variable name |
| 公共方法 | [ReadStringAsync](fe9ebb7f-56ba-e6a2-4827-8306e360bb11.htm) | 读取Kuka机器人的所有的数据信息，返回字符串信息，解码方式为UTF8，需要指定变量名称  Read all the data information of the Kuka robot, return the string information, decode by ANSI, need to specify the variable name |
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
| 公共方法 | [ResetProgram](a77a2d41-d618-14d4-67db-d6bc64fade87.htm) | 复位当前的程序  Reset current program |
| 公共方法 | [ResetProgramAsync](0eb606b6-3e63-2c56-7884-1294f3fd6908.htm) | 复位当前的程序  Reset current program |
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
| 公共方法 | [SetLoginAccount](725afc25-6e3f-bd7d-25ed-580270c8fc56.htm) | 设置当前的登录的账户名和密码信息，并启用账户验证的功能，账户名为空时设置不生效  Set the current login account name and password information, and enable the account verification function. The account name setting will not take effect when it is empty (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法代码示例 | [SetPersistentConnection](d8e8df9f-cbe9-6863-54e9-aa3b38906b26.htm) | 在读取数据之前可以调用本方法将客户端设置为长连接模式，相当于跳过了ConnectServer的结果验证，对异形客户端无效，当第一次进行通信时再进行创建连接请求。  Before reading the data, you can call this method to set the client to the long connection mode, which is equivalent to skipping the result verification of ConnectServer, and it is invalid for the alien client. When the first communication is performed, the connection creation request is performed. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | [SetPipeSocket](b6a7a243-2a6d-c4ea-45b7-bf980450a2bd.htm) | 设置一个新的网络管道，一般来说不需要调用本方法，当多个网口设备共用一个网络连接时才需要使用本方法进行设置共享的管道。  To set up a new network channel, generally speaking, you do not need to call this method. This method is only needed to set up a shared channel when multiple network port devices share a network connection. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | [StartProgram](bdc9cf2e-32d7-bf04-4b36-59d801353fa4.htm) | 启动机器人的指定的程序  Start the specified program of the robot |
| 公共方法 | [StartProgramAsync](23f3a9ea-bb62-8464-c373-0f39402faace.htm) | 启动机器人的指定的程序  Start the specified program of the robot |
| 公共方法 | [StopProgram](88961e52-cac2-7d1f-44eb-50020640dbb7.htm) | 停止当前的程序  Stop current program |
| 公共方法 | [StopProgramAsync](b7f37f2f-2078-4f47-1cda-786714052fe4.htm) | 停止当前的程序  Stop current program |
| 公共方法 | [ToString](5d2084ed-0b08-b8e6-a666-194365eb5ab8.htm) | (重写 [NetworkDoubleBaseToString](85a8532f-8fc5-886c-9d76-70844e0a7d64.htm).) |
| 公共方法 | [UnpackResponseContent](32bad9a0-b59d-c171-ce85-74055edfc597.htm) | 根据对方返回的报文命令，对命令进行基本的拆包，例如各种Modbus协议拆包为统一的核心报文，还支持对报文的验证  According to the message command returned by the other party, the command is basically unpacked, for example, various Modbus protocols are unpacked into a unified core message, and the verification of the message is also supported (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | [Write(String, Byte)](9d930e57-aae7-23c4-5d09-bf5d102fc8cb.htm) | 根据Kuka机器人的变量名称，写入原始的数据内容  Write the original data content according to the variable name of the Kuka robot |
| 公共方法 | [Write(String, String)](34bb38f9-eea8-fbd3-93c7-820a88a93b5f.htm) | 根据Kuka机器人的变量名称，写入UTF8编码的字符串数据信息  Writes ansi-encoded string data information based on the variable name of the Kuka robot |
| 公共方法 | [Write(String, String)](fb0e5a1c-1c0a-1788-bf73-63fd58b018c3.htm) | 根据Kuka机器人的变量名称，写入多个UTF8编码的字符串数据信息  Write multiple UTF8 encoded string data information according to the variable name of the Kuka robot |
| 公共方法 | [WriteAsync(String, Byte)](72461f2e-93fb-b952-68ac-1e49b7eee5c9.htm) | 根据Kuka机器人的变量名称，写入原始的数据内容  Write the original data content according to the variable name of the Kuka robot |
| 公共方法 | [WriteAsync(String, String)](d409afe2-c4f1-065d-4d50-96a71e5bbfa8.htm) | 根据Kuka机器人的变量名称，写入UTF8编码的字符串数据信息  Writes ansi-encoded string data information based on the variable name of the Kuka robot |
| 公共方法 | [WriteAsync(String, String)](501ef94b-8805-e676-59b0-2fa5501d9691.htm) | 根据Kuka机器人的变量名称，写入多个UTF8编码的字符串数据信息  Write multiple UTF8 encoded string data information according to the variable name of the Kuka robot |
| 受保护的方法 | [WriteStreamFromSocket](97bf95ea-51fe-8f2e-5da6-1c83766aaa78.htm) | 从套接字中接收所有的数据然后写入到指定的流当中去，需要指定数据的长度，支持按照百分比进行进度报告  Receives all data from the socket and writes it to the specified stream. The length of the data needs to be specified, and progress reporting is supported in percentage. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [WriteStreamFromSocketAsync](b4cf7f6d-1e0c-a266-4c49-683fdb14a922.htm) | 从套接字中接收所有的数据然后写入到指定的流当中去，需要指定数据的长度，支持按照百分比进行进度报告  Receives all data from the socket and writes it to the specified stream. The length of the data needs to be specified, and progress reporting is supported in percentage. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |

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

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.Robot.KUKA 命名空间](185b17e5-ca64-0bf5-9521-90aab446dd88.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## KukaTcpNet 构造函数 

[原文連結](http://api.hslcommunication.cn/html/cee3d8db-13ef-cb45-3b84-777bfe79c535.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.KUKA](../html/185b17e5-ca64-0bf5-9521-90aab446dd88.htm "HslCommunication.Robot.KUKA")

[KukaTcpNet 类](../html/12f49624-a008-ac45-d356-6af762d5cc4c.htm "KukaTcpNet 类")

[KukaTcpNet 构造函数](../html/cee3d8db-13ef-cb45-3b84-777bfe79c535.htm "KukaTcpNet 构造函数 ")

[KukaTcpNet 构造函数](../html/15e39d65-1299-147c-dba7-8e60ec37c0fa.htm "KukaTcpNet 构造函数 ")

[KukaTcpNet 构造函数 (String, Int32)](../html/95666452-5792-d70d-28ca-1f333b2f089c.htm "KukaTcpNet 构造函数 (String, Int32)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| KukaTcpNet 构造函数 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [KukaTcpNet](15e39d65-1299-147c-dba7-8e60ec37c0fa.htm) | 实例化一个默认的对象  Instantiate a default object |
| 公共方法 | [KukaTcpNet(String, Int32)](95666452-5792-d70d-28ca-1f333b2f089c.htm) | 实例化一个默认的Kuka机器人对象，并指定IP地址和端口号，端口号通常为9999  Instantiate a default Kuka robot object and specify the IP address and port number, usually 9999 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[KukaTcpNet 类](12f49624-a008-ac45-d356-6af762d5cc4c.htm)

[HslCommunication.Robot.KUKA 命名空间](185b17e5-ca64-0bf5-9521-90aab446dd88.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## KukaTcpNet 构造函数 

[原文連結](http://api.hslcommunication.cn/html/15e39d65-1299-147c-dba7-8e60ec37c0fa.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.KUKA](../html/185b17e5-ca64-0bf5-9521-90aab446dd88.htm "HslCommunication.Robot.KUKA")

[KukaTcpNet 类](../html/12f49624-a008-ac45-d356-6af762d5cc4c.htm "KukaTcpNet 类")

[KukaTcpNet 构造函数](../html/cee3d8db-13ef-cb45-3b84-777bfe79c535.htm "KukaTcpNet 构造函数 ")

[KukaTcpNet 构造函数](../html/15e39d65-1299-147c-dba7-8e60ec37c0fa.htm "KukaTcpNet 构造函数 ")

[KukaTcpNet 构造函数 (String, Int32)](../html/95666452-5792-d70d-28ca-1f333b2f089c.htm "KukaTcpNet 构造函数 (String, Int32)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| KukaTcpNet 构造函数 |

实例化一个默认的对象  
Instantiate a default object

**命名空间：**
 [HslCommunication.Robot.KUKA](185b17e5-ca64-0bf5-9521-90aab446dd88.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public KukaTcpNet()
```

```
Public Sub New
```

```
public:
KukaTcpNet()
```

```
new : unit -> KukaTcpNet
```

![](../icons/SectionExpanded.png)参见

#### 引用

[KukaTcpNet 类](12f49624-a008-ac45-d356-6af762d5cc4c.htm)

[KukaTcpNet 重载](cee3d8db-13ef-cb45-3b84-777bfe79c535.htm)

[HslCommunication.Robot.KUKA 命名空间](185b17e5-ca64-0bf5-9521-90aab446dd88.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## KukaTcpNet 构造函数 (String, Int32)

[原文連結](http://api.hslcommunication.cn/html/95666452-5792-d70d-28ca-1f333b2f089c.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.KUKA](../html/185b17e5-ca64-0bf5-9521-90aab446dd88.htm "HslCommunication.Robot.KUKA")

[KukaTcpNet 类](../html/12f49624-a008-ac45-d356-6af762d5cc4c.htm "KukaTcpNet 类")

[KukaTcpNet 构造函数](../html/cee3d8db-13ef-cb45-3b84-777bfe79c535.htm "KukaTcpNet 构造函数 ")

[KukaTcpNet 构造函数](../html/15e39d65-1299-147c-dba7-8e60ec37c0fa.htm "KukaTcpNet 构造函数 ")

[KukaTcpNet 构造函数 (String, Int32)](../html/95666452-5792-d70d-28ca-1f333b2f089c.htm "KukaTcpNet 构造函数 (String, Int32)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| KukaTcpNet 构造函数 (String, Int32) |

实例化一个默认的Kuka机器人对象，并指定IP地址和端口号，端口号通常为9999  
Instantiate a default Kuka robot object and specify the IP address and port number, usually 9999

**命名空间：**
 [HslCommunication.Robot.KUKA](185b17e5-ca64-0bf5-9521-90aab446dd88.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public KukaTcpNet(
	string ipAddress,
	int port
)
```

```
Public Sub New ( 
	ipAddress As String,
	port As Integer
)
```

```
public:
KukaTcpNet(
	String^ ipAddress, 
	int port
)
```

```
new : 
        ipAddress : string * 
        port : int -> KukaTcpNet
```

#### 参数

ipAddress
:   类型：SystemString  
    Ip地址

port
:   类型：SystemInt32  
    端口号

![](../icons/SectionExpanded.png)参见

#### 引用

[KukaTcpNet 类](12f49624-a008-ac45-d356-6af762d5cc4c.htm)

[KukaTcpNet 重载](cee3d8db-13ef-cb45-3b84-777bfe79c535.htm)

[HslCommunication.Robot.KUKA 命名空间](185b17e5-ca64-0bf5-9521-90aab446dd88.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## KukaTcpNet 属性

[原文連結](http://api.hslcommunication.cn/html/6bd4fdff-252d-96e2-d206-4600b799da8b.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.KUKA](../html/185b17e5-ca64-0bf5-9521-90aab446dd88.htm "HslCommunication.Robot.KUKA")

[KukaTcpNet 类](../html/12f49624-a008-ac45-d356-6af762d5cc4c.htm "KukaTcpNet 类")

[KukaTcpNet 构造函数](../html/cee3d8db-13ef-cb45-3b84-777bfe79c535.htm "KukaTcpNet 构造函数 ")

[KukaTcpNet 属性](../html/6bd4fdff-252d-96e2-d206-4600b799da8b.htm "KukaTcpNet 属性")

[KukaTcpNet 方法](../html/b0df9dbe-319e-b403-b1a2-b18b345d0728.htm "KukaTcpNet 方法")

[KukaTcpNet 字段](../html/36242a04-fc58-ecd5-a572-4db7621d6d78.htm "KukaTcpNet 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| KukaTcpNet 属性 |

[KukaTcpNet](12f49624-a008-ac45-d356-6af762d5cc4c.htm) 类型公开以下成员。

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

[KukaTcpNet 类](12f49624-a008-ac45-d356-6af762d5cc4c.htm)

[HslCommunication.Robot.KUKA 命名空间](185b17e5-ca64-0bf5-9521-90aab446dd88.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## BuildReadCommands 方法 

[原文連結](http://api.hslcommunication.cn/html/55e34028-68f5-6cdd-447a-cf6a3584113d.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.KUKA](../html/185b17e5-ca64-0bf5-9521-90aab446dd88.htm "HslCommunication.Robot.KUKA")

[KukaTcpNet 类](../html/12f49624-a008-ac45-d356-6af762d5cc4c.htm "KukaTcpNet 类")

[KukaTcpNet 方法](../html/b0df9dbe-319e-b403-b1a2-b18b345d0728.htm "KukaTcpNet 方法")

[BuildReadCommands 方法](../html/55e34028-68f5-6cdd-447a-cf6a3584113d.htm "BuildReadCommands 方法 ")

[BuildReadCommands 方法 (String)](../html/000b663f-b787-0b88-4487-ba8853851296.htm "BuildReadCommands 方法 (String)")

[BuildReadCommands 方法 (String[])](../html/a3e56815-ccd1-0a07-ab03-09391f01a81a.htm "BuildReadCommands 方法 (String[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| KukaTcpNetBuildReadCommands 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法静态成员 | [BuildReadCommands(String)](000b663f-b787-0b88-4487-ba8853851296.htm) | 构建读取变量的报文命令 |
| 公共方法静态成员 | [BuildReadCommands(String)](a3e56815-ccd1-0a07-ab03-09391f01a81a.htm) | 构建读取变量的报文命令 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[KukaTcpNet 类](12f49624-a008-ac45-d356-6af762d5cc4c.htm)

[HslCommunication.Robot.KUKA 命名空间](185b17e5-ca64-0bf5-9521-90aab446dd88.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## BuildReadCommands 方法 (String)

[原文連結](http://api.hslcommunication.cn/html/000b663f-b787-0b88-4487-ba8853851296.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.KUKA](../html/185b17e5-ca64-0bf5-9521-90aab446dd88.htm "HslCommunication.Robot.KUKA")

[KukaTcpNet 类](../html/12f49624-a008-ac45-d356-6af762d5cc4c.htm "KukaTcpNet 类")

[KukaTcpNet 方法](../html/b0df9dbe-319e-b403-b1a2-b18b345d0728.htm "KukaTcpNet 方法")

[BuildReadCommands 方法](../html/55e34028-68f5-6cdd-447a-cf6a3584113d.htm "BuildReadCommands 方法 ")

[BuildReadCommands 方法 (String)](../html/000b663f-b787-0b88-4487-ba8853851296.htm "BuildReadCommands 方法 (String)")

[BuildReadCommands 方法 (String[])](../html/a3e56815-ccd1-0a07-ab03-09391f01a81a.htm "BuildReadCommands 方法 (String[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| KukaTcpNetBuildReadCommands 方法 (String) |

构建读取变量的报文命令

**命名空间：**
 [HslCommunication.Robot.KUKA](185b17e5-ca64-0bf5-9521-90aab446dd88.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static string BuildReadCommands(
	string address
)
```

```
Public Shared Function BuildReadCommands ( 
	address As String
) As String
```

```
public:
static String^ BuildReadCommands(
	String^ address
)
```

```
static member BuildReadCommands : 
        address : string -> string 
```

#### 参数

address
:   类型：SystemString  
    地址信息

#### 返回值

类型：String  
报文内容

![](../icons/SectionExpanded.png)参见

#### 引用

[KukaTcpNet 类](12f49624-a008-ac45-d356-6af762d5cc4c.htm)

[BuildReadCommands 重载](55e34028-68f5-6cdd-447a-cf6a3584113d.htm)

[HslCommunication.Robot.KUKA 命名空间](185b17e5-ca64-0bf5-9521-90aab446dd88.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## BuildReadCommands 方法 (String[])

[原文連結](http://api.hslcommunication.cn/html/a3e56815-ccd1-0a07-ab03-09391f01a81a.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.KUKA](../html/185b17e5-ca64-0bf5-9521-90aab446dd88.htm "HslCommunication.Robot.KUKA")

[KukaTcpNet 类](../html/12f49624-a008-ac45-d356-6af762d5cc4c.htm "KukaTcpNet 类")

[KukaTcpNet 方法](../html/b0df9dbe-319e-b403-b1a2-b18b345d0728.htm "KukaTcpNet 方法")

[BuildReadCommands 方法](../html/55e34028-68f5-6cdd-447a-cf6a3584113d.htm "BuildReadCommands 方法 ")

[BuildReadCommands 方法 (String)](../html/000b663f-b787-0b88-4487-ba8853851296.htm "BuildReadCommands 方法 (String)")

[BuildReadCommands 方法 (String[])](../html/a3e56815-ccd1-0a07-ab03-09391f01a81a.htm "BuildReadCommands 方法 (String[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| KukaTcpNetBuildReadCommands 方法 (String) |

构建读取变量的报文命令

**命名空间：**
 [HslCommunication.Robot.KUKA](185b17e5-ca64-0bf5-9521-90aab446dd88.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static string BuildReadCommands(
	string[] address
)
```

```
Public Shared Function BuildReadCommands ( 
	address As String()
) As String
```

```
public:
static String^ BuildReadCommands(
	array<String^>^ address
)
```

```
static member BuildReadCommands : 
        address : string[] -> string 
```

#### 参数

address
:   类型：SystemString  
    地址信息

#### 返回值

类型：String  
报文内容

![](../icons/SectionExpanded.png)参见

#### 引用

[KukaTcpNet 类](12f49624-a008-ac45-d356-6af762d5cc4c.htm)

[BuildReadCommands 重载](55e34028-68f5-6cdd-447a-cf6a3584113d.htm)

[HslCommunication.Robot.KUKA 命名空间](185b17e5-ca64-0bf5-9521-90aab446dd88.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## BuildWriteCommands 方法 

[原文連結](http://api.hslcommunication.cn/html/33a5a02a-d180-6fe4-4c71-1c54babcc316.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.KUKA](../html/185b17e5-ca64-0bf5-9521-90aab446dd88.htm "HslCommunication.Robot.KUKA")

[KukaTcpNet 类](../html/12f49624-a008-ac45-d356-6af762d5cc4c.htm "KukaTcpNet 类")

[KukaTcpNet 方法](../html/b0df9dbe-319e-b403-b1a2-b18b345d0728.htm "KukaTcpNet 方法")

[BuildWriteCommands 方法](../html/33a5a02a-d180-6fe4-4c71-1c54babcc316.htm "BuildWriteCommands 方法 ")

[BuildWriteCommands 方法 (String, String)](../html/d709da01-1bc3-2179-4886-412a92997137.htm "BuildWriteCommands 方法 (String, String)")

[BuildWriteCommands 方法 (String[], String[])](../html/bca4cdb6-0084-3dc5-30e1-f0f1f3c3fbaf.htm "BuildWriteCommands 方法 (String[], String[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| KukaTcpNetBuildWriteCommands 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法静态成员 | [BuildWriteCommands(String, String)](d709da01-1bc3-2179-4886-412a92997137.htm) | 构建写入变量的报文命令 |
| 公共方法静态成员 | [BuildWriteCommands(String, String)](bca4cdb6-0084-3dc5-30e1-f0f1f3c3fbaf.htm) | 构建写入变量的报文命令 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[KukaTcpNet 类](12f49624-a008-ac45-d356-6af762d5cc4c.htm)

[HslCommunication.Robot.KUKA 命名空间](185b17e5-ca64-0bf5-9521-90aab446dd88.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)