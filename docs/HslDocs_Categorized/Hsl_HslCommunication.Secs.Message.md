# HslCommunication - HslCommunication.Secs.Message

> 分類頁數: 8



---
## HslCommunication.Secs.Message

[原文連結](http://api.hslcommunication.cn/html/e51b7aa4-2a26-de1d-4009-eae157262c4c.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Secs.Message](../html/e51b7aa4-2a26-de1d-4009-eae157262c4c.htm "HslCommunication.Secs.Message")

[SecsHsmsMessage 类](../html/37532ceb-c98f-741c-69ca-cf154a5ad9df.htm "SecsHsmsMessage 类")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslCommunication.Secs.Message 命名空间 |

[缺少 "N:HslCommunication.Secs.Message" 的 <summary> 文档]

![](../icons/SectionExpanded.png)类

|  | 类 | 说明 |
| --- | --- | --- |
| 公共类 | [SecsHsmsMessage](37532ceb-c98f-741c-69ca-cf154a5ad9df.htm) | Hsms协议的消息定义 |

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## SecsHsmsMessage 类

[原文連結](http://api.hslcommunication.cn/html/37532ceb-c98f-741c-69ca-cf154a5ad9df.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Secs.Message](../html/e51b7aa4-2a26-de1d-4009-eae157262c4c.htm "HslCommunication.Secs.Message")

[SecsHsmsMessage 类](../html/37532ceb-c98f-741c-69ca-cf154a5ad9df.htm "SecsHsmsMessage 类")

[SecsHsmsMessage 构造函数](../html/c6707d1f-0ca5-0586-e3a9-8082780acc92.htm "SecsHsmsMessage 构造函数 ")

[SecsHsmsMessage 属性](../html/4250bbfb-44b7-d8a4-9531-a88c16fbc4c7.htm "SecsHsmsMessage 属性")

[SecsHsmsMessage 方法](../html/0370daef-e5fa-2b9a-b44f-be4d8e4f91a7.htm "SecsHsmsMessage 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| SecsHsmsMessage 类 |

Hsms协议的消息定义

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  [HslCommunication.Core.IMessageNetMessageBase](55e9c192-5d9b-9d92-31be-b05cd76d6903.htm)  
    HslCommunication.Secs.MessageSecsHsmsMessage

**命名空间：**
 [HslCommunication.Secs.Message](e51b7aa4-2a26-de1d-4009-eae157262c4c.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class SecsHsmsMessage : NetMessageBase, 
	INetMessage
```

```
Public Class SecsHsmsMessage
	Inherits NetMessageBase
	Implements INetMessage
```

```
public ref class SecsHsmsMessage : public NetMessageBase, 
	INetMessage
```

```
type SecsHsmsMessage =  
    class
        inherit NetMessageBase
        interface INetMessage
    end
```

SecsHsmsMessage 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [SecsHsmsMessage](c6707d1f-0ca5-0586-e3a9-8082780acc92.htm) | 初始化 SecsHsmsMessage 类的一个新实例 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [ContentBytes](27f34824-6adc-b126-b622-6f9eeaabe44f.htm) | 消息内容字节  Message content byte (继承自 [NetMessageBase](55e9c192-5d9b-9d92-31be-b05cd76d6903.htm)。) |
| 公共属性 | [HeadBytes](ee605527-7bfc-97a3-52b8-620d839d3d4c.htm) | 消息头字节  Message header byte (继承自 [NetMessageBase](55e9c192-5d9b-9d92-31be-b05cd76d6903.htm)。) |
| 公共属性 | [ProtocolHeadBytesLength](e06b4a1f-8f9b-55b8-8cd9-7bff446c8636.htm) | 消息头的指令长度，第一次接受数据的长度  Instruction length of the message header, the length of the first received data |
| 公共属性 | [SendBytes](41476212-a254-3722-3445-9e636263d0f6.htm) | 发送的字节信息  Byte information sent (继承自 [NetMessageBase](55e9c192-5d9b-9d92-31be-b05cd76d6903.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [CheckHeadBytesLegal](8c83a112-609f-e2e8-8f8e-8bbe42acf991.htm) | 检查头子节的合法性  Check the legitimacy of the head subsection (重写 [NetMessageBaseCheckHeadBytesLegal(Byte)](3085e36e-54a0-a2ff-cc81-099c11d5b238.htm).) |
| 公共方法 | [CheckMessageMatch](857a414e-7059-b002-b34c-10255390d6b7.htm) | 检查发送的接收的报文是否是匹配的，如果匹配，则返回 1, 如果不匹配且直接返回错误，则返回 0，如果不匹配继续接收，直到匹配或是超时，则返回 -1  If the packet is matched, 1 is returned. If the packet is not matched and an error is returned, 0 is returned. If the packet is not matched, -1 is returned until the packet is matched or times out (继承自 [NetMessageBase](55e9c192-5d9b-9d92-31be-b05cd76d6903.htm)。) |
| 公共方法 | [CheckReceiveDataComplete](909effc7-04ee-4941-593e-2da3f57e6201.htm) | 当消息头报文的长度定义为-1的时候，则使用动态的长度信息，可以使用本方法来判断一个消息是否处于完整的状态。  If the length of the message header is defined as -1, this method can be used to determine whether a message is in the complete state by using dynamic length information. (继承自 [NetMessageBase](55e9c192-5d9b-9d92-31be-b05cd76d6903.htm)。) |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | [GetContentLengthByHeadBytes](bb3f82a8-d286-dab8-03da-b829d4417640.htm) | 从当前的头子节文件中提取出接下来需要接收的数据长度  Extract the length of the data to be received from the current header file |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法 | [GetHeadBytesIdentity](61d7c255-4473-31e7-7f84-dfd957a8826f.htm) | 获取头子节里的消息标识  Get the message ID in the header subsection (继承自 [NetMessageBase](55e9c192-5d9b-9d92-31be-b05cd76d6903.htm)。) |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法 | [PependedUselesByteLength](91167f2b-bebd-0102-29d4-b514cd3aa0a0.htm) | 在接收头报文的时候，返回前置无效的报文头字节长度，默认为0，不处理  When receiving a header message, return the header byte length of the invalid header, the default is 0, and no processing is performed. (继承自 [NetMessageBase](55e9c192-5d9b-9d92-31be-b05cd76d6903.htm)。) |
| 公共方法 | [ToString](762f299e-5af6-4f6a-292c-2ae472bfe3b5.htm) | (继承自 [NetMessageBase](55e9c192-5d9b-9d92-31be-b05cd76d6903.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.Secs.Message 命名空间](e51b7aa4-2a26-de1d-4009-eae157262c4c.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## SecsHsmsMessage 构造函数 

[原文連結](http://api.hslcommunication.cn/html/c6707d1f-0ca5-0586-e3a9-8082780acc92.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Secs.Message](../html/e51b7aa4-2a26-de1d-4009-eae157262c4c.htm "HslCommunication.Secs.Message")

[SecsHsmsMessage 类](../html/37532ceb-c98f-741c-69ca-cf154a5ad9df.htm "SecsHsmsMessage 类")

[SecsHsmsMessage 构造函数](../html/c6707d1f-0ca5-0586-e3a9-8082780acc92.htm "SecsHsmsMessage 构造函数 ")

[SecsHsmsMessage 属性](../html/4250bbfb-44b7-d8a4-9531-a88c16fbc4c7.htm "SecsHsmsMessage 属性")

[SecsHsmsMessage 方法](../html/0370daef-e5fa-2b9a-b44f-be4d8e4f91a7.htm "SecsHsmsMessage 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| SecsHsmsMessage 构造函数 |

初始化 [SecsHsmsMessage](37532ceb-c98f-741c-69ca-cf154a5ad9df.htm) 类的一个新实例

**命名空间：**
 [HslCommunication.Secs.Message](e51b7aa4-2a26-de1d-4009-eae157262c4c.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public SecsHsmsMessage()
```

```
Public Sub New
```

```
public:
SecsHsmsMessage()
```

```
new : unit -> SecsHsmsMessage
```

![](../icons/SectionExpanded.png)参见

#### 引用

[SecsHsmsMessage 类](37532ceb-c98f-741c-69ca-cf154a5ad9df.htm)

[HslCommunication.Secs.Message 命名空间](e51b7aa4-2a26-de1d-4009-eae157262c4c.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## SecsHsmsMessage 属性

[原文連結](http://api.hslcommunication.cn/html/4250bbfb-44b7-d8a4-9531-a88c16fbc4c7.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Secs.Message](../html/e51b7aa4-2a26-de1d-4009-eae157262c4c.htm "HslCommunication.Secs.Message")

[SecsHsmsMessage 类](../html/37532ceb-c98f-741c-69ca-cf154a5ad9df.htm "SecsHsmsMessage 类")

[SecsHsmsMessage 属性](../html/4250bbfb-44b7-d8a4-9531-a88c16fbc4c7.htm "SecsHsmsMessage 属性")

[ProtocolHeadBytesLength 属性](../html/e06b4a1f-8f9b-55b8-8cd9-7bff446c8636.htm "ProtocolHeadBytesLength 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| SecsHsmsMessage 属性 |

[SecsHsmsMessage](37532ceb-c98f-741c-69ca-cf154a5ad9df.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [ContentBytes](27f34824-6adc-b126-b622-6f9eeaabe44f.htm) | 消息内容字节  Message content byte (继承自 [NetMessageBase](55e9c192-5d9b-9d92-31be-b05cd76d6903.htm)。) |
| 公共属性 | [HeadBytes](ee605527-7bfc-97a3-52b8-620d839d3d4c.htm) | 消息头字节  Message header byte (继承自 [NetMessageBase](55e9c192-5d9b-9d92-31be-b05cd76d6903.htm)。) |
| 公共属性 | [ProtocolHeadBytesLength](e06b4a1f-8f9b-55b8-8cd9-7bff446c8636.htm) | 消息头的指令长度，第一次接受数据的长度  Instruction length of the message header, the length of the first received data |
| 公共属性 | [SendBytes](41476212-a254-3722-3445-9e636263d0f6.htm) | 发送的字节信息  Byte information sent (继承自 [NetMessageBase](55e9c192-5d9b-9d92-31be-b05cd76d6903.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[SecsHsmsMessage 类](37532ceb-c98f-741c-69ca-cf154a5ad9df.htm)

[HslCommunication.Secs.Message 命名空间](e51b7aa4-2a26-de1d-4009-eae157262c4c.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ProtocolHeadBytesLength 属性 

[原文連結](http://api.hslcommunication.cn/html/e06b4a1f-8f9b-55b8-8cd9-7bff446c8636.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Secs.Message](../html/e51b7aa4-2a26-de1d-4009-eae157262c4c.htm "HslCommunication.Secs.Message")

[SecsHsmsMessage 类](../html/37532ceb-c98f-741c-69ca-cf154a5ad9df.htm "SecsHsmsMessage 类")

[SecsHsmsMessage 属性](../html/4250bbfb-44b7-d8a4-9531-a88c16fbc4c7.htm "SecsHsmsMessage 属性")

[ProtocolHeadBytesLength 属性](../html/e06b4a1f-8f9b-55b8-8cd9-7bff446c8636.htm "ProtocolHeadBytesLength 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| SecsHsmsMessageProtocolHeadBytesLength 属性 |

消息头的指令长度，第一次接受数据的长度  
Instruction length of the message header, the length of the first received data

**命名空间：**
 [HslCommunication.Secs.Message](e51b7aa4-2a26-de1d-4009-eae157262c4c.htm)  
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

[SecsHsmsMessage 类](37532ceb-c98f-741c-69ca-cf154a5ad9df.htm)

[HslCommunication.Secs.Message 命名空间](e51b7aa4-2a26-de1d-4009-eae157262c4c.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## SecsHsmsMessage 方法

[原文連結](http://api.hslcommunication.cn/html/0370daef-e5fa-2b9a-b44f-be4d8e4f91a7.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Secs.Message](../html/e51b7aa4-2a26-de1d-4009-eae157262c4c.htm "HslCommunication.Secs.Message")

[SecsHsmsMessage 类](../html/37532ceb-c98f-741c-69ca-cf154a5ad9df.htm "SecsHsmsMessage 类")

[SecsHsmsMessage 方法](../html/0370daef-e5fa-2b9a-b44f-be4d8e4f91a7.htm "SecsHsmsMessage 方法")

[CheckHeadBytesLegal 方法](../html/8c83a112-609f-e2e8-8f8e-8bbe42acf991.htm "CheckHeadBytesLegal 方法 ")

[GetContentLengthByHeadBytes 方法](../html/bb3f82a8-d286-dab8-03da-b829d4417640.htm "GetContentLengthByHeadBytes 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| SecsHsmsMessage 方法 |

[SecsHsmsMessage](37532ceb-c98f-741c-69ca-cf154a5ad9df.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [CheckHeadBytesLegal](8c83a112-609f-e2e8-8f8e-8bbe42acf991.htm) | 检查头子节的合法性  Check the legitimacy of the head subsection (重写 [NetMessageBaseCheckHeadBytesLegal(Byte)](3085e36e-54a0-a2ff-cc81-099c11d5b238.htm).) |
| 公共方法 | [CheckMessageMatch](857a414e-7059-b002-b34c-10255390d6b7.htm) | 检查发送的接收的报文是否是匹配的，如果匹配，则返回 1, 如果不匹配且直接返回错误，则返回 0，如果不匹配继续接收，直到匹配或是超时，则返回 -1  If the packet is matched, 1 is returned. If the packet is not matched and an error is returned, 0 is returned. If the packet is not matched, -1 is returned until the packet is matched or times out (继承自 [NetMessageBase](55e9c192-5d9b-9d92-31be-b05cd76d6903.htm)。) |
| 公共方法 | [CheckReceiveDataComplete](909effc7-04ee-4941-593e-2da3f57e6201.htm) | 当消息头报文的长度定义为-1的时候，则使用动态的长度信息，可以使用本方法来判断一个消息是否处于完整的状态。  If the length of the message header is defined as -1, this method can be used to determine whether a message is in the complete state by using dynamic length information. (继承自 [NetMessageBase](55e9c192-5d9b-9d92-31be-b05cd76d6903.htm)。) |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | [GetContentLengthByHeadBytes](bb3f82a8-d286-dab8-03da-b829d4417640.htm) | 从当前的头子节文件中提取出接下来需要接收的数据长度  Extract the length of the data to be received from the current header file |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法 | [GetHeadBytesIdentity](61d7c255-4473-31e7-7f84-dfd957a8826f.htm) | 获取头子节里的消息标识  Get the message ID in the header subsection (继承自 [NetMessageBase](55e9c192-5d9b-9d92-31be-b05cd76d6903.htm)。) |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法 | [PependedUselesByteLength](91167f2b-bebd-0102-29d4-b514cd3aa0a0.htm) | 在接收头报文的时候，返回前置无效的报文头字节长度，默认为0，不处理  When receiving a header message, return the header byte length of the invalid header, the default is 0, and no processing is performed. (继承自 [NetMessageBase](55e9c192-5d9b-9d92-31be-b05cd76d6903.htm)。) |
| 公共方法 | [ToString](762f299e-5af6-4f6a-292c-2ae472bfe3b5.htm) | (继承自 [NetMessageBase](55e9c192-5d9b-9d92-31be-b05cd76d6903.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[SecsHsmsMessage 类](37532ceb-c98f-741c-69ca-cf154a5ad9df.htm)

[HslCommunication.Secs.Message 命名空间](e51b7aa4-2a26-de1d-4009-eae157262c4c.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## CheckHeadBytesLegal 方法 

[原文連結](http://api.hslcommunication.cn/html/8c83a112-609f-e2e8-8f8e-8bbe42acf991.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Secs.Message](../html/e51b7aa4-2a26-de1d-4009-eae157262c4c.htm "HslCommunication.Secs.Message")

[SecsHsmsMessage 类](../html/37532ceb-c98f-741c-69ca-cf154a5ad9df.htm "SecsHsmsMessage 类")

[SecsHsmsMessage 方法](../html/0370daef-e5fa-2b9a-b44f-be4d8e4f91a7.htm "SecsHsmsMessage 方法")

[CheckHeadBytesLegal 方法](../html/8c83a112-609f-e2e8-8f8e-8bbe42acf991.htm "CheckHeadBytesLegal 方法 ")

[GetContentLengthByHeadBytes 方法](../html/bb3f82a8-d286-dab8-03da-b829d4417640.htm "GetContentLengthByHeadBytes 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| SecsHsmsMessageCheckHeadBytesLegal 方法 |

检查头子节的合法性  
Check the legitimacy of the head subsection

**命名空间：**
 [HslCommunication.Secs.Message](e51b7aa4-2a26-de1d-4009-eae157262c4c.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public override bool CheckHeadBytesLegal(
	byte[] token
)
```

```
Public Overrides Function CheckHeadBytesLegal ( 
	token As Byte()
) As Boolean
```

```
public:
virtual bool CheckHeadBytesLegal(
	array<unsigned char>^ token
) override
```

```
abstract CheckHeadBytesLegal : 
        token : byte[] -> bool 
override CheckHeadBytesLegal : 
        token : byte[] -> bool
```

#### 参数

token
:   类型：SystemByte  
    特殊的令牌，有些特殊消息的验证

#### 返回值

类型：Boolean  
是否成功的结果

#### 实现

[INetMessageCheckHeadBytesLegal(Byte)](9fd4bd5c-a68f-3560-d8a7-699dc6fe676e.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[SecsHsmsMessage 类](37532ceb-c98f-741c-69ca-cf154a5ad9df.htm)

[HslCommunication.Secs.Message 命名空间](e51b7aa4-2a26-de1d-4009-eae157262c4c.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GetContentLengthByHeadBytes 方法 

[原文連結](http://api.hslcommunication.cn/html/bb3f82a8-d286-dab8-03da-b829d4417640.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Secs.Message](../html/e51b7aa4-2a26-de1d-4009-eae157262c4c.htm "HslCommunication.Secs.Message")

[SecsHsmsMessage 类](../html/37532ceb-c98f-741c-69ca-cf154a5ad9df.htm "SecsHsmsMessage 类")

[SecsHsmsMessage 方法](../html/0370daef-e5fa-2b9a-b44f-be4d8e4f91a7.htm "SecsHsmsMessage 方法")

[CheckHeadBytesLegal 方法](../html/8c83a112-609f-e2e8-8f8e-8bbe42acf991.htm "CheckHeadBytesLegal 方法 ")

[GetContentLengthByHeadBytes 方法](../html/bb3f82a8-d286-dab8-03da-b829d4417640.htm "GetContentLengthByHeadBytes 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| SecsHsmsMessageGetContentLengthByHeadBytes 方法 |

从当前的头子节文件中提取出接下来需要接收的数据长度  
Extract the length of the data to be received from the current header file

**命名空间：**
 [HslCommunication.Secs.Message](e51b7aa4-2a26-de1d-4009-eae157262c4c.htm)  
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

[SecsHsmsMessage 类](37532ceb-c98f-741c-69ca-cf154a5ad9df.htm)

[HslCommunication.Secs.Message 命名空间](e51b7aa4-2a26-de1d-4009-eae157262c4c.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)