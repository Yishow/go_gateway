# HslCommunication - HslCommunication.ModBus

> 分類頁數: 30



---
## HslCommunication.ModBus

[原文連結](http://api.hslcommunication.cn/html/da9ac567-746e-47e5-70e1-b44577d09ce4.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.ModBus](../html/da9ac567-746e-47e5-70e1-b44577d09ce4.htm "HslCommunication.ModBus")

[IModbus 接口](../html/1d3ec4c8-6dd5-984d-e7ac-c107336985a5.htm "IModbus 接口")

[ModbusAscii 类](../html/873a8425-d84d-1a6f-ed5c-79e8d056cc76.htm "ModbusAscii 类")

[ModbusAsciiOverTcp 类](../html/619c4d60-2614-da3c-4d4e-901324d3b881.htm "ModbusAsciiOverTcp 类")

[ModbusInfo 类](../html/286c61d2-bd8e-219d-1abd-d82d0274e148.htm "ModbusInfo 类")

[ModbusMappingAddress 类](../html/7e34db06-e502-ab5c-1cd0-701b4622679a.htm "ModbusMappingAddress 类")

[ModBusMonitorAddress 类](../html/5a5ef4fb-6c8d-529b-976a-b8aba17b22f7.htm "ModBusMonitorAddress 类")

[ModbusRtu 类](../html/a0944e49-821b-30ff-a5a0-0d8c71c64280.htm "ModbusRtu 类")

[ModbusRtuOverTcp 类](../html/122ebebd-f8a8-4214-b165-280215a9543e.htm "ModbusRtuOverTcp 类")

[ModbusTcpNet 类](../html/8b0255a2-7125-a4ed-7a9c-14314e33603a.htm "ModbusTcpNet 类")

[ModbusTcpServer 类](../html/00dd6b64-ce7f-97c4-4216-20d43df665c7.htm "ModbusTcpServer 类")

[ModbusUdpNet 类](../html/ab8d6e2c-21e6-4f67-e846-315bb6345d79.htm "ModbusUdpNet 类")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslCommunication.ModBus 命名空间 |

[缺少 "N:HslCommunication.ModBus" 的 <summary> 文档]

![](../icons/SectionExpanded.png)类

|  | 类 | 说明 |
| --- | --- | --- |
| 公共类代码示例 | [ModbusAscii](873a8425-d84d-1a6f-ed5c-79e8d056cc76.htm) | Modbus-Ascii通讯协议的类库，基于rtu类库完善过来，支持标准的功能码，也支持扩展的功能码实现，地址采用富文本的形式，详细见备注说明  The client communication class of Modbus-Ascii protocol is convenient for data interaction with the server. It supports standard function codes and also supports extended function codes. The address is in rich text. For details, see the remarks. |
| 公共类代码示例 | [ModbusAsciiOverTcp](619c4d60-2614-da3c-4d4e-901324d3b881.htm) | Modbus-Ascii通讯协议的网口透传类，基于rtu类库完善过来，支持标准的功能码，也支持扩展的功能码实现，地址采用富文本的形式，详细见备注说明  The client communication class of Modbus-Ascii protocol is convenient for data interaction with the server. It supports standard function codes and also supports extended function codes. The address is in rich text. For details, see the remarks. |
| 公共类 | [ModbusInfo](286c61d2-bd8e-219d-1abd-d82d0274e148.htm) | Modbus协议相关的一些信息，包括功能码定义，报文的生成的定义等等信息  Some information related to Modbus protocol, including function code definition, definition of message generation, etc. |
| 公共类 | [ModbusMappingAddress](7e34db06-e502-ab5c-1cd0-701b4622679a.htm) | Modbus的地址映射类 |
| 公共类 | [ModBusMonitorAddress](5a5ef4fb-6c8d-529b-976a-b8aba17b22f7.htm) | 服务器端提供的数据监视服务 |
| 公共类代码示例 | [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm) | Modbus-Rtu通讯协议的类库，多项式码0xA001，支持标准的功能码，也支持扩展的功能码实现，地址采用富文本的形式，详细见备注说明  Modbus-Rtu communication protocol class library, polynomial code 0xA001, supports standard function codes, and also supports extended function code implementation. The address is in rich text. For details, see the remark |
| 公共类代码示例 | [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm) | Modbus-Rtu通讯协议的类库，多项式码0xA001，支持标准的功能码，也支持扩展的功能码实现，地址采用富文本的形式，详细见备注说明  Modbus-Rtu communication protocol class library, polynomial code 0xA001, supports standard function codes, and also supports extended function code implementation. The address is in rich text. For details, see the remark |
| 公共类代码示例 | [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm) | Modbus-Tcp协议的客户端通讯类，方便的和服务器进行数据交互，支持标准的功能码，也支持扩展的功能码实现，地址采用富文本的形式，详细见API文档说明  The client communication class of Modbus-Tcp protocol is convenient for data interaction with the server. It supports standard function codes and also supports extended function codes. The address is in rich text. For details, see the remarks. |
| 公共类代码示例 | [ModbusTcpServer](00dd6b64-ce7f-97c4-4216-20d43df665c7.htm) | Modbus的虚拟服务器，同时支持Tcp，Rtu，Ascii的机制，支持线圈，离散输入，寄存器和输入寄存器的读写操作，同时支持掩码写入功能，可以用来当做系统的数据交换池  Modbus virtual server supports Tcp and Rtu mechanisms at the same time, supports read and write operations of coils, discrete inputs, r egisters and input registers, and supports mask write function, which can be used as a system data exchange pool |
| 公共类代码示例 | [ModbusUdpNet](ab8d6e2c-21e6-4f67-e846-315bb6345d79.htm) | Modbus-Udp协议的客户端通讯类，方便的和服务器进行数据交互，支持标准的功能码，也支持扩展的功能码实现，地址采用富文本的形式，详细见备注说明  The client communication class of Modbus-Udp protocol is convenient for data interaction with the server. It supports standard function codes and also supports extended function codes. The address is in rich text. For details, see the remarks. |

![](../icons/SectionExpanded.png)接口

|  | 接口 | 说明 |
| --- | --- | --- |
| 公共接口 | [IModbus](1d3ec4c8-6dd5-984d-e7ac-c107336985a5.htm) | Modbus设备的接口，用来表示Modbus相关的设备对象，[ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm), [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm), [ModbusAscii](873a8425-d84d-1a6f-ed5c-79e8d056cc76.htm),[ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm),[ModbusUdpNet](ab8d6e2c-21e6-4f67-e846-315bb6345d79.htm)均实现了该接口信息  Modbus device interface, used to represent Modbus-related device objects, [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm), [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm),[ModbusAscii](873a8425-d84d-1a6f-ed5c-79e8d056cc76.htm),[ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm),[ModbusUdpNet](ab8d6e2c-21e6-4f67-e846-315bb6345d79.htm) all implement the interface information |

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## IModbus 接口

[原文連結](http://api.hslcommunication.cn/html/1d3ec4c8-6dd5-984d-e7ac-c107336985a5.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.ModBus](../html/da9ac567-746e-47e5-70e1-b44577d09ce4.htm "HslCommunication.ModBus")

[IModbus 接口](../html/1d3ec4c8-6dd5-984d-e7ac-c107336985a5.htm "IModbus 接口")

[IModbus 属性](../html/63ba248d-f09c-a4f5-9b92-875f94c9eb87.htm "IModbus 属性")

[IModbus 方法](../html/6a76dccd-76d8-56d4-6b83-0eb4f03215f8.htm "IModbus 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IModbus 接口 |

Modbus设备的接口，用来表示Modbus相关的设备对象，[ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm), [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm),
[ModbusAscii](873a8425-d84d-1a6f-ed5c-79e8d056cc76.htm),[ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm),[ModbusUdpNet](ab8d6e2c-21e6-4f67-e846-315bb6345d79.htm)均实现了该接口信息  
Modbus device interface, used to represent Modbus-related device objects, [ModbusTcpNet](8b0255a2-7125-a4ed-7a9c-14314e33603a.htm),
[ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm),[ModbusAscii](873a8425-d84d-1a6f-ed5c-79e8d056cc76.htm),[ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm),[ModbusUdpNet](ab8d6e2c-21e6-4f67-e846-315bb6345d79.htm) all implement the interface information

**命名空间：**
 [HslCommunication.ModBus](da9ac567-746e-47e5-70e1-b44577d09ce4.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public interface IModbus : IReadWriteDevice, 
	IReadWriteNet
```

```
Public Interface IModbus
	Inherits IReadWriteDevice, IReadWriteNet
```

```
public interface class IModbus : IReadWriteDevice, 
	IReadWriteNet
```

```
type IModbus =  
    interface
        interface IReadWriteDevice
        interface IReadWriteNet
    end
```

IModbus 类型公开以下成员。

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [AddressStartWithZero](fd8b02a4-9c23-b3d0-9b63-9550061e62b4.htm) | 获取或设置起始的地址是否从0开始，默认为True  Gets or sets whether the starting address starts from 0. The default is True |
| 公共属性 | [BroadcastStation](7cf6c140-dc5f-f15e-d8bc-1940e85ef6de.htm) | 获取或是设置当前广播模式对应的站号，广播模式意味着不接收设备方的数据返回操作，默认为 -1，表示不使用广播模式。  Gets or sets the station number corresponding to the current broadcast mode. Broadcast mode means that the data return operation of the device is not received. The default value is -1, indicating that broadcast mode is not used. |
| 公共属性代码示例 | [ByteTransform](0836ec33-d464-ed39-37ff-b2ada0d286fc.htm) | 当前的数据变换机制，当你需要从字节数据转换类型数据的时候需要。  The current data transformation mechanism is required when you need to convert type data from byte data. (继承自 [IReadWriteDevice](d987afcc-f99c-6aed-3847-5b1890504635.htm)。) |
| 公共属性 | [ConnectionId](20787592-e3ec-e0d0-27c2-abfa24654bd2.htm) | 当前连接的唯一ID号，默认为长度20的guid码加随机数组成，方便列表管理，也可以自己指定  The unique ID number of the current connection. The default is a 20-digit guid code plus a random number. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共属性 | [DataFormat](2db5f659-ab77-7f72-6129-6261b2aa58d1.htm) | 获取或设置数据解析的格式，可选ABCD, BADC，CDAB，DCBA格式，对int,uint,float,double,long,ulong类型有作用  Get or set the format of the data analysis, optional ABCD, BADC, CDAB, DCBA format, effective for int, uint, float, double, long, ulong type |
| 公共属性 | [DisableFunctionCode06](5216856f-35e3-559e-87a2-49a66b684b3b.htm) | 获取或设置是否禁用功能码05的写入操作，设置为true时，表示不允许使用功能码 06 进行写入操作。写入单个的字数据也使用 0x10 功能码  Get or set whether to disable the write operation of function code 06. When set to true, it indicates that the write operation of function code 05 is not allowed. Writing individual character data also uses the 0x10 function code |
| 公共属性 | [EnableWriteMaskCode](fb5d1fd0-4751-090d-58c3-6e54ba35e060.htm) | 获取或设置当前掩码写入的功能码是否激活状态，设置为 false 时，再执行写入位时，会通过读字，修改位，写字的方式来间接实现。  When the function code is set to false, and then the write bit is executed, it will be indirectly implemented by reading, modifying the bit, and writing the word. |
| 公共属性 | [IsStringReverse](69e67063-1e88-b0c1-e146-579f27dd881a.htm) | 字符串数据是否按照字来反转，默认为False  Whether the string data is reversed according to words. The default is False. |
| 公共属性代码示例 | [LogNet](38760c1f-22e3-f8fe-7512-15849a4b6eb0.htm) | 组件的日志工具，支持日志记录，只要实例化后，当前网络的基本信息，就以[DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm)等级进行输出  The component's logging tool supports logging. As long as the instantiation of the basic network information, the output will be output at [DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm) (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共属性 | [Station](00e1dba2-2e8a-2cf2-9849-c6dd0a5b1891.htm) | 获取或者重新修改服务器的默认站号信息，当然，你可以再读写的时候动态指定，例如地址 "s=2;100", 更详细的例子可以参考DEMO界面上的地址示例。  Get or re-modify the server's default station number information, of course, you can dynamically specify when reading and writing, such as the address "s=2; 100", a more detailed example can refer to the DEMO interface address example. |
| 公共属性 | [WordReadBatchLength](5ef2a51a-5b09-2430-5a21-8293b60de43f.htm) | 当进行字读取的时候，商业授权用户支持分批读取，本属性表示分批的长度信息，默认为 120，可以根据实际情况来设置其他值。  When reading of word type, commercial authorization users support batch reading. This attribute represents the length information of the batch, with a default of 120. Other values can be set according to the actual situation. |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [Read(String, UInt16)](b8d4187f-ab9a-79e9-c7ca-8a01ec0df233.htm) | 批量读取字节数组信息，需要指定地址和长度，返回原始的字节数组  Batch read byte array information, need to specify the address and length, return the original byte array (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadT](1be4d687-47e0-f01f-2969-570e5b057c4a.htm) | 读取支持Hsl特性的数据内容，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明，或是官网的详细文档  Read the data content of the Hsl attribute. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [ReadAsync(String, UInt16)](59839c27-5a76-2f53-730c-c9d2bb41c2c4.htm) | 异步批量读取字节数组信息，需要指定地址和长度，返回原始的字节数组  Asynchronous batch read byte array information, need to specify the address and length, return the original byte array (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadAsyncT](2c18ca3b-cee6-e5ac-8962-48d1f7a99ff4.htm) | 异步读取支持Hsl特性的数据内容，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Asynchronously read the data content of the Hsl attribute. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [ReadBool(String)](67c3a6fd-ad81-62bf-a72d-74b99522a93e.htm) | 读取单个的Boolean数据信息  Read a single Boolean data message (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [ReadBool(String, UInt16)](16e00e01-28d8-bfb7-e3bb-9d8b8778d7b9.htm) | 批量读取Boolean数组信息，需要指定地址和长度，返回Boolean 数组  Batch read Boolean array information, need to specify the address and length, return Boolean array (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [ReadBoolAsync(String)](f99fec07-3c4c-82a3-4427-4cc26281c871.htm) | 异步读取单个的Boolean数据信息  Asynchronously read a single Boolean data message (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [ReadBoolAsync(String, UInt16)](74e40b23-198a-944b-6ed7-3f58ca51da5d.htm) | 异步批量读取Boolean数组信息，需要指定地址和长度，返回Boolean 数组  Asynchronously batch read Boolean array information, need to specify the address and length, return Boolean array (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadCustomerT(String)](04bd0f87-9d33-c4d1-9420-7ca995687eb3.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，返回一个新的类型的实例对象。  To read a custom data type, you need to inherit from the IDataTransfer interface and return an instance object of a new type. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadCustomerT(String, T)](72f8e1a8-fd49-efc7-3485-fb2684f1d948.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，传入一个实例，对这个实例进行赋值，并返回该实例的对象。  To read a custom data type, you need to inherit from the IDataTransfer interface, pass in an instance, assign a value to this instance, and return the object of the instance. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadCustomerAsyncT(String)](482ef33f-8bb1-bea7-2bad-5c23bafcd173.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，返回一个新的类型的实例对象。  To read a custom data type, you need to inherit from the IDataTransfer interface and return an instance object of a new type. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadCustomerAsyncT(String, T)](0405df60-0745-1be6-0101-59edaf4e4d38.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，传入一个实例，对这个实例进行赋值，并返回该实例的对象。  To read a custom data type, you need to inherit from the IDataTransfer interface, pass in an instance, assign a value to this instance, and return the object of the instance. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadDouble(String)](18936f56-2b3c-59cc-788b-63490168e4e2.htm) | 读取双浮点的数据  Read double floating point data (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadDouble(String, UInt16)](a8468555-105b-337a-196f-b0a667856edd.htm) | 读取双浮点数据的数组  Read double floating point data array (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadDoubleAsync(String)](d902471c-124a-67a2-be35-e2feebbb1298.htm) | 异步读取双浮点的数据  Asynchronously read double floating point data (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadDoubleAsync(String, UInt16)](3a436322-c7f5-bb1d-e97b-dc6c0d6990d2.htm) | 异步读取双浮点数据的数组  Asynchronously read double floating point data array (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [ReadFile](97062c02-4031-5a46-9d12-744755610942.htm) | 使用0x14功能码读取文件的数据内容，需要指定文件编号，地址，长度，返回读取的结果数据，读取的长度为任意长度，内部自动根据实际情况来切割操作。 |
| 公共方法代码示例 | [ReadFloat(String)](95ca7390-d065-ef28-cfdb-e8c161437ecd.htm) | 读取单浮点数据  Read single floating point data (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadFloat(String, UInt16)](a8ae302d-6269-8baa-cacc-f35df9fcd131.htm) | 读取单浮点精度的数组  Read single floating point array (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadFloatAsync(String)](de191b71-b9e7-6f9c-e7a6-9a117e13fb05.htm) | 异步读取单浮点数据  Asynchronously read single floating point data (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadFloatAsync(String, UInt16)](26a13b36-422b-0520-3cb0-66a235e3d7a8.htm) | 异步读取单浮点精度的数组  Asynchronously read single floating point array (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [ReadFromCoreServer(Byte)](413a68ef-2f4b-638f-8bad-973673920de4.htm) | 将当前的数据报文发送到设备去，具体使用什么通信方式取决于设备信息，然后从设备接收数据回来，并返回给调用者。  Send the current data message to the device, the specific communication method used depends on the device information, and then receive the data back from the device and return it to the caller. (继承自 [IReadWriteDevice](d987afcc-f99c-6aed-3847-5b1890504635.htm)。) |
| 公共方法 | [ReadFromCoreServer(IEnumerableByte)](236bfe99-f977-a092-d283-c313ab8aaa7c.htm) | 将多个数据报文按顺序发到设备，并从设备接收返回的数据内容，然后拼接成一个Byte[]信息，需要重写[UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm)方法才能返回正确的结果。  Send multiple data packets to the device in sequence, and receive the returned data content from the device, and then splicing them into a Byte[] message, you need to rewrite [UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm) method to return the correct result. (继承自 [IReadWriteDevice](d987afcc-f99c-6aed-3847-5b1890504635.htm)。) |
| 公共方法 | [ReadFromCoreServerAsync(Byte)](646de5ac-9ca6-1fec-f9de-a318dbc34dc1.htm) | 将当前的数据报文发送到设备去，具体使用什么通信方式取决于设备信息，然后从设备接收数据回来，并返回给调用者。  Send the current data message to the device, the specific communication method used depends on the device information, and then receive the data back from the device and return it to the caller. (继承自 [IReadWriteDevice](d987afcc-f99c-6aed-3847-5b1890504635.htm)。) |
| 公共方法 | [ReadFromCoreServerAsync(IEnumerableByte)](90bb695c-5fc3-eb72-39f6-195819857394.htm) | 将多个数据报文按顺序发到设备，并从设备接收返回的数据内容，然后拼接成一个Byte[]信息，需要重写[UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm)方法才能返回正确的结果。  Send multiple data packets to the device in sequence, and receive the returned data content from the device, and then splicing them into a Byte[] message, you need to rewrite [UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm) method to return the correct result. (继承自 [IReadWriteDevice](d987afcc-f99c-6aed-3847-5b1890504635.htm)。) |
| 公共方法代码示例 | [ReadInt16(String)](c5308422-4397-2d19-f2de-6d4213394e3d.htm) | 读取16位的有符号的整型数据  Read 16-bit signed integer data (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadInt16(String, UInt16)](3173b3fb-e0ac-1b58-ffff-3741aa992458.htm) | 读取16位的有符号整型数组  Read 16-bit signed integer array (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadInt16Async(String)](bef4e722-646e-3196-0fc7-bf12cfa7eb57.htm) | 异步读取16位的有符号的整型数据  Asynchronously read 16-bit signed integer data (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadInt16Async(String, UInt16)](1d194744-a756-5f40-90a8-ec59633a3970.htm) | 异步读取16位的有符号整型数组  Asynchronously read 16-bit signed integer array (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadInt32(String)](51ba7f65-d19b-20dd-e81e-f3913226b5cc.htm) | 读取32位的有符号整型  Read 32-bit signed integer (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadInt32(String, UInt16)](39024553-e8a3-ceb6-946c-0812c211c218.htm) | 读取32位有符号整型数组  Read 32-bit signed integer array (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadInt32Async(String)](fcf027ad-23d6-a6a4-c715-b0c1dcf18b31.htm) | 异步读取32位的有符号整型  Asynchronously read 32-bit signed integer (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadInt32Async(String, UInt16)](f873bfc0-32e5-3303-ec0a-b0c5515754c6.htm) | 异步读取32位有符号整型数组  Asynchronously read 32-bit signed integer array (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadInt64(String)](a3bea048-fd83-2158-4580-231419212efe.htm) | 读取64位的有符号整型  Read 64-bit signed integer (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadInt64(String, UInt16)](32190a15-c77a-7831-9d59-f2419a7f8bb2.htm) | 读取64位的有符号整型数组  Read 64-bit signed integer array (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadInt64Async(String)](d2220dae-6131-a6aa-8533-a5cce335e2c4.htm) | 异步读取64位的有符号整型  Asynchronously read 64-bit signed integer (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadInt64Async(String, UInt16)](b511a119-ee98-1abc-b578-edcac930fcd4.htm) | 异步读取64位的有符号整型数组  Asynchronously read 64-bit signed integer array (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadString(String, UInt16)](a8ca8b30-6612-307e-c046-d8ff01b30df1.htm) | 读取字符串数据，默认为最常见的ASCII编码  Read string data, default is the most common ASCII encoding (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadString(String, UInt16, Encoding)](fc289def-e53b-b47d-22fe-7e02b5e354e0.htm) | 使用指定的编码，读取字符串数据  Reads string data using the specified encoding (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadStringAsync(String, UInt16)](960c5f9b-c87b-d7e9-a4e0-048e4f2f5981.htm) | 异步读取字符串数据，默认为最常见的ASCII编码  Asynchronously read string data, default is the most common ASCII encoding (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadStringAsync(String, UInt16, Encoding)](db06baf4-e674-22de-e89e-cb26646413f8.htm) | 异步使用指定的编码，读取字符串数据  Asynchronously reads string data using the specified encoding (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadStructT](50df0eb8-8f9c-ad9c-4dc9-bf75b4438d46.htm) | 读取结构体类型的数据，根据结构体自身的定义，读取原始字节数组，然后解析出实际的结构体数据，结构体需要实现[HslStructAttribute](87e427a8-0c7b-43a8-c421-c27264eda90f.htm)特性 (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadStructAsyncT](af626b7e-2adb-ecd1-f2a0-82aaa6e8b68d.htm) | 读取结构体类型的数据，根据结构体自身的定义，读取原始字节数组，然后解析出实际的结构体数据，结构体需要实现[HslStructAttribute](87e427a8-0c7b-43a8-c421-c27264eda90f.htm)特性 (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadUInt16(String)](66238795-7948-2cf9-c639-d1bb3da126dd.htm) | 读取16位的无符号整型  Read 16-bit unsigned integer (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadUInt16(String, UInt16)](28888661-1fcc-f941-5746-9124eb539fe8.htm) | 读取16位的无符号整型数组  Read 16-bit unsigned integer array (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadUInt16Async(String)](1744a25b-c440-7d1e-dee5-606b0ff485bc.htm) | 异步读取16位的无符号整型  Asynchronously read 16-bit unsigned integer (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadUInt16Async(String, UInt16)](0748f071-d6a7-581b-ed7a-d73317b16762.htm) | 异步读取16位的无符号整型数组  Asynchronously read 16-bit unsigned integer array (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadUInt32(String)](b252ad01-4035-bc3e-e297-7fd97c3ac7af.htm) | 读取32位的无符号整型  Read 32-bit unsigned integer (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadUInt32(String, UInt16)](64d8381b-05b6-a990-c03e-2c8a1b7cfb7f.htm) | 读取32位的无符号整型数组  Read 32-bit unsigned integer array (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadUInt32Async(String)](a5b848e0-82bd-2bb5-2dd9-6634063e7b48.htm) | 异步读取32位的无符号整型  Asynchronously read 32-bit unsigned integer (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadUInt32Async(String, UInt16)](59ba0cbe-08f8-5a30-6029-0a9f718c7d76.htm) | 异步读取32位的无符号整型数组  Asynchronously read 32-bit unsigned integer array (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadUInt64(String)](529c82ba-1d5e-f33b-2eb9-c862799f952f.htm) | 读取64位的无符号整型  Read 64-bit unsigned integer (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadUInt64(String, UInt16)](2af999cc-c8d4-e25a-68f1-d621e5d84af7.htm) | 读取64位的无符号整型的数组  Read 64-bit unsigned integer array (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadUInt64Async(String)](0c521a11-248c-ba7a-fc27-0e5a645611b5.htm) | 异步读取64位的无符号整型  Asynchronously read 64-bit unsigned integer (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadUInt64Async(String, UInt16)](b510070f-beb5-959b-7c6b-20f212a68dd1.htm) | 异步读取64位的无符号整型的数组  Asynchronously read 64-bit unsigned integer array (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [ReadWrite](6296fde9-8c50-c741-29da-00c5d0054a3b.htm) | 使用0x17功能码来实现同时写入并读取数据的操作，使用一条报文来实现，需要指定读取的地址，长度，写入的地址，写入的数据信息，返回读取的结果数据。  Use 0x17 function code to write and read data at the same time, and use a message to implement it, you need to specify the read address, length, written address, written data information, and return the read result data. |
| 公共方法 | [RegisteredAddressMapping](53b0eecc-8513-81c0-4a8b-a715d134e29c.htm) | 注册一个新的地址映射关系，注册地址映射关系后，就可以使用新的地址来读写Modbus数据了，通常用于其他的支持Modbus协议的PLC。  After registering a new address mapping, you can use the new address to read and write Modbus data, which is usually used for other PLCs that support the Modbus protocol. |
| 公共方法 | [TranslateToModbusAddress](68338057-dac7-e74e-f841-de013cea8da3.htm) | 将当前的地址信息转换成Modbus格式的地址，如果转换失败，返回失败的消息。默认不进行任何的转换。  Convert the current address information into a Modbus format address. If the conversion fails, a failure message will be returned. No conversion is performed by default. |
| 公共方法 | [Wait(String, Boolean, Int32, Int32)](c8c01cf3-bf99-761c-23d4-82dd54081952.htm) | 等待指定地址的Boolean值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the Boolean value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [Wait(String, Int16, Int32, Int32)](f8a820b0-cf90-668c-6fba-524bd0659042.htm) | 等待指定地址的Int16值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the Int16 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [Wait(String, Int32, Int32, Int32)](96ad6909-b660-2894-9c17-a7f4f4d4a2bd.htm) | 等待指定地址的Int32值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the Int32 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [Wait(String, Int64, Int32, Int32)](28ba4678-2aeb-7502-95e8-d6b855bbe971.htm) | 等待指定地址的Int64值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the Int64 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [Wait(String, UInt16, Int32, Int32)](44af3173-f7f4-04f5-753b-7855581a171a.htm) | 等待指定地址的UInt16值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the UInt16 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [Wait(String, UInt32, Int32, Int32)](0e877e55-9b76-f662-6841-9102a9e9b9b8.htm) | 等待指定地址的UInt32值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the UInt32 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [Wait(String, UInt64, Int32, Int32)](6801831d-c349-800d-e3d3-34d2fde09e0f.htm) | 等待指定地址的UInt64值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the UInt64 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [WaitAsync(String, Boolean, Int32, Int32)](b51daff5-9c09-c614-9ea3-f14e907a36b8.htm) | 等待指定地址的Boolean值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the Boolean value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [WaitAsync(String, Int16, Int32, Int32)](cfd8af62-371b-84e9-5695-87539fbb204b.htm) | 等待指定地址的Int16值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the Int16 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [WaitAsync(String, Int32, Int32, Int32)](639f2d1b-9a41-0cf2-0335-7f9641fbb0c8.htm) | 等待指定地址的Int32值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the Int32 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [WaitAsync(String, Int64, Int32, Int32)](4355bcb4-e923-9226-ac82-55dd1325f34f.htm) | 等待指定地址的Int64值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the Int64 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [WaitAsync(String, UInt16, Int32, Int32)](291b6f29-9cbf-6fc3-4828-226f412e9cc6.htm) | 等待指定地址的UInt16值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the UInt16 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [WaitAsync(String, UInt32, Int32, Int32)](e1009c58-9f02-660e-225a-b6861eb8781f.htm) | 等待指定地址的UInt32值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the UInt32 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [WaitAsync(String, UInt64, Int32, Int32)](4236aa40-6c9d-2a0c-8470-c55dab497074.htm) | 等待指定地址的UInt64值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the UInt64 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [Write(String, Boolean)](4d742738-6864-465c-5382-bd49032eae37.htm) | 写入单个的Boolean数据，返回是否成功  Write a single Boolean data, and return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [Write(String, Boolean)](2a99122c-1496-d766-f0aa-44a18918fcd3.htm) | 批量写入Boolean数组数据，返回是否成功  Batch write Boolean array data, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [Write(String, Byte)](b8d8d94e-300d-c7b7-2e80-6859c62479fb.htm) | 写入原始的byte数组数据到指定的地址，返回是否写入成功  Write the original byte array data to the specified address, and return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [Write(String, Double)](34e1c967-7214-bb8e-3daf-fdb93fc1db59.htm) | 写入double数据，返回是否成功  Write double data, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [Write(String, Double)](f2a8cc81-ee85-4d3a-1eaf-af33b15de5e6.htm) | 写入double数组，返回是否成功  Write double array, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [Write(String, Int16)](53b27553-763c-2dd9-b8a1-df2ebf68ab2c.htm) | 写入short数据，返回是否成功  Write short data, returns whether success (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [Write(String, Int16)](0b8fff2d-0da9-d73a-bc11-866a501e6112.htm) | 写入short数组，返回是否成功  Write short array, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [Write(String, Int32)](281ce6f6-5c53-da62-9a09-9d439390a972.htm) | 写入int数据，返回是否成功  Write int data, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [Write(String, Int32)](496e02c2-62f2-4f72-aff1-1d691cc87726.htm) | 写入int[]数组，返回是否成功  Write int array, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [Write(String, Int64)](7ec919a0-1d6b-b7d7-1f2f-eb1644beb912.htm) | 写入long数据，返回是否成功  Write long data, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [Write(String, Int64)](026100a5-d52e-b640-f239-a140c903201d.htm) | 写入long数组，返回是否成功  Write long array, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [Write(String, Single)](7c72c919-367f-9bc6-567c-550813459297.htm) | 写入float数据，返回是否成功  Write float data, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [Write(String, Single)](195c858e-1e2d-863a-8a2c-dc49cb2063d3.htm) | 写入float数组，返回是否成功  Write float array, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [Write(String, String)](fc988c8f-b2d6-18f9-3274-9b57d4a04e1c.htm) | 写入字符串信息，编码为ASCII  Write string information, encoded as ASCII (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [Write(String, UInt16)](7084bc37-6bda-e506-4118-74aff979efff.htm) | 写入ushort数据，返回是否成功  Write ushort data, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [Write(String, UInt16)](9c21745f-7c5f-cae3-7f2a-2bcf260e54a9.htm) | 写入ushort数组，返回是否成功  Write ushort array, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [Write(String, UInt32)](ed79fea9-bc6e-3ef3-5294-134b872141cd.htm) | 写入uint数据，返回是否成功  Write uint data, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [Write(String, UInt32)](3ae38bac-66cc-5cdf-be68-483362b17ac3.htm) | 写入uint[]数组，返回是否成功  Write uint array, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [Write(String, UInt64)](3927095f-0705-e823-bb24-74c9363fab36.htm) | 写入ulong数据，返回是否成功  Write ulong data, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [Write(String, UInt64)](2eb4a3ce-c350-9949-405b-ccb6b5d9786e.htm) | 写入ulong数组，返回是否成功  Write ulong array, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [Write(String, String, Int32)](1ed07021-8f14-743a-071c-c348025634f5.htm) | 写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为ASCII  Write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, it is filled with 0 and the encoding is ASCII. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [Write(String, String, Encoding)](c66a9dfb-840a-3ac4-96ef-48e23fb85ec1.htm) | 写入字符串信息，需要指定的编码信息  Write string information, need to specify the encoding information (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [Write(String, String, Int32, Encoding)](61396f95-43cd-34fb-9055-37676c4e436d.htm) | 写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为指定的编码信息  Write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, then the operation is complemented with 0 , you should specified the encoding information (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteT(T)](ec80557a-6d74-6030-7e09-57c0c6a63cc4.htm) | 写入支持Hsl特性的数据，返回是否写入成功，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Write data that supports the Hsl attribute, and return whether the write was successful. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [WriteAsync(String, Boolean)](559924da-b298-9c81-794f-83d53aafc602.htm) | 异步批量写入Boolean数组数据，返回是否成功  Asynchronously batch write Boolean array data, return success (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [WriteAsync(String, Boolean)](cf8f8e56-df3d-607b-1c52-c1a76378db39.htm) | 异步批量写入Boolean数组数据，返回是否成功  Asynchronously batch write Boolean array data, return success (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Byte)](1a2e83c8-7782-847a-5d89-08d5556ac4ad.htm) | 异步写入原始的byte数组数据到指定的地址，返回是否写入成功  Asynchronously writes the original byte array data to the specified address, and returns whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Double)](ca502ab0-a05c-9dc7-c8b0-9afd3a24cb63.htm) | 异步写入double数据，返回是否成功  Asynchronously write double data, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Double)](33c58466-53fe-e221-1d73-6650c434e2b1.htm) | 异步写入double数组，返回是否成功  Asynchronously write double array, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int16)](3e923476-059a-0de9-4fc5-6df5dc7a93a5.htm) | 异步写入short数据，返回是否成功  Asynchronously write short data, returns whether success (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int16)](8955c5b3-7709-57ce-51fa-a58369f1147f.htm) | 异步写入short数组，返回是否成功  Asynchronously write short array, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int32)](448b905c-dfee-7f97-4dd9-8caf3539eef4.htm) | 异步写入int数据，返回是否成功  Asynchronously write int data, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int32)](cd5f66ad-6ffc-5394-4554-43abb07d977a.htm) | 异步写入int[]数组，返回是否成功  Asynchronously write int array, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int64)](972d6a25-bb1a-3c92-c12a-4b6942fc026f.htm) | 异步写入long数据，返回是否成功  Asynchronously write long data, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int64)](89e51562-f954-65b5-0c20-d1602d8e78a8.htm) | 异步写入long数组，返回是否成功  Asynchronously write long array, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Single)](213794a0-4dc8-36b4-3b3f-d1ee74e35155.htm) | 异步写入float数据，返回是否成功  Asynchronously write float data, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Single)](dca9a471-a3cb-3242-b3fe-edeb04d4f786.htm) | 异步写入float数组，返回是否成功  Asynchronously write float array, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, String)](265dfc0f-8630-cc0a-4a25-cd31d50f6056.htm) | 异步写入字符串信息，编码为ASCII  Asynchronously write string information, encoded as ASCII (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt16)](ee35584e-d8c8-720e-7a05-24a2d3011021.htm) | 异步写入ushort数据，返回是否成功  Asynchronously write ushort data, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt16)](9df3d402-f0c7-5888-65d0-f527a6ada1b2.htm) | 异步写入ushort数组，返回是否成功  Asynchronously write ushort array, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt32)](76ac4f61-3a88-384c-fae1-857066a74ac0.htm) | 异步写入uint数据，返回是否成功  Asynchronously write uint data, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt32)](aae8db84-1d63-31a3-e84c-f03a5f784ff0.htm) | 异步写入uint[]数组，返回是否成功  Asynchronously write uint array, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt64)](81bbb9e0-b222-067d-7dd9-a879277859c4.htm) | 异步写入ulong数据，返回是否成功  Asynchronously write ulong data, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt64)](f56c270e-dc8d-9adc-ea60-36d47e2c6e50.htm) | 异步写入ulong数组，返回是否成功  Asynchronously write ulong array, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [WriteAsync(String, String, Int32)](c9efea95-68b3-ab09-5176-d388576844e8.htm) | 异步写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为ASCII  Asynchronously write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, it is filled with 0 and the encoding is ASCII. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, String, Encoding)](cda78ec9-ab0c-b6f1-db38-0e51e6aa06d0.htm) | 异步写入字符串信息，需要指定的编码信息  Asynchronously write string information, need to specify the encoding information (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [WriteAsync(String, String, Int32, Encoding)](c077e3ad-d170-aea3-0058-2888f2b27126.htm) | 异步写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为指定的编码信息  Asynchronously write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, then the operation is complemented with 0 , you should specified the encoding information (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteAsyncT(T)](da0f692b-4d73-0520-6e17-7890b6ef51cb.htm) | 异步写入支持Hsl特性的数据，返回是否写入成功，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Asynchronously write data that supports the Hsl attribute, and return whether the write was successful. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteCustomerT](466c49b4-95ac-148b-c18d-a419c5f6dab3.htm) | 写入自定义类型的数据，该类型必须继承自IDataTransfer接口  Write data of a custom type, which must inherit from the IDataTransfer interface (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteCustomerAsyncT](82aa8f35-3c06-d4f6-e8c5-affb6ab6ce41.htm) | 写入自定义类型的数据，该类型必须继承自IDataTransfer接口  Write data of a custom type, which must inherit from the IDataTransfer interface (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [WriteFile](63cb39f3-f1de-824f-1671-5839e496c80f.htm) | 使用0x15功能码写入文件的数据内容，需要指定文件编号，地址，数据，返回写入的结果数据。写入的数据字节长度不能大于248个字节，如果想写入更多的字节，就需要切片写入操作 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.ModBus 命名空间](da9ac567-746e-47e5-70e1-b44577d09ce4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## IModbus 属性

[原文連結](http://api.hslcommunication.cn/html/63ba248d-f09c-a4f5-9b92-875f94c9eb87.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.ModBus](../html/da9ac567-746e-47e5-70e1-b44577d09ce4.htm "HslCommunication.ModBus")

[IModbus 接口](../html/1d3ec4c8-6dd5-984d-e7ac-c107336985a5.htm "IModbus 接口")

[IModbus 属性](../html/63ba248d-f09c-a4f5-9b92-875f94c9eb87.htm "IModbus 属性")

[AddressStartWithZero 属性](../html/fd8b02a4-9c23-b3d0-9b63-9550061e62b4.htm "AddressStartWithZero 属性 ")

[BroadcastStation 属性](../html/7cf6c140-dc5f-f15e-d8bc-1940e85ef6de.htm "BroadcastStation 属性 ")

[DataFormat 属性](../html/2db5f659-ab77-7f72-6129-6261b2aa58d1.htm "DataFormat 属性 ")

[DisableFunctionCode06 属性](../html/5216856f-35e3-559e-87a2-49a66b684b3b.htm "DisableFunctionCode06 属性 ")

[EnableWriteMaskCode 属性](../html/fb5d1fd0-4751-090d-58c3-6e54ba35e060.htm "EnableWriteMaskCode 属性 ")

[IsStringReverse 属性](../html/69e67063-1e88-b0c1-e146-579f27dd881a.htm "IsStringReverse 属性 ")

[Station 属性](../html/00e1dba2-2e8a-2cf2-9849-c6dd0a5b1891.htm "Station 属性 ")

[WordReadBatchLength 属性](../html/5ef2a51a-5b09-2430-5a21-8293b60de43f.htm "WordReadBatchLength 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IModbus 属性 |

[IModbus](1d3ec4c8-6dd5-984d-e7ac-c107336985a5.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [AddressStartWithZero](fd8b02a4-9c23-b3d0-9b63-9550061e62b4.htm) | 获取或设置起始的地址是否从0开始，默认为True  Gets or sets whether the starting address starts from 0. The default is True |
| 公共属性 | [BroadcastStation](7cf6c140-dc5f-f15e-d8bc-1940e85ef6de.htm) | 获取或是设置当前广播模式对应的站号，广播模式意味着不接收设备方的数据返回操作，默认为 -1，表示不使用广播模式。  Gets or sets the station number corresponding to the current broadcast mode. Broadcast mode means that the data return operation of the device is not received. The default value is -1, indicating that broadcast mode is not used. |
| 公共属性代码示例 | [ByteTransform](0836ec33-d464-ed39-37ff-b2ada0d286fc.htm) | 当前的数据变换机制，当你需要从字节数据转换类型数据的时候需要。  The current data transformation mechanism is required when you need to convert type data from byte data. (继承自 [IReadWriteDevice](d987afcc-f99c-6aed-3847-5b1890504635.htm)。) |
| 公共属性 | [ConnectionId](20787592-e3ec-e0d0-27c2-abfa24654bd2.htm) | 当前连接的唯一ID号，默认为长度20的guid码加随机数组成，方便列表管理，也可以自己指定  The unique ID number of the current connection. The default is a 20-digit guid code plus a random number. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共属性 | [DataFormat](2db5f659-ab77-7f72-6129-6261b2aa58d1.htm) | 获取或设置数据解析的格式，可选ABCD, BADC，CDAB，DCBA格式，对int,uint,float,double,long,ulong类型有作用  Get or set the format of the data analysis, optional ABCD, BADC, CDAB, DCBA format, effective for int, uint, float, double, long, ulong type |
| 公共属性 | [DisableFunctionCode06](5216856f-35e3-559e-87a2-49a66b684b3b.htm) | 获取或设置是否禁用功能码05的写入操作，设置为true时，表示不允许使用功能码 06 进行写入操作。写入单个的字数据也使用 0x10 功能码  Get or set whether to disable the write operation of function code 06. When set to true, it indicates that the write operation of function code 05 is not allowed. Writing individual character data also uses the 0x10 function code |
| 公共属性 | [EnableWriteMaskCode](fb5d1fd0-4751-090d-58c3-6e54ba35e060.htm) | 获取或设置当前掩码写入的功能码是否激活状态，设置为 false 时，再执行写入位时，会通过读字，修改位，写字的方式来间接实现。  When the function code is set to false, and then the write bit is executed, it will be indirectly implemented by reading, modifying the bit, and writing the word. |
| 公共属性 | [IsStringReverse](69e67063-1e88-b0c1-e146-579f27dd881a.htm) | 字符串数据是否按照字来反转，默认为False  Whether the string data is reversed according to words. The default is False. |
| 公共属性代码示例 | [LogNet](38760c1f-22e3-f8fe-7512-15849a4b6eb0.htm) | 组件的日志工具，支持日志记录，只要实例化后，当前网络的基本信息，就以[DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm)等级进行输出  The component's logging tool supports logging. As long as the instantiation of the basic network information, the output will be output at [DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm) (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共属性 | [Station](00e1dba2-2e8a-2cf2-9849-c6dd0a5b1891.htm) | 获取或者重新修改服务器的默认站号信息，当然，你可以再读写的时候动态指定，例如地址 "s=2;100", 更详细的例子可以参考DEMO界面上的地址示例。  Get or re-modify the server's default station number information, of course, you can dynamically specify when reading and writing, such as the address "s=2; 100", a more detailed example can refer to the DEMO interface address example. |
| 公共属性 | [WordReadBatchLength](5ef2a51a-5b09-2430-5a21-8293b60de43f.htm) | 当进行字读取的时候，商业授权用户支持分批读取，本属性表示分批的长度信息，默认为 120，可以根据实际情况来设置其他值。  When reading of word type, commercial authorization users support batch reading. This attribute represents the length information of the batch, with a default of 120. Other values can be set according to the actual situation. |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[IModbus 接口](1d3ec4c8-6dd5-984d-e7ac-c107336985a5.htm)

[HslCommunication.ModBus 命名空间](da9ac567-746e-47e5-70e1-b44577d09ce4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## AddressStartWithZero 属性 

[原文連結](http://api.hslcommunication.cn/html/fd8b02a4-9c23-b3d0-9b63-9550061e62b4.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.ModBus](../html/da9ac567-746e-47e5-70e1-b44577d09ce4.htm "HslCommunication.ModBus")

[IModbus 接口](../html/1d3ec4c8-6dd5-984d-e7ac-c107336985a5.htm "IModbus 接口")

[IModbus 属性](../html/63ba248d-f09c-a4f5-9b92-875f94c9eb87.htm "IModbus 属性")

[AddressStartWithZero 属性](../html/fd8b02a4-9c23-b3d0-9b63-9550061e62b4.htm "AddressStartWithZero 属性 ")

[BroadcastStation 属性](../html/7cf6c140-dc5f-f15e-d8bc-1940e85ef6de.htm "BroadcastStation 属性 ")

[DataFormat 属性](../html/2db5f659-ab77-7f72-6129-6261b2aa58d1.htm "DataFormat 属性 ")

[DisableFunctionCode06 属性](../html/5216856f-35e3-559e-87a2-49a66b684b3b.htm "DisableFunctionCode06 属性 ")

[EnableWriteMaskCode 属性](../html/fb5d1fd0-4751-090d-58c3-6e54ba35e060.htm "EnableWriteMaskCode 属性 ")

[IsStringReverse 属性](../html/69e67063-1e88-b0c1-e146-579f27dd881a.htm "IsStringReverse 属性 ")

[Station 属性](../html/00e1dba2-2e8a-2cf2-9849-c6dd0a5b1891.htm "Station 属性 ")

[WordReadBatchLength 属性](../html/5ef2a51a-5b09-2430-5a21-8293b60de43f.htm "WordReadBatchLength 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IModbusAddressStartWithZero 属性 |

获取或设置起始的地址是否从0开始，默认为True  
Gets or sets whether the starting address starts from 0. The default is True

**命名空间：**
 [HslCommunication.ModBus](da9ac567-746e-47e5-70e1-b44577d09ce4.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
bool AddressStartWithZero { get; set; }
```

```
Property AddressStartWithZero As Boolean
	Get
	Set
```

```
property bool AddressStartWithZero {
	bool get ();
	void set (bool value);
}
```

```
abstract AddressStartWithZero : bool with get, set
```

#### 属性值

类型：Boolean

![](../icons/SectionExpanded.png)备注

| 警告 **警告：** |
| --- |
| 因为有些设备的起始地址是从1开始的，就要设置本属性为False |

![](../icons/SectionExpanded.png)参见

#### 引用

[IModbus 接口](1d3ec4c8-6dd5-984d-e7ac-c107336985a5.htm)

[HslCommunication.ModBus 命名空间](da9ac567-746e-47e5-70e1-b44577d09ce4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## BroadcastStation 属性 

[原文連結](http://api.hslcommunication.cn/html/7cf6c140-dc5f-f15e-d8bc-1940e85ef6de.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.ModBus](../html/da9ac567-746e-47e5-70e1-b44577d09ce4.htm "HslCommunication.ModBus")

[IModbus 接口](../html/1d3ec4c8-6dd5-984d-e7ac-c107336985a5.htm "IModbus 接口")

[IModbus 属性](../html/63ba248d-f09c-a4f5-9b92-875f94c9eb87.htm "IModbus 属性")

[AddressStartWithZero 属性](../html/fd8b02a4-9c23-b3d0-9b63-9550061e62b4.htm "AddressStartWithZero 属性 ")

[BroadcastStation 属性](../html/7cf6c140-dc5f-f15e-d8bc-1940e85ef6de.htm "BroadcastStation 属性 ")

[DataFormat 属性](../html/2db5f659-ab77-7f72-6129-6261b2aa58d1.htm "DataFormat 属性 ")

[DisableFunctionCode06 属性](../html/5216856f-35e3-559e-87a2-49a66b684b3b.htm "DisableFunctionCode06 属性 ")

[EnableWriteMaskCode 属性](../html/fb5d1fd0-4751-090d-58c3-6e54ba35e060.htm "EnableWriteMaskCode 属性 ")

[IsStringReverse 属性](../html/69e67063-1e88-b0c1-e146-579f27dd881a.htm "IsStringReverse 属性 ")

[Station 属性](../html/00e1dba2-2e8a-2cf2-9849-c6dd0a5b1891.htm "Station 属性 ")

[WordReadBatchLength 属性](../html/5ef2a51a-5b09-2430-5a21-8293b60de43f.htm "WordReadBatchLength 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IModbusBroadcastStation 属性 |

获取或是设置当前广播模式对应的站号，广播模式意味着不接收设备方的数据返回操作，默认为 -1，表示不使用广播模式。  
Gets or sets the station number corresponding to the current broadcast mode. Broadcast mode means that the data return operation of the device is not received. The default value is -1, indicating that broadcast mode is not used.

**命名空间：**
 [HslCommunication.ModBus](da9ac567-746e-47e5-70e1-b44577d09ce4.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
int BroadcastStation { get; set; }
```

```
Property BroadcastStation As Integer
	Get
	Set
```

```
property int BroadcastStation {
	int get ();
	void set (int value);
}
```

```
abstract BroadcastStation : int with get, set
```

#### 属性值

类型：Int32

![](../icons/SectionExpanded.png)参见

#### 引用

[IModbus 接口](1d3ec4c8-6dd5-984d-e7ac-c107336985a5.htm)

[HslCommunication.ModBus 命名空间](da9ac567-746e-47e5-70e1-b44577d09ce4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## DataFormat 属性 

[原文連結](http://api.hslcommunication.cn/html/2db5f659-ab77-7f72-6129-6261b2aa58d1.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.ModBus](../html/da9ac567-746e-47e5-70e1-b44577d09ce4.htm "HslCommunication.ModBus")

[IModbus 接口](../html/1d3ec4c8-6dd5-984d-e7ac-c107336985a5.htm "IModbus 接口")

[IModbus 属性](../html/63ba248d-f09c-a4f5-9b92-875f94c9eb87.htm "IModbus 属性")

[AddressStartWithZero 属性](../html/fd8b02a4-9c23-b3d0-9b63-9550061e62b4.htm "AddressStartWithZero 属性 ")

[BroadcastStation 属性](../html/7cf6c140-dc5f-f15e-d8bc-1940e85ef6de.htm "BroadcastStation 属性 ")

[DataFormat 属性](../html/2db5f659-ab77-7f72-6129-6261b2aa58d1.htm "DataFormat 属性 ")

[DisableFunctionCode06 属性](../html/5216856f-35e3-559e-87a2-49a66b684b3b.htm "DisableFunctionCode06 属性 ")

[EnableWriteMaskCode 属性](../html/fb5d1fd0-4751-090d-58c3-6e54ba35e060.htm "EnableWriteMaskCode 属性 ")

[IsStringReverse 属性](../html/69e67063-1e88-b0c1-e146-579f27dd881a.htm "IsStringReverse 属性 ")

[Station 属性](../html/00e1dba2-2e8a-2cf2-9849-c6dd0a5b1891.htm "Station 属性 ")

[WordReadBatchLength 属性](../html/5ef2a51a-5b09-2430-5a21-8293b60de43f.htm "WordReadBatchLength 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IModbusDataFormat 属性 |

获取或设置数据解析的格式，可选ABCD, BADC，CDAB，DCBA格式，对int,uint,float,double,long,ulong类型有作用  
Get or set the format of the data analysis, optional ABCD, BADC, CDAB, DCBA format, effective for int, uint, float, double, long, ulong type

**命名空间：**
 [HslCommunication.ModBus](da9ac567-746e-47e5-70e1-b44577d09ce4.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
DataFormat DataFormat { get; set; }
```

```
Property DataFormat As DataFormat
	Get
	Set
```

```
property DataFormat DataFormat {
	DataFormat get ();
	void set (DataFormat value);
}
```

```
abstract DataFormat : DataFormat with get, set
```

#### 属性值

类型：[DataFormat](4e3d381b-5b6d-f012-7658-19a29e76e059.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[IModbus 接口](1d3ec4c8-6dd5-984d-e7ac-c107336985a5.htm)

[HslCommunication.ModBus 命名空间](da9ac567-746e-47e5-70e1-b44577d09ce4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## DisableFunctionCode06 属性 

[原文連結](http://api.hslcommunication.cn/html/5216856f-35e3-559e-87a2-49a66b684b3b.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.ModBus](../html/da9ac567-746e-47e5-70e1-b44577d09ce4.htm "HslCommunication.ModBus")

[IModbus 接口](../html/1d3ec4c8-6dd5-984d-e7ac-c107336985a5.htm "IModbus 接口")

[IModbus 属性](../html/63ba248d-f09c-a4f5-9b92-875f94c9eb87.htm "IModbus 属性")

[AddressStartWithZero 属性](../html/fd8b02a4-9c23-b3d0-9b63-9550061e62b4.htm "AddressStartWithZero 属性 ")

[BroadcastStation 属性](../html/7cf6c140-dc5f-f15e-d8bc-1940e85ef6de.htm "BroadcastStation 属性 ")

[DataFormat 属性](../html/2db5f659-ab77-7f72-6129-6261b2aa58d1.htm "DataFormat 属性 ")

[DisableFunctionCode06 属性](../html/5216856f-35e3-559e-87a2-49a66b684b3b.htm "DisableFunctionCode06 属性 ")

[EnableWriteMaskCode 属性](../html/fb5d1fd0-4751-090d-58c3-6e54ba35e060.htm "EnableWriteMaskCode 属性 ")

[IsStringReverse 属性](../html/69e67063-1e88-b0c1-e146-579f27dd881a.htm "IsStringReverse 属性 ")

[Station 属性](../html/00e1dba2-2e8a-2cf2-9849-c6dd0a5b1891.htm "Station 属性 ")

[WordReadBatchLength 属性](../html/5ef2a51a-5b09-2430-5a21-8293b60de43f.htm "WordReadBatchLength 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IModbusDisableFunctionCode06 属性 |

获取或设置是否禁用功能码05的写入操作，设置为true时，表示不允许使用功能码 06 进行写入操作。写入单个的字数据也使用 0x10 功能码  
Get or set whether to disable the write operation of function code 06. When set to true, it indicates that the write operation of function code 05 is not allowed. Writing individual character data also uses the 0x10 function code

**命名空间：**
 [HslCommunication.ModBus](da9ac567-746e-47e5-70e1-b44577d09ce4.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
bool DisableFunctionCode06 { get; set; }
```

```
Property DisableFunctionCode06 As Boolean
	Get
	Set
```

```
property bool DisableFunctionCode06 {
	bool get ();
	void set (bool value);
}
```

```
abstract DisableFunctionCode06 : bool with get, set
```

#### 属性值

类型：Boolean

![](../icons/SectionExpanded.png)参见

#### 引用

[IModbus 接口](1d3ec4c8-6dd5-984d-e7ac-c107336985a5.htm)

[HslCommunication.ModBus 命名空间](da9ac567-746e-47e5-70e1-b44577d09ce4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## EnableWriteMaskCode 属性 

[原文連結](http://api.hslcommunication.cn/html/fb5d1fd0-4751-090d-58c3-6e54ba35e060.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.ModBus](../html/da9ac567-746e-47e5-70e1-b44577d09ce4.htm "HslCommunication.ModBus")

[IModbus 接口](../html/1d3ec4c8-6dd5-984d-e7ac-c107336985a5.htm "IModbus 接口")

[IModbus 属性](../html/63ba248d-f09c-a4f5-9b92-875f94c9eb87.htm "IModbus 属性")

[AddressStartWithZero 属性](../html/fd8b02a4-9c23-b3d0-9b63-9550061e62b4.htm "AddressStartWithZero 属性 ")

[BroadcastStation 属性](../html/7cf6c140-dc5f-f15e-d8bc-1940e85ef6de.htm "BroadcastStation 属性 ")

[DataFormat 属性](../html/2db5f659-ab77-7f72-6129-6261b2aa58d1.htm "DataFormat 属性 ")

[DisableFunctionCode06 属性](../html/5216856f-35e3-559e-87a2-49a66b684b3b.htm "DisableFunctionCode06 属性 ")

[EnableWriteMaskCode 属性](../html/fb5d1fd0-4751-090d-58c3-6e54ba35e060.htm "EnableWriteMaskCode 属性 ")

[IsStringReverse 属性](../html/69e67063-1e88-b0c1-e146-579f27dd881a.htm "IsStringReverse 属性 ")

[Station 属性](../html/00e1dba2-2e8a-2cf2-9849-c6dd0a5b1891.htm "Station 属性 ")

[WordReadBatchLength 属性](../html/5ef2a51a-5b09-2430-5a21-8293b60de43f.htm "WordReadBatchLength 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IModbusEnableWriteMaskCode 属性 |

获取或设置当前掩码写入的功能码是否激活状态，设置为 false 时，再执行写入位时，会通过读字，修改位，写字的方式来间接实现。  
When the function code is set to false, and then the write bit is executed, it will be indirectly implemented by reading, modifying the bit, and writing the word.

**命名空间：**
 [HslCommunication.ModBus](da9ac567-746e-47e5-70e1-b44577d09ce4.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
bool EnableWriteMaskCode { get; set; }
```

```
Property EnableWriteMaskCode As Boolean
	Get
	Set
```

```
property bool EnableWriteMaskCode {
	bool get ();
	void set (bool value);
}
```

```
abstract EnableWriteMaskCode : bool with get, set
```

#### 属性值

类型：Boolean

![](../icons/SectionExpanded.png)参见

#### 引用

[IModbus 接口](1d3ec4c8-6dd5-984d-e7ac-c107336985a5.htm)

[HslCommunication.ModBus 命名空间](da9ac567-746e-47e5-70e1-b44577d09ce4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## IsStringReverse 属性 

[原文連結](http://api.hslcommunication.cn/html/69e67063-1e88-b0c1-e146-579f27dd881a.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.ModBus](../html/da9ac567-746e-47e5-70e1-b44577d09ce4.htm "HslCommunication.ModBus")

[IModbus 接口](../html/1d3ec4c8-6dd5-984d-e7ac-c107336985a5.htm "IModbus 接口")

[IModbus 属性](../html/63ba248d-f09c-a4f5-9b92-875f94c9eb87.htm "IModbus 属性")

[AddressStartWithZero 属性](../html/fd8b02a4-9c23-b3d0-9b63-9550061e62b4.htm "AddressStartWithZero 属性 ")

[BroadcastStation 属性](../html/7cf6c140-dc5f-f15e-d8bc-1940e85ef6de.htm "BroadcastStation 属性 ")

[DataFormat 属性](../html/2db5f659-ab77-7f72-6129-6261b2aa58d1.htm "DataFormat 属性 ")

[DisableFunctionCode06 属性](../html/5216856f-35e3-559e-87a2-49a66b684b3b.htm "DisableFunctionCode06 属性 ")

[EnableWriteMaskCode 属性](../html/fb5d1fd0-4751-090d-58c3-6e54ba35e060.htm "EnableWriteMaskCode 属性 ")

[IsStringReverse 属性](../html/69e67063-1e88-b0c1-e146-579f27dd881a.htm "IsStringReverse 属性 ")

[Station 属性](../html/00e1dba2-2e8a-2cf2-9849-c6dd0a5b1891.htm "Station 属性 ")

[WordReadBatchLength 属性](../html/5ef2a51a-5b09-2430-5a21-8293b60de43f.htm "WordReadBatchLength 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IModbusIsStringReverse 属性 |

字符串数据是否按照字来反转，默认为False  
Whether the string data is reversed according to words. The default is False.

**命名空间：**
 [HslCommunication.ModBus](da9ac567-746e-47e5-70e1-b44577d09ce4.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
bool IsStringReverse { get; set; }
```

```
Property IsStringReverse As Boolean
	Get
	Set
```

```
property bool IsStringReverse {
	bool get ();
	void set (bool value);
}
```

```
abstract IsStringReverse : bool with get, set
```

#### 属性值

类型：Boolean

![](../icons/SectionExpanded.png)备注

字符串按照2个字节的排列进行颠倒，根据实际情况进行设置

![](../icons/SectionExpanded.png)参见

#### 引用

[IModbus 接口](1d3ec4c8-6dd5-984d-e7ac-c107336985a5.htm)

[HslCommunication.ModBus 命名空间](da9ac567-746e-47e5-70e1-b44577d09ce4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Station 属性 

[原文連結](http://api.hslcommunication.cn/html/00e1dba2-2e8a-2cf2-9849-c6dd0a5b1891.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.ModBus](../html/da9ac567-746e-47e5-70e1-b44577d09ce4.htm "HslCommunication.ModBus")

[IModbus 接口](../html/1d3ec4c8-6dd5-984d-e7ac-c107336985a5.htm "IModbus 接口")

[IModbus 属性](../html/63ba248d-f09c-a4f5-9b92-875f94c9eb87.htm "IModbus 属性")

[AddressStartWithZero 属性](../html/fd8b02a4-9c23-b3d0-9b63-9550061e62b4.htm "AddressStartWithZero 属性 ")

[BroadcastStation 属性](../html/7cf6c140-dc5f-f15e-d8bc-1940e85ef6de.htm "BroadcastStation 属性 ")

[DataFormat 属性](../html/2db5f659-ab77-7f72-6129-6261b2aa58d1.htm "DataFormat 属性 ")

[DisableFunctionCode06 属性](../html/5216856f-35e3-559e-87a2-49a66b684b3b.htm "DisableFunctionCode06 属性 ")

[EnableWriteMaskCode 属性](../html/fb5d1fd0-4751-090d-58c3-6e54ba35e060.htm "EnableWriteMaskCode 属性 ")

[IsStringReverse 属性](../html/69e67063-1e88-b0c1-e146-579f27dd881a.htm "IsStringReverse 属性 ")

[Station 属性](../html/00e1dba2-2e8a-2cf2-9849-c6dd0a5b1891.htm "Station 属性 ")

[WordReadBatchLength 属性](../html/5ef2a51a-5b09-2430-5a21-8293b60de43f.htm "WordReadBatchLength 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IModbusStation 属性 |

获取或者重新修改服务器的默认站号信息，当然，你可以再读写的时候动态指定，例如地址 "s=2;100", 更详细的例子可以参考DEMO界面上的地址示例。  
Get or re-modify the server's default station number information, of course, you can dynamically specify when reading and writing,
such as the address "s=2; 100", a more detailed example can refer to the DEMO interface address example.

**命名空间：**
 [HslCommunication.ModBus](da9ac567-746e-47e5-70e1-b44577d09ce4.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
byte Station { get; set; }
```

```
Property Station As Byte
	Get
	Set
```

```
property unsigned char Station {
	unsigned char get ();
	void set (unsigned char value);
}
```

```
abstract Station : byte with get, set
```

#### 属性值

类型：Byte

![](../icons/SectionExpanded.png)备注

当你调用 ReadCoil("100") 时，对应的站号就是本属性的值，当你调用 ReadCoil("s=2;100") 时，就忽略本属性的值，读写寄存器的时候同理

![](../icons/SectionExpanded.png)参见

#### 引用

[IModbus 接口](1d3ec4c8-6dd5-984d-e7ac-c107336985a5.htm)

[HslCommunication.ModBus 命名空间](da9ac567-746e-47e5-70e1-b44577d09ce4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## WordReadBatchLength 属性 

[原文連結](http://api.hslcommunication.cn/html/5ef2a51a-5b09-2430-5a21-8293b60de43f.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.ModBus](../html/da9ac567-746e-47e5-70e1-b44577d09ce4.htm "HslCommunication.ModBus")

[IModbus 接口](../html/1d3ec4c8-6dd5-984d-e7ac-c107336985a5.htm "IModbus 接口")

[IModbus 属性](../html/63ba248d-f09c-a4f5-9b92-875f94c9eb87.htm "IModbus 属性")

[AddressStartWithZero 属性](../html/fd8b02a4-9c23-b3d0-9b63-9550061e62b4.htm "AddressStartWithZero 属性 ")

[BroadcastStation 属性](../html/7cf6c140-dc5f-f15e-d8bc-1940e85ef6de.htm "BroadcastStation 属性 ")

[DataFormat 属性](../html/2db5f659-ab77-7f72-6129-6261b2aa58d1.htm "DataFormat 属性 ")

[DisableFunctionCode06 属性](../html/5216856f-35e3-559e-87a2-49a66b684b3b.htm "DisableFunctionCode06 属性 ")

[EnableWriteMaskCode 属性](../html/fb5d1fd0-4751-090d-58c3-6e54ba35e060.htm "EnableWriteMaskCode 属性 ")

[IsStringReverse 属性](../html/69e67063-1e88-b0c1-e146-579f27dd881a.htm "IsStringReverse 属性 ")

[Station 属性](../html/00e1dba2-2e8a-2cf2-9849-c6dd0a5b1891.htm "Station 属性 ")

[WordReadBatchLength 属性](../html/5ef2a51a-5b09-2430-5a21-8293b60de43f.htm "WordReadBatchLength 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IModbusWordReadBatchLength 属性 |

当进行字读取的时候，商业授权用户支持分批读取，本属性表示分批的长度信息，默认为 120，可以根据实际情况来设置其他值。  
When reading of word type, commercial authorization users support batch reading. This attribute represents the length information of the batch, with a default of 120. Other values can be set according to the actual situation.

**命名空间：**
 [HslCommunication.ModBus](da9ac567-746e-47e5-70e1-b44577d09ce4.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
int WordReadBatchLength { get; set; }
```

```
Property WordReadBatchLength As Integer
	Get
	Set
```

```
property int WordReadBatchLength {
	int get ();
	void set (int value);
}
```

```
abstract WordReadBatchLength : int with get, set
```

#### 属性值

类型：Int32

![](../icons/SectionExpanded.png)参见

#### 引用

[IModbus 接口](1d3ec4c8-6dd5-984d-e7ac-c107336985a5.htm)

[HslCommunication.ModBus 命名空间](da9ac567-746e-47e5-70e1-b44577d09ce4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## IModbus 方法

[原文連結](http://api.hslcommunication.cn/html/6a76dccd-76d8-56d4-6b83-0eb4f03215f8.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.ModBus](../html/da9ac567-746e-47e5-70e1-b44577d09ce4.htm "HslCommunication.ModBus")

[IModbus 接口](../html/1d3ec4c8-6dd5-984d-e7ac-c107336985a5.htm "IModbus 接口")

[IModbus 方法](../html/6a76dccd-76d8-56d4-6b83-0eb4f03215f8.htm "IModbus 方法")

[ReadFile 方法](../html/97062c02-4031-5a46-9d12-744755610942.htm "ReadFile 方法 ")

[ReadWrite 方法](../html/6296fde9-8c50-c741-29da-00c5d0054a3b.htm "ReadWrite 方法 ")

[RegisteredAddressMapping 方法](../html/53b0eecc-8513-81c0-4a8b-a715d134e29c.htm "RegisteredAddressMapping 方法 ")

[TranslateToModbusAddress 方法](../html/68338057-dac7-e74e-f841-de013cea8da3.htm "TranslateToModbusAddress 方法 ")

[WriteFile 方法](../html/63cb39f3-f1de-824f-1671-5839e496c80f.htm "WriteFile 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IModbus 方法 |

[IModbus](1d3ec4c8-6dd5-984d-e7ac-c107336985a5.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [Read(String, UInt16)](b8d4187f-ab9a-79e9-c7ca-8a01ec0df233.htm) | 批量读取字节数组信息，需要指定地址和长度，返回原始的字节数组  Batch read byte array information, need to specify the address and length, return the original byte array (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadT](1be4d687-47e0-f01f-2969-570e5b057c4a.htm) | 读取支持Hsl特性的数据内容，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明，或是官网的详细文档  Read the data content of the Hsl attribute. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [ReadAsync(String, UInt16)](59839c27-5a76-2f53-730c-c9d2bb41c2c4.htm) | 异步批量读取字节数组信息，需要指定地址和长度，返回原始的字节数组  Asynchronous batch read byte array information, need to specify the address and length, return the original byte array (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadAsyncT](2c18ca3b-cee6-e5ac-8962-48d1f7a99ff4.htm) | 异步读取支持Hsl特性的数据内容，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Asynchronously read the data content of the Hsl attribute. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [ReadBool(String)](67c3a6fd-ad81-62bf-a72d-74b99522a93e.htm) | 读取单个的Boolean数据信息  Read a single Boolean data message (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [ReadBool(String, UInt16)](16e00e01-28d8-bfb7-e3bb-9d8b8778d7b9.htm) | 批量读取Boolean数组信息，需要指定地址和长度，返回Boolean 数组  Batch read Boolean array information, need to specify the address and length, return Boolean array (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [ReadBoolAsync(String)](f99fec07-3c4c-82a3-4427-4cc26281c871.htm) | 异步读取单个的Boolean数据信息  Asynchronously read a single Boolean data message (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [ReadBoolAsync(String, UInt16)](74e40b23-198a-944b-6ed7-3f58ca51da5d.htm) | 异步批量读取Boolean数组信息，需要指定地址和长度，返回Boolean 数组  Asynchronously batch read Boolean array information, need to specify the address and length, return Boolean array (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadCustomerT(String)](04bd0f87-9d33-c4d1-9420-7ca995687eb3.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，返回一个新的类型的实例对象。  To read a custom data type, you need to inherit from the IDataTransfer interface and return an instance object of a new type. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadCustomerT(String, T)](72f8e1a8-fd49-efc7-3485-fb2684f1d948.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，传入一个实例，对这个实例进行赋值，并返回该实例的对象。  To read a custom data type, you need to inherit from the IDataTransfer interface, pass in an instance, assign a value to this instance, and return the object of the instance. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadCustomerAsyncT(String)](482ef33f-8bb1-bea7-2bad-5c23bafcd173.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，返回一个新的类型的实例对象。  To read a custom data type, you need to inherit from the IDataTransfer interface and return an instance object of a new type. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadCustomerAsyncT(String, T)](0405df60-0745-1be6-0101-59edaf4e4d38.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，传入一个实例，对这个实例进行赋值，并返回该实例的对象。  To read a custom data type, you need to inherit from the IDataTransfer interface, pass in an instance, assign a value to this instance, and return the object of the instance. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadDouble(String)](18936f56-2b3c-59cc-788b-63490168e4e2.htm) | 读取双浮点的数据  Read double floating point data (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadDouble(String, UInt16)](a8468555-105b-337a-196f-b0a667856edd.htm) | 读取双浮点数据的数组  Read double floating point data array (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadDoubleAsync(String)](d902471c-124a-67a2-be35-e2feebbb1298.htm) | 异步读取双浮点的数据  Asynchronously read double floating point data (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadDoubleAsync(String, UInt16)](3a436322-c7f5-bb1d-e97b-dc6c0d6990d2.htm) | 异步读取双浮点数据的数组  Asynchronously read double floating point data array (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [ReadFile](97062c02-4031-5a46-9d12-744755610942.htm) | 使用0x14功能码读取文件的数据内容，需要指定文件编号，地址，长度，返回读取的结果数据，读取的长度为任意长度，内部自动根据实际情况来切割操作。 |
| 公共方法代码示例 | [ReadFloat(String)](95ca7390-d065-ef28-cfdb-e8c161437ecd.htm) | 读取单浮点数据  Read single floating point data (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadFloat(String, UInt16)](a8ae302d-6269-8baa-cacc-f35df9fcd131.htm) | 读取单浮点精度的数组  Read single floating point array (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadFloatAsync(String)](de191b71-b9e7-6f9c-e7a6-9a117e13fb05.htm) | 异步读取单浮点数据  Asynchronously read single floating point data (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadFloatAsync(String, UInt16)](26a13b36-422b-0520-3cb0-66a235e3d7a8.htm) | 异步读取单浮点精度的数组  Asynchronously read single floating point array (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [ReadFromCoreServer(Byte)](413a68ef-2f4b-638f-8bad-973673920de4.htm) | 将当前的数据报文发送到设备去，具体使用什么通信方式取决于设备信息，然后从设备接收数据回来，并返回给调用者。  Send the current data message to the device, the specific communication method used depends on the device information, and then receive the data back from the device and return it to the caller. (继承自 [IReadWriteDevice](d987afcc-f99c-6aed-3847-5b1890504635.htm)。) |
| 公共方法 | [ReadFromCoreServer(IEnumerableByte)](236bfe99-f977-a092-d283-c313ab8aaa7c.htm) | 将多个数据报文按顺序发到设备，并从设备接收返回的数据内容，然后拼接成一个Byte[]信息，需要重写[UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm)方法才能返回正确的结果。  Send multiple data packets to the device in sequence, and receive the returned data content from the device, and then splicing them into a Byte[] message, you need to rewrite [UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm) method to return the correct result. (继承自 [IReadWriteDevice](d987afcc-f99c-6aed-3847-5b1890504635.htm)。) |
| 公共方法 | [ReadFromCoreServerAsync(Byte)](646de5ac-9ca6-1fec-f9de-a318dbc34dc1.htm) | 将当前的数据报文发送到设备去，具体使用什么通信方式取决于设备信息，然后从设备接收数据回来，并返回给调用者。  Send the current data message to the device, the specific communication method used depends on the device information, and then receive the data back from the device and return it to the caller. (继承自 [IReadWriteDevice](d987afcc-f99c-6aed-3847-5b1890504635.htm)。) |
| 公共方法 | [ReadFromCoreServerAsync(IEnumerableByte)](90bb695c-5fc3-eb72-39f6-195819857394.htm) | 将多个数据报文按顺序发到设备，并从设备接收返回的数据内容，然后拼接成一个Byte[]信息，需要重写[UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm)方法才能返回正确的结果。  Send multiple data packets to the device in sequence, and receive the returned data content from the device, and then splicing them into a Byte[] message, you need to rewrite [UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm) method to return the correct result. (继承自 [IReadWriteDevice](d987afcc-f99c-6aed-3847-5b1890504635.htm)。) |
| 公共方法代码示例 | [ReadInt16(String)](c5308422-4397-2d19-f2de-6d4213394e3d.htm) | 读取16位的有符号的整型数据  Read 16-bit signed integer data (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadInt16(String, UInt16)](3173b3fb-e0ac-1b58-ffff-3741aa992458.htm) | 读取16位的有符号整型数组  Read 16-bit signed integer array (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadInt16Async(String)](bef4e722-646e-3196-0fc7-bf12cfa7eb57.htm) | 异步读取16位的有符号的整型数据  Asynchronously read 16-bit signed integer data (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadInt16Async(String, UInt16)](1d194744-a756-5f40-90a8-ec59633a3970.htm) | 异步读取16位的有符号整型数组  Asynchronously read 16-bit signed integer array (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadInt32(String)](51ba7f65-d19b-20dd-e81e-f3913226b5cc.htm) | 读取32位的有符号整型  Read 32-bit signed integer (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadInt32(String, UInt16)](39024553-e8a3-ceb6-946c-0812c211c218.htm) | 读取32位有符号整型数组  Read 32-bit signed integer array (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadInt32Async(String)](fcf027ad-23d6-a6a4-c715-b0c1dcf18b31.htm) | 异步读取32位的有符号整型  Asynchronously read 32-bit signed integer (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadInt32Async(String, UInt16)](f873bfc0-32e5-3303-ec0a-b0c5515754c6.htm) | 异步读取32位有符号整型数组  Asynchronously read 32-bit signed integer array (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadInt64(String)](a3bea048-fd83-2158-4580-231419212efe.htm) | 读取64位的有符号整型  Read 64-bit signed integer (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadInt64(String, UInt16)](32190a15-c77a-7831-9d59-f2419a7f8bb2.htm) | 读取64位的有符号整型数组  Read 64-bit signed integer array (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadInt64Async(String)](d2220dae-6131-a6aa-8533-a5cce335e2c4.htm) | 异步读取64位的有符号整型  Asynchronously read 64-bit signed integer (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadInt64Async(String, UInt16)](b511a119-ee98-1abc-b578-edcac930fcd4.htm) | 异步读取64位的有符号整型数组  Asynchronously read 64-bit signed integer array (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadString(String, UInt16)](a8ca8b30-6612-307e-c046-d8ff01b30df1.htm) | 读取字符串数据，默认为最常见的ASCII编码  Read string data, default is the most common ASCII encoding (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadString(String, UInt16, Encoding)](fc289def-e53b-b47d-22fe-7e02b5e354e0.htm) | 使用指定的编码，读取字符串数据  Reads string data using the specified encoding (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadStringAsync(String, UInt16)](960c5f9b-c87b-d7e9-a4e0-048e4f2f5981.htm) | 异步读取字符串数据，默认为最常见的ASCII编码  Asynchronously read string data, default is the most common ASCII encoding (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadStringAsync(String, UInt16, Encoding)](db06baf4-e674-22de-e89e-cb26646413f8.htm) | 异步使用指定的编码，读取字符串数据  Asynchronously reads string data using the specified encoding (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadStructT](50df0eb8-8f9c-ad9c-4dc9-bf75b4438d46.htm) | 读取结构体类型的数据，根据结构体自身的定义，读取原始字节数组，然后解析出实际的结构体数据，结构体需要实现[HslStructAttribute](87e427a8-0c7b-43a8-c421-c27264eda90f.htm)特性 (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadStructAsyncT](af626b7e-2adb-ecd1-f2a0-82aaa6e8b68d.htm) | 读取结构体类型的数据，根据结构体自身的定义，读取原始字节数组，然后解析出实际的结构体数据，结构体需要实现[HslStructAttribute](87e427a8-0c7b-43a8-c421-c27264eda90f.htm)特性 (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadUInt16(String)](66238795-7948-2cf9-c639-d1bb3da126dd.htm) | 读取16位的无符号整型  Read 16-bit unsigned integer (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadUInt16(String, UInt16)](28888661-1fcc-f941-5746-9124eb539fe8.htm) | 读取16位的无符号整型数组  Read 16-bit unsigned integer array (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadUInt16Async(String)](1744a25b-c440-7d1e-dee5-606b0ff485bc.htm) | 异步读取16位的无符号整型  Asynchronously read 16-bit unsigned integer (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadUInt16Async(String, UInt16)](0748f071-d6a7-581b-ed7a-d73317b16762.htm) | 异步读取16位的无符号整型数组  Asynchronously read 16-bit unsigned integer array (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadUInt32(String)](b252ad01-4035-bc3e-e297-7fd97c3ac7af.htm) | 读取32位的无符号整型  Read 32-bit unsigned integer (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadUInt32(String, UInt16)](64d8381b-05b6-a990-c03e-2c8a1b7cfb7f.htm) | 读取32位的无符号整型数组  Read 32-bit unsigned integer array (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadUInt32Async(String)](a5b848e0-82bd-2bb5-2dd9-6634063e7b48.htm) | 异步读取32位的无符号整型  Asynchronously read 32-bit unsigned integer (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadUInt32Async(String, UInt16)](59ba0cbe-08f8-5a30-6029-0a9f718c7d76.htm) | 异步读取32位的无符号整型数组  Asynchronously read 32-bit unsigned integer array (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadUInt64(String)](529c82ba-1d5e-f33b-2eb9-c862799f952f.htm) | 读取64位的无符号整型  Read 64-bit unsigned integer (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadUInt64(String, UInt16)](2af999cc-c8d4-e25a-68f1-d621e5d84af7.htm) | 读取64位的无符号整型的数组  Read 64-bit unsigned integer array (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadUInt64Async(String)](0c521a11-248c-ba7a-fc27-0e5a645611b5.htm) | 异步读取64位的无符号整型  Asynchronously read 64-bit unsigned integer (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadUInt64Async(String, UInt16)](b510070f-beb5-959b-7c6b-20f212a68dd1.htm) | 异步读取64位的无符号整型的数组  Asynchronously read 64-bit unsigned integer array (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [ReadWrite](6296fde9-8c50-c741-29da-00c5d0054a3b.htm) | 使用0x17功能码来实现同时写入并读取数据的操作，使用一条报文来实现，需要指定读取的地址，长度，写入的地址，写入的数据信息，返回读取的结果数据。  Use 0x17 function code to write and read data at the same time, and use a message to implement it, you need to specify the read address, length, written address, written data information, and return the read result data. |
| 公共方法 | [RegisteredAddressMapping](53b0eecc-8513-81c0-4a8b-a715d134e29c.htm) | 注册一个新的地址映射关系，注册地址映射关系后，就可以使用新的地址来读写Modbus数据了，通常用于其他的支持Modbus协议的PLC。  After registering a new address mapping, you can use the new address to read and write Modbus data, which is usually used for other PLCs that support the Modbus protocol. |
| 公共方法 | [TranslateToModbusAddress](68338057-dac7-e74e-f841-de013cea8da3.htm) | 将当前的地址信息转换成Modbus格式的地址，如果转换失败，返回失败的消息。默认不进行任何的转换。  Convert the current address information into a Modbus format address. If the conversion fails, a failure message will be returned. No conversion is performed by default. |
| 公共方法 | [Wait(String, Boolean, Int32, Int32)](c8c01cf3-bf99-761c-23d4-82dd54081952.htm) | 等待指定地址的Boolean值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the Boolean value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [Wait(String, Int16, Int32, Int32)](f8a820b0-cf90-668c-6fba-524bd0659042.htm) | 等待指定地址的Int16值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the Int16 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [Wait(String, Int32, Int32, Int32)](96ad6909-b660-2894-9c17-a7f4f4d4a2bd.htm) | 等待指定地址的Int32值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the Int32 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [Wait(String, Int64, Int32, Int32)](28ba4678-2aeb-7502-95e8-d6b855bbe971.htm) | 等待指定地址的Int64值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the Int64 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [Wait(String, UInt16, Int32, Int32)](44af3173-f7f4-04f5-753b-7855581a171a.htm) | 等待指定地址的UInt16值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the UInt16 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [Wait(String, UInt32, Int32, Int32)](0e877e55-9b76-f662-6841-9102a9e9b9b8.htm) | 等待指定地址的UInt32值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the UInt32 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [Wait(String, UInt64, Int32, Int32)](6801831d-c349-800d-e3d3-34d2fde09e0f.htm) | 等待指定地址的UInt64值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the UInt64 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [WaitAsync(String, Boolean, Int32, Int32)](b51daff5-9c09-c614-9ea3-f14e907a36b8.htm) | 等待指定地址的Boolean值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the Boolean value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [WaitAsync(String, Int16, Int32, Int32)](cfd8af62-371b-84e9-5695-87539fbb204b.htm) | 等待指定地址的Int16值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the Int16 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [WaitAsync(String, Int32, Int32, Int32)](639f2d1b-9a41-0cf2-0335-7f9641fbb0c8.htm) | 等待指定地址的Int32值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the Int32 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [WaitAsync(String, Int64, Int32, Int32)](4355bcb4-e923-9226-ac82-55dd1325f34f.htm) | 等待指定地址的Int64值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the Int64 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [WaitAsync(String, UInt16, Int32, Int32)](291b6f29-9cbf-6fc3-4828-226f412e9cc6.htm) | 等待指定地址的UInt16值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the UInt16 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [WaitAsync(String, UInt32, Int32, Int32)](e1009c58-9f02-660e-225a-b6861eb8781f.htm) | 等待指定地址的UInt32值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the UInt32 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [WaitAsync(String, UInt64, Int32, Int32)](4236aa40-6c9d-2a0c-8470-c55dab497074.htm) | 等待指定地址的UInt64值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the UInt64 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [Write(String, Boolean)](4d742738-6864-465c-5382-bd49032eae37.htm) | 写入单个的Boolean数据，返回是否成功  Write a single Boolean data, and return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [Write(String, Boolean)](2a99122c-1496-d766-f0aa-44a18918fcd3.htm) | 批量写入Boolean数组数据，返回是否成功  Batch write Boolean array data, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [Write(String, Byte)](b8d8d94e-300d-c7b7-2e80-6859c62479fb.htm) | 写入原始的byte数组数据到指定的地址，返回是否写入成功  Write the original byte array data to the specified address, and return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [Write(String, Double)](34e1c967-7214-bb8e-3daf-fdb93fc1db59.htm) | 写入double数据，返回是否成功  Write double data, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [Write(String, Double)](f2a8cc81-ee85-4d3a-1eaf-af33b15de5e6.htm) | 写入double数组，返回是否成功  Write double array, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [Write(String, Int16)](53b27553-763c-2dd9-b8a1-df2ebf68ab2c.htm) | 写入short数据，返回是否成功  Write short data, returns whether success (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [Write(String, Int16)](0b8fff2d-0da9-d73a-bc11-866a501e6112.htm) | 写入short数组，返回是否成功  Write short array, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [Write(String, Int32)](281ce6f6-5c53-da62-9a09-9d439390a972.htm) | 写入int数据，返回是否成功  Write int data, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [Write(String, Int32)](496e02c2-62f2-4f72-aff1-1d691cc87726.htm) | 写入int[]数组，返回是否成功  Write int array, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [Write(String, Int64)](7ec919a0-1d6b-b7d7-1f2f-eb1644beb912.htm) | 写入long数据，返回是否成功  Write long data, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [Write(String, Int64)](026100a5-d52e-b640-f239-a140c903201d.htm) | 写入long数组，返回是否成功  Write long array, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [Write(String, Single)](7c72c919-367f-9bc6-567c-550813459297.htm) | 写入float数据，返回是否成功  Write float data, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [Write(String, Single)](195c858e-1e2d-863a-8a2c-dc49cb2063d3.htm) | 写入float数组，返回是否成功  Write float array, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [Write(String, String)](fc988c8f-b2d6-18f9-3274-9b57d4a04e1c.htm) | 写入字符串信息，编码为ASCII  Write string information, encoded as ASCII (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [Write(String, UInt16)](7084bc37-6bda-e506-4118-74aff979efff.htm) | 写入ushort数据，返回是否成功  Write ushort data, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [Write(String, UInt16)](9c21745f-7c5f-cae3-7f2a-2bcf260e54a9.htm) | 写入ushort数组，返回是否成功  Write ushort array, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [Write(String, UInt32)](ed79fea9-bc6e-3ef3-5294-134b872141cd.htm) | 写入uint数据，返回是否成功  Write uint data, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [Write(String, UInt32)](3ae38bac-66cc-5cdf-be68-483362b17ac3.htm) | 写入uint[]数组，返回是否成功  Write uint array, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [Write(String, UInt64)](3927095f-0705-e823-bb24-74c9363fab36.htm) | 写入ulong数据，返回是否成功  Write ulong data, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [Write(String, UInt64)](2eb4a3ce-c350-9949-405b-ccb6b5d9786e.htm) | 写入ulong数组，返回是否成功  Write ulong array, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [Write(String, String, Int32)](1ed07021-8f14-743a-071c-c348025634f5.htm) | 写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为ASCII  Write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, it is filled with 0 and the encoding is ASCII. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [Write(String, String, Encoding)](c66a9dfb-840a-3ac4-96ef-48e23fb85ec1.htm) | 写入字符串信息，需要指定的编码信息  Write string information, need to specify the encoding information (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [Write(String, String, Int32, Encoding)](61396f95-43cd-34fb-9055-37676c4e436d.htm) | 写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为指定的编码信息  Write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, then the operation is complemented with 0 , you should specified the encoding information (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteT(T)](ec80557a-6d74-6030-7e09-57c0c6a63cc4.htm) | 写入支持Hsl特性的数据，返回是否写入成功，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Write data that supports the Hsl attribute, and return whether the write was successful. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [WriteAsync(String, Boolean)](559924da-b298-9c81-794f-83d53aafc602.htm) | 异步批量写入Boolean数组数据，返回是否成功  Asynchronously batch write Boolean array data, return success (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [WriteAsync(String, Boolean)](cf8f8e56-df3d-607b-1c52-c1a76378db39.htm) | 异步批量写入Boolean数组数据，返回是否成功  Asynchronously batch write Boolean array data, return success (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Byte)](1a2e83c8-7782-847a-5d89-08d5556ac4ad.htm) | 异步写入原始的byte数组数据到指定的地址，返回是否写入成功  Asynchronously writes the original byte array data to the specified address, and returns whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Double)](ca502ab0-a05c-9dc7-c8b0-9afd3a24cb63.htm) | 异步写入double数据，返回是否成功  Asynchronously write double data, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Double)](33c58466-53fe-e221-1d73-6650c434e2b1.htm) | 异步写入double数组，返回是否成功  Asynchronously write double array, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int16)](3e923476-059a-0de9-4fc5-6df5dc7a93a5.htm) | 异步写入short数据，返回是否成功  Asynchronously write short data, returns whether success (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int16)](8955c5b3-7709-57ce-51fa-a58369f1147f.htm) | 异步写入short数组，返回是否成功  Asynchronously write short array, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int32)](448b905c-dfee-7f97-4dd9-8caf3539eef4.htm) | 异步写入int数据，返回是否成功  Asynchronously write int data, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int32)](cd5f66ad-6ffc-5394-4554-43abb07d977a.htm) | 异步写入int[]数组，返回是否成功  Asynchronously write int array, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int64)](972d6a25-bb1a-3c92-c12a-4b6942fc026f.htm) | 异步写入long数据，返回是否成功  Asynchronously write long data, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int64)](89e51562-f954-65b5-0c20-d1602d8e78a8.htm) | 异步写入long数组，返回是否成功  Asynchronously write long array, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Single)](213794a0-4dc8-36b4-3b3f-d1ee74e35155.htm) | 异步写入float数据，返回是否成功  Asynchronously write float data, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Single)](dca9a471-a3cb-3242-b3fe-edeb04d4f786.htm) | 异步写入float数组，返回是否成功  Asynchronously write float array, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, String)](265dfc0f-8630-cc0a-4a25-cd31d50f6056.htm) | 异步写入字符串信息，编码为ASCII  Asynchronously write string information, encoded as ASCII (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt16)](ee35584e-d8c8-720e-7a05-24a2d3011021.htm) | 异步写入ushort数据，返回是否成功  Asynchronously write ushort data, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt16)](9df3d402-f0c7-5888-65d0-f527a6ada1b2.htm) | 异步写入ushort数组，返回是否成功  Asynchronously write ushort array, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt32)](76ac4f61-3a88-384c-fae1-857066a74ac0.htm) | 异步写入uint数据，返回是否成功  Asynchronously write uint data, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt32)](aae8db84-1d63-31a3-e84c-f03a5f784ff0.htm) | 异步写入uint[]数组，返回是否成功  Asynchronously write uint array, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt64)](81bbb9e0-b222-067d-7dd9-a879277859c4.htm) | 异步写入ulong数据，返回是否成功  Asynchronously write ulong data, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt64)](f56c270e-dc8d-9adc-ea60-36d47e2c6e50.htm) | 异步写入ulong数组，返回是否成功  Asynchronously write ulong array, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [WriteAsync(String, String, Int32)](c9efea95-68b3-ab09-5176-d388576844e8.htm) | 异步写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为ASCII  Asynchronously write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, it is filled with 0 and the encoding is ASCII. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, String, Encoding)](cda78ec9-ab0c-b6f1-db38-0e51e6aa06d0.htm) | 异步写入字符串信息，需要指定的编码信息  Asynchronously write string information, need to specify the encoding information (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [WriteAsync(String, String, Int32, Encoding)](c077e3ad-d170-aea3-0058-2888f2b27126.htm) | 异步写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为指定的编码信息  Asynchronously write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, then the operation is complemented with 0 , you should specified the encoding information (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteAsyncT(T)](da0f692b-4d73-0520-6e17-7890b6ef51cb.htm) | 异步写入支持Hsl特性的数据，返回是否写入成功，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Asynchronously write data that supports the Hsl attribute, and return whether the write was successful. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteCustomerT](466c49b4-95ac-148b-c18d-a419c5f6dab3.htm) | 写入自定义类型的数据，该类型必须继承自IDataTransfer接口  Write data of a custom type, which must inherit from the IDataTransfer interface (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteCustomerAsyncT](82aa8f35-3c06-d4f6-e8c5-affb6ab6ce41.htm) | 写入自定义类型的数据，该类型必须继承自IDataTransfer接口  Write data of a custom type, which must inherit from the IDataTransfer interface (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [WriteFile](63cb39f3-f1de-824f-1671-5839e496c80f.htm) | 使用0x15功能码写入文件的数据内容，需要指定文件编号，地址，数据，返回写入的结果数据。写入的数据字节长度不能大于248个字节，如果想写入更多的字节，就需要切片写入操作 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[IModbus 接口](1d3ec4c8-6dd5-984d-e7ac-c107336985a5.htm)

[HslCommunication.ModBus 命名空间](da9ac567-746e-47e5-70e1-b44577d09ce4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadFile 方法 

[原文連結](http://api.hslcommunication.cn/html/97062c02-4031-5a46-9d12-744755610942.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.ModBus](../html/da9ac567-746e-47e5-70e1-b44577d09ce4.htm "HslCommunication.ModBus")

[IModbus 接口](../html/1d3ec4c8-6dd5-984d-e7ac-c107336985a5.htm "IModbus 接口")

[IModbus 方法](../html/6a76dccd-76d8-56d4-6b83-0eb4f03215f8.htm "IModbus 方法")

[ReadFile 方法](../html/97062c02-4031-5a46-9d12-744755610942.htm "ReadFile 方法 ")

[ReadWrite 方法](../html/6296fde9-8c50-c741-29da-00c5d0054a3b.htm "ReadWrite 方法 ")

[RegisteredAddressMapping 方法](../html/53b0eecc-8513-81c0-4a8b-a715d134e29c.htm "RegisteredAddressMapping 方法 ")

[TranslateToModbusAddress 方法](../html/68338057-dac7-e74e-f841-de013cea8da3.htm "TranslateToModbusAddress 方法 ")

[WriteFile 方法](../html/63cb39f3-f1de-824f-1671-5839e496c80f.htm "WriteFile 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IModbusReadFile 方法 |

使用0x14功能码读取文件的数据内容，需要指定文件编号，地址，长度，返回读取的结果数据，读取的长度为任意长度，内部自动根据实际情况来切割操作。

**命名空间：**
 [HslCommunication.ModBus](da9ac567-746e-47e5-70e1-b44577d09ce4.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
OperateResult<byte[]> ReadFile(
	ushort fileNumber,
	ushort address,
	ushort length
)
```

```
Function ReadFile ( 
	fileNumber As UShort,
	address As UShort,
	length As UShort
) As OperateResult(Of Byte())
```

```
OperateResult<array<unsigned char>^>^ ReadFile(
	unsigned short fileNumber, 
	unsigned short address, 
	unsigned short length
)
```

```
abstract ReadFile : 
        fileNumber : uint16 * 
        address : uint16 * 
        length : uint16 -> OperateResult<byte[]> 
```

#### 参数

fileNumber
:   类型：SystemUInt16  
    文件号

address
:   类型：SystemUInt16  
    起始地址

length
:   类型：SystemUInt16  
    读取长度

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
返回结果

![](../icons/SectionExpanded.png)参见

#### 引用

[IModbus 接口](1d3ec4c8-6dd5-984d-e7ac-c107336985a5.htm)

[HslCommunication.ModBus 命名空间](da9ac567-746e-47e5-70e1-b44577d09ce4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadWrite 方法 

[原文連結](http://api.hslcommunication.cn/html/6296fde9-8c50-c741-29da-00c5d0054a3b.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.ModBus](../html/da9ac567-746e-47e5-70e1-b44577d09ce4.htm "HslCommunication.ModBus")

[IModbus 接口](../html/1d3ec4c8-6dd5-984d-e7ac-c107336985a5.htm "IModbus 接口")

[IModbus 方法](../html/6a76dccd-76d8-56d4-6b83-0eb4f03215f8.htm "IModbus 方法")

[ReadFile 方法](../html/97062c02-4031-5a46-9d12-744755610942.htm "ReadFile 方法 ")

[ReadWrite 方法](../html/6296fde9-8c50-c741-29da-00c5d0054a3b.htm "ReadWrite 方法 ")

[RegisteredAddressMapping 方法](../html/53b0eecc-8513-81c0-4a8b-a715d134e29c.htm "RegisteredAddressMapping 方法 ")

[TranslateToModbusAddress 方法](../html/68338057-dac7-e74e-f841-de013cea8da3.htm "TranslateToModbusAddress 方法 ")

[WriteFile 方法](../html/63cb39f3-f1de-824f-1671-5839e496c80f.htm "WriteFile 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IModbusReadWrite 方法 |

使用0x17功能码来实现同时写入并读取数据的操作，使用一条报文来实现，需要指定读取的地址，长度，写入的地址，写入的数据信息，返回读取的结果数据。  
Use 0x17 function code to write and read data at the same time, and use a message to implement it,
you need to specify the read address, length, written address, written data information, and return the read result data.

**命名空间：**
 [HslCommunication.ModBus](da9ac567-746e-47e5-70e1-b44577d09ce4.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
OperateResult<byte[]> ReadWrite(
	string readAddress,
	ushort length,
	string writeAddress,
	byte[] value
)
```

```
Function ReadWrite ( 
	readAddress As String,
	length As UShort,
	writeAddress As String,
	value As Byte()
) As OperateResult(Of Byte())
```

```
OperateResult<array<unsigned char>^>^ ReadWrite(
	String^ readAddress, 
	unsigned short length, 
	String^ writeAddress, 
	array<unsigned char>^ value
)
```

```
abstract ReadWrite : 
        readAddress : string * 
        length : uint16 * 
        writeAddress : string * 
        value : byte[] -> OperateResult<byte[]> 
```

#### 参数

readAddress
:   类型：SystemString  
    读取的地址信息

length
:   类型：SystemUInt16  
    读取的长度信息

writeAddress
:   类型：SystemString  
    写入的地址信息

value
:   类型：SystemByte  
    写入的字节数据信息

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
读取的结果对象

![](../icons/SectionExpanded.png)参见

#### 引用

[IModbus 接口](1d3ec4c8-6dd5-984d-e7ac-c107336985a5.htm)

[HslCommunication.ModBus 命名空间](da9ac567-746e-47e5-70e1-b44577d09ce4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## RegisteredAddressMapping 方法 

[原文連結](http://api.hslcommunication.cn/html/53b0eecc-8513-81c0-4a8b-a715d134e29c.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.ModBus](../html/da9ac567-746e-47e5-70e1-b44577d09ce4.htm "HslCommunication.ModBus")

[IModbus 接口](../html/1d3ec4c8-6dd5-984d-e7ac-c107336985a5.htm "IModbus 接口")

[IModbus 方法](../html/6a76dccd-76d8-56d4-6b83-0eb4f03215f8.htm "IModbus 方法")

[ReadFile 方法](../html/97062c02-4031-5a46-9d12-744755610942.htm "ReadFile 方法 ")

[ReadWrite 方法](../html/6296fde9-8c50-c741-29da-00c5d0054a3b.htm "ReadWrite 方法 ")

[RegisteredAddressMapping 方法](../html/53b0eecc-8513-81c0-4a8b-a715d134e29c.htm "RegisteredAddressMapping 方法 ")

[TranslateToModbusAddress 方法](../html/68338057-dac7-e74e-f841-de013cea8da3.htm "TranslateToModbusAddress 方法 ")

[WriteFile 方法](../html/63cb39f3-f1de-824f-1671-5839e496c80f.htm "WriteFile 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IModbusRegisteredAddressMapping 方法 |

注册一个新的地址映射关系，注册地址映射关系后，就可以使用新的地址来读写Modbus数据了，通常用于其他的支持Modbus协议的PLC。  
After registering a new address mapping, you can use the new address to read and write Modbus data, which is usually used for other PLCs that support the Modbus protocol.

**命名空间：**
 [HslCommunication.ModBus](da9ac567-746e-47e5-70e1-b44577d09ce4.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
void RegisteredAddressMapping(
	Func<string, byte, OperateResult<string>> mapping
)
```

```
Sub RegisteredAddressMapping ( 
	mapping As Func(Of String, Byte, OperateResult(Of String))
)
```

```
void RegisteredAddressMapping(
	Func<String^, unsigned char, OperateResult<String^>^>^ mapping
)
```

```
abstract RegisteredAddressMapping : 
        mapping : Func<string, byte, OperateResult<string>> -> unit 
```

#### 参数

mapping
:   类型：SystemFuncString, Byte, [OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)String  
    地址映射关系信息

![](../icons/SectionExpanded.png)参见

#### 引用

[IModbus 接口](1d3ec4c8-6dd5-984d-e7ac-c107336985a5.htm)

[HslCommunication.ModBus 命名空间](da9ac567-746e-47e5-70e1-b44577d09ce4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## TranslateToModbusAddress 方法 

[原文連結](http://api.hslcommunication.cn/html/68338057-dac7-e74e-f841-de013cea8da3.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.ModBus](../html/da9ac567-746e-47e5-70e1-b44577d09ce4.htm "HslCommunication.ModBus")

[IModbus 接口](../html/1d3ec4c8-6dd5-984d-e7ac-c107336985a5.htm "IModbus 接口")

[IModbus 方法](../html/6a76dccd-76d8-56d4-6b83-0eb4f03215f8.htm "IModbus 方法")

[ReadFile 方法](../html/97062c02-4031-5a46-9d12-744755610942.htm "ReadFile 方法 ")

[ReadWrite 方法](../html/6296fde9-8c50-c741-29da-00c5d0054a3b.htm "ReadWrite 方法 ")

[RegisteredAddressMapping 方法](../html/53b0eecc-8513-81c0-4a8b-a715d134e29c.htm "RegisteredAddressMapping 方法 ")

[TranslateToModbusAddress 方法](../html/68338057-dac7-e74e-f841-de013cea8da3.htm "TranslateToModbusAddress 方法 ")

[WriteFile 方法](../html/63cb39f3-f1de-824f-1671-5839e496c80f.htm "WriteFile 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IModbusTranslateToModbusAddress 方法 |

将当前的地址信息转换成Modbus格式的地址，如果转换失败，返回失败的消息。默认不进行任何的转换。  
Convert the current address information into a Modbus format address. If the conversion fails, a failure message will be returned. No conversion is performed by default.

**命名空间：**
 [HslCommunication.ModBus](da9ac567-746e-47e5-70e1-b44577d09ce4.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
OperateResult<string> TranslateToModbusAddress(
	string address,
	byte modbusCode
)
```

```
Function TranslateToModbusAddress ( 
	address As String,
	modbusCode As Byte
) As OperateResult(Of String)
```

```
OperateResult<String^>^ TranslateToModbusAddress(
	String^ address, 
	unsigned char modbusCode
)
```

```
abstract TranslateToModbusAddress : 
        address : string * 
        modbusCode : byte -> OperateResult<string> 
```

#### 参数

address
:   类型：SystemString  
    传入的地址

modbusCode
:   类型：SystemByte  
    Modbus的功能码

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)String  
转换之后Modbus的地址

![](../icons/SectionExpanded.png)参见

#### 引用

[IModbus 接口](1d3ec4c8-6dd5-984d-e7ac-c107336985a5.htm)

[HslCommunication.ModBus 命名空间](da9ac567-746e-47e5-70e1-b44577d09ce4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## WriteFile 方法 

[原文連結](http://api.hslcommunication.cn/html/63cb39f3-f1de-824f-1671-5839e496c80f.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.ModBus](../html/da9ac567-746e-47e5-70e1-b44577d09ce4.htm "HslCommunication.ModBus")

[IModbus 接口](../html/1d3ec4c8-6dd5-984d-e7ac-c107336985a5.htm "IModbus 接口")

[IModbus 方法](../html/6a76dccd-76d8-56d4-6b83-0eb4f03215f8.htm "IModbus 方法")

[ReadFile 方法](../html/97062c02-4031-5a46-9d12-744755610942.htm "ReadFile 方法 ")

[ReadWrite 方法](../html/6296fde9-8c50-c741-29da-00c5d0054a3b.htm "ReadWrite 方法 ")

[RegisteredAddressMapping 方法](../html/53b0eecc-8513-81c0-4a8b-a715d134e29c.htm "RegisteredAddressMapping 方法 ")

[TranslateToModbusAddress 方法](../html/68338057-dac7-e74e-f841-de013cea8da3.htm "TranslateToModbusAddress 方法 ")

[WriteFile 方法](../html/63cb39f3-f1de-824f-1671-5839e496c80f.htm "WriteFile 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IModbusWriteFile 方法 |

使用0x15功能码写入文件的数据内容，需要指定文件编号，地址，数据，返回写入的结果数据。写入的数据字节长度不能大于248个字节，如果想写入更多的字节，就需要切片写入操作

**命名空间：**
 [HslCommunication.ModBus](da9ac567-746e-47e5-70e1-b44577d09ce4.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
OperateResult WriteFile(
	ushort fileNumber,
	ushort address,
	byte[] data
)
```

```
Function WriteFile ( 
	fileNumber As UShort,
	address As UShort,
	data As Byte()
) As OperateResult
```

```
OperateResult^ WriteFile(
	unsigned short fileNumber, 
	unsigned short address, 
	array<unsigned char>^ data
)
```

```
abstract WriteFile : 
        fileNumber : uint16 * 
        address : uint16 * 
        data : byte[] -> OperateResult 
```

#### 参数

fileNumber
:   类型：SystemUInt16  
    文件号

address
:   类型：SystemUInt16  
    起始地址

data
:   类型：SystemByte  
    写入的数据信息

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否写入成功

![](../icons/SectionExpanded.png)参见

#### 引用

[IModbus 接口](1d3ec4c8-6dd5-984d-e7ac-c107336985a5.htm)

[HslCommunication.ModBus 命名空间](da9ac567-746e-47e5-70e1-b44577d09ce4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ModbusAscii 类

[原文連結](http://api.hslcommunication.cn/html/873a8425-d84d-1a6f-ed5c-79e8d056cc76.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.ModBus](../html/da9ac567-746e-47e5-70e1-b44577d09ce4.htm "HslCommunication.ModBus")

[ModbusAscii 类](../html/873a8425-d84d-1a6f-ed5c-79e8d056cc76.htm "ModbusAscii 类")

[ModbusAscii 构造函数](../html/223428e7-bc23-5a0f-0fa6-ac355d422c71.htm "ModbusAscii 构造函数 ")

[ModbusAscii 属性](../html/fde6aaa6-18a7-1888-9bf8-dd886a14100b.htm "ModbusAscii 属性")

[ModbusAscii 方法](../html/7819f1d9-ca1b-2c0b-6ae3-9a83dfa03b21.htm "ModbusAscii 方法")

[ModbusAscii 字段](../html/62529de8-1a77-85ed-6b9c-9167bdfd5459.htm "ModbusAscii 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ModbusAscii 类 |

Modbus-Ascii通讯协议的类库，基于rtu类库完善过来，支持标准的功能码，也支持扩展的功能码实现，地址采用富文本的形式，详细见备注说明  
The client communication class of Modbus-Ascii protocol is convenient for data interaction with the server. It supports standard function codes and also supports extended function codes.
The address is in rich text. For details, see the remarks.

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  [HslCommunication.Core.NetBinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)  
    [HslCommunication.Core.DeviceDeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)  
      [HslCommunication.Core.DeviceDeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)  
        [HslCommunication.ModBusModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)  
          HslCommunication.ModBusModbusAscii  
            [HslCommunication.Profinet.DeltaDeltaSerialAscii](9c7335b9-7dac-db72-d707-7402bf190719.htm)

**命名空间：**
 [HslCommunication.ModBus](da9ac567-746e-47e5-70e1-b44577d09ce4.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class ModbusAscii : ModbusRtu
```

```
Public Class ModbusAscii
	Inherits ModbusRtu
```

```
public ref class ModbusAscii : public ModbusRtu
```

```
type ModbusAscii =  
    class
        inherit ModbusRtu
    end
```

ModbusAscii 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [ModbusAscii](29c891cf-fb6d-b443-ff97-12d8c955c41f.htm) | 实例化一个Modbus-ascii协议的客户端对象  Instantiate a client object of the Modbus-ascii protocol |
| 公共方法 | [ModbusAscii(Byte)](b6fbcbed-323d-9f61-fd2b-9b8dac246caa.htm) | 指定Modbus从站的站号来初始化  Specify the station number of the Modbus slave to initialize |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [AddressStartWithZero](4d17309a-509d-62c1-9119-3a5735c8315d.htm) | 获取或设置起始的地址是否从0开始，默认为True  Gets or sets whether the starting address starts from 0. The default is True (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共属性 | [BaudRate](5554f563-5e29-699c-435b-c8e7ad5fa0e4.htm) | 当前连接串口信息的波特率  Baud rate of current connection serial port information (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 公共属性 | [BroadcastStation](27196ecf-f9c7-30ac-d439-bf70d26dc328.htm) | 获取或是设置当前广播模式对应的站号，广播模式意味着不接收设备方的数据返回操作，默认为 -1，表示不使用广播模式。  Gets or sets the station number corresponding to the current broadcast mode. Broadcast mode means that the data return operation of the device is not received. The default value is -1, indicating that broadcast mode is not used. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共属性代码示例 | [ByteTransform](f35118c9-3691-f752-ae6f-2c5549714faa.htm) | 当前的数据变换机制，当你需要从字节数据转换类型数据的时候需要。  The current data transformation mechanism is required when you need to convert type data from byte data. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共属性 | [CommunicationPipe](a4fcb477-64d2-8da7-c712-687f0261755b.htm) | 获取或设置当前的管道信息，管道类型为[CommunicationPipe](ad8e903b-d289-21b0-077e-15ede8005221.htm)的继承类，内置了[PipeTcpNet](26b94bd5-2695-36e7-6378-c3d41846be8f.htm)管道，[PipeUdpNet](5a0b810e-187e-826b-0a13-34fd5576aa0e.htm)管道，[PipeSerialPort](9d924607-3dfc-3916-b532-c5dec48fe347.htm)管道等  Get or set the current pipeline information, the pipeline type is [CommunicationPipe](ad8e903b-d289-21b0-077e-15ede8005221.htm) inheritance class, [PipeTcpNet](26b94bd5-2695-36e7-6378-c3d41846be8f.htm) pipeline, [PipeUdpNet](5a0b810e-187e-826b-0a13-34fd5576aa0e.htm) pipeline, [PipeSerialPort](9d924607-3dfc-3916-b532-c5dec48fe347.htm) pipeline, etc (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 公共属性 | [ConnectionId](ce1efb4b-fbc9-4683-163e-0a8e71e99d6b.htm) | 当前连接的唯一ID号，默认为长度20的guid码加随机数组成，方便列表管理，也可以自己指定  The unique ID number of the current connection. The default is a 20-digit guid code plus a random number. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [Crc16CheckEnable](8b8b3de6-eaf7-0484-604f-657dc210ab61.htm) | 获取或设置是否启用CRC16校验码的检查功能，默认启用，如果需要忽略检查CRC16，则设置为 false 即可。  Gets or sets whether to enable the check function of CRC16 check code. It is enabled by default. If you need to ignore the check of CRC16, you can set it to false. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共属性 | [DataFormat](eeb434e5-c320-09f1-6cf4-ff4798c2505a.htm) | 获取或设置数据解析的格式，可选ABCD, BADC，CDAB，DCBA格式，对int,uint,float,double,long,ulong类型有作用  Get or set the format of the data analysis, optional ABCD, BADC, CDAB, DCBA format, effective for int, uint, float, double, long, ulong type (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共属性 | [DisableFunctionCode06](096f15a9-323c-69db-65a1-be3723d914d1.htm) | 获取或设置是否禁用功能码05的写入操作，设置为true时，表示不允许使用功能码 06 进行写入操作。写入单个的字数据也使用 0x10 功能码  Get or set whether to disable the write operation of function code 06. When set to true, it indicates that the write operation of function code 05 is not allowed. Writing individual character data also uses the 0x10 function code (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共属性 | [EnableWriteMaskCode](85372ceb-e693-1067-0285-581b163f8da4.htm) | 获取或设置当前掩码写入的功能码是否激活状态，设置为 false 时，再执行写入位时，会通过读字，修改位，写字的方式来间接实现。  When the function code is set to false, and then the write bit is executed, it will be indirectly implemented by reading, modifying the bit, and writing the word. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共属性 | [IsClearCacheBeforeRead](81dc971f-1f86-8535-ae9a-586b6590a584.htm) | 是否在发送数据前清空缓冲数据，默认是false  Whether to empty the buffer before sending data, the default is false (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 公共属性 | [IsStringReverse](c063a2e6-4c63-cb3f-8956-9164f3c853f0.htm) | 字符串数据是否按照字来反转，默认为False  Whether the string data is reversed according to words. The default is False. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共属性代码示例 | [LogNet](8c56f246-88ac-afb8-8c06-972fafe4f78a.htm) | 组件的日志工具，支持日志记录，只要实例化后，当前网络的基本信息，就以[DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm)等级进行输出  The component's logging tool supports logging. As long as the instantiation of the basic network information, the output will be output at [DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm) (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [PortName](f2b8e4e3-7079-5240-4ff1-14c6c9b376e6.htm) | 当前连接串口信息的端口号名称  The port name of the current connection serial port information (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 公共属性 | [ReceiveEmptyDataCount](f800fd56-7d49-87c4-985b-103a38330078.htm) | 获取或设置连续接收空的数据次数，在数据接收完成时有效，每个单位消耗的时间为[SleepTime](1e3a83ef-b6fc-bc8c-8950-ea3f6ae5c19b.htm)。  Obtain or set the number of consecutive times to receive empty data, which is valid when the data is received, and the time consumed by each unit is [SleepTime](1e3a83ef-b6fc-bc8c-8950-ea3f6ae5c19b.htm) (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 公共属性代码示例 | [ReceiveTimeOut](250560a7-332b-ce1b-5eb0-e3abe4ee6012.htm) | 获取或设置接收服务器反馈的时间，如果为负数，则不接收反馈   Gets or sets the time to receive server feedback, and if it is a negative number, does not receive feedback (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [RtsEnable](dd4057ef-88fe-f7d0-d2c4-f6e15488fe78.htm) | 获取或设置一个值，该值指示在串行通信中是否启用请求发送 (RTS) 信号。  Gets or sets a value indicating whether the request sending (RTS) signal is enabled in serial communication. (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 公共属性 | [SendBeforeHex](72aa1de4-116f-abfe-4adb-814a0be0e560.htm) | 获取或设置在发送通信报文前追加发送的字节信息，HEX格式，通常用于lora组网时，需要携带 00 00 00 02 四个字节的站地址功能。  Obtain or set the byte information sent before sending communication packets, HEX format, usually used for LORA networking, you need to carry 00 00 00 02 four-byte station address function. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [SleepTime](86a5fbb9-3456-6beb-7113-5bbad5eabbc5.htm) | 获取或设置在正式接收对方返回数据前的时候，需要休息的时间，当设置为0的时候，不需要休息。  Get or set the time required to rest before officially receiving the data from the other party. When it is set to 0, no rest is required. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [Station](063d55ef-6207-33e0-9073-a725af469563.htm) | 获取或者重新修改服务器的默认站号信息，当然，你可以再读写的时候动态指定，例如地址 "s=2;100", 更详细的例子可以参考DEMO界面上的地址示例。  Get or re-modify the server's default station number information, of course, you can dynamically specify when reading and writing, such as the address "s=2; 100", a more detailed example can refer to the DEMO interface address example. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共属性 | [StationCheckMatch](03c10d53-c7a2-9256-df12-d780b8efd32c.htm) | 获取或设置当前是否启用站号检查的功能，默认启用，读写数据时将进行站号确认操作，如果需要忽略站号，则设置为 false 即可。  Gets or sets whether the station ID check function is enabled. It is enabled by default. When reading or writing data, the station ID is confirmed. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 受保护的属性 | [WordLength](f1a505cc-e85c-6c8b-b189-58a1ab5753f9.htm) | 一个字单位的数据表示的地址长度，西门子为2，三菱，欧姆龙，modbusTcp就为1，AB PLC无效  The address length represented by one word of data, Siemens is 2, Mitsubishi, Omron, modbusTcp is 1, AB PLC is invalid (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共属性 | [WordReadBatchLength](75150ff2-55eb-9881-6d35-d916a737a46d.htm) | 当进行字读取的时候，商业授权用户支持分批读取，本属性表示分批的长度信息，默认为 120，可以根据实际情况来设置其他值。  When reading of word type, commercial authorization users support batch reading. This attribute represents the length information of the batch, with a default of 120. Other values can be set according to the actual situation. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
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
| 受保护的方法 | [GetNewNetMessage](88dbd66f-58bf-6e1b-c2ab-34c111364480.htm) | 获取一个新的消息对象的方法，需要在继承类里面进行重写  The method to get a new message object needs to be overridden in the inheritance class (重写 [ModbusRtuGetNewNetMessage](11f57be7-edd6-53dd-f91f-74c38a04578c.htm).) |
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
| 公共方法 | [PackCommandWithHeader](73a56f28-e750-b76d-3ebd-bd5e55ad9aae.htm) | 对当前的命令进行打包处理，通常是携带命令头内容，标记当前的命令的长度信息，需要进行重写，否则默认不打包  The current command is packaged, usually carrying the content of the command header, marking the length of the current command, and it needs to be rewritten, otherwise it is not packaged by default (重写 [ModbusRtuPackCommandWithHeader(Byte)](ad370105-3cc4-77f9-797d-0ca5edafd761.htm).) |
| 公共方法代码示例 | [Read(String, UInt16)](6571f657-29f0-1869-7e0c-4d4add79fc2b.htm) | 从Modbus服务器批量读取寄存器的信息，需要指定起始地址，读取长度，如果富文本地址不指定，默认使用的功能码是 0x03，如果需要使用04功能码，那么地址就写成 x=4;100  To read the register information from the Modbus server in batches, you need to specify the start address and read length. If the rich text address is not specified, the default function code is 0x03. If you need to use the 04 function code, the address is written as x = 4; 100 (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [ReadT](95abe064-5594-5ae7-a849-98c53f88fee0.htm) | 读取支持Hsl特性的数据内容，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明，或是官网的详细文档  Read the data content of the Hsl attribute. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadAsync(String, UInt16)](d8c14d37-56dc-2a31-1a19-3e2b5b65197d.htm) | 异步批量读取字节数组信息，需要指定地址和长度，返回原始的字节数组  Asynchronous batch read byte array information, need to specify the address and length, return the original byte array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadAsyncT](0522e740-4bf0-2e31-c4c5-e41384321b37.htm) | 异步读取支持Hsl特性的数据内容，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Asynchronously read the data content of the Hsl attribute. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadBool(String)](9b2faa2e-eaa9-7ace-c49f-7a0318f1df1a.htm) | 读取单个的Boolean数据信息  Read a single Boolean data message (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadBool(String, UInt16)](2d125640-18a5-f9a9-29b9-b121915c77ce.htm) | 批量读取线圈或是离散的数据信息，需要指定地址和长度，具体的结果取决于实现，如果富文本地址不指定，默认使用的功能码是 0x01  To read coils or discrete data in batches, you need to specify the address and length. The specific result depends on the implementation. If the rich text address is not specified, the default function code is 0x01. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [ReadBoolAsync(String)](cc37da1f-d673-fe67-78ee-4e465d29a2bd.htm) | 异步读取单个的Boolean数据信息  Asynchronously read a single Boolean data message (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadBoolAsync(String, UInt16)](a45d383d-0095-3ea7-32f8-e75af88d60f8.htm) | 异步批量读取Boolean数组信息，需要指定地址和长度，返回Boolean 数组  Asynchronously batch read Boolean array information, need to specify the address and length, return Boolean array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadCoil(String)](2d2adf84-ab2c-3b95-269f-3a8ddc659018.htm) | 读取线圈，需要指定起始地址，如果富文本地址不指定，默认使用的功能码是 0x01  To read the coil, you need to specify the start address. If the rich text address is not specified, the default function code is 0x01. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [ReadCoil(String, UInt16)](2869f250-557a-4888-5ebc-ea5e98028aed.htm) | 批量的读取线圈，需要指定起始地址，读取长度，如果富文本地址不指定，默认使用的功能码是 0x01  For batch reading coils, you need to specify the start address and read length. If the rich text address is not specified, the default function code is 0x01. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [ReadCoilAsync(String)](d5239965-a5cf-8e54-b4be-e4cb86b38cc0.htm) | 读取线圈，需要指定起始地址，如果富文本地址不指定，默认使用的功能码是 0x01  To read the coil, you need to specify the start address. If the rich text address is not specified, the default function code is 0x01. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [ReadCoilAsync(String, UInt16)](dfc142d3-3390-5eec-1128-6e19cac1ea12.htm) | 批量的读取线圈，需要指定起始地址，读取长度，如果富文本地址不指定，默认使用的功能码是 0x01  For batch reading coils, you need to specify the start address and read length. If the rich text address is not specified, the default function code is 0x01. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [ReadCustomerT(String)](ca069688-d8f0-7dbd-e962-86a10dce83cf.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，返回一个新的类型的实例对象。  To read a custom data type, you need to inherit from the IDataTransfer interface and return an instance object of a new type. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadCustomerT(String, T)](ca7dc1a1-c97d-689f-9819-4dae4af67206.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，传入一个实例，对这个实例进行赋值，并返回该实例的对象。  To read a custom data type, you need to inherit from the IDataTransfer interface, pass in an instance, assign a value to this instance, and return the object of the instance. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadCustomerAsyncT(String)](f3fcda2b-255d-aadd-567e-c7b3df3636d9.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，返回一个新的类型的实例对象。  To read a custom data type, you need to inherit from the IDataTransfer interface and return an instance object of a new type. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadCustomerAsyncT(String, T)](527c9fae-4743-6ebd-7b9a-5ef3a8b7c455.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，传入一个实例，对这个实例进行赋值，并返回该实例的对象。  To read a custom data type, you need to inherit from the IDataTransfer interface, pass in an instance, assign a value to this instance, and return the object of the instance. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadDiscrete(String)](cffe0847-c797-bc7f-71f1-f93092565ec6.htm) | 读取输入线圈，需要指定起始地址，如果富文本地址不指定，默认使用的功能码是 0x02  To read the input coil, you need to specify the start address. If the rich text address is not specified, the default function code is 0x02. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [ReadDiscrete(String, UInt16)](181cbf55-3a34-a988-9aa5-18432b9da97a.htm) | 批量的读取输入点，需要指定起始地址，读取长度，如果富文本地址不指定，默认使用的功能码是 0x02  To read input points in batches, you need to specify the start address and read length. If the rich text address is not specified, the default function code is 0x02 (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [ReadDiscreteAsync(String)](9d11737c-b5d9-5c24-09a1-0b70a87bdfc3.htm) | 读取输入线圈，需要指定起始地址，如果富文本地址不指定，默认使用的功能码是 0x02  To read the input coil, you need to specify the start address. If the rich text address is not specified, the default function code is 0x02. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [ReadDiscreteAsync(String, UInt16)](9c34359d-3d18-fdbf-d6c3-35134eec3b9d.htm) | 批量的读取输入点，需要指定起始地址，读取长度，如果富文本地址不指定，默认使用的功能码是 0x02  To read input points in batches, you need to specify the start address and read length. If the rich text address is not specified, the default function code is 0x02 (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [ReadDouble(String)](e7e31c43-71b4-e692-6bdc-ba9533e1d83d.htm) | 读取双浮点的数据  Read double floating point data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadDouble(String, UInt16)](9df438ac-8869-5002-7ace-bfb7c8f7aa6b.htm) | 读取双浮点数据的数组  Read double floating point data array (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [ReadDoubleAsync(String)](2cf895fa-1033-1def-f048-0e8caefaa71d.htm) | 异步读取双浮点的数据  Asynchronously read double floating point data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadDoubleAsync(String, UInt16)](e5c602ad-d3ce-a677-908b-03d72a4388cd.htm) | 异步读取双浮点数据的数组  Asynchronously read double floating point data array (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [ReadFile](702ddeff-3ba5-952d-abce-99a988a1d3cc.htm) | 使用0x14功能码读取文件的数据内容，需要指定文件编号，地址，长度，返回读取的结果数据，读取的长度为任意长度，内部自动根据实际情况来切割操作。 (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [ReadFloat(String)](bc2801ae-6600-e845-7c6c-522c9156eb0e.htm) | 读取单浮点数据  Read single floating point data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadFloat(String, UInt16)](909486dc-cc84-c34f-2de5-ec17b3d04182.htm) | 读取单浮点精度的数组  Read single floating point array (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [ReadFloatAsync(String)](89f743f8-4b3a-7085-1e49-a1d7adb10052.htm) | 异步读取单浮点数据  Asynchronously read single floating point data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadFloatAsync(String, UInt16)](eaeb9795-d78c-f884-fa54-c698c5e37404.htm) | 异步读取单浮点精度的数组  Asynchronously read single floating point array (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [ReadFromCoreServer(IEnumerableByte)](044cfdd9-ec76-b89f-2f2b-5f439025006c.htm) | 将多个数据报文按顺序发到设备，并从设备接收返回的数据内容，然后拼接成一个Byte[]信息，需要重写[UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm)方法才能返回正确的结果。  Send multiple data packets to the device in sequence, and receive the returned data content from the device, and then splicing them into a Byte[] message, you need to rewrite [UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm) method to return the correct result. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServer(Byte)](ff061b13-df60-6617-08c2-40b81cde38ea.htm) | 将Modbus报文数据发送到当前的通道中，并从通道中接收Modbus的报文，通道将根据当前连接自动获取，本方法是线程安全的。  Send Modbus message data to the current channel, and receive Modbus messages from the channel. The channel will automatically obtain it according to the current connection. This method is thread-safe. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServer(Byte, Boolean, Boolean)](b32ce8e8-a959-fcc2-6736-14b39fdb9a40.htm) | 将二进制的数据发送到管道中去，然后从管道里接收二进制的数据回来，并返回是否成功的结果对象。  Send binary data to the pipeline, and then receive binary data back from the pipeline, and return whether the success of the result object (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServer(CommunicationPipe, Byte, Boolean, Boolean)](ef36d56b-df09-f727-e8d3-316d37a13ee9.htm) | 使用指定的管道来进行数据通信，发送原始数据到管道，然后从管道接收相关的数据返回，本方法无锁 (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServerAsync(IEnumerableByte)](bcdb59b7-a70a-d5c0-1fc6-9ececcd39ea2.htm) | 将多个数据报文按顺序发到设备，并从设备接收返回的数据内容，然后拼接成一个Byte[]信息，需要重写[UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm)方法才能返回正确的结果。  Send multiple data packets to the device in sequence, and receive the returned data content from the device, and then splicing them into a Byte[] message, you need to rewrite [UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm) method to return the correct result. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServerAsync(Byte)](61427407-a616-84ae-62c4-abea8bc7a451.htm) | 将Modbus报文数据发送到当前的通道中，并从通道中接收Modbus的报文，通道将根据当前连接自动获取，本方法是线程安全的。  Send Modbus message data to the current channel, and receive Modbus messages from the channel. The channel will automatically obtain it according to the current connection. This method is thread-safe. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServerAsync(Byte, Boolean, Boolean)](1c79c204-5ab5-9342-e65d-7455c20ea577.htm) | 将二进制的数据发送到管道中去，然后从管道里接收二进制的数据回来，并返回是否成功的结果对象。  Send binary data to the pipeline, and then receive binary data back from the pipeline, and return whether the success of the result object (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServerAsync(CommunicationPipe, Byte, Boolean, Boolean)](a254a19d-3c29-147b-335b-8adccdf4ee2f.htm) | 使用指定的管道来进行数据通信，发送原始数据到管道，然后从管道接收相关的数据返回，本方法无锁 (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法代码示例 | [ReadInt16(String)](57a16003-2fd5-b2bd-09bf-14366be899b8.htm) | 读取16位的有符号的整型数据  Read 16-bit signed integer data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt16(String, UInt16)](b2bde158-f51c-93fc-66ca-5e4cf6ea1fb7.htm) | 读取16位的有符号整型数组  Read 16-bit signed integer array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt16Async(String)](ac021793-890d-ceab-eae4-47ff763390c1.htm) | 异步读取16位的有符号的整型数据  Asynchronously read 16-bit signed integer data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt16Async(String, UInt16)](f30e02fc-8078-f1b3-c226-5991e3e23caf.htm) | 异步读取16位的有符号整型数组  Asynchronously read 16-bit signed integer array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt32(String)](bfbea035-9a83-40e3-084c-28122041b37a.htm) | 读取32位的有符号整型  Read 32-bit signed integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt32(String, UInt16)](ca0b8024-b31c-4e45-ad3e-916e188efebe.htm) | 读取32位有符号整型数组  Read 32-bit signed integer array (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [ReadInt32Async(String)](96726ea0-82aa-1229-21f6-43cd62abca43.htm) | 异步读取32位的有符号整型  Asynchronously read 32-bit signed integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt32Async(String, UInt16)](4d8ec145-5fcb-0e5e-6065-957de3434321.htm) | 异步读取32位有符号整型数组  Asynchronously read 32-bit signed integer array (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [ReadInt64(String)](0d0c2284-2da1-6f9a-b353-537561cd5e95.htm) | 读取64位的有符号整型  Read 64-bit signed integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt64(String, UInt16)](662b75d9-9f79-a959-bd6f-a797741d84c8.htm) | 读取64位的有符号整型数组  Read 64-bit signed integer array (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [ReadInt64Async(String)](26d46c8a-3775-3abd-d3ec-5cd1758a317e.htm) | 异步读取64位的有符号整型  Asynchronously read 64-bit signed integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt64Async(String, UInt16)](c3fa5d8b-cf3e-2424-4541-99bb11f7b8d7.htm) | 异步读取64位的有符号整型数组  Asynchronously read 64-bit signed integer array (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
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
| 公共方法代码示例 | [ReadUInt32(String, UInt16)](7ef202c8-7f3a-596b-12c0-bf5e5a5abd92.htm) | 读取32位的无符号整型数组  Read 32-bit unsigned integer array (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [ReadUInt32Async(String)](2343953f-6350-331e-6ba7-18a4449f0575.htm) | 异步读取32位的无符号整型  Asynchronously read 32-bit unsigned integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt32Async(String, UInt16)](ca94295c-48ce-dbdf-01eb-4d01701a9f50.htm) | 异步读取32位的无符号整型数组  Asynchronously read 32-bit unsigned integer array (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [ReadUInt64(String)](947487dd-d4a3-72d8-d9fb-2224aad70e12.htm) | 读取64位的无符号整型  Read 64-bit unsigned integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt64(String, UInt16)](130d0d35-93a1-d72d-c2f4-53b525c11f02.htm) | 读取64位的无符号整型的数组  Read 64-bit unsigned integer array (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [ReadUInt64Async(String)](3411ff3a-772f-dbe1-8931-75e00a731579.htm) | 异步读取64位的无符号整型  Asynchronously read 64-bit unsigned integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt64Async(String, UInt16)](6f22e2db-a0eb-b4b7-f43b-030555667ade.htm) | 异步读取64位的无符号整型的数组  Asynchronously read 64-bit unsigned integer array (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [ReadWrite](455d5de4-9bc6-c0c0-db5c-b6a334800126.htm) | 使用0x17功能码来实现同时写入并读取数据的操作，使用一条报文来实现，需要指定读取的地址，长度，写入的地址，写入的数据信息，返回读取的结果数据。  Use 0x17 function code to write and read data at the same time, and use a message to implement it, you need to specify the read address, length, written address, written data information, and return the read result data. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [RegisteredAddressMapping](a00ba5f5-373c-6e60-845a-3d8ad093568e.htm) | 注册一个新的地址映射关系，注册地址映射关系后，就可以使用新的地址来读写Modbus数据了，通常用于其他的支持Modbus协议的PLC。  After registering a new address mapping, you can use the new address to read and write Modbus data, which is usually used for other PLCs that support the Modbus protocol. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [SerialPortInni(ActionSerialPort)](3728a0f2-4553-0bc8-5c92-253c13e84f05.htm) | 根据自定义初始化方法进行初始化串口信息  Initialize the serial port information according to the custom initialization method (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 公共方法 | [SerialPortInni(String)](f82e0ba4-af2f-5165-dcde-e047c3351e74.htm) | 初始化串口信息，9600波特率，8位数据位，1位停止位，无奇偶校验  Initial serial port information, 9600 baud rate, 8 data bits, 1 stop bit, no parity (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 公共方法 | [SerialPortInni(String, Int32)](433523a8-1606-f029-50db-68ef68034046.htm) | 初始化串口信息，波特率，8位数据位，1位停止位，无奇偶校验  Initializes serial port information, baud rate, 8-bit data bit, 1-bit stop bit, no parity (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 公共方法 | [SerialPortInni(String, Int32, Int32, StopBits, Parity)](cefa7cd2-9c1f-b79d-9ea3-d4d27d7ddeae.htm) | 初始化串口信息，波特率，数据位，停止位，奇偶校验需要全部自己来指定  Start serial port information, baud rate, data bit, stop bit, parity all need to be specified (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 公共方法 | [SetDtuPipe](f765d1d7-234a-b229-9457-97e25e024dc9.htm) | 将当前的通信对象设置DTU模式，允许传入现成的管道，并返回初始化结果，如果该设备重写了握手报文，就是返回握手结果  Set the current communication object to DTU mode, allow the existing pipe to be passed in, and return the initialization result, if the device rewrites the handshake packet, the handshake result is returned (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [SetDtuPipeAsync](9d1fd12d-b804-3828-a627-a2c0eef1d23f.htm) | 将当前的通信对象设置DTU模式，允许传入现成的管道，并返回初始化结果，如果该设备重写了握手报文，就是返回握手结果  Set the current communication object to DTU mode, allow the existing pipe to be passed in, and return the initialization result, if the device rewrites the handshake packet, the handshake result is returned (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ToString](9e0bb9d9-c4ac-4e81-3d89-f9eb55d88718.htm) | (重写 [ModbusRtuToString](11f9cb41-d3e4-d45d-abd1-36e7bb2577e6.htm).) |
| 公共方法 | [TranslateToModbusAddress](43c76b4c-2f7e-c60a-8b90-ee2c3898ce42.htm) | 将当前的地址信息转换成Modbus格式的地址，如果转换失败，返回失败的消息。默认不进行任何的转换。  Convert the current address information into a Modbus format address. If the conversion fails, a failure message will be returned. No conversion is performed by default. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [UnpackResponseContent](00a9f9a8-e32c-331f-cddd-d6d676a819f3.htm) | 根据对方返回的报文命令，对命令进行基本的拆包，例如各种Modbus协议拆包为统一的核心报文，还支持对报文的验证  According to the message command returned by the other party, the command is basically unpacked, for example, various Modbus protocols are unpacked into a unified core message, and the verification of the message is also supported (重写 [ModbusRtuUnpackResponseContent(Byte, Byte)](2af5ade0-337b-6bff-777e-99f624a94606.htm).) |
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
| 公共方法代码示例 | [Write(String, Double)](64c44576-7e95-cbcb-a8c9-e6b7ff7e3c50.htm) | 写入double数据，返回是否成功  Write double data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Int16)](c7cfd67d-ca31-3087-c684-ee4c10780ab1.htm) | 写入short数组，返回是否成功  Write short array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Int32)](ba80d1a5-c8d8-7e0b-5c6b-9435cf20cadb.htm) | 写入int数据，返回是否成功  Write int data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Int64)](ba246b33-5a12-a22c-6d40-0969a64b508b.htm) | 写入long数据，返回是否成功  Write long data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Single)](a2265fa7-5b56-a5aa-9d21-a0dc5458cbf9.htm) | 写入float数据，返回是否成功  Write float data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, String)](4e0ec363-ba64-7ca7-b8d3-2e9a8087bd55.htm) | 写入字符串信息，编码为ASCII  Write string information, encoded as ASCII (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, UInt16)](65544962-82ac-bb47-d72f-5dd0c3d4ac6c.htm) | 写入ushort数组，返回是否成功  Write ushort array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, UInt32)](c99d8246-f95e-4f79-738b-16e622072eb4.htm) | 写入uint数据，返回是否成功  Write uint data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, UInt64)](4924086e-9b36-d446-4ae4-6682134a6ad1.htm) | 写入ulong数据，返回是否成功  Write ulong data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [Write(String, Boolean)](a40c7e8d-5770-9ad2-e263-66a92c45f61c.htm) | 向线圈中写入bool数值，返回是否写入成功，如果富文本地址不指定，默认使用的功能码是 0x05， 如果你的地址为字地址，例如100.2，那么将使用0x16的功能码，通过掩码的方式来修改寄存器的某一位，需要Modbus服务器支持，否则自动切换为读取字数据，修改位，在写入字的方式。  Write bool value to the coil and return whether the writing is successful. If the rich text address is not specified, the default function code is 0x05. If your address is a word address, such as 100.2, then you will use the function code of 0x16 to modify a bit of the register through a mask. It needs Modbus server support, Otherwise, it automatically switches to read the word data, modifies the bits, and writes the word in the way. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [Write(String, Boolean)](fc664b2b-f304-683e-573c-239379b83eb0.htm) | 向线圈中写入bool数组，返回是否写入成功，如果富文本地址不指定，默认使用的功能码是 0x0F  Write the bool array to the coil, and return whether the writing is successful. If the rich text address is not specified, the default function code is 0x0F. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [Write(String, Byte)](ec533f87-fd4a-c833-0256-fdd83da3609f.htm) | 将数据写入到Modbus的寄存器上去，需要指定起始地址和数据内容，如果富文本地址不指定，默认使用的功能码是 0x10  To write data to Modbus registers, you need to specify the start address and data content. If the rich text address is not specified, the default function code is 0x10 (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [Write(String, Double)](e4621d28-0353-fd60-d600-f457a46a661e.htm) | 写入double数组，返回是否成功  Write double array, return whether the write was successful (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [Write(String, Int16)](87345acc-1644-e27d-f5b3-51a01e0d5b02.htm) | 将数据写入到Modbus的单个寄存器上去，需要指定起始地址和数据值，如果富文本地址不指定，默认使用的功能码是 0x06  To write data to a single register of Modbus, you need to specify the start address and data value. If the rich text address is not specified, the default function code is 0x06. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [Write(String, Int32)](cb0f528c-c39e-a439-6093-65c7c8e8ca10.htm) | 写入int[]数组，返回是否成功  Write int array, return whether the write was successful (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [Write(String, Int64)](d3322d92-64ca-0080-27c6-b8693206ad96.htm) | 写入long数组，返回是否成功  Write long array, return whether the write was successful (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [Write(String, Single)](08a2fa24-4d18-9e7d-1fc9-873e8870ecc2.htm) | 写入float数组，返回是否成功  Write float array, return whether the write was successful (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [Write(String, UInt16)](4c43fcf9-d5c0-e4e4-70c2-c2d1aa8bfe4b.htm) | 将数据写入到Modbus的单个寄存器上去，需要指定起始地址和数据值，如果富文本地址不指定，默认使用的功能码是 0x06  To write data to a single register of Modbus, you need to specify the start address and data value. If the rich text address is not specified, the default function code is 0x06. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [Write(String, UInt32)](3113766d-8db6-7f55-c9c9-0580b9b1b756.htm) | 写入uint[]数组，返回是否成功  Write uint array, return whether the write was successful (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [Write(String, UInt64)](414c37ec-7955-1077-addb-8a7c0a861215.htm) | 写入ulong数组，返回是否成功  Write ulong array, return whether the write was successful (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [Write(String, String, Int32)](02fe2ca9-49a3-ede4-d1af-8c6da1f644da.htm) | 写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为ASCII  Write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, it is filled with 0 and the encoding is ASCII. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, String, Encoding)](7b602661-902d-2b70-1fdb-ba87ecd307bf.htm) | 写入字符串信息，需要指定的编码信息  Write string information, need to specify the encoding information (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, String, Int32, Encoding)](e8cb30d9-4c5e-6c68-b9af-1b1b6c2dc714.htm) | 写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为指定的编码信息  Write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, then the operation is complemented with 0 , you should specified the encoding information (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteT(T)](9f48853f-fbeb-c37d-25a2-1c76d1dee407.htm) | 写入支持Hsl特性的数据，返回是否写入成功，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Write data that supports the Hsl attribute, and return whether the write was successful. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WriteAsync(String, Boolean)](fe872dae-7ef9-12c1-d883-e92ba8676416.htm) | 异步批量写入Boolean数组数据，返回是否成功  Asynchronously batch write Boolean array data, return success (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Byte)](8f1cec30-8b86-6749-c8f3-7d39168d59df.htm) | 异步写入原始的byte数组数据到指定的地址，返回是否写入成功  Asynchronously writes the original byte array data to the specified address, and returns whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Double)](068adaa4-d1de-dffc-8376-6f483e2c1d96.htm) | 异步写入double数据，返回是否成功  Asynchronously write double data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int16)](2625e868-7096-e526-ce25-b29afc639b41.htm) | 异步写入short数组，返回是否成功  Asynchronously write short array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int32)](df91d629-2080-328d-67cd-fbefb9a45404.htm) | 异步写入int数据，返回是否成功  Asynchronously write int data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int64)](ec6c1e24-a0d8-f1f5-56af-109b0c546fc8.htm) | 异步写入long数据，返回是否成功  Asynchronously write long data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Single)](910189df-13f9-b7d0-a626-b5b25ce5d32a.htm) | 异步写入float数据，返回是否成功  Asynchronously write float data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, String)](fe839c9a-46f3-e317-0155-b10a8a83a9b1.htm) | 异步写入字符串信息，编码为ASCII  Asynchronously write string information, encoded as ASCII (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt16)](34a1a321-972e-06a8-76fb-794dd1229005.htm) | 异步写入ushort数组，返回是否成功  Asynchronously write ushort array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt32)](abad215c-bc99-ac90-b9f3-7adb3a7cca8f.htm) | 异步写入uint数据，返回是否成功  Asynchronously write uint data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt64)](e396074d-6cf1-b45d-5e8e-01026f6e391b.htm) | 异步写入ulong数据，返回是否成功  Asynchronously write ulong data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WriteAsync(String, Boolean)](90b9711b-1891-383c-0e56-19dd224b7aab.htm) | 向线圈中写入bool数值，返回是否写入成功，如果富文本地址不指定，默认使用的功能码是 0x05， 如果你的地址为字地址，例如100.2，那么将使用0x16的功能码，通过掩码的方式来修改寄存器的某一位，需要Modbus服务器支持，否则自动切换为读取字数据，修改位，在写入字的方式。  Write bool value to the coil and return whether the writing is successful. If the rich text address is not specified, the default function code is 0x05. If your address is a word address, such as 100.2, then you will use the function code of 0x16 to modify a bit of the register through a mask. It needs Modbus server support, Otherwise, it automatically switches to read the word data, modifies the bits, and writes the word in the way. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Double)](d67e053b-fae0-fbbb-d688-b514f417bbd7.htm) | 异步写入double数组，返回是否成功  Asynchronously write double array, return whether the write was successful (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [WriteAsync(String, Int16)](7e819772-f315-0f56-8c1e-8e2bd3aa2396.htm) | 将数据写入到Modbus的单个寄存器上去，需要指定起始地址和数据值，如果富文本地址不指定，默认使用的功能码是 0x06  To write data to a single register of Modbus, you need to specify the start address and data value. If the rich text address is not specified, the default function code is 0x06. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int32)](9e2c12ba-21f7-1c91-71ae-1eb2f2f853c3.htm) | 异步写入int[]数组，返回是否成功  Asynchronously write int array, return whether the write was successful (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int64)](3ef3dd46-318a-f55d-7296-57d303aafc93.htm) | 异步写入long数组，返回是否成功  Asynchronously write long array, return whether the write was successful (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Single)](ffa2fa0f-8321-2c1f-3916-c2cd90882f3c.htm) | 异步写入float数组，返回是否成功  Asynchronously write float array, return whether the write was successful (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [WriteAsync(String, UInt16)](4bb55e75-5cb4-4b22-c8b7-2983eae63e31.htm) | 将数据写入到Modbus的单个寄存器上去，需要指定起始地址和数据值，如果富文本地址不指定，默认使用的功能码是 0x06  To write data to a single register of Modbus, you need to specify the start address and data value. If the rich text address is not specified, the default function code is 0x06. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt32)](e108f604-e11a-2599-b88c-d69e630212b5.htm) | 异步写入uint[]数组，返回是否成功  Asynchronously write uint array, return whether the write was successful (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt64)](d304bc30-22bb-24bc-dd75-f767cfba2c35.htm) | 异步写入ulong数组，返回是否成功  Asynchronously write ulong array, return whether the write was successful (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [WriteAsync(String, String, Int32)](b739b983-b03c-8ecf-2713-75f19233ae2f.htm) | 异步写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为ASCII  Asynchronously write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, it is filled with 0 and the encoding is ASCII. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, String, Encoding)](876e9800-0406-6407-bdf8-8ab1cf28df38.htm) | 异步写入字符串信息，需要指定的编码信息  Asynchronously write string information, need to specify the encoding information (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WriteAsync(String, String, Int32, Encoding)](ae9fb2a3-fddb-48eb-e9ad-6cc0ebbac1f6.htm) | 异步写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为指定的编码信息  Asynchronously write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, then the operation is complemented with 0 , you should specified the encoding information (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsyncT(T)](25d96e77-d1b2-918c-33f9-e5e4cc57c9cd.htm) | 异步写入支持Hsl特性的数据，返回是否写入成功，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Asynchronously write data that supports the Hsl attribute, and return whether the write was successful. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteCustomerT](84a37f82-3f67-a0f8-fd3f-0630772565e4.htm) | 写入自定义类型的数据，该类型必须继承自IDataTransfer接口  Write data of a custom type, which must inherit from the IDataTransfer interface (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteCustomerAsyncT](04628397-cfa1-d785-7790-4c02cb377ab6.htm) | 写入自定义类型的数据，该类型必须继承自IDataTransfer接口  Write data of a custom type, which must inherit from the IDataTransfer interface (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WriteFile](8b9201b8-f363-7860-6b09-e63811574239.htm) | 使用0x15功能码写入文件的数据内容，需要指定文件编号，地址，数据，返回写入的结果数据。写入的数据字节长度不能大于248个字节，如果想写入更多的字节，就需要切片写入操作 (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [WriteMask](6e5e9e07-cc0e-f60c-07be-eb4a63249863.htm) | 向设备写入掩码数据，使用0x16功能码，需要确认对方是否支持相关的操作，掩码数据的操作主要针对寄存器。  To write mask data to the server, using the 0x16 function code, you need to confirm whether the other party supports related operations. The operation of mask data is mainly directed to the register. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [WriteMaskAsync](fd2cff63-8ffc-bf53-a7b6-24253b445aff.htm) | 向设备写入掩码数据，使用0x16功能码，需要确认对方是否支持相关的操作，掩码数据的操作主要针对寄存器。  To write mask data to the server, using the 0x16 function code, you need to confirm whether the other party supports related operations. The operation of mask data is mainly directed to the register. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [WriteOneRegister(String, Int16)](ed003503-3828-ed79-d4b4-08cfc417c8ab.htm) | 将数据写入到Modbus的单个寄存器上去，需要指定起始地址和数据值，如果富文本地址不指定，默认使用的功能码是 0x06  To write data to a single register of Modbus, you need to specify the start address and data value. If the rich text address is not specified, the default function code is 0x06. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [WriteOneRegister(String, UInt16)](9603651d-3202-1560-9141-b2c33949e06c.htm) | 将数据写入到Modbus的单个寄存器上去，需要指定起始地址和数据值，如果富文本地址不指定，默认使用的功能码是 0x06  To write data to a single register of Modbus, you need to specify the start address and data value. If the rich text address is not specified, the default function code is 0x06. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [WriteOneRegisterAsync(String, Int16)](b5f960df-3437-910f-aba4-f4dcbbfedad8.htm) | 将数据写入到Modbus的单个寄存器上去，需要指定起始地址和数据值，如果富文本地址不指定，默认使用的功能码是 0x06  To write data to a single register of Modbus, you need to specify the start address and data value. If the rich text address is not specified, the default function code is 0x06. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [WriteOneRegisterAsync(String, UInt16)](66dad72a-c201-d696-fdf2-bcdac070fc57.htm) | 将数据写入到Modbus的单个寄存器上去，需要指定起始地址和数据值，如果富文本地址不指定，默认使用的功能码是 0x06  To write data to a single register of Modbus, you need to specify the start address and data value. If the rich text address is not specified, the default function code is 0x06. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |

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

本客户端支持的标准的modbus协议，Modbus-Tcp及Modbus-Udp内置的消息号会进行自增，地址支持富文本格式，具体参考示例代码。  
读取线圈，输入线圈，寄存器，输入寄存器的方法中的读取长度对商业授权用户不限制，内部自动切割读取，结果合并。

![](../icons/SectionExpanded.png)示例

基本的用法请参照下面的代码示例，初始化部分的代码省略

Modbus示例

[复制](# "复制")

```
// 本类支持的读写操作提供了非常多的重载方法，总有你想要的方法
private ModbusAscii modbus = new ModbusAscii( );   // 实例化



private void CoilExample( )
{
    // 读取线圈示例
    bool coil100 = modbus.ReadCoil( "100" ).Content;

    // 判断是否读取成功
    OperateResult<bool> result_coil100 = modbus.ReadCoil( "100" );
    if (result_coil100.IsSuccess)
    {
        // success
        bool value = result_coil100.Content;
    }
    else
    {
        // failed
    }


    // 假设读取站号10的线圈100的值
    bool coil_station_ten_100 = modbus.ReadCoil( "s=10;100" ).Content;



    // =============================================================================================
    // 写入也是同理，线圈100写通
    modbus.Write( "100", true );

    // 站号10的线圈写通
    modbus.Write( "s=10;100", true );

    // 想要判断是否写入成功
    if (modbus.Write( "s=10;100", true ).IsSuccess)
    {
        // success
    }
    else
    {
        // failed
    }



    // ===========================================================================================
    // 批量读写也是类似，批量的读取
    bool[] coil10_19 = modbus.ReadCoil( "100", 10 ).Content;

    // 写入也是同理
    modbus.Write( "100", new bool[] { true, false, true, false, false, false, true, false, false, false } );


    // 离散输入的数据读取同理
}


private void RegisterExample( )
{
    // 读取寄存器100的值
    short register100 = modbus.ReadInt16( "100" ).Content;

    // 批量读取寄存器100-109的值
    short[] register100_109 = modbus.ReadInt16( "100", 10 ).Content;

    // 写入寄存器100的值，注意，一定要强制转换short类型
    modbus.Write( "100", (short)123 );

    // 批量写
    modbus.Write( "100", new short[] { 123, -123, 4244 } );


    // ==============================================================================================
    // 以下是一些常规的操作，不再对是否成功的结果进行判断
    // 读取操作

    bool coil100 = modbus.ReadCoil( "100" ).Content;   // 读取线圈100的通断
    short short100 = modbus.ReadInt16( "100" ).Content; // 读取寄存器100的short值
    ushort ushort100 = modbus.ReadUInt16( "100" ).Content; // 读取寄存器100的ushort值
    int int100 = modbus.ReadInt32( "100" ).Content;      // 读取寄存器100-101的int值
    uint uint100 = modbus.ReadUInt32( "100" ).Content;   // 读取寄存器100-101的uint值
    float float100 = modbus.ReadFloat( "100" ).Content; // 读取寄存器100-101的float值
    long long100 = modbus.ReadInt64( "100" ).Content;    // 读取寄存器100-103的long值
    ulong ulong100 = modbus.ReadUInt64( "100" ).Content; // 读取寄存器100-103的ulong值
    double double100 = modbus.ReadDouble( "100" ).Content; // 读取寄存器100-103的double值
    string str100 = modbus.ReadString( "100", 5 ).Content;// 读取100到104共10个字符的字符串

    // 写入操作
    modbus.Write( "100", true );// 写入线圈100为通
    modbus.Write( "100", (short)12345 );// 写入寄存器100为12345
    modbus.Write( "100", (ushort)45678 );// 写入寄存器100为45678
    modbus.Write( "100", 123456789 );// 写入寄存器100-101为123456789
    modbus.Write( "100", (uint)123456778 );// 写入寄存器100-101为123456778
    modbus.Write( "100", 123.456 );// 写入寄存器100-101为123.456
    modbus.Write( "100", 12312312312414L );//写入寄存器100-103为一个大数据
    modbus.Write( "100", 12634534534543656UL );// 写入寄存器100-103为一个大数据
    modbus.Write( "100", 123.456d );// 写入寄存器100-103为一个双精度的数据
    modbus.Write( "100", "K123456789" );

    // ===============================================================================================
    // 读取输入寄存器
    short input_short100 = modbus.ReadInt16( "x=4;100" ).Content; // 读取寄存器100的short值
    ushort input_ushort100 = modbus.ReadUInt16( "x=4;100" ).Content; // 读取寄存器100的ushort值
    int input_int100 = modbus.ReadInt32( "x=4;100" ).Content;      // 读取寄存器100-101的int值
    uint input_uint100 = modbus.ReadUInt32( "x=4;100" ).Content;   // 读取寄存器100-101的uint值
    float input_float100 = modbus.ReadFloat( "x=4;100" ).Content; // 读取寄存器100-101的float值
    long input_long100 = modbus.ReadInt64( "x=4;100" ).Content;    // 读取寄存器100-103的long值
    ulong input_ulong100 = modbus.ReadUInt64( "x=4;100" ).Content; // 读取寄存器100-103的ulong值
    double input_double100 = modbus.ReadDouble( "x=4;100" ).Content; // 读取寄存器100-103的double值
    string input_str100 = modbus.ReadString( "x=4;100", 5 ).Content;// 读取100到104共10个字符的字符串
}
```

复杂的读取数据的代码示例如下：

read示例

[复制](# "复制")

```
ModbusRtu modbus = new ModbusRtu( );   // 实例化
// 此处忽略初始化
// modbus.SerialPortInni( "COM3" );

// 假设100存储了short的报警，101,102存储了float的温度，103，104存储了int的产量
OperateResult<byte[]> read = modbus.Read( "100", 5 );
if (read.IsSuccess)
{
    // 共计10个字节的结果内容
    short alarm = modbus.ByteTransform.TransInt16( read.Content, 0 );
    float temp = modbus.ByteTransform.TransSingle( read.Content, 2 );
    int product = modbus.ByteTransform.TransInt32( read.Content, 6 );
}
else
{
    // failed
}
```

写入数据的代码如下：

write示例

[复制](# "复制")

```
ModbusRtu modbus = new ModbusRtu( );   // 实例化
// 此处忽略初始化
// modbus.SerialPortInni( "COM3" );

// 假设100存储了short的报警，101,102存储了float的温度，103，104存储了int的产量
byte[] buffer = new byte[10];
modbus.ByteTransform.TransByte( (short)1 ).CopyTo( buffer, 0 );
modbus.ByteTransform.TransByte( 123.456f ).CopyTo( buffer, 2 );
modbus.ByteTransform.TransByte( 45678922 ).CopyTo( buffer, 6 );

OperateResult write = modbus.Write( "100", buffer );
if (write.IsSuccess)
{
    // success
}
else
{
    // failed
}
```

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.ModBus 命名空间](da9ac567-746e-47e5-70e1-b44577d09ce4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ModbusAscii 构造函数 

[原文連結](http://api.hslcommunication.cn/html/223428e7-bc23-5a0f-0fa6-ac355d422c71.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.ModBus](../html/da9ac567-746e-47e5-70e1-b44577d09ce4.htm "HslCommunication.ModBus")

[ModbusAscii 类](../html/873a8425-d84d-1a6f-ed5c-79e8d056cc76.htm "ModbusAscii 类")

[ModbusAscii 构造函数](../html/223428e7-bc23-5a0f-0fa6-ac355d422c71.htm "ModbusAscii 构造函数 ")

[ModbusAscii 构造函数](../html/29c891cf-fb6d-b443-ff97-12d8c955c41f.htm "ModbusAscii 构造函数 ")

[ModbusAscii 构造函数 (Byte)](../html/b6fbcbed-323d-9f61-fd2b-9b8dac246caa.htm "ModbusAscii 构造函数 (Byte)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ModbusAscii 构造函数 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [ModbusAscii](29c891cf-fb6d-b443-ff97-12d8c955c41f.htm) | 实例化一个Modbus-ascii协议的客户端对象  Instantiate a client object of the Modbus-ascii protocol |
| 公共方法 | [ModbusAscii(Byte)](b6fbcbed-323d-9f61-fd2b-9b8dac246caa.htm) | 指定Modbus从站的站号来初始化  Specify the station number of the Modbus slave to initialize |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[ModbusAscii 类](873a8425-d84d-1a6f-ed5c-79e8d056cc76.htm)

[HslCommunication.ModBus 命名空间](da9ac567-746e-47e5-70e1-b44577d09ce4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ModbusAscii 构造函数 

[原文連結](http://api.hslcommunication.cn/html/29c891cf-fb6d-b443-ff97-12d8c955c41f.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.ModBus](../html/da9ac567-746e-47e5-70e1-b44577d09ce4.htm "HslCommunication.ModBus")

[ModbusAscii 类](../html/873a8425-d84d-1a6f-ed5c-79e8d056cc76.htm "ModbusAscii 类")

[ModbusAscii 构造函数](../html/223428e7-bc23-5a0f-0fa6-ac355d422c71.htm "ModbusAscii 构造函数 ")

[ModbusAscii 构造函数](../html/29c891cf-fb6d-b443-ff97-12d8c955c41f.htm "ModbusAscii 构造函数 ")

[ModbusAscii 构造函数 (Byte)](../html/b6fbcbed-323d-9f61-fd2b-9b8dac246caa.htm "ModbusAscii 构造函数 (Byte)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ModbusAscii 构造函数 |

实例化一个Modbus-ascii协议的客户端对象  
Instantiate a client object of the Modbus-ascii protocol

**命名空间：**
 [HslCommunication.ModBus](da9ac567-746e-47e5-70e1-b44577d09ce4.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public ModbusAscii()
```

```
Public Sub New
```

```
public:
ModbusAscii()
```

```
new : unit -> ModbusAscii
```

![](../icons/SectionExpanded.png)参见

#### 引用

[ModbusAscii 类](873a8425-d84d-1a6f-ed5c-79e8d056cc76.htm)

[ModbusAscii 重载](223428e7-bc23-5a0f-0fa6-ac355d422c71.htm)

[HslCommunication.ModBus 命名空间](da9ac567-746e-47e5-70e1-b44577d09ce4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ModbusAscii 构造函数 (Byte)

[原文連結](http://api.hslcommunication.cn/html/b6fbcbed-323d-9f61-fd2b-9b8dac246caa.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.ModBus](../html/da9ac567-746e-47e5-70e1-b44577d09ce4.htm "HslCommunication.ModBus")

[ModbusAscii 类](../html/873a8425-d84d-1a6f-ed5c-79e8d056cc76.htm "ModbusAscii 类")

[ModbusAscii 构造函数](../html/223428e7-bc23-5a0f-0fa6-ac355d422c71.htm "ModbusAscii 构造函数 ")

[ModbusAscii 构造函数](../html/29c891cf-fb6d-b443-ff97-12d8c955c41f.htm "ModbusAscii 构造函数 ")

[ModbusAscii 构造函数 (Byte)](../html/b6fbcbed-323d-9f61-fd2b-9b8dac246caa.htm "ModbusAscii 构造函数 (Byte)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ModbusAscii 构造函数 (Byte) |

指定Modbus从站的站号来初始化  
Specify the station number of the Modbus slave to initialize

**命名空间：**
 [HslCommunication.ModBus](da9ac567-746e-47e5-70e1-b44577d09ce4.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public ModbusAscii(
	byte station = 1
)
```

```
Public Sub New ( 
	Optional station As Byte = 1
)
```

```
public:
ModbusAscii(
	unsigned char station = 1
)
```

```
new : 
        ?station : byte 
(* Defaults:
        let _station = defaultArg station 1
*)
-> ModbusAscii
```

#### 参数

station (Optional)
:   类型：SystemByte  
    Modbus从站的站号

![](../icons/SectionExpanded.png)参见

#### 引用

[ModbusAscii 类](873a8425-d84d-1a6f-ed5c-79e8d056cc76.htm)

[ModbusAscii 重载](223428e7-bc23-5a0f-0fa6-ac355d422c71.htm)

[HslCommunication.ModBus 命名空间](da9ac567-746e-47e5-70e1-b44577d09ce4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ModbusAscii 属性

[原文連結](http://api.hslcommunication.cn/html/fde6aaa6-18a7-1888-9bf8-dd886a14100b.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.ModBus](../html/da9ac567-746e-47e5-70e1-b44577d09ce4.htm "HslCommunication.ModBus")

[ModbusAscii 类](../html/873a8425-d84d-1a6f-ed5c-79e8d056cc76.htm "ModbusAscii 类")

[ModbusAscii 构造函数](../html/223428e7-bc23-5a0f-0fa6-ac355d422c71.htm "ModbusAscii 构造函数 ")

[ModbusAscii 属性](../html/fde6aaa6-18a7-1888-9bf8-dd886a14100b.htm "ModbusAscii 属性")

[ModbusAscii 方法](../html/7819f1d9-ca1b-2c0b-6ae3-9a83dfa03b21.htm "ModbusAscii 方法")

[ModbusAscii 字段](../html/62529de8-1a77-85ed-6b9c-9167bdfd5459.htm "ModbusAscii 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ModbusAscii 属性 |

[ModbusAscii](873a8425-d84d-1a6f-ed5c-79e8d056cc76.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [AddressStartWithZero](4d17309a-509d-62c1-9119-3a5735c8315d.htm) | 获取或设置起始的地址是否从0开始，默认为True  Gets or sets whether the starting address starts from 0. The default is True (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共属性 | [BaudRate](5554f563-5e29-699c-435b-c8e7ad5fa0e4.htm) | 当前连接串口信息的波特率  Baud rate of current connection serial port information (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 公共属性 | [BroadcastStation](27196ecf-f9c7-30ac-d439-bf70d26dc328.htm) | 获取或是设置当前广播模式对应的站号，广播模式意味着不接收设备方的数据返回操作，默认为 -1，表示不使用广播模式。  Gets or sets the station number corresponding to the current broadcast mode. Broadcast mode means that the data return operation of the device is not received. The default value is -1, indicating that broadcast mode is not used. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共属性代码示例 | [ByteTransform](f35118c9-3691-f752-ae6f-2c5549714faa.htm) | 当前的数据变换机制，当你需要从字节数据转换类型数据的时候需要。  The current data transformation mechanism is required when you need to convert type data from byte data. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共属性 | [CommunicationPipe](a4fcb477-64d2-8da7-c712-687f0261755b.htm) | 获取或设置当前的管道信息，管道类型为[CommunicationPipe](ad8e903b-d289-21b0-077e-15ede8005221.htm)的继承类，内置了[PipeTcpNet](26b94bd5-2695-36e7-6378-c3d41846be8f.htm)管道，[PipeUdpNet](5a0b810e-187e-826b-0a13-34fd5576aa0e.htm)管道，[PipeSerialPort](9d924607-3dfc-3916-b532-c5dec48fe347.htm)管道等  Get or set the current pipeline information, the pipeline type is [CommunicationPipe](ad8e903b-d289-21b0-077e-15ede8005221.htm) inheritance class, [PipeTcpNet](26b94bd5-2695-36e7-6378-c3d41846be8f.htm) pipeline, [PipeUdpNet](5a0b810e-187e-826b-0a13-34fd5576aa0e.htm) pipeline, [PipeSerialPort](9d924607-3dfc-3916-b532-c5dec48fe347.htm) pipeline, etc (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 公共属性 | [ConnectionId](ce1efb4b-fbc9-4683-163e-0a8e71e99d6b.htm) | 当前连接的唯一ID号，默认为长度20的guid码加随机数组成，方便列表管理，也可以自己指定  The unique ID number of the current connection. The default is a 20-digit guid code plus a random number. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [Crc16CheckEnable](8b8b3de6-eaf7-0484-604f-657dc210ab61.htm) | 获取或设置是否启用CRC16校验码的检查功能，默认启用，如果需要忽略检查CRC16，则设置为 false 即可。  Gets or sets whether to enable the check function of CRC16 check code. It is enabled by default. If you need to ignore the check of CRC16, you can set it to false. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共属性 | [DataFormat](eeb434e5-c320-09f1-6cf4-ff4798c2505a.htm) | 获取或设置数据解析的格式，可选ABCD, BADC，CDAB，DCBA格式，对int,uint,float,double,long,ulong类型有作用  Get or set the format of the data analysis, optional ABCD, BADC, CDAB, DCBA format, effective for int, uint, float, double, long, ulong type (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共属性 | [DisableFunctionCode06](096f15a9-323c-69db-65a1-be3723d914d1.htm) | 获取或设置是否禁用功能码05的写入操作，设置为true时，表示不允许使用功能码 06 进行写入操作。写入单个的字数据也使用 0x10 功能码  Get or set whether to disable the write operation of function code 06. When set to true, it indicates that the write operation of function code 05 is not allowed. Writing individual character data also uses the 0x10 function code (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共属性 | [EnableWriteMaskCode](85372ceb-e693-1067-0285-581b163f8da4.htm) | 获取或设置当前掩码写入的功能码是否激活状态，设置为 false 时，再执行写入位时，会通过读字，修改位，写字的方式来间接实现。  When the function code is set to false, and then the write bit is executed, it will be indirectly implemented by reading, modifying the bit, and writing the word. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共属性 | [IsClearCacheBeforeRead](81dc971f-1f86-8535-ae9a-586b6590a584.htm) | 是否在发送数据前清空缓冲数据，默认是false  Whether to empty the buffer before sending data, the default is false (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 公共属性 | [IsStringReverse](c063a2e6-4c63-cb3f-8956-9164f3c853f0.htm) | 字符串数据是否按照字来反转，默认为False  Whether the string data is reversed according to words. The default is False. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共属性代码示例 | [LogNet](8c56f246-88ac-afb8-8c06-972fafe4f78a.htm) | 组件的日志工具，支持日志记录，只要实例化后，当前网络的基本信息，就以[DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm)等级进行输出  The component's logging tool supports logging. As long as the instantiation of the basic network information, the output will be output at [DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm) (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [PortName](f2b8e4e3-7079-5240-4ff1-14c6c9b376e6.htm) | 当前连接串口信息的端口号名称  The port name of the current connection serial port information (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 公共属性 | [ReceiveEmptyDataCount](f800fd56-7d49-87c4-985b-103a38330078.htm) | 获取或设置连续接收空的数据次数，在数据接收完成时有效，每个单位消耗的时间为[SleepTime](1e3a83ef-b6fc-bc8c-8950-ea3f6ae5c19b.htm)。  Obtain or set the number of consecutive times to receive empty data, which is valid when the data is received, and the time consumed by each unit is [SleepTime](1e3a83ef-b6fc-bc8c-8950-ea3f6ae5c19b.htm) (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 公共属性代码示例 | [ReceiveTimeOut](250560a7-332b-ce1b-5eb0-e3abe4ee6012.htm) | 获取或设置接收服务器反馈的时间，如果为负数，则不接收反馈   Gets or sets the time to receive server feedback, and if it is a negative number, does not receive feedback (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [RtsEnable](dd4057ef-88fe-f7d0-d2c4-f6e15488fe78.htm) | 获取或设置一个值，该值指示在串行通信中是否启用请求发送 (RTS) 信号。  Gets or sets a value indicating whether the request sending (RTS) signal is enabled in serial communication. (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 公共属性 | [SendBeforeHex](72aa1de4-116f-abfe-4adb-814a0be0e560.htm) | 获取或设置在发送通信报文前追加发送的字节信息，HEX格式，通常用于lora组网时，需要携带 00 00 00 02 四个字节的站地址功能。  Obtain or set the byte information sent before sending communication packets, HEX format, usually used for LORA networking, you need to carry 00 00 00 02 four-byte station address function. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [SleepTime](86a5fbb9-3456-6beb-7113-5bbad5eabbc5.htm) | 获取或设置在正式接收对方返回数据前的时候，需要休息的时间，当设置为0的时候，不需要休息。  Get or set the time required to rest before officially receiving the data from the other party. When it is set to 0, no rest is required. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [Station](063d55ef-6207-33e0-9073-a725af469563.htm) | 获取或者重新修改服务器的默认站号信息，当然，你可以再读写的时候动态指定，例如地址 "s=2;100", 更详细的例子可以参考DEMO界面上的地址示例。  Get or re-modify the server's default station number information, of course, you can dynamically specify when reading and writing, such as the address "s=2; 100", a more detailed example can refer to the DEMO interface address example. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共属性 | [StationCheckMatch](03c10d53-c7a2-9256-df12-d780b8efd32c.htm) | 获取或设置当前是否启用站号检查的功能，默认启用，读写数据时将进行站号确认操作，如果需要忽略站号，则设置为 false 即可。  Gets or sets whether the station ID check function is enabled. It is enabled by default. When reading or writing data, the station ID is confirmed. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 受保护的属性 | [WordLength](f1a505cc-e85c-6c8b-b189-58a1ab5753f9.htm) | 一个字单位的数据表示的地址长度，西门子为2，三菱，欧姆龙，modbusTcp就为1，AB PLC无效  The address length represented by one word of data, Siemens is 2, Mitsubishi, Omron, modbusTcp is 1, AB PLC is invalid (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共属性 | [WordReadBatchLength](75150ff2-55eb-9881-6d35-d916a737a46d.htm) | 当进行字读取的时候，商业授权用户支持分批读取，本属性表示分批的长度信息，默认为 120，可以根据实际情况来设置其他值。  When reading of word type, commercial authorization users support batch reading. This attribute represents the length information of the batch, with a default of 120. Other values can be set according to the actual situation. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[ModbusAscii 类](873a8425-d84d-1a6f-ed5c-79e8d056cc76.htm)

[HslCommunication.ModBus 命名空间](da9ac567-746e-47e5-70e1-b44577d09ce4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ModbusAscii 方法

[原文連結](http://api.hslcommunication.cn/html/7819f1d9-ca1b-2c0b-6ae3-9a83dfa03b21.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.ModBus](../html/da9ac567-746e-47e5-70e1-b44577d09ce4.htm "HslCommunication.ModBus")

[ModbusAscii 类](../html/873a8425-d84d-1a6f-ed5c-79e8d056cc76.htm "ModbusAscii 类")

[ModbusAscii 方法](../html/7819f1d9-ca1b-2c0b-6ae3-9a83dfa03b21.htm "ModbusAscii 方法")

[GetNewNetMessage 方法](../html/88dbd66f-58bf-6e1b-c2ab-34c111364480.htm "GetNewNetMessage 方法 ")

[PackCommandWithHeader 方法](../html/73a56f28-e750-b76d-3ebd-bd5e55ad9aae.htm "PackCommandWithHeader 方法 ")

[ToString 方法](../html/9e0bb9d9-c4ac-4e81-3d89-f9eb55d88718.htm "ToString 方法 ")

[UnpackResponseContent 方法](../html/00a9f9a8-e32c-331f-cddd-d6d676a819f3.htm "UnpackResponseContent 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ModbusAscii 方法 |

[ModbusAscii](873a8425-d84d-1a6f-ed5c-79e8d056cc76.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
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
| 受保护的方法 | [GetNewNetMessage](88dbd66f-58bf-6e1b-c2ab-34c111364480.htm) | 获取一个新的消息对象的方法，需要在继承类里面进行重写  The method to get a new message object needs to be overridden in the inheritance class (重写 [ModbusRtuGetNewNetMessage](11f57be7-edd6-53dd-f91f-74c38a04578c.htm).) |
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
| 公共方法 | [PackCommandWithHeader](73a56f28-e750-b76d-3ebd-bd5e55ad9aae.htm) | 对当前的命令进行打包处理，通常是携带命令头内容，标记当前的命令的长度信息，需要进行重写，否则默认不打包  The current command is packaged, usually carrying the content of the command header, marking the length of the current command, and it needs to be rewritten, otherwise it is not packaged by default (重写 [ModbusRtuPackCommandWithHeader(Byte)](ad370105-3cc4-77f9-797d-0ca5edafd761.htm).) |
| 公共方法代码示例 | [Read(String, UInt16)](6571f657-29f0-1869-7e0c-4d4add79fc2b.htm) | 从Modbus服务器批量读取寄存器的信息，需要指定起始地址，读取长度，如果富文本地址不指定，默认使用的功能码是 0x03，如果需要使用04功能码，那么地址就写成 x=4;100  To read the register information from the Modbus server in batches, you need to specify the start address and read length. If the rich text address is not specified, the default function code is 0x03. If you need to use the 04 function code, the address is written as x = 4; 100 (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [ReadT](95abe064-5594-5ae7-a849-98c53f88fee0.htm) | 读取支持Hsl特性的数据内容，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明，或是官网的详细文档  Read the data content of the Hsl attribute. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadAsync(String, UInt16)](d8c14d37-56dc-2a31-1a19-3e2b5b65197d.htm) | 异步批量读取字节数组信息，需要指定地址和长度，返回原始的字节数组  Asynchronous batch read byte array information, need to specify the address and length, return the original byte array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadAsyncT](0522e740-4bf0-2e31-c4c5-e41384321b37.htm) | 异步读取支持Hsl特性的数据内容，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Asynchronously read the data content of the Hsl attribute. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadBool(String)](9b2faa2e-eaa9-7ace-c49f-7a0318f1df1a.htm) | 读取单个的Boolean数据信息  Read a single Boolean data message (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadBool(String, UInt16)](2d125640-18a5-f9a9-29b9-b121915c77ce.htm) | 批量读取线圈或是离散的数据信息，需要指定地址和长度，具体的结果取决于实现，如果富文本地址不指定，默认使用的功能码是 0x01  To read coils or discrete data in batches, you need to specify the address and length. The specific result depends on the implementation. If the rich text address is not specified, the default function code is 0x01. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [ReadBoolAsync(String)](cc37da1f-d673-fe67-78ee-4e465d29a2bd.htm) | 异步读取单个的Boolean数据信息  Asynchronously read a single Boolean data message (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadBoolAsync(String, UInt16)](a45d383d-0095-3ea7-32f8-e75af88d60f8.htm) | 异步批量读取Boolean数组信息，需要指定地址和长度，返回Boolean 数组  Asynchronously batch read Boolean array information, need to specify the address and length, return Boolean array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadCoil(String)](2d2adf84-ab2c-3b95-269f-3a8ddc659018.htm) | 读取线圈，需要指定起始地址，如果富文本地址不指定，默认使用的功能码是 0x01  To read the coil, you need to specify the start address. If the rich text address is not specified, the default function code is 0x01. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [ReadCoil(String, UInt16)](2869f250-557a-4888-5ebc-ea5e98028aed.htm) | 批量的读取线圈，需要指定起始地址，读取长度，如果富文本地址不指定，默认使用的功能码是 0x01  For batch reading coils, you need to specify the start address and read length. If the rich text address is not specified, the default function code is 0x01. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [ReadCoilAsync(String)](d5239965-a5cf-8e54-b4be-e4cb86b38cc0.htm) | 读取线圈，需要指定起始地址，如果富文本地址不指定，默认使用的功能码是 0x01  To read the coil, you need to specify the start address. If the rich text address is not specified, the default function code is 0x01. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [ReadCoilAsync(String, UInt16)](dfc142d3-3390-5eec-1128-6e19cac1ea12.htm) | 批量的读取线圈，需要指定起始地址，读取长度，如果富文本地址不指定，默认使用的功能码是 0x01  For batch reading coils, you need to specify the start address and read length. If the rich text address is not specified, the default function code is 0x01. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [ReadCustomerT(String)](ca069688-d8f0-7dbd-e962-86a10dce83cf.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，返回一个新的类型的实例对象。  To read a custom data type, you need to inherit from the IDataTransfer interface and return an instance object of a new type. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadCustomerT(String, T)](ca7dc1a1-c97d-689f-9819-4dae4af67206.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，传入一个实例，对这个实例进行赋值，并返回该实例的对象。  To read a custom data type, you need to inherit from the IDataTransfer interface, pass in an instance, assign a value to this instance, and return the object of the instance. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadCustomerAsyncT(String)](f3fcda2b-255d-aadd-567e-c7b3df3636d9.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，返回一个新的类型的实例对象。  To read a custom data type, you need to inherit from the IDataTransfer interface and return an instance object of a new type. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadCustomerAsyncT(String, T)](527c9fae-4743-6ebd-7b9a-5ef3a8b7c455.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，传入一个实例，对这个实例进行赋值，并返回该实例的对象。  To read a custom data type, you need to inherit from the IDataTransfer interface, pass in an instance, assign a value to this instance, and return the object of the instance. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadDiscrete(String)](cffe0847-c797-bc7f-71f1-f93092565ec6.htm) | 读取输入线圈，需要指定起始地址，如果富文本地址不指定，默认使用的功能码是 0x02  To read the input coil, you need to specify the start address. If the rich text address is not specified, the default function code is 0x02. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [ReadDiscrete(String, UInt16)](181cbf55-3a34-a988-9aa5-18432b9da97a.htm) | 批量的读取输入点，需要指定起始地址，读取长度，如果富文本地址不指定，默认使用的功能码是 0x02  To read input points in batches, you need to specify the start address and read length. If the rich text address is not specified, the default function code is 0x02 (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [ReadDiscreteAsync(String)](9d11737c-b5d9-5c24-09a1-0b70a87bdfc3.htm) | 读取输入线圈，需要指定起始地址，如果富文本地址不指定，默认使用的功能码是 0x02  To read the input coil, you need to specify the start address. If the rich text address is not specified, the default function code is 0x02. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [ReadDiscreteAsync(String, UInt16)](9c34359d-3d18-fdbf-d6c3-35134eec3b9d.htm) | 批量的读取输入点，需要指定起始地址，读取长度，如果富文本地址不指定，默认使用的功能码是 0x02  To read input points in batches, you need to specify the start address and read length. If the rich text address is not specified, the default function code is 0x02 (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [ReadDouble(String)](e7e31c43-71b4-e692-6bdc-ba9533e1d83d.htm) | 读取双浮点的数据  Read double floating point data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadDouble(String, UInt16)](9df438ac-8869-5002-7ace-bfb7c8f7aa6b.htm) | 读取双浮点数据的数组  Read double floating point data array (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [ReadDoubleAsync(String)](2cf895fa-1033-1def-f048-0e8caefaa71d.htm) | 异步读取双浮点的数据  Asynchronously read double floating point data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadDoubleAsync(String, UInt16)](e5c602ad-d3ce-a677-908b-03d72a4388cd.htm) | 异步读取双浮点数据的数组  Asynchronously read double floating point data array (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [ReadFile](702ddeff-3ba5-952d-abce-99a988a1d3cc.htm) | 使用0x14功能码读取文件的数据内容，需要指定文件编号，地址，长度，返回读取的结果数据，读取的长度为任意长度，内部自动根据实际情况来切割操作。 (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [ReadFloat(String)](bc2801ae-6600-e845-7c6c-522c9156eb0e.htm) | 读取单浮点数据  Read single floating point data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadFloat(String, UInt16)](909486dc-cc84-c34f-2de5-ec17b3d04182.htm) | 读取单浮点精度的数组  Read single floating point array (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [ReadFloatAsync(String)](89f743f8-4b3a-7085-1e49-a1d7adb10052.htm) | 异步读取单浮点数据  Asynchronously read single floating point data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadFloatAsync(String, UInt16)](eaeb9795-d78c-f884-fa54-c698c5e37404.htm) | 异步读取单浮点精度的数组  Asynchronously read single floating point array (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [ReadFromCoreServer(IEnumerableByte)](044cfdd9-ec76-b89f-2f2b-5f439025006c.htm) | 将多个数据报文按顺序发到设备，并从设备接收返回的数据内容，然后拼接成一个Byte[]信息，需要重写[UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm)方法才能返回正确的结果。  Send multiple data packets to the device in sequence, and receive the returned data content from the device, and then splicing them into a Byte[] message, you need to rewrite [UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm) method to return the correct result. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServer(Byte)](ff061b13-df60-6617-08c2-40b81cde38ea.htm) | 将Modbus报文数据发送到当前的通道中，并从通道中接收Modbus的报文，通道将根据当前连接自动获取，本方法是线程安全的。  Send Modbus message data to the current channel, and receive Modbus messages from the channel. The channel will automatically obtain it according to the current connection. This method is thread-safe. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServer(Byte, Boolean, Boolean)](b32ce8e8-a959-fcc2-6736-14b39fdb9a40.htm) | 将二进制的数据发送到管道中去，然后从管道里接收二进制的数据回来，并返回是否成功的结果对象。  Send binary data to the pipeline, and then receive binary data back from the pipeline, and return whether the success of the result object (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServer(CommunicationPipe, Byte, Boolean, Boolean)](ef36d56b-df09-f727-e8d3-316d37a13ee9.htm) | 使用指定的管道来进行数据通信，发送原始数据到管道，然后从管道接收相关的数据返回，本方法无锁 (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServerAsync(IEnumerableByte)](bcdb59b7-a70a-d5c0-1fc6-9ececcd39ea2.htm) | 将多个数据报文按顺序发到设备，并从设备接收返回的数据内容，然后拼接成一个Byte[]信息，需要重写[UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm)方法才能返回正确的结果。  Send multiple data packets to the device in sequence, and receive the returned data content from the device, and then splicing them into a Byte[] message, you need to rewrite [UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm) method to return the correct result. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServerAsync(Byte)](61427407-a616-84ae-62c4-abea8bc7a451.htm) | 将Modbus报文数据发送到当前的通道中，并从通道中接收Modbus的报文，通道将根据当前连接自动获取，本方法是线程安全的。  Send Modbus message data to the current channel, and receive Modbus messages from the channel. The channel will automatically obtain it according to the current connection. This method is thread-safe. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServerAsync(Byte, Boolean, Boolean)](1c79c204-5ab5-9342-e65d-7455c20ea577.htm) | 将二进制的数据发送到管道中去，然后从管道里接收二进制的数据回来，并返回是否成功的结果对象。  Send binary data to the pipeline, and then receive binary data back from the pipeline, and return whether the success of the result object (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServerAsync(CommunicationPipe, Byte, Boolean, Boolean)](a254a19d-3c29-147b-335b-8adccdf4ee2f.htm) | 使用指定的管道来进行数据通信，发送原始数据到管道，然后从管道接收相关的数据返回，本方法无锁 (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法代码示例 | [ReadInt16(String)](57a16003-2fd5-b2bd-09bf-14366be899b8.htm) | 读取16位的有符号的整型数据  Read 16-bit signed integer data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt16(String, UInt16)](b2bde158-f51c-93fc-66ca-5e4cf6ea1fb7.htm) | 读取16位的有符号整型数组  Read 16-bit signed integer array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt16Async(String)](ac021793-890d-ceab-eae4-47ff763390c1.htm) | 异步读取16位的有符号的整型数据  Asynchronously read 16-bit signed integer data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt16Async(String, UInt16)](f30e02fc-8078-f1b3-c226-5991e3e23caf.htm) | 异步读取16位的有符号整型数组  Asynchronously read 16-bit signed integer array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt32(String)](bfbea035-9a83-40e3-084c-28122041b37a.htm) | 读取32位的有符号整型  Read 32-bit signed integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt32(String, UInt16)](ca0b8024-b31c-4e45-ad3e-916e188efebe.htm) | 读取32位有符号整型数组  Read 32-bit signed integer array (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [ReadInt32Async(String)](96726ea0-82aa-1229-21f6-43cd62abca43.htm) | 异步读取32位的有符号整型  Asynchronously read 32-bit signed integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt32Async(String, UInt16)](4d8ec145-5fcb-0e5e-6065-957de3434321.htm) | 异步读取32位有符号整型数组  Asynchronously read 32-bit signed integer array (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [ReadInt64(String)](0d0c2284-2da1-6f9a-b353-537561cd5e95.htm) | 读取64位的有符号整型  Read 64-bit signed integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt64(String, UInt16)](662b75d9-9f79-a959-bd6f-a797741d84c8.htm) | 读取64位的有符号整型数组  Read 64-bit signed integer array (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [ReadInt64Async(String)](26d46c8a-3775-3abd-d3ec-5cd1758a317e.htm) | 异步读取64位的有符号整型  Asynchronously read 64-bit signed integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt64Async(String, UInt16)](c3fa5d8b-cf3e-2424-4541-99bb11f7b8d7.htm) | 异步读取64位的有符号整型数组  Asynchronously read 64-bit signed integer array (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
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
| 公共方法代码示例 | [ReadUInt32(String, UInt16)](7ef202c8-7f3a-596b-12c0-bf5e5a5abd92.htm) | 读取32位的无符号整型数组  Read 32-bit unsigned integer array (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [ReadUInt32Async(String)](2343953f-6350-331e-6ba7-18a4449f0575.htm) | 异步读取32位的无符号整型  Asynchronously read 32-bit unsigned integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt32Async(String, UInt16)](ca94295c-48ce-dbdf-01eb-4d01701a9f50.htm) | 异步读取32位的无符号整型数组  Asynchronously read 32-bit unsigned integer array (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [ReadUInt64(String)](947487dd-d4a3-72d8-d9fb-2224aad70e12.htm) | 读取64位的无符号整型  Read 64-bit unsigned integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt64(String, UInt16)](130d0d35-93a1-d72d-c2f4-53b525c11f02.htm) | 读取64位的无符号整型的数组  Read 64-bit unsigned integer array (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [ReadUInt64Async(String)](3411ff3a-772f-dbe1-8931-75e00a731579.htm) | 异步读取64位的无符号整型  Asynchronously read 64-bit unsigned integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt64Async(String, UInt16)](6f22e2db-a0eb-b4b7-f43b-030555667ade.htm) | 异步读取64位的无符号整型的数组  Asynchronously read 64-bit unsigned integer array (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [ReadWrite](455d5de4-9bc6-c0c0-db5c-b6a334800126.htm) | 使用0x17功能码来实现同时写入并读取数据的操作，使用一条报文来实现，需要指定读取的地址，长度，写入的地址，写入的数据信息，返回读取的结果数据。  Use 0x17 function code to write and read data at the same time, and use a message to implement it, you need to specify the read address, length, written address, written data information, and return the read result data. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [RegisteredAddressMapping](a00ba5f5-373c-6e60-845a-3d8ad093568e.htm) | 注册一个新的地址映射关系，注册地址映射关系后，就可以使用新的地址来读写Modbus数据了，通常用于其他的支持Modbus协议的PLC。  After registering a new address mapping, you can use the new address to read and write Modbus data, which is usually used for other PLCs that support the Modbus protocol. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [SerialPortInni(ActionSerialPort)](3728a0f2-4553-0bc8-5c92-253c13e84f05.htm) | 根据自定义初始化方法进行初始化串口信息  Initialize the serial port information according to the custom initialization method (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 公共方法 | [SerialPortInni(String)](f82e0ba4-af2f-5165-dcde-e047c3351e74.htm) | 初始化串口信息，9600波特率，8位数据位，1位停止位，无奇偶校验  Initial serial port information, 9600 baud rate, 8 data bits, 1 stop bit, no parity (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 公共方法 | [SerialPortInni(String, Int32)](433523a8-1606-f029-50db-68ef68034046.htm) | 初始化串口信息，波特率，8位数据位，1位停止位，无奇偶校验  Initializes serial port information, baud rate, 8-bit data bit, 1-bit stop bit, no parity (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 公共方法 | [SerialPortInni(String, Int32, Int32, StopBits, Parity)](cefa7cd2-9c1f-b79d-9ea3-d4d27d7ddeae.htm) | 初始化串口信息，波特率，数据位，停止位，奇偶校验需要全部自己来指定  Start serial port information, baud rate, data bit, stop bit, parity all need to be specified (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 公共方法 | [SetDtuPipe](f765d1d7-234a-b229-9457-97e25e024dc9.htm) | 将当前的通信对象设置DTU模式，允许传入现成的管道，并返回初始化结果，如果该设备重写了握手报文，就是返回握手结果  Set the current communication object to DTU mode, allow the existing pipe to be passed in, and return the initialization result, if the device rewrites the handshake packet, the handshake result is returned (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [SetDtuPipeAsync](9d1fd12d-b804-3828-a627-a2c0eef1d23f.htm) | 将当前的通信对象设置DTU模式，允许传入现成的管道，并返回初始化结果，如果该设备重写了握手报文，就是返回握手结果  Set the current communication object to DTU mode, allow the existing pipe to be passed in, and return the initialization result, if the device rewrites the handshake packet, the handshake result is returned (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ToString](9e0bb9d9-c4ac-4e81-3d89-f9eb55d88718.htm) | (重写 [ModbusRtuToString](11f9cb41-d3e4-d45d-abd1-36e7bb2577e6.htm).) |
| 公共方法 | [TranslateToModbusAddress](43c76b4c-2f7e-c60a-8b90-ee2c3898ce42.htm) | 将当前的地址信息转换成Modbus格式的地址，如果转换失败，返回失败的消息。默认不进行任何的转换。  Convert the current address information into a Modbus format address. If the conversion fails, a failure message will be returned. No conversion is performed by default. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [UnpackResponseContent](00a9f9a8-e32c-331f-cddd-d6d676a819f3.htm) | 根据对方返回的报文命令，对命令进行基本的拆包，例如各种Modbus协议拆包为统一的核心报文，还支持对报文的验证  According to the message command returned by the other party, the command is basically unpacked, for example, various Modbus protocols are unpacked into a unified core message, and the verification of the message is also supported (重写 [ModbusRtuUnpackResponseContent(Byte, Byte)](2af5ade0-337b-6bff-777e-99f624a94606.htm).) |
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
| 公共方法代码示例 | [Write(String, Double)](64c44576-7e95-cbcb-a8c9-e6b7ff7e3c50.htm) | 写入double数据，返回是否成功  Write double data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Int16)](c7cfd67d-ca31-3087-c684-ee4c10780ab1.htm) | 写入short数组，返回是否成功  Write short array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Int32)](ba80d1a5-c8d8-7e0b-5c6b-9435cf20cadb.htm) | 写入int数据，返回是否成功  Write int data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Int64)](ba246b33-5a12-a22c-6d40-0969a64b508b.htm) | 写入long数据，返回是否成功  Write long data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Single)](a2265fa7-5b56-a5aa-9d21-a0dc5458cbf9.htm) | 写入float数据，返回是否成功  Write float data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, String)](4e0ec363-ba64-7ca7-b8d3-2e9a8087bd55.htm) | 写入字符串信息，编码为ASCII  Write string information, encoded as ASCII (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, UInt16)](65544962-82ac-bb47-d72f-5dd0c3d4ac6c.htm) | 写入ushort数组，返回是否成功  Write ushort array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, UInt32)](c99d8246-f95e-4f79-738b-16e622072eb4.htm) | 写入uint数据，返回是否成功  Write uint data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, UInt64)](4924086e-9b36-d446-4ae4-6682134a6ad1.htm) | 写入ulong数据，返回是否成功  Write ulong data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [Write(String, Boolean)](a40c7e8d-5770-9ad2-e263-66a92c45f61c.htm) | 向线圈中写入bool数值，返回是否写入成功，如果富文本地址不指定，默认使用的功能码是 0x05， 如果你的地址为字地址，例如100.2，那么将使用0x16的功能码，通过掩码的方式来修改寄存器的某一位，需要Modbus服务器支持，否则自动切换为读取字数据，修改位，在写入字的方式。  Write bool value to the coil and return whether the writing is successful. If the rich text address is not specified, the default function code is 0x05. If your address is a word address, such as 100.2, then you will use the function code of 0x16 to modify a bit of the register through a mask. It needs Modbus server support, Otherwise, it automatically switches to read the word data, modifies the bits, and writes the word in the way. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [Write(String, Boolean)](fc664b2b-f304-683e-573c-239379b83eb0.htm) | 向线圈中写入bool数组，返回是否写入成功，如果富文本地址不指定，默认使用的功能码是 0x0F  Write the bool array to the coil, and return whether the writing is successful. If the rich text address is not specified, the default function code is 0x0F. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [Write(String, Byte)](ec533f87-fd4a-c833-0256-fdd83da3609f.htm) | 将数据写入到Modbus的寄存器上去，需要指定起始地址和数据内容，如果富文本地址不指定，默认使用的功能码是 0x10  To write data to Modbus registers, you need to specify the start address and data content. If the rich text address is not specified, the default function code is 0x10 (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [Write(String, Double)](e4621d28-0353-fd60-d600-f457a46a661e.htm) | 写入double数组，返回是否成功  Write double array, return whether the write was successful (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [Write(String, Int16)](87345acc-1644-e27d-f5b3-51a01e0d5b02.htm) | 将数据写入到Modbus的单个寄存器上去，需要指定起始地址和数据值，如果富文本地址不指定，默认使用的功能码是 0x06  To write data to a single register of Modbus, you need to specify the start address and data value. If the rich text address is not specified, the default function code is 0x06. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [Write(String, Int32)](cb0f528c-c39e-a439-6093-65c7c8e8ca10.htm) | 写入int[]数组，返回是否成功  Write int array, return whether the write was successful (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [Write(String, Int64)](d3322d92-64ca-0080-27c6-b8693206ad96.htm) | 写入long数组，返回是否成功  Write long array, return whether the write was successful (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [Write(String, Single)](08a2fa24-4d18-9e7d-1fc9-873e8870ecc2.htm) | 写入float数组，返回是否成功  Write float array, return whether the write was successful (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [Write(String, UInt16)](4c43fcf9-d5c0-e4e4-70c2-c2d1aa8bfe4b.htm) | 将数据写入到Modbus的单个寄存器上去，需要指定起始地址和数据值，如果富文本地址不指定，默认使用的功能码是 0x06  To write data to a single register of Modbus, you need to specify the start address and data value. If the rich text address is not specified, the default function code is 0x06. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [Write(String, UInt32)](3113766d-8db6-7f55-c9c9-0580b9b1b756.htm) | 写入uint[]数组，返回是否成功  Write uint array, return whether the write was successful (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [Write(String, UInt64)](414c37ec-7955-1077-addb-8a7c0a861215.htm) | 写入ulong数组，返回是否成功  Write ulong array, return whether the write was successful (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [Write(String, String, Int32)](02fe2ca9-49a3-ede4-d1af-8c6da1f644da.htm) | 写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为ASCII  Write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, it is filled with 0 and the encoding is ASCII. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, String, Encoding)](7b602661-902d-2b70-1fdb-ba87ecd307bf.htm) | 写入字符串信息，需要指定的编码信息  Write string information, need to specify the encoding information (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, String, Int32, Encoding)](e8cb30d9-4c5e-6c68-b9af-1b1b6c2dc714.htm) | 写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为指定的编码信息  Write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, then the operation is complemented with 0 , you should specified the encoding information (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteT(T)](9f48853f-fbeb-c37d-25a2-1c76d1dee407.htm) | 写入支持Hsl特性的数据，返回是否写入成功，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Write data that supports the Hsl attribute, and return whether the write was successful. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WriteAsync(String, Boolean)](fe872dae-7ef9-12c1-d883-e92ba8676416.htm) | 异步批量写入Boolean数组数据，返回是否成功  Asynchronously batch write Boolean array data, return success (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Byte)](8f1cec30-8b86-6749-c8f3-7d39168d59df.htm) | 异步写入原始的byte数组数据到指定的地址，返回是否写入成功  Asynchronously writes the original byte array data to the specified address, and returns whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Double)](068adaa4-d1de-dffc-8376-6f483e2c1d96.htm) | 异步写入double数据，返回是否成功  Asynchronously write double data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int16)](2625e868-7096-e526-ce25-b29afc639b41.htm) | 异步写入short数组，返回是否成功  Asynchronously write short array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int32)](df91d629-2080-328d-67cd-fbefb9a45404.htm) | 异步写入int数据，返回是否成功  Asynchronously write int data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int64)](ec6c1e24-a0d8-f1f5-56af-109b0c546fc8.htm) | 异步写入long数据，返回是否成功  Asynchronously write long data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Single)](910189df-13f9-b7d0-a626-b5b25ce5d32a.htm) | 异步写入float数据，返回是否成功  Asynchronously write float data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, String)](fe839c9a-46f3-e317-0155-b10a8a83a9b1.htm) | 异步写入字符串信息，编码为ASCII  Asynchronously write string information, encoded as ASCII (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt16)](34a1a321-972e-06a8-76fb-794dd1229005.htm) | 异步写入ushort数组，返回是否成功  Asynchronously write ushort array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt32)](abad215c-bc99-ac90-b9f3-7adb3a7cca8f.htm) | 异步写入uint数据，返回是否成功  Asynchronously write uint data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt64)](e396074d-6cf1-b45d-5e8e-01026f6e391b.htm) | 异步写入ulong数据，返回是否成功  Asynchronously write ulong data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WriteAsync(String, Boolean)](90b9711b-1891-383c-0e56-19dd224b7aab.htm) | 向线圈中写入bool数值，返回是否写入成功，如果富文本地址不指定，默认使用的功能码是 0x05， 如果你的地址为字地址，例如100.2，那么将使用0x16的功能码，通过掩码的方式来修改寄存器的某一位，需要Modbus服务器支持，否则自动切换为读取字数据，修改位，在写入字的方式。  Write bool value to the coil and return whether the writing is successful. If the rich text address is not specified, the default function code is 0x05. If your address is a word address, such as 100.2, then you will use the function code of 0x16 to modify a bit of the register through a mask. It needs Modbus server support, Otherwise, it automatically switches to read the word data, modifies the bits, and writes the word in the way. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Double)](d67e053b-fae0-fbbb-d688-b514f417bbd7.htm) | 异步写入double数组，返回是否成功  Asynchronously write double array, return whether the write was successful (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [WriteAsync(String, Int16)](7e819772-f315-0f56-8c1e-8e2bd3aa2396.htm) | 将数据写入到Modbus的单个寄存器上去，需要指定起始地址和数据值，如果富文本地址不指定，默认使用的功能码是 0x06  To write data to a single register of Modbus, you need to specify the start address and data value. If the rich text address is not specified, the default function code is 0x06. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int32)](9e2c12ba-21f7-1c91-71ae-1eb2f2f853c3.htm) | 异步写入int[]数组，返回是否成功  Asynchronously write int array, return whether the write was successful (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int64)](3ef3dd46-318a-f55d-7296-57d303aafc93.htm) | 异步写入long数组，返回是否成功  Asynchronously write long array, return whether the write was successful (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Single)](ffa2fa0f-8321-2c1f-3916-c2cd90882f3c.htm) | 异步写入float数组，返回是否成功  Asynchronously write float array, return whether the write was successful (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [WriteAsync(String, UInt16)](4bb55e75-5cb4-4b22-c8b7-2983eae63e31.htm) | 将数据写入到Modbus的单个寄存器上去，需要指定起始地址和数据值，如果富文本地址不指定，默认使用的功能码是 0x06  To write data to a single register of Modbus, you need to specify the start address and data value. If the rich text address is not specified, the default function code is 0x06. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt32)](e108f604-e11a-2599-b88c-d69e630212b5.htm) | 异步写入uint[]数组，返回是否成功  Asynchronously write uint array, return whether the write was successful (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt64)](d304bc30-22bb-24bc-dd75-f767cfba2c35.htm) | 异步写入ulong数组，返回是否成功  Asynchronously write ulong array, return whether the write was successful (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [WriteAsync(String, String, Int32)](b739b983-b03c-8ecf-2713-75f19233ae2f.htm) | 异步写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为ASCII  Asynchronously write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, it is filled with 0 and the encoding is ASCII. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, String, Encoding)](876e9800-0406-6407-bdf8-8ab1cf28df38.htm) | 异步写入字符串信息，需要指定的编码信息  Asynchronously write string information, need to specify the encoding information (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WriteAsync(String, String, Int32, Encoding)](ae9fb2a3-fddb-48eb-e9ad-6cc0ebbac1f6.htm) | 异步写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为指定的编码信息  Asynchronously write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, then the operation is complemented with 0 , you should specified the encoding information (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsyncT(T)](25d96e77-d1b2-918c-33f9-e5e4cc57c9cd.htm) | 异步写入支持Hsl特性的数据，返回是否写入成功，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Asynchronously write data that supports the Hsl attribute, and return whether the write was successful. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteCustomerT](84a37f82-3f67-a0f8-fd3f-0630772565e4.htm) | 写入自定义类型的数据，该类型必须继承自IDataTransfer接口  Write data of a custom type, which must inherit from the IDataTransfer interface (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteCustomerAsyncT](04628397-cfa1-d785-7790-4c02cb377ab6.htm) | 写入自定义类型的数据，该类型必须继承自IDataTransfer接口  Write data of a custom type, which must inherit from the IDataTransfer interface (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WriteFile](8b9201b8-f363-7860-6b09-e63811574239.htm) | 使用0x15功能码写入文件的数据内容，需要指定文件编号，地址，数据，返回写入的结果数据。写入的数据字节长度不能大于248个字节，如果想写入更多的字节，就需要切片写入操作 (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [WriteMask](6e5e9e07-cc0e-f60c-07be-eb4a63249863.htm) | 向设备写入掩码数据，使用0x16功能码，需要确认对方是否支持相关的操作，掩码数据的操作主要针对寄存器。  To write mask data to the server, using the 0x16 function code, you need to confirm whether the other party supports related operations. The operation of mask data is mainly directed to the register. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [WriteMaskAsync](fd2cff63-8ffc-bf53-a7b6-24253b445aff.htm) | 向设备写入掩码数据，使用0x16功能码，需要确认对方是否支持相关的操作，掩码数据的操作主要针对寄存器。  To write mask data to the server, using the 0x16 function code, you need to confirm whether the other party supports related operations. The operation of mask data is mainly directed to the register. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [WriteOneRegister(String, Int16)](ed003503-3828-ed79-d4b4-08cfc417c8ab.htm) | 将数据写入到Modbus的单个寄存器上去，需要指定起始地址和数据值，如果富文本地址不指定，默认使用的功能码是 0x06  To write data to a single register of Modbus, you need to specify the start address and data value. If the rich text address is not specified, the default function code is 0x06. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [WriteOneRegister(String, UInt16)](9603651d-3202-1560-9141-b2c33949e06c.htm) | 将数据写入到Modbus的单个寄存器上去，需要指定起始地址和数据值，如果富文本地址不指定，默认使用的功能码是 0x06  To write data to a single register of Modbus, you need to specify the start address and data value. If the rich text address is not specified, the default function code is 0x06. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [WriteOneRegisterAsync(String, Int16)](b5f960df-3437-910f-aba4-f4dcbbfedad8.htm) | 将数据写入到Modbus的单个寄存器上去，需要指定起始地址和数据值，如果富文本地址不指定，默认使用的功能码是 0x06  To write data to a single register of Modbus, you need to specify the start address and data value. If the rich text address is not specified, the default function code is 0x06. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [WriteOneRegisterAsync(String, UInt16)](66dad72a-c201-d696-fdf2-bcdac070fc57.htm) | 将数据写入到Modbus的单个寄存器上去，需要指定起始地址和数据值，如果富文本地址不指定，默认使用的功能码是 0x06  To write data to a single register of Modbus, you need to specify the start address and data value. If the rich text address is not specified, the default function code is 0x06. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[ModbusAscii 类](873a8425-d84d-1a6f-ed5c-79e8d056cc76.htm)

[HslCommunication.ModBus 命名空间](da9ac567-746e-47e5-70e1-b44577d09ce4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GetNewNetMessage 方法 

[原文連結](http://api.hslcommunication.cn/html/88dbd66f-58bf-6e1b-c2ab-34c111364480.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.ModBus](../html/da9ac567-746e-47e5-70e1-b44577d09ce4.htm "HslCommunication.ModBus")

[ModbusAscii 类](../html/873a8425-d84d-1a6f-ed5c-79e8d056cc76.htm "ModbusAscii 类")

[ModbusAscii 方法](../html/7819f1d9-ca1b-2c0b-6ae3-9a83dfa03b21.htm "ModbusAscii 方法")

[GetNewNetMessage 方法](../html/88dbd66f-58bf-6e1b-c2ab-34c111364480.htm "GetNewNetMessage 方法 ")

[PackCommandWithHeader 方法](../html/73a56f28-e750-b76d-3ebd-bd5e55ad9aae.htm "PackCommandWithHeader 方法 ")

[ToString 方法](../html/9e0bb9d9-c4ac-4e81-3d89-f9eb55d88718.htm "ToString 方法 ")

[UnpackResponseContent 方法](../html/00a9f9a8-e32c-331f-cddd-d6d676a819f3.htm "UnpackResponseContent 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ModbusAsciiGetNewNetMessage 方法 |

获取一个新的消息对象的方法，需要在继承类里面进行重写  
The method to get a new message object needs to be overridden in the inheritance class

**命名空间：**
 [HslCommunication.ModBus](da9ac567-746e-47e5-70e1-b44577d09ce4.htm)  
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

[ModbusAscii 类](873a8425-d84d-1a6f-ed5c-79e8d056cc76.htm)

[HslCommunication.ModBus 命名空间](da9ac567-746e-47e5-70e1-b44577d09ce4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## PackCommandWithHeader 方法 

[原文連結](http://api.hslcommunication.cn/html/73a56f28-e750-b76d-3ebd-bd5e55ad9aae.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.ModBus](../html/da9ac567-746e-47e5-70e1-b44577d09ce4.htm "HslCommunication.ModBus")

[ModbusAscii 类](../html/873a8425-d84d-1a6f-ed5c-79e8d056cc76.htm "ModbusAscii 类")

[ModbusAscii 方法](../html/7819f1d9-ca1b-2c0b-6ae3-9a83dfa03b21.htm "ModbusAscii 方法")

[GetNewNetMessage 方法](../html/88dbd66f-58bf-6e1b-c2ab-34c111364480.htm "GetNewNetMessage 方法 ")

[PackCommandWithHeader 方法](../html/73a56f28-e750-b76d-3ebd-bd5e55ad9aae.htm "PackCommandWithHeader 方法 ")

[ToString 方法](../html/9e0bb9d9-c4ac-4e81-3d89-f9eb55d88718.htm "ToString 方法 ")

[UnpackResponseContent 方法](../html/00a9f9a8-e32c-331f-cddd-d6d676a819f3.htm "UnpackResponseContent 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ModbusAsciiPackCommandWithHeader 方法 |

对当前的命令进行打包处理，通常是携带命令头内容，标记当前的命令的长度信息，需要进行重写，否则默认不打包  
The current command is packaged, usually carrying the content of the command header, marking the length of the current command,
and it needs to be rewritten, otherwise it is not packaged by default

**命名空间：**
 [HslCommunication.ModBus](da9ac567-746e-47e5-70e1-b44577d09ce4.htm)  
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

[ModbusAscii 类](873a8425-d84d-1a6f-ed5c-79e8d056cc76.htm)

[HslCommunication.ModBus 命名空间](da9ac567-746e-47e5-70e1-b44577d09ce4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ToString 方法 

[原文連結](http://api.hslcommunication.cn/html/9e0bb9d9-c4ac-4e81-3d89-f9eb55d88718.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.ModBus](../html/da9ac567-746e-47e5-70e1-b44577d09ce4.htm "HslCommunication.ModBus")

[ModbusAscii 类](../html/873a8425-d84d-1a6f-ed5c-79e8d056cc76.htm "ModbusAscii 类")

[ModbusAscii 方法](../html/7819f1d9-ca1b-2c0b-6ae3-9a83dfa03b21.htm "ModbusAscii 方法")

[GetNewNetMessage 方法](../html/88dbd66f-58bf-6e1b-c2ab-34c111364480.htm "GetNewNetMessage 方法 ")

[PackCommandWithHeader 方法](../html/73a56f28-e750-b76d-3ebd-bd5e55ad9aae.htm "PackCommandWithHeader 方法 ")

[ToString 方法](../html/9e0bb9d9-c4ac-4e81-3d89-f9eb55d88718.htm "ToString 方法 ")

[UnpackResponseContent 方法](../html/00a9f9a8-e32c-331f-cddd-d6d676a819f3.htm "UnpackResponseContent 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ModbusAsciiToString 方法 |

[缺少 "M:HslCommunication.ModBus.ModbusAscii.ToString" 的 <summary> 文档]

**命名空间：**
 [HslCommunication.ModBus](da9ac567-746e-47e5-70e1-b44577d09ce4.htm)  
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

[缺少 "M:HslCommunication.ModBus.ModbusAscii.ToString" 的 <returns> 文档]

![](../icons/SectionExpanded.png)参见

#### 引用

[ModbusAscii 类](873a8425-d84d-1a6f-ed5c-79e8d056cc76.htm)

[HslCommunication.ModBus 命名空间](da9ac567-746e-47e5-70e1-b44577d09ce4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## UnpackResponseContent 方法 

[原文連結](http://api.hslcommunication.cn/html/00a9f9a8-e32c-331f-cddd-d6d676a819f3.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.ModBus](../html/da9ac567-746e-47e5-70e1-b44577d09ce4.htm "HslCommunication.ModBus")

[ModbusAscii 类](../html/873a8425-d84d-1a6f-ed5c-79e8d056cc76.htm "ModbusAscii 类")

[ModbusAscii 方法](../html/7819f1d9-ca1b-2c0b-6ae3-9a83dfa03b21.htm "ModbusAscii 方法")

[GetNewNetMessage 方法](../html/88dbd66f-58bf-6e1b-c2ab-34c111364480.htm "GetNewNetMessage 方法 ")

[PackCommandWithHeader 方法](../html/73a56f28-e750-b76d-3ebd-bd5e55ad9aae.htm "PackCommandWithHeader 方法 ")

[ToString 方法](../html/9e0bb9d9-c4ac-4e81-3d89-f9eb55d88718.htm "ToString 方法 ")

[UnpackResponseContent 方法](../html/00a9f9a8-e32c-331f-cddd-d6d676a819f3.htm "UnpackResponseContent 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ModbusAsciiUnpackResponseContent 方法 |

根据对方返回的报文命令，对命令进行基本的拆包，例如各种Modbus协议拆包为统一的核心报文，还支持对报文的验证  
According to the message command returned by the other party, the command is basically unpacked, for example,
various Modbus protocols are unpacked into a unified core message, and the verification of the message is also supported

**命名空间：**
 [HslCommunication.ModBus](da9ac567-746e-47e5-70e1-b44577d09ce4.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public override OperateResult<byte[]> UnpackResponseContent(
	byte[] send,
	byte[] response
)
```

```
Public Overrides Function UnpackResponseContent ( 
	send As Byte(),
	response As Byte()
) As OperateResult(Of Byte())
```

```
public:
virtual OperateResult<array<unsigned char>^>^ UnpackResponseContent(
	array<unsigned char>^ send, 
	array<unsigned char>^ response
) override
```

```
abstract UnpackResponseContent : 
        send : byte[] * 
        response : byte[] -> OperateResult<byte[]> 
override UnpackResponseContent : 
        send : byte[] * 
        response : byte[] -> OperateResult<byte[]>
```

#### 参数

send
:   类型：SystemByte  
    发送的原始报文数据

response
:   类型：SystemByte  
    设备方反馈的原始报文内容

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
返回拆包之后的报文信息，默认不进行任何的拆包操作

![](../icons/SectionExpanded.png)备注

在实际解包的操作过程中，通常对状态码，错误码等消息进行判断，如果校验不通过，将携带错误消息返回  
During the actual unpacking operation, the status code, error code and other messages are usually judged. If the verification fails, the error message will be returned.

![](../icons/SectionExpanded.png)参见

#### 引用

[ModbusAscii 类](873a8425-d84d-1a6f-ed5c-79e8d056cc76.htm)

[HslCommunication.ModBus 命名空间](da9ac567-746e-47e5-70e1-b44577d09ce4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ModbusAscii 字段

[原文連結](http://api.hslcommunication.cn/html/62529de8-1a77-85ed-6b9c-9167bdfd5459.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.ModBus](../html/da9ac567-746e-47e5-70e1-b44577d09ce4.htm "HslCommunication.ModBus")

[ModbusAscii 类](../html/873a8425-d84d-1a6f-ed5c-79e8d056cc76.htm "ModbusAscii 类")

[ModbusAscii 构造函数](../html/223428e7-bc23-5a0f-0fa6-ac355d422c71.htm "ModbusAscii 构造函数 ")

[ModbusAscii 属性](../html/fde6aaa6-18a7-1888-9bf8-dd886a14100b.htm "ModbusAscii 属性")

[ModbusAscii 方法](../html/7819f1d9-ca1b-2c0b-6ae3-9a83dfa03b21.htm "ModbusAscii 方法")

[ModbusAscii 字段](../html/62529de8-1a77-85ed-6b9c-9167bdfd5459.htm "ModbusAscii 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ModbusAscii 字段 |

[ModbusAscii](873a8425-d84d-1a6f-ed5c-79e8d056cc76.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)字段

|  | 名称 | 说明 |
| --- | --- | --- |
| 受保护的字段 | [LogMsgFormatBinary](2ed2ff87-caec-0e44-f08b-e2b015225a58.htm) | 设置日志记录报文是否二进制，如果为False，那就使用ASCII码  Set whether the log message is binary, if it is False, then use ASCII code (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[ModbusAscii 类](873a8425-d84d-1a6f-ed5c-79e8d056cc76.htm)

[HslCommunication.ModBus 命名空间](da9ac567-746e-47e5-70e1-b44577d09ce4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ModbusAsciiOverTcp 类

[原文連結](http://api.hslcommunication.cn/html/619c4d60-2614-da3c-4d4e-901324d3b881.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.ModBus](../html/da9ac567-746e-47e5-70e1-b44577d09ce4.htm "HslCommunication.ModBus")

[ModbusAsciiOverTcp 类](../html/619c4d60-2614-da3c-4d4e-901324d3b881.htm "ModbusAsciiOverTcp 类")

[ModbusAsciiOverTcp 构造函数](../html/e2016355-b196-96a9-1e3b-8b32562755a5.htm "ModbusAsciiOverTcp 构造函数 ")

[ModbusAsciiOverTcp 属性](../html/849553d2-15a1-27f1-45c7-b171067a053c.htm "ModbusAsciiOverTcp 属性")

[ModbusAsciiOverTcp 方法](../html/81e01ae3-be89-0d47-cea4-ff6e6eac5c1b.htm "ModbusAsciiOverTcp 方法")

[ModbusAsciiOverTcp 字段](../html/ed438a7a-bfa7-b597-d1e0-05b072f1dcdd.htm "ModbusAsciiOverTcp 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ModbusAsciiOverTcp 类 |

Modbus-Ascii通讯协议的网口透传类，基于rtu类库完善过来，支持标准的功能码，也支持扩展的功能码实现，地址采用富文本的形式，详细见备注说明  
The client communication class of Modbus-Ascii protocol is convenient for data interaction with the server. It supports standard function codes and also supports extended function codes.
The address is in rich text. For details, see the remarks.

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  [HslCommunication.Core.NetBinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)  
    [HslCommunication.Core.DeviceDeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)  
      [HslCommunication.Core.DeviceDeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)  
        [HslCommunication.ModBusModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)  
          HslCommunication.ModBusModbusAsciiOverTcp  
            [HslCommunication.Profinet.DeltaDeltaSerialAsciiOverTcp](52cafdff-0770-ecab-c904-3feb59cd6573.htm)

**命名空间：**
 [HslCommunication.ModBus](da9ac567-746e-47e5-70e1-b44577d09ce4.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class ModbusAsciiOverTcp : ModbusRtuOverTcp
```

```
Public Class ModbusAsciiOverTcp
	Inherits ModbusRtuOverTcp
```

```
public ref class ModbusAsciiOverTcp : public ModbusRtuOverTcp
```

```
type ModbusAsciiOverTcp =  
    class
        inherit ModbusRtuOverTcp
    end
```

ModbusAsciiOverTcp 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [ModbusAsciiOverTcp](8a82007a-e21d-0722-935c-bb55f3f06a01.htm) | 实例化一个Modbus-ascii协议的客户端对象  Instantiate a client object of the Modbus-ascii protocol |
| 公共方法 | [ModbusAsciiOverTcp(String, Int32, Byte)](57ac5644-74b3-3db0-83c4-5abbdd906a17.htm) | 指定服务器地址，端口号，客户端自己的站号来初始化  Specify the server address, port number, and client's own station number to initialize |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [AddressStartWithZero](f4122661-5369-da8f-7c8c-84c85d70a480.htm) | 获取或设置起始的地址是否从0开始，默认为True  Gets or sets whether the starting address starts from 0. The default is True (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共属性 | [BroadcastStation](f307b9fd-b01f-cb79-bb97-0d1d5e8c7d29.htm) | 获取或是设置当前广播模式对应的站号，广播模式意味着不接收设备方的数据返回操作，默认为 -1，表示不使用广播模式。  Gets or sets the station number corresponding to the current broadcast mode. Broadcast mode means that the data return operation of the device is not received. The default value is -1, indicating that broadcast mode is not used. (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共属性代码示例 | [ByteTransform](f35118c9-3691-f752-ae6f-2c5549714faa.htm) | 当前的数据变换机制，当你需要从字节数据转换类型数据的时候需要。  The current data transformation mechanism is required when you need to convert type data from byte data. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共属性 | [CommunicationPipe](ad8e903b-d289-21b0-077e-15ede8005221.htm) | 获取或设置当前的管道信息，管道类型为[CommunicationPipe](ad8e903b-d289-21b0-077e-15ede8005221.htm)的继承类，内置了[PipeTcpNet](26b94bd5-2695-36e7-6378-c3d41846be8f.htm)管道，[PipeUdpNet](5a0b810e-187e-826b-0a13-34fd5576aa0e.htm)管道，[PipeSerialPort](9d924607-3dfc-3916-b532-c5dec48fe347.htm)管道等  Get or set the current pipeline information, the pipeline type is [CommunicationPipe](ad8e903b-d289-21b0-077e-15ede8005221.htm) inheritance class, [PipeTcpNet](26b94bd5-2695-36e7-6378-c3d41846be8f.htm) pipeline, [PipeUdpNet](5a0b810e-187e-826b-0a13-34fd5576aa0e.htm) pipeline, [PipeSerialPort](9d924607-3dfc-3916-b532-c5dec48fe347.htm) pipeline, etc (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [ConnectionId](ce1efb4b-fbc9-4683-163e-0a8e71e99d6b.htm) | 当前连接的唯一ID号，默认为长度20的guid码加随机数组成，方便列表管理，也可以自己指定  The unique ID number of the current connection. The default is a 20-digit guid code plus a random number. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性代码示例 | [ConnectTimeOut](71a651e5-737d-6e88-75a2-fb891b9af6cd.htm) | 获取或设置接收服务器反馈的时间，如果为负数，则不接收反馈   Gets or sets the time to receive server feedback, and if it is a negative number, does not receive feedback (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共属性 | [Crc16CheckEnable](fe1f096a-5f06-b612-7dab-4e56be636a76.htm) | 获取或设置是否启用CRC16校验码的检查功能，默认启用，如果需要忽略检查CRC16，则设置为 false 即可。  Gets or sets whether to enable the check function of CRC16 check code. It is enabled by default. If you need to ignore the check of CRC16, you can set it to false. (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共属性 | [DataFormat](ee89a24c-d5b0-82e0-9c40-9cc4ab57f2b7.htm) | 获取或设置数据解析的格式，可选ABCD, BADC，CDAB，DCBA格式，对int,uint,float,double,long,ulong类型有作用  Get or set the format of the data analysis, optional ABCD, BADC, CDAB, DCBA format, effective for int, uint, float, double, long, ulong type (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共属性 | [DisableFunctionCode06](a1be5198-fdfa-3e32-32ed-0edbdaf577ce.htm) | 获取或设置是否禁用功能码05的写入操作，设置为true时，表示不允许使用功能码 06 进行写入操作。写入单个的字数据也使用 0x10 功能码  Get or set whether to disable the write operation of function code 06. When set to true, it indicates that the write operation of function code 05 is not allowed. Writing individual character data also uses the 0x10 function code (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共属性 | [EnableWriteMaskCode](358ac1fd-3000-898e-1c63-6213649d851a.htm) | 获取或设置当前掩码写入的功能码是否激活状态，设置为 false 时，再执行写入位时，会通过读字，修改位，写字的方式来间接实现。  When the function code is set to false, and then the write bit is executed, it will be indirectly implemented by reading, modifying the bit, and writing the word. (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共属性代码示例 | [IpAddress](c378f4e8-3fdf-1932-7239-ed44bc5af584.htm) | 获取或是设置远程服务器的IP地址，如果是本机测试，那么需要设置为127.0.0.1   Get or set the IP address of the remote server. If it is a local test, then it needs to be set to 127.0.0.1 (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共属性 | [IsStringReverse](2b20320d-43f0-d660-a603-580f923e9e43.htm) | 字符串数据是否按照字来反转，默认为False  Whether the string data is reversed according to words. The default is False. (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共属性 | [LocalBinding](622b418f-5969-3726-5dcc-04456c4b97a7.htm) | 获取或设置绑定的本地的IP地址和端口号信息，如果端口设置为0，代表任何可用的端口  Get or set the bound local IP address and port number information, if the port is set to 0, it means any available port (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共属性代码示例 | [LogNet](8c56f246-88ac-afb8-8c06-972fafe4f78a.htm) | 组件的日志工具，支持日志记录，只要实例化后，当前网络的基本信息，就以[DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm)等级进行输出  The component's logging tool supports logging. As long as the instantiation of the basic network information, the output will be output at [DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm) (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性代码示例 | [Port](6f1ad24a-c916-5b67-aa17-a8c7dc2acb2e.htm) | 获取或设置服务器的端口号，具体的值需要取决于对方的配置  Gets or sets the port number of the server. The specific value depends on the configuration of the other party. (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共属性代码示例 | [ReceiveTimeOut](250560a7-332b-ce1b-5eb0-e3abe4ee6012.htm) | 获取或设置接收服务器反馈的时间，如果为负数，则不接收反馈   Gets or sets the time to receive server feedback, and if it is a negative number, does not receive feedback (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [SendBeforeHex](72aa1de4-116f-abfe-4adb-814a0be0e560.htm) | 获取或设置在发送通信报文前追加发送的字节信息，HEX格式，通常用于lora组网时，需要携带 00 00 00 02 四个字节的站地址功能。  Obtain or set the byte information sent before sending communication packets, HEX format, usually used for LORA networking, you need to carry 00 00 00 02 four-byte station address function. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [SleepTime](86a5fbb9-3456-6beb-7113-5bbad5eabbc5.htm) | 获取或设置在正式接收对方返回数据前的时候，需要休息的时间，当设置为0的时候，不需要休息。  Get or set the time required to rest before officially receiving the data from the other party. When it is set to 0, no rest is required. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [SocketKeepAliveTime](96ab7703-8c4d-b4c2-ddd9-0dc93ebcfc8d.htm) | 获取或设置客户端的Socket的心跳时间信息，这个是Socket底层自动实现的心跳包，不基于协议层实现。默认小于0，不开启心跳检测，如果需要开启，设置 60\_000 比较合适，单位毫秒  Get or set the heartbeat time information of the Socket of the client. This is the heartbeat packet automatically implemented by the bottom layer of the Socket, not based on the protocol layer. The default value is less than 0, and heartbeat detection is not enabled. If you need to enable it, it is more appropriate to set 60\_000, in milliseconds. (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共属性 | [Station](9e7022b4-45cc-58e8-8606-f9866a9e4e66.htm) | 获取或者重新修改服务器的默认站号信息，当然，你可以再读写的时候动态指定，例如地址 "s=2;100", 更详细的例子可以参考DEMO界面上的地址示例。  Get or re-modify the server's default station number information, of course, you can dynamically specify when reading and writing, such as the address "s=2; 100", a more detailed example can refer to the DEMO interface address example. (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共属性 | [StationCheckMacth](861ed947-785a-3866-afd5-ec8042e7bebc.htm) | 获取或设置当前是否启用站号检查的功能，默认启用，读写数据时将进行站号确认操作，如果需要忽略站号，则设置为 false 即可。  Gets or sets whether the station ID check function is enabled. It is enabled by default. When reading or writing data, the station ID is confirmed. (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 受保护的属性 | [WordLength](f1a505cc-e85c-6c8b-b189-58a1ab5753f9.htm) | 一个字单位的数据表示的地址长度，西门子为2，三菱，欧姆龙，modbusTcp就为1，AB PLC无效  The address length represented by one word of data, Siemens is 2, Mitsubishi, Omron, modbusTcp is 1, AB PLC is invalid (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共属性 | [WordReadBatchLength](1d4b2647-54fd-22b6-dddf-e88356c270a5.htm) | 当进行字读取的时候，商业授权用户支持分批读取，本属性表示分批的长度信息，默认为 120，可以根据实际情况来设置其他值。  When reading of word type, commercial authorization users support batch reading. This attribute represents the length information of the batch, with a default of 120. Other values can be set according to the actual situation. (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法代码示例 | [ConnectClose](569a9c0e-7628-b8d5-e124-0853452747b8.htm) | 手动断开与远程服务器的连接，如果当前是长连接模式，那么就会切换到短连接模式  Manually disconnect from the remote server, if it is currently in long connection mode, it will switch to short connection mode (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共方法代码示例 | [ConnectCloseAsync](76be8150-ff22-95b9-16b1-7c3440c54ae8.htm) | 手动断开与远程服务器的连接，如果当前是长连接模式，那么就会切换到短连接模式  Manually disconnect from the remote server, if it is currently in long connection mode, it will switch to short connection mode (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共方法代码示例 | [ConnectServer](6299cc1f-2fab-b754-4597-1eddc2b21308.htm) | 尝试连接远程的服务器，如果连接成功，就切换短连接模式到长连接模式，后面的每次请求都共享一个通道，使得通讯速度更快速  Try to connect to a remote server. If the connection is successful, switch the short connection mode to the long connection mode. Each subsequent request will share a channel, making the communication speed faster. (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共方法代码示例 | [ConnectServerAsync](e7eff5c0-8e2a-3b6c-5027-16f4d8bcf5ae.htm) | 尝试连接远程的服务器，如果连接成功，就切换短连接模式到长连接模式，后面的每次请求都共享一个通道，使得通讯速度更快速  Try to connect to a remote server. If the connection is successful, switch the short connection mode to the long connection mode. Each subsequent request will share a channel, making the communication speed faster. (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
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
| 受保护的方法 | [GetNewNetMessage](a95198f1-fac7-8c00-05af-cd30b0b44299.htm) | 获取一个新的消息对象的方法，需要在继承类里面进行重写  The method to get a new message object needs to be overridden in the inheritance class (重写 [ModbusRtuOverTcpGetNewNetMessage](c08212c0-06b8-a8c2-2fa9-b1b36ef1212d.htm).) |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法 | [GetWordLength](c27821ea-1739-8eb9-a848-23a00828fec8.htm) | 一个字单位的数据表示的地址长度，西门子为2，三菱，欧姆龙，modbusTcp就为1，AB PLC无效  The address length represented by one word of data, Siemens is 2, Mitsubishi, Omron, modbusTcp is 1, AB PLC is invalid (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 受保护的方法代码示例 | [InitializationOnConnect](4ab6989f-a5b5-05fe-8b44-bfa45824a2a7.htm) | 根据实际的协议选择是否重写本方法，有些协议在创建连接之后，需要进行一些初始化的信号握手，才能最终建立网络通道。  Whether to rewrite this method is based on the actual protocol. Some protocols require some initial signal handshake to establish a network channel after the connection is created. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法代码示例 | [InitializationOnConnectAsync](7da2af04-d4fc-5045-19a4-4b0d1722eef3.htm) | 根据实际的协议选择是否重写本方法，有些协议在创建连接之后，需要进行一些初始化的信号握手，才能最终建立网络通道。  Whether to rewrite this method is based on the actual protocol. Some protocols require some initial signal handshake to establish a network channel after the connection is created. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [IpAddressPing](0ac126b4-a830-a822-4ef5-855474942cd7.htm) | 对当前设备的IP地址进行PING的操作，返回PING的结果，正常来说，返回Success  PING the IP address of the current device and return the PING result. Normally, it returns Success (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 受保护的方法 | [LogRevcMessage(Byte)](dab1657f-1e75-7fe1-7003-ccedfc69855d.htm) | 使用日志记录一个接收的报文信息  Logs are used to record information about a received packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [LogRevcMessage(Byte, PipeSession)](0f51ae45-575e-4228-f978-01b37d6d5df1.htm) | 使用日志记录一个接收的报文信息  Logs are used to record information about a received packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [LogSendMessage(Byte)](4f90eecc-68a2-1ef2-44ef-369f19bed3ec.htm) | 使用日志记录一个发送的报文信息  Logs are used to record information about a send packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [LogSendMessage(Byte, PipeSession)](209caadf-4f0f-40aa-04c2-ebc61f8dab13.htm) | 使用日志记录一个发送的报文信息  Logs are used to record information about a send packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法 | [PackCommandWithHeader](bd949a79-bda2-9cea-ab5a-76ff90f7258d.htm) | 对当前的命令进行打包处理，通常是携带命令头内容，标记当前的命令的长度信息，需要进行重写，否则默认不打包  The current command is packaged, usually carrying the content of the command header, marking the length of the current command, and it needs to be rewritten, otherwise it is not packaged by default (重写 [ModbusRtuOverTcpPackCommandWithHeader(Byte)](ca0c7547-4769-2e5c-b3e5-97e68f1bb00b.htm).) |
| 公共方法代码示例 | [Read(String, UInt16)](a061c449-5af4-7b69-e6ed-76f1740e178f.htm) | 从Modbus服务器批量读取寄存器的信息，需要指定起始地址，读取长度，如果富文本地址不指定，默认使用的功能码是 0x03，如果需要使用04功能码，那么地址就写成 x=4;100  To read the register information from the Modbus server in batches, you need to specify the start address and read length. If the rich text address is not specified, the default function code is 0x03. If you need to use the 04 function code, the address is written as x = 4; 100 (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法代码示例 | [ReadT](95abe064-5594-5ae7-a849-98c53f88fee0.htm) | 读取支持Hsl特性的数据内容，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明，或是官网的详细文档  Read the data content of the Hsl attribute. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadAsync(String, UInt16)](fb1da34a-99cc-17bc-6e2a-e5155a6b9015.htm) | 从Modbus服务器批量读取寄存器的信息，需要指定起始地址，读取长度，如果富文本地址不指定，默认使用的功能码是 0x03，如果需要使用04功能码，那么地址就写成 x=4;100  To read the register information from the Modbus server in batches, you need to specify the start address and read length. If the rich text address is not specified, the default function code is 0x03. If you need to use the 04 function code, the address is written as x = 4; 100 (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法代码示例 | [ReadAsyncT](0522e740-4bf0-2e31-c4c5-e41384321b37.htm) | 异步读取支持Hsl特性的数据内容，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Asynchronously read the data content of the Hsl attribute. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadBool(String)](9b2faa2e-eaa9-7ace-c49f-7a0318f1df1a.htm) | 读取单个的Boolean数据信息  Read a single Boolean data message (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadBool(String, UInt16)](b9d9664e-4a71-41e9-dd8a-bc6167dcac6c.htm) | 批量读取线圈或是离散的数据信息，需要指定地址和长度，具体的结果取决于实现，如果富文本地址不指定，默认使用的功能码是 0x01  To read coils or discrete data in batches, you need to specify the address and length. The specific result depends on the implementation. If the rich text address is not specified, the default function code is 0x01. (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法 | [ReadBoolAsync(String)](cc37da1f-d673-fe67-78ee-4e465d29a2bd.htm) | 异步读取单个的Boolean数据信息  Asynchronously read a single Boolean data message (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadBoolAsync(String, UInt16)](4ab1379f-3369-4060-6837-413f93bc424a.htm) | 批量读取线圈或是离散的数据信息，需要指定地址和长度，具体的结果取决于实现，如果富文本地址不指定，默认使用的功能码是 0x01  To read coils or discrete data in batches, you need to specify the address and length. The specific result depends on the implementation. If the rich text address is not specified, the default function code is 0x01. (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法 | [ReadCoil(String)](9b750e8c-d4e3-2200-8536-7fba9a12b819.htm) | 读取线圈，需要指定起始地址，如果富文本地址不指定，默认使用的功能码是 0x01  To read the coil, you need to specify the start address. If the rich text address is not specified, the default function code is 0x01. (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法 | [ReadCoil(String, UInt16)](4633c27f-4e6f-7ed2-88df-cf6b88256496.htm) | 批量的读取线圈，需要指定起始地址，读取长度，如果富文本地址不指定，默认使用的功能码是 0x01  For batch reading coils, you need to specify the start address and read length. If the rich text address is not specified, the default function code is 0x01. (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法 | [ReadCoilAsync(String)](f744db67-b882-392c-1824-71070c6d19a4.htm) | 读取线圈，需要指定起始地址，如果富文本地址不指定，默认使用的功能码是 0x01  To read the coil, you need to specify the start address. If the rich text address is not specified, the default function code is 0x01. (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法 | [ReadCoilAsync(String, UInt16)](2e8b3480-8188-ed91-cca1-7e53cc641ca2.htm) | 批量的读取线圈，需要指定起始地址，读取长度，如果富文本地址不指定，默认使用的功能码是 0x01  For batch reading coils, you need to specify the start address and read length. If the rich text address is not specified, the default function code is 0x01. (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法代码示例 | [ReadCustomerT(String)](ca069688-d8f0-7dbd-e962-86a10dce83cf.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，返回一个新的类型的实例对象。  To read a custom data type, you need to inherit from the IDataTransfer interface and return an instance object of a new type. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadCustomerT(String, T)](ca7dc1a1-c97d-689f-9819-4dae4af67206.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，传入一个实例，对这个实例进行赋值，并返回该实例的对象。  To read a custom data type, you need to inherit from the IDataTransfer interface, pass in an instance, assign a value to this instance, and return the object of the instance. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadCustomerAsyncT(String)](f3fcda2b-255d-aadd-567e-c7b3df3636d9.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，返回一个新的类型的实例对象。  To read a custom data type, you need to inherit from the IDataTransfer interface and return an instance object of a new type. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadCustomerAsyncT(String, T)](527c9fae-4743-6ebd-7b9a-5ef3a8b7c455.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，传入一个实例，对这个实例进行赋值，并返回该实例的对象。  To read a custom data type, you need to inherit from the IDataTransfer interface, pass in an instance, assign a value to this instance, and return the object of the instance. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadDiscrete(String)](4886013f-4945-b8b5-7a85-8fd979367a5b.htm) | 读取输入线圈，需要指定起始地址，如果富文本地址不指定，默认使用的功能码是 0x02  To read the input coil, you need to specify the start address. If the rich text address is not specified, the default function code is 0x02. (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法 | [ReadDiscrete(String, UInt16)](f907fa38-a034-8af3-5b1d-99ccd6062667.htm) | 批量的读取输入点，需要指定起始地址，读取长度，如果富文本地址不指定，默认使用的功能码是 0x02  To read input points in batches, you need to specify the start address and read length. If the rich text address is not specified, the default function code is 0x02 (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法 | [ReadDiscreteAsync(String)](abc82d59-8329-5160-b68d-e19588fe4a20.htm) | 读取输入线圈，需要指定起始地址，如果富文本地址不指定，默认使用的功能码是 0x02  To read the input coil, you need to specify the start address. If the rich text address is not specified, the default function code is 0x02. (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法 | [ReadDiscreteAsync(String, UInt16)](70c766bd-2a88-287f-7c06-c2da7870f010.htm) | 批量的读取输入点，需要指定起始地址，读取长度，如果富文本地址不指定，默认使用的功能码是 0x02  To read input points in batches, you need to specify the start address and read length. If the rich text address is not specified, the default function code is 0x02 (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法代码示例 | [ReadDouble(String)](e7e31c43-71b4-e692-6bdc-ba9533e1d83d.htm) | 读取双浮点的数据  Read double floating point data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadDouble(String, UInt16)](98be3669-7b4b-fa15-efd7-12cdda2a1ac2.htm) | 读取双浮点数据的数组  Read double floating point data array (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法代码示例 | [ReadDoubleAsync(String)](2cf895fa-1033-1def-f048-0e8caefaa71d.htm) | 异步读取双浮点的数据  Asynchronously read double floating point data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadDoubleAsync(String, UInt16)](82fc6208-2830-7463-9003-fd7a33ea6dab.htm) | 异步读取双浮点数据的数组  Asynchronously read double floating point data array (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法 | [ReadFile](6ff35712-0972-5094-2aee-05f62228db24.htm) | 使用0x14功能码读取文件的数据内容，需要指定文件编号，地址，长度，返回读取的结果数据，读取的长度为任意长度，内部自动根据实际情况来切割操作。 (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法代码示例 | [ReadFloat(String)](bc2801ae-6600-e845-7c6c-522c9156eb0e.htm) | 读取单浮点数据  Read single floating point data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadFloat(String, UInt16)](de69a6c8-ee84-464e-d735-aa479a2fbe4f.htm) | 读取单浮点精度的数组  Read single floating point array (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法代码示例 | [ReadFloatAsync(String)](89f743f8-4b3a-7085-1e49-a1d7adb10052.htm) | 异步读取单浮点数据  Asynchronously read single floating point data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadFloatAsync(String, UInt16)](ffac5e66-8ba9-5e95-b6f9-21ac11013093.htm) | 异步读取单浮点精度的数组  Asynchronously read single floating point array (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法 | [ReadFromCoreServer(IEnumerableByte)](044cfdd9-ec76-b89f-2f2b-5f439025006c.htm) | 将多个数据报文按顺序发到设备，并从设备接收返回的数据内容，然后拼接成一个Byte[]信息，需要重写[UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm)方法才能返回正确的结果。  Send multiple data packets to the device in sequence, and receive the returned data content from the device, and then splicing them into a Byte[] message, you need to rewrite [UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm) method to return the correct result. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServer(Byte)](b516dde4-8446-eaaa-f60c-35500a2c9a25.htm) | 将二进制的数据发送到管道中去，然后从管道里接收二进制的数据回来，并返回是否成功的结果对象。  Send binary data to the pipeline, and then receive binary data back from the pipeline, and return whether the success of the result object (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServer(Byte, Boolean, Boolean)](b32ce8e8-a959-fcc2-6736-14b39fdb9a40.htm) | 将二进制的数据发送到管道中去，然后从管道里接收二进制的数据回来，并返回是否成功的结果对象。  Send binary data to the pipeline, and then receive binary data back from the pipeline, and return whether the success of the result object (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServer(CommunicationPipe, Byte, Boolean, Boolean)](ef36d56b-df09-f727-e8d3-316d37a13ee9.htm) | 使用指定的管道来进行数据通信，发送原始数据到管道，然后从管道接收相关的数据返回，本方法无锁 (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServerAsync(IEnumerableByte)](bcdb59b7-a70a-d5c0-1fc6-9ececcd39ea2.htm) | 将多个数据报文按顺序发到设备，并从设备接收返回的数据内容，然后拼接成一个Byte[]信息，需要重写[UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm)方法才能返回正确的结果。  Send multiple data packets to the device in sequence, and receive the returned data content from the device, and then splicing them into a Byte[] message, you need to rewrite [UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm) method to return the correct result. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServerAsync(Byte)](90b7ccb6-1139-1b6f-edd1-d5d0b845eb6c.htm) | 将二进制的数据发送到管道中去，然后从管道里接收二进制的数据回来，并返回是否成功的结果对象。  Send binary data to the pipeline, and then receive binary data back from the pipeline, and return whether the success of the result object (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServerAsync(Byte, Boolean, Boolean)](1c79c204-5ab5-9342-e65d-7455c20ea577.htm) | 将二进制的数据发送到管道中去，然后从管道里接收二进制的数据回来，并返回是否成功的结果对象。  Send binary data to the pipeline, and then receive binary data back from the pipeline, and return whether the success of the result object (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServerAsync(CommunicationPipe, Byte, Boolean, Boolean)](a254a19d-3c29-147b-335b-8adccdf4ee2f.htm) | 使用指定的管道来进行数据通信，发送原始数据到管道，然后从管道接收相关的数据返回，本方法无锁 (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法代码示例 | [ReadInt16(String)](57a16003-2fd5-b2bd-09bf-14366be899b8.htm) | 读取16位的有符号的整型数据  Read 16-bit signed integer data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt16(String, UInt16)](b2bde158-f51c-93fc-66ca-5e4cf6ea1fb7.htm) | 读取16位的有符号整型数组  Read 16-bit signed integer array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt16Async(String)](ac021793-890d-ceab-eae4-47ff763390c1.htm) | 异步读取16位的有符号的整型数据  Asynchronously read 16-bit signed integer data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt16Async(String, UInt16)](f30e02fc-8078-f1b3-c226-5991e3e23caf.htm) | 异步读取16位的有符号整型数组  Asynchronously read 16-bit signed integer array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt32(String)](bfbea035-9a83-40e3-084c-28122041b37a.htm) | 读取32位的有符号整型  Read 32-bit signed integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt32(String, UInt16)](44156b71-97d1-917a-782a-5e716ffd286c.htm) | 读取32位有符号整型数组  Read 32-bit signed integer array (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法代码示例 | [ReadInt32Async(String)](96726ea0-82aa-1229-21f6-43cd62abca43.htm) | 异步读取32位的有符号整型  Asynchronously read 32-bit signed integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt32Async(String, UInt16)](4957686b-6205-802f-177e-b5322e1eabda.htm) | 异步读取32位有符号整型数组  Asynchronously read 32-bit signed integer array (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法代码示例 | [ReadInt64(String)](0d0c2284-2da1-6f9a-b353-537561cd5e95.htm) | 读取64位的有符号整型  Read 64-bit signed integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt64(String, UInt16)](2a1dd46e-c6be-563e-6ae7-86d7432b64f6.htm) | 读取64位的有符号整型数组  Read 64-bit signed integer array (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法代码示例 | [ReadInt64Async(String)](26d46c8a-3775-3abd-d3ec-5cd1758a317e.htm) | 异步读取64位的有符号整型  Asynchronously read 64-bit signed integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt64Async(String, UInt16)](5a7ebc54-70a1-50a1-6c3a-9c4b3c4f9af2.htm) | 异步读取64位的有符号整型数组  Asynchronously read 64-bit signed integer array (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
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
| 公共方法代码示例 | [ReadUInt32(String, UInt16)](c5543613-3abd-2cb0-dc1d-eb89b3da3a1e.htm) | 读取32位的无符号整型数组  Read 32-bit unsigned integer array (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法代码示例 | [ReadUInt32Async(String)](2343953f-6350-331e-6ba7-18a4449f0575.htm) | 异步读取32位的无符号整型  Asynchronously read 32-bit unsigned integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt32Async(String, UInt16)](272c5014-a400-5e74-69f5-1de1237e4ae0.htm) | 异步读取32位的无符号整型数组  Asynchronously read 32-bit unsigned integer array (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法代码示例 | [ReadUInt64(String)](947487dd-d4a3-72d8-d9fb-2224aad70e12.htm) | 读取64位的无符号整型  Read 64-bit unsigned integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt64(String, UInt16)](f3a9d180-d7bc-aacc-b07d-0edc4608547c.htm) | 读取64位的无符号整型的数组  Read 64-bit unsigned integer array (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法代码示例 | [ReadUInt64Async(String)](3411ff3a-772f-dbe1-8931-75e00a731579.htm) | 异步读取64位的无符号整型  Asynchronously read 64-bit unsigned integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt64Async(String, UInt16)](20fab556-ee24-9f7b-9a2c-dfdbd2b2dca4.htm) | 异步读取64位的无符号整型的数组  Asynchronously read 64-bit unsigned integer array (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法 | [ReadWrite](4d80e7e8-42a3-5005-e8f0-5dd028431ccd.htm) | 使用0x17功能码来实现同时写入并读取数据的操作，使用一条报文来实现，需要指定读取的地址，长度，写入的地址，写入的数据信息，返回读取的结果数据。  Use 0x17 function code to write and read data at the same time, and use a message to implement it, you need to specify the read address, length, written address, written data information, and return the read result data. (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法 | [RegisteredAddressMapping](5686f821-9dcb-c44c-5a99-663d0b75748a.htm) | 注册一个新的地址映射关系，注册地址映射关系后，就可以使用新的地址来读写Modbus数据了，通常用于其他的支持Modbus协议的PLC。  After registering a new address mapping, you can use the new address to read and write Modbus data, which is usually used for other PLCs that support the Modbus protocol. (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法 | [SetDtuPipe](f765d1d7-234a-b229-9457-97e25e024dc9.htm) | 将当前的通信对象设置DTU模式，允许传入现成的管道，并返回初始化结果，如果该设备重写了握手报文，就是返回握手结果  Set the current communication object to DTU mode, allow the existing pipe to be passed in, and return the initialization result, if the device rewrites the handshake packet, the handshake result is returned (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [SetDtuPipeAsync](9d1fd12d-b804-3828-a627-a2c0eef1d23f.htm) | 将当前的通信对象设置DTU模式，允许传入现成的管道，并返回初始化结果，如果该设备重写了握手报文，就是返回握手结果  Set the current communication object to DTU mode, allow the existing pipe to be passed in, and return the initialization result, if the device rewrites the handshake packet, the handshake result is returned (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [SetPersistentConnection](f903dae0-4a4e-c7db-8350-50e36b257ec5.htm) | **已过时。** V11版本及之前设置长连接的方法，在V12版本以上中没有任何效果，默认长连接，删除调用即可，此处保留方法是为了部分用户保持兼容性升级。  The method of setting the long connection in V11 and before, has no effect in V12 and above. this method can be deleted. The method is retained here to maintain compatibility upgrades for some users. (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共方法 | [ToString](fcffa0b7-3ff8-42db-f970-82cacc1eb62a.htm) | (重写 [ModbusRtuOverTcpToString](c409cc5d-0b8d-7cab-dc72-b323757b8449.htm).) |
| 公共方法 | [TranslateToModbusAddress](fd4a9aee-b09c-08e3-dff3-51a3dbdd6087.htm) | 将当前的地址信息转换成Modbus格式的地址，如果转换失败，返回失败的消息。默认不进行任何的转换。  Convert the current address information into a Modbus format address. If the conversion fails, a failure message will be returned. No conversion is performed by default. (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法 | [UnpackResponseContent](d7cb5512-8c9d-3247-0d89-c89cc3011799.htm) | 根据对方返回的报文命令，对命令进行基本的拆包，例如各种Modbus协议拆包为统一的核心报文，还支持对报文的验证  According to the message command returned by the other party, the command is basically unpacked, for example, various Modbus protocols are unpacked into a unified core message, and the verification of the message is also supported (重写 [ModbusRtuOverTcpUnpackResponseContent(Byte, Byte)](4dc1fa2e-6465-742e-bf0f-f133a6686523.htm).) |
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
| 公共方法代码示例 | [Write(String, Double)](64c44576-7e95-cbcb-a8c9-e6b7ff7e3c50.htm) | 写入double数据，返回是否成功  Write double data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Int16)](c7cfd67d-ca31-3087-c684-ee4c10780ab1.htm) | 写入short数组，返回是否成功  Write short array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Int32)](ba80d1a5-c8d8-7e0b-5c6b-9435cf20cadb.htm) | 写入int数据，返回是否成功  Write int data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Int64)](ba246b33-5a12-a22c-6d40-0969a64b508b.htm) | 写入long数据，返回是否成功  Write long data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Single)](a2265fa7-5b56-a5aa-9d21-a0dc5458cbf9.htm) | 写入float数据，返回是否成功  Write float data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, String)](4e0ec363-ba64-7ca7-b8d3-2e9a8087bd55.htm) | 写入字符串信息，编码为ASCII  Write string information, encoded as ASCII (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, UInt16)](65544962-82ac-bb47-d72f-5dd0c3d4ac6c.htm) | 写入ushort数组，返回是否成功  Write ushort array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, UInt32)](c99d8246-f95e-4f79-738b-16e622072eb4.htm) | 写入uint数据，返回是否成功  Write uint data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, UInt64)](4924086e-9b36-d446-4ae4-6682134a6ad1.htm) | 写入ulong数据，返回是否成功  Write ulong data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [Write(String, Boolean)](0be65147-ace7-ba53-504d-d4b923a30333.htm) | 向线圈中写入bool数值，返回是否写入成功，如果富文本地址不指定，默认使用的功能码是 0x05， 如果你的地址为字地址，例如100.2，那么将使用0x16的功能码，通过掩码的方式来修改寄存器的某一位，需要Modbus服务器支持，否则自动切换为读取字数据，修改位，在写入字的方式。  Write bool value to the coil and return whether the writing is successful. If the rich text address is not specified, the default function code is 0x05. If your address is a word address, such as 100.2, then you will use the function code of 0x16 to modify a bit of the register through a mask. It needs Modbus server support, Otherwise, it automatically switches to read the word data, modifies the bits, and writes the word in the way. (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法 | [Write(String, Boolean)](4fe07f04-9b89-7291-d2d7-1f474293f4de.htm) | 向线圈中写入bool数组，返回是否写入成功，如果富文本地址不指定，默认使用的功能码是 0x0F  Write the bool array to the coil, and return whether the writing is successful. If the rich text address is not specified, the default function code is 0x0F. (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法代码示例 | [Write(String, Byte)](fdfef281-9b77-397a-d6df-078cd26afd17.htm) | 将数据写入到Modbus的寄存器上去，需要指定起始地址和数据内容，如果富文本地址不指定，默认使用的功能码是 0x10  To write data to Modbus registers, you need to specify the start address and data content. If the rich text address is not specified, the default function code is 0x10 (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法代码示例 | [Write(String, Double)](4d6f1238-0d16-f592-924e-c35cbd755e83.htm) | 写入double数组，返回是否成功  Write double array, return whether the write was successful (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法 | [Write(String, Int16)](77d3ad67-4385-bef7-4223-cc2340afbb7e.htm) | 将数据写入到Modbus的单个寄存器上去，需要指定起始地址和数据值，如果富文本地址不指定，默认使用的功能码是 0x06  To write data to a single register of Modbus, you need to specify the start address and data value. If the rich text address is not specified, the default function code is 0x06. (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法代码示例 | [Write(String, Int32)](ed9597f2-87de-c1ef-4b83-4267527a38c9.htm) | 写入int[]数组，返回是否成功  Write int array, return whether the write was successful (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法代码示例 | [Write(String, Int64)](72aef2fe-26dd-dc40-55df-9fa4aa2ae6bf.htm) | 写入long数组，返回是否成功  Write long array, return whether the write was successful (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法代码示例 | [Write(String, Single)](9e904462-199f-488f-2e9f-3c972ca9955e.htm) | 写入float数组，返回是否成功  Write float array, return whether the write was successful (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法 | [Write(String, UInt16)](a3316efc-b1cc-2b7f-ece1-f945b3b0182d.htm) | 将数据写入到Modbus的单个寄存器上去，需要指定起始地址和数据值，如果富文本地址不指定，默认使用的功能码是 0x06  To write data to a single register of Modbus, you need to specify the start address and data value. If the rich text address is not specified, the default function code is 0x06. (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法代码示例 | [Write(String, UInt32)](3f968111-808b-5681-7126-74765225c918.htm) | 写入uint[]数组，返回是否成功  Write uint array, return whether the write was successful (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法代码示例 | [Write(String, UInt64)](a0df0554-5053-ebed-bfab-5f8a3c331a48.htm) | 写入ulong数组，返回是否成功  Write ulong array, return whether the write was successful (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法代码示例 | [Write(String, String, Int32)](02fe2ca9-49a3-ede4-d1af-8c6da1f644da.htm) | 写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为ASCII  Write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, it is filled with 0 and the encoding is ASCII. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, String, Encoding)](7b602661-902d-2b70-1fdb-ba87ecd307bf.htm) | 写入字符串信息，需要指定的编码信息  Write string information, need to specify the encoding information (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, String, Int32, Encoding)](e8cb30d9-4c5e-6c68-b9af-1b1b6c2dc714.htm) | 写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为指定的编码信息  Write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, then the operation is complemented with 0 , you should specified the encoding information (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteT(T)](9f48853f-fbeb-c37d-25a2-1c76d1dee407.htm) | 写入支持Hsl特性的数据，返回是否写入成功，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Write data that supports the Hsl attribute, and return whether the write was successful. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Double)](068adaa4-d1de-dffc-8376-6f483e2c1d96.htm) | 异步写入double数据，返回是否成功  Asynchronously write double data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int16)](2625e868-7096-e526-ce25-b29afc639b41.htm) | 异步写入short数组，返回是否成功  Asynchronously write short array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int32)](df91d629-2080-328d-67cd-fbefb9a45404.htm) | 异步写入int数据，返回是否成功  Asynchronously write int data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int64)](ec6c1e24-a0d8-f1f5-56af-109b0c546fc8.htm) | 异步写入long数据，返回是否成功  Asynchronously write long data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Single)](910189df-13f9-b7d0-a626-b5b25ce5d32a.htm) | 异步写入float数据，返回是否成功  Asynchronously write float data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, String)](fe839c9a-46f3-e317-0155-b10a8a83a9b1.htm) | 异步写入字符串信息，编码为ASCII  Asynchronously write string information, encoded as ASCII (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt16)](34a1a321-972e-06a8-76fb-794dd1229005.htm) | 异步写入ushort数组，返回是否成功  Asynchronously write ushort array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt32)](abad215c-bc99-ac90-b9f3-7adb3a7cca8f.htm) | 异步写入uint数据，返回是否成功  Asynchronously write uint data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt64)](e396074d-6cf1-b45d-5e8e-01026f6e391b.htm) | 异步写入ulong数据，返回是否成功  Asynchronously write ulong data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WriteAsync(String, Boolean)](62f1700b-e043-2ec6-f2aa-47aa41548cec.htm) | 向线圈中写入bool数值，返回是否写入成功，如果富文本地址不指定，默认使用的功能码是 0x05， 如果你的地址为字地址，例如100.2，那么将使用0x16的功能码，通过掩码的方式来修改寄存器的某一位，需要Modbus服务器支持，否则自动切换为读取字数据，修改位，在写入字的方式。  Write bool value to the coil and return whether the writing is successful. If the rich text address is not specified, the default function code is 0x05. If your address is a word address, such as 100.2, then you will use the function code of 0x16 to modify a bit of the register through a mask. It needs Modbus server support, Otherwise, it automatically switches to read the word data, modifies the bits, and writes the word in the way. (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法 | [WriteAsync(String, Boolean)](dbc7b234-d200-d7e8-e308-fed865d7aa95.htm) | 向线圈中写入bool数组，返回是否写入成功，如果富文本地址不指定，默认使用的功能码是 0x0F  Write the bool array to the coil, and return whether the writing is successful. If the rich text address is not specified, the default function code is 0x0F. (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Byte)](104f3bf6-6822-f73c-3419-e8c99d5a3c15.htm) | 将数据写入到Modbus的寄存器上去，需要指定起始地址和数据内容，如果富文本地址不指定，默认使用的功能码是 0x10  To write data to Modbus registers, you need to specify the start address and data content. If the rich text address is not specified, the default function code is 0x10 (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Double)](910251ea-28ae-052b-b95b-c6caf8fe6483.htm) | 异步写入double数组，返回是否成功  Asynchronously write double array, return whether the write was successful (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法 | [WriteAsync(String, Int16)](17151d4f-bde5-cab5-e749-0992ad67de8b.htm) | 将数据写入到Modbus的单个寄存器上去，需要指定起始地址和数据值，如果富文本地址不指定，默认使用的功能码是 0x06  To write data to a single register of Modbus, you need to specify the start address and data value. If the rich text address is not specified, the default function code is 0x06. (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int32)](7bdf041e-af53-261c-ec07-0b1b85c1c735.htm) | 异步写入int[]数组，返回是否成功  Asynchronously write int array, return whether the write was successful (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int64)](d8069128-a5ea-a0c1-2221-3b151b626c53.htm) | 异步写入long数组，返回是否成功  Asynchronously write long array, return whether the write was successful (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Single)](2af0147b-de42-f977-b99f-e91f585e0c9a.htm) | 异步写入float数组，返回是否成功  Asynchronously write float array, return whether the write was successful (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法 | [WriteAsync(String, UInt16)](d08d6403-4da8-13fe-3762-86932dab62de.htm) | 将数据写入到Modbus的单个寄存器上去，需要指定起始地址和数据值，如果富文本地址不指定，默认使用的功能码是 0x06  To write data to a single register of Modbus, you need to specify the start address and data value. If the rich text address is not specified, the default function code is 0x06. (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt32)](eb068ca1-b343-5288-05a9-f8481d323bdf.htm) | 异步写入uint[]数组，返回是否成功  Asynchronously write uint array, return whether the write was successful (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt64)](e19ff622-f678-4c52-1662-2af481317995.htm) | 异步写入ulong数组，返回是否成功  Asynchronously write ulong array, return whether the write was successful (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法 | [WriteAsync(String, String, Int32)](b739b983-b03c-8ecf-2713-75f19233ae2f.htm) | 异步写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为ASCII  Asynchronously write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, it is filled with 0 and the encoding is ASCII. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, String, Encoding)](876e9800-0406-6407-bdf8-8ab1cf28df38.htm) | 异步写入字符串信息，需要指定的编码信息  Asynchronously write string information, need to specify the encoding information (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WriteAsync(String, String, Int32, Encoding)](ae9fb2a3-fddb-48eb-e9ad-6cc0ebbac1f6.htm) | 异步写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为指定的编码信息  Asynchronously write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, then the operation is complemented with 0 , you should specified the encoding information (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsyncT(T)](25d96e77-d1b2-918c-33f9-e5e4cc57c9cd.htm) | 异步写入支持Hsl特性的数据，返回是否写入成功，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Asynchronously write data that supports the Hsl attribute, and return whether the write was successful. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteCustomerT](84a37f82-3f67-a0f8-fd3f-0630772565e4.htm) | 写入自定义类型的数据，该类型必须继承自IDataTransfer接口  Write data of a custom type, which must inherit from the IDataTransfer interface (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteCustomerAsyncT](04628397-cfa1-d785-7790-4c02cb377ab6.htm) | 写入自定义类型的数据，该类型必须继承自IDataTransfer接口  Write data of a custom type, which must inherit from the IDataTransfer interface (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WriteFile](c888070e-6866-6228-2d68-533c2b238fd8.htm) | 使用0x15功能码写入文件的数据内容，需要指定文件编号，地址，数据，返回写入的结果数据。写入的数据字节长度不能大于248个字节，如果想写入更多的字节，就需要切片写入操作 (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法 | [WriteMask](5cce8a41-41b1-90d8-b44e-07c85b38041a.htm) | 向设备写入掩码数据，使用0x16功能码，需要确认对方是否支持相关的操作，掩码数据的操作主要针对寄存器。  To write mask data to the server, using the 0x16 function code, you need to confirm whether the other party supports related operations. The operation of mask data is mainly directed to the register. (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法 | [WriteMaskAsync](24048bc4-090d-04c2-b762-85ac5c836c8d.htm) | 向设备写入掩码数据，使用0x16功能码，需要确认对方是否支持相关的操作，掩码数据的操作主要针对寄存器。  To write mask data to the server, using the 0x16 function code, you need to confirm whether the other party supports related operations. The operation of mask data is mainly directed to the register. (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法 | [WriteOneRegister(String, Int16)](685792aa-f82c-7400-8319-d62d56af0aba.htm) | 将数据写入到Modbus的单个寄存器上去，需要指定起始地址和数据值，如果富文本地址不指定，默认使用的功能码是 0x06  To write data to a single register of Modbus, you need to specify the start address and data value. If the rich text address is not specified, the default function code is 0x06. (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法 | [WriteOneRegister(String, UInt16)](38dca771-988c-7bbb-1963-8c354830cc3f.htm) | 将数据写入到Modbus的单个寄存器上去，需要指定起始地址和数据值，如果富文本地址不指定，默认使用的功能码是 0x06  To write data to a single register of Modbus, you need to specify the start address and data value. If the rich text address is not specified, the default function code is 0x06. (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法 | [WriteOneRegisterAsync(String, Int16)](7c9fb70d-655c-0544-fa34-086ffe8ad357.htm) | 将数据写入到Modbus的单个寄存器上去，需要指定起始地址和数据值，如果富文本地址不指定，默认使用的功能码是 0x06  To write data to a single register of Modbus, you need to specify the start address and data value. If the rich text address is not specified, the default function code is 0x06. (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法 | [WriteOneRegisterAsync(String, UInt16)](5d3d62d9-e9fc-a807-36d6-5771b9bac661.htm) | 将数据写入到Modbus的单个寄存器上去，需要指定起始地址和数据值，如果富文本地址不指定，默认使用的功能码是 0x06  To write data to a single register of Modbus, you need to specify the start address and data value. If the rich text address is not specified, the default function code is 0x06. (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |

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

本客户端支持的标准的modbus-ascii协议，地址支持富文本格式，具体参考示例代码。  
读取线圈，输入线圈，寄存器，输入寄存器的方法中的读取长度对商业授权用户不限制，内部自动切割读取，结果合并。

![](../icons/SectionExpanded.png)示例

基本的用法请参照下面的代码示例，初始化部分的代码省略

Modbus示例

[复制](# "复制")

```
// 本类支持的读写操作提供了非常多的重载方法，总有你想要的方法
private ModbusAscii modbus = new ModbusAscii( );   // 实例化



private void CoilExample( )
{
    // 读取线圈示例
    bool coil100 = modbus.ReadCoil( "100" ).Content;

    // 判断是否读取成功
    OperateResult<bool> result_coil100 = modbus.ReadCoil( "100" );
    if (result_coil100.IsSuccess)
    {
        // success
        bool value = result_coil100.Content;
    }
    else
    {
        // failed
    }


    // 假设读取站号10的线圈100的值
    bool coil_station_ten_100 = modbus.ReadCoil( "s=10;100" ).Content;



    // =============================================================================================
    // 写入也是同理，线圈100写通
    modbus.Write( "100", true );

    // 站号10的线圈写通
    modbus.Write( "s=10;100", true );

    // 想要判断是否写入成功
    if (modbus.Write( "s=10;100", true ).IsSuccess)
    {
        // success
    }
    else
    {
        // failed
    }



    // ===========================================================================================
    // 批量读写也是类似，批量的读取
    bool[] coil10_19 = modbus.ReadCoil( "100", 10 ).Content;

    // 写入也是同理
    modbus.Write( "100", new bool[] { true, false, true, false, false, false, true, false, false, false } );


    // 离散输入的数据读取同理
}


private void RegisterExample( )
{
    // 读取寄存器100的值
    short register100 = modbus.ReadInt16( "100" ).Content;

    // 批量读取寄存器100-109的值
    short[] register100_109 = modbus.ReadInt16( "100", 10 ).Content;

    // 写入寄存器100的值，注意，一定要强制转换short类型
    modbus.Write( "100", (short)123 );

    // 批量写
    modbus.Write( "100", new short[] { 123, -123, 4244 } );


    // ==============================================================================================
    // 以下是一些常规的操作，不再对是否成功的结果进行判断
    // 读取操作

    bool coil100 = modbus.ReadCoil( "100" ).Content;   // 读取线圈100的通断
    short short100 = modbus.ReadInt16( "100" ).Content; // 读取寄存器100的short值
    ushort ushort100 = modbus.ReadUInt16( "100" ).Content; // 读取寄存器100的ushort值
    int int100 = modbus.ReadInt32( "100" ).Content;      // 读取寄存器100-101的int值
    uint uint100 = modbus.ReadUInt32( "100" ).Content;   // 读取寄存器100-101的uint值
    float float100 = modbus.ReadFloat( "100" ).Content; // 读取寄存器100-101的float值
    long long100 = modbus.ReadInt64( "100" ).Content;    // 读取寄存器100-103的long值
    ulong ulong100 = modbus.ReadUInt64( "100" ).Content; // 读取寄存器100-103的ulong值
    double double100 = modbus.ReadDouble( "100" ).Content; // 读取寄存器100-103的double值
    string str100 = modbus.ReadString( "100", 5 ).Content;// 读取100到104共10个字符的字符串

    // 写入操作
    modbus.Write( "100", true );// 写入线圈100为通
    modbus.Write( "100", (short)12345 );// 写入寄存器100为12345
    modbus.Write( "100", (ushort)45678 );// 写入寄存器100为45678
    modbus.Write( "100", 123456789 );// 写入寄存器100-101为123456789
    modbus.Write( "100", (uint)123456778 );// 写入寄存器100-101为123456778
    modbus.Write( "100", 123.456 );// 写入寄存器100-101为123.456
    modbus.Write( "100", 12312312312414L );//写入寄存器100-103为一个大数据
    modbus.Write( "100", 12634534534543656UL );// 写入寄存器100-103为一个大数据
    modbus.Write( "100", 123.456d );// 写入寄存器100-103为一个双精度的数据
    modbus.Write( "100", "K123456789" );

    // ===============================================================================================
    // 读取输入寄存器
    short input_short100 = modbus.ReadInt16( "x=4;100" ).Content; // 读取寄存器100的short值
    ushort input_ushort100 = modbus.ReadUInt16( "x=4;100" ).Content; // 读取寄存器100的ushort值
    int input_int100 = modbus.ReadInt32( "x=4;100" ).Content;      // 读取寄存器100-101的int值
    uint input_uint100 = modbus.ReadUInt32( "x=4;100" ).Content;   // 读取寄存器100-101的uint值
    float input_float100 = modbus.ReadFloat( "x=4;100" ).Content; // 读取寄存器100-101的float值
    long input_long100 = modbus.ReadInt64( "x=4;100" ).Content;    // 读取寄存器100-103的long值
    ulong input_ulong100 = modbus.ReadUInt64( "x=4;100" ).Content; // 读取寄存器100-103的ulong值
    double input_double100 = modbus.ReadDouble( "x=4;100" ).Content; // 读取寄存器100-103的double值
    string input_str100 = modbus.ReadString( "x=4;100", 5 ).Content;// 读取100到104共10个字符的字符串
}
```

复杂的读取数据的代码示例如下：

read示例

[复制](# "复制")

```
ModbusRtu modbus = new ModbusRtu( );   // 实例化
// 此处忽略初始化
// modbus.SerialPortInni( "COM3" );

// 假设100存储了short的报警，101,102存储了float的温度，103，104存储了int的产量
OperateResult<byte[]> read = modbus.Read( "100", 5 );
if (read.IsSuccess)
{
    // 共计10个字节的结果内容
    short alarm = modbus.ByteTransform.TransInt16( read.Content, 0 );
    float temp = modbus.ByteTransform.TransSingle( read.Content, 2 );
    int product = modbus.ByteTransform.TransInt32( read.Content, 6 );
}
else
{
    // failed
}
```

写入数据的代码如下：

write示例

[复制](# "复制")

```
ModbusRtu modbus = new ModbusRtu( );   // 实例化
// 此处忽略初始化
// modbus.SerialPortInni( "COM3" );

// 假设100存储了short的报警，101,102存储了float的温度，103，104存储了int的产量
byte[] buffer = new byte[10];
modbus.ByteTransform.TransByte( (short)1 ).CopyTo( buffer, 0 );
modbus.ByteTransform.TransByte( 123.456f ).CopyTo( buffer, 2 );
modbus.ByteTransform.TransByte( 45678922 ).CopyTo( buffer, 6 );

OperateResult write = modbus.Write( "100", buffer );
if (write.IsSuccess)
{
    // success
}
else
{
    // failed
}
```

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.ModBus 命名空间](da9ac567-746e-47e5-70e1-b44577d09ce4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)