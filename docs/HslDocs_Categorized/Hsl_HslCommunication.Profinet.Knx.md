# HslCommunication - HslCommunication.Profinet.Knx

> 分類頁數: 30



---
## HslCommunication.Profinet.Knx

[原文連結](http://api.hslcommunication.cn/html/e9ca9cc4-27ce-cc6f-ea0b-badd3d84fe43.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Knx](../html/e9ca9cc4-27ce-cc6f-ea0b-badd3d84fe43.htm "HslCommunication.Profinet.Knx")

[KnxCode 类](../html/34d5ccbc-483b-c8ae-ce2d-c2c37b1fe2a3.htm "KnxCode 类")

[KnxCode.GetData 委托](../html/b3609506-bec0-cd66-ed44-d55b914630cc.htm "KnxCode.GetData 委托")

[KnxCode.ReturnData 委托](../html/48e7711d-0df1-c42b-455a-c0da28af40c3.htm "KnxCode.ReturnData 委托")

[KnxUdp 类](../html/17f9544a-78c3-1c9f-72df-11e3d0fd6606.htm "KnxUdp 类")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslCommunication.Profinet.Knx 命名空间 |

[缺少 "N:HslCommunication.Profinet.Knx" 的 <summary> 文档]

![](../icons/SectionExpanded.png)类

|  | 类 | 说明 |
| --- | --- | --- |
| 公共类 | [KnxCode](34d5ccbc-483b-c8ae-ce2d-c2c37b1fe2a3.htm) | Knx协议 |
| 公共类 | [KnxUdp](17f9544a-78c3-1c9f-72df-11e3d0fd6606.htm) | Knx驱动，具体的用法参照demo |

![](../icons/SectionExpanded.png)委托

|  | 委托 | 说明 |
| --- | --- | --- |
| 公共委托 | [KnxCodeGetData](b3609506-bec0-cd66-ed44-d55b914630cc.htm) | 获取数据的委托 |
| 公共委托 | [KnxCodeReturnData](48e7711d-0df1-c42b-455a-c0da28af40c3.htm) | 返回数据的委托 |

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## KnxCode 类

[原文連結](http://api.hslcommunication.cn/html/34d5ccbc-483b-c8ae-ce2d-c2c37b1fe2a3.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Knx](../html/e9ca9cc4-27ce-cc6f-ea0b-badd3d84fe43.htm "HslCommunication.Profinet.Knx")

[KnxCode 类](../html/34d5ccbc-483b-c8ae-ce2d-c2c37b1fe2a3.htm "KnxCode 类")

[KnxCode 构造函数](../html/8c6b524b-4ea2-0fd9-4f68-e47bdb0d4072.htm "KnxCode 构造函数 ")

[KnxCode 属性](../html/27521f8e-229d-d074-2045-93701b5ba6c8.htm "KnxCode 属性")

[KnxCode 方法](../html/e884e5e7-ee06-735b-4cfd-b5d4d3defcda.htm "KnxCode 方法")

[KnxCode 事件](../html/ec9a1a8a-f847-6cdc-2ead-fb076720c249.htm "KnxCode 事件")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| KnxCode 类 |

Knx协议

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  HslCommunication.Profinet.KnxKnxCode

**命名空间：**
 [HslCommunication.Profinet.Knx](e9ca9cc4-27ce-cc6f-ea0b-badd3d84fe43.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class KnxCode
```

```
Public Class KnxCode
```

```
public ref class KnxCode
```

```
type KnxCode =  class end
```

KnxCode 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [KnxCode](8c6b524b-4ea2-0fd9-4f68-e47bdb0d4072.htm) | 初始化 KnxCode 类的一个新实例 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [Channel](becd6d47-af9a-05b5-23e9-efe9f0484a94.htm) | 通道 |
| 公共属性 | [IsConnect](4a7f2116-fa2e-c1c5-cebd-731e738e18aa.htm) | 连接状态 |
| 公共属性 | [SequenceCounter](6c59acad-fc68-435b-b63a-48466dccb9ef.htm) | 序号计数 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [Disconnect\_knx](6630651a-0b26-e2a0-6496-c96961f75d70.htm) | 关闭KNX连接 |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | [Get\_knx\_addr](beebc704-05ac-65b8-ef51-3cd12bfb1bb0.htm) | 暂时没有注释 |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法 | GetType | (继承自 Object。) |
| 公共方法 | [Handshake](d472b44f-a82a-9ef7-a428-be6cad9f2c9c.htm) | 返回握手报文 |
| 公共方法 | [KNX\_check](74f3fe4e-38d9-a3c2-4471-c69bdaea78fa.htm) | KNX报文解析 |
| 公共方法 | [Knx\_Resd\_step1](632e2d6b-58fb-26d3-a82d-5eaead6ac470.htm) | 从KNX获取数据 |
| 公共方法 | [knx\_server\_is\_real](5eefbb10-f6ed-aa89-6e5b-ab9ba54ef8f0.htm) | 连接保持（每隔1s发送一次到设备） |
| 公共方法 | [Knx\_Write](1689ffa4-bb72-b1f5-4d46-1fbbafedcdab.htm) | 写入数据到KNX系统 |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法 | ToString | (继承自 Object。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)事件

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共事件 | [GetData\_msg](37ef9ff8-39ef-8201-e610-1730857abe98.htm) | 返回从knx系统得到的数据 |
| 公共事件 | [Return\_data\_msg](e7f128c6-9453-d8ce-4612-0bf738d92506.htm) | 返回需要写入KNX总线的应答报文（应答数据） |
| 公共事件 | [Set\_knx\_data](ade93b1f-a70e-46f9-ac3b-91da56a7a373.htm) | 返回需要写入的KNX系统的报文（写入数据） |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.Profinet.Knx 命名空间](e9ca9cc4-27ce-cc6f-ea0b-badd3d84fe43.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## KnxCode 构造函数 

[原文連結](http://api.hslcommunication.cn/html/8c6b524b-4ea2-0fd9-4f68-e47bdb0d4072.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Knx](../html/e9ca9cc4-27ce-cc6f-ea0b-badd3d84fe43.htm "HslCommunication.Profinet.Knx")

[KnxCode 类](../html/34d5ccbc-483b-c8ae-ce2d-c2c37b1fe2a3.htm "KnxCode 类")

[KnxCode 构造函数](../html/8c6b524b-4ea2-0fd9-4f68-e47bdb0d4072.htm "KnxCode 构造函数 ")

[KnxCode 属性](../html/27521f8e-229d-d074-2045-93701b5ba6c8.htm "KnxCode 属性")

[KnxCode 方法](../html/e884e5e7-ee06-735b-4cfd-b5d4d3defcda.htm "KnxCode 方法")

[KnxCode 事件](../html/ec9a1a8a-f847-6cdc-2ead-fb076720c249.htm "KnxCode 事件")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| KnxCode 构造函数 |

初始化 [KnxCode](34d5ccbc-483b-c8ae-ce2d-c2c37b1fe2a3.htm) 类的一个新实例

**命名空间：**
 [HslCommunication.Profinet.Knx](e9ca9cc4-27ce-cc6f-ea0b-badd3d84fe43.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public KnxCode()
```

```
Public Sub New
```

```
public:
KnxCode()
```

```
new : unit -> KnxCode
```

![](../icons/SectionExpanded.png)参见

#### 引用

[KnxCode 类](34d5ccbc-483b-c8ae-ce2d-c2c37b1fe2a3.htm)

[HslCommunication.Profinet.Knx 命名空间](e9ca9cc4-27ce-cc6f-ea0b-badd3d84fe43.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## KnxCode 属性

[原文連結](http://api.hslcommunication.cn/html/27521f8e-229d-d074-2045-93701b5ba6c8.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Knx](../html/e9ca9cc4-27ce-cc6f-ea0b-badd3d84fe43.htm "HslCommunication.Profinet.Knx")

[KnxCode 类](../html/34d5ccbc-483b-c8ae-ce2d-c2c37b1fe2a3.htm "KnxCode 类")

[KnxCode 属性](../html/27521f8e-229d-d074-2045-93701b5ba6c8.htm "KnxCode 属性")

[Channel 属性](../html/becd6d47-af9a-05b5-23e9-efe9f0484a94.htm "Channel 属性 ")

[IsConnect 属性](../html/4a7f2116-fa2e-c1c5-cebd-731e738e18aa.htm "IsConnect 属性 ")

[SequenceCounter 属性](../html/6c59acad-fc68-435b-b63a-48466dccb9ef.htm "SequenceCounter 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| KnxCode 属性 |

[KnxCode](34d5ccbc-483b-c8ae-ce2d-c2c37b1fe2a3.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [Channel](becd6d47-af9a-05b5-23e9-efe9f0484a94.htm) | 通道 |
| 公共属性 | [IsConnect](4a7f2116-fa2e-c1c5-cebd-731e738e18aa.htm) | 连接状态 |
| 公共属性 | [SequenceCounter](6c59acad-fc68-435b-b63a-48466dccb9ef.htm) | 序号计数 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[KnxCode 类](34d5ccbc-483b-c8ae-ce2d-c2c37b1fe2a3.htm)

[HslCommunication.Profinet.Knx 命名空间](e9ca9cc4-27ce-cc6f-ea0b-badd3d84fe43.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Channel 属性 

[原文連結](http://api.hslcommunication.cn/html/becd6d47-af9a-05b5-23e9-efe9f0484a94.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Knx](../html/e9ca9cc4-27ce-cc6f-ea0b-badd3d84fe43.htm "HslCommunication.Profinet.Knx")

[KnxCode 类](../html/34d5ccbc-483b-c8ae-ce2d-c2c37b1fe2a3.htm "KnxCode 类")

[KnxCode 属性](../html/27521f8e-229d-d074-2045-93701b5ba6c8.htm "KnxCode 属性")

[Channel 属性](../html/becd6d47-af9a-05b5-23e9-efe9f0484a94.htm "Channel 属性 ")

[IsConnect 属性](../html/4a7f2116-fa2e-c1c5-cebd-731e738e18aa.htm "IsConnect 属性 ")

[SequenceCounter 属性](../html/6c59acad-fc68-435b-b63a-48466dccb9ef.htm "SequenceCounter 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| KnxCodeChannel 属性 |

通道

**命名空间：**
 [HslCommunication.Profinet.Knx](e9ca9cc4-27ce-cc6f-ea0b-badd3d84fe43.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public byte Channel { get; set; }
```

```
Public Property Channel As Byte
	Get
	Set
```

```
public:
property unsigned char Channel {
	unsigned char get ();
	void set (unsigned char value);
}
```

```
member Channel : byte with get, set
```

#### 属性值

类型：Byte

![](../icons/SectionExpanded.png)参见

#### 引用

[KnxCode 类](34d5ccbc-483b-c8ae-ce2d-c2c37b1fe2a3.htm)

[HslCommunication.Profinet.Knx 命名空间](e9ca9cc4-27ce-cc6f-ea0b-badd3d84fe43.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## IsConnect 属性 

[原文連結](http://api.hslcommunication.cn/html/4a7f2116-fa2e-c1c5-cebd-731e738e18aa.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Knx](../html/e9ca9cc4-27ce-cc6f-ea0b-badd3d84fe43.htm "HslCommunication.Profinet.Knx")

[KnxCode 类](../html/34d5ccbc-483b-c8ae-ce2d-c2c37b1fe2a3.htm "KnxCode 类")

[KnxCode 属性](../html/27521f8e-229d-d074-2045-93701b5ba6c8.htm "KnxCode 属性")

[Channel 属性](../html/becd6d47-af9a-05b5-23e9-efe9f0484a94.htm "Channel 属性 ")

[IsConnect 属性](../html/4a7f2116-fa2e-c1c5-cebd-731e738e18aa.htm "IsConnect 属性 ")

[SequenceCounter 属性](../html/6c59acad-fc68-435b-b63a-48466dccb9ef.htm "SequenceCounter 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| KnxCodeIsConnect 属性 |

连接状态

**命名空间：**
 [HslCommunication.Profinet.Knx](e9ca9cc4-27ce-cc6f-ea0b-badd3d84fe43.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public bool IsConnect { get; }
```

```
Public ReadOnly Property IsConnect As Boolean
	Get
```

```
public:
property bool IsConnect {
	bool get ();
}
```

```
member IsConnect : bool with get
```

#### 属性值

类型：Boolean

![](../icons/SectionExpanded.png)参见

#### 引用

[KnxCode 类](34d5ccbc-483b-c8ae-ce2d-c2c37b1fe2a3.htm)

[HslCommunication.Profinet.Knx 命名空间](e9ca9cc4-27ce-cc6f-ea0b-badd3d84fe43.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## SequenceCounter 属性 

[原文連結](http://api.hslcommunication.cn/html/6c59acad-fc68-435b-b63a-48466dccb9ef.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Knx](../html/e9ca9cc4-27ce-cc6f-ea0b-badd3d84fe43.htm "HslCommunication.Profinet.Knx")

[KnxCode 类](../html/34d5ccbc-483b-c8ae-ce2d-c2c37b1fe2a3.htm "KnxCode 类")

[KnxCode 属性](../html/27521f8e-229d-d074-2045-93701b5ba6c8.htm "KnxCode 属性")

[Channel 属性](../html/becd6d47-af9a-05b5-23e9-efe9f0484a94.htm "Channel 属性 ")

[IsConnect 属性](../html/4a7f2116-fa2e-c1c5-cebd-731e738e18aa.htm "IsConnect 属性 ")

[SequenceCounter 属性](../html/6c59acad-fc68-435b-b63a-48466dccb9ef.htm "SequenceCounter 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| KnxCodeSequenceCounter 属性 |

序号计数

**命名空间：**
 [HslCommunication.Profinet.Knx](e9ca9cc4-27ce-cc6f-ea0b-badd3d84fe43.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public byte SequenceCounter { get; set; }
```

```
Public Property SequenceCounter As Byte
	Get
	Set
```

```
public:
property unsigned char SequenceCounter {
	unsigned char get ();
	void set (unsigned char value);
}
```

```
member SequenceCounter : byte with get, set
```

#### 属性值

类型：Byte

![](../icons/SectionExpanded.png)参见

#### 引用

[KnxCode 类](34d5ccbc-483b-c8ae-ce2d-c2c37b1fe2a3.htm)

[HslCommunication.Profinet.Knx 命名空间](e9ca9cc4-27ce-cc6f-ea0b-badd3d84fe43.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## KnxCode 方法

[原文連結](http://api.hslcommunication.cn/html/e884e5e7-ee06-735b-4cfd-b5d4d3defcda.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Knx](../html/e9ca9cc4-27ce-cc6f-ea0b-badd3d84fe43.htm "HslCommunication.Profinet.Knx")

[KnxCode 类](../html/34d5ccbc-483b-c8ae-ce2d-c2c37b1fe2a3.htm "KnxCode 类")

[KnxCode 方法](../html/e884e5e7-ee06-735b-4cfd-b5d4d3defcda.htm "KnxCode 方法")

[Disconnect\_knx 方法](../html/6630651a-0b26-e2a0-6496-c96961f75d70.htm "Disconnect_knx 方法 ")

[Get\_knx\_addr 方法](../html/beebc704-05ac-65b8-ef51-3cd12bfb1bb0.htm "Get_knx_addr 方法 ")

[Handshake 方法](../html/d472b44f-a82a-9ef7-a428-be6cad9f2c9c.htm "Handshake 方法 ")

[KNX\_check 方法](../html/74f3fe4e-38d9-a3c2-4471-c69bdaea78fa.htm "KNX_check 方法 ")

[Knx\_Resd\_step1 方法](../html/632e2d6b-58fb-26d3-a82d-5eaead6ac470.htm "Knx_Resd_step1 方法 ")

[knx\_server\_is\_real 方法](../html/5eefbb10-f6ed-aa89-6e5b-ab9ba54ef8f0.htm "knx_server_is_real 方法 ")

[Knx\_Write 方法](../html/1689ffa4-bb72-b1f5-4d46-1fbbafedcdab.htm "Knx_Write 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| KnxCode 方法 |

[KnxCode](34d5ccbc-483b-c8ae-ce2d-c2c37b1fe2a3.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [Disconnect\_knx](6630651a-0b26-e2a0-6496-c96961f75d70.htm) | 关闭KNX连接 |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | [Get\_knx\_addr](beebc704-05ac-65b8-ef51-3cd12bfb1bb0.htm) | 暂时没有注释 |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法 | GetType | (继承自 Object。) |
| 公共方法 | [Handshake](d472b44f-a82a-9ef7-a428-be6cad9f2c9c.htm) | 返回握手报文 |
| 公共方法 | [KNX\_check](74f3fe4e-38d9-a3c2-4471-c69bdaea78fa.htm) | KNX报文解析 |
| 公共方法 | [Knx\_Resd\_step1](632e2d6b-58fb-26d3-a82d-5eaead6ac470.htm) | 从KNX获取数据 |
| 公共方法 | [knx\_server\_is\_real](5eefbb10-f6ed-aa89-6e5b-ab9ba54ef8f0.htm) | 连接保持（每隔1s发送一次到设备） |
| 公共方法 | [Knx\_Write](1689ffa4-bb72-b1f5-4d46-1fbbafedcdab.htm) | 写入数据到KNX系统 |
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

[KnxCode 类](34d5ccbc-483b-c8ae-ce2d-c2c37b1fe2a3.htm)

[HslCommunication.Profinet.Knx 命名空间](e9ca9cc4-27ce-cc6f-ea0b-badd3d84fe43.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Disconnect_knx 方法 

[原文連結](http://api.hslcommunication.cn/html/6630651a-0b26-e2a0-6496-c96961f75d70.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Knx](../html/e9ca9cc4-27ce-cc6f-ea0b-badd3d84fe43.htm "HslCommunication.Profinet.Knx")

[KnxCode 类](../html/34d5ccbc-483b-c8ae-ce2d-c2c37b1fe2a3.htm "KnxCode 类")

[KnxCode 方法](../html/e884e5e7-ee06-735b-4cfd-b5d4d3defcda.htm "KnxCode 方法")

[Disconnect\_knx 方法](../html/6630651a-0b26-e2a0-6496-c96961f75d70.htm "Disconnect_knx 方法 ")

[Get\_knx\_addr 方法](../html/beebc704-05ac-65b8-ef51-3cd12bfb1bb0.htm "Get_knx_addr 方法 ")

[Handshake 方法](../html/d472b44f-a82a-9ef7-a428-be6cad9f2c9c.htm "Handshake 方法 ")

[KNX\_check 方法](../html/74f3fe4e-38d9-a3c2-4471-c69bdaea78fa.htm "KNX_check 方法 ")

[Knx\_Resd\_step1 方法](../html/632e2d6b-58fb-26d3-a82d-5eaead6ac470.htm "Knx_Resd_step1 方法 ")

[knx\_server\_is\_real 方法](../html/5eefbb10-f6ed-aa89-6e5b-ab9ba54ef8f0.htm "knx_server_is_real 方法 ")

[Knx\_Write 方法](../html/1689ffa4-bb72-b1f5-4d46-1fbbafedcdab.htm "Knx_Write 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| KnxCodeDisconnect\_knx 方法 |

关闭KNX连接

**命名空间：**
 [HslCommunication.Profinet.Knx](e9ca9cc4-27ce-cc6f-ea0b-badd3d84fe43.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public byte[] Disconnect_knx(
	byte channel,
	IPEndPoint IP_PROT
)
```

```
Public Function Disconnect_knx ( 
	channel As Byte,
	IP_PROT As IPEndPoint
) As Byte()
```

```
public:
array<unsigned char>^ Disconnect_knx(
	unsigned char channel, 
	IPEndPoint^ IP_PROT
)
```

```
member Disconnect_knx : 
        channel : byte * 
        IP_PROT : IPEndPoint -> byte[] 
```

#### 参数

channel
:   类型：SystemByte  
    通道号

IP\_PROT
:   类型：System.NetIPEndPoint  
    本机IP

#### 返回值

类型：Byte  

[缺少 "M:HslCommunication.Profinet.Knx.KnxCode.Disconnect\_knx(System.Byte,System.Net.IPEndPoint)" 的 <returns> 文档]

![](../icons/SectionExpanded.png)参见

#### 引用

[KnxCode 类](34d5ccbc-483b-c8ae-ce2d-c2c37b1fe2a3.htm)

[HslCommunication.Profinet.Knx 命名空间](e9ca9cc4-27ce-cc6f-ea0b-badd3d84fe43.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Get_knx_addr 方法 

[原文連結](http://api.hslcommunication.cn/html/beebc704-05ac-65b8-ef51-3cd12bfb1bb0.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Knx](../html/e9ca9cc4-27ce-cc6f-ea0b-badd3d84fe43.htm "HslCommunication.Profinet.Knx")

[KnxCode 类](../html/34d5ccbc-483b-c8ae-ce2d-c2c37b1fe2a3.htm "KnxCode 类")

[KnxCode 方法](../html/e884e5e7-ee06-735b-4cfd-b5d4d3defcda.htm "KnxCode 方法")

[Disconnect\_knx 方法](../html/6630651a-0b26-e2a0-6496-c96961f75d70.htm "Disconnect_knx 方法 ")

[Get\_knx\_addr 方法](../html/beebc704-05ac-65b8-ef51-3cd12bfb1bb0.htm "Get_knx_addr 方法 ")

[Handshake 方法](../html/d472b44f-a82a-9ef7-a428-be6cad9f2c9c.htm "Handshake 方法 ")

[KNX\_check 方法](../html/74f3fe4e-38d9-a3c2-4471-c69bdaea78fa.htm "KNX_check 方法 ")

[Knx\_Resd\_step1 方法](../html/632e2d6b-58fb-26d3-a82d-5eaead6ac470.htm "Knx_Resd_step1 方法 ")

[knx\_server\_is\_real 方法](../html/5eefbb10-f6ed-aa89-6e5b-ab9ba54ef8f0.htm "knx_server_is_real 方法 ")

[Knx\_Write 方法](../html/1689ffa4-bb72-b1f5-4d46-1fbbafedcdab.htm "Knx_Write 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| KnxCodeGet\_knx\_addr 方法 |

暂时没有注释

**命名空间：**
 [HslCommunication.Profinet.Knx](e9ca9cc4-27ce-cc6f-ea0b-badd3d84fe43.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public short Get_knx_addr(
	string addr,
	out bool is_ok
)
```

```
Public Function Get_knx_addr ( 
	addr As String,
	<OutAttribute> ByRef is_ok As Boolean
) As Short
```

```
public:
short Get_knx_addr(
	String^ addr, 
	[OutAttribute] bool% is_ok
)
```

```
member Get_knx_addr : 
        addr : string * 
        is_ok : bool byref -> int16 
```

#### 参数

addr
:   类型：SystemString  

    [缺少 "M:HslCommunication.Profinet.Knx.KnxCode.Get\_knx\_addr(System.String,System.Boolean@)" 的 <param name="addr"/> 文档]

is\_ok
:   类型：SystemBoolean  

    [缺少 "M:HslCommunication.Profinet.Knx.KnxCode.Get\_knx\_addr(System.String,System.Boolean@)" 的 <param name="is\_ok"/> 文档]

#### 返回值

类型：Int16  

[缺少 "M:HslCommunication.Profinet.Knx.KnxCode.Get\_knx\_addr(System.String,System.Boolean@)" 的 <returns> 文档]

![](../icons/SectionExpanded.png)参见

#### 引用

[KnxCode 类](34d5ccbc-483b-c8ae-ce2d-c2c37b1fe2a3.htm)

[HslCommunication.Profinet.Knx 命名空间](e9ca9cc4-27ce-cc6f-ea0b-badd3d84fe43.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Handshake 方法 

[原文連結](http://api.hslcommunication.cn/html/d472b44f-a82a-9ef7-a428-be6cad9f2c9c.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Knx](../html/e9ca9cc4-27ce-cc6f-ea0b-badd3d84fe43.htm "HslCommunication.Profinet.Knx")

[KnxCode 类](../html/34d5ccbc-483b-c8ae-ce2d-c2c37b1fe2a3.htm "KnxCode 类")

[KnxCode 方法](../html/e884e5e7-ee06-735b-4cfd-b5d4d3defcda.htm "KnxCode 方法")

[Disconnect\_knx 方法](../html/6630651a-0b26-e2a0-6496-c96961f75d70.htm "Disconnect_knx 方法 ")

[Get\_knx\_addr 方法](../html/beebc704-05ac-65b8-ef51-3cd12bfb1bb0.htm "Get_knx_addr 方法 ")

[Handshake 方法](../html/d472b44f-a82a-9ef7-a428-be6cad9f2c9c.htm "Handshake 方法 ")

[KNX\_check 方法](../html/74f3fe4e-38d9-a3c2-4471-c69bdaea78fa.htm "KNX_check 方法 ")

[Knx\_Resd\_step1 方法](../html/632e2d6b-58fb-26d3-a82d-5eaead6ac470.htm "Knx_Resd_step1 方法 ")

[knx\_server\_is\_real 方法](../html/5eefbb10-f6ed-aa89-6e5b-ab9ba54ef8f0.htm "knx_server_is_real 方法 ")

[Knx\_Write 方法](../html/1689ffa4-bb72-b1f5-4d46-1fbbafedcdab.htm "Knx_Write 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| KnxCodeHandshake 方法 |

返回握手报文

**命名空间：**
 [HslCommunication.Profinet.Knx](e9ca9cc4-27ce-cc6f-ea0b-badd3d84fe43.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public byte[] Handshake(
	IPEndPoint IP_PROT
)
```

```
Public Function Handshake ( 
	IP_PROT As IPEndPoint
) As Byte()
```

```
public:
array<unsigned char>^ Handshake(
	IPEndPoint^ IP_PROT
)
```

```
member Handshake : 
        IP_PROT : IPEndPoint -> byte[] 
```

#### 参数

IP\_PROT
:   类型：System.NetIPEndPoint  
    本机ip地址

#### 返回值

类型：Byte  

[缺少 "M:HslCommunication.Profinet.Knx.KnxCode.Handshake(System.Net.IPEndPoint)" 的 <returns> 文档]

![](../icons/SectionExpanded.png)参见

#### 引用

[KnxCode 类](34d5ccbc-483b-c8ae-ce2d-c2c37b1fe2a3.htm)

[HslCommunication.Profinet.Knx 命名空间](e9ca9cc4-27ce-cc6f-ea0b-badd3d84fe43.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## KNX_check 方法 

[原文連結](http://api.hslcommunication.cn/html/74f3fe4e-38d9-a3c2-4471-c69bdaea78fa.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Knx](../html/e9ca9cc4-27ce-cc6f-ea0b-badd3d84fe43.htm "HslCommunication.Profinet.Knx")

[KnxCode 类](../html/34d5ccbc-483b-c8ae-ce2d-c2c37b1fe2a3.htm "KnxCode 类")

[KnxCode 方法](../html/e884e5e7-ee06-735b-4cfd-b5d4d3defcda.htm "KnxCode 方法")

[Disconnect\_knx 方法](../html/6630651a-0b26-e2a0-6496-c96961f75d70.htm "Disconnect_knx 方法 ")

[Get\_knx\_addr 方法](../html/beebc704-05ac-65b8-ef51-3cd12bfb1bb0.htm "Get_knx_addr 方法 ")

[Handshake 方法](../html/d472b44f-a82a-9ef7-a428-be6cad9f2c9c.htm "Handshake 方法 ")

[KNX\_check 方法](../html/74f3fe4e-38d9-a3c2-4471-c69bdaea78fa.htm "KNX_check 方法 ")

[Knx\_Resd\_step1 方法](../html/632e2d6b-58fb-26d3-a82d-5eaead6ac470.htm "Knx_Resd_step1 方法 ")

[knx\_server\_is\_real 方法](../html/5eefbb10-f6ed-aa89-6e5b-ab9ba54ef8f0.htm "knx_server_is_real 方法 ")

[Knx\_Write 方法](../html/1689ffa4-bb72-b1f5-4d46-1fbbafedcdab.htm "Knx_Write 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| KnxCodeKNX\_check 方法 |

KNX报文解析

**命名空间：**
 [HslCommunication.Profinet.Knx](e9ca9cc4-27ce-cc6f-ea0b-badd3d84fe43.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public void KNX_check(
	byte[] in_data
)
```

```
Public Sub KNX_check ( 
	in_data As Byte()
)
```

```
public:
void KNX_check(
	array<unsigned char>^ in_data
)
```

```
member KNX_check : 
        in_data : byte[] -> unit 
```

#### 参数

in\_data
:   类型：SystemByte  

    [缺少 "M:HslCommunication.Profinet.Knx.KnxCode.KNX\_check(System.Byte[])" 的 <param name="in\_data"/> 文档]

![](../icons/SectionExpanded.png)参见

#### 引用

[KnxCode 类](34d5ccbc-483b-c8ae-ce2d-c2c37b1fe2a3.htm)

[HslCommunication.Profinet.Knx 命名空间](e9ca9cc4-27ce-cc6f-ea0b-badd3d84fe43.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Knx_Resd_step1 方法 

[原文連結](http://api.hslcommunication.cn/html/632e2d6b-58fb-26d3-a82d-5eaead6ac470.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Knx](../html/e9ca9cc4-27ce-cc6f-ea0b-badd3d84fe43.htm "HslCommunication.Profinet.Knx")

[KnxCode 类](../html/34d5ccbc-483b-c8ae-ce2d-c2c37b1fe2a3.htm "KnxCode 类")

[KnxCode 方法](../html/e884e5e7-ee06-735b-4cfd-b5d4d3defcda.htm "KnxCode 方法")

[Disconnect\_knx 方法](../html/6630651a-0b26-e2a0-6496-c96961f75d70.htm "Disconnect_knx 方法 ")

[Get\_knx\_addr 方法](../html/beebc704-05ac-65b8-ef51-3cd12bfb1bb0.htm "Get_knx_addr 方法 ")

[Handshake 方法](../html/d472b44f-a82a-9ef7-a428-be6cad9f2c9c.htm "Handshake 方法 ")

[KNX\_check 方法](../html/74f3fe4e-38d9-a3c2-4471-c69bdaea78fa.htm "KNX_check 方法 ")

[Knx\_Resd\_step1 方法](../html/632e2d6b-58fb-26d3-a82d-5eaead6ac470.htm "Knx_Resd_step1 方法 ")

[knx\_server\_is\_real 方法](../html/5eefbb10-f6ed-aa89-6e5b-ab9ba54ef8f0.htm "knx_server_is_real 方法 ")

[Knx\_Write 方法](../html/1689ffa4-bb72-b1f5-4d46-1fbbafedcdab.htm "Knx_Write 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| KnxCodeKnx\_Resd\_step1 方法 |

从KNX获取数据

**命名空间：**
 [HslCommunication.Profinet.Knx](e9ca9cc4-27ce-cc6f-ea0b-badd3d84fe43.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public void Knx_Resd_step1(
	short addr
)
```

```
Public Sub Knx_Resd_step1 ( 
	addr As Short
)
```

```
public:
void Knx_Resd_step1(
	short addr
)
```

```
member Knx_Resd_step1 : 
        addr : int16 -> unit 
```

#### 参数

addr
:   类型：SystemInt16  

    [缺少 "M:HslCommunication.Profinet.Knx.KnxCode.Knx\_Resd\_step1(System.Int16)" 的 <param name="addr"/> 文档]

![](../icons/SectionExpanded.png)参见

#### 引用

[KnxCode 类](34d5ccbc-483b-c8ae-ce2d-c2c37b1fe2a3.htm)

[HslCommunication.Profinet.Knx 命名空间](e9ca9cc4-27ce-cc6f-ea0b-badd3d84fe43.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## knx_server_is_real 方法 

[原文連結](http://api.hslcommunication.cn/html/5eefbb10-f6ed-aa89-6e5b-ab9ba54ef8f0.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Knx](../html/e9ca9cc4-27ce-cc6f-ea0b-badd3d84fe43.htm "HslCommunication.Profinet.Knx")

[KnxCode 类](../html/34d5ccbc-483b-c8ae-ce2d-c2c37b1fe2a3.htm "KnxCode 类")

[KnxCode 方法](../html/e884e5e7-ee06-735b-4cfd-b5d4d3defcda.htm "KnxCode 方法")

[Disconnect\_knx 方法](../html/6630651a-0b26-e2a0-6496-c96961f75d70.htm "Disconnect_knx 方法 ")

[Get\_knx\_addr 方法](../html/beebc704-05ac-65b8-ef51-3cd12bfb1bb0.htm "Get_knx_addr 方法 ")

[Handshake 方法](../html/d472b44f-a82a-9ef7-a428-be6cad9f2c9c.htm "Handshake 方法 ")

[KNX\_check 方法](../html/74f3fe4e-38d9-a3c2-4471-c69bdaea78fa.htm "KNX_check 方法 ")

[Knx\_Resd\_step1 方法](../html/632e2d6b-58fb-26d3-a82d-5eaead6ac470.htm "Knx_Resd_step1 方法 ")

[knx\_server\_is\_real 方法](../html/5eefbb10-f6ed-aa89-6e5b-ab9ba54ef8f0.htm "knx_server_is_real 方法 ")

[Knx\_Write 方法](../html/1689ffa4-bb72-b1f5-4d46-1fbbafedcdab.htm "Knx_Write 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| KnxCodeknx\_server\_is\_real 方法 |

连接保持（每隔1s发送一次到设备）

**命名空间：**
 [HslCommunication.Profinet.Knx](e9ca9cc4-27ce-cc6f-ea0b-badd3d84fe43.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public void knx_server_is_real(
	IPEndPoint IP_PROT
)
```

```
Public Sub knx_server_is_real ( 
	IP_PROT As IPEndPoint
)
```

```
public:
void knx_server_is_real(
	IPEndPoint^ IP_PROT
)
```

```
member knx_server_is_real : 
        IP_PROT : IPEndPoint -> unit 
```

#### 参数

IP\_PROT
:   类型：System.NetIPEndPoint  

    [缺少 "M:HslCommunication.Profinet.Knx.KnxCode.knx\_server\_is\_real(System.Net.IPEndPoint)" 的 <param name="IP\_PROT"/> 文档]

![](../icons/SectionExpanded.png)参见

#### 引用

[KnxCode 类](34d5ccbc-483b-c8ae-ce2d-c2c37b1fe2a3.htm)

[HslCommunication.Profinet.Knx 命名空间](e9ca9cc4-27ce-cc6f-ea0b-badd3d84fe43.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Knx_Write 方法 

[原文連結](http://api.hslcommunication.cn/html/1689ffa4-bb72-b1f5-4d46-1fbbafedcdab.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Knx](../html/e9ca9cc4-27ce-cc6f-ea0b-badd3d84fe43.htm "HslCommunication.Profinet.Knx")

[KnxCode 类](../html/34d5ccbc-483b-c8ae-ce2d-c2c37b1fe2a3.htm "KnxCode 类")

[KnxCode 方法](../html/e884e5e7-ee06-735b-4cfd-b5d4d3defcda.htm "KnxCode 方法")

[Disconnect\_knx 方法](../html/6630651a-0b26-e2a0-6496-c96961f75d70.htm "Disconnect_knx 方法 ")

[Get\_knx\_addr 方法](../html/beebc704-05ac-65b8-ef51-3cd12bfb1bb0.htm "Get_knx_addr 方法 ")

[Handshake 方法](../html/d472b44f-a82a-9ef7-a428-be6cad9f2c9c.htm "Handshake 方法 ")

[KNX\_check 方法](../html/74f3fe4e-38d9-a3c2-4471-c69bdaea78fa.htm "KNX_check 方法 ")

[Knx\_Resd\_step1 方法](../html/632e2d6b-58fb-26d3-a82d-5eaead6ac470.htm "Knx_Resd_step1 方法 ")

[knx\_server\_is\_real 方法](../html/5eefbb10-f6ed-aa89-6e5b-ab9ba54ef8f0.htm "knx_server_is_real 方法 ")

[Knx\_Write 方法](../html/1689ffa4-bb72-b1f5-4d46-1fbbafedcdab.htm "Knx_Write 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| KnxCodeKnx\_Write 方法 |

写入数据到KNX系统

**命名空间：**
 [HslCommunication.Profinet.Knx](e9ca9cc4-27ce-cc6f-ea0b-badd3d84fe43.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public void Knx_Write(
	short addr,
	byte len,
	byte[] data
)
```

```
Public Sub Knx_Write ( 
	addr As Short,
	len As Byte,
	data As Byte()
)
```

```
public:
void Knx_Write(
	short addr, 
	unsigned char len, 
	array<unsigned char>^ data
)
```

```
member Knx_Write : 
        addr : int16 * 
        len : byte * 
        data : byte[] -> unit 
```

#### 参数

addr
:   类型：SystemInt16  
    地址

len
:   类型：SystemByte  
    长度

data
:   类型：SystemByte  
    数据

![](../icons/SectionExpanded.png)参见

#### 引用

[KnxCode 类](34d5ccbc-483b-c8ae-ce2d-c2c37b1fe2a3.htm)

[HslCommunication.Profinet.Knx 命名空间](e9ca9cc4-27ce-cc6f-ea0b-badd3d84fe43.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## KnxCode 事件

[原文連結](http://api.hslcommunication.cn/html/ec9a1a8a-f847-6cdc-2ead-fb076720c249.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Knx](../html/e9ca9cc4-27ce-cc6f-ea0b-badd3d84fe43.htm "HslCommunication.Profinet.Knx")

[KnxCode 类](../html/34d5ccbc-483b-c8ae-ce2d-c2c37b1fe2a3.htm "KnxCode 类")

[KnxCode 事件](../html/ec9a1a8a-f847-6cdc-2ead-fb076720c249.htm "KnxCode 事件")

[GetData\_msg 事件](../html/37ef9ff8-39ef-8201-e610-1730857abe98.htm "GetData_msg 事件")

[Return\_data\_msg 事件](../html/e7f128c6-9453-d8ce-4612-0bf738d92506.htm "Return_data_msg 事件")

[Set\_knx\_data 事件](../html/ade93b1f-a70e-46f9-ac3b-91da56a7a373.htm "Set_knx_data 事件")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| KnxCode 事件 |

[KnxCode](34d5ccbc-483b-c8ae-ce2d-c2c37b1fe2a3.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)事件

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共事件 | [GetData\_msg](37ef9ff8-39ef-8201-e610-1730857abe98.htm) | 返回从knx系统得到的数据 |
| 公共事件 | [Return\_data\_msg](e7f128c6-9453-d8ce-4612-0bf738d92506.htm) | 返回需要写入KNX总线的应答报文（应答数据） |
| 公共事件 | [Set\_knx\_data](ade93b1f-a70e-46f9-ac3b-91da56a7a373.htm) | 返回需要写入的KNX系统的报文（写入数据） |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[KnxCode 类](34d5ccbc-483b-c8ae-ce2d-c2c37b1fe2a3.htm)

[HslCommunication.Profinet.Knx 命名空间](e9ca9cc4-27ce-cc6f-ea0b-badd3d84fe43.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GetData_msg 事件

[原文連結](http://api.hslcommunication.cn/html/37ef9ff8-39ef-8201-e610-1730857abe98.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Knx](../html/e9ca9cc4-27ce-cc6f-ea0b-badd3d84fe43.htm "HslCommunication.Profinet.Knx")

[KnxCode 类](../html/34d5ccbc-483b-c8ae-ce2d-c2c37b1fe2a3.htm "KnxCode 类")

[KnxCode 事件](../html/ec9a1a8a-f847-6cdc-2ead-fb076720c249.htm "KnxCode 事件")

[GetData\_msg 事件](../html/37ef9ff8-39ef-8201-e610-1730857abe98.htm "GetData_msg 事件")

[Return\_data\_msg 事件](../html/e7f128c6-9453-d8ce-4612-0bf738d92506.htm "Return_data_msg 事件")

[Set\_knx\_data 事件](../html/ade93b1f-a70e-46f9-ac3b-91da56a7a373.htm "Set_knx_data 事件")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| KnxCodeGetData\_msg 事件 |

返回从knx系统得到的数据

**命名空间：**
 [HslCommunication.Profinet.Knx](e9ca9cc4-27ce-cc6f-ea0b-badd3d84fe43.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public event KnxCodeGetData GetData_msg
```

```
Public Event GetData_msg As KnxCodeGetData
```

```
public:
 event KnxCodeGetData^ GetData_msg {
	void add (KnxCodeGetData^ value);
	void remove (KnxCodeGetData^ value);
}
```

```
member GetData_msg : IEvent<KnxCodeGetData,
    EventArgs>
```

#### 值

类型：[HslCommunication.Profinet.KnxKnxCodeGetData](b3609506-bec0-cd66-ed44-d55b914630cc.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[KnxCode 类](34d5ccbc-483b-c8ae-ce2d-c2c37b1fe2a3.htm)

[HslCommunication.Profinet.Knx 命名空间](e9ca9cc4-27ce-cc6f-ea0b-badd3d84fe43.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Return_data_msg 事件

[原文連結](http://api.hslcommunication.cn/html/e7f128c6-9453-d8ce-4612-0bf738d92506.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Knx](../html/e9ca9cc4-27ce-cc6f-ea0b-badd3d84fe43.htm "HslCommunication.Profinet.Knx")

[KnxCode 类](../html/34d5ccbc-483b-c8ae-ce2d-c2c37b1fe2a3.htm "KnxCode 类")

[KnxCode 事件](../html/ec9a1a8a-f847-6cdc-2ead-fb076720c249.htm "KnxCode 事件")

[GetData\_msg 事件](../html/37ef9ff8-39ef-8201-e610-1730857abe98.htm "GetData_msg 事件")

[Return\_data\_msg 事件](../html/e7f128c6-9453-d8ce-4612-0bf738d92506.htm "Return_data_msg 事件")

[Set\_knx\_data 事件](../html/ade93b1f-a70e-46f9-ac3b-91da56a7a373.htm "Set_knx_data 事件")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| KnxCodeReturn\_data\_msg 事件 |

返回需要写入KNX总线的应答报文（应答数据）

**命名空间：**
 [HslCommunication.Profinet.Knx](e9ca9cc4-27ce-cc6f-ea0b-badd3d84fe43.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public event KnxCodeReturnData Return_data_msg
```

```
Public Event Return_data_msg As KnxCodeReturnData
```

```
public:
 event KnxCodeReturnData^ Return_data_msg {
	void add (KnxCodeReturnData^ value);
	void remove (KnxCodeReturnData^ value);
}
```

```
member Return_data_msg : IEvent<KnxCodeReturnData,
    EventArgs>
```

#### 值

类型：[HslCommunication.Profinet.KnxKnxCodeReturnData](48e7711d-0df1-c42b-455a-c0da28af40c3.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[KnxCode 类](34d5ccbc-483b-c8ae-ce2d-c2c37b1fe2a3.htm)

[HslCommunication.Profinet.Knx 命名空间](e9ca9cc4-27ce-cc6f-ea0b-badd3d84fe43.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Set_knx_data 事件

[原文連結](http://api.hslcommunication.cn/html/ade93b1f-a70e-46f9-ac3b-91da56a7a373.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Knx](../html/e9ca9cc4-27ce-cc6f-ea0b-badd3d84fe43.htm "HslCommunication.Profinet.Knx")

[KnxCode 类](../html/34d5ccbc-483b-c8ae-ce2d-c2c37b1fe2a3.htm "KnxCode 类")

[KnxCode 事件](../html/ec9a1a8a-f847-6cdc-2ead-fb076720c249.htm "KnxCode 事件")

[GetData\_msg 事件](../html/37ef9ff8-39ef-8201-e610-1730857abe98.htm "GetData_msg 事件")

[Return\_data\_msg 事件](../html/e7f128c6-9453-d8ce-4612-0bf738d92506.htm "Return_data_msg 事件")

[Set\_knx\_data 事件](../html/ade93b1f-a70e-46f9-ac3b-91da56a7a373.htm "Set_knx_data 事件")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| KnxCodeSet\_knx\_data 事件 |

返回需要写入的KNX系统的报文（写入数据）

**命名空间：**
 [HslCommunication.Profinet.Knx](e9ca9cc4-27ce-cc6f-ea0b-badd3d84fe43.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public event KnxCodeReturnData Set_knx_data
```

```
Public Event Set_knx_data As KnxCodeReturnData
```

```
public:
 event KnxCodeReturnData^ Set_knx_data {
	void add (KnxCodeReturnData^ value);
	void remove (KnxCodeReturnData^ value);
}
```

```
member Set_knx_data : IEvent<KnxCodeReturnData,
    EventArgs>
```

#### 值

类型：[HslCommunication.Profinet.KnxKnxCodeReturnData](48e7711d-0df1-c42b-455a-c0da28af40c3.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[KnxCode 类](34d5ccbc-483b-c8ae-ce2d-c2c37b1fe2a3.htm)

[HslCommunication.Profinet.Knx 命名空间](e9ca9cc4-27ce-cc6f-ea0b-badd3d84fe43.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## KnxCode.GetData 委托

[原文連結](http://api.hslcommunication.cn/html/b3609506-bec0-cd66-ed44-d55b914630cc.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Knx](../html/e9ca9cc4-27ce-cc6f-ea0b-badd3d84fe43.htm "HslCommunication.Profinet.Knx")

[KnxCode 类](../html/34d5ccbc-483b-c8ae-ce2d-c2c37b1fe2a3.htm "KnxCode 类")

[KnxCode.GetData 委托](../html/b3609506-bec0-cd66-ed44-d55b914630cc.htm "KnxCode.GetData 委托")

[KnxCode.ReturnData 委托](../html/48e7711d-0df1-c42b-455a-c0da28af40c3.htm "KnxCode.ReturnData 委托")

[KnxUdp 类](../html/17f9544a-78c3-1c9f-72df-11e3d0fd6606.htm "KnxUdp 类")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| KnxCodeGetData 委托 |

获取数据的委托

**命名空间：**
 [HslCommunication.Profinet.Knx](e9ca9cc4-27ce-cc6f-ea0b-badd3d84fe43.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public delegate void GetData(
	short addr,
	byte len,
	byte[] data
)
```

```
Public Delegate Sub GetData ( 
	addr As Short,
	len As Byte,
	data As Byte()
)
```

```
public delegate void GetData(
	short addr, 
	unsigned char len, 
	array<unsigned char>^ data
)
```

```
type GetData = 
    delegate of 
        addr : int16 * 
        len : byte * 
        data : byte[] -> unit
```

#### 参数

addr
:   类型：SystemInt16

len
:   类型：SystemByte

data
:   类型：SystemByte

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.Profinet.Knx 命名空间](e9ca9cc4-27ce-cc6f-ea0b-badd3d84fe43.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## KnxCode.ReturnData 委托

[原文連結](http://api.hslcommunication.cn/html/48e7711d-0df1-c42b-455a-c0da28af40c3.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Knx](../html/e9ca9cc4-27ce-cc6f-ea0b-badd3d84fe43.htm "HslCommunication.Profinet.Knx")

[KnxCode 类](../html/34d5ccbc-483b-c8ae-ce2d-c2c37b1fe2a3.htm "KnxCode 类")

[KnxCode.GetData 委托](../html/b3609506-bec0-cd66-ed44-d55b914630cc.htm "KnxCode.GetData 委托")

[KnxCode.ReturnData 委托](../html/48e7711d-0df1-c42b-455a-c0da28af40c3.htm "KnxCode.ReturnData 委托")

[KnxUdp 类](../html/17f9544a-78c3-1c9f-72df-11e3d0fd6606.htm "KnxUdp 类")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| KnxCodeReturnData 委托 |

返回数据的委托

**命名空间：**
 [HslCommunication.Profinet.Knx](e9ca9cc4-27ce-cc6f-ea0b-badd3d84fe43.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public delegate void ReturnData(
	byte[] data
)
```

```
Public Delegate Sub ReturnData ( 
	data As Byte()
)
```

```
public delegate void ReturnData(
	array<unsigned char>^ data
)
```

```
type ReturnData = 
    delegate of 
        data : byte[] -> unit
```

#### 参数

data
:   类型：SystemByte

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.Profinet.Knx 命名空间](e9ca9cc4-27ce-cc6f-ea0b-badd3d84fe43.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## KnxUdp 类

[原文連結](http://api.hslcommunication.cn/html/17f9544a-78c3-1c9f-72df-11e3d0fd6606.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Knx](../html/e9ca9cc4-27ce-cc6f-ea0b-badd3d84fe43.htm "HslCommunication.Profinet.Knx")

[KnxUdp 类](../html/17f9544a-78c3-1c9f-72df-11e3d0fd6606.htm "KnxUdp 类")

[KnxUdp 构造函数](../html/afb00997-41ec-4539-00de-808a6238c044.htm "KnxUdp 构造函数 ")

[KnxUdp 属性](../html/3816a091-7361-922a-40c7-cb8a598ad4a1.htm "KnxUdp 属性")

[KnxUdp 方法](../html/aabb6ab6-b731-94c5-c252-dc03e513a741.htm "KnxUdp 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| KnxUdp 类 |

Knx驱动，具体的用法参照demo

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  HslCommunication.Profinet.KnxKnxUdp

**命名空间：**
 [HslCommunication.Profinet.Knx](e9ca9cc4-27ce-cc6f-ea0b-badd3d84fe43.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class KnxUdp
```

```
Public Class KnxUdp
```

```
public ref class KnxUdp
```

```
type KnxUdp =  class end
```

KnxUdp 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [KnxUdp](afb00997-41ec-4539-00de-808a6238c044.htm) | 实例化一个默认的对象 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [Channel](c03c4e81-e0c6-45b3-c5b6-990a7530174f.htm) | 通道号（由设备发来） |
| 公共属性 | [IsConnect](a1247873-ea2b-4177-6366-d9808e8adbdf.htm) | 当前的状态是否连接中 |
| 公共属性 | [KnxCode](1f3bb20c-1084-d8cb-5080-5923b57573d0.htm) | 通信指令类 |
| 公共属性 | [LocalEndpoint](b7e4b069-6bf2-6d8e-6502-e04a60cec773.htm) | 本机IP地址 |
| 公共属性 | [LogNet](af906e51-4f30-5076-0933-3e822a48b96f.htm) | 系统的日志信息 |
| 公共属性 | [RouEndpoint](a20342d5-ed6a-7803-931d-bbad42b7a1dc.htm) | 远程ip地址 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [ConnectKnx](348fa77a-8991-443d-2492-67e45c3abc17.htm) | 和KNX网络进行握手并开始监听 |
| 公共方法 | [DisConnectKnx](a656da6e-6ac5-0277-ef0e-ef2d4daacb82.htm) | 关闭连接 |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法 | GetType | (继承自 Object。) |
| 公共方法 | [KeepConnection](a0040ab8-aeba-3476-3b79-3e9cee6fc9de.htm) | 保持KNX连接 |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法 | [ReadKnxData](e1007297-a037-d259-fed1-e6735fead9fc.htm) | 读取指定KNX组地址 |
| 公共方法 | [SetKnxData](b6cac092-9681-f2a3-c3f1-4ecdcc0842be.htm) | 将报文写入KNX系统 |
| 公共方法 | ToString | (继承自 Object。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)备注

感谢上海NULL提供的技术支持

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.Profinet.Knx 命名空间](e9ca9cc4-27ce-cc6f-ea0b-badd3d84fe43.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## KnxUdp 构造函数 

[原文連結](http://api.hslcommunication.cn/html/afb00997-41ec-4539-00de-808a6238c044.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Knx](../html/e9ca9cc4-27ce-cc6f-ea0b-badd3d84fe43.htm "HslCommunication.Profinet.Knx")

[KnxUdp 类](../html/17f9544a-78c3-1c9f-72df-11e3d0fd6606.htm "KnxUdp 类")

[KnxUdp 构造函数](../html/afb00997-41ec-4539-00de-808a6238c044.htm "KnxUdp 构造函数 ")

[KnxUdp 属性](../html/3816a091-7361-922a-40c7-cb8a598ad4a1.htm "KnxUdp 属性")

[KnxUdp 方法](../html/aabb6ab6-b731-94c5-c252-dc03e513a741.htm "KnxUdp 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| KnxUdp 构造函数 |

实例化一个默认的对象

**命名空间：**
 [HslCommunication.Profinet.Knx](e9ca9cc4-27ce-cc6f-ea0b-badd3d84fe43.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public KnxUdp()
```

```
Public Sub New
```

```
public:
KnxUdp()
```

```
new : unit -> KnxUdp
```

![](../icons/SectionExpanded.png)参见

#### 引用

[KnxUdp 类](17f9544a-78c3-1c9f-72df-11e3d0fd6606.htm)

[HslCommunication.Profinet.Knx 命名空间](e9ca9cc4-27ce-cc6f-ea0b-badd3d84fe43.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## KnxUdp 属性

[原文連結](http://api.hslcommunication.cn/html/3816a091-7361-922a-40c7-cb8a598ad4a1.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Knx](../html/e9ca9cc4-27ce-cc6f-ea0b-badd3d84fe43.htm "HslCommunication.Profinet.Knx")

[KnxUdp 类](../html/17f9544a-78c3-1c9f-72df-11e3d0fd6606.htm "KnxUdp 类")

[KnxUdp 属性](../html/3816a091-7361-922a-40c7-cb8a598ad4a1.htm "KnxUdp 属性")

[Channel 属性](../html/c03c4e81-e0c6-45b3-c5b6-990a7530174f.htm "Channel 属性 ")

[IsConnect 属性](../html/a1247873-ea2b-4177-6366-d9808e8adbdf.htm "IsConnect 属性 ")

[KnxCode 属性](../html/1f3bb20c-1084-d8cb-5080-5923b57573d0.htm "KnxCode 属性 ")

[LocalEndpoint 属性](../html/b7e4b069-6bf2-6d8e-6502-e04a60cec773.htm "LocalEndpoint 属性 ")

[LogNet 属性](../html/af906e51-4f30-5076-0933-3e822a48b96f.htm "LogNet 属性 ")

[RouEndpoint 属性](../html/a20342d5-ed6a-7803-931d-bbad42b7a1dc.htm "RouEndpoint 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| KnxUdp 属性 |

[KnxUdp](17f9544a-78c3-1c9f-72df-11e3d0fd6606.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [Channel](c03c4e81-e0c6-45b3-c5b6-990a7530174f.htm) | 通道号（由设备发来） |
| 公共属性 | [IsConnect](a1247873-ea2b-4177-6366-d9808e8adbdf.htm) | 当前的状态是否连接中 |
| 公共属性 | [KnxCode](1f3bb20c-1084-d8cb-5080-5923b57573d0.htm) | 通信指令类 |
| 公共属性 | [LocalEndpoint](b7e4b069-6bf2-6d8e-6502-e04a60cec773.htm) | 本机IP地址 |
| 公共属性 | [LogNet](af906e51-4f30-5076-0933-3e822a48b96f.htm) | 系统的日志信息 |
| 公共属性 | [RouEndpoint](a20342d5-ed6a-7803-931d-bbad42b7a1dc.htm) | 远程ip地址 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[KnxUdp 类](17f9544a-78c3-1c9f-72df-11e3d0fd6606.htm)

[HslCommunication.Profinet.Knx 命名空间](e9ca9cc4-27ce-cc6f-ea0b-badd3d84fe43.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Channel 属性 

[原文連結](http://api.hslcommunication.cn/html/c03c4e81-e0c6-45b3-c5b6-990a7530174f.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Knx](../html/e9ca9cc4-27ce-cc6f-ea0b-badd3d84fe43.htm "HslCommunication.Profinet.Knx")

[KnxUdp 类](../html/17f9544a-78c3-1c9f-72df-11e3d0fd6606.htm "KnxUdp 类")

[KnxUdp 属性](../html/3816a091-7361-922a-40c7-cb8a598ad4a1.htm "KnxUdp 属性")

[Channel 属性](../html/c03c4e81-e0c6-45b3-c5b6-990a7530174f.htm "Channel 属性 ")

[IsConnect 属性](../html/a1247873-ea2b-4177-6366-d9808e8adbdf.htm "IsConnect 属性 ")

[KnxCode 属性](../html/1f3bb20c-1084-d8cb-5080-5923b57573d0.htm "KnxCode 属性 ")

[LocalEndpoint 属性](../html/b7e4b069-6bf2-6d8e-6502-e04a60cec773.htm "LocalEndpoint 属性 ")

[LogNet 属性](../html/af906e51-4f30-5076-0933-3e822a48b96f.htm "LogNet 属性 ")

[RouEndpoint 属性](../html/a20342d5-ed6a-7803-931d-bbad42b7a1dc.htm "RouEndpoint 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| KnxUdpChannel 属性 |

通道号（由设备发来）

**命名空间：**
 [HslCommunication.Profinet.Knx](e9ca9cc4-27ce-cc6f-ea0b-badd3d84fe43.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public byte Channel { get; set; }
```

```
Public Property Channel As Byte
	Get
	Set
```

```
public:
property unsigned char Channel {
	unsigned char get ();
	void set (unsigned char value);
}
```

```
member Channel : byte with get, set
```

#### 属性值

类型：Byte

![](../icons/SectionExpanded.png)参见

#### 引用

[KnxUdp 类](17f9544a-78c3-1c9f-72df-11e3d0fd6606.htm)

[HslCommunication.Profinet.Knx 命名空间](e9ca9cc4-27ce-cc6f-ea0b-badd3d84fe43.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## IsConnect 属性 

[原文連結](http://api.hslcommunication.cn/html/a1247873-ea2b-4177-6366-d9808e8adbdf.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Knx](../html/e9ca9cc4-27ce-cc6f-ea0b-badd3d84fe43.htm "HslCommunication.Profinet.Knx")

[KnxUdp 类](../html/17f9544a-78c3-1c9f-72df-11e3d0fd6606.htm "KnxUdp 类")

[KnxUdp 属性](../html/3816a091-7361-922a-40c7-cb8a598ad4a1.htm "KnxUdp 属性")

[Channel 属性](../html/c03c4e81-e0c6-45b3-c5b6-990a7530174f.htm "Channel 属性 ")

[IsConnect 属性](../html/a1247873-ea2b-4177-6366-d9808e8adbdf.htm "IsConnect 属性 ")

[KnxCode 属性](../html/1f3bb20c-1084-d8cb-5080-5923b57573d0.htm "KnxCode 属性 ")

[LocalEndpoint 属性](../html/b7e4b069-6bf2-6d8e-6502-e04a60cec773.htm "LocalEndpoint 属性 ")

[LogNet 属性](../html/af906e51-4f30-5076-0933-3e822a48b96f.htm "LogNet 属性 ")

[RouEndpoint 属性](../html/a20342d5-ed6a-7803-931d-bbad42b7a1dc.htm "RouEndpoint 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| KnxUdpIsConnect 属性 |

当前的状态是否连接中

**命名空间：**
 [HslCommunication.Profinet.Knx](e9ca9cc4-27ce-cc6f-ea0b-badd3d84fe43.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public bool IsConnect { get; }
```

```
Public ReadOnly Property IsConnect As Boolean
	Get
```

```
public:
property bool IsConnect {
	bool get ();
}
```

```
member IsConnect : bool with get
```

#### 属性值

类型：Boolean

![](../icons/SectionExpanded.png)参见

#### 引用

[KnxUdp 类](17f9544a-78c3-1c9f-72df-11e3d0fd6606.htm)

[HslCommunication.Profinet.Knx 命名空间](e9ca9cc4-27ce-cc6f-ea0b-badd3d84fe43.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## KnxCode 属性 

[原文連結](http://api.hslcommunication.cn/html/1f3bb20c-1084-d8cb-5080-5923b57573d0.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Knx](../html/e9ca9cc4-27ce-cc6f-ea0b-badd3d84fe43.htm "HslCommunication.Profinet.Knx")

[KnxUdp 类](../html/17f9544a-78c3-1c9f-72df-11e3d0fd6606.htm "KnxUdp 类")

[KnxUdp 属性](../html/3816a091-7361-922a-40c7-cb8a598ad4a1.htm "KnxUdp 属性")

[Channel 属性](../html/c03c4e81-e0c6-45b3-c5b6-990a7530174f.htm "Channel 属性 ")

[IsConnect 属性](../html/a1247873-ea2b-4177-6366-d9808e8adbdf.htm "IsConnect 属性 ")

[KnxCode 属性](../html/1f3bb20c-1084-d8cb-5080-5923b57573d0.htm "KnxCode 属性 ")

[LocalEndpoint 属性](../html/b7e4b069-6bf2-6d8e-6502-e04a60cec773.htm "LocalEndpoint 属性 ")

[LogNet 属性](../html/af906e51-4f30-5076-0933-3e822a48b96f.htm "LogNet 属性 ")

[RouEndpoint 属性](../html/a20342d5-ed6a-7803-931d-bbad42b7a1dc.htm "RouEndpoint 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| KnxUdpKnxCode 属性 |

通信指令类

**命名空间：**
 [HslCommunication.Profinet.Knx](e9ca9cc4-27ce-cc6f-ea0b-badd3d84fe43.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public KnxCode KnxCode { get; }
```

```
Public ReadOnly Property KnxCode As KnxCode
	Get
```

```
public:
property KnxCode^ KnxCode {
	KnxCode^ get ();
}
```

```
member KnxCode : KnxCode with get
```

#### 属性值

类型：[KnxCode](34d5ccbc-483b-c8ae-ce2d-c2c37b1fe2a3.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[KnxUdp 类](17f9544a-78c3-1c9f-72df-11e3d0fd6606.htm)

[HslCommunication.Profinet.Knx 命名空间](e9ca9cc4-27ce-cc6f-ea0b-badd3d84fe43.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## LocalEndpoint 属性 

[原文連結](http://api.hslcommunication.cn/html/b7e4b069-6bf2-6d8e-6502-e04a60cec773.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Knx](../html/e9ca9cc4-27ce-cc6f-ea0b-badd3d84fe43.htm "HslCommunication.Profinet.Knx")

[KnxUdp 类](../html/17f9544a-78c3-1c9f-72df-11e3d0fd6606.htm "KnxUdp 类")

[KnxUdp 属性](../html/3816a091-7361-922a-40c7-cb8a598ad4a1.htm "KnxUdp 属性")

[Channel 属性](../html/c03c4e81-e0c6-45b3-c5b6-990a7530174f.htm "Channel 属性 ")

[IsConnect 属性](../html/a1247873-ea2b-4177-6366-d9808e8adbdf.htm "IsConnect 属性 ")

[KnxCode 属性](../html/1f3bb20c-1084-d8cb-5080-5923b57573d0.htm "KnxCode 属性 ")

[LocalEndpoint 属性](../html/b7e4b069-6bf2-6d8e-6502-e04a60cec773.htm "LocalEndpoint 属性 ")

[LogNet 属性](../html/af906e51-4f30-5076-0933-3e822a48b96f.htm "LogNet 属性 ")

[RouEndpoint 属性](../html/a20342d5-ed6a-7803-931d-bbad42b7a1dc.htm "RouEndpoint 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| KnxUdpLocalEndpoint 属性 |

本机IP地址

**命名空间：**
 [HslCommunication.Profinet.Knx](e9ca9cc4-27ce-cc6f-ea0b-badd3d84fe43.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public IPEndPoint LocalEndpoint { get; set; }
```

```
Public Property LocalEndpoint As IPEndPoint
	Get
	Set
```

```
public:
property IPEndPoint^ LocalEndpoint {
	IPEndPoint^ get ();
	void set (IPEndPoint^ value);
}
```

```
member LocalEndpoint : IPEndPoint with get, set
```

#### 属性值

类型：IPEndPoint

![](../icons/SectionExpanded.png)参见

#### 引用

[KnxUdp 类](17f9544a-78c3-1c9f-72df-11e3d0fd6606.htm)

[HslCommunication.Profinet.Knx 命名空间](e9ca9cc4-27ce-cc6f-ea0b-badd3d84fe43.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## LogNet 属性 

[原文連結](http://api.hslcommunication.cn/html/af906e51-4f30-5076-0933-3e822a48b96f.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Knx](../html/e9ca9cc4-27ce-cc6f-ea0b-badd3d84fe43.htm "HslCommunication.Profinet.Knx")

[KnxUdp 类](../html/17f9544a-78c3-1c9f-72df-11e3d0fd6606.htm "KnxUdp 类")

[KnxUdp 属性](../html/3816a091-7361-922a-40c7-cb8a598ad4a1.htm "KnxUdp 属性")

[Channel 属性](../html/c03c4e81-e0c6-45b3-c5b6-990a7530174f.htm "Channel 属性 ")

[IsConnect 属性](../html/a1247873-ea2b-4177-6366-d9808e8adbdf.htm "IsConnect 属性 ")

[KnxCode 属性](../html/1f3bb20c-1084-d8cb-5080-5923b57573d0.htm "KnxCode 属性 ")

[LocalEndpoint 属性](../html/b7e4b069-6bf2-6d8e-6502-e04a60cec773.htm "LocalEndpoint 属性 ")

[LogNet 属性](../html/af906e51-4f30-5076-0933-3e822a48b96f.htm "LogNet 属性 ")

[RouEndpoint 属性](../html/a20342d5-ed6a-7803-931d-bbad42b7a1dc.htm "RouEndpoint 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| KnxUdpLogNet 属性 |

系统的日志信息

**命名空间：**
 [HslCommunication.Profinet.Knx](e9ca9cc4-27ce-cc6f-ea0b-badd3d84fe43.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public ILogNet LogNet { get; set; }
```

```
Public Property LogNet As ILogNet
	Get
	Set
```

```
public:
property ILogNet^ LogNet {
	ILogNet^ get ();
	void set (ILogNet^ value);
}
```

```
member LogNet : ILogNet with get, set
```

#### 属性值

类型：[ILogNet](d95d4704-db0d-b379-41e6-03879927a543.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[KnxUdp 类](17f9544a-78c3-1c9f-72df-11e3d0fd6606.htm)

[HslCommunication.Profinet.Knx 命名空间](e9ca9cc4-27ce-cc6f-ea0b-badd3d84fe43.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## RouEndpoint 属性 

[原文連結](http://api.hslcommunication.cn/html/a20342d5-ed6a-7803-931d-bbad42b7a1dc.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Knx](../html/e9ca9cc4-27ce-cc6f-ea0b-badd3d84fe43.htm "HslCommunication.Profinet.Knx")

[KnxUdp 类](../html/17f9544a-78c3-1c9f-72df-11e3d0fd6606.htm "KnxUdp 类")

[KnxUdp 属性](../html/3816a091-7361-922a-40c7-cb8a598ad4a1.htm "KnxUdp 属性")

[Channel 属性](../html/c03c4e81-e0c6-45b3-c5b6-990a7530174f.htm "Channel 属性 ")

[IsConnect 属性](../html/a1247873-ea2b-4177-6366-d9808e8adbdf.htm "IsConnect 属性 ")

[KnxCode 属性](../html/1f3bb20c-1084-d8cb-5080-5923b57573d0.htm "KnxCode 属性 ")

[LocalEndpoint 属性](../html/b7e4b069-6bf2-6d8e-6502-e04a60cec773.htm "LocalEndpoint 属性 ")

[LogNet 属性](../html/af906e51-4f30-5076-0933-3e822a48b96f.htm "LogNet 属性 ")

[RouEndpoint 属性](../html/a20342d5-ed6a-7803-931d-bbad42b7a1dc.htm "RouEndpoint 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| KnxUdpRouEndpoint 属性 |

远程ip地址

**命名空间：**
 [HslCommunication.Profinet.Knx](e9ca9cc4-27ce-cc6f-ea0b-badd3d84fe43.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public IPEndPoint RouEndpoint { get; set; }
```

```
Public Property RouEndpoint As IPEndPoint
	Get
	Set
```

```
public:
property IPEndPoint^ RouEndpoint {
	IPEndPoint^ get ();
	void set (IPEndPoint^ value);
}
```

```
member RouEndpoint : IPEndPoint with get, set
```

#### 属性值

类型：IPEndPoint

![](../icons/SectionExpanded.png)参见

#### 引用

[KnxUdp 类](17f9544a-78c3-1c9f-72df-11e3d0fd6606.htm)

[HslCommunication.Profinet.Knx 命名空间](e9ca9cc4-27ce-cc6f-ea0b-badd3d84fe43.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)