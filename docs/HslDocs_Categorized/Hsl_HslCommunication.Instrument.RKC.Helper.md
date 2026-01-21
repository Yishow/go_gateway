# HslCommunication - HslCommunication.Instrument.RKC.Helper

> 分類頁數: 10



---
## HslCommunication.Instrument.RKC.Helper

[原文連結](http://api.hslcommunication.cn/html/2cdad1f4-ee6a-ad8b-6007-ad356b31e63b.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.RKC.Helper](../html/2cdad1f4-ee6a-ad8b-6007-ad356b31e63b.htm "HslCommunication.Instrument.RKC.Helper")

[TemperatureControllerHelper 类](../html/c22605f7-e27d-4eb4-08a6-280ba90b5c06.htm "TemperatureControllerHelper 类")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslCommunication.Instrument.RKC.Helper 命名空间 |

[缺少 "N:HslCommunication.Instrument.RKC.Helper" 的 <summary> 文档]

![](../icons/SectionExpanded.png)类

|  | 类 | 说明 |
| --- | --- | --- |
| 公共类 | [TemperatureControllerHelper](c22605f7-e27d-4eb4-08a6-280ba90b5c06.htm) | RKC温度控制器的辅助类信息，提供了报文的生成，读写的实现的方法  Auxiliary information of the RKC temperature controller provides a method for message generation and reading and writing |

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## TemperatureControllerHelper 类

[原文連結](http://api.hslcommunication.cn/html/c22605f7-e27d-4eb4-08a6-280ba90b5c06.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.RKC.Helper](../html/2cdad1f4-ee6a-ad8b-6007-ad356b31e63b.htm "HslCommunication.Instrument.RKC.Helper")

[TemperatureControllerHelper 类](../html/c22605f7-e27d-4eb4-08a6-280ba90b5c06.htm "TemperatureControllerHelper 类")

[TemperatureControllerHelper 构造函数](../html/5f1d7ede-7805-a0f4-500f-cd0383c33932.htm "TemperatureControllerHelper 构造函数 ")

[TemperatureControllerHelper 方法](../html/1daf4895-ebdc-0ba3-30d1-fed7d20f7241.htm "TemperatureControllerHelper 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| TemperatureControllerHelper 类 |

RKC温度控制器的辅助类信息，提供了报文的生成，读写的实现的方法  
Auxiliary information of the RKC temperature controller provides a method for message generation and reading and writing

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  HslCommunication.Instrument.RKC.HelperTemperatureControllerHelper

**命名空间：**
 [HslCommunication.Instrument.RKC.Helper](2cdad1f4-ee6a-ad8b-6007-ad356b31e63b.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class TemperatureControllerHelper
```

```
Public Class TemperatureControllerHelper
```

```
public ref class TemperatureControllerHelper
```

```
type TemperatureControllerHelper =  class end
```

TemperatureControllerHelper 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [TemperatureControllerHelper](5f1d7ede-7805-a0f4-500f-cd0383c33932.htm) | 初始化 TemperatureControllerHelper 类的一个新实例 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法静态成员 | [BuildReadCommand](19eaab23-8b3c-1692-cf9c-4130906deac8.htm) | 构建读取的报文命令，需要指定站号信息，数据地址 |
| 公共方法静态成员 | [BuildWriteCommand](895624be-f7b3-8ad0-6826-d4e86f056e4e.htm) | 构建一个写入的报文信息 |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法静态成员 | [ReadDouble](4c2d2ecc-855c-c8ef-8b3d-36d88be820d9.htm) | 从RKC设备读取Double类型的数据信息，地址示例：M1,M2,M3,AA,AB,B1,ER等，更详细的地址及具体含义需要参考API文档  Read Double type data information from RKC device. Examples of addresses: M1, M2, M3, AA, AB, B1, ER, etc. For more detailed addresses and specific meanings, please refer to the API documentation |
| 公共方法静态成员 | [ReadDoubleAsync](3342bb17-8a4b-0f24-8f61-8be463e0b501.htm) | 从RKC设备读取Double类型的数据信息，地址示例：M1,M2,M3,AA,AB,B1,ER等，更详细的地址及具体含义需要参考API文档  Read Double type data information from RKC device. Examples of addresses: M1, M2, M3, AA, AB, B1, ER, etc. For more detailed addresses and specific meanings, please refer to the API documentation |
| 公共方法 | ToString | (继承自 Object。) |
| 公共方法静态成员 | [Write](2066c275-d854-23ee-0b2f-e615e63122db.htm) | 将Double类型的数据写入到RKC设备中去，地址示例：M1,M2,M3,AA,AB,B1,ER等，更详细的地址及具体含义需要参考API文档  Write Double type data to the RKC device. Examples of addresses: M1, M2, M3, AA, AB, B1, ER, etc. For more detailed addresses and specific meanings, please refer to the API documentation |
| 公共方法静态成员 | [WriteAsync](518286a5-1b2f-cd3e-0d20-0757d8c7d69a.htm) | 将Double类型的数据写入到RKC设备中去，地址示例：M1,M2,M3,AA,AB,B1,ER等，更详细的地址及具体含义需要参考API文档  Write Double type data to the RKC device. Examples of addresses: M1, M2, M3, AA, AB, B1, ER, etc. For more detailed addresses and specific meanings, please refer to the API documentation |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.Instrument.RKC.Helper 命名空间](2cdad1f4-ee6a-ad8b-6007-ad356b31e63b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## TemperatureControllerHelper 构造函数 

[原文連結](http://api.hslcommunication.cn/html/5f1d7ede-7805-a0f4-500f-cd0383c33932.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.RKC.Helper](../html/2cdad1f4-ee6a-ad8b-6007-ad356b31e63b.htm "HslCommunication.Instrument.RKC.Helper")

[TemperatureControllerHelper 类](../html/c22605f7-e27d-4eb4-08a6-280ba90b5c06.htm "TemperatureControllerHelper 类")

[TemperatureControllerHelper 构造函数](../html/5f1d7ede-7805-a0f4-500f-cd0383c33932.htm "TemperatureControllerHelper 构造函数 ")

[TemperatureControllerHelper 方法](../html/1daf4895-ebdc-0ba3-30d1-fed7d20f7241.htm "TemperatureControllerHelper 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| TemperatureControllerHelper 构造函数 |

初始化 [TemperatureControllerHelper](c22605f7-e27d-4eb4-08a6-280ba90b5c06.htm) 类的一个新实例

**命名空间：**
 [HslCommunication.Instrument.RKC.Helper](2cdad1f4-ee6a-ad8b-6007-ad356b31e63b.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public TemperatureControllerHelper()
```

```
Public Sub New
```

```
public:
TemperatureControllerHelper()
```

```
new : unit -> TemperatureControllerHelper
```

![](../icons/SectionExpanded.png)参见

#### 引用

[TemperatureControllerHelper 类](c22605f7-e27d-4eb4-08a6-280ba90b5c06.htm)

[HslCommunication.Instrument.RKC.Helper 命名空间](2cdad1f4-ee6a-ad8b-6007-ad356b31e63b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## TemperatureControllerHelper 方法

[原文連結](http://api.hslcommunication.cn/html/1daf4895-ebdc-0ba3-30d1-fed7d20f7241.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.RKC.Helper](../html/2cdad1f4-ee6a-ad8b-6007-ad356b31e63b.htm "HslCommunication.Instrument.RKC.Helper")

[TemperatureControllerHelper 类](../html/c22605f7-e27d-4eb4-08a6-280ba90b5c06.htm "TemperatureControllerHelper 类")

[TemperatureControllerHelper 方法](../html/1daf4895-ebdc-0ba3-30d1-fed7d20f7241.htm "TemperatureControllerHelper 方法")

[BuildReadCommand 方法](../html/19eaab23-8b3c-1692-cf9c-4130906deac8.htm "BuildReadCommand 方法 ")

[BuildWriteCommand 方法](../html/895624be-f7b3-8ad0-6826-d4e86f056e4e.htm "BuildWriteCommand 方法 ")

[ReadDouble 方法](../html/4c2d2ecc-855c-c8ef-8b3d-36d88be820d9.htm "ReadDouble 方法 ")

[ReadDoubleAsync 方法](../html/3342bb17-8a4b-0f24-8f61-8be463e0b501.htm "ReadDoubleAsync 方法 ")

[Write 方法](../html/2066c275-d854-23ee-0b2f-e615e63122db.htm "Write 方法 ")

[WriteAsync 方法](../html/518286a5-1b2f-cd3e-0d20-0757d8c7d69a.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| TemperatureControllerHelper 方法 |

[TemperatureControllerHelper](c22605f7-e27d-4eb4-08a6-280ba90b5c06.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法静态成员 | [BuildReadCommand](19eaab23-8b3c-1692-cf9c-4130906deac8.htm) | 构建读取的报文命令，需要指定站号信息，数据地址 |
| 公共方法静态成员 | [BuildWriteCommand](895624be-f7b3-8ad0-6826-d4e86f056e4e.htm) | 构建一个写入的报文信息 |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法静态成员 | [ReadDouble](4c2d2ecc-855c-c8ef-8b3d-36d88be820d9.htm) | 从RKC设备读取Double类型的数据信息，地址示例：M1,M2,M3,AA,AB,B1,ER等，更详细的地址及具体含义需要参考API文档  Read Double type data information from RKC device. Examples of addresses: M1, M2, M3, AA, AB, B1, ER, etc. For more detailed addresses and specific meanings, please refer to the API documentation |
| 公共方法静态成员 | [ReadDoubleAsync](3342bb17-8a4b-0f24-8f61-8be463e0b501.htm) | 从RKC设备读取Double类型的数据信息，地址示例：M1,M2,M3,AA,AB,B1,ER等，更详细的地址及具体含义需要参考API文档  Read Double type data information from RKC device. Examples of addresses: M1, M2, M3, AA, AB, B1, ER, etc. For more detailed addresses and specific meanings, please refer to the API documentation |
| 公共方法 | ToString | (继承自 Object。) |
| 公共方法静态成员 | [Write](2066c275-d854-23ee-0b2f-e615e63122db.htm) | 将Double类型的数据写入到RKC设备中去，地址示例：M1,M2,M3,AA,AB,B1,ER等，更详细的地址及具体含义需要参考API文档  Write Double type data to the RKC device. Examples of addresses: M1, M2, M3, AA, AB, B1, ER, etc. For more detailed addresses and specific meanings, please refer to the API documentation |
| 公共方法静态成员 | [WriteAsync](518286a5-1b2f-cd3e-0d20-0757d8c7d69a.htm) | 将Double类型的数据写入到RKC设备中去，地址示例：M1,M2,M3,AA,AB,B1,ER等，更详细的地址及具体含义需要参考API文档  Write Double type data to the RKC device. Examples of addresses: M1, M2, M3, AA, AB, B1, ER, etc. For more detailed addresses and specific meanings, please refer to the API documentation |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[TemperatureControllerHelper 类](c22605f7-e27d-4eb4-08a6-280ba90b5c06.htm)

[HslCommunication.Instrument.RKC.Helper 命名空间](2cdad1f4-ee6a-ad8b-6007-ad356b31e63b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## BuildReadCommand 方法 

[原文連結](http://api.hslcommunication.cn/html/19eaab23-8b3c-1692-cf9c-4130906deac8.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.RKC.Helper](../html/2cdad1f4-ee6a-ad8b-6007-ad356b31e63b.htm "HslCommunication.Instrument.RKC.Helper")

[TemperatureControllerHelper 类](../html/c22605f7-e27d-4eb4-08a6-280ba90b5c06.htm "TemperatureControllerHelper 类")

[TemperatureControllerHelper 方法](../html/1daf4895-ebdc-0ba3-30d1-fed7d20f7241.htm "TemperatureControllerHelper 方法")

[BuildReadCommand 方法](../html/19eaab23-8b3c-1692-cf9c-4130906deac8.htm "BuildReadCommand 方法 ")

[BuildWriteCommand 方法](../html/895624be-f7b3-8ad0-6826-d4e86f056e4e.htm "BuildWriteCommand 方法 ")

[ReadDouble 方法](../html/4c2d2ecc-855c-c8ef-8b3d-36d88be820d9.htm "ReadDouble 方法 ")

[ReadDoubleAsync 方法](../html/3342bb17-8a4b-0f24-8f61-8be463e0b501.htm "ReadDoubleAsync 方法 ")

[Write 方法](../html/2066c275-d854-23ee-0b2f-e615e63122db.htm "Write 方法 ")

[WriteAsync 方法](../html/518286a5-1b2f-cd3e-0d20-0757d8c7d69a.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| TemperatureControllerHelperBuildReadCommand 方法 |

构建读取的报文命令，需要指定站号信息，数据地址

**命名空间：**
 [HslCommunication.Instrument.RKC.Helper](2cdad1f4-ee6a-ad8b-6007-ad356b31e63b.htm)  
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
	string address
)
```

```
Public Shared Function BuildReadCommand ( 
	station As Byte,
	address As String
) As OperateResult(Of Byte())
```

```
public:
static OperateResult<array<unsigned char>^>^ BuildReadCommand(
	unsigned char station, 
	String^ address
)
```

```
static member BuildReadCommand : 
        station : byte * 
        address : string -> OperateResult<byte[]> 
```

#### 参数

station
:   类型：SystemByte  
    站号信息

address
:   类型：SystemString  
    数据的地址

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
是否成功

![](../icons/SectionExpanded.png)参见

#### 引用

[TemperatureControllerHelper 类](c22605f7-e27d-4eb4-08a6-280ba90b5c06.htm)

[HslCommunication.Instrument.RKC.Helper 命名空间](2cdad1f4-ee6a-ad8b-6007-ad356b31e63b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## BuildWriteCommand 方法 

[原文連結](http://api.hslcommunication.cn/html/895624be-f7b3-8ad0-6826-d4e86f056e4e.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.RKC.Helper](../html/2cdad1f4-ee6a-ad8b-6007-ad356b31e63b.htm "HslCommunication.Instrument.RKC.Helper")

[TemperatureControllerHelper 类](../html/c22605f7-e27d-4eb4-08a6-280ba90b5c06.htm "TemperatureControllerHelper 类")

[TemperatureControllerHelper 方法](../html/1daf4895-ebdc-0ba3-30d1-fed7d20f7241.htm "TemperatureControllerHelper 方法")

[BuildReadCommand 方法](../html/19eaab23-8b3c-1692-cf9c-4130906deac8.htm "BuildReadCommand 方法 ")

[BuildWriteCommand 方法](../html/895624be-f7b3-8ad0-6826-d4e86f056e4e.htm "BuildWriteCommand 方法 ")

[ReadDouble 方法](../html/4c2d2ecc-855c-c8ef-8b3d-36d88be820d9.htm "ReadDouble 方法 ")

[ReadDoubleAsync 方法](../html/3342bb17-8a4b-0f24-8f61-8be463e0b501.htm "ReadDoubleAsync 方法 ")

[Write 方法](../html/2066c275-d854-23ee-0b2f-e615e63122db.htm "Write 方法 ")

[WriteAsync 方法](../html/518286a5-1b2f-cd3e-0d20-0757d8c7d69a.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| TemperatureControllerHelperBuildWriteCommand 方法 |

构建一个写入的报文信息

**命名空间：**
 [HslCommunication.Instrument.RKC.Helper](2cdad1f4-ee6a-ad8b-6007-ad356b31e63b.htm)  
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
	double value
)
```

```
Public Shared Function BuildWriteCommand ( 
	station As Byte,
	address As String,
	value As Double
) As OperateResult(Of Byte())
```

```
public:
static OperateResult<array<unsigned char>^>^ BuildWriteCommand(
	unsigned char station, 
	String^ address, 
	double value
)
```

```
static member BuildWriteCommand : 
        station : byte * 
        address : string * 
        value : float -> OperateResult<byte[]> 
```

#### 参数

station
:   类型：SystemByte  
    站号信息

address
:   类型：SystemString  
    地址信息

value
:   类型：SystemDouble  
    等待写入的值

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
是否成功的结果报文

![](../icons/SectionExpanded.png)参见

#### 引用

[TemperatureControllerHelper 类](c22605f7-e27d-4eb4-08a6-280ba90b5c06.htm)

[HslCommunication.Instrument.RKC.Helper 命名空间](2cdad1f4-ee6a-ad8b-6007-ad356b31e63b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadDouble 方法 

[原文連結](http://api.hslcommunication.cn/html/4c2d2ecc-855c-c8ef-8b3d-36d88be820d9.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.RKC.Helper](../html/2cdad1f4-ee6a-ad8b-6007-ad356b31e63b.htm "HslCommunication.Instrument.RKC.Helper")

[TemperatureControllerHelper 类](../html/c22605f7-e27d-4eb4-08a6-280ba90b5c06.htm "TemperatureControllerHelper 类")

[TemperatureControllerHelper 方法](../html/1daf4895-ebdc-0ba3-30d1-fed7d20f7241.htm "TemperatureControllerHelper 方法")

[BuildReadCommand 方法](../html/19eaab23-8b3c-1692-cf9c-4130906deac8.htm "BuildReadCommand 方法 ")

[BuildWriteCommand 方法](../html/895624be-f7b3-8ad0-6826-d4e86f056e4e.htm "BuildWriteCommand 方法 ")

[ReadDouble 方法](../html/4c2d2ecc-855c-c8ef-8b3d-36d88be820d9.htm "ReadDouble 方法 ")

[ReadDoubleAsync 方法](../html/3342bb17-8a4b-0f24-8f61-8be463e0b501.htm "ReadDoubleAsync 方法 ")

[Write 方法](../html/2066c275-d854-23ee-0b2f-e615e63122db.htm "Write 方法 ")

[WriteAsync 方法](../html/518286a5-1b2f-cd3e-0d20-0757d8c7d69a.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| TemperatureControllerHelperReadDouble 方法 |

从RKC设备读取Double类型的数据信息，地址示例：M1,M2,M3,AA,AB,B1,ER等，更详细的地址及具体含义需要参考API文档  
Read Double type data information from RKC device. Examples of addresses: M1, M2, M3, AA, AB, B1, ER, etc.
For more detailed addresses and specific meanings, please refer to the API documentation

**命名空间：**
 [HslCommunication.Instrument.RKC.Helper](2cdad1f4-ee6a-ad8b-6007-ad356b31e63b.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult<double> ReadDouble(
	IReadWriteDevice device,
	byte station,
	string address
)
```

```
Public Shared Function ReadDouble ( 
	device As IReadWriteDevice,
	station As Byte,
	address As String
) As OperateResult(Of Double)
```

```
public:
static OperateResult<double>^ ReadDouble(
	IReadWriteDevice^ device, 
	unsigned char station, 
	String^ address
)
```

```
static member ReadDouble : 
        device : IReadWriteDevice * 
        station : byte * 
        address : string -> OperateResult<float> 
```

#### 参数

device
:   类型：[HslCommunication.CoreIReadWriteDevice](d987afcc-f99c-6aed-3847-5b1890504635.htm)  
    设备通信对象

station
:   类型：SystemByte  
    表号信息，也叫站号信息

address
:   类型：SystemString  
    数据地址信息，地址示例：M1,M2,M3,AA,AB,B1,ER等

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Double  
结果数据对象信息

![](../icons/SectionExpanded.png)参见

#### 引用

[TemperatureControllerHelper 类](c22605f7-e27d-4eb4-08a6-280ba90b5c06.htm)

[HslCommunication.Instrument.RKC.Helper 命名空间](2cdad1f4-ee6a-ad8b-6007-ad356b31e63b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadDoubleAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/3342bb17-8a4b-0f24-8f61-8be463e0b501.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.RKC.Helper](../html/2cdad1f4-ee6a-ad8b-6007-ad356b31e63b.htm "HslCommunication.Instrument.RKC.Helper")

[TemperatureControllerHelper 类](../html/c22605f7-e27d-4eb4-08a6-280ba90b5c06.htm "TemperatureControllerHelper 类")

[TemperatureControllerHelper 方法](../html/1daf4895-ebdc-0ba3-30d1-fed7d20f7241.htm "TemperatureControllerHelper 方法")

[BuildReadCommand 方法](../html/19eaab23-8b3c-1692-cf9c-4130906deac8.htm "BuildReadCommand 方法 ")

[BuildWriteCommand 方法](../html/895624be-f7b3-8ad0-6826-d4e86f056e4e.htm "BuildWriteCommand 方法 ")

[ReadDouble 方法](../html/4c2d2ecc-855c-c8ef-8b3d-36d88be820d9.htm "ReadDouble 方法 ")

[ReadDoubleAsync 方法](../html/3342bb17-8a4b-0f24-8f61-8be463e0b501.htm "ReadDoubleAsync 方法 ")

[Write 方法](../html/2066c275-d854-23ee-0b2f-e615e63122db.htm "Write 方法 ")

[WriteAsync 方法](../html/518286a5-1b2f-cd3e-0d20-0757d8c7d69a.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| TemperatureControllerHelperReadDoubleAsync 方法 |

从RKC设备读取Double类型的数据信息，地址示例：M1,M2,M3,AA,AB,B1,ER等，更详细的地址及具体含义需要参考API文档  
Read Double type data information from RKC device. Examples of addresses: M1, M2, M3, AA, AB, B1, ER, etc.
For more detailed addresses and specific meanings, please refer to the API documentation

**命名空间：**
 [HslCommunication.Instrument.RKC.Helper](2cdad1f4-ee6a-ad8b-6007-ad356b31e63b.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static Task<OperateResult<double>> ReadDoubleAsync(
	IReadWriteDevice device,
	byte station,
	string address
)
```

```
Public Shared Function ReadDoubleAsync ( 
	device As IReadWriteDevice,
	station As Byte,
	address As String
) As Task(Of OperateResult(Of Double))
```

```
public:
static Task<OperateResult<double>^>^ ReadDoubleAsync(
	IReadWriteDevice^ device, 
	unsigned char station, 
	String^ address
)
```

```
static member ReadDoubleAsync : 
        device : IReadWriteDevice * 
        station : byte * 
        address : string -> Task<OperateResult<float>> 
```

#### 参数

device
:   类型：[HslCommunication.CoreIReadWriteDevice](d987afcc-f99c-6aed-3847-5b1890504635.htm)  
    设备通信对象

station
:   类型：SystemByte  
    表号信息，也叫站号信息

address
:   类型：SystemString  
    数据地址信息，地址示例：M1,M2,M3,AA,AB,B1,ER等

#### 返回值

类型：Task[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Double  
结果数据对象信息

![](../icons/SectionExpanded.png)参见

#### 引用

[TemperatureControllerHelper 类](c22605f7-e27d-4eb4-08a6-280ba90b5c06.htm)

[HslCommunication.Instrument.RKC.Helper 命名空间](2cdad1f4-ee6a-ad8b-6007-ad356b31e63b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Write 方法 

[原文連結](http://api.hslcommunication.cn/html/2066c275-d854-23ee-0b2f-e615e63122db.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.RKC.Helper](../html/2cdad1f4-ee6a-ad8b-6007-ad356b31e63b.htm "HslCommunication.Instrument.RKC.Helper")

[TemperatureControllerHelper 类](../html/c22605f7-e27d-4eb4-08a6-280ba90b5c06.htm "TemperatureControllerHelper 类")

[TemperatureControllerHelper 方法](../html/1daf4895-ebdc-0ba3-30d1-fed7d20f7241.htm "TemperatureControllerHelper 方法")

[BuildReadCommand 方法](../html/19eaab23-8b3c-1692-cf9c-4130906deac8.htm "BuildReadCommand 方法 ")

[BuildWriteCommand 方法](../html/895624be-f7b3-8ad0-6826-d4e86f056e4e.htm "BuildWriteCommand 方法 ")

[ReadDouble 方法](../html/4c2d2ecc-855c-c8ef-8b3d-36d88be820d9.htm "ReadDouble 方法 ")

[ReadDoubleAsync 方法](../html/3342bb17-8a4b-0f24-8f61-8be463e0b501.htm "ReadDoubleAsync 方法 ")

[Write 方法](../html/2066c275-d854-23ee-0b2f-e615e63122db.htm "Write 方法 ")

[WriteAsync 方法](../html/518286a5-1b2f-cd3e-0d20-0757d8c7d69a.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| TemperatureControllerHelperWrite 方法 |

将Double类型的数据写入到RKC设备中去，地址示例：M1,M2,M3,AA,AB,B1,ER等，更详细的地址及具体含义需要参考API文档  
Write Double type data to the RKC device. Examples of addresses: M1, M2, M3, AA, AB, B1, ER, etc.
For more detailed addresses and specific meanings, please refer to the API documentation

**命名空间：**
 [HslCommunication.Instrument.RKC.Helper](2cdad1f4-ee6a-ad8b-6007-ad356b31e63b.htm)  
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
	IReadWriteDevice device,
	byte station,
	string address,
	double value
)
```

```
Public Shared Function Write ( 
	device As IReadWriteDevice,
	station As Byte,
	address As String,
	value As Double
) As OperateResult
```

```
public:
static OperateResult^ Write(
	IReadWriteDevice^ device, 
	unsigned char station, 
	String^ address, 
	double value
)
```

```
static member Write : 
        device : IReadWriteDevice * 
        station : byte * 
        address : string * 
        value : float -> OperateResult 
```

#### 参数

device
:   类型：[HslCommunication.CoreIReadWriteDevice](d987afcc-f99c-6aed-3847-5b1890504635.htm)  
    设备通信对象

station
:   类型：SystemByte  
    表号信息，也叫站号信息

address
:   类型：SystemString  
    数据的地址信息，地址示例：M1,M2,M3,AA,AB,B1,ER等

value
:   类型：SystemDouble  
    等待写入的值

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否写入成功

![](../icons/SectionExpanded.png)参见

#### 引用

[TemperatureControllerHelper 类](c22605f7-e27d-4eb4-08a6-280ba90b5c06.htm)

[HslCommunication.Instrument.RKC.Helper 命名空间](2cdad1f4-ee6a-ad8b-6007-ad356b31e63b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## WriteAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/518286a5-1b2f-cd3e-0d20-0757d8c7d69a.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.RKC.Helper](../html/2cdad1f4-ee6a-ad8b-6007-ad356b31e63b.htm "HslCommunication.Instrument.RKC.Helper")

[TemperatureControllerHelper 类](../html/c22605f7-e27d-4eb4-08a6-280ba90b5c06.htm "TemperatureControllerHelper 类")

[TemperatureControllerHelper 方法](../html/1daf4895-ebdc-0ba3-30d1-fed7d20f7241.htm "TemperatureControllerHelper 方法")

[BuildReadCommand 方法](../html/19eaab23-8b3c-1692-cf9c-4130906deac8.htm "BuildReadCommand 方法 ")

[BuildWriteCommand 方法](../html/895624be-f7b3-8ad0-6826-d4e86f056e4e.htm "BuildWriteCommand 方法 ")

[ReadDouble 方法](../html/4c2d2ecc-855c-c8ef-8b3d-36d88be820d9.htm "ReadDouble 方法 ")

[ReadDoubleAsync 方法](../html/3342bb17-8a4b-0f24-8f61-8be463e0b501.htm "ReadDoubleAsync 方法 ")

[Write 方法](../html/2066c275-d854-23ee-0b2f-e615e63122db.htm "Write 方法 ")

[WriteAsync 方法](../html/518286a5-1b2f-cd3e-0d20-0757d8c7d69a.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| TemperatureControllerHelperWriteAsync 方法 |

将Double类型的数据写入到RKC设备中去，地址示例：M1,M2,M3,AA,AB,B1,ER等，更详细的地址及具体含义需要参考API文档  
Write Double type data to the RKC device. Examples of addresses: M1, M2, M3, AA, AB, B1, ER, etc.
For more detailed addresses and specific meanings, please refer to the API documentation

**命名空间：**
 [HslCommunication.Instrument.RKC.Helper](2cdad1f4-ee6a-ad8b-6007-ad356b31e63b.htm)  
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
	IReadWriteDevice device,
	byte station,
	string address,
	double value
)
```

```
Public Shared Function WriteAsync ( 
	device As IReadWriteDevice,
	station As Byte,
	address As String,
	value As Double
) As Task(Of OperateResult)
```

```
public:
static Task<OperateResult^>^ WriteAsync(
	IReadWriteDevice^ device, 
	unsigned char station, 
	String^ address, 
	double value
)
```

```
static member WriteAsync : 
        device : IReadWriteDevice * 
        station : byte * 
        address : string * 
        value : float -> Task<OperateResult> 
```

#### 参数

device
:   类型：[HslCommunication.CoreIReadWriteDevice](d987afcc-f99c-6aed-3847-5b1890504635.htm)  
    设备通信对象

station
:   类型：SystemByte  
    表号信息，也叫站号信息

address
:   类型：SystemString  
    数据的地址信息，地址示例：M1,M2,M3,AA,AB,B1,ER等

value
:   类型：SystemDouble  
    等待写入的值

#### 返回值

类型：Task[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否写入成功

![](../icons/SectionExpanded.png)参见

#### 引用

[TemperatureControllerHelper 类](c22605f7-e27d-4eb4-08a6-280ba90b5c06.htm)

[HslCommunication.Instrument.RKC.Helper 命名空间](2cdad1f4-ee6a-ad8b-6007-ad356b31e63b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)