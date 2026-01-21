# HslCommunication - HslCommunication.Profinet.Melsec.Helper

> 分類頁數: 30



---
## HslCommunication.Profinet.Melsec.Helper

[原文連結](http://api.hslcommunication.cn/html/2922ab51-7617-dc63-33ec-7505340acd85.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Melsec.Helper](../html/2922ab51-7617-dc63-33ec-7505340acd85.htm "HslCommunication.Profinet.Melsec.Helper")

[IMelsecFxSerial 接口](../html/ad624299-4861-4650-998f-7ae37d5a2c0f.htm "IMelsecFxSerial 接口")

[IReadWriteA3C 接口](../html/08311113-f4d5-8321-79ab-3aa3b26d9a6f.htm "IReadWriteA3C 接口")

[IReadWriteFxLinks 接口](../html/44bfc17d-ccc2-b900-f4ec-1713ef97627f.htm "IReadWriteFxLinks 接口")

[IReadWriteMc 接口](../html/6ed75f4f-6587-8991-17e0-3899347e38d8.htm "IReadWriteMc 接口")

[McAsciiHelper 类](../html/34069267-3020-34f6-7e24-08cb0df10659.htm "McAsciiHelper 类")

[McBinaryHelper 类](../html/998e07f1-0394-d262-40b5-f0c4c56b364f.htm "McBinaryHelper 类")

[McHelper 类](../html/b99ac294-1134-5418-5be8-d0cf7aec7b8e.htm "McHelper 类")

[McType 枚举](../html/26e256cf-18a7-bca3-2fe9-544737862be1.htm "McType 枚举")

[MelsecA3CNetHelper 类](../html/585f912f-7ec6-abcf-48e4-2aad8e8efd03.htm "MelsecA3CNetHelper 类")

[MelsecFxLinksHelper 类](../html/c745e738-f2bc-5715-0bc3-bad096a560bb.htm "MelsecFxLinksHelper 类")

[MelsecFxSerialHelper 类](../html/b262f62b-de0b-b459-bb6e-b193170c9cf8.htm "MelsecFxSerialHelper 类")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslCommunication.Profinet.Melsec.Helper 命名空间 |

[缺少 "N:HslCommunication.Profinet.Melsec.Helper" 的 <summary> 文档]

![](../icons/SectionExpanded.png)类

|  | 类 | 说明 |
| --- | --- | --- |
| 公共类 | [McAsciiHelper](34069267-3020-34f6-7e24-08cb0df10659.htm) | 基于MC协议的ASCII格式的辅助类 |
| 公共类 | [McBinaryHelper](998e07f1-0394-d262-40b5-f0c4c56b364f.htm) | 三菱PLC，二进制的辅助类对象 |
| 公共类 | [McHelper](b99ac294-1134-5418-5be8-d0cf7aec7b8e.htm) | MC协议的辅助类对象，提供了MC协议的读写操作的基本支持 |
| 公共类 | [MelsecA3CNetHelper](585f912f-7ec6-abcf-48e4-2aad8e8efd03.htm) | MelsecA3CNet1协议通信的辅助类 |
| 公共类 | [MelsecFxLinksHelper](c745e738-f2bc-5715-0bc3-bad096a560bb.htm) | 三菱的FxLinks的辅助方法信息 |
| 公共类 | [MelsecFxSerialHelper](b262f62b-de0b-b459-bb6e-b193170c9cf8.htm) | 三菱编程口协议的辅助方法，定义了如何读写bool数据，以及读写原始字节的数据。  The auxiliary method of Mitsubishi programming port protocol defines how to read and write bool data and read and write raw byte data. |

![](../icons/SectionExpanded.png)接口

|  | 接口 | 说明 |
| --- | --- | --- |
| 公共接口 | [IMelsecFxSerial](ad624299-4861-4650-998f-7ae37d5a2c0f.htm) | 三菱的串口的接口类对象 |
| 公共接口 | [IReadWriteA3C](08311113-f4d5-8321-79ab-3aa3b26d9a6f.htm) | 三菱的A3C协议类接口对象，具有站号，是否和校验的属性  Mitsubishi's A3C protocol interface object, which has the attributes of station number, and checksum |
| 公共接口 | [IReadWriteFxLinks](44bfc17d-ccc2-b900-f4ec-1713ef97627f.htm) | 三菱的FxLink协议接口的设备信息 |
| 公共接口 | [IReadWriteMc](6ed75f4f-6587-8991-17e0-3899347e38d8.htm) | 基于MC协议的标准的设备接口，适用任何基于MC协议的PLC设备，主要是三菱，基恩士，松下的PLC设备。  The standard equipment interface based on MC protocol is suitable for any PLC equipment based on MC protocol, mainly PLC equipment from Mitsubishi, Keyence, and Panasonic. |

![](../icons/SectionExpanded.png)枚举

|  | 枚举 | 说明 |
| --- | --- | --- |
| 公共枚举 | [McType](26e256cf-18a7-bca3-2fe9-544737862be1.htm) | MC协议的类型 |

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## IMelsecFxSerial 接口

[原文連結](http://api.hslcommunication.cn/html/ad624299-4861-4650-998f-7ae37d5a2c0f.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Melsec.Helper](../html/2922ab51-7617-dc63-33ec-7505340acd85.htm "HslCommunication.Profinet.Melsec.Helper")

[IMelsecFxSerial 接口](../html/ad624299-4861-4650-998f-7ae37d5a2c0f.htm "IMelsecFxSerial 接口")

[IMelsecFxSerial 属性](../html/abbd947b-9481-9bb9-9e27-d46b3870708f.htm "IMelsecFxSerial 属性")

[IMelsecFxSerial 方法](../html/9e9b202f-00ed-53f2-299b-3a4c12e37c35.htm "IMelsecFxSerial 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IMelsecFxSerial 接口 |

三菱的串口的接口类对象

**命名空间：**
 [HslCommunication.Profinet.Melsec.Helper](2922ab51-7617-dc63-33ec-7505340acd85.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public interface IMelsecFxSerial : IReadWriteNet
```

```
Public Interface IMelsecFxSerial
	Inherits IReadWriteNet
```

```
public interface class IMelsecFxSerial : IReadWriteNet
```

```
type IMelsecFxSerial =  
    interface
        interface IReadWriteNet
    end
```

IMelsecFxSerial 类型公开以下成员。

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [ConnectionId](20787592-e3ec-e0d0-27c2-abfa24654bd2.htm) | 当前连接的唯一ID号，默认为长度20的guid码加随机数组成，方便列表管理，也可以自己指定  The unique ID number of the current connection. The default is a 20-digit guid code plus a random number. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共属性代码示例 | [LogNet](38760c1f-22e3-f8fe-7512-15849a4b6eb0.htm) | 组件的日志工具，支持日志记录，只要实例化后，当前网络的基本信息，就以[DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm)等级进行输出  The component's logging tool supports logging. As long as the instantiation of the basic network information, the output will be output at [DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm) (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [ActivePlc](806adee4-2d65-790f-391c-f7f3fec37772.htm) | 激活PLC的接收状态，需要再和PLC交互之前进行调用，之后就需要再调用了。 |
| 公共方法 | [ActivePlcAsync](c0ac8f63-0111-265f-9887-c3e6b2eee8a4.htm) | 激活PLC的接收状态，需要再和PLC交互之前进行调用，之后就需要再调用了。 |
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

[HslCommunication.Profinet.Melsec.Helper 命名空间](2922ab51-7617-dc63-33ec-7505340acd85.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## IMelsecFxSerial 属性

[原文連結](http://api.hslcommunication.cn/html/abbd947b-9481-9bb9-9e27-d46b3870708f.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Melsec.Helper](../html/2922ab51-7617-dc63-33ec-7505340acd85.htm "HslCommunication.Profinet.Melsec.Helper")

[IMelsecFxSerial 接口](../html/ad624299-4861-4650-998f-7ae37d5a2c0f.htm "IMelsecFxSerial 接口")

[IMelsecFxSerial 属性](../html/abbd947b-9481-9bb9-9e27-d46b3870708f.htm "IMelsecFxSerial 属性")

[IMelsecFxSerial 方法](../html/9e9b202f-00ed-53f2-299b-3a4c12e37c35.htm "IMelsecFxSerial 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IMelsecFxSerial 属性 |

[IMelsecFxSerial](ad624299-4861-4650-998f-7ae37d5a2c0f.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [ConnectionId](20787592-e3ec-e0d0-27c2-abfa24654bd2.htm) | 当前连接的唯一ID号，默认为长度20的guid码加随机数组成，方便列表管理，也可以自己指定  The unique ID number of the current connection. The default is a 20-digit guid code plus a random number. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共属性代码示例 | [LogNet](38760c1f-22e3-f8fe-7512-15849a4b6eb0.htm) | 组件的日志工具，支持日志记录，只要实例化后，当前网络的基本信息，就以[DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm)等级进行输出  The component's logging tool supports logging. As long as the instantiation of the basic network information, the output will be output at [DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm) (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[IMelsecFxSerial 接口](ad624299-4861-4650-998f-7ae37d5a2c0f.htm)

[HslCommunication.Profinet.Melsec.Helper 命名空间](2922ab51-7617-dc63-33ec-7505340acd85.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## IMelsecFxSerial 方法

[原文連結](http://api.hslcommunication.cn/html/9e9b202f-00ed-53f2-299b-3a4c12e37c35.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Melsec.Helper](../html/2922ab51-7617-dc63-33ec-7505340acd85.htm "HslCommunication.Profinet.Melsec.Helper")

[IMelsecFxSerial 接口](../html/ad624299-4861-4650-998f-7ae37d5a2c0f.htm "IMelsecFxSerial 接口")

[IMelsecFxSerial 方法](../html/9e9b202f-00ed-53f2-299b-3a4c12e37c35.htm "IMelsecFxSerial 方法")

[ActivePlc 方法](../html/806adee4-2d65-790f-391c-f7f3fec37772.htm "ActivePlc 方法 ")

[ActivePlcAsync 方法](../html/c0ac8f63-0111-265f-9887-c3e6b2eee8a4.htm "ActivePlcAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IMelsecFxSerial 方法 |

[IMelsecFxSerial](ad624299-4861-4650-998f-7ae37d5a2c0f.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [ActivePlc](806adee4-2d65-790f-391c-f7f3fec37772.htm) | 激活PLC的接收状态，需要再和PLC交互之前进行调用，之后就需要再调用了。 |
| 公共方法 | [ActivePlcAsync](c0ac8f63-0111-265f-9887-c3e6b2eee8a4.htm) | 激活PLC的接收状态，需要再和PLC交互之前进行调用，之后就需要再调用了。 |
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

[IMelsecFxSerial 接口](ad624299-4861-4650-998f-7ae37d5a2c0f.htm)

[HslCommunication.Profinet.Melsec.Helper 命名空间](2922ab51-7617-dc63-33ec-7505340acd85.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ActivePlc 方法 

[原文連結](http://api.hslcommunication.cn/html/806adee4-2d65-790f-391c-f7f3fec37772.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Melsec.Helper](../html/2922ab51-7617-dc63-33ec-7505340acd85.htm "HslCommunication.Profinet.Melsec.Helper")

[IMelsecFxSerial 接口](../html/ad624299-4861-4650-998f-7ae37d5a2c0f.htm "IMelsecFxSerial 接口")

[IMelsecFxSerial 方法](../html/9e9b202f-00ed-53f2-299b-3a4c12e37c35.htm "IMelsecFxSerial 方法")

[ActivePlc 方法](../html/806adee4-2d65-790f-391c-f7f3fec37772.htm "ActivePlc 方法 ")

[ActivePlcAsync 方法](../html/c0ac8f63-0111-265f-9887-c3e6b2eee8a4.htm "ActivePlcAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IMelsecFxSerialActivePlc 方法 |

激活PLC的接收状态，需要再和PLC交互之前进行调用，之后就需要再调用了。

**命名空间：**
 [HslCommunication.Profinet.Melsec.Helper](2922ab51-7617-dc63-33ec-7505340acd85.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
OperateResult ActivePlc()
```

```
Function ActivePlc As OperateResult
```

```
OperateResult^ ActivePlc()
```

```
abstract ActivePlc : unit -> OperateResult 
```

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否激活成功

![](../icons/SectionExpanded.png)参见

#### 引用

[IMelsecFxSerial 接口](ad624299-4861-4650-998f-7ae37d5a2c0f.htm)

[HslCommunication.Profinet.Melsec.Helper 命名空间](2922ab51-7617-dc63-33ec-7505340acd85.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ActivePlcAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/c0ac8f63-0111-265f-9887-c3e6b2eee8a4.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Melsec.Helper](../html/2922ab51-7617-dc63-33ec-7505340acd85.htm "HslCommunication.Profinet.Melsec.Helper")

[IMelsecFxSerial 接口](../html/ad624299-4861-4650-998f-7ae37d5a2c0f.htm "IMelsecFxSerial 接口")

[IMelsecFxSerial 方法](../html/9e9b202f-00ed-53f2-299b-3a4c12e37c35.htm "IMelsecFxSerial 方法")

[ActivePlc 方法](../html/806adee4-2d65-790f-391c-f7f3fec37772.htm "ActivePlc 方法 ")

[ActivePlcAsync 方法](../html/c0ac8f63-0111-265f-9887-c3e6b2eee8a4.htm "ActivePlcAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IMelsecFxSerialActivePlcAsync 方法 |

激活PLC的接收状态，需要再和PLC交互之前进行调用，之后就需要再调用了。

**命名空间：**
 [HslCommunication.Profinet.Melsec.Helper](2922ab51-7617-dc63-33ec-7505340acd85.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
Task<OperateResult> ActivePlcAsync()
```

```
Function ActivePlcAsync As Task(Of OperateResult)
```

```
Task<OperateResult^>^ ActivePlcAsync()
```

```
abstract ActivePlcAsync : unit -> Task<OperateResult> 
```

#### 返回值

类型：Task[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否激活成功

![](../icons/SectionExpanded.png)参见

#### 引用

[IMelsecFxSerial 接口](ad624299-4861-4650-998f-7ae37d5a2c0f.htm)

[HslCommunication.Profinet.Melsec.Helper 命名空间](2922ab51-7617-dc63-33ec-7505340acd85.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## IReadWriteA3C 接口

[原文連結](http://api.hslcommunication.cn/html/08311113-f4d5-8321-79ab-3aa3b26d9a6f.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Melsec.Helper](../html/2922ab51-7617-dc63-33ec-7505340acd85.htm "HslCommunication.Profinet.Melsec.Helper")

[IReadWriteA3C 接口](../html/08311113-f4d5-8321-79ab-3aa3b26d9a6f.htm "IReadWriteA3C 接口")

[IReadWriteA3C 属性](../html/a2ec1b91-2087-dd88-e1ab-d9b9023b7790.htm "IReadWriteA3C 属性")

[IReadWriteA3C 方法](../html/e72bcc42-1c30-a900-a811-6adbdf3fe2ca.htm "IReadWriteA3C 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IReadWriteA3C 接口 |

三菱的A3C协议类接口对象，具有站号，是否和校验的属性  
Mitsubishi's A3C protocol interface object, which has the attributes of station number, and checksum

**命名空间：**
 [HslCommunication.Profinet.Melsec.Helper](2922ab51-7617-dc63-33ec-7505340acd85.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public interface IReadWriteA3C : IReadWriteDevice, 
	IReadWriteNet
```

```
Public Interface IReadWriteA3C
	Inherits IReadWriteDevice, IReadWriteNet
```

```
public interface class IReadWriteA3C : IReadWriteDevice, 
	IReadWriteNet
```

```
type IReadWriteA3C =  
    interface
        interface IReadWriteDevice
        interface IReadWriteNet
    end
```

IReadWriteA3C 类型公开以下成员。

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性代码示例 | [ByteTransform](0836ec33-d464-ed39-37ff-b2ada0d286fc.htm) | 当前的数据变换机制，当你需要从字节数据转换类型数据的时候需要。  The current data transformation mechanism is required when you need to convert type data from byte data. (继承自 [IReadWriteDevice](d987afcc-f99c-6aed-3847-5b1890504635.htm)。) |
| 公共属性 | [ConnectionId](20787592-e3ec-e0d0-27c2-abfa24654bd2.htm) | 当前连接的唯一ID号，默认为长度20的guid码加随机数组成，方便列表管理，也可以自己指定  The unique ID number of the current connection. The default is a 20-digit guid code plus a random number. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共属性 | [EnableWriteBitToWordRegister](68ca52a7-8d7e-ce6d-11aa-fd8ee0bcb61b.htm) | 是否开启支持写入位到字寄存器的功能，该功能先读取字寄存器的字数据，然后修改其中的位，再写入回去，可能存在脏数据的风险  Whether to enable supporting the function of writing bit-to-word registers, which first reads the word data of the word register, then modifies the bits in it, and then writes back, there may be a risk of dirty data |
| 公共属性 | [Format](bb00688f-68ac-927d-3f53-c021515ccb63.htm) | 当前的A3C协议的格式信息，可选格式1，2，3，4，默认格式1  Format information of the current A3C protocol, optional format 1, 2, 3, 4, default format 1 |
| 公共属性代码示例 | [LogNet](38760c1f-22e3-f8fe-7512-15849a4b6eb0.htm) | 组件的日志工具，支持日志记录，只要实例化后，当前网络的基本信息，就以[DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm)等级进行输出  The component's logging tool supports logging. As long as the instantiation of the basic network information, the output will be output at [DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm) (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共属性 | [Station](e6637ca9-5a39-00bd-68de-6a6c1da07af9.htm) | 当前A3C协议的站编号信息  Station number information of the current A3C protocol |
| 公共属性 | [SumCheck](165338dc-5988-fe3e-320f-fb46ce851b1b.htm) | 当前的A3C协议是否使用和校验，默认使用  Whether the current A3C protocol uses sum check, it is used by default |

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

[HslCommunication.Profinet.Melsec.Helper 命名空间](2922ab51-7617-dc63-33ec-7505340acd85.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## IReadWriteA3C 属性

[原文連結](http://api.hslcommunication.cn/html/a2ec1b91-2087-dd88-e1ab-d9b9023b7790.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Melsec.Helper](../html/2922ab51-7617-dc63-33ec-7505340acd85.htm "HslCommunication.Profinet.Melsec.Helper")

[IReadWriteA3C 接口](../html/08311113-f4d5-8321-79ab-3aa3b26d9a6f.htm "IReadWriteA3C 接口")

[IReadWriteA3C 属性](../html/a2ec1b91-2087-dd88-e1ab-d9b9023b7790.htm "IReadWriteA3C 属性")

[EnableWriteBitToWordRegister 属性](../html/68ca52a7-8d7e-ce6d-11aa-fd8ee0bcb61b.htm "EnableWriteBitToWordRegister 属性 ")

[Format 属性](../html/bb00688f-68ac-927d-3f53-c021515ccb63.htm "Format 属性 ")

[Station 属性](../html/e6637ca9-5a39-00bd-68de-6a6c1da07af9.htm "Station 属性 ")

[SumCheck 属性](../html/165338dc-5988-fe3e-320f-fb46ce851b1b.htm "SumCheck 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IReadWriteA3C 属性 |

[IReadWriteA3C](08311113-f4d5-8321-79ab-3aa3b26d9a6f.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性代码示例 | [ByteTransform](0836ec33-d464-ed39-37ff-b2ada0d286fc.htm) | 当前的数据变换机制，当你需要从字节数据转换类型数据的时候需要。  The current data transformation mechanism is required when you need to convert type data from byte data. (继承自 [IReadWriteDevice](d987afcc-f99c-6aed-3847-5b1890504635.htm)。) |
| 公共属性 | [ConnectionId](20787592-e3ec-e0d0-27c2-abfa24654bd2.htm) | 当前连接的唯一ID号，默认为长度20的guid码加随机数组成，方便列表管理，也可以自己指定  The unique ID number of the current connection. The default is a 20-digit guid code plus a random number. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共属性 | [EnableWriteBitToWordRegister](68ca52a7-8d7e-ce6d-11aa-fd8ee0bcb61b.htm) | 是否开启支持写入位到字寄存器的功能，该功能先读取字寄存器的字数据，然后修改其中的位，再写入回去，可能存在脏数据的风险  Whether to enable supporting the function of writing bit-to-word registers, which first reads the word data of the word register, then modifies the bits in it, and then writes back, there may be a risk of dirty data |
| 公共属性 | [Format](bb00688f-68ac-927d-3f53-c021515ccb63.htm) | 当前的A3C协议的格式信息，可选格式1，2，3，4，默认格式1  Format information of the current A3C protocol, optional format 1, 2, 3, 4, default format 1 |
| 公共属性代码示例 | [LogNet](38760c1f-22e3-f8fe-7512-15849a4b6eb0.htm) | 组件的日志工具，支持日志记录，只要实例化后，当前网络的基本信息，就以[DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm)等级进行输出  The component's logging tool supports logging. As long as the instantiation of the basic network information, the output will be output at [DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm) (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共属性 | [Station](e6637ca9-5a39-00bd-68de-6a6c1da07af9.htm) | 当前A3C协议的站编号信息  Station number information of the current A3C protocol |
| 公共属性 | [SumCheck](165338dc-5988-fe3e-320f-fb46ce851b1b.htm) | 当前的A3C协议是否使用和校验，默认使用  Whether the current A3C protocol uses sum check, it is used by default |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[IReadWriteA3C 接口](08311113-f4d5-8321-79ab-3aa3b26d9a6f.htm)

[HslCommunication.Profinet.Melsec.Helper 命名空间](2922ab51-7617-dc63-33ec-7505340acd85.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## EnableWriteBitToWordRegister 属性 

[原文連結](http://api.hslcommunication.cn/html/68ca52a7-8d7e-ce6d-11aa-fd8ee0bcb61b.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Melsec.Helper](../html/2922ab51-7617-dc63-33ec-7505340acd85.htm "HslCommunication.Profinet.Melsec.Helper")

[IReadWriteA3C 接口](../html/08311113-f4d5-8321-79ab-3aa3b26d9a6f.htm "IReadWriteA3C 接口")

[IReadWriteA3C 属性](../html/a2ec1b91-2087-dd88-e1ab-d9b9023b7790.htm "IReadWriteA3C 属性")

[EnableWriteBitToWordRegister 属性](../html/68ca52a7-8d7e-ce6d-11aa-fd8ee0bcb61b.htm "EnableWriteBitToWordRegister 属性 ")

[Format 属性](../html/bb00688f-68ac-927d-3f53-c021515ccb63.htm "Format 属性 ")

[Station 属性](../html/e6637ca9-5a39-00bd-68de-6a6c1da07af9.htm "Station 属性 ")

[SumCheck 属性](../html/165338dc-5988-fe3e-320f-fb46ce851b1b.htm "SumCheck 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IReadWriteA3CEnableWriteBitToWordRegister 属性 |

是否开启支持写入位到字寄存器的功能，该功能先读取字寄存器的字数据，然后修改其中的位，再写入回去，可能存在脏数据的风险  
Whether to enable supporting the function of writing bit-to-word registers, which first reads the word data of the word register, then modifies the bits in it, and then writes back, there may be a risk of dirty data

**命名空间：**
 [HslCommunication.Profinet.Melsec.Helper](2922ab51-7617-dc63-33ec-7505340acd85.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
bool EnableWriteBitToWordRegister { get; set; }
```

```
Property EnableWriteBitToWordRegister As Boolean
	Get
	Set
```

```
property bool EnableWriteBitToWordRegister {
	bool get ();
	void set (bool value);
}
```

```
abstract EnableWriteBitToWordRegister : bool with get, set
```

#### 属性值

类型：Boolean

![](../icons/SectionExpanded.png)备注

关于脏数据风险：从读取数据，修改位，再次写入数据时，大概需要经过3ms~10ms不等的时间，如果此期间内PLC修改了该字寄存器的其他位，再次写入数据时会恢复该点位的数据到读取时的初始值，可能引发设备故障，请谨慎开启此功能。  
About dirty data risk: from reading data, modifying bits, writing data again, it takes about 3ms ~ 10ms to pass a time, if the PLC modifies other bits of the word register during this period, when writing data again,
it will restore the data at the point to the initial value when reading, which may cause equipment failure, please turn on this function carefully.

![](../icons/SectionExpanded.png)参见

#### 引用

[IReadWriteA3C 接口](08311113-f4d5-8321-79ab-3aa3b26d9a6f.htm)

[HslCommunication.Profinet.Melsec.Helper 命名空间](2922ab51-7617-dc63-33ec-7505340acd85.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Format 属性 

[原文連結](http://api.hslcommunication.cn/html/bb00688f-68ac-927d-3f53-c021515ccb63.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Melsec.Helper](../html/2922ab51-7617-dc63-33ec-7505340acd85.htm "HslCommunication.Profinet.Melsec.Helper")

[IReadWriteA3C 接口](../html/08311113-f4d5-8321-79ab-3aa3b26d9a6f.htm "IReadWriteA3C 接口")

[IReadWriteA3C 属性](../html/a2ec1b91-2087-dd88-e1ab-d9b9023b7790.htm "IReadWriteA3C 属性")

[EnableWriteBitToWordRegister 属性](../html/68ca52a7-8d7e-ce6d-11aa-fd8ee0bcb61b.htm "EnableWriteBitToWordRegister 属性 ")

[Format 属性](../html/bb00688f-68ac-927d-3f53-c021515ccb63.htm "Format 属性 ")

[Station 属性](../html/e6637ca9-5a39-00bd-68de-6a6c1da07af9.htm "Station 属性 ")

[SumCheck 属性](../html/165338dc-5988-fe3e-320f-fb46ce851b1b.htm "SumCheck 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IReadWriteA3CFormat 属性 |

当前的A3C协议的格式信息，可选格式1，2，3，4，默认格式1  
Format information of the current A3C protocol, optional format 1, 2, 3, 4, default format 1

**命名空间：**
 [HslCommunication.Profinet.Melsec.Helper](2922ab51-7617-dc63-33ec-7505340acd85.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
int Format { get; set; }
```

```
Property Format As Integer
	Get
	Set
```

```
property int Format {
	int get ();
	void set (int value);
}
```

```
abstract Format : int with get, set
```

#### 属性值

类型：Int32

![](../icons/SectionExpanded.png)参见

#### 引用

[IReadWriteA3C 接口](08311113-f4d5-8321-79ab-3aa3b26d9a6f.htm)

[HslCommunication.Profinet.Melsec.Helper 命名空间](2922ab51-7617-dc63-33ec-7505340acd85.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Station 属性 

[原文連結](http://api.hslcommunication.cn/html/e6637ca9-5a39-00bd-68de-6a6c1da07af9.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Melsec.Helper](../html/2922ab51-7617-dc63-33ec-7505340acd85.htm "HslCommunication.Profinet.Melsec.Helper")

[IReadWriteA3C 接口](../html/08311113-f4d5-8321-79ab-3aa3b26d9a6f.htm "IReadWriteA3C 接口")

[IReadWriteA3C 属性](../html/a2ec1b91-2087-dd88-e1ab-d9b9023b7790.htm "IReadWriteA3C 属性")

[EnableWriteBitToWordRegister 属性](../html/68ca52a7-8d7e-ce6d-11aa-fd8ee0bcb61b.htm "EnableWriteBitToWordRegister 属性 ")

[Format 属性](../html/bb00688f-68ac-927d-3f53-c021515ccb63.htm "Format 属性 ")

[Station 属性](../html/e6637ca9-5a39-00bd-68de-6a6c1da07af9.htm "Station 属性 ")

[SumCheck 属性](../html/165338dc-5988-fe3e-320f-fb46ce851b1b.htm "SumCheck 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IReadWriteA3CStation 属性 |

当前A3C协议的站编号信息  
Station number information of the current A3C protocol

**命名空间：**
 [HslCommunication.Profinet.Melsec.Helper](2922ab51-7617-dc63-33ec-7505340acd85.htm)  
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

![](../icons/SectionExpanded.png)参见

#### 引用

[IReadWriteA3C 接口](08311113-f4d5-8321-79ab-3aa3b26d9a6f.htm)

[HslCommunication.Profinet.Melsec.Helper 命名空间](2922ab51-7617-dc63-33ec-7505340acd85.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## SumCheck 属性 

[原文連結](http://api.hslcommunication.cn/html/165338dc-5988-fe3e-320f-fb46ce851b1b.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Melsec.Helper](../html/2922ab51-7617-dc63-33ec-7505340acd85.htm "HslCommunication.Profinet.Melsec.Helper")

[IReadWriteA3C 接口](../html/08311113-f4d5-8321-79ab-3aa3b26d9a6f.htm "IReadWriteA3C 接口")

[IReadWriteA3C 属性](../html/a2ec1b91-2087-dd88-e1ab-d9b9023b7790.htm "IReadWriteA3C 属性")

[EnableWriteBitToWordRegister 属性](../html/68ca52a7-8d7e-ce6d-11aa-fd8ee0bcb61b.htm "EnableWriteBitToWordRegister 属性 ")

[Format 属性](../html/bb00688f-68ac-927d-3f53-c021515ccb63.htm "Format 属性 ")

[Station 属性](../html/e6637ca9-5a39-00bd-68de-6a6c1da07af9.htm "Station 属性 ")

[SumCheck 属性](../html/165338dc-5988-fe3e-320f-fb46ce851b1b.htm "SumCheck 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IReadWriteA3CSumCheck 属性 |

当前的A3C协议是否使用和校验，默认使用  
Whether the current A3C protocol uses sum check, it is used by default

**命名空间：**
 [HslCommunication.Profinet.Melsec.Helper](2922ab51-7617-dc63-33ec-7505340acd85.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
bool SumCheck { get; set; }
```

```
Property SumCheck As Boolean
	Get
	Set
```

```
property bool SumCheck {
	bool get ();
	void set (bool value);
}
```

```
abstract SumCheck : bool with get, set
```

#### 属性值

类型：Boolean

![](../icons/SectionExpanded.png)参见

#### 引用

[IReadWriteA3C 接口](08311113-f4d5-8321-79ab-3aa3b26d9a6f.htm)

[HslCommunication.Profinet.Melsec.Helper 命名空间](2922ab51-7617-dc63-33ec-7505340acd85.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## IReadWriteA3C 方法

[原文連結](http://api.hslcommunication.cn/html/e72bcc42-1c30-a900-a811-6adbdf3fe2ca.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Melsec.Helper](../html/2922ab51-7617-dc63-33ec-7505340acd85.htm "HslCommunication.Profinet.Melsec.Helper")

[IReadWriteA3C 接口](../html/08311113-f4d5-8321-79ab-3aa3b26d9a6f.htm "IReadWriteA3C 接口")

[IReadWriteA3C 属性](../html/a2ec1b91-2087-dd88-e1ab-d9b9023b7790.htm "IReadWriteA3C 属性")

[IReadWriteA3C 方法](../html/e72bcc42-1c30-a900-a811-6adbdf3fe2ca.htm "IReadWriteA3C 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IReadWriteA3C 方法 |

[IReadWriteA3C](08311113-f4d5-8321-79ab-3aa3b26d9a6f.htm) 类型公开以下成员。

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

[IReadWriteA3C 接口](08311113-f4d5-8321-79ab-3aa3b26d9a6f.htm)

[HslCommunication.Profinet.Melsec.Helper 命名空间](2922ab51-7617-dc63-33ec-7505340acd85.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## IReadWriteFxLinks 接口

[原文連結](http://api.hslcommunication.cn/html/44bfc17d-ccc2-b900-f4ec-1713ef97627f.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Melsec.Helper](../html/2922ab51-7617-dc63-33ec-7505340acd85.htm "HslCommunication.Profinet.Melsec.Helper")

[IReadWriteFxLinks 接口](../html/44bfc17d-ccc2-b900-f4ec-1713ef97627f.htm "IReadWriteFxLinks 接口")

[IReadWriteFxLinks 属性](../html/b0bf94f9-973b-fcc7-239c-9982053c8bcc.htm "IReadWriteFxLinks 属性")

[IReadWriteFxLinks 方法](../html/0465799c-8cd1-228e-b435-00bcdc4931e6.htm "IReadWriteFxLinks 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IReadWriteFxLinks 接口 |

三菱的FxLink协议接口的设备信息

**命名空间：**
 [HslCommunication.Profinet.Melsec.Helper](2922ab51-7617-dc63-33ec-7505340acd85.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public interface IReadWriteFxLinks : IReadWriteDevice, 
	IReadWriteNet
```

```
Public Interface IReadWriteFxLinks
	Inherits IReadWriteDevice, IReadWriteNet
```

```
public interface class IReadWriteFxLinks : IReadWriteDevice, 
	IReadWriteNet
```

```
type IReadWriteFxLinks =  
    interface
        interface IReadWriteDevice
        interface IReadWriteNet
    end
```

IReadWriteFxLinks 类型公开以下成员。

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性代码示例 | [ByteTransform](0836ec33-d464-ed39-37ff-b2ada0d286fc.htm) | 当前的数据变换机制，当你需要从字节数据转换类型数据的时候需要。  The current data transformation mechanism is required when you need to convert type data from byte data. (继承自 [IReadWriteDevice](d987afcc-f99c-6aed-3847-5b1890504635.htm)。) |
| 公共属性 | [ConnectionId](20787592-e3ec-e0d0-27c2-abfa24654bd2.htm) | 当前连接的唯一ID号，默认为长度20的guid码加随机数组成，方便列表管理，也可以自己指定  The unique ID number of the current connection. The default is a 20-digit guid code plus a random number. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共属性 | [Format](7dec45f5-a4e1-8497-c852-c5de39321b27.htm) | 当前的PLC的Fxlinks协议格式，通常是格式1，或是格式4，所以此处可以设置1，或者是4  The current PLC Fxlinks protocol format is usually format 1 or format 4, so it can be set to 1 or 4 here |
| 公共属性代码示例 | [LogNet](38760c1f-22e3-f8fe-7512-15849a4b6eb0.htm) | 组件的日志工具，支持日志记录，只要实例化后，当前网络的基本信息，就以[DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm)等级进行输出  The component's logging tool supports logging. As long as the instantiation of the basic network information, the output will be output at [DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm) (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共属性 | [Station](60bb7345-0b70-9563-a737-edb4253b9dd2.htm) | PLC的当前的站号，需要根据实际的值来设定，默认是0  The current station number of the PLC needs to be set according to the actual value. The default is 0. |
| 公共属性 | [SumCheck](5438825c-e258-eba7-1872-f8845f779cc4.htm) | 是否启动和校验  Whether to start and sum verify |
| 公共属性 | [WaittingTime](9be693fe-8dbd-cbe0-79fe-73d090755e23.htm) | 报文等待时间，单位10ms，设置范围为0-15  Message waiting time, unit is 10ms, setting range is 0-15 |

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
| 公共方法 | [ReadPlcType](31ce67d2-6daf-2611-b772-5a492c329779.htm) | **[商业授权]** 读取PLC的型号信息，可以携带额外的参数信息，指定站号。举例：s=2; 注意：分号是必须的。 **[Authorization]** Read the PLC model information, you can carry additional parameter information, and specify the station number. Example: s=2; Note: The semicolon is required. |
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
| 公共方法 | [StartPLC](cd88fcc3-a14a-98a7-7d71-60ee8df77c8b.htm) | **[商业授权]** 启动PLC的操作，可以携带额外的参数信息，指定站号。举例：s=2; 注意：分号是必须的。 **[Authorization]** Start the PLC operation, you can carry additional parameter information and specify the station number. Example: s=2; Note: The semicolon is required. |
| 公共方法 | [StopPLC](7b157d89-f30c-90de-b8c0-75e4045456a0.htm) | **[商业授权]** 停止PLC的操作，可以携带额外的参数信息，指定站号。举例：s=2; 注意：分号是必须的。 **[Authorization]** Stop PLC operation, you can carry additional parameter information and specify the station number. Example: s=2; Note: The semicolon is required. |
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

[HslCommunication.Profinet.Melsec.Helper 命名空间](2922ab51-7617-dc63-33ec-7505340acd85.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## IReadWriteFxLinks 属性

[原文連結](http://api.hslcommunication.cn/html/b0bf94f9-973b-fcc7-239c-9982053c8bcc.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Melsec.Helper](../html/2922ab51-7617-dc63-33ec-7505340acd85.htm "HslCommunication.Profinet.Melsec.Helper")

[IReadWriteFxLinks 接口](../html/44bfc17d-ccc2-b900-f4ec-1713ef97627f.htm "IReadWriteFxLinks 接口")

[IReadWriteFxLinks 属性](../html/b0bf94f9-973b-fcc7-239c-9982053c8bcc.htm "IReadWriteFxLinks 属性")

[Format 属性](../html/7dec45f5-a4e1-8497-c852-c5de39321b27.htm "Format 属性 ")

[Station 属性](../html/60bb7345-0b70-9563-a737-edb4253b9dd2.htm "Station 属性 ")

[SumCheck 属性](../html/5438825c-e258-eba7-1872-f8845f779cc4.htm "SumCheck 属性 ")

[WaittingTime 属性](../html/9be693fe-8dbd-cbe0-79fe-73d090755e23.htm "WaittingTime 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IReadWriteFxLinks 属性 |

[IReadWriteFxLinks](44bfc17d-ccc2-b900-f4ec-1713ef97627f.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性代码示例 | [ByteTransform](0836ec33-d464-ed39-37ff-b2ada0d286fc.htm) | 当前的数据变换机制，当你需要从字节数据转换类型数据的时候需要。  The current data transformation mechanism is required when you need to convert type data from byte data. (继承自 [IReadWriteDevice](d987afcc-f99c-6aed-3847-5b1890504635.htm)。) |
| 公共属性 | [ConnectionId](20787592-e3ec-e0d0-27c2-abfa24654bd2.htm) | 当前连接的唯一ID号，默认为长度20的guid码加随机数组成，方便列表管理，也可以自己指定  The unique ID number of the current connection. The default is a 20-digit guid code plus a random number. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共属性 | [Format](7dec45f5-a4e1-8497-c852-c5de39321b27.htm) | 当前的PLC的Fxlinks协议格式，通常是格式1，或是格式4，所以此处可以设置1，或者是4  The current PLC Fxlinks protocol format is usually format 1 or format 4, so it can be set to 1 or 4 here |
| 公共属性代码示例 | [LogNet](38760c1f-22e3-f8fe-7512-15849a4b6eb0.htm) | 组件的日志工具，支持日志记录，只要实例化后，当前网络的基本信息，就以[DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm)等级进行输出  The component's logging tool supports logging. As long as the instantiation of the basic network information, the output will be output at [DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm) (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共属性 | [Station](60bb7345-0b70-9563-a737-edb4253b9dd2.htm) | PLC的当前的站号，需要根据实际的值来设定，默认是0  The current station number of the PLC needs to be set according to the actual value. The default is 0. |
| 公共属性 | [SumCheck](5438825c-e258-eba7-1872-f8845f779cc4.htm) | 是否启动和校验  Whether to start and sum verify |
| 公共属性 | [WaittingTime](9be693fe-8dbd-cbe0-79fe-73d090755e23.htm) | 报文等待时间，单位10ms，设置范围为0-15  Message waiting time, unit is 10ms, setting range is 0-15 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[IReadWriteFxLinks 接口](44bfc17d-ccc2-b900-f4ec-1713ef97627f.htm)

[HslCommunication.Profinet.Melsec.Helper 命名空间](2922ab51-7617-dc63-33ec-7505340acd85.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Format 属性 

[原文連結](http://api.hslcommunication.cn/html/7dec45f5-a4e1-8497-c852-c5de39321b27.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Melsec.Helper](../html/2922ab51-7617-dc63-33ec-7505340acd85.htm "HslCommunication.Profinet.Melsec.Helper")

[IReadWriteFxLinks 接口](../html/44bfc17d-ccc2-b900-f4ec-1713ef97627f.htm "IReadWriteFxLinks 接口")

[IReadWriteFxLinks 属性](../html/b0bf94f9-973b-fcc7-239c-9982053c8bcc.htm "IReadWriteFxLinks 属性")

[Format 属性](../html/7dec45f5-a4e1-8497-c852-c5de39321b27.htm "Format 属性 ")

[Station 属性](../html/60bb7345-0b70-9563-a737-edb4253b9dd2.htm "Station 属性 ")

[SumCheck 属性](../html/5438825c-e258-eba7-1872-f8845f779cc4.htm "SumCheck 属性 ")

[WaittingTime 属性](../html/9be693fe-8dbd-cbe0-79fe-73d090755e23.htm "WaittingTime 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IReadWriteFxLinksFormat 属性 |

当前的PLC的Fxlinks协议格式，通常是格式1，或是格式4，所以此处可以设置1，或者是4  
The current PLC Fxlinks protocol format is usually format 1 or format 4, so it can be set to 1 or 4 here

**命名空间：**
 [HslCommunication.Profinet.Melsec.Helper](2922ab51-7617-dc63-33ec-7505340acd85.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
int Format { get; set; }
```

```
Property Format As Integer
	Get
	Set
```

```
property int Format {
	int get ();
	void set (int value);
}
```

```
abstract Format : int with get, set
```

#### 属性值

类型：Int32

![](../icons/SectionExpanded.png)参见

#### 引用

[IReadWriteFxLinks 接口](44bfc17d-ccc2-b900-f4ec-1713ef97627f.htm)

[HslCommunication.Profinet.Melsec.Helper 命名空间](2922ab51-7617-dc63-33ec-7505340acd85.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Station 属性 

[原文連結](http://api.hslcommunication.cn/html/60bb7345-0b70-9563-a737-edb4253b9dd2.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Melsec.Helper](../html/2922ab51-7617-dc63-33ec-7505340acd85.htm "HslCommunication.Profinet.Melsec.Helper")

[IReadWriteFxLinks 接口](../html/44bfc17d-ccc2-b900-f4ec-1713ef97627f.htm "IReadWriteFxLinks 接口")

[IReadWriteFxLinks 属性](../html/b0bf94f9-973b-fcc7-239c-9982053c8bcc.htm "IReadWriteFxLinks 属性")

[Format 属性](../html/7dec45f5-a4e1-8497-c852-c5de39321b27.htm "Format 属性 ")

[Station 属性](../html/60bb7345-0b70-9563-a737-edb4253b9dd2.htm "Station 属性 ")

[SumCheck 属性](../html/5438825c-e258-eba7-1872-f8845f779cc4.htm "SumCheck 属性 ")

[WaittingTime 属性](../html/9be693fe-8dbd-cbe0-79fe-73d090755e23.htm "WaittingTime 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IReadWriteFxLinksStation 属性 |

PLC的当前的站号，需要根据实际的值来设定，默认是0  
The current station number of the PLC needs to be set according to the actual value. The default is 0.

**命名空间：**
 [HslCommunication.Profinet.Melsec.Helper](2922ab51-7617-dc63-33ec-7505340acd85.htm)  
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

![](../icons/SectionExpanded.png)参见

#### 引用

[IReadWriteFxLinks 接口](44bfc17d-ccc2-b900-f4ec-1713ef97627f.htm)

[HslCommunication.Profinet.Melsec.Helper 命名空间](2922ab51-7617-dc63-33ec-7505340acd85.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## SumCheck 属性 

[原文連結](http://api.hslcommunication.cn/html/5438825c-e258-eba7-1872-f8845f779cc4.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Melsec.Helper](../html/2922ab51-7617-dc63-33ec-7505340acd85.htm "HslCommunication.Profinet.Melsec.Helper")

[IReadWriteFxLinks 接口](../html/44bfc17d-ccc2-b900-f4ec-1713ef97627f.htm "IReadWriteFxLinks 接口")

[IReadWriteFxLinks 属性](../html/b0bf94f9-973b-fcc7-239c-9982053c8bcc.htm "IReadWriteFxLinks 属性")

[Format 属性](../html/7dec45f5-a4e1-8497-c852-c5de39321b27.htm "Format 属性 ")

[Station 属性](../html/60bb7345-0b70-9563-a737-edb4253b9dd2.htm "Station 属性 ")

[SumCheck 属性](../html/5438825c-e258-eba7-1872-f8845f779cc4.htm "SumCheck 属性 ")

[WaittingTime 属性](../html/9be693fe-8dbd-cbe0-79fe-73d090755e23.htm "WaittingTime 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IReadWriteFxLinksSumCheck 属性 |

是否启动和校验  
Whether to start and sum verify

**命名空间：**
 [HslCommunication.Profinet.Melsec.Helper](2922ab51-7617-dc63-33ec-7505340acd85.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
bool SumCheck { get; set; }
```

```
Property SumCheck As Boolean
	Get
	Set
```

```
property bool SumCheck {
	bool get ();
	void set (bool value);
}
```

```
abstract SumCheck : bool with get, set
```

#### 属性值

类型：Boolean

![](../icons/SectionExpanded.png)参见

#### 引用

[IReadWriteFxLinks 接口](44bfc17d-ccc2-b900-f4ec-1713ef97627f.htm)

[HslCommunication.Profinet.Melsec.Helper 命名空间](2922ab51-7617-dc63-33ec-7505340acd85.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## WaittingTime 属性 

[原文連結](http://api.hslcommunication.cn/html/9be693fe-8dbd-cbe0-79fe-73d090755e23.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Melsec.Helper](../html/2922ab51-7617-dc63-33ec-7505340acd85.htm "HslCommunication.Profinet.Melsec.Helper")

[IReadWriteFxLinks 接口](../html/44bfc17d-ccc2-b900-f4ec-1713ef97627f.htm "IReadWriteFxLinks 接口")

[IReadWriteFxLinks 属性](../html/b0bf94f9-973b-fcc7-239c-9982053c8bcc.htm "IReadWriteFxLinks 属性")

[Format 属性](../html/7dec45f5-a4e1-8497-c852-c5de39321b27.htm "Format 属性 ")

[Station 属性](../html/60bb7345-0b70-9563-a737-edb4253b9dd2.htm "Station 属性 ")

[SumCheck 属性](../html/5438825c-e258-eba7-1872-f8845f779cc4.htm "SumCheck 属性 ")

[WaittingTime 属性](../html/9be693fe-8dbd-cbe0-79fe-73d090755e23.htm "WaittingTime 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IReadWriteFxLinksWaittingTime 属性 |

报文等待时间，单位10ms，设置范围为0-15  
Message waiting time, unit is 10ms, setting range is 0-15

**命名空间：**
 [HslCommunication.Profinet.Melsec.Helper](2922ab51-7617-dc63-33ec-7505340acd85.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
byte WaittingTime { get; set; }
```

```
Property WaittingTime As Byte
	Get
	Set
```

```
property unsigned char WaittingTime {
	unsigned char get ();
	void set (unsigned char value);
}
```

```
abstract WaittingTime : byte with get, set
```

#### 属性值

类型：Byte

![](../icons/SectionExpanded.png)参见

#### 引用

[IReadWriteFxLinks 接口](44bfc17d-ccc2-b900-f4ec-1713ef97627f.htm)

[HslCommunication.Profinet.Melsec.Helper 命名空间](2922ab51-7617-dc63-33ec-7505340acd85.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## IReadWriteFxLinks 方法

[原文連結](http://api.hslcommunication.cn/html/0465799c-8cd1-228e-b435-00bcdc4931e6.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Melsec.Helper](../html/2922ab51-7617-dc63-33ec-7505340acd85.htm "HslCommunication.Profinet.Melsec.Helper")

[IReadWriteFxLinks 接口](../html/44bfc17d-ccc2-b900-f4ec-1713ef97627f.htm "IReadWriteFxLinks 接口")

[IReadWriteFxLinks 方法](../html/0465799c-8cd1-228e-b435-00bcdc4931e6.htm "IReadWriteFxLinks 方法")

[ReadPlcType 方法](../html/31ce67d2-6daf-2611-b772-5a492c329779.htm "ReadPlcType 方法 ")

[StartPLC 方法](../html/cd88fcc3-a14a-98a7-7d71-60ee8df77c8b.htm "StartPLC 方法 ")

[StopPLC 方法](../html/7b157d89-f30c-90de-b8c0-75e4045456a0.htm "StopPLC 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IReadWriteFxLinks 方法 |

[IReadWriteFxLinks](44bfc17d-ccc2-b900-f4ec-1713ef97627f.htm) 类型公开以下成员。

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
| 公共方法 | [ReadPlcType](31ce67d2-6daf-2611-b772-5a492c329779.htm) | **[商业授权]** 读取PLC的型号信息，可以携带额外的参数信息，指定站号。举例：s=2; 注意：分号是必须的。 **[Authorization]** Read the PLC model information, you can carry additional parameter information, and specify the station number. Example: s=2; Note: The semicolon is required. |
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
| 公共方法 | [StartPLC](cd88fcc3-a14a-98a7-7d71-60ee8df77c8b.htm) | **[商业授权]** 启动PLC的操作，可以携带额外的参数信息，指定站号。举例：s=2; 注意：分号是必须的。 **[Authorization]** Start the PLC operation, you can carry additional parameter information and specify the station number. Example: s=2; Note: The semicolon is required. |
| 公共方法 | [StopPLC](7b157d89-f30c-90de-b8c0-75e4045456a0.htm) | **[商业授权]** 停止PLC的操作，可以携带额外的参数信息，指定站号。举例：s=2; 注意：分号是必须的。 **[Authorization]** Stop PLC operation, you can carry additional parameter information and specify the station number. Example: s=2; Note: The semicolon is required. |
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

[IReadWriteFxLinks 接口](44bfc17d-ccc2-b900-f4ec-1713ef97627f.htm)

[HslCommunication.Profinet.Melsec.Helper 命名空间](2922ab51-7617-dc63-33ec-7505340acd85.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadPlcType 方法 

[原文連結](http://api.hslcommunication.cn/html/31ce67d2-6daf-2611-b772-5a492c329779.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Melsec.Helper](../html/2922ab51-7617-dc63-33ec-7505340acd85.htm "HslCommunication.Profinet.Melsec.Helper")

[IReadWriteFxLinks 接口](../html/44bfc17d-ccc2-b900-f4ec-1713ef97627f.htm "IReadWriteFxLinks 接口")

[IReadWriteFxLinks 方法](../html/0465799c-8cd1-228e-b435-00bcdc4931e6.htm "IReadWriteFxLinks 方法")

[ReadPlcType 方法](../html/31ce67d2-6daf-2611-b772-5a492c329779.htm "ReadPlcType 方法 ")

[StartPLC 方法](../html/cd88fcc3-a14a-98a7-7d71-60ee8df77c8b.htm "StartPLC 方法 ")

[StopPLC 方法](../html/7b157d89-f30c-90de-b8c0-75e4045456a0.htm "StopPLC 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IReadWriteFxLinksReadPlcType 方法 |

**[商业授权]** 读取PLC的型号信息，可以携带额外的参数信息，指定站号。举例：s=2; 注意：分号是必须的。  
**[Authorization]** Read the PLC model information, you can carry additional parameter information, and specify the station number. Example: s=2; Note: The semicolon is required.

**命名空间：**
 [HslCommunication.Profinet.Melsec.Helper](2922ab51-7617-dc63-33ec-7505340acd85.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
OperateResult<string> ReadPlcType(
	string parameter = ""
)
```

```
Function ReadPlcType ( 
	Optional parameter As String = ""
) As OperateResult(Of String)
```

```
OperateResult<String^>^ ReadPlcType(
	String^ parameter = L""
)
```

```
abstract ReadPlcType : 
        ?parameter : string 
(* Defaults:
        let _parameter = defaultArg parameter ""
*)
-> OperateResult<string> 
```

#### 参数

parameter (Optional)
:   类型：SystemString  
    允许携带的参数信息，例如s=2; 也可以为空

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)String  
带PLC型号的结果信息

![](../icons/SectionExpanded.png)参见

#### 引用

[IReadWriteFxLinks 接口](44bfc17d-ccc2-b900-f4ec-1713ef97627f.htm)

[HslCommunication.Profinet.Melsec.Helper 命名空间](2922ab51-7617-dc63-33ec-7505340acd85.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## StartPLC 方法 

[原文連結](http://api.hslcommunication.cn/html/cd88fcc3-a14a-98a7-7d71-60ee8df77c8b.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Melsec.Helper](../html/2922ab51-7617-dc63-33ec-7505340acd85.htm "HslCommunication.Profinet.Melsec.Helper")

[IReadWriteFxLinks 接口](../html/44bfc17d-ccc2-b900-f4ec-1713ef97627f.htm "IReadWriteFxLinks 接口")

[IReadWriteFxLinks 方法](../html/0465799c-8cd1-228e-b435-00bcdc4931e6.htm "IReadWriteFxLinks 方法")

[ReadPlcType 方法](../html/31ce67d2-6daf-2611-b772-5a492c329779.htm "ReadPlcType 方法 ")

[StartPLC 方法](../html/cd88fcc3-a14a-98a7-7d71-60ee8df77c8b.htm "StartPLC 方法 ")

[StopPLC 方法](../html/7b157d89-f30c-90de-b8c0-75e4045456a0.htm "StopPLC 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IReadWriteFxLinksStartPLC 方法 |

**[商业授权]** 启动PLC的操作，可以携带额外的参数信息，指定站号。举例：s=2; 注意：分号是必须的。  
**[Authorization]** Start the PLC operation, you can carry additional parameter information and specify the station number. Example: s=2; Note: The semicolon is required.

**命名空间：**
 [HslCommunication.Profinet.Melsec.Helper](2922ab51-7617-dc63-33ec-7505340acd85.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
OperateResult StartPLC(
	string parameter = ""
)
```

```
Function StartPLC ( 
	Optional parameter As String = ""
) As OperateResult
```

```
OperateResult^ StartPLC(
	String^ parameter = L""
)
```

```
abstract StartPLC : 
        ?parameter : string 
(* Defaults:
        let _parameter = defaultArg parameter ""
*)
-> OperateResult 
```

#### 参数

parameter (Optional)
:   类型：SystemString  
    允许携带的参数信息，例如s=2; 也可以为空

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否启动成功

![](../icons/SectionExpanded.png)参见

#### 引用

[IReadWriteFxLinks 接口](44bfc17d-ccc2-b900-f4ec-1713ef97627f.htm)

[HslCommunication.Profinet.Melsec.Helper 命名空间](2922ab51-7617-dc63-33ec-7505340acd85.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## StopPLC 方法 

[原文連結](http://api.hslcommunication.cn/html/7b157d89-f30c-90de-b8c0-75e4045456a0.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Melsec.Helper](../html/2922ab51-7617-dc63-33ec-7505340acd85.htm "HslCommunication.Profinet.Melsec.Helper")

[IReadWriteFxLinks 接口](../html/44bfc17d-ccc2-b900-f4ec-1713ef97627f.htm "IReadWriteFxLinks 接口")

[IReadWriteFxLinks 方法](../html/0465799c-8cd1-228e-b435-00bcdc4931e6.htm "IReadWriteFxLinks 方法")

[ReadPlcType 方法](../html/31ce67d2-6daf-2611-b772-5a492c329779.htm "ReadPlcType 方法 ")

[StartPLC 方法](../html/cd88fcc3-a14a-98a7-7d71-60ee8df77c8b.htm "StartPLC 方法 ")

[StopPLC 方法](../html/7b157d89-f30c-90de-b8c0-75e4045456a0.htm "StopPLC 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IReadWriteFxLinksStopPLC 方法 |

**[商业授权]** 停止PLC的操作，可以携带额外的参数信息，指定站号。举例：s=2; 注意：分号是必须的。  
**[Authorization]** Stop PLC operation, you can carry additional parameter information and specify the station number. Example: s=2; Note: The semicolon is required.

**命名空间：**
 [HslCommunication.Profinet.Melsec.Helper](2922ab51-7617-dc63-33ec-7505340acd85.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
OperateResult StopPLC(
	string parameter = ""
)
```

```
Function StopPLC ( 
	Optional parameter As String = ""
) As OperateResult
```

```
OperateResult^ StopPLC(
	String^ parameter = L""
)
```

```
abstract StopPLC : 
        ?parameter : string 
(* Defaults:
        let _parameter = defaultArg parameter ""
*)
-> OperateResult 
```

#### 参数

parameter (Optional)
:   类型：SystemString  
    允许携带的参数信息，例如s=2; 也可以为空

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否停止成功

![](../icons/SectionExpanded.png)参见

#### 引用

[IReadWriteFxLinks 接口](44bfc17d-ccc2-b900-f4ec-1713ef97627f.htm)

[HslCommunication.Profinet.Melsec.Helper 命名空间](2922ab51-7617-dc63-33ec-7505340acd85.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## IReadWriteMc 接口

[原文連結](http://api.hslcommunication.cn/html/6ed75f4f-6587-8991-17e0-3899347e38d8.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Melsec.Helper](../html/2922ab51-7617-dc63-33ec-7505340acd85.htm "HslCommunication.Profinet.Melsec.Helper")

[IReadWriteMc 接口](../html/6ed75f4f-6587-8991-17e0-3899347e38d8.htm "IReadWriteMc 接口")

[IReadWriteMc 属性](../html/7f688e25-c3cf-a8ce-bbf0-5e98efcaf36e.htm "IReadWriteMc 属性")

[IReadWriteMc 方法](../html/a888d766-46f8-0f62-1085-6dda227cee59.htm "IReadWriteMc 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IReadWriteMc 接口 |

基于MC协议的标准的设备接口，适用任何基于MC协议的PLC设备，主要是三菱，基恩士，松下的PLC设备。  
The standard equipment interface based on MC protocol is suitable for any PLC equipment based on MC protocol,
mainly PLC equipment from Mitsubishi, Keyence, and Panasonic.

**命名空间：**
 [HslCommunication.Profinet.Melsec.Helper](2922ab51-7617-dc63-33ec-7505340acd85.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public interface IReadWriteMc : IReadWriteDevice, 
	IReadWriteNet
```

```
Public Interface IReadWriteMc
	Inherits IReadWriteDevice, IReadWriteNet
```

```
public interface class IReadWriteMc : IReadWriteDevice, 
	IReadWriteNet
```

```
type IReadWriteMc =  
    interface
        interface IReadWriteDevice
        interface IReadWriteNet
    end
```

IReadWriteMc 类型公开以下成员。

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性代码示例 | [ByteTransform](0836ec33-d464-ed39-37ff-b2ada0d286fc.htm) | 当前的数据变换机制，当你需要从字节数据转换类型数据的时候需要。  The current data transformation mechanism is required when you need to convert type data from byte data. (继承自 [IReadWriteDevice](d987afcc-f99c-6aed-3847-5b1890504635.htm)。) |
| 公共属性 | [ConnectionId](20787592-e3ec-e0d0-27c2-abfa24654bd2.htm) | 当前连接的唯一ID号，默认为长度20的guid码加随机数组成，方便列表管理，也可以自己指定  The unique ID number of the current connection. The default is a 20-digit guid code plus a random number. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共属性 | [EnableWriteBitToWordRegister](c9f3d44f-f8f1-5d1e-f113-49e4d9235164.htm) | 是否开启支持写入位到字寄存器的功能，该功能先读取字寄存器的字数据，然后修改其中的位，再写入回去，可能存在脏数据的风险  Whether to enable supporting the function of writing bit-to-word registers, which first reads the word data of the word register, then modifies the bits in it, and then writes back, there may be a risk of dirty data |
| 公共属性代码示例 | [LogNet](38760c1f-22e3-f8fe-7512-15849a4b6eb0.htm) | 组件的日志工具，支持日志记录，只要实例化后，当前网络的基本信息，就以[DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm)等级进行输出  The component's logging tool supports logging. As long as the instantiation of the basic network information, the output will be output at [DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm) (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共属性 | [McType](7cd3931e-bf0d-7fa5-7c96-3246fef747e4.htm) | 当前的MC协议的格式类型  The format type of the current MC protocol |
| 公共属性 | [NetworkNumber](e0049030-86a8-66d7-f05c-26366acc93ae.htm) | 网络号，通常为0  Network number, usually 0 |
| 公共属性 | [NetworkStationNumber](9db5ba75-b9d7-381c-63cf-42a853820955.htm) | 网络站号，通常为0  Network station number, usually 0 |
| 公共属性 | [PLCNumber](94aab2ad-9fc2-adfc-a907-03d760372070.htm) | PLC编号，如果是本站信息，则是 0xFF 值，其他站则根据实际的情况指定。  The PLC number, if it is the information of this station, is the 0xFF value, and other stations are specified according to the actual situation. |
| 公共属性 | [TargetIOStation](b4eb330f-c026-d7dc-b2a0-b1f9cbaf6284.htm) | 请求目标模块的IO编号，默认是管理CPU，也就是 0x03FF 的值，如果需要访问其他的非管理CPU的时候，请参考手册进行配置相关的值  The IO number of the target module is the default management CPU, that is, the value of 0x03FF, if you need to access other non-management CPUs, please refer to the manual to configure the relevant values |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [ErrorStateReset](a69af7d7-1002-67be-adf6-f63dc0134d58.htm) | LED 熄灭 出错代码初始化  LED off Error code initialization |
| 公共方法 | [ErrorStateResetAsync](a7e2996a-5868-f673-0309-6bad1ac30f0e.htm) | LED 熄灭 出错代码初始化  LED off Error code initialization |
| 公共方法 | [ExtractActualData](bdddaba6-236f-f033-3897-d30097554eda.htm) | 从PLC反馈的数据中提取出实际的数据内容，需要传入反馈数据，是否位读取 |
| 公共方法 | [McAnalysisAddress](83ac01c3-c0ca-d422-8656-482d6ef6821e.htm) | 当前MC协议的分析地址的方法，对传入的字符串格式的地址进行数据解析。  The current MC protocol's address analysis method performs data parsing on the address of the incoming string format. |
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
| 公共方法 | [ReadPlcType](472b1ae2-4412-3ef8-e36a-a683f311e284.htm) | 读取PLC的型号信息，例如 Q02HCPU  Read PLC model information, such as Q02HCPU |
| 公共方法 | [ReadPlcTypeAsync](0e706bb1-5383-d4f7-b6c9-f2ae868b2f51.htm) | 读取PLC的型号信息，例如 Q02HCPU  Read PLC model information, such as Q02HCPU |
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
| 公共方法 | [RemoteReset](5272208a-7575-7843-854b-95bd4cfedc33.htm) | 远程Reset操作  Remote Reset Operation |
| 公共方法 | [RemoteResetAsync](977bbd3a-2629-2974-0aee-b36273173173.htm) | 远程Reset操作  Remote Reset Operation |
| 公共方法 | [RemoteRun](8ec8c923-619b-e03b-84cc-cf2fa3591458.htm) | 远程Run操作  Remote Run Operation |
| 公共方法 | [RemoteRunAsync](2fedee07-72cb-3279-dcd3-de49f0c836a1.htm) | 远程Run操作  Remote Run Operation |
| 公共方法 | [RemoteStop](b1f47349-a2fa-8f69-4d25-cd0bbce8cbfb.htm) | 远程Stop操作  Remote Stop operation |
| 公共方法 | [RemoteStopAsync](4ccf313c-8a40-f131-e408-0f6240f2ed3a.htm) | 远程Stop操作  Remote Stop operation |
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

[HslCommunication.Profinet.Melsec.Helper 命名空间](2922ab51-7617-dc63-33ec-7505340acd85.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## IReadWriteMc 属性

[原文連結](http://api.hslcommunication.cn/html/7f688e25-c3cf-a8ce-bbf0-5e98efcaf36e.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Melsec.Helper](../html/2922ab51-7617-dc63-33ec-7505340acd85.htm "HslCommunication.Profinet.Melsec.Helper")

[IReadWriteMc 接口](../html/6ed75f4f-6587-8991-17e0-3899347e38d8.htm "IReadWriteMc 接口")

[IReadWriteMc 属性](../html/7f688e25-c3cf-a8ce-bbf0-5e98efcaf36e.htm "IReadWriteMc 属性")

[EnableWriteBitToWordRegister 属性](../html/c9f3d44f-f8f1-5d1e-f113-49e4d9235164.htm "EnableWriteBitToWordRegister 属性 ")

[McType 属性](../html/7cd3931e-bf0d-7fa5-7c96-3246fef747e4.htm "McType 属性 ")

[NetworkNumber 属性](../html/e0049030-86a8-66d7-f05c-26366acc93ae.htm "NetworkNumber 属性 ")

[NetworkStationNumber 属性](../html/9db5ba75-b9d7-381c-63cf-42a853820955.htm "NetworkStationNumber 属性 ")

[PLCNumber 属性](../html/94aab2ad-9fc2-adfc-a907-03d760372070.htm "PLCNumber 属性 ")

[TargetIOStation 属性](../html/b4eb330f-c026-d7dc-b2a0-b1f9cbaf6284.htm "TargetIOStation 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IReadWriteMc 属性 |

[IReadWriteMc](6ed75f4f-6587-8991-17e0-3899347e38d8.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性代码示例 | [ByteTransform](0836ec33-d464-ed39-37ff-b2ada0d286fc.htm) | 当前的数据变换机制，当你需要从字节数据转换类型数据的时候需要。  The current data transformation mechanism is required when you need to convert type data from byte data. (继承自 [IReadWriteDevice](d987afcc-f99c-6aed-3847-5b1890504635.htm)。) |
| 公共属性 | [ConnectionId](20787592-e3ec-e0d0-27c2-abfa24654bd2.htm) | 当前连接的唯一ID号，默认为长度20的guid码加随机数组成，方便列表管理，也可以自己指定  The unique ID number of the current connection. The default is a 20-digit guid code plus a random number. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共属性 | [EnableWriteBitToWordRegister](c9f3d44f-f8f1-5d1e-f113-49e4d9235164.htm) | 是否开启支持写入位到字寄存器的功能，该功能先读取字寄存器的字数据，然后修改其中的位，再写入回去，可能存在脏数据的风险  Whether to enable supporting the function of writing bit-to-word registers, which first reads the word data of the word register, then modifies the bits in it, and then writes back, there may be a risk of dirty data |
| 公共属性代码示例 | [LogNet](38760c1f-22e3-f8fe-7512-15849a4b6eb0.htm) | 组件的日志工具，支持日志记录，只要实例化后，当前网络的基本信息，就以[DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm)等级进行输出  The component's logging tool supports logging. As long as the instantiation of the basic network information, the output will be output at [DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm) (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共属性 | [McType](7cd3931e-bf0d-7fa5-7c96-3246fef747e4.htm) | 当前的MC协议的格式类型  The format type of the current MC protocol |
| 公共属性 | [NetworkNumber](e0049030-86a8-66d7-f05c-26366acc93ae.htm) | 网络号，通常为0  Network number, usually 0 |
| 公共属性 | [NetworkStationNumber](9db5ba75-b9d7-381c-63cf-42a853820955.htm) | 网络站号，通常为0  Network station number, usually 0 |
| 公共属性 | [PLCNumber](94aab2ad-9fc2-adfc-a907-03d760372070.htm) | PLC编号，如果是本站信息，则是 0xFF 值，其他站则根据实际的情况指定。  The PLC number, if it is the information of this station, is the 0xFF value, and other stations are specified according to the actual situation. |
| 公共属性 | [TargetIOStation](b4eb330f-c026-d7dc-b2a0-b1f9cbaf6284.htm) | 请求目标模块的IO编号，默认是管理CPU，也就是 0x03FF 的值，如果需要访问其他的非管理CPU的时候，请参考手册进行配置相关的值  The IO number of the target module is the default management CPU, that is, the value of 0x03FF, if you need to access other non-management CPUs, please refer to the manual to configure the relevant values |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[IReadWriteMc 接口](6ed75f4f-6587-8991-17e0-3899347e38d8.htm)

[HslCommunication.Profinet.Melsec.Helper 命名空间](2922ab51-7617-dc63-33ec-7505340acd85.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## EnableWriteBitToWordRegister 属性 

[原文連結](http://api.hslcommunication.cn/html/c9f3d44f-f8f1-5d1e-f113-49e4d9235164.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Melsec.Helper](../html/2922ab51-7617-dc63-33ec-7505340acd85.htm "HslCommunication.Profinet.Melsec.Helper")

[IReadWriteMc 接口](../html/6ed75f4f-6587-8991-17e0-3899347e38d8.htm "IReadWriteMc 接口")

[IReadWriteMc 属性](../html/7f688e25-c3cf-a8ce-bbf0-5e98efcaf36e.htm "IReadWriteMc 属性")

[EnableWriteBitToWordRegister 属性](../html/c9f3d44f-f8f1-5d1e-f113-49e4d9235164.htm "EnableWriteBitToWordRegister 属性 ")

[McType 属性](../html/7cd3931e-bf0d-7fa5-7c96-3246fef747e4.htm "McType 属性 ")

[NetworkNumber 属性](../html/e0049030-86a8-66d7-f05c-26366acc93ae.htm "NetworkNumber 属性 ")

[NetworkStationNumber 属性](../html/9db5ba75-b9d7-381c-63cf-42a853820955.htm "NetworkStationNumber 属性 ")

[PLCNumber 属性](../html/94aab2ad-9fc2-adfc-a907-03d760372070.htm "PLCNumber 属性 ")

[TargetIOStation 属性](../html/b4eb330f-c026-d7dc-b2a0-b1f9cbaf6284.htm "TargetIOStation 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IReadWriteMcEnableWriteBitToWordRegister 属性 |

是否开启支持写入位到字寄存器的功能，该功能先读取字寄存器的字数据，然后修改其中的位，再写入回去，可能存在脏数据的风险  
Whether to enable supporting the function of writing bit-to-word registers, which first reads the word data of the word register, then modifies the bits in it, and then writes back, there may be a risk of dirty data

**命名空间：**
 [HslCommunication.Profinet.Melsec.Helper](2922ab51-7617-dc63-33ec-7505340acd85.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
bool EnableWriteBitToWordRegister { get; set; }
```

```
Property EnableWriteBitToWordRegister As Boolean
	Get
	Set
```

```
property bool EnableWriteBitToWordRegister {
	bool get ();
	void set (bool value);
}
```

```
abstract EnableWriteBitToWordRegister : bool with get, set
```

#### 属性值

类型：Boolean

![](../icons/SectionExpanded.png)备注

关于脏数据风险：从读取数据，修改位，再次写入数据时，大概需要经过3ms~10ms不等的时间，如果此期间内PLC修改了该字寄存器的其他位，再次写入数据时会恢复该点位的数据到读取时的初始值，可能引发设备故障，请谨慎开启此功能。  
About dirty data risk: from reading data, modifying bits, writing data again, it takes about 3ms ~ 10ms to pass a time, if the PLC modifies other bits of the word register during this period, when writing data again,
it will restore the data at the point to the initial value when reading, which may cause equipment failure, please turn on this function carefully.

![](../icons/SectionExpanded.png)参见

#### 引用

[IReadWriteMc 接口](6ed75f4f-6587-8991-17e0-3899347e38d8.htm)

[HslCommunication.Profinet.Melsec.Helper 命名空间](2922ab51-7617-dc63-33ec-7505340acd85.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## McType 属性 

[原文連結](http://api.hslcommunication.cn/html/7cd3931e-bf0d-7fa5-7c96-3246fef747e4.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Melsec.Helper](../html/2922ab51-7617-dc63-33ec-7505340acd85.htm "HslCommunication.Profinet.Melsec.Helper")

[IReadWriteMc 接口](../html/6ed75f4f-6587-8991-17e0-3899347e38d8.htm "IReadWriteMc 接口")

[IReadWriteMc 属性](../html/7f688e25-c3cf-a8ce-bbf0-5e98efcaf36e.htm "IReadWriteMc 属性")

[EnableWriteBitToWordRegister 属性](../html/c9f3d44f-f8f1-5d1e-f113-49e4d9235164.htm "EnableWriteBitToWordRegister 属性 ")

[McType 属性](../html/7cd3931e-bf0d-7fa5-7c96-3246fef747e4.htm "McType 属性 ")

[NetworkNumber 属性](../html/e0049030-86a8-66d7-f05c-26366acc93ae.htm "NetworkNumber 属性 ")

[NetworkStationNumber 属性](../html/9db5ba75-b9d7-381c-63cf-42a853820955.htm "NetworkStationNumber 属性 ")

[PLCNumber 属性](../html/94aab2ad-9fc2-adfc-a907-03d760372070.htm "PLCNumber 属性 ")

[TargetIOStation 属性](../html/b4eb330f-c026-d7dc-b2a0-b1f9cbaf6284.htm "TargetIOStation 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IReadWriteMcMcType 属性 |

当前的MC协议的格式类型  
The format type of the current MC protocol

**命名空间：**
 [HslCommunication.Profinet.Melsec.Helper](2922ab51-7617-dc63-33ec-7505340acd85.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
McType McType { get; }
```

```
ReadOnly Property McType As McType
	Get
```

```
property McType McType {
	McType get ();
}
```

```
abstract McType : McType with get
```

#### 属性值

类型：[McType](26e256cf-18a7-bca3-2fe9-544737862be1.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[IReadWriteMc 接口](6ed75f4f-6587-8991-17e0-3899347e38d8.htm)

[HslCommunication.Profinet.Melsec.Helper 命名空间](2922ab51-7617-dc63-33ec-7505340acd85.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## NetworkNumber 属性 

[原文連結](http://api.hslcommunication.cn/html/e0049030-86a8-66d7-f05c-26366acc93ae.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Melsec.Helper](../html/2922ab51-7617-dc63-33ec-7505340acd85.htm "HslCommunication.Profinet.Melsec.Helper")

[IReadWriteMc 接口](../html/6ed75f4f-6587-8991-17e0-3899347e38d8.htm "IReadWriteMc 接口")

[IReadWriteMc 属性](../html/7f688e25-c3cf-a8ce-bbf0-5e98efcaf36e.htm "IReadWriteMc 属性")

[EnableWriteBitToWordRegister 属性](../html/c9f3d44f-f8f1-5d1e-f113-49e4d9235164.htm "EnableWriteBitToWordRegister 属性 ")

[McType 属性](../html/7cd3931e-bf0d-7fa5-7c96-3246fef747e4.htm "McType 属性 ")

[NetworkNumber 属性](../html/e0049030-86a8-66d7-f05c-26366acc93ae.htm "NetworkNumber 属性 ")

[NetworkStationNumber 属性](../html/9db5ba75-b9d7-381c-63cf-42a853820955.htm "NetworkStationNumber 属性 ")

[PLCNumber 属性](../html/94aab2ad-9fc2-adfc-a907-03d760372070.htm "PLCNumber 属性 ")

[TargetIOStation 属性](../html/b4eb330f-c026-d7dc-b2a0-b1f9cbaf6284.htm "TargetIOStation 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IReadWriteMcNetworkNumber 属性 |

网络号，通常为0  
Network number, usually 0

**命名空间：**
 [HslCommunication.Profinet.Melsec.Helper](2922ab51-7617-dc63-33ec-7505340acd85.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
byte NetworkNumber { get; set; }
```

```
Property NetworkNumber As Byte
	Get
	Set
```

```
property unsigned char NetworkNumber {
	unsigned char get ();
	void set (unsigned char value);
}
```

```
abstract NetworkNumber : byte with get, set
```

#### 属性值

类型：Byte

![](../icons/SectionExpanded.png)备注

依据PLC的配置而配置，如果PLC配置了1，那么此处也填0，如果PLC配置了2，此处就填2，测试不通的话，继续测试0

![](../icons/SectionExpanded.png)参见

#### 引用

[IReadWriteMc 接口](6ed75f4f-6587-8991-17e0-3899347e38d8.htm)

[HslCommunication.Profinet.Melsec.Helper 命名空间](2922ab51-7617-dc63-33ec-7505340acd85.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## NetworkStationNumber 属性 

[原文連結](http://api.hslcommunication.cn/html/9db5ba75-b9d7-381c-63cf-42a853820955.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Melsec.Helper](../html/2922ab51-7617-dc63-33ec-7505340acd85.htm "HslCommunication.Profinet.Melsec.Helper")

[IReadWriteMc 接口](../html/6ed75f4f-6587-8991-17e0-3899347e38d8.htm "IReadWriteMc 接口")

[IReadWriteMc 属性](../html/7f688e25-c3cf-a8ce-bbf0-5e98efcaf36e.htm "IReadWriteMc 属性")

[EnableWriteBitToWordRegister 属性](../html/c9f3d44f-f8f1-5d1e-f113-49e4d9235164.htm "EnableWriteBitToWordRegister 属性 ")

[McType 属性](../html/7cd3931e-bf0d-7fa5-7c96-3246fef747e4.htm "McType 属性 ")

[NetworkNumber 属性](../html/e0049030-86a8-66d7-f05c-26366acc93ae.htm "NetworkNumber 属性 ")

[NetworkStationNumber 属性](../html/9db5ba75-b9d7-381c-63cf-42a853820955.htm "NetworkStationNumber 属性 ")

[PLCNumber 属性](../html/94aab2ad-9fc2-adfc-a907-03d760372070.htm "PLCNumber 属性 ")

[TargetIOStation 属性](../html/b4eb330f-c026-d7dc-b2a0-b1f9cbaf6284.htm "TargetIOStation 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IReadWriteMcNetworkStationNumber 属性 |

网络站号，通常为0  
Network station number, usually 0

**命名空间：**
 [HslCommunication.Profinet.Melsec.Helper](2922ab51-7617-dc63-33ec-7505340acd85.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
byte NetworkStationNumber { get; set; }
```

```
Property NetworkStationNumber As Byte
	Get
	Set
```

```
property unsigned char NetworkStationNumber {
	unsigned char get ();
	void set (unsigned char value);
}
```

```
abstract NetworkStationNumber : byte with get, set
```

#### 属性值

类型：Byte

![](../icons/SectionExpanded.png)备注

依据PLC的配置而配置，如果PLC配置了1，那么此处也填0，如果PLC配置了2，此处就填2，测试不通的话，继续测试0

![](../icons/SectionExpanded.png)参见

#### 引用

[IReadWriteMc 接口](6ed75f4f-6587-8991-17e0-3899347e38d8.htm)

[HslCommunication.Profinet.Melsec.Helper 命名空间](2922ab51-7617-dc63-33ec-7505340acd85.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## PLCNumber 属性 

[原文連結](http://api.hslcommunication.cn/html/94aab2ad-9fc2-adfc-a907-03d760372070.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Melsec.Helper](../html/2922ab51-7617-dc63-33ec-7505340acd85.htm "HslCommunication.Profinet.Melsec.Helper")

[IReadWriteMc 接口](../html/6ed75f4f-6587-8991-17e0-3899347e38d8.htm "IReadWriteMc 接口")

[IReadWriteMc 属性](../html/7f688e25-c3cf-a8ce-bbf0-5e98efcaf36e.htm "IReadWriteMc 属性")

[EnableWriteBitToWordRegister 属性](../html/c9f3d44f-f8f1-5d1e-f113-49e4d9235164.htm "EnableWriteBitToWordRegister 属性 ")

[McType 属性](../html/7cd3931e-bf0d-7fa5-7c96-3246fef747e4.htm "McType 属性 ")

[NetworkNumber 属性](../html/e0049030-86a8-66d7-f05c-26366acc93ae.htm "NetworkNumber 属性 ")

[NetworkStationNumber 属性](../html/9db5ba75-b9d7-381c-63cf-42a853820955.htm "NetworkStationNumber 属性 ")

[PLCNumber 属性](../html/94aab2ad-9fc2-adfc-a907-03d760372070.htm "PLCNumber 属性 ")

[TargetIOStation 属性](../html/b4eb330f-c026-d7dc-b2a0-b1f9cbaf6284.htm "TargetIOStation 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IReadWriteMcPLCNumber 属性 |

PLC编号，如果是本站信息，则是 0xFF 值，其他站则根据实际的情况指定。  
The PLC number, if it is the information of this station, is the 0xFF value, and other stations are specified according to the actual situation.

**命名空间：**
 [HslCommunication.Profinet.Melsec.Helper](2922ab51-7617-dc63-33ec-7505340acd85.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
byte PLCNumber { get; set; }
```

```
Property PLCNumber As Byte
	Get
	Set
```

```
property unsigned char PLCNumber {
	unsigned char get ();
	void set (unsigned char value);
}
```

```
abstract PLCNumber : byte with get, set
```

#### 属性值

类型：Byte

![](../icons/SectionExpanded.png)参见

#### 引用

[IReadWriteMc 接口](6ed75f4f-6587-8991-17e0-3899347e38d8.htm)

[HslCommunication.Profinet.Melsec.Helper 命名空间](2922ab51-7617-dc63-33ec-7505340acd85.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)