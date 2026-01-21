# HslCommunication - HslCommunication.Robot.YASKAWA

> 分類頁數: 30



---
## HslCommunication.Robot.YASKAWA

[原文連結](http://api.hslcommunication.cn/html/3d67023f-78aa-c5d9-85b6-f18f5dfbc176.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.YASKAWA](../html/3d67023f-78aa-c5d9-85b6-f18f5dfbc176.htm "HslCommunication.Robot.YASKAWA")

[YRC1000TcpNet 类](../html/cdcdc1cc-17cf-da27-7d20-ca6d6db1cd7e.htm "YRC1000TcpNet 类")

[YRCAlarmItem 类](../html/047b663a-2219-ab12-9934-1eab37060551.htm "YRCAlarmItem 类")

[YRCHighEthernet 类](../html/d16ad666-07d9-50cc-8295-4e76e789c768.htm "YRCHighEthernet 类")

[YRCRobotData 类](../html/233ef65f-4a8d-20e6-03ed-67c1cbefd2d6.htm "YRCRobotData 类")

[YRCType 枚举](../html/38ab618c-dbb6-4896-0126-5f10e97180ba.htm "YRCType 枚举")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslCommunication.Robot.YASKAWA 命名空间 |

[缺少 "N:HslCommunication.Robot.YASKAWA" 的 <summary> 文档]

![](../icons/SectionExpanded.png)类

|  | 类 | 说明 |
| --- | --- | --- |
| 公共类 | [YRC1000TcpNet](cdcdc1cc-17cf-da27-7d20-ca6d6db1cd7e.htm) | 安川机器人的Ethernet 服务器功能对应的客户端通讯类  Yaskawa robot's Ethernet server features a communication class |
| 公共类 | [YRCAlarmItem](047b663a-2219-ab12-9934-1eab37060551.htm) | 安川的报警信息 |
| 公共类 | [YRCHighEthernet](d16ad666-07d9-50cc-8295-4e76e789c768.htm) | 安川机器人的通信类，基于高速以太网的通信，基于UDP协议实现，默认端口10040，支持读写一些数据地址 |
| 公共类 | [YRCRobotData](233ef65f-4a8d-20e6-03ed-67c1cbefd2d6.htm) | 安川机器人的数据信息，其中 Re只在YRC100中有效，没有外部轴的系统， 7-12外部轴的值设定为「0」 |

![](../icons/SectionExpanded.png)枚举

|  | 枚举 | 说明 |
| --- | --- | --- |
| 公共枚举 | [YRCType](38ab618c-dbb6-4896-0126-5f10e97180ba.htm) | YRC机器人的类型 |

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## YRC1000TcpNet 类

[原文連結](http://api.hslcommunication.cn/html/cdcdc1cc-17cf-da27-7d20-ca6d6db1cd7e.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.YASKAWA](../html/3d67023f-78aa-c5d9-85b6-f18f5dfbc176.htm "HslCommunication.Robot.YASKAWA")

[YRC1000TcpNet 类](../html/cdcdc1cc-17cf-da27-7d20-ca6d6db1cd7e.htm "YRC1000TcpNet 类")

[YRC1000TcpNet 构造函数](../html/ebe3130d-b046-c1e9-5c6e-de2669433a40.htm "YRC1000TcpNet 构造函数 ")

[YRC1000TcpNet 属性](../html/bae7bc3e-546f-f9ee-b28d-7ad7c6eea531.htm "YRC1000TcpNet 属性")

[YRC1000TcpNet 方法](../html/7a26c354-eee0-b601-1d8a-51198e40e3fa.htm "YRC1000TcpNet 方法")

[YRC1000TcpNet 字段](../html/296507d2-3c83-4a11-77d9-21def100a3eb.htm "YRC1000TcpNet 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| YRC1000TcpNet 类 |

安川机器人的Ethernet 服务器功能对应的客户端通讯类  
Yaskawa robot's Ethernet server features a communication class

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  [HslCommunication.Core.NetNetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)  
    [HslCommunication.Core.NetNetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)  
      HslCommunication.Robot.YASKAWAYRC1000TcpNet

**命名空间：**
 [HslCommunication.Robot.YASKAWA](3d67023f-78aa-c5d9-85b6-f18f5dfbc176.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class YRC1000TcpNet : NetworkDoubleBase, 
	IRobotNet
```

```
Public Class YRC1000TcpNet
	Inherits NetworkDoubleBase
	Implements IRobotNet
```

```
public ref class YRC1000TcpNet : public NetworkDoubleBase, 
	IRobotNet
```

```
type YRC1000TcpNet =  
    class
        inherit NetworkDoubleBase
        interface IRobotNet
    end
```

YRC1000TcpNet 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [YRC1000TcpNet](ebe3130d-b046-c1e9-5c6e-de2669433a40.htm) | 指定机器人的ip地址及端口号来实例化对象  Specify the robot's IP address and port number to instantiate the object |

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
| 公共属性 | [Type](33ca3ea7-1f17-6eb5-1e22-38eb57f85d4b.htm) | 获取或设置当前的机器人类型，默认为 [YRC1000](38ab618c-dbb6-4896-0126-5f10e97180ba.htm) Get or set the current robot type, the default is [YRC1000](38ab618c-dbb6-4896-0126-5f10e97180ba.htm) |
| 受保护的属性 | [UseServerActivePush](0e77fcc5-1744-cedf-a8e1-958031151826.htm) | 获取或设置当前的连接是否激活从服务器主动推送的功能 (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 受保护的方法 | [AccountCertificate](e9ef8320-5daf-4629-dcc6-4076e7944a9a.htm) | 认证账号，根据已经设置的用户名和密码，进行发送服务器进行账号认证。  Authentication account, according to the user name and password that have been set, sending server for account authentication. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 受保护的方法 | [AccountCertificateAsync](2dc6eb5a-79c8-1e51-7aa8-c3d29769033a.htm) | 认证账号，根据已经设置的用户名和密码，进行发送服务器进行账号认证。  Authentication account, according to the user name and password that have been set, sending server for account authentication. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | [Cancel](be66b6de-fe15-95b5-6969-70518f086388.htm) | 进行错误取消  Make an error cancellation |
| 公共方法 | [CancelAsync](fedb71cd-919c-c002-7556-551ae30fb10e.htm) | 进行错误取消  Make an error cancellation |
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
| 公共方法 | [Cycle](5a144fc7-5bc5-7250-6c13-826c172bd2ce.htm) | 选择循环。循环编号 1:步骤，2:1循环，3:连续自动  Choose loop. Cycle number 1: step, 2:1 cycle, 3: continuous automatic |
| 公共方法 | [CycleAsync](82bfe8a9-00b3-b34e-2a82-6d022becdf9a.htm) | 选择循环。循环编号 1:步骤，2:1循环，3:连续自动  Choose loop. Cycle number 1: step, 2:1 cycle, 3: continuous automatic |
| 受保护的方法 | [DecideWhetherQAMessage](77091dac-82d7-8c67-d691-31ae5b2aaee4.htm) | 决定当前的消息是否是应答机制的消息内容，需要在客户端进行重写实现，如果是应答机制，返回 True, 否则返回 False  To determine whether the current message is the message content of the response mechanism, it needs to be rewritten on the client side. If it is the response mechanism, return True, otherwise return False (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | [Delete](c7bfb022-46f4-8cdb-1b50-ffc1e5eb716b.htm) | 删除指定的程序。指定「\*」 时， 删除当前登录的所有程序。指定「 删除程序名称」 时，仅删除指定的程序。  Delete the specified program. When "\*" is specified, all currently registered programs will be deleted. When "delete program name" is specified, only the specified program will be deleted. |
| 公共方法 | [DeleteAsync](c5a62bdd-3432-abcd-9d03-1ff7e30bd447.htm) | 删除指定的程序。指定「\*」 时， 删除当前登录的所有程序。指定「 删除程序名称」 时，仅删除指定的程序。  Delete the specified program. When "\*" is specified, all currently registered programs will be deleted. When "delete program name" is specified, only the specified program will be deleted. |
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
| 受保护的方法 | [GetNewNetMessage](888815ce-da5c-0cb3-c6f9-800ec37a19ab.htm) | 获取一个新的消息对象的方法，需要在继承类里面进行重写  The method to get a new message object needs to be overridden in the inheritance class (重写 [NetworkDoubleBaseGetNewNetMessage](481a40e0-8539-5fc8-3e2b-c25bcc952697.htm).) |
| 公共方法 | [GetPipeSocket](19904680-372e-e0a6-1b27-b4180e759668.htm) | 获取当前用于通信的管道信息  Get the current pipe information used for communication (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | GetType | (继承自 Object。) |
| 公共方法 | [HLock](9c4ed34f-89fd-62e2-f707-08d1348c7541.htm) | 设定示教编程器和 I/O的操作信号的联锁。 状态参数 False: OFF，True: ON  Set the interlock between the programming pendant and the operation signal of I/O. Status parameter False: OFF，True: ON |
| 公共方法 | [HLockAsync](f0eb6982-a92b-4481-6435-33fa2b841b3c.htm) | 设定示教编程器和 I/O的操作信号的联锁。 状态参数 False: OFF，True: ON  Set the interlock between the programming pendant and the operation signal of I/O. Status parameter False: OFF，True: ON |
| 公共方法 | [Hold](bf4ccc0a-1576-83ca-49e3-fde86b3a3738.htm) | 进行HOLD 的 ON/OFF 操作，状态参数 False: OFF操作，True: ON操作  Perform HOLD ON operation, False: OFF，True: ON |
| 公共方法 | [HoldAsync](8c84a7d8-3677-0b6b-bdb5-44aa61ed1fcd.htm) | 进行HOLD 的 ON/OFF 操作，状态参数 False: OFF操作，True: ON操作  Perform HOLD ON operation, False: OFF，True: ON |
| 受保护的方法 | [InitializationOnConnect](5bd0274d-b3e7-f797-2dc0-c0dd2f330694.htm) | before read data , the connection should be Initialized (重写 [NetworkDoubleBaseInitializationOnConnect(Socket)](40d5eb26-9bca-11cf-42a1-55af6ea5aebd.htm).) |
| 受保护的方法 | [InitializationOnConnectAsync](e2f434c6-2ad9-30e8-8f2c-dcab95f5ec2b.htm) | before read data , the connection should be Initialized (重写 [NetworkDoubleBaseInitializationOnConnectAsync(Socket)](14cac21b-4001-da7d-09a6-941eba868e14.htm).) |
| 公共方法 | [IORead](4539f381-2df3-4315-084c-7fcc81d523cb.htm) | 读取I/O 信号。 I/O 数据是每8个点输出，所以读出接点数是8的倍数。  Read I/O signal. I/O data is output every 8 points, so the number of read contacts is a multiple of 8. |
| 公共方法 | [IOReadAsync](5d550dc5-c02c-7245-122e-c3624e8d9c14.htm) | 读取I/O 信号。 I/O 数据是每8个点输出，所以读出接点数是8的倍数。  Read I/O signal. I/O data is output every 8 points, so the number of read contacts is a multiple of 8. |
| 公共方法 | [IOWrite](0dfecba5-f497-4b0f-55a5-339c7202d742.htm) | 写入I/O信号状态，写入接点数请指定8的倍数。IO 信号的网络写入仅可是（ #27010 ～ #29567）。  To write I/O signal status, please specify a multiple of 8 for the number of write contacts. The network write of IO signal is only available (#27010 to #29567). |
| 公共方法 | [IOWriteAsync](64b5c202-29ff-0737-9c62-a4735bc00911.htm) | 写入I/O信号状态，写入接点数请指定8的倍数。IO 信号的网络写入仅可是（ #27010 ～ #29567）。  To write I/O signal status, please specify a multiple of 8 for the number of write contacts. The network write of IO signal is only available (#27010 to #29567). |
| 公共方法 | [IpAddressPing](d9611e18-4033-de95-444e-b91ae0a81709.htm) | 对当前设备的IP地址进行PING的操作，返回PING的结果，正常来说，返回Success  PING the IP address of the current device and return the PING result. Normally, it returns Success (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | [JSeq](7057228b-aa47-3318-248d-1594588783bd.htm) | 设定执行程序的名称和行编号。  Set the name and line number of the executed program. |
| 公共方法 | [JSeqAsync](c4eb7232-63e1-cd10-f89e-7f58072e3e9c.htm) | 设定执行程序的名称和行编号。  Set the name and line number of the executed program. |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法 | [Mode](cabcebca-60ed-7689-e44a-e227936dcc51.htm) | 选择模式。模式编号为1:示教模式，2:再现模式  Choose a mode. The mode number is 1: teaching mode, 2: reproduction mode |
| 公共方法 | [ModeAsync](e0956a5d-5ae1-f884-52dd-0e91b85828b5.htm) | 选择模式。模式编号为1:示教模式，2:再现模式  Choose a mode. The mode number is 1: teaching mode, 2: reproduction mode |
| 公共方法 | [MoveJ](10990eab-b2f9-e4b8-3c88-5bffc78aa34a.htm) | 向指定的坐标系位置进行关节动作。其中没有外部轴的系统， 7-12外部轴的值设定为「0」  Perform joint motions to the specified coordinate system position. where there is no external axis system, the value of 7-12 external axis is set to "0" |
| 公共方法 | [MoveJAsync](5b05904a-27a0-b4f9-adc7-1fe2cc5b14a1.htm) | 向指定的坐标系位置进行关节动作。其中没有外部轴的系统， 7-12外部轴的值设定为「0」  Perform joint motions to the specified coordinate system position. where there is no external axis system, the value of 7-12 external axis is set to "0" |
| 公共方法 | [MSDP](1e78f8b5-7336-7fdf-a7e0-e40af1d3de7f.htm) | 接受消息数据时， 在YRC1000的示教编程器的远程画面下显示消息若。若不是远程画面时，强制切换到远程画面。显示MDSP命令的消息。  When receiving message data, a message is displayed on the remote screen of the YRC1000 programming pendant. If it is not a remote screen, it is forced to switch to the remote screen. Display the message of the MDSP command. |
| 公共方法 | [MSDPAsync](c2a18a60-0dae-112b-af16-0e028e6a203a.htm) | 接受消息数据时， 在YRC1000的示教编程器的远程画面下显示消息若。若不是远程画面时，强制切换到远程画面。显示MDSP命令的消息。  When receiving message data, a message is displayed on the remote screen of the YRC1000 programming pendant. If it is not a remote screen, it is forced to switch to the remote screen. Display the message of the MDSP command. |
| 公共方法 | [PackCommandWithHeader](f96881d6-e59f-9e16-07e3-1188eee1a934.htm) | 对当前的命令进行打包处理，通常是携带命令头内容，标记当前的命令的长度信息，需要进行重写，否则默认不打包  The current command is packaged, usually carrying the content of the command header, marking the length of the current command, and it needs to be rewritten, otherwise it is not packaged by default (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | [Read](69082cc9-c365-8fc4-36df-082c0c6c22c3.htm) | 根据地址读取机器人的原始的字节数据信息  Read the robot's original byte data information according to the address |
| 公共方法 | [ReadALARM](b56d10b7-3ebc-effd-c8f4-1dbf2d9f35a0.htm) | 读取机器人的报警信息  Read the alarm information of the robot |
| 公共方法 | [ReadALARMAsync](6bd50e3f-6923-2c34-0e8d-d08ba82840f1.htm) | 读取机器人的报警信息  Read the alarm information of the robot |
| 公共方法 | [ReadAsync](185f22c2-a4f1-985e-dd8e-4bf7107c5ae7.htm) | 根据地址读取机器人的原始的字节数据信息  Read the robot's original byte data information according to the address |
| 公共方法 | [ReadByCommand](80dbd712-4782-b205-7e73-977c37f5a214.htm) | 根据指令来读取设备的信息，如果命令数据为空，则传入null即可，注意，所有的命令不带换行符  Read the device information according to the instructions. If the command data is empty, pass in null. Note that all commands do not have a newline character |
| 公共方法 | [ReadByCommandAsync](2ac7651a-e36e-1635-0a61-b855722693a1.htm) | 根据指令来读取设备的信息，如果命令数据为空，则传入null即可，注意，所有的命令不带换行符  Read the device information according to the instructions. If the command data is empty, pass in null. Note that all commands do not have a newline character |
| 公共方法 | [ReadByteVariable](e22f2610-10f3-8472-92c2-3f9e904ebb2b.htm) | 读取机器人的字节型变量的数据，需要传入变量的编号  To read the data of the byte variable of the robot, the number of the variable needs to be passed in |
| 公共方法 | [ReadByteVariableAsync](d01b8a15-6079-9c57-6936-e0c569cf5194.htm) | 读取机器人的字节型变量的数据，需要传入变量的编号  To read the data of the byte variable of the robot, the number of the variable needs to be passed in |
| 公共方法 | [ReadDoubleIntegerVariable](beab7b65-77db-3c62-6f95-ccf796ecac03.htm) | 读取机器人的双精度整型变量的数据，需要传入变量的编号  To read the data of the double integer variable of the robot, the number of the variable needs to be passed in |
| 公共方法 | [ReadDoubleIntegerVariableAsync](c6845ca5-1d25-d613-89f9-28da2515cf24.htm) | 读取机器人的双精度整型变量的数据，需要传入变量的编号  To read the data of the double integer variable of the robot, the number of the variable needs to be passed in |
| 公共方法代码示例 | [ReadFromCoreServer(Byte)](b3e14431-d773-0d1a-f453-a3ad60648c29.htm) | 将数据发送到当前的网络通道中，并从网络通道中接收一个[INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm)指定的完整的报文，网络通道将根据[GetAvailableSocket](e86a6fc1-eb91-8de9-5ff8-799b0ea3786b.htm)方法自动获取，本方法是线程安全的。  Send data to the current network channel and receive a complete message specified by [INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm) from the network channel. The network channel will be automatically obtained according to the [GetAvailableSocket](e86a6fc1-eb91-8de9-5ff8-799b0ea3786b.htm) method This method is thread-safe. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | [ReadFromCoreServer(IEnumerableByte)](580dae22-fad5-5816-d58f-43068e89e419.htm) | 将多个数据报文按顺序发到设备，并从设备接收返回的数据内容，然后拼接成一个Byte[]信息，需要重写[UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm)方法才能返回正确的结果。  Send multiple data packets to the device in sequence, and receive the returned data content from the device, and then splicing them into a Byte[] message, you need to rewrite [UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm) method to return the correct result. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 受保护的方法 | [ReadFromCoreServer(Socket, String)](6f6ef8c0-ab75-9e60-470d-b1deb46ea8c4.htm) | Read string value from socket |
| 公共方法代码示例 | [ReadFromCoreServer(Byte, Boolean, Boolean)](1a33d974-805d-f149-3665-eb8d8a99ea84.htm) | 将数据发送到当前的网络通道中，并从网络通道中接收一个[INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm)指定的完整的报文，网络通道将根据[GetAvailableSocket](e86a6fc1-eb91-8de9-5ff8-799b0ea3786b.htm)方法自动获取，本方法是线程安全的。  Send data to the current network channel and receive a complete message specified by [INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm) from the network channel. The network channel will be automatically obtained according to the [GetAvailableSocket](e86a6fc1-eb91-8de9-5ff8-799b0ea3786b.htm) method This method is thread-safe. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServer(Socket, Byte, Boolean, Boolean)](3b72b289-b201-7548-0981-ac049c9a4135.htm) | 将数据报文发送指定的网络通道上，根据当前指定的[INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm)类型，返回一条完整的数据指令  Sends a data message to the specified network channel, and returns a complete data command according to the currently specified [INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm) type (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServerAsync(Byte)](ec51a18e-0f5d-dbe1-c229-d62e8760ac97.htm) | 将数据发送到当前的网络通道中，并从网络通道中接收一个[INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm)指定的完整的报文，网络通道将根据[GetAvailableSocket](e86a6fc1-eb91-8de9-5ff8-799b0ea3786b.htm)方法自动获取，本方法是线程安全的。  Send data to the current network channel and receive a complete message specified by [INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm) from the network channel. The network channel will be automatically obtained according to the [GetAvailableSocket](e86a6fc1-eb91-8de9-5ff8-799b0ea3786b.htm) method This method is thread-safe. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | [ReadFromCoreServerAsync(IEnumerableByte)](0fce97fd-8441-6b54-0877-b8dfacc738e3.htm) | 将多个数据报文按顺序发到设备，并从设备接收返回的数据内容，然后拼接成一个Byte[]信息，需要重写[UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm)方法才能返回正确的结果。  Send multiple data packets to the device in sequence, and receive the returned data content from the device, and then splicing them into a Byte[] message, you need to rewrite [UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm) method to return the correct result. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 受保护的方法 | [ReadFromCoreServerAsync(Socket, String)](064572e5-938f-bc9a-a9d0-b2744e21cbe5.htm) | Read string value from socket |
| 公共方法代码示例 | [ReadFromCoreServerAsync(Byte, Boolean, Boolean)](7525b379-04c2-8dac-b451-a8b5413d6747.htm) | 将数据发送到当前的网络通道中，并从网络通道中接收一个[INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm)指定的完整的报文，网络通道将根据[GetAvailableSocket](e86a6fc1-eb91-8de9-5ff8-799b0ea3786b.htm)方法自动获取，本方法是线程安全的。  Send data to the current network channel and receive a complete message specified by [INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm) from the network channel. The network channel will be automatically obtained according to the [GetAvailableSocket](e86a6fc1-eb91-8de9-5ff8-799b0ea3786b.htm) method This method is thread-safe. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServerAsync(Socket, Byte, Boolean, Boolean)](07db4a00-b548-8bd6-5de6-12deb91f5311.htm) | 将数据报文发送指定的网络通道上，根据当前指定的[INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm)类型，返回一条完整的数据指令  Sends a data message to the specified network channel, and returns a complete data command according to the currently specified [INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm) type (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | [ReadIntegerVariable](df3171c4-e18b-a52c-9032-bd789e7388a0.htm) | 读取机器人的整型变量的数据，需要传入变量的编号  To read the data of the integer variable of the robot, the number of the variable needs to be passed in |
| 公共方法 | [ReadIntegerVariableAsync](954c86a5-37e1-f32b-3a6e-7fc4274cbb6f.htm) | 读取机器人的整型变量的数据，需要传入变量的编号  To read the data of the integer variable of the robot, the number of the variable needs to be passed in |
| 公共方法 | [ReadJSeq](d7a72f10-2923-e045-8c71-e8b43c551c5a.htm) | 读取当前的程序名，行编号，步编号。  Read the current program name, line number, and step number. |
| 公共方法 | [ReadJSeqAsync](ddfc3321-88c7-9d2a-b71e-9315ce77864d.htm) | 读取当前的程序名，行编号，步编号。  Read the current program name, line number, and step number. |
| 公共方法 | [ReadPOSC](fea57b6b-1d2f-0547-73de-6120ee5de705.htm) | 指定坐标系的当前值读取。并且可以指定外部轴的有无。  The current value of the specified coordinate system is read. And you can specify the presence or absence of an external axis. |
| 公共方法 | [ReadPOSCAsync](5f5731c4-a35a-affa-3d76-7b5f50fe63a1.htm) | 指定坐标系的当前值读取。并且可以指定外部轴的有无。  The current value of the specified coordinate system is read. And you can specify the presence or absence of an external axis. |
| 公共方法 | [ReadPOSJ](14045444-1d70-3bdf-9848-a886c204da46.htm) | 关节坐标系的坐标位置读取。  Read the coordinate data information of the robot |
| 公共方法 | [ReadPOSJAsync](0ed6f125-5c82-6a30-1a4e-a6a773dda5d8.htm) | 关节坐标系的坐标位置读取。  Read the coordinate data information of the robot |
| 公共方法 | [ReadRealVariable](5f6c460a-0bb9-32fb-8aa8-33b99dc8be5e.htm) | 读取机器人的实数变量的数据，需要传入变量的编号  To read the data of the real variable of the robot, the number of the variable needs to be passed in |
| 公共方法 | [ReadRealVariableAsync](996ae64e-84f9-ccf6-fa2a-de4e4a1b13c9.htm) | 读取机器人的实数变量的数据，需要传入变量的编号  To read the data of the real variable of the robot, the number of the variable needs to be passed in |
| 公共方法 | [ReadStats](3a24171c-f204-8c63-cdbc-8cf23a965b75.htm) | 模式状态，循环状态，动作状态，报警错误状态，伺服状态的读取。  Reading of mode status, cycle status, action status, alarm error status, and servo status. |
| 公共方法 | [ReadStatsAsync](43a75b5c-c5af-e81f-286c-a4aa2a7a898b.htm) | 模式状态，循环状态，动作状态，报警错误状态，伺服状态的读取。  Reading of mode status, cycle status, action status, alarm error status, and servo status. |
| 公共方法 | [ReadString](dbc09266-42e8-b939-aea5-540d528d19ee.htm) | 根据地址读取机器人的字符串的数据信息  Read the string data information of the robot based on the address |
| 公共方法 | [ReadStringAsync](8fa97613-eb8d-c1e1-33f6-925dcf1a1ac0.htm) | 根据地址读取机器人的字符串的数据信息  Read the string data information of the robot based on the address |
| 公共方法 | [ReadStringVariable](5dc9095d-86e8-d434-0027-83797a1d8c0a.htm) | 读取机器人的字符串变量的数据，需要传入变量的编号  To read the data of the string variable of the robot, the number of the variable needs to be passed in |
| 公共方法 | [ReadStringVariableAsync](f850e762-bf8a-6ee7-c18c-3496d70244fe.htm) | 读取机器人的字符串变量的数据，需要传入变量的编号  To read the data of the string variable of the robot, the number of the variable needs to be passed in |
| 公共方法 | [ReadUFrame](100a1955-fc0a-871a-98f9-bb361c422a3f.htm) | 读取指定用户的坐标数据。  Read the coordinate data of the specified user. |
| 公共方法 | [ReadUFrameAsync](f63d72f5-eaf4-ae5d-e80a-d881f93a878e.htm) | 读取指定用户的坐标数据。  Read the coordinate data of the specified user. |
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
| 公共方法 | [Reset](bf61707b-875c-4b42-2732-9b2fb69535ac.htm) | 对机械手的报警进行复位  Reset the alarm of the manipulator |
| 公共方法 | [ResetAsync](f135af7f-0250-eb16-48fe-959dd0aeeb72.htm) | 对机械手的报警进行复位  Reset the alarm of the manipulator |
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
| 公共方法 | [SetMJ](9cd8c255-f908-6f4e-da07-33edc743dbf8.htm) | 指定的程序设定为主程序。设定主程序的同时执行程序也被设定。  The specified program is set as the main program. The execution program is also set when the main program is set. |
| 公共方法 | [SetMJAsync](74b1cb5d-83f7-61b7-8eb0-0e10a675b964.htm) | 指定的程序设定为主程序。设定主程序的同时执行程序也被设定。  The specified program is set as the main program. The execution program is also set when the main program is set. |
| 公共方法代码示例 | [SetPersistentConnection](d8e8df9f-cbe9-6863-54e9-aa3b38906b26.htm) | 在读取数据之前可以调用本方法将客户端设置为长连接模式，相当于跳过了ConnectServer的结果验证，对异形客户端无效，当第一次进行通信时再进行创建连接请求。  Before reading the data, you can call this method to set the client to the long connection mode, which is equivalent to skipping the result verification of ConnectServer, and it is invalid for the alien client. When the first communication is performed, the connection creation request is performed. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | [SetPipeSocket](b6a7a243-2a6d-c4ea-45b7-bf980450a2bd.htm) | 设置一个新的网络管道，一般来说不需要调用本方法，当多个网口设备共用一个网络连接时才需要使用本方法进行设置共享的管道。  To set up a new network channel, generally speaking, you do not need to call this method. This method is only needed to set up a shared channel when multiple network port devices share a network connection. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | [Start](39764e6e-fde9-c99b-9034-99a64232aef6.htm) | 开始程序。操作时指定程序名时，此程序能附带对应主程序，则从该程序的开头开始执行。如果没有指定，则从前行开始执行  Start the program. When the program name is specified during operation, the program can be accompanied by the corresponding main program, and the execution starts from the beginning of the program. If not specified, execute from the previous line |
| 公共方法 | [StartAsync](ddefed38-f8b8-4ae1-73e3-eecbfae54043.htm) | 开始程序。操作时指定程序名时，此程序能附带对应主程序，则从该程序的开头开始执行。如果没有指定，则从前行开始执行  Start the program. When the program name is specified during operation, the program can be accompanied by the corresponding main program, and the execution starts from the beginning of the program. If not specified, execute from the previous line |
| 公共方法 | [Svon](1f83c8cb-fea3-d99f-c7ad-f12f39bd9dee.htm) | 进行伺服电源的ON/OFF操作，状态参数 False: OFF，True: ON  Carry out the ON/OFF operation of the servo power, the status parameter False: OFF，True: ON |
| 公共方法 | [SvonAsync](9b51697e-6164-280a-d37c-1cc16f467b2d.htm) | 进行伺服电源的ON/OFF操作，状态参数 False: OFF，True: ON  Carry out the ON/OFF operation of the servo power, the status parameter False: OFF，True: ON |
| 公共方法 | [ToString](1a5b666e-cae2-1664-908f-56866da693e5.htm) | (重写 [NetworkDoubleBaseToString](85a8532f-8fc5-886c-9d76-70844e0a7d64.htm).) |
| 公共方法 | [UnpackResponseContent](32bad9a0-b59d-c171-ce85-74055edfc597.htm) | 根据对方返回的报文命令，对命令进行基本的拆包，例如各种Modbus协议拆包为统一的核心报文，还支持对报文的验证  According to the message command returned by the other party, the command is basically unpacked, for example, various Modbus protocols are unpacked into a unified core message, and the verification of the message is also supported (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | [Write(String, Byte)](e940746c-c324-0cdd-d2cb-0be92dec38c4.htm) | 根据地址，来写入设备的相关的字节数据  According to the address, to write the device related bytes data |
| 公共方法 | [Write(String, String)](cad64999-c93f-b5fe-5d49-6da85c0e08a6.htm) | 根据地址，来写入设备相关的字符串数据  According to the address, to write the device related string data |
| 公共方法 | [WriteAsync(String, Byte)](a1a50677-47c1-eb7d-96c2-2930f998ba8e.htm) | 根据地址，来写入设备的相关的字节数据  According to the address, to write the device related bytes data |
| 公共方法 | [WriteAsync(String, String)](ef5e384d-acac-ee75-14a2-369eed504d6b.htm) | 根据地址，来写入设备相关的字符串数据  According to the address, to write the device related string data |
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

要想成功的通信，有两个至关重要的前提。  
1. 开启以太网服务器，[系统]-[设定]-[选项功能]-[网络功能设定]启用网络功能。  
2. 开启远程的命令，[输入输出]-[模拟输入]-[远程命令选择] 激活远程命令

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.Robot.YASKAWA 命名空间](3d67023f-78aa-c5d9-85b6-f18f5dfbc176.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## YRC1000TcpNet 构造函数 

[原文連結](http://api.hslcommunication.cn/html/ebe3130d-b046-c1e9-5c6e-de2669433a40.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.YASKAWA](../html/3d67023f-78aa-c5d9-85b6-f18f5dfbc176.htm "HslCommunication.Robot.YASKAWA")

[YRC1000TcpNet 类](../html/cdcdc1cc-17cf-da27-7d20-ca6d6db1cd7e.htm "YRC1000TcpNet 类")

[YRC1000TcpNet 构造函数](../html/ebe3130d-b046-c1e9-5c6e-de2669433a40.htm "YRC1000TcpNet 构造函数 ")

[YRC1000TcpNet 属性](../html/bae7bc3e-546f-f9ee-b28d-7ad7c6eea531.htm "YRC1000TcpNet 属性")

[YRC1000TcpNet 方法](../html/7a26c354-eee0-b601-1d8a-51198e40e3fa.htm "YRC1000TcpNet 方法")

[YRC1000TcpNet 字段](../html/296507d2-3c83-4a11-77d9-21def100a3eb.htm "YRC1000TcpNet 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| YRC1000TcpNet 构造函数 |

指定机器人的ip地址及端口号来实例化对象  
Specify the robot's IP address and port number to instantiate the object

**命名空间：**
 [HslCommunication.Robot.YASKAWA](3d67023f-78aa-c5d9-85b6-f18f5dfbc176.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public YRC1000TcpNet(
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
YRC1000TcpNet(
	String^ ipAddress, 
	int port
)
```

```
new : 
        ipAddress : string * 
        port : int -> YRC1000TcpNet
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

[YRC1000TcpNet 类](cdcdc1cc-17cf-da27-7d20-ca6d6db1cd7e.htm)

[HslCommunication.Robot.YASKAWA 命名空间](3d67023f-78aa-c5d9-85b6-f18f5dfbc176.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## YRC1000TcpNet 属性

[原文連結](http://api.hslcommunication.cn/html/bae7bc3e-546f-f9ee-b28d-7ad7c6eea531.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.YASKAWA](../html/3d67023f-78aa-c5d9-85b6-f18f5dfbc176.htm "HslCommunication.Robot.YASKAWA")

[YRC1000TcpNet 类](../html/cdcdc1cc-17cf-da27-7d20-ca6d6db1cd7e.htm "YRC1000TcpNet 类")

[YRC1000TcpNet 属性](../html/bae7bc3e-546f-f9ee-b28d-7ad7c6eea531.htm "YRC1000TcpNet 属性")

[Type 属性](../html/33ca3ea7-1f17-6eb5-1e22-38eb57f85d4b.htm "Type 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| YRC1000TcpNet 属性 |

[YRC1000TcpNet](cdcdc1cc-17cf-da27-7d20-ca6d6db1cd7e.htm) 类型公开以下成员。

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
| 公共属性 | [Type](33ca3ea7-1f17-6eb5-1e22-38eb57f85d4b.htm) | 获取或设置当前的机器人类型，默认为 [YRC1000](38ab618c-dbb6-4896-0126-5f10e97180ba.htm) Get or set the current robot type, the default is [YRC1000](38ab618c-dbb6-4896-0126-5f10e97180ba.htm) |
| 受保护的属性 | [UseServerActivePush](0e77fcc5-1744-cedf-a8e1-958031151826.htm) | 获取或设置当前的连接是否激活从服务器主动推送的功能 (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[YRC1000TcpNet 类](cdcdc1cc-17cf-da27-7d20-ca6d6db1cd7e.htm)

[HslCommunication.Robot.YASKAWA 命名空间](3d67023f-78aa-c5d9-85b6-f18f5dfbc176.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Type 属性 

[原文連結](http://api.hslcommunication.cn/html/33ca3ea7-1f17-6eb5-1e22-38eb57f85d4b.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.YASKAWA](../html/3d67023f-78aa-c5d9-85b6-f18f5dfbc176.htm "HslCommunication.Robot.YASKAWA")

[YRC1000TcpNet 类](../html/cdcdc1cc-17cf-da27-7d20-ca6d6db1cd7e.htm "YRC1000TcpNet 类")

[YRC1000TcpNet 属性](../html/bae7bc3e-546f-f9ee-b28d-7ad7c6eea531.htm "YRC1000TcpNet 属性")

[Type 属性](../html/33ca3ea7-1f17-6eb5-1e22-38eb57f85d4b.htm "Type 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| YRC1000TcpNetType 属性 |

获取或设置当前的机器人类型，默认为 [YRC1000](38ab618c-dbb6-4896-0126-5f10e97180ba.htm)
Get or set the current robot type, the default is [YRC1000](38ab618c-dbb6-4896-0126-5f10e97180ba.htm)

**命名空间：**
 [HslCommunication.Robot.YASKAWA](3d67023f-78aa-c5d9-85b6-f18f5dfbc176.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public YRCType Type { get; set; }
```

```
Public Property Type As YRCType
	Get
	Set
```

```
public:
property YRCType Type {
	YRCType get ();
	void set (YRCType value);
}
```

```
member Type : YRCType with get, set
```

#### 属性值

类型：[YRCType](38ab618c-dbb6-4896-0126-5f10e97180ba.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[YRC1000TcpNet 类](cdcdc1cc-17cf-da27-7d20-ca6d6db1cd7e.htm)

[HslCommunication.Robot.YASKAWA 命名空间](3d67023f-78aa-c5d9-85b6-f18f5dfbc176.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## YRC1000TcpNet 方法

[原文連結](http://api.hslcommunication.cn/html/7a26c354-eee0-b601-1d8a-51198e40e3fa.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.YASKAWA](../html/3d67023f-78aa-c5d9-85b6-f18f5dfbc176.htm "HslCommunication.Robot.YASKAWA")

[YRC1000TcpNet 类](../html/cdcdc1cc-17cf-da27-7d20-ca6d6db1cd7e.htm "YRC1000TcpNet 类")

[YRC1000TcpNet 方法](../html/7a26c354-eee0-b601-1d8a-51198e40e3fa.htm "YRC1000TcpNet 方法")

[Cancel 方法](../html/be66b6de-fe15-95b5-6969-70518f086388.htm "Cancel 方法 ")

[CancelAsync 方法](../html/fedb71cd-919c-c002-7556-551ae30fb10e.htm "CancelAsync 方法 ")

[Cycle 方法](../html/5a144fc7-5bc5-7250-6c13-826c172bd2ce.htm "Cycle 方法 ")

[CycleAsync 方法](../html/82bfe8a9-00b3-b34e-2a82-6d022becdf9a.htm "CycleAsync 方法 ")

[Delete 方法](../html/c7bfb022-46f4-8cdb-1b50-ffc1e5eb716b.htm "Delete 方法 ")

[DeleteAsync 方法](../html/c5a62bdd-3432-abcd-9d03-1ff7e30bd447.htm "DeleteAsync 方法 ")

[GetNewNetMessage 方法](../html/888815ce-da5c-0cb3-c6f9-800ec37a19ab.htm "GetNewNetMessage 方法 ")

[HLock 方法](../html/9c4ed34f-89fd-62e2-f707-08d1348c7541.htm "HLock 方法 ")

[HLockAsync 方法](../html/f0eb6982-a92b-4481-6435-33fa2b841b3c.htm "HLockAsync 方法 ")

[Hold 方法](../html/bf4ccc0a-1576-83ca-49e3-fde86b3a3738.htm "Hold 方法 ")

[HoldAsync 方法](../html/8c84a7d8-3677-0b6b-bdb5-44aa61ed1fcd.htm "HoldAsync 方法 ")

[InitializationOnConnect 方法](../html/5bd0274d-b3e7-f797-2dc0-c0dd2f330694.htm "InitializationOnConnect 方法 ")

[InitializationOnConnectAsync 方法](../html/e2f434c6-2ad9-30e8-8f2c-dcab95f5ec2b.htm "InitializationOnConnectAsync 方法 ")

[IORead 方法](../html/4539f381-2df3-4315-084c-7fcc81d523cb.htm "IORead 方法 ")

[IOReadAsync 方法](../html/5d550dc5-c02c-7245-122e-c3624e8d9c14.htm "IOReadAsync 方法 ")

[IOWrite 方法](../html/0dfecba5-f497-4b0f-55a5-339c7202d742.htm "IOWrite 方法 ")

[IOWriteAsync 方法](../html/64b5c202-29ff-0737-9c62-a4735bc00911.htm "IOWriteAsync 方法 ")

[JSeq 方法](../html/7057228b-aa47-3318-248d-1594588783bd.htm "JSeq 方法 ")

[JSeqAsync 方法](../html/c4eb7232-63e1-cd10-f89e-7f58072e3e9c.htm "JSeqAsync 方法 ")

[Mode 方法](../html/cabcebca-60ed-7689-e44a-e227936dcc51.htm "Mode 方法 ")

[ModeAsync 方法](../html/e0956a5d-5ae1-f884-52dd-0e91b85828b5.htm "ModeAsync 方法 ")

[MoveJ 方法](../html/10990eab-b2f9-e4b8-3c88-5bffc78aa34a.htm "MoveJ 方法 ")

[MoveJAsync 方法](../html/5b05904a-27a0-b4f9-adc7-1fe2cc5b14a1.htm "MoveJAsync 方法 ")

[MSDP 方法](../html/1e78f8b5-7336-7fdf-a7e0-e40af1d3de7f.htm "MSDP 方法 ")

[MSDPAsync 方法](../html/c2a18a60-0dae-112b-af16-0e028e6a203a.htm "MSDPAsync 方法 ")

[Read 方法](../html/69082cc9-c365-8fc4-36df-082c0c6c22c3.htm "Read 方法 ")

[ReadALARM 方法](../html/b56d10b7-3ebc-effd-c8f4-1dbf2d9f35a0.htm "ReadALARM 方法 ")

[ReadALARMAsync 方法](../html/6bd50e3f-6923-2c34-0e8d-d08ba82840f1.htm "ReadALARMAsync 方法 ")

[ReadAsync 方法](../html/185f22c2-a4f1-985e-dd8e-4bf7107c5ae7.htm "ReadAsync 方法 ")

[ReadByCommand 方法](../html/80dbd712-4782-b205-7e73-977c37f5a214.htm "ReadByCommand 方法 ")

[ReadByCommandAsync 方法](../html/2ac7651a-e36e-1635-0a61-b855722693a1.htm "ReadByCommandAsync 方法 ")

[ReadByteVariable 方法](../html/e22f2610-10f3-8472-92c2-3f9e904ebb2b.htm "ReadByteVariable 方法 ")

[ReadByteVariableAsync 方法](../html/d01b8a15-6079-9c57-6936-e0c569cf5194.htm "ReadByteVariableAsync 方法 ")

[ReadDoubleIntegerVariable 方法](../html/beab7b65-77db-3c62-6f95-ccf796ecac03.htm "ReadDoubleIntegerVariable 方法 ")

[ReadDoubleIntegerVariableAsync 方法](../html/c6845ca5-1d25-d613-89f9-28da2515cf24.htm "ReadDoubleIntegerVariableAsync 方法 ")

[ReadFromCoreServer 方法](../html/4e4be36a-8ce8-a41b-6ec4-4f82479c6f23.htm "ReadFromCoreServer 方法 ")

[ReadFromCoreServerAsync 方法](../html/e4cb7ebe-d42a-ea17-8972-b246a42e9a90.htm "ReadFromCoreServerAsync 方法 ")

[ReadIntegerVariable 方法](../html/df3171c4-e18b-a52c-9032-bd789e7388a0.htm "ReadIntegerVariable 方法 ")

[ReadIntegerVariableAsync 方法](../html/954c86a5-37e1-f32b-3a6e-7fc4274cbb6f.htm "ReadIntegerVariableAsync 方法 ")

[ReadJSeq 方法](../html/d7a72f10-2923-e045-8c71-e8b43c551c5a.htm "ReadJSeq 方法 ")

[ReadJSeqAsync 方法](../html/ddfc3321-88c7-9d2a-b71e-9315ce77864d.htm "ReadJSeqAsync 方法 ")

[ReadPOSC 方法](../html/fea57b6b-1d2f-0547-73de-6120ee5de705.htm "ReadPOSC 方法 ")

[ReadPOSCAsync 方法](../html/5f5731c4-a35a-affa-3d76-7b5f50fe63a1.htm "ReadPOSCAsync 方法 ")

[ReadPOSJ 方法](../html/14045444-1d70-3bdf-9848-a886c204da46.htm "ReadPOSJ 方法 ")

[ReadPOSJAsync 方法](../html/0ed6f125-5c82-6a30-1a4e-a6a773dda5d8.htm "ReadPOSJAsync 方法 ")

[ReadRealVariable 方法](../html/5f6c460a-0bb9-32fb-8aa8-33b99dc8be5e.htm "ReadRealVariable 方法 ")

[ReadRealVariableAsync 方法](../html/996ae64e-84f9-ccf6-fa2a-de4e4a1b13c9.htm "ReadRealVariableAsync 方法 ")

[ReadStats 方法](../html/3a24171c-f204-8c63-cdbc-8cf23a965b75.htm "ReadStats 方法 ")

[ReadStatsAsync 方法](../html/43a75b5c-c5af-e81f-286c-a4aa2a7a898b.htm "ReadStatsAsync 方法 ")

[ReadString 方法](../html/dbc09266-42e8-b939-aea5-540d528d19ee.htm "ReadString 方法 ")

[ReadStringAsync 方法](../html/8fa97613-eb8d-c1e1-33f6-925dcf1a1ac0.htm "ReadStringAsync 方法 ")

[ReadStringVariable 方法](../html/5dc9095d-86e8-d434-0027-83797a1d8c0a.htm "ReadStringVariable 方法 ")

[ReadStringVariableAsync 方法](../html/f850e762-bf8a-6ee7-c18c-3496d70244fe.htm "ReadStringVariableAsync 方法 ")

[ReadUFrame 方法](../html/100a1955-fc0a-871a-98f9-bb361c422a3f.htm "ReadUFrame 方法 ")

[ReadUFrameAsync 方法](../html/f63d72f5-eaf4-ae5d-e80a-d881f93a878e.htm "ReadUFrameAsync 方法 ")

[Reset 方法](../html/bf61707b-875c-4b42-2732-9b2fb69535ac.htm "Reset 方法 ")

[ResetAsync 方法](../html/f135af7f-0250-eb16-48fe-959dd0aeeb72.htm "ResetAsync 方法 ")

[SetMJ 方法](../html/9cd8c255-f908-6f4e-da07-33edc743dbf8.htm "SetMJ 方法 ")

[SetMJAsync 方法](../html/74b1cb5d-83f7-61b7-8eb0-0e10a675b964.htm "SetMJAsync 方法 ")

[Start 方法](../html/39764e6e-fde9-c99b-9034-99a64232aef6.htm "Start 方法 ")

[StartAsync 方法](../html/ddefed38-f8b8-4ae1-73e3-eecbfae54043.htm "StartAsync 方法 ")

[Svon 方法](../html/1f83c8cb-fea3-d99f-c7ad-f12f39bd9dee.htm "Svon 方法 ")

[SvonAsync 方法](../html/9b51697e-6164-280a-d37c-1cc16f467b2d.htm "SvonAsync 方法 ")

[ToString 方法](../html/1a5b666e-cae2-1664-908f-56866da693e5.htm "ToString 方法 ")

[Write 方法](../html/27720779-dc92-a6b4-c792-940ec2a10d22.htm "Write 方法 ")

[WriteAsync 方法](../html/e1646488-a935-0a5b-ec6a-4d04a3e69d38.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| YRC1000TcpNet 方法 |

[YRC1000TcpNet](cdcdc1cc-17cf-da27-7d20-ca6d6db1cd7e.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 受保护的方法 | [AccountCertificate](e9ef8320-5daf-4629-dcc6-4076e7944a9a.htm) | 认证账号，根据已经设置的用户名和密码，进行发送服务器进行账号认证。  Authentication account, according to the user name and password that have been set, sending server for account authentication. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 受保护的方法 | [AccountCertificateAsync](2dc6eb5a-79c8-1e51-7aa8-c3d29769033a.htm) | 认证账号，根据已经设置的用户名和密码，进行发送服务器进行账号认证。  Authentication account, according to the user name and password that have been set, sending server for account authentication. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | [Cancel](be66b6de-fe15-95b5-6969-70518f086388.htm) | 进行错误取消  Make an error cancellation |
| 公共方法 | [CancelAsync](fedb71cd-919c-c002-7556-551ae30fb10e.htm) | 进行错误取消  Make an error cancellation |
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
| 公共方法 | [Cycle](5a144fc7-5bc5-7250-6c13-826c172bd2ce.htm) | 选择循环。循环编号 1:步骤，2:1循环，3:连续自动  Choose loop. Cycle number 1: step, 2:1 cycle, 3: continuous automatic |
| 公共方法 | [CycleAsync](82bfe8a9-00b3-b34e-2a82-6d022becdf9a.htm) | 选择循环。循环编号 1:步骤，2:1循环，3:连续自动  Choose loop. Cycle number 1: step, 2:1 cycle, 3: continuous automatic |
| 受保护的方法 | [DecideWhetherQAMessage](77091dac-82d7-8c67-d691-31ae5b2aaee4.htm) | 决定当前的消息是否是应答机制的消息内容，需要在客户端进行重写实现，如果是应答机制，返回 True, 否则返回 False  To determine whether the current message is the message content of the response mechanism, it needs to be rewritten on the client side. If it is the response mechanism, return True, otherwise return False (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | [Delete](c7bfb022-46f4-8cdb-1b50-ffc1e5eb716b.htm) | 删除指定的程序。指定「\*」 时， 删除当前登录的所有程序。指定「 删除程序名称」 时，仅删除指定的程序。  Delete the specified program. When "\*" is specified, all currently registered programs will be deleted. When "delete program name" is specified, only the specified program will be deleted. |
| 公共方法 | [DeleteAsync](c5a62bdd-3432-abcd-9d03-1ff7e30bd447.htm) | 删除指定的程序。指定「\*」 时， 删除当前登录的所有程序。指定「 删除程序名称」 时，仅删除指定的程序。  Delete the specified program. When "\*" is specified, all currently registered programs will be deleted. When "delete program name" is specified, only the specified program will be deleted. |
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
| 受保护的方法 | [GetNewNetMessage](888815ce-da5c-0cb3-c6f9-800ec37a19ab.htm) | 获取一个新的消息对象的方法，需要在继承类里面进行重写  The method to get a new message object needs to be overridden in the inheritance class (重写 [NetworkDoubleBaseGetNewNetMessage](481a40e0-8539-5fc8-3e2b-c25bcc952697.htm).) |
| 公共方法 | [GetPipeSocket](19904680-372e-e0a6-1b27-b4180e759668.htm) | 获取当前用于通信的管道信息  Get the current pipe information used for communication (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | GetType | (继承自 Object。) |
| 公共方法 | [HLock](9c4ed34f-89fd-62e2-f707-08d1348c7541.htm) | 设定示教编程器和 I/O的操作信号的联锁。 状态参数 False: OFF，True: ON  Set the interlock between the programming pendant and the operation signal of I/O. Status parameter False: OFF，True: ON |
| 公共方法 | [HLockAsync](f0eb6982-a92b-4481-6435-33fa2b841b3c.htm) | 设定示教编程器和 I/O的操作信号的联锁。 状态参数 False: OFF，True: ON  Set the interlock between the programming pendant and the operation signal of I/O. Status parameter False: OFF，True: ON |
| 公共方法 | [Hold](bf4ccc0a-1576-83ca-49e3-fde86b3a3738.htm) | 进行HOLD 的 ON/OFF 操作，状态参数 False: OFF操作，True: ON操作  Perform HOLD ON operation, False: OFF，True: ON |
| 公共方法 | [HoldAsync](8c84a7d8-3677-0b6b-bdb5-44aa61ed1fcd.htm) | 进行HOLD 的 ON/OFF 操作，状态参数 False: OFF操作，True: ON操作  Perform HOLD ON operation, False: OFF，True: ON |
| 受保护的方法 | [InitializationOnConnect](5bd0274d-b3e7-f797-2dc0-c0dd2f330694.htm) | before read data , the connection should be Initialized (重写 [NetworkDoubleBaseInitializationOnConnect(Socket)](40d5eb26-9bca-11cf-42a1-55af6ea5aebd.htm).) |
| 受保护的方法 | [InitializationOnConnectAsync](e2f434c6-2ad9-30e8-8f2c-dcab95f5ec2b.htm) | before read data , the connection should be Initialized (重写 [NetworkDoubleBaseInitializationOnConnectAsync(Socket)](14cac21b-4001-da7d-09a6-941eba868e14.htm).) |
| 公共方法 | [IORead](4539f381-2df3-4315-084c-7fcc81d523cb.htm) | 读取I/O 信号。 I/O 数据是每8个点输出，所以读出接点数是8的倍数。  Read I/O signal. I/O data is output every 8 points, so the number of read contacts is a multiple of 8. |
| 公共方法 | [IOReadAsync](5d550dc5-c02c-7245-122e-c3624e8d9c14.htm) | 读取I/O 信号。 I/O 数据是每8个点输出，所以读出接点数是8的倍数。  Read I/O signal. I/O data is output every 8 points, so the number of read contacts is a multiple of 8. |
| 公共方法 | [IOWrite](0dfecba5-f497-4b0f-55a5-339c7202d742.htm) | 写入I/O信号状态，写入接点数请指定8的倍数。IO 信号的网络写入仅可是（ #27010 ～ #29567）。  To write I/O signal status, please specify a multiple of 8 for the number of write contacts. The network write of IO signal is only available (#27010 to #29567). |
| 公共方法 | [IOWriteAsync](64b5c202-29ff-0737-9c62-a4735bc00911.htm) | 写入I/O信号状态，写入接点数请指定8的倍数。IO 信号的网络写入仅可是（ #27010 ～ #29567）。  To write I/O signal status, please specify a multiple of 8 for the number of write contacts. The network write of IO signal is only available (#27010 to #29567). |
| 公共方法 | [IpAddressPing](d9611e18-4033-de95-444e-b91ae0a81709.htm) | 对当前设备的IP地址进行PING的操作，返回PING的结果，正常来说，返回Success  PING the IP address of the current device and return the PING result. Normally, it returns Success (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | [JSeq](7057228b-aa47-3318-248d-1594588783bd.htm) | 设定执行程序的名称和行编号。  Set the name and line number of the executed program. |
| 公共方法 | [JSeqAsync](c4eb7232-63e1-cd10-f89e-7f58072e3e9c.htm) | 设定执行程序的名称和行编号。  Set the name and line number of the executed program. |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法 | [Mode](cabcebca-60ed-7689-e44a-e227936dcc51.htm) | 选择模式。模式编号为1:示教模式，2:再现模式  Choose a mode. The mode number is 1: teaching mode, 2: reproduction mode |
| 公共方法 | [ModeAsync](e0956a5d-5ae1-f884-52dd-0e91b85828b5.htm) | 选择模式。模式编号为1:示教模式，2:再现模式  Choose a mode. The mode number is 1: teaching mode, 2: reproduction mode |
| 公共方法 | [MoveJ](10990eab-b2f9-e4b8-3c88-5bffc78aa34a.htm) | 向指定的坐标系位置进行关节动作。其中没有外部轴的系统， 7-12外部轴的值设定为「0」  Perform joint motions to the specified coordinate system position. where there is no external axis system, the value of 7-12 external axis is set to "0" |
| 公共方法 | [MoveJAsync](5b05904a-27a0-b4f9-adc7-1fe2cc5b14a1.htm) | 向指定的坐标系位置进行关节动作。其中没有外部轴的系统， 7-12外部轴的值设定为「0」  Perform joint motions to the specified coordinate system position. where there is no external axis system, the value of 7-12 external axis is set to "0" |
| 公共方法 | [MSDP](1e78f8b5-7336-7fdf-a7e0-e40af1d3de7f.htm) | 接受消息数据时， 在YRC1000的示教编程器的远程画面下显示消息若。若不是远程画面时，强制切换到远程画面。显示MDSP命令的消息。  When receiving message data, a message is displayed on the remote screen of the YRC1000 programming pendant. If it is not a remote screen, it is forced to switch to the remote screen. Display the message of the MDSP command. |
| 公共方法 | [MSDPAsync](c2a18a60-0dae-112b-af16-0e028e6a203a.htm) | 接受消息数据时， 在YRC1000的示教编程器的远程画面下显示消息若。若不是远程画面时，强制切换到远程画面。显示MDSP命令的消息。  When receiving message data, a message is displayed on the remote screen of the YRC1000 programming pendant. If it is not a remote screen, it is forced to switch to the remote screen. Display the message of the MDSP command. |
| 公共方法 | [PackCommandWithHeader](f96881d6-e59f-9e16-07e3-1188eee1a934.htm) | 对当前的命令进行打包处理，通常是携带命令头内容，标记当前的命令的长度信息，需要进行重写，否则默认不打包  The current command is packaged, usually carrying the content of the command header, marking the length of the current command, and it needs to be rewritten, otherwise it is not packaged by default (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | [Read](69082cc9-c365-8fc4-36df-082c0c6c22c3.htm) | 根据地址读取机器人的原始的字节数据信息  Read the robot's original byte data information according to the address |
| 公共方法 | [ReadALARM](b56d10b7-3ebc-effd-c8f4-1dbf2d9f35a0.htm) | 读取机器人的报警信息  Read the alarm information of the robot |
| 公共方法 | [ReadALARMAsync](6bd50e3f-6923-2c34-0e8d-d08ba82840f1.htm) | 读取机器人的报警信息  Read the alarm information of the robot |
| 公共方法 | [ReadAsync](185f22c2-a4f1-985e-dd8e-4bf7107c5ae7.htm) | 根据地址读取机器人的原始的字节数据信息  Read the robot's original byte data information according to the address |
| 公共方法 | [ReadByCommand](80dbd712-4782-b205-7e73-977c37f5a214.htm) | 根据指令来读取设备的信息，如果命令数据为空，则传入null即可，注意，所有的命令不带换行符  Read the device information according to the instructions. If the command data is empty, pass in null. Note that all commands do not have a newline character |
| 公共方法 | [ReadByCommandAsync](2ac7651a-e36e-1635-0a61-b855722693a1.htm) | 根据指令来读取设备的信息，如果命令数据为空，则传入null即可，注意，所有的命令不带换行符  Read the device information according to the instructions. If the command data is empty, pass in null. Note that all commands do not have a newline character |
| 公共方法 | [ReadByteVariable](e22f2610-10f3-8472-92c2-3f9e904ebb2b.htm) | 读取机器人的字节型变量的数据，需要传入变量的编号  To read the data of the byte variable of the robot, the number of the variable needs to be passed in |
| 公共方法 | [ReadByteVariableAsync](d01b8a15-6079-9c57-6936-e0c569cf5194.htm) | 读取机器人的字节型变量的数据，需要传入变量的编号  To read the data of the byte variable of the robot, the number of the variable needs to be passed in |
| 公共方法 | [ReadDoubleIntegerVariable](beab7b65-77db-3c62-6f95-ccf796ecac03.htm) | 读取机器人的双精度整型变量的数据，需要传入变量的编号  To read the data of the double integer variable of the robot, the number of the variable needs to be passed in |
| 公共方法 | [ReadDoubleIntegerVariableAsync](c6845ca5-1d25-d613-89f9-28da2515cf24.htm) | 读取机器人的双精度整型变量的数据，需要传入变量的编号  To read the data of the double integer variable of the robot, the number of the variable needs to be passed in |
| 公共方法代码示例 | [ReadFromCoreServer(Byte)](b3e14431-d773-0d1a-f453-a3ad60648c29.htm) | 将数据发送到当前的网络通道中，并从网络通道中接收一个[INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm)指定的完整的报文，网络通道将根据[GetAvailableSocket](e86a6fc1-eb91-8de9-5ff8-799b0ea3786b.htm)方法自动获取，本方法是线程安全的。  Send data to the current network channel and receive a complete message specified by [INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm) from the network channel. The network channel will be automatically obtained according to the [GetAvailableSocket](e86a6fc1-eb91-8de9-5ff8-799b0ea3786b.htm) method This method is thread-safe. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | [ReadFromCoreServer(IEnumerableByte)](580dae22-fad5-5816-d58f-43068e89e419.htm) | 将多个数据报文按顺序发到设备，并从设备接收返回的数据内容，然后拼接成一个Byte[]信息，需要重写[UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm)方法才能返回正确的结果。  Send multiple data packets to the device in sequence, and receive the returned data content from the device, and then splicing them into a Byte[] message, you need to rewrite [UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm) method to return the correct result. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 受保护的方法 | [ReadFromCoreServer(Socket, String)](6f6ef8c0-ab75-9e60-470d-b1deb46ea8c4.htm) | Read string value from socket |
| 公共方法代码示例 | [ReadFromCoreServer(Byte, Boolean, Boolean)](1a33d974-805d-f149-3665-eb8d8a99ea84.htm) | 将数据发送到当前的网络通道中，并从网络通道中接收一个[INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm)指定的完整的报文，网络通道将根据[GetAvailableSocket](e86a6fc1-eb91-8de9-5ff8-799b0ea3786b.htm)方法自动获取，本方法是线程安全的。  Send data to the current network channel and receive a complete message specified by [INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm) from the network channel. The network channel will be automatically obtained according to the [GetAvailableSocket](e86a6fc1-eb91-8de9-5ff8-799b0ea3786b.htm) method This method is thread-safe. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServer(Socket, Byte, Boolean, Boolean)](3b72b289-b201-7548-0981-ac049c9a4135.htm) | 将数据报文发送指定的网络通道上，根据当前指定的[INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm)类型，返回一条完整的数据指令  Sends a data message to the specified network channel, and returns a complete data command according to the currently specified [INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm) type (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServerAsync(Byte)](ec51a18e-0f5d-dbe1-c229-d62e8760ac97.htm) | 将数据发送到当前的网络通道中，并从网络通道中接收一个[INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm)指定的完整的报文，网络通道将根据[GetAvailableSocket](e86a6fc1-eb91-8de9-5ff8-799b0ea3786b.htm)方法自动获取，本方法是线程安全的。  Send data to the current network channel and receive a complete message specified by [INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm) from the network channel. The network channel will be automatically obtained according to the [GetAvailableSocket](e86a6fc1-eb91-8de9-5ff8-799b0ea3786b.htm) method This method is thread-safe. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | [ReadFromCoreServerAsync(IEnumerableByte)](0fce97fd-8441-6b54-0877-b8dfacc738e3.htm) | 将多个数据报文按顺序发到设备，并从设备接收返回的数据内容，然后拼接成一个Byte[]信息，需要重写[UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm)方法才能返回正确的结果。  Send multiple data packets to the device in sequence, and receive the returned data content from the device, and then splicing them into a Byte[] message, you need to rewrite [UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm) method to return the correct result. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 受保护的方法 | [ReadFromCoreServerAsync(Socket, String)](064572e5-938f-bc9a-a9d0-b2744e21cbe5.htm) | Read string value from socket |
| 公共方法代码示例 | [ReadFromCoreServerAsync(Byte, Boolean, Boolean)](7525b379-04c2-8dac-b451-a8b5413d6747.htm) | 将数据发送到当前的网络通道中，并从网络通道中接收一个[INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm)指定的完整的报文，网络通道将根据[GetAvailableSocket](e86a6fc1-eb91-8de9-5ff8-799b0ea3786b.htm)方法自动获取，本方法是线程安全的。  Send data to the current network channel and receive a complete message specified by [INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm) from the network channel. The network channel will be automatically obtained according to the [GetAvailableSocket](e86a6fc1-eb91-8de9-5ff8-799b0ea3786b.htm) method This method is thread-safe. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServerAsync(Socket, Byte, Boolean, Boolean)](07db4a00-b548-8bd6-5de6-12deb91f5311.htm) | 将数据报文发送指定的网络通道上，根据当前指定的[INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm)类型，返回一条完整的数据指令  Sends a data message to the specified network channel, and returns a complete data command according to the currently specified [INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm) type (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | [ReadIntegerVariable](df3171c4-e18b-a52c-9032-bd789e7388a0.htm) | 读取机器人的整型变量的数据，需要传入变量的编号  To read the data of the integer variable of the robot, the number of the variable needs to be passed in |
| 公共方法 | [ReadIntegerVariableAsync](954c86a5-37e1-f32b-3a6e-7fc4274cbb6f.htm) | 读取机器人的整型变量的数据，需要传入变量的编号  To read the data of the integer variable of the robot, the number of the variable needs to be passed in |
| 公共方法 | [ReadJSeq](d7a72f10-2923-e045-8c71-e8b43c551c5a.htm) | 读取当前的程序名，行编号，步编号。  Read the current program name, line number, and step number. |
| 公共方法 | [ReadJSeqAsync](ddfc3321-88c7-9d2a-b71e-9315ce77864d.htm) | 读取当前的程序名，行编号，步编号。  Read the current program name, line number, and step number. |
| 公共方法 | [ReadPOSC](fea57b6b-1d2f-0547-73de-6120ee5de705.htm) | 指定坐标系的当前值读取。并且可以指定外部轴的有无。  The current value of the specified coordinate system is read. And you can specify the presence or absence of an external axis. |
| 公共方法 | [ReadPOSCAsync](5f5731c4-a35a-affa-3d76-7b5f50fe63a1.htm) | 指定坐标系的当前值读取。并且可以指定外部轴的有无。  The current value of the specified coordinate system is read. And you can specify the presence or absence of an external axis. |
| 公共方法 | [ReadPOSJ](14045444-1d70-3bdf-9848-a886c204da46.htm) | 关节坐标系的坐标位置读取。  Read the coordinate data information of the robot |
| 公共方法 | [ReadPOSJAsync](0ed6f125-5c82-6a30-1a4e-a6a773dda5d8.htm) | 关节坐标系的坐标位置读取。  Read the coordinate data information of the robot |
| 公共方法 | [ReadRealVariable](5f6c460a-0bb9-32fb-8aa8-33b99dc8be5e.htm) | 读取机器人的实数变量的数据，需要传入变量的编号  To read the data of the real variable of the robot, the number of the variable needs to be passed in |
| 公共方法 | [ReadRealVariableAsync](996ae64e-84f9-ccf6-fa2a-de4e4a1b13c9.htm) | 读取机器人的实数变量的数据，需要传入变量的编号  To read the data of the real variable of the robot, the number of the variable needs to be passed in |
| 公共方法 | [ReadStats](3a24171c-f204-8c63-cdbc-8cf23a965b75.htm) | 模式状态，循环状态，动作状态，报警错误状态，伺服状态的读取。  Reading of mode status, cycle status, action status, alarm error status, and servo status. |
| 公共方法 | [ReadStatsAsync](43a75b5c-c5af-e81f-286c-a4aa2a7a898b.htm) | 模式状态，循环状态，动作状态，报警错误状态，伺服状态的读取。  Reading of mode status, cycle status, action status, alarm error status, and servo status. |
| 公共方法 | [ReadString](dbc09266-42e8-b939-aea5-540d528d19ee.htm) | 根据地址读取机器人的字符串的数据信息  Read the string data information of the robot based on the address |
| 公共方法 | [ReadStringAsync](8fa97613-eb8d-c1e1-33f6-925dcf1a1ac0.htm) | 根据地址读取机器人的字符串的数据信息  Read the string data information of the robot based on the address |
| 公共方法 | [ReadStringVariable](5dc9095d-86e8-d434-0027-83797a1d8c0a.htm) | 读取机器人的字符串变量的数据，需要传入变量的编号  To read the data of the string variable of the robot, the number of the variable needs to be passed in |
| 公共方法 | [ReadStringVariableAsync](f850e762-bf8a-6ee7-c18c-3496d70244fe.htm) | 读取机器人的字符串变量的数据，需要传入变量的编号  To read the data of the string variable of the robot, the number of the variable needs to be passed in |
| 公共方法 | [ReadUFrame](100a1955-fc0a-871a-98f9-bb361c422a3f.htm) | 读取指定用户的坐标数据。  Read the coordinate data of the specified user. |
| 公共方法 | [ReadUFrameAsync](f63d72f5-eaf4-ae5d-e80a-d881f93a878e.htm) | 读取指定用户的坐标数据。  Read the coordinate data of the specified user. |
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
| 公共方法 | [Reset](bf61707b-875c-4b42-2732-9b2fb69535ac.htm) | 对机械手的报警进行复位  Reset the alarm of the manipulator |
| 公共方法 | [ResetAsync](f135af7f-0250-eb16-48fe-959dd0aeeb72.htm) | 对机械手的报警进行复位  Reset the alarm of the manipulator |
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
| 公共方法 | [SetMJ](9cd8c255-f908-6f4e-da07-33edc743dbf8.htm) | 指定的程序设定为主程序。设定主程序的同时执行程序也被设定。  The specified program is set as the main program. The execution program is also set when the main program is set. |
| 公共方法 | [SetMJAsync](74b1cb5d-83f7-61b7-8eb0-0e10a675b964.htm) | 指定的程序设定为主程序。设定主程序的同时执行程序也被设定。  The specified program is set as the main program. The execution program is also set when the main program is set. |
| 公共方法代码示例 | [SetPersistentConnection](d8e8df9f-cbe9-6863-54e9-aa3b38906b26.htm) | 在读取数据之前可以调用本方法将客户端设置为长连接模式，相当于跳过了ConnectServer的结果验证，对异形客户端无效，当第一次进行通信时再进行创建连接请求。  Before reading the data, you can call this method to set the client to the long connection mode, which is equivalent to skipping the result verification of ConnectServer, and it is invalid for the alien client. When the first communication is performed, the connection creation request is performed. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | [SetPipeSocket](b6a7a243-2a6d-c4ea-45b7-bf980450a2bd.htm) | 设置一个新的网络管道，一般来说不需要调用本方法，当多个网口设备共用一个网络连接时才需要使用本方法进行设置共享的管道。  To set up a new network channel, generally speaking, you do not need to call this method. This method is only needed to set up a shared channel when multiple network port devices share a network connection. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | [Start](39764e6e-fde9-c99b-9034-99a64232aef6.htm) | 开始程序。操作时指定程序名时，此程序能附带对应主程序，则从该程序的开头开始执行。如果没有指定，则从前行开始执行  Start the program. When the program name is specified during operation, the program can be accompanied by the corresponding main program, and the execution starts from the beginning of the program. If not specified, execute from the previous line |
| 公共方法 | [StartAsync](ddefed38-f8b8-4ae1-73e3-eecbfae54043.htm) | 开始程序。操作时指定程序名时，此程序能附带对应主程序，则从该程序的开头开始执行。如果没有指定，则从前行开始执行  Start the program. When the program name is specified during operation, the program can be accompanied by the corresponding main program, and the execution starts from the beginning of the program. If not specified, execute from the previous line |
| 公共方法 | [Svon](1f83c8cb-fea3-d99f-c7ad-f12f39bd9dee.htm) | 进行伺服电源的ON/OFF操作，状态参数 False: OFF，True: ON  Carry out the ON/OFF operation of the servo power, the status parameter False: OFF，True: ON |
| 公共方法 | [SvonAsync](9b51697e-6164-280a-d37c-1cc16f467b2d.htm) | 进行伺服电源的ON/OFF操作，状态参数 False: OFF，True: ON  Carry out the ON/OFF operation of the servo power, the status parameter False: OFF，True: ON |
| 公共方法 | [ToString](1a5b666e-cae2-1664-908f-56866da693e5.htm) | (重写 [NetworkDoubleBaseToString](85a8532f-8fc5-886c-9d76-70844e0a7d64.htm).) |
| 公共方法 | [UnpackResponseContent](32bad9a0-b59d-c171-ce85-74055edfc597.htm) | 根据对方返回的报文命令，对命令进行基本的拆包，例如各种Modbus协议拆包为统一的核心报文，还支持对报文的验证  According to the message command returned by the other party, the command is basically unpacked, for example, various Modbus protocols are unpacked into a unified core message, and the verification of the message is also supported (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | [Write(String, Byte)](e940746c-c324-0cdd-d2cb-0be92dec38c4.htm) | 根据地址，来写入设备的相关的字节数据  According to the address, to write the device related bytes data |
| 公共方法 | [Write(String, String)](cad64999-c93f-b5fe-5d49-6da85c0e08a6.htm) | 根据地址，来写入设备相关的字符串数据  According to the address, to write the device related string data |
| 公共方法 | [WriteAsync(String, Byte)](a1a50677-47c1-eb7d-96c2-2930f998ba8e.htm) | 根据地址，来写入设备的相关的字节数据  According to the address, to write the device related bytes data |
| 公共方法 | [WriteAsync(String, String)](ef5e384d-acac-ee75-14a2-369eed504d6b.htm) | 根据地址，来写入设备相关的字符串数据  According to the address, to write the device related string data |
| 受保护的方法 | [WriteStreamFromSocket](97bf95ea-51fe-8f2e-5da6-1c83766aaa78.htm) | 从套接字中接收所有的数据然后写入到指定的流当中去，需要指定数据的长度，支持按照百分比进行进度报告  Receives all data from the socket and writes it to the specified stream. The length of the data needs to be specified, and progress reporting is supported in percentage. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [WriteStreamFromSocketAsync](b4cf7f6d-1e0c-a266-4c49-683fdb14a922.htm) | 从套接字中接收所有的数据然后写入到指定的流当中去，需要指定数据的长度，支持按照百分比进行进度报告  Receives all data from the socket and writes it to the specified stream. The length of the data needs to be specified, and progress reporting is supported in percentage. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[YRC1000TcpNet 类](cdcdc1cc-17cf-da27-7d20-ca6d6db1cd7e.htm)

[HslCommunication.Robot.YASKAWA 命名空间](3d67023f-78aa-c5d9-85b6-f18f5dfbc176.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Cancel 方法 

[原文連結](http://api.hslcommunication.cn/html/be66b6de-fe15-95b5-6969-70518f086388.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.YASKAWA](../html/3d67023f-78aa-c5d9-85b6-f18f5dfbc176.htm "HslCommunication.Robot.YASKAWA")

[YRC1000TcpNet 类](../html/cdcdc1cc-17cf-da27-7d20-ca6d6db1cd7e.htm "YRC1000TcpNet 类")

[YRC1000TcpNet 方法](../html/7a26c354-eee0-b601-1d8a-51198e40e3fa.htm "YRC1000TcpNet 方法")

[Cancel 方法](../html/be66b6de-fe15-95b5-6969-70518f086388.htm "Cancel 方法 ")

[CancelAsync 方法](../html/fedb71cd-919c-c002-7556-551ae30fb10e.htm "CancelAsync 方法 ")

[Cycle 方法](../html/5a144fc7-5bc5-7250-6c13-826c172bd2ce.htm "Cycle 方法 ")

[CycleAsync 方法](../html/82bfe8a9-00b3-b34e-2a82-6d022becdf9a.htm "CycleAsync 方法 ")

[Delete 方法](../html/c7bfb022-46f4-8cdb-1b50-ffc1e5eb716b.htm "Delete 方法 ")

[DeleteAsync 方法](../html/c5a62bdd-3432-abcd-9d03-1ff7e30bd447.htm "DeleteAsync 方法 ")

[GetNewNetMessage 方法](../html/888815ce-da5c-0cb3-c6f9-800ec37a19ab.htm "GetNewNetMessage 方法 ")

[HLock 方法](../html/9c4ed34f-89fd-62e2-f707-08d1348c7541.htm "HLock 方法 ")

[HLockAsync 方法](../html/f0eb6982-a92b-4481-6435-33fa2b841b3c.htm "HLockAsync 方法 ")

[Hold 方法](../html/bf4ccc0a-1576-83ca-49e3-fde86b3a3738.htm "Hold 方法 ")

[HoldAsync 方法](../html/8c84a7d8-3677-0b6b-bdb5-44aa61ed1fcd.htm "HoldAsync 方法 ")

[InitializationOnConnect 方法](../html/5bd0274d-b3e7-f797-2dc0-c0dd2f330694.htm "InitializationOnConnect 方法 ")

[InitializationOnConnectAsync 方法](../html/e2f434c6-2ad9-30e8-8f2c-dcab95f5ec2b.htm "InitializationOnConnectAsync 方法 ")

[IORead 方法](../html/4539f381-2df3-4315-084c-7fcc81d523cb.htm "IORead 方法 ")

[IOReadAsync 方法](../html/5d550dc5-c02c-7245-122e-c3624e8d9c14.htm "IOReadAsync 方法 ")

[IOWrite 方法](../html/0dfecba5-f497-4b0f-55a5-339c7202d742.htm "IOWrite 方法 ")

[IOWriteAsync 方法](../html/64b5c202-29ff-0737-9c62-a4735bc00911.htm "IOWriteAsync 方法 ")

[JSeq 方法](../html/7057228b-aa47-3318-248d-1594588783bd.htm "JSeq 方法 ")

[JSeqAsync 方法](../html/c4eb7232-63e1-cd10-f89e-7f58072e3e9c.htm "JSeqAsync 方法 ")

[Mode 方法](../html/cabcebca-60ed-7689-e44a-e227936dcc51.htm "Mode 方法 ")

[ModeAsync 方法](../html/e0956a5d-5ae1-f884-52dd-0e91b85828b5.htm "ModeAsync 方法 ")

[MoveJ 方法](../html/10990eab-b2f9-e4b8-3c88-5bffc78aa34a.htm "MoveJ 方法 ")

[MoveJAsync 方法](../html/5b05904a-27a0-b4f9-adc7-1fe2cc5b14a1.htm "MoveJAsync 方法 ")

[MSDP 方法](../html/1e78f8b5-7336-7fdf-a7e0-e40af1d3de7f.htm "MSDP 方法 ")

[MSDPAsync 方法](../html/c2a18a60-0dae-112b-af16-0e028e6a203a.htm "MSDPAsync 方法 ")

[Read 方法](../html/69082cc9-c365-8fc4-36df-082c0c6c22c3.htm "Read 方法 ")

[ReadALARM 方法](../html/b56d10b7-3ebc-effd-c8f4-1dbf2d9f35a0.htm "ReadALARM 方法 ")

[ReadALARMAsync 方法](../html/6bd50e3f-6923-2c34-0e8d-d08ba82840f1.htm "ReadALARMAsync 方法 ")

[ReadAsync 方法](../html/185f22c2-a4f1-985e-dd8e-4bf7107c5ae7.htm "ReadAsync 方法 ")

[ReadByCommand 方法](../html/80dbd712-4782-b205-7e73-977c37f5a214.htm "ReadByCommand 方法 ")

[ReadByCommandAsync 方法](../html/2ac7651a-e36e-1635-0a61-b855722693a1.htm "ReadByCommandAsync 方法 ")

[ReadByteVariable 方法](../html/e22f2610-10f3-8472-92c2-3f9e904ebb2b.htm "ReadByteVariable 方法 ")

[ReadByteVariableAsync 方法](../html/d01b8a15-6079-9c57-6936-e0c569cf5194.htm "ReadByteVariableAsync 方法 ")

[ReadDoubleIntegerVariable 方法](../html/beab7b65-77db-3c62-6f95-ccf796ecac03.htm "ReadDoubleIntegerVariable 方法 ")

[ReadDoubleIntegerVariableAsync 方法](../html/c6845ca5-1d25-d613-89f9-28da2515cf24.htm "ReadDoubleIntegerVariableAsync 方法 ")

[ReadFromCoreServer 方法](../html/4e4be36a-8ce8-a41b-6ec4-4f82479c6f23.htm "ReadFromCoreServer 方法 ")

[ReadFromCoreServerAsync 方法](../html/e4cb7ebe-d42a-ea17-8972-b246a42e9a90.htm "ReadFromCoreServerAsync 方法 ")

[ReadIntegerVariable 方法](../html/df3171c4-e18b-a52c-9032-bd789e7388a0.htm "ReadIntegerVariable 方法 ")

[ReadIntegerVariableAsync 方法](../html/954c86a5-37e1-f32b-3a6e-7fc4274cbb6f.htm "ReadIntegerVariableAsync 方法 ")

[ReadJSeq 方法](../html/d7a72f10-2923-e045-8c71-e8b43c551c5a.htm "ReadJSeq 方法 ")

[ReadJSeqAsync 方法](../html/ddfc3321-88c7-9d2a-b71e-9315ce77864d.htm "ReadJSeqAsync 方法 ")

[ReadPOSC 方法](../html/fea57b6b-1d2f-0547-73de-6120ee5de705.htm "ReadPOSC 方法 ")

[ReadPOSCAsync 方法](../html/5f5731c4-a35a-affa-3d76-7b5f50fe63a1.htm "ReadPOSCAsync 方法 ")

[ReadPOSJ 方法](../html/14045444-1d70-3bdf-9848-a886c204da46.htm "ReadPOSJ 方法 ")

[ReadPOSJAsync 方法](../html/0ed6f125-5c82-6a30-1a4e-a6a773dda5d8.htm "ReadPOSJAsync 方法 ")

[ReadRealVariable 方法](../html/5f6c460a-0bb9-32fb-8aa8-33b99dc8be5e.htm "ReadRealVariable 方法 ")

[ReadRealVariableAsync 方法](../html/996ae64e-84f9-ccf6-fa2a-de4e4a1b13c9.htm "ReadRealVariableAsync 方法 ")

[ReadStats 方法](../html/3a24171c-f204-8c63-cdbc-8cf23a965b75.htm "ReadStats 方法 ")

[ReadStatsAsync 方法](../html/43a75b5c-c5af-e81f-286c-a4aa2a7a898b.htm "ReadStatsAsync 方法 ")

[ReadString 方法](../html/dbc09266-42e8-b939-aea5-540d528d19ee.htm "ReadString 方法 ")

[ReadStringAsync 方法](../html/8fa97613-eb8d-c1e1-33f6-925dcf1a1ac0.htm "ReadStringAsync 方法 ")

[ReadStringVariable 方法](../html/5dc9095d-86e8-d434-0027-83797a1d8c0a.htm "ReadStringVariable 方法 ")

[ReadStringVariableAsync 方法](../html/f850e762-bf8a-6ee7-c18c-3496d70244fe.htm "ReadStringVariableAsync 方法 ")

[ReadUFrame 方法](../html/100a1955-fc0a-871a-98f9-bb361c422a3f.htm "ReadUFrame 方法 ")

[ReadUFrameAsync 方法](../html/f63d72f5-eaf4-ae5d-e80a-d881f93a878e.htm "ReadUFrameAsync 方法 ")

[Reset 方法](../html/bf61707b-875c-4b42-2732-9b2fb69535ac.htm "Reset 方法 ")

[ResetAsync 方法](../html/f135af7f-0250-eb16-48fe-959dd0aeeb72.htm "ResetAsync 方法 ")

[SetMJ 方法](../html/9cd8c255-f908-6f4e-da07-33edc743dbf8.htm "SetMJ 方法 ")

[SetMJAsync 方法](../html/74b1cb5d-83f7-61b7-8eb0-0e10a675b964.htm "SetMJAsync 方法 ")

[Start 方法](../html/39764e6e-fde9-c99b-9034-99a64232aef6.htm "Start 方法 ")

[StartAsync 方法](../html/ddefed38-f8b8-4ae1-73e3-eecbfae54043.htm "StartAsync 方法 ")

[Svon 方法](../html/1f83c8cb-fea3-d99f-c7ad-f12f39bd9dee.htm "Svon 方法 ")

[SvonAsync 方法](../html/9b51697e-6164-280a-d37c-1cc16f467b2d.htm "SvonAsync 方法 ")

[ToString 方法](../html/1a5b666e-cae2-1664-908f-56866da693e5.htm "ToString 方法 ")

[Write 方法](../html/27720779-dc92-a6b4-c792-940ec2a10d22.htm "Write 方法 ")

[WriteAsync 方法](../html/e1646488-a935-0a5b-ec6a-4d04a3e69d38.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| YRC1000TcpNetCancel 方法 |

进行错误取消  
Make an error cancellation

**命名空间：**
 [HslCommunication.Robot.YASKAWA](3d67023f-78aa-c5d9-85b6-f18f5dfbc176.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public OperateResult Cancel()
```

```
Public Function Cancel As OperateResult
```

```
public:
OperateResult^ Cancel()
```

```
member Cancel : unit -> OperateResult 
```

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否取消成功

![](../icons/SectionExpanded.png)参见

#### 引用

[YRC1000TcpNet 类](cdcdc1cc-17cf-da27-7d20-ca6d6db1cd7e.htm)

[HslCommunication.Robot.YASKAWA 命名空间](3d67023f-78aa-c5d9-85b6-f18f5dfbc176.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## CancelAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/fedb71cd-919c-c002-7556-551ae30fb10e.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.YASKAWA](../html/3d67023f-78aa-c5d9-85b6-f18f5dfbc176.htm "HslCommunication.Robot.YASKAWA")

[YRC1000TcpNet 类](../html/cdcdc1cc-17cf-da27-7d20-ca6d6db1cd7e.htm "YRC1000TcpNet 类")

[YRC1000TcpNet 方法](../html/7a26c354-eee0-b601-1d8a-51198e40e3fa.htm "YRC1000TcpNet 方法")

[Cancel 方法](../html/be66b6de-fe15-95b5-6969-70518f086388.htm "Cancel 方法 ")

[CancelAsync 方法](../html/fedb71cd-919c-c002-7556-551ae30fb10e.htm "CancelAsync 方法 ")

[Cycle 方法](../html/5a144fc7-5bc5-7250-6c13-826c172bd2ce.htm "Cycle 方法 ")

[CycleAsync 方法](../html/82bfe8a9-00b3-b34e-2a82-6d022becdf9a.htm "CycleAsync 方法 ")

[Delete 方法](../html/c7bfb022-46f4-8cdb-1b50-ffc1e5eb716b.htm "Delete 方法 ")

[DeleteAsync 方法](../html/c5a62bdd-3432-abcd-9d03-1ff7e30bd447.htm "DeleteAsync 方法 ")

[GetNewNetMessage 方法](../html/888815ce-da5c-0cb3-c6f9-800ec37a19ab.htm "GetNewNetMessage 方法 ")

[HLock 方法](../html/9c4ed34f-89fd-62e2-f707-08d1348c7541.htm "HLock 方法 ")

[HLockAsync 方法](../html/f0eb6982-a92b-4481-6435-33fa2b841b3c.htm "HLockAsync 方法 ")

[Hold 方法](../html/bf4ccc0a-1576-83ca-49e3-fde86b3a3738.htm "Hold 方法 ")

[HoldAsync 方法](../html/8c84a7d8-3677-0b6b-bdb5-44aa61ed1fcd.htm "HoldAsync 方法 ")

[InitializationOnConnect 方法](../html/5bd0274d-b3e7-f797-2dc0-c0dd2f330694.htm "InitializationOnConnect 方法 ")

[InitializationOnConnectAsync 方法](../html/e2f434c6-2ad9-30e8-8f2c-dcab95f5ec2b.htm "InitializationOnConnectAsync 方法 ")

[IORead 方法](../html/4539f381-2df3-4315-084c-7fcc81d523cb.htm "IORead 方法 ")

[IOReadAsync 方法](../html/5d550dc5-c02c-7245-122e-c3624e8d9c14.htm "IOReadAsync 方法 ")

[IOWrite 方法](../html/0dfecba5-f497-4b0f-55a5-339c7202d742.htm "IOWrite 方法 ")

[IOWriteAsync 方法](../html/64b5c202-29ff-0737-9c62-a4735bc00911.htm "IOWriteAsync 方法 ")

[JSeq 方法](../html/7057228b-aa47-3318-248d-1594588783bd.htm "JSeq 方法 ")

[JSeqAsync 方法](../html/c4eb7232-63e1-cd10-f89e-7f58072e3e9c.htm "JSeqAsync 方法 ")

[Mode 方法](../html/cabcebca-60ed-7689-e44a-e227936dcc51.htm "Mode 方法 ")

[ModeAsync 方法](../html/e0956a5d-5ae1-f884-52dd-0e91b85828b5.htm "ModeAsync 方法 ")

[MoveJ 方法](../html/10990eab-b2f9-e4b8-3c88-5bffc78aa34a.htm "MoveJ 方法 ")

[MoveJAsync 方法](../html/5b05904a-27a0-b4f9-adc7-1fe2cc5b14a1.htm "MoveJAsync 方法 ")

[MSDP 方法](../html/1e78f8b5-7336-7fdf-a7e0-e40af1d3de7f.htm "MSDP 方法 ")

[MSDPAsync 方法](../html/c2a18a60-0dae-112b-af16-0e028e6a203a.htm "MSDPAsync 方法 ")

[Read 方法](../html/69082cc9-c365-8fc4-36df-082c0c6c22c3.htm "Read 方法 ")

[ReadALARM 方法](../html/b56d10b7-3ebc-effd-c8f4-1dbf2d9f35a0.htm "ReadALARM 方法 ")

[ReadALARMAsync 方法](../html/6bd50e3f-6923-2c34-0e8d-d08ba82840f1.htm "ReadALARMAsync 方法 ")

[ReadAsync 方法](../html/185f22c2-a4f1-985e-dd8e-4bf7107c5ae7.htm "ReadAsync 方法 ")

[ReadByCommand 方法](../html/80dbd712-4782-b205-7e73-977c37f5a214.htm "ReadByCommand 方法 ")

[ReadByCommandAsync 方法](../html/2ac7651a-e36e-1635-0a61-b855722693a1.htm "ReadByCommandAsync 方法 ")

[ReadByteVariable 方法](../html/e22f2610-10f3-8472-92c2-3f9e904ebb2b.htm "ReadByteVariable 方法 ")

[ReadByteVariableAsync 方法](../html/d01b8a15-6079-9c57-6936-e0c569cf5194.htm "ReadByteVariableAsync 方法 ")

[ReadDoubleIntegerVariable 方法](../html/beab7b65-77db-3c62-6f95-ccf796ecac03.htm "ReadDoubleIntegerVariable 方法 ")

[ReadDoubleIntegerVariableAsync 方法](../html/c6845ca5-1d25-d613-89f9-28da2515cf24.htm "ReadDoubleIntegerVariableAsync 方法 ")

[ReadFromCoreServer 方法](../html/4e4be36a-8ce8-a41b-6ec4-4f82479c6f23.htm "ReadFromCoreServer 方法 ")

[ReadFromCoreServerAsync 方法](../html/e4cb7ebe-d42a-ea17-8972-b246a42e9a90.htm "ReadFromCoreServerAsync 方法 ")

[ReadIntegerVariable 方法](../html/df3171c4-e18b-a52c-9032-bd789e7388a0.htm "ReadIntegerVariable 方法 ")

[ReadIntegerVariableAsync 方法](../html/954c86a5-37e1-f32b-3a6e-7fc4274cbb6f.htm "ReadIntegerVariableAsync 方法 ")

[ReadJSeq 方法](../html/d7a72f10-2923-e045-8c71-e8b43c551c5a.htm "ReadJSeq 方法 ")

[ReadJSeqAsync 方法](../html/ddfc3321-88c7-9d2a-b71e-9315ce77864d.htm "ReadJSeqAsync 方法 ")

[ReadPOSC 方法](../html/fea57b6b-1d2f-0547-73de-6120ee5de705.htm "ReadPOSC 方法 ")

[ReadPOSCAsync 方法](../html/5f5731c4-a35a-affa-3d76-7b5f50fe63a1.htm "ReadPOSCAsync 方法 ")

[ReadPOSJ 方法](../html/14045444-1d70-3bdf-9848-a886c204da46.htm "ReadPOSJ 方法 ")

[ReadPOSJAsync 方法](../html/0ed6f125-5c82-6a30-1a4e-a6a773dda5d8.htm "ReadPOSJAsync 方法 ")

[ReadRealVariable 方法](../html/5f6c460a-0bb9-32fb-8aa8-33b99dc8be5e.htm "ReadRealVariable 方法 ")

[ReadRealVariableAsync 方法](../html/996ae64e-84f9-ccf6-fa2a-de4e4a1b13c9.htm "ReadRealVariableAsync 方法 ")

[ReadStats 方法](../html/3a24171c-f204-8c63-cdbc-8cf23a965b75.htm "ReadStats 方法 ")

[ReadStatsAsync 方法](../html/43a75b5c-c5af-e81f-286c-a4aa2a7a898b.htm "ReadStatsAsync 方法 ")

[ReadString 方法](../html/dbc09266-42e8-b939-aea5-540d528d19ee.htm "ReadString 方法 ")

[ReadStringAsync 方法](../html/8fa97613-eb8d-c1e1-33f6-925dcf1a1ac0.htm "ReadStringAsync 方法 ")

[ReadStringVariable 方法](../html/5dc9095d-86e8-d434-0027-83797a1d8c0a.htm "ReadStringVariable 方法 ")

[ReadStringVariableAsync 方法](../html/f850e762-bf8a-6ee7-c18c-3496d70244fe.htm "ReadStringVariableAsync 方法 ")

[ReadUFrame 方法](../html/100a1955-fc0a-871a-98f9-bb361c422a3f.htm "ReadUFrame 方法 ")

[ReadUFrameAsync 方法](../html/f63d72f5-eaf4-ae5d-e80a-d881f93a878e.htm "ReadUFrameAsync 方法 ")

[Reset 方法](../html/bf61707b-875c-4b42-2732-9b2fb69535ac.htm "Reset 方法 ")

[ResetAsync 方法](../html/f135af7f-0250-eb16-48fe-959dd0aeeb72.htm "ResetAsync 方法 ")

[SetMJ 方法](../html/9cd8c255-f908-6f4e-da07-33edc743dbf8.htm "SetMJ 方法 ")

[SetMJAsync 方法](../html/74b1cb5d-83f7-61b7-8eb0-0e10a675b964.htm "SetMJAsync 方法 ")

[Start 方法](../html/39764e6e-fde9-c99b-9034-99a64232aef6.htm "Start 方法 ")

[StartAsync 方法](../html/ddefed38-f8b8-4ae1-73e3-eecbfae54043.htm "StartAsync 方法 ")

[Svon 方法](../html/1f83c8cb-fea3-d99f-c7ad-f12f39bd9dee.htm "Svon 方法 ")

[SvonAsync 方法](../html/9b51697e-6164-280a-d37c-1cc16f467b2d.htm "SvonAsync 方法 ")

[ToString 方法](../html/1a5b666e-cae2-1664-908f-56866da693e5.htm "ToString 方法 ")

[Write 方法](../html/27720779-dc92-a6b4-c792-940ec2a10d22.htm "Write 方法 ")

[WriteAsync 方法](../html/e1646488-a935-0a5b-ec6a-4d04a3e69d38.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| YRC1000TcpNetCancelAsync 方法 |

进行错误取消  
Make an error cancellation

**命名空间：**
 [HslCommunication.Robot.YASKAWA](3d67023f-78aa-c5d9-85b6-f18f5dfbc176.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public Task<OperateResult> CancelAsync()
```

```
Public Function CancelAsync As Task(Of OperateResult)
```

```
public:
Task<OperateResult^>^ CancelAsync()
```

```
member CancelAsync : unit -> Task<OperateResult> 
```

#### 返回值

类型：Task[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否取消成功

![](../icons/SectionExpanded.png)参见

#### 引用

[YRC1000TcpNet 类](cdcdc1cc-17cf-da27-7d20-ca6d6db1cd7e.htm)

[HslCommunication.Robot.YASKAWA 命名空间](3d67023f-78aa-c5d9-85b6-f18f5dfbc176.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Cycle 方法 

[原文連結](http://api.hslcommunication.cn/html/5a144fc7-5bc5-7250-6c13-826c172bd2ce.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.YASKAWA](../html/3d67023f-78aa-c5d9-85b6-f18f5dfbc176.htm "HslCommunication.Robot.YASKAWA")

[YRC1000TcpNet 类](../html/cdcdc1cc-17cf-da27-7d20-ca6d6db1cd7e.htm "YRC1000TcpNet 类")

[YRC1000TcpNet 方法](../html/7a26c354-eee0-b601-1d8a-51198e40e3fa.htm "YRC1000TcpNet 方法")

[Cancel 方法](../html/be66b6de-fe15-95b5-6969-70518f086388.htm "Cancel 方法 ")

[CancelAsync 方法](../html/fedb71cd-919c-c002-7556-551ae30fb10e.htm "CancelAsync 方法 ")

[Cycle 方法](../html/5a144fc7-5bc5-7250-6c13-826c172bd2ce.htm "Cycle 方法 ")

[CycleAsync 方法](../html/82bfe8a9-00b3-b34e-2a82-6d022becdf9a.htm "CycleAsync 方法 ")

[Delete 方法](../html/c7bfb022-46f4-8cdb-1b50-ffc1e5eb716b.htm "Delete 方法 ")

[DeleteAsync 方法](../html/c5a62bdd-3432-abcd-9d03-1ff7e30bd447.htm "DeleteAsync 方法 ")

[GetNewNetMessage 方法](../html/888815ce-da5c-0cb3-c6f9-800ec37a19ab.htm "GetNewNetMessage 方法 ")

[HLock 方法](../html/9c4ed34f-89fd-62e2-f707-08d1348c7541.htm "HLock 方法 ")

[HLockAsync 方法](../html/f0eb6982-a92b-4481-6435-33fa2b841b3c.htm "HLockAsync 方法 ")

[Hold 方法](../html/bf4ccc0a-1576-83ca-49e3-fde86b3a3738.htm "Hold 方法 ")

[HoldAsync 方法](../html/8c84a7d8-3677-0b6b-bdb5-44aa61ed1fcd.htm "HoldAsync 方法 ")

[InitializationOnConnect 方法](../html/5bd0274d-b3e7-f797-2dc0-c0dd2f330694.htm "InitializationOnConnect 方法 ")

[InitializationOnConnectAsync 方法](../html/e2f434c6-2ad9-30e8-8f2c-dcab95f5ec2b.htm "InitializationOnConnectAsync 方法 ")

[IORead 方法](../html/4539f381-2df3-4315-084c-7fcc81d523cb.htm "IORead 方法 ")

[IOReadAsync 方法](../html/5d550dc5-c02c-7245-122e-c3624e8d9c14.htm "IOReadAsync 方法 ")

[IOWrite 方法](../html/0dfecba5-f497-4b0f-55a5-339c7202d742.htm "IOWrite 方法 ")

[IOWriteAsync 方法](../html/64b5c202-29ff-0737-9c62-a4735bc00911.htm "IOWriteAsync 方法 ")

[JSeq 方法](../html/7057228b-aa47-3318-248d-1594588783bd.htm "JSeq 方法 ")

[JSeqAsync 方法](../html/c4eb7232-63e1-cd10-f89e-7f58072e3e9c.htm "JSeqAsync 方法 ")

[Mode 方法](../html/cabcebca-60ed-7689-e44a-e227936dcc51.htm "Mode 方法 ")

[ModeAsync 方法](../html/e0956a5d-5ae1-f884-52dd-0e91b85828b5.htm "ModeAsync 方法 ")

[MoveJ 方法](../html/10990eab-b2f9-e4b8-3c88-5bffc78aa34a.htm "MoveJ 方法 ")

[MoveJAsync 方法](../html/5b05904a-27a0-b4f9-adc7-1fe2cc5b14a1.htm "MoveJAsync 方法 ")

[MSDP 方法](../html/1e78f8b5-7336-7fdf-a7e0-e40af1d3de7f.htm "MSDP 方法 ")

[MSDPAsync 方法](../html/c2a18a60-0dae-112b-af16-0e028e6a203a.htm "MSDPAsync 方法 ")

[Read 方法](../html/69082cc9-c365-8fc4-36df-082c0c6c22c3.htm "Read 方法 ")

[ReadALARM 方法](../html/b56d10b7-3ebc-effd-c8f4-1dbf2d9f35a0.htm "ReadALARM 方法 ")

[ReadALARMAsync 方法](../html/6bd50e3f-6923-2c34-0e8d-d08ba82840f1.htm "ReadALARMAsync 方法 ")

[ReadAsync 方法](../html/185f22c2-a4f1-985e-dd8e-4bf7107c5ae7.htm "ReadAsync 方法 ")

[ReadByCommand 方法](../html/80dbd712-4782-b205-7e73-977c37f5a214.htm "ReadByCommand 方法 ")

[ReadByCommandAsync 方法](../html/2ac7651a-e36e-1635-0a61-b855722693a1.htm "ReadByCommandAsync 方法 ")

[ReadByteVariable 方法](../html/e22f2610-10f3-8472-92c2-3f9e904ebb2b.htm "ReadByteVariable 方法 ")

[ReadByteVariableAsync 方法](../html/d01b8a15-6079-9c57-6936-e0c569cf5194.htm "ReadByteVariableAsync 方法 ")

[ReadDoubleIntegerVariable 方法](../html/beab7b65-77db-3c62-6f95-ccf796ecac03.htm "ReadDoubleIntegerVariable 方法 ")

[ReadDoubleIntegerVariableAsync 方法](../html/c6845ca5-1d25-d613-89f9-28da2515cf24.htm "ReadDoubleIntegerVariableAsync 方法 ")

[ReadFromCoreServer 方法](../html/4e4be36a-8ce8-a41b-6ec4-4f82479c6f23.htm "ReadFromCoreServer 方法 ")

[ReadFromCoreServerAsync 方法](../html/e4cb7ebe-d42a-ea17-8972-b246a42e9a90.htm "ReadFromCoreServerAsync 方法 ")

[ReadIntegerVariable 方法](../html/df3171c4-e18b-a52c-9032-bd789e7388a0.htm "ReadIntegerVariable 方法 ")

[ReadIntegerVariableAsync 方法](../html/954c86a5-37e1-f32b-3a6e-7fc4274cbb6f.htm "ReadIntegerVariableAsync 方法 ")

[ReadJSeq 方法](../html/d7a72f10-2923-e045-8c71-e8b43c551c5a.htm "ReadJSeq 方法 ")

[ReadJSeqAsync 方法](../html/ddfc3321-88c7-9d2a-b71e-9315ce77864d.htm "ReadJSeqAsync 方法 ")

[ReadPOSC 方法](../html/fea57b6b-1d2f-0547-73de-6120ee5de705.htm "ReadPOSC 方法 ")

[ReadPOSCAsync 方法](../html/5f5731c4-a35a-affa-3d76-7b5f50fe63a1.htm "ReadPOSCAsync 方法 ")

[ReadPOSJ 方法](../html/14045444-1d70-3bdf-9848-a886c204da46.htm "ReadPOSJ 方法 ")

[ReadPOSJAsync 方法](../html/0ed6f125-5c82-6a30-1a4e-a6a773dda5d8.htm "ReadPOSJAsync 方法 ")

[ReadRealVariable 方法](../html/5f6c460a-0bb9-32fb-8aa8-33b99dc8be5e.htm "ReadRealVariable 方法 ")

[ReadRealVariableAsync 方法](../html/996ae64e-84f9-ccf6-fa2a-de4e4a1b13c9.htm "ReadRealVariableAsync 方法 ")

[ReadStats 方法](../html/3a24171c-f204-8c63-cdbc-8cf23a965b75.htm "ReadStats 方法 ")

[ReadStatsAsync 方法](../html/43a75b5c-c5af-e81f-286c-a4aa2a7a898b.htm "ReadStatsAsync 方法 ")

[ReadString 方法](../html/dbc09266-42e8-b939-aea5-540d528d19ee.htm "ReadString 方法 ")

[ReadStringAsync 方法](../html/8fa97613-eb8d-c1e1-33f6-925dcf1a1ac0.htm "ReadStringAsync 方法 ")

[ReadStringVariable 方法](../html/5dc9095d-86e8-d434-0027-83797a1d8c0a.htm "ReadStringVariable 方法 ")

[ReadStringVariableAsync 方法](../html/f850e762-bf8a-6ee7-c18c-3496d70244fe.htm "ReadStringVariableAsync 方法 ")

[ReadUFrame 方法](../html/100a1955-fc0a-871a-98f9-bb361c422a3f.htm "ReadUFrame 方法 ")

[ReadUFrameAsync 方法](../html/f63d72f5-eaf4-ae5d-e80a-d881f93a878e.htm "ReadUFrameAsync 方法 ")

[Reset 方法](../html/bf61707b-875c-4b42-2732-9b2fb69535ac.htm "Reset 方法 ")

[ResetAsync 方法](../html/f135af7f-0250-eb16-48fe-959dd0aeeb72.htm "ResetAsync 方法 ")

[SetMJ 方法](../html/9cd8c255-f908-6f4e-da07-33edc743dbf8.htm "SetMJ 方法 ")

[SetMJAsync 方法](../html/74b1cb5d-83f7-61b7-8eb0-0e10a675b964.htm "SetMJAsync 方法 ")

[Start 方法](../html/39764e6e-fde9-c99b-9034-99a64232aef6.htm "Start 方法 ")

[StartAsync 方法](../html/ddefed38-f8b8-4ae1-73e3-eecbfae54043.htm "StartAsync 方法 ")

[Svon 方法](../html/1f83c8cb-fea3-d99f-c7ad-f12f39bd9dee.htm "Svon 方法 ")

[SvonAsync 方法](../html/9b51697e-6164-280a-d37c-1cc16f467b2d.htm "SvonAsync 方法 ")

[ToString 方法](../html/1a5b666e-cae2-1664-908f-56866da693e5.htm "ToString 方法 ")

[Write 方法](../html/27720779-dc92-a6b4-c792-940ec2a10d22.htm "Write 方法 ")

[WriteAsync 方法](../html/e1646488-a935-0a5b-ec6a-4d04a3e69d38.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| YRC1000TcpNetCycle 方法 |

选择循环。循环编号 1:步骤，2:1循环，3:连续自动  
Choose loop. Cycle number 1: step, 2:1 cycle, 3: continuous automatic

**命名空间：**
 [HslCommunication.Robot.YASKAWA](3d67023f-78aa-c5d9-85b6-f18f5dfbc176.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public OperateResult Cycle(
	int number
)
```

```
Public Function Cycle ( 
	number As Integer
) As OperateResult
```

```
public:
OperateResult^ Cycle(
	int number
)
```

```
member Cycle : 
        number : int -> OperateResult 
```

#### 参数

number
:   类型：SystemInt32  
    循环编号 1:步骤，2:1循环，3:连续自动

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
循环是否选择成功

![](../icons/SectionExpanded.png)参见

#### 引用

[YRC1000TcpNet 类](cdcdc1cc-17cf-da27-7d20-ca6d6db1cd7e.htm)

[HslCommunication.Robot.YASKAWA 命名空间](3d67023f-78aa-c5d9-85b6-f18f5dfbc176.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## CycleAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/82bfe8a9-00b3-b34e-2a82-6d022becdf9a.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.YASKAWA](../html/3d67023f-78aa-c5d9-85b6-f18f5dfbc176.htm "HslCommunication.Robot.YASKAWA")

[YRC1000TcpNet 类](../html/cdcdc1cc-17cf-da27-7d20-ca6d6db1cd7e.htm "YRC1000TcpNet 类")

[YRC1000TcpNet 方法](../html/7a26c354-eee0-b601-1d8a-51198e40e3fa.htm "YRC1000TcpNet 方法")

[Cancel 方法](../html/be66b6de-fe15-95b5-6969-70518f086388.htm "Cancel 方法 ")

[CancelAsync 方法](../html/fedb71cd-919c-c002-7556-551ae30fb10e.htm "CancelAsync 方法 ")

[Cycle 方法](../html/5a144fc7-5bc5-7250-6c13-826c172bd2ce.htm "Cycle 方法 ")

[CycleAsync 方法](../html/82bfe8a9-00b3-b34e-2a82-6d022becdf9a.htm "CycleAsync 方法 ")

[Delete 方法](../html/c7bfb022-46f4-8cdb-1b50-ffc1e5eb716b.htm "Delete 方法 ")

[DeleteAsync 方法](../html/c5a62bdd-3432-abcd-9d03-1ff7e30bd447.htm "DeleteAsync 方法 ")

[GetNewNetMessage 方法](../html/888815ce-da5c-0cb3-c6f9-800ec37a19ab.htm "GetNewNetMessage 方法 ")

[HLock 方法](../html/9c4ed34f-89fd-62e2-f707-08d1348c7541.htm "HLock 方法 ")

[HLockAsync 方法](../html/f0eb6982-a92b-4481-6435-33fa2b841b3c.htm "HLockAsync 方法 ")

[Hold 方法](../html/bf4ccc0a-1576-83ca-49e3-fde86b3a3738.htm "Hold 方法 ")

[HoldAsync 方法](../html/8c84a7d8-3677-0b6b-bdb5-44aa61ed1fcd.htm "HoldAsync 方法 ")

[InitializationOnConnect 方法](../html/5bd0274d-b3e7-f797-2dc0-c0dd2f330694.htm "InitializationOnConnect 方法 ")

[InitializationOnConnectAsync 方法](../html/e2f434c6-2ad9-30e8-8f2c-dcab95f5ec2b.htm "InitializationOnConnectAsync 方法 ")

[IORead 方法](../html/4539f381-2df3-4315-084c-7fcc81d523cb.htm "IORead 方法 ")

[IOReadAsync 方法](../html/5d550dc5-c02c-7245-122e-c3624e8d9c14.htm "IOReadAsync 方法 ")

[IOWrite 方法](../html/0dfecba5-f497-4b0f-55a5-339c7202d742.htm "IOWrite 方法 ")

[IOWriteAsync 方法](../html/64b5c202-29ff-0737-9c62-a4735bc00911.htm "IOWriteAsync 方法 ")

[JSeq 方法](../html/7057228b-aa47-3318-248d-1594588783bd.htm "JSeq 方法 ")

[JSeqAsync 方法](../html/c4eb7232-63e1-cd10-f89e-7f58072e3e9c.htm "JSeqAsync 方法 ")

[Mode 方法](../html/cabcebca-60ed-7689-e44a-e227936dcc51.htm "Mode 方法 ")

[ModeAsync 方法](../html/e0956a5d-5ae1-f884-52dd-0e91b85828b5.htm "ModeAsync 方法 ")

[MoveJ 方法](../html/10990eab-b2f9-e4b8-3c88-5bffc78aa34a.htm "MoveJ 方法 ")

[MoveJAsync 方法](../html/5b05904a-27a0-b4f9-adc7-1fe2cc5b14a1.htm "MoveJAsync 方法 ")

[MSDP 方法](../html/1e78f8b5-7336-7fdf-a7e0-e40af1d3de7f.htm "MSDP 方法 ")

[MSDPAsync 方法](../html/c2a18a60-0dae-112b-af16-0e028e6a203a.htm "MSDPAsync 方法 ")

[Read 方法](../html/69082cc9-c365-8fc4-36df-082c0c6c22c3.htm "Read 方法 ")

[ReadALARM 方法](../html/b56d10b7-3ebc-effd-c8f4-1dbf2d9f35a0.htm "ReadALARM 方法 ")

[ReadALARMAsync 方法](../html/6bd50e3f-6923-2c34-0e8d-d08ba82840f1.htm "ReadALARMAsync 方法 ")

[ReadAsync 方法](../html/185f22c2-a4f1-985e-dd8e-4bf7107c5ae7.htm "ReadAsync 方法 ")

[ReadByCommand 方法](../html/80dbd712-4782-b205-7e73-977c37f5a214.htm "ReadByCommand 方法 ")

[ReadByCommandAsync 方法](../html/2ac7651a-e36e-1635-0a61-b855722693a1.htm "ReadByCommandAsync 方法 ")

[ReadByteVariable 方法](../html/e22f2610-10f3-8472-92c2-3f9e904ebb2b.htm "ReadByteVariable 方法 ")

[ReadByteVariableAsync 方法](../html/d01b8a15-6079-9c57-6936-e0c569cf5194.htm "ReadByteVariableAsync 方法 ")

[ReadDoubleIntegerVariable 方法](../html/beab7b65-77db-3c62-6f95-ccf796ecac03.htm "ReadDoubleIntegerVariable 方法 ")

[ReadDoubleIntegerVariableAsync 方法](../html/c6845ca5-1d25-d613-89f9-28da2515cf24.htm "ReadDoubleIntegerVariableAsync 方法 ")

[ReadFromCoreServer 方法](../html/4e4be36a-8ce8-a41b-6ec4-4f82479c6f23.htm "ReadFromCoreServer 方法 ")

[ReadFromCoreServerAsync 方法](../html/e4cb7ebe-d42a-ea17-8972-b246a42e9a90.htm "ReadFromCoreServerAsync 方法 ")

[ReadIntegerVariable 方法](../html/df3171c4-e18b-a52c-9032-bd789e7388a0.htm "ReadIntegerVariable 方法 ")

[ReadIntegerVariableAsync 方法](../html/954c86a5-37e1-f32b-3a6e-7fc4274cbb6f.htm "ReadIntegerVariableAsync 方法 ")

[ReadJSeq 方法](../html/d7a72f10-2923-e045-8c71-e8b43c551c5a.htm "ReadJSeq 方法 ")

[ReadJSeqAsync 方法](../html/ddfc3321-88c7-9d2a-b71e-9315ce77864d.htm "ReadJSeqAsync 方法 ")

[ReadPOSC 方法](../html/fea57b6b-1d2f-0547-73de-6120ee5de705.htm "ReadPOSC 方法 ")

[ReadPOSCAsync 方法](../html/5f5731c4-a35a-affa-3d76-7b5f50fe63a1.htm "ReadPOSCAsync 方法 ")

[ReadPOSJ 方法](../html/14045444-1d70-3bdf-9848-a886c204da46.htm "ReadPOSJ 方法 ")

[ReadPOSJAsync 方法](../html/0ed6f125-5c82-6a30-1a4e-a6a773dda5d8.htm "ReadPOSJAsync 方法 ")

[ReadRealVariable 方法](../html/5f6c460a-0bb9-32fb-8aa8-33b99dc8be5e.htm "ReadRealVariable 方法 ")

[ReadRealVariableAsync 方法](../html/996ae64e-84f9-ccf6-fa2a-de4e4a1b13c9.htm "ReadRealVariableAsync 方法 ")

[ReadStats 方法](../html/3a24171c-f204-8c63-cdbc-8cf23a965b75.htm "ReadStats 方法 ")

[ReadStatsAsync 方法](../html/43a75b5c-c5af-e81f-286c-a4aa2a7a898b.htm "ReadStatsAsync 方法 ")

[ReadString 方法](../html/dbc09266-42e8-b939-aea5-540d528d19ee.htm "ReadString 方法 ")

[ReadStringAsync 方法](../html/8fa97613-eb8d-c1e1-33f6-925dcf1a1ac0.htm "ReadStringAsync 方法 ")

[ReadStringVariable 方法](../html/5dc9095d-86e8-d434-0027-83797a1d8c0a.htm "ReadStringVariable 方法 ")

[ReadStringVariableAsync 方法](../html/f850e762-bf8a-6ee7-c18c-3496d70244fe.htm "ReadStringVariableAsync 方法 ")

[ReadUFrame 方法](../html/100a1955-fc0a-871a-98f9-bb361c422a3f.htm "ReadUFrame 方法 ")

[ReadUFrameAsync 方法](../html/f63d72f5-eaf4-ae5d-e80a-d881f93a878e.htm "ReadUFrameAsync 方法 ")

[Reset 方法](../html/bf61707b-875c-4b42-2732-9b2fb69535ac.htm "Reset 方法 ")

[ResetAsync 方法](../html/f135af7f-0250-eb16-48fe-959dd0aeeb72.htm "ResetAsync 方法 ")

[SetMJ 方法](../html/9cd8c255-f908-6f4e-da07-33edc743dbf8.htm "SetMJ 方法 ")

[SetMJAsync 方法](../html/74b1cb5d-83f7-61b7-8eb0-0e10a675b964.htm "SetMJAsync 方法 ")

[Start 方法](../html/39764e6e-fde9-c99b-9034-99a64232aef6.htm "Start 方法 ")

[StartAsync 方法](../html/ddefed38-f8b8-4ae1-73e3-eecbfae54043.htm "StartAsync 方法 ")

[Svon 方法](../html/1f83c8cb-fea3-d99f-c7ad-f12f39bd9dee.htm "Svon 方法 ")

[SvonAsync 方法](../html/9b51697e-6164-280a-d37c-1cc16f467b2d.htm "SvonAsync 方法 ")

[ToString 方法](../html/1a5b666e-cae2-1664-908f-56866da693e5.htm "ToString 方法 ")

[Write 方法](../html/27720779-dc92-a6b4-c792-940ec2a10d22.htm "Write 方法 ")

[WriteAsync 方法](../html/e1646488-a935-0a5b-ec6a-4d04a3e69d38.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| YRC1000TcpNetCycleAsync 方法 |

选择循环。循环编号 1:步骤，2:1循环，3:连续自动  
Choose loop. Cycle number 1: step, 2:1 cycle, 3: continuous automatic

**命名空间：**
 [HslCommunication.Robot.YASKAWA](3d67023f-78aa-c5d9-85b6-f18f5dfbc176.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public Task<OperateResult> CycleAsync(
	int number
)
```

```
Public Function CycleAsync ( 
	number As Integer
) As Task(Of OperateResult)
```

```
public:
Task<OperateResult^>^ CycleAsync(
	int number
)
```

```
member CycleAsync : 
        number : int -> Task<OperateResult> 
```

#### 参数

number
:   类型：SystemInt32  
    循环编号 1:步骤，2:1循环，3:连续自动

#### 返回值

类型：Task[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
循环是否选择成功

![](../icons/SectionExpanded.png)参见

#### 引用

[YRC1000TcpNet 类](cdcdc1cc-17cf-da27-7d20-ca6d6db1cd7e.htm)

[HslCommunication.Robot.YASKAWA 命名空间](3d67023f-78aa-c5d9-85b6-f18f5dfbc176.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Delete 方法 

[原文連結](http://api.hslcommunication.cn/html/c7bfb022-46f4-8cdb-1b50-ffc1e5eb716b.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.YASKAWA](../html/3d67023f-78aa-c5d9-85b6-f18f5dfbc176.htm "HslCommunication.Robot.YASKAWA")

[YRC1000TcpNet 类](../html/cdcdc1cc-17cf-da27-7d20-ca6d6db1cd7e.htm "YRC1000TcpNet 类")

[YRC1000TcpNet 方法](../html/7a26c354-eee0-b601-1d8a-51198e40e3fa.htm "YRC1000TcpNet 方法")

[Cancel 方法](../html/be66b6de-fe15-95b5-6969-70518f086388.htm "Cancel 方法 ")

[CancelAsync 方法](../html/fedb71cd-919c-c002-7556-551ae30fb10e.htm "CancelAsync 方法 ")

[Cycle 方法](../html/5a144fc7-5bc5-7250-6c13-826c172bd2ce.htm "Cycle 方法 ")

[CycleAsync 方法](../html/82bfe8a9-00b3-b34e-2a82-6d022becdf9a.htm "CycleAsync 方法 ")

[Delete 方法](../html/c7bfb022-46f4-8cdb-1b50-ffc1e5eb716b.htm "Delete 方法 ")

[DeleteAsync 方法](../html/c5a62bdd-3432-abcd-9d03-1ff7e30bd447.htm "DeleteAsync 方法 ")

[GetNewNetMessage 方法](../html/888815ce-da5c-0cb3-c6f9-800ec37a19ab.htm "GetNewNetMessage 方法 ")

[HLock 方法](../html/9c4ed34f-89fd-62e2-f707-08d1348c7541.htm "HLock 方法 ")

[HLockAsync 方法](../html/f0eb6982-a92b-4481-6435-33fa2b841b3c.htm "HLockAsync 方法 ")

[Hold 方法](../html/bf4ccc0a-1576-83ca-49e3-fde86b3a3738.htm "Hold 方法 ")

[HoldAsync 方法](../html/8c84a7d8-3677-0b6b-bdb5-44aa61ed1fcd.htm "HoldAsync 方法 ")

[InitializationOnConnect 方法](../html/5bd0274d-b3e7-f797-2dc0-c0dd2f330694.htm "InitializationOnConnect 方法 ")

[InitializationOnConnectAsync 方法](../html/e2f434c6-2ad9-30e8-8f2c-dcab95f5ec2b.htm "InitializationOnConnectAsync 方法 ")

[IORead 方法](../html/4539f381-2df3-4315-084c-7fcc81d523cb.htm "IORead 方法 ")

[IOReadAsync 方法](../html/5d550dc5-c02c-7245-122e-c3624e8d9c14.htm "IOReadAsync 方法 ")

[IOWrite 方法](../html/0dfecba5-f497-4b0f-55a5-339c7202d742.htm "IOWrite 方法 ")

[IOWriteAsync 方法](../html/64b5c202-29ff-0737-9c62-a4735bc00911.htm "IOWriteAsync 方法 ")

[JSeq 方法](../html/7057228b-aa47-3318-248d-1594588783bd.htm "JSeq 方法 ")

[JSeqAsync 方法](../html/c4eb7232-63e1-cd10-f89e-7f58072e3e9c.htm "JSeqAsync 方法 ")

[Mode 方法](../html/cabcebca-60ed-7689-e44a-e227936dcc51.htm "Mode 方法 ")

[ModeAsync 方法](../html/e0956a5d-5ae1-f884-52dd-0e91b85828b5.htm "ModeAsync 方法 ")

[MoveJ 方法](../html/10990eab-b2f9-e4b8-3c88-5bffc78aa34a.htm "MoveJ 方法 ")

[MoveJAsync 方法](../html/5b05904a-27a0-b4f9-adc7-1fe2cc5b14a1.htm "MoveJAsync 方法 ")

[MSDP 方法](../html/1e78f8b5-7336-7fdf-a7e0-e40af1d3de7f.htm "MSDP 方法 ")

[MSDPAsync 方法](../html/c2a18a60-0dae-112b-af16-0e028e6a203a.htm "MSDPAsync 方法 ")

[Read 方法](../html/69082cc9-c365-8fc4-36df-082c0c6c22c3.htm "Read 方法 ")

[ReadALARM 方法](../html/b56d10b7-3ebc-effd-c8f4-1dbf2d9f35a0.htm "ReadALARM 方法 ")

[ReadALARMAsync 方法](../html/6bd50e3f-6923-2c34-0e8d-d08ba82840f1.htm "ReadALARMAsync 方法 ")

[ReadAsync 方法](../html/185f22c2-a4f1-985e-dd8e-4bf7107c5ae7.htm "ReadAsync 方法 ")

[ReadByCommand 方法](../html/80dbd712-4782-b205-7e73-977c37f5a214.htm "ReadByCommand 方法 ")

[ReadByCommandAsync 方法](../html/2ac7651a-e36e-1635-0a61-b855722693a1.htm "ReadByCommandAsync 方法 ")

[ReadByteVariable 方法](../html/e22f2610-10f3-8472-92c2-3f9e904ebb2b.htm "ReadByteVariable 方法 ")

[ReadByteVariableAsync 方法](../html/d01b8a15-6079-9c57-6936-e0c569cf5194.htm "ReadByteVariableAsync 方法 ")

[ReadDoubleIntegerVariable 方法](../html/beab7b65-77db-3c62-6f95-ccf796ecac03.htm "ReadDoubleIntegerVariable 方法 ")

[ReadDoubleIntegerVariableAsync 方法](../html/c6845ca5-1d25-d613-89f9-28da2515cf24.htm "ReadDoubleIntegerVariableAsync 方法 ")

[ReadFromCoreServer 方法](../html/4e4be36a-8ce8-a41b-6ec4-4f82479c6f23.htm "ReadFromCoreServer 方法 ")

[ReadFromCoreServerAsync 方法](../html/e4cb7ebe-d42a-ea17-8972-b246a42e9a90.htm "ReadFromCoreServerAsync 方法 ")

[ReadIntegerVariable 方法](../html/df3171c4-e18b-a52c-9032-bd789e7388a0.htm "ReadIntegerVariable 方法 ")

[ReadIntegerVariableAsync 方法](../html/954c86a5-37e1-f32b-3a6e-7fc4274cbb6f.htm "ReadIntegerVariableAsync 方法 ")

[ReadJSeq 方法](../html/d7a72f10-2923-e045-8c71-e8b43c551c5a.htm "ReadJSeq 方法 ")

[ReadJSeqAsync 方法](../html/ddfc3321-88c7-9d2a-b71e-9315ce77864d.htm "ReadJSeqAsync 方法 ")

[ReadPOSC 方法](../html/fea57b6b-1d2f-0547-73de-6120ee5de705.htm "ReadPOSC 方法 ")

[ReadPOSCAsync 方法](../html/5f5731c4-a35a-affa-3d76-7b5f50fe63a1.htm "ReadPOSCAsync 方法 ")

[ReadPOSJ 方法](../html/14045444-1d70-3bdf-9848-a886c204da46.htm "ReadPOSJ 方法 ")

[ReadPOSJAsync 方法](../html/0ed6f125-5c82-6a30-1a4e-a6a773dda5d8.htm "ReadPOSJAsync 方法 ")

[ReadRealVariable 方法](../html/5f6c460a-0bb9-32fb-8aa8-33b99dc8be5e.htm "ReadRealVariable 方法 ")

[ReadRealVariableAsync 方法](../html/996ae64e-84f9-ccf6-fa2a-de4e4a1b13c9.htm "ReadRealVariableAsync 方法 ")

[ReadStats 方法](../html/3a24171c-f204-8c63-cdbc-8cf23a965b75.htm "ReadStats 方法 ")

[ReadStatsAsync 方法](../html/43a75b5c-c5af-e81f-286c-a4aa2a7a898b.htm "ReadStatsAsync 方法 ")

[ReadString 方法](../html/dbc09266-42e8-b939-aea5-540d528d19ee.htm "ReadString 方法 ")

[ReadStringAsync 方法](../html/8fa97613-eb8d-c1e1-33f6-925dcf1a1ac0.htm "ReadStringAsync 方法 ")

[ReadStringVariable 方法](../html/5dc9095d-86e8-d434-0027-83797a1d8c0a.htm "ReadStringVariable 方法 ")

[ReadStringVariableAsync 方法](../html/f850e762-bf8a-6ee7-c18c-3496d70244fe.htm "ReadStringVariableAsync 方法 ")

[ReadUFrame 方法](../html/100a1955-fc0a-871a-98f9-bb361c422a3f.htm "ReadUFrame 方法 ")

[ReadUFrameAsync 方法](../html/f63d72f5-eaf4-ae5d-e80a-d881f93a878e.htm "ReadUFrameAsync 方法 ")

[Reset 方法](../html/bf61707b-875c-4b42-2732-9b2fb69535ac.htm "Reset 方法 ")

[ResetAsync 方法](../html/f135af7f-0250-eb16-48fe-959dd0aeeb72.htm "ResetAsync 方法 ")

[SetMJ 方法](../html/9cd8c255-f908-6f4e-da07-33edc743dbf8.htm "SetMJ 方法 ")

[SetMJAsync 方法](../html/74b1cb5d-83f7-61b7-8eb0-0e10a675b964.htm "SetMJAsync 方法 ")

[Start 方法](../html/39764e6e-fde9-c99b-9034-99a64232aef6.htm "Start 方法 ")

[StartAsync 方法](../html/ddefed38-f8b8-4ae1-73e3-eecbfae54043.htm "StartAsync 方法 ")

[Svon 方法](../html/1f83c8cb-fea3-d99f-c7ad-f12f39bd9dee.htm "Svon 方法 ")

[SvonAsync 方法](../html/9b51697e-6164-280a-d37c-1cc16f467b2d.htm "SvonAsync 方法 ")

[ToString 方法](../html/1a5b666e-cae2-1664-908f-56866da693e5.htm "ToString 方法 ")

[Write 方法](../html/27720779-dc92-a6b4-c792-940ec2a10d22.htm "Write 方法 ")

[WriteAsync 方法](../html/e1646488-a935-0a5b-ec6a-4d04a3e69d38.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| YRC1000TcpNetDelete 方法 |

删除指定的程序。指定「\*」 时， 删除当前登录的所有程序。指定「 删除程序名称」 时，仅删除指定的程序。  
Delete the specified program. When "\*" is specified, all currently registered programs will be deleted.
When "delete program name" is specified, only the specified program will be deleted.

**命名空间：**
 [HslCommunication.Robot.YASKAWA](3d67023f-78aa-c5d9-85b6-f18f5dfbc176.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public OperateResult Delete(
	string programName = null
)
```

```
Public Function Delete ( 
	Optional programName As String = Nothing
) As OperateResult
```

```
public:
OperateResult^ Delete(
	String^ programName = nullptr
)
```

```
member Delete : 
        ?programName : string 
(* Defaults:
        let _programName = defaultArg programName null
*)
-> OperateResult 
```

#### 参数

programName (Optional)
:   类型：SystemString  
    删除的程序名称，如果设置为「\*」时，删除当前登录的所有程序。

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否删除成功

![](../icons/SectionExpanded.png)参见

#### 引用

[YRC1000TcpNet 类](cdcdc1cc-17cf-da27-7d20-ca6d6db1cd7e.htm)

[HslCommunication.Robot.YASKAWA 命名空间](3d67023f-78aa-c5d9-85b6-f18f5dfbc176.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## DeleteAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/c5a62bdd-3432-abcd-9d03-1ff7e30bd447.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.YASKAWA](../html/3d67023f-78aa-c5d9-85b6-f18f5dfbc176.htm "HslCommunication.Robot.YASKAWA")

[YRC1000TcpNet 类](../html/cdcdc1cc-17cf-da27-7d20-ca6d6db1cd7e.htm "YRC1000TcpNet 类")

[YRC1000TcpNet 方法](../html/7a26c354-eee0-b601-1d8a-51198e40e3fa.htm "YRC1000TcpNet 方法")

[Cancel 方法](../html/be66b6de-fe15-95b5-6969-70518f086388.htm "Cancel 方法 ")

[CancelAsync 方法](../html/fedb71cd-919c-c002-7556-551ae30fb10e.htm "CancelAsync 方法 ")

[Cycle 方法](../html/5a144fc7-5bc5-7250-6c13-826c172bd2ce.htm "Cycle 方法 ")

[CycleAsync 方法](../html/82bfe8a9-00b3-b34e-2a82-6d022becdf9a.htm "CycleAsync 方法 ")

[Delete 方法](../html/c7bfb022-46f4-8cdb-1b50-ffc1e5eb716b.htm "Delete 方法 ")

[DeleteAsync 方法](../html/c5a62bdd-3432-abcd-9d03-1ff7e30bd447.htm "DeleteAsync 方法 ")

[GetNewNetMessage 方法](../html/888815ce-da5c-0cb3-c6f9-800ec37a19ab.htm "GetNewNetMessage 方法 ")

[HLock 方法](../html/9c4ed34f-89fd-62e2-f707-08d1348c7541.htm "HLock 方法 ")

[HLockAsync 方法](../html/f0eb6982-a92b-4481-6435-33fa2b841b3c.htm "HLockAsync 方法 ")

[Hold 方法](../html/bf4ccc0a-1576-83ca-49e3-fde86b3a3738.htm "Hold 方法 ")

[HoldAsync 方法](../html/8c84a7d8-3677-0b6b-bdb5-44aa61ed1fcd.htm "HoldAsync 方法 ")

[InitializationOnConnect 方法](../html/5bd0274d-b3e7-f797-2dc0-c0dd2f330694.htm "InitializationOnConnect 方法 ")

[InitializationOnConnectAsync 方法](../html/e2f434c6-2ad9-30e8-8f2c-dcab95f5ec2b.htm "InitializationOnConnectAsync 方法 ")

[IORead 方法](../html/4539f381-2df3-4315-084c-7fcc81d523cb.htm "IORead 方法 ")

[IOReadAsync 方法](../html/5d550dc5-c02c-7245-122e-c3624e8d9c14.htm "IOReadAsync 方法 ")

[IOWrite 方法](../html/0dfecba5-f497-4b0f-55a5-339c7202d742.htm "IOWrite 方法 ")

[IOWriteAsync 方法](../html/64b5c202-29ff-0737-9c62-a4735bc00911.htm "IOWriteAsync 方法 ")

[JSeq 方法](../html/7057228b-aa47-3318-248d-1594588783bd.htm "JSeq 方法 ")

[JSeqAsync 方法](../html/c4eb7232-63e1-cd10-f89e-7f58072e3e9c.htm "JSeqAsync 方法 ")

[Mode 方法](../html/cabcebca-60ed-7689-e44a-e227936dcc51.htm "Mode 方法 ")

[ModeAsync 方法](../html/e0956a5d-5ae1-f884-52dd-0e91b85828b5.htm "ModeAsync 方法 ")

[MoveJ 方法](../html/10990eab-b2f9-e4b8-3c88-5bffc78aa34a.htm "MoveJ 方法 ")

[MoveJAsync 方法](../html/5b05904a-27a0-b4f9-adc7-1fe2cc5b14a1.htm "MoveJAsync 方法 ")

[MSDP 方法](../html/1e78f8b5-7336-7fdf-a7e0-e40af1d3de7f.htm "MSDP 方法 ")

[MSDPAsync 方法](../html/c2a18a60-0dae-112b-af16-0e028e6a203a.htm "MSDPAsync 方法 ")

[Read 方法](../html/69082cc9-c365-8fc4-36df-082c0c6c22c3.htm "Read 方法 ")

[ReadALARM 方法](../html/b56d10b7-3ebc-effd-c8f4-1dbf2d9f35a0.htm "ReadALARM 方法 ")

[ReadALARMAsync 方法](../html/6bd50e3f-6923-2c34-0e8d-d08ba82840f1.htm "ReadALARMAsync 方法 ")

[ReadAsync 方法](../html/185f22c2-a4f1-985e-dd8e-4bf7107c5ae7.htm "ReadAsync 方法 ")

[ReadByCommand 方法](../html/80dbd712-4782-b205-7e73-977c37f5a214.htm "ReadByCommand 方法 ")

[ReadByCommandAsync 方法](../html/2ac7651a-e36e-1635-0a61-b855722693a1.htm "ReadByCommandAsync 方法 ")

[ReadByteVariable 方法](../html/e22f2610-10f3-8472-92c2-3f9e904ebb2b.htm "ReadByteVariable 方法 ")

[ReadByteVariableAsync 方法](../html/d01b8a15-6079-9c57-6936-e0c569cf5194.htm "ReadByteVariableAsync 方法 ")

[ReadDoubleIntegerVariable 方法](../html/beab7b65-77db-3c62-6f95-ccf796ecac03.htm "ReadDoubleIntegerVariable 方法 ")

[ReadDoubleIntegerVariableAsync 方法](../html/c6845ca5-1d25-d613-89f9-28da2515cf24.htm "ReadDoubleIntegerVariableAsync 方法 ")

[ReadFromCoreServer 方法](../html/4e4be36a-8ce8-a41b-6ec4-4f82479c6f23.htm "ReadFromCoreServer 方法 ")

[ReadFromCoreServerAsync 方法](../html/e4cb7ebe-d42a-ea17-8972-b246a42e9a90.htm "ReadFromCoreServerAsync 方法 ")

[ReadIntegerVariable 方法](../html/df3171c4-e18b-a52c-9032-bd789e7388a0.htm "ReadIntegerVariable 方法 ")

[ReadIntegerVariableAsync 方法](../html/954c86a5-37e1-f32b-3a6e-7fc4274cbb6f.htm "ReadIntegerVariableAsync 方法 ")

[ReadJSeq 方法](../html/d7a72f10-2923-e045-8c71-e8b43c551c5a.htm "ReadJSeq 方法 ")

[ReadJSeqAsync 方法](../html/ddfc3321-88c7-9d2a-b71e-9315ce77864d.htm "ReadJSeqAsync 方法 ")

[ReadPOSC 方法](../html/fea57b6b-1d2f-0547-73de-6120ee5de705.htm "ReadPOSC 方法 ")

[ReadPOSCAsync 方法](../html/5f5731c4-a35a-affa-3d76-7b5f50fe63a1.htm "ReadPOSCAsync 方法 ")

[ReadPOSJ 方法](../html/14045444-1d70-3bdf-9848-a886c204da46.htm "ReadPOSJ 方法 ")

[ReadPOSJAsync 方法](../html/0ed6f125-5c82-6a30-1a4e-a6a773dda5d8.htm "ReadPOSJAsync 方法 ")

[ReadRealVariable 方法](../html/5f6c460a-0bb9-32fb-8aa8-33b99dc8be5e.htm "ReadRealVariable 方法 ")

[ReadRealVariableAsync 方法](../html/996ae64e-84f9-ccf6-fa2a-de4e4a1b13c9.htm "ReadRealVariableAsync 方法 ")

[ReadStats 方法](../html/3a24171c-f204-8c63-cdbc-8cf23a965b75.htm "ReadStats 方法 ")

[ReadStatsAsync 方法](../html/43a75b5c-c5af-e81f-286c-a4aa2a7a898b.htm "ReadStatsAsync 方法 ")

[ReadString 方法](../html/dbc09266-42e8-b939-aea5-540d528d19ee.htm "ReadString 方法 ")

[ReadStringAsync 方法](../html/8fa97613-eb8d-c1e1-33f6-925dcf1a1ac0.htm "ReadStringAsync 方法 ")

[ReadStringVariable 方法](../html/5dc9095d-86e8-d434-0027-83797a1d8c0a.htm "ReadStringVariable 方法 ")

[ReadStringVariableAsync 方法](../html/f850e762-bf8a-6ee7-c18c-3496d70244fe.htm "ReadStringVariableAsync 方法 ")

[ReadUFrame 方法](../html/100a1955-fc0a-871a-98f9-bb361c422a3f.htm "ReadUFrame 方法 ")

[ReadUFrameAsync 方法](../html/f63d72f5-eaf4-ae5d-e80a-d881f93a878e.htm "ReadUFrameAsync 方法 ")

[Reset 方法](../html/bf61707b-875c-4b42-2732-9b2fb69535ac.htm "Reset 方法 ")

[ResetAsync 方法](../html/f135af7f-0250-eb16-48fe-959dd0aeeb72.htm "ResetAsync 方法 ")

[SetMJ 方法](../html/9cd8c255-f908-6f4e-da07-33edc743dbf8.htm "SetMJ 方法 ")

[SetMJAsync 方法](../html/74b1cb5d-83f7-61b7-8eb0-0e10a675b964.htm "SetMJAsync 方法 ")

[Start 方法](../html/39764e6e-fde9-c99b-9034-99a64232aef6.htm "Start 方法 ")

[StartAsync 方法](../html/ddefed38-f8b8-4ae1-73e3-eecbfae54043.htm "StartAsync 方法 ")

[Svon 方法](../html/1f83c8cb-fea3-d99f-c7ad-f12f39bd9dee.htm "Svon 方法 ")

[SvonAsync 方法](../html/9b51697e-6164-280a-d37c-1cc16f467b2d.htm "SvonAsync 方法 ")

[ToString 方法](../html/1a5b666e-cae2-1664-908f-56866da693e5.htm "ToString 方法 ")

[Write 方法](../html/27720779-dc92-a6b4-c792-940ec2a10d22.htm "Write 方法 ")

[WriteAsync 方法](../html/e1646488-a935-0a5b-ec6a-4d04a3e69d38.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| YRC1000TcpNetDeleteAsync 方法 |

删除指定的程序。指定「\*」 时， 删除当前登录的所有程序。指定「 删除程序名称」 时，仅删除指定的程序。  
Delete the specified program. When "\*" is specified, all currently registered programs will be deleted.
When "delete program name" is specified, only the specified program will be deleted.

**命名空间：**
 [HslCommunication.Robot.YASKAWA](3d67023f-78aa-c5d9-85b6-f18f5dfbc176.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public Task<OperateResult> DeleteAsync(
	string programName = null
)
```

```
Public Function DeleteAsync ( 
	Optional programName As String = Nothing
) As Task(Of OperateResult)
```

```
public:
Task<OperateResult^>^ DeleteAsync(
	String^ programName = nullptr
)
```

```
member DeleteAsync : 
        ?programName : string 
(* Defaults:
        let _programName = defaultArg programName null
*)
-> Task<OperateResult> 
```

#### 参数

programName (Optional)
:   类型：SystemString  
    删除的程序名称，如果设置为「\*」时，删除当前登录的所有程序。

#### 返回值

类型：Task[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否删除成功

![](../icons/SectionExpanded.png)参见

#### 引用

[YRC1000TcpNet 类](cdcdc1cc-17cf-da27-7d20-ca6d6db1cd7e.htm)

[HslCommunication.Robot.YASKAWA 命名空间](3d67023f-78aa-c5d9-85b6-f18f5dfbc176.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GetNewNetMessage 方法 

[原文連結](http://api.hslcommunication.cn/html/888815ce-da5c-0cb3-c6f9-800ec37a19ab.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.YASKAWA](../html/3d67023f-78aa-c5d9-85b6-f18f5dfbc176.htm "HslCommunication.Robot.YASKAWA")

[YRC1000TcpNet 类](../html/cdcdc1cc-17cf-da27-7d20-ca6d6db1cd7e.htm "YRC1000TcpNet 类")

[YRC1000TcpNet 方法](../html/7a26c354-eee0-b601-1d8a-51198e40e3fa.htm "YRC1000TcpNet 方法")

[Cancel 方法](../html/be66b6de-fe15-95b5-6969-70518f086388.htm "Cancel 方法 ")

[CancelAsync 方法](../html/fedb71cd-919c-c002-7556-551ae30fb10e.htm "CancelAsync 方法 ")

[Cycle 方法](../html/5a144fc7-5bc5-7250-6c13-826c172bd2ce.htm "Cycle 方法 ")

[CycleAsync 方法](../html/82bfe8a9-00b3-b34e-2a82-6d022becdf9a.htm "CycleAsync 方法 ")

[Delete 方法](../html/c7bfb022-46f4-8cdb-1b50-ffc1e5eb716b.htm "Delete 方法 ")

[DeleteAsync 方法](../html/c5a62bdd-3432-abcd-9d03-1ff7e30bd447.htm "DeleteAsync 方法 ")

[GetNewNetMessage 方法](../html/888815ce-da5c-0cb3-c6f9-800ec37a19ab.htm "GetNewNetMessage 方法 ")

[HLock 方法](../html/9c4ed34f-89fd-62e2-f707-08d1348c7541.htm "HLock 方法 ")

[HLockAsync 方法](../html/f0eb6982-a92b-4481-6435-33fa2b841b3c.htm "HLockAsync 方法 ")

[Hold 方法](../html/bf4ccc0a-1576-83ca-49e3-fde86b3a3738.htm "Hold 方法 ")

[HoldAsync 方法](../html/8c84a7d8-3677-0b6b-bdb5-44aa61ed1fcd.htm "HoldAsync 方法 ")

[InitializationOnConnect 方法](../html/5bd0274d-b3e7-f797-2dc0-c0dd2f330694.htm "InitializationOnConnect 方法 ")

[InitializationOnConnectAsync 方法](../html/e2f434c6-2ad9-30e8-8f2c-dcab95f5ec2b.htm "InitializationOnConnectAsync 方法 ")

[IORead 方法](../html/4539f381-2df3-4315-084c-7fcc81d523cb.htm "IORead 方法 ")

[IOReadAsync 方法](../html/5d550dc5-c02c-7245-122e-c3624e8d9c14.htm "IOReadAsync 方法 ")

[IOWrite 方法](../html/0dfecba5-f497-4b0f-55a5-339c7202d742.htm "IOWrite 方法 ")

[IOWriteAsync 方法](../html/64b5c202-29ff-0737-9c62-a4735bc00911.htm "IOWriteAsync 方法 ")

[JSeq 方法](../html/7057228b-aa47-3318-248d-1594588783bd.htm "JSeq 方法 ")

[JSeqAsync 方法](../html/c4eb7232-63e1-cd10-f89e-7f58072e3e9c.htm "JSeqAsync 方法 ")

[Mode 方法](../html/cabcebca-60ed-7689-e44a-e227936dcc51.htm "Mode 方法 ")

[ModeAsync 方法](../html/e0956a5d-5ae1-f884-52dd-0e91b85828b5.htm "ModeAsync 方法 ")

[MoveJ 方法](../html/10990eab-b2f9-e4b8-3c88-5bffc78aa34a.htm "MoveJ 方法 ")

[MoveJAsync 方法](../html/5b05904a-27a0-b4f9-adc7-1fe2cc5b14a1.htm "MoveJAsync 方法 ")

[MSDP 方法](../html/1e78f8b5-7336-7fdf-a7e0-e40af1d3de7f.htm "MSDP 方法 ")

[MSDPAsync 方法](../html/c2a18a60-0dae-112b-af16-0e028e6a203a.htm "MSDPAsync 方法 ")

[Read 方法](../html/69082cc9-c365-8fc4-36df-082c0c6c22c3.htm "Read 方法 ")

[ReadALARM 方法](../html/b56d10b7-3ebc-effd-c8f4-1dbf2d9f35a0.htm "ReadALARM 方法 ")

[ReadALARMAsync 方法](../html/6bd50e3f-6923-2c34-0e8d-d08ba82840f1.htm "ReadALARMAsync 方法 ")

[ReadAsync 方法](../html/185f22c2-a4f1-985e-dd8e-4bf7107c5ae7.htm "ReadAsync 方法 ")

[ReadByCommand 方法](../html/80dbd712-4782-b205-7e73-977c37f5a214.htm "ReadByCommand 方法 ")

[ReadByCommandAsync 方法](../html/2ac7651a-e36e-1635-0a61-b855722693a1.htm "ReadByCommandAsync 方法 ")

[ReadByteVariable 方法](../html/e22f2610-10f3-8472-92c2-3f9e904ebb2b.htm "ReadByteVariable 方法 ")

[ReadByteVariableAsync 方法](../html/d01b8a15-6079-9c57-6936-e0c569cf5194.htm "ReadByteVariableAsync 方法 ")

[ReadDoubleIntegerVariable 方法](../html/beab7b65-77db-3c62-6f95-ccf796ecac03.htm "ReadDoubleIntegerVariable 方法 ")

[ReadDoubleIntegerVariableAsync 方法](../html/c6845ca5-1d25-d613-89f9-28da2515cf24.htm "ReadDoubleIntegerVariableAsync 方法 ")

[ReadFromCoreServer 方法](../html/4e4be36a-8ce8-a41b-6ec4-4f82479c6f23.htm "ReadFromCoreServer 方法 ")

[ReadFromCoreServerAsync 方法](../html/e4cb7ebe-d42a-ea17-8972-b246a42e9a90.htm "ReadFromCoreServerAsync 方法 ")

[ReadIntegerVariable 方法](../html/df3171c4-e18b-a52c-9032-bd789e7388a0.htm "ReadIntegerVariable 方法 ")

[ReadIntegerVariableAsync 方法](../html/954c86a5-37e1-f32b-3a6e-7fc4274cbb6f.htm "ReadIntegerVariableAsync 方法 ")

[ReadJSeq 方法](../html/d7a72f10-2923-e045-8c71-e8b43c551c5a.htm "ReadJSeq 方法 ")

[ReadJSeqAsync 方法](../html/ddfc3321-88c7-9d2a-b71e-9315ce77864d.htm "ReadJSeqAsync 方法 ")

[ReadPOSC 方法](../html/fea57b6b-1d2f-0547-73de-6120ee5de705.htm "ReadPOSC 方法 ")

[ReadPOSCAsync 方法](../html/5f5731c4-a35a-affa-3d76-7b5f50fe63a1.htm "ReadPOSCAsync 方法 ")

[ReadPOSJ 方法](../html/14045444-1d70-3bdf-9848-a886c204da46.htm "ReadPOSJ 方法 ")

[ReadPOSJAsync 方法](../html/0ed6f125-5c82-6a30-1a4e-a6a773dda5d8.htm "ReadPOSJAsync 方法 ")

[ReadRealVariable 方法](../html/5f6c460a-0bb9-32fb-8aa8-33b99dc8be5e.htm "ReadRealVariable 方法 ")

[ReadRealVariableAsync 方法](../html/996ae64e-84f9-ccf6-fa2a-de4e4a1b13c9.htm "ReadRealVariableAsync 方法 ")

[ReadStats 方法](../html/3a24171c-f204-8c63-cdbc-8cf23a965b75.htm "ReadStats 方法 ")

[ReadStatsAsync 方法](../html/43a75b5c-c5af-e81f-286c-a4aa2a7a898b.htm "ReadStatsAsync 方法 ")

[ReadString 方法](../html/dbc09266-42e8-b939-aea5-540d528d19ee.htm "ReadString 方法 ")

[ReadStringAsync 方法](../html/8fa97613-eb8d-c1e1-33f6-925dcf1a1ac0.htm "ReadStringAsync 方法 ")

[ReadStringVariable 方法](../html/5dc9095d-86e8-d434-0027-83797a1d8c0a.htm "ReadStringVariable 方法 ")

[ReadStringVariableAsync 方法](../html/f850e762-bf8a-6ee7-c18c-3496d70244fe.htm "ReadStringVariableAsync 方法 ")

[ReadUFrame 方法](../html/100a1955-fc0a-871a-98f9-bb361c422a3f.htm "ReadUFrame 方法 ")

[ReadUFrameAsync 方法](../html/f63d72f5-eaf4-ae5d-e80a-d881f93a878e.htm "ReadUFrameAsync 方法 ")

[Reset 方法](../html/bf61707b-875c-4b42-2732-9b2fb69535ac.htm "Reset 方法 ")

[ResetAsync 方法](../html/f135af7f-0250-eb16-48fe-959dd0aeeb72.htm "ResetAsync 方法 ")

[SetMJ 方法](../html/9cd8c255-f908-6f4e-da07-33edc743dbf8.htm "SetMJ 方法 ")

[SetMJAsync 方法](../html/74b1cb5d-83f7-61b7-8eb0-0e10a675b964.htm "SetMJAsync 方法 ")

[Start 方法](../html/39764e6e-fde9-c99b-9034-99a64232aef6.htm "Start 方法 ")

[StartAsync 方法](../html/ddefed38-f8b8-4ae1-73e3-eecbfae54043.htm "StartAsync 方法 ")

[Svon 方法](../html/1f83c8cb-fea3-d99f-c7ad-f12f39bd9dee.htm "Svon 方法 ")

[SvonAsync 方法](../html/9b51697e-6164-280a-d37c-1cc16f467b2d.htm "SvonAsync 方法 ")

[ToString 方法](../html/1a5b666e-cae2-1664-908f-56866da693e5.htm "ToString 方法 ")

[Write 方法](../html/27720779-dc92-a6b4-c792-940ec2a10d22.htm "Write 方法 ")

[WriteAsync 方法](../html/e1646488-a935-0a5b-ec6a-4d04a3e69d38.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| YRC1000TcpNetGetNewNetMessage 方法 |

获取一个新的消息对象的方法，需要在继承类里面进行重写  
The method to get a new message object needs to be overridden in the inheritance class

**命名空间：**
 [HslCommunication.Robot.YASKAWA](3d67023f-78aa-c5d9-85b6-f18f5dfbc176.htm)  
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

[YRC1000TcpNet 类](cdcdc1cc-17cf-da27-7d20-ca6d6db1cd7e.htm)

[HslCommunication.Robot.YASKAWA 命名空间](3d67023f-78aa-c5d9-85b6-f18f5dfbc176.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## HLock 方法 

[原文連結](http://api.hslcommunication.cn/html/9c4ed34f-89fd-62e2-f707-08d1348c7541.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.YASKAWA](../html/3d67023f-78aa-c5d9-85b6-f18f5dfbc176.htm "HslCommunication.Robot.YASKAWA")

[YRC1000TcpNet 类](../html/cdcdc1cc-17cf-da27-7d20-ca6d6db1cd7e.htm "YRC1000TcpNet 类")

[YRC1000TcpNet 方法](../html/7a26c354-eee0-b601-1d8a-51198e40e3fa.htm "YRC1000TcpNet 方法")

[Cancel 方法](../html/be66b6de-fe15-95b5-6969-70518f086388.htm "Cancel 方法 ")

[CancelAsync 方法](../html/fedb71cd-919c-c002-7556-551ae30fb10e.htm "CancelAsync 方法 ")

[Cycle 方法](../html/5a144fc7-5bc5-7250-6c13-826c172bd2ce.htm "Cycle 方法 ")

[CycleAsync 方法](../html/82bfe8a9-00b3-b34e-2a82-6d022becdf9a.htm "CycleAsync 方法 ")

[Delete 方法](../html/c7bfb022-46f4-8cdb-1b50-ffc1e5eb716b.htm "Delete 方法 ")

[DeleteAsync 方法](../html/c5a62bdd-3432-abcd-9d03-1ff7e30bd447.htm "DeleteAsync 方法 ")

[GetNewNetMessage 方法](../html/888815ce-da5c-0cb3-c6f9-800ec37a19ab.htm "GetNewNetMessage 方法 ")

[HLock 方法](../html/9c4ed34f-89fd-62e2-f707-08d1348c7541.htm "HLock 方法 ")

[HLockAsync 方法](../html/f0eb6982-a92b-4481-6435-33fa2b841b3c.htm "HLockAsync 方法 ")

[Hold 方法](../html/bf4ccc0a-1576-83ca-49e3-fde86b3a3738.htm "Hold 方法 ")

[HoldAsync 方法](../html/8c84a7d8-3677-0b6b-bdb5-44aa61ed1fcd.htm "HoldAsync 方法 ")

[InitializationOnConnect 方法](../html/5bd0274d-b3e7-f797-2dc0-c0dd2f330694.htm "InitializationOnConnect 方法 ")

[InitializationOnConnectAsync 方法](../html/e2f434c6-2ad9-30e8-8f2c-dcab95f5ec2b.htm "InitializationOnConnectAsync 方法 ")

[IORead 方法](../html/4539f381-2df3-4315-084c-7fcc81d523cb.htm "IORead 方法 ")

[IOReadAsync 方法](../html/5d550dc5-c02c-7245-122e-c3624e8d9c14.htm "IOReadAsync 方法 ")

[IOWrite 方法](../html/0dfecba5-f497-4b0f-55a5-339c7202d742.htm "IOWrite 方法 ")

[IOWriteAsync 方法](../html/64b5c202-29ff-0737-9c62-a4735bc00911.htm "IOWriteAsync 方法 ")

[JSeq 方法](../html/7057228b-aa47-3318-248d-1594588783bd.htm "JSeq 方法 ")

[JSeqAsync 方法](../html/c4eb7232-63e1-cd10-f89e-7f58072e3e9c.htm "JSeqAsync 方法 ")

[Mode 方法](../html/cabcebca-60ed-7689-e44a-e227936dcc51.htm "Mode 方法 ")

[ModeAsync 方法](../html/e0956a5d-5ae1-f884-52dd-0e91b85828b5.htm "ModeAsync 方法 ")

[MoveJ 方法](../html/10990eab-b2f9-e4b8-3c88-5bffc78aa34a.htm "MoveJ 方法 ")

[MoveJAsync 方法](../html/5b05904a-27a0-b4f9-adc7-1fe2cc5b14a1.htm "MoveJAsync 方法 ")

[MSDP 方法](../html/1e78f8b5-7336-7fdf-a7e0-e40af1d3de7f.htm "MSDP 方法 ")

[MSDPAsync 方法](../html/c2a18a60-0dae-112b-af16-0e028e6a203a.htm "MSDPAsync 方法 ")

[Read 方法](../html/69082cc9-c365-8fc4-36df-082c0c6c22c3.htm "Read 方法 ")

[ReadALARM 方法](../html/b56d10b7-3ebc-effd-c8f4-1dbf2d9f35a0.htm "ReadALARM 方法 ")

[ReadALARMAsync 方法](../html/6bd50e3f-6923-2c34-0e8d-d08ba82840f1.htm "ReadALARMAsync 方法 ")

[ReadAsync 方法](../html/185f22c2-a4f1-985e-dd8e-4bf7107c5ae7.htm "ReadAsync 方法 ")

[ReadByCommand 方法](../html/80dbd712-4782-b205-7e73-977c37f5a214.htm "ReadByCommand 方法 ")

[ReadByCommandAsync 方法](../html/2ac7651a-e36e-1635-0a61-b855722693a1.htm "ReadByCommandAsync 方法 ")

[ReadByteVariable 方法](../html/e22f2610-10f3-8472-92c2-3f9e904ebb2b.htm "ReadByteVariable 方法 ")

[ReadByteVariableAsync 方法](../html/d01b8a15-6079-9c57-6936-e0c569cf5194.htm "ReadByteVariableAsync 方法 ")

[ReadDoubleIntegerVariable 方法](../html/beab7b65-77db-3c62-6f95-ccf796ecac03.htm "ReadDoubleIntegerVariable 方法 ")

[ReadDoubleIntegerVariableAsync 方法](../html/c6845ca5-1d25-d613-89f9-28da2515cf24.htm "ReadDoubleIntegerVariableAsync 方法 ")

[ReadFromCoreServer 方法](../html/4e4be36a-8ce8-a41b-6ec4-4f82479c6f23.htm "ReadFromCoreServer 方法 ")

[ReadFromCoreServerAsync 方法](../html/e4cb7ebe-d42a-ea17-8972-b246a42e9a90.htm "ReadFromCoreServerAsync 方法 ")

[ReadIntegerVariable 方法](../html/df3171c4-e18b-a52c-9032-bd789e7388a0.htm "ReadIntegerVariable 方法 ")

[ReadIntegerVariableAsync 方法](../html/954c86a5-37e1-f32b-3a6e-7fc4274cbb6f.htm "ReadIntegerVariableAsync 方法 ")

[ReadJSeq 方法](../html/d7a72f10-2923-e045-8c71-e8b43c551c5a.htm "ReadJSeq 方法 ")

[ReadJSeqAsync 方法](../html/ddfc3321-88c7-9d2a-b71e-9315ce77864d.htm "ReadJSeqAsync 方法 ")

[ReadPOSC 方法](../html/fea57b6b-1d2f-0547-73de-6120ee5de705.htm "ReadPOSC 方法 ")

[ReadPOSCAsync 方法](../html/5f5731c4-a35a-affa-3d76-7b5f50fe63a1.htm "ReadPOSCAsync 方法 ")

[ReadPOSJ 方法](../html/14045444-1d70-3bdf-9848-a886c204da46.htm "ReadPOSJ 方法 ")

[ReadPOSJAsync 方法](../html/0ed6f125-5c82-6a30-1a4e-a6a773dda5d8.htm "ReadPOSJAsync 方法 ")

[ReadRealVariable 方法](../html/5f6c460a-0bb9-32fb-8aa8-33b99dc8be5e.htm "ReadRealVariable 方法 ")

[ReadRealVariableAsync 方法](../html/996ae64e-84f9-ccf6-fa2a-de4e4a1b13c9.htm "ReadRealVariableAsync 方法 ")

[ReadStats 方法](../html/3a24171c-f204-8c63-cdbc-8cf23a965b75.htm "ReadStats 方法 ")

[ReadStatsAsync 方法](../html/43a75b5c-c5af-e81f-286c-a4aa2a7a898b.htm "ReadStatsAsync 方法 ")

[ReadString 方法](../html/dbc09266-42e8-b939-aea5-540d528d19ee.htm "ReadString 方法 ")

[ReadStringAsync 方法](../html/8fa97613-eb8d-c1e1-33f6-925dcf1a1ac0.htm "ReadStringAsync 方法 ")

[ReadStringVariable 方法](../html/5dc9095d-86e8-d434-0027-83797a1d8c0a.htm "ReadStringVariable 方法 ")

[ReadStringVariableAsync 方法](../html/f850e762-bf8a-6ee7-c18c-3496d70244fe.htm "ReadStringVariableAsync 方法 ")

[ReadUFrame 方法](../html/100a1955-fc0a-871a-98f9-bb361c422a3f.htm "ReadUFrame 方法 ")

[ReadUFrameAsync 方法](../html/f63d72f5-eaf4-ae5d-e80a-d881f93a878e.htm "ReadUFrameAsync 方法 ")

[Reset 方法](../html/bf61707b-875c-4b42-2732-9b2fb69535ac.htm "Reset 方法 ")

[ResetAsync 方法](../html/f135af7f-0250-eb16-48fe-959dd0aeeb72.htm "ResetAsync 方法 ")

[SetMJ 方法](../html/9cd8c255-f908-6f4e-da07-33edc743dbf8.htm "SetMJ 方法 ")

[SetMJAsync 方法](../html/74b1cb5d-83f7-61b7-8eb0-0e10a675b964.htm "SetMJAsync 方法 ")

[Start 方法](../html/39764e6e-fde9-c99b-9034-99a64232aef6.htm "Start 方法 ")

[StartAsync 方法](../html/ddefed38-f8b8-4ae1-73e3-eecbfae54043.htm "StartAsync 方法 ")

[Svon 方法](../html/1f83c8cb-fea3-d99f-c7ad-f12f39bd9dee.htm "Svon 方法 ")

[SvonAsync 方法](../html/9b51697e-6164-280a-d37c-1cc16f467b2d.htm "SvonAsync 方法 ")

[ToString 方法](../html/1a5b666e-cae2-1664-908f-56866da693e5.htm "ToString 方法 ")

[Write 方法](../html/27720779-dc92-a6b4-c792-940ec2a10d22.htm "Write 方法 ")

[WriteAsync 方法](../html/e1646488-a935-0a5b-ec6a-4d04a3e69d38.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| YRC1000TcpNetHLock 方法 |

设定示教编程器和 I/O的操作信号的联锁。 状态参数 False: OFF，True: ON  
Set the interlock between the programming pendant and the operation signal of I/O. Status parameter False: OFF，True: ON

**命名空间：**
 [HslCommunication.Robot.YASKAWA](3d67023f-78aa-c5d9-85b6-f18f5dfbc176.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public OperateResult HLock(
	bool status
)
```

```
Public Function HLock ( 
	status As Boolean
) As OperateResult
```

```
public:
OperateResult^ HLock(
	bool status
)
```

```
member HLock : 
        status : bool -> OperateResult 
```

#### 参数

status
:   类型：SystemBoolean  
    状态参数 False: OFF，True: ON

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否设定成功

![](../icons/SectionExpanded.png)备注

联锁为ON时，仅可执行以下操作。

1. 示教编程器的非常停止
2. Ｉ /O 的模式切换， 外部启动， 外部伺服ON，循环切换， I/O 禁止、 PP/PANEL 禁止、 主程序调出以外的输入信号

示教编程器在编辑中或者通过其他的功能访问文件时，不能使用HLOCK.

![](../icons/SectionExpanded.png)参见

#### 引用

[YRC1000TcpNet 类](cdcdc1cc-17cf-da27-7d20-ca6d6db1cd7e.htm)

[HslCommunication.Robot.YASKAWA 命名空间](3d67023f-78aa-c5d9-85b6-f18f5dfbc176.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## HLockAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/f0eb6982-a92b-4481-6435-33fa2b841b3c.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.YASKAWA](../html/3d67023f-78aa-c5d9-85b6-f18f5dfbc176.htm "HslCommunication.Robot.YASKAWA")

[YRC1000TcpNet 类](../html/cdcdc1cc-17cf-da27-7d20-ca6d6db1cd7e.htm "YRC1000TcpNet 类")

[YRC1000TcpNet 方法](../html/7a26c354-eee0-b601-1d8a-51198e40e3fa.htm "YRC1000TcpNet 方法")

[Cancel 方法](../html/be66b6de-fe15-95b5-6969-70518f086388.htm "Cancel 方法 ")

[CancelAsync 方法](../html/fedb71cd-919c-c002-7556-551ae30fb10e.htm "CancelAsync 方法 ")

[Cycle 方法](../html/5a144fc7-5bc5-7250-6c13-826c172bd2ce.htm "Cycle 方法 ")

[CycleAsync 方法](../html/82bfe8a9-00b3-b34e-2a82-6d022becdf9a.htm "CycleAsync 方法 ")

[Delete 方法](../html/c7bfb022-46f4-8cdb-1b50-ffc1e5eb716b.htm "Delete 方法 ")

[DeleteAsync 方法](../html/c5a62bdd-3432-abcd-9d03-1ff7e30bd447.htm "DeleteAsync 方法 ")

[GetNewNetMessage 方法](../html/888815ce-da5c-0cb3-c6f9-800ec37a19ab.htm "GetNewNetMessage 方法 ")

[HLock 方法](../html/9c4ed34f-89fd-62e2-f707-08d1348c7541.htm "HLock 方法 ")

[HLockAsync 方法](../html/f0eb6982-a92b-4481-6435-33fa2b841b3c.htm "HLockAsync 方法 ")

[Hold 方法](../html/bf4ccc0a-1576-83ca-49e3-fde86b3a3738.htm "Hold 方法 ")

[HoldAsync 方法](../html/8c84a7d8-3677-0b6b-bdb5-44aa61ed1fcd.htm "HoldAsync 方法 ")

[InitializationOnConnect 方法](../html/5bd0274d-b3e7-f797-2dc0-c0dd2f330694.htm "InitializationOnConnect 方法 ")

[InitializationOnConnectAsync 方法](../html/e2f434c6-2ad9-30e8-8f2c-dcab95f5ec2b.htm "InitializationOnConnectAsync 方法 ")

[IORead 方法](../html/4539f381-2df3-4315-084c-7fcc81d523cb.htm "IORead 方法 ")

[IOReadAsync 方法](../html/5d550dc5-c02c-7245-122e-c3624e8d9c14.htm "IOReadAsync 方法 ")

[IOWrite 方法](../html/0dfecba5-f497-4b0f-55a5-339c7202d742.htm "IOWrite 方法 ")

[IOWriteAsync 方法](../html/64b5c202-29ff-0737-9c62-a4735bc00911.htm "IOWriteAsync 方法 ")

[JSeq 方法](../html/7057228b-aa47-3318-248d-1594588783bd.htm "JSeq 方法 ")

[JSeqAsync 方法](../html/c4eb7232-63e1-cd10-f89e-7f58072e3e9c.htm "JSeqAsync 方法 ")

[Mode 方法](../html/cabcebca-60ed-7689-e44a-e227936dcc51.htm "Mode 方法 ")

[ModeAsync 方法](../html/e0956a5d-5ae1-f884-52dd-0e91b85828b5.htm "ModeAsync 方法 ")

[MoveJ 方法](../html/10990eab-b2f9-e4b8-3c88-5bffc78aa34a.htm "MoveJ 方法 ")

[MoveJAsync 方法](../html/5b05904a-27a0-b4f9-adc7-1fe2cc5b14a1.htm "MoveJAsync 方法 ")

[MSDP 方法](../html/1e78f8b5-7336-7fdf-a7e0-e40af1d3de7f.htm "MSDP 方法 ")

[MSDPAsync 方法](../html/c2a18a60-0dae-112b-af16-0e028e6a203a.htm "MSDPAsync 方法 ")

[Read 方法](../html/69082cc9-c365-8fc4-36df-082c0c6c22c3.htm "Read 方法 ")

[ReadALARM 方法](../html/b56d10b7-3ebc-effd-c8f4-1dbf2d9f35a0.htm "ReadALARM 方法 ")

[ReadALARMAsync 方法](../html/6bd50e3f-6923-2c34-0e8d-d08ba82840f1.htm "ReadALARMAsync 方法 ")

[ReadAsync 方法](../html/185f22c2-a4f1-985e-dd8e-4bf7107c5ae7.htm "ReadAsync 方法 ")

[ReadByCommand 方法](../html/80dbd712-4782-b205-7e73-977c37f5a214.htm "ReadByCommand 方法 ")

[ReadByCommandAsync 方法](../html/2ac7651a-e36e-1635-0a61-b855722693a1.htm "ReadByCommandAsync 方法 ")

[ReadByteVariable 方法](../html/e22f2610-10f3-8472-92c2-3f9e904ebb2b.htm "ReadByteVariable 方法 ")

[ReadByteVariableAsync 方法](../html/d01b8a15-6079-9c57-6936-e0c569cf5194.htm "ReadByteVariableAsync 方法 ")

[ReadDoubleIntegerVariable 方法](../html/beab7b65-77db-3c62-6f95-ccf796ecac03.htm "ReadDoubleIntegerVariable 方法 ")

[ReadDoubleIntegerVariableAsync 方法](../html/c6845ca5-1d25-d613-89f9-28da2515cf24.htm "ReadDoubleIntegerVariableAsync 方法 ")

[ReadFromCoreServer 方法](../html/4e4be36a-8ce8-a41b-6ec4-4f82479c6f23.htm "ReadFromCoreServer 方法 ")

[ReadFromCoreServerAsync 方法](../html/e4cb7ebe-d42a-ea17-8972-b246a42e9a90.htm "ReadFromCoreServerAsync 方法 ")

[ReadIntegerVariable 方法](../html/df3171c4-e18b-a52c-9032-bd789e7388a0.htm "ReadIntegerVariable 方法 ")

[ReadIntegerVariableAsync 方法](../html/954c86a5-37e1-f32b-3a6e-7fc4274cbb6f.htm "ReadIntegerVariableAsync 方法 ")

[ReadJSeq 方法](../html/d7a72f10-2923-e045-8c71-e8b43c551c5a.htm "ReadJSeq 方法 ")

[ReadJSeqAsync 方法](../html/ddfc3321-88c7-9d2a-b71e-9315ce77864d.htm "ReadJSeqAsync 方法 ")

[ReadPOSC 方法](../html/fea57b6b-1d2f-0547-73de-6120ee5de705.htm "ReadPOSC 方法 ")

[ReadPOSCAsync 方法](../html/5f5731c4-a35a-affa-3d76-7b5f50fe63a1.htm "ReadPOSCAsync 方法 ")

[ReadPOSJ 方法](../html/14045444-1d70-3bdf-9848-a886c204da46.htm "ReadPOSJ 方法 ")

[ReadPOSJAsync 方法](../html/0ed6f125-5c82-6a30-1a4e-a6a773dda5d8.htm "ReadPOSJAsync 方法 ")

[ReadRealVariable 方法](../html/5f6c460a-0bb9-32fb-8aa8-33b99dc8be5e.htm "ReadRealVariable 方法 ")

[ReadRealVariableAsync 方法](../html/996ae64e-84f9-ccf6-fa2a-de4e4a1b13c9.htm "ReadRealVariableAsync 方法 ")

[ReadStats 方法](../html/3a24171c-f204-8c63-cdbc-8cf23a965b75.htm "ReadStats 方法 ")

[ReadStatsAsync 方法](../html/43a75b5c-c5af-e81f-286c-a4aa2a7a898b.htm "ReadStatsAsync 方法 ")

[ReadString 方法](../html/dbc09266-42e8-b939-aea5-540d528d19ee.htm "ReadString 方法 ")

[ReadStringAsync 方法](../html/8fa97613-eb8d-c1e1-33f6-925dcf1a1ac0.htm "ReadStringAsync 方法 ")

[ReadStringVariable 方法](../html/5dc9095d-86e8-d434-0027-83797a1d8c0a.htm "ReadStringVariable 方法 ")

[ReadStringVariableAsync 方法](../html/f850e762-bf8a-6ee7-c18c-3496d70244fe.htm "ReadStringVariableAsync 方法 ")

[ReadUFrame 方法](../html/100a1955-fc0a-871a-98f9-bb361c422a3f.htm "ReadUFrame 方法 ")

[ReadUFrameAsync 方法](../html/f63d72f5-eaf4-ae5d-e80a-d881f93a878e.htm "ReadUFrameAsync 方法 ")

[Reset 方法](../html/bf61707b-875c-4b42-2732-9b2fb69535ac.htm "Reset 方法 ")

[ResetAsync 方法](../html/f135af7f-0250-eb16-48fe-959dd0aeeb72.htm "ResetAsync 方法 ")

[SetMJ 方法](../html/9cd8c255-f908-6f4e-da07-33edc743dbf8.htm "SetMJ 方法 ")

[SetMJAsync 方法](../html/74b1cb5d-83f7-61b7-8eb0-0e10a675b964.htm "SetMJAsync 方法 ")

[Start 方法](../html/39764e6e-fde9-c99b-9034-99a64232aef6.htm "Start 方法 ")

[StartAsync 方法](../html/ddefed38-f8b8-4ae1-73e3-eecbfae54043.htm "StartAsync 方法 ")

[Svon 方法](../html/1f83c8cb-fea3-d99f-c7ad-f12f39bd9dee.htm "Svon 方法 ")

[SvonAsync 方法](../html/9b51697e-6164-280a-d37c-1cc16f467b2d.htm "SvonAsync 方法 ")

[ToString 方法](../html/1a5b666e-cae2-1664-908f-56866da693e5.htm "ToString 方法 ")

[Write 方法](../html/27720779-dc92-a6b4-c792-940ec2a10d22.htm "Write 方法 ")

[WriteAsync 方法](../html/e1646488-a935-0a5b-ec6a-4d04a3e69d38.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| YRC1000TcpNetHLockAsync 方法 |

设定示教编程器和 I/O的操作信号的联锁。 状态参数 False: OFF，True: ON  
Set the interlock between the programming pendant and the operation signal of I/O. Status parameter False: OFF，True: ON

**命名空间：**
 [HslCommunication.Robot.YASKAWA](3d67023f-78aa-c5d9-85b6-f18f5dfbc176.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public Task<OperateResult> HLockAsync(
	bool status
)
```

```
Public Function HLockAsync ( 
	status As Boolean
) As Task(Of OperateResult)
```

```
public:
Task<OperateResult^>^ HLockAsync(
	bool status
)
```

```
member HLockAsync : 
        status : bool -> Task<OperateResult> 
```

#### 参数

status
:   类型：SystemBoolean  
    状态参数 False: OFF，True: ON

#### 返回值

类型：Task[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否设定成功

![](../icons/SectionExpanded.png)备注

联锁为ON时，仅可执行以下操作。

1. 示教编程器的非常停止
2. Ｉ /O 的模式切换， 外部启动， 外部伺服ON，循环切换， I/O 禁止、 PP/PANEL 禁止、 主程序调出以外的输入信号

示教编程器在编辑中或者通过其他的功能访问文件时，不能使用HLOCK.

![](../icons/SectionExpanded.png)参见

#### 引用

[YRC1000TcpNet 类](cdcdc1cc-17cf-da27-7d20-ca6d6db1cd7e.htm)

[HslCommunication.Robot.YASKAWA 命名空间](3d67023f-78aa-c5d9-85b6-f18f5dfbc176.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Hold 方法 

[原文連結](http://api.hslcommunication.cn/html/bf4ccc0a-1576-83ca-49e3-fde86b3a3738.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.YASKAWA](../html/3d67023f-78aa-c5d9-85b6-f18f5dfbc176.htm "HslCommunication.Robot.YASKAWA")

[YRC1000TcpNet 类](../html/cdcdc1cc-17cf-da27-7d20-ca6d6db1cd7e.htm "YRC1000TcpNet 类")

[YRC1000TcpNet 方法](../html/7a26c354-eee0-b601-1d8a-51198e40e3fa.htm "YRC1000TcpNet 方法")

[Cancel 方法](../html/be66b6de-fe15-95b5-6969-70518f086388.htm "Cancel 方法 ")

[CancelAsync 方法](../html/fedb71cd-919c-c002-7556-551ae30fb10e.htm "CancelAsync 方法 ")

[Cycle 方法](../html/5a144fc7-5bc5-7250-6c13-826c172bd2ce.htm "Cycle 方法 ")

[CycleAsync 方法](../html/82bfe8a9-00b3-b34e-2a82-6d022becdf9a.htm "CycleAsync 方法 ")

[Delete 方法](../html/c7bfb022-46f4-8cdb-1b50-ffc1e5eb716b.htm "Delete 方法 ")

[DeleteAsync 方法](../html/c5a62bdd-3432-abcd-9d03-1ff7e30bd447.htm "DeleteAsync 方法 ")

[GetNewNetMessage 方法](../html/888815ce-da5c-0cb3-c6f9-800ec37a19ab.htm "GetNewNetMessage 方法 ")

[HLock 方法](../html/9c4ed34f-89fd-62e2-f707-08d1348c7541.htm "HLock 方法 ")

[HLockAsync 方法](../html/f0eb6982-a92b-4481-6435-33fa2b841b3c.htm "HLockAsync 方法 ")

[Hold 方法](../html/bf4ccc0a-1576-83ca-49e3-fde86b3a3738.htm "Hold 方法 ")

[HoldAsync 方法](../html/8c84a7d8-3677-0b6b-bdb5-44aa61ed1fcd.htm "HoldAsync 方法 ")

[InitializationOnConnect 方法](../html/5bd0274d-b3e7-f797-2dc0-c0dd2f330694.htm "InitializationOnConnect 方法 ")

[InitializationOnConnectAsync 方法](../html/e2f434c6-2ad9-30e8-8f2c-dcab95f5ec2b.htm "InitializationOnConnectAsync 方法 ")

[IORead 方法](../html/4539f381-2df3-4315-084c-7fcc81d523cb.htm "IORead 方法 ")

[IOReadAsync 方法](../html/5d550dc5-c02c-7245-122e-c3624e8d9c14.htm "IOReadAsync 方法 ")

[IOWrite 方法](../html/0dfecba5-f497-4b0f-55a5-339c7202d742.htm "IOWrite 方法 ")

[IOWriteAsync 方法](../html/64b5c202-29ff-0737-9c62-a4735bc00911.htm "IOWriteAsync 方法 ")

[JSeq 方法](../html/7057228b-aa47-3318-248d-1594588783bd.htm "JSeq 方法 ")

[JSeqAsync 方法](../html/c4eb7232-63e1-cd10-f89e-7f58072e3e9c.htm "JSeqAsync 方法 ")

[Mode 方法](../html/cabcebca-60ed-7689-e44a-e227936dcc51.htm "Mode 方法 ")

[ModeAsync 方法](../html/e0956a5d-5ae1-f884-52dd-0e91b85828b5.htm "ModeAsync 方法 ")

[MoveJ 方法](../html/10990eab-b2f9-e4b8-3c88-5bffc78aa34a.htm "MoveJ 方法 ")

[MoveJAsync 方法](../html/5b05904a-27a0-b4f9-adc7-1fe2cc5b14a1.htm "MoveJAsync 方法 ")

[MSDP 方法](../html/1e78f8b5-7336-7fdf-a7e0-e40af1d3de7f.htm "MSDP 方法 ")

[MSDPAsync 方法](../html/c2a18a60-0dae-112b-af16-0e028e6a203a.htm "MSDPAsync 方法 ")

[Read 方法](../html/69082cc9-c365-8fc4-36df-082c0c6c22c3.htm "Read 方法 ")

[ReadALARM 方法](../html/b56d10b7-3ebc-effd-c8f4-1dbf2d9f35a0.htm "ReadALARM 方法 ")

[ReadALARMAsync 方法](../html/6bd50e3f-6923-2c34-0e8d-d08ba82840f1.htm "ReadALARMAsync 方法 ")

[ReadAsync 方法](../html/185f22c2-a4f1-985e-dd8e-4bf7107c5ae7.htm "ReadAsync 方法 ")

[ReadByCommand 方法](../html/80dbd712-4782-b205-7e73-977c37f5a214.htm "ReadByCommand 方法 ")

[ReadByCommandAsync 方法](../html/2ac7651a-e36e-1635-0a61-b855722693a1.htm "ReadByCommandAsync 方法 ")

[ReadByteVariable 方法](../html/e22f2610-10f3-8472-92c2-3f9e904ebb2b.htm "ReadByteVariable 方法 ")

[ReadByteVariableAsync 方法](../html/d01b8a15-6079-9c57-6936-e0c569cf5194.htm "ReadByteVariableAsync 方法 ")

[ReadDoubleIntegerVariable 方法](../html/beab7b65-77db-3c62-6f95-ccf796ecac03.htm "ReadDoubleIntegerVariable 方法 ")

[ReadDoubleIntegerVariableAsync 方法](../html/c6845ca5-1d25-d613-89f9-28da2515cf24.htm "ReadDoubleIntegerVariableAsync 方法 ")

[ReadFromCoreServer 方法](../html/4e4be36a-8ce8-a41b-6ec4-4f82479c6f23.htm "ReadFromCoreServer 方法 ")

[ReadFromCoreServerAsync 方法](../html/e4cb7ebe-d42a-ea17-8972-b246a42e9a90.htm "ReadFromCoreServerAsync 方法 ")

[ReadIntegerVariable 方法](../html/df3171c4-e18b-a52c-9032-bd789e7388a0.htm "ReadIntegerVariable 方法 ")

[ReadIntegerVariableAsync 方法](../html/954c86a5-37e1-f32b-3a6e-7fc4274cbb6f.htm "ReadIntegerVariableAsync 方法 ")

[ReadJSeq 方法](../html/d7a72f10-2923-e045-8c71-e8b43c551c5a.htm "ReadJSeq 方法 ")

[ReadJSeqAsync 方法](../html/ddfc3321-88c7-9d2a-b71e-9315ce77864d.htm "ReadJSeqAsync 方法 ")

[ReadPOSC 方法](../html/fea57b6b-1d2f-0547-73de-6120ee5de705.htm "ReadPOSC 方法 ")

[ReadPOSCAsync 方法](../html/5f5731c4-a35a-affa-3d76-7b5f50fe63a1.htm "ReadPOSCAsync 方法 ")

[ReadPOSJ 方法](../html/14045444-1d70-3bdf-9848-a886c204da46.htm "ReadPOSJ 方法 ")

[ReadPOSJAsync 方法](../html/0ed6f125-5c82-6a30-1a4e-a6a773dda5d8.htm "ReadPOSJAsync 方法 ")

[ReadRealVariable 方法](../html/5f6c460a-0bb9-32fb-8aa8-33b99dc8be5e.htm "ReadRealVariable 方法 ")

[ReadRealVariableAsync 方法](../html/996ae64e-84f9-ccf6-fa2a-de4e4a1b13c9.htm "ReadRealVariableAsync 方法 ")

[ReadStats 方法](../html/3a24171c-f204-8c63-cdbc-8cf23a965b75.htm "ReadStats 方法 ")

[ReadStatsAsync 方法](../html/43a75b5c-c5af-e81f-286c-a4aa2a7a898b.htm "ReadStatsAsync 方法 ")

[ReadString 方法](../html/dbc09266-42e8-b939-aea5-540d528d19ee.htm "ReadString 方法 ")

[ReadStringAsync 方法](../html/8fa97613-eb8d-c1e1-33f6-925dcf1a1ac0.htm "ReadStringAsync 方法 ")

[ReadStringVariable 方法](../html/5dc9095d-86e8-d434-0027-83797a1d8c0a.htm "ReadStringVariable 方法 ")

[ReadStringVariableAsync 方法](../html/f850e762-bf8a-6ee7-c18c-3496d70244fe.htm "ReadStringVariableAsync 方法 ")

[ReadUFrame 方法](../html/100a1955-fc0a-871a-98f9-bb361c422a3f.htm "ReadUFrame 方法 ")

[ReadUFrameAsync 方法](../html/f63d72f5-eaf4-ae5d-e80a-d881f93a878e.htm "ReadUFrameAsync 方法 ")

[Reset 方法](../html/bf61707b-875c-4b42-2732-9b2fb69535ac.htm "Reset 方法 ")

[ResetAsync 方法](../html/f135af7f-0250-eb16-48fe-959dd0aeeb72.htm "ResetAsync 方法 ")

[SetMJ 方法](../html/9cd8c255-f908-6f4e-da07-33edc743dbf8.htm "SetMJ 方法 ")

[SetMJAsync 方法](../html/74b1cb5d-83f7-61b7-8eb0-0e10a675b964.htm "SetMJAsync 方法 ")

[Start 方法](../html/39764e6e-fde9-c99b-9034-99a64232aef6.htm "Start 方法 ")

[StartAsync 方法](../html/ddefed38-f8b8-4ae1-73e3-eecbfae54043.htm "StartAsync 方法 ")

[Svon 方法](../html/1f83c8cb-fea3-d99f-c7ad-f12f39bd9dee.htm "Svon 方法 ")

[SvonAsync 方法](../html/9b51697e-6164-280a-d37c-1cc16f467b2d.htm "SvonAsync 方法 ")

[ToString 方法](../html/1a5b666e-cae2-1664-908f-56866da693e5.htm "ToString 方法 ")

[Write 方法](../html/27720779-dc92-a6b4-c792-940ec2a10d22.htm "Write 方法 ")

[WriteAsync 方法](../html/e1646488-a935-0a5b-ec6a-4d04a3e69d38.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| YRC1000TcpNetHold 方法 |

进行HOLD 的 ON/OFF 操作，状态参数 False: OFF操作，True: ON操作  
Perform HOLD ON operation, False: OFF，True: ON

**命名空间：**
 [HslCommunication.Robot.YASKAWA](3d67023f-78aa-c5d9-85b6-f18f5dfbc176.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public OperateResult Hold(
	bool status
)
```

```
Public Function Hold ( 
	status As Boolean
) As OperateResult
```

```
public:
OperateResult^ Hold(
	bool status
)
```

```
member Hold : 
        status : bool -> OperateResult 
```

#### 参数

status
:   类型：SystemBoolean  
    状态参数 False: OFF操作，True: ON操作

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否成功的HOLD操作

![](../icons/SectionExpanded.png)参见

#### 引用

[YRC1000TcpNet 类](cdcdc1cc-17cf-da27-7d20-ca6d6db1cd7e.htm)

[HslCommunication.Robot.YASKAWA 命名空间](3d67023f-78aa-c5d9-85b6-f18f5dfbc176.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## HoldAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/8c84a7d8-3677-0b6b-bdb5-44aa61ed1fcd.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.YASKAWA](../html/3d67023f-78aa-c5d9-85b6-f18f5dfbc176.htm "HslCommunication.Robot.YASKAWA")

[YRC1000TcpNet 类](../html/cdcdc1cc-17cf-da27-7d20-ca6d6db1cd7e.htm "YRC1000TcpNet 类")

[YRC1000TcpNet 方法](../html/7a26c354-eee0-b601-1d8a-51198e40e3fa.htm "YRC1000TcpNet 方法")

[Cancel 方法](../html/be66b6de-fe15-95b5-6969-70518f086388.htm "Cancel 方法 ")

[CancelAsync 方法](../html/fedb71cd-919c-c002-7556-551ae30fb10e.htm "CancelAsync 方法 ")

[Cycle 方法](../html/5a144fc7-5bc5-7250-6c13-826c172bd2ce.htm "Cycle 方法 ")

[CycleAsync 方法](../html/82bfe8a9-00b3-b34e-2a82-6d022becdf9a.htm "CycleAsync 方法 ")

[Delete 方法](../html/c7bfb022-46f4-8cdb-1b50-ffc1e5eb716b.htm "Delete 方法 ")

[DeleteAsync 方法](../html/c5a62bdd-3432-abcd-9d03-1ff7e30bd447.htm "DeleteAsync 方法 ")

[GetNewNetMessage 方法](../html/888815ce-da5c-0cb3-c6f9-800ec37a19ab.htm "GetNewNetMessage 方法 ")

[HLock 方法](../html/9c4ed34f-89fd-62e2-f707-08d1348c7541.htm "HLock 方法 ")

[HLockAsync 方法](../html/f0eb6982-a92b-4481-6435-33fa2b841b3c.htm "HLockAsync 方法 ")

[Hold 方法](../html/bf4ccc0a-1576-83ca-49e3-fde86b3a3738.htm "Hold 方法 ")

[HoldAsync 方法](../html/8c84a7d8-3677-0b6b-bdb5-44aa61ed1fcd.htm "HoldAsync 方法 ")

[InitializationOnConnect 方法](../html/5bd0274d-b3e7-f797-2dc0-c0dd2f330694.htm "InitializationOnConnect 方法 ")

[InitializationOnConnectAsync 方法](../html/e2f434c6-2ad9-30e8-8f2c-dcab95f5ec2b.htm "InitializationOnConnectAsync 方法 ")

[IORead 方法](../html/4539f381-2df3-4315-084c-7fcc81d523cb.htm "IORead 方法 ")

[IOReadAsync 方法](../html/5d550dc5-c02c-7245-122e-c3624e8d9c14.htm "IOReadAsync 方法 ")

[IOWrite 方法](../html/0dfecba5-f497-4b0f-55a5-339c7202d742.htm "IOWrite 方法 ")

[IOWriteAsync 方法](../html/64b5c202-29ff-0737-9c62-a4735bc00911.htm "IOWriteAsync 方法 ")

[JSeq 方法](../html/7057228b-aa47-3318-248d-1594588783bd.htm "JSeq 方法 ")

[JSeqAsync 方法](../html/c4eb7232-63e1-cd10-f89e-7f58072e3e9c.htm "JSeqAsync 方法 ")

[Mode 方法](../html/cabcebca-60ed-7689-e44a-e227936dcc51.htm "Mode 方法 ")

[ModeAsync 方法](../html/e0956a5d-5ae1-f884-52dd-0e91b85828b5.htm "ModeAsync 方法 ")

[MoveJ 方法](../html/10990eab-b2f9-e4b8-3c88-5bffc78aa34a.htm "MoveJ 方法 ")

[MoveJAsync 方法](../html/5b05904a-27a0-b4f9-adc7-1fe2cc5b14a1.htm "MoveJAsync 方法 ")

[MSDP 方法](../html/1e78f8b5-7336-7fdf-a7e0-e40af1d3de7f.htm "MSDP 方法 ")

[MSDPAsync 方法](../html/c2a18a60-0dae-112b-af16-0e028e6a203a.htm "MSDPAsync 方法 ")

[Read 方法](../html/69082cc9-c365-8fc4-36df-082c0c6c22c3.htm "Read 方法 ")

[ReadALARM 方法](../html/b56d10b7-3ebc-effd-c8f4-1dbf2d9f35a0.htm "ReadALARM 方法 ")

[ReadALARMAsync 方法](../html/6bd50e3f-6923-2c34-0e8d-d08ba82840f1.htm "ReadALARMAsync 方法 ")

[ReadAsync 方法](../html/185f22c2-a4f1-985e-dd8e-4bf7107c5ae7.htm "ReadAsync 方法 ")

[ReadByCommand 方法](../html/80dbd712-4782-b205-7e73-977c37f5a214.htm "ReadByCommand 方法 ")

[ReadByCommandAsync 方法](../html/2ac7651a-e36e-1635-0a61-b855722693a1.htm "ReadByCommandAsync 方法 ")

[ReadByteVariable 方法](../html/e22f2610-10f3-8472-92c2-3f9e904ebb2b.htm "ReadByteVariable 方法 ")

[ReadByteVariableAsync 方法](../html/d01b8a15-6079-9c57-6936-e0c569cf5194.htm "ReadByteVariableAsync 方法 ")

[ReadDoubleIntegerVariable 方法](../html/beab7b65-77db-3c62-6f95-ccf796ecac03.htm "ReadDoubleIntegerVariable 方法 ")

[ReadDoubleIntegerVariableAsync 方法](../html/c6845ca5-1d25-d613-89f9-28da2515cf24.htm "ReadDoubleIntegerVariableAsync 方法 ")

[ReadFromCoreServer 方法](../html/4e4be36a-8ce8-a41b-6ec4-4f82479c6f23.htm "ReadFromCoreServer 方法 ")

[ReadFromCoreServerAsync 方法](../html/e4cb7ebe-d42a-ea17-8972-b246a42e9a90.htm "ReadFromCoreServerAsync 方法 ")

[ReadIntegerVariable 方法](../html/df3171c4-e18b-a52c-9032-bd789e7388a0.htm "ReadIntegerVariable 方法 ")

[ReadIntegerVariableAsync 方法](../html/954c86a5-37e1-f32b-3a6e-7fc4274cbb6f.htm "ReadIntegerVariableAsync 方法 ")

[ReadJSeq 方法](../html/d7a72f10-2923-e045-8c71-e8b43c551c5a.htm "ReadJSeq 方法 ")

[ReadJSeqAsync 方法](../html/ddfc3321-88c7-9d2a-b71e-9315ce77864d.htm "ReadJSeqAsync 方法 ")

[ReadPOSC 方法](../html/fea57b6b-1d2f-0547-73de-6120ee5de705.htm "ReadPOSC 方法 ")

[ReadPOSCAsync 方法](../html/5f5731c4-a35a-affa-3d76-7b5f50fe63a1.htm "ReadPOSCAsync 方法 ")

[ReadPOSJ 方法](../html/14045444-1d70-3bdf-9848-a886c204da46.htm "ReadPOSJ 方法 ")

[ReadPOSJAsync 方法](../html/0ed6f125-5c82-6a30-1a4e-a6a773dda5d8.htm "ReadPOSJAsync 方法 ")

[ReadRealVariable 方法](../html/5f6c460a-0bb9-32fb-8aa8-33b99dc8be5e.htm "ReadRealVariable 方法 ")

[ReadRealVariableAsync 方法](../html/996ae64e-84f9-ccf6-fa2a-de4e4a1b13c9.htm "ReadRealVariableAsync 方法 ")

[ReadStats 方法](../html/3a24171c-f204-8c63-cdbc-8cf23a965b75.htm "ReadStats 方法 ")

[ReadStatsAsync 方法](../html/43a75b5c-c5af-e81f-286c-a4aa2a7a898b.htm "ReadStatsAsync 方法 ")

[ReadString 方法](../html/dbc09266-42e8-b939-aea5-540d528d19ee.htm "ReadString 方法 ")

[ReadStringAsync 方法](../html/8fa97613-eb8d-c1e1-33f6-925dcf1a1ac0.htm "ReadStringAsync 方法 ")

[ReadStringVariable 方法](../html/5dc9095d-86e8-d434-0027-83797a1d8c0a.htm "ReadStringVariable 方法 ")

[ReadStringVariableAsync 方法](../html/f850e762-bf8a-6ee7-c18c-3496d70244fe.htm "ReadStringVariableAsync 方法 ")

[ReadUFrame 方法](../html/100a1955-fc0a-871a-98f9-bb361c422a3f.htm "ReadUFrame 方法 ")

[ReadUFrameAsync 方法](../html/f63d72f5-eaf4-ae5d-e80a-d881f93a878e.htm "ReadUFrameAsync 方法 ")

[Reset 方法](../html/bf61707b-875c-4b42-2732-9b2fb69535ac.htm "Reset 方法 ")

[ResetAsync 方法](../html/f135af7f-0250-eb16-48fe-959dd0aeeb72.htm "ResetAsync 方法 ")

[SetMJ 方法](../html/9cd8c255-f908-6f4e-da07-33edc743dbf8.htm "SetMJ 方法 ")

[SetMJAsync 方法](../html/74b1cb5d-83f7-61b7-8eb0-0e10a675b964.htm "SetMJAsync 方法 ")

[Start 方法](../html/39764e6e-fde9-c99b-9034-99a64232aef6.htm "Start 方法 ")

[StartAsync 方法](../html/ddefed38-f8b8-4ae1-73e3-eecbfae54043.htm "StartAsync 方法 ")

[Svon 方法](../html/1f83c8cb-fea3-d99f-c7ad-f12f39bd9dee.htm "Svon 方法 ")

[SvonAsync 方法](../html/9b51697e-6164-280a-d37c-1cc16f467b2d.htm "SvonAsync 方法 ")

[ToString 方法](../html/1a5b666e-cae2-1664-908f-56866da693e5.htm "ToString 方法 ")

[Write 方法](../html/27720779-dc92-a6b4-c792-940ec2a10d22.htm "Write 方法 ")

[WriteAsync 方法](../html/e1646488-a935-0a5b-ec6a-4d04a3e69d38.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| YRC1000TcpNetHoldAsync 方法 |

进行HOLD 的 ON/OFF 操作，状态参数 False: OFF操作，True: ON操作  
Perform HOLD ON operation, False: OFF，True: ON

**命名空间：**
 [HslCommunication.Robot.YASKAWA](3d67023f-78aa-c5d9-85b6-f18f5dfbc176.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public Task<OperateResult> HoldAsync(
	bool status
)
```

```
Public Function HoldAsync ( 
	status As Boolean
) As Task(Of OperateResult)
```

```
public:
Task<OperateResult^>^ HoldAsync(
	bool status
)
```

```
member HoldAsync : 
        status : bool -> Task<OperateResult> 
```

#### 参数

status
:   类型：SystemBoolean  
    状态参数 False: OFF操作，True: ON操作

#### 返回值

类型：Task[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否成功的HOLD操作

![](../icons/SectionExpanded.png)参见

#### 引用

[YRC1000TcpNet 类](cdcdc1cc-17cf-da27-7d20-ca6d6db1cd7e.htm)

[HslCommunication.Robot.YASKAWA 命名空间](3d67023f-78aa-c5d9-85b6-f18f5dfbc176.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## InitializationOnConnect 方法 

[原文連結](http://api.hslcommunication.cn/html/5bd0274d-b3e7-f797-2dc0-c0dd2f330694.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.YASKAWA](../html/3d67023f-78aa-c5d9-85b6-f18f5dfbc176.htm "HslCommunication.Robot.YASKAWA")

[YRC1000TcpNet 类](../html/cdcdc1cc-17cf-da27-7d20-ca6d6db1cd7e.htm "YRC1000TcpNet 类")

[YRC1000TcpNet 方法](../html/7a26c354-eee0-b601-1d8a-51198e40e3fa.htm "YRC1000TcpNet 方法")

[Cancel 方法](../html/be66b6de-fe15-95b5-6969-70518f086388.htm "Cancel 方法 ")

[CancelAsync 方法](../html/fedb71cd-919c-c002-7556-551ae30fb10e.htm "CancelAsync 方法 ")

[Cycle 方法](../html/5a144fc7-5bc5-7250-6c13-826c172bd2ce.htm "Cycle 方法 ")

[CycleAsync 方法](../html/82bfe8a9-00b3-b34e-2a82-6d022becdf9a.htm "CycleAsync 方法 ")

[Delete 方法](../html/c7bfb022-46f4-8cdb-1b50-ffc1e5eb716b.htm "Delete 方法 ")

[DeleteAsync 方法](../html/c5a62bdd-3432-abcd-9d03-1ff7e30bd447.htm "DeleteAsync 方法 ")

[GetNewNetMessage 方法](../html/888815ce-da5c-0cb3-c6f9-800ec37a19ab.htm "GetNewNetMessage 方法 ")

[HLock 方法](../html/9c4ed34f-89fd-62e2-f707-08d1348c7541.htm "HLock 方法 ")

[HLockAsync 方法](../html/f0eb6982-a92b-4481-6435-33fa2b841b3c.htm "HLockAsync 方法 ")

[Hold 方法](../html/bf4ccc0a-1576-83ca-49e3-fde86b3a3738.htm "Hold 方法 ")

[HoldAsync 方法](../html/8c84a7d8-3677-0b6b-bdb5-44aa61ed1fcd.htm "HoldAsync 方法 ")

[InitializationOnConnect 方法](../html/5bd0274d-b3e7-f797-2dc0-c0dd2f330694.htm "InitializationOnConnect 方法 ")

[InitializationOnConnectAsync 方法](../html/e2f434c6-2ad9-30e8-8f2c-dcab95f5ec2b.htm "InitializationOnConnectAsync 方法 ")

[IORead 方法](../html/4539f381-2df3-4315-084c-7fcc81d523cb.htm "IORead 方法 ")

[IOReadAsync 方法](../html/5d550dc5-c02c-7245-122e-c3624e8d9c14.htm "IOReadAsync 方法 ")

[IOWrite 方法](../html/0dfecba5-f497-4b0f-55a5-339c7202d742.htm "IOWrite 方法 ")

[IOWriteAsync 方法](../html/64b5c202-29ff-0737-9c62-a4735bc00911.htm "IOWriteAsync 方法 ")

[JSeq 方法](../html/7057228b-aa47-3318-248d-1594588783bd.htm "JSeq 方法 ")

[JSeqAsync 方法](../html/c4eb7232-63e1-cd10-f89e-7f58072e3e9c.htm "JSeqAsync 方法 ")

[Mode 方法](../html/cabcebca-60ed-7689-e44a-e227936dcc51.htm "Mode 方法 ")

[ModeAsync 方法](../html/e0956a5d-5ae1-f884-52dd-0e91b85828b5.htm "ModeAsync 方法 ")

[MoveJ 方法](../html/10990eab-b2f9-e4b8-3c88-5bffc78aa34a.htm "MoveJ 方法 ")

[MoveJAsync 方法](../html/5b05904a-27a0-b4f9-adc7-1fe2cc5b14a1.htm "MoveJAsync 方法 ")

[MSDP 方法](../html/1e78f8b5-7336-7fdf-a7e0-e40af1d3de7f.htm "MSDP 方法 ")

[MSDPAsync 方法](../html/c2a18a60-0dae-112b-af16-0e028e6a203a.htm "MSDPAsync 方法 ")

[Read 方法](../html/69082cc9-c365-8fc4-36df-082c0c6c22c3.htm "Read 方法 ")

[ReadALARM 方法](../html/b56d10b7-3ebc-effd-c8f4-1dbf2d9f35a0.htm "ReadALARM 方法 ")

[ReadALARMAsync 方法](../html/6bd50e3f-6923-2c34-0e8d-d08ba82840f1.htm "ReadALARMAsync 方法 ")

[ReadAsync 方法](../html/185f22c2-a4f1-985e-dd8e-4bf7107c5ae7.htm "ReadAsync 方法 ")

[ReadByCommand 方法](../html/80dbd712-4782-b205-7e73-977c37f5a214.htm "ReadByCommand 方法 ")

[ReadByCommandAsync 方法](../html/2ac7651a-e36e-1635-0a61-b855722693a1.htm "ReadByCommandAsync 方法 ")

[ReadByteVariable 方法](../html/e22f2610-10f3-8472-92c2-3f9e904ebb2b.htm "ReadByteVariable 方法 ")

[ReadByteVariableAsync 方法](../html/d01b8a15-6079-9c57-6936-e0c569cf5194.htm "ReadByteVariableAsync 方法 ")

[ReadDoubleIntegerVariable 方法](../html/beab7b65-77db-3c62-6f95-ccf796ecac03.htm "ReadDoubleIntegerVariable 方法 ")

[ReadDoubleIntegerVariableAsync 方法](../html/c6845ca5-1d25-d613-89f9-28da2515cf24.htm "ReadDoubleIntegerVariableAsync 方法 ")

[ReadFromCoreServer 方法](../html/4e4be36a-8ce8-a41b-6ec4-4f82479c6f23.htm "ReadFromCoreServer 方法 ")

[ReadFromCoreServerAsync 方法](../html/e4cb7ebe-d42a-ea17-8972-b246a42e9a90.htm "ReadFromCoreServerAsync 方法 ")

[ReadIntegerVariable 方法](../html/df3171c4-e18b-a52c-9032-bd789e7388a0.htm "ReadIntegerVariable 方法 ")

[ReadIntegerVariableAsync 方法](../html/954c86a5-37e1-f32b-3a6e-7fc4274cbb6f.htm "ReadIntegerVariableAsync 方法 ")

[ReadJSeq 方法](../html/d7a72f10-2923-e045-8c71-e8b43c551c5a.htm "ReadJSeq 方法 ")

[ReadJSeqAsync 方法](../html/ddfc3321-88c7-9d2a-b71e-9315ce77864d.htm "ReadJSeqAsync 方法 ")

[ReadPOSC 方法](../html/fea57b6b-1d2f-0547-73de-6120ee5de705.htm "ReadPOSC 方法 ")

[ReadPOSCAsync 方法](../html/5f5731c4-a35a-affa-3d76-7b5f50fe63a1.htm "ReadPOSCAsync 方法 ")

[ReadPOSJ 方法](../html/14045444-1d70-3bdf-9848-a886c204da46.htm "ReadPOSJ 方法 ")

[ReadPOSJAsync 方法](../html/0ed6f125-5c82-6a30-1a4e-a6a773dda5d8.htm "ReadPOSJAsync 方法 ")

[ReadRealVariable 方法](../html/5f6c460a-0bb9-32fb-8aa8-33b99dc8be5e.htm "ReadRealVariable 方法 ")

[ReadRealVariableAsync 方法](../html/996ae64e-84f9-ccf6-fa2a-de4e4a1b13c9.htm "ReadRealVariableAsync 方法 ")

[ReadStats 方法](../html/3a24171c-f204-8c63-cdbc-8cf23a965b75.htm "ReadStats 方法 ")

[ReadStatsAsync 方法](../html/43a75b5c-c5af-e81f-286c-a4aa2a7a898b.htm "ReadStatsAsync 方法 ")

[ReadString 方法](../html/dbc09266-42e8-b939-aea5-540d528d19ee.htm "ReadString 方法 ")

[ReadStringAsync 方法](../html/8fa97613-eb8d-c1e1-33f6-925dcf1a1ac0.htm "ReadStringAsync 方法 ")

[ReadStringVariable 方法](../html/5dc9095d-86e8-d434-0027-83797a1d8c0a.htm "ReadStringVariable 方法 ")

[ReadStringVariableAsync 方法](../html/f850e762-bf8a-6ee7-c18c-3496d70244fe.htm "ReadStringVariableAsync 方法 ")

[ReadUFrame 方法](../html/100a1955-fc0a-871a-98f9-bb361c422a3f.htm "ReadUFrame 方法 ")

[ReadUFrameAsync 方法](../html/f63d72f5-eaf4-ae5d-e80a-d881f93a878e.htm "ReadUFrameAsync 方法 ")

[Reset 方法](../html/bf61707b-875c-4b42-2732-9b2fb69535ac.htm "Reset 方法 ")

[ResetAsync 方法](../html/f135af7f-0250-eb16-48fe-959dd0aeeb72.htm "ResetAsync 方法 ")

[SetMJ 方法](../html/9cd8c255-f908-6f4e-da07-33edc743dbf8.htm "SetMJ 方法 ")

[SetMJAsync 方法](../html/74b1cb5d-83f7-61b7-8eb0-0e10a675b964.htm "SetMJAsync 方法 ")

[Start 方法](../html/39764e6e-fde9-c99b-9034-99a64232aef6.htm "Start 方法 ")

[StartAsync 方法](../html/ddefed38-f8b8-4ae1-73e3-eecbfae54043.htm "StartAsync 方法 ")

[Svon 方法](../html/1f83c8cb-fea3-d99f-c7ad-f12f39bd9dee.htm "Svon 方法 ")

[SvonAsync 方法](../html/9b51697e-6164-280a-d37c-1cc16f467b2d.htm "SvonAsync 方法 ")

[ToString 方法](../html/1a5b666e-cae2-1664-908f-56866da693e5.htm "ToString 方法 ")

[Write 方法](../html/27720779-dc92-a6b4-c792-940ec2a10d22.htm "Write 方法 ")

[WriteAsync 方法](../html/e1646488-a935-0a5b-ec6a-4d04a3e69d38.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| YRC1000TcpNetInitializationOnConnect 方法 |

before read data , the connection should be Initialized

**命名空间：**
 [HslCommunication.Robot.YASKAWA](3d67023f-78aa-c5d9-85b6-f18f5dfbc176.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
protected override OperateResult InitializationOnConnect(
	Socket socket
)
```

```
Protected Overrides Function InitializationOnConnect ( 
	socket As Socket
) As OperateResult
```

```
protected:
virtual OperateResult^ InitializationOnConnect(
	Socket^ socket
) override
```

```
abstract InitializationOnConnect : 
        socket : Socket -> OperateResult 
override InitializationOnConnect : 
        socket : Socket -> OperateResult
```

#### 参数

socket
:   类型：System.Net.SocketsSocket  
    connected socket

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
whether is the Initialization is success.

![](../icons/SectionExpanded.png)参见

#### 引用

[YRC1000TcpNet 类](cdcdc1cc-17cf-da27-7d20-ca6d6db1cd7e.htm)

[HslCommunication.Robot.YASKAWA 命名空间](3d67023f-78aa-c5d9-85b6-f18f5dfbc176.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## InitializationOnConnectAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/e2f434c6-2ad9-30e8-8f2c-dcab95f5ec2b.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.YASKAWA](../html/3d67023f-78aa-c5d9-85b6-f18f5dfbc176.htm "HslCommunication.Robot.YASKAWA")

[YRC1000TcpNet 类](../html/cdcdc1cc-17cf-da27-7d20-ca6d6db1cd7e.htm "YRC1000TcpNet 类")

[YRC1000TcpNet 方法](../html/7a26c354-eee0-b601-1d8a-51198e40e3fa.htm "YRC1000TcpNet 方法")

[Cancel 方法](../html/be66b6de-fe15-95b5-6969-70518f086388.htm "Cancel 方法 ")

[CancelAsync 方法](../html/fedb71cd-919c-c002-7556-551ae30fb10e.htm "CancelAsync 方法 ")

[Cycle 方法](../html/5a144fc7-5bc5-7250-6c13-826c172bd2ce.htm "Cycle 方法 ")

[CycleAsync 方法](../html/82bfe8a9-00b3-b34e-2a82-6d022becdf9a.htm "CycleAsync 方法 ")

[Delete 方法](../html/c7bfb022-46f4-8cdb-1b50-ffc1e5eb716b.htm "Delete 方法 ")

[DeleteAsync 方法](../html/c5a62bdd-3432-abcd-9d03-1ff7e30bd447.htm "DeleteAsync 方法 ")

[GetNewNetMessage 方法](../html/888815ce-da5c-0cb3-c6f9-800ec37a19ab.htm "GetNewNetMessage 方法 ")

[HLock 方法](../html/9c4ed34f-89fd-62e2-f707-08d1348c7541.htm "HLock 方法 ")

[HLockAsync 方法](../html/f0eb6982-a92b-4481-6435-33fa2b841b3c.htm "HLockAsync 方法 ")

[Hold 方法](../html/bf4ccc0a-1576-83ca-49e3-fde86b3a3738.htm "Hold 方法 ")

[HoldAsync 方法](../html/8c84a7d8-3677-0b6b-bdb5-44aa61ed1fcd.htm "HoldAsync 方法 ")

[InitializationOnConnect 方法](../html/5bd0274d-b3e7-f797-2dc0-c0dd2f330694.htm "InitializationOnConnect 方法 ")

[InitializationOnConnectAsync 方法](../html/e2f434c6-2ad9-30e8-8f2c-dcab95f5ec2b.htm "InitializationOnConnectAsync 方法 ")

[IORead 方法](../html/4539f381-2df3-4315-084c-7fcc81d523cb.htm "IORead 方法 ")

[IOReadAsync 方法](../html/5d550dc5-c02c-7245-122e-c3624e8d9c14.htm "IOReadAsync 方法 ")

[IOWrite 方法](../html/0dfecba5-f497-4b0f-55a5-339c7202d742.htm "IOWrite 方法 ")

[IOWriteAsync 方法](../html/64b5c202-29ff-0737-9c62-a4735bc00911.htm "IOWriteAsync 方法 ")

[JSeq 方法](../html/7057228b-aa47-3318-248d-1594588783bd.htm "JSeq 方法 ")

[JSeqAsync 方法](../html/c4eb7232-63e1-cd10-f89e-7f58072e3e9c.htm "JSeqAsync 方法 ")

[Mode 方法](../html/cabcebca-60ed-7689-e44a-e227936dcc51.htm "Mode 方法 ")

[ModeAsync 方法](../html/e0956a5d-5ae1-f884-52dd-0e91b85828b5.htm "ModeAsync 方法 ")

[MoveJ 方法](../html/10990eab-b2f9-e4b8-3c88-5bffc78aa34a.htm "MoveJ 方法 ")

[MoveJAsync 方法](../html/5b05904a-27a0-b4f9-adc7-1fe2cc5b14a1.htm "MoveJAsync 方法 ")

[MSDP 方法](../html/1e78f8b5-7336-7fdf-a7e0-e40af1d3de7f.htm "MSDP 方法 ")

[MSDPAsync 方法](../html/c2a18a60-0dae-112b-af16-0e028e6a203a.htm "MSDPAsync 方法 ")

[Read 方法](../html/69082cc9-c365-8fc4-36df-082c0c6c22c3.htm "Read 方法 ")

[ReadALARM 方法](../html/b56d10b7-3ebc-effd-c8f4-1dbf2d9f35a0.htm "ReadALARM 方法 ")

[ReadALARMAsync 方法](../html/6bd50e3f-6923-2c34-0e8d-d08ba82840f1.htm "ReadALARMAsync 方法 ")

[ReadAsync 方法](../html/185f22c2-a4f1-985e-dd8e-4bf7107c5ae7.htm "ReadAsync 方法 ")

[ReadByCommand 方法](../html/80dbd712-4782-b205-7e73-977c37f5a214.htm "ReadByCommand 方法 ")

[ReadByCommandAsync 方法](../html/2ac7651a-e36e-1635-0a61-b855722693a1.htm "ReadByCommandAsync 方法 ")

[ReadByteVariable 方法](../html/e22f2610-10f3-8472-92c2-3f9e904ebb2b.htm "ReadByteVariable 方法 ")

[ReadByteVariableAsync 方法](../html/d01b8a15-6079-9c57-6936-e0c569cf5194.htm "ReadByteVariableAsync 方法 ")

[ReadDoubleIntegerVariable 方法](../html/beab7b65-77db-3c62-6f95-ccf796ecac03.htm "ReadDoubleIntegerVariable 方法 ")

[ReadDoubleIntegerVariableAsync 方法](../html/c6845ca5-1d25-d613-89f9-28da2515cf24.htm "ReadDoubleIntegerVariableAsync 方法 ")

[ReadFromCoreServer 方法](../html/4e4be36a-8ce8-a41b-6ec4-4f82479c6f23.htm "ReadFromCoreServer 方法 ")

[ReadFromCoreServerAsync 方法](../html/e4cb7ebe-d42a-ea17-8972-b246a42e9a90.htm "ReadFromCoreServerAsync 方法 ")

[ReadIntegerVariable 方法](../html/df3171c4-e18b-a52c-9032-bd789e7388a0.htm "ReadIntegerVariable 方法 ")

[ReadIntegerVariableAsync 方法](../html/954c86a5-37e1-f32b-3a6e-7fc4274cbb6f.htm "ReadIntegerVariableAsync 方法 ")

[ReadJSeq 方法](../html/d7a72f10-2923-e045-8c71-e8b43c551c5a.htm "ReadJSeq 方法 ")

[ReadJSeqAsync 方法](../html/ddfc3321-88c7-9d2a-b71e-9315ce77864d.htm "ReadJSeqAsync 方法 ")

[ReadPOSC 方法](../html/fea57b6b-1d2f-0547-73de-6120ee5de705.htm "ReadPOSC 方法 ")

[ReadPOSCAsync 方法](../html/5f5731c4-a35a-affa-3d76-7b5f50fe63a1.htm "ReadPOSCAsync 方法 ")

[ReadPOSJ 方法](../html/14045444-1d70-3bdf-9848-a886c204da46.htm "ReadPOSJ 方法 ")

[ReadPOSJAsync 方法](../html/0ed6f125-5c82-6a30-1a4e-a6a773dda5d8.htm "ReadPOSJAsync 方法 ")

[ReadRealVariable 方法](../html/5f6c460a-0bb9-32fb-8aa8-33b99dc8be5e.htm "ReadRealVariable 方法 ")

[ReadRealVariableAsync 方法](../html/996ae64e-84f9-ccf6-fa2a-de4e4a1b13c9.htm "ReadRealVariableAsync 方法 ")

[ReadStats 方法](../html/3a24171c-f204-8c63-cdbc-8cf23a965b75.htm "ReadStats 方法 ")

[ReadStatsAsync 方法](../html/43a75b5c-c5af-e81f-286c-a4aa2a7a898b.htm "ReadStatsAsync 方法 ")

[ReadString 方法](../html/dbc09266-42e8-b939-aea5-540d528d19ee.htm "ReadString 方法 ")

[ReadStringAsync 方法](../html/8fa97613-eb8d-c1e1-33f6-925dcf1a1ac0.htm "ReadStringAsync 方法 ")

[ReadStringVariable 方法](../html/5dc9095d-86e8-d434-0027-83797a1d8c0a.htm "ReadStringVariable 方法 ")

[ReadStringVariableAsync 方法](../html/f850e762-bf8a-6ee7-c18c-3496d70244fe.htm "ReadStringVariableAsync 方法 ")

[ReadUFrame 方法](../html/100a1955-fc0a-871a-98f9-bb361c422a3f.htm "ReadUFrame 方法 ")

[ReadUFrameAsync 方法](../html/f63d72f5-eaf4-ae5d-e80a-d881f93a878e.htm "ReadUFrameAsync 方法 ")

[Reset 方法](../html/bf61707b-875c-4b42-2732-9b2fb69535ac.htm "Reset 方法 ")

[ResetAsync 方法](../html/f135af7f-0250-eb16-48fe-959dd0aeeb72.htm "ResetAsync 方法 ")

[SetMJ 方法](../html/9cd8c255-f908-6f4e-da07-33edc743dbf8.htm "SetMJ 方法 ")

[SetMJAsync 方法](../html/74b1cb5d-83f7-61b7-8eb0-0e10a675b964.htm "SetMJAsync 方法 ")

[Start 方法](../html/39764e6e-fde9-c99b-9034-99a64232aef6.htm "Start 方法 ")

[StartAsync 方法](../html/ddefed38-f8b8-4ae1-73e3-eecbfae54043.htm "StartAsync 方法 ")

[Svon 方法](../html/1f83c8cb-fea3-d99f-c7ad-f12f39bd9dee.htm "Svon 方法 ")

[SvonAsync 方法](../html/9b51697e-6164-280a-d37c-1cc16f467b2d.htm "SvonAsync 方法 ")

[ToString 方法](../html/1a5b666e-cae2-1664-908f-56866da693e5.htm "ToString 方法 ")

[Write 方法](../html/27720779-dc92-a6b4-c792-940ec2a10d22.htm "Write 方法 ")

[WriteAsync 方法](../html/e1646488-a935-0a5b-ec6a-4d04a3e69d38.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| YRC1000TcpNetInitializationOnConnectAsync 方法 |

before read data , the connection should be Initialized

**命名空间：**
 [HslCommunication.Robot.YASKAWA](3d67023f-78aa-c5d9-85b6-f18f5dfbc176.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
protected override Task<OperateResult> InitializationOnConnectAsync(
	Socket socket
)
```

```
Protected Overrides Function InitializationOnConnectAsync ( 
	socket As Socket
) As Task(Of OperateResult)
```

```
protected:
virtual Task<OperateResult^>^ InitializationOnConnectAsync(
	Socket^ socket
) override
```

```
abstract InitializationOnConnectAsync : 
        socket : Socket -> Task<OperateResult> 
override InitializationOnConnectAsync : 
        socket : Socket -> Task<OperateResult>
```

#### 参数

socket
:   类型：System.Net.SocketsSocket  
    connected socket

#### 返回值

类型：Task[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
whether is the Initialization is success.

![](../icons/SectionExpanded.png)参见

#### 引用

[YRC1000TcpNet 类](cdcdc1cc-17cf-da27-7d20-ca6d6db1cd7e.htm)

[HslCommunication.Robot.YASKAWA 命名空间](3d67023f-78aa-c5d9-85b6-f18f5dfbc176.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## IORead 方法 

[原文連結](http://api.hslcommunication.cn/html/4539f381-2df3-4315-084c-7fcc81d523cb.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.YASKAWA](../html/3d67023f-78aa-c5d9-85b6-f18f5dfbc176.htm "HslCommunication.Robot.YASKAWA")

[YRC1000TcpNet 类](../html/cdcdc1cc-17cf-da27-7d20-ca6d6db1cd7e.htm "YRC1000TcpNet 类")

[YRC1000TcpNet 方法](../html/7a26c354-eee0-b601-1d8a-51198e40e3fa.htm "YRC1000TcpNet 方法")

[Cancel 方法](../html/be66b6de-fe15-95b5-6969-70518f086388.htm "Cancel 方法 ")

[CancelAsync 方法](../html/fedb71cd-919c-c002-7556-551ae30fb10e.htm "CancelAsync 方法 ")

[Cycle 方法](../html/5a144fc7-5bc5-7250-6c13-826c172bd2ce.htm "Cycle 方法 ")

[CycleAsync 方法](../html/82bfe8a9-00b3-b34e-2a82-6d022becdf9a.htm "CycleAsync 方法 ")

[Delete 方法](../html/c7bfb022-46f4-8cdb-1b50-ffc1e5eb716b.htm "Delete 方法 ")

[DeleteAsync 方法](../html/c5a62bdd-3432-abcd-9d03-1ff7e30bd447.htm "DeleteAsync 方法 ")

[GetNewNetMessage 方法](../html/888815ce-da5c-0cb3-c6f9-800ec37a19ab.htm "GetNewNetMessage 方法 ")

[HLock 方法](../html/9c4ed34f-89fd-62e2-f707-08d1348c7541.htm "HLock 方法 ")

[HLockAsync 方法](../html/f0eb6982-a92b-4481-6435-33fa2b841b3c.htm "HLockAsync 方法 ")

[Hold 方法](../html/bf4ccc0a-1576-83ca-49e3-fde86b3a3738.htm "Hold 方法 ")

[HoldAsync 方法](../html/8c84a7d8-3677-0b6b-bdb5-44aa61ed1fcd.htm "HoldAsync 方法 ")

[InitializationOnConnect 方法](../html/5bd0274d-b3e7-f797-2dc0-c0dd2f330694.htm "InitializationOnConnect 方法 ")

[InitializationOnConnectAsync 方法](../html/e2f434c6-2ad9-30e8-8f2c-dcab95f5ec2b.htm "InitializationOnConnectAsync 方法 ")

[IORead 方法](../html/4539f381-2df3-4315-084c-7fcc81d523cb.htm "IORead 方法 ")

[IOReadAsync 方法](../html/5d550dc5-c02c-7245-122e-c3624e8d9c14.htm "IOReadAsync 方法 ")

[IOWrite 方法](../html/0dfecba5-f497-4b0f-55a5-339c7202d742.htm "IOWrite 方法 ")

[IOWriteAsync 方法](../html/64b5c202-29ff-0737-9c62-a4735bc00911.htm "IOWriteAsync 方法 ")

[JSeq 方法](../html/7057228b-aa47-3318-248d-1594588783bd.htm "JSeq 方法 ")

[JSeqAsync 方法](../html/c4eb7232-63e1-cd10-f89e-7f58072e3e9c.htm "JSeqAsync 方法 ")

[Mode 方法](../html/cabcebca-60ed-7689-e44a-e227936dcc51.htm "Mode 方法 ")

[ModeAsync 方法](../html/e0956a5d-5ae1-f884-52dd-0e91b85828b5.htm "ModeAsync 方法 ")

[MoveJ 方法](../html/10990eab-b2f9-e4b8-3c88-5bffc78aa34a.htm "MoveJ 方法 ")

[MoveJAsync 方法](../html/5b05904a-27a0-b4f9-adc7-1fe2cc5b14a1.htm "MoveJAsync 方法 ")

[MSDP 方法](../html/1e78f8b5-7336-7fdf-a7e0-e40af1d3de7f.htm "MSDP 方法 ")

[MSDPAsync 方法](../html/c2a18a60-0dae-112b-af16-0e028e6a203a.htm "MSDPAsync 方法 ")

[Read 方法](../html/69082cc9-c365-8fc4-36df-082c0c6c22c3.htm "Read 方法 ")

[ReadALARM 方法](../html/b56d10b7-3ebc-effd-c8f4-1dbf2d9f35a0.htm "ReadALARM 方法 ")

[ReadALARMAsync 方法](../html/6bd50e3f-6923-2c34-0e8d-d08ba82840f1.htm "ReadALARMAsync 方法 ")

[ReadAsync 方法](../html/185f22c2-a4f1-985e-dd8e-4bf7107c5ae7.htm "ReadAsync 方法 ")

[ReadByCommand 方法](../html/80dbd712-4782-b205-7e73-977c37f5a214.htm "ReadByCommand 方法 ")

[ReadByCommandAsync 方法](../html/2ac7651a-e36e-1635-0a61-b855722693a1.htm "ReadByCommandAsync 方法 ")

[ReadByteVariable 方法](../html/e22f2610-10f3-8472-92c2-3f9e904ebb2b.htm "ReadByteVariable 方法 ")

[ReadByteVariableAsync 方法](../html/d01b8a15-6079-9c57-6936-e0c569cf5194.htm "ReadByteVariableAsync 方法 ")

[ReadDoubleIntegerVariable 方法](../html/beab7b65-77db-3c62-6f95-ccf796ecac03.htm "ReadDoubleIntegerVariable 方法 ")

[ReadDoubleIntegerVariableAsync 方法](../html/c6845ca5-1d25-d613-89f9-28da2515cf24.htm "ReadDoubleIntegerVariableAsync 方法 ")

[ReadFromCoreServer 方法](../html/4e4be36a-8ce8-a41b-6ec4-4f82479c6f23.htm "ReadFromCoreServer 方法 ")

[ReadFromCoreServerAsync 方法](../html/e4cb7ebe-d42a-ea17-8972-b246a42e9a90.htm "ReadFromCoreServerAsync 方法 ")

[ReadIntegerVariable 方法](../html/df3171c4-e18b-a52c-9032-bd789e7388a0.htm "ReadIntegerVariable 方法 ")

[ReadIntegerVariableAsync 方法](../html/954c86a5-37e1-f32b-3a6e-7fc4274cbb6f.htm "ReadIntegerVariableAsync 方法 ")

[ReadJSeq 方法](../html/d7a72f10-2923-e045-8c71-e8b43c551c5a.htm "ReadJSeq 方法 ")

[ReadJSeqAsync 方法](../html/ddfc3321-88c7-9d2a-b71e-9315ce77864d.htm "ReadJSeqAsync 方法 ")

[ReadPOSC 方法](../html/fea57b6b-1d2f-0547-73de-6120ee5de705.htm "ReadPOSC 方法 ")

[ReadPOSCAsync 方法](../html/5f5731c4-a35a-affa-3d76-7b5f50fe63a1.htm "ReadPOSCAsync 方法 ")

[ReadPOSJ 方法](../html/14045444-1d70-3bdf-9848-a886c204da46.htm "ReadPOSJ 方法 ")

[ReadPOSJAsync 方法](../html/0ed6f125-5c82-6a30-1a4e-a6a773dda5d8.htm "ReadPOSJAsync 方法 ")

[ReadRealVariable 方法](../html/5f6c460a-0bb9-32fb-8aa8-33b99dc8be5e.htm "ReadRealVariable 方法 ")

[ReadRealVariableAsync 方法](../html/996ae64e-84f9-ccf6-fa2a-de4e4a1b13c9.htm "ReadRealVariableAsync 方法 ")

[ReadStats 方法](../html/3a24171c-f204-8c63-cdbc-8cf23a965b75.htm "ReadStats 方法 ")

[ReadStatsAsync 方法](../html/43a75b5c-c5af-e81f-286c-a4aa2a7a898b.htm "ReadStatsAsync 方法 ")

[ReadString 方法](../html/dbc09266-42e8-b939-aea5-540d528d19ee.htm "ReadString 方法 ")

[ReadStringAsync 方法](../html/8fa97613-eb8d-c1e1-33f6-925dcf1a1ac0.htm "ReadStringAsync 方法 ")

[ReadStringVariable 方法](../html/5dc9095d-86e8-d434-0027-83797a1d8c0a.htm "ReadStringVariable 方法 ")

[ReadStringVariableAsync 方法](../html/f850e762-bf8a-6ee7-c18c-3496d70244fe.htm "ReadStringVariableAsync 方法 ")

[ReadUFrame 方法](../html/100a1955-fc0a-871a-98f9-bb361c422a3f.htm "ReadUFrame 方法 ")

[ReadUFrameAsync 方法](../html/f63d72f5-eaf4-ae5d-e80a-d881f93a878e.htm "ReadUFrameAsync 方法 ")

[Reset 方法](../html/bf61707b-875c-4b42-2732-9b2fb69535ac.htm "Reset 方法 ")

[ResetAsync 方法](../html/f135af7f-0250-eb16-48fe-959dd0aeeb72.htm "ResetAsync 方法 ")

[SetMJ 方法](../html/9cd8c255-f908-6f4e-da07-33edc743dbf8.htm "SetMJ 方法 ")

[SetMJAsync 方法](../html/74b1cb5d-83f7-61b7-8eb0-0e10a675b964.htm "SetMJAsync 方法 ")

[Start 方法](../html/39764e6e-fde9-c99b-9034-99a64232aef6.htm "Start 方法 ")

[StartAsync 方法](../html/ddefed38-f8b8-4ae1-73e3-eecbfae54043.htm "StartAsync 方法 ")

[Svon 方法](../html/1f83c8cb-fea3-d99f-c7ad-f12f39bd9dee.htm "Svon 方法 ")

[SvonAsync 方法](../html/9b51697e-6164-280a-d37c-1cc16f467b2d.htm "SvonAsync 方法 ")

[ToString 方法](../html/1a5b666e-cae2-1664-908f-56866da693e5.htm "ToString 方法 ")

[Write 方法](../html/27720779-dc92-a6b4-c792-940ec2a10d22.htm "Write 方法 ")

[WriteAsync 方法](../html/e1646488-a935-0a5b-ec6a-4d04a3e69d38.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| YRC1000TcpNetIORead 方法 |

读取I/O 信号。 I/O 数据是每8个点输出，所以读出接点数是8的倍数。  
Read I/O signal. I/O data is output every 8 points, so the number of read contacts is a multiple of 8.

**命名空间：**
 [HslCommunication.Robot.YASKAWA](3d67023f-78aa-c5d9-85b6-f18f5dfbc176.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public OperateResult<bool[]> IORead(
	int address,
	int length
)
```

```
Public Function IORead ( 
	address As Integer,
	length As Integer
) As OperateResult(Of Boolean())
```

```
public:
OperateResult<array<bool>^>^ IORead(
	int address, 
	int length
)
```

```
member IORead : 
        address : int * 
        length : int -> OperateResult<bool[]> 
```

#### 参数

address
:   类型：SystemInt32  
    读出开始点编号

length
:   类型：SystemInt32  
    读出的接点数

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Boolean  
读取的结果点位信息

![](../icons/SectionExpanded.png)参见

#### 引用

[YRC1000TcpNet 类](cdcdc1cc-17cf-da27-7d20-ca6d6db1cd7e.htm)

[HslCommunication.Robot.YASKAWA 命名空间](3d67023f-78aa-c5d9-85b6-f18f5dfbc176.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## IOReadAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/5d550dc5-c02c-7245-122e-c3624e8d9c14.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.YASKAWA](../html/3d67023f-78aa-c5d9-85b6-f18f5dfbc176.htm "HslCommunication.Robot.YASKAWA")

[YRC1000TcpNet 类](../html/cdcdc1cc-17cf-da27-7d20-ca6d6db1cd7e.htm "YRC1000TcpNet 类")

[YRC1000TcpNet 方法](../html/7a26c354-eee0-b601-1d8a-51198e40e3fa.htm "YRC1000TcpNet 方法")

[Cancel 方法](../html/be66b6de-fe15-95b5-6969-70518f086388.htm "Cancel 方法 ")

[CancelAsync 方法](../html/fedb71cd-919c-c002-7556-551ae30fb10e.htm "CancelAsync 方法 ")

[Cycle 方法](../html/5a144fc7-5bc5-7250-6c13-826c172bd2ce.htm "Cycle 方法 ")

[CycleAsync 方法](../html/82bfe8a9-00b3-b34e-2a82-6d022becdf9a.htm "CycleAsync 方法 ")

[Delete 方法](../html/c7bfb022-46f4-8cdb-1b50-ffc1e5eb716b.htm "Delete 方法 ")

[DeleteAsync 方法](../html/c5a62bdd-3432-abcd-9d03-1ff7e30bd447.htm "DeleteAsync 方法 ")

[GetNewNetMessage 方法](../html/888815ce-da5c-0cb3-c6f9-800ec37a19ab.htm "GetNewNetMessage 方法 ")

[HLock 方法](../html/9c4ed34f-89fd-62e2-f707-08d1348c7541.htm "HLock 方法 ")

[HLockAsync 方法](../html/f0eb6982-a92b-4481-6435-33fa2b841b3c.htm "HLockAsync 方法 ")

[Hold 方法](../html/bf4ccc0a-1576-83ca-49e3-fde86b3a3738.htm "Hold 方法 ")

[HoldAsync 方法](../html/8c84a7d8-3677-0b6b-bdb5-44aa61ed1fcd.htm "HoldAsync 方法 ")

[InitializationOnConnect 方法](../html/5bd0274d-b3e7-f797-2dc0-c0dd2f330694.htm "InitializationOnConnect 方法 ")

[InitializationOnConnectAsync 方法](../html/e2f434c6-2ad9-30e8-8f2c-dcab95f5ec2b.htm "InitializationOnConnectAsync 方法 ")

[IORead 方法](../html/4539f381-2df3-4315-084c-7fcc81d523cb.htm "IORead 方法 ")

[IOReadAsync 方法](../html/5d550dc5-c02c-7245-122e-c3624e8d9c14.htm "IOReadAsync 方法 ")

[IOWrite 方法](../html/0dfecba5-f497-4b0f-55a5-339c7202d742.htm "IOWrite 方法 ")

[IOWriteAsync 方法](../html/64b5c202-29ff-0737-9c62-a4735bc00911.htm "IOWriteAsync 方法 ")

[JSeq 方法](../html/7057228b-aa47-3318-248d-1594588783bd.htm "JSeq 方法 ")

[JSeqAsync 方法](../html/c4eb7232-63e1-cd10-f89e-7f58072e3e9c.htm "JSeqAsync 方法 ")

[Mode 方法](../html/cabcebca-60ed-7689-e44a-e227936dcc51.htm "Mode 方法 ")

[ModeAsync 方法](../html/e0956a5d-5ae1-f884-52dd-0e91b85828b5.htm "ModeAsync 方法 ")

[MoveJ 方法](../html/10990eab-b2f9-e4b8-3c88-5bffc78aa34a.htm "MoveJ 方法 ")

[MoveJAsync 方法](../html/5b05904a-27a0-b4f9-adc7-1fe2cc5b14a1.htm "MoveJAsync 方法 ")

[MSDP 方法](../html/1e78f8b5-7336-7fdf-a7e0-e40af1d3de7f.htm "MSDP 方法 ")

[MSDPAsync 方法](../html/c2a18a60-0dae-112b-af16-0e028e6a203a.htm "MSDPAsync 方法 ")

[Read 方法](../html/69082cc9-c365-8fc4-36df-082c0c6c22c3.htm "Read 方法 ")

[ReadALARM 方法](../html/b56d10b7-3ebc-effd-c8f4-1dbf2d9f35a0.htm "ReadALARM 方法 ")

[ReadALARMAsync 方法](../html/6bd50e3f-6923-2c34-0e8d-d08ba82840f1.htm "ReadALARMAsync 方法 ")

[ReadAsync 方法](../html/185f22c2-a4f1-985e-dd8e-4bf7107c5ae7.htm "ReadAsync 方法 ")

[ReadByCommand 方法](../html/80dbd712-4782-b205-7e73-977c37f5a214.htm "ReadByCommand 方法 ")

[ReadByCommandAsync 方法](../html/2ac7651a-e36e-1635-0a61-b855722693a1.htm "ReadByCommandAsync 方法 ")

[ReadByteVariable 方法](../html/e22f2610-10f3-8472-92c2-3f9e904ebb2b.htm "ReadByteVariable 方法 ")

[ReadByteVariableAsync 方法](../html/d01b8a15-6079-9c57-6936-e0c569cf5194.htm "ReadByteVariableAsync 方法 ")

[ReadDoubleIntegerVariable 方法](../html/beab7b65-77db-3c62-6f95-ccf796ecac03.htm "ReadDoubleIntegerVariable 方法 ")

[ReadDoubleIntegerVariableAsync 方法](../html/c6845ca5-1d25-d613-89f9-28da2515cf24.htm "ReadDoubleIntegerVariableAsync 方法 ")

[ReadFromCoreServer 方法](../html/4e4be36a-8ce8-a41b-6ec4-4f82479c6f23.htm "ReadFromCoreServer 方法 ")

[ReadFromCoreServerAsync 方法](../html/e4cb7ebe-d42a-ea17-8972-b246a42e9a90.htm "ReadFromCoreServerAsync 方法 ")

[ReadIntegerVariable 方法](../html/df3171c4-e18b-a52c-9032-bd789e7388a0.htm "ReadIntegerVariable 方法 ")

[ReadIntegerVariableAsync 方法](../html/954c86a5-37e1-f32b-3a6e-7fc4274cbb6f.htm "ReadIntegerVariableAsync 方法 ")

[ReadJSeq 方法](../html/d7a72f10-2923-e045-8c71-e8b43c551c5a.htm "ReadJSeq 方法 ")

[ReadJSeqAsync 方法](../html/ddfc3321-88c7-9d2a-b71e-9315ce77864d.htm "ReadJSeqAsync 方法 ")

[ReadPOSC 方法](../html/fea57b6b-1d2f-0547-73de-6120ee5de705.htm "ReadPOSC 方法 ")

[ReadPOSCAsync 方法](../html/5f5731c4-a35a-affa-3d76-7b5f50fe63a1.htm "ReadPOSCAsync 方法 ")

[ReadPOSJ 方法](../html/14045444-1d70-3bdf-9848-a886c204da46.htm "ReadPOSJ 方法 ")

[ReadPOSJAsync 方法](../html/0ed6f125-5c82-6a30-1a4e-a6a773dda5d8.htm "ReadPOSJAsync 方法 ")

[ReadRealVariable 方法](../html/5f6c460a-0bb9-32fb-8aa8-33b99dc8be5e.htm "ReadRealVariable 方法 ")

[ReadRealVariableAsync 方法](../html/996ae64e-84f9-ccf6-fa2a-de4e4a1b13c9.htm "ReadRealVariableAsync 方法 ")

[ReadStats 方法](../html/3a24171c-f204-8c63-cdbc-8cf23a965b75.htm "ReadStats 方法 ")

[ReadStatsAsync 方法](../html/43a75b5c-c5af-e81f-286c-a4aa2a7a898b.htm "ReadStatsAsync 方法 ")

[ReadString 方法](../html/dbc09266-42e8-b939-aea5-540d528d19ee.htm "ReadString 方法 ")

[ReadStringAsync 方法](../html/8fa97613-eb8d-c1e1-33f6-925dcf1a1ac0.htm "ReadStringAsync 方法 ")

[ReadStringVariable 方法](../html/5dc9095d-86e8-d434-0027-83797a1d8c0a.htm "ReadStringVariable 方法 ")

[ReadStringVariableAsync 方法](../html/f850e762-bf8a-6ee7-c18c-3496d70244fe.htm "ReadStringVariableAsync 方法 ")

[ReadUFrame 方法](../html/100a1955-fc0a-871a-98f9-bb361c422a3f.htm "ReadUFrame 方法 ")

[ReadUFrameAsync 方法](../html/f63d72f5-eaf4-ae5d-e80a-d881f93a878e.htm "ReadUFrameAsync 方法 ")

[Reset 方法](../html/bf61707b-875c-4b42-2732-9b2fb69535ac.htm "Reset 方法 ")

[ResetAsync 方法](../html/f135af7f-0250-eb16-48fe-959dd0aeeb72.htm "ResetAsync 方法 ")

[SetMJ 方法](../html/9cd8c255-f908-6f4e-da07-33edc743dbf8.htm "SetMJ 方法 ")

[SetMJAsync 方法](../html/74b1cb5d-83f7-61b7-8eb0-0e10a675b964.htm "SetMJAsync 方法 ")

[Start 方法](../html/39764e6e-fde9-c99b-9034-99a64232aef6.htm "Start 方法 ")

[StartAsync 方法](../html/ddefed38-f8b8-4ae1-73e3-eecbfae54043.htm "StartAsync 方法 ")

[Svon 方法](../html/1f83c8cb-fea3-d99f-c7ad-f12f39bd9dee.htm "Svon 方法 ")

[SvonAsync 方法](../html/9b51697e-6164-280a-d37c-1cc16f467b2d.htm "SvonAsync 方法 ")

[ToString 方法](../html/1a5b666e-cae2-1664-908f-56866da693e5.htm "ToString 方法 ")

[Write 方法](../html/27720779-dc92-a6b4-c792-940ec2a10d22.htm "Write 方法 ")

[WriteAsync 方法](../html/e1646488-a935-0a5b-ec6a-4d04a3e69d38.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| YRC1000TcpNetIOReadAsync 方法 |

读取I/O 信号。 I/O 数据是每8个点输出，所以读出接点数是8的倍数。  
Read I/O signal. I/O data is output every 8 points, so the number of read contacts is a multiple of 8.

**命名空间：**
 [HslCommunication.Robot.YASKAWA](3d67023f-78aa-c5d9-85b6-f18f5dfbc176.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public Task<OperateResult<bool[]>> IOReadAsync(
	int address,
	int length
)
```

```
Public Function IOReadAsync ( 
	address As Integer,
	length As Integer
) As Task(Of OperateResult(Of Boolean()))
```

```
public:
Task<OperateResult<array<bool>^>^>^ IOReadAsync(
	int address, 
	int length
)
```

```
member IOReadAsync : 
        address : int * 
        length : int -> Task<OperateResult<bool[]>> 
```

#### 参数

address
:   类型：SystemInt32  
    读出开始点编号

length
:   类型：SystemInt32  
    读出的接点数

#### 返回值

类型：Task[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Boolean  
读取的结果点位信息

![](../icons/SectionExpanded.png)参见

#### 引用

[YRC1000TcpNet 类](cdcdc1cc-17cf-da27-7d20-ca6d6db1cd7e.htm)

[HslCommunication.Robot.YASKAWA 命名空间](3d67023f-78aa-c5d9-85b6-f18f5dfbc176.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## IOWrite 方法 

[原文連結](http://api.hslcommunication.cn/html/0dfecba5-f497-4b0f-55a5-339c7202d742.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.YASKAWA](../html/3d67023f-78aa-c5d9-85b6-f18f5dfbc176.htm "HslCommunication.Robot.YASKAWA")

[YRC1000TcpNet 类](../html/cdcdc1cc-17cf-da27-7d20-ca6d6db1cd7e.htm "YRC1000TcpNet 类")

[YRC1000TcpNet 方法](../html/7a26c354-eee0-b601-1d8a-51198e40e3fa.htm "YRC1000TcpNet 方法")

[Cancel 方法](../html/be66b6de-fe15-95b5-6969-70518f086388.htm "Cancel 方法 ")

[CancelAsync 方法](../html/fedb71cd-919c-c002-7556-551ae30fb10e.htm "CancelAsync 方法 ")

[Cycle 方法](../html/5a144fc7-5bc5-7250-6c13-826c172bd2ce.htm "Cycle 方法 ")

[CycleAsync 方法](../html/82bfe8a9-00b3-b34e-2a82-6d022becdf9a.htm "CycleAsync 方法 ")

[Delete 方法](../html/c7bfb022-46f4-8cdb-1b50-ffc1e5eb716b.htm "Delete 方法 ")

[DeleteAsync 方法](../html/c5a62bdd-3432-abcd-9d03-1ff7e30bd447.htm "DeleteAsync 方法 ")

[GetNewNetMessage 方法](../html/888815ce-da5c-0cb3-c6f9-800ec37a19ab.htm "GetNewNetMessage 方法 ")

[HLock 方法](../html/9c4ed34f-89fd-62e2-f707-08d1348c7541.htm "HLock 方法 ")

[HLockAsync 方法](../html/f0eb6982-a92b-4481-6435-33fa2b841b3c.htm "HLockAsync 方法 ")

[Hold 方法](../html/bf4ccc0a-1576-83ca-49e3-fde86b3a3738.htm "Hold 方法 ")

[HoldAsync 方法](../html/8c84a7d8-3677-0b6b-bdb5-44aa61ed1fcd.htm "HoldAsync 方法 ")

[InitializationOnConnect 方法](../html/5bd0274d-b3e7-f797-2dc0-c0dd2f330694.htm "InitializationOnConnect 方法 ")

[InitializationOnConnectAsync 方法](../html/e2f434c6-2ad9-30e8-8f2c-dcab95f5ec2b.htm "InitializationOnConnectAsync 方法 ")

[IORead 方法](../html/4539f381-2df3-4315-084c-7fcc81d523cb.htm "IORead 方法 ")

[IOReadAsync 方法](../html/5d550dc5-c02c-7245-122e-c3624e8d9c14.htm "IOReadAsync 方法 ")

[IOWrite 方法](../html/0dfecba5-f497-4b0f-55a5-339c7202d742.htm "IOWrite 方法 ")

[IOWriteAsync 方法](../html/64b5c202-29ff-0737-9c62-a4735bc00911.htm "IOWriteAsync 方法 ")

[JSeq 方法](../html/7057228b-aa47-3318-248d-1594588783bd.htm "JSeq 方法 ")

[JSeqAsync 方法](../html/c4eb7232-63e1-cd10-f89e-7f58072e3e9c.htm "JSeqAsync 方法 ")

[Mode 方法](../html/cabcebca-60ed-7689-e44a-e227936dcc51.htm "Mode 方法 ")

[ModeAsync 方法](../html/e0956a5d-5ae1-f884-52dd-0e91b85828b5.htm "ModeAsync 方法 ")

[MoveJ 方法](../html/10990eab-b2f9-e4b8-3c88-5bffc78aa34a.htm "MoveJ 方法 ")

[MoveJAsync 方法](../html/5b05904a-27a0-b4f9-adc7-1fe2cc5b14a1.htm "MoveJAsync 方法 ")

[MSDP 方法](../html/1e78f8b5-7336-7fdf-a7e0-e40af1d3de7f.htm "MSDP 方法 ")

[MSDPAsync 方法](../html/c2a18a60-0dae-112b-af16-0e028e6a203a.htm "MSDPAsync 方法 ")

[Read 方法](../html/69082cc9-c365-8fc4-36df-082c0c6c22c3.htm "Read 方法 ")

[ReadALARM 方法](../html/b56d10b7-3ebc-effd-c8f4-1dbf2d9f35a0.htm "ReadALARM 方法 ")

[ReadALARMAsync 方法](../html/6bd50e3f-6923-2c34-0e8d-d08ba82840f1.htm "ReadALARMAsync 方法 ")

[ReadAsync 方法](../html/185f22c2-a4f1-985e-dd8e-4bf7107c5ae7.htm "ReadAsync 方法 ")

[ReadByCommand 方法](../html/80dbd712-4782-b205-7e73-977c37f5a214.htm "ReadByCommand 方法 ")

[ReadByCommandAsync 方法](../html/2ac7651a-e36e-1635-0a61-b855722693a1.htm "ReadByCommandAsync 方法 ")

[ReadByteVariable 方法](../html/e22f2610-10f3-8472-92c2-3f9e904ebb2b.htm "ReadByteVariable 方法 ")

[ReadByteVariableAsync 方法](../html/d01b8a15-6079-9c57-6936-e0c569cf5194.htm "ReadByteVariableAsync 方法 ")

[ReadDoubleIntegerVariable 方法](../html/beab7b65-77db-3c62-6f95-ccf796ecac03.htm "ReadDoubleIntegerVariable 方法 ")

[ReadDoubleIntegerVariableAsync 方法](../html/c6845ca5-1d25-d613-89f9-28da2515cf24.htm "ReadDoubleIntegerVariableAsync 方法 ")

[ReadFromCoreServer 方法](../html/4e4be36a-8ce8-a41b-6ec4-4f82479c6f23.htm "ReadFromCoreServer 方法 ")

[ReadFromCoreServerAsync 方法](../html/e4cb7ebe-d42a-ea17-8972-b246a42e9a90.htm "ReadFromCoreServerAsync 方法 ")

[ReadIntegerVariable 方法](../html/df3171c4-e18b-a52c-9032-bd789e7388a0.htm "ReadIntegerVariable 方法 ")

[ReadIntegerVariableAsync 方法](../html/954c86a5-37e1-f32b-3a6e-7fc4274cbb6f.htm "ReadIntegerVariableAsync 方法 ")

[ReadJSeq 方法](../html/d7a72f10-2923-e045-8c71-e8b43c551c5a.htm "ReadJSeq 方法 ")

[ReadJSeqAsync 方法](../html/ddfc3321-88c7-9d2a-b71e-9315ce77864d.htm "ReadJSeqAsync 方法 ")

[ReadPOSC 方法](../html/fea57b6b-1d2f-0547-73de-6120ee5de705.htm "ReadPOSC 方法 ")

[ReadPOSCAsync 方法](../html/5f5731c4-a35a-affa-3d76-7b5f50fe63a1.htm "ReadPOSCAsync 方法 ")

[ReadPOSJ 方法](../html/14045444-1d70-3bdf-9848-a886c204da46.htm "ReadPOSJ 方法 ")

[ReadPOSJAsync 方法](../html/0ed6f125-5c82-6a30-1a4e-a6a773dda5d8.htm "ReadPOSJAsync 方法 ")

[ReadRealVariable 方法](../html/5f6c460a-0bb9-32fb-8aa8-33b99dc8be5e.htm "ReadRealVariable 方法 ")

[ReadRealVariableAsync 方法](../html/996ae64e-84f9-ccf6-fa2a-de4e4a1b13c9.htm "ReadRealVariableAsync 方法 ")

[ReadStats 方法](../html/3a24171c-f204-8c63-cdbc-8cf23a965b75.htm "ReadStats 方法 ")

[ReadStatsAsync 方法](../html/43a75b5c-c5af-e81f-286c-a4aa2a7a898b.htm "ReadStatsAsync 方法 ")

[ReadString 方法](../html/dbc09266-42e8-b939-aea5-540d528d19ee.htm "ReadString 方法 ")

[ReadStringAsync 方法](../html/8fa97613-eb8d-c1e1-33f6-925dcf1a1ac0.htm "ReadStringAsync 方法 ")

[ReadStringVariable 方法](../html/5dc9095d-86e8-d434-0027-83797a1d8c0a.htm "ReadStringVariable 方法 ")

[ReadStringVariableAsync 方法](../html/f850e762-bf8a-6ee7-c18c-3496d70244fe.htm "ReadStringVariableAsync 方法 ")

[ReadUFrame 方法](../html/100a1955-fc0a-871a-98f9-bb361c422a3f.htm "ReadUFrame 方法 ")

[ReadUFrameAsync 方法](../html/f63d72f5-eaf4-ae5d-e80a-d881f93a878e.htm "ReadUFrameAsync 方法 ")

[Reset 方法](../html/bf61707b-875c-4b42-2732-9b2fb69535ac.htm "Reset 方法 ")

[ResetAsync 方法](../html/f135af7f-0250-eb16-48fe-959dd0aeeb72.htm "ResetAsync 方法 ")

[SetMJ 方法](../html/9cd8c255-f908-6f4e-da07-33edc743dbf8.htm "SetMJ 方法 ")

[SetMJAsync 方法](../html/74b1cb5d-83f7-61b7-8eb0-0e10a675b964.htm "SetMJAsync 方法 ")

[Start 方法](../html/39764e6e-fde9-c99b-9034-99a64232aef6.htm "Start 方法 ")

[StartAsync 方法](../html/ddefed38-f8b8-4ae1-73e3-eecbfae54043.htm "StartAsync 方法 ")

[Svon 方法](../html/1f83c8cb-fea3-d99f-c7ad-f12f39bd9dee.htm "Svon 方法 ")

[SvonAsync 方法](../html/9b51697e-6164-280a-d37c-1cc16f467b2d.htm "SvonAsync 方法 ")

[ToString 方法](../html/1a5b666e-cae2-1664-908f-56866da693e5.htm "ToString 方法 ")

[Write 方法](../html/27720779-dc92-a6b4-c792-940ec2a10d22.htm "Write 方法 ")

[WriteAsync 方法](../html/e1646488-a935-0a5b-ec6a-4d04a3e69d38.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| YRC1000TcpNetIOWrite 方法 |

写入I/O信号状态，写入接点数请指定8的倍数。IO 信号的网络写入仅可是（ #27010 ～ #29567）。  
To write I/O signal status, please specify a multiple of 8 for the number of write contacts.
The network write of IO signal is only available (#27010 to #29567).

**命名空间：**
 [HslCommunication.Robot.YASKAWA](3d67023f-78aa-c5d9-85b6-f18f5dfbc176.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public OperateResult IOWrite(
	int address,
	bool[] value
)
```

```
Public Function IOWrite ( 
	address As Integer,
	value As Boolean()
) As OperateResult
```

```
public:
OperateResult^ IOWrite(
	int address, 
	array<bool>^ value
)
```

```
member IOWrite : 
        address : int * 
        value : bool[] -> OperateResult 
```

#### 参数

address
:   类型：SystemInt32  
    写入开始接点编号

value
:   类型：SystemBoolean  
    写入的bool值，写入接点数请指定8的倍数。

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否写入成功

![](../icons/SectionExpanded.png)参见

#### 引用

[YRC1000TcpNet 类](cdcdc1cc-17cf-da27-7d20-ca6d6db1cd7e.htm)

[HslCommunication.Robot.YASKAWA 命名空间](3d67023f-78aa-c5d9-85b6-f18f5dfbc176.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## IOWriteAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/64b5c202-29ff-0737-9c62-a4735bc00911.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.YASKAWA](../html/3d67023f-78aa-c5d9-85b6-f18f5dfbc176.htm "HslCommunication.Robot.YASKAWA")

[YRC1000TcpNet 类](../html/cdcdc1cc-17cf-da27-7d20-ca6d6db1cd7e.htm "YRC1000TcpNet 类")

[YRC1000TcpNet 方法](../html/7a26c354-eee0-b601-1d8a-51198e40e3fa.htm "YRC1000TcpNet 方法")

[Cancel 方法](../html/be66b6de-fe15-95b5-6969-70518f086388.htm "Cancel 方法 ")

[CancelAsync 方法](../html/fedb71cd-919c-c002-7556-551ae30fb10e.htm "CancelAsync 方法 ")

[Cycle 方法](../html/5a144fc7-5bc5-7250-6c13-826c172bd2ce.htm "Cycle 方法 ")

[CycleAsync 方法](../html/82bfe8a9-00b3-b34e-2a82-6d022becdf9a.htm "CycleAsync 方法 ")

[Delete 方法](../html/c7bfb022-46f4-8cdb-1b50-ffc1e5eb716b.htm "Delete 方法 ")

[DeleteAsync 方法](../html/c5a62bdd-3432-abcd-9d03-1ff7e30bd447.htm "DeleteAsync 方法 ")

[GetNewNetMessage 方法](../html/888815ce-da5c-0cb3-c6f9-800ec37a19ab.htm "GetNewNetMessage 方法 ")

[HLock 方法](../html/9c4ed34f-89fd-62e2-f707-08d1348c7541.htm "HLock 方法 ")

[HLockAsync 方法](../html/f0eb6982-a92b-4481-6435-33fa2b841b3c.htm "HLockAsync 方法 ")

[Hold 方法](../html/bf4ccc0a-1576-83ca-49e3-fde86b3a3738.htm "Hold 方法 ")

[HoldAsync 方法](../html/8c84a7d8-3677-0b6b-bdb5-44aa61ed1fcd.htm "HoldAsync 方法 ")

[InitializationOnConnect 方法](../html/5bd0274d-b3e7-f797-2dc0-c0dd2f330694.htm "InitializationOnConnect 方法 ")

[InitializationOnConnectAsync 方法](../html/e2f434c6-2ad9-30e8-8f2c-dcab95f5ec2b.htm "InitializationOnConnectAsync 方法 ")

[IORead 方法](../html/4539f381-2df3-4315-084c-7fcc81d523cb.htm "IORead 方法 ")

[IOReadAsync 方法](../html/5d550dc5-c02c-7245-122e-c3624e8d9c14.htm "IOReadAsync 方法 ")

[IOWrite 方法](../html/0dfecba5-f497-4b0f-55a5-339c7202d742.htm "IOWrite 方法 ")

[IOWriteAsync 方法](../html/64b5c202-29ff-0737-9c62-a4735bc00911.htm "IOWriteAsync 方法 ")

[JSeq 方法](../html/7057228b-aa47-3318-248d-1594588783bd.htm "JSeq 方法 ")

[JSeqAsync 方法](../html/c4eb7232-63e1-cd10-f89e-7f58072e3e9c.htm "JSeqAsync 方法 ")

[Mode 方法](../html/cabcebca-60ed-7689-e44a-e227936dcc51.htm "Mode 方法 ")

[ModeAsync 方法](../html/e0956a5d-5ae1-f884-52dd-0e91b85828b5.htm "ModeAsync 方法 ")

[MoveJ 方法](../html/10990eab-b2f9-e4b8-3c88-5bffc78aa34a.htm "MoveJ 方法 ")

[MoveJAsync 方法](../html/5b05904a-27a0-b4f9-adc7-1fe2cc5b14a1.htm "MoveJAsync 方法 ")

[MSDP 方法](../html/1e78f8b5-7336-7fdf-a7e0-e40af1d3de7f.htm "MSDP 方法 ")

[MSDPAsync 方法](../html/c2a18a60-0dae-112b-af16-0e028e6a203a.htm "MSDPAsync 方法 ")

[Read 方法](../html/69082cc9-c365-8fc4-36df-082c0c6c22c3.htm "Read 方法 ")

[ReadALARM 方法](../html/b56d10b7-3ebc-effd-c8f4-1dbf2d9f35a0.htm "ReadALARM 方法 ")

[ReadALARMAsync 方法](../html/6bd50e3f-6923-2c34-0e8d-d08ba82840f1.htm "ReadALARMAsync 方法 ")

[ReadAsync 方法](../html/185f22c2-a4f1-985e-dd8e-4bf7107c5ae7.htm "ReadAsync 方法 ")

[ReadByCommand 方法](../html/80dbd712-4782-b205-7e73-977c37f5a214.htm "ReadByCommand 方法 ")

[ReadByCommandAsync 方法](../html/2ac7651a-e36e-1635-0a61-b855722693a1.htm "ReadByCommandAsync 方法 ")

[ReadByteVariable 方法](../html/e22f2610-10f3-8472-92c2-3f9e904ebb2b.htm "ReadByteVariable 方法 ")

[ReadByteVariableAsync 方法](../html/d01b8a15-6079-9c57-6936-e0c569cf5194.htm "ReadByteVariableAsync 方法 ")

[ReadDoubleIntegerVariable 方法](../html/beab7b65-77db-3c62-6f95-ccf796ecac03.htm "ReadDoubleIntegerVariable 方法 ")

[ReadDoubleIntegerVariableAsync 方法](../html/c6845ca5-1d25-d613-89f9-28da2515cf24.htm "ReadDoubleIntegerVariableAsync 方法 ")

[ReadFromCoreServer 方法](../html/4e4be36a-8ce8-a41b-6ec4-4f82479c6f23.htm "ReadFromCoreServer 方法 ")

[ReadFromCoreServerAsync 方法](../html/e4cb7ebe-d42a-ea17-8972-b246a42e9a90.htm "ReadFromCoreServerAsync 方法 ")

[ReadIntegerVariable 方法](../html/df3171c4-e18b-a52c-9032-bd789e7388a0.htm "ReadIntegerVariable 方法 ")

[ReadIntegerVariableAsync 方法](../html/954c86a5-37e1-f32b-3a6e-7fc4274cbb6f.htm "ReadIntegerVariableAsync 方法 ")

[ReadJSeq 方法](../html/d7a72f10-2923-e045-8c71-e8b43c551c5a.htm "ReadJSeq 方法 ")

[ReadJSeqAsync 方法](../html/ddfc3321-88c7-9d2a-b71e-9315ce77864d.htm "ReadJSeqAsync 方法 ")

[ReadPOSC 方法](../html/fea57b6b-1d2f-0547-73de-6120ee5de705.htm "ReadPOSC 方法 ")

[ReadPOSCAsync 方法](../html/5f5731c4-a35a-affa-3d76-7b5f50fe63a1.htm "ReadPOSCAsync 方法 ")

[ReadPOSJ 方法](../html/14045444-1d70-3bdf-9848-a886c204da46.htm "ReadPOSJ 方法 ")

[ReadPOSJAsync 方法](../html/0ed6f125-5c82-6a30-1a4e-a6a773dda5d8.htm "ReadPOSJAsync 方法 ")

[ReadRealVariable 方法](../html/5f6c460a-0bb9-32fb-8aa8-33b99dc8be5e.htm "ReadRealVariable 方法 ")

[ReadRealVariableAsync 方法](../html/996ae64e-84f9-ccf6-fa2a-de4e4a1b13c9.htm "ReadRealVariableAsync 方法 ")

[ReadStats 方法](../html/3a24171c-f204-8c63-cdbc-8cf23a965b75.htm "ReadStats 方法 ")

[ReadStatsAsync 方法](../html/43a75b5c-c5af-e81f-286c-a4aa2a7a898b.htm "ReadStatsAsync 方法 ")

[ReadString 方法](../html/dbc09266-42e8-b939-aea5-540d528d19ee.htm "ReadString 方法 ")

[ReadStringAsync 方法](../html/8fa97613-eb8d-c1e1-33f6-925dcf1a1ac0.htm "ReadStringAsync 方法 ")

[ReadStringVariable 方法](../html/5dc9095d-86e8-d434-0027-83797a1d8c0a.htm "ReadStringVariable 方法 ")

[ReadStringVariableAsync 方法](../html/f850e762-bf8a-6ee7-c18c-3496d70244fe.htm "ReadStringVariableAsync 方法 ")

[ReadUFrame 方法](../html/100a1955-fc0a-871a-98f9-bb361c422a3f.htm "ReadUFrame 方法 ")

[ReadUFrameAsync 方法](../html/f63d72f5-eaf4-ae5d-e80a-d881f93a878e.htm "ReadUFrameAsync 方法 ")

[Reset 方法](../html/bf61707b-875c-4b42-2732-9b2fb69535ac.htm "Reset 方法 ")

[ResetAsync 方法](../html/f135af7f-0250-eb16-48fe-959dd0aeeb72.htm "ResetAsync 方法 ")

[SetMJ 方法](../html/9cd8c255-f908-6f4e-da07-33edc743dbf8.htm "SetMJ 方法 ")

[SetMJAsync 方法](../html/74b1cb5d-83f7-61b7-8eb0-0e10a675b964.htm "SetMJAsync 方法 ")

[Start 方法](../html/39764e6e-fde9-c99b-9034-99a64232aef6.htm "Start 方法 ")

[StartAsync 方法](../html/ddefed38-f8b8-4ae1-73e3-eecbfae54043.htm "StartAsync 方法 ")

[Svon 方法](../html/1f83c8cb-fea3-d99f-c7ad-f12f39bd9dee.htm "Svon 方法 ")

[SvonAsync 方法](../html/9b51697e-6164-280a-d37c-1cc16f467b2d.htm "SvonAsync 方法 ")

[ToString 方法](../html/1a5b666e-cae2-1664-908f-56866da693e5.htm "ToString 方法 ")

[Write 方法](../html/27720779-dc92-a6b4-c792-940ec2a10d22.htm "Write 方法 ")

[WriteAsync 方法](../html/e1646488-a935-0a5b-ec6a-4d04a3e69d38.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| YRC1000TcpNetIOWriteAsync 方法 |

写入I/O信号状态，写入接点数请指定8的倍数。IO 信号的网络写入仅可是（ #27010 ～ #29567）。  
To write I/O signal status, please specify a multiple of 8 for the number of write contacts.
The network write of IO signal is only available (#27010 to #29567).

**命名空间：**
 [HslCommunication.Robot.YASKAWA](3d67023f-78aa-c5d9-85b6-f18f5dfbc176.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public Task<OperateResult> IOWriteAsync(
	int address,
	bool[] value
)
```

```
Public Function IOWriteAsync ( 
	address As Integer,
	value As Boolean()
) As Task(Of OperateResult)
```

```
public:
Task<OperateResult^>^ IOWriteAsync(
	int address, 
	array<bool>^ value
)
```

```
member IOWriteAsync : 
        address : int * 
        value : bool[] -> Task<OperateResult> 
```

#### 参数

address
:   类型：SystemInt32  
    写入开始接点编号

value
:   类型：SystemBoolean  
    写入的bool值，写入接点数请指定8的倍数。

#### 返回值

类型：Task[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否写入成功

![](../icons/SectionExpanded.png)参见

#### 引用

[YRC1000TcpNet 类](cdcdc1cc-17cf-da27-7d20-ca6d6db1cd7e.htm)

[HslCommunication.Robot.YASKAWA 命名空间](3d67023f-78aa-c5d9-85b6-f18f5dfbc176.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## JSeq 方法 

[原文連結](http://api.hslcommunication.cn/html/7057228b-aa47-3318-248d-1594588783bd.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.YASKAWA](../html/3d67023f-78aa-c5d9-85b6-f18f5dfbc176.htm "HslCommunication.Robot.YASKAWA")

[YRC1000TcpNet 类](../html/cdcdc1cc-17cf-da27-7d20-ca6d6db1cd7e.htm "YRC1000TcpNet 类")

[YRC1000TcpNet 方法](../html/7a26c354-eee0-b601-1d8a-51198e40e3fa.htm "YRC1000TcpNet 方法")

[Cancel 方法](../html/be66b6de-fe15-95b5-6969-70518f086388.htm "Cancel 方法 ")

[CancelAsync 方法](../html/fedb71cd-919c-c002-7556-551ae30fb10e.htm "CancelAsync 方法 ")

[Cycle 方法](../html/5a144fc7-5bc5-7250-6c13-826c172bd2ce.htm "Cycle 方法 ")

[CycleAsync 方法](../html/82bfe8a9-00b3-b34e-2a82-6d022becdf9a.htm "CycleAsync 方法 ")

[Delete 方法](../html/c7bfb022-46f4-8cdb-1b50-ffc1e5eb716b.htm "Delete 方法 ")

[DeleteAsync 方法](../html/c5a62bdd-3432-abcd-9d03-1ff7e30bd447.htm "DeleteAsync 方法 ")

[GetNewNetMessage 方法](../html/888815ce-da5c-0cb3-c6f9-800ec37a19ab.htm "GetNewNetMessage 方法 ")

[HLock 方法](../html/9c4ed34f-89fd-62e2-f707-08d1348c7541.htm "HLock 方法 ")

[HLockAsync 方法](../html/f0eb6982-a92b-4481-6435-33fa2b841b3c.htm "HLockAsync 方法 ")

[Hold 方法](../html/bf4ccc0a-1576-83ca-49e3-fde86b3a3738.htm "Hold 方法 ")

[HoldAsync 方法](../html/8c84a7d8-3677-0b6b-bdb5-44aa61ed1fcd.htm "HoldAsync 方法 ")

[InitializationOnConnect 方法](../html/5bd0274d-b3e7-f797-2dc0-c0dd2f330694.htm "InitializationOnConnect 方法 ")

[InitializationOnConnectAsync 方法](../html/e2f434c6-2ad9-30e8-8f2c-dcab95f5ec2b.htm "InitializationOnConnectAsync 方法 ")

[IORead 方法](../html/4539f381-2df3-4315-084c-7fcc81d523cb.htm "IORead 方法 ")

[IOReadAsync 方法](../html/5d550dc5-c02c-7245-122e-c3624e8d9c14.htm "IOReadAsync 方法 ")

[IOWrite 方法](../html/0dfecba5-f497-4b0f-55a5-339c7202d742.htm "IOWrite 方法 ")

[IOWriteAsync 方法](../html/64b5c202-29ff-0737-9c62-a4735bc00911.htm "IOWriteAsync 方法 ")

[JSeq 方法](../html/7057228b-aa47-3318-248d-1594588783bd.htm "JSeq 方法 ")

[JSeqAsync 方法](../html/c4eb7232-63e1-cd10-f89e-7f58072e3e9c.htm "JSeqAsync 方法 ")

[Mode 方法](../html/cabcebca-60ed-7689-e44a-e227936dcc51.htm "Mode 方法 ")

[ModeAsync 方法](../html/e0956a5d-5ae1-f884-52dd-0e91b85828b5.htm "ModeAsync 方法 ")

[MoveJ 方法](../html/10990eab-b2f9-e4b8-3c88-5bffc78aa34a.htm "MoveJ 方法 ")

[MoveJAsync 方法](../html/5b05904a-27a0-b4f9-adc7-1fe2cc5b14a1.htm "MoveJAsync 方法 ")

[MSDP 方法](../html/1e78f8b5-7336-7fdf-a7e0-e40af1d3de7f.htm "MSDP 方法 ")

[MSDPAsync 方法](../html/c2a18a60-0dae-112b-af16-0e028e6a203a.htm "MSDPAsync 方法 ")

[Read 方法](../html/69082cc9-c365-8fc4-36df-082c0c6c22c3.htm "Read 方法 ")

[ReadALARM 方法](../html/b56d10b7-3ebc-effd-c8f4-1dbf2d9f35a0.htm "ReadALARM 方法 ")

[ReadALARMAsync 方法](../html/6bd50e3f-6923-2c34-0e8d-d08ba82840f1.htm "ReadALARMAsync 方法 ")

[ReadAsync 方法](../html/185f22c2-a4f1-985e-dd8e-4bf7107c5ae7.htm "ReadAsync 方法 ")

[ReadByCommand 方法](../html/80dbd712-4782-b205-7e73-977c37f5a214.htm "ReadByCommand 方法 ")

[ReadByCommandAsync 方法](../html/2ac7651a-e36e-1635-0a61-b855722693a1.htm "ReadByCommandAsync 方法 ")

[ReadByteVariable 方法](../html/e22f2610-10f3-8472-92c2-3f9e904ebb2b.htm "ReadByteVariable 方法 ")

[ReadByteVariableAsync 方法](../html/d01b8a15-6079-9c57-6936-e0c569cf5194.htm "ReadByteVariableAsync 方法 ")

[ReadDoubleIntegerVariable 方法](../html/beab7b65-77db-3c62-6f95-ccf796ecac03.htm "ReadDoubleIntegerVariable 方法 ")

[ReadDoubleIntegerVariableAsync 方法](../html/c6845ca5-1d25-d613-89f9-28da2515cf24.htm "ReadDoubleIntegerVariableAsync 方法 ")

[ReadFromCoreServer 方法](../html/4e4be36a-8ce8-a41b-6ec4-4f82479c6f23.htm "ReadFromCoreServer 方法 ")

[ReadFromCoreServerAsync 方法](../html/e4cb7ebe-d42a-ea17-8972-b246a42e9a90.htm "ReadFromCoreServerAsync 方法 ")

[ReadIntegerVariable 方法](../html/df3171c4-e18b-a52c-9032-bd789e7388a0.htm "ReadIntegerVariable 方法 ")

[ReadIntegerVariableAsync 方法](../html/954c86a5-37e1-f32b-3a6e-7fc4274cbb6f.htm "ReadIntegerVariableAsync 方法 ")

[ReadJSeq 方法](../html/d7a72f10-2923-e045-8c71-e8b43c551c5a.htm "ReadJSeq 方法 ")

[ReadJSeqAsync 方法](../html/ddfc3321-88c7-9d2a-b71e-9315ce77864d.htm "ReadJSeqAsync 方法 ")

[ReadPOSC 方法](../html/fea57b6b-1d2f-0547-73de-6120ee5de705.htm "ReadPOSC 方法 ")

[ReadPOSCAsync 方法](../html/5f5731c4-a35a-affa-3d76-7b5f50fe63a1.htm "ReadPOSCAsync 方法 ")

[ReadPOSJ 方法](../html/14045444-1d70-3bdf-9848-a886c204da46.htm "ReadPOSJ 方法 ")

[ReadPOSJAsync 方法](../html/0ed6f125-5c82-6a30-1a4e-a6a773dda5d8.htm "ReadPOSJAsync 方法 ")

[ReadRealVariable 方法](../html/5f6c460a-0bb9-32fb-8aa8-33b99dc8be5e.htm "ReadRealVariable 方法 ")

[ReadRealVariableAsync 方法](../html/996ae64e-84f9-ccf6-fa2a-de4e4a1b13c9.htm "ReadRealVariableAsync 方法 ")

[ReadStats 方法](../html/3a24171c-f204-8c63-cdbc-8cf23a965b75.htm "ReadStats 方法 ")

[ReadStatsAsync 方法](../html/43a75b5c-c5af-e81f-286c-a4aa2a7a898b.htm "ReadStatsAsync 方法 ")

[ReadString 方法](../html/dbc09266-42e8-b939-aea5-540d528d19ee.htm "ReadString 方法 ")

[ReadStringAsync 方法](../html/8fa97613-eb8d-c1e1-33f6-925dcf1a1ac0.htm "ReadStringAsync 方法 ")

[ReadStringVariable 方法](../html/5dc9095d-86e8-d434-0027-83797a1d8c0a.htm "ReadStringVariable 方法 ")

[ReadStringVariableAsync 方法](../html/f850e762-bf8a-6ee7-c18c-3496d70244fe.htm "ReadStringVariableAsync 方法 ")

[ReadUFrame 方法](../html/100a1955-fc0a-871a-98f9-bb361c422a3f.htm "ReadUFrame 方法 ")

[ReadUFrameAsync 方法](../html/f63d72f5-eaf4-ae5d-e80a-d881f93a878e.htm "ReadUFrameAsync 方法 ")

[Reset 方法](../html/bf61707b-875c-4b42-2732-9b2fb69535ac.htm "Reset 方法 ")

[ResetAsync 方法](../html/f135af7f-0250-eb16-48fe-959dd0aeeb72.htm "ResetAsync 方法 ")

[SetMJ 方法](../html/9cd8c255-f908-6f4e-da07-33edc743dbf8.htm "SetMJ 方法 ")

[SetMJAsync 方法](../html/74b1cb5d-83f7-61b7-8eb0-0e10a675b964.htm "SetMJAsync 方法 ")

[Start 方法](../html/39764e6e-fde9-c99b-9034-99a64232aef6.htm "Start 方法 ")

[StartAsync 方法](../html/ddefed38-f8b8-4ae1-73e3-eecbfae54043.htm "StartAsync 方法 ")

[Svon 方法](../html/1f83c8cb-fea3-d99f-c7ad-f12f39bd9dee.htm "Svon 方法 ")

[SvonAsync 方法](../html/9b51697e-6164-280a-d37c-1cc16f467b2d.htm "SvonAsync 方法 ")

[ToString 方法](../html/1a5b666e-cae2-1664-908f-56866da693e5.htm "ToString 方法 ")

[Write 方法](../html/27720779-dc92-a6b4-c792-940ec2a10d22.htm "Write 方法 ")

[WriteAsync 方法](../html/e1646488-a935-0a5b-ec6a-4d04a3e69d38.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| YRC1000TcpNetJSeq 方法 |

设定执行程序的名称和行编号。  
Set the name and line number of the executed program.

**命名空间：**
 [HslCommunication.Robot.YASKAWA](3d67023f-78aa-c5d9-85b6-f18f5dfbc176.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public OperateResult JSeq(
	string programName,
	int line
)
```

```
Public Function JSeq ( 
	programName As String,
	line As Integer
) As OperateResult
```

```
public:
OperateResult^ JSeq(
	String^ programName, 
	int line
)
```

```
member JSeq : 
        programName : string * 
        line : int -> OperateResult 
```

#### 参数

programName
:   类型：SystemString  
    设定程序名称

line
:   类型：SystemInt32  
    设定行编号（ 0 ～ 9999）

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否设定成功

![](../icons/SectionExpanded.png)参见

#### 引用

[YRC1000TcpNet 类](cdcdc1cc-17cf-da27-7d20-ca6d6db1cd7e.htm)

[HslCommunication.Robot.YASKAWA 命名空间](3d67023f-78aa-c5d9-85b6-f18f5dfbc176.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## JSeqAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/c4eb7232-63e1-cd10-f89e-7f58072e3e9c.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.YASKAWA](../html/3d67023f-78aa-c5d9-85b6-f18f5dfbc176.htm "HslCommunication.Robot.YASKAWA")

[YRC1000TcpNet 类](../html/cdcdc1cc-17cf-da27-7d20-ca6d6db1cd7e.htm "YRC1000TcpNet 类")

[YRC1000TcpNet 方法](../html/7a26c354-eee0-b601-1d8a-51198e40e3fa.htm "YRC1000TcpNet 方法")

[Cancel 方法](../html/be66b6de-fe15-95b5-6969-70518f086388.htm "Cancel 方法 ")

[CancelAsync 方法](../html/fedb71cd-919c-c002-7556-551ae30fb10e.htm "CancelAsync 方法 ")

[Cycle 方法](../html/5a144fc7-5bc5-7250-6c13-826c172bd2ce.htm "Cycle 方法 ")

[CycleAsync 方法](../html/82bfe8a9-00b3-b34e-2a82-6d022becdf9a.htm "CycleAsync 方法 ")

[Delete 方法](../html/c7bfb022-46f4-8cdb-1b50-ffc1e5eb716b.htm "Delete 方法 ")

[DeleteAsync 方法](../html/c5a62bdd-3432-abcd-9d03-1ff7e30bd447.htm "DeleteAsync 方法 ")

[GetNewNetMessage 方法](../html/888815ce-da5c-0cb3-c6f9-800ec37a19ab.htm "GetNewNetMessage 方法 ")

[HLock 方法](../html/9c4ed34f-89fd-62e2-f707-08d1348c7541.htm "HLock 方法 ")

[HLockAsync 方法](../html/f0eb6982-a92b-4481-6435-33fa2b841b3c.htm "HLockAsync 方法 ")

[Hold 方法](../html/bf4ccc0a-1576-83ca-49e3-fde86b3a3738.htm "Hold 方法 ")

[HoldAsync 方法](../html/8c84a7d8-3677-0b6b-bdb5-44aa61ed1fcd.htm "HoldAsync 方法 ")

[InitializationOnConnect 方法](../html/5bd0274d-b3e7-f797-2dc0-c0dd2f330694.htm "InitializationOnConnect 方法 ")

[InitializationOnConnectAsync 方法](../html/e2f434c6-2ad9-30e8-8f2c-dcab95f5ec2b.htm "InitializationOnConnectAsync 方法 ")

[IORead 方法](../html/4539f381-2df3-4315-084c-7fcc81d523cb.htm "IORead 方法 ")

[IOReadAsync 方法](../html/5d550dc5-c02c-7245-122e-c3624e8d9c14.htm "IOReadAsync 方法 ")

[IOWrite 方法](../html/0dfecba5-f497-4b0f-55a5-339c7202d742.htm "IOWrite 方法 ")

[IOWriteAsync 方法](../html/64b5c202-29ff-0737-9c62-a4735bc00911.htm "IOWriteAsync 方法 ")

[JSeq 方法](../html/7057228b-aa47-3318-248d-1594588783bd.htm "JSeq 方法 ")

[JSeqAsync 方法](../html/c4eb7232-63e1-cd10-f89e-7f58072e3e9c.htm "JSeqAsync 方法 ")

[Mode 方法](../html/cabcebca-60ed-7689-e44a-e227936dcc51.htm "Mode 方法 ")

[ModeAsync 方法](../html/e0956a5d-5ae1-f884-52dd-0e91b85828b5.htm "ModeAsync 方法 ")

[MoveJ 方法](../html/10990eab-b2f9-e4b8-3c88-5bffc78aa34a.htm "MoveJ 方法 ")

[MoveJAsync 方法](../html/5b05904a-27a0-b4f9-adc7-1fe2cc5b14a1.htm "MoveJAsync 方法 ")

[MSDP 方法](../html/1e78f8b5-7336-7fdf-a7e0-e40af1d3de7f.htm "MSDP 方法 ")

[MSDPAsync 方法](../html/c2a18a60-0dae-112b-af16-0e028e6a203a.htm "MSDPAsync 方法 ")

[Read 方法](../html/69082cc9-c365-8fc4-36df-082c0c6c22c3.htm "Read 方法 ")

[ReadALARM 方法](../html/b56d10b7-3ebc-effd-c8f4-1dbf2d9f35a0.htm "ReadALARM 方法 ")

[ReadALARMAsync 方法](../html/6bd50e3f-6923-2c34-0e8d-d08ba82840f1.htm "ReadALARMAsync 方法 ")

[ReadAsync 方法](../html/185f22c2-a4f1-985e-dd8e-4bf7107c5ae7.htm "ReadAsync 方法 ")

[ReadByCommand 方法](../html/80dbd712-4782-b205-7e73-977c37f5a214.htm "ReadByCommand 方法 ")

[ReadByCommandAsync 方法](../html/2ac7651a-e36e-1635-0a61-b855722693a1.htm "ReadByCommandAsync 方法 ")

[ReadByteVariable 方法](../html/e22f2610-10f3-8472-92c2-3f9e904ebb2b.htm "ReadByteVariable 方法 ")

[ReadByteVariableAsync 方法](../html/d01b8a15-6079-9c57-6936-e0c569cf5194.htm "ReadByteVariableAsync 方法 ")

[ReadDoubleIntegerVariable 方法](../html/beab7b65-77db-3c62-6f95-ccf796ecac03.htm "ReadDoubleIntegerVariable 方法 ")

[ReadDoubleIntegerVariableAsync 方法](../html/c6845ca5-1d25-d613-89f9-28da2515cf24.htm "ReadDoubleIntegerVariableAsync 方法 ")

[ReadFromCoreServer 方法](../html/4e4be36a-8ce8-a41b-6ec4-4f82479c6f23.htm "ReadFromCoreServer 方法 ")

[ReadFromCoreServerAsync 方法](../html/e4cb7ebe-d42a-ea17-8972-b246a42e9a90.htm "ReadFromCoreServerAsync 方法 ")

[ReadIntegerVariable 方法](../html/df3171c4-e18b-a52c-9032-bd789e7388a0.htm "ReadIntegerVariable 方法 ")

[ReadIntegerVariableAsync 方法](../html/954c86a5-37e1-f32b-3a6e-7fc4274cbb6f.htm "ReadIntegerVariableAsync 方法 ")

[ReadJSeq 方法](../html/d7a72f10-2923-e045-8c71-e8b43c551c5a.htm "ReadJSeq 方法 ")

[ReadJSeqAsync 方法](../html/ddfc3321-88c7-9d2a-b71e-9315ce77864d.htm "ReadJSeqAsync 方法 ")

[ReadPOSC 方法](../html/fea57b6b-1d2f-0547-73de-6120ee5de705.htm "ReadPOSC 方法 ")

[ReadPOSCAsync 方法](../html/5f5731c4-a35a-affa-3d76-7b5f50fe63a1.htm "ReadPOSCAsync 方法 ")

[ReadPOSJ 方法](../html/14045444-1d70-3bdf-9848-a886c204da46.htm "ReadPOSJ 方法 ")

[ReadPOSJAsync 方法](../html/0ed6f125-5c82-6a30-1a4e-a6a773dda5d8.htm "ReadPOSJAsync 方法 ")

[ReadRealVariable 方法](../html/5f6c460a-0bb9-32fb-8aa8-33b99dc8be5e.htm "ReadRealVariable 方法 ")

[ReadRealVariableAsync 方法](../html/996ae64e-84f9-ccf6-fa2a-de4e4a1b13c9.htm "ReadRealVariableAsync 方法 ")

[ReadStats 方法](../html/3a24171c-f204-8c63-cdbc-8cf23a965b75.htm "ReadStats 方法 ")

[ReadStatsAsync 方法](../html/43a75b5c-c5af-e81f-286c-a4aa2a7a898b.htm "ReadStatsAsync 方法 ")

[ReadString 方法](../html/dbc09266-42e8-b939-aea5-540d528d19ee.htm "ReadString 方法 ")

[ReadStringAsync 方法](../html/8fa97613-eb8d-c1e1-33f6-925dcf1a1ac0.htm "ReadStringAsync 方法 ")

[ReadStringVariable 方法](../html/5dc9095d-86e8-d434-0027-83797a1d8c0a.htm "ReadStringVariable 方法 ")

[ReadStringVariableAsync 方法](../html/f850e762-bf8a-6ee7-c18c-3496d70244fe.htm "ReadStringVariableAsync 方法 ")

[ReadUFrame 方法](../html/100a1955-fc0a-871a-98f9-bb361c422a3f.htm "ReadUFrame 方法 ")

[ReadUFrameAsync 方法](../html/f63d72f5-eaf4-ae5d-e80a-d881f93a878e.htm "ReadUFrameAsync 方法 ")

[Reset 方法](../html/bf61707b-875c-4b42-2732-9b2fb69535ac.htm "Reset 方法 ")

[ResetAsync 方法](../html/f135af7f-0250-eb16-48fe-959dd0aeeb72.htm "ResetAsync 方法 ")

[SetMJ 方法](../html/9cd8c255-f908-6f4e-da07-33edc743dbf8.htm "SetMJ 方法 ")

[SetMJAsync 方法](../html/74b1cb5d-83f7-61b7-8eb0-0e10a675b964.htm "SetMJAsync 方法 ")

[Start 方法](../html/39764e6e-fde9-c99b-9034-99a64232aef6.htm "Start 方法 ")

[StartAsync 方法](../html/ddefed38-f8b8-4ae1-73e3-eecbfae54043.htm "StartAsync 方法 ")

[Svon 方法](../html/1f83c8cb-fea3-d99f-c7ad-f12f39bd9dee.htm "Svon 方法 ")

[SvonAsync 方法](../html/9b51697e-6164-280a-d37c-1cc16f467b2d.htm "SvonAsync 方法 ")

[ToString 方法](../html/1a5b666e-cae2-1664-908f-56866da693e5.htm "ToString 方法 ")

[Write 方法](../html/27720779-dc92-a6b4-c792-940ec2a10d22.htm "Write 方法 ")

[WriteAsync 方法](../html/e1646488-a935-0a5b-ec6a-4d04a3e69d38.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| YRC1000TcpNetJSeqAsync 方法 |

设定执行程序的名称和行编号。  
Set the name and line number of the executed program.

**命名空间：**
 [HslCommunication.Robot.YASKAWA](3d67023f-78aa-c5d9-85b6-f18f5dfbc176.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public Task<OperateResult> JSeqAsync(
	string programName,
	int line
)
```

```
Public Function JSeqAsync ( 
	programName As String,
	line As Integer
) As Task(Of OperateResult)
```

```
public:
Task<OperateResult^>^ JSeqAsync(
	String^ programName, 
	int line
)
```

```
member JSeqAsync : 
        programName : string * 
        line : int -> Task<OperateResult> 
```

#### 参数

programName
:   类型：SystemString  
    设定程序名称

line
:   类型：SystemInt32  
    设定行编号（ 0 ～ 9999）

#### 返回值

类型：Task[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否设定成功

![](../icons/SectionExpanded.png)参见

#### 引用

[YRC1000TcpNet 类](cdcdc1cc-17cf-da27-7d20-ca6d6db1cd7e.htm)

[HslCommunication.Robot.YASKAWA 命名空间](3d67023f-78aa-c5d9-85b6-f18f5dfbc176.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Mode 方法 

[原文連結](http://api.hslcommunication.cn/html/cabcebca-60ed-7689-e44a-e227936dcc51.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.YASKAWA](../html/3d67023f-78aa-c5d9-85b6-f18f5dfbc176.htm "HslCommunication.Robot.YASKAWA")

[YRC1000TcpNet 类](../html/cdcdc1cc-17cf-da27-7d20-ca6d6db1cd7e.htm "YRC1000TcpNet 类")

[YRC1000TcpNet 方法](../html/7a26c354-eee0-b601-1d8a-51198e40e3fa.htm "YRC1000TcpNet 方法")

[Cancel 方法](../html/be66b6de-fe15-95b5-6969-70518f086388.htm "Cancel 方法 ")

[CancelAsync 方法](../html/fedb71cd-919c-c002-7556-551ae30fb10e.htm "CancelAsync 方法 ")

[Cycle 方法](../html/5a144fc7-5bc5-7250-6c13-826c172bd2ce.htm "Cycle 方法 ")

[CycleAsync 方法](../html/82bfe8a9-00b3-b34e-2a82-6d022becdf9a.htm "CycleAsync 方法 ")

[Delete 方法](../html/c7bfb022-46f4-8cdb-1b50-ffc1e5eb716b.htm "Delete 方法 ")

[DeleteAsync 方法](../html/c5a62bdd-3432-abcd-9d03-1ff7e30bd447.htm "DeleteAsync 方法 ")

[GetNewNetMessage 方法](../html/888815ce-da5c-0cb3-c6f9-800ec37a19ab.htm "GetNewNetMessage 方法 ")

[HLock 方法](../html/9c4ed34f-89fd-62e2-f707-08d1348c7541.htm "HLock 方法 ")

[HLockAsync 方法](../html/f0eb6982-a92b-4481-6435-33fa2b841b3c.htm "HLockAsync 方法 ")

[Hold 方法](../html/bf4ccc0a-1576-83ca-49e3-fde86b3a3738.htm "Hold 方法 ")

[HoldAsync 方法](../html/8c84a7d8-3677-0b6b-bdb5-44aa61ed1fcd.htm "HoldAsync 方法 ")

[InitializationOnConnect 方法](../html/5bd0274d-b3e7-f797-2dc0-c0dd2f330694.htm "InitializationOnConnect 方法 ")

[InitializationOnConnectAsync 方法](../html/e2f434c6-2ad9-30e8-8f2c-dcab95f5ec2b.htm "InitializationOnConnectAsync 方法 ")

[IORead 方法](../html/4539f381-2df3-4315-084c-7fcc81d523cb.htm "IORead 方法 ")

[IOReadAsync 方法](../html/5d550dc5-c02c-7245-122e-c3624e8d9c14.htm "IOReadAsync 方法 ")

[IOWrite 方法](../html/0dfecba5-f497-4b0f-55a5-339c7202d742.htm "IOWrite 方法 ")

[IOWriteAsync 方法](../html/64b5c202-29ff-0737-9c62-a4735bc00911.htm "IOWriteAsync 方法 ")

[JSeq 方法](../html/7057228b-aa47-3318-248d-1594588783bd.htm "JSeq 方法 ")

[JSeqAsync 方法](../html/c4eb7232-63e1-cd10-f89e-7f58072e3e9c.htm "JSeqAsync 方法 ")

[Mode 方法](../html/cabcebca-60ed-7689-e44a-e227936dcc51.htm "Mode 方法 ")

[ModeAsync 方法](../html/e0956a5d-5ae1-f884-52dd-0e91b85828b5.htm "ModeAsync 方法 ")

[MoveJ 方法](../html/10990eab-b2f9-e4b8-3c88-5bffc78aa34a.htm "MoveJ 方法 ")

[MoveJAsync 方法](../html/5b05904a-27a0-b4f9-adc7-1fe2cc5b14a1.htm "MoveJAsync 方法 ")

[MSDP 方法](../html/1e78f8b5-7336-7fdf-a7e0-e40af1d3de7f.htm "MSDP 方法 ")

[MSDPAsync 方法](../html/c2a18a60-0dae-112b-af16-0e028e6a203a.htm "MSDPAsync 方法 ")

[Read 方法](../html/69082cc9-c365-8fc4-36df-082c0c6c22c3.htm "Read 方法 ")

[ReadALARM 方法](../html/b56d10b7-3ebc-effd-c8f4-1dbf2d9f35a0.htm "ReadALARM 方法 ")

[ReadALARMAsync 方法](../html/6bd50e3f-6923-2c34-0e8d-d08ba82840f1.htm "ReadALARMAsync 方法 ")

[ReadAsync 方法](../html/185f22c2-a4f1-985e-dd8e-4bf7107c5ae7.htm "ReadAsync 方法 ")

[ReadByCommand 方法](../html/80dbd712-4782-b205-7e73-977c37f5a214.htm "ReadByCommand 方法 ")

[ReadByCommandAsync 方法](../html/2ac7651a-e36e-1635-0a61-b855722693a1.htm "ReadByCommandAsync 方法 ")

[ReadByteVariable 方法](../html/e22f2610-10f3-8472-92c2-3f9e904ebb2b.htm "ReadByteVariable 方法 ")

[ReadByteVariableAsync 方法](../html/d01b8a15-6079-9c57-6936-e0c569cf5194.htm "ReadByteVariableAsync 方法 ")

[ReadDoubleIntegerVariable 方法](../html/beab7b65-77db-3c62-6f95-ccf796ecac03.htm "ReadDoubleIntegerVariable 方法 ")

[ReadDoubleIntegerVariableAsync 方法](../html/c6845ca5-1d25-d613-89f9-28da2515cf24.htm "ReadDoubleIntegerVariableAsync 方法 ")

[ReadFromCoreServer 方法](../html/4e4be36a-8ce8-a41b-6ec4-4f82479c6f23.htm "ReadFromCoreServer 方法 ")

[ReadFromCoreServerAsync 方法](../html/e4cb7ebe-d42a-ea17-8972-b246a42e9a90.htm "ReadFromCoreServerAsync 方法 ")

[ReadIntegerVariable 方法](../html/df3171c4-e18b-a52c-9032-bd789e7388a0.htm "ReadIntegerVariable 方法 ")

[ReadIntegerVariableAsync 方法](../html/954c86a5-37e1-f32b-3a6e-7fc4274cbb6f.htm "ReadIntegerVariableAsync 方法 ")

[ReadJSeq 方法](../html/d7a72f10-2923-e045-8c71-e8b43c551c5a.htm "ReadJSeq 方法 ")

[ReadJSeqAsync 方法](../html/ddfc3321-88c7-9d2a-b71e-9315ce77864d.htm "ReadJSeqAsync 方法 ")

[ReadPOSC 方法](../html/fea57b6b-1d2f-0547-73de-6120ee5de705.htm "ReadPOSC 方法 ")

[ReadPOSCAsync 方法](../html/5f5731c4-a35a-affa-3d76-7b5f50fe63a1.htm "ReadPOSCAsync 方法 ")

[ReadPOSJ 方法](../html/14045444-1d70-3bdf-9848-a886c204da46.htm "ReadPOSJ 方法 ")

[ReadPOSJAsync 方法](../html/0ed6f125-5c82-6a30-1a4e-a6a773dda5d8.htm "ReadPOSJAsync 方法 ")

[ReadRealVariable 方法](../html/5f6c460a-0bb9-32fb-8aa8-33b99dc8be5e.htm "ReadRealVariable 方法 ")

[ReadRealVariableAsync 方法](../html/996ae64e-84f9-ccf6-fa2a-de4e4a1b13c9.htm "ReadRealVariableAsync 方法 ")

[ReadStats 方法](../html/3a24171c-f204-8c63-cdbc-8cf23a965b75.htm "ReadStats 方法 ")

[ReadStatsAsync 方法](../html/43a75b5c-c5af-e81f-286c-a4aa2a7a898b.htm "ReadStatsAsync 方法 ")

[ReadString 方法](../html/dbc09266-42e8-b939-aea5-540d528d19ee.htm "ReadString 方法 ")

[ReadStringAsync 方法](../html/8fa97613-eb8d-c1e1-33f6-925dcf1a1ac0.htm "ReadStringAsync 方法 ")

[ReadStringVariable 方法](../html/5dc9095d-86e8-d434-0027-83797a1d8c0a.htm "ReadStringVariable 方法 ")

[ReadStringVariableAsync 方法](../html/f850e762-bf8a-6ee7-c18c-3496d70244fe.htm "ReadStringVariableAsync 方法 ")

[ReadUFrame 方法](../html/100a1955-fc0a-871a-98f9-bb361c422a3f.htm "ReadUFrame 方法 ")

[ReadUFrameAsync 方法](../html/f63d72f5-eaf4-ae5d-e80a-d881f93a878e.htm "ReadUFrameAsync 方法 ")

[Reset 方法](../html/bf61707b-875c-4b42-2732-9b2fb69535ac.htm "Reset 方法 ")

[ResetAsync 方法](../html/f135af7f-0250-eb16-48fe-959dd0aeeb72.htm "ResetAsync 方法 ")

[SetMJ 方法](../html/9cd8c255-f908-6f4e-da07-33edc743dbf8.htm "SetMJ 方法 ")

[SetMJAsync 方法](../html/74b1cb5d-83f7-61b7-8eb0-0e10a675b964.htm "SetMJAsync 方法 ")

[Start 方法](../html/39764e6e-fde9-c99b-9034-99a64232aef6.htm "Start 方法 ")

[StartAsync 方法](../html/ddefed38-f8b8-4ae1-73e3-eecbfae54043.htm "StartAsync 方法 ")

[Svon 方法](../html/1f83c8cb-fea3-d99f-c7ad-f12f39bd9dee.htm "Svon 方法 ")

[SvonAsync 方法](../html/9b51697e-6164-280a-d37c-1cc16f467b2d.htm "SvonAsync 方法 ")

[ToString 方法](../html/1a5b666e-cae2-1664-908f-56866da693e5.htm "ToString 方法 ")

[Write 方法](../html/27720779-dc92-a6b4-c792-940ec2a10d22.htm "Write 方法 ")

[WriteAsync 方法](../html/e1646488-a935-0a5b-ec6a-4d04a3e69d38.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| YRC1000TcpNetMode 方法 |

选择模式。模式编号为1:示教模式，2:再现模式  
Choose a mode. The mode number is 1: teaching mode, 2: reproduction mode

**命名空间：**
 [HslCommunication.Robot.YASKAWA](3d67023f-78aa-c5d9-85b6-f18f5dfbc176.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public OperateResult Mode(
	int number
)
```

```
Public Function Mode ( 
	number As Integer
) As OperateResult
```

```
public:
OperateResult^ Mode(
	int number
)
```

```
member Mode : 
        number : int -> OperateResult 
```

#### 参数

number
:   类型：SystemInt32  
    模式编号为1:示教模式，2:再现模式

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
模式是否选择成功

![](../icons/SectionExpanded.png)备注

MODE 命令，是在「操作条件」 画面中获得外部模式切换的许可后可以使用。

![](../icons/SectionExpanded.png)参见

#### 引用

[YRC1000TcpNet 类](cdcdc1cc-17cf-da27-7d20-ca6d6db1cd7e.htm)

[HslCommunication.Robot.YASKAWA 命名空间](3d67023f-78aa-c5d9-85b6-f18f5dfbc176.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ModeAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/e0956a5d-5ae1-f884-52dd-0e91b85828b5.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.YASKAWA](../html/3d67023f-78aa-c5d9-85b6-f18f5dfbc176.htm "HslCommunication.Robot.YASKAWA")

[YRC1000TcpNet 类](../html/cdcdc1cc-17cf-da27-7d20-ca6d6db1cd7e.htm "YRC1000TcpNet 类")

[YRC1000TcpNet 方法](../html/7a26c354-eee0-b601-1d8a-51198e40e3fa.htm "YRC1000TcpNet 方法")

[Cancel 方法](../html/be66b6de-fe15-95b5-6969-70518f086388.htm "Cancel 方法 ")

[CancelAsync 方法](../html/fedb71cd-919c-c002-7556-551ae30fb10e.htm "CancelAsync 方法 ")

[Cycle 方法](../html/5a144fc7-5bc5-7250-6c13-826c172bd2ce.htm "Cycle 方法 ")

[CycleAsync 方法](../html/82bfe8a9-00b3-b34e-2a82-6d022becdf9a.htm "CycleAsync 方法 ")

[Delete 方法](../html/c7bfb022-46f4-8cdb-1b50-ffc1e5eb716b.htm "Delete 方法 ")

[DeleteAsync 方法](../html/c5a62bdd-3432-abcd-9d03-1ff7e30bd447.htm "DeleteAsync 方法 ")

[GetNewNetMessage 方法](../html/888815ce-da5c-0cb3-c6f9-800ec37a19ab.htm "GetNewNetMessage 方法 ")

[HLock 方法](../html/9c4ed34f-89fd-62e2-f707-08d1348c7541.htm "HLock 方法 ")

[HLockAsync 方法](../html/f0eb6982-a92b-4481-6435-33fa2b841b3c.htm "HLockAsync 方法 ")

[Hold 方法](../html/bf4ccc0a-1576-83ca-49e3-fde86b3a3738.htm "Hold 方法 ")

[HoldAsync 方法](../html/8c84a7d8-3677-0b6b-bdb5-44aa61ed1fcd.htm "HoldAsync 方法 ")

[InitializationOnConnect 方法](../html/5bd0274d-b3e7-f797-2dc0-c0dd2f330694.htm "InitializationOnConnect 方法 ")

[InitializationOnConnectAsync 方法](../html/e2f434c6-2ad9-30e8-8f2c-dcab95f5ec2b.htm "InitializationOnConnectAsync 方法 ")

[IORead 方法](../html/4539f381-2df3-4315-084c-7fcc81d523cb.htm "IORead 方法 ")

[IOReadAsync 方法](../html/5d550dc5-c02c-7245-122e-c3624e8d9c14.htm "IOReadAsync 方法 ")

[IOWrite 方法](../html/0dfecba5-f497-4b0f-55a5-339c7202d742.htm "IOWrite 方法 ")

[IOWriteAsync 方法](../html/64b5c202-29ff-0737-9c62-a4735bc00911.htm "IOWriteAsync 方法 ")

[JSeq 方法](../html/7057228b-aa47-3318-248d-1594588783bd.htm "JSeq 方法 ")

[JSeqAsync 方法](../html/c4eb7232-63e1-cd10-f89e-7f58072e3e9c.htm "JSeqAsync 方法 ")

[Mode 方法](../html/cabcebca-60ed-7689-e44a-e227936dcc51.htm "Mode 方法 ")

[ModeAsync 方法](../html/e0956a5d-5ae1-f884-52dd-0e91b85828b5.htm "ModeAsync 方法 ")

[MoveJ 方法](../html/10990eab-b2f9-e4b8-3c88-5bffc78aa34a.htm "MoveJ 方法 ")

[MoveJAsync 方法](../html/5b05904a-27a0-b4f9-adc7-1fe2cc5b14a1.htm "MoveJAsync 方法 ")

[MSDP 方法](../html/1e78f8b5-7336-7fdf-a7e0-e40af1d3de7f.htm "MSDP 方法 ")

[MSDPAsync 方法](../html/c2a18a60-0dae-112b-af16-0e028e6a203a.htm "MSDPAsync 方法 ")

[Read 方法](../html/69082cc9-c365-8fc4-36df-082c0c6c22c3.htm "Read 方法 ")

[ReadALARM 方法](../html/b56d10b7-3ebc-effd-c8f4-1dbf2d9f35a0.htm "ReadALARM 方法 ")

[ReadALARMAsync 方法](../html/6bd50e3f-6923-2c34-0e8d-d08ba82840f1.htm "ReadALARMAsync 方法 ")

[ReadAsync 方法](../html/185f22c2-a4f1-985e-dd8e-4bf7107c5ae7.htm "ReadAsync 方法 ")

[ReadByCommand 方法](../html/80dbd712-4782-b205-7e73-977c37f5a214.htm "ReadByCommand 方法 ")

[ReadByCommandAsync 方法](../html/2ac7651a-e36e-1635-0a61-b855722693a1.htm "ReadByCommandAsync 方法 ")

[ReadByteVariable 方法](../html/e22f2610-10f3-8472-92c2-3f9e904ebb2b.htm "ReadByteVariable 方法 ")

[ReadByteVariableAsync 方法](../html/d01b8a15-6079-9c57-6936-e0c569cf5194.htm "ReadByteVariableAsync 方法 ")

[ReadDoubleIntegerVariable 方法](../html/beab7b65-77db-3c62-6f95-ccf796ecac03.htm "ReadDoubleIntegerVariable 方法 ")

[ReadDoubleIntegerVariableAsync 方法](../html/c6845ca5-1d25-d613-89f9-28da2515cf24.htm "ReadDoubleIntegerVariableAsync 方法 ")

[ReadFromCoreServer 方法](../html/4e4be36a-8ce8-a41b-6ec4-4f82479c6f23.htm "ReadFromCoreServer 方法 ")

[ReadFromCoreServerAsync 方法](../html/e4cb7ebe-d42a-ea17-8972-b246a42e9a90.htm "ReadFromCoreServerAsync 方法 ")

[ReadIntegerVariable 方法](../html/df3171c4-e18b-a52c-9032-bd789e7388a0.htm "ReadIntegerVariable 方法 ")

[ReadIntegerVariableAsync 方法](../html/954c86a5-37e1-f32b-3a6e-7fc4274cbb6f.htm "ReadIntegerVariableAsync 方法 ")

[ReadJSeq 方法](../html/d7a72f10-2923-e045-8c71-e8b43c551c5a.htm "ReadJSeq 方法 ")

[ReadJSeqAsync 方法](../html/ddfc3321-88c7-9d2a-b71e-9315ce77864d.htm "ReadJSeqAsync 方法 ")

[ReadPOSC 方法](../html/fea57b6b-1d2f-0547-73de-6120ee5de705.htm "ReadPOSC 方法 ")

[ReadPOSCAsync 方法](../html/5f5731c4-a35a-affa-3d76-7b5f50fe63a1.htm "ReadPOSCAsync 方法 ")

[ReadPOSJ 方法](../html/14045444-1d70-3bdf-9848-a886c204da46.htm "ReadPOSJ 方法 ")

[ReadPOSJAsync 方法](../html/0ed6f125-5c82-6a30-1a4e-a6a773dda5d8.htm "ReadPOSJAsync 方法 ")

[ReadRealVariable 方法](../html/5f6c460a-0bb9-32fb-8aa8-33b99dc8be5e.htm "ReadRealVariable 方法 ")

[ReadRealVariableAsync 方法](../html/996ae64e-84f9-ccf6-fa2a-de4e4a1b13c9.htm "ReadRealVariableAsync 方法 ")

[ReadStats 方法](../html/3a24171c-f204-8c63-cdbc-8cf23a965b75.htm "ReadStats 方法 ")

[ReadStatsAsync 方法](../html/43a75b5c-c5af-e81f-286c-a4aa2a7a898b.htm "ReadStatsAsync 方法 ")

[ReadString 方法](../html/dbc09266-42e8-b939-aea5-540d528d19ee.htm "ReadString 方法 ")

[ReadStringAsync 方法](../html/8fa97613-eb8d-c1e1-33f6-925dcf1a1ac0.htm "ReadStringAsync 方法 ")

[ReadStringVariable 方法](../html/5dc9095d-86e8-d434-0027-83797a1d8c0a.htm "ReadStringVariable 方法 ")

[ReadStringVariableAsync 方法](../html/f850e762-bf8a-6ee7-c18c-3496d70244fe.htm "ReadStringVariableAsync 方法 ")

[ReadUFrame 方法](../html/100a1955-fc0a-871a-98f9-bb361c422a3f.htm "ReadUFrame 方法 ")

[ReadUFrameAsync 方法](../html/f63d72f5-eaf4-ae5d-e80a-d881f93a878e.htm "ReadUFrameAsync 方法 ")

[Reset 方法](../html/bf61707b-875c-4b42-2732-9b2fb69535ac.htm "Reset 方法 ")

[ResetAsync 方法](../html/f135af7f-0250-eb16-48fe-959dd0aeeb72.htm "ResetAsync 方法 ")

[SetMJ 方法](../html/9cd8c255-f908-6f4e-da07-33edc743dbf8.htm "SetMJ 方法 ")

[SetMJAsync 方法](../html/74b1cb5d-83f7-61b7-8eb0-0e10a675b964.htm "SetMJAsync 方法 ")

[Start 方法](../html/39764e6e-fde9-c99b-9034-99a64232aef6.htm "Start 方法 ")

[StartAsync 方法](../html/ddefed38-f8b8-4ae1-73e3-eecbfae54043.htm "StartAsync 方法 ")

[Svon 方法](../html/1f83c8cb-fea3-d99f-c7ad-f12f39bd9dee.htm "Svon 方法 ")

[SvonAsync 方法](../html/9b51697e-6164-280a-d37c-1cc16f467b2d.htm "SvonAsync 方法 ")

[ToString 方法](../html/1a5b666e-cae2-1664-908f-56866da693e5.htm "ToString 方法 ")

[Write 方法](../html/27720779-dc92-a6b4-c792-940ec2a10d22.htm "Write 方法 ")

[WriteAsync 方法](../html/e1646488-a935-0a5b-ec6a-4d04a3e69d38.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| YRC1000TcpNetModeAsync 方法 |

选择模式。模式编号为1:示教模式，2:再现模式  
Choose a mode. The mode number is 1: teaching mode, 2: reproduction mode

**命名空间：**
 [HslCommunication.Robot.YASKAWA](3d67023f-78aa-c5d9-85b6-f18f5dfbc176.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public Task<OperateResult> ModeAsync(
	int number
)
```

```
Public Function ModeAsync ( 
	number As Integer
) As Task(Of OperateResult)
```

```
public:
Task<OperateResult^>^ ModeAsync(
	int number
)
```

```
member ModeAsync : 
        number : int -> Task<OperateResult> 
```

#### 参数

number
:   类型：SystemInt32  
    模式编号为1:示教模式，2:再现模式

#### 返回值

类型：Task[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
模式是否选择成功

![](../icons/SectionExpanded.png)备注

MODE 命令，是在「操作条件」 画面中获得外部模式切换的许可后可以使用。

![](../icons/SectionExpanded.png)参见

#### 引用

[YRC1000TcpNet 类](cdcdc1cc-17cf-da27-7d20-ca6d6db1cd7e.htm)

[HslCommunication.Robot.YASKAWA 命名空间](3d67023f-78aa-c5d9-85b6-f18f5dfbc176.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## MoveJ 方法 

[原文連結](http://api.hslcommunication.cn/html/10990eab-b2f9-e4b8-3c88-5bffc78aa34a.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.YASKAWA](../html/3d67023f-78aa-c5d9-85b6-f18f5dfbc176.htm "HslCommunication.Robot.YASKAWA")

[YRC1000TcpNet 类](../html/cdcdc1cc-17cf-da27-7d20-ca6d6db1cd7e.htm "YRC1000TcpNet 类")

[YRC1000TcpNet 方法](../html/7a26c354-eee0-b601-1d8a-51198e40e3fa.htm "YRC1000TcpNet 方法")

[Cancel 方法](../html/be66b6de-fe15-95b5-6969-70518f086388.htm "Cancel 方法 ")

[CancelAsync 方法](../html/fedb71cd-919c-c002-7556-551ae30fb10e.htm "CancelAsync 方法 ")

[Cycle 方法](../html/5a144fc7-5bc5-7250-6c13-826c172bd2ce.htm "Cycle 方法 ")

[CycleAsync 方法](../html/82bfe8a9-00b3-b34e-2a82-6d022becdf9a.htm "CycleAsync 方法 ")

[Delete 方法](../html/c7bfb022-46f4-8cdb-1b50-ffc1e5eb716b.htm "Delete 方法 ")

[DeleteAsync 方法](../html/c5a62bdd-3432-abcd-9d03-1ff7e30bd447.htm "DeleteAsync 方法 ")

[GetNewNetMessage 方法](../html/888815ce-da5c-0cb3-c6f9-800ec37a19ab.htm "GetNewNetMessage 方法 ")

[HLock 方法](../html/9c4ed34f-89fd-62e2-f707-08d1348c7541.htm "HLock 方法 ")

[HLockAsync 方法](../html/f0eb6982-a92b-4481-6435-33fa2b841b3c.htm "HLockAsync 方法 ")

[Hold 方法](../html/bf4ccc0a-1576-83ca-49e3-fde86b3a3738.htm "Hold 方法 ")

[HoldAsync 方法](../html/8c84a7d8-3677-0b6b-bdb5-44aa61ed1fcd.htm "HoldAsync 方法 ")

[InitializationOnConnect 方法](../html/5bd0274d-b3e7-f797-2dc0-c0dd2f330694.htm "InitializationOnConnect 方法 ")

[InitializationOnConnectAsync 方法](../html/e2f434c6-2ad9-30e8-8f2c-dcab95f5ec2b.htm "InitializationOnConnectAsync 方法 ")

[IORead 方法](../html/4539f381-2df3-4315-084c-7fcc81d523cb.htm "IORead 方法 ")

[IOReadAsync 方法](../html/5d550dc5-c02c-7245-122e-c3624e8d9c14.htm "IOReadAsync 方法 ")

[IOWrite 方法](../html/0dfecba5-f497-4b0f-55a5-339c7202d742.htm "IOWrite 方法 ")

[IOWriteAsync 方法](../html/64b5c202-29ff-0737-9c62-a4735bc00911.htm "IOWriteAsync 方法 ")

[JSeq 方法](../html/7057228b-aa47-3318-248d-1594588783bd.htm "JSeq 方法 ")

[JSeqAsync 方法](../html/c4eb7232-63e1-cd10-f89e-7f58072e3e9c.htm "JSeqAsync 方法 ")

[Mode 方法](../html/cabcebca-60ed-7689-e44a-e227936dcc51.htm "Mode 方法 ")

[ModeAsync 方法](../html/e0956a5d-5ae1-f884-52dd-0e91b85828b5.htm "ModeAsync 方法 ")

[MoveJ 方法](../html/10990eab-b2f9-e4b8-3c88-5bffc78aa34a.htm "MoveJ 方法 ")

[MoveJAsync 方法](../html/5b05904a-27a0-b4f9-adc7-1fe2cc5b14a1.htm "MoveJAsync 方法 ")

[MSDP 方法](../html/1e78f8b5-7336-7fdf-a7e0-e40af1d3de7f.htm "MSDP 方法 ")

[MSDPAsync 方法](../html/c2a18a60-0dae-112b-af16-0e028e6a203a.htm "MSDPAsync 方法 ")

[Read 方法](../html/69082cc9-c365-8fc4-36df-082c0c6c22c3.htm "Read 方法 ")

[ReadALARM 方法](../html/b56d10b7-3ebc-effd-c8f4-1dbf2d9f35a0.htm "ReadALARM 方法 ")

[ReadALARMAsync 方法](../html/6bd50e3f-6923-2c34-0e8d-d08ba82840f1.htm "ReadALARMAsync 方法 ")

[ReadAsync 方法](../html/185f22c2-a4f1-985e-dd8e-4bf7107c5ae7.htm "ReadAsync 方法 ")

[ReadByCommand 方法](../html/80dbd712-4782-b205-7e73-977c37f5a214.htm "ReadByCommand 方法 ")

[ReadByCommandAsync 方法](../html/2ac7651a-e36e-1635-0a61-b855722693a1.htm "ReadByCommandAsync 方法 ")

[ReadByteVariable 方法](../html/e22f2610-10f3-8472-92c2-3f9e904ebb2b.htm "ReadByteVariable 方法 ")

[ReadByteVariableAsync 方法](../html/d01b8a15-6079-9c57-6936-e0c569cf5194.htm "ReadByteVariableAsync 方法 ")

[ReadDoubleIntegerVariable 方法](../html/beab7b65-77db-3c62-6f95-ccf796ecac03.htm "ReadDoubleIntegerVariable 方法 ")

[ReadDoubleIntegerVariableAsync 方法](../html/c6845ca5-1d25-d613-89f9-28da2515cf24.htm "ReadDoubleIntegerVariableAsync 方法 ")

[ReadFromCoreServer 方法](../html/4e4be36a-8ce8-a41b-6ec4-4f82479c6f23.htm "ReadFromCoreServer 方法 ")

[ReadFromCoreServerAsync 方法](../html/e4cb7ebe-d42a-ea17-8972-b246a42e9a90.htm "ReadFromCoreServerAsync 方法 ")

[ReadIntegerVariable 方法](../html/df3171c4-e18b-a52c-9032-bd789e7388a0.htm "ReadIntegerVariable 方法 ")

[ReadIntegerVariableAsync 方法](../html/954c86a5-37e1-f32b-3a6e-7fc4274cbb6f.htm "ReadIntegerVariableAsync 方法 ")

[ReadJSeq 方法](../html/d7a72f10-2923-e045-8c71-e8b43c551c5a.htm "ReadJSeq 方法 ")

[ReadJSeqAsync 方法](../html/ddfc3321-88c7-9d2a-b71e-9315ce77864d.htm "ReadJSeqAsync 方法 ")

[ReadPOSC 方法](../html/fea57b6b-1d2f-0547-73de-6120ee5de705.htm "ReadPOSC 方法 ")

[ReadPOSCAsync 方法](../html/5f5731c4-a35a-affa-3d76-7b5f50fe63a1.htm "ReadPOSCAsync 方法 ")

[ReadPOSJ 方法](../html/14045444-1d70-3bdf-9848-a886c204da46.htm "ReadPOSJ 方法 ")

[ReadPOSJAsync 方法](../html/0ed6f125-5c82-6a30-1a4e-a6a773dda5d8.htm "ReadPOSJAsync 方法 ")

[ReadRealVariable 方法](../html/5f6c460a-0bb9-32fb-8aa8-33b99dc8be5e.htm "ReadRealVariable 方法 ")

[ReadRealVariableAsync 方法](../html/996ae64e-84f9-ccf6-fa2a-de4e4a1b13c9.htm "ReadRealVariableAsync 方法 ")

[ReadStats 方法](../html/3a24171c-f204-8c63-cdbc-8cf23a965b75.htm "ReadStats 方法 ")

[ReadStatsAsync 方法](../html/43a75b5c-c5af-e81f-286c-a4aa2a7a898b.htm "ReadStatsAsync 方法 ")

[ReadString 方法](../html/dbc09266-42e8-b939-aea5-540d528d19ee.htm "ReadString 方法 ")

[ReadStringAsync 方法](../html/8fa97613-eb8d-c1e1-33f6-925dcf1a1ac0.htm "ReadStringAsync 方法 ")

[ReadStringVariable 方法](../html/5dc9095d-86e8-d434-0027-83797a1d8c0a.htm "ReadStringVariable 方法 ")

[ReadStringVariableAsync 方法](../html/f850e762-bf8a-6ee7-c18c-3496d70244fe.htm "ReadStringVariableAsync 方法 ")

[ReadUFrame 方法](../html/100a1955-fc0a-871a-98f9-bb361c422a3f.htm "ReadUFrame 方法 ")

[ReadUFrameAsync 方法](../html/f63d72f5-eaf4-ae5d-e80a-d881f93a878e.htm "ReadUFrameAsync 方法 ")

[Reset 方法](../html/bf61707b-875c-4b42-2732-9b2fb69535ac.htm "Reset 方法 ")

[ResetAsync 方法](../html/f135af7f-0250-eb16-48fe-959dd0aeeb72.htm "ResetAsync 方法 ")

[SetMJ 方法](../html/9cd8c255-f908-6f4e-da07-33edc743dbf8.htm "SetMJ 方法 ")

[SetMJAsync 方法](../html/74b1cb5d-83f7-61b7-8eb0-0e10a675b964.htm "SetMJAsync 方法 ")

[Start 方法](../html/39764e6e-fde9-c99b-9034-99a64232aef6.htm "Start 方法 ")

[StartAsync 方法](../html/ddefed38-f8b8-4ae1-73e3-eecbfae54043.htm "StartAsync 方法 ")

[Svon 方法](../html/1f83c8cb-fea3-d99f-c7ad-f12f39bd9dee.htm "Svon 方法 ")

[SvonAsync 方法](../html/9b51697e-6164-280a-d37c-1cc16f467b2d.htm "SvonAsync 方法 ")

[ToString 方法](../html/1a5b666e-cae2-1664-908f-56866da693e5.htm "ToString 方法 ")

[Write 方法](../html/27720779-dc92-a6b4-c792-940ec2a10d22.htm "Write 方法 ")

[WriteAsync 方法](../html/e1646488-a935-0a5b-ec6a-4d04a3e69d38.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| YRC1000TcpNetMoveJ 方法 |

向指定的坐标系位置进行关节动作。其中没有外部轴的系统， 7-12外部轴的值设定为「0」  
Perform joint motions to the specified coordinate system position.
where there is no external axis system, the value of 7-12 external axis is set to "0"

**命名空间：**
 [HslCommunication.Robot.YASKAWA](3d67023f-78aa-c5d9-85b6-f18f5dfbc176.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public OperateResult MoveJ(
	YRCRobotData robotData
)
```

```
Public Function MoveJ ( 
	robotData As YRCRobotData
) As OperateResult
```

```
public:
OperateResult^ MoveJ(
	YRCRobotData^ robotData
)
```

```
member MoveJ : 
        robotData : YRCRobotData -> OperateResult 
```

#### 参数

robotData
:   类型：[HslCommunication.Robot.YASKAWAYRCRobotData](233ef65f-4a8d-20e6-03ed-67c1cbefd2d6.htm)  
    机器的的数据信息

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否动作成功

![](../icons/SectionExpanded.png)备注

其中形态数据由6个bool数组组成，每个bool含义参考参数说明，0表示 False，1表示 True

![](../icons/SectionExpanded.png)参见

#### 引用

[YRC1000TcpNet 类](cdcdc1cc-17cf-da27-7d20-ca6d6db1cd7e.htm)

[HslCommunication.Robot.YASKAWA 命名空间](3d67023f-78aa-c5d9-85b6-f18f5dfbc176.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## MoveJAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/5b05904a-27a0-b4f9-adc7-1fe2cc5b14a1.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.YASKAWA](../html/3d67023f-78aa-c5d9-85b6-f18f5dfbc176.htm "HslCommunication.Robot.YASKAWA")

[YRC1000TcpNet 类](../html/cdcdc1cc-17cf-da27-7d20-ca6d6db1cd7e.htm "YRC1000TcpNet 类")

[YRC1000TcpNet 方法](../html/7a26c354-eee0-b601-1d8a-51198e40e3fa.htm "YRC1000TcpNet 方法")

[Cancel 方法](../html/be66b6de-fe15-95b5-6969-70518f086388.htm "Cancel 方法 ")

[CancelAsync 方法](../html/fedb71cd-919c-c002-7556-551ae30fb10e.htm "CancelAsync 方法 ")

[Cycle 方法](../html/5a144fc7-5bc5-7250-6c13-826c172bd2ce.htm "Cycle 方法 ")

[CycleAsync 方法](../html/82bfe8a9-00b3-b34e-2a82-6d022becdf9a.htm "CycleAsync 方法 ")

[Delete 方法](../html/c7bfb022-46f4-8cdb-1b50-ffc1e5eb716b.htm "Delete 方法 ")

[DeleteAsync 方法](../html/c5a62bdd-3432-abcd-9d03-1ff7e30bd447.htm "DeleteAsync 方法 ")

[GetNewNetMessage 方法](../html/888815ce-da5c-0cb3-c6f9-800ec37a19ab.htm "GetNewNetMessage 方法 ")

[HLock 方法](../html/9c4ed34f-89fd-62e2-f707-08d1348c7541.htm "HLock 方法 ")

[HLockAsync 方法](../html/f0eb6982-a92b-4481-6435-33fa2b841b3c.htm "HLockAsync 方法 ")

[Hold 方法](../html/bf4ccc0a-1576-83ca-49e3-fde86b3a3738.htm "Hold 方法 ")

[HoldAsync 方法](../html/8c84a7d8-3677-0b6b-bdb5-44aa61ed1fcd.htm "HoldAsync 方法 ")

[InitializationOnConnect 方法](../html/5bd0274d-b3e7-f797-2dc0-c0dd2f330694.htm "InitializationOnConnect 方法 ")

[InitializationOnConnectAsync 方法](../html/e2f434c6-2ad9-30e8-8f2c-dcab95f5ec2b.htm "InitializationOnConnectAsync 方法 ")

[IORead 方法](../html/4539f381-2df3-4315-084c-7fcc81d523cb.htm "IORead 方法 ")

[IOReadAsync 方法](../html/5d550dc5-c02c-7245-122e-c3624e8d9c14.htm "IOReadAsync 方法 ")

[IOWrite 方法](../html/0dfecba5-f497-4b0f-55a5-339c7202d742.htm "IOWrite 方法 ")

[IOWriteAsync 方法](../html/64b5c202-29ff-0737-9c62-a4735bc00911.htm "IOWriteAsync 方法 ")

[JSeq 方法](../html/7057228b-aa47-3318-248d-1594588783bd.htm "JSeq 方法 ")

[JSeqAsync 方法](../html/c4eb7232-63e1-cd10-f89e-7f58072e3e9c.htm "JSeqAsync 方法 ")

[Mode 方法](../html/cabcebca-60ed-7689-e44a-e227936dcc51.htm "Mode 方法 ")

[ModeAsync 方法](../html/e0956a5d-5ae1-f884-52dd-0e91b85828b5.htm "ModeAsync 方法 ")

[MoveJ 方法](../html/10990eab-b2f9-e4b8-3c88-5bffc78aa34a.htm "MoveJ 方法 ")

[MoveJAsync 方法](../html/5b05904a-27a0-b4f9-adc7-1fe2cc5b14a1.htm "MoveJAsync 方法 ")

[MSDP 方法](../html/1e78f8b5-7336-7fdf-a7e0-e40af1d3de7f.htm "MSDP 方法 ")

[MSDPAsync 方法](../html/c2a18a60-0dae-112b-af16-0e028e6a203a.htm "MSDPAsync 方法 ")

[Read 方法](../html/69082cc9-c365-8fc4-36df-082c0c6c22c3.htm "Read 方法 ")

[ReadALARM 方法](../html/b56d10b7-3ebc-effd-c8f4-1dbf2d9f35a0.htm "ReadALARM 方法 ")

[ReadALARMAsync 方法](../html/6bd50e3f-6923-2c34-0e8d-d08ba82840f1.htm "ReadALARMAsync 方法 ")

[ReadAsync 方法](../html/185f22c2-a4f1-985e-dd8e-4bf7107c5ae7.htm "ReadAsync 方法 ")

[ReadByCommand 方法](../html/80dbd712-4782-b205-7e73-977c37f5a214.htm "ReadByCommand 方法 ")

[ReadByCommandAsync 方法](../html/2ac7651a-e36e-1635-0a61-b855722693a1.htm "ReadByCommandAsync 方法 ")

[ReadByteVariable 方法](../html/e22f2610-10f3-8472-92c2-3f9e904ebb2b.htm "ReadByteVariable 方法 ")

[ReadByteVariableAsync 方法](../html/d01b8a15-6079-9c57-6936-e0c569cf5194.htm "ReadByteVariableAsync 方法 ")

[ReadDoubleIntegerVariable 方法](../html/beab7b65-77db-3c62-6f95-ccf796ecac03.htm "ReadDoubleIntegerVariable 方法 ")

[ReadDoubleIntegerVariableAsync 方法](../html/c6845ca5-1d25-d613-89f9-28da2515cf24.htm "ReadDoubleIntegerVariableAsync 方法 ")

[ReadFromCoreServer 方法](../html/4e4be36a-8ce8-a41b-6ec4-4f82479c6f23.htm "ReadFromCoreServer 方法 ")

[ReadFromCoreServerAsync 方法](../html/e4cb7ebe-d42a-ea17-8972-b246a42e9a90.htm "ReadFromCoreServerAsync 方法 ")

[ReadIntegerVariable 方法](../html/df3171c4-e18b-a52c-9032-bd789e7388a0.htm "ReadIntegerVariable 方法 ")

[ReadIntegerVariableAsync 方法](../html/954c86a5-37e1-f32b-3a6e-7fc4274cbb6f.htm "ReadIntegerVariableAsync 方法 ")

[ReadJSeq 方法](../html/d7a72f10-2923-e045-8c71-e8b43c551c5a.htm "ReadJSeq 方法 ")

[ReadJSeqAsync 方法](../html/ddfc3321-88c7-9d2a-b71e-9315ce77864d.htm "ReadJSeqAsync 方法 ")

[ReadPOSC 方法](../html/fea57b6b-1d2f-0547-73de-6120ee5de705.htm "ReadPOSC 方法 ")

[ReadPOSCAsync 方法](../html/5f5731c4-a35a-affa-3d76-7b5f50fe63a1.htm "ReadPOSCAsync 方法 ")

[ReadPOSJ 方法](../html/14045444-1d70-3bdf-9848-a886c204da46.htm "ReadPOSJ 方法 ")

[ReadPOSJAsync 方法](../html/0ed6f125-5c82-6a30-1a4e-a6a773dda5d8.htm "ReadPOSJAsync 方法 ")

[ReadRealVariable 方法](../html/5f6c460a-0bb9-32fb-8aa8-33b99dc8be5e.htm "ReadRealVariable 方法 ")

[ReadRealVariableAsync 方法](../html/996ae64e-84f9-ccf6-fa2a-de4e4a1b13c9.htm "ReadRealVariableAsync 方法 ")

[ReadStats 方法](../html/3a24171c-f204-8c63-cdbc-8cf23a965b75.htm "ReadStats 方法 ")

[ReadStatsAsync 方法](../html/43a75b5c-c5af-e81f-286c-a4aa2a7a898b.htm "ReadStatsAsync 方法 ")

[ReadString 方法](../html/dbc09266-42e8-b939-aea5-540d528d19ee.htm "ReadString 方法 ")

[ReadStringAsync 方法](../html/8fa97613-eb8d-c1e1-33f6-925dcf1a1ac0.htm "ReadStringAsync 方法 ")

[ReadStringVariable 方法](../html/5dc9095d-86e8-d434-0027-83797a1d8c0a.htm "ReadStringVariable 方法 ")

[ReadStringVariableAsync 方法](../html/f850e762-bf8a-6ee7-c18c-3496d70244fe.htm "ReadStringVariableAsync 方法 ")

[ReadUFrame 方法](../html/100a1955-fc0a-871a-98f9-bb361c422a3f.htm "ReadUFrame 方法 ")

[ReadUFrameAsync 方法](../html/f63d72f5-eaf4-ae5d-e80a-d881f93a878e.htm "ReadUFrameAsync 方法 ")

[Reset 方法](../html/bf61707b-875c-4b42-2732-9b2fb69535ac.htm "Reset 方法 ")

[ResetAsync 方法](../html/f135af7f-0250-eb16-48fe-959dd0aeeb72.htm "ResetAsync 方法 ")

[SetMJ 方法](../html/9cd8c255-f908-6f4e-da07-33edc743dbf8.htm "SetMJ 方法 ")

[SetMJAsync 方法](../html/74b1cb5d-83f7-61b7-8eb0-0e10a675b964.htm "SetMJAsync 方法 ")

[Start 方法](../html/39764e6e-fde9-c99b-9034-99a64232aef6.htm "Start 方法 ")

[StartAsync 方法](../html/ddefed38-f8b8-4ae1-73e3-eecbfae54043.htm "StartAsync 方法 ")

[Svon 方法](../html/1f83c8cb-fea3-d99f-c7ad-f12f39bd9dee.htm "Svon 方法 ")

[SvonAsync 方法](../html/9b51697e-6164-280a-d37c-1cc16f467b2d.htm "SvonAsync 方法 ")

[ToString 方法](../html/1a5b666e-cae2-1664-908f-56866da693e5.htm "ToString 方法 ")

[Write 方法](../html/27720779-dc92-a6b4-c792-940ec2a10d22.htm "Write 方法 ")

[WriteAsync 方法](../html/e1646488-a935-0a5b-ec6a-4d04a3e69d38.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| YRC1000TcpNetMoveJAsync 方法 |

向指定的坐标系位置进行关节动作。其中没有外部轴的系统， 7-12外部轴的值设定为「0」  
Perform joint motions to the specified coordinate system position.
where there is no external axis system, the value of 7-12 external axis is set to "0"

**命名空间：**
 [HslCommunication.Robot.YASKAWA](3d67023f-78aa-c5d9-85b6-f18f5dfbc176.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public Task<OperateResult> MoveJAsync(
	YRCRobotData robotData
)
```

```
Public Function MoveJAsync ( 
	robotData As YRCRobotData
) As Task(Of OperateResult)
```

```
public:
Task<OperateResult^>^ MoveJAsync(
	YRCRobotData^ robotData
)
```

```
member MoveJAsync : 
        robotData : YRCRobotData -> Task<OperateResult> 
```

#### 参数

robotData
:   类型：[HslCommunication.Robot.YASKAWAYRCRobotData](233ef65f-4a8d-20e6-03ed-67c1cbefd2d6.htm)  
    机器的的数据信息

#### 返回值

类型：Task[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否动作成功

![](../icons/SectionExpanded.png)备注

其中形态数据由6个bool数组组成，每个bool含义参考参数说明，0表示 False，1表示 True

![](../icons/SectionExpanded.png)参见

#### 引用

[YRC1000TcpNet 类](cdcdc1cc-17cf-da27-7d20-ca6d6db1cd7e.htm)

[HslCommunication.Robot.YASKAWA 命名空间](3d67023f-78aa-c5d9-85b6-f18f5dfbc176.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## MSDP 方法 

[原文連結](http://api.hslcommunication.cn/html/1e78f8b5-7336-7fdf-a7e0-e40af1d3de7f.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.YASKAWA](../html/3d67023f-78aa-c5d9-85b6-f18f5dfbc176.htm "HslCommunication.Robot.YASKAWA")

[YRC1000TcpNet 类](../html/cdcdc1cc-17cf-da27-7d20-ca6d6db1cd7e.htm "YRC1000TcpNet 类")

[YRC1000TcpNet 方法](../html/7a26c354-eee0-b601-1d8a-51198e40e3fa.htm "YRC1000TcpNet 方法")

[Cancel 方法](../html/be66b6de-fe15-95b5-6969-70518f086388.htm "Cancel 方法 ")

[CancelAsync 方法](../html/fedb71cd-919c-c002-7556-551ae30fb10e.htm "CancelAsync 方法 ")

[Cycle 方法](../html/5a144fc7-5bc5-7250-6c13-826c172bd2ce.htm "Cycle 方法 ")

[CycleAsync 方法](../html/82bfe8a9-00b3-b34e-2a82-6d022becdf9a.htm "CycleAsync 方法 ")

[Delete 方法](../html/c7bfb022-46f4-8cdb-1b50-ffc1e5eb716b.htm "Delete 方法 ")

[DeleteAsync 方法](../html/c5a62bdd-3432-abcd-9d03-1ff7e30bd447.htm "DeleteAsync 方法 ")

[GetNewNetMessage 方法](../html/888815ce-da5c-0cb3-c6f9-800ec37a19ab.htm "GetNewNetMessage 方法 ")

[HLock 方法](../html/9c4ed34f-89fd-62e2-f707-08d1348c7541.htm "HLock 方法 ")

[HLockAsync 方法](../html/f0eb6982-a92b-4481-6435-33fa2b841b3c.htm "HLockAsync 方法 ")

[Hold 方法](../html/bf4ccc0a-1576-83ca-49e3-fde86b3a3738.htm "Hold 方法 ")

[HoldAsync 方法](../html/8c84a7d8-3677-0b6b-bdb5-44aa61ed1fcd.htm "HoldAsync 方法 ")

[InitializationOnConnect 方法](../html/5bd0274d-b3e7-f797-2dc0-c0dd2f330694.htm "InitializationOnConnect 方法 ")

[InitializationOnConnectAsync 方法](../html/e2f434c6-2ad9-30e8-8f2c-dcab95f5ec2b.htm "InitializationOnConnectAsync 方法 ")

[IORead 方法](../html/4539f381-2df3-4315-084c-7fcc81d523cb.htm "IORead 方法 ")

[IOReadAsync 方法](../html/5d550dc5-c02c-7245-122e-c3624e8d9c14.htm "IOReadAsync 方法 ")

[IOWrite 方法](../html/0dfecba5-f497-4b0f-55a5-339c7202d742.htm "IOWrite 方法 ")

[IOWriteAsync 方法](../html/64b5c202-29ff-0737-9c62-a4735bc00911.htm "IOWriteAsync 方法 ")

[JSeq 方法](../html/7057228b-aa47-3318-248d-1594588783bd.htm "JSeq 方法 ")

[JSeqAsync 方法](../html/c4eb7232-63e1-cd10-f89e-7f58072e3e9c.htm "JSeqAsync 方法 ")

[Mode 方法](../html/cabcebca-60ed-7689-e44a-e227936dcc51.htm "Mode 方法 ")

[ModeAsync 方法](../html/e0956a5d-5ae1-f884-52dd-0e91b85828b5.htm "ModeAsync 方法 ")

[MoveJ 方法](../html/10990eab-b2f9-e4b8-3c88-5bffc78aa34a.htm "MoveJ 方法 ")

[MoveJAsync 方法](../html/5b05904a-27a0-b4f9-adc7-1fe2cc5b14a1.htm "MoveJAsync 方法 ")

[MSDP 方法](../html/1e78f8b5-7336-7fdf-a7e0-e40af1d3de7f.htm "MSDP 方法 ")

[MSDPAsync 方法](../html/c2a18a60-0dae-112b-af16-0e028e6a203a.htm "MSDPAsync 方法 ")

[Read 方法](../html/69082cc9-c365-8fc4-36df-082c0c6c22c3.htm "Read 方法 ")

[ReadALARM 方法](../html/b56d10b7-3ebc-effd-c8f4-1dbf2d9f35a0.htm "ReadALARM 方法 ")

[ReadALARMAsync 方法](../html/6bd50e3f-6923-2c34-0e8d-d08ba82840f1.htm "ReadALARMAsync 方法 ")

[ReadAsync 方法](../html/185f22c2-a4f1-985e-dd8e-4bf7107c5ae7.htm "ReadAsync 方法 ")

[ReadByCommand 方法](../html/80dbd712-4782-b205-7e73-977c37f5a214.htm "ReadByCommand 方法 ")

[ReadByCommandAsync 方法](../html/2ac7651a-e36e-1635-0a61-b855722693a1.htm "ReadByCommandAsync 方法 ")

[ReadByteVariable 方法](../html/e22f2610-10f3-8472-92c2-3f9e904ebb2b.htm "ReadByteVariable 方法 ")

[ReadByteVariableAsync 方法](../html/d01b8a15-6079-9c57-6936-e0c569cf5194.htm "ReadByteVariableAsync 方法 ")

[ReadDoubleIntegerVariable 方法](../html/beab7b65-77db-3c62-6f95-ccf796ecac03.htm "ReadDoubleIntegerVariable 方法 ")

[ReadDoubleIntegerVariableAsync 方法](../html/c6845ca5-1d25-d613-89f9-28da2515cf24.htm "ReadDoubleIntegerVariableAsync 方法 ")

[ReadFromCoreServer 方法](../html/4e4be36a-8ce8-a41b-6ec4-4f82479c6f23.htm "ReadFromCoreServer 方法 ")

[ReadFromCoreServerAsync 方法](../html/e4cb7ebe-d42a-ea17-8972-b246a42e9a90.htm "ReadFromCoreServerAsync 方法 ")

[ReadIntegerVariable 方法](../html/df3171c4-e18b-a52c-9032-bd789e7388a0.htm "ReadIntegerVariable 方法 ")

[ReadIntegerVariableAsync 方法](../html/954c86a5-37e1-f32b-3a6e-7fc4274cbb6f.htm "ReadIntegerVariableAsync 方法 ")

[ReadJSeq 方法](../html/d7a72f10-2923-e045-8c71-e8b43c551c5a.htm "ReadJSeq 方法 ")

[ReadJSeqAsync 方法](../html/ddfc3321-88c7-9d2a-b71e-9315ce77864d.htm "ReadJSeqAsync 方法 ")

[ReadPOSC 方法](../html/fea57b6b-1d2f-0547-73de-6120ee5de705.htm "ReadPOSC 方法 ")

[ReadPOSCAsync 方法](../html/5f5731c4-a35a-affa-3d76-7b5f50fe63a1.htm "ReadPOSCAsync 方法 ")

[ReadPOSJ 方法](../html/14045444-1d70-3bdf-9848-a886c204da46.htm "ReadPOSJ 方法 ")

[ReadPOSJAsync 方法](../html/0ed6f125-5c82-6a30-1a4e-a6a773dda5d8.htm "ReadPOSJAsync 方法 ")

[ReadRealVariable 方法](../html/5f6c460a-0bb9-32fb-8aa8-33b99dc8be5e.htm "ReadRealVariable 方法 ")

[ReadRealVariableAsync 方法](../html/996ae64e-84f9-ccf6-fa2a-de4e4a1b13c9.htm "ReadRealVariableAsync 方法 ")

[ReadStats 方法](../html/3a24171c-f204-8c63-cdbc-8cf23a965b75.htm "ReadStats 方法 ")

[ReadStatsAsync 方法](../html/43a75b5c-c5af-e81f-286c-a4aa2a7a898b.htm "ReadStatsAsync 方法 ")

[ReadString 方法](../html/dbc09266-42e8-b939-aea5-540d528d19ee.htm "ReadString 方法 ")

[ReadStringAsync 方法](../html/8fa97613-eb8d-c1e1-33f6-925dcf1a1ac0.htm "ReadStringAsync 方法 ")

[ReadStringVariable 方法](../html/5dc9095d-86e8-d434-0027-83797a1d8c0a.htm "ReadStringVariable 方法 ")

[ReadStringVariableAsync 方法](../html/f850e762-bf8a-6ee7-c18c-3496d70244fe.htm "ReadStringVariableAsync 方法 ")

[ReadUFrame 方法](../html/100a1955-fc0a-871a-98f9-bb361c422a3f.htm "ReadUFrame 方法 ")

[ReadUFrameAsync 方法](../html/f63d72f5-eaf4-ae5d-e80a-d881f93a878e.htm "ReadUFrameAsync 方法 ")

[Reset 方法](../html/bf61707b-875c-4b42-2732-9b2fb69535ac.htm "Reset 方法 ")

[ResetAsync 方法](../html/f135af7f-0250-eb16-48fe-959dd0aeeb72.htm "ResetAsync 方法 ")

[SetMJ 方法](../html/9cd8c255-f908-6f4e-da07-33edc743dbf8.htm "SetMJ 方法 ")

[SetMJAsync 方法](../html/74b1cb5d-83f7-61b7-8eb0-0e10a675b964.htm "SetMJAsync 方法 ")

[Start 方法](../html/39764e6e-fde9-c99b-9034-99a64232aef6.htm "Start 方法 ")

[StartAsync 方法](../html/ddefed38-f8b8-4ae1-73e3-eecbfae54043.htm "StartAsync 方法 ")

[Svon 方法](../html/1f83c8cb-fea3-d99f-c7ad-f12f39bd9dee.htm "Svon 方法 ")

[SvonAsync 方法](../html/9b51697e-6164-280a-d37c-1cc16f467b2d.htm "SvonAsync 方法 ")

[ToString 方法](../html/1a5b666e-cae2-1664-908f-56866da693e5.htm "ToString 方法 ")

[Write 方法](../html/27720779-dc92-a6b4-c792-940ec2a10d22.htm "Write 方法 ")

[WriteAsync 方法](../html/e1646488-a935-0a5b-ec6a-4d04a3e69d38.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| YRC1000TcpNetMSDP 方法 |

接受消息数据时， 在YRC1000的示教编程器的远程画面下显示消息若。若不是远程画面时，强制切换到远程画面。显示MDSP命令的消息。  
When receiving message data, a message is displayed on the remote screen of the YRC1000 programming pendant.
If it is not a remote screen, it is forced to switch to the remote screen. Display the message of the MDSP command.

**命名空间：**
 [HslCommunication.Robot.YASKAWA](3d67023f-78aa-c5d9-85b6-f18f5dfbc176.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public OperateResult MSDP(
	string message
)
```

```
Public Function MSDP ( 
	message As String
) As OperateResult
```

```
public:
OperateResult^ MSDP(
	String^ message
)
```

```
member MSDP : 
        message : string -> OperateResult 
```

#### 参数

message
:   类型：SystemString  
    显示信息（最大 30byte 字符串）

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否显示成功

![](../icons/SectionExpanded.png)参见

#### 引用

[YRC1000TcpNet 类](cdcdc1cc-17cf-da27-7d20-ca6d6db1cd7e.htm)

[HslCommunication.Robot.YASKAWA 命名空间](3d67023f-78aa-c5d9-85b6-f18f5dfbc176.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)