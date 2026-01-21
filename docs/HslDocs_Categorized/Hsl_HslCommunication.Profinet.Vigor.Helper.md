# HslCommunication - HslCommunication.Profinet.Vigor.Helper

> 分類頁數: 21



---
## HslCommunication.Profinet.Vigor.Helper

[原文連結](http://api.hslcommunication.cn/html/e525e2ab-bd42-dfb9-2389-84669c014cbb.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Vigor.Helper](../html/e525e2ab-bd42-dfb9-2389-84669c014cbb.htm "HslCommunication.Profinet.Vigor.Helper")

[VigorHelper 类](../html/204e8b8a-20c9-5818-14e3-c361aa47fbed.htm "VigorHelper 类")

[VigorVsHelper 类](../html/383f1b2a-980f-adaa-16af-3ccaca87cb84.htm "VigorVsHelper 类")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslCommunication.Profinet.Vigor.Helper 命名空间 |

[缺少 "N:HslCommunication.Profinet.Vigor.Helper" 的 <summary> 文档]

![](../icons/SectionExpanded.png)类

|  | 类 | 说明 |
| --- | --- | --- |
| 公共类 | [VigorHelper](204e8b8a-20c9-5818-14e3-c361aa47fbed.htm) | 丰炜PLC的辅助方法 |
| 公共类 | [VigorVsHelper](383f1b2a-980f-adaa-16af-3ccaca87cb84.htm) | 丰炜PLC的辅助类对象 |

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## VigorHelper 类

[原文連結](http://api.hslcommunication.cn/html/204e8b8a-20c9-5818-14e3-c361aa47fbed.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Vigor.Helper](../html/e525e2ab-bd42-dfb9-2389-84669c014cbb.htm "HslCommunication.Profinet.Vigor.Helper")

[VigorHelper 类](../html/204e8b8a-20c9-5818-14e3-c361aa47fbed.htm "VigorHelper 类")

[VigorHelper 构造函数](../html/2f7012bb-f477-83c1-bdc8-37a272284079.htm "VigorHelper 构造函数 ")

[VigorHelper 方法](../html/09b5496d-304b-b53c-2968-c981725c968c.htm "VigorHelper 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| VigorHelper 类 |

丰炜PLC的辅助方法

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  HslCommunication.Profinet.Vigor.HelperVigorHelper

**命名空间：**
 [HslCommunication.Profinet.Vigor.Helper](e525e2ab-bd42-dfb9-2389-84669c014cbb.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class VigorHelper
```

```
Public Class VigorHelper
```

```
public ref class VigorHelper
```

```
type VigorHelper =  class end
```

VigorHelper 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [VigorHelper](2f7012bb-f477-83c1-bdc8-37a272284079.htm) | 初始化 VigorHelper 类的一个新实例 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法静态成员 | [Read](82c88c10-df59-2fe5-9091-0b208fe5bdbf.htm) | 批量读取字节数组信息，需要指定地址和长度，返回原始的字节数组  Batch read byte array information, need to specify the address and length, return the original byte array |
| 公共方法静态成员 | [ReadAsync](12d8903c-f936-e4a9-765a-203a6a1a7d65.htm) | 批量读取字节数组信息，需要指定地址和长度，返回原始的字节数组  Batch read byte array information, need to specify the address and length, return the original byte array |
| 公共方法静态成员 | [ReadBool](59c6ed00-7b11-a3b5-f709-94f354c31569.htm) | 批量读取Boolean数组信息，需要指定地址和长度，返回Boolean 数组  Batch read Boolean array information, need to specify the address and length, return Boolean array |
| 公共方法静态成员 | [ReadBoolAsync](74b8830b-2a50-fff9-53f9-802d4ac13ad1.htm) | 批量读取Boolean数组信息，需要指定地址和长度，返回Boolean 数组  Batch read Boolean array information, need to specify the address and length, return Boolean array |
| 公共方法 | ToString | (继承自 Object。) |
| 公共方法静态成员 | [Write(IReadWriteDevice, Byte, String, Boolean)](033e8f76-6479-2b13-9892-a4c67509b982.htm) | 批量写入Boolean数组数据，返回是否成功  Batch write Boolean array data, return whether the write was successful |
| 公共方法静态成员 | [Write(IReadWriteDevice, Byte, String, Byte)](71aeee72-e80e-2f0f-b069-841341d39f79.htm) | 写入原始的byte数组数据到指定的地址，返回是否写入成功  Write the original byte array data to the specified address, and return whether the write was successful |
| 公共方法静态成员 | [WriteAsync(IReadWriteDevice, Byte, String, Boolean)](1009fb21-0e63-755c-9810-694fade4f4f8.htm) | 批量写入Boolean数组数据，返回是否成功  Batch write Boolean array data, return whether the write was successful |
| 公共方法静态成员 | [WriteAsync(IReadWriteDevice, Byte, String, Byte)](e64fe734-a64f-8bb1-d6f0-4284aa42312b.htm) | 写入原始的byte数组数据到指定的地址，返回是否写入成功  Write the original byte array data to the specified address, and return whether the write was successful |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.Profinet.Vigor.Helper 命名空间](e525e2ab-bd42-dfb9-2389-84669c014cbb.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## VigorHelper 构造函数 

[原文連結](http://api.hslcommunication.cn/html/2f7012bb-f477-83c1-bdc8-37a272284079.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Vigor.Helper](../html/e525e2ab-bd42-dfb9-2389-84669c014cbb.htm "HslCommunication.Profinet.Vigor.Helper")

[VigorHelper 类](../html/204e8b8a-20c9-5818-14e3-c361aa47fbed.htm "VigorHelper 类")

[VigorHelper 构造函数](../html/2f7012bb-f477-83c1-bdc8-37a272284079.htm "VigorHelper 构造函数 ")

[VigorHelper 方法](../html/09b5496d-304b-b53c-2968-c981725c968c.htm "VigorHelper 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| VigorHelper 构造函数 |

初始化 [VigorHelper](204e8b8a-20c9-5818-14e3-c361aa47fbed.htm) 类的一个新实例

**命名空间：**
 [HslCommunication.Profinet.Vigor.Helper](e525e2ab-bd42-dfb9-2389-84669c014cbb.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public VigorHelper()
```

```
Public Sub New
```

```
public:
VigorHelper()
```

```
new : unit -> VigorHelper
```

![](../icons/SectionExpanded.png)参见

#### 引用

[VigorHelper 类](204e8b8a-20c9-5818-14e3-c361aa47fbed.htm)

[HslCommunication.Profinet.Vigor.Helper 命名空间](e525e2ab-bd42-dfb9-2389-84669c014cbb.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## VigorHelper 方法

[原文連結](http://api.hslcommunication.cn/html/09b5496d-304b-b53c-2968-c981725c968c.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Vigor.Helper](../html/e525e2ab-bd42-dfb9-2389-84669c014cbb.htm "HslCommunication.Profinet.Vigor.Helper")

[VigorHelper 类](../html/204e8b8a-20c9-5818-14e3-c361aa47fbed.htm "VigorHelper 类")

[VigorHelper 方法](../html/09b5496d-304b-b53c-2968-c981725c968c.htm "VigorHelper 方法")

[Read 方法](../html/82c88c10-df59-2fe5-9091-0b208fe5bdbf.htm "Read 方法 ")

[ReadAsync 方法](../html/12d8903c-f936-e4a9-765a-203a6a1a7d65.htm "ReadAsync 方法 ")

[ReadBool 方法](../html/59c6ed00-7b11-a3b5-f709-94f354c31569.htm "ReadBool 方法 ")

[ReadBoolAsync 方法](../html/74b8830b-2a50-fff9-53f9-802d4ac13ad1.htm "ReadBoolAsync 方法 ")

[Write 方法](../html/79d68534-f5bc-d4fd-5b73-04626d0001c1.htm "Write 方法 ")

[WriteAsync 方法](../html/5451f32d-1f77-73c9-38a9-41845add2dd3.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| VigorHelper 方法 |

[VigorHelper](204e8b8a-20c9-5818-14e3-c361aa47fbed.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法静态成员 | [Read](82c88c10-df59-2fe5-9091-0b208fe5bdbf.htm) | 批量读取字节数组信息，需要指定地址和长度，返回原始的字节数组  Batch read byte array information, need to specify the address and length, return the original byte array |
| 公共方法静态成员 | [ReadAsync](12d8903c-f936-e4a9-765a-203a6a1a7d65.htm) | 批量读取字节数组信息，需要指定地址和长度，返回原始的字节数组  Batch read byte array information, need to specify the address and length, return the original byte array |
| 公共方法静态成员 | [ReadBool](59c6ed00-7b11-a3b5-f709-94f354c31569.htm) | 批量读取Boolean数组信息，需要指定地址和长度，返回Boolean 数组  Batch read Boolean array information, need to specify the address and length, return Boolean array |
| 公共方法静态成员 | [ReadBoolAsync](74b8830b-2a50-fff9-53f9-802d4ac13ad1.htm) | 批量读取Boolean数组信息，需要指定地址和长度，返回Boolean 数组  Batch read Boolean array information, need to specify the address and length, return Boolean array |
| 公共方法 | ToString | (继承自 Object。) |
| 公共方法静态成员 | [Write(IReadWriteDevice, Byte, String, Boolean)](033e8f76-6479-2b13-9892-a4c67509b982.htm) | 批量写入Boolean数组数据，返回是否成功  Batch write Boolean array data, return whether the write was successful |
| 公共方法静态成员 | [Write(IReadWriteDevice, Byte, String, Byte)](71aeee72-e80e-2f0f-b069-841341d39f79.htm) | 写入原始的byte数组数据到指定的地址，返回是否写入成功  Write the original byte array data to the specified address, and return whether the write was successful |
| 公共方法静态成员 | [WriteAsync(IReadWriteDevice, Byte, String, Boolean)](1009fb21-0e63-755c-9810-694fade4f4f8.htm) | 批量写入Boolean数组数据，返回是否成功  Batch write Boolean array data, return whether the write was successful |
| 公共方法静态成员 | [WriteAsync(IReadWriteDevice, Byte, String, Byte)](e64fe734-a64f-8bb1-d6f0-4284aa42312b.htm) | 写入原始的byte数组数据到指定的地址，返回是否写入成功  Write the original byte array data to the specified address, and return whether the write was successful |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[VigorHelper 类](204e8b8a-20c9-5818-14e3-c361aa47fbed.htm)

[HslCommunication.Profinet.Vigor.Helper 命名空间](e525e2ab-bd42-dfb9-2389-84669c014cbb.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Read 方法 

[原文連結](http://api.hslcommunication.cn/html/82c88c10-df59-2fe5-9091-0b208fe5bdbf.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Vigor.Helper](../html/e525e2ab-bd42-dfb9-2389-84669c014cbb.htm "HslCommunication.Profinet.Vigor.Helper")

[VigorHelper 类](../html/204e8b8a-20c9-5818-14e3-c361aa47fbed.htm "VigorHelper 类")

[VigorHelper 方法](../html/09b5496d-304b-b53c-2968-c981725c968c.htm "VigorHelper 方法")

[Read 方法](../html/82c88c10-df59-2fe5-9091-0b208fe5bdbf.htm "Read 方法 ")

[ReadAsync 方法](../html/12d8903c-f936-e4a9-765a-203a6a1a7d65.htm "ReadAsync 方法 ")

[ReadBool 方法](../html/59c6ed00-7b11-a3b5-f709-94f354c31569.htm "ReadBool 方法 ")

[ReadBoolAsync 方法](../html/74b8830b-2a50-fff9-53f9-802d4ac13ad1.htm "ReadBoolAsync 方法 ")

[Write 方法](../html/79d68534-f5bc-d4fd-5b73-04626d0001c1.htm "Write 方法 ")

[WriteAsync 方法](../html/5451f32d-1f77-73c9-38a9-41845add2dd3.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| VigorHelperRead 方法 |

批量读取字节数组信息，需要指定地址和长度，返回原始的字节数组  
Batch read byte array information, need to specify the address and length, return the original byte array

**命名空间：**
 [HslCommunication.Profinet.Vigor.Helper](e525e2ab-bd42-dfb9-2389-84669c014cbb.htm)  
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
	byte station,
	string address,
	ushort length
)
```

```
Public Shared Function Read ( 
	plc As IReadWriteDevice,
	station As Byte,
	address As String,
	length As UShort
) As OperateResult(Of Byte())
```

```
public:
static OperateResult<array<unsigned char>^>^ Read(
	IReadWriteDevice^ plc, 
	unsigned char station, 
	String^ address, 
	unsigned short length
)
```

```
static member Read : 
        plc : IReadWriteDevice * 
        station : byte * 
        address : string * 
        length : uint16 -> OperateResult<byte[]> 
```

#### 参数

plc
:   类型：[HslCommunication.CoreIReadWriteDevice](d987afcc-f99c-6aed-3847-5b1890504635.htm)  

    [缺少 "M:HslCommunication.Profinet.Vigor.Helper.VigorHelper.Read(HslCommunication.Core.IReadWriteDevice,System.Byte,System.String,System.UInt16)" 的 <param name="plc"/> 文档]

station
:   类型：SystemByte  

    [缺少 "M:HslCommunication.Profinet.Vigor.Helper.VigorHelper.Read(HslCommunication.Core.IReadWriteDevice,System.Byte,System.String,System.UInt16)" 的 <param name="station"/> 文档]

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

支持字地址，单次最多读取64字节，支持D,SD,R,T,C的数据读取，同时地址支持携带站号信息，s=2;D100

![](../icons/SectionExpanded.png)参见

#### 引用

[VigorHelper 类](204e8b8a-20c9-5818-14e3-c361aa47fbed.htm)

[HslCommunication.Profinet.Vigor.Helper 命名空间](e525e2ab-bd42-dfb9-2389-84669c014cbb.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/12d8903c-f936-e4a9-765a-203a6a1a7d65.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Vigor.Helper](../html/e525e2ab-bd42-dfb9-2389-84669c014cbb.htm "HslCommunication.Profinet.Vigor.Helper")

[VigorHelper 类](../html/204e8b8a-20c9-5818-14e3-c361aa47fbed.htm "VigorHelper 类")

[VigorHelper 方法](../html/09b5496d-304b-b53c-2968-c981725c968c.htm "VigorHelper 方法")

[Read 方法](../html/82c88c10-df59-2fe5-9091-0b208fe5bdbf.htm "Read 方法 ")

[ReadAsync 方法](../html/12d8903c-f936-e4a9-765a-203a6a1a7d65.htm "ReadAsync 方法 ")

[ReadBool 方法](../html/59c6ed00-7b11-a3b5-f709-94f354c31569.htm "ReadBool 方法 ")

[ReadBoolAsync 方法](../html/74b8830b-2a50-fff9-53f9-802d4ac13ad1.htm "ReadBoolAsync 方法 ")

[Write 方法](../html/79d68534-f5bc-d4fd-5b73-04626d0001c1.htm "Write 方法 ")

[WriteAsync 方法](../html/5451f32d-1f77-73c9-38a9-41845add2dd3.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| VigorHelperReadAsync 方法 |

批量读取字节数组信息，需要指定地址和长度，返回原始的字节数组  
Batch read byte array information, need to specify the address and length, return the original byte array

**命名空间：**
 [HslCommunication.Profinet.Vigor.Helper](e525e2ab-bd42-dfb9-2389-84669c014cbb.htm)  
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
	byte station,
	string address,
	ushort length
)
```

```
Public Shared Function ReadAsync ( 
	plc As IReadWriteDevice,
	station As Byte,
	address As String,
	length As UShort
) As Task(Of OperateResult(Of Byte()))
```

```
public:
static Task<OperateResult<array<unsigned char>^>^>^ ReadAsync(
	IReadWriteDevice^ plc, 
	unsigned char station, 
	String^ address, 
	unsigned short length
)
```

```
static member ReadAsync : 
        plc : IReadWriteDevice * 
        station : byte * 
        address : string * 
        length : uint16 -> Task<OperateResult<byte[]>> 
```

#### 参数

plc
:   类型：[HslCommunication.CoreIReadWriteDevice](d987afcc-f99c-6aed-3847-5b1890504635.htm)  

    [缺少 "M:HslCommunication.Profinet.Vigor.Helper.VigorHelper.ReadAsync(HslCommunication.Core.IReadWriteDevice,System.Byte,System.String,System.UInt16)" 的 <param name="plc"/> 文档]

station
:   类型：SystemByte  

    [缺少 "M:HslCommunication.Profinet.Vigor.Helper.VigorHelper.ReadAsync(HslCommunication.Core.IReadWriteDevice,System.Byte,System.String,System.UInt16)" 的 <param name="station"/> 文档]

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

支持字地址，单次最多读取64字节，支持D,SD,R,T,C的数据读取，同时地址支持携带站号信息，s=2;D100

![](../icons/SectionExpanded.png)参见

#### 引用

[VigorHelper 类](204e8b8a-20c9-5818-14e3-c361aa47fbed.htm)

[HslCommunication.Profinet.Vigor.Helper 命名空间](e525e2ab-bd42-dfb9-2389-84669c014cbb.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadBool 方法 

[原文連結](http://api.hslcommunication.cn/html/59c6ed00-7b11-a3b5-f709-94f354c31569.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Vigor.Helper](../html/e525e2ab-bd42-dfb9-2389-84669c014cbb.htm "HslCommunication.Profinet.Vigor.Helper")

[VigorHelper 类](../html/204e8b8a-20c9-5818-14e3-c361aa47fbed.htm "VigorHelper 类")

[VigorHelper 方法](../html/09b5496d-304b-b53c-2968-c981725c968c.htm "VigorHelper 方法")

[Read 方法](../html/82c88c10-df59-2fe5-9091-0b208fe5bdbf.htm "Read 方法 ")

[ReadAsync 方法](../html/12d8903c-f936-e4a9-765a-203a6a1a7d65.htm "ReadAsync 方法 ")

[ReadBool 方法](../html/59c6ed00-7b11-a3b5-f709-94f354c31569.htm "ReadBool 方法 ")

[ReadBoolAsync 方法](../html/74b8830b-2a50-fff9-53f9-802d4ac13ad1.htm "ReadBoolAsync 方法 ")

[Write 方法](../html/79d68534-f5bc-d4fd-5b73-04626d0001c1.htm "Write 方法 ")

[WriteAsync 方法](../html/5451f32d-1f77-73c9-38a9-41845add2dd3.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| VigorHelperReadBool 方法 |

批量读取Boolean数组信息，需要指定地址和长度，返回Boolean 数组  
Batch read Boolean array information, need to specify the address and length, return Boolean array

**命名空间：**
 [HslCommunication.Profinet.Vigor.Helper](e525e2ab-bd42-dfb9-2389-84669c014cbb.htm)  
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
	byte station,
	string address,
	ushort length
)
```

```
Public Shared Function ReadBool ( 
	plc As IReadWriteDevice,
	station As Byte,
	address As String,
	length As UShort
) As OperateResult(Of Boolean())
```

```
public:
static OperateResult<array<bool>^>^ ReadBool(
	IReadWriteDevice^ plc, 
	unsigned char station, 
	String^ address, 
	unsigned short length
)
```

```
static member ReadBool : 
        plc : IReadWriteDevice * 
        station : byte * 
        address : string * 
        length : uint16 -> OperateResult<bool[]> 
```

#### 参数

plc
:   类型：[HslCommunication.CoreIReadWriteDevice](d987afcc-f99c-6aed-3847-5b1890504635.htm)  

    [缺少 "M:HslCommunication.Profinet.Vigor.Helper.VigorHelper.ReadBool(HslCommunication.Core.IReadWriteDevice,System.Byte,System.String,System.UInt16)" 的 <param name="plc"/> 文档]

station
:   类型：SystemByte  

    [缺少 "M:HslCommunication.Profinet.Vigor.Helper.VigorHelper.ReadBool(HslCommunication.Core.IReadWriteDevice,System.Byte,System.String,System.UInt16)" 的 <param name="station"/> 文档]

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

需要输入位地址，最多读取1024位，支持X,Y,M,SM,S,TS(定时器触点),TC（定时器线圈）,CS(计数器触点),CC（计数器线圈）

![](../icons/SectionExpanded.png)参见

#### 引用

[VigorHelper 类](204e8b8a-20c9-5818-14e3-c361aa47fbed.htm)

[HslCommunication.Profinet.Vigor.Helper 命名空间](e525e2ab-bd42-dfb9-2389-84669c014cbb.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadBoolAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/74b8830b-2a50-fff9-53f9-802d4ac13ad1.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Vigor.Helper](../html/e525e2ab-bd42-dfb9-2389-84669c014cbb.htm "HslCommunication.Profinet.Vigor.Helper")

[VigorHelper 类](../html/204e8b8a-20c9-5818-14e3-c361aa47fbed.htm "VigorHelper 类")

[VigorHelper 方法](../html/09b5496d-304b-b53c-2968-c981725c968c.htm "VigorHelper 方法")

[Read 方法](../html/82c88c10-df59-2fe5-9091-0b208fe5bdbf.htm "Read 方法 ")

[ReadAsync 方法](../html/12d8903c-f936-e4a9-765a-203a6a1a7d65.htm "ReadAsync 方法 ")

[ReadBool 方法](../html/59c6ed00-7b11-a3b5-f709-94f354c31569.htm "ReadBool 方法 ")

[ReadBoolAsync 方法](../html/74b8830b-2a50-fff9-53f9-802d4ac13ad1.htm "ReadBoolAsync 方法 ")

[Write 方法](../html/79d68534-f5bc-d4fd-5b73-04626d0001c1.htm "Write 方法 ")

[WriteAsync 方法](../html/5451f32d-1f77-73c9-38a9-41845add2dd3.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| VigorHelperReadBoolAsync 方法 |

批量读取Boolean数组信息，需要指定地址和长度，返回Boolean 数组  
Batch read Boolean array information, need to specify the address and length, return Boolean array

**命名空间：**
 [HslCommunication.Profinet.Vigor.Helper](e525e2ab-bd42-dfb9-2389-84669c014cbb.htm)  
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
	IReadWriteDevice plc,
	byte station,
	string address,
	ushort length
)
```

```
Public Shared Function ReadBoolAsync ( 
	plc As IReadWriteDevice,
	station As Byte,
	address As String,
	length As UShort
) As Task(Of OperateResult(Of Boolean()))
```

```
public:
static Task<OperateResult<array<bool>^>^>^ ReadBoolAsync(
	IReadWriteDevice^ plc, 
	unsigned char station, 
	String^ address, 
	unsigned short length
)
```

```
static member ReadBoolAsync : 
        plc : IReadWriteDevice * 
        station : byte * 
        address : string * 
        length : uint16 -> Task<OperateResult<bool[]>> 
```

#### 参数

plc
:   类型：[HslCommunication.CoreIReadWriteDevice](d987afcc-f99c-6aed-3847-5b1890504635.htm)  

    [缺少 "M:HslCommunication.Profinet.Vigor.Helper.VigorHelper.ReadBoolAsync(HslCommunication.Core.IReadWriteDevice,System.Byte,System.String,System.UInt16)" 的 <param name="plc"/> 文档]

station
:   类型：SystemByte  

    [缺少 "M:HslCommunication.Profinet.Vigor.Helper.VigorHelper.ReadBoolAsync(HslCommunication.Core.IReadWriteDevice,System.Byte,System.String,System.UInt16)" 的 <param name="station"/> 文档]

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

需要输入位地址，最多读取1024位，支持X,Y,M,SM,S,TS(定时器触点),TC（定时器线圈）,CS(计数器触点),CC（计数器线圈）

![](../icons/SectionExpanded.png)参见

#### 引用

[VigorHelper 类](204e8b8a-20c9-5818-14e3-c361aa47fbed.htm)

[HslCommunication.Profinet.Vigor.Helper 命名空间](e525e2ab-bd42-dfb9-2389-84669c014cbb.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Write 方法 

[原文連結](http://api.hslcommunication.cn/html/79d68534-f5bc-d4fd-5b73-04626d0001c1.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Vigor.Helper](../html/e525e2ab-bd42-dfb9-2389-84669c014cbb.htm "HslCommunication.Profinet.Vigor.Helper")

[VigorHelper 类](../html/204e8b8a-20c9-5818-14e3-c361aa47fbed.htm "VigorHelper 类")

[VigorHelper 方法](../html/09b5496d-304b-b53c-2968-c981725c968c.htm "VigorHelper 方法")

[Write 方法](../html/79d68534-f5bc-d4fd-5b73-04626d0001c1.htm "Write 方法 ")

[Write 方法 (IReadWriteDevice, Byte, String, Boolean[])](../html/033e8f76-6479-2b13-9892-a4c67509b982.htm "Write 方法 (IReadWriteDevice, Byte, String, Boolean[])")

[Write 方法 (IReadWriteDevice, Byte, String, Byte[])](../html/71aeee72-e80e-2f0f-b069-841341d39f79.htm "Write 方法 (IReadWriteDevice, Byte, String, Byte[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| VigorHelperWrite 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法静态成员 | [Write(IReadWriteDevice, Byte, String, Boolean)](033e8f76-6479-2b13-9892-a4c67509b982.htm) | 批量写入Boolean数组数据，返回是否成功  Batch write Boolean array data, return whether the write was successful |
| 公共方法静态成员 | [Write(IReadWriteDevice, Byte, String, Byte)](71aeee72-e80e-2f0f-b069-841341d39f79.htm) | 写入原始的byte数组数据到指定的地址，返回是否写入成功  Write the original byte array data to the specified address, and return whether the write was successful |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[VigorHelper 类](204e8b8a-20c9-5818-14e3-c361aa47fbed.htm)

[HslCommunication.Profinet.Vigor.Helper 命名空间](e525e2ab-bd42-dfb9-2389-84669c014cbb.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Write 方法 (IReadWriteDevice, Byte, String, Boolean[])

[原文連結](http://api.hslcommunication.cn/html/033e8f76-6479-2b13-9892-a4c67509b982.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Vigor.Helper](../html/e525e2ab-bd42-dfb9-2389-84669c014cbb.htm "HslCommunication.Profinet.Vigor.Helper")

[VigorHelper 类](../html/204e8b8a-20c9-5818-14e3-c361aa47fbed.htm "VigorHelper 类")

[VigorHelper 方法](../html/09b5496d-304b-b53c-2968-c981725c968c.htm "VigorHelper 方法")

[Write 方法](../html/79d68534-f5bc-d4fd-5b73-04626d0001c1.htm "Write 方法 ")

[Write 方法 (IReadWriteDevice, Byte, String, Boolean[])](../html/033e8f76-6479-2b13-9892-a4c67509b982.htm "Write 方法 (IReadWriteDevice, Byte, String, Boolean[])")

[Write 方法 (IReadWriteDevice, Byte, String, Byte[])](../html/71aeee72-e80e-2f0f-b069-841341d39f79.htm "Write 方法 (IReadWriteDevice, Byte, String, Byte[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| VigorHelperWrite 方法 (IReadWriteDevice, Byte, String, Boolean) |

批量写入Boolean数组数据，返回是否成功  
Batch write Boolean array data, return whether the write was successful

**命名空间：**
 [HslCommunication.Profinet.Vigor.Helper](e525e2ab-bd42-dfb9-2389-84669c014cbb.htm)  
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
	byte station,
	string address,
	bool[] value
)
```

```
Public Shared Function Write ( 
	plc As IReadWriteDevice,
	station As Byte,
	address As String,
	value As Boolean()
) As OperateResult
```

```
public:
static OperateResult^ Write(
	IReadWriteDevice^ plc, 
	unsigned char station, 
	String^ address, 
	array<bool>^ value
)
```

```
static member Write : 
        plc : IReadWriteDevice * 
        station : byte * 
        address : string * 
        value : bool[] -> OperateResult 
```

#### 参数

plc
:   类型：[HslCommunication.CoreIReadWriteDevice](d987afcc-f99c-6aed-3847-5b1890504635.htm)  

    [缺少 "M:HslCommunication.Profinet.Vigor.Helper.VigorHelper.Write(HslCommunication.Core.IReadWriteDevice,System.Byte,System.String,System.Boolean[])" 的 <param name="plc"/> 文档]

station
:   类型：SystemByte  

    [缺少 "M:HslCommunication.Profinet.Vigor.Helper.VigorHelper.Write(HslCommunication.Core.IReadWriteDevice,System.Byte,System.String,System.Boolean[])" 的 <param name="station"/> 文档]

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

支持位地址的写入，支持X,Y,M,SM,S,TS(定时器触点),TC（定时器线圈）,CS(计数器触点),CC（计数器线圈）

![](../icons/SectionExpanded.png)参见

#### 引用

[VigorHelper 类](204e8b8a-20c9-5818-14e3-c361aa47fbed.htm)

[Write 重载](79d68534-f5bc-d4fd-5b73-04626d0001c1.htm)

[HslCommunication.Profinet.Vigor.Helper 命名空间](e525e2ab-bd42-dfb9-2389-84669c014cbb.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Write 方法 (IReadWriteDevice, Byte, String, Byte[])

[原文連結](http://api.hslcommunication.cn/html/71aeee72-e80e-2f0f-b069-841341d39f79.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Vigor.Helper](../html/e525e2ab-bd42-dfb9-2389-84669c014cbb.htm "HslCommunication.Profinet.Vigor.Helper")

[VigorHelper 类](../html/204e8b8a-20c9-5818-14e3-c361aa47fbed.htm "VigorHelper 类")

[VigorHelper 方法](../html/09b5496d-304b-b53c-2968-c981725c968c.htm "VigorHelper 方法")

[Write 方法](../html/79d68534-f5bc-d4fd-5b73-04626d0001c1.htm "Write 方法 ")

[Write 方法 (IReadWriteDevice, Byte, String, Boolean[])](../html/033e8f76-6479-2b13-9892-a4c67509b982.htm "Write 方法 (IReadWriteDevice, Byte, String, Boolean[])")

[Write 方法 (IReadWriteDevice, Byte, String, Byte[])](../html/71aeee72-e80e-2f0f-b069-841341d39f79.htm "Write 方法 (IReadWriteDevice, Byte, String, Byte[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| VigorHelperWrite 方法 (IReadWriteDevice, Byte, String, Byte) |

写入原始的byte数组数据到指定的地址，返回是否写入成功  
Write the original byte array data to the specified address, and return whether the write was successful

**命名空间：**
 [HslCommunication.Profinet.Vigor.Helper](e525e2ab-bd42-dfb9-2389-84669c014cbb.htm)  
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
	byte station,
	string address,
	byte[] value
)
```

```
Public Shared Function Write ( 
	plc As IReadWriteDevice,
	station As Byte,
	address As String,
	value As Byte()
) As OperateResult
```

```
public:
static OperateResult^ Write(
	IReadWriteDevice^ plc, 
	unsigned char station, 
	String^ address, 
	array<unsigned char>^ value
)
```

```
static member Write : 
        plc : IReadWriteDevice * 
        station : byte * 
        address : string * 
        value : byte[] -> OperateResult 
```

#### 参数

plc
:   类型：[HslCommunication.CoreIReadWriteDevice](d987afcc-f99c-6aed-3847-5b1890504635.htm)  

    [缺少 "M:HslCommunication.Profinet.Vigor.Helper.VigorHelper.Write(HslCommunication.Core.IReadWriteDevice,System.Byte,System.String,System.Byte[])" 的 <param name="plc"/> 文档]

station
:   类型：SystemByte  

    [缺少 "M:HslCommunication.Profinet.Vigor.Helper.VigorHelper.Write(HslCommunication.Core.IReadWriteDevice,System.Byte,System.String,System.Byte[])" 的 <param name="station"/> 文档]

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

支持字地址，单次最多读取64字节，支持D,SD,R,T,C的数据写入，其中C199~C200不能连续写入，前者是16位计数器，后者是32位计数器

![](../icons/SectionExpanded.png)参见

#### 引用

[VigorHelper 类](204e8b8a-20c9-5818-14e3-c361aa47fbed.htm)

[Write 重载](79d68534-f5bc-d4fd-5b73-04626d0001c1.htm)

[HslCommunication.Profinet.Vigor.Helper 命名空间](e525e2ab-bd42-dfb9-2389-84669c014cbb.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## WriteAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/5451f32d-1f77-73c9-38a9-41845add2dd3.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Vigor.Helper](../html/e525e2ab-bd42-dfb9-2389-84669c014cbb.htm "HslCommunication.Profinet.Vigor.Helper")

[VigorHelper 类](../html/204e8b8a-20c9-5818-14e3-c361aa47fbed.htm "VigorHelper 类")

[VigorHelper 方法](../html/09b5496d-304b-b53c-2968-c981725c968c.htm "VigorHelper 方法")

[WriteAsync 方法](../html/5451f32d-1f77-73c9-38a9-41845add2dd3.htm "WriteAsync 方法 ")

[WriteAsync 方法 (IReadWriteDevice, Byte, String, Boolean[])](../html/1009fb21-0e63-755c-9810-694fade4f4f8.htm "WriteAsync 方法 (IReadWriteDevice, Byte, String, Boolean[])")

[WriteAsync 方法 (IReadWriteDevice, Byte, String, Byte[])](../html/e64fe734-a64f-8bb1-d6f0-4284aa42312b.htm "WriteAsync 方法 (IReadWriteDevice, Byte, String, Byte[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| VigorHelperWriteAsync 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法静态成员 | [WriteAsync(IReadWriteDevice, Byte, String, Boolean)](1009fb21-0e63-755c-9810-694fade4f4f8.htm) | 批量写入Boolean数组数据，返回是否成功  Batch write Boolean array data, return whether the write was successful |
| 公共方法静态成员 | [WriteAsync(IReadWriteDevice, Byte, String, Byte)](e64fe734-a64f-8bb1-d6f0-4284aa42312b.htm) | 写入原始的byte数组数据到指定的地址，返回是否写入成功  Write the original byte array data to the specified address, and return whether the write was successful |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[VigorHelper 类](204e8b8a-20c9-5818-14e3-c361aa47fbed.htm)

[HslCommunication.Profinet.Vigor.Helper 命名空间](e525e2ab-bd42-dfb9-2389-84669c014cbb.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## WriteAsync 方法 (IReadWriteDevice, Byte, String, Boolean[])

[原文連結](http://api.hslcommunication.cn/html/1009fb21-0e63-755c-9810-694fade4f4f8.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Vigor.Helper](../html/e525e2ab-bd42-dfb9-2389-84669c014cbb.htm "HslCommunication.Profinet.Vigor.Helper")

[VigorHelper 类](../html/204e8b8a-20c9-5818-14e3-c361aa47fbed.htm "VigorHelper 类")

[VigorHelper 方法](../html/09b5496d-304b-b53c-2968-c981725c968c.htm "VigorHelper 方法")

[WriteAsync 方法](../html/5451f32d-1f77-73c9-38a9-41845add2dd3.htm "WriteAsync 方法 ")

[WriteAsync 方法 (IReadWriteDevice, Byte, String, Boolean[])](../html/1009fb21-0e63-755c-9810-694fade4f4f8.htm "WriteAsync 方法 (IReadWriteDevice, Byte, String, Boolean[])")

[WriteAsync 方法 (IReadWriteDevice, Byte, String, Byte[])](../html/e64fe734-a64f-8bb1-d6f0-4284aa42312b.htm "WriteAsync 方法 (IReadWriteDevice, Byte, String, Byte[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| VigorHelperWriteAsync 方法 (IReadWriteDevice, Byte, String, Boolean) |

批量写入Boolean数组数据，返回是否成功  
Batch write Boolean array data, return whether the write was successful

**命名空间：**
 [HslCommunication.Profinet.Vigor.Helper](e525e2ab-bd42-dfb9-2389-84669c014cbb.htm)  
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
	IReadWriteDevice plc,
	byte station,
	string address,
	bool[] value
)
```

```
Public Shared Function WriteAsync ( 
	plc As IReadWriteDevice,
	station As Byte,
	address As String,
	value As Boolean()
) As Task(Of OperateResult)
```

```
public:
static Task<OperateResult^>^ WriteAsync(
	IReadWriteDevice^ plc, 
	unsigned char station, 
	String^ address, 
	array<bool>^ value
)
```

```
static member WriteAsync : 
        plc : IReadWriteDevice * 
        station : byte * 
        address : string * 
        value : bool[] -> Task<OperateResult> 
```

#### 参数

plc
:   类型：[HslCommunication.CoreIReadWriteDevice](d987afcc-f99c-6aed-3847-5b1890504635.htm)  

    [缺少 "M:HslCommunication.Profinet.Vigor.Helper.VigorHelper.WriteAsync(HslCommunication.Core.IReadWriteDevice,System.Byte,System.String,System.Boolean[])" 的 <param name="plc"/> 文档]

station
:   类型：SystemByte  

    [缺少 "M:HslCommunication.Profinet.Vigor.Helper.VigorHelper.WriteAsync(HslCommunication.Core.IReadWriteDevice,System.Byte,System.String,System.Boolean[])" 的 <param name="station"/> 文档]

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

支持位地址的写入，支持X,Y,M,SM,S,TS(定时器触点),TC（定时器线圈）,CS(计数器触点),CC（计数器线圈）

![](../icons/SectionExpanded.png)参见

#### 引用

[VigorHelper 类](204e8b8a-20c9-5818-14e3-c361aa47fbed.htm)

[WriteAsync 重载](5451f32d-1f77-73c9-38a9-41845add2dd3.htm)

[HslCommunication.Profinet.Vigor.Helper 命名空间](e525e2ab-bd42-dfb9-2389-84669c014cbb.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## WriteAsync 方法 (IReadWriteDevice, Byte, String, Byte[])

[原文連結](http://api.hslcommunication.cn/html/e64fe734-a64f-8bb1-d6f0-4284aa42312b.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Vigor.Helper](../html/e525e2ab-bd42-dfb9-2389-84669c014cbb.htm "HslCommunication.Profinet.Vigor.Helper")

[VigorHelper 类](../html/204e8b8a-20c9-5818-14e3-c361aa47fbed.htm "VigorHelper 类")

[VigorHelper 方法](../html/09b5496d-304b-b53c-2968-c981725c968c.htm "VigorHelper 方法")

[WriteAsync 方法](../html/5451f32d-1f77-73c9-38a9-41845add2dd3.htm "WriteAsync 方法 ")

[WriteAsync 方法 (IReadWriteDevice, Byte, String, Boolean[])](../html/1009fb21-0e63-755c-9810-694fade4f4f8.htm "WriteAsync 方法 (IReadWriteDevice, Byte, String, Boolean[])")

[WriteAsync 方法 (IReadWriteDevice, Byte, String, Byte[])](../html/e64fe734-a64f-8bb1-d6f0-4284aa42312b.htm "WriteAsync 方法 (IReadWriteDevice, Byte, String, Byte[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| VigorHelperWriteAsync 方法 (IReadWriteDevice, Byte, String, Byte) |

写入原始的byte数组数据到指定的地址，返回是否写入成功  
Write the original byte array data to the specified address, and return whether the write was successful

**命名空间：**
 [HslCommunication.Profinet.Vigor.Helper](e525e2ab-bd42-dfb9-2389-84669c014cbb.htm)  
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
	IReadWriteDevice plc,
	byte station,
	string address,
	byte[] value
)
```

```
Public Shared Function WriteAsync ( 
	plc As IReadWriteDevice,
	station As Byte,
	address As String,
	value As Byte()
) As Task(Of OperateResult)
```

```
public:
static Task<OperateResult^>^ WriteAsync(
	IReadWriteDevice^ plc, 
	unsigned char station, 
	String^ address, 
	array<unsigned char>^ value
)
```

```
static member WriteAsync : 
        plc : IReadWriteDevice * 
        station : byte * 
        address : string * 
        value : byte[] -> Task<OperateResult> 
```

#### 参数

plc
:   类型：[HslCommunication.CoreIReadWriteDevice](d987afcc-f99c-6aed-3847-5b1890504635.htm)  

    [缺少 "M:HslCommunication.Profinet.Vigor.Helper.VigorHelper.WriteAsync(HslCommunication.Core.IReadWriteDevice,System.Byte,System.String,System.Byte[])" 的 <param name="plc"/> 文档]

station
:   类型：SystemByte  

    [缺少 "M:HslCommunication.Profinet.Vigor.Helper.VigorHelper.WriteAsync(HslCommunication.Core.IReadWriteDevice,System.Byte,System.String,System.Byte[])" 的 <param name="station"/> 文档]

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

支持字地址，单次最多读取64字节，支持D,SD,R,T,C的数据写入，其中C199~C200不能连续写入，前者是16位计数器，后者是32位计数器

![](../icons/SectionExpanded.png)参见

#### 引用

[VigorHelper 类](204e8b8a-20c9-5818-14e3-c361aa47fbed.htm)

[WriteAsync 重载](5451f32d-1f77-73c9-38a9-41845add2dd3.htm)

[HslCommunication.Profinet.Vigor.Helper 命名空间](e525e2ab-bd42-dfb9-2389-84669c014cbb.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## VigorVsHelper 类

[原文連結](http://api.hslcommunication.cn/html/383f1b2a-980f-adaa-16af-3ccaca87cb84.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Vigor.Helper](../html/e525e2ab-bd42-dfb9-2389-84669c014cbb.htm "HslCommunication.Profinet.Vigor.Helper")

[VigorVsHelper 类](../html/383f1b2a-980f-adaa-16af-3ccaca87cb84.htm "VigorVsHelper 类")

[VigorVsHelper 构造函数](../html/5393f9ba-4114-1ba2-36b3-f39c13e093fe.htm "VigorVsHelper 构造函数 ")

[VigorVsHelper 方法](../html/72a44e2b-4a70-1fab-4f2b-7ac002491982.htm "VigorVsHelper 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| VigorVsHelper 类 |

丰炜PLC的辅助类对象

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  HslCommunication.Profinet.Vigor.HelperVigorVsHelper

**命名空间：**
 [HslCommunication.Profinet.Vigor.Helper](e525e2ab-bd42-dfb9-2389-84669c014cbb.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class VigorVsHelper
```

```
Public Class VigorVsHelper
```

```
public ref class VigorVsHelper
```

```
type VigorVsHelper =  class end
```

VigorVsHelper 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [VigorVsHelper](5393f9ba-4114-1ba2-36b3-f39c13e093fe.htm) | 初始化 VigorVsHelper 类的一个新实例 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法静态成员 | [BuildReadCommand](b91ecb44-b6fd-7d55-8d4b-9b0af18c9e50.htm) | 构建读取的报文命令，对于字地址，单次最多读取64字节，支持D,SD,R,T,C的数据读取，对于位地址，最多读取1024位，支持X,Y,M,SM,S,TS(定时器触点),TC（定时器线圈）,CS(计数器触点),CC（计数器线圈）  Construct a read message command. For word addresses, up to 64 bytes can be read at a time, and data reading of D, SD, R, T, and C is supported. For bit addresses, up to 1024 bits are read, and X, Y are supported. , M, SM, S, TS (timer contact), TC (timer coil), CS (counter contact), CC (counter coil) |
| 公共方法静态成员 | [BuildWriteBoolCommand](9add88b4-2c3f-152b-b77d-f0a564f581ee.htm) | 构建以位单位写入的报文，单次最多写入1024bit，支持X,Y,M,SM,S,TS(定时器触点),TC（定时器线圈）,CS(计数器触点),CC（计数器线圈） |
| 公共方法静态成员 | [BuildWriteWordCommand](eb8fb0b1-24d2-75ab-02f2-ecabe5f52407.htm) | 构建以字单位写入的报文，单次最多写入64个word，地址支持 D,SD,R,T,C，对于C200~C255,是属于32位的计数器  Construct a message written in word units, and write up to 64 words in a single time. The address supports D, SD, R, T, C. For C200~C255, it is a 32-bit counter |
| 公共方法静态成员 | [CheckResponseContent](896619a2-3b79-b155-7c62-4c217ede52df.htm) | 检查从PLC返回的报文是否正确，以及提取出正确的结果数据 |
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

[HslCommunication.Profinet.Vigor.Helper 命名空间](e525e2ab-bd42-dfb9-2389-84669c014cbb.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## VigorVsHelper 构造函数 

[原文連結](http://api.hslcommunication.cn/html/5393f9ba-4114-1ba2-36b3-f39c13e093fe.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Vigor.Helper](../html/e525e2ab-bd42-dfb9-2389-84669c014cbb.htm "HslCommunication.Profinet.Vigor.Helper")

[VigorVsHelper 类](../html/383f1b2a-980f-adaa-16af-3ccaca87cb84.htm "VigorVsHelper 类")

[VigorVsHelper 构造函数](../html/5393f9ba-4114-1ba2-36b3-f39c13e093fe.htm "VigorVsHelper 构造函数 ")

[VigorVsHelper 方法](../html/72a44e2b-4a70-1fab-4f2b-7ac002491982.htm "VigorVsHelper 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| VigorVsHelper 构造函数 |

初始化 [VigorVsHelper](383f1b2a-980f-adaa-16af-3ccaca87cb84.htm) 类的一个新实例

**命名空间：**
 [HslCommunication.Profinet.Vigor.Helper](e525e2ab-bd42-dfb9-2389-84669c014cbb.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public VigorVsHelper()
```

```
Public Sub New
```

```
public:
VigorVsHelper()
```

```
new : unit -> VigorVsHelper
```

![](../icons/SectionExpanded.png)参见

#### 引用

[VigorVsHelper 类](383f1b2a-980f-adaa-16af-3ccaca87cb84.htm)

[HslCommunication.Profinet.Vigor.Helper 命名空间](e525e2ab-bd42-dfb9-2389-84669c014cbb.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## VigorVsHelper 方法

[原文連結](http://api.hslcommunication.cn/html/72a44e2b-4a70-1fab-4f2b-7ac002491982.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Vigor.Helper](../html/e525e2ab-bd42-dfb9-2389-84669c014cbb.htm "HslCommunication.Profinet.Vigor.Helper")

[VigorVsHelper 类](../html/383f1b2a-980f-adaa-16af-3ccaca87cb84.htm "VigorVsHelper 类")

[VigorVsHelper 方法](../html/72a44e2b-4a70-1fab-4f2b-7ac002491982.htm "VigorVsHelper 方法")

[BuildReadCommand 方法](../html/b91ecb44-b6fd-7d55-8d4b-9b0af18c9e50.htm "BuildReadCommand 方法 ")

[BuildWriteBoolCommand 方法](../html/9add88b4-2c3f-152b-b77d-f0a564f581ee.htm "BuildWriteBoolCommand 方法 ")

[BuildWriteWordCommand 方法](../html/eb8fb0b1-24d2-75ab-02f2-ecabe5f52407.htm "BuildWriteWordCommand 方法 ")

[CheckResponseContent 方法](../html/896619a2-3b79-b155-7c62-4c217ede52df.htm "CheckResponseContent 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| VigorVsHelper 方法 |

[VigorVsHelper](383f1b2a-980f-adaa-16af-3ccaca87cb84.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法静态成员 | [BuildReadCommand](b91ecb44-b6fd-7d55-8d4b-9b0af18c9e50.htm) | 构建读取的报文命令，对于字地址，单次最多读取64字节，支持D,SD,R,T,C的数据读取，对于位地址，最多读取1024位，支持X,Y,M,SM,S,TS(定时器触点),TC（定时器线圈）,CS(计数器触点),CC（计数器线圈）  Construct a read message command. For word addresses, up to 64 bytes can be read at a time, and data reading of D, SD, R, T, and C is supported. For bit addresses, up to 1024 bits are read, and X, Y are supported. , M, SM, S, TS (timer contact), TC (timer coil), CS (counter contact), CC (counter coil) |
| 公共方法静态成员 | [BuildWriteBoolCommand](9add88b4-2c3f-152b-b77d-f0a564f581ee.htm) | 构建以位单位写入的报文，单次最多写入1024bit，支持X,Y,M,SM,S,TS(定时器触点),TC（定时器线圈）,CS(计数器触点),CC（计数器线圈） |
| 公共方法静态成员 | [BuildWriteWordCommand](eb8fb0b1-24d2-75ab-02f2-ecabe5f52407.htm) | 构建以字单位写入的报文，单次最多写入64个word，地址支持 D,SD,R,T,C，对于C200~C255,是属于32位的计数器  Construct a message written in word units, and write up to 64 words in a single time. The address supports D, SD, R, T, C. For C200~C255, it is a 32-bit counter |
| 公共方法静态成员 | [CheckResponseContent](896619a2-3b79-b155-7c62-4c217ede52df.htm) | 检查从PLC返回的报文是否正确，以及提取出正确的结果数据 |
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

[VigorVsHelper 类](383f1b2a-980f-adaa-16af-3ccaca87cb84.htm)

[HslCommunication.Profinet.Vigor.Helper 命名空间](e525e2ab-bd42-dfb9-2389-84669c014cbb.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## BuildReadCommand 方法 

[原文連結](http://api.hslcommunication.cn/html/b91ecb44-b6fd-7d55-8d4b-9b0af18c9e50.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Vigor.Helper](../html/e525e2ab-bd42-dfb9-2389-84669c014cbb.htm "HslCommunication.Profinet.Vigor.Helper")

[VigorVsHelper 类](../html/383f1b2a-980f-adaa-16af-3ccaca87cb84.htm "VigorVsHelper 类")

[VigorVsHelper 方法](../html/72a44e2b-4a70-1fab-4f2b-7ac002491982.htm "VigorVsHelper 方法")

[BuildReadCommand 方法](../html/b91ecb44-b6fd-7d55-8d4b-9b0af18c9e50.htm "BuildReadCommand 方法 ")

[BuildWriteBoolCommand 方法](../html/9add88b4-2c3f-152b-b77d-f0a564f581ee.htm "BuildWriteBoolCommand 方法 ")

[BuildWriteWordCommand 方法](../html/eb8fb0b1-24d2-75ab-02f2-ecabe5f52407.htm "BuildWriteWordCommand 方法 ")

[CheckResponseContent 方法](../html/896619a2-3b79-b155-7c62-4c217ede52df.htm "CheckResponseContent 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| VigorVsHelperBuildReadCommand 方法 |

构建读取的报文命令，对于字地址，单次最多读取64字节，支持D,SD,R,T,C的数据读取，对于位地址，最多读取1024位，支持X,Y,M,SM,S,TS(定时器触点),TC（定时器线圈）,CS(计数器触点),CC（计数器线圈）  
Construct a read message command. For word addresses, up to 64 bytes can be read at a time, and data reading of D, SD, R, T, and C is supported. For bit addresses,
up to 1024 bits are read, and X, Y are supported. , M, SM, S, TS (timer contact), TC (timer coil), CS (counter contact), CC (counter coil)

**命名空间：**
 [HslCommunication.Profinet.Vigor.Helper](e525e2ab-bd42-dfb9-2389-84669c014cbb.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult<List<byte[]>> BuildReadCommand(
	byte station,
	string address,
	ushort length,
	bool isBool
)
```

```
Public Shared Function BuildReadCommand ( 
	station As Byte,
	address As String,
	length As UShort,
	isBool As Boolean
) As OperateResult(Of List(Of Byte()))
```

```
public:
static OperateResult<List<array<unsigned char>^>^>^ BuildReadCommand(
	unsigned char station, 
	String^ address, 
	unsigned short length, 
	bool isBool
)
```

```
static member BuildReadCommand : 
        station : byte * 
        address : string * 
        length : uint16 * 
        isBool : bool -> OperateResult<List<byte[]>> 
```

#### 参数

station
:   类型：SystemByte  
    站号信息

address
:   类型：SystemString  
    PLC的数据地址

length
:   类型：SystemUInt16  
    读取的长度

isBool
:   类型：SystemBoolean  
    是否进行位读取

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)ListByte  
完整的读取的报文信息

![](../icons/SectionExpanded.png)参见

#### 引用

[VigorVsHelper 类](383f1b2a-980f-adaa-16af-3ccaca87cb84.htm)

[HslCommunication.Profinet.Vigor.Helper 命名空间](e525e2ab-bd42-dfb9-2389-84669c014cbb.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## BuildWriteBoolCommand 方法 

[原文連結](http://api.hslcommunication.cn/html/9add88b4-2c3f-152b-b77d-f0a564f581ee.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Vigor.Helper](../html/e525e2ab-bd42-dfb9-2389-84669c014cbb.htm "HslCommunication.Profinet.Vigor.Helper")

[VigorVsHelper 类](../html/383f1b2a-980f-adaa-16af-3ccaca87cb84.htm "VigorVsHelper 类")

[VigorVsHelper 方法](../html/72a44e2b-4a70-1fab-4f2b-7ac002491982.htm "VigorVsHelper 方法")

[BuildReadCommand 方法](../html/b91ecb44-b6fd-7d55-8d4b-9b0af18c9e50.htm "BuildReadCommand 方法 ")

[BuildWriteBoolCommand 方法](../html/9add88b4-2c3f-152b-b77d-f0a564f581ee.htm "BuildWriteBoolCommand 方法 ")

[BuildWriteWordCommand 方法](../html/eb8fb0b1-24d2-75ab-02f2-ecabe5f52407.htm "BuildWriteWordCommand 方法 ")

[CheckResponseContent 方法](../html/896619a2-3b79-b155-7c62-4c217ede52df.htm "CheckResponseContent 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| VigorVsHelperBuildWriteBoolCommand 方法 |

构建以位单位写入的报文，单次最多写入1024bit，支持X,Y,M,SM,S,TS(定时器触点),TC（定时器线圈）,CS(计数器触点),CC（计数器线圈）

**命名空间：**
 [HslCommunication.Profinet.Vigor.Helper](e525e2ab-bd42-dfb9-2389-84669c014cbb.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult<byte[]> BuildWriteBoolCommand(
	byte station,
	string address,
	bool[] value
)
```

```
Public Shared Function BuildWriteBoolCommand ( 
	station As Byte,
	address As String,
	value As Boolean()
) As OperateResult(Of Byte())
```

```
public:
static OperateResult<array<unsigned char>^>^ BuildWriteBoolCommand(
	unsigned char station, 
	String^ address, 
	array<bool>^ value
)
```

```
static member BuildWriteBoolCommand : 
        station : byte * 
        address : string * 
        value : bool[] -> OperateResult<byte[]> 
```

#### 参数

station
:   类型：SystemByte  
    站号信息

address
:   类型：SystemString  
    PLC的地址

value
:   类型：SystemBoolean  
    等待写入的bool数组

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
写入位数据的完整报文信息

![](../icons/SectionExpanded.png)参见

#### 引用

[VigorVsHelper 类](383f1b2a-980f-adaa-16af-3ccaca87cb84.htm)

[HslCommunication.Profinet.Vigor.Helper 命名空间](e525e2ab-bd42-dfb9-2389-84669c014cbb.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## BuildWriteWordCommand 方法 

[原文連結](http://api.hslcommunication.cn/html/eb8fb0b1-24d2-75ab-02f2-ecabe5f52407.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Vigor.Helper](../html/e525e2ab-bd42-dfb9-2389-84669c014cbb.htm "HslCommunication.Profinet.Vigor.Helper")

[VigorVsHelper 类](../html/383f1b2a-980f-adaa-16af-3ccaca87cb84.htm "VigorVsHelper 类")

[VigorVsHelper 方法](../html/72a44e2b-4a70-1fab-4f2b-7ac002491982.htm "VigorVsHelper 方法")

[BuildReadCommand 方法](../html/b91ecb44-b6fd-7d55-8d4b-9b0af18c9e50.htm "BuildReadCommand 方法 ")

[BuildWriteBoolCommand 方法](../html/9add88b4-2c3f-152b-b77d-f0a564f581ee.htm "BuildWriteBoolCommand 方法 ")

[BuildWriteWordCommand 方法](../html/eb8fb0b1-24d2-75ab-02f2-ecabe5f52407.htm "BuildWriteWordCommand 方法 ")

[CheckResponseContent 方法](../html/896619a2-3b79-b155-7c62-4c217ede52df.htm "CheckResponseContent 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| VigorVsHelperBuildWriteWordCommand 方法 |

构建以字单位写入的报文，单次最多写入64个word，地址支持 D,SD,R,T,C，对于C200~C255,是属于32位的计数器  
Construct a message written in word units, and write up to 64 words in a single time. The address supports D, SD, R, T, C. For C200~C255, it is a 32-bit counter

**命名空间：**
 [HslCommunication.Profinet.Vigor.Helper](e525e2ab-bd42-dfb9-2389-84669c014cbb.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult<byte[]> BuildWriteWordCommand(
	byte station,
	string address,
	byte[] value
)
```

```
Public Shared Function BuildWriteWordCommand ( 
	station As Byte,
	address As String,
	value As Byte()
) As OperateResult(Of Byte())
```

```
public:
static OperateResult<array<unsigned char>^>^ BuildWriteWordCommand(
	unsigned char station, 
	String^ address, 
	array<unsigned char>^ value
)
```

```
static member BuildWriteWordCommand : 
        station : byte * 
        address : string * 
        value : byte[] -> OperateResult<byte[]> 
```

#### 参数

station
:   类型：SystemByte  
    站号信息

address
:   类型：SystemString  
    PLC的地址

value
:   类型：SystemByte  
    写入的原始数据

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
写入命令的完整报文

![](../icons/SectionExpanded.png)参见

#### 引用

[VigorVsHelper 类](383f1b2a-980f-adaa-16af-3ccaca87cb84.htm)

[HslCommunication.Profinet.Vigor.Helper 命名空间](e525e2ab-bd42-dfb9-2389-84669c014cbb.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## CheckResponseContent 方法 

[原文連結](http://api.hslcommunication.cn/html/896619a2-3b79-b155-7c62-4c217ede52df.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Vigor.Helper](../html/e525e2ab-bd42-dfb9-2389-84669c014cbb.htm "HslCommunication.Profinet.Vigor.Helper")

[VigorVsHelper 类](../html/383f1b2a-980f-adaa-16af-3ccaca87cb84.htm "VigorVsHelper 类")

[VigorVsHelper 方法](../html/72a44e2b-4a70-1fab-4f2b-7ac002491982.htm "VigorVsHelper 方法")

[BuildReadCommand 方法](../html/b91ecb44-b6fd-7d55-8d4b-9b0af18c9e50.htm "BuildReadCommand 方法 ")

[BuildWriteBoolCommand 方法](../html/9add88b4-2c3f-152b-b77d-f0a564f581ee.htm "BuildWriteBoolCommand 方法 ")

[BuildWriteWordCommand 方法](../html/eb8fb0b1-24d2-75ab-02f2-ecabe5f52407.htm "BuildWriteWordCommand 方法 ")

[CheckResponseContent 方法](../html/896619a2-3b79-b155-7c62-4c217ede52df.htm "CheckResponseContent 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| VigorVsHelperCheckResponseContent 方法 |

检查从PLC返回的报文是否正确，以及提取出正确的结果数据

**命名空间：**
 [HslCommunication.Profinet.Vigor.Helper](e525e2ab-bd42-dfb9-2389-84669c014cbb.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult<byte[]> CheckResponseContent(
	byte[] response
)
```

```
Public Shared Function CheckResponseContent ( 
	response As Byte()
) As OperateResult(Of Byte())
```

```
public:
static OperateResult<array<unsigned char>^>^ CheckResponseContent(
	array<unsigned char>^ response
)
```

```
static member CheckResponseContent : 
        response : byte[] -> OperateResult<byte[]> 
```

#### 参数

response
:   类型：SystemByte  
    PLC返回的报文

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
提取的结果数据内容

![](../icons/SectionExpanded.png)参见

#### 引用

[VigorVsHelper 类](383f1b2a-980f-adaa-16af-3ccaca87cb84.htm)

[HslCommunication.Profinet.Vigor.Helper 命名空间](e525e2ab-bd42-dfb9-2389-84669c014cbb.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)