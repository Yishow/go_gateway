# HslCommunication - HslCommunication.CNC.Fanuc

> 分類頁數: 30



---
## HslCommunication.CNC.Fanuc

[原文連結](http://api.hslcommunication.cn/html/e11b249e-4a41-3205-3921-50085ef8c1c6.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.CNC.Fanuc](../html/e11b249e-4a41-3205-3921-50085ef8c1c6.htm "HslCommunication.CNC.Fanuc")

[CNCFanucSeriesMessage 类](../html/67bc0fef-2075-8eea-8ef0-bbcaea12932d.htm "CNCFanucSeriesMessage 类")

[CNCRunStatus 枚举](../html/375e9d94-4983-d199-3adc-00efde97fe5c.htm "CNCRunStatus 枚举")

[CNCWorkMode 枚举](../html/fac6694f-6128-43b9-6056-e4a07201a496.htm "CNCWorkMode 枚举")

[CutterInfo 类](../html/41023477-96be-7b95-3000-fb98b3d28e50.htm "CutterInfo 类")

[FanucOperatorMessage 类](../html/dc495da4-b0c1-f6b7-c86c-2f470d1871d6.htm "FanucOperatorMessage 类")

[FanucSeries0i 类](../html/25cac34a-b90e-b589-124d-e0d86bcf6fe3.htm "FanucSeries0i 类")

[FanucSysInfo 类](../html/aca9fd72-2c26-6b6f-4f88-f8139659f31e.htm "FanucSysInfo 类")

[FileDirInfo 类](../html/017b5d5b-06d6-861a-e780-c1b4f3435622.htm "FileDirInfo 类")

[SysAlarm 类](../html/e2923ffd-34d2-a8a5-69f3-95d3aad14cf5.htm "SysAlarm 类")

[SysAllCoors 类](../html/68c07730-3736-56a9-a8ad-2306d7704af1.htm "SysAllCoors 类")

[SysStatusInfo 类](../html/3976aa62-a2d4-15cf-e8cc-0cde91ab2213.htm "SysStatusInfo 类")

[ToolInformation 类](../html/f4315993-25ca-bd5e-d41d-f571a01c0167.htm "ToolInformation 类")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslCommunication.CNC.Fanuc 命名空间 |

[缺少 "N:HslCommunication.CNC.Fanuc" 的 <summary> 文档]

![](../icons/SectionExpanded.png)类

|  | 类 | 说明 |
| --- | --- | --- |
| 公共类 | [CNCFanucSeriesMessage](67bc0fef-2075-8eea-8ef0-bbcaea12932d.htm) | Fanuc床子的消息对象 |
| 公共类 | [CutterInfo](41023477-96be-7b95-3000-fb98b3d28e50.htm) | 刀具信息 |
| 公共类 | [FanucOperatorMessage](dc495da4-b0c1-f6b7-c86c-2f470d1871d6.htm) | Fanuc机床的操作信息 |
| 公共类 | [FanucSeries0i](25cac34a-b90e-b589-124d-e0d86bcf6fe3.htm) | 一个FANUC的机床通信类对象 |
| 公共类 | [FanucSysInfo](aca9fd72-2c26-6b6f-4f88-f8139659f31e.htm) | Fanuc的系统信息 |
| 公共类 | [FileDirInfo](017b5d5b-06d6-861a-e780-c1b4f3435622.htm) | 文件或是文件夹的信息 |
| 公共类 | [SysAlarm](e2923ffd-34d2-a8a5-69f3-95d3aad14cf5.htm) | 当前机床的报警信息 |
| 公共类 | [SysAllCoors](68c07730-3736-56a9-a8ad-2306d7704af1.htm) | 系统的坐标信息 |
| 公共类 | [SysStatusInfo](3976aa62-a2d4-15cf-e8cc-0cde91ab2213.htm) | 系统状态信息 |
| 公共类 | [ToolInformation](f4315993-25ca-bd5e-d41d-f571a01c0167.htm) | 刀具信息 |

![](../icons/SectionExpanded.png)枚举

|  | 枚举 | 说明 |
| --- | --- | --- |
| 公共枚举 | [CNCRunStatus](375e9d94-4983-d199-3adc-00efde97fe5c.htm) | CNC的运行状态 |
| 公共枚举 | [CNCWorkMode](fac6694f-6128-43b9-6056-e4a07201a496.htm) | 设备的工作模式 |

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## CNCFanucSeriesMessage 类

[原文連結](http://api.hslcommunication.cn/html/67bc0fef-2075-8eea-8ef0-bbcaea12932d.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.CNC.Fanuc](../html/e11b249e-4a41-3205-3921-50085ef8c1c6.htm "HslCommunication.CNC.Fanuc")

[CNCFanucSeriesMessage 类](../html/67bc0fef-2075-8eea-8ef0-bbcaea12932d.htm "CNCFanucSeriesMessage 类")

[CNCFanucSeriesMessage 构造函数](../html/cfafaddd-31a6-25ed-58f8-330b887236fc.htm "CNCFanucSeriesMessage 构造函数 ")

[CNCFanucSeriesMessage 属性](../html/f9805794-4b15-2bc9-c3aa-11ac70995a6a.htm "CNCFanucSeriesMessage 属性")

[CNCFanucSeriesMessage 方法](../html/065eacd7-1523-b6ab-24d5-3bbe8abec9a7.htm "CNCFanucSeriesMessage 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CNCFanucSeriesMessage 类 |

Fanuc床子的消息对象

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  [HslCommunication.Core.IMessageNetMessageBase](55e9c192-5d9b-9d92-31be-b05cd76d6903.htm)  
    HslCommunication.CNC.FanucCNCFanucSeriesMessage

**命名空间：**
 [HslCommunication.CNC.Fanuc](e11b249e-4a41-3205-3921-50085ef8c1c6.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class CNCFanucSeriesMessage : NetMessageBase, 
	INetMessage
```

```
Public Class CNCFanucSeriesMessage
	Inherits NetMessageBase
	Implements INetMessage
```

```
public ref class CNCFanucSeriesMessage : public NetMessageBase, 
	INetMessage
```

```
type CNCFanucSeriesMessage =  
    class
        inherit NetMessageBase
        interface INetMessage
    end
```

CNCFanucSeriesMessage 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [CNCFanucSeriesMessage](cfafaddd-31a6-25ed-58f8-330b887236fc.htm) | 初始化 CNCFanucSeriesMessage 类的一个新实例 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [ContentBytes](27f34824-6adc-b126-b622-6f9eeaabe44f.htm) | 消息内容字节  Message content byte (继承自 [NetMessageBase](55e9c192-5d9b-9d92-31be-b05cd76d6903.htm)。) |
| 公共属性 | [HeadBytes](ee605527-7bfc-97a3-52b8-620d839d3d4c.htm) | 消息头字节  Message header byte (继承自 [NetMessageBase](55e9c192-5d9b-9d92-31be-b05cd76d6903.htm)。) |
| 公共属性 | [ProtocolHeadBytesLength](c72f0070-8447-8e90-b01c-dd905acb27b0.htm) | 消息头的指令长度，第一次接受数据的长度  Instruction length of the message header, the length of the first received data |
| 公共属性 | [SendBytes](41476212-a254-3722-3445-9e636263d0f6.htm) | 发送的字节信息  Byte information sent (继承自 [NetMessageBase](55e9c192-5d9b-9d92-31be-b05cd76d6903.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [CheckHeadBytesLegal](3085e36e-54a0-a2ff-cc81-099c11d5b238.htm) | 检查头子节的合法性  Check the legitimacy of the head subsection (继承自 [NetMessageBase](55e9c192-5d9b-9d92-31be-b05cd76d6903.htm)。) |
| 公共方法 | [CheckMessageMatch](857a414e-7059-b002-b34c-10255390d6b7.htm) | 检查发送的接收的报文是否是匹配的，如果匹配，则返回 1, 如果不匹配且直接返回错误，则返回 0，如果不匹配继续接收，直到匹配或是超时，则返回 -1  If the packet is matched, 1 is returned. If the packet is not matched and an error is returned, 0 is returned. If the packet is not matched, -1 is returned until the packet is matched or times out (继承自 [NetMessageBase](55e9c192-5d9b-9d92-31be-b05cd76d6903.htm)。) |
| 公共方法 | [CheckReceiveDataComplete](909effc7-04ee-4941-593e-2da3f57e6201.htm) | 当消息头报文的长度定义为-1的时候，则使用动态的长度信息，可以使用本方法来判断一个消息是否处于完整的状态。  If the length of the message header is defined as -1, this method can be used to determine whether a message is in the complete state by using dynamic length information. (继承自 [NetMessageBase](55e9c192-5d9b-9d92-31be-b05cd76d6903.htm)。) |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | [GetContentLengthByHeadBytes](29d4fd2e-4546-0096-2631-5f534b0c710d.htm) | 从当前的头子节文件中提取出接下来需要接收的数据长度  Extract the length of the data to be received from the current header file |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法 | [GetHeadBytesIdentity](61d7c255-4473-31e7-7f84-dfd957a8826f.htm) | 获取头子节里的消息标识  Get the message ID in the header subsection (继承自 [NetMessageBase](55e9c192-5d9b-9d92-31be-b05cd76d6903.htm)。) |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法 | [PependedUselesByteLength](91167f2b-bebd-0102-29d4-b514cd3aa0a0.htm) | 在接收头报文的时候，返回前置无效的报文头字节长度，默认为0，不处理  When receiving a header message, return the header byte length of the invalid header, the default is 0, and no processing is performed. (继承自 [NetMessageBase](55e9c192-5d9b-9d92-31be-b05cd76d6903.htm)。) |
| 公共方法 | [ToString](d19532d7-e39b-e64a-972e-fa81b8b9a7a0.htm) | (重写 [NetMessageBaseToString](762f299e-5af6-4f6a-292c-2ae472bfe3b5.htm).) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.CNC.Fanuc 命名空间](e11b249e-4a41-3205-3921-50085ef8c1c6.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## CNCFanucSeriesMessage 构造函数 

[原文連結](http://api.hslcommunication.cn/html/cfafaddd-31a6-25ed-58f8-330b887236fc.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.CNC.Fanuc](../html/e11b249e-4a41-3205-3921-50085ef8c1c6.htm "HslCommunication.CNC.Fanuc")

[CNCFanucSeriesMessage 类](../html/67bc0fef-2075-8eea-8ef0-bbcaea12932d.htm "CNCFanucSeriesMessage 类")

[CNCFanucSeriesMessage 构造函数](../html/cfafaddd-31a6-25ed-58f8-330b887236fc.htm "CNCFanucSeriesMessage 构造函数 ")

[CNCFanucSeriesMessage 属性](../html/f9805794-4b15-2bc9-c3aa-11ac70995a6a.htm "CNCFanucSeriesMessage 属性")

[CNCFanucSeriesMessage 方法](../html/065eacd7-1523-b6ab-24d5-3bbe8abec9a7.htm "CNCFanucSeriesMessage 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CNCFanucSeriesMessage 构造函数 |

初始化 [CNCFanucSeriesMessage](67bc0fef-2075-8eea-8ef0-bbcaea12932d.htm) 类的一个新实例

**命名空间：**
 [HslCommunication.CNC.Fanuc](e11b249e-4a41-3205-3921-50085ef8c1c6.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public CNCFanucSeriesMessage()
```

```
Public Sub New
```

```
public:
CNCFanucSeriesMessage()
```

```
new : unit -> CNCFanucSeriesMessage
```

![](../icons/SectionExpanded.png)参见

#### 引用

[CNCFanucSeriesMessage 类](67bc0fef-2075-8eea-8ef0-bbcaea12932d.htm)

[HslCommunication.CNC.Fanuc 命名空间](e11b249e-4a41-3205-3921-50085ef8c1c6.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## CNCFanucSeriesMessage 属性

[原文連結](http://api.hslcommunication.cn/html/f9805794-4b15-2bc9-c3aa-11ac70995a6a.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.CNC.Fanuc](../html/e11b249e-4a41-3205-3921-50085ef8c1c6.htm "HslCommunication.CNC.Fanuc")

[CNCFanucSeriesMessage 类](../html/67bc0fef-2075-8eea-8ef0-bbcaea12932d.htm "CNCFanucSeriesMessage 类")

[CNCFanucSeriesMessage 属性](../html/f9805794-4b15-2bc9-c3aa-11ac70995a6a.htm "CNCFanucSeriesMessage 属性")

[ProtocolHeadBytesLength 属性](../html/c72f0070-8447-8e90-b01c-dd905acb27b0.htm "ProtocolHeadBytesLength 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CNCFanucSeriesMessage 属性 |

[CNCFanucSeriesMessage](67bc0fef-2075-8eea-8ef0-bbcaea12932d.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [ContentBytes](27f34824-6adc-b126-b622-6f9eeaabe44f.htm) | 消息内容字节  Message content byte (继承自 [NetMessageBase](55e9c192-5d9b-9d92-31be-b05cd76d6903.htm)。) |
| 公共属性 | [HeadBytes](ee605527-7bfc-97a3-52b8-620d839d3d4c.htm) | 消息头字节  Message header byte (继承自 [NetMessageBase](55e9c192-5d9b-9d92-31be-b05cd76d6903.htm)。) |
| 公共属性 | [ProtocolHeadBytesLength](c72f0070-8447-8e90-b01c-dd905acb27b0.htm) | 消息头的指令长度，第一次接受数据的长度  Instruction length of the message header, the length of the first received data |
| 公共属性 | [SendBytes](41476212-a254-3722-3445-9e636263d0f6.htm) | 发送的字节信息  Byte information sent (继承自 [NetMessageBase](55e9c192-5d9b-9d92-31be-b05cd76d6903.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[CNCFanucSeriesMessage 类](67bc0fef-2075-8eea-8ef0-bbcaea12932d.htm)

[HslCommunication.CNC.Fanuc 命名空间](e11b249e-4a41-3205-3921-50085ef8c1c6.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ProtocolHeadBytesLength 属性 

[原文連結](http://api.hslcommunication.cn/html/c72f0070-8447-8e90-b01c-dd905acb27b0.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.CNC.Fanuc](../html/e11b249e-4a41-3205-3921-50085ef8c1c6.htm "HslCommunication.CNC.Fanuc")

[CNCFanucSeriesMessage 类](../html/67bc0fef-2075-8eea-8ef0-bbcaea12932d.htm "CNCFanucSeriesMessage 类")

[CNCFanucSeriesMessage 属性](../html/f9805794-4b15-2bc9-c3aa-11ac70995a6a.htm "CNCFanucSeriesMessage 属性")

[ProtocolHeadBytesLength 属性](../html/c72f0070-8447-8e90-b01c-dd905acb27b0.htm "ProtocolHeadBytesLength 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CNCFanucSeriesMessageProtocolHeadBytesLength 属性 |

消息头的指令长度，第一次接受数据的长度  
Instruction length of the message header, the length of the first received data

**命名空间：**
 [HslCommunication.CNC.Fanuc](e11b249e-4a41-3205-3921-50085ef8c1c6.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public int ProtocolHeadBytesLength { get; }
```

```
Public ReadOnly Property ProtocolHeadBytesLength As Integer
	Get
```

```
public:
virtual property int ProtocolHeadBytesLength {
	int get () sealed;
}
```

```
abstract ProtocolHeadBytesLength : int with get
override ProtocolHeadBytesLength : int with get
```

#### 属性值

类型：Int32

#### 实现

[INetMessageProtocolHeadBytesLength](dec23ad8-cb9e-b1e2-86f5-fd4ecf5c4291.htm)

![](../icons/SectionExpanded.png)备注

当最高位字节的最高位为1时，第0和1位为校验的字符数量，第二高位字节表示结束字符之后的剩余字符长度信息，因为一个int占用四个字节，所以最多可以判断2个结束的字符信息。  
When the highest bit of the highest-order byte is 1, the 0th and 1st bits are the number of characters to be checked,
and the second high-order byte indicates the length information of the remaining characters after the end character.
Because one int occupies four bytes, the maximum It is possible to judge the character information of 2 ends.

![](../icons/SectionExpanded.png)参见

#### 引用

[CNCFanucSeriesMessage 类](67bc0fef-2075-8eea-8ef0-bbcaea12932d.htm)

[HslCommunication.CNC.Fanuc 命名空间](e11b249e-4a41-3205-3921-50085ef8c1c6.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## CNCFanucSeriesMessage 方法

[原文連結](http://api.hslcommunication.cn/html/065eacd7-1523-b6ab-24d5-3bbe8abec9a7.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.CNC.Fanuc](../html/e11b249e-4a41-3205-3921-50085ef8c1c6.htm "HslCommunication.CNC.Fanuc")

[CNCFanucSeriesMessage 类](../html/67bc0fef-2075-8eea-8ef0-bbcaea12932d.htm "CNCFanucSeriesMessage 类")

[CNCFanucSeriesMessage 方法](../html/065eacd7-1523-b6ab-24d5-3bbe8abec9a7.htm "CNCFanucSeriesMessage 方法")

[GetContentLengthByHeadBytes 方法](../html/29d4fd2e-4546-0096-2631-5f534b0c710d.htm "GetContentLengthByHeadBytes 方法 ")

[ToString 方法](../html/d19532d7-e39b-e64a-972e-fa81b8b9a7a0.htm "ToString 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CNCFanucSeriesMessage 方法 |

[CNCFanucSeriesMessage](67bc0fef-2075-8eea-8ef0-bbcaea12932d.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [CheckHeadBytesLegal](3085e36e-54a0-a2ff-cc81-099c11d5b238.htm) | 检查头子节的合法性  Check the legitimacy of the head subsection (继承自 [NetMessageBase](55e9c192-5d9b-9d92-31be-b05cd76d6903.htm)。) |
| 公共方法 | [CheckMessageMatch](857a414e-7059-b002-b34c-10255390d6b7.htm) | 检查发送的接收的报文是否是匹配的，如果匹配，则返回 1, 如果不匹配且直接返回错误，则返回 0，如果不匹配继续接收，直到匹配或是超时，则返回 -1  If the packet is matched, 1 is returned. If the packet is not matched and an error is returned, 0 is returned. If the packet is not matched, -1 is returned until the packet is matched or times out (继承自 [NetMessageBase](55e9c192-5d9b-9d92-31be-b05cd76d6903.htm)。) |
| 公共方法 | [CheckReceiveDataComplete](909effc7-04ee-4941-593e-2da3f57e6201.htm) | 当消息头报文的长度定义为-1的时候，则使用动态的长度信息，可以使用本方法来判断一个消息是否处于完整的状态。  If the length of the message header is defined as -1, this method can be used to determine whether a message is in the complete state by using dynamic length information. (继承自 [NetMessageBase](55e9c192-5d9b-9d92-31be-b05cd76d6903.htm)。) |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | [GetContentLengthByHeadBytes](29d4fd2e-4546-0096-2631-5f534b0c710d.htm) | 从当前的头子节文件中提取出接下来需要接收的数据长度  Extract the length of the data to be received from the current header file |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法 | [GetHeadBytesIdentity](61d7c255-4473-31e7-7f84-dfd957a8826f.htm) | 获取头子节里的消息标识  Get the message ID in the header subsection (继承自 [NetMessageBase](55e9c192-5d9b-9d92-31be-b05cd76d6903.htm)。) |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法 | [PependedUselesByteLength](91167f2b-bebd-0102-29d4-b514cd3aa0a0.htm) | 在接收头报文的时候，返回前置无效的报文头字节长度，默认为0，不处理  When receiving a header message, return the header byte length of the invalid header, the default is 0, and no processing is performed. (继承自 [NetMessageBase](55e9c192-5d9b-9d92-31be-b05cd76d6903.htm)。) |
| 公共方法 | [ToString](d19532d7-e39b-e64a-972e-fa81b8b9a7a0.htm) | (重写 [NetMessageBaseToString](762f299e-5af6-4f6a-292c-2ae472bfe3b5.htm).) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[CNCFanucSeriesMessage 类](67bc0fef-2075-8eea-8ef0-bbcaea12932d.htm)

[HslCommunication.CNC.Fanuc 命名空间](e11b249e-4a41-3205-3921-50085ef8c1c6.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GetContentLengthByHeadBytes 方法 

[原文連結](http://api.hslcommunication.cn/html/29d4fd2e-4546-0096-2631-5f534b0c710d.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.CNC.Fanuc](../html/e11b249e-4a41-3205-3921-50085ef8c1c6.htm "HslCommunication.CNC.Fanuc")

[CNCFanucSeriesMessage 类](../html/67bc0fef-2075-8eea-8ef0-bbcaea12932d.htm "CNCFanucSeriesMessage 类")

[CNCFanucSeriesMessage 方法](../html/065eacd7-1523-b6ab-24d5-3bbe8abec9a7.htm "CNCFanucSeriesMessage 方法")

[GetContentLengthByHeadBytes 方法](../html/29d4fd2e-4546-0096-2631-5f534b0c710d.htm "GetContentLengthByHeadBytes 方法 ")

[ToString 方法](../html/d19532d7-e39b-e64a-972e-fa81b8b9a7a0.htm "ToString 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CNCFanucSeriesMessageGetContentLengthByHeadBytes 方法 |

从当前的头子节文件中提取出接下来需要接收的数据长度  
Extract the length of the data to be received from the current header file

**命名空间：**
 [HslCommunication.CNC.Fanuc](e11b249e-4a41-3205-3921-50085ef8c1c6.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public int GetContentLengthByHeadBytes()
```

```
Public Function GetContentLengthByHeadBytes As Integer
```

```
public:
virtual int GetContentLengthByHeadBytes() sealed
```

```
abstract GetContentLengthByHeadBytes : unit -> int 
override GetContentLengthByHeadBytes : unit -> int
```

#### 返回值

类型：Int32  
返回接下来的数据内容长度

#### 实现

[INetMessageGetContentLengthByHeadBytes](a2e476e1-16e3-938c-6c58-419d08e8b40c.htm)

![](../icons/SectionExpanded.png)备注

如果剩余字节的长度小于0，则表示消息头数据还没有接收完整，还需要接收一定的长度(返回值的绝对值)，然后再判断剩余字节长度是否小于0，直到结果大于等于0为止，最多判断的次数为16次，超过16次将返回失败  
If the length of the remaining bytes is less than 0, it means that the message header data has not been received completely, and a certain length (absolute value of the return value) needs to be received,
and then it is judged whether the length of the remaining bytes is less than 0 until the result is greater than or equal to 0,
the maximum number of judgments is 16, more than 16 times will return failure

![](../icons/SectionExpanded.png)参见

#### 引用

[CNCFanucSeriesMessage 类](67bc0fef-2075-8eea-8ef0-bbcaea12932d.htm)

[HslCommunication.CNC.Fanuc 命名空间](e11b249e-4a41-3205-3921-50085ef8c1c6.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ToString 方法 

[原文連結](http://api.hslcommunication.cn/html/d19532d7-e39b-e64a-972e-fa81b8b9a7a0.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.CNC.Fanuc](../html/e11b249e-4a41-3205-3921-50085ef8c1c6.htm "HslCommunication.CNC.Fanuc")

[CNCFanucSeriesMessage 类](../html/67bc0fef-2075-8eea-8ef0-bbcaea12932d.htm "CNCFanucSeriesMessage 类")

[CNCFanucSeriesMessage 方法](../html/065eacd7-1523-b6ab-24d5-3bbe8abec9a7.htm "CNCFanucSeriesMessage 方法")

[GetContentLengthByHeadBytes 方法](../html/29d4fd2e-4546-0096-2631-5f534b0c710d.htm "GetContentLengthByHeadBytes 方法 ")

[ToString 方法](../html/d19532d7-e39b-e64a-972e-fa81b8b9a7a0.htm "ToString 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CNCFanucSeriesMessageToString 方法 |

[缺少 "M:HslCommunication.CNC.Fanuc.CNCFanucSeriesMessage.ToString" 的 <summary> 文档]

**命名空间：**
 [HslCommunication.CNC.Fanuc](e11b249e-4a41-3205-3921-50085ef8c1c6.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public override string ToString()
```

```
Public Overrides Function ToString As String
```

```
public:
virtual String^ ToString() override
```

```
abstract ToString : unit -> string 
override ToString : unit -> string
```

#### 返回值

类型：String  

[缺少 "M:HslCommunication.CNC.Fanuc.CNCFanucSeriesMessage.ToString" 的 <returns> 文档]

![](../icons/SectionExpanded.png)参见

#### 引用

[CNCFanucSeriesMessage 类](67bc0fef-2075-8eea-8ef0-bbcaea12932d.htm)

[HslCommunication.CNC.Fanuc 命名空间](e11b249e-4a41-3205-3921-50085ef8c1c6.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## CNCRunStatus 枚举

[原文連結](http://api.hslcommunication.cn/html/375e9d94-4983-d199-3adc-00efde97fe5c.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.CNC.Fanuc](../html/e11b249e-4a41-3205-3921-50085ef8c1c6.htm "HslCommunication.CNC.Fanuc")

[CNCFanucSeriesMessage 类](../html/67bc0fef-2075-8eea-8ef0-bbcaea12932d.htm "CNCFanucSeriesMessage 类")

[CNCRunStatus 枚举](../html/375e9d94-4983-d199-3adc-00efde97fe5c.htm "CNCRunStatus 枚举")

[CNCWorkMode 枚举](../html/fac6694f-6128-43b9-6056-e4a07201a496.htm "CNCWorkMode 枚举")

[CutterInfo 类](../html/41023477-96be-7b95-3000-fb98b3d28e50.htm "CutterInfo 类")

[FanucOperatorMessage 类](../html/dc495da4-b0c1-f6b7-c86c-2f470d1871d6.htm "FanucOperatorMessage 类")

[FanucSeries0i 类](../html/25cac34a-b90e-b589-124d-e0d86bcf6fe3.htm "FanucSeries0i 类")

[FanucSysInfo 类](../html/aca9fd72-2c26-6b6f-4f88-f8139659f31e.htm "FanucSysInfo 类")

[FileDirInfo 类](../html/017b5d5b-06d6-861a-e780-c1b4f3435622.htm "FileDirInfo 类")

[SysAlarm 类](../html/e2923ffd-34d2-a8a5-69f3-95d3aad14cf5.htm "SysAlarm 类")

[SysAllCoors 类](../html/68c07730-3736-56a9-a8ad-2306d7704af1.htm "SysAllCoors 类")

[SysStatusInfo 类](../html/3976aa62-a2d4-15cf-e8cc-0cde91ab2213.htm "SysStatusInfo 类")

[ToolInformation 类](../html/f4315993-25ca-bd5e-d41d-f571a01c0167.htm "ToolInformation 类")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CNCRunStatus 枚举 |

CNC的运行状态

**命名空间：**
 [HslCommunication.CNC.Fanuc](e11b249e-4a41-3205-3921-50085ef8c1c6.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public enum CNCRunStatus
```

```
Public Enumeration CNCRunStatus
```

```
public enum class CNCRunStatus
```

```
type CNCRunStatus
```

![](../icons/SectionExpanded.png)成员

|  | 成员名称 | 值 | 说明 |
| --- | --- | --- | --- |
|  | RESET | 0 | 重置 |
|  | STOP | 1 | 停止 |
|  | HOLD | 2 | 等待 |
|  | START | 3 | 启动 |
|  | MSTR | 4 | MSTR |

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.CNC.Fanuc 命名空间](e11b249e-4a41-3205-3921-50085ef8c1c6.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## CNCWorkMode 枚举

[原文連結](http://api.hslcommunication.cn/html/fac6694f-6128-43b9-6056-e4a07201a496.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.CNC.Fanuc](../html/e11b249e-4a41-3205-3921-50085ef8c1c6.htm "HslCommunication.CNC.Fanuc")

[CNCFanucSeriesMessage 类](../html/67bc0fef-2075-8eea-8ef0-bbcaea12932d.htm "CNCFanucSeriesMessage 类")

[CNCRunStatus 枚举](../html/375e9d94-4983-d199-3adc-00efde97fe5c.htm "CNCRunStatus 枚举")

[CNCWorkMode 枚举](../html/fac6694f-6128-43b9-6056-e4a07201a496.htm "CNCWorkMode 枚举")

[CutterInfo 类](../html/41023477-96be-7b95-3000-fb98b3d28e50.htm "CutterInfo 类")

[FanucOperatorMessage 类](../html/dc495da4-b0c1-f6b7-c86c-2f470d1871d6.htm "FanucOperatorMessage 类")

[FanucSeries0i 类](../html/25cac34a-b90e-b589-124d-e0d86bcf6fe3.htm "FanucSeries0i 类")

[FanucSysInfo 类](../html/aca9fd72-2c26-6b6f-4f88-f8139659f31e.htm "FanucSysInfo 类")

[FileDirInfo 类](../html/017b5d5b-06d6-861a-e780-c1b4f3435622.htm "FileDirInfo 类")

[SysAlarm 类](../html/e2923ffd-34d2-a8a5-69f3-95d3aad14cf5.htm "SysAlarm 类")

[SysAllCoors 类](../html/68c07730-3736-56a9-a8ad-2306d7704af1.htm "SysAllCoors 类")

[SysStatusInfo 类](../html/3976aa62-a2d4-15cf-e8cc-0cde91ab2213.htm "SysStatusInfo 类")

[ToolInformation 类](../html/f4315993-25ca-bd5e-d41d-f571a01c0167.htm "ToolInformation 类")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CNCWorkMode 枚举 |

设备的工作模式

**命名空间：**
 [HslCommunication.CNC.Fanuc](e11b249e-4a41-3205-3921-50085ef8c1c6.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public enum CNCWorkMode
```

```
Public Enumeration CNCWorkMode
```

```
public enum class CNCWorkMode
```

```
type CNCWorkMode
```

![](../icons/SectionExpanded.png)成员

|  | 成员名称 | 值 | 说明 |
| --- | --- | --- | --- |
|  | MDI | 0 | 手动输入 |
|  | AUTO | 1 | 自动循环 |
|  | EDIT | 3 | 程序编辑 |
|  | HANDLE | 4 | ×100 |
|  | JOG | 5 | 连续进给 |
|  | TeachInJOG | 6 | ??? |
|  | TeachInHandle | 7 | 示教 |
|  | INCfeed | 8 | ??? |
|  | REFerence | 9 | 机床回零 |
|  | ReMoTe | 10 | ??? |

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.CNC.Fanuc 命名空间](e11b249e-4a41-3205-3921-50085ef8c1c6.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## CutterInfo 类

[原文連結](http://api.hslcommunication.cn/html/41023477-96be-7b95-3000-fb98b3d28e50.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.CNC.Fanuc](../html/e11b249e-4a41-3205-3921-50085ef8c1c6.htm "HslCommunication.CNC.Fanuc")

[CutterInfo 类](../html/41023477-96be-7b95-3000-fb98b3d28e50.htm "CutterInfo 类")

[CutterInfo 构造函数](../html/db1562d8-ed25-1073-d95f-5b7f52e57c40.htm "CutterInfo 构造函数 ")

[CutterInfo 属性](../html/474119ab-6e0e-3b29-d1d9-f3ec883d47fd.htm "CutterInfo 属性")

[CutterInfo 方法](../html/dbc9c79b-50bc-fd8f-fe74-951e59d13222.htm "CutterInfo 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CutterInfo 类 |

刀具信息

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  HslCommunication.CNC.FanucCutterInfo

**命名空间：**
 [HslCommunication.CNC.Fanuc](e11b249e-4a41-3205-3921-50085ef8c1c6.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class CutterInfo
```

```
Public Class CutterInfo
```

```
public ref class CutterInfo
```

```
type CutterInfo =  class end
```

CutterInfo 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [CutterInfo](db1562d8-ed25-1073-d95f-5b7f52e57c40.htm) | 初始化 CutterInfo 类的一个新实例 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [LengthSharpOffset](0270102e-15a3-9b9c-142c-86cb0458f42b.htm) | 长度形状补偿 |
| 公共属性 | [LengthWearOffset](8d964911-252f-a649-ad4c-2bcc515b1bdd.htm) | 长度磨损补偿 |
| 公共属性 | [RadiusSharpOffset](4857b4aa-a44f-a04d-e87b-8a42cbdf4c61.htm) | 半径形状补偿 |
| 公共属性 | [RadiusWearOffset](f89ab2a6-b3cb-cca0-bad6-df8edf560440.htm) | 半径磨损补偿 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法 | [ToString](80661f5b-cea4-22a7-b661-21afdfb6e182.htm) | (重写 ObjectToString.) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.CNC.Fanuc 命名空间](e11b249e-4a41-3205-3921-50085ef8c1c6.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## CutterInfo 构造函数 

[原文連結](http://api.hslcommunication.cn/html/db1562d8-ed25-1073-d95f-5b7f52e57c40.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.CNC.Fanuc](../html/e11b249e-4a41-3205-3921-50085ef8c1c6.htm "HslCommunication.CNC.Fanuc")

[CutterInfo 类](../html/41023477-96be-7b95-3000-fb98b3d28e50.htm "CutterInfo 类")

[CutterInfo 构造函数](../html/db1562d8-ed25-1073-d95f-5b7f52e57c40.htm "CutterInfo 构造函数 ")

[CutterInfo 属性](../html/474119ab-6e0e-3b29-d1d9-f3ec883d47fd.htm "CutterInfo 属性")

[CutterInfo 方法](../html/dbc9c79b-50bc-fd8f-fe74-951e59d13222.htm "CutterInfo 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CutterInfo 构造函数 |

初始化 [CutterInfo](41023477-96be-7b95-3000-fb98b3d28e50.htm) 类的一个新实例

**命名空间：**
 [HslCommunication.CNC.Fanuc](e11b249e-4a41-3205-3921-50085ef8c1c6.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public CutterInfo()
```

```
Public Sub New
```

```
public:
CutterInfo()
```

```
new : unit -> CutterInfo
```

![](../icons/SectionExpanded.png)参见

#### 引用

[CutterInfo 类](41023477-96be-7b95-3000-fb98b3d28e50.htm)

[HslCommunication.CNC.Fanuc 命名空间](e11b249e-4a41-3205-3921-50085ef8c1c6.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## CutterInfo 属性

[原文連結](http://api.hslcommunication.cn/html/474119ab-6e0e-3b29-d1d9-f3ec883d47fd.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.CNC.Fanuc](../html/e11b249e-4a41-3205-3921-50085ef8c1c6.htm "HslCommunication.CNC.Fanuc")

[CutterInfo 类](../html/41023477-96be-7b95-3000-fb98b3d28e50.htm "CutterInfo 类")

[CutterInfo 属性](../html/474119ab-6e0e-3b29-d1d9-f3ec883d47fd.htm "CutterInfo 属性")

[LengthSharpOffset 属性](../html/0270102e-15a3-9b9c-142c-86cb0458f42b.htm "LengthSharpOffset 属性 ")

[LengthWearOffset 属性](../html/8d964911-252f-a649-ad4c-2bcc515b1bdd.htm "LengthWearOffset 属性 ")

[RadiusSharpOffset 属性](../html/4857b4aa-a44f-a04d-e87b-8a42cbdf4c61.htm "RadiusSharpOffset 属性 ")

[RadiusWearOffset 属性](../html/f89ab2a6-b3cb-cca0-bad6-df8edf560440.htm "RadiusWearOffset 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CutterInfo 属性 |

[CutterInfo](41023477-96be-7b95-3000-fb98b3d28e50.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [LengthSharpOffset](0270102e-15a3-9b9c-142c-86cb0458f42b.htm) | 长度形状补偿 |
| 公共属性 | [LengthWearOffset](8d964911-252f-a649-ad4c-2bcc515b1bdd.htm) | 长度磨损补偿 |
| 公共属性 | [RadiusSharpOffset](4857b4aa-a44f-a04d-e87b-8a42cbdf4c61.htm) | 半径形状补偿 |
| 公共属性 | [RadiusWearOffset](f89ab2a6-b3cb-cca0-bad6-df8edf560440.htm) | 半径磨损补偿 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[CutterInfo 类](41023477-96be-7b95-3000-fb98b3d28e50.htm)

[HslCommunication.CNC.Fanuc 命名空间](e11b249e-4a41-3205-3921-50085ef8c1c6.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## LengthSharpOffset 属性 

[原文連結](http://api.hslcommunication.cn/html/0270102e-15a3-9b9c-142c-86cb0458f42b.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.CNC.Fanuc](../html/e11b249e-4a41-3205-3921-50085ef8c1c6.htm "HslCommunication.CNC.Fanuc")

[CutterInfo 类](../html/41023477-96be-7b95-3000-fb98b3d28e50.htm "CutterInfo 类")

[CutterInfo 属性](../html/474119ab-6e0e-3b29-d1d9-f3ec883d47fd.htm "CutterInfo 属性")

[LengthSharpOffset 属性](../html/0270102e-15a3-9b9c-142c-86cb0458f42b.htm "LengthSharpOffset 属性 ")

[LengthWearOffset 属性](../html/8d964911-252f-a649-ad4c-2bcc515b1bdd.htm "LengthWearOffset 属性 ")

[RadiusSharpOffset 属性](../html/4857b4aa-a44f-a04d-e87b-8a42cbdf4c61.htm "RadiusSharpOffset 属性 ")

[RadiusWearOffset 属性](../html/f89ab2a6-b3cb-cca0-bad6-df8edf560440.htm "RadiusWearOffset 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CutterInfoLengthSharpOffset 属性 |

长度形状补偿

**命名空间：**
 [HslCommunication.CNC.Fanuc](e11b249e-4a41-3205-3921-50085ef8c1c6.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public double LengthSharpOffset { get; set; }
```

```
Public Property LengthSharpOffset As Double
	Get
	Set
```

```
public:
property double LengthSharpOffset {
	double get ();
	void set (double value);
}
```

```
member LengthSharpOffset : float with get, set
```

#### 属性值

类型：Double

![](../icons/SectionExpanded.png)参见

#### 引用

[CutterInfo 类](41023477-96be-7b95-3000-fb98b3d28e50.htm)

[HslCommunication.CNC.Fanuc 命名空间](e11b249e-4a41-3205-3921-50085ef8c1c6.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## LengthWearOffset 属性 

[原文連結](http://api.hslcommunication.cn/html/8d964911-252f-a649-ad4c-2bcc515b1bdd.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.CNC.Fanuc](../html/e11b249e-4a41-3205-3921-50085ef8c1c6.htm "HslCommunication.CNC.Fanuc")

[CutterInfo 类](../html/41023477-96be-7b95-3000-fb98b3d28e50.htm "CutterInfo 类")

[CutterInfo 属性](../html/474119ab-6e0e-3b29-d1d9-f3ec883d47fd.htm "CutterInfo 属性")

[LengthSharpOffset 属性](../html/0270102e-15a3-9b9c-142c-86cb0458f42b.htm "LengthSharpOffset 属性 ")

[LengthWearOffset 属性](../html/8d964911-252f-a649-ad4c-2bcc515b1bdd.htm "LengthWearOffset 属性 ")

[RadiusSharpOffset 属性](../html/4857b4aa-a44f-a04d-e87b-8a42cbdf4c61.htm "RadiusSharpOffset 属性 ")

[RadiusWearOffset 属性](../html/f89ab2a6-b3cb-cca0-bad6-df8edf560440.htm "RadiusWearOffset 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CutterInfoLengthWearOffset 属性 |

长度磨损补偿

**命名空间：**
 [HslCommunication.CNC.Fanuc](e11b249e-4a41-3205-3921-50085ef8c1c6.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public double LengthWearOffset { get; set; }
```

```
Public Property LengthWearOffset As Double
	Get
	Set
```

```
public:
property double LengthWearOffset {
	double get ();
	void set (double value);
}
```

```
member LengthWearOffset : float with get, set
```

#### 属性值

类型：Double

![](../icons/SectionExpanded.png)参见

#### 引用

[CutterInfo 类](41023477-96be-7b95-3000-fb98b3d28e50.htm)

[HslCommunication.CNC.Fanuc 命名空间](e11b249e-4a41-3205-3921-50085ef8c1c6.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## RadiusSharpOffset 属性 

[原文連結](http://api.hslcommunication.cn/html/4857b4aa-a44f-a04d-e87b-8a42cbdf4c61.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.CNC.Fanuc](../html/e11b249e-4a41-3205-3921-50085ef8c1c6.htm "HslCommunication.CNC.Fanuc")

[CutterInfo 类](../html/41023477-96be-7b95-3000-fb98b3d28e50.htm "CutterInfo 类")

[CutterInfo 属性](../html/474119ab-6e0e-3b29-d1d9-f3ec883d47fd.htm "CutterInfo 属性")

[LengthSharpOffset 属性](../html/0270102e-15a3-9b9c-142c-86cb0458f42b.htm "LengthSharpOffset 属性 ")

[LengthWearOffset 属性](../html/8d964911-252f-a649-ad4c-2bcc515b1bdd.htm "LengthWearOffset 属性 ")

[RadiusSharpOffset 属性](../html/4857b4aa-a44f-a04d-e87b-8a42cbdf4c61.htm "RadiusSharpOffset 属性 ")

[RadiusWearOffset 属性](../html/f89ab2a6-b3cb-cca0-bad6-df8edf560440.htm "RadiusWearOffset 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CutterInfoRadiusSharpOffset 属性 |

半径形状补偿

**命名空间：**
 [HslCommunication.CNC.Fanuc](e11b249e-4a41-3205-3921-50085ef8c1c6.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public double RadiusSharpOffset { get; set; }
```

```
Public Property RadiusSharpOffset As Double
	Get
	Set
```

```
public:
property double RadiusSharpOffset {
	double get ();
	void set (double value);
}
```

```
member RadiusSharpOffset : float with get, set
```

#### 属性值

类型：Double

![](../icons/SectionExpanded.png)参见

#### 引用

[CutterInfo 类](41023477-96be-7b95-3000-fb98b3d28e50.htm)

[HslCommunication.CNC.Fanuc 命名空间](e11b249e-4a41-3205-3921-50085ef8c1c6.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## RadiusWearOffset 属性 

[原文連結](http://api.hslcommunication.cn/html/f89ab2a6-b3cb-cca0-bad6-df8edf560440.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.CNC.Fanuc](../html/e11b249e-4a41-3205-3921-50085ef8c1c6.htm "HslCommunication.CNC.Fanuc")

[CutterInfo 类](../html/41023477-96be-7b95-3000-fb98b3d28e50.htm "CutterInfo 类")

[CutterInfo 属性](../html/474119ab-6e0e-3b29-d1d9-f3ec883d47fd.htm "CutterInfo 属性")

[LengthSharpOffset 属性](../html/0270102e-15a3-9b9c-142c-86cb0458f42b.htm "LengthSharpOffset 属性 ")

[LengthWearOffset 属性](../html/8d964911-252f-a649-ad4c-2bcc515b1bdd.htm "LengthWearOffset 属性 ")

[RadiusSharpOffset 属性](../html/4857b4aa-a44f-a04d-e87b-8a42cbdf4c61.htm "RadiusSharpOffset 属性 ")

[RadiusWearOffset 属性](../html/f89ab2a6-b3cb-cca0-bad6-df8edf560440.htm "RadiusWearOffset 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CutterInfoRadiusWearOffset 属性 |

半径磨损补偿

**命名空间：**
 [HslCommunication.CNC.Fanuc](e11b249e-4a41-3205-3921-50085ef8c1c6.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public double RadiusWearOffset { get; set; }
```

```
Public Property RadiusWearOffset As Double
	Get
	Set
```

```
public:
property double RadiusWearOffset {
	double get ();
	void set (double value);
}
```

```
member RadiusWearOffset : float with get, set
```

#### 属性值

类型：Double

![](../icons/SectionExpanded.png)参见

#### 引用

[CutterInfo 类](41023477-96be-7b95-3000-fb98b3d28e50.htm)

[HslCommunication.CNC.Fanuc 命名空间](e11b249e-4a41-3205-3921-50085ef8c1c6.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## CutterInfo 方法

[原文連結](http://api.hslcommunication.cn/html/dbc9c79b-50bc-fd8f-fe74-951e59d13222.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.CNC.Fanuc](../html/e11b249e-4a41-3205-3921-50085ef8c1c6.htm "HslCommunication.CNC.Fanuc")

[CutterInfo 类](../html/41023477-96be-7b95-3000-fb98b3d28e50.htm "CutterInfo 类")

[CutterInfo 方法](../html/dbc9c79b-50bc-fd8f-fe74-951e59d13222.htm "CutterInfo 方法")

[ToString 方法](../html/80661f5b-cea4-22a7-b661-21afdfb6e182.htm "ToString 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CutterInfo 方法 |

[CutterInfo](41023477-96be-7b95-3000-fb98b3d28e50.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法 | [ToString](80661f5b-cea4-22a7-b661-21afdfb6e182.htm) | (重写 ObjectToString.) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[CutterInfo 类](41023477-96be-7b95-3000-fb98b3d28e50.htm)

[HslCommunication.CNC.Fanuc 命名空间](e11b249e-4a41-3205-3921-50085ef8c1c6.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ToString 方法 

[原文連結](http://api.hslcommunication.cn/html/80661f5b-cea4-22a7-b661-21afdfb6e182.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.CNC.Fanuc](../html/e11b249e-4a41-3205-3921-50085ef8c1c6.htm "HslCommunication.CNC.Fanuc")

[CutterInfo 类](../html/41023477-96be-7b95-3000-fb98b3d28e50.htm "CutterInfo 类")

[CutterInfo 方法](../html/dbc9c79b-50bc-fd8f-fe74-951e59d13222.htm "CutterInfo 方法")

[ToString 方法](../html/80661f5b-cea4-22a7-b661-21afdfb6e182.htm "ToString 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CutterInfoToString 方法 |

[缺少 "M:HslCommunication.CNC.Fanuc.CutterInfo.ToString" 的 <summary> 文档]

**命名空间：**
 [HslCommunication.CNC.Fanuc](e11b249e-4a41-3205-3921-50085ef8c1c6.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public override string ToString()
```

```
Public Overrides Function ToString As String
```

```
public:
virtual String^ ToString() override
```

```
abstract ToString : unit -> string 
override ToString : unit -> string
```

#### 返回值

类型：String  

[缺少 "M:HslCommunication.CNC.Fanuc.CutterInfo.ToString" 的 <returns> 文档]

![](../icons/SectionExpanded.png)参见

#### 引用

[CutterInfo 类](41023477-96be-7b95-3000-fb98b3d28e50.htm)

[HslCommunication.CNC.Fanuc 命名空间](e11b249e-4a41-3205-3921-50085ef8c1c6.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## FanucOperatorMessage 类

[原文連結](http://api.hslcommunication.cn/html/dc495da4-b0c1-f6b7-c86c-2f470d1871d6.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.CNC.Fanuc](../html/e11b249e-4a41-3205-3921-50085ef8c1c6.htm "HslCommunication.CNC.Fanuc")

[FanucOperatorMessage 类](../html/dc495da4-b0c1-f6b7-c86c-2f470d1871d6.htm "FanucOperatorMessage 类")

[FanucOperatorMessage 构造函数](../html/9fc98131-29a2-c0aa-f2aa-966c55aaa7b5.htm "FanucOperatorMessage 构造函数 ")

[FanucOperatorMessage 属性](../html/bd1370d1-8417-85e5-6b3f-a9e6618cc9a2.htm "FanucOperatorMessage 属性")

[FanucOperatorMessage 方法](../html/d6567167-ba49-b19f-033b-702ac65377eb.htm "FanucOperatorMessage 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FanucOperatorMessage 类 |

Fanuc机床的操作信息

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  HslCommunication.CNC.FanucFanucOperatorMessage

**命名空间：**
 [HslCommunication.CNC.Fanuc](e11b249e-4a41-3205-3921-50085ef8c1c6.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class FanucOperatorMessage
```

```
Public Class FanucOperatorMessage
```

```
public ref class FanucOperatorMessage
```

```
type FanucOperatorMessage =  class end
```

FanucOperatorMessage 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [FanucOperatorMessage](9fc98131-29a2-c0aa-f2aa-966c55aaa7b5.htm) | 初始化 FanucOperatorMessage 类的一个新实例 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [Data](eda7cd7d-718b-054c-168a-c9608a470d6a.htm) | Operator's message strings |
| 公共属性 | [Number](f5be5033-553b-fddb-2bcd-ecc76c2ded6d.htm) | Number of operator's message |
| 公共属性 | [Type](af82fc58-95db-f877-35f6-08b0de0de4f5.htm) | Kind of operator's message |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法静态成员 | [CreateMessage](99684894-618d-fedf-45ba-f80369fd8ce7.htm) | 创建一个fanuc的操作消息对象 |
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

[HslCommunication.CNC.Fanuc 命名空间](e11b249e-4a41-3205-3921-50085ef8c1c6.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## FanucOperatorMessage 构造函数 

[原文連結](http://api.hslcommunication.cn/html/9fc98131-29a2-c0aa-f2aa-966c55aaa7b5.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.CNC.Fanuc](../html/e11b249e-4a41-3205-3921-50085ef8c1c6.htm "HslCommunication.CNC.Fanuc")

[FanucOperatorMessage 类](../html/dc495da4-b0c1-f6b7-c86c-2f470d1871d6.htm "FanucOperatorMessage 类")

[FanucOperatorMessage 构造函数](../html/9fc98131-29a2-c0aa-f2aa-966c55aaa7b5.htm "FanucOperatorMessage 构造函数 ")

[FanucOperatorMessage 属性](../html/bd1370d1-8417-85e5-6b3f-a9e6618cc9a2.htm "FanucOperatorMessage 属性")

[FanucOperatorMessage 方法](../html/d6567167-ba49-b19f-033b-702ac65377eb.htm "FanucOperatorMessage 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FanucOperatorMessage 构造函数 |

初始化 [FanucOperatorMessage](dc495da4-b0c1-f6b7-c86c-2f470d1871d6.htm) 类的一个新实例

**命名空间：**
 [HslCommunication.CNC.Fanuc](e11b249e-4a41-3205-3921-50085ef8c1c6.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public FanucOperatorMessage()
```

```
Public Sub New
```

```
public:
FanucOperatorMessage()
```

```
new : unit -> FanucOperatorMessage
```

![](../icons/SectionExpanded.png)参见

#### 引用

[FanucOperatorMessage 类](dc495da4-b0c1-f6b7-c86c-2f470d1871d6.htm)

[HslCommunication.CNC.Fanuc 命名空间](e11b249e-4a41-3205-3921-50085ef8c1c6.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## FanucOperatorMessage 属性

[原文連結](http://api.hslcommunication.cn/html/bd1370d1-8417-85e5-6b3f-a9e6618cc9a2.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.CNC.Fanuc](../html/e11b249e-4a41-3205-3921-50085ef8c1c6.htm "HslCommunication.CNC.Fanuc")

[FanucOperatorMessage 类](../html/dc495da4-b0c1-f6b7-c86c-2f470d1871d6.htm "FanucOperatorMessage 类")

[FanucOperatorMessage 属性](../html/bd1370d1-8417-85e5-6b3f-a9e6618cc9a2.htm "FanucOperatorMessage 属性")

[Data 属性](../html/eda7cd7d-718b-054c-168a-c9608a470d6a.htm "Data 属性 ")

[Number 属性](../html/f5be5033-553b-fddb-2bcd-ecc76c2ded6d.htm "Number 属性 ")

[Type 属性](../html/af82fc58-95db-f877-35f6-08b0de0de4f5.htm "Type 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FanucOperatorMessage 属性 |

[FanucOperatorMessage](dc495da4-b0c1-f6b7-c86c-2f470d1871d6.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [Data](eda7cd7d-718b-054c-168a-c9608a470d6a.htm) | Operator's message strings |
| 公共属性 | [Number](f5be5033-553b-fddb-2bcd-ecc76c2ded6d.htm) | Number of operator's message |
| 公共属性 | [Type](af82fc58-95db-f877-35f6-08b0de0de4f5.htm) | Kind of operator's message |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[FanucOperatorMessage 类](dc495da4-b0c1-f6b7-c86c-2f470d1871d6.htm)

[HslCommunication.CNC.Fanuc 命名空间](e11b249e-4a41-3205-3921-50085ef8c1c6.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Data 属性 

[原文連結](http://api.hslcommunication.cn/html/eda7cd7d-718b-054c-168a-c9608a470d6a.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.CNC.Fanuc](../html/e11b249e-4a41-3205-3921-50085ef8c1c6.htm "HslCommunication.CNC.Fanuc")

[FanucOperatorMessage 类](../html/dc495da4-b0c1-f6b7-c86c-2f470d1871d6.htm "FanucOperatorMessage 类")

[FanucOperatorMessage 属性](../html/bd1370d1-8417-85e5-6b3f-a9e6618cc9a2.htm "FanucOperatorMessage 属性")

[Data 属性](../html/eda7cd7d-718b-054c-168a-c9608a470d6a.htm "Data 属性 ")

[Number 属性](../html/f5be5033-553b-fddb-2bcd-ecc76c2ded6d.htm "Number 属性 ")

[Type 属性](../html/af82fc58-95db-f877-35f6-08b0de0de4f5.htm "Type 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FanucOperatorMessageData 属性 |

Operator's message strings

**命名空间：**
 [HslCommunication.CNC.Fanuc](e11b249e-4a41-3205-3921-50085ef8c1c6.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public string Data { get; set; }
```

```
Public Property Data As String
	Get
	Set
```

```
public:
property String^ Data {
	String^ get ();
	void set (String^ value);
}
```

```
member Data : string with get, set
```

#### 属性值

类型：String

![](../icons/SectionExpanded.png)参见

#### 引用

[FanucOperatorMessage 类](dc495da4-b0c1-f6b7-c86c-2f470d1871d6.htm)

[HslCommunication.CNC.Fanuc 命名空间](e11b249e-4a41-3205-3921-50085ef8c1c6.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Number 属性 

[原文連結](http://api.hslcommunication.cn/html/f5be5033-553b-fddb-2bcd-ecc76c2ded6d.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.CNC.Fanuc](../html/e11b249e-4a41-3205-3921-50085ef8c1c6.htm "HslCommunication.CNC.Fanuc")

[FanucOperatorMessage 类](../html/dc495da4-b0c1-f6b7-c86c-2f470d1871d6.htm "FanucOperatorMessage 类")

[FanucOperatorMessage 属性](../html/bd1370d1-8417-85e5-6b3f-a9e6618cc9a2.htm "FanucOperatorMessage 属性")

[Data 属性](../html/eda7cd7d-718b-054c-168a-c9608a470d6a.htm "Data 属性 ")

[Number 属性](../html/f5be5033-553b-fddb-2bcd-ecc76c2ded6d.htm "Number 属性 ")

[Type 属性](../html/af82fc58-95db-f877-35f6-08b0de0de4f5.htm "Type 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FanucOperatorMessageNumber 属性 |

Number of operator's message

**命名空间：**
 [HslCommunication.CNC.Fanuc](e11b249e-4a41-3205-3921-50085ef8c1c6.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public short Number { get; set; }
```

```
Public Property Number As Short
	Get
	Set
```

```
public:
property short Number {
	short get ();
	void set (short value);
}
```

```
member Number : int16 with get, set
```

#### 属性值

类型：Int16

![](../icons/SectionExpanded.png)参见

#### 引用

[FanucOperatorMessage 类](dc495da4-b0c1-f6b7-c86c-2f470d1871d6.htm)

[HslCommunication.CNC.Fanuc 命名空间](e11b249e-4a41-3205-3921-50085ef8c1c6.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Type 属性 

[原文連結](http://api.hslcommunication.cn/html/af82fc58-95db-f877-35f6-08b0de0de4f5.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.CNC.Fanuc](../html/e11b249e-4a41-3205-3921-50085ef8c1c6.htm "HslCommunication.CNC.Fanuc")

[FanucOperatorMessage 类](../html/dc495da4-b0c1-f6b7-c86c-2f470d1871d6.htm "FanucOperatorMessage 类")

[FanucOperatorMessage 属性](../html/bd1370d1-8417-85e5-6b3f-a9e6618cc9a2.htm "FanucOperatorMessage 属性")

[Data 属性](../html/eda7cd7d-718b-054c-168a-c9608a470d6a.htm "Data 属性 ")

[Number 属性](../html/f5be5033-553b-fddb-2bcd-ecc76c2ded6d.htm "Number 属性 ")

[Type 属性](../html/af82fc58-95db-f877-35f6-08b0de0de4f5.htm "Type 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FanucOperatorMessageType 属性 |

Kind of operator's message

**命名空间：**
 [HslCommunication.CNC.Fanuc](e11b249e-4a41-3205-3921-50085ef8c1c6.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public short Type { get; set; }
```

```
Public Property Type As Short
	Get
	Set
```

```
public:
property short Type {
	short get ();
	void set (short value);
}
```

```
member Type : int16 with get, set
```

#### 属性值

类型：Int16

![](../icons/SectionExpanded.png)参见

#### 引用

[FanucOperatorMessage 类](dc495da4-b0c1-f6b7-c86c-2f470d1871d6.htm)

[HslCommunication.CNC.Fanuc 命名空间](e11b249e-4a41-3205-3921-50085ef8c1c6.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## FanucOperatorMessage 方法

[原文連結](http://api.hslcommunication.cn/html/d6567167-ba49-b19f-033b-702ac65377eb.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.CNC.Fanuc](../html/e11b249e-4a41-3205-3921-50085ef8c1c6.htm "HslCommunication.CNC.Fanuc")

[FanucOperatorMessage 类](../html/dc495da4-b0c1-f6b7-c86c-2f470d1871d6.htm "FanucOperatorMessage 类")

[FanucOperatorMessage 方法](../html/d6567167-ba49-b19f-033b-702ac65377eb.htm "FanucOperatorMessage 方法")

[CreateMessage 方法](../html/99684894-618d-fedf-45ba-f80369fd8ce7.htm "CreateMessage 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FanucOperatorMessage 方法 |

[FanucOperatorMessage](dc495da4-b0c1-f6b7-c86c-2f470d1871d6.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法静态成员 | [CreateMessage](99684894-618d-fedf-45ba-f80369fd8ce7.htm) | 创建一个fanuc的操作消息对象 |
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

[FanucOperatorMessage 类](dc495da4-b0c1-f6b7-c86c-2f470d1871d6.htm)

[HslCommunication.CNC.Fanuc 命名空间](e11b249e-4a41-3205-3921-50085ef8c1c6.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## CreateMessage 方法 

[原文連結](http://api.hslcommunication.cn/html/99684894-618d-fedf-45ba-f80369fd8ce7.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.CNC.Fanuc](../html/e11b249e-4a41-3205-3921-50085ef8c1c6.htm "HslCommunication.CNC.Fanuc")

[FanucOperatorMessage 类](../html/dc495da4-b0c1-f6b7-c86c-2f470d1871d6.htm "FanucOperatorMessage 类")

[FanucOperatorMessage 方法](../html/d6567167-ba49-b19f-033b-702ac65377eb.htm "FanucOperatorMessage 方法")

[CreateMessage 方法](../html/99684894-618d-fedf-45ba-f80369fd8ce7.htm "CreateMessage 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FanucOperatorMessageCreateMessage 方法 |

创建一个fanuc的操作消息对象

**命名空间：**
 [HslCommunication.CNC.Fanuc](e11b249e-4a41-3205-3921-50085ef8c1c6.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static FanucOperatorMessage CreateMessage(
	IByteTransform byteTransform,
	byte[] buffer,
	Encoding encoding
)
```

```
Public Shared Function CreateMessage ( 
	byteTransform As IByteTransform,
	buffer As Byte(),
	encoding As Encoding
) As FanucOperatorMessage
```

```
public:
static FanucOperatorMessage^ CreateMessage(
	IByteTransform^ byteTransform, 
	array<unsigned char>^ buffer, 
	Encoding^ encoding
)
```

```
static member CreateMessage : 
        byteTransform : IByteTransform * 
        buffer : byte[] * 
        encoding : Encoding -> FanucOperatorMessage 
```

#### 参数

byteTransform
:   类型：[HslCommunication.CoreIByteTransform](56c55574-bb2a-fe66-e7d5-1332a1bc18f0.htm)  
    数据变换对象

buffer
:   类型：SystemByte  
    读取的数据缓存信息

encoding
:   类型：System.TextEncoding  
    解析的编码信息

#### 返回值

类型：[FanucOperatorMessage](dc495da4-b0c1-f6b7-c86c-2f470d1871d6.htm)  
fanuc设备的操作信息

![](../icons/SectionExpanded.png)参见

#### 引用

[FanucOperatorMessage 类](dc495da4-b0c1-f6b7-c86c-2f470d1871d6.htm)

[HslCommunication.CNC.Fanuc 命名空间](e11b249e-4a41-3205-3921-50085ef8c1c6.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## FanucSeries0i 类

[原文連結](http://api.hslcommunication.cn/html/25cac34a-b90e-b589-124d-e0d86bcf6fe3.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.CNC.Fanuc](../html/e11b249e-4a41-3205-3921-50085ef8c1c6.htm "HslCommunication.CNC.Fanuc")

[FanucSeries0i 类](../html/25cac34a-b90e-b589-124d-e0d86bcf6fe3.htm "FanucSeries0i 类")

[FanucSeries0i 构造函数](../html/46e14f26-eede-f5a7-e341-8d3a14840fad.htm "FanucSeries0i 构造函数 ")

[FanucSeries0i 属性](../html/f3c3f1dd-9cbe-2925-2bdc-f80342a8c050.htm "FanucSeries0i 属性")

[FanucSeries0i 方法](../html/a30f58ab-c363-f5d5-a62a-4faac4956b44.htm "FanucSeries0i 方法")

[FanucSeries0i 字段](../html/60abcab1-f90f-0e0b-7d99-537d4978551d.htm "FanucSeries0i 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FanucSeries0i 类 |

一个FANUC的机床通信类对象

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  [HslCommunication.Core.NetNetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)  
    [HslCommunication.Core.NetNetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)  
      HslCommunication.CNC.FanucFanucSeries0i

**命名空间：**
 [HslCommunication.CNC.Fanuc](e11b249e-4a41-3205-3921-50085ef8c1c6.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class FanucSeries0i : NetworkDoubleBase
```

```
Public Class FanucSeries0i
	Inherits NetworkDoubleBase
```

```
public ref class FanucSeries0i : public NetworkDoubleBase
```

```
type FanucSeries0i =  
    class
        inherit NetworkDoubleBase
    end
```

FanucSeries0i 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [FanucSeries0i](46e14f26-eede-f5a7-e341-8d3a14840fad.htm) | 根据IP及端口来实例化一个对象内容 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [AlienSession](5e30a9ab-7102-79a2-c3dc-450163b93107.htm) | 当前的异形连接对象，如果设置了异形连接的话，仅用于异形模式的情况使用  The current alien connection object, if alien connection is set, is only used in the case of alien mode (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共属性代码示例 | [ByteTransform](720ae40b-1a12-b12f-7845-11c3bf80de8a.htm) | 当前的数据变换机制，当你需要从字节数据转换类型数据的时候需要。  The current data transformation mechanism is required when you need to convert type data from byte data. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共属性 | [ConnectionId](902454e3-c22d-b58e-171c-c59aa43ab52b.htm) | 当前连接的唯一ID号，默认为长度20的guid码加随机数组成，方便列表管理，也可以自己指定  The unique ID number of the current connection. The default is a 20-digit guid code plus a random number. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共属性代码示例 | [ConnectTimeOut](8fc841ea-c377-ce38-79d9-0252e0ffb28c.htm) | 获取或设置连接的超时时间，单位是毫秒   Gets or sets the timeout for the connection, in milliseconds (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共属性代码示例 | [IpAddress](bea866c8-a585-e1c1-8d93-517eb2505454.htm) | 获取或是设置远程服务器的IP地址，如果是本机测试，那么需要设置为127.0.0.1   Get or set the IP address of the remote server. If it is a local test, then it needs to be set to 127.0.0.1 (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共属性 | [LocalBinding](4c347174-2cb7-c637-f92a-1b7502844749.htm) | 获取或设置绑定的本地的IP地址和端口号信息，如果端口设置为0，代表任何可用的端口  Get or set the bound local IP address and port number information, if the port is set to 0, it means any available port (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共属性代码示例 | [LogNet](46dfe698-46c6-73a2-5458-8fd568d264ce.htm) | 组件的日志工具，支持日志记录，只要实例化后，当前网络的基本信息，就以[DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm)等级进行输出  The component's logging tool supports logging. As long as the instantiation of the basic network information, the output will be output at [DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm) (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 公共属性 | [OperatePath](2d441798-c782-22eb-c0af-57fe12260042.htm) | 获取或设置当前操作的路径信息，默认为1，如果机床支持多路径的，可以设置为其他值。  Gets or sets the path information for the current operation, the default is 1, if the machine supports multipathing, it can be set to other values. |
| 公共属性代码示例 | [Port](45e6a4b4-1190-6d1a-e1e6-d8c010e15fe4.htm) | 获取或设置服务器的端口号，具体的值需要取决于对方的配置  Gets or sets the port number of the server. The specific value depends on the configuration of the other party. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共属性代码示例 | [ReceiveTimeOut](995db577-ab5c-6c9c-6f6f-253e26600472.htm) | 获取或设置接收服务器反馈的时间，如果为负数，则不接收反馈   Gets or sets the time to receive server feedback, and if it is a negative number, does not receive feedback (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共属性 | [SendBeforeHex](59a68d6e-a784-a757-40cf-138d6b05099b.htm) | 获取或设置在发送通信报文前追加发送的字节信息，HEX格式，通常用于lora组网时，需要携带 00 00 00 02 四个字节的站地址功能。  Obtain or set the byte information sent before sending communication packets, HEX format, usually used for LORA networking, you need to carry 00 00 00 02 four-byte station address function. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共属性 | [SleepTime](d275ea69-962f-70f9-b486-cb3dc3382b74.htm) | 获取或设置在正式接收对方返回数据前的时候，需要休息的时间，当设置为0的时候，不需要休息。  Get or set the time required to rest before officially receiving the data from the other party. When it is set to 0, no rest is required. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共属性 | [SocketKeepAliveTime](7f4b5ca3-e4db-76e4-ef2b-4c93bf9447ac.htm) | 获取或设置客户端的Socket的心跳时间信息，这个是Socket底层自动实现的心跳包，不基于协议层实现。默认小于0，不开启心跳检测，如果需要开启，设置 60\_000 比较合适，单位毫秒  Get or set the heartbeat time information of the Socket of the client. This is the heartbeat packet automatically implemented by the bottom layer of the Socket, not based on the protocol layer. The default value is less than 0, and heartbeat detection is not enabled. If you need to enable it, it is more appropriate to set 60\_000, in milliseconds. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共属性 | [TextEncoding](ed60b460-0bf8-67cf-427e-d4f421530779.htm) | 获取或设置当前的文本的字符编码信息，如果你不清楚，可以调用[ReadLanguage](07effb14-3a6a-142e-b4e2-94e3ca92b12b.htm)方法来自动匹配。  Get or set the character encoding information of the current text. If you are not sure, you can call the [ReadLanguage](07effb14-3a6a-142e-b4e2-94e3ca92b12b.htm) method to automatically match. |
| 公共属性代码示例 | [Token](80db00a7-aaa5-164f-e87c-af5400dc80fa.htm) | 网络类的身份令牌，在hsl协议的模式下会有效，在和设备进行通信的时候是无效的  Network-type identity tokens will be valid in the hsl protocol mode and will not be valid when communicating with the device (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的属性 | [UseServerActivePush](0e77fcc5-1744-cedf-a8e1-958031151826.htm) | 获取或设置当前的连接是否激活从服务器主动推送的功能 (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 受保护的方法 | [AccountCertificate](e9ef8320-5daf-4629-dcc6-4076e7944a9a.htm) | 认证账号，根据已经设置的用户名和密码，进行发送服务器进行账号认证。  Authentication account, according to the user name and password that have been set, sending server for account authentication. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 受保护的方法 | [AccountCertificateAsync](2dc6eb5a-79c8-1e51-7aa8-c3d29769033a.htm) | 认证账号，根据已经设置的用户名和密码，进行发送服务器进行账号认证。  Authentication account, according to the user name and password that have been set, sending server for account authentication. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 受保护的方法 | [CheckReceiveDataComplete](ea6ac6fd-de29-39d0-1098-6d90d5a1fb5d.htm) | 检查当前从网口接收的数据是否是完整的，如果是完整的，则需要返回 True，表示数据接收立即完成，默认返回 True  Check whether the data currently received from the network port is complete, and if it is complete, you need to return True, indicating that the data reception is completed immediately, and the default value is True (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 受保护的方法 | [CheckRemoteToken](1d231984-101b-0fff-0daf-1d18581a4241.htm) | 检查当前的头子节信息的令牌是否是正确的，仅用于某些特殊的协议实现  Check whether the token of the current header subsection information is correct, only for some special protocol implementations (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 公共方法 | [ClearToolGroup](fea3cfde-5d38-e3a4-7dd4-6f5f82147c95.htm) | 清除刀组号信息  Clear the knife group number information |
| 公共方法代码示例 | [ConnectClose](e235581f-1c77-3b52-6f49-a6d1c79e559a.htm) | 手动断开与远程服务器的连接，如果当前是长连接模式，那么就会切换到短连接模式  Manually disconnect from the remote server, if it is currently in long connection mode, it will switch to short connection mode (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法代码示例 | [ConnectCloseAsync](ecec49bf-ad46-e74b-69ce-211464eaea6d.htm) | 手动断开与远程服务器的连接，如果当前是长连接模式，那么就会切换到短连接模式  Manually disconnect from the remote server, if it is currently in long connection mode, it will switch to short connection mode (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法代码示例 | [ConnectServer](ea436d31-7950-42df-a9c0-3c749c47e31d.htm) | 尝试连接远程的服务器，如果连接成功，就切换短连接模式到长连接模式，后面的每次请求都共享一个通道，使得通讯速度更快速  Try to connect to a remote server. If the connection is successful, switch the short connection mode to the long connection mode. Each subsequent request will share a channel, making the communication speed faster. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法代码示例 | [ConnectServer(AlienSession)](a5e9ec5a-548c-3d66-2509-fafa5ec21c3b.htm) | 使用指定的套接字创建异形客户端，在异形客户端的模式下，网络通道需要被动创建。  Use the specified socket to create the alien client. In the alien client mode, the network channel needs to be created passively. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | [ConnectServer(MqttClient, String, String)](3ae46111-33b7-c7dc-02d8-f689c1084b5a.htm) | 使用一个MQTT中转服务器来连接设备对象，并进行相关的读取操作 (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法代码示例 | [ConnectServerAsync](7eeb4d56-e784-e48e-b49a-8498a8fe8384.htm) | 尝试连接远程的服务器，如果连接成功，就切换短连接模式到长连接模式，后面的每次请求都共享一个通道，使得通讯速度更快速  Try to connect to a remote server. If the connection is successful, switch the short connection mode to the long connection mode. Each subsequent request will share a channel, making the communication speed faster. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 受保护的方法代码示例 | [CreateSocketAndConnect(String, Int32)](a153ae27-92e3-4fe7-b2e7-c8d207bebbb0.htm) | 创建一个新的socket对象并连接到远程的地址，默认超时时间为10秒钟，需要指定ip地址以及端口号信息  Create a new socket object and connect to the remote address. The default timeout is 10 seconds. You need to specify the IP address and port number. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法代码示例 | [CreateSocketAndConnect(IPEndPoint, Int32, IPEndPoint)](7ce6288f-c6db-e968-d676-f851d2b05c77.htm) | 创建一个新的socket对象并连接到远程的地址，需要指定远程终结点，超时时间（单位是毫秒），如果需要绑定本地的IP或是端口，传入 local对象  To create a new socket object and connect to the remote address, you need to specify the remote endpoint, the timeout period (in milliseconds), if you need to bind the local IP or port, pass in the local object (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法代码示例 | [CreateSocketAndConnect(String, Int32, Int32)](d0e14a5f-1e97-e59d-078f-4fbf0412b865.htm) | 创建一个新的socket对象并连接到远程的地址，需要指定ip地址以及端口号信息，还有超时时间，单位是毫秒  To create a new socket object and connect to a remote address, you need to specify the IP address and port number information, and the timeout period in milliseconds (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法代码示例 | [CreateSocketAndConnectAsync(String, Int32)](0227bc60-1c7f-3c1d-0cd6-e49bb904d6df.htm) | 创建一个新的socket对象并连接到远程的地址，默认超时时间为10秒钟，需要指定ip地址以及端口号信息  Create a new socket object and connect to the remote address. The default timeout is 10 seconds. You need to specify the IP address and port number. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法代码示例 | [CreateSocketAndConnectAsync(IPEndPoint, Int32, IPEndPoint)](ccec07bf-cfa4-b302-210e-441821e095a4.htm) | 创建一个新的socket对象并连接到远程的地址，需要指定远程终结点，超时时间（单位是毫秒），如果需要绑定本地的IP或是端口，传入 local对象  To create a new socket object and connect to the remote address, you need to specify the remote endpoint, the timeout period (in milliseconds), if you need to bind the local IP or port, pass in the local object (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法代码示例 | [CreateSocketAndConnectAsync(String, Int32, Int32)](7ec25883-60e2-c7ac-2f90-0ea1ae0c4945.htm) | 创建一个新的socket对象并连接到远程的地址，需要指定ip地址以及端口号信息，还有超时时间，单位是毫秒  To create a new socket object and connect to a remote address, you need to specify the IP address and port number information, and the timeout period in milliseconds (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [DecideWhetherQAMessage](77091dac-82d7-8c67-d691-31ae5b2aaee4.htm) | 决定当前的消息是否是应答机制的消息内容，需要在客户端进行重写实现，如果是应答机制，返回 True, 否则返回 False  To determine whether the current message is the message content of the response mechanism, it needs to be rewritten on the client side. If it is the response mechanism, return True, otherwise return False (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | [DeleteFile](d234c829-05fa-cbc9-c6c9-00de4155986d.htm) | 根据指定的文件名来删除文件，如果是路径，则必须 '/' 结尾，如果是文件，则需要输入完整的文件名，例如：//CNC\_MEM/USER/PATH2/O12  Delete the file according to the specified file name, if it is a path, it must end with '/', if it is a file, you need to enter the complete file name, for example: //CNC\_MEM/USER/PATH2/O12 |
| 公共方法 | [DeleteFileAsync](57d466e0-d331-983c-d7dc-649fddb423e3.htm) | 根据指定的文件名来删除文件，如果是路径，则必须 '/' 结尾，如果是文件，则需要输入完整的文件名，例如：//CNC\_MEM/USER/PATH2/O12  Delete the file according to the specified file name, if it is a path, it must end with '/', if it is a file, you need to enter the complete file name, for example: //CNC\_MEM/USER/PATH2/O12 |
| 受保护的方法 | [DeleteFileByName](59afa4a4-3f3c-9831-26ee-c7c13a213bf7.htm) | 删除一个指定的文件，如果文件不存在，直接返回 True，如果文件存在则直接删除，删除成功返回 True，如果发生了异常，返回False  Delete a specified file, if the file does not exist, return True directly, if the file exists, delete it directly, if the deletion is successful, return True, if an exception occurs, return  False (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 公共方法 | [DeleteProgram](1a13c1e7-af34-60b6-443b-05abeecf8666.htm) | 根据指定的程序号信息，删除当前的程序信息  According to the designated program number information, delete the current program information |
| 公共方法 | [DeleteProgramAsync](f726c396-920a-6196-e4d1-b8e0faed2042.htm) | 根据指定的程序号信息，删除当前的程序信息  According to the designated program number information, delete the current program information |
| 公共方法 | [Dispose](8a86a200-fbe0-3b45-0fed-c551405c9b59.htm) | 释放当前的资源，如果调用了本方法，那么该对象再使用的时候，需要重新实例化。  Release the current resource. If this method is called, the object needs to be instantiated again when it is used again. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 受保护的方法 | [Dispose(Boolean)](32f03275-9ae0-525c-84bc-6e2515d11c12.htm) | 释放当前的资源，并自动关闭长连接，如果设置了的话 (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | [ExtraAfterReadFromCoreServer](9f8357e9-dff4-1002-9368-0744cdedff29.htm) | 和服务器交互完成的时候调用的方法，可以根据读写结果进行一些额外的操作，具体的操作需要根据实际的需求来重写实现  The method called when the interaction with the server is completed can perform some additional operations based on the read and write results. The specific operations need to be rewritten according to actual needs. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 受保护的方法代码示例 | [ExtraOnDisconnect](01e6fcae-9597-8fdc-0e5f-f2da2b47cb34.htm) | 根据实际的协议选择是否重写本方法，有些协议在断开连接之前，需要发送一些报文来关闭当前的网络通道  Select whether to rewrite this method according to the actual protocol. Some protocols need to send some packets to close the current network channel before disconnecting. (重写 [NetworkDoubleBaseExtraOnDisconnect(Socket)](b8569c75-5e55-cd55-5251-207465a8ba39.htm).) |
| 受保护的方法代码示例 | [ExtraOnDisconnectAsync](d7e49cac-3ffc-5913-d0e5-e4756842b539.htm) | 根据实际的协议选择是否重写本方法，有些协议在断开连接之前，需要发送一些报文来关闭当前的网络通道  Select whether to rewrite this method according to the actual protocol. Some protocols need to send some packets to close the current network channel before disconnecting. (重写 [NetworkDoubleBaseExtraOnDisconnectAsync(Socket)](86aa00d8-b2d6-ada9-5926-2b7a2237fe6b.htm).) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 受保护的方法 | [GetAvailableSocket](e86a6fc1-eb91-8de9-5ff8-799b0ea3786b.htm) | 获取本次操作的可用的网络通道，如果是短连接，就重新生成一个新的网络通道，如果是长连接，就复用当前的网络通道。  Obtain the available network channels for this operation. If it is a short connection, a new network channel is regenerated. If it is a long connection, the current network channel is reused. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 受保护的方法 | [GetAvailableSocketAsync](40cf9ce7-a5e6-5431-2e4b-56e3b7c98346.htm) | 获取本次操作的可用的网络通道，如果是短连接，就重新生成一个新的网络通道，如果是长连接，就复用当前的网络通道。  Obtain the available network channels for this operation. If it is a short connection, a new network channel is regenerated. If it is a long connection, the current network channel is reused. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 受保护的方法 | [GetNewNetMessage](08a968e5-c4fa-4645-ecf5-3f0aafffe8b1.htm) | 获取一个新的消息对象的方法，需要在继承类里面进行重写  The method to get a new message object needs to be overridden in the inheritance class (重写 [NetworkDoubleBaseGetNewNetMessage](481a40e0-8539-5fc8-3e2b-c25bcc952697.htm).) |
| 公共方法 | [GetPipeSocket](19904680-372e-e0a6-1b27-b4180e759668.htm) | 获取当前用于通信的管道信息  Get the current pipe information used for communication (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法代码示例 | [InitializationOnConnect](e152ce0c-5e19-f246-04bc-5ad871bbc639.htm) | 根据实际的协议选择是否重写本方法，有些协议在创建连接之后，需要进行一些初始化的信号握手，才能最终建立网络通道。  Whether to rewrite this method is based on the actual protocol. Some protocols require some initial signal handshake to establish a network channel after the connection is created. (重写 [NetworkDoubleBaseInitializationOnConnect(Socket)](40d5eb26-9bca-11cf-42a1-55af6ea5aebd.htm).) |
| 受保护的方法代码示例 | [InitializationOnConnectAsync](05c559d5-f14b-5ffb-cdcf-f6373945e297.htm) | 根据实际的协议选择是否重写本方法，有些协议在创建连接之后，需要进行一些初始化的信号握手，才能最终建立网络通道。  Whether to rewrite this method is based on the actual protocol. Some protocols require some initial signal handshake to establish a network channel after the connection is created. (重写 [NetworkDoubleBaseInitializationOnConnectAsync(Socket)](14cac21b-4001-da7d-09a6-941eba868e14.htm).) |
| 公共方法 | [IpAddressPing](d9611e18-4033-de95-444e-b91ae0a81709.htm) | 对当前设备的IP地址进行PING的操作，返回PING的结果，正常来说，返回Success  PING the IP address of the current device and return the PING result. Normally, it returns Success (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法 | [PackCommandWithHeader](f96881d6-e59f-9e16-07e3-1188eee1a934.htm) | 对当前的命令进行打包处理，通常是携带命令头内容，标记当前的命令的长度信息，需要进行重写，否则默认不打包  The current command is packaged, usually carrying the content of the command header, marking the length of the current command, and it needs to be rewritten, otherwise it is not packaged by default (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | [ReadAlarmStatus](3cc5fed9-583d-beb3-5824-e2babfcce366.htm) | 读取报警状态信息  Read alarm status information |
| 公共方法 | [ReadAlarmStatusAsync](68e9cf98-dcf8-81e8-c237-20955af81fa3.htm) | 读取报警状态信息  Read alarm status information |
| 公共方法 | [ReadAllDirectoryAndFile](5dc883b9-7864-46e0-0937-452e3aa8f559.htm) | 读取指定路径下的所有的子路径和文件的信息，路径信息，例如 "//CNC\_MEM/USER/" |
| 公共方法 | [ReadAllDirectoryAndFileAsync](5ac5b3d4-debe-ee45-111e-1941465711cc.htm) | 读取指定路径下的所有的子路径和文件的信息，路径信息，例如 "//CNC\_MEM/USER/" |
| 公共方法 | [ReadAllDirectoryAndFileCount](962a550b-0cbf-c2aa-46f6-582b6e431129.htm) | 获取指定的路径里所有的文件夹数量和文件数量之和，路径示例：例如 "//CNC\_MEM/USER/"， "//CNC\_MEM/USER/PATH1/" |
| 公共方法 | [ReadAllDirectoryAndFileCountAsync](4da8f996-e4d0-e237-adbd-d2d25ad39df1.htm) | 获取指定的路径里所有的文件夹数量和文件数量之和，路径示例：例如 "//CNC\_MEM/USER/"， "//CNC\_MEM/USER/PATH1/" |
| 公共方法 | [ReadAxisNames](2f8e1af8-904f-b48d-d94e-d406d5245877.htm) | 获取系统的轴名称信息，数组的长度表示有几个轴  Gets the axis name information of the system, and the length of the array indicates how many axes there are |
| 公共方法 | [ReadAxisNamesAsync](ea52754d-2f7c-2fe1-5504-5f1cc46a69f6.htm) | 获取系统的轴名称信息，数组的长度表示有几个轴  Gets the axis name information of the system, and the length of the array indicates how many axes there are |
| 公共方法 | [ReadCurrentDateTime](c45b3d5f-167b-3057-c8f9-1a0a93625f29.htm) | 读取机床的当前时间信息  Read the current time information of the machine tool |
| 公共方法 | [ReadCurrentDateTimeAsync](f218e95b-20e9-6fd5-9fd5-61cfb0e8b166.htm) | 读取机床的当前时间信息  Read the current time information of the machine tool |
| 公共方法 | [ReadCurrentForegroundDir](401c10a0-a22c-3274-ca35-c8a3e7a33d68.htm) | 读取当前程序的前台路径  Read the foreground path of the current program |
| 公共方法 | [ReadCurrentForegroundDirAsync](f84b8000-ddbf-416e-fe1b-312b645393aa.htm) | 读取当前程序的前台路径  Read the foreground path of the current program |
| 公共方法 | [ReadCurrentProduceCount](c6bb0240-438c-f32a-ea99-1a9370e5c6fa.htm) | 读取当前的已加工的零件数量  Read the current number of processed parts |
| 公共方法 | [ReadCurrentProduceCountAsync](6425b7cc-dc5e-187d-3757-978c9a6b4ca3.htm) | 读取当前的已加工的零件数量  Read the current number of processed parts |
| 公共方法 | [ReadCurrentProgram](92fff23c-abe6-c062-de76-f8e5701d8580.htm) | 读取当前的程序内容，只能读取程序的片段，返回程序内容。  Read the current program content, only read the program fragments, and return the program content. |
| 公共方法 | [ReadCurrentProgramAsync](8c48757e-1e89-36f8-3aee-351e70d24947.htm) | 读取当前的程序内容，只能读取程序的片段，返回程序内容。  Read the current program content, only read the program fragments, and return the program content. |
| 公共方法 | [ReadCutterInfos](3bd33285-d3ab-3960-23ed-78a1c8d010d9.htm) | 读取当前的刀具补偿信息  Read current tool compensation information |
| 公共方法 | [ReadCutterInfosAsync](66731090-3ae5-29a6-9e3c-facd7c9ae96c.htm) | 读取当前的刀具补偿信息  Read current tool compensation information |
| 公共方法 | [ReadCutterNumber](47060091-fd92-0786-aff1-afa49916a2a2.htm) | 读取当前的正在使用的刀具号  Read the tool number currently in use |
| 公共方法 | [ReadCutterNumberAsync](dc07d184-7b07-cd5c-18e4-f72b4275acf0.htm) | 读取当前的正在使用的刀具号  Read the tool number currently in use |
| 公共方法 | [ReadData](3da8b8e9-38f1-42f5-4345-e21d398f5b75.htm) | 读取寄存器的数据信息，需要传入寄存器的代码，起始地址，结束地址信息  To read the data information of the register, you need to pass in the code of the register, the start address, and the end address information |
| 公共方法 | [ReadDataAsync](413324b4-a6a5-22a5-a0bd-20acb833b245.htm) | 读取寄存器的数据信息，需要传入寄存器的代码，起始地址，结束地址信息  To read the data information of the register, you need to pass in the code of the register, the start address, and the end address information |
| 公共方法 | [ReadDeviceWorkPiecesSize](e5f4403c-c5ce-7002-16f0-95bb9915c184.htm) | 读取工件尺寸  Read workpiece size |
| 公共方法 | [ReadDeviceWorkPiecesSizeAsync](bb52c035-0c74-da3b-d609-1c5d33c3019b.htm) | 读取工件尺寸  Read workpiece size |
| 公共方法 | [ReadDiagnoss](043323d6-6e31-fc8b-8256-f46295d5e27e.htm) | 读取调试信息，需要指定调式编号，轴编号 |
| 公共方法 | [ReadDiagnossAsync](963bdc7e-e472-44a3-f837-8c12cc35d09e.htm) | 读取调试信息，需要指定调式编号，轴编号 |
| 公共方法 | [ReadExpectProduceCount](ba6d7588-c8e2-d475-4441-b06ed9f317a0.htm) | 读取期望的加工的零件数量  Read the expected number of processed parts |
| 公共方法 | [ReadExpectProduceCountAsync](c9b4d57e-ee6a-1749-5a2d-c580d2e8713e.htm) | 读取期望的加工的零件数量  Read the expected number of processed parts |
| 公共方法 | [ReadFanucAxisLoad](261cf950-cf24-17e0-28a1-44762de8bf50.htm) | 读取伺服负载  Read servo load |
| 公共方法 | [ReadFanucAxisLoadAsync](69ec1acb-a852-9529-a26d-f33b1ff5d37b.htm) | 读取伺服负载  Read servo load |
| 公共方法 | [ReadFeedRate](60dd5705-40f6-85f9-de2f-05876af13c82.htm) | 读取进给倍率  Read feedrate override |
| 公共方法代码示例 | [ReadFromCoreServer(Byte)](b3e14431-d773-0d1a-f453-a3ad60648c29.htm) | 将数据发送到当前的网络通道中，并从网络通道中接收一个[INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm)指定的完整的报文，网络通道将根据[GetAvailableSocket](e86a6fc1-eb91-8de9-5ff8-799b0ea3786b.htm)方法自动获取，本方法是线程安全的。  Send data to the current network channel and receive a complete message specified by [INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm) from the network channel. The network channel will be automatically obtained according to the [GetAvailableSocket](e86a6fc1-eb91-8de9-5ff8-799b0ea3786b.htm) method This method is thread-safe. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | [ReadFromCoreServer(IEnumerableByte)](580dae22-fad5-5816-d58f-43068e89e419.htm) | 将多个数据报文按顺序发到设备，并从设备接收返回的数据内容，然后拼接成一个Byte[]信息，需要重写[UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm)方法才能返回正确的结果。  Send multiple data packets to the device in sequence, and receive the returned data content from the device, and then splicing them into a Byte[] message, you need to rewrite [UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm) method to return the correct result. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServer(Byte, Boolean, Boolean)](1a33d974-805d-f149-3665-eb8d8a99ea84.htm) | 将数据发送到当前的网络通道中，并从网络通道中接收一个[INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm)指定的完整的报文，网络通道将根据[GetAvailableSocket](e86a6fc1-eb91-8de9-5ff8-799b0ea3786b.htm)方法自动获取，本方法是线程安全的。  Send data to the current network channel and receive a complete message specified by [INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm) from the network channel. The network channel will be automatically obtained according to the [GetAvailableSocket](e86a6fc1-eb91-8de9-5ff8-799b0ea3786b.htm) method This method is thread-safe. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServer(Socket, Byte, Boolean, Boolean)](3b72b289-b201-7548-0981-ac049c9a4135.htm) | 将数据报文发送指定的网络通道上，根据当前指定的[INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm)类型，返回一条完整的数据指令  Sends a data message to the specified network channel, and returns a complete data command according to the currently specified [INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm) type (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServerAsync(Byte)](ec51a18e-0f5d-dbe1-c229-d62e8760ac97.htm) | 将数据发送到当前的网络通道中，并从网络通道中接收一个[INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm)指定的完整的报文，网络通道将根据[GetAvailableSocket](e86a6fc1-eb91-8de9-5ff8-799b0ea3786b.htm)方法自动获取，本方法是线程安全的。  Send data to the current network channel and receive a complete message specified by [INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm) from the network channel. The network channel will be automatically obtained according to the [GetAvailableSocket](e86a6fc1-eb91-8de9-5ff8-799b0ea3786b.htm) method This method is thread-safe. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | [ReadFromCoreServerAsync(IEnumerableByte)](0fce97fd-8441-6b54-0877-b8dfacc738e3.htm) | 将多个数据报文按顺序发到设备，并从设备接收返回的数据内容，然后拼接成一个Byte[]信息，需要重写[UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm)方法才能返回正确的结果。  Send multiple data packets to the device in sequence, and receive the returned data content from the device, and then splicing them into a Byte[] message, you need to rewrite [UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm) method to return the correct result. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServerAsync(Byte, Boolean, Boolean)](7525b379-04c2-8dac-b451-a8b5413d6747.htm) | 将数据发送到当前的网络通道中，并从网络通道中接收一个[INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm)指定的完整的报文，网络通道将根据[GetAvailableSocket](e86a6fc1-eb91-8de9-5ff8-799b0ea3786b.htm)方法自动获取，本方法是线程安全的。  Send data to the current network channel and receive a complete message specified by [INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm) from the network channel. The network channel will be automatically obtained according to the [GetAvailableSocket](e86a6fc1-eb91-8de9-5ff8-799b0ea3786b.htm) method This method is thread-safe. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServerAsync(Socket, Byte, Boolean, Boolean)](07db4a00-b548-8bd6-5de6-12deb91f5311.htm) | 将数据报文发送指定的网络通道上，根据当前指定的[INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm)类型，返回一条完整的数据指令  Sends a data message to the specified network channel, and returns a complete data command according to the currently specified [INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm) type (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | [ReadLanguage](07effb14-3a6a-142e-b4e2-94e3ca92b12b.htm) | 读取机床的语言设定信息，具体值的含义参照API文档说明  Read the language setting information of the machine tool, refer to the API documentation for the meaning of the specific values |
| 公共方法 | [ReadLanguageAsync](8b03a867-9b52-e29b-8f5a-fed8add0a5f9.htm) | 读取机床的语言设定信息，具体值的含义参照API文档说明  Read the language setting information of the machine tool, refer to the API documentation for the meaning of the specific values |
| 公共方法 | [ReadOperatorMessage](cf6ac234-d624-0e0b-f005-1ec20c27a316.htm) | 读取机床的操作信息  Read machine operation information |
| 公共方法 | [ReadPMCData](5a9501be-1ebf-6473-81f6-f28cb792fbe7.htm) | 读取PMC数据，需要传入起始地址和结束地址，返回byte[]数据信息  To read PMC data, you need to pass in the start address and length, and return byte[] data information |
| 公共方法 | [ReadPMCDataAsync](2bd1304f-48a5-fc88-32d7-a03e4c9dc962.htm) | 读取PMC数据，需要传入起始地址和结束地址，返回byte[]数据信息  To read PMC data, you need to pass in the start address and length, and return byte[] data information |
| 公共方法 | [ReadProgram(Int32, String)](b338bd69-419a-4e25-affa-26a8b2f6ed8b.htm) | **[商业授权]** 读取指定程序号的程序内容，可以指定路径信息，路径默认为空就是主路径，//CNC\_MEM/USER/PATH1/ ，也可以指定其他路径 **[Authorization]** Read the program content of the specified program number |
| 公共方法 | [ReadProgram(String, String)](0f42ec4b-ac66-44ea-0af4-e148a7b9b480.htm) | **[商业授权]** 读取指定程序号的程序内容，可以指定路径信息，路径默认为空就是主路径，//CNC\_MEM/USER/PATH1/ ，也可以指定其他路径 **[Authorization]** Read the program content of the specified program number |
| 公共方法 | [ReadProgramAsync(Int32, String)](7b591dad-8161-ba9e-0af2-c099c9c37faa.htm) | **[商业授权]** 读取指定程序号的程序内容，可以指定路径信息，路径默认为空就是主路径，//CNC\_MEM/USER/PATH1/ ，也可以指定其他路径 **[Authorization]** Read the program content of the specified program number |
| 公共方法 | [ReadProgramAsync(String, String)](f8211b10-191e-5a5c-9fd2-cf36d8a92f01.htm) | **[商业授权]** 读取指定程序号的程序内容，可以指定路径信息，路径默认为空就是主路径，//CNC\_MEM/USER/PATH1/ ，也可以指定其他路径 **[Authorization]** Read the program content of the specified program number |
| 公共方法 | [ReadProgramList](393d3c09-5084-4e1b-1238-57d0f8eaf94e.htm) | 读取设备的程序列表  Read the program list of the device |
| 公共方法 | [ReadProgramListAsync](adffc5e5-d0aa-b0b9-f595-09d7669896b0.htm) | 读取设备的程序列表  Read the program list of the device |
| 公共方法 | [ReadProgramNumber](45ee871c-0c6a-74c8-cc8b-0f3b33e951d1.htm) | 读取程序号信息  Read program number |
| 公共方法 | [ReadProgramNumberAsync](f0a4dc17-e4fe-1be8-2ab7-8cd677219563.htm) | 读取程序号信息  Read program number |
| 公共方法 | [ReadSpindleLoad](9e18d250-7a2d-4239-8ab4-d5db45d82394.htm) | 读取主轴负载  Read spindle load |
| 公共方法 | [ReadSpindleLoadAsync](140d0998-4895-017e-6d73-f128426bd11b.htm) | 读取主轴负载  Read spindle load |
| 公共方法 | [ReadSpindleNames](23ca2375-7b7e-78d0-ffe7-6ca53ddba0c7.htm) | 读取系统的主轴名称信息，返回的数组长度表示有几个主轴  Reads the system's spindle name information, and the returned array length indicates how many spindles there are |
| 公共方法 | [ReadSpindleNamesAsync](eb27bbe0-8b84-fa84-fa78-8e0bdbe26c6b.htm) | 读取系统的主轴名称信息，返回的数组长度表示有几个主轴  Reads the system's spindle name information, and the returned array length indicates how many spindles there are |
| 公共方法 | [ReadSpindleRate](db0420d4-2ff3-34c7-b9ce-188e4375f13f.htm) | 读取主轴倍率  Read spindle override |
| 公共方法 | [ReadSpindleSpeedAndFeedRate](6f2accd4-1100-64e8-7829-d1dc77a88d31.htm) | 主轴转速及进给倍率  Spindle speed and feedrate override |
| 公共方法 | [ReadSpindleSpeedAndFeedRateAsync](816f6e8f-db8f-e9f0-43e2-ed042011912d.htm) | 主轴转速及进给倍率  Spindle speed and feedrate override |
| 公共方法 | [ReadSysAllCoors](787f447e-1c00-438a-4356-bc1ae5d3b872.htm) | 读取机床的坐标，包括机械坐标，绝对坐标，相对坐标  Read the coordinates of the machine tool, including mechanical coordinates, absolute coordinates, and relative coordinates |
| 公共方法 | [ReadSysAllCoorsAsync](0b23143e-adb4-10aa-c364-67de74a1abf9.htm) | 读取机床的坐标，包括机械坐标，绝对坐标，相对坐标  Read the coordinates of the machine tool, including mechanical coordinates, absolute coordinates, and relative coordinates |
| 公共方法 | [ReadSysInfo](462b6e58-7101-c235-1cb9-e586e96124a9.htm) | 获取fanuc机床设备的基本信息，型号，轴数量等等。  Get basic information about fanuc machines, models, number of axes and much more |
| 公共方法 | [ReadSysStatusInfo](b39a7199-4fb3-5ce8-7d4a-ef4b264ef4c8.htm) | 读取系统的基本信息状态，工作模式，运行状态，是否急停等等操作  Read the basic information status of the system, working mode, running status, emergency stop, etc. |
| 公共方法 | [ReadSysStatusInfoAsync](de515b11-3f34-26bf-1523-a7da36dd76c1.htm) | 读取系统的基本信息状态，工作模式，运行状态，是否急停等等操作  Read the basic information status of the system, working mode, running status, emergency stop, etc. |
| 公共方法 | [ReadSystemAlarm](6f8dd08b-e79e-0b01-a8d9-a124f501ff77.htm) | 读取报警信息  Read alarm information |
| 公共方法 | [ReadSystemAlarmAsync](32772719-d3f3-88d2-4f0b-614a30e9ac20.htm) | 读取报警信息  Read alarm information |
| 公共方法 | [ReadSystemMacroValue(Int32)](17b3bae1-83ca-2a9e-f920-f84b300f5436.htm) | 读取宏变量，可以用来读取刀具号  Read macro variable, can be used to read tool number |
| 公共方法 | [ReadSystemMacroValue(Int32, Int32)](162e7904-f5fb-03f6-826c-adcaec7ae6ac.htm) | 读取宏变量，可以用来读取刀具号  Read macro variable, can be used to read tool number |
| 公共方法 | [ReadSystemMacroValueAsync(Int32)](dd30c938-ab49-7ad2-90d5-3fe7d56daf42.htm) | 读取宏变量，可以用来读取刀具号  Read macro variable, can be used to read tool number |
| 公共方法 | [ReadSystemMacroValueAsync(Int32, Int32)](7a590f54-ecb7-716f-5d23-64538f8a5fb6.htm) | 读取宏变量，可以用来读取刀具号  Read macro variable, can be used to read tool number |
| 公共方法 | [ReadSystemProgramCurrent](32c0dc2c-8617-d1a6-ecc0-be4084347f2b.htm) | 读取程序名及程序号  Read program name and program number |
| 公共方法 | [ReadSystemProgramCurrentAsync](31b59dfe-2fc4-11e1-8876-6c4bb4d122d7.htm) | 读取程序名及程序号  Read program name and program number |
| 公共方法 | [ReadTimeData](73fa77dd-eca1-3ccb-9435-f2eeae9cde26.htm) | 读取fanuc机床的时间，0是开机时间，1是运行时间，2是切割时间，3是循环时间，4是空闲时间，返回秒为单位的信息  Read the time of the fanuc machine tool, 0 is the boot time, 1 is the running time, 2 is the cutting time, 3 is the cycle time, 4 is the idle time, and returns the information in seconds. |
| 公共方法 | [ReadTimeDataAsync](f5e0bf55-c64c-2ab7-0ba9-fa1bc0eb00b5.htm) | 读取fanuc机床的时间，0是开机时间，1是运行时间，2是切割时间，3是循环时间，4是空闲时间，返回秒为单位的信息  Read the time of the fanuc machine tool, 0 is the boot time, 1 is the running time, 2 is the cutting time, 3 is the cycle time, 4 is the idle time, and returns the information in seconds. |
| 公共方法 | [ReadToolInfoByGroup](b4e9a0bc-906f-b7d3-68c7-c7fa0a37c924.htm) | 根据刀组号读取刀具信息，包括寿命及使用次数。  Tool information is read according to the tool group number, including life and number of uses. |
| 公共方法 | [ReadUseToolGroupId](c9993369-e3f3-b6e0-a5de-737ca0630531.htm) | 读取当前正在使用的刀组号  Reads the knife group number that is currently in use |
| 受保护的方法 | [Receive(SslStream, Int32, Int32, ActionInt64, Int64)](50ad65ee-3ff1-6a08-31f3-a29090797796.htm) | 接收固定长度的字节数组，允许指定超时时间，默认为60秒，当length大于0时，接收固定长度的数据内容，当length小于0时，接收不大于2048长度的随机数据信息  Receiving a fixed-length byte array, allowing a specified timeout time. The default is 60 seconds. When length is greater than 0, fixed-length data content is received. When length is less than 0, random data information of a length not greater than 2048 is received. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [Receive(Socket, Int32, Int32, ActionInt64, Int64)](28c887ec-7a68-e90a-f531-b2298ded707b.htm) | 接收固定长度的字节数组，允许指定超时时间，默认为60秒，当length大于0时，接收固定长度的数据内容，当length小于0时，接收不大于2048长度的随机数据信息  Receiving a fixed-length byte array, allowing a specified timeout time. The default is 60 seconds. When length is greater than 0, fixed-length data content is received. When length is less than 0, random data information of a length not greater than 2048 is received. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [Receive(SslStream, Byte, Int32, Int32, Int32, ActionInt64, Int64)](d6fbd69f-3aa1-9f84-139a-04003a9ce7c0.htm) | 接收固定长度的字节数组，允许指定超时时间，默认为60秒，当length大于0时，接收固定长度的数据内容，当length小于0时，buffer长度的缓存数据  Receiving a fixed-length byte array, allowing a specified timeout time. The default is 60 seconds. When length is greater than 0, fixed-length data content is received. When length is less than 0, random data information of a length not greater than 2048 is received. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [Receive(Socket, Byte, Int32, Int32, Int32, ActionInt64, Int64)](403209a0-350c-ede4-a17c-8ab8ad5ddc5b.htm) | 接收固定长度的字节数组，允许指定超时时间，默认为60秒，当length大于0时，接收固定长度的数据内容，当length小于0时，buffer长度的缓存数据  Receiving a fixed-length byte array, allowing a specified timeout time. The default is 60 seconds. When length is greater than 0, fixed-length data content is received. When length is less than 0, random data information of a length not greater than 2048 is received. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [ReceiveAndCheckBytes](1d627617-63a7-7078-86e8-cfb10c3ad500.htm) | [自校验] 接收一条完整的同步数据，包含头子节和内容字节，基础的数据，如果结果异常，则结束通讯  [Self-checking] Receive a complete synchronization data, including header subsection and content bytes, basic data, if the result is abnormal, the communication ends (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [ReceiveAndCheckBytesAsync](a9e13d42-804a-afee-4c86-d59ad03ec469.htm) | [自校验] 接收一条完整的同步数据，包含头子节和内容字节，基础的数据，如果结果异常，则结束通讯  [Self-checking] Receive a complete synchronization data, including header subsection and content bytes, basic data, if the result is abnormal, the communication ends (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [ReceiveAsync(SslStream, Int32, Int32, ActionInt64, Int64)](548ecb57-8a13-a71d-0a29-a71c282f94d8.htm) | 接收固定长度的字节数组，允许指定超时时间，默认为60秒，当length大于0时，接收固定长度的数据内容，当length小于0时，接收不大于2048长度的随机数据信息  Receiving a fixed-length byte array, allowing a specified timeout time. The default is 60 seconds. When length is greater than 0, fixed-length data content is received. When length is less than 0, random data information of a length not greater than 2048 is received. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [ReceiveAsync(Socket, Int32, Int32, ActionInt64, Int64)](6edaea53-6855-bfcd-33e6-5ca3df268636.htm) | 接收固定长度的字节数组，允许指定超时时间，默认为60秒，当length大于0时，接收固定长度的数据内容，当length小于0时，接收不大于2048长度的随机数据信息  Receiving a fixed-length byte array, allowing a specified timeout time. The default is 60 seconds. When length is greater than 0, fixed-length data content is received. When length is less than 0, random data information of a length not greater than 2048 is received. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [ReceiveAsync(SslStream, Byte, Int32, Int32, Int32, ActionInt64, Int64)](6104ebb8-3044-1bc6-3972-add79588e90c.htm) | 接收固定长度的字节数组，允许指定超时时间，默认为60秒，当length大于0时，接收固定长度的数据内容，当length小于0时，buffer长度的缓存数据  Receiving a fixed-length byte array, allowing a specified timeout time. The default is 60 seconds. When length is greater than 0, fixed-length data content is received. When length is less than 0, random data information of a length not greater than 2048 is received. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [ReceiveAsync(Socket, Byte, Int32, Int32, Int32, ActionInt64, Int64)](5d183a4a-19ba-46b2-9338-a4a88f2a7c70.htm) | 接收固定长度的字节数组，允许指定超时时间，默认为60秒，当length大于0时，接收固定长度的数据内容，当length小于0时，buffer长度的缓存数据  Receiving a fixed-length byte array, allowing a specified timeout time. The default is 60 seconds. When length is greater than 0, fixed-length data content is received. When length is less than 0, random data information of a length not greater than 2048 is received. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [ReceiveByMessage](33d5e9f5-67f0-c2c9-4ad3-c59ebfdde4d7.htm) | 接收一条完整的 数据内容，需要指定超时时间，单位为毫秒。   Receive a complete data content, Need to specify a timeout period in milliseconds (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [ReceiveByMessageAsync](e75c8b31-a6df-7f85-cd9f-d5e6a095a983.htm) | 接收一条完整的 数据内容，需要指定超时时间，单位为毫秒。   Receive a complete data content, Need to specify a timeout period in milliseconds (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [ReceiveBytesContentFromSocket](190b4c27-554b-713d-0f49-7942516e96f0.htm) | [自校验] 从网络中接收一串字节数据，如果结果异常，则结束通讯  [Self-checking] Receive a string of byte data from the network. If the result is abnormal, the communication ends. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [ReceiveBytesContentFromSocketAsync](2e3e6ee6-4d7e-a5c9-f642-13969c4c3b78.htm) | [自校验] 从网络中接收一串字节数据，如果结果异常，则结束通讯  [Self-checking] Receive a string of byte data from the network. If the result is abnormal, the communication ends. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [ReceiveCommandLineFromSocket(Socket, Byte, Int32)](f6fbd8fd-5f1f-f4c2-2ab4-3b7f1166dd4c.htm) | 接收一行命令数据，需要自己指定这个结束符，默认超时时间为60秒，也即是60000，单位是毫秒  To receive a line of command data, you need to specify the terminator yourself. The default timeout is 60 seconds, which is 60,000, in milliseconds. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [ReceiveCommandLineFromSocket(Socket, Byte, Byte, Int32)](9fb845a9-a46c-2285-2076-33bcaf62c16a.htm) | 接收一行命令数据，需要自己指定这个结束符，默认超时时间为60秒，也即是60000，单位是毫秒  To receive a line of command data, you need to specify the terminator yourself. The default timeout is 60 seconds, which is 60,000, in milliseconds. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [ReceiveCommandLineFromSocketAsync(Socket, Byte, Int32)](6f21495a-3b08-501e-3219-30f00b831a36.htm) | 接收一行命令数据，需要自己指定这个结束符，默认超时时间为60秒，也即是60000，单位是毫秒  To receive a line of command data, you need to specify the terminator yourself. The default timeout is 60 seconds, which is 60,000, in milliseconds. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [ReceiveCommandLineFromSocketAsync(Socket, Byte, Byte, Int32)](bcc92f0b-de8e-c78f-4aba-fff2cdd34c82.htm) | 接收一行命令数据，需要自己指定这个结束符，默认超时时间为60秒，也即是60000，单位是毫秒  To receive a line of command data, you need to specify the terminator yourself. The default timeout is 60 seconds, which is 60,000, in milliseconds. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [ReceiveHslMessage](9cd1bc29-c99a-5426-1fb1-047268154fe3.htm) | 接收一条hsl协议的数据信息，自动解析，解压，解码操作，获取最后的实际的数据，接收结果依次为暗号，用户码，负载数据  Receive a piece of hsl protocol data information, automatically parse, decompress, and decode operations to obtain the last actual data. The result is a opCode, user code, and payload data in order. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [ReceiveHslMessageAsync](63e3bd3f-efd6-fc4b-9361-5caab03c47d9.htm) | 接收一条hsl协议的数据信息，自动解析，解压，解码操作，获取最后的实际的数据，接收结果依次为暗号，用户码，负载数据  Receive a piece of hsl protocol data information, automatically parse, decompress, and decode operations to obtain the last actual data. The result is a opCode, user code, and payload data in order. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [ReceiveMqttFile](a529cfa6-f57f-bede-b266-4f8942c4934a.htm) | 使用MQTT协议从网络接收字节数组，然后写入文件或流中，支持进度报告  Use MQTT protocol to receive byte array from the network, and then write it to file or stream, support progress report (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [ReceiveMqttFileAsync](2b725eb6-1051-c1c8-5e31-9787204d9ffa.htm) | 使用MQTT协议从网络接收字节数组，然后写入文件或流中，支持进度报告  Use MQTT protocol to receive byte array from the network, and then write it to file or stream, support progress report (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [ReceiveMqttMessage(SslStream, Int32, ActionInt64, Int64)](c268dc7d-d3d1-3806-8e2c-8de3a7736cc2.htm) | 接收一条完整的MQTT协议的报文信息，包含控制码和负载数据  Receive a message of a completed MQTT protocol, including control code and payload data (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [ReceiveMqttMessage(Socket, Int32, ActionInt64, Int64)](81c27af3-012e-665f-129c-7aae05c77021.htm) | 接收一条完整的MQTT协议的报文信息，包含控制码和负载数据  Receive a message of a completed MQTT protocol, including control code and payload data (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [ReceiveMqttMessageAsync(SslStream, Int32, ActionInt64, Int64)](5fa070b0-e679-7201-df00-bf7c34811100.htm) | 接收一条完整的MQTT协议的报文信息，包含控制码和负载数据  Receive a message of a completed MQTT protocol, including control code and payload data (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [ReceiveMqttMessageAsync(Socket, Int32, ActionInt64, Int64)](f6ec22af-2c46-7ace-716c-a198e9a70a82.htm) | 接收一条完整的MQTT协议的报文信息，包含控制码和负载数据  Receive a message of a completed MQTT protocol, including control code and payload data (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [ReceiveMqttStream](dcd504a7-4231-7f98-2abc-dd6aed3ed1e2.htm) | 使用MQTT协议从socket接收指定长度的字节数组，然后全部写入到流中，可以指定进度报告  Use the MQTT protocol to receive a byte array of specified length from the socket, and then write all of them to the stream, and you can specify a progress report (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [ReceiveMqttStreamAsync](8e562525-0bdf-6592-82ea-c27cb6aa339b.htm) | 使用MQTT协议从socket接收指定长度的字节数组，然后全部写入到流中，可以指定进度报告  Use the MQTT protocol to receive a byte array of specified length from the socket, and then write all of them to the stream, and you can specify a progress report (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [ReceiveRedisCommand](896bd785-c06b-c7f6-475f-cc78ca50b201.htm) | 从网络接收一条完整的redis报文的消息  Receive a complete redis message from the network (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [ReceiveRedisCommandAsync](7da184d3-e92d-fb42-9c44-a3c00f990157.htm) | 从网络接收一条完整的redis报文的消息  Receive a complete redis message from the network (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [ReceiveRedisCommandString](3299524f-73b1-476e-18f3-cddd811a5b2e.htm) | 接收一行基于redis协议的字符串的信息，需要指定固定的长度  Receive a line of information based on the redis protocol string, you need to specify a fixed length (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [ReceiveRedisCommandStringAsync](9bdf24f8-847b-5177-d45d-2524e58774d9.htm) | 接收一行基于redis协议的字符串的信息，需要指定固定的长度  Receive a line of information based on the redis protocol string, you need to specify a fixed length (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [ReceiveStringArrayContentFromSocket](d32e83aa-c62b-0c81-ebe6-1482fe8a5bfa.htm) | [自校验] 从网络中接收一个字符串数组，如果结果异常，则结束通讯  [Self-check] Receive an array of strings from the network. If the result is abnormal, the communication ends. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [ReceiveStringArrayContentFromSocketAsync](ca8a25f7-200d-b95c-8b31-d1374c77c280.htm) | [自校验] 从网络中接收一个字符串数组，如果结果异常，则结束通讯  [Self-check] Receive an array of strings from the network. If the result is abnormal, the communication ends. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [ReceiveStringContentFromSocket](ad3b2915-8cfb-ff45-f24f-1eb255bd289e.htm) | [自校验] 从网络中接收一个字符串数据，如果结果异常，则结束通讯  [Self-checking] Receive a string of data from the network. If the result is abnormal, the communication ends. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [ReceiveStringContentFromSocketAsync](3cd827b7-1fca-6184-0067-85c016d3eb15.htm) | [自校验] 从网络中接收一个字符串数据，如果结果异常，则结束通讯  [Self-checking] Receive a string of data from the network. If the result is abnormal, the communication ends. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [Send(SslStream, Byte)](8fad1450-4dc9-4505-2873-09fd128d0f09.htm) | 发送消息给套接字，直到完成的时候返回，经过测试，本方法是线程安全的。  Send a message to the socket until it returns when completed. After testing, this method is thread-safe. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [Send(Socket, Byte)](d98c31b7-b055-011c-0546-93434bd77af5.htm) | 发送消息给套接字，直到完成的时候返回，经过测试，本方法是线程安全的。  Send a message to the socket until it returns when completed. After testing, this method is thread-safe. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [Send(SslStream, Byte, Int32, Int32)](b701a27a-0c44-1cf4-7bdf-73518ea1f2e6.htm) | 发送消息给套接字，直到完成的时候返回，经过测试，本方法是线程安全的。  Send a message to the socket until it returns when completed. After testing, this method is thread-safe. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [Send(Socket, Byte, Int32, Int32)](09e7c105-bfea-0186-03fe-981d8124c208.htm) | 发送消息给套接字，直到完成的时候返回，经过测试，本方法是线程安全的。  Send a message to the socket until it returns when completed. After testing, this method is thread-safe. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [SendAccountAndCheckReceive](86806e1d-1b59-4333-9c64-65fc445a64e1.htm) | [自校验] 直接发送字符串数组并确认对方接收完成数据，如果结果异常，则结束通讯  [Self-check] Send string array directly and confirm that the other party has received the completed data. If the result is abnormal, the communication ends. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [SendAccountAndCheckReceiveAsync](cb6aa97e-a3f7-e64a-e68a-636756d28207.htm) | [自校验] 直接发送字符串数组并确认对方接收完成数据，如果结果异常，则结束通讯  [Self-check] Send string array directly and confirm that the other party has received the completed data. If the result is abnormal, the communication ends. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [SendAsync(SslStream, Byte)](9ab7bc43-5a28-5e39-ec77-811d091507a5.htm) | 发送消息给套接字，直到完成的时候返回，经过测试，本方法是线程安全的。  Send a message to the socket until it returns when completed. After testing, this method is thread-safe. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [SendAsync(Socket, Byte)](dec1ad4d-89b9-fc85-1426-8e61acb444aa.htm) | 发送消息给套接字，直到完成的时候返回，经过测试，本方法是线程安全的。  Send a message to the socket until it returns when completed. After testing, this method is thread-safe. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [SendAsync(SslStream, Byte, Int32, Int32)](4ce38aad-0f1a-da9b-aa13-a755fa236ac7.htm) | 发送消息给套接字，直到完成的时候返回，经过测试，本方法是线程安全的。  Send a message to the socket until it returns when completed. After testing, this method is thread-safe. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [SendAsync(Socket, Byte, Int32, Int32)](e18a57a9-9caf-193f-9971-fbd41bba9e61.htm) | 发送消息给套接字，直到完成的时候返回，经过测试，本方法是线程安全的。  Send a message to the socket until it returns when completed. After testing, this method is thread-safe. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [SendBaseAndCheckReceive](3ab1c2a3-2a33-b6f8-b648-98ca72b4f139.htm) | [自校验] 发送字节数据并确认对方接收完成数据，如果结果异常，则结束通讯  [Self-check] Send the byte data and confirm that the other party has received the completed data. If the result is abnormal, the communication ends. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [SendBaseAndCheckReceiveAsync](b6344eac-286c-c80f-ed85-56c5784065a5.htm) | [自校验] 发送字节数据并确认对方接收完成数据，如果结果异常，则结束通讯  [Self-check] Send the byte data and confirm that the other party has received the completed data. If the result is abnormal, the communication ends. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [SendBytesAndCheckReceive](ddbf44f9-7dff-964d-3e8c-d4198887183f.htm) | [自校验] 发送字节数据并确认对方接收完成数据，如果结果异常，则结束通讯  [Self-check] Send the byte data and confirm that the other party has received the completed data. If the result is abnormal, the communication ends. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [SendBytesAndCheckReceiveAsync](038563fa-a4e0-2103-4f93-f52f1ee2ded3.htm) | [自校验] 发送字节数据并确认对方接收完成数据，如果结果异常，则结束通讯  [Self-check] Send the byte data and confirm that the other party has received the completed data. If the result is abnormal, the communication ends. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [SendMqttFile(Socket, Stream, String, String, ActionInt64, Int64, AesCryptography, HslCancelToken)](46ad6f88-4010-ff7d-12d7-196516140b31.htm) | 使用MQTT协议将一个数据流发送到网络上去，需要保存的文件名，可选指定文件描述信息，进度报告  Use the MQTT protocol to send a data stream to the network, the file name that needs to be saved, optional file description information, progress report (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [SendMqttFile(Socket, String, String, String, ActionInt64, Int64, AesCryptography, HslCancelToken)](29af3d39-9ebe-f5a6-b203-69243b1ec55e.htm) | 使用MQTT协议将一个文件发送到网络上去，需要指定文件名，保存的文件名，可选指定文件描述信息，进度报告  To send a file to the network using the MQTT protocol, you need to specify the file name, the saved file name, optionally specify the file description information, and the progress report (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [SendMqttFileAsync(Socket, Stream, String, String, ActionInt64, Int64, AesCryptography, HslCancelToken)](38401fac-1a74-bbfb-2847-27fa314ef1e7.htm) | 使用MQTT协议将一个数据流发送到网络上去，需要保存的文件名，可选指定文件描述信息，进度报告  Use the MQTT protocol to send a data stream to the network, the file name that needs to be saved, optional file description information, progress report (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [SendMqttFileAsync(Socket, String, String, String, ActionInt64, Int64, AesCryptography, HslCancelToken)](92b852f5-9536-9588-4919-6faf4e6b835a.htm) | 使用MQTT协议将一个文件发送到网络上去，需要指定文件名，保存的文件名，可选指定文件描述信息，进度报告  To send a file to the network using the MQTT protocol, you need to specify the file name, the saved file name, optionally specify the file description information, and the progress report (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [SendMqttStream](12db1776-8fb3-7d64-4331-2df07128c59a.htm) | 使用MQTT协议将流中的数据读取到字节数组，然后都写入到socket里面，可以指定进度报告，主要用于将文件发送到网络。  Use the MQTT protocol to read the data in the stream into a byte array, and then write them all into the socket. You can specify a progress report, which is mainly used to send files to the network. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [SendMqttStreamAsync](8493fd7e-4ba3-cc9c-efbf-18747cae084d.htm) | 使用MQTT协议将流中的数据读取到字节数组，然后都写入到socket里面，可以指定进度报告，主要用于将文件发送到网络。  Use the MQTT protocol to read the data in the stream into a byte array, and then write them all into the socket. You can specify a progress report, which is mainly used to send files to the network. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [SendStreamToSocket](f6220302-c53a-4f8a-ab69-7827a9779519.htm) | 发送一个流的所有数据到指定的网络套接字，需要指定发送的数据长度，支持按照百分比的进度报告  Send all the data of a stream to the specified network socket. You need to specify the length of the data to be sent. It supports the progress report in percentage. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [SendStreamToSocketAsync](e4e82bab-973b-99b8-a3da-098c009d244f.htm) | 发送一个流的所有数据到指定的网络套接字，需要指定发送的数据长度，支持按照百分比的进度报告  Send all the data of a stream to the specified network socket. You need to specify the length of the data to be sent. It supports the progress report in percentage. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [SendStringAndCheckReceive(Socket, Int32, String)](4aa4b2d4-98fd-33ce-9a08-4a4f9529f950.htm) | [自校验] 直接发送字符串数据并确认对方接收完成数据，如果结果异常，则结束通讯  [Self-checking] Send string data directly and confirm that the other party has received the completed data. If the result is abnormal, the communication ends. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [SendStringAndCheckReceive(Socket, Int32, String)](b7bdf13e-1a1d-cbe5-f511-e6a5e0be4c69.htm) | [自校验] 直接发送字符串数组并确认对方接收完成数据，如果结果异常，则结束通讯  [Self-check] Send string array directly and confirm that the other party has received the completed data. If the result is abnormal, the communication ends. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [SendStringAndCheckReceiveAsync(Socket, Int32, String)](41c7589d-8116-66a0-5311-685001239af6.htm) | [自校验] 直接发送字符串数据并确认对方接收完成数据，如果结果异常，则结束通讯  [Self-checking] Send string data directly and confirm that the other party has received the completed data. If the result is abnormal, the communication ends. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [SendStringAndCheckReceiveAsync(Socket, Int32, String)](814826a1-27c4-ca81-424c-ae50df6eb0cb.htm) | [自校验] 直接发送字符串数组并确认对方接收完成数据，如果结果异常，则结束通讯  [Self-check] Send string array directly and confirm that the other party has received the completed data. If the result is abnormal, the communication ends. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 公共方法 | [SetCurrentProgram](0482b5d2-3417-6a32-d10b-ce411c7f3495.htm) | 设置指定的程序号为当前的主程序，如果程序号不存在，返回错误信息  Set the specified program number as the current main program, if the program number does not exist, an error message will be returned |
| 公共方法 | [SetCurrentProgramAsync](12bbe257-ec1c-9d76-332d-d19d7f2eaf74.htm) | 设置指定的程序号为当前的主程序，如果程序号不存在，返回错误信息  Set the specified program number as the current main program, if the program number does not exist, an error message will be returned |
| 公共方法 | [SetDeviceProgsCurr](60eed3dc-8c04-1dd7-d8b6-4f8f6cd783dd.htm) | 设置指定路径为当前路径  Set the specified path as the current path |
| 公共方法 | [SetDeviceProgsCurrAsync](fc4f2d6a-d174-0109-d45d-d2a8205a97f7.htm) | 设置指定路径为当前路径  Set the specified path as the current path |
| 公共方法 | [SetLoginAccount](725afc25-6e3f-bd7d-25ed-580270c8fc56.htm) | 设置当前的登录的账户名和密码信息，并启用账户验证的功能，账户名为空时设置不生效  Set the current login account name and password information, and enable the account verification function. The account name setting will not take effect when it is empty (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法代码示例 | [SetPersistentConnection](d8e8df9f-cbe9-6863-54e9-aa3b38906b26.htm) | 在读取数据之前可以调用本方法将客户端设置为长连接模式，相当于跳过了ConnectServer的结果验证，对异形客户端无效，当第一次进行通信时再进行创建连接请求。  Before reading the data, you can call this method to set the client to the long connection mode, which is equivalent to skipping the result verification of ConnectServer, and it is invalid for the alien client. When the first communication is performed, the connection creation request is performed. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | [SetPipeSocket](b6a7a243-2a6d-c4ea-45b7-bf980450a2bd.htm) | 设置一个新的网络管道，一般来说不需要调用本方法，当多个网口设备共用一个网络连接时才需要使用本方法进行设置共享的管道。  To set up a new network channel, generally speaking, you do not need to call this method. This method is only needed to set up a shared channel when multiple network port devices share a network connection. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | [StartProcessing](1b710a44-7f0d-2840-dd32-9e4383edf25a.htm) | 启动加工程序  Start the processing program |
| 公共方法 | [StartProcessingAsync](150ff8b5-aff6-3a77-02aa-848409110795.htm) | 启动加工程序  Start the processing program |
| 公共方法 | [ToString](75a85f6d-b145-6c8b-7703-fe875009ce0c.htm) | (重写 [NetworkDoubleBaseToString](85a8532f-8fc5-886c-9d76-70844e0a7d64.htm).) |
| 公共方法 | [UnpackResponseContent](32bad9a0-b59d-c171-ce85-74055edfc597.htm) | 根据对方返回的报文命令，对命令进行基本的拆包，例如各种Modbus协议拆包为统一的核心报文，还支持对报文的验证  According to the message command returned by the other party, the command is basically unpacked, for example, various Modbus protocols are unpacked into a unified core message, and the verification of the message is also supported (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | [WriteCutterLengthShapeOffset](d59ac9fb-e8d8-b3c3-8b1f-9fd00b10bfbb.htm) | 根据刀具号写入长度形状补偿，刀具号为1-24  Write length shape compensation according to the tool number, the tool number is 1-24 |
| 公共方法 | [WriteCutterLengthSharpOffsetAsync](4e9bdbdd-2a35-67fa-e446-047d1224b110.htm) | 根据刀具号写入长度形状补偿，刀具号为1-24  Write length shape compensation according to the tool number, the tool number is 1-24 |
| 公共方法 | [WriteCutterLengthWearOffset](04b5bcf5-d966-5a52-7acf-7bf950192057.htm) | 根据刀具号写入长度磨损补偿，刀具号为1-24  Write length wear compensation according to the tool number, the tool number is 1-24 |
| 公共方法 | [WriteCutterLengthWearOffsetAsync](36c9c202-2810-8981-45d0-5ec1c565ca66.htm) | 根据刀具号写入长度磨损补偿，刀具号为1-24  Write length wear compensation according to the tool number, the tool number is 1-24 |
| 公共方法 | [WriteCutterRadiusShapeOffset](3d098a46-32a0-16a1-a684-da5f3f7b314c.htm) | 根据刀具号写入半径形状补偿，刀具号为1-24  Write radius shape compensation according to the tool number, the tool number is 1-24 |
| 公共方法 | [WriteCutterRadiusSharpOffsetAsync](6e4d4858-2b49-2edd-226f-6067152fd55d.htm) | 根据刀具号写入半径形状补偿，刀具号为1-24  Write radius shape compensation according to the tool number, the tool number is 1-24 |
| 公共方法 | [WriteCutterRadiusWearOffset](46e900e8-a978-cc47-847a-d2495c68ca26.htm) | 根据刀具号写入半径磨损补偿，刀具号为1-24  Write radius wear compensation according to the tool number, the tool number is 1-24 |
| 公共方法 | [WriteCutterRadiusWearOffsetAsync](871fbf74-42f7-dd68-abb2-b293de093760.htm) | 根据刀具号写入半径磨损补偿，刀具号为1-24  Write radius wear compensation according to the tool number, the tool number is 1-24 |
| 公共方法 | [WriteData](83c49fc3-e111-27fa-1507-071305810814.htm) | 将原始字节的数据写入到指定的寄存器里，需要传入寄存器的代码，起始地址，原始的字节数据信息  To write the original byte data into the specified register, you need to pass in the code of the register, the starting address, and the original byte data information |
| 公共方法 | [WriteDataAsync](a797695f-d279-582c-e220-c3aac0eafe49.htm) | 将原始字节的数据写入到指定的寄存器里，需要传入寄存器的代码，起始地址，原始的字节数据信息  To write the original byte data into the specified register, you need to pass in the code of the register, the starting address, and the original byte data information |
| 公共方法 | [WritePMCData](2ef60dde-887e-55d1-cc95-ad2c0578693e.htm) | 写入PMC数据，需要传入起始地址和，以及等待写入的byte[]数据信息  To write PMC data, you need to pass in the start address, as well as the byte[] data information waiting to be written |
| 公共方法 | [WritePMCDataAsync](92d8b19a-061b-155b-87c5-178f9ab282b8.htm) | 写入PMC数据，需要传入起始地址和，以及等待写入的byte[]数据信息  To write PMC data, you need to pass in the start address, as well as the byte[] data information waiting to be written |
| 公共方法 | [WriteProgramContent](fa630137-8341-f0e4-a576-b9607f70193e.htm) | **[商业授权]** 将指定程序内容的NC加工程序，写入到数控机床里，返回是否下载成功 **[Authorization]** Download the NC machining program to the CNC machine tool, and return whether the download is successful |
| 公共方法 | [WriteProgramContentAsync](e5660b56-e296-acc3-cff4-c2c8cb3fece9.htm) | **[商业授权]** 将指定程序内容的NC加工程序，写入到数控机床里，返回是否下载成功 **[Authorization]** Download the NC machining program to the CNC machine tool, and return whether the download is successful |
| 公共方法 | [WriteProgramFile](fd6f83e8-fd6b-9316-4deb-1161905a4074.htm) | **[商业授权]** 将指定文件的NC加工程序，下载到数控机床里，返回是否下载成功 **[Authorization]** Download the NC machining program of the specified file to the CNC machine tool, and return whether the download is successful |
| 公共方法 | [WriteProgramFileAsync](443219d4-d894-3d2a-f88a-918ffbdfce4e.htm) | **[商业授权]** 将指定文件的NC加工程序，下载到数控机床里，返回是否下载成功 **[Authorization]** Download the NC machining program of the specified file to the CNC machine tool, and return whether the download is successful |
| 受保护的方法 | [WriteStreamFromSocket](97bf95ea-51fe-8f2e-5da6-1c83766aaa78.htm) | 从套接字中接收所有的数据然后写入到指定的流当中去，需要指定数据的长度，支持按照百分比进行进度报告  Receives all data from the socket and writes it to the specified stream. The length of the data needs to be specified, and progress reporting is supported in percentage. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [WriteStreamFromSocketAsync](b4cf7f6d-1e0c-a266-4c49-683fdb14a922.htm) | 从套接字中接收所有的数据然后写入到指定的流当中去，需要指定数据的长度，支持按照百分比进行进度报告  Receives all data from the socket and writes it to the specified stream. The length of the data needs to be specified, and progress reporting is supported in percentage. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 公共方法 | [WriteSystemMacroValue](6c34115e-c976-a644-6370-33489aeac5f6.htm) | 写宏变量，需要指定地址及写入的数据  Write macro variable, need to specify the address and write data |
| 公共方法 | [WriteSystemMacroValueAsync](e5431223-fdae-ebf7-8313-f375ec6f9203.htm) | 写宏变量，需要指定地址及写入的数据  Write macro variable, need to specify the address and write data |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)字段

|  | 名称 | 说明 |
| --- | --- | --- |
| 受保护的字段 | [fileCacheSize](2328fb89-cfda-7fa9-88d0-98f1ae3c7c54.htm) | 文件传输的时候的缓存大小，直接影响传输的速度，值越大，传输速度越快，越占内存，默认为100K大小  The size of the cache during file transfer directly affects the speed of the transfer. The larger the value, the faster the transfer speed and the more memory it takes. The default size is 100K. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的字段 | [isPersistentConn](17c281db-b635-a496-ef0e-8e91770c94be.htm) | 是否是长连接的状态  Whether it is a long connection state (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 受保护的字段 | [isUseAccountCertificate](c3f9fdea-29af-2f88-8f6a-f4651318a8a6.htm) | 是否使用账号登录，这个账户登录的功能是HSL组件创建的服务器特有的功能。  Whether to log in using an account. The function of this account login is a server-specific function created by the  HSL  component. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 受保护的字段 | [LogMsgFormatBinary](a00be7ba-6c1e-fd46-c68c-e004433341e5.htm) | 设置日志记录报文是否二进制，如果为False，那就使用ASCII码  Set whether the log message is binary, if it is False, then use ASCII code (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 受保护的字段 | [pipeSocket](d4efd9b0-ad04-ce88-5dde-843794e4e3fa.htm) | 当前的网络的管道信息 (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.CNC.Fanuc 命名空间](e11b249e-4a41-3205-3921-50085ef8c1c6.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## FanucSeries0i 构造函数 

[原文連結](http://api.hslcommunication.cn/html/46e14f26-eede-f5a7-e341-8d3a14840fad.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.CNC.Fanuc](../html/e11b249e-4a41-3205-3921-50085ef8c1c6.htm "HslCommunication.CNC.Fanuc")

[FanucSeries0i 类](../html/25cac34a-b90e-b589-124d-e0d86bcf6fe3.htm "FanucSeries0i 类")

[FanucSeries0i 构造函数](../html/46e14f26-eede-f5a7-e341-8d3a14840fad.htm "FanucSeries0i 构造函数 ")

[FanucSeries0i 属性](../html/f3c3f1dd-9cbe-2925-2bdc-f80342a8c050.htm "FanucSeries0i 属性")

[FanucSeries0i 方法](../html/a30f58ab-c363-f5d5-a62a-4faac4956b44.htm "FanucSeries0i 方法")

[FanucSeries0i 字段](../html/60abcab1-f90f-0e0b-7d99-537d4978551d.htm "FanucSeries0i 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FanucSeries0i 构造函数 |

根据IP及端口来实例化一个对象内容

**命名空间：**
 [HslCommunication.CNC.Fanuc](e11b249e-4a41-3205-3921-50085ef8c1c6.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public FanucSeries0i(
	string ipAddress,
	int port = 8193
)
```

```
Public Sub New ( 
	ipAddress As String,
	Optional port As Integer = 8193
)
```

```
public:
FanucSeries0i(
	String^ ipAddress, 
	int port = 8193
)
```

```
new : 
        ipAddress : string * 
        ?port : int 
(* Defaults:
        let _port = defaultArg port 8193
*)
-> FanucSeries0i
```

#### 参数

ipAddress
:   类型：SystemString  
    Ip地址信息

port (Optional)
:   类型：SystemInt32  
    端口号

![](../icons/SectionExpanded.png)参见

#### 引用

[FanucSeries0i 类](25cac34a-b90e-b589-124d-e0d86bcf6fe3.htm)

[HslCommunication.CNC.Fanuc 命名空间](e11b249e-4a41-3205-3921-50085ef8c1c6.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## FanucSeries0i 属性

[原文連結](http://api.hslcommunication.cn/html/f3c3f1dd-9cbe-2925-2bdc-f80342a8c050.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.CNC.Fanuc](../html/e11b249e-4a41-3205-3921-50085ef8c1c6.htm "HslCommunication.CNC.Fanuc")

[FanucSeries0i 类](../html/25cac34a-b90e-b589-124d-e0d86bcf6fe3.htm "FanucSeries0i 类")

[FanucSeries0i 属性](../html/f3c3f1dd-9cbe-2925-2bdc-f80342a8c050.htm "FanucSeries0i 属性")

[OperatePath 属性](../html/2d441798-c782-22eb-c0af-57fe12260042.htm "OperatePath 属性 ")

[TextEncoding 属性](../html/ed60b460-0bf8-67cf-427e-d4f421530779.htm "TextEncoding 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| FanucSeries0i 属性 |

[FanucSeries0i](25cac34a-b90e-b589-124d-e0d86bcf6fe3.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [AlienSession](5e30a9ab-7102-79a2-c3dc-450163b93107.htm) | 当前的异形连接对象，如果设置了异形连接的话，仅用于异形模式的情况使用  The current alien connection object, if alien connection is set, is only used in the case of alien mode (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共属性代码示例 | [ByteTransform](720ae40b-1a12-b12f-7845-11c3bf80de8a.htm) | 当前的数据变换机制，当你需要从字节数据转换类型数据的时候需要。  The current data transformation mechanism is required when you need to convert type data from byte data. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共属性 | [ConnectionId](902454e3-c22d-b58e-171c-c59aa43ab52b.htm) | 当前连接的唯一ID号，默认为长度20的guid码加随机数组成，方便列表管理，也可以自己指定  The unique ID number of the current connection. The default is a 20-digit guid code plus a random number. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共属性代码示例 | [ConnectTimeOut](8fc841ea-c377-ce38-79d9-0252e0ffb28c.htm) | 获取或设置连接的超时时间，单位是毫秒   Gets or sets the timeout for the connection, in milliseconds (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共属性代码示例 | [IpAddress](bea866c8-a585-e1c1-8d93-517eb2505454.htm) | 获取或是设置远程服务器的IP地址，如果是本机测试，那么需要设置为127.0.0.1   Get or set the IP address of the remote server. If it is a local test, then it needs to be set to 127.0.0.1 (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共属性 | [LocalBinding](4c347174-2cb7-c637-f92a-1b7502844749.htm) | 获取或设置绑定的本地的IP地址和端口号信息，如果端口设置为0，代表任何可用的端口  Get or set the bound local IP address and port number information, if the port is set to 0, it means any available port (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共属性代码示例 | [LogNet](46dfe698-46c6-73a2-5458-8fd568d264ce.htm) | 组件的日志工具，支持日志记录，只要实例化后，当前网络的基本信息，就以[DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm)等级进行输出  The component's logging tool supports logging. As long as the instantiation of the basic network information, the output will be output at [DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm) (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 公共属性 | [OperatePath](2d441798-c782-22eb-c0af-57fe12260042.htm) | 获取或设置当前操作的路径信息，默认为1，如果机床支持多路径的，可以设置为其他值。  Gets or sets the path information for the current operation, the default is 1, if the machine supports multipathing, it can be set to other values. |
| 公共属性代码示例 | [Port](45e6a4b4-1190-6d1a-e1e6-d8c010e15fe4.htm) | 获取或设置服务器的端口号，具体的值需要取决于对方的配置  Gets or sets the port number of the server. The specific value depends on the configuration of the other party. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共属性代码示例 | [ReceiveTimeOut](995db577-ab5c-6c9c-6f6f-253e26600472.htm) | 获取或设置接收服务器反馈的时间，如果为负数，则不接收反馈   Gets or sets the time to receive server feedback, and if it is a negative number, does not receive feedback (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共属性 | [SendBeforeHex](59a68d6e-a784-a757-40cf-138d6b05099b.htm) | 获取或设置在发送通信报文前追加发送的字节信息，HEX格式，通常用于lora组网时，需要携带 00 00 00 02 四个字节的站地址功能。  Obtain or set the byte information sent before sending communication packets, HEX format, usually used for LORA networking, you need to carry 00 00 00 02 four-byte station address function. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共属性 | [SleepTime](d275ea69-962f-70f9-b486-cb3dc3382b74.htm) | 获取或设置在正式接收对方返回数据前的时候，需要休息的时间，当设置为0的时候，不需要休息。  Get or set the time required to rest before officially receiving the data from the other party. When it is set to 0, no rest is required. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共属性 | [SocketKeepAliveTime](7f4b5ca3-e4db-76e4-ef2b-4c93bf9447ac.htm) | 获取或设置客户端的Socket的心跳时间信息，这个是Socket底层自动实现的心跳包，不基于协议层实现。默认小于0，不开启心跳检测，如果需要开启，设置 60\_000 比较合适，单位毫秒  Get or set the heartbeat time information of the Socket of the client. This is the heartbeat packet automatically implemented by the bottom layer of the Socket, not based on the protocol layer. The default value is less than 0, and heartbeat detection is not enabled. If you need to enable it, it is more appropriate to set 60\_000, in milliseconds. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共属性 | [TextEncoding](ed60b460-0bf8-67cf-427e-d4f421530779.htm) | 获取或设置当前的文本的字符编码信息，如果你不清楚，可以调用[ReadLanguage](07effb14-3a6a-142e-b4e2-94e3ca92b12b.htm)方法来自动匹配。  Get or set the character encoding information of the current text. If you are not sure, you can call the [ReadLanguage](07effb14-3a6a-142e-b4e2-94e3ca92b12b.htm) method to automatically match. |
| 公共属性代码示例 | [Token](80db00a7-aaa5-164f-e87c-af5400dc80fa.htm) | 网络类的身份令牌，在hsl协议的模式下会有效，在和设备进行通信的时候是无效的  Network-type identity tokens will be valid in the hsl protocol mode and will not be valid when communicating with the device (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的属性 | [UseServerActivePush](0e77fcc5-1744-cedf-a8e1-958031151826.htm) | 获取或设置当前的连接是否激活从服务器主动推送的功能 (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[FanucSeries0i 类](25cac34a-b90e-b589-124d-e0d86bcf6fe3.htm)

[HslCommunication.CNC.Fanuc 命名空间](e11b249e-4a41-3205-3921-50085ef8c1c6.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)