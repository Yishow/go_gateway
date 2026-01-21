# HslCommunication - HslCommunication.Profinet.GE

> 分類頁數: 30



---
## HslCommunication.Profinet.GE

[原文連結](http://api.hslcommunication.cn/html/db594ebf-c4ae-a11c-a82c-72625c3b6ab4.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.GE](../html/db594ebf-c4ae-a11c-a82c-72625c3b6ab4.htm "HslCommunication.Profinet.GE")

[GeHelper 类](../html/3c60bc6c-29f1-ae7f-d703-4bfc7c463e98.htm "GeHelper 类")

[GeSRTPNet 类](../html/52c0166b-3420-af47-8dec-f25e74515116.htm "GeSRTPNet 类")

[GeSRTPServer 类](../html/5d0019d1-c525-e17f-1b80-8bf8607878c2.htm "GeSRTPServer 类")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslCommunication.Profinet.GE 命名空间 |

[缺少 "N:HslCommunication.Profinet.GE" 的 <summary> 文档]

![](../icons/SectionExpanded.png)类

|  | 类 | 说明 |
| --- | --- | --- |
| 公共类 | [GeHelper](3c60bc6c-29f1-ae7f-d703-4bfc7c463e98.htm) | GE plc相关的辅助类对象 |
| 公共类代码示例 | [GeSRTPNet](52c0166b-3420-af47-8dec-f25e74515116.htm) | Ge(通用电气)集团的SRTP协议实现的客户端，支持 I,Q,M,T,SA,SB,SC,S,G 的位和字节读写，支持 AI,AQ,R 的字读写操作，支持读取PLC时间，程序名操作  The client implemented by Ge (General Electric) Group's SRTP protocol, supports I, Q, M, T, SA, SB, SC, S, G bit and byte read and write, supports AI, AQ, R word read and write Operation, support reading PLC time, program name operation |
| 公共类代码示例 | [GeSRTPServer](5d0019d1-c525-e17f-1b80-8bf8607878c2.htm) | Ge的SRTP协议实现的虚拟PLC，支持I,Q,M,T,SA,SB,SC,S,G的位和字节读写，支持AI,AQ,R的字读写操作，支持读取当前时间及程序名称。  Virtual PLC implemented by Ge's SRTP protocol, supports bit and byte read and write of I, Q, M, T, SA, SB, SC, S, G, supports word read and write operations of AI, AQ, R, and supports reading Current time and program name. |

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GeHelper 类

[原文連結](http://api.hslcommunication.cn/html/3c60bc6c-29f1-ae7f-d703-4bfc7c463e98.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.GE](../html/db594ebf-c4ae-a11c-a82c-72625c3b6ab4.htm "HslCommunication.Profinet.GE")

[GeHelper 类](../html/3c60bc6c-29f1-ae7f-d703-4bfc7c463e98.htm "GeHelper 类")

[GeHelper 构造函数](../html/b5f05274-a0e3-4885-2077-725dc4bc47dd.htm "GeHelper 构造函数 ")

[GeHelper 方法](../html/59587f2f-8da2-a1fb-dc9d-5f031122efbe.htm "GeHelper 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| GeHelper 类 |

GE plc相关的辅助类对象

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  HslCommunication.Profinet.GEGeHelper

**命名空间：**
 [HslCommunication.Profinet.GE](db594ebf-c4ae-a11c-a82c-72625c3b6ab4.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class GeHelper
```

```
Public Class GeHelper
```

```
public ref class GeHelper
```

```
type GeHelper =  class end
```

GeHelper 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [GeHelper](b5f05274-a0e3-4885-2077-725dc4bc47dd.htm) | 初始化 GeHelper 类的一个新实例 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法静态成员 | [BuildReadCommand(Int64, GeSRTPAddress)](da60a578-4c2f-3c9e-8e0d-7f5a0700e498.htm) | 构建一个读取数据的报文命令，需要指定消息号，读取的 GE 地址信息  To construct a message command to read data, you need to specify the message number and read GE address information |
| 公共方法静态成员 | [BuildReadCommand(Int64, String, UInt16, Boolean)](d0787b47-fcba-7a3a-632f-1fe47bdc2c9e.htm) | 构建一个读取数据的报文命令，需要指定消息号，地址，长度，是否位读取，返回完整的报文信息。  To construct a message command to read data, you need to specify the message number, address, length, whether to read in bits, and return the complete message information. |
| 公共方法静态成员 | [BuildReadCoreCommand](243b844b-95fe-ce64-3e8d-29f5e9638098.htm) | 构建一个读取数据的报文信息，需要指定操作的数据代码，读取的参数信息  To construct a message information for reading data, you need to specify the data code of the operation and the parameter information to be read |
| 公共方法静态成员 | [BuildWriteCommand(Int64, GeSRTPAddress, Byte)](e4d4aae5-bf38-f48e-0f40-694198d4e944.htm) | 构建一个批量写入 byte 数组变量的报文，需要指定消息号，写入的地址，地址参照 [GeSRTPNet](52c0166b-3420-af47-8dec-f25e74515116.htm) 说明。  To construct a message to be written into byte array variables in batches, you need to specify the message number and write address. For the address, refer to the description of [GeSRTPNet](52c0166b-3420-af47-8dec-f25e74515116.htm). |
| 公共方法静态成员 | [BuildWriteCommand(Int64, String, Boolean)](969ce7ae-8486-af96-782e-31df7ef63f54.htm) | 构建一个批量写入 bool 数组变量的报文，需要指定消息号，写入的地址，地址参照 [GeSRTPNet](52c0166b-3420-af47-8dec-f25e74515116.htm) 说明。  To construct a message to be written into bool array variables in batches, you need to specify the message number and write address. For the address, refer to the description of [GeSRTPNet](52c0166b-3420-af47-8dec-f25e74515116.htm). |
| 公共方法静态成员 | [BuildWriteCommand(Int64, String, Byte)](bee10f7e-5595-5100-4d76-3501ae2139dc.htm) | 构建一个批量写入 byte 数组变量的报文，需要指定消息号，写入的地址，地址参照 [GeSRTPNet](52c0166b-3420-af47-8dec-f25e74515116.htm) 说明。  To construct a message to be written into byte array variables in batches, you need to specify the message number and write address. For the address, refer to the description of [GeSRTPNet](52c0166b-3420-af47-8dec-f25e74515116.htm). |
| 公共方法 | Equals | (继承自 Object。) |
| 公共方法静态成员 | [ExtraDateTime](25f30152-fad8-5f60-5176-7dca1edfccea.htm) | 从实际的时间的字节数组里解析出C#格式的时间对象，这个时间可能是时区0的时间，需要自行转化本地时间。  Analyze the time object in C# format from the actual time byte array. This time may be the time in time zone 0, and you need to convert the local time yourself. |
| 公共方法静态成员 | [ExtraProgramName](bddf3623-efe4-07e3-130a-54c019072cc7.htm) | 从实际的时间的字节数组里解析出PLC的程序的名称。  Parse the name of the PLC program from the actual time byte array |
| 公共方法静态成员 | [ExtraResponseContent](62346044-4399-25c6-f5e1-62fb52f02957.htm) | 从PLC返回的数据中，提取出实际的数据内容，最少6个字节的数据。超出实际的数据长度的部分没有任何意义。  From the data returned by the PLC, extract the actual data content, at least 6 bytes of data. The part beyond the actual data length has no meaning. |
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

[HslCommunication.Profinet.GE 命名空间](db594ebf-c4ae-a11c-a82c-72625c3b6ab4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GeHelper 构造函数 

[原文連結](http://api.hslcommunication.cn/html/b5f05274-a0e3-4885-2077-725dc4bc47dd.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.GE](../html/db594ebf-c4ae-a11c-a82c-72625c3b6ab4.htm "HslCommunication.Profinet.GE")

[GeHelper 类](../html/3c60bc6c-29f1-ae7f-d703-4bfc7c463e98.htm "GeHelper 类")

[GeHelper 构造函数](../html/b5f05274-a0e3-4885-2077-725dc4bc47dd.htm "GeHelper 构造函数 ")

[GeHelper 方法](../html/59587f2f-8da2-a1fb-dc9d-5f031122efbe.htm "GeHelper 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| GeHelper 构造函数 |

初始化 [GeHelper](3c60bc6c-29f1-ae7f-d703-4bfc7c463e98.htm) 类的一个新实例

**命名空间：**
 [HslCommunication.Profinet.GE](db594ebf-c4ae-a11c-a82c-72625c3b6ab4.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public GeHelper()
```

```
Public Sub New
```

```
public:
GeHelper()
```

```
new : unit -> GeHelper
```

![](../icons/SectionExpanded.png)参见

#### 引用

[GeHelper 类](3c60bc6c-29f1-ae7f-d703-4bfc7c463e98.htm)

[HslCommunication.Profinet.GE 命名空间](db594ebf-c4ae-a11c-a82c-72625c3b6ab4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GeHelper 方法

[原文連結](http://api.hslcommunication.cn/html/59587f2f-8da2-a1fb-dc9d-5f031122efbe.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.GE](../html/db594ebf-c4ae-a11c-a82c-72625c3b6ab4.htm "HslCommunication.Profinet.GE")

[GeHelper 类](../html/3c60bc6c-29f1-ae7f-d703-4bfc7c463e98.htm "GeHelper 类")

[GeHelper 方法](../html/59587f2f-8da2-a1fb-dc9d-5f031122efbe.htm "GeHelper 方法")

[BuildReadCommand 方法](../html/9ff3c313-5f82-78f1-ee0c-ebd401fc06eb.htm "BuildReadCommand 方法 ")

[BuildReadCoreCommand 方法](../html/243b844b-95fe-ce64-3e8d-29f5e9638098.htm "BuildReadCoreCommand 方法 ")

[BuildWriteCommand 方法](../html/f1b77cc0-230f-bef6-7f79-e621a76b5d1d.htm "BuildWriteCommand 方法 ")

[ExtraDateTime 方法](../html/25f30152-fad8-5f60-5176-7dca1edfccea.htm "ExtraDateTime 方法 ")

[ExtraProgramName 方法](../html/bddf3623-efe4-07e3-130a-54c019072cc7.htm "ExtraProgramName 方法 ")

[ExtraResponseContent 方法](../html/62346044-4399-25c6-f5e1-62fb52f02957.htm "ExtraResponseContent 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| GeHelper 方法 |

[GeHelper](3c60bc6c-29f1-ae7f-d703-4bfc7c463e98.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法静态成员 | [BuildReadCommand(Int64, GeSRTPAddress)](da60a578-4c2f-3c9e-8e0d-7f5a0700e498.htm) | 构建一个读取数据的报文命令，需要指定消息号，读取的 GE 地址信息  To construct a message command to read data, you need to specify the message number and read GE address information |
| 公共方法静态成员 | [BuildReadCommand(Int64, String, UInt16, Boolean)](d0787b47-fcba-7a3a-632f-1fe47bdc2c9e.htm) | 构建一个读取数据的报文命令，需要指定消息号，地址，长度，是否位读取，返回完整的报文信息。  To construct a message command to read data, you need to specify the message number, address, length, whether to read in bits, and return the complete message information. |
| 公共方法静态成员 | [BuildReadCoreCommand](243b844b-95fe-ce64-3e8d-29f5e9638098.htm) | 构建一个读取数据的报文信息，需要指定操作的数据代码，读取的参数信息  To construct a message information for reading data, you need to specify the data code of the operation and the parameter information to be read |
| 公共方法静态成员 | [BuildWriteCommand(Int64, GeSRTPAddress, Byte)](e4d4aae5-bf38-f48e-0f40-694198d4e944.htm) | 构建一个批量写入 byte 数组变量的报文，需要指定消息号，写入的地址，地址参照 [GeSRTPNet](52c0166b-3420-af47-8dec-f25e74515116.htm) 说明。  To construct a message to be written into byte array variables in batches, you need to specify the message number and write address. For the address, refer to the description of [GeSRTPNet](52c0166b-3420-af47-8dec-f25e74515116.htm). |
| 公共方法静态成员 | [BuildWriteCommand(Int64, String, Boolean)](969ce7ae-8486-af96-782e-31df7ef63f54.htm) | 构建一个批量写入 bool 数组变量的报文，需要指定消息号，写入的地址，地址参照 [GeSRTPNet](52c0166b-3420-af47-8dec-f25e74515116.htm) 说明。  To construct a message to be written into bool array variables in batches, you need to specify the message number and write address. For the address, refer to the description of [GeSRTPNet](52c0166b-3420-af47-8dec-f25e74515116.htm). |
| 公共方法静态成员 | [BuildWriteCommand(Int64, String, Byte)](bee10f7e-5595-5100-4d76-3501ae2139dc.htm) | 构建一个批量写入 byte 数组变量的报文，需要指定消息号，写入的地址，地址参照 [GeSRTPNet](52c0166b-3420-af47-8dec-f25e74515116.htm) 说明。  To construct a message to be written into byte array variables in batches, you need to specify the message number and write address. For the address, refer to the description of [GeSRTPNet](52c0166b-3420-af47-8dec-f25e74515116.htm). |
| 公共方法 | Equals | (继承自 Object。) |
| 公共方法静态成员 | [ExtraDateTime](25f30152-fad8-5f60-5176-7dca1edfccea.htm) | 从实际的时间的字节数组里解析出C#格式的时间对象，这个时间可能是时区0的时间，需要自行转化本地时间。  Analyze the time object in C# format from the actual time byte array. This time may be the time in time zone 0, and you need to convert the local time yourself. |
| 公共方法静态成员 | [ExtraProgramName](bddf3623-efe4-07e3-130a-54c019072cc7.htm) | 从实际的时间的字节数组里解析出PLC的程序的名称。  Parse the name of the PLC program from the actual time byte array |
| 公共方法静态成员 | [ExtraResponseContent](62346044-4399-25c6-f5e1-62fb52f02957.htm) | 从PLC返回的数据中，提取出实际的数据内容，最少6个字节的数据。超出实际的数据长度的部分没有任何意义。  From the data returned by the PLC, extract the actual data content, at least 6 bytes of data. The part beyond the actual data length has no meaning. |
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

[GeHelper 类](3c60bc6c-29f1-ae7f-d703-4bfc7c463e98.htm)

[HslCommunication.Profinet.GE 命名空间](db594ebf-c4ae-a11c-a82c-72625c3b6ab4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## BuildReadCommand 方法 

[原文連結](http://api.hslcommunication.cn/html/9ff3c313-5f82-78f1-ee0c-ebd401fc06eb.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.GE](../html/db594ebf-c4ae-a11c-a82c-72625c3b6ab4.htm "HslCommunication.Profinet.GE")

[GeHelper 类](../html/3c60bc6c-29f1-ae7f-d703-4bfc7c463e98.htm "GeHelper 类")

[GeHelper 方法](../html/59587f2f-8da2-a1fb-dc9d-5f031122efbe.htm "GeHelper 方法")

[BuildReadCommand 方法](../html/9ff3c313-5f82-78f1-ee0c-ebd401fc06eb.htm "BuildReadCommand 方法 ")

[BuildReadCommand 方法 (Int64, GeSRTPAddress)](../html/da60a578-4c2f-3c9e-8e0d-7f5a0700e498.htm "BuildReadCommand 方法 (Int64, GeSRTPAddress)")

[BuildReadCommand 方法 (Int64, String, UInt16, Boolean)](../html/d0787b47-fcba-7a3a-632f-1fe47bdc2c9e.htm "BuildReadCommand 方法 (Int64, String, UInt16, Boolean)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| GeHelperBuildReadCommand 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法静态成员 | [BuildReadCommand(Int64, GeSRTPAddress)](da60a578-4c2f-3c9e-8e0d-7f5a0700e498.htm) | 构建一个读取数据的报文命令，需要指定消息号，读取的 GE 地址信息  To construct a message command to read data, you need to specify the message number and read GE address information |
| 公共方法静态成员 | [BuildReadCommand(Int64, String, UInt16, Boolean)](d0787b47-fcba-7a3a-632f-1fe47bdc2c9e.htm) | 构建一个读取数据的报文命令，需要指定消息号，地址，长度，是否位读取，返回完整的报文信息。  To construct a message command to read data, you need to specify the message number, address, length, whether to read in bits, and return the complete message information. |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[GeHelper 类](3c60bc6c-29f1-ae7f-d703-4bfc7c463e98.htm)

[HslCommunication.Profinet.GE 命名空间](db594ebf-c4ae-a11c-a82c-72625c3b6ab4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## BuildReadCommand 方法 (Int64, GeSRTPAddress)

[原文連結](http://api.hslcommunication.cn/html/da60a578-4c2f-3c9e-8e0d-7f5a0700e498.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.GE](../html/db594ebf-c4ae-a11c-a82c-72625c3b6ab4.htm "HslCommunication.Profinet.GE")

[GeHelper 类](../html/3c60bc6c-29f1-ae7f-d703-4bfc7c463e98.htm "GeHelper 类")

[GeHelper 方法](../html/59587f2f-8da2-a1fb-dc9d-5f031122efbe.htm "GeHelper 方法")

[BuildReadCommand 方法](../html/9ff3c313-5f82-78f1-ee0c-ebd401fc06eb.htm "BuildReadCommand 方法 ")

[BuildReadCommand 方法 (Int64, GeSRTPAddress)](../html/da60a578-4c2f-3c9e-8e0d-7f5a0700e498.htm "BuildReadCommand 方法 (Int64, GeSRTPAddress)")

[BuildReadCommand 方法 (Int64, String, UInt16, Boolean)](../html/d0787b47-fcba-7a3a-632f-1fe47bdc2c9e.htm "BuildReadCommand 方法 (Int64, String, UInt16, Boolean)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| GeHelperBuildReadCommand 方法 (Int64, GeSRTPAddress) |

构建一个读取数据的报文命令，需要指定消息号，读取的 GE 地址信息  
To construct a message command to read data, you need to specify the message number and read GE address information

**命名空间：**
 [HslCommunication.Profinet.GE](db594ebf-c4ae-a11c-a82c-72625c3b6ab4.htm)  
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
	long id,
	GeSRTPAddress address
)
```

```
Public Shared Function BuildReadCommand ( 
	id As Long,
	address As GeSRTPAddress
) As OperateResult(Of Byte())
```

```
public:
static OperateResult<array<unsigned char>^>^ BuildReadCommand(
	long long id, 
	GeSRTPAddress^ address
)
```

```
static member BuildReadCommand : 
        id : int64 * 
        address : GeSRTPAddress -> OperateResult<byte[]> 
```

#### 参数

id
:   类型：SystemInt64  
    消息号

address
:   类型：[HslCommunication.Core.AddressGeSRTPAddress](7c65a869-5fea-2bf6-526d-6d2013323256.htm)  
    GE 的地址

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
包含是否成功的报文信息

![](../icons/SectionExpanded.png)参见

#### 引用

[GeHelper 类](3c60bc6c-29f1-ae7f-d703-4bfc7c463e98.htm)

[BuildReadCommand 重载](9ff3c313-5f82-78f1-ee0c-ebd401fc06eb.htm)

[HslCommunication.Profinet.GE 命名空间](db594ebf-c4ae-a11c-a82c-72625c3b6ab4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## BuildReadCommand 方法 (Int64, String, UInt16, Boolean)

[原文連結](http://api.hslcommunication.cn/html/d0787b47-fcba-7a3a-632f-1fe47bdc2c9e.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.GE](../html/db594ebf-c4ae-a11c-a82c-72625c3b6ab4.htm "HslCommunication.Profinet.GE")

[GeHelper 类](../html/3c60bc6c-29f1-ae7f-d703-4bfc7c463e98.htm "GeHelper 类")

[GeHelper 方法](../html/59587f2f-8da2-a1fb-dc9d-5f031122efbe.htm "GeHelper 方法")

[BuildReadCommand 方法](../html/9ff3c313-5f82-78f1-ee0c-ebd401fc06eb.htm "BuildReadCommand 方法 ")

[BuildReadCommand 方法 (Int64, GeSRTPAddress)](../html/da60a578-4c2f-3c9e-8e0d-7f5a0700e498.htm "BuildReadCommand 方法 (Int64, GeSRTPAddress)")

[BuildReadCommand 方法 (Int64, String, UInt16, Boolean)](../html/d0787b47-fcba-7a3a-632f-1fe47bdc2c9e.htm "BuildReadCommand 方法 (Int64, String, UInt16, Boolean)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| GeHelperBuildReadCommand 方法 (Int64, String, UInt16, Boolean) |

构建一个读取数据的报文命令，需要指定消息号，地址，长度，是否位读取，返回完整的报文信息。  
To construct a message command to read data, you need to specify the message number,
address, length, whether to read in bits, and return the complete message information.

**命名空间：**
 [HslCommunication.Profinet.GE](db594ebf-c4ae-a11c-a82c-72625c3b6ab4.htm)  
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
	long id,
	string address,
	ushort length,
	bool isBit
)
```

```
Public Shared Function BuildReadCommand ( 
	id As Long,
	address As String,
	length As UShort,
	isBit As Boolean
) As OperateResult(Of Byte())
```

```
public:
static OperateResult<array<unsigned char>^>^ BuildReadCommand(
	long long id, 
	String^ address, 
	unsigned short length, 
	bool isBit
)
```

```
static member BuildReadCommand : 
        id : int64 * 
        address : string * 
        length : uint16 * 
        isBit : bool -> OperateResult<byte[]> 
```

#### 参数

id
:   类型：SystemInt64  
    消息号

address
:   类型：SystemString  
    地址

length
:   类型：SystemUInt16  
    读取的长度

isBit
:   类型：SystemBoolean  

    [缺少 "M:HslCommunication.Profinet.GE.GeHelper.BuildReadCommand(System.Int64,System.String,System.UInt16,System.Boolean)" 的 <param name="isBit"/> 文档]

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
包含是否成功的报文对象

![](../icons/SectionExpanded.png)参见

#### 引用

[GeHelper 类](3c60bc6c-29f1-ae7f-d703-4bfc7c463e98.htm)

[BuildReadCommand 重载](9ff3c313-5f82-78f1-ee0c-ebd401fc06eb.htm)

[HslCommunication.Profinet.GE 命名空间](db594ebf-c4ae-a11c-a82c-72625c3b6ab4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## BuildReadCoreCommand 方法 

[原文連結](http://api.hslcommunication.cn/html/243b844b-95fe-ce64-3e8d-29f5e9638098.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.GE](../html/db594ebf-c4ae-a11c-a82c-72625c3b6ab4.htm "HslCommunication.Profinet.GE")

[GeHelper 类](../html/3c60bc6c-29f1-ae7f-d703-4bfc7c463e98.htm "GeHelper 类")

[GeHelper 方法](../html/59587f2f-8da2-a1fb-dc9d-5f031122efbe.htm "GeHelper 方法")

[BuildReadCommand 方法](../html/9ff3c313-5f82-78f1-ee0c-ebd401fc06eb.htm "BuildReadCommand 方法 ")

[BuildReadCoreCommand 方法](../html/243b844b-95fe-ce64-3e8d-29f5e9638098.htm "BuildReadCoreCommand 方法 ")

[BuildWriteCommand 方法](../html/f1b77cc0-230f-bef6-7f79-e621a76b5d1d.htm "BuildWriteCommand 方法 ")

[ExtraDateTime 方法](../html/25f30152-fad8-5f60-5176-7dca1edfccea.htm "ExtraDateTime 方法 ")

[ExtraProgramName 方法](../html/bddf3623-efe4-07e3-130a-54c019072cc7.htm "ExtraProgramName 方法 ")

[ExtraResponseContent 方法](../html/62346044-4399-25c6-f5e1-62fb52f02957.htm "ExtraResponseContent 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| GeHelperBuildReadCoreCommand 方法 |

构建一个读取数据的报文信息，需要指定操作的数据代码，读取的参数信息  
To construct a message information for reading data, you need to specify the data code of the operation and the parameter information to be read

**命名空间：**
 [HslCommunication.Profinet.GE](db594ebf-c4ae-a11c-a82c-72625c3b6ab4.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult<byte[]> BuildReadCoreCommand(
	long id,
	byte code,
	byte[] data
)
```

```
Public Shared Function BuildReadCoreCommand ( 
	id As Long,
	code As Byte,
	data As Byte()
) As OperateResult(Of Byte())
```

```
public:
static OperateResult<array<unsigned char>^>^ BuildReadCoreCommand(
	long long id, 
	unsigned char code, 
	array<unsigned char>^ data
)
```

```
static member BuildReadCoreCommand : 
        id : int64 * 
        code : byte * 
        data : byte[] -> OperateResult<byte[]> 
```

#### 参数

id
:   类型：SystemInt64  
    消息号

code
:   类型：SystemByte  
    操作代码

data
:   类型：SystemByte  
    数据参数

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
包含是否成功的报文信息

![](../icons/SectionExpanded.png)参见

#### 引用

[GeHelper 类](3c60bc6c-29f1-ae7f-d703-4bfc7c463e98.htm)

[HslCommunication.Profinet.GE 命名空间](db594ebf-c4ae-a11c-a82c-72625c3b6ab4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## BuildWriteCommand 方法 

[原文連結](http://api.hslcommunication.cn/html/f1b77cc0-230f-bef6-7f79-e621a76b5d1d.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.GE](../html/db594ebf-c4ae-a11c-a82c-72625c3b6ab4.htm "HslCommunication.Profinet.GE")

[GeHelper 类](../html/3c60bc6c-29f1-ae7f-d703-4bfc7c463e98.htm "GeHelper 类")

[GeHelper 方法](../html/59587f2f-8da2-a1fb-dc9d-5f031122efbe.htm "GeHelper 方法")

[BuildWriteCommand 方法](../html/f1b77cc0-230f-bef6-7f79-e621a76b5d1d.htm "BuildWriteCommand 方法 ")

[BuildWriteCommand 方法 (Int64, GeSRTPAddress, Byte[])](../html/e4d4aae5-bf38-f48e-0f40-694198d4e944.htm "BuildWriteCommand 方法 (Int64, GeSRTPAddress, Byte[])")

[BuildWriteCommand 方法 (Int64, String, Boolean[])](../html/969ce7ae-8486-af96-782e-31df7ef63f54.htm "BuildWriteCommand 方法 (Int64, String, Boolean[])")

[BuildWriteCommand 方法 (Int64, String, Byte[])](../html/bee10f7e-5595-5100-4d76-3501ae2139dc.htm "BuildWriteCommand 方法 (Int64, String, Byte[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| GeHelperBuildWriteCommand 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法静态成员 | [BuildWriteCommand(Int64, GeSRTPAddress, Byte)](e4d4aae5-bf38-f48e-0f40-694198d4e944.htm) | 构建一个批量写入 byte 数组变量的报文，需要指定消息号，写入的地址，地址参照 [GeSRTPNet](52c0166b-3420-af47-8dec-f25e74515116.htm) 说明。  To construct a message to be written into byte array variables in batches, you need to specify the message number and write address. For the address, refer to the description of [GeSRTPNet](52c0166b-3420-af47-8dec-f25e74515116.htm). |
| 公共方法静态成员 | [BuildWriteCommand(Int64, String, Boolean)](969ce7ae-8486-af96-782e-31df7ef63f54.htm) | 构建一个批量写入 bool 数组变量的报文，需要指定消息号，写入的地址，地址参照 [GeSRTPNet](52c0166b-3420-af47-8dec-f25e74515116.htm) 说明。  To construct a message to be written into bool array variables in batches, you need to specify the message number and write address. For the address, refer to the description of [GeSRTPNet](52c0166b-3420-af47-8dec-f25e74515116.htm). |
| 公共方法静态成员 | [BuildWriteCommand(Int64, String, Byte)](bee10f7e-5595-5100-4d76-3501ae2139dc.htm) | 构建一个批量写入 byte 数组变量的报文，需要指定消息号，写入的地址，地址参照 [GeSRTPNet](52c0166b-3420-af47-8dec-f25e74515116.htm) 说明。  To construct a message to be written into byte array variables in batches, you need to specify the message number and write address. For the address, refer to the description of [GeSRTPNet](52c0166b-3420-af47-8dec-f25e74515116.htm). |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[GeHelper 类](3c60bc6c-29f1-ae7f-d703-4bfc7c463e98.htm)

[HslCommunication.Profinet.GE 命名空间](db594ebf-c4ae-a11c-a82c-72625c3b6ab4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## BuildWriteCommand 方法 (Int64, GeSRTPAddress, Byte[])

[原文連結](http://api.hslcommunication.cn/html/e4d4aae5-bf38-f48e-0f40-694198d4e944.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.GE](../html/db594ebf-c4ae-a11c-a82c-72625c3b6ab4.htm "HslCommunication.Profinet.GE")

[GeHelper 类](../html/3c60bc6c-29f1-ae7f-d703-4bfc7c463e98.htm "GeHelper 类")

[GeHelper 方法](../html/59587f2f-8da2-a1fb-dc9d-5f031122efbe.htm "GeHelper 方法")

[BuildWriteCommand 方法](../html/f1b77cc0-230f-bef6-7f79-e621a76b5d1d.htm "BuildWriteCommand 方法 ")

[BuildWriteCommand 方法 (Int64, GeSRTPAddress, Byte[])](../html/e4d4aae5-bf38-f48e-0f40-694198d4e944.htm "BuildWriteCommand 方法 (Int64, GeSRTPAddress, Byte[])")

[BuildWriteCommand 方法 (Int64, String, Boolean[])](../html/969ce7ae-8486-af96-782e-31df7ef63f54.htm "BuildWriteCommand 方法 (Int64, String, Boolean[])")

[BuildWriteCommand 方法 (Int64, String, Byte[])](../html/bee10f7e-5595-5100-4d76-3501ae2139dc.htm "BuildWriteCommand 方法 (Int64, String, Byte[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| GeHelperBuildWriteCommand 方法 (Int64, GeSRTPAddress, Byte) |

构建一个批量写入 byte 数组变量的报文，需要指定消息号，写入的地址，地址参照 [GeSRTPNet](52c0166b-3420-af47-8dec-f25e74515116.htm) 说明。  
To construct a message to be written into byte array variables in batches,
you need to specify the message number and write address. For the address, refer to the description of [GeSRTPNet](52c0166b-3420-af47-8dec-f25e74515116.htm).

**命名空间：**
 [HslCommunication.Profinet.GE](db594ebf-c4ae-a11c-a82c-72625c3b6ab4.htm)  
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
	long id,
	GeSRTPAddress address,
	byte[] value
)
```

```
Public Shared Function BuildWriteCommand ( 
	id As Long,
	address As GeSRTPAddress,
	value As Byte()
) As OperateResult(Of Byte())
```

```
public:
static OperateResult<array<unsigned char>^>^ BuildWriteCommand(
	long long id, 
	GeSRTPAddress^ address, 
	array<unsigned char>^ value
)
```

```
static member BuildWriteCommand : 
        id : int64 * 
        address : GeSRTPAddress * 
        value : byte[] -> OperateResult<byte[]> 
```

#### 参数

id
:   类型：SystemInt64  
    消息的序号

address
:   类型：[HslCommunication.Core.AddressGeSRTPAddress](7c65a869-5fea-2bf6-526d-6d2013323256.htm)  
    地址信息

value
:   类型：SystemByte  
    byte数组的原始数据

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
包含结果信息的报文内容

![](../icons/SectionExpanded.png)参见

#### 引用

[GeHelper 类](3c60bc6c-29f1-ae7f-d703-4bfc7c463e98.htm)

[BuildWriteCommand 重载](f1b77cc0-230f-bef6-7f79-e621a76b5d1d.htm)

[HslCommunication.Profinet.GE 命名空间](db594ebf-c4ae-a11c-a82c-72625c3b6ab4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## BuildWriteCommand 方法 (Int64, String, Boolean[])

[原文連結](http://api.hslcommunication.cn/html/969ce7ae-8486-af96-782e-31df7ef63f54.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.GE](../html/db594ebf-c4ae-a11c-a82c-72625c3b6ab4.htm "HslCommunication.Profinet.GE")

[GeHelper 类](../html/3c60bc6c-29f1-ae7f-d703-4bfc7c463e98.htm "GeHelper 类")

[GeHelper 方法](../html/59587f2f-8da2-a1fb-dc9d-5f031122efbe.htm "GeHelper 方法")

[BuildWriteCommand 方法](../html/f1b77cc0-230f-bef6-7f79-e621a76b5d1d.htm "BuildWriteCommand 方法 ")

[BuildWriteCommand 方法 (Int64, GeSRTPAddress, Byte[])](../html/e4d4aae5-bf38-f48e-0f40-694198d4e944.htm "BuildWriteCommand 方法 (Int64, GeSRTPAddress, Byte[])")

[BuildWriteCommand 方法 (Int64, String, Boolean[])](../html/969ce7ae-8486-af96-782e-31df7ef63f54.htm "BuildWriteCommand 方法 (Int64, String, Boolean[])")

[BuildWriteCommand 方法 (Int64, String, Byte[])](../html/bee10f7e-5595-5100-4d76-3501ae2139dc.htm "BuildWriteCommand 方法 (Int64, String, Byte[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| GeHelperBuildWriteCommand 方法 (Int64, String, Boolean) |

构建一个批量写入 bool 数组变量的报文，需要指定消息号，写入的地址，地址参照 [GeSRTPNet](52c0166b-3420-af47-8dec-f25e74515116.htm) 说明。  
To construct a message to be written into bool array variables in batches,
you need to specify the message number and write address. For the address, refer to the description of [GeSRTPNet](52c0166b-3420-af47-8dec-f25e74515116.htm).

**命名空间：**
 [HslCommunication.Profinet.GE](db594ebf-c4ae-a11c-a82c-72625c3b6ab4.htm)  
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
	long id,
	string address,
	bool[] value
)
```

```
Public Shared Function BuildWriteCommand ( 
	id As Long,
	address As String,
	value As Boolean()
) As OperateResult(Of Byte())
```

```
public:
static OperateResult<array<unsigned char>^>^ BuildWriteCommand(
	long long id, 
	String^ address, 
	array<bool>^ value
)
```

```
static member BuildWriteCommand : 
        id : int64 * 
        address : string * 
        value : bool[] -> OperateResult<byte[]> 
```

#### 参数

id
:   类型：SystemInt64  
    消息的序号

address
:   类型：SystemString  
    地址信息

value
:   类型：SystemBoolean  
    bool数组

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
包含结果信息的报文内容

![](../icons/SectionExpanded.png)参见

#### 引用

[GeHelper 类](3c60bc6c-29f1-ae7f-d703-4bfc7c463e98.htm)

[BuildWriteCommand 重载](f1b77cc0-230f-bef6-7f79-e621a76b5d1d.htm)

[HslCommunication.Profinet.GE 命名空间](db594ebf-c4ae-a11c-a82c-72625c3b6ab4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## BuildWriteCommand 方法 (Int64, String, Byte[])

[原文連結](http://api.hslcommunication.cn/html/bee10f7e-5595-5100-4d76-3501ae2139dc.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.GE](../html/db594ebf-c4ae-a11c-a82c-72625c3b6ab4.htm "HslCommunication.Profinet.GE")

[GeHelper 类](../html/3c60bc6c-29f1-ae7f-d703-4bfc7c463e98.htm "GeHelper 类")

[GeHelper 方法](../html/59587f2f-8da2-a1fb-dc9d-5f031122efbe.htm "GeHelper 方法")

[BuildWriteCommand 方法](../html/f1b77cc0-230f-bef6-7f79-e621a76b5d1d.htm "BuildWriteCommand 方法 ")

[BuildWriteCommand 方法 (Int64, GeSRTPAddress, Byte[])](../html/e4d4aae5-bf38-f48e-0f40-694198d4e944.htm "BuildWriteCommand 方法 (Int64, GeSRTPAddress, Byte[])")

[BuildWriteCommand 方法 (Int64, String, Boolean[])](../html/969ce7ae-8486-af96-782e-31df7ef63f54.htm "BuildWriteCommand 方法 (Int64, String, Boolean[])")

[BuildWriteCommand 方法 (Int64, String, Byte[])](../html/bee10f7e-5595-5100-4d76-3501ae2139dc.htm "BuildWriteCommand 方法 (Int64, String, Byte[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| GeHelperBuildWriteCommand 方法 (Int64, String, Byte) |

构建一个批量写入 byte 数组变量的报文，需要指定消息号，写入的地址，地址参照 [GeSRTPNet](52c0166b-3420-af47-8dec-f25e74515116.htm) 说明。  
To construct a message to be written into byte array variables in batches,
you need to specify the message number and write address. For the address, refer to the description of [GeSRTPNet](52c0166b-3420-af47-8dec-f25e74515116.htm).

**命名空间：**
 [HslCommunication.Profinet.GE](db594ebf-c4ae-a11c-a82c-72625c3b6ab4.htm)  
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
	long id,
	string address,
	byte[] value
)
```

```
Public Shared Function BuildWriteCommand ( 
	id As Long,
	address As String,
	value As Byte()
) As OperateResult(Of Byte())
```

```
public:
static OperateResult<array<unsigned char>^>^ BuildWriteCommand(
	long long id, 
	String^ address, 
	array<unsigned char>^ value
)
```

```
static member BuildWriteCommand : 
        id : int64 * 
        address : string * 
        value : byte[] -> OperateResult<byte[]> 
```

#### 参数

id
:   类型：SystemInt64  
    消息的序号

address
:   类型：SystemString  
    地址信息

value
:   类型：SystemByte  
    byte数组的原始数据

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
包含结果信息的报文内容

![](../icons/SectionExpanded.png)参见

#### 引用

[GeHelper 类](3c60bc6c-29f1-ae7f-d703-4bfc7c463e98.htm)

[BuildWriteCommand 重载](f1b77cc0-230f-bef6-7f79-e621a76b5d1d.htm)

[HslCommunication.Profinet.GE 命名空间](db594ebf-c4ae-a11c-a82c-72625c3b6ab4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ExtraDateTime 方法 

[原文連結](http://api.hslcommunication.cn/html/25f30152-fad8-5f60-5176-7dca1edfccea.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.GE](../html/db594ebf-c4ae-a11c-a82c-72625c3b6ab4.htm "HslCommunication.Profinet.GE")

[GeHelper 类](../html/3c60bc6c-29f1-ae7f-d703-4bfc7c463e98.htm "GeHelper 类")

[GeHelper 方法](../html/59587f2f-8da2-a1fb-dc9d-5f031122efbe.htm "GeHelper 方法")

[BuildReadCommand 方法](../html/9ff3c313-5f82-78f1-ee0c-ebd401fc06eb.htm "BuildReadCommand 方法 ")

[BuildReadCoreCommand 方法](../html/243b844b-95fe-ce64-3e8d-29f5e9638098.htm "BuildReadCoreCommand 方法 ")

[BuildWriteCommand 方法](../html/f1b77cc0-230f-bef6-7f79-e621a76b5d1d.htm "BuildWriteCommand 方法 ")

[ExtraDateTime 方法](../html/25f30152-fad8-5f60-5176-7dca1edfccea.htm "ExtraDateTime 方法 ")

[ExtraProgramName 方法](../html/bddf3623-efe4-07e3-130a-54c019072cc7.htm "ExtraProgramName 方法 ")

[ExtraResponseContent 方法](../html/62346044-4399-25c6-f5e1-62fb52f02957.htm "ExtraResponseContent 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| GeHelperExtraDateTime 方法 |

从实际的时间的字节数组里解析出C#格式的时间对象，这个时间可能是时区0的时间，需要自行转化本地时间。  
Analyze the time object in C# format from the actual time byte array.
This time may be the time in time zone 0, and you need to convert the local time yourself.

**命名空间：**
 [HslCommunication.Profinet.GE](db594ebf-c4ae-a11c-a82c-72625c3b6ab4.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult<DateTime> ExtraDateTime(
	byte[] content
)
```

```
Public Shared Function ExtraDateTime ( 
	content As Byte()
) As OperateResult(Of DateTime)
```

```
public:
static OperateResult<DateTime>^ ExtraDateTime(
	array<unsigned char>^ content
)
```

```
static member ExtraDateTime : 
        content : byte[] -> OperateResult<DateTime> 
```

#### 参数

content
:   类型：SystemByte  
    字节数组

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)DateTime  
包含是否成功的结果对象

![](../icons/SectionExpanded.png)参见

#### 引用

[GeHelper 类](3c60bc6c-29f1-ae7f-d703-4bfc7c463e98.htm)

[HslCommunication.Profinet.GE 命名空间](db594ebf-c4ae-a11c-a82c-72625c3b6ab4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ExtraProgramName 方法 

[原文連結](http://api.hslcommunication.cn/html/bddf3623-efe4-07e3-130a-54c019072cc7.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.GE](../html/db594ebf-c4ae-a11c-a82c-72625c3b6ab4.htm "HslCommunication.Profinet.GE")

[GeHelper 类](../html/3c60bc6c-29f1-ae7f-d703-4bfc7c463e98.htm "GeHelper 类")

[GeHelper 方法](../html/59587f2f-8da2-a1fb-dc9d-5f031122efbe.htm "GeHelper 方法")

[BuildReadCommand 方法](../html/9ff3c313-5f82-78f1-ee0c-ebd401fc06eb.htm "BuildReadCommand 方法 ")

[BuildReadCoreCommand 方法](../html/243b844b-95fe-ce64-3e8d-29f5e9638098.htm "BuildReadCoreCommand 方法 ")

[BuildWriteCommand 方法](../html/f1b77cc0-230f-bef6-7f79-e621a76b5d1d.htm "BuildWriteCommand 方法 ")

[ExtraDateTime 方法](../html/25f30152-fad8-5f60-5176-7dca1edfccea.htm "ExtraDateTime 方法 ")

[ExtraProgramName 方法](../html/bddf3623-efe4-07e3-130a-54c019072cc7.htm "ExtraProgramName 方法 ")

[ExtraResponseContent 方法](../html/62346044-4399-25c6-f5e1-62fb52f02957.htm "ExtraResponseContent 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| GeHelperExtraProgramName 方法 |

从实际的时间的字节数组里解析出PLC的程序的名称。  
Parse the name of the PLC program from the actual time byte array

**命名空间：**
 [HslCommunication.Profinet.GE](db594ebf-c4ae-a11c-a82c-72625c3b6ab4.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult<string> ExtraProgramName(
	byte[] content
)
```

```
Public Shared Function ExtraProgramName ( 
	content As Byte()
) As OperateResult(Of String)
```

```
public:
static OperateResult<String^>^ ExtraProgramName(
	array<unsigned char>^ content
)
```

```
static member ExtraProgramName : 
        content : byte[] -> OperateResult<string> 
```

#### 参数

content
:   类型：SystemByte  
    字节数组

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)String  
包含是否成功的结果对象

![](../icons/SectionExpanded.png)参见

#### 引用

[GeHelper 类](3c60bc6c-29f1-ae7f-d703-4bfc7c463e98.htm)

[HslCommunication.Profinet.GE 命名空间](db594ebf-c4ae-a11c-a82c-72625c3b6ab4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ExtraResponseContent 方法 

[原文連結](http://api.hslcommunication.cn/html/62346044-4399-25c6-f5e1-62fb52f02957.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.GE](../html/db594ebf-c4ae-a11c-a82c-72625c3b6ab4.htm "HslCommunication.Profinet.GE")

[GeHelper 类](../html/3c60bc6c-29f1-ae7f-d703-4bfc7c463e98.htm "GeHelper 类")

[GeHelper 方法](../html/59587f2f-8da2-a1fb-dc9d-5f031122efbe.htm "GeHelper 方法")

[BuildReadCommand 方法](../html/9ff3c313-5f82-78f1-ee0c-ebd401fc06eb.htm "BuildReadCommand 方法 ")

[BuildReadCoreCommand 方法](../html/243b844b-95fe-ce64-3e8d-29f5e9638098.htm "BuildReadCoreCommand 方法 ")

[BuildWriteCommand 方法](../html/f1b77cc0-230f-bef6-7f79-e621a76b5d1d.htm "BuildWriteCommand 方法 ")

[ExtraDateTime 方法](../html/25f30152-fad8-5f60-5176-7dca1edfccea.htm "ExtraDateTime 方法 ")

[ExtraProgramName 方法](../html/bddf3623-efe4-07e3-130a-54c019072cc7.htm "ExtraProgramName 方法 ")

[ExtraResponseContent 方法](../html/62346044-4399-25c6-f5e1-62fb52f02957.htm "ExtraResponseContent 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| GeHelperExtraResponseContent 方法 |

从PLC返回的数据中，提取出实际的数据内容，最少6个字节的数据。超出实际的数据长度的部分没有任何意义。  
From the data returned by the PLC, extract the actual data content, at least 6 bytes of data. The part beyond the actual data length has no meaning.

**命名空间：**
 [HslCommunication.Profinet.GE](db594ebf-c4ae-a11c-a82c-72625c3b6ab4.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult<byte[]> ExtraResponseContent(
	byte[] content
)
```

```
Public Shared Function ExtraResponseContent ( 
	content As Byte()
) As OperateResult(Of Byte())
```

```
public:
static OperateResult<array<unsigned char>^>^ ExtraResponseContent(
	array<unsigned char>^ content
)
```

```
static member ExtraResponseContent : 
        content : byte[] -> OperateResult<byte[]> 
```

#### 参数

content
:   类型：SystemByte  
    PLC返回的数据信息

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
解析后的实际数据内容

![](../icons/SectionExpanded.png)参见

#### 引用

[GeHelper 类](3c60bc6c-29f1-ae7f-d703-4bfc7c463e98.htm)

[HslCommunication.Profinet.GE 命名空间](db594ebf-c4ae-a11c-a82c-72625c3b6ab4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GeSRTPNet 类

[原文連結](http://api.hslcommunication.cn/html/52c0166b-3420-af47-8dec-f25e74515116.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.GE](../html/db594ebf-c4ae-a11c-a82c-72625c3b6ab4.htm "HslCommunication.Profinet.GE")

[GeSRTPNet 类](../html/52c0166b-3420-af47-8dec-f25e74515116.htm "GeSRTPNet 类")

[GeSRTPNet 构造函数](../html/ac855f75-4e0d-5f51-4931-522830c9ec8e.htm "GeSRTPNet 构造函数 ")

[GeSRTPNet 属性](../html/e7d24106-6d02-0567-161e-a5fc6449bd27.htm "GeSRTPNet 属性")

[GeSRTPNet 方法](../html/a589d191-964f-8e25-3d96-9d59e2371b59.htm "GeSRTPNet 方法")

[GeSRTPNet 字段](../html/a1067a5b-0a99-e58d-2ded-087018368937.htm "GeSRTPNet 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| GeSRTPNet 类 |

Ge(通用电气)集团的SRTP协议实现的客户端，支持 I,Q,M,T,SA,SB,SC,S,G 的位和字节读写，支持 AI,AQ,R 的字读写操作，支持读取PLC时间，程序名操作  
The client implemented by Ge (General Electric) Group's SRTP protocol, supports I, Q, M, T, SA, SB, SC, S, G bit and byte read and write,
supports AI, AQ, R word read and write Operation, support reading PLC time, program name operation

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  [HslCommunication.Core.NetBinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)  
    [HslCommunication.Core.DeviceDeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)  
      [HslCommunication.Core.DeviceDeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)  
        HslCommunication.Profinet.GEGeSRTPNet

**命名空间：**
 [HslCommunication.Profinet.GE](db594ebf-c4ae-a11c-a82c-72625c3b6ab4.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class GeSRTPNet : DeviceTcpNet
```

```
Public Class GeSRTPNet
	Inherits DeviceTcpNet
```

```
public ref class GeSRTPNet : public DeviceTcpNet
```

```
type GeSRTPNet =  
    class
        inherit DeviceTcpNet
    end
```

GeSRTPNet 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [GeSRTPNet](de142e24-e78d-4190-7648-20a4eff02c8c.htm) | 实例化一个默认的对象  Instantiate a default object |
| 公共方法 | [GeSRTPNet(String, Int32)](1c89aa71-d5c2-6650-b157-917100827ccc.htm) | 指定IP地址和端口号来实例化一个对象  Specify the IP address and port number to instantiate an object |

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
| 受保护的属性 | [WordLength](f1a505cc-e85c-6c8b-b189-58a1ab5753f9.htm) | 一个字单位的数据表示的地址长度，西门子为2，三菱，欧姆龙，modbusTcp就为1，AB PLC无效  The address length represented by one word of data, Siemens is 2, Mitsubishi, Omron, modbusTcp is 1, AB PLC is invalid (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |

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
| 受保护的方法 | [GetNewNetMessage](925f55de-2df6-75c6-df9f-753de861bebc.htm) | 获取一个新的消息对象的方法，需要在继承类里面进行重写  The method to get a new message object needs to be overridden in the inheritance class (重写 [BinaryCommunicationGetNewNetMessage](16636a1b-a6f0-9c0f-9d9c-c94e26c67459.htm).) |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法 | [GetWordLength](c27821ea-1739-8eb9-a848-23a00828fec8.htm) | 一个字单位的数据表示的地址长度，西门子为2，三菱，欧姆龙，modbusTcp就为1，AB PLC无效  The address length represented by one word of data, Siemens is 2, Mitsubishi, Omron, modbusTcp is 1, AB PLC is invalid (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 受保护的方法代码示例 | [InitializationOnConnect](9364af5f-5b71-458f-44b6-1b26ac09fac7.htm) | 根据实际的协议选择是否重写本方法，有些协议在创建连接之后，需要进行一些初始化的信号握手，才能最终建立网络通道。  Whether to rewrite this method is based on the actual protocol. Some protocols require some initial signal handshake to establish a network channel after the connection is created. (重写 [BinaryCommunicationInitializationOnConnect](4ab6989f-a5b5-05fe-8b44-bfa45824a2a7.htm).) |
| 受保护的方法代码示例 | [InitializationOnConnectAsync](23faa44d-2cc1-80c9-b6bc-995a2bdb8ab2.htm) | 根据实际的协议选择是否重写本方法，有些协议在创建连接之后，需要进行一些初始化的信号握手，才能最终建立网络通道。  Whether to rewrite this method is based on the actual protocol. Some protocols require some initial signal handshake to establish a network channel after the connection is created. (重写 [BinaryCommunicationInitializationOnConnectAsync](7da2af04-d4fc-5045-19a4-4b0d1722eef3.htm).) |
| 公共方法 | [IpAddressPing](0ac126b4-a830-a822-4ef5-855474942cd7.htm) | 对当前设备的IP地址进行PING的操作，返回PING的结果，正常来说，返回Success  PING the IP address of the current device and return the PING result. Normally, it returns Success (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 受保护的方法 | [LogRevcMessage(Byte)](dab1657f-1e75-7fe1-7003-ccedfc69855d.htm) | 使用日志记录一个接收的报文信息  Logs are used to record information about a received packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [LogRevcMessage(Byte, PipeSession)](0f51ae45-575e-4228-f978-01b37d6d5df1.htm) | 使用日志记录一个接收的报文信息  Logs are used to record information about a received packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [LogSendMessage(Byte)](4f90eecc-68a2-1ef2-44ef-369f19bed3ec.htm) | 使用日志记录一个发送的报文信息  Logs are used to record information about a send packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [LogSendMessage(Byte, PipeSession)](209caadf-4f0f-40aa-04c2-ebc61f8dab13.htm) | 使用日志记录一个发送的报文信息  Logs are used to record information about a send packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法 | [PackCommandWithHeader](d0b2cd76-c306-11fc-38c7-af8ad306cce7.htm) | 对当前的命令进行打包处理，通常是携带命令头内容，标记当前的命令的长度信息，需要进行重写，否则默认不打包  The current command is packaged, usually carrying the content of the command header, marking the length of the current command, and it needs to be rewritten, otherwise it is not packaged by default (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [Read(String, UInt16)](92534a5c-499e-278f-005e-57f05207e689.htm) | 批量读取字节数组信息，需要指定地址和长度，返回原始的字节数组，支持 I,Q,M,T,SA,SB,SC,S,G 的位和字节读写，支持 AI,AQ,R 的字读写操作，地址示例：R1,M1  Batch read byte array information, need to specify the address and length, return the original byte array. Support I, Q, M, T, SA, SB, SC, S, G bit and byte read and write, support AI, AQ, R word read and write operations, address examples: R1, M1 (重写 [DeviceCommunicationRead(String, UInt16)](9dd09452-c200-bc65-75ff-69aa34b366d4.htm).) |
| 公共方法代码示例 | [ReadT](95abe064-5594-5ae7-a849-98c53f88fee0.htm) | 读取支持Hsl特性的数据内容，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明，或是官网的详细文档  Read the data content of the Hsl attribute. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadAsync(String, UInt16)](176d000d-4655-22b6-5610-1348a997f8e5.htm) | 批量读取字节数组信息，需要指定地址和长度，返回原始的字节数组，支持 I,Q,M,T,SA,SB,SC,S,G 的位和字节读写，支持 AI,AQ,R 的字读写操作，地址示例：R1,M1  Batch read byte array information, need to specify the address and length, return the original byte array. Support I, Q, M, T, SA, SB, SC, S, G bit and byte read and write, support AI, AQ, R word read and write operations, address examples: R1, M1 (重写 [DeviceCommunicationReadAsync(String, UInt16)](d8c14d37-56dc-2a31-1a19-3e2b5b65197d.htm).) |
| 公共方法代码示例 | [ReadAsyncT](0522e740-4bf0-2e31-c4c5-e41384321b37.htm) | 异步读取支持Hsl特性的数据内容，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Asynchronously read the data content of the Hsl attribute. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadBool(String)](9b2faa2e-eaa9-7ace-c49f-7a0318f1df1a.htm) | 读取单个的Boolean数据信息  Read a single Boolean data message (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadBool(String, UInt16)](cfc5e2c8-62cc-a442-48cf-3acc4a2d8261.htm) | 按照位为单位，批量从指定的地址里读取 bool 数组数据，不支持 AI，AQ，R 地址类型，地址比如从1开始，例如 I1,Q1,M1,T1,SA1,SB1,SC1,S1,G1  In units of bits, read bool array data from the specified address in batches. AI, AQ, R address types are not supported. For example, the address starts from 1, such as I1, Q1, M1, T1, SA1, SB1, SC1, S1, G1 (重写 [DeviceCommunicationReadBool(String, UInt16)](69d0d24f-34a7-5b41-5c1b-619edc7e6b54.htm).) |
| 公共方法 | [ReadBoolAsync(String)](cc37da1f-d673-fe67-78ee-4e465d29a2bd.htm) | 异步读取单个的Boolean数据信息  Asynchronously read a single Boolean data message (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadBoolAsync(String, UInt16)](b079c820-bf03-3d5d-bff8-fd980f4c68f1.htm) | 按照位为单位，批量从指定的地址里读取 bool 数组数据，不支持 AI，AQ，R 地址类型，地址比如从1开始，例如 I1,Q1,M1,T1,SA1,SB1,SC1,S1,G1  In units of bits, read bool array data from the specified address in batches. AI, AQ, R address types are not supported. For example, the address starts from 1, such as I1, Q1, M1, T1, SA1, SB1, SC1, S1, G1 (重写 [DeviceCommunicationReadBoolAsync(String, UInt16)](a45d383d-0095-3ea7-32f8-e75af88d60f8.htm).) |
| 公共方法 | [ReadByte](c477231b-11fd-ddde-e293-5d5a3b7cd473.htm) | 根据指定的地址来读取一个字节的数据，按照字节为单位，例如 M1 字节，就是指 M1-M8 位组成的字节，M2 字节就是 M9-M16 组成的字节。不支持对 AI,AQ,R 寄存器的字节读取  A byte of data is read according to the specified address, and the unit is byte. For example, M1 byte refers to a byte composed of M1-M8 bits, and M2 byte is a byte composed of M9-M16. Does not support byte reading of AI, AQ, R registers |
| 公共方法 | [ReadByteAsync](e6e6e840-c1a3-b92e-25bf-609016104961.htm) | 根据指定的地址来读取一个字节的数据，按照字节为单位，例如 M1 字节，就是指 M1-M8 位组成的字节，M2 字节就是 M9-M16 组成的字节。不支持对 AI,AQ,R 寄存器的字节读取  A byte of data is read according to the specified address, and the unit is byte. For example, M1 byte refers to a byte composed of M1-M8 bits, and M2 byte is a byte composed of M9-M16. Does not support byte reading of AI, AQ, R registers |
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
| 公共方法 | [ReadPLCTime](cd7a4ead-ed5c-4cab-9821-e98f2c6e80aa.htm) | **[商业授权]** 读取PLC当前的时间，这个时间可能是不包含时区的，需要自己转换成本地的时间。 **[Authorization]** Read the current time of the PLC, this time may not include the time zone, you need to convert the local time yourself. |
| 公共方法 | [ReadPLCTimeAsync](194f0501-71c1-2117-5b63-696759d21a80.htm) | **[商业授权]** 读取PLC当前的时间，这个时间可能是不包含时区的，需要自己转换成本地的时间。 **[Authorization]** Read the current time of the PLC, this time may not include the time zone, you need to convert the local time yourself. |
| 公共方法 | [ReadProgramName](394d566d-cf52-338e-8f3e-e8c63d28f4a4.htm) | **[商业授权]** 读取PLC当前的程序的名称 **[Authorization]** Read the name of the current program of the PLC |
| 公共方法 | [ReadProgramNameAsync](0f5ad8db-9e52-7c29-7626-dc894dfec490.htm) | **[商业授权]** 读取PLC当前的程序的名称 **[Authorization]** Read the name of the current program of the PLC |
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
| 公共方法 | [SetDtuPipe](f765d1d7-234a-b229-9457-97e25e024dc9.htm) | 将当前的通信对象设置DTU模式，允许传入现成的管道，并返回初始化结果，如果该设备重写了握手报文，就是返回握手结果  Set the current communication object to DTU mode, allow the existing pipe to be passed in, and return the initialization result, if the device rewrites the handshake packet, the handshake result is returned (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [SetDtuPipeAsync](9d1fd12d-b804-3828-a627-a2c0eef1d23f.htm) | 将当前的通信对象设置DTU模式，允许传入现成的管道，并返回初始化结果，如果该设备重写了握手报文，就是返回握手结果  Set the current communication object to DTU mode, allow the existing pipe to be passed in, and return the initialization result, if the device rewrites the handshake packet, the handshake result is returned (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [SetPersistentConnection](f903dae0-4a4e-c7db-8350-50e36b257ec5.htm) | **已过时。** V11版本及之前设置长连接的方法，在V12版本以上中没有任何效果，默认长连接，删除调用即可，此处保留方法是为了部分用户保持兼容性升级。  The method of setting the long connection in V11 and before, has no effect in V12 and above. this method can be deleted. The method is retained here to maintain compatibility upgrades for some users. (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共方法 | [ToString](ef49b930-7974-80fa-c392-84830ae0ee83.htm) | (重写 [DeviceTcpNetToString](209a196b-90ea-2b73-e915-ce4b11f1263d.htm).) |
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
| 公共方法 | [Write(String, Boolean)](18c50ec3-9b8d-7cde-0429-5577a95562f5.htm) | 按照位为单位，批量写入 bool 数组到指定的地址里，不支持 AI，AQ，R 地址类型，地址比如从1开始，例如 I1,Q1,M1,T1,SA1,SB1,SC1,S1,G1  In units of bits, write bool arrays in batches to the specified addresses. AI, AQ, and R address types are not supported. For example, the address starts from 1, such as I1, Q1, M1, T1, SA1, SB1, SC1, S1, G1 (重写 [DeviceCommunicationWrite(String, Boolean)](5b046c82-3d8e-550b-fdf3-41ca5558169b.htm).) |
| 公共方法 | [Write(String, Byte)](9a0ec8e1-6d5c-3e0e-1e0a-3d3534043d62.htm) | 向PLC中写入byte数据，返回是否写入成功  Write byte data to PLC and return whether the writing is successful |
| 公共方法 | [Write(String, Byte)](f958948c-4214-7418-1525-b7bf0a236a9e.htm) | 写入原始的byte数组数据到指定的地址，返回是否写入成功  Write the original byte array data to the specified address, and return whether the write was successful (重写 [DeviceCommunicationWrite(String, Byte)](b5536dd6-ccde-6883-ac27-bff9e4b37806.htm).) |
| 公共方法代码示例 | [Write(String, String, Int32)](02fe2ca9-49a3-ede4-d1af-8c6da1f644da.htm) | 写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为ASCII  Write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, it is filled with 0 and the encoding is ASCII. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, String, Encoding)](7b602661-902d-2b70-1fdb-ba87ecd307bf.htm) | 写入字符串信息，需要指定的编码信息  Write string information, need to specify the encoding information (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, String, Int32, Encoding)](e8cb30d9-4c5e-6c68-b9af-1b1b6c2dc714.htm) | 写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为指定的编码信息  Write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, then the operation is complemented with 0 , you should specified the encoding information (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteT(T)](9f48853f-fbeb-c37d-25a2-1c76d1dee407.htm) | 写入支持Hsl特性的数据，返回是否写入成功，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Write data that supports the Hsl attribute, and return whether the write was successful. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WriteAsync(String, Boolean)](10fe669c-805e-1e5a-f016-f319ad363e90.htm) | 异步批量写入Boolean数组数据，返回是否成功  Asynchronously batch write Boolean array data, return success (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
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
| 公共方法 | [WriteAsync(String, Boolean)](c10ddba4-b701-b34c-2d58-cd94e308de96.htm) | 按照位为单位，批量写入 bool 数组到指定的地址里，不支持 AI，AQ，R 地址类型，地址比如从1开始，例如 I1,Q1,M1,T1,SA1,SB1,SC1,S1,G1  In units of bits, write bool arrays in batches to the specified addresses. AI, AQ, and R address types are not supported. For example, the address starts from 1, such as I1, Q1, M1, T1, SA1, SB1, SC1, S1, G1 (重写 [DeviceCommunicationWriteAsync(String, Boolean)](fe872dae-7ef9-12c1-d883-e92ba8676416.htm).) |
| 公共方法 | [WriteAsync(String, Byte)](604f4cfc-50c1-0119-1320-fe91e83a1850.htm) | 向PLC中写入byte数据，返回是否写入成功  Write byte data to PLC and return whether the writing is successful |
| 公共方法 | [WriteAsync(String, Byte)](98119c8f-0d6b-fa09-7403-10fedd11e528.htm) | 写入原始的byte数组数据到指定的地址，返回是否写入成功  Write the original byte array data to the specified address, and return whether the write was successful (重写 [DeviceCommunicationWriteAsync(String, Byte)](8f1cec30-8b86-6749-c8f3-7d39168d59df.htm).) |
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

![](../icons/SectionExpanded.png)备注

PLC的端口号默认18245，其中读取R，AI，AQ寄存器的原始字节时，传入的长度参数为字节长度。  
对其他寄存器而言，M1-M8的位读取，相当于 M1的字节读取。写入也是同理。

![](../icons/SectionExpanded.png)示例

| 重要事项 **重要事项** |
| --- |
| 对于 AI,AQ,R 类型的数据来说，底层是以字为单位的，但是调用 [Read(String, UInt16)](92534a5c-499e-278f-005e-57f05207e689.htm) 来读取原始数据时，需要传入2被的字长度，对于其他方法不影响 |

简单的短连接使用

[复制](# "复制")

```
// 实例化对象，指定PLC的ip地址
GeSRTPNet ge = new GeSRTPNet( "192.168.1.110" );
// 举例读取R1的值
short R1 = ge.ReadInt16( "R1" ).Content;
```

简单的长连接使用

[复制](# "复制")

```
// 实例化对象，指定PLC的ip地址和端口号
GeSRTPNet ge = new GeSRTPNet( "192.168.1.110" );

// 连接对象
OperateResult connect = ge.ConnectServer( );
if (!connect.IsSuccess)
{
    Console.WriteLine( "connect failed:" + connect.Message );
    return;
}

// 举例读取M100的值
bool M100 = ge.ReadBool( "M100" ).Content;

// 如果需要判断是否读取成功
OperateResult<bool> read = ge.ReadBool( "M100" );
if (!read.IsSuccess)
{
    Console.WriteLine( "read failed:" + read.Message );
    return;
}
Console.WriteLine( "read value:" + read.Content );

// 写入的原理是一样的
OperateResult write = ge.Write( "M100", true );
if (!write.IsSuccess)
{
    Console.WriteLine( "write failed:" + write.Message );
    return;
}

// 如果我读取 M100 的字节这个数据是什么呢？
OperateResult<byte> readByte = ge.ReadByte( "M100" );
if (!readByte.IsSuccess)
{
    Console.WriteLine( "read failed:" + readByte.Message );
    return;
}

// 按照字节为单位读取 M100 的值，那按照位来说，就是 (100-1)*8+1= M793 ~ M800 这8个位组成的byte值
Console.WriteLine( "read value:" + readByte.Content );
// 如果我要重新解析成bool数组，方法参照 HslCommunication 提供的扩展方法 ToBoolArray
bool[] boolarray = new byte[] { readByte.Content }.ToBoolArray( );
bool m793 = boolarray[0];
bool m794 = boolarray[1];
bool m795 = boolarray[2];
bool m796 = boolarray[3];
bool m797 = boolarray[4];
bool m798 = boolarray[5];
bool m799 = boolarray[6];
bool m800 = boolarray[7];
// 反之，写入也是同样的原理


// 使用完毕之后，关闭连接
ge.ConnectClose( );
```

假设起始地址为R1，R1存储了温度，100.6℃值为1006，R2存储了压力，1.23Mpa值为123，R3-R4存储了产量计数，读取如下：

Read示例

[复制](# "复制")

```
GeSRTPNet ge = new GeSRTPNet( "192.168.1.110" );

OperateResult<byte[]> read = ge.Read( "R1", 8 );
if (read.IsSuccess)
{
    float temp  = ge.ByteTransform.TransInt16( read.Content, 0 ) / 10f;
    float press = ge.ByteTransform.TransInt16( read.Content, 2 ) / 100f;
    int count   = ge.ByteTransform.TransInt32( read.Content, 2 );

    // do something
}
else
{
    // failed
}
```

以下是读取不同类型数据的示例

Read示例

[复制](# "复制")

```
GeSRTPNet ge = new GeSRTPNet( "192.168.1.110" );

// 此处以R1寄存器作为示例
short short_R1    = ge.ReadInt16( "R1" ).Content;     // 读取R1组成的字
ushort ushort_R1  = ge.ReadUInt16( "R1" ).Content;    // 读取R1组成的无符号的值
int int_R1        = ge.ReadInt32( "R1" ).Content;     // 读取R1-R2组成的有符号的数据
uint uint_R1      = ge.ReadUInt32( "R1" ).Content;    // 读取R1-R2组成的无符号的值
float float_R1    = ge.ReadFloat( "R1" ).Content;     // 读取R1-R2组成的单精度值
long long_R1      = ge.ReadInt64( "R1" ).Content;     // 读取R1-R4组成的大数据值
ulong ulong_R1    = ge.ReadUInt64( "R1" ).Content;    // 读取R1-R4组成的无符号大数据
double double_R1  = ge.ReadDouble( "R1" ).Content;    // 读取R1-R4组成的双精度值
string str_R1     = ge.ReadString( "R1", 10 ).Content;// 读取R1-R5组成的ASCII字符串数据

// 读取数组
short[] short_R1_array      = ge.ReadInt16( "R1", 10 ).Content;   // 读取R1-R10组成的数组
ushort[] ushort_R1_array    = ge.ReadUInt16( "R1", 10 ).Content;  // 读取R1-R10组成的无符号的值的数组
int[] int_R1_array          = ge.ReadInt32( "R1", 10 ).Content;   // 读取R1-R20组成的有符号的数据的数组
uint[] uint_R1_array        = ge.ReadUInt32( "R1", 10 ).Content;  // 读取R1-R20组成的无符号的值的数组
float[] float_R1_array      = ge.ReadFloat( "R1", 10 ).Content;   // 读取R1-R20组成的单精度值的数组
long[] long_R1_array        = ge.ReadInt64( "R1", 10 ).Content;   // 读取R1-R40组成的大数据值的数组
ulong[] ulong_R1_array      = ge.ReadUInt64( "R1", 10 ).Content;  // 读取R1-R40组成的无符号大数据的数组
double[] double_R1_array    = ge.ReadDouble( "R1", 10 ).Content;  // 读取R1-R40组成的双精度值的数组
```

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.Profinet.GE 命名空间](db594ebf-c4ae-a11c-a82c-72625c3b6ab4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GeSRTPNet 构造函数 

[原文連結](http://api.hslcommunication.cn/html/ac855f75-4e0d-5f51-4931-522830c9ec8e.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.GE](../html/db594ebf-c4ae-a11c-a82c-72625c3b6ab4.htm "HslCommunication.Profinet.GE")

[GeSRTPNet 类](../html/52c0166b-3420-af47-8dec-f25e74515116.htm "GeSRTPNet 类")

[GeSRTPNet 构造函数](../html/ac855f75-4e0d-5f51-4931-522830c9ec8e.htm "GeSRTPNet 构造函数 ")

[GeSRTPNet 构造函数](../html/de142e24-e78d-4190-7648-20a4eff02c8c.htm "GeSRTPNet 构造函数 ")

[GeSRTPNet 构造函数 (String, Int32)](../html/1c89aa71-d5c2-6650-b157-917100827ccc.htm "GeSRTPNet 构造函数 (String, Int32)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| GeSRTPNet 构造函数 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [GeSRTPNet](de142e24-e78d-4190-7648-20a4eff02c8c.htm) | 实例化一个默认的对象  Instantiate a default object |
| 公共方法 | [GeSRTPNet(String, Int32)](1c89aa71-d5c2-6650-b157-917100827ccc.htm) | 指定IP地址和端口号来实例化一个对象  Specify the IP address and port number to instantiate an object |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[GeSRTPNet 类](52c0166b-3420-af47-8dec-f25e74515116.htm)

[HslCommunication.Profinet.GE 命名空间](db594ebf-c4ae-a11c-a82c-72625c3b6ab4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GeSRTPNet 构造函数 

[原文連結](http://api.hslcommunication.cn/html/de142e24-e78d-4190-7648-20a4eff02c8c.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.GE](../html/db594ebf-c4ae-a11c-a82c-72625c3b6ab4.htm "HslCommunication.Profinet.GE")

[GeSRTPNet 类](../html/52c0166b-3420-af47-8dec-f25e74515116.htm "GeSRTPNet 类")

[GeSRTPNet 构造函数](../html/ac855f75-4e0d-5f51-4931-522830c9ec8e.htm "GeSRTPNet 构造函数 ")

[GeSRTPNet 构造函数](../html/de142e24-e78d-4190-7648-20a4eff02c8c.htm "GeSRTPNet 构造函数 ")

[GeSRTPNet 构造函数 (String, Int32)](../html/1c89aa71-d5c2-6650-b157-917100827ccc.htm "GeSRTPNet 构造函数 (String, Int32)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| GeSRTPNet 构造函数 |

实例化一个默认的对象  
Instantiate a default object

**命名空间：**
 [HslCommunication.Profinet.GE](db594ebf-c4ae-a11c-a82c-72625c3b6ab4.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public GeSRTPNet()
```

```
Public Sub New
```

```
public:
GeSRTPNet()
```

```
new : unit -> GeSRTPNet
```

![](../icons/SectionExpanded.png)参见

#### 引用

[GeSRTPNet 类](52c0166b-3420-af47-8dec-f25e74515116.htm)

[GeSRTPNet 重载](ac855f75-4e0d-5f51-4931-522830c9ec8e.htm)

[HslCommunication.Profinet.GE 命名空间](db594ebf-c4ae-a11c-a82c-72625c3b6ab4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GeSRTPNet 构造函数 (String, Int32)

[原文連結](http://api.hslcommunication.cn/html/1c89aa71-d5c2-6650-b157-917100827ccc.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.GE](../html/db594ebf-c4ae-a11c-a82c-72625c3b6ab4.htm "HslCommunication.Profinet.GE")

[GeSRTPNet 类](../html/52c0166b-3420-af47-8dec-f25e74515116.htm "GeSRTPNet 类")

[GeSRTPNet 构造函数](../html/ac855f75-4e0d-5f51-4931-522830c9ec8e.htm "GeSRTPNet 构造函数 ")

[GeSRTPNet 构造函数](../html/de142e24-e78d-4190-7648-20a4eff02c8c.htm "GeSRTPNet 构造函数 ")

[GeSRTPNet 构造函数 (String, Int32)](../html/1c89aa71-d5c2-6650-b157-917100827ccc.htm "GeSRTPNet 构造函数 (String, Int32)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| GeSRTPNet 构造函数 (String, Int32) |

指定IP地址和端口号来实例化一个对象  
Specify the IP address and port number to instantiate an object

**命名空间：**
 [HslCommunication.Profinet.GE](db594ebf-c4ae-a11c-a82c-72625c3b6ab4.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public GeSRTPNet(
	string ipAddress,
	int port = 18245
)
```

```
Public Sub New ( 
	ipAddress As String,
	Optional port As Integer = 18245
)
```

```
public:
GeSRTPNet(
	String^ ipAddress, 
	int port = 18245
)
```

```
new : 
        ipAddress : string * 
        ?port : int 
(* Defaults:
        let _port = defaultArg port 18245
*)
-> GeSRTPNet
```

#### 参数

ipAddress
:   类型：SystemString  
    Ip地址

port (Optional)
:   类型：SystemInt32  
    端口号

![](../icons/SectionExpanded.png)参见

#### 引用

[GeSRTPNet 类](52c0166b-3420-af47-8dec-f25e74515116.htm)

[GeSRTPNet 重载](ac855f75-4e0d-5f51-4931-522830c9ec8e.htm)

[HslCommunication.Profinet.GE 命名空间](db594ebf-c4ae-a11c-a82c-72625c3b6ab4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GeSRTPNet 属性

[原文連結](http://api.hslcommunication.cn/html/e7d24106-6d02-0567-161e-a5fc6449bd27.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.GE](../html/db594ebf-c4ae-a11c-a82c-72625c3b6ab4.htm "HslCommunication.Profinet.GE")

[GeSRTPNet 类](../html/52c0166b-3420-af47-8dec-f25e74515116.htm "GeSRTPNet 类")

[GeSRTPNet 构造函数](../html/ac855f75-4e0d-5f51-4931-522830c9ec8e.htm "GeSRTPNet 构造函数 ")

[GeSRTPNet 属性](../html/e7d24106-6d02-0567-161e-a5fc6449bd27.htm "GeSRTPNet 属性")

[GeSRTPNet 方法](../html/a589d191-964f-8e25-3d96-9d59e2371b59.htm "GeSRTPNet 方法")

[GeSRTPNet 字段](../html/a1067a5b-0a99-e58d-2ded-087018368937.htm "GeSRTPNet 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| GeSRTPNet 属性 |

[GeSRTPNet](52c0166b-3420-af47-8dec-f25e74515116.htm) 类型公开以下成员。

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
| 受保护的属性 | [WordLength](f1a505cc-e85c-6c8b-b189-58a1ab5753f9.htm) | 一个字单位的数据表示的地址长度，西门子为2，三菱，欧姆龙，modbusTcp就为1，AB PLC无效  The address length represented by one word of data, Siemens is 2, Mitsubishi, Omron, modbusTcp is 1, AB PLC is invalid (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[GeSRTPNet 类](52c0166b-3420-af47-8dec-f25e74515116.htm)

[HslCommunication.Profinet.GE 命名空间](db594ebf-c4ae-a11c-a82c-72625c3b6ab4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GeSRTPNet 方法

[原文連結](http://api.hslcommunication.cn/html/a589d191-964f-8e25-3d96-9d59e2371b59.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.GE](../html/db594ebf-c4ae-a11c-a82c-72625c3b6ab4.htm "HslCommunication.Profinet.GE")

[GeSRTPNet 类](../html/52c0166b-3420-af47-8dec-f25e74515116.htm "GeSRTPNet 类")

[GeSRTPNet 方法](../html/a589d191-964f-8e25-3d96-9d59e2371b59.htm "GeSRTPNet 方法")

[GetNewNetMessage 方法](../html/925f55de-2df6-75c6-df9f-753de861bebc.htm "GetNewNetMessage 方法 ")

[InitializationOnConnect 方法](../html/9364af5f-5b71-458f-44b6-1b26ac09fac7.htm "InitializationOnConnect 方法 ")

[InitializationOnConnectAsync 方法](../html/23faa44d-2cc1-80c9-b6bc-995a2bdb8ab2.htm "InitializationOnConnectAsync 方法 ")

[Read 方法](../html/1c2c8c8a-c1d3-7969-d6fa-afa2eb43193c.htm "Read 方法 ")

[ReadAsync 方法](../html/10de74b6-5396-f9c5-584a-7feef7dd0409.htm "ReadAsync 方法 ")

[ReadBool 方法](../html/9b6c8b8f-c3f1-be92-7de3-24a35f22d5e4.htm "ReadBool 方法 ")

[ReadBoolAsync 方法](../html/84e03db1-bbfe-55df-a5de-3fcc91e24336.htm "ReadBoolAsync 方法 ")

[ReadByte 方法](../html/c477231b-11fd-ddde-e293-5d5a3b7cd473.htm "ReadByte 方法 ")

[ReadByteAsync 方法](../html/e6e6e840-c1a3-b92e-25bf-609016104961.htm "ReadByteAsync 方法 ")

[ReadPLCTime 方法](../html/cd7a4ead-ed5c-4cab-9821-e98f2c6e80aa.htm "ReadPLCTime 方法 ")

[ReadPLCTimeAsync 方法](../html/194f0501-71c1-2117-5b63-696759d21a80.htm "ReadPLCTimeAsync 方法 ")

[ReadProgramName 方法](../html/394d566d-cf52-338e-8f3e-e8c63d28f4a4.htm "ReadProgramName 方法 ")

[ReadProgramNameAsync 方法](../html/0f5ad8db-9e52-7c29-7626-dc894dfec490.htm "ReadProgramNameAsync 方法 ")

[ToString 方法](../html/ef49b930-7974-80fa-c392-84830ae0ee83.htm "ToString 方法 ")

[Write 方法](../html/01b0b2e3-f297-36b6-0be3-9225afcce5f7.htm "Write 方法 ")

[WriteAsync 方法](../html/0d2be398-b4a1-001c-653d-b30b94b027a1.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| GeSRTPNet 方法 |

[GeSRTPNet](52c0166b-3420-af47-8dec-f25e74515116.htm) 类型公开以下成员。

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
| 受保护的方法 | [GetNewNetMessage](925f55de-2df6-75c6-df9f-753de861bebc.htm) | 获取一个新的消息对象的方法，需要在继承类里面进行重写  The method to get a new message object needs to be overridden in the inheritance class (重写 [BinaryCommunicationGetNewNetMessage](16636a1b-a6f0-9c0f-9d9c-c94e26c67459.htm).) |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法 | [GetWordLength](c27821ea-1739-8eb9-a848-23a00828fec8.htm) | 一个字单位的数据表示的地址长度，西门子为2，三菱，欧姆龙，modbusTcp就为1，AB PLC无效  The address length represented by one word of data, Siemens is 2, Mitsubishi, Omron, modbusTcp is 1, AB PLC is invalid (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 受保护的方法代码示例 | [InitializationOnConnect](9364af5f-5b71-458f-44b6-1b26ac09fac7.htm) | 根据实际的协议选择是否重写本方法，有些协议在创建连接之后，需要进行一些初始化的信号握手，才能最终建立网络通道。  Whether to rewrite this method is based on the actual protocol. Some protocols require some initial signal handshake to establish a network channel after the connection is created. (重写 [BinaryCommunicationInitializationOnConnect](4ab6989f-a5b5-05fe-8b44-bfa45824a2a7.htm).) |
| 受保护的方法代码示例 | [InitializationOnConnectAsync](23faa44d-2cc1-80c9-b6bc-995a2bdb8ab2.htm) | 根据实际的协议选择是否重写本方法，有些协议在创建连接之后，需要进行一些初始化的信号握手，才能最终建立网络通道。  Whether to rewrite this method is based on the actual protocol. Some protocols require some initial signal handshake to establish a network channel after the connection is created. (重写 [BinaryCommunicationInitializationOnConnectAsync](7da2af04-d4fc-5045-19a4-4b0d1722eef3.htm).) |
| 公共方法 | [IpAddressPing](0ac126b4-a830-a822-4ef5-855474942cd7.htm) | 对当前设备的IP地址进行PING的操作，返回PING的结果，正常来说，返回Success  PING the IP address of the current device and return the PING result. Normally, it returns Success (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 受保护的方法 | [LogRevcMessage(Byte)](dab1657f-1e75-7fe1-7003-ccedfc69855d.htm) | 使用日志记录一个接收的报文信息  Logs are used to record information about a received packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [LogRevcMessage(Byte, PipeSession)](0f51ae45-575e-4228-f978-01b37d6d5df1.htm) | 使用日志记录一个接收的报文信息  Logs are used to record information about a received packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [LogSendMessage(Byte)](4f90eecc-68a2-1ef2-44ef-369f19bed3ec.htm) | 使用日志记录一个发送的报文信息  Logs are used to record information about a send packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [LogSendMessage(Byte, PipeSession)](209caadf-4f0f-40aa-04c2-ebc61f8dab13.htm) | 使用日志记录一个发送的报文信息  Logs are used to record information about a send packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法 | [PackCommandWithHeader](d0b2cd76-c306-11fc-38c7-af8ad306cce7.htm) | 对当前的命令进行打包处理，通常是携带命令头内容，标记当前的命令的长度信息，需要进行重写，否则默认不打包  The current command is packaged, usually carrying the content of the command header, marking the length of the current command, and it needs to be rewritten, otherwise it is not packaged by default (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [Read(String, UInt16)](92534a5c-499e-278f-005e-57f05207e689.htm) | 批量读取字节数组信息，需要指定地址和长度，返回原始的字节数组，支持 I,Q,M,T,SA,SB,SC,S,G 的位和字节读写，支持 AI,AQ,R 的字读写操作，地址示例：R1,M1  Batch read byte array information, need to specify the address and length, return the original byte array. Support I, Q, M, T, SA, SB, SC, S, G bit and byte read and write, support AI, AQ, R word read and write operations, address examples: R1, M1 (重写 [DeviceCommunicationRead(String, UInt16)](9dd09452-c200-bc65-75ff-69aa34b366d4.htm).) |
| 公共方法代码示例 | [ReadT](95abe064-5594-5ae7-a849-98c53f88fee0.htm) | 读取支持Hsl特性的数据内容，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明，或是官网的详细文档  Read the data content of the Hsl attribute. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadAsync(String, UInt16)](176d000d-4655-22b6-5610-1348a997f8e5.htm) | 批量读取字节数组信息，需要指定地址和长度，返回原始的字节数组，支持 I,Q,M,T,SA,SB,SC,S,G 的位和字节读写，支持 AI,AQ,R 的字读写操作，地址示例：R1,M1  Batch read byte array information, need to specify the address and length, return the original byte array. Support I, Q, M, T, SA, SB, SC, S, G bit and byte read and write, support AI, AQ, R word read and write operations, address examples: R1, M1 (重写 [DeviceCommunicationReadAsync(String, UInt16)](d8c14d37-56dc-2a31-1a19-3e2b5b65197d.htm).) |
| 公共方法代码示例 | [ReadAsyncT](0522e740-4bf0-2e31-c4c5-e41384321b37.htm) | 异步读取支持Hsl特性的数据内容，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Asynchronously read the data content of the Hsl attribute. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadBool(String)](9b2faa2e-eaa9-7ace-c49f-7a0318f1df1a.htm) | 读取单个的Boolean数据信息  Read a single Boolean data message (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadBool(String, UInt16)](cfc5e2c8-62cc-a442-48cf-3acc4a2d8261.htm) | 按照位为单位，批量从指定的地址里读取 bool 数组数据，不支持 AI，AQ，R 地址类型，地址比如从1开始，例如 I1,Q1,M1,T1,SA1,SB1,SC1,S1,G1  In units of bits, read bool array data from the specified address in batches. AI, AQ, R address types are not supported. For example, the address starts from 1, such as I1, Q1, M1, T1, SA1, SB1, SC1, S1, G1 (重写 [DeviceCommunicationReadBool(String, UInt16)](69d0d24f-34a7-5b41-5c1b-619edc7e6b54.htm).) |
| 公共方法 | [ReadBoolAsync(String)](cc37da1f-d673-fe67-78ee-4e465d29a2bd.htm) | 异步读取单个的Boolean数据信息  Asynchronously read a single Boolean data message (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadBoolAsync(String, UInt16)](b079c820-bf03-3d5d-bff8-fd980f4c68f1.htm) | 按照位为单位，批量从指定的地址里读取 bool 数组数据，不支持 AI，AQ，R 地址类型，地址比如从1开始，例如 I1,Q1,M1,T1,SA1,SB1,SC1,S1,G1  In units of bits, read bool array data from the specified address in batches. AI, AQ, R address types are not supported. For example, the address starts from 1, such as I1, Q1, M1, T1, SA1, SB1, SC1, S1, G1 (重写 [DeviceCommunicationReadBoolAsync(String, UInt16)](a45d383d-0095-3ea7-32f8-e75af88d60f8.htm).) |
| 公共方法 | [ReadByte](c477231b-11fd-ddde-e293-5d5a3b7cd473.htm) | 根据指定的地址来读取一个字节的数据，按照字节为单位，例如 M1 字节，就是指 M1-M8 位组成的字节，M2 字节就是 M9-M16 组成的字节。不支持对 AI,AQ,R 寄存器的字节读取  A byte of data is read according to the specified address, and the unit is byte. For example, M1 byte refers to a byte composed of M1-M8 bits, and M2 byte is a byte composed of M9-M16. Does not support byte reading of AI, AQ, R registers |
| 公共方法 | [ReadByteAsync](e6e6e840-c1a3-b92e-25bf-609016104961.htm) | 根据指定的地址来读取一个字节的数据，按照字节为单位，例如 M1 字节，就是指 M1-M8 位组成的字节，M2 字节就是 M9-M16 组成的字节。不支持对 AI,AQ,R 寄存器的字节读取  A byte of data is read according to the specified address, and the unit is byte. For example, M1 byte refers to a byte composed of M1-M8 bits, and M2 byte is a byte composed of M9-M16. Does not support byte reading of AI, AQ, R registers |
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
| 公共方法 | [ReadPLCTime](cd7a4ead-ed5c-4cab-9821-e98f2c6e80aa.htm) | **[商业授权]** 读取PLC当前的时间，这个时间可能是不包含时区的，需要自己转换成本地的时间。 **[Authorization]** Read the current time of the PLC, this time may not include the time zone, you need to convert the local time yourself. |
| 公共方法 | [ReadPLCTimeAsync](194f0501-71c1-2117-5b63-696759d21a80.htm) | **[商业授权]** 读取PLC当前的时间，这个时间可能是不包含时区的，需要自己转换成本地的时间。 **[Authorization]** Read the current time of the PLC, this time may not include the time zone, you need to convert the local time yourself. |
| 公共方法 | [ReadProgramName](394d566d-cf52-338e-8f3e-e8c63d28f4a4.htm) | **[商业授权]** 读取PLC当前的程序的名称 **[Authorization]** Read the name of the current program of the PLC |
| 公共方法 | [ReadProgramNameAsync](0f5ad8db-9e52-7c29-7626-dc894dfec490.htm) | **[商业授权]** 读取PLC当前的程序的名称 **[Authorization]** Read the name of the current program of the PLC |
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
| 公共方法 | [SetDtuPipe](f765d1d7-234a-b229-9457-97e25e024dc9.htm) | 将当前的通信对象设置DTU模式，允许传入现成的管道，并返回初始化结果，如果该设备重写了握手报文，就是返回握手结果  Set the current communication object to DTU mode, allow the existing pipe to be passed in, and return the initialization result, if the device rewrites the handshake packet, the handshake result is returned (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [SetDtuPipeAsync](9d1fd12d-b804-3828-a627-a2c0eef1d23f.htm) | 将当前的通信对象设置DTU模式，允许传入现成的管道，并返回初始化结果，如果该设备重写了握手报文，就是返回握手结果  Set the current communication object to DTU mode, allow the existing pipe to be passed in, and return the initialization result, if the device rewrites the handshake packet, the handshake result is returned (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [SetPersistentConnection](f903dae0-4a4e-c7db-8350-50e36b257ec5.htm) | **已过时。** V11版本及之前设置长连接的方法，在V12版本以上中没有任何效果，默认长连接，删除调用即可，此处保留方法是为了部分用户保持兼容性升级。  The method of setting the long connection in V11 and before, has no effect in V12 and above. this method can be deleted. The method is retained here to maintain compatibility upgrades for some users. (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共方法 | [ToString](ef49b930-7974-80fa-c392-84830ae0ee83.htm) | (重写 [DeviceTcpNetToString](209a196b-90ea-2b73-e915-ce4b11f1263d.htm).) |
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
| 公共方法 | [Write(String, Boolean)](18c50ec3-9b8d-7cde-0429-5577a95562f5.htm) | 按照位为单位，批量写入 bool 数组到指定的地址里，不支持 AI，AQ，R 地址类型，地址比如从1开始，例如 I1,Q1,M1,T1,SA1,SB1,SC1,S1,G1  In units of bits, write bool arrays in batches to the specified addresses. AI, AQ, and R address types are not supported. For example, the address starts from 1, such as I1, Q1, M1, T1, SA1, SB1, SC1, S1, G1 (重写 [DeviceCommunicationWrite(String, Boolean)](5b046c82-3d8e-550b-fdf3-41ca5558169b.htm).) |
| 公共方法 | [Write(String, Byte)](9a0ec8e1-6d5c-3e0e-1e0a-3d3534043d62.htm) | 向PLC中写入byte数据，返回是否写入成功  Write byte data to PLC and return whether the writing is successful |
| 公共方法 | [Write(String, Byte)](f958948c-4214-7418-1525-b7bf0a236a9e.htm) | 写入原始的byte数组数据到指定的地址，返回是否写入成功  Write the original byte array data to the specified address, and return whether the write was successful (重写 [DeviceCommunicationWrite(String, Byte)](b5536dd6-ccde-6883-ac27-bff9e4b37806.htm).) |
| 公共方法代码示例 | [Write(String, String, Int32)](02fe2ca9-49a3-ede4-d1af-8c6da1f644da.htm) | 写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为ASCII  Write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, it is filled with 0 and the encoding is ASCII. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, String, Encoding)](7b602661-902d-2b70-1fdb-ba87ecd307bf.htm) | 写入字符串信息，需要指定的编码信息  Write string information, need to specify the encoding information (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, String, Int32, Encoding)](e8cb30d9-4c5e-6c68-b9af-1b1b6c2dc714.htm) | 写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为指定的编码信息  Write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, then the operation is complemented with 0 , you should specified the encoding information (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteT(T)](9f48853f-fbeb-c37d-25a2-1c76d1dee407.htm) | 写入支持Hsl特性的数据，返回是否写入成功，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Write data that supports the Hsl attribute, and return whether the write was successful. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WriteAsync(String, Boolean)](10fe669c-805e-1e5a-f016-f319ad363e90.htm) | 异步批量写入Boolean数组数据，返回是否成功  Asynchronously batch write Boolean array data, return success (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
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
| 公共方法 | [WriteAsync(String, Boolean)](c10ddba4-b701-b34c-2d58-cd94e308de96.htm) | 按照位为单位，批量写入 bool 数组到指定的地址里，不支持 AI，AQ，R 地址类型，地址比如从1开始，例如 I1,Q1,M1,T1,SA1,SB1,SC1,S1,G1  In units of bits, write bool arrays in batches to the specified addresses. AI, AQ, and R address types are not supported. For example, the address starts from 1, such as I1, Q1, M1, T1, SA1, SB1, SC1, S1, G1 (重写 [DeviceCommunicationWriteAsync(String, Boolean)](fe872dae-7ef9-12c1-d883-e92ba8676416.htm).) |
| 公共方法 | [WriteAsync(String, Byte)](604f4cfc-50c1-0119-1320-fe91e83a1850.htm) | 向PLC中写入byte数据，返回是否写入成功  Write byte data to PLC and return whether the writing is successful |
| 公共方法 | [WriteAsync(String, Byte)](98119c8f-0d6b-fa09-7403-10fedd11e528.htm) | 写入原始的byte数组数据到指定的地址，返回是否写入成功  Write the original byte array data to the specified address, and return whether the write was successful (重写 [DeviceCommunicationWriteAsync(String, Byte)](8f1cec30-8b86-6749-c8f3-7d39168d59df.htm).) |
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

[GeSRTPNet 类](52c0166b-3420-af47-8dec-f25e74515116.htm)

[HslCommunication.Profinet.GE 命名空间](db594ebf-c4ae-a11c-a82c-72625c3b6ab4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GetNewNetMessage 方法 

[原文連結](http://api.hslcommunication.cn/html/925f55de-2df6-75c6-df9f-753de861bebc.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.GE](../html/db594ebf-c4ae-a11c-a82c-72625c3b6ab4.htm "HslCommunication.Profinet.GE")

[GeSRTPNet 类](../html/52c0166b-3420-af47-8dec-f25e74515116.htm "GeSRTPNet 类")

[GeSRTPNet 方法](../html/a589d191-964f-8e25-3d96-9d59e2371b59.htm "GeSRTPNet 方法")

[GetNewNetMessage 方法](../html/925f55de-2df6-75c6-df9f-753de861bebc.htm "GetNewNetMessage 方法 ")

[InitializationOnConnect 方法](../html/9364af5f-5b71-458f-44b6-1b26ac09fac7.htm "InitializationOnConnect 方法 ")

[InitializationOnConnectAsync 方法](../html/23faa44d-2cc1-80c9-b6bc-995a2bdb8ab2.htm "InitializationOnConnectAsync 方法 ")

[Read 方法](../html/1c2c8c8a-c1d3-7969-d6fa-afa2eb43193c.htm "Read 方法 ")

[ReadAsync 方法](../html/10de74b6-5396-f9c5-584a-7feef7dd0409.htm "ReadAsync 方法 ")

[ReadBool 方法](../html/9b6c8b8f-c3f1-be92-7de3-24a35f22d5e4.htm "ReadBool 方法 ")

[ReadBoolAsync 方法](../html/84e03db1-bbfe-55df-a5de-3fcc91e24336.htm "ReadBoolAsync 方法 ")

[ReadByte 方法](../html/c477231b-11fd-ddde-e293-5d5a3b7cd473.htm "ReadByte 方法 ")

[ReadByteAsync 方法](../html/e6e6e840-c1a3-b92e-25bf-609016104961.htm "ReadByteAsync 方法 ")

[ReadPLCTime 方法](../html/cd7a4ead-ed5c-4cab-9821-e98f2c6e80aa.htm "ReadPLCTime 方法 ")

[ReadPLCTimeAsync 方法](../html/194f0501-71c1-2117-5b63-696759d21a80.htm "ReadPLCTimeAsync 方法 ")

[ReadProgramName 方法](../html/394d566d-cf52-338e-8f3e-e8c63d28f4a4.htm "ReadProgramName 方法 ")

[ReadProgramNameAsync 方法](../html/0f5ad8db-9e52-7c29-7626-dc894dfec490.htm "ReadProgramNameAsync 方法 ")

[ToString 方法](../html/ef49b930-7974-80fa-c392-84830ae0ee83.htm "ToString 方法 ")

[Write 方法](../html/01b0b2e3-f297-36b6-0be3-9225afcce5f7.htm "Write 方法 ")

[WriteAsync 方法](../html/0d2be398-b4a1-001c-653d-b30b94b027a1.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| GeSRTPNetGetNewNetMessage 方法 |

获取一个新的消息对象的方法，需要在继承类里面进行重写  
The method to get a new message object needs to be overridden in the inheritance class

**命名空间：**
 [HslCommunication.Profinet.GE](db594ebf-c4ae-a11c-a82c-72625c3b6ab4.htm)  
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

[GeSRTPNet 类](52c0166b-3420-af47-8dec-f25e74515116.htm)

[HslCommunication.Profinet.GE 命名空间](db594ebf-c4ae-a11c-a82c-72625c3b6ab4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## InitializationOnConnect 方法 

[原文連結](http://api.hslcommunication.cn/html/9364af5f-5b71-458f-44b6-1b26ac09fac7.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.GE](../html/db594ebf-c4ae-a11c-a82c-72625c3b6ab4.htm "HslCommunication.Profinet.GE")

[GeSRTPNet 类](../html/52c0166b-3420-af47-8dec-f25e74515116.htm "GeSRTPNet 类")

[GeSRTPNet 方法](../html/a589d191-964f-8e25-3d96-9d59e2371b59.htm "GeSRTPNet 方法")

[GetNewNetMessage 方法](../html/925f55de-2df6-75c6-df9f-753de861bebc.htm "GetNewNetMessage 方法 ")

[InitializationOnConnect 方法](../html/9364af5f-5b71-458f-44b6-1b26ac09fac7.htm "InitializationOnConnect 方法 ")

[InitializationOnConnectAsync 方法](../html/23faa44d-2cc1-80c9-b6bc-995a2bdb8ab2.htm "InitializationOnConnectAsync 方法 ")

[Read 方法](../html/1c2c8c8a-c1d3-7969-d6fa-afa2eb43193c.htm "Read 方法 ")

[ReadAsync 方法](../html/10de74b6-5396-f9c5-584a-7feef7dd0409.htm "ReadAsync 方法 ")

[ReadBool 方法](../html/9b6c8b8f-c3f1-be92-7de3-24a35f22d5e4.htm "ReadBool 方法 ")

[ReadBoolAsync 方法](../html/84e03db1-bbfe-55df-a5de-3fcc91e24336.htm "ReadBoolAsync 方法 ")

[ReadByte 方法](../html/c477231b-11fd-ddde-e293-5d5a3b7cd473.htm "ReadByte 方法 ")

[ReadByteAsync 方法](../html/e6e6e840-c1a3-b92e-25bf-609016104961.htm "ReadByteAsync 方法 ")

[ReadPLCTime 方法](../html/cd7a4ead-ed5c-4cab-9821-e98f2c6e80aa.htm "ReadPLCTime 方法 ")

[ReadPLCTimeAsync 方法](../html/194f0501-71c1-2117-5b63-696759d21a80.htm "ReadPLCTimeAsync 方法 ")

[ReadProgramName 方法](../html/394d566d-cf52-338e-8f3e-e8c63d28f4a4.htm "ReadProgramName 方法 ")

[ReadProgramNameAsync 方法](../html/0f5ad8db-9e52-7c29-7626-dc894dfec490.htm "ReadProgramNameAsync 方法 ")

[ToString 方法](../html/ef49b930-7974-80fa-c392-84830ae0ee83.htm "ToString 方法 ")

[Write 方法](../html/01b0b2e3-f297-36b6-0be3-9225afcce5f7.htm "Write 方法 ")

[WriteAsync 方法](../html/0d2be398-b4a1-001c-653d-b30b94b027a1.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| GeSRTPNetInitializationOnConnect 方法 |

根据实际的协议选择是否重写本方法，有些协议在创建连接之后，需要进行一些初始化的信号握手，才能最终建立网络通道。  
Whether to rewrite this method is based on the actual protocol. Some protocols require some initial signal handshake to establish a network channel after the connection is created.

**命名空间：**
 [HslCommunication.Profinet.GE](db594ebf-c4ae-a11c-a82c-72625c3b6ab4.htm)  
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

[GeSRTPNet 类](52c0166b-3420-af47-8dec-f25e74515116.htm)

[HslCommunication.Profinet.GE 命名空间](db594ebf-c4ae-a11c-a82c-72625c3b6ab4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## InitializationOnConnectAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/23faa44d-2cc1-80c9-b6bc-995a2bdb8ab2.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.GE](../html/db594ebf-c4ae-a11c-a82c-72625c3b6ab4.htm "HslCommunication.Profinet.GE")

[GeSRTPNet 类](../html/52c0166b-3420-af47-8dec-f25e74515116.htm "GeSRTPNet 类")

[GeSRTPNet 方法](../html/a589d191-964f-8e25-3d96-9d59e2371b59.htm "GeSRTPNet 方法")

[GetNewNetMessage 方法](../html/925f55de-2df6-75c6-df9f-753de861bebc.htm "GetNewNetMessage 方法 ")

[InitializationOnConnect 方法](../html/9364af5f-5b71-458f-44b6-1b26ac09fac7.htm "InitializationOnConnect 方法 ")

[InitializationOnConnectAsync 方法](../html/23faa44d-2cc1-80c9-b6bc-995a2bdb8ab2.htm "InitializationOnConnectAsync 方法 ")

[Read 方法](../html/1c2c8c8a-c1d3-7969-d6fa-afa2eb43193c.htm "Read 方法 ")

[ReadAsync 方法](../html/10de74b6-5396-f9c5-584a-7feef7dd0409.htm "ReadAsync 方法 ")

[ReadBool 方法](../html/9b6c8b8f-c3f1-be92-7de3-24a35f22d5e4.htm "ReadBool 方法 ")

[ReadBoolAsync 方法](../html/84e03db1-bbfe-55df-a5de-3fcc91e24336.htm "ReadBoolAsync 方法 ")

[ReadByte 方法](../html/c477231b-11fd-ddde-e293-5d5a3b7cd473.htm "ReadByte 方法 ")

[ReadByteAsync 方法](../html/e6e6e840-c1a3-b92e-25bf-609016104961.htm "ReadByteAsync 方法 ")

[ReadPLCTime 方法](../html/cd7a4ead-ed5c-4cab-9821-e98f2c6e80aa.htm "ReadPLCTime 方法 ")

[ReadPLCTimeAsync 方法](../html/194f0501-71c1-2117-5b63-696759d21a80.htm "ReadPLCTimeAsync 方法 ")

[ReadProgramName 方法](../html/394d566d-cf52-338e-8f3e-e8c63d28f4a4.htm "ReadProgramName 方法 ")

[ReadProgramNameAsync 方法](../html/0f5ad8db-9e52-7c29-7626-dc894dfec490.htm "ReadProgramNameAsync 方法 ")

[ToString 方法](../html/ef49b930-7974-80fa-c392-84830ae0ee83.htm "ToString 方法 ")

[Write 方法](../html/01b0b2e3-f297-36b6-0be3-9225afcce5f7.htm "Write 方法 ")

[WriteAsync 方法](../html/0d2be398-b4a1-001c-653d-b30b94b027a1.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| GeSRTPNetInitializationOnConnectAsync 方法 |

根据实际的协议选择是否重写本方法，有些协议在创建连接之后，需要进行一些初始化的信号握手，才能最终建立网络通道。  
Whether to rewrite this method is based on the actual protocol. Some protocols require some initial signal handshake to establish a network channel after the connection is created.

**命名空间：**
 [HslCommunication.Profinet.GE](db594ebf-c4ae-a11c-a82c-72625c3b6ab4.htm)  
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

[GeSRTPNet 类](52c0166b-3420-af47-8dec-f25e74515116.htm)

[HslCommunication.Profinet.GE 命名空间](db594ebf-c4ae-a11c-a82c-72625c3b6ab4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Read 方法 

[原文連結](http://api.hslcommunication.cn/html/1c2c8c8a-c1d3-7969-d6fa-afa2eb43193c.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.GE](../html/db594ebf-c4ae-a11c-a82c-72625c3b6ab4.htm "HslCommunication.Profinet.GE")

[GeSRTPNet 类](../html/52c0166b-3420-af47-8dec-f25e74515116.htm "GeSRTPNet 类")

[GeSRTPNet 方法](../html/a589d191-964f-8e25-3d96-9d59e2371b59.htm "GeSRTPNet 方法")

[Read 方法](../html/1c2c8c8a-c1d3-7969-d6fa-afa2eb43193c.htm "Read 方法 ")

[Read 方法 (String, UInt16)](../html/92534a5c-499e-278f-005e-57f05207e689.htm "Read 方法 (String, UInt16)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| GeSRTPNetRead 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法代码示例 | [ReadT](95abe064-5594-5ae7-a849-98c53f88fee0.htm) | 读取支持Hsl特性的数据内容，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明，或是官网的详细文档  Read the data content of the Hsl attribute. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [Read(String, UInt16)](92534a5c-499e-278f-005e-57f05207e689.htm) | 批量读取字节数组信息，需要指定地址和长度，返回原始的字节数组，支持 I,Q,M,T,SA,SB,SC,S,G 的位和字节读写，支持 AI,AQ,R 的字读写操作，地址示例：R1,M1  Batch read byte array information, need to specify the address and length, return the original byte array. Support I, Q, M, T, SA, SB, SC, S, G bit and byte read and write, support AI, AQ, R word read and write operations, address examples: R1, M1 (重写 [DeviceCommunicationRead(String, UInt16)](9dd09452-c200-bc65-75ff-69aa34b366d4.htm).) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[GeSRTPNet 类](52c0166b-3420-af47-8dec-f25e74515116.htm)

[HslCommunication.Profinet.GE 命名空间](db594ebf-c4ae-a11c-a82c-72625c3b6ab4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Read 方法 (String, UInt16)

[原文連結](http://api.hslcommunication.cn/html/92534a5c-499e-278f-005e-57f05207e689.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.GE](../html/db594ebf-c4ae-a11c-a82c-72625c3b6ab4.htm "HslCommunication.Profinet.GE")

[GeSRTPNet 类](../html/52c0166b-3420-af47-8dec-f25e74515116.htm "GeSRTPNet 类")

[GeSRTPNet 方法](../html/a589d191-964f-8e25-3d96-9d59e2371b59.htm "GeSRTPNet 方法")

[Read 方法](../html/1c2c8c8a-c1d3-7969-d6fa-afa2eb43193c.htm "Read 方法 ")

[Read 方法 (String, UInt16)](../html/92534a5c-499e-278f-005e-57f05207e689.htm "Read 方法 (String, UInt16)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| GeSRTPNetRead 方法 (String, UInt16) |

批量读取字节数组信息，需要指定地址和长度，返回原始的字节数组，支持 I,Q,M,T,SA,SB,SC,S,G 的位和字节读写，支持 AI,AQ,R 的字读写操作，地址示例：R1,M1  
Batch read byte array information, need to specify the address and length, return the original byte array.
Support I, Q, M, T, SA, SB, SC, S, G bit and byte read and write, support AI, AQ, R word read and write operations, address examples: R1, M1

**命名空间：**
 [HslCommunication.Profinet.GE](db594ebf-c4ae-a11c-a82c-72625c3b6ab4.htm)  
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
    数据地址

length
:   类型：SystemUInt16  
    数据长度

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
带有成功标识的byte[]数组

#### 实现

[IReadWriteNetRead(String, UInt16)](b8d4187f-ab9a-79e9-c7ca-8a01ec0df233.htm)

![](../icons/SectionExpanded.png)备注

其中读取R，AI，AQ寄存器的原始字节时，传入的长度参数为字节长度。长度为10，返回10个字节数组信息，如果返回长度不满6个字节的，一律返回6个字节的数据

![](../icons/SectionExpanded.png)参见

#### 引用

[GeSRTPNet 类](52c0166b-3420-af47-8dec-f25e74515116.htm)

[Read 重载](1c2c8c8a-c1d3-7969-d6fa-afa2eb43193c.htm)

[HslCommunication.Profinet.GE 命名空间](db594ebf-c4ae-a11c-a82c-72625c3b6ab4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/10de74b6-5396-f9c5-584a-7feef7dd0409.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.GE](../html/db594ebf-c4ae-a11c-a82c-72625c3b6ab4.htm "HslCommunication.Profinet.GE")

[GeSRTPNet 类](../html/52c0166b-3420-af47-8dec-f25e74515116.htm "GeSRTPNet 类")

[GeSRTPNet 方法](../html/a589d191-964f-8e25-3d96-9d59e2371b59.htm "GeSRTPNet 方法")

[ReadAsync 方法](../html/10de74b6-5396-f9c5-584a-7feef7dd0409.htm "ReadAsync 方法 ")

[ReadAsync 方法 (String, UInt16)](../html/176d000d-4655-22b6-5610-1348a997f8e5.htm "ReadAsync 方法 (String, UInt16)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| GeSRTPNetReadAsync 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法代码示例 | [ReadAsyncT](0522e740-4bf0-2e31-c4c5-e41384321b37.htm) | 异步读取支持Hsl特性的数据内容，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Asynchronously read the data content of the Hsl attribute. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadAsync(String, UInt16)](176d000d-4655-22b6-5610-1348a997f8e5.htm) | 批量读取字节数组信息，需要指定地址和长度，返回原始的字节数组，支持 I,Q,M,T,SA,SB,SC,S,G 的位和字节读写，支持 AI,AQ,R 的字读写操作，地址示例：R1,M1  Batch read byte array information, need to specify the address and length, return the original byte array. Support I, Q, M, T, SA, SB, SC, S, G bit and byte read and write, support AI, AQ, R word read and write operations, address examples: R1, M1 (重写 [DeviceCommunicationReadAsync(String, UInt16)](d8c14d37-56dc-2a31-1a19-3e2b5b65197d.htm).) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[GeSRTPNet 类](52c0166b-3420-af47-8dec-f25e74515116.htm)

[HslCommunication.Profinet.GE 命名空间](db594ebf-c4ae-a11c-a82c-72625c3b6ab4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadAsync 方法 (String, UInt16)

[原文連結](http://api.hslcommunication.cn/html/176d000d-4655-22b6-5610-1348a997f8e5.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.GE](../html/db594ebf-c4ae-a11c-a82c-72625c3b6ab4.htm "HslCommunication.Profinet.GE")

[GeSRTPNet 类](../html/52c0166b-3420-af47-8dec-f25e74515116.htm "GeSRTPNet 类")

[GeSRTPNet 方法](../html/a589d191-964f-8e25-3d96-9d59e2371b59.htm "GeSRTPNet 方法")

[ReadAsync 方法](../html/10de74b6-5396-f9c5-584a-7feef7dd0409.htm "ReadAsync 方法 ")

[ReadAsync 方法 (String, UInt16)](../html/176d000d-4655-22b6-5610-1348a997f8e5.htm "ReadAsync 方法 (String, UInt16)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| GeSRTPNetReadAsync 方法 (String, UInt16) |

批量读取字节数组信息，需要指定地址和长度，返回原始的字节数组，支持 I,Q,M,T,SA,SB,SC,S,G 的位和字节读写，支持 AI,AQ,R 的字读写操作，地址示例：R1,M1  
Batch read byte array information, need to specify the address and length, return the original byte array.
Support I, Q, M, T, SA, SB, SC, S, G bit and byte read and write, support AI, AQ, R word read and write operations, address examples: R1, M1

**命名空间：**
 [HslCommunication.Profinet.GE](db594ebf-c4ae-a11c-a82c-72625c3b6ab4.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public override Task<OperateResult<byte[]>> ReadAsync(
	string address,
	ushort length
)
```

```
Public Overrides Function ReadAsync ( 
	address As String,
	length As UShort
) As Task(Of OperateResult(Of Byte()))
```

```
public:
virtual Task<OperateResult<array<unsigned char>^>^>^ ReadAsync(
	String^ address, 
	unsigned short length
) override
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

![](../icons/SectionExpanded.png)备注

其中读取R，AI，AQ寄存器的原始字节时，传入的长度参数为字节长度。长度为10，返回10个字节数组信息，如果返回长度不满6个字节的，一律返回6个字节的数据

![](../icons/SectionExpanded.png)参见

#### 引用

[GeSRTPNet 类](52c0166b-3420-af47-8dec-f25e74515116.htm)

[ReadAsync 重载](10de74b6-5396-f9c5-584a-7feef7dd0409.htm)

[HslCommunication.Profinet.GE 命名空间](db594ebf-c4ae-a11c-a82c-72625c3b6ab4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadBool 方法 

[原文連結](http://api.hslcommunication.cn/html/9b6c8b8f-c3f1-be92-7de3-24a35f22d5e4.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.GE](../html/db594ebf-c4ae-a11c-a82c-72625c3b6ab4.htm "HslCommunication.Profinet.GE")

[GeSRTPNet 类](../html/52c0166b-3420-af47-8dec-f25e74515116.htm "GeSRTPNet 类")

[GeSRTPNet 方法](../html/a589d191-964f-8e25-3d96-9d59e2371b59.htm "GeSRTPNet 方法")

[ReadBool 方法](../html/9b6c8b8f-c3f1-be92-7de3-24a35f22d5e4.htm "ReadBool 方法 ")

[ReadBool 方法 (String, UInt16)](../html/cfc5e2c8-62cc-a442-48cf-3acc4a2d8261.htm "ReadBool 方法 (String, UInt16)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| GeSRTPNetReadBool 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [ReadBool(String)](9b2faa2e-eaa9-7ace-c49f-7a0318f1df1a.htm) | 读取单个的Boolean数据信息  Read a single Boolean data message (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadBool(String, UInt16)](cfc5e2c8-62cc-a442-48cf-3acc4a2d8261.htm) | 按照位为单位，批量从指定的地址里读取 bool 数组数据，不支持 AI，AQ，R 地址类型，地址比如从1开始，例如 I1,Q1,M1,T1,SA1,SB1,SC1,S1,G1  In units of bits, read bool array data from the specified address in batches. AI, AQ, R address types are not supported. For example, the address starts from 1, such as I1, Q1, M1, T1, SA1, SB1, SC1, S1, G1 (重写 [DeviceCommunicationReadBool(String, UInt16)](69d0d24f-34a7-5b41-5c1b-619edc7e6b54.htm).) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[GeSRTPNet 类](52c0166b-3420-af47-8dec-f25e74515116.htm)

[HslCommunication.Profinet.GE 命名空间](db594ebf-c4ae-a11c-a82c-72625c3b6ab4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadBool 方法 (String, UInt16)

[原文連結](http://api.hslcommunication.cn/html/cfc5e2c8-62cc-a442-48cf-3acc4a2d8261.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.GE](../html/db594ebf-c4ae-a11c-a82c-72625c3b6ab4.htm "HslCommunication.Profinet.GE")

[GeSRTPNet 类](../html/52c0166b-3420-af47-8dec-f25e74515116.htm "GeSRTPNet 类")

[GeSRTPNet 方法](../html/a589d191-964f-8e25-3d96-9d59e2371b59.htm "GeSRTPNet 方法")

[ReadBool 方法](../html/9b6c8b8f-c3f1-be92-7de3-24a35f22d5e4.htm "ReadBool 方法 ")

[ReadBool 方法 (String, UInt16)](../html/cfc5e2c8-62cc-a442-48cf-3acc4a2d8261.htm "ReadBool 方法 (String, UInt16)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| GeSRTPNetReadBool 方法 (String, UInt16) |

按照位为单位，批量从指定的地址里读取 bool 数组数据，不支持 AI，AQ，R 地址类型，地址比如从1开始，例如 I1,Q1,M1,T1,SA1,SB1,SC1,S1,G1  
In units of bits, read bool array data from the specified address in batches. AI, AQ,
R address types are not supported. For example, the address starts from 1, such as I1, Q1, M1, T1, SA1, SB1, SC1, S1, G1

**命名空间：**
 [HslCommunication.Profinet.GE](db594ebf-c4ae-a11c-a82c-72625c3b6ab4.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public override OperateResult<bool[]> ReadBool(
	string address,
	ushort length
)
```

```
Public Overrides Function ReadBool ( 
	address As String,
	length As UShort
) As OperateResult(Of Boolean())
```

```
public:
virtual OperateResult<array<bool>^>^ ReadBool(
	String^ address, 
	unsigned short length
) override
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
    PLC的地址信息，例如 M1, G1

length
:   类型：SystemUInt16  
    读取的长度信息，按照位为单位

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Boolean  
包含是否读取成功的bool数组

#### 实现

[IReadWriteNetReadBool(String, UInt16)](16e00e01-28d8-bfb7-e3bb-9d8b8778d7b9.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[GeSRTPNet 类](52c0166b-3420-af47-8dec-f25e74515116.htm)

[ReadBool 重载](9b6c8b8f-c3f1-be92-7de3-24a35f22d5e4.htm)

[HslCommunication.Profinet.GE 命名空间](db594ebf-c4ae-a11c-a82c-72625c3b6ab4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)