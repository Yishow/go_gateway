# HslCommunication - HslCommunication.Profinet.Sick

> 分類頁數: 10



---
## HslCommunication.Profinet.Sick

[原文連結](http://api.hslcommunication.cn/html/0dd7d6fd-55ea-53e6-68fc-8aa45401ad37.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Sick](../html/0dd7d6fd-55ea-53e6-68fc-8aa45401ad37.htm "HslCommunication.Profinet.Sick")

[SickIcrTcpServer 类](../html/54f2dac1-62dc-1efb-aa3c-68b0d45a33cc.htm "SickIcrTcpServer 类")

[SickIcrTcpServer.ReceivedBarCodeDelegate 委托](../html/121fcfd3-7f91-ed26-cccf-98db3fa02780.htm "SickIcrTcpServer.ReceivedBarCodeDelegate 委托")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslCommunication.Profinet.Sick 命名空间 |

[缺少 "N:HslCommunication.Profinet.Sick" 的 <summary> 文档]

![](../icons/SectionExpanded.png)类

|  | 类 | 说明 |
| --- | --- | --- |
| 公共类 | [SickIcrTcpServer](54f2dac1-62dc-1efb-aa3c-68b0d45a33cc.htm) | Sick的扫码器的服务器信息，只要启动服务器之后，扫码器配置将条码发送到PC的指定端口上来即可，就可以持续的接收条码信息，同样也适用于海康，基恩士，DATELOGIC 。  The server information of Sick's code scanner, as long as the server is started, the code scanner is configured to send the barcode to the designated port of the PC, and it can continuously receive the barcode information. |

![](../icons/SectionExpanded.png)委托

|  | 委托 | 说明 |
| --- | --- | --- |
| 公共委托 | [SickIcrTcpServerReceivedBarCodeDelegate](121fcfd3-7f91-ed26-cccf-98db3fa02780.htm) | 接收条码数据的委托信息  Entrusted information to receive barcode data |

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## SickIcrTcpServer 类

[原文連結](http://api.hslcommunication.cn/html/54f2dac1-62dc-1efb-aa3c-68b0d45a33cc.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Sick](../html/0dd7d6fd-55ea-53e6-68fc-8aa45401ad37.htm "HslCommunication.Profinet.Sick")

[SickIcrTcpServer 类](../html/54f2dac1-62dc-1efb-aa3c-68b0d45a33cc.htm "SickIcrTcpServer 类")

[SickIcrTcpServer 构造函数](../html/60dfe1fa-3a09-7b39-b644-e5f58760b8d3.htm "SickIcrTcpServer 构造函数 ")

[SickIcrTcpServer 属性](../html/fe9ef976-5561-aeb7-3156-02ae5bd248b1.htm "SickIcrTcpServer 属性")

[SickIcrTcpServer 方法](../html/bdb23ae3-79a4-5717-4c42-e29fb17e11e5.htm "SickIcrTcpServer 方法")

[SickIcrTcpServer 事件](../html/375747ce-c440-cf51-e045-2bb8b064ff64.htm "SickIcrTcpServer 事件")

[SickIcrTcpServer 字段](../html/fff1c148-1845-faa4-5860-217826955f40.htm "SickIcrTcpServer 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| SickIcrTcpServer 类 |

Sick的扫码器的服务器信息，只要启动服务器之后，扫码器配置将条码发送到PC的指定端口上来即可，就可以持续的接收条码信息，同样也适用于海康，基恩士，DATELOGIC 。  
The server information of Sick's code scanner, as long as the server is started, the code scanner is configured to send the barcode to the designated port of the PC, and it can continuously receive the barcode information.

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  [HslCommunication.Core.NetCommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)  
    [HslCommunication.Core.NetCommunicationServer](5a3cb700-c0f6-e35d-6807-dfe8bfecf747.htm)  
      HslCommunication.Profinet.SickSickIcrTcpServer

**命名空间：**
 [HslCommunication.Profinet.Sick](0dd7d6fd-55ea-53e6-68fc-8aa45401ad37.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class SickIcrTcpServer : CommunicationServer
```

```
Public Class SickIcrTcpServer
	Inherits CommunicationServer
```

```
public ref class SickIcrTcpServer : public CommunicationServer
```

```
type SickIcrTcpServer =  
    class
        inherit CommunicationServer
    end
```

SickIcrTcpServer 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [SickIcrTcpServer](60dfe1fa-3a09-7b39-b644-e5f58760b8d3.htm) | 实例化一个默认的服务器对象  Instantiate a default server object |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [CheckSerialDataComplete](bef5c4d1-8fac-1f94-42ab-155ad17f76d5.htm) | 检查串口接收到的数据是否完整 (继承自 [CommunicationServer](5a3cb700-c0f6-e35d-6807-dfe8bfecf747.htm)。) |
| 公共属性 | [CreateNewMessage](19bf00bd-dfa5-8afb-6716-a985f22da99a.htm) | 获取一个新的消息对象的方法，需要在继承类里面进行重写  The method to get a new message object needs to be overridden in the inheritance class (继承自 [CommunicationServer](5a3cb700-c0f6-e35d-6807-dfe8bfecf747.htm)。) |
| 公共属性 | [CreatePipeSession](5bebf8f4-912a-0e90-f877-7ae0fe9f02ed.htm) | 创建会话状态的委托对象，也就可以自己指定创建自定义的会话 (继承自 [CommunicationServer](5a3cb700-c0f6-e35d-6807-dfe8bfecf747.htm)。) |
| 公共属性 | [EnableIPv6](f15cbcb3-5dcf-a658-0a93-3f0921383c32.htm) | 获取或设置服务器是否支持IPv6的地址协议信息  Get or set whether the server supports IPv6 address protocol information (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |
| 公共属性 | [ForceSerialReceiveOnce](4e1ee30d-f0c6-72fd-161a-83cc6c2e8e5e.htm) | 获取或设置当前的服务器接收串口数据时候，是否强制只接收一次数据，默认为false，适合点对点通信，如果你总线形式的连接，则需要设置 True  Get or set whether to force the data to be received only once when the current server receives serial port data. The default value is false, which is suitable for point-to-point communication. If you have a bus connection, you need to set True (继承自 [CommunicationServer](5a3cb700-c0f6-e35d-6807-dfe8bfecf747.htm)。) |
| 公共属性 | [IsStarted](8d6786a8-d416-4f45-158c-de6be0634aca.htm) | 服务器引擎是否启动  Whether the server engine is started (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |
| 公共属性 | [IsUseSSL](f4daa1e8-f433-0ee2-8284-16c4754eba25.htm) | 获取当前的服务器是否使用了SSL证书功能 (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |
| 公共属性 | [LocalAddress](c0effb49-cfe6-3267-9aef-afdab27a41c3.htm) | 获取或设置服务器绑定的本地IP地址，默认为空，使用本地所有可用的ip地址  Obtain or set the local IP address bound to the server, which is empty by default, and use all available IP addresses locally (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |
| 公共属性 | [LogDebugMessage](8a6db355-1130-2af8-cc6a-aeccc1ac0b29.htm) | 记录一些调试日志的委托，将会进行输出调试文本。  The delegate that records some debug logs will output debug text. (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |
| 公共属性代码示例 | [LogNet](746b148d-1ee3-1ee7-ee5d-3938eb72b2df.htm) | 组件的日志工具，支持日志记录，只要实例化后，当前网络的基本信息，就以[DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm)等级进行输出  The component's logging tool supports logging. As long as the instantiation of the basic network information, the output will be output at [DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm) (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |
| 公共属性 | [Port](6f700f7a-bd7e-c4b3-1275-a30491469e8f.htm) | 获取或设置服务器的端口号，如果是设置，需要在服务器启动前设置完成，才能生效。  Gets or sets the port number of the server. If it is set, it needs to be set before the server starts to take effect. (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |
| 公共属性 | [SerialReceiveAtleastTime](5802707e-1dd6-f7f9-e293-dd387ae817fd.htm) | 获取或设置串口模式下，接收一条数据最短的时间要求，当设备发送的数据非常慢的时候，或是分割发送数据的时候，就需要将本值设置的大一点，默认为20ms  Get or set the shortest time required to receive a piece of data in serial port mode. When the data sent by the device is very slow, or when the data is divided and sent, you need to set this value to a larger value, the default is 20ms (继承自 [CommunicationServer](5a3cb700-c0f6-e35d-6807-dfe8bfecf747.htm)。) |
| 公共属性 | [SessionsMax](c65bac2d-a238-f93d-02f2-a7ed1297cf7a.htm) | 获取或设置当前允许登录的最大客户端数量，默认为 uint.MaxValue = 4294967295  Gets or sets the maximum number of clients that are currently allowed to log in, which defaults to uint.MaxValue = 4294967295 (继承自 [CommunicationServer](5a3cb700-c0f6-e35d-6807-dfe8bfecf747.htm)。) |
| 公共属性 | [SocketKeepAliveTime](7d77e3df-9bbe-1a6e-a051-c87d9516c3a2.htm) | 获取或设置客户端的Socket的心跳时间信息，这个是Socket底层自动实现的心跳包，不基于协议层实现。默认小于0，不开启心跳检测，如果需要开启，设置 60\_000 比较合适，单位毫秒  Get or set the heartbeat time information of the Socket of the client. This is the heartbeat packet automatically implemented by the bottom layer of the Socket, not based on the protocol layer. The default value is less than 0, and heartbeat detection is not enabled. If you need to enable it, it is more appropriate to set 60\_000, in milliseconds. (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |
| 公共属性 | [ThreadPoolLoginAfterClientCheck](e926b0f9-3849-07d2-ed1d-0a249fad29ae.htm) | 当线程检查后，进行登录之前的检查，通常用于自定义的握手包校验操作。仅对TCP通信的时候有效。 (继承自 [CommunicationServer](5a3cb700-c0f6-e35d-6807-dfe8bfecf747.htm)。) |
| 公共属性 | [UdpBufferSize](762072d5-f18f-8422-fc5f-19fb0eae7eff.htm) | 获取或设置一次接收时的数据长度，默认2KB数据长度，特殊情况的时候需要调整  Gets or sets the length of data received at a time. The default length is 2KB (继承自 [CommunicationServer](5a3cb700-c0f6-e35d-6807-dfe8bfecf747.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [AddSession](22a49542-4dd7-d360-ed5d-85266665b102.htm) | 新增加一个管道会话信息  A new pipeline session information has been added (继承自 [CommunicationServer](5a3cb700-c0f6-e35d-6807-dfe8bfecf747.htm)。) |
| 受保护的方法 | [CheckSerialReceiveDataComplete](a66e7f2a-c605-d949-b1ec-9271b75df1b6.htm) | 检查串口接收的数据是否完成的方法，如果接收完成，则返回True (继承自 [CommunicationServer](5a3cb700-c0f6-e35d-6807-dfe8bfecf747.htm)。) |
| 公共方法 | [CloseSerialSlave](f5a8b6f0-0481-b9b3-c31b-90d9a945b416.htm) | 关闭提供从机服务的串口对象  Close the serial port object that provides slave services (继承自 [CommunicationServer](5a3cb700-c0f6-e35d-6807-dfe8bfecf747.htm)。) |
| 公共方法 | [ConnectHslAlientClient](b98a142a-9a96-4cc8-4fdc-704dbfc7b438.htm) | 创建一个指定的异形客户端连接，使用Hsl协议来发送注册包  Create a specified profiled client connection and use the Hsl protocol to send registration packets (继承自 [CommunicationServer](5a3cb700-c0f6-e35d-6807-dfe8bfecf747.htm)。) |
| 公共方法 | [ConnectRemoteServer](5946c84a-5591-14d2-2c75-42b95914959e.htm) | 新增一个主动连接的请求，将不会收到是否连接成功的信息，当网络中断及奔溃之后，会自动重新连接。  A new active connection request will not receive a message whether the connection is successful. When the network is interrupted and crashed, it will automatically reconnect. (继承自 [CommunicationServer](5a3cb700-c0f6-e35d-6807-dfe8bfecf747.htm)。) |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | [ExtraOnClose](ede241ce-62f0-59aa-45b0-1cc5622a6acf.htm) | 关闭的时候额外执行的功能代码 (继承自 [CommunicationServer](5a3cb700-c0f6-e35d-6807-dfe8bfecf747.htm)。) |
| 受保护的方法 | [ExtraOnStart](2397069b-cc49-1ba5-41c6-ffe84e25219c.htm) | 服务器启动的时候额外执行的功能代码 (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 受保护的方法 | [GetNewNetMessage](89004772-77de-0bf5-785b-46d2a4ce3a3d.htm) | 获取一个新的消息对象的方法，需要在继承类里面进行重写  The method to get a new message object needs to be overridden in the inheritance class (继承自 [CommunicationServer](5a3cb700-c0f6-e35d-6807-dfe8bfecf747.htm)。) |
| 公共方法 | [GetPipeSessions](fbbc274f-38c2-ad5a-dadf-93ac5a7eacff.htm) | 获取管道会话的列表  Get a list of pipeline sessions (继承自 [CommunicationServer](5a3cb700-c0f6-e35d-6807-dfe8bfecf747.htm)。) |
| 公共方法 | [GetTrustedClients](1a468d3b-7c50-7a36-f1f9-76a1d1488a8a.htm) | 获取受信任的客户端列表  Get a list of trusted clients (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法 | [LogDebugMsg](3dc4b43f-2fdc-1d8e-ef40-66319b9f7055.htm) | 记录当前的日志信息 (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法 | [RemoveSession(TimeSpan)](863eb441-97a5-7c9a-2545-0665cd0e46b8.htm) | 指定超时时间移除当前的会话列表，只有是TCP的管道（[PipeTcpNet](26b94bd5-2695-36e7-6378-c3d41846be8f.htm)）才需要被移除。  Specify a timeout to remove the current session list. Only TCP pipe ([PipeTcpNet](26b94bd5-2695-36e7-6378-c3d41846be8f.htm)) need to be removed. (继承自 [CommunicationServer](5a3cb700-c0f6-e35d-6807-dfe8bfecf747.htm)。) |
| 公共方法 | [RemoveSession(PipeSession, String)](92e9b2f6-96c3-b586-b554-f8ecd73df580.htm) | 移除一个管道会话  Remove a pipeline session (继承自 [CommunicationServer](5a3cb700-c0f6-e35d-6807-dfe8bfecf747.htm)。) |
| 公共方法 | [ServerClose](04a70d3a-d427-ed8d-53cf-3d3be56bdbd3.htm) | 关闭服务器的引擎  Shut down the server's engine (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |
| 公共方法 | [ServerStart](6aa1f94d-768c-4711-80e5-a529a5040f98.htm) | 指定端口号来启动服务器的引擎  Specify the port number to start the server's engine (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |
| 公共方法 | [ServerStart(Int32)](06e14b74-3deb-8e46-81c6-688f078375d9.htm) | 指定端口号来启动服务器的引擎  Specify the port number to start the server's engine (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |
| 公共方法 | [ServerStart(Int32, Boolean)](af75f437-3c8e-ac5c-96ab-96ea7bb4cdb4.htm) | 指定端口号来启动服务器的引擎  Specify the port number to start the server's engine (继承自 [CommunicationServer](5a3cb700-c0f6-e35d-6807-dfe8bfecf747.htm)。) |
| 公共方法 | [ServerStart(Int32, Int32)](1feb9368-f6bf-85e6-122a-d0c0d5a67bed.htm) | 指定一个TCP端口及UDP端口，同时启动两种模式的服务器  Specify a TCP port and a UDP port to start the server in both modes at the same time (继承自 [CommunicationServer](5a3cb700-c0f6-e35d-6807-dfe8bfecf747.htm)。) |
| 公共方法 | [SetNetMessage](836e8de1-0461-348c-b56b-2765e37c0ec5.htm) | 设置当前的服务器接收的消息信息 (继承自 [CommunicationServer](5a3cb700-c0f6-e35d-6807-dfe8bfecf747.htm)。) |
| 公共方法 | [SetSslPipeAction](5ecd4347-5dae-126d-f43c-76d4c7390705.htm) | 设置一个SSL的管道操作对象，在管道实例化之后，可以进行一些初始化的属性设置，例如自定义 SslProtocols 枚举 (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |
| 公共方法 | [SetTrustedIpAddress](6204b3e7-af37-7c20-98db-08114431b5ac.htm) | 设置并启动受信任的客户端登录并读写，如果为null，将关闭对客户端的ip验证  Set and start the trusted client login and read and write, if it is null, the client's IP verification will be turned off (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |
| 受保护的方法 | [SocketAcceptExtraCheck](3ed24041-5b60-57a7-7629-f7c49baf8bf9.htm) | 当客户端的socket登录的时候额外检查的操作，并返回操作的结果信息。  The operation is additionally checked when the client's socket logs in, and the result information of the operation is returned. (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |
| 公共方法 | [StartSerialSlave(ActionSerialPort)](dde390ea-45f6-bc3f-be69-3902f6a5641b.htm) | 启动串口的从机服务，使用自定义的初始化方法初始化串口的参数  Start the slave service of serial and initialize the parameters of the serial port using a custom initialization method (继承自 [CommunicationServer](5a3cb700-c0f6-e35d-6807-dfe8bfecf747.htm)。) |
| 公共方法 | [StartSerialSlave(String)](cfed6d77-50f0-eb64-37d0-7c466cc2d7bf.htm) | 启动串口的从机服务，使用默认的参数进行初始化串口，9600波特率，8位数据位，无奇偶校验，1位停止位  Start the slave service of serial, initialize the serial port with default parameters, 9600 baud rate, 8 data bits, no parity, 1 stop bit (继承自 [CommunicationServer](5a3cb700-c0f6-e35d-6807-dfe8bfecf747.htm)。) |
| 公共方法 | [StartSerialSlave(String, Int32)](155bfb40-a6bc-8ca0-3ebc-14307221ba59.htm) | 启动串口的从机服务，使用默认的参数进行初始化串口，8位数据位，无奇偶校验，1位停止位  Start the slave service of serial, initialize the serial port with default parameters, 8 data bits, no parity, 1 stop bit (继承自 [CommunicationServer](5a3cb700-c0f6-e35d-6807-dfe8bfecf747.htm)。) |
| 公共方法 | [StartSerialSlave(String, Int32, Int32, Parity, StopBits)](7045bee6-9c7f-65e0-191a-205c360dd42c.htm) | 启动串口的从机服务，使用指定的参数进行初始化串口，指定数据位，指定奇偶校验，指定停止位 (继承自 [CommunicationServer](5a3cb700-c0f6-e35d-6807-dfe8bfecf747.htm)。) |
| 受保护的方法 | [ThreadPoolLogin](7ebe2e21-f51e-3a08-a499-624462a0fede.htm) | 当客户端连接到服务器，并听过额外的检查后，进行回调的方法  Callback method when the client connects to the server and has heard additional checks (继承自 [CommunicationServer](5a3cb700-c0f6-e35d-6807-dfe8bfecf747.htm)。) |
| 公共方法 | [ToString](f4649b7b-3f2c-6b17-214a-ecde54810eab.htm) | (重写 [CommunicationServerToString](955ea032-e1fd-11ff-9ec8-72cb0572d36a.htm).) |
| 公共方法 | [UseSSL(X509Certificate)](3577f66a-4337-6d6c-0f73-1f904f6bed38.htm) | 使用SSL通信，传递一个证书的对象 (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |
| 公共方法 | [UseSSL(String, String)](e6f37bac-c743-e05a-729a-0e4a64daf486.htm) | 使用SSL通信，传递一个证书的路径，以及证书的密码 (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)事件

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共事件 | [OnClientOffline](d79ec6d7-e08e-9afe-3cc0-3a34f3f07017.htm) | 当客户端下线时候的触发的事件  Event triggered when the client goes offline (继承自 [CommunicationServer](5a3cb700-c0f6-e35d-6807-dfe8bfecf747.htm)。) |
| 公共事件 | [OnClientOnline](d60efb73-a438-eca9-3cf4-7bc26a24e165.htm) | 当客户端上线时候的触发的事件  Event triggered when the client goes online (继承自 [CommunicationServer](5a3cb700-c0f6-e35d-6807-dfe8bfecf747.htm)。) |
| 公共事件 | [OnPipeMessageReceived](539eb796-8339-4d86-7fa2-c9b6a50f533f.htm) | 当管道接收到消息时触发 (继承自 [CommunicationServer](5a3cb700-c0f6-e35d-6807-dfe8bfecf747.htm)。) |
| 公共事件 | [OnReceivedBarCode](20f767de-bcf2-0b81-e456-149f87ecaa92.htm) | 当接收到条码数据的时候触发  Triggered when barcode data is received |

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

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.Profinet.Sick 命名空间](0dd7d6fd-55ea-53e6-68fc-8aa45401ad37.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## SickIcrTcpServer 构造函数 

[原文連結](http://api.hslcommunication.cn/html/60dfe1fa-3a09-7b39-b644-e5f58760b8d3.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Sick](../html/0dd7d6fd-55ea-53e6-68fc-8aa45401ad37.htm "HslCommunication.Profinet.Sick")

[SickIcrTcpServer 类](../html/54f2dac1-62dc-1efb-aa3c-68b0d45a33cc.htm "SickIcrTcpServer 类")

[SickIcrTcpServer 构造函数](../html/60dfe1fa-3a09-7b39-b644-e5f58760b8d3.htm "SickIcrTcpServer 构造函数 ")

[SickIcrTcpServer 属性](../html/fe9ef976-5561-aeb7-3156-02ae5bd248b1.htm "SickIcrTcpServer 属性")

[SickIcrTcpServer 方法](../html/bdb23ae3-79a4-5717-4c42-e29fb17e11e5.htm "SickIcrTcpServer 方法")

[SickIcrTcpServer 事件](../html/375747ce-c440-cf51-e045-2bb8b064ff64.htm "SickIcrTcpServer 事件")

[SickIcrTcpServer 字段](../html/fff1c148-1845-faa4-5860-217826955f40.htm "SickIcrTcpServer 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| SickIcrTcpServer 构造函数 |

实例化一个默认的服务器对象  
Instantiate a default server object

**命名空间：**
 [HslCommunication.Profinet.Sick](0dd7d6fd-55ea-53e6-68fc-8aa45401ad37.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public SickIcrTcpServer()
```

```
Public Sub New
```

```
public:
SickIcrTcpServer()
```

```
new : unit -> SickIcrTcpServer
```

![](../icons/SectionExpanded.png)参见

#### 引用

[SickIcrTcpServer 类](54f2dac1-62dc-1efb-aa3c-68b0d45a33cc.htm)

[HslCommunication.Profinet.Sick 命名空间](0dd7d6fd-55ea-53e6-68fc-8aa45401ad37.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## SickIcrTcpServer 属性

[原文連結](http://api.hslcommunication.cn/html/fe9ef976-5561-aeb7-3156-02ae5bd248b1.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Sick](../html/0dd7d6fd-55ea-53e6-68fc-8aa45401ad37.htm "HslCommunication.Profinet.Sick")

[SickIcrTcpServer 类](../html/54f2dac1-62dc-1efb-aa3c-68b0d45a33cc.htm "SickIcrTcpServer 类")

[SickIcrTcpServer 构造函数](../html/60dfe1fa-3a09-7b39-b644-e5f58760b8d3.htm "SickIcrTcpServer 构造函数 ")

[SickIcrTcpServer 属性](../html/fe9ef976-5561-aeb7-3156-02ae5bd248b1.htm "SickIcrTcpServer 属性")

[SickIcrTcpServer 方法](../html/bdb23ae3-79a4-5717-4c42-e29fb17e11e5.htm "SickIcrTcpServer 方法")

[SickIcrTcpServer 事件](../html/375747ce-c440-cf51-e045-2bb8b064ff64.htm "SickIcrTcpServer 事件")

[SickIcrTcpServer 字段](../html/fff1c148-1845-faa4-5860-217826955f40.htm "SickIcrTcpServer 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| SickIcrTcpServer 属性 |

[SickIcrTcpServer](54f2dac1-62dc-1efb-aa3c-68b0d45a33cc.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [CheckSerialDataComplete](bef5c4d1-8fac-1f94-42ab-155ad17f76d5.htm) | 检查串口接收到的数据是否完整 (继承自 [CommunicationServer](5a3cb700-c0f6-e35d-6807-dfe8bfecf747.htm)。) |
| 公共属性 | [CreateNewMessage](19bf00bd-dfa5-8afb-6716-a985f22da99a.htm) | 获取一个新的消息对象的方法，需要在继承类里面进行重写  The method to get a new message object needs to be overridden in the inheritance class (继承自 [CommunicationServer](5a3cb700-c0f6-e35d-6807-dfe8bfecf747.htm)。) |
| 公共属性 | [CreatePipeSession](5bebf8f4-912a-0e90-f877-7ae0fe9f02ed.htm) | 创建会话状态的委托对象，也就可以自己指定创建自定义的会话 (继承自 [CommunicationServer](5a3cb700-c0f6-e35d-6807-dfe8bfecf747.htm)。) |
| 公共属性 | [EnableIPv6](f15cbcb3-5dcf-a658-0a93-3f0921383c32.htm) | 获取或设置服务器是否支持IPv6的地址协议信息  Get or set whether the server supports IPv6 address protocol information (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |
| 公共属性 | [ForceSerialReceiveOnce](4e1ee30d-f0c6-72fd-161a-83cc6c2e8e5e.htm) | 获取或设置当前的服务器接收串口数据时候，是否强制只接收一次数据，默认为false，适合点对点通信，如果你总线形式的连接，则需要设置 True  Get or set whether to force the data to be received only once when the current server receives serial port data. The default value is false, which is suitable for point-to-point communication. If you have a bus connection, you need to set True (继承自 [CommunicationServer](5a3cb700-c0f6-e35d-6807-dfe8bfecf747.htm)。) |
| 公共属性 | [IsStarted](8d6786a8-d416-4f45-158c-de6be0634aca.htm) | 服务器引擎是否启动  Whether the server engine is started (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |
| 公共属性 | [IsUseSSL](f4daa1e8-f433-0ee2-8284-16c4754eba25.htm) | 获取当前的服务器是否使用了SSL证书功能 (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |
| 公共属性 | [LocalAddress](c0effb49-cfe6-3267-9aef-afdab27a41c3.htm) | 获取或设置服务器绑定的本地IP地址，默认为空，使用本地所有可用的ip地址  Obtain or set the local IP address bound to the server, which is empty by default, and use all available IP addresses locally (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |
| 公共属性 | [LogDebugMessage](8a6db355-1130-2af8-cc6a-aeccc1ac0b29.htm) | 记录一些调试日志的委托，将会进行输出调试文本。  The delegate that records some debug logs will output debug text. (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |
| 公共属性代码示例 | [LogNet](746b148d-1ee3-1ee7-ee5d-3938eb72b2df.htm) | 组件的日志工具，支持日志记录，只要实例化后，当前网络的基本信息，就以[DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm)等级进行输出  The component's logging tool supports logging. As long as the instantiation of the basic network information, the output will be output at [DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm) (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |
| 公共属性 | [Port](6f700f7a-bd7e-c4b3-1275-a30491469e8f.htm) | 获取或设置服务器的端口号，如果是设置，需要在服务器启动前设置完成，才能生效。  Gets or sets the port number of the server. If it is set, it needs to be set before the server starts to take effect. (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |
| 公共属性 | [SerialReceiveAtleastTime](5802707e-1dd6-f7f9-e293-dd387ae817fd.htm) | 获取或设置串口模式下，接收一条数据最短的时间要求，当设备发送的数据非常慢的时候，或是分割发送数据的时候，就需要将本值设置的大一点，默认为20ms  Get or set the shortest time required to receive a piece of data in serial port mode. When the data sent by the device is very slow, or when the data is divided and sent, you need to set this value to a larger value, the default is 20ms (继承自 [CommunicationServer](5a3cb700-c0f6-e35d-6807-dfe8bfecf747.htm)。) |
| 公共属性 | [SessionsMax](c65bac2d-a238-f93d-02f2-a7ed1297cf7a.htm) | 获取或设置当前允许登录的最大客户端数量，默认为 uint.MaxValue = 4294967295  Gets or sets the maximum number of clients that are currently allowed to log in, which defaults to uint.MaxValue = 4294967295 (继承自 [CommunicationServer](5a3cb700-c0f6-e35d-6807-dfe8bfecf747.htm)。) |
| 公共属性 | [SocketKeepAliveTime](7d77e3df-9bbe-1a6e-a051-c87d9516c3a2.htm) | 获取或设置客户端的Socket的心跳时间信息，这个是Socket底层自动实现的心跳包，不基于协议层实现。默认小于0，不开启心跳检测，如果需要开启，设置 60\_000 比较合适，单位毫秒  Get or set the heartbeat time information of the Socket of the client. This is the heartbeat packet automatically implemented by the bottom layer of the Socket, not based on the protocol layer. The default value is less than 0, and heartbeat detection is not enabled. If you need to enable it, it is more appropriate to set 60\_000, in milliseconds. (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |
| 公共属性 | [ThreadPoolLoginAfterClientCheck](e926b0f9-3849-07d2-ed1d-0a249fad29ae.htm) | 当线程检查后，进行登录之前的检查，通常用于自定义的握手包校验操作。仅对TCP通信的时候有效。 (继承自 [CommunicationServer](5a3cb700-c0f6-e35d-6807-dfe8bfecf747.htm)。) |
| 公共属性 | [UdpBufferSize](762072d5-f18f-8422-fc5f-19fb0eae7eff.htm) | 获取或设置一次接收时的数据长度，默认2KB数据长度，特殊情况的时候需要调整  Gets or sets the length of data received at a time. The default length is 2KB (继承自 [CommunicationServer](5a3cb700-c0f6-e35d-6807-dfe8bfecf747.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[SickIcrTcpServer 类](54f2dac1-62dc-1efb-aa3c-68b0d45a33cc.htm)

[HslCommunication.Profinet.Sick 命名空间](0dd7d6fd-55ea-53e6-68fc-8aa45401ad37.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## SickIcrTcpServer 方法

[原文連結](http://api.hslcommunication.cn/html/bdb23ae3-79a4-5717-4c42-e29fb17e11e5.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Sick](../html/0dd7d6fd-55ea-53e6-68fc-8aa45401ad37.htm "HslCommunication.Profinet.Sick")

[SickIcrTcpServer 类](../html/54f2dac1-62dc-1efb-aa3c-68b0d45a33cc.htm "SickIcrTcpServer 类")

[SickIcrTcpServer 方法](../html/bdb23ae3-79a4-5717-4c42-e29fb17e11e5.htm "SickIcrTcpServer 方法")

[ToString 方法](../html/f4649b7b-3f2c-6b17-214a-ecde54810eab.htm "ToString 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| SickIcrTcpServer 方法 |

[SickIcrTcpServer](54f2dac1-62dc-1efb-aa3c-68b0d45a33cc.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [AddSession](22a49542-4dd7-d360-ed5d-85266665b102.htm) | 新增加一个管道会话信息  A new pipeline session information has been added (继承自 [CommunicationServer](5a3cb700-c0f6-e35d-6807-dfe8bfecf747.htm)。) |
| 受保护的方法 | [CheckSerialReceiveDataComplete](a66e7f2a-c605-d949-b1ec-9271b75df1b6.htm) | 检查串口接收的数据是否完成的方法，如果接收完成，则返回True (继承自 [CommunicationServer](5a3cb700-c0f6-e35d-6807-dfe8bfecf747.htm)。) |
| 公共方法 | [CloseSerialSlave](f5a8b6f0-0481-b9b3-c31b-90d9a945b416.htm) | 关闭提供从机服务的串口对象  Close the serial port object that provides slave services (继承自 [CommunicationServer](5a3cb700-c0f6-e35d-6807-dfe8bfecf747.htm)。) |
| 公共方法 | [ConnectHslAlientClient](b98a142a-9a96-4cc8-4fdc-704dbfc7b438.htm) | 创建一个指定的异形客户端连接，使用Hsl协议来发送注册包  Create a specified profiled client connection and use the Hsl protocol to send registration packets (继承自 [CommunicationServer](5a3cb700-c0f6-e35d-6807-dfe8bfecf747.htm)。) |
| 公共方法 | [ConnectRemoteServer](5946c84a-5591-14d2-2c75-42b95914959e.htm) | 新增一个主动连接的请求，将不会收到是否连接成功的信息，当网络中断及奔溃之后，会自动重新连接。  A new active connection request will not receive a message whether the connection is successful. When the network is interrupted and crashed, it will automatically reconnect. (继承自 [CommunicationServer](5a3cb700-c0f6-e35d-6807-dfe8bfecf747.htm)。) |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | [ExtraOnClose](ede241ce-62f0-59aa-45b0-1cc5622a6acf.htm) | 关闭的时候额外执行的功能代码 (继承自 [CommunicationServer](5a3cb700-c0f6-e35d-6807-dfe8bfecf747.htm)。) |
| 受保护的方法 | [ExtraOnStart](2397069b-cc49-1ba5-41c6-ffe84e25219c.htm) | 服务器启动的时候额外执行的功能代码 (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 受保护的方法 | [GetNewNetMessage](89004772-77de-0bf5-785b-46d2a4ce3a3d.htm) | 获取一个新的消息对象的方法，需要在继承类里面进行重写  The method to get a new message object needs to be overridden in the inheritance class (继承自 [CommunicationServer](5a3cb700-c0f6-e35d-6807-dfe8bfecf747.htm)。) |
| 公共方法 | [GetPipeSessions](fbbc274f-38c2-ad5a-dadf-93ac5a7eacff.htm) | 获取管道会话的列表  Get a list of pipeline sessions (继承自 [CommunicationServer](5a3cb700-c0f6-e35d-6807-dfe8bfecf747.htm)。) |
| 公共方法 | [GetTrustedClients](1a468d3b-7c50-7a36-f1f9-76a1d1488a8a.htm) | 获取受信任的客户端列表  Get a list of trusted clients (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法 | [LogDebugMsg](3dc4b43f-2fdc-1d8e-ef40-66319b9f7055.htm) | 记录当前的日志信息 (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法 | [RemoveSession(TimeSpan)](863eb441-97a5-7c9a-2545-0665cd0e46b8.htm) | 指定超时时间移除当前的会话列表，只有是TCP的管道（[PipeTcpNet](26b94bd5-2695-36e7-6378-c3d41846be8f.htm)）才需要被移除。  Specify a timeout to remove the current session list. Only TCP pipe ([PipeTcpNet](26b94bd5-2695-36e7-6378-c3d41846be8f.htm)) need to be removed. (继承自 [CommunicationServer](5a3cb700-c0f6-e35d-6807-dfe8bfecf747.htm)。) |
| 公共方法 | [RemoveSession(PipeSession, String)](92e9b2f6-96c3-b586-b554-f8ecd73df580.htm) | 移除一个管道会话  Remove a pipeline session (继承自 [CommunicationServer](5a3cb700-c0f6-e35d-6807-dfe8bfecf747.htm)。) |
| 公共方法 | [ServerClose](04a70d3a-d427-ed8d-53cf-3d3be56bdbd3.htm) | 关闭服务器的引擎  Shut down the server's engine (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |
| 公共方法 | [ServerStart](6aa1f94d-768c-4711-80e5-a529a5040f98.htm) | 指定端口号来启动服务器的引擎  Specify the port number to start the server's engine (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |
| 公共方法 | [ServerStart(Int32)](06e14b74-3deb-8e46-81c6-688f078375d9.htm) | 指定端口号来启动服务器的引擎  Specify the port number to start the server's engine (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |
| 公共方法 | [ServerStart(Int32, Boolean)](af75f437-3c8e-ac5c-96ab-96ea7bb4cdb4.htm) | 指定端口号来启动服务器的引擎  Specify the port number to start the server's engine (继承自 [CommunicationServer](5a3cb700-c0f6-e35d-6807-dfe8bfecf747.htm)。) |
| 公共方法 | [ServerStart(Int32, Int32)](1feb9368-f6bf-85e6-122a-d0c0d5a67bed.htm) | 指定一个TCP端口及UDP端口，同时启动两种模式的服务器  Specify a TCP port and a UDP port to start the server in both modes at the same time (继承自 [CommunicationServer](5a3cb700-c0f6-e35d-6807-dfe8bfecf747.htm)。) |
| 公共方法 | [SetNetMessage](836e8de1-0461-348c-b56b-2765e37c0ec5.htm) | 设置当前的服务器接收的消息信息 (继承自 [CommunicationServer](5a3cb700-c0f6-e35d-6807-dfe8bfecf747.htm)。) |
| 公共方法 | [SetSslPipeAction](5ecd4347-5dae-126d-f43c-76d4c7390705.htm) | 设置一个SSL的管道操作对象，在管道实例化之后，可以进行一些初始化的属性设置，例如自定义 SslProtocols 枚举 (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |
| 公共方法 | [SetTrustedIpAddress](6204b3e7-af37-7c20-98db-08114431b5ac.htm) | 设置并启动受信任的客户端登录并读写，如果为null，将关闭对客户端的ip验证  Set and start the trusted client login and read and write, if it is null, the client's IP verification will be turned off (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |
| 受保护的方法 | [SocketAcceptExtraCheck](3ed24041-5b60-57a7-7629-f7c49baf8bf9.htm) | 当客户端的socket登录的时候额外检查的操作，并返回操作的结果信息。  The operation is additionally checked when the client's socket logs in, and the result information of the operation is returned. (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |
| 公共方法 | [StartSerialSlave(ActionSerialPort)](dde390ea-45f6-bc3f-be69-3902f6a5641b.htm) | 启动串口的从机服务，使用自定义的初始化方法初始化串口的参数  Start the slave service of serial and initialize the parameters of the serial port using a custom initialization method (继承自 [CommunicationServer](5a3cb700-c0f6-e35d-6807-dfe8bfecf747.htm)。) |
| 公共方法 | [StartSerialSlave(String)](cfed6d77-50f0-eb64-37d0-7c466cc2d7bf.htm) | 启动串口的从机服务，使用默认的参数进行初始化串口，9600波特率，8位数据位，无奇偶校验，1位停止位  Start the slave service of serial, initialize the serial port with default parameters, 9600 baud rate, 8 data bits, no parity, 1 stop bit (继承自 [CommunicationServer](5a3cb700-c0f6-e35d-6807-dfe8bfecf747.htm)。) |
| 公共方法 | [StartSerialSlave(String, Int32)](155bfb40-a6bc-8ca0-3ebc-14307221ba59.htm) | 启动串口的从机服务，使用默认的参数进行初始化串口，8位数据位，无奇偶校验，1位停止位  Start the slave service of serial, initialize the serial port with default parameters, 8 data bits, no parity, 1 stop bit (继承自 [CommunicationServer](5a3cb700-c0f6-e35d-6807-dfe8bfecf747.htm)。) |
| 公共方法 | [StartSerialSlave(String, Int32, Int32, Parity, StopBits)](7045bee6-9c7f-65e0-191a-205c360dd42c.htm) | 启动串口的从机服务，使用指定的参数进行初始化串口，指定数据位，指定奇偶校验，指定停止位 (继承自 [CommunicationServer](5a3cb700-c0f6-e35d-6807-dfe8bfecf747.htm)。) |
| 受保护的方法 | [ThreadPoolLogin](7ebe2e21-f51e-3a08-a499-624462a0fede.htm) | 当客户端连接到服务器，并听过额外的检查后，进行回调的方法  Callback method when the client connects to the server and has heard additional checks (继承自 [CommunicationServer](5a3cb700-c0f6-e35d-6807-dfe8bfecf747.htm)。) |
| 公共方法 | [ToString](f4649b7b-3f2c-6b17-214a-ecde54810eab.htm) | (重写 [CommunicationServerToString](955ea032-e1fd-11ff-9ec8-72cb0572d36a.htm).) |
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

[SickIcrTcpServer 类](54f2dac1-62dc-1efb-aa3c-68b0d45a33cc.htm)

[HslCommunication.Profinet.Sick 命名空间](0dd7d6fd-55ea-53e6-68fc-8aa45401ad37.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ToString 方法 

[原文連結](http://api.hslcommunication.cn/html/f4649b7b-3f2c-6b17-214a-ecde54810eab.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Sick](../html/0dd7d6fd-55ea-53e6-68fc-8aa45401ad37.htm "HslCommunication.Profinet.Sick")

[SickIcrTcpServer 类](../html/54f2dac1-62dc-1efb-aa3c-68b0d45a33cc.htm "SickIcrTcpServer 类")

[SickIcrTcpServer 方法](../html/bdb23ae3-79a4-5717-4c42-e29fb17e11e5.htm "SickIcrTcpServer 方法")

[ToString 方法](../html/f4649b7b-3f2c-6b17-214a-ecde54810eab.htm "ToString 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| SickIcrTcpServerToString 方法 |

[缺少 "M:HslCommunication.Profinet.Sick.SickIcrTcpServer.ToString" 的 <summary> 文档]

**命名空间：**
 [HslCommunication.Profinet.Sick](0dd7d6fd-55ea-53e6-68fc-8aa45401ad37.htm)  
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

[缺少 "M:HslCommunication.Profinet.Sick.SickIcrTcpServer.ToString" 的 <returns> 文档]

![](../icons/SectionExpanded.png)参见

#### 引用

[SickIcrTcpServer 类](54f2dac1-62dc-1efb-aa3c-68b0d45a33cc.htm)

[HslCommunication.Profinet.Sick 命名空间](0dd7d6fd-55ea-53e6-68fc-8aa45401ad37.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## SickIcrTcpServer 事件

[原文連結](http://api.hslcommunication.cn/html/375747ce-c440-cf51-e045-2bb8b064ff64.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Sick](../html/0dd7d6fd-55ea-53e6-68fc-8aa45401ad37.htm "HslCommunication.Profinet.Sick")

[SickIcrTcpServer 类](../html/54f2dac1-62dc-1efb-aa3c-68b0d45a33cc.htm "SickIcrTcpServer 类")

[SickIcrTcpServer 事件](../html/375747ce-c440-cf51-e045-2bb8b064ff64.htm "SickIcrTcpServer 事件")

[OnReceivedBarCode 事件](../html/20f767de-bcf2-0b81-e456-149f87ecaa92.htm "OnReceivedBarCode 事件")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| SickIcrTcpServer 事件 |

[SickIcrTcpServer](54f2dac1-62dc-1efb-aa3c-68b0d45a33cc.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)事件

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共事件 | [OnClientOffline](d79ec6d7-e08e-9afe-3cc0-3a34f3f07017.htm) | 当客户端下线时候的触发的事件  Event triggered when the client goes offline (继承自 [CommunicationServer](5a3cb700-c0f6-e35d-6807-dfe8bfecf747.htm)。) |
| 公共事件 | [OnClientOnline](d60efb73-a438-eca9-3cf4-7bc26a24e165.htm) | 当客户端上线时候的触发的事件  Event triggered when the client goes online (继承自 [CommunicationServer](5a3cb700-c0f6-e35d-6807-dfe8bfecf747.htm)。) |
| 公共事件 | [OnPipeMessageReceived](539eb796-8339-4d86-7fa2-c9b6a50f533f.htm) | 当管道接收到消息时触发 (继承自 [CommunicationServer](5a3cb700-c0f6-e35d-6807-dfe8bfecf747.htm)。) |
| 公共事件 | [OnReceivedBarCode](20f767de-bcf2-0b81-e456-149f87ecaa92.htm) | 当接收到条码数据的时候触发  Triggered when barcode data is received |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[SickIcrTcpServer 类](54f2dac1-62dc-1efb-aa3c-68b0d45a33cc.htm)

[HslCommunication.Profinet.Sick 命名空间](0dd7d6fd-55ea-53e6-68fc-8aa45401ad37.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## OnReceivedBarCode 事件

[原文連結](http://api.hslcommunication.cn/html/20f767de-bcf2-0b81-e456-149f87ecaa92.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Sick](../html/0dd7d6fd-55ea-53e6-68fc-8aa45401ad37.htm "HslCommunication.Profinet.Sick")

[SickIcrTcpServer 类](../html/54f2dac1-62dc-1efb-aa3c-68b0d45a33cc.htm "SickIcrTcpServer 类")

[SickIcrTcpServer 事件](../html/375747ce-c440-cf51-e045-2bb8b064ff64.htm "SickIcrTcpServer 事件")

[OnReceivedBarCode 事件](../html/20f767de-bcf2-0b81-e456-149f87ecaa92.htm "OnReceivedBarCode 事件")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| SickIcrTcpServerOnReceivedBarCode 事件 |

当接收到条码数据的时候触发  
Triggered when barcode data is received

**命名空间：**
 [HslCommunication.Profinet.Sick](0dd7d6fd-55ea-53e6-68fc-8aa45401ad37.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public event SickIcrTcpServerReceivedBarCodeDelegate OnReceivedBarCode
```

```
Public Event OnReceivedBarCode As SickIcrTcpServerReceivedBarCodeDelegate
```

```
public:
 event SickIcrTcpServerReceivedBarCodeDelegate^ OnReceivedBarCode {
	void add (SickIcrTcpServerReceivedBarCodeDelegate^ value);
	void remove (SickIcrTcpServerReceivedBarCodeDelegate^ value);
}
```

```
member OnReceivedBarCode : IEvent<SickIcrTcpServerReceivedBarCodeDelegate,
    EventArgs>
```

#### 值

类型：[HslCommunication.Profinet.SickSickIcrTcpServerReceivedBarCodeDelegate](121fcfd3-7f91-ed26-cccf-98db3fa02780.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[SickIcrTcpServer 类](54f2dac1-62dc-1efb-aa3c-68b0d45a33cc.htm)

[HslCommunication.Profinet.Sick 命名空间](0dd7d6fd-55ea-53e6-68fc-8aa45401ad37.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## SickIcrTcpServer 字段

[原文連結](http://api.hslcommunication.cn/html/fff1c148-1845-faa4-5860-217826955f40.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Sick](../html/0dd7d6fd-55ea-53e6-68fc-8aa45401ad37.htm "HslCommunication.Profinet.Sick")

[SickIcrTcpServer 类](../html/54f2dac1-62dc-1efb-aa3c-68b0d45a33cc.htm "SickIcrTcpServer 类")

[SickIcrTcpServer 构造函数](../html/60dfe1fa-3a09-7b39-b644-e5f58760b8d3.htm "SickIcrTcpServer 构造函数 ")

[SickIcrTcpServer 属性](../html/fe9ef976-5561-aeb7-3156-02ae5bd248b1.htm "SickIcrTcpServer 属性")

[SickIcrTcpServer 方法](../html/bdb23ae3-79a4-5717-4c42-e29fb17e11e5.htm "SickIcrTcpServer 方法")

[SickIcrTcpServer 事件](../html/375747ce-c440-cf51-e045-2bb8b064ff64.htm "SickIcrTcpServer 事件")

[SickIcrTcpServer 字段](../html/fff1c148-1845-faa4-5860-217826955f40.htm "SickIcrTcpServer 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| SickIcrTcpServer 字段 |

[SickIcrTcpServer](54f2dac1-62dc-1efb-aa3c-68b0d45a33cc.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)字段

|  | 名称 | 说明 |
| --- | --- | --- |
| 受保护的字段 | [socketServer](f8c04282-bfb8-502e-730b-8d4127164ecb.htm) | 核心的socket服务器 (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[SickIcrTcpServer 类](54f2dac1-62dc-1efb-aa3c-68b0d45a33cc.htm)

[HslCommunication.Profinet.Sick 命名空间](0dd7d6fd-55ea-53e6-68fc-8aa45401ad37.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## SickIcrTcpServer.ReceivedBarCodeDelegate 委托

[原文連結](http://api.hslcommunication.cn/html/121fcfd3-7f91-ed26-cccf-98db3fa02780.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Sick](../html/0dd7d6fd-55ea-53e6-68fc-8aa45401ad37.htm "HslCommunication.Profinet.Sick")

[SickIcrTcpServer 类](../html/54f2dac1-62dc-1efb-aa3c-68b0d45a33cc.htm "SickIcrTcpServer 类")

[SickIcrTcpServer.ReceivedBarCodeDelegate 委托](../html/121fcfd3-7f91-ed26-cccf-98db3fa02780.htm "SickIcrTcpServer.ReceivedBarCodeDelegate 委托")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| SickIcrTcpServerReceivedBarCodeDelegate 委托 |

接收条码数据的委托信息  
Entrusted information to receive barcode data

**命名空间：**
 [HslCommunication.Profinet.Sick](0dd7d6fd-55ea-53e6-68fc-8aa45401ad37.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public delegate void ReceivedBarCodeDelegate(
	string ipAddress,
	string barCode
)
```

```
Public Delegate Sub ReceivedBarCodeDelegate ( 
	ipAddress As String,
	barCode As String
)
```

```
public delegate void ReceivedBarCodeDelegate(
	String^ ipAddress, 
	String^ barCode
)
```

```
type ReceivedBarCodeDelegate = 
    delegate of 
        ipAddress : string * 
        barCode : string -> unit
```

#### 参数

ipAddress
:   类型：SystemString  
    Ip地址信息

barCode
:   类型：SystemString  
    条码信息

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.Profinet.Sick 命名空间](0dd7d6fd-55ea-53e6-68fc-8aa45401ad37.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)