# HslCommunication - HslCommunication.Profinet.FATEK.Helper

> 分類頁數: 30



---
## HslCommunication.Profinet.FATEK.Helper

[原文連結](http://api.hslcommunication.cn/html/4ac2b39a-b5b1-9027-86a7-eadd294bcffa.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.FATEK.Helper](../html/4ac2b39a-b5b1-9027-86a7-eadd294bcffa.htm "HslCommunication.Profinet.FATEK.Helper")

[FatekProgramHelper 类](../html/5aa26a25-055f-eaea-5db8-af91b52b1e7b.htm "FatekProgramHelper 类")

[IFatekProgram 接口](../html/de0bc2e8-1b3c-afbb-7c20-5f803e24d8fa.htm "IFatekProgram 接口")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslCommunication.Profinet.FATEK.Helper 命名空间 |

[缺少 "N:HslCommunication.Profinet.FATEK.Helper" 的 <summary> 文档]

![](../icons/SectionExpanded.png)类

|  | 类 | 说明 |
| --- | --- | --- |
| 公共类 | [FatekProgramHelper](5aa26a25-055f-eaea-5db8-af91b52b1e7b.htm) | FatekProgram相关的辅助方法，例如报文构建，核心读写支持 |

![](../icons/SectionExpanded.png)接口

|  | 接口 | 说明 |
| --- | --- | --- |
| 公共接口 | [IFatekProgram](de0bc2e8-1b3c-afbb-7c20-5f803e24d8fa.htm) | FatekProgram协议的接口 |

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## FatekProgramHelper 类

[原文連結](http://api.hslcommunication.cn/html/5aa26a25-055f-eaea-5db8-af91b52b1e7b.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.FATEK.Helper](../html/4ac2b39a-b5b1-9027-86a7-eadd294bcffa.htm "HslCommunication.Profinet.FATEK.Helper")

[FatekProgramHelper 类](../html/5aa26a25-055f-eaea-5db8-af91b52b1e7b.htm "FatekProgramHelper 类")

[FatekProgramHelper 构造函数](../html/353097a6-5ae3-89b0-0566-25957c7b6b20.htm "FatekProgramHelper 构造函数 ")

[FatekProgramHelper 方法](../html/481082d5-89ec-7b3f-860f-ced4eec5c714.htm "FatekProgramHelper 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FatekProgramHelper 类 |

FatekProgram相关的辅助方法，例如报文构建，核心读写支持

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  HslCommunication.Profinet.FATEK.HelperFatekProgramHelper

**命名空间：**
 [HslCommunication.Profinet.FATEK.Helper](4ac2b39a-b5b1-9027-86a7-eadd294bcffa.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class FatekProgramHelper
```

```
Public Class FatekProgramHelper
```

```
public ref class FatekProgramHelper
```

```
type FatekProgramHelper =  class end
```

FatekProgramHelper 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [FatekProgramHelper](353097a6-5ae3-89b0-0566-25957c7b6b20.htm) | 初始化 FatekProgramHelper 类的一个新实例 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法静态成员 | [BuildReadBoolCommand](1f0e3930-99a4-b02f-d0b6-e139660f883a.htm) | 创建一条读取的指令信息，需要指定一些参数 |
| 公共方法静态成员 | [BuildReadWordCommand](805ff280-0163-7ccb-25eb-2458424cc1ec.htm) | 创建一条读取的指令信息，需要指定一些参数 |
| 公共方法静态成员 | [BuildWriteBoolCommand](eecc01ad-74fc-3610-4939-9bbf87fc483d.htm) | 创建一条别入bool数据的指令信息，需要指定一些参数 |
| 公共方法静态成员 | [BuildWriteByteCommand](7b615285-6728-1344-384e-bf7cf918abf7.htm) | 创建一条别入byte数据的指令信息，需要指定一些参数，按照字单位 |
| 公共方法静态成员 | [CalculateAcc](49b24ae3-78e3-e4d6-4743-f531b867da87.htm) | 计算指令的和校验码 |
| 公共方法静态成员 | [CheckReceiveDataComplete](b200e06c-38b7-9656-a340-54aa2373bb0a.htm) | 检查当前的串口的数据接收是否完整 |
| 公共方法静态成员 | [CheckResponse](1b4e1f3e-399f-4381-f912-86708faee605.htm) | 检查PLC反馈的报文是否正确，如果不正确，返回错误消息 |
| 公共方法 | Equals | (继承自 Object。) |
| 公共方法静态成员 | [ExtraResponse](ad08cd70-4590-fc5d-a80d-a445450f9631.htm) | 提取当前的结果数据信息，针对的是字单位的方式 |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法静态成员 | [GetErrorDescriptionFromCode](fca0a278-c68e-cfe1-fcf4-816efcc878cb.htm) | 根据错误码获取到真实的文本信息 |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法静态成员 | [PackFatekCommand](fc22707e-222f-8fec-8d57-83742d4f17c1.htm) | 将Fatek的基本命令打包成可以发送PLC的电文消息 |
| 公共方法静态成员 | [Read](cc1b8d93-c06c-41b9-748f-aa824b5f223b.htm) | 批量读取PLC的字节数据，以字为单位，支持读取X,Y,M,S,D,T,C,R,RT,RC具体的地址范围需要根据PLC型号来确认，地址可以携带站号信息，例如 s=2;D100  Read PLC byte data in batches, in word units. Supports reading X, Y, M, S, D, T, C, R, RT, RC. The specific address range needs to be confirmed according to the PLC model, The address can carry station number information, such as s=2;D100 |
| 公共方法静态成员 | [ReadAsync](c1a5498b-193d-f3bb-3e3c-7ab9585164ed.htm) | 批量读取PLC的字节数据，以字为单位，支持读取X,Y,M,S,D,T,C,R,RT,RC具体的地址范围需要根据PLC型号来确认，地址可以携带站号信息，例如 s=2;D100  Read PLC byte data in batches, in word units. Supports reading X, Y, M, S, D, T, C, R, RT, RC. The specific address range needs to be confirmed according to the PLC model, The address can carry station number information, such as s=2;D100 |
| 公共方法静态成员 | [ReadBool](5f136473-54c4-ce95-34f0-39c1028372f2.htm) | 批量读取bool类型数据，支持的类型为X,Y,M,S,T,C，具体的地址范围取决于PLC的类型，地址可以携带站号信息，例如 s=2;M100  Read bool data in batches. The supported types are X, Y, M, S, T, C. The specific address range depends on the type of PLC, The address can carry station number information, such as s=2;M100 |
| 公共方法静态成员 | [ReadBoolAsync](997a3300-c2ad-8816-bdfc-a2e0f9532bb4.htm) | 批量读取bool类型数据，支持的类型为X,Y,M,S,T,C，具体的地址范围取决于PLC的类型，地址可以携带站号信息，例如 s=2;M100  Read bool data in batches. The supported types are X, Y, M, S, T, C. The specific address range depends on the type of PLC, The address can carry station number information, such as s=2;M100 |
| 公共方法静态成员 | [ReadStatus](a73fd644-bf59-ca9a-d0ff-45d177c75e64.htm) | 读取当前PLC的状态信息，返回一个bool数组，同时包含了几种电量信息，分别为 0: RUN/STOP, 1: BAT LOW/正常, 2: Ladder checksum error/正常, 3: 使用ROM PACK/未使用, 4: WDT Timeout/正常, 5: 设定ID/未设ID， 6： 紧急停机/正常  Read the status information of the current PLC and return a bool array, which also contains several power information, 0: RUN/STOP, 1: BAT LOW/normal, 2: Ladder checksum error/normal, 3: Use ROM PACK/ Not used, 4: WDT Timeout/Normal, 5: ID set/ID not set, 6: Emergency stop/Normal |
| 公共方法静态成员 | [ReadStatusAsync](62d31aab-0b53-ce2f-ffa5-09b3e4f85675.htm) | 使PLC处于STOP状态 |
| 公共方法静态成员 | [Run](d19f36a9-9645-fc5f-9bf8-7de6a709ed47.htm) | 使PLC处于RUN的状态 |
| 公共方法静态成员 | [RunAsync](0c1b1b0e-cecd-d2a2-5fb1-46a8a0321c39.htm) | 使PLC处于RUN的状态 |
| 公共方法静态成员 | [Stop](2fc1d7f4-d7f2-ee5a-b271-69c33c1c3479.htm) | 使PLC处于STOP状态 |
| 公共方法静态成员 | [StopAsync](cfee21d4-ecbd-ea41-89f0-03d4f71102c7.htm) | 使PLC处于STOP状态 |
| 公共方法 | ToString | (继承自 Object。) |
| 公共方法静态成员 | [Write(IReadWriteDevice, Byte, String, Boolean)](dfabf156-9540-d08c-da8c-954c185417ed.htm) | 批量写入bool类型的数组，支持的类型为X,Y,M,S,T,C，具体的地址范围取决于PLC的类型，地址可以携带站号信息，例如 s=2;M100  Write arrays of type bool in batches. The supported types are X, Y, M, S, T, C. The specific address range depends on the type of PLC, The address can carry station number information, such as s=2;M100 |
| 公共方法静态成员 | [Write(IReadWriteDevice, Byte, String, Byte)](1d873a33-51ee-e065-1aac-efeb6b37a452.htm) | 批量写入PLC的数据，以字为单位，也就是说最少2个字节信息，支持X,Y,M,S,D,T,C,R,RT,RC具体的地址范围需要根据PLC型号来确认，地址可以携带站号信息，例如 s=2;D100  The data written to the PLC in batches, in units of words, that is, at least 2 bytes of information, supporting X, Y, M, S, D, T, C, R, RT, and RC. The specific address range needs to be based on the PLC model To confirm, The address can carry station number information, such as s=2;D100 |
| 公共方法静态成员 | [WriteAsync(IReadWriteDevice, Byte, String, Boolean)](e415bcb2-0cbc-6f93-151a-93d96112f9d9.htm) | 批量写入bool类型的数组，支持的类型为X,Y,M,S,T,C，具体的地址范围取决于PLC的类型，地址可以携带站号信息，例如 s=2;M100  Write arrays of type bool in batches. The supported types are X, Y, M, S, T, C. The specific address range depends on the type of PLC, The address can carry station number information, such as s=2;M100 |
| 公共方法静态成员 | [WriteAsync(IReadWriteDevice, Byte, String, Byte)](3754fd86-bb2c-0ca3-e736-279e77905479.htm) | 批量写入PLC的数据，以字为单位，也就是说最少2个字节信息，支持X,Y,M,S,D,T,C,R,RT,RC具体的地址范围需要根据PLC型号来确认，地址可以携带站号信息，例如 s=2;D100  The data written to the PLC in batches, in units of words, that is, at least 2 bytes of information, supporting X, Y, M, S, D, T, C, R, RT, and RC. The specific address range needs to be based on the PLC model To confirm, The address can carry station number information, such as s=2;D100 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.Profinet.FATEK.Helper 命名空间](4ac2b39a-b5b1-9027-86a7-eadd294bcffa.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## FatekProgramHelper 构造函数 

[原文連結](http://api.hslcommunication.cn/html/353097a6-5ae3-89b0-0566-25957c7b6b20.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.FATEK.Helper](../html/4ac2b39a-b5b1-9027-86a7-eadd294bcffa.htm "HslCommunication.Profinet.FATEK.Helper")

[FatekProgramHelper 类](../html/5aa26a25-055f-eaea-5db8-af91b52b1e7b.htm "FatekProgramHelper 类")

[FatekProgramHelper 构造函数](../html/353097a6-5ae3-89b0-0566-25957c7b6b20.htm "FatekProgramHelper 构造函数 ")

[FatekProgramHelper 方法](../html/481082d5-89ec-7b3f-860f-ced4eec5c714.htm "FatekProgramHelper 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FatekProgramHelper 构造函数 |

初始化 [FatekProgramHelper](5aa26a25-055f-eaea-5db8-af91b52b1e7b.htm) 类的一个新实例

**命名空间：**
 [HslCommunication.Profinet.FATEK.Helper](4ac2b39a-b5b1-9027-86a7-eadd294bcffa.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public FatekProgramHelper()
```

```
Public Sub New
```

```
public:
FatekProgramHelper()
```

```
new : unit -> FatekProgramHelper
```

![](../icons/SectionExpanded.png)参见

#### 引用

[FatekProgramHelper 类](5aa26a25-055f-eaea-5db8-af91b52b1e7b.htm)

[HslCommunication.Profinet.FATEK.Helper 命名空间](4ac2b39a-b5b1-9027-86a7-eadd294bcffa.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## FatekProgramHelper 方法

[原文連結](http://api.hslcommunication.cn/html/481082d5-89ec-7b3f-860f-ced4eec5c714.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.FATEK.Helper](../html/4ac2b39a-b5b1-9027-86a7-eadd294bcffa.htm "HslCommunication.Profinet.FATEK.Helper")

[FatekProgramHelper 类](../html/5aa26a25-055f-eaea-5db8-af91b52b1e7b.htm "FatekProgramHelper 类")

[FatekProgramHelper 方法](../html/481082d5-89ec-7b3f-860f-ced4eec5c714.htm "FatekProgramHelper 方法")

[BuildReadBoolCommand 方法](../html/1f0e3930-99a4-b02f-d0b6-e139660f883a.htm "BuildReadBoolCommand 方法 ")

[BuildReadWordCommand 方法](../html/805ff280-0163-7ccb-25eb-2458424cc1ec.htm "BuildReadWordCommand 方法 ")

[BuildWriteBoolCommand 方法](../html/eecc01ad-74fc-3610-4939-9bbf87fc483d.htm "BuildWriteBoolCommand 方法 ")

[BuildWriteByteCommand 方法](../html/7b615285-6728-1344-384e-bf7cf918abf7.htm "BuildWriteByteCommand 方法 ")

[CalculateAcc 方法](../html/49b24ae3-78e3-e4d6-4743-f531b867da87.htm "CalculateAcc 方法 ")

[CheckReceiveDataComplete 方法](../html/b200e06c-38b7-9656-a340-54aa2373bb0a.htm "CheckReceiveDataComplete 方法 ")

[CheckResponse 方法](../html/1b4e1f3e-399f-4381-f912-86708faee605.htm "CheckResponse 方法 ")

[ExtraResponse 方法](../html/ad08cd70-4590-fc5d-a80d-a445450f9631.htm "ExtraResponse 方法 ")

[GetErrorDescriptionFromCode 方法](../html/fca0a278-c68e-cfe1-fcf4-816efcc878cb.htm "GetErrorDescriptionFromCode 方法 ")

[PackFatekCommand 方法](../html/fc22707e-222f-8fec-8d57-83742d4f17c1.htm "PackFatekCommand 方法 ")

[Read 方法](../html/cc1b8d93-c06c-41b9-748f-aa824b5f223b.htm "Read 方法 ")

[ReadAsync 方法](../html/c1a5498b-193d-f3bb-3e3c-7ab9585164ed.htm "ReadAsync 方法 ")

[ReadBool 方法](../html/5f136473-54c4-ce95-34f0-39c1028372f2.htm "ReadBool 方法 ")

[ReadBoolAsync 方法](../html/997a3300-c2ad-8816-bdfc-a2e0f9532bb4.htm "ReadBoolAsync 方法 ")

[ReadStatus 方法](../html/a73fd644-bf59-ca9a-d0ff-45d177c75e64.htm "ReadStatus 方法 ")

[ReadStatusAsync 方法](../html/62d31aab-0b53-ce2f-ffa5-09b3e4f85675.htm "ReadStatusAsync 方法 ")

[Run 方法](../html/d19f36a9-9645-fc5f-9bf8-7de6a709ed47.htm "Run 方法 ")

[RunAsync 方法](../html/0c1b1b0e-cecd-d2a2-5fb1-46a8a0321c39.htm "RunAsync 方法 ")

[Stop 方法](../html/2fc1d7f4-d7f2-ee5a-b271-69c33c1c3479.htm "Stop 方法 ")

[StopAsync 方法](../html/cfee21d4-ecbd-ea41-89f0-03d4f71102c7.htm "StopAsync 方法 ")

[Write 方法](../html/5adce02f-a40e-f29a-98db-23fb5791b508.htm "Write 方法 ")

[WriteAsync 方法](../html/4b5310b6-d2c7-d1ac-d43d-0c29bbbfea3b.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FatekProgramHelper 方法 |

[FatekProgramHelper](5aa26a25-055f-eaea-5db8-af91b52b1e7b.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法静态成员 | [BuildReadBoolCommand](1f0e3930-99a4-b02f-d0b6-e139660f883a.htm) | 创建一条读取的指令信息，需要指定一些参数 |
| 公共方法静态成员 | [BuildReadWordCommand](805ff280-0163-7ccb-25eb-2458424cc1ec.htm) | 创建一条读取的指令信息，需要指定一些参数 |
| 公共方法静态成员 | [BuildWriteBoolCommand](eecc01ad-74fc-3610-4939-9bbf87fc483d.htm) | 创建一条别入bool数据的指令信息，需要指定一些参数 |
| 公共方法静态成员 | [BuildWriteByteCommand](7b615285-6728-1344-384e-bf7cf918abf7.htm) | 创建一条别入byte数据的指令信息，需要指定一些参数，按照字单位 |
| 公共方法静态成员 | [CalculateAcc](49b24ae3-78e3-e4d6-4743-f531b867da87.htm) | 计算指令的和校验码 |
| 公共方法静态成员 | [CheckReceiveDataComplete](b200e06c-38b7-9656-a340-54aa2373bb0a.htm) | 检查当前的串口的数据接收是否完整 |
| 公共方法静态成员 | [CheckResponse](1b4e1f3e-399f-4381-f912-86708faee605.htm) | 检查PLC反馈的报文是否正确，如果不正确，返回错误消息 |
| 公共方法 | Equals | (继承自 Object。) |
| 公共方法静态成员 | [ExtraResponse](ad08cd70-4590-fc5d-a80d-a445450f9631.htm) | 提取当前的结果数据信息，针对的是字单位的方式 |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法静态成员 | [GetErrorDescriptionFromCode](fca0a278-c68e-cfe1-fcf4-816efcc878cb.htm) | 根据错误码获取到真实的文本信息 |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法静态成员 | [PackFatekCommand](fc22707e-222f-8fec-8d57-83742d4f17c1.htm) | 将Fatek的基本命令打包成可以发送PLC的电文消息 |
| 公共方法静态成员 | [Read](cc1b8d93-c06c-41b9-748f-aa824b5f223b.htm) | 批量读取PLC的字节数据，以字为单位，支持读取X,Y,M,S,D,T,C,R,RT,RC具体的地址范围需要根据PLC型号来确认，地址可以携带站号信息，例如 s=2;D100  Read PLC byte data in batches, in word units. Supports reading X, Y, M, S, D, T, C, R, RT, RC. The specific address range needs to be confirmed according to the PLC model, The address can carry station number information, such as s=2;D100 |
| 公共方法静态成员 | [ReadAsync](c1a5498b-193d-f3bb-3e3c-7ab9585164ed.htm) | 批量读取PLC的字节数据，以字为单位，支持读取X,Y,M,S,D,T,C,R,RT,RC具体的地址范围需要根据PLC型号来确认，地址可以携带站号信息，例如 s=2;D100  Read PLC byte data in batches, in word units. Supports reading X, Y, M, S, D, T, C, R, RT, RC. The specific address range needs to be confirmed according to the PLC model, The address can carry station number information, such as s=2;D100 |
| 公共方法静态成员 | [ReadBool](5f136473-54c4-ce95-34f0-39c1028372f2.htm) | 批量读取bool类型数据，支持的类型为X,Y,M,S,T,C，具体的地址范围取决于PLC的类型，地址可以携带站号信息，例如 s=2;M100  Read bool data in batches. The supported types are X, Y, M, S, T, C. The specific address range depends on the type of PLC, The address can carry station number information, such as s=2;M100 |
| 公共方法静态成员 | [ReadBoolAsync](997a3300-c2ad-8816-bdfc-a2e0f9532bb4.htm) | 批量读取bool类型数据，支持的类型为X,Y,M,S,T,C，具体的地址范围取决于PLC的类型，地址可以携带站号信息，例如 s=2;M100  Read bool data in batches. The supported types are X, Y, M, S, T, C. The specific address range depends on the type of PLC, The address can carry station number information, such as s=2;M100 |
| 公共方法静态成员 | [ReadStatus](a73fd644-bf59-ca9a-d0ff-45d177c75e64.htm) | 读取当前PLC的状态信息，返回一个bool数组，同时包含了几种电量信息，分别为 0: RUN/STOP, 1: BAT LOW/正常, 2: Ladder checksum error/正常, 3: 使用ROM PACK/未使用, 4: WDT Timeout/正常, 5: 设定ID/未设ID， 6： 紧急停机/正常  Read the status information of the current PLC and return a bool array, which also contains several power information, 0: RUN/STOP, 1: BAT LOW/normal, 2: Ladder checksum error/normal, 3: Use ROM PACK/ Not used, 4: WDT Timeout/Normal, 5: ID set/ID not set, 6: Emergency stop/Normal |
| 公共方法静态成员 | [ReadStatusAsync](62d31aab-0b53-ce2f-ffa5-09b3e4f85675.htm) | 使PLC处于STOP状态 |
| 公共方法静态成员 | [Run](d19f36a9-9645-fc5f-9bf8-7de6a709ed47.htm) | 使PLC处于RUN的状态 |
| 公共方法静态成员 | [RunAsync](0c1b1b0e-cecd-d2a2-5fb1-46a8a0321c39.htm) | 使PLC处于RUN的状态 |
| 公共方法静态成员 | [Stop](2fc1d7f4-d7f2-ee5a-b271-69c33c1c3479.htm) | 使PLC处于STOP状态 |
| 公共方法静态成员 | [StopAsync](cfee21d4-ecbd-ea41-89f0-03d4f71102c7.htm) | 使PLC处于STOP状态 |
| 公共方法 | ToString | (继承自 Object。) |
| 公共方法静态成员 | [Write(IReadWriteDevice, Byte, String, Boolean)](dfabf156-9540-d08c-da8c-954c185417ed.htm) | 批量写入bool类型的数组，支持的类型为X,Y,M,S,T,C，具体的地址范围取决于PLC的类型，地址可以携带站号信息，例如 s=2;M100  Write arrays of type bool in batches. The supported types are X, Y, M, S, T, C. The specific address range depends on the type of PLC, The address can carry station number information, such as s=2;M100 |
| 公共方法静态成员 | [Write(IReadWriteDevice, Byte, String, Byte)](1d873a33-51ee-e065-1aac-efeb6b37a452.htm) | 批量写入PLC的数据，以字为单位，也就是说最少2个字节信息，支持X,Y,M,S,D,T,C,R,RT,RC具体的地址范围需要根据PLC型号来确认，地址可以携带站号信息，例如 s=2;D100  The data written to the PLC in batches, in units of words, that is, at least 2 bytes of information, supporting X, Y, M, S, D, T, C, R, RT, and RC. The specific address range needs to be based on the PLC model To confirm, The address can carry station number information, such as s=2;D100 |
| 公共方法静态成员 | [WriteAsync(IReadWriteDevice, Byte, String, Boolean)](e415bcb2-0cbc-6f93-151a-93d96112f9d9.htm) | 批量写入bool类型的数组，支持的类型为X,Y,M,S,T,C，具体的地址范围取决于PLC的类型，地址可以携带站号信息，例如 s=2;M100  Write arrays of type bool in batches. The supported types are X, Y, M, S, T, C. The specific address range depends on the type of PLC, The address can carry station number information, such as s=2;M100 |
| 公共方法静态成员 | [WriteAsync(IReadWriteDevice, Byte, String, Byte)](3754fd86-bb2c-0ca3-e736-279e77905479.htm) | 批量写入PLC的数据，以字为单位，也就是说最少2个字节信息，支持X,Y,M,S,D,T,C,R,RT,RC具体的地址范围需要根据PLC型号来确认，地址可以携带站号信息，例如 s=2;D100  The data written to the PLC in batches, in units of words, that is, at least 2 bytes of information, supporting X, Y, M, S, D, T, C, R, RT, and RC. The specific address range needs to be based on the PLC model To confirm, The address can carry station number information, such as s=2;D100 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[FatekProgramHelper 类](5aa26a25-055f-eaea-5db8-af91b52b1e7b.htm)

[HslCommunication.Profinet.FATEK.Helper 命名空间](4ac2b39a-b5b1-9027-86a7-eadd294bcffa.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## BuildReadBoolCommand 方法 

[原文連結](http://api.hslcommunication.cn/html/1f0e3930-99a4-b02f-d0b6-e139660f883a.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.FATEK.Helper](../html/4ac2b39a-b5b1-9027-86a7-eadd294bcffa.htm "HslCommunication.Profinet.FATEK.Helper")

[FatekProgramHelper 类](../html/5aa26a25-055f-eaea-5db8-af91b52b1e7b.htm "FatekProgramHelper 类")

[FatekProgramHelper 方法](../html/481082d5-89ec-7b3f-860f-ced4eec5c714.htm "FatekProgramHelper 方法")

[BuildReadBoolCommand 方法](../html/1f0e3930-99a4-b02f-d0b6-e139660f883a.htm "BuildReadBoolCommand 方法 ")

[BuildReadWordCommand 方法](../html/805ff280-0163-7ccb-25eb-2458424cc1ec.htm "BuildReadWordCommand 方法 ")

[BuildWriteBoolCommand 方法](../html/eecc01ad-74fc-3610-4939-9bbf87fc483d.htm "BuildWriteBoolCommand 方法 ")

[BuildWriteByteCommand 方法](../html/7b615285-6728-1344-384e-bf7cf918abf7.htm "BuildWriteByteCommand 方法 ")

[CalculateAcc 方法](../html/49b24ae3-78e3-e4d6-4743-f531b867da87.htm "CalculateAcc 方法 ")

[CheckReceiveDataComplete 方法](../html/b200e06c-38b7-9656-a340-54aa2373bb0a.htm "CheckReceiveDataComplete 方法 ")

[CheckResponse 方法](../html/1b4e1f3e-399f-4381-f912-86708faee605.htm "CheckResponse 方法 ")

[ExtraResponse 方法](../html/ad08cd70-4590-fc5d-a80d-a445450f9631.htm "ExtraResponse 方法 ")

[GetErrorDescriptionFromCode 方法](../html/fca0a278-c68e-cfe1-fcf4-816efcc878cb.htm "GetErrorDescriptionFromCode 方法 ")

[PackFatekCommand 方法](../html/fc22707e-222f-8fec-8d57-83742d4f17c1.htm "PackFatekCommand 方法 ")

[Read 方法](../html/cc1b8d93-c06c-41b9-748f-aa824b5f223b.htm "Read 方法 ")

[ReadAsync 方法](../html/c1a5498b-193d-f3bb-3e3c-7ab9585164ed.htm "ReadAsync 方法 ")

[ReadBool 方法](../html/5f136473-54c4-ce95-34f0-39c1028372f2.htm "ReadBool 方法 ")

[ReadBoolAsync 方法](../html/997a3300-c2ad-8816-bdfc-a2e0f9532bb4.htm "ReadBoolAsync 方法 ")

[ReadStatus 方法](../html/a73fd644-bf59-ca9a-d0ff-45d177c75e64.htm "ReadStatus 方法 ")

[ReadStatusAsync 方法](../html/62d31aab-0b53-ce2f-ffa5-09b3e4f85675.htm "ReadStatusAsync 方法 ")

[Run 方法](../html/d19f36a9-9645-fc5f-9bf8-7de6a709ed47.htm "Run 方法 ")

[RunAsync 方法](../html/0c1b1b0e-cecd-d2a2-5fb1-46a8a0321c39.htm "RunAsync 方法 ")

[Stop 方法](../html/2fc1d7f4-d7f2-ee5a-b271-69c33c1c3479.htm "Stop 方法 ")

[StopAsync 方法](../html/cfee21d4-ecbd-ea41-89f0-03d4f71102c7.htm "StopAsync 方法 ")

[Write 方法](../html/5adce02f-a40e-f29a-98db-23fb5791b508.htm "Write 方法 ")

[WriteAsync 方法](../html/4b5310b6-d2c7-d1ac-d43d-0c29bbbfea3b.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FatekProgramHelperBuildReadBoolCommand 方法 |

创建一条读取的指令信息，需要指定一些参数

**命名空间：**
 [HslCommunication.Profinet.FATEK.Helper](4ac2b39a-b5b1-9027-86a7-eadd294bcffa.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult<List<byte[]>> BuildReadBoolCommand(
	byte station,
	string address,
	ushort length
)
```

```
Public Shared Function BuildReadBoolCommand ( 
	station As Byte,
	address As String,
	length As UShort
) As OperateResult(Of List(Of Byte()))
```

```
public:
static OperateResult<List<array<unsigned char>^>^>^ BuildReadBoolCommand(
	unsigned char station, 
	String^ address, 
	unsigned short length
)
```

```
static member BuildReadBoolCommand : 
        station : byte * 
        address : string * 
        length : uint16 -> OperateResult<List<byte[]>> 
```

#### 参数

station
:   类型：SystemByte  
    PLC的站号

address
:   类型：SystemString  
    地址信息

length
:   类型：SystemUInt16  
    数据长度

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)ListByte  
是否成功的结果对象

![](../icons/SectionExpanded.png)参见

#### 引用

[FatekProgramHelper 类](5aa26a25-055f-eaea-5db8-af91b52b1e7b.htm)

[HslCommunication.Profinet.FATEK.Helper 命名空间](4ac2b39a-b5b1-9027-86a7-eadd294bcffa.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## BuildReadWordCommand 方法 

[原文連結](http://api.hslcommunication.cn/html/805ff280-0163-7ccb-25eb-2458424cc1ec.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.FATEK.Helper](../html/4ac2b39a-b5b1-9027-86a7-eadd294bcffa.htm "HslCommunication.Profinet.FATEK.Helper")

[FatekProgramHelper 类](../html/5aa26a25-055f-eaea-5db8-af91b52b1e7b.htm "FatekProgramHelper 类")

[FatekProgramHelper 方法](../html/481082d5-89ec-7b3f-860f-ced4eec5c714.htm "FatekProgramHelper 方法")

[BuildReadBoolCommand 方法](../html/1f0e3930-99a4-b02f-d0b6-e139660f883a.htm "BuildReadBoolCommand 方法 ")

[BuildReadWordCommand 方法](../html/805ff280-0163-7ccb-25eb-2458424cc1ec.htm "BuildReadWordCommand 方法 ")

[BuildWriteBoolCommand 方法](../html/eecc01ad-74fc-3610-4939-9bbf87fc483d.htm "BuildWriteBoolCommand 方法 ")

[BuildWriteByteCommand 方法](../html/7b615285-6728-1344-384e-bf7cf918abf7.htm "BuildWriteByteCommand 方法 ")

[CalculateAcc 方法](../html/49b24ae3-78e3-e4d6-4743-f531b867da87.htm "CalculateAcc 方法 ")

[CheckReceiveDataComplete 方法](../html/b200e06c-38b7-9656-a340-54aa2373bb0a.htm "CheckReceiveDataComplete 方法 ")

[CheckResponse 方法](../html/1b4e1f3e-399f-4381-f912-86708faee605.htm "CheckResponse 方法 ")

[ExtraResponse 方法](../html/ad08cd70-4590-fc5d-a80d-a445450f9631.htm "ExtraResponse 方法 ")

[GetErrorDescriptionFromCode 方法](../html/fca0a278-c68e-cfe1-fcf4-816efcc878cb.htm "GetErrorDescriptionFromCode 方法 ")

[PackFatekCommand 方法](../html/fc22707e-222f-8fec-8d57-83742d4f17c1.htm "PackFatekCommand 方法 ")

[Read 方法](../html/cc1b8d93-c06c-41b9-748f-aa824b5f223b.htm "Read 方法 ")

[ReadAsync 方法](../html/c1a5498b-193d-f3bb-3e3c-7ab9585164ed.htm "ReadAsync 方法 ")

[ReadBool 方法](../html/5f136473-54c4-ce95-34f0-39c1028372f2.htm "ReadBool 方法 ")

[ReadBoolAsync 方法](../html/997a3300-c2ad-8816-bdfc-a2e0f9532bb4.htm "ReadBoolAsync 方法 ")

[ReadStatus 方法](../html/a73fd644-bf59-ca9a-d0ff-45d177c75e64.htm "ReadStatus 方法 ")

[ReadStatusAsync 方法](../html/62d31aab-0b53-ce2f-ffa5-09b3e4f85675.htm "ReadStatusAsync 方法 ")

[Run 方法](../html/d19f36a9-9645-fc5f-9bf8-7de6a709ed47.htm "Run 方法 ")

[RunAsync 方法](../html/0c1b1b0e-cecd-d2a2-5fb1-46a8a0321c39.htm "RunAsync 方法 ")

[Stop 方法](../html/2fc1d7f4-d7f2-ee5a-b271-69c33c1c3479.htm "Stop 方法 ")

[StopAsync 方法](../html/cfee21d4-ecbd-ea41-89f0-03d4f71102c7.htm "StopAsync 方法 ")

[Write 方法](../html/5adce02f-a40e-f29a-98db-23fb5791b508.htm "Write 方法 ")

[WriteAsync 方法](../html/4b5310b6-d2c7-d1ac-d43d-0c29bbbfea3b.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FatekProgramHelperBuildReadWordCommand 方法 |

创建一条读取的指令信息，需要指定一些参数

**命名空间：**
 [HslCommunication.Profinet.FATEK.Helper](4ac2b39a-b5b1-9027-86a7-eadd294bcffa.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult<List<byte[]>> BuildReadWordCommand(
	byte station,
	string address,
	ushort length
)
```

```
Public Shared Function BuildReadWordCommand ( 
	station As Byte,
	address As String,
	length As UShort
) As OperateResult(Of List(Of Byte()))
```

```
public:
static OperateResult<List<array<unsigned char>^>^>^ BuildReadWordCommand(
	unsigned char station, 
	String^ address, 
	unsigned short length
)
```

```
static member BuildReadWordCommand : 
        station : byte * 
        address : string * 
        length : uint16 -> OperateResult<List<byte[]>> 
```

#### 参数

station
:   类型：SystemByte  
    PLC的站号

address
:   类型：SystemString  
    地址信息

length
:   类型：SystemUInt16  
    数据长度

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)ListByte  
是否成功的结果对象

![](../icons/SectionExpanded.png)参见

#### 引用

[FatekProgramHelper 类](5aa26a25-055f-eaea-5db8-af91b52b1e7b.htm)

[HslCommunication.Profinet.FATEK.Helper 命名空间](4ac2b39a-b5b1-9027-86a7-eadd294bcffa.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## BuildWriteBoolCommand 方法 

[原文連結](http://api.hslcommunication.cn/html/eecc01ad-74fc-3610-4939-9bbf87fc483d.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.FATEK.Helper](../html/4ac2b39a-b5b1-9027-86a7-eadd294bcffa.htm "HslCommunication.Profinet.FATEK.Helper")

[FatekProgramHelper 类](../html/5aa26a25-055f-eaea-5db8-af91b52b1e7b.htm "FatekProgramHelper 类")

[FatekProgramHelper 方法](../html/481082d5-89ec-7b3f-860f-ced4eec5c714.htm "FatekProgramHelper 方法")

[BuildReadBoolCommand 方法](../html/1f0e3930-99a4-b02f-d0b6-e139660f883a.htm "BuildReadBoolCommand 方法 ")

[BuildReadWordCommand 方法](../html/805ff280-0163-7ccb-25eb-2458424cc1ec.htm "BuildReadWordCommand 方法 ")

[BuildWriteBoolCommand 方法](../html/eecc01ad-74fc-3610-4939-9bbf87fc483d.htm "BuildWriteBoolCommand 方法 ")

[BuildWriteByteCommand 方法](../html/7b615285-6728-1344-384e-bf7cf918abf7.htm "BuildWriteByteCommand 方法 ")

[CalculateAcc 方法](../html/49b24ae3-78e3-e4d6-4743-f531b867da87.htm "CalculateAcc 方法 ")

[CheckReceiveDataComplete 方法](../html/b200e06c-38b7-9656-a340-54aa2373bb0a.htm "CheckReceiveDataComplete 方法 ")

[CheckResponse 方法](../html/1b4e1f3e-399f-4381-f912-86708faee605.htm "CheckResponse 方法 ")

[ExtraResponse 方法](../html/ad08cd70-4590-fc5d-a80d-a445450f9631.htm "ExtraResponse 方法 ")

[GetErrorDescriptionFromCode 方法](../html/fca0a278-c68e-cfe1-fcf4-816efcc878cb.htm "GetErrorDescriptionFromCode 方法 ")

[PackFatekCommand 方法](../html/fc22707e-222f-8fec-8d57-83742d4f17c1.htm "PackFatekCommand 方法 ")

[Read 方法](../html/cc1b8d93-c06c-41b9-748f-aa824b5f223b.htm "Read 方法 ")

[ReadAsync 方法](../html/c1a5498b-193d-f3bb-3e3c-7ab9585164ed.htm "ReadAsync 方法 ")

[ReadBool 方法](../html/5f136473-54c4-ce95-34f0-39c1028372f2.htm "ReadBool 方法 ")

[ReadBoolAsync 方法](../html/997a3300-c2ad-8816-bdfc-a2e0f9532bb4.htm "ReadBoolAsync 方法 ")

[ReadStatus 方法](../html/a73fd644-bf59-ca9a-d0ff-45d177c75e64.htm "ReadStatus 方法 ")

[ReadStatusAsync 方法](../html/62d31aab-0b53-ce2f-ffa5-09b3e4f85675.htm "ReadStatusAsync 方法 ")

[Run 方法](../html/d19f36a9-9645-fc5f-9bf8-7de6a709ed47.htm "Run 方法 ")

[RunAsync 方法](../html/0c1b1b0e-cecd-d2a2-5fb1-46a8a0321c39.htm "RunAsync 方法 ")

[Stop 方法](../html/2fc1d7f4-d7f2-ee5a-b271-69c33c1c3479.htm "Stop 方法 ")

[StopAsync 方法](../html/cfee21d4-ecbd-ea41-89f0-03d4f71102c7.htm "StopAsync 方法 ")

[Write 方法](../html/5adce02f-a40e-f29a-98db-23fb5791b508.htm "Write 方法 ")

[WriteAsync 方法](../html/4b5310b6-d2c7-d1ac-d43d-0c29bbbfea3b.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FatekProgramHelperBuildWriteBoolCommand 方法 |

创建一条别入bool数据的指令信息，需要指定一些参数

**命名空间：**
 [HslCommunication.Profinet.FATEK.Helper](4ac2b39a-b5b1-9027-86a7-eadd294bcffa.htm)  
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
    站号

address
:   类型：SystemString  
    地址

value
:   类型：SystemBoolean  
    数组值

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
是否创建成功

![](../icons/SectionExpanded.png)参见

#### 引用

[FatekProgramHelper 类](5aa26a25-055f-eaea-5db8-af91b52b1e7b.htm)

[HslCommunication.Profinet.FATEK.Helper 命名空间](4ac2b39a-b5b1-9027-86a7-eadd294bcffa.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## BuildWriteByteCommand 方法 

[原文連結](http://api.hslcommunication.cn/html/7b615285-6728-1344-384e-bf7cf918abf7.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.FATEK.Helper](../html/4ac2b39a-b5b1-9027-86a7-eadd294bcffa.htm "HslCommunication.Profinet.FATEK.Helper")

[FatekProgramHelper 类](../html/5aa26a25-055f-eaea-5db8-af91b52b1e7b.htm "FatekProgramHelper 类")

[FatekProgramHelper 方法](../html/481082d5-89ec-7b3f-860f-ced4eec5c714.htm "FatekProgramHelper 方法")

[BuildReadBoolCommand 方法](../html/1f0e3930-99a4-b02f-d0b6-e139660f883a.htm "BuildReadBoolCommand 方法 ")

[BuildReadWordCommand 方法](../html/805ff280-0163-7ccb-25eb-2458424cc1ec.htm "BuildReadWordCommand 方法 ")

[BuildWriteBoolCommand 方法](../html/eecc01ad-74fc-3610-4939-9bbf87fc483d.htm "BuildWriteBoolCommand 方法 ")

[BuildWriteByteCommand 方法](../html/7b615285-6728-1344-384e-bf7cf918abf7.htm "BuildWriteByteCommand 方法 ")

[CalculateAcc 方法](../html/49b24ae3-78e3-e4d6-4743-f531b867da87.htm "CalculateAcc 方法 ")

[CheckReceiveDataComplete 方法](../html/b200e06c-38b7-9656-a340-54aa2373bb0a.htm "CheckReceiveDataComplete 方法 ")

[CheckResponse 方法](../html/1b4e1f3e-399f-4381-f912-86708faee605.htm "CheckResponse 方法 ")

[ExtraResponse 方法](../html/ad08cd70-4590-fc5d-a80d-a445450f9631.htm "ExtraResponse 方法 ")

[GetErrorDescriptionFromCode 方法](../html/fca0a278-c68e-cfe1-fcf4-816efcc878cb.htm "GetErrorDescriptionFromCode 方法 ")

[PackFatekCommand 方法](../html/fc22707e-222f-8fec-8d57-83742d4f17c1.htm "PackFatekCommand 方法 ")

[Read 方法](../html/cc1b8d93-c06c-41b9-748f-aa824b5f223b.htm "Read 方法 ")

[ReadAsync 方法](../html/c1a5498b-193d-f3bb-3e3c-7ab9585164ed.htm "ReadAsync 方法 ")

[ReadBool 方法](../html/5f136473-54c4-ce95-34f0-39c1028372f2.htm "ReadBool 方法 ")

[ReadBoolAsync 方法](../html/997a3300-c2ad-8816-bdfc-a2e0f9532bb4.htm "ReadBoolAsync 方法 ")

[ReadStatus 方法](../html/a73fd644-bf59-ca9a-d0ff-45d177c75e64.htm "ReadStatus 方法 ")

[ReadStatusAsync 方法](../html/62d31aab-0b53-ce2f-ffa5-09b3e4f85675.htm "ReadStatusAsync 方法 ")

[Run 方法](../html/d19f36a9-9645-fc5f-9bf8-7de6a709ed47.htm "Run 方法 ")

[RunAsync 方法](../html/0c1b1b0e-cecd-d2a2-5fb1-46a8a0321c39.htm "RunAsync 方法 ")

[Stop 方法](../html/2fc1d7f4-d7f2-ee5a-b271-69c33c1c3479.htm "Stop 方法 ")

[StopAsync 方法](../html/cfee21d4-ecbd-ea41-89f0-03d4f71102c7.htm "StopAsync 方法 ")

[Write 方法](../html/5adce02f-a40e-f29a-98db-23fb5791b508.htm "Write 方法 ")

[WriteAsync 方法](../html/4b5310b6-d2c7-d1ac-d43d-0c29bbbfea3b.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FatekProgramHelperBuildWriteByteCommand 方法 |

创建一条别入byte数据的指令信息，需要指定一些参数，按照字单位

**命名空间：**
 [HslCommunication.Profinet.FATEK.Helper](4ac2b39a-b5b1-9027-86a7-eadd294bcffa.htm)  
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
    站号

address
:   类型：SystemString  
    地址

value
:   类型：SystemByte  
    数组值

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
是否创建成功

![](../icons/SectionExpanded.png)参见

#### 引用

[FatekProgramHelper 类](5aa26a25-055f-eaea-5db8-af91b52b1e7b.htm)

[HslCommunication.Profinet.FATEK.Helper 命名空间](4ac2b39a-b5b1-9027-86a7-eadd294bcffa.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## CalculateAcc 方法 

[原文連結](http://api.hslcommunication.cn/html/49b24ae3-78e3-e4d6-4743-f531b867da87.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.FATEK.Helper](../html/4ac2b39a-b5b1-9027-86a7-eadd294bcffa.htm "HslCommunication.Profinet.FATEK.Helper")

[FatekProgramHelper 类](../html/5aa26a25-055f-eaea-5db8-af91b52b1e7b.htm "FatekProgramHelper 类")

[FatekProgramHelper 方法](../html/481082d5-89ec-7b3f-860f-ced4eec5c714.htm "FatekProgramHelper 方法")

[BuildReadBoolCommand 方法](../html/1f0e3930-99a4-b02f-d0b6-e139660f883a.htm "BuildReadBoolCommand 方法 ")

[BuildReadWordCommand 方法](../html/805ff280-0163-7ccb-25eb-2458424cc1ec.htm "BuildReadWordCommand 方法 ")

[BuildWriteBoolCommand 方法](../html/eecc01ad-74fc-3610-4939-9bbf87fc483d.htm "BuildWriteBoolCommand 方法 ")

[BuildWriteByteCommand 方法](../html/7b615285-6728-1344-384e-bf7cf918abf7.htm "BuildWriteByteCommand 方法 ")

[CalculateAcc 方法](../html/49b24ae3-78e3-e4d6-4743-f531b867da87.htm "CalculateAcc 方法 ")

[CheckReceiveDataComplete 方法](../html/b200e06c-38b7-9656-a340-54aa2373bb0a.htm "CheckReceiveDataComplete 方法 ")

[CheckResponse 方法](../html/1b4e1f3e-399f-4381-f912-86708faee605.htm "CheckResponse 方法 ")

[ExtraResponse 方法](../html/ad08cd70-4590-fc5d-a80d-a445450f9631.htm "ExtraResponse 方法 ")

[GetErrorDescriptionFromCode 方法](../html/fca0a278-c68e-cfe1-fcf4-816efcc878cb.htm "GetErrorDescriptionFromCode 方法 ")

[PackFatekCommand 方法](../html/fc22707e-222f-8fec-8d57-83742d4f17c1.htm "PackFatekCommand 方法 ")

[Read 方法](../html/cc1b8d93-c06c-41b9-748f-aa824b5f223b.htm "Read 方法 ")

[ReadAsync 方法](../html/c1a5498b-193d-f3bb-3e3c-7ab9585164ed.htm "ReadAsync 方法 ")

[ReadBool 方法](../html/5f136473-54c4-ce95-34f0-39c1028372f2.htm "ReadBool 方法 ")

[ReadBoolAsync 方法](../html/997a3300-c2ad-8816-bdfc-a2e0f9532bb4.htm "ReadBoolAsync 方法 ")

[ReadStatus 方法](../html/a73fd644-bf59-ca9a-d0ff-45d177c75e64.htm "ReadStatus 方法 ")

[ReadStatusAsync 方法](../html/62d31aab-0b53-ce2f-ffa5-09b3e4f85675.htm "ReadStatusAsync 方法 ")

[Run 方法](../html/d19f36a9-9645-fc5f-9bf8-7de6a709ed47.htm "Run 方法 ")

[RunAsync 方法](../html/0c1b1b0e-cecd-d2a2-5fb1-46a8a0321c39.htm "RunAsync 方法 ")

[Stop 方法](../html/2fc1d7f4-d7f2-ee5a-b271-69c33c1c3479.htm "Stop 方法 ")

[StopAsync 方法](../html/cfee21d4-ecbd-ea41-89f0-03d4f71102c7.htm "StopAsync 方法 ")

[Write 方法](../html/5adce02f-a40e-f29a-98db-23fb5791b508.htm "Write 方法 ")

[WriteAsync 方法](../html/4b5310b6-d2c7-d1ac-d43d-0c29bbbfea3b.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FatekProgramHelperCalculateAcc 方法 |

计算指令的和校验码

**命名空间：**
 [HslCommunication.Profinet.FATEK.Helper](4ac2b39a-b5b1-9027-86a7-eadd294bcffa.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static string CalculateAcc(
	string data
)
```

```
Public Shared Function CalculateAcc ( 
	data As String
) As String
```

```
public:
static String^ CalculateAcc(
	String^ data
)
```

```
static member CalculateAcc : 
        data : string -> string 
```

#### 参数

data
:   类型：SystemString  
    指令

#### 返回值

类型：String  
校验之后的信息

![](../icons/SectionExpanded.png)参见

#### 引用

[FatekProgramHelper 类](5aa26a25-055f-eaea-5db8-af91b52b1e7b.htm)

[HslCommunication.Profinet.FATEK.Helper 命名空间](4ac2b39a-b5b1-9027-86a7-eadd294bcffa.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## CheckReceiveDataComplete 方法 

[原文連結](http://api.hslcommunication.cn/html/b200e06c-38b7-9656-a340-54aa2373bb0a.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.FATEK.Helper](../html/4ac2b39a-b5b1-9027-86a7-eadd294bcffa.htm "HslCommunication.Profinet.FATEK.Helper")

[FatekProgramHelper 类](../html/5aa26a25-055f-eaea-5db8-af91b52b1e7b.htm "FatekProgramHelper 类")

[FatekProgramHelper 方法](../html/481082d5-89ec-7b3f-860f-ced4eec5c714.htm "FatekProgramHelper 方法")

[BuildReadBoolCommand 方法](../html/1f0e3930-99a4-b02f-d0b6-e139660f883a.htm "BuildReadBoolCommand 方法 ")

[BuildReadWordCommand 方法](../html/805ff280-0163-7ccb-25eb-2458424cc1ec.htm "BuildReadWordCommand 方法 ")

[BuildWriteBoolCommand 方法](../html/eecc01ad-74fc-3610-4939-9bbf87fc483d.htm "BuildWriteBoolCommand 方法 ")

[BuildWriteByteCommand 方法](../html/7b615285-6728-1344-384e-bf7cf918abf7.htm "BuildWriteByteCommand 方法 ")

[CalculateAcc 方法](../html/49b24ae3-78e3-e4d6-4743-f531b867da87.htm "CalculateAcc 方法 ")

[CheckReceiveDataComplete 方法](../html/b200e06c-38b7-9656-a340-54aa2373bb0a.htm "CheckReceiveDataComplete 方法 ")

[CheckResponse 方法](../html/1b4e1f3e-399f-4381-f912-86708faee605.htm "CheckResponse 方法 ")

[ExtraResponse 方法](../html/ad08cd70-4590-fc5d-a80d-a445450f9631.htm "ExtraResponse 方法 ")

[GetErrorDescriptionFromCode 方法](../html/fca0a278-c68e-cfe1-fcf4-816efcc878cb.htm "GetErrorDescriptionFromCode 方法 ")

[PackFatekCommand 方法](../html/fc22707e-222f-8fec-8d57-83742d4f17c1.htm "PackFatekCommand 方法 ")

[Read 方法](../html/cc1b8d93-c06c-41b9-748f-aa824b5f223b.htm "Read 方法 ")

[ReadAsync 方法](../html/c1a5498b-193d-f3bb-3e3c-7ab9585164ed.htm "ReadAsync 方法 ")

[ReadBool 方法](../html/5f136473-54c4-ce95-34f0-39c1028372f2.htm "ReadBool 方法 ")

[ReadBoolAsync 方法](../html/997a3300-c2ad-8816-bdfc-a2e0f9532bb4.htm "ReadBoolAsync 方法 ")

[ReadStatus 方法](../html/a73fd644-bf59-ca9a-d0ff-45d177c75e64.htm "ReadStatus 方法 ")

[ReadStatusAsync 方法](../html/62d31aab-0b53-ce2f-ffa5-09b3e4f85675.htm "ReadStatusAsync 方法 ")

[Run 方法](../html/d19f36a9-9645-fc5f-9bf8-7de6a709ed47.htm "Run 方法 ")

[RunAsync 方法](../html/0c1b1b0e-cecd-d2a2-5fb1-46a8a0321c39.htm "RunAsync 方法 ")

[Stop 方法](../html/2fc1d7f4-d7f2-ee5a-b271-69c33c1c3479.htm "Stop 方法 ")

[StopAsync 方法](../html/cfee21d4-ecbd-ea41-89f0-03d4f71102c7.htm "StopAsync 方法 ")

[Write 方法](../html/5adce02f-a40e-f29a-98db-23fb5791b508.htm "Write 方法 ")

[WriteAsync 方法](../html/4b5310b6-d2c7-d1ac-d43d-0c29bbbfea3b.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FatekProgramHelperCheckReceiveDataComplete 方法 |

检查当前的串口的数据接收是否完整

**命名空间：**
 [HslCommunication.Profinet.FATEK.Helper](4ac2b39a-b5b1-9027-86a7-eadd294bcffa.htm)  
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
    数据流

#### 返回值

类型：Boolean  
是否数据接收完成

![](../icons/SectionExpanded.png)参见

#### 引用

[FatekProgramHelper 类](5aa26a25-055f-eaea-5db8-af91b52b1e7b.htm)

[HslCommunication.Profinet.FATEK.Helper 命名空间](4ac2b39a-b5b1-9027-86a7-eadd294bcffa.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## CheckResponse 方法 

[原文連結](http://api.hslcommunication.cn/html/1b4e1f3e-399f-4381-f912-86708faee605.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.FATEK.Helper](../html/4ac2b39a-b5b1-9027-86a7-eadd294bcffa.htm "HslCommunication.Profinet.FATEK.Helper")

[FatekProgramHelper 类](../html/5aa26a25-055f-eaea-5db8-af91b52b1e7b.htm "FatekProgramHelper 类")

[FatekProgramHelper 方法](../html/481082d5-89ec-7b3f-860f-ced4eec5c714.htm "FatekProgramHelper 方法")

[BuildReadBoolCommand 方法](../html/1f0e3930-99a4-b02f-d0b6-e139660f883a.htm "BuildReadBoolCommand 方法 ")

[BuildReadWordCommand 方法](../html/805ff280-0163-7ccb-25eb-2458424cc1ec.htm "BuildReadWordCommand 方法 ")

[BuildWriteBoolCommand 方法](../html/eecc01ad-74fc-3610-4939-9bbf87fc483d.htm "BuildWriteBoolCommand 方法 ")

[BuildWriteByteCommand 方法](../html/7b615285-6728-1344-384e-bf7cf918abf7.htm "BuildWriteByteCommand 方法 ")

[CalculateAcc 方法](../html/49b24ae3-78e3-e4d6-4743-f531b867da87.htm "CalculateAcc 方法 ")

[CheckReceiveDataComplete 方法](../html/b200e06c-38b7-9656-a340-54aa2373bb0a.htm "CheckReceiveDataComplete 方法 ")

[CheckResponse 方法](../html/1b4e1f3e-399f-4381-f912-86708faee605.htm "CheckResponse 方法 ")

[ExtraResponse 方法](../html/ad08cd70-4590-fc5d-a80d-a445450f9631.htm "ExtraResponse 方法 ")

[GetErrorDescriptionFromCode 方法](../html/fca0a278-c68e-cfe1-fcf4-816efcc878cb.htm "GetErrorDescriptionFromCode 方法 ")

[PackFatekCommand 方法](../html/fc22707e-222f-8fec-8d57-83742d4f17c1.htm "PackFatekCommand 方法 ")

[Read 方法](../html/cc1b8d93-c06c-41b9-748f-aa824b5f223b.htm "Read 方法 ")

[ReadAsync 方法](../html/c1a5498b-193d-f3bb-3e3c-7ab9585164ed.htm "ReadAsync 方法 ")

[ReadBool 方法](../html/5f136473-54c4-ce95-34f0-39c1028372f2.htm "ReadBool 方法 ")

[ReadBoolAsync 方法](../html/997a3300-c2ad-8816-bdfc-a2e0f9532bb4.htm "ReadBoolAsync 方法 ")

[ReadStatus 方法](../html/a73fd644-bf59-ca9a-d0ff-45d177c75e64.htm "ReadStatus 方法 ")

[ReadStatusAsync 方法](../html/62d31aab-0b53-ce2f-ffa5-09b3e4f85675.htm "ReadStatusAsync 方法 ")

[Run 方法](../html/d19f36a9-9645-fc5f-9bf8-7de6a709ed47.htm "Run 方法 ")

[RunAsync 方法](../html/0c1b1b0e-cecd-d2a2-5fb1-46a8a0321c39.htm "RunAsync 方法 ")

[Stop 方法](../html/2fc1d7f4-d7f2-ee5a-b271-69c33c1c3479.htm "Stop 方法 ")

[StopAsync 方法](../html/cfee21d4-ecbd-ea41-89f0-03d4f71102c7.htm "StopAsync 方法 ")

[Write 方法](../html/5adce02f-a40e-f29a-98db-23fb5791b508.htm "Write 方法 ")

[WriteAsync 方法](../html/4b5310b6-d2c7-d1ac-d43d-0c29bbbfea3b.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FatekProgramHelperCheckResponse 方法 |

检查PLC反馈的报文是否正确，如果不正确，返回错误消息

**命名空间：**
 [HslCommunication.Profinet.FATEK.Helper](4ac2b39a-b5b1-9027-86a7-eadd294bcffa.htm)  
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
    PLC反馈的报文信息

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
反馈的报文是否正确

![](../icons/SectionExpanded.png)参见

#### 引用

[FatekProgramHelper 类](5aa26a25-055f-eaea-5db8-af91b52b1e7b.htm)

[HslCommunication.Profinet.FATEK.Helper 命名空间](4ac2b39a-b5b1-9027-86a7-eadd294bcffa.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ExtraResponse 方法 

[原文連結](http://api.hslcommunication.cn/html/ad08cd70-4590-fc5d-a80d-a445450f9631.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.FATEK.Helper](../html/4ac2b39a-b5b1-9027-86a7-eadd294bcffa.htm "HslCommunication.Profinet.FATEK.Helper")

[FatekProgramHelper 类](../html/5aa26a25-055f-eaea-5db8-af91b52b1e7b.htm "FatekProgramHelper 类")

[FatekProgramHelper 方法](../html/481082d5-89ec-7b3f-860f-ced4eec5c714.htm "FatekProgramHelper 方法")

[BuildReadBoolCommand 方法](../html/1f0e3930-99a4-b02f-d0b6-e139660f883a.htm "BuildReadBoolCommand 方法 ")

[BuildReadWordCommand 方法](../html/805ff280-0163-7ccb-25eb-2458424cc1ec.htm "BuildReadWordCommand 方法 ")

[BuildWriteBoolCommand 方法](../html/eecc01ad-74fc-3610-4939-9bbf87fc483d.htm "BuildWriteBoolCommand 方法 ")

[BuildWriteByteCommand 方法](../html/7b615285-6728-1344-384e-bf7cf918abf7.htm "BuildWriteByteCommand 方法 ")

[CalculateAcc 方法](../html/49b24ae3-78e3-e4d6-4743-f531b867da87.htm "CalculateAcc 方法 ")

[CheckReceiveDataComplete 方法](../html/b200e06c-38b7-9656-a340-54aa2373bb0a.htm "CheckReceiveDataComplete 方法 ")

[CheckResponse 方法](../html/1b4e1f3e-399f-4381-f912-86708faee605.htm "CheckResponse 方法 ")

[ExtraResponse 方法](../html/ad08cd70-4590-fc5d-a80d-a445450f9631.htm "ExtraResponse 方法 ")

[GetErrorDescriptionFromCode 方法](../html/fca0a278-c68e-cfe1-fcf4-816efcc878cb.htm "GetErrorDescriptionFromCode 方法 ")

[PackFatekCommand 方法](../html/fc22707e-222f-8fec-8d57-83742d4f17c1.htm "PackFatekCommand 方法 ")

[Read 方法](../html/cc1b8d93-c06c-41b9-748f-aa824b5f223b.htm "Read 方法 ")

[ReadAsync 方法](../html/c1a5498b-193d-f3bb-3e3c-7ab9585164ed.htm "ReadAsync 方法 ")

[ReadBool 方法](../html/5f136473-54c4-ce95-34f0-39c1028372f2.htm "ReadBool 方法 ")

[ReadBoolAsync 方法](../html/997a3300-c2ad-8816-bdfc-a2e0f9532bb4.htm "ReadBoolAsync 方法 ")

[ReadStatus 方法](../html/a73fd644-bf59-ca9a-d0ff-45d177c75e64.htm "ReadStatus 方法 ")

[ReadStatusAsync 方法](../html/62d31aab-0b53-ce2f-ffa5-09b3e4f85675.htm "ReadStatusAsync 方法 ")

[Run 方法](../html/d19f36a9-9645-fc5f-9bf8-7de6a709ed47.htm "Run 方法 ")

[RunAsync 方法](../html/0c1b1b0e-cecd-d2a2-5fb1-46a8a0321c39.htm "RunAsync 方法 ")

[Stop 方法](../html/2fc1d7f4-d7f2-ee5a-b271-69c33c1c3479.htm "Stop 方法 ")

[StopAsync 方法](../html/cfee21d4-ecbd-ea41-89f0-03d4f71102c7.htm "StopAsync 方法 ")

[Write 方法](../html/5adce02f-a40e-f29a-98db-23fb5791b508.htm "Write 方法 ")

[WriteAsync 方法](../html/4b5310b6-d2c7-d1ac-d43d-0c29bbbfea3b.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FatekProgramHelperExtraResponse 方法 |

提取当前的结果数据信息，针对的是字单位的方式

**命名空间：**
 [HslCommunication.Profinet.FATEK.Helper](4ac2b39a-b5b1-9027-86a7-eadd294bcffa.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static byte[] ExtraResponse(
	byte[] response,
	ushort length
)
```

```
Public Shared Function ExtraResponse ( 
	response As Byte(),
	length As UShort
) As Byte()
```

```
public:
static array<unsigned char>^ ExtraResponse(
	array<unsigned char>^ response, 
	unsigned short length
)
```

```
static member ExtraResponse : 
        response : byte[] * 
        length : uint16 -> byte[] 
```

#### 参数

response
:   类型：SystemByte  
    PLC返回的数据信息

length
:   类型：SystemUInt16  
    读取的长度内容

#### 返回值

类型：Byte  
结果数组

![](../icons/SectionExpanded.png)参见

#### 引用

[FatekProgramHelper 类](5aa26a25-055f-eaea-5db8-af91b52b1e7b.htm)

[HslCommunication.Profinet.FATEK.Helper 命名空间](4ac2b39a-b5b1-9027-86a7-eadd294bcffa.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GetErrorDescriptionFromCode 方法 

[原文連結](http://api.hslcommunication.cn/html/fca0a278-c68e-cfe1-fcf4-816efcc878cb.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.FATEK.Helper](../html/4ac2b39a-b5b1-9027-86a7-eadd294bcffa.htm "HslCommunication.Profinet.FATEK.Helper")

[FatekProgramHelper 类](../html/5aa26a25-055f-eaea-5db8-af91b52b1e7b.htm "FatekProgramHelper 类")

[FatekProgramHelper 方法](../html/481082d5-89ec-7b3f-860f-ced4eec5c714.htm "FatekProgramHelper 方法")

[BuildReadBoolCommand 方法](../html/1f0e3930-99a4-b02f-d0b6-e139660f883a.htm "BuildReadBoolCommand 方法 ")

[BuildReadWordCommand 方法](../html/805ff280-0163-7ccb-25eb-2458424cc1ec.htm "BuildReadWordCommand 方法 ")

[BuildWriteBoolCommand 方法](../html/eecc01ad-74fc-3610-4939-9bbf87fc483d.htm "BuildWriteBoolCommand 方法 ")

[BuildWriteByteCommand 方法](../html/7b615285-6728-1344-384e-bf7cf918abf7.htm "BuildWriteByteCommand 方法 ")

[CalculateAcc 方法](../html/49b24ae3-78e3-e4d6-4743-f531b867da87.htm "CalculateAcc 方法 ")

[CheckReceiveDataComplete 方法](../html/b200e06c-38b7-9656-a340-54aa2373bb0a.htm "CheckReceiveDataComplete 方法 ")

[CheckResponse 方法](../html/1b4e1f3e-399f-4381-f912-86708faee605.htm "CheckResponse 方法 ")

[ExtraResponse 方法](../html/ad08cd70-4590-fc5d-a80d-a445450f9631.htm "ExtraResponse 方法 ")

[GetErrorDescriptionFromCode 方法](../html/fca0a278-c68e-cfe1-fcf4-816efcc878cb.htm "GetErrorDescriptionFromCode 方法 ")

[PackFatekCommand 方法](../html/fc22707e-222f-8fec-8d57-83742d4f17c1.htm "PackFatekCommand 方法 ")

[Read 方法](../html/cc1b8d93-c06c-41b9-748f-aa824b5f223b.htm "Read 方法 ")

[ReadAsync 方法](../html/c1a5498b-193d-f3bb-3e3c-7ab9585164ed.htm "ReadAsync 方法 ")

[ReadBool 方法](../html/5f136473-54c4-ce95-34f0-39c1028372f2.htm "ReadBool 方法 ")

[ReadBoolAsync 方法](../html/997a3300-c2ad-8816-bdfc-a2e0f9532bb4.htm "ReadBoolAsync 方法 ")

[ReadStatus 方法](../html/a73fd644-bf59-ca9a-d0ff-45d177c75e64.htm "ReadStatus 方法 ")

[ReadStatusAsync 方法](../html/62d31aab-0b53-ce2f-ffa5-09b3e4f85675.htm "ReadStatusAsync 方法 ")

[Run 方法](../html/d19f36a9-9645-fc5f-9bf8-7de6a709ed47.htm "Run 方法 ")

[RunAsync 方法](../html/0c1b1b0e-cecd-d2a2-5fb1-46a8a0321c39.htm "RunAsync 方法 ")

[Stop 方法](../html/2fc1d7f4-d7f2-ee5a-b271-69c33c1c3479.htm "Stop 方法 ")

[StopAsync 方法](../html/cfee21d4-ecbd-ea41-89f0-03d4f71102c7.htm "StopAsync 方法 ")

[Write 方法](../html/5adce02f-a40e-f29a-98db-23fb5791b508.htm "Write 方法 ")

[WriteAsync 方法](../html/4b5310b6-d2c7-d1ac-d43d-0c29bbbfea3b.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FatekProgramHelperGetErrorDescriptionFromCode 方法 |

根据错误码获取到真实的文本信息

**命名空间：**
 [HslCommunication.Profinet.FATEK.Helper](4ac2b39a-b5b1-9027-86a7-eadd294bcffa.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static string GetErrorDescriptionFromCode(
	char code
)
```

```
Public Shared Function GetErrorDescriptionFromCode ( 
	code As Char
) As String
```

```
public:
static String^ GetErrorDescriptionFromCode(
	wchar_t code
)
```

```
static member GetErrorDescriptionFromCode : 
        code : char -> string 
```

#### 参数

code
:   类型：SystemChar  
    错误码

#### 返回值

类型：String  
错误的文本描述

![](../icons/SectionExpanded.png)参见

#### 引用

[FatekProgramHelper 类](5aa26a25-055f-eaea-5db8-af91b52b1e7b.htm)

[HslCommunication.Profinet.FATEK.Helper 命名空间](4ac2b39a-b5b1-9027-86a7-eadd294bcffa.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## PackFatekCommand 方法 

[原文連結](http://api.hslcommunication.cn/html/fc22707e-222f-8fec-8d57-83742d4f17c1.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.FATEK.Helper](../html/4ac2b39a-b5b1-9027-86a7-eadd294bcffa.htm "HslCommunication.Profinet.FATEK.Helper")

[FatekProgramHelper 类](../html/5aa26a25-055f-eaea-5db8-af91b52b1e7b.htm "FatekProgramHelper 类")

[FatekProgramHelper 方法](../html/481082d5-89ec-7b3f-860f-ced4eec5c714.htm "FatekProgramHelper 方法")

[BuildReadBoolCommand 方法](../html/1f0e3930-99a4-b02f-d0b6-e139660f883a.htm "BuildReadBoolCommand 方法 ")

[BuildReadWordCommand 方法](../html/805ff280-0163-7ccb-25eb-2458424cc1ec.htm "BuildReadWordCommand 方法 ")

[BuildWriteBoolCommand 方法](../html/eecc01ad-74fc-3610-4939-9bbf87fc483d.htm "BuildWriteBoolCommand 方法 ")

[BuildWriteByteCommand 方法](../html/7b615285-6728-1344-384e-bf7cf918abf7.htm "BuildWriteByteCommand 方法 ")

[CalculateAcc 方法](../html/49b24ae3-78e3-e4d6-4743-f531b867da87.htm "CalculateAcc 方法 ")

[CheckReceiveDataComplete 方法](../html/b200e06c-38b7-9656-a340-54aa2373bb0a.htm "CheckReceiveDataComplete 方法 ")

[CheckResponse 方法](../html/1b4e1f3e-399f-4381-f912-86708faee605.htm "CheckResponse 方法 ")

[ExtraResponse 方法](../html/ad08cd70-4590-fc5d-a80d-a445450f9631.htm "ExtraResponse 方法 ")

[GetErrorDescriptionFromCode 方法](../html/fca0a278-c68e-cfe1-fcf4-816efcc878cb.htm "GetErrorDescriptionFromCode 方法 ")

[PackFatekCommand 方法](../html/fc22707e-222f-8fec-8d57-83742d4f17c1.htm "PackFatekCommand 方法 ")

[Read 方法](../html/cc1b8d93-c06c-41b9-748f-aa824b5f223b.htm "Read 方法 ")

[ReadAsync 方法](../html/c1a5498b-193d-f3bb-3e3c-7ab9585164ed.htm "ReadAsync 方法 ")

[ReadBool 方法](../html/5f136473-54c4-ce95-34f0-39c1028372f2.htm "ReadBool 方法 ")

[ReadBoolAsync 方法](../html/997a3300-c2ad-8816-bdfc-a2e0f9532bb4.htm "ReadBoolAsync 方法 ")

[ReadStatus 方法](../html/a73fd644-bf59-ca9a-d0ff-45d177c75e64.htm "ReadStatus 方法 ")

[ReadStatusAsync 方法](../html/62d31aab-0b53-ce2f-ffa5-09b3e4f85675.htm "ReadStatusAsync 方法 ")

[Run 方法](../html/d19f36a9-9645-fc5f-9bf8-7de6a709ed47.htm "Run 方法 ")

[RunAsync 方法](../html/0c1b1b0e-cecd-d2a2-5fb1-46a8a0321c39.htm "RunAsync 方法 ")

[Stop 方法](../html/2fc1d7f4-d7f2-ee5a-b271-69c33c1c3479.htm "Stop 方法 ")

[StopAsync 方法](../html/cfee21d4-ecbd-ea41-89f0-03d4f71102c7.htm "StopAsync 方法 ")

[Write 方法](../html/5adce02f-a40e-f29a-98db-23fb5791b508.htm "Write 方法 ")

[WriteAsync 方法](../html/4b5310b6-d2c7-d1ac-d43d-0c29bbbfea3b.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FatekProgramHelperPackFatekCommand 方法 |

将Fatek的基本命令打包成可以发送PLC的电文消息

**命名空间：**
 [HslCommunication.Profinet.FATEK.Helper](4ac2b39a-b5b1-9027-86a7-eadd294bcffa.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static byte[] PackFatekCommand(
	byte station,
	string cmd
)
```

```
Public Shared Function PackFatekCommand ( 
	station As Byte,
	cmd As String
) As Byte()
```

```
public:
static array<unsigned char>^ PackFatekCommand(
	unsigned char station, 
	String^ cmd
)
```

```
static member PackFatekCommand : 
        station : byte * 
        cmd : string -> byte[] 
```

#### 参数

station
:   类型：SystemByte  
    PLC的站号信息

cmd
:   类型：SystemString  
    基本命令信息

#### 返回值

类型：Byte  
发送PLC的电文消息

![](../icons/SectionExpanded.png)参见

#### 引用

[FatekProgramHelper 类](5aa26a25-055f-eaea-5db8-af91b52b1e7b.htm)

[HslCommunication.Profinet.FATEK.Helper 命名空间](4ac2b39a-b5b1-9027-86a7-eadd294bcffa.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Read 方法 

[原文連結](http://api.hslcommunication.cn/html/cc1b8d93-c06c-41b9-748f-aa824b5f223b.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.FATEK.Helper](../html/4ac2b39a-b5b1-9027-86a7-eadd294bcffa.htm "HslCommunication.Profinet.FATEK.Helper")

[FatekProgramHelper 类](../html/5aa26a25-055f-eaea-5db8-af91b52b1e7b.htm "FatekProgramHelper 类")

[FatekProgramHelper 方法](../html/481082d5-89ec-7b3f-860f-ced4eec5c714.htm "FatekProgramHelper 方法")

[BuildReadBoolCommand 方法](../html/1f0e3930-99a4-b02f-d0b6-e139660f883a.htm "BuildReadBoolCommand 方法 ")

[BuildReadWordCommand 方法](../html/805ff280-0163-7ccb-25eb-2458424cc1ec.htm "BuildReadWordCommand 方法 ")

[BuildWriteBoolCommand 方法](../html/eecc01ad-74fc-3610-4939-9bbf87fc483d.htm "BuildWriteBoolCommand 方法 ")

[BuildWriteByteCommand 方法](../html/7b615285-6728-1344-384e-bf7cf918abf7.htm "BuildWriteByteCommand 方法 ")

[CalculateAcc 方法](../html/49b24ae3-78e3-e4d6-4743-f531b867da87.htm "CalculateAcc 方法 ")

[CheckReceiveDataComplete 方法](../html/b200e06c-38b7-9656-a340-54aa2373bb0a.htm "CheckReceiveDataComplete 方法 ")

[CheckResponse 方法](../html/1b4e1f3e-399f-4381-f912-86708faee605.htm "CheckResponse 方法 ")

[ExtraResponse 方法](../html/ad08cd70-4590-fc5d-a80d-a445450f9631.htm "ExtraResponse 方法 ")

[GetErrorDescriptionFromCode 方法](../html/fca0a278-c68e-cfe1-fcf4-816efcc878cb.htm "GetErrorDescriptionFromCode 方法 ")

[PackFatekCommand 方法](../html/fc22707e-222f-8fec-8d57-83742d4f17c1.htm "PackFatekCommand 方法 ")

[Read 方法](../html/cc1b8d93-c06c-41b9-748f-aa824b5f223b.htm "Read 方法 ")

[ReadAsync 方法](../html/c1a5498b-193d-f3bb-3e3c-7ab9585164ed.htm "ReadAsync 方法 ")

[ReadBool 方法](../html/5f136473-54c4-ce95-34f0-39c1028372f2.htm "ReadBool 方法 ")

[ReadBoolAsync 方法](../html/997a3300-c2ad-8816-bdfc-a2e0f9532bb4.htm "ReadBoolAsync 方法 ")

[ReadStatus 方法](../html/a73fd644-bf59-ca9a-d0ff-45d177c75e64.htm "ReadStatus 方法 ")

[ReadStatusAsync 方法](../html/62d31aab-0b53-ce2f-ffa5-09b3e4f85675.htm "ReadStatusAsync 方法 ")

[Run 方法](../html/d19f36a9-9645-fc5f-9bf8-7de6a709ed47.htm "Run 方法 ")

[RunAsync 方法](../html/0c1b1b0e-cecd-d2a2-5fb1-46a8a0321c39.htm "RunAsync 方法 ")

[Stop 方法](../html/2fc1d7f4-d7f2-ee5a-b271-69c33c1c3479.htm "Stop 方法 ")

[StopAsync 方法](../html/cfee21d4-ecbd-ea41-89f0-03d4f71102c7.htm "StopAsync 方法 ")

[Write 方法](../html/5adce02f-a40e-f29a-98db-23fb5791b508.htm "Write 方法 ")

[WriteAsync 方法](../html/4b5310b6-d2c7-d1ac-d43d-0c29bbbfea3b.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FatekProgramHelperRead 方法 |

批量读取PLC的字节数据，以字为单位，支持读取X,Y,M,S,D,T,C,R,RT,RC具体的地址范围需要根据PLC型号来确认，地址可以携带站号信息，例如 s=2;D100  
Read PLC byte data in batches, in word units. Supports reading X, Y, M, S, D, T, C, R, RT, RC.
The specific address range needs to be confirmed according to the PLC model, The address can carry station number information, such as s=2;D100

**命名空间：**
 [HslCommunication.Profinet.FATEK.Helper](4ac2b39a-b5b1-9027-86a7-eadd294bcffa.htm)  
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
	IReadWriteDevice device,
	byte station,
	string address,
	ushort length
)
```

```
Public Shared Function Read ( 
	device As IReadWriteDevice,
	station As Byte,
	address As String,
	length As UShort
) As OperateResult(Of Byte())
```

```
public:
static OperateResult<array<unsigned char>^>^ Read(
	IReadWriteDevice^ device, 
	unsigned char station, 
	String^ address, 
	unsigned short length
)
```

```
static member Read : 
        device : IReadWriteDevice * 
        station : byte * 
        address : string * 
        length : uint16 -> OperateResult<byte[]> 
```

#### 参数

device
:   类型：[HslCommunication.CoreIReadWriteDevice](d987afcc-f99c-6aed-3847-5b1890504635.htm)  
    PLC通信的对象

station
:   类型：SystemByte  
    设备的站点信息

address
:   类型：SystemString  
    地址信息

length
:   类型：SystemUInt16  
    数据长度

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
读取结果信息

![](../icons/SectionExpanded.png)参见

#### 引用

[FatekProgramHelper 类](5aa26a25-055f-eaea-5db8-af91b52b1e7b.htm)

[HslCommunication.Profinet.FATEK.Helper 命名空间](4ac2b39a-b5b1-9027-86a7-eadd294bcffa.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/c1a5498b-193d-f3bb-3e3c-7ab9585164ed.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.FATEK.Helper](../html/4ac2b39a-b5b1-9027-86a7-eadd294bcffa.htm "HslCommunication.Profinet.FATEK.Helper")

[FatekProgramHelper 类](../html/5aa26a25-055f-eaea-5db8-af91b52b1e7b.htm "FatekProgramHelper 类")

[FatekProgramHelper 方法](../html/481082d5-89ec-7b3f-860f-ced4eec5c714.htm "FatekProgramHelper 方法")

[BuildReadBoolCommand 方法](../html/1f0e3930-99a4-b02f-d0b6-e139660f883a.htm "BuildReadBoolCommand 方法 ")

[BuildReadWordCommand 方法](../html/805ff280-0163-7ccb-25eb-2458424cc1ec.htm "BuildReadWordCommand 方法 ")

[BuildWriteBoolCommand 方法](../html/eecc01ad-74fc-3610-4939-9bbf87fc483d.htm "BuildWriteBoolCommand 方法 ")

[BuildWriteByteCommand 方法](../html/7b615285-6728-1344-384e-bf7cf918abf7.htm "BuildWriteByteCommand 方法 ")

[CalculateAcc 方法](../html/49b24ae3-78e3-e4d6-4743-f531b867da87.htm "CalculateAcc 方法 ")

[CheckReceiveDataComplete 方法](../html/b200e06c-38b7-9656-a340-54aa2373bb0a.htm "CheckReceiveDataComplete 方法 ")

[CheckResponse 方法](../html/1b4e1f3e-399f-4381-f912-86708faee605.htm "CheckResponse 方法 ")

[ExtraResponse 方法](../html/ad08cd70-4590-fc5d-a80d-a445450f9631.htm "ExtraResponse 方法 ")

[GetErrorDescriptionFromCode 方法](../html/fca0a278-c68e-cfe1-fcf4-816efcc878cb.htm "GetErrorDescriptionFromCode 方法 ")

[PackFatekCommand 方法](../html/fc22707e-222f-8fec-8d57-83742d4f17c1.htm "PackFatekCommand 方法 ")

[Read 方法](../html/cc1b8d93-c06c-41b9-748f-aa824b5f223b.htm "Read 方法 ")

[ReadAsync 方法](../html/c1a5498b-193d-f3bb-3e3c-7ab9585164ed.htm "ReadAsync 方法 ")

[ReadBool 方法](../html/5f136473-54c4-ce95-34f0-39c1028372f2.htm "ReadBool 方法 ")

[ReadBoolAsync 方法](../html/997a3300-c2ad-8816-bdfc-a2e0f9532bb4.htm "ReadBoolAsync 方法 ")

[ReadStatus 方法](../html/a73fd644-bf59-ca9a-d0ff-45d177c75e64.htm "ReadStatus 方法 ")

[ReadStatusAsync 方法](../html/62d31aab-0b53-ce2f-ffa5-09b3e4f85675.htm "ReadStatusAsync 方法 ")

[Run 方法](../html/d19f36a9-9645-fc5f-9bf8-7de6a709ed47.htm "Run 方法 ")

[RunAsync 方法](../html/0c1b1b0e-cecd-d2a2-5fb1-46a8a0321c39.htm "RunAsync 方法 ")

[Stop 方法](../html/2fc1d7f4-d7f2-ee5a-b271-69c33c1c3479.htm "Stop 方法 ")

[StopAsync 方法](../html/cfee21d4-ecbd-ea41-89f0-03d4f71102c7.htm "StopAsync 方法 ")

[Write 方法](../html/5adce02f-a40e-f29a-98db-23fb5791b508.htm "Write 方法 ")

[WriteAsync 方法](../html/4b5310b6-d2c7-d1ac-d43d-0c29bbbfea3b.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FatekProgramHelperReadAsync 方法 |

批量读取PLC的字节数据，以字为单位，支持读取X,Y,M,S,D,T,C,R,RT,RC具体的地址范围需要根据PLC型号来确认，地址可以携带站号信息，例如 s=2;D100  
Read PLC byte data in batches, in word units. Supports reading X, Y, M, S, D, T, C, R, RT, RC.
The specific address range needs to be confirmed according to the PLC model, The address can carry station number information, such as s=2;D100

**命名空间：**
 [HslCommunication.Profinet.FATEK.Helper](4ac2b39a-b5b1-9027-86a7-eadd294bcffa.htm)  
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
	IReadWriteDevice device,
	byte station,
	string address,
	ushort length
)
```

```
Public Shared Function ReadAsync ( 
	device As IReadWriteDevice,
	station As Byte,
	address As String,
	length As UShort
) As Task(Of OperateResult(Of Byte()))
```

```
public:
static Task<OperateResult<array<unsigned char>^>^>^ ReadAsync(
	IReadWriteDevice^ device, 
	unsigned char station, 
	String^ address, 
	unsigned short length
)
```

```
static member ReadAsync : 
        device : IReadWriteDevice * 
        station : byte * 
        address : string * 
        length : uint16 -> Task<OperateResult<byte[]>> 
```

#### 参数

device
:   类型：[HslCommunication.CoreIReadWriteDevice](d987afcc-f99c-6aed-3847-5b1890504635.htm)  
    PLC通信的对象

station
:   类型：SystemByte  
    设备的站点信息

address
:   类型：SystemString  
    地址信息

length
:   类型：SystemUInt16  
    数据长度

#### 返回值

类型：Task[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
读取结果信息

![](../icons/SectionExpanded.png)参见

#### 引用

[FatekProgramHelper 类](5aa26a25-055f-eaea-5db8-af91b52b1e7b.htm)

[HslCommunication.Profinet.FATEK.Helper 命名空间](4ac2b39a-b5b1-9027-86a7-eadd294bcffa.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadBool 方法 

[原文連結](http://api.hslcommunication.cn/html/5f136473-54c4-ce95-34f0-39c1028372f2.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.FATEK.Helper](../html/4ac2b39a-b5b1-9027-86a7-eadd294bcffa.htm "HslCommunication.Profinet.FATEK.Helper")

[FatekProgramHelper 类](../html/5aa26a25-055f-eaea-5db8-af91b52b1e7b.htm "FatekProgramHelper 类")

[FatekProgramHelper 方法](../html/481082d5-89ec-7b3f-860f-ced4eec5c714.htm "FatekProgramHelper 方法")

[BuildReadBoolCommand 方法](../html/1f0e3930-99a4-b02f-d0b6-e139660f883a.htm "BuildReadBoolCommand 方法 ")

[BuildReadWordCommand 方法](../html/805ff280-0163-7ccb-25eb-2458424cc1ec.htm "BuildReadWordCommand 方法 ")

[BuildWriteBoolCommand 方法](../html/eecc01ad-74fc-3610-4939-9bbf87fc483d.htm "BuildWriteBoolCommand 方法 ")

[BuildWriteByteCommand 方法](../html/7b615285-6728-1344-384e-bf7cf918abf7.htm "BuildWriteByteCommand 方法 ")

[CalculateAcc 方法](../html/49b24ae3-78e3-e4d6-4743-f531b867da87.htm "CalculateAcc 方法 ")

[CheckReceiveDataComplete 方法](../html/b200e06c-38b7-9656-a340-54aa2373bb0a.htm "CheckReceiveDataComplete 方法 ")

[CheckResponse 方法](../html/1b4e1f3e-399f-4381-f912-86708faee605.htm "CheckResponse 方法 ")

[ExtraResponse 方法](../html/ad08cd70-4590-fc5d-a80d-a445450f9631.htm "ExtraResponse 方法 ")

[GetErrorDescriptionFromCode 方法](../html/fca0a278-c68e-cfe1-fcf4-816efcc878cb.htm "GetErrorDescriptionFromCode 方法 ")

[PackFatekCommand 方法](../html/fc22707e-222f-8fec-8d57-83742d4f17c1.htm "PackFatekCommand 方法 ")

[Read 方法](../html/cc1b8d93-c06c-41b9-748f-aa824b5f223b.htm "Read 方法 ")

[ReadAsync 方法](../html/c1a5498b-193d-f3bb-3e3c-7ab9585164ed.htm "ReadAsync 方法 ")

[ReadBool 方法](../html/5f136473-54c4-ce95-34f0-39c1028372f2.htm "ReadBool 方法 ")

[ReadBoolAsync 方法](../html/997a3300-c2ad-8816-bdfc-a2e0f9532bb4.htm "ReadBoolAsync 方法 ")

[ReadStatus 方法](../html/a73fd644-bf59-ca9a-d0ff-45d177c75e64.htm "ReadStatus 方法 ")

[ReadStatusAsync 方法](../html/62d31aab-0b53-ce2f-ffa5-09b3e4f85675.htm "ReadStatusAsync 方法 ")

[Run 方法](../html/d19f36a9-9645-fc5f-9bf8-7de6a709ed47.htm "Run 方法 ")

[RunAsync 方法](../html/0c1b1b0e-cecd-d2a2-5fb1-46a8a0321c39.htm "RunAsync 方法 ")

[Stop 方法](../html/2fc1d7f4-d7f2-ee5a-b271-69c33c1c3479.htm "Stop 方法 ")

[StopAsync 方法](../html/cfee21d4-ecbd-ea41-89f0-03d4f71102c7.htm "StopAsync 方法 ")

[Write 方法](../html/5adce02f-a40e-f29a-98db-23fb5791b508.htm "Write 方法 ")

[WriteAsync 方法](../html/4b5310b6-d2c7-d1ac-d43d-0c29bbbfea3b.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FatekProgramHelperReadBool 方法 |

批量读取bool类型数据，支持的类型为X,Y,M,S,T,C，具体的地址范围取决于PLC的类型，地址可以携带站号信息，例如 s=2;M100  
Read bool data in batches. The supported types are X, Y, M, S, T, C. The specific address range depends on the type of PLC,
The address can carry station number information, such as s=2;M100

**命名空间：**
 [HslCommunication.Profinet.FATEK.Helper](4ac2b39a-b5b1-9027-86a7-eadd294bcffa.htm)  
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
	IReadWriteDevice device,
	byte station,
	string address,
	ushort length
)
```

```
Public Shared Function ReadBool ( 
	device As IReadWriteDevice,
	station As Byte,
	address As String,
	length As UShort
) As OperateResult(Of Boolean())
```

```
public:
static OperateResult<array<bool>^>^ ReadBool(
	IReadWriteDevice^ device, 
	unsigned char station, 
	String^ address, 
	unsigned short length
)
```

```
static member ReadBool : 
        device : IReadWriteDevice * 
        station : byte * 
        address : string * 
        length : uint16 -> OperateResult<bool[]> 
```

#### 参数

device
:   类型：[HslCommunication.CoreIReadWriteDevice](d987afcc-f99c-6aed-3847-5b1890504635.htm)  
    PLC通信对象

station
:   类型：SystemByte  
    设备的站号信息

address
:   类型：SystemString  
    地址信息，比如X10，Y17，M100

length
:   类型：SystemUInt16  
    读取的长度

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Boolean  
读取结果信息

![](../icons/SectionExpanded.png)参见

#### 引用

[FatekProgramHelper 类](5aa26a25-055f-eaea-5db8-af91b52b1e7b.htm)

[HslCommunication.Profinet.FATEK.Helper 命名空间](4ac2b39a-b5b1-9027-86a7-eadd294bcffa.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadBoolAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/997a3300-c2ad-8816-bdfc-a2e0f9532bb4.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.FATEK.Helper](../html/4ac2b39a-b5b1-9027-86a7-eadd294bcffa.htm "HslCommunication.Profinet.FATEK.Helper")

[FatekProgramHelper 类](../html/5aa26a25-055f-eaea-5db8-af91b52b1e7b.htm "FatekProgramHelper 类")

[FatekProgramHelper 方法](../html/481082d5-89ec-7b3f-860f-ced4eec5c714.htm "FatekProgramHelper 方法")

[BuildReadBoolCommand 方法](../html/1f0e3930-99a4-b02f-d0b6-e139660f883a.htm "BuildReadBoolCommand 方法 ")

[BuildReadWordCommand 方法](../html/805ff280-0163-7ccb-25eb-2458424cc1ec.htm "BuildReadWordCommand 方法 ")

[BuildWriteBoolCommand 方法](../html/eecc01ad-74fc-3610-4939-9bbf87fc483d.htm "BuildWriteBoolCommand 方法 ")

[BuildWriteByteCommand 方法](../html/7b615285-6728-1344-384e-bf7cf918abf7.htm "BuildWriteByteCommand 方法 ")

[CalculateAcc 方法](../html/49b24ae3-78e3-e4d6-4743-f531b867da87.htm "CalculateAcc 方法 ")

[CheckReceiveDataComplete 方法](../html/b200e06c-38b7-9656-a340-54aa2373bb0a.htm "CheckReceiveDataComplete 方法 ")

[CheckResponse 方法](../html/1b4e1f3e-399f-4381-f912-86708faee605.htm "CheckResponse 方法 ")

[ExtraResponse 方法](../html/ad08cd70-4590-fc5d-a80d-a445450f9631.htm "ExtraResponse 方法 ")

[GetErrorDescriptionFromCode 方法](../html/fca0a278-c68e-cfe1-fcf4-816efcc878cb.htm "GetErrorDescriptionFromCode 方法 ")

[PackFatekCommand 方法](../html/fc22707e-222f-8fec-8d57-83742d4f17c1.htm "PackFatekCommand 方法 ")

[Read 方法](../html/cc1b8d93-c06c-41b9-748f-aa824b5f223b.htm "Read 方法 ")

[ReadAsync 方法](../html/c1a5498b-193d-f3bb-3e3c-7ab9585164ed.htm "ReadAsync 方法 ")

[ReadBool 方法](../html/5f136473-54c4-ce95-34f0-39c1028372f2.htm "ReadBool 方法 ")

[ReadBoolAsync 方法](../html/997a3300-c2ad-8816-bdfc-a2e0f9532bb4.htm "ReadBoolAsync 方法 ")

[ReadStatus 方法](../html/a73fd644-bf59-ca9a-d0ff-45d177c75e64.htm "ReadStatus 方法 ")

[ReadStatusAsync 方法](../html/62d31aab-0b53-ce2f-ffa5-09b3e4f85675.htm "ReadStatusAsync 方法 ")

[Run 方法](../html/d19f36a9-9645-fc5f-9bf8-7de6a709ed47.htm "Run 方法 ")

[RunAsync 方法](../html/0c1b1b0e-cecd-d2a2-5fb1-46a8a0321c39.htm "RunAsync 方法 ")

[Stop 方法](../html/2fc1d7f4-d7f2-ee5a-b271-69c33c1c3479.htm "Stop 方法 ")

[StopAsync 方法](../html/cfee21d4-ecbd-ea41-89f0-03d4f71102c7.htm "StopAsync 方法 ")

[Write 方法](../html/5adce02f-a40e-f29a-98db-23fb5791b508.htm "Write 方法 ")

[WriteAsync 方法](../html/4b5310b6-d2c7-d1ac-d43d-0c29bbbfea3b.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FatekProgramHelperReadBoolAsync 方法 |

批量读取bool类型数据，支持的类型为X,Y,M,S,T,C，具体的地址范围取决于PLC的类型，地址可以携带站号信息，例如 s=2;M100  
Read bool data in batches. The supported types are X, Y, M, S, T, C. The specific address range depends on the type of PLC,
The address can carry station number information, such as s=2;M100

**命名空间：**
 [HslCommunication.Profinet.FATEK.Helper](4ac2b39a-b5b1-9027-86a7-eadd294bcffa.htm)  
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
	IReadWriteDevice device,
	byte station,
	string address,
	ushort length
)
```

```
Public Shared Function ReadBoolAsync ( 
	device As IReadWriteDevice,
	station As Byte,
	address As String,
	length As UShort
) As Task(Of OperateResult(Of Boolean()))
```

```
public:
static Task<OperateResult<array<bool>^>^>^ ReadBoolAsync(
	IReadWriteDevice^ device, 
	unsigned char station, 
	String^ address, 
	unsigned short length
)
```

```
static member ReadBoolAsync : 
        device : IReadWriteDevice * 
        station : byte * 
        address : string * 
        length : uint16 -> Task<OperateResult<bool[]>> 
```

#### 参数

device
:   类型：[HslCommunication.CoreIReadWriteDevice](d987afcc-f99c-6aed-3847-5b1890504635.htm)  
    PLC通信对象

station
:   类型：SystemByte  
    设备的站号信息

address
:   类型：SystemString  
    地址信息，比如X10，Y17，M100

length
:   类型：SystemUInt16  
    读取的长度

#### 返回值

类型：Task[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Boolean  
读取结果信息

![](../icons/SectionExpanded.png)参见

#### 引用

[FatekProgramHelper 类](5aa26a25-055f-eaea-5db8-af91b52b1e7b.htm)

[HslCommunication.Profinet.FATEK.Helper 命名空间](4ac2b39a-b5b1-9027-86a7-eadd294bcffa.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadStatus 方法 

[原文連結](http://api.hslcommunication.cn/html/a73fd644-bf59-ca9a-d0ff-45d177c75e64.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.FATEK.Helper](../html/4ac2b39a-b5b1-9027-86a7-eadd294bcffa.htm "HslCommunication.Profinet.FATEK.Helper")

[FatekProgramHelper 类](../html/5aa26a25-055f-eaea-5db8-af91b52b1e7b.htm "FatekProgramHelper 类")

[FatekProgramHelper 方法](../html/481082d5-89ec-7b3f-860f-ced4eec5c714.htm "FatekProgramHelper 方法")

[BuildReadBoolCommand 方法](../html/1f0e3930-99a4-b02f-d0b6-e139660f883a.htm "BuildReadBoolCommand 方法 ")

[BuildReadWordCommand 方法](../html/805ff280-0163-7ccb-25eb-2458424cc1ec.htm "BuildReadWordCommand 方法 ")

[BuildWriteBoolCommand 方法](../html/eecc01ad-74fc-3610-4939-9bbf87fc483d.htm "BuildWriteBoolCommand 方法 ")

[BuildWriteByteCommand 方法](../html/7b615285-6728-1344-384e-bf7cf918abf7.htm "BuildWriteByteCommand 方法 ")

[CalculateAcc 方法](../html/49b24ae3-78e3-e4d6-4743-f531b867da87.htm "CalculateAcc 方法 ")

[CheckReceiveDataComplete 方法](../html/b200e06c-38b7-9656-a340-54aa2373bb0a.htm "CheckReceiveDataComplete 方法 ")

[CheckResponse 方法](../html/1b4e1f3e-399f-4381-f912-86708faee605.htm "CheckResponse 方法 ")

[ExtraResponse 方法](../html/ad08cd70-4590-fc5d-a80d-a445450f9631.htm "ExtraResponse 方法 ")

[GetErrorDescriptionFromCode 方法](../html/fca0a278-c68e-cfe1-fcf4-816efcc878cb.htm "GetErrorDescriptionFromCode 方法 ")

[PackFatekCommand 方法](../html/fc22707e-222f-8fec-8d57-83742d4f17c1.htm "PackFatekCommand 方法 ")

[Read 方法](../html/cc1b8d93-c06c-41b9-748f-aa824b5f223b.htm "Read 方法 ")

[ReadAsync 方法](../html/c1a5498b-193d-f3bb-3e3c-7ab9585164ed.htm "ReadAsync 方法 ")

[ReadBool 方法](../html/5f136473-54c4-ce95-34f0-39c1028372f2.htm "ReadBool 方法 ")

[ReadBoolAsync 方法](../html/997a3300-c2ad-8816-bdfc-a2e0f9532bb4.htm "ReadBoolAsync 方法 ")

[ReadStatus 方法](../html/a73fd644-bf59-ca9a-d0ff-45d177c75e64.htm "ReadStatus 方法 ")

[ReadStatusAsync 方法](../html/62d31aab-0b53-ce2f-ffa5-09b3e4f85675.htm "ReadStatusAsync 方法 ")

[Run 方法](../html/d19f36a9-9645-fc5f-9bf8-7de6a709ed47.htm "Run 方法 ")

[RunAsync 方法](../html/0c1b1b0e-cecd-d2a2-5fb1-46a8a0321c39.htm "RunAsync 方法 ")

[Stop 方法](../html/2fc1d7f4-d7f2-ee5a-b271-69c33c1c3479.htm "Stop 方法 ")

[StopAsync 方法](../html/cfee21d4-ecbd-ea41-89f0-03d4f71102c7.htm "StopAsync 方法 ")

[Write 方法](../html/5adce02f-a40e-f29a-98db-23fb5791b508.htm "Write 方法 ")

[WriteAsync 方法](../html/4b5310b6-d2c7-d1ac-d43d-0c29bbbfea3b.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FatekProgramHelperReadStatus 方法 |

读取当前PLC的状态信息，返回一个bool数组，同时包含了几种电量信息，分别为 0: RUN/STOP, 1: BAT LOW/正常, 2: Ladder checksum error/正常, 3: 使用ROM PACK/未使用,
4: WDT Timeout/正常, 5: 设定ID/未设ID， 6： 紧急停机/正常  
Read the status information of the current PLC and return a bool array, which also contains several power information, 0: RUN/STOP, 1: BAT LOW/normal,
2: Ladder checksum error/normal, 3: Use ROM PACK/ Not used, 4: WDT Timeout/Normal, 5: ID set/ID not set, 6: Emergency stop/Normal

**命名空间：**
 [HslCommunication.Profinet.FATEK.Helper](4ac2b39a-b5b1-9027-86a7-eadd294bcffa.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult<bool[]> ReadStatus(
	IReadWriteDevice device,
	byte station
)
```

```
Public Shared Function ReadStatus ( 
	device As IReadWriteDevice,
	station As Byte
) As OperateResult(Of Boolean())
```

```
public:
static OperateResult<array<bool>^>^ ReadStatus(
	IReadWriteDevice^ device, 
	unsigned char station
)
```

```
static member ReadStatus : 
        device : IReadWriteDevice * 
        station : byte -> OperateResult<bool[]> 
```

#### 参数

device
:   类型：[HslCommunication.CoreIReadWriteDevice](d987afcc-f99c-6aed-3847-5b1890504635.htm)  
    PLC通信对象

station
:   类型：SystemByte  
    站号信息

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Boolean  
状态结果信息

![](../icons/SectionExpanded.png)参见

#### 引用

[FatekProgramHelper 类](5aa26a25-055f-eaea-5db8-af91b52b1e7b.htm)

[HslCommunication.Profinet.FATEK.Helper 命名空间](4ac2b39a-b5b1-9027-86a7-eadd294bcffa.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadStatusAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/62d31aab-0b53-ce2f-ffa5-09b3e4f85675.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.FATEK.Helper](../html/4ac2b39a-b5b1-9027-86a7-eadd294bcffa.htm "HslCommunication.Profinet.FATEK.Helper")

[FatekProgramHelper 类](../html/5aa26a25-055f-eaea-5db8-af91b52b1e7b.htm "FatekProgramHelper 类")

[FatekProgramHelper 方法](../html/481082d5-89ec-7b3f-860f-ced4eec5c714.htm "FatekProgramHelper 方法")

[BuildReadBoolCommand 方法](../html/1f0e3930-99a4-b02f-d0b6-e139660f883a.htm "BuildReadBoolCommand 方法 ")

[BuildReadWordCommand 方法](../html/805ff280-0163-7ccb-25eb-2458424cc1ec.htm "BuildReadWordCommand 方法 ")

[BuildWriteBoolCommand 方法](../html/eecc01ad-74fc-3610-4939-9bbf87fc483d.htm "BuildWriteBoolCommand 方法 ")

[BuildWriteByteCommand 方法](../html/7b615285-6728-1344-384e-bf7cf918abf7.htm "BuildWriteByteCommand 方法 ")

[CalculateAcc 方法](../html/49b24ae3-78e3-e4d6-4743-f531b867da87.htm "CalculateAcc 方法 ")

[CheckReceiveDataComplete 方法](../html/b200e06c-38b7-9656-a340-54aa2373bb0a.htm "CheckReceiveDataComplete 方法 ")

[CheckResponse 方法](../html/1b4e1f3e-399f-4381-f912-86708faee605.htm "CheckResponse 方法 ")

[ExtraResponse 方法](../html/ad08cd70-4590-fc5d-a80d-a445450f9631.htm "ExtraResponse 方法 ")

[GetErrorDescriptionFromCode 方法](../html/fca0a278-c68e-cfe1-fcf4-816efcc878cb.htm "GetErrorDescriptionFromCode 方法 ")

[PackFatekCommand 方法](../html/fc22707e-222f-8fec-8d57-83742d4f17c1.htm "PackFatekCommand 方法 ")

[Read 方法](../html/cc1b8d93-c06c-41b9-748f-aa824b5f223b.htm "Read 方法 ")

[ReadAsync 方法](../html/c1a5498b-193d-f3bb-3e3c-7ab9585164ed.htm "ReadAsync 方法 ")

[ReadBool 方法](../html/5f136473-54c4-ce95-34f0-39c1028372f2.htm "ReadBool 方法 ")

[ReadBoolAsync 方法](../html/997a3300-c2ad-8816-bdfc-a2e0f9532bb4.htm "ReadBoolAsync 方法 ")

[ReadStatus 方法](../html/a73fd644-bf59-ca9a-d0ff-45d177c75e64.htm "ReadStatus 方法 ")

[ReadStatusAsync 方法](../html/62d31aab-0b53-ce2f-ffa5-09b3e4f85675.htm "ReadStatusAsync 方法 ")

[Run 方法](../html/d19f36a9-9645-fc5f-9bf8-7de6a709ed47.htm "Run 方法 ")

[RunAsync 方法](../html/0c1b1b0e-cecd-d2a2-5fb1-46a8a0321c39.htm "RunAsync 方法 ")

[Stop 方法](../html/2fc1d7f4-d7f2-ee5a-b271-69c33c1c3479.htm "Stop 方法 ")

[StopAsync 方法](../html/cfee21d4-ecbd-ea41-89f0-03d4f71102c7.htm "StopAsync 方法 ")

[Write 方法](../html/5adce02f-a40e-f29a-98db-23fb5791b508.htm "Write 方法 ")

[WriteAsync 方法](../html/4b5310b6-d2c7-d1ac-d43d-0c29bbbfea3b.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FatekProgramHelperReadStatusAsync 方法 |

使PLC处于STOP状态

**命名空间：**
 [HslCommunication.Profinet.FATEK.Helper](4ac2b39a-b5b1-9027-86a7-eadd294bcffa.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static Task<OperateResult<bool[]>> ReadStatusAsync(
	IReadWriteDevice device,
	byte station
)
```

```
Public Shared Function ReadStatusAsync ( 
	device As IReadWriteDevice,
	station As Byte
) As Task(Of OperateResult(Of Boolean()))
```

```
public:
static Task<OperateResult<array<bool>^>^>^ ReadStatusAsync(
	IReadWriteDevice^ device, 
	unsigned char station
)
```

```
static member ReadStatusAsync : 
        device : IReadWriteDevice * 
        station : byte -> Task<OperateResult<bool[]>> 
```

#### 参数

device
:   类型：[HslCommunication.CoreIReadWriteDevice](d987afcc-f99c-6aed-3847-5b1890504635.htm)  
    PLC通信对象

station
:   类型：SystemByte  
    站号信息

#### 返回值

类型：Task[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Boolean  
是否操作成功

![](../icons/SectionExpanded.png)参见

#### 引用

[FatekProgramHelper 类](5aa26a25-055f-eaea-5db8-af91b52b1e7b.htm)

[HslCommunication.Profinet.FATEK.Helper 命名空间](4ac2b39a-b5b1-9027-86a7-eadd294bcffa.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Run 方法 

[原文連結](http://api.hslcommunication.cn/html/d19f36a9-9645-fc5f-9bf8-7de6a709ed47.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.FATEK.Helper](../html/4ac2b39a-b5b1-9027-86a7-eadd294bcffa.htm "HslCommunication.Profinet.FATEK.Helper")

[FatekProgramHelper 类](../html/5aa26a25-055f-eaea-5db8-af91b52b1e7b.htm "FatekProgramHelper 类")

[FatekProgramHelper 方法](../html/481082d5-89ec-7b3f-860f-ced4eec5c714.htm "FatekProgramHelper 方法")

[BuildReadBoolCommand 方法](../html/1f0e3930-99a4-b02f-d0b6-e139660f883a.htm "BuildReadBoolCommand 方法 ")

[BuildReadWordCommand 方法](../html/805ff280-0163-7ccb-25eb-2458424cc1ec.htm "BuildReadWordCommand 方法 ")

[BuildWriteBoolCommand 方法](../html/eecc01ad-74fc-3610-4939-9bbf87fc483d.htm "BuildWriteBoolCommand 方法 ")

[BuildWriteByteCommand 方法](../html/7b615285-6728-1344-384e-bf7cf918abf7.htm "BuildWriteByteCommand 方法 ")

[CalculateAcc 方法](../html/49b24ae3-78e3-e4d6-4743-f531b867da87.htm "CalculateAcc 方法 ")

[CheckReceiveDataComplete 方法](../html/b200e06c-38b7-9656-a340-54aa2373bb0a.htm "CheckReceiveDataComplete 方法 ")

[CheckResponse 方法](../html/1b4e1f3e-399f-4381-f912-86708faee605.htm "CheckResponse 方法 ")

[ExtraResponse 方法](../html/ad08cd70-4590-fc5d-a80d-a445450f9631.htm "ExtraResponse 方法 ")

[GetErrorDescriptionFromCode 方法](../html/fca0a278-c68e-cfe1-fcf4-816efcc878cb.htm "GetErrorDescriptionFromCode 方法 ")

[PackFatekCommand 方法](../html/fc22707e-222f-8fec-8d57-83742d4f17c1.htm "PackFatekCommand 方法 ")

[Read 方法](../html/cc1b8d93-c06c-41b9-748f-aa824b5f223b.htm "Read 方法 ")

[ReadAsync 方法](../html/c1a5498b-193d-f3bb-3e3c-7ab9585164ed.htm "ReadAsync 方法 ")

[ReadBool 方法](../html/5f136473-54c4-ce95-34f0-39c1028372f2.htm "ReadBool 方法 ")

[ReadBoolAsync 方法](../html/997a3300-c2ad-8816-bdfc-a2e0f9532bb4.htm "ReadBoolAsync 方法 ")

[ReadStatus 方法](../html/a73fd644-bf59-ca9a-d0ff-45d177c75e64.htm "ReadStatus 方法 ")

[ReadStatusAsync 方法](../html/62d31aab-0b53-ce2f-ffa5-09b3e4f85675.htm "ReadStatusAsync 方法 ")

[Run 方法](../html/d19f36a9-9645-fc5f-9bf8-7de6a709ed47.htm "Run 方法 ")

[RunAsync 方法](../html/0c1b1b0e-cecd-d2a2-5fb1-46a8a0321c39.htm "RunAsync 方法 ")

[Stop 方法](../html/2fc1d7f4-d7f2-ee5a-b271-69c33c1c3479.htm "Stop 方法 ")

[StopAsync 方法](../html/cfee21d4-ecbd-ea41-89f0-03d4f71102c7.htm "StopAsync 方法 ")

[Write 方法](../html/5adce02f-a40e-f29a-98db-23fb5791b508.htm "Write 方法 ")

[WriteAsync 方法](../html/4b5310b6-d2c7-d1ac-d43d-0c29bbbfea3b.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FatekProgramHelperRun 方法 |

使PLC处于RUN的状态

**命名空间：**
 [HslCommunication.Profinet.FATEK.Helper](4ac2b39a-b5b1-9027-86a7-eadd294bcffa.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult Run(
	IReadWriteDevice device,
	byte station
)
```

```
Public Shared Function Run ( 
	device As IReadWriteDevice,
	station As Byte
) As OperateResult
```

```
public:
static OperateResult^ Run(
	IReadWriteDevice^ device, 
	unsigned char station
)
```

```
static member Run : 
        device : IReadWriteDevice * 
        station : byte -> OperateResult 
```

#### 参数

device
:   类型：[HslCommunication.CoreIReadWriteDevice](d987afcc-f99c-6aed-3847-5b1890504635.htm)  
    PLC通信对象

station
:   类型：SystemByte  
    站号信息

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否操作成功

![](../icons/SectionExpanded.png)参见

#### 引用

[FatekProgramHelper 类](5aa26a25-055f-eaea-5db8-af91b52b1e7b.htm)

[HslCommunication.Profinet.FATEK.Helper 命名空间](4ac2b39a-b5b1-9027-86a7-eadd294bcffa.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## RunAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/0c1b1b0e-cecd-d2a2-5fb1-46a8a0321c39.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.FATEK.Helper](../html/4ac2b39a-b5b1-9027-86a7-eadd294bcffa.htm "HslCommunication.Profinet.FATEK.Helper")

[FatekProgramHelper 类](../html/5aa26a25-055f-eaea-5db8-af91b52b1e7b.htm "FatekProgramHelper 类")

[FatekProgramHelper 方法](../html/481082d5-89ec-7b3f-860f-ced4eec5c714.htm "FatekProgramHelper 方法")

[BuildReadBoolCommand 方法](../html/1f0e3930-99a4-b02f-d0b6-e139660f883a.htm "BuildReadBoolCommand 方法 ")

[BuildReadWordCommand 方法](../html/805ff280-0163-7ccb-25eb-2458424cc1ec.htm "BuildReadWordCommand 方法 ")

[BuildWriteBoolCommand 方法](../html/eecc01ad-74fc-3610-4939-9bbf87fc483d.htm "BuildWriteBoolCommand 方法 ")

[BuildWriteByteCommand 方法](../html/7b615285-6728-1344-384e-bf7cf918abf7.htm "BuildWriteByteCommand 方法 ")

[CalculateAcc 方法](../html/49b24ae3-78e3-e4d6-4743-f531b867da87.htm "CalculateAcc 方法 ")

[CheckReceiveDataComplete 方法](../html/b200e06c-38b7-9656-a340-54aa2373bb0a.htm "CheckReceiveDataComplete 方法 ")

[CheckResponse 方法](../html/1b4e1f3e-399f-4381-f912-86708faee605.htm "CheckResponse 方法 ")

[ExtraResponse 方法](../html/ad08cd70-4590-fc5d-a80d-a445450f9631.htm "ExtraResponse 方法 ")

[GetErrorDescriptionFromCode 方法](../html/fca0a278-c68e-cfe1-fcf4-816efcc878cb.htm "GetErrorDescriptionFromCode 方法 ")

[PackFatekCommand 方法](../html/fc22707e-222f-8fec-8d57-83742d4f17c1.htm "PackFatekCommand 方法 ")

[Read 方法](../html/cc1b8d93-c06c-41b9-748f-aa824b5f223b.htm "Read 方法 ")

[ReadAsync 方法](../html/c1a5498b-193d-f3bb-3e3c-7ab9585164ed.htm "ReadAsync 方法 ")

[ReadBool 方法](../html/5f136473-54c4-ce95-34f0-39c1028372f2.htm "ReadBool 方法 ")

[ReadBoolAsync 方法](../html/997a3300-c2ad-8816-bdfc-a2e0f9532bb4.htm "ReadBoolAsync 方法 ")

[ReadStatus 方法](../html/a73fd644-bf59-ca9a-d0ff-45d177c75e64.htm "ReadStatus 方法 ")

[ReadStatusAsync 方法](../html/62d31aab-0b53-ce2f-ffa5-09b3e4f85675.htm "ReadStatusAsync 方法 ")

[Run 方法](../html/d19f36a9-9645-fc5f-9bf8-7de6a709ed47.htm "Run 方法 ")

[RunAsync 方法](../html/0c1b1b0e-cecd-d2a2-5fb1-46a8a0321c39.htm "RunAsync 方法 ")

[Stop 方法](../html/2fc1d7f4-d7f2-ee5a-b271-69c33c1c3479.htm "Stop 方法 ")

[StopAsync 方法](../html/cfee21d4-ecbd-ea41-89f0-03d4f71102c7.htm "StopAsync 方法 ")

[Write 方法](../html/5adce02f-a40e-f29a-98db-23fb5791b508.htm "Write 方法 ")

[WriteAsync 方法](../html/4b5310b6-d2c7-d1ac-d43d-0c29bbbfea3b.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FatekProgramHelperRunAsync 方法 |

使PLC处于RUN的状态

**命名空间：**
 [HslCommunication.Profinet.FATEK.Helper](4ac2b39a-b5b1-9027-86a7-eadd294bcffa.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static Task<OperateResult> RunAsync(
	IReadWriteDevice device,
	byte station
)
```

```
Public Shared Function RunAsync ( 
	device As IReadWriteDevice,
	station As Byte
) As Task(Of OperateResult)
```

```
public:
static Task<OperateResult^>^ RunAsync(
	IReadWriteDevice^ device, 
	unsigned char station
)
```

```
static member RunAsync : 
        device : IReadWriteDevice * 
        station : byte -> Task<OperateResult> 
```

#### 参数

device
:   类型：[HslCommunication.CoreIReadWriteDevice](d987afcc-f99c-6aed-3847-5b1890504635.htm)  
    PLC通信对象

station
:   类型：SystemByte  
    站号信息

#### 返回值

类型：Task[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否操作成功

![](../icons/SectionExpanded.png)参见

#### 引用

[FatekProgramHelper 类](5aa26a25-055f-eaea-5db8-af91b52b1e7b.htm)

[HslCommunication.Profinet.FATEK.Helper 命名空间](4ac2b39a-b5b1-9027-86a7-eadd294bcffa.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Stop 方法 

[原文連結](http://api.hslcommunication.cn/html/2fc1d7f4-d7f2-ee5a-b271-69c33c1c3479.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.FATEK.Helper](../html/4ac2b39a-b5b1-9027-86a7-eadd294bcffa.htm "HslCommunication.Profinet.FATEK.Helper")

[FatekProgramHelper 类](../html/5aa26a25-055f-eaea-5db8-af91b52b1e7b.htm "FatekProgramHelper 类")

[FatekProgramHelper 方法](../html/481082d5-89ec-7b3f-860f-ced4eec5c714.htm "FatekProgramHelper 方法")

[BuildReadBoolCommand 方法](../html/1f0e3930-99a4-b02f-d0b6-e139660f883a.htm "BuildReadBoolCommand 方法 ")

[BuildReadWordCommand 方法](../html/805ff280-0163-7ccb-25eb-2458424cc1ec.htm "BuildReadWordCommand 方法 ")

[BuildWriteBoolCommand 方法](../html/eecc01ad-74fc-3610-4939-9bbf87fc483d.htm "BuildWriteBoolCommand 方法 ")

[BuildWriteByteCommand 方法](../html/7b615285-6728-1344-384e-bf7cf918abf7.htm "BuildWriteByteCommand 方法 ")

[CalculateAcc 方法](../html/49b24ae3-78e3-e4d6-4743-f531b867da87.htm "CalculateAcc 方法 ")

[CheckReceiveDataComplete 方法](../html/b200e06c-38b7-9656-a340-54aa2373bb0a.htm "CheckReceiveDataComplete 方法 ")

[CheckResponse 方法](../html/1b4e1f3e-399f-4381-f912-86708faee605.htm "CheckResponse 方法 ")

[ExtraResponse 方法](../html/ad08cd70-4590-fc5d-a80d-a445450f9631.htm "ExtraResponse 方法 ")

[GetErrorDescriptionFromCode 方法](../html/fca0a278-c68e-cfe1-fcf4-816efcc878cb.htm "GetErrorDescriptionFromCode 方法 ")

[PackFatekCommand 方法](../html/fc22707e-222f-8fec-8d57-83742d4f17c1.htm "PackFatekCommand 方法 ")

[Read 方法](../html/cc1b8d93-c06c-41b9-748f-aa824b5f223b.htm "Read 方法 ")

[ReadAsync 方法](../html/c1a5498b-193d-f3bb-3e3c-7ab9585164ed.htm "ReadAsync 方法 ")

[ReadBool 方法](../html/5f136473-54c4-ce95-34f0-39c1028372f2.htm "ReadBool 方法 ")

[ReadBoolAsync 方法](../html/997a3300-c2ad-8816-bdfc-a2e0f9532bb4.htm "ReadBoolAsync 方法 ")

[ReadStatus 方法](../html/a73fd644-bf59-ca9a-d0ff-45d177c75e64.htm "ReadStatus 方法 ")

[ReadStatusAsync 方法](../html/62d31aab-0b53-ce2f-ffa5-09b3e4f85675.htm "ReadStatusAsync 方法 ")

[Run 方法](../html/d19f36a9-9645-fc5f-9bf8-7de6a709ed47.htm "Run 方法 ")

[RunAsync 方法](../html/0c1b1b0e-cecd-d2a2-5fb1-46a8a0321c39.htm "RunAsync 方法 ")

[Stop 方法](../html/2fc1d7f4-d7f2-ee5a-b271-69c33c1c3479.htm "Stop 方法 ")

[StopAsync 方法](../html/cfee21d4-ecbd-ea41-89f0-03d4f71102c7.htm "StopAsync 方法 ")

[Write 方法](../html/5adce02f-a40e-f29a-98db-23fb5791b508.htm "Write 方法 ")

[WriteAsync 方法](../html/4b5310b6-d2c7-d1ac-d43d-0c29bbbfea3b.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FatekProgramHelperStop 方法 |

使PLC处于STOP状态

**命名空间：**
 [HslCommunication.Profinet.FATEK.Helper](4ac2b39a-b5b1-9027-86a7-eadd294bcffa.htm)  
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
	IReadWriteDevice device,
	byte station
)
```

```
Public Shared Function Stop ( 
	device As IReadWriteDevice,
	station As Byte
) As OperateResult
```

```
public:
static OperateResult^ Stop(
	IReadWriteDevice^ device, 
	unsigned char station
)
```

```
static member Stop : 
        device : IReadWriteDevice * 
        station : byte -> OperateResult 
```

#### 参数

device
:   类型：[HslCommunication.CoreIReadWriteDevice](d987afcc-f99c-6aed-3847-5b1890504635.htm)  
    PLC通信对象

station
:   类型：SystemByte  
    站号信息

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否操作成功

![](../icons/SectionExpanded.png)参见

#### 引用

[FatekProgramHelper 类](5aa26a25-055f-eaea-5db8-af91b52b1e7b.htm)

[HslCommunication.Profinet.FATEK.Helper 命名空间](4ac2b39a-b5b1-9027-86a7-eadd294bcffa.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## StopAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/cfee21d4-ecbd-ea41-89f0-03d4f71102c7.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.FATEK.Helper](../html/4ac2b39a-b5b1-9027-86a7-eadd294bcffa.htm "HslCommunication.Profinet.FATEK.Helper")

[FatekProgramHelper 类](../html/5aa26a25-055f-eaea-5db8-af91b52b1e7b.htm "FatekProgramHelper 类")

[FatekProgramHelper 方法](../html/481082d5-89ec-7b3f-860f-ced4eec5c714.htm "FatekProgramHelper 方法")

[BuildReadBoolCommand 方法](../html/1f0e3930-99a4-b02f-d0b6-e139660f883a.htm "BuildReadBoolCommand 方法 ")

[BuildReadWordCommand 方法](../html/805ff280-0163-7ccb-25eb-2458424cc1ec.htm "BuildReadWordCommand 方法 ")

[BuildWriteBoolCommand 方法](../html/eecc01ad-74fc-3610-4939-9bbf87fc483d.htm "BuildWriteBoolCommand 方法 ")

[BuildWriteByteCommand 方法](../html/7b615285-6728-1344-384e-bf7cf918abf7.htm "BuildWriteByteCommand 方法 ")

[CalculateAcc 方法](../html/49b24ae3-78e3-e4d6-4743-f531b867da87.htm "CalculateAcc 方法 ")

[CheckReceiveDataComplete 方法](../html/b200e06c-38b7-9656-a340-54aa2373bb0a.htm "CheckReceiveDataComplete 方法 ")

[CheckResponse 方法](../html/1b4e1f3e-399f-4381-f912-86708faee605.htm "CheckResponse 方法 ")

[ExtraResponse 方法](../html/ad08cd70-4590-fc5d-a80d-a445450f9631.htm "ExtraResponse 方法 ")

[GetErrorDescriptionFromCode 方法](../html/fca0a278-c68e-cfe1-fcf4-816efcc878cb.htm "GetErrorDescriptionFromCode 方法 ")

[PackFatekCommand 方法](../html/fc22707e-222f-8fec-8d57-83742d4f17c1.htm "PackFatekCommand 方法 ")

[Read 方法](../html/cc1b8d93-c06c-41b9-748f-aa824b5f223b.htm "Read 方法 ")

[ReadAsync 方法](../html/c1a5498b-193d-f3bb-3e3c-7ab9585164ed.htm "ReadAsync 方法 ")

[ReadBool 方法](../html/5f136473-54c4-ce95-34f0-39c1028372f2.htm "ReadBool 方法 ")

[ReadBoolAsync 方法](../html/997a3300-c2ad-8816-bdfc-a2e0f9532bb4.htm "ReadBoolAsync 方法 ")

[ReadStatus 方法](../html/a73fd644-bf59-ca9a-d0ff-45d177c75e64.htm "ReadStatus 方法 ")

[ReadStatusAsync 方法](../html/62d31aab-0b53-ce2f-ffa5-09b3e4f85675.htm "ReadStatusAsync 方法 ")

[Run 方法](../html/d19f36a9-9645-fc5f-9bf8-7de6a709ed47.htm "Run 方法 ")

[RunAsync 方法](../html/0c1b1b0e-cecd-d2a2-5fb1-46a8a0321c39.htm "RunAsync 方法 ")

[Stop 方法](../html/2fc1d7f4-d7f2-ee5a-b271-69c33c1c3479.htm "Stop 方法 ")

[StopAsync 方法](../html/cfee21d4-ecbd-ea41-89f0-03d4f71102c7.htm "StopAsync 方法 ")

[Write 方法](../html/5adce02f-a40e-f29a-98db-23fb5791b508.htm "Write 方法 ")

[WriteAsync 方法](../html/4b5310b6-d2c7-d1ac-d43d-0c29bbbfea3b.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FatekProgramHelperStopAsync 方法 |

使PLC处于STOP状态

**命名空间：**
 [HslCommunication.Profinet.FATEK.Helper](4ac2b39a-b5b1-9027-86a7-eadd294bcffa.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static Task<OperateResult> StopAsync(
	IReadWriteDevice device,
	byte station
)
```

```
Public Shared Function StopAsync ( 
	device As IReadWriteDevice,
	station As Byte
) As Task(Of OperateResult)
```

```
public:
static Task<OperateResult^>^ StopAsync(
	IReadWriteDevice^ device, 
	unsigned char station
)
```

```
static member StopAsync : 
        device : IReadWriteDevice * 
        station : byte -> Task<OperateResult> 
```

#### 参数

device
:   类型：[HslCommunication.CoreIReadWriteDevice](d987afcc-f99c-6aed-3847-5b1890504635.htm)  
    PLC通信对象

station
:   类型：SystemByte  
    站号信息

#### 返回值

类型：Task[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否操作成功

![](../icons/SectionExpanded.png)参见

#### 引用

[FatekProgramHelper 类](5aa26a25-055f-eaea-5db8-af91b52b1e7b.htm)

[HslCommunication.Profinet.FATEK.Helper 命名空间](4ac2b39a-b5b1-9027-86a7-eadd294bcffa.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Write 方法 

[原文連結](http://api.hslcommunication.cn/html/5adce02f-a40e-f29a-98db-23fb5791b508.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.FATEK.Helper](../html/4ac2b39a-b5b1-9027-86a7-eadd294bcffa.htm "HslCommunication.Profinet.FATEK.Helper")

[FatekProgramHelper 类](../html/5aa26a25-055f-eaea-5db8-af91b52b1e7b.htm "FatekProgramHelper 类")

[FatekProgramHelper 方法](../html/481082d5-89ec-7b3f-860f-ced4eec5c714.htm "FatekProgramHelper 方法")

[Write 方法](../html/5adce02f-a40e-f29a-98db-23fb5791b508.htm "Write 方法 ")

[Write 方法 (IReadWriteDevice, Byte, String, Boolean[])](../html/dfabf156-9540-d08c-da8c-954c185417ed.htm "Write 方法 (IReadWriteDevice, Byte, String, Boolean[])")

[Write 方法 (IReadWriteDevice, Byte, String, Byte[])](../html/1d873a33-51ee-e065-1aac-efeb6b37a452.htm "Write 方法 (IReadWriteDevice, Byte, String, Byte[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FatekProgramHelperWrite 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法静态成员 | [Write(IReadWriteDevice, Byte, String, Boolean)](dfabf156-9540-d08c-da8c-954c185417ed.htm) | 批量写入bool类型的数组，支持的类型为X,Y,M,S,T,C，具体的地址范围取决于PLC的类型，地址可以携带站号信息，例如 s=2;M100  Write arrays of type bool in batches. The supported types are X, Y, M, S, T, C. The specific address range depends on the type of PLC, The address can carry station number information, such as s=2;M100 |
| 公共方法静态成员 | [Write(IReadWriteDevice, Byte, String, Byte)](1d873a33-51ee-e065-1aac-efeb6b37a452.htm) | 批量写入PLC的数据，以字为单位，也就是说最少2个字节信息，支持X,Y,M,S,D,T,C,R,RT,RC具体的地址范围需要根据PLC型号来确认，地址可以携带站号信息，例如 s=2;D100  The data written to the PLC in batches, in units of words, that is, at least 2 bytes of information, supporting X, Y, M, S, D, T, C, R, RT, and RC. The specific address range needs to be based on the PLC model To confirm, The address can carry station number information, such as s=2;D100 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[FatekProgramHelper 类](5aa26a25-055f-eaea-5db8-af91b52b1e7b.htm)

[HslCommunication.Profinet.FATEK.Helper 命名空间](4ac2b39a-b5b1-9027-86a7-eadd294bcffa.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Write 方法 (IReadWriteDevice, Byte, String, Boolean[])

[原文連結](http://api.hslcommunication.cn/html/dfabf156-9540-d08c-da8c-954c185417ed.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.FATEK.Helper](../html/4ac2b39a-b5b1-9027-86a7-eadd294bcffa.htm "HslCommunication.Profinet.FATEK.Helper")

[FatekProgramHelper 类](../html/5aa26a25-055f-eaea-5db8-af91b52b1e7b.htm "FatekProgramHelper 类")

[FatekProgramHelper 方法](../html/481082d5-89ec-7b3f-860f-ced4eec5c714.htm "FatekProgramHelper 方法")

[Write 方法](../html/5adce02f-a40e-f29a-98db-23fb5791b508.htm "Write 方法 ")

[Write 方法 (IReadWriteDevice, Byte, String, Boolean[])](../html/dfabf156-9540-d08c-da8c-954c185417ed.htm "Write 方法 (IReadWriteDevice, Byte, String, Boolean[])")

[Write 方法 (IReadWriteDevice, Byte, String, Byte[])](../html/1d873a33-51ee-e065-1aac-efeb6b37a452.htm "Write 方法 (IReadWriteDevice, Byte, String, Byte[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FatekProgramHelperWrite 方法 (IReadWriteDevice, Byte, String, Boolean) |

批量写入bool类型的数组，支持的类型为X,Y,M,S,T,C，具体的地址范围取决于PLC的类型，地址可以携带站号信息，例如 s=2;M100  
Write arrays of type bool in batches. The supported types are X, Y, M, S, T, C. The specific address range depends on the type of PLC,
The address can carry station number information, such as s=2;M100

**命名空间：**
 [HslCommunication.Profinet.FATEK.Helper](4ac2b39a-b5b1-9027-86a7-eadd294bcffa.htm)  
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
	bool[] value
)
```

```
Public Shared Function Write ( 
	device As IReadWriteDevice,
	station As Byte,
	address As String,
	value As Boolean()
) As OperateResult
```

```
public:
static OperateResult^ Write(
	IReadWriteDevice^ device, 
	unsigned char station, 
	String^ address, 
	array<bool>^ value
)
```

```
static member Write : 
        device : IReadWriteDevice * 
        station : byte * 
        address : string * 
        value : bool[] -> OperateResult 
```

#### 参数

device
:   类型：[HslCommunication.CoreIReadWriteDevice](d987afcc-f99c-6aed-3847-5b1890504635.htm)  
    PLC通信对象

station
:   类型：SystemByte  
    站号信息

address
:   类型：SystemString  
    PLC的地址信息

value
:   类型：SystemBoolean  
    数据信息

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否写入成功

![](../icons/SectionExpanded.png)参见

#### 引用

[FatekProgramHelper 类](5aa26a25-055f-eaea-5db8-af91b52b1e7b.htm)

[Write 重载](5adce02f-a40e-f29a-98db-23fb5791b508.htm)

[HslCommunication.Profinet.FATEK.Helper 命名空间](4ac2b39a-b5b1-9027-86a7-eadd294bcffa.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Write 方法 (IReadWriteDevice, Byte, String, Byte[])

[原文連結](http://api.hslcommunication.cn/html/1d873a33-51ee-e065-1aac-efeb6b37a452.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.FATEK.Helper](../html/4ac2b39a-b5b1-9027-86a7-eadd294bcffa.htm "HslCommunication.Profinet.FATEK.Helper")

[FatekProgramHelper 类](../html/5aa26a25-055f-eaea-5db8-af91b52b1e7b.htm "FatekProgramHelper 类")

[FatekProgramHelper 方法](../html/481082d5-89ec-7b3f-860f-ced4eec5c714.htm "FatekProgramHelper 方法")

[Write 方法](../html/5adce02f-a40e-f29a-98db-23fb5791b508.htm "Write 方法 ")

[Write 方法 (IReadWriteDevice, Byte, String, Boolean[])](../html/dfabf156-9540-d08c-da8c-954c185417ed.htm "Write 方法 (IReadWriteDevice, Byte, String, Boolean[])")

[Write 方法 (IReadWriteDevice, Byte, String, Byte[])](../html/1d873a33-51ee-e065-1aac-efeb6b37a452.htm "Write 方法 (IReadWriteDevice, Byte, String, Byte[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FatekProgramHelperWrite 方法 (IReadWriteDevice, Byte, String, Byte) |

批量写入PLC的数据，以字为单位，也就是说最少2个字节信息，支持X,Y,M,S,D,T,C,R,RT,RC具体的地址范围需要根据PLC型号来确认，地址可以携带站号信息，例如 s=2;D100  
The data written to the PLC in batches, in units of words, that is, at least 2 bytes of information,
supporting X, Y, M, S, D, T, C, R, RT, and RC. The specific address range needs to be based on the PLC model To confirm, The address can carry station number information, such as s=2;D100

**命名空间：**
 [HslCommunication.Profinet.FATEK.Helper](4ac2b39a-b5b1-9027-86a7-eadd294bcffa.htm)  
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
	byte[] value
)
```

```
Public Shared Function Write ( 
	device As IReadWriteDevice,
	station As Byte,
	address As String,
	value As Byte()
) As OperateResult
```

```
public:
static OperateResult^ Write(
	IReadWriteDevice^ device, 
	unsigned char station, 
	String^ address, 
	array<unsigned char>^ value
)
```

```
static member Write : 
        device : IReadWriteDevice * 
        station : byte * 
        address : string * 
        value : byte[] -> OperateResult 
```

#### 参数

device
:   类型：[HslCommunication.CoreIReadWriteDevice](d987afcc-f99c-6aed-3847-5b1890504635.htm)  
    PLC通信的对象

station
:   类型：SystemByte  
    设备的站号信息

address
:   类型：SystemString  
    地址信息，举例，D100，R200，RC100，RT200

value
:   类型：SystemByte  
    数据值

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否写入成功

![](../icons/SectionExpanded.png)参见

#### 引用

[FatekProgramHelper 类](5aa26a25-055f-eaea-5db8-af91b52b1e7b.htm)

[Write 重载](5adce02f-a40e-f29a-98db-23fb5791b508.htm)

[HslCommunication.Profinet.FATEK.Helper 命名空间](4ac2b39a-b5b1-9027-86a7-eadd294bcffa.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## WriteAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/4b5310b6-d2c7-d1ac-d43d-0c29bbbfea3b.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.FATEK.Helper](../html/4ac2b39a-b5b1-9027-86a7-eadd294bcffa.htm "HslCommunication.Profinet.FATEK.Helper")

[FatekProgramHelper 类](../html/5aa26a25-055f-eaea-5db8-af91b52b1e7b.htm "FatekProgramHelper 类")

[FatekProgramHelper 方法](../html/481082d5-89ec-7b3f-860f-ced4eec5c714.htm "FatekProgramHelper 方法")

[WriteAsync 方法](../html/4b5310b6-d2c7-d1ac-d43d-0c29bbbfea3b.htm "WriteAsync 方法 ")

[WriteAsync 方法 (IReadWriteDevice, Byte, String, Boolean[])](../html/e415bcb2-0cbc-6f93-151a-93d96112f9d9.htm "WriteAsync 方法 (IReadWriteDevice, Byte, String, Boolean[])")

[WriteAsync 方法 (IReadWriteDevice, Byte, String, Byte[])](../html/3754fd86-bb2c-0ca3-e736-279e77905479.htm "WriteAsync 方法 (IReadWriteDevice, Byte, String, Byte[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FatekProgramHelperWriteAsync 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法静态成员 | [WriteAsync(IReadWriteDevice, Byte, String, Boolean)](e415bcb2-0cbc-6f93-151a-93d96112f9d9.htm) | 批量写入bool类型的数组，支持的类型为X,Y,M,S,T,C，具体的地址范围取决于PLC的类型，地址可以携带站号信息，例如 s=2;M100  Write arrays of type bool in batches. The supported types are X, Y, M, S, T, C. The specific address range depends on the type of PLC, The address can carry station number information, such as s=2;M100 |
| 公共方法静态成员 | [WriteAsync(IReadWriteDevice, Byte, String, Byte)](3754fd86-bb2c-0ca3-e736-279e77905479.htm) | 批量写入PLC的数据，以字为单位，也就是说最少2个字节信息，支持X,Y,M,S,D,T,C,R,RT,RC具体的地址范围需要根据PLC型号来确认，地址可以携带站号信息，例如 s=2;D100  The data written to the PLC in batches, in units of words, that is, at least 2 bytes of information, supporting X, Y, M, S, D, T, C, R, RT, and RC. The specific address range needs to be based on the PLC model To confirm, The address can carry station number information, such as s=2;D100 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[FatekProgramHelper 类](5aa26a25-055f-eaea-5db8-af91b52b1e7b.htm)

[HslCommunication.Profinet.FATEK.Helper 命名空间](4ac2b39a-b5b1-9027-86a7-eadd294bcffa.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## WriteAsync 方法 (IReadWriteDevice, Byte, String, Boolean[])

[原文連結](http://api.hslcommunication.cn/html/e415bcb2-0cbc-6f93-151a-93d96112f9d9.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.FATEK.Helper](../html/4ac2b39a-b5b1-9027-86a7-eadd294bcffa.htm "HslCommunication.Profinet.FATEK.Helper")

[FatekProgramHelper 类](../html/5aa26a25-055f-eaea-5db8-af91b52b1e7b.htm "FatekProgramHelper 类")

[FatekProgramHelper 方法](../html/481082d5-89ec-7b3f-860f-ced4eec5c714.htm "FatekProgramHelper 方法")

[WriteAsync 方法](../html/4b5310b6-d2c7-d1ac-d43d-0c29bbbfea3b.htm "WriteAsync 方法 ")

[WriteAsync 方法 (IReadWriteDevice, Byte, String, Boolean[])](../html/e415bcb2-0cbc-6f93-151a-93d96112f9d9.htm "WriteAsync 方法 (IReadWriteDevice, Byte, String, Boolean[])")

[WriteAsync 方法 (IReadWriteDevice, Byte, String, Byte[])](../html/3754fd86-bb2c-0ca3-e736-279e77905479.htm "WriteAsync 方法 (IReadWriteDevice, Byte, String, Byte[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FatekProgramHelperWriteAsync 方法 (IReadWriteDevice, Byte, String, Boolean) |

批量写入bool类型的数组，支持的类型为X,Y,M,S,T,C，具体的地址范围取决于PLC的类型，地址可以携带站号信息，例如 s=2;M100  
Write arrays of type bool in batches. The supported types are X, Y, M, S, T, C. The specific address range depends on the type of PLC,
The address can carry station number information, such as s=2;M100

**命名空间：**
 [HslCommunication.Profinet.FATEK.Helper](4ac2b39a-b5b1-9027-86a7-eadd294bcffa.htm)  
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
	bool[] value
)
```

```
Public Shared Function WriteAsync ( 
	device As IReadWriteDevice,
	station As Byte,
	address As String,
	value As Boolean()
) As Task(Of OperateResult)
```

```
public:
static Task<OperateResult^>^ WriteAsync(
	IReadWriteDevice^ device, 
	unsigned char station, 
	String^ address, 
	array<bool>^ value
)
```

```
static member WriteAsync : 
        device : IReadWriteDevice * 
        station : byte * 
        address : string * 
        value : bool[] -> Task<OperateResult> 
```

#### 参数

device
:   类型：[HslCommunication.CoreIReadWriteDevice](d987afcc-f99c-6aed-3847-5b1890504635.htm)  
    PLC通信对象

station
:   类型：SystemByte  
    站号信息

address
:   类型：SystemString  
    PLC的地址信息

value
:   类型：SystemBoolean  
    数据信息

#### 返回值

类型：Task[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否写入成功

![](../icons/SectionExpanded.png)参见

#### 引用

[FatekProgramHelper 类](5aa26a25-055f-eaea-5db8-af91b52b1e7b.htm)

[WriteAsync 重载](4b5310b6-d2c7-d1ac-d43d-0c29bbbfea3b.htm)

[HslCommunication.Profinet.FATEK.Helper 命名空间](4ac2b39a-b5b1-9027-86a7-eadd294bcffa.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## WriteAsync 方法 (IReadWriteDevice, Byte, String, Byte[])

[原文連結](http://api.hslcommunication.cn/html/3754fd86-bb2c-0ca3-e736-279e77905479.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.FATEK.Helper](../html/4ac2b39a-b5b1-9027-86a7-eadd294bcffa.htm "HslCommunication.Profinet.FATEK.Helper")

[FatekProgramHelper 类](../html/5aa26a25-055f-eaea-5db8-af91b52b1e7b.htm "FatekProgramHelper 类")

[FatekProgramHelper 方法](../html/481082d5-89ec-7b3f-860f-ced4eec5c714.htm "FatekProgramHelper 方法")

[WriteAsync 方法](../html/4b5310b6-d2c7-d1ac-d43d-0c29bbbfea3b.htm "WriteAsync 方法 ")

[WriteAsync 方法 (IReadWriteDevice, Byte, String, Boolean[])](../html/e415bcb2-0cbc-6f93-151a-93d96112f9d9.htm "WriteAsync 方法 (IReadWriteDevice, Byte, String, Boolean[])")

[WriteAsync 方法 (IReadWriteDevice, Byte, String, Byte[])](../html/3754fd86-bb2c-0ca3-e736-279e77905479.htm "WriteAsync 方法 (IReadWriteDevice, Byte, String, Byte[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FatekProgramHelperWriteAsync 方法 (IReadWriteDevice, Byte, String, Byte) |

批量写入PLC的数据，以字为单位，也就是说最少2个字节信息，支持X,Y,M,S,D,T,C,R,RT,RC具体的地址范围需要根据PLC型号来确认，地址可以携带站号信息，例如 s=2;D100  
The data written to the PLC in batches, in units of words, that is, at least 2 bytes of information,
supporting X, Y, M, S, D, T, C, R, RT, and RC. The specific address range needs to be based on the PLC model To confirm, The address can carry station number information, such as s=2;D100

**命名空间：**
 [HslCommunication.Profinet.FATEK.Helper](4ac2b39a-b5b1-9027-86a7-eadd294bcffa.htm)  
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
	byte[] value
)
```

```
Public Shared Function WriteAsync ( 
	device As IReadWriteDevice,
	station As Byte,
	address As String,
	value As Byte()
) As Task(Of OperateResult)
```

```
public:
static Task<OperateResult^>^ WriteAsync(
	IReadWriteDevice^ device, 
	unsigned char station, 
	String^ address, 
	array<unsigned char>^ value
)
```

```
static member WriteAsync : 
        device : IReadWriteDevice * 
        station : byte * 
        address : string * 
        value : byte[] -> Task<OperateResult> 
```

#### 参数

device
:   类型：[HslCommunication.CoreIReadWriteDevice](d987afcc-f99c-6aed-3847-5b1890504635.htm)  
    PLC通信的对象

station
:   类型：SystemByte  
    设备的站号信息

address
:   类型：SystemString  
    地址信息，举例，D100，R200，RC100，RT200

value
:   类型：SystemByte  
    数据值

#### 返回值

类型：Task[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否写入成功

![](../icons/SectionExpanded.png)参见

#### 引用

[FatekProgramHelper 类](5aa26a25-055f-eaea-5db8-af91b52b1e7b.htm)

[WriteAsync 重载](4b5310b6-d2c7-d1ac-d43d-0c29bbbfea3b.htm)

[HslCommunication.Profinet.FATEK.Helper 命名空间](4ac2b39a-b5b1-9027-86a7-eadd294bcffa.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)