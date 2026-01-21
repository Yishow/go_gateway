# HslCommunication - HslCommunication.Profinet.Yamatake.Helper

> 分類頁數: 9



---
## HslCommunication.Profinet.Yamatake.Helper

[原文連結](http://api.hslcommunication.cn/html/86528ff3-8972-57d5-e768-c195b2888ccd.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Yamatake.Helper](../html/86528ff3-8972-57d5-e768-c195b2888ccd.htm "HslCommunication.Profinet.Yamatake.Helper")

[DigitronCPLHelper 类](../html/fa8a8b36-d2ae-5f51-1721-64f5163ade4e.htm "DigitronCPLHelper 类")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslCommunication.Profinet.Yamatake.Helper 命名空间 |

[缺少 "N:HslCommunication.Profinet.Yamatake.Helper" 的 <summary> 文档]

![](../icons/SectionExpanded.png)类

|  | 类 | 说明 |
| --- | --- | --- |
| 公共类 | [DigitronCPLHelper](fa8a8b36-d2ae-5f51-1721-64f5163ade4e.htm) | 辅助类方法 |

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## DigitronCPLHelper 类

[原文連結](http://api.hslcommunication.cn/html/fa8a8b36-d2ae-5f51-1721-64f5163ade4e.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Yamatake.Helper](../html/86528ff3-8972-57d5-e768-c195b2888ccd.htm "HslCommunication.Profinet.Yamatake.Helper")

[DigitronCPLHelper 类](../html/fa8a8b36-d2ae-5f51-1721-64f5163ade4e.htm "DigitronCPLHelper 类")

[DigitronCPLHelper 构造函数](../html/62212f46-a1bf-d210-0298-3c7155e09640.htm "DigitronCPLHelper 构造函数 ")

[DigitronCPLHelper 方法](../html/48064c28-7c39-a2bf-ad5e-b8f29d0f9dd1.htm "DigitronCPLHelper 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DigitronCPLHelper 类 |

辅助类方法

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  HslCommunication.Profinet.Yamatake.HelperDigitronCPLHelper

**命名空间：**
 [HslCommunication.Profinet.Yamatake.Helper](86528ff3-8972-57d5-e768-c195b2888ccd.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class DigitronCPLHelper
```

```
Public Class DigitronCPLHelper
```

```
public ref class DigitronCPLHelper
```

```
type DigitronCPLHelper =  class end
```

DigitronCPLHelper 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [DigitronCPLHelper](62212f46-a1bf-d210-0298-3c7155e09640.htm) | 初始化 DigitronCPLHelper 类的一个新实例 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法静态成员 | [BuildReadCommand](b4d41cf6-2776-3844-9da8-5593814f79c5.htm) | 构建写入操作的报文信息 |
| 公共方法静态成员 | [BuildWriteCommand](41a2e62a-7d9d-7b58-ee00-8f73ad538114.htm) | 构建写入操作的命令报文 |
| 公共方法 | Equals | (继承自 Object。) |
| 公共方法静态成员 | [ExtraActualResponse](1806d7e3-d804-ebda-4571-0e5823b3d00d.htm) | 从反馈的数据内容中解析出真实的数据信息 |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法静态成员 | [GetErrorText](abe4aa67-d36b-4f18-3d45-96142d78120b.htm) | 根据错误码获取到相关的错误代号信息 |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法静态成员 | [PackResponseContent](04df21d0-d08d-24ab-02db-0a1f259581a3.htm) | 用于服务器反馈的数据的报文打包操作 |
| 公共方法 | ToString | (继承自 Object。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.Profinet.Yamatake.Helper 命名空间](86528ff3-8972-57d5-e768-c195b2888ccd.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## DigitronCPLHelper 构造函数 

[原文連結](http://api.hslcommunication.cn/html/62212f46-a1bf-d210-0298-3c7155e09640.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Yamatake.Helper](../html/86528ff3-8972-57d5-e768-c195b2888ccd.htm "HslCommunication.Profinet.Yamatake.Helper")

[DigitronCPLHelper 类](../html/fa8a8b36-d2ae-5f51-1721-64f5163ade4e.htm "DigitronCPLHelper 类")

[DigitronCPLHelper 构造函数](../html/62212f46-a1bf-d210-0298-3c7155e09640.htm "DigitronCPLHelper 构造函数 ")

[DigitronCPLHelper 方法](../html/48064c28-7c39-a2bf-ad5e-b8f29d0f9dd1.htm "DigitronCPLHelper 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DigitronCPLHelper 构造函数 |

初始化 [DigitronCPLHelper](fa8a8b36-d2ae-5f51-1721-64f5163ade4e.htm) 类的一个新实例

**命名空间：**
 [HslCommunication.Profinet.Yamatake.Helper](86528ff3-8972-57d5-e768-c195b2888ccd.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public DigitronCPLHelper()
```

```
Public Sub New
```

```
public:
DigitronCPLHelper()
```

```
new : unit -> DigitronCPLHelper
```

![](../icons/SectionExpanded.png)参见

#### 引用

[DigitronCPLHelper 类](fa8a8b36-d2ae-5f51-1721-64f5163ade4e.htm)

[HslCommunication.Profinet.Yamatake.Helper 命名空间](86528ff3-8972-57d5-e768-c195b2888ccd.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## DigitronCPLHelper 方法

[原文連結](http://api.hslcommunication.cn/html/48064c28-7c39-a2bf-ad5e-b8f29d0f9dd1.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Yamatake.Helper](../html/86528ff3-8972-57d5-e768-c195b2888ccd.htm "HslCommunication.Profinet.Yamatake.Helper")

[DigitronCPLHelper 类](../html/fa8a8b36-d2ae-5f51-1721-64f5163ade4e.htm "DigitronCPLHelper 类")

[DigitronCPLHelper 方法](../html/48064c28-7c39-a2bf-ad5e-b8f29d0f9dd1.htm "DigitronCPLHelper 方法")

[BuildReadCommand 方法](../html/b4d41cf6-2776-3844-9da8-5593814f79c5.htm "BuildReadCommand 方法 ")

[BuildWriteCommand 方法](../html/41a2e62a-7d9d-7b58-ee00-8f73ad538114.htm "BuildWriteCommand 方法 ")

[ExtraActualResponse 方法](../html/1806d7e3-d804-ebda-4571-0e5823b3d00d.htm "ExtraActualResponse 方法 ")

[GetErrorText 方法](../html/abe4aa67-d36b-4f18-3d45-96142d78120b.htm "GetErrorText 方法 ")

[PackResponseContent 方法](../html/04df21d0-d08d-24ab-02db-0a1f259581a3.htm "PackResponseContent 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DigitronCPLHelper 方法 |

[DigitronCPLHelper](fa8a8b36-d2ae-5f51-1721-64f5163ade4e.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法静态成员 | [BuildReadCommand](b4d41cf6-2776-3844-9da8-5593814f79c5.htm) | 构建写入操作的报文信息 |
| 公共方法静态成员 | [BuildWriteCommand](41a2e62a-7d9d-7b58-ee00-8f73ad538114.htm) | 构建写入操作的命令报文 |
| 公共方法 | Equals | (继承自 Object。) |
| 公共方法静态成员 | [ExtraActualResponse](1806d7e3-d804-ebda-4571-0e5823b3d00d.htm) | 从反馈的数据内容中解析出真实的数据信息 |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法静态成员 | [GetErrorText](abe4aa67-d36b-4f18-3d45-96142d78120b.htm) | 根据错误码获取到相关的错误代号信息 |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法静态成员 | [PackResponseContent](04df21d0-d08d-24ab-02db-0a1f259581a3.htm) | 用于服务器反馈的数据的报文打包操作 |
| 公共方法 | ToString | (继承自 Object。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[DigitronCPLHelper 类](fa8a8b36-d2ae-5f51-1721-64f5163ade4e.htm)

[HslCommunication.Profinet.Yamatake.Helper 命名空间](86528ff3-8972-57d5-e768-c195b2888ccd.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## BuildReadCommand 方法 

[原文連結](http://api.hslcommunication.cn/html/b4d41cf6-2776-3844-9da8-5593814f79c5.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Yamatake.Helper](../html/86528ff3-8972-57d5-e768-c195b2888ccd.htm "HslCommunication.Profinet.Yamatake.Helper")

[DigitronCPLHelper 类](../html/fa8a8b36-d2ae-5f51-1721-64f5163ade4e.htm "DigitronCPLHelper 类")

[DigitronCPLHelper 方法](../html/48064c28-7c39-a2bf-ad5e-b8f29d0f9dd1.htm "DigitronCPLHelper 方法")

[BuildReadCommand 方法](../html/b4d41cf6-2776-3844-9da8-5593814f79c5.htm "BuildReadCommand 方法 ")

[BuildWriteCommand 方法](../html/41a2e62a-7d9d-7b58-ee00-8f73ad538114.htm "BuildWriteCommand 方法 ")

[ExtraActualResponse 方法](../html/1806d7e3-d804-ebda-4571-0e5823b3d00d.htm "ExtraActualResponse 方法 ")

[GetErrorText 方法](../html/abe4aa67-d36b-4f18-3d45-96142d78120b.htm "GetErrorText 方法 ")

[PackResponseContent 方法](../html/04df21d0-d08d-24ab-02db-0a1f259581a3.htm "PackResponseContent 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DigitronCPLHelperBuildReadCommand 方法 |

构建写入操作的报文信息

**命名空间：**
 [HslCommunication.Profinet.Yamatake.Helper](86528ff3-8972-57d5-e768-c195b2888ccd.htm)  
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
	ushort length
)
```

```
Public Shared Function BuildReadCommand ( 
	station As Byte,
	address As String,
	length As UShort
) As OperateResult(Of Byte())
```

```
public:
static OperateResult<array<unsigned char>^>^ BuildReadCommand(
	unsigned char station, 
	String^ address, 
	unsigned short length
)
```

```
static member BuildReadCommand : 
        station : byte * 
        address : string * 
        length : uint16 -> OperateResult<byte[]> 
```

#### 参数

station
:   类型：SystemByte  
    站号

address
:   类型：SystemString  
    地址

length
:   类型：SystemUInt16  
    长度的长度

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
报文内容

![](../icons/SectionExpanded.png)参见

#### 引用

[DigitronCPLHelper 类](fa8a8b36-d2ae-5f51-1721-64f5163ade4e.htm)

[HslCommunication.Profinet.Yamatake.Helper 命名空间](86528ff3-8972-57d5-e768-c195b2888ccd.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## BuildWriteCommand 方法 

[原文連結](http://api.hslcommunication.cn/html/41a2e62a-7d9d-7b58-ee00-8f73ad538114.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Yamatake.Helper](../html/86528ff3-8972-57d5-e768-c195b2888ccd.htm "HslCommunication.Profinet.Yamatake.Helper")

[DigitronCPLHelper 类](../html/fa8a8b36-d2ae-5f51-1721-64f5163ade4e.htm "DigitronCPLHelper 类")

[DigitronCPLHelper 方法](../html/48064c28-7c39-a2bf-ad5e-b8f29d0f9dd1.htm "DigitronCPLHelper 方法")

[BuildReadCommand 方法](../html/b4d41cf6-2776-3844-9da8-5593814f79c5.htm "BuildReadCommand 方法 ")

[BuildWriteCommand 方法](../html/41a2e62a-7d9d-7b58-ee00-8f73ad538114.htm "BuildWriteCommand 方法 ")

[ExtraActualResponse 方法](../html/1806d7e3-d804-ebda-4571-0e5823b3d00d.htm "ExtraActualResponse 方法 ")

[GetErrorText 方法](../html/abe4aa67-d36b-4f18-3d45-96142d78120b.htm "GetErrorText 方法 ")

[PackResponseContent 方法](../html/04df21d0-d08d-24ab-02db-0a1f259581a3.htm "PackResponseContent 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DigitronCPLHelperBuildWriteCommand 方法 |

构建写入操作的命令报文

**命名空间：**
 [HslCommunication.Profinet.Yamatake.Helper](86528ff3-8972-57d5-e768-c195b2888ccd.htm)  
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
    站号信息

address
:   类型：SystemString  
    数据的地址

value
:   类型：SystemByte  
    等待写入的值

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
写入的报文命令

![](../icons/SectionExpanded.png)参见

#### 引用

[DigitronCPLHelper 类](fa8a8b36-d2ae-5f51-1721-64f5163ade4e.htm)

[HslCommunication.Profinet.Yamatake.Helper 命名空间](86528ff3-8972-57d5-e768-c195b2888ccd.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ExtraActualResponse 方法 

[原文連結](http://api.hslcommunication.cn/html/1806d7e3-d804-ebda-4571-0e5823b3d00d.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Yamatake.Helper](../html/86528ff3-8972-57d5-e768-c195b2888ccd.htm "HslCommunication.Profinet.Yamatake.Helper")

[DigitronCPLHelper 类](../html/fa8a8b36-d2ae-5f51-1721-64f5163ade4e.htm "DigitronCPLHelper 类")

[DigitronCPLHelper 方法](../html/48064c28-7c39-a2bf-ad5e-b8f29d0f9dd1.htm "DigitronCPLHelper 方法")

[BuildReadCommand 方法](../html/b4d41cf6-2776-3844-9da8-5593814f79c5.htm "BuildReadCommand 方法 ")

[BuildWriteCommand 方法](../html/41a2e62a-7d9d-7b58-ee00-8f73ad538114.htm "BuildWriteCommand 方法 ")

[ExtraActualResponse 方法](../html/1806d7e3-d804-ebda-4571-0e5823b3d00d.htm "ExtraActualResponse 方法 ")

[GetErrorText 方法](../html/abe4aa67-d36b-4f18-3d45-96142d78120b.htm "GetErrorText 方法 ")

[PackResponseContent 方法](../html/04df21d0-d08d-24ab-02db-0a1f259581a3.htm "PackResponseContent 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DigitronCPLHelperExtraActualResponse 方法 |

从反馈的数据内容中解析出真实的数据信息

**命名空间：**
 [HslCommunication.Profinet.Yamatake.Helper](86528ff3-8972-57d5-e768-c195b2888ccd.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult<byte[]> ExtraActualResponse(
	byte[] response
)
```

```
Public Shared Function ExtraActualResponse ( 
	response As Byte()
) As OperateResult(Of Byte())
```

```
public:
static OperateResult<array<unsigned char>^>^ ExtraActualResponse(
	array<unsigned char>^ response
)
```

```
static member ExtraActualResponse : 
        response : byte[] -> OperateResult<byte[]> 
```

#### 参数

response
:   类型：SystemByte  
    仪表反馈的真实的数据信息

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
解析之后的实际数据信息

![](../icons/SectionExpanded.png)参见

#### 引用

[DigitronCPLHelper 类](fa8a8b36-d2ae-5f51-1721-64f5163ade4e.htm)

[HslCommunication.Profinet.Yamatake.Helper 命名空间](86528ff3-8972-57d5-e768-c195b2888ccd.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GetErrorText 方法 

[原文連結](http://api.hslcommunication.cn/html/abe4aa67-d36b-4f18-3d45-96142d78120b.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Yamatake.Helper](../html/86528ff3-8972-57d5-e768-c195b2888ccd.htm "HslCommunication.Profinet.Yamatake.Helper")

[DigitronCPLHelper 类](../html/fa8a8b36-d2ae-5f51-1721-64f5163ade4e.htm "DigitronCPLHelper 类")

[DigitronCPLHelper 方法](../html/48064c28-7c39-a2bf-ad5e-b8f29d0f9dd1.htm "DigitronCPLHelper 方法")

[BuildReadCommand 方法](../html/b4d41cf6-2776-3844-9da8-5593814f79c5.htm "BuildReadCommand 方法 ")

[BuildWriteCommand 方法](../html/41a2e62a-7d9d-7b58-ee00-8f73ad538114.htm "BuildWriteCommand 方法 ")

[ExtraActualResponse 方法](../html/1806d7e3-d804-ebda-4571-0e5823b3d00d.htm "ExtraActualResponse 方法 ")

[GetErrorText 方法](../html/abe4aa67-d36b-4f18-3d45-96142d78120b.htm "GetErrorText 方法 ")

[PackResponseContent 方法](../html/04df21d0-d08d-24ab-02db-0a1f259581a3.htm "PackResponseContent 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DigitronCPLHelperGetErrorText 方法 |

根据错误码获取到相关的错误代号信息

**命名空间：**
 [HslCommunication.Profinet.Yamatake.Helper](86528ff3-8972-57d5-e768-c195b2888ccd.htm)  
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
    错误码

#### 返回值

类型：String  
错误码对应的文本描述信息

![](../icons/SectionExpanded.png)参见

#### 引用

[DigitronCPLHelper 类](fa8a8b36-d2ae-5f51-1721-64f5163ade4e.htm)

[HslCommunication.Profinet.Yamatake.Helper 命名空间](86528ff3-8972-57d5-e768-c195b2888ccd.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## PackResponseContent 方法 

[原文連結](http://api.hslcommunication.cn/html/04df21d0-d08d-24ab-02db-0a1f259581a3.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Yamatake.Helper](../html/86528ff3-8972-57d5-e768-c195b2888ccd.htm "HslCommunication.Profinet.Yamatake.Helper")

[DigitronCPLHelper 类](../html/fa8a8b36-d2ae-5f51-1721-64f5163ade4e.htm "DigitronCPLHelper 类")

[DigitronCPLHelper 方法](../html/48064c28-7c39-a2bf-ad5e-b8f29d0f9dd1.htm "DigitronCPLHelper 方法")

[BuildReadCommand 方法](../html/b4d41cf6-2776-3844-9da8-5593814f79c5.htm "BuildReadCommand 方法 ")

[BuildWriteCommand 方法](../html/41a2e62a-7d9d-7b58-ee00-8f73ad538114.htm "BuildWriteCommand 方法 ")

[ExtraActualResponse 方法](../html/1806d7e3-d804-ebda-4571-0e5823b3d00d.htm "ExtraActualResponse 方法 ")

[GetErrorText 方法](../html/abe4aa67-d36b-4f18-3d45-96142d78120b.htm "GetErrorText 方法 ")

[PackResponseContent 方法](../html/04df21d0-d08d-24ab-02db-0a1f259581a3.htm "PackResponseContent 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DigitronCPLHelperPackResponseContent 方法 |

用于服务器反馈的数据的报文打包操作

**命名空间：**
 [HslCommunication.Profinet.Yamatake.Helper](86528ff3-8972-57d5-e768-c195b2888ccd.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static byte[] PackResponseContent(
	byte station,
	int err,
	byte[] value,
	byte dataType
)
```

```
Public Shared Function PackResponseContent ( 
	station As Byte,
	err As Integer,
	value As Byte(),
	dataType As Byte
) As Byte()
```

```
public:
static array<unsigned char>^ PackResponseContent(
	unsigned char station, 
	int err, 
	array<unsigned char>^ value, 
	unsigned char dataType
)
```

```
static member PackResponseContent : 
        station : byte * 
        err : int * 
        value : byte[] * 
        dataType : byte -> byte[] 
```

#### 参数

station
:   类型：SystemByte  
    站号

err
:   类型：SystemInt32  
    错误码，如果为0则表示正常

value
:   类型：SystemByte  
    原始数据值信息

dataType
:   类型：SystemByte  
    数据类型

#### 返回值

类型：Byte  
打包的报文数据信息

![](../icons/SectionExpanded.png)参见

#### 引用

[DigitronCPLHelper 类](fa8a8b36-d2ae-5f51-1721-64f5163ade4e.htm)

[HslCommunication.Profinet.Yamatake.Helper 命名空间](86528ff3-8972-57d5-e768-c195b2888ccd.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)