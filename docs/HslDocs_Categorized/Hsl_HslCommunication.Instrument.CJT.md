# HslCommunication - HslCommunication.Instrument.CJT

> 分類頁數: 30



---
## HslCommunication.Instrument.CJT

[原文連結](http://api.hslcommunication.cn/html/32a7bbe4-7443-0a40-8e16-c1df11e7a933.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.CJT](../html/32a7bbe4-7443-0a40-8e16-c1df11e7a933.htm "HslCommunication.Instrument.CJT")

[CJT188 类](../html/7845b81d-5321-3a2e-346d-ac4fb15e3f2a.htm "CJT188 类")

[CJT188OverTcp 类](../html/d4e97b13-9198-62d1-6eb0-8aeb8c19eda4.htm "CJT188OverTcp 类")

[CjtFlowRate 类](../html/e7588170-6c94-a0e9-5ad4-437a682d797a.htm "CjtFlowRate 类")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslCommunication.Instrument.CJT 命名空间 |

[缺少 "N:HslCommunication.Instrument.CJT" 的 <summary> 文档]

![](../icons/SectionExpanded.png)类

|  | 类 | 说明 |
| --- | --- | --- |
| 公共类 | [CJT188](7845b81d-5321-3a2e-346d-ac4fb15e3f2a.htm) | 城市建设部的188协议，基于DJ/T188-2004实现的协议 |
| 公共类 | [CJT188OverTcp](d4e97b13-9198-62d1-6eb0-8aeb8c19eda4.htm) | CJT188串口透传协议 |
| 公共类 | [CjtFlowRate](e7588170-6c94-a0e9-5ad4-437a682d797a.htm) | CJT协议的流量数据，主要是用来获取水表流量及燃气流量的 |

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## CJT188 类

[原文連結](http://api.hslcommunication.cn/html/7845b81d-5321-3a2e-346d-ac4fb15e3f2a.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.CJT](../html/32a7bbe4-7443-0a40-8e16-c1df11e7a933.htm "HslCommunication.Instrument.CJT")

[CJT188 类](../html/7845b81d-5321-3a2e-346d-ac4fb15e3f2a.htm "CJT188 类")

[CJT188 构造函数](../html/72a09742-43f7-3f9c-7d8f-0872d228bb0f.htm "CJT188 构造函数 ")

[CJT188 属性](../html/be37af45-b2d6-9a32-9d71-e9b361f46ab2.htm "CJT188 属性")

[CJT188 方法](../html/04972e44-8c91-c7ca-0ab7-28b426216410.htm "CJT188 方法")

[CJT188 字段](../html/e0c4f1ed-6183-52ba-c1f5-cdf5993cd2cc.htm "CJT188 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CJT188 类 |

城市建设部的188协议，基于DJ/T188-2004实现的协议

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  [HslCommunication.Core.NetBinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)  
    [HslCommunication.Core.DeviceDeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)  
      [HslCommunication.Core.DeviceDeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)  
        HslCommunication.Instrument.CJTCJT188

**命名空间：**
 [HslCommunication.Instrument.CJT](32a7bbe4-7443-0a40-8e16-c1df11e7a933.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class CJT188 : DeviceSerialPort, ICjt188, 
	IReadWriteDevice, IReadWriteNet
```

```
Public Class CJT188
	Inherits DeviceSerialPort
	Implements ICjt188, IReadWriteDevice, IReadWriteNet
```

```
public ref class CJT188 : public DeviceSerialPort, 
	ICjt188, IReadWriteDevice, IReadWriteNet
```

```
type CJT188 =  
    class
        inherit DeviceSerialPort
        interface ICjt188
        interface IReadWriteDevice
        interface IReadWriteNet
    end
```

CJT188 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [CJT188](72a09742-43f7-3f9c-7d8f-0872d228bb0f.htm) | 指定地址域来实例化一个对象，地址域是一个14个字符的BCD码，例如：14910000729011  Specify the address field, to instantiate an object, which address field is a 14-character BCD code, for example: 14910000729011 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [BaudRate](5554f563-5e29-699c-435b-c8e7ad5fa0e4.htm) | 当前连接串口信息的波特率  Baud rate of current connection serial port information (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 公共属性代码示例 | [ByteTransform](f35118c9-3691-f752-ae6f-2c5549714faa.htm) | 当前的数据变换机制，当你需要从字节数据转换类型数据的时候需要。  The current data transformation mechanism is required when you need to convert type data from byte data. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共属性 | [CommunicationPipe](a4fcb477-64d2-8da7-c712-687f0261755b.htm) | 获取或设置当前的管道信息，管道类型为[CommunicationPipe](ad8e903b-d289-21b0-077e-15ede8005221.htm)的继承类，内置了[PipeTcpNet](26b94bd5-2695-36e7-6378-c3d41846be8f.htm)管道，[PipeUdpNet](5a0b810e-187e-826b-0a13-34fd5576aa0e.htm)管道，[PipeSerialPort](9d924607-3dfc-3916-b532-c5dec48fe347.htm)管道等  Get or set the current pipeline information, the pipeline type is [CommunicationPipe](ad8e903b-d289-21b0-077e-15ede8005221.htm) inheritance class, [PipeTcpNet](26b94bd5-2695-36e7-6378-c3d41846be8f.htm) pipeline, [PipeUdpNet](5a0b810e-187e-826b-0a13-34fd5576aa0e.htm) pipeline, [PipeSerialPort](9d924607-3dfc-3916-b532-c5dec48fe347.htm) pipeline, etc (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 公共属性 | [ConnectionId](ce1efb4b-fbc9-4683-163e-0a8e71e99d6b.htm) | 当前连接的唯一ID号，默认为长度20的guid码加随机数组成，方便列表管理，也可以自己指定  The unique ID number of the current connection. The default is a 20-digit guid code plus a random number. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [EnableCodeFE](971299bb-737a-a525-a122-aab17ba9e425.htm) | 获取或设置是否在每一次的报文通信时，增加"FE FE FE FE"的命令头  Get or set whether to add the command header of "FE FE FE FE" in each message communication |
| 公共属性 | [InstrumentType](ec507b3f-7e2e-1aaa-0165-6cc29562c0fe.htm) | 获取或设置仪表的类型，通常是 0x10:冷水水表 0x11:生活热水水表 0x12:直饮水水表 0x13:中水水表 0x20:热量表(热量) 0x21:热量表(冷量) 0x30:燃气表 0x40:电度表 |
| 公共属性 | [IsClearCacheBeforeRead](81dc971f-1f86-8535-ae9a-586b6590a584.htm) | 是否在发送数据前清空缓冲数据，默认是false  Whether to empty the buffer before sending data, the default is false (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 公共属性代码示例 | [LogNet](8c56f246-88ac-afb8-8c06-972fafe4f78a.htm) | 组件的日志工具，支持日志记录，只要实例化后，当前网络的基本信息，就以[DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm)等级进行输出  The component's logging tool supports logging. As long as the instantiation of the basic network information, the output will be output at [DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm) (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [PortName](f2b8e4e3-7079-5240-4ff1-14c6c9b376e6.htm) | 当前连接串口信息的端口号名称  The port name of the current connection serial port information (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 公共属性 | [ReceiveEmptyDataCount](f800fd56-7d49-87c4-985b-103a38330078.htm) | 获取或设置连续接收空的数据次数，在数据接收完成时有效，每个单位消耗的时间为[SleepTime](1e3a83ef-b6fc-bc8c-8950-ea3f6ae5c19b.htm)。  Obtain or set the number of consecutive times to receive empty data, which is valid when the data is received, and the time consumed by each unit is [SleepTime](1e3a83ef-b6fc-bc8c-8950-ea3f6ae5c19b.htm) (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 公共属性代码示例 | [ReceiveTimeOut](250560a7-332b-ce1b-5eb0-e3abe4ee6012.htm) | 获取或设置接收服务器反馈的时间，如果为负数，则不接收反馈   Gets or sets the time to receive server feedback, and if it is a negative number, does not receive feedback (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [RtsEnable](dd4057ef-88fe-f7d0-d2c4-f6e15488fe78.htm) | 获取或设置一个值，该值指示在串行通信中是否启用请求发送 (RTS) 信号。  Gets or sets a value indicating whether the request sending (RTS) signal is enabled in serial communication. (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 公共属性 | [SendBeforeHex](72aa1de4-116f-abfe-4adb-814a0be0e560.htm) | 获取或设置在发送通信报文前追加发送的字节信息，HEX格式，通常用于lora组网时，需要携带 00 00 00 02 四个字节的站地址功能。  Obtain or set the byte information sent before sending communication packets, HEX format, usually used for LORA networking, you need to carry 00 00 00 02 four-byte station address function. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [SleepTime](86a5fbb9-3456-6beb-7113-5bbad5eabbc5.htm) | 获取或设置在正式接收对方返回数据前的时候，需要休息的时间，当设置为0的时候，不需要休息。  Get or set the time required to rest before officially receiving the data from the other party. When it is set to 0, no rest is required. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [Station](4fefd98f-2903-133f-88e3-11d7f62706db.htm) | 获取或设置当前的地址域信息，是一个14个字符的BCD码，例如：14910000729011  Get or set the current address domain information, which is a 14-character BCD code, for example: 14910000729011 |
| 公共属性 | [StationMatch](c4053c8a-3cf4-4a54-cb18-134072642a84.htm) | 获取或设置是否验证匹配接收到的站号信息  Gets or sets whether to verify that the received station number information is matched |
| 受保护的属性 | [WordLength](f1a505cc-e85c-6c8b-b189-58a1ab5753f9.htm) | 一个字单位的数据表示的地址长度，西门子为2，三菱，欧姆龙，modbusTcp就为1，AB PLC无效  The address length represented by one word of data, Siemens is 2, Mitsubishi, Omron, modbusTcp is 1, AB PLC is invalid (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [ActiveDeveice](6de9a868-a602-a896-330c-e56955f8eb96.htm) | 激活设备的命令，只发送数据到设备，不等待设备数据返回  The command to activate the device, only send data to the device, do not wait for the device data to return |
| 公共方法 | [Close](d3a2ebfb-4e6b-eeda-8bc8-c0972e21b5ff.htm) | 关闭当前的串口连接  Close the current serial connection (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 受保护的方法 | [DecideWhetherQAMessage](acc40cda-e310-8537-6685-59b67eeb16ff.htm) | 决定当前的消息是否是用于问答机制返回的消息，默认直接返回 true, 实际的情况需要根据协议进行重写方法  To determine whether the current message is the message returned by the question answering mechanism, the default is true. In actual cases, the rewriting method needs to be performed according to the protocol (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [Dispose](b9dda6bf-d342-254e-689a-93fa0bb265d1.htm) | (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 受保护的方法 | [Dispose(Boolean)](0bbb6ed3-111e-069c-e9f5-8cbc3f7f7b77.htm) | (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | [ExtraAfterReadFromCoreServer](06a49a15-a3ce-5e39-28c7-2b78851eca2e.htm) | 和服务器交互完成的时候调用的方法，可以根据读写结果进行一些额外的操作，具体的操作需要根据实际的需求来重写实现  The method called when the interaction with the server is completed can perform some additional operations based on the read and write results. The specific operations need to be rewritten according to actual needs. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法代码示例 | [ExtraOnDisconnect](de9dbba2-3b2c-1c5e-53c8-4e6f8a44e2c4.htm) | 根据实际的协议选择是否重写本方法，有些协议在断开连接之前，需要发送一些报文来关闭当前的网络通道  Select whether to rewrite this method according to the actual protocol. Some protocols need to send some packets to close the current network channel before disconnecting. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法代码示例 | [ExtraOnDisconnectAsync](210f593e-16d8-77c4-1178-58fbcfe9b27b.htm) | 根据实际的协议选择是否重写本方法，有些协议在断开连接之前，需要发送一些报文来关闭当前的网络通道  Select whether to rewrite this method according to the actual protocol. Some protocols need to send some packets to close the current network channel before disconnecting. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 受保护的方法 | [GetLogTextFromBinary](67f3b2f7-957c-fa0b-8320-6799237157db.htm) | 获取当前的报文进行日志记录的时候，是否使用二进制的格式记录，默认返回 [LogMsgFormatBinary](2ed2ff87-caec-0e44-f08b-e2b015225a58.htm)，重写可以根据session对象分别返回不同记录模式  Whether to log the current packet in binary format, the default return is [LogMsgFormatBinary](2ed2ff87-caec-0e44-f08b-e2b015225a58.htm). If you want to override it, different recording modes can be returned according to session (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [GetNewNetMessage](d78644ed-1883-a2f5-91ad-741d56356d12.htm) | 获取一个新的消息对象的方法，需要在继承类里面进行重写  The method to get a new message object needs to be overridden in the inheritance class (重写 [BinaryCommunicationGetNewNetMessage](16636a1b-a6f0-9c0f-9d9c-c94e26c67459.htm).) |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法 | [GetWordLength](c27821ea-1739-8eb9-a848-23a00828fec8.htm) | 一个字单位的数据表示的地址长度，西门子为2，三菱，欧姆龙，modbusTcp就为1，AB PLC无效  The address length represented by one word of data, Siemens is 2, Mitsubishi, Omron, modbusTcp is 1, AB PLC is invalid (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 受保护的方法代码示例 | [InitializationOnConnect](4ab6989f-a5b5-05fe-8b44-bfa45824a2a7.htm) | 根据实际的协议选择是否重写本方法，有些协议在创建连接之后，需要进行一些初始化的信号握手，才能最终建立网络通道。  Whether to rewrite this method is based on the actual protocol. Some protocols require some initial signal handshake to establish a network channel after the connection is created. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法代码示例 | [InitializationOnConnectAsync](7da2af04-d4fc-5045-19a4-4b0d1722eef3.htm) | 根据实际的协议选择是否重写本方法，有些协议在创建连接之后，需要进行一些初始化的信号握手，才能最终建立网络通道。  Whether to rewrite this method is based on the actual protocol. Some protocols require some initial signal handshake to establish a network channel after the connection is created. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [IsOpen](57cc9454-4b5a-37de-df7a-7184463ad5ef.htm) | 获取一个值，指示串口是否处于打开状态  Gets a value indicating whether the serial port is open (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 受保护的方法 | [LogRevcMessage(Byte)](dab1657f-1e75-7fe1-7003-ccedfc69855d.htm) | 使用日志记录一个接收的报文信息  Logs are used to record information about a received packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [LogRevcMessage(Byte, PipeSession)](0f51ae45-575e-4228-f978-01b37d6d5df1.htm) | 使用日志记录一个接收的报文信息  Logs are used to record information about a received packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [LogSendMessage(Byte)](4f90eecc-68a2-1ef2-44ef-369f19bed3ec.htm) | 使用日志记录一个发送的报文信息  Logs are used to record information about a send packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [LogSendMessage(Byte, PipeSession)](209caadf-4f0f-40aa-04c2-ebc61f8dab13.htm) | 使用日志记录一个发送的报文信息  Logs are used to record information about a send packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法 | [Open](e6840d02-f668-0c1a-d187-4d379de79908.htm) | 打开一个新的串行端口连接  Open a new serial port connection (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 公共方法 | [PackCommandWithHeader](84add1a8-4ab7-c216-b67e-0312b85d665e.htm) | 对当前的命令进行打包处理，通常是携带命令头内容，标记当前的命令的长度信息，需要进行重写，否则默认不打包  The current command is packaged, usually carrying the content of the command header, marking the length of the current command, and it needs to be rewritten, otherwise it is not packaged by default (重写 [BinaryCommunicationPackCommandWithHeader(Byte)](d0b2cd76-c306-11fc-38c7-af8ad306cce7.htm).) |
| 公共方法 | [Read(String, UInt16)](7b7961dc-c496-8545-5d9a-eac041553faa.htm) | 根据指定的数据标识来读取相关的原始数据信息，地址标识根据手册来，从高位到地位，例如 91-1F，分割符可以任意特殊字符或是没有分隔符。  Read the relevant original data information according to the specified data identifier. The address identifier is based on the manual, from high to position, such as 91-1F. The separator can be any special character or no separator. (重写 [DeviceCommunicationRead(String, UInt16)](9dd09452-c200-bc65-75ff-69aa34b366d4.htm).) |
| 公共方法代码示例 | [ReadT](95abe064-5594-5ae7-a849-98c53f88fee0.htm) | 读取支持Hsl特性的数据内容，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明，或是官网的详细文档  Read the data content of the Hsl attribute. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadAddress](81a31692-1162-4850-5404-ed02f888f324.htm) | 读取设备的通信地址，仅支持点对点通讯的情况，返回地址域数据，例如：14910000729012  Read the communication address of the device, only support point-to-point communication, and return the address field data, for example: 149100007290 |
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
| 公共方法代码示例 | [ReadDouble(String, UInt16)](4eceecff-16be-93ac-cfc9-da43972e3f19.htm) | 读取双浮点数据的数组  Read double floating point data array (重写 [DeviceCommunicationReadDouble(String, UInt16)](4eee0f05-f861-4b48-3f06-99737fb24aee.htm).) |
| 公共方法代码示例 | [ReadDoubleAsync(String)](2cf895fa-1033-1def-f048-0e8caefaa71d.htm) | 异步读取双浮点的数据  Asynchronously read double floating point data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadDoubleAsync(String, UInt16)](51c76fdf-a6e4-f9e7-975c-ef8bf7fa1e7a.htm) | 读取双浮点数据的数组  Read double floating point data array (重写 [DeviceCommunicationReadDoubleAsync(String, UInt16)](2256d406-d4d5-da54-210b-ccef58a5e914.htm).) |
| 公共方法代码示例 | [ReadFloat(String)](bc2801ae-6600-e845-7c6c-522c9156eb0e.htm) | 读取单浮点数据  Read single floating point data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadFloat(String, UInt16)](91bc4f0a-56d6-51d8-ddda-d776174aded3.htm) | 读取单浮点精度的数组  Read single floating point array (重写 [DeviceCommunicationReadFloat(String, UInt16)](6a9fb098-b13c-313b-3b27-20d3c787945d.htm).) |
| 公共方法代码示例 | [ReadFloatAsync(String)](89f743f8-4b3a-7085-1e49-a1d7adb10052.htm) | 异步读取单浮点数据  Asynchronously read single floating point data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadFloatAsync(String, UInt16)](86291e27-0fb1-9cff-9d92-fb62f3194724.htm) | 读取单浮点精度的数组  Read single floating point array (重写 [DeviceCommunicationReadFloatAsync(String, UInt16)](cb6d5c5d-9b5c-2467-3e6b-5efa9d20cd72.htm).) |
| 公共方法 | [ReadFromCoreServer(IEnumerableByte)](044cfdd9-ec76-b89f-2f2b-5f439025006c.htm) | 将多个数据报文按顺序发到设备，并从设备接收返回的数据内容，然后拼接成一个Byte[]信息，需要重写[UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm)方法才能返回正确的结果。  Send multiple data packets to the device in sequence, and receive the returned data content from the device, and then splicing them into a Byte[] message, you need to rewrite [UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm) method to return the correct result. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServer(Byte)](010eab80-ec07-2570-5860-38ad4ef731c8.htm) | 将二进制的数据发送到管道中去，然后从管道里接收二进制的数据回来，并返回是否成功的结果对象。  Send binary data to the pipeline, and then receive binary data back from the pipeline, and return whether the success of the result object (重写 [BinaryCommunicationReadFromCoreServer(Byte)](8146dc0c-1edc-3ca7-428a-a3bdbd35e6c5.htm).) |
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
| 公共方法代码示例 | [ReadString(String, UInt16, Encoding)](d4dc35c1-5ab3-44f1-df6d-972add664fe3.htm) | 使用指定的编码，读取字符串数据  Reads string data using the specified encoding (重写 [DeviceCommunicationReadString(String, UInt16, Encoding)](b34c71a6-58f3-6bf4-0b93-471be4f2ed31.htm).) |
| 公共方法 | [ReadStringArray](c793adda-9b2c-6746-467e-b07a4f88e64a.htm) | 读取指定地址的所有的字符串数据信息，一般来说，一个地址只有一个数据 |
| 公共方法代码示例 | [ReadStringAsync(String, UInt16)](1835da81-b811-d151-0b53-907d8d87d0a8.htm) | 异步读取字符串数据，默认为最常见的ASCII编码  Asynchronously read string data, default is the most common ASCII encoding (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadStringAsync(String, UInt16, Encoding)](a6830c2f-45d7-c561-2237-f511b3d9c320.htm) | 异步使用指定的编码，读取字符串数据  Asynchronously reads string data using the specified encoding (重写 [DeviceCommunicationReadStringAsync(String, UInt16, Encoding)](091dfc67-3d25-f33e-3a55-a0228863815e.htm).) |
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
| 公共方法 | [SerialPortInni(ActionSerialPort)](3728a0f2-4553-0bc8-5c92-253c13e84f05.htm) | 根据自定义初始化方法进行初始化串口信息  Initialize the serial port information according to the custom initialization method (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 公共方法 | [SerialPortInni(String)](f82e0ba4-af2f-5165-dcde-e047c3351e74.htm) | 初始化串口信息，9600波特率，8位数据位，1位停止位，无奇偶校验  Initial serial port information, 9600 baud rate, 8 data bits, 1 stop bit, no parity (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 公共方法 | [SerialPortInni(String, Int32)](433523a8-1606-f029-50db-68ef68034046.htm) | 初始化串口信息，波特率，8位数据位，1位停止位，无奇偶校验  Initializes serial port information, baud rate, 8-bit data bit, 1-bit stop bit, no parity (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 公共方法 | [SerialPortInni(String, Int32, Int32, StopBits, Parity)](cefa7cd2-9c1f-b79d-9ea3-d4d27d7ddeae.htm) | 初始化串口信息，波特率，数据位，停止位，奇偶校验需要全部自己来指定  Start serial port information, baud rate, data bit, stop bit, parity all need to be specified (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 公共方法 | [SetDtuPipe](f765d1d7-234a-b229-9457-97e25e024dc9.htm) | 将当前的通信对象设置DTU模式，允许传入现成的管道，并返回初始化结果，如果该设备重写了握手报文，就是返回握手结果  Set the current communication object to DTU mode, allow the existing pipe to be passed in, and return the initialization result, if the device rewrites the handshake packet, the handshake result is returned (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [SetDtuPipeAsync](9d1fd12d-b804-3828-a627-a2c0eef1d23f.htm) | 将当前的通信对象设置DTU模式，允许传入现成的管道，并返回初始化结果，如果该设备重写了握手报文，就是返回握手结果  Set the current communication object to DTU mode, allow the existing pipe to be passed in, and return the initialization result, if the device rewrites the handshake packet, the handshake result is returned (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ToString](8bff78ed-c0ad-15e4-e369-ad46dbff3241.htm) | (重写 [DeviceSerialPortToString](774e7804-97ce-fb07-3764-39fb1ca2e1e7.htm).) |
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
| 公共方法 | [Write(String, Boolean)](06d0b9af-7f28-0b9c-098f-4da28e69504a.htm) | 写入单个的Boolean数据，返回是否成功  Write a single Boolean data, and return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [Write(String, Boolean)](5b046c82-3d8e-550b-fdf3-41ca5558169b.htm) | 批量写入Boolean数组数据，返回是否成功  Batch write Boolean array data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
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
| 公共方法 | [Write(String, Byte)](4036166e-1056-e46c-d98d-b5e3f1e761b0.htm) | 根据指定的数据标识来写入相关的原始数据信息，地址标识根据手册来，从高位到地位，例如 90-1F，分割符可以任意特殊字符或是没有分隔符。  Read the relevant original data information according to the specified data identifier. The address identifier is based on the manual, from high to position, such as 90-1F. The separator can be any special character or no separator. (重写 [DeviceCommunicationWrite(String, Byte)](b5536dd6-ccde-6883-ac27-bff9e4b37806.htm).) |
| 公共方法代码示例 | [Write(String, String, Int32)](02fe2ca9-49a3-ede4-d1af-8c6da1f644da.htm) | 写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为ASCII  Write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, it is filled with 0 and the encoding is ASCII. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, String, Encoding)](7b602661-902d-2b70-1fdb-ba87ecd307bf.htm) | 写入字符串信息，需要指定的编码信息  Write string information, need to specify the encoding information (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, String, Int32, Encoding)](e8cb30d9-4c5e-6c68-b9af-1b1b6c2dc714.htm) | 写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为指定的编码信息  Write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, then the operation is complemented with 0 , you should specified the encoding information (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteT(T)](9f48853f-fbeb-c37d-25a2-1c76d1dee407.htm) | 写入支持Hsl特性的数据，返回是否写入成功，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Write data that supports the Hsl attribute, and return whether the write was successful. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WriteAddress](37e7ec52-2c99-b84f-3062-f23d3541a67a.htm) | 写入设备的地址域信息，仅支持点对点通讯的情况，需要指定地址域信息，例如：14910000729011  Write the address domain information of the device, only support point-to-point communication, you need to specify the address domain information, for example: 14910000729011 |
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

[HslCommunication.Instrument.CJT 命名空间](32a7bbe4-7443-0a40-8e16-c1df11e7a933.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## CJT188 构造函数 

[原文連結](http://api.hslcommunication.cn/html/72a09742-43f7-3f9c-7d8f-0872d228bb0f.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.CJT](../html/32a7bbe4-7443-0a40-8e16-c1df11e7a933.htm "HslCommunication.Instrument.CJT")

[CJT188 类](../html/7845b81d-5321-3a2e-346d-ac4fb15e3f2a.htm "CJT188 类")

[CJT188 构造函数](../html/72a09742-43f7-3f9c-7d8f-0872d228bb0f.htm "CJT188 构造函数 ")

[CJT188 属性](../html/be37af45-b2d6-9a32-9d71-e9b361f46ab2.htm "CJT188 属性")

[CJT188 方法](../html/04972e44-8c91-c7ca-0ab7-28b426216410.htm "CJT188 方法")

[CJT188 字段](../html/e0c4f1ed-6183-52ba-c1f5-cdf5993cd2cc.htm "CJT188 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CJT188 构造函数 |

指定地址域来实例化一个对象，地址域是一个14个字符的BCD码，例如：14910000729011  
Specify the address field, to instantiate an object, which address field is a 14-character BCD code, for example: 14910000729011

**命名空间：**
 [HslCommunication.Instrument.CJT](32a7bbe4-7443-0a40-8e16-c1df11e7a933.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public CJT188(
	string station
)
```

```
Public Sub New ( 
	station As String
)
```

```
public:
CJT188(
	String^ station
)
```

```
new : 
        station : string -> CJT188
```

#### 参数

station
:   类型：SystemString  
    设备的地址信息，是一个14字符的BCD码

![](../icons/SectionExpanded.png)参见

#### 引用

[CJT188 类](7845b81d-5321-3a2e-346d-ac4fb15e3f2a.htm)

[HslCommunication.Instrument.CJT 命名空间](32a7bbe4-7443-0a40-8e16-c1df11e7a933.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## CJT188 属性

[原文連結](http://api.hslcommunication.cn/html/be37af45-b2d6-9a32-9d71-e9b361f46ab2.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.CJT](../html/32a7bbe4-7443-0a40-8e16-c1df11e7a933.htm "HslCommunication.Instrument.CJT")

[CJT188 类](../html/7845b81d-5321-3a2e-346d-ac4fb15e3f2a.htm "CJT188 类")

[CJT188 属性](../html/be37af45-b2d6-9a32-9d71-e9b361f46ab2.htm "CJT188 属性")

[EnableCodeFE 属性](../html/971299bb-737a-a525-a122-aab17ba9e425.htm "EnableCodeFE 属性 ")

[InstrumentType 属性](../html/ec507b3f-7e2e-1aaa-0165-6cc29562c0fe.htm "InstrumentType 属性 ")

[Station 属性](../html/4fefd98f-2903-133f-88e3-11d7f62706db.htm "Station 属性 ")

[StationMatch 属性](../html/c4053c8a-3cf4-4a54-cb18-134072642a84.htm "StationMatch 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CJT188 属性 |

[CJT188](7845b81d-5321-3a2e-346d-ac4fb15e3f2a.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [BaudRate](5554f563-5e29-699c-435b-c8e7ad5fa0e4.htm) | 当前连接串口信息的波特率  Baud rate of current connection serial port information (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 公共属性代码示例 | [ByteTransform](f35118c9-3691-f752-ae6f-2c5549714faa.htm) | 当前的数据变换机制，当你需要从字节数据转换类型数据的时候需要。  The current data transformation mechanism is required when you need to convert type data from byte data. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共属性 | [CommunicationPipe](a4fcb477-64d2-8da7-c712-687f0261755b.htm) | 获取或设置当前的管道信息，管道类型为[CommunicationPipe](ad8e903b-d289-21b0-077e-15ede8005221.htm)的继承类，内置了[PipeTcpNet](26b94bd5-2695-36e7-6378-c3d41846be8f.htm)管道，[PipeUdpNet](5a0b810e-187e-826b-0a13-34fd5576aa0e.htm)管道，[PipeSerialPort](9d924607-3dfc-3916-b532-c5dec48fe347.htm)管道等  Get or set the current pipeline information, the pipeline type is [CommunicationPipe](ad8e903b-d289-21b0-077e-15ede8005221.htm) inheritance class, [PipeTcpNet](26b94bd5-2695-36e7-6378-c3d41846be8f.htm) pipeline, [PipeUdpNet](5a0b810e-187e-826b-0a13-34fd5576aa0e.htm) pipeline, [PipeSerialPort](9d924607-3dfc-3916-b532-c5dec48fe347.htm) pipeline, etc (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 公共属性 | [ConnectionId](ce1efb4b-fbc9-4683-163e-0a8e71e99d6b.htm) | 当前连接的唯一ID号，默认为长度20的guid码加随机数组成，方便列表管理，也可以自己指定  The unique ID number of the current connection. The default is a 20-digit guid code plus a random number. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [EnableCodeFE](971299bb-737a-a525-a122-aab17ba9e425.htm) | 获取或设置是否在每一次的报文通信时，增加"FE FE FE FE"的命令头  Get or set whether to add the command header of "FE FE FE FE" in each message communication |
| 公共属性 | [InstrumentType](ec507b3f-7e2e-1aaa-0165-6cc29562c0fe.htm) | 获取或设置仪表的类型，通常是 0x10:冷水水表 0x11:生活热水水表 0x12:直饮水水表 0x13:中水水表 0x20:热量表(热量) 0x21:热量表(冷量) 0x30:燃气表 0x40:电度表 |
| 公共属性 | [IsClearCacheBeforeRead](81dc971f-1f86-8535-ae9a-586b6590a584.htm) | 是否在发送数据前清空缓冲数据，默认是false  Whether to empty the buffer before sending data, the default is false (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 公共属性代码示例 | [LogNet](8c56f246-88ac-afb8-8c06-972fafe4f78a.htm) | 组件的日志工具，支持日志记录，只要实例化后，当前网络的基本信息，就以[DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm)等级进行输出  The component's logging tool supports logging. As long as the instantiation of the basic network information, the output will be output at [DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm) (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [PortName](f2b8e4e3-7079-5240-4ff1-14c6c9b376e6.htm) | 当前连接串口信息的端口号名称  The port name of the current connection serial port information (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 公共属性 | [ReceiveEmptyDataCount](f800fd56-7d49-87c4-985b-103a38330078.htm) | 获取或设置连续接收空的数据次数，在数据接收完成时有效，每个单位消耗的时间为[SleepTime](1e3a83ef-b6fc-bc8c-8950-ea3f6ae5c19b.htm)。  Obtain or set the number of consecutive times to receive empty data, which is valid when the data is received, and the time consumed by each unit is [SleepTime](1e3a83ef-b6fc-bc8c-8950-ea3f6ae5c19b.htm) (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 公共属性代码示例 | [ReceiveTimeOut](250560a7-332b-ce1b-5eb0-e3abe4ee6012.htm) | 获取或设置接收服务器反馈的时间，如果为负数，则不接收反馈   Gets or sets the time to receive server feedback, and if it is a negative number, does not receive feedback (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [RtsEnable](dd4057ef-88fe-f7d0-d2c4-f6e15488fe78.htm) | 获取或设置一个值，该值指示在串行通信中是否启用请求发送 (RTS) 信号。  Gets or sets a value indicating whether the request sending (RTS) signal is enabled in serial communication. (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 公共属性 | [SendBeforeHex](72aa1de4-116f-abfe-4adb-814a0be0e560.htm) | 获取或设置在发送通信报文前追加发送的字节信息，HEX格式，通常用于lora组网时，需要携带 00 00 00 02 四个字节的站地址功能。  Obtain or set the byte information sent before sending communication packets, HEX format, usually used for LORA networking, you need to carry 00 00 00 02 four-byte station address function. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [SleepTime](86a5fbb9-3456-6beb-7113-5bbad5eabbc5.htm) | 获取或设置在正式接收对方返回数据前的时候，需要休息的时间，当设置为0的时候，不需要休息。  Get or set the time required to rest before officially receiving the data from the other party. When it is set to 0, no rest is required. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [Station](4fefd98f-2903-133f-88e3-11d7f62706db.htm) | 获取或设置当前的地址域信息，是一个14个字符的BCD码，例如：14910000729011  Get or set the current address domain information, which is a 14-character BCD code, for example: 14910000729011 |
| 公共属性 | [StationMatch](c4053c8a-3cf4-4a54-cb18-134072642a84.htm) | 获取或设置是否验证匹配接收到的站号信息  Gets or sets whether to verify that the received station number information is matched |
| 受保护的属性 | [WordLength](f1a505cc-e85c-6c8b-b189-58a1ab5753f9.htm) | 一个字单位的数据表示的地址长度，西门子为2，三菱，欧姆龙，modbusTcp就为1，AB PLC无效  The address length represented by one word of data, Siemens is 2, Mitsubishi, Omron, modbusTcp is 1, AB PLC is invalid (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[CJT188 类](7845b81d-5321-3a2e-346d-ac4fb15e3f2a.htm)

[HslCommunication.Instrument.CJT 命名空间](32a7bbe4-7443-0a40-8e16-c1df11e7a933.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## EnableCodeFE 属性 

[原文連結](http://api.hslcommunication.cn/html/971299bb-737a-a525-a122-aab17ba9e425.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.CJT](../html/32a7bbe4-7443-0a40-8e16-c1df11e7a933.htm "HslCommunication.Instrument.CJT")

[CJT188 类](../html/7845b81d-5321-3a2e-346d-ac4fb15e3f2a.htm "CJT188 类")

[CJT188 属性](../html/be37af45-b2d6-9a32-9d71-e9b361f46ab2.htm "CJT188 属性")

[EnableCodeFE 属性](../html/971299bb-737a-a525-a122-aab17ba9e425.htm "EnableCodeFE 属性 ")

[InstrumentType 属性](../html/ec507b3f-7e2e-1aaa-0165-6cc29562c0fe.htm "InstrumentType 属性 ")

[Station 属性](../html/4fefd98f-2903-133f-88e3-11d7f62706db.htm "Station 属性 ")

[StationMatch 属性](../html/c4053c8a-3cf4-4a54-cb18-134072642a84.htm "StationMatch 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CJT188EnableCodeFE 属性 |

获取或设置是否在每一次的报文通信时，增加"FE FE FE FE"的命令头  
Get or set whether to add the command header of "FE FE FE FE" in each message communication

**命名空间：**
 [HslCommunication.Instrument.CJT](32a7bbe4-7443-0a40-8e16-c1df11e7a933.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public bool EnableCodeFE { get; set; }
```

```
Public Property EnableCodeFE As Boolean
	Get
	Set
```

```
public:
virtual property bool EnableCodeFE {
	bool get () sealed;
	void set (bool value) sealed;
}
```

```
abstract EnableCodeFE : bool with get, set
override EnableCodeFE : bool with get, set
```

#### 属性值

类型：Boolean

#### 实现

[ICjt188EnableCodeFE](2420a2be-c28d-4ddf-63e7-05e2eb050af8.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[CJT188 类](7845b81d-5321-3a2e-346d-ac4fb15e3f2a.htm)

[HslCommunication.Instrument.CJT 命名空间](32a7bbe4-7443-0a40-8e16-c1df11e7a933.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## InstrumentType 属性 

[原文連結](http://api.hslcommunication.cn/html/ec507b3f-7e2e-1aaa-0165-6cc29562c0fe.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.CJT](../html/32a7bbe4-7443-0a40-8e16-c1df11e7a933.htm "HslCommunication.Instrument.CJT")

[CJT188 类](../html/7845b81d-5321-3a2e-346d-ac4fb15e3f2a.htm "CJT188 类")

[CJT188 属性](../html/be37af45-b2d6-9a32-9d71-e9b361f46ab2.htm "CJT188 属性")

[EnableCodeFE 属性](../html/971299bb-737a-a525-a122-aab17ba9e425.htm "EnableCodeFE 属性 ")

[InstrumentType 属性](../html/ec507b3f-7e2e-1aaa-0165-6cc29562c0fe.htm "InstrumentType 属性 ")

[Station 属性](../html/4fefd98f-2903-133f-88e3-11d7f62706db.htm "Station 属性 ")

[StationMatch 属性](../html/c4053c8a-3cf4-4a54-cb18-134072642a84.htm "StationMatch 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CJT188InstrumentType 属性 |

获取或设置仪表的类型，通常是 0x10:冷水水表 0x11:生活热水水表 0x12:直饮水水表 0x13:中水水表 0x20:热量表(热量) 0x21:热量表(冷量) 0x30:燃气表 0x40:电度表

**命名空间：**
 [HslCommunication.Instrument.CJT](32a7bbe4-7443-0a40-8e16-c1df11e7a933.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public byte InstrumentType { get; set; }
```

```
Public Property InstrumentType As Byte
	Get
	Set
```

```
public:
virtual property unsigned char InstrumentType {
	unsigned char get () sealed;
	void set (unsigned char value) sealed;
}
```

```
abstract InstrumentType : byte with get, set
override InstrumentType : byte with get, set
```

#### 属性值

类型：Byte

#### 实现

[ICjt188InstrumentType](54d205f6-517b-1b53-bb21-afa097006156.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[CJT188 类](7845b81d-5321-3a2e-346d-ac4fb15e3f2a.htm)

[HslCommunication.Instrument.CJT 命名空间](32a7bbe4-7443-0a40-8e16-c1df11e7a933.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Station 属性 

[原文連結](http://api.hslcommunication.cn/html/4fefd98f-2903-133f-88e3-11d7f62706db.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.CJT](../html/32a7bbe4-7443-0a40-8e16-c1df11e7a933.htm "HslCommunication.Instrument.CJT")

[CJT188 类](../html/7845b81d-5321-3a2e-346d-ac4fb15e3f2a.htm "CJT188 类")

[CJT188 属性](../html/be37af45-b2d6-9a32-9d71-e9b361f46ab2.htm "CJT188 属性")

[EnableCodeFE 属性](../html/971299bb-737a-a525-a122-aab17ba9e425.htm "EnableCodeFE 属性 ")

[InstrumentType 属性](../html/ec507b3f-7e2e-1aaa-0165-6cc29562c0fe.htm "InstrumentType 属性 ")

[Station 属性](../html/4fefd98f-2903-133f-88e3-11d7f62706db.htm "Station 属性 ")

[StationMatch 属性](../html/c4053c8a-3cf4-4a54-cb18-134072642a84.htm "StationMatch 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CJT188Station 属性 |

获取或设置当前的地址域信息，是一个14个字符的BCD码，例如：14910000729011  
Get or set the current address domain information, which is a 14-character BCD code, for example: 14910000729011

**命名空间：**
 [HslCommunication.Instrument.CJT](32a7bbe4-7443-0a40-8e16-c1df11e7a933.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public string Station { get; set; }
```

```
Public Property Station As String
	Get
	Set
```

```
public:
virtual property String^ Station {
	String^ get () sealed;
	void set (String^ value) sealed;
}
```

```
abstract Station : string with get, set
override Station : string with get, set
```

#### 属性值

类型：String

#### 实现

[ICjt188Station](36f1a758-e626-a565-2d13-3a0cdea61613.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[CJT188 类](7845b81d-5321-3a2e-346d-ac4fb15e3f2a.htm)

[HslCommunication.Instrument.CJT 命名空间](32a7bbe4-7443-0a40-8e16-c1df11e7a933.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## StationMatch 属性 

[原文連結](http://api.hslcommunication.cn/html/c4053c8a-3cf4-4a54-cb18-134072642a84.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.CJT](../html/32a7bbe4-7443-0a40-8e16-c1df11e7a933.htm "HslCommunication.Instrument.CJT")

[CJT188 类](../html/7845b81d-5321-3a2e-346d-ac4fb15e3f2a.htm "CJT188 类")

[CJT188 属性](../html/be37af45-b2d6-9a32-9d71-e9b361f46ab2.htm "CJT188 属性")

[EnableCodeFE 属性](../html/971299bb-737a-a525-a122-aab17ba9e425.htm "EnableCodeFE 属性 ")

[InstrumentType 属性](../html/ec507b3f-7e2e-1aaa-0165-6cc29562c0fe.htm "InstrumentType 属性 ")

[Station 属性](../html/4fefd98f-2903-133f-88e3-11d7f62706db.htm "Station 属性 ")

[StationMatch 属性](../html/c4053c8a-3cf4-4a54-cb18-134072642a84.htm "StationMatch 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CJT188StationMatch 属性 |

获取或设置是否验证匹配接收到的站号信息  
Gets or sets whether to verify that the received station number information is matched

**命名空间：**
 [HslCommunication.Instrument.CJT](32a7bbe4-7443-0a40-8e16-c1df11e7a933.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public bool StationMatch { get; set; }
```

```
Public Property StationMatch As Boolean
	Get
	Set
```

```
public:
property bool StationMatch {
	bool get ();
	void set (bool value);
}
```

```
member StationMatch : bool with get, set
```

#### 属性值

类型：Boolean

![](../icons/SectionExpanded.png)参见

#### 引用

[CJT188 类](7845b81d-5321-3a2e-346d-ac4fb15e3f2a.htm)

[HslCommunication.Instrument.CJT 命名空间](32a7bbe4-7443-0a40-8e16-c1df11e7a933.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## CJT188 方法

[原文連結](http://api.hslcommunication.cn/html/04972e44-8c91-c7ca-0ab7-28b426216410.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.CJT](../html/32a7bbe4-7443-0a40-8e16-c1df11e7a933.htm "HslCommunication.Instrument.CJT")

[CJT188 类](../html/7845b81d-5321-3a2e-346d-ac4fb15e3f2a.htm "CJT188 类")

[CJT188 方法](../html/04972e44-8c91-c7ca-0ab7-28b426216410.htm "CJT188 方法")

[ActiveDeveice 方法](../html/6de9a868-a602-a896-330c-e56955f8eb96.htm "ActiveDeveice 方法 ")

[GetNewNetMessage 方法](../html/d78644ed-1883-a2f5-91ad-741d56356d12.htm "GetNewNetMessage 方法 ")

[PackCommandWithHeader 方法](../html/84add1a8-4ab7-c216-b67e-0312b85d665e.htm "PackCommandWithHeader 方法 ")

[Read 方法](../html/1f6fea44-369c-71f2-4883-631528d6f3e5.htm "Read 方法 ")

[ReadAddress 方法](../html/81a31692-1162-4850-5404-ed02f888f324.htm "ReadAddress 方法 ")

[ReadDouble 方法](../html/b5fa76a1-104a-bdb6-c6d6-7472db8f735f.htm "ReadDouble 方法 ")

[ReadDoubleAsync 方法](../html/e53e8ff4-31bd-336c-4218-2c8fc79855c7.htm "ReadDoubleAsync 方法 ")

[ReadFloat 方法](../html/a5eeaaa8-4c0b-8d81-4420-3c7ce42a39f2.htm "ReadFloat 方法 ")

[ReadFloatAsync 方法](../html/7dd8d20e-35b5-8258-4550-df204a38e81f.htm "ReadFloatAsync 方法 ")

[ReadFromCoreServer 方法](../html/42325ee8-6793-9ad0-6b53-58f5487de0b1.htm "ReadFromCoreServer 方法 ")

[ReadString 方法](../html/e0ea9f3d-4c56-6078-2d10-e14579abb235.htm "ReadString 方法 ")

[ReadStringArray 方法](../html/c793adda-9b2c-6746-467e-b07a4f88e64a.htm "ReadStringArray 方法 ")

[ReadStringAsync 方法](../html/947017c8-6aa1-d267-3e9d-3b880b5021be.htm "ReadStringAsync 方法 ")

[ToString 方法](../html/8bff78ed-c0ad-15e4-e369-ad46dbff3241.htm "ToString 方法 ")

[Write 方法](../html/efa8452a-5af5-3ac6-b3fd-8f9c33845943.htm "Write 方法 ")

[WriteAddress 方法](../html/37e7ec52-2c99-b84f-3062-f23d3541a67a.htm "WriteAddress 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CJT188 方法 |

[CJT188](7845b81d-5321-3a2e-346d-ac4fb15e3f2a.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [ActiveDeveice](6de9a868-a602-a896-330c-e56955f8eb96.htm) | 激活设备的命令，只发送数据到设备，不等待设备数据返回  The command to activate the device, only send data to the device, do not wait for the device data to return |
| 公共方法 | [Close](d3a2ebfb-4e6b-eeda-8bc8-c0972e21b5ff.htm) | 关闭当前的串口连接  Close the current serial connection (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 受保护的方法 | [DecideWhetherQAMessage](acc40cda-e310-8537-6685-59b67eeb16ff.htm) | 决定当前的消息是否是用于问答机制返回的消息，默认直接返回 true, 实际的情况需要根据协议进行重写方法  To determine whether the current message is the message returned by the question answering mechanism, the default is true. In actual cases, the rewriting method needs to be performed according to the protocol (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [Dispose](b9dda6bf-d342-254e-689a-93fa0bb265d1.htm) | (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 受保护的方法 | [Dispose(Boolean)](0bbb6ed3-111e-069c-e9f5-8cbc3f7f7b77.htm) | (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | [ExtraAfterReadFromCoreServer](06a49a15-a3ce-5e39-28c7-2b78851eca2e.htm) | 和服务器交互完成的时候调用的方法，可以根据读写结果进行一些额外的操作，具体的操作需要根据实际的需求来重写实现  The method called when the interaction with the server is completed can perform some additional operations based on the read and write results. The specific operations need to be rewritten according to actual needs. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法代码示例 | [ExtraOnDisconnect](de9dbba2-3b2c-1c5e-53c8-4e6f8a44e2c4.htm) | 根据实际的协议选择是否重写本方法，有些协议在断开连接之前，需要发送一些报文来关闭当前的网络通道  Select whether to rewrite this method according to the actual protocol. Some protocols need to send some packets to close the current network channel before disconnecting. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法代码示例 | [ExtraOnDisconnectAsync](210f593e-16d8-77c4-1178-58fbcfe9b27b.htm) | 根据实际的协议选择是否重写本方法，有些协议在断开连接之前，需要发送一些报文来关闭当前的网络通道  Select whether to rewrite this method according to the actual protocol. Some protocols need to send some packets to close the current network channel before disconnecting. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 受保护的方法 | [GetLogTextFromBinary](67f3b2f7-957c-fa0b-8320-6799237157db.htm) | 获取当前的报文进行日志记录的时候，是否使用二进制的格式记录，默认返回 [LogMsgFormatBinary](2ed2ff87-caec-0e44-f08b-e2b015225a58.htm)，重写可以根据session对象分别返回不同记录模式  Whether to log the current packet in binary format, the default return is [LogMsgFormatBinary](2ed2ff87-caec-0e44-f08b-e2b015225a58.htm). If you want to override it, different recording modes can be returned according to session (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [GetNewNetMessage](d78644ed-1883-a2f5-91ad-741d56356d12.htm) | 获取一个新的消息对象的方法，需要在继承类里面进行重写  The method to get a new message object needs to be overridden in the inheritance class (重写 [BinaryCommunicationGetNewNetMessage](16636a1b-a6f0-9c0f-9d9c-c94e26c67459.htm).) |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法 | [GetWordLength](c27821ea-1739-8eb9-a848-23a00828fec8.htm) | 一个字单位的数据表示的地址长度，西门子为2，三菱，欧姆龙，modbusTcp就为1，AB PLC无效  The address length represented by one word of data, Siemens is 2, Mitsubishi, Omron, modbusTcp is 1, AB PLC is invalid (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 受保护的方法代码示例 | [InitializationOnConnect](4ab6989f-a5b5-05fe-8b44-bfa45824a2a7.htm) | 根据实际的协议选择是否重写本方法，有些协议在创建连接之后，需要进行一些初始化的信号握手，才能最终建立网络通道。  Whether to rewrite this method is based on the actual protocol. Some protocols require some initial signal handshake to establish a network channel after the connection is created. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法代码示例 | [InitializationOnConnectAsync](7da2af04-d4fc-5045-19a4-4b0d1722eef3.htm) | 根据实际的协议选择是否重写本方法，有些协议在创建连接之后，需要进行一些初始化的信号握手，才能最终建立网络通道。  Whether to rewrite this method is based on the actual protocol. Some protocols require some initial signal handshake to establish a network channel after the connection is created. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [IsOpen](57cc9454-4b5a-37de-df7a-7184463ad5ef.htm) | 获取一个值，指示串口是否处于打开状态  Gets a value indicating whether the serial port is open (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 受保护的方法 | [LogRevcMessage(Byte)](dab1657f-1e75-7fe1-7003-ccedfc69855d.htm) | 使用日志记录一个接收的报文信息  Logs are used to record information about a received packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [LogRevcMessage(Byte, PipeSession)](0f51ae45-575e-4228-f978-01b37d6d5df1.htm) | 使用日志记录一个接收的报文信息  Logs are used to record information about a received packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [LogSendMessage(Byte)](4f90eecc-68a2-1ef2-44ef-369f19bed3ec.htm) | 使用日志记录一个发送的报文信息  Logs are used to record information about a send packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [LogSendMessage(Byte, PipeSession)](209caadf-4f0f-40aa-04c2-ebc61f8dab13.htm) | 使用日志记录一个发送的报文信息  Logs are used to record information about a send packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法 | [Open](e6840d02-f668-0c1a-d187-4d379de79908.htm) | 打开一个新的串行端口连接  Open a new serial port connection (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 公共方法 | [PackCommandWithHeader](84add1a8-4ab7-c216-b67e-0312b85d665e.htm) | 对当前的命令进行打包处理，通常是携带命令头内容，标记当前的命令的长度信息，需要进行重写，否则默认不打包  The current command is packaged, usually carrying the content of the command header, marking the length of the current command, and it needs to be rewritten, otherwise it is not packaged by default (重写 [BinaryCommunicationPackCommandWithHeader(Byte)](d0b2cd76-c306-11fc-38c7-af8ad306cce7.htm).) |
| 公共方法 | [Read(String, UInt16)](7b7961dc-c496-8545-5d9a-eac041553faa.htm) | 根据指定的数据标识来读取相关的原始数据信息，地址标识根据手册来，从高位到地位，例如 91-1F，分割符可以任意特殊字符或是没有分隔符。  Read the relevant original data information according to the specified data identifier. The address identifier is based on the manual, from high to position, such as 91-1F. The separator can be any special character or no separator. (重写 [DeviceCommunicationRead(String, UInt16)](9dd09452-c200-bc65-75ff-69aa34b366d4.htm).) |
| 公共方法代码示例 | [ReadT](95abe064-5594-5ae7-a849-98c53f88fee0.htm) | 读取支持Hsl特性的数据内容，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明，或是官网的详细文档  Read the data content of the Hsl attribute. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadAddress](81a31692-1162-4850-5404-ed02f888f324.htm) | 读取设备的通信地址，仅支持点对点通讯的情况，返回地址域数据，例如：14910000729012  Read the communication address of the device, only support point-to-point communication, and return the address field data, for example: 149100007290 |
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
| 公共方法代码示例 | [ReadDouble(String, UInt16)](4eceecff-16be-93ac-cfc9-da43972e3f19.htm) | 读取双浮点数据的数组  Read double floating point data array (重写 [DeviceCommunicationReadDouble(String, UInt16)](4eee0f05-f861-4b48-3f06-99737fb24aee.htm).) |
| 公共方法代码示例 | [ReadDoubleAsync(String)](2cf895fa-1033-1def-f048-0e8caefaa71d.htm) | 异步读取双浮点的数据  Asynchronously read double floating point data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadDoubleAsync(String, UInt16)](51c76fdf-a6e4-f9e7-975c-ef8bf7fa1e7a.htm) | 读取双浮点数据的数组  Read double floating point data array (重写 [DeviceCommunicationReadDoubleAsync(String, UInt16)](2256d406-d4d5-da54-210b-ccef58a5e914.htm).) |
| 公共方法代码示例 | [ReadFloat(String)](bc2801ae-6600-e845-7c6c-522c9156eb0e.htm) | 读取单浮点数据  Read single floating point data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadFloat(String, UInt16)](91bc4f0a-56d6-51d8-ddda-d776174aded3.htm) | 读取单浮点精度的数组  Read single floating point array (重写 [DeviceCommunicationReadFloat(String, UInt16)](6a9fb098-b13c-313b-3b27-20d3c787945d.htm).) |
| 公共方法代码示例 | [ReadFloatAsync(String)](89f743f8-4b3a-7085-1e49-a1d7adb10052.htm) | 异步读取单浮点数据  Asynchronously read single floating point data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadFloatAsync(String, UInt16)](86291e27-0fb1-9cff-9d92-fb62f3194724.htm) | 读取单浮点精度的数组  Read single floating point array (重写 [DeviceCommunicationReadFloatAsync(String, UInt16)](cb6d5c5d-9b5c-2467-3e6b-5efa9d20cd72.htm).) |
| 公共方法 | [ReadFromCoreServer(IEnumerableByte)](044cfdd9-ec76-b89f-2f2b-5f439025006c.htm) | 将多个数据报文按顺序发到设备，并从设备接收返回的数据内容，然后拼接成一个Byte[]信息，需要重写[UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm)方法才能返回正确的结果。  Send multiple data packets to the device in sequence, and receive the returned data content from the device, and then splicing them into a Byte[] message, you need to rewrite [UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm) method to return the correct result. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServer(Byte)](010eab80-ec07-2570-5860-38ad4ef731c8.htm) | 将二进制的数据发送到管道中去，然后从管道里接收二进制的数据回来，并返回是否成功的结果对象。  Send binary data to the pipeline, and then receive binary data back from the pipeline, and return whether the success of the result object (重写 [BinaryCommunicationReadFromCoreServer(Byte)](8146dc0c-1edc-3ca7-428a-a3bdbd35e6c5.htm).) |
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
| 公共方法代码示例 | [ReadString(String, UInt16, Encoding)](d4dc35c1-5ab3-44f1-df6d-972add664fe3.htm) | 使用指定的编码，读取字符串数据  Reads string data using the specified encoding (重写 [DeviceCommunicationReadString(String, UInt16, Encoding)](b34c71a6-58f3-6bf4-0b93-471be4f2ed31.htm).) |
| 公共方法 | [ReadStringArray](c793adda-9b2c-6746-467e-b07a4f88e64a.htm) | 读取指定地址的所有的字符串数据信息，一般来说，一个地址只有一个数据 |
| 公共方法代码示例 | [ReadStringAsync(String, UInt16)](1835da81-b811-d151-0b53-907d8d87d0a8.htm) | 异步读取字符串数据，默认为最常见的ASCII编码  Asynchronously read string data, default is the most common ASCII encoding (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadStringAsync(String, UInt16, Encoding)](a6830c2f-45d7-c561-2237-f511b3d9c320.htm) | 异步使用指定的编码，读取字符串数据  Asynchronously reads string data using the specified encoding (重写 [DeviceCommunicationReadStringAsync(String, UInt16, Encoding)](091dfc67-3d25-f33e-3a55-a0228863815e.htm).) |
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
| 公共方法 | [SerialPortInni(ActionSerialPort)](3728a0f2-4553-0bc8-5c92-253c13e84f05.htm) | 根据自定义初始化方法进行初始化串口信息  Initialize the serial port information according to the custom initialization method (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 公共方法 | [SerialPortInni(String)](f82e0ba4-af2f-5165-dcde-e047c3351e74.htm) | 初始化串口信息，9600波特率，8位数据位，1位停止位，无奇偶校验  Initial serial port information, 9600 baud rate, 8 data bits, 1 stop bit, no parity (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 公共方法 | [SerialPortInni(String, Int32)](433523a8-1606-f029-50db-68ef68034046.htm) | 初始化串口信息，波特率，8位数据位，1位停止位，无奇偶校验  Initializes serial port information, baud rate, 8-bit data bit, 1-bit stop bit, no parity (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 公共方法 | [SerialPortInni(String, Int32, Int32, StopBits, Parity)](cefa7cd2-9c1f-b79d-9ea3-d4d27d7ddeae.htm) | 初始化串口信息，波特率，数据位，停止位，奇偶校验需要全部自己来指定  Start serial port information, baud rate, data bit, stop bit, parity all need to be specified (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 公共方法 | [SetDtuPipe](f765d1d7-234a-b229-9457-97e25e024dc9.htm) | 将当前的通信对象设置DTU模式，允许传入现成的管道，并返回初始化结果，如果该设备重写了握手报文，就是返回握手结果  Set the current communication object to DTU mode, allow the existing pipe to be passed in, and return the initialization result, if the device rewrites the handshake packet, the handshake result is returned (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [SetDtuPipeAsync](9d1fd12d-b804-3828-a627-a2c0eef1d23f.htm) | 将当前的通信对象设置DTU模式，允许传入现成的管道，并返回初始化结果，如果该设备重写了握手报文，就是返回握手结果  Set the current communication object to DTU mode, allow the existing pipe to be passed in, and return the initialization result, if the device rewrites the handshake packet, the handshake result is returned (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ToString](8bff78ed-c0ad-15e4-e369-ad46dbff3241.htm) | (重写 [DeviceSerialPortToString](774e7804-97ce-fb07-3764-39fb1ca2e1e7.htm).) |
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
| 公共方法 | [Write(String, Boolean)](06d0b9af-7f28-0b9c-098f-4da28e69504a.htm) | 写入单个的Boolean数据，返回是否成功  Write a single Boolean data, and return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [Write(String, Boolean)](5b046c82-3d8e-550b-fdf3-41ca5558169b.htm) | 批量写入Boolean数组数据，返回是否成功  Batch write Boolean array data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
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
| 公共方法 | [Write(String, Byte)](4036166e-1056-e46c-d98d-b5e3f1e761b0.htm) | 根据指定的数据标识来写入相关的原始数据信息，地址标识根据手册来，从高位到地位，例如 90-1F，分割符可以任意特殊字符或是没有分隔符。  Read the relevant original data information according to the specified data identifier. The address identifier is based on the manual, from high to position, such as 90-1F. The separator can be any special character or no separator. (重写 [DeviceCommunicationWrite(String, Byte)](b5536dd6-ccde-6883-ac27-bff9e4b37806.htm).) |
| 公共方法代码示例 | [Write(String, String, Int32)](02fe2ca9-49a3-ede4-d1af-8c6da1f644da.htm) | 写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为ASCII  Write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, it is filled with 0 and the encoding is ASCII. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, String, Encoding)](7b602661-902d-2b70-1fdb-ba87ecd307bf.htm) | 写入字符串信息，需要指定的编码信息  Write string information, need to specify the encoding information (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, String, Int32, Encoding)](e8cb30d9-4c5e-6c68-b9af-1b1b6c2dc714.htm) | 写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为指定的编码信息  Write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, then the operation is complemented with 0 , you should specified the encoding information (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteT(T)](9f48853f-fbeb-c37d-25a2-1c76d1dee407.htm) | 写入支持Hsl特性的数据，返回是否写入成功，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Write data that supports the Hsl attribute, and return whether the write was successful. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WriteAddress](37e7ec52-2c99-b84f-3062-f23d3541a67a.htm) | 写入设备的地址域信息，仅支持点对点通讯的情况，需要指定地址域信息，例如：14910000729011  Write the address domain information of the device, only support point-to-point communication, you need to specify the address domain information, for example: 14910000729011 |
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

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[CJT188 类](7845b81d-5321-3a2e-346d-ac4fb15e3f2a.htm)

[HslCommunication.Instrument.CJT 命名空间](32a7bbe4-7443-0a40-8e16-c1df11e7a933.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ActiveDeveice 方法 

[原文連結](http://api.hslcommunication.cn/html/6de9a868-a602-a896-330c-e56955f8eb96.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.CJT](../html/32a7bbe4-7443-0a40-8e16-c1df11e7a933.htm "HslCommunication.Instrument.CJT")

[CJT188 类](../html/7845b81d-5321-3a2e-346d-ac4fb15e3f2a.htm "CJT188 类")

[CJT188 方法](../html/04972e44-8c91-c7ca-0ab7-28b426216410.htm "CJT188 方法")

[ActiveDeveice 方法](../html/6de9a868-a602-a896-330c-e56955f8eb96.htm "ActiveDeveice 方法 ")

[GetNewNetMessage 方法](../html/d78644ed-1883-a2f5-91ad-741d56356d12.htm "GetNewNetMessage 方法 ")

[PackCommandWithHeader 方法](../html/84add1a8-4ab7-c216-b67e-0312b85d665e.htm "PackCommandWithHeader 方法 ")

[Read 方法](../html/1f6fea44-369c-71f2-4883-631528d6f3e5.htm "Read 方法 ")

[ReadAddress 方法](../html/81a31692-1162-4850-5404-ed02f888f324.htm "ReadAddress 方法 ")

[ReadDouble 方法](../html/b5fa76a1-104a-bdb6-c6d6-7472db8f735f.htm "ReadDouble 方法 ")

[ReadDoubleAsync 方法](../html/e53e8ff4-31bd-336c-4218-2c8fc79855c7.htm "ReadDoubleAsync 方法 ")

[ReadFloat 方法](../html/a5eeaaa8-4c0b-8d81-4420-3c7ce42a39f2.htm "ReadFloat 方法 ")

[ReadFloatAsync 方法](../html/7dd8d20e-35b5-8258-4550-df204a38e81f.htm "ReadFloatAsync 方法 ")

[ReadFromCoreServer 方法](../html/42325ee8-6793-9ad0-6b53-58f5487de0b1.htm "ReadFromCoreServer 方法 ")

[ReadString 方法](../html/e0ea9f3d-4c56-6078-2d10-e14579abb235.htm "ReadString 方法 ")

[ReadStringArray 方法](../html/c793adda-9b2c-6746-467e-b07a4f88e64a.htm "ReadStringArray 方法 ")

[ReadStringAsync 方法](../html/947017c8-6aa1-d267-3e9d-3b880b5021be.htm "ReadStringAsync 方法 ")

[ToString 方法](../html/8bff78ed-c0ad-15e4-e369-ad46dbff3241.htm "ToString 方法 ")

[Write 方法](../html/efa8452a-5af5-3ac6-b3fd-8f9c33845943.htm "Write 方法 ")

[WriteAddress 方法](../html/37e7ec52-2c99-b84f-3062-f23d3541a67a.htm "WriteAddress 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CJT188ActiveDeveice 方法 |

激活设备的命令，只发送数据到设备，不等待设备数据返回  
The command to activate the device, only send data to the device, do not wait for the device data to return

**命名空间：**
 [HslCommunication.Instrument.CJT](32a7bbe4-7443-0a40-8e16-c1df11e7a933.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public OperateResult ActiveDeveice()
```

```
Public Function ActiveDeveice As OperateResult
```

```
public:
virtual OperateResult^ ActiveDeveice() sealed
```

```
abstract ActiveDeveice : unit -> OperateResult 
override ActiveDeveice : unit -> OperateResult
```

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否发送成功

#### 实现

[ICjt188ActiveDeveice](f00c3a89-859f-c8ba-5c0a-b58cd4162eaa.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[CJT188 类](7845b81d-5321-3a2e-346d-ac4fb15e3f2a.htm)

[HslCommunication.Instrument.CJT 命名空间](32a7bbe4-7443-0a40-8e16-c1df11e7a933.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GetNewNetMessage 方法 

[原文連結](http://api.hslcommunication.cn/html/d78644ed-1883-a2f5-91ad-741d56356d12.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.CJT](../html/32a7bbe4-7443-0a40-8e16-c1df11e7a933.htm "HslCommunication.Instrument.CJT")

[CJT188 类](../html/7845b81d-5321-3a2e-346d-ac4fb15e3f2a.htm "CJT188 类")

[CJT188 方法](../html/04972e44-8c91-c7ca-0ab7-28b426216410.htm "CJT188 方法")

[ActiveDeveice 方法](../html/6de9a868-a602-a896-330c-e56955f8eb96.htm "ActiveDeveice 方法 ")

[GetNewNetMessage 方法](../html/d78644ed-1883-a2f5-91ad-741d56356d12.htm "GetNewNetMessage 方法 ")

[PackCommandWithHeader 方法](../html/84add1a8-4ab7-c216-b67e-0312b85d665e.htm "PackCommandWithHeader 方法 ")

[Read 方法](../html/1f6fea44-369c-71f2-4883-631528d6f3e5.htm "Read 方法 ")

[ReadAddress 方法](../html/81a31692-1162-4850-5404-ed02f888f324.htm "ReadAddress 方法 ")

[ReadDouble 方法](../html/b5fa76a1-104a-bdb6-c6d6-7472db8f735f.htm "ReadDouble 方法 ")

[ReadDoubleAsync 方法](../html/e53e8ff4-31bd-336c-4218-2c8fc79855c7.htm "ReadDoubleAsync 方法 ")

[ReadFloat 方法](../html/a5eeaaa8-4c0b-8d81-4420-3c7ce42a39f2.htm "ReadFloat 方法 ")

[ReadFloatAsync 方法](../html/7dd8d20e-35b5-8258-4550-df204a38e81f.htm "ReadFloatAsync 方法 ")

[ReadFromCoreServer 方法](../html/42325ee8-6793-9ad0-6b53-58f5487de0b1.htm "ReadFromCoreServer 方法 ")

[ReadString 方法](../html/e0ea9f3d-4c56-6078-2d10-e14579abb235.htm "ReadString 方法 ")

[ReadStringArray 方法](../html/c793adda-9b2c-6746-467e-b07a4f88e64a.htm "ReadStringArray 方法 ")

[ReadStringAsync 方法](../html/947017c8-6aa1-d267-3e9d-3b880b5021be.htm "ReadStringAsync 方法 ")

[ToString 方法](../html/8bff78ed-c0ad-15e4-e369-ad46dbff3241.htm "ToString 方法 ")

[Write 方法](../html/efa8452a-5af5-3ac6-b3fd-8f9c33845943.htm "Write 方法 ")

[WriteAddress 方法](../html/37e7ec52-2c99-b84f-3062-f23d3541a67a.htm "WriteAddress 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CJT188GetNewNetMessage 方法 |

获取一个新的消息对象的方法，需要在继承类里面进行重写  
The method to get a new message object needs to be overridden in the inheritance class

**命名空间：**
 [HslCommunication.Instrument.CJT](32a7bbe4-7443-0a40-8e16-c1df11e7a933.htm)  
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

[CJT188 类](7845b81d-5321-3a2e-346d-ac4fb15e3f2a.htm)

[HslCommunication.Instrument.CJT 命名空间](32a7bbe4-7443-0a40-8e16-c1df11e7a933.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## PackCommandWithHeader 方法 

[原文連結](http://api.hslcommunication.cn/html/84add1a8-4ab7-c216-b67e-0312b85d665e.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.CJT](../html/32a7bbe4-7443-0a40-8e16-c1df11e7a933.htm "HslCommunication.Instrument.CJT")

[CJT188 类](../html/7845b81d-5321-3a2e-346d-ac4fb15e3f2a.htm "CJT188 类")

[CJT188 方法](../html/04972e44-8c91-c7ca-0ab7-28b426216410.htm "CJT188 方法")

[ActiveDeveice 方法](../html/6de9a868-a602-a896-330c-e56955f8eb96.htm "ActiveDeveice 方法 ")

[GetNewNetMessage 方法](../html/d78644ed-1883-a2f5-91ad-741d56356d12.htm "GetNewNetMessage 方法 ")

[PackCommandWithHeader 方法](../html/84add1a8-4ab7-c216-b67e-0312b85d665e.htm "PackCommandWithHeader 方法 ")

[Read 方法](../html/1f6fea44-369c-71f2-4883-631528d6f3e5.htm "Read 方法 ")

[ReadAddress 方法](../html/81a31692-1162-4850-5404-ed02f888f324.htm "ReadAddress 方法 ")

[ReadDouble 方法](../html/b5fa76a1-104a-bdb6-c6d6-7472db8f735f.htm "ReadDouble 方法 ")

[ReadDoubleAsync 方法](../html/e53e8ff4-31bd-336c-4218-2c8fc79855c7.htm "ReadDoubleAsync 方法 ")

[ReadFloat 方法](../html/a5eeaaa8-4c0b-8d81-4420-3c7ce42a39f2.htm "ReadFloat 方法 ")

[ReadFloatAsync 方法](../html/7dd8d20e-35b5-8258-4550-df204a38e81f.htm "ReadFloatAsync 方法 ")

[ReadFromCoreServer 方法](../html/42325ee8-6793-9ad0-6b53-58f5487de0b1.htm "ReadFromCoreServer 方法 ")

[ReadString 方法](../html/e0ea9f3d-4c56-6078-2d10-e14579abb235.htm "ReadString 方法 ")

[ReadStringArray 方法](../html/c793adda-9b2c-6746-467e-b07a4f88e64a.htm "ReadStringArray 方法 ")

[ReadStringAsync 方法](../html/947017c8-6aa1-d267-3e9d-3b880b5021be.htm "ReadStringAsync 方法 ")

[ToString 方法](../html/8bff78ed-c0ad-15e4-e369-ad46dbff3241.htm "ToString 方法 ")

[Write 方法](../html/efa8452a-5af5-3ac6-b3fd-8f9c33845943.htm "Write 方法 ")

[WriteAddress 方法](../html/37e7ec52-2c99-b84f-3062-f23d3541a67a.htm "WriteAddress 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CJT188PackCommandWithHeader 方法 |

对当前的命令进行打包处理，通常是携带命令头内容，标记当前的命令的长度信息，需要进行重写，否则默认不打包  
The current command is packaged, usually carrying the content of the command header, marking the length of the current command,
and it needs to be rewritten, otherwise it is not packaged by default

**命名空间：**
 [HslCommunication.Instrument.CJT](32a7bbe4-7443-0a40-8e16-c1df11e7a933.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public override byte[] PackCommandWithHeader(
	byte[] command
)
```

```
Public Overrides Function PackCommandWithHeader ( 
	command As Byte()
) As Byte()
```

```
public:
virtual array<unsigned char>^ PackCommandWithHeader(
	array<unsigned char>^ command
) override
```

```
abstract PackCommandWithHeader : 
        command : byte[] -> byte[] 
override PackCommandWithHeader : 
        command : byte[] -> byte[]
```

#### 参数

command
:   类型：SystemByte  
    发送的数据命令内容

#### 返回值

类型：Byte  
打包之后的数据结果信息

![](../icons/SectionExpanded.png)备注

对发送的命令打包之后，直接发送给真实的对方设备了，例如在AB-PLC里面，就重写了打包方法，将当前的会话ID参数传递给PLC设备  
After packaging the sent command, it is directly sent to the real counterpart device. For example, in AB-PLC,
the packaging method is rewritten and the current session ID parameter is passed to the PLC device.

![](../icons/SectionExpanded.png)参见

#### 引用

[CJT188 类](7845b81d-5321-3a2e-346d-ac4fb15e3f2a.htm)

[HslCommunication.Instrument.CJT 命名空间](32a7bbe4-7443-0a40-8e16-c1df11e7a933.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Read 方法 

[原文連結](http://api.hslcommunication.cn/html/1f6fea44-369c-71f2-4883-631528d6f3e5.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.CJT](../html/32a7bbe4-7443-0a40-8e16-c1df11e7a933.htm "HslCommunication.Instrument.CJT")

[CJT188 类](../html/7845b81d-5321-3a2e-346d-ac4fb15e3f2a.htm "CJT188 类")

[CJT188 方法](../html/04972e44-8c91-c7ca-0ab7-28b426216410.htm "CJT188 方法")

[Read 方法](../html/1f6fea44-369c-71f2-4883-631528d6f3e5.htm "Read 方法 ")

[Read 方法 (String, UInt16)](../html/7b7961dc-c496-8545-5d9a-eac041553faa.htm "Read 方法 (String, UInt16)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CJT188Read 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法代码示例 | [ReadT](95abe064-5594-5ae7-a849-98c53f88fee0.htm) | 读取支持Hsl特性的数据内容，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明，或是官网的详细文档  Read the data content of the Hsl attribute. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [Read(String, UInt16)](7b7961dc-c496-8545-5d9a-eac041553faa.htm) | 根据指定的数据标识来读取相关的原始数据信息，地址标识根据手册来，从高位到地位，例如 91-1F，分割符可以任意特殊字符或是没有分隔符。  Read the relevant original data information according to the specified data identifier. The address identifier is based on the manual, from high to position, such as 91-1F. The separator can be any special character or no separator. (重写 [DeviceCommunicationRead(String, UInt16)](9dd09452-c200-bc65-75ff-69aa34b366d4.htm).) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[CJT188 类](7845b81d-5321-3a2e-346d-ac4fb15e3f2a.htm)

[HslCommunication.Instrument.CJT 命名空间](32a7bbe4-7443-0a40-8e16-c1df11e7a933.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Read 方法 (String, UInt16)

[原文連結](http://api.hslcommunication.cn/html/7b7961dc-c496-8545-5d9a-eac041553faa.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.CJT](../html/32a7bbe4-7443-0a40-8e16-c1df11e7a933.htm "HslCommunication.Instrument.CJT")

[CJT188 类](../html/7845b81d-5321-3a2e-346d-ac4fb15e3f2a.htm "CJT188 类")

[CJT188 方法](../html/04972e44-8c91-c7ca-0ab7-28b426216410.htm "CJT188 方法")

[Read 方法](../html/1f6fea44-369c-71f2-4883-631528d6f3e5.htm "Read 方法 ")

[Read 方法 (String, UInt16)](../html/7b7961dc-c496-8545-5d9a-eac041553faa.htm "Read 方法 (String, UInt16)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CJT188Read 方法 (String, UInt16) |

根据指定的数据标识来读取相关的原始数据信息，地址标识根据手册来，从高位到地位，例如 91-1F，分割符可以任意特殊字符或是没有分隔符。  
Read the relevant original data information according to the specified data identifier. The address identifier is based on the manual,
from high to position, such as 91-1F. The separator can be any special character or no separator.

**命名空间：**
 [HslCommunication.Instrument.CJT](32a7bbe4-7443-0a40-8e16-c1df11e7a933.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public override OperateResult<byte[]> Read(
	string address,
	ushort length
)
```

```
Public Overrides Function Read ( 
	address As String,
	length As UShort
) As OperateResult(Of Byte())
```

```
public:
virtual OperateResult<array<unsigned char>^>^ Read(
	String^ address, 
	unsigned short length
) override
```

```
abstract Read : 
        address : string * 
        length : uint16 -> OperateResult<byte[]> 
override Read : 
        address : string * 
        length : uint16 -> OperateResult<byte[]>
```

#### 参数

address
:   类型：SystemString  
    数据标识，具体需要查找手册来对应

length
:   类型：SystemUInt16  
    数据长度信息

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
结果信息

#### 实现

[IReadWriteNetRead(String, UInt16)](b8d4187f-ab9a-79e9-c7ca-8a01ec0df233.htm)  
[IReadWriteNetRead(String, UInt16)](b8d4187f-ab9a-79e9-c7ca-8a01ec0df233.htm)

![](../icons/SectionExpanded.png)备注

地址可以携带地址域信息，例如 "s=2;90-1F" 或是 "s=100000;90-1F"，关于数据域信息，需要查找手册，例如:D1-20 表示： 上一月结算日累积流量

![](../icons/SectionExpanded.png)参见

#### 引用

[CJT188 类](7845b81d-5321-3a2e-346d-ac4fb15e3f2a.htm)

[Read 重载](1f6fea44-369c-71f2-4883-631528d6f3e5.htm)

[HslCommunication.Instrument.CJT 命名空间](32a7bbe4-7443-0a40-8e16-c1df11e7a933.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadAddress 方法 

[原文連結](http://api.hslcommunication.cn/html/81a31692-1162-4850-5404-ed02f888f324.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.CJT](../html/32a7bbe4-7443-0a40-8e16-c1df11e7a933.htm "HslCommunication.Instrument.CJT")

[CJT188 类](../html/7845b81d-5321-3a2e-346d-ac4fb15e3f2a.htm "CJT188 类")

[CJT188 方法](../html/04972e44-8c91-c7ca-0ab7-28b426216410.htm "CJT188 方法")

[ActiveDeveice 方法](../html/6de9a868-a602-a896-330c-e56955f8eb96.htm "ActiveDeveice 方法 ")

[GetNewNetMessage 方法](../html/d78644ed-1883-a2f5-91ad-741d56356d12.htm "GetNewNetMessage 方法 ")

[PackCommandWithHeader 方法](../html/84add1a8-4ab7-c216-b67e-0312b85d665e.htm "PackCommandWithHeader 方法 ")

[Read 方法](../html/1f6fea44-369c-71f2-4883-631528d6f3e5.htm "Read 方法 ")

[ReadAddress 方法](../html/81a31692-1162-4850-5404-ed02f888f324.htm "ReadAddress 方法 ")

[ReadDouble 方法](../html/b5fa76a1-104a-bdb6-c6d6-7472db8f735f.htm "ReadDouble 方法 ")

[ReadDoubleAsync 方法](../html/e53e8ff4-31bd-336c-4218-2c8fc79855c7.htm "ReadDoubleAsync 方法 ")

[ReadFloat 方法](../html/a5eeaaa8-4c0b-8d81-4420-3c7ce42a39f2.htm "ReadFloat 方法 ")

[ReadFloatAsync 方法](../html/7dd8d20e-35b5-8258-4550-df204a38e81f.htm "ReadFloatAsync 方法 ")

[ReadFromCoreServer 方法](../html/42325ee8-6793-9ad0-6b53-58f5487de0b1.htm "ReadFromCoreServer 方法 ")

[ReadString 方法](../html/e0ea9f3d-4c56-6078-2d10-e14579abb235.htm "ReadString 方法 ")

[ReadStringArray 方法](../html/c793adda-9b2c-6746-467e-b07a4f88e64a.htm "ReadStringArray 方法 ")

[ReadStringAsync 方法](../html/947017c8-6aa1-d267-3e9d-3b880b5021be.htm "ReadStringAsync 方法 ")

[ToString 方法](../html/8bff78ed-c0ad-15e4-e369-ad46dbff3241.htm "ToString 方法 ")

[Write 方法](../html/efa8452a-5af5-3ac6-b3fd-8f9c33845943.htm "Write 方法 ")

[WriteAddress 方法](../html/37e7ec52-2c99-b84f-3062-f23d3541a67a.htm "WriteAddress 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CJT188ReadAddress 方法 |

读取设备的通信地址，仅支持点对点通讯的情况，返回地址域数据，例如：14910000729012  
Read the communication address of the device, only support point-to-point communication, and return the address field data, for example: 149100007290

**命名空间：**
 [HslCommunication.Instrument.CJT](32a7bbe4-7443-0a40-8e16-c1df11e7a933.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public OperateResult<string> ReadAddress()
```

```
Public Function ReadAddress As OperateResult(Of String)
```

```
public:
virtual OperateResult<String^>^ ReadAddress() sealed
```

```
abstract ReadAddress : unit -> OperateResult<string> 
override ReadAddress : unit -> OperateResult<string>
```

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)String  
设备的通信地址

#### 实现

[ICjt188ReadAddress](a669de5c-80f0-74df-235c-2eccb0ff9c38.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[CJT188 类](7845b81d-5321-3a2e-346d-ac4fb15e3f2a.htm)

[HslCommunication.Instrument.CJT 命名空间](32a7bbe4-7443-0a40-8e16-c1df11e7a933.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadDouble 方法 

[原文連結](http://api.hslcommunication.cn/html/b5fa76a1-104a-bdb6-c6d6-7472db8f735f.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.CJT](../html/32a7bbe4-7443-0a40-8e16-c1df11e7a933.htm "HslCommunication.Instrument.CJT")

[CJT188 类](../html/7845b81d-5321-3a2e-346d-ac4fb15e3f2a.htm "CJT188 类")

[CJT188 方法](../html/04972e44-8c91-c7ca-0ab7-28b426216410.htm "CJT188 方法")

[ReadDouble 方法](../html/b5fa76a1-104a-bdb6-c6d6-7472db8f735f.htm "ReadDouble 方法 ")

[ReadDouble 方法 (String, UInt16)](../html/4eceecff-16be-93ac-cfc9-da43972e3f19.htm "ReadDouble 方法 (String, UInt16)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CJT188ReadDouble 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法代码示例 | [ReadDouble(String)](e7e31c43-71b4-e692-6bdc-ba9533e1d83d.htm) | 读取双浮点的数据  Read double floating point data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadDouble(String, UInt16)](4eceecff-16be-93ac-cfc9-da43972e3f19.htm) | 读取双浮点数据的数组  Read double floating point data array (重写 [DeviceCommunicationReadDouble(String, UInt16)](4eee0f05-f861-4b48-3f06-99737fb24aee.htm).) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[CJT188 类](7845b81d-5321-3a2e-346d-ac4fb15e3f2a.htm)

[HslCommunication.Instrument.CJT 命名空间](32a7bbe4-7443-0a40-8e16-c1df11e7a933.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadDouble 方法 (String, UInt16)

[原文連結](http://api.hslcommunication.cn/html/4eceecff-16be-93ac-cfc9-da43972e3f19.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.CJT](../html/32a7bbe4-7443-0a40-8e16-c1df11e7a933.htm "HslCommunication.Instrument.CJT")

[CJT188 类](../html/7845b81d-5321-3a2e-346d-ac4fb15e3f2a.htm "CJT188 类")

[CJT188 方法](../html/04972e44-8c91-c7ca-0ab7-28b426216410.htm "CJT188 方法")

[ReadDouble 方法](../html/b5fa76a1-104a-bdb6-c6d6-7472db8f735f.htm "ReadDouble 方法 ")

[ReadDouble 方法 (String, UInt16)](../html/4eceecff-16be-93ac-cfc9-da43972e3f19.htm "ReadDouble 方法 (String, UInt16)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CJT188ReadDouble 方法 (String, UInt16) |

读取双浮点数据的数组  
Read double floating point data array

**命名空间：**
 [HslCommunication.Instrument.CJT](32a7bbe4-7443-0a40-8e16-c1df11e7a933.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public override OperateResult<double[]> ReadDouble(
	string address,
	ushort length
)
```

```
Public Overrides Function ReadDouble ( 
	address As String,
	length As UShort
) As OperateResult(Of Double())
```

```
public:
virtual OperateResult<array<double>^>^ ReadDouble(
	String^ address, 
	unsigned short length
) override
```

```
abstract ReadDouble : 
        address : string * 
        length : uint16 -> OperateResult<float[]> 
override ReadDouble : 
        address : string * 
        length : uint16 -> OperateResult<float[]>
```

#### 参数

address
:   类型：SystemString  
    起始地址

length
:   类型：SystemUInt16  
    数组长度

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Double  
带有成功标识的double数组

#### 实现

[IReadWriteNetReadDouble(String, UInt16)](a8468555-105b-337a-196f-b0a667856edd.htm)  
[IReadWriteNetReadDouble(String, UInt16)](a8468555-105b-337a-196f-b0a667856edd.htm)

![](../icons/SectionExpanded.png)示例

以下为三菱的连接对象示例，其他的设备读写情况参照下面的代码：

Double类型示例

[复制](# "复制")

```
MelsecMcNet melsec_net = new MelsecMcNet( "192.168.0.100", 6000 );

// 以下是简单的读取，没有仔细校验的方式
double[] d100_139 = melsec_net.ReadDouble( "D100", 10 ).Content;

// 如果需要判断是否读取成功

OperateResult<double[]> R_d100_139 = melsec_net.ReadDouble( "D100", 10 );
if (R_d100_139.IsSuccess)
{
    double value_d100 = R_d100_139.Content[0];
    double value_d104 = R_d100_139.Content[1];
    double value_d108 = R_d100_139.Content[2];
    double value_d112 = R_d100_139.Content[3];
    double value_d116 = R_d100_139.Content[4];
    double value_d120 = R_d100_139.Content[5];
    double value_d124 = R_d100_139.Content[6];
    double value_d128 = R_d100_139.Content[7];
    double value_d132 = R_d100_139.Content[8];
    double value_d136 = R_d100_139.Content[9];
}
else
{
    // failed
}
```

![](../icons/SectionExpanded.png)参见

#### 引用

[CJT188 类](7845b81d-5321-3a2e-346d-ac4fb15e3f2a.htm)

[ReadDouble 重载](b5fa76a1-104a-bdb6-c6d6-7472db8f735f.htm)

[HslCommunication.Instrument.CJT 命名空间](32a7bbe4-7443-0a40-8e16-c1df11e7a933.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadDoubleAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/e53e8ff4-31bd-336c-4218-2c8fc79855c7.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.CJT](../html/32a7bbe4-7443-0a40-8e16-c1df11e7a933.htm "HslCommunication.Instrument.CJT")

[CJT188 类](../html/7845b81d-5321-3a2e-346d-ac4fb15e3f2a.htm "CJT188 类")

[CJT188 方法](../html/04972e44-8c91-c7ca-0ab7-28b426216410.htm "CJT188 方法")

[ReadDoubleAsync 方法](../html/e53e8ff4-31bd-336c-4218-2c8fc79855c7.htm "ReadDoubleAsync 方法 ")

[ReadDoubleAsync 方法 (String, UInt16)](../html/51c76fdf-a6e4-f9e7-975c-ef8bf7fa1e7a.htm "ReadDoubleAsync 方法 (String, UInt16)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CJT188ReadDoubleAsync 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法代码示例 | [ReadDoubleAsync(String)](2cf895fa-1033-1def-f048-0e8caefaa71d.htm) | 异步读取双浮点的数据  Asynchronously read double floating point data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadDoubleAsync(String, UInt16)](51c76fdf-a6e4-f9e7-975c-ef8bf7fa1e7a.htm) | 读取双浮点数据的数组  Read double floating point data array (重写 [DeviceCommunicationReadDoubleAsync(String, UInt16)](2256d406-d4d5-da54-210b-ccef58a5e914.htm).) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[CJT188 类](7845b81d-5321-3a2e-346d-ac4fb15e3f2a.htm)

[HslCommunication.Instrument.CJT 命名空间](32a7bbe4-7443-0a40-8e16-c1df11e7a933.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadDoubleAsync 方法 (String, UInt16)

[原文連結](http://api.hslcommunication.cn/html/51c76fdf-a6e4-f9e7-975c-ef8bf7fa1e7a.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.CJT](../html/32a7bbe4-7443-0a40-8e16-c1df11e7a933.htm "HslCommunication.Instrument.CJT")

[CJT188 类](../html/7845b81d-5321-3a2e-346d-ac4fb15e3f2a.htm "CJT188 类")

[CJT188 方法](../html/04972e44-8c91-c7ca-0ab7-28b426216410.htm "CJT188 方法")

[ReadDoubleAsync 方法](../html/e53e8ff4-31bd-336c-4218-2c8fc79855c7.htm "ReadDoubleAsync 方法 ")

[ReadDoubleAsync 方法 (String, UInt16)](../html/51c76fdf-a6e4-f9e7-975c-ef8bf7fa1e7a.htm "ReadDoubleAsync 方法 (String, UInt16)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CJT188ReadDoubleAsync 方法 (String, UInt16) |

读取双浮点数据的数组  
Read double floating point data array

**命名空间：**
 [HslCommunication.Instrument.CJT](32a7bbe4-7443-0a40-8e16-c1df11e7a933.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public override Task<OperateResult<double[]>> ReadDoubleAsync(
	string address,
	ushort length
)
```

```
Public Overrides Function ReadDoubleAsync ( 
	address As String,
	length As UShort
) As Task(Of OperateResult(Of Double()))
```

```
public:
virtual Task<OperateResult<array<double>^>^>^ ReadDoubleAsync(
	String^ address, 
	unsigned short length
) override
```

```
abstract ReadDoubleAsync : 
        address : string * 
        length : uint16 -> Task<OperateResult<float[]>> 
override ReadDoubleAsync : 
        address : string * 
        length : uint16 -> Task<OperateResult<float[]>>
```

#### 参数

address
:   类型：SystemString  
    起始地址

length
:   类型：SystemUInt16  
    数组长度

#### 返回值

类型：Task[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Double  
带有成功标识的double数组

#### 实现

[IReadWriteNetReadDoubleAsync(String, UInt16)](3a436322-c7f5-bb1d-e97b-dc6c0d6990d2.htm)  
[IReadWriteNetReadDoubleAsync(String, UInt16)](3a436322-c7f5-bb1d-e97b-dc6c0d6990d2.htm)

![](../icons/SectionExpanded.png)示例

以下为三菱的连接对象示例，其他的设备读写情况参照下面的代码：

Double类型示例

[复制](# "复制")

```
MelsecMcNet melsec_net = new MelsecMcNet( "192.168.0.100", 6000 );

// 以下是简单的读取，没有仔细校验的方式
double[] d100_139 = melsec_net.ReadDouble( "D100", 10 ).Content;

// 如果需要判断是否读取成功

OperateResult<double[]> R_d100_139 = melsec_net.ReadDouble( "D100", 10 );
if (R_d100_139.IsSuccess)
{
    double value_d100 = R_d100_139.Content[0];
    double value_d104 = R_d100_139.Content[1];
    double value_d108 = R_d100_139.Content[2];
    double value_d112 = R_d100_139.Content[3];
    double value_d116 = R_d100_139.Content[4];
    double value_d120 = R_d100_139.Content[5];
    double value_d124 = R_d100_139.Content[6];
    double value_d128 = R_d100_139.Content[7];
    double value_d132 = R_d100_139.Content[8];
    double value_d136 = R_d100_139.Content[9];
}
else
{
    // failed
}
```

![](../icons/SectionExpanded.png)参见

#### 引用

[CJT188 类](7845b81d-5321-3a2e-346d-ac4fb15e3f2a.htm)

[ReadDoubleAsync 重载](e53e8ff4-31bd-336c-4218-2c8fc79855c7.htm)

[HslCommunication.Instrument.CJT 命名空间](32a7bbe4-7443-0a40-8e16-c1df11e7a933.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadFloat 方法 

[原文連結](http://api.hslcommunication.cn/html/a5eeaaa8-4c0b-8d81-4420-3c7ce42a39f2.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.CJT](../html/32a7bbe4-7443-0a40-8e16-c1df11e7a933.htm "HslCommunication.Instrument.CJT")

[CJT188 类](../html/7845b81d-5321-3a2e-346d-ac4fb15e3f2a.htm "CJT188 类")

[CJT188 方法](../html/04972e44-8c91-c7ca-0ab7-28b426216410.htm "CJT188 方法")

[ReadFloat 方法](../html/a5eeaaa8-4c0b-8d81-4420-3c7ce42a39f2.htm "ReadFloat 方法 ")

[ReadFloat 方法 (String, UInt16)](../html/91bc4f0a-56d6-51d8-ddda-d776174aded3.htm "ReadFloat 方法 (String, UInt16)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CJT188ReadFloat 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法代码示例 | [ReadFloat(String)](bc2801ae-6600-e845-7c6c-522c9156eb0e.htm) | 读取单浮点数据  Read single floating point data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadFloat(String, UInt16)](91bc4f0a-56d6-51d8-ddda-d776174aded3.htm) | 读取单浮点精度的数组  Read single floating point array (重写 [DeviceCommunicationReadFloat(String, UInt16)](6a9fb098-b13c-313b-3b27-20d3c787945d.htm).) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[CJT188 类](7845b81d-5321-3a2e-346d-ac4fb15e3f2a.htm)

[HslCommunication.Instrument.CJT 命名空间](32a7bbe4-7443-0a40-8e16-c1df11e7a933.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadFloat 方法 (String, UInt16)

[原文連結](http://api.hslcommunication.cn/html/91bc4f0a-56d6-51d8-ddda-d776174aded3.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.CJT](../html/32a7bbe4-7443-0a40-8e16-c1df11e7a933.htm "HslCommunication.Instrument.CJT")

[CJT188 类](../html/7845b81d-5321-3a2e-346d-ac4fb15e3f2a.htm "CJT188 类")

[CJT188 方法](../html/04972e44-8c91-c7ca-0ab7-28b426216410.htm "CJT188 方法")

[ReadFloat 方法](../html/a5eeaaa8-4c0b-8d81-4420-3c7ce42a39f2.htm "ReadFloat 方法 ")

[ReadFloat 方法 (String, UInt16)](../html/91bc4f0a-56d6-51d8-ddda-d776174aded3.htm "ReadFloat 方法 (String, UInt16)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CJT188ReadFloat 方法 (String, UInt16) |

读取单浮点精度的数组  
Read single floating point array

**命名空间：**
 [HslCommunication.Instrument.CJT](32a7bbe4-7443-0a40-8e16-c1df11e7a933.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public override OperateResult<float[]> ReadFloat(
	string address,
	ushort length
)
```

```
Public Overrides Function ReadFloat ( 
	address As String,
	length As UShort
) As OperateResult(Of Single())
```

```
public:
virtual OperateResult<array<float>^>^ ReadFloat(
	String^ address, 
	unsigned short length
) override
```

```
abstract ReadFloat : 
        address : string * 
        length : uint16 -> OperateResult<float32[]> 
override ReadFloat : 
        address : string * 
        length : uint16 -> OperateResult<float32[]>
```

#### 参数

address
:   类型：SystemString  
    起始地址

length
:   类型：SystemUInt16  
    数组长度

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Single  
带有成功标识的float数组

#### 实现

[IReadWriteNetReadFloat(String, UInt16)](a8ae302d-6269-8baa-cacc-f35df9fcd131.htm)  
[IReadWriteNetReadFloat(String, UInt16)](a8ae302d-6269-8baa-cacc-f35df9fcd131.htm)

![](../icons/SectionExpanded.png)示例

以下为三菱的连接对象示例，其他的设备读写情况参照下面的代码：

Float类型示例

[复制](# "复制")

```
MelsecMcNet melsec_net = new MelsecMcNet( "192.168.0.100", 6000 );

// 以下是简单的读取，没有仔细校验的方式
float[] d100_119 = melsec_net.ReadFloat( "D100", 10 ).Content;

// 如果需要判断是否读取成功

OperateResult<float[]> R_d100_119 = melsec_net.ReadFloat( "D100", 10 );
if (R_d100_119.IsSuccess)
{
    float value_d100 = R_d100_119.Content[0];
    float value_d102 = R_d100_119.Content[1];
    float value_d104 = R_d100_119.Content[2];
    float value_d106 = R_d100_119.Content[3];
    float value_d108 = R_d100_119.Content[4];
    float value_d110 = R_d100_119.Content[5];
    float value_d112 = R_d100_119.Content[6];
    float value_d114 = R_d100_119.Content[7];
    float value_d116 = R_d100_119.Content[8];
    float value_d118 = R_d100_119.Content[9];
}
else
{
    // failed
}
```

![](../icons/SectionExpanded.png)参见

#### 引用

[CJT188 类](7845b81d-5321-3a2e-346d-ac4fb15e3f2a.htm)

[ReadFloat 重载](a5eeaaa8-4c0b-8d81-4420-3c7ce42a39f2.htm)

[HslCommunication.Instrument.CJT 命名空间](32a7bbe4-7443-0a40-8e16-c1df11e7a933.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadFloatAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/7dd8d20e-35b5-8258-4550-df204a38e81f.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.CJT](../html/32a7bbe4-7443-0a40-8e16-c1df11e7a933.htm "HslCommunication.Instrument.CJT")

[CJT188 类](../html/7845b81d-5321-3a2e-346d-ac4fb15e3f2a.htm "CJT188 类")

[CJT188 方法](../html/04972e44-8c91-c7ca-0ab7-28b426216410.htm "CJT188 方法")

[ReadFloatAsync 方法](../html/7dd8d20e-35b5-8258-4550-df204a38e81f.htm "ReadFloatAsync 方法 ")

[ReadFloatAsync 方法 (String, UInt16)](../html/86291e27-0fb1-9cff-9d92-fb62f3194724.htm "ReadFloatAsync 方法 (String, UInt16)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CJT188ReadFloatAsync 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法代码示例 | [ReadFloatAsync(String)](89f743f8-4b3a-7085-1e49-a1d7adb10052.htm) | 异步读取单浮点数据  Asynchronously read single floating point data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadFloatAsync(String, UInt16)](86291e27-0fb1-9cff-9d92-fb62f3194724.htm) | 读取单浮点精度的数组  Read single floating point array (重写 [DeviceCommunicationReadFloatAsync(String, UInt16)](cb6d5c5d-9b5c-2467-3e6b-5efa9d20cd72.htm).) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[CJT188 类](7845b81d-5321-3a2e-346d-ac4fb15e3f2a.htm)

[HslCommunication.Instrument.CJT 命名空间](32a7bbe4-7443-0a40-8e16-c1df11e7a933.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadFloatAsync 方法 (String, UInt16)

[原文連結](http://api.hslcommunication.cn/html/86291e27-0fb1-9cff-9d92-fb62f3194724.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.CJT](../html/32a7bbe4-7443-0a40-8e16-c1df11e7a933.htm "HslCommunication.Instrument.CJT")

[CJT188 类](../html/7845b81d-5321-3a2e-346d-ac4fb15e3f2a.htm "CJT188 类")

[CJT188 方法](../html/04972e44-8c91-c7ca-0ab7-28b426216410.htm "CJT188 方法")

[ReadFloatAsync 方法](../html/7dd8d20e-35b5-8258-4550-df204a38e81f.htm "ReadFloatAsync 方法 ")

[ReadFloatAsync 方法 (String, UInt16)](../html/86291e27-0fb1-9cff-9d92-fb62f3194724.htm "ReadFloatAsync 方法 (String, UInt16)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CJT188ReadFloatAsync 方法 (String, UInt16) |

读取单浮点精度的数组  
Read single floating point array

**命名空间：**
 [HslCommunication.Instrument.CJT](32a7bbe4-7443-0a40-8e16-c1df11e7a933.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public override Task<OperateResult<float[]>> ReadFloatAsync(
	string address,
	ushort length
)
```

```
Public Overrides Function ReadFloatAsync ( 
	address As String,
	length As UShort
) As Task(Of OperateResult(Of Single()))
```

```
public:
virtual Task<OperateResult<array<float>^>^>^ ReadFloatAsync(
	String^ address, 
	unsigned short length
) override
```

```
abstract ReadFloatAsync : 
        address : string * 
        length : uint16 -> Task<OperateResult<float32[]>> 
override ReadFloatAsync : 
        address : string * 
        length : uint16 -> Task<OperateResult<float32[]>>
```

#### 参数

address
:   类型：SystemString  
    起始地址

length
:   类型：SystemUInt16  
    数组长度

#### 返回值

类型：Task[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Single  
带有成功标识的float数组

#### 实现

[IReadWriteNetReadFloatAsync(String, UInt16)](26a13b36-422b-0520-3cb0-66a235e3d7a8.htm)  
[IReadWriteNetReadFloatAsync(String, UInt16)](26a13b36-422b-0520-3cb0-66a235e3d7a8.htm)

![](../icons/SectionExpanded.png)示例

以下为三菱的连接对象示例，其他的设备读写情况参照下面的代码：

Float类型示例

[复制](# "复制")

```
MelsecMcNet melsec_net = new MelsecMcNet( "192.168.0.100", 6000 );

// 以下是简单的读取，没有仔细校验的方式
float[] d100_119 = melsec_net.ReadFloat( "D100", 10 ).Content;

// 如果需要判断是否读取成功

OperateResult<float[]> R_d100_119 = melsec_net.ReadFloat( "D100", 10 );
if (R_d100_119.IsSuccess)
{
    float value_d100 = R_d100_119.Content[0];
    float value_d102 = R_d100_119.Content[1];
    float value_d104 = R_d100_119.Content[2];
    float value_d106 = R_d100_119.Content[3];
    float value_d108 = R_d100_119.Content[4];
    float value_d110 = R_d100_119.Content[5];
    float value_d112 = R_d100_119.Content[6];
    float value_d114 = R_d100_119.Content[7];
    float value_d116 = R_d100_119.Content[8];
    float value_d118 = R_d100_119.Content[9];
}
else
{
    // failed
}
```

![](../icons/SectionExpanded.png)参见

#### 引用

[CJT188 类](7845b81d-5321-3a2e-346d-ac4fb15e3f2a.htm)

[ReadFloatAsync 重载](7dd8d20e-35b5-8258-4550-df204a38e81f.htm)

[HslCommunication.Instrument.CJT 命名空间](32a7bbe4-7443-0a40-8e16-c1df11e7a933.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadFromCoreServer 方法 

[原文連結](http://api.hslcommunication.cn/html/42325ee8-6793-9ad0-6b53-58f5487de0b1.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.CJT](../html/32a7bbe4-7443-0a40-8e16-c1df11e7a933.htm "HslCommunication.Instrument.CJT")

[CJT188 类](../html/7845b81d-5321-3a2e-346d-ac4fb15e3f2a.htm "CJT188 类")

[CJT188 方法](../html/04972e44-8c91-c7ca-0ab7-28b426216410.htm "CJT188 方法")

[ReadFromCoreServer 方法](../html/42325ee8-6793-9ad0-6b53-58f5487de0b1.htm "ReadFromCoreServer 方法 ")

[ReadFromCoreServer 方法 (Byte[])](../html/010eab80-ec07-2570-5860-38ad4ef731c8.htm "ReadFromCoreServer 方法 (Byte[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CJT188ReadFromCoreServer 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [ReadFromCoreServer(IEnumerableByte)](044cfdd9-ec76-b89f-2f2b-5f439025006c.htm) | 将多个数据报文按顺序发到设备，并从设备接收返回的数据内容，然后拼接成一个Byte[]信息，需要重写[UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm)方法才能返回正确的结果。  Send multiple data packets to the device in sequence, and receive the returned data content from the device, and then splicing them into a Byte[] message, you need to rewrite [UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm) method to return the correct result. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServer(Byte)](010eab80-ec07-2570-5860-38ad4ef731c8.htm) | 将二进制的数据发送到管道中去，然后从管道里接收二进制的数据回来，并返回是否成功的结果对象。  Send binary data to the pipeline, and then receive binary data back from the pipeline, and return whether the success of the result object (重写 [BinaryCommunicationReadFromCoreServer(Byte)](8146dc0c-1edc-3ca7-428a-a3bdbd35e6c5.htm).) |
| 公共方法代码示例 | [ReadFromCoreServer(Byte, Boolean, Boolean)](b32ce8e8-a959-fcc2-6736-14b39fdb9a40.htm) | 将二进制的数据发送到管道中去，然后从管道里接收二进制的数据回来，并返回是否成功的结果对象。  Send binary data to the pipeline, and then receive binary data back from the pipeline, and return whether the success of the result object (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServer(CommunicationPipe, Byte, Boolean, Boolean)](ef36d56b-df09-f727-e8d3-316d37a13ee9.htm) | 使用指定的管道来进行数据通信，发送原始数据到管道，然后从管道接收相关的数据返回，本方法无锁 (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[CJT188 类](7845b81d-5321-3a2e-346d-ac4fb15e3f2a.htm)

[HslCommunication.Instrument.CJT 命名空间](32a7bbe4-7443-0a40-8e16-c1df11e7a933.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadFromCoreServer 方法 (Byte[])

[原文連結](http://api.hslcommunication.cn/html/010eab80-ec07-2570-5860-38ad4ef731c8.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.CJT](../html/32a7bbe4-7443-0a40-8e16-c1df11e7a933.htm "HslCommunication.Instrument.CJT")

[CJT188 类](../html/7845b81d-5321-3a2e-346d-ac4fb15e3f2a.htm "CJT188 类")

[CJT188 方法](../html/04972e44-8c91-c7ca-0ab7-28b426216410.htm "CJT188 方法")

[ReadFromCoreServer 方法](../html/42325ee8-6793-9ad0-6b53-58f5487de0b1.htm "ReadFromCoreServer 方法 ")

[ReadFromCoreServer 方法 (Byte[])](../html/010eab80-ec07-2570-5860-38ad4ef731c8.htm "ReadFromCoreServer 方法 (Byte[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CJT188ReadFromCoreServer 方法 (Byte) |

将二进制的数据发送到管道中去，然后从管道里接收二进制的数据回来，并返回是否成功的结果对象。  
Send binary data to the pipeline, and then receive binary data back from the pipeline, and return whether the success of the result object

**命名空间：**
 [HslCommunication.Instrument.CJT](32a7bbe4-7443-0a40-8e16-c1df11e7a933.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public override OperateResult<byte[]> ReadFromCoreServer(
	byte[] send
)
```

```
Public Overrides Function ReadFromCoreServer ( 
	send As Byte()
) As OperateResult(Of Byte())
```

```
public:
virtual OperateResult<array<unsigned char>^>^ ReadFromCoreServer(
	array<unsigned char>^ send
) override
```

```
abstract ReadFromCoreServer : 
        send : byte[] -> OperateResult<byte[]> 
override ReadFromCoreServer : 
        send : byte[] -> OperateResult<byte[]>
```

#### 参数

send
:   类型：SystemByte  
    发送的完整的报文信息

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
接收的完整的报文信息

#### 实现

[IReadWriteDeviceReadFromCoreServer(Byte)](413a68ef-2f4b-638f-8bad-973673920de4.htm)

![](../icons/SectionExpanded.png)备注

本方法用于实现本组件还未实现的一些报文功能，例如有些modbus服务器会有一些特殊的功能码支持，需要收发特殊的报文，详细请看示例

![](../icons/SectionExpanded.png)示例

此处举例有个modbus服务器，有个特殊的功能码0x09，后面携带子数据0x01即可，发送字节为 0x00 0x00 0x00 0x00 0x00 0x03 0x01 0x09 0x01

ReadFromCoreServer示例

[复制](# "复制")

```
ModbusTcpNet modbus = new ModbusTcpNet( "192.168.0.100" );

// 此处举例实现特殊的modbus功能码
OperateResult<byte[]> read = modbus.ReadFromCoreServer( SoftBasic.HexStringToBytes( "0x00 0x00 0x00 0x00 0x00 0x03 0x01 0x09 0x01" ) );
if (read.IsSuccess)
{
    // 成功，开始解析从服务器返回的数据，是一条完整的报文信息
    Console.WriteLine( SoftBasic.ByteToHexString( read.Content, ' ' ) );
}
else
{
    // 失败
}
```

![](../icons/SectionExpanded.png)参见

#### 引用

[CJT188 类](7845b81d-5321-3a2e-346d-ac4fb15e3f2a.htm)

[ReadFromCoreServer 重载](42325ee8-6793-9ad0-6b53-58f5487de0b1.htm)

[HslCommunication.Instrument.CJT 命名空间](32a7bbe4-7443-0a40-8e16-c1df11e7a933.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadString 方法 

[原文連結](http://api.hslcommunication.cn/html/e0ea9f3d-4c56-6078-2d10-e14579abb235.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.CJT](../html/32a7bbe4-7443-0a40-8e16-c1df11e7a933.htm "HslCommunication.Instrument.CJT")

[CJT188 类](../html/7845b81d-5321-3a2e-346d-ac4fb15e3f2a.htm "CJT188 类")

[CJT188 方法](../html/04972e44-8c91-c7ca-0ab7-28b426216410.htm "CJT188 方法")

[ReadString 方法](../html/e0ea9f3d-4c56-6078-2d10-e14579abb235.htm "ReadString 方法 ")

[ReadString 方法 (String, UInt16, Encoding)](../html/d4dc35c1-5ab3-44f1-df6d-972add664fe3.htm "ReadString 方法 (String, UInt16, Encoding)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CJT188ReadString 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法代码示例 | [ReadString(String, UInt16)](07812d76-415a-6400-c8a2-1c21da51e7a3.htm) | 读取字符串数据，默认为最常见的ASCII编码  Read string data, default is the most common ASCII encoding (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadString(String, UInt16, Encoding)](d4dc35c1-5ab3-44f1-df6d-972add664fe3.htm) | 使用指定的编码，读取字符串数据  Reads string data using the specified encoding (重写 [DeviceCommunicationReadString(String, UInt16, Encoding)](b34c71a6-58f3-6bf4-0b93-471be4f2ed31.htm).) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[CJT188 类](7845b81d-5321-3a2e-346d-ac4fb15e3f2a.htm)

[HslCommunication.Instrument.CJT 命名空间](32a7bbe4-7443-0a40-8e16-c1df11e7a933.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadString 方法 (String, UInt16, Encoding)

[原文連結](http://api.hslcommunication.cn/html/d4dc35c1-5ab3-44f1-df6d-972add664fe3.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.CJT](../html/32a7bbe4-7443-0a40-8e16-c1df11e7a933.htm "HslCommunication.Instrument.CJT")

[CJT188 类](../html/7845b81d-5321-3a2e-346d-ac4fb15e3f2a.htm "CJT188 类")

[CJT188 方法](../html/04972e44-8c91-c7ca-0ab7-28b426216410.htm "CJT188 方法")

[ReadString 方法](../html/e0ea9f3d-4c56-6078-2d10-e14579abb235.htm "ReadString 方法 ")

[ReadString 方法 (String, UInt16, Encoding)](../html/d4dc35c1-5ab3-44f1-df6d-972add664fe3.htm "ReadString 方法 (String, UInt16, Encoding)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CJT188ReadString 方法 (String, UInt16, Encoding) |

使用指定的编码，读取字符串数据  
Reads string data using the specified encoding

**命名空间：**
 [HslCommunication.Instrument.CJT](32a7bbe4-7443-0a40-8e16-c1df11e7a933.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public override OperateResult<string> ReadString(
	string address,
	ushort length,
	Encoding encoding
)
```

```
Public Overrides Function ReadString ( 
	address As String,
	length As UShort,
	encoding As Encoding
) As OperateResult(Of String)
```

```
public:
virtual OperateResult<String^>^ ReadString(
	String^ address, 
	unsigned short length, 
	Encoding^ encoding
) override
```

```
abstract ReadString : 
        address : string * 
        length : uint16 * 
        encoding : Encoding -> OperateResult<string> 
override ReadString : 
        address : string * 
        length : uint16 * 
        encoding : Encoding -> OperateResult<string>
```

#### 参数

address
:   类型：SystemString  
    起始地址

length
:   类型：SystemUInt16  
    数据长度

encoding
:   类型：System.TextEncoding  
    指定的自定义的编码

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)String  
带有成功标识的string数据

#### 实现

[IReadWriteNetReadString(String, UInt16, Encoding)](fc289def-e53b-b47d-22fe-7e02b5e354e0.htm)  
[IReadWriteNetReadString(String, UInt16, Encoding)](fc289def-e53b-b47d-22fe-7e02b5e354e0.htm)

![](../icons/SectionExpanded.png)示例

以下为三菱的连接对象示例，其他的设备读写情况参照下面的代码：

String类型示例

[复制](# "复制")

```
MelsecMcNet melsec_net = new MelsecMcNet( "192.168.0.100", 6000 );

// 以下是简单的读取，没有仔细校验的方式
string d100_value = melsec_net.ReadString( "D100", 5 ).Content;

// 如果需要判断是否读取成功，使用 Unicode 编码即可读取中文，如果还是乱码，就需要自己指定编码来实现
OperateResult<string> R_d100_value = melsec_net.ReadString( "D100", 5, Encoding.Unicode );
if (R_d100_value.IsSuccess)
{
    // success
    string value = R_d100_value.Content;
}
else
{
    // failed
}
```

![](../icons/SectionExpanded.png)参见

#### 引用

[CJT188 类](7845b81d-5321-3a2e-346d-ac4fb15e3f2a.htm)

[ReadString 重载](e0ea9f3d-4c56-6078-2d10-e14579abb235.htm)

[HslCommunication.Instrument.CJT 命名空间](32a7bbe4-7443-0a40-8e16-c1df11e7a933.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadStringArray 方法 

[原文連結](http://api.hslcommunication.cn/html/c793adda-9b2c-6746-467e-b07a4f88e64a.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.CJT](../html/32a7bbe4-7443-0a40-8e16-c1df11e7a933.htm "HslCommunication.Instrument.CJT")

[CJT188 类](../html/7845b81d-5321-3a2e-346d-ac4fb15e3f2a.htm "CJT188 类")

[CJT188 方法](../html/04972e44-8c91-c7ca-0ab7-28b426216410.htm "CJT188 方法")

[ActiveDeveice 方法](../html/6de9a868-a602-a896-330c-e56955f8eb96.htm "ActiveDeveice 方法 ")

[GetNewNetMessage 方法](../html/d78644ed-1883-a2f5-91ad-741d56356d12.htm "GetNewNetMessage 方法 ")

[PackCommandWithHeader 方法](../html/84add1a8-4ab7-c216-b67e-0312b85d665e.htm "PackCommandWithHeader 方法 ")

[Read 方法](../html/1f6fea44-369c-71f2-4883-631528d6f3e5.htm "Read 方法 ")

[ReadAddress 方法](../html/81a31692-1162-4850-5404-ed02f888f324.htm "ReadAddress 方法 ")

[ReadDouble 方法](../html/b5fa76a1-104a-bdb6-c6d6-7472db8f735f.htm "ReadDouble 方法 ")

[ReadDoubleAsync 方法](../html/e53e8ff4-31bd-336c-4218-2c8fc79855c7.htm "ReadDoubleAsync 方法 ")

[ReadFloat 方法](../html/a5eeaaa8-4c0b-8d81-4420-3c7ce42a39f2.htm "ReadFloat 方法 ")

[ReadFloatAsync 方法](../html/7dd8d20e-35b5-8258-4550-df204a38e81f.htm "ReadFloatAsync 方法 ")

[ReadFromCoreServer 方法](../html/42325ee8-6793-9ad0-6b53-58f5487de0b1.htm "ReadFromCoreServer 方法 ")

[ReadString 方法](../html/e0ea9f3d-4c56-6078-2d10-e14579abb235.htm "ReadString 方法 ")

[ReadStringArray 方法](../html/c793adda-9b2c-6746-467e-b07a4f88e64a.htm "ReadStringArray 方法 ")

[ReadStringAsync 方法](../html/947017c8-6aa1-d267-3e9d-3b880b5021be.htm "ReadStringAsync 方法 ")

[ToString 方法](../html/8bff78ed-c0ad-15e4-e369-ad46dbff3241.htm "ToString 方法 ")

[Write 方法](../html/efa8452a-5af5-3ac6-b3fd-8f9c33845943.htm "Write 方法 ")

[WriteAddress 方法](../html/37e7ec52-2c99-b84f-3062-f23d3541a67a.htm "WriteAddress 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CJT188ReadStringArray 方法 |

读取指定地址的所有的字符串数据信息，一般来说，一个地址只有一个数据

**命名空间：**
 [HslCommunication.Instrument.CJT](32a7bbe4-7443-0a40-8e16-c1df11e7a933.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public OperateResult<string[]> ReadStringArray(
	string address
)
```

```
Public Function ReadStringArray ( 
	address As String
) As OperateResult(Of String())
```

```
public:
virtual OperateResult<array<String^>^>^ ReadStringArray(
	String^ address
) sealed
```

```
abstract ReadStringArray : 
        address : string -> OperateResult<string[]> 
override ReadStringArray : 
        address : string -> OperateResult<string[]>
```

#### 参数

address
:   类型：SystemString  
    数据标识，具体需要查找手册来对应

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)String  
字符串数组信息

#### 实现

[ICjt188ReadStringArray(String)](226fd33b-1900-9236-7bb7-ef5b195d6101.htm)

![](../icons/SectionExpanded.png)备注

地址可以携带地址域信息，例如 "s=2;90-1F" 或是 "s=100000;90-1F"

![](../icons/SectionExpanded.png)参见

#### 引用

[CJT188 类](7845b81d-5321-3a2e-346d-ac4fb15e3f2a.htm)

[HslCommunication.Instrument.CJT 命名空间](32a7bbe4-7443-0a40-8e16-c1df11e7a933.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadStringAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/947017c8-6aa1-d267-3e9d-3b880b5021be.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.CJT](../html/32a7bbe4-7443-0a40-8e16-c1df11e7a933.htm "HslCommunication.Instrument.CJT")

[CJT188 类](../html/7845b81d-5321-3a2e-346d-ac4fb15e3f2a.htm "CJT188 类")

[CJT188 方法](../html/04972e44-8c91-c7ca-0ab7-28b426216410.htm "CJT188 方法")

[ReadStringAsync 方法](../html/947017c8-6aa1-d267-3e9d-3b880b5021be.htm "ReadStringAsync 方法 ")

[ReadStringAsync 方法 (String, UInt16, Encoding)](../html/a6830c2f-45d7-c561-2237-f511b3d9c320.htm "ReadStringAsync 方法 (String, UInt16, Encoding)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CJT188ReadStringAsync 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法代码示例 | [ReadStringAsync(String, UInt16)](1835da81-b811-d151-0b53-907d8d87d0a8.htm) | 异步读取字符串数据，默认为最常见的ASCII编码  Asynchronously read string data, default is the most common ASCII encoding (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadStringAsync(String, UInt16, Encoding)](a6830c2f-45d7-c561-2237-f511b3d9c320.htm) | 异步使用指定的编码，读取字符串数据  Asynchronously reads string data using the specified encoding (重写 [DeviceCommunicationReadStringAsync(String, UInt16, Encoding)](091dfc67-3d25-f33e-3a55-a0228863815e.htm).) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[CJT188 类](7845b81d-5321-3a2e-346d-ac4fb15e3f2a.htm)

[HslCommunication.Instrument.CJT 命名空间](32a7bbe4-7443-0a40-8e16-c1df11e7a933.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadStringAsync 方法 (String, UInt16, Encoding)

[原文連結](http://api.hslcommunication.cn/html/a6830c2f-45d7-c561-2237-f511b3d9c320.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.CJT](../html/32a7bbe4-7443-0a40-8e16-c1df11e7a933.htm "HslCommunication.Instrument.CJT")

[CJT188 类](../html/7845b81d-5321-3a2e-346d-ac4fb15e3f2a.htm "CJT188 类")

[CJT188 方法](../html/04972e44-8c91-c7ca-0ab7-28b426216410.htm "CJT188 方法")

[ReadStringAsync 方法](../html/947017c8-6aa1-d267-3e9d-3b880b5021be.htm "ReadStringAsync 方法 ")

[ReadStringAsync 方法 (String, UInt16, Encoding)](../html/a6830c2f-45d7-c561-2237-f511b3d9c320.htm "ReadStringAsync 方法 (String, UInt16, Encoding)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CJT188ReadStringAsync 方法 (String, UInt16, Encoding) |

异步使用指定的编码，读取字符串数据  
Asynchronously reads string data using the specified encoding

**命名空间：**
 [HslCommunication.Instrument.CJT](32a7bbe4-7443-0a40-8e16-c1df11e7a933.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public override Task<OperateResult<string>> ReadStringAsync(
	string address,
	ushort length,
	Encoding encoding
)
```

```
Public Overrides Function ReadStringAsync ( 
	address As String,
	length As UShort,
	encoding As Encoding
) As Task(Of OperateResult(Of String))
```

```
public:
virtual Task<OperateResult<String^>^>^ ReadStringAsync(
	String^ address, 
	unsigned short length, 
	Encoding^ encoding
) override
```

```
abstract ReadStringAsync : 
        address : string * 
        length : uint16 * 
        encoding : Encoding -> Task<OperateResult<string>> 
override ReadStringAsync : 
        address : string * 
        length : uint16 * 
        encoding : Encoding -> Task<OperateResult<string>>
```

#### 参数

address
:   类型：SystemString  
    起始地址

length
:   类型：SystemUInt16  
    数据长度

encoding
:   类型：System.TextEncoding  
    指定的自定义的编码

#### 返回值

类型：Task[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)String  
带有成功标识的string数据

#### 实现

[IReadWriteNetReadStringAsync(String, UInt16, Encoding)](db06baf4-e674-22de-e89e-cb26646413f8.htm)  
[IReadWriteNetReadStringAsync(String, UInt16, Encoding)](db06baf4-e674-22de-e89e-cb26646413f8.htm)

![](../icons/SectionExpanded.png)示例

以下为三菱的连接对象示例，其他的设备读写情况参照下面的代码：

String类型示例

[复制](# "复制")

```
MelsecMcNet melsec_net = new MelsecMcNet( "192.168.0.100", 6000 );

// 以下是简单的读取，没有仔细校验的方式
string d100_value = melsec_net.ReadString( "D100", 5 ).Content;

// 如果需要判断是否读取成功，使用 Unicode 编码即可读取中文，如果还是乱码，就需要自己指定编码来实现
OperateResult<string> R_d100_value = melsec_net.ReadString( "D100", 5, Encoding.Unicode );
if (R_d100_value.IsSuccess)
{
    // success
    string value = R_d100_value.Content;
}
else
{
    // failed
}
```

![](../icons/SectionExpanded.png)参见

#### 引用

[CJT188 类](7845b81d-5321-3a2e-346d-ac4fb15e3f2a.htm)

[ReadStringAsync 重载](947017c8-6aa1-d267-3e9d-3b880b5021be.htm)

[HslCommunication.Instrument.CJT 命名空间](32a7bbe4-7443-0a40-8e16-c1df11e7a933.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)