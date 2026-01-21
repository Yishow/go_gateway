# HslCommunication - HslCommunication.Profinet.LSIS.Helper

> 分類頁數: 30



---
## HslCommunication.Profinet.LSIS.Helper

[原文連結](http://api.hslcommunication.cn/html/0200f8be-eba3-dea9-4eb6-b51d8c230c6b.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.LSIS.Helper](../html/0200f8be-eba3-dea9-4eb6-b51d8c230c6b.htm "HslCommunication.Profinet.LSIS.Helper")

[LSCnetHelper 类](../html/c874de34-e25a-6c9c-8d0e-ddffe0651200.htm "LSCnetHelper 类")

[LSCpuHelper 类](../html/6aac48e7-8ea0-56d0-c53c-eb3f51f736ba.htm "LSCpuHelper 类")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslCommunication.Profinet.LSIS.Helper 命名空间 |

[缺少 "N:HslCommunication.Profinet.LSIS.Helper" 的 <summary> 文档]

![](../icons/SectionExpanded.png)类

|  | 类 | 说明 |
| --- | --- | --- |
| 公共类 | [LSCnetHelper](c874de34-e25a-6c9c-8d0e-ddffe0651200.htm) | Cnet的辅助类 |
| 公共类 | [LSCpuHelper](6aac48e7-8ea0-56d0-c53c-eb3f51f736ba.htm) | [LSCpu](70a769a2-3ce3-ef6b-d7e2-38132b33a8bf.htm)相关的辅助类，提供了一些辅助的静态方法信息 |

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## LSCnetHelper 类

[原文連結](http://api.hslcommunication.cn/html/c874de34-e25a-6c9c-8d0e-ddffe0651200.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.LSIS.Helper](../html/0200f8be-eba3-dea9-4eb6-b51d8c230c6b.htm "HslCommunication.Profinet.LSIS.Helper")

[LSCnetHelper 类](../html/c874de34-e25a-6c9c-8d0e-ddffe0651200.htm "LSCnetHelper 类")

[LSCnetHelper 构造函数](../html/390c8142-86d6-bc82-208a-d215c39146fb.htm "LSCnetHelper 构造函数 ")

[LSCnetHelper 方法](../html/d8e60529-bf76-1a12-c052-e3bfeaa4ec33.htm "LSCnetHelper 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| LSCnetHelper 类 |

Cnet的辅助类

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  HslCommunication.Profinet.LSIS.HelperLSCnetHelper

**命名空间：**
 [HslCommunication.Profinet.LSIS.Helper](0200f8be-eba3-dea9-4eb6-b51d8c230c6b.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class LSCnetHelper
```

```
Public Class LSCnetHelper
```

```
public ref class LSCnetHelper
```

```
type LSCnetHelper =  class end
```

LSCnetHelper 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [LSCnetHelper](390c8142-86d6-bc82-208a-d215c39146fb.htm) | 初始化 LSCnetHelper 类的一个新实例 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法静态成员 | [AnalysisAddress](15ee9fd3-6c4b-6775-db86-370121c062ff.htm) | 从输入的地址里解析出真实的可以直接放到协议的地址信息，如果是X的地址，自动转换带小数点的表示方式到位地址，如果是其他类型地址，则一律统一转化为字节为单位的地址  The real address information that can be directly put into the protocol is parsed from the input address. If it is the address of X, it will automatically convert the representation with a decimal point to the address. If it is an address of other types, it will be uniformly converted into a unit of bytes. address |
| 公共方法静态成员 | [BuildReadByteCommand](187eecaa-60f5-ed67-f0c6-ce15a59e1e52.htm) | reading address Type of ReadByte |
| 公共方法静态成员 | [BuildReadIndividualCommand(Byte, String)](6804f38b-02ff-c012-e0bb-5cf53c817cc3.htm) | Multi reading address Type of Read Individual |
| 公共方法静态成员 | [BuildReadIndividualCommand(Byte, String)](9747205d-7347-ab27-dcf9-fab8e8efbe9a.htm) | Multi reading address Type of Read Individual |
| 公共方法静态成员 | [BuildWriteByteCommand](62a3477c-e064-8c71-9687-4f1cc45ebc16.htm) | write data to address Type of ReadByte |
| 公共方法静态成员 | [BuildWriteCommand](2873d537-6750-08f6-ac89-0bf2ea205de3.htm) | write data to address Type of ReadByte |
| 公共方法静态成员 | [BuildWriteOneCommand](c59dc3f6-675c-8f9d-9dce-3c0857887b3b.htm) | write data to address Type of One |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法静态成员 | [GetAddressOfU\_Q\_I](c3212593-efad-0a6e-22a6-11be146a83d9.htm) |  |
| 公共方法静态成员 | [GetErrorText](d8d3953e-a826-5558-6773-8809cf6d9dbe.htm) | 根据错误号，获取到真实的错误描述信息  According to the error number, get the real error description information |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法静态成员 | [Read(IReadWriteDevice, Int32, String)](5af4122b-f6be-660b-080d-c930b67f90e9.htm) | 从PLC设备读取多个地址的数据信息，返回连续的字节数组，需要按照实际情况进行按顺序解析。  Read the data information of multiple addresses from the PLC device and return a continuous byte array, which needs to be parsed in order according to the actual situation. |
| 公共方法静态成员 | [Read(IReadWriteDeviceStation, String, UInt16)](c2735b0c-d09d-bfee-e23b-628145ca9709.htm) | 从PLC的指定地址读取原始的字节数据信息，地址示例：MB100, MW100, MD100, 如果输入了M100等同于MB100  Read the original byte data information from the designated address of the PLC. Examples of addresses: MB100, MW100, MD100, if the input M100 is equivalent to MB100 |
| 公共方法静态成员 | [ReadAsync(IReadWriteDevice, Int32, String)](c685c2aa-fd95-7051-a432-edf9d1eff0a5.htm) | 从PLC设备读取多个地址的数据信息，返回连续的字节数组，需要按照实际情况进行按顺序解析。  Read the data information of multiple addresses from the PLC device and return a continuous byte array, which needs to be parsed in order according to the actual situation. |
| 公共方法静态成员 | [ReadAsync(IReadWriteDeviceStation, String, UInt16)](3c28f44d-b751-2afc-e6c4-0a66a210160b.htm) | 从PLC的指定地址读取原始的字节数据信息，地址示例：MB100, MW100, MD100, 如果输入了M100等同于MB100  Read the original byte data information from the designated address of the PLC. Examples of addresses: MB100, MW100, MD100, if the input M100 is equivalent to MB100 |
| 公共方法静态成员 | [ReadBool(IReadWriteDevice, Int32, String)](6c90e993-92e2-4b41-a94f-9233ae54991d.htm) | 从PLC的指定地址读取原始的位数据信息，地址示例：MX100, MX10A  Read the original bool data information from the designated address of the PLC. Examples of addresses: MX100, MX10A |
| 公共方法静态成员 | [ReadBool(IReadWriteDeviceStation, String, UInt16)](84c5849a-9b04-82f1-c992-c67d92c7218a.htm) | 从PLC的指定地址读取原始的位数据信息，地址示例：MB100.0, MW100.0  Read the original bool data information from the designated address of the PLC. Examples of addresses: MB100.0, MW100.0 |
| 公共方法静态成员 | [ReadBoolAsync(IReadWriteDevice, Int32, String)](e8d25674-d63b-00d8-dbd9-15736a738497.htm) | 从PLC的指定地址读取原始的位数据信息，地址示例：MX100, MX10A  Read the original bool data information from the designated address of the PLC. Examples of addresses: MX100, MX10A |
| 公共方法静态成员 | [ReadBoolAsync(IReadWriteDeviceStation, String, UInt16)](74cb9308-0f0a-23f2-9d31-73fe16baab20.htm) | 从PLC的指定地址读取原始的位数据信息，地址示例：MB100.0, MW100.0  Read the original bool data information from the designated address of the PLC. Examples of addresses: MB100.0, MW100.0 |
| 公共方法 | ToString | (继承自 Object。) |
| 公共方法静态成员 | [UnpackResponseContent](ef40b9f3-0fc8-b1dd-945b-8d2011f1b6a9.htm) | 根据对方返回的报文命令，对命令进行基本的拆包，例如各种Modbus协议拆包为统一的核心报文，还支持对报文的验证  According to the message command returned by the other party, the command is basically unpacked, for example, various Modbus protocols are unpacked into a unified core message, and the verification of the message is also supported |
| 公共方法静态成员 | [Write(IReadWriteDevice, Int32, String, Boolean)](e253354d-f080-9498-b34e-710eb2ff95c3.htm) | 将bool数据写入到PLC的指定的地址里，地址示例：MX100, MX10A  Write the bool data to the designated address of the PLC. Examples of addresses: MX100, MX10A |
| 公共方法静态成员 | [Write(IReadWriteDevice, Int32, String, Byte)](9fb7ffe6-6132-bbd7-3287-b65dc8956570.htm) | 将原始数据写入到PLC的指定的地址里，地址示例：MB100, MW100, MD100, 如果输入了M100等同于MB100  Write the original data to the designated address of the PLC. Examples of addresses: MB100, MW100, MD100, if input M100 is equivalent to MB100 |
| 公共方法静态成员 | [WriteAsync(IReadWriteDevice, Int32, String, Boolean)](318c6b22-8e31-1bd1-2e6c-4d66959d0fb5.htm) | 将bool数据写入到PLC的指定的地址里，地址示例：MX100, MX10A  Write the bool data to the designated address of the PLC. Examples of addresses: MX100, MX10A |
| 公共方法静态成员 | [WriteAsync(IReadWriteDevice, Int32, String, Byte)](724ba1f9-5ac5-c55d-803c-bf25a3ec9e1f.htm) | 将原始数据写入到PLC的指定的地址里，地址示例：MB100, MW100, MD100, 如果输入了M100等同于MB100  Write the original data to the designated address of the PLC. Examples of addresses: MB100, MW100, MD100, if input M100 is equivalent to MB100 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.Profinet.LSIS.Helper 命名空间](0200f8be-eba3-dea9-4eb6-b51d8c230c6b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## LSCnetHelper 构造函数 

[原文連結](http://api.hslcommunication.cn/html/390c8142-86d6-bc82-208a-d215c39146fb.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.LSIS.Helper](../html/0200f8be-eba3-dea9-4eb6-b51d8c230c6b.htm "HslCommunication.Profinet.LSIS.Helper")

[LSCnetHelper 类](../html/c874de34-e25a-6c9c-8d0e-ddffe0651200.htm "LSCnetHelper 类")

[LSCnetHelper 构造函数](../html/390c8142-86d6-bc82-208a-d215c39146fb.htm "LSCnetHelper 构造函数 ")

[LSCnetHelper 方法](../html/d8e60529-bf76-1a12-c052-e3bfeaa4ec33.htm "LSCnetHelper 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| LSCnetHelper 构造函数 |

初始化 [LSCnetHelper](c874de34-e25a-6c9c-8d0e-ddffe0651200.htm) 类的一个新实例

**命名空间：**
 [HslCommunication.Profinet.LSIS.Helper](0200f8be-eba3-dea9-4eb6-b51d8c230c6b.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public LSCnetHelper()
```

```
Public Sub New
```

```
public:
LSCnetHelper()
```

```
new : unit -> LSCnetHelper
```

![](../icons/SectionExpanded.png)参见

#### 引用

[LSCnetHelper 类](c874de34-e25a-6c9c-8d0e-ddffe0651200.htm)

[HslCommunication.Profinet.LSIS.Helper 命名空间](0200f8be-eba3-dea9-4eb6-b51d8c230c6b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## LSCnetHelper 方法

[原文連結](http://api.hslcommunication.cn/html/d8e60529-bf76-1a12-c052-e3bfeaa4ec33.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.LSIS.Helper](../html/0200f8be-eba3-dea9-4eb6-b51d8c230c6b.htm "HslCommunication.Profinet.LSIS.Helper")

[LSCnetHelper 类](../html/c874de34-e25a-6c9c-8d0e-ddffe0651200.htm "LSCnetHelper 类")

[LSCnetHelper 方法](../html/d8e60529-bf76-1a12-c052-e3bfeaa4ec33.htm "LSCnetHelper 方法")

[AnalysisAddress 方法](../html/15ee9fd3-6c4b-6775-db86-370121c062ff.htm "AnalysisAddress 方法 ")

[BuildReadByteCommand 方法](../html/187eecaa-60f5-ed67-f0c6-ce15a59e1e52.htm "BuildReadByteCommand 方法 ")

[BuildReadIndividualCommand 方法](../html/883ada55-fd41-162f-09a9-2d34f92ac605.htm "BuildReadIndividualCommand 方法 ")

[BuildWriteByteCommand 方法](../html/62a3477c-e064-8c71-9687-4f1cc45ebc16.htm "BuildWriteByteCommand 方法 ")

[BuildWriteCommand 方法](../html/2873d537-6750-08f6-ac89-0bf2ea205de3.htm "BuildWriteCommand 方法 ")

[BuildWriteOneCommand 方法](../html/c59dc3f6-675c-8f9d-9dce-3c0857887b3b.htm "BuildWriteOneCommand 方法 ")

[GetAddressOfU\_Q\_I 方法](../html/c3212593-efad-0a6e-22a6-11be146a83d9.htm "GetAddressOfU_Q_I 方法 ")

[GetErrorText 方法](../html/d8d3953e-a826-5558-6773-8809cf6d9dbe.htm "GetErrorText 方法 ")

[Read 方法](../html/7c298c2b-c4d6-de5b-8b67-1ddb7fb1f180.htm "Read 方法 ")

[ReadAsync 方法](../html/f3d85a1f-1978-d789-f6a0-f1dd0fb17e8c.htm "ReadAsync 方法 ")

[ReadBool 方法](../html/3359728b-1b82-547f-13ea-05d524a7bacc.htm "ReadBool 方法 ")

[ReadBoolAsync 方法](../html/6b11ae2c-5cb7-a003-2852-48d340e12c62.htm "ReadBoolAsync 方法 ")

[UnpackResponseContent 方法](../html/ef40b9f3-0fc8-b1dd-945b-8d2011f1b6a9.htm "UnpackResponseContent 方法 ")

[Write 方法](../html/2792b231-69a8-6fae-dc81-e778e73172e1.htm "Write 方法 ")

[WriteAsync 方法](../html/a8b58372-0a53-0c10-6340-30584a178a1b.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| LSCnetHelper 方法 |

[LSCnetHelper](c874de34-e25a-6c9c-8d0e-ddffe0651200.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法静态成员 | [AnalysisAddress](15ee9fd3-6c4b-6775-db86-370121c062ff.htm) | 从输入的地址里解析出真实的可以直接放到协议的地址信息，如果是X的地址，自动转换带小数点的表示方式到位地址，如果是其他类型地址，则一律统一转化为字节为单位的地址  The real address information that can be directly put into the protocol is parsed from the input address. If it is the address of X, it will automatically convert the representation with a decimal point to the address. If it is an address of other types, it will be uniformly converted into a unit of bytes. address |
| 公共方法静态成员 | [BuildReadByteCommand](187eecaa-60f5-ed67-f0c6-ce15a59e1e52.htm) | reading address Type of ReadByte |
| 公共方法静态成员 | [BuildReadIndividualCommand(Byte, String)](6804f38b-02ff-c012-e0bb-5cf53c817cc3.htm) | Multi reading address Type of Read Individual |
| 公共方法静态成员 | [BuildReadIndividualCommand(Byte, String)](9747205d-7347-ab27-dcf9-fab8e8efbe9a.htm) | Multi reading address Type of Read Individual |
| 公共方法静态成员 | [BuildWriteByteCommand](62a3477c-e064-8c71-9687-4f1cc45ebc16.htm) | write data to address Type of ReadByte |
| 公共方法静态成员 | [BuildWriteCommand](2873d537-6750-08f6-ac89-0bf2ea205de3.htm) | write data to address Type of ReadByte |
| 公共方法静态成员 | [BuildWriteOneCommand](c59dc3f6-675c-8f9d-9dce-3c0857887b3b.htm) | write data to address Type of One |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法静态成员 | [GetAddressOfU\_Q\_I](c3212593-efad-0a6e-22a6-11be146a83d9.htm) |  |
| 公共方法静态成员 | [GetErrorText](d8d3953e-a826-5558-6773-8809cf6d9dbe.htm) | 根据错误号，获取到真实的错误描述信息  According to the error number, get the real error description information |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法静态成员 | [Read(IReadWriteDevice, Int32, String)](5af4122b-f6be-660b-080d-c930b67f90e9.htm) | 从PLC设备读取多个地址的数据信息，返回连续的字节数组，需要按照实际情况进行按顺序解析。  Read the data information of multiple addresses from the PLC device and return a continuous byte array, which needs to be parsed in order according to the actual situation. |
| 公共方法静态成员 | [Read(IReadWriteDeviceStation, String, UInt16)](c2735b0c-d09d-bfee-e23b-628145ca9709.htm) | 从PLC的指定地址读取原始的字节数据信息，地址示例：MB100, MW100, MD100, 如果输入了M100等同于MB100  Read the original byte data information from the designated address of the PLC. Examples of addresses: MB100, MW100, MD100, if the input M100 is equivalent to MB100 |
| 公共方法静态成员 | [ReadAsync(IReadWriteDevice, Int32, String)](c685c2aa-fd95-7051-a432-edf9d1eff0a5.htm) | 从PLC设备读取多个地址的数据信息，返回连续的字节数组，需要按照实际情况进行按顺序解析。  Read the data information of multiple addresses from the PLC device and return a continuous byte array, which needs to be parsed in order according to the actual situation. |
| 公共方法静态成员 | [ReadAsync(IReadWriteDeviceStation, String, UInt16)](3c28f44d-b751-2afc-e6c4-0a66a210160b.htm) | 从PLC的指定地址读取原始的字节数据信息，地址示例：MB100, MW100, MD100, 如果输入了M100等同于MB100  Read the original byte data information from the designated address of the PLC. Examples of addresses: MB100, MW100, MD100, if the input M100 is equivalent to MB100 |
| 公共方法静态成员 | [ReadBool(IReadWriteDevice, Int32, String)](6c90e993-92e2-4b41-a94f-9233ae54991d.htm) | 从PLC的指定地址读取原始的位数据信息，地址示例：MX100, MX10A  Read the original bool data information from the designated address of the PLC. Examples of addresses: MX100, MX10A |
| 公共方法静态成员 | [ReadBool(IReadWriteDeviceStation, String, UInt16)](84c5849a-9b04-82f1-c992-c67d92c7218a.htm) | 从PLC的指定地址读取原始的位数据信息，地址示例：MB100.0, MW100.0  Read the original bool data information from the designated address of the PLC. Examples of addresses: MB100.0, MW100.0 |
| 公共方法静态成员 | [ReadBoolAsync(IReadWriteDevice, Int32, String)](e8d25674-d63b-00d8-dbd9-15736a738497.htm) | 从PLC的指定地址读取原始的位数据信息，地址示例：MX100, MX10A  Read the original bool data information from the designated address of the PLC. Examples of addresses: MX100, MX10A |
| 公共方法静态成员 | [ReadBoolAsync(IReadWriteDeviceStation, String, UInt16)](74cb9308-0f0a-23f2-9d31-73fe16baab20.htm) | 从PLC的指定地址读取原始的位数据信息，地址示例：MB100.0, MW100.0  Read the original bool data information from the designated address of the PLC. Examples of addresses: MB100.0, MW100.0 |
| 公共方法 | ToString | (继承自 Object。) |
| 公共方法静态成员 | [UnpackResponseContent](ef40b9f3-0fc8-b1dd-945b-8d2011f1b6a9.htm) | 根据对方返回的报文命令，对命令进行基本的拆包，例如各种Modbus协议拆包为统一的核心报文，还支持对报文的验证  According to the message command returned by the other party, the command is basically unpacked, for example, various Modbus protocols are unpacked into a unified core message, and the verification of the message is also supported |
| 公共方法静态成员 | [Write(IReadWriteDevice, Int32, String, Boolean)](e253354d-f080-9498-b34e-710eb2ff95c3.htm) | 将bool数据写入到PLC的指定的地址里，地址示例：MX100, MX10A  Write the bool data to the designated address of the PLC. Examples of addresses: MX100, MX10A |
| 公共方法静态成员 | [Write(IReadWriteDevice, Int32, String, Byte)](9fb7ffe6-6132-bbd7-3287-b65dc8956570.htm) | 将原始数据写入到PLC的指定的地址里，地址示例：MB100, MW100, MD100, 如果输入了M100等同于MB100  Write the original data to the designated address of the PLC. Examples of addresses: MB100, MW100, MD100, if input M100 is equivalent to MB100 |
| 公共方法静态成员 | [WriteAsync(IReadWriteDevice, Int32, String, Boolean)](318c6b22-8e31-1bd1-2e6c-4d66959d0fb5.htm) | 将bool数据写入到PLC的指定的地址里，地址示例：MX100, MX10A  Write the bool data to the designated address of the PLC. Examples of addresses: MX100, MX10A |
| 公共方法静态成员 | [WriteAsync(IReadWriteDevice, Int32, String, Byte)](724ba1f9-5ac5-c55d-803c-bf25a3ec9e1f.htm) | 将原始数据写入到PLC的指定的地址里，地址示例：MB100, MW100, MD100, 如果输入了M100等同于MB100  Write the original data to the designated address of the PLC. Examples of addresses: MB100, MW100, MD100, if input M100 is equivalent to MB100 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[LSCnetHelper 类](c874de34-e25a-6c9c-8d0e-ddffe0651200.htm)

[HslCommunication.Profinet.LSIS.Helper 命名空间](0200f8be-eba3-dea9-4eb6-b51d8c230c6b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## AnalysisAddress 方法 

[原文連結](http://api.hslcommunication.cn/html/15ee9fd3-6c4b-6775-db86-370121c062ff.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.LSIS.Helper](../html/0200f8be-eba3-dea9-4eb6-b51d8c230c6b.htm "HslCommunication.Profinet.LSIS.Helper")

[LSCnetHelper 类](../html/c874de34-e25a-6c9c-8d0e-ddffe0651200.htm "LSCnetHelper 类")

[LSCnetHelper 方法](../html/d8e60529-bf76-1a12-c052-e3bfeaa4ec33.htm "LSCnetHelper 方法")

[AnalysisAddress 方法](../html/15ee9fd3-6c4b-6775-db86-370121c062ff.htm "AnalysisAddress 方法 ")

[BuildReadByteCommand 方法](../html/187eecaa-60f5-ed67-f0c6-ce15a59e1e52.htm "BuildReadByteCommand 方法 ")

[BuildReadIndividualCommand 方法](../html/883ada55-fd41-162f-09a9-2d34f92ac605.htm "BuildReadIndividualCommand 方法 ")

[BuildWriteByteCommand 方法](../html/62a3477c-e064-8c71-9687-4f1cc45ebc16.htm "BuildWriteByteCommand 方法 ")

[BuildWriteCommand 方法](../html/2873d537-6750-08f6-ac89-0bf2ea205de3.htm "BuildWriteCommand 方法 ")

[BuildWriteOneCommand 方法](../html/c59dc3f6-675c-8f9d-9dce-3c0857887b3b.htm "BuildWriteOneCommand 方法 ")

[GetAddressOfU\_Q\_I 方法](../html/c3212593-efad-0a6e-22a6-11be146a83d9.htm "GetAddressOfU_Q_I 方法 ")

[GetErrorText 方法](../html/d8d3953e-a826-5558-6773-8809cf6d9dbe.htm "GetErrorText 方法 ")

[Read 方法](../html/7c298c2b-c4d6-de5b-8b67-1ddb7fb1f180.htm "Read 方法 ")

[ReadAsync 方法](../html/f3d85a1f-1978-d789-f6a0-f1dd0fb17e8c.htm "ReadAsync 方法 ")

[ReadBool 方法](../html/3359728b-1b82-547f-13ea-05d524a7bacc.htm "ReadBool 方法 ")

[ReadBoolAsync 方法](../html/6b11ae2c-5cb7-a003-2852-48d340e12c62.htm "ReadBoolAsync 方法 ")

[UnpackResponseContent 方法](../html/ef40b9f3-0fc8-b1dd-945b-8d2011f1b6a9.htm "UnpackResponseContent 方法 ")

[Write 方法](../html/2792b231-69a8-6fae-dc81-e778e73172e1.htm "Write 方法 ")

[WriteAsync 方法](../html/a8b58372-0a53-0c10-6340-30584a178a1b.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| LSCnetHelperAnalysisAddress 方法 |

从输入的地址里解析出真实的可以直接放到协议的地址信息，如果是X的地址，自动转换带小数点的表示方式到位地址，如果是其他类型地址，则一律统一转化为字节为单位的地址  
The real address information that can be directly put into the protocol is parsed from the input address. If it is the address of X,
it will automatically convert the representation with a decimal point to the address. If it is an address of other types, it will be uniformly converted into a unit of bytes. address

**命名空间：**
 [HslCommunication.Profinet.LSIS.Helper](0200f8be-eba3-dea9-4eb6-b51d8c230c6b.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult<string> AnalysisAddress(
	string address,
	bool IsWrite = false
)
```

```
Public Shared Function AnalysisAddress ( 
	address As String,
	Optional IsWrite As Boolean = false
) As OperateResult(Of String)
```

```
public:
static OperateResult<String^>^ AnalysisAddress(
	String^ address, 
	bool IsWrite = false
)
```

```
static member AnalysisAddress : 
        address : string * 
        ?IsWrite : bool 
(* Defaults:
        let _IsWrite = defaultArg IsWrite false
*)
-> OperateResult<string> 
```

#### 参数

address
:   类型：SystemString  
    输入的起始偏移地址

IsWrite (Optional)
:   类型：SystemBoolean  
    是否转换为bool地址

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)String  
analysis result

![](../icons/SectionExpanded.png)参见

#### 引用

[LSCnetHelper 类](c874de34-e25a-6c9c-8d0e-ddffe0651200.htm)

[HslCommunication.Profinet.LSIS.Helper 命名空间](0200f8be-eba3-dea9-4eb6-b51d8c230c6b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## BuildReadByteCommand 方法 

[原文連結](http://api.hslcommunication.cn/html/187eecaa-60f5-ed67-f0c6-ce15a59e1e52.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.LSIS.Helper](../html/0200f8be-eba3-dea9-4eb6-b51d8c230c6b.htm "HslCommunication.Profinet.LSIS.Helper")

[LSCnetHelper 类](../html/c874de34-e25a-6c9c-8d0e-ddffe0651200.htm "LSCnetHelper 类")

[LSCnetHelper 方法](../html/d8e60529-bf76-1a12-c052-e3bfeaa4ec33.htm "LSCnetHelper 方法")

[AnalysisAddress 方法](../html/15ee9fd3-6c4b-6775-db86-370121c062ff.htm "AnalysisAddress 方法 ")

[BuildReadByteCommand 方法](../html/187eecaa-60f5-ed67-f0c6-ce15a59e1e52.htm "BuildReadByteCommand 方法 ")

[BuildReadIndividualCommand 方法](../html/883ada55-fd41-162f-09a9-2d34f92ac605.htm "BuildReadIndividualCommand 方法 ")

[BuildWriteByteCommand 方法](../html/62a3477c-e064-8c71-9687-4f1cc45ebc16.htm "BuildWriteByteCommand 方法 ")

[BuildWriteCommand 方法](../html/2873d537-6750-08f6-ac89-0bf2ea205de3.htm "BuildWriteCommand 方法 ")

[BuildWriteOneCommand 方法](../html/c59dc3f6-675c-8f9d-9dce-3c0857887b3b.htm "BuildWriteOneCommand 方法 ")

[GetAddressOfU\_Q\_I 方法](../html/c3212593-efad-0a6e-22a6-11be146a83d9.htm "GetAddressOfU_Q_I 方法 ")

[GetErrorText 方法](../html/d8d3953e-a826-5558-6773-8809cf6d9dbe.htm "GetErrorText 方法 ")

[Read 方法](../html/7c298c2b-c4d6-de5b-8b67-1ddb7fb1f180.htm "Read 方法 ")

[ReadAsync 方法](../html/f3d85a1f-1978-d789-f6a0-f1dd0fb17e8c.htm "ReadAsync 方法 ")

[ReadBool 方法](../html/3359728b-1b82-547f-13ea-05d524a7bacc.htm "ReadBool 方法 ")

[ReadBoolAsync 方法](../html/6b11ae2c-5cb7-a003-2852-48d340e12c62.htm "ReadBoolAsync 方法 ")

[UnpackResponseContent 方法](../html/ef40b9f3-0fc8-b1dd-945b-8d2011f1b6a9.htm "UnpackResponseContent 方法 ")

[Write 方法](../html/2792b231-69a8-6fae-dc81-e778e73172e1.htm "Write 方法 ")

[WriteAsync 方法](../html/a8b58372-0a53-0c10-6340-30584a178a1b.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| LSCnetHelperBuildReadByteCommand 方法 |

reading address Type of ReadByte

**命名空间：**
 [HslCommunication.Profinet.LSIS.Helper](0200f8be-eba3-dea9-4eb6-b51d8c230c6b.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult<List<byte[]>> BuildReadByteCommand(
	byte station,
	string address,
	ushort length
)
```

```
Public Shared Function BuildReadByteCommand ( 
	station As Byte,
	address As String,
	length As UShort
) As OperateResult(Of List(Of Byte()))
```

```
public:
static OperateResult<List<array<unsigned char>^>^>^ BuildReadByteCommand(
	unsigned char station, 
	String^ address, 
	unsigned short length
)
```

```
static member BuildReadByteCommand : 
        station : byte * 
        address : string * 
        length : uint16 -> OperateResult<List<byte[]>> 
```

#### 参数

station
:   类型：SystemByte  
    plc station

address
:   类型：SystemString  
    address, for example: M100, D100, DW100

length
:   类型：SystemUInt16  
    read length

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)ListByte  
command bytes

![](../icons/SectionExpanded.png)参见

#### 引用

[LSCnetHelper 类](c874de34-e25a-6c9c-8d0e-ddffe0651200.htm)

[HslCommunication.Profinet.LSIS.Helper 命名空间](0200f8be-eba3-dea9-4eb6-b51d8c230c6b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## BuildReadIndividualCommand 方法 

[原文連結](http://api.hslcommunication.cn/html/883ada55-fd41-162f-09a9-2d34f92ac605.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.LSIS.Helper](../html/0200f8be-eba3-dea9-4eb6-b51d8c230c6b.htm "HslCommunication.Profinet.LSIS.Helper")

[LSCnetHelper 类](../html/c874de34-e25a-6c9c-8d0e-ddffe0651200.htm "LSCnetHelper 类")

[LSCnetHelper 方法](../html/d8e60529-bf76-1a12-c052-e3bfeaa4ec33.htm "LSCnetHelper 方法")

[BuildReadIndividualCommand 方法](../html/883ada55-fd41-162f-09a9-2d34f92ac605.htm "BuildReadIndividualCommand 方法 ")

[BuildReadIndividualCommand 方法 (Byte, String)](../html/6804f38b-02ff-c012-e0bb-5cf53c817cc3.htm "BuildReadIndividualCommand 方法 (Byte, String)")

[BuildReadIndividualCommand 方法 (Byte, String[])](../html/9747205d-7347-ab27-dcf9-fab8e8efbe9a.htm "BuildReadIndividualCommand 方法 (Byte, String[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| LSCnetHelperBuildReadIndividualCommand 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法静态成员 | [BuildReadIndividualCommand(Byte, String)](6804f38b-02ff-c012-e0bb-5cf53c817cc3.htm) | Multi reading address Type of Read Individual |
| 公共方法静态成员 | [BuildReadIndividualCommand(Byte, String)](9747205d-7347-ab27-dcf9-fab8e8efbe9a.htm) | Multi reading address Type of Read Individual |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[LSCnetHelper 类](c874de34-e25a-6c9c-8d0e-ddffe0651200.htm)

[HslCommunication.Profinet.LSIS.Helper 命名空间](0200f8be-eba3-dea9-4eb6-b51d8c230c6b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## BuildReadIndividualCommand 方法 (Byte, String)

[原文連結](http://api.hslcommunication.cn/html/6804f38b-02ff-c012-e0bb-5cf53c817cc3.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.LSIS.Helper](../html/0200f8be-eba3-dea9-4eb6-b51d8c230c6b.htm "HslCommunication.Profinet.LSIS.Helper")

[LSCnetHelper 类](../html/c874de34-e25a-6c9c-8d0e-ddffe0651200.htm "LSCnetHelper 类")

[LSCnetHelper 方法](../html/d8e60529-bf76-1a12-c052-e3bfeaa4ec33.htm "LSCnetHelper 方法")

[BuildReadIndividualCommand 方法](../html/883ada55-fd41-162f-09a9-2d34f92ac605.htm "BuildReadIndividualCommand 方法 ")

[BuildReadIndividualCommand 方法 (Byte, String)](../html/6804f38b-02ff-c012-e0bb-5cf53c817cc3.htm "BuildReadIndividualCommand 方法 (Byte, String)")

[BuildReadIndividualCommand 方法 (Byte, String[])](../html/9747205d-7347-ab27-dcf9-fab8e8efbe9a.htm "BuildReadIndividualCommand 方法 (Byte, String[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| LSCnetHelperBuildReadIndividualCommand 方法 (Byte, String) |

Multi reading address Type of Read Individual

**命名空间：**
 [HslCommunication.Profinet.LSIS.Helper](0200f8be-eba3-dea9-4eb6-b51d8c230c6b.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult<byte[]> BuildReadIndividualCommand(
	byte station,
	string address
)
```

```
Public Shared Function BuildReadIndividualCommand ( 
	station As Byte,
	address As String
) As OperateResult(Of Byte())
```

```
public:
static OperateResult<array<unsigned char>^>^ BuildReadIndividualCommand(
	unsigned char station, 
	String^ address
)
```

```
static member BuildReadIndividualCommand : 
        station : byte * 
        address : string -> OperateResult<byte[]> 
```

#### 参数

station
:   类型：SystemByte  
    plc station

address
:   类型：SystemString  

    [缺少 "M:HslCommunication.Profinet.LSIS.Helper.LSCnetHelper.BuildReadIndividualCommand(System.Byte,System.String)" 的 <param name="address"/> 文档]

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  

[缺少 "M:HslCommunication.Profinet.LSIS.Helper.LSCnetHelper.BuildReadIndividualCommand(System.Byte,System.String)" 的 <returns> 文档]

![](../icons/SectionExpanded.png)参见

#### 引用

[LSCnetHelper 类](c874de34-e25a-6c9c-8d0e-ddffe0651200.htm)

[BuildReadIndividualCommand 重载](883ada55-fd41-162f-09a9-2d34f92ac605.htm)

[HslCommunication.Profinet.LSIS.Helper 命名空间](0200f8be-eba3-dea9-4eb6-b51d8c230c6b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## BuildReadIndividualCommand 方法 (Byte, String[])

[原文連結](http://api.hslcommunication.cn/html/9747205d-7347-ab27-dcf9-fab8e8efbe9a.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.LSIS.Helper](../html/0200f8be-eba3-dea9-4eb6-b51d8c230c6b.htm "HslCommunication.Profinet.LSIS.Helper")

[LSCnetHelper 类](../html/c874de34-e25a-6c9c-8d0e-ddffe0651200.htm "LSCnetHelper 类")

[LSCnetHelper 方法](../html/d8e60529-bf76-1a12-c052-e3bfeaa4ec33.htm "LSCnetHelper 方法")

[BuildReadIndividualCommand 方法](../html/883ada55-fd41-162f-09a9-2d34f92ac605.htm "BuildReadIndividualCommand 方法 ")

[BuildReadIndividualCommand 方法 (Byte, String)](../html/6804f38b-02ff-c012-e0bb-5cf53c817cc3.htm "BuildReadIndividualCommand 方法 (Byte, String)")

[BuildReadIndividualCommand 方法 (Byte, String[])](../html/9747205d-7347-ab27-dcf9-fab8e8efbe9a.htm "BuildReadIndividualCommand 方法 (Byte, String[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| LSCnetHelperBuildReadIndividualCommand 方法 (Byte, String) |

Multi reading address Type of Read Individual

**命名空间：**
 [HslCommunication.Profinet.LSIS.Helper](0200f8be-eba3-dea9-4eb6-b51d8c230c6b.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult<byte[]> BuildReadIndividualCommand(
	byte station,
	string[] addresses
)
```

```
Public Shared Function BuildReadIndividualCommand ( 
	station As Byte,
	addresses As String()
) As OperateResult(Of Byte())
```

```
public:
static OperateResult<array<unsigned char>^>^ BuildReadIndividualCommand(
	unsigned char station, 
	array<String^>^ addresses
)
```

```
static member BuildReadIndividualCommand : 
        station : byte * 
        addresses : string[] -> OperateResult<byte[]> 
```

#### 参数

station
:   类型：SystemByte  
    plc station

addresses
:   类型：SystemString  
    address, for example: MX100, PX100

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  

[缺少 "M:HslCommunication.Profinet.LSIS.Helper.LSCnetHelper.BuildReadIndividualCommand(System.Byte,System.String[])" 的 <returns> 文档]

![](../icons/SectionExpanded.png)参见

#### 引用

[LSCnetHelper 类](c874de34-e25a-6c9c-8d0e-ddffe0651200.htm)

[BuildReadIndividualCommand 重载](883ada55-fd41-162f-09a9-2d34f92ac605.htm)

[HslCommunication.Profinet.LSIS.Helper 命名空间](0200f8be-eba3-dea9-4eb6-b51d8c230c6b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## BuildWriteByteCommand 方法 

[原文連結](http://api.hslcommunication.cn/html/62a3477c-e064-8c71-9687-4f1cc45ebc16.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.LSIS.Helper](../html/0200f8be-eba3-dea9-4eb6-b51d8c230c6b.htm "HslCommunication.Profinet.LSIS.Helper")

[LSCnetHelper 类](../html/c874de34-e25a-6c9c-8d0e-ddffe0651200.htm "LSCnetHelper 类")

[LSCnetHelper 方法](../html/d8e60529-bf76-1a12-c052-e3bfeaa4ec33.htm "LSCnetHelper 方法")

[AnalysisAddress 方法](../html/15ee9fd3-6c4b-6775-db86-370121c062ff.htm "AnalysisAddress 方法 ")

[BuildReadByteCommand 方法](../html/187eecaa-60f5-ed67-f0c6-ce15a59e1e52.htm "BuildReadByteCommand 方法 ")

[BuildReadIndividualCommand 方法](../html/883ada55-fd41-162f-09a9-2d34f92ac605.htm "BuildReadIndividualCommand 方法 ")

[BuildWriteByteCommand 方法](../html/62a3477c-e064-8c71-9687-4f1cc45ebc16.htm "BuildWriteByteCommand 方法 ")

[BuildWriteCommand 方法](../html/2873d537-6750-08f6-ac89-0bf2ea205de3.htm "BuildWriteCommand 方法 ")

[BuildWriteOneCommand 方法](../html/c59dc3f6-675c-8f9d-9dce-3c0857887b3b.htm "BuildWriteOneCommand 方法 ")

[GetAddressOfU\_Q\_I 方法](../html/c3212593-efad-0a6e-22a6-11be146a83d9.htm "GetAddressOfU_Q_I 方法 ")

[GetErrorText 方法](../html/d8d3953e-a826-5558-6773-8809cf6d9dbe.htm "GetErrorText 方法 ")

[Read 方法](../html/7c298c2b-c4d6-de5b-8b67-1ddb7fb1f180.htm "Read 方法 ")

[ReadAsync 方法](../html/f3d85a1f-1978-d789-f6a0-f1dd0fb17e8c.htm "ReadAsync 方法 ")

[ReadBool 方法](../html/3359728b-1b82-547f-13ea-05d524a7bacc.htm "ReadBool 方法 ")

[ReadBoolAsync 方法](../html/6b11ae2c-5cb7-a003-2852-48d340e12c62.htm "ReadBoolAsync 方法 ")

[UnpackResponseContent 方法](../html/ef40b9f3-0fc8-b1dd-945b-8d2011f1b6a9.htm "UnpackResponseContent 方法 ")

[Write 方法](../html/2792b231-69a8-6fae-dc81-e778e73172e1.htm "Write 方法 ")

[WriteAsync 方法](../html/a8b58372-0a53-0c10-6340-30584a178a1b.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| LSCnetHelperBuildWriteByteCommand 方法 |

write data to address Type of ReadByte

**命名空间：**
 [HslCommunication.Profinet.LSIS.Helper](0200f8be-eba3-dea9-4eb6-b51d8c230c6b.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult<byte[]> BuildWriteByteCommand(
	byte station,
	string address,
	byte[] value
)
```

```
Public Shared Function BuildWriteByteCommand ( 
	station As Byte,
	address As String,
	value As Byte()
) As OperateResult(Of Byte())
```

```
public:
static OperateResult<array<unsigned char>^>^ BuildWriteByteCommand(
	unsigned char station, 
	String^ address, 
	array<unsigned char>^ value
)
```

```
static member BuildWriteByteCommand : 
        station : byte * 
        address : string * 
        value : byte[] -> OperateResult<byte[]> 
```

#### 参数

station
:   类型：SystemByte  
    plc station

address
:   类型：SystemString  
    address, for example: M100, D100, DW100

value
:   类型：SystemByte  
    source value

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
command bytes

![](../icons/SectionExpanded.png)参见

#### 引用

[LSCnetHelper 类](c874de34-e25a-6c9c-8d0e-ddffe0651200.htm)

[HslCommunication.Profinet.LSIS.Helper 命名空间](0200f8be-eba3-dea9-4eb6-b51d8c230c6b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## BuildWriteCommand 方法 

[原文連結](http://api.hslcommunication.cn/html/2873d537-6750-08f6-ac89-0bf2ea205de3.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.LSIS.Helper](../html/0200f8be-eba3-dea9-4eb6-b51d8c230c6b.htm "HslCommunication.Profinet.LSIS.Helper")

[LSCnetHelper 类](../html/c874de34-e25a-6c9c-8d0e-ddffe0651200.htm "LSCnetHelper 类")

[LSCnetHelper 方法](../html/d8e60529-bf76-1a12-c052-e3bfeaa4ec33.htm "LSCnetHelper 方法")

[AnalysisAddress 方法](../html/15ee9fd3-6c4b-6775-db86-370121c062ff.htm "AnalysisAddress 方法 ")

[BuildReadByteCommand 方法](../html/187eecaa-60f5-ed67-f0c6-ce15a59e1e52.htm "BuildReadByteCommand 方法 ")

[BuildReadIndividualCommand 方法](../html/883ada55-fd41-162f-09a9-2d34f92ac605.htm "BuildReadIndividualCommand 方法 ")

[BuildWriteByteCommand 方法](../html/62a3477c-e064-8c71-9687-4f1cc45ebc16.htm "BuildWriteByteCommand 方法 ")

[BuildWriteCommand 方法](../html/2873d537-6750-08f6-ac89-0bf2ea205de3.htm "BuildWriteCommand 方法 ")

[BuildWriteOneCommand 方法](../html/c59dc3f6-675c-8f9d-9dce-3c0857887b3b.htm "BuildWriteOneCommand 方法 ")

[GetAddressOfU\_Q\_I 方法](../html/c3212593-efad-0a6e-22a6-11be146a83d9.htm "GetAddressOfU_Q_I 方法 ")

[GetErrorText 方法](../html/d8d3953e-a826-5558-6773-8809cf6d9dbe.htm "GetErrorText 方法 ")

[Read 方法](../html/7c298c2b-c4d6-de5b-8b67-1ddb7fb1f180.htm "Read 方法 ")

[ReadAsync 方法](../html/f3d85a1f-1978-d789-f6a0-f1dd0fb17e8c.htm "ReadAsync 方法 ")

[ReadBool 方法](../html/3359728b-1b82-547f-13ea-05d524a7bacc.htm "ReadBool 方法 ")

[ReadBoolAsync 方法](../html/6b11ae2c-5cb7-a003-2852-48d340e12c62.htm "ReadBoolAsync 方法 ")

[UnpackResponseContent 方法](../html/ef40b9f3-0fc8-b1dd-945b-8d2011f1b6a9.htm "UnpackResponseContent 方法 ")

[Write 方法](../html/2792b231-69a8-6fae-dc81-e778e73172e1.htm "Write 方法 ")

[WriteAsync 方法](../html/a8b58372-0a53-0c10-6340-30584a178a1b.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| LSCnetHelperBuildWriteCommand 方法 |

write data to address Type of ReadByte

**命名空间：**
 [HslCommunication.Profinet.LSIS.Helper](0200f8be-eba3-dea9-4eb6-b51d8c230c6b.htm)  
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
	byte[] value
)
```

```
Public Shared Function BuildWriteCommand ( 
	station As Byte,
	address As String,
	value As Byte()
) As OperateResult(Of Byte())
```

```
public:
static OperateResult<array<unsigned char>^>^ BuildWriteCommand(
	unsigned char station, 
	String^ address, 
	array<unsigned char>^ value
)
```

```
static member BuildWriteCommand : 
        station : byte * 
        address : string * 
        value : byte[] -> OperateResult<byte[]> 
```

#### 参数

station
:   类型：SystemByte  
    plc station

address
:   类型：SystemString  
    address, for example: M100, D100, DW100

value
:   类型：SystemByte  
    source value

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
command bytes

![](../icons/SectionExpanded.png)参见

#### 引用

[LSCnetHelper 类](c874de34-e25a-6c9c-8d0e-ddffe0651200.htm)

[HslCommunication.Profinet.LSIS.Helper 命名空间](0200f8be-eba3-dea9-4eb6-b51d8c230c6b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## BuildWriteOneCommand 方法 

[原文連結](http://api.hslcommunication.cn/html/c59dc3f6-675c-8f9d-9dce-3c0857887b3b.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.LSIS.Helper](../html/0200f8be-eba3-dea9-4eb6-b51d8c230c6b.htm "HslCommunication.Profinet.LSIS.Helper")

[LSCnetHelper 类](../html/c874de34-e25a-6c9c-8d0e-ddffe0651200.htm "LSCnetHelper 类")

[LSCnetHelper 方法](../html/d8e60529-bf76-1a12-c052-e3bfeaa4ec33.htm "LSCnetHelper 方法")

[AnalysisAddress 方法](../html/15ee9fd3-6c4b-6775-db86-370121c062ff.htm "AnalysisAddress 方法 ")

[BuildReadByteCommand 方法](../html/187eecaa-60f5-ed67-f0c6-ce15a59e1e52.htm "BuildReadByteCommand 方法 ")

[BuildReadIndividualCommand 方法](../html/883ada55-fd41-162f-09a9-2d34f92ac605.htm "BuildReadIndividualCommand 方法 ")

[BuildWriteByteCommand 方法](../html/62a3477c-e064-8c71-9687-4f1cc45ebc16.htm "BuildWriteByteCommand 方法 ")

[BuildWriteCommand 方法](../html/2873d537-6750-08f6-ac89-0bf2ea205de3.htm "BuildWriteCommand 方法 ")

[BuildWriteOneCommand 方法](../html/c59dc3f6-675c-8f9d-9dce-3c0857887b3b.htm "BuildWriteOneCommand 方法 ")

[GetAddressOfU\_Q\_I 方法](../html/c3212593-efad-0a6e-22a6-11be146a83d9.htm "GetAddressOfU_Q_I 方法 ")

[GetErrorText 方法](../html/d8d3953e-a826-5558-6773-8809cf6d9dbe.htm "GetErrorText 方法 ")

[Read 方法](../html/7c298c2b-c4d6-de5b-8b67-1ddb7fb1f180.htm "Read 方法 ")

[ReadAsync 方法](../html/f3d85a1f-1978-d789-f6a0-f1dd0fb17e8c.htm "ReadAsync 方法 ")

[ReadBool 方法](../html/3359728b-1b82-547f-13ea-05d524a7bacc.htm "ReadBool 方法 ")

[ReadBoolAsync 方法](../html/6b11ae2c-5cb7-a003-2852-48d340e12c62.htm "ReadBoolAsync 方法 ")

[UnpackResponseContent 方法](../html/ef40b9f3-0fc8-b1dd-945b-8d2011f1b6a9.htm "UnpackResponseContent 方法 ")

[Write 方法](../html/2792b231-69a8-6fae-dc81-e778e73172e1.htm "Write 方法 ")

[WriteAsync 方法](../html/a8b58372-0a53-0c10-6340-30584a178a1b.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| LSCnetHelperBuildWriteOneCommand 方法 |

write data to address Type of One

**命名空间：**
 [HslCommunication.Profinet.LSIS.Helper](0200f8be-eba3-dea9-4eb6-b51d8c230c6b.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult<byte[]> BuildWriteOneCommand(
	byte station,
	string address,
	byte[] value
)
```

```
Public Shared Function BuildWriteOneCommand ( 
	station As Byte,
	address As String,
	value As Byte()
) As OperateResult(Of Byte())
```

```
public:
static OperateResult<array<unsigned char>^>^ BuildWriteOneCommand(
	unsigned char station, 
	String^ address, 
	array<unsigned char>^ value
)
```

```
static member BuildWriteOneCommand : 
        station : byte * 
        address : string * 
        value : byte[] -> OperateResult<byte[]> 
```

#### 参数

station
:   类型：SystemByte  
    plc station

address
:   类型：SystemString  
    address, for example: M100, D100, DW100

value
:   类型：SystemByte  
    source value

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
command bytes

![](../icons/SectionExpanded.png)参见

#### 引用

[LSCnetHelper 类](c874de34-e25a-6c9c-8d0e-ddffe0651200.htm)

[HslCommunication.Profinet.LSIS.Helper 命名空间](0200f8be-eba3-dea9-4eb6-b51d8c230c6b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GetAddressOfU_Q_I 方法 

[原文連結](http://api.hslcommunication.cn/html/c3212593-efad-0a6e-22a6-11be146a83d9.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.LSIS.Helper](../html/0200f8be-eba3-dea9-4eb6-b51d8c230c6b.htm "HslCommunication.Profinet.LSIS.Helper")

[LSCnetHelper 类](../html/c874de34-e25a-6c9c-8d0e-ddffe0651200.htm "LSCnetHelper 类")

[LSCnetHelper 方法](../html/d8e60529-bf76-1a12-c052-e3bfeaa4ec33.htm "LSCnetHelper 方法")

[AnalysisAddress 方法](../html/15ee9fd3-6c4b-6775-db86-370121c062ff.htm "AnalysisAddress 方法 ")

[BuildReadByteCommand 方法](../html/187eecaa-60f5-ed67-f0c6-ce15a59e1e52.htm "BuildReadByteCommand 方法 ")

[BuildReadIndividualCommand 方法](../html/883ada55-fd41-162f-09a9-2d34f92ac605.htm "BuildReadIndividualCommand 方法 ")

[BuildWriteByteCommand 方法](../html/62a3477c-e064-8c71-9687-4f1cc45ebc16.htm "BuildWriteByteCommand 方法 ")

[BuildWriteCommand 方法](../html/2873d537-6750-08f6-ac89-0bf2ea205de3.htm "BuildWriteCommand 方法 ")

[BuildWriteOneCommand 方法](../html/c59dc3f6-675c-8f9d-9dce-3c0857887b3b.htm "BuildWriteOneCommand 方法 ")

[GetAddressOfU\_Q\_I 方法](../html/c3212593-efad-0a6e-22a6-11be146a83d9.htm "GetAddressOfU_Q_I 方法 ")

[GetErrorText 方法](../html/d8d3953e-a826-5558-6773-8809cf6d9dbe.htm "GetErrorText 方法 ")

[Read 方法](../html/7c298c2b-c4d6-de5b-8b67-1ddb7fb1f180.htm "Read 方法 ")

[ReadAsync 方法](../html/f3d85a1f-1978-d789-f6a0-f1dd0fb17e8c.htm "ReadAsync 方法 ")

[ReadBool 方法](../html/3359728b-1b82-547f-13ea-05d524a7bacc.htm "ReadBool 方法 ")

[ReadBoolAsync 方法](../html/6b11ae2c-5cb7-a003-2852-48d340e12c62.htm "ReadBoolAsync 方法 ")

[UnpackResponseContent 方法](../html/ef40b9f3-0fc8-b1dd-945b-8d2011f1b6a9.htm "UnpackResponseContent 方法 ")

[Write 方法](../html/2792b231-69a8-6fae-dc81-e778e73172e1.htm "Write 方法 ")

[WriteAsync 方法](../html/a8b58372-0a53-0c10-6340-30584a178a1b.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| LSCnetHelperGetAddressOfU\_Q\_I 方法 |

[缺少 "M:HslCommunication.Profinet.LSIS.Helper.LSCnetHelper.GetAddressOfU\_Q\_I(System.String,System.Boolean)" 的 <summary> 文档]

**命名空间：**
 [HslCommunication.Profinet.LSIS.Helper](0200f8be-eba3-dea9-4eb6-b51d8c230c6b.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static string GetAddressOfU_Q_I(
	string Address,
	bool IsWrite = false
)
```

```
Public Shared Function GetAddressOfU_Q_I ( 
	Address As String,
	Optional IsWrite As Boolean = false
) As String
```

```
public:
static String^ GetAddressOfU_Q_I(
	String^ Address, 
	bool IsWrite = false
)
```

```
static member GetAddressOfU_Q_I : 
        Address : string * 
        ?IsWrite : bool 
(* Defaults:
        let _IsWrite = defaultArg IsWrite false
*)
-> string 
```

#### 参数

Address
:   类型：SystemString  

    [缺少 "M:HslCommunication.Profinet.LSIS.Helper.LSCnetHelper.GetAddressOfU\_Q\_I(System.String,System.Boolean)" 的 <param name="Address"/> 文档]

IsWrite (Optional)
:   类型：SystemBoolean  

    [缺少 "M:HslCommunication.Profinet.LSIS.Helper.LSCnetHelper.GetAddressOfU\_Q\_I(System.String,System.Boolean)" 的 <param name="IsWrite"/> 文档]

#### 返回值

类型：String  

[缺少 "M:HslCommunication.Profinet.LSIS.Helper.LSCnetHelper.GetAddressOfU\_Q\_I(System.String,System.Boolean)" 的 <returns> 文档]

![](../icons/SectionExpanded.png)参见

#### 引用

[LSCnetHelper 类](c874de34-e25a-6c9c-8d0e-ddffe0651200.htm)

[HslCommunication.Profinet.LSIS.Helper 命名空间](0200f8be-eba3-dea9-4eb6-b51d8c230c6b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GetErrorText 方法 

[原文連結](http://api.hslcommunication.cn/html/d8d3953e-a826-5558-6773-8809cf6d9dbe.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.LSIS.Helper](../html/0200f8be-eba3-dea9-4eb6-b51d8c230c6b.htm "HslCommunication.Profinet.LSIS.Helper")

[LSCnetHelper 类](../html/c874de34-e25a-6c9c-8d0e-ddffe0651200.htm "LSCnetHelper 类")

[LSCnetHelper 方法](../html/d8e60529-bf76-1a12-c052-e3bfeaa4ec33.htm "LSCnetHelper 方法")

[AnalysisAddress 方法](../html/15ee9fd3-6c4b-6775-db86-370121c062ff.htm "AnalysisAddress 方法 ")

[BuildReadByteCommand 方法](../html/187eecaa-60f5-ed67-f0c6-ce15a59e1e52.htm "BuildReadByteCommand 方法 ")

[BuildReadIndividualCommand 方法](../html/883ada55-fd41-162f-09a9-2d34f92ac605.htm "BuildReadIndividualCommand 方法 ")

[BuildWriteByteCommand 方法](../html/62a3477c-e064-8c71-9687-4f1cc45ebc16.htm "BuildWriteByteCommand 方法 ")

[BuildWriteCommand 方法](../html/2873d537-6750-08f6-ac89-0bf2ea205de3.htm "BuildWriteCommand 方法 ")

[BuildWriteOneCommand 方法](../html/c59dc3f6-675c-8f9d-9dce-3c0857887b3b.htm "BuildWriteOneCommand 方法 ")

[GetAddressOfU\_Q\_I 方法](../html/c3212593-efad-0a6e-22a6-11be146a83d9.htm "GetAddressOfU_Q_I 方法 ")

[GetErrorText 方法](../html/d8d3953e-a826-5558-6773-8809cf6d9dbe.htm "GetErrorText 方法 ")

[Read 方法](../html/7c298c2b-c4d6-de5b-8b67-1ddb7fb1f180.htm "Read 方法 ")

[ReadAsync 方法](../html/f3d85a1f-1978-d789-f6a0-f1dd0fb17e8c.htm "ReadAsync 方法 ")

[ReadBool 方法](../html/3359728b-1b82-547f-13ea-05d524a7bacc.htm "ReadBool 方法 ")

[ReadBoolAsync 方法](../html/6b11ae2c-5cb7-a003-2852-48d340e12c62.htm "ReadBoolAsync 方法 ")

[UnpackResponseContent 方法](../html/ef40b9f3-0fc8-b1dd-945b-8d2011f1b6a9.htm "UnpackResponseContent 方法 ")

[Write 方法](../html/2792b231-69a8-6fae-dc81-e778e73172e1.htm "Write 方法 ")

[WriteAsync 方法](../html/a8b58372-0a53-0c10-6340-30584a178a1b.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| LSCnetHelperGetErrorText 方法 |

根据错误号，获取到真实的错误描述信息  
According to the error number, get the real error description information

**命名空间：**
 [HslCommunication.Profinet.LSIS.Helper](0200f8be-eba3-dea9-4eb6-b51d8c230c6b.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static string GetErrorText(
	int err
)
```

```
Public Shared Function GetErrorText ( 
	err As Integer
) As String
```

```
public:
static String^ GetErrorText(
	int err
)
```

```
static member GetErrorText : 
        err : int -> string 
```

#### 参数

err
:   类型：SystemInt32  
    错误号

#### 返回值

类型：String  
真实的错误描述信息

![](../icons/SectionExpanded.png)参见

#### 引用

[LSCnetHelper 类](c874de34-e25a-6c9c-8d0e-ddffe0651200.htm)

[HslCommunication.Profinet.LSIS.Helper 命名空间](0200f8be-eba3-dea9-4eb6-b51d8c230c6b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Read 方法 

[原文連結](http://api.hslcommunication.cn/html/7c298c2b-c4d6-de5b-8b67-1ddb7fb1f180.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.LSIS.Helper](../html/0200f8be-eba3-dea9-4eb6-b51d8c230c6b.htm "HslCommunication.Profinet.LSIS.Helper")

[LSCnetHelper 类](../html/c874de34-e25a-6c9c-8d0e-ddffe0651200.htm "LSCnetHelper 类")

[LSCnetHelper 方法](../html/d8e60529-bf76-1a12-c052-e3bfeaa4ec33.htm "LSCnetHelper 方法")

[Read 方法](../html/7c298c2b-c4d6-de5b-8b67-1ddb7fb1f180.htm "Read 方法 ")

[Read 方法 (IReadWriteDevice, Int32, String[])](../html/5af4122b-f6be-660b-080d-c930b67f90e9.htm "Read 方法 (IReadWriteDevice, Int32, String[])")

[Read 方法 (IReadWriteDeviceStation, String, UInt16)](../html/c2735b0c-d09d-bfee-e23b-628145ca9709.htm "Read 方法 (IReadWriteDeviceStation, String, UInt16)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| LSCnetHelperRead 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法静态成员 | [Read(IReadWriteDevice, Int32, String)](5af4122b-f6be-660b-080d-c930b67f90e9.htm) | 从PLC设备读取多个地址的数据信息，返回连续的字节数组，需要按照实际情况进行按顺序解析。  Read the data information of multiple addresses from the PLC device and return a continuous byte array, which needs to be parsed in order according to the actual situation. |
| 公共方法静态成员 | [Read(IReadWriteDeviceStation, String, UInt16)](c2735b0c-d09d-bfee-e23b-628145ca9709.htm) | 从PLC的指定地址读取原始的字节数据信息，地址示例：MB100, MW100, MD100, 如果输入了M100等同于MB100  Read the original byte data information from the designated address of the PLC. Examples of addresses: MB100, MW100, MD100, if the input M100 is equivalent to MB100 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[LSCnetHelper 类](c874de34-e25a-6c9c-8d0e-ddffe0651200.htm)

[HslCommunication.Profinet.LSIS.Helper 命名空间](0200f8be-eba3-dea9-4eb6-b51d8c230c6b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Read 方法 (IReadWriteDevice, Int32, String[])

[原文連結](http://api.hslcommunication.cn/html/5af4122b-f6be-660b-080d-c930b67f90e9.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.LSIS.Helper](../html/0200f8be-eba3-dea9-4eb6-b51d8c230c6b.htm "HslCommunication.Profinet.LSIS.Helper")

[LSCnetHelper 类](../html/c874de34-e25a-6c9c-8d0e-ddffe0651200.htm "LSCnetHelper 类")

[LSCnetHelper 方法](../html/d8e60529-bf76-1a12-c052-e3bfeaa4ec33.htm "LSCnetHelper 方法")

[Read 方法](../html/7c298c2b-c4d6-de5b-8b67-1ddb7fb1f180.htm "Read 方法 ")

[Read 方法 (IReadWriteDevice, Int32, String[])](../html/5af4122b-f6be-660b-080d-c930b67f90e9.htm "Read 方法 (IReadWriteDevice, Int32, String[])")

[Read 方法 (IReadWriteDeviceStation, String, UInt16)](../html/c2735b0c-d09d-bfee-e23b-628145ca9709.htm "Read 方法 (IReadWriteDeviceStation, String, UInt16)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| LSCnetHelperRead 方法 (IReadWriteDevice, Int32, String) |

从PLC设备读取多个地址的数据信息，返回连续的字节数组，需要按照实际情况进行按顺序解析。  
Read the data information of multiple addresses from the PLC device and return a continuous byte array, which needs to be parsed in order according to the actual situation.

**命名空间：**
 [HslCommunication.Profinet.LSIS.Helper](0200f8be-eba3-dea9-4eb6-b51d8c230c6b.htm)  
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
	int station,
	string[] address
)
```

```
Public Shared Function Read ( 
	plc As IReadWriteDevice,
	station As Integer,
	address As String()
) As OperateResult(Of Byte())
```

```
public:
static OperateResult<array<unsigned char>^>^ Read(
	IReadWriteDevice^ plc, 
	int station, 
	array<String^>^ address
)
```

```
static member Read : 
        plc : IReadWriteDevice * 
        station : int * 
        address : string[] -> OperateResult<byte[]> 
```

#### 参数

plc
:   类型：[HslCommunication.CoreIReadWriteDevice](d987afcc-f99c-6aed-3847-5b1890504635.htm)  
    PLC通信对象

station
:   类型：SystemInt32  
    站号信息

address
:   类型：SystemString  
    PLC的地址信息，例如 M100, MB100, MW100, MD100

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
结果对象数据

![](../icons/SectionExpanded.png)备注

按照每16个地址长度进行自动的切割，支持任意的多的长度地址

![](../icons/SectionExpanded.png)参见

#### 引用

[LSCnetHelper 类](c874de34-e25a-6c9c-8d0e-ddffe0651200.htm)

[Read 重载](7c298c2b-c4d6-de5b-8b67-1ddb7fb1f180.htm)

[HslCommunication.Profinet.LSIS.Helper 命名空间](0200f8be-eba3-dea9-4eb6-b51d8c230c6b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Read 方法 (IReadWriteDeviceStation, String, UInt16)

[原文連結](http://api.hslcommunication.cn/html/c2735b0c-d09d-bfee-e23b-628145ca9709.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.LSIS.Helper](../html/0200f8be-eba3-dea9-4eb6-b51d8c230c6b.htm "HslCommunication.Profinet.LSIS.Helper")

[LSCnetHelper 类](../html/c874de34-e25a-6c9c-8d0e-ddffe0651200.htm "LSCnetHelper 类")

[LSCnetHelper 方法](../html/d8e60529-bf76-1a12-c052-e3bfeaa4ec33.htm "LSCnetHelper 方法")

[Read 方法](../html/7c298c2b-c4d6-de5b-8b67-1ddb7fb1f180.htm "Read 方法 ")

[Read 方法 (IReadWriteDevice, Int32, String[])](../html/5af4122b-f6be-660b-080d-c930b67f90e9.htm "Read 方法 (IReadWriteDevice, Int32, String[])")

[Read 方法 (IReadWriteDeviceStation, String, UInt16)](../html/c2735b0c-d09d-bfee-e23b-628145ca9709.htm "Read 方法 (IReadWriteDeviceStation, String, UInt16)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| LSCnetHelperRead 方法 (IReadWriteDeviceStation, String, UInt16) |

从PLC的指定地址读取原始的字节数据信息，地址示例：MB100, MW100, MD100, 如果输入了M100等同于MB100  
Read the original byte data information from the designated address of the PLC.
Examples of addresses: MB100, MW100, MD100, if the input M100 is equivalent to MB100

**命名空间：**
 [HslCommunication.Profinet.LSIS.Helper](0200f8be-eba3-dea9-4eb6-b51d8c230c6b.htm)  
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
	IReadWriteDeviceStation plc,
	string address,
	ushort length
)
```

```
Public Shared Function Read ( 
	plc As IReadWriteDeviceStation,
	address As String,
	length As UShort
) As OperateResult(Of Byte())
```

```
public:
static OperateResult<array<unsigned char>^>^ Read(
	IReadWriteDeviceStation^ plc, 
	String^ address, 
	unsigned short length
)
```

```
static member Read : 
        plc : IReadWriteDeviceStation * 
        address : string * 
        length : uint16 -> OperateResult<byte[]> 
```

#### 参数

plc
:   类型：[HslCommunication.Core.NetIReadWriteDeviceStation](0a167412-aa6f-f505-8c8e-75fa32832055.htm)  
    PLC通信对象

address
:   类型：SystemString  
    PLC的地址信息，例如 M100, MB100, MW100, MD100

length
:   类型：SystemUInt16  
    读取的长度信息

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
返回是否读取成功的结果对象

![](../icons/SectionExpanded.png)备注

地址类型支持 P,M,L,K,F,T,C,D,R,I,Q,W, 支持携带站号的形式，例如 s=2;MW100

![](../icons/SectionExpanded.png)参见

#### 引用

[LSCnetHelper 类](c874de34-e25a-6c9c-8d0e-ddffe0651200.htm)

[Read 重载](7c298c2b-c4d6-de5b-8b67-1ddb7fb1f180.htm)

[HslCommunication.Profinet.LSIS.Helper 命名空间](0200f8be-eba3-dea9-4eb6-b51d8c230c6b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/f3d85a1f-1978-d789-f6a0-f1dd0fb17e8c.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.LSIS.Helper](../html/0200f8be-eba3-dea9-4eb6-b51d8c230c6b.htm "HslCommunication.Profinet.LSIS.Helper")

[LSCnetHelper 类](../html/c874de34-e25a-6c9c-8d0e-ddffe0651200.htm "LSCnetHelper 类")

[LSCnetHelper 方法](../html/d8e60529-bf76-1a12-c052-e3bfeaa4ec33.htm "LSCnetHelper 方法")

[ReadAsync 方法](../html/f3d85a1f-1978-d789-f6a0-f1dd0fb17e8c.htm "ReadAsync 方法 ")

[ReadAsync 方法 (IReadWriteDevice, Int32, String[])](../html/c685c2aa-fd95-7051-a432-edf9d1eff0a5.htm "ReadAsync 方法 (IReadWriteDevice, Int32, String[])")

[ReadAsync 方法 (IReadWriteDeviceStation, String, UInt16)](../html/3c28f44d-b751-2afc-e6c4-0a66a210160b.htm "ReadAsync 方法 (IReadWriteDeviceStation, String, UInt16)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| LSCnetHelperReadAsync 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法静态成员 | [ReadAsync(IReadWriteDevice, Int32, String)](c685c2aa-fd95-7051-a432-edf9d1eff0a5.htm) | 从PLC设备读取多个地址的数据信息，返回连续的字节数组，需要按照实际情况进行按顺序解析。  Read the data information of multiple addresses from the PLC device and return a continuous byte array, which needs to be parsed in order according to the actual situation. |
| 公共方法静态成员 | [ReadAsync(IReadWriteDeviceStation, String, UInt16)](3c28f44d-b751-2afc-e6c4-0a66a210160b.htm) | 从PLC的指定地址读取原始的字节数据信息，地址示例：MB100, MW100, MD100, 如果输入了M100等同于MB100  Read the original byte data information from the designated address of the PLC. Examples of addresses: MB100, MW100, MD100, if the input M100 is equivalent to MB100 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[LSCnetHelper 类](c874de34-e25a-6c9c-8d0e-ddffe0651200.htm)

[HslCommunication.Profinet.LSIS.Helper 命名空间](0200f8be-eba3-dea9-4eb6-b51d8c230c6b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadAsync 方法 (IReadWriteDevice, Int32, String[])

[原文連結](http://api.hslcommunication.cn/html/c685c2aa-fd95-7051-a432-edf9d1eff0a5.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.LSIS.Helper](../html/0200f8be-eba3-dea9-4eb6-b51d8c230c6b.htm "HslCommunication.Profinet.LSIS.Helper")

[LSCnetHelper 类](../html/c874de34-e25a-6c9c-8d0e-ddffe0651200.htm "LSCnetHelper 类")

[LSCnetHelper 方法](../html/d8e60529-bf76-1a12-c052-e3bfeaa4ec33.htm "LSCnetHelper 方法")

[ReadAsync 方法](../html/f3d85a1f-1978-d789-f6a0-f1dd0fb17e8c.htm "ReadAsync 方法 ")

[ReadAsync 方法 (IReadWriteDevice, Int32, String[])](../html/c685c2aa-fd95-7051-a432-edf9d1eff0a5.htm "ReadAsync 方法 (IReadWriteDevice, Int32, String[])")

[ReadAsync 方法 (IReadWriteDeviceStation, String, UInt16)](../html/3c28f44d-b751-2afc-e6c4-0a66a210160b.htm "ReadAsync 方法 (IReadWriteDeviceStation, String, UInt16)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| LSCnetHelperReadAsync 方法 (IReadWriteDevice, Int32, String) |

从PLC设备读取多个地址的数据信息，返回连续的字节数组，需要按照实际情况进行按顺序解析。  
Read the data information of multiple addresses from the PLC device and return a continuous byte array, which needs to be parsed in order according to the actual situation.

**命名空间：**
 [HslCommunication.Profinet.LSIS.Helper](0200f8be-eba3-dea9-4eb6-b51d8c230c6b.htm)  
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
	IReadWriteDevice plc,
	int station,
	string[] address
)
```

```
Public Shared Function ReadAsync ( 
	plc As IReadWriteDevice,
	station As Integer,
	address As String()
) As Task(Of OperateResult(Of Byte()))
```

```
public:
static Task<OperateResult<array<unsigned char>^>^>^ ReadAsync(
	IReadWriteDevice^ plc, 
	int station, 
	array<String^>^ address
)
```

```
static member ReadAsync : 
        plc : IReadWriteDevice * 
        station : int * 
        address : string[] -> Task<OperateResult<byte[]>> 
```

#### 参数

plc
:   类型：[HslCommunication.CoreIReadWriteDevice](d987afcc-f99c-6aed-3847-5b1890504635.htm)  
    PLC通信对象

station
:   类型：SystemInt32  
    站号信息

address
:   类型：SystemString  
    PLC的地址信息，例如 M100, MB100, MW100, MD100

#### 返回值

类型：Task[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
结果对象数据

![](../icons/SectionExpanded.png)备注

按照每16个地址长度进行自动的切割，支持任意的多的长度地址

![](../icons/SectionExpanded.png)参见

#### 引用

[LSCnetHelper 类](c874de34-e25a-6c9c-8d0e-ddffe0651200.htm)

[ReadAsync 重载](f3d85a1f-1978-d789-f6a0-f1dd0fb17e8c.htm)

[HslCommunication.Profinet.LSIS.Helper 命名空间](0200f8be-eba3-dea9-4eb6-b51d8c230c6b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadAsync 方法 (IReadWriteDeviceStation, String, UInt16)

[原文連結](http://api.hslcommunication.cn/html/3c28f44d-b751-2afc-e6c4-0a66a210160b.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.LSIS.Helper](../html/0200f8be-eba3-dea9-4eb6-b51d8c230c6b.htm "HslCommunication.Profinet.LSIS.Helper")

[LSCnetHelper 类](../html/c874de34-e25a-6c9c-8d0e-ddffe0651200.htm "LSCnetHelper 类")

[LSCnetHelper 方法](../html/d8e60529-bf76-1a12-c052-e3bfeaa4ec33.htm "LSCnetHelper 方法")

[ReadAsync 方法](../html/f3d85a1f-1978-d789-f6a0-f1dd0fb17e8c.htm "ReadAsync 方法 ")

[ReadAsync 方法 (IReadWriteDevice, Int32, String[])](../html/c685c2aa-fd95-7051-a432-edf9d1eff0a5.htm "ReadAsync 方法 (IReadWriteDevice, Int32, String[])")

[ReadAsync 方法 (IReadWriteDeviceStation, String, UInt16)](../html/3c28f44d-b751-2afc-e6c4-0a66a210160b.htm "ReadAsync 方法 (IReadWriteDeviceStation, String, UInt16)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| LSCnetHelperReadAsync 方法 (IReadWriteDeviceStation, String, UInt16) |

从PLC的指定地址读取原始的字节数据信息，地址示例：MB100, MW100, MD100, 如果输入了M100等同于MB100  
Read the original byte data information from the designated address of the PLC.
Examples of addresses: MB100, MW100, MD100, if the input M100 is equivalent to MB100

**命名空间：**
 [HslCommunication.Profinet.LSIS.Helper](0200f8be-eba3-dea9-4eb6-b51d8c230c6b.htm)  
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
	IReadWriteDeviceStation plc,
	string address,
	ushort length
)
```

```
Public Shared Function ReadAsync ( 
	plc As IReadWriteDeviceStation,
	address As String,
	length As UShort
) As Task(Of OperateResult(Of Byte()))
```

```
public:
static Task<OperateResult<array<unsigned char>^>^>^ ReadAsync(
	IReadWriteDeviceStation^ plc, 
	String^ address, 
	unsigned short length
)
```

```
static member ReadAsync : 
        plc : IReadWriteDeviceStation * 
        address : string * 
        length : uint16 -> Task<OperateResult<byte[]>> 
```

#### 参数

plc
:   类型：[HslCommunication.Core.NetIReadWriteDeviceStation](0a167412-aa6f-f505-8c8e-75fa32832055.htm)  
    PLC通信对象

address
:   类型：SystemString  
    PLC的地址信息，例如 M100, MB100, MW100, MD100

length
:   类型：SystemUInt16  
    读取的长度信息

#### 返回值

类型：Task[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
返回是否读取成功的结果对象

![](../icons/SectionExpanded.png)备注

地址类型支持 P,M,L,K,F,T,C,D,R,I,Q,W, 支持携带站号的形式，例如 s=2;MW100

![](../icons/SectionExpanded.png)参见

#### 引用

[LSCnetHelper 类](c874de34-e25a-6c9c-8d0e-ddffe0651200.htm)

[ReadAsync 重载](f3d85a1f-1978-d789-f6a0-f1dd0fb17e8c.htm)

[HslCommunication.Profinet.LSIS.Helper 命名空间](0200f8be-eba3-dea9-4eb6-b51d8c230c6b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadBool 方法 

[原文連結](http://api.hslcommunication.cn/html/3359728b-1b82-547f-13ea-05d524a7bacc.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.LSIS.Helper](../html/0200f8be-eba3-dea9-4eb6-b51d8c230c6b.htm "HslCommunication.Profinet.LSIS.Helper")

[LSCnetHelper 类](../html/c874de34-e25a-6c9c-8d0e-ddffe0651200.htm "LSCnetHelper 类")

[LSCnetHelper 方法](../html/d8e60529-bf76-1a12-c052-e3bfeaa4ec33.htm "LSCnetHelper 方法")

[ReadBool 方法](../html/3359728b-1b82-547f-13ea-05d524a7bacc.htm "ReadBool 方法 ")

[ReadBool 方法 (IReadWriteDevice, Int32, String)](../html/6c90e993-92e2-4b41-a94f-9233ae54991d.htm "ReadBool 方法 (IReadWriteDevice, Int32, String)")

[ReadBool 方法 (IReadWriteDeviceStation, String, UInt16)](../html/84c5849a-9b04-82f1-c992-c67d92c7218a.htm "ReadBool 方法 (IReadWriteDeviceStation, String, UInt16)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| LSCnetHelperReadBool 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法静态成员 | [ReadBool(IReadWriteDevice, Int32, String)](6c90e993-92e2-4b41-a94f-9233ae54991d.htm) | 从PLC的指定地址读取原始的位数据信息，地址示例：MX100, MX10A  Read the original bool data information from the designated address of the PLC. Examples of addresses: MX100, MX10A |
| 公共方法静态成员 | [ReadBool(IReadWriteDeviceStation, String, UInt16)](84c5849a-9b04-82f1-c992-c67d92c7218a.htm) | 从PLC的指定地址读取原始的位数据信息，地址示例：MB100.0, MW100.0  Read the original bool data information from the designated address of the PLC. Examples of addresses: MB100.0, MW100.0 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[LSCnetHelper 类](c874de34-e25a-6c9c-8d0e-ddffe0651200.htm)

[HslCommunication.Profinet.LSIS.Helper 命名空间](0200f8be-eba3-dea9-4eb6-b51d8c230c6b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadBool 方法 (IReadWriteDevice, Int32, String)

[原文連結](http://api.hslcommunication.cn/html/6c90e993-92e2-4b41-a94f-9233ae54991d.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.LSIS.Helper](../html/0200f8be-eba3-dea9-4eb6-b51d8c230c6b.htm "HslCommunication.Profinet.LSIS.Helper")

[LSCnetHelper 类](../html/c874de34-e25a-6c9c-8d0e-ddffe0651200.htm "LSCnetHelper 类")

[LSCnetHelper 方法](../html/d8e60529-bf76-1a12-c052-e3bfeaa4ec33.htm "LSCnetHelper 方法")

[ReadBool 方法](../html/3359728b-1b82-547f-13ea-05d524a7bacc.htm "ReadBool 方法 ")

[ReadBool 方法 (IReadWriteDevice, Int32, String)](../html/6c90e993-92e2-4b41-a94f-9233ae54991d.htm "ReadBool 方法 (IReadWriteDevice, Int32, String)")

[ReadBool 方法 (IReadWriteDeviceStation, String, UInt16)](../html/84c5849a-9b04-82f1-c992-c67d92c7218a.htm "ReadBool 方法 (IReadWriteDeviceStation, String, UInt16)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| LSCnetHelperReadBool 方法 (IReadWriteDevice, Int32, String) |

从PLC的指定地址读取原始的位数据信息，地址示例：MX100, MX10A  
Read the original bool data information from the designated address of the PLC.
Examples of addresses: MX100, MX10A

**命名空间：**
 [HslCommunication.Profinet.LSIS.Helper](0200f8be-eba3-dea9-4eb6-b51d8c230c6b.htm)  
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
	int station,
	string address
)
```

```
Public Shared Function ReadBool ( 
	plc As IReadWriteDevice,
	station As Integer,
	address As String
) As OperateResult(Of Boolean)
```

```
public:
static OperateResult<bool>^ ReadBool(
	IReadWriteDevice^ plc, 
	int station, 
	String^ address
)
```

```
static member ReadBool : 
        plc : IReadWriteDevice * 
        station : int * 
        address : string -> OperateResult<bool> 
```

#### 参数

plc
:   类型：[HslCommunication.CoreIReadWriteDevice](d987afcc-f99c-6aed-3847-5b1890504635.htm)  
    PLC通信对象

station
:   类型：SystemInt32  
    站号信息

address
:   类型：SystemString  
    PLC的地址信息，例如 MX100, MX10A

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Boolean  
返回是否读取成功的结果对象

![](../icons/SectionExpanded.png)备注

地址类型支持 P,M,L,K,F,T,C,D,R,I,Q,W, 支持携带站号的形式，例如 s=2;MX100

![](../icons/SectionExpanded.png)参见

#### 引用

[LSCnetHelper 类](c874de34-e25a-6c9c-8d0e-ddffe0651200.htm)

[ReadBool 重载](3359728b-1b82-547f-13ea-05d524a7bacc.htm)

[HslCommunication.Profinet.LSIS.Helper 命名空间](0200f8be-eba3-dea9-4eb6-b51d8c230c6b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadBool 方法 (IReadWriteDeviceStation, String, UInt16)

[原文連結](http://api.hslcommunication.cn/html/84c5849a-9b04-82f1-c992-c67d92c7218a.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.LSIS.Helper](../html/0200f8be-eba3-dea9-4eb6-b51d8c230c6b.htm "HslCommunication.Profinet.LSIS.Helper")

[LSCnetHelper 类](../html/c874de34-e25a-6c9c-8d0e-ddffe0651200.htm "LSCnetHelper 类")

[LSCnetHelper 方法](../html/d8e60529-bf76-1a12-c052-e3bfeaa4ec33.htm "LSCnetHelper 方法")

[ReadBool 方法](../html/3359728b-1b82-547f-13ea-05d524a7bacc.htm "ReadBool 方法 ")

[ReadBool 方法 (IReadWriteDevice, Int32, String)](../html/6c90e993-92e2-4b41-a94f-9233ae54991d.htm "ReadBool 方法 (IReadWriteDevice, Int32, String)")

[ReadBool 方法 (IReadWriteDeviceStation, String, UInt16)](../html/84c5849a-9b04-82f1-c992-c67d92c7218a.htm "ReadBool 方法 (IReadWriteDeviceStation, String, UInt16)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| LSCnetHelperReadBool 方法 (IReadWriteDeviceStation, String, UInt16) |

从PLC的指定地址读取原始的位数据信息，地址示例：MB100.0, MW100.0  
Read the original bool data information from the designated address of the PLC.
Examples of addresses: MB100.0, MW100.0

**命名空间：**
 [HslCommunication.Profinet.LSIS.Helper](0200f8be-eba3-dea9-4eb6-b51d8c230c6b.htm)  
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
	IReadWriteDeviceStation plc,
	string address,
	ushort length
)
```

```
Public Shared Function ReadBool ( 
	plc As IReadWriteDeviceStation,
	address As String,
	length As UShort
) As OperateResult(Of Boolean())
```

```
public:
static OperateResult<array<bool>^>^ ReadBool(
	IReadWriteDeviceStation^ plc, 
	String^ address, 
	unsigned short length
)
```

```
static member ReadBool : 
        plc : IReadWriteDeviceStation * 
        address : string * 
        length : uint16 -> OperateResult<bool[]> 
```

#### 参数

plc
:   类型：[HslCommunication.Core.NetIReadWriteDeviceStation](0a167412-aa6f-f505-8c8e-75fa32832055.htm)  
    PLC通信对象

address
:   类型：SystemString  
    PLC的地址信息，例如 MB100.0, MW100.0

length
:   类型：SystemUInt16  
    读取的长度信息

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Boolean  
返回是否读取成功的结果对象

![](../icons/SectionExpanded.png)备注

地址类型支持 P,M,L,K,F,T,C,D,R,I,Q,W, 支持携带站号的形式，例如 s=2;MB100.0

![](../icons/SectionExpanded.png)参见

#### 引用

[LSCnetHelper 类](c874de34-e25a-6c9c-8d0e-ddffe0651200.htm)

[ReadBool 重载](3359728b-1b82-547f-13ea-05d524a7bacc.htm)

[HslCommunication.Profinet.LSIS.Helper 命名空间](0200f8be-eba3-dea9-4eb6-b51d8c230c6b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadBoolAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/6b11ae2c-5cb7-a003-2852-48d340e12c62.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.LSIS.Helper](../html/0200f8be-eba3-dea9-4eb6-b51d8c230c6b.htm "HslCommunication.Profinet.LSIS.Helper")

[LSCnetHelper 类](../html/c874de34-e25a-6c9c-8d0e-ddffe0651200.htm "LSCnetHelper 类")

[LSCnetHelper 方法](../html/d8e60529-bf76-1a12-c052-e3bfeaa4ec33.htm "LSCnetHelper 方法")

[ReadBoolAsync 方法](../html/6b11ae2c-5cb7-a003-2852-48d340e12c62.htm "ReadBoolAsync 方法 ")

[ReadBoolAsync 方法 (IReadWriteDevice, Int32, String)](../html/e8d25674-d63b-00d8-dbd9-15736a738497.htm "ReadBoolAsync 方法 (IReadWriteDevice, Int32, String)")

[ReadBoolAsync 方法 (IReadWriteDeviceStation, String, UInt16)](../html/74cb9308-0f0a-23f2-9d31-73fe16baab20.htm "ReadBoolAsync 方法 (IReadWriteDeviceStation, String, UInt16)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| LSCnetHelperReadBoolAsync 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法静态成员 | [ReadBoolAsync(IReadWriteDevice, Int32, String)](e8d25674-d63b-00d8-dbd9-15736a738497.htm) | 从PLC的指定地址读取原始的位数据信息，地址示例：MX100, MX10A  Read the original bool data information from the designated address of the PLC. Examples of addresses: MX100, MX10A |
| 公共方法静态成员 | [ReadBoolAsync(IReadWriteDeviceStation, String, UInt16)](74cb9308-0f0a-23f2-9d31-73fe16baab20.htm) | 从PLC的指定地址读取原始的位数据信息，地址示例：MB100.0, MW100.0  Read the original bool data information from the designated address of the PLC. Examples of addresses: MB100.0, MW100.0 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[LSCnetHelper 类](c874de34-e25a-6c9c-8d0e-ddffe0651200.htm)

[HslCommunication.Profinet.LSIS.Helper 命名空间](0200f8be-eba3-dea9-4eb6-b51d8c230c6b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadBoolAsync 方法 (IReadWriteDevice, Int32, String)

[原文連結](http://api.hslcommunication.cn/html/e8d25674-d63b-00d8-dbd9-15736a738497.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.LSIS.Helper](../html/0200f8be-eba3-dea9-4eb6-b51d8c230c6b.htm "HslCommunication.Profinet.LSIS.Helper")

[LSCnetHelper 类](../html/c874de34-e25a-6c9c-8d0e-ddffe0651200.htm "LSCnetHelper 类")

[LSCnetHelper 方法](../html/d8e60529-bf76-1a12-c052-e3bfeaa4ec33.htm "LSCnetHelper 方法")

[ReadBoolAsync 方法](../html/6b11ae2c-5cb7-a003-2852-48d340e12c62.htm "ReadBoolAsync 方法 ")

[ReadBoolAsync 方法 (IReadWriteDevice, Int32, String)](../html/e8d25674-d63b-00d8-dbd9-15736a738497.htm "ReadBoolAsync 方法 (IReadWriteDevice, Int32, String)")

[ReadBoolAsync 方法 (IReadWriteDeviceStation, String, UInt16)](../html/74cb9308-0f0a-23f2-9d31-73fe16baab20.htm "ReadBoolAsync 方法 (IReadWriteDeviceStation, String, UInt16)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| LSCnetHelperReadBoolAsync 方法 (IReadWriteDevice, Int32, String) |

从PLC的指定地址读取原始的位数据信息，地址示例：MX100, MX10A  
Read the original bool data information from the designated address of the PLC.
Examples of addresses: MX100, MX10A

**命名空间：**
 [HslCommunication.Profinet.LSIS.Helper](0200f8be-eba3-dea9-4eb6-b51d8c230c6b.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static Task<OperateResult<bool>> ReadBoolAsync(
	IReadWriteDevice plc,
	int station,
	string address
)
```

```
Public Shared Function ReadBoolAsync ( 
	plc As IReadWriteDevice,
	station As Integer,
	address As String
) As Task(Of OperateResult(Of Boolean))
```

```
public:
static Task<OperateResult<bool>^>^ ReadBoolAsync(
	IReadWriteDevice^ plc, 
	int station, 
	String^ address
)
```

```
static member ReadBoolAsync : 
        plc : IReadWriteDevice * 
        station : int * 
        address : string -> Task<OperateResult<bool>> 
```

#### 参数

plc
:   类型：[HslCommunication.CoreIReadWriteDevice](d987afcc-f99c-6aed-3847-5b1890504635.htm)  
    PLC通信对象

station
:   类型：SystemInt32  
    站号信息

address
:   类型：SystemString  
    PLC的地址信息，例如 MX100, MX10A

#### 返回值

类型：Task[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Boolean  
返回是否读取成功的结果对象

![](../icons/SectionExpanded.png)备注

地址类型支持 P,M,L,K,F,T,C,D,R,I,Q,W, 支持携带站号的形式，例如 s=2;MX100

![](../icons/SectionExpanded.png)参见

#### 引用

[LSCnetHelper 类](c874de34-e25a-6c9c-8d0e-ddffe0651200.htm)

[ReadBoolAsync 重载](6b11ae2c-5cb7-a003-2852-48d340e12c62.htm)

[HslCommunication.Profinet.LSIS.Helper 命名空间](0200f8be-eba3-dea9-4eb6-b51d8c230c6b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadBoolAsync 方法 (IReadWriteDeviceStation, String, UInt16)

[原文連結](http://api.hslcommunication.cn/html/74cb9308-0f0a-23f2-9d31-73fe16baab20.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.LSIS.Helper](../html/0200f8be-eba3-dea9-4eb6-b51d8c230c6b.htm "HslCommunication.Profinet.LSIS.Helper")

[LSCnetHelper 类](../html/c874de34-e25a-6c9c-8d0e-ddffe0651200.htm "LSCnetHelper 类")

[LSCnetHelper 方法](../html/d8e60529-bf76-1a12-c052-e3bfeaa4ec33.htm "LSCnetHelper 方法")

[ReadBoolAsync 方法](../html/6b11ae2c-5cb7-a003-2852-48d340e12c62.htm "ReadBoolAsync 方法 ")

[ReadBoolAsync 方法 (IReadWriteDevice, Int32, String)](../html/e8d25674-d63b-00d8-dbd9-15736a738497.htm "ReadBoolAsync 方法 (IReadWriteDevice, Int32, String)")

[ReadBoolAsync 方法 (IReadWriteDeviceStation, String, UInt16)](../html/74cb9308-0f0a-23f2-9d31-73fe16baab20.htm "ReadBoolAsync 方法 (IReadWriteDeviceStation, String, UInt16)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| LSCnetHelperReadBoolAsync 方法 (IReadWriteDeviceStation, String, UInt16) |

从PLC的指定地址读取原始的位数据信息，地址示例：MB100.0, MW100.0  
Read the original bool data information from the designated address of the PLC.
Examples of addresses: MB100.0, MW100.0

**命名空间：**
 [HslCommunication.Profinet.LSIS.Helper](0200f8be-eba3-dea9-4eb6-b51d8c230c6b.htm)  
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
	IReadWriteDeviceStation plc,
	string address,
	ushort length
)
```

```
Public Shared Function ReadBoolAsync ( 
	plc As IReadWriteDeviceStation,
	address As String,
	length As UShort
) As Task(Of OperateResult(Of Boolean()))
```

```
public:
static Task<OperateResult<array<bool>^>^>^ ReadBoolAsync(
	IReadWriteDeviceStation^ plc, 
	String^ address, 
	unsigned short length
)
```

```
static member ReadBoolAsync : 
        plc : IReadWriteDeviceStation * 
        address : string * 
        length : uint16 -> Task<OperateResult<bool[]>> 
```

#### 参数

plc
:   类型：[HslCommunication.Core.NetIReadWriteDeviceStation](0a167412-aa6f-f505-8c8e-75fa32832055.htm)  
    PLC通信对象

address
:   类型：SystemString  
    PLC的地址信息，例如 MB100.0, MW100.0

length
:   类型：SystemUInt16  
    读取的长度信息

#### 返回值

类型：Task[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Boolean  
返回是否读取成功的结果对象

![](../icons/SectionExpanded.png)备注

地址类型支持 P,M,L,K,F,T,C,D,R,I,Q,W, 支持携带站号的形式，例如 s=2;MB100.0

![](../icons/SectionExpanded.png)参见

#### 引用

[LSCnetHelper 类](c874de34-e25a-6c9c-8d0e-ddffe0651200.htm)

[ReadBoolAsync 重载](6b11ae2c-5cb7-a003-2852-48d340e12c62.htm)

[HslCommunication.Profinet.LSIS.Helper 命名空间](0200f8be-eba3-dea9-4eb6-b51d8c230c6b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## UnpackResponseContent 方法 

[原文連結](http://api.hslcommunication.cn/html/ef40b9f3-0fc8-b1dd-945b-8d2011f1b6a9.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.LSIS.Helper](../html/0200f8be-eba3-dea9-4eb6-b51d8c230c6b.htm "HslCommunication.Profinet.LSIS.Helper")

[LSCnetHelper 类](../html/c874de34-e25a-6c9c-8d0e-ddffe0651200.htm "LSCnetHelper 类")

[LSCnetHelper 方法](../html/d8e60529-bf76-1a12-c052-e3bfeaa4ec33.htm "LSCnetHelper 方法")

[AnalysisAddress 方法](../html/15ee9fd3-6c4b-6775-db86-370121c062ff.htm "AnalysisAddress 方法 ")

[BuildReadByteCommand 方法](../html/187eecaa-60f5-ed67-f0c6-ce15a59e1e52.htm "BuildReadByteCommand 方法 ")

[BuildReadIndividualCommand 方法](../html/883ada55-fd41-162f-09a9-2d34f92ac605.htm "BuildReadIndividualCommand 方法 ")

[BuildWriteByteCommand 方法](../html/62a3477c-e064-8c71-9687-4f1cc45ebc16.htm "BuildWriteByteCommand 方法 ")

[BuildWriteCommand 方法](../html/2873d537-6750-08f6-ac89-0bf2ea205de3.htm "BuildWriteCommand 方法 ")

[BuildWriteOneCommand 方法](../html/c59dc3f6-675c-8f9d-9dce-3c0857887b3b.htm "BuildWriteOneCommand 方法 ")

[GetAddressOfU\_Q\_I 方法](../html/c3212593-efad-0a6e-22a6-11be146a83d9.htm "GetAddressOfU_Q_I 方法 ")

[GetErrorText 方法](../html/d8d3953e-a826-5558-6773-8809cf6d9dbe.htm "GetErrorText 方法 ")

[Read 方法](../html/7c298c2b-c4d6-de5b-8b67-1ddb7fb1f180.htm "Read 方法 ")

[ReadAsync 方法](../html/f3d85a1f-1978-d789-f6a0-f1dd0fb17e8c.htm "ReadAsync 方法 ")

[ReadBool 方法](../html/3359728b-1b82-547f-13ea-05d524a7bacc.htm "ReadBool 方法 ")

[ReadBoolAsync 方法](../html/6b11ae2c-5cb7-a003-2852-48d340e12c62.htm "ReadBoolAsync 方法 ")

[UnpackResponseContent 方法](../html/ef40b9f3-0fc8-b1dd-945b-8d2011f1b6a9.htm "UnpackResponseContent 方法 ")

[Write 方法](../html/2792b231-69a8-6fae-dc81-e778e73172e1.htm "Write 方法 ")

[WriteAsync 方法](../html/a8b58372-0a53-0c10-6340-30584a178a1b.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| LSCnetHelperUnpackResponseContent 方法 |

根据对方返回的报文命令，对命令进行基本的拆包，例如各种Modbus协议拆包为统一的核心报文，还支持对报文的验证  
According to the message command returned by the other party, the command is basically unpacked, for example,
various Modbus protocols are unpacked into a unified core message, and the verification of the message is also supported

**命名空间：**
 [HslCommunication.Profinet.LSIS.Helper](0200f8be-eba3-dea9-4eb6-b51d8c230c6b.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult<byte[]> UnpackResponseContent(
	byte[] send,
	byte[] response
)
```

```
Public Shared Function UnpackResponseContent ( 
	send As Byte(),
	response As Byte()
) As OperateResult(Of Byte())
```

```
public:
static OperateResult<array<unsigned char>^>^ UnpackResponseContent(
	array<unsigned char>^ send, 
	array<unsigned char>^ response
)
```

```
static member UnpackResponseContent : 
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

[LSCnetHelper 类](c874de34-e25a-6c9c-8d0e-ddffe0651200.htm)

[HslCommunication.Profinet.LSIS.Helper 命名空间](0200f8be-eba3-dea9-4eb6-b51d8c230c6b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Write 方法 

[原文連結](http://api.hslcommunication.cn/html/2792b231-69a8-6fae-dc81-e778e73172e1.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.LSIS.Helper](../html/0200f8be-eba3-dea9-4eb6-b51d8c230c6b.htm "HslCommunication.Profinet.LSIS.Helper")

[LSCnetHelper 类](../html/c874de34-e25a-6c9c-8d0e-ddffe0651200.htm "LSCnetHelper 类")

[LSCnetHelper 方法](../html/d8e60529-bf76-1a12-c052-e3bfeaa4ec33.htm "LSCnetHelper 方法")

[Write 方法](../html/2792b231-69a8-6fae-dc81-e778e73172e1.htm "Write 方法 ")

[Write 方法 (IReadWriteDevice, Int32, String, Boolean)](../html/e253354d-f080-9498-b34e-710eb2ff95c3.htm "Write 方法 (IReadWriteDevice, Int32, String, Boolean)")

[Write 方法 (IReadWriteDevice, Int32, String, Byte[])](../html/9fb7ffe6-6132-bbd7-3287-b65dc8956570.htm "Write 方法 (IReadWriteDevice, Int32, String, Byte[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| LSCnetHelperWrite 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法静态成员 | [Write(IReadWriteDevice, Int32, String, Boolean)](e253354d-f080-9498-b34e-710eb2ff95c3.htm) | 将bool数据写入到PLC的指定的地址里，地址示例：MX100, MX10A  Write the bool data to the designated address of the PLC. Examples of addresses: MX100, MX10A |
| 公共方法静态成员 | [Write(IReadWriteDevice, Int32, String, Byte)](9fb7ffe6-6132-bbd7-3287-b65dc8956570.htm) | 将原始数据写入到PLC的指定的地址里，地址示例：MB100, MW100, MD100, 如果输入了M100等同于MB100  Write the original data to the designated address of the PLC. Examples of addresses: MB100, MW100, MD100, if input M100 is equivalent to MB100 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[LSCnetHelper 类](c874de34-e25a-6c9c-8d0e-ddffe0651200.htm)

[HslCommunication.Profinet.LSIS.Helper 命名空间](0200f8be-eba3-dea9-4eb6-b51d8c230c6b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Write 方法 (IReadWriteDevice, Int32, String, Boolean)

[原文連結](http://api.hslcommunication.cn/html/e253354d-f080-9498-b34e-710eb2ff95c3.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.LSIS.Helper](../html/0200f8be-eba3-dea9-4eb6-b51d8c230c6b.htm "HslCommunication.Profinet.LSIS.Helper")

[LSCnetHelper 类](../html/c874de34-e25a-6c9c-8d0e-ddffe0651200.htm "LSCnetHelper 类")

[LSCnetHelper 方法](../html/d8e60529-bf76-1a12-c052-e3bfeaa4ec33.htm "LSCnetHelper 方法")

[Write 方法](../html/2792b231-69a8-6fae-dc81-e778e73172e1.htm "Write 方法 ")

[Write 方法 (IReadWriteDevice, Int32, String, Boolean)](../html/e253354d-f080-9498-b34e-710eb2ff95c3.htm "Write 方法 (IReadWriteDevice, Int32, String, Boolean)")

[Write 方法 (IReadWriteDevice, Int32, String, Byte[])](../html/9fb7ffe6-6132-bbd7-3287-b65dc8956570.htm "Write 方法 (IReadWriteDevice, Int32, String, Byte[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| LSCnetHelperWrite 方法 (IReadWriteDevice, Int32, String, Boolean) |

将bool数据写入到PLC的指定的地址里，地址示例：MX100, MX10A  
Write the bool data to the designated address of the PLC. Examples of addresses: MX100, MX10A

**命名空间：**
 [HslCommunication.Profinet.LSIS.Helper](0200f8be-eba3-dea9-4eb6-b51d8c230c6b.htm)  
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
	IReadWriteDevice plc,
	int station,
	string address,
	bool value
)
```

```
Public Shared Function Write ( 
	plc As IReadWriteDevice,
	station As Integer,
	address As String,
	value As Boolean
) As OperateResult
```

```
public:
static OperateResult^ Write(
	IReadWriteDevice^ plc, 
	int station, 
	String^ address, 
	bool value
)
```

```
static member Write : 
        plc : IReadWriteDevice * 
        station : int * 
        address : string * 
        value : bool -> OperateResult 
```

#### 参数

plc
:   类型：[HslCommunication.CoreIReadWriteDevice](d987afcc-f99c-6aed-3847-5b1890504635.htm)  
    PLC通信对象

station
:   类型：SystemInt32  
    站号信息

address
:   类型：SystemString  
    PLC的地址信息，例如 MX100, MX10A

value
:   类型：SystemBoolean  
    bool值信息

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
返回是否读取成功的结果对象

![](../icons/SectionExpanded.png)备注

地址类型支持 P,M,L,K,F,T,C,D,R,I,Q,W, 支持携带站号的形式，例如 s=2;MX100

![](../icons/SectionExpanded.png)参见

#### 引用

[LSCnetHelper 类](c874de34-e25a-6c9c-8d0e-ddffe0651200.htm)

[Write 重载](2792b231-69a8-6fae-dc81-e778e73172e1.htm)

[HslCommunication.Profinet.LSIS.Helper 命名空间](0200f8be-eba3-dea9-4eb6-b51d8c230c6b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Write 方法 (IReadWriteDevice, Int32, String, Byte[])

[原文連結](http://api.hslcommunication.cn/html/9fb7ffe6-6132-bbd7-3287-b65dc8956570.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.LSIS.Helper](../html/0200f8be-eba3-dea9-4eb6-b51d8c230c6b.htm "HslCommunication.Profinet.LSIS.Helper")

[LSCnetHelper 类](../html/c874de34-e25a-6c9c-8d0e-ddffe0651200.htm "LSCnetHelper 类")

[LSCnetHelper 方法](../html/d8e60529-bf76-1a12-c052-e3bfeaa4ec33.htm "LSCnetHelper 方法")

[Write 方法](../html/2792b231-69a8-6fae-dc81-e778e73172e1.htm "Write 方法 ")

[Write 方法 (IReadWriteDevice, Int32, String, Boolean)](../html/e253354d-f080-9498-b34e-710eb2ff95c3.htm "Write 方法 (IReadWriteDevice, Int32, String, Boolean)")

[Write 方法 (IReadWriteDevice, Int32, String, Byte[])](../html/9fb7ffe6-6132-bbd7-3287-b65dc8956570.htm "Write 方法 (IReadWriteDevice, Int32, String, Byte[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| LSCnetHelperWrite 方法 (IReadWriteDevice, Int32, String, Byte) |

将原始数据写入到PLC的指定的地址里，地址示例：MB100, MW100, MD100, 如果输入了M100等同于MB100  
Write the original data to the designated address of the PLC.
Examples of addresses: MB100, MW100, MD100, if input M100 is equivalent to MB100

**命名空间：**
 [HslCommunication.Profinet.LSIS.Helper](0200f8be-eba3-dea9-4eb6-b51d8c230c6b.htm)  
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
	IReadWriteDevice plc,
	int station,
	string address,
	byte[] value
)
```

```
Public Shared Function Write ( 
	plc As IReadWriteDevice,
	station As Integer,
	address As String,
	value As Byte()
) As OperateResult
```

```
public:
static OperateResult^ Write(
	IReadWriteDevice^ plc, 
	int station, 
	String^ address, 
	array<unsigned char>^ value
)
```

```
static member Write : 
        plc : IReadWriteDevice * 
        station : int * 
        address : string * 
        value : byte[] -> OperateResult 
```

#### 参数

plc
:   类型：[HslCommunication.CoreIReadWriteDevice](d987afcc-f99c-6aed-3847-5b1890504635.htm)  
    PLC通信对象

station
:   类型：SystemInt32  
    站号信息

address
:   类型：SystemString  
    PLC的地址信息，例如 M100, MB100, MW100, MD100

value
:   类型：SystemByte  
    等待写入的原始数据内容

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否写入成功

![](../icons/SectionExpanded.png)参见

#### 引用

[LSCnetHelper 类](c874de34-e25a-6c9c-8d0e-ddffe0651200.htm)

[Write 重载](2792b231-69a8-6fae-dc81-e778e73172e1.htm)

[HslCommunication.Profinet.LSIS.Helper 命名空间](0200f8be-eba3-dea9-4eb6-b51d8c230c6b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)