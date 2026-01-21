# HslCommunication - HslCommunication.Profinet.Keyence

> 分類頁數: 30



---
## HslCommunication.Profinet.Keyence

[原文連結](http://api.hslcommunication.cn/html/baaebe57-0dd8-6393-ef35-013cbae668ce.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Keyence](../html/baaebe57-0dd8-6393-ef35-013cbae668ce.htm "HslCommunication.Profinet.Keyence")

[KeyenceDataType 类](../html/8f6cf1c4-f1e8-10dc-3f3c-127cbbdd8458.htm "KeyenceDataType 类")

[KeyenceDLEN1 类](../html/f70c1a08-5573-e24a-5e62-6160f8d51d85.htm "KeyenceDLEN1 类")

[KeyenceKvOld 类](../html/fd02e6bc-ac12-307f-2bff-53d8db6d94f1.htm "KeyenceKvOld 类")

[KeyenceMcAsciiNet 类](../html/a28c120b-fe5c-4d0b-c615-023a5aec93bb.htm "KeyenceMcAsciiNet 类")

[KeyenceMcNet 类](../html/04bd2a21-7ab0-2fb9-f7f1-e0f0ecaf9227.htm "KeyenceMcNet 类")

[KeyenceNanoHelper 类](../html/32735db0-bcfa-ea9c-7d43-0f8753431c6f.htm "KeyenceNanoHelper 类")

[KeyenceNanoSerial 类](../html/8aac53dc-bb22-1ce9-1611-39dfd3ef309b.htm "KeyenceNanoSerial 类")

[KeyenceNanoSerialOverTcp 类](../html/364ddbdb-06f0-d75e-3f56-5df4bad940a8.htm "KeyenceNanoSerialOverTcp 类")

[KeyenceNanoServer 类](../html/fd119063-0ea1-76d5-8975-a37d5e682ffe.htm "KeyenceNanoServer 类")

[KeyencePLCS 枚举](../html/e59e6dae-084d-60db-fe36-834b14311ee9.htm "KeyencePLCS 枚举")

[KeyenceSR2000Serial 类](../html/77f8d8da-8005-e4ca-598a-7ea8c5e9d40c.htm "KeyenceSR2000Serial 类")

[KeyenceSR2000SeriesTcp 类](../html/4430e548-fc58-8c52-e767-63500f424f91.htm "KeyenceSR2000SeriesTcp 类")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslCommunication.Profinet.Keyence 命名空间 |

[缺少 "N:HslCommunication.Profinet.Keyence" 的 <summary> 文档]

![](../icons/SectionExpanded.png)类

|  | 类 | 说明 |
| --- | --- | --- |
| 公共类 | [KeyenceDataType](8f6cf1c4-f1e8-10dc-3f3c-127cbbdd8458.htm) | Keyence PLC的数据类型，此处包含了几个常用的类型 |
| 公共类 | [KeyenceDLEN1](f70c1a08-5573-e24a-5e62-6160f8d51d85.htm) | 基恩士的数字传感器的以太网模块，可以同时连接并读取多个传感器模块的功能代码 |
| 公共类 | [KeyenceKvOld](fd02e6bc-ac12-307f-2bff-53d8db6d94f1.htm) | 老型号的基恩士类，适用于 Kv-10xx, 16xx, 24xx, 40xx, kv-80, kv-300，其中 xx 表示 AR/AT/DR/DT  Old Keenes class, suitable for Kv-10xx, 16xx, 24xx, 40xx, kv-80, kv-300, xx means AR/AT/DR/DT |
| 公共类 | [KeyenceMcAsciiNet](a28c120b-fe5c-4d0b-c615-023a5aec93bb.htm) | 基恩士PLC的数据通信类，使用QnA兼容3E帧的通信协议实现，使用ASCII的格式，地址格式需要进行转换成三菱的格式，详细参照备注说明  Keyence PLC's data communication class is implemented using QnA compatible 3E frame communication protocol. It uses ascii format. The address format needs to be converted to Mitsubishi format. |
| 公共类 | [KeyenceMcNet](04bd2a21-7ab0-2fb9-f7f1-e0f0ecaf9227.htm) | 基恩士PLC的数据通信类，使用QnA兼容3E帧的通信协议实现，使用二进制的格式，地址同时支持三菱的地址格式及基恩士自身的地址格式，详细参照备注说明  The data communication class of KEYENCE PLC is implemented using a QnA-compatible 3E frame communication protocol, using binary format, and the address supports both Mitsubishi's address format and Keyence's own address format, please refer to the remarks for details |
| 公共类 | [KeyenceNanoHelper](32735db0-bcfa-ea9c-7d43-0f8753431c6f.htm) | KeyenceNano的基本辅助方法 |
| 公共类 | [KeyenceNanoSerial](8aac53dc-bb22-1ce9-1611-39dfd3ef309b.htm) | 基恩士KV上位链路串口通信的对象,适用于Nano系列串口数据,KV1000以及L20V通信模块，地址格式参考api文档  Keyence KV upper link serial communication object, suitable for Nano series serial data, and L20V communication module, please refer to api document for address format |
| 公共类 | [KeyenceNanoSerialOverTcp](364ddbdb-06f0-d75e-3f56-5df4bad940a8.htm) | 基恩士KV上位链路协议的通信对象,适用于KV5000/5500/3000,KV1000,KV700,以及L20V通信模块，本类是基于tcp通信  The communication object of KEYENCE KV upper link protocol is suitable for KV5000/5500/3000, KV1000, KV700, and L20V communication modules. This type is based on tcp communication |
| 公共类 | [KeyenceNanoServer](fd119063-0ea1-76d5-8975-a37d5e682ffe.htm) | 基恩士的上位链路协议的虚拟服务器 |
| 公共类 | [KeyenceSR2000Serial](77f8d8da-8005-e4ca-598a-7ea8c5e9d40c.htm) | 基恩士的SR2000的扫码设备，可以进行简单的交互 |
| 公共类 | [KeyenceSR2000SeriesTcp](4430e548-fc58-8c52-e767-63500f424f91.htm) | 基恩士的SR2000的扫码设备，可以进行简单的交互 |

![](../icons/SectionExpanded.png)枚举

|  | 枚举 | 说明 |
| --- | --- | --- |
| 公共枚举 | [KeyencePLCS](e59e6dae-084d-60db-fe36-834b14311ee9.htm) | 基恩士PLC的各种系列选择 |

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## KeyenceDataType 类

[原文連結](http://api.hslcommunication.cn/html/8f6cf1c4-f1e8-10dc-3f3c-127cbbdd8458.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Keyence](../html/baaebe57-0dd8-6393-ef35-013cbae668ce.htm "HslCommunication.Profinet.Keyence")

[KeyenceDataType 类](../html/8f6cf1c4-f1e8-10dc-3f3c-127cbbdd8458.htm "KeyenceDataType 类")

[KeyenceDataType 构造函数](../html/56d011d9-b961-b147-ddc7-5d88f0e5edb3.htm "KeyenceDataType 构造函数 ")

[KeyenceDataType 属性](../html/2f5f201d-f0f3-7c00-a5fc-e03a079a6c3c.htm "KeyenceDataType 属性")

[KeyenceDataType 方法](../html/c2189cf3-c6cc-e650-a5ec-d428fb4971f0.htm "KeyenceDataType 方法")

[KeyenceDataType 字段](../html/f015109d-1e12-3cc5-e7ca-0a6ca880f9c3.htm "KeyenceDataType 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| KeyenceDataType 类 |

Keyence PLC的数据类型，此处包含了几个常用的类型

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  HslCommunication.Profinet.KeyenceKeyenceDataType

**命名空间：**
 [HslCommunication.Profinet.Keyence](baaebe57-0dd8-6393-ef35-013cbae668ce.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class KeyenceDataType
```

```
Public Class KeyenceDataType
```

```
public ref class KeyenceDataType
```

```
type KeyenceDataType =  class end
```

KeyenceDataType 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [KeyenceDataType](56d011d9-b961-b147-ddc7-5d88f0e5edb3.htm) | 如果您清楚类型代号，可以根据值进行扩展 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [AsciiCode](1643513b-6aab-f698-0c0a-68072a461915.htm) | 当以ASCII格式通讯时的类型描述 |
| 公共属性 | [DataCode](1907d7d7-3ca3-964d-8726-b5d76dba616c.htm) | 类型的代号值 |
| 公共属性 | [DataType](c2000d77-86ff-8b0c-92f3-49cb61533fd4.htm) | 数据的类型，0代表按字，1代表按位 |
| 公共属性 | [FromBase](7d0b17d8-3508-f567-2b16-e1faaa72c0ca.htm) | 指示地址是10进制，还是16进制的 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法 | ToString | (继承自 Object。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)字段

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共字段静态成员 | [B](742136fc-a835-2b05-cb52-8487888b30ad.htm) | 链接继电器 |
| 公共字段静态成员 | [CN](cc4e540f-e45a-8809-0aca-56c7cd565636.htm) | 计数器（当前值） |
| 公共字段静态成员 | [CS](6a9b8328-dc8d-3a7e-c53c-25918c990d16.htm) | 计数器（接点） |
| 公共字段静态成员 | [D](0c4edcf3-95ac-4648-2215-97819edfde79.htm) | 数据存储器 |
| 公共字段静态成员 | [L](7083f821-41f7-5b22-061b-98d732b32fbf.htm) | 锁存继电器 |
| 公共字段静态成员 | [M](fe5f1aaa-506a-dffa-9951-46b62d61e008.htm) | 内部辅助继电器 |
| 公共字段静态成员 | [R](4d28ca66-e1f7-4d28-2afa-b532d1720b06.htm) | 文件寄存器 |
| 公共字段静态成员 | [SD](fe4e7372-7f6d-8c31-7cb9-2a96d589f25b.htm) | 控制存储器 |
| 公共字段静态成员 | [SM](7f39bf89-55d4-ca3a-12eb-d6abce4ebe02.htm) | 控制继电器 |
| 公共字段静态成员 | [TN](7af1a704-0c31-70f4-7126-6eda3675137e.htm) | 计时器（当前值） |
| 公共字段静态成员 | [TS](0bb3d647-793c-84e9-3a36-9074798194f9.htm) | 计时器（接点） |
| 公共字段静态成员 | [W](40cd4bfc-d820-d6d1-fcd7-783e83dd585b.htm) | 链路寄存器 |
| 公共字段静态成员 | [X](d61985b8-6990-aa2e-6abc-427459a7725f.htm) | X输入继电器 |
| 公共字段静态成员 | [Y](335aa99c-cdd4-82a1-bd60-00ffbe9c254e.htm) | Y输出继电器 |
| 公共字段静态成员 | [ZR](4905c593-832b-f615-89cd-e3431fb4a66d.htm) | 文件寄存器 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.Profinet.Keyence 命名空间](baaebe57-0dd8-6393-ef35-013cbae668ce.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## KeyenceDataType 构造函数 

[原文連結](http://api.hslcommunication.cn/html/56d011d9-b961-b147-ddc7-5d88f0e5edb3.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Keyence](../html/baaebe57-0dd8-6393-ef35-013cbae668ce.htm "HslCommunication.Profinet.Keyence")

[KeyenceDataType 类](../html/8f6cf1c4-f1e8-10dc-3f3c-127cbbdd8458.htm "KeyenceDataType 类")

[KeyenceDataType 构造函数](../html/56d011d9-b961-b147-ddc7-5d88f0e5edb3.htm "KeyenceDataType 构造函数 ")

[KeyenceDataType 属性](../html/2f5f201d-f0f3-7c00-a5fc-e03a079a6c3c.htm "KeyenceDataType 属性")

[KeyenceDataType 方法](../html/c2189cf3-c6cc-e650-a5ec-d428fb4971f0.htm "KeyenceDataType 方法")

[KeyenceDataType 字段](../html/f015109d-1e12-3cc5-e7ca-0a6ca880f9c3.htm "KeyenceDataType 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| KeyenceDataType 构造函数 |

如果您清楚类型代号，可以根据值进行扩展

**命名空间：**
 [HslCommunication.Profinet.Keyence](baaebe57-0dd8-6393-ef35-013cbae668ce.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public KeyenceDataType(
	byte code,
	byte type,
	string asciiCode,
	int fromBase
)
```

```
Public Sub New ( 
	code As Byte,
	type As Byte,
	asciiCode As String,
	fromBase As Integer
)
```

```
public:
KeyenceDataType(
	unsigned char code, 
	unsigned char type, 
	String^ asciiCode, 
	int fromBase
)
```

```
new : 
        code : byte * 
        type : byte * 
        asciiCode : string * 
        fromBase : int -> KeyenceDataType
```

#### 参数

code
:   类型：SystemByte  
    数据类型的代号

type
:   类型：SystemByte  
    0或1，默认为0

asciiCode
:   类型：SystemString  
    ASCII格式的类型信息

fromBase
:   类型：SystemInt32  
    指示地址的多少进制的，10或是16

![](../icons/SectionExpanded.png)参见

#### 引用

[KeyenceDataType 类](8f6cf1c4-f1e8-10dc-3f3c-127cbbdd8458.htm)

[HslCommunication.Profinet.Keyence 命名空间](baaebe57-0dd8-6393-ef35-013cbae668ce.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## KeyenceDataType 属性

[原文連結](http://api.hslcommunication.cn/html/2f5f201d-f0f3-7c00-a5fc-e03a079a6c3c.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Keyence](../html/baaebe57-0dd8-6393-ef35-013cbae668ce.htm "HslCommunication.Profinet.Keyence")

[KeyenceDataType 类](../html/8f6cf1c4-f1e8-10dc-3f3c-127cbbdd8458.htm "KeyenceDataType 类")

[KeyenceDataType 属性](../html/2f5f201d-f0f3-7c00-a5fc-e03a079a6c3c.htm "KeyenceDataType 属性")

[AsciiCode 属性](../html/1643513b-6aab-f698-0c0a-68072a461915.htm "AsciiCode 属性 ")

[DataCode 属性](../html/1907d7d7-3ca3-964d-8726-b5d76dba616c.htm "DataCode 属性 ")

[DataType 属性](../html/c2000d77-86ff-8b0c-92f3-49cb61533fd4.htm "DataType 属性 ")

[FromBase 属性](../html/7d0b17d8-3508-f567-2b16-e1faaa72c0ca.htm "FromBase 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| KeyenceDataType 属性 |

[KeyenceDataType](8f6cf1c4-f1e8-10dc-3f3c-127cbbdd8458.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [AsciiCode](1643513b-6aab-f698-0c0a-68072a461915.htm) | 当以ASCII格式通讯时的类型描述 |
| 公共属性 | [DataCode](1907d7d7-3ca3-964d-8726-b5d76dba616c.htm) | 类型的代号值 |
| 公共属性 | [DataType](c2000d77-86ff-8b0c-92f3-49cb61533fd4.htm) | 数据的类型，0代表按字，1代表按位 |
| 公共属性 | [FromBase](7d0b17d8-3508-f567-2b16-e1faaa72c0ca.htm) | 指示地址是10进制，还是16进制的 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[KeyenceDataType 类](8f6cf1c4-f1e8-10dc-3f3c-127cbbdd8458.htm)

[HslCommunication.Profinet.Keyence 命名空间](baaebe57-0dd8-6393-ef35-013cbae668ce.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## AsciiCode 属性 

[原文連結](http://api.hslcommunication.cn/html/1643513b-6aab-f698-0c0a-68072a461915.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Keyence](../html/baaebe57-0dd8-6393-ef35-013cbae668ce.htm "HslCommunication.Profinet.Keyence")

[KeyenceDataType 类](../html/8f6cf1c4-f1e8-10dc-3f3c-127cbbdd8458.htm "KeyenceDataType 类")

[KeyenceDataType 属性](../html/2f5f201d-f0f3-7c00-a5fc-e03a079a6c3c.htm "KeyenceDataType 属性")

[AsciiCode 属性](../html/1643513b-6aab-f698-0c0a-68072a461915.htm "AsciiCode 属性 ")

[DataCode 属性](../html/1907d7d7-3ca3-964d-8726-b5d76dba616c.htm "DataCode 属性 ")

[DataType 属性](../html/c2000d77-86ff-8b0c-92f3-49cb61533fd4.htm "DataType 属性 ")

[FromBase 属性](../html/7d0b17d8-3508-f567-2b16-e1faaa72c0ca.htm "FromBase 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| KeyenceDataTypeAsciiCode 属性 |

当以ASCII格式通讯时的类型描述

**命名空间：**
 [HslCommunication.Profinet.Keyence](baaebe57-0dd8-6393-ef35-013cbae668ce.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public string AsciiCode { get; }
```

```
Public ReadOnly Property AsciiCode As String
	Get
```

```
public:
property String^ AsciiCode {
	String^ get ();
}
```

```
member AsciiCode : string with get
```

#### 属性值

类型：String

![](../icons/SectionExpanded.png)参见

#### 引用

[KeyenceDataType 类](8f6cf1c4-f1e8-10dc-3f3c-127cbbdd8458.htm)

[HslCommunication.Profinet.Keyence 命名空间](baaebe57-0dd8-6393-ef35-013cbae668ce.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## DataCode 属性 

[原文連結](http://api.hslcommunication.cn/html/1907d7d7-3ca3-964d-8726-b5d76dba616c.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Keyence](../html/baaebe57-0dd8-6393-ef35-013cbae668ce.htm "HslCommunication.Profinet.Keyence")

[KeyenceDataType 类](../html/8f6cf1c4-f1e8-10dc-3f3c-127cbbdd8458.htm "KeyenceDataType 类")

[KeyenceDataType 属性](../html/2f5f201d-f0f3-7c00-a5fc-e03a079a6c3c.htm "KeyenceDataType 属性")

[AsciiCode 属性](../html/1643513b-6aab-f698-0c0a-68072a461915.htm "AsciiCode 属性 ")

[DataCode 属性](../html/1907d7d7-3ca3-964d-8726-b5d76dba616c.htm "DataCode 属性 ")

[DataType 属性](../html/c2000d77-86ff-8b0c-92f3-49cb61533fd4.htm "DataType 属性 ")

[FromBase 属性](../html/7d0b17d8-3508-f567-2b16-e1faaa72c0ca.htm "FromBase 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| KeyenceDataTypeDataCode 属性 |

类型的代号值

**命名空间：**
 [HslCommunication.Profinet.Keyence](baaebe57-0dd8-6393-ef35-013cbae668ce.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public byte DataCode { get; }
```

```
Public ReadOnly Property DataCode As Byte
	Get
```

```
public:
property unsigned char DataCode {
	unsigned char get ();
}
```

```
member DataCode : byte with get
```

#### 属性值

类型：Byte

![](../icons/SectionExpanded.png)参见

#### 引用

[KeyenceDataType 类](8f6cf1c4-f1e8-10dc-3f3c-127cbbdd8458.htm)

[HslCommunication.Profinet.Keyence 命名空间](baaebe57-0dd8-6393-ef35-013cbae668ce.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## DataType 属性 

[原文連結](http://api.hslcommunication.cn/html/c2000d77-86ff-8b0c-92f3-49cb61533fd4.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Keyence](../html/baaebe57-0dd8-6393-ef35-013cbae668ce.htm "HslCommunication.Profinet.Keyence")

[KeyenceDataType 类](../html/8f6cf1c4-f1e8-10dc-3f3c-127cbbdd8458.htm "KeyenceDataType 类")

[KeyenceDataType 属性](../html/2f5f201d-f0f3-7c00-a5fc-e03a079a6c3c.htm "KeyenceDataType 属性")

[AsciiCode 属性](../html/1643513b-6aab-f698-0c0a-68072a461915.htm "AsciiCode 属性 ")

[DataCode 属性](../html/1907d7d7-3ca3-964d-8726-b5d76dba616c.htm "DataCode 属性 ")

[DataType 属性](../html/c2000d77-86ff-8b0c-92f3-49cb61533fd4.htm "DataType 属性 ")

[FromBase 属性](../html/7d0b17d8-3508-f567-2b16-e1faaa72c0ca.htm "FromBase 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| KeyenceDataTypeDataType 属性 |

数据的类型，0代表按字，1代表按位

**命名空间：**
 [HslCommunication.Profinet.Keyence](baaebe57-0dd8-6393-ef35-013cbae668ce.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public byte DataType { get; }
```

```
Public ReadOnly Property DataType As Byte
	Get
```

```
public:
property unsigned char DataType {
	unsigned char get ();
}
```

```
member DataType : byte with get
```

#### 属性值

类型：Byte

![](../icons/SectionExpanded.png)参见

#### 引用

[KeyenceDataType 类](8f6cf1c4-f1e8-10dc-3f3c-127cbbdd8458.htm)

[HslCommunication.Profinet.Keyence 命名空间](baaebe57-0dd8-6393-ef35-013cbae668ce.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## FromBase 属性 

[原文連結](http://api.hslcommunication.cn/html/7d0b17d8-3508-f567-2b16-e1faaa72c0ca.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Keyence](../html/baaebe57-0dd8-6393-ef35-013cbae668ce.htm "HslCommunication.Profinet.Keyence")

[KeyenceDataType 类](../html/8f6cf1c4-f1e8-10dc-3f3c-127cbbdd8458.htm "KeyenceDataType 类")

[KeyenceDataType 属性](../html/2f5f201d-f0f3-7c00-a5fc-e03a079a6c3c.htm "KeyenceDataType 属性")

[AsciiCode 属性](../html/1643513b-6aab-f698-0c0a-68072a461915.htm "AsciiCode 属性 ")

[DataCode 属性](../html/1907d7d7-3ca3-964d-8726-b5d76dba616c.htm "DataCode 属性 ")

[DataType 属性](../html/c2000d77-86ff-8b0c-92f3-49cb61533fd4.htm "DataType 属性 ")

[FromBase 属性](../html/7d0b17d8-3508-f567-2b16-e1faaa72c0ca.htm "FromBase 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| KeyenceDataTypeFromBase 属性 |

指示地址是10进制，还是16进制的

**命名空间：**
 [HslCommunication.Profinet.Keyence](baaebe57-0dd8-6393-ef35-013cbae668ce.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public int FromBase { get; }
```

```
Public ReadOnly Property FromBase As Integer
	Get
```

```
public:
property int FromBase {
	int get ();
}
```

```
member FromBase : int with get
```

#### 属性值

类型：Int32

![](../icons/SectionExpanded.png)参见

#### 引用

[KeyenceDataType 类](8f6cf1c4-f1e8-10dc-3f3c-127cbbdd8458.htm)

[HslCommunication.Profinet.Keyence 命名空间](baaebe57-0dd8-6393-ef35-013cbae668ce.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## KeyenceDataType 方法

[原文連結](http://api.hslcommunication.cn/html/c2189cf3-c6cc-e650-a5ec-d428fb4971f0.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Keyence](../html/baaebe57-0dd8-6393-ef35-013cbae668ce.htm "HslCommunication.Profinet.Keyence")

[KeyenceDataType 类](../html/8f6cf1c4-f1e8-10dc-3f3c-127cbbdd8458.htm "KeyenceDataType 类")

[KeyenceDataType 构造函数](../html/56d011d9-b961-b147-ddc7-5d88f0e5edb3.htm "KeyenceDataType 构造函数 ")

[KeyenceDataType 属性](../html/2f5f201d-f0f3-7c00-a5fc-e03a079a6c3c.htm "KeyenceDataType 属性")

[KeyenceDataType 方法](../html/c2189cf3-c6cc-e650-a5ec-d428fb4971f0.htm "KeyenceDataType 方法")

[KeyenceDataType 字段](../html/f015109d-1e12-3cc5-e7ca-0a6ca880f9c3.htm "KeyenceDataType 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| KeyenceDataType 方法 |

[KeyenceDataType](8f6cf1c4-f1e8-10dc-3f3c-127cbbdd8458.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
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

[KeyenceDataType 类](8f6cf1c4-f1e8-10dc-3f3c-127cbbdd8458.htm)

[HslCommunication.Profinet.Keyence 命名空间](baaebe57-0dd8-6393-ef35-013cbae668ce.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## KeyenceDataType 字段

[原文連結](http://api.hslcommunication.cn/html/f015109d-1e12-3cc5-e7ca-0a6ca880f9c3.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Keyence](../html/baaebe57-0dd8-6393-ef35-013cbae668ce.htm "HslCommunication.Profinet.Keyence")

[KeyenceDataType 类](../html/8f6cf1c4-f1e8-10dc-3f3c-127cbbdd8458.htm "KeyenceDataType 类")

[KeyenceDataType 字段](../html/f015109d-1e12-3cc5-e7ca-0a6ca880f9c3.htm "KeyenceDataType 字段")

[B 字段](../html/742136fc-a835-2b05-cb52-8487888b30ad.htm "B 字段")

[CN 字段](../html/cc4e540f-e45a-8809-0aca-56c7cd565636.htm "CN 字段")

[CS 字段](../html/6a9b8328-dc8d-3a7e-c53c-25918c990d16.htm "CS 字段")

[D 字段](../html/0c4edcf3-95ac-4648-2215-97819edfde79.htm "D 字段")

[L 字段](../html/7083f821-41f7-5b22-061b-98d732b32fbf.htm "L 字段")

[M 字段](../html/fe5f1aaa-506a-dffa-9951-46b62d61e008.htm "M 字段")

[R 字段](../html/4d28ca66-e1f7-4d28-2afa-b532d1720b06.htm "R 字段")

[SD 字段](../html/fe4e7372-7f6d-8c31-7cb9-2a96d589f25b.htm "SD 字段")

[SM 字段](../html/7f39bf89-55d4-ca3a-12eb-d6abce4ebe02.htm "SM 字段")

[TN 字段](../html/7af1a704-0c31-70f4-7126-6eda3675137e.htm "TN 字段")

[TS 字段](../html/0bb3d647-793c-84e9-3a36-9074798194f9.htm "TS 字段")

[W 字段](../html/40cd4bfc-d820-d6d1-fcd7-783e83dd585b.htm "W 字段")

[X 字段](../html/d61985b8-6990-aa2e-6abc-427459a7725f.htm "X 字段")

[Y 字段](../html/335aa99c-cdd4-82a1-bd60-00ffbe9c254e.htm "Y 字段")

[ZR 字段](../html/4905c593-832b-f615-89cd-e3431fb4a66d.htm "ZR 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| KeyenceDataType 字段 |

[KeyenceDataType](8f6cf1c4-f1e8-10dc-3f3c-127cbbdd8458.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)字段

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共字段静态成员 | [B](742136fc-a835-2b05-cb52-8487888b30ad.htm) | 链接继电器 |
| 公共字段静态成员 | [CN](cc4e540f-e45a-8809-0aca-56c7cd565636.htm) | 计数器（当前值） |
| 公共字段静态成员 | [CS](6a9b8328-dc8d-3a7e-c53c-25918c990d16.htm) | 计数器（接点） |
| 公共字段静态成员 | [D](0c4edcf3-95ac-4648-2215-97819edfde79.htm) | 数据存储器 |
| 公共字段静态成员 | [L](7083f821-41f7-5b22-061b-98d732b32fbf.htm) | 锁存继电器 |
| 公共字段静态成员 | [M](fe5f1aaa-506a-dffa-9951-46b62d61e008.htm) | 内部辅助继电器 |
| 公共字段静态成员 | [R](4d28ca66-e1f7-4d28-2afa-b532d1720b06.htm) | 文件寄存器 |
| 公共字段静态成员 | [SD](fe4e7372-7f6d-8c31-7cb9-2a96d589f25b.htm) | 控制存储器 |
| 公共字段静态成员 | [SM](7f39bf89-55d4-ca3a-12eb-d6abce4ebe02.htm) | 控制继电器 |
| 公共字段静态成员 | [TN](7af1a704-0c31-70f4-7126-6eda3675137e.htm) | 计时器（当前值） |
| 公共字段静态成员 | [TS](0bb3d647-793c-84e9-3a36-9074798194f9.htm) | 计时器（接点） |
| 公共字段静态成员 | [W](40cd4bfc-d820-d6d1-fcd7-783e83dd585b.htm) | 链路寄存器 |
| 公共字段静态成员 | [X](d61985b8-6990-aa2e-6abc-427459a7725f.htm) | X输入继电器 |
| 公共字段静态成员 | [Y](335aa99c-cdd4-82a1-bd60-00ffbe9c254e.htm) | Y输出继电器 |
| 公共字段静态成员 | [ZR](4905c593-832b-f615-89cd-e3431fb4a66d.htm) | 文件寄存器 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[KeyenceDataType 类](8f6cf1c4-f1e8-10dc-3f3c-127cbbdd8458.htm)

[HslCommunication.Profinet.Keyence 命名空间](baaebe57-0dd8-6393-ef35-013cbae668ce.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## B 字段

[原文連結](http://api.hslcommunication.cn/html/742136fc-a835-2b05-cb52-8487888b30ad.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Keyence](../html/baaebe57-0dd8-6393-ef35-013cbae668ce.htm "HslCommunication.Profinet.Keyence")

[KeyenceDataType 类](../html/8f6cf1c4-f1e8-10dc-3f3c-127cbbdd8458.htm "KeyenceDataType 类")

[KeyenceDataType 字段](../html/f015109d-1e12-3cc5-e7ca-0a6ca880f9c3.htm "KeyenceDataType 字段")

[B 字段](../html/742136fc-a835-2b05-cb52-8487888b30ad.htm "B 字段")

[CN 字段](../html/cc4e540f-e45a-8809-0aca-56c7cd565636.htm "CN 字段")

[CS 字段](../html/6a9b8328-dc8d-3a7e-c53c-25918c990d16.htm "CS 字段")

[D 字段](../html/0c4edcf3-95ac-4648-2215-97819edfde79.htm "D 字段")

[L 字段](../html/7083f821-41f7-5b22-061b-98d732b32fbf.htm "L 字段")

[M 字段](../html/fe5f1aaa-506a-dffa-9951-46b62d61e008.htm "M 字段")

[R 字段](../html/4d28ca66-e1f7-4d28-2afa-b532d1720b06.htm "R 字段")

[SD 字段](../html/fe4e7372-7f6d-8c31-7cb9-2a96d589f25b.htm "SD 字段")

[SM 字段](../html/7f39bf89-55d4-ca3a-12eb-d6abce4ebe02.htm "SM 字段")

[TN 字段](../html/7af1a704-0c31-70f4-7126-6eda3675137e.htm "TN 字段")

[TS 字段](../html/0bb3d647-793c-84e9-3a36-9074798194f9.htm "TS 字段")

[W 字段](../html/40cd4bfc-d820-d6d1-fcd7-783e83dd585b.htm "W 字段")

[X 字段](../html/d61985b8-6990-aa2e-6abc-427459a7725f.htm "X 字段")

[Y 字段](../html/335aa99c-cdd4-82a1-bd60-00ffbe9c254e.htm "Y 字段")

[ZR 字段](../html/4905c593-832b-f615-89cd-e3431fb4a66d.htm "ZR 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| KeyenceDataTypeB 字段 |

链接继电器

**命名空间：**
 [HslCommunication.Profinet.Keyence](baaebe57-0dd8-6393-ef35-013cbae668ce.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static readonly KeyenceDataType B
```

```
Public Shared ReadOnly B As KeyenceDataType
```

```
public:
static initonly KeyenceDataType^ B
```

```
static val B: KeyenceDataType
```

#### 字段值

类型：[KeyenceDataType](8f6cf1c4-f1e8-10dc-3f3c-127cbbdd8458.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[KeyenceDataType 类](8f6cf1c4-f1e8-10dc-3f3c-127cbbdd8458.htm)

[HslCommunication.Profinet.Keyence 命名空间](baaebe57-0dd8-6393-ef35-013cbae668ce.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## CN 字段

[原文連結](http://api.hslcommunication.cn/html/cc4e540f-e45a-8809-0aca-56c7cd565636.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Keyence](../html/baaebe57-0dd8-6393-ef35-013cbae668ce.htm "HslCommunication.Profinet.Keyence")

[KeyenceDataType 类](../html/8f6cf1c4-f1e8-10dc-3f3c-127cbbdd8458.htm "KeyenceDataType 类")

[KeyenceDataType 字段](../html/f015109d-1e12-3cc5-e7ca-0a6ca880f9c3.htm "KeyenceDataType 字段")

[B 字段](../html/742136fc-a835-2b05-cb52-8487888b30ad.htm "B 字段")

[CN 字段](../html/cc4e540f-e45a-8809-0aca-56c7cd565636.htm "CN 字段")

[CS 字段](../html/6a9b8328-dc8d-3a7e-c53c-25918c990d16.htm "CS 字段")

[D 字段](../html/0c4edcf3-95ac-4648-2215-97819edfde79.htm "D 字段")

[L 字段](../html/7083f821-41f7-5b22-061b-98d732b32fbf.htm "L 字段")

[M 字段](../html/fe5f1aaa-506a-dffa-9951-46b62d61e008.htm "M 字段")

[R 字段](../html/4d28ca66-e1f7-4d28-2afa-b532d1720b06.htm "R 字段")

[SD 字段](../html/fe4e7372-7f6d-8c31-7cb9-2a96d589f25b.htm "SD 字段")

[SM 字段](../html/7f39bf89-55d4-ca3a-12eb-d6abce4ebe02.htm "SM 字段")

[TN 字段](../html/7af1a704-0c31-70f4-7126-6eda3675137e.htm "TN 字段")

[TS 字段](../html/0bb3d647-793c-84e9-3a36-9074798194f9.htm "TS 字段")

[W 字段](../html/40cd4bfc-d820-d6d1-fcd7-783e83dd585b.htm "W 字段")

[X 字段](../html/d61985b8-6990-aa2e-6abc-427459a7725f.htm "X 字段")

[Y 字段](../html/335aa99c-cdd4-82a1-bd60-00ffbe9c254e.htm "Y 字段")

[ZR 字段](../html/4905c593-832b-f615-89cd-e3431fb4a66d.htm "ZR 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| KeyenceDataTypeCN 字段 |

计数器（当前值）

**命名空间：**
 [HslCommunication.Profinet.Keyence](baaebe57-0dd8-6393-ef35-013cbae668ce.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static readonly KeyenceDataType CN
```

```
Public Shared ReadOnly CN As KeyenceDataType
```

```
public:
static initonly KeyenceDataType^ CN
```

```
static val CN: KeyenceDataType
```

#### 字段值

类型：[KeyenceDataType](8f6cf1c4-f1e8-10dc-3f3c-127cbbdd8458.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[KeyenceDataType 类](8f6cf1c4-f1e8-10dc-3f3c-127cbbdd8458.htm)

[HslCommunication.Profinet.Keyence 命名空间](baaebe57-0dd8-6393-ef35-013cbae668ce.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## CS 字段

[原文連結](http://api.hslcommunication.cn/html/6a9b8328-dc8d-3a7e-c53c-25918c990d16.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Keyence](../html/baaebe57-0dd8-6393-ef35-013cbae668ce.htm "HslCommunication.Profinet.Keyence")

[KeyenceDataType 类](../html/8f6cf1c4-f1e8-10dc-3f3c-127cbbdd8458.htm "KeyenceDataType 类")

[KeyenceDataType 字段](../html/f015109d-1e12-3cc5-e7ca-0a6ca880f9c3.htm "KeyenceDataType 字段")

[B 字段](../html/742136fc-a835-2b05-cb52-8487888b30ad.htm "B 字段")

[CN 字段](../html/cc4e540f-e45a-8809-0aca-56c7cd565636.htm "CN 字段")

[CS 字段](../html/6a9b8328-dc8d-3a7e-c53c-25918c990d16.htm "CS 字段")

[D 字段](../html/0c4edcf3-95ac-4648-2215-97819edfde79.htm "D 字段")

[L 字段](../html/7083f821-41f7-5b22-061b-98d732b32fbf.htm "L 字段")

[M 字段](../html/fe5f1aaa-506a-dffa-9951-46b62d61e008.htm "M 字段")

[R 字段](../html/4d28ca66-e1f7-4d28-2afa-b532d1720b06.htm "R 字段")

[SD 字段](../html/fe4e7372-7f6d-8c31-7cb9-2a96d589f25b.htm "SD 字段")

[SM 字段](../html/7f39bf89-55d4-ca3a-12eb-d6abce4ebe02.htm "SM 字段")

[TN 字段](../html/7af1a704-0c31-70f4-7126-6eda3675137e.htm "TN 字段")

[TS 字段](../html/0bb3d647-793c-84e9-3a36-9074798194f9.htm "TS 字段")

[W 字段](../html/40cd4bfc-d820-d6d1-fcd7-783e83dd585b.htm "W 字段")

[X 字段](../html/d61985b8-6990-aa2e-6abc-427459a7725f.htm "X 字段")

[Y 字段](../html/335aa99c-cdd4-82a1-bd60-00ffbe9c254e.htm "Y 字段")

[ZR 字段](../html/4905c593-832b-f615-89cd-e3431fb4a66d.htm "ZR 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| KeyenceDataTypeCS 字段 |

计数器（接点）

**命名空间：**
 [HslCommunication.Profinet.Keyence](baaebe57-0dd8-6393-ef35-013cbae668ce.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static readonly KeyenceDataType CS
```

```
Public Shared ReadOnly CS As KeyenceDataType
```

```
public:
static initonly KeyenceDataType^ CS
```

```
static val CS: KeyenceDataType
```

#### 字段值

类型：[KeyenceDataType](8f6cf1c4-f1e8-10dc-3f3c-127cbbdd8458.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[KeyenceDataType 类](8f6cf1c4-f1e8-10dc-3f3c-127cbbdd8458.htm)

[HslCommunication.Profinet.Keyence 命名空间](baaebe57-0dd8-6393-ef35-013cbae668ce.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## D 字段

[原文連結](http://api.hslcommunication.cn/html/0c4edcf3-95ac-4648-2215-97819edfde79.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Keyence](../html/baaebe57-0dd8-6393-ef35-013cbae668ce.htm "HslCommunication.Profinet.Keyence")

[KeyenceDataType 类](../html/8f6cf1c4-f1e8-10dc-3f3c-127cbbdd8458.htm "KeyenceDataType 类")

[KeyenceDataType 字段](../html/f015109d-1e12-3cc5-e7ca-0a6ca880f9c3.htm "KeyenceDataType 字段")

[B 字段](../html/742136fc-a835-2b05-cb52-8487888b30ad.htm "B 字段")

[CN 字段](../html/cc4e540f-e45a-8809-0aca-56c7cd565636.htm "CN 字段")

[CS 字段](../html/6a9b8328-dc8d-3a7e-c53c-25918c990d16.htm "CS 字段")

[D 字段](../html/0c4edcf3-95ac-4648-2215-97819edfde79.htm "D 字段")

[L 字段](../html/7083f821-41f7-5b22-061b-98d732b32fbf.htm "L 字段")

[M 字段](../html/fe5f1aaa-506a-dffa-9951-46b62d61e008.htm "M 字段")

[R 字段](../html/4d28ca66-e1f7-4d28-2afa-b532d1720b06.htm "R 字段")

[SD 字段](../html/fe4e7372-7f6d-8c31-7cb9-2a96d589f25b.htm "SD 字段")

[SM 字段](../html/7f39bf89-55d4-ca3a-12eb-d6abce4ebe02.htm "SM 字段")

[TN 字段](../html/7af1a704-0c31-70f4-7126-6eda3675137e.htm "TN 字段")

[TS 字段](../html/0bb3d647-793c-84e9-3a36-9074798194f9.htm "TS 字段")

[W 字段](../html/40cd4bfc-d820-d6d1-fcd7-783e83dd585b.htm "W 字段")

[X 字段](../html/d61985b8-6990-aa2e-6abc-427459a7725f.htm "X 字段")

[Y 字段](../html/335aa99c-cdd4-82a1-bd60-00ffbe9c254e.htm "Y 字段")

[ZR 字段](../html/4905c593-832b-f615-89cd-e3431fb4a66d.htm "ZR 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| KeyenceDataTypeD 字段 |

数据存储器

**命名空间：**
 [HslCommunication.Profinet.Keyence](baaebe57-0dd8-6393-ef35-013cbae668ce.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static readonly KeyenceDataType D
```

```
Public Shared ReadOnly D As KeyenceDataType
```

```
public:
static initonly KeyenceDataType^ D
```

```
static val D: KeyenceDataType
```

#### 字段值

类型：[KeyenceDataType](8f6cf1c4-f1e8-10dc-3f3c-127cbbdd8458.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[KeyenceDataType 类](8f6cf1c4-f1e8-10dc-3f3c-127cbbdd8458.htm)

[HslCommunication.Profinet.Keyence 命名空间](baaebe57-0dd8-6393-ef35-013cbae668ce.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## L 字段

[原文連結](http://api.hslcommunication.cn/html/7083f821-41f7-5b22-061b-98d732b32fbf.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Keyence](../html/baaebe57-0dd8-6393-ef35-013cbae668ce.htm "HslCommunication.Profinet.Keyence")

[KeyenceDataType 类](../html/8f6cf1c4-f1e8-10dc-3f3c-127cbbdd8458.htm "KeyenceDataType 类")

[KeyenceDataType 字段](../html/f015109d-1e12-3cc5-e7ca-0a6ca880f9c3.htm "KeyenceDataType 字段")

[B 字段](../html/742136fc-a835-2b05-cb52-8487888b30ad.htm "B 字段")

[CN 字段](../html/cc4e540f-e45a-8809-0aca-56c7cd565636.htm "CN 字段")

[CS 字段](../html/6a9b8328-dc8d-3a7e-c53c-25918c990d16.htm "CS 字段")

[D 字段](../html/0c4edcf3-95ac-4648-2215-97819edfde79.htm "D 字段")

[L 字段](../html/7083f821-41f7-5b22-061b-98d732b32fbf.htm "L 字段")

[M 字段](../html/fe5f1aaa-506a-dffa-9951-46b62d61e008.htm "M 字段")

[R 字段](../html/4d28ca66-e1f7-4d28-2afa-b532d1720b06.htm "R 字段")

[SD 字段](../html/fe4e7372-7f6d-8c31-7cb9-2a96d589f25b.htm "SD 字段")

[SM 字段](../html/7f39bf89-55d4-ca3a-12eb-d6abce4ebe02.htm "SM 字段")

[TN 字段](../html/7af1a704-0c31-70f4-7126-6eda3675137e.htm "TN 字段")

[TS 字段](../html/0bb3d647-793c-84e9-3a36-9074798194f9.htm "TS 字段")

[W 字段](../html/40cd4bfc-d820-d6d1-fcd7-783e83dd585b.htm "W 字段")

[X 字段](../html/d61985b8-6990-aa2e-6abc-427459a7725f.htm "X 字段")

[Y 字段](../html/335aa99c-cdd4-82a1-bd60-00ffbe9c254e.htm "Y 字段")

[ZR 字段](../html/4905c593-832b-f615-89cd-e3431fb4a66d.htm "ZR 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| KeyenceDataTypeL 字段 |

锁存继电器

**命名空间：**
 [HslCommunication.Profinet.Keyence](baaebe57-0dd8-6393-ef35-013cbae668ce.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static readonly KeyenceDataType L
```

```
Public Shared ReadOnly L As KeyenceDataType
```

```
public:
static initonly KeyenceDataType^ L
```

```
static val L: KeyenceDataType
```

#### 字段值

类型：[KeyenceDataType](8f6cf1c4-f1e8-10dc-3f3c-127cbbdd8458.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[KeyenceDataType 类](8f6cf1c4-f1e8-10dc-3f3c-127cbbdd8458.htm)

[HslCommunication.Profinet.Keyence 命名空间](baaebe57-0dd8-6393-ef35-013cbae668ce.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## M 字段

[原文連結](http://api.hslcommunication.cn/html/fe5f1aaa-506a-dffa-9951-46b62d61e008.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Keyence](../html/baaebe57-0dd8-6393-ef35-013cbae668ce.htm "HslCommunication.Profinet.Keyence")

[KeyenceDataType 类](../html/8f6cf1c4-f1e8-10dc-3f3c-127cbbdd8458.htm "KeyenceDataType 类")

[KeyenceDataType 字段](../html/f015109d-1e12-3cc5-e7ca-0a6ca880f9c3.htm "KeyenceDataType 字段")

[B 字段](../html/742136fc-a835-2b05-cb52-8487888b30ad.htm "B 字段")

[CN 字段](../html/cc4e540f-e45a-8809-0aca-56c7cd565636.htm "CN 字段")

[CS 字段](../html/6a9b8328-dc8d-3a7e-c53c-25918c990d16.htm "CS 字段")

[D 字段](../html/0c4edcf3-95ac-4648-2215-97819edfde79.htm "D 字段")

[L 字段](../html/7083f821-41f7-5b22-061b-98d732b32fbf.htm "L 字段")

[M 字段](../html/fe5f1aaa-506a-dffa-9951-46b62d61e008.htm "M 字段")

[R 字段](../html/4d28ca66-e1f7-4d28-2afa-b532d1720b06.htm "R 字段")

[SD 字段](../html/fe4e7372-7f6d-8c31-7cb9-2a96d589f25b.htm "SD 字段")

[SM 字段](../html/7f39bf89-55d4-ca3a-12eb-d6abce4ebe02.htm "SM 字段")

[TN 字段](../html/7af1a704-0c31-70f4-7126-6eda3675137e.htm "TN 字段")

[TS 字段](../html/0bb3d647-793c-84e9-3a36-9074798194f9.htm "TS 字段")

[W 字段](../html/40cd4bfc-d820-d6d1-fcd7-783e83dd585b.htm "W 字段")

[X 字段](../html/d61985b8-6990-aa2e-6abc-427459a7725f.htm "X 字段")

[Y 字段](../html/335aa99c-cdd4-82a1-bd60-00ffbe9c254e.htm "Y 字段")

[ZR 字段](../html/4905c593-832b-f615-89cd-e3431fb4a66d.htm "ZR 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| KeyenceDataTypeM 字段 |

内部辅助继电器

**命名空间：**
 [HslCommunication.Profinet.Keyence](baaebe57-0dd8-6393-ef35-013cbae668ce.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static readonly KeyenceDataType M
```

```
Public Shared ReadOnly M As KeyenceDataType
```

```
public:
static initonly KeyenceDataType^ M
```

```
static val M: KeyenceDataType
```

#### 字段值

类型：[KeyenceDataType](8f6cf1c4-f1e8-10dc-3f3c-127cbbdd8458.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[KeyenceDataType 类](8f6cf1c4-f1e8-10dc-3f3c-127cbbdd8458.htm)

[HslCommunication.Profinet.Keyence 命名空间](baaebe57-0dd8-6393-ef35-013cbae668ce.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## R 字段

[原文連結](http://api.hslcommunication.cn/html/4d28ca66-e1f7-4d28-2afa-b532d1720b06.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Keyence](../html/baaebe57-0dd8-6393-ef35-013cbae668ce.htm "HslCommunication.Profinet.Keyence")

[KeyenceDataType 类](../html/8f6cf1c4-f1e8-10dc-3f3c-127cbbdd8458.htm "KeyenceDataType 类")

[KeyenceDataType 字段](../html/f015109d-1e12-3cc5-e7ca-0a6ca880f9c3.htm "KeyenceDataType 字段")

[B 字段](../html/742136fc-a835-2b05-cb52-8487888b30ad.htm "B 字段")

[CN 字段](../html/cc4e540f-e45a-8809-0aca-56c7cd565636.htm "CN 字段")

[CS 字段](../html/6a9b8328-dc8d-3a7e-c53c-25918c990d16.htm "CS 字段")

[D 字段](../html/0c4edcf3-95ac-4648-2215-97819edfde79.htm "D 字段")

[L 字段](../html/7083f821-41f7-5b22-061b-98d732b32fbf.htm "L 字段")

[M 字段](../html/fe5f1aaa-506a-dffa-9951-46b62d61e008.htm "M 字段")

[R 字段](../html/4d28ca66-e1f7-4d28-2afa-b532d1720b06.htm "R 字段")

[SD 字段](../html/fe4e7372-7f6d-8c31-7cb9-2a96d589f25b.htm "SD 字段")

[SM 字段](../html/7f39bf89-55d4-ca3a-12eb-d6abce4ebe02.htm "SM 字段")

[TN 字段](../html/7af1a704-0c31-70f4-7126-6eda3675137e.htm "TN 字段")

[TS 字段](../html/0bb3d647-793c-84e9-3a36-9074798194f9.htm "TS 字段")

[W 字段](../html/40cd4bfc-d820-d6d1-fcd7-783e83dd585b.htm "W 字段")

[X 字段](../html/d61985b8-6990-aa2e-6abc-427459a7725f.htm "X 字段")

[Y 字段](../html/335aa99c-cdd4-82a1-bd60-00ffbe9c254e.htm "Y 字段")

[ZR 字段](../html/4905c593-832b-f615-89cd-e3431fb4a66d.htm "ZR 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| KeyenceDataTypeR 字段 |

文件寄存器

**命名空间：**
 [HslCommunication.Profinet.Keyence](baaebe57-0dd8-6393-ef35-013cbae668ce.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static readonly KeyenceDataType R
```

```
Public Shared ReadOnly R As KeyenceDataType
```

```
public:
static initonly KeyenceDataType^ R
```

```
static val R: KeyenceDataType
```

#### 字段值

类型：[KeyenceDataType](8f6cf1c4-f1e8-10dc-3f3c-127cbbdd8458.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[KeyenceDataType 类](8f6cf1c4-f1e8-10dc-3f3c-127cbbdd8458.htm)

[HslCommunication.Profinet.Keyence 命名空间](baaebe57-0dd8-6393-ef35-013cbae668ce.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## SD 字段

[原文連結](http://api.hslcommunication.cn/html/fe4e7372-7f6d-8c31-7cb9-2a96d589f25b.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Keyence](../html/baaebe57-0dd8-6393-ef35-013cbae668ce.htm "HslCommunication.Profinet.Keyence")

[KeyenceDataType 类](../html/8f6cf1c4-f1e8-10dc-3f3c-127cbbdd8458.htm "KeyenceDataType 类")

[KeyenceDataType 字段](../html/f015109d-1e12-3cc5-e7ca-0a6ca880f9c3.htm "KeyenceDataType 字段")

[B 字段](../html/742136fc-a835-2b05-cb52-8487888b30ad.htm "B 字段")

[CN 字段](../html/cc4e540f-e45a-8809-0aca-56c7cd565636.htm "CN 字段")

[CS 字段](../html/6a9b8328-dc8d-3a7e-c53c-25918c990d16.htm "CS 字段")

[D 字段](../html/0c4edcf3-95ac-4648-2215-97819edfde79.htm "D 字段")

[L 字段](../html/7083f821-41f7-5b22-061b-98d732b32fbf.htm "L 字段")

[M 字段](../html/fe5f1aaa-506a-dffa-9951-46b62d61e008.htm "M 字段")

[R 字段](../html/4d28ca66-e1f7-4d28-2afa-b532d1720b06.htm "R 字段")

[SD 字段](../html/fe4e7372-7f6d-8c31-7cb9-2a96d589f25b.htm "SD 字段")

[SM 字段](../html/7f39bf89-55d4-ca3a-12eb-d6abce4ebe02.htm "SM 字段")

[TN 字段](../html/7af1a704-0c31-70f4-7126-6eda3675137e.htm "TN 字段")

[TS 字段](../html/0bb3d647-793c-84e9-3a36-9074798194f9.htm "TS 字段")

[W 字段](../html/40cd4bfc-d820-d6d1-fcd7-783e83dd585b.htm "W 字段")

[X 字段](../html/d61985b8-6990-aa2e-6abc-427459a7725f.htm "X 字段")

[Y 字段](../html/335aa99c-cdd4-82a1-bd60-00ffbe9c254e.htm "Y 字段")

[ZR 字段](../html/4905c593-832b-f615-89cd-e3431fb4a66d.htm "ZR 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| KeyenceDataTypeSD 字段 |

控制存储器

**命名空间：**
 [HslCommunication.Profinet.Keyence](baaebe57-0dd8-6393-ef35-013cbae668ce.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static readonly KeyenceDataType SD
```

```
Public Shared ReadOnly SD As KeyenceDataType
```

```
public:
static initonly KeyenceDataType^ SD
```

```
static val SD: KeyenceDataType
```

#### 字段值

类型：[KeyenceDataType](8f6cf1c4-f1e8-10dc-3f3c-127cbbdd8458.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[KeyenceDataType 类](8f6cf1c4-f1e8-10dc-3f3c-127cbbdd8458.htm)

[HslCommunication.Profinet.Keyence 命名空间](baaebe57-0dd8-6393-ef35-013cbae668ce.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## SM 字段

[原文連結](http://api.hslcommunication.cn/html/7f39bf89-55d4-ca3a-12eb-d6abce4ebe02.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Keyence](../html/baaebe57-0dd8-6393-ef35-013cbae668ce.htm "HslCommunication.Profinet.Keyence")

[KeyenceDataType 类](../html/8f6cf1c4-f1e8-10dc-3f3c-127cbbdd8458.htm "KeyenceDataType 类")

[KeyenceDataType 字段](../html/f015109d-1e12-3cc5-e7ca-0a6ca880f9c3.htm "KeyenceDataType 字段")

[B 字段](../html/742136fc-a835-2b05-cb52-8487888b30ad.htm "B 字段")

[CN 字段](../html/cc4e540f-e45a-8809-0aca-56c7cd565636.htm "CN 字段")

[CS 字段](../html/6a9b8328-dc8d-3a7e-c53c-25918c990d16.htm "CS 字段")

[D 字段](../html/0c4edcf3-95ac-4648-2215-97819edfde79.htm "D 字段")

[L 字段](../html/7083f821-41f7-5b22-061b-98d732b32fbf.htm "L 字段")

[M 字段](../html/fe5f1aaa-506a-dffa-9951-46b62d61e008.htm "M 字段")

[R 字段](../html/4d28ca66-e1f7-4d28-2afa-b532d1720b06.htm "R 字段")

[SD 字段](../html/fe4e7372-7f6d-8c31-7cb9-2a96d589f25b.htm "SD 字段")

[SM 字段](../html/7f39bf89-55d4-ca3a-12eb-d6abce4ebe02.htm "SM 字段")

[TN 字段](../html/7af1a704-0c31-70f4-7126-6eda3675137e.htm "TN 字段")

[TS 字段](../html/0bb3d647-793c-84e9-3a36-9074798194f9.htm "TS 字段")

[W 字段](../html/40cd4bfc-d820-d6d1-fcd7-783e83dd585b.htm "W 字段")

[X 字段](../html/d61985b8-6990-aa2e-6abc-427459a7725f.htm "X 字段")

[Y 字段](../html/335aa99c-cdd4-82a1-bd60-00ffbe9c254e.htm "Y 字段")

[ZR 字段](../html/4905c593-832b-f615-89cd-e3431fb4a66d.htm "ZR 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| KeyenceDataTypeSM 字段 |

控制继电器

**命名空间：**
 [HslCommunication.Profinet.Keyence](baaebe57-0dd8-6393-ef35-013cbae668ce.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static readonly KeyenceDataType SM
```

```
Public Shared ReadOnly SM As KeyenceDataType
```

```
public:
static initonly KeyenceDataType^ SM
```

```
static val SM: KeyenceDataType
```

#### 字段值

类型：[KeyenceDataType](8f6cf1c4-f1e8-10dc-3f3c-127cbbdd8458.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[KeyenceDataType 类](8f6cf1c4-f1e8-10dc-3f3c-127cbbdd8458.htm)

[HslCommunication.Profinet.Keyence 命名空间](baaebe57-0dd8-6393-ef35-013cbae668ce.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## TN 字段

[原文連結](http://api.hslcommunication.cn/html/7af1a704-0c31-70f4-7126-6eda3675137e.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Keyence](../html/baaebe57-0dd8-6393-ef35-013cbae668ce.htm "HslCommunication.Profinet.Keyence")

[KeyenceDataType 类](../html/8f6cf1c4-f1e8-10dc-3f3c-127cbbdd8458.htm "KeyenceDataType 类")

[KeyenceDataType 字段](../html/f015109d-1e12-3cc5-e7ca-0a6ca880f9c3.htm "KeyenceDataType 字段")

[B 字段](../html/742136fc-a835-2b05-cb52-8487888b30ad.htm "B 字段")

[CN 字段](../html/cc4e540f-e45a-8809-0aca-56c7cd565636.htm "CN 字段")

[CS 字段](../html/6a9b8328-dc8d-3a7e-c53c-25918c990d16.htm "CS 字段")

[D 字段](../html/0c4edcf3-95ac-4648-2215-97819edfde79.htm "D 字段")

[L 字段](../html/7083f821-41f7-5b22-061b-98d732b32fbf.htm "L 字段")

[M 字段](../html/fe5f1aaa-506a-dffa-9951-46b62d61e008.htm "M 字段")

[R 字段](../html/4d28ca66-e1f7-4d28-2afa-b532d1720b06.htm "R 字段")

[SD 字段](../html/fe4e7372-7f6d-8c31-7cb9-2a96d589f25b.htm "SD 字段")

[SM 字段](../html/7f39bf89-55d4-ca3a-12eb-d6abce4ebe02.htm "SM 字段")

[TN 字段](../html/7af1a704-0c31-70f4-7126-6eda3675137e.htm "TN 字段")

[TS 字段](../html/0bb3d647-793c-84e9-3a36-9074798194f9.htm "TS 字段")

[W 字段](../html/40cd4bfc-d820-d6d1-fcd7-783e83dd585b.htm "W 字段")

[X 字段](../html/d61985b8-6990-aa2e-6abc-427459a7725f.htm "X 字段")

[Y 字段](../html/335aa99c-cdd4-82a1-bd60-00ffbe9c254e.htm "Y 字段")

[ZR 字段](../html/4905c593-832b-f615-89cd-e3431fb4a66d.htm "ZR 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| KeyenceDataTypeTN 字段 |

计时器（当前值）

**命名空间：**
 [HslCommunication.Profinet.Keyence](baaebe57-0dd8-6393-ef35-013cbae668ce.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static readonly KeyenceDataType TN
```

```
Public Shared ReadOnly TN As KeyenceDataType
```

```
public:
static initonly KeyenceDataType^ TN
```

```
static val TN: KeyenceDataType
```

#### 字段值

类型：[KeyenceDataType](8f6cf1c4-f1e8-10dc-3f3c-127cbbdd8458.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[KeyenceDataType 类](8f6cf1c4-f1e8-10dc-3f3c-127cbbdd8458.htm)

[HslCommunication.Profinet.Keyence 命名空间](baaebe57-0dd8-6393-ef35-013cbae668ce.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## TS 字段

[原文連結](http://api.hslcommunication.cn/html/0bb3d647-793c-84e9-3a36-9074798194f9.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Keyence](../html/baaebe57-0dd8-6393-ef35-013cbae668ce.htm "HslCommunication.Profinet.Keyence")

[KeyenceDataType 类](../html/8f6cf1c4-f1e8-10dc-3f3c-127cbbdd8458.htm "KeyenceDataType 类")

[KeyenceDataType 字段](../html/f015109d-1e12-3cc5-e7ca-0a6ca880f9c3.htm "KeyenceDataType 字段")

[B 字段](../html/742136fc-a835-2b05-cb52-8487888b30ad.htm "B 字段")

[CN 字段](../html/cc4e540f-e45a-8809-0aca-56c7cd565636.htm "CN 字段")

[CS 字段](../html/6a9b8328-dc8d-3a7e-c53c-25918c990d16.htm "CS 字段")

[D 字段](../html/0c4edcf3-95ac-4648-2215-97819edfde79.htm "D 字段")

[L 字段](../html/7083f821-41f7-5b22-061b-98d732b32fbf.htm "L 字段")

[M 字段](../html/fe5f1aaa-506a-dffa-9951-46b62d61e008.htm "M 字段")

[R 字段](../html/4d28ca66-e1f7-4d28-2afa-b532d1720b06.htm "R 字段")

[SD 字段](../html/fe4e7372-7f6d-8c31-7cb9-2a96d589f25b.htm "SD 字段")

[SM 字段](../html/7f39bf89-55d4-ca3a-12eb-d6abce4ebe02.htm "SM 字段")

[TN 字段](../html/7af1a704-0c31-70f4-7126-6eda3675137e.htm "TN 字段")

[TS 字段](../html/0bb3d647-793c-84e9-3a36-9074798194f9.htm "TS 字段")

[W 字段](../html/40cd4bfc-d820-d6d1-fcd7-783e83dd585b.htm "W 字段")

[X 字段](../html/d61985b8-6990-aa2e-6abc-427459a7725f.htm "X 字段")

[Y 字段](../html/335aa99c-cdd4-82a1-bd60-00ffbe9c254e.htm "Y 字段")

[ZR 字段](../html/4905c593-832b-f615-89cd-e3431fb4a66d.htm "ZR 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| KeyenceDataTypeTS 字段 |

计时器（接点）

**命名空间：**
 [HslCommunication.Profinet.Keyence](baaebe57-0dd8-6393-ef35-013cbae668ce.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static readonly KeyenceDataType TS
```

```
Public Shared ReadOnly TS As KeyenceDataType
```

```
public:
static initonly KeyenceDataType^ TS
```

```
static val TS: KeyenceDataType
```

#### 字段值

类型：[KeyenceDataType](8f6cf1c4-f1e8-10dc-3f3c-127cbbdd8458.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[KeyenceDataType 类](8f6cf1c4-f1e8-10dc-3f3c-127cbbdd8458.htm)

[HslCommunication.Profinet.Keyence 命名空间](baaebe57-0dd8-6393-ef35-013cbae668ce.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## W 字段

[原文連結](http://api.hslcommunication.cn/html/40cd4bfc-d820-d6d1-fcd7-783e83dd585b.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Keyence](../html/baaebe57-0dd8-6393-ef35-013cbae668ce.htm "HslCommunication.Profinet.Keyence")

[KeyenceDataType 类](../html/8f6cf1c4-f1e8-10dc-3f3c-127cbbdd8458.htm "KeyenceDataType 类")

[KeyenceDataType 字段](../html/f015109d-1e12-3cc5-e7ca-0a6ca880f9c3.htm "KeyenceDataType 字段")

[B 字段](../html/742136fc-a835-2b05-cb52-8487888b30ad.htm "B 字段")

[CN 字段](../html/cc4e540f-e45a-8809-0aca-56c7cd565636.htm "CN 字段")

[CS 字段](../html/6a9b8328-dc8d-3a7e-c53c-25918c990d16.htm "CS 字段")

[D 字段](../html/0c4edcf3-95ac-4648-2215-97819edfde79.htm "D 字段")

[L 字段](../html/7083f821-41f7-5b22-061b-98d732b32fbf.htm "L 字段")

[M 字段](../html/fe5f1aaa-506a-dffa-9951-46b62d61e008.htm "M 字段")

[R 字段](../html/4d28ca66-e1f7-4d28-2afa-b532d1720b06.htm "R 字段")

[SD 字段](../html/fe4e7372-7f6d-8c31-7cb9-2a96d589f25b.htm "SD 字段")

[SM 字段](../html/7f39bf89-55d4-ca3a-12eb-d6abce4ebe02.htm "SM 字段")

[TN 字段](../html/7af1a704-0c31-70f4-7126-6eda3675137e.htm "TN 字段")

[TS 字段](../html/0bb3d647-793c-84e9-3a36-9074798194f9.htm "TS 字段")

[W 字段](../html/40cd4bfc-d820-d6d1-fcd7-783e83dd585b.htm "W 字段")

[X 字段](../html/d61985b8-6990-aa2e-6abc-427459a7725f.htm "X 字段")

[Y 字段](../html/335aa99c-cdd4-82a1-bd60-00ffbe9c254e.htm "Y 字段")

[ZR 字段](../html/4905c593-832b-f615-89cd-e3431fb4a66d.htm "ZR 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| KeyenceDataTypeW 字段 |

链路寄存器

**命名空间：**
 [HslCommunication.Profinet.Keyence](baaebe57-0dd8-6393-ef35-013cbae668ce.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static readonly KeyenceDataType W
```

```
Public Shared ReadOnly W As KeyenceDataType
```

```
public:
static initonly KeyenceDataType^ W
```

```
static val W: KeyenceDataType
```

#### 字段值

类型：[KeyenceDataType](8f6cf1c4-f1e8-10dc-3f3c-127cbbdd8458.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[KeyenceDataType 类](8f6cf1c4-f1e8-10dc-3f3c-127cbbdd8458.htm)

[HslCommunication.Profinet.Keyence 命名空间](baaebe57-0dd8-6393-ef35-013cbae668ce.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## X 字段

[原文連結](http://api.hslcommunication.cn/html/d61985b8-6990-aa2e-6abc-427459a7725f.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Keyence](../html/baaebe57-0dd8-6393-ef35-013cbae668ce.htm "HslCommunication.Profinet.Keyence")

[KeyenceDataType 类](../html/8f6cf1c4-f1e8-10dc-3f3c-127cbbdd8458.htm "KeyenceDataType 类")

[KeyenceDataType 字段](../html/f015109d-1e12-3cc5-e7ca-0a6ca880f9c3.htm "KeyenceDataType 字段")

[B 字段](../html/742136fc-a835-2b05-cb52-8487888b30ad.htm "B 字段")

[CN 字段](../html/cc4e540f-e45a-8809-0aca-56c7cd565636.htm "CN 字段")

[CS 字段](../html/6a9b8328-dc8d-3a7e-c53c-25918c990d16.htm "CS 字段")

[D 字段](../html/0c4edcf3-95ac-4648-2215-97819edfde79.htm "D 字段")

[L 字段](../html/7083f821-41f7-5b22-061b-98d732b32fbf.htm "L 字段")

[M 字段](../html/fe5f1aaa-506a-dffa-9951-46b62d61e008.htm "M 字段")

[R 字段](../html/4d28ca66-e1f7-4d28-2afa-b532d1720b06.htm "R 字段")

[SD 字段](../html/fe4e7372-7f6d-8c31-7cb9-2a96d589f25b.htm "SD 字段")

[SM 字段](../html/7f39bf89-55d4-ca3a-12eb-d6abce4ebe02.htm "SM 字段")

[TN 字段](../html/7af1a704-0c31-70f4-7126-6eda3675137e.htm "TN 字段")

[TS 字段](../html/0bb3d647-793c-84e9-3a36-9074798194f9.htm "TS 字段")

[W 字段](../html/40cd4bfc-d820-d6d1-fcd7-783e83dd585b.htm "W 字段")

[X 字段](../html/d61985b8-6990-aa2e-6abc-427459a7725f.htm "X 字段")

[Y 字段](../html/335aa99c-cdd4-82a1-bd60-00ffbe9c254e.htm "Y 字段")

[ZR 字段](../html/4905c593-832b-f615-89cd-e3431fb4a66d.htm "ZR 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| KeyenceDataTypeX 字段 |

X输入继电器

**命名空间：**
 [HslCommunication.Profinet.Keyence](baaebe57-0dd8-6393-ef35-013cbae668ce.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static readonly KeyenceDataType X
```

```
Public Shared ReadOnly X As KeyenceDataType
```

```
public:
static initonly KeyenceDataType^ X
```

```
static val X: KeyenceDataType
```

#### 字段值

类型：[KeyenceDataType](8f6cf1c4-f1e8-10dc-3f3c-127cbbdd8458.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[KeyenceDataType 类](8f6cf1c4-f1e8-10dc-3f3c-127cbbdd8458.htm)

[HslCommunication.Profinet.Keyence 命名空间](baaebe57-0dd8-6393-ef35-013cbae668ce.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Y 字段

[原文連結](http://api.hslcommunication.cn/html/335aa99c-cdd4-82a1-bd60-00ffbe9c254e.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Keyence](../html/baaebe57-0dd8-6393-ef35-013cbae668ce.htm "HslCommunication.Profinet.Keyence")

[KeyenceDataType 类](../html/8f6cf1c4-f1e8-10dc-3f3c-127cbbdd8458.htm "KeyenceDataType 类")

[KeyenceDataType 字段](../html/f015109d-1e12-3cc5-e7ca-0a6ca880f9c3.htm "KeyenceDataType 字段")

[B 字段](../html/742136fc-a835-2b05-cb52-8487888b30ad.htm "B 字段")

[CN 字段](../html/cc4e540f-e45a-8809-0aca-56c7cd565636.htm "CN 字段")

[CS 字段](../html/6a9b8328-dc8d-3a7e-c53c-25918c990d16.htm "CS 字段")

[D 字段](../html/0c4edcf3-95ac-4648-2215-97819edfde79.htm "D 字段")

[L 字段](../html/7083f821-41f7-5b22-061b-98d732b32fbf.htm "L 字段")

[M 字段](../html/fe5f1aaa-506a-dffa-9951-46b62d61e008.htm "M 字段")

[R 字段](../html/4d28ca66-e1f7-4d28-2afa-b532d1720b06.htm "R 字段")

[SD 字段](../html/fe4e7372-7f6d-8c31-7cb9-2a96d589f25b.htm "SD 字段")

[SM 字段](../html/7f39bf89-55d4-ca3a-12eb-d6abce4ebe02.htm "SM 字段")

[TN 字段](../html/7af1a704-0c31-70f4-7126-6eda3675137e.htm "TN 字段")

[TS 字段](../html/0bb3d647-793c-84e9-3a36-9074798194f9.htm "TS 字段")

[W 字段](../html/40cd4bfc-d820-d6d1-fcd7-783e83dd585b.htm "W 字段")

[X 字段](../html/d61985b8-6990-aa2e-6abc-427459a7725f.htm "X 字段")

[Y 字段](../html/335aa99c-cdd4-82a1-bd60-00ffbe9c254e.htm "Y 字段")

[ZR 字段](../html/4905c593-832b-f615-89cd-e3431fb4a66d.htm "ZR 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| KeyenceDataTypeY 字段 |

Y输出继电器

**命名空间：**
 [HslCommunication.Profinet.Keyence](baaebe57-0dd8-6393-ef35-013cbae668ce.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static readonly KeyenceDataType Y
```

```
Public Shared ReadOnly Y As KeyenceDataType
```

```
public:
static initonly KeyenceDataType^ Y
```

```
static val Y: KeyenceDataType
```

#### 字段值

类型：[KeyenceDataType](8f6cf1c4-f1e8-10dc-3f3c-127cbbdd8458.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[KeyenceDataType 类](8f6cf1c4-f1e8-10dc-3f3c-127cbbdd8458.htm)

[HslCommunication.Profinet.Keyence 命名空间](baaebe57-0dd8-6393-ef35-013cbae668ce.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ZR 字段

[原文連結](http://api.hslcommunication.cn/html/4905c593-832b-f615-89cd-e3431fb4a66d.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Keyence](../html/baaebe57-0dd8-6393-ef35-013cbae668ce.htm "HslCommunication.Profinet.Keyence")

[KeyenceDataType 类](../html/8f6cf1c4-f1e8-10dc-3f3c-127cbbdd8458.htm "KeyenceDataType 类")

[KeyenceDataType 字段](../html/f015109d-1e12-3cc5-e7ca-0a6ca880f9c3.htm "KeyenceDataType 字段")

[B 字段](../html/742136fc-a835-2b05-cb52-8487888b30ad.htm "B 字段")

[CN 字段](../html/cc4e540f-e45a-8809-0aca-56c7cd565636.htm "CN 字段")

[CS 字段](../html/6a9b8328-dc8d-3a7e-c53c-25918c990d16.htm "CS 字段")

[D 字段](../html/0c4edcf3-95ac-4648-2215-97819edfde79.htm "D 字段")

[L 字段](../html/7083f821-41f7-5b22-061b-98d732b32fbf.htm "L 字段")

[M 字段](../html/fe5f1aaa-506a-dffa-9951-46b62d61e008.htm "M 字段")

[R 字段](../html/4d28ca66-e1f7-4d28-2afa-b532d1720b06.htm "R 字段")

[SD 字段](../html/fe4e7372-7f6d-8c31-7cb9-2a96d589f25b.htm "SD 字段")

[SM 字段](../html/7f39bf89-55d4-ca3a-12eb-d6abce4ebe02.htm "SM 字段")

[TN 字段](../html/7af1a704-0c31-70f4-7126-6eda3675137e.htm "TN 字段")

[TS 字段](../html/0bb3d647-793c-84e9-3a36-9074798194f9.htm "TS 字段")

[W 字段](../html/40cd4bfc-d820-d6d1-fcd7-783e83dd585b.htm "W 字段")

[X 字段](../html/d61985b8-6990-aa2e-6abc-427459a7725f.htm "X 字段")

[Y 字段](../html/335aa99c-cdd4-82a1-bd60-00ffbe9c254e.htm "Y 字段")

[ZR 字段](../html/4905c593-832b-f615-89cd-e3431fb4a66d.htm "ZR 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| KeyenceDataTypeZR 字段 |

文件寄存器

**命名空间：**
 [HslCommunication.Profinet.Keyence](baaebe57-0dd8-6393-ef35-013cbae668ce.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static readonly KeyenceDataType ZR
```

```
Public Shared ReadOnly ZR As KeyenceDataType
```

```
public:
static initonly KeyenceDataType^ ZR
```

```
static val ZR: KeyenceDataType
```

#### 字段值

类型：[KeyenceDataType](8f6cf1c4-f1e8-10dc-3f3c-127cbbdd8458.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[KeyenceDataType 类](8f6cf1c4-f1e8-10dc-3f3c-127cbbdd8458.htm)

[HslCommunication.Profinet.Keyence 命名空间](baaebe57-0dd8-6393-ef35-013cbae668ce.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## KeyenceDLEN1 类

[原文連結](http://api.hslcommunication.cn/html/f70c1a08-5573-e24a-5e62-6160f8d51d85.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Keyence](../html/baaebe57-0dd8-6393-ef35-013cbae668ce.htm "HslCommunication.Profinet.Keyence")

[KeyenceDLEN1 类](../html/f70c1a08-5573-e24a-5e62-6160f8d51d85.htm "KeyenceDLEN1 类")

[KeyenceDLEN1 构造函数](../html/2b448611-e563-81ea-bb8f-f19f2b770227.htm "KeyenceDLEN1 构造函数 ")

[KeyenceDLEN1 属性](../html/0e1b2426-1733-156b-88f3-7afbfa6f6fb6.htm "KeyenceDLEN1 属性")

[KeyenceDLEN1 方法](../html/37026385-faa9-45e7-5144-a4640a4cf549.htm "KeyenceDLEN1 方法")

[KeyenceDLEN1 字段](../html/f2fd7589-523e-969b-1681-6693b51c9e3b.htm "KeyenceDLEN1 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| KeyenceDLEN1 类 |

基恩士的数字传感器的以太网模块，可以同时连接并读取多个传感器模块的功能代码

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  [HslCommunication.Core.NetNetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)  
    [HslCommunication.Core.NetNetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)  
      HslCommunication.Profinet.KeyenceKeyenceDLEN1

**命名空间：**
 [HslCommunication.Profinet.Keyence](baaebe57-0dd8-6393-ef35-013cbae668ce.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class KeyenceDLEN1 : NetworkDoubleBase
```

```
Public Class KeyenceDLEN1
	Inherits NetworkDoubleBase
```

```
public ref class KeyenceDLEN1 : public NetworkDoubleBase
```

```
type KeyenceDLEN1 =  
    class
        inherit NetworkDoubleBase
    end
```

KeyenceDLEN1 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [KeyenceDLEN1](b89d07a3-5efd-6043-9a49-97e1a38f3f0e.htm) | 实例化基恩士的Qna兼容3E帧协议的通讯对象  Instantiate Keyence Qna compatible 3E frame protocol communication object |
| 公共方法 | [KeyenceDLEN1(String, Int32)](91afc05b-9a98-6583-d40b-90e55a351607.htm) | 指定ip地址及端口号来实例化一个基恩士的Qna兼容3E帧协议的通讯对象  Specify an IP address and port number to instantiate a Keynes Qna compatible 3E frame protocol communication object |

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
| 公共属性代码示例 | [Port](45e6a4b4-1190-6d1a-e1e6-d8c010e15fe4.htm) | 获取或设置服务器的端口号，具体的值需要取决于对方的配置  Gets or sets the port number of the server. The specific value depends on the configuration of the other party. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共属性代码示例 | [ReceiveTimeOut](995db577-ab5c-6c9c-6f6f-253e26600472.htm) | 获取或设置接收服务器反馈的时间，如果为负数，则不接收反馈   Gets or sets the time to receive server feedback, and if it is a negative number, does not receive feedback (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共属性 | [SendBeforeHex](59a68d6e-a784-a757-40cf-138d6b05099b.htm) | 获取或设置在发送通信报文前追加发送的字节信息，HEX格式，通常用于lora组网时，需要携带 00 00 00 02 四个字节的站地址功能。  Obtain or set the byte information sent before sending communication packets, HEX format, usually used for LORA networking, you need to carry 00 00 00 02 four-byte station address function. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共属性 | [SleepTime](d275ea69-962f-70f9-b486-cb3dc3382b74.htm) | 获取或设置在正式接收对方返回数据前的时候，需要休息的时间，当设置为0的时候，不需要休息。  Get or set the time required to rest before officially receiving the data from the other party. When it is set to 0, no rest is required. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共属性 | [SocketKeepAliveTime](7f4b5ca3-e4db-76e4-ef2b-4c93bf9447ac.htm) | 获取或设置客户端的Socket的心跳时间信息，这个是Socket底层自动实现的心跳包，不基于协议层实现。默认小于0，不开启心跳检测，如果需要开启，设置 60\_000 比较合适，单位毫秒  Get or set the heartbeat time information of the Socket of the client. This is the heartbeat packet automatically implemented by the bottom layer of the Socket, not based on the protocol layer. The default value is less than 0, and heartbeat detection is not enabled. If you need to enable it, it is more appropriate to set 60\_000, in milliseconds. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
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
| 公共方法静态成员 | [CheckResponse](fecce584-b3e4-48e8-ac9c-d7a0cda5b75f.htm) | 坚持设备的返回的数据，并校验是否成功 |
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
| 受保护的方法 | [DeleteFileByName](59afa4a4-3f3c-9831-26ee-c7c13a213bf7.htm) | 删除一个指定的文件，如果文件不存在，直接返回 True，如果文件存在则直接删除，删除成功返回 True，如果发生了异常，返回False  Delete a specified file, if the file does not exist, return True directly, if the file exists, delete it directly, if the deletion is successful, return True, if an exception occurs, return  False (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 公共方法 | [Dispose](8a86a200-fbe0-3b45-0fed-c551405c9b59.htm) | 释放当前的资源，如果调用了本方法，那么该对象再使用的时候，需要重新实例化。  Release the current resource. If this method is called, the object needs to be instantiated again when it is used again. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 受保护的方法 | [Dispose(Boolean)](32f03275-9ae0-525c-84bc-6e2515d11c12.htm) | 释放当前的资源，并自动关闭长连接，如果设置了的话 (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | [ExtraAfterReadFromCoreServer](9f8357e9-dff4-1002-9368-0744cdedff29.htm) | 和服务器交互完成的时候调用的方法，可以根据读写结果进行一些额外的操作，具体的操作需要根据实际的需求来重写实现  The method called when the interaction with the server is completed can perform some additional operations based on the read and write results. The specific operations need to be rewritten according to actual needs. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 受保护的方法代码示例 | [ExtraOnDisconnect](b8569c75-5e55-cd55-5251-207465a8ba39.htm) | 根据实际的协议选择是否重写本方法，有些协议在断开连接之前，需要发送一些报文来关闭当前的网络通道  Select whether to rewrite this method according to the actual protocol. Some protocols need to send some packets to close the current network channel before disconnecting. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 受保护的方法代码示例 | [ExtraOnDisconnectAsync](86aa00d8-b2d6-ada9-5926-2b7a2237fe6b.htm) | 根据实际的协议选择是否重写本方法，有些协议在断开连接之前，需要发送一些报文来关闭当前的网络通道  Select whether to rewrite this method according to the actual protocol. Some protocols need to send some packets to close the current network channel before disconnecting. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 受保护的方法 | [GetAvailableSocket](e86a6fc1-eb91-8de9-5ff8-799b0ea3786b.htm) | 获取本次操作的可用的网络通道，如果是短连接，就重新生成一个新的网络通道，如果是长连接，就复用当前的网络通道。  Obtain the available network channels for this operation. If it is a short connection, a new network channel is regenerated. If it is a long connection, the current network channel is reused. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 受保护的方法 | [GetAvailableSocketAsync](40cf9ce7-a5e6-5431-2e4b-56e3b7c98346.htm) | 获取本次操作的可用的网络通道，如果是短连接，就重新生成一个新的网络通道，如果是长连接，就复用当前的网络通道。  Obtain the available network channels for this operation. If it is a short connection, a new network channel is regenerated. If it is a long connection, the current network channel is reused. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 受保护的方法 | [GetNewNetMessage](9227a14d-2786-82d3-705d-42f941cbc4f0.htm) | 获取一个新的消息对象的方法，需要在继承类里面进行重写  The method to get a new message object needs to be overridden in the inheritance class (重写 [NetworkDoubleBaseGetNewNetMessage](481a40e0-8539-5fc8-3e2b-c25bcc952697.htm).) |
| 公共方法 | [GetPipeSocket](19904680-372e-e0a6-1b27-b4180e759668.htm) | 获取当前用于通信的管道信息  Get the current pipe information used for communication (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法代码示例 | [InitializationOnConnect](40d5eb26-9bca-11cf-42a1-55af6ea5aebd.htm) | 根据实际的协议选择是否重写本方法，有些协议在创建连接之后，需要进行一些初始化的信号握手，才能最终建立网络通道。  Whether to rewrite this method is based on the actual protocol. Some protocols require some initial signal handshake to establish a network channel after the connection is created. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 受保护的方法代码示例 | [InitializationOnConnectAsync](14cac21b-4001-da7d-09a6-941eba868e14.htm) | 根据实际的协议选择是否重写本方法，有些协议在创建连接之后，需要进行一些初始化的信号握手，才能最终建立网络通道。  Whether to rewrite this method is based on the actual protocol. Some protocols require some initial signal handshake to establish a network channel after the connection is created. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | [IpAddressPing](d9611e18-4033-de95-444e-b91ae0a81709.htm) | 对当前设备的IP地址进行PING的操作，返回PING的结果，正常来说，返回Success  PING the IP address of the current device and return the PING result. Normally, it returns Success (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法 | [PackCommandWithHeader](f96881d6-e59f-9e16-07e3-1188eee1a934.htm) | 对当前的命令进行打包处理，通常是携带命令头内容，标记当前的命令的长度信息，需要进行重写，否则默认不打包  The current command is packaged, usually carrying the content of the command header, marking the length of the current command, and it needs to be rewritten, otherwise it is not packaged by default (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | [ReadByCommand](3f86d309-e309-31be-0b29-983a30ec3d9f.htm) | 使用M0命令读取所有的传感器的数据信息 |
| 公共方法代码示例 | [ReadFromCoreServer(Byte)](b3e14431-d773-0d1a-f453-a3ad60648c29.htm) | 将数据发送到当前的网络通道中，并从网络通道中接收一个[INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm)指定的完整的报文，网络通道将根据[GetAvailableSocket](e86a6fc1-eb91-8de9-5ff8-799b0ea3786b.htm)方法自动获取，本方法是线程安全的。  Send data to the current network channel and receive a complete message specified by [INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm) from the network channel. The network channel will be automatically obtained according to the [GetAvailableSocket](e86a6fc1-eb91-8de9-5ff8-799b0ea3786b.htm) method This method is thread-safe. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | [ReadFromCoreServer(IEnumerableByte)](580dae22-fad5-5816-d58f-43068e89e419.htm) | 将多个数据报文按顺序发到设备，并从设备接收返回的数据内容，然后拼接成一个Byte[]信息，需要重写[UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm)方法才能返回正确的结果。  Send multiple data packets to the device in sequence, and receive the returned data content from the device, and then splicing them into a Byte[] message, you need to rewrite [UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm) method to return the correct result. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServer(Byte, Boolean, Boolean)](1a33d974-805d-f149-3665-eb8d8a99ea84.htm) | 将数据发送到当前的网络通道中，并从网络通道中接收一个[INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm)指定的完整的报文，网络通道将根据[GetAvailableSocket](e86a6fc1-eb91-8de9-5ff8-799b0ea3786b.htm)方法自动获取，本方法是线程安全的。  Send data to the current network channel and receive a complete message specified by [INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm) from the network channel. The network channel will be automatically obtained according to the [GetAvailableSocket](e86a6fc1-eb91-8de9-5ff8-799b0ea3786b.htm) method This method is thread-safe. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServer(Socket, Byte, Boolean, Boolean)](3b72b289-b201-7548-0981-ac049c9a4135.htm) | 将数据报文发送指定的网络通道上，根据当前指定的[INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm)类型，返回一条完整的数据指令  Sends a data message to the specified network channel, and returns a complete data command according to the currently specified [INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm) type (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServerAsync(Byte)](ec51a18e-0f5d-dbe1-c229-d62e8760ac97.htm) | 将数据发送到当前的网络通道中，并从网络通道中接收一个[INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm)指定的完整的报文，网络通道将根据[GetAvailableSocket](e86a6fc1-eb91-8de9-5ff8-799b0ea3786b.htm)方法自动获取，本方法是线程安全的。  Send data to the current network channel and receive a complete message specified by [INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm) from the network channel. The network channel will be automatically obtained according to the [GetAvailableSocket](e86a6fc1-eb91-8de9-5ff8-799b0ea3786b.htm) method This method is thread-safe. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | [ReadFromCoreServerAsync(IEnumerableByte)](0fce97fd-8441-6b54-0877-b8dfacc738e3.htm) | 将多个数据报文按顺序发到设备，并从设备接收返回的数据内容，然后拼接成一个Byte[]信息，需要重写[UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm)方法才能返回正确的结果。  Send multiple data packets to the device in sequence, and receive the returned data content from the device, and then splicing them into a Byte[] message, you need to rewrite [UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm) method to return the correct result. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServerAsync(Byte, Boolean, Boolean)](7525b379-04c2-8dac-b451-a8b5413d6747.htm) | 将数据发送到当前的网络通道中，并从网络通道中接收一个[INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm)指定的完整的报文，网络通道将根据[GetAvailableSocket](e86a6fc1-eb91-8de9-5ff8-799b0ea3786b.htm)方法自动获取，本方法是线程安全的。  Send data to the current network channel and receive a complete message specified by [INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm) from the network channel. The network channel will be automatically obtained according to the [GetAvailableSocket](e86a6fc1-eb91-8de9-5ff8-799b0ea3786b.htm) method This method is thread-safe. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServerAsync(Socket, Byte, Boolean, Boolean)](07db4a00-b548-8bd6-5de6-12deb91f5311.htm) | 将数据报文发送指定的网络通道上，根据当前指定的[INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm)类型，返回一条完整的数据指令  Sends a data message to the specified network channel, and returns a complete data command according to the currently specified [INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm) type (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
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
| 公共方法 | [SetLoginAccount](725afc25-6e3f-bd7d-25ed-580270c8fc56.htm) | 设置当前的登录的账户名和密码信息，并启用账户验证的功能，账户名为空时设置不生效  Set the current login account name and password information, and enable the account verification function. The account name setting will not take effect when it is empty (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法代码示例 | [SetPersistentConnection](d8e8df9f-cbe9-6863-54e9-aa3b38906b26.htm) | 在读取数据之前可以调用本方法将客户端设置为长连接模式，相当于跳过了ConnectServer的结果验证，对异形客户端无效，当第一次进行通信时再进行创建连接请求。  Before reading the data, you can call this method to set the client to the long connection mode, which is equivalent to skipping the result verification of ConnectServer, and it is invalid for the alien client. When the first communication is performed, the connection creation request is performed. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | [SetPipeSocket](b6a7a243-2a6d-c4ea-45b7-bf980450a2bd.htm) | 设置一个新的网络管道，一般来说不需要调用本方法，当多个网口设备共用一个网络连接时才需要使用本方法进行设置共享的管道。  To set up a new network channel, generally speaking, you do not need to call this method. This method is only needed to set up a shared channel when multiple network port devices share a network connection. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共方法 | [ToString](d796efa2-af69-082d-2f83-df50d7d65586.htm) | (重写 [NetworkDoubleBaseToString](85a8532f-8fc5-886c-9d76-70844e0a7d64.htm).) |
| 公共方法 | [UnpackResponseContent](2164ab7a-ce41-3ee1-db04-af22ee11b3de.htm) | 根据对方返回的报文命令，对命令进行基本的拆包，例如各种Modbus协议拆包为统一的核心报文，还支持对报文的验证  According to the message command returned by the other party, the command is basically unpacked, for example, various Modbus protocols are unpacked into a unified core message, and the verification of the message is also supported (重写 [NetworkDoubleBaseUnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm).) |
| 受保护的方法 | [WriteStreamFromSocket](97bf95ea-51fe-8f2e-5da6-1c83766aaa78.htm) | 从套接字中接收所有的数据然后写入到指定的流当中去，需要指定数据的长度，支持按照百分比进行进度报告  Receives all data from the socket and writes it to the specified stream. The length of the data needs to be specified, and progress reporting is supported in percentage. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [WriteStreamFromSocketAsync](b4cf7f6d-1e0c-a266-4c49-683fdb14a922.htm) | 从套接字中接收所有的数据然后写入到指定的流当中去，需要指定数据的长度，支持按照百分比进行进度报告  Receives all data from the socket and writes it to the specified stream. The length of the data needs to be specified, and progress reporting is supported in percentage. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |

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

[HslCommunication.Profinet.Keyence 命名空间](baaebe57-0dd8-6393-ef35-013cbae668ce.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## KeyenceDLEN1 构造函数 

[原文連結](http://api.hslcommunication.cn/html/2b448611-e563-81ea-bb8f-f19f2b770227.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Keyence](../html/baaebe57-0dd8-6393-ef35-013cbae668ce.htm "HslCommunication.Profinet.Keyence")

[KeyenceDLEN1 类](../html/f70c1a08-5573-e24a-5e62-6160f8d51d85.htm "KeyenceDLEN1 类")

[KeyenceDLEN1 构造函数](../html/2b448611-e563-81ea-bb8f-f19f2b770227.htm "KeyenceDLEN1 构造函数 ")

[KeyenceDLEN1 构造函数](../html/b89d07a3-5efd-6043-9a49-97e1a38f3f0e.htm "KeyenceDLEN1 构造函数 ")

[KeyenceDLEN1 构造函数 (String, Int32)](../html/91afc05b-9a98-6583-d40b-90e55a351607.htm "KeyenceDLEN1 构造函数 (String, Int32)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| KeyenceDLEN1 构造函数 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [KeyenceDLEN1](b89d07a3-5efd-6043-9a49-97e1a38f3f0e.htm) | 实例化基恩士的Qna兼容3E帧协议的通讯对象  Instantiate Keyence Qna compatible 3E frame protocol communication object |
| 公共方法 | [KeyenceDLEN1(String, Int32)](91afc05b-9a98-6583-d40b-90e55a351607.htm) | 指定ip地址及端口号来实例化一个基恩士的Qna兼容3E帧协议的通讯对象  Specify an IP address and port number to instantiate a Keynes Qna compatible 3E frame protocol communication object |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[KeyenceDLEN1 类](f70c1a08-5573-e24a-5e62-6160f8d51d85.htm)

[HslCommunication.Profinet.Keyence 命名空间](baaebe57-0dd8-6393-ef35-013cbae668ce.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## KeyenceDLEN1 构造函数 

[原文連結](http://api.hslcommunication.cn/html/b89d07a3-5efd-6043-9a49-97e1a38f3f0e.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Keyence](../html/baaebe57-0dd8-6393-ef35-013cbae668ce.htm "HslCommunication.Profinet.Keyence")

[KeyenceDLEN1 类](../html/f70c1a08-5573-e24a-5e62-6160f8d51d85.htm "KeyenceDLEN1 类")

[KeyenceDLEN1 构造函数](../html/2b448611-e563-81ea-bb8f-f19f2b770227.htm "KeyenceDLEN1 构造函数 ")

[KeyenceDLEN1 构造函数](../html/b89d07a3-5efd-6043-9a49-97e1a38f3f0e.htm "KeyenceDLEN1 构造函数 ")

[KeyenceDLEN1 构造函数 (String, Int32)](../html/91afc05b-9a98-6583-d40b-90e55a351607.htm "KeyenceDLEN1 构造函数 (String, Int32)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| KeyenceDLEN1 构造函数 |

实例化基恩士的Qna兼容3E帧协议的通讯对象  
Instantiate Keyence Qna compatible 3E frame protocol communication object

**命名空间：**
 [HslCommunication.Profinet.Keyence](baaebe57-0dd8-6393-ef35-013cbae668ce.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public KeyenceDLEN1()
```

```
Public Sub New
```

```
public:
KeyenceDLEN1()
```

```
new : unit -> KeyenceDLEN1
```

![](../icons/SectionExpanded.png)参见

#### 引用

[KeyenceDLEN1 类](f70c1a08-5573-e24a-5e62-6160f8d51d85.htm)

[KeyenceDLEN1 重载](2b448611-e563-81ea-bb8f-f19f2b770227.htm)

[HslCommunication.Profinet.Keyence 命名空间](baaebe57-0dd8-6393-ef35-013cbae668ce.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## KeyenceDLEN1 构造函数 (String, Int32)

[原文連結](http://api.hslcommunication.cn/html/91afc05b-9a98-6583-d40b-90e55a351607.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Keyence](../html/baaebe57-0dd8-6393-ef35-013cbae668ce.htm "HslCommunication.Profinet.Keyence")

[KeyenceDLEN1 类](../html/f70c1a08-5573-e24a-5e62-6160f8d51d85.htm "KeyenceDLEN1 类")

[KeyenceDLEN1 构造函数](../html/2b448611-e563-81ea-bb8f-f19f2b770227.htm "KeyenceDLEN1 构造函数 ")

[KeyenceDLEN1 构造函数](../html/b89d07a3-5efd-6043-9a49-97e1a38f3f0e.htm "KeyenceDLEN1 构造函数 ")

[KeyenceDLEN1 构造函数 (String, Int32)](../html/91afc05b-9a98-6583-d40b-90e55a351607.htm "KeyenceDLEN1 构造函数 (String, Int32)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| KeyenceDLEN1 构造函数 (String, Int32) |

指定ip地址及端口号来实例化一个基恩士的Qna兼容3E帧协议的通讯对象  
Specify an IP address and port number to instantiate a Keynes Qna compatible 3E frame protocol communication object

**命名空间：**
 [HslCommunication.Profinet.Keyence](baaebe57-0dd8-6393-ef35-013cbae668ce.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public KeyenceDLEN1(
	string ipAddress,
	int port
)
```

```
Public Sub New ( 
	ipAddress As String,
	port As Integer
)
```

```
public:
KeyenceDLEN1(
	String^ ipAddress, 
	int port
)
```

```
new : 
        ipAddress : string * 
        port : int -> KeyenceDLEN1
```

#### 参数

ipAddress
:   类型：SystemString  
    PLC的Ip地址

port
:   类型：SystemInt32  
    PLC的端口

![](../icons/SectionExpanded.png)参见

#### 引用

[KeyenceDLEN1 类](f70c1a08-5573-e24a-5e62-6160f8d51d85.htm)

[KeyenceDLEN1 重载](2b448611-e563-81ea-bb8f-f19f2b770227.htm)

[HslCommunication.Profinet.Keyence 命名空间](baaebe57-0dd8-6393-ef35-013cbae668ce.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## KeyenceDLEN1 属性

[原文連結](http://api.hslcommunication.cn/html/0e1b2426-1733-156b-88f3-7afbfa6f6fb6.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Keyence](../html/baaebe57-0dd8-6393-ef35-013cbae668ce.htm "HslCommunication.Profinet.Keyence")

[KeyenceDLEN1 类](../html/f70c1a08-5573-e24a-5e62-6160f8d51d85.htm "KeyenceDLEN1 类")

[KeyenceDLEN1 构造函数](../html/2b448611-e563-81ea-bb8f-f19f2b770227.htm "KeyenceDLEN1 构造函数 ")

[KeyenceDLEN1 属性](../html/0e1b2426-1733-156b-88f3-7afbfa6f6fb6.htm "KeyenceDLEN1 属性")

[KeyenceDLEN1 方法](../html/37026385-faa9-45e7-5144-a4640a4cf549.htm "KeyenceDLEN1 方法")

[KeyenceDLEN1 字段](../html/f2fd7589-523e-969b-1681-6693b51c9e3b.htm "KeyenceDLEN1 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| KeyenceDLEN1 属性 |

[KeyenceDLEN1](f70c1a08-5573-e24a-5e62-6160f8d51d85.htm) 类型公开以下成员。

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
| 公共属性代码示例 | [Port](45e6a4b4-1190-6d1a-e1e6-d8c010e15fe4.htm) | 获取或设置服务器的端口号，具体的值需要取决于对方的配置  Gets or sets the port number of the server. The specific value depends on the configuration of the other party. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共属性代码示例 | [ReceiveTimeOut](995db577-ab5c-6c9c-6f6f-253e26600472.htm) | 获取或设置接收服务器反馈的时间，如果为负数，则不接收反馈   Gets or sets the time to receive server feedback, and if it is a negative number, does not receive feedback (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共属性 | [SendBeforeHex](59a68d6e-a784-a757-40cf-138d6b05099b.htm) | 获取或设置在发送通信报文前追加发送的字节信息，HEX格式，通常用于lora组网时，需要携带 00 00 00 02 四个字节的站地址功能。  Obtain or set the byte information sent before sending communication packets, HEX format, usually used for LORA networking, you need to carry 00 00 00 02 four-byte station address function. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共属性 | [SleepTime](d275ea69-962f-70f9-b486-cb3dc3382b74.htm) | 获取或设置在正式接收对方返回数据前的时候，需要休息的时间，当设置为0的时候，不需要休息。  Get or set the time required to rest before officially receiving the data from the other party. When it is set to 0, no rest is required. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共属性 | [SocketKeepAliveTime](7f4b5ca3-e4db-76e4-ef2b-4c93bf9447ac.htm) | 获取或设置客户端的Socket的心跳时间信息，这个是Socket底层自动实现的心跳包，不基于协议层实现。默认小于0，不开启心跳检测，如果需要开启，设置 60\_000 比较合适，单位毫秒  Get or set the heartbeat time information of the Socket of the client. This is the heartbeat packet automatically implemented by the bottom layer of the Socket, not based on the protocol layer. The default value is less than 0, and heartbeat detection is not enabled. If you need to enable it, it is more appropriate to set 60\_000, in milliseconds. (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |
| 公共属性代码示例 | [Token](80db00a7-aaa5-164f-e87c-af5400dc80fa.htm) | 网络类的身份令牌，在hsl协议的模式下会有效，在和设备进行通信的时候是无效的  Network-type identity tokens will be valid in the hsl protocol mode and will not be valid when communicating with the device (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的属性 | [UseServerActivePush](0e77fcc5-1744-cedf-a8e1-958031151826.htm) | 获取或设置当前的连接是否激活从服务器主动推送的功能 (继承自 [NetworkDoubleBase](0ebfedc3-835f-c7a0-94ef-e2ed3ba953d8.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[KeyenceDLEN1 类](f70c1a08-5573-e24a-5e62-6160f8d51d85.htm)

[HslCommunication.Profinet.Keyence 命名空间](baaebe57-0dd8-6393-ef35-013cbae668ce.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)