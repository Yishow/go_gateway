# HslCommunication - HslCommunication.Profinet.Siemens.Helper

> 分類頁數: 30



---
## HslCommunication.Profinet.Siemens.Helper

[原文連結](http://api.hslcommunication.cn/html/8dfc52a4-e8ce-f3d2-0647-4d18de8416d3.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Siemens.Helper](../html/8dfc52a4-e8ce-f3d2-0647-4d18de8416d3.htm "HslCommunication.Profinet.Siemens.Helper")

[ISiemensPPI 接口](../html/b9b7bddf-d584-6c46-fb23-670c72182056.htm "ISiemensPPI 接口")

[SiemensPPIHelper 类](../html/24deecb2-7eb3-9aba-2da3-137adacdcbd9.htm "SiemensPPIHelper 类")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslCommunication.Profinet.Siemens.Helper 命名空间 |

[缺少 "N:HslCommunication.Profinet.Siemens.Helper" 的 <summary> 文档]

![](../icons/SectionExpanded.png)类

|  | 类 | 说明 |
| --- | --- | --- |
| 公共类 | [SiemensPPIHelper](24deecb2-7eb3-9aba-2da3-137adacdcbd9.htm) | 西门子PPI协议的辅助类对象 |

![](../icons/SectionExpanded.png)接口

|  | 接口 | 说明 |
| --- | --- | --- |
| 公共接口 | [ISiemensPPI](b9b7bddf-d584-6c46-fb23-670c72182056.htm) | 西门子PPI的公用接口信息 |

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ISiemensPPI 接口

[原文連結](http://api.hslcommunication.cn/html/b9b7bddf-d584-6c46-fb23-670c72182056.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Siemens.Helper](../html/8dfc52a4-e8ce-f3d2-0647-4d18de8416d3.htm "HslCommunication.Profinet.Siemens.Helper")

[ISiemensPPI 接口](../html/b9b7bddf-d584-6c46-fb23-670c72182056.htm "ISiemensPPI 接口")

[ISiemensPPI 属性](../html/074a2952-b08b-e25d-b18c-09677149cf63.htm "ISiemensPPI 属性")

[ISiemensPPI 方法](../html/c77ec8ca-32af-590e-a9e8-7cc7cd2ad723.htm "ISiemensPPI 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ISiemensPPI 接口 |

西门子PPI的公用接口信息

**命名空间：**
 [HslCommunication.Profinet.Siemens.Helper](8dfc52a4-e8ce-f3d2-0647-4d18de8416d3.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public interface ISiemensPPI : IReadWriteNet
```

```
Public Interface ISiemensPPI
	Inherits IReadWriteNet
```

```
public interface class ISiemensPPI : IReadWriteNet
```

```
type ISiemensPPI =  
    interface
        interface IReadWriteNet
    end
```

ISiemensPPI 类型公开以下成员。

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [ConnectionId](20787592-e3ec-e0d0-27c2-abfa24654bd2.htm) | 当前连接的唯一ID号，默认为长度20的guid码加随机数组成，方便列表管理，也可以自己指定  The unique ID number of the current connection. The default is a 20-digit guid code plus a random number. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
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
| 公共方法 | [ReadPlcType](86b53070-da19-05b6-2b50-acab61ca51f1.htm) | 读取西门子PLC的型号信息，参数信息可以携带站号信息 "s=2;", 注意，分号是必须的。  Read the model information of Siemens PLC, the parameter information can carry the station number information "s=2;", note that the semicolon is required. |
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
| 公共方法 | [Start](35128174-f97e-ad50-b559-71de1ac32caa.htm) | 启动西门子PLC为RUN模式，参数信息可以携带站号信息 "s=2;", 注意，分号是必须的。  Start Siemens PLC in RUN mode, parameter information can carry station number information "s=2;", note that the semicolon is required. |
| 公共方法 | [Stop](d59f2c81-06f6-743c-2a09-b184b5623996.htm) | 停止西门子PLC，切换为Stop模式，参数信息可以携带站号信息 "s=2;", 注意，分号是必须的。  Stop Siemens PLC and switch to Stop mode, parameter information can carry station number information "s=2;", note that the semicolon is required. |
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

[HslCommunication.Profinet.Siemens.Helper 命名空间](8dfc52a4-e8ce-f3d2-0647-4d18de8416d3.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ISiemensPPI 属性

[原文連結](http://api.hslcommunication.cn/html/074a2952-b08b-e25d-b18c-09677149cf63.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Siemens.Helper](../html/8dfc52a4-e8ce-f3d2-0647-4d18de8416d3.htm "HslCommunication.Profinet.Siemens.Helper")

[ISiemensPPI 接口](../html/b9b7bddf-d584-6c46-fb23-670c72182056.htm "ISiemensPPI 接口")

[ISiemensPPI 属性](../html/074a2952-b08b-e25d-b18c-09677149cf63.htm "ISiemensPPI 属性")

[ISiemensPPI 方法](../html/c77ec8ca-32af-590e-a9e8-7cc7cd2ad723.htm "ISiemensPPI 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ISiemensPPI 属性 |

[ISiemensPPI](b9b7bddf-d584-6c46-fb23-670c72182056.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [ConnectionId](20787592-e3ec-e0d0-27c2-abfa24654bd2.htm) | 当前连接的唯一ID号，默认为长度20的guid码加随机数组成，方便列表管理，也可以自己指定  The unique ID number of the current connection. The default is a 20-digit guid code plus a random number. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共属性代码示例 | [LogNet](38760c1f-22e3-f8fe-7512-15849a4b6eb0.htm) | 组件的日志工具，支持日志记录，只要实例化后，当前网络的基本信息，就以[DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm)等级进行输出  The component's logging tool supports logging. As long as the instantiation of the basic network information, the output will be output at [DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm) (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[ISiemensPPI 接口](b9b7bddf-d584-6c46-fb23-670c72182056.htm)

[HslCommunication.Profinet.Siemens.Helper 命名空间](8dfc52a4-e8ce-f3d2-0647-4d18de8416d3.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ISiemensPPI 方法

[原文連結](http://api.hslcommunication.cn/html/c77ec8ca-32af-590e-a9e8-7cc7cd2ad723.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Siemens.Helper](../html/8dfc52a4-e8ce-f3d2-0647-4d18de8416d3.htm "HslCommunication.Profinet.Siemens.Helper")

[ISiemensPPI 接口](../html/b9b7bddf-d584-6c46-fb23-670c72182056.htm "ISiemensPPI 接口")

[ISiemensPPI 方法](../html/c77ec8ca-32af-590e-a9e8-7cc7cd2ad723.htm "ISiemensPPI 方法")

[ReadPlcType 方法](../html/86b53070-da19-05b6-2b50-acab61ca51f1.htm "ReadPlcType 方法 ")

[Start 方法](../html/35128174-f97e-ad50-b559-71de1ac32caa.htm "Start 方法 ")

[Stop 方法](../html/d59f2c81-06f6-743c-2a09-b184b5623996.htm "Stop 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ISiemensPPI 方法 |

[ISiemensPPI](b9b7bddf-d584-6c46-fb23-670c72182056.htm) 类型公开以下成员。

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
| 公共方法 | [ReadPlcType](86b53070-da19-05b6-2b50-acab61ca51f1.htm) | 读取西门子PLC的型号信息，参数信息可以携带站号信息 "s=2;", 注意，分号是必须的。  Read the model information of Siemens PLC, the parameter information can carry the station number information "s=2;", note that the semicolon is required. |
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
| 公共方法 | [Start](35128174-f97e-ad50-b559-71de1ac32caa.htm) | 启动西门子PLC为RUN模式，参数信息可以携带站号信息 "s=2;", 注意，分号是必须的。  Start Siemens PLC in RUN mode, parameter information can carry station number information "s=2;", note that the semicolon is required. |
| 公共方法 | [Stop](d59f2c81-06f6-743c-2a09-b184b5623996.htm) | 停止西门子PLC，切换为Stop模式，参数信息可以携带站号信息 "s=2;", 注意，分号是必须的。  Stop Siemens PLC and switch to Stop mode, parameter information can carry station number information "s=2;", note that the semicolon is required. |
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

[ISiemensPPI 接口](b9b7bddf-d584-6c46-fb23-670c72182056.htm)

[HslCommunication.Profinet.Siemens.Helper 命名空间](8dfc52a4-e8ce-f3d2-0647-4d18de8416d3.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadPlcType 方法 

[原文連結](http://api.hslcommunication.cn/html/86b53070-da19-05b6-2b50-acab61ca51f1.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Siemens.Helper](../html/8dfc52a4-e8ce-f3d2-0647-4d18de8416d3.htm "HslCommunication.Profinet.Siemens.Helper")

[ISiemensPPI 接口](../html/b9b7bddf-d584-6c46-fb23-670c72182056.htm "ISiemensPPI 接口")

[ISiemensPPI 方法](../html/c77ec8ca-32af-590e-a9e8-7cc7cd2ad723.htm "ISiemensPPI 方法")

[ReadPlcType 方法](../html/86b53070-da19-05b6-2b50-acab61ca51f1.htm "ReadPlcType 方法 ")

[Start 方法](../html/35128174-f97e-ad50-b559-71de1ac32caa.htm "Start 方法 ")

[Stop 方法](../html/d59f2c81-06f6-743c-2a09-b184b5623996.htm "Stop 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ISiemensPPIReadPlcType 方法 |

读取西门子PLC的型号信息，参数信息可以携带站号信息 "s=2;", 注意，分号是必须的。  
Read the model information of Siemens PLC, the parameter information can carry the station number information "s=2;", note that the semicolon is required.

**命名空间：**
 [HslCommunication.Profinet.Siemens.Helper](8dfc52a4-e8ce-f3d2-0647-4d18de8416d3.htm)  
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
    额外的参数信息，例如可以携带站号信息 "s=2;", 注意，分号是必须的。

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)String  
包含是否成功的结果对象

![](../icons/SectionExpanded.png)参见

#### 引用

[ISiemensPPI 接口](b9b7bddf-d584-6c46-fb23-670c72182056.htm)

[HslCommunication.Profinet.Siemens.Helper 命名空间](8dfc52a4-e8ce-f3d2-0647-4d18de8416d3.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Start 方法 

[原文連結](http://api.hslcommunication.cn/html/35128174-f97e-ad50-b559-71de1ac32caa.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Siemens.Helper](../html/8dfc52a4-e8ce-f3d2-0647-4d18de8416d3.htm "HslCommunication.Profinet.Siemens.Helper")

[ISiemensPPI 接口](../html/b9b7bddf-d584-6c46-fb23-670c72182056.htm "ISiemensPPI 接口")

[ISiemensPPI 方法](../html/c77ec8ca-32af-590e-a9e8-7cc7cd2ad723.htm "ISiemensPPI 方法")

[ReadPlcType 方法](../html/86b53070-da19-05b6-2b50-acab61ca51f1.htm "ReadPlcType 方法 ")

[Start 方法](../html/35128174-f97e-ad50-b559-71de1ac32caa.htm "Start 方法 ")

[Stop 方法](../html/d59f2c81-06f6-743c-2a09-b184b5623996.htm "Stop 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ISiemensPPIStart 方法 |

启动西门子PLC为RUN模式，参数信息可以携带站号信息 "s=2;", 注意，分号是必须的。  
Start Siemens PLC in RUN mode, parameter information can carry station number information "s=2;", note that the semicolon is required.

**命名空间：**
 [HslCommunication.Profinet.Siemens.Helper](8dfc52a4-e8ce-f3d2-0647-4d18de8416d3.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
OperateResult Start(
	string parameter = ""
)
```

```
Function Start ( 
	Optional parameter As String = ""
) As OperateResult
```

```
OperateResult^ Start(
	String^ parameter = L""
)
```

```
abstract Start : 
        ?parameter : string 
(* Defaults:
        let _parameter = defaultArg parameter ""
*)
-> OperateResult 
```

#### 参数

parameter (Optional)
:   类型：SystemString  
    额外的参数信息，例如可以携带站号信息 "s=2;", 注意，分号是必须的。

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否启动成功

![](../icons/SectionExpanded.png)参见

#### 引用

[ISiemensPPI 接口](b9b7bddf-d584-6c46-fb23-670c72182056.htm)

[HslCommunication.Profinet.Siemens.Helper 命名空间](8dfc52a4-e8ce-f3d2-0647-4d18de8416d3.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Stop 方法 

[原文連結](http://api.hslcommunication.cn/html/d59f2c81-06f6-743c-2a09-b184b5623996.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Siemens.Helper](../html/8dfc52a4-e8ce-f3d2-0647-4d18de8416d3.htm "HslCommunication.Profinet.Siemens.Helper")

[ISiemensPPI 接口](../html/b9b7bddf-d584-6c46-fb23-670c72182056.htm "ISiemensPPI 接口")

[ISiemensPPI 方法](../html/c77ec8ca-32af-590e-a9e8-7cc7cd2ad723.htm "ISiemensPPI 方法")

[ReadPlcType 方法](../html/86b53070-da19-05b6-2b50-acab61ca51f1.htm "ReadPlcType 方法 ")

[Start 方法](../html/35128174-f97e-ad50-b559-71de1ac32caa.htm "Start 方法 ")

[Stop 方法](../html/d59f2c81-06f6-743c-2a09-b184b5623996.htm "Stop 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ISiemensPPIStop 方法 |

停止西门子PLC，切换为Stop模式，参数信息可以携带站号信息 "s=2;", 注意，分号是必须的。  
Stop Siemens PLC and switch to Stop mode, parameter information can carry station number information "s=2;", note that the semicolon is required.

**命名空间：**
 [HslCommunication.Profinet.Siemens.Helper](8dfc52a4-e8ce-f3d2-0647-4d18de8416d3.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
OperateResult Stop(
	string parameter = ""
)
```

```
Function Stop ( 
	Optional parameter As String = ""
) As OperateResult
```

```
OperateResult^ Stop(
	String^ parameter = L""
)
```

```
abstract Stop : 
        ?parameter : string 
(* Defaults:
        let _parameter = defaultArg parameter ""
*)
-> OperateResult 
```

#### 参数

parameter (Optional)
:   类型：SystemString  
    额外的参数信息，例如可以携带站号信息 "s=2;", 注意，分号是必须的。

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否停止成功

![](../icons/SectionExpanded.png)参见

#### 引用

[ISiemensPPI 接口](b9b7bddf-d584-6c46-fb23-670c72182056.htm)

[HslCommunication.Profinet.Siemens.Helper 命名空间](8dfc52a4-e8ce-f3d2-0647-4d18de8416d3.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## SiemensPPIHelper 类

[原文連結](http://api.hslcommunication.cn/html/24deecb2-7eb3-9aba-2da3-137adacdcbd9.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Siemens.Helper](../html/8dfc52a4-e8ce-f3d2-0647-4d18de8416d3.htm "HslCommunication.Profinet.Siemens.Helper")

[SiemensPPIHelper 类](../html/24deecb2-7eb3-9aba-2da3-137adacdcbd9.htm "SiemensPPIHelper 类")

[SiemensPPIHelper 构造函数](../html/5bdf4455-8da0-a3e1-f13f-b5984a4269c8.htm "SiemensPPIHelper 构造函数 ")

[SiemensPPIHelper 方法](../html/4b399f3b-1a4b-7dc1-3738-de0332b4df28.htm "SiemensPPIHelper 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| SiemensPPIHelper 类 |

西门子PPI协议的辅助类对象

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  HslCommunication.Profinet.Siemens.HelperSiemensPPIHelper

**命名空间：**
 [HslCommunication.Profinet.Siemens.Helper](8dfc52a4-e8ce-f3d2-0647-4d18de8416d3.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class SiemensPPIHelper
```

```
Public Class SiemensPPIHelper
```

```
public ref class SiemensPPIHelper
```

```
type SiemensPPIHelper =  class end
```

SiemensPPIHelper 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [SiemensPPIHelper](5bdf4455-8da0-a3e1-f13f-b5984a4269c8.htm) | 初始化 SiemensPPIHelper 类的一个新实例 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法静态成员 | [AnalysisAddress](c546b2ce-0787-8f2a-cb22-af4a20e2cc2e.htm) | 解析数据地址，解析出地址类型，起始地址，DB块的地址  Parse data address, parse out address type, start address, db block address |
| 公共方法静态成员 | [BuildReadCommand(Byte, S7AddressData, UInt16, Boolean)](b2261f10-e47c-1d30-dfc9-2a7f764f6cc9.htm) | 生成一个读取字数据指令头的通用方法  A general method for generating a command header to read a Word data |
| 公共方法静态成员 | [BuildReadCommand(Byte, String, UInt16, Boolean)](2fa2c432-d875-9b18-3751-cb7cbb6a307f.htm) | 生成一个读取字数据指令头的通用方法  A general method for generating a command header to read a Word data |
| 公共方法静态成员 | [BuildWriteCommand(Byte, String, Boolean)](4caeaf99-db0a-473a-eda9-7a45022b9196.htm) | 创建写入PLC的bool类型数据报文指令 |
| 公共方法静态成员 | [BuildWriteCommand(Byte, String, Byte)](3434783a-9022-a456-66ac-e6c1ea965c5a.htm) | 生成一个写入PLC数据信息的报文内容 |
| 公共方法静态成员 | [CheckReceiveDataComplete](8117ecc5-5e23-bfd4-4a64-2cf40d93daf3.htm) | 当消息头报文的长度定义为-1的时候，则使用动态的长度信息，可以使用本方法来判断一个消息是否处于完整的状态。  If the length of the message header is defined as -1, this method can be used to determine whether a message is in the complete state by using dynamic length information. |
| 公共方法静态成员 | [CheckResponse](4b18f602-f1bd-ee62-851a-3b68b37e83d8.htm) | 检查西门子PLC的返回的数据和合法性，对反馈的数据进行初步的校验 |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法静态成员 | [GetExecuteConfirm](12e4e1b2-3109-b243-f4d3-95cebad7fd95.htm) | 根据站号信息获取命令二次确认的报文信息 |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法静态成员 | [GetMsgFromStatus(Byte)](b0b45e44-eb56-30d8-87e6-bff63e4d8f10.htm) | 根据错误代号信息，获取到指定的文本信息  According to the error code information, get the specified text information |
| 公共方法静态成员 | [GetMsgFromStatus(Byte, Byte)](51965764-1bd1-e1b8-fc52-a1093f616d45.htm) | 根据错误信息，获取到文本信息 |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法静态成员 | [Read](779f6bef-32cb-56c1-ccd8-efabb01df7b6.htm) | 从西门子的PLC中读取数据信息，地址为"M100","AI100","I0","Q0","V100","S100"等  Read data information from Siemens PLC with addresses "M100", "AI100", "I0", "Q0", "V100", "S100", etc. |
| 公共方法静态成员 | [ReadBool(IReadWriteDevice, String, Byte, Object)](40320454-2b00-ba26-7783-85381ed2c28c.htm) | 从西门子的PLC中读取bool数据信息，地址为"M100.0","AI100.1","I0.3","Q0.6","V100.4","S100"等  Read bool data information from Siemens PLC, the addresses are "M100.0", "AI100.1", "I0.3", "Q0.6", "V100.4", "S100", etc. |
| 公共方法静态成员 | [ReadBool(IReadWriteDevice, String, UInt16, Byte, Object)](a00bc8a8-5dad-bd8b-0507-efd55bbd8e7a.htm) | 从西门子的PLC中读取bool数据信息，地址为"M100.0","AI100.1","I0.3","Q0.6","V100.4","S100"等  Read bool data information from Siemens PLC, the addresses are "M100.0", "AI100.1", "I0.3", "Q0.6", "V100.4", "S100", etc. |
| 公共方法静态成员 | [ReadPlcType](1c804c7a-7c9f-ca3c-96e0-5b8f327e1018.htm) | 读取西门子PLC的型号信息，参数信息可以携带站号信息 "s=2;", 注意，分号是必须的。  Read the model information of Siemens PLC, the parameter information can carry the station number information "s=2;", note that the semicolon is required. |
| 公共方法静态成员 | [Start](2c8531c5-1c74-4603-321d-cd96b90d899a.htm) | 启动西门子PLC为RUN模式，参数信息可以携带站号信息 "s=2;", 注意，分号是必须的。  Start Siemens PLC in RUN mode, parameter information can carry station number information "s=2;", note that the semicolon is required. |
| 公共方法静态成员 | [Stop](c36b817b-fa4d-df99-eb13-d0aa9e83a5df.htm) | 停止西门子PLC，切换为Stop模式，参数信息可以携带站号信息 "s=2;", 注意，分号是必须的。  Stop Siemens PLC and switch to Stop mode, parameter information can carry station number information "s=2;", note that the semicolon is required. |
| 公共方法 | ToString | (继承自 Object。) |
| 公共方法静态成员 | [Write(IReadWriteDevice, String, Boolean, Byte, Object)](7136caed-d967-caff-b692-67381753c19e.htm) | 将bool数据写入到西门子PLC中，地址为"M100.0","AI100.1","I0.3","Q0.6","V100.4","S100"等  Write the bool data to Siemens PLC with the addresses "M100.0", "AI100.1", "I0.3", "Q0.6", "V100.4", "S100", etc. |
| 公共方法静态成员 | [Write(IReadWriteDevice, String, Byte, Byte, Object)](093dc8fb-e7a9-3589-704b-39d4afdebdab.htm) | 将字节数据写入到西门子PLC中，地址为"M100.0","AI100.1","I0.3","Q0.6","V100.4","S100"等  Write byte data to Siemens PLC with addresses "M100.0", "AI100.1", "I0.3", "Q0.6", "V100.4", "S100", etc. |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.Profinet.Siemens.Helper 命名空间](8dfc52a4-e8ce-f3d2-0647-4d18de8416d3.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## SiemensPPIHelper 构造函数 

[原文連結](http://api.hslcommunication.cn/html/5bdf4455-8da0-a3e1-f13f-b5984a4269c8.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Siemens.Helper](../html/8dfc52a4-e8ce-f3d2-0647-4d18de8416d3.htm "HslCommunication.Profinet.Siemens.Helper")

[SiemensPPIHelper 类](../html/24deecb2-7eb3-9aba-2da3-137adacdcbd9.htm "SiemensPPIHelper 类")

[SiemensPPIHelper 构造函数](../html/5bdf4455-8da0-a3e1-f13f-b5984a4269c8.htm "SiemensPPIHelper 构造函数 ")

[SiemensPPIHelper 方法](../html/4b399f3b-1a4b-7dc1-3738-de0332b4df28.htm "SiemensPPIHelper 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| SiemensPPIHelper 构造函数 |

初始化 [SiemensPPIHelper](24deecb2-7eb3-9aba-2da3-137adacdcbd9.htm) 类的一个新实例

**命名空间：**
 [HslCommunication.Profinet.Siemens.Helper](8dfc52a4-e8ce-f3d2-0647-4d18de8416d3.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public SiemensPPIHelper()
```

```
Public Sub New
```

```
public:
SiemensPPIHelper()
```

```
new : unit -> SiemensPPIHelper
```

![](../icons/SectionExpanded.png)参见

#### 引用

[SiemensPPIHelper 类](24deecb2-7eb3-9aba-2da3-137adacdcbd9.htm)

[HslCommunication.Profinet.Siemens.Helper 命名空间](8dfc52a4-e8ce-f3d2-0647-4d18de8416d3.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## SiemensPPIHelper 方法

[原文連結](http://api.hslcommunication.cn/html/4b399f3b-1a4b-7dc1-3738-de0332b4df28.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Siemens.Helper](../html/8dfc52a4-e8ce-f3d2-0647-4d18de8416d3.htm "HslCommunication.Profinet.Siemens.Helper")

[SiemensPPIHelper 类](../html/24deecb2-7eb3-9aba-2da3-137adacdcbd9.htm "SiemensPPIHelper 类")

[SiemensPPIHelper 方法](../html/4b399f3b-1a4b-7dc1-3738-de0332b4df28.htm "SiemensPPIHelper 方法")

[AnalysisAddress 方法](../html/c546b2ce-0787-8f2a-cb22-af4a20e2cc2e.htm "AnalysisAddress 方法 ")

[BuildReadCommand 方法](../html/42b0fc83-2709-ba05-78dc-3a1c90da0bd6.htm "BuildReadCommand 方法 ")

[BuildWriteCommand 方法](../html/98fcc6b4-0da5-b298-7dcf-1b581a10b258.htm "BuildWriteCommand 方法 ")

[CheckReceiveDataComplete 方法](../html/8117ecc5-5e23-bfd4-4a64-2cf40d93daf3.htm "CheckReceiveDataComplete 方法 ")

[CheckResponse 方法](../html/4b18f602-f1bd-ee62-851a-3b68b37e83d8.htm "CheckResponse 方法 ")

[GetExecuteConfirm 方法](../html/12e4e1b2-3109-b243-f4d3-95cebad7fd95.htm "GetExecuteConfirm 方法 ")

[GetMsgFromStatus 方法](../html/5e092d96-9b66-8a86-c4a1-26fab2f0c209.htm "GetMsgFromStatus 方法 ")

[Read 方法](../html/779f6bef-32cb-56c1-ccd8-efabb01df7b6.htm "Read 方法 ")

[ReadBool 方法](../html/8f7c2d91-686b-c59c-54a8-cafee7433947.htm "ReadBool 方法 ")

[ReadPlcType 方法](../html/1c804c7a-7c9f-ca3c-96e0-5b8f327e1018.htm "ReadPlcType 方法 ")

[Start 方法](../html/2c8531c5-1c74-4603-321d-cd96b90d899a.htm "Start 方法 ")

[Stop 方法](../html/c36b817b-fa4d-df99-eb13-d0aa9e83a5df.htm "Stop 方法 ")

[Write 方法](../html/6ab7778a-7ce6-1601-a45d-fb2ad3eac080.htm "Write 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| SiemensPPIHelper 方法 |

[SiemensPPIHelper](24deecb2-7eb3-9aba-2da3-137adacdcbd9.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法静态成员 | [AnalysisAddress](c546b2ce-0787-8f2a-cb22-af4a20e2cc2e.htm) | 解析数据地址，解析出地址类型，起始地址，DB块的地址  Parse data address, parse out address type, start address, db block address |
| 公共方法静态成员 | [BuildReadCommand(Byte, S7AddressData, UInt16, Boolean)](b2261f10-e47c-1d30-dfc9-2a7f764f6cc9.htm) | 生成一个读取字数据指令头的通用方法  A general method for generating a command header to read a Word data |
| 公共方法静态成员 | [BuildReadCommand(Byte, String, UInt16, Boolean)](2fa2c432-d875-9b18-3751-cb7cbb6a307f.htm) | 生成一个读取字数据指令头的通用方法  A general method for generating a command header to read a Word data |
| 公共方法静态成员 | [BuildWriteCommand(Byte, String, Boolean)](4caeaf99-db0a-473a-eda9-7a45022b9196.htm) | 创建写入PLC的bool类型数据报文指令 |
| 公共方法静态成员 | [BuildWriteCommand(Byte, String, Byte)](3434783a-9022-a456-66ac-e6c1ea965c5a.htm) | 生成一个写入PLC数据信息的报文内容 |
| 公共方法静态成员 | [CheckReceiveDataComplete](8117ecc5-5e23-bfd4-4a64-2cf40d93daf3.htm) | 当消息头报文的长度定义为-1的时候，则使用动态的长度信息，可以使用本方法来判断一个消息是否处于完整的状态。  If the length of the message header is defined as -1, this method can be used to determine whether a message is in the complete state by using dynamic length information. |
| 公共方法静态成员 | [CheckResponse](4b18f602-f1bd-ee62-851a-3b68b37e83d8.htm) | 检查西门子PLC的返回的数据和合法性，对反馈的数据进行初步的校验 |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法静态成员 | [GetExecuteConfirm](12e4e1b2-3109-b243-f4d3-95cebad7fd95.htm) | 根据站号信息获取命令二次确认的报文信息 |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法静态成员 | [GetMsgFromStatus(Byte)](b0b45e44-eb56-30d8-87e6-bff63e4d8f10.htm) | 根据错误代号信息，获取到指定的文本信息  According to the error code information, get the specified text information |
| 公共方法静态成员 | [GetMsgFromStatus(Byte, Byte)](51965764-1bd1-e1b8-fc52-a1093f616d45.htm) | 根据错误信息，获取到文本信息 |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法静态成员 | [Read](779f6bef-32cb-56c1-ccd8-efabb01df7b6.htm) | 从西门子的PLC中读取数据信息，地址为"M100","AI100","I0","Q0","V100","S100"等  Read data information from Siemens PLC with addresses "M100", "AI100", "I0", "Q0", "V100", "S100", etc. |
| 公共方法静态成员 | [ReadBool(IReadWriteDevice, String, Byte, Object)](40320454-2b00-ba26-7783-85381ed2c28c.htm) | 从西门子的PLC中读取bool数据信息，地址为"M100.0","AI100.1","I0.3","Q0.6","V100.4","S100"等  Read bool data information from Siemens PLC, the addresses are "M100.0", "AI100.1", "I0.3", "Q0.6", "V100.4", "S100", etc. |
| 公共方法静态成员 | [ReadBool(IReadWriteDevice, String, UInt16, Byte, Object)](a00bc8a8-5dad-bd8b-0507-efd55bbd8e7a.htm) | 从西门子的PLC中读取bool数据信息，地址为"M100.0","AI100.1","I0.3","Q0.6","V100.4","S100"等  Read bool data information from Siemens PLC, the addresses are "M100.0", "AI100.1", "I0.3", "Q0.6", "V100.4", "S100", etc. |
| 公共方法静态成员 | [ReadPlcType](1c804c7a-7c9f-ca3c-96e0-5b8f327e1018.htm) | 读取西门子PLC的型号信息，参数信息可以携带站号信息 "s=2;", 注意，分号是必须的。  Read the model information of Siemens PLC, the parameter information can carry the station number information "s=2;", note that the semicolon is required. |
| 公共方法静态成员 | [Start](2c8531c5-1c74-4603-321d-cd96b90d899a.htm) | 启动西门子PLC为RUN模式，参数信息可以携带站号信息 "s=2;", 注意，分号是必须的。  Start Siemens PLC in RUN mode, parameter information can carry station number information "s=2;", note that the semicolon is required. |
| 公共方法静态成员 | [Stop](c36b817b-fa4d-df99-eb13-d0aa9e83a5df.htm) | 停止西门子PLC，切换为Stop模式，参数信息可以携带站号信息 "s=2;", 注意，分号是必须的。  Stop Siemens PLC and switch to Stop mode, parameter information can carry station number information "s=2;", note that the semicolon is required. |
| 公共方法 | ToString | (继承自 Object。) |
| 公共方法静态成员 | [Write(IReadWriteDevice, String, Boolean, Byte, Object)](7136caed-d967-caff-b692-67381753c19e.htm) | 将bool数据写入到西门子PLC中，地址为"M100.0","AI100.1","I0.3","Q0.6","V100.4","S100"等  Write the bool data to Siemens PLC with the addresses "M100.0", "AI100.1", "I0.3", "Q0.6", "V100.4", "S100", etc. |
| 公共方法静态成员 | [Write(IReadWriteDevice, String, Byte, Byte, Object)](093dc8fb-e7a9-3589-704b-39d4afdebdab.htm) | 将字节数据写入到西门子PLC中，地址为"M100.0","AI100.1","I0.3","Q0.6","V100.4","S100"等  Write byte data to Siemens PLC with addresses "M100.0", "AI100.1", "I0.3", "Q0.6", "V100.4", "S100", etc. |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[SiemensPPIHelper 类](24deecb2-7eb3-9aba-2da3-137adacdcbd9.htm)

[HslCommunication.Profinet.Siemens.Helper 命名空间](8dfc52a4-e8ce-f3d2-0647-4d18de8416d3.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## AnalysisAddress 方法 

[原文連結](http://api.hslcommunication.cn/html/c546b2ce-0787-8f2a-cb22-af4a20e2cc2e.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Siemens.Helper](../html/8dfc52a4-e8ce-f3d2-0647-4d18de8416d3.htm "HslCommunication.Profinet.Siemens.Helper")

[SiemensPPIHelper 类](../html/24deecb2-7eb3-9aba-2da3-137adacdcbd9.htm "SiemensPPIHelper 类")

[SiemensPPIHelper 方法](../html/4b399f3b-1a4b-7dc1-3738-de0332b4df28.htm "SiemensPPIHelper 方法")

[AnalysisAddress 方法](../html/c546b2ce-0787-8f2a-cb22-af4a20e2cc2e.htm "AnalysisAddress 方法 ")

[BuildReadCommand 方法](../html/42b0fc83-2709-ba05-78dc-3a1c90da0bd6.htm "BuildReadCommand 方法 ")

[BuildWriteCommand 方法](../html/98fcc6b4-0da5-b298-7dcf-1b581a10b258.htm "BuildWriteCommand 方法 ")

[CheckReceiveDataComplete 方法](../html/8117ecc5-5e23-bfd4-4a64-2cf40d93daf3.htm "CheckReceiveDataComplete 方法 ")

[CheckResponse 方法](../html/4b18f602-f1bd-ee62-851a-3b68b37e83d8.htm "CheckResponse 方法 ")

[GetExecuteConfirm 方法](../html/12e4e1b2-3109-b243-f4d3-95cebad7fd95.htm "GetExecuteConfirm 方法 ")

[GetMsgFromStatus 方法](../html/5e092d96-9b66-8a86-c4a1-26fab2f0c209.htm "GetMsgFromStatus 方法 ")

[Read 方法](../html/779f6bef-32cb-56c1-ccd8-efabb01df7b6.htm "Read 方法 ")

[ReadBool 方法](../html/8f7c2d91-686b-c59c-54a8-cafee7433947.htm "ReadBool 方法 ")

[ReadPlcType 方法](../html/1c804c7a-7c9f-ca3c-96e0-5b8f327e1018.htm "ReadPlcType 方法 ")

[Start 方法](../html/2c8531c5-1c74-4603-321d-cd96b90d899a.htm "Start 方法 ")

[Stop 方法](../html/c36b817b-fa4d-df99-eb13-d0aa9e83a5df.htm "Stop 方法 ")

[Write 方法](../html/6ab7778a-7ce6-1601-a45d-fb2ad3eac080.htm "Write 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| SiemensPPIHelperAnalysisAddress 方法 |

解析数据地址，解析出地址类型，起始地址，DB块的地址  
Parse data address, parse out address type, start address, db block address

**命名空间：**
 [HslCommunication.Profinet.Siemens.Helper](8dfc52a4-e8ce-f3d2-0647-4d18de8416d3.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult<S7AddressData> AnalysisAddress(
	string address
)
```

```
Public Shared Function AnalysisAddress ( 
	address As String
) As OperateResult(Of S7AddressData)
```

```
public:
static OperateResult<S7AddressData^>^ AnalysisAddress(
	String^ address
)
```

```
static member AnalysisAddress : 
        address : string -> OperateResult<S7AddressData> 
```

#### 参数

address
:   类型：SystemString  
    起始地址，例如M100，I0，Q0，V100 ->
    Start address, such as M100,I0,Q0,V100

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)[S7AddressData](1ca2fe47-3d14-8511-61cb-b2d4f79d938a.htm)  
解析数据地址，解析出地址类型，起始地址，DB块的地址 ->
Parse data address, parse out address type, start address, db block address

![](../icons/SectionExpanded.png)参见

#### 引用

[SiemensPPIHelper 类](24deecb2-7eb3-9aba-2da3-137adacdcbd9.htm)

[HslCommunication.Profinet.Siemens.Helper 命名空间](8dfc52a4-e8ce-f3d2-0647-4d18de8416d3.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## BuildReadCommand 方法 

[原文連結](http://api.hslcommunication.cn/html/42b0fc83-2709-ba05-78dc-3a1c90da0bd6.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Siemens.Helper](../html/8dfc52a4-e8ce-f3d2-0647-4d18de8416d3.htm "HslCommunication.Profinet.Siemens.Helper")

[SiemensPPIHelper 类](../html/24deecb2-7eb3-9aba-2da3-137adacdcbd9.htm "SiemensPPIHelper 类")

[SiemensPPIHelper 方法](../html/4b399f3b-1a4b-7dc1-3738-de0332b4df28.htm "SiemensPPIHelper 方法")

[BuildReadCommand 方法](../html/42b0fc83-2709-ba05-78dc-3a1c90da0bd6.htm "BuildReadCommand 方法 ")

[BuildReadCommand 方法 (Byte, S7AddressData, UInt16, Boolean)](../html/b2261f10-e47c-1d30-dfc9-2a7f764f6cc9.htm "BuildReadCommand 方法 (Byte, S7AddressData, UInt16, Boolean)")

[BuildReadCommand 方法 (Byte, String, UInt16, Boolean)](../html/2fa2c432-d875-9b18-3751-cb7cbb6a307f.htm "BuildReadCommand 方法 (Byte, String, UInt16, Boolean)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| SiemensPPIHelperBuildReadCommand 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法静态成员 | [BuildReadCommand(Byte, S7AddressData, UInt16, Boolean)](b2261f10-e47c-1d30-dfc9-2a7f764f6cc9.htm) | 生成一个读取字数据指令头的通用方法  A general method for generating a command header to read a Word data |
| 公共方法静态成员 | [BuildReadCommand(Byte, String, UInt16, Boolean)](2fa2c432-d875-9b18-3751-cb7cbb6a307f.htm) | 生成一个读取字数据指令头的通用方法  A general method for generating a command header to read a Word data |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[SiemensPPIHelper 类](24deecb2-7eb3-9aba-2da3-137adacdcbd9.htm)

[HslCommunication.Profinet.Siemens.Helper 命名空间](8dfc52a4-e8ce-f3d2-0647-4d18de8416d3.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## BuildReadCommand 方法 (Byte, S7AddressData, UInt16, Boolean)

[原文連結](http://api.hslcommunication.cn/html/b2261f10-e47c-1d30-dfc9-2a7f764f6cc9.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Siemens.Helper](../html/8dfc52a4-e8ce-f3d2-0647-4d18de8416d3.htm "HslCommunication.Profinet.Siemens.Helper")

[SiemensPPIHelper 类](../html/24deecb2-7eb3-9aba-2da3-137adacdcbd9.htm "SiemensPPIHelper 类")

[SiemensPPIHelper 方法](../html/4b399f3b-1a4b-7dc1-3738-de0332b4df28.htm "SiemensPPIHelper 方法")

[BuildReadCommand 方法](../html/42b0fc83-2709-ba05-78dc-3a1c90da0bd6.htm "BuildReadCommand 方法 ")

[BuildReadCommand 方法 (Byte, S7AddressData, UInt16, Boolean)](../html/b2261f10-e47c-1d30-dfc9-2a7f764f6cc9.htm "BuildReadCommand 方法 (Byte, S7AddressData, UInt16, Boolean)")

[BuildReadCommand 方法 (Byte, String, UInt16, Boolean)](../html/2fa2c432-d875-9b18-3751-cb7cbb6a307f.htm "BuildReadCommand 方法 (Byte, String, UInt16, Boolean)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| SiemensPPIHelperBuildReadCommand 方法 (Byte, S7AddressData, UInt16, Boolean) |

生成一个读取字数据指令头的通用方法  
A general method for generating a command header to read a Word data

**命名空间：**
 [HslCommunication.Profinet.Siemens.Helper](8dfc52a4-e8ce-f3d2-0647-4d18de8416d3.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult<byte[]> BuildReadCommand(
	byte station,
	S7AddressData address,
	ushort length,
	bool isBit
)
```

```
Public Shared Function BuildReadCommand ( 
	station As Byte,
	address As S7AddressData,
	length As UShort,
	isBit As Boolean
) As OperateResult(Of Byte())
```

```
public:
static OperateResult<array<unsigned char>^>^ BuildReadCommand(
	unsigned char station, 
	S7AddressData^ address, 
	unsigned short length, 
	bool isBit
)
```

```
static member BuildReadCommand : 
        station : byte * 
        address : S7AddressData * 
        length : uint16 * 
        isBit : bool -> OperateResult<byte[]> 
```

#### 参数

station
:   类型：SystemByte  
    设备的站号信息 -> Station number information for the device

address
:   类型：[HslCommunication.Core.AddressS7AddressData](1ca2fe47-3d14-8511-61cb-b2d4f79d938a.htm)  
    起始地址，例如M100，I0，Q0，V100 ->
    Start address, such as M100,I0,Q0,V100

length
:   类型：SystemUInt16  
    读取数据长度 -> Read Data length

isBit
:   类型：SystemBoolean  
    是否为位读取

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
包含结果对象的报文 -> Message containing the result object

![](../icons/SectionExpanded.png)参见

#### 引用

[SiemensPPIHelper 类](24deecb2-7eb3-9aba-2da3-137adacdcbd9.htm)

[BuildReadCommand 重载](42b0fc83-2709-ba05-78dc-3a1c90da0bd6.htm)

[HslCommunication.Profinet.Siemens.Helper 命名空间](8dfc52a4-e8ce-f3d2-0647-4d18de8416d3.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## BuildReadCommand 方法 (Byte, String, UInt16, Boolean)

[原文連結](http://api.hslcommunication.cn/html/2fa2c432-d875-9b18-3751-cb7cbb6a307f.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Siemens.Helper](../html/8dfc52a4-e8ce-f3d2-0647-4d18de8416d3.htm "HslCommunication.Profinet.Siemens.Helper")

[SiemensPPIHelper 类](../html/24deecb2-7eb3-9aba-2da3-137adacdcbd9.htm "SiemensPPIHelper 类")

[SiemensPPIHelper 方法](../html/4b399f3b-1a4b-7dc1-3738-de0332b4df28.htm "SiemensPPIHelper 方法")

[BuildReadCommand 方法](../html/42b0fc83-2709-ba05-78dc-3a1c90da0bd6.htm "BuildReadCommand 方法 ")

[BuildReadCommand 方法 (Byte, S7AddressData, UInt16, Boolean)](../html/b2261f10-e47c-1d30-dfc9-2a7f764f6cc9.htm "BuildReadCommand 方法 (Byte, S7AddressData, UInt16, Boolean)")

[BuildReadCommand 方法 (Byte, String, UInt16, Boolean)](../html/2fa2c432-d875-9b18-3751-cb7cbb6a307f.htm "BuildReadCommand 方法 (Byte, String, UInt16, Boolean)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| SiemensPPIHelperBuildReadCommand 方法 (Byte, String, UInt16, Boolean) |

生成一个读取字数据指令头的通用方法  
A general method for generating a command header to read a Word data

**命名空间：**
 [HslCommunication.Profinet.Siemens.Helper](8dfc52a4-e8ce-f3d2-0647-4d18de8416d3.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult<byte[]> BuildReadCommand(
	byte station,
	string address,
	ushort length,
	bool isBit
)
```

```
Public Shared Function BuildReadCommand ( 
	station As Byte,
	address As String,
	length As UShort,
	isBit As Boolean
) As OperateResult(Of Byte())
```

```
public:
static OperateResult<array<unsigned char>^>^ BuildReadCommand(
	unsigned char station, 
	String^ address, 
	unsigned short length, 
	bool isBit
)
```

```
static member BuildReadCommand : 
        station : byte * 
        address : string * 
        length : uint16 * 
        isBit : bool -> OperateResult<byte[]> 
```

#### 参数

station
:   类型：SystemByte  
    设备的站号信息 -> Station number information for the device

address
:   类型：SystemString  
    起始地址，例如M100，I0，Q0，V100 ->
    Start address, such as M100,I0,Q0,V100

length
:   类型：SystemUInt16  
    读取数据长度 -> Read Data length

isBit
:   类型：SystemBoolean  
    是否为位读取

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
包含结果对象的报文 -> Message containing the result object

![](../icons/SectionExpanded.png)参见

#### 引用

[SiemensPPIHelper 类](24deecb2-7eb3-9aba-2da3-137adacdcbd9.htm)

[BuildReadCommand 重载](42b0fc83-2709-ba05-78dc-3a1c90da0bd6.htm)

[HslCommunication.Profinet.Siemens.Helper 命名空间](8dfc52a4-e8ce-f3d2-0647-4d18de8416d3.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## BuildWriteCommand 方法 

[原文連結](http://api.hslcommunication.cn/html/98fcc6b4-0da5-b298-7dcf-1b581a10b258.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Siemens.Helper](../html/8dfc52a4-e8ce-f3d2-0647-4d18de8416d3.htm "HslCommunication.Profinet.Siemens.Helper")

[SiemensPPIHelper 类](../html/24deecb2-7eb3-9aba-2da3-137adacdcbd9.htm "SiemensPPIHelper 类")

[SiemensPPIHelper 方法](../html/4b399f3b-1a4b-7dc1-3738-de0332b4df28.htm "SiemensPPIHelper 方法")

[BuildWriteCommand 方法](../html/98fcc6b4-0da5-b298-7dcf-1b581a10b258.htm "BuildWriteCommand 方法 ")

[BuildWriteCommand 方法 (Byte, String, Boolean[])](../html/4caeaf99-db0a-473a-eda9-7a45022b9196.htm "BuildWriteCommand 方法 (Byte, String, Boolean[])")

[BuildWriteCommand 方法 (Byte, String, Byte[])](../html/3434783a-9022-a456-66ac-e6c1ea965c5a.htm "BuildWriteCommand 方法 (Byte, String, Byte[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| SiemensPPIHelperBuildWriteCommand 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法静态成员 | [BuildWriteCommand(Byte, String, Boolean)](4caeaf99-db0a-473a-eda9-7a45022b9196.htm) | 创建写入PLC的bool类型数据报文指令 |
| 公共方法静态成员 | [BuildWriteCommand(Byte, String, Byte)](3434783a-9022-a456-66ac-e6c1ea965c5a.htm) | 生成一个写入PLC数据信息的报文内容 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[SiemensPPIHelper 类](24deecb2-7eb3-9aba-2da3-137adacdcbd9.htm)

[HslCommunication.Profinet.Siemens.Helper 命名空间](8dfc52a4-e8ce-f3d2-0647-4d18de8416d3.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## BuildWriteCommand 方法 (Byte, String, Boolean[])

[原文連結](http://api.hslcommunication.cn/html/4caeaf99-db0a-473a-eda9-7a45022b9196.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Siemens.Helper](../html/8dfc52a4-e8ce-f3d2-0647-4d18de8416d3.htm "HslCommunication.Profinet.Siemens.Helper")

[SiemensPPIHelper 类](../html/24deecb2-7eb3-9aba-2da3-137adacdcbd9.htm "SiemensPPIHelper 类")

[SiemensPPIHelper 方法](../html/4b399f3b-1a4b-7dc1-3738-de0332b4df28.htm "SiemensPPIHelper 方法")

[BuildWriteCommand 方法](../html/98fcc6b4-0da5-b298-7dcf-1b581a10b258.htm "BuildWriteCommand 方法 ")

[BuildWriteCommand 方法 (Byte, String, Boolean[])](../html/4caeaf99-db0a-473a-eda9-7a45022b9196.htm "BuildWriteCommand 方法 (Byte, String, Boolean[])")

[BuildWriteCommand 方法 (Byte, String, Byte[])](../html/3434783a-9022-a456-66ac-e6c1ea965c5a.htm "BuildWriteCommand 方法 (Byte, String, Byte[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| SiemensPPIHelperBuildWriteCommand 方法 (Byte, String, Boolean) |

创建写入PLC的bool类型数据报文指令

**命名空间：**
 [HslCommunication.Profinet.Siemens.Helper](8dfc52a4-e8ce-f3d2-0647-4d18de8416d3.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult<byte[]> BuildWriteCommand(
	byte station,
	string address,
	bool[] values
)
```

```
Public Shared Function BuildWriteCommand ( 
	station As Byte,
	address As String,
	values As Boolean()
) As OperateResult(Of Byte())
```

```
public:
static OperateResult<array<unsigned char>^>^ BuildWriteCommand(
	unsigned char station, 
	String^ address, 
	array<bool>^ values
)
```

```
static member BuildWriteCommand : 
        station : byte * 
        address : string * 
        values : bool[] -> OperateResult<byte[]> 
```

#### 参数

station
:   类型：SystemByte  
    PLC的站号信息

address
:   类型：SystemString  
    地址信息

values
:   类型：SystemBoolean  
    bool[]数据值

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
带有成功标识的结果对象

![](../icons/SectionExpanded.png)参见

#### 引用

[SiemensPPIHelper 类](24deecb2-7eb3-9aba-2da3-137adacdcbd9.htm)

[BuildWriteCommand 重载](98fcc6b4-0da5-b298-7dcf-1b581a10b258.htm)

[HslCommunication.Profinet.Siemens.Helper 命名空间](8dfc52a4-e8ce-f3d2-0647-4d18de8416d3.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## BuildWriteCommand 方法 (Byte, String, Byte[])

[原文連結](http://api.hslcommunication.cn/html/3434783a-9022-a456-66ac-e6c1ea965c5a.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Siemens.Helper](../html/8dfc52a4-e8ce-f3d2-0647-4d18de8416d3.htm "HslCommunication.Profinet.Siemens.Helper")

[SiemensPPIHelper 类](../html/24deecb2-7eb3-9aba-2da3-137adacdcbd9.htm "SiemensPPIHelper 类")

[SiemensPPIHelper 方法](../html/4b399f3b-1a4b-7dc1-3738-de0332b4df28.htm "SiemensPPIHelper 方法")

[BuildWriteCommand 方法](../html/98fcc6b4-0da5-b298-7dcf-1b581a10b258.htm "BuildWriteCommand 方法 ")

[BuildWriteCommand 方法 (Byte, String, Boolean[])](../html/4caeaf99-db0a-473a-eda9-7a45022b9196.htm "BuildWriteCommand 方法 (Byte, String, Boolean[])")

[BuildWriteCommand 方法 (Byte, String, Byte[])](../html/3434783a-9022-a456-66ac-e6c1ea965c5a.htm "BuildWriteCommand 方法 (Byte, String, Byte[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| SiemensPPIHelperBuildWriteCommand 方法 (Byte, String, Byte) |

生成一个写入PLC数据信息的报文内容

**命名空间：**
 [HslCommunication.Profinet.Siemens.Helper](8dfc52a4-e8ce-f3d2-0647-4d18de8416d3.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult<byte[]> BuildWriteCommand(
	byte station,
	string address,
	byte[] values
)
```

```
Public Shared Function BuildWriteCommand ( 
	station As Byte,
	address As String,
	values As Byte()
) As OperateResult(Of Byte())
```

```
public:
static OperateResult<array<unsigned char>^>^ BuildWriteCommand(
	unsigned char station, 
	String^ address, 
	array<unsigned char>^ values
)
```

```
static member BuildWriteCommand : 
        station : byte * 
        address : string * 
        values : byte[] -> OperateResult<byte[]> 
```

#### 参数

station
:   类型：SystemByte  
    PLC的站号

address
:   类型：SystemString  
    地址

values
:   类型：SystemByte  
    数据值

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
是否写入成功

![](../icons/SectionExpanded.png)参见

#### 引用

[SiemensPPIHelper 类](24deecb2-7eb3-9aba-2da3-137adacdcbd9.htm)

[BuildWriteCommand 重载](98fcc6b4-0da5-b298-7dcf-1b581a10b258.htm)

[HslCommunication.Profinet.Siemens.Helper 命名空间](8dfc52a4-e8ce-f3d2-0647-4d18de8416d3.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## CheckReceiveDataComplete 方法 

[原文連結](http://api.hslcommunication.cn/html/8117ecc5-5e23-bfd4-4a64-2cf40d93daf3.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Siemens.Helper](../html/8dfc52a4-e8ce-f3d2-0647-4d18de8416d3.htm "HslCommunication.Profinet.Siemens.Helper")

[SiemensPPIHelper 类](../html/24deecb2-7eb3-9aba-2da3-137adacdcbd9.htm "SiemensPPIHelper 类")

[SiemensPPIHelper 方法](../html/4b399f3b-1a4b-7dc1-3738-de0332b4df28.htm "SiemensPPIHelper 方法")

[AnalysisAddress 方法](../html/c546b2ce-0787-8f2a-cb22-af4a20e2cc2e.htm "AnalysisAddress 方法 ")

[BuildReadCommand 方法](../html/42b0fc83-2709-ba05-78dc-3a1c90da0bd6.htm "BuildReadCommand 方法 ")

[BuildWriteCommand 方法](../html/98fcc6b4-0da5-b298-7dcf-1b581a10b258.htm "BuildWriteCommand 方法 ")

[CheckReceiveDataComplete 方法](../html/8117ecc5-5e23-bfd4-4a64-2cf40d93daf3.htm "CheckReceiveDataComplete 方法 ")

[CheckResponse 方法](../html/4b18f602-f1bd-ee62-851a-3b68b37e83d8.htm "CheckResponse 方法 ")

[GetExecuteConfirm 方法](../html/12e4e1b2-3109-b243-f4d3-95cebad7fd95.htm "GetExecuteConfirm 方法 ")

[GetMsgFromStatus 方法](../html/5e092d96-9b66-8a86-c4a1-26fab2f0c209.htm "GetMsgFromStatus 方法 ")

[Read 方法](../html/779f6bef-32cb-56c1-ccd8-efabb01df7b6.htm "Read 方法 ")

[ReadBool 方法](../html/8f7c2d91-686b-c59c-54a8-cafee7433947.htm "ReadBool 方法 ")

[ReadPlcType 方法](../html/1c804c7a-7c9f-ca3c-96e0-5b8f327e1018.htm "ReadPlcType 方法 ")

[Start 方法](../html/2c8531c5-1c74-4603-321d-cd96b90d899a.htm "Start 方法 ")

[Stop 方法](../html/c36b817b-fa4d-df99-eb13-d0aa9e83a5df.htm "Stop 方法 ")

[Write 方法](../html/6ab7778a-7ce6-1601-a45d-fb2ad3eac080.htm "Write 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| SiemensPPIHelperCheckReceiveDataComplete 方法 |

当消息头报文的长度定义为-1的时候，则使用动态的长度信息，可以使用本方法来判断一个消息是否处于完整的状态。  
If the length of the message header is defined as -1, this method can be used to determine whether a message is in the complete state by using dynamic length information.

**命名空间：**
 [HslCommunication.Profinet.Siemens.Helper](8dfc52a4-e8ce-f3d2-0647-4d18de8416d3.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static bool CheckReceiveDataComplete(
	MemoryStream ms
)
```

```
Public Shared Function CheckReceiveDataComplete ( 
	ms As MemoryStream
) As Boolean
```

```
public:
static bool CheckReceiveDataComplete(
	MemoryStream^ ms
)
```

```
static member CheckReceiveDataComplete : 
        ms : MemoryStream -> bool 
```

#### 参数

ms
:   类型：System.IOMemoryStream  
    接收到的数据内容

#### 返回值

类型：Boolean  
是否是完整的消息

![](../icons/SectionExpanded.png)参见

#### 引用

[SiemensPPIHelper 类](24deecb2-7eb3-9aba-2da3-137adacdcbd9.htm)

[HslCommunication.Profinet.Siemens.Helper 命名空间](8dfc52a4-e8ce-f3d2-0647-4d18de8416d3.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## CheckResponse 方法 

[原文連結](http://api.hslcommunication.cn/html/4b18f602-f1bd-ee62-851a-3b68b37e83d8.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Siemens.Helper](../html/8dfc52a4-e8ce-f3d2-0647-4d18de8416d3.htm "HslCommunication.Profinet.Siemens.Helper")

[SiemensPPIHelper 类](../html/24deecb2-7eb3-9aba-2da3-137adacdcbd9.htm "SiemensPPIHelper 类")

[SiemensPPIHelper 方法](../html/4b399f3b-1a4b-7dc1-3738-de0332b4df28.htm "SiemensPPIHelper 方法")

[AnalysisAddress 方法](../html/c546b2ce-0787-8f2a-cb22-af4a20e2cc2e.htm "AnalysisAddress 方法 ")

[BuildReadCommand 方法](../html/42b0fc83-2709-ba05-78dc-3a1c90da0bd6.htm "BuildReadCommand 方法 ")

[BuildWriteCommand 方法](../html/98fcc6b4-0da5-b298-7dcf-1b581a10b258.htm "BuildWriteCommand 方法 ")

[CheckReceiveDataComplete 方法](../html/8117ecc5-5e23-bfd4-4a64-2cf40d93daf3.htm "CheckReceiveDataComplete 方法 ")

[CheckResponse 方法](../html/4b18f602-f1bd-ee62-851a-3b68b37e83d8.htm "CheckResponse 方法 ")

[GetExecuteConfirm 方法](../html/12e4e1b2-3109-b243-f4d3-95cebad7fd95.htm "GetExecuteConfirm 方法 ")

[GetMsgFromStatus 方法](../html/5e092d96-9b66-8a86-c4a1-26fab2f0c209.htm "GetMsgFromStatus 方法 ")

[Read 方法](../html/779f6bef-32cb-56c1-ccd8-efabb01df7b6.htm "Read 方法 ")

[ReadBool 方法](../html/8f7c2d91-686b-c59c-54a8-cafee7433947.htm "ReadBool 方法 ")

[ReadPlcType 方法](../html/1c804c7a-7c9f-ca3c-96e0-5b8f327e1018.htm "ReadPlcType 方法 ")

[Start 方法](../html/2c8531c5-1c74-4603-321d-cd96b90d899a.htm "Start 方法 ")

[Stop 方法](../html/c36b817b-fa4d-df99-eb13-d0aa9e83a5df.htm "Stop 方法 ")

[Write 方法](../html/6ab7778a-7ce6-1601-a45d-fb2ad3eac080.htm "Write 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| SiemensPPIHelperCheckResponse 方法 |

检查西门子PLC的返回的数据和合法性，对反馈的数据进行初步的校验

**命名空间：**
 [HslCommunication.Profinet.Siemens.Helper](8dfc52a4-e8ce-f3d2-0647-4d18de8416d3.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult CheckResponse(
	byte[] content
)
```

```
Public Shared Function CheckResponse ( 
	content As Byte()
) As OperateResult
```

```
public:
static OperateResult^ CheckResponse(
	array<unsigned char>^ content
)
```

```
static member CheckResponse : 
        content : byte[] -> OperateResult 
```

#### 参数

content
:   类型：SystemByte  
    服务器返回的原始的数据内容

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否校验成功

![](../icons/SectionExpanded.png)参见

#### 引用

[SiemensPPIHelper 类](24deecb2-7eb3-9aba-2da3-137adacdcbd9.htm)

[HslCommunication.Profinet.Siemens.Helper 命名空间](8dfc52a4-e8ce-f3d2-0647-4d18de8416d3.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GetExecuteConfirm 方法 

[原文連結](http://api.hslcommunication.cn/html/12e4e1b2-3109-b243-f4d3-95cebad7fd95.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Siemens.Helper](../html/8dfc52a4-e8ce-f3d2-0647-4d18de8416d3.htm "HslCommunication.Profinet.Siemens.Helper")

[SiemensPPIHelper 类](../html/24deecb2-7eb3-9aba-2da3-137adacdcbd9.htm "SiemensPPIHelper 类")

[SiemensPPIHelper 方法](../html/4b399f3b-1a4b-7dc1-3738-de0332b4df28.htm "SiemensPPIHelper 方法")

[AnalysisAddress 方法](../html/c546b2ce-0787-8f2a-cb22-af4a20e2cc2e.htm "AnalysisAddress 方法 ")

[BuildReadCommand 方法](../html/42b0fc83-2709-ba05-78dc-3a1c90da0bd6.htm "BuildReadCommand 方法 ")

[BuildWriteCommand 方法](../html/98fcc6b4-0da5-b298-7dcf-1b581a10b258.htm "BuildWriteCommand 方法 ")

[CheckReceiveDataComplete 方法](../html/8117ecc5-5e23-bfd4-4a64-2cf40d93daf3.htm "CheckReceiveDataComplete 方法 ")

[CheckResponse 方法](../html/4b18f602-f1bd-ee62-851a-3b68b37e83d8.htm "CheckResponse 方法 ")

[GetExecuteConfirm 方法](../html/12e4e1b2-3109-b243-f4d3-95cebad7fd95.htm "GetExecuteConfirm 方法 ")

[GetMsgFromStatus 方法](../html/5e092d96-9b66-8a86-c4a1-26fab2f0c209.htm "GetMsgFromStatus 方法 ")

[Read 方法](../html/779f6bef-32cb-56c1-ccd8-efabb01df7b6.htm "Read 方法 ")

[ReadBool 方法](../html/8f7c2d91-686b-c59c-54a8-cafee7433947.htm "ReadBool 方法 ")

[ReadPlcType 方法](../html/1c804c7a-7c9f-ca3c-96e0-5b8f327e1018.htm "ReadPlcType 方法 ")

[Start 方法](../html/2c8531c5-1c74-4603-321d-cd96b90d899a.htm "Start 方法 ")

[Stop 方法](../html/c36b817b-fa4d-df99-eb13-d0aa9e83a5df.htm "Stop 方法 ")

[Write 方法](../html/6ab7778a-7ce6-1601-a45d-fb2ad3eac080.htm "Write 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| SiemensPPIHelperGetExecuteConfirm 方法 |

根据站号信息获取命令二次确认的报文信息

**命名空间：**
 [HslCommunication.Profinet.Siemens.Helper](8dfc52a4-e8ce-f3d2-0647-4d18de8416d3.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static byte[] GetExecuteConfirm(
	byte station
)
```

```
Public Shared Function GetExecuteConfirm ( 
	station As Byte
) As Byte()
```

```
public:
static array<unsigned char>^ GetExecuteConfirm(
	unsigned char station
)
```

```
static member GetExecuteConfirm : 
        station : byte -> byte[] 
```

#### 参数

station
:   类型：SystemByte  
    站号信息

#### 返回值

类型：Byte  
二次命令确认的报文

![](../icons/SectionExpanded.png)参见

#### 引用

[SiemensPPIHelper 类](24deecb2-7eb3-9aba-2da3-137adacdcbd9.htm)

[HslCommunication.Profinet.Siemens.Helper 命名空间](8dfc52a4-e8ce-f3d2-0647-4d18de8416d3.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GetMsgFromStatus 方法 

[原文連結](http://api.hslcommunication.cn/html/5e092d96-9b66-8a86-c4a1-26fab2f0c209.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Siemens.Helper](../html/8dfc52a4-e8ce-f3d2-0647-4d18de8416d3.htm "HslCommunication.Profinet.Siemens.Helper")

[SiemensPPIHelper 类](../html/24deecb2-7eb3-9aba-2da3-137adacdcbd9.htm "SiemensPPIHelper 类")

[SiemensPPIHelper 方法](../html/4b399f3b-1a4b-7dc1-3738-de0332b4df28.htm "SiemensPPIHelper 方法")

[GetMsgFromStatus 方法](../html/5e092d96-9b66-8a86-c4a1-26fab2f0c209.htm "GetMsgFromStatus 方法 ")

[GetMsgFromStatus 方法 (Byte)](../html/b0b45e44-eb56-30d8-87e6-bff63e4d8f10.htm "GetMsgFromStatus 方法 (Byte)")

[GetMsgFromStatus 方法 (Byte, Byte)](../html/51965764-1bd1-e1b8-fc52-a1093f616d45.htm "GetMsgFromStatus 方法 (Byte, Byte)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| SiemensPPIHelperGetMsgFromStatus 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法静态成员 | [GetMsgFromStatus(Byte)](b0b45e44-eb56-30d8-87e6-bff63e4d8f10.htm) | 根据错误代号信息，获取到指定的文本信息  According to the error code information, get the specified text information |
| 公共方法静态成员 | [GetMsgFromStatus(Byte, Byte)](51965764-1bd1-e1b8-fc52-a1093f616d45.htm) | 根据错误信息，获取到文本信息 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[SiemensPPIHelper 类](24deecb2-7eb3-9aba-2da3-137adacdcbd9.htm)

[HslCommunication.Profinet.Siemens.Helper 命名空间](8dfc52a4-e8ce-f3d2-0647-4d18de8416d3.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GetMsgFromStatus 方法 (Byte)

[原文連結](http://api.hslcommunication.cn/html/b0b45e44-eb56-30d8-87e6-bff63e4d8f10.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Siemens.Helper](../html/8dfc52a4-e8ce-f3d2-0647-4d18de8416d3.htm "HslCommunication.Profinet.Siemens.Helper")

[SiemensPPIHelper 类](../html/24deecb2-7eb3-9aba-2da3-137adacdcbd9.htm "SiemensPPIHelper 类")

[SiemensPPIHelper 方法](../html/4b399f3b-1a4b-7dc1-3738-de0332b4df28.htm "SiemensPPIHelper 方法")

[GetMsgFromStatus 方法](../html/5e092d96-9b66-8a86-c4a1-26fab2f0c209.htm "GetMsgFromStatus 方法 ")

[GetMsgFromStatus 方法 (Byte)](../html/b0b45e44-eb56-30d8-87e6-bff63e4d8f10.htm "GetMsgFromStatus 方法 (Byte)")

[GetMsgFromStatus 方法 (Byte, Byte)](../html/51965764-1bd1-e1b8-fc52-a1093f616d45.htm "GetMsgFromStatus 方法 (Byte, Byte)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| SiemensPPIHelperGetMsgFromStatus 方法 (Byte) |

根据错误代号信息，获取到指定的文本信息  
According to the error code information, get the specified text information

**命名空间：**
 [HslCommunication.Profinet.Siemens.Helper](8dfc52a4-e8ce-f3d2-0647-4d18de8416d3.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static string GetMsgFromStatus(
	byte code
)
```

```
Public Shared Function GetMsgFromStatus ( 
	code As Byte
) As String
```

```
public:
static String^ GetMsgFromStatus(
	unsigned char code
)
```

```
static member GetMsgFromStatus : 
        code : byte -> string 
```

#### 参数

code
:   类型：SystemByte  
    错误状态信息

#### 返回值

类型：String  
消息文本

![](../icons/SectionExpanded.png)参见

#### 引用

[SiemensPPIHelper 类](24deecb2-7eb3-9aba-2da3-137adacdcbd9.htm)

[GetMsgFromStatus 重载](5e092d96-9b66-8a86-c4a1-26fab2f0c209.htm)

[HslCommunication.Profinet.Siemens.Helper 命名空间](8dfc52a4-e8ce-f3d2-0647-4d18de8416d3.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GetMsgFromStatus 方法 (Byte, Byte)

[原文連結](http://api.hslcommunication.cn/html/51965764-1bd1-e1b8-fc52-a1093f616d45.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Siemens.Helper](../html/8dfc52a4-e8ce-f3d2-0647-4d18de8416d3.htm "HslCommunication.Profinet.Siemens.Helper")

[SiemensPPIHelper 类](../html/24deecb2-7eb3-9aba-2da3-137adacdcbd9.htm "SiemensPPIHelper 类")

[SiemensPPIHelper 方法](../html/4b399f3b-1a4b-7dc1-3738-de0332b4df28.htm "SiemensPPIHelper 方法")

[GetMsgFromStatus 方法](../html/5e092d96-9b66-8a86-c4a1-26fab2f0c209.htm "GetMsgFromStatus 方法 ")

[GetMsgFromStatus 方法 (Byte)](../html/b0b45e44-eb56-30d8-87e6-bff63e4d8f10.htm "GetMsgFromStatus 方法 (Byte)")

[GetMsgFromStatus 方法 (Byte, Byte)](../html/51965764-1bd1-e1b8-fc52-a1093f616d45.htm "GetMsgFromStatus 方法 (Byte, Byte)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| SiemensPPIHelperGetMsgFromStatus 方法 (Byte, Byte) |

根据错误信息，获取到文本信息

**命名空间：**
 [HslCommunication.Profinet.Siemens.Helper](8dfc52a4-e8ce-f3d2-0647-4d18de8416d3.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static string GetMsgFromStatus(
	byte errorClass,
	byte errorCode
)
```

```
Public Shared Function GetMsgFromStatus ( 
	errorClass As Byte,
	errorCode As Byte
) As String
```

```
public:
static String^ GetMsgFromStatus(
	unsigned char errorClass, 
	unsigned char errorCode
)
```

```
static member GetMsgFromStatus : 
        errorClass : byte * 
        errorCode : byte -> string 
```

#### 参数

errorClass
:   类型：SystemByte  
    错误类型

errorCode
:   类型：SystemByte  
    错误代码

#### 返回值

类型：String  
错误信息

![](../icons/SectionExpanded.png)参见

#### 引用

[SiemensPPIHelper 类](24deecb2-7eb3-9aba-2da3-137adacdcbd9.htm)

[GetMsgFromStatus 重载](5e092d96-9b66-8a86-c4a1-26fab2f0c209.htm)

[HslCommunication.Profinet.Siemens.Helper 命名空间](8dfc52a4-e8ce-f3d2-0647-4d18de8416d3.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Read 方法 

[原文連結](http://api.hslcommunication.cn/html/779f6bef-32cb-56c1-ccd8-efabb01df7b6.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Siemens.Helper](../html/8dfc52a4-e8ce-f3d2-0647-4d18de8416d3.htm "HslCommunication.Profinet.Siemens.Helper")

[SiemensPPIHelper 类](../html/24deecb2-7eb3-9aba-2da3-137adacdcbd9.htm "SiemensPPIHelper 类")

[SiemensPPIHelper 方法](../html/4b399f3b-1a4b-7dc1-3738-de0332b4df28.htm "SiemensPPIHelper 方法")

[AnalysisAddress 方法](../html/c546b2ce-0787-8f2a-cb22-af4a20e2cc2e.htm "AnalysisAddress 方法 ")

[BuildReadCommand 方法](../html/42b0fc83-2709-ba05-78dc-3a1c90da0bd6.htm "BuildReadCommand 方法 ")

[BuildWriteCommand 方法](../html/98fcc6b4-0da5-b298-7dcf-1b581a10b258.htm "BuildWriteCommand 方法 ")

[CheckReceiveDataComplete 方法](../html/8117ecc5-5e23-bfd4-4a64-2cf40d93daf3.htm "CheckReceiveDataComplete 方法 ")

[CheckResponse 方法](../html/4b18f602-f1bd-ee62-851a-3b68b37e83d8.htm "CheckResponse 方法 ")

[GetExecuteConfirm 方法](../html/12e4e1b2-3109-b243-f4d3-95cebad7fd95.htm "GetExecuteConfirm 方法 ")

[GetMsgFromStatus 方法](../html/5e092d96-9b66-8a86-c4a1-26fab2f0c209.htm "GetMsgFromStatus 方法 ")

[Read 方法](../html/779f6bef-32cb-56c1-ccd8-efabb01df7b6.htm "Read 方法 ")

[ReadBool 方法](../html/8f7c2d91-686b-c59c-54a8-cafee7433947.htm "ReadBool 方法 ")

[ReadPlcType 方法](../html/1c804c7a-7c9f-ca3c-96e0-5b8f327e1018.htm "ReadPlcType 方法 ")

[Start 方法](../html/2c8531c5-1c74-4603-321d-cd96b90d899a.htm "Start 方法 ")

[Stop 方法](../html/c36b817b-fa4d-df99-eb13-d0aa9e83a5df.htm "Stop 方法 ")

[Write 方法](../html/6ab7778a-7ce6-1601-a45d-fb2ad3eac080.htm "Write 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| SiemensPPIHelperRead 方法 |

从西门子的PLC中读取数据信息，地址为"M100","AI100","I0","Q0","V100","S100"等  
Read data information from Siemens PLC with addresses "M100", "AI100", "I0", "Q0", "V100", "S100", etc.

**命名空间：**
 [HslCommunication.Profinet.Siemens.Helper](8dfc52a4-e8ce-f3d2-0647-4d18de8416d3.htm)  
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
	IReadWriteDevice plc,
	string address,
	ushort length,
	byte station,
	Object communicationLock
)
```

```
Public Shared Function Read ( 
	plc As IReadWriteDevice,
	address As String,
	length As UShort,
	station As Byte,
	communicationLock As Object
) As OperateResult(Of Byte())
```

```
public:
static OperateResult<array<unsigned char>^>^ Read(
	IReadWriteDevice^ plc, 
	String^ address, 
	unsigned short length, 
	unsigned char station, 
	Object^ communicationLock
)
```

```
static member Read : 
        plc : IReadWriteDevice * 
        address : string * 
        length : uint16 * 
        station : byte * 
        communicationLock : Object -> OperateResult<byte[]> 
```

#### 参数

plc
:   类型：[HslCommunication.CoreIReadWriteDevice](d987afcc-f99c-6aed-3847-5b1890504635.htm)  
    PLC的通信对象

address
:   类型：SystemString  
    西门子的地址数据信息

length
:   类型：SystemUInt16  
    数据长度

station
:   类型：SystemByte  
    当前的站号信息

communicationLock
:   类型：SystemObject  
    当前的同通信锁

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
带返回结果的结果对象

![](../icons/SectionExpanded.png)参见

#### 引用

[SiemensPPIHelper 类](24deecb2-7eb3-9aba-2da3-137adacdcbd9.htm)

[HslCommunication.Profinet.Siemens.Helper 命名空间](8dfc52a4-e8ce-f3d2-0647-4d18de8416d3.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadBool 方法 

[原文連結](http://api.hslcommunication.cn/html/8f7c2d91-686b-c59c-54a8-cafee7433947.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Siemens.Helper](../html/8dfc52a4-e8ce-f3d2-0647-4d18de8416d3.htm "HslCommunication.Profinet.Siemens.Helper")

[SiemensPPIHelper 类](../html/24deecb2-7eb3-9aba-2da3-137adacdcbd9.htm "SiemensPPIHelper 类")

[SiemensPPIHelper 方法](../html/4b399f3b-1a4b-7dc1-3738-de0332b4df28.htm "SiemensPPIHelper 方法")

[ReadBool 方法](../html/8f7c2d91-686b-c59c-54a8-cafee7433947.htm "ReadBool 方法 ")

[ReadBool 方法 (IReadWriteDevice, String, Byte, Object)](../html/40320454-2b00-ba26-7783-85381ed2c28c.htm "ReadBool 方法 (IReadWriteDevice, String, Byte, Object)")

[ReadBool 方法 (IReadWriteDevice, String, UInt16, Byte, Object)](../html/a00bc8a8-5dad-bd8b-0507-efd55bbd8e7a.htm "ReadBool 方法 (IReadWriteDevice, String, UInt16, Byte, Object)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| SiemensPPIHelperReadBool 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法静态成员 | [ReadBool(IReadWriteDevice, String, Byte, Object)](40320454-2b00-ba26-7783-85381ed2c28c.htm) | 从西门子的PLC中读取bool数据信息，地址为"M100.0","AI100.1","I0.3","Q0.6","V100.4","S100"等  Read bool data information from Siemens PLC, the addresses are "M100.0", "AI100.1", "I0.3", "Q0.6", "V100.4", "S100", etc. |
| 公共方法静态成员 | [ReadBool(IReadWriteDevice, String, UInt16, Byte, Object)](a00bc8a8-5dad-bd8b-0507-efd55bbd8e7a.htm) | 从西门子的PLC中读取bool数据信息，地址为"M100.0","AI100.1","I0.3","Q0.6","V100.4","S100"等  Read bool data information from Siemens PLC, the addresses are "M100.0", "AI100.1", "I0.3", "Q0.6", "V100.4", "S100", etc. |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[SiemensPPIHelper 类](24deecb2-7eb3-9aba-2da3-137adacdcbd9.htm)

[HslCommunication.Profinet.Siemens.Helper 命名空间](8dfc52a4-e8ce-f3d2-0647-4d18de8416d3.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadBool 方法 (IReadWriteDevice, String, Byte, Object)

[原文連結](http://api.hslcommunication.cn/html/40320454-2b00-ba26-7783-85381ed2c28c.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Siemens.Helper](../html/8dfc52a4-e8ce-f3d2-0647-4d18de8416d3.htm "HslCommunication.Profinet.Siemens.Helper")

[SiemensPPIHelper 类](../html/24deecb2-7eb3-9aba-2da3-137adacdcbd9.htm "SiemensPPIHelper 类")

[SiemensPPIHelper 方法](../html/4b399f3b-1a4b-7dc1-3738-de0332b4df28.htm "SiemensPPIHelper 方法")

[ReadBool 方法](../html/8f7c2d91-686b-c59c-54a8-cafee7433947.htm "ReadBool 方法 ")

[ReadBool 方法 (IReadWriteDevice, String, Byte, Object)](../html/40320454-2b00-ba26-7783-85381ed2c28c.htm "ReadBool 方法 (IReadWriteDevice, String, Byte, Object)")

[ReadBool 方法 (IReadWriteDevice, String, UInt16, Byte, Object)](../html/a00bc8a8-5dad-bd8b-0507-efd55bbd8e7a.htm "ReadBool 方法 (IReadWriteDevice, String, UInt16, Byte, Object)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| SiemensPPIHelperReadBool 方法 (IReadWriteDevice, String, Byte, Object) |

从西门子的PLC中读取bool数据信息，地址为"M100.0","AI100.1","I0.3","Q0.6","V100.4","S100"等  
Read bool data information from Siemens PLC, the addresses are "M100.0", "AI100.1", "I0.3", "Q0.6", "V100.4", "S100", etc.

**命名空间：**
 [HslCommunication.Profinet.Siemens.Helper](8dfc52a4-e8ce-f3d2-0647-4d18de8416d3.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult<bool> ReadBool(
	IReadWriteDevice plc,
	string address,
	byte station,
	Object communicationLock
)
```

```
Public Shared Function ReadBool ( 
	plc As IReadWriteDevice,
	address As String,
	station As Byte,
	communicationLock As Object
) As OperateResult(Of Boolean)
```

```
public:
static OperateResult<bool>^ ReadBool(
	IReadWriteDevice^ plc, 
	String^ address, 
	unsigned char station, 
	Object^ communicationLock
)
```

```
static member ReadBool : 
        plc : IReadWriteDevice * 
        address : string * 
        station : byte * 
        communicationLock : Object -> OperateResult<bool> 
```

#### 参数

plc
:   类型：[HslCommunication.CoreIReadWriteDevice](d987afcc-f99c-6aed-3847-5b1890504635.htm)  
    PLC的通信对象

address
:   类型：SystemString  
    西门子的地址数据信息

station
:   类型：SystemByte  
    当前的站号信息

communicationLock
:   类型：SystemObject  
    当前的同通信锁

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Boolean  
带返回结果的结果对象

![](../icons/SectionExpanded.png)参见

#### 引用

[SiemensPPIHelper 类](24deecb2-7eb3-9aba-2da3-137adacdcbd9.htm)

[ReadBool 重载](8f7c2d91-686b-c59c-54a8-cafee7433947.htm)

[HslCommunication.Profinet.Siemens.Helper 命名空间](8dfc52a4-e8ce-f3d2-0647-4d18de8416d3.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadBool 方法 (IReadWriteDevice, String, UInt16, Byte, Object)

[原文連結](http://api.hslcommunication.cn/html/a00bc8a8-5dad-bd8b-0507-efd55bbd8e7a.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Siemens.Helper](../html/8dfc52a4-e8ce-f3d2-0647-4d18de8416d3.htm "HslCommunication.Profinet.Siemens.Helper")

[SiemensPPIHelper 类](../html/24deecb2-7eb3-9aba-2da3-137adacdcbd9.htm "SiemensPPIHelper 类")

[SiemensPPIHelper 方法](../html/4b399f3b-1a4b-7dc1-3738-de0332b4df28.htm "SiemensPPIHelper 方法")

[ReadBool 方法](../html/8f7c2d91-686b-c59c-54a8-cafee7433947.htm "ReadBool 方法 ")

[ReadBool 方法 (IReadWriteDevice, String, Byte, Object)](../html/40320454-2b00-ba26-7783-85381ed2c28c.htm "ReadBool 方法 (IReadWriteDevice, String, Byte, Object)")

[ReadBool 方法 (IReadWriteDevice, String, UInt16, Byte, Object)](../html/a00bc8a8-5dad-bd8b-0507-efd55bbd8e7a.htm "ReadBool 方法 (IReadWriteDevice, String, UInt16, Byte, Object)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| SiemensPPIHelperReadBool 方法 (IReadWriteDevice, String, UInt16, Byte, Object) |

从西门子的PLC中读取bool数据信息，地址为"M100.0","AI100.1","I0.3","Q0.6","V100.4","S100"等  
Read bool data information from Siemens PLC, the addresses are "M100.0", "AI100.1", "I0.3", "Q0.6", "V100.4", "S100", etc.

**命名空间：**
 [HslCommunication.Profinet.Siemens.Helper](8dfc52a4-e8ce-f3d2-0647-4d18de8416d3.htm)  
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
	IReadWriteDevice plc,
	string address,
	ushort length,
	byte station,
	Object communicationLock
)
```

```
Public Shared Function ReadBool ( 
	plc As IReadWriteDevice,
	address As String,
	length As UShort,
	station As Byte,
	communicationLock As Object
) As OperateResult(Of Boolean())
```

```
public:
static OperateResult<array<bool>^>^ ReadBool(
	IReadWriteDevice^ plc, 
	String^ address, 
	unsigned short length, 
	unsigned char station, 
	Object^ communicationLock
)
```

```
static member ReadBool : 
        plc : IReadWriteDevice * 
        address : string * 
        length : uint16 * 
        station : byte * 
        communicationLock : Object -> OperateResult<bool[]> 
```

#### 参数

plc
:   类型：[HslCommunication.CoreIReadWriteDevice](d987afcc-f99c-6aed-3847-5b1890504635.htm)  
    PLC的通信对象

address
:   类型：SystemString  
    西门子的地址数据信息

length
:   类型：SystemUInt16  

    [缺少 "M:HslCommunication.Profinet.Siemens.Helper.SiemensPPIHelper.ReadBool(HslCommunication.Core.IReadWriteDevice,System.String,System.UInt16,System.Byte,System.Object)" 的 <param name="length"/> 文档]

station
:   类型：SystemByte  
    当前的站号信息

communicationLock
:   类型：SystemObject  
    当前的同通信锁

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Boolean  
带返回结果的结果对象

![](../icons/SectionExpanded.png)参见

#### 引用

[SiemensPPIHelper 类](24deecb2-7eb3-9aba-2da3-137adacdcbd9.htm)

[ReadBool 重载](8f7c2d91-686b-c59c-54a8-cafee7433947.htm)

[HslCommunication.Profinet.Siemens.Helper 命名空间](8dfc52a4-e8ce-f3d2-0647-4d18de8416d3.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadPlcType 方法 

[原文連結](http://api.hslcommunication.cn/html/1c804c7a-7c9f-ca3c-96e0-5b8f327e1018.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Siemens.Helper](../html/8dfc52a4-e8ce-f3d2-0647-4d18de8416d3.htm "HslCommunication.Profinet.Siemens.Helper")

[SiemensPPIHelper 类](../html/24deecb2-7eb3-9aba-2da3-137adacdcbd9.htm "SiemensPPIHelper 类")

[SiemensPPIHelper 方法](../html/4b399f3b-1a4b-7dc1-3738-de0332b4df28.htm "SiemensPPIHelper 方法")

[AnalysisAddress 方法](../html/c546b2ce-0787-8f2a-cb22-af4a20e2cc2e.htm "AnalysisAddress 方法 ")

[BuildReadCommand 方法](../html/42b0fc83-2709-ba05-78dc-3a1c90da0bd6.htm "BuildReadCommand 方法 ")

[BuildWriteCommand 方法](../html/98fcc6b4-0da5-b298-7dcf-1b581a10b258.htm "BuildWriteCommand 方法 ")

[CheckReceiveDataComplete 方法](../html/8117ecc5-5e23-bfd4-4a64-2cf40d93daf3.htm "CheckReceiveDataComplete 方法 ")

[CheckResponse 方法](../html/4b18f602-f1bd-ee62-851a-3b68b37e83d8.htm "CheckResponse 方法 ")

[GetExecuteConfirm 方法](../html/12e4e1b2-3109-b243-f4d3-95cebad7fd95.htm "GetExecuteConfirm 方法 ")

[GetMsgFromStatus 方法](../html/5e092d96-9b66-8a86-c4a1-26fab2f0c209.htm "GetMsgFromStatus 方法 ")

[Read 方法](../html/779f6bef-32cb-56c1-ccd8-efabb01df7b6.htm "Read 方法 ")

[ReadBool 方法](../html/8f7c2d91-686b-c59c-54a8-cafee7433947.htm "ReadBool 方法 ")

[ReadPlcType 方法](../html/1c804c7a-7c9f-ca3c-96e0-5b8f327e1018.htm "ReadPlcType 方法 ")

[Start 方法](../html/2c8531c5-1c74-4603-321d-cd96b90d899a.htm "Start 方法 ")

[Stop 方法](../html/c36b817b-fa4d-df99-eb13-d0aa9e83a5df.htm "Stop 方法 ")

[Write 方法](../html/6ab7778a-7ce6-1601-a45d-fb2ad3eac080.htm "Write 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| SiemensPPIHelperReadPlcType 方法 |

读取西门子PLC的型号信息，参数信息可以携带站号信息 "s=2;", 注意，分号是必须的。  
Read the model information of Siemens PLC, the parameter information can carry the station number information "s=2;", note that the semicolon is required.

**命名空间：**
 [HslCommunication.Profinet.Siemens.Helper](8dfc52a4-e8ce-f3d2-0647-4d18de8416d3.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult<string> ReadPlcType(
	IReadWriteDevice plc,
	string parameter,
	byte station,
	Object communicationLock
)
```

```
Public Shared Function ReadPlcType ( 
	plc As IReadWriteDevice,
	parameter As String,
	station As Byte,
	communicationLock As Object
) As OperateResult(Of String)
```

```
public:
static OperateResult<String^>^ ReadPlcType(
	IReadWriteDevice^ plc, 
	String^ parameter, 
	unsigned char station, 
	Object^ communicationLock
)
```

```
static member ReadPlcType : 
        plc : IReadWriteDevice * 
        parameter : string * 
        station : byte * 
        communicationLock : Object -> OperateResult<string> 
```

#### 参数

plc
:   类型：[HslCommunication.CoreIReadWriteDevice](d987afcc-f99c-6aed-3847-5b1890504635.htm)  
    PLC的通信对象

parameter
:   类型：SystemString  
    额外的参数信息，例如可以携带站号信息 "s=2;", 注意，分号是必须的。

station
:   类型：SystemByte  
    当前的站号信息

communicationLock
:   类型：SystemObject  
    当前的同通信锁

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)String  
包含是否成功的结果对象

![](../icons/SectionExpanded.png)参见

#### 引用

[SiemensPPIHelper 类](24deecb2-7eb3-9aba-2da3-137adacdcbd9.htm)

[HslCommunication.Profinet.Siemens.Helper 命名空间](8dfc52a4-e8ce-f3d2-0647-4d18de8416d3.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Start 方法 

[原文連結](http://api.hslcommunication.cn/html/2c8531c5-1c74-4603-321d-cd96b90d899a.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Siemens.Helper](../html/8dfc52a4-e8ce-f3d2-0647-4d18de8416d3.htm "HslCommunication.Profinet.Siemens.Helper")

[SiemensPPIHelper 类](../html/24deecb2-7eb3-9aba-2da3-137adacdcbd9.htm "SiemensPPIHelper 类")

[SiemensPPIHelper 方法](../html/4b399f3b-1a4b-7dc1-3738-de0332b4df28.htm "SiemensPPIHelper 方法")

[AnalysisAddress 方法](../html/c546b2ce-0787-8f2a-cb22-af4a20e2cc2e.htm "AnalysisAddress 方法 ")

[BuildReadCommand 方法](../html/42b0fc83-2709-ba05-78dc-3a1c90da0bd6.htm "BuildReadCommand 方法 ")

[BuildWriteCommand 方法](../html/98fcc6b4-0da5-b298-7dcf-1b581a10b258.htm "BuildWriteCommand 方法 ")

[CheckReceiveDataComplete 方法](../html/8117ecc5-5e23-bfd4-4a64-2cf40d93daf3.htm "CheckReceiveDataComplete 方法 ")

[CheckResponse 方法](../html/4b18f602-f1bd-ee62-851a-3b68b37e83d8.htm "CheckResponse 方法 ")

[GetExecuteConfirm 方法](../html/12e4e1b2-3109-b243-f4d3-95cebad7fd95.htm "GetExecuteConfirm 方法 ")

[GetMsgFromStatus 方法](../html/5e092d96-9b66-8a86-c4a1-26fab2f0c209.htm "GetMsgFromStatus 方法 ")

[Read 方法](../html/779f6bef-32cb-56c1-ccd8-efabb01df7b6.htm "Read 方法 ")

[ReadBool 方法](../html/8f7c2d91-686b-c59c-54a8-cafee7433947.htm "ReadBool 方法 ")

[ReadPlcType 方法](../html/1c804c7a-7c9f-ca3c-96e0-5b8f327e1018.htm "ReadPlcType 方法 ")

[Start 方法](../html/2c8531c5-1c74-4603-321d-cd96b90d899a.htm "Start 方法 ")

[Stop 方法](../html/c36b817b-fa4d-df99-eb13-d0aa9e83a5df.htm "Stop 方法 ")

[Write 方法](../html/6ab7778a-7ce6-1601-a45d-fb2ad3eac080.htm "Write 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| SiemensPPIHelperStart 方法 |

启动西门子PLC为RUN模式，参数信息可以携带站号信息 "s=2;", 注意，分号是必须的。  
Start Siemens PLC in RUN mode, parameter information can carry station number information "s=2;", note that the semicolon is required.

**命名空间：**
 [HslCommunication.Profinet.Siemens.Helper](8dfc52a4-e8ce-f3d2-0647-4d18de8416d3.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult Start(
	IReadWriteDevice plc,
	string parameter,
	byte station,
	Object communicationLock
)
```

```
Public Shared Function Start ( 
	plc As IReadWriteDevice,
	parameter As String,
	station As Byte,
	communicationLock As Object
) As OperateResult
```

```
public:
static OperateResult^ Start(
	IReadWriteDevice^ plc, 
	String^ parameter, 
	unsigned char station, 
	Object^ communicationLock
)
```

```
static member Start : 
        plc : IReadWriteDevice * 
        parameter : string * 
        station : byte * 
        communicationLock : Object -> OperateResult 
```

#### 参数

plc
:   类型：[HslCommunication.CoreIReadWriteDevice](d987afcc-f99c-6aed-3847-5b1890504635.htm)  
    PLC的通信对象

parameter
:   类型：SystemString  
    额外的参数信息，例如可以携带站号信息 "s=2;", 注意，分号是必须的。

station
:   类型：SystemByte  
    当前的站号信息

communicationLock
:   类型：SystemObject  
    当前的同通信锁

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否启动成功

![](../icons/SectionExpanded.png)参见

#### 引用

[SiemensPPIHelper 类](24deecb2-7eb3-9aba-2da3-137adacdcbd9.htm)

[HslCommunication.Profinet.Siemens.Helper 命名空间](8dfc52a4-e8ce-f3d2-0647-4d18de8416d3.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Stop 方法 

[原文連結](http://api.hslcommunication.cn/html/c36b817b-fa4d-df99-eb13-d0aa9e83a5df.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Siemens.Helper](../html/8dfc52a4-e8ce-f3d2-0647-4d18de8416d3.htm "HslCommunication.Profinet.Siemens.Helper")

[SiemensPPIHelper 类](../html/24deecb2-7eb3-9aba-2da3-137adacdcbd9.htm "SiemensPPIHelper 类")

[SiemensPPIHelper 方法](../html/4b399f3b-1a4b-7dc1-3738-de0332b4df28.htm "SiemensPPIHelper 方法")

[AnalysisAddress 方法](../html/c546b2ce-0787-8f2a-cb22-af4a20e2cc2e.htm "AnalysisAddress 方法 ")

[BuildReadCommand 方法](../html/42b0fc83-2709-ba05-78dc-3a1c90da0bd6.htm "BuildReadCommand 方法 ")

[BuildWriteCommand 方法](../html/98fcc6b4-0da5-b298-7dcf-1b581a10b258.htm "BuildWriteCommand 方法 ")

[CheckReceiveDataComplete 方法](../html/8117ecc5-5e23-bfd4-4a64-2cf40d93daf3.htm "CheckReceiveDataComplete 方法 ")

[CheckResponse 方法](../html/4b18f602-f1bd-ee62-851a-3b68b37e83d8.htm "CheckResponse 方法 ")

[GetExecuteConfirm 方法](../html/12e4e1b2-3109-b243-f4d3-95cebad7fd95.htm "GetExecuteConfirm 方法 ")

[GetMsgFromStatus 方法](../html/5e092d96-9b66-8a86-c4a1-26fab2f0c209.htm "GetMsgFromStatus 方法 ")

[Read 方法](../html/779f6bef-32cb-56c1-ccd8-efabb01df7b6.htm "Read 方法 ")

[ReadBool 方法](../html/8f7c2d91-686b-c59c-54a8-cafee7433947.htm "ReadBool 方法 ")

[ReadPlcType 方法](../html/1c804c7a-7c9f-ca3c-96e0-5b8f327e1018.htm "ReadPlcType 方法 ")

[Start 方法](../html/2c8531c5-1c74-4603-321d-cd96b90d899a.htm "Start 方法 ")

[Stop 方法](../html/c36b817b-fa4d-df99-eb13-d0aa9e83a5df.htm "Stop 方法 ")

[Write 方法](../html/6ab7778a-7ce6-1601-a45d-fb2ad3eac080.htm "Write 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| SiemensPPIHelperStop 方法 |

停止西门子PLC，切换为Stop模式，参数信息可以携带站号信息 "s=2;", 注意，分号是必须的。  
Stop Siemens PLC and switch to Stop mode, parameter information can carry station number information "s=2;", note that the semicolon is required.

**命名空间：**
 [HslCommunication.Profinet.Siemens.Helper](8dfc52a4-e8ce-f3d2-0647-4d18de8416d3.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult Stop(
	IReadWriteDevice plc,
	string parameter,
	byte station,
	Object communicationLock
)
```

```
Public Shared Function Stop ( 
	plc As IReadWriteDevice,
	parameter As String,
	station As Byte,
	communicationLock As Object
) As OperateResult
```

```
public:
static OperateResult^ Stop(
	IReadWriteDevice^ plc, 
	String^ parameter, 
	unsigned char station, 
	Object^ communicationLock
)
```

```
static member Stop : 
        plc : IReadWriteDevice * 
        parameter : string * 
        station : byte * 
        communicationLock : Object -> OperateResult 
```

#### 参数

plc
:   类型：[HslCommunication.CoreIReadWriteDevice](d987afcc-f99c-6aed-3847-5b1890504635.htm)  
    PLC的通信对象

parameter
:   类型：SystemString  
    额外的参数信息，例如可以携带站号信息 "s=2;", 注意，分号是必须的。

station
:   类型：SystemByte  
    当前的站号信息

communicationLock
:   类型：SystemObject  
    当前的同通信锁

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否停止成功

![](../icons/SectionExpanded.png)参见

#### 引用

[SiemensPPIHelper 类](24deecb2-7eb3-9aba-2da3-137adacdcbd9.htm)

[HslCommunication.Profinet.Siemens.Helper 命名空间](8dfc52a4-e8ce-f3d2-0647-4d18de8416d3.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)