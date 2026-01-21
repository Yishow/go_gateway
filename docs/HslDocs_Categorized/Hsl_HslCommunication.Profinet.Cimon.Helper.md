# HslCommunication - HslCommunication.Profinet.Cimon.Helper

> 分類頁數: 14



---
## HslCommunication.Profinet.Cimon.Helper

[原文連結](http://api.hslcommunication.cn/html/b2a94f85-ae94-43a8-6837-ab54c3ff8472.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Cimon.Helper](../html/b2a94f85-ae94-43a8-6837-ab54c3ff8472.htm "HslCommunication.Profinet.Cimon.Helper")

[CimonHelper 类](../html/905927a8-4f9b-57a8-d154-285d8392501c.htm "CimonHelper 类")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslCommunication.Profinet.Cimon.Helper 命名空间 |

[缺少 "N:HslCommunication.Profinet.Cimon.Helper" 的 <summary> 文档]

![](../icons/SectionExpanded.png)类

|  | 类 | 说明 |
| --- | --- | --- |
| 公共类 | [CimonHelper](905927a8-4f9b-57a8-d154-285d8392501c.htm) | 辅助帮助类对象 |

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## CimonHelper 类

[原文連結](http://api.hslcommunication.cn/html/905927a8-4f9b-57a8-d154-285d8392501c.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Cimon.Helper](../html/b2a94f85-ae94-43a8-6837-ab54c3ff8472.htm "HslCommunication.Profinet.Cimon.Helper")

[CimonHelper 类](../html/905927a8-4f9b-57a8-d154-285d8392501c.htm "CimonHelper 类")

[CimonHelper 构造函数](../html/0024c101-9b2e-b568-285b-1ad683d93744.htm "CimonHelper 构造函数 ")

[CimonHelper 方法](../html/064d3ebd-bc26-34b1-890e-42aba4f37907.htm "CimonHelper 方法")

[CimonHelper 字段](../html/000bf2d7-caf4-8652-2eec-2b7d9394fea8.htm "CimonHelper 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CimonHelper 类 |

辅助帮助类对象

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  HslCommunication.Profinet.Cimon.HelperCimonHelper

**命名空间：**
 [HslCommunication.Profinet.Cimon.Helper](b2a94f85-ae94-43a8-6837-ab54c3ff8472.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class CimonHelper
```

```
Public Class CimonHelper
```

```
public ref class CimonHelper
```

```
type CimonHelper =  class end
```

CimonHelper 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [CimonHelper](0024c101-9b2e-b568-285b-1ad683d93744.htm) | 初始化 CimonHelper 类的一个新实例 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法静态成员 | [BuildReadBitCommand](58556fca-1e03-0768-38db-275bdb113f8c.htm) | 构建一个读取位的报文信息 |
| 公共方法静态成员 | [BuildReadByteCommand](cdae7c83-0b68-efdc-f11f-517a53529f57.htm) | 构建读取的报文信息 |
| 公共方法静态成员 | [BuildWriteBitCommand](34f882ef-724a-dd1b-7dbc-bb1aa83f8558.htm) | 构建写入bool数组的报文信息 |
| 公共方法静态成员 | [BuildWriteByteCommand](e8406da3-061c-674a-a8ac-dbf17a2102aa.htm) | 构建写入的报文信息 |
| 公共方法 | Equals | (继承自 Object。) |
| 公共方法静态成员 | [ExtractActualData](24c9311c-bdf6-f7a2-7c37-b13661056ae9.htm) | 解析出实际的数据信息 |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法静态成员 | [GetErrorText](6d81fea1-065b-d76c-74fc-d27fd25ea5c1.htm) | 根据错误码获取错误文本信息 |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法静态成员 | [PackEntireCommand](0ec61d6f-f2c7-e826-6ece-11ba2924571c.htm) | 打包完整的命令信息 |
| 公共方法静态成员 | [PackErrorResponse](0fae543a-102f-8736-b861-126605c8ea85.htm) | 打包错误信息的报文 |
| 公共方法 | ToString | (继承自 Object。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)字段

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共字段静态成员 | [AddressTypes](b239d8c0-5228-ce83-1cd3-bb8bd355fbee.htm) | 所有支持的地址信息 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.Profinet.Cimon.Helper 命名空间](b2a94f85-ae94-43a8-6837-ab54c3ff8472.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## CimonHelper 构造函数 

[原文連結](http://api.hslcommunication.cn/html/0024c101-9b2e-b568-285b-1ad683d93744.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Cimon.Helper](../html/b2a94f85-ae94-43a8-6837-ab54c3ff8472.htm "HslCommunication.Profinet.Cimon.Helper")

[CimonHelper 类](../html/905927a8-4f9b-57a8-d154-285d8392501c.htm "CimonHelper 类")

[CimonHelper 构造函数](../html/0024c101-9b2e-b568-285b-1ad683d93744.htm "CimonHelper 构造函数 ")

[CimonHelper 方法](../html/064d3ebd-bc26-34b1-890e-42aba4f37907.htm "CimonHelper 方法")

[CimonHelper 字段](../html/000bf2d7-caf4-8652-2eec-2b7d9394fea8.htm "CimonHelper 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CimonHelper 构造函数 |

初始化 [CimonHelper](905927a8-4f9b-57a8-d154-285d8392501c.htm) 类的一个新实例

**命名空间：**
 [HslCommunication.Profinet.Cimon.Helper](b2a94f85-ae94-43a8-6837-ab54c3ff8472.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public CimonHelper()
```

```
Public Sub New
```

```
public:
CimonHelper()
```

```
new : unit -> CimonHelper
```

![](../icons/SectionExpanded.png)参见

#### 引用

[CimonHelper 类](905927a8-4f9b-57a8-d154-285d8392501c.htm)

[HslCommunication.Profinet.Cimon.Helper 命名空间](b2a94f85-ae94-43a8-6837-ab54c3ff8472.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## CimonHelper 方法

[原文連結](http://api.hslcommunication.cn/html/064d3ebd-bc26-34b1-890e-42aba4f37907.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Cimon.Helper](../html/b2a94f85-ae94-43a8-6837-ab54c3ff8472.htm "HslCommunication.Profinet.Cimon.Helper")

[CimonHelper 类](../html/905927a8-4f9b-57a8-d154-285d8392501c.htm "CimonHelper 类")

[CimonHelper 方法](../html/064d3ebd-bc26-34b1-890e-42aba4f37907.htm "CimonHelper 方法")

[BuildReadBitCommand 方法](../html/58556fca-1e03-0768-38db-275bdb113f8c.htm "BuildReadBitCommand 方法 ")

[BuildReadByteCommand 方法](../html/cdae7c83-0b68-efdc-f11f-517a53529f57.htm "BuildReadByteCommand 方法 ")

[BuildWriteBitCommand 方法](../html/34f882ef-724a-dd1b-7dbc-bb1aa83f8558.htm "BuildWriteBitCommand 方法 ")

[BuildWriteByteCommand 方法](../html/e8406da3-061c-674a-a8ac-dbf17a2102aa.htm "BuildWriteByteCommand 方法 ")

[ExtractActualData 方法](../html/24c9311c-bdf6-f7a2-7c37-b13661056ae9.htm "ExtractActualData 方法 ")

[GetErrorText 方法](../html/6d81fea1-065b-d76c-74fc-d27fd25ea5c1.htm "GetErrorText 方法 ")

[PackEntireCommand 方法](../html/0ec61d6f-f2c7-e826-6ece-11ba2924571c.htm "PackEntireCommand 方法 ")

[PackErrorResponse 方法](../html/0fae543a-102f-8736-b861-126605c8ea85.htm "PackErrorResponse 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CimonHelper 方法 |

[CimonHelper](905927a8-4f9b-57a8-d154-285d8392501c.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法静态成员 | [BuildReadBitCommand](58556fca-1e03-0768-38db-275bdb113f8c.htm) | 构建一个读取位的报文信息 |
| 公共方法静态成员 | [BuildReadByteCommand](cdae7c83-0b68-efdc-f11f-517a53529f57.htm) | 构建读取的报文信息 |
| 公共方法静态成员 | [BuildWriteBitCommand](34f882ef-724a-dd1b-7dbc-bb1aa83f8558.htm) | 构建写入bool数组的报文信息 |
| 公共方法静态成员 | [BuildWriteByteCommand](e8406da3-061c-674a-a8ac-dbf17a2102aa.htm) | 构建写入的报文信息 |
| 公共方法 | Equals | (继承自 Object。) |
| 公共方法静态成员 | [ExtractActualData](24c9311c-bdf6-f7a2-7c37-b13661056ae9.htm) | 解析出实际的数据信息 |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法静态成员 | [GetErrorText](6d81fea1-065b-d76c-74fc-d27fd25ea5c1.htm) | 根据错误码获取错误文本信息 |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法静态成员 | [PackEntireCommand](0ec61d6f-f2c7-e826-6ece-11ba2924571c.htm) | 打包完整的命令信息 |
| 公共方法静态成员 | [PackErrorResponse](0fae543a-102f-8736-b861-126605c8ea85.htm) | 打包错误信息的报文 |
| 公共方法 | ToString | (继承自 Object。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[CimonHelper 类](905927a8-4f9b-57a8-d154-285d8392501c.htm)

[HslCommunication.Profinet.Cimon.Helper 命名空间](b2a94f85-ae94-43a8-6837-ab54c3ff8472.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## BuildReadBitCommand 方法 

[原文連結](http://api.hslcommunication.cn/html/58556fca-1e03-0768-38db-275bdb113f8c.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Cimon.Helper](../html/b2a94f85-ae94-43a8-6837-ab54c3ff8472.htm "HslCommunication.Profinet.Cimon.Helper")

[CimonHelper 类](../html/905927a8-4f9b-57a8-d154-285d8392501c.htm "CimonHelper 类")

[CimonHelper 方法](../html/064d3ebd-bc26-34b1-890e-42aba4f37907.htm "CimonHelper 方法")

[BuildReadBitCommand 方法](../html/58556fca-1e03-0768-38db-275bdb113f8c.htm "BuildReadBitCommand 方法 ")

[BuildReadByteCommand 方法](../html/cdae7c83-0b68-efdc-f11f-517a53529f57.htm "BuildReadByteCommand 方法 ")

[BuildWriteBitCommand 方法](../html/34f882ef-724a-dd1b-7dbc-bb1aa83f8558.htm "BuildWriteBitCommand 方法 ")

[BuildWriteByteCommand 方法](../html/e8406da3-061c-674a-a8ac-dbf17a2102aa.htm "BuildWriteByteCommand 方法 ")

[ExtractActualData 方法](../html/24c9311c-bdf6-f7a2-7c37-b13661056ae9.htm "ExtractActualData 方法 ")

[GetErrorText 方法](../html/6d81fea1-065b-d76c-74fc-d27fd25ea5c1.htm "GetErrorText 方法 ")

[PackEntireCommand 方法](../html/0ec61d6f-f2c7-e826-6ece-11ba2924571c.htm "PackEntireCommand 方法 ")

[PackErrorResponse 方法](../html/0fae543a-102f-8736-b861-126605c8ea85.htm "PackErrorResponse 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CimonHelperBuildReadBitCommand 方法 |

构建一个读取位的报文信息

**命名空间：**
 [HslCommunication.Profinet.Cimon.Helper](b2a94f85-ae94-43a8-6837-ab54c3ff8472.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult<byte[]> BuildReadBitCommand(
	byte frameNo,
	string address,
	int length
)
```

```
Public Shared Function BuildReadBitCommand ( 
	frameNo As Byte,
	address As String,
	length As Integer
) As OperateResult(Of Byte())
```

```
public:
static OperateResult<array<unsigned char>^>^ BuildReadBitCommand(
	unsigned char frameNo, 
	String^ address, 
	int length
)
```

```
static member BuildReadBitCommand : 
        frameNo : byte * 
        address : string * 
        length : int -> OperateResult<byte[]> 
```

#### 参数

frameNo
:   类型：SystemByte  
    站号信息

address
:   类型：SystemString  

    [缺少 "M:HslCommunication.Profinet.Cimon.Helper.CimonHelper.BuildReadBitCommand(System.Byte,System.String,System.Int32)" 的 <param name="address"/> 文档]

length
:   类型：SystemInt32  

    [缺少 "M:HslCommunication.Profinet.Cimon.Helper.CimonHelper.BuildReadBitCommand(System.Byte,System.String,System.Int32)" 的 <param name="length"/> 文档]

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  

[缺少 "M:HslCommunication.Profinet.Cimon.Helper.CimonHelper.BuildReadBitCommand(System.Byte,System.String,System.Int32)" 的 <returns> 文档]

![](../icons/SectionExpanded.png)参见

#### 引用

[CimonHelper 类](905927a8-4f9b-57a8-d154-285d8392501c.htm)

[HslCommunication.Profinet.Cimon.Helper 命名空间](b2a94f85-ae94-43a8-6837-ab54c3ff8472.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## BuildReadByteCommand 方法 

[原文連結](http://api.hslcommunication.cn/html/cdae7c83-0b68-efdc-f11f-517a53529f57.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Cimon.Helper](../html/b2a94f85-ae94-43a8-6837-ab54c3ff8472.htm "HslCommunication.Profinet.Cimon.Helper")

[CimonHelper 类](../html/905927a8-4f9b-57a8-d154-285d8392501c.htm "CimonHelper 类")

[CimonHelper 方法](../html/064d3ebd-bc26-34b1-890e-42aba4f37907.htm "CimonHelper 方法")

[BuildReadBitCommand 方法](../html/58556fca-1e03-0768-38db-275bdb113f8c.htm "BuildReadBitCommand 方法 ")

[BuildReadByteCommand 方法](../html/cdae7c83-0b68-efdc-f11f-517a53529f57.htm "BuildReadByteCommand 方法 ")

[BuildWriteBitCommand 方法](../html/34f882ef-724a-dd1b-7dbc-bb1aa83f8558.htm "BuildWriteBitCommand 方法 ")

[BuildWriteByteCommand 方法](../html/e8406da3-061c-674a-a8ac-dbf17a2102aa.htm "BuildWriteByteCommand 方法 ")

[ExtractActualData 方法](../html/24c9311c-bdf6-f7a2-7c37-b13661056ae9.htm "ExtractActualData 方法 ")

[GetErrorText 方法](../html/6d81fea1-065b-d76c-74fc-d27fd25ea5c1.htm "GetErrorText 方法 ")

[PackEntireCommand 方法](../html/0ec61d6f-f2c7-e826-6ece-11ba2924571c.htm "PackEntireCommand 方法 ")

[PackErrorResponse 方法](../html/0fae543a-102f-8736-b861-126605c8ea85.htm "PackErrorResponse 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CimonHelperBuildReadByteCommand 方法 |

构建读取的报文信息

**命名空间：**
 [HslCommunication.Profinet.Cimon.Helper](b2a94f85-ae94-43a8-6837-ab54c3ff8472.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult<byte[]> BuildReadByteCommand(
	byte frameNo,
	string address,
	int length
)
```

```
Public Shared Function BuildReadByteCommand ( 
	frameNo As Byte,
	address As String,
	length As Integer
) As OperateResult(Of Byte())
```

```
public:
static OperateResult<array<unsigned char>^>^ BuildReadByteCommand(
	unsigned char frameNo, 
	String^ address, 
	int length
)
```

```
static member BuildReadByteCommand : 
        frameNo : byte * 
        address : string * 
        length : int -> OperateResult<byte[]> 
```

#### 参数

frameNo
:   类型：SystemByte  
    站号信息

address
:   类型：SystemString  
    PLC的地址

length
:   类型：SystemInt32  
    读取的长度

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
读取指定地址的报文

![](../icons/SectionExpanded.png)参见

#### 引用

[CimonHelper 类](905927a8-4f9b-57a8-d154-285d8392501c.htm)

[HslCommunication.Profinet.Cimon.Helper 命名空间](b2a94f85-ae94-43a8-6837-ab54c3ff8472.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## BuildWriteBitCommand 方法 

[原文連結](http://api.hslcommunication.cn/html/34f882ef-724a-dd1b-7dbc-bb1aa83f8558.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Cimon.Helper](../html/b2a94f85-ae94-43a8-6837-ab54c3ff8472.htm "HslCommunication.Profinet.Cimon.Helper")

[CimonHelper 类](../html/905927a8-4f9b-57a8-d154-285d8392501c.htm "CimonHelper 类")

[CimonHelper 方法](../html/064d3ebd-bc26-34b1-890e-42aba4f37907.htm "CimonHelper 方法")

[BuildReadBitCommand 方法](../html/58556fca-1e03-0768-38db-275bdb113f8c.htm "BuildReadBitCommand 方法 ")

[BuildReadByteCommand 方法](../html/cdae7c83-0b68-efdc-f11f-517a53529f57.htm "BuildReadByteCommand 方法 ")

[BuildWriteBitCommand 方法](../html/34f882ef-724a-dd1b-7dbc-bb1aa83f8558.htm "BuildWriteBitCommand 方法 ")

[BuildWriteByteCommand 方法](../html/e8406da3-061c-674a-a8ac-dbf17a2102aa.htm "BuildWriteByteCommand 方法 ")

[ExtractActualData 方法](../html/24c9311c-bdf6-f7a2-7c37-b13661056ae9.htm "ExtractActualData 方法 ")

[GetErrorText 方法](../html/6d81fea1-065b-d76c-74fc-d27fd25ea5c1.htm "GetErrorText 方法 ")

[PackEntireCommand 方法](../html/0ec61d6f-f2c7-e826-6ece-11ba2924571c.htm "PackEntireCommand 方法 ")

[PackErrorResponse 方法](../html/0fae543a-102f-8736-b861-126605c8ea85.htm "PackErrorResponse 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CimonHelperBuildWriteBitCommand 方法 |

构建写入bool数组的报文信息

**命名空间：**
 [HslCommunication.Profinet.Cimon.Helper](b2a94f85-ae94-43a8-6837-ab54c3ff8472.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult<byte[]> BuildWriteBitCommand(
	byte frameNo,
	string address,
	bool[] data
)
```

```
Public Shared Function BuildWriteBitCommand ( 
	frameNo As Byte,
	address As String,
	data As Boolean()
) As OperateResult(Of Byte())
```

```
public:
static OperateResult<array<unsigned char>^>^ BuildWriteBitCommand(
	unsigned char frameNo, 
	String^ address, 
	array<bool>^ data
)
```

```
static member BuildWriteBitCommand : 
        frameNo : byte * 
        address : string * 
        data : bool[] -> OperateResult<byte[]> 
```

#### 参数

frameNo
:   类型：SystemByte  
    站号信息

address
:   类型：SystemString  
    写入的PLC地址

data
:   类型：SystemBoolean  
    Bool数组值

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
写入数据的报文信息

![](../icons/SectionExpanded.png)参见

#### 引用

[CimonHelper 类](905927a8-4f9b-57a8-d154-285d8392501c.htm)

[HslCommunication.Profinet.Cimon.Helper 命名空间](b2a94f85-ae94-43a8-6837-ab54c3ff8472.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## BuildWriteByteCommand 方法 

[原文連結](http://api.hslcommunication.cn/html/e8406da3-061c-674a-a8ac-dbf17a2102aa.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Cimon.Helper](../html/b2a94f85-ae94-43a8-6837-ab54c3ff8472.htm "HslCommunication.Profinet.Cimon.Helper")

[CimonHelper 类](../html/905927a8-4f9b-57a8-d154-285d8392501c.htm "CimonHelper 类")

[CimonHelper 方法](../html/064d3ebd-bc26-34b1-890e-42aba4f37907.htm "CimonHelper 方法")

[BuildReadBitCommand 方法](../html/58556fca-1e03-0768-38db-275bdb113f8c.htm "BuildReadBitCommand 方法 ")

[BuildReadByteCommand 方法](../html/cdae7c83-0b68-efdc-f11f-517a53529f57.htm "BuildReadByteCommand 方法 ")

[BuildWriteBitCommand 方法](../html/34f882ef-724a-dd1b-7dbc-bb1aa83f8558.htm "BuildWriteBitCommand 方法 ")

[BuildWriteByteCommand 方法](../html/e8406da3-061c-674a-a8ac-dbf17a2102aa.htm "BuildWriteByteCommand 方法 ")

[ExtractActualData 方法](../html/24c9311c-bdf6-f7a2-7c37-b13661056ae9.htm "ExtractActualData 方法 ")

[GetErrorText 方法](../html/6d81fea1-065b-d76c-74fc-d27fd25ea5c1.htm "GetErrorText 方法 ")

[PackEntireCommand 方法](../html/0ec61d6f-f2c7-e826-6ece-11ba2924571c.htm "PackEntireCommand 方法 ")

[PackErrorResponse 方法](../html/0fae543a-102f-8736-b861-126605c8ea85.htm "PackErrorResponse 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CimonHelperBuildWriteByteCommand 方法 |

构建写入的报文信息

**命名空间：**
 [HslCommunication.Profinet.Cimon.Helper](b2a94f85-ae94-43a8-6837-ab54c3ff8472.htm)  
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
	byte frameNo,
	string address,
	byte[] data
)
```

```
Public Shared Function BuildWriteByteCommand ( 
	frameNo As Byte,
	address As String,
	data As Byte()
) As OperateResult(Of Byte())
```

```
public:
static OperateResult<array<unsigned char>^>^ BuildWriteByteCommand(
	unsigned char frameNo, 
	String^ address, 
	array<unsigned char>^ data
)
```

```
static member BuildWriteByteCommand : 
        frameNo : byte * 
        address : string * 
        data : byte[] -> OperateResult<byte[]> 
```

#### 参数

frameNo
:   类型：SystemByte  
    站号信息

address
:   类型：SystemString  
    PLC的地址

data
:   类型：SystemByte  
    写入的数据

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
写入数据到地址的报文

![](../icons/SectionExpanded.png)参见

#### 引用

[CimonHelper 类](905927a8-4f9b-57a8-d154-285d8392501c.htm)

[HslCommunication.Profinet.Cimon.Helper 命名空间](b2a94f85-ae94-43a8-6837-ab54c3ff8472.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ExtractActualData 方法 

[原文連結](http://api.hslcommunication.cn/html/24c9311c-bdf6-f7a2-7c37-b13661056ae9.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Cimon.Helper](../html/b2a94f85-ae94-43a8-6837-ab54c3ff8472.htm "HslCommunication.Profinet.Cimon.Helper")

[CimonHelper 类](../html/905927a8-4f9b-57a8-d154-285d8392501c.htm "CimonHelper 类")

[CimonHelper 方法](../html/064d3ebd-bc26-34b1-890e-42aba4f37907.htm "CimonHelper 方法")

[BuildReadBitCommand 方法](../html/58556fca-1e03-0768-38db-275bdb113f8c.htm "BuildReadBitCommand 方法 ")

[BuildReadByteCommand 方法](../html/cdae7c83-0b68-efdc-f11f-517a53529f57.htm "BuildReadByteCommand 方法 ")

[BuildWriteBitCommand 方法](../html/34f882ef-724a-dd1b-7dbc-bb1aa83f8558.htm "BuildWriteBitCommand 方法 ")

[BuildWriteByteCommand 方法](../html/e8406da3-061c-674a-a8ac-dbf17a2102aa.htm "BuildWriteByteCommand 方法 ")

[ExtractActualData 方法](../html/24c9311c-bdf6-f7a2-7c37-b13661056ae9.htm "ExtractActualData 方法 ")

[GetErrorText 方法](../html/6d81fea1-065b-d76c-74fc-d27fd25ea5c1.htm "GetErrorText 方法 ")

[PackEntireCommand 方法](../html/0ec61d6f-f2c7-e826-6ece-11ba2924571c.htm "PackEntireCommand 方法 ")

[PackErrorResponse 方法](../html/0fae543a-102f-8736-b861-126605c8ea85.htm "PackErrorResponse 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CimonHelperExtractActualData 方法 |

解析出实际的数据信息

**命名空间：**
 [HslCommunication.Profinet.Cimon.Helper](b2a94f85-ae94-43a8-6837-ab54c3ff8472.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult<byte[]> ExtractActualData(
	byte[] response
)
```

```
Public Shared Function ExtractActualData ( 
	response As Byte()
) As OperateResult(Of Byte())
```

```
public:
static OperateResult<array<unsigned char>^>^ ExtractActualData(
	array<unsigned char>^ response
)
```

```
static member ExtractActualData : 
        response : byte[] -> OperateResult<byte[]> 
```

#### 参数

response
:   类型：SystemByte  
    PLC返回的报文数据

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
提炼的真实的数据

![](../icons/SectionExpanded.png)参见

#### 引用

[CimonHelper 类](905927a8-4f9b-57a8-d154-285d8392501c.htm)

[HslCommunication.Profinet.Cimon.Helper 命名空间](b2a94f85-ae94-43a8-6837-ab54c3ff8472.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GetErrorText 方法 

[原文連結](http://api.hslcommunication.cn/html/6d81fea1-065b-d76c-74fc-d27fd25ea5c1.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Cimon.Helper](../html/b2a94f85-ae94-43a8-6837-ab54c3ff8472.htm "HslCommunication.Profinet.Cimon.Helper")

[CimonHelper 类](../html/905927a8-4f9b-57a8-d154-285d8392501c.htm "CimonHelper 类")

[CimonHelper 方法](../html/064d3ebd-bc26-34b1-890e-42aba4f37907.htm "CimonHelper 方法")

[BuildReadBitCommand 方法](../html/58556fca-1e03-0768-38db-275bdb113f8c.htm "BuildReadBitCommand 方法 ")

[BuildReadByteCommand 方法](../html/cdae7c83-0b68-efdc-f11f-517a53529f57.htm "BuildReadByteCommand 方法 ")

[BuildWriteBitCommand 方法](../html/34f882ef-724a-dd1b-7dbc-bb1aa83f8558.htm "BuildWriteBitCommand 方法 ")

[BuildWriteByteCommand 方法](../html/e8406da3-061c-674a-a8ac-dbf17a2102aa.htm "BuildWriteByteCommand 方法 ")

[ExtractActualData 方法](../html/24c9311c-bdf6-f7a2-7c37-b13661056ae9.htm "ExtractActualData 方法 ")

[GetErrorText 方法](../html/6d81fea1-065b-d76c-74fc-d27fd25ea5c1.htm "GetErrorText 方法 ")

[PackEntireCommand 方法](../html/0ec61d6f-f2c7-e826-6ece-11ba2924571c.htm "PackEntireCommand 方法 ")

[PackErrorResponse 方法](../html/0fae543a-102f-8736-b861-126605c8ea85.htm "PackErrorResponse 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CimonHelperGetErrorText 方法 |

根据错误码获取错误文本信息

**命名空间：**
 [HslCommunication.Profinet.Cimon.Helper](b2a94f85-ae94-43a8-6837-ab54c3ff8472.htm)  
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
	int error
)
```

```
Public Shared Function GetErrorText ( 
	error As Integer
) As String
```

```
public:
static String^ GetErrorText(
	int error
)
```

```
static member GetErrorText : 
        error : int -> string 
```

#### 参数

error
:   类型：SystemInt32  
    错误信息

#### 返回值

类型：String  
错误消息

![](../icons/SectionExpanded.png)参见

#### 引用

[CimonHelper 类](905927a8-4f9b-57a8-d154-285d8392501c.htm)

[HslCommunication.Profinet.Cimon.Helper 命名空间](b2a94f85-ae94-43a8-6837-ab54c3ff8472.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## PackEntireCommand 方法 

[原文連結](http://api.hslcommunication.cn/html/0ec61d6f-f2c7-e826-6ece-11ba2924571c.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Cimon.Helper](../html/b2a94f85-ae94-43a8-6837-ab54c3ff8472.htm "HslCommunication.Profinet.Cimon.Helper")

[CimonHelper 类](../html/905927a8-4f9b-57a8-d154-285d8392501c.htm "CimonHelper 类")

[CimonHelper 方法](../html/064d3ebd-bc26-34b1-890e-42aba4f37907.htm "CimonHelper 方法")

[BuildReadBitCommand 方法](../html/58556fca-1e03-0768-38db-275bdb113f8c.htm "BuildReadBitCommand 方法 ")

[BuildReadByteCommand 方法](../html/cdae7c83-0b68-efdc-f11f-517a53529f57.htm "BuildReadByteCommand 方法 ")

[BuildWriteBitCommand 方法](../html/34f882ef-724a-dd1b-7dbc-bb1aa83f8558.htm "BuildWriteBitCommand 方法 ")

[BuildWriteByteCommand 方法](../html/e8406da3-061c-674a-a8ac-dbf17a2102aa.htm "BuildWriteByteCommand 方法 ")

[ExtractActualData 方法](../html/24c9311c-bdf6-f7a2-7c37-b13661056ae9.htm "ExtractActualData 方法 ")

[GetErrorText 方法](../html/6d81fea1-065b-d76c-74fc-d27fd25ea5c1.htm "GetErrorText 方法 ")

[PackEntireCommand 方法](../html/0ec61d6f-f2c7-e826-6ece-11ba2924571c.htm "PackEntireCommand 方法 ")

[PackErrorResponse 方法](../html/0fae543a-102f-8736-b861-126605c8ea85.htm "PackErrorResponse 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CimonHelperPackEntireCommand 方法 |

打包完整的命令信息

**命名空间：**
 [HslCommunication.Profinet.Cimon.Helper](b2a94f85-ae94-43a8-6837-ab54c3ff8472.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static byte[] PackEntireCommand(
	bool response,
	byte frame,
	byte cmd,
	byte[] data
)
```

```
Public Shared Function PackEntireCommand ( 
	response As Boolean,
	frame As Byte,
	cmd As Byte,
	data As Byte()
) As Byte()
```

```
public:
static array<unsigned char>^ PackEntireCommand(
	bool response, 
	unsigned char frame, 
	unsigned char cmd, 
	array<unsigned char>^ data
)
```

```
static member PackEntireCommand : 
        response : bool * 
        frame : byte * 
        cmd : byte * 
        data : byte[] -> byte[] 
```

#### 参数

response
:   类型：SystemBoolean  
    是否响应的报文

frame
:   类型：SystemByte  
    站号信息

cmd
:   类型：SystemByte  
    命令码

data
:   类型：SystemByte  
    返回的数据信息

#### 返回值

类型：Byte  
完成的报文

![](../icons/SectionExpanded.png)参见

#### 引用

[CimonHelper 类](905927a8-4f9b-57a8-d154-285d8392501c.htm)

[HslCommunication.Profinet.Cimon.Helper 命名空间](b2a94f85-ae94-43a8-6837-ab54c3ff8472.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## PackErrorResponse 方法 

[原文連結](http://api.hslcommunication.cn/html/0fae543a-102f-8736-b861-126605c8ea85.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Cimon.Helper](../html/b2a94f85-ae94-43a8-6837-ab54c3ff8472.htm "HslCommunication.Profinet.Cimon.Helper")

[CimonHelper 类](../html/905927a8-4f9b-57a8-d154-285d8392501c.htm "CimonHelper 类")

[CimonHelper 方法](../html/064d3ebd-bc26-34b1-890e-42aba4f37907.htm "CimonHelper 方法")

[BuildReadBitCommand 方法](../html/58556fca-1e03-0768-38db-275bdb113f8c.htm "BuildReadBitCommand 方法 ")

[BuildReadByteCommand 方法](../html/cdae7c83-0b68-efdc-f11f-517a53529f57.htm "BuildReadByteCommand 方法 ")

[BuildWriteBitCommand 方法](../html/34f882ef-724a-dd1b-7dbc-bb1aa83f8558.htm "BuildWriteBitCommand 方法 ")

[BuildWriteByteCommand 方法](../html/e8406da3-061c-674a-a8ac-dbf17a2102aa.htm "BuildWriteByteCommand 方法 ")

[ExtractActualData 方法](../html/24c9311c-bdf6-f7a2-7c37-b13661056ae9.htm "ExtractActualData 方法 ")

[GetErrorText 方法](../html/6d81fea1-065b-d76c-74fc-d27fd25ea5c1.htm "GetErrorText 方法 ")

[PackEntireCommand 方法](../html/0ec61d6f-f2c7-e826-6ece-11ba2924571c.htm "PackEntireCommand 方法 ")

[PackErrorResponse 方法](../html/0fae543a-102f-8736-b861-126605c8ea85.htm "PackErrorResponse 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CimonHelperPackErrorResponse 方法 |

打包错误信息的报文

**命名空间：**
 [HslCommunication.Profinet.Cimon.Helper](b2a94f85-ae94-43a8-6837-ab54c3ff8472.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static byte[] PackErrorResponse(
	byte frame,
	int err
)
```

```
Public Shared Function PackErrorResponse ( 
	frame As Byte,
	err As Integer
) As Byte()
```

```
public:
static array<unsigned char>^ PackErrorResponse(
	unsigned char frame, 
	int err
)
```

```
static member PackErrorResponse : 
        frame : byte * 
        err : int -> byte[] 
```

#### 参数

frame
:   类型：SystemByte  

    [缺少 "M:HslCommunication.Profinet.Cimon.Helper.CimonHelper.PackErrorResponse(System.Byte,System.Int32)" 的 <param name="frame"/> 文档]

err
:   类型：SystemInt32  

    [缺少 "M:HslCommunication.Profinet.Cimon.Helper.CimonHelper.PackErrorResponse(System.Byte,System.Int32)" 的 <param name="err"/> 文档]

#### 返回值

类型：Byte  

[缺少 "M:HslCommunication.Profinet.Cimon.Helper.CimonHelper.PackErrorResponse(System.Byte,System.Int32)" 的 <returns> 文档]

![](../icons/SectionExpanded.png)参见

#### 引用

[CimonHelper 类](905927a8-4f9b-57a8-d154-285d8392501c.htm)

[HslCommunication.Profinet.Cimon.Helper 命名空间](b2a94f85-ae94-43a8-6837-ab54c3ff8472.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## CimonHelper 字段

[原文連結](http://api.hslcommunication.cn/html/000bf2d7-caf4-8652-2eec-2b7d9394fea8.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Cimon.Helper](../html/b2a94f85-ae94-43a8-6837-ab54c3ff8472.htm "HslCommunication.Profinet.Cimon.Helper")

[CimonHelper 类](../html/905927a8-4f9b-57a8-d154-285d8392501c.htm "CimonHelper 类")

[CimonHelper 字段](../html/000bf2d7-caf4-8652-2eec-2b7d9394fea8.htm "CimonHelper 字段")

[AddressTypes 字段](../html/b239d8c0-5228-ce83-1cd3-bb8bd355fbee.htm "AddressTypes 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CimonHelper 字段 |

[CimonHelper](905927a8-4f9b-57a8-d154-285d8392501c.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)字段

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共字段静态成员 | [AddressTypes](b239d8c0-5228-ce83-1cd3-bb8bd355fbee.htm) | 所有支持的地址信息 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[CimonHelper 类](905927a8-4f9b-57a8-d154-285d8392501c.htm)

[HslCommunication.Profinet.Cimon.Helper 命名空间](b2a94f85-ae94-43a8-6837-ab54c3ff8472.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## AddressTypes 字段

[原文連結](http://api.hslcommunication.cn/html/b239d8c0-5228-ce83-1cd3-bb8bd355fbee.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Cimon.Helper](../html/b2a94f85-ae94-43a8-6837-ab54c3ff8472.htm "HslCommunication.Profinet.Cimon.Helper")

[CimonHelper 类](../html/905927a8-4f9b-57a8-d154-285d8392501c.htm "CimonHelper 类")

[CimonHelper 字段](../html/000bf2d7-caf4-8652-2eec-2b7d9394fea8.htm "CimonHelper 字段")

[AddressTypes 字段](../html/b239d8c0-5228-ce83-1cd3-bb8bd355fbee.htm "AddressTypes 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CimonHelperAddressTypes 字段 |

所有支持的地址信息

**命名空间：**
 [HslCommunication.Profinet.Cimon.Helper](b2a94f85-ae94-43a8-6837-ab54c3ff8472.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public const string AddressTypes = "YMLKFTCDSX"
```

```
Public Const AddressTypes As String = "YMLKFTCDSX"
```

```
public:
literal String^ AddressTypes = "YMLKFTCDSX"
```

```
static val mutable AddressTypes: string
```

#### 字段值

类型：String

![](../icons/SectionExpanded.png)参见

#### 引用

[CimonHelper 类](905927a8-4f9b-57a8-d154-285d8392501c.htm)

[HslCommunication.Profinet.Cimon.Helper 命名空间](b2a94f85-ae94-43a8-6837-ab54c3ff8472.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)