# HslCommunication - HslCommunication.Core.Device

> 分類頁數: 30



---
## HslCommunication.Core.Device

[原文連結](http://api.hslcommunication.cn/html/01bea659-9e5c-248c-f839-060b8da639c8.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Core.Device](../html/01bea659-9e5c-248c-f839-060b8da639c8.htm "HslCommunication.Core.Device")

[DeviceCommunication 类](../html/d345c980-d488-c83d-72ce-f958039275db.htm "DeviceCommunication 类")

[DeviceSerialPort 类](../html/591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm "DeviceSerialPort 类")

[DeviceServer 类](../html/bad891c2-0527-acc1-3bcb-0eecb1a393a2.htm "DeviceServer 类")

[DeviceServer.DataReceivedDelegate 委托](../html/ab35401b-94ee-bb45-f94a-c14b952442ef.htm "DeviceServer.DataReceivedDelegate 委托")

[DeviceServer.DataSendDelegate 委托](../html/63ccd485-841f-ae5c-d4bd-597382878bdd.htm "DeviceServer.DataSendDelegate 委托")

[DeviceTcpNet 类](../html/09da6b9b-3727-e84e-5c46-19061382088e.htm "DeviceTcpNet 类")

[DeviceUdpNet 类](../html/225773c5-8ad0-51c5-91bc-51d651421601.htm "DeviceUdpNet 类")

[DeviceWebApi 类](../html/af832e8e-9b2b-8d0b-a923-719eafaee568.htm "DeviceWebApi 类")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslCommunication.Core.Device 命名空间 |

[缺少 "N:HslCommunication.Core.Device" 的 <summary> 文档]

![](../icons/SectionExpanded.png)类

|  | 类 | 说明 |
| --- | --- | --- |
| 公共类 | [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm) | 所有设备的基类信息 |
| 公共类 | [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm) | 串口的设备类对象信息 |
| 公共类 | [DeviceServer](bad891c2-0527-acc1-3bcb-0eecb1a393a2.htm) | 设备服务器类 |
| 公共类 | [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm) | 基于TCP管道的设备基类信息 |
| 公共类 | [DeviceUdpNet](225773c5-8ad0-51c5-91bc-51d651421601.htm) | 基于UDP/IP管道的设备基类信息 |
| 公共类 | [DeviceWebApi](af832e8e-9b2b-8d0b-a923-719eafaee568.htm) | 基于WebApi接口的设备基类 |

![](../icons/SectionExpanded.png)委托

|  | 委托 | 说明 |
| --- | --- | --- |
| 公共委托 | [DeviceServerDataReceivedDelegate](ab35401b-94ee-bb45-f94a-c14b952442ef.htm) | 当接收到来自客户的数据信息时触发的对象，该数据可能来自tcp或是串口  The object that is triggered when receiving data information from the customer, the data may come from tcp or serial port |
| 公共委托 | [DeviceServerDataSendDelegate](63ccd485-841f-ae5c-d4bd-597382878bdd.htm) | 数据发送的时候委托  Show DataSend To PLC |

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## DeviceCommunication 类

[原文連結](http://api.hslcommunication.cn/html/d345c980-d488-c83d-72ce-f958039275db.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Core.Device](../html/01bea659-9e5c-248c-f839-060b8da639c8.htm "HslCommunication.Core.Device")

[DeviceCommunication 类](../html/d345c980-d488-c83d-72ce-f958039275db.htm "DeviceCommunication 类")

[DeviceCommunication 构造函数](../html/366c4552-5077-647e-da2d-4449ff82ae09.htm "DeviceCommunication 构造函数 ")

[DeviceCommunication 属性](../html/a4a01035-8948-11b1-8a61-f2043d31db47.htm "DeviceCommunication 属性")

[DeviceCommunication 方法](../html/ec3d93f5-da3b-d238-ded9-aab29ab5a132.htm "DeviceCommunication 方法")

[DeviceCommunication 字段](../html/174ff885-ae87-fd73-9345-73ea5fd9a1f8.htm "DeviceCommunication 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DeviceCommunication 类 |

所有设备的基类信息

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  [HslCommunication.Core.NetBinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)  
    HslCommunication.Core.DeviceDeviceCommunication  
      [HslCommunication.Core.DeviceDeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)  
      [HslCommunication.Core.DeviceDeviceServer](bad891c2-0527-acc1-3bcb-0eecb1a393a2.htm)  
      [HslCommunication.Core.DeviceDeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)  
      [HslCommunication.Core.DeviceDeviceUdpNet](225773c5-8ad0-51c5-91bc-51d651421601.htm)  
      [HslCommunication.Core.DeviceDeviceWebApi](af832e8e-9b2b-8d0b-a923-719eafaee568.htm)

**命名空间：**
 [HslCommunication.Core.Device](01bea659-9e5c-248c-f839-060b8da639c8.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class DeviceCommunication : BinaryCommunication, 
	IReadWriteDevice, IReadWriteNet, IDisposable
```

```
Public Class DeviceCommunication
	Inherits BinaryCommunication
	Implements IReadWriteDevice, IReadWriteNet, IDisposable
```

```
public ref class DeviceCommunication : public BinaryCommunication, 
	IReadWriteDevice, IReadWriteNet, IDisposable
```

```
type DeviceCommunication =  
    class
        inherit BinaryCommunication
        interface IReadWriteDevice
        interface IReadWriteNet
        interface IDisposable
    end
```

DeviceCommunication 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [DeviceCommunication](366c4552-5077-647e-da2d-4449ff82ae09.htm) | 默认的无参构造函数   Default no-parameter constructor |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性代码示例 | [ByteTransform](f35118c9-3691-f752-ae6f-2c5549714faa.htm) | 当前的数据变换机制，当你需要从字节数据转换类型数据的时候需要。  The current data transformation mechanism is required when you need to convert type data from byte data. |
| 公共属性 | [CommunicationPipe](ad8e903b-d289-21b0-077e-15ede8005221.htm) | 获取或设置当前的管道信息，管道类型为[CommunicationPipe](ad8e903b-d289-21b0-077e-15ede8005221.htm)的继承类，内置了[PipeTcpNet](26b94bd5-2695-36e7-6378-c3d41846be8f.htm)管道，[PipeUdpNet](5a0b810e-187e-826b-0a13-34fd5576aa0e.htm)管道，[PipeSerialPort](9d924607-3dfc-3916-b532-c5dec48fe347.htm)管道等  Get or set the current pipeline information, the pipeline type is [CommunicationPipe](ad8e903b-d289-21b0-077e-15ede8005221.htm) inheritance class, [PipeTcpNet](26b94bd5-2695-36e7-6378-c3d41846be8f.htm) pipeline, [PipeUdpNet](5a0b810e-187e-826b-0a13-34fd5576aa0e.htm) pipeline, [PipeSerialPort](9d924607-3dfc-3916-b532-c5dec48fe347.htm) pipeline, etc (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [ConnectionId](ce1efb4b-fbc9-4683-163e-0a8e71e99d6b.htm) | 当前连接的唯一ID号，默认为长度20的guid码加随机数组成，方便列表管理，也可以自己指定  The unique ID number of the current connection. The default is a 20-digit guid code plus a random number. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性代码示例 | [LogNet](8c56f246-88ac-afb8-8c06-972fafe4f78a.htm) | 组件的日志工具，支持日志记录，只要实例化后，当前网络的基本信息，就以[DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm)等级进行输出  The component's logging tool supports logging. As long as the instantiation of the basic network information, the output will be output at [DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm) (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性代码示例 | [ReceiveTimeOut](250560a7-332b-ce1b-5eb0-e3abe4ee6012.htm) | 获取或设置接收服务器反馈的时间，如果为负数，则不接收反馈   Gets or sets the time to receive server feedback, and if it is a negative number, does not receive feedback (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [SendBeforeHex](72aa1de4-116f-abfe-4adb-814a0be0e560.htm) | 获取或设置在发送通信报文前追加发送的字节信息，HEX格式，通常用于lora组网时，需要携带 00 00 00 02 四个字节的站地址功能。  Obtain or set the byte information sent before sending communication packets, HEX format, usually used for LORA networking, you need to carry 00 00 00 02 four-byte station address function. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [SleepTime](86a5fbb9-3456-6beb-7113-5bbad5eabbc5.htm) | 获取或设置在正式接收对方返回数据前的时候，需要休息的时间，当设置为0的时候，不需要休息。  Get or set the time required to rest before officially receiving the data from the other party. When it is set to 0, no rest is required. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的属性 | [WordLength](f1a505cc-e85c-6c8b-b189-58a1ab5753f9.htm) | 一个字单位的数据表示的地址长度，西门子为2，三菱，欧姆龙，modbusTcp就为1，AB PLC无效  The address length represented by one word of data, Siemens is 2, Mitsubishi, Omron, modbusTcp is 1, AB PLC is invalid |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 受保护的方法 | [DecideWhetherQAMessage](acc40cda-e310-8537-6685-59b67eeb16ff.htm) | 决定当前的消息是否是用于问答机制返回的消息，默认直接返回 true, 实际的情况需要根据协议进行重写方法  To determine whether the current message is the message returned by the question answering mechanism, the default is true. In actual cases, the rewriting method needs to be performed according to the protocol (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [Dispose](b9dda6bf-d342-254e-689a-93fa0bb265d1.htm) | 释放被 DeviceCommunication 使用的所有资源 |
| 受保护的方法 | [Dispose(Boolean)](0bbb6ed3-111e-069c-e9f5-8cbc3f7f7b77.htm) | 释放被 DeviceCommunication 使用的非托管资源，并且是否托管资源（可选） |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | [ExtraAfterReadFromCoreServer](06a49a15-a3ce-5e39-28c7-2b78851eca2e.htm) | 和服务器交互完成的时候调用的方法，可以根据读写结果进行一些额外的操作，具体的操作需要根据实际的需求来重写实现  The method called when the interaction with the server is completed can perform some additional operations based on the read and write results. The specific operations need to be rewritten according to actual needs. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法代码示例 | [ExtraOnDisconnect](de9dbba2-3b2c-1c5e-53c8-4e6f8a44e2c4.htm) | 根据实际的协议选择是否重写本方法，有些协议在断开连接之前，需要发送一些报文来关闭当前的网络通道  Select whether to rewrite this method according to the actual protocol. Some protocols need to send some packets to close the current network channel before disconnecting. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法代码示例 | [ExtraOnDisconnectAsync](210f593e-16d8-77c4-1178-58fbcfe9b27b.htm) | 根据实际的协议选择是否重写本方法，有些协议在断开连接之前，需要发送一些报文来关闭当前的网络通道  Select whether to rewrite this method according to the actual protocol. Some protocols need to send some packets to close the current network channel before disconnecting. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 受保护的方法 | [GetLogTextFromBinary](67f3b2f7-957c-fa0b-8320-6799237157db.htm) | 获取当前的报文进行日志记录的时候，是否使用二进制的格式记录，默认返回 [LogMsgFormatBinary](2ed2ff87-caec-0e44-f08b-e2b015225a58.htm)，重写可以根据session对象分别返回不同记录模式  Whether to log the current packet in binary format, the default return is [LogMsgFormatBinary](2ed2ff87-caec-0e44-f08b-e2b015225a58.htm). If you want to override it, different recording modes can be returned according to session (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [GetNewNetMessage](16636a1b-a6f0-9c0f-9d9c-c94e26c67459.htm) | 获取一个新的消息对象的方法，需要在继承类里面进行重写  The method to get a new message object needs to be overridden in the inheritance class (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法 | [GetWordLength](c27821ea-1739-8eb9-a848-23a00828fec8.htm) | 一个字单位的数据表示的地址长度，西门子为2，三菱，欧姆龙，modbusTcp就为1，AB PLC无效  The address length represented by one word of data, Siemens is 2, Mitsubishi, Omron, modbusTcp is 1, AB PLC is invalid |
| 受保护的方法代码示例 | [InitializationOnConnect](4ab6989f-a5b5-05fe-8b44-bfa45824a2a7.htm) | 根据实际的协议选择是否重写本方法，有些协议在创建连接之后，需要进行一些初始化的信号握手，才能最终建立网络通道。  Whether to rewrite this method is based on the actual protocol. Some protocols require some initial signal handshake to establish a network channel after the connection is created. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法代码示例 | [InitializationOnConnectAsync](7da2af04-d4fc-5045-19a4-4b0d1722eef3.htm) | 根据实际的协议选择是否重写本方法，有些协议在创建连接之后，需要进行一些初始化的信号握手，才能最终建立网络通道。  Whether to rewrite this method is based on the actual protocol. Some protocols require some initial signal handshake to establish a network channel after the connection is created. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [LogRevcMessage(Byte)](dab1657f-1e75-7fe1-7003-ccedfc69855d.htm) | 使用日志记录一个接收的报文信息  Logs are used to record information about a received packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [LogRevcMessage(Byte, PipeSession)](0f51ae45-575e-4228-f978-01b37d6d5df1.htm) | 使用日志记录一个接收的报文信息  Logs are used to record information about a received packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [LogSendMessage(Byte)](4f90eecc-68a2-1ef2-44ef-369f19bed3ec.htm) | 使用日志记录一个发送的报文信息  Logs are used to record information about a send packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [LogSendMessage(Byte, PipeSession)](209caadf-4f0f-40aa-04c2-ebc61f8dab13.htm) | 使用日志记录一个发送的报文信息  Logs are used to record information about a send packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法 | [PackCommandWithHeader](d0b2cd76-c306-11fc-38c7-af8ad306cce7.htm) | 对当前的命令进行打包处理，通常是携带命令头内容，标记当前的命令的长度信息，需要进行重写，否则默认不打包  The current command is packaged, usually carrying the content of the command header, marking the length of the current command, and it needs to be rewritten, otherwise it is not packaged by default (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [Read(String, UInt16)](9dd09452-c200-bc65-75ff-69aa34b366d4.htm) | 批量读取字节数组信息，需要指定地址和长度，返回原始的字节数组  Batch read byte array information, need to specify the address and length, return the original byte array |
| 公共方法代码示例 | [ReadT](95abe064-5594-5ae7-a849-98c53f88fee0.htm) | 读取支持Hsl特性的数据内容，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明，或是官网的详细文档  Read the data content of the Hsl attribute. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. |
| 公共方法 | [ReadAsync(String, UInt16)](d8c14d37-56dc-2a31-1a19-3e2b5b65197d.htm) | 异步批量读取字节数组信息，需要指定地址和长度，返回原始的字节数组  Asynchronous batch read byte array information, need to specify the address and length, return the original byte array |
| 公共方法代码示例 | [ReadAsyncT](0522e740-4bf0-2e31-c4c5-e41384321b37.htm) | 异步读取支持Hsl特性的数据内容，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Asynchronously read the data content of the Hsl attribute. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. |
| 公共方法 | [ReadBool(String)](9b2faa2e-eaa9-7ace-c49f-7a0318f1df1a.htm) | 读取单个的Boolean数据信息  Read a single Boolean data message |
| 公共方法 | [ReadBool(String, UInt16)](69d0d24f-34a7-5b41-5c1b-619edc7e6b54.htm) | 批量读取Boolean数组信息，需要指定地址和长度，返回Boolean 数组  Batch read Boolean array information, need to specify the address and length, return Boolean array |
| 公共方法 | [ReadBoolAsync(String)](cc37da1f-d673-fe67-78ee-4e465d29a2bd.htm) | 异步读取单个的Boolean数据信息  Asynchronously read a single Boolean data message |
| 公共方法 | [ReadBoolAsync(String, UInt16)](a45d383d-0095-3ea7-32f8-e75af88d60f8.htm) | 异步批量读取Boolean数组信息，需要指定地址和长度，返回Boolean 数组  Asynchronously batch read Boolean array information, need to specify the address and length, return Boolean array |
| 公共方法代码示例 | [ReadCustomerT(String)](ca069688-d8f0-7dbd-e962-86a10dce83cf.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，返回一个新的类型的实例对象。  To read a custom data type, you need to inherit from the IDataTransfer interface and return an instance object of a new type. |
| 公共方法代码示例 | [ReadCustomerT(String, T)](ca7dc1a1-c97d-689f-9819-4dae4af67206.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，传入一个实例，对这个实例进行赋值，并返回该实例的对象。  To read a custom data type, you need to inherit from the IDataTransfer interface, pass in an instance, assign a value to this instance, and return the object of the instance. |
| 公共方法代码示例 | [ReadCustomerAsyncT(String)](f3fcda2b-255d-aadd-567e-c7b3df3636d9.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，返回一个新的类型的实例对象。  To read a custom data type, you need to inherit from the IDataTransfer interface and return an instance object of a new type. |
| 公共方法代码示例 | [ReadCustomerAsyncT(String, T)](527c9fae-4743-6ebd-7b9a-5ef3a8b7c455.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，传入一个实例，对这个实例进行赋值，并返回该实例的对象。  To read a custom data type, you need to inherit from the IDataTransfer interface, pass in an instance, assign a value to this instance, and return the object of the instance. |
| 公共方法代码示例 | [ReadDouble(String)](e7e31c43-71b4-e692-6bdc-ba9533e1d83d.htm) | 读取双浮点的数据  Read double floating point data |
| 公共方法代码示例 | [ReadDouble(String, UInt16)](4eee0f05-f861-4b48-3f06-99737fb24aee.htm) | 读取双浮点数据的数组  Read double floating point data array |
| 公共方法代码示例 | [ReadDoubleAsync(String)](2cf895fa-1033-1def-f048-0e8caefaa71d.htm) | 异步读取双浮点的数据  Asynchronously read double floating point data |
| 公共方法代码示例 | [ReadDoubleAsync(String, UInt16)](2256d406-d4d5-da54-210b-ccef58a5e914.htm) | 异步读取双浮点数据的数组  Asynchronously read double floating point data array |
| 公共方法代码示例 | [ReadFloat(String)](bc2801ae-6600-e845-7c6c-522c9156eb0e.htm) | 读取单浮点数据  Read single floating point data |
| 公共方法代码示例 | [ReadFloat(String, UInt16)](6a9fb098-b13c-313b-3b27-20d3c787945d.htm) | 读取单浮点精度的数组  Read single floating point array |
| 公共方法代码示例 | [ReadFloatAsync(String)](89f743f8-4b3a-7085-1e49-a1d7adb10052.htm) | 异步读取单浮点数据  Asynchronously read single floating point data |
| 公共方法代码示例 | [ReadFloatAsync(String, UInt16)](cb6d5c5d-9b5c-2467-3e6b-5efa9d20cd72.htm) | 异步读取单浮点精度的数组  Asynchronously read single floating point array |
| 公共方法代码示例 | [ReadFromCoreServer(Byte)](8146dc0c-1edc-3ca7-428a-a3bdbd35e6c5.htm) | 将二进制的数据发送到管道中去，然后从管道里接收二进制的数据回来，并返回是否成功的结果对象。  Send binary data to the pipeline, and then receive binary data back from the pipeline, and return whether the success of the result object (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServer(IEnumerableByte)](044cfdd9-ec76-b89f-2f2b-5f439025006c.htm) | 将多个数据报文按顺序发到设备，并从设备接收返回的数据内容，然后拼接成一个Byte[]信息，需要重写[UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm)方法才能返回正确的结果。  Send multiple data packets to the device in sequence, and receive the returned data content from the device, and then splicing them into a Byte[] message, you need to rewrite [UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm) method to return the correct result. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServer(Byte, Boolean, Boolean)](b32ce8e8-a959-fcc2-6736-14b39fdb9a40.htm) | 将二进制的数据发送到管道中去，然后从管道里接收二进制的数据回来，并返回是否成功的结果对象。  Send binary data to the pipeline, and then receive binary data back from the pipeline, and return whether the success of the result object (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServer(CommunicationPipe, Byte, Boolean, Boolean)](ef36d56b-df09-f727-e8d3-316d37a13ee9.htm) | 使用指定的管道来进行数据通信，发送原始数据到管道，然后从管道接收相关的数据返回，本方法无锁 (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServerAsync(Byte)](a4de5a24-c67a-7b02-db05-33b85eaf6e1e.htm) | 将二进制的数据发送到管道中去，然后从管道里接收二进制的数据回来，并返回是否成功的结果对象。  Send binary data to the pipeline, and then receive binary data back from the pipeline, and return whether the success of the result object (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServerAsync(IEnumerableByte)](bcdb59b7-a70a-d5c0-1fc6-9ececcd39ea2.htm) | 将多个数据报文按顺序发到设备，并从设备接收返回的数据内容，然后拼接成一个Byte[]信息，需要重写[UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm)方法才能返回正确的结果。  Send multiple data packets to the device in sequence, and receive the returned data content from the device, and then splicing them into a Byte[] message, you need to rewrite [UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm) method to return the correct result. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServerAsync(Byte, Boolean, Boolean)](1c79c204-5ab5-9342-e65d-7455c20ea577.htm) | 将二进制的数据发送到管道中去，然后从管道里接收二进制的数据回来，并返回是否成功的结果对象。  Send binary data to the pipeline, and then receive binary data back from the pipeline, and return whether the success of the result object (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServerAsync(CommunicationPipe, Byte, Boolean, Boolean)](a254a19d-3c29-147b-335b-8adccdf4ee2f.htm) | 使用指定的管道来进行数据通信，发送原始数据到管道，然后从管道接收相关的数据返回，本方法无锁 (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法代码示例 | [ReadInt16(String)](57a16003-2fd5-b2bd-09bf-14366be899b8.htm) | 读取16位的有符号的整型数据  Read 16-bit signed integer data |
| 公共方法代码示例 | [ReadInt16(String, UInt16)](b2bde158-f51c-93fc-66ca-5e4cf6ea1fb7.htm) | 读取16位的有符号整型数组  Read 16-bit signed integer array |
| 公共方法代码示例 | [ReadInt16Async(String)](ac021793-890d-ceab-eae4-47ff763390c1.htm) | 异步读取16位的有符号的整型数据  Asynchronously read 16-bit signed integer data |
| 公共方法代码示例 | [ReadInt16Async(String, UInt16)](f30e02fc-8078-f1b3-c226-5991e3e23caf.htm) | 异步读取16位的有符号整型数组  Asynchronously read 16-bit signed integer array |
| 公共方法代码示例 | [ReadInt32(String)](bfbea035-9a83-40e3-084c-28122041b37a.htm) | 读取32位的有符号整型  Read 32-bit signed integer |
| 公共方法代码示例 | [ReadInt32(String, UInt16)](095240a3-bc0f-8647-4867-88e40cc089e3.htm) | 读取32位有符号整型数组  Read 32-bit signed integer array |
| 公共方法代码示例 | [ReadInt32Async(String)](96726ea0-82aa-1229-21f6-43cd62abca43.htm) | 异步读取32位的有符号整型  Asynchronously read 32-bit signed integer |
| 公共方法代码示例 | [ReadInt32Async(String, UInt16)](d6356095-a289-18eb-136c-b7e6ae099616.htm) | 异步读取32位有符号整型数组  Asynchronously read 32-bit signed integer array |
| 公共方法代码示例 | [ReadInt64(String)](0d0c2284-2da1-6f9a-b353-537561cd5e95.htm) | 读取64位的有符号整型  Read 64-bit signed integer |
| 公共方法代码示例 | [ReadInt64(String, UInt16)](35a55019-efb1-bf2d-be1b-01c8eed99243.htm) | 读取64位的有符号整型数组  Read 64-bit signed integer array |
| 公共方法代码示例 | [ReadInt64Async(String)](26d46c8a-3775-3abd-d3ec-5cd1758a317e.htm) | 异步读取64位的有符号整型  Asynchronously read 64-bit signed integer |
| 公共方法代码示例 | [ReadInt64Async(String, UInt16)](8dc97c4a-8904-9a34-9ddc-76c1a9477ee5.htm) | 异步读取64位的有符号整型数组  Asynchronously read 64-bit signed integer array |
| 公共方法代码示例 | [ReadString(String, UInt16)](07812d76-415a-6400-c8a2-1c21da51e7a3.htm) | 读取字符串数据，默认为最常见的ASCII编码  Read string data, default is the most common ASCII encoding |
| 公共方法代码示例 | [ReadString(String, UInt16, Encoding)](b34c71a6-58f3-6bf4-0b93-471be4f2ed31.htm) | 使用指定的编码，读取字符串数据  Reads string data using the specified encoding |
| 公共方法代码示例 | [ReadStringAsync(String, UInt16)](1835da81-b811-d151-0b53-907d8d87d0a8.htm) | 异步读取字符串数据，默认为最常见的ASCII编码  Asynchronously read string data, default is the most common ASCII encoding |
| 公共方法代码示例 | [ReadStringAsync(String, UInt16, Encoding)](091dfc67-3d25-f33e-3a55-a0228863815e.htm) | 异步使用指定的编码，读取字符串数据  Asynchronously reads string data using the specified encoding |
| 公共方法代码示例 | [ReadStructT](49853082-fbad-52e5-d756-0615cedb4b83.htm) | 读取结构体类型的数据，根据结构体自身的定义，读取原始字节数组，然后解析出实际的结构体数据，结构体需要实现[HslStructAttribute](87e427a8-0c7b-43a8-c421-c27264eda90f.htm)特性 |
| 公共方法代码示例 | [ReadStructAsyncT](36461e61-d0ed-62cf-79ce-dcc7f71a3fa2.htm) | 读取结构体类型的数据，根据结构体自身的定义，读取原始字节数组，然后解析出实际的结构体数据，结构体需要实现[HslStructAttribute](87e427a8-0c7b-43a8-c421-c27264eda90f.htm)特性 |
| 公共方法代码示例 | [ReadUInt16(String)](736379d1-9560-5052-658a-13c126d56d78.htm) | 读取16位的无符号整型  Read 16-bit unsigned integer |
| 公共方法代码示例 | [ReadUInt16(String, UInt16)](e284c4bd-d3f8-ff30-eb42-0aa035cfb4f6.htm) | 读取16位的无符号整型数组  Read 16-bit unsigned integer array |
| 公共方法代码示例 | [ReadUInt16Async(String)](361acc93-43b5-a9a5-6fa7-fa8728b7b96c.htm) | 异步读取16位的无符号整型  Asynchronously read 16-bit unsigned integer |
| 公共方法代码示例 | [ReadUInt16Async(String, UInt16)](73dd24b8-b501-df53-1de6-ccb2708f656d.htm) | 异步读取16位的无符号整型数组  Asynchronously read 16-bit unsigned integer array |
| 公共方法代码示例 | [ReadUInt32(String)](ca73c958-39c7-9b34-aa44-3db2d9b0cad3.htm) | 读取32位的无符号整型  Read 32-bit unsigned integer |
| 公共方法代码示例 | [ReadUInt32(String, UInt16)](878b4709-0278-fcdf-3b65-afbf0f496c1e.htm) | 读取32位的无符号整型数组  Read 32-bit unsigned integer array |
| 公共方法代码示例 | [ReadUInt32Async(String)](2343953f-6350-331e-6ba7-18a4449f0575.htm) | 异步读取32位的无符号整型  Asynchronously read 32-bit unsigned integer |
| 公共方法代码示例 | [ReadUInt32Async(String, UInt16)](e0eba6c1-56a2-5190-ba03-66116bcafe3d.htm) | 异步读取32位的无符号整型数组  Asynchronously read 32-bit unsigned integer array |
| 公共方法代码示例 | [ReadUInt64(String)](947487dd-d4a3-72d8-d9fb-2224aad70e12.htm) | 读取64位的无符号整型  Read 64-bit unsigned integer |
| 公共方法代码示例 | [ReadUInt64(String, UInt16)](858766a9-9d62-87a8-f727-a00838f3d2b6.htm) | 读取64位的无符号整型的数组  Read 64-bit unsigned integer array |
| 公共方法代码示例 | [ReadUInt64Async(String)](3411ff3a-772f-dbe1-8931-75e00a731579.htm) | 异步读取64位的无符号整型  Asynchronously read 64-bit unsigned integer |
| 公共方法代码示例 | [ReadUInt64Async(String, UInt16)](c2095a6e-ff89-dd98-8a21-fe72f82dbbc6.htm) | 异步读取64位的无符号整型的数组  Asynchronously read 64-bit unsigned integer array |
| 公共方法 | [SetDtuPipe](f765d1d7-234a-b229-9457-97e25e024dc9.htm) | 将当前的通信对象设置DTU模式，允许传入现成的管道，并返回初始化结果，如果该设备重写了握手报文，就是返回握手结果  Set the current communication object to DTU mode, allow the existing pipe to be passed in, and return the initialization result, if the device rewrites the handshake packet, the handshake result is returned (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [SetDtuPipeAsync](9d1fd12d-b804-3828-a627-a2c0eef1d23f.htm) | 将当前的通信对象设置DTU模式，允许传入现成的管道，并返回初始化结果，如果该设备重写了握手报文，就是返回握手结果  Set the current communication object to DTU mode, allow the existing pipe to be passed in, and return the initialization result, if the device rewrites the handshake packet, the handshake result is returned (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ToString](5667665a-da16-d227-f350-5b3f0ab9ab04.htm) | (重写 ObjectToString.) |
| 公共方法 | [UnpackResponseContent](208017a7-8da2-33f8-3273-d7e7e850100f.htm) | 根据对方返回的报文命令，对命令进行基本的拆包，例如各种Modbus协议拆包为统一的核心报文，还支持对报文的验证  According to the message command returned by the other party, the command is basically unpacked, for example, various Modbus protocols are unpacked into a unified core message, and the verification of the message is also supported (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [Wait(String, Boolean, Int32, Int32)](0e4db26e-1177-2b18-9492-ba05afb6fcc4.htm) | 等待指定地址的Boolean值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the Boolean value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. |
| 公共方法 | [Wait(String, Int16, Int32, Int32)](e06fe098-8c8b-83f9-73ca-2a7cce19d43d.htm) | 等待指定地址的Int16值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the Int16 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. |
| 公共方法 | [Wait(String, Int32, Int32, Int32)](1bce0aae-27e9-f1d6-e8f3-25ae6b47554e.htm) | 等待指定地址的Int32值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the Int32 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. |
| 公共方法 | [Wait(String, Int64, Int32, Int32)](073354e9-2508-2158-1b01-881c3a34844b.htm) | 等待指定地址的Int64值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the Int64 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. |
| 公共方法 | [Wait(String, UInt16, Int32, Int32)](d97f8d26-1c1a-df69-93f0-f60209597ced.htm) | 等待指定地址的UInt16值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the UInt16 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. |
| 公共方法 | [Wait(String, UInt32, Int32, Int32)](1e61c395-b07f-e66c-5f9c-2708093150d9.htm) | 等待指定地址的UInt32值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the UInt32 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. |
| 公共方法 | [Wait(String, UInt64, Int32, Int32)](d659729b-5766-c497-d418-856c34ec2b24.htm) | 等待指定地址的UInt64值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the UInt64 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. |
| 公共方法 | [WaitAsync(String, Boolean, Int32, Int32)](65399e87-9d38-f819-ccd3-005f68c2c3b3.htm) | 等待指定地址的Boolean值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the Boolean value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. |
| 公共方法 | [WaitAsync(String, Int16, Int32, Int32)](01277e39-afb1-86a4-3aa3-560160e66fb9.htm) | 等待指定地址的Int16值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the Int16 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. |
| 公共方法 | [WaitAsync(String, Int32, Int32, Int32)](4a6d380f-6b2a-1270-07a7-39b88a2712fd.htm) | 等待指定地址的Int32值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the Int32 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. |
| 公共方法 | [WaitAsync(String, Int64, Int32, Int32)](d7eede9b-4d44-479b-fb88-3ff2e1bb7568.htm) | 等待指定地址的Int64值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the Int64 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. |
| 公共方法 | [WaitAsync(String, UInt16, Int32, Int32)](0a8395c7-dc7f-4d5f-4316-98e18fe70f4c.htm) | 等待指定地址的UInt16值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the UInt16 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. |
| 公共方法 | [WaitAsync(String, UInt32, Int32, Int32)](19e9074e-99a2-99fd-91ef-a5314a115397.htm) | 等待指定地址的UInt32值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the UInt32 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. |
| 公共方法 | [WaitAsync(String, UInt64, Int32, Int32)](5109dac0-0ba2-cbb9-e383-1a869c6dc22e.htm) | 等待指定地址的UInt64值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the UInt64 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. |
| 公共方法 | [Write(String, Boolean)](06d0b9af-7f28-0b9c-098f-4da28e69504a.htm) | 写入单个的Boolean数据，返回是否成功  Write a single Boolean data, and return whether the write was successful |
| 公共方法 | [Write(String, Boolean)](5b046c82-3d8e-550b-fdf3-41ca5558169b.htm) | 批量写入Boolean数组数据，返回是否成功  Batch write Boolean array data, return whether the write was successful |
| 公共方法 | [Write(String, Byte)](b5536dd6-ccde-6883-ac27-bff9e4b37806.htm) | 写入原始的byte数组数据到指定的地址，返回是否写入成功  Write the original byte array data to the specified address, and return whether the write was successful |
| 公共方法代码示例 | [Write(String, Double)](64c44576-7e95-cbcb-a8c9-e6b7ff7e3c50.htm) | 写入double数据，返回是否成功  Write double data, return whether the write was successful |
| 公共方法代码示例 | [Write(String, Double)](df56e6b1-0688-ece7-3bf1-9db32fec9f6e.htm) | 写入double数组，返回是否成功  Write double array, return whether the write was successful |
| 公共方法代码示例 | [Write(String, Int16)](0e23506e-b604-e67c-b5eb-06271b6c6335.htm) | 写入short数据，返回是否成功  Write short data, returns whether success |
| 公共方法代码示例 | [Write(String, Int16)](c7cfd67d-ca31-3087-c684-ee4c10780ab1.htm) | 写入short数组，返回是否成功  Write short array, return whether the write was successful |
| 公共方法代码示例 | [Write(String, Int32)](ba80d1a5-c8d8-7e0b-5c6b-9435cf20cadb.htm) | 写入int数据，返回是否成功  Write int data, return whether the write was successful |
| 公共方法代码示例 | [Write(String, Int32)](57dac21c-8b28-de7c-c213-77cda465cccb.htm) | 写入int[]数组，返回是否成功  Write int array, return whether the write was successful |
| 公共方法代码示例 | [Write(String, Int64)](ba246b33-5a12-a22c-6d40-0969a64b508b.htm) | 写入long数据，返回是否成功  Write long data, return whether the write was successful |
| 公共方法代码示例 | [Write(String, Int64)](5b910f9c-6e6b-bb9d-6842-1b496640c1f4.htm) | 写入long数组，返回是否成功  Write long array, return whether the write was successful |
| 公共方法代码示例 | [Write(String, Single)](a2265fa7-5b56-a5aa-9d21-a0dc5458cbf9.htm) | 写入float数据，返回是否成功  Write float data, return whether the write was successful |
| 公共方法代码示例 | [Write(String, Single)](bb76227e-3e95-8ca2-5dcf-c087d90b8a01.htm) | 写入float数组，返回是否成功  Write float array, return whether the write was successful |
| 公共方法代码示例 | [Write(String, String)](4e0ec363-ba64-7ca7-b8d3-2e9a8087bd55.htm) | 写入字符串信息，编码为ASCII  Write string information, encoded as ASCII |
| 公共方法代码示例 | [Write(String, UInt16)](e5bc9466-65b4-fc64-e0d0-1867b1113130.htm) | 写入ushort数据，返回是否成功  Write ushort data, return whether the write was successful |
| 公共方法代码示例 | [Write(String, UInt16)](65544962-82ac-bb47-d72f-5dd0c3d4ac6c.htm) | 写入ushort数组，返回是否成功  Write ushort array, return whether the write was successful |
| 公共方法代码示例 | [Write(String, UInt32)](c99d8246-f95e-4f79-738b-16e622072eb4.htm) | 写入uint数据，返回是否成功  Write uint data, return whether the write was successful |
| 公共方法代码示例 | [Write(String, UInt32)](b80b3099-b34d-fde8-6a65-837da3bdc5ab.htm) | 写入uint[]数组，返回是否成功  Write uint array, return whether the write was successful |
| 公共方法代码示例 | [Write(String, UInt64)](4924086e-9b36-d446-4ae4-6682134a6ad1.htm) | 写入ulong数据，返回是否成功  Write ulong data, return whether the write was successful |
| 公共方法代码示例 | [Write(String, UInt64)](b1131bdd-8b51-ec6d-dd0e-fec1f6cb51b6.htm) | 写入ulong数组，返回是否成功  Write ulong array, return whether the write was successful |
| 公共方法代码示例 | [Write(String, String, Int32)](02fe2ca9-49a3-ede4-d1af-8c6da1f644da.htm) | 写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为ASCII  Write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, it is filled with 0 and the encoding is ASCII. |
| 公共方法代码示例 | [Write(String, String, Encoding)](7b602661-902d-2b70-1fdb-ba87ecd307bf.htm) | 写入字符串信息，需要指定的编码信息  Write string information, need to specify the encoding information |
| 公共方法代码示例 | [Write(String, String, Int32, Encoding)](e8cb30d9-4c5e-6c68-b9af-1b1b6c2dc714.htm) | 写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为指定的编码信息  Write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, then the operation is complemented with 0 , you should specified the encoding information |
| 公共方法代码示例 | [WriteT(T)](9f48853f-fbeb-c37d-25a2-1c76d1dee407.htm) | 写入支持Hsl特性的数据，返回是否写入成功，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Write data that supports the Hsl attribute, and return whether the write was successful. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. |
| 公共方法 | [WriteAsync(String, Boolean)](10fe669c-805e-1e5a-f016-f319ad363e90.htm) | 异步批量写入Boolean数组数据，返回是否成功  Asynchronously batch write Boolean array data, return success |
| 公共方法 | [WriteAsync(String, Boolean)](fe872dae-7ef9-12c1-d883-e92ba8676416.htm) | 异步批量写入Boolean数组数据，返回是否成功  Asynchronously batch write Boolean array data, return success |
| 公共方法代码示例 | [WriteAsync(String, Byte)](8f1cec30-8b86-6749-c8f3-7d39168d59df.htm) | 异步写入原始的byte数组数据到指定的地址，返回是否写入成功  Asynchronously writes the original byte array data to the specified address, and returns whether the write was successful |
| 公共方法代码示例 | [WriteAsync(String, Double)](068adaa4-d1de-dffc-8376-6f483e2c1d96.htm) | 异步写入double数据，返回是否成功  Asynchronously write double data, return whether the write was successful |
| 公共方法代码示例 | [WriteAsync(String, Double)](36199b83-c4ea-b02f-3675-fb6ac48d0c59.htm) | 异步写入double数组，返回是否成功  Asynchronously write double array, return whether the write was successful |
| 公共方法代码示例 | [WriteAsync(String, Int16)](ccd95ccc-10a3-af67-62e8-c2312b4871b0.htm) | 异步写入short数据，返回是否成功  Asynchronously write short data, returns whether success |
| 公共方法代码示例 | [WriteAsync(String, Int16)](2625e868-7096-e526-ce25-b29afc639b41.htm) | 异步写入short数组，返回是否成功  Asynchronously write short array, return whether the write was successful |
| 公共方法代码示例 | [WriteAsync(String, Int32)](df91d629-2080-328d-67cd-fbefb9a45404.htm) | 异步写入int数据，返回是否成功  Asynchronously write int data, return whether the write was successful |
| 公共方法代码示例 | [WriteAsync(String, Int32)](99962a08-9624-cea6-6009-5e6f207236b7.htm) | 异步写入int[]数组，返回是否成功  Asynchronously write int array, return whether the write was successful |
| 公共方法代码示例 | [WriteAsync(String, Int64)](ec6c1e24-a0d8-f1f5-56af-109b0c546fc8.htm) | 异步写入long数据，返回是否成功  Asynchronously write long data, return whether the write was successful |
| 公共方法代码示例 | [WriteAsync(String, Int64)](4b2ab94c-2891-e0a9-f10b-b67c63bd1cf9.htm) | 异步写入long数组，返回是否成功  Asynchronously write long array, return whether the write was successful |
| 公共方法代码示例 | [WriteAsync(String, Single)](910189df-13f9-b7d0-a626-b5b25ce5d32a.htm) | 异步写入float数据，返回是否成功  Asynchronously write float data, return whether the write was successful |
| 公共方法代码示例 | [WriteAsync(String, Single)](dbb2359a-fb54-cf93-4413-20d5a4ef9d36.htm) | 异步写入float数组，返回是否成功  Asynchronously write float array, return whether the write was successful |
| 公共方法代码示例 | [WriteAsync(String, String)](fe839c9a-46f3-e317-0155-b10a8a83a9b1.htm) | 异步写入字符串信息，编码为ASCII  Asynchronously write string information, encoded as ASCII |
| 公共方法代码示例 | [WriteAsync(String, UInt16)](a04f3aa2-a51e-26cd-974e-bf70a2db359d.htm) | 异步写入ushort数据，返回是否成功  Asynchronously write ushort data, return whether the write was successful |
| 公共方法代码示例 | [WriteAsync(String, UInt16)](34a1a321-972e-06a8-76fb-794dd1229005.htm) | 异步写入ushort数组，返回是否成功  Asynchronously write ushort array, return whether the write was successful |
| 公共方法代码示例 | [WriteAsync(String, UInt32)](abad215c-bc99-ac90-b9f3-7adb3a7cca8f.htm) | 异步写入uint数据，返回是否成功  Asynchronously write uint data, return whether the write was successful |
| 公共方法代码示例 | [WriteAsync(String, UInt32)](5b655bf3-5ba7-ad31-6897-89c59051ac60.htm) | 异步写入uint[]数组，返回是否成功  Asynchronously write uint array, return whether the write was successful |
| 公共方法代码示例 | [WriteAsync(String, UInt64)](e396074d-6cf1-b45d-5e8e-01026f6e391b.htm) | 异步写入ulong数据，返回是否成功  Asynchronously write ulong data, return whether the write was successful |
| 公共方法代码示例 | [WriteAsync(String, UInt64)](ae3df8ec-29c0-42ff-3d6a-68f496997b3d.htm) | 异步写入ulong数组，返回是否成功  Asynchronously write ulong array, return whether the write was successful |
| 公共方法 | [WriteAsync(String, String, Int32)](b739b983-b03c-8ecf-2713-75f19233ae2f.htm) | 异步写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为ASCII  Asynchronously write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, it is filled with 0 and the encoding is ASCII. |
| 公共方法代码示例 | [WriteAsync(String, String, Encoding)](876e9800-0406-6407-bdf8-8ab1cf28df38.htm) | 异步写入字符串信息，需要指定的编码信息  Asynchronously write string information, need to specify the encoding information |
| 公共方法 | [WriteAsync(String, String, Int32, Encoding)](ae9fb2a3-fddb-48eb-e9ad-6cc0ebbac1f6.htm) | 异步写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为指定的编码信息  Asynchronously write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, then the operation is complemented with 0 , you should specified the encoding information |
| 公共方法代码示例 | [WriteAsyncT(T)](25d96e77-d1b2-918c-33f9-e5e4cc57c9cd.htm) | 异步写入支持Hsl特性的数据，返回是否写入成功，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Asynchronously write data that supports the Hsl attribute, and return whether the write was successful. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. |
| 公共方法代码示例 | [WriteCustomerT](84a37f82-3f67-a0f8-fd3f-0630772565e4.htm) | 写入自定义类型的数据，该类型必须继承自IDataTransfer接口  Write data of a custom type, which must inherit from the IDataTransfer interface |
| 公共方法代码示例 | [WriteCustomerAsyncT](04628397-cfa1-d785-7790-4c02cb377ab6.htm) | 写入自定义类型的数据，该类型必须继承自IDataTransfer接口  Write data of a custom type, which must inherit from the IDataTransfer interface |

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

[HslCommunication.Core.Device 命名空间](01bea659-9e5c-248c-f839-060b8da639c8.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## DeviceCommunication 构造函数 

[原文連結](http://api.hslcommunication.cn/html/366c4552-5077-647e-da2d-4449ff82ae09.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Core.Device](../html/01bea659-9e5c-248c-f839-060b8da639c8.htm "HslCommunication.Core.Device")

[DeviceCommunication 类](../html/d345c980-d488-c83d-72ce-f958039275db.htm "DeviceCommunication 类")

[DeviceCommunication 构造函数](../html/366c4552-5077-647e-da2d-4449ff82ae09.htm "DeviceCommunication 构造函数 ")

[DeviceCommunication 属性](../html/a4a01035-8948-11b1-8a61-f2043d31db47.htm "DeviceCommunication 属性")

[DeviceCommunication 方法](../html/ec3d93f5-da3b-d238-ded9-aab29ab5a132.htm "DeviceCommunication 方法")

[DeviceCommunication 字段](../html/174ff885-ae87-fd73-9345-73ea5fd9a1f8.htm "DeviceCommunication 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DeviceCommunication 构造函数 |

默认的无参构造函数   
Default no-parameter constructor

**命名空间：**
 [HslCommunication.Core.Device](01bea659-9e5c-248c-f839-060b8da639c8.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public DeviceCommunication()
```

```
Public Sub New
```

```
public:
DeviceCommunication()
```

```
new : unit -> DeviceCommunication
```

![](../icons/SectionExpanded.png)参见

#### 引用

[DeviceCommunication 类](d345c980-d488-c83d-72ce-f958039275db.htm)

[HslCommunication.Core.Device 命名空间](01bea659-9e5c-248c-f839-060b8da639c8.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## DeviceCommunication 属性

[原文連結](http://api.hslcommunication.cn/html/a4a01035-8948-11b1-8a61-f2043d31db47.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Core.Device](../html/01bea659-9e5c-248c-f839-060b8da639c8.htm "HslCommunication.Core.Device")

[DeviceCommunication 类](../html/d345c980-d488-c83d-72ce-f958039275db.htm "DeviceCommunication 类")

[DeviceCommunication 属性](../html/a4a01035-8948-11b1-8a61-f2043d31db47.htm "DeviceCommunication 属性")

[ByteTransform 属性](../html/f35118c9-3691-f752-ae6f-2c5549714faa.htm "ByteTransform 属性 ")

[WordLength 属性](../html/f1a505cc-e85c-6c8b-b189-58a1ab5753f9.htm "WordLength 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DeviceCommunication 属性 |

[DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性代码示例 | [ByteTransform](f35118c9-3691-f752-ae6f-2c5549714faa.htm) | 当前的数据变换机制，当你需要从字节数据转换类型数据的时候需要。  The current data transformation mechanism is required when you need to convert type data from byte data. |
| 公共属性 | [CommunicationPipe](ad8e903b-d289-21b0-077e-15ede8005221.htm) | 获取或设置当前的管道信息，管道类型为[CommunicationPipe](ad8e903b-d289-21b0-077e-15ede8005221.htm)的继承类，内置了[PipeTcpNet](26b94bd5-2695-36e7-6378-c3d41846be8f.htm)管道，[PipeUdpNet](5a0b810e-187e-826b-0a13-34fd5576aa0e.htm)管道，[PipeSerialPort](9d924607-3dfc-3916-b532-c5dec48fe347.htm)管道等  Get or set the current pipeline information, the pipeline type is [CommunicationPipe](ad8e903b-d289-21b0-077e-15ede8005221.htm) inheritance class, [PipeTcpNet](26b94bd5-2695-36e7-6378-c3d41846be8f.htm) pipeline, [PipeUdpNet](5a0b810e-187e-826b-0a13-34fd5576aa0e.htm) pipeline, [PipeSerialPort](9d924607-3dfc-3916-b532-c5dec48fe347.htm) pipeline, etc (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [ConnectionId](ce1efb4b-fbc9-4683-163e-0a8e71e99d6b.htm) | 当前连接的唯一ID号，默认为长度20的guid码加随机数组成，方便列表管理，也可以自己指定  The unique ID number of the current connection. The default is a 20-digit guid code plus a random number. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性代码示例 | [LogNet](8c56f246-88ac-afb8-8c06-972fafe4f78a.htm) | 组件的日志工具，支持日志记录，只要实例化后，当前网络的基本信息，就以[DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm)等级进行输出  The component's logging tool supports logging. As long as the instantiation of the basic network information, the output will be output at [DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm) (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性代码示例 | [ReceiveTimeOut](250560a7-332b-ce1b-5eb0-e3abe4ee6012.htm) | 获取或设置接收服务器反馈的时间，如果为负数，则不接收反馈   Gets or sets the time to receive server feedback, and if it is a negative number, does not receive feedback (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [SendBeforeHex](72aa1de4-116f-abfe-4adb-814a0be0e560.htm) | 获取或设置在发送通信报文前追加发送的字节信息，HEX格式，通常用于lora组网时，需要携带 00 00 00 02 四个字节的站地址功能。  Obtain or set the byte information sent before sending communication packets, HEX format, usually used for LORA networking, you need to carry 00 00 00 02 four-byte station address function. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [SleepTime](86a5fbb9-3456-6beb-7113-5bbad5eabbc5.htm) | 获取或设置在正式接收对方返回数据前的时候，需要休息的时间，当设置为0的时候，不需要休息。  Get or set the time required to rest before officially receiving the data from the other party. When it is set to 0, no rest is required. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的属性 | [WordLength](f1a505cc-e85c-6c8b-b189-58a1ab5753f9.htm) | 一个字单位的数据表示的地址长度，西门子为2，三菱，欧姆龙，modbusTcp就为1，AB PLC无效  The address length represented by one word of data, Siemens is 2, Mitsubishi, Omron, modbusTcp is 1, AB PLC is invalid |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[DeviceCommunication 类](d345c980-d488-c83d-72ce-f958039275db.htm)

[HslCommunication.Core.Device 命名空间](01bea659-9e5c-248c-f839-060b8da639c8.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ByteTransform 属性 

[原文連結](http://api.hslcommunication.cn/html/f35118c9-3691-f752-ae6f-2c5549714faa.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Core.Device](../html/01bea659-9e5c-248c-f839-060b8da639c8.htm "HslCommunication.Core.Device")

[DeviceCommunication 类](../html/d345c980-d488-c83d-72ce-f958039275db.htm "DeviceCommunication 类")

[DeviceCommunication 属性](../html/a4a01035-8948-11b1-8a61-f2043d31db47.htm "DeviceCommunication 属性")

[ByteTransform 属性](../html/f35118c9-3691-f752-ae6f-2c5549714faa.htm "ByteTransform 属性 ")

[WordLength 属性](../html/f1a505cc-e85c-6c8b-b189-58a1ab5753f9.htm "WordLength 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DeviceCommunicationByteTransform 属性 |

当前的数据变换机制，当你需要从字节数据转换类型数据的时候需要。  
The current data transformation mechanism is required when you need to convert type data from byte data.

**命名空间：**
 [HslCommunication.Core.Device](01bea659-9e5c-248c-f839-060b8da639c8.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public IByteTransform ByteTransform { get; set; }
```

```
Public Property ByteTransform As IByteTransform
	Get
	Set
```

```
public:
virtual property IByteTransform^ ByteTransform {
	IByteTransform^ get () sealed;
	void set (IByteTransform^ value) sealed;
}
```

```
abstract ByteTransform : IByteTransform with get, set
override ByteTransform : IByteTransform with get, set
```

#### 属性值

类型：[IByteTransform](56c55574-bb2a-fe66-e7d5-1332a1bc18f0.htm)

#### 实现

[IReadWriteDeviceByteTransform](0836ec33-d464-ed39-37ff-b2ada0d286fc.htm)

![](../icons/SectionExpanded.png)备注

在HSL里提供了三种数据变换机制，分别是 [RegularByteTransform](4a51ffeb-6068-0331-b842-19f7ec0133f0.htm), [ReverseBytesTransform](deb76ea7-f43b-2772-aca1-7044b05710c8.htm),
[ReverseWordTransform](0541a2c1-e245-0d53-775b-ebd826119edc.htm)，各自的[DataFormat](4e3d381b-5b6d-f012-7658-19a29e76e059.htm)属性也可以自定调整，基本满足所有的情况使用。  
Three data transformation mechanisms are provided in HSL, namely [RegularByteTransform](4a51ffeb-6068-0331-b842-19f7ec0133f0.htm), [ReverseBytesTransform](deb76ea7-f43b-2772-aca1-7044b05710c8.htm),
[ReverseWordTransform](0541a2c1-e245-0d53-775b-ebd826119edc.htm), and their respective [DataFormat](4e3d381b-5b6d-f012-7658-19a29e76e059.htm) property can also be adjusted by itself, basically satisfying all situations.

![](../icons/SectionExpanded.png)示例

主要是用来转换数据类型的，下面仅仅演示了2个方法，其他的类型转换，类似处理。

ByteTransform示例

[复制](# "复制")

```
// 假设buffer是client从设备读取的数据内容
byte[] buffer = new byte[8];
// 转化为4个short
short[] short_value = client.ByteTransform.TransInt16( buffer, 0, 4 );
// 转化为2个float
float[] float_value = client.ByteTransform.TransSingle( buffer, 0, 2 );
// 其他的类型转换是类似的
```

![](../icons/SectionExpanded.png)参见

#### 引用

[DeviceCommunication 类](d345c980-d488-c83d-72ce-f958039275db.htm)

[HslCommunication.Core.Device 命名空间](01bea659-9e5c-248c-f839-060b8da639c8.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## WordLength 属性 

[原文連結](http://api.hslcommunication.cn/html/f1a505cc-e85c-6c8b-b189-58a1ab5753f9.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Core.Device](../html/01bea659-9e5c-248c-f839-060b8da639c8.htm "HslCommunication.Core.Device")

[DeviceCommunication 类](../html/d345c980-d488-c83d-72ce-f958039275db.htm "DeviceCommunication 类")

[DeviceCommunication 属性](../html/a4a01035-8948-11b1-8a61-f2043d31db47.htm "DeviceCommunication 属性")

[ByteTransform 属性](../html/f35118c9-3691-f752-ae6f-2c5549714faa.htm "ByteTransform 属性 ")

[WordLength 属性](../html/f1a505cc-e85c-6c8b-b189-58a1ab5753f9.htm "WordLength 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DeviceCommunicationWordLength 属性 |

一个字单位的数据表示的地址长度，西门子为2，三菱，欧姆龙，modbusTcp就为1，AB PLC无效  
The address length represented by one word of data, Siemens is 2, Mitsubishi, Omron, modbusTcp is 1, AB PLC is invalid

**命名空间：**
 [HslCommunication.Core.Device](01bea659-9e5c-248c-f839-060b8da639c8.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
protected ushort WordLength { get; set; }
```

```
Protected Property WordLength As UShort
	Get
	Set
```

```
protected:
property unsigned short WordLength {
	unsigned short get ();
	void set (unsigned short value);
}
```

```
member WordLength : uint16 with get, set
```

#### 属性值

类型：UInt16

![](../icons/SectionExpanded.png)备注

对设备来说，一个地址的数据对应的字节数，或是1个字节或是2个字节，4个字节，通常是这四个选择，当设置为0时，则表示4字节的地址长度信息  
For the device, the number of bytes corresponding to the data of an address, either 1 byte or 2 bytes, 4 bytes, usually these four choices, when set to 0, it means 4 words Section address length information

![](../icons/SectionExpanded.png)参见

#### 引用

[DeviceCommunication 类](d345c980-d488-c83d-72ce-f958039275db.htm)

[HslCommunication.Core.Device 命名空间](01bea659-9e5c-248c-f839-060b8da639c8.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## DeviceCommunication 方法

[原文連結](http://api.hslcommunication.cn/html/ec3d93f5-da3b-d238-ded9-aab29ab5a132.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Core.Device](../html/01bea659-9e5c-248c-f839-060b8da639c8.htm "HslCommunication.Core.Device")

[DeviceCommunication 类](../html/d345c980-d488-c83d-72ce-f958039275db.htm "DeviceCommunication 类")

[DeviceCommunication 方法](../html/ec3d93f5-da3b-d238-ded9-aab29ab5a132.htm "DeviceCommunication 方法")

[Dispose 方法](../html/b23db6a9-5432-b84f-52d0-8c05a7d15c5f.htm "Dispose 方法 ")

[GetWordLength 方法](../html/c27821ea-1739-8eb9-a848-23a00828fec8.htm "GetWordLength 方法 ")

[Read 方法](../html/b0f83386-169b-7cf5-0c21-7226d7fda785.htm "Read 方法 ")

[ReadAsync 方法](../html/be30ee7c-174c-0bc1-d252-c53c3cd1b327.htm "ReadAsync 方法 ")

[ReadBool 方法](../html/92f161ca-fcc4-cc3a-b4b1-565066bb5d22.htm "ReadBool 方法 ")

[ReadBoolAsync 方法](../html/4dba6cf8-3f49-f03d-c0aa-baf079420e7b.htm "ReadBoolAsync 方法 ")

[ReadCustomer 方法](../html/dd38653e-aa47-424e-3feb-fe58886dd68a.htm "ReadCustomer 方法 ")

[ReadCustomerAsync 方法](../html/61a05bce-10a8-f47d-5113-26dcc5e5f2e5.htm "ReadCustomerAsync 方法 ")

[ReadDouble 方法](../html/3881ce17-f5a3-879f-372c-c641da32a34e.htm "ReadDouble 方法 ")

[ReadDoubleAsync 方法](../html/906912f8-200d-7976-944a-6191e738d405.htm "ReadDoubleAsync 方法 ")

[ReadFloat 方法](../html/7503b7f9-448a-33b2-7c4f-29d7e1bfd9f8.htm "ReadFloat 方法 ")

[ReadFloatAsync 方法](../html/558433f6-26ac-20cd-c040-a41a14a247d0.htm "ReadFloatAsync 方法 ")

[ReadInt16 方法](../html/ab92baaf-8694-565d-f2a8-b77307018056.htm "ReadInt16 方法 ")

[ReadInt16Async 方法](../html/f7764fb5-63f9-1e68-56db-a56df5d7c695.htm "ReadInt16Async 方法 ")

[ReadInt32 方法](../html/85a34c86-d647-9b5f-6cf4-b09e1d404529.htm "ReadInt32 方法 ")

[ReadInt32Async 方法](../html/8a2ad471-3f2b-3661-eeda-1133ecc9c34b.htm "ReadInt32Async 方法 ")

[ReadInt64 方法](../html/165d276d-a440-7859-2259-653f0adf9d53.htm "ReadInt64 方法 ")

[ReadInt64Async 方法](../html/7ae41fb5-16e5-d49f-d9ef-c9254f45de4b.htm "ReadInt64Async 方法 ")

[ReadString 方法](../html/d0a56ace-6194-4c01-65e3-b06080c68db4.htm "ReadString 方法 ")

[ReadStringAsync 方法](../html/8e81c96c-0715-09da-19eb-66b70236c885.htm "ReadStringAsync 方法 ")

[ReadStruct(T) 方法](../html/49853082-fbad-52e5-d756-0615cedb4b83.htm "ReadStruct(T) 方法 ")

[ReadStructAsync(T) 方法](../html/36461e61-d0ed-62cf-79ce-dcc7f71a3fa2.htm "ReadStructAsync(T) 方法 ")

[ReadUInt16 方法](../html/0fbc3a25-289c-2634-addd-437531faba42.htm "ReadUInt16 方法 ")

[ReadUInt16Async 方法](../html/4bb9a651-fc5e-a366-3b59-9e386f077341.htm "ReadUInt16Async 方法 ")

[ReadUInt32 方法](../html/eb4c948b-5f45-15cd-6997-af31f02fbebd.htm "ReadUInt32 方法 ")

[ReadUInt32Async 方法](../html/a2344c05-31e4-c854-e91d-52c7be199c88.htm "ReadUInt32Async 方法 ")

[ReadUInt64 方法](../html/1119525d-d082-efd2-62fd-eee85a6c0ebc.htm "ReadUInt64 方法 ")

[ReadUInt64Async 方法](../html/39184b30-e043-4c2d-ff3f-2da818e3c2a4.htm "ReadUInt64Async 方法 ")

[ToString 方法](../html/5667665a-da16-d227-f350-5b3f0ab9ab04.htm "ToString 方法 ")

[Wait 方法](../html/ca240694-8f9c-62e2-7edd-f249956c15d1.htm "Wait 方法 ")

[WaitAsync 方法](../html/7d5d9d45-f6bf-7f0b-4ed5-38337c5526c9.htm "WaitAsync 方法 ")

[Write 方法](../html/0a2456d7-ff9e-39ff-972d-fb1f06bf0b87.htm "Write 方法 ")

[WriteAsync 方法](../html/0d8cf9b4-de1f-abd5-77eb-b472bb7a0ab6.htm "WriteAsync 方法 ")

[WriteCustomer(T) 方法](../html/84a37f82-3f67-a0f8-fd3f-0630772565e4.htm "WriteCustomer(T) 方法 ")

[WriteCustomerAsync(T) 方法](../html/04628397-cfa1-d785-7790-4c02cb377ab6.htm "WriteCustomerAsync(T) 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DeviceCommunication 方法 |

[DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 受保护的方法 | [DecideWhetherQAMessage](acc40cda-e310-8537-6685-59b67eeb16ff.htm) | 决定当前的消息是否是用于问答机制返回的消息，默认直接返回 true, 实际的情况需要根据协议进行重写方法  To determine whether the current message is the message returned by the question answering mechanism, the default is true. In actual cases, the rewriting method needs to be performed according to the protocol (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [Dispose](b9dda6bf-d342-254e-689a-93fa0bb265d1.htm) |  |
| 受保护的方法 | [Dispose(Boolean)](0bbb6ed3-111e-069c-e9f5-8cbc3f7f7b77.htm) |  |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | [ExtraAfterReadFromCoreServer](06a49a15-a3ce-5e39-28c7-2b78851eca2e.htm) | 和服务器交互完成的时候调用的方法，可以根据读写结果进行一些额外的操作，具体的操作需要根据实际的需求来重写实现  The method called when the interaction with the server is completed can perform some additional operations based on the read and write results. The specific operations need to be rewritten according to actual needs. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法代码示例 | [ExtraOnDisconnect](de9dbba2-3b2c-1c5e-53c8-4e6f8a44e2c4.htm) | 根据实际的协议选择是否重写本方法，有些协议在断开连接之前，需要发送一些报文来关闭当前的网络通道  Select whether to rewrite this method according to the actual protocol. Some protocols need to send some packets to close the current network channel before disconnecting. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法代码示例 | [ExtraOnDisconnectAsync](210f593e-16d8-77c4-1178-58fbcfe9b27b.htm) | 根据实际的协议选择是否重写本方法，有些协议在断开连接之前，需要发送一些报文来关闭当前的网络通道  Select whether to rewrite this method according to the actual protocol. Some protocols need to send some packets to close the current network channel before disconnecting. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 受保护的方法 | [GetLogTextFromBinary](67f3b2f7-957c-fa0b-8320-6799237157db.htm) | 获取当前的报文进行日志记录的时候，是否使用二进制的格式记录，默认返回 [LogMsgFormatBinary](2ed2ff87-caec-0e44-f08b-e2b015225a58.htm)，重写可以根据session对象分别返回不同记录模式  Whether to log the current packet in binary format, the default return is [LogMsgFormatBinary](2ed2ff87-caec-0e44-f08b-e2b015225a58.htm). If you want to override it, different recording modes can be returned according to session (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [GetNewNetMessage](16636a1b-a6f0-9c0f-9d9c-c94e26c67459.htm) | 获取一个新的消息对象的方法，需要在继承类里面进行重写  The method to get a new message object needs to be overridden in the inheritance class (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法 | [GetWordLength](c27821ea-1739-8eb9-a848-23a00828fec8.htm) | 一个字单位的数据表示的地址长度，西门子为2，三菱，欧姆龙，modbusTcp就为1，AB PLC无效  The address length represented by one word of data, Siemens is 2, Mitsubishi, Omron, modbusTcp is 1, AB PLC is invalid |
| 受保护的方法代码示例 | [InitializationOnConnect](4ab6989f-a5b5-05fe-8b44-bfa45824a2a7.htm) | 根据实际的协议选择是否重写本方法，有些协议在创建连接之后，需要进行一些初始化的信号握手，才能最终建立网络通道。  Whether to rewrite this method is based on the actual protocol. Some protocols require some initial signal handshake to establish a network channel after the connection is created. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法代码示例 | [InitializationOnConnectAsync](7da2af04-d4fc-5045-19a4-4b0d1722eef3.htm) | 根据实际的协议选择是否重写本方法，有些协议在创建连接之后，需要进行一些初始化的信号握手，才能最终建立网络通道。  Whether to rewrite this method is based on the actual protocol. Some protocols require some initial signal handshake to establish a network channel after the connection is created. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [LogRevcMessage(Byte)](dab1657f-1e75-7fe1-7003-ccedfc69855d.htm) | 使用日志记录一个接收的报文信息  Logs are used to record information about a received packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [LogRevcMessage(Byte, PipeSession)](0f51ae45-575e-4228-f978-01b37d6d5df1.htm) | 使用日志记录一个接收的报文信息  Logs are used to record information about a received packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [LogSendMessage(Byte)](4f90eecc-68a2-1ef2-44ef-369f19bed3ec.htm) | 使用日志记录一个发送的报文信息  Logs are used to record information about a send packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [LogSendMessage(Byte, PipeSession)](209caadf-4f0f-40aa-04c2-ebc61f8dab13.htm) | 使用日志记录一个发送的报文信息  Logs are used to record information about a send packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法 | [PackCommandWithHeader](d0b2cd76-c306-11fc-38c7-af8ad306cce7.htm) | 对当前的命令进行打包处理，通常是携带命令头内容，标记当前的命令的长度信息，需要进行重写，否则默认不打包  The current command is packaged, usually carrying the content of the command header, marking the length of the current command, and it needs to be rewritten, otherwise it is not packaged by default (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [Read(String, UInt16)](9dd09452-c200-bc65-75ff-69aa34b366d4.htm) | 批量读取字节数组信息，需要指定地址和长度，返回原始的字节数组  Batch read byte array information, need to specify the address and length, return the original byte array |
| 公共方法代码示例 | [ReadT](95abe064-5594-5ae7-a849-98c53f88fee0.htm) | 读取支持Hsl特性的数据内容，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明，或是官网的详细文档  Read the data content of the Hsl attribute. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. |
| 公共方法 | [ReadAsync(String, UInt16)](d8c14d37-56dc-2a31-1a19-3e2b5b65197d.htm) | 异步批量读取字节数组信息，需要指定地址和长度，返回原始的字节数组  Asynchronous batch read byte array information, need to specify the address and length, return the original byte array |
| 公共方法代码示例 | [ReadAsyncT](0522e740-4bf0-2e31-c4c5-e41384321b37.htm) | 异步读取支持Hsl特性的数据内容，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Asynchronously read the data content of the Hsl attribute. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. |
| 公共方法 | [ReadBool(String)](9b2faa2e-eaa9-7ace-c49f-7a0318f1df1a.htm) | 读取单个的Boolean数据信息  Read a single Boolean data message |
| 公共方法 | [ReadBool(String, UInt16)](69d0d24f-34a7-5b41-5c1b-619edc7e6b54.htm) | 批量读取Boolean数组信息，需要指定地址和长度，返回Boolean 数组  Batch read Boolean array information, need to specify the address and length, return Boolean array |
| 公共方法 | [ReadBoolAsync(String)](cc37da1f-d673-fe67-78ee-4e465d29a2bd.htm) | 异步读取单个的Boolean数据信息  Asynchronously read a single Boolean data message |
| 公共方法 | [ReadBoolAsync(String, UInt16)](a45d383d-0095-3ea7-32f8-e75af88d60f8.htm) | 异步批量读取Boolean数组信息，需要指定地址和长度，返回Boolean 数组  Asynchronously batch read Boolean array information, need to specify the address and length, return Boolean array |
| 公共方法代码示例 | [ReadCustomerT(String)](ca069688-d8f0-7dbd-e962-86a10dce83cf.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，返回一个新的类型的实例对象。  To read a custom data type, you need to inherit from the IDataTransfer interface and return an instance object of a new type. |
| 公共方法代码示例 | [ReadCustomerT(String, T)](ca7dc1a1-c97d-689f-9819-4dae4af67206.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，传入一个实例，对这个实例进行赋值，并返回该实例的对象。  To read a custom data type, you need to inherit from the IDataTransfer interface, pass in an instance, assign a value to this instance, and return the object of the instance. |
| 公共方法代码示例 | [ReadCustomerAsyncT(String)](f3fcda2b-255d-aadd-567e-c7b3df3636d9.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，返回一个新的类型的实例对象。  To read a custom data type, you need to inherit from the IDataTransfer interface and return an instance object of a new type. |
| 公共方法代码示例 | [ReadCustomerAsyncT(String, T)](527c9fae-4743-6ebd-7b9a-5ef3a8b7c455.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，传入一个实例，对这个实例进行赋值，并返回该实例的对象。  To read a custom data type, you need to inherit from the IDataTransfer interface, pass in an instance, assign a value to this instance, and return the object of the instance. |
| 公共方法代码示例 | [ReadDouble(String)](e7e31c43-71b4-e692-6bdc-ba9533e1d83d.htm) | 读取双浮点的数据  Read double floating point data |
| 公共方法代码示例 | [ReadDouble(String, UInt16)](4eee0f05-f861-4b48-3f06-99737fb24aee.htm) | 读取双浮点数据的数组  Read double floating point data array |
| 公共方法代码示例 | [ReadDoubleAsync(String)](2cf895fa-1033-1def-f048-0e8caefaa71d.htm) | 异步读取双浮点的数据  Asynchronously read double floating point data |
| 公共方法代码示例 | [ReadDoubleAsync(String, UInt16)](2256d406-d4d5-da54-210b-ccef58a5e914.htm) | 异步读取双浮点数据的数组  Asynchronously read double floating point data array |
| 公共方法代码示例 | [ReadFloat(String)](bc2801ae-6600-e845-7c6c-522c9156eb0e.htm) | 读取单浮点数据  Read single floating point data |
| 公共方法代码示例 | [ReadFloat(String, UInt16)](6a9fb098-b13c-313b-3b27-20d3c787945d.htm) | 读取单浮点精度的数组  Read single floating point array |
| 公共方法代码示例 | [ReadFloatAsync(String)](89f743f8-4b3a-7085-1e49-a1d7adb10052.htm) | 异步读取单浮点数据  Asynchronously read single floating point data |
| 公共方法代码示例 | [ReadFloatAsync(String, UInt16)](cb6d5c5d-9b5c-2467-3e6b-5efa9d20cd72.htm) | 异步读取单浮点精度的数组  Asynchronously read single floating point array |
| 公共方法代码示例 | [ReadFromCoreServer(Byte)](8146dc0c-1edc-3ca7-428a-a3bdbd35e6c5.htm) | 将二进制的数据发送到管道中去，然后从管道里接收二进制的数据回来，并返回是否成功的结果对象。  Send binary data to the pipeline, and then receive binary data back from the pipeline, and return whether the success of the result object (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServer(IEnumerableByte)](044cfdd9-ec76-b89f-2f2b-5f439025006c.htm) | 将多个数据报文按顺序发到设备，并从设备接收返回的数据内容，然后拼接成一个Byte[]信息，需要重写[UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm)方法才能返回正确的结果。  Send multiple data packets to the device in sequence, and receive the returned data content from the device, and then splicing them into a Byte[] message, you need to rewrite [UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm) method to return the correct result. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServer(Byte, Boolean, Boolean)](b32ce8e8-a959-fcc2-6736-14b39fdb9a40.htm) | 将二进制的数据发送到管道中去，然后从管道里接收二进制的数据回来，并返回是否成功的结果对象。  Send binary data to the pipeline, and then receive binary data back from the pipeline, and return whether the success of the result object (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServer(CommunicationPipe, Byte, Boolean, Boolean)](ef36d56b-df09-f727-e8d3-316d37a13ee9.htm) | 使用指定的管道来进行数据通信，发送原始数据到管道，然后从管道接收相关的数据返回，本方法无锁 (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServerAsync(Byte)](a4de5a24-c67a-7b02-db05-33b85eaf6e1e.htm) | 将二进制的数据发送到管道中去，然后从管道里接收二进制的数据回来，并返回是否成功的结果对象。  Send binary data to the pipeline, and then receive binary data back from the pipeline, and return whether the success of the result object (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServerAsync(IEnumerableByte)](bcdb59b7-a70a-d5c0-1fc6-9ececcd39ea2.htm) | 将多个数据报文按顺序发到设备，并从设备接收返回的数据内容，然后拼接成一个Byte[]信息，需要重写[UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm)方法才能返回正确的结果。  Send multiple data packets to the device in sequence, and receive the returned data content from the device, and then splicing them into a Byte[] message, you need to rewrite [UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm) method to return the correct result. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServerAsync(Byte, Boolean, Boolean)](1c79c204-5ab5-9342-e65d-7455c20ea577.htm) | 将二进制的数据发送到管道中去，然后从管道里接收二进制的数据回来，并返回是否成功的结果对象。  Send binary data to the pipeline, and then receive binary data back from the pipeline, and return whether the success of the result object (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServerAsync(CommunicationPipe, Byte, Boolean, Boolean)](a254a19d-3c29-147b-335b-8adccdf4ee2f.htm) | 使用指定的管道来进行数据通信，发送原始数据到管道，然后从管道接收相关的数据返回，本方法无锁 (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法代码示例 | [ReadInt16(String)](57a16003-2fd5-b2bd-09bf-14366be899b8.htm) | 读取16位的有符号的整型数据  Read 16-bit signed integer data |
| 公共方法代码示例 | [ReadInt16(String, UInt16)](b2bde158-f51c-93fc-66ca-5e4cf6ea1fb7.htm) | 读取16位的有符号整型数组  Read 16-bit signed integer array |
| 公共方法代码示例 | [ReadInt16Async(String)](ac021793-890d-ceab-eae4-47ff763390c1.htm) | 异步读取16位的有符号的整型数据  Asynchronously read 16-bit signed integer data |
| 公共方法代码示例 | [ReadInt16Async(String, UInt16)](f30e02fc-8078-f1b3-c226-5991e3e23caf.htm) | 异步读取16位的有符号整型数组  Asynchronously read 16-bit signed integer array |
| 公共方法代码示例 | [ReadInt32(String)](bfbea035-9a83-40e3-084c-28122041b37a.htm) | 读取32位的有符号整型  Read 32-bit signed integer |
| 公共方法代码示例 | [ReadInt32(String, UInt16)](095240a3-bc0f-8647-4867-88e40cc089e3.htm) | 读取32位有符号整型数组  Read 32-bit signed integer array |
| 公共方法代码示例 | [ReadInt32Async(String)](96726ea0-82aa-1229-21f6-43cd62abca43.htm) | 异步读取32位的有符号整型  Asynchronously read 32-bit signed integer |
| 公共方法代码示例 | [ReadInt32Async(String, UInt16)](d6356095-a289-18eb-136c-b7e6ae099616.htm) | 异步读取32位有符号整型数组  Asynchronously read 32-bit signed integer array |
| 公共方法代码示例 | [ReadInt64(String)](0d0c2284-2da1-6f9a-b353-537561cd5e95.htm) | 读取64位的有符号整型  Read 64-bit signed integer |
| 公共方法代码示例 | [ReadInt64(String, UInt16)](35a55019-efb1-bf2d-be1b-01c8eed99243.htm) | 读取64位的有符号整型数组  Read 64-bit signed integer array |
| 公共方法代码示例 | [ReadInt64Async(String)](26d46c8a-3775-3abd-d3ec-5cd1758a317e.htm) | 异步读取64位的有符号整型  Asynchronously read 64-bit signed integer |
| 公共方法代码示例 | [ReadInt64Async(String, UInt16)](8dc97c4a-8904-9a34-9ddc-76c1a9477ee5.htm) | 异步读取64位的有符号整型数组  Asynchronously read 64-bit signed integer array |
| 公共方法代码示例 | [ReadString(String, UInt16)](07812d76-415a-6400-c8a2-1c21da51e7a3.htm) | 读取字符串数据，默认为最常见的ASCII编码  Read string data, default is the most common ASCII encoding |
| 公共方法代码示例 | [ReadString(String, UInt16, Encoding)](b34c71a6-58f3-6bf4-0b93-471be4f2ed31.htm) | 使用指定的编码，读取字符串数据  Reads string data using the specified encoding |
| 公共方法代码示例 | [ReadStringAsync(String, UInt16)](1835da81-b811-d151-0b53-907d8d87d0a8.htm) | 异步读取字符串数据，默认为最常见的ASCII编码  Asynchronously read string data, default is the most common ASCII encoding |
| 公共方法代码示例 | [ReadStringAsync(String, UInt16, Encoding)](091dfc67-3d25-f33e-3a55-a0228863815e.htm) | 异步使用指定的编码，读取字符串数据  Asynchronously reads string data using the specified encoding |
| 公共方法代码示例 | [ReadStructT](49853082-fbad-52e5-d756-0615cedb4b83.htm) | 读取结构体类型的数据，根据结构体自身的定义，读取原始字节数组，然后解析出实际的结构体数据，结构体需要实现[HslStructAttribute](87e427a8-0c7b-43a8-c421-c27264eda90f.htm)特性 |
| 公共方法代码示例 | [ReadStructAsyncT](36461e61-d0ed-62cf-79ce-dcc7f71a3fa2.htm) | 读取结构体类型的数据，根据结构体自身的定义，读取原始字节数组，然后解析出实际的结构体数据，结构体需要实现[HslStructAttribute](87e427a8-0c7b-43a8-c421-c27264eda90f.htm)特性 |
| 公共方法代码示例 | [ReadUInt16(String)](736379d1-9560-5052-658a-13c126d56d78.htm) | 读取16位的无符号整型  Read 16-bit unsigned integer |
| 公共方法代码示例 | [ReadUInt16(String, UInt16)](e284c4bd-d3f8-ff30-eb42-0aa035cfb4f6.htm) | 读取16位的无符号整型数组  Read 16-bit unsigned integer array |
| 公共方法代码示例 | [ReadUInt16Async(String)](361acc93-43b5-a9a5-6fa7-fa8728b7b96c.htm) | 异步读取16位的无符号整型  Asynchronously read 16-bit unsigned integer |
| 公共方法代码示例 | [ReadUInt16Async(String, UInt16)](73dd24b8-b501-df53-1de6-ccb2708f656d.htm) | 异步读取16位的无符号整型数组  Asynchronously read 16-bit unsigned integer array |
| 公共方法代码示例 | [ReadUInt32(String)](ca73c958-39c7-9b34-aa44-3db2d9b0cad3.htm) | 读取32位的无符号整型  Read 32-bit unsigned integer |
| 公共方法代码示例 | [ReadUInt32(String, UInt16)](878b4709-0278-fcdf-3b65-afbf0f496c1e.htm) | 读取32位的无符号整型数组  Read 32-bit unsigned integer array |
| 公共方法代码示例 | [ReadUInt32Async(String)](2343953f-6350-331e-6ba7-18a4449f0575.htm) | 异步读取32位的无符号整型  Asynchronously read 32-bit unsigned integer |
| 公共方法代码示例 | [ReadUInt32Async(String, UInt16)](e0eba6c1-56a2-5190-ba03-66116bcafe3d.htm) | 异步读取32位的无符号整型数组  Asynchronously read 32-bit unsigned integer array |
| 公共方法代码示例 | [ReadUInt64(String)](947487dd-d4a3-72d8-d9fb-2224aad70e12.htm) | 读取64位的无符号整型  Read 64-bit unsigned integer |
| 公共方法代码示例 | [ReadUInt64(String, UInt16)](858766a9-9d62-87a8-f727-a00838f3d2b6.htm) | 读取64位的无符号整型的数组  Read 64-bit unsigned integer array |
| 公共方法代码示例 | [ReadUInt64Async(String)](3411ff3a-772f-dbe1-8931-75e00a731579.htm) | 异步读取64位的无符号整型  Asynchronously read 64-bit unsigned integer |
| 公共方法代码示例 | [ReadUInt64Async(String, UInt16)](c2095a6e-ff89-dd98-8a21-fe72f82dbbc6.htm) | 异步读取64位的无符号整型的数组  Asynchronously read 64-bit unsigned integer array |
| 公共方法 | [SetDtuPipe](f765d1d7-234a-b229-9457-97e25e024dc9.htm) | 将当前的通信对象设置DTU模式，允许传入现成的管道，并返回初始化结果，如果该设备重写了握手报文，就是返回握手结果  Set the current communication object to DTU mode, allow the existing pipe to be passed in, and return the initialization result, if the device rewrites the handshake packet, the handshake result is returned (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [SetDtuPipeAsync](9d1fd12d-b804-3828-a627-a2c0eef1d23f.htm) | 将当前的通信对象设置DTU模式，允许传入现成的管道，并返回初始化结果，如果该设备重写了握手报文，就是返回握手结果  Set the current communication object to DTU mode, allow the existing pipe to be passed in, and return the initialization result, if the device rewrites the handshake packet, the handshake result is returned (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ToString](5667665a-da16-d227-f350-5b3f0ab9ab04.htm) | (重写 ObjectToString.) |
| 公共方法 | [UnpackResponseContent](208017a7-8da2-33f8-3273-d7e7e850100f.htm) | 根据对方返回的报文命令，对命令进行基本的拆包，例如各种Modbus协议拆包为统一的核心报文，还支持对报文的验证  According to the message command returned by the other party, the command is basically unpacked, for example, various Modbus protocols are unpacked into a unified core message, and the verification of the message is also supported (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [Wait(String, Boolean, Int32, Int32)](0e4db26e-1177-2b18-9492-ba05afb6fcc4.htm) | 等待指定地址的Boolean值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the Boolean value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. |
| 公共方法 | [Wait(String, Int16, Int32, Int32)](e06fe098-8c8b-83f9-73ca-2a7cce19d43d.htm) | 等待指定地址的Int16值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the Int16 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. |
| 公共方法 | [Wait(String, Int32, Int32, Int32)](1bce0aae-27e9-f1d6-e8f3-25ae6b47554e.htm) | 等待指定地址的Int32值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the Int32 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. |
| 公共方法 | [Wait(String, Int64, Int32, Int32)](073354e9-2508-2158-1b01-881c3a34844b.htm) | 等待指定地址的Int64值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the Int64 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. |
| 公共方法 | [Wait(String, UInt16, Int32, Int32)](d97f8d26-1c1a-df69-93f0-f60209597ced.htm) | 等待指定地址的UInt16值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the UInt16 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. |
| 公共方法 | [Wait(String, UInt32, Int32, Int32)](1e61c395-b07f-e66c-5f9c-2708093150d9.htm) | 等待指定地址的UInt32值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the UInt32 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. |
| 公共方法 | [Wait(String, UInt64, Int32, Int32)](d659729b-5766-c497-d418-856c34ec2b24.htm) | 等待指定地址的UInt64值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the UInt64 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. |
| 公共方法 | [WaitAsync(String, Boolean, Int32, Int32)](65399e87-9d38-f819-ccd3-005f68c2c3b3.htm) | 等待指定地址的Boolean值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the Boolean value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. |
| 公共方法 | [WaitAsync(String, Int16, Int32, Int32)](01277e39-afb1-86a4-3aa3-560160e66fb9.htm) | 等待指定地址的Int16值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the Int16 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. |
| 公共方法 | [WaitAsync(String, Int32, Int32, Int32)](4a6d380f-6b2a-1270-07a7-39b88a2712fd.htm) | 等待指定地址的Int32值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the Int32 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. |
| 公共方法 | [WaitAsync(String, Int64, Int32, Int32)](d7eede9b-4d44-479b-fb88-3ff2e1bb7568.htm) | 等待指定地址的Int64值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the Int64 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. |
| 公共方法 | [WaitAsync(String, UInt16, Int32, Int32)](0a8395c7-dc7f-4d5f-4316-98e18fe70f4c.htm) | 等待指定地址的UInt16值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the UInt16 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. |
| 公共方法 | [WaitAsync(String, UInt32, Int32, Int32)](19e9074e-99a2-99fd-91ef-a5314a115397.htm) | 等待指定地址的UInt32值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the UInt32 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. |
| 公共方法 | [WaitAsync(String, UInt64, Int32, Int32)](5109dac0-0ba2-cbb9-e383-1a869c6dc22e.htm) | 等待指定地址的UInt64值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the UInt64 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. |
| 公共方法 | [Write(String, Boolean)](06d0b9af-7f28-0b9c-098f-4da28e69504a.htm) | 写入单个的Boolean数据，返回是否成功  Write a single Boolean data, and return whether the write was successful |
| 公共方法 | [Write(String, Boolean)](5b046c82-3d8e-550b-fdf3-41ca5558169b.htm) | 批量写入Boolean数组数据，返回是否成功  Batch write Boolean array data, return whether the write was successful |
| 公共方法 | [Write(String, Byte)](b5536dd6-ccde-6883-ac27-bff9e4b37806.htm) | 写入原始的byte数组数据到指定的地址，返回是否写入成功  Write the original byte array data to the specified address, and return whether the write was successful |
| 公共方法代码示例 | [Write(String, Double)](64c44576-7e95-cbcb-a8c9-e6b7ff7e3c50.htm) | 写入double数据，返回是否成功  Write double data, return whether the write was successful |
| 公共方法代码示例 | [Write(String, Double)](df56e6b1-0688-ece7-3bf1-9db32fec9f6e.htm) | 写入double数组，返回是否成功  Write double array, return whether the write was successful |
| 公共方法代码示例 | [Write(String, Int16)](0e23506e-b604-e67c-b5eb-06271b6c6335.htm) | 写入short数据，返回是否成功  Write short data, returns whether success |
| 公共方法代码示例 | [Write(String, Int16)](c7cfd67d-ca31-3087-c684-ee4c10780ab1.htm) | 写入short数组，返回是否成功  Write short array, return whether the write was successful |
| 公共方法代码示例 | [Write(String, Int32)](ba80d1a5-c8d8-7e0b-5c6b-9435cf20cadb.htm) | 写入int数据，返回是否成功  Write int data, return whether the write was successful |
| 公共方法代码示例 | [Write(String, Int32)](57dac21c-8b28-de7c-c213-77cda465cccb.htm) | 写入int[]数组，返回是否成功  Write int array, return whether the write was successful |
| 公共方法代码示例 | [Write(String, Int64)](ba246b33-5a12-a22c-6d40-0969a64b508b.htm) | 写入long数据，返回是否成功  Write long data, return whether the write was successful |
| 公共方法代码示例 | [Write(String, Int64)](5b910f9c-6e6b-bb9d-6842-1b496640c1f4.htm) | 写入long数组，返回是否成功  Write long array, return whether the write was successful |
| 公共方法代码示例 | [Write(String, Single)](a2265fa7-5b56-a5aa-9d21-a0dc5458cbf9.htm) | 写入float数据，返回是否成功  Write float data, return whether the write was successful |
| 公共方法代码示例 | [Write(String, Single)](bb76227e-3e95-8ca2-5dcf-c087d90b8a01.htm) | 写入float数组，返回是否成功  Write float array, return whether the write was successful |
| 公共方法代码示例 | [Write(String, String)](4e0ec363-ba64-7ca7-b8d3-2e9a8087bd55.htm) | 写入字符串信息，编码为ASCII  Write string information, encoded as ASCII |
| 公共方法代码示例 | [Write(String, UInt16)](e5bc9466-65b4-fc64-e0d0-1867b1113130.htm) | 写入ushort数据，返回是否成功  Write ushort data, return whether the write was successful |
| 公共方法代码示例 | [Write(String, UInt16)](65544962-82ac-bb47-d72f-5dd0c3d4ac6c.htm) | 写入ushort数组，返回是否成功  Write ushort array, return whether the write was successful |
| 公共方法代码示例 | [Write(String, UInt32)](c99d8246-f95e-4f79-738b-16e622072eb4.htm) | 写入uint数据，返回是否成功  Write uint data, return whether the write was successful |
| 公共方法代码示例 | [Write(String, UInt32)](b80b3099-b34d-fde8-6a65-837da3bdc5ab.htm) | 写入uint[]数组，返回是否成功  Write uint array, return whether the write was successful |
| 公共方法代码示例 | [Write(String, UInt64)](4924086e-9b36-d446-4ae4-6682134a6ad1.htm) | 写入ulong数据，返回是否成功  Write ulong data, return whether the write was successful |
| 公共方法代码示例 | [Write(String, UInt64)](b1131bdd-8b51-ec6d-dd0e-fec1f6cb51b6.htm) | 写入ulong数组，返回是否成功  Write ulong array, return whether the write was successful |
| 公共方法代码示例 | [Write(String, String, Int32)](02fe2ca9-49a3-ede4-d1af-8c6da1f644da.htm) | 写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为ASCII  Write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, it is filled with 0 and the encoding is ASCII. |
| 公共方法代码示例 | [Write(String, String, Encoding)](7b602661-902d-2b70-1fdb-ba87ecd307bf.htm) | 写入字符串信息，需要指定的编码信息  Write string information, need to specify the encoding information |
| 公共方法代码示例 | [Write(String, String, Int32, Encoding)](e8cb30d9-4c5e-6c68-b9af-1b1b6c2dc714.htm) | 写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为指定的编码信息  Write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, then the operation is complemented with 0 , you should specified the encoding information |
| 公共方法代码示例 | [WriteT(T)](9f48853f-fbeb-c37d-25a2-1c76d1dee407.htm) | 写入支持Hsl特性的数据，返回是否写入成功，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Write data that supports the Hsl attribute, and return whether the write was successful. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. |
| 公共方法 | [WriteAsync(String, Boolean)](10fe669c-805e-1e5a-f016-f319ad363e90.htm) | 异步批量写入Boolean数组数据，返回是否成功  Asynchronously batch write Boolean array data, return success |
| 公共方法 | [WriteAsync(String, Boolean)](fe872dae-7ef9-12c1-d883-e92ba8676416.htm) | 异步批量写入Boolean数组数据，返回是否成功  Asynchronously batch write Boolean array data, return success |
| 公共方法代码示例 | [WriteAsync(String, Byte)](8f1cec30-8b86-6749-c8f3-7d39168d59df.htm) | 异步写入原始的byte数组数据到指定的地址，返回是否写入成功  Asynchronously writes the original byte array data to the specified address, and returns whether the write was successful |
| 公共方法代码示例 | [WriteAsync(String, Double)](068adaa4-d1de-dffc-8376-6f483e2c1d96.htm) | 异步写入double数据，返回是否成功  Asynchronously write double data, return whether the write was successful |
| 公共方法代码示例 | [WriteAsync(String, Double)](36199b83-c4ea-b02f-3675-fb6ac48d0c59.htm) | 异步写入double数组，返回是否成功  Asynchronously write double array, return whether the write was successful |
| 公共方法代码示例 | [WriteAsync(String, Int16)](ccd95ccc-10a3-af67-62e8-c2312b4871b0.htm) | 异步写入short数据，返回是否成功  Asynchronously write short data, returns whether success |
| 公共方法代码示例 | [WriteAsync(String, Int16)](2625e868-7096-e526-ce25-b29afc639b41.htm) | 异步写入short数组，返回是否成功  Asynchronously write short array, return whether the write was successful |
| 公共方法代码示例 | [WriteAsync(String, Int32)](df91d629-2080-328d-67cd-fbefb9a45404.htm) | 异步写入int数据，返回是否成功  Asynchronously write int data, return whether the write was successful |
| 公共方法代码示例 | [WriteAsync(String, Int32)](99962a08-9624-cea6-6009-5e6f207236b7.htm) | 异步写入int[]数组，返回是否成功  Asynchronously write int array, return whether the write was successful |
| 公共方法代码示例 | [WriteAsync(String, Int64)](ec6c1e24-a0d8-f1f5-56af-109b0c546fc8.htm) | 异步写入long数据，返回是否成功  Asynchronously write long data, return whether the write was successful |
| 公共方法代码示例 | [WriteAsync(String, Int64)](4b2ab94c-2891-e0a9-f10b-b67c63bd1cf9.htm) | 异步写入long数组，返回是否成功  Asynchronously write long array, return whether the write was successful |
| 公共方法代码示例 | [WriteAsync(String, Single)](910189df-13f9-b7d0-a626-b5b25ce5d32a.htm) | 异步写入float数据，返回是否成功  Asynchronously write float data, return whether the write was successful |
| 公共方法代码示例 | [WriteAsync(String, Single)](dbb2359a-fb54-cf93-4413-20d5a4ef9d36.htm) | 异步写入float数组，返回是否成功  Asynchronously write float array, return whether the write was successful |
| 公共方法代码示例 | [WriteAsync(String, String)](fe839c9a-46f3-e317-0155-b10a8a83a9b1.htm) | 异步写入字符串信息，编码为ASCII  Asynchronously write string information, encoded as ASCII |
| 公共方法代码示例 | [WriteAsync(String, UInt16)](a04f3aa2-a51e-26cd-974e-bf70a2db359d.htm) | 异步写入ushort数据，返回是否成功  Asynchronously write ushort data, return whether the write was successful |
| 公共方法代码示例 | [WriteAsync(String, UInt16)](34a1a321-972e-06a8-76fb-794dd1229005.htm) | 异步写入ushort数组，返回是否成功  Asynchronously write ushort array, return whether the write was successful |
| 公共方法代码示例 | [WriteAsync(String, UInt32)](abad215c-bc99-ac90-b9f3-7adb3a7cca8f.htm) | 异步写入uint数据，返回是否成功  Asynchronously write uint data, return whether the write was successful |
| 公共方法代码示例 | [WriteAsync(String, UInt32)](5b655bf3-5ba7-ad31-6897-89c59051ac60.htm) | 异步写入uint[]数组，返回是否成功  Asynchronously write uint array, return whether the write was successful |
| 公共方法代码示例 | [WriteAsync(String, UInt64)](e396074d-6cf1-b45d-5e8e-01026f6e391b.htm) | 异步写入ulong数据，返回是否成功  Asynchronously write ulong data, return whether the write was successful |
| 公共方法代码示例 | [WriteAsync(String, UInt64)](ae3df8ec-29c0-42ff-3d6a-68f496997b3d.htm) | 异步写入ulong数组，返回是否成功  Asynchronously write ulong array, return whether the write was successful |
| 公共方法 | [WriteAsync(String, String, Int32)](b739b983-b03c-8ecf-2713-75f19233ae2f.htm) | 异步写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为ASCII  Asynchronously write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, it is filled with 0 and the encoding is ASCII. |
| 公共方法代码示例 | [WriteAsync(String, String, Encoding)](876e9800-0406-6407-bdf8-8ab1cf28df38.htm) | 异步写入字符串信息，需要指定的编码信息  Asynchronously write string information, need to specify the encoding information |
| 公共方法 | [WriteAsync(String, String, Int32, Encoding)](ae9fb2a3-fddb-48eb-e9ad-6cc0ebbac1f6.htm) | 异步写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为指定的编码信息  Asynchronously write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, then the operation is complemented with 0 , you should specified the encoding information |
| 公共方法代码示例 | [WriteAsyncT(T)](25d96e77-d1b2-918c-33f9-e5e4cc57c9cd.htm) | 异步写入支持Hsl特性的数据，返回是否写入成功，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Asynchronously write data that supports the Hsl attribute, and return whether the write was successful. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. |
| 公共方法代码示例 | [WriteCustomerT](84a37f82-3f67-a0f8-fd3f-0630772565e4.htm) | 写入自定义类型的数据，该类型必须继承自IDataTransfer接口  Write data of a custom type, which must inherit from the IDataTransfer interface |
| 公共方法代码示例 | [WriteCustomerAsyncT](04628397-cfa1-d785-7790-4c02cb377ab6.htm) | 写入自定义类型的数据，该类型必须继承自IDataTransfer接口  Write data of a custom type, which must inherit from the IDataTransfer interface |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[DeviceCommunication 类](d345c980-d488-c83d-72ce-f958039275db.htm)

[HslCommunication.Core.Device 命名空间](01bea659-9e5c-248c-f839-060b8da639c8.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Dispose 方法 

[原文連結](http://api.hslcommunication.cn/html/b23db6a9-5432-b84f-52d0-8c05a7d15c5f.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Core.Device](../html/01bea659-9e5c-248c-f839-060b8da639c8.htm "HslCommunication.Core.Device")

[DeviceCommunication 类](../html/d345c980-d488-c83d-72ce-f958039275db.htm "DeviceCommunication 类")

[DeviceCommunication 方法](../html/ec3d93f5-da3b-d238-ded9-aab29ab5a132.htm "DeviceCommunication 方法")

[Dispose 方法](../html/b23db6a9-5432-b84f-52d0-8c05a7d15c5f.htm "Dispose 方法 ")

[Dispose 方法](../html/b9dda6bf-d342-254e-689a-93fa0bb265d1.htm "Dispose 方法 ")

[Dispose 方法 (Boolean)](../html/0bbb6ed3-111e-069c-e9f5-8cbc3f7f7b77.htm "Dispose 方法 (Boolean)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DeviceCommunicationDispose 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [Dispose](b9dda6bf-d342-254e-689a-93fa0bb265d1.htm) | 释放被 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm) 使用的所有资源 |
| 受保护的方法 | [Dispose(Boolean)](0bbb6ed3-111e-069c-e9f5-8cbc3f7f7b77.htm) | 释放被 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm) 使用的非托管资源，并且是否托管资源（可选） |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[DeviceCommunication 类](d345c980-d488-c83d-72ce-f958039275db.htm)

[HslCommunication.Core.Device 命名空间](01bea659-9e5c-248c-f839-060b8da639c8.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Dispose 方法 

[原文連結](http://api.hslcommunication.cn/html/b9dda6bf-d342-254e-689a-93fa0bb265d1.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Core.Device](../html/01bea659-9e5c-248c-f839-060b8da639c8.htm "HslCommunication.Core.Device")

[DeviceCommunication 类](../html/d345c980-d488-c83d-72ce-f958039275db.htm "DeviceCommunication 类")

[DeviceCommunication 方法](../html/ec3d93f5-da3b-d238-ded9-aab29ab5a132.htm "DeviceCommunication 方法")

[Dispose 方法](../html/b23db6a9-5432-b84f-52d0-8c05a7d15c5f.htm "Dispose 方法 ")

[Dispose 方法](../html/b9dda6bf-d342-254e-689a-93fa0bb265d1.htm "Dispose 方法 ")

[Dispose 方法 (Boolean)](../html/0bbb6ed3-111e-069c-e9f5-8cbc3f7f7b77.htm "Dispose 方法 (Boolean)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DeviceCommunicationDispose 方法 |

释放被 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm) 使用的所有资源

**命名空间：**
 [HslCommunication.Core.Device](01bea659-9e5c-248c-f839-060b8da639c8.htm)  
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

[DeviceCommunication 类](d345c980-d488-c83d-72ce-f958039275db.htm)

[Dispose 重载](b23db6a9-5432-b84f-52d0-8c05a7d15c5f.htm)

[HslCommunication.Core.Device 命名空间](01bea659-9e5c-248c-f839-060b8da639c8.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Dispose 方法 (Boolean)

[原文連結](http://api.hslcommunication.cn/html/0bbb6ed3-111e-069c-e9f5-8cbc3f7f7b77.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Core.Device](../html/01bea659-9e5c-248c-f839-060b8da639c8.htm "HslCommunication.Core.Device")

[DeviceCommunication 类](../html/d345c980-d488-c83d-72ce-f958039275db.htm "DeviceCommunication 类")

[DeviceCommunication 方法](../html/ec3d93f5-da3b-d238-ded9-aab29ab5a132.htm "DeviceCommunication 方法")

[Dispose 方法](../html/b23db6a9-5432-b84f-52d0-8c05a7d15c5f.htm "Dispose 方法 ")

[Dispose 方法](../html/b9dda6bf-d342-254e-689a-93fa0bb265d1.htm "Dispose 方法 ")

[Dispose 方法 (Boolean)](../html/0bbb6ed3-111e-069c-e9f5-8cbc3f7f7b77.htm "Dispose 方法 (Boolean)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DeviceCommunicationDispose 方法 (Boolean) |

释放被 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm) 使用的非托管资源，并且是否托管资源（可选）

**命名空间：**
 [HslCommunication.Core.Device](01bea659-9e5c-248c-f839-060b8da639c8.htm)  
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
    为 true 则同时释放托管资源和非托管资源；为 false 则只释放非托管资源

![](../icons/SectionExpanded.png)参见

#### 引用

[DeviceCommunication 类](d345c980-d488-c83d-72ce-f958039275db.htm)

[Dispose 重载](b23db6a9-5432-b84f-52d0-8c05a7d15c5f.htm)

[HslCommunication.Core.Device 命名空间](01bea659-9e5c-248c-f839-060b8da639c8.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GetWordLength 方法 

[原文連結](http://api.hslcommunication.cn/html/c27821ea-1739-8eb9-a848-23a00828fec8.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Core.Device](../html/01bea659-9e5c-248c-f839-060b8da639c8.htm "HslCommunication.Core.Device")

[DeviceCommunication 类](../html/d345c980-d488-c83d-72ce-f958039275db.htm "DeviceCommunication 类")

[DeviceCommunication 方法](../html/ec3d93f5-da3b-d238-ded9-aab29ab5a132.htm "DeviceCommunication 方法")

[Dispose 方法](../html/b23db6a9-5432-b84f-52d0-8c05a7d15c5f.htm "Dispose 方法 ")

[GetWordLength 方法](../html/c27821ea-1739-8eb9-a848-23a00828fec8.htm "GetWordLength 方法 ")

[Read 方法](../html/b0f83386-169b-7cf5-0c21-7226d7fda785.htm "Read 方法 ")

[ReadAsync 方法](../html/be30ee7c-174c-0bc1-d252-c53c3cd1b327.htm "ReadAsync 方法 ")

[ReadBool 方法](../html/92f161ca-fcc4-cc3a-b4b1-565066bb5d22.htm "ReadBool 方法 ")

[ReadBoolAsync 方法](../html/4dba6cf8-3f49-f03d-c0aa-baf079420e7b.htm "ReadBoolAsync 方法 ")

[ReadCustomer 方法](../html/dd38653e-aa47-424e-3feb-fe58886dd68a.htm "ReadCustomer 方法 ")

[ReadCustomerAsync 方法](../html/61a05bce-10a8-f47d-5113-26dcc5e5f2e5.htm "ReadCustomerAsync 方法 ")

[ReadDouble 方法](../html/3881ce17-f5a3-879f-372c-c641da32a34e.htm "ReadDouble 方法 ")

[ReadDoubleAsync 方法](../html/906912f8-200d-7976-944a-6191e738d405.htm "ReadDoubleAsync 方法 ")

[ReadFloat 方法](../html/7503b7f9-448a-33b2-7c4f-29d7e1bfd9f8.htm "ReadFloat 方法 ")

[ReadFloatAsync 方法](../html/558433f6-26ac-20cd-c040-a41a14a247d0.htm "ReadFloatAsync 方法 ")

[ReadInt16 方法](../html/ab92baaf-8694-565d-f2a8-b77307018056.htm "ReadInt16 方法 ")

[ReadInt16Async 方法](../html/f7764fb5-63f9-1e68-56db-a56df5d7c695.htm "ReadInt16Async 方法 ")

[ReadInt32 方法](../html/85a34c86-d647-9b5f-6cf4-b09e1d404529.htm "ReadInt32 方法 ")

[ReadInt32Async 方法](../html/8a2ad471-3f2b-3661-eeda-1133ecc9c34b.htm "ReadInt32Async 方法 ")

[ReadInt64 方法](../html/165d276d-a440-7859-2259-653f0adf9d53.htm "ReadInt64 方法 ")

[ReadInt64Async 方法](../html/7ae41fb5-16e5-d49f-d9ef-c9254f45de4b.htm "ReadInt64Async 方法 ")

[ReadString 方法](../html/d0a56ace-6194-4c01-65e3-b06080c68db4.htm "ReadString 方法 ")

[ReadStringAsync 方法](../html/8e81c96c-0715-09da-19eb-66b70236c885.htm "ReadStringAsync 方法 ")

[ReadStruct(T) 方法](../html/49853082-fbad-52e5-d756-0615cedb4b83.htm "ReadStruct(T) 方法 ")

[ReadStructAsync(T) 方法](../html/36461e61-d0ed-62cf-79ce-dcc7f71a3fa2.htm "ReadStructAsync(T) 方法 ")

[ReadUInt16 方法](../html/0fbc3a25-289c-2634-addd-437531faba42.htm "ReadUInt16 方法 ")

[ReadUInt16Async 方法](../html/4bb9a651-fc5e-a366-3b59-9e386f077341.htm "ReadUInt16Async 方法 ")

[ReadUInt32 方法](../html/eb4c948b-5f45-15cd-6997-af31f02fbebd.htm "ReadUInt32 方法 ")

[ReadUInt32Async 方法](../html/a2344c05-31e4-c854-e91d-52c7be199c88.htm "ReadUInt32Async 方法 ")

[ReadUInt64 方法](../html/1119525d-d082-efd2-62fd-eee85a6c0ebc.htm "ReadUInt64 方法 ")

[ReadUInt64Async 方法](../html/39184b30-e043-4c2d-ff3f-2da818e3c2a4.htm "ReadUInt64Async 方法 ")

[ToString 方法](../html/5667665a-da16-d227-f350-5b3f0ab9ab04.htm "ToString 方法 ")

[Wait 方法](../html/ca240694-8f9c-62e2-7edd-f249956c15d1.htm "Wait 方法 ")

[WaitAsync 方法](../html/7d5d9d45-f6bf-7f0b-4ed5-38337c5526c9.htm "WaitAsync 方法 ")

[Write 方法](../html/0a2456d7-ff9e-39ff-972d-fb1f06bf0b87.htm "Write 方法 ")

[WriteAsync 方法](../html/0d8cf9b4-de1f-abd5-77eb-b472bb7a0ab6.htm "WriteAsync 方法 ")

[WriteCustomer(T) 方法](../html/84a37f82-3f67-a0f8-fd3f-0630772565e4.htm "WriteCustomer(T) 方法 ")

[WriteCustomerAsync(T) 方法](../html/04628397-cfa1-d785-7790-4c02cb377ab6.htm "WriteCustomerAsync(T) 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DeviceCommunicationGetWordLength 方法 |

一个字单位的数据表示的地址长度，西门子为2，三菱，欧姆龙，modbusTcp就为1，AB PLC无效  
The address length represented by one word of data, Siemens is 2, Mitsubishi, Omron, modbusTcp is 1, AB PLC is invalid

**命名空间：**
 [HslCommunication.Core.Device](01bea659-9e5c-248c-f839-060b8da639c8.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
protected virtual ushort GetWordLength(
	string address,
	int length,
	int dataTypeLength
)
```

```
Protected Overridable Function GetWordLength ( 
	address As String,
	length As Integer,
	dataTypeLength As Integer
) As UShort
```

```
protected:
virtual unsigned short GetWordLength(
	String^ address, 
	int length, 
	int dataTypeLength
)
```

```
abstract GetWordLength : 
        address : string * 
        length : int * 
        dataTypeLength : int -> uint16 
override GetWordLength : 
        address : string * 
        length : int * 
        dataTypeLength : int -> uint16
```

#### 参数

address
:   类型：SystemString  
    读取的设备的地址信息

length
:   类型：SystemInt32  
    读取的数据长度信息

dataTypeLength
:   类型：SystemInt32  
    数据类型的字节长度信息，比如short, 就是2，int,float就是4

#### 返回值

类型：UInt16  

[缺少 "M:HslCommunication.Core.Device.DeviceCommunication.GetWordLength(System.String,System.Int32,System.Int32)" 的 <returns> 文档]

![](../icons/SectionExpanded.png)备注

对设备来说，一个地址的数据对应的字节数，或是1个字节或是2个字节，通常是这两个选择。  
当前也可以重写来根据不同的地址动态控制不同的地址长度，比如有的地址是一个地址一个字节的，有的地址是一个地址两个字节的

![](../icons/SectionExpanded.png)参见

#### 引用

[DeviceCommunication 类](d345c980-d488-c83d-72ce-f958039275db.htm)

[HslCommunication.Core.Device 命名空间](01bea659-9e5c-248c-f839-060b8da639c8.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Read 方法 

[原文連結](http://api.hslcommunication.cn/html/b0f83386-169b-7cf5-0c21-7226d7fda785.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Core.Device](../html/01bea659-9e5c-248c-f839-060b8da639c8.htm "HslCommunication.Core.Device")

[DeviceCommunication 类](../html/d345c980-d488-c83d-72ce-f958039275db.htm "DeviceCommunication 类")

[DeviceCommunication 方法](../html/ec3d93f5-da3b-d238-ded9-aab29ab5a132.htm "DeviceCommunication 方法")

[Read 方法](../html/b0f83386-169b-7cf5-0c21-7226d7fda785.htm "Read 方法 ")

[Read(T) 方法](../html/95abe064-5594-5ae7-a849-98c53f88fee0.htm "Read(T) 方法 ")

[Read 方法 (String, UInt16)](../html/9dd09452-c200-bc65-75ff-69aa34b366d4.htm "Read 方法 (String, UInt16)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DeviceCommunicationRead 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法代码示例 | [ReadT](95abe064-5594-5ae7-a849-98c53f88fee0.htm) | 读取支持Hsl特性的数据内容，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明，或是官网的详细文档  Read the data content of the Hsl attribute. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. |
| 公共方法 | [Read(String, UInt16)](9dd09452-c200-bc65-75ff-69aa34b366d4.htm) | 批量读取字节数组信息，需要指定地址和长度，返回原始的字节数组  Batch read byte array information, need to specify the address and length, return the original byte array |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[DeviceCommunication 类](d345c980-d488-c83d-72ce-f958039275db.htm)

[HslCommunication.Core.Device 命名空间](01bea659-9e5c-248c-f839-060b8da639c8.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Read(T) 方法 

[原文連結](http://api.hslcommunication.cn/html/95abe064-5594-5ae7-a849-98c53f88fee0.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Core.Device](../html/01bea659-9e5c-248c-f839-060b8da639c8.htm "HslCommunication.Core.Device")

[DeviceCommunication 类](../html/d345c980-d488-c83d-72ce-f958039275db.htm "DeviceCommunication 类")

[DeviceCommunication 方法](../html/ec3d93f5-da3b-d238-ded9-aab29ab5a132.htm "DeviceCommunication 方法")

[Read 方法](../html/b0f83386-169b-7cf5-0c21-7226d7fda785.htm "Read 方法 ")

[Read(T) 方法](../html/95abe064-5594-5ae7-a849-98c53f88fee0.htm "Read(T) 方法 ")

[Read 方法 (String, UInt16)](../html/9dd09452-c200-bc65-75ff-69aa34b366d4.htm "Read 方法 (String, UInt16)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DeviceCommunicationReadT 方法 |

读取支持Hsl特性的数据内容，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明，或是官网的详细文档  
Read the data content of the Hsl attribute. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details.

**命名空间：**
 [HslCommunication.Core.Device](01bea659-9e5c-248c-f839-060b8da639c8.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public virtual OperateResult<T> Read<T>()
where T : class, new()
```

```
Public Overridable Function Read(Of T As {Class, New}) As OperateResult(Of T)
```

```
public:
generic<typename T>
where T : ref class, gcnew()
virtual OperateResult<T>^ Read()
```

```
abstract Read : unit -> OperateResult<'T>  when 'T : not struct, new()
override Read : unit -> OperateResult<'T>  when 'T : not struct, new()
```

#### 类型参数

T
:   自定义的数据类型对象

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)T  
包含是否成功的结果对象

#### 实现

[IReadWriteNetReadT](1be4d687-47e0-f01f-2969-570e5b057c4a.htm)

![](../icons/SectionExpanded.png)示例

此处演示西门子的读取示例，先定义一个类，重点是将需要读取的数据，写入到属性的特性中去。

特性实现示例

[复制](# "复制")

```
// 假设你要读取几个数据的情况，我们把需要读取的数据定义成一个个的数量，本示例既适合单个读取，也适合批量读取，以下就是混搭的情况。
// 我们假设，我们要读取的PLC是西门子PLC，地址数据的假设如下
// 我们假设 设备是否启动是 M0.0
// 产量是 M10 开始的2个地址数据
// 温度信息是 DB1.0开始的4个地址数据
// 报警的IO信息是 M200 开始，5个字节，共计40个IO点信息
// 那么我们可以做如下的定义

public class DataExample
{
    /// <summary>
    /// 设备是否启动
    /// </summary>
    [HslDeviceAddress( "M0.0" )]
    public bool Enable { get; set; }

    /// <summary>
    /// 产量信息
    /// </summary>
    [HslDeviceAddress( "M10" )]
    public short Production { get; set; }

    /// <summary>
    /// 温度信息
    /// </summary>
    [HslDeviceAddress( "DB1.0" )]
    public float Temperature { get; set; }

    /// <summary>
    /// 连续的位报警信息
    /// </summary>
    [HslDeviceAddress( "M200", 5 )]
    public byte[] AlarmStatus { get; set; }
}
```

接下来就可以实现数据的读取了

ReadObject示例

[复制](# "复制")

```
SiemensS7Net plc = new SiemensS7Net( SiemensPLCS.S1200, "192.168.0.100" );

// 此处需要注意的是，凡是带有 HslDeviceAddress 特性的属性都会被读取出来
OperateResult<DataExample> read = plc.Read<DataExample>( );
if (read.IsSuccess)
{
    // success
    DataExample data = read.Content;
}
else
{
    // failed
    Console.WriteLine( "读取失败：" + read.Message );
}
```

![](../icons/SectionExpanded.png)参见

#### 引用

[DeviceCommunication 类](d345c980-d488-c83d-72ce-f958039275db.htm)

[Read 重载](b0f83386-169b-7cf5-0c21-7226d7fda785.htm)

[HslCommunication.Core.Device 命名空间](01bea659-9e5c-248c-f839-060b8da639c8.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Read 方法 (String, UInt16)

[原文連結](http://api.hslcommunication.cn/html/9dd09452-c200-bc65-75ff-69aa34b366d4.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Core.Device](../html/01bea659-9e5c-248c-f839-060b8da639c8.htm "HslCommunication.Core.Device")

[DeviceCommunication 类](../html/d345c980-d488-c83d-72ce-f958039275db.htm "DeviceCommunication 类")

[DeviceCommunication 方法](../html/ec3d93f5-da3b-d238-ded9-aab29ab5a132.htm "DeviceCommunication 方法")

[Read 方法](../html/b0f83386-169b-7cf5-0c21-7226d7fda785.htm "Read 方法 ")

[Read(T) 方法](../html/95abe064-5594-5ae7-a849-98c53f88fee0.htm "Read(T) 方法 ")

[Read 方法 (String, UInt16)](../html/9dd09452-c200-bc65-75ff-69aa34b366d4.htm "Read 方法 (String, UInt16)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DeviceCommunicationRead 方法 (String, UInt16) |

批量读取字节数组信息，需要指定地址和长度，返回原始的字节数组  
Batch read byte array information, need to specify the address and length, return the original byte array

**命名空间：**
 [HslCommunication.Core.Device](01bea659-9e5c-248c-f839-060b8da639c8.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public virtual OperateResult<byte[]> Read(
	string address,
	ushort length
)
```

```
Public Overridable Function Read ( 
	address As String,
	length As UShort
) As OperateResult(Of Byte())
```

```
public:
virtual OperateResult<array<unsigned char>^>^ Read(
	String^ address, 
	unsigned short length
)
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
    数据地址

length
:   类型：SystemUInt16  
    数据长度

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
带有成功标识的byte[]数组

#### 实现

[IReadWriteNetRead(String, UInt16)](b8d4187f-ab9a-79e9-c7ca-8a01ec0df233.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[DeviceCommunication 类](d345c980-d488-c83d-72ce-f958039275db.htm)

[Read 重载](b0f83386-169b-7cf5-0c21-7226d7fda785.htm)

[HslCommunication.Core.Device 命名空间](01bea659-9e5c-248c-f839-060b8da639c8.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/be30ee7c-174c-0bc1-d252-c53c3cd1b327.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Core.Device](../html/01bea659-9e5c-248c-f839-060b8da639c8.htm "HslCommunication.Core.Device")

[DeviceCommunication 类](../html/d345c980-d488-c83d-72ce-f958039275db.htm "DeviceCommunication 类")

[DeviceCommunication 方法](../html/ec3d93f5-da3b-d238-ded9-aab29ab5a132.htm "DeviceCommunication 方法")

[ReadAsync 方法](../html/be30ee7c-174c-0bc1-d252-c53c3cd1b327.htm "ReadAsync 方法 ")

[ReadAsync(T) 方法](../html/0522e740-4bf0-2e31-c4c5-e41384321b37.htm "ReadAsync(T) 方法 ")

[ReadAsync 方法 (String, UInt16)](../html/d8c14d37-56dc-2a31-1a19-3e2b5b65197d.htm "ReadAsync 方法 (String, UInt16)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DeviceCommunicationReadAsync 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法代码示例 | [ReadAsyncT](0522e740-4bf0-2e31-c4c5-e41384321b37.htm) | 异步读取支持Hsl特性的数据内容，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Asynchronously read the data content of the Hsl attribute. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. |
| 公共方法 | [ReadAsync(String, UInt16)](d8c14d37-56dc-2a31-1a19-3e2b5b65197d.htm) | 异步批量读取字节数组信息，需要指定地址和长度，返回原始的字节数组  Asynchronous batch read byte array information, need to specify the address and length, return the original byte array |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[DeviceCommunication 类](d345c980-d488-c83d-72ce-f958039275db.htm)

[HslCommunication.Core.Device 命名空间](01bea659-9e5c-248c-f839-060b8da639c8.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadAsync(T) 方法 

[原文連結](http://api.hslcommunication.cn/html/0522e740-4bf0-2e31-c4c5-e41384321b37.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Core.Device](../html/01bea659-9e5c-248c-f839-060b8da639c8.htm "HslCommunication.Core.Device")

[DeviceCommunication 类](../html/d345c980-d488-c83d-72ce-f958039275db.htm "DeviceCommunication 类")

[DeviceCommunication 方法](../html/ec3d93f5-da3b-d238-ded9-aab29ab5a132.htm "DeviceCommunication 方法")

[ReadAsync 方法](../html/be30ee7c-174c-0bc1-d252-c53c3cd1b327.htm "ReadAsync 方法 ")

[ReadAsync(T) 方法](../html/0522e740-4bf0-2e31-c4c5-e41384321b37.htm "ReadAsync(T) 方法 ")

[ReadAsync 方法 (String, UInt16)](../html/d8c14d37-56dc-2a31-1a19-3e2b5b65197d.htm "ReadAsync 方法 (String, UInt16)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DeviceCommunicationReadAsyncT 方法 |

异步读取支持Hsl特性的数据内容，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  
Asynchronously read the data content of the Hsl attribute. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details.

**命名空间：**
 [HslCommunication.Core.Device](01bea659-9e5c-248c-f839-060b8da639c8.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public virtual Task<OperateResult<T>> ReadAsync<T>()
where T : class, new()
```

```
Public Overridable Function ReadAsync(Of T As {Class, New}) As Task(Of OperateResult(Of T))
```

```
public:
generic<typename T>
where T : ref class, gcnew()
virtual Task<OperateResult<T>^>^ ReadAsync()
```

```
abstract ReadAsync : unit -> Task<OperateResult<'T>>  when 'T : not struct, new()
override ReadAsync : unit -> Task<OperateResult<'T>>  when 'T : not struct, new()
```

#### 类型参数

T
:   自定义的数据类型对象

#### 返回值

类型：Task[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)T  
包含是否成功的结果对象

#### 实现

[IReadWriteNetReadAsyncT](2c18ca3b-cee6-e5ac-8962-48d1f7a99ff4.htm)

![](../icons/SectionExpanded.png)示例

此处演示西门子的读取示例，先定义一个类，重点是将需要读取的数据，写入到属性的特性中去。

特性实现示例

[复制](# "复制")

```
// 假设你要读取几个数据的情况，我们把需要读取的数据定义成一个个的数量，本示例既适合单个读取，也适合批量读取，以下就是混搭的情况。
// 我们假设，我们要读取的PLC是西门子PLC，地址数据的假设如下
// 我们假设 设备是否启动是 M0.0
// 产量是 M10 开始的2个地址数据
// 温度信息是 DB1.0开始的4个地址数据
// 报警的IO信息是 M200 开始，5个字节，共计40个IO点信息
// 那么我们可以做如下的定义

public class DataExample
{
    /// <summary>
    /// 设备是否启动
    /// </summary>
    [HslDeviceAddress( "M0.0" )]
    public bool Enable { get; set; }

    /// <summary>
    /// 产量信息
    /// </summary>
    [HslDeviceAddress( "M10" )]
    public short Production { get; set; }

    /// <summary>
    /// 温度信息
    /// </summary>
    [HslDeviceAddress( "DB1.0" )]
    public float Temperature { get; set; }

    /// <summary>
    /// 连续的位报警信息
    /// </summary>
    [HslDeviceAddress( "M200", 5 )]
    public byte[] AlarmStatus { get; set; }
}
```

接下来就可以实现数据的读取了

ReadObjectAsync示例

[复制](# "复制")

```
SiemensS7Net plc = new SiemensS7Net( SiemensPLCS.S1200, "192.168.0.100" );

// 此处需要注意的是，凡是带有 HslDeviceAddress 特性的属性都会被读取出来
OperateResult<DataExample> read = await plc.ReadAsync<DataExample>( );
if (read.IsSuccess)
{
    // success
    DataExample data = read.Content;
}
else
{
    // failed
    Console.WriteLine( "读取失败：" + read.Message );
}
```

![](../icons/SectionExpanded.png)参见

#### 引用

[DeviceCommunication 类](d345c980-d488-c83d-72ce-f958039275db.htm)

[ReadAsync 重载](be30ee7c-174c-0bc1-d252-c53c3cd1b327.htm)

[HslCommunication.Core.Device 命名空间](01bea659-9e5c-248c-f839-060b8da639c8.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadAsync 方法 (String, UInt16)

[原文連結](http://api.hslcommunication.cn/html/d8c14d37-56dc-2a31-1a19-3e2b5b65197d.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Core.Device](../html/01bea659-9e5c-248c-f839-060b8da639c8.htm "HslCommunication.Core.Device")

[DeviceCommunication 类](../html/d345c980-d488-c83d-72ce-f958039275db.htm "DeviceCommunication 类")

[DeviceCommunication 方法](../html/ec3d93f5-da3b-d238-ded9-aab29ab5a132.htm "DeviceCommunication 方法")

[ReadAsync 方法](../html/be30ee7c-174c-0bc1-d252-c53c3cd1b327.htm "ReadAsync 方法 ")

[ReadAsync(T) 方法](../html/0522e740-4bf0-2e31-c4c5-e41384321b37.htm "ReadAsync(T) 方法 ")

[ReadAsync 方法 (String, UInt16)](../html/d8c14d37-56dc-2a31-1a19-3e2b5b65197d.htm "ReadAsync 方法 (String, UInt16)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DeviceCommunicationReadAsync 方法 (String, UInt16) |

异步批量读取字节数组信息，需要指定地址和长度，返回原始的字节数组  
Asynchronous batch read byte array information, need to specify the address and length, return the original byte array

**命名空间：**
 [HslCommunication.Core.Device](01bea659-9e5c-248c-f839-060b8da639c8.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public virtual Task<OperateResult<byte[]>> ReadAsync(
	string address,
	ushort length
)
```

```
Public Overridable Function ReadAsync ( 
	address As String,
	length As UShort
) As Task(Of OperateResult(Of Byte()))
```

```
public:
virtual Task<OperateResult<array<unsigned char>^>^>^ ReadAsync(
	String^ address, 
	unsigned short length
)
```

```
abstract ReadAsync : 
        address : string * 
        length : uint16 -> Task<OperateResult<byte[]>> 
override ReadAsync : 
        address : string * 
        length : uint16 -> Task<OperateResult<byte[]>>
```

#### 参数

address
:   类型：SystemString  
    数据地址

length
:   类型：SystemUInt16  
    数据长度

#### 返回值

类型：Task[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
带有成功标识的byte[]数组

#### 实现

[IReadWriteNetReadAsync(String, UInt16)](59839c27-5a76-2f53-730c-c9d2bb41c2c4.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[DeviceCommunication 类](d345c980-d488-c83d-72ce-f958039275db.htm)

[ReadAsync 重载](be30ee7c-174c-0bc1-d252-c53c3cd1b327.htm)

[HslCommunication.Core.Device 命名空间](01bea659-9e5c-248c-f839-060b8da639c8.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadBool 方法 

[原文連結](http://api.hslcommunication.cn/html/92f161ca-fcc4-cc3a-b4b1-565066bb5d22.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Core.Device](../html/01bea659-9e5c-248c-f839-060b8da639c8.htm "HslCommunication.Core.Device")

[DeviceCommunication 类](../html/d345c980-d488-c83d-72ce-f958039275db.htm "DeviceCommunication 类")

[DeviceCommunication 方法](../html/ec3d93f5-da3b-d238-ded9-aab29ab5a132.htm "DeviceCommunication 方法")

[ReadBool 方法](../html/92f161ca-fcc4-cc3a-b4b1-565066bb5d22.htm "ReadBool 方法 ")

[ReadBool 方法 (String)](../html/9b2faa2e-eaa9-7ace-c49f-7a0318f1df1a.htm "ReadBool 方法 (String)")

[ReadBool 方法 (String, UInt16)](../html/69d0d24f-34a7-5b41-5c1b-619edc7e6b54.htm "ReadBool 方法 (String, UInt16)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DeviceCommunicationReadBool 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [ReadBool(String)](9b2faa2e-eaa9-7ace-c49f-7a0318f1df1a.htm) | 读取单个的Boolean数据信息  Read a single Boolean data message |
| 公共方法 | [ReadBool(String, UInt16)](69d0d24f-34a7-5b41-5c1b-619edc7e6b54.htm) | 批量读取Boolean数组信息，需要指定地址和长度，返回Boolean 数组  Batch read Boolean array information, need to specify the address and length, return Boolean array |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[DeviceCommunication 类](d345c980-d488-c83d-72ce-f958039275db.htm)

[HslCommunication.Core.Device 命名空间](01bea659-9e5c-248c-f839-060b8da639c8.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadBool 方法 (String)

[原文連結](http://api.hslcommunication.cn/html/9b2faa2e-eaa9-7ace-c49f-7a0318f1df1a.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Core.Device](../html/01bea659-9e5c-248c-f839-060b8da639c8.htm "HslCommunication.Core.Device")

[DeviceCommunication 类](../html/d345c980-d488-c83d-72ce-f958039275db.htm "DeviceCommunication 类")

[DeviceCommunication 方法](../html/ec3d93f5-da3b-d238-ded9-aab29ab5a132.htm "DeviceCommunication 方法")

[ReadBool 方法](../html/92f161ca-fcc4-cc3a-b4b1-565066bb5d22.htm "ReadBool 方法 ")

[ReadBool 方法 (String)](../html/9b2faa2e-eaa9-7ace-c49f-7a0318f1df1a.htm "ReadBool 方法 (String)")

[ReadBool 方法 (String, UInt16)](../html/69d0d24f-34a7-5b41-5c1b-619edc7e6b54.htm "ReadBool 方法 (String, UInt16)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DeviceCommunicationReadBool 方法 (String) |

读取单个的Boolean数据信息  
Read a single Boolean data message

**命名空间：**
 [HslCommunication.Core.Device](01bea659-9e5c-248c-f839-060b8da639c8.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public virtual OperateResult<bool> ReadBool(
	string address
)
```

```
Public Overridable Function ReadBool ( 
	address As String
) As OperateResult(Of Boolean)
```

```
public:
virtual OperateResult<bool>^ ReadBool(
	String^ address
)
```

```
abstract ReadBool : 
        address : string -> OperateResult<bool> 
override ReadBool : 
        address : string -> OperateResult<bool>
```

#### 参数

address
:   类型：SystemString  
    数据地址

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Boolean  
带有成功标识的 bool 值

#### 实现

[IReadWriteNetReadBool(String)](67c3a6fd-ad81-62bf-a72d-74b99522a93e.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[DeviceCommunication 类](d345c980-d488-c83d-72ce-f958039275db.htm)

[ReadBool 重载](92f161ca-fcc4-cc3a-b4b1-565066bb5d22.htm)

[HslCommunication.Core.Device 命名空间](01bea659-9e5c-248c-f839-060b8da639c8.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadBool 方法 (String, UInt16)

[原文連結](http://api.hslcommunication.cn/html/69d0d24f-34a7-5b41-5c1b-619edc7e6b54.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Core.Device](../html/01bea659-9e5c-248c-f839-060b8da639c8.htm "HslCommunication.Core.Device")

[DeviceCommunication 类](../html/d345c980-d488-c83d-72ce-f958039275db.htm "DeviceCommunication 类")

[DeviceCommunication 方法](../html/ec3d93f5-da3b-d238-ded9-aab29ab5a132.htm "DeviceCommunication 方法")

[ReadBool 方法](../html/92f161ca-fcc4-cc3a-b4b1-565066bb5d22.htm "ReadBool 方法 ")

[ReadBool 方法 (String)](../html/9b2faa2e-eaa9-7ace-c49f-7a0318f1df1a.htm "ReadBool 方法 (String)")

[ReadBool 方法 (String, UInt16)](../html/69d0d24f-34a7-5b41-5c1b-619edc7e6b54.htm "ReadBool 方法 (String, UInt16)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DeviceCommunicationReadBool 方法 (String, UInt16) |

批量读取Boolean数组信息，需要指定地址和长度，返回Boolean 数组  
Batch read Boolean array information, need to specify the address and length, return Boolean array

**命名空间：**
 [HslCommunication.Core.Device](01bea659-9e5c-248c-f839-060b8da639c8.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public virtual OperateResult<bool[]> ReadBool(
	string address,
	ushort length
)
```

```
Public Overridable Function ReadBool ( 
	address As String,
	length As UShort
) As OperateResult(Of Boolean())
```

```
public:
virtual OperateResult<array<bool>^>^ ReadBool(
	String^ address, 
	unsigned short length
)
```

```
abstract ReadBool : 
        address : string * 
        length : uint16 -> OperateResult<bool[]> 
override ReadBool : 
        address : string * 
        length : uint16 -> OperateResult<bool[]>
```

#### 参数

address
:   类型：SystemString  
    数据地址

length
:   类型：SystemUInt16  
    数据长度

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Boolean  
带有成功标识的 bool[] 数组

#### 实现

[IReadWriteNetReadBool(String, UInt16)](16e00e01-28d8-bfb7-e3bb-9d8b8778d7b9.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[DeviceCommunication 类](d345c980-d488-c83d-72ce-f958039275db.htm)

[ReadBool 重载](92f161ca-fcc4-cc3a-b4b1-565066bb5d22.htm)

[HslCommunication.Core.Device 命名空间](01bea659-9e5c-248c-f839-060b8da639c8.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadBoolAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/4dba6cf8-3f49-f03d-c0aa-baf079420e7b.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Core.Device](../html/01bea659-9e5c-248c-f839-060b8da639c8.htm "HslCommunication.Core.Device")

[DeviceCommunication 类](../html/d345c980-d488-c83d-72ce-f958039275db.htm "DeviceCommunication 类")

[DeviceCommunication 方法](../html/ec3d93f5-da3b-d238-ded9-aab29ab5a132.htm "DeviceCommunication 方法")

[ReadBoolAsync 方法](../html/4dba6cf8-3f49-f03d-c0aa-baf079420e7b.htm "ReadBoolAsync 方法 ")

[ReadBoolAsync 方法 (String)](../html/cc37da1f-d673-fe67-78ee-4e465d29a2bd.htm "ReadBoolAsync 方法 (String)")

[ReadBoolAsync 方法 (String, UInt16)](../html/a45d383d-0095-3ea7-32f8-e75af88d60f8.htm "ReadBoolAsync 方法 (String, UInt16)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DeviceCommunicationReadBoolAsync 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [ReadBoolAsync(String)](cc37da1f-d673-fe67-78ee-4e465d29a2bd.htm) | 异步读取单个的Boolean数据信息  Asynchronously read a single Boolean data message |
| 公共方法 | [ReadBoolAsync(String, UInt16)](a45d383d-0095-3ea7-32f8-e75af88d60f8.htm) | 异步批量读取Boolean数组信息，需要指定地址和长度，返回Boolean 数组  Asynchronously batch read Boolean array information, need to specify the address and length, return Boolean array |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[DeviceCommunication 类](d345c980-d488-c83d-72ce-f958039275db.htm)

[HslCommunication.Core.Device 命名空间](01bea659-9e5c-248c-f839-060b8da639c8.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadBoolAsync 方法 (String)

[原文連結](http://api.hslcommunication.cn/html/cc37da1f-d673-fe67-78ee-4e465d29a2bd.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Core.Device](../html/01bea659-9e5c-248c-f839-060b8da639c8.htm "HslCommunication.Core.Device")

[DeviceCommunication 类](../html/d345c980-d488-c83d-72ce-f958039275db.htm "DeviceCommunication 类")

[DeviceCommunication 方法](../html/ec3d93f5-da3b-d238-ded9-aab29ab5a132.htm "DeviceCommunication 方法")

[ReadBoolAsync 方法](../html/4dba6cf8-3f49-f03d-c0aa-baf079420e7b.htm "ReadBoolAsync 方法 ")

[ReadBoolAsync 方法 (String)](../html/cc37da1f-d673-fe67-78ee-4e465d29a2bd.htm "ReadBoolAsync 方法 (String)")

[ReadBoolAsync 方法 (String, UInt16)](../html/a45d383d-0095-3ea7-32f8-e75af88d60f8.htm "ReadBoolAsync 方法 (String, UInt16)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DeviceCommunicationReadBoolAsync 方法 (String) |

异步读取单个的Boolean数据信息  
Asynchronously read a single Boolean data message

**命名空间：**
 [HslCommunication.Core.Device](01bea659-9e5c-248c-f839-060b8da639c8.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public virtual Task<OperateResult<bool>> ReadBoolAsync(
	string address
)
```

```
Public Overridable Function ReadBoolAsync ( 
	address As String
) As Task(Of OperateResult(Of Boolean))
```

```
public:
virtual Task<OperateResult<bool>^>^ ReadBoolAsync(
	String^ address
)
```

```
abstract ReadBoolAsync : 
        address : string -> Task<OperateResult<bool>> 
override ReadBoolAsync : 
        address : string -> Task<OperateResult<bool>>
```

#### 参数

address
:   类型：SystemString  
    数据地址

#### 返回值

类型：Task[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Boolean  
带有成功标识的byte[]数组

#### 实现

[IReadWriteNetReadBoolAsync(String)](f99fec07-3c4c-82a3-4427-4cc26281c871.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[DeviceCommunication 类](d345c980-d488-c83d-72ce-f958039275db.htm)

[ReadBoolAsync 重载](4dba6cf8-3f49-f03d-c0aa-baf079420e7b.htm)

[HslCommunication.Core.Device 命名空间](01bea659-9e5c-248c-f839-060b8da639c8.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadBoolAsync 方法 (String, UInt16)

[原文連結](http://api.hslcommunication.cn/html/a45d383d-0095-3ea7-32f8-e75af88d60f8.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Core.Device](../html/01bea659-9e5c-248c-f839-060b8da639c8.htm "HslCommunication.Core.Device")

[DeviceCommunication 类](../html/d345c980-d488-c83d-72ce-f958039275db.htm "DeviceCommunication 类")

[DeviceCommunication 方法](../html/ec3d93f5-da3b-d238-ded9-aab29ab5a132.htm "DeviceCommunication 方法")

[ReadBoolAsync 方法](../html/4dba6cf8-3f49-f03d-c0aa-baf079420e7b.htm "ReadBoolAsync 方法 ")

[ReadBoolAsync 方法 (String)](../html/cc37da1f-d673-fe67-78ee-4e465d29a2bd.htm "ReadBoolAsync 方法 (String)")

[ReadBoolAsync 方法 (String, UInt16)](../html/a45d383d-0095-3ea7-32f8-e75af88d60f8.htm "ReadBoolAsync 方法 (String, UInt16)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DeviceCommunicationReadBoolAsync 方法 (String, UInt16) |

异步批量读取Boolean数组信息，需要指定地址和长度，返回Boolean 数组  
Asynchronously batch read Boolean array information, need to specify the address and length, return Boolean array

**命名空间：**
 [HslCommunication.Core.Device](01bea659-9e5c-248c-f839-060b8da639c8.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public virtual Task<OperateResult<bool[]>> ReadBoolAsync(
	string address,
	ushort length
)
```

```
Public Overridable Function ReadBoolAsync ( 
	address As String,
	length As UShort
) As Task(Of OperateResult(Of Boolean()))
```

```
public:
virtual Task<OperateResult<array<bool>^>^>^ ReadBoolAsync(
	String^ address, 
	unsigned short length
)
```

```
abstract ReadBoolAsync : 
        address : string * 
        length : uint16 -> Task<OperateResult<bool[]>> 
override ReadBoolAsync : 
        address : string * 
        length : uint16 -> Task<OperateResult<bool[]>>
```

#### 参数

address
:   类型：SystemString  
    数据地址

length
:   类型：SystemUInt16  
    数据长度

#### 返回值

类型：Task[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Boolean  
带有成功标识的byte[]数组

#### 实现

[IReadWriteNetReadBoolAsync(String, UInt16)](74e40b23-198a-944b-6ed7-3f58ca51da5d.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[DeviceCommunication 类](d345c980-d488-c83d-72ce-f958039275db.htm)

[ReadBoolAsync 重载](4dba6cf8-3f49-f03d-c0aa-baf079420e7b.htm)

[HslCommunication.Core.Device 命名空间](01bea659-9e5c-248c-f839-060b8da639c8.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadCustomer 方法 

[原文連結](http://api.hslcommunication.cn/html/dd38653e-aa47-424e-3feb-fe58886dd68a.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Core.Device](../html/01bea659-9e5c-248c-f839-060b8da639c8.htm "HslCommunication.Core.Device")

[DeviceCommunication 类](../html/d345c980-d488-c83d-72ce-f958039275db.htm "DeviceCommunication 类")

[DeviceCommunication 方法](../html/ec3d93f5-da3b-d238-ded9-aab29ab5a132.htm "DeviceCommunication 方法")

[ReadCustomer 方法](../html/dd38653e-aa47-424e-3feb-fe58886dd68a.htm "ReadCustomer 方法 ")

[ReadCustomer(T) 方法 (String)](../html/ca069688-d8f0-7dbd-e962-86a10dce83cf.htm "ReadCustomer(T) 方法 (String)")

[ReadCustomer(T) 方法 (String, T)](../html/ca7dc1a1-c97d-689f-9819-4dae4af67206.htm "ReadCustomer(T) 方法 (String, T)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DeviceCommunicationReadCustomer 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法代码示例 | [ReadCustomerT(String)](ca069688-d8f0-7dbd-e962-86a10dce83cf.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，返回一个新的类型的实例对象。  To read a custom data type, you need to inherit from the IDataTransfer interface and return an instance object of a new type. |
| 公共方法代码示例 | [ReadCustomerT(String, T)](ca7dc1a1-c97d-689f-9819-4dae4af67206.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，传入一个实例，对这个实例进行赋值，并返回该实例的对象。  To read a custom data type, you need to inherit from the IDataTransfer interface, pass in an instance, assign a value to this instance, and return the object of the instance. |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[DeviceCommunication 类](d345c980-d488-c83d-72ce-f958039275db.htm)

[HslCommunication.Core.Device 命名空间](01bea659-9e5c-248c-f839-060b8da639c8.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadCustomer(T) 方法 (String)

[原文連結](http://api.hslcommunication.cn/html/ca069688-d8f0-7dbd-e962-86a10dce83cf.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Core.Device](../html/01bea659-9e5c-248c-f839-060b8da639c8.htm "HslCommunication.Core.Device")

[DeviceCommunication 类](../html/d345c980-d488-c83d-72ce-f958039275db.htm "DeviceCommunication 类")

[DeviceCommunication 方法](../html/ec3d93f5-da3b-d238-ded9-aab29ab5a132.htm "DeviceCommunication 方法")

[ReadCustomer 方法](../html/dd38653e-aa47-424e-3feb-fe58886dd68a.htm "ReadCustomer 方法 ")

[ReadCustomer(T) 方法 (String)](../html/ca069688-d8f0-7dbd-e962-86a10dce83cf.htm "ReadCustomer(T) 方法 (String)")

[ReadCustomer(T) 方法 (String, T)](../html/ca7dc1a1-c97d-689f-9819-4dae4af67206.htm "ReadCustomer(T) 方法 (String, T)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DeviceCommunicationReadCustomerT 方法 (String) |

读取自定义的数据类型，需要继承自IDataTransfer接口，返回一个新的类型的实例对象。  
To read a custom data type, you need to inherit from the IDataTransfer interface and return an instance object of a new type.

**命名空间：**
 [HslCommunication.Core.Device](01bea659-9e5c-248c-f839-060b8da639c8.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public OperateResult<T> ReadCustomer<T>(
	string address
)
where T : new(), IDataTransfer
```

```
Public Function ReadCustomer(Of T As {New, IDataTransfer}) ( 
	address As String
) As OperateResult(Of T)
```

```
public:
generic<typename T>
where T : gcnew(), IDataTransfer
virtual OperateResult<T>^ ReadCustomer(
	String^ address
) sealed
```

```
abstract ReadCustomer : 
        address : string -> OperateResult<'T>  when 'T : new() and IDataTransfer
override ReadCustomer : 
        address : string -> OperateResult<'T>  when 'T : new() and IDataTransfer
```

#### 参数

address
:   类型：SystemString  
    起始地址

#### 类型参数

T
:   自定义的类型

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)T  
带有成功标识的自定义类型数据

#### 实现

[IReadWriteNetReadCustomerT(String)](04bd0f87-9d33-c4d1-9420-7ca995687eb3.htm)

![](../icons/SectionExpanded.png)备注

需要是定义一个类，选择好相对于的ByteTransform实例，才能调用该方法。

![](../icons/SectionExpanded.png)示例

此处演示三菱的读取示例，先定义一个类，实现[IDataTransfer](90fc4cd3-fb3a-043c-2bd6-109124d6a78c.htm)接口

DataMy示例

[复制](# "复制")

```
public class DataMy : IDataTransfer
{
    // 根据对应的设备选择对应的实例化
    // 三菱 RegularByteTransform
    // 西门子 ReverseBytesTransform
    // Modbus及欧姆龙 ReverseWordTransform
    private IByteTransform byteTransform = new RegularByteTransform( );

    public ushort ReadCount => 5;


    public short temperature = 0;  // 温度
    public float press = 0f;       // 压力
    public int others = 0;         // 自定义的其他信息



    public void ParseSource( byte[] Content )
    {
        temperature = byteTransform.TransInt16( Content, 0 );
        press = byteTransform.TransSingle( Content, 2 );
        others = byteTransform.TransInt32( Content, 6 );
    }

    public byte[] ToSource( )
    {
        byte[] buffer = new byte[10];
        byteTransform.TransByte( temperature ).CopyTo( buffer, 0 );
        byteTransform.TransByte( press ).CopyTo( buffer, 2 );
        byteTransform.TransByte( others ).CopyTo( buffer, 6 );
        return buffer;
    }
}
```

接下来就可以实现数据的读取了

ReadCustomer示例

[复制](# "复制")

```
MelsecMcNet melsec = new MelsecMcNet( "192.168.0.100", 6000 );
OperateResult<DataMy> read = melsec.ReadCustomer<DataMy>( "M100" );
if (read.IsSuccess)
{
    // success
    DataMy data = read.Content;
}
else
{
    // failed
    Console.WriteLine( "读取失败：" + read.Message );
}
```

![](../icons/SectionExpanded.png)参见

#### 引用

[DeviceCommunication 类](d345c980-d488-c83d-72ce-f958039275db.htm)

[ReadCustomer 重载](dd38653e-aa47-424e-3feb-fe58886dd68a.htm)

[HslCommunication.Core.Device 命名空间](01bea659-9e5c-248c-f839-060b8da639c8.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadCustomer(T) 方法 (String, T)

[原文連結](http://api.hslcommunication.cn/html/ca7dc1a1-c97d-689f-9819-4dae4af67206.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Core.Device](../html/01bea659-9e5c-248c-f839-060b8da639c8.htm "HslCommunication.Core.Device")

[DeviceCommunication 类](../html/d345c980-d488-c83d-72ce-f958039275db.htm "DeviceCommunication 类")

[DeviceCommunication 方法](../html/ec3d93f5-da3b-d238-ded9-aab29ab5a132.htm "DeviceCommunication 方法")

[ReadCustomer 方法](../html/dd38653e-aa47-424e-3feb-fe58886dd68a.htm "ReadCustomer 方法 ")

[ReadCustomer(T) 方法 (String)](../html/ca069688-d8f0-7dbd-e962-86a10dce83cf.htm "ReadCustomer(T) 方法 (String)")

[ReadCustomer(T) 方法 (String, T)](../html/ca7dc1a1-c97d-689f-9819-4dae4af67206.htm "ReadCustomer(T) 方法 (String, T)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DeviceCommunicationReadCustomerT 方法 (String, T) |

读取自定义的数据类型，需要继承自IDataTransfer接口，传入一个实例，对这个实例进行赋值，并返回该实例的对象。  
To read a custom data type, you need to inherit from the IDataTransfer interface, pass in an instance,
assign a value to this instance, and return the object of the instance.

**命名空间：**
 [HslCommunication.Core.Device](01bea659-9e5c-248c-f839-060b8da639c8.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public OperateResult<T> ReadCustomer<T>(
	string address,
	T obj
)
where T : new(), IDataTransfer
```

```
Public Function ReadCustomer(Of T As {New, IDataTransfer}) ( 
	address As String,
	obj As T
) As OperateResult(Of T)
```

```
public:
generic<typename T>
where T : gcnew(), IDataTransfer
virtual OperateResult<T>^ ReadCustomer(
	String^ address, 
	T obj
) sealed
```

```
abstract ReadCustomer : 
        address : string * 
        obj : 'T -> OperateResult<'T>  when 'T : new() and IDataTransfer
override ReadCustomer : 
        address : string * 
        obj : 'T -> OperateResult<'T>  when 'T : new() and IDataTransfer
```

#### 参数

address
:   类型：SystemString  
    起始地址

obj
:   类型：T  
    实例

#### 类型参数

T
:   自定义的类型

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)T  
带有成功标识的自定义类型数据

#### 实现

[IReadWriteNetReadCustomerT(String, T)](72f8e1a8-fd49-efc7-3485-fb2684f1d948.htm)

![](../icons/SectionExpanded.png)备注

需要是定义一个类，选择好相对于的ByteTransform实例，才能调用该方法。

![](../icons/SectionExpanded.png)示例

此处演示三菱的读取示例，先定义一个类，实现[IDataTransfer](90fc4cd3-fb3a-043c-2bd6-109124d6a78c.htm)接口

DataMy示例

[复制](# "复制")

```
public class DataMy : IDataTransfer
{
    // 根据对应的设备选择对应的实例化
    // 三菱 RegularByteTransform
    // 西门子 ReverseBytesTransform
    // Modbus及欧姆龙 ReverseWordTransform
    private IByteTransform byteTransform = new RegularByteTransform( );

    public ushort ReadCount => 5;


    public short temperature = 0;  // 温度
    public float press = 0f;       // 压力
    public int others = 0;         // 自定义的其他信息



    public void ParseSource( byte[] Content )
    {
        temperature = byteTransform.TransInt16( Content, 0 );
        press = byteTransform.TransSingle( Content, 2 );
        others = byteTransform.TransInt32( Content, 6 );
    }

    public byte[] ToSource( )
    {
        byte[] buffer = new byte[10];
        byteTransform.TransByte( temperature ).CopyTo( buffer, 0 );
        byteTransform.TransByte( press ).CopyTo( buffer, 2 );
        byteTransform.TransByte( others ).CopyTo( buffer, 6 );
        return buffer;
    }
}
```

接下来就可以实现数据的读取了

ReadCustomer示例

[复制](# "复制")

```
MelsecMcNet melsec = new MelsecMcNet( "192.168.0.100", 6000 );
DataMy dataMy = new DataMy( );
OperateResult<DataMy> read = melsec.ReadCustomer<DataMy>( "M100", dataMy );
if (read.IsSuccess)
{
    // success
    Console.WriteLine( dataMy.temperature.ToString( ) );
}
else
{
    // failed
    Console.WriteLine( "读取失败：" + read.Message );
}
```

![](../icons/SectionExpanded.png)参见

#### 引用

[DeviceCommunication 类](d345c980-d488-c83d-72ce-f958039275db.htm)

[ReadCustomer 重载](dd38653e-aa47-424e-3feb-fe58886dd68a.htm)

[HslCommunication.Core.Device 命名空间](01bea659-9e5c-248c-f839-060b8da639c8.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadCustomerAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/61a05bce-10a8-f47d-5113-26dcc5e5f2e5.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Core.Device](../html/01bea659-9e5c-248c-f839-060b8da639c8.htm "HslCommunication.Core.Device")

[DeviceCommunication 类](../html/d345c980-d488-c83d-72ce-f958039275db.htm "DeviceCommunication 类")

[DeviceCommunication 方法](../html/ec3d93f5-da3b-d238-ded9-aab29ab5a132.htm "DeviceCommunication 方法")

[ReadCustomerAsync 方法](../html/61a05bce-10a8-f47d-5113-26dcc5e5f2e5.htm "ReadCustomerAsync 方法 ")

[ReadCustomerAsync(T) 方法 (String)](../html/f3fcda2b-255d-aadd-567e-c7b3df3636d9.htm "ReadCustomerAsync(T) 方法 (String)")

[ReadCustomerAsync(T) 方法 (String, T)](../html/527c9fae-4743-6ebd-7b9a-5ef3a8b7c455.htm "ReadCustomerAsync(T) 方法 (String, T)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DeviceCommunicationReadCustomerAsync 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法代码示例 | [ReadCustomerAsyncT(String)](f3fcda2b-255d-aadd-567e-c7b3df3636d9.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，返回一个新的类型的实例对象。  To read a custom data type, you need to inherit from the IDataTransfer interface and return an instance object of a new type. |
| 公共方法代码示例 | [ReadCustomerAsyncT(String, T)](527c9fae-4743-6ebd-7b9a-5ef3a8b7c455.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，传入一个实例，对这个实例进行赋值，并返回该实例的对象。  To read a custom data type, you need to inherit from the IDataTransfer interface, pass in an instance, assign a value to this instance, and return the object of the instance. |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[DeviceCommunication 类](d345c980-d488-c83d-72ce-f958039275db.htm)

[HslCommunication.Core.Device 命名空间](01bea659-9e5c-248c-f839-060b8da639c8.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadCustomerAsync(T) 方法 (String)

[原文連結](http://api.hslcommunication.cn/html/f3fcda2b-255d-aadd-567e-c7b3df3636d9.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Core.Device](../html/01bea659-9e5c-248c-f839-060b8da639c8.htm "HslCommunication.Core.Device")

[DeviceCommunication 类](../html/d345c980-d488-c83d-72ce-f958039275db.htm "DeviceCommunication 类")

[DeviceCommunication 方法](../html/ec3d93f5-da3b-d238-ded9-aab29ab5a132.htm "DeviceCommunication 方法")

[ReadCustomerAsync 方法](../html/61a05bce-10a8-f47d-5113-26dcc5e5f2e5.htm "ReadCustomerAsync 方法 ")

[ReadCustomerAsync(T) 方法 (String)](../html/f3fcda2b-255d-aadd-567e-c7b3df3636d9.htm "ReadCustomerAsync(T) 方法 (String)")

[ReadCustomerAsync(T) 方法 (String, T)](../html/527c9fae-4743-6ebd-7b9a-5ef3a8b7c455.htm "ReadCustomerAsync(T) 方法 (String, T)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DeviceCommunicationReadCustomerAsyncT 方法 (String) |

读取自定义的数据类型，需要继承自IDataTransfer接口，返回一个新的类型的实例对象。  
To read a custom data type, you need to inherit from the IDataTransfer interface and return an instance object of a new type.

**命名空间：**
 [HslCommunication.Core.Device](01bea659-9e5c-248c-f839-060b8da639c8.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public Task<OperateResult<T>> ReadCustomerAsync<T>(
	string address
)
where T : new(), IDataTransfer
```

```
Public Function ReadCustomerAsync(Of T As {New, IDataTransfer}) ( 
	address As String
) As Task(Of OperateResult(Of T))
```

```
public:
generic<typename T>
where T : gcnew(), IDataTransfer
virtual Task<OperateResult<T>^>^ ReadCustomerAsync(
	String^ address
) sealed
```

```
abstract ReadCustomerAsync : 
        address : string -> Task<OperateResult<'T>>  when 'T : new() and IDataTransfer
override ReadCustomerAsync : 
        address : string -> Task<OperateResult<'T>>  when 'T : new() and IDataTransfer
```

#### 参数

address
:   类型：SystemString  
    起始地址

#### 类型参数

T
:   自定义的类型

#### 返回值

类型：Task[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)T  
带有成功标识的自定义类型数据

#### 实现

[IReadWriteNetReadCustomerAsyncT(String)](482ef33f-8bb1-bea7-2bad-5c23bafcd173.htm)

![](../icons/SectionExpanded.png)备注

需要是定义一个类，选择好相对于的ByteTransform实例，才能调用该方法。

![](../icons/SectionExpanded.png)示例

此处演示三菱的读取示例，先定义一个类，实现[IDataTransfer](90fc4cd3-fb3a-043c-2bd6-109124d6a78c.htm)接口

DataMy示例

[复制](# "复制")

```
public class DataMy : IDataTransfer
{
    // 根据对应的设备选择对应的实例化
    // 三菱 RegularByteTransform
    // 西门子 ReverseBytesTransform
    // Modbus及欧姆龙 ReverseWordTransform
    private IByteTransform byteTransform = new RegularByteTransform( );

    public ushort ReadCount => 5;


    public short temperature = 0;  // 温度
    public float press = 0f;       // 压力
    public int others = 0;         // 自定义的其他信息



    public void ParseSource( byte[] Content )
    {
        temperature = byteTransform.TransInt16( Content, 0 );
        press = byteTransform.TransSingle( Content, 2 );
        others = byteTransform.TransInt32( Content, 6 );
    }

    public byte[] ToSource( )
    {
        byte[] buffer = new byte[10];
        byteTransform.TransByte( temperature ).CopyTo( buffer, 0 );
        byteTransform.TransByte( press ).CopyTo( buffer, 2 );
        byteTransform.TransByte( others ).CopyTo( buffer, 6 );
        return buffer;
    }
}
```

接下来就可以实现数据的读取了

ReadCustomerAsync示例

[复制](# "复制")

```
MelsecMcNet melsec = new MelsecMcNet( "192.168.0.100", 6000 );
OperateResult<DataMy> read = await melsec.ReadCustomerAsync<DataMy>( "M100" );
if (read.IsSuccess)
{
    // success
    DataMy data = read.Content;
}
else
{
    // failed
    Console.WriteLine( "读取失败：" + read.Message );
}
```

![](../icons/SectionExpanded.png)参见

#### 引用

[DeviceCommunication 类](d345c980-d488-c83d-72ce-f958039275db.htm)

[ReadCustomerAsync 重载](61a05bce-10a8-f47d-5113-26dcc5e5f2e5.htm)

[HslCommunication.Core.Device 命名空间](01bea659-9e5c-248c-f839-060b8da639c8.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadCustomerAsync(T) 方法 (String, T)

[原文連結](http://api.hslcommunication.cn/html/527c9fae-4743-6ebd-7b9a-5ef3a8b7c455.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Core.Device](../html/01bea659-9e5c-248c-f839-060b8da639c8.htm "HslCommunication.Core.Device")

[DeviceCommunication 类](../html/d345c980-d488-c83d-72ce-f958039275db.htm "DeviceCommunication 类")

[DeviceCommunication 方法](../html/ec3d93f5-da3b-d238-ded9-aab29ab5a132.htm "DeviceCommunication 方法")

[ReadCustomerAsync 方法](../html/61a05bce-10a8-f47d-5113-26dcc5e5f2e5.htm "ReadCustomerAsync 方法 ")

[ReadCustomerAsync(T) 方法 (String)](../html/f3fcda2b-255d-aadd-567e-c7b3df3636d9.htm "ReadCustomerAsync(T) 方法 (String)")

[ReadCustomerAsync(T) 方法 (String, T)](../html/527c9fae-4743-6ebd-7b9a-5ef3a8b7c455.htm "ReadCustomerAsync(T) 方法 (String, T)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DeviceCommunicationReadCustomerAsyncT 方法 (String, T) |

读取自定义的数据类型，需要继承自IDataTransfer接口，传入一个实例，对这个实例进行赋值，并返回该实例的对象。  
To read a custom data type, you need to inherit from the IDataTransfer interface, pass in an instance,
assign a value to this instance, and return the object of the instance.

**命名空间：**
 [HslCommunication.Core.Device](01bea659-9e5c-248c-f839-060b8da639c8.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public Task<OperateResult<T>> ReadCustomerAsync<T>(
	string address,
	T obj
)
where T : new(), IDataTransfer
```

```
Public Function ReadCustomerAsync(Of T As {New, IDataTransfer}) ( 
	address As String,
	obj As T
) As Task(Of OperateResult(Of T))
```

```
public:
generic<typename T>
where T : gcnew(), IDataTransfer
virtual Task<OperateResult<T>^>^ ReadCustomerAsync(
	String^ address, 
	T obj
) sealed
```

```
abstract ReadCustomerAsync : 
        address : string * 
        obj : 'T -> Task<OperateResult<'T>>  when 'T : new() and IDataTransfer
override ReadCustomerAsync : 
        address : string * 
        obj : 'T -> Task<OperateResult<'T>>  when 'T : new() and IDataTransfer
```

#### 参数

address
:   类型：SystemString  
    起始地址

obj
:   类型：T  
    实例

#### 类型参数

T
:   自定义的类型

#### 返回值

类型：Task[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)T  
带有成功标识的自定义类型数据

#### 实现

[IReadWriteNetReadCustomerAsyncT(String, T)](0405df60-0745-1be6-0101-59edaf4e4d38.htm)

![](../icons/SectionExpanded.png)备注

需要是定义一个类，选择好相对于的ByteTransform实例，才能调用该方法。

![](../icons/SectionExpanded.png)示例

此处演示三菱的读取示例，先定义一个类，实现[IDataTransfer](90fc4cd3-fb3a-043c-2bd6-109124d6a78c.htm)接口

DataMy示例

[复制](# "复制")

```
public class DataMy : IDataTransfer
{
    // 根据对应的设备选择对应的实例化
    // 三菱 RegularByteTransform
    // 西门子 ReverseBytesTransform
    // Modbus及欧姆龙 ReverseWordTransform
    private IByteTransform byteTransform = new RegularByteTransform( );

    public ushort ReadCount => 5;


    public short temperature = 0;  // 温度
    public float press = 0f;       // 压力
    public int others = 0;         // 自定义的其他信息



    public void ParseSource( byte[] Content )
    {
        temperature = byteTransform.TransInt16( Content, 0 );
        press = byteTransform.TransSingle( Content, 2 );
        others = byteTransform.TransInt32( Content, 6 );
    }

    public byte[] ToSource( )
    {
        byte[] buffer = new byte[10];
        byteTransform.TransByte( temperature ).CopyTo( buffer, 0 );
        byteTransform.TransByte( press ).CopyTo( buffer, 2 );
        byteTransform.TransByte( others ).CopyTo( buffer, 6 );
        return buffer;
    }
}
```

接下来就可以实现数据的读取了

ReadCustomerAsync示例

[复制](# "复制")

```
MelsecMcNet melsec = new MelsecMcNet( "192.168.0.100", 6000 );
DataMy dataMy = new DataMy( );
OperateResult<DataMy> read = await melsec.ReadCustomerAsync<DataMy>( "M100", dataMy );
if (read.IsSuccess)
{
    // success
    Console.WriteLine( dataMy.temperature.ToString( ) );
}
else
{
    // failed
    Console.WriteLine( "读取失败：" + read.Message );
}
```

![](../icons/SectionExpanded.png)参见

#### 引用

[DeviceCommunication 类](d345c980-d488-c83d-72ce-f958039275db.htm)

[ReadCustomerAsync 重载](61a05bce-10a8-f47d-5113-26dcc5e5f2e5.htm)

[HslCommunication.Core.Device 命名空间](01bea659-9e5c-248c-f839-060b8da639c8.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadDouble 方法 

[原文連結](http://api.hslcommunication.cn/html/3881ce17-f5a3-879f-372c-c641da32a34e.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Core.Device](../html/01bea659-9e5c-248c-f839-060b8da639c8.htm "HslCommunication.Core.Device")

[DeviceCommunication 类](../html/d345c980-d488-c83d-72ce-f958039275db.htm "DeviceCommunication 类")

[DeviceCommunication 方法](../html/ec3d93f5-da3b-d238-ded9-aab29ab5a132.htm "DeviceCommunication 方法")

[ReadDouble 方法](../html/3881ce17-f5a3-879f-372c-c641da32a34e.htm "ReadDouble 方法 ")

[ReadDouble 方法 (String)](../html/e7e31c43-71b4-e692-6bdc-ba9533e1d83d.htm "ReadDouble 方法 (String)")

[ReadDouble 方法 (String, UInt16)](../html/4eee0f05-f861-4b48-3f06-99737fb24aee.htm "ReadDouble 方法 (String, UInt16)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DeviceCommunicationReadDouble 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法代码示例 | [ReadDouble(String)](e7e31c43-71b4-e692-6bdc-ba9533e1d83d.htm) | 读取双浮点的数据  Read double floating point data |
| 公共方法代码示例 | [ReadDouble(String, UInt16)](4eee0f05-f861-4b48-3f06-99737fb24aee.htm) | 读取双浮点数据的数组  Read double floating point data array |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[DeviceCommunication 类](d345c980-d488-c83d-72ce-f958039275db.htm)

[HslCommunication.Core.Device 命名空间](01bea659-9e5c-248c-f839-060b8da639c8.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)