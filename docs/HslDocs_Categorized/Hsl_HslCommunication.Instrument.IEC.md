# HslCommunication - HslCommunication.Instrument.IEC

> 分類頁數: 30



---
## HslCommunication.Instrument.IEC

[原文連結](http://api.hslcommunication.cn/html/4f75643b-b529-650f-8587-78378826eec2.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.IEC](../html/4f75643b-b529-650f-8587-78378826eec2.htm "HslCommunication.Instrument.IEC")

[IEC104 类](../html/8695140f-80b4-71fc-7a09-eff994df68ca.htm "IEC104 类")

[IEC104.TypeID 类](../html/634a3bbf-dbf7-b538-6361-94ebfde09a1c.htm "IEC104.TypeID 类")

[IEC104MessageEventArgs 类](../html/34f93396-e777-2b46-1db2-20c3f8ec3d23.htm "IEC104MessageEventArgs 类")

[IEC104Server 类](../html/f956b272-09e9-282f-5756-98f2249f61e2.htm "IEC104Server 类")

[IecValueObject(T) 类](../html/dd68297c-f009-b216-89a4-112d4eecb574.htm "IecValueObject(T) 类")

[IecValueServerObject(T) 类](../html/10ccf1c2-45a3-9052-4ad9-d10b82269f20.htm "IecValueServerObject(T) 类")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslCommunication.Instrument.IEC 命名空间 |

[缺少 "N:HslCommunication.Instrument.IEC" 的 <summary> 文档]

![](../icons/SectionExpanded.png)类

|  | 类 | 说明 |
| --- | --- | --- |
| 公共类代码示例 | [IEC104](8695140f-80b4-71fc-7a09-eff994df68ca.htm) | IEC104规约实现的电力协议，支持总召唤，支持回调任意的突发上传事件 |
| 公共类 | [IEC104TypeID](634a3bbf-dbf7-b538-6361-94ebfde09a1c.htm) | 类型信息资源 |
| 公共类 | [IEC104MessageEventArgs](34f93396-e777-2b46-1db2-20c3f8ec3d23.htm) | IEC104的消息事件 |
| 公共类 | [IEC104Server](f956b272-09e9-282f-5756-98f2249f61e2.htm) | IEC104的虚拟服务器，支持总召唤操作，数据区包含单点遥信，双点遥信，归一化遥测值，标度化遥测值，短浮点遥测值，32位比特串，默认地址数量100，也可以手动设置更大的地址范围  The virtual server of IEC104 supports total call operation. The data area includes single-point remote signaling, dual-point remote signaling, normalized telemetry values, scaled telemetry values, short floating-point telemetry values, 32-bit bit strings, and the default address range is 100. A larger address range can also be manually set |
| 公共类 | [IecValueObjectT](dd68297c-f009-b216-89a4-112d4eecb574.htm) | IEC的数据对象，带值，品质信息，地址信息，时标信息 |
| 公共类 | [IecValueServerObjectT](10ccf1c2-45a3-9052-4ad9-d10b82269f20.htm) | 服务器端的IEC数据对象信息，可以实例化不连续地址，或是连续地址信息 |

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## IEC104 类

[原文連結](http://api.hslcommunication.cn/html/8695140f-80b4-71fc-7a09-eff994df68ca.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.IEC](../html/4f75643b-b529-650f-8587-78378826eec2.htm "HslCommunication.Instrument.IEC")

[IEC104 类](../html/8695140f-80b4-71fc-7a09-eff994df68ca.htm "IEC104 类")

[IEC104 构造函数](../html/22a80517-47fd-8960-d826-72c0a07e021f.htm "IEC104 构造函数 ")

[IEC104 属性](../html/1da629b7-909c-dfb0-74a8-3a845a94fe2a.htm "IEC104 属性")

[IEC104 方法](../html/8eb8a1e0-92fb-2fdb-83ac-2cab7f59d96b.htm "IEC104 方法")

[IEC104 事件](../html/98e01792-5618-ac77-38c5-20d50d856869.htm "IEC104 事件")

[IEC104 字段](../html/f6e20d95-3950-c9c9-a87e-21f561d422d0.htm "IEC104 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IEC104 类 |

IEC104规约实现的电力协议，支持总召唤，支持回调任意的突发上传事件

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  [HslCommunication.Core.NetBinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)  
    [HslCommunication.Core.DeviceDeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)  
      [HslCommunication.Core.DeviceDeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)  
        HslCommunication.Instrument.IECIEC104

**命名空间：**
 [HslCommunication.Instrument.IEC](4f75643b-b529-650f-8587-78378826eec2.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class IEC104 : DeviceTcpNet
```

```
Public Class IEC104
	Inherits DeviceTcpNet
```

```
public ref class IEC104 : public DeviceTcpNet
```

```
type IEC104 =  
    class
        inherit DeviceTcpNet
    end
```

IEC104 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [IEC104](4d456238-2b09-50c6-42c0-cd6893cb5c7a.htm) | 实例化IEC104协议的通讯对象  Instantiate the communication object of the IEC104 protocol |
| 公共方法 | [IEC104(String, Int32)](125f8e8d-2907-db1f-c4b1-7b76e85f102d.htm) | 指定ip地址和端口号来实例化一个默认的对象  Specify the IP address and port number to instantiate a default object |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性代码示例 | [ByteTransform](f35118c9-3691-f752-ae6f-2c5549714faa.htm) | 当前的数据变换机制，当你需要从字节数据转换类型数据的时候需要。  The current data transformation mechanism is required when you need to convert type data from byte data. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共属性 | [CommunicationPipe](ad8e903b-d289-21b0-077e-15ede8005221.htm) | 获取或设置当前的管道信息，管道类型为[CommunicationPipe](ad8e903b-d289-21b0-077e-15ede8005221.htm)的继承类，内置了[PipeTcpNet](26b94bd5-2695-36e7-6378-c3d41846be8f.htm)管道，[PipeUdpNet](5a0b810e-187e-826b-0a13-34fd5576aa0e.htm)管道，[PipeSerialPort](9d924607-3dfc-3916-b532-c5dec48fe347.htm)管道等  Get or set the current pipeline information, the pipeline type is [CommunicationPipe](ad8e903b-d289-21b0-077e-15ede8005221.htm) inheritance class, [PipeTcpNet](26b94bd5-2695-36e7-6378-c3d41846be8f.htm) pipeline, [PipeUdpNet](5a0b810e-187e-826b-0a13-34fd5576aa0e.htm) pipeline, [PipeSerialPort](9d924607-3dfc-3916-b532-c5dec48fe347.htm) pipeline, etc (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [ConnectionId](ce1efb4b-fbc9-4683-163e-0a8e71e99d6b.htm) | 当前连接的唯一ID号，默认为长度20的guid码加随机数组成，方便列表管理，也可以自己指定  The unique ID number of the current connection. The default is a 20-digit guid code plus a random number. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性代码示例 | [ConnectTimeOut](71a651e5-737d-6e88-75a2-fb891b9af6cd.htm) | 获取或设置接收服务器反馈的时间，如果为负数，则不接收反馈   Gets or sets the time to receive server feedback, and if it is a negative number, does not receive feedback (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共属性代码示例 | [IpAddress](c378f4e8-3fdf-1932-7239-ed44bc5af584.htm) | 获取或是设置远程服务器的IP地址，如果是本机测试，那么需要设置为127.0.0.1   Get or set the IP address of the remote server. If it is a local test, then it needs to be set to 127.0.0.1 (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共属性 | [LocalBinding](622b418f-5969-3726-5dcc-04456c4b97a7.htm) | 获取或设置绑定的本地的IP地址和端口号信息，如果端口设置为0，代表任何可用的端口  Get or set the bound local IP address and port number information, if the port is set to 0, it means any available port (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共属性代码示例 | [LogNet](8c56f246-88ac-afb8-8c06-972fafe4f78a.htm) | 组件的日志工具，支持日志记录，只要实例化后，当前网络的基本信息，就以[DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm)等级进行输出  The component's logging tool supports logging. As long as the instantiation of the basic network information, the output will be output at [DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm) (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性代码示例 | [Port](6f1ad24a-c916-5b67-aa17-a8c7dc2acb2e.htm) | 获取或设置服务器的端口号，具体的值需要取决于对方的配置  Gets or sets the port number of the server. The specific value depends on the configuration of the other party. (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共属性代码示例 | [ReceiveTimeOut](250560a7-332b-ce1b-5eb0-e3abe4ee6012.htm) | 获取或设置接收服务器反馈的时间，如果为负数，则不接收反馈   Gets or sets the time to receive server feedback, and if it is a negative number, does not receive feedback (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [SendBeforeHex](72aa1de4-116f-abfe-4adb-814a0be0e560.htm) | 获取或设置在发送通信报文前追加发送的字节信息，HEX格式，通常用于lora组网时，需要携带 00 00 00 02 四个字节的站地址功能。  Obtain or set the byte information sent before sending communication packets, HEX format, usually used for LORA networking, you need to carry 00 00 00 02 four-byte station address function. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [SleepTime](86a5fbb9-3456-6beb-7113-5bbad5eabbc5.htm) | 获取或设置在正式接收对方返回数据前的时候，需要休息的时间，当设置为0的时候，不需要休息。  Get or set the time required to rest before officially receiving the data from the other party. When it is set to 0, no rest is required. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [SocketKeepAliveTime](96ab7703-8c4d-b4c2-ddd9-0dc93ebcfc8d.htm) | 获取或设置客户端的Socket的心跳时间信息，这个是Socket底层自动实现的心跳包，不基于协议层实现。默认小于0，不开启心跳检测，如果需要开启，设置 60\_000 比较合适，单位毫秒  Get or set the heartbeat time information of the Socket of the client. This is the heartbeat packet automatically implemented by the bottom layer of the Socket, not based on the protocol layer. The default value is less than 0, and heartbeat detection is not enabled. If you need to enable it, it is more appropriate to set 60\_000, in milliseconds. (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共属性 | [Station](9fddc8b9-c09a-2499-b762-17943593b07a.htm) | 公共单元地址，在低版本远动程序中，站地址范围一般在1~254，255表示全局地址。新远动程序则可支持1~65534为站地址，而65535为全局地址。默认值为1 |
| 受保护的属性 | [WordLength](f1a505cc-e85c-6c8b-b189-58a1ab5753f9.htm) | 一个字单位的数据表示的地址长度，西门子为2，三菱，欧姆龙，modbusTcp就为1，AB PLC无效  The address length represented by one word of data, Siemens is 2, Mitsubishi, Omron, modbusTcp is 1, AB PLC is invalid (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法代码示例 | [ConnectClose](569a9c0e-7628-b8d5-e124-0853452747b8.htm) | 手动断开与远程服务器的连接，如果当前是长连接模式，那么就会切换到短连接模式  Manually disconnect from the remote server, if it is currently in long connection mode, it will switch to short connection mode (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共方法代码示例 | [ConnectCloseAsync](76be8150-ff22-95b9-16b1-7c3440c54ae8.htm) | 手动断开与远程服务器的连接，如果当前是长连接模式，那么就会切换到短连接模式  Manually disconnect from the remote server, if it is currently in long connection mode, it will switch to short connection mode (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共方法代码示例 | [ConnectServer](6299cc1f-2fab-b754-4597-1eddc2b21308.htm) | 尝试连接远程的服务器，如果连接成功，就切换短连接模式到长连接模式，后面的每次请求都共享一个通道，使得通讯速度更快速  Try to connect to a remote server. If the connection is successful, switch the short connection mode to the long connection mode. Each subsequent request will share a channel, making the communication speed faster. (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共方法代码示例 | [ConnectServerAsync](e7eff5c0-8e2a-3b6c-5027-16f4d8bcf5ae.htm) | 尝试连接远程的服务器，如果连接成功，就切换短连接模式到长连接模式，后面的每次请求都共享一个通道，使得通讯速度更快速  Try to connect to a remote server. If the connection is successful, switch the short connection mode to the long connection mode. Each subsequent request will share a channel, making the communication speed faster. (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 受保护的方法 | [DecideWhetherQAMessage](a1d1b211-1f5a-e398-8d60-90c959313942.htm) | 决定当前的消息是否是用于问答机制返回的消息，默认直接返回 true, 实际的情况需要根据协议进行重写方法  To determine whether the current message is the message returned by the question answering mechanism, the default is true. In actual cases, the rewriting method needs to be performed according to the protocol (重写 [BinaryCommunicationDecideWhetherQAMessage(CommunicationPipe, OperateResultByte)](acc40cda-e310-8537-6685-59b67eeb16ff.htm).) |
| 公共方法 | [Dispose](b9dda6bf-d342-254e-689a-93fa0bb265d1.htm) | (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 受保护的方法 | [Dispose(Boolean)](0bbb6ed3-111e-069c-e9f5-8cbc3f7f7b77.htm) | (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | [ExtraAfterReadFromCoreServer](06a49a15-a3ce-5e39-28c7-2b78851eca2e.htm) | 和服务器交互完成的时候调用的方法，可以根据读写结果进行一些额外的操作，具体的操作需要根据实际的需求来重写实现  The method called when the interaction with the server is completed can perform some additional operations based on the read and write results. The specific operations need to be rewritten according to actual needs. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法代码示例 | [ExtraOnDisconnect](de9dbba2-3b2c-1c5e-53c8-4e6f8a44e2c4.htm) | 根据实际的协议选择是否重写本方法，有些协议在断开连接之前，需要发送一些报文来关闭当前的网络通道  Select whether to rewrite this method according to the actual protocol. Some protocols need to send some packets to close the current network channel before disconnecting. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法代码示例 | [ExtraOnDisconnectAsync](210f593e-16d8-77c4-1178-58fbcfe9b27b.htm) | 根据实际的协议选择是否重写本方法，有些协议在断开连接之前，需要发送一些报文来关闭当前的网络通道  Select whether to rewrite this method according to the actual protocol. Some protocols need to send some packets to close the current network channel before disconnecting. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 受保护的方法 | [GetLogTextFromBinary](67f3b2f7-957c-fa0b-8320-6799237157db.htm) | 获取当前的报文进行日志记录的时候，是否使用二进制的格式记录，默认返回 [LogMsgFormatBinary](2ed2ff87-caec-0e44-f08b-e2b015225a58.htm)，重写可以根据session对象分别返回不同记录模式  Whether to log the current packet in binary format, the default return is [LogMsgFormatBinary](2ed2ff87-caec-0e44-f08b-e2b015225a58.htm). If you want to override it, different recording modes can be returned according to session (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [GetNewNetMessage](5272711a-7847-2bfd-ef08-bf54632cbef2.htm) | 获取一个新的消息对象的方法，需要在继承类里面进行重写  The method to get a new message object needs to be overridden in the inheritance class (重写 [BinaryCommunicationGetNewNetMessage](16636a1b-a6f0-9c0f-9d9c-c94e26c67459.htm).) |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法 | [GetWordLength](c27821ea-1739-8eb9-a848-23a00828fec8.htm) | 一个字单位的数据表示的地址长度，西门子为2，三菱，欧姆龙，modbusTcp就为1，AB PLC无效  The address length represented by one word of data, Siemens is 2, Mitsubishi, Omron, modbusTcp is 1, AB PLC is invalid (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 受保护的方法代码示例 | [InitializationOnConnect](3acd2dc3-e3a2-b827-23ec-08a3fb8a6d6d.htm) | 根据实际的协议选择是否重写本方法，有些协议在创建连接之后，需要进行一些初始化的信号握手，才能最终建立网络通道。  Whether to rewrite this method is based on the actual protocol. Some protocols require some initial signal handshake to establish a network channel after the connection is created. (重写 [BinaryCommunicationInitializationOnConnect](4ab6989f-a5b5-05fe-8b44-bfa45824a2a7.htm).) |
| 受保护的方法代码示例 | [InitializationOnConnectAsync](deaf517e-2c31-d73d-562f-371bdeb4d6c2.htm) | 根据实际的协议选择是否重写本方法，有些协议在创建连接之后，需要进行一些初始化的信号握手，才能最终建立网络通道。  Whether to rewrite this method is based on the actual protocol. Some protocols require some initial signal handshake to establish a network channel after the connection is created. (重写 [BinaryCommunicationInitializationOnConnectAsync](7da2af04-d4fc-5045-19a4-4b0d1722eef3.htm).) |
| 公共方法 | [IpAddressPing](0ac126b4-a830-a822-4ef5-855474942cd7.htm) | 对当前设备的IP地址进行PING的操作，返回PING的结果，正常来说，返回Success  PING the IP address of the current device and return the PING result. Normally, it returns Success (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 受保护的方法 | [LogRevcMessage(Byte)](dab1657f-1e75-7fe1-7003-ccedfc69855d.htm) | 使用日志记录一个接收的报文信息  Logs are used to record information about a received packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [LogRevcMessage(Byte, PipeSession)](0f51ae45-575e-4228-f978-01b37d6d5df1.htm) | 使用日志记录一个接收的报文信息  Logs are used to record information about a received packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [LogSendMessage(Byte)](4f90eecc-68a2-1ef2-44ef-369f19bed3ec.htm) | 使用日志记录一个发送的报文信息  Logs are used to record information about a send packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [LogSendMessage(Byte, PipeSession)](209caadf-4f0f-40aa-04c2-ebc61f8dab13.htm) | 使用日志记录一个发送的报文信息  Logs are used to record information about a send packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法 | [PackCommandWithHeader](d0b2cd76-c306-11fc-38c7-af8ad306cce7.htm) | 对当前的命令进行打包处理，通常是携带命令头内容，标记当前的命令的长度信息，需要进行重写，否则默认不打包  The current command is packaged, usually carrying the content of the command header, marking the length of the current command, and it needs to be rewritten, otherwise it is not packaged by default (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [Read(String, UInt16)](9dd09452-c200-bc65-75ff-69aa34b366d4.htm) | 批量读取字节数组信息，需要指定地址和长度，返回原始的字节数组  Batch read byte array information, need to specify the address and length, return the original byte array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadT](95abe064-5594-5ae7-a849-98c53f88fee0.htm) | 读取支持Hsl特性的数据内容，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明，或是官网的详细文档  Read the data content of the Hsl attribute. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadAsync(String, UInt16)](d8c14d37-56dc-2a31-1a19-3e2b5b65197d.htm) | 异步批量读取字节数组信息，需要指定地址和长度，返回原始的字节数组  Asynchronous batch read byte array information, need to specify the address and length, return the original byte array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadAsyncT](0522e740-4bf0-2e31-c4c5-e41384321b37.htm) | 异步读取支持Hsl特性的数据内容，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Asynchronously read the data content of the Hsl attribute. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadBool(String)](9b2faa2e-eaa9-7ace-c49f-7a0318f1df1a.htm) | 读取单个的Boolean数据信息  Read a single Boolean data message (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadBool(String, UInt16)](69d0d24f-34a7-5b41-5c1b-619edc7e6b54.htm) | 批量读取Boolean数组信息，需要指定地址和长度，返回Boolean 数组  Batch read Boolean array information, need to specify the address and length, return Boolean array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadBoolAsync(String)](cc37da1f-d673-fe67-78ee-4e465d29a2bd.htm) | 异步读取单个的Boolean数据信息  Asynchronously read a single Boolean data message (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadBoolAsync(String, UInt16)](a45d383d-0095-3ea7-32f8-e75af88d60f8.htm) | 异步批量读取Boolean数组信息，需要指定地址和长度，返回Boolean 数组  Asynchronously batch read Boolean array information, need to specify the address and length, return Boolean array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadCustomerT(String)](ca069688-d8f0-7dbd-e962-86a10dce83cf.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，返回一个新的类型的实例对象。  To read a custom data type, you need to inherit from the IDataTransfer interface and return an instance object of a new type. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadCustomerT(String, T)](ca7dc1a1-c97d-689f-9819-4dae4af67206.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，传入一个实例，对这个实例进行赋值，并返回该实例的对象。  To read a custom data type, you need to inherit from the IDataTransfer interface, pass in an instance, assign a value to this instance, and return the object of the instance. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadCustomerAsyncT(String)](f3fcda2b-255d-aadd-567e-c7b3df3636d9.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，返回一个新的类型的实例对象。  To read a custom data type, you need to inherit from the IDataTransfer interface and return an instance object of a new type. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadCustomerAsyncT(String, T)](527c9fae-4743-6ebd-7b9a-5ef3a8b7c455.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，传入一个实例，对这个实例进行赋值，并返回该实例的对象。  To read a custom data type, you need to inherit from the IDataTransfer interface, pass in an instance, assign a value to this instance, and return the object of the instance. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadDouble(String)](e7e31c43-71b4-e692-6bdc-ba9533e1d83d.htm) | 读取双浮点的数据  Read double floating point data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadDouble(String, UInt16)](4eee0f05-f861-4b48-3f06-99737fb24aee.htm) | 读取双浮点数据的数组  Read double floating point data array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadDoubleAsync(String)](2cf895fa-1033-1def-f048-0e8caefaa71d.htm) | 异步读取双浮点的数据  Asynchronously read double floating point data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadDoubleAsync(String, UInt16)](2256d406-d4d5-da54-210b-ccef58a5e914.htm) | 异步读取双浮点数据的数组  Asynchronously read double floating point data array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadFloat(String)](bc2801ae-6600-e845-7c6c-522c9156eb0e.htm) | 读取单浮点数据  Read single floating point data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadFloat(String, UInt16)](6a9fb098-b13c-313b-3b27-20d3c787945d.htm) | 读取单浮点精度的数组  Read single floating point array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadFloatAsync(String)](89f743f8-4b3a-7085-1e49-a1d7adb10052.htm) | 异步读取单浮点数据  Asynchronously read single floating point data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadFloatAsync(String, UInt16)](cb6d5c5d-9b5c-2467-3e6b-5efa9d20cd72.htm) | 异步读取单浮点精度的数组  Asynchronously read single floating point array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServer(Byte)](8146dc0c-1edc-3ca7-428a-a3bdbd35e6c5.htm) | 将二进制的数据发送到管道中去，然后从管道里接收二进制的数据回来，并返回是否成功的结果对象。  Send binary data to the pipeline, and then receive binary data back from the pipeline, and return whether the success of the result object (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServer(IEnumerableByte)](044cfdd9-ec76-b89f-2f2b-5f439025006c.htm) | 将多个数据报文按顺序发到设备，并从设备接收返回的数据内容，然后拼接成一个Byte[]信息，需要重写[UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm)方法才能返回正确的结果。  Send multiple data packets to the device in sequence, and receive the returned data content from the device, and then splicing them into a Byte[] message, you need to rewrite [UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm) method to return the correct result. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServer(Byte, Boolean, Boolean)](b32ce8e8-a959-fcc2-6736-14b39fdb9a40.htm) | 将二进制的数据发送到管道中去，然后从管道里接收二进制的数据回来，并返回是否成功的结果对象。  Send binary data to the pipeline, and then receive binary data back from the pipeline, and return whether the success of the result object (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServer(CommunicationPipe, Byte, Boolean, Boolean)](ef36d56b-df09-f727-e8d3-316d37a13ee9.htm) | 使用指定的管道来进行数据通信，发送原始数据到管道，然后从管道接收相关的数据返回，本方法无锁 (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServerAsync(Byte)](a4de5a24-c67a-7b02-db05-33b85eaf6e1e.htm) | 将二进制的数据发送到管道中去，然后从管道里接收二进制的数据回来，并返回是否成功的结果对象。  Send binary data to the pipeline, and then receive binary data back from the pipeline, and return whether the success of the result object (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServerAsync(IEnumerableByte)](bcdb59b7-a70a-d5c0-1fc6-9ececcd39ea2.htm) | 将多个数据报文按顺序发到设备，并从设备接收返回的数据内容，然后拼接成一个Byte[]信息，需要重写[UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm)方法才能返回正确的结果。  Send multiple data packets to the device in sequence, and receive the returned data content from the device, and then splicing them into a Byte[] message, you need to rewrite [UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm) method to return the correct result. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServerAsync(Byte, Boolean, Boolean)](1c79c204-5ab5-9342-e65d-7455c20ea577.htm) | 将二进制的数据发送到管道中去，然后从管道里接收二进制的数据回来，并返回是否成功的结果对象。  Send binary data to the pipeline, and then receive binary data back from the pipeline, and return whether the success of the result object (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServerAsync(CommunicationPipe, Byte, Boolean, Boolean)](a254a19d-3c29-147b-335b-8adccdf4ee2f.htm) | 使用指定的管道来进行数据通信，发送原始数据到管道，然后从管道接收相关的数据返回，本方法无锁 (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法代码示例 | [ReadInt16(String)](57a16003-2fd5-b2bd-09bf-14366be899b8.htm) | 读取16位的有符号的整型数据  Read 16-bit signed integer data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt16(String, UInt16)](b2bde158-f51c-93fc-66ca-5e4cf6ea1fb7.htm) | 读取16位的有符号整型数组  Read 16-bit signed integer array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt16Async(String)](ac021793-890d-ceab-eae4-47ff763390c1.htm) | 异步读取16位的有符号的整型数据  Asynchronously read 16-bit signed integer data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt16Async(String, UInt16)](f30e02fc-8078-f1b3-c226-5991e3e23caf.htm) | 异步读取16位的有符号整型数组  Asynchronously read 16-bit signed integer array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt32(String)](bfbea035-9a83-40e3-084c-28122041b37a.htm) | 读取32位的有符号整型  Read 32-bit signed integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt32(String, UInt16)](095240a3-bc0f-8647-4867-88e40cc089e3.htm) | 读取32位有符号整型数组  Read 32-bit signed integer array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt32Async(String)](96726ea0-82aa-1229-21f6-43cd62abca43.htm) | 异步读取32位的有符号整型  Asynchronously read 32-bit signed integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt32Async(String, UInt16)](d6356095-a289-18eb-136c-b7e6ae099616.htm) | 异步读取32位有符号整型数组  Asynchronously read 32-bit signed integer array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt64(String)](0d0c2284-2da1-6f9a-b353-537561cd5e95.htm) | 读取64位的有符号整型  Read 64-bit signed integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt64(String, UInt16)](35a55019-efb1-bf2d-be1b-01c8eed99243.htm) | 读取64位的有符号整型数组  Read 64-bit signed integer array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt64Async(String)](26d46c8a-3775-3abd-d3ec-5cd1758a317e.htm) | 异步读取64位的有符号整型  Asynchronously read 64-bit signed integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt64Async(String, UInt16)](8dc97c4a-8904-9a34-9ddc-76c1a9477ee5.htm) | 异步读取64位的有符号整型数组  Asynchronously read 64-bit signed integer array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadString(String, UInt16)](07812d76-415a-6400-c8a2-1c21da51e7a3.htm) | 读取字符串数据，默认为最常见的ASCII编码  Read string data, default is the most common ASCII encoding (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadString(String, UInt16, Encoding)](b34c71a6-58f3-6bf4-0b93-471be4f2ed31.htm) | 使用指定的编码，读取字符串数据  Reads string data using the specified encoding (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadStringAsync(String, UInt16)](1835da81-b811-d151-0b53-907d8d87d0a8.htm) | 异步读取字符串数据，默认为最常见的ASCII编码  Asynchronously read string data, default is the most common ASCII encoding (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadStringAsync(String, UInt16, Encoding)](091dfc67-3d25-f33e-3a55-a0228863815e.htm) | 异步使用指定的编码，读取字符串数据  Asynchronously reads string data using the specified encoding (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadStructT](49853082-fbad-52e5-d756-0615cedb4b83.htm) | 读取结构体类型的数据，根据结构体自身的定义，读取原始字节数组，然后解析出实际的结构体数据，结构体需要实现[HslStructAttribute](87e427a8-0c7b-43a8-c421-c27264eda90f.htm)特性 (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadStructAsyncT](36461e61-d0ed-62cf-79ce-dcc7f71a3fa2.htm) | 读取结构体类型的数据，根据结构体自身的定义，读取原始字节数组，然后解析出实际的结构体数据，结构体需要实现[HslStructAttribute](87e427a8-0c7b-43a8-c421-c27264eda90f.htm)特性 (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt16(String)](736379d1-9560-5052-658a-13c126d56d78.htm) | 读取16位的无符号整型  Read 16-bit unsigned integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt16(String, UInt16)](e284c4bd-d3f8-ff30-eb42-0aa035cfb4f6.htm) | 读取16位的无符号整型数组  Read 16-bit unsigned integer array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt16Async(String)](361acc93-43b5-a9a5-6fa7-fa8728b7b96c.htm) | 异步读取16位的无符号整型  Asynchronously read 16-bit unsigned integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt16Async(String, UInt16)](73dd24b8-b501-df53-1de6-ccb2708f656d.htm) | 异步读取16位的无符号整型数组  Asynchronously read 16-bit unsigned integer array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt32(String)](ca73c958-39c7-9b34-aa44-3db2d9b0cad3.htm) | 读取32位的无符号整型  Read 32-bit unsigned integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt32(String, UInt16)](878b4709-0278-fcdf-3b65-afbf0f496c1e.htm) | 读取32位的无符号整型数组  Read 32-bit unsigned integer array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt32Async(String)](2343953f-6350-331e-6ba7-18a4449f0575.htm) | 异步读取32位的无符号整型  Asynchronously read 32-bit unsigned integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt32Async(String, UInt16)](e0eba6c1-56a2-5190-ba03-66116bcafe3d.htm) | 异步读取32位的无符号整型数组  Asynchronously read 32-bit unsigned integer array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt64(String)](947487dd-d4a3-72d8-d9fb-2224aad70e12.htm) | 读取64位的无符号整型  Read 64-bit unsigned integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt64(String, UInt16)](858766a9-9d62-87a8-f727-a00838f3d2b6.htm) | 读取64位的无符号整型的数组  Read 64-bit unsigned integer array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt64Async(String)](3411ff3a-772f-dbe1-8931-75e00a731579.htm) | 异步读取64位的无符号整型  Asynchronously read 64-bit unsigned integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt64Async(String, UInt16)](c2095a6e-ff89-dd98-8a21-fe72f82dbbc6.htm) | 异步读取64位的无符号整型的数组  Asynchronously read 64-bit unsigned integer array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [SendFrameIMessage](bc1483b2-9ec7-cfb8-d087-174ac72f0bd5.htm) | 以I消息的格式发送传入的原始字节数据，传入的消息为asdu信息 |
| 公共方法 | [SendFrameUMessage](68e6eb25-e268-6630-c89f-dfc08716f062.htm) | 向设备发送U帧消息的报文信息，传入功能码，STARTDT: 0x07, STOPDT: 0x13; TESTFR: 0x43 |
| 公共方法 | [SetDtuPipe](f765d1d7-234a-b229-9457-97e25e024dc9.htm) | 将当前的通信对象设置DTU模式，允许传入现成的管道，并返回初始化结果，如果该设备重写了握手报文，就是返回握手结果  Set the current communication object to DTU mode, allow the existing pipe to be passed in, and return the initialization result, if the device rewrites the handshake packet, the handshake result is returned (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [SetDtuPipeAsync](9d1fd12d-b804-3828-a627-a2c0eef1d23f.htm) | 将当前的通信对象设置DTU模式，允许传入现成的管道，并返回初始化结果，如果该设备重写了握手报文，就是返回握手结果  Set the current communication object to DTU mode, allow the existing pipe to be passed in, and return the initialization result, if the device rewrites the handshake packet, the handshake result is returned (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [SetPersistentConnection](f903dae0-4a4e-c7db-8350-50e36b257ec5.htm) | **已过时。** V11版本及之前设置长连接的方法，在V12版本以上中没有任何效果，默认长连接，删除调用即可，此处保留方法是为了部分用户保持兼容性升级。  The method of setting the long connection in V11 and before, has no effect in V12 and above. this method can be deleted. The method is retained here to maintain compatibility upgrades for some users. (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共方法 | [ToString](173ac89a-23b2-e945-1453-9842f68e8242.htm) | (重写 [DeviceTcpNetToString](209a196b-90ea-2b73-e915-ce4b11f1263d.htm).) |
| 公共方法 | [TotalSubscriptions(Byte)](9313534c-3b2a-c4e0-50c2-4c5f459032d5.htm) | 总召唤，支持设置总召唤代码，0x01: 定时主动上送, 0x14: 相应总召唤上送, 0x03: 越限上送 |
| 公共方法 | [TotalSubscriptions(Byte, Byte)](c907771b-d52a-ed4b-50ad-1838d7120c51.htm) | 总召唤，支持设置总召唤代码，0x01: 定时主动上送, 0x14: 相应总召唤上送, 0x03: 越限上送 |
| 公共方法 | [TotalSubscriptions(UInt16, Byte, Byte)](b1ebf604-869a-1264-da60-eb36acf8dbf2.htm) | 总召唤，支持设置总召唤代码，0x01: 定时主动上送, 0x14: 相应总召唤上送, 0x03: 越限上送 |
| 公共方法 | [UnpackResponseContent](208017a7-8da2-33f8-3273-d7e7e850100f.htm) | 根据对方返回的报文命令，对命令进行基本的拆包，例如各种Modbus协议拆包为统一的核心报文，还支持对报文的验证  According to the message command returned by the other party, the command is basically unpacked, for example, various Modbus protocols are unpacked into a unified core message, and the verification of the message is also supported (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [Wait(String, Boolean, Int32, Int32)](0e4db26e-1177-2b18-9492-ba05afb6fcc4.htm) | 等待指定地址的Boolean值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the Boolean value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [Wait(String, Int16, Int32, Int32)](e06fe098-8c8b-83f9-73ca-2a7cce19d43d.htm) | 等待指定地址的Int16值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the Int16 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [Wait(String, Int32, Int32, Int32)](1bce0aae-27e9-f1d6-e8f3-25ae6b47554e.htm) | 等待指定地址的Int32值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the Int32 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [Wait(String, Int64, Int32, Int32)](073354e9-2508-2158-1b01-881c3a34844b.htm) | 等待指定地址的Int64值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the Int64 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [Wait(String, UInt16, Int32, Int32)](d97f8d26-1c1a-df69-93f0-f60209597ced.htm) | 等待指定地址的UInt16值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the UInt16 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [Wait(String, UInt32, Int32, Int32)](1e61c395-b07f-e66c-5f9c-2708093150d9.htm) | 等待指定地址的UInt32值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the UInt32 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [Wait(String, UInt64, Int32, Int32)](d659729b-5766-c497-d418-856c34ec2b24.htm) | 等待指定地址的UInt64值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the UInt64 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WaitAsync(String, Boolean, Int32, Int32)](65399e87-9d38-f819-ccd3-005f68c2c3b3.htm) | 等待指定地址的Boolean值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the Boolean value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WaitAsync(String, Int16, Int32, Int32)](01277e39-afb1-86a4-3aa3-560160e66fb9.htm) | 等待指定地址的Int16值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the Int16 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WaitAsync(String, Int32, Int32, Int32)](4a6d380f-6b2a-1270-07a7-39b88a2712fd.htm) | 等待指定地址的Int32值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the Int32 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WaitAsync(String, Int64, Int32, Int32)](d7eede9b-4d44-479b-fb88-3ff2e1bb7568.htm) | 等待指定地址的Int64值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the Int64 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WaitAsync(String, UInt16, Int32, Int32)](0a8395c7-dc7f-4d5f-4316-98e18fe70f4c.htm) | 等待指定地址的UInt16值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the UInt16 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WaitAsync(String, UInt32, Int32, Int32)](19e9074e-99a2-99fd-91ef-a5314a115397.htm) | 等待指定地址的UInt32值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the UInt32 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WaitAsync(String, UInt64, Int32, Int32)](5109dac0-0ba2-cbb9-e383-1a869c6dc22e.htm) | 等待指定地址的UInt64值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the UInt64 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [Write(String, Boolean)](5b046c82-3d8e-550b-fdf3-41ca5558169b.htm) | 批量写入Boolean数组数据，返回是否成功  Batch write Boolean array data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [Write(String, Byte)](b5536dd6-ccde-6883-ac27-bff9e4b37806.htm) | 写入原始的byte数组数据到指定的地址，返回是否写入成功  Write the original byte array data to the specified address, and return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Double)](64c44576-7e95-cbcb-a8c9-e6b7ff7e3c50.htm) | 写入double数据，返回是否成功  Write double data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Double)](df56e6b1-0688-ece7-3bf1-9db32fec9f6e.htm) | 写入double数组，返回是否成功  Write double array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Int16)](0e23506e-b604-e67c-b5eb-06271b6c6335.htm) | 写入short数据，返回是否成功  Write short data, returns whether success (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Int16)](c7cfd67d-ca31-3087-c684-ee4c10780ab1.htm) | 写入short数组，返回是否成功  Write short array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Int32)](ba80d1a5-c8d8-7e0b-5c6b-9435cf20cadb.htm) | 写入int数据，返回是否成功  Write int data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Int32)](57dac21c-8b28-de7c-c213-77cda465cccb.htm) | 写入int[]数组，返回是否成功  Write int array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Int64)](ba246b33-5a12-a22c-6d40-0969a64b508b.htm) | 写入long数据，返回是否成功  Write long data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Int64)](5b910f9c-6e6b-bb9d-6842-1b496640c1f4.htm) | 写入long数组，返回是否成功  Write long array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Single)](a2265fa7-5b56-a5aa-9d21-a0dc5458cbf9.htm) | 写入float数据，返回是否成功  Write float data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Single)](bb76227e-3e95-8ca2-5dcf-c087d90b8a01.htm) | 写入float数组，返回是否成功  Write float array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, String)](4e0ec363-ba64-7ca7-b8d3-2e9a8087bd55.htm) | 写入字符串信息，编码为ASCII  Write string information, encoded as ASCII (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, UInt16)](e5bc9466-65b4-fc64-e0d0-1867b1113130.htm) | 写入ushort数据，返回是否成功  Write ushort data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, UInt16)](65544962-82ac-bb47-d72f-5dd0c3d4ac6c.htm) | 写入ushort数组，返回是否成功  Write ushort array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, UInt32)](c99d8246-f95e-4f79-738b-16e622072eb4.htm) | 写入uint数据，返回是否成功  Write uint data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, UInt32)](b80b3099-b34d-fde8-6a65-837da3bdc5ab.htm) | 写入uint[]数组，返回是否成功  Write uint array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, UInt64)](4924086e-9b36-d446-4ae4-6682134a6ad1.htm) | 写入ulong数据，返回是否成功  Write ulong data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, UInt64)](b1131bdd-8b51-ec6d-dd0e-fec1f6cb51b6.htm) | 写入ulong数组，返回是否成功  Write ulong array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [Write(String, Boolean)](44adb540-1a09-8257-dce4-a87dfef039c7.htm) | 写入单个的Boolean数据，返回是否成功  Write a single Boolean data, and return whether the write was successful (重写 [DeviceCommunicationWrite(String, Boolean)](06d0b9af-7f28-0b9c-098f-4da28e69504a.htm).) |
| 公共方法代码示例 | [Write(String, String, Int32)](02fe2ca9-49a3-ede4-d1af-8c6da1f644da.htm) | 写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为ASCII  Write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, it is filled with 0 and the encoding is ASCII. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, String, Encoding)](7b602661-902d-2b70-1fdb-ba87ecd307bf.htm) | 写入字符串信息，需要指定的编码信息  Write string information, need to specify the encoding information (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, String, Int32, Encoding)](e8cb30d9-4c5e-6c68-b9af-1b1b6c2dc714.htm) | 写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为指定的编码信息  Write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, then the operation is complemented with 0 , you should specified the encoding information (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteT(T)](9f48853f-fbeb-c37d-25a2-1c76d1dee407.htm) | 写入支持Hsl特性的数据，返回是否写入成功，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Write data that supports the Hsl attribute, and return whether the write was successful. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WriteAsync(String, Boolean)](10fe669c-805e-1e5a-f016-f319ad363e90.htm) | 异步批量写入Boolean数组数据，返回是否成功  Asynchronously batch write Boolean array data, return success (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WriteAsync(String, Boolean)](fe872dae-7ef9-12c1-d883-e92ba8676416.htm) | 异步批量写入Boolean数组数据，返回是否成功  Asynchronously batch write Boolean array data, return success (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Byte)](8f1cec30-8b86-6749-c8f3-7d39168d59df.htm) | 异步写入原始的byte数组数据到指定的地址，返回是否写入成功  Asynchronously writes the original byte array data to the specified address, and returns whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Double)](068adaa4-d1de-dffc-8376-6f483e2c1d96.htm) | 异步写入double数据，返回是否成功  Asynchronously write double data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Double)](36199b83-c4ea-b02f-3675-fb6ac48d0c59.htm) | 异步写入double数组，返回是否成功  Asynchronously write double array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int16)](ccd95ccc-10a3-af67-62e8-c2312b4871b0.htm) | 异步写入short数据，返回是否成功  Asynchronously write short data, returns whether success (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int16)](2625e868-7096-e526-ce25-b29afc639b41.htm) | 异步写入short数组，返回是否成功  Asynchronously write short array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int32)](df91d629-2080-328d-67cd-fbefb9a45404.htm) | 异步写入int数据，返回是否成功  Asynchronously write int data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int32)](99962a08-9624-cea6-6009-5e6f207236b7.htm) | 异步写入int[]数组，返回是否成功  Asynchronously write int array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int64)](ec6c1e24-a0d8-f1f5-56af-109b0c546fc8.htm) | 异步写入long数据，返回是否成功  Asynchronously write long data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int64)](4b2ab94c-2891-e0a9-f10b-b67c63bd1cf9.htm) | 异步写入long数组，返回是否成功  Asynchronously write long array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Single)](910189df-13f9-b7d0-a626-b5b25ce5d32a.htm) | 异步写入float数据，返回是否成功  Asynchronously write float data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Single)](dbb2359a-fb54-cf93-4413-20d5a4ef9d36.htm) | 异步写入float数组，返回是否成功  Asynchronously write float array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, String)](fe839c9a-46f3-e317-0155-b10a8a83a9b1.htm) | 异步写入字符串信息，编码为ASCII  Asynchronously write string information, encoded as ASCII (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt16)](a04f3aa2-a51e-26cd-974e-bf70a2db359d.htm) | 异步写入ushort数据，返回是否成功  Asynchronously write ushort data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt16)](34a1a321-972e-06a8-76fb-794dd1229005.htm) | 异步写入ushort数组，返回是否成功  Asynchronously write ushort array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt32)](abad215c-bc99-ac90-b9f3-7adb3a7cca8f.htm) | 异步写入uint数据，返回是否成功  Asynchronously write uint data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt32)](5b655bf3-5ba7-ad31-6897-89c59051ac60.htm) | 异步写入uint[]数组，返回是否成功  Asynchronously write uint array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt64)](e396074d-6cf1-b45d-5e8e-01026f6e391b.htm) | 异步写入ulong数据，返回是否成功  Asynchronously write ulong data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt64)](ae3df8ec-29c0-42ff-3d6a-68f496997b3d.htm) | 异步写入ulong数组，返回是否成功  Asynchronously write ulong array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WriteAsync(String, String, Int32)](b739b983-b03c-8ecf-2713-75f19233ae2f.htm) | 异步写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为ASCII  Asynchronously write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, it is filled with 0 and the encoding is ASCII. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, String, Encoding)](876e9800-0406-6407-bdf8-8ab1cf28df38.htm) | 异步写入字符串信息，需要指定的编码信息  Asynchronously write string information, need to specify the encoding information (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WriteAsync(String, String, Int32, Encoding)](ae9fb2a3-fddb-48eb-e9ad-6cc0ebbac1f6.htm) | 异步写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为指定的编码信息  Asynchronously write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, then the operation is complemented with 0 , you should specified the encoding information (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsyncT(T)](25d96e77-d1b2-918c-33f9-e5e4cc57c9cd.htm) | 异步写入支持Hsl特性的数据，返回是否写入成功，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Asynchronously write data that supports the Hsl attribute, and return whether the write was successful. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteCustomerT](84a37f82-3f67-a0f8-fd3f-0630772565e4.htm) | 写入自定义类型的数据，该类型必须继承自IDataTransfer接口  Write data of a custom type, which must inherit from the IDataTransfer interface (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteCustomerAsyncT](04628397-cfa1-d785-7790-4c02cb377ab6.htm) | 写入自定义类型的数据，该类型必须继承自IDataTransfer接口  Write data of a custom type, which must inherit from the IDataTransfer interface (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WriteIec(Byte, UInt16, UInt16, Boolean)](5deaa15f-bc22-08b9-dffe-dd8b16c8be8c.htm) | 指定类型标识，传送原因，数据地址，数据值信息来写入到IEC的仪表中去 |
| 公共方法 | [WriteIec(Byte, UInt16, UInt16, Byte)](fe546efd-9483-5c5b-2a89-e9c433ae3ed4.htm) | 指定类型标识，传送原因，数据地址，数据值信息来写入到IEC的仪表中去 |
| 公共方法 | [WriteIec(Byte, UInt16, UInt16, Int16)](fa5408d3-9794-f5db-0ac9-f700d71a2259.htm) | 指定类型标识，传送原因，数据地址，数据值信息来写入到IEC的仪表中去 |
| 公共方法 | [WriteIec(Byte, UInt16, UInt16, Single)](207e2965-2c89-32e9-cb00-934333cfdc5a.htm) | 指定类型标识，传送原因，数据地址，数据值信息来写入到IEC的仪表中去 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)事件

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共事件 | [OnIEC104MessageReceived](d5456fa4-f675-29ed-96d7-3026316bc386.htm) | 当接收到了IEC104的消息触发的事件 |

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

![](../icons/SectionExpanded.png)备注

手动绑定事件 [OnIEC104MessageReceived](d5456fa4-f675-29ed-96d7-3026316bc386.htm)，当收到设备的数据时就可以触发，如果需要总召唤，调用方法 [TotalSubscriptions(Byte)](9313534c-3b2a-c4e0-50c2-4c5f459032d5.htm)

![](../icons/SectionExpanded.png)示例

一般来说，实例化之后，绑定事件，连接即可，然后定时调用总召唤（定时调用时还可以判断是否连接成功）

自定义的使用

[复制](# "复制")

```
private IEC104 iec;

private void IecInitialize( )
{
    iec = new IEC104( "127.0.0.1", 2404 );
    iec.OnIEC104MessageReceived += Iec_OnIEC104MessageReceived;     // 接收设备的数据的事件
    OperateResult connect = iec.ConnectServer( );
    if ( connect.IsSuccess)
    {
        Console.WriteLine( "Connect IEC success" );
    }
    else
    {
        Console.WriteLine( "Connect IEC failed: " + connect.Message );
        return;
    }

    // 发起一次总召唤
    iec.TotalSubscriptions( );
}

private void Iec_OnIEC104MessageReceived( object sender, IEC104MessageEventArgs e )
{
    // 响应总召唤，或是设备的点位变化上报，都会触发这里的方法，具体的代码都在这里开发
    if (e.TypeID == IEC104.TypeID.M_SP_NA_1)
    {
        if (e.IsAddressContinuous)
        {
            // 例如这里收到了一个连续地址的 单点遥信，报文如下
            // ASDU: 01 A4 14 00 01 00 01 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
            // 含义                    [地址]   [0-36]个字节，表示36个开关量加品质
            ushort addressStart = BitConverter.ToUInt16( e.Body, 0 );     // Body 除去了ASDU头部6个字节 (类型，类型限定，传送原因，单元地址)
            bool[] values = new bool[e.InfoObjectCount];      // 总计有这么多个数据
            for (int i = 0; i < values.Length; i++)
            {
                values[i] = e.Body[3 + i].GetBoolByIndex( 0 );     // 单点遥控时，最低位表示通断
            }

            // addressStart 起始地址，values 就是通断值
        }
        else
        {
            // 如果是突发的变化，会收到
            // ASDU 01 01 03 00 01 00 01 00 00 01
            // 含义                   [地址]   [值]

            for (int i = 0; i < e.InfoObjectCount; i++)  // 可能有多个，一般是1个
            {
                ushort addressStart = BitConverter.ToUInt16( e.Body, 3 * i );
                bool value = e.Body[3 * i + 2].GetBoolByIndex( 0 );
            }

        }
    }
}
```

当然可以使用hsl辅助解析数据

便捷的解析

[复制](# "复制")

```
private void Iec_OnIEC104MessageReceived( object sender, IEC104MessageEventArgs e )
{
    // 响应总召唤，或是设备的点位变化上报，都会触发这里的方法，具体的代码都在这里开发
    if (e.TypeID == IEC104.TypeID.M_SP_NA_1)
    {
        // 无论这里收到了一个连续地址的还是非连续的 单点遥信，有36个值的话，就会解析成36个数据
        List<IecValueObject<byte>> list = IecValueObject<byte>.ParseYaoXinValue( e );
        for (int i = 0; i < list.Count; i++)
        {
            IecValueObject<byte> value = list[i];
            // value.Value  值，0关，1开
            // value.Quality  品质
            // value.Address  地址
        }
    }
    else if (e.TypeID == IEC104.TypeID.M_ME_ND_1)
    {
        // 归一化标度值
        List<IecValueObject<short>> list = IecValueObject<short>.ParseInt16Value( e );
        for (int i = 0; i < list.Count; i++)
        {
            IecValueObject<short> value = list[i];
            // value.Value  值
            // value.Quality  品质
            // value.Address  地址
        }
    }
    else if (e.TypeID == IEC104.TypeID.M_ME_TF_1 || e.TypeID == IEC104.TypeID.M_ME_NC_1)
    {
        // 短浮点遥测值，带不带时标都处理
        List<IecValueObject<float>> list = IecValueObject<float>.ParseFloatValue( e );
        for (int i = 0; i < list.Count; i++)
        {
            IecValueObject<float> value = list[i];
            // value.Value  值
            // value.Quality  品质
            // value.Address  地址
            // value.Time  如果带时标的话
        }
    }
}
```

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.Instrument.IEC 命名空间](4f75643b-b529-650f-8587-78378826eec2.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## IEC104 构造函数 

[原文連結](http://api.hslcommunication.cn/html/22a80517-47fd-8960-d826-72c0a07e021f.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.IEC](../html/4f75643b-b529-650f-8587-78378826eec2.htm "HslCommunication.Instrument.IEC")

[IEC104 类](../html/8695140f-80b4-71fc-7a09-eff994df68ca.htm "IEC104 类")

[IEC104 构造函数](../html/22a80517-47fd-8960-d826-72c0a07e021f.htm "IEC104 构造函数 ")

[IEC104 构造函数](../html/4d456238-2b09-50c6-42c0-cd6893cb5c7a.htm "IEC104 构造函数 ")

[IEC104 构造函数 (String, Int32)](../html/125f8e8d-2907-db1f-c4b1-7b76e85f102d.htm "IEC104 构造函数 (String, Int32)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IEC104 构造函数 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [IEC104](4d456238-2b09-50c6-42c0-cd6893cb5c7a.htm) | 实例化IEC104协议的通讯对象  Instantiate the communication object of the IEC104 protocol |
| 公共方法 | [IEC104(String, Int32)](125f8e8d-2907-db1f-c4b1-7b76e85f102d.htm) | 指定ip地址和端口号来实例化一个默认的对象  Specify the IP address and port number to instantiate a default object |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[IEC104 类](8695140f-80b4-71fc-7a09-eff994df68ca.htm)

[HslCommunication.Instrument.IEC 命名空间](4f75643b-b529-650f-8587-78378826eec2.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## IEC104 构造函数 

[原文連結](http://api.hslcommunication.cn/html/4d456238-2b09-50c6-42c0-cd6893cb5c7a.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.IEC](../html/4f75643b-b529-650f-8587-78378826eec2.htm "HslCommunication.Instrument.IEC")

[IEC104 类](../html/8695140f-80b4-71fc-7a09-eff994df68ca.htm "IEC104 类")

[IEC104 构造函数](../html/22a80517-47fd-8960-d826-72c0a07e021f.htm "IEC104 构造函数 ")

[IEC104 构造函数](../html/4d456238-2b09-50c6-42c0-cd6893cb5c7a.htm "IEC104 构造函数 ")

[IEC104 构造函数 (String, Int32)](../html/125f8e8d-2907-db1f-c4b1-7b76e85f102d.htm "IEC104 构造函数 (String, Int32)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IEC104 构造函数 |

实例化IEC104协议的通讯对象  
Instantiate the communication object of the IEC104 protocol

**命名空间：**
 [HslCommunication.Instrument.IEC](4f75643b-b529-650f-8587-78378826eec2.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public IEC104()
```

```
Public Sub New
```

```
public:
IEC104()
```

```
new : unit -> IEC104
```

![](../icons/SectionExpanded.png)参见

#### 引用

[IEC104 类](8695140f-80b4-71fc-7a09-eff994df68ca.htm)

[IEC104 重载](22a80517-47fd-8960-d826-72c0a07e021f.htm)

[HslCommunication.Instrument.IEC 命名空间](4f75643b-b529-650f-8587-78378826eec2.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## IEC104 构造函数 (String, Int32)

[原文連結](http://api.hslcommunication.cn/html/125f8e8d-2907-db1f-c4b1-7b76e85f102d.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.IEC](../html/4f75643b-b529-650f-8587-78378826eec2.htm "HslCommunication.Instrument.IEC")

[IEC104 类](../html/8695140f-80b4-71fc-7a09-eff994df68ca.htm "IEC104 类")

[IEC104 构造函数](../html/22a80517-47fd-8960-d826-72c0a07e021f.htm "IEC104 构造函数 ")

[IEC104 构造函数](../html/4d456238-2b09-50c6-42c0-cd6893cb5c7a.htm "IEC104 构造函数 ")

[IEC104 构造函数 (String, Int32)](../html/125f8e8d-2907-db1f-c4b1-7b76e85f102d.htm "IEC104 构造函数 (String, Int32)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IEC104 构造函数 (String, Int32) |

指定ip地址和端口号来实例化一个默认的对象  
Specify the IP address and port number to instantiate a default object

**命名空间：**
 [HslCommunication.Instrument.IEC](4f75643b-b529-650f-8587-78378826eec2.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public IEC104(
	string ipAddress,
	int port = 2404
)
```

```
Public Sub New ( 
	ipAddress As String,
	Optional port As Integer = 2404
)
```

```
public:
IEC104(
	String^ ipAddress, 
	int port = 2404
)
```

```
new : 
        ipAddress : string * 
        ?port : int 
(* Defaults:
        let _port = defaultArg port 2404
*)
-> IEC104
```

#### 参数

ipAddress
:   类型：SystemString  
    IEC104的Ip地址

port (Optional)
:   类型：SystemInt32  
    IEC104的端口, 默认是2404端口

![](../icons/SectionExpanded.png)参见

#### 引用

[IEC104 类](8695140f-80b4-71fc-7a09-eff994df68ca.htm)

[IEC104 重载](22a80517-47fd-8960-d826-72c0a07e021f.htm)

[HslCommunication.Instrument.IEC 命名空间](4f75643b-b529-650f-8587-78378826eec2.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## IEC104 属性

[原文連結](http://api.hslcommunication.cn/html/1da629b7-909c-dfb0-74a8-3a845a94fe2a.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.IEC](../html/4f75643b-b529-650f-8587-78378826eec2.htm "HslCommunication.Instrument.IEC")

[IEC104 类](../html/8695140f-80b4-71fc-7a09-eff994df68ca.htm "IEC104 类")

[IEC104 属性](../html/1da629b7-909c-dfb0-74a8-3a845a94fe2a.htm "IEC104 属性")

[Station 属性](../html/9fddc8b9-c09a-2499-b762-17943593b07a.htm "Station 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IEC104 属性 |

[IEC104](8695140f-80b4-71fc-7a09-eff994df68ca.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性代码示例 | [ByteTransform](f35118c9-3691-f752-ae6f-2c5549714faa.htm) | 当前的数据变换机制，当你需要从字节数据转换类型数据的时候需要。  The current data transformation mechanism is required when you need to convert type data from byte data. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共属性 | [CommunicationPipe](ad8e903b-d289-21b0-077e-15ede8005221.htm) | 获取或设置当前的管道信息，管道类型为[CommunicationPipe](ad8e903b-d289-21b0-077e-15ede8005221.htm)的继承类，内置了[PipeTcpNet](26b94bd5-2695-36e7-6378-c3d41846be8f.htm)管道，[PipeUdpNet](5a0b810e-187e-826b-0a13-34fd5576aa0e.htm)管道，[PipeSerialPort](9d924607-3dfc-3916-b532-c5dec48fe347.htm)管道等  Get or set the current pipeline information, the pipeline type is [CommunicationPipe](ad8e903b-d289-21b0-077e-15ede8005221.htm) inheritance class, [PipeTcpNet](26b94bd5-2695-36e7-6378-c3d41846be8f.htm) pipeline, [PipeUdpNet](5a0b810e-187e-826b-0a13-34fd5576aa0e.htm) pipeline, [PipeSerialPort](9d924607-3dfc-3916-b532-c5dec48fe347.htm) pipeline, etc (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [ConnectionId](ce1efb4b-fbc9-4683-163e-0a8e71e99d6b.htm) | 当前连接的唯一ID号，默认为长度20的guid码加随机数组成，方便列表管理，也可以自己指定  The unique ID number of the current connection. The default is a 20-digit guid code plus a random number. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性代码示例 | [ConnectTimeOut](71a651e5-737d-6e88-75a2-fb891b9af6cd.htm) | 获取或设置接收服务器反馈的时间，如果为负数，则不接收反馈   Gets or sets the time to receive server feedback, and if it is a negative number, does not receive feedback (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共属性代码示例 | [IpAddress](c378f4e8-3fdf-1932-7239-ed44bc5af584.htm) | 获取或是设置远程服务器的IP地址，如果是本机测试，那么需要设置为127.0.0.1   Get or set the IP address of the remote server. If it is a local test, then it needs to be set to 127.0.0.1 (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共属性 | [LocalBinding](622b418f-5969-3726-5dcc-04456c4b97a7.htm) | 获取或设置绑定的本地的IP地址和端口号信息，如果端口设置为0，代表任何可用的端口  Get or set the bound local IP address and port number information, if the port is set to 0, it means any available port (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共属性代码示例 | [LogNet](8c56f246-88ac-afb8-8c06-972fafe4f78a.htm) | 组件的日志工具，支持日志记录，只要实例化后，当前网络的基本信息，就以[DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm)等级进行输出  The component's logging tool supports logging. As long as the instantiation of the basic network information, the output will be output at [DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm) (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性代码示例 | [Port](6f1ad24a-c916-5b67-aa17-a8c7dc2acb2e.htm) | 获取或设置服务器的端口号，具体的值需要取决于对方的配置  Gets or sets the port number of the server. The specific value depends on the configuration of the other party. (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共属性代码示例 | [ReceiveTimeOut](250560a7-332b-ce1b-5eb0-e3abe4ee6012.htm) | 获取或设置接收服务器反馈的时间，如果为负数，则不接收反馈   Gets or sets the time to receive server feedback, and if it is a negative number, does not receive feedback (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [SendBeforeHex](72aa1de4-116f-abfe-4adb-814a0be0e560.htm) | 获取或设置在发送通信报文前追加发送的字节信息，HEX格式，通常用于lora组网时，需要携带 00 00 00 02 四个字节的站地址功能。  Obtain or set the byte information sent before sending communication packets, HEX format, usually used for LORA networking, you need to carry 00 00 00 02 four-byte station address function. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [SleepTime](86a5fbb9-3456-6beb-7113-5bbad5eabbc5.htm) | 获取或设置在正式接收对方返回数据前的时候，需要休息的时间，当设置为0的时候，不需要休息。  Get or set the time required to rest before officially receiving the data from the other party. When it is set to 0, no rest is required. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [SocketKeepAliveTime](96ab7703-8c4d-b4c2-ddd9-0dc93ebcfc8d.htm) | 获取或设置客户端的Socket的心跳时间信息，这个是Socket底层自动实现的心跳包，不基于协议层实现。默认小于0，不开启心跳检测，如果需要开启，设置 60\_000 比较合适，单位毫秒  Get or set the heartbeat time information of the Socket of the client. This is the heartbeat packet automatically implemented by the bottom layer of the Socket, not based on the protocol layer. The default value is less than 0, and heartbeat detection is not enabled. If you need to enable it, it is more appropriate to set 60\_000, in milliseconds. (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共属性 | [Station](9fddc8b9-c09a-2499-b762-17943593b07a.htm) | 公共单元地址，在低版本远动程序中，站地址范围一般在1~254，255表示全局地址。新远动程序则可支持1~65534为站地址，而65535为全局地址。默认值为1 |
| 受保护的属性 | [WordLength](f1a505cc-e85c-6c8b-b189-58a1ab5753f9.htm) | 一个字单位的数据表示的地址长度，西门子为2，三菱，欧姆龙，modbusTcp就为1，AB PLC无效  The address length represented by one word of data, Siemens is 2, Mitsubishi, Omron, modbusTcp is 1, AB PLC is invalid (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[IEC104 类](8695140f-80b4-71fc-7a09-eff994df68ca.htm)

[HslCommunication.Instrument.IEC 命名空间](4f75643b-b529-650f-8587-78378826eec2.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Station 属性 

[原文連結](http://api.hslcommunication.cn/html/9fddc8b9-c09a-2499-b762-17943593b07a.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.IEC](../html/4f75643b-b529-650f-8587-78378826eec2.htm "HslCommunication.Instrument.IEC")

[IEC104 类](../html/8695140f-80b4-71fc-7a09-eff994df68ca.htm "IEC104 类")

[IEC104 属性](../html/1da629b7-909c-dfb0-74a8-3a845a94fe2a.htm "IEC104 属性")

[Station 属性](../html/9fddc8b9-c09a-2499-b762-17943593b07a.htm "Station 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IEC104Station 属性 |

公共单元地址，在低版本远动程序中，站地址范围一般在1~254，255表示全局地址。新远动程序则可支持1~65534为站地址，而65535为全局地址。默认值为1

**命名空间：**
 [HslCommunication.Instrument.IEC](4f75643b-b529-650f-8587-78378826eec2.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public int Station { get; set; }
```

```
Public Property Station As Integer
	Get
	Set
```

```
public:
property int Station {
	int get ();
	void set (int value);
}
```

```
member Station : int with get, set
```

#### 属性值

类型：Int32

![](../icons/SectionExpanded.png)参见

#### 引用

[IEC104 类](8695140f-80b4-71fc-7a09-eff994df68ca.htm)

[HslCommunication.Instrument.IEC 命名空间](4f75643b-b529-650f-8587-78378826eec2.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## IEC104 方法

[原文連結](http://api.hslcommunication.cn/html/8eb8a1e0-92fb-2fdb-83ac-2cab7f59d96b.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.IEC](../html/4f75643b-b529-650f-8587-78378826eec2.htm "HslCommunication.Instrument.IEC")

[IEC104 类](../html/8695140f-80b4-71fc-7a09-eff994df68ca.htm "IEC104 类")

[IEC104 方法](../html/8eb8a1e0-92fb-2fdb-83ac-2cab7f59d96b.htm "IEC104 方法")

[DecideWhetherQAMessage 方法](../html/a1d1b211-1f5a-e398-8d60-90c959313942.htm "DecideWhetherQAMessage 方法 ")

[GetNewNetMessage 方法](../html/5272711a-7847-2bfd-ef08-bf54632cbef2.htm "GetNewNetMessage 方法 ")

[InitializationOnConnect 方法](../html/3acd2dc3-e3a2-b827-23ec-08a3fb8a6d6d.htm "InitializationOnConnect 方法 ")

[InitializationOnConnectAsync 方法](../html/deaf517e-2c31-d73d-562f-371bdeb4d6c2.htm "InitializationOnConnectAsync 方法 ")

[SendFrameIMessage 方法](../html/bc1483b2-9ec7-cfb8-d087-174ac72f0bd5.htm "SendFrameIMessage 方法 ")

[SendFrameUMessage 方法](../html/68e6eb25-e268-6630-c89f-dfc08716f062.htm "SendFrameUMessage 方法 ")

[ToString 方法](../html/173ac89a-23b2-e945-1453-9842f68e8242.htm "ToString 方法 ")

[TotalSubscriptions 方法](../html/260e2e2f-ecf3-963e-f919-792283695b35.htm "TotalSubscriptions 方法 ")

[Write 方法](../html/011ee6ef-ce1d-e012-5c77-11b63f1243cd.htm "Write 方法 ")

[WriteIec 方法](../html/6ad1a56e-d20c-5001-6574-3d159b60f874.htm "WriteIec 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IEC104 方法 |

[IEC104](8695140f-80b4-71fc-7a09-eff994df68ca.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法代码示例 | [ConnectClose](569a9c0e-7628-b8d5-e124-0853452747b8.htm) | 手动断开与远程服务器的连接，如果当前是长连接模式，那么就会切换到短连接模式  Manually disconnect from the remote server, if it is currently in long connection mode, it will switch to short connection mode (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共方法代码示例 | [ConnectCloseAsync](76be8150-ff22-95b9-16b1-7c3440c54ae8.htm) | 手动断开与远程服务器的连接，如果当前是长连接模式，那么就会切换到短连接模式  Manually disconnect from the remote server, if it is currently in long connection mode, it will switch to short connection mode (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共方法代码示例 | [ConnectServer](6299cc1f-2fab-b754-4597-1eddc2b21308.htm) | 尝试连接远程的服务器，如果连接成功，就切换短连接模式到长连接模式，后面的每次请求都共享一个通道，使得通讯速度更快速  Try to connect to a remote server. If the connection is successful, switch the short connection mode to the long connection mode. Each subsequent request will share a channel, making the communication speed faster. (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共方法代码示例 | [ConnectServerAsync](e7eff5c0-8e2a-3b6c-5027-16f4d8bcf5ae.htm) | 尝试连接远程的服务器，如果连接成功，就切换短连接模式到长连接模式，后面的每次请求都共享一个通道，使得通讯速度更快速  Try to connect to a remote server. If the connection is successful, switch the short connection mode to the long connection mode. Each subsequent request will share a channel, making the communication speed faster. (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 受保护的方法 | [DecideWhetherQAMessage](a1d1b211-1f5a-e398-8d60-90c959313942.htm) | 决定当前的消息是否是用于问答机制返回的消息，默认直接返回 true, 实际的情况需要根据协议进行重写方法  To determine whether the current message is the message returned by the question answering mechanism, the default is true. In actual cases, the rewriting method needs to be performed according to the protocol (重写 [BinaryCommunicationDecideWhetherQAMessage(CommunicationPipe, OperateResultByte)](acc40cda-e310-8537-6685-59b67eeb16ff.htm).) |
| 公共方法 | [Dispose](b9dda6bf-d342-254e-689a-93fa0bb265d1.htm) | (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 受保护的方法 | [Dispose(Boolean)](0bbb6ed3-111e-069c-e9f5-8cbc3f7f7b77.htm) | (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | [ExtraAfterReadFromCoreServer](06a49a15-a3ce-5e39-28c7-2b78851eca2e.htm) | 和服务器交互完成的时候调用的方法，可以根据读写结果进行一些额外的操作，具体的操作需要根据实际的需求来重写实现  The method called when the interaction with the server is completed can perform some additional operations based on the read and write results. The specific operations need to be rewritten according to actual needs. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法代码示例 | [ExtraOnDisconnect](de9dbba2-3b2c-1c5e-53c8-4e6f8a44e2c4.htm) | 根据实际的协议选择是否重写本方法，有些协议在断开连接之前，需要发送一些报文来关闭当前的网络通道  Select whether to rewrite this method according to the actual protocol. Some protocols need to send some packets to close the current network channel before disconnecting. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法代码示例 | [ExtraOnDisconnectAsync](210f593e-16d8-77c4-1178-58fbcfe9b27b.htm) | 根据实际的协议选择是否重写本方法，有些协议在断开连接之前，需要发送一些报文来关闭当前的网络通道  Select whether to rewrite this method according to the actual protocol. Some protocols need to send some packets to close the current network channel before disconnecting. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 受保护的方法 | [GetLogTextFromBinary](67f3b2f7-957c-fa0b-8320-6799237157db.htm) | 获取当前的报文进行日志记录的时候，是否使用二进制的格式记录，默认返回 [LogMsgFormatBinary](2ed2ff87-caec-0e44-f08b-e2b015225a58.htm)，重写可以根据session对象分别返回不同记录模式  Whether to log the current packet in binary format, the default return is [LogMsgFormatBinary](2ed2ff87-caec-0e44-f08b-e2b015225a58.htm). If you want to override it, different recording modes can be returned according to session (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [GetNewNetMessage](5272711a-7847-2bfd-ef08-bf54632cbef2.htm) | 获取一个新的消息对象的方法，需要在继承类里面进行重写  The method to get a new message object needs to be overridden in the inheritance class (重写 [BinaryCommunicationGetNewNetMessage](16636a1b-a6f0-9c0f-9d9c-c94e26c67459.htm).) |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法 | [GetWordLength](c27821ea-1739-8eb9-a848-23a00828fec8.htm) | 一个字单位的数据表示的地址长度，西门子为2，三菱，欧姆龙，modbusTcp就为1，AB PLC无效  The address length represented by one word of data, Siemens is 2, Mitsubishi, Omron, modbusTcp is 1, AB PLC is invalid (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 受保护的方法代码示例 | [InitializationOnConnect](3acd2dc3-e3a2-b827-23ec-08a3fb8a6d6d.htm) | 根据实际的协议选择是否重写本方法，有些协议在创建连接之后，需要进行一些初始化的信号握手，才能最终建立网络通道。  Whether to rewrite this method is based on the actual protocol. Some protocols require some initial signal handshake to establish a network channel after the connection is created. (重写 [BinaryCommunicationInitializationOnConnect](4ab6989f-a5b5-05fe-8b44-bfa45824a2a7.htm).) |
| 受保护的方法代码示例 | [InitializationOnConnectAsync](deaf517e-2c31-d73d-562f-371bdeb4d6c2.htm) | 根据实际的协议选择是否重写本方法，有些协议在创建连接之后，需要进行一些初始化的信号握手，才能最终建立网络通道。  Whether to rewrite this method is based on the actual protocol. Some protocols require some initial signal handshake to establish a network channel after the connection is created. (重写 [BinaryCommunicationInitializationOnConnectAsync](7da2af04-d4fc-5045-19a4-4b0d1722eef3.htm).) |
| 公共方法 | [IpAddressPing](0ac126b4-a830-a822-4ef5-855474942cd7.htm) | 对当前设备的IP地址进行PING的操作，返回PING的结果，正常来说，返回Success  PING the IP address of the current device and return the PING result. Normally, it returns Success (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 受保护的方法 | [LogRevcMessage(Byte)](dab1657f-1e75-7fe1-7003-ccedfc69855d.htm) | 使用日志记录一个接收的报文信息  Logs are used to record information about a received packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [LogRevcMessage(Byte, PipeSession)](0f51ae45-575e-4228-f978-01b37d6d5df1.htm) | 使用日志记录一个接收的报文信息  Logs are used to record information about a received packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [LogSendMessage(Byte)](4f90eecc-68a2-1ef2-44ef-369f19bed3ec.htm) | 使用日志记录一个发送的报文信息  Logs are used to record information about a send packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [LogSendMessage(Byte, PipeSession)](209caadf-4f0f-40aa-04c2-ebc61f8dab13.htm) | 使用日志记录一个发送的报文信息  Logs are used to record information about a send packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法 | [PackCommandWithHeader](d0b2cd76-c306-11fc-38c7-af8ad306cce7.htm) | 对当前的命令进行打包处理，通常是携带命令头内容，标记当前的命令的长度信息，需要进行重写，否则默认不打包  The current command is packaged, usually carrying the content of the command header, marking the length of the current command, and it needs to be rewritten, otherwise it is not packaged by default (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [Read(String, UInt16)](9dd09452-c200-bc65-75ff-69aa34b366d4.htm) | 批量读取字节数组信息，需要指定地址和长度，返回原始的字节数组  Batch read byte array information, need to specify the address and length, return the original byte array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadT](95abe064-5594-5ae7-a849-98c53f88fee0.htm) | 读取支持Hsl特性的数据内容，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明，或是官网的详细文档  Read the data content of the Hsl attribute. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadAsync(String, UInt16)](d8c14d37-56dc-2a31-1a19-3e2b5b65197d.htm) | 异步批量读取字节数组信息，需要指定地址和长度，返回原始的字节数组  Asynchronous batch read byte array information, need to specify the address and length, return the original byte array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadAsyncT](0522e740-4bf0-2e31-c4c5-e41384321b37.htm) | 异步读取支持Hsl特性的数据内容，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Asynchronously read the data content of the Hsl attribute. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadBool(String)](9b2faa2e-eaa9-7ace-c49f-7a0318f1df1a.htm) | 读取单个的Boolean数据信息  Read a single Boolean data message (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadBool(String, UInt16)](69d0d24f-34a7-5b41-5c1b-619edc7e6b54.htm) | 批量读取Boolean数组信息，需要指定地址和长度，返回Boolean 数组  Batch read Boolean array information, need to specify the address and length, return Boolean array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadBoolAsync(String)](cc37da1f-d673-fe67-78ee-4e465d29a2bd.htm) | 异步读取单个的Boolean数据信息  Asynchronously read a single Boolean data message (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadBoolAsync(String, UInt16)](a45d383d-0095-3ea7-32f8-e75af88d60f8.htm) | 异步批量读取Boolean数组信息，需要指定地址和长度，返回Boolean 数组  Asynchronously batch read Boolean array information, need to specify the address and length, return Boolean array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadCustomerT(String)](ca069688-d8f0-7dbd-e962-86a10dce83cf.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，返回一个新的类型的实例对象。  To read a custom data type, you need to inherit from the IDataTransfer interface and return an instance object of a new type. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadCustomerT(String, T)](ca7dc1a1-c97d-689f-9819-4dae4af67206.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，传入一个实例，对这个实例进行赋值，并返回该实例的对象。  To read a custom data type, you need to inherit from the IDataTransfer interface, pass in an instance, assign a value to this instance, and return the object of the instance. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadCustomerAsyncT(String)](f3fcda2b-255d-aadd-567e-c7b3df3636d9.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，返回一个新的类型的实例对象。  To read a custom data type, you need to inherit from the IDataTransfer interface and return an instance object of a new type. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadCustomerAsyncT(String, T)](527c9fae-4743-6ebd-7b9a-5ef3a8b7c455.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，传入一个实例，对这个实例进行赋值，并返回该实例的对象。  To read a custom data type, you need to inherit from the IDataTransfer interface, pass in an instance, assign a value to this instance, and return the object of the instance. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadDouble(String)](e7e31c43-71b4-e692-6bdc-ba9533e1d83d.htm) | 读取双浮点的数据  Read double floating point data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadDouble(String, UInt16)](4eee0f05-f861-4b48-3f06-99737fb24aee.htm) | 读取双浮点数据的数组  Read double floating point data array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadDoubleAsync(String)](2cf895fa-1033-1def-f048-0e8caefaa71d.htm) | 异步读取双浮点的数据  Asynchronously read double floating point data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadDoubleAsync(String, UInt16)](2256d406-d4d5-da54-210b-ccef58a5e914.htm) | 异步读取双浮点数据的数组  Asynchronously read double floating point data array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadFloat(String)](bc2801ae-6600-e845-7c6c-522c9156eb0e.htm) | 读取单浮点数据  Read single floating point data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadFloat(String, UInt16)](6a9fb098-b13c-313b-3b27-20d3c787945d.htm) | 读取单浮点精度的数组  Read single floating point array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadFloatAsync(String)](89f743f8-4b3a-7085-1e49-a1d7adb10052.htm) | 异步读取单浮点数据  Asynchronously read single floating point data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadFloatAsync(String, UInt16)](cb6d5c5d-9b5c-2467-3e6b-5efa9d20cd72.htm) | 异步读取单浮点精度的数组  Asynchronously read single floating point array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServer(Byte)](8146dc0c-1edc-3ca7-428a-a3bdbd35e6c5.htm) | 将二进制的数据发送到管道中去，然后从管道里接收二进制的数据回来，并返回是否成功的结果对象。  Send binary data to the pipeline, and then receive binary data back from the pipeline, and return whether the success of the result object (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServer(IEnumerableByte)](044cfdd9-ec76-b89f-2f2b-5f439025006c.htm) | 将多个数据报文按顺序发到设备，并从设备接收返回的数据内容，然后拼接成一个Byte[]信息，需要重写[UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm)方法才能返回正确的结果。  Send multiple data packets to the device in sequence, and receive the returned data content from the device, and then splicing them into a Byte[] message, you need to rewrite [UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm) method to return the correct result. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServer(Byte, Boolean, Boolean)](b32ce8e8-a959-fcc2-6736-14b39fdb9a40.htm) | 将二进制的数据发送到管道中去，然后从管道里接收二进制的数据回来，并返回是否成功的结果对象。  Send binary data to the pipeline, and then receive binary data back from the pipeline, and return whether the success of the result object (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServer(CommunicationPipe, Byte, Boolean, Boolean)](ef36d56b-df09-f727-e8d3-316d37a13ee9.htm) | 使用指定的管道来进行数据通信，发送原始数据到管道，然后从管道接收相关的数据返回，本方法无锁 (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServerAsync(Byte)](a4de5a24-c67a-7b02-db05-33b85eaf6e1e.htm) | 将二进制的数据发送到管道中去，然后从管道里接收二进制的数据回来，并返回是否成功的结果对象。  Send binary data to the pipeline, and then receive binary data back from the pipeline, and return whether the success of the result object (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServerAsync(IEnumerableByte)](bcdb59b7-a70a-d5c0-1fc6-9ececcd39ea2.htm) | 将多个数据报文按顺序发到设备，并从设备接收返回的数据内容，然后拼接成一个Byte[]信息，需要重写[UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm)方法才能返回正确的结果。  Send multiple data packets to the device in sequence, and receive the returned data content from the device, and then splicing them into a Byte[] message, you need to rewrite [UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm) method to return the correct result. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServerAsync(Byte, Boolean, Boolean)](1c79c204-5ab5-9342-e65d-7455c20ea577.htm) | 将二进制的数据发送到管道中去，然后从管道里接收二进制的数据回来，并返回是否成功的结果对象。  Send binary data to the pipeline, and then receive binary data back from the pipeline, and return whether the success of the result object (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServerAsync(CommunicationPipe, Byte, Boolean, Boolean)](a254a19d-3c29-147b-335b-8adccdf4ee2f.htm) | 使用指定的管道来进行数据通信，发送原始数据到管道，然后从管道接收相关的数据返回，本方法无锁 (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法代码示例 | [ReadInt16(String)](57a16003-2fd5-b2bd-09bf-14366be899b8.htm) | 读取16位的有符号的整型数据  Read 16-bit signed integer data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt16(String, UInt16)](b2bde158-f51c-93fc-66ca-5e4cf6ea1fb7.htm) | 读取16位的有符号整型数组  Read 16-bit signed integer array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt16Async(String)](ac021793-890d-ceab-eae4-47ff763390c1.htm) | 异步读取16位的有符号的整型数据  Asynchronously read 16-bit signed integer data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt16Async(String, UInt16)](f30e02fc-8078-f1b3-c226-5991e3e23caf.htm) | 异步读取16位的有符号整型数组  Asynchronously read 16-bit signed integer array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt32(String)](bfbea035-9a83-40e3-084c-28122041b37a.htm) | 读取32位的有符号整型  Read 32-bit signed integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt32(String, UInt16)](095240a3-bc0f-8647-4867-88e40cc089e3.htm) | 读取32位有符号整型数组  Read 32-bit signed integer array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt32Async(String)](96726ea0-82aa-1229-21f6-43cd62abca43.htm) | 异步读取32位的有符号整型  Asynchronously read 32-bit signed integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt32Async(String, UInt16)](d6356095-a289-18eb-136c-b7e6ae099616.htm) | 异步读取32位有符号整型数组  Asynchronously read 32-bit signed integer array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt64(String)](0d0c2284-2da1-6f9a-b353-537561cd5e95.htm) | 读取64位的有符号整型  Read 64-bit signed integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt64(String, UInt16)](35a55019-efb1-bf2d-be1b-01c8eed99243.htm) | 读取64位的有符号整型数组  Read 64-bit signed integer array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt64Async(String)](26d46c8a-3775-3abd-d3ec-5cd1758a317e.htm) | 异步读取64位的有符号整型  Asynchronously read 64-bit signed integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt64Async(String, UInt16)](8dc97c4a-8904-9a34-9ddc-76c1a9477ee5.htm) | 异步读取64位的有符号整型数组  Asynchronously read 64-bit signed integer array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadString(String, UInt16)](07812d76-415a-6400-c8a2-1c21da51e7a3.htm) | 读取字符串数据，默认为最常见的ASCII编码  Read string data, default is the most common ASCII encoding (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadString(String, UInt16, Encoding)](b34c71a6-58f3-6bf4-0b93-471be4f2ed31.htm) | 使用指定的编码，读取字符串数据  Reads string data using the specified encoding (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadStringAsync(String, UInt16)](1835da81-b811-d151-0b53-907d8d87d0a8.htm) | 异步读取字符串数据，默认为最常见的ASCII编码  Asynchronously read string data, default is the most common ASCII encoding (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadStringAsync(String, UInt16, Encoding)](091dfc67-3d25-f33e-3a55-a0228863815e.htm) | 异步使用指定的编码，读取字符串数据  Asynchronously reads string data using the specified encoding (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadStructT](49853082-fbad-52e5-d756-0615cedb4b83.htm) | 读取结构体类型的数据，根据结构体自身的定义，读取原始字节数组，然后解析出实际的结构体数据，结构体需要实现[HslStructAttribute](87e427a8-0c7b-43a8-c421-c27264eda90f.htm)特性 (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadStructAsyncT](36461e61-d0ed-62cf-79ce-dcc7f71a3fa2.htm) | 读取结构体类型的数据，根据结构体自身的定义，读取原始字节数组，然后解析出实际的结构体数据，结构体需要实现[HslStructAttribute](87e427a8-0c7b-43a8-c421-c27264eda90f.htm)特性 (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt16(String)](736379d1-9560-5052-658a-13c126d56d78.htm) | 读取16位的无符号整型  Read 16-bit unsigned integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt16(String, UInt16)](e284c4bd-d3f8-ff30-eb42-0aa035cfb4f6.htm) | 读取16位的无符号整型数组  Read 16-bit unsigned integer array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt16Async(String)](361acc93-43b5-a9a5-6fa7-fa8728b7b96c.htm) | 异步读取16位的无符号整型  Asynchronously read 16-bit unsigned integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt16Async(String, UInt16)](73dd24b8-b501-df53-1de6-ccb2708f656d.htm) | 异步读取16位的无符号整型数组  Asynchronously read 16-bit unsigned integer array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt32(String)](ca73c958-39c7-9b34-aa44-3db2d9b0cad3.htm) | 读取32位的无符号整型  Read 32-bit unsigned integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt32(String, UInt16)](878b4709-0278-fcdf-3b65-afbf0f496c1e.htm) | 读取32位的无符号整型数组  Read 32-bit unsigned integer array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt32Async(String)](2343953f-6350-331e-6ba7-18a4449f0575.htm) | 异步读取32位的无符号整型  Asynchronously read 32-bit unsigned integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt32Async(String, UInt16)](e0eba6c1-56a2-5190-ba03-66116bcafe3d.htm) | 异步读取32位的无符号整型数组  Asynchronously read 32-bit unsigned integer array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt64(String)](947487dd-d4a3-72d8-d9fb-2224aad70e12.htm) | 读取64位的无符号整型  Read 64-bit unsigned integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt64(String, UInt16)](858766a9-9d62-87a8-f727-a00838f3d2b6.htm) | 读取64位的无符号整型的数组  Read 64-bit unsigned integer array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt64Async(String)](3411ff3a-772f-dbe1-8931-75e00a731579.htm) | 异步读取64位的无符号整型  Asynchronously read 64-bit unsigned integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt64Async(String, UInt16)](c2095a6e-ff89-dd98-8a21-fe72f82dbbc6.htm) | 异步读取64位的无符号整型的数组  Asynchronously read 64-bit unsigned integer array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [SendFrameIMessage](bc1483b2-9ec7-cfb8-d087-174ac72f0bd5.htm) | 以I消息的格式发送传入的原始字节数据，传入的消息为asdu信息 |
| 公共方法 | [SendFrameUMessage](68e6eb25-e268-6630-c89f-dfc08716f062.htm) | 向设备发送U帧消息的报文信息，传入功能码，STARTDT: 0x07, STOPDT: 0x13; TESTFR: 0x43 |
| 公共方法 | [SetDtuPipe](f765d1d7-234a-b229-9457-97e25e024dc9.htm) | 将当前的通信对象设置DTU模式，允许传入现成的管道，并返回初始化结果，如果该设备重写了握手报文，就是返回握手结果  Set the current communication object to DTU mode, allow the existing pipe to be passed in, and return the initialization result, if the device rewrites the handshake packet, the handshake result is returned (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [SetDtuPipeAsync](9d1fd12d-b804-3828-a627-a2c0eef1d23f.htm) | 将当前的通信对象设置DTU模式，允许传入现成的管道，并返回初始化结果，如果该设备重写了握手报文，就是返回握手结果  Set the current communication object to DTU mode, allow the existing pipe to be passed in, and return the initialization result, if the device rewrites the handshake packet, the handshake result is returned (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [SetPersistentConnection](f903dae0-4a4e-c7db-8350-50e36b257ec5.htm) | **已过时。** V11版本及之前设置长连接的方法，在V12版本以上中没有任何效果，默认长连接，删除调用即可，此处保留方法是为了部分用户保持兼容性升级。  The method of setting the long connection in V11 and before, has no effect in V12 and above. this method can be deleted. The method is retained here to maintain compatibility upgrades for some users. (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共方法 | [ToString](173ac89a-23b2-e945-1453-9842f68e8242.htm) | (重写 [DeviceTcpNetToString](209a196b-90ea-2b73-e915-ce4b11f1263d.htm).) |
| 公共方法 | [TotalSubscriptions(Byte)](9313534c-3b2a-c4e0-50c2-4c5f459032d5.htm) | 总召唤，支持设置总召唤代码，0x01: 定时主动上送, 0x14: 相应总召唤上送, 0x03: 越限上送 |
| 公共方法 | [TotalSubscriptions(Byte, Byte)](c907771b-d52a-ed4b-50ad-1838d7120c51.htm) | 总召唤，支持设置总召唤代码，0x01: 定时主动上送, 0x14: 相应总召唤上送, 0x03: 越限上送 |
| 公共方法 | [TotalSubscriptions(UInt16, Byte, Byte)](b1ebf604-869a-1264-da60-eb36acf8dbf2.htm) | 总召唤，支持设置总召唤代码，0x01: 定时主动上送, 0x14: 相应总召唤上送, 0x03: 越限上送 |
| 公共方法 | [UnpackResponseContent](208017a7-8da2-33f8-3273-d7e7e850100f.htm) | 根据对方返回的报文命令，对命令进行基本的拆包，例如各种Modbus协议拆包为统一的核心报文，还支持对报文的验证  According to the message command returned by the other party, the command is basically unpacked, for example, various Modbus protocols are unpacked into a unified core message, and the verification of the message is also supported (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [Wait(String, Boolean, Int32, Int32)](0e4db26e-1177-2b18-9492-ba05afb6fcc4.htm) | 等待指定地址的Boolean值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the Boolean value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [Wait(String, Int16, Int32, Int32)](e06fe098-8c8b-83f9-73ca-2a7cce19d43d.htm) | 等待指定地址的Int16值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the Int16 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [Wait(String, Int32, Int32, Int32)](1bce0aae-27e9-f1d6-e8f3-25ae6b47554e.htm) | 等待指定地址的Int32值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the Int32 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [Wait(String, Int64, Int32, Int32)](073354e9-2508-2158-1b01-881c3a34844b.htm) | 等待指定地址的Int64值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the Int64 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [Wait(String, UInt16, Int32, Int32)](d97f8d26-1c1a-df69-93f0-f60209597ced.htm) | 等待指定地址的UInt16值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the UInt16 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [Wait(String, UInt32, Int32, Int32)](1e61c395-b07f-e66c-5f9c-2708093150d9.htm) | 等待指定地址的UInt32值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the UInt32 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [Wait(String, UInt64, Int32, Int32)](d659729b-5766-c497-d418-856c34ec2b24.htm) | 等待指定地址的UInt64值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the UInt64 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WaitAsync(String, Boolean, Int32, Int32)](65399e87-9d38-f819-ccd3-005f68c2c3b3.htm) | 等待指定地址的Boolean值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the Boolean value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WaitAsync(String, Int16, Int32, Int32)](01277e39-afb1-86a4-3aa3-560160e66fb9.htm) | 等待指定地址的Int16值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the Int16 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WaitAsync(String, Int32, Int32, Int32)](4a6d380f-6b2a-1270-07a7-39b88a2712fd.htm) | 等待指定地址的Int32值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the Int32 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WaitAsync(String, Int64, Int32, Int32)](d7eede9b-4d44-479b-fb88-3ff2e1bb7568.htm) | 等待指定地址的Int64值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the Int64 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WaitAsync(String, UInt16, Int32, Int32)](0a8395c7-dc7f-4d5f-4316-98e18fe70f4c.htm) | 等待指定地址的UInt16值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the UInt16 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WaitAsync(String, UInt32, Int32, Int32)](19e9074e-99a2-99fd-91ef-a5314a115397.htm) | 等待指定地址的UInt32值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the UInt32 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WaitAsync(String, UInt64, Int32, Int32)](5109dac0-0ba2-cbb9-e383-1a869c6dc22e.htm) | 等待指定地址的UInt64值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the UInt64 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [Write(String, Boolean)](5b046c82-3d8e-550b-fdf3-41ca5558169b.htm) | 批量写入Boolean数组数据，返回是否成功  Batch write Boolean array data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [Write(String, Byte)](b5536dd6-ccde-6883-ac27-bff9e4b37806.htm) | 写入原始的byte数组数据到指定的地址，返回是否写入成功  Write the original byte array data to the specified address, and return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Double)](64c44576-7e95-cbcb-a8c9-e6b7ff7e3c50.htm) | 写入double数据，返回是否成功  Write double data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Double)](df56e6b1-0688-ece7-3bf1-9db32fec9f6e.htm) | 写入double数组，返回是否成功  Write double array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Int16)](0e23506e-b604-e67c-b5eb-06271b6c6335.htm) | 写入short数据，返回是否成功  Write short data, returns whether success (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Int16)](c7cfd67d-ca31-3087-c684-ee4c10780ab1.htm) | 写入short数组，返回是否成功  Write short array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Int32)](ba80d1a5-c8d8-7e0b-5c6b-9435cf20cadb.htm) | 写入int数据，返回是否成功  Write int data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Int32)](57dac21c-8b28-de7c-c213-77cda465cccb.htm) | 写入int[]数组，返回是否成功  Write int array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Int64)](ba246b33-5a12-a22c-6d40-0969a64b508b.htm) | 写入long数据，返回是否成功  Write long data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Int64)](5b910f9c-6e6b-bb9d-6842-1b496640c1f4.htm) | 写入long数组，返回是否成功  Write long array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Single)](a2265fa7-5b56-a5aa-9d21-a0dc5458cbf9.htm) | 写入float数据，返回是否成功  Write float data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Single)](bb76227e-3e95-8ca2-5dcf-c087d90b8a01.htm) | 写入float数组，返回是否成功  Write float array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, String)](4e0ec363-ba64-7ca7-b8d3-2e9a8087bd55.htm) | 写入字符串信息，编码为ASCII  Write string information, encoded as ASCII (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, UInt16)](e5bc9466-65b4-fc64-e0d0-1867b1113130.htm) | 写入ushort数据，返回是否成功  Write ushort data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, UInt16)](65544962-82ac-bb47-d72f-5dd0c3d4ac6c.htm) | 写入ushort数组，返回是否成功  Write ushort array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, UInt32)](c99d8246-f95e-4f79-738b-16e622072eb4.htm) | 写入uint数据，返回是否成功  Write uint data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, UInt32)](b80b3099-b34d-fde8-6a65-837da3bdc5ab.htm) | 写入uint[]数组，返回是否成功  Write uint array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, UInt64)](4924086e-9b36-d446-4ae4-6682134a6ad1.htm) | 写入ulong数据，返回是否成功  Write ulong data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, UInt64)](b1131bdd-8b51-ec6d-dd0e-fec1f6cb51b6.htm) | 写入ulong数组，返回是否成功  Write ulong array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [Write(String, Boolean)](44adb540-1a09-8257-dce4-a87dfef039c7.htm) | 写入单个的Boolean数据，返回是否成功  Write a single Boolean data, and return whether the write was successful (重写 [DeviceCommunicationWrite(String, Boolean)](06d0b9af-7f28-0b9c-098f-4da28e69504a.htm).) |
| 公共方法代码示例 | [Write(String, String, Int32)](02fe2ca9-49a3-ede4-d1af-8c6da1f644da.htm) | 写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为ASCII  Write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, it is filled with 0 and the encoding is ASCII. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, String, Encoding)](7b602661-902d-2b70-1fdb-ba87ecd307bf.htm) | 写入字符串信息，需要指定的编码信息  Write string information, need to specify the encoding information (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, String, Int32, Encoding)](e8cb30d9-4c5e-6c68-b9af-1b1b6c2dc714.htm) | 写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为指定的编码信息  Write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, then the operation is complemented with 0 , you should specified the encoding information (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteT(T)](9f48853f-fbeb-c37d-25a2-1c76d1dee407.htm) | 写入支持Hsl特性的数据，返回是否写入成功，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Write data that supports the Hsl attribute, and return whether the write was successful. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WriteAsync(String, Boolean)](10fe669c-805e-1e5a-f016-f319ad363e90.htm) | 异步批量写入Boolean数组数据，返回是否成功  Asynchronously batch write Boolean array data, return success (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WriteAsync(String, Boolean)](fe872dae-7ef9-12c1-d883-e92ba8676416.htm) | 异步批量写入Boolean数组数据，返回是否成功  Asynchronously batch write Boolean array data, return success (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Byte)](8f1cec30-8b86-6749-c8f3-7d39168d59df.htm) | 异步写入原始的byte数组数据到指定的地址，返回是否写入成功  Asynchronously writes the original byte array data to the specified address, and returns whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Double)](068adaa4-d1de-dffc-8376-6f483e2c1d96.htm) | 异步写入double数据，返回是否成功  Asynchronously write double data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Double)](36199b83-c4ea-b02f-3675-fb6ac48d0c59.htm) | 异步写入double数组，返回是否成功  Asynchronously write double array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int16)](ccd95ccc-10a3-af67-62e8-c2312b4871b0.htm) | 异步写入short数据，返回是否成功  Asynchronously write short data, returns whether success (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int16)](2625e868-7096-e526-ce25-b29afc639b41.htm) | 异步写入short数组，返回是否成功  Asynchronously write short array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int32)](df91d629-2080-328d-67cd-fbefb9a45404.htm) | 异步写入int数据，返回是否成功  Asynchronously write int data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int32)](99962a08-9624-cea6-6009-5e6f207236b7.htm) | 异步写入int[]数组，返回是否成功  Asynchronously write int array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int64)](ec6c1e24-a0d8-f1f5-56af-109b0c546fc8.htm) | 异步写入long数据，返回是否成功  Asynchronously write long data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int64)](4b2ab94c-2891-e0a9-f10b-b67c63bd1cf9.htm) | 异步写入long数组，返回是否成功  Asynchronously write long array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Single)](910189df-13f9-b7d0-a626-b5b25ce5d32a.htm) | 异步写入float数据，返回是否成功  Asynchronously write float data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Single)](dbb2359a-fb54-cf93-4413-20d5a4ef9d36.htm) | 异步写入float数组，返回是否成功  Asynchronously write float array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, String)](fe839c9a-46f3-e317-0155-b10a8a83a9b1.htm) | 异步写入字符串信息，编码为ASCII  Asynchronously write string information, encoded as ASCII (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt16)](a04f3aa2-a51e-26cd-974e-bf70a2db359d.htm) | 异步写入ushort数据，返回是否成功  Asynchronously write ushort data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt16)](34a1a321-972e-06a8-76fb-794dd1229005.htm) | 异步写入ushort数组，返回是否成功  Asynchronously write ushort array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt32)](abad215c-bc99-ac90-b9f3-7adb3a7cca8f.htm) | 异步写入uint数据，返回是否成功  Asynchronously write uint data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt32)](5b655bf3-5ba7-ad31-6897-89c59051ac60.htm) | 异步写入uint[]数组，返回是否成功  Asynchronously write uint array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt64)](e396074d-6cf1-b45d-5e8e-01026f6e391b.htm) | 异步写入ulong数据，返回是否成功  Asynchronously write ulong data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt64)](ae3df8ec-29c0-42ff-3d6a-68f496997b3d.htm) | 异步写入ulong数组，返回是否成功  Asynchronously write ulong array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WriteAsync(String, String, Int32)](b739b983-b03c-8ecf-2713-75f19233ae2f.htm) | 异步写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为ASCII  Asynchronously write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, it is filled with 0 and the encoding is ASCII. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, String, Encoding)](876e9800-0406-6407-bdf8-8ab1cf28df38.htm) | 异步写入字符串信息，需要指定的编码信息  Asynchronously write string information, need to specify the encoding information (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WriteAsync(String, String, Int32, Encoding)](ae9fb2a3-fddb-48eb-e9ad-6cc0ebbac1f6.htm) | 异步写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为指定的编码信息  Asynchronously write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, then the operation is complemented with 0 , you should specified the encoding information (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsyncT(T)](25d96e77-d1b2-918c-33f9-e5e4cc57c9cd.htm) | 异步写入支持Hsl特性的数据，返回是否写入成功，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Asynchronously write data that supports the Hsl attribute, and return whether the write was successful. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteCustomerT](84a37f82-3f67-a0f8-fd3f-0630772565e4.htm) | 写入自定义类型的数据，该类型必须继承自IDataTransfer接口  Write data of a custom type, which must inherit from the IDataTransfer interface (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteCustomerAsyncT](04628397-cfa1-d785-7790-4c02cb377ab6.htm) | 写入自定义类型的数据，该类型必须继承自IDataTransfer接口  Write data of a custom type, which must inherit from the IDataTransfer interface (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WriteIec(Byte, UInt16, UInt16, Boolean)](5deaa15f-bc22-08b9-dffe-dd8b16c8be8c.htm) | 指定类型标识，传送原因，数据地址，数据值信息来写入到IEC的仪表中去 |
| 公共方法 | [WriteIec(Byte, UInt16, UInt16, Byte)](fe546efd-9483-5c5b-2a89-e9c433ae3ed4.htm) | 指定类型标识，传送原因，数据地址，数据值信息来写入到IEC的仪表中去 |
| 公共方法 | [WriteIec(Byte, UInt16, UInt16, Int16)](fa5408d3-9794-f5db-0ac9-f700d71a2259.htm) | 指定类型标识，传送原因，数据地址，数据值信息来写入到IEC的仪表中去 |
| 公共方法 | [WriteIec(Byte, UInt16, UInt16, Single)](207e2965-2c89-32e9-cb00-934333cfdc5a.htm) | 指定类型标识，传送原因，数据地址，数据值信息来写入到IEC的仪表中去 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[IEC104 类](8695140f-80b4-71fc-7a09-eff994df68ca.htm)

[HslCommunication.Instrument.IEC 命名空间](4f75643b-b529-650f-8587-78378826eec2.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## DecideWhetherQAMessage 方法 

[原文連結](http://api.hslcommunication.cn/html/a1d1b211-1f5a-e398-8d60-90c959313942.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.IEC](../html/4f75643b-b529-650f-8587-78378826eec2.htm "HslCommunication.Instrument.IEC")

[IEC104 类](../html/8695140f-80b4-71fc-7a09-eff994df68ca.htm "IEC104 类")

[IEC104 方法](../html/8eb8a1e0-92fb-2fdb-83ac-2cab7f59d96b.htm "IEC104 方法")

[DecideWhetherQAMessage 方法](../html/a1d1b211-1f5a-e398-8d60-90c959313942.htm "DecideWhetherQAMessage 方法 ")

[GetNewNetMessage 方法](../html/5272711a-7847-2bfd-ef08-bf54632cbef2.htm "GetNewNetMessage 方法 ")

[InitializationOnConnect 方法](../html/3acd2dc3-e3a2-b827-23ec-08a3fb8a6d6d.htm "InitializationOnConnect 方法 ")

[InitializationOnConnectAsync 方法](../html/deaf517e-2c31-d73d-562f-371bdeb4d6c2.htm "InitializationOnConnectAsync 方法 ")

[SendFrameIMessage 方法](../html/bc1483b2-9ec7-cfb8-d087-174ac72f0bd5.htm "SendFrameIMessage 方法 ")

[SendFrameUMessage 方法](../html/68e6eb25-e268-6630-c89f-dfc08716f062.htm "SendFrameUMessage 方法 ")

[ToString 方法](../html/173ac89a-23b2-e945-1453-9842f68e8242.htm "ToString 方法 ")

[TotalSubscriptions 方法](../html/260e2e2f-ecf3-963e-f919-792283695b35.htm "TotalSubscriptions 方法 ")

[Write 方法](../html/011ee6ef-ce1d-e012-5c77-11b63f1243cd.htm "Write 方法 ")

[WriteIec 方法](../html/6ad1a56e-d20c-5001-6574-3d159b60f874.htm "WriteIec 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IEC104DecideWhetherQAMessage 方法 |

决定当前的消息是否是用于问答机制返回的消息，默认直接返回 true, 实际的情况需要根据协议进行重写方法  
To determine whether the current message is the message returned by the question answering mechanism,
the default is true. In actual cases, the rewriting method needs to be performed according to the protocol

**命名空间：**
 [HslCommunication.Instrument.IEC](4f75643b-b529-650f-8587-78378826eec2.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
protected override bool DecideWhetherQAMessage(
	CommunicationPipe pipe,
	OperateResult<byte[]> receive
)
```

```
Protected Overrides Function DecideWhetherQAMessage ( 
	pipe As CommunicationPipe,
	receive As OperateResult(Of Byte())
) As Boolean
```

```
protected:
virtual bool DecideWhetherQAMessage(
	CommunicationPipe^ pipe, 
	OperateResult<array<unsigned char>^>^ receive
) override
```

```
abstract DecideWhetherQAMessage : 
        pipe : CommunicationPipe * 
        receive : OperateResult<byte[]> -> bool 
override DecideWhetherQAMessage : 
        pipe : CommunicationPipe * 
        receive : OperateResult<byte[]> -> bool
```

#### 参数

pipe
:   类型：[HslCommunication.Core.PipeCommunicationPipe](ab3f9c5c-144b-4209-4999-b382a1d94450.htm)  
    管道信息

receive
:   类型：[HslCommunicationOperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
    接收的数据信息

#### 返回值

类型：Boolean  
是否是问答的数据

![](../icons/SectionExpanded.png)参见

#### 引用

[IEC104 类](8695140f-80b4-71fc-7a09-eff994df68ca.htm)

[HslCommunication.Instrument.IEC 命名空间](4f75643b-b529-650f-8587-78378826eec2.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GetNewNetMessage 方法 

[原文連結](http://api.hslcommunication.cn/html/5272711a-7847-2bfd-ef08-bf54632cbef2.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.IEC](../html/4f75643b-b529-650f-8587-78378826eec2.htm "HslCommunication.Instrument.IEC")

[IEC104 类](../html/8695140f-80b4-71fc-7a09-eff994df68ca.htm "IEC104 类")

[IEC104 方法](../html/8eb8a1e0-92fb-2fdb-83ac-2cab7f59d96b.htm "IEC104 方法")

[DecideWhetherQAMessage 方法](../html/a1d1b211-1f5a-e398-8d60-90c959313942.htm "DecideWhetherQAMessage 方法 ")

[GetNewNetMessage 方法](../html/5272711a-7847-2bfd-ef08-bf54632cbef2.htm "GetNewNetMessage 方法 ")

[InitializationOnConnect 方法](../html/3acd2dc3-e3a2-b827-23ec-08a3fb8a6d6d.htm "InitializationOnConnect 方法 ")

[InitializationOnConnectAsync 方法](../html/deaf517e-2c31-d73d-562f-371bdeb4d6c2.htm "InitializationOnConnectAsync 方法 ")

[SendFrameIMessage 方法](../html/bc1483b2-9ec7-cfb8-d087-174ac72f0bd5.htm "SendFrameIMessage 方法 ")

[SendFrameUMessage 方法](../html/68e6eb25-e268-6630-c89f-dfc08716f062.htm "SendFrameUMessage 方法 ")

[ToString 方法](../html/173ac89a-23b2-e945-1453-9842f68e8242.htm "ToString 方法 ")

[TotalSubscriptions 方法](../html/260e2e2f-ecf3-963e-f919-792283695b35.htm "TotalSubscriptions 方法 ")

[Write 方法](../html/011ee6ef-ce1d-e012-5c77-11b63f1243cd.htm "Write 方法 ")

[WriteIec 方法](../html/6ad1a56e-d20c-5001-6574-3d159b60f874.htm "WriteIec 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IEC104GetNewNetMessage 方法 |

获取一个新的消息对象的方法，需要在继承类里面进行重写  
The method to get a new message object needs to be overridden in the inheritance class

**命名空间：**
 [HslCommunication.Instrument.IEC](4f75643b-b529-650f-8587-78378826eec2.htm)  
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

[IEC104 类](8695140f-80b4-71fc-7a09-eff994df68ca.htm)

[HslCommunication.Instrument.IEC 命名空间](4f75643b-b529-650f-8587-78378826eec2.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## InitializationOnConnect 方法 

[原文連結](http://api.hslcommunication.cn/html/3acd2dc3-e3a2-b827-23ec-08a3fb8a6d6d.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.IEC](../html/4f75643b-b529-650f-8587-78378826eec2.htm "HslCommunication.Instrument.IEC")

[IEC104 类](../html/8695140f-80b4-71fc-7a09-eff994df68ca.htm "IEC104 类")

[IEC104 方法](../html/8eb8a1e0-92fb-2fdb-83ac-2cab7f59d96b.htm "IEC104 方法")

[DecideWhetherQAMessage 方法](../html/a1d1b211-1f5a-e398-8d60-90c959313942.htm "DecideWhetherQAMessage 方法 ")

[GetNewNetMessage 方法](../html/5272711a-7847-2bfd-ef08-bf54632cbef2.htm "GetNewNetMessage 方法 ")

[InitializationOnConnect 方法](../html/3acd2dc3-e3a2-b827-23ec-08a3fb8a6d6d.htm "InitializationOnConnect 方法 ")

[InitializationOnConnectAsync 方法](../html/deaf517e-2c31-d73d-562f-371bdeb4d6c2.htm "InitializationOnConnectAsync 方法 ")

[SendFrameIMessage 方法](../html/bc1483b2-9ec7-cfb8-d087-174ac72f0bd5.htm "SendFrameIMessage 方法 ")

[SendFrameUMessage 方法](../html/68e6eb25-e268-6630-c89f-dfc08716f062.htm "SendFrameUMessage 方法 ")

[ToString 方法](../html/173ac89a-23b2-e945-1453-9842f68e8242.htm "ToString 方法 ")

[TotalSubscriptions 方法](../html/260e2e2f-ecf3-963e-f919-792283695b35.htm "TotalSubscriptions 方法 ")

[Write 方法](../html/011ee6ef-ce1d-e012-5c77-11b63f1243cd.htm "Write 方法 ")

[WriteIec 方法](../html/6ad1a56e-d20c-5001-6574-3d159b60f874.htm "WriteIec 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IEC104InitializationOnConnect 方法 |

根据实际的协议选择是否重写本方法，有些协议在创建连接之后，需要进行一些初始化的信号握手，才能最终建立网络通道。  
Whether to rewrite this method is based on the actual protocol. Some protocols require some initial signal handshake to establish a network channel after the connection is created.

**命名空间：**
 [HslCommunication.Instrument.IEC](4f75643b-b529-650f-8587-78378826eec2.htm)  
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

[IEC104 类](8695140f-80b4-71fc-7a09-eff994df68ca.htm)

[HslCommunication.Instrument.IEC 命名空间](4f75643b-b529-650f-8587-78378826eec2.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## InitializationOnConnectAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/deaf517e-2c31-d73d-562f-371bdeb4d6c2.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.IEC](../html/4f75643b-b529-650f-8587-78378826eec2.htm "HslCommunication.Instrument.IEC")

[IEC104 类](../html/8695140f-80b4-71fc-7a09-eff994df68ca.htm "IEC104 类")

[IEC104 方法](../html/8eb8a1e0-92fb-2fdb-83ac-2cab7f59d96b.htm "IEC104 方法")

[DecideWhetherQAMessage 方法](../html/a1d1b211-1f5a-e398-8d60-90c959313942.htm "DecideWhetherQAMessage 方法 ")

[GetNewNetMessage 方法](../html/5272711a-7847-2bfd-ef08-bf54632cbef2.htm "GetNewNetMessage 方法 ")

[InitializationOnConnect 方法](../html/3acd2dc3-e3a2-b827-23ec-08a3fb8a6d6d.htm "InitializationOnConnect 方法 ")

[InitializationOnConnectAsync 方法](../html/deaf517e-2c31-d73d-562f-371bdeb4d6c2.htm "InitializationOnConnectAsync 方法 ")

[SendFrameIMessage 方法](../html/bc1483b2-9ec7-cfb8-d087-174ac72f0bd5.htm "SendFrameIMessage 方法 ")

[SendFrameUMessage 方法](../html/68e6eb25-e268-6630-c89f-dfc08716f062.htm "SendFrameUMessage 方法 ")

[ToString 方法](../html/173ac89a-23b2-e945-1453-9842f68e8242.htm "ToString 方法 ")

[TotalSubscriptions 方法](../html/260e2e2f-ecf3-963e-f919-792283695b35.htm "TotalSubscriptions 方法 ")

[Write 方法](../html/011ee6ef-ce1d-e012-5c77-11b63f1243cd.htm "Write 方法 ")

[WriteIec 方法](../html/6ad1a56e-d20c-5001-6574-3d159b60f874.htm "WriteIec 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IEC104InitializationOnConnectAsync 方法 |

根据实际的协议选择是否重写本方法，有些协议在创建连接之后，需要进行一些初始化的信号握手，才能最终建立网络通道。  
Whether to rewrite this method is based on the actual protocol. Some protocols require some initial signal handshake to establish a network channel after the connection is created.

**命名空间：**
 [HslCommunication.Instrument.IEC](4f75643b-b529-650f-8587-78378826eec2.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
protected override Task<OperateResult> InitializationOnConnectAsync()
```

```
Protected Overrides Function InitializationOnConnectAsync As Task(Of OperateResult)
```

```
protected:
virtual Task<OperateResult^>^ InitializationOnConnectAsync() override
```

```
abstract InitializationOnConnectAsync : unit -> Task<OperateResult> 
override InitializationOnConnectAsync : unit -> Task<OperateResult>
```

#### 返回值

类型：Task[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
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

[IEC104 类](8695140f-80b4-71fc-7a09-eff994df68ca.htm)

[HslCommunication.Instrument.IEC 命名空间](4f75643b-b529-650f-8587-78378826eec2.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## SendFrameIMessage 方法 

[原文連結](http://api.hslcommunication.cn/html/bc1483b2-9ec7-cfb8-d087-174ac72f0bd5.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.IEC](../html/4f75643b-b529-650f-8587-78378826eec2.htm "HslCommunication.Instrument.IEC")

[IEC104 类](../html/8695140f-80b4-71fc-7a09-eff994df68ca.htm "IEC104 类")

[IEC104 方法](../html/8eb8a1e0-92fb-2fdb-83ac-2cab7f59d96b.htm "IEC104 方法")

[DecideWhetherQAMessage 方法](../html/a1d1b211-1f5a-e398-8d60-90c959313942.htm "DecideWhetherQAMessage 方法 ")

[GetNewNetMessage 方法](../html/5272711a-7847-2bfd-ef08-bf54632cbef2.htm "GetNewNetMessage 方法 ")

[InitializationOnConnect 方法](../html/3acd2dc3-e3a2-b827-23ec-08a3fb8a6d6d.htm "InitializationOnConnect 方法 ")

[InitializationOnConnectAsync 方法](../html/deaf517e-2c31-d73d-562f-371bdeb4d6c2.htm "InitializationOnConnectAsync 方法 ")

[SendFrameIMessage 方法](../html/bc1483b2-9ec7-cfb8-d087-174ac72f0bd5.htm "SendFrameIMessage 方法 ")

[SendFrameUMessage 方法](../html/68e6eb25-e268-6630-c89f-dfc08716f062.htm "SendFrameUMessage 方法 ")

[ToString 方法](../html/173ac89a-23b2-e945-1453-9842f68e8242.htm "ToString 方法 ")

[TotalSubscriptions 方法](../html/260e2e2f-ecf3-963e-f919-792283695b35.htm "TotalSubscriptions 方法 ")

[Write 方法](../html/011ee6ef-ce1d-e012-5c77-11b63f1243cd.htm "Write 方法 ")

[WriteIec 方法](../html/6ad1a56e-d20c-5001-6574-3d159b60f874.htm "WriteIec 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IEC104SendFrameIMessage 方法 |

以I消息的格式发送传入的原始字节数据，传入的消息为asdu信息

**命名空间：**
 [HslCommunication.Instrument.IEC](4f75643b-b529-650f-8587-78378826eec2.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public OperateResult SendFrameIMessage(
	byte[] asdu
)
```

```
Public Function SendFrameIMessage ( 
	asdu As Byte()
) As OperateResult
```

```
public:
OperateResult^ SendFrameIMessage(
	array<unsigned char>^ asdu
)
```

```
member SendFrameIMessage : 
        asdu : byte[] -> OperateResult 
```

#### 参数

asdu
:   类型：SystemByte  
    ASDU报文信息

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否发送成功

![](../icons/SectionExpanded.png)参见

#### 引用

[IEC104 类](8695140f-80b4-71fc-7a09-eff994df68ca.htm)

[HslCommunication.Instrument.IEC 命名空间](4f75643b-b529-650f-8587-78378826eec2.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## SendFrameUMessage 方法 

[原文連結](http://api.hslcommunication.cn/html/68e6eb25-e268-6630-c89f-dfc08716f062.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.IEC](../html/4f75643b-b529-650f-8587-78378826eec2.htm "HslCommunication.Instrument.IEC")

[IEC104 类](../html/8695140f-80b4-71fc-7a09-eff994df68ca.htm "IEC104 类")

[IEC104 方法](../html/8eb8a1e0-92fb-2fdb-83ac-2cab7f59d96b.htm "IEC104 方法")

[DecideWhetherQAMessage 方法](../html/a1d1b211-1f5a-e398-8d60-90c959313942.htm "DecideWhetherQAMessage 方法 ")

[GetNewNetMessage 方法](../html/5272711a-7847-2bfd-ef08-bf54632cbef2.htm "GetNewNetMessage 方法 ")

[InitializationOnConnect 方法](../html/3acd2dc3-e3a2-b827-23ec-08a3fb8a6d6d.htm "InitializationOnConnect 方法 ")

[InitializationOnConnectAsync 方法](../html/deaf517e-2c31-d73d-562f-371bdeb4d6c2.htm "InitializationOnConnectAsync 方法 ")

[SendFrameIMessage 方法](../html/bc1483b2-9ec7-cfb8-d087-174ac72f0bd5.htm "SendFrameIMessage 方法 ")

[SendFrameUMessage 方法](../html/68e6eb25-e268-6630-c89f-dfc08716f062.htm "SendFrameUMessage 方法 ")

[ToString 方法](../html/173ac89a-23b2-e945-1453-9842f68e8242.htm "ToString 方法 ")

[TotalSubscriptions 方法](../html/260e2e2f-ecf3-963e-f919-792283695b35.htm "TotalSubscriptions 方法 ")

[Write 方法](../html/011ee6ef-ce1d-e012-5c77-11b63f1243cd.htm "Write 方法 ")

[WriteIec 方法](../html/6ad1a56e-d20c-5001-6574-3d159b60f874.htm "WriteIec 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IEC104SendFrameUMessage 方法 |

向设备发送U帧消息的报文信息，传入功能码，STARTDT: 0x07, STOPDT: 0x13; TESTFR: 0x43

**命名空间：**
 [HslCommunication.Instrument.IEC](4f75643b-b529-650f-8587-78378826eec2.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public OperateResult SendFrameUMessage(
	byte controlField
)
```

```
Public Function SendFrameUMessage ( 
	controlField As Byte
) As OperateResult
```

```
public:
OperateResult^ SendFrameUMessage(
	unsigned char controlField
)
```

```
member SendFrameUMessage : 
        controlField : byte -> OperateResult 
```

#### 参数

controlField
:   类型：SystemByte  
    功能码，STARTDT: 0x07, STOPDT: 0x13; TESTFR: 0x43

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否发送成功

![](../icons/SectionExpanded.png)参见

#### 引用

[IEC104 类](8695140f-80b4-71fc-7a09-eff994df68ca.htm)

[HslCommunication.Instrument.IEC 命名空间](4f75643b-b529-650f-8587-78378826eec2.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ToString 方法 

[原文連結](http://api.hslcommunication.cn/html/173ac89a-23b2-e945-1453-9842f68e8242.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.IEC](../html/4f75643b-b529-650f-8587-78378826eec2.htm "HslCommunication.Instrument.IEC")

[IEC104 类](../html/8695140f-80b4-71fc-7a09-eff994df68ca.htm "IEC104 类")

[IEC104 方法](../html/8eb8a1e0-92fb-2fdb-83ac-2cab7f59d96b.htm "IEC104 方法")

[DecideWhetherQAMessage 方法](../html/a1d1b211-1f5a-e398-8d60-90c959313942.htm "DecideWhetherQAMessage 方法 ")

[GetNewNetMessage 方法](../html/5272711a-7847-2bfd-ef08-bf54632cbef2.htm "GetNewNetMessage 方法 ")

[InitializationOnConnect 方法](../html/3acd2dc3-e3a2-b827-23ec-08a3fb8a6d6d.htm "InitializationOnConnect 方法 ")

[InitializationOnConnectAsync 方法](../html/deaf517e-2c31-d73d-562f-371bdeb4d6c2.htm "InitializationOnConnectAsync 方法 ")

[SendFrameIMessage 方法](../html/bc1483b2-9ec7-cfb8-d087-174ac72f0bd5.htm "SendFrameIMessage 方法 ")

[SendFrameUMessage 方法](../html/68e6eb25-e268-6630-c89f-dfc08716f062.htm "SendFrameUMessage 方法 ")

[ToString 方法](../html/173ac89a-23b2-e945-1453-9842f68e8242.htm "ToString 方法 ")

[TotalSubscriptions 方法](../html/260e2e2f-ecf3-963e-f919-792283695b35.htm "TotalSubscriptions 方法 ")

[Write 方法](../html/011ee6ef-ce1d-e012-5c77-11b63f1243cd.htm "Write 方法 ")

[WriteIec 方法](../html/6ad1a56e-d20c-5001-6574-3d159b60f874.htm "WriteIec 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IEC104ToString 方法 |

[缺少 "M:HslCommunication.Instrument.IEC.IEC104.ToString" 的 <summary> 文档]

**命名空间：**
 [HslCommunication.Instrument.IEC](4f75643b-b529-650f-8587-78378826eec2.htm)  
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

[缺少 "M:HslCommunication.Instrument.IEC.IEC104.ToString" 的 <returns> 文档]

![](../icons/SectionExpanded.png)参见

#### 引用

[IEC104 类](8695140f-80b4-71fc-7a09-eff994df68ca.htm)

[HslCommunication.Instrument.IEC 命名空间](4f75643b-b529-650f-8587-78378826eec2.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## TotalSubscriptions 方法 

[原文連結](http://api.hslcommunication.cn/html/260e2e2f-ecf3-963e-f919-792283695b35.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.IEC](../html/4f75643b-b529-650f-8587-78378826eec2.htm "HslCommunication.Instrument.IEC")

[IEC104 类](../html/8695140f-80b4-71fc-7a09-eff994df68ca.htm "IEC104 类")

[IEC104 方法](../html/8eb8a1e0-92fb-2fdb-83ac-2cab7f59d96b.htm "IEC104 方法")

[TotalSubscriptions 方法](../html/260e2e2f-ecf3-963e-f919-792283695b35.htm "TotalSubscriptions 方法 ")

[TotalSubscriptions 方法 (Byte)](../html/9313534c-3b2a-c4e0-50c2-4c5f459032d5.htm "TotalSubscriptions 方法 (Byte)")

[TotalSubscriptions 方法 (Byte, Byte)](../html/c907771b-d52a-ed4b-50ad-1838d7120c51.htm "TotalSubscriptions 方法 (Byte, Byte)")

[TotalSubscriptions 方法 (UInt16, Byte, Byte)](../html/b1ebf604-869a-1264-da60-eb36acf8dbf2.htm "TotalSubscriptions 方法 (UInt16, Byte, Byte)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IEC104TotalSubscriptions 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [TotalSubscriptions(Byte)](9313534c-3b2a-c4e0-50c2-4c5f459032d5.htm) | 总召唤，支持设置总召唤代码，0x01: 定时主动上送, 0x14: 相应总召唤上送, 0x03: 越限上送 |
| 公共方法 | [TotalSubscriptions(Byte, Byte)](c907771b-d52a-ed4b-50ad-1838d7120c51.htm) | 总召唤，支持设置总召唤代码，0x01: 定时主动上送, 0x14: 相应总召唤上送, 0x03: 越限上送 |
| 公共方法 | [TotalSubscriptions(UInt16, Byte, Byte)](b1ebf604-869a-1264-da60-eb36acf8dbf2.htm) | 总召唤，支持设置总召唤代码，0x01: 定时主动上送, 0x14: 相应总召唤上送, 0x03: 越限上送 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[IEC104 类](8695140f-80b4-71fc-7a09-eff994df68ca.htm)

[HslCommunication.Instrument.IEC 命名空间](4f75643b-b529-650f-8587-78378826eec2.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## TotalSubscriptions 方法 (Byte)

[原文連結](http://api.hslcommunication.cn/html/9313534c-3b2a-c4e0-50c2-4c5f459032d5.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.IEC](../html/4f75643b-b529-650f-8587-78378826eec2.htm "HslCommunication.Instrument.IEC")

[IEC104 类](../html/8695140f-80b4-71fc-7a09-eff994df68ca.htm "IEC104 类")

[IEC104 方法](../html/8eb8a1e0-92fb-2fdb-83ac-2cab7f59d96b.htm "IEC104 方法")

[TotalSubscriptions 方法](../html/260e2e2f-ecf3-963e-f919-792283695b35.htm "TotalSubscriptions 方法 ")

[TotalSubscriptions 方法 (Byte)](../html/9313534c-3b2a-c4e0-50c2-4c5f459032d5.htm "TotalSubscriptions 方法 (Byte)")

[TotalSubscriptions 方法 (Byte, Byte)](../html/c907771b-d52a-ed4b-50ad-1838d7120c51.htm "TotalSubscriptions 方法 (Byte, Byte)")

[TotalSubscriptions 方法 (UInt16, Byte, Byte)](../html/b1ebf604-869a-1264-da60-eb36acf8dbf2.htm "TotalSubscriptions 方法 (UInt16, Byte, Byte)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IEC104TotalSubscriptions 方法 (Byte) |

总召唤，支持设置总召唤代码，0x01: 定时主动上送, 0x14: 相应总召唤上送, 0x03: 越限上送

**命名空间：**
 [HslCommunication.Instrument.IEC](4f75643b-b529-650f-8587-78378826eec2.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public OperateResult TotalSubscriptions(
	byte code = 20
)
```

```
Public Function TotalSubscriptions ( 
	Optional code As Byte = 20
) As OperateResult
```

```
public:
OperateResult^ TotalSubscriptions(
	unsigned char code = 20
)
```

```
member TotalSubscriptions : 
        ?code : byte 
(* Defaults:
        let _code = defaultArg code 20
*)
-> OperateResult 
```

#### 参数

code (Optional)
:   类型：SystemByte  
    总召唤的代码，0x01: 定时主动上送, 0x14: 相应总召唤上送, 0x03: 越限上送

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否召唤成功

![](../icons/SectionExpanded.png)参见

#### 引用

[IEC104 类](8695140f-80b4-71fc-7a09-eff994df68ca.htm)

[TotalSubscriptions 重载](260e2e2f-ecf3-963e-f919-792283695b35.htm)

[HslCommunication.Instrument.IEC 命名空间](4f75643b-b529-650f-8587-78378826eec2.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## TotalSubscriptions 方法 (Byte, Byte)

[原文連結](http://api.hslcommunication.cn/html/c907771b-d52a-ed4b-50ad-1838d7120c51.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.IEC](../html/4f75643b-b529-650f-8587-78378826eec2.htm "HslCommunication.Instrument.IEC")

[IEC104 类](../html/8695140f-80b4-71fc-7a09-eff994df68ca.htm "IEC104 类")

[IEC104 方法](../html/8eb8a1e0-92fb-2fdb-83ac-2cab7f59d96b.htm "IEC104 方法")

[TotalSubscriptions 方法](../html/260e2e2f-ecf3-963e-f919-792283695b35.htm "TotalSubscriptions 方法 ")

[TotalSubscriptions 方法 (Byte)](../html/9313534c-3b2a-c4e0-50c2-4c5f459032d5.htm "TotalSubscriptions 方法 (Byte)")

[TotalSubscriptions 方法 (Byte, Byte)](../html/c907771b-d52a-ed4b-50ad-1838d7120c51.htm "TotalSubscriptions 方法 (Byte, Byte)")

[TotalSubscriptions 方法 (UInt16, Byte, Byte)](../html/b1ebf604-869a-1264-da60-eb36acf8dbf2.htm "TotalSubscriptions 方法 (UInt16, Byte, Byte)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IEC104TotalSubscriptions 方法 (Byte, Byte) |

总召唤，支持设置总召唤代码，0x01: 定时主动上送, 0x14: 相应总召唤上送, 0x03: 越限上送

**命名空间：**
 [HslCommunication.Instrument.IEC](4f75643b-b529-650f-8587-78378826eec2.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public OperateResult TotalSubscriptions(
	byte code,
	byte reason
)
```

```
Public Function TotalSubscriptions ( 
	code As Byte,
	reason As Byte
) As OperateResult
```

```
public:
OperateResult^ TotalSubscriptions(
	unsigned char code, 
	unsigned char reason
)
```

```
member TotalSubscriptions : 
        code : byte * 
        reason : byte -> OperateResult 
```

#### 参数

code
:   类型：SystemByte  
    总召唤的代码，0x01: 定时主动上送, 0x14: 相应总召唤上送, 0x03: 越限上送

reason
:   类型：SystemByte  
    传送原因，通常可选 06:激活，07:激活确认

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否总召唤成功

![](../icons/SectionExpanded.png)参见

#### 引用

[IEC104 类](8695140f-80b4-71fc-7a09-eff994df68ca.htm)

[TotalSubscriptions 重载](260e2e2f-ecf3-963e-f919-792283695b35.htm)

[HslCommunication.Instrument.IEC 命名空间](4f75643b-b529-650f-8587-78378826eec2.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## TotalSubscriptions 方法 (UInt16, Byte, Byte)

[原文連結](http://api.hslcommunication.cn/html/b1ebf604-869a-1264-da60-eb36acf8dbf2.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.IEC](../html/4f75643b-b529-650f-8587-78378826eec2.htm "HslCommunication.Instrument.IEC")

[IEC104 类](../html/8695140f-80b4-71fc-7a09-eff994df68ca.htm "IEC104 类")

[IEC104 方法](../html/8eb8a1e0-92fb-2fdb-83ac-2cab7f59d96b.htm "IEC104 方法")

[TotalSubscriptions 方法](../html/260e2e2f-ecf3-963e-f919-792283695b35.htm "TotalSubscriptions 方法 ")

[TotalSubscriptions 方法 (Byte)](../html/9313534c-3b2a-c4e0-50c2-4c5f459032d5.htm "TotalSubscriptions 方法 (Byte)")

[TotalSubscriptions 方法 (Byte, Byte)](../html/c907771b-d52a-ed4b-50ad-1838d7120c51.htm "TotalSubscriptions 方法 (Byte, Byte)")

[TotalSubscriptions 方法 (UInt16, Byte, Byte)](../html/b1ebf604-869a-1264-da60-eb36acf8dbf2.htm "TotalSubscriptions 方法 (UInt16, Byte, Byte)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IEC104TotalSubscriptions 方法 (UInt16, Byte, Byte) |

总召唤，支持设置总召唤代码，0x01: 定时主动上送, 0x14: 相应总召唤上送, 0x03: 越限上送

**命名空间：**
 [HslCommunication.Instrument.IEC](4f75643b-b529-650f-8587-78378826eec2.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public OperateResult TotalSubscriptions(
	ushort station,
	byte code,
	byte reason
)
```

```
Public Function TotalSubscriptions ( 
	station As UShort,
	code As Byte,
	reason As Byte
) As OperateResult
```

```
public:
OperateResult^ TotalSubscriptions(
	unsigned short station, 
	unsigned char code, 
	unsigned char reason
)
```

```
member TotalSubscriptions : 
        station : uint16 * 
        code : byte * 
        reason : byte -> OperateResult 
```

#### 参数

station
:   类型：SystemUInt16  
    公共单元地址

code
:   类型：SystemByte  
    总召唤的代码，0x01: 定时主动上送, 0x14: 相应总召唤上送, 0x03: 越限上送

reason
:   类型：SystemByte  
    传送原因，通常可选 06:激活，07:激活确认

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否总召唤成功

![](../icons/SectionExpanded.png)参见

#### 引用

[IEC104 类](8695140f-80b4-71fc-7a09-eff994df68ca.htm)

[TotalSubscriptions 重载](260e2e2f-ecf3-963e-f919-792283695b35.htm)

[HslCommunication.Instrument.IEC 命名空间](4f75643b-b529-650f-8587-78378826eec2.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Write 方法 

[原文連結](http://api.hslcommunication.cn/html/011ee6ef-ce1d-e012-5c77-11b63f1243cd.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.IEC](../html/4f75643b-b529-650f-8587-78378826eec2.htm "HslCommunication.Instrument.IEC")

[IEC104 类](../html/8695140f-80b4-71fc-7a09-eff994df68ca.htm "IEC104 类")

[IEC104 方法](../html/8eb8a1e0-92fb-2fdb-83ac-2cab7f59d96b.htm "IEC104 方法")

[Write 方法](../html/011ee6ef-ce1d-e012-5c77-11b63f1243cd.htm "Write 方法 ")

[Write 方法 (String, Boolean)](../html/44adb540-1a09-8257-dce4-a87dfef039c7.htm "Write 方法 (String, Boolean)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IEC104Write 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法代码示例 | [WriteT(T)](9f48853f-fbeb-c37d-25a2-1c76d1dee407.htm) | 写入支持Hsl特性的数据，返回是否写入成功，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Write data that supports the Hsl attribute, and return whether the write was successful. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [Write(String, Boolean)](5b046c82-3d8e-550b-fdf3-41ca5558169b.htm) | 批量写入Boolean数组数据，返回是否成功  Batch write Boolean array data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [Write(String, Byte)](b5536dd6-ccde-6883-ac27-bff9e4b37806.htm) | 写入原始的byte数组数据到指定的地址，返回是否写入成功  Write the original byte array data to the specified address, and return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Double)](64c44576-7e95-cbcb-a8c9-e6b7ff7e3c50.htm) | 写入double数据，返回是否成功  Write double data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Double)](df56e6b1-0688-ece7-3bf1-9db32fec9f6e.htm) | 写入double数组，返回是否成功  Write double array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Int16)](0e23506e-b604-e67c-b5eb-06271b6c6335.htm) | 写入short数据，返回是否成功  Write short data, returns whether success (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Int16)](c7cfd67d-ca31-3087-c684-ee4c10780ab1.htm) | 写入short数组，返回是否成功  Write short array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Int32)](ba80d1a5-c8d8-7e0b-5c6b-9435cf20cadb.htm) | 写入int数据，返回是否成功  Write int data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Int32)](57dac21c-8b28-de7c-c213-77cda465cccb.htm) | 写入int[]数组，返回是否成功  Write int array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Int64)](ba246b33-5a12-a22c-6d40-0969a64b508b.htm) | 写入long数据，返回是否成功  Write long data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Int64)](5b910f9c-6e6b-bb9d-6842-1b496640c1f4.htm) | 写入long数组，返回是否成功  Write long array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Single)](a2265fa7-5b56-a5aa-9d21-a0dc5458cbf9.htm) | 写入float数据，返回是否成功  Write float data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Single)](bb76227e-3e95-8ca2-5dcf-c087d90b8a01.htm) | 写入float数组，返回是否成功  Write float array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, String)](4e0ec363-ba64-7ca7-b8d3-2e9a8087bd55.htm) | 写入字符串信息，编码为ASCII  Write string information, encoded as ASCII (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, UInt16)](e5bc9466-65b4-fc64-e0d0-1867b1113130.htm) | 写入ushort数据，返回是否成功  Write ushort data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, UInt16)](65544962-82ac-bb47-d72f-5dd0c3d4ac6c.htm) | 写入ushort数组，返回是否成功  Write ushort array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, UInt32)](c99d8246-f95e-4f79-738b-16e622072eb4.htm) | 写入uint数据，返回是否成功  Write uint data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, UInt32)](b80b3099-b34d-fde8-6a65-837da3bdc5ab.htm) | 写入uint[]数组，返回是否成功  Write uint array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, UInt64)](4924086e-9b36-d446-4ae4-6682134a6ad1.htm) | 写入ulong数据，返回是否成功  Write ulong data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, UInt64)](b1131bdd-8b51-ec6d-dd0e-fec1f6cb51b6.htm) | 写入ulong数组，返回是否成功  Write ulong array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [Write(String, Boolean)](44adb540-1a09-8257-dce4-a87dfef039c7.htm) | 写入单个的Boolean数据，返回是否成功  Write a single Boolean data, and return whether the write was successful (重写 [DeviceCommunicationWrite(String, Boolean)](06d0b9af-7f28-0b9c-098f-4da28e69504a.htm).) |
| 公共方法代码示例 | [Write(String, String, Int32)](02fe2ca9-49a3-ede4-d1af-8c6da1f644da.htm) | 写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为ASCII  Write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, it is filled with 0 and the encoding is ASCII. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, String, Encoding)](7b602661-902d-2b70-1fdb-ba87ecd307bf.htm) | 写入字符串信息，需要指定的编码信息  Write string information, need to specify the encoding information (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, String, Int32, Encoding)](e8cb30d9-4c5e-6c68-b9af-1b1b6c2dc714.htm) | 写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为指定的编码信息  Write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, then the operation is complemented with 0 , you should specified the encoding information (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[IEC104 类](8695140f-80b4-71fc-7a09-eff994df68ca.htm)

[HslCommunication.Instrument.IEC 命名空间](4f75643b-b529-650f-8587-78378826eec2.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Write 方法 (String, Boolean)

[原文連結](http://api.hslcommunication.cn/html/44adb540-1a09-8257-dce4-a87dfef039c7.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.IEC](../html/4f75643b-b529-650f-8587-78378826eec2.htm "HslCommunication.Instrument.IEC")

[IEC104 类](../html/8695140f-80b4-71fc-7a09-eff994df68ca.htm "IEC104 类")

[IEC104 方法](../html/8eb8a1e0-92fb-2fdb-83ac-2cab7f59d96b.htm "IEC104 方法")

[Write 方法](../html/011ee6ef-ce1d-e012-5c77-11b63f1243cd.htm "Write 方法 ")

[Write 方法 (String, Boolean)](../html/44adb540-1a09-8257-dce4-a87dfef039c7.htm "Write 方法 (String, Boolean)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IEC104Write 方法 (String, Boolean) |

写入单个的Boolean数据，返回是否成功  
Write a single Boolean data, and return whether the write was successful

**命名空间：**
 [HslCommunication.Instrument.IEC](4f75643b-b529-650f-8587-78378826eec2.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public override OperateResult Write(
	string address,
	bool value
)
```

```
Public Overrides Function Write ( 
	address As String,
	value As Boolean
) As OperateResult
```

```
public:
virtual OperateResult^ Write(
	String^ address, 
	bool value
) override
```

```
abstract Write : 
        address : string * 
        value : bool -> OperateResult 
override Write : 
        address : string * 
        value : bool -> OperateResult
```

#### 参数

address
:   类型：SystemString  
    起始地址

value
:   类型：SystemBoolean  
    写入值

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否写入成功

#### 实现

[IReadWriteNetWrite(String, Boolean)](4d742738-6864-465c-5382-bd49032eae37.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[IEC104 类](8695140f-80b4-71fc-7a09-eff994df68ca.htm)

[Write 重载](011ee6ef-ce1d-e012-5c77-11b63f1243cd.htm)

[HslCommunication.Instrument.IEC 命名空间](4f75643b-b529-650f-8587-78378826eec2.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## WriteIec 方法 

[原文連結](http://api.hslcommunication.cn/html/6ad1a56e-d20c-5001-6574-3d159b60f874.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.IEC](../html/4f75643b-b529-650f-8587-78378826eec2.htm "HslCommunication.Instrument.IEC")

[IEC104 类](../html/8695140f-80b4-71fc-7a09-eff994df68ca.htm "IEC104 类")

[IEC104 方法](../html/8eb8a1e0-92fb-2fdb-83ac-2cab7f59d96b.htm "IEC104 方法")

[WriteIec 方法](../html/6ad1a56e-d20c-5001-6574-3d159b60f874.htm "WriteIec 方法 ")

[WriteIec 方法 (Byte, UInt16, UInt16, Boolean)](../html/5deaa15f-bc22-08b9-dffe-dd8b16c8be8c.htm "WriteIec 方法 (Byte, UInt16, UInt16, Boolean)")

[WriteIec 方法 (Byte, UInt16, UInt16, Byte[])](../html/fe546efd-9483-5c5b-2a89-e9c433ae3ed4.htm "WriteIec 方法 (Byte, UInt16, UInt16, Byte[])")

[WriteIec 方法 (Byte, UInt16, UInt16, Int16)](../html/fa5408d3-9794-f5db-0ac9-f700d71a2259.htm "WriteIec 方法 (Byte, UInt16, UInt16, Int16)")

[WriteIec 方法 (Byte, UInt16, UInt16, Single)](../html/207e2965-2c89-32e9-cb00-934333cfdc5a.htm "WriteIec 方法 (Byte, UInt16, UInt16, Single)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IEC104WriteIec 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [WriteIec(Byte, UInt16, UInt16, Boolean)](5deaa15f-bc22-08b9-dffe-dd8b16c8be8c.htm) | 指定类型标识，传送原因，数据地址，数据值信息来写入到IEC的仪表中去 |
| 公共方法 | [WriteIec(Byte, UInt16, UInt16, Byte)](fe546efd-9483-5c5b-2a89-e9c433ae3ed4.htm) | 指定类型标识，传送原因，数据地址，数据值信息来写入到IEC的仪表中去 |
| 公共方法 | [WriteIec(Byte, UInt16, UInt16, Int16)](fa5408d3-9794-f5db-0ac9-f700d71a2259.htm) | 指定类型标识，传送原因，数据地址，数据值信息来写入到IEC的仪表中去 |
| 公共方法 | [WriteIec(Byte, UInt16, UInt16, Single)](207e2965-2c89-32e9-cb00-934333cfdc5a.htm) | 指定类型标识，传送原因，数据地址，数据值信息来写入到IEC的仪表中去 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[IEC104 类](8695140f-80b4-71fc-7a09-eff994df68ca.htm)

[HslCommunication.Instrument.IEC 命名空间](4f75643b-b529-650f-8587-78378826eec2.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## WriteIec 方法 (Byte, UInt16, UInt16, Boolean)

[原文連結](http://api.hslcommunication.cn/html/5deaa15f-bc22-08b9-dffe-dd8b16c8be8c.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.IEC](../html/4f75643b-b529-650f-8587-78378826eec2.htm "HslCommunication.Instrument.IEC")

[IEC104 类](../html/8695140f-80b4-71fc-7a09-eff994df68ca.htm "IEC104 类")

[IEC104 方法](../html/8eb8a1e0-92fb-2fdb-83ac-2cab7f59d96b.htm "IEC104 方法")

[WriteIec 方法](../html/6ad1a56e-d20c-5001-6574-3d159b60f874.htm "WriteIec 方法 ")

[WriteIec 方法 (Byte, UInt16, UInt16, Boolean)](../html/5deaa15f-bc22-08b9-dffe-dd8b16c8be8c.htm "WriteIec 方法 (Byte, UInt16, UInt16, Boolean)")

[WriteIec 方法 (Byte, UInt16, UInt16, Byte[])](../html/fe546efd-9483-5c5b-2a89-e9c433ae3ed4.htm "WriteIec 方法 (Byte, UInt16, UInt16, Byte[])")

[WriteIec 方法 (Byte, UInt16, UInt16, Int16)](../html/fa5408d3-9794-f5db-0ac9-f700d71a2259.htm "WriteIec 方法 (Byte, UInt16, UInt16, Int16)")

[WriteIec 方法 (Byte, UInt16, UInt16, Single)](../html/207e2965-2c89-32e9-cb00-934333cfdc5a.htm "WriteIec 方法 (Byte, UInt16, UInt16, Single)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IEC104WriteIec 方法 (Byte, UInt16, UInt16, Boolean) |

指定类型标识，传送原因，数据地址，数据值信息来写入到IEC的仪表中去

**命名空间：**
 [HslCommunication.Instrument.IEC](4f75643b-b529-650f-8587-78378826eec2.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public OperateResult WriteIec(
	byte type,
	ushort reason,
	ushort address,
	bool value
)
```

```
Public Function WriteIec ( 
	type As Byte,
	reason As UShort,
	address As UShort,
	value As Boolean
) As OperateResult
```

```
public:
OperateResult^ WriteIec(
	unsigned char type, 
	unsigned short reason, 
	unsigned short address, 
	bool value
)
```

```
member WriteIec : 
        type : byte * 
        reason : uint16 * 
        address : uint16 * 
        value : bool -> OperateResult 
```

#### 参数

type
:   类型：SystemByte  
    类型标识

reason
:   类型：SystemUInt16  
    传送原因

address
:   类型：SystemUInt16  
    信息对象地址

value
:   类型：SystemBoolean  
    值信息

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否发送到仪表成功

![](../icons/SectionExpanded.png)参见

#### 引用

[IEC104 类](8695140f-80b4-71fc-7a09-eff994df68ca.htm)

[WriteIec 重载](6ad1a56e-d20c-5001-6574-3d159b60f874.htm)

[HslCommunication.Instrument.IEC 命名空间](4f75643b-b529-650f-8587-78378826eec2.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## WriteIec 方法 (Byte, UInt16, UInt16, Byte[])

[原文連結](http://api.hslcommunication.cn/html/fe546efd-9483-5c5b-2a89-e9c433ae3ed4.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.IEC](../html/4f75643b-b529-650f-8587-78378826eec2.htm "HslCommunication.Instrument.IEC")

[IEC104 类](../html/8695140f-80b4-71fc-7a09-eff994df68ca.htm "IEC104 类")

[IEC104 方法](../html/8eb8a1e0-92fb-2fdb-83ac-2cab7f59d96b.htm "IEC104 方法")

[WriteIec 方法](../html/6ad1a56e-d20c-5001-6574-3d159b60f874.htm "WriteIec 方法 ")

[WriteIec 方法 (Byte, UInt16, UInt16, Boolean)](../html/5deaa15f-bc22-08b9-dffe-dd8b16c8be8c.htm "WriteIec 方法 (Byte, UInt16, UInt16, Boolean)")

[WriteIec 方法 (Byte, UInt16, UInt16, Byte[])](../html/fe546efd-9483-5c5b-2a89-e9c433ae3ed4.htm "WriteIec 方法 (Byte, UInt16, UInt16, Byte[])")

[WriteIec 方法 (Byte, UInt16, UInt16, Int16)](../html/fa5408d3-9794-f5db-0ac9-f700d71a2259.htm "WriteIec 方法 (Byte, UInt16, UInt16, Int16)")

[WriteIec 方法 (Byte, UInt16, UInt16, Single)](../html/207e2965-2c89-32e9-cb00-934333cfdc5a.htm "WriteIec 方法 (Byte, UInt16, UInt16, Single)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IEC104WriteIec 方法 (Byte, UInt16, UInt16, Byte) |

指定类型标识，传送原因，数据地址，数据值信息来写入到IEC的仪表中去

**命名空间：**
 [HslCommunication.Instrument.IEC](4f75643b-b529-650f-8587-78378826eec2.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public OperateResult WriteIec(
	byte type,
	ushort reason,
	ushort address,
	byte[] value
)
```

```
Public Function WriteIec ( 
	type As Byte,
	reason As UShort,
	address As UShort,
	value As Byte()
) As OperateResult
```

```
public:
OperateResult^ WriteIec(
	unsigned char type, 
	unsigned short reason, 
	unsigned short address, 
	array<unsigned char>^ value
)
```

```
member WriteIec : 
        type : byte * 
        reason : uint16 * 
        address : uint16 * 
        value : byte[] -> OperateResult 
```

#### 参数

type
:   类型：SystemByte  
    类型标识

reason
:   类型：SystemUInt16  
    传送原因

address
:   类型：SystemUInt16  
    信息对象地址

value
:   类型：SystemByte  
    值信息

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否发送到仪表成功

![](../icons/SectionExpanded.png)参见

#### 引用

[IEC104 类](8695140f-80b4-71fc-7a09-eff994df68ca.htm)

[WriteIec 重载](6ad1a56e-d20c-5001-6574-3d159b60f874.htm)

[HslCommunication.Instrument.IEC 命名空间](4f75643b-b529-650f-8587-78378826eec2.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## WriteIec 方法 (Byte, UInt16, UInt16, Int16)

[原文連結](http://api.hslcommunication.cn/html/fa5408d3-9794-f5db-0ac9-f700d71a2259.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.IEC](../html/4f75643b-b529-650f-8587-78378826eec2.htm "HslCommunication.Instrument.IEC")

[IEC104 类](../html/8695140f-80b4-71fc-7a09-eff994df68ca.htm "IEC104 类")

[IEC104 方法](../html/8eb8a1e0-92fb-2fdb-83ac-2cab7f59d96b.htm "IEC104 方法")

[WriteIec 方法](../html/6ad1a56e-d20c-5001-6574-3d159b60f874.htm "WriteIec 方法 ")

[WriteIec 方法 (Byte, UInt16, UInt16, Boolean)](../html/5deaa15f-bc22-08b9-dffe-dd8b16c8be8c.htm "WriteIec 方法 (Byte, UInt16, UInt16, Boolean)")

[WriteIec 方法 (Byte, UInt16, UInt16, Byte[])](../html/fe546efd-9483-5c5b-2a89-e9c433ae3ed4.htm "WriteIec 方法 (Byte, UInt16, UInt16, Byte[])")

[WriteIec 方法 (Byte, UInt16, UInt16, Int16)](../html/fa5408d3-9794-f5db-0ac9-f700d71a2259.htm "WriteIec 方法 (Byte, UInt16, UInt16, Int16)")

[WriteIec 方法 (Byte, UInt16, UInt16, Single)](../html/207e2965-2c89-32e9-cb00-934333cfdc5a.htm "WriteIec 方法 (Byte, UInt16, UInt16, Single)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IEC104WriteIec 方法 (Byte, UInt16, UInt16, Int16) |

指定类型标识，传送原因，数据地址，数据值信息来写入到IEC的仪表中去

**命名空间：**
 [HslCommunication.Instrument.IEC](4f75643b-b529-650f-8587-78378826eec2.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public OperateResult WriteIec(
	byte type,
	ushort reason,
	ushort address,
	short value
)
```

```
Public Function WriteIec ( 
	type As Byte,
	reason As UShort,
	address As UShort,
	value As Short
) As OperateResult
```

```
public:
OperateResult^ WriteIec(
	unsigned char type, 
	unsigned short reason, 
	unsigned short address, 
	short value
)
```

```
member WriteIec : 
        type : byte * 
        reason : uint16 * 
        address : uint16 * 
        value : int16 -> OperateResult 
```

#### 参数

type
:   类型：SystemByte  
    类型标识

reason
:   类型：SystemUInt16  
    传送原因

address
:   类型：SystemUInt16  
    信息对象地址

value
:   类型：SystemInt16  
    值信息

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否发送到仪表成功

![](../icons/SectionExpanded.png)参见

#### 引用

[IEC104 类](8695140f-80b4-71fc-7a09-eff994df68ca.htm)

[WriteIec 重载](6ad1a56e-d20c-5001-6574-3d159b60f874.htm)

[HslCommunication.Instrument.IEC 命名空间](4f75643b-b529-650f-8587-78378826eec2.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## WriteIec 方法 (Byte, UInt16, UInt16, Single)

[原文連結](http://api.hslcommunication.cn/html/207e2965-2c89-32e9-cb00-934333cfdc5a.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.IEC](../html/4f75643b-b529-650f-8587-78378826eec2.htm "HslCommunication.Instrument.IEC")

[IEC104 类](../html/8695140f-80b4-71fc-7a09-eff994df68ca.htm "IEC104 类")

[IEC104 方法](../html/8eb8a1e0-92fb-2fdb-83ac-2cab7f59d96b.htm "IEC104 方法")

[WriteIec 方法](../html/6ad1a56e-d20c-5001-6574-3d159b60f874.htm "WriteIec 方法 ")

[WriteIec 方法 (Byte, UInt16, UInt16, Boolean)](../html/5deaa15f-bc22-08b9-dffe-dd8b16c8be8c.htm "WriteIec 方法 (Byte, UInt16, UInt16, Boolean)")

[WriteIec 方法 (Byte, UInt16, UInt16, Byte[])](../html/fe546efd-9483-5c5b-2a89-e9c433ae3ed4.htm "WriteIec 方法 (Byte, UInt16, UInt16, Byte[])")

[WriteIec 方法 (Byte, UInt16, UInt16, Int16)](../html/fa5408d3-9794-f5db-0ac9-f700d71a2259.htm "WriteIec 方法 (Byte, UInt16, UInt16, Int16)")

[WriteIec 方法 (Byte, UInt16, UInt16, Single)](../html/207e2965-2c89-32e9-cb00-934333cfdc5a.htm "WriteIec 方法 (Byte, UInt16, UInt16, Single)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IEC104WriteIec 方法 (Byte, UInt16, UInt16, Single) |

指定类型标识，传送原因，数据地址，数据值信息来写入到IEC的仪表中去

**命名空间：**
 [HslCommunication.Instrument.IEC](4f75643b-b529-650f-8587-78378826eec2.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public OperateResult WriteIec(
	byte type,
	ushort reason,
	ushort address,
	float value
)
```

```
Public Function WriteIec ( 
	type As Byte,
	reason As UShort,
	address As UShort,
	value As Single
) As OperateResult
```

```
public:
OperateResult^ WriteIec(
	unsigned char type, 
	unsigned short reason, 
	unsigned short address, 
	float value
)
```

```
member WriteIec : 
        type : byte * 
        reason : uint16 * 
        address : uint16 * 
        value : float32 -> OperateResult 
```

#### 参数

type
:   类型：SystemByte  
    类型标识

reason
:   类型：SystemUInt16  
    传送原因

address
:   类型：SystemUInt16  
    信息对象地址

value
:   类型：SystemSingle  
    值信息

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否发送到仪表成功

![](../icons/SectionExpanded.png)参见

#### 引用

[IEC104 类](8695140f-80b4-71fc-7a09-eff994df68ca.htm)

[WriteIec 重载](6ad1a56e-d20c-5001-6574-3d159b60f874.htm)

[HslCommunication.Instrument.IEC 命名空间](4f75643b-b529-650f-8587-78378826eec2.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## IEC104 事件

[原文連結](http://api.hslcommunication.cn/html/98e01792-5618-ac77-38c5-20d50d856869.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.IEC](../html/4f75643b-b529-650f-8587-78378826eec2.htm "HslCommunication.Instrument.IEC")

[IEC104 类](../html/8695140f-80b4-71fc-7a09-eff994df68ca.htm "IEC104 类")

[IEC104 事件](../html/98e01792-5618-ac77-38c5-20d50d856869.htm "IEC104 事件")

[OnIEC104MessageReceived 事件](../html/d5456fa4-f675-29ed-96d7-3026316bc386.htm "OnIEC104MessageReceived 事件")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IEC104 事件 |

[IEC104](8695140f-80b4-71fc-7a09-eff994df68ca.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)事件

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共事件 | [OnIEC104MessageReceived](d5456fa4-f675-29ed-96d7-3026316bc386.htm) | 当接收到了IEC104的消息触发的事件 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[IEC104 类](8695140f-80b4-71fc-7a09-eff994df68ca.htm)

[HslCommunication.Instrument.IEC 命名空间](4f75643b-b529-650f-8587-78378826eec2.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## OnIEC104MessageReceived 事件

[原文連結](http://api.hslcommunication.cn/html/d5456fa4-f675-29ed-96d7-3026316bc386.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.IEC](../html/4f75643b-b529-650f-8587-78378826eec2.htm "HslCommunication.Instrument.IEC")

[IEC104 类](../html/8695140f-80b4-71fc-7a09-eff994df68ca.htm "IEC104 类")

[IEC104 事件](../html/98e01792-5618-ac77-38c5-20d50d856869.htm "IEC104 事件")

[OnIEC104MessageReceived 事件](../html/d5456fa4-f675-29ed-96d7-3026316bc386.htm "OnIEC104MessageReceived 事件")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IEC104OnIEC104MessageReceived 事件 |

当接收到了IEC104的消息触发的事件

**命名空间：**
 [HslCommunication.Instrument.IEC](4f75643b-b529-650f-8587-78378826eec2.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public event EventHandler<IEC104MessageEventArgs> OnIEC104MessageReceived
```

```
Public Event OnIEC104MessageReceived As EventHandler(Of IEC104MessageEventArgs)
```

```
public:
 event EventHandler<IEC104MessageEventArgs^>^ OnIEC104MessageReceived {
	void add (EventHandler<IEC104MessageEventArgs^>^ value);
	void remove (EventHandler<IEC104MessageEventArgs^>^ value);
}
```

```
member OnIEC104MessageReceived : IEvent<EventHandler<IEC104MessageEventArgs>,
    IEC104MessageEventArgs>
```

#### 值

类型：SystemEventHandler[IEC104MessageEventArgs](34f93396-e777-2b46-1db2-20c3f8ec3d23.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[IEC104 类](8695140f-80b4-71fc-7a09-eff994df68ca.htm)

[HslCommunication.Instrument.IEC 命名空间](4f75643b-b529-650f-8587-78378826eec2.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## IEC104 字段

[原文連結](http://api.hslcommunication.cn/html/f6e20d95-3950-c9c9-a87e-21f561d422d0.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.IEC](../html/4f75643b-b529-650f-8587-78378826eec2.htm "HslCommunication.Instrument.IEC")

[IEC104 类](../html/8695140f-80b4-71fc-7a09-eff994df68ca.htm "IEC104 类")

[IEC104 构造函数](../html/22a80517-47fd-8960-d826-72c0a07e021f.htm "IEC104 构造函数 ")

[IEC104 属性](../html/1da629b7-909c-dfb0-74a8-3a845a94fe2a.htm "IEC104 属性")

[IEC104 方法](../html/8eb8a1e0-92fb-2fdb-83ac-2cab7f59d96b.htm "IEC104 方法")

[IEC104 事件](../html/98e01792-5618-ac77-38c5-20d50d856869.htm "IEC104 事件")

[IEC104 字段](../html/f6e20d95-3950-c9c9-a87e-21f561d422d0.htm "IEC104 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IEC104 字段 |

[IEC104](8695140f-80b4-71fc-7a09-eff994df68ca.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)字段

|  | 名称 | 说明 |
| --- | --- | --- |
| 受保护的字段 | [LogMsgFormatBinary](2ed2ff87-caec-0e44-f08b-e2b015225a58.htm) | 设置日志记录报文是否二进制，如果为False，那就使用ASCII码  Set whether the log message is binary, if it is False, then use ASCII code (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[IEC104 类](8695140f-80b4-71fc-7a09-eff994df68ca.htm)

[HslCommunication.Instrument.IEC 命名空间](4f75643b-b529-650f-8587-78378826eec2.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## IEC104.TypeID 类

[原文連結](http://api.hslcommunication.cn/html/634a3bbf-dbf7-b538-6361-94ebfde09a1c.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.IEC](../html/4f75643b-b529-650f-8587-78378826eec2.htm "HslCommunication.Instrument.IEC")

[IEC104.TypeID 类](../html/634a3bbf-dbf7-b538-6361-94ebfde09a1c.htm "IEC104.TypeID 类")

[IEC104.TypeID 构造函数](../html/d8da4291-c3ef-6fe0-53ec-6c0cd37f91f7.htm "IEC104.TypeID 构造函数 ")

[TypeID 方法](../html/a42ebce3-1b83-30fc-3e4e-973e6c6d4f42.htm "TypeID 方法")

[TypeID 字段](../html/9a396918-532b-ba2e-0fc0-92dd37439c8c.htm "TypeID 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IEC104TypeID 类 |

类型信息资源

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  HslCommunication.Instrument.IECIEC104TypeID

**命名空间：**
 [HslCommunication.Instrument.IEC](4f75643b-b529-650f-8587-78378826eec2.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class TypeID
```

```
Public Class TypeID
```

```
public ref class TypeID
```

```
type TypeID =  class end
```

IEC104TypeID 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [IEC104TypeID](d8da4291-c3ef-6fe0-53ec-6c0cd37f91f7.htm) | 初始化 IEC104TypeID 类的一个新实例 |

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

![](../icons/SectionExpanded.png)字段

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共字段静态成员 | [C\_CI\_NA\_1](0f6a8a93-8722-e529-1854-7cbc0f50c7bf.htm) | 累计量召唤，带不同的限定词可以用于组召唤 |
| 公共字段静态成员 | [C\_CS\_NA\_1](a648a447-7390-5499-d6b4-876fd78be04b.htm) | 时钟同步命令，需要通过测量通道延迟加以校正 |
| 公共字段静态成员 | [C\_DC\_NA\_1](18f0c138-6edf-e9ef-e3d6-1794449466af.htm) | 双点遥控，一个报文只有一个遥控信息体，不带时标 |
| 公共字段静态成员 | [C\_IC\_NA\_1](e8af3d26-3b4f-3897-bc9a-f95beff9426d.htm) | 总召唤，带不同的限定词可以用于组召唤 |
| 公共字段静态成员 | [C\_RC\_NA\_1](0ad191e6-3e7a-9170-1d16-d1686a4ffec5.htm) | 升降遥控，一个报文只有一个遥控信息体，不带时标 |
| 公共字段静态成员 | [C\_RD\_NA\_1](422d8fe3-f0ce-490f-cc5c-8e9de3654c30.htm) | 读命令，读取单个的信息对象值 |
| 公共字段静态成员 | [C\_RS\_NA\_1](369f4de9-3821-db05-5c0a-a781e7398c39.htm) | 复位进程命令，使用前需要双方验证 |
| 公共字段静态成员 | [C\_SC\_NA\_1](0417a9a8-d5f4-db05-63e9-bae048fd3238.htm) | 单点遥控，一个报文只有一个遥控信息体，不带时标 |
| 公共字段静态成员 | [C\_SE\_NA\_1](e04304b8-dc6f-9fcc-efc2-72f604443ab5.htm) | 归一化设定值，一个报文只有一个设定值，不带时标 |
| 公共字段静态成员 | [C\_SE\_NB\_1](3ccbaf22-79ae-5520-77c9-1971aafed9d3.htm) | 标度化设定值，一个报文只有一个设定值，不带时标 |
| 公共字段静态成员 | [C\_SE\_NC\_1](ded43c3b-8fa4-f9f6-616f-4a3496fc0d8a.htm) | 短浮点设定值，一个报文只有一个设定值，不带时标 |
| 公共字段静态成员 | [C\_SE\_ND\_1](cc60ae55-8951-098f-4b3a-74188453a191.htm) | 32比特串设定，一个报文只有一个设定值，不带时标 |
| 公共字段静态成员 | [C\_SE\_NE\_1](c64d5ce4-800e-cb54-c578-bee4b23d95ad.htm) | 归一化设定值，一个报文可以包含多个设定值，不带时标 |
| 公共字段静态成员 | [C\_SE\_TA\_1](a668486e-cfa3-6f81-f5ae-c27bea84e8d9.htm) | 单点遥控，一个报文只有一个遥控信息体，带时标 |
| 公共字段静态成员 | [C\_SE\_TB\_1](b910b8b6-ca71-d377-5a83-45e5276b1302.htm) | 双点遥控，一个报文只有一个遥控信息体，带时标 |
| 公共字段静态成员 | [C\_SE\_TC\_1](3135cf57-4aca-b8af-a8b5-9d0722284323.htm) | 升降遥控，一个报文只有一个遥控信息体，带时标 |
| 公共字段静态成员 | [C\_SE\_TD\_1](977034f4-3850-4ca4-389f-1b5e1862dc65.htm) | 归一化设定值，一个报文只有一个设定值，带时标 |
| 公共字段静态成员 | [C\_SE\_TE\_1](40f16971-4365-ac10-bb6e-2256c8021755.htm) | 标度化设定值，一个报文只有一个设定值，带时标 |
| 公共字段静态成员 | [C\_SE\_TF\_1](8ade1932-71ba-5bae-7a4b-b1e7e6b6d2b1.htm) | 短浮点设定值，一个报文只有一个设定值，带时标 |
| 公共字段静态成员 | [C\_SE\_TG\_1](38572d98-012b-f8ec-5430-576f1ee874f7.htm) | 32比特串设定，一个报文只有一个设定值，带时标 |
| 公共字段静态成员 | [C\_TS\_TA\_1](2b236b3b-9ed0-3469-e0d5-2601bb1e7fb3.htm) | 带时标的测试命令 |
| 公共字段静态成员 | [M\_BO\_NA\_1](256662bb-34cb-6b24-8ac4-15a2efe84002.htm) | 32比特串，带品质描述，不带时标 |
| 公共字段静态成员 | [M\_BO\_TB\_1](539aee06-598b-96ac-b956-df0ff1afd531.htm) | 32比特串，带品质描述，带绝对时标 |
| 公共字段静态成员 | [M\_DP\_NA\_1](98ae8bfc-ebec-e05d-9817-f3ea93691e54.htm) | 双点遥信，带品质描述，不带时标 |
| 公共字段静态成员 | [M\_DP\_TB\_1](8d3dafe9-6d94-7aa1-b935-6ea5e47fd873.htm) | 双点遥信，带品质描述，带绝对时标 |
| 公共字段静态成员 | [M\_EI\_NA\_1](10b98b81-8319-c76f-15a8-eb35fad5cea9.htm) | 初始化结束，报告厂站初始化完成 |
| 公共字段静态成员 | [M\_IT\_NA\_1](8f6d48b0-ba84-6af7-5544-4128575e77b3.htm) | 累计量，带品质描述，不带时标 |
| 公共字段静态成员 | [M\_IT\_TB\_1](4bb81633-fcc7-5eb8-1cd7-8cebaef5e590.htm) | 累计量，带品质描述，带绝对时标 |
| 公共字段静态成员 | [M\_ME\_NA\_1](3669d8e0-2520-f235-c753-f06469961f74.htm) | 归一化遥测值，带品质描述，不带时标 |
| 公共字段静态成员 | [M\_ME\_NB\_1](1b0de6ce-e9bc-8870-6dc0-77c7c60d3ba4.htm) | 标度化遥测值，带品质描述，不带时标 |
| 公共字段静态成员 | [M\_ME\_NC\_1](9b65d075-ee7f-17ee-c138-88310bacfa26.htm) | 短浮点遥测值，带品质描述，不带时标 |
| 公共字段静态成员 | [M\_ME\_ND\_1](60890c4c-bcf1-2343-2c5e-2d3f95dee93f.htm) | 归一化遥测值，不带品质描述，不带时标 |
| 公共字段静态成员 | [M\_ME\_TD\_1](69aee490-3bc5-9e4a-81ff-05d19913da48.htm) | 归一化遥测值，带品质描述，带绝对时标 |
| 公共字段静态成员 | [M\_ME\_TE\_1](0b4e705d-6ffe-34c4-bcb1-c07a0964ae93.htm) | 标度化遥测值，带品质描述，带绝对时标 |
| 公共字段静态成员 | [M\_ME\_TF\_1](4338d9a7-bafa-3c04-7cc8-9cb4b2d4aff6.htm) | 短浮点遥测值，带品质描述，带绝对时标 |
| 公共字段静态成员 | [M\_PS\_NA\_1](ad515fc9-bff7-7171-dfb9-77b2906e0c2e.htm) | 成组单点遥信，只带变位标志 |
| 公共字段静态成员 | [M\_SP\_NA\_1](9e5f226b-86c9-5b95-603b-1c679e7720d1.htm) | 单点遥信，带品质描述，不带时标 |
| 公共字段静态成员 | [M\_SP\_TB\_1](8e372e65-d66f-0842-3a4f-f58ab60625b4.htm) | 单点遥信，带品质描述，带绝对时标 |
| 公共字段静态成员 | [M\_ST\_NA\_1](2a268ffe-e18b-9f7d-bf92-f846eee79dc1.htm) | 步位置信息，带品质描述，不带时标 |
| 公共字段静态成员 | [M\_ST\_TB\_1](f2f53c1d-7fa4-eee2-8f60-32d2d64de18f.htm) | 步位置信息，带品质描述，带绝对时标 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.Instrument.IEC 命名空间](4f75643b-b529-650f-8587-78378826eec2.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)