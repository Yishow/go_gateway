# HslCommunication - HslCommunication.Serial

> 分類頁數: 30



---
## HslCommunication.Serial

[原文連結](http://api.hslcommunication.cn/html/8d02123a-6db7-7a98-75db-8e3b63e89789.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Serial](../html/8d02123a-6db7-7a98-75db-8e3b63e89789.htm "HslCommunication.Serial")

[CheckType 枚举](../html/4aa8e7b7-fed9-5aa0-66b9-14e7c677fc3e.htm "CheckType 枚举")

[SerialBase 类](../html/e877e20b-ef65-266c-05a5-5db58e6bad0c.htm "SerialBase 类")

[SoftCRC16 类](../html/b01f7f5f-55ee-cbf8-faf3-bb5ca2d26e57.htm "SoftCRC16 类")

[SoftLRC 类](../html/7fa2a303-52da-0d46-ddd1-94c14f91f3bd.htm "SoftLRC 类")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslCommunication.Serial 命名空间 |

[缺少 "N:HslCommunication.Serial" 的 <summary> 文档]

![](../icons/SectionExpanded.png)类

|  | 类 | 说明 |
| --- | --- | --- |
| 公共类 | [SerialBase](e877e20b-ef65-266c-05a5-5db58e6bad0c.htm) | 所有串行通信类的基类，提供了一些基础的服务，核心的通信实现  The base class of all serial communication classes provides some basic services for the core communication implementation |
| 公共类代码示例 | [SoftCRC16](b01f7f5f-55ee-cbf8-faf3-bb5ca2d26e57.htm) | 用于CRC16验证的类，提供了标准的验证方法，可以方便快速的对数据进行CRC校验  The class for CRC16 validation provides a standard validation method that makes it easy to CRC data quickly |
| 公共类 | [SoftLRC](7fa2a303-52da-0d46-ddd1-94c14f91f3bd.htm) | 用于LRC验证的类，提供了标准的验证方法  The class used for LRC verification provides a standard verification method |

![](../icons/SectionExpanded.png)枚举

|  | 枚举 | 说明 |
| --- | --- | --- |
| 公共枚举 | [CheckType](4aa8e7b7-fed9-5aa0-66b9-14e7c677fc3e.htm) | 校验方式 |

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## CheckType 枚举

[原文連結](http://api.hslcommunication.cn/html/4aa8e7b7-fed9-5aa0-66b9-14e7c677fc3e.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Serial](../html/8d02123a-6db7-7a98-75db-8e3b63e89789.htm "HslCommunication.Serial")

[CheckType 枚举](../html/4aa8e7b7-fed9-5aa0-66b9-14e7c677fc3e.htm "CheckType 枚举")

[SerialBase 类](../html/e877e20b-ef65-266c-05a5-5db58e6bad0c.htm "SerialBase 类")

[SoftCRC16 类](../html/b01f7f5f-55ee-cbf8-faf3-bb5ca2d26e57.htm "SoftCRC16 类")

[SoftLRC 类](../html/7fa2a303-52da-0d46-ddd1-94c14f91f3bd.htm "SoftLRC 类")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CheckType 枚举 |

校验方式

**命名空间：**
 [HslCommunication.Serial](8d02123a-6db7-7a98-75db-8e3b63e89789.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public enum CheckType
```

```
Public Enumeration CheckType
```

```
public enum class CheckType
```

```
type CheckType
```

![](../icons/SectionExpanded.png)成员

|  | 成员名称 | 值 | 说明 |
| --- | --- | --- | --- |
|  | BCC | 0 | 和校验 |
|  | CRC16 | 1 | CRC校验的方式 |

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.Serial 命名空间](8d02123a-6db7-7a98-75db-8e3b63e89789.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## SerialBase 类

[原文連結](http://api.hslcommunication.cn/html/e877e20b-ef65-266c-05a5-5db58e6bad0c.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Serial](../html/8d02123a-6db7-7a98-75db-8e3b63e89789.htm "HslCommunication.Serial")

[SerialBase 类](../html/e877e20b-ef65-266c-05a5-5db58e6bad0c.htm "SerialBase 类")

[SerialBase 构造函数](../html/c67cfb2a-45e7-d93d-8c47-8ac0aa342052.htm "SerialBase 构造函数 ")

[SerialBase 属性](../html/e82084a0-6e96-85cd-f4dc-0aa932aa0ae4.htm "SerialBase 属性")

[SerialBase 方法](../html/e28d440f-d77d-55fd-4869-070495fa1a3d.htm "SerialBase 方法")

[SerialBase 字段](../html/f6820884-a072-a4e4-125e-4d91fc28953e.htm "SerialBase 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| SerialBase 类 |

所有串行通信类的基类，提供了一些基础的服务，核心的通信实现  
The base class of all serial communication classes provides some basic services for the core communication implementation

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  [HslCommunication.Core.NetBinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)  
    HslCommunication.SerialSerialBase  
      [HslCommunication.Instrument.LightShineInLightSourceController](5de4a57b-1202-4bcd-4650-a524ba7166df.htm)  
      [HslCommunication.Profinet.IDCardSAMSerial](48f6a433-86f0-a4f8-1751-462e12ad6923.htm)  
      [HslCommunication.Profinet.KeyenceKeyenceSR2000Serial](77f8d8da-8005-e4ca-598a-7ea8c5e9d40c.htm)

**命名空间：**
 [HslCommunication.Serial](8d02123a-6db7-7a98-75db-8e3b63e89789.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class SerialBase : BinaryCommunication, IDisposable
```

```
Public Class SerialBase
	Inherits BinaryCommunication
	Implements IDisposable
```

```
public ref class SerialBase : public BinaryCommunication, 
	IDisposable
```

```
type SerialBase =  
    class
        inherit BinaryCommunication
        interface IDisposable
    end
```

SerialBase 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [SerialBase](c67cfb2a-45e7-d93d-8c47-8ac0aa342052.htm) | 实例化一个无参的构造方法  Instantiate a parameterless constructor |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [BaudRate](1c33c44a-d2b9-badd-df76-09e155d26cc1.htm) | 当前连接串口信息的波特率  Baud rate of current connection serial port information |
| 公共属性 | [CommunicationPipe](cb528f9b-fe2c-5ac7-6912-2c3c5c1fdf2a.htm) | 获取或设置当前的管道信息，管道类型为[CommunicationPipe](ad8e903b-d289-21b0-077e-15ede8005221.htm)的继承类，内置了[PipeTcpNet](26b94bd5-2695-36e7-6378-c3d41846be8f.htm)管道，[PipeUdpNet](5a0b810e-187e-826b-0a13-34fd5576aa0e.htm)管道，[PipeSerialPort](9d924607-3dfc-3916-b532-c5dec48fe347.htm)管道等  Get or set the current pipeline information, the pipeline type is [CommunicationPipe](ad8e903b-d289-21b0-077e-15ede8005221.htm) inheritance class, [PipeTcpNet](26b94bd5-2695-36e7-6378-c3d41846be8f.htm) pipeline, [PipeUdpNet](5a0b810e-187e-826b-0a13-34fd5576aa0e.htm) pipeline, [PipeSerialPort](9d924607-3dfc-3916-b532-c5dec48fe347.htm) pipeline, etc (重写 [BinaryCommunicationCommunicationPipe](ad8e903b-d289-21b0-077e-15ede8005221.htm).) |
| 公共属性 | [ConnectionId](ce1efb4b-fbc9-4683-163e-0a8e71e99d6b.htm) | 当前连接的唯一ID号，默认为长度20的guid码加随机数组成，方便列表管理，也可以自己指定  The unique ID number of the current connection. The default is a 20-digit guid code plus a random number. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [IsClearCacheBeforeRead](cf2aaeff-e4eb-849a-4783-d10e5e965515.htm) | 是否在发送数据前清空缓冲数据，默认是false  Whether to empty the buffer before sending data, the default is false |
| 公共属性代码示例 | [LogNet](8c56f246-88ac-afb8-8c06-972fafe4f78a.htm) | 组件的日志工具，支持日志记录，只要实例化后，当前网络的基本信息，就以[DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm)等级进行输出  The component's logging tool supports logging. As long as the instantiation of the basic network information, the output will be output at [DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm) (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [PortName](250b8746-0bb7-4caf-2a43-11109f2bda20.htm) | 当前连接串口信息的端口号名称  The port name of the current connection serial port information |
| 公共属性 | [ReceiveEmptyDataCount](480b2206-1de4-6a74-6f8e-7c30a2ca05f9.htm) | 获取或设置连续接收空的数据次数，在数据接收完成时有效，每个单位消耗的时间为[SleepTime](1e3a83ef-b6fc-bc8c-8950-ea3f6ae5c19b.htm)。  Obtain or set the number of consecutive times to receive empty data, which is valid when the data is received, and the time consumed by each unit is [SleepTime](1e3a83ef-b6fc-bc8c-8950-ea3f6ae5c19b.htm) |
| 公共属性代码示例 | [ReceiveTimeOut](250560a7-332b-ce1b-5eb0-e3abe4ee6012.htm) | 获取或设置接收服务器反馈的时间，如果为负数，则不接收反馈   Gets or sets the time to receive server feedback, and if it is a negative number, does not receive feedback (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [RtsEnable](18c16d98-49b0-5ae2-5c2b-0c6862801c92.htm) | 获取或设置一个值，该值指示在串行通信中是否启用请求发送 (RTS) 信号。  Gets or sets a value indicating whether the request sending (RTS) signal is enabled in serial communication. |
| 公共属性 | [SendBeforeHex](72aa1de4-116f-abfe-4adb-814a0be0e560.htm) | 获取或设置在发送通信报文前追加发送的字节信息，HEX格式，通常用于lora组网时，需要携带 00 00 00 02 四个字节的站地址功能。  Obtain or set the byte information sent before sending communication packets, HEX format, usually used for LORA networking, you need to carry 00 00 00 02 four-byte station address function. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [SleepTime](86a5fbb9-3456-6beb-7113-5bbad5eabbc5.htm) | 获取或设置在正式接收对方返回数据前的时候，需要休息的时间，当设置为0的时候，不需要休息。  Get or set the time required to rest before officially receiving the data from the other party. When it is set to 0, no rest is required. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [Close](992f28b3-b99b-2196-274e-2f25e9bd3c1d.htm) | 关闭当前的串口连接  Close the current serial connection |
| 受保护的方法 | [DecideWhetherQAMessage](acc40cda-e310-8537-6685-59b67eeb16ff.htm) | 决定当前的消息是否是用于问答机制返回的消息，默认直接返回 true, 实际的情况需要根据协议进行重写方法  To determine whether the current message is the message returned by the question answering mechanism, the default is true. In actual cases, the rewriting method needs to be performed according to the protocol (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [Dispose](c83e3a27-3a09-f67c-9c7f-96b7a5a4c8b0.htm) | 释放当前的对象 |
| 受保护的方法 | [Dispose(Boolean)](c509106d-35a7-2ff2-7dd8-ca67e16d7950.htm) | 释放当前的对象 |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | [ExtraAfterReadFromCoreServer](06a49a15-a3ce-5e39-28c7-2b78851eca2e.htm) | 和服务器交互完成的时候调用的方法，可以根据读写结果进行一些额外的操作，具体的操作需要根据实际的需求来重写实现  The method called when the interaction with the server is completed can perform some additional operations based on the read and write results. The specific operations need to be rewritten according to actual needs. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法代码示例 | [ExtraOnDisconnect](de9dbba2-3b2c-1c5e-53c8-4e6f8a44e2c4.htm) | 根据实际的协议选择是否重写本方法，有些协议在断开连接之前，需要发送一些报文来关闭当前的网络通道  Select whether to rewrite this method according to the actual protocol. Some protocols need to send some packets to close the current network channel before disconnecting. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法代码示例 | [ExtraOnDisconnectAsync](210f593e-16d8-77c4-1178-58fbcfe9b27b.htm) | 根据实际的协议选择是否重写本方法，有些协议在断开连接之前，需要发送一些报文来关闭当前的网络通道  Select whether to rewrite this method according to the actual protocol. Some protocols need to send some packets to close the current network channel before disconnecting. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 受保护的方法 | [GetLogTextFromBinary](67f3b2f7-957c-fa0b-8320-6799237157db.htm) | 获取当前的报文进行日志记录的时候，是否使用二进制的格式记录，默认返回 [LogMsgFormatBinary](2ed2ff87-caec-0e44-f08b-e2b015225a58.htm)，重写可以根据session对象分别返回不同记录模式  Whether to log the current packet in binary format, the default return is [LogMsgFormatBinary](2ed2ff87-caec-0e44-f08b-e2b015225a58.htm). If you want to override it, different recording modes can be returned according to session (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [GetNewNetMessage](16636a1b-a6f0-9c0f-9d9c-c94e26c67459.htm) | 获取一个新的消息对象的方法，需要在继承类里面进行重写  The method to get a new message object needs to be overridden in the inheritance class (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法代码示例 | [InitializationOnConnect](4ab6989f-a5b5-05fe-8b44-bfa45824a2a7.htm) | 根据实际的协议选择是否重写本方法，有些协议在创建连接之后，需要进行一些初始化的信号握手，才能最终建立网络通道。  Whether to rewrite this method is based on the actual protocol. Some protocols require some initial signal handshake to establish a network channel after the connection is created. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法代码示例 | [InitializationOnConnectAsync](7da2af04-d4fc-5045-19a4-4b0d1722eef3.htm) | 根据实际的协议选择是否重写本方法，有些协议在创建连接之后，需要进行一些初始化的信号握手，才能最终建立网络通道。  Whether to rewrite this method is based on the actual protocol. Some protocols require some initial signal handshake to establish a network channel after the connection is created. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [IsOpen](24a0e05a-dd59-2e3c-e4c5-d3bb6ebc8881.htm) | 获取一个值，指示串口是否处于打开状态  Gets a value indicating whether the serial port is open |
| 受保护的方法 | [LogRevcMessage(Byte)](dab1657f-1e75-7fe1-7003-ccedfc69855d.htm) | 使用日志记录一个接收的报文信息  Logs are used to record information about a received packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [LogRevcMessage(Byte, PipeSession)](0f51ae45-575e-4228-f978-01b37d6d5df1.htm) | 使用日志记录一个接收的报文信息  Logs are used to record information about a received packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [LogSendMessage(Byte)](4f90eecc-68a2-1ef2-44ef-369f19bed3ec.htm) | 使用日志记录一个发送的报文信息  Logs are used to record information about a send packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [LogSendMessage(Byte, PipeSession)](209caadf-4f0f-40aa-04c2-ebc61f8dab13.htm) | 使用日志记录一个发送的报文信息  Logs are used to record information about a send packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法 | [Open](3501fe5b-65b5-a973-5f25-3de45d0054d6.htm) | 打开当前的管道信息，返回是否成功打开的结果对象，并通过属性 "Content" 指示当前是否为新创建的连接对象，如果是，则该值为 true  Opens the current pipe information, returns whether the result object was successfully opened, and indicates whether the current connection object is a newly created connection object through the property "Content", if so, the value is true |
| 公共方法 | [PackCommandWithHeader](d0b2cd76-c306-11fc-38c7-af8ad306cce7.htm) | 对当前的命令进行打包处理，通常是携带命令头内容，标记当前的命令的长度信息，需要进行重写，否则默认不打包  The current command is packaged, usually carrying the content of the command header, marking the length of the current command, and it needs to be rewritten, otherwise it is not packaged by default (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServer(Byte)](8146dc0c-1edc-3ca7-428a-a3bdbd35e6c5.htm) | 将二进制的数据发送到管道中去，然后从管道里接收二进制的数据回来，并返回是否成功的结果对象。  Send binary data to the pipeline, and then receive binary data back from the pipeline, and return whether the success of the result object (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServer(IEnumerableByte)](044cfdd9-ec76-b89f-2f2b-5f439025006c.htm) | 将多个数据报文按顺序发到设备，并从设备接收返回的数据内容，然后拼接成一个Byte[]信息，需要重写[UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm)方法才能返回正确的结果。  Send multiple data packets to the device in sequence, and receive the returned data content from the device, and then splicing them into a Byte[] message, you need to rewrite [UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm) method to return the correct result. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServer(Byte, Boolean, Boolean)](b32ce8e8-a959-fcc2-6736-14b39fdb9a40.htm) | 将二进制的数据发送到管道中去，然后从管道里接收二进制的数据回来，并返回是否成功的结果对象。  Send binary data to the pipeline, and then receive binary data back from the pipeline, and return whether the success of the result object (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServer(CommunicationPipe, Byte, Boolean, Boolean)](ef36d56b-df09-f727-e8d3-316d37a13ee9.htm) | 使用指定的管道来进行数据通信，发送原始数据到管道，然后从管道接收相关的数据返回，本方法无锁 (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServerAsync(Byte)](a4de5a24-c67a-7b02-db05-33b85eaf6e1e.htm) | 将二进制的数据发送到管道中去，然后从管道里接收二进制的数据回来，并返回是否成功的结果对象。  Send binary data to the pipeline, and then receive binary data back from the pipeline, and return whether the success of the result object (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServerAsync(IEnumerableByte)](bcdb59b7-a70a-d5c0-1fc6-9ececcd39ea2.htm) | 将多个数据报文按顺序发到设备，并从设备接收返回的数据内容，然后拼接成一个Byte[]信息，需要重写[UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm)方法才能返回正确的结果。  Send multiple data packets to the device in sequence, and receive the returned data content from the device, and then splicing them into a Byte[] message, you need to rewrite [UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm) method to return the correct result. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServerAsync(Byte, Boolean, Boolean)](1c79c204-5ab5-9342-e65d-7455c20ea577.htm) | 将二进制的数据发送到管道中去，然后从管道里接收二进制的数据回来，并返回是否成功的结果对象。  Send binary data to the pipeline, and then receive binary data back from the pipeline, and return whether the success of the result object (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServerAsync(CommunicationPipe, Byte, Boolean, Boolean)](a254a19d-3c29-147b-335b-8adccdf4ee2f.htm) | 使用指定的管道来进行数据通信，发送原始数据到管道，然后从管道接收相关的数据返回，本方法无锁 (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [SerialPortInni(ActionSerialPort)](53585409-6243-9359-1846-bf8b4e266708.htm) | 根据自定义初始化方法进行初始化串口信息  Initialize the serial port information according to the custom initialization method |
| 公共方法 | [SerialPortInni(String)](7dff5e26-6594-291e-4792-b13db9c5d613.htm) | 初始化串口信息，9600波特率，8位数据位，1位停止位，无奇偶校验  Initial serial port information, 9600 baud rate, 8 data bits, 1 stop bit, no parity |
| 公共方法 | [SerialPortInni(String, Int32)](0c7aa163-1d81-54fd-f4f5-103a7e9a341d.htm) | 初始化串口信息，波特率，8位数据位，1位停止位，无奇偶校验  Initializes serial port information, baud rate, 8-bit data bit, 1-bit stop bit, no parity |
| 公共方法 | [SerialPortInni(String, Int32, Int32, StopBits, Parity)](420c818b-48e9-fa02-59e2-f4d5897d6802.htm) | 初始化串口信息，波特率，数据位，停止位，奇偶校验需要全部自己来指定  Start serial port information, baud rate, data bit, stop bit, parity all need to be specified |
| 公共方法 | [SetDtuPipe](f765d1d7-234a-b229-9457-97e25e024dc9.htm) | 将当前的通信对象设置DTU模式，允许传入现成的管道，并返回初始化结果，如果该设备重写了握手报文，就是返回握手结果  Set the current communication object to DTU mode, allow the existing pipe to be passed in, and return the initialization result, if the device rewrites the handshake packet, the handshake result is returned (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [SetDtuPipeAsync](9d1fd12d-b804-3828-a627-a2c0eef1d23f.htm) | 将当前的通信对象设置DTU模式，允许传入现成的管道，并返回初始化结果，如果该设备重写了握手报文，就是返回握手结果  Set the current communication object to DTU mode, allow the existing pipe to be passed in, and return the initialization result, if the device rewrites the handshake packet, the handshake result is returned (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ToString](1ac09324-4d19-863d-d683-cbe804ef2124.htm) | (重写 ObjectToString.) |
| 公共方法 | [UnpackResponseContent](208017a7-8da2-33f8-3273-d7e7e850100f.htm) | 根据对方返回的报文命令，对命令进行基本的拆包，例如各种Modbus协议拆包为统一的核心报文，还支持对报文的验证  According to the message command returned by the other party, the command is basically unpacked, for example, various Modbus protocols are unpacked into a unified core message, and the verification of the message is also supported (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |

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

[HslCommunication.Serial 命名空间](8d02123a-6db7-7a98-75db-8e3b63e89789.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## SerialBase 构造函数 

[原文連結](http://api.hslcommunication.cn/html/c67cfb2a-45e7-d93d-8c47-8ac0aa342052.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Serial](../html/8d02123a-6db7-7a98-75db-8e3b63e89789.htm "HslCommunication.Serial")

[SerialBase 类](../html/e877e20b-ef65-266c-05a5-5db58e6bad0c.htm "SerialBase 类")

[SerialBase 构造函数](../html/c67cfb2a-45e7-d93d-8c47-8ac0aa342052.htm "SerialBase 构造函数 ")

[SerialBase 属性](../html/e82084a0-6e96-85cd-f4dc-0aa932aa0ae4.htm "SerialBase 属性")

[SerialBase 方法](../html/e28d440f-d77d-55fd-4869-070495fa1a3d.htm "SerialBase 方法")

[SerialBase 字段](../html/f6820884-a072-a4e4-125e-4d91fc28953e.htm "SerialBase 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| SerialBase 构造函数 |

实例化一个无参的构造方法  
Instantiate a parameterless constructor

**命名空间：**
 [HslCommunication.Serial](8d02123a-6db7-7a98-75db-8e3b63e89789.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public SerialBase()
```

```
Public Sub New
```

```
public:
SerialBase()
```

```
new : unit -> SerialBase
```

![](../icons/SectionExpanded.png)参见

#### 引用

[SerialBase 类](e877e20b-ef65-266c-05a5-5db58e6bad0c.htm)

[HslCommunication.Serial 命名空间](8d02123a-6db7-7a98-75db-8e3b63e89789.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## SerialBase 属性

[原文連結](http://api.hslcommunication.cn/html/e82084a0-6e96-85cd-f4dc-0aa932aa0ae4.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Serial](../html/8d02123a-6db7-7a98-75db-8e3b63e89789.htm "HslCommunication.Serial")

[SerialBase 类](../html/e877e20b-ef65-266c-05a5-5db58e6bad0c.htm "SerialBase 类")

[SerialBase 属性](../html/e82084a0-6e96-85cd-f4dc-0aa932aa0ae4.htm "SerialBase 属性")

[BaudRate 属性](../html/1c33c44a-d2b9-badd-df76-09e155d26cc1.htm "BaudRate 属性 ")

[CommunicationPipe 属性](../html/cb528f9b-fe2c-5ac7-6912-2c3c5c1fdf2a.htm "CommunicationPipe 属性 ")

[IsClearCacheBeforeRead 属性](../html/cf2aaeff-e4eb-849a-4783-d10e5e965515.htm "IsClearCacheBeforeRead 属性 ")

[PortName 属性](../html/250b8746-0bb7-4caf-2a43-11109f2bda20.htm "PortName 属性 ")

[ReceiveEmptyDataCount 属性](../html/480b2206-1de4-6a74-6f8e-7c30a2ca05f9.htm "ReceiveEmptyDataCount 属性 ")

[RtsEnable 属性](../html/18c16d98-49b0-5ae2-5c2b-0c6862801c92.htm "RtsEnable 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| SerialBase 属性 |

[SerialBase](e877e20b-ef65-266c-05a5-5db58e6bad0c.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [BaudRate](1c33c44a-d2b9-badd-df76-09e155d26cc1.htm) | 当前连接串口信息的波特率  Baud rate of current connection serial port information |
| 公共属性 | [CommunicationPipe](cb528f9b-fe2c-5ac7-6912-2c3c5c1fdf2a.htm) | 获取或设置当前的管道信息，管道类型为[CommunicationPipe](ad8e903b-d289-21b0-077e-15ede8005221.htm)的继承类，内置了[PipeTcpNet](26b94bd5-2695-36e7-6378-c3d41846be8f.htm)管道，[PipeUdpNet](5a0b810e-187e-826b-0a13-34fd5576aa0e.htm)管道，[PipeSerialPort](9d924607-3dfc-3916-b532-c5dec48fe347.htm)管道等  Get or set the current pipeline information, the pipeline type is [CommunicationPipe](ad8e903b-d289-21b0-077e-15ede8005221.htm) inheritance class, [PipeTcpNet](26b94bd5-2695-36e7-6378-c3d41846be8f.htm) pipeline, [PipeUdpNet](5a0b810e-187e-826b-0a13-34fd5576aa0e.htm) pipeline, [PipeSerialPort](9d924607-3dfc-3916-b532-c5dec48fe347.htm) pipeline, etc (重写 [BinaryCommunicationCommunicationPipe](ad8e903b-d289-21b0-077e-15ede8005221.htm).) |
| 公共属性 | [ConnectionId](ce1efb4b-fbc9-4683-163e-0a8e71e99d6b.htm) | 当前连接的唯一ID号，默认为长度20的guid码加随机数组成，方便列表管理，也可以自己指定  The unique ID number of the current connection. The default is a 20-digit guid code plus a random number. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [IsClearCacheBeforeRead](cf2aaeff-e4eb-849a-4783-d10e5e965515.htm) | 是否在发送数据前清空缓冲数据，默认是false  Whether to empty the buffer before sending data, the default is false |
| 公共属性代码示例 | [LogNet](8c56f246-88ac-afb8-8c06-972fafe4f78a.htm) | 组件的日志工具，支持日志记录，只要实例化后，当前网络的基本信息，就以[DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm)等级进行输出  The component's logging tool supports logging. As long as the instantiation of the basic network information, the output will be output at [DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm) (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [PortName](250b8746-0bb7-4caf-2a43-11109f2bda20.htm) | 当前连接串口信息的端口号名称  The port name of the current connection serial port information |
| 公共属性 | [ReceiveEmptyDataCount](480b2206-1de4-6a74-6f8e-7c30a2ca05f9.htm) | 获取或设置连续接收空的数据次数，在数据接收完成时有效，每个单位消耗的时间为[SleepTime](1e3a83ef-b6fc-bc8c-8950-ea3f6ae5c19b.htm)。  Obtain or set the number of consecutive times to receive empty data, which is valid when the data is received, and the time consumed by each unit is [SleepTime](1e3a83ef-b6fc-bc8c-8950-ea3f6ae5c19b.htm) |
| 公共属性代码示例 | [ReceiveTimeOut](250560a7-332b-ce1b-5eb0-e3abe4ee6012.htm) | 获取或设置接收服务器反馈的时间，如果为负数，则不接收反馈   Gets or sets the time to receive server feedback, and if it is a negative number, does not receive feedback (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [RtsEnable](18c16d98-49b0-5ae2-5c2b-0c6862801c92.htm) | 获取或设置一个值，该值指示在串行通信中是否启用请求发送 (RTS) 信号。  Gets or sets a value indicating whether the request sending (RTS) signal is enabled in serial communication. |
| 公共属性 | [SendBeforeHex](72aa1de4-116f-abfe-4adb-814a0be0e560.htm) | 获取或设置在发送通信报文前追加发送的字节信息，HEX格式，通常用于lora组网时，需要携带 00 00 00 02 四个字节的站地址功能。  Obtain or set the byte information sent before sending communication packets, HEX format, usually used for LORA networking, you need to carry 00 00 00 02 four-byte station address function. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [SleepTime](86a5fbb9-3456-6beb-7113-5bbad5eabbc5.htm) | 获取或设置在正式接收对方返回数据前的时候，需要休息的时间，当设置为0的时候，不需要休息。  Get or set the time required to rest before officially receiving the data from the other party. When it is set to 0, no rest is required. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[SerialBase 类](e877e20b-ef65-266c-05a5-5db58e6bad0c.htm)

[HslCommunication.Serial 命名空间](8d02123a-6db7-7a98-75db-8e3b63e89789.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## BaudRate 属性 

[原文連結](http://api.hslcommunication.cn/html/1c33c44a-d2b9-badd-df76-09e155d26cc1.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Serial](../html/8d02123a-6db7-7a98-75db-8e3b63e89789.htm "HslCommunication.Serial")

[SerialBase 类](../html/e877e20b-ef65-266c-05a5-5db58e6bad0c.htm "SerialBase 类")

[SerialBase 属性](../html/e82084a0-6e96-85cd-f4dc-0aa932aa0ae4.htm "SerialBase 属性")

[BaudRate 属性](../html/1c33c44a-d2b9-badd-df76-09e155d26cc1.htm "BaudRate 属性 ")

[CommunicationPipe 属性](../html/cb528f9b-fe2c-5ac7-6912-2c3c5c1fdf2a.htm "CommunicationPipe 属性 ")

[IsClearCacheBeforeRead 属性](../html/cf2aaeff-e4eb-849a-4783-d10e5e965515.htm "IsClearCacheBeforeRead 属性 ")

[PortName 属性](../html/250b8746-0bb7-4caf-2a43-11109f2bda20.htm "PortName 属性 ")

[ReceiveEmptyDataCount 属性](../html/480b2206-1de4-6a74-6f8e-7c30a2ca05f9.htm "ReceiveEmptyDataCount 属性 ")

[RtsEnable 属性](../html/18c16d98-49b0-5ae2-5c2b-0c6862801c92.htm "RtsEnable 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| SerialBaseBaudRate 属性 |

当前连接串口信息的波特率  
Baud rate of current connection serial port information

**命名空间：**
 [HslCommunication.Serial](8d02123a-6db7-7a98-75db-8e3b63e89789.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public int BaudRate { get; }
```

```
Public ReadOnly Property BaudRate As Integer
	Get
```

```
public:
property int BaudRate {
	int get ();
}
```

```
member BaudRate : int with get
```

#### 属性值

类型：Int32

![](../icons/SectionExpanded.png)参见

#### 引用

[SerialBase 类](e877e20b-ef65-266c-05a5-5db58e6bad0c.htm)

[HslCommunication.Serial 命名空间](8d02123a-6db7-7a98-75db-8e3b63e89789.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## CommunicationPipe 属性 

[原文連結](http://api.hslcommunication.cn/html/cb528f9b-fe2c-5ac7-6912-2c3c5c1fdf2a.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Serial](../html/8d02123a-6db7-7a98-75db-8e3b63e89789.htm "HslCommunication.Serial")

[SerialBase 类](../html/e877e20b-ef65-266c-05a5-5db58e6bad0c.htm "SerialBase 类")

[SerialBase 属性](../html/e82084a0-6e96-85cd-f4dc-0aa932aa0ae4.htm "SerialBase 属性")

[BaudRate 属性](../html/1c33c44a-d2b9-badd-df76-09e155d26cc1.htm "BaudRate 属性 ")

[CommunicationPipe 属性](../html/cb528f9b-fe2c-5ac7-6912-2c3c5c1fdf2a.htm "CommunicationPipe 属性 ")

[IsClearCacheBeforeRead 属性](../html/cf2aaeff-e4eb-849a-4783-d10e5e965515.htm "IsClearCacheBeforeRead 属性 ")

[PortName 属性](../html/250b8746-0bb7-4caf-2a43-11109f2bda20.htm "PortName 属性 ")

[ReceiveEmptyDataCount 属性](../html/480b2206-1de4-6a74-6f8e-7c30a2ca05f9.htm "ReceiveEmptyDataCount 属性 ")

[RtsEnable 属性](../html/18c16d98-49b0-5ae2-5c2b-0c6862801c92.htm "RtsEnable 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| SerialBaseCommunicationPipe 属性 |

获取或设置当前的管道信息，管道类型为[CommunicationPipe](ad8e903b-d289-21b0-077e-15ede8005221.htm)的继承类，内置了[PipeTcpNet](26b94bd5-2695-36e7-6378-c3d41846be8f.htm)管道，[PipeUdpNet](5a0b810e-187e-826b-0a13-34fd5576aa0e.htm)管道，[PipeSerialPort](9d924607-3dfc-3916-b532-c5dec48fe347.htm)管道等  
Get or set the current pipeline information, the pipeline type is [CommunicationPipe](ad8e903b-d289-21b0-077e-15ede8005221.htm) inheritance class, [PipeTcpNet](26b94bd5-2695-36e7-6378-c3d41846be8f.htm) pipeline, [PipeUdpNet](5a0b810e-187e-826b-0a13-34fd5576aa0e.htm) pipeline, [PipeSerialPort](9d924607-3dfc-3916-b532-c5dec48fe347.htm) pipeline, etc

**命名空间：**
 [HslCommunication.Serial](8d02123a-6db7-7a98-75db-8e3b63e89789.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public override CommunicationPipe CommunicationPipe { get; set; }
```

```
Public Overrides Property CommunicationPipe As CommunicationPipe
	Get
	Set
```

```
public:
virtual property CommunicationPipe^ CommunicationPipe {
	CommunicationPipe^ get () override;
	void set (CommunicationPipe^ value) override;
}
```

```
abstract CommunicationPipe : CommunicationPipe with get, set
override CommunicationPipe : CommunicationPipe with get, set
```

#### 属性值

类型：[CommunicationPipe](ab3f9c5c-144b-4209-4999-b382a1d94450.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[SerialBase 类](e877e20b-ef65-266c-05a5-5db58e6bad0c.htm)

[HslCommunication.Serial 命名空间](8d02123a-6db7-7a98-75db-8e3b63e89789.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## IsClearCacheBeforeRead 属性 

[原文連結](http://api.hslcommunication.cn/html/cf2aaeff-e4eb-849a-4783-d10e5e965515.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Serial](../html/8d02123a-6db7-7a98-75db-8e3b63e89789.htm "HslCommunication.Serial")

[SerialBase 类](../html/e877e20b-ef65-266c-05a5-5db58e6bad0c.htm "SerialBase 类")

[SerialBase 属性](../html/e82084a0-6e96-85cd-f4dc-0aa932aa0ae4.htm "SerialBase 属性")

[BaudRate 属性](../html/1c33c44a-d2b9-badd-df76-09e155d26cc1.htm "BaudRate 属性 ")

[CommunicationPipe 属性](../html/cb528f9b-fe2c-5ac7-6912-2c3c5c1fdf2a.htm "CommunicationPipe 属性 ")

[IsClearCacheBeforeRead 属性](../html/cf2aaeff-e4eb-849a-4783-d10e5e965515.htm "IsClearCacheBeforeRead 属性 ")

[PortName 属性](../html/250b8746-0bb7-4caf-2a43-11109f2bda20.htm "PortName 属性 ")

[ReceiveEmptyDataCount 属性](../html/480b2206-1de4-6a74-6f8e-7c30a2ca05f9.htm "ReceiveEmptyDataCount 属性 ")

[RtsEnable 属性](../html/18c16d98-49b0-5ae2-5c2b-0c6862801c92.htm "RtsEnable 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| SerialBaseIsClearCacheBeforeRead 属性 |

是否在发送数据前清空缓冲数据，默认是false  
Whether to empty the buffer before sending data, the default is false

**命名空间：**
 [HslCommunication.Serial](8d02123a-6db7-7a98-75db-8e3b63e89789.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public bool IsClearCacheBeforeRead { get; set; }
```

```
Public Property IsClearCacheBeforeRead As Boolean
	Get
	Set
```

```
public:
property bool IsClearCacheBeforeRead {
	bool get ();
	void set (bool value);
}
```

```
member IsClearCacheBeforeRead : bool with get, set
```

#### 属性值

类型：Boolean

![](../icons/SectionExpanded.png)参见

#### 引用

[SerialBase 类](e877e20b-ef65-266c-05a5-5db58e6bad0c.htm)

[HslCommunication.Serial 命名空间](8d02123a-6db7-7a98-75db-8e3b63e89789.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## PortName 属性 

[原文連結](http://api.hslcommunication.cn/html/250b8746-0bb7-4caf-2a43-11109f2bda20.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Serial](../html/8d02123a-6db7-7a98-75db-8e3b63e89789.htm "HslCommunication.Serial")

[SerialBase 类](../html/e877e20b-ef65-266c-05a5-5db58e6bad0c.htm "SerialBase 类")

[SerialBase 属性](../html/e82084a0-6e96-85cd-f4dc-0aa932aa0ae4.htm "SerialBase 属性")

[BaudRate 属性](../html/1c33c44a-d2b9-badd-df76-09e155d26cc1.htm "BaudRate 属性 ")

[CommunicationPipe 属性](../html/cb528f9b-fe2c-5ac7-6912-2c3c5c1fdf2a.htm "CommunicationPipe 属性 ")

[IsClearCacheBeforeRead 属性](../html/cf2aaeff-e4eb-849a-4783-d10e5e965515.htm "IsClearCacheBeforeRead 属性 ")

[PortName 属性](../html/250b8746-0bb7-4caf-2a43-11109f2bda20.htm "PortName 属性 ")

[ReceiveEmptyDataCount 属性](../html/480b2206-1de4-6a74-6f8e-7c30a2ca05f9.htm "ReceiveEmptyDataCount 属性 ")

[RtsEnable 属性](../html/18c16d98-49b0-5ae2-5c2b-0c6862801c92.htm "RtsEnable 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| SerialBasePortName 属性 |

当前连接串口信息的端口号名称  
The port name of the current connection serial port information

**命名空间：**
 [HslCommunication.Serial](8d02123a-6db7-7a98-75db-8e3b63e89789.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public string PortName { get; }
```

```
Public ReadOnly Property PortName As String
	Get
```

```
public:
property String^ PortName {
	String^ get ();
}
```

```
member PortName : string with get
```

#### 属性值

类型：String

![](../icons/SectionExpanded.png)参见

#### 引用

[SerialBase 类](e877e20b-ef65-266c-05a5-5db58e6bad0c.htm)

[HslCommunication.Serial 命名空间](8d02123a-6db7-7a98-75db-8e3b63e89789.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReceiveEmptyDataCount 属性 

[原文連結](http://api.hslcommunication.cn/html/480b2206-1de4-6a74-6f8e-7c30a2ca05f9.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Serial](../html/8d02123a-6db7-7a98-75db-8e3b63e89789.htm "HslCommunication.Serial")

[SerialBase 类](../html/e877e20b-ef65-266c-05a5-5db58e6bad0c.htm "SerialBase 类")

[SerialBase 属性](../html/e82084a0-6e96-85cd-f4dc-0aa932aa0ae4.htm "SerialBase 属性")

[BaudRate 属性](../html/1c33c44a-d2b9-badd-df76-09e155d26cc1.htm "BaudRate 属性 ")

[CommunicationPipe 属性](../html/cb528f9b-fe2c-5ac7-6912-2c3c5c1fdf2a.htm "CommunicationPipe 属性 ")

[IsClearCacheBeforeRead 属性](../html/cf2aaeff-e4eb-849a-4783-d10e5e965515.htm "IsClearCacheBeforeRead 属性 ")

[PortName 属性](../html/250b8746-0bb7-4caf-2a43-11109f2bda20.htm "PortName 属性 ")

[ReceiveEmptyDataCount 属性](../html/480b2206-1de4-6a74-6f8e-7c30a2ca05f9.htm "ReceiveEmptyDataCount 属性 ")

[RtsEnable 属性](../html/18c16d98-49b0-5ae2-5c2b-0c6862801c92.htm "RtsEnable 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| SerialBaseReceiveEmptyDataCount 属性 |

获取或设置连续接收空的数据次数，在数据接收完成时有效，每个单位消耗的时间为[SleepTime](1e3a83ef-b6fc-bc8c-8950-ea3f6ae5c19b.htm)。  
Obtain or set the number of consecutive times to receive empty data, which is valid when the data is received, and the time consumed by each unit is [SleepTime](1e3a83ef-b6fc-bc8c-8950-ea3f6ae5c19b.htm)

**命名空间：**
 [HslCommunication.Serial](8d02123a-6db7-7a98-75db-8e3b63e89789.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public int ReceiveEmptyDataCount { get; set; }
```

```
Public Property ReceiveEmptyDataCount As Integer
	Get
	Set
```

```
public:
property int ReceiveEmptyDataCount {
	int get ();
	void set (int value);
}
```

```
member ReceiveEmptyDataCount : int with get, set
```

#### 属性值

类型：Int32

![](../icons/SectionExpanded.png)参见

#### 引用

[SerialBase 类](e877e20b-ef65-266c-05a5-5db58e6bad0c.htm)

[HslCommunication.Serial 命名空间](8d02123a-6db7-7a98-75db-8e3b63e89789.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## RtsEnable 属性 

[原文連結](http://api.hslcommunication.cn/html/18c16d98-49b0-5ae2-5c2b-0c6862801c92.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Serial](../html/8d02123a-6db7-7a98-75db-8e3b63e89789.htm "HslCommunication.Serial")

[SerialBase 类](../html/e877e20b-ef65-266c-05a5-5db58e6bad0c.htm "SerialBase 类")

[SerialBase 属性](../html/e82084a0-6e96-85cd-f4dc-0aa932aa0ae4.htm "SerialBase 属性")

[BaudRate 属性](../html/1c33c44a-d2b9-badd-df76-09e155d26cc1.htm "BaudRate 属性 ")

[CommunicationPipe 属性](../html/cb528f9b-fe2c-5ac7-6912-2c3c5c1fdf2a.htm "CommunicationPipe 属性 ")

[IsClearCacheBeforeRead 属性](../html/cf2aaeff-e4eb-849a-4783-d10e5e965515.htm "IsClearCacheBeforeRead 属性 ")

[PortName 属性](../html/250b8746-0bb7-4caf-2a43-11109f2bda20.htm "PortName 属性 ")

[ReceiveEmptyDataCount 属性](../html/480b2206-1de4-6a74-6f8e-7c30a2ca05f9.htm "ReceiveEmptyDataCount 属性 ")

[RtsEnable 属性](../html/18c16d98-49b0-5ae2-5c2b-0c6862801c92.htm "RtsEnable 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| SerialBaseRtsEnable 属性 |

获取或设置一个值，该值指示在串行通信中是否启用请求发送 (RTS) 信号。  
Gets or sets a value indicating whether the request sending (RTS) signal is enabled in serial communication.

**命名空间：**
 [HslCommunication.Serial](8d02123a-6db7-7a98-75db-8e3b63e89789.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public bool RtsEnable { get; set; }
```

```
Public Property RtsEnable As Boolean
	Get
	Set
```

```
public:
property bool RtsEnable {
	bool get ();
	void set (bool value);
}
```

```
member RtsEnable : bool with get, set
```

#### 属性值

类型：Boolean

![](../icons/SectionExpanded.png)参见

#### 引用

[SerialBase 类](e877e20b-ef65-266c-05a5-5db58e6bad0c.htm)

[HslCommunication.Serial 命名空间](8d02123a-6db7-7a98-75db-8e3b63e89789.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## SerialBase 方法

[原文連結](http://api.hslcommunication.cn/html/e28d440f-d77d-55fd-4869-070495fa1a3d.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Serial](../html/8d02123a-6db7-7a98-75db-8e3b63e89789.htm "HslCommunication.Serial")

[SerialBase 类](../html/e877e20b-ef65-266c-05a5-5db58e6bad0c.htm "SerialBase 类")

[SerialBase 方法](../html/e28d440f-d77d-55fd-4869-070495fa1a3d.htm "SerialBase 方法")

[Close 方法](../html/992f28b3-b99b-2196-274e-2f25e9bd3c1d.htm "Close 方法 ")

[Dispose 方法](../html/6c344567-ed89-8218-dbf2-b13298737f98.htm "Dispose 方法 ")

[IsOpen 方法](../html/24a0e05a-dd59-2e3c-e4c5-d3bb6ebc8881.htm "IsOpen 方法 ")

[Open 方法](../html/3501fe5b-65b5-a973-5f25-3de45d0054d6.htm "Open 方法 ")

[SerialPortInni 方法](../html/a386ba4a-83eb-3046-54fa-3ecf08001455.htm "SerialPortInni 方法 ")

[ToString 方法](../html/1ac09324-4d19-863d-d683-cbe804ef2124.htm "ToString 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| SerialBase 方法 |

[SerialBase](e877e20b-ef65-266c-05a5-5db58e6bad0c.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [Close](992f28b3-b99b-2196-274e-2f25e9bd3c1d.htm) | 关闭当前的串口连接  Close the current serial connection |
| 受保护的方法 | [DecideWhetherQAMessage](acc40cda-e310-8537-6685-59b67eeb16ff.htm) | 决定当前的消息是否是用于问答机制返回的消息，默认直接返回 true, 实际的情况需要根据协议进行重写方法  To determine whether the current message is the message returned by the question answering mechanism, the default is true. In actual cases, the rewriting method needs to be performed according to the protocol (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [Dispose](c83e3a27-3a09-f67c-9c7f-96b7a5a4c8b0.htm) | 释放当前的对象 |
| 受保护的方法 | [Dispose(Boolean)](c509106d-35a7-2ff2-7dd8-ca67e16d7950.htm) | 释放当前的对象 |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | [ExtraAfterReadFromCoreServer](06a49a15-a3ce-5e39-28c7-2b78851eca2e.htm) | 和服务器交互完成的时候调用的方法，可以根据读写结果进行一些额外的操作，具体的操作需要根据实际的需求来重写实现  The method called when the interaction with the server is completed can perform some additional operations based on the read and write results. The specific operations need to be rewritten according to actual needs. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法代码示例 | [ExtraOnDisconnect](de9dbba2-3b2c-1c5e-53c8-4e6f8a44e2c4.htm) | 根据实际的协议选择是否重写本方法，有些协议在断开连接之前，需要发送一些报文来关闭当前的网络通道  Select whether to rewrite this method according to the actual protocol. Some protocols need to send some packets to close the current network channel before disconnecting. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法代码示例 | [ExtraOnDisconnectAsync](210f593e-16d8-77c4-1178-58fbcfe9b27b.htm) | 根据实际的协议选择是否重写本方法，有些协议在断开连接之前，需要发送一些报文来关闭当前的网络通道  Select whether to rewrite this method according to the actual protocol. Some protocols need to send some packets to close the current network channel before disconnecting. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 受保护的方法 | [GetLogTextFromBinary](67f3b2f7-957c-fa0b-8320-6799237157db.htm) | 获取当前的报文进行日志记录的时候，是否使用二进制的格式记录，默认返回 [LogMsgFormatBinary](2ed2ff87-caec-0e44-f08b-e2b015225a58.htm)，重写可以根据session对象分别返回不同记录模式  Whether to log the current packet in binary format, the default return is [LogMsgFormatBinary](2ed2ff87-caec-0e44-f08b-e2b015225a58.htm). If you want to override it, different recording modes can be returned according to session (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [GetNewNetMessage](16636a1b-a6f0-9c0f-9d9c-c94e26c67459.htm) | 获取一个新的消息对象的方法，需要在继承类里面进行重写  The method to get a new message object needs to be overridden in the inheritance class (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法代码示例 | [InitializationOnConnect](4ab6989f-a5b5-05fe-8b44-bfa45824a2a7.htm) | 根据实际的协议选择是否重写本方法，有些协议在创建连接之后，需要进行一些初始化的信号握手，才能最终建立网络通道。  Whether to rewrite this method is based on the actual protocol. Some protocols require some initial signal handshake to establish a network channel after the connection is created. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法代码示例 | [InitializationOnConnectAsync](7da2af04-d4fc-5045-19a4-4b0d1722eef3.htm) | 根据实际的协议选择是否重写本方法，有些协议在创建连接之后，需要进行一些初始化的信号握手，才能最终建立网络通道。  Whether to rewrite this method is based on the actual protocol. Some protocols require some initial signal handshake to establish a network channel after the connection is created. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [IsOpen](24a0e05a-dd59-2e3c-e4c5-d3bb6ebc8881.htm) | 获取一个值，指示串口是否处于打开状态  Gets a value indicating whether the serial port is open |
| 受保护的方法 | [LogRevcMessage(Byte)](dab1657f-1e75-7fe1-7003-ccedfc69855d.htm) | 使用日志记录一个接收的报文信息  Logs are used to record information about a received packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [LogRevcMessage(Byte, PipeSession)](0f51ae45-575e-4228-f978-01b37d6d5df1.htm) | 使用日志记录一个接收的报文信息  Logs are used to record information about a received packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [LogSendMessage(Byte)](4f90eecc-68a2-1ef2-44ef-369f19bed3ec.htm) | 使用日志记录一个发送的报文信息  Logs are used to record information about a send packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [LogSendMessage(Byte, PipeSession)](209caadf-4f0f-40aa-04c2-ebc61f8dab13.htm) | 使用日志记录一个发送的报文信息  Logs are used to record information about a send packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法 | [Open](3501fe5b-65b5-a973-5f25-3de45d0054d6.htm) | 打开当前的管道信息，返回是否成功打开的结果对象，并通过属性 "Content" 指示当前是否为新创建的连接对象，如果是，则该值为 true  Opens the current pipe information, returns whether the result object was successfully opened, and indicates whether the current connection object is a newly created connection object through the property "Content", if so, the value is true |
| 公共方法 | [PackCommandWithHeader](d0b2cd76-c306-11fc-38c7-af8ad306cce7.htm) | 对当前的命令进行打包处理，通常是携带命令头内容，标记当前的命令的长度信息，需要进行重写，否则默认不打包  The current command is packaged, usually carrying the content of the command header, marking the length of the current command, and it needs to be rewritten, otherwise it is not packaged by default (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServer(Byte)](8146dc0c-1edc-3ca7-428a-a3bdbd35e6c5.htm) | 将二进制的数据发送到管道中去，然后从管道里接收二进制的数据回来，并返回是否成功的结果对象。  Send binary data to the pipeline, and then receive binary data back from the pipeline, and return whether the success of the result object (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServer(IEnumerableByte)](044cfdd9-ec76-b89f-2f2b-5f439025006c.htm) | 将多个数据报文按顺序发到设备，并从设备接收返回的数据内容，然后拼接成一个Byte[]信息，需要重写[UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm)方法才能返回正确的结果。  Send multiple data packets to the device in sequence, and receive the returned data content from the device, and then splicing them into a Byte[] message, you need to rewrite [UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm) method to return the correct result. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServer(Byte, Boolean, Boolean)](b32ce8e8-a959-fcc2-6736-14b39fdb9a40.htm) | 将二进制的数据发送到管道中去，然后从管道里接收二进制的数据回来，并返回是否成功的结果对象。  Send binary data to the pipeline, and then receive binary data back from the pipeline, and return whether the success of the result object (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServer(CommunicationPipe, Byte, Boolean, Boolean)](ef36d56b-df09-f727-e8d3-316d37a13ee9.htm) | 使用指定的管道来进行数据通信，发送原始数据到管道，然后从管道接收相关的数据返回，本方法无锁 (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServerAsync(Byte)](a4de5a24-c67a-7b02-db05-33b85eaf6e1e.htm) | 将二进制的数据发送到管道中去，然后从管道里接收二进制的数据回来，并返回是否成功的结果对象。  Send binary data to the pipeline, and then receive binary data back from the pipeline, and return whether the success of the result object (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServerAsync(IEnumerableByte)](bcdb59b7-a70a-d5c0-1fc6-9ececcd39ea2.htm) | 将多个数据报文按顺序发到设备，并从设备接收返回的数据内容，然后拼接成一个Byte[]信息，需要重写[UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm)方法才能返回正确的结果。  Send multiple data packets to the device in sequence, and receive the returned data content from the device, and then splicing them into a Byte[] message, you need to rewrite [UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm) method to return the correct result. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServerAsync(Byte, Boolean, Boolean)](1c79c204-5ab5-9342-e65d-7455c20ea577.htm) | 将二进制的数据发送到管道中去，然后从管道里接收二进制的数据回来，并返回是否成功的结果对象。  Send binary data to the pipeline, and then receive binary data back from the pipeline, and return whether the success of the result object (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServerAsync(CommunicationPipe, Byte, Boolean, Boolean)](a254a19d-3c29-147b-335b-8adccdf4ee2f.htm) | 使用指定的管道来进行数据通信，发送原始数据到管道，然后从管道接收相关的数据返回，本方法无锁 (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [SerialPortInni(ActionSerialPort)](53585409-6243-9359-1846-bf8b4e266708.htm) | 根据自定义初始化方法进行初始化串口信息  Initialize the serial port information according to the custom initialization method |
| 公共方法 | [SerialPortInni(String)](7dff5e26-6594-291e-4792-b13db9c5d613.htm) | 初始化串口信息，9600波特率，8位数据位，1位停止位，无奇偶校验  Initial serial port information, 9600 baud rate, 8 data bits, 1 stop bit, no parity |
| 公共方法 | [SerialPortInni(String, Int32)](0c7aa163-1d81-54fd-f4f5-103a7e9a341d.htm) | 初始化串口信息，波特率，8位数据位，1位停止位，无奇偶校验  Initializes serial port information, baud rate, 8-bit data bit, 1-bit stop bit, no parity |
| 公共方法 | [SerialPortInni(String, Int32, Int32, StopBits, Parity)](420c818b-48e9-fa02-59e2-f4d5897d6802.htm) | 初始化串口信息，波特率，数据位，停止位，奇偶校验需要全部自己来指定  Start serial port information, baud rate, data bit, stop bit, parity all need to be specified |
| 公共方法 | [SetDtuPipe](f765d1d7-234a-b229-9457-97e25e024dc9.htm) | 将当前的通信对象设置DTU模式，允许传入现成的管道，并返回初始化结果，如果该设备重写了握手报文，就是返回握手结果  Set the current communication object to DTU mode, allow the existing pipe to be passed in, and return the initialization result, if the device rewrites the handshake packet, the handshake result is returned (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [SetDtuPipeAsync](9d1fd12d-b804-3828-a627-a2c0eef1d23f.htm) | 将当前的通信对象设置DTU模式，允许传入现成的管道，并返回初始化结果，如果该设备重写了握手报文，就是返回握手结果  Set the current communication object to DTU mode, allow the existing pipe to be passed in, and return the initialization result, if the device rewrites the handshake packet, the handshake result is returned (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ToString](1ac09324-4d19-863d-d683-cbe804ef2124.htm) | (重写 ObjectToString.) |
| 公共方法 | [UnpackResponseContent](208017a7-8da2-33f8-3273-d7e7e850100f.htm) | 根据对方返回的报文命令，对命令进行基本的拆包，例如各种Modbus协议拆包为统一的核心报文，还支持对报文的验证  According to the message command returned by the other party, the command is basically unpacked, for example, various Modbus protocols are unpacked into a unified core message, and the verification of the message is also supported (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[SerialBase 类](e877e20b-ef65-266c-05a5-5db58e6bad0c.htm)

[HslCommunication.Serial 命名空间](8d02123a-6db7-7a98-75db-8e3b63e89789.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Close 方法 

[原文連結](http://api.hslcommunication.cn/html/992f28b3-b99b-2196-274e-2f25e9bd3c1d.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Serial](../html/8d02123a-6db7-7a98-75db-8e3b63e89789.htm "HslCommunication.Serial")

[SerialBase 类](../html/e877e20b-ef65-266c-05a5-5db58e6bad0c.htm "SerialBase 类")

[SerialBase 方法](../html/e28d440f-d77d-55fd-4869-070495fa1a3d.htm "SerialBase 方法")

[Close 方法](../html/992f28b3-b99b-2196-274e-2f25e9bd3c1d.htm "Close 方法 ")

[Dispose 方法](../html/6c344567-ed89-8218-dbf2-b13298737f98.htm "Dispose 方法 ")

[IsOpen 方法](../html/24a0e05a-dd59-2e3c-e4c5-d3bb6ebc8881.htm "IsOpen 方法 ")

[Open 方法](../html/3501fe5b-65b5-a973-5f25-3de45d0054d6.htm "Open 方法 ")

[SerialPortInni 方法](../html/a386ba4a-83eb-3046-54fa-3ecf08001455.htm "SerialPortInni 方法 ")

[ToString 方法](../html/1ac09324-4d19-863d-d683-cbe804ef2124.htm "ToString 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| SerialBaseClose 方法 |

关闭当前的串口连接  
Close the current serial connection

**命名空间：**
 [HslCommunication.Serial](8d02123a-6db7-7a98-75db-8e3b63e89789.htm)  
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
void Close()
```

```
member Close : unit -> unit 
```

![](../icons/SectionExpanded.png)参见

#### 引用

[SerialBase 类](e877e20b-ef65-266c-05a5-5db58e6bad0c.htm)

[HslCommunication.Serial 命名空间](8d02123a-6db7-7a98-75db-8e3b63e89789.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Dispose 方法 

[原文連結](http://api.hslcommunication.cn/html/6c344567-ed89-8218-dbf2-b13298737f98.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Serial](../html/8d02123a-6db7-7a98-75db-8e3b63e89789.htm "HslCommunication.Serial")

[SerialBase 类](../html/e877e20b-ef65-266c-05a5-5db58e6bad0c.htm "SerialBase 类")

[SerialBase 方法](../html/e28d440f-d77d-55fd-4869-070495fa1a3d.htm "SerialBase 方法")

[Dispose 方法](../html/6c344567-ed89-8218-dbf2-b13298737f98.htm "Dispose 方法 ")

[Dispose 方法](../html/c83e3a27-3a09-f67c-9c7f-96b7a5a4c8b0.htm "Dispose 方法 ")

[Dispose 方法 (Boolean)](../html/c509106d-35a7-2ff2-7dd8-ca67e16d7950.htm "Dispose 方法 (Boolean)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| SerialBaseDispose 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [Dispose](c83e3a27-3a09-f67c-9c7f-96b7a5a4c8b0.htm) | 释放当前的对象 |
| 受保护的方法 | [Dispose(Boolean)](c509106d-35a7-2ff2-7dd8-ca67e16d7950.htm) | 释放当前的对象 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[SerialBase 类](e877e20b-ef65-266c-05a5-5db58e6bad0c.htm)

[HslCommunication.Serial 命名空间](8d02123a-6db7-7a98-75db-8e3b63e89789.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Dispose 方法 

[原文連結](http://api.hslcommunication.cn/html/c83e3a27-3a09-f67c-9c7f-96b7a5a4c8b0.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Serial](../html/8d02123a-6db7-7a98-75db-8e3b63e89789.htm "HslCommunication.Serial")

[SerialBase 类](../html/e877e20b-ef65-266c-05a5-5db58e6bad0c.htm "SerialBase 类")

[SerialBase 方法](../html/e28d440f-d77d-55fd-4869-070495fa1a3d.htm "SerialBase 方法")

[Dispose 方法](../html/6c344567-ed89-8218-dbf2-b13298737f98.htm "Dispose 方法 ")

[Dispose 方法](../html/c83e3a27-3a09-f67c-9c7f-96b7a5a4c8b0.htm "Dispose 方法 ")

[Dispose 方法 (Boolean)](../html/c509106d-35a7-2ff2-7dd8-ca67e16d7950.htm "Dispose 方法 (Boolean)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| SerialBaseDispose 方法 |

释放当前的对象

**命名空间：**
 [HslCommunication.Serial](8d02123a-6db7-7a98-75db-8e3b63e89789.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public void Dispose()
```

```
Public Sub Dispose
```

```
public:
virtual void Dispose() sealed
```

```
abstract Dispose : unit -> unit 
override Dispose : unit -> unit
```

#### 实现

IDisposableDispose

![](../icons/SectionExpanded.png)参见

#### 引用

[SerialBase 类](e877e20b-ef65-266c-05a5-5db58e6bad0c.htm)

[Dispose 重载](6c344567-ed89-8218-dbf2-b13298737f98.htm)

[HslCommunication.Serial 命名空间](8d02123a-6db7-7a98-75db-8e3b63e89789.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Dispose 方法 (Boolean)

[原文連結](http://api.hslcommunication.cn/html/c509106d-35a7-2ff2-7dd8-ca67e16d7950.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Serial](../html/8d02123a-6db7-7a98-75db-8e3b63e89789.htm "HslCommunication.Serial")

[SerialBase 类](../html/e877e20b-ef65-266c-05a5-5db58e6bad0c.htm "SerialBase 类")

[SerialBase 方法](../html/e28d440f-d77d-55fd-4869-070495fa1a3d.htm "SerialBase 方法")

[Dispose 方法](../html/6c344567-ed89-8218-dbf2-b13298737f98.htm "Dispose 方法 ")

[Dispose 方法](../html/c83e3a27-3a09-f67c-9c7f-96b7a5a4c8b0.htm "Dispose 方法 ")

[Dispose 方法 (Boolean)](../html/c509106d-35a7-2ff2-7dd8-ca67e16d7950.htm "Dispose 方法 (Boolean)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| SerialBaseDispose 方法 (Boolean) |

释放当前的对象

**命名空间：**
 [HslCommunication.Serial](8d02123a-6db7-7a98-75db-8e3b63e89789.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
protected virtual void Dispose(
	bool disposing
)
```

```
Protected Overridable Sub Dispose ( 
	disposing As Boolean
)
```

```
protected:
virtual void Dispose(
	bool disposing
)
```

```
abstract Dispose : 
        disposing : bool -> unit 
override Dispose : 
        disposing : bool -> unit
```

#### 参数

disposing
:   类型：SystemBoolean  
    是否在

![](../icons/SectionExpanded.png)参见

#### 引用

[SerialBase 类](e877e20b-ef65-266c-05a5-5db58e6bad0c.htm)

[Dispose 重载](6c344567-ed89-8218-dbf2-b13298737f98.htm)

[HslCommunication.Serial 命名空间](8d02123a-6db7-7a98-75db-8e3b63e89789.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## IsOpen 方法 

[原文連結](http://api.hslcommunication.cn/html/24a0e05a-dd59-2e3c-e4c5-d3bb6ebc8881.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Serial](../html/8d02123a-6db7-7a98-75db-8e3b63e89789.htm "HslCommunication.Serial")

[SerialBase 类](../html/e877e20b-ef65-266c-05a5-5db58e6bad0c.htm "SerialBase 类")

[SerialBase 方法](../html/e28d440f-d77d-55fd-4869-070495fa1a3d.htm "SerialBase 方法")

[Close 方法](../html/992f28b3-b99b-2196-274e-2f25e9bd3c1d.htm "Close 方法 ")

[Dispose 方法](../html/6c344567-ed89-8218-dbf2-b13298737f98.htm "Dispose 方法 ")

[IsOpen 方法](../html/24a0e05a-dd59-2e3c-e4c5-d3bb6ebc8881.htm "IsOpen 方法 ")

[Open 方法](../html/3501fe5b-65b5-a973-5f25-3de45d0054d6.htm "Open 方法 ")

[SerialPortInni 方法](../html/a386ba4a-83eb-3046-54fa-3ecf08001455.htm "SerialPortInni 方法 ")

[ToString 方法](../html/1ac09324-4d19-863d-d683-cbe804ef2124.htm "ToString 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| SerialBaseIsOpen 方法 |

获取一个值，指示串口是否处于打开状态  
Gets a value indicating whether the serial port is open

**命名空间：**
 [HslCommunication.Serial](8d02123a-6db7-7a98-75db-8e3b63e89789.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public bool IsOpen()
```

```
Public Function IsOpen As Boolean
```

```
public:
bool IsOpen()
```

```
member IsOpen : unit -> bool 
```

#### 返回值

类型：Boolean  
是或否

![](../icons/SectionExpanded.png)参见

#### 引用

[SerialBase 类](e877e20b-ef65-266c-05a5-5db58e6bad0c.htm)

[HslCommunication.Serial 命名空间](8d02123a-6db7-7a98-75db-8e3b63e89789.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Open 方法 

[原文連結](http://api.hslcommunication.cn/html/3501fe5b-65b5-a973-5f25-3de45d0054d6.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Serial](../html/8d02123a-6db7-7a98-75db-8e3b63e89789.htm "HslCommunication.Serial")

[SerialBase 类](../html/e877e20b-ef65-266c-05a5-5db58e6bad0c.htm "SerialBase 类")

[SerialBase 方法](../html/e28d440f-d77d-55fd-4869-070495fa1a3d.htm "SerialBase 方法")

[Close 方法](../html/992f28b3-b99b-2196-274e-2f25e9bd3c1d.htm "Close 方法 ")

[Dispose 方法](../html/6c344567-ed89-8218-dbf2-b13298737f98.htm "Dispose 方法 ")

[IsOpen 方法](../html/24a0e05a-dd59-2e3c-e4c5-d3bb6ebc8881.htm "IsOpen 方法 ")

[Open 方法](../html/3501fe5b-65b5-a973-5f25-3de45d0054d6.htm "Open 方法 ")

[SerialPortInni 方法](../html/a386ba4a-83eb-3046-54fa-3ecf08001455.htm "SerialPortInni 方法 ")

[ToString 方法](../html/1ac09324-4d19-863d-d683-cbe804ef2124.htm "ToString 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| SerialBaseOpen 方法 |

打开当前的管道信息，返回是否成功打开的结果对象，并通过属性 "Content" 指示当前是否为新创建的连接对象，如果是，则该值为 true  
Opens the current pipe information, returns whether the result object was successfully opened,
and indicates whether the current connection object is a newly created connection object through the property "Content", if so, the value is true

**命名空间：**
 [HslCommunication.Serial](8d02123a-6db7-7a98-75db-8e3b63e89789.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public virtual OperateResult Open()
```

```
Public Overridable Function Open As OperateResult
```

```
public:
virtual OperateResult^ Open()
```

```
abstract Open : unit -> OperateResult 
override Open : unit -> OperateResult
```

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否打开成功的结果对象

![](../icons/SectionExpanded.png)参见

#### 引用

[SerialBase 类](e877e20b-ef65-266c-05a5-5db58e6bad0c.htm)

[HslCommunication.Serial 命名空间](8d02123a-6db7-7a98-75db-8e3b63e89789.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## SerialPortInni 方法 

[原文連結](http://api.hslcommunication.cn/html/a386ba4a-83eb-3046-54fa-3ecf08001455.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Serial](../html/8d02123a-6db7-7a98-75db-8e3b63e89789.htm "HslCommunication.Serial")

[SerialBase 类](../html/e877e20b-ef65-266c-05a5-5db58e6bad0c.htm "SerialBase 类")

[SerialBase 方法](../html/e28d440f-d77d-55fd-4869-070495fa1a3d.htm "SerialBase 方法")

[SerialPortInni 方法](../html/a386ba4a-83eb-3046-54fa-3ecf08001455.htm "SerialPortInni 方法 ")

[SerialPortInni 方法 (Action(SerialPort))](../html/53585409-6243-9359-1846-bf8b4e266708.htm "SerialPortInni 方法 (Action(SerialPort))")

[SerialPortInni 方法 (String)](../html/7dff5e26-6594-291e-4792-b13db9c5d613.htm "SerialPortInni 方法 (String)")

[SerialPortInni 方法 (String, Int32)](../html/0c7aa163-1d81-54fd-f4f5-103a7e9a341d.htm "SerialPortInni 方法 (String, Int32)")

[SerialPortInni 方法 (String, Int32, Int32, StopBits, Parity)](../html/420c818b-48e9-fa02-59e2-f4d5897d6802.htm "SerialPortInni 方法 (String, Int32, Int32, StopBits, Parity)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| SerialBaseSerialPortInni 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [SerialPortInni(ActionSerialPort)](53585409-6243-9359-1846-bf8b4e266708.htm) | 根据自定义初始化方法进行初始化串口信息  Initialize the serial port information according to the custom initialization method |
| 公共方法 | [SerialPortInni(String)](7dff5e26-6594-291e-4792-b13db9c5d613.htm) | 初始化串口信息，9600波特率，8位数据位，1位停止位，无奇偶校验  Initial serial port information, 9600 baud rate, 8 data bits, 1 stop bit, no parity |
| 公共方法 | [SerialPortInni(String, Int32)](0c7aa163-1d81-54fd-f4f5-103a7e9a341d.htm) | 初始化串口信息，波特率，8位数据位，1位停止位，无奇偶校验  Initializes serial port information, baud rate, 8-bit data bit, 1-bit stop bit, no parity |
| 公共方法 | [SerialPortInni(String, Int32, Int32, StopBits, Parity)](420c818b-48e9-fa02-59e2-f4d5897d6802.htm) | 初始化串口信息，波特率，数据位，停止位，奇偶校验需要全部自己来指定  Start serial port information, baud rate, data bit, stop bit, parity all need to be specified |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[SerialBase 类](e877e20b-ef65-266c-05a5-5db58e6bad0c.htm)

[HslCommunication.Serial 命名空间](8d02123a-6db7-7a98-75db-8e3b63e89789.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## SerialPortInni 方法 (Action(SerialPort))

[原文連結](http://api.hslcommunication.cn/html/53585409-6243-9359-1846-bf8b4e266708.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Serial](../html/8d02123a-6db7-7a98-75db-8e3b63e89789.htm "HslCommunication.Serial")

[SerialBase 类](../html/e877e20b-ef65-266c-05a5-5db58e6bad0c.htm "SerialBase 类")

[SerialBase 方法](../html/e28d440f-d77d-55fd-4869-070495fa1a3d.htm "SerialBase 方法")

[SerialPortInni 方法](../html/a386ba4a-83eb-3046-54fa-3ecf08001455.htm "SerialPortInni 方法 ")

[SerialPortInni 方法 (Action(SerialPort))](../html/53585409-6243-9359-1846-bf8b4e266708.htm "SerialPortInni 方法 (Action(SerialPort))")

[SerialPortInni 方法 (String)](../html/7dff5e26-6594-291e-4792-b13db9c5d613.htm "SerialPortInni 方法 (String)")

[SerialPortInni 方法 (String, Int32)](../html/0c7aa163-1d81-54fd-f4f5-103a7e9a341d.htm "SerialPortInni 方法 (String, Int32)")

[SerialPortInni 方法 (String, Int32, Int32, StopBits, Parity)](../html/420c818b-48e9-fa02-59e2-f4d5897d6802.htm "SerialPortInni 方法 (String, Int32, Int32, StopBits, Parity)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| SerialBaseSerialPortInni 方法 (ActionSerialPort) |

根据自定义初始化方法进行初始化串口信息  
Initialize the serial port information according to the custom initialization method

**命名空间：**
 [HslCommunication.Serial](8d02123a-6db7-7a98-75db-8e3b63e89789.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public void SerialPortInni(
	Action<SerialPort> initi
)
```

```
Public Sub SerialPortInni ( 
	initi As Action(Of SerialPort)
)
```

```
public:
void SerialPortInni(
	Action<SerialPort^>^ initi
)
```

```
member SerialPortInni : 
        initi : Action<SerialPort> -> unit 
```

#### 参数

initi
:   类型：SystemActionSerialPort  
    初始化的委托方法

![](../icons/SectionExpanded.png)参见

#### 引用

[SerialBase 类](e877e20b-ef65-266c-05a5-5db58e6bad0c.htm)

[SerialPortInni 重载](a386ba4a-83eb-3046-54fa-3ecf08001455.htm)

[HslCommunication.Serial 命名空间](8d02123a-6db7-7a98-75db-8e3b63e89789.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## SerialPortInni 方法 (String)

[原文連結](http://api.hslcommunication.cn/html/7dff5e26-6594-291e-4792-b13db9c5d613.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Serial](../html/8d02123a-6db7-7a98-75db-8e3b63e89789.htm "HslCommunication.Serial")

[SerialBase 类](../html/e877e20b-ef65-266c-05a5-5db58e6bad0c.htm "SerialBase 类")

[SerialBase 方法](../html/e28d440f-d77d-55fd-4869-070495fa1a3d.htm "SerialBase 方法")

[SerialPortInni 方法](../html/a386ba4a-83eb-3046-54fa-3ecf08001455.htm "SerialPortInni 方法 ")

[SerialPortInni 方法 (Action(SerialPort))](../html/53585409-6243-9359-1846-bf8b4e266708.htm "SerialPortInni 方法 (Action(SerialPort))")

[SerialPortInni 方法 (String)](../html/7dff5e26-6594-291e-4792-b13db9c5d613.htm "SerialPortInni 方法 (String)")

[SerialPortInni 方法 (String, Int32)](../html/0c7aa163-1d81-54fd-f4f5-103a7e9a341d.htm "SerialPortInni 方法 (String, Int32)")

[SerialPortInni 方法 (String, Int32, Int32, StopBits, Parity)](../html/420c818b-48e9-fa02-59e2-f4d5897d6802.htm "SerialPortInni 方法 (String, Int32, Int32, StopBits, Parity)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| SerialBaseSerialPortInni 方法 (String) |

初始化串口信息，9600波特率，8位数据位，1位停止位，无奇偶校验  
Initial serial port information, 9600 baud rate, 8 data bits, 1 stop bit, no parity

**命名空间：**
 [HslCommunication.Serial](8d02123a-6db7-7a98-75db-8e3b63e89789.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public virtual void SerialPortInni(
	string portName
)
```

```
Public Overridable Sub SerialPortInni ( 
	portName As String
)
```

```
public:
virtual void SerialPortInni(
	String^ portName
)
```

```
abstract SerialPortInni : 
        portName : string -> unit 
override SerialPortInni : 
        portName : string -> unit
```

#### 参数

portName
:   类型：SystemString  
    端口号信息，例如"COM3"

![](../icons/SectionExpanded.png)备注

portName 支持格式化的方式，例如输入 COM3-9600-8-N-1，COM5-19200-7-E-2，其中奇偶校验的字母可选，N:无校验，O：奇校验，E:偶校验，停止位可选 0, 1, 2, 1.5 四种选项

![](../icons/SectionExpanded.png)参见

#### 引用

[SerialBase 类](e877e20b-ef65-266c-05a5-5db58e6bad0c.htm)

[SerialPortInni 重载](a386ba4a-83eb-3046-54fa-3ecf08001455.htm)

[HslCommunication.Serial 命名空间](8d02123a-6db7-7a98-75db-8e3b63e89789.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## SerialPortInni 方法 (String, Int32)

[原文連結](http://api.hslcommunication.cn/html/0c7aa163-1d81-54fd-f4f5-103a7e9a341d.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Serial](../html/8d02123a-6db7-7a98-75db-8e3b63e89789.htm "HslCommunication.Serial")

[SerialBase 类](../html/e877e20b-ef65-266c-05a5-5db58e6bad0c.htm "SerialBase 类")

[SerialBase 方法](../html/e28d440f-d77d-55fd-4869-070495fa1a3d.htm "SerialBase 方法")

[SerialPortInni 方法](../html/a386ba4a-83eb-3046-54fa-3ecf08001455.htm "SerialPortInni 方法 ")

[SerialPortInni 方法 (Action(SerialPort))](../html/53585409-6243-9359-1846-bf8b4e266708.htm "SerialPortInni 方法 (Action(SerialPort))")

[SerialPortInni 方法 (String)](../html/7dff5e26-6594-291e-4792-b13db9c5d613.htm "SerialPortInni 方法 (String)")

[SerialPortInni 方法 (String, Int32)](../html/0c7aa163-1d81-54fd-f4f5-103a7e9a341d.htm "SerialPortInni 方法 (String, Int32)")

[SerialPortInni 方法 (String, Int32, Int32, StopBits, Parity)](../html/420c818b-48e9-fa02-59e2-f4d5897d6802.htm "SerialPortInni 方法 (String, Int32, Int32, StopBits, Parity)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| SerialBaseSerialPortInni 方法 (String, Int32) |

初始化串口信息，波特率，8位数据位，1位停止位，无奇偶校验  
Initializes serial port information, baud rate, 8-bit data bit, 1-bit stop bit, no parity

**命名空间：**
 [HslCommunication.Serial](8d02123a-6db7-7a98-75db-8e3b63e89789.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public virtual void SerialPortInni(
	string portName,
	int baudRate
)
```

```
Public Overridable Sub SerialPortInni ( 
	portName As String,
	baudRate As Integer
)
```

```
public:
virtual void SerialPortInni(
	String^ portName, 
	int baudRate
)
```

```
abstract SerialPortInni : 
        portName : string * 
        baudRate : int -> unit 
override SerialPortInni : 
        portName : string * 
        baudRate : int -> unit
```

#### 参数

portName
:   类型：SystemString  
    端口号信息，例如"COM3"

baudRate
:   类型：SystemInt32  
    波特率

![](../icons/SectionExpanded.png)参见

#### 引用

[SerialBase 类](e877e20b-ef65-266c-05a5-5db58e6bad0c.htm)

[SerialPortInni 重载](a386ba4a-83eb-3046-54fa-3ecf08001455.htm)

[HslCommunication.Serial 命名空间](8d02123a-6db7-7a98-75db-8e3b63e89789.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## SerialPortInni 方法 (String, Int32, Int32, StopBits, Parity)

[原文連結](http://api.hslcommunication.cn/html/420c818b-48e9-fa02-59e2-f4d5897d6802.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Serial](../html/8d02123a-6db7-7a98-75db-8e3b63e89789.htm "HslCommunication.Serial")

[SerialBase 类](../html/e877e20b-ef65-266c-05a5-5db58e6bad0c.htm "SerialBase 类")

[SerialBase 方法](../html/e28d440f-d77d-55fd-4869-070495fa1a3d.htm "SerialBase 方法")

[SerialPortInni 方法](../html/a386ba4a-83eb-3046-54fa-3ecf08001455.htm "SerialPortInni 方法 ")

[SerialPortInni 方法 (Action(SerialPort))](../html/53585409-6243-9359-1846-bf8b4e266708.htm "SerialPortInni 方法 (Action(SerialPort))")

[SerialPortInni 方法 (String)](../html/7dff5e26-6594-291e-4792-b13db9c5d613.htm "SerialPortInni 方法 (String)")

[SerialPortInni 方法 (String, Int32)](../html/0c7aa163-1d81-54fd-f4f5-103a7e9a341d.htm "SerialPortInni 方法 (String, Int32)")

[SerialPortInni 方法 (String, Int32, Int32, StopBits, Parity)](../html/420c818b-48e9-fa02-59e2-f4d5897d6802.htm "SerialPortInni 方法 (String, Int32, Int32, StopBits, Parity)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| SerialBaseSerialPortInni 方法 (String, Int32, Int32, StopBits, Parity) |

初始化串口信息，波特率，数据位，停止位，奇偶校验需要全部自己来指定  
Start serial port information, baud rate, data bit, stop bit, parity all need to be specified

**命名空间：**
 [HslCommunication.Serial](8d02123a-6db7-7a98-75db-8e3b63e89789.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public virtual void SerialPortInni(
	string portName,
	int baudRate,
	int dataBits,
	StopBits stopBits,
	Parity parity
)
```

```
Public Overridable Sub SerialPortInni ( 
	portName As String,
	baudRate As Integer,
	dataBits As Integer,
	stopBits As StopBits,
	parity As Parity
)
```

```
public:
virtual void SerialPortInni(
	String^ portName, 
	int baudRate, 
	int dataBits, 
	StopBits stopBits, 
	Parity parity
)
```

```
abstract SerialPortInni : 
        portName : string * 
        baudRate : int * 
        dataBits : int * 
        stopBits : StopBits * 
        parity : Parity -> unit 
override SerialPortInni : 
        portName : string * 
        baudRate : int * 
        dataBits : int * 
        stopBits : StopBits * 
        parity : Parity -> unit
```

#### 参数

portName
:   类型：SystemString  
    端口号信息，例如"COM3"

baudRate
:   类型：SystemInt32  
    波特率

dataBits
:   类型：SystemInt32  
    数据位

stopBits
:   类型：System.IO.PortsStopBits  
    停止位

parity
:   类型：System.IO.PortsParity  
    奇偶校验

![](../icons/SectionExpanded.png)参见

#### 引用

[SerialBase 类](e877e20b-ef65-266c-05a5-5db58e6bad0c.htm)

[SerialPortInni 重载](a386ba4a-83eb-3046-54fa-3ecf08001455.htm)

[HslCommunication.Serial 命名空间](8d02123a-6db7-7a98-75db-8e3b63e89789.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ToString 方法 

[原文連結](http://api.hslcommunication.cn/html/1ac09324-4d19-863d-d683-cbe804ef2124.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Serial](../html/8d02123a-6db7-7a98-75db-8e3b63e89789.htm "HslCommunication.Serial")

[SerialBase 类](../html/e877e20b-ef65-266c-05a5-5db58e6bad0c.htm "SerialBase 类")

[SerialBase 方法](../html/e28d440f-d77d-55fd-4869-070495fa1a3d.htm "SerialBase 方法")

[Close 方法](../html/992f28b3-b99b-2196-274e-2f25e9bd3c1d.htm "Close 方法 ")

[Dispose 方法](../html/6c344567-ed89-8218-dbf2-b13298737f98.htm "Dispose 方法 ")

[IsOpen 方法](../html/24a0e05a-dd59-2e3c-e4c5-d3bb6ebc8881.htm "IsOpen 方法 ")

[Open 方法](../html/3501fe5b-65b5-a973-5f25-3de45d0054d6.htm "Open 方法 ")

[SerialPortInni 方法](../html/a386ba4a-83eb-3046-54fa-3ecf08001455.htm "SerialPortInni 方法 ")

[ToString 方法](../html/1ac09324-4d19-863d-d683-cbe804ef2124.htm "ToString 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| SerialBaseToString 方法 |

[缺少 "M:HslCommunication.Serial.SerialBase.ToString" 的 <summary> 文档]

**命名空间：**
 [HslCommunication.Serial](8d02123a-6db7-7a98-75db-8e3b63e89789.htm)  
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

[缺少 "M:HslCommunication.Serial.SerialBase.ToString" 的 <returns> 文档]

![](../icons/SectionExpanded.png)参见

#### 引用

[SerialBase 类](e877e20b-ef65-266c-05a5-5db58e6bad0c.htm)

[HslCommunication.Serial 命名空间](8d02123a-6db7-7a98-75db-8e3b63e89789.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## SerialBase 字段

[原文連結](http://api.hslcommunication.cn/html/f6820884-a072-a4e4-125e-4d91fc28953e.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Serial](../html/8d02123a-6db7-7a98-75db-8e3b63e89789.htm "HslCommunication.Serial")

[SerialBase 类](../html/e877e20b-ef65-266c-05a5-5db58e6bad0c.htm "SerialBase 类")

[SerialBase 构造函数](../html/c67cfb2a-45e7-d93d-8c47-8ac0aa342052.htm "SerialBase 构造函数 ")

[SerialBase 属性](../html/e82084a0-6e96-85cd-f4dc-0aa932aa0ae4.htm "SerialBase 属性")

[SerialBase 方法](../html/e28d440f-d77d-55fd-4869-070495fa1a3d.htm "SerialBase 方法")

[SerialBase 字段](../html/f6820884-a072-a4e4-125e-4d91fc28953e.htm "SerialBase 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| SerialBase 字段 |

[SerialBase](e877e20b-ef65-266c-05a5-5db58e6bad0c.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)字段

|  | 名称 | 说明 |
| --- | --- | --- |
| 受保护的字段 | [LogMsgFormatBinary](2ed2ff87-caec-0e44-f08b-e2b015225a58.htm) | 设置日志记录报文是否二进制，如果为False，那就使用ASCII码  Set whether the log message is binary, if it is False, then use ASCII code (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[SerialBase 类](e877e20b-ef65-266c-05a5-5db58e6bad0c.htm)

[HslCommunication.Serial 命名空间](8d02123a-6db7-7a98-75db-8e3b63e89789.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## SoftCRC16 类

[原文連結](http://api.hslcommunication.cn/html/b01f7f5f-55ee-cbf8-faf3-bb5ca2d26e57.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Serial](../html/8d02123a-6db7-7a98-75db-8e3b63e89789.htm "HslCommunication.Serial")

[SoftCRC16 类](../html/b01f7f5f-55ee-cbf8-faf3-bb5ca2d26e57.htm "SoftCRC16 类")

[SoftCRC16 构造函数](../html/30112f17-18ab-9959-89fc-d7f34429662b.htm "SoftCRC16 构造函数 ")

[SoftCRC16 方法](../html/3507e884-fbdf-d41d-f15d-a2ada0d8761e.htm "SoftCRC16 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| SoftCRC16 类 |

用于CRC16验证的类，提供了标准的验证方法，可以方便快速的对数据进行CRC校验  
The class for CRC16 validation provides a standard validation method that makes it easy to CRC data quickly

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  HslCommunication.SerialSoftCRC16

**命名空间：**
 [HslCommunication.Serial](8d02123a-6db7-7a98-75db-8e3b63e89789.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class SoftCRC16
```

```
Public Class SoftCRC16
```

```
public ref class SoftCRC16
```

```
type SoftCRC16 =  class end
```

SoftCRC16 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [SoftCRC16](30112f17-18ab-9959-89fc-d7f34429662b.htm) | 初始化 SoftCRC16 类的一个新实例 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法静态成员 | [CheckCRC16(Byte)](41b3f829-0a78-428b-6f56-c25ef3ec69a6.htm) | 来校验对应的接收数据的CRC校验码，默认多项式码为0xA001  To verify the CRC check code corresponding to the received data, the default polynomial code is 0xA001 |
| 公共方法静态成员 | [CheckCRC16(Byte, Byte, Byte)](3fef9122-70f9-9b89-a276-912f8fd57dc0.htm) | 指定多项式码来校验对应的接收数据的CRC校验码  Specifies a polynomial code to validate the corresponding CRC check code for the received data |
| 公共方法静态成员 | [CRC16(Byte)](a1b697a9-589d-45e3-f316-bacda86403a4.htm) | 获取对应的数据的CRC校验码，默认多项式码为0xA001  Get the CRC check code of the corresponding data, the default polynomial code is 0xA001 |
| 公共方法静态成员 | [CRC16(Byte, Byte, Byte, Byte, Byte)](df1b5b2e-45f0-71f9-faa1-0d990942a99a.htm) | 通过指定多项式码来获取对应的数据的CRC校验码  The CRC check code of the corresponding data is obtained by specifying the polynomial code |
| 公共方法静态成员 | [CRC16Only](6dfc8cf2-c069-2acf-a327-f728f1597d23.htm) | 通过指定多项式码来获取对应的数据的CRC校验码  The CRC check code of the corresponding data is obtained by specifying the polynomial code |
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

![](../icons/SectionExpanded.png)备注

本类提供了几个静态的方法，用来进行CRC16码的计算和验证的，多项式码可以自己指定配置，但是预置的寄存器为0xFF 0xFF

![](../icons/SectionExpanded.png)示例

先演示如何校验一串数据的CRC码

SoftCRC16示例

[复制](# "复制")

```
// 进行CRC校验，例如从modbus接收的数据 01 03 00 00 00 02 C4 0B

bool check = HslCommunication.Serial.SoftCRC16.CheckCRC16(
    HslCommunication.BasicFramework.SoftBasic.HexStringToBytes( "01 03 00 00 00 02 C4 0B" ) );

if (check)
{
    Console.WriteLine( "check success!" );       // 此处success
}
else
{
    Console.WriteLine( "check failed!" );
}

// 上述的代码是使用了多项式码 A0 01，检验成功，如果您的多项式不是这个，比如 B8 08 那么就需要按照如下的方式
check = HslCommunication.Serial.SoftCRC16.CheckCRC16(
    HslCommunication.BasicFramework.SoftBasic.HexStringToBytes( "01 03 00 00 00 02 C4 0B" ), 0xB8, 0x08 );

if (check)
{
    Console.WriteLine( "check success!" );
}
else
{
    Console.WriteLine( "check failed!" );    // 此处failed
}
```

然后下面是如何生成你自己的CRC校验码

SoftCRC16示例

[复制](# "复制")

```
// 计算CRC码，比如我要给"01 03 00 00 00 02"增加crc校验

byte[] buffer = HslCommunication.Serial.SoftCRC16.CRC16(
    HslCommunication.BasicFramework.SoftBasic.HexStringToBytes( "01 03 00 00 00 02" ) );

// buffer 就是 "01 03 00 00 00 02 C4 0B" 然后就可以发送到modbus的串口了

// 如果需要自己指定多项式码 B8 08
buffer = HslCommunication.Serial.SoftCRC16.CRC16(
    HslCommunication.BasicFramework.SoftBasic.HexStringToBytes( "01 03 00 00 00 02" ), 0xB8, 0x08 );
```

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.Serial 命名空间](8d02123a-6db7-7a98-75db-8e3b63e89789.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## SoftCRC16 构造函数 

[原文連結](http://api.hslcommunication.cn/html/30112f17-18ab-9959-89fc-d7f34429662b.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Serial](../html/8d02123a-6db7-7a98-75db-8e3b63e89789.htm "HslCommunication.Serial")

[SoftCRC16 类](../html/b01f7f5f-55ee-cbf8-faf3-bb5ca2d26e57.htm "SoftCRC16 类")

[SoftCRC16 构造函数](../html/30112f17-18ab-9959-89fc-d7f34429662b.htm "SoftCRC16 构造函数 ")

[SoftCRC16 方法](../html/3507e884-fbdf-d41d-f15d-a2ada0d8761e.htm "SoftCRC16 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| SoftCRC16 构造函数 |

初始化 [SoftCRC16](b01f7f5f-55ee-cbf8-faf3-bb5ca2d26e57.htm) 类的一个新实例

**命名空间：**
 [HslCommunication.Serial](8d02123a-6db7-7a98-75db-8e3b63e89789.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public SoftCRC16()
```

```
Public Sub New
```

```
public:
SoftCRC16()
```

```
new : unit -> SoftCRC16
```

![](../icons/SectionExpanded.png)参见

#### 引用

[SoftCRC16 类](b01f7f5f-55ee-cbf8-faf3-bb5ca2d26e57.htm)

[HslCommunication.Serial 命名空间](8d02123a-6db7-7a98-75db-8e3b63e89789.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## SoftCRC16 方法

[原文連結](http://api.hslcommunication.cn/html/3507e884-fbdf-d41d-f15d-a2ada0d8761e.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Serial](../html/8d02123a-6db7-7a98-75db-8e3b63e89789.htm "HslCommunication.Serial")

[SoftCRC16 类](../html/b01f7f5f-55ee-cbf8-faf3-bb5ca2d26e57.htm "SoftCRC16 类")

[SoftCRC16 方法](../html/3507e884-fbdf-d41d-f15d-a2ada0d8761e.htm "SoftCRC16 方法")

[CheckCRC16 方法](../html/a8299ec2-70a9-ee66-c40d-6eb6fbdff043.htm "CheckCRC16 方法 ")

[CRC16 方法](../html/f27109a7-52f2-c320-c229-61bbb963edcd.htm "CRC16 方法 ")

[CRC16Only 方法](../html/6dfc8cf2-c069-2acf-a327-f728f1597d23.htm "CRC16Only 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| SoftCRC16 方法 |

[SoftCRC16](b01f7f5f-55ee-cbf8-faf3-bb5ca2d26e57.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法静态成员 | [CheckCRC16(Byte)](41b3f829-0a78-428b-6f56-c25ef3ec69a6.htm) | 来校验对应的接收数据的CRC校验码，默认多项式码为0xA001  To verify the CRC check code corresponding to the received data, the default polynomial code is 0xA001 |
| 公共方法静态成员 | [CheckCRC16(Byte, Byte, Byte)](3fef9122-70f9-9b89-a276-912f8fd57dc0.htm) | 指定多项式码来校验对应的接收数据的CRC校验码  Specifies a polynomial code to validate the corresponding CRC check code for the received data |
| 公共方法静态成员 | [CRC16(Byte)](a1b697a9-589d-45e3-f316-bacda86403a4.htm) | 获取对应的数据的CRC校验码，默认多项式码为0xA001  Get the CRC check code of the corresponding data, the default polynomial code is 0xA001 |
| 公共方法静态成员 | [CRC16(Byte, Byte, Byte, Byte, Byte)](df1b5b2e-45f0-71f9-faa1-0d990942a99a.htm) | 通过指定多项式码来获取对应的数据的CRC校验码  The CRC check code of the corresponding data is obtained by specifying the polynomial code |
| 公共方法静态成员 | [CRC16Only](6dfc8cf2-c069-2acf-a327-f728f1597d23.htm) | 通过指定多项式码来获取对应的数据的CRC校验码  The CRC check code of the corresponding data is obtained by specifying the polynomial code |
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

[SoftCRC16 类](b01f7f5f-55ee-cbf8-faf3-bb5ca2d26e57.htm)

[HslCommunication.Serial 命名空间](8d02123a-6db7-7a98-75db-8e3b63e89789.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## CheckCRC16 方法 

[原文連結](http://api.hslcommunication.cn/html/a8299ec2-70a9-ee66-c40d-6eb6fbdff043.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Serial](../html/8d02123a-6db7-7a98-75db-8e3b63e89789.htm "HslCommunication.Serial")

[SoftCRC16 类](../html/b01f7f5f-55ee-cbf8-faf3-bb5ca2d26e57.htm "SoftCRC16 类")

[SoftCRC16 方法](../html/3507e884-fbdf-d41d-f15d-a2ada0d8761e.htm "SoftCRC16 方法")

[CheckCRC16 方法](../html/a8299ec2-70a9-ee66-c40d-6eb6fbdff043.htm "CheckCRC16 方法 ")

[CheckCRC16 方法 (Byte[])](../html/41b3f829-0a78-428b-6f56-c25ef3ec69a6.htm "CheckCRC16 方法 (Byte[])")

[CheckCRC16 方法 (Byte[], Byte, Byte)](../html/3fef9122-70f9-9b89-a276-912f8fd57dc0.htm "CheckCRC16 方法 (Byte[], Byte, Byte)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| SoftCRC16CheckCRC16 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法静态成员 | [CheckCRC16(Byte)](41b3f829-0a78-428b-6f56-c25ef3ec69a6.htm) | 来校验对应的接收数据的CRC校验码，默认多项式码为0xA001  To verify the CRC check code corresponding to the received data, the default polynomial code is 0xA001 |
| 公共方法静态成员 | [CheckCRC16(Byte, Byte, Byte)](3fef9122-70f9-9b89-a276-912f8fd57dc0.htm) | 指定多项式码来校验对应的接收数据的CRC校验码  Specifies a polynomial code to validate the corresponding CRC check code for the received data |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[SoftCRC16 类](b01f7f5f-55ee-cbf8-faf3-bb5ca2d26e57.htm)

[HslCommunication.Serial 命名空间](8d02123a-6db7-7a98-75db-8e3b63e89789.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## CheckCRC16 方法 (Byte[])

[原文連結](http://api.hslcommunication.cn/html/41b3f829-0a78-428b-6f56-c25ef3ec69a6.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Serial](../html/8d02123a-6db7-7a98-75db-8e3b63e89789.htm "HslCommunication.Serial")

[SoftCRC16 类](../html/b01f7f5f-55ee-cbf8-faf3-bb5ca2d26e57.htm "SoftCRC16 类")

[SoftCRC16 方法](../html/3507e884-fbdf-d41d-f15d-a2ada0d8761e.htm "SoftCRC16 方法")

[CheckCRC16 方法](../html/a8299ec2-70a9-ee66-c40d-6eb6fbdff043.htm "CheckCRC16 方法 ")

[CheckCRC16 方法 (Byte[])](../html/41b3f829-0a78-428b-6f56-c25ef3ec69a6.htm "CheckCRC16 方法 (Byte[])")

[CheckCRC16 方法 (Byte[], Byte, Byte)](../html/3fef9122-70f9-9b89-a276-912f8fd57dc0.htm "CheckCRC16 方法 (Byte[], Byte, Byte)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| SoftCRC16CheckCRC16 方法 (Byte) |

来校验对应的接收数据的CRC校验码，默认多项式码为0xA001  
To verify the CRC check code corresponding to the received data, the default polynomial code is 0xA001

**命名空间：**
 [HslCommunication.Serial](8d02123a-6db7-7a98-75db-8e3b63e89789.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static bool CheckCRC16(
	byte[] value
)
```

```
Public Shared Function CheckCRC16 ( 
	value As Byte()
) As Boolean
```

```
public:
static bool CheckCRC16(
	array<unsigned char>^ value
)
```

```
static member CheckCRC16 : 
        value : byte[] -> bool 
```

#### 参数

value
:   类型：SystemByte  
    需要校验的数据，带CRC校验码

#### 返回值

类型：Boolean  
返回校验成功与否

![](../icons/SectionExpanded.png)参见

#### 引用

[SoftCRC16 类](b01f7f5f-55ee-cbf8-faf3-bb5ca2d26e57.htm)

[CheckCRC16 重载](a8299ec2-70a9-ee66-c40d-6eb6fbdff043.htm)

[HslCommunication.Serial 命名空间](8d02123a-6db7-7a98-75db-8e3b63e89789.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)