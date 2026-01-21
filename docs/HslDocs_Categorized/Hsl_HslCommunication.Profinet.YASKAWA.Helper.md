# HslCommunication - HslCommunication.Profinet.YASKAWA.Helper

> 分類頁數: 30



---
## HslCommunication.Profinet.YASKAWA.Helper

[原文連結](http://api.hslcommunication.cn/html/310f2e91-5abf-4755-25f7-9225e316f6fe.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.YASKAWA.Helper](../html/310f2e91-5abf-4755-25f7-9225e316f6fe.htm "HslCommunication.Profinet.YASKAWA.Helper")

[IMemobus 接口](../html/0f5ca9bf-1080-d913-ed93-7c9a46291f73.htm "IMemobus 接口")

[MemobusHelper 类](../html/254c50cc-c87c-57ad-eb8f-bda323fe9527.htm "MemobusHelper 类")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslCommunication.Profinet.YASKAWA.Helper 命名空间 |

[缺少 "N:HslCommunication.Profinet.YASKAWA.Helper" 的 <summary> 文档]

![](../icons/SectionExpanded.png)类

|  | 类 | 说明 |
| --- | --- | --- |
| 公共类 | [MemobusHelper](254c50cc-c87c-57ad-eb8f-bda323fe9527.htm) | Memobus的辅助类对象 |

![](../icons/SectionExpanded.png)接口

|  | 接口 | 说明 |
| --- | --- | --- |
| 公共接口 | [IMemobus](0f5ca9bf-1080-d913-ed93-7c9a46291f73.htm) | Memobus协议的接口信息 |

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## IMemobus 接口

[原文連結](http://api.hslcommunication.cn/html/0f5ca9bf-1080-d913-ed93-7c9a46291f73.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.YASKAWA.Helper](../html/310f2e91-5abf-4755-25f7-9225e316f6fe.htm "HslCommunication.Profinet.YASKAWA.Helper")

[IMemobus 接口](../html/0f5ca9bf-1080-d913-ed93-7c9a46291f73.htm "IMemobus 接口")

[IMemobus 属性](../html/df4ac10e-fc77-944a-49a5-6ab842b0e1fb.htm "IMemobus 属性")

[IMemobus 方法](../html/59bfc9d8-9df6-cffb-4a6a-9125b9d3be0c.htm "IMemobus 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IMemobus 接口 |

Memobus协议的接口信息

**命名空间：**
 [HslCommunication.Profinet.YASKAWA.Helper](310f2e91-5abf-4755-25f7-9225e316f6fe.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public interface IMemobus : IReadWriteDevice, 
	IReadWriteNet
```

```
Public Interface IMemobus
	Inherits IReadWriteDevice, IReadWriteNet
```

```
public interface class IMemobus : IReadWriteDevice, 
	IReadWriteNet
```

```
type IMemobus =  
    interface
        interface IReadWriteDevice
        interface IReadWriteNet
    end
```

IMemobus 类型公开以下成员。

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性代码示例 | [ByteTransform](0836ec33-d464-ed39-37ff-b2ada0d286fc.htm) | 当前的数据变换机制，当你需要从字节数据转换类型数据的时候需要。  The current data transformation mechanism is required when you need to convert type data from byte data. (继承自 [IReadWriteDevice](d987afcc-f99c-6aed-3847-5b1890504635.htm)。) |
| 公共属性 | [ConnectionId](20787592-e3ec-e0d0-27c2-abfa24654bd2.htm) | 当前连接的唯一ID号，默认为长度20的guid码加随机数组成，方便列表管理，也可以自己指定  The unique ID number of the current connection. The default is a 20-digit guid code plus a random number. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共属性 | [CpuFrom](c901404e-6f8c-fe52-3650-888c03c9c414.htm) | 获取或设置发送源的CPU的编号信息，默认为 1  Get or set the number information of the sending source CPU, the default is 1 |
| 公共属性 | [CpuTo](ee8252b3-b1e4-802c-52f2-e266bfb6faed.htm) | 获取或设置发送目标的CPU的编号信息，默认为 2  Get or set the CPU number information of the sending destination, the default is 2 |
| 公共属性代码示例 | [LogNet](38760c1f-22e3-f8fe-7512-15849a4b6eb0.htm) | 组件的日志工具，支持日志记录，只要实例化后，当前网络的基本信息，就以[DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm)等级进行输出  The component's logging tool supports logging. As long as the instantiation of the basic network information, the output will be output at [DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm) (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |

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

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.Profinet.YASKAWA.Helper 命名空间](310f2e91-5abf-4755-25f7-9225e316f6fe.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## IMemobus 属性

[原文連結](http://api.hslcommunication.cn/html/df4ac10e-fc77-944a-49a5-6ab842b0e1fb.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.YASKAWA.Helper](../html/310f2e91-5abf-4755-25f7-9225e316f6fe.htm "HslCommunication.Profinet.YASKAWA.Helper")

[IMemobus 接口](../html/0f5ca9bf-1080-d913-ed93-7c9a46291f73.htm "IMemobus 接口")

[IMemobus 属性](../html/df4ac10e-fc77-944a-49a5-6ab842b0e1fb.htm "IMemobus 属性")

[CpuFrom 属性](../html/c901404e-6f8c-fe52-3650-888c03c9c414.htm "CpuFrom 属性 ")

[CpuTo 属性](../html/ee8252b3-b1e4-802c-52f2-e266bfb6faed.htm "CpuTo 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IMemobus 属性 |

[IMemobus](0f5ca9bf-1080-d913-ed93-7c9a46291f73.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性代码示例 | [ByteTransform](0836ec33-d464-ed39-37ff-b2ada0d286fc.htm) | 当前的数据变换机制，当你需要从字节数据转换类型数据的时候需要。  The current data transformation mechanism is required when you need to convert type data from byte data. (继承自 [IReadWriteDevice](d987afcc-f99c-6aed-3847-5b1890504635.htm)。) |
| 公共属性 | [ConnectionId](20787592-e3ec-e0d0-27c2-abfa24654bd2.htm) | 当前连接的唯一ID号，默认为长度20的guid码加随机数组成，方便列表管理，也可以自己指定  The unique ID number of the current connection. The default is a 20-digit guid code plus a random number. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共属性 | [CpuFrom](c901404e-6f8c-fe52-3650-888c03c9c414.htm) | 获取或设置发送源的CPU的编号信息，默认为 1  Get or set the number information of the sending source CPU, the default is 1 |
| 公共属性 | [CpuTo](ee8252b3-b1e4-802c-52f2-e266bfb6faed.htm) | 获取或设置发送目标的CPU的编号信息，默认为 2  Get or set the CPU number information of the sending destination, the default is 2 |
| 公共属性代码示例 | [LogNet](38760c1f-22e3-f8fe-7512-15849a4b6eb0.htm) | 组件的日志工具，支持日志记录，只要实例化后，当前网络的基本信息，就以[DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm)等级进行输出  The component's logging tool supports logging. As long as the instantiation of the basic network information, the output will be output at [DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm) (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[IMemobus 接口](0f5ca9bf-1080-d913-ed93-7c9a46291f73.htm)

[HslCommunication.Profinet.YASKAWA.Helper 命名空间](310f2e91-5abf-4755-25f7-9225e316f6fe.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## CpuFrom 属性 

[原文連結](http://api.hslcommunication.cn/html/c901404e-6f8c-fe52-3650-888c03c9c414.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.YASKAWA.Helper](../html/310f2e91-5abf-4755-25f7-9225e316f6fe.htm "HslCommunication.Profinet.YASKAWA.Helper")

[IMemobus 接口](../html/0f5ca9bf-1080-d913-ed93-7c9a46291f73.htm "IMemobus 接口")

[IMemobus 属性](../html/df4ac10e-fc77-944a-49a5-6ab842b0e1fb.htm "IMemobus 属性")

[CpuFrom 属性](../html/c901404e-6f8c-fe52-3650-888c03c9c414.htm "CpuFrom 属性 ")

[CpuTo 属性](../html/ee8252b3-b1e4-802c-52f2-e266bfb6faed.htm "CpuTo 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IMemobusCpuFrom 属性 |

获取或设置发送源的CPU的编号信息，默认为 1  
Get or set the number information of the sending source CPU, the default is 1

**命名空间：**
 [HslCommunication.Profinet.YASKAWA.Helper](310f2e91-5abf-4755-25f7-9225e316f6fe.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
byte CpuFrom { get; set; }
```

```
Property CpuFrom As Byte
	Get
	Set
```

```
property unsigned char CpuFrom {
	unsigned char get ();
	void set (unsigned char value);
}
```

```
abstract CpuFrom : byte with get, set
```

#### 属性值

类型：Byte

![](../icons/SectionExpanded.png)参见

#### 引用

[IMemobus 接口](0f5ca9bf-1080-d913-ed93-7c9a46291f73.htm)

[HslCommunication.Profinet.YASKAWA.Helper 命名空间](310f2e91-5abf-4755-25f7-9225e316f6fe.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## CpuTo 属性 

[原文連結](http://api.hslcommunication.cn/html/ee8252b3-b1e4-802c-52f2-e266bfb6faed.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.YASKAWA.Helper](../html/310f2e91-5abf-4755-25f7-9225e316f6fe.htm "HslCommunication.Profinet.YASKAWA.Helper")

[IMemobus 接口](../html/0f5ca9bf-1080-d913-ed93-7c9a46291f73.htm "IMemobus 接口")

[IMemobus 属性](../html/df4ac10e-fc77-944a-49a5-6ab842b0e1fb.htm "IMemobus 属性")

[CpuFrom 属性](../html/c901404e-6f8c-fe52-3650-888c03c9c414.htm "CpuFrom 属性 ")

[CpuTo 属性](../html/ee8252b3-b1e4-802c-52f2-e266bfb6faed.htm "CpuTo 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IMemobusCpuTo 属性 |

获取或设置发送目标的CPU的编号信息，默认为 2  
Get or set the CPU number information of the sending destination, the default is 2

**命名空间：**
 [HslCommunication.Profinet.YASKAWA.Helper](310f2e91-5abf-4755-25f7-9225e316f6fe.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
byte CpuTo { get; set; }
```

```
Property CpuTo As Byte
	Get
	Set
```

```
property unsigned char CpuTo {
	unsigned char get ();
	void set (unsigned char value);
}
```

```
abstract CpuTo : byte with get, set
```

#### 属性值

类型：Byte

![](../icons/SectionExpanded.png)参见

#### 引用

[IMemobus 接口](0f5ca9bf-1080-d913-ed93-7c9a46291f73.htm)

[HslCommunication.Profinet.YASKAWA.Helper 命名空间](310f2e91-5abf-4755-25f7-9225e316f6fe.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## IMemobus 方法

[原文連結](http://api.hslcommunication.cn/html/59bfc9d8-9df6-cffb-4a6a-9125b9d3be0c.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.YASKAWA.Helper](../html/310f2e91-5abf-4755-25f7-9225e316f6fe.htm "HslCommunication.Profinet.YASKAWA.Helper")

[IMemobus 接口](../html/0f5ca9bf-1080-d913-ed93-7c9a46291f73.htm "IMemobus 接口")

[IMemobus 属性](../html/df4ac10e-fc77-944a-49a5-6ab842b0e1fb.htm "IMemobus 属性")

[IMemobus 方法](../html/59bfc9d8-9df6-cffb-4a6a-9125b9d3be0c.htm "IMemobus 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IMemobus 方法 |

[IMemobus](0f5ca9bf-1080-d913-ed93-7c9a46291f73.htm) 类型公开以下成员。

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

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[IMemobus 接口](0f5ca9bf-1080-d913-ed93-7c9a46291f73.htm)

[HslCommunication.Profinet.YASKAWA.Helper 命名空间](310f2e91-5abf-4755-25f7-9225e316f6fe.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## MemobusHelper 类

[原文連結](http://api.hslcommunication.cn/html/254c50cc-c87c-57ad-eb8f-bda323fe9527.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.YASKAWA.Helper](../html/310f2e91-5abf-4755-25f7-9225e316f6fe.htm "HslCommunication.Profinet.YASKAWA.Helper")

[MemobusHelper 类](../html/254c50cc-c87c-57ad-eb8f-bda323fe9527.htm "MemobusHelper 类")

[MemobusHelper 构造函数](../html/ff33d620-e945-02fa-13b4-2b1f4b2f652d.htm "MemobusHelper 构造函数 ")

[MemobusHelper 方法](../html/e43bdd34-e6b5-c55a-deed-0d20ed628891.htm "MemobusHelper 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MemobusHelper 类 |

Memobus的辅助类对象

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  HslCommunication.Profinet.YASKAWA.HelperMemobusHelper

**命名空间：**
 [HslCommunication.Profinet.YASKAWA.Helper](310f2e91-5abf-4755-25f7-9225e316f6fe.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class MemobusHelper
```

```
Public Class MemobusHelper
```

```
public ref class MemobusHelper
```

```
type MemobusHelper =  class end
```

MemobusHelper 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [MemobusHelper](ff33d620-e945-02fa-13b4-2b1f4b2f652d.htm) | 初始化 MemobusHelper 类的一个新实例 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法静态成员 | [Read](edf5ccc5-454d-c522-2d86-234f53054551.htm) | 批量读取字节数组信息，需要指定地址和长度，返回原始的字节数组  Batch read byte array information, need to specify the address and length, return the original byte array |
| 公共方法静态成员 | [ReadAsync](41d84bd7-1ae9-01bb-4813-5c50b2dfd419.htm) | 批量读取字节数组信息，需要指定地址和长度，返回原始的字节数组  Batch read byte array information, need to specify the address and length, return the original byte array |
| 公共方法静态成员 | [ReadBool](7f5d3261-4c7a-ef65-1226-f31cc4b3e4da.htm) | 批量读取Boolean数组信息，需要指定地址和长度，返回Boolean 数组  Batch read Boolean array information, need to specify the address and length, return Boolean array |
| 公共方法静态成员 | [ReadBoolAsync](fbe714bf-2e73-6ffc-6aee-c45e3f583d43.htm) | 批量读取Boolean数组信息，需要指定地址和长度，返回Boolean 数组  Batch read Boolean array information, need to specify the address and length, return Boolean array |
| 公共方法静态成员 | [ReadRandom(IMemobus, String)](bdf0757b-18b2-0b9e-366a-a012a5c88955.htm) | 随机读取扩展的保持寄存器的内容，也即读取不连续地址的字数据，可以指定多个地址，然后一次性读取所有的数据，然后解析出实际的数据  Randomly read the contents of the extended hold register, that is, read word data of discontinuous addresses, you can specify multiple addresses, then read all the data at once, and then parse out the actual data |
| 公共方法静态成员 | [ReadRandom(IMemobus, UInt16)](26f4b67e-78ef-6494-615a-9db55bfb15ae.htm) | 随机读取扩展的保持寄存器的内容，也即读取不连续地址的字数据，可以指定多个地址，然后一次性读取所有的数据，然后解析出实际的数据  Randomly read the contents of the extended hold register, that is, read word data of discontinuous addresses, you can specify multiple addresses, then read all the data at once, and then parse out the actual data |
| 公共方法静态成员 | [ReadRandomAsync(IMemobus, String)](0edac1c0-22c9-5164-527f-d4533556bf61.htm) | 随机读取扩展的保持寄存器的内容，也即读取不连续地址的字数据，可以指定多个地址，然后一次性读取所有的数据，然后解析出实际的数据  Randomly read the contents of the extended hold register, that is, read word data of discontinuous addresses, you can specify multiple addresses, then read all the data at once, and then parse out the actual data |
| 公共方法静态成员 | [ReadRandomAsync(IMemobus, UInt16)](8aa741a5-4c85-127e-d8be-6677ba4be1cf.htm) | 随机读取扩展的保持寄存器的内容，也即读取不连续地址的字数据，可以指定多个地址，然后一次性读取所有的数据，然后解析出实际的数据  Randomly read the contents of the extended hold register, that is, read word data of discontinuous addresses, you can specify multiple addresses, then read all the data at once, and then parse out the actual data |
| 公共方法 | ToString | (继承自 Object。) |
| 公共方法静态成员 | [Write(IMemobus, String, Boolean)](1f35cb0d-f4ae-b48b-6830-6f58f66e94d5.htm) | 写入单个的Boolean数据，返回是否成功  Write a single Boolean data, and return whether the write was successful |
| 公共方法静态成员 | [Write(IMemobus, String, Boolean)](7cf177a0-b8e8-78da-fb1f-dcd8b1bec795.htm) | 批量写入Boolean数组数据，返回是否成功  Batch write Boolean array data, return whether the write was successful |
| 公共方法静态成员 | [Write(IMemobus, String, Byte)](a304cab6-f89e-e9e8-ffb6-8906b7a15382.htm) | 写入原始的byte数组数据到指定的地址，返回是否写入成功  Write the original byte array data to the specified address, and return whether the write was successful |
| 公共方法静态成员代码示例 | [Write(IMemobus, String, Int16, FuncString, Int16, OperateResult)](ac26cdc6-3666-0f1e-0d3c-245e81c57aff.htm) | 写入short数据，返回是否成功  Write short data, returns whether success |
| 公共方法静态成员代码示例 | [Write(IMemobus, String, UInt16, FuncString, UInt16, OperateResult)](b37601c7-5f78-2779-1479-d6e47642f8d1.htm) | 写入ushort数据，返回是否成功  Write ushort data, return whether the write was successful |
| 公共方法静态成员 | [WriteAsync(IMemobus, String, Boolean)](1428b3bb-6eb9-ef9d-3af5-d0965535426f.htm) | 写入单个的Boolean数据，返回是否成功  Write a single Boolean data, and return whether the write was successful |
| 公共方法静态成员 | [WriteAsync(IMemobus, String, Boolean)](1724f1cb-ddbe-e3c2-c1f0-acbd5b105ddd.htm) | 批量写入Boolean数组数据，返回是否成功  Batch write Boolean array data, return whether the write was successful |
| 公共方法静态成员 | [WriteAsync(IMemobus, String, Byte)](0391768c-0e3f-e53a-d516-566224a49be8.htm) | 写入原始的byte数组数据到指定的地址，返回是否写入成功  Write the original byte array data to the specified address, and return whether the write was successful |
| 公共方法静态成员代码示例 | [WriteAsync(IMemobus, String, Int16, FuncString, Int16, TaskOperateResult)](fcb2e4f3-8e61-75b5-c416-2be479b4e69d.htm) | 写入short数据，返回是否成功  Write short data, returns whether success |
| 公共方法静态成员代码示例 | [WriteAsync(IMemobus, String, UInt16, FuncString, UInt16, TaskOperateResult)](df1c6846-765e-58ee-cf5e-cc360a65f483.htm) | 写入ushort数据，返回是否成功  Write ushort data, return whether the write was successful |
| 公共方法静态成员 | [WriteRandom](53f86b45-0abc-47c2-a1c0-58247705efa1.htm) | 随机写入扩展的保持寄存器的内容，也即写入不连续的地址的字数据，字节数组的长度必须为地址数组长度的两倍，才能正确写入。  Write the contents of the extended hold registers randomly, that is, write word data for discontinuous addresses, and the byte array must be twice the length of the address array to be written correctly. |
| 公共方法静态成员 | [WriteRandomAsync](768d6601-f297-57f7-c3a9-2deed351fcbf.htm) | 随机写入扩展的保持寄存器的内容，也即写入不连续的地址的字数据，字节数组的长度必须为地址数组长度的两倍，才能正确写入。  Write the contents of the extended hold registers randomly, that is, write word data for discontinuous addresses, and the byte array must be twice the length of the address array to be written correctly. |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.Profinet.YASKAWA.Helper 命名空间](310f2e91-5abf-4755-25f7-9225e316f6fe.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## MemobusHelper 构造函数 

[原文連結](http://api.hslcommunication.cn/html/ff33d620-e945-02fa-13b4-2b1f4b2f652d.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.YASKAWA.Helper](../html/310f2e91-5abf-4755-25f7-9225e316f6fe.htm "HslCommunication.Profinet.YASKAWA.Helper")

[MemobusHelper 类](../html/254c50cc-c87c-57ad-eb8f-bda323fe9527.htm "MemobusHelper 类")

[MemobusHelper 构造函数](../html/ff33d620-e945-02fa-13b4-2b1f4b2f652d.htm "MemobusHelper 构造函数 ")

[MemobusHelper 方法](../html/e43bdd34-e6b5-c55a-deed-0d20ed628891.htm "MemobusHelper 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MemobusHelper 构造函数 |

初始化 [MemobusHelper](254c50cc-c87c-57ad-eb8f-bda323fe9527.htm) 类的一个新实例

**命名空间：**
 [HslCommunication.Profinet.YASKAWA.Helper](310f2e91-5abf-4755-25f7-9225e316f6fe.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public MemobusHelper()
```

```
Public Sub New
```

```
public:
MemobusHelper()
```

```
new : unit -> MemobusHelper
```

![](../icons/SectionExpanded.png)参见

#### 引用

[MemobusHelper 类](254c50cc-c87c-57ad-eb8f-bda323fe9527.htm)

[HslCommunication.Profinet.YASKAWA.Helper 命名空间](310f2e91-5abf-4755-25f7-9225e316f6fe.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## MemobusHelper 方法

[原文連結](http://api.hslcommunication.cn/html/e43bdd34-e6b5-c55a-deed-0d20ed628891.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.YASKAWA.Helper](../html/310f2e91-5abf-4755-25f7-9225e316f6fe.htm "HslCommunication.Profinet.YASKAWA.Helper")

[MemobusHelper 类](../html/254c50cc-c87c-57ad-eb8f-bda323fe9527.htm "MemobusHelper 类")

[MemobusHelper 方法](../html/e43bdd34-e6b5-c55a-deed-0d20ed628891.htm "MemobusHelper 方法")

[Read 方法](../html/edf5ccc5-454d-c522-2d86-234f53054551.htm "Read 方法 ")

[ReadAsync 方法](../html/41d84bd7-1ae9-01bb-4813-5c50b2dfd419.htm "ReadAsync 方法 ")

[ReadBool 方法](../html/7f5d3261-4c7a-ef65-1226-f31cc4b3e4da.htm "ReadBool 方法 ")

[ReadBoolAsync 方法](../html/fbe714bf-2e73-6ffc-6aee-c45e3f583d43.htm "ReadBoolAsync 方法 ")

[ReadRandom 方法](../html/d19df820-e22b-5da4-280e-13c38e12cd37.htm "ReadRandom 方法 ")

[ReadRandomAsync 方法](../html/b237b1ea-3ff4-b07a-c31f-8e2686ff566b.htm "ReadRandomAsync 方法 ")

[Write 方法](../html/fb613f76-96a5-5566-eb37-f400dbd0a93c.htm "Write 方法 ")

[WriteAsync 方法](../html/bbc31e49-e3cd-881d-75b1-42ec61c72b34.htm "WriteAsync 方法 ")

[WriteRandom 方法](../html/53f86b45-0abc-47c2-a1c0-58247705efa1.htm "WriteRandom 方法 ")

[WriteRandomAsync 方法](../html/768d6601-f297-57f7-c3a9-2deed351fcbf.htm "WriteRandomAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MemobusHelper 方法 |

[MemobusHelper](254c50cc-c87c-57ad-eb8f-bda323fe9527.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法静态成员 | [Read](edf5ccc5-454d-c522-2d86-234f53054551.htm) | 批量读取字节数组信息，需要指定地址和长度，返回原始的字节数组  Batch read byte array information, need to specify the address and length, return the original byte array |
| 公共方法静态成员 | [ReadAsync](41d84bd7-1ae9-01bb-4813-5c50b2dfd419.htm) | 批量读取字节数组信息，需要指定地址和长度，返回原始的字节数组  Batch read byte array information, need to specify the address and length, return the original byte array |
| 公共方法静态成员 | [ReadBool](7f5d3261-4c7a-ef65-1226-f31cc4b3e4da.htm) | 批量读取Boolean数组信息，需要指定地址和长度，返回Boolean 数组  Batch read Boolean array information, need to specify the address and length, return Boolean array |
| 公共方法静态成员 | [ReadBoolAsync](fbe714bf-2e73-6ffc-6aee-c45e3f583d43.htm) | 批量读取Boolean数组信息，需要指定地址和长度，返回Boolean 数组  Batch read Boolean array information, need to specify the address and length, return Boolean array |
| 公共方法静态成员 | [ReadRandom(IMemobus, String)](bdf0757b-18b2-0b9e-366a-a012a5c88955.htm) | 随机读取扩展的保持寄存器的内容，也即读取不连续地址的字数据，可以指定多个地址，然后一次性读取所有的数据，然后解析出实际的数据  Randomly read the contents of the extended hold register, that is, read word data of discontinuous addresses, you can specify multiple addresses, then read all the data at once, and then parse out the actual data |
| 公共方法静态成员 | [ReadRandom(IMemobus, UInt16)](26f4b67e-78ef-6494-615a-9db55bfb15ae.htm) | 随机读取扩展的保持寄存器的内容，也即读取不连续地址的字数据，可以指定多个地址，然后一次性读取所有的数据，然后解析出实际的数据  Randomly read the contents of the extended hold register, that is, read word data of discontinuous addresses, you can specify multiple addresses, then read all the data at once, and then parse out the actual data |
| 公共方法静态成员 | [ReadRandomAsync(IMemobus, String)](0edac1c0-22c9-5164-527f-d4533556bf61.htm) | 随机读取扩展的保持寄存器的内容，也即读取不连续地址的字数据，可以指定多个地址，然后一次性读取所有的数据，然后解析出实际的数据  Randomly read the contents of the extended hold register, that is, read word data of discontinuous addresses, you can specify multiple addresses, then read all the data at once, and then parse out the actual data |
| 公共方法静态成员 | [ReadRandomAsync(IMemobus, UInt16)](8aa741a5-4c85-127e-d8be-6677ba4be1cf.htm) | 随机读取扩展的保持寄存器的内容，也即读取不连续地址的字数据，可以指定多个地址，然后一次性读取所有的数据，然后解析出实际的数据  Randomly read the contents of the extended hold register, that is, read word data of discontinuous addresses, you can specify multiple addresses, then read all the data at once, and then parse out the actual data |
| 公共方法 | ToString | (继承自 Object。) |
| 公共方法静态成员 | [Write(IMemobus, String, Boolean)](1f35cb0d-f4ae-b48b-6830-6f58f66e94d5.htm) | 写入单个的Boolean数据，返回是否成功  Write a single Boolean data, and return whether the write was successful |
| 公共方法静态成员 | [Write(IMemobus, String, Boolean)](7cf177a0-b8e8-78da-fb1f-dcd8b1bec795.htm) | 批量写入Boolean数组数据，返回是否成功  Batch write Boolean array data, return whether the write was successful |
| 公共方法静态成员 | [Write(IMemobus, String, Byte)](a304cab6-f89e-e9e8-ffb6-8906b7a15382.htm) | 写入原始的byte数组数据到指定的地址，返回是否写入成功  Write the original byte array data to the specified address, and return whether the write was successful |
| 公共方法静态成员代码示例 | [Write(IMemobus, String, Int16, FuncString, Int16, OperateResult)](ac26cdc6-3666-0f1e-0d3c-245e81c57aff.htm) | 写入short数据，返回是否成功  Write short data, returns whether success |
| 公共方法静态成员代码示例 | [Write(IMemobus, String, UInt16, FuncString, UInt16, OperateResult)](b37601c7-5f78-2779-1479-d6e47642f8d1.htm) | 写入ushort数据，返回是否成功  Write ushort data, return whether the write was successful |
| 公共方法静态成员 | [WriteAsync(IMemobus, String, Boolean)](1428b3bb-6eb9-ef9d-3af5-d0965535426f.htm) | 写入单个的Boolean数据，返回是否成功  Write a single Boolean data, and return whether the write was successful |
| 公共方法静态成员 | [WriteAsync(IMemobus, String, Boolean)](1724f1cb-ddbe-e3c2-c1f0-acbd5b105ddd.htm) | 批量写入Boolean数组数据，返回是否成功  Batch write Boolean array data, return whether the write was successful |
| 公共方法静态成员 | [WriteAsync(IMemobus, String, Byte)](0391768c-0e3f-e53a-d516-566224a49be8.htm) | 写入原始的byte数组数据到指定的地址，返回是否写入成功  Write the original byte array data to the specified address, and return whether the write was successful |
| 公共方法静态成员代码示例 | [WriteAsync(IMemobus, String, Int16, FuncString, Int16, TaskOperateResult)](fcb2e4f3-8e61-75b5-c416-2be479b4e69d.htm) | 写入short数据，返回是否成功  Write short data, returns whether success |
| 公共方法静态成员代码示例 | [WriteAsync(IMemobus, String, UInt16, FuncString, UInt16, TaskOperateResult)](df1c6846-765e-58ee-cf5e-cc360a65f483.htm) | 写入ushort数据，返回是否成功  Write ushort data, return whether the write was successful |
| 公共方法静态成员 | [WriteRandom](53f86b45-0abc-47c2-a1c0-58247705efa1.htm) | 随机写入扩展的保持寄存器的内容，也即写入不连续的地址的字数据，字节数组的长度必须为地址数组长度的两倍，才能正确写入。  Write the contents of the extended hold registers randomly, that is, write word data for discontinuous addresses, and the byte array must be twice the length of the address array to be written correctly. |
| 公共方法静态成员 | [WriteRandomAsync](768d6601-f297-57f7-c3a9-2deed351fcbf.htm) | 随机写入扩展的保持寄存器的内容，也即写入不连续的地址的字数据，字节数组的长度必须为地址数组长度的两倍，才能正确写入。  Write the contents of the extended hold registers randomly, that is, write word data for discontinuous addresses, and the byte array must be twice the length of the address array to be written correctly. |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[MemobusHelper 类](254c50cc-c87c-57ad-eb8f-bda323fe9527.htm)

[HslCommunication.Profinet.YASKAWA.Helper 命名空间](310f2e91-5abf-4755-25f7-9225e316f6fe.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Read 方法 

[原文連結](http://api.hslcommunication.cn/html/edf5ccc5-454d-c522-2d86-234f53054551.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.YASKAWA.Helper](../html/310f2e91-5abf-4755-25f7-9225e316f6fe.htm "HslCommunication.Profinet.YASKAWA.Helper")

[MemobusHelper 类](../html/254c50cc-c87c-57ad-eb8f-bda323fe9527.htm "MemobusHelper 类")

[MemobusHelper 方法](../html/e43bdd34-e6b5-c55a-deed-0d20ed628891.htm "MemobusHelper 方法")

[Read 方法](../html/edf5ccc5-454d-c522-2d86-234f53054551.htm "Read 方法 ")

[ReadAsync 方法](../html/41d84bd7-1ae9-01bb-4813-5c50b2dfd419.htm "ReadAsync 方法 ")

[ReadBool 方法](../html/7f5d3261-4c7a-ef65-1226-f31cc4b3e4da.htm "ReadBool 方法 ")

[ReadBoolAsync 方法](../html/fbe714bf-2e73-6ffc-6aee-c45e3f583d43.htm "ReadBoolAsync 方法 ")

[ReadRandom 方法](../html/d19df820-e22b-5da4-280e-13c38e12cd37.htm "ReadRandom 方法 ")

[ReadRandomAsync 方法](../html/b237b1ea-3ff4-b07a-c31f-8e2686ff566b.htm "ReadRandomAsync 方法 ")

[Write 方法](../html/fb613f76-96a5-5566-eb37-f400dbd0a93c.htm "Write 方法 ")

[WriteAsync 方法](../html/bbc31e49-e3cd-881d-75b1-42ec61c72b34.htm "WriteAsync 方法 ")

[WriteRandom 方法](../html/53f86b45-0abc-47c2-a1c0-58247705efa1.htm "WriteRandom 方法 ")

[WriteRandomAsync 方法](../html/768d6601-f297-57f7-c3a9-2deed351fcbf.htm "WriteRandomAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MemobusHelperRead 方法 |

批量读取字节数组信息，需要指定地址和长度，返回原始的字节数组  
Batch read byte array information, need to specify the address and length, return the original byte array

**命名空间：**
 [HslCommunication.Profinet.YASKAWA.Helper](310f2e91-5abf-4755-25f7-9225e316f6fe.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult<byte[]> Read(
	IMemobus memobus,
	string address,
	ushort length
)
```

```
Public Shared Function Read ( 
	memobus As IMemobus,
	address As String,
	length As UShort
) As OperateResult(Of Byte())
```

```
public:
static OperateResult<array<unsigned char>^>^ Read(
	IMemobus^ memobus, 
	String^ address, 
	unsigned short length
)
```

```
static member Read : 
        memobus : IMemobus * 
        address : string * 
        length : uint16 -> OperateResult<byte[]> 
```

#### 参数

memobus
:   类型：[HslCommunication.Profinet.YASKAWA.HelperIMemobus](0f5ca9bf-1080-d913-ed93-7c9a46291f73.htm)  

    [缺少 "M:HslCommunication.Profinet.YASKAWA.Helper.MemobusHelper.Read(HslCommunication.Profinet.YASKAWA.Helper.IMemobus,System.String,System.UInt16)" 的 <param name="memobus"/> 文档]

address
:   类型：SystemString  
    数据地址

length
:   类型：SystemUInt16  
    数据长度

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
带有成功标识的byte[]数组

![](../icons/SectionExpanded.png)备注

地址默认使用功能码03，如果需要指定其他的功能码地址，需要手动指定功能码，例如：x=4;100, x=9;100, x=10;100, 当然也可以写成 x=0x0A;100  
The address uses function code 03 by default. If you need to specify other function code addresses,
you need to manually specify the function code, for example: x=4;100, x=9;100, x=10;100, of course, it can also be written as x=0x0A; 100

![](../icons/SectionExpanded.png)参见

#### 引用

[MemobusHelper 类](254c50cc-c87c-57ad-eb8f-bda323fe9527.htm)

[HslCommunication.Profinet.YASKAWA.Helper 命名空间](310f2e91-5abf-4755-25f7-9225e316f6fe.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/41d84bd7-1ae9-01bb-4813-5c50b2dfd419.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.YASKAWA.Helper](../html/310f2e91-5abf-4755-25f7-9225e316f6fe.htm "HslCommunication.Profinet.YASKAWA.Helper")

[MemobusHelper 类](../html/254c50cc-c87c-57ad-eb8f-bda323fe9527.htm "MemobusHelper 类")

[MemobusHelper 方法](../html/e43bdd34-e6b5-c55a-deed-0d20ed628891.htm "MemobusHelper 方法")

[Read 方法](../html/edf5ccc5-454d-c522-2d86-234f53054551.htm "Read 方法 ")

[ReadAsync 方法](../html/41d84bd7-1ae9-01bb-4813-5c50b2dfd419.htm "ReadAsync 方法 ")

[ReadBool 方法](../html/7f5d3261-4c7a-ef65-1226-f31cc4b3e4da.htm "ReadBool 方法 ")

[ReadBoolAsync 方法](../html/fbe714bf-2e73-6ffc-6aee-c45e3f583d43.htm "ReadBoolAsync 方法 ")

[ReadRandom 方法](../html/d19df820-e22b-5da4-280e-13c38e12cd37.htm "ReadRandom 方法 ")

[ReadRandomAsync 方法](../html/b237b1ea-3ff4-b07a-c31f-8e2686ff566b.htm "ReadRandomAsync 方法 ")

[Write 方法](../html/fb613f76-96a5-5566-eb37-f400dbd0a93c.htm "Write 方法 ")

[WriteAsync 方法](../html/bbc31e49-e3cd-881d-75b1-42ec61c72b34.htm "WriteAsync 方法 ")

[WriteRandom 方法](../html/53f86b45-0abc-47c2-a1c0-58247705efa1.htm "WriteRandom 方法 ")

[WriteRandomAsync 方法](../html/768d6601-f297-57f7-c3a9-2deed351fcbf.htm "WriteRandomAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MemobusHelperReadAsync 方法 |

批量读取字节数组信息，需要指定地址和长度，返回原始的字节数组  
Batch read byte array information, need to specify the address and length, return the original byte array

**命名空间：**
 [HslCommunication.Profinet.YASKAWA.Helper](310f2e91-5abf-4755-25f7-9225e316f6fe.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static Task<OperateResult<byte[]>> ReadAsync(
	IMemobus memobus,
	string address,
	ushort length
)
```

```
Public Shared Function ReadAsync ( 
	memobus As IMemobus,
	address As String,
	length As UShort
) As Task(Of OperateResult(Of Byte()))
```

```
public:
static Task<OperateResult<array<unsigned char>^>^>^ ReadAsync(
	IMemobus^ memobus, 
	String^ address, 
	unsigned short length
)
```

```
static member ReadAsync : 
        memobus : IMemobus * 
        address : string * 
        length : uint16 -> Task<OperateResult<byte[]>> 
```

#### 参数

memobus
:   类型：[HslCommunication.Profinet.YASKAWA.HelperIMemobus](0f5ca9bf-1080-d913-ed93-7c9a46291f73.htm)  

    [缺少 "M:HslCommunication.Profinet.YASKAWA.Helper.MemobusHelper.ReadAsync(HslCommunication.Profinet.YASKAWA.Helper.IMemobus,System.String,System.UInt16)" 的 <param name="memobus"/> 文档]

address
:   类型：SystemString  
    数据地址

length
:   类型：SystemUInt16  
    数据长度

#### 返回值

类型：Task[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
带有成功标识的byte[]数组

![](../icons/SectionExpanded.png)备注

地址默认使用功能码03，如果需要指定其他的功能码地址，需要手动指定功能码，例如：x=4;100, x=9;100, x=10;100, 当然也可以写成 x=0x0A;100  
The address uses function code 03 by default. If you need to specify other function code addresses,
you need to manually specify the function code, for example: x=4;100, x=9;100, x=10;100, of course, it can also be written as x=0x0A; 100

![](../icons/SectionExpanded.png)参见

#### 引用

[MemobusHelper 类](254c50cc-c87c-57ad-eb8f-bda323fe9527.htm)

[HslCommunication.Profinet.YASKAWA.Helper 命名空间](310f2e91-5abf-4755-25f7-9225e316f6fe.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadBool 方法 

[原文連結](http://api.hslcommunication.cn/html/7f5d3261-4c7a-ef65-1226-f31cc4b3e4da.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.YASKAWA.Helper](../html/310f2e91-5abf-4755-25f7-9225e316f6fe.htm "HslCommunication.Profinet.YASKAWA.Helper")

[MemobusHelper 类](../html/254c50cc-c87c-57ad-eb8f-bda323fe9527.htm "MemobusHelper 类")

[MemobusHelper 方法](../html/e43bdd34-e6b5-c55a-deed-0d20ed628891.htm "MemobusHelper 方法")

[Read 方法](../html/edf5ccc5-454d-c522-2d86-234f53054551.htm "Read 方法 ")

[ReadAsync 方法](../html/41d84bd7-1ae9-01bb-4813-5c50b2dfd419.htm "ReadAsync 方法 ")

[ReadBool 方法](../html/7f5d3261-4c7a-ef65-1226-f31cc4b3e4da.htm "ReadBool 方法 ")

[ReadBoolAsync 方法](../html/fbe714bf-2e73-6ffc-6aee-c45e3f583d43.htm "ReadBoolAsync 方法 ")

[ReadRandom 方法](../html/d19df820-e22b-5da4-280e-13c38e12cd37.htm "ReadRandom 方法 ")

[ReadRandomAsync 方法](../html/b237b1ea-3ff4-b07a-c31f-8e2686ff566b.htm "ReadRandomAsync 方法 ")

[Write 方法](../html/fb613f76-96a5-5566-eb37-f400dbd0a93c.htm "Write 方法 ")

[WriteAsync 方法](../html/bbc31e49-e3cd-881d-75b1-42ec61c72b34.htm "WriteAsync 方法 ")

[WriteRandom 方法](../html/53f86b45-0abc-47c2-a1c0-58247705efa1.htm "WriteRandom 方法 ")

[WriteRandomAsync 方法](../html/768d6601-f297-57f7-c3a9-2deed351fcbf.htm "WriteRandomAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MemobusHelperReadBool 方法 |

批量读取Boolean数组信息，需要指定地址和长度，返回Boolean 数组  
Batch read Boolean array information, need to specify the address and length, return Boolean array

**命名空间：**
 [HslCommunication.Profinet.YASKAWA.Helper](310f2e91-5abf-4755-25f7-9225e316f6fe.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult<bool[]> ReadBool(
	IMemobus memobus,
	string address,
	ushort length
)
```

```
Public Shared Function ReadBool ( 
	memobus As IMemobus,
	address As String,
	length As UShort
) As OperateResult(Of Boolean())
```

```
public:
static OperateResult<array<bool>^>^ ReadBool(
	IMemobus^ memobus, 
	String^ address, 
	unsigned short length
)
```

```
static member ReadBool : 
        memobus : IMemobus * 
        address : string * 
        length : uint16 -> OperateResult<bool[]> 
```

#### 参数

memobus
:   类型：[HslCommunication.Profinet.YASKAWA.HelperIMemobus](0f5ca9bf-1080-d913-ed93-7c9a46291f73.htm)  

    [缺少 "M:HslCommunication.Profinet.YASKAWA.Helper.MemobusHelper.ReadBool(HslCommunication.Profinet.YASKAWA.Helper.IMemobus,System.String,System.UInt16)" 的 <param name="memobus"/> 文档]

address
:   类型：SystemString  
    数据地址

length
:   类型：SystemUInt16  
    数据长度

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Boolean  
带有成功标识的 bool[] 数组

![](../icons/SectionExpanded.png)备注

默认使用功能码01，读取线圈操作，如果需要指定读取输入线圈，地址需要携带额外的参数，例如 x=2;100  
The function code 01 is used by default to read the coil operation. If you need to specify the read input coil, the address needs to carry additional parameters, such as x=2;100

![](../icons/SectionExpanded.png)参见

#### 引用

[MemobusHelper 类](254c50cc-c87c-57ad-eb8f-bda323fe9527.htm)

[HslCommunication.Profinet.YASKAWA.Helper 命名空间](310f2e91-5abf-4755-25f7-9225e316f6fe.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadBoolAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/fbe714bf-2e73-6ffc-6aee-c45e3f583d43.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.YASKAWA.Helper](../html/310f2e91-5abf-4755-25f7-9225e316f6fe.htm "HslCommunication.Profinet.YASKAWA.Helper")

[MemobusHelper 类](../html/254c50cc-c87c-57ad-eb8f-bda323fe9527.htm "MemobusHelper 类")

[MemobusHelper 方法](../html/e43bdd34-e6b5-c55a-deed-0d20ed628891.htm "MemobusHelper 方法")

[Read 方法](../html/edf5ccc5-454d-c522-2d86-234f53054551.htm "Read 方法 ")

[ReadAsync 方法](../html/41d84bd7-1ae9-01bb-4813-5c50b2dfd419.htm "ReadAsync 方法 ")

[ReadBool 方法](../html/7f5d3261-4c7a-ef65-1226-f31cc4b3e4da.htm "ReadBool 方法 ")

[ReadBoolAsync 方法](../html/fbe714bf-2e73-6ffc-6aee-c45e3f583d43.htm "ReadBoolAsync 方法 ")

[ReadRandom 方法](../html/d19df820-e22b-5da4-280e-13c38e12cd37.htm "ReadRandom 方法 ")

[ReadRandomAsync 方法](../html/b237b1ea-3ff4-b07a-c31f-8e2686ff566b.htm "ReadRandomAsync 方法 ")

[Write 方法](../html/fb613f76-96a5-5566-eb37-f400dbd0a93c.htm "Write 方法 ")

[WriteAsync 方法](../html/bbc31e49-e3cd-881d-75b1-42ec61c72b34.htm "WriteAsync 方法 ")

[WriteRandom 方法](../html/53f86b45-0abc-47c2-a1c0-58247705efa1.htm "WriteRandom 方法 ")

[WriteRandomAsync 方法](../html/768d6601-f297-57f7-c3a9-2deed351fcbf.htm "WriteRandomAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MemobusHelperReadBoolAsync 方法 |

批量读取Boolean数组信息，需要指定地址和长度，返回Boolean 数组  
Batch read Boolean array information, need to specify the address and length, return Boolean array

**命名空间：**
 [HslCommunication.Profinet.YASKAWA.Helper](310f2e91-5abf-4755-25f7-9225e316f6fe.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static Task<OperateResult<bool[]>> ReadBoolAsync(
	IMemobus memobus,
	string address,
	ushort length
)
```

```
Public Shared Function ReadBoolAsync ( 
	memobus As IMemobus,
	address As String,
	length As UShort
) As Task(Of OperateResult(Of Boolean()))
```

```
public:
static Task<OperateResult<array<bool>^>^>^ ReadBoolAsync(
	IMemobus^ memobus, 
	String^ address, 
	unsigned short length
)
```

```
static member ReadBoolAsync : 
        memobus : IMemobus * 
        address : string * 
        length : uint16 -> Task<OperateResult<bool[]>> 
```

#### 参数

memobus
:   类型：[HslCommunication.Profinet.YASKAWA.HelperIMemobus](0f5ca9bf-1080-d913-ed93-7c9a46291f73.htm)  

    [缺少 "M:HslCommunication.Profinet.YASKAWA.Helper.MemobusHelper.ReadBoolAsync(HslCommunication.Profinet.YASKAWA.Helper.IMemobus,System.String,System.UInt16)" 的 <param name="memobus"/> 文档]

address
:   类型：SystemString  
    数据地址

length
:   类型：SystemUInt16  
    数据长度

#### 返回值

类型：Task[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Boolean  
带有成功标识的 bool[] 数组

![](../icons/SectionExpanded.png)备注

默认使用功能码01，读取线圈操作，如果需要指定读取输入线圈，地址需要携带额外的参数，例如 x=2;100  
The function code 01 is used by default to read the coil operation. If you need to specify the read input coil, the address needs to carry additional parameters, such as x=2;100

![](../icons/SectionExpanded.png)参见

#### 引用

[MemobusHelper 类](254c50cc-c87c-57ad-eb8f-bda323fe9527.htm)

[HslCommunication.Profinet.YASKAWA.Helper 命名空间](310f2e91-5abf-4755-25f7-9225e316f6fe.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadRandom 方法 

[原文連結](http://api.hslcommunication.cn/html/d19df820-e22b-5da4-280e-13c38e12cd37.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.YASKAWA.Helper](../html/310f2e91-5abf-4755-25f7-9225e316f6fe.htm "HslCommunication.Profinet.YASKAWA.Helper")

[MemobusHelper 类](../html/254c50cc-c87c-57ad-eb8f-bda323fe9527.htm "MemobusHelper 类")

[MemobusHelper 方法](../html/e43bdd34-e6b5-c55a-deed-0d20ed628891.htm "MemobusHelper 方法")

[ReadRandom 方法](../html/d19df820-e22b-5da4-280e-13c38e12cd37.htm "ReadRandom 方法 ")

[ReadRandom 方法 (IMemobus, String[])](../html/bdf0757b-18b2-0b9e-366a-a012a5c88955.htm "ReadRandom 方法 (IMemobus, String[])")

[ReadRandom 方法 (IMemobus, UInt16[])](../html/26f4b67e-78ef-6494-615a-9db55bfb15ae.htm "ReadRandom 方法 (IMemobus, UInt16[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MemobusHelperReadRandom 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法静态成员 | [ReadRandom(IMemobus, String)](bdf0757b-18b2-0b9e-366a-a012a5c88955.htm) | 随机读取扩展的保持寄存器的内容，也即读取不连续地址的字数据，可以指定多个地址，然后一次性读取所有的数据，然后解析出实际的数据  Randomly read the contents of the extended hold register, that is, read word data of discontinuous addresses, you can specify multiple addresses, then read all the data at once, and then parse out the actual data |
| 公共方法静态成员 | [ReadRandom(IMemobus, UInt16)](26f4b67e-78ef-6494-615a-9db55bfb15ae.htm) | 随机读取扩展的保持寄存器的内容，也即读取不连续地址的字数据，可以指定多个地址，然后一次性读取所有的数据，然后解析出实际的数据  Randomly read the contents of the extended hold register, that is, read word data of discontinuous addresses, you can specify multiple addresses, then read all the data at once, and then parse out the actual data |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[MemobusHelper 类](254c50cc-c87c-57ad-eb8f-bda323fe9527.htm)

[HslCommunication.Profinet.YASKAWA.Helper 命名空间](310f2e91-5abf-4755-25f7-9225e316f6fe.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadRandom 方法 (IMemobus, String[])

[原文連結](http://api.hslcommunication.cn/html/bdf0757b-18b2-0b9e-366a-a012a5c88955.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.YASKAWA.Helper](../html/310f2e91-5abf-4755-25f7-9225e316f6fe.htm "HslCommunication.Profinet.YASKAWA.Helper")

[MemobusHelper 类](../html/254c50cc-c87c-57ad-eb8f-bda323fe9527.htm "MemobusHelper 类")

[MemobusHelper 方法](../html/e43bdd34-e6b5-c55a-deed-0d20ed628891.htm "MemobusHelper 方法")

[ReadRandom 方法](../html/d19df820-e22b-5da4-280e-13c38e12cd37.htm "ReadRandom 方法 ")

[ReadRandom 方法 (IMemobus, String[])](../html/bdf0757b-18b2-0b9e-366a-a012a5c88955.htm "ReadRandom 方法 (IMemobus, String[])")

[ReadRandom 方法 (IMemobus, UInt16[])](../html/26f4b67e-78ef-6494-615a-9db55bfb15ae.htm "ReadRandom 方法 (IMemobus, UInt16[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MemobusHelperReadRandom 方法 (IMemobus, String) |

随机读取扩展的保持寄存器的内容，也即读取不连续地址的字数据，可以指定多个地址，然后一次性读取所有的数据，然后解析出实际的数据  
Randomly read the contents of the extended hold register, that is, read word data of discontinuous addresses,
you can specify multiple addresses, then read all the data at once, and then parse out the actual data

**命名空间：**
 [HslCommunication.Profinet.YASKAWA.Helper](310f2e91-5abf-4755-25f7-9225e316f6fe.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult<byte[]> ReadRandom(
	IMemobus memobus,
	string[] address
)
```

```
Public Shared Function ReadRandom ( 
	memobus As IMemobus,
	address As String()
) As OperateResult(Of Byte())
```

```
public:
static OperateResult<array<unsigned char>^>^ ReadRandom(
	IMemobus^ memobus, 
	array<String^>^ address
)
```

```
static member ReadRandom : 
        memobus : IMemobus * 
        address : string[] -> OperateResult<byte[]> 
```

#### 参数

memobus
:   类型：[HslCommunication.Profinet.YASKAWA.HelperIMemobus](0f5ca9bf-1080-d913-ed93-7c9a46291f73.htm)  
    PLC通信对象

address
:   类型：SystemString  
    地址信息

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
读取的原始字节结果信息

![](../icons/SectionExpanded.png)备注

本方法的地址支持 保持寄存器 M100, 数据寄存器 G100, 输入寄存器 I100, 输出寄存器 O100, 系统寄存器 S100,

![](../icons/SectionExpanded.png)参见

#### 引用

[MemobusHelper 类](254c50cc-c87c-57ad-eb8f-bda323fe9527.htm)

[ReadRandom 重载](d19df820-e22b-5da4-280e-13c38e12cd37.htm)

[HslCommunication.Profinet.YASKAWA.Helper 命名空间](310f2e91-5abf-4755-25f7-9225e316f6fe.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadRandom 方法 (IMemobus, UInt16[])

[原文連結](http://api.hslcommunication.cn/html/26f4b67e-78ef-6494-615a-9db55bfb15ae.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.YASKAWA.Helper](../html/310f2e91-5abf-4755-25f7-9225e316f6fe.htm "HslCommunication.Profinet.YASKAWA.Helper")

[MemobusHelper 类](../html/254c50cc-c87c-57ad-eb8f-bda323fe9527.htm "MemobusHelper 类")

[MemobusHelper 方法](../html/e43bdd34-e6b5-c55a-deed-0d20ed628891.htm "MemobusHelper 方法")

[ReadRandom 方法](../html/d19df820-e22b-5da4-280e-13c38e12cd37.htm "ReadRandom 方法 ")

[ReadRandom 方法 (IMemobus, String[])](../html/bdf0757b-18b2-0b9e-366a-a012a5c88955.htm "ReadRandom 方法 (IMemobus, String[])")

[ReadRandom 方法 (IMemobus, UInt16[])](../html/26f4b67e-78ef-6494-615a-9db55bfb15ae.htm "ReadRandom 方法 (IMemobus, UInt16[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MemobusHelperReadRandom 方法 (IMemobus, UInt16) |

随机读取扩展的保持寄存器的内容，也即读取不连续地址的字数据，可以指定多个地址，然后一次性读取所有的数据，然后解析出实际的数据  
Randomly read the contents of the extended hold register, that is, read word data of discontinuous addresses,
you can specify multiple addresses, then read all the data at once, and then parse out the actual data

**命名空间：**
 [HslCommunication.Profinet.YASKAWA.Helper](310f2e91-5abf-4755-25f7-9225e316f6fe.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult<byte[]> ReadRandom(
	IMemobus memobus,
	ushort[] address
)
```

```
Public Shared Function ReadRandom ( 
	memobus As IMemobus,
	address As UShort()
) As OperateResult(Of Byte())
```

```
public:
static OperateResult<array<unsigned char>^>^ ReadRandom(
	IMemobus^ memobus, 
	array<unsigned short>^ address
)
```

```
static member ReadRandom : 
        memobus : IMemobus * 
        address : uint16[] -> OperateResult<byte[]> 
```

#### 参数

memobus
:   类型：[HslCommunication.Profinet.YASKAWA.HelperIMemobus](0f5ca9bf-1080-d913-ed93-7c9a46291f73.htm)  
    PLC通信对象

address
:   类型：SystemUInt16  
    地址信息

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
读取的原始字节结果信息

![](../icons/SectionExpanded.png)备注

注意，本方法只能针对扩展的保持寄存器进行读取

![](../icons/SectionExpanded.png)参见

#### 引用

[MemobusHelper 类](254c50cc-c87c-57ad-eb8f-bda323fe9527.htm)

[ReadRandom 重载](d19df820-e22b-5da4-280e-13c38e12cd37.htm)

[HslCommunication.Profinet.YASKAWA.Helper 命名空间](310f2e91-5abf-4755-25f7-9225e316f6fe.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadRandomAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/b237b1ea-3ff4-b07a-c31f-8e2686ff566b.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.YASKAWA.Helper](../html/310f2e91-5abf-4755-25f7-9225e316f6fe.htm "HslCommunication.Profinet.YASKAWA.Helper")

[MemobusHelper 类](../html/254c50cc-c87c-57ad-eb8f-bda323fe9527.htm "MemobusHelper 类")

[MemobusHelper 方法](../html/e43bdd34-e6b5-c55a-deed-0d20ed628891.htm "MemobusHelper 方法")

[ReadRandomAsync 方法](../html/b237b1ea-3ff4-b07a-c31f-8e2686ff566b.htm "ReadRandomAsync 方法 ")

[ReadRandomAsync 方法 (IMemobus, String[])](../html/0edac1c0-22c9-5164-527f-d4533556bf61.htm "ReadRandomAsync 方法 (IMemobus, String[])")

[ReadRandomAsync 方法 (IMemobus, UInt16[])](../html/8aa741a5-4c85-127e-d8be-6677ba4be1cf.htm "ReadRandomAsync 方法 (IMemobus, UInt16[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MemobusHelperReadRandomAsync 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法静态成员 | [ReadRandomAsync(IMemobus, String)](0edac1c0-22c9-5164-527f-d4533556bf61.htm) | 随机读取扩展的保持寄存器的内容，也即读取不连续地址的字数据，可以指定多个地址，然后一次性读取所有的数据，然后解析出实际的数据  Randomly read the contents of the extended hold register, that is, read word data of discontinuous addresses, you can specify multiple addresses, then read all the data at once, and then parse out the actual data |
| 公共方法静态成员 | [ReadRandomAsync(IMemobus, UInt16)](8aa741a5-4c85-127e-d8be-6677ba4be1cf.htm) | 随机读取扩展的保持寄存器的内容，也即读取不连续地址的字数据，可以指定多个地址，然后一次性读取所有的数据，然后解析出实际的数据  Randomly read the contents of the extended hold register, that is, read word data of discontinuous addresses, you can specify multiple addresses, then read all the data at once, and then parse out the actual data |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[MemobusHelper 类](254c50cc-c87c-57ad-eb8f-bda323fe9527.htm)

[HslCommunication.Profinet.YASKAWA.Helper 命名空间](310f2e91-5abf-4755-25f7-9225e316f6fe.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadRandomAsync 方法 (IMemobus, String[])

[原文連結](http://api.hslcommunication.cn/html/0edac1c0-22c9-5164-527f-d4533556bf61.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.YASKAWA.Helper](../html/310f2e91-5abf-4755-25f7-9225e316f6fe.htm "HslCommunication.Profinet.YASKAWA.Helper")

[MemobusHelper 类](../html/254c50cc-c87c-57ad-eb8f-bda323fe9527.htm "MemobusHelper 类")

[MemobusHelper 方法](../html/e43bdd34-e6b5-c55a-deed-0d20ed628891.htm "MemobusHelper 方法")

[ReadRandomAsync 方法](../html/b237b1ea-3ff4-b07a-c31f-8e2686ff566b.htm "ReadRandomAsync 方法 ")

[ReadRandomAsync 方法 (IMemobus, String[])](../html/0edac1c0-22c9-5164-527f-d4533556bf61.htm "ReadRandomAsync 方法 (IMemobus, String[])")

[ReadRandomAsync 方法 (IMemobus, UInt16[])](../html/8aa741a5-4c85-127e-d8be-6677ba4be1cf.htm "ReadRandomAsync 方法 (IMemobus, UInt16[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MemobusHelperReadRandomAsync 方法 (IMemobus, String) |

随机读取扩展的保持寄存器的内容，也即读取不连续地址的字数据，可以指定多个地址，然后一次性读取所有的数据，然后解析出实际的数据  
Randomly read the contents of the extended hold register, that is, read word data of discontinuous addresses,
you can specify multiple addresses, then read all the data at once, and then parse out the actual data

**命名空间：**
 [HslCommunication.Profinet.YASKAWA.Helper](310f2e91-5abf-4755-25f7-9225e316f6fe.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static Task<OperateResult<byte[]>> ReadRandomAsync(
	IMemobus memobus,
	string[] address
)
```

```
Public Shared Function ReadRandomAsync ( 
	memobus As IMemobus,
	address As String()
) As Task(Of OperateResult(Of Byte()))
```

```
public:
static Task<OperateResult<array<unsigned char>^>^>^ ReadRandomAsync(
	IMemobus^ memobus, 
	array<String^>^ address
)
```

```
static member ReadRandomAsync : 
        memobus : IMemobus * 
        address : string[] -> Task<OperateResult<byte[]>> 
```

#### 参数

memobus
:   类型：[HslCommunication.Profinet.YASKAWA.HelperIMemobus](0f5ca9bf-1080-d913-ed93-7c9a46291f73.htm)  
    PLC通信对象

address
:   类型：SystemString  
    地址信息

#### 返回值

类型：Task[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
读取的原始字节结果信息

![](../icons/SectionExpanded.png)备注

注意，本方法只能针对扩展的保持寄存器进行读取

![](../icons/SectionExpanded.png)参见

#### 引用

[MemobusHelper 类](254c50cc-c87c-57ad-eb8f-bda323fe9527.htm)

[ReadRandomAsync 重载](b237b1ea-3ff4-b07a-c31f-8e2686ff566b.htm)

[HslCommunication.Profinet.YASKAWA.Helper 命名空间](310f2e91-5abf-4755-25f7-9225e316f6fe.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadRandomAsync 方法 (IMemobus, UInt16[])

[原文連結](http://api.hslcommunication.cn/html/8aa741a5-4c85-127e-d8be-6677ba4be1cf.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.YASKAWA.Helper](../html/310f2e91-5abf-4755-25f7-9225e316f6fe.htm "HslCommunication.Profinet.YASKAWA.Helper")

[MemobusHelper 类](../html/254c50cc-c87c-57ad-eb8f-bda323fe9527.htm "MemobusHelper 类")

[MemobusHelper 方法](../html/e43bdd34-e6b5-c55a-deed-0d20ed628891.htm "MemobusHelper 方法")

[ReadRandomAsync 方法](../html/b237b1ea-3ff4-b07a-c31f-8e2686ff566b.htm "ReadRandomAsync 方法 ")

[ReadRandomAsync 方法 (IMemobus, String[])](../html/0edac1c0-22c9-5164-527f-d4533556bf61.htm "ReadRandomAsync 方法 (IMemobus, String[])")

[ReadRandomAsync 方法 (IMemobus, UInt16[])](../html/8aa741a5-4c85-127e-d8be-6677ba4be1cf.htm "ReadRandomAsync 方法 (IMemobus, UInt16[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MemobusHelperReadRandomAsync 方法 (IMemobus, UInt16) |

随机读取扩展的保持寄存器的内容，也即读取不连续地址的字数据，可以指定多个地址，然后一次性读取所有的数据，然后解析出实际的数据  
Randomly read the contents of the extended hold register, that is, read word data of discontinuous addresses,
you can specify multiple addresses, then read all the data at once, and then parse out the actual data

**命名空间：**
 [HslCommunication.Profinet.YASKAWA.Helper](310f2e91-5abf-4755-25f7-9225e316f6fe.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static Task<OperateResult<byte[]>> ReadRandomAsync(
	IMemobus memobus,
	ushort[] address
)
```

```
Public Shared Function ReadRandomAsync ( 
	memobus As IMemobus,
	address As UShort()
) As Task(Of OperateResult(Of Byte()))
```

```
public:
static Task<OperateResult<array<unsigned char>^>^>^ ReadRandomAsync(
	IMemobus^ memobus, 
	array<unsigned short>^ address
)
```

```
static member ReadRandomAsync : 
        memobus : IMemobus * 
        address : uint16[] -> Task<OperateResult<byte[]>> 
```

#### 参数

memobus
:   类型：[HslCommunication.Profinet.YASKAWA.HelperIMemobus](0f5ca9bf-1080-d913-ed93-7c9a46291f73.htm)  
    PLC通信对象

address
:   类型：SystemUInt16  
    地址信息

#### 返回值

类型：Task[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
读取的原始字节结果信息

![](../icons/SectionExpanded.png)备注

注意，本方法只能针对扩展的保持寄存器进行读取

![](../icons/SectionExpanded.png)参见

#### 引用

[MemobusHelper 类](254c50cc-c87c-57ad-eb8f-bda323fe9527.htm)

[ReadRandomAsync 重载](b237b1ea-3ff4-b07a-c31f-8e2686ff566b.htm)

[HslCommunication.Profinet.YASKAWA.Helper 命名空间](310f2e91-5abf-4755-25f7-9225e316f6fe.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Write 方法 

[原文連結](http://api.hslcommunication.cn/html/fb613f76-96a5-5566-eb37-f400dbd0a93c.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.YASKAWA.Helper](../html/310f2e91-5abf-4755-25f7-9225e316f6fe.htm "HslCommunication.Profinet.YASKAWA.Helper")

[MemobusHelper 类](../html/254c50cc-c87c-57ad-eb8f-bda323fe9527.htm "MemobusHelper 类")

[MemobusHelper 方法](../html/e43bdd34-e6b5-c55a-deed-0d20ed628891.htm "MemobusHelper 方法")

[Write 方法](../html/fb613f76-96a5-5566-eb37-f400dbd0a93c.htm "Write 方法 ")

[Write 方法 (IMemobus, String, Boolean)](../html/1f35cb0d-f4ae-b48b-6830-6f58f66e94d5.htm "Write 方法 (IMemobus, String, Boolean)")

[Write 方法 (IMemobus, String, Boolean[])](../html/7cf177a0-b8e8-78da-fb1f-dcd8b1bec795.htm "Write 方法 (IMemobus, String, Boolean[])")

[Write 方法 (IMemobus, String, Byte[])](../html/a304cab6-f89e-e9e8-ffb6-8906b7a15382.htm "Write 方法 (IMemobus, String, Byte[])")

[Write 方法 (IMemobus, String, Int16, Func(String, Int16, OperateResult))](../html/ac26cdc6-3666-0f1e-0d3c-245e81c57aff.htm "Write 方法 (IMemobus, String, Int16, Func(String, Int16, OperateResult))")

[Write 方法 (IMemobus, String, UInt16, Func(String, UInt16, OperateResult))](../html/b37601c7-5f78-2779-1479-d6e47642f8d1.htm "Write 方法 (IMemobus, String, UInt16, Func(String, UInt16, OperateResult))")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MemobusHelperWrite 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法静态成员 | [Write(IMemobus, String, Boolean)](1f35cb0d-f4ae-b48b-6830-6f58f66e94d5.htm) | 写入单个的Boolean数据，返回是否成功  Write a single Boolean data, and return whether the write was successful |
| 公共方法静态成员 | [Write(IMemobus, String, Boolean)](7cf177a0-b8e8-78da-fb1f-dcd8b1bec795.htm) | 批量写入Boolean数组数据，返回是否成功  Batch write Boolean array data, return whether the write was successful |
| 公共方法静态成员 | [Write(IMemobus, String, Byte)](a304cab6-f89e-e9e8-ffb6-8906b7a15382.htm) | 写入原始的byte数组数据到指定的地址，返回是否写入成功  Write the original byte array data to the specified address, and return whether the write was successful |
| 公共方法静态成员代码示例 | [Write(IMemobus, String, Int16, FuncString, Int16, OperateResult)](ac26cdc6-3666-0f1e-0d3c-245e81c57aff.htm) | 写入short数据，返回是否成功  Write short data, returns whether success |
| 公共方法静态成员代码示例 | [Write(IMemobus, String, UInt16, FuncString, UInt16, OperateResult)](b37601c7-5f78-2779-1479-d6e47642f8d1.htm) | 写入ushort数据，返回是否成功  Write ushort data, return whether the write was successful |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[MemobusHelper 类](254c50cc-c87c-57ad-eb8f-bda323fe9527.htm)

[HslCommunication.Profinet.YASKAWA.Helper 命名空间](310f2e91-5abf-4755-25f7-9225e316f6fe.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Write 方法 (IMemobus, String, Boolean)

[原文連結](http://api.hslcommunication.cn/html/1f35cb0d-f4ae-b48b-6830-6f58f66e94d5.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.YASKAWA.Helper](../html/310f2e91-5abf-4755-25f7-9225e316f6fe.htm "HslCommunication.Profinet.YASKAWA.Helper")

[MemobusHelper 类](../html/254c50cc-c87c-57ad-eb8f-bda323fe9527.htm "MemobusHelper 类")

[MemobusHelper 方法](../html/e43bdd34-e6b5-c55a-deed-0d20ed628891.htm "MemobusHelper 方法")

[Write 方法](../html/fb613f76-96a5-5566-eb37-f400dbd0a93c.htm "Write 方法 ")

[Write 方法 (IMemobus, String, Boolean)](../html/1f35cb0d-f4ae-b48b-6830-6f58f66e94d5.htm "Write 方法 (IMemobus, String, Boolean)")

[Write 方法 (IMemobus, String, Boolean[])](../html/7cf177a0-b8e8-78da-fb1f-dcd8b1bec795.htm "Write 方法 (IMemobus, String, Boolean[])")

[Write 方法 (IMemobus, String, Byte[])](../html/a304cab6-f89e-e9e8-ffb6-8906b7a15382.htm "Write 方法 (IMemobus, String, Byte[])")

[Write 方法 (IMemobus, String, Int16, Func(String, Int16, OperateResult))](../html/ac26cdc6-3666-0f1e-0d3c-245e81c57aff.htm "Write 方法 (IMemobus, String, Int16, Func(String, Int16, OperateResult))")

[Write 方法 (IMemobus, String, UInt16, Func(String, UInt16, OperateResult))](../html/b37601c7-5f78-2779-1479-d6e47642f8d1.htm "Write 方法 (IMemobus, String, UInt16, Func(String, UInt16, OperateResult))")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MemobusHelperWrite 方法 (IMemobus, String, Boolean) |

写入单个的Boolean数据，返回是否成功  
Write a single Boolean data, and return whether the write was successful

**命名空间：**
 [HslCommunication.Profinet.YASKAWA.Helper](310f2e91-5abf-4755-25f7-9225e316f6fe.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult Write(
	IMemobus memobus,
	string address,
	bool value
)
```

```
Public Shared Function Write ( 
	memobus As IMemobus,
	address As String,
	value As Boolean
) As OperateResult
```

```
public:
static OperateResult^ Write(
	IMemobus^ memobus, 
	String^ address, 
	bool value
)
```

```
static member Write : 
        memobus : IMemobus * 
        address : string * 
        value : bool -> OperateResult 
```

#### 参数

memobus
:   类型：[HslCommunication.Profinet.YASKAWA.HelperIMemobus](0f5ca9bf-1080-d913-ed93-7c9a46291f73.htm)  

    [缺少 "M:HslCommunication.Profinet.YASKAWA.Helper.MemobusHelper.Write(HslCommunication.Profinet.YASKAWA.Helper.IMemobus,System.String,System.Boolean)" 的 <param name="memobus"/> 文档]

address
:   类型：SystemString  
    起始地址

value
:   类型：SystemBoolean  
    写入值

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否写入成功

![](../icons/SectionExpanded.png)备注

单一线圈的状态变更，使用的主功能码为0x20, 子功能码为0x05  
The status of a single coil is changed, the main function code used is 0x20, and the sub function code is 0x05

![](../icons/SectionExpanded.png)参见

#### 引用

[MemobusHelper 类](254c50cc-c87c-57ad-eb8f-bda323fe9527.htm)

[Write 重载](fb613f76-96a5-5566-eb37-f400dbd0a93c.htm)

[HslCommunication.Profinet.YASKAWA.Helper 命名空间](310f2e91-5abf-4755-25f7-9225e316f6fe.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Write 方法 (IMemobus, String, Boolean[])

[原文連結](http://api.hslcommunication.cn/html/7cf177a0-b8e8-78da-fb1f-dcd8b1bec795.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.YASKAWA.Helper](../html/310f2e91-5abf-4755-25f7-9225e316f6fe.htm "HslCommunication.Profinet.YASKAWA.Helper")

[MemobusHelper 类](../html/254c50cc-c87c-57ad-eb8f-bda323fe9527.htm "MemobusHelper 类")

[MemobusHelper 方法](../html/e43bdd34-e6b5-c55a-deed-0d20ed628891.htm "MemobusHelper 方法")

[Write 方法](../html/fb613f76-96a5-5566-eb37-f400dbd0a93c.htm "Write 方法 ")

[Write 方法 (IMemobus, String, Boolean)](../html/1f35cb0d-f4ae-b48b-6830-6f58f66e94d5.htm "Write 方法 (IMemobus, String, Boolean)")

[Write 方法 (IMemobus, String, Boolean[])](../html/7cf177a0-b8e8-78da-fb1f-dcd8b1bec795.htm "Write 方法 (IMemobus, String, Boolean[])")

[Write 方法 (IMemobus, String, Byte[])](../html/a304cab6-f89e-e9e8-ffb6-8906b7a15382.htm "Write 方法 (IMemobus, String, Byte[])")

[Write 方法 (IMemobus, String, Int16, Func(String, Int16, OperateResult))](../html/ac26cdc6-3666-0f1e-0d3c-245e81c57aff.htm "Write 方法 (IMemobus, String, Int16, Func(String, Int16, OperateResult))")

[Write 方法 (IMemobus, String, UInt16, Func(String, UInt16, OperateResult))](../html/b37601c7-5f78-2779-1479-d6e47642f8d1.htm "Write 方法 (IMemobus, String, UInt16, Func(String, UInt16, OperateResult))")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MemobusHelperWrite 方法 (IMemobus, String, Boolean) |

批量写入Boolean数组数据，返回是否成功  
Batch write Boolean array data, return whether the write was successful

**命名空间：**
 [HslCommunication.Profinet.YASKAWA.Helper](310f2e91-5abf-4755-25f7-9225e316f6fe.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult Write(
	IMemobus memobus,
	string address,
	bool[] value
)
```

```
Public Shared Function Write ( 
	memobus As IMemobus,
	address As String,
	value As Boolean()
) As OperateResult
```

```
public:
static OperateResult^ Write(
	IMemobus^ memobus, 
	String^ address, 
	array<bool>^ value
)
```

```
static member Write : 
        memobus : IMemobus * 
        address : string * 
        value : bool[] -> OperateResult 
```

#### 参数

memobus
:   类型：[HslCommunication.Profinet.YASKAWA.HelperIMemobus](0f5ca9bf-1080-d913-ed93-7c9a46291f73.htm)  

    [缺少 "M:HslCommunication.Profinet.YASKAWA.Helper.MemobusHelper.Write(HslCommunication.Profinet.YASKAWA.Helper.IMemobus,System.String,System.Boolean[])" 的 <param name="memobus"/> 文档]

address
:   类型：SystemString  
    起始地址

value
:   类型：SystemBoolean  
    写入值

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否写入成功

![](../icons/SectionExpanded.png)备注

多个线圈的状态更改，默认使用的是 0x0f 子功能码。  
The status of multiple coils is changed, and the sub-function code 0x0f is used by default.

![](../icons/SectionExpanded.png)参见

#### 引用

[MemobusHelper 类](254c50cc-c87c-57ad-eb8f-bda323fe9527.htm)

[Write 重载](fb613f76-96a5-5566-eb37-f400dbd0a93c.htm)

[HslCommunication.Profinet.YASKAWA.Helper 命名空间](310f2e91-5abf-4755-25f7-9225e316f6fe.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Write 方法 (IMemobus, String, Byte[])

[原文連結](http://api.hslcommunication.cn/html/a304cab6-f89e-e9e8-ffb6-8906b7a15382.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.YASKAWA.Helper](../html/310f2e91-5abf-4755-25f7-9225e316f6fe.htm "HslCommunication.Profinet.YASKAWA.Helper")

[MemobusHelper 类](../html/254c50cc-c87c-57ad-eb8f-bda323fe9527.htm "MemobusHelper 类")

[MemobusHelper 方法](../html/e43bdd34-e6b5-c55a-deed-0d20ed628891.htm "MemobusHelper 方法")

[Write 方法](../html/fb613f76-96a5-5566-eb37-f400dbd0a93c.htm "Write 方法 ")

[Write 方法 (IMemobus, String, Boolean)](../html/1f35cb0d-f4ae-b48b-6830-6f58f66e94d5.htm "Write 方法 (IMemobus, String, Boolean)")

[Write 方法 (IMemobus, String, Boolean[])](../html/7cf177a0-b8e8-78da-fb1f-dcd8b1bec795.htm "Write 方法 (IMemobus, String, Boolean[])")

[Write 方法 (IMemobus, String, Byte[])](../html/a304cab6-f89e-e9e8-ffb6-8906b7a15382.htm "Write 方法 (IMemobus, String, Byte[])")

[Write 方法 (IMemobus, String, Int16, Func(String, Int16, OperateResult))](../html/ac26cdc6-3666-0f1e-0d3c-245e81c57aff.htm "Write 方法 (IMemobus, String, Int16, Func(String, Int16, OperateResult))")

[Write 方法 (IMemobus, String, UInt16, Func(String, UInt16, OperateResult))](../html/b37601c7-5f78-2779-1479-d6e47642f8d1.htm "Write 方法 (IMemobus, String, UInt16, Func(String, UInt16, OperateResult))")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MemobusHelperWrite 方法 (IMemobus, String, Byte) |

写入原始的byte数组数据到指定的地址，返回是否写入成功  
Write the original byte array data to the specified address, and return whether the write was successful

**命名空间：**
 [HslCommunication.Profinet.YASKAWA.Helper](310f2e91-5abf-4755-25f7-9225e316f6fe.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult Write(
	IMemobus memobus,
	string address,
	byte[] value
)
```

```
Public Shared Function Write ( 
	memobus As IMemobus,
	address As String,
	value As Byte()
) As OperateResult
```

```
public:
static OperateResult^ Write(
	IMemobus^ memobus, 
	String^ address, 
	array<unsigned char>^ value
)
```

```
static member Write : 
        memobus : IMemobus * 
        address : string * 
        value : byte[] -> OperateResult 
```

#### 参数

memobus
:   类型：[HslCommunication.Profinet.YASKAWA.HelperIMemobus](0f5ca9bf-1080-d913-ed93-7c9a46291f73.htm)  

    [缺少 "M:HslCommunication.Profinet.YASKAWA.Helper.MemobusHelper.Write(HslCommunication.Profinet.YASKAWA.Helper.IMemobus,System.String,System.Byte[])" 的 <param name="memobus"/> 文档]

address
:   类型：SystemString  
    起始地址

value
:   类型：SystemByte  
    写入值

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否写入成功

![](../icons/SectionExpanded.png)备注

连续的寄存器写入操作，默认功能码是0x10，如果需要写入扩展的寄存器，使用 x=0xA;100 或是 x=10;100 即可。  
For continuous register write operation, the default function code is 0x10. If you need to write an extended register, use x=0xA;100 or x=10;100.

![](../icons/SectionExpanded.png)参见

#### 引用

[MemobusHelper 类](254c50cc-c87c-57ad-eb8f-bda323fe9527.htm)

[Write 重载](fb613f76-96a5-5566-eb37-f400dbd0a93c.htm)

[HslCommunication.Profinet.YASKAWA.Helper 命名空间](310f2e91-5abf-4755-25f7-9225e316f6fe.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Write 方法 (IMemobus, String, Int16, Func(String, Int16, OperateResult))

[原文連結](http://api.hslcommunication.cn/html/ac26cdc6-3666-0f1e-0d3c-245e81c57aff.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.YASKAWA.Helper](../html/310f2e91-5abf-4755-25f7-9225e316f6fe.htm "HslCommunication.Profinet.YASKAWA.Helper")

[MemobusHelper 类](../html/254c50cc-c87c-57ad-eb8f-bda323fe9527.htm "MemobusHelper 类")

[MemobusHelper 方法](../html/e43bdd34-e6b5-c55a-deed-0d20ed628891.htm "MemobusHelper 方法")

[Write 方法](../html/fb613f76-96a5-5566-eb37-f400dbd0a93c.htm "Write 方法 ")

[Write 方法 (IMemobus, String, Boolean)](../html/1f35cb0d-f4ae-b48b-6830-6f58f66e94d5.htm "Write 方法 (IMemobus, String, Boolean)")

[Write 方法 (IMemobus, String, Boolean[])](../html/7cf177a0-b8e8-78da-fb1f-dcd8b1bec795.htm "Write 方法 (IMemobus, String, Boolean[])")

[Write 方法 (IMemobus, String, Byte[])](../html/a304cab6-f89e-e9e8-ffb6-8906b7a15382.htm "Write 方法 (IMemobus, String, Byte[])")

[Write 方法 (IMemobus, String, Int16, Func(String, Int16, OperateResult))](../html/ac26cdc6-3666-0f1e-0d3c-245e81c57aff.htm "Write 方法 (IMemobus, String, Int16, Func(String, Int16, OperateResult))")

[Write 方法 (IMemobus, String, UInt16, Func(String, UInt16, OperateResult))](../html/b37601c7-5f78-2779-1479-d6e47642f8d1.htm "Write 方法 (IMemobus, String, UInt16, Func(String, UInt16, OperateResult))")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MemobusHelperWrite 方法 (IMemobus, String, Int16, FuncString, Int16, OperateResult) |

写入short数据，返回是否成功  
Write short data, returns whether success

**命名空间：**
 [HslCommunication.Profinet.YASKAWA.Helper](310f2e91-5abf-4755-25f7-9225e316f6fe.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult Write(
	IMemobus memobus,
	string address,
	short value,
	Func<string, short, OperateResult> writeShort
)
```

```
Public Shared Function Write ( 
	memobus As IMemobus,
	address As String,
	value As Short,
	writeShort As Func(Of String, Short, OperateResult)
) As OperateResult
```

```
public:
static OperateResult^ Write(
	IMemobus^ memobus, 
	String^ address, 
	short value, 
	Func<String^, short, OperateResult^>^ writeShort
)
```

```
static member Write : 
        memobus : IMemobus * 
        address : string * 
        value : int16 * 
        writeShort : Func<string, int16, OperateResult> -> OperateResult 
```

#### 参数

memobus
:   类型：[HslCommunication.Profinet.YASKAWA.HelperIMemobus](0f5ca9bf-1080-d913-ed93-7c9a46291f73.htm)  

    [缺少 "M:HslCommunication.Profinet.YASKAWA.Helper.MemobusHelper.Write(HslCommunication.Profinet.YASKAWA.Helper.IMemobus,System.String,System.Int16,System.Func{System.String,System.Int16,HslCommunication.OperateResult})" 的 <param name="memobus"/> 文档]

address
:   类型：SystemString  
    起始地址

value
:   类型：SystemInt16  
    写入值

writeShort
:   类型：SystemFuncString, Int16, [OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  

    [缺少 "M:HslCommunication.Profinet.YASKAWA.Helper.MemobusHelper.Write(HslCommunication.Profinet.YASKAWA.Helper.IMemobus,System.String,System.Int16,System.Func{System.String,System.Int16,HslCommunication.OperateResult})" 的 <param name="writeShort"/> 文档]

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
带有成功标识的结果类对象

![](../icons/SectionExpanded.png)备注

单一保持寄存器的值变更，使用的主功能码为0x20, 默认子功能码为0x06，也可以写入扩展的保持型寄存器，子功能码为0x0B  
The value of a single hold register is changed, using a primary function code of 0x20 and a default subfunction code of 0x06, or an extended holding register with a subfunction code of 0x0B

![](../icons/SectionExpanded.png)示例

以下为三菱的连接对象示例，其他的设备读写情况参照下面的代码：

Int16类型示例

[复制](# "复制")

```
MelsecMcNet melsec_net = new MelsecMcNet( "192.168.0.100", 6000 );

// 简单的写入
melsec_net.Write( "D100", (short)123 );

// 如果想要判断是否写入成功
OperateResult write = melsec_net.Write( "D100", (short)123 );
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

[MemobusHelper 类](254c50cc-c87c-57ad-eb8f-bda323fe9527.htm)

[Write 重载](fb613f76-96a5-5566-eb37-f400dbd0a93c.htm)

[HslCommunication.Profinet.YASKAWA.Helper 命名空间](310f2e91-5abf-4755-25f7-9225e316f6fe.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Write 方法 (IMemobus, String, UInt16, Func(String, UInt16, OperateResult))

[原文連結](http://api.hslcommunication.cn/html/b37601c7-5f78-2779-1479-d6e47642f8d1.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.YASKAWA.Helper](../html/310f2e91-5abf-4755-25f7-9225e316f6fe.htm "HslCommunication.Profinet.YASKAWA.Helper")

[MemobusHelper 类](../html/254c50cc-c87c-57ad-eb8f-bda323fe9527.htm "MemobusHelper 类")

[MemobusHelper 方法](../html/e43bdd34-e6b5-c55a-deed-0d20ed628891.htm "MemobusHelper 方法")

[Write 方法](../html/fb613f76-96a5-5566-eb37-f400dbd0a93c.htm "Write 方法 ")

[Write 方法 (IMemobus, String, Boolean)](../html/1f35cb0d-f4ae-b48b-6830-6f58f66e94d5.htm "Write 方法 (IMemobus, String, Boolean)")

[Write 方法 (IMemobus, String, Boolean[])](../html/7cf177a0-b8e8-78da-fb1f-dcd8b1bec795.htm "Write 方法 (IMemobus, String, Boolean[])")

[Write 方法 (IMemobus, String, Byte[])](../html/a304cab6-f89e-e9e8-ffb6-8906b7a15382.htm "Write 方法 (IMemobus, String, Byte[])")

[Write 方法 (IMemobus, String, Int16, Func(String, Int16, OperateResult))](../html/ac26cdc6-3666-0f1e-0d3c-245e81c57aff.htm "Write 方法 (IMemobus, String, Int16, Func(String, Int16, OperateResult))")

[Write 方法 (IMemobus, String, UInt16, Func(String, UInt16, OperateResult))](../html/b37601c7-5f78-2779-1479-d6e47642f8d1.htm "Write 方法 (IMemobus, String, UInt16, Func(String, UInt16, OperateResult))")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MemobusHelperWrite 方法 (IMemobus, String, UInt16, FuncString, UInt16, OperateResult) |

写入ushort数据，返回是否成功  
Write ushort data, return whether the write was successful

**命名空间：**
 [HslCommunication.Profinet.YASKAWA.Helper](310f2e91-5abf-4755-25f7-9225e316f6fe.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult Write(
	IMemobus memobus,
	string address,
	ushort value,
	Func<string, ushort, OperateResult> writeUShort
)
```

```
Public Shared Function Write ( 
	memobus As IMemobus,
	address As String,
	value As UShort,
	writeUShort As Func(Of String, UShort, OperateResult)
) As OperateResult
```

```
public:
static OperateResult^ Write(
	IMemobus^ memobus, 
	String^ address, 
	unsigned short value, 
	Func<String^, unsigned short, OperateResult^>^ writeUShort
)
```

```
static member Write : 
        memobus : IMemobus * 
        address : string * 
        value : uint16 * 
        writeUShort : Func<string, uint16, OperateResult> -> OperateResult 
```

#### 参数

memobus
:   类型：[HslCommunication.Profinet.YASKAWA.HelperIMemobus](0f5ca9bf-1080-d913-ed93-7c9a46291f73.htm)  

    [缺少 "M:HslCommunication.Profinet.YASKAWA.Helper.MemobusHelper.Write(HslCommunication.Profinet.YASKAWA.Helper.IMemobus,System.String,System.UInt16,System.Func{System.String,System.UInt16,HslCommunication.OperateResult})" 的 <param name="memobus"/> 文档]

address
:   类型：SystemString  
    起始地址

value
:   类型：SystemUInt16  
    写入值

writeUShort
:   类型：SystemFuncString, UInt16, [OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  

    [缺少 "M:HslCommunication.Profinet.YASKAWA.Helper.MemobusHelper.Write(HslCommunication.Profinet.YASKAWA.Helper.IMemobus,System.String,System.UInt16,System.Func{System.String,System.UInt16,HslCommunication.OperateResult})" 的 <param name="writeUShort"/> 文档]

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
带有成功标识的结果类对象

![](../icons/SectionExpanded.png)备注

单一保持寄存器的值变更，使用的主功能码为0x20, 默认子功能码为0x06  
The value of a single hold register changes, using a primary function code of 0x20 and a default subfunction code of 0x06

![](../icons/SectionExpanded.png)示例

以下为三菱的连接对象示例，其他的设备读写情况参照下面的代码：

UInt16类型示例

[复制](# "复制")

```
MelsecMcNet melsec_net = new MelsecMcNet( "192.168.0.100", 6000 );

// 简单的写入
melsec_net.Write( "D100", (ushort)123 );

// 如果想要判断是否写入成功
OperateResult write = melsec_net.Write( "D100", (ushort)123 );
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

[MemobusHelper 类](254c50cc-c87c-57ad-eb8f-bda323fe9527.htm)

[Write 重载](fb613f76-96a5-5566-eb37-f400dbd0a93c.htm)

[HslCommunication.Profinet.YASKAWA.Helper 命名空间](310f2e91-5abf-4755-25f7-9225e316f6fe.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## WriteAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/bbc31e49-e3cd-881d-75b1-42ec61c72b34.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.YASKAWA.Helper](../html/310f2e91-5abf-4755-25f7-9225e316f6fe.htm "HslCommunication.Profinet.YASKAWA.Helper")

[MemobusHelper 类](../html/254c50cc-c87c-57ad-eb8f-bda323fe9527.htm "MemobusHelper 类")

[MemobusHelper 方法](../html/e43bdd34-e6b5-c55a-deed-0d20ed628891.htm "MemobusHelper 方法")

[WriteAsync 方法](../html/bbc31e49-e3cd-881d-75b1-42ec61c72b34.htm "WriteAsync 方法 ")

[WriteAsync 方法 (IMemobus, String, Boolean)](../html/1428b3bb-6eb9-ef9d-3af5-d0965535426f.htm "WriteAsync 方法 (IMemobus, String, Boolean)")

[WriteAsync 方法 (IMemobus, String, Boolean[])](../html/1724f1cb-ddbe-e3c2-c1f0-acbd5b105ddd.htm "WriteAsync 方法 (IMemobus, String, Boolean[])")

[WriteAsync 方法 (IMemobus, String, Byte[])](../html/0391768c-0e3f-e53a-d516-566224a49be8.htm "WriteAsync 方法 (IMemobus, String, Byte[])")

[WriteAsync 方法 (IMemobus, String, Int16, Func(String, Int16, Task(OperateResult)))](../html/fcb2e4f3-8e61-75b5-c416-2be479b4e69d.htm "WriteAsync 方法 (IMemobus, String, Int16, Func(String, Int16, Task(OperateResult)))")

[WriteAsync 方法 (IMemobus, String, UInt16, Func(String, UInt16, Task(OperateResult)))](../html/df1c6846-765e-58ee-cf5e-cc360a65f483.htm "WriteAsync 方法 (IMemobus, String, UInt16, Func(String, UInt16, Task(OperateResult)))")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MemobusHelperWriteAsync 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法静态成员 | [WriteAsync(IMemobus, String, Boolean)](1428b3bb-6eb9-ef9d-3af5-d0965535426f.htm) | 写入单个的Boolean数据，返回是否成功  Write a single Boolean data, and return whether the write was successful |
| 公共方法静态成员 | [WriteAsync(IMemobus, String, Boolean)](1724f1cb-ddbe-e3c2-c1f0-acbd5b105ddd.htm) | 批量写入Boolean数组数据，返回是否成功  Batch write Boolean array data, return whether the write was successful |
| 公共方法静态成员 | [WriteAsync(IMemobus, String, Byte)](0391768c-0e3f-e53a-d516-566224a49be8.htm) | 写入原始的byte数组数据到指定的地址，返回是否写入成功  Write the original byte array data to the specified address, and return whether the write was successful |
| 公共方法静态成员代码示例 | [WriteAsync(IMemobus, String, Int16, FuncString, Int16, TaskOperateResult)](fcb2e4f3-8e61-75b5-c416-2be479b4e69d.htm) | 写入short数据，返回是否成功  Write short data, returns whether success |
| 公共方法静态成员代码示例 | [WriteAsync(IMemobus, String, UInt16, FuncString, UInt16, TaskOperateResult)](df1c6846-765e-58ee-cf5e-cc360a65f483.htm) | 写入ushort数据，返回是否成功  Write ushort data, return whether the write was successful |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[MemobusHelper 类](254c50cc-c87c-57ad-eb8f-bda323fe9527.htm)

[HslCommunication.Profinet.YASKAWA.Helper 命名空间](310f2e91-5abf-4755-25f7-9225e316f6fe.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## WriteAsync 方法 (IMemobus, String, Boolean)

[原文連結](http://api.hslcommunication.cn/html/1428b3bb-6eb9-ef9d-3af5-d0965535426f.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.YASKAWA.Helper](../html/310f2e91-5abf-4755-25f7-9225e316f6fe.htm "HslCommunication.Profinet.YASKAWA.Helper")

[MemobusHelper 类](../html/254c50cc-c87c-57ad-eb8f-bda323fe9527.htm "MemobusHelper 类")

[MemobusHelper 方法](../html/e43bdd34-e6b5-c55a-deed-0d20ed628891.htm "MemobusHelper 方法")

[WriteAsync 方法](../html/bbc31e49-e3cd-881d-75b1-42ec61c72b34.htm "WriteAsync 方法 ")

[WriteAsync 方法 (IMemobus, String, Boolean)](../html/1428b3bb-6eb9-ef9d-3af5-d0965535426f.htm "WriteAsync 方法 (IMemobus, String, Boolean)")

[WriteAsync 方法 (IMemobus, String, Boolean[])](../html/1724f1cb-ddbe-e3c2-c1f0-acbd5b105ddd.htm "WriteAsync 方法 (IMemobus, String, Boolean[])")

[WriteAsync 方法 (IMemobus, String, Byte[])](../html/0391768c-0e3f-e53a-d516-566224a49be8.htm "WriteAsync 方法 (IMemobus, String, Byte[])")

[WriteAsync 方法 (IMemobus, String, Int16, Func(String, Int16, Task(OperateResult)))](../html/fcb2e4f3-8e61-75b5-c416-2be479b4e69d.htm "WriteAsync 方法 (IMemobus, String, Int16, Func(String, Int16, Task(OperateResult)))")

[WriteAsync 方法 (IMemobus, String, UInt16, Func(String, UInt16, Task(OperateResult)))](../html/df1c6846-765e-58ee-cf5e-cc360a65f483.htm "WriteAsync 方法 (IMemobus, String, UInt16, Func(String, UInt16, Task(OperateResult)))")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MemobusHelperWriteAsync 方法 (IMemobus, String, Boolean) |

写入单个的Boolean数据，返回是否成功  
Write a single Boolean data, and return whether the write was successful

**命名空间：**
 [HslCommunication.Profinet.YASKAWA.Helper](310f2e91-5abf-4755-25f7-9225e316f6fe.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static Task<OperateResult> WriteAsync(
	IMemobus memobus,
	string address,
	bool value
)
```

```
Public Shared Function WriteAsync ( 
	memobus As IMemobus,
	address As String,
	value As Boolean
) As Task(Of OperateResult)
```

```
public:
static Task<OperateResult^>^ WriteAsync(
	IMemobus^ memobus, 
	String^ address, 
	bool value
)
```

```
static member WriteAsync : 
        memobus : IMemobus * 
        address : string * 
        value : bool -> Task<OperateResult> 
```

#### 参数

memobus
:   类型：[HslCommunication.Profinet.YASKAWA.HelperIMemobus](0f5ca9bf-1080-d913-ed93-7c9a46291f73.htm)  

    [缺少 "M:HslCommunication.Profinet.YASKAWA.Helper.MemobusHelper.WriteAsync(HslCommunication.Profinet.YASKAWA.Helper.IMemobus,System.String,System.Boolean)" 的 <param name="memobus"/> 文档]

address
:   类型：SystemString  
    起始地址

value
:   类型：SystemBoolean  
    写入值

#### 返回值

类型：Task[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否写入成功

![](../icons/SectionExpanded.png)备注

单一线圈的状态变更，使用的主功能码为0x20, 子功能码为0x05  
The status of a single coil is changed, the main function code used is 0x20, and the sub function code is 0x05

![](../icons/SectionExpanded.png)参见

#### 引用

[MemobusHelper 类](254c50cc-c87c-57ad-eb8f-bda323fe9527.htm)

[WriteAsync 重载](bbc31e49-e3cd-881d-75b1-42ec61c72b34.htm)

[HslCommunication.Profinet.YASKAWA.Helper 命名空间](310f2e91-5abf-4755-25f7-9225e316f6fe.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## WriteAsync 方法 (IMemobus, String, Boolean[])

[原文連結](http://api.hslcommunication.cn/html/1724f1cb-ddbe-e3c2-c1f0-acbd5b105ddd.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.YASKAWA.Helper](../html/310f2e91-5abf-4755-25f7-9225e316f6fe.htm "HslCommunication.Profinet.YASKAWA.Helper")

[MemobusHelper 类](../html/254c50cc-c87c-57ad-eb8f-bda323fe9527.htm "MemobusHelper 类")

[MemobusHelper 方法](../html/e43bdd34-e6b5-c55a-deed-0d20ed628891.htm "MemobusHelper 方法")

[WriteAsync 方法](../html/bbc31e49-e3cd-881d-75b1-42ec61c72b34.htm "WriteAsync 方法 ")

[WriteAsync 方法 (IMemobus, String, Boolean)](../html/1428b3bb-6eb9-ef9d-3af5-d0965535426f.htm "WriteAsync 方法 (IMemobus, String, Boolean)")

[WriteAsync 方法 (IMemobus, String, Boolean[])](../html/1724f1cb-ddbe-e3c2-c1f0-acbd5b105ddd.htm "WriteAsync 方法 (IMemobus, String, Boolean[])")

[WriteAsync 方法 (IMemobus, String, Byte[])](../html/0391768c-0e3f-e53a-d516-566224a49be8.htm "WriteAsync 方法 (IMemobus, String, Byte[])")

[WriteAsync 方法 (IMemobus, String, Int16, Func(String, Int16, Task(OperateResult)))](../html/fcb2e4f3-8e61-75b5-c416-2be479b4e69d.htm "WriteAsync 方法 (IMemobus, String, Int16, Func(String, Int16, Task(OperateResult)))")

[WriteAsync 方法 (IMemobus, String, UInt16, Func(String, UInt16, Task(OperateResult)))](../html/df1c6846-765e-58ee-cf5e-cc360a65f483.htm "WriteAsync 方法 (IMemobus, String, UInt16, Func(String, UInt16, Task(OperateResult)))")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MemobusHelperWriteAsync 方法 (IMemobus, String, Boolean) |

批量写入Boolean数组数据，返回是否成功  
Batch write Boolean array data, return whether the write was successful

**命名空间：**
 [HslCommunication.Profinet.YASKAWA.Helper](310f2e91-5abf-4755-25f7-9225e316f6fe.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static Task<OperateResult> WriteAsync(
	IMemobus memobus,
	string address,
	bool[] value
)
```

```
Public Shared Function WriteAsync ( 
	memobus As IMemobus,
	address As String,
	value As Boolean()
) As Task(Of OperateResult)
```

```
public:
static Task<OperateResult^>^ WriteAsync(
	IMemobus^ memobus, 
	String^ address, 
	array<bool>^ value
)
```

```
static member WriteAsync : 
        memobus : IMemobus * 
        address : string * 
        value : bool[] -> Task<OperateResult> 
```

#### 参数

memobus
:   类型：[HslCommunication.Profinet.YASKAWA.HelperIMemobus](0f5ca9bf-1080-d913-ed93-7c9a46291f73.htm)  

    [缺少 "M:HslCommunication.Profinet.YASKAWA.Helper.MemobusHelper.WriteAsync(HslCommunication.Profinet.YASKAWA.Helper.IMemobus,System.String,System.Boolean[])" 的 <param name="memobus"/> 文档]

address
:   类型：SystemString  
    起始地址

value
:   类型：SystemBoolean  
    写入值

#### 返回值

类型：Task[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否写入成功

![](../icons/SectionExpanded.png)备注

多个线圈的状态更改，默认使用的是 0x0f 子功能码。  
The status of multiple coils is changed, and the sub-function code 0x0f is used by default.

![](../icons/SectionExpanded.png)参见

#### 引用

[MemobusHelper 类](254c50cc-c87c-57ad-eb8f-bda323fe9527.htm)

[WriteAsync 重载](bbc31e49-e3cd-881d-75b1-42ec61c72b34.htm)

[HslCommunication.Profinet.YASKAWA.Helper 命名空间](310f2e91-5abf-4755-25f7-9225e316f6fe.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## WriteAsync 方法 (IMemobus, String, Byte[])

[原文連結](http://api.hslcommunication.cn/html/0391768c-0e3f-e53a-d516-566224a49be8.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.YASKAWA.Helper](../html/310f2e91-5abf-4755-25f7-9225e316f6fe.htm "HslCommunication.Profinet.YASKAWA.Helper")

[MemobusHelper 类](../html/254c50cc-c87c-57ad-eb8f-bda323fe9527.htm "MemobusHelper 类")

[MemobusHelper 方法](../html/e43bdd34-e6b5-c55a-deed-0d20ed628891.htm "MemobusHelper 方法")

[WriteAsync 方法](../html/bbc31e49-e3cd-881d-75b1-42ec61c72b34.htm "WriteAsync 方法 ")

[WriteAsync 方法 (IMemobus, String, Boolean)](../html/1428b3bb-6eb9-ef9d-3af5-d0965535426f.htm "WriteAsync 方法 (IMemobus, String, Boolean)")

[WriteAsync 方法 (IMemobus, String, Boolean[])](../html/1724f1cb-ddbe-e3c2-c1f0-acbd5b105ddd.htm "WriteAsync 方法 (IMemobus, String, Boolean[])")

[WriteAsync 方法 (IMemobus, String, Byte[])](../html/0391768c-0e3f-e53a-d516-566224a49be8.htm "WriteAsync 方法 (IMemobus, String, Byte[])")

[WriteAsync 方法 (IMemobus, String, Int16, Func(String, Int16, Task(OperateResult)))](../html/fcb2e4f3-8e61-75b5-c416-2be479b4e69d.htm "WriteAsync 方法 (IMemobus, String, Int16, Func(String, Int16, Task(OperateResult)))")

[WriteAsync 方法 (IMemobus, String, UInt16, Func(String, UInt16, Task(OperateResult)))](../html/df1c6846-765e-58ee-cf5e-cc360a65f483.htm "WriteAsync 方法 (IMemobus, String, UInt16, Func(String, UInt16, Task(OperateResult)))")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MemobusHelperWriteAsync 方法 (IMemobus, String, Byte) |

写入原始的byte数组数据到指定的地址，返回是否写入成功  
Write the original byte array data to the specified address, and return whether the write was successful

**命名空间：**
 [HslCommunication.Profinet.YASKAWA.Helper](310f2e91-5abf-4755-25f7-9225e316f6fe.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static Task<OperateResult> WriteAsync(
	IMemobus memobus,
	string address,
	byte[] value
)
```

```
Public Shared Function WriteAsync ( 
	memobus As IMemobus,
	address As String,
	value As Byte()
) As Task(Of OperateResult)
```

```
public:
static Task<OperateResult^>^ WriteAsync(
	IMemobus^ memobus, 
	String^ address, 
	array<unsigned char>^ value
)
```

```
static member WriteAsync : 
        memobus : IMemobus * 
        address : string * 
        value : byte[] -> Task<OperateResult> 
```

#### 参数

memobus
:   类型：[HslCommunication.Profinet.YASKAWA.HelperIMemobus](0f5ca9bf-1080-d913-ed93-7c9a46291f73.htm)  

    [缺少 "M:HslCommunication.Profinet.YASKAWA.Helper.MemobusHelper.WriteAsync(HslCommunication.Profinet.YASKAWA.Helper.IMemobus,System.String,System.Byte[])" 的 <param name="memobus"/> 文档]

address
:   类型：SystemString  
    起始地址

value
:   类型：SystemByte  
    写入值

#### 返回值

类型：Task[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否写入成功

![](../icons/SectionExpanded.png)备注

连续的寄存器写入操作，默认功能码是0x10，如果需要写入扩展的寄存器，使用 x=0xA;100 或是 x=10;100 即可。  
For continuous register write operation, the default function code is 0x10. If you need to write an extended register, use x=0xA;100 or x=10;100.

![](../icons/SectionExpanded.png)参见

#### 引用

[MemobusHelper 类](254c50cc-c87c-57ad-eb8f-bda323fe9527.htm)

[WriteAsync 重载](bbc31e49-e3cd-881d-75b1-42ec61c72b34.htm)

[HslCommunication.Profinet.YASKAWA.Helper 命名空间](310f2e91-5abf-4755-25f7-9225e316f6fe.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## WriteAsync 方法 (IMemobus, String, Int16, Func(String, Int16, Task(OperateResult)))

[原文連結](http://api.hslcommunication.cn/html/fcb2e4f3-8e61-75b5-c416-2be479b4e69d.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.YASKAWA.Helper](../html/310f2e91-5abf-4755-25f7-9225e316f6fe.htm "HslCommunication.Profinet.YASKAWA.Helper")

[MemobusHelper 类](../html/254c50cc-c87c-57ad-eb8f-bda323fe9527.htm "MemobusHelper 类")

[MemobusHelper 方法](../html/e43bdd34-e6b5-c55a-deed-0d20ed628891.htm "MemobusHelper 方法")

[WriteAsync 方法](../html/bbc31e49-e3cd-881d-75b1-42ec61c72b34.htm "WriteAsync 方法 ")

[WriteAsync 方法 (IMemobus, String, Boolean)](../html/1428b3bb-6eb9-ef9d-3af5-d0965535426f.htm "WriteAsync 方法 (IMemobus, String, Boolean)")

[WriteAsync 方法 (IMemobus, String, Boolean[])](../html/1724f1cb-ddbe-e3c2-c1f0-acbd5b105ddd.htm "WriteAsync 方法 (IMemobus, String, Boolean[])")

[WriteAsync 方法 (IMemobus, String, Byte[])](../html/0391768c-0e3f-e53a-d516-566224a49be8.htm "WriteAsync 方法 (IMemobus, String, Byte[])")

[WriteAsync 方法 (IMemobus, String, Int16, Func(String, Int16, Task(OperateResult)))](../html/fcb2e4f3-8e61-75b5-c416-2be479b4e69d.htm "WriteAsync 方法 (IMemobus, String, Int16, Func(String, Int16, Task(OperateResult)))")

[WriteAsync 方法 (IMemobus, String, UInt16, Func(String, UInt16, Task(OperateResult)))](../html/df1c6846-765e-58ee-cf5e-cc360a65f483.htm "WriteAsync 方法 (IMemobus, String, UInt16, Func(String, UInt16, Task(OperateResult)))")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| MemobusHelperWriteAsync 方法 (IMemobus, String, Int16, FuncString, Int16, TaskOperateResult) |

写入short数据，返回是否成功  
Write short data, returns whether success

**命名空间：**
 [HslCommunication.Profinet.YASKAWA.Helper](310f2e91-5abf-4755-25f7-9225e316f6fe.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static Task<OperateResult> WriteAsync(
	IMemobus memobus,
	string address,
	short value,
	Func<string, short, Task<OperateResult>> writeShort
)
```

```
Public Shared Function WriteAsync ( 
	memobus As IMemobus,
	address As String,
	value As Short,
	writeShort As Func(Of String, Short, Task(Of OperateResult))
) As Task(Of OperateResult)
```

```
public:
static Task<OperateResult^>^ WriteAsync(
	IMemobus^ memobus, 
	String^ address, 
	short value, 
	Func<String^, short, Task<OperateResult^>^>^ writeShort
)
```

```
static member WriteAsync : 
        memobus : IMemobus * 
        address : string * 
        value : int16 * 
        writeShort : Func<string, int16, Task<OperateResult>> -> Task<OperateResult> 
```

#### 参数

memobus
:   类型：[HslCommunication.Profinet.YASKAWA.HelperIMemobus](0f5ca9bf-1080-d913-ed93-7c9a46291f73.htm)  

    [缺少 "M:HslCommunication.Profinet.YASKAWA.Helper.MemobusHelper.WriteAsync(HslCommunication.Profinet.YASKAWA.Helper.IMemobus,System.String,System.Int16,System.Func{System.String,System.Int16,System.Threading.Tasks.Task{HslCommunication.OperateResult}})" 的 <param name="memobus"/> 文档]

address
:   类型：SystemString  
    起始地址

value
:   类型：SystemInt16  
    写入值

writeShort
:   类型：SystemFuncString, Int16, Task[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  

    [缺少 "M:HslCommunication.Profinet.YASKAWA.Helper.MemobusHelper.WriteAsync(HslCommunication.Profinet.YASKAWA.Helper.IMemobus,System.String,System.Int16,System.Func{System.String,System.Int16,System.Threading.Tasks.Task{HslCommunication.OperateResult}})" 的 <param name="writeShort"/> 文档]

#### 返回值

类型：Task[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
带有成功标识的结果类对象

![](../icons/SectionExpanded.png)备注

单一保持寄存器的值变更，使用的主功能码为0x20, 默认子功能码为0x06，也可以写入扩展的保持型寄存器，子功能码为0x0B  
The value of a single hold register is changed, using a primary function code of 0x20 and a default subfunction code of 0x06, or an extended holding register with a subfunction code of 0x0B

![](../icons/SectionExpanded.png)示例

以下为三菱的连接对象示例，其他的设备读写情况参照下面的代码：

Int16类型示例

[复制](# "复制")

```
MelsecMcNet melsec_net = new MelsecMcNet( "192.168.0.100", 6000 );

// 简单的写入
melsec_net.Write( "D100", (short)123 );

// 如果想要判断是否写入成功
OperateResult write = melsec_net.Write( "D100", (short)123 );
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

[MemobusHelper 类](254c50cc-c87c-57ad-eb8f-bda323fe9527.htm)

[WriteAsync 重载](bbc31e49-e3cd-881d-75b1-42ec61c72b34.htm)

[HslCommunication.Profinet.YASKAWA.Helper 命名空间](310f2e91-5abf-4755-25f7-9225e316f6fe.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)