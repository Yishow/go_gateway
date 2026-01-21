# HslCommunication - HslCommunication.Enthernet.Ftp

> 分類頁數: 30



---
## HslCommunication.Enthernet.Ftp

[原文連結](http://api.hslcommunication.cn/html/05837d2e-e0a8-4564-6d5e-a834d56ce5ce.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Enthernet.Ftp](../html/05837d2e-e0a8-4564-6d5e-a834d56ce5ce.htm "HslCommunication.Enthernet.Ftp")

[FtpClient 类](../html/98289259-aeed-7095-b875-9d3ef78237fe.htm "FtpClient 类")

[FtpFileItem 类](../html/97ec5ae6-c5a1-0122-982b-cb02f1031c96.htm "FtpFileItem 类")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslCommunication.Enthernet.Ftp 命名空间 |

[缺少 "N:HslCommunication.Enthernet.Ftp" 的 <summary> 文档]

![](../icons/SectionExpanded.png)类

|  | 类 | 说明 |
| --- | --- | --- |
| 公共类 | [FtpClient](98289259-aeed-7095-b875-9d3ef78237fe.htm) | FTP协议的客户端对象 |
| 公共类 | [FtpFileItem](97ec5ae6-c5a1-0122-982b-cb02f1031c96.htm) | FTP文件的信息对象 |

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## FtpClient 类

[原文連結](http://api.hslcommunication.cn/html/98289259-aeed-7095-b875-9d3ef78237fe.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Enthernet.Ftp](../html/05837d2e-e0a8-4564-6d5e-a834d56ce5ce.htm "HslCommunication.Enthernet.Ftp")

[FtpClient 类](../html/98289259-aeed-7095-b875-9d3ef78237fe.htm "FtpClient 类")

[FtpClient 构造函数](../html/034f536b-56cd-7b2b-7f55-c493aa15c924.htm "FtpClient 构造函数 ")

[FtpClient 属性](../html/69a325b9-c5bb-9bf7-f4b4-847f840c54d4.htm "FtpClient 属性")

[FtpClient 方法](../html/c26ace03-ff28-7070-8a9e-62b8abf31430.htm "FtpClient 方法")

[FtpClient 字段](../html/2f2894ac-2a80-2e5c-3af6-6258dfebafdf.htm "FtpClient 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FtpClient 类 |

FTP协议的客户端对象

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  [HslCommunication.Core.NetBinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)  
    [HslCommunication.Core.NetTcpNetCommunication](01e0c54c-8c1b-d3eb-4947-0d134d446266.htm)  
      HslCommunication.Enthernet.FtpFtpClient

**命名空间：**
 [HslCommunication.Enthernet.Ftp](05837d2e-e0a8-4564-6d5e-a834d56ce5ce.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class FtpClient : TcpNetCommunication
```

```
Public Class FtpClient
	Inherits TcpNetCommunication
```

```
public ref class FtpClient : public TcpNetCommunication
```

```
type FtpClient =  
    class
        inherit TcpNetCommunication
    end
```

FtpClient 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [FtpClient](034f536b-56cd-7b2b-7f55-c493aa15c924.htm) | 实例化一个默认的对象 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [CommunicationPipe](ad8e903b-d289-21b0-077e-15ede8005221.htm) | 获取或设置当前的管道信息，管道类型为[CommunicationPipe](ad8e903b-d289-21b0-077e-15ede8005221.htm)的继承类，内置了[PipeTcpNet](26b94bd5-2695-36e7-6378-c3d41846be8f.htm)管道，[PipeUdpNet](5a0b810e-187e-826b-0a13-34fd5576aa0e.htm)管道，[PipeSerialPort](9d924607-3dfc-3916-b532-c5dec48fe347.htm)管道等  Get or set the current pipeline information, the pipeline type is [CommunicationPipe](ad8e903b-d289-21b0-077e-15ede8005221.htm) inheritance class, [PipeTcpNet](26b94bd5-2695-36e7-6378-c3d41846be8f.htm) pipeline, [PipeUdpNet](5a0b810e-187e-826b-0a13-34fd5576aa0e.htm) pipeline, [PipeSerialPort](9d924607-3dfc-3916-b532-c5dec48fe347.htm) pipeline, etc (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [ConnectionId](ce1efb4b-fbc9-4683-163e-0a8e71e99d6b.htm) | 当前连接的唯一ID号，默认为长度20的guid码加随机数组成，方便列表管理，也可以自己指定  The unique ID number of the current connection. The default is a 20-digit guid code plus a random number. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性代码示例 | [ConnectTimeOut](b40f14d5-6462-45c2-477a-5c1494446993.htm) | 获取或设置接收服务器反馈的时间，如果为负数，则不接收反馈   Gets or sets the time to receive server feedback, and if it is a negative number, does not receive feedback (继承自 [TcpNetCommunication](01e0c54c-8c1b-d3eb-4947-0d134d446266.htm)。) |
| 公共属性代码示例 | [IpAddress](f082e96b-09f5-d3c7-03cd-3926dc8f8b04.htm) | 获取或是设置远程服务器的IP地址，如果是本机测试，那么需要设置为127.0.0.1   Get or set the IP address of the remote server. If it is a local test, then it needs to be set to 127.0.0.1 (继承自 [TcpNetCommunication](01e0c54c-8c1b-d3eb-4947-0d134d446266.htm)。) |
| 公共属性 | [LocalBinding](3736875e-eb89-6d38-4617-c759e0a4d059.htm) | 获取或设置绑定的本地的IP地址和端口号信息，如果端口设置为0，代表任何可用的端口  Get or set the bound local IP address and port number information, if the port is set to 0, it means any available port (继承自 [TcpNetCommunication](01e0c54c-8c1b-d3eb-4947-0d134d446266.htm)。) |
| 公共属性代码示例 | [LogNet](8c56f246-88ac-afb8-8c06-972fafe4f78a.htm) | 组件的日志工具，支持日志记录，只要实例化后，当前网络的基本信息，就以[DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm)等级进行输出  The component's logging tool supports logging. As long as the instantiation of the basic network information, the output will be output at [DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm) (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [Password](ec1a8da1-fb8a-54ef-aac4-53a04498488d.htm) | 获取或设置登录的密码信息，默认值为 hslcommunication@example.com |
| 公共属性代码示例 | [Port](9bb71d3f-301e-a868-8021-fb5e027da9a0.htm) | 获取或设置服务器的端口号，具体的值需要取决于对方的配置  Gets or sets the port number of the server. The specific value depends on the configuration of the other party. (继承自 [TcpNetCommunication](01e0c54c-8c1b-d3eb-4947-0d134d446266.htm)。) |
| 公共属性代码示例 | [ReceiveTimeOut](250560a7-332b-ce1b-5eb0-e3abe4ee6012.htm) | 获取或设置接收服务器反馈的时间，如果为负数，则不接收反馈   Gets or sets the time to receive server feedback, and if it is a negative number, does not receive feedback (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [SendBeforeHex](72aa1de4-116f-abfe-4adb-814a0be0e560.htm) | 获取或设置在发送通信报文前追加发送的字节信息，HEX格式，通常用于lora组网时，需要携带 00 00 00 02 四个字节的站地址功能。  Obtain or set the byte information sent before sending communication packets, HEX format, usually used for LORA networking, you need to carry 00 00 00 02 four-byte station address function. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [SleepTime](86a5fbb9-3456-6beb-7113-5bbad5eabbc5.htm) | 获取或设置在正式接收对方返回数据前的时候，需要休息的时间，当设置为0的时候，不需要休息。  Get or set the time required to rest before officially receiving the data from the other party. When it is set to 0, no rest is required. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [SocketKeepAliveTime](53f6a35c-7bf4-3c34-4964-35363a5ac351.htm) | 获取或设置客户端的Socket的心跳时间信息，这个是Socket底层自动实现的心跳包，不基于协议层实现。默认小于0，不开启心跳检测，如果需要开启，设置 60\_000 比较合适，单位毫秒  Get or set the heartbeat time information of the Socket of the client. This is the heartbeat packet automatically implemented by the bottom layer of the Socket, not based on the protocol layer. The default value is less than 0, and heartbeat detection is not enabled. If you need to enable it, it is more appropriate to set 60\_000, in milliseconds. (继承自 [TcpNetCommunication](01e0c54c-8c1b-d3eb-4947-0d134d446266.htm)。) |
| 公共属性 | [SystemType](147c3982-4761-334d-9d78-41d745a8b753.htm) | 系统的类型信息 |
| 公共属性 | [Username](bf1063c3-0f40-b74b-4323-29156b0d8f4e.htm) | 获取或设置登录的用户名信息，默认值为anonymous |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法代码示例 | [ConnectClose](5ff1efa3-0f43-7515-add2-48681749ca7b.htm) | 手动断开与远程服务器的连接，如果当前是长连接模式，那么就会切换到短连接模式  Manually disconnect from the remote server, if it is currently in long connection mode, it will switch to short connection mode (继承自 [TcpNetCommunication](01e0c54c-8c1b-d3eb-4947-0d134d446266.htm)。) |
| 公共方法代码示例 | [ConnectCloseAsync](3b03af5f-46f8-e1a9-1964-228bac33c510.htm) | 手动断开与远程服务器的连接，如果当前是长连接模式，那么就会切换到短连接模式  Manually disconnect from the remote server, if it is currently in long connection mode, it will switch to short connection mode (继承自 [TcpNetCommunication](01e0c54c-8c1b-d3eb-4947-0d134d446266.htm)。) |
| 公共方法代码示例 | [ConnectServer](1dc20e7a-8929-bad2-d034-b2454d55a906.htm) | 尝试连接远程的服务器，如果连接成功，就切换短连接模式到长连接模式，后面的每次请求都共享一个通道，使得通讯速度更快速  Try to connect to a remote server. If the connection is successful, switch the short connection mode to the long connection mode. Each subsequent request will share a channel, making the communication speed faster. (继承自 [TcpNetCommunication](01e0c54c-8c1b-d3eb-4947-0d134d446266.htm)。) |
| 公共方法代码示例 | [ConnectServerAsync](9e07b3d8-a288-803d-c27c-1d75f4f92aa7.htm) | 尝试连接远程的服务器，如果连接成功，就切换短连接模式到长连接模式，后面的每次请求都共享一个通道，使得通讯速度更快速  Try to connect to a remote server. If the connection is successful, switch the short connection mode to the long connection mode. Each subsequent request will share a channel, making the communication speed faster. (继承自 [TcpNetCommunication](01e0c54c-8c1b-d3eb-4947-0d134d446266.htm)。) |
| 公共方法 | [CreateDirectory](82ac3dc9-8139-1130-727c-3df3f3d5ce75.htm) | 在FTP服务器上创建指定的目录  Create the specified directory on the FTP server. |
| 受保护的方法 | [DecideWhetherQAMessage](acc40cda-e310-8537-6685-59b67eeb16ff.htm) | 决定当前的消息是否是用于问答机制返回的消息，默认直接返回 true, 实际的情况需要根据协议进行重写方法  To determine whether the current message is the message returned by the question answering mechanism, the default is true. In actual cases, the rewriting method needs to be performed according to the protocol (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [DeleteDirectory](ad8fd03c-4875-6c6b-83a3-e7590c7891a3.htm) | 删除FTP服务器上的指定目录，并且可以指定是否删除所有子文件  Delete the specified directory on the FTP server, with the option to specify whether to delete all subfiles. |
| 公共方法 | [DeleteFile](e86754cd-59fd-e54b-be3d-18d3de267d61.htm) | 删除FTP服务器上的指定文件  Delete the specified file(s) on the FTP server. |
| 公共方法 | [DownloadFile](41dcfcde-d66a-0549-86f2-0063edee69f3.htm) | 从FTP服务器下载文件到本地文件，支持进度汇报和取消操作  Download files from an FTP server to local files, with support for progress reporting and cancellation operations. |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | [ExtraAfterReadFromCoreServer](06a49a15-a3ce-5e39-28c7-2b78851eca2e.htm) | 和服务器交互完成的时候调用的方法，可以根据读写结果进行一些额外的操作，具体的操作需要根据实际的需求来重写实现  The method called when the interaction with the server is completed can perform some additional operations based on the read and write results. The specific operations need to be rewritten according to actual needs. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法代码示例 | [ExtraOnDisconnect](8113d5c7-4f60-fba4-0e43-c0dd4a778d59.htm) | 根据实际的协议选择是否重写本方法，有些协议在断开连接之前，需要发送一些报文来关闭当前的网络通道  Select whether to rewrite this method according to the actual protocol. Some protocols need to send some packets to close the current network channel before disconnecting. (重写 [BinaryCommunicationExtraOnDisconnect](de9dbba2-3b2c-1c5e-53c8-4e6f8a44e2c4.htm).) |
| 受保护的方法代码示例 | [ExtraOnDisconnectAsync](210f593e-16d8-77c4-1178-58fbcfe9b27b.htm) | 根据实际的协议选择是否重写本方法，有些协议在断开连接之前，需要发送一些报文来关闭当前的网络通道  Select whether to rewrite this method according to the actual protocol. Some protocols need to send some packets to close the current network channel before disconnecting. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 受保护的方法 | [GetLogTextFromBinary](051c0c92-9f31-e96b-26f7-3e1b1a63f272.htm) | 获取当前的报文进行日志记录的时候，是否使用二进制的格式记录，默认返回 [LogMsgFormatBinary](2ed2ff87-caec-0e44-f08b-e2b015225a58.htm)，重写可以根据session对象分别返回不同记录模式  Whether to log the current packet in binary format, the default return is [LogMsgFormatBinary](2ed2ff87-caec-0e44-f08b-e2b015225a58.htm). If you want to override it, different recording modes can be returned according to session (重写 [BinaryCommunicationGetLogTextFromBinary(PipeSession, Byte)](67f3b2f7-957c-fa0b-8320-6799237157db.htm).) |
| 受保护的方法 | [GetNewNetMessage](9cb70414-4298-bcf3-fb20-79e4cb0393ad.htm) | 获取一个新的消息对象的方法，需要在继承类里面进行重写  The method to get a new message object needs to be overridden in the inheritance class (重写 [BinaryCommunicationGetNewNetMessage](16636a1b-a6f0-9c0f-9d9c-c94e26c67459.htm).) |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法代码示例 | [InitializationOnConnect](56cbb726-5571-5c28-4c66-e1820eef5df8.htm) | 根据实际的协议选择是否重写本方法，有些协议在创建连接之后，需要进行一些初始化的信号握手，才能最终建立网络通道。  Whether to rewrite this method is based on the actual protocol. Some protocols require some initial signal handshake to establish a network channel after the connection is created. (重写 [BinaryCommunicationInitializationOnConnect](4ab6989f-a5b5-05fe-8b44-bfa45824a2a7.htm).) |
| 受保护的方法代码示例 | [InitializationOnConnectAsync](7da2af04-d4fc-5045-19a4-4b0d1722eef3.htm) | 根据实际的协议选择是否重写本方法，有些协议在创建连接之后，需要进行一些初始化的信号握手，才能最终建立网络通道。  Whether to rewrite this method is based on the actual protocol. Some protocols require some initial signal handshake to establish a network channel after the connection is created. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [IpAddressPing](1011593d-fb48-acd1-44c2-27fd565a1180.htm) | 对当前设备的IP地址进行PING的操作，返回PING的结果，正常来说，返回Success  PING the IP address of the current device and return the PING result. Normally, it returns Success (继承自 [TcpNetCommunication](01e0c54c-8c1b-d3eb-4947-0d134d446266.htm)。) |
| 公共方法 | [ListFiles](c0783d51-47c7-c3d2-bf08-fa9afa5d50c9.htm) | 浏览指定目录下的文件及文件夹信息列表，默认根目录为 /   Browse the list of file and folder information in the specified directory. The default root directory is / |
| 受保护的方法 | [LogRevcMessage(Byte)](dab1657f-1e75-7fe1-7003-ccedfc69855d.htm) | 使用日志记录一个接收的报文信息  Logs are used to record information about a received packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [LogRevcMessage(Byte, PipeSession)](0f51ae45-575e-4228-f978-01b37d6d5df1.htm) | 使用日志记录一个接收的报文信息  Logs are used to record information about a received packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [LogSendMessage(Byte)](4f90eecc-68a2-1ef2-44ef-369f19bed3ec.htm) | 使用日志记录一个发送的报文信息  Logs are used to record information about a send packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [LogSendMessage(Byte, PipeSession)](209caadf-4f0f-40aa-04c2-ebc61f8dab13.htm) | 使用日志记录一个发送的报文信息  Logs are used to record information about a send packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法 | [PackCommandWithHeader](d0b2cd76-c306-11fc-38c7-af8ad306cce7.htm) | 对当前的命令进行打包处理，通常是携带命令头内容，标记当前的命令的长度信息，需要进行重写，否则默认不打包  The current command is packaged, usually carrying the content of the command header, marking the length of the current command, and it needs to be rewritten, otherwise it is not packaged by default (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServer(Byte)](8146dc0c-1edc-3ca7-428a-a3bdbd35e6c5.htm) | 将二进制的数据发送到管道中去，然后从管道里接收二进制的数据回来，并返回是否成功的结果对象。  Send binary data to the pipeline, and then receive binary data back from the pipeline, and return whether the success of the result object (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServer(IEnumerableByte)](044cfdd9-ec76-b89f-2f2b-5f439025006c.htm) | 将多个数据报文按顺序发到设备，并从设备接收返回的数据内容，然后拼接成一个Byte[]信息，需要重写[UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm)方法才能返回正确的结果。  Send multiple data packets to the device in sequence, and receive the returned data content from the device, and then splicing them into a Byte[] message, you need to rewrite [UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm) method to return the correct result. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServer(Byte, Boolean, Boolean)](b32ce8e8-a959-fcc2-6736-14b39fdb9a40.htm) | 将二进制的数据发送到管道中去，然后从管道里接收二进制的数据回来，并返回是否成功的结果对象。  Send binary data to the pipeline, and then receive binary data back from the pipeline, and return whether the success of the result object (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServer(CommunicationPipe, Byte, Boolean, Boolean)](ef36d56b-df09-f727-e8d3-316d37a13ee9.htm) | 使用指定的管道来进行数据通信，发送原始数据到管道，然后从管道接收相关的数据返回，本方法无锁 (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServerAsync(Byte)](a4de5a24-c67a-7b02-db05-33b85eaf6e1e.htm) | 将二进制的数据发送到管道中去，然后从管道里接收二进制的数据回来，并返回是否成功的结果对象。  Send binary data to the pipeline, and then receive binary data back from the pipeline, and return whether the success of the result object (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServerAsync(IEnumerableByte)](bcdb59b7-a70a-d5c0-1fc6-9ececcd39ea2.htm) | 将多个数据报文按顺序发到设备，并从设备接收返回的数据内容，然后拼接成一个Byte[]信息，需要重写[UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm)方法才能返回正确的结果。  Send multiple data packets to the device in sequence, and receive the returned data content from the device, and then splicing them into a Byte[] message, you need to rewrite [UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm) method to return the correct result. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServerAsync(Byte, Boolean, Boolean)](1c79c204-5ab5-9342-e65d-7455c20ea577.htm) | 将二进制的数据发送到管道中去，然后从管道里接收二进制的数据回来，并返回是否成功的结果对象。  Send binary data to the pipeline, and then receive binary data back from the pipeline, and return whether the success of the result object (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServerAsync(CommunicationPipe, Byte, Boolean, Boolean)](a254a19d-3c29-147b-335b-8adccdf4ee2f.htm) | 使用指定的管道来进行数据通信，发送原始数据到管道，然后从管道接收相关的数据返回，本方法无锁 (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [RenameFile](d6aeb935-d84d-2f14-e447-201605bdaa27.htm) | 重命名文件的名称  Rename the name of the file. |
| 公共方法 | [RenameFolder](9348c583-a8a6-89f0-57e1-31ccae07c876.htm) | 重命名目录，不支持跨目录重命名  Rename a directory; cross-directory renaming is not supported. |
| 公共方法 | [SetDtuPipe](f765d1d7-234a-b229-9457-97e25e024dc9.htm) | 将当前的通信对象设置DTU模式，允许传入现成的管道，并返回初始化结果，如果该设备重写了握手报文，就是返回握手结果  Set the current communication object to DTU mode, allow the existing pipe to be passed in, and return the initialization result, if the device rewrites the handshake packet, the handshake result is returned (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [SetDtuPipeAsync](9d1fd12d-b804-3828-a627-a2c0eef1d23f.htm) | 将当前的通信对象设置DTU模式，允许传入现成的管道，并返回初始化结果，如果该设备重写了握手报文，就是返回握手结果  Set the current communication object to DTU mode, allow the existing pipe to be passed in, and return the initialization result, if the device rewrites the handshake packet, the handshake result is returned (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ToString](83ab8b3b-8124-ea01-71fa-3bdbea2ba010.htm) | (重写 [TcpNetCommunicationToString](c1a8a99d-7115-55ca-57df-ece8e93e460c.htm).) |
| 公共方法 | [UnpackResponseContent](208017a7-8da2-33f8-3273-d7e7e850100f.htm) | 根据对方返回的报文命令，对命令进行基本的拆包，例如各种Modbus协议拆包为统一的核心报文，还支持对报文的验证  According to the message command returned by the other party, the command is basically unpacked, for example, various Modbus protocols are unpacked into a unified core message, and the verification of the message is also supported (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [UploadFile](9b1e30fb-8019-f6fd-8891-6d05d1425af8.htm) | 上传本地文件到FTP服务器，支持进度汇报和取消操作  Upload local files to an FTP server, with support for progress reporting and cancellation operations. |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)字段

|  | 名称 | 说明 |
| --- | --- | --- |
| 受保护的字段 | [LogMsgFormatBinary](2ed2ff87-caec-0e44-f08b-e2b015225a58.htm) | 设置日志记录报文是否二进制，如果为False，那就使用ASCII码  Set whether the log message is binary, if it is False, then use ASCII code (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.Enthernet.Ftp 命名空间](05837d2e-e0a8-4564-6d5e-a834d56ce5ce.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## FtpClient 构造函数 

[原文連結](http://api.hslcommunication.cn/html/034f536b-56cd-7b2b-7f55-c493aa15c924.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Enthernet.Ftp](../html/05837d2e-e0a8-4564-6d5e-a834d56ce5ce.htm "HslCommunication.Enthernet.Ftp")

[FtpClient 类](../html/98289259-aeed-7095-b875-9d3ef78237fe.htm "FtpClient 类")

[FtpClient 构造函数](../html/034f536b-56cd-7b2b-7f55-c493aa15c924.htm "FtpClient 构造函数 ")

[FtpClient 属性](../html/69a325b9-c5bb-9bf7-f4b4-847f840c54d4.htm "FtpClient 属性")

[FtpClient 方法](../html/c26ace03-ff28-7070-8a9e-62b8abf31430.htm "FtpClient 方法")

[FtpClient 字段](../html/2f2894ac-2a80-2e5c-3af6-6258dfebafdf.htm "FtpClient 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FtpClient 构造函数 |

实例化一个默认的对象

**命名空间：**
 [HslCommunication.Enthernet.Ftp](05837d2e-e0a8-4564-6d5e-a834d56ce5ce.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public FtpClient()
```

```
Public Sub New
```

```
public:
FtpClient()
```

```
new : unit -> FtpClient
```

![](../icons/SectionExpanded.png)参见

#### 引用

[FtpClient 类](98289259-aeed-7095-b875-9d3ef78237fe.htm)

[HslCommunication.Enthernet.Ftp 命名空间](05837d2e-e0a8-4564-6d5e-a834d56ce5ce.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## FtpClient 属性

[原文連結](http://api.hslcommunication.cn/html/69a325b9-c5bb-9bf7-f4b4-847f840c54d4.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Enthernet.Ftp](../html/05837d2e-e0a8-4564-6d5e-a834d56ce5ce.htm "HslCommunication.Enthernet.Ftp")

[FtpClient 类](../html/98289259-aeed-7095-b875-9d3ef78237fe.htm "FtpClient 类")

[FtpClient 属性](../html/69a325b9-c5bb-9bf7-f4b4-847f840c54d4.htm "FtpClient 属性")

[Password 属性](../html/ec1a8da1-fb8a-54ef-aac4-53a04498488d.htm "Password 属性 ")

[SystemType 属性](../html/147c3982-4761-334d-9d78-41d745a8b753.htm "SystemType 属性 ")

[Username 属性](../html/bf1063c3-0f40-b74b-4323-29156b0d8f4e.htm "Username 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FtpClient 属性 |

[FtpClient](98289259-aeed-7095-b875-9d3ef78237fe.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [CommunicationPipe](ad8e903b-d289-21b0-077e-15ede8005221.htm) | 获取或设置当前的管道信息，管道类型为[CommunicationPipe](ad8e903b-d289-21b0-077e-15ede8005221.htm)的继承类，内置了[PipeTcpNet](26b94bd5-2695-36e7-6378-c3d41846be8f.htm)管道，[PipeUdpNet](5a0b810e-187e-826b-0a13-34fd5576aa0e.htm)管道，[PipeSerialPort](9d924607-3dfc-3916-b532-c5dec48fe347.htm)管道等  Get or set the current pipeline information, the pipeline type is [CommunicationPipe](ad8e903b-d289-21b0-077e-15ede8005221.htm) inheritance class, [PipeTcpNet](26b94bd5-2695-36e7-6378-c3d41846be8f.htm) pipeline, [PipeUdpNet](5a0b810e-187e-826b-0a13-34fd5576aa0e.htm) pipeline, [PipeSerialPort](9d924607-3dfc-3916-b532-c5dec48fe347.htm) pipeline, etc (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [ConnectionId](ce1efb4b-fbc9-4683-163e-0a8e71e99d6b.htm) | 当前连接的唯一ID号，默认为长度20的guid码加随机数组成，方便列表管理，也可以自己指定  The unique ID number of the current connection. The default is a 20-digit guid code plus a random number. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性代码示例 | [ConnectTimeOut](b40f14d5-6462-45c2-477a-5c1494446993.htm) | 获取或设置接收服务器反馈的时间，如果为负数，则不接收反馈   Gets or sets the time to receive server feedback, and if it is a negative number, does not receive feedback (继承自 [TcpNetCommunication](01e0c54c-8c1b-d3eb-4947-0d134d446266.htm)。) |
| 公共属性代码示例 | [IpAddress](f082e96b-09f5-d3c7-03cd-3926dc8f8b04.htm) | 获取或是设置远程服务器的IP地址，如果是本机测试，那么需要设置为127.0.0.1   Get or set the IP address of the remote server. If it is a local test, then it needs to be set to 127.0.0.1 (继承自 [TcpNetCommunication](01e0c54c-8c1b-d3eb-4947-0d134d446266.htm)。) |
| 公共属性 | [LocalBinding](3736875e-eb89-6d38-4617-c759e0a4d059.htm) | 获取或设置绑定的本地的IP地址和端口号信息，如果端口设置为0，代表任何可用的端口  Get or set the bound local IP address and port number information, if the port is set to 0, it means any available port (继承自 [TcpNetCommunication](01e0c54c-8c1b-d3eb-4947-0d134d446266.htm)。) |
| 公共属性代码示例 | [LogNet](8c56f246-88ac-afb8-8c06-972fafe4f78a.htm) | 组件的日志工具，支持日志记录，只要实例化后，当前网络的基本信息，就以[DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm)等级进行输出  The component's logging tool supports logging. As long as the instantiation of the basic network information, the output will be output at [DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm) (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [Password](ec1a8da1-fb8a-54ef-aac4-53a04498488d.htm) | 获取或设置登录的密码信息，默认值为 hslcommunication@example.com |
| 公共属性代码示例 | [Port](9bb71d3f-301e-a868-8021-fb5e027da9a0.htm) | 获取或设置服务器的端口号，具体的值需要取决于对方的配置  Gets or sets the port number of the server. The specific value depends on the configuration of the other party. (继承自 [TcpNetCommunication](01e0c54c-8c1b-d3eb-4947-0d134d446266.htm)。) |
| 公共属性代码示例 | [ReceiveTimeOut](250560a7-332b-ce1b-5eb0-e3abe4ee6012.htm) | 获取或设置接收服务器反馈的时间，如果为负数，则不接收反馈   Gets or sets the time to receive server feedback, and if it is a negative number, does not receive feedback (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [SendBeforeHex](72aa1de4-116f-abfe-4adb-814a0be0e560.htm) | 获取或设置在发送通信报文前追加发送的字节信息，HEX格式，通常用于lora组网时，需要携带 00 00 00 02 四个字节的站地址功能。  Obtain or set the byte information sent before sending communication packets, HEX format, usually used for LORA networking, you need to carry 00 00 00 02 four-byte station address function. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [SleepTime](86a5fbb9-3456-6beb-7113-5bbad5eabbc5.htm) | 获取或设置在正式接收对方返回数据前的时候，需要休息的时间，当设置为0的时候，不需要休息。  Get or set the time required to rest before officially receiving the data from the other party. When it is set to 0, no rest is required. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [SocketKeepAliveTime](53f6a35c-7bf4-3c34-4964-35363a5ac351.htm) | 获取或设置客户端的Socket的心跳时间信息，这个是Socket底层自动实现的心跳包，不基于协议层实现。默认小于0，不开启心跳检测，如果需要开启，设置 60\_000 比较合适，单位毫秒  Get or set the heartbeat time information of the Socket of the client. This is the heartbeat packet automatically implemented by the bottom layer of the Socket, not based on the protocol layer. The default value is less than 0, and heartbeat detection is not enabled. If you need to enable it, it is more appropriate to set 60\_000, in milliseconds. (继承自 [TcpNetCommunication](01e0c54c-8c1b-d3eb-4947-0d134d446266.htm)。) |
| 公共属性 | [SystemType](147c3982-4761-334d-9d78-41d745a8b753.htm) | 系统的类型信息 |
| 公共属性 | [Username](bf1063c3-0f40-b74b-4323-29156b0d8f4e.htm) | 获取或设置登录的用户名信息，默认值为anonymous |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[FtpClient 类](98289259-aeed-7095-b875-9d3ef78237fe.htm)

[HslCommunication.Enthernet.Ftp 命名空间](05837d2e-e0a8-4564-6d5e-a834d56ce5ce.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Password 属性 

[原文連結](http://api.hslcommunication.cn/html/ec1a8da1-fb8a-54ef-aac4-53a04498488d.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Enthernet.Ftp](../html/05837d2e-e0a8-4564-6d5e-a834d56ce5ce.htm "HslCommunication.Enthernet.Ftp")

[FtpClient 类](../html/98289259-aeed-7095-b875-9d3ef78237fe.htm "FtpClient 类")

[FtpClient 属性](../html/69a325b9-c5bb-9bf7-f4b4-847f840c54d4.htm "FtpClient 属性")

[Password 属性](../html/ec1a8da1-fb8a-54ef-aac4-53a04498488d.htm "Password 属性 ")

[SystemType 属性](../html/147c3982-4761-334d-9d78-41d745a8b753.htm "SystemType 属性 ")

[Username 属性](../html/bf1063c3-0f40-b74b-4323-29156b0d8f4e.htm "Username 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FtpClientPassword 属性 |

获取或设置登录的密码信息，默认值为 hslcommunication@example.com

**命名空间：**
 [HslCommunication.Enthernet.Ftp](05837d2e-e0a8-4564-6d5e-a834d56ce5ce.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public string Password { get; set; }
```

```
Public Property Password As String
	Get
	Set
```

```
public:
property String^ Password {
	String^ get ();
	void set (String^ value);
}
```

```
member Password : string with get, set
```

#### 属性值

类型：String

![](../icons/SectionExpanded.png)参见

#### 引用

[FtpClient 类](98289259-aeed-7095-b875-9d3ef78237fe.htm)

[HslCommunication.Enthernet.Ftp 命名空间](05837d2e-e0a8-4564-6d5e-a834d56ce5ce.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## SystemType 属性 

[原文連結](http://api.hslcommunication.cn/html/147c3982-4761-334d-9d78-41d745a8b753.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Enthernet.Ftp](../html/05837d2e-e0a8-4564-6d5e-a834d56ce5ce.htm "HslCommunication.Enthernet.Ftp")

[FtpClient 类](../html/98289259-aeed-7095-b875-9d3ef78237fe.htm "FtpClient 类")

[FtpClient 属性](../html/69a325b9-c5bb-9bf7-f4b4-847f840c54d4.htm "FtpClient 属性")

[Password 属性](../html/ec1a8da1-fb8a-54ef-aac4-53a04498488d.htm "Password 属性 ")

[SystemType 属性](../html/147c3982-4761-334d-9d78-41d745a8b753.htm "SystemType 属性 ")

[Username 属性](../html/bf1063c3-0f40-b74b-4323-29156b0d8f4e.htm "Username 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FtpClientSystemType 属性 |

系统的类型信息

**命名空间：**
 [HslCommunication.Enthernet.Ftp](05837d2e-e0a8-4564-6d5e-a834d56ce5ce.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public string SystemType { get; }
```

```
Public ReadOnly Property SystemType As String
	Get
```

```
public:
property String^ SystemType {
	String^ get ();
}
```

```
member SystemType : string with get
```

#### 属性值

类型：String

![](../icons/SectionExpanded.png)参见

#### 引用

[FtpClient 类](98289259-aeed-7095-b875-9d3ef78237fe.htm)

[HslCommunication.Enthernet.Ftp 命名空间](05837d2e-e0a8-4564-6d5e-a834d56ce5ce.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Username 属性 

[原文連結](http://api.hslcommunication.cn/html/bf1063c3-0f40-b74b-4323-29156b0d8f4e.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Enthernet.Ftp](../html/05837d2e-e0a8-4564-6d5e-a834d56ce5ce.htm "HslCommunication.Enthernet.Ftp")

[FtpClient 类](../html/98289259-aeed-7095-b875-9d3ef78237fe.htm "FtpClient 类")

[FtpClient 属性](../html/69a325b9-c5bb-9bf7-f4b4-847f840c54d4.htm "FtpClient 属性")

[Password 属性](../html/ec1a8da1-fb8a-54ef-aac4-53a04498488d.htm "Password 属性 ")

[SystemType 属性](../html/147c3982-4761-334d-9d78-41d745a8b753.htm "SystemType 属性 ")

[Username 属性](../html/bf1063c3-0f40-b74b-4323-29156b0d8f4e.htm "Username 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FtpClientUsername 属性 |

获取或设置登录的用户名信息，默认值为anonymous

**命名空间：**
 [HslCommunication.Enthernet.Ftp](05837d2e-e0a8-4564-6d5e-a834d56ce5ce.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public string Username { get; set; }
```

```
Public Property Username As String
	Get
	Set
```

```
public:
property String^ Username {
	String^ get ();
	void set (String^ value);
}
```

```
member Username : string with get, set
```

#### 属性值

类型：String

![](../icons/SectionExpanded.png)参见

#### 引用

[FtpClient 类](98289259-aeed-7095-b875-9d3ef78237fe.htm)

[HslCommunication.Enthernet.Ftp 命名空间](05837d2e-e0a8-4564-6d5e-a834d56ce5ce.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## FtpClient 方法

[原文連結](http://api.hslcommunication.cn/html/c26ace03-ff28-7070-8a9e-62b8abf31430.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Enthernet.Ftp](../html/05837d2e-e0a8-4564-6d5e-a834d56ce5ce.htm "HslCommunication.Enthernet.Ftp")

[FtpClient 类](../html/98289259-aeed-7095-b875-9d3ef78237fe.htm "FtpClient 类")

[FtpClient 方法](../html/c26ace03-ff28-7070-8a9e-62b8abf31430.htm "FtpClient 方法")

[CreateDirectory 方法](../html/82ac3dc9-8139-1130-727c-3df3f3d5ce75.htm "CreateDirectory 方法 ")

[DeleteDirectory 方法](../html/ad8fd03c-4875-6c6b-83a3-e7590c7891a3.htm "DeleteDirectory 方法 ")

[DeleteFile 方法](../html/e86754cd-59fd-e54b-be3d-18d3de267d61.htm "DeleteFile 方法 ")

[DownloadFile 方法](../html/41dcfcde-d66a-0549-86f2-0063edee69f3.htm "DownloadFile 方法 ")

[ExtraOnDisconnect 方法](../html/8113d5c7-4f60-fba4-0e43-c0dd4a778d59.htm "ExtraOnDisconnect 方法 ")

[GetLogTextFromBinary 方法](../html/051c0c92-9f31-e96b-26f7-3e1b1a63f272.htm "GetLogTextFromBinary 方法 ")

[GetNewNetMessage 方法](../html/9cb70414-4298-bcf3-fb20-79e4cb0393ad.htm "GetNewNetMessage 方法 ")

[InitializationOnConnect 方法](../html/56cbb726-5571-5c28-4c66-e1820eef5df8.htm "InitializationOnConnect 方法 ")

[ListFiles 方法](../html/c0783d51-47c7-c3d2-bf08-fa9afa5d50c9.htm "ListFiles 方法 ")

[RenameFile 方法](../html/d6aeb935-d84d-2f14-e447-201605bdaa27.htm "RenameFile 方法 ")

[RenameFolder 方法](../html/9348c583-a8a6-89f0-57e1-31ccae07c876.htm "RenameFolder 方法 ")

[ToString 方法](../html/83ab8b3b-8124-ea01-71fa-3bdbea2ba010.htm "ToString 方法 ")

[UploadFile 方法](../html/9b1e30fb-8019-f6fd-8891-6d05d1425af8.htm "UploadFile 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FtpClient 方法 |

[FtpClient](98289259-aeed-7095-b875-9d3ef78237fe.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法代码示例 | [ConnectClose](5ff1efa3-0f43-7515-add2-48681749ca7b.htm) | 手动断开与远程服务器的连接，如果当前是长连接模式，那么就会切换到短连接模式  Manually disconnect from the remote server, if it is currently in long connection mode, it will switch to short connection mode (继承自 [TcpNetCommunication](01e0c54c-8c1b-d3eb-4947-0d134d446266.htm)。) |
| 公共方法代码示例 | [ConnectCloseAsync](3b03af5f-46f8-e1a9-1964-228bac33c510.htm) | 手动断开与远程服务器的连接，如果当前是长连接模式，那么就会切换到短连接模式  Manually disconnect from the remote server, if it is currently in long connection mode, it will switch to short connection mode (继承自 [TcpNetCommunication](01e0c54c-8c1b-d3eb-4947-0d134d446266.htm)。) |
| 公共方法代码示例 | [ConnectServer](1dc20e7a-8929-bad2-d034-b2454d55a906.htm) | 尝试连接远程的服务器，如果连接成功，就切换短连接模式到长连接模式，后面的每次请求都共享一个通道，使得通讯速度更快速  Try to connect to a remote server. If the connection is successful, switch the short connection mode to the long connection mode. Each subsequent request will share a channel, making the communication speed faster. (继承自 [TcpNetCommunication](01e0c54c-8c1b-d3eb-4947-0d134d446266.htm)。) |
| 公共方法代码示例 | [ConnectServerAsync](9e07b3d8-a288-803d-c27c-1d75f4f92aa7.htm) | 尝试连接远程的服务器，如果连接成功，就切换短连接模式到长连接模式，后面的每次请求都共享一个通道，使得通讯速度更快速  Try to connect to a remote server. If the connection is successful, switch the short connection mode to the long connection mode. Each subsequent request will share a channel, making the communication speed faster. (继承自 [TcpNetCommunication](01e0c54c-8c1b-d3eb-4947-0d134d446266.htm)。) |
| 公共方法 | [CreateDirectory](82ac3dc9-8139-1130-727c-3df3f3d5ce75.htm) | 在FTP服务器上创建指定的目录  Create the specified directory on the FTP server. |
| 受保护的方法 | [DecideWhetherQAMessage](acc40cda-e310-8537-6685-59b67eeb16ff.htm) | 决定当前的消息是否是用于问答机制返回的消息，默认直接返回 true, 实际的情况需要根据协议进行重写方法  To determine whether the current message is the message returned by the question answering mechanism, the default is true. In actual cases, the rewriting method needs to be performed according to the protocol (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [DeleteDirectory](ad8fd03c-4875-6c6b-83a3-e7590c7891a3.htm) | 删除FTP服务器上的指定目录，并且可以指定是否删除所有子文件  Delete the specified directory on the FTP server, with the option to specify whether to delete all subfiles. |
| 公共方法 | [DeleteFile](e86754cd-59fd-e54b-be3d-18d3de267d61.htm) | 删除FTP服务器上的指定文件  Delete the specified file(s) on the FTP server. |
| 公共方法 | [DownloadFile](41dcfcde-d66a-0549-86f2-0063edee69f3.htm) | 从FTP服务器下载文件到本地文件，支持进度汇报和取消操作  Download files from an FTP server to local files, with support for progress reporting and cancellation operations. |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | [ExtraAfterReadFromCoreServer](06a49a15-a3ce-5e39-28c7-2b78851eca2e.htm) | 和服务器交互完成的时候调用的方法，可以根据读写结果进行一些额外的操作，具体的操作需要根据实际的需求来重写实现  The method called when the interaction with the server is completed can perform some additional operations based on the read and write results. The specific operations need to be rewritten according to actual needs. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法代码示例 | [ExtraOnDisconnect](8113d5c7-4f60-fba4-0e43-c0dd4a778d59.htm) | 根据实际的协议选择是否重写本方法，有些协议在断开连接之前，需要发送一些报文来关闭当前的网络通道  Select whether to rewrite this method according to the actual protocol. Some protocols need to send some packets to close the current network channel before disconnecting. (重写 [BinaryCommunicationExtraOnDisconnect](de9dbba2-3b2c-1c5e-53c8-4e6f8a44e2c4.htm).) |
| 受保护的方法代码示例 | [ExtraOnDisconnectAsync](210f593e-16d8-77c4-1178-58fbcfe9b27b.htm) | 根据实际的协议选择是否重写本方法，有些协议在断开连接之前，需要发送一些报文来关闭当前的网络通道  Select whether to rewrite this method according to the actual protocol. Some protocols need to send some packets to close the current network channel before disconnecting. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 受保护的方法 | [GetLogTextFromBinary](051c0c92-9f31-e96b-26f7-3e1b1a63f272.htm) | 获取当前的报文进行日志记录的时候，是否使用二进制的格式记录，默认返回 [LogMsgFormatBinary](2ed2ff87-caec-0e44-f08b-e2b015225a58.htm)，重写可以根据session对象分别返回不同记录模式  Whether to log the current packet in binary format, the default return is [LogMsgFormatBinary](2ed2ff87-caec-0e44-f08b-e2b015225a58.htm). If you want to override it, different recording modes can be returned according to session (重写 [BinaryCommunicationGetLogTextFromBinary(PipeSession, Byte)](67f3b2f7-957c-fa0b-8320-6799237157db.htm).) |
| 受保护的方法 | [GetNewNetMessage](9cb70414-4298-bcf3-fb20-79e4cb0393ad.htm) | 获取一个新的消息对象的方法，需要在继承类里面进行重写  The method to get a new message object needs to be overridden in the inheritance class (重写 [BinaryCommunicationGetNewNetMessage](16636a1b-a6f0-9c0f-9d9c-c94e26c67459.htm).) |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法代码示例 | [InitializationOnConnect](56cbb726-5571-5c28-4c66-e1820eef5df8.htm) | 根据实际的协议选择是否重写本方法，有些协议在创建连接之后，需要进行一些初始化的信号握手，才能最终建立网络通道。  Whether to rewrite this method is based on the actual protocol. Some protocols require some initial signal handshake to establish a network channel after the connection is created. (重写 [BinaryCommunicationInitializationOnConnect](4ab6989f-a5b5-05fe-8b44-bfa45824a2a7.htm).) |
| 受保护的方法代码示例 | [InitializationOnConnectAsync](7da2af04-d4fc-5045-19a4-4b0d1722eef3.htm) | 根据实际的协议选择是否重写本方法，有些协议在创建连接之后，需要进行一些初始化的信号握手，才能最终建立网络通道。  Whether to rewrite this method is based on the actual protocol. Some protocols require some initial signal handshake to establish a network channel after the connection is created. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [IpAddressPing](1011593d-fb48-acd1-44c2-27fd565a1180.htm) | 对当前设备的IP地址进行PING的操作，返回PING的结果，正常来说，返回Success  PING the IP address of the current device and return the PING result. Normally, it returns Success (继承自 [TcpNetCommunication](01e0c54c-8c1b-d3eb-4947-0d134d446266.htm)。) |
| 公共方法 | [ListFiles](c0783d51-47c7-c3d2-bf08-fa9afa5d50c9.htm) | 浏览指定目录下的文件及文件夹信息列表，默认根目录为 /   Browse the list of file and folder information in the specified directory. The default root directory is / |
| 受保护的方法 | [LogRevcMessage(Byte)](dab1657f-1e75-7fe1-7003-ccedfc69855d.htm) | 使用日志记录一个接收的报文信息  Logs are used to record information about a received packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [LogRevcMessage(Byte, PipeSession)](0f51ae45-575e-4228-f978-01b37d6d5df1.htm) | 使用日志记录一个接收的报文信息  Logs are used to record information about a received packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [LogSendMessage(Byte)](4f90eecc-68a2-1ef2-44ef-369f19bed3ec.htm) | 使用日志记录一个发送的报文信息  Logs are used to record information about a send packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [LogSendMessage(Byte, PipeSession)](209caadf-4f0f-40aa-04c2-ebc61f8dab13.htm) | 使用日志记录一个发送的报文信息  Logs are used to record information about a send packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法 | [PackCommandWithHeader](d0b2cd76-c306-11fc-38c7-af8ad306cce7.htm) | 对当前的命令进行打包处理，通常是携带命令头内容，标记当前的命令的长度信息，需要进行重写，否则默认不打包  The current command is packaged, usually carrying the content of the command header, marking the length of the current command, and it needs to be rewritten, otherwise it is not packaged by default (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServer(Byte)](8146dc0c-1edc-3ca7-428a-a3bdbd35e6c5.htm) | 将二进制的数据发送到管道中去，然后从管道里接收二进制的数据回来，并返回是否成功的结果对象。  Send binary data to the pipeline, and then receive binary data back from the pipeline, and return whether the success of the result object (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServer(IEnumerableByte)](044cfdd9-ec76-b89f-2f2b-5f439025006c.htm) | 将多个数据报文按顺序发到设备，并从设备接收返回的数据内容，然后拼接成一个Byte[]信息，需要重写[UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm)方法才能返回正确的结果。  Send multiple data packets to the device in sequence, and receive the returned data content from the device, and then splicing them into a Byte[] message, you need to rewrite [UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm) method to return the correct result. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServer(Byte, Boolean, Boolean)](b32ce8e8-a959-fcc2-6736-14b39fdb9a40.htm) | 将二进制的数据发送到管道中去，然后从管道里接收二进制的数据回来，并返回是否成功的结果对象。  Send binary data to the pipeline, and then receive binary data back from the pipeline, and return whether the success of the result object (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServer(CommunicationPipe, Byte, Boolean, Boolean)](ef36d56b-df09-f727-e8d3-316d37a13ee9.htm) | 使用指定的管道来进行数据通信，发送原始数据到管道，然后从管道接收相关的数据返回，本方法无锁 (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServerAsync(Byte)](a4de5a24-c67a-7b02-db05-33b85eaf6e1e.htm) | 将二进制的数据发送到管道中去，然后从管道里接收二进制的数据回来，并返回是否成功的结果对象。  Send binary data to the pipeline, and then receive binary data back from the pipeline, and return whether the success of the result object (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServerAsync(IEnumerableByte)](bcdb59b7-a70a-d5c0-1fc6-9ececcd39ea2.htm) | 将多个数据报文按顺序发到设备，并从设备接收返回的数据内容，然后拼接成一个Byte[]信息，需要重写[UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm)方法才能返回正确的结果。  Send multiple data packets to the device in sequence, and receive the returned data content from the device, and then splicing them into a Byte[] message, you need to rewrite [UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm) method to return the correct result. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServerAsync(Byte, Boolean, Boolean)](1c79c204-5ab5-9342-e65d-7455c20ea577.htm) | 将二进制的数据发送到管道中去，然后从管道里接收二进制的数据回来，并返回是否成功的结果对象。  Send binary data to the pipeline, and then receive binary data back from the pipeline, and return whether the success of the result object (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServerAsync(CommunicationPipe, Byte, Boolean, Boolean)](a254a19d-3c29-147b-335b-8adccdf4ee2f.htm) | 使用指定的管道来进行数据通信，发送原始数据到管道，然后从管道接收相关的数据返回，本方法无锁 (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [RenameFile](d6aeb935-d84d-2f14-e447-201605bdaa27.htm) | 重命名文件的名称  Rename the name of the file. |
| 公共方法 | [RenameFolder](9348c583-a8a6-89f0-57e1-31ccae07c876.htm) | 重命名目录，不支持跨目录重命名  Rename a directory; cross-directory renaming is not supported. |
| 公共方法 | [SetDtuPipe](f765d1d7-234a-b229-9457-97e25e024dc9.htm) | 将当前的通信对象设置DTU模式，允许传入现成的管道，并返回初始化结果，如果该设备重写了握手报文，就是返回握手结果  Set the current communication object to DTU mode, allow the existing pipe to be passed in, and return the initialization result, if the device rewrites the handshake packet, the handshake result is returned (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [SetDtuPipeAsync](9d1fd12d-b804-3828-a627-a2c0eef1d23f.htm) | 将当前的通信对象设置DTU模式，允许传入现成的管道，并返回初始化结果，如果该设备重写了握手报文，就是返回握手结果  Set the current communication object to DTU mode, allow the existing pipe to be passed in, and return the initialization result, if the device rewrites the handshake packet, the handshake result is returned (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ToString](83ab8b3b-8124-ea01-71fa-3bdbea2ba010.htm) | (重写 [TcpNetCommunicationToString](c1a8a99d-7115-55ca-57df-ece8e93e460c.htm).) |
| 公共方法 | [UnpackResponseContent](208017a7-8da2-33f8-3273-d7e7e850100f.htm) | 根据对方返回的报文命令，对命令进行基本的拆包，例如各种Modbus协议拆包为统一的核心报文，还支持对报文的验证  According to the message command returned by the other party, the command is basically unpacked, for example, various Modbus protocols are unpacked into a unified core message, and the verification of the message is also supported (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [UploadFile](9b1e30fb-8019-f6fd-8891-6d05d1425af8.htm) | 上传本地文件到FTP服务器，支持进度汇报和取消操作  Upload local files to an FTP server, with support for progress reporting and cancellation operations. |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[FtpClient 类](98289259-aeed-7095-b875-9d3ef78237fe.htm)

[HslCommunication.Enthernet.Ftp 命名空间](05837d2e-e0a8-4564-6d5e-a834d56ce5ce.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## CreateDirectory 方法 

[原文連結](http://api.hslcommunication.cn/html/82ac3dc9-8139-1130-727c-3df3f3d5ce75.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Enthernet.Ftp](../html/05837d2e-e0a8-4564-6d5e-a834d56ce5ce.htm "HslCommunication.Enthernet.Ftp")

[FtpClient 类](../html/98289259-aeed-7095-b875-9d3ef78237fe.htm "FtpClient 类")

[FtpClient 方法](../html/c26ace03-ff28-7070-8a9e-62b8abf31430.htm "FtpClient 方法")

[CreateDirectory 方法](../html/82ac3dc9-8139-1130-727c-3df3f3d5ce75.htm "CreateDirectory 方法 ")

[DeleteDirectory 方法](../html/ad8fd03c-4875-6c6b-83a3-e7590c7891a3.htm "DeleteDirectory 方法 ")

[DeleteFile 方法](../html/e86754cd-59fd-e54b-be3d-18d3de267d61.htm "DeleteFile 方法 ")

[DownloadFile 方法](../html/41dcfcde-d66a-0549-86f2-0063edee69f3.htm "DownloadFile 方法 ")

[ExtraOnDisconnect 方法](../html/8113d5c7-4f60-fba4-0e43-c0dd4a778d59.htm "ExtraOnDisconnect 方法 ")

[GetLogTextFromBinary 方法](../html/051c0c92-9f31-e96b-26f7-3e1b1a63f272.htm "GetLogTextFromBinary 方法 ")

[GetNewNetMessage 方法](../html/9cb70414-4298-bcf3-fb20-79e4cb0393ad.htm "GetNewNetMessage 方法 ")

[InitializationOnConnect 方法](../html/56cbb726-5571-5c28-4c66-e1820eef5df8.htm "InitializationOnConnect 方法 ")

[ListFiles 方法](../html/c0783d51-47c7-c3d2-bf08-fa9afa5d50c9.htm "ListFiles 方法 ")

[RenameFile 方法](../html/d6aeb935-d84d-2f14-e447-201605bdaa27.htm "RenameFile 方法 ")

[RenameFolder 方法](../html/9348c583-a8a6-89f0-57e1-31ccae07c876.htm "RenameFolder 方法 ")

[ToString 方法](../html/83ab8b3b-8124-ea01-71fa-3bdbea2ba010.htm "ToString 方法 ")

[UploadFile 方法](../html/9b1e30fb-8019-f6fd-8891-6d05d1425af8.htm "UploadFile 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FtpClientCreateDirectory 方法 |

在FTP服务器上创建指定的目录  
Create the specified directory on the FTP server.

**命名空间：**
 [HslCommunication.Enthernet.Ftp](05837d2e-e0a8-4564-6d5e-a834d56ce5ce.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public OperateResult CreateDirectory(
	string path
)
```

```
Public Function CreateDirectory ( 
	path As String
) As OperateResult
```

```
public:
OperateResult^ CreateDirectory(
	String^ path
)
```

```
member CreateDirectory : 
        path : string -> OperateResult 
```

#### 参数

path
:   类型：SystemString  
    指定的目录

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否创建目录成功

![](../icons/SectionExpanded.png)参见

#### 引用

[FtpClient 类](98289259-aeed-7095-b875-9d3ef78237fe.htm)

[HslCommunication.Enthernet.Ftp 命名空间](05837d2e-e0a8-4564-6d5e-a834d56ce5ce.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## DeleteDirectory 方法 

[原文連結](http://api.hslcommunication.cn/html/ad8fd03c-4875-6c6b-83a3-e7590c7891a3.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Enthernet.Ftp](../html/05837d2e-e0a8-4564-6d5e-a834d56ce5ce.htm "HslCommunication.Enthernet.Ftp")

[FtpClient 类](../html/98289259-aeed-7095-b875-9d3ef78237fe.htm "FtpClient 类")

[FtpClient 方法](../html/c26ace03-ff28-7070-8a9e-62b8abf31430.htm "FtpClient 方法")

[CreateDirectory 方法](../html/82ac3dc9-8139-1130-727c-3df3f3d5ce75.htm "CreateDirectory 方法 ")

[DeleteDirectory 方法](../html/ad8fd03c-4875-6c6b-83a3-e7590c7891a3.htm "DeleteDirectory 方法 ")

[DeleteFile 方法](../html/e86754cd-59fd-e54b-be3d-18d3de267d61.htm "DeleteFile 方法 ")

[DownloadFile 方法](../html/41dcfcde-d66a-0549-86f2-0063edee69f3.htm "DownloadFile 方法 ")

[ExtraOnDisconnect 方法](../html/8113d5c7-4f60-fba4-0e43-c0dd4a778d59.htm "ExtraOnDisconnect 方法 ")

[GetLogTextFromBinary 方法](../html/051c0c92-9f31-e96b-26f7-3e1b1a63f272.htm "GetLogTextFromBinary 方法 ")

[GetNewNetMessage 方法](../html/9cb70414-4298-bcf3-fb20-79e4cb0393ad.htm "GetNewNetMessage 方法 ")

[InitializationOnConnect 方法](../html/56cbb726-5571-5c28-4c66-e1820eef5df8.htm "InitializationOnConnect 方法 ")

[ListFiles 方法](../html/c0783d51-47c7-c3d2-bf08-fa9afa5d50c9.htm "ListFiles 方法 ")

[RenameFile 方法](../html/d6aeb935-d84d-2f14-e447-201605bdaa27.htm "RenameFile 方法 ")

[RenameFolder 方法](../html/9348c583-a8a6-89f0-57e1-31ccae07c876.htm "RenameFolder 方法 ")

[ToString 方法](../html/83ab8b3b-8124-ea01-71fa-3bdbea2ba010.htm "ToString 方法 ")

[UploadFile 方法](../html/9b1e30fb-8019-f6fd-8891-6d05d1425af8.htm "UploadFile 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FtpClientDeleteDirectory 方法 |

删除FTP服务器上的指定目录，并且可以指定是否删除所有子文件  
Delete the specified directory on the FTP server, with the option to specify whether to delete all subfiles.

**命名空间：**
 [HslCommunication.Enthernet.Ftp](05837d2e-e0a8-4564-6d5e-a834d56ce5ce.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public OperateResult DeleteDirectory(
	string path,
	bool clearFiles = false
)
```

```
Public Function DeleteDirectory ( 
	path As String,
	Optional clearFiles As Boolean = false
) As OperateResult
```

```
public:
OperateResult^ DeleteDirectory(
	String^ path, 
	bool clearFiles = false
)
```

```
member DeleteDirectory : 
        path : string * 
        ?clearFiles : bool 
(* Defaults:
        let _clearFiles = defaultArg clearFiles false
*)
-> OperateResult 
```

#### 参数

path
:   类型：SystemString  
    Ftp服务器上的目录

clearFiles (Optional)
:   类型：SystemBoolean  
    是否删除所有的文件

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否删除成功

![](../icons/SectionExpanded.png)参见

#### 引用

[FtpClient 类](98289259-aeed-7095-b875-9d3ef78237fe.htm)

[HslCommunication.Enthernet.Ftp 命名空间](05837d2e-e0a8-4564-6d5e-a834d56ce5ce.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## DeleteFile 方法 

[原文連結](http://api.hslcommunication.cn/html/e86754cd-59fd-e54b-be3d-18d3de267d61.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Enthernet.Ftp](../html/05837d2e-e0a8-4564-6d5e-a834d56ce5ce.htm "HslCommunication.Enthernet.Ftp")

[FtpClient 类](../html/98289259-aeed-7095-b875-9d3ef78237fe.htm "FtpClient 类")

[FtpClient 方法](../html/c26ace03-ff28-7070-8a9e-62b8abf31430.htm "FtpClient 方法")

[CreateDirectory 方法](../html/82ac3dc9-8139-1130-727c-3df3f3d5ce75.htm "CreateDirectory 方法 ")

[DeleteDirectory 方法](../html/ad8fd03c-4875-6c6b-83a3-e7590c7891a3.htm "DeleteDirectory 方法 ")

[DeleteFile 方法](../html/e86754cd-59fd-e54b-be3d-18d3de267d61.htm "DeleteFile 方法 ")

[DownloadFile 方法](../html/41dcfcde-d66a-0549-86f2-0063edee69f3.htm "DownloadFile 方法 ")

[ExtraOnDisconnect 方法](../html/8113d5c7-4f60-fba4-0e43-c0dd4a778d59.htm "ExtraOnDisconnect 方法 ")

[GetLogTextFromBinary 方法](../html/051c0c92-9f31-e96b-26f7-3e1b1a63f272.htm "GetLogTextFromBinary 方法 ")

[GetNewNetMessage 方法](../html/9cb70414-4298-bcf3-fb20-79e4cb0393ad.htm "GetNewNetMessage 方法 ")

[InitializationOnConnect 方法](../html/56cbb726-5571-5c28-4c66-e1820eef5df8.htm "InitializationOnConnect 方法 ")

[ListFiles 方法](../html/c0783d51-47c7-c3d2-bf08-fa9afa5d50c9.htm "ListFiles 方法 ")

[RenameFile 方法](../html/d6aeb935-d84d-2f14-e447-201605bdaa27.htm "RenameFile 方法 ")

[RenameFolder 方法](../html/9348c583-a8a6-89f0-57e1-31ccae07c876.htm "RenameFolder 方法 ")

[ToString 方法](../html/83ab8b3b-8124-ea01-71fa-3bdbea2ba010.htm "ToString 方法 ")

[UploadFile 方法](../html/9b1e30fb-8019-f6fd-8891-6d05d1425af8.htm "UploadFile 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FtpClientDeleteFile 方法 |

删除FTP服务器上的指定文件  
Delete the specified file(s) on the FTP server.

**命名空间：**
 [HslCommunication.Enthernet.Ftp](05837d2e-e0a8-4564-6d5e-a834d56ce5ce.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public OperateResult DeleteFile(
	string remoteFile
)
```

```
Public Function DeleteFile ( 
	remoteFile As String
) As OperateResult
```

```
public:
OperateResult^ DeleteFile(
	String^ remoteFile
)
```

```
member DeleteFile : 
        remoteFile : string -> OperateResult 
```

#### 参数

remoteFile
:   类型：SystemString  
    服务器上的文件，例如 /A/123.txt

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否删除成功

![](../icons/SectionExpanded.png)参见

#### 引用

[FtpClient 类](98289259-aeed-7095-b875-9d3ef78237fe.htm)

[HslCommunication.Enthernet.Ftp 命名空间](05837d2e-e0a8-4564-6d5e-a834d56ce5ce.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## DownloadFile 方法 

[原文連結](http://api.hslcommunication.cn/html/41dcfcde-d66a-0549-86f2-0063edee69f3.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Enthernet.Ftp](../html/05837d2e-e0a8-4564-6d5e-a834d56ce5ce.htm "HslCommunication.Enthernet.Ftp")

[FtpClient 类](../html/98289259-aeed-7095-b875-9d3ef78237fe.htm "FtpClient 类")

[FtpClient 方法](../html/c26ace03-ff28-7070-8a9e-62b8abf31430.htm "FtpClient 方法")

[CreateDirectory 方法](../html/82ac3dc9-8139-1130-727c-3df3f3d5ce75.htm "CreateDirectory 方法 ")

[DeleteDirectory 方法](../html/ad8fd03c-4875-6c6b-83a3-e7590c7891a3.htm "DeleteDirectory 方法 ")

[DeleteFile 方法](../html/e86754cd-59fd-e54b-be3d-18d3de267d61.htm "DeleteFile 方法 ")

[DownloadFile 方法](../html/41dcfcde-d66a-0549-86f2-0063edee69f3.htm "DownloadFile 方法 ")

[ExtraOnDisconnect 方法](../html/8113d5c7-4f60-fba4-0e43-c0dd4a778d59.htm "ExtraOnDisconnect 方法 ")

[GetLogTextFromBinary 方法](../html/051c0c92-9f31-e96b-26f7-3e1b1a63f272.htm "GetLogTextFromBinary 方法 ")

[GetNewNetMessage 方法](../html/9cb70414-4298-bcf3-fb20-79e4cb0393ad.htm "GetNewNetMessage 方法 ")

[InitializationOnConnect 方法](../html/56cbb726-5571-5c28-4c66-e1820eef5df8.htm "InitializationOnConnect 方法 ")

[ListFiles 方法](../html/c0783d51-47c7-c3d2-bf08-fa9afa5d50c9.htm "ListFiles 方法 ")

[RenameFile 方法](../html/d6aeb935-d84d-2f14-e447-201605bdaa27.htm "RenameFile 方法 ")

[RenameFolder 方法](../html/9348c583-a8a6-89f0-57e1-31ccae07c876.htm "RenameFolder 方法 ")

[ToString 方法](../html/83ab8b3b-8124-ea01-71fa-3bdbea2ba010.htm "ToString 方法 ")

[UploadFile 方法](../html/9b1e30fb-8019-f6fd-8891-6d05d1425af8.htm "UploadFile 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FtpClientDownloadFile 方法 |

从FTP服务器下载文件到本地文件，支持进度汇报和取消操作  
Download files from an FTP server to local files, with support for progress reporting and cancellation operations.

**命名空间：**
 [HslCommunication.Enthernet.Ftp](05837d2e-e0a8-4564-6d5e-a834d56ce5ce.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public OperateResult DownloadFile(
	string remoteFile,
	string localFile,
	Action<long, long> processReport = null,
	HslCancelToken cancelToken = null
)
```

```
Public Function DownloadFile ( 
	remoteFile As String,
	localFile As String,
	Optional processReport As Action(Of Long, Long) = Nothing,
	Optional cancelToken As HslCancelToken = Nothing
) As OperateResult
```

```
public:
OperateResult^ DownloadFile(
	String^ remoteFile, 
	String^ localFile, 
	Action<long long, long long>^ processReport = nullptr, 
	HslCancelToken^ cancelToken = nullptr
)
```

```
member DownloadFile : 
        remoteFile : string * 
        localFile : string * 
        ?processReport : Action<int64, int64> * 
        ?cancelToken : HslCancelToken 
(* Defaults:
        let _processReport = defaultArg processReport null
        let _cancelToken = defaultArg cancelToken null
*)
-> OperateResult 
```

#### 参数

remoteFile
:   类型：SystemString  
    服务器的文件，例如 /A/123.txt

localFile
:   类型：SystemString  
    本地保存的文件名

processReport (Optional)
:   类型：SystemActionInt64, Int64  
    进度报告

cancelToken (Optional)
:   类型：[HslCommunication.CoreHslCancelToken](072498cc-b288-38cf-c2d2-7d9452ac6a3e.htm)  
    取消操作

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否下载成功

![](../icons/SectionExpanded.png)参见

#### 引用

[FtpClient 类](98289259-aeed-7095-b875-9d3ef78237fe.htm)

[HslCommunication.Enthernet.Ftp 命名空间](05837d2e-e0a8-4564-6d5e-a834d56ce5ce.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ExtraOnDisconnect 方法 

[原文連結](http://api.hslcommunication.cn/html/8113d5c7-4f60-fba4-0e43-c0dd4a778d59.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Enthernet.Ftp](../html/05837d2e-e0a8-4564-6d5e-a834d56ce5ce.htm "HslCommunication.Enthernet.Ftp")

[FtpClient 类](../html/98289259-aeed-7095-b875-9d3ef78237fe.htm "FtpClient 类")

[FtpClient 方法](../html/c26ace03-ff28-7070-8a9e-62b8abf31430.htm "FtpClient 方法")

[CreateDirectory 方法](../html/82ac3dc9-8139-1130-727c-3df3f3d5ce75.htm "CreateDirectory 方法 ")

[DeleteDirectory 方法](../html/ad8fd03c-4875-6c6b-83a3-e7590c7891a3.htm "DeleteDirectory 方法 ")

[DeleteFile 方法](../html/e86754cd-59fd-e54b-be3d-18d3de267d61.htm "DeleteFile 方法 ")

[DownloadFile 方法](../html/41dcfcde-d66a-0549-86f2-0063edee69f3.htm "DownloadFile 方法 ")

[ExtraOnDisconnect 方法](../html/8113d5c7-4f60-fba4-0e43-c0dd4a778d59.htm "ExtraOnDisconnect 方法 ")

[GetLogTextFromBinary 方法](../html/051c0c92-9f31-e96b-26f7-3e1b1a63f272.htm "GetLogTextFromBinary 方法 ")

[GetNewNetMessage 方法](../html/9cb70414-4298-bcf3-fb20-79e4cb0393ad.htm "GetNewNetMessage 方法 ")

[InitializationOnConnect 方法](../html/56cbb726-5571-5c28-4c66-e1820eef5df8.htm "InitializationOnConnect 方法 ")

[ListFiles 方法](../html/c0783d51-47c7-c3d2-bf08-fa9afa5d50c9.htm "ListFiles 方法 ")

[RenameFile 方法](../html/d6aeb935-d84d-2f14-e447-201605bdaa27.htm "RenameFile 方法 ")

[RenameFolder 方法](../html/9348c583-a8a6-89f0-57e1-31ccae07c876.htm "RenameFolder 方法 ")

[ToString 方法](../html/83ab8b3b-8124-ea01-71fa-3bdbea2ba010.htm "ToString 方法 ")

[UploadFile 方法](../html/9b1e30fb-8019-f6fd-8891-6d05d1425af8.htm "UploadFile 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FtpClientExtraOnDisconnect 方法 |

根据实际的协议选择是否重写本方法，有些协议在断开连接之前，需要发送一些报文来关闭当前的网络通道  
Select whether to rewrite this method according to the actual protocol. Some protocols need to send some packets to close the current network channel before disconnecting.

**命名空间：**
 [HslCommunication.Enthernet.Ftp](05837d2e-e0a8-4564-6d5e-a834d56ce5ce.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
protected override OperateResult ExtraOnDisconnect()
```

```
Protected Overrides Function ExtraOnDisconnect As OperateResult
```

```
protected:
virtual OperateResult^ ExtraOnDisconnect() override
```

```
abstract ExtraOnDisconnect : unit -> OperateResult 
override ExtraOnDisconnect : unit -> OperateResult
```

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
当断开连接时额外的操作结果

![](../icons/SectionExpanded.png)示例

目前暂无相关的示例，组件支持的协议都不用实现这个方法。

![](../icons/SectionExpanded.png)参见

#### 引用

[FtpClient 类](98289259-aeed-7095-b875-9d3ef78237fe.htm)

[HslCommunication.Enthernet.Ftp 命名空间](05837d2e-e0a8-4564-6d5e-a834d56ce5ce.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GetLogTextFromBinary 方法 

[原文連結](http://api.hslcommunication.cn/html/051c0c92-9f31-e96b-26f7-3e1b1a63f272.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Enthernet.Ftp](../html/05837d2e-e0a8-4564-6d5e-a834d56ce5ce.htm "HslCommunication.Enthernet.Ftp")

[FtpClient 类](../html/98289259-aeed-7095-b875-9d3ef78237fe.htm "FtpClient 类")

[FtpClient 方法](../html/c26ace03-ff28-7070-8a9e-62b8abf31430.htm "FtpClient 方法")

[CreateDirectory 方法](../html/82ac3dc9-8139-1130-727c-3df3f3d5ce75.htm "CreateDirectory 方法 ")

[DeleteDirectory 方法](../html/ad8fd03c-4875-6c6b-83a3-e7590c7891a3.htm "DeleteDirectory 方法 ")

[DeleteFile 方法](../html/e86754cd-59fd-e54b-be3d-18d3de267d61.htm "DeleteFile 方法 ")

[DownloadFile 方法](../html/41dcfcde-d66a-0549-86f2-0063edee69f3.htm "DownloadFile 方法 ")

[ExtraOnDisconnect 方法](../html/8113d5c7-4f60-fba4-0e43-c0dd4a778d59.htm "ExtraOnDisconnect 方法 ")

[GetLogTextFromBinary 方法](../html/051c0c92-9f31-e96b-26f7-3e1b1a63f272.htm "GetLogTextFromBinary 方法 ")

[GetNewNetMessage 方法](../html/9cb70414-4298-bcf3-fb20-79e4cb0393ad.htm "GetNewNetMessage 方法 ")

[InitializationOnConnect 方法](../html/56cbb726-5571-5c28-4c66-e1820eef5df8.htm "InitializationOnConnect 方法 ")

[ListFiles 方法](../html/c0783d51-47c7-c3d2-bf08-fa9afa5d50c9.htm "ListFiles 方法 ")

[RenameFile 方法](../html/d6aeb935-d84d-2f14-e447-201605bdaa27.htm "RenameFile 方法 ")

[RenameFolder 方法](../html/9348c583-a8a6-89f0-57e1-31ccae07c876.htm "RenameFolder 方法 ")

[ToString 方法](../html/83ab8b3b-8124-ea01-71fa-3bdbea2ba010.htm "ToString 方法 ")

[UploadFile 方法](../html/9b1e30fb-8019-f6fd-8891-6d05d1425af8.htm "UploadFile 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FtpClientGetLogTextFromBinary 方法 |

获取当前的报文进行日志记录的时候，是否使用二进制的格式记录，默认返回 [LogMsgFormatBinary](2ed2ff87-caec-0e44-f08b-e2b015225a58.htm)，重写可以根据session对象分别返回不同记录模式  
Whether to log the current packet in binary format, the default return is [LogMsgFormatBinary](2ed2ff87-caec-0e44-f08b-e2b015225a58.htm). If you want to override it,
different recording modes can be returned according to session

**命名空间：**
 [HslCommunication.Enthernet.Ftp](05837d2e-e0a8-4564-6d5e-a834d56ce5ce.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
protected override string GetLogTextFromBinary(
	PipeSession session,
	byte[] content
)
```

```
Protected Overrides Function GetLogTextFromBinary ( 
	session As PipeSession,
	content As Byte()
) As String
```

```
protected:
virtual String^ GetLogTextFromBinary(
	PipeSession^ session, 
	array<unsigned char>^ content
) override
```

```
abstract GetLogTextFromBinary : 
        session : PipeSession * 
        content : byte[] -> string 
override GetLogTextFromBinary : 
        session : PipeSession * 
        content : byte[] -> string
```

#### 参数

session
:   类型：[HslCommunication.Core.NetPipeSession](bb565655-08a0-c9eb-9c75-26ad6a263256.htm)  
    会话对象

content
:   类型：SystemByte  
    等待记录的字节消息内容

#### 返回值

类型：String  
是否二进制记录报文格式

![](../icons/SectionExpanded.png)参见

#### 引用

[FtpClient 类](98289259-aeed-7095-b875-9d3ef78237fe.htm)

[HslCommunication.Enthernet.Ftp 命名空间](05837d2e-e0a8-4564-6d5e-a834d56ce5ce.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GetNewNetMessage 方法 

[原文連結](http://api.hslcommunication.cn/html/9cb70414-4298-bcf3-fb20-79e4cb0393ad.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Enthernet.Ftp](../html/05837d2e-e0a8-4564-6d5e-a834d56ce5ce.htm "HslCommunication.Enthernet.Ftp")

[FtpClient 类](../html/98289259-aeed-7095-b875-9d3ef78237fe.htm "FtpClient 类")

[FtpClient 方法](../html/c26ace03-ff28-7070-8a9e-62b8abf31430.htm "FtpClient 方法")

[CreateDirectory 方法](../html/82ac3dc9-8139-1130-727c-3df3f3d5ce75.htm "CreateDirectory 方法 ")

[DeleteDirectory 方法](../html/ad8fd03c-4875-6c6b-83a3-e7590c7891a3.htm "DeleteDirectory 方法 ")

[DeleteFile 方法](../html/e86754cd-59fd-e54b-be3d-18d3de267d61.htm "DeleteFile 方法 ")

[DownloadFile 方法](../html/41dcfcde-d66a-0549-86f2-0063edee69f3.htm "DownloadFile 方法 ")

[ExtraOnDisconnect 方法](../html/8113d5c7-4f60-fba4-0e43-c0dd4a778d59.htm "ExtraOnDisconnect 方法 ")

[GetLogTextFromBinary 方法](../html/051c0c92-9f31-e96b-26f7-3e1b1a63f272.htm "GetLogTextFromBinary 方法 ")

[GetNewNetMessage 方法](../html/9cb70414-4298-bcf3-fb20-79e4cb0393ad.htm "GetNewNetMessage 方法 ")

[InitializationOnConnect 方法](../html/56cbb726-5571-5c28-4c66-e1820eef5df8.htm "InitializationOnConnect 方法 ")

[ListFiles 方法](../html/c0783d51-47c7-c3d2-bf08-fa9afa5d50c9.htm "ListFiles 方法 ")

[RenameFile 方法](../html/d6aeb935-d84d-2f14-e447-201605bdaa27.htm "RenameFile 方法 ")

[RenameFolder 方法](../html/9348c583-a8a6-89f0-57e1-31ccae07c876.htm "RenameFolder 方法 ")

[ToString 方法](../html/83ab8b3b-8124-ea01-71fa-3bdbea2ba010.htm "ToString 方法 ")

[UploadFile 方法](../html/9b1e30fb-8019-f6fd-8891-6d05d1425af8.htm "UploadFile 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FtpClientGetNewNetMessage 方法 |

获取一个新的消息对象的方法，需要在继承类里面进行重写  
The method to get a new message object needs to be overridden in the inheritance class

**命名空间：**
 [HslCommunication.Enthernet.Ftp](05837d2e-e0a8-4564-6d5e-a834d56ce5ce.htm)  
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

[FtpClient 类](98289259-aeed-7095-b875-9d3ef78237fe.htm)

[HslCommunication.Enthernet.Ftp 命名空间](05837d2e-e0a8-4564-6d5e-a834d56ce5ce.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## InitializationOnConnect 方法 

[原文連結](http://api.hslcommunication.cn/html/56cbb726-5571-5c28-4c66-e1820eef5df8.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Enthernet.Ftp](../html/05837d2e-e0a8-4564-6d5e-a834d56ce5ce.htm "HslCommunication.Enthernet.Ftp")

[FtpClient 类](../html/98289259-aeed-7095-b875-9d3ef78237fe.htm "FtpClient 类")

[FtpClient 方法](../html/c26ace03-ff28-7070-8a9e-62b8abf31430.htm "FtpClient 方法")

[CreateDirectory 方法](../html/82ac3dc9-8139-1130-727c-3df3f3d5ce75.htm "CreateDirectory 方法 ")

[DeleteDirectory 方法](../html/ad8fd03c-4875-6c6b-83a3-e7590c7891a3.htm "DeleteDirectory 方法 ")

[DeleteFile 方法](../html/e86754cd-59fd-e54b-be3d-18d3de267d61.htm "DeleteFile 方法 ")

[DownloadFile 方法](../html/41dcfcde-d66a-0549-86f2-0063edee69f3.htm "DownloadFile 方法 ")

[ExtraOnDisconnect 方法](../html/8113d5c7-4f60-fba4-0e43-c0dd4a778d59.htm "ExtraOnDisconnect 方法 ")

[GetLogTextFromBinary 方法](../html/051c0c92-9f31-e96b-26f7-3e1b1a63f272.htm "GetLogTextFromBinary 方法 ")

[GetNewNetMessage 方法](../html/9cb70414-4298-bcf3-fb20-79e4cb0393ad.htm "GetNewNetMessage 方法 ")

[InitializationOnConnect 方法](../html/56cbb726-5571-5c28-4c66-e1820eef5df8.htm "InitializationOnConnect 方法 ")

[ListFiles 方法](../html/c0783d51-47c7-c3d2-bf08-fa9afa5d50c9.htm "ListFiles 方法 ")

[RenameFile 方法](../html/d6aeb935-d84d-2f14-e447-201605bdaa27.htm "RenameFile 方法 ")

[RenameFolder 方法](../html/9348c583-a8a6-89f0-57e1-31ccae07c876.htm "RenameFolder 方法 ")

[ToString 方法](../html/83ab8b3b-8124-ea01-71fa-3bdbea2ba010.htm "ToString 方法 ")

[UploadFile 方法](../html/9b1e30fb-8019-f6fd-8891-6d05d1425af8.htm "UploadFile 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FtpClientInitializationOnConnect 方法 |

根据实际的协议选择是否重写本方法，有些协议在创建连接之后，需要进行一些初始化的信号握手，才能最终建立网络通道。  
Whether to rewrite this method is based on the actual protocol. Some protocols require some initial signal handshake to establish a network channel after the connection is created.

**命名空间：**
 [HslCommunication.Enthernet.Ftp](05837d2e-e0a8-4564-6d5e-a834d56ce5ce.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
protected override OperateResult InitializationOnConnect()
```

```
Protected Overrides Function InitializationOnConnect As OperateResult
```

```
protected:
virtual OperateResult^ InitializationOnConnect() override
```

```
abstract InitializationOnConnect : unit -> OperateResult 
override InitializationOnConnect : unit -> OperateResult
```

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否初始化成功，依据具体的协议进行重写

![](../icons/SectionExpanded.png)示例

有些协议不需要握手信号，比如三菱的MC协议，Modbus协议，西门子和欧姆龙就存在握手信息，此处的例子是继承本类后重写的西门子的协议示例

西门子重连示例

[复制](# "复制")

```
        /// <inheritdoc/>
        public override OperateResult<byte[]> ReadFromCoreServer( CommunicationPipe pipe, byte[] send, bool hasResponseData, bool usePackAndUnpack )
        {
            while (true)
            {
                OperateResult<byte[]> read = base.ReadFromCoreServer( pipe, send, hasResponseData, usePackAndUnpack );
                if (!read.IsSuccess) return read;

                if (read.Content?.Length >= 4)
                {
                    if ((read.Content[2] * 256 + read.Content[3]) != 0x07) return read;
                }
            }
        }

        /// <inheritdoc/>
        protected override OperateResult InitializationOnConnect( )
        {
            // 第一次握手 -> First handshake
            OperateResult<byte[]> read_first = ReadFromCoreServer( this.CommunicationPipe, plcHead1, hasResponseData: true, usePackAndUnpack: true );
            if (!read_first.IsSuccess) return read_first;

            // 第二次握手 -> Second handshake
            OperateResult<byte[]> read_second = ReadFromCoreServer( this.CommunicationPipe, plcHead2, hasResponseData: true, usePackAndUnpack: true );
            if (!read_second.IsSuccess) return read_second;

            // 调整单次接收的pdu长度信息
            pdu_length = ByteTransform.TransUInt16( read_second.Content.SelectLast( 2 ), 0 ) - 28;
            if (pdu_length < 200) pdu_length = 200;

            incrementCount = new SoftIncrementCount( ushort.MaxValue, 1 );
            // 返回成功的信号 -> Return a successful signal
            return OperateResult.CreateSuccessResult( );
        }

#if !NET35 && !NET20
        /// <inheritdoc/>
        public async override Task<OperateResult<byte[]>> ReadFromCoreServerAsync( CommunicationPipe pipe, byte[] send, bool hasResponseData, bool usePackAndUnpack )
        {
            while (true)
            {
                OperateResult<byte[]> read = await base.ReadFromCoreServerAsync( pipe, send, hasResponseData, usePackAndUnpack );
                if (!read.IsSuccess) return read;

                if (read.Content?.Length >= 4)
                {
                    if ((read.Content[2] * 256 + read.Content[3]) != 0x07) return read;
                }
            }
        }

        /// <inheritdoc/>
        protected override async Task<OperateResult> InitializationOnConnectAsync( )
        {
            // 第一次握手 -> First handshake
            OperateResult<byte[]> read_first = await ReadFromCoreServerAsync( this.CommunicationPipe, plcHead1, hasResponseData: true, usePackAndUnpack: true );
            if (!read_first.IsSuccess) return read_first;

            // 第二次握手 -> Second handshake
            OperateResult<byte[]> read_second = await ReadFromCoreServerAsync( this.CommunicationPipe, plcHead2, hasResponseData: true, usePackAndUnpack: true );
            if (!read_second.IsSuccess) return read_second;

            // 调整单次接收的pdu长度信息
            pdu_length = ByteTransform.TransUInt16( read_second.Content.SelectLast( 2 ), 0 ) - 28;
            if (pdu_length < 200) pdu_length = 200;

            incrementCount = new SoftIncrementCount( ushort.MaxValue, 1 );
            // 返回成功的信号 -> Return a successful signal
            return OperateResult.CreateSuccessResult( );
        }
#endif
```

![](../icons/SectionExpanded.png)参见

#### 引用

[FtpClient 类](98289259-aeed-7095-b875-9d3ef78237fe.htm)

[HslCommunication.Enthernet.Ftp 命名空间](05837d2e-e0a8-4564-6d5e-a834d56ce5ce.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ListFiles 方法 

[原文連結](http://api.hslcommunication.cn/html/c0783d51-47c7-c3d2-bf08-fa9afa5d50c9.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Enthernet.Ftp](../html/05837d2e-e0a8-4564-6d5e-a834d56ce5ce.htm "HslCommunication.Enthernet.Ftp")

[FtpClient 类](../html/98289259-aeed-7095-b875-9d3ef78237fe.htm "FtpClient 类")

[FtpClient 方法](../html/c26ace03-ff28-7070-8a9e-62b8abf31430.htm "FtpClient 方法")

[CreateDirectory 方法](../html/82ac3dc9-8139-1130-727c-3df3f3d5ce75.htm "CreateDirectory 方法 ")

[DeleteDirectory 方法](../html/ad8fd03c-4875-6c6b-83a3-e7590c7891a3.htm "DeleteDirectory 方法 ")

[DeleteFile 方法](../html/e86754cd-59fd-e54b-be3d-18d3de267d61.htm "DeleteFile 方法 ")

[DownloadFile 方法](../html/41dcfcde-d66a-0549-86f2-0063edee69f3.htm "DownloadFile 方法 ")

[ExtraOnDisconnect 方法](../html/8113d5c7-4f60-fba4-0e43-c0dd4a778d59.htm "ExtraOnDisconnect 方法 ")

[GetLogTextFromBinary 方法](../html/051c0c92-9f31-e96b-26f7-3e1b1a63f272.htm "GetLogTextFromBinary 方法 ")

[GetNewNetMessage 方法](../html/9cb70414-4298-bcf3-fb20-79e4cb0393ad.htm "GetNewNetMessage 方法 ")

[InitializationOnConnect 方法](../html/56cbb726-5571-5c28-4c66-e1820eef5df8.htm "InitializationOnConnect 方法 ")

[ListFiles 方法](../html/c0783d51-47c7-c3d2-bf08-fa9afa5d50c9.htm "ListFiles 方法 ")

[RenameFile 方法](../html/d6aeb935-d84d-2f14-e447-201605bdaa27.htm "RenameFile 方法 ")

[RenameFolder 方法](../html/9348c583-a8a6-89f0-57e1-31ccae07c876.htm "RenameFolder 方法 ")

[ToString 方法](../html/83ab8b3b-8124-ea01-71fa-3bdbea2ba010.htm "ToString 方法 ")

[UploadFile 方法](../html/9b1e30fb-8019-f6fd-8891-6d05d1425af8.htm "UploadFile 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FtpClientListFiles 方法 |

浏览指定目录下的文件及文件夹信息列表，默认根目录为 /   
Browse the list of file and folder information in the specified directory. The default root directory is /

**命名空间：**
 [HslCommunication.Enthernet.Ftp](05837d2e-e0a8-4564-6d5e-a834d56ce5ce.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public OperateResult<FtpFileItem[]> ListFiles(
	string parent = "/"
)
```

```
Public Function ListFiles ( 
	Optional parent As String = "/"
) As OperateResult(Of FtpFileItem())
```

```
public:
OperateResult<array<FtpFileItem^>^>^ ListFiles(
	String^ parent = L"/"
)
```

```
member ListFiles : 
        ?parent : string 
(* Defaults:
        let _parent = defaultArg parent "/"
*)
-> OperateResult<FtpFileItem[]> 
```

#### 参数

parent (Optional)
:   类型：SystemString  
    指定的目录

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)[FtpFileItem](97ec5ae6-c5a1-0122-982b-cb02f1031c96.htm)  
文件列表信息

![](../icons/SectionExpanded.png)参见

#### 引用

[FtpClient 类](98289259-aeed-7095-b875-9d3ef78237fe.htm)

[HslCommunication.Enthernet.Ftp 命名空间](05837d2e-e0a8-4564-6d5e-a834d56ce5ce.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## RenameFile 方法 

[原文連結](http://api.hslcommunication.cn/html/d6aeb935-d84d-2f14-e447-201605bdaa27.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Enthernet.Ftp](../html/05837d2e-e0a8-4564-6d5e-a834d56ce5ce.htm "HslCommunication.Enthernet.Ftp")

[FtpClient 类](../html/98289259-aeed-7095-b875-9d3ef78237fe.htm "FtpClient 类")

[FtpClient 方法](../html/c26ace03-ff28-7070-8a9e-62b8abf31430.htm "FtpClient 方法")

[CreateDirectory 方法](../html/82ac3dc9-8139-1130-727c-3df3f3d5ce75.htm "CreateDirectory 方法 ")

[DeleteDirectory 方法](../html/ad8fd03c-4875-6c6b-83a3-e7590c7891a3.htm "DeleteDirectory 方法 ")

[DeleteFile 方法](../html/e86754cd-59fd-e54b-be3d-18d3de267d61.htm "DeleteFile 方法 ")

[DownloadFile 方法](../html/41dcfcde-d66a-0549-86f2-0063edee69f3.htm "DownloadFile 方法 ")

[ExtraOnDisconnect 方法](../html/8113d5c7-4f60-fba4-0e43-c0dd4a778d59.htm "ExtraOnDisconnect 方法 ")

[GetLogTextFromBinary 方法](../html/051c0c92-9f31-e96b-26f7-3e1b1a63f272.htm "GetLogTextFromBinary 方法 ")

[GetNewNetMessage 方法](../html/9cb70414-4298-bcf3-fb20-79e4cb0393ad.htm "GetNewNetMessage 方法 ")

[InitializationOnConnect 方法](../html/56cbb726-5571-5c28-4c66-e1820eef5df8.htm "InitializationOnConnect 方法 ")

[ListFiles 方法](../html/c0783d51-47c7-c3d2-bf08-fa9afa5d50c9.htm "ListFiles 方法 ")

[RenameFile 方法](../html/d6aeb935-d84d-2f14-e447-201605bdaa27.htm "RenameFile 方法 ")

[RenameFolder 方法](../html/9348c583-a8a6-89f0-57e1-31ccae07c876.htm "RenameFolder 方法 ")

[ToString 方法](../html/83ab8b3b-8124-ea01-71fa-3bdbea2ba010.htm "ToString 方法 ")

[UploadFile 方法](../html/9b1e30fb-8019-f6fd-8891-6d05d1425af8.htm "UploadFile 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FtpClientRenameFile 方法 |

重命名文件的名称  
Rename the name of the file.

**命名空间：**
 [HslCommunication.Enthernet.Ftp](05837d2e-e0a8-4564-6d5e-a834d56ce5ce.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public OperateResult RenameFile(
	string oldName,
	string newName
)
```

```
Public Function RenameFile ( 
	oldName As String,
	newName As String
) As OperateResult
```

```
public:
OperateResult^ RenameFile(
	String^ oldName, 
	String^ newName
)
```

```
member RenameFile : 
        oldName : string * 
        newName : string -> OperateResult 
```

#### 参数

oldName
:   类型：SystemString  
    旧的文件名，例如 /A/123.txt

newName
:   类型：SystemString  
    新文件名，例如 124.txt

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否重命名成功

![](../icons/SectionExpanded.png)参见

#### 引用

[FtpClient 类](98289259-aeed-7095-b875-9d3ef78237fe.htm)

[HslCommunication.Enthernet.Ftp 命名空间](05837d2e-e0a8-4564-6d5e-a834d56ce5ce.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## RenameFolder 方法 

[原文連結](http://api.hslcommunication.cn/html/9348c583-a8a6-89f0-57e1-31ccae07c876.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Enthernet.Ftp](../html/05837d2e-e0a8-4564-6d5e-a834d56ce5ce.htm "HslCommunication.Enthernet.Ftp")

[FtpClient 类](../html/98289259-aeed-7095-b875-9d3ef78237fe.htm "FtpClient 类")

[FtpClient 方法](../html/c26ace03-ff28-7070-8a9e-62b8abf31430.htm "FtpClient 方法")

[CreateDirectory 方法](../html/82ac3dc9-8139-1130-727c-3df3f3d5ce75.htm "CreateDirectory 方法 ")

[DeleteDirectory 方法](../html/ad8fd03c-4875-6c6b-83a3-e7590c7891a3.htm "DeleteDirectory 方法 ")

[DeleteFile 方法](../html/e86754cd-59fd-e54b-be3d-18d3de267d61.htm "DeleteFile 方法 ")

[DownloadFile 方法](../html/41dcfcde-d66a-0549-86f2-0063edee69f3.htm "DownloadFile 方法 ")

[ExtraOnDisconnect 方法](../html/8113d5c7-4f60-fba4-0e43-c0dd4a778d59.htm "ExtraOnDisconnect 方法 ")

[GetLogTextFromBinary 方法](../html/051c0c92-9f31-e96b-26f7-3e1b1a63f272.htm "GetLogTextFromBinary 方法 ")

[GetNewNetMessage 方法](../html/9cb70414-4298-bcf3-fb20-79e4cb0393ad.htm "GetNewNetMessage 方法 ")

[InitializationOnConnect 方法](../html/56cbb726-5571-5c28-4c66-e1820eef5df8.htm "InitializationOnConnect 方法 ")

[ListFiles 方法](../html/c0783d51-47c7-c3d2-bf08-fa9afa5d50c9.htm "ListFiles 方法 ")

[RenameFile 方法](../html/d6aeb935-d84d-2f14-e447-201605bdaa27.htm "RenameFile 方法 ")

[RenameFolder 方法](../html/9348c583-a8a6-89f0-57e1-31ccae07c876.htm "RenameFolder 方法 ")

[ToString 方法](../html/83ab8b3b-8124-ea01-71fa-3bdbea2ba010.htm "ToString 方法 ")

[UploadFile 方法](../html/9b1e30fb-8019-f6fd-8891-6d05d1425af8.htm "UploadFile 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FtpClientRenameFolder 方法 |

重命名目录，不支持跨目录重命名  
Rename a directory; cross-directory renaming is not supported.

**命名空间：**
 [HslCommunication.Enthernet.Ftp](05837d2e-e0a8-4564-6d5e-a834d56ce5ce.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public OperateResult RenameFolder(
	string oldPath,
	string newPath
)
```

```
Public Function RenameFolder ( 
	oldPath As String,
	newPath As String
) As OperateResult
```

```
public:
OperateResult^ RenameFolder(
	String^ oldPath, 
	String^ newPath
)
```

```
member RenameFolder : 
        oldPath : string * 
        newPath : string -> OperateResult 
```

#### 参数

oldPath
:   类型：SystemString  
    例如 /A/B

newPath
:   类型：SystemString  
    例如 C

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否重命名成功

![](../icons/SectionExpanded.png)参见

#### 引用

[FtpClient 类](98289259-aeed-7095-b875-9d3ef78237fe.htm)

[HslCommunication.Enthernet.Ftp 命名空间](05837d2e-e0a8-4564-6d5e-a834d56ce5ce.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ToString 方法 

[原文連結](http://api.hslcommunication.cn/html/83ab8b3b-8124-ea01-71fa-3bdbea2ba010.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Enthernet.Ftp](../html/05837d2e-e0a8-4564-6d5e-a834d56ce5ce.htm "HslCommunication.Enthernet.Ftp")

[FtpClient 类](../html/98289259-aeed-7095-b875-9d3ef78237fe.htm "FtpClient 类")

[FtpClient 方法](../html/c26ace03-ff28-7070-8a9e-62b8abf31430.htm "FtpClient 方法")

[CreateDirectory 方法](../html/82ac3dc9-8139-1130-727c-3df3f3d5ce75.htm "CreateDirectory 方法 ")

[DeleteDirectory 方法](../html/ad8fd03c-4875-6c6b-83a3-e7590c7891a3.htm "DeleteDirectory 方法 ")

[DeleteFile 方法](../html/e86754cd-59fd-e54b-be3d-18d3de267d61.htm "DeleteFile 方法 ")

[DownloadFile 方法](../html/41dcfcde-d66a-0549-86f2-0063edee69f3.htm "DownloadFile 方法 ")

[ExtraOnDisconnect 方法](../html/8113d5c7-4f60-fba4-0e43-c0dd4a778d59.htm "ExtraOnDisconnect 方法 ")

[GetLogTextFromBinary 方法](../html/051c0c92-9f31-e96b-26f7-3e1b1a63f272.htm "GetLogTextFromBinary 方法 ")

[GetNewNetMessage 方法](../html/9cb70414-4298-bcf3-fb20-79e4cb0393ad.htm "GetNewNetMessage 方法 ")

[InitializationOnConnect 方法](../html/56cbb726-5571-5c28-4c66-e1820eef5df8.htm "InitializationOnConnect 方法 ")

[ListFiles 方法](../html/c0783d51-47c7-c3d2-bf08-fa9afa5d50c9.htm "ListFiles 方法 ")

[RenameFile 方法](../html/d6aeb935-d84d-2f14-e447-201605bdaa27.htm "RenameFile 方法 ")

[RenameFolder 方法](../html/9348c583-a8a6-89f0-57e1-31ccae07c876.htm "RenameFolder 方法 ")

[ToString 方法](../html/83ab8b3b-8124-ea01-71fa-3bdbea2ba010.htm "ToString 方法 ")

[UploadFile 方法](../html/9b1e30fb-8019-f6fd-8891-6d05d1425af8.htm "UploadFile 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FtpClientToString 方法 |

[缺少 "M:HslCommunication.Enthernet.Ftp.FtpClient.ToString" 的 <summary> 文档]

**命名空间：**
 [HslCommunication.Enthernet.Ftp](05837d2e-e0a8-4564-6d5e-a834d56ce5ce.htm)  
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

[缺少 "M:HslCommunication.Enthernet.Ftp.FtpClient.ToString" 的 <returns> 文档]

![](../icons/SectionExpanded.png)参见

#### 引用

[FtpClient 类](98289259-aeed-7095-b875-9d3ef78237fe.htm)

[HslCommunication.Enthernet.Ftp 命名空间](05837d2e-e0a8-4564-6d5e-a834d56ce5ce.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## UploadFile 方法 

[原文連結](http://api.hslcommunication.cn/html/9b1e30fb-8019-f6fd-8891-6d05d1425af8.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Enthernet.Ftp](../html/05837d2e-e0a8-4564-6d5e-a834d56ce5ce.htm "HslCommunication.Enthernet.Ftp")

[FtpClient 类](../html/98289259-aeed-7095-b875-9d3ef78237fe.htm "FtpClient 类")

[FtpClient 方法](../html/c26ace03-ff28-7070-8a9e-62b8abf31430.htm "FtpClient 方法")

[CreateDirectory 方法](../html/82ac3dc9-8139-1130-727c-3df3f3d5ce75.htm "CreateDirectory 方法 ")

[DeleteDirectory 方法](../html/ad8fd03c-4875-6c6b-83a3-e7590c7891a3.htm "DeleteDirectory 方法 ")

[DeleteFile 方法](../html/e86754cd-59fd-e54b-be3d-18d3de267d61.htm "DeleteFile 方法 ")

[DownloadFile 方法](../html/41dcfcde-d66a-0549-86f2-0063edee69f3.htm "DownloadFile 方法 ")

[ExtraOnDisconnect 方法](../html/8113d5c7-4f60-fba4-0e43-c0dd4a778d59.htm "ExtraOnDisconnect 方法 ")

[GetLogTextFromBinary 方法](../html/051c0c92-9f31-e96b-26f7-3e1b1a63f272.htm "GetLogTextFromBinary 方法 ")

[GetNewNetMessage 方法](../html/9cb70414-4298-bcf3-fb20-79e4cb0393ad.htm "GetNewNetMessage 方法 ")

[InitializationOnConnect 方法](../html/56cbb726-5571-5c28-4c66-e1820eef5df8.htm "InitializationOnConnect 方法 ")

[ListFiles 方法](../html/c0783d51-47c7-c3d2-bf08-fa9afa5d50c9.htm "ListFiles 方法 ")

[RenameFile 方法](../html/d6aeb935-d84d-2f14-e447-201605bdaa27.htm "RenameFile 方法 ")

[RenameFolder 方法](../html/9348c583-a8a6-89f0-57e1-31ccae07c876.htm "RenameFolder 方法 ")

[ToString 方法](../html/83ab8b3b-8124-ea01-71fa-3bdbea2ba010.htm "ToString 方法 ")

[UploadFile 方法](../html/9b1e30fb-8019-f6fd-8891-6d05d1425af8.htm "UploadFile 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FtpClientUploadFile 方法 |

上传本地文件到FTP服务器，支持进度汇报和取消操作  
Upload local files to an FTP server, with support for progress reporting and cancellation operations.

**命名空间：**
 [HslCommunication.Enthernet.Ftp](05837d2e-e0a8-4564-6d5e-a834d56ce5ce.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public OperateResult UploadFile(
	string localFile,
	string remoteFile,
	Action<long, long> processReport = null,
	HslCancelToken cancelToken = null
)
```

```
Public Function UploadFile ( 
	localFile As String,
	remoteFile As String,
	Optional processReport As Action(Of Long, Long) = Nothing,
	Optional cancelToken As HslCancelToken = Nothing
) As OperateResult
```

```
public:
OperateResult^ UploadFile(
	String^ localFile, 
	String^ remoteFile, 
	Action<long long, long long>^ processReport = nullptr, 
	HslCancelToken^ cancelToken = nullptr
)
```

```
member UploadFile : 
        localFile : string * 
        remoteFile : string * 
        ?processReport : Action<int64, int64> * 
        ?cancelToken : HslCancelToken 
(* Defaults:
        let _processReport = defaultArg processReport null
        let _cancelToken = defaultArg cancelToken null
*)
-> OperateResult 
```

#### 参数

localFile
:   类型：SystemString  
    本地的文件

remoteFile
:   类型：SystemString  
    服务器保存的路径，例如 /A/123.txt

processReport (Optional)
:   类型：SystemActionInt64, Int64  
    上传的进度

cancelToken (Optional)
:   类型：[HslCommunication.CoreHslCancelToken](072498cc-b288-38cf-c2d2-7d9452ac6a3e.htm)  
    取消操作对象

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否上传成功

![](../icons/SectionExpanded.png)参见

#### 引用

[FtpClient 类](98289259-aeed-7095-b875-9d3ef78237fe.htm)

[HslCommunication.Enthernet.Ftp 命名空间](05837d2e-e0a8-4564-6d5e-a834d56ce5ce.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## FtpClient 字段

[原文連結](http://api.hslcommunication.cn/html/2f2894ac-2a80-2e5c-3af6-6258dfebafdf.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Enthernet.Ftp](../html/05837d2e-e0a8-4564-6d5e-a834d56ce5ce.htm "HslCommunication.Enthernet.Ftp")

[FtpClient 类](../html/98289259-aeed-7095-b875-9d3ef78237fe.htm "FtpClient 类")

[FtpClient 构造函数](../html/034f536b-56cd-7b2b-7f55-c493aa15c924.htm "FtpClient 构造函数 ")

[FtpClient 属性](../html/69a325b9-c5bb-9bf7-f4b4-847f840c54d4.htm "FtpClient 属性")

[FtpClient 方法](../html/c26ace03-ff28-7070-8a9e-62b8abf31430.htm "FtpClient 方法")

[FtpClient 字段](../html/2f2894ac-2a80-2e5c-3af6-6258dfebafdf.htm "FtpClient 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FtpClient 字段 |

[FtpClient](98289259-aeed-7095-b875-9d3ef78237fe.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)字段

|  | 名称 | 说明 |
| --- | --- | --- |
| 受保护的字段 | [LogMsgFormatBinary](2ed2ff87-caec-0e44-f08b-e2b015225a58.htm) | 设置日志记录报文是否二进制，如果为False，那就使用ASCII码  Set whether the log message is binary, if it is False, then use ASCII code (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[FtpClient 类](98289259-aeed-7095-b875-9d3ef78237fe.htm)

[HslCommunication.Enthernet.Ftp 命名空间](05837d2e-e0a8-4564-6d5e-a834d56ce5ce.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## FtpFileItem 类

[原文連結](http://api.hslcommunication.cn/html/97ec5ae6-c5a1-0122-982b-cb02f1031c96.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Enthernet.Ftp](../html/05837d2e-e0a8-4564-6d5e-a834d56ce5ce.htm "HslCommunication.Enthernet.Ftp")

[FtpFileItem 类](../html/97ec5ae6-c5a1-0122-982b-cb02f1031c96.htm "FtpFileItem 类")

[FtpFileItem 构造函数](../html/d2e60541-69f6-ac12-c2c1-0c5d1d5aab00.htm "FtpFileItem 构造函数 ")

[FtpFileItem 属性](../html/964711d1-f8ed-a250-7df4-b5dc69f09165.htm "FtpFileItem 属性")

[FtpFileItem 方法](../html/ba27823a-f9f5-c1ca-fb4b-b1b824cdfd1e.htm "FtpFileItem 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FtpFileItem 类 |

FTP文件的信息对象

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  HslCommunication.Enthernet.FtpFtpFileItem

**命名空间：**
 [HslCommunication.Enthernet.Ftp](05837d2e-e0a8-4564-6d5e-a834d56ce5ce.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class FtpFileItem
```

```
Public Class FtpFileItem
```

```
public ref class FtpFileItem
```

```
type FtpFileItem =  class end
```

FtpFileItem 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [FtpFileItem](01290781-4927-4ebe-f013-b46b660e8134.htm) | 实例化一个默认的对象 |
| 公共方法 | [FtpFileItem(String)](1a3dfa69-1638-d5fb-cf9d-6cc1598980e2.htm) | 指定行内容解析出对象信息 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [CreateTime](c49b1168-b67b-2a2a-b6bd-1cecb1cb65b7.htm) | 文件的创建时间 |
| 公共属性 | [IsDirectory](54b4edce-7441-29e1-70e8-08f50250a854.htm) | 是否文件夹的标识 |
| 公共属性 | [Name](f3bd5e40-549c-a932-0c5e-658698759680.htm) | 文件名或是文件夹名 |
| 公共属性 | [Size](ad3e9a2f-8c66-7d69-f5d3-6b44606830ad.htm) | 文件的大小，字节为单位 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法 | ToString | (继承自 Object。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.Enthernet.Ftp 命名空间](05837d2e-e0a8-4564-6d5e-a834d56ce5ce.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## FtpFileItem 构造函数 

[原文連結](http://api.hslcommunication.cn/html/d2e60541-69f6-ac12-c2c1-0c5d1d5aab00.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Enthernet.Ftp](../html/05837d2e-e0a8-4564-6d5e-a834d56ce5ce.htm "HslCommunication.Enthernet.Ftp")

[FtpFileItem 类](../html/97ec5ae6-c5a1-0122-982b-cb02f1031c96.htm "FtpFileItem 类")

[FtpFileItem 构造函数](../html/d2e60541-69f6-ac12-c2c1-0c5d1d5aab00.htm "FtpFileItem 构造函数 ")

[FtpFileItem 构造函数](../html/01290781-4927-4ebe-f013-b46b660e8134.htm "FtpFileItem 构造函数 ")

[FtpFileItem 构造函数 (String)](../html/1a3dfa69-1638-d5fb-cf9d-6cc1598980e2.htm "FtpFileItem 构造函数 (String)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FtpFileItem 构造函数 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [FtpFileItem](01290781-4927-4ebe-f013-b46b660e8134.htm) | 实例化一个默认的对象 |
| 公共方法 | [FtpFileItem(String)](1a3dfa69-1638-d5fb-cf9d-6cc1598980e2.htm) | 指定行内容解析出对象信息 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[FtpFileItem 类](97ec5ae6-c5a1-0122-982b-cb02f1031c96.htm)

[HslCommunication.Enthernet.Ftp 命名空间](05837d2e-e0a8-4564-6d5e-a834d56ce5ce.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## FtpFileItem 构造函数 

[原文連結](http://api.hslcommunication.cn/html/01290781-4927-4ebe-f013-b46b660e8134.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Enthernet.Ftp](../html/05837d2e-e0a8-4564-6d5e-a834d56ce5ce.htm "HslCommunication.Enthernet.Ftp")

[FtpFileItem 类](../html/97ec5ae6-c5a1-0122-982b-cb02f1031c96.htm "FtpFileItem 类")

[FtpFileItem 构造函数](../html/d2e60541-69f6-ac12-c2c1-0c5d1d5aab00.htm "FtpFileItem 构造函数 ")

[FtpFileItem 构造函数](../html/01290781-4927-4ebe-f013-b46b660e8134.htm "FtpFileItem 构造函数 ")

[FtpFileItem 构造函数 (String)](../html/1a3dfa69-1638-d5fb-cf9d-6cc1598980e2.htm "FtpFileItem 构造函数 (String)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FtpFileItem 构造函数 |

实例化一个默认的对象

**命名空间：**
 [HslCommunication.Enthernet.Ftp](05837d2e-e0a8-4564-6d5e-a834d56ce5ce.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public FtpFileItem()
```

```
Public Sub New
```

```
public:
FtpFileItem()
```

```
new : unit -> FtpFileItem
```

![](../icons/SectionExpanded.png)参见

#### 引用

[FtpFileItem 类](97ec5ae6-c5a1-0122-982b-cb02f1031c96.htm)

[FtpFileItem 重载](d2e60541-69f6-ac12-c2c1-0c5d1d5aab00.htm)

[HslCommunication.Enthernet.Ftp 命名空间](05837d2e-e0a8-4564-6d5e-a834d56ce5ce.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## FtpFileItem 构造函数 (String)

[原文連結](http://api.hslcommunication.cn/html/1a3dfa69-1638-d5fb-cf9d-6cc1598980e2.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Enthernet.Ftp](../html/05837d2e-e0a8-4564-6d5e-a834d56ce5ce.htm "HslCommunication.Enthernet.Ftp")

[FtpFileItem 类](../html/97ec5ae6-c5a1-0122-982b-cb02f1031c96.htm "FtpFileItem 类")

[FtpFileItem 构造函数](../html/d2e60541-69f6-ac12-c2c1-0c5d1d5aab00.htm "FtpFileItem 构造函数 ")

[FtpFileItem 构造函数](../html/01290781-4927-4ebe-f013-b46b660e8134.htm "FtpFileItem 构造函数 ")

[FtpFileItem 构造函数 (String)](../html/1a3dfa69-1638-d5fb-cf9d-6cc1598980e2.htm "FtpFileItem 构造函数 (String)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FtpFileItem 构造函数 (String) |

指定行内容解析出对象信息

**命名空间：**
 [HslCommunication.Enthernet.Ftp](05837d2e-e0a8-4564-6d5e-a834d56ce5ce.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public FtpFileItem(
	string content
)
```

```
Public Sub New ( 
	content As String
)
```

```
public:
FtpFileItem(
	String^ content
)
```

```
new : 
        content : string -> FtpFileItem
```

#### 参数

content
:   类型：SystemString  
    原始的数据信息

![](../icons/SectionExpanded.png)参见

#### 引用

[FtpFileItem 类](97ec5ae6-c5a1-0122-982b-cb02f1031c96.htm)

[FtpFileItem 重载](d2e60541-69f6-ac12-c2c1-0c5d1d5aab00.htm)

[HslCommunication.Enthernet.Ftp 命名空间](05837d2e-e0a8-4564-6d5e-a834d56ce5ce.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## FtpFileItem 属性

[原文連結](http://api.hslcommunication.cn/html/964711d1-f8ed-a250-7df4-b5dc69f09165.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Enthernet.Ftp](../html/05837d2e-e0a8-4564-6d5e-a834d56ce5ce.htm "HslCommunication.Enthernet.Ftp")

[FtpFileItem 类](../html/97ec5ae6-c5a1-0122-982b-cb02f1031c96.htm "FtpFileItem 类")

[FtpFileItem 属性](../html/964711d1-f8ed-a250-7df4-b5dc69f09165.htm "FtpFileItem 属性")

[CreateTime 属性](../html/c49b1168-b67b-2a2a-b6bd-1cecb1cb65b7.htm "CreateTime 属性 ")

[IsDirectory 属性](../html/54b4edce-7441-29e1-70e8-08f50250a854.htm "IsDirectory 属性 ")

[Name 属性](../html/f3bd5e40-549c-a932-0c5e-658698759680.htm "Name 属性 ")

[Size 属性](../html/ad3e9a2f-8c66-7d69-f5d3-6b44606830ad.htm "Size 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FtpFileItem 属性 |

[FtpFileItem](97ec5ae6-c5a1-0122-982b-cb02f1031c96.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [CreateTime](c49b1168-b67b-2a2a-b6bd-1cecb1cb65b7.htm) | 文件的创建时间 |
| 公共属性 | [IsDirectory](54b4edce-7441-29e1-70e8-08f50250a854.htm) | 是否文件夹的标识 |
| 公共属性 | [Name](f3bd5e40-549c-a932-0c5e-658698759680.htm) | 文件名或是文件夹名 |
| 公共属性 | [Size](ad3e9a2f-8c66-7d69-f5d3-6b44606830ad.htm) | 文件的大小，字节为单位 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[FtpFileItem 类](97ec5ae6-c5a1-0122-982b-cb02f1031c96.htm)

[HslCommunication.Enthernet.Ftp 命名空间](05837d2e-e0a8-4564-6d5e-a834d56ce5ce.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## CreateTime 属性 

[原文連結](http://api.hslcommunication.cn/html/c49b1168-b67b-2a2a-b6bd-1cecb1cb65b7.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Enthernet.Ftp](../html/05837d2e-e0a8-4564-6d5e-a834d56ce5ce.htm "HslCommunication.Enthernet.Ftp")

[FtpFileItem 类](../html/97ec5ae6-c5a1-0122-982b-cb02f1031c96.htm "FtpFileItem 类")

[FtpFileItem 属性](../html/964711d1-f8ed-a250-7df4-b5dc69f09165.htm "FtpFileItem 属性")

[CreateTime 属性](../html/c49b1168-b67b-2a2a-b6bd-1cecb1cb65b7.htm "CreateTime 属性 ")

[IsDirectory 属性](../html/54b4edce-7441-29e1-70e8-08f50250a854.htm "IsDirectory 属性 ")

[Name 属性](../html/f3bd5e40-549c-a932-0c5e-658698759680.htm "Name 属性 ")

[Size 属性](../html/ad3e9a2f-8c66-7d69-f5d3-6b44606830ad.htm "Size 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FtpFileItemCreateTime 属性 |

文件的创建时间

**命名空间：**
 [HslCommunication.Enthernet.Ftp](05837d2e-e0a8-4564-6d5e-a834d56ce5ce.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public DateTime CreateTime { get; set; }
```

```
Public Property CreateTime As DateTime
	Get
	Set
```

```
public:
property DateTime CreateTime {
	DateTime get ();
	void set (DateTime value);
}
```

```
member CreateTime : DateTime with get, set
```

#### 属性值

类型：DateTime

![](../icons/SectionExpanded.png)参见

#### 引用

[FtpFileItem 类](97ec5ae6-c5a1-0122-982b-cb02f1031c96.htm)

[HslCommunication.Enthernet.Ftp 命名空间](05837d2e-e0a8-4564-6d5e-a834d56ce5ce.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## IsDirectory 属性 

[原文連結](http://api.hslcommunication.cn/html/54b4edce-7441-29e1-70e8-08f50250a854.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Enthernet.Ftp](../html/05837d2e-e0a8-4564-6d5e-a834d56ce5ce.htm "HslCommunication.Enthernet.Ftp")

[FtpFileItem 类](../html/97ec5ae6-c5a1-0122-982b-cb02f1031c96.htm "FtpFileItem 类")

[FtpFileItem 属性](../html/964711d1-f8ed-a250-7df4-b5dc69f09165.htm "FtpFileItem 属性")

[CreateTime 属性](../html/c49b1168-b67b-2a2a-b6bd-1cecb1cb65b7.htm "CreateTime 属性 ")

[IsDirectory 属性](../html/54b4edce-7441-29e1-70e8-08f50250a854.htm "IsDirectory 属性 ")

[Name 属性](../html/f3bd5e40-549c-a932-0c5e-658698759680.htm "Name 属性 ")

[Size 属性](../html/ad3e9a2f-8c66-7d69-f5d3-6b44606830ad.htm "Size 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FtpFileItemIsDirectory 属性 |

是否文件夹的标识

**命名空间：**
 [HslCommunication.Enthernet.Ftp](05837d2e-e0a8-4564-6d5e-a834d56ce5ce.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public bool IsDirectory { get; set; }
```

```
Public Property IsDirectory As Boolean
	Get
	Set
```

```
public:
property bool IsDirectory {
	bool get ();
	void set (bool value);
}
```

```
member IsDirectory : bool with get, set
```

#### 属性值

类型：Boolean

![](../icons/SectionExpanded.png)参见

#### 引用

[FtpFileItem 类](97ec5ae6-c5a1-0122-982b-cb02f1031c96.htm)

[HslCommunication.Enthernet.Ftp 命名空间](05837d2e-e0a8-4564-6d5e-a834d56ce5ce.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Name 属性 

[原文連結](http://api.hslcommunication.cn/html/f3bd5e40-549c-a932-0c5e-658698759680.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Enthernet.Ftp](../html/05837d2e-e0a8-4564-6d5e-a834d56ce5ce.htm "HslCommunication.Enthernet.Ftp")

[FtpFileItem 类](../html/97ec5ae6-c5a1-0122-982b-cb02f1031c96.htm "FtpFileItem 类")

[FtpFileItem 属性](../html/964711d1-f8ed-a250-7df4-b5dc69f09165.htm "FtpFileItem 属性")

[CreateTime 属性](../html/c49b1168-b67b-2a2a-b6bd-1cecb1cb65b7.htm "CreateTime 属性 ")

[IsDirectory 属性](../html/54b4edce-7441-29e1-70e8-08f50250a854.htm "IsDirectory 属性 ")

[Name 属性](../html/f3bd5e40-549c-a932-0c5e-658698759680.htm "Name 属性 ")

[Size 属性](../html/ad3e9a2f-8c66-7d69-f5d3-6b44606830ad.htm "Size 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FtpFileItemName 属性 |

文件名或是文件夹名

**命名空间：**
 [HslCommunication.Enthernet.Ftp](05837d2e-e0a8-4564-6d5e-a834d56ce5ce.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public string Name { get; set; }
```

```
Public Property Name As String
	Get
	Set
```

```
public:
property String^ Name {
	String^ get ();
	void set (String^ value);
}
```

```
member Name : string with get, set
```

#### 属性值

类型：String

![](../icons/SectionExpanded.png)参见

#### 引用

[FtpFileItem 类](97ec5ae6-c5a1-0122-982b-cb02f1031c96.htm)

[HslCommunication.Enthernet.Ftp 命名空间](05837d2e-e0a8-4564-6d5e-a834d56ce5ce.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)