# HslCommunication - HslCommunication.Instrument.IEC.Helper

> 分類頁數: 25



---
## HslCommunication.Instrument.IEC.Helper

[原文連結](http://api.hslcommunication.cn/html/acdb00cf-8087-ba54-8ab1-947a92f2b0b4.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.IEC.Helper](../html/acdb00cf-8087-ba54-8ab1-947a92f2b0b4.htm "HslCommunication.Instrument.IEC.Helper")

[IECHelper 类](../html/b37c3cea-b292-985a-f89b-eb97e38d44fe.htm "IECHelper 类")

[IECSessionInfo 类](../html/aa145fdf-e2f1-8eab-1375-9d9c8f5a5308.htm "IECSessionInfo 类")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslCommunication.Instrument.IEC.Helper 命名空间 |

[缺少 "N:HslCommunication.Instrument.IEC.Helper" 的 <summary> 文档]

![](../icons/SectionExpanded.png)类

|  | 类 | 说明 |
| --- | --- | --- |
| 公共类 | [IECHelper](b37c3cea-b292-985a-f89b-eb97e38d44fe.htm) | IEC协议的辅助类信息 |
| 公共类 | [IECSessionInfo](aa145fdf-e2f1-8eab-1375-9d9c8f5a5308.htm) | IEC会话的额外信息，主要包含了这个会话的收发消息ID |

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## IECHelper 类

[原文連結](http://api.hslcommunication.cn/html/b37c3cea-b292-985a-f89b-eb97e38d44fe.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.IEC.Helper](../html/acdb00cf-8087-ba54-8ab1-947a92f2b0b4.htm "HslCommunication.Instrument.IEC.Helper")

[IECHelper 类](../html/b37c3cea-b292-985a-f89b-eb97e38d44fe.htm "IECHelper 类")

[IECHelper 构造函数](../html/d127c687-70be-3195-2a23-374ad8b2ff17.htm "IECHelper 构造函数 ")

[IECHelper 方法](../html/bfa98ce5-b54e-bed4-e8cb-ec8cd035ce87.htm "IECHelper 方法")

[IECHelper 字段](../html/3c7bb720-5c4e-776a-8dba-47d8548c956f.htm "IECHelper 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IECHelper 类 |

IEC协议的辅助类信息

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  HslCommunication.Instrument.IEC.HelperIECHelper

**命名空间：**
 [HslCommunication.Instrument.IEC.Helper](acdb00cf-8087-ba54-8ab1-947a92f2b0b4.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class IECHelper
```

```
Public Class IECHelper
```

```
public ref class IECHelper
```

```
type IECHelper =  class end
```

IECHelper 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [IECHelper](d127c687-70be-3195-2a23-374ad8b2ff17.htm) | 初始化 IECHelper 类的一个新实例 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法静态成员 | [BuildFrameIMessage](e202c7d8-644d-7f16-c4f3-062e9b14752e.htm) | 构建一个I帧消息的报文信息，传入相关的参数信息，返回完整的104消息报文 |
| 公共方法静态成员 | [BuildFrameSMessage](c15c3a7c-41e1-674d-16c6-3431d0d4cc79.htm) | 构建一个S帧协议的内容，需要传入接收需要信息 |
| 公共方法静态成员 | [BuildFrameUMessage](5f188dfb-499f-8e9f-0702-2a9a776610b9.htm) | 构建一个U帧消息的报文信息，传入功能码，STARTDT: 0x07, STOPDT: 0x13; TESTFR: 0x43 |
| 公共方法静态成员 | [BuildWriteIec](7c8079e5-a515-5bb9-a2d5-b725a9a715af.htm) | 构建写入IEC仪表的报文信息 |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法静态成员 | [GetAbsoluteTimeScale](8e2291f8-ea47-1adc-13d8-d9dcf20fb8eb.htm) | 根据给定的时间，获取绝对时标的报文数据信息 |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法静态成员 | [GetTransmissionReasonText](6f2107bf-dba0-3de5-28c1-8867f2adc193.htm) | 获取在传送原因的文本描述信息 |
| 公共方法 | GetType | (继承自 Object。) |
| 公共方法静态成员 | [GetTypeIDText](16d411e8-2d6c-c3e7-9a73-13afdf067b3c.htm) | 获取当前类型的文本描述信息 |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法静态成员 | [PackIEC104Message](0aee7163-21a3-dafa-cb9e-4ce71754840c.htm) | 将IEC104的报文打包成完整的IEC104标准的协议报文 |
| 公共方法静态成员 | [ParseYaoCeValueT](8c7521dc-7e8e-69c3-ba1f-3dff8f9f53fa.htm) | 解析遥信值的方法 |
| 公共方法静态成员 | [PraseTimeFromAbsoluteTimeScale](5644ad87-778b-57e1-ce7d-9e5df893bcb2.htm) | 根据给定的绝对时标的原始内容，解析出实际的时间信息。 |
| 公共方法 | ToString | (继承自 Object。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)字段

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共字段静态成员 | [IEC104ControlStartDT](322f0c4b-400f-0657-a33c-0f7a3c961616.htm) | U帧协议里，启动的功能 |
| 公共字段静态成员 | [IEC104ControlStopDT](57c8b03e-dbec-6985-da98-28158a5f61c8.htm) | U帧协议里，停止的功能 |
| 公共字段静态成员 | [IEC104ControlTestFR](e6bc2d0e-63ba-cb18-d33f-dc91af9d966b.htm) | U帧协议里，测试的功能，主站和子站均可发出 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.Instrument.IEC.Helper 命名空间](acdb00cf-8087-ba54-8ab1-947a92f2b0b4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## IECHelper 构造函数 

[原文連結](http://api.hslcommunication.cn/html/d127c687-70be-3195-2a23-374ad8b2ff17.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.IEC.Helper](../html/acdb00cf-8087-ba54-8ab1-947a92f2b0b4.htm "HslCommunication.Instrument.IEC.Helper")

[IECHelper 类](../html/b37c3cea-b292-985a-f89b-eb97e38d44fe.htm "IECHelper 类")

[IECHelper 构造函数](../html/d127c687-70be-3195-2a23-374ad8b2ff17.htm "IECHelper 构造函数 ")

[IECHelper 方法](../html/bfa98ce5-b54e-bed4-e8cb-ec8cd035ce87.htm "IECHelper 方法")

[IECHelper 字段](../html/3c7bb720-5c4e-776a-8dba-47d8548c956f.htm "IECHelper 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IECHelper 构造函数 |

初始化 [IECHelper](b37c3cea-b292-985a-f89b-eb97e38d44fe.htm) 类的一个新实例

**命名空间：**
 [HslCommunication.Instrument.IEC.Helper](acdb00cf-8087-ba54-8ab1-947a92f2b0b4.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public IECHelper()
```

```
Public Sub New
```

```
public:
IECHelper()
```

```
new : unit -> IECHelper
```

![](../icons/SectionExpanded.png)参见

#### 引用

[IECHelper 类](b37c3cea-b292-985a-f89b-eb97e38d44fe.htm)

[HslCommunication.Instrument.IEC.Helper 命名空间](acdb00cf-8087-ba54-8ab1-947a92f2b0b4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## IECHelper 方法

[原文連結](http://api.hslcommunication.cn/html/bfa98ce5-b54e-bed4-e8cb-ec8cd035ce87.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.IEC.Helper](../html/acdb00cf-8087-ba54-8ab1-947a92f2b0b4.htm "HslCommunication.Instrument.IEC.Helper")

[IECHelper 类](../html/b37c3cea-b292-985a-f89b-eb97e38d44fe.htm "IECHelper 类")

[IECHelper 方法](../html/bfa98ce5-b54e-bed4-e8cb-ec8cd035ce87.htm "IECHelper 方法")

[BuildFrameIMessage 方法](../html/e202c7d8-644d-7f16-c4f3-062e9b14752e.htm "BuildFrameIMessage 方法 ")

[BuildFrameSMessage 方法](../html/c15c3a7c-41e1-674d-16c6-3431d0d4cc79.htm "BuildFrameSMessage 方法 ")

[BuildFrameUMessage 方法](../html/5f188dfb-499f-8e9f-0702-2a9a776610b9.htm "BuildFrameUMessage 方法 ")

[BuildWriteIec 方法](../html/7c8079e5-a515-5bb9-a2d5-b725a9a715af.htm "BuildWriteIec 方法 ")

[GetAbsoluteTimeScale 方法](../html/8e2291f8-ea47-1adc-13d8-d9dcf20fb8eb.htm "GetAbsoluteTimeScale 方法 ")

[GetTransmissionReasonText 方法](../html/6f2107bf-dba0-3de5-28c1-8867f2adc193.htm "GetTransmissionReasonText 方法 ")

[GetTypeIDText 方法](../html/16d411e8-2d6c-c3e7-9a73-13afdf067b3c.htm "GetTypeIDText 方法 ")

[PackIEC104Message 方法](../html/0aee7163-21a3-dafa-cb9e-4ce71754840c.htm "PackIEC104Message 方法 ")

[ParseYaoCeValue(T) 方法](../html/8c7521dc-7e8e-69c3-ba1f-3dff8f9f53fa.htm "ParseYaoCeValue(T) 方法 ")

[PraseTimeFromAbsoluteTimeScale 方法](../html/5644ad87-778b-57e1-ce7d-9e5df893bcb2.htm "PraseTimeFromAbsoluteTimeScale 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IECHelper 方法 |

[IECHelper](b37c3cea-b292-985a-f89b-eb97e38d44fe.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法静态成员 | [BuildFrameIMessage](e202c7d8-644d-7f16-c4f3-062e9b14752e.htm) | 构建一个I帧消息的报文信息，传入相关的参数信息，返回完整的104消息报文 |
| 公共方法静态成员 | [BuildFrameSMessage](c15c3a7c-41e1-674d-16c6-3431d0d4cc79.htm) | 构建一个S帧协议的内容，需要传入接收需要信息 |
| 公共方法静态成员 | [BuildFrameUMessage](5f188dfb-499f-8e9f-0702-2a9a776610b9.htm) | 构建一个U帧消息的报文信息，传入功能码，STARTDT: 0x07, STOPDT: 0x13; TESTFR: 0x43 |
| 公共方法静态成员 | [BuildWriteIec](7c8079e5-a515-5bb9-a2d5-b725a9a715af.htm) | 构建写入IEC仪表的报文信息 |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法静态成员 | [GetAbsoluteTimeScale](8e2291f8-ea47-1adc-13d8-d9dcf20fb8eb.htm) | 根据给定的时间，获取绝对时标的报文数据信息 |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法静态成员 | [GetTransmissionReasonText](6f2107bf-dba0-3de5-28c1-8867f2adc193.htm) | 获取在传送原因的文本描述信息 |
| 公共方法 | GetType | (继承自 Object。) |
| 公共方法静态成员 | [GetTypeIDText](16d411e8-2d6c-c3e7-9a73-13afdf067b3c.htm) | 获取当前类型的文本描述信息 |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法静态成员 | [PackIEC104Message](0aee7163-21a3-dafa-cb9e-4ce71754840c.htm) | 将IEC104的报文打包成完整的IEC104标准的协议报文 |
| 公共方法静态成员 | [ParseYaoCeValueT](8c7521dc-7e8e-69c3-ba1f-3dff8f9f53fa.htm) | 解析遥信值的方法 |
| 公共方法静态成员 | [PraseTimeFromAbsoluteTimeScale](5644ad87-778b-57e1-ce7d-9e5df893bcb2.htm) | 根据给定的绝对时标的原始内容，解析出实际的时间信息。 |
| 公共方法 | ToString | (继承自 Object。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[IECHelper 类](b37c3cea-b292-985a-f89b-eb97e38d44fe.htm)

[HslCommunication.Instrument.IEC.Helper 命名空间](acdb00cf-8087-ba54-8ab1-947a92f2b0b4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## BuildFrameIMessage 方法 

[原文連結](http://api.hslcommunication.cn/html/e202c7d8-644d-7f16-c4f3-062e9b14752e.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.IEC.Helper](../html/acdb00cf-8087-ba54-8ab1-947a92f2b0b4.htm "HslCommunication.Instrument.IEC.Helper")

[IECHelper 类](../html/b37c3cea-b292-985a-f89b-eb97e38d44fe.htm "IECHelper 类")

[IECHelper 方法](../html/bfa98ce5-b54e-bed4-e8cb-ec8cd035ce87.htm "IECHelper 方法")

[BuildFrameIMessage 方法](../html/e202c7d8-644d-7f16-c4f3-062e9b14752e.htm "BuildFrameIMessage 方法 ")

[BuildFrameSMessage 方法](../html/c15c3a7c-41e1-674d-16c6-3431d0d4cc79.htm "BuildFrameSMessage 方法 ")

[BuildFrameUMessage 方法](../html/5f188dfb-499f-8e9f-0702-2a9a776610b9.htm "BuildFrameUMessage 方法 ")

[BuildWriteIec 方法](../html/7c8079e5-a515-5bb9-a2d5-b725a9a715af.htm "BuildWriteIec 方法 ")

[GetAbsoluteTimeScale 方法](../html/8e2291f8-ea47-1adc-13d8-d9dcf20fb8eb.htm "GetAbsoluteTimeScale 方法 ")

[GetTransmissionReasonText 方法](../html/6f2107bf-dba0-3de5-28c1-8867f2adc193.htm "GetTransmissionReasonText 方法 ")

[GetTypeIDText 方法](../html/16d411e8-2d6c-c3e7-9a73-13afdf067b3c.htm "GetTypeIDText 方法 ")

[PackIEC104Message 方法](../html/0aee7163-21a3-dafa-cb9e-4ce71754840c.htm "PackIEC104Message 方法 ")

[ParseYaoCeValue(T) 方法](../html/8c7521dc-7e8e-69c3-ba1f-3dff8f9f53fa.htm "ParseYaoCeValue(T) 方法 ")

[PraseTimeFromAbsoluteTimeScale 方法](../html/5644ad87-778b-57e1-ce7d-9e5df893bcb2.htm "PraseTimeFromAbsoluteTimeScale 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IECHelperBuildFrameIMessage 方法 |

构建一个I帧消息的报文信息，传入相关的参数信息，返回完整的104消息报文

**命名空间：**
 [HslCommunication.Instrument.IEC.Helper](acdb00cf-8087-ba54-8ab1-947a92f2b0b4.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static byte[] BuildFrameIMessage(
	int sendID,
	int receiveID,
	byte typeId,
	byte variableStructureQualifier,
	ushort reason,
	ushort station,
	byte[] body
)
```

```
Public Shared Function BuildFrameIMessage ( 
	sendID As Integer,
	receiveID As Integer,
	typeId As Byte,
	variableStructureQualifier As Byte,
	reason As UShort,
	station As UShort,
	body As Byte()
) As Byte()
```

```
public:
static array<unsigned char>^ BuildFrameIMessage(
	int sendID, 
	int receiveID, 
	unsigned char typeId, 
	unsigned char variableStructureQualifier, 
	unsigned short reason, 
	unsigned short station, 
	array<unsigned char>^ body
)
```

```
static member BuildFrameIMessage : 
        sendID : int * 
        receiveID : int * 
        typeId : byte * 
        variableStructureQualifier : byte * 
        reason : uint16 * 
        station : uint16 * 
        body : byte[] -> byte[] 
```

#### 参数

sendID
:   类型：SystemInt32  
    发送的序列号

receiveID
:   类型：SystemInt32  
    接收的序列号

typeId
:   类型：SystemByte  
    类型标识

variableStructureQualifier
:   类型：SystemByte  
    可变结构限定词

reason
:   类型：SystemUInt16  
    传送原因

station
:   类型：SystemUInt16  
    应用服务数据单元公共地址

body
:   类型：SystemByte  
    信息体，最大243个字节的长度

#### 返回值

类型：Byte  
用于发送的104报文信息

![](../icons/SectionExpanded.png)参见

#### 引用

[IECHelper 类](b37c3cea-b292-985a-f89b-eb97e38d44fe.htm)

[HslCommunication.Instrument.IEC.Helper 命名空间](acdb00cf-8087-ba54-8ab1-947a92f2b0b4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## BuildFrameSMessage 方法 

[原文連結](http://api.hslcommunication.cn/html/c15c3a7c-41e1-674d-16c6-3431d0d4cc79.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.IEC.Helper](../html/acdb00cf-8087-ba54-8ab1-947a92f2b0b4.htm "HslCommunication.Instrument.IEC.Helper")

[IECHelper 类](../html/b37c3cea-b292-985a-f89b-eb97e38d44fe.htm "IECHelper 类")

[IECHelper 方法](../html/bfa98ce5-b54e-bed4-e8cb-ec8cd035ce87.htm "IECHelper 方法")

[BuildFrameIMessage 方法](../html/e202c7d8-644d-7f16-c4f3-062e9b14752e.htm "BuildFrameIMessage 方法 ")

[BuildFrameSMessage 方法](../html/c15c3a7c-41e1-674d-16c6-3431d0d4cc79.htm "BuildFrameSMessage 方法 ")

[BuildFrameUMessage 方法](../html/5f188dfb-499f-8e9f-0702-2a9a776610b9.htm "BuildFrameUMessage 方法 ")

[BuildWriteIec 方法](../html/7c8079e5-a515-5bb9-a2d5-b725a9a715af.htm "BuildWriteIec 方法 ")

[GetAbsoluteTimeScale 方法](../html/8e2291f8-ea47-1adc-13d8-d9dcf20fb8eb.htm "GetAbsoluteTimeScale 方法 ")

[GetTransmissionReasonText 方法](../html/6f2107bf-dba0-3de5-28c1-8867f2adc193.htm "GetTransmissionReasonText 方法 ")

[GetTypeIDText 方法](../html/16d411e8-2d6c-c3e7-9a73-13afdf067b3c.htm "GetTypeIDText 方法 ")

[PackIEC104Message 方法](../html/0aee7163-21a3-dafa-cb9e-4ce71754840c.htm "PackIEC104Message 方法 ")

[ParseYaoCeValue(T) 方法](../html/8c7521dc-7e8e-69c3-ba1f-3dff8f9f53fa.htm "ParseYaoCeValue(T) 方法 ")

[PraseTimeFromAbsoluteTimeScale 方法](../html/5644ad87-778b-57e1-ce7d-9e5df893bcb2.htm "PraseTimeFromAbsoluteTimeScale 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IECHelperBuildFrameSMessage 方法 |

构建一个S帧协议的内容，需要传入接收需要信息

**命名空间：**
 [HslCommunication.Instrument.IEC.Helper](acdb00cf-8087-ba54-8ab1-947a92f2b0b4.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static byte[] BuildFrameSMessage(
	int receiveID
)
```

```
Public Shared Function BuildFrameSMessage ( 
	receiveID As Integer
) As Byte()
```

```
public:
static array<unsigned char>^ BuildFrameSMessage(
	int receiveID
)
```

```
static member BuildFrameSMessage : 
        receiveID : int -> byte[] 
```

#### 参数

receiveID
:   类型：SystemInt32  
    接收序号信息

#### 返回值

类型：Byte  
S帧协议的报文信息

![](../icons/SectionExpanded.png)参见

#### 引用

[IECHelper 类](b37c3cea-b292-985a-f89b-eb97e38d44fe.htm)

[HslCommunication.Instrument.IEC.Helper 命名空间](acdb00cf-8087-ba54-8ab1-947a92f2b0b4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## BuildFrameUMessage 方法 

[原文連結](http://api.hslcommunication.cn/html/5f188dfb-499f-8e9f-0702-2a9a776610b9.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.IEC.Helper](../html/acdb00cf-8087-ba54-8ab1-947a92f2b0b4.htm "HslCommunication.Instrument.IEC.Helper")

[IECHelper 类](../html/b37c3cea-b292-985a-f89b-eb97e38d44fe.htm "IECHelper 类")

[IECHelper 方法](../html/bfa98ce5-b54e-bed4-e8cb-ec8cd035ce87.htm "IECHelper 方法")

[BuildFrameIMessage 方法](../html/e202c7d8-644d-7f16-c4f3-062e9b14752e.htm "BuildFrameIMessage 方法 ")

[BuildFrameSMessage 方法](../html/c15c3a7c-41e1-674d-16c6-3431d0d4cc79.htm "BuildFrameSMessage 方法 ")

[BuildFrameUMessage 方法](../html/5f188dfb-499f-8e9f-0702-2a9a776610b9.htm "BuildFrameUMessage 方法 ")

[BuildWriteIec 方法](../html/7c8079e5-a515-5bb9-a2d5-b725a9a715af.htm "BuildWriteIec 方法 ")

[GetAbsoluteTimeScale 方法](../html/8e2291f8-ea47-1adc-13d8-d9dcf20fb8eb.htm "GetAbsoluteTimeScale 方法 ")

[GetTransmissionReasonText 方法](../html/6f2107bf-dba0-3de5-28c1-8867f2adc193.htm "GetTransmissionReasonText 方法 ")

[GetTypeIDText 方法](../html/16d411e8-2d6c-c3e7-9a73-13afdf067b3c.htm "GetTypeIDText 方法 ")

[PackIEC104Message 方法](../html/0aee7163-21a3-dafa-cb9e-4ce71754840c.htm "PackIEC104Message 方法 ")

[ParseYaoCeValue(T) 方法](../html/8c7521dc-7e8e-69c3-ba1f-3dff8f9f53fa.htm "ParseYaoCeValue(T) 方法 ")

[PraseTimeFromAbsoluteTimeScale 方法](../html/5644ad87-778b-57e1-ce7d-9e5df893bcb2.htm "PraseTimeFromAbsoluteTimeScale 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IECHelperBuildFrameUMessage 方法 |

构建一个U帧消息的报文信息，传入功能码，STARTDT: 0x07, STOPDT: 0x13; TESTFR: 0x43

**命名空间：**
 [HslCommunication.Instrument.IEC.Helper](acdb00cf-8087-ba54-8ab1-947a92f2b0b4.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static byte[] BuildFrameUMessage(
	byte controlField
)
```

```
Public Shared Function BuildFrameUMessage ( 
	controlField As Byte
) As Byte()
```

```
public:
static array<unsigned char>^ BuildFrameUMessage(
	unsigned char controlField
)
```

```
static member BuildFrameUMessage : 
        controlField : byte -> byte[] 
```

#### 参数

controlField
:   类型：SystemByte  
    控制码信息

#### 返回值

类型：Byte  
U帧的报文信息

![](../icons/SectionExpanded.png)参见

#### 引用

[IECHelper 类](b37c3cea-b292-985a-f89b-eb97e38d44fe.htm)

[HslCommunication.Instrument.IEC.Helper 命名空间](acdb00cf-8087-ba54-8ab1-947a92f2b0b4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## BuildWriteIec 方法 

[原文連結](http://api.hslcommunication.cn/html/7c8079e5-a515-5bb9-a2d5-b725a9a715af.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.IEC.Helper](../html/acdb00cf-8087-ba54-8ab1-947a92f2b0b4.htm "HslCommunication.Instrument.IEC.Helper")

[IECHelper 类](../html/b37c3cea-b292-985a-f89b-eb97e38d44fe.htm "IECHelper 类")

[IECHelper 方法](../html/bfa98ce5-b54e-bed4-e8cb-ec8cd035ce87.htm "IECHelper 方法")

[BuildFrameIMessage 方法](../html/e202c7d8-644d-7f16-c4f3-062e9b14752e.htm "BuildFrameIMessage 方法 ")

[BuildFrameSMessage 方法](../html/c15c3a7c-41e1-674d-16c6-3431d0d4cc79.htm "BuildFrameSMessage 方法 ")

[BuildFrameUMessage 方法](../html/5f188dfb-499f-8e9f-0702-2a9a776610b9.htm "BuildFrameUMessage 方法 ")

[BuildWriteIec 方法](../html/7c8079e5-a515-5bb9-a2d5-b725a9a715af.htm "BuildWriteIec 方法 ")

[GetAbsoluteTimeScale 方法](../html/8e2291f8-ea47-1adc-13d8-d9dcf20fb8eb.htm "GetAbsoluteTimeScale 方法 ")

[GetTransmissionReasonText 方法](../html/6f2107bf-dba0-3de5-28c1-8867f2adc193.htm "GetTransmissionReasonText 方法 ")

[GetTypeIDText 方法](../html/16d411e8-2d6c-c3e7-9a73-13afdf067b3c.htm "GetTypeIDText 方法 ")

[PackIEC104Message 方法](../html/0aee7163-21a3-dafa-cb9e-4ce71754840c.htm "PackIEC104Message 方法 ")

[ParseYaoCeValue(T) 方法](../html/8c7521dc-7e8e-69c3-ba1f-3dff8f9f53fa.htm "ParseYaoCeValue(T) 方法 ")

[PraseTimeFromAbsoluteTimeScale 方法](../html/5644ad87-778b-57e1-ce7d-9e5df893bcb2.htm "PraseTimeFromAbsoluteTimeScale 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IECHelperBuildWriteIec 方法 |

构建写入IEC仪表的报文信息

**命名空间：**
 [HslCommunication.Instrument.IEC.Helper](acdb00cf-8087-ba54-8ab1-947a92f2b0b4.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static byte[] BuildWriteIec(
	byte type,
	ushort reason,
	ushort station,
	ushort address,
	byte[] value
)
```

```
Public Shared Function BuildWriteIec ( 
	type As Byte,
	reason As UShort,
	station As UShort,
	address As UShort,
	value As Byte()
) As Byte()
```

```
public:
static array<unsigned char>^ BuildWriteIec(
	unsigned char type, 
	unsigned short reason, 
	unsigned short station, 
	unsigned short address, 
	array<unsigned char>^ value
)
```

```
static member BuildWriteIec : 
        type : byte * 
        reason : uint16 * 
        station : uint16 * 
        address : uint16 * 
        value : byte[] -> byte[] 
```

#### 参数

type
:   类型：SystemByte  
    指令类型信息

reason
:   类型：SystemUInt16  
    原因信息

station
:   类型：SystemUInt16  
    公共单元地址

address
:   类型：SystemUInt16  
    信息对象地址

value
:   类型：SystemByte  
    值数据

#### 返回值

类型：Byte  
发送仪表的报文

![](../icons/SectionExpanded.png)参见

#### 引用

[IECHelper 类](b37c3cea-b292-985a-f89b-eb97e38d44fe.htm)

[HslCommunication.Instrument.IEC.Helper 命名空间](acdb00cf-8087-ba54-8ab1-947a92f2b0b4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GetAbsoluteTimeScale 方法 

[原文連結](http://api.hslcommunication.cn/html/8e2291f8-ea47-1adc-13d8-d9dcf20fb8eb.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.IEC.Helper](../html/acdb00cf-8087-ba54-8ab1-947a92f2b0b4.htm "HslCommunication.Instrument.IEC.Helper")

[IECHelper 类](../html/b37c3cea-b292-985a-f89b-eb97e38d44fe.htm "IECHelper 类")

[IECHelper 方法](../html/bfa98ce5-b54e-bed4-e8cb-ec8cd035ce87.htm "IECHelper 方法")

[BuildFrameIMessage 方法](../html/e202c7d8-644d-7f16-c4f3-062e9b14752e.htm "BuildFrameIMessage 方法 ")

[BuildFrameSMessage 方法](../html/c15c3a7c-41e1-674d-16c6-3431d0d4cc79.htm "BuildFrameSMessage 方法 ")

[BuildFrameUMessage 方法](../html/5f188dfb-499f-8e9f-0702-2a9a776610b9.htm "BuildFrameUMessage 方法 ")

[BuildWriteIec 方法](../html/7c8079e5-a515-5bb9-a2d5-b725a9a715af.htm "BuildWriteIec 方法 ")

[GetAbsoluteTimeScale 方法](../html/8e2291f8-ea47-1adc-13d8-d9dcf20fb8eb.htm "GetAbsoluteTimeScale 方法 ")

[GetTransmissionReasonText 方法](../html/6f2107bf-dba0-3de5-28c1-8867f2adc193.htm "GetTransmissionReasonText 方法 ")

[GetTypeIDText 方法](../html/16d411e8-2d6c-c3e7-9a73-13afdf067b3c.htm "GetTypeIDText 方法 ")

[PackIEC104Message 方法](../html/0aee7163-21a3-dafa-cb9e-4ce71754840c.htm "PackIEC104Message 方法 ")

[ParseYaoCeValue(T) 方法](../html/8c7521dc-7e8e-69c3-ba1f-3dff8f9f53fa.htm "ParseYaoCeValue(T) 方法 ")

[PraseTimeFromAbsoluteTimeScale 方法](../html/5644ad87-778b-57e1-ce7d-9e5df893bcb2.htm "PraseTimeFromAbsoluteTimeScale 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IECHelperGetAbsoluteTimeScale 方法 |

根据给定的时间，获取绝对时标的报文数据信息

**命名空间：**
 [HslCommunication.Instrument.IEC.Helper](acdb00cf-8087-ba54-8ab1-947a92f2b0b4.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static byte[] GetAbsoluteTimeScale(
	DateTime dateTime,
	bool valid
)
```

```
Public Shared Function GetAbsoluteTimeScale ( 
	dateTime As DateTime,
	valid As Boolean
) As Byte()
```

```
public:
static array<unsigned char>^ GetAbsoluteTimeScale(
	DateTime dateTime, 
	bool valid
)
```

```
static member GetAbsoluteTimeScale : 
        dateTime : DateTime * 
        valid : bool -> byte[] 
```

#### 参数

dateTime
:   类型：SystemDateTime  
    时间信息

valid
:   类型：SystemBoolean  
    时标是否有效

#### 返回值

类型：Byte  
可用于发送的绝对时标的报文

![](../icons/SectionExpanded.png)参见

#### 引用

[IECHelper 类](b37c3cea-b292-985a-f89b-eb97e38d44fe.htm)

[HslCommunication.Instrument.IEC.Helper 命名空间](acdb00cf-8087-ba54-8ab1-947a92f2b0b4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GetTransmissionReasonText 方法 

[原文連結](http://api.hslcommunication.cn/html/6f2107bf-dba0-3de5-28c1-8867f2adc193.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.IEC.Helper](../html/acdb00cf-8087-ba54-8ab1-947a92f2b0b4.htm "HslCommunication.Instrument.IEC.Helper")

[IECHelper 类](../html/b37c3cea-b292-985a-f89b-eb97e38d44fe.htm "IECHelper 类")

[IECHelper 方法](../html/bfa98ce5-b54e-bed4-e8cb-ec8cd035ce87.htm "IECHelper 方法")

[BuildFrameIMessage 方法](../html/e202c7d8-644d-7f16-c4f3-062e9b14752e.htm "BuildFrameIMessage 方法 ")

[BuildFrameSMessage 方法](../html/c15c3a7c-41e1-674d-16c6-3431d0d4cc79.htm "BuildFrameSMessage 方法 ")

[BuildFrameUMessage 方法](../html/5f188dfb-499f-8e9f-0702-2a9a776610b9.htm "BuildFrameUMessage 方法 ")

[BuildWriteIec 方法](../html/7c8079e5-a515-5bb9-a2d5-b725a9a715af.htm "BuildWriteIec 方法 ")

[GetAbsoluteTimeScale 方法](../html/8e2291f8-ea47-1adc-13d8-d9dcf20fb8eb.htm "GetAbsoluteTimeScale 方法 ")

[GetTransmissionReasonText 方法](../html/6f2107bf-dba0-3de5-28c1-8867f2adc193.htm "GetTransmissionReasonText 方法 ")

[GetTypeIDText 方法](../html/16d411e8-2d6c-c3e7-9a73-13afdf067b3c.htm "GetTypeIDText 方法 ")

[PackIEC104Message 方法](../html/0aee7163-21a3-dafa-cb9e-4ce71754840c.htm "PackIEC104Message 方法 ")

[ParseYaoCeValue(T) 方法](../html/8c7521dc-7e8e-69c3-ba1f-3dff8f9f53fa.htm "ParseYaoCeValue(T) 方法 ")

[PraseTimeFromAbsoluteTimeScale 方法](../html/5644ad87-778b-57e1-ce7d-9e5df893bcb2.htm "PraseTimeFromAbsoluteTimeScale 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IECHelperGetTransmissionReasonText 方法 |

获取在传送原因的文本描述信息

**命名空间：**
 [HslCommunication.Instrument.IEC.Helper](acdb00cf-8087-ba54-8ab1-947a92f2b0b4.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static string GetTransmissionReasonText(
	int transmissionReason
)
```

```
Public Shared Function GetTransmissionReasonText ( 
	transmissionReason As Integer
) As String
```

```
public:
static String^ GetTransmissionReasonText(
	int transmissionReason
)
```

```
static member GetTransmissionReasonText : 
        transmissionReason : int -> string 
```

#### 参数

transmissionReason
:   类型：SystemInt32  

    [缺少 "M:HslCommunication.Instrument.IEC.Helper.IECHelper.GetTransmissionReasonText(System.Int32)" 的 <param name="transmissionReason"/> 文档]

#### 返回值

类型：String  
字符串值

![](../icons/SectionExpanded.png)参见

#### 引用

[IECHelper 类](b37c3cea-b292-985a-f89b-eb97e38d44fe.htm)

[HslCommunication.Instrument.IEC.Helper 命名空间](acdb00cf-8087-ba54-8ab1-947a92f2b0b4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GetTypeIDText 方法 

[原文連結](http://api.hslcommunication.cn/html/16d411e8-2d6c-c3e7-9a73-13afdf067b3c.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.IEC.Helper](../html/acdb00cf-8087-ba54-8ab1-947a92f2b0b4.htm "HslCommunication.Instrument.IEC.Helper")

[IECHelper 类](../html/b37c3cea-b292-985a-f89b-eb97e38d44fe.htm "IECHelper 类")

[IECHelper 方法](../html/bfa98ce5-b54e-bed4-e8cb-ec8cd035ce87.htm "IECHelper 方法")

[BuildFrameIMessage 方法](../html/e202c7d8-644d-7f16-c4f3-062e9b14752e.htm "BuildFrameIMessage 方法 ")

[BuildFrameSMessage 方法](../html/c15c3a7c-41e1-674d-16c6-3431d0d4cc79.htm "BuildFrameSMessage 方法 ")

[BuildFrameUMessage 方法](../html/5f188dfb-499f-8e9f-0702-2a9a776610b9.htm "BuildFrameUMessage 方法 ")

[BuildWriteIec 方法](../html/7c8079e5-a515-5bb9-a2d5-b725a9a715af.htm "BuildWriteIec 方法 ")

[GetAbsoluteTimeScale 方法](../html/8e2291f8-ea47-1adc-13d8-d9dcf20fb8eb.htm "GetAbsoluteTimeScale 方法 ")

[GetTransmissionReasonText 方法](../html/6f2107bf-dba0-3de5-28c1-8867f2adc193.htm "GetTransmissionReasonText 方法 ")

[GetTypeIDText 方法](../html/16d411e8-2d6c-c3e7-9a73-13afdf067b3c.htm "GetTypeIDText 方法 ")

[PackIEC104Message 方法](../html/0aee7163-21a3-dafa-cb9e-4ce71754840c.htm "PackIEC104Message 方法 ")

[ParseYaoCeValue(T) 方法](../html/8c7521dc-7e8e-69c3-ba1f-3dff8f9f53fa.htm "ParseYaoCeValue(T) 方法 ")

[PraseTimeFromAbsoluteTimeScale 方法](../html/5644ad87-778b-57e1-ce7d-9e5df893bcb2.htm "PraseTimeFromAbsoluteTimeScale 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IECHelperGetTypeIDText 方法 |

获取当前类型的文本描述信息

**命名空间：**
 [HslCommunication.Instrument.IEC.Helper](acdb00cf-8087-ba54-8ab1-947a92f2b0b4.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static string GetTypeIDText(
	int typeId
)
```

```
Public Shared Function GetTypeIDText ( 
	typeId As Integer
) As String
```

```
public:
static String^ GetTypeIDText(
	int typeId
)
```

```
static member GetTypeIDText : 
        typeId : int -> string 
```

#### 参数

typeId
:   类型：SystemInt32  

    [缺少 "M:HslCommunication.Instrument.IEC.Helper.IECHelper.GetTypeIDText(System.Int32)" 的 <param name="typeId"/> 文档]

#### 返回值

类型：String  
文本值

![](../icons/SectionExpanded.png)参见

#### 引用

[IECHelper 类](b37c3cea-b292-985a-f89b-eb97e38d44fe.htm)

[HslCommunication.Instrument.IEC.Helper 命名空间](acdb00cf-8087-ba54-8ab1-947a92f2b0b4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## PackIEC104Message 方法 

[原文連結](http://api.hslcommunication.cn/html/0aee7163-21a3-dafa-cb9e-4ce71754840c.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.IEC.Helper](../html/acdb00cf-8087-ba54-8ab1-947a92f2b0b4.htm "HslCommunication.Instrument.IEC.Helper")

[IECHelper 类](../html/b37c3cea-b292-985a-f89b-eb97e38d44fe.htm "IECHelper 类")

[IECHelper 方法](../html/bfa98ce5-b54e-bed4-e8cb-ec8cd035ce87.htm "IECHelper 方法")

[BuildFrameIMessage 方法](../html/e202c7d8-644d-7f16-c4f3-062e9b14752e.htm "BuildFrameIMessage 方法 ")

[BuildFrameSMessage 方法](../html/c15c3a7c-41e1-674d-16c6-3431d0d4cc79.htm "BuildFrameSMessage 方法 ")

[BuildFrameUMessage 方法](../html/5f188dfb-499f-8e9f-0702-2a9a776610b9.htm "BuildFrameUMessage 方法 ")

[BuildWriteIec 方法](../html/7c8079e5-a515-5bb9-a2d5-b725a9a715af.htm "BuildWriteIec 方法 ")

[GetAbsoluteTimeScale 方法](../html/8e2291f8-ea47-1adc-13d8-d9dcf20fb8eb.htm "GetAbsoluteTimeScale 方法 ")

[GetTransmissionReasonText 方法](../html/6f2107bf-dba0-3de5-28c1-8867f2adc193.htm "GetTransmissionReasonText 方法 ")

[GetTypeIDText 方法](../html/16d411e8-2d6c-c3e7-9a73-13afdf067b3c.htm "GetTypeIDText 方法 ")

[PackIEC104Message 方法](../html/0aee7163-21a3-dafa-cb9e-4ce71754840c.htm "PackIEC104Message 方法 ")

[ParseYaoCeValue(T) 方法](../html/8c7521dc-7e8e-69c3-ba1f-3dff8f9f53fa.htm "ParseYaoCeValue(T) 方法 ")

[PraseTimeFromAbsoluteTimeScale 方法](../html/5644ad87-778b-57e1-ce7d-9e5df893bcb2.htm "PraseTimeFromAbsoluteTimeScale 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IECHelperPackIEC104Message 方法 |

将IEC104的报文打包成完整的IEC104标准的协议报文

**命名空间：**
 [HslCommunication.Instrument.IEC.Helper](acdb00cf-8087-ba54-8ab1-947a92f2b0b4.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static byte[] PackIEC104Message(
	byte controlField1,
	byte controlField2,
	byte controlField3,
	byte controlField4,
	byte[] asdu
)
```

```
Public Shared Function PackIEC104Message ( 
	controlField1 As Byte,
	controlField2 As Byte,
	controlField3 As Byte,
	controlField4 As Byte,
	asdu As Byte()
) As Byte()
```

```
public:
static array<unsigned char>^ PackIEC104Message(
	unsigned char controlField1, 
	unsigned char controlField2, 
	unsigned char controlField3, 
	unsigned char controlField4, 
	array<unsigned char>^ asdu
)
```

```
static member PackIEC104Message : 
        controlField1 : byte * 
        controlField2 : byte * 
        controlField3 : byte * 
        controlField4 : byte * 
        asdu : byte[] -> byte[] 
```

#### 参数

controlField1
:   类型：SystemByte  
    控制域1

controlField2
:   类型：SystemByte  
    控制域2

controlField3
:   类型：SystemByte  
    控制域3

controlField4
:   类型：SystemByte  
    控制域4

asdu
:   类型：SystemByte  
    ASDU报文，包含类型标识，可变结构限定词，传送原因，应用服务器数据单元公共地址，信息体

#### 返回值

类型：Byte  
完整的报文消息

![](../icons/SectionExpanded.png)参见

#### 引用

[IECHelper 类](b37c3cea-b292-985a-f89b-eb97e38d44fe.htm)

[HslCommunication.Instrument.IEC.Helper 命名空间](acdb00cf-8087-ba54-8ab1-947a92f2b0b4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ParseYaoCeValue(T) 方法 

[原文連結](http://api.hslcommunication.cn/html/8c7521dc-7e8e-69c3-ba1f-3dff8f9f53fa.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.IEC.Helper](../html/acdb00cf-8087-ba54-8ab1-947a92f2b0b4.htm "HslCommunication.Instrument.IEC.Helper")

[IECHelper 类](../html/b37c3cea-b292-985a-f89b-eb97e38d44fe.htm "IECHelper 类")

[IECHelper 方法](../html/bfa98ce5-b54e-bed4-e8cb-ec8cd035ce87.htm "IECHelper 方法")

[BuildFrameIMessage 方法](../html/e202c7d8-644d-7f16-c4f3-062e9b14752e.htm "BuildFrameIMessage 方法 ")

[BuildFrameSMessage 方法](../html/c15c3a7c-41e1-674d-16c6-3431d0d4cc79.htm "BuildFrameSMessage 方法 ")

[BuildFrameUMessage 方法](../html/5f188dfb-499f-8e9f-0702-2a9a776610b9.htm "BuildFrameUMessage 方法 ")

[BuildWriteIec 方法](../html/7c8079e5-a515-5bb9-a2d5-b725a9a715af.htm "BuildWriteIec 方法 ")

[GetAbsoluteTimeScale 方法](../html/8e2291f8-ea47-1adc-13d8-d9dcf20fb8eb.htm "GetAbsoluteTimeScale 方法 ")

[GetTransmissionReasonText 方法](../html/6f2107bf-dba0-3de5-28c1-8867f2adc193.htm "GetTransmissionReasonText 方法 ")

[GetTypeIDText 方法](../html/16d411e8-2d6c-c3e7-9a73-13afdf067b3c.htm "GetTypeIDText 方法 ")

[PackIEC104Message 方法](../html/0aee7163-21a3-dafa-cb9e-4ce71754840c.htm "PackIEC104Message 方法 ")

[ParseYaoCeValue(T) 方法](../html/8c7521dc-7e8e-69c3-ba1f-3dff8f9f53fa.htm "ParseYaoCeValue(T) 方法 ")

[PraseTimeFromAbsoluteTimeScale 方法](../html/5644ad87-778b-57e1-ce7d-9e5df893bcb2.htm "PraseTimeFromAbsoluteTimeScale 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IECHelperParseYaoCeValueT 方法 |

解析遥信值的方法

**命名空间：**
 [HslCommunication.Instrument.IEC.Helper](acdb00cf-8087-ba54-8ab1-947a92f2b0b4.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static List<IecValueObject<T>> ParseYaoCeValue<T>(
	IEC104MessageEventArgs message,
	Func<byte[], int, T> trans,
	int unitLength
)
```

```
Public Shared Function ParseYaoCeValue(Of T) ( 
	message As IEC104MessageEventArgs,
	trans As Func(Of Byte(), Integer, T),
	unitLength As Integer
) As List(Of IecValueObject(Of T))
```

```
public:
generic<typename T>
static List<IecValueObject<T>^>^ ParseYaoCeValue(
	IEC104MessageEventArgs^ message, 
	Func<array<unsigned char>^, int, T>^ trans, 
	int unitLength
)
```

```
static member ParseYaoCeValue : 
        message : IEC104MessageEventArgs * 
        trans : Func<byte[], int, 'T> * 
        unitLength : int -> List<IecValueObject<'T>> 
```

#### 参数

message
:   类型：[HslCommunication.Instrument.IECIEC104MessageEventArgs](34f93396-e777-2b46-1db2-20c3f8ec3d23.htm)  
    IEC104的消息

trans
:   类型：SystemFuncByte, Int32, T  
    从实际的字节数据转换指定类型的方法

unitLength
:   类型：SystemInt32  
    数据类型的字节长度信息

#### 类型参数

T
:   转换后的类型信息

#### 返回值

类型：List[IecValueObject](dd68297c-f009-b216-89a4-112d4eecb574.htm)T  
列表值

![](../icons/SectionExpanded.png)参见

#### 引用

[IECHelper 类](b37c3cea-b292-985a-f89b-eb97e38d44fe.htm)

[HslCommunication.Instrument.IEC.Helper 命名空间](acdb00cf-8087-ba54-8ab1-947a92f2b0b4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## PraseTimeFromAbsoluteTimeScale 方法 

[原文連結](http://api.hslcommunication.cn/html/5644ad87-778b-57e1-ce7d-9e5df893bcb2.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.IEC.Helper](../html/acdb00cf-8087-ba54-8ab1-947a92f2b0b4.htm "HslCommunication.Instrument.IEC.Helper")

[IECHelper 类](../html/b37c3cea-b292-985a-f89b-eb97e38d44fe.htm "IECHelper 类")

[IECHelper 方法](../html/bfa98ce5-b54e-bed4-e8cb-ec8cd035ce87.htm "IECHelper 方法")

[BuildFrameIMessage 方法](../html/e202c7d8-644d-7f16-c4f3-062e9b14752e.htm "BuildFrameIMessage 方法 ")

[BuildFrameSMessage 方法](../html/c15c3a7c-41e1-674d-16c6-3431d0d4cc79.htm "BuildFrameSMessage 方法 ")

[BuildFrameUMessage 方法](../html/5f188dfb-499f-8e9f-0702-2a9a776610b9.htm "BuildFrameUMessage 方法 ")

[BuildWriteIec 方法](../html/7c8079e5-a515-5bb9-a2d5-b725a9a715af.htm "BuildWriteIec 方法 ")

[GetAbsoluteTimeScale 方法](../html/8e2291f8-ea47-1adc-13d8-d9dcf20fb8eb.htm "GetAbsoluteTimeScale 方法 ")

[GetTransmissionReasonText 方法](../html/6f2107bf-dba0-3de5-28c1-8867f2adc193.htm "GetTransmissionReasonText 方法 ")

[GetTypeIDText 方法](../html/16d411e8-2d6c-c3e7-9a73-13afdf067b3c.htm "GetTypeIDText 方法 ")

[PackIEC104Message 方法](../html/0aee7163-21a3-dafa-cb9e-4ce71754840c.htm "PackIEC104Message 方法 ")

[ParseYaoCeValue(T) 方法](../html/8c7521dc-7e8e-69c3-ba1f-3dff8f9f53fa.htm "ParseYaoCeValue(T) 方法 ")

[PraseTimeFromAbsoluteTimeScale 方法](../html/5644ad87-778b-57e1-ce7d-9e5df893bcb2.htm "PraseTimeFromAbsoluteTimeScale 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IECHelperPraseTimeFromAbsoluteTimeScale 方法 |

根据给定的绝对时标的原始内容，解析出实际的时间信息。

**命名空间：**
 [HslCommunication.Instrument.IEC.Helper](acdb00cf-8087-ba54-8ab1-947a92f2b0b4.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static DateTime PraseTimeFromAbsoluteTimeScale(
	byte[] source,
	int index
)
```

```
Public Shared Function PraseTimeFromAbsoluteTimeScale ( 
	source As Byte(),
	index As Integer
) As DateTime
```

```
public:
static DateTime PraseTimeFromAbsoluteTimeScale(
	array<unsigned char>^ source, 
	int index
)
```

```
static member PraseTimeFromAbsoluteTimeScale : 
        source : byte[] * 
        index : int -> DateTime 
```

#### 参数

source
:   类型：SystemByte  
    原始字节

index
:   类型：SystemInt32  
    数据的偏移索引

#### 返回值

类型：DateTime  
时间信息

![](../icons/SectionExpanded.png)参见

#### 引用

[IECHelper 类](b37c3cea-b292-985a-f89b-eb97e38d44fe.htm)

[HslCommunication.Instrument.IEC.Helper 命名空间](acdb00cf-8087-ba54-8ab1-947a92f2b0b4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## IECHelper 字段

[原文連結](http://api.hslcommunication.cn/html/3c7bb720-5c4e-776a-8dba-47d8548c956f.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.IEC.Helper](../html/acdb00cf-8087-ba54-8ab1-947a92f2b0b4.htm "HslCommunication.Instrument.IEC.Helper")

[IECHelper 类](../html/b37c3cea-b292-985a-f89b-eb97e38d44fe.htm "IECHelper 类")

[IECHelper 字段](../html/3c7bb720-5c4e-776a-8dba-47d8548c956f.htm "IECHelper 字段")

[IEC104ControlStartDT 字段](../html/322f0c4b-400f-0657-a33c-0f7a3c961616.htm "IEC104ControlStartDT 字段")

[IEC104ControlStopDT 字段](../html/57c8b03e-dbec-6985-da98-28158a5f61c8.htm "IEC104ControlStopDT 字段")

[IEC104ControlTestFR 字段](../html/e6bc2d0e-63ba-cb18-d33f-dc91af9d966b.htm "IEC104ControlTestFR 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IECHelper 字段 |

[IECHelper](b37c3cea-b292-985a-f89b-eb97e38d44fe.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)字段

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共字段静态成员 | [IEC104ControlStartDT](322f0c4b-400f-0657-a33c-0f7a3c961616.htm) | U帧协议里，启动的功能 |
| 公共字段静态成员 | [IEC104ControlStopDT](57c8b03e-dbec-6985-da98-28158a5f61c8.htm) | U帧协议里，停止的功能 |
| 公共字段静态成员 | [IEC104ControlTestFR](e6bc2d0e-63ba-cb18-d33f-dc91af9d966b.htm) | U帧协议里，测试的功能，主站和子站均可发出 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[IECHelper 类](b37c3cea-b292-985a-f89b-eb97e38d44fe.htm)

[HslCommunication.Instrument.IEC.Helper 命名空间](acdb00cf-8087-ba54-8ab1-947a92f2b0b4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## IEC104ControlStartDT 字段

[原文連結](http://api.hslcommunication.cn/html/322f0c4b-400f-0657-a33c-0f7a3c961616.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.IEC.Helper](../html/acdb00cf-8087-ba54-8ab1-947a92f2b0b4.htm "HslCommunication.Instrument.IEC.Helper")

[IECHelper 类](../html/b37c3cea-b292-985a-f89b-eb97e38d44fe.htm "IECHelper 类")

[IECHelper 字段](../html/3c7bb720-5c4e-776a-8dba-47d8548c956f.htm "IECHelper 字段")

[IEC104ControlStartDT 字段](../html/322f0c4b-400f-0657-a33c-0f7a3c961616.htm "IEC104ControlStartDT 字段")

[IEC104ControlStopDT 字段](../html/57c8b03e-dbec-6985-da98-28158a5f61c8.htm "IEC104ControlStopDT 字段")

[IEC104ControlTestFR 字段](../html/e6bc2d0e-63ba-cb18-d33f-dc91af9d966b.htm "IEC104ControlTestFR 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IECHelperIEC104ControlStartDT 字段 |

U帧协议里，启动的功能

**命名空间：**
 [HslCommunication.Instrument.IEC.Helper](acdb00cf-8087-ba54-8ab1-947a92f2b0b4.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public const byte IEC104ControlStartDT = 7
```

```
Public Const IEC104ControlStartDT As Byte = 7
```

```
public:
literal unsigned char IEC104ControlStartDT = 7
```

```
static val mutable IEC104ControlStartDT: byte
```

#### 字段值

类型：Byte

![](../icons/SectionExpanded.png)参见

#### 引用

[IECHelper 类](b37c3cea-b292-985a-f89b-eb97e38d44fe.htm)

[HslCommunication.Instrument.IEC.Helper 命名空间](acdb00cf-8087-ba54-8ab1-947a92f2b0b4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## IEC104ControlStopDT 字段

[原文連結](http://api.hslcommunication.cn/html/57c8b03e-dbec-6985-da98-28158a5f61c8.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.IEC.Helper](../html/acdb00cf-8087-ba54-8ab1-947a92f2b0b4.htm "HslCommunication.Instrument.IEC.Helper")

[IECHelper 类](../html/b37c3cea-b292-985a-f89b-eb97e38d44fe.htm "IECHelper 类")

[IECHelper 字段](../html/3c7bb720-5c4e-776a-8dba-47d8548c956f.htm "IECHelper 字段")

[IEC104ControlStartDT 字段](../html/322f0c4b-400f-0657-a33c-0f7a3c961616.htm "IEC104ControlStartDT 字段")

[IEC104ControlStopDT 字段](../html/57c8b03e-dbec-6985-da98-28158a5f61c8.htm "IEC104ControlStopDT 字段")

[IEC104ControlTestFR 字段](../html/e6bc2d0e-63ba-cb18-d33f-dc91af9d966b.htm "IEC104ControlTestFR 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IECHelperIEC104ControlStopDT 字段 |

U帧协议里，停止的功能

**命名空间：**
 [HslCommunication.Instrument.IEC.Helper](acdb00cf-8087-ba54-8ab1-947a92f2b0b4.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public const byte IEC104ControlStopDT = 19
```

```
Public Const IEC104ControlStopDT As Byte = 19
```

```
public:
literal unsigned char IEC104ControlStopDT = 19
```

```
static val mutable IEC104ControlStopDT: byte
```

#### 字段值

类型：Byte

![](../icons/SectionExpanded.png)参见

#### 引用

[IECHelper 类](b37c3cea-b292-985a-f89b-eb97e38d44fe.htm)

[HslCommunication.Instrument.IEC.Helper 命名空间](acdb00cf-8087-ba54-8ab1-947a92f2b0b4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## IEC104ControlTestFR 字段

[原文連結](http://api.hslcommunication.cn/html/e6bc2d0e-63ba-cb18-d33f-dc91af9d966b.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.IEC.Helper](../html/acdb00cf-8087-ba54-8ab1-947a92f2b0b4.htm "HslCommunication.Instrument.IEC.Helper")

[IECHelper 类](../html/b37c3cea-b292-985a-f89b-eb97e38d44fe.htm "IECHelper 类")

[IECHelper 字段](../html/3c7bb720-5c4e-776a-8dba-47d8548c956f.htm "IECHelper 字段")

[IEC104ControlStartDT 字段](../html/322f0c4b-400f-0657-a33c-0f7a3c961616.htm "IEC104ControlStartDT 字段")

[IEC104ControlStopDT 字段](../html/57c8b03e-dbec-6985-da98-28158a5f61c8.htm "IEC104ControlStopDT 字段")

[IEC104ControlTestFR 字段](../html/e6bc2d0e-63ba-cb18-d33f-dc91af9d966b.htm "IEC104ControlTestFR 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IECHelperIEC104ControlTestFR 字段 |

U帧协议里，测试的功能，主站和子站均可发出

**命名空间：**
 [HslCommunication.Instrument.IEC.Helper](acdb00cf-8087-ba54-8ab1-947a92f2b0b4.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public const byte IEC104ControlTestFR = 67
```

```
Public Const IEC104ControlTestFR As Byte = 67
```

```
public:
literal unsigned char IEC104ControlTestFR = 67
```

```
static val mutable IEC104ControlTestFR: byte
```

#### 字段值

类型：Byte

![](../icons/SectionExpanded.png)参见

#### 引用

[IECHelper 类](b37c3cea-b292-985a-f89b-eb97e38d44fe.htm)

[HslCommunication.Instrument.IEC.Helper 命名空间](acdb00cf-8087-ba54-8ab1-947a92f2b0b4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## IECSessionInfo 类

[原文連結](http://api.hslcommunication.cn/html/aa145fdf-e2f1-8eab-1375-9d9c8f5a5308.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.IEC.Helper](../html/acdb00cf-8087-ba54-8ab1-947a92f2b0b4.htm "HslCommunication.Instrument.IEC.Helper")

[IECSessionInfo 类](../html/aa145fdf-e2f1-8eab-1375-9d9c8f5a5308.htm "IECSessionInfo 类")

[IECSessionInfo 构造函数](../html/c940a99e-5e7b-140c-50eb-024787b9a16d.htm "IECSessionInfo 构造函数 ")

[IECSessionInfo 属性](../html/14596a48-dcfc-88e5-c405-4428a32a21da.htm "IECSessionInfo 属性")

[IECSessionInfo 方法](../html/a35cdf73-ffbc-c86a-8252-75cee01ce174.htm "IECSessionInfo 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IECSessionInfo 类 |

IEC会话的额外信息，主要包含了这个会话的收发消息ID

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  HslCommunication.Instrument.IEC.HelperIECSessionInfo

**命名空间：**
 [HslCommunication.Instrument.IEC.Helper](acdb00cf-8087-ba54-8ab1-947a92f2b0b4.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class IECSessionInfo
```

```
Public Class IECSessionInfo
```

```
public ref class IECSessionInfo
```

```
type IECSessionInfo =  class end
```

IECSessionInfo 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [IECSessionInfo](c940a99e-5e7b-140c-50eb-024787b9a16d.htm) | 初始化 IECSessionInfo 类的一个新实例 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [RecvMessageID](c59a51c9-9cf0-099e-6a64-06b0d20562cf.htm) | 当前接收的消息ID |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法 | [GetSendMessageID](b024f3ee-2aa7-83ff-462a-a54abf9ab240.htm) | 获取发送的消息ID信息，并进行自增操作 |
| 公共方法 | GetType | (继承自 Object。) |
| 公共方法 | [IncrRecvMessageID](1e8a2869-04b4-3e33-957c-1263a629ab6f.htm) | 自增一个接收到的消息ID |
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

[HslCommunication.Instrument.IEC.Helper 命名空间](acdb00cf-8087-ba54-8ab1-947a92f2b0b4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## IECSessionInfo 构造函数 

[原文連結](http://api.hslcommunication.cn/html/c940a99e-5e7b-140c-50eb-024787b9a16d.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.IEC.Helper](../html/acdb00cf-8087-ba54-8ab1-947a92f2b0b4.htm "HslCommunication.Instrument.IEC.Helper")

[IECSessionInfo 类](../html/aa145fdf-e2f1-8eab-1375-9d9c8f5a5308.htm "IECSessionInfo 类")

[IECSessionInfo 构造函数](../html/c940a99e-5e7b-140c-50eb-024787b9a16d.htm "IECSessionInfo 构造函数 ")

[IECSessionInfo 属性](../html/14596a48-dcfc-88e5-c405-4428a32a21da.htm "IECSessionInfo 属性")

[IECSessionInfo 方法](../html/a35cdf73-ffbc-c86a-8252-75cee01ce174.htm "IECSessionInfo 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IECSessionInfo 构造函数 |

初始化 [IECSessionInfo](aa145fdf-e2f1-8eab-1375-9d9c8f5a5308.htm) 类的一个新实例

**命名空间：**
 [HslCommunication.Instrument.IEC.Helper](acdb00cf-8087-ba54-8ab1-947a92f2b0b4.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public IECSessionInfo()
```

```
Public Sub New
```

```
public:
IECSessionInfo()
```

```
new : unit -> IECSessionInfo
```

![](../icons/SectionExpanded.png)参见

#### 引用

[IECSessionInfo 类](aa145fdf-e2f1-8eab-1375-9d9c8f5a5308.htm)

[HslCommunication.Instrument.IEC.Helper 命名空间](acdb00cf-8087-ba54-8ab1-947a92f2b0b4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## IECSessionInfo 属性

[原文連結](http://api.hslcommunication.cn/html/14596a48-dcfc-88e5-c405-4428a32a21da.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.IEC.Helper](../html/acdb00cf-8087-ba54-8ab1-947a92f2b0b4.htm "HslCommunication.Instrument.IEC.Helper")

[IECSessionInfo 类](../html/aa145fdf-e2f1-8eab-1375-9d9c8f5a5308.htm "IECSessionInfo 类")

[IECSessionInfo 属性](../html/14596a48-dcfc-88e5-c405-4428a32a21da.htm "IECSessionInfo 属性")

[RecvMessageID 属性](../html/c59a51c9-9cf0-099e-6a64-06b0d20562cf.htm "RecvMessageID 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IECSessionInfo 属性 |

[IECSessionInfo](aa145fdf-e2f1-8eab-1375-9d9c8f5a5308.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [RecvMessageID](c59a51c9-9cf0-099e-6a64-06b0d20562cf.htm) | 当前接收的消息ID |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[IECSessionInfo 类](aa145fdf-e2f1-8eab-1375-9d9c8f5a5308.htm)

[HslCommunication.Instrument.IEC.Helper 命名空间](acdb00cf-8087-ba54-8ab1-947a92f2b0b4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## RecvMessageID 属性 

[原文連結](http://api.hslcommunication.cn/html/c59a51c9-9cf0-099e-6a64-06b0d20562cf.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.IEC.Helper](../html/acdb00cf-8087-ba54-8ab1-947a92f2b0b4.htm "HslCommunication.Instrument.IEC.Helper")

[IECSessionInfo 类](../html/aa145fdf-e2f1-8eab-1375-9d9c8f5a5308.htm "IECSessionInfo 类")

[IECSessionInfo 属性](../html/14596a48-dcfc-88e5-c405-4428a32a21da.htm "IECSessionInfo 属性")

[RecvMessageID 属性](../html/c59a51c9-9cf0-099e-6a64-06b0d20562cf.htm "RecvMessageID 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IECSessionInfoRecvMessageID 属性 |

当前接收的消息ID

**命名空间：**
 [HslCommunication.Instrument.IEC.Helper](acdb00cf-8087-ba54-8ab1-947a92f2b0b4.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public int RecvMessageID { get; }
```

```
Public ReadOnly Property RecvMessageID As Integer
	Get
```

```
public:
property int RecvMessageID {
	int get ();
}
```

```
member RecvMessageID : int with get
```

#### 属性值

类型：Int32

![](../icons/SectionExpanded.png)参见

#### 引用

[IECSessionInfo 类](aa145fdf-e2f1-8eab-1375-9d9c8f5a5308.htm)

[HslCommunication.Instrument.IEC.Helper 命名空间](acdb00cf-8087-ba54-8ab1-947a92f2b0b4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## IECSessionInfo 方法

[原文連結](http://api.hslcommunication.cn/html/a35cdf73-ffbc-c86a-8252-75cee01ce174.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.IEC.Helper](../html/acdb00cf-8087-ba54-8ab1-947a92f2b0b4.htm "HslCommunication.Instrument.IEC.Helper")

[IECSessionInfo 类](../html/aa145fdf-e2f1-8eab-1375-9d9c8f5a5308.htm "IECSessionInfo 类")

[IECSessionInfo 方法](../html/a35cdf73-ffbc-c86a-8252-75cee01ce174.htm "IECSessionInfo 方法")

[GetSendMessageID 方法](../html/b024f3ee-2aa7-83ff-462a-a54abf9ab240.htm "GetSendMessageID 方法 ")

[IncrRecvMessageID 方法](../html/1e8a2869-04b4-3e33-957c-1263a629ab6f.htm "IncrRecvMessageID 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IECSessionInfo 方法 |

[IECSessionInfo](aa145fdf-e2f1-8eab-1375-9d9c8f5a5308.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法 | [GetSendMessageID](b024f3ee-2aa7-83ff-462a-a54abf9ab240.htm) | 获取发送的消息ID信息，并进行自增操作 |
| 公共方法 | GetType | (继承自 Object。) |
| 公共方法 | [IncrRecvMessageID](1e8a2869-04b4-3e33-957c-1263a629ab6f.htm) | 自增一个接收到的消息ID |
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

[IECSessionInfo 类](aa145fdf-e2f1-8eab-1375-9d9c8f5a5308.htm)

[HslCommunication.Instrument.IEC.Helper 命名空间](acdb00cf-8087-ba54-8ab1-947a92f2b0b4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GetSendMessageID 方法 

[原文連結](http://api.hslcommunication.cn/html/b024f3ee-2aa7-83ff-462a-a54abf9ab240.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.IEC.Helper](../html/acdb00cf-8087-ba54-8ab1-947a92f2b0b4.htm "HslCommunication.Instrument.IEC.Helper")

[IECSessionInfo 类](../html/aa145fdf-e2f1-8eab-1375-9d9c8f5a5308.htm "IECSessionInfo 类")

[IECSessionInfo 方法](../html/a35cdf73-ffbc-c86a-8252-75cee01ce174.htm "IECSessionInfo 方法")

[GetSendMessageID 方法](../html/b024f3ee-2aa7-83ff-462a-a54abf9ab240.htm "GetSendMessageID 方法 ")

[IncrRecvMessageID 方法](../html/1e8a2869-04b4-3e33-957c-1263a629ab6f.htm "IncrRecvMessageID 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IECSessionInfoGetSendMessageID 方法 |

获取发送的消息ID信息，并进行自增操作

**命名空间：**
 [HslCommunication.Instrument.IEC.Helper](acdb00cf-8087-ba54-8ab1-947a92f2b0b4.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public int GetSendMessageID()
```

```
Public Function GetSendMessageID As Integer
```

```
public:
int GetSendMessageID()
```

```
member GetSendMessageID : unit -> int 
```

#### 返回值

类型：Int32  

[缺少 "M:HslCommunication.Instrument.IEC.Helper.IECSessionInfo.GetSendMessageID" 的 <returns> 文档]

![](../icons/SectionExpanded.png)参见

#### 引用

[IECSessionInfo 类](aa145fdf-e2f1-8eab-1375-9d9c8f5a5308.htm)

[HslCommunication.Instrument.IEC.Helper 命名空间](acdb00cf-8087-ba54-8ab1-947a92f2b0b4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## IncrRecvMessageID 方法 

[原文連結](http://api.hslcommunication.cn/html/1e8a2869-04b4-3e33-957c-1263a629ab6f.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.IEC.Helper](../html/acdb00cf-8087-ba54-8ab1-947a92f2b0b4.htm "HslCommunication.Instrument.IEC.Helper")

[IECSessionInfo 类](../html/aa145fdf-e2f1-8eab-1375-9d9c8f5a5308.htm "IECSessionInfo 类")

[IECSessionInfo 方法](../html/a35cdf73-ffbc-c86a-8252-75cee01ce174.htm "IECSessionInfo 方法")

[GetSendMessageID 方法](../html/b024f3ee-2aa7-83ff-462a-a54abf9ab240.htm "GetSendMessageID 方法 ")

[IncrRecvMessageID 方法](../html/1e8a2869-04b4-3e33-957c-1263a629ab6f.htm "IncrRecvMessageID 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IECSessionInfoIncrRecvMessageID 方法 |

自增一个接收到的消息ID

**命名空间：**
 [HslCommunication.Instrument.IEC.Helper](acdb00cf-8087-ba54-8ab1-947a92f2b0b4.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public int IncrRecvMessageID()
```

```
Public Function IncrRecvMessageID As Integer
```

```
public:
int IncrRecvMessageID()
```

```
member IncrRecvMessageID : unit -> int 
```

#### 返回值

类型：Int32  

[缺少 "M:HslCommunication.Instrument.IEC.Helper.IECSessionInfo.IncrRecvMessageID" 的 <returns> 文档]

![](../icons/SectionExpanded.png)参见

#### 引用

[IECSessionInfo 类](aa145fdf-e2f1-8eab-1375-9d9c8f5a5308.htm)

[HslCommunication.Instrument.IEC.Helper 命名空间](acdb00cf-8087-ba54-8ab1-947a92f2b0b4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)