# HslCommunication - HslCommunication.Profinet.IDCard

> 分類頁數: 30



---
## HslCommunication.Profinet.IDCard

[原文連結](http://api.hslcommunication.cn/html/54fb02ce-3bfd-9c1d-96f8-511375556490.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.IDCard](../html/54fb02ce-3bfd-9c1d-96f8-511375556490.htm "HslCommunication.Profinet.IDCard")

[IdentityCard 类](../html/596fe590-4a15-954f-47da-25d8386a8cde.htm "IdentityCard 类")

[SAMSerial 类](../html/48f6a433-86f0-a4f8-1751-462e12ad6923.htm "SAMSerial 类")

[SAMTcpNet 类](../html/e33ad4bf-ecaf-8484-5f34-e472f905402f.htm "SAMTcpNet 类")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslCommunication.Profinet.IDCard 命名空间 |

[缺少 "N:HslCommunication.Profinet.IDCard" 的 <summary> 文档]

![](../icons/SectionExpanded.png)类

|  | 类 | 说明 |
| --- | --- | --- |
| 公共类 | [IdentityCard](596fe590-4a15-954f-47da-25d8386a8cde.htm) | 身份证的信息类 |
| 公共类代码示例 | [SAMSerial](48f6a433-86f0-a4f8-1751-462e12ad6923.htm) | 基于SAM协议的串口通信类，支持读取身份证的数据信息，详细参见API文档  Network class implemented by Tcp based on the SAM protocol, which supports reading ID card data information, see API documentation for details |
| 公共类代码示例 | [SAMTcpNet](e33ad4bf-ecaf-8484-5f34-e472f905402f.htm) | 基于SAM协议的Tcp实现的网络类，支持读取身份证的数据信息，通过透传的形式实现，除了初始化和串口类不一致，调用方法是几乎一模一样的，详细参见API文档  The network class implemented by Tcp based on the SAM protocol supports reading ID card data information and is implemented in the form of transparent transmission. Except for the inconsistency between the initialization and the serial port class, the calling method is almost the same. See the API documentation for details |

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## IdentityCard 类

[原文連結](http://api.hslcommunication.cn/html/596fe590-4a15-954f-47da-25d8386a8cde.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.IDCard](../html/54fb02ce-3bfd-9c1d-96f8-511375556490.htm "HslCommunication.Profinet.IDCard")

[IdentityCard 类](../html/596fe590-4a15-954f-47da-25d8386a8cde.htm "IdentityCard 类")

[IdentityCard 构造函数](../html/13a00ca8-7b8a-ce44-0303-84c31355656a.htm "IdentityCard 构造函数 ")

[IdentityCard 属性](../html/a6854f05-9e49-fd16-f424-f9c494cd2366.htm "IdentityCard 属性")

[IdentityCard 方法](../html/3320b06f-433d-42d4-ca92-7fb6247c2668.htm "IdentityCard 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IdentityCard 类 |

身份证的信息类

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  HslCommunication.Profinet.IDCardIdentityCard

**命名空间：**
 [HslCommunication.Profinet.IDCard](54fb02ce-3bfd-9c1d-96f8-511375556490.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class IdentityCard
```

```
Public Class IdentityCard
```

```
public ref class IdentityCard
```

```
type IdentityCard =  class end
```

IdentityCard 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [IdentityCard](13a00ca8-7b8a-ce44-0303-84c31355656a.htm) | 初始化 IdentityCard 类的一个新实例 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [Address](41a5131c-5317-fe1a-9043-6d7408869a9b.htm) | 地址 |
| 公共属性 | [Birthday](8f3ca62a-f3b1-716a-e34f-18f277479d82.htm) | 生日 |
| 公共属性 | [Id](e27bdaf4-fe99-f947-0698-bde89c3bf8ed.htm) | 身份证号 |
| 公共属性 | [Name](eb2adfaa-795d-7be7-3dda-8828eed3a8b6.htm) | 名字 |
| 公共属性 | [Nation](cbce8863-68df-3356-453d-d1eee371161d.htm) | 民族 |
| 公共属性 | [Organ](f58b4a79-0f09-507e-a9f7-ecfaa11b07ae.htm) | 发证机关 |
| 公共属性 | [Portrait](2d066987-4c10-478f-a23d-324f1ab322c4.htm) | 头像信息 |
| 公共属性 | [Sex](971c6b7b-1089-1ba3-71f6-1fcf4c6ed117.htm) | 性别 |
| 公共属性 | [ValidityEndDate](7792f6f5-d2af-35a3-be13-5214dcc953f1.htm) | 有效期日期的结束日期 |
| 公共属性 | [ValidityStartDate](28128ab6-34b5-67e4-e81f-c169bfcae121.htm) | 有效期日期的起始日期 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法 | [ToString](4c2ab48a-b5fa-2cb7-1ca5-e5f65404eee5.htm) | 返回表示当前对象的字符串 (重写 ObjectToString.) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.Profinet.IDCard 命名空间](54fb02ce-3bfd-9c1d-96f8-511375556490.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## IdentityCard 构造函数 

[原文連結](http://api.hslcommunication.cn/html/13a00ca8-7b8a-ce44-0303-84c31355656a.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.IDCard](../html/54fb02ce-3bfd-9c1d-96f8-511375556490.htm "HslCommunication.Profinet.IDCard")

[IdentityCard 类](../html/596fe590-4a15-954f-47da-25d8386a8cde.htm "IdentityCard 类")

[IdentityCard 构造函数](../html/13a00ca8-7b8a-ce44-0303-84c31355656a.htm "IdentityCard 构造函数 ")

[IdentityCard 属性](../html/a6854f05-9e49-fd16-f424-f9c494cd2366.htm "IdentityCard 属性")

[IdentityCard 方法](../html/3320b06f-433d-42d4-ca92-7fb6247c2668.htm "IdentityCard 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IdentityCard 构造函数 |

初始化 [IdentityCard](596fe590-4a15-954f-47da-25d8386a8cde.htm) 类的一个新实例

**命名空间：**
 [HslCommunication.Profinet.IDCard](54fb02ce-3bfd-9c1d-96f8-511375556490.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public IdentityCard()
```

```
Public Sub New
```

```
public:
IdentityCard()
```

```
new : unit -> IdentityCard
```

![](../icons/SectionExpanded.png)参见

#### 引用

[IdentityCard 类](596fe590-4a15-954f-47da-25d8386a8cde.htm)

[HslCommunication.Profinet.IDCard 命名空间](54fb02ce-3bfd-9c1d-96f8-511375556490.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## IdentityCard 属性

[原文連結](http://api.hslcommunication.cn/html/a6854f05-9e49-fd16-f424-f9c494cd2366.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.IDCard](../html/54fb02ce-3bfd-9c1d-96f8-511375556490.htm "HslCommunication.Profinet.IDCard")

[IdentityCard 类](../html/596fe590-4a15-954f-47da-25d8386a8cde.htm "IdentityCard 类")

[IdentityCard 属性](../html/a6854f05-9e49-fd16-f424-f9c494cd2366.htm "IdentityCard 属性")

[Address 属性](../html/41a5131c-5317-fe1a-9043-6d7408869a9b.htm "Address 属性 ")

[Birthday 属性](../html/8f3ca62a-f3b1-716a-e34f-18f277479d82.htm "Birthday 属性 ")

[Id 属性](../html/e27bdaf4-fe99-f947-0698-bde89c3bf8ed.htm "Id 属性 ")

[Name 属性](../html/eb2adfaa-795d-7be7-3dda-8828eed3a8b6.htm "Name 属性 ")

[Nation 属性](../html/cbce8863-68df-3356-453d-d1eee371161d.htm "Nation 属性 ")

[Organ 属性](../html/f58b4a79-0f09-507e-a9f7-ecfaa11b07ae.htm "Organ 属性 ")

[Portrait 属性](../html/2d066987-4c10-478f-a23d-324f1ab322c4.htm "Portrait 属性 ")

[Sex 属性](../html/971c6b7b-1089-1ba3-71f6-1fcf4c6ed117.htm "Sex 属性 ")

[ValidityEndDate 属性](../html/7792f6f5-d2af-35a3-be13-5214dcc953f1.htm "ValidityEndDate 属性 ")

[ValidityStartDate 属性](../html/28128ab6-34b5-67e4-e81f-c169bfcae121.htm "ValidityStartDate 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IdentityCard 属性 |

[IdentityCard](596fe590-4a15-954f-47da-25d8386a8cde.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [Address](41a5131c-5317-fe1a-9043-6d7408869a9b.htm) | 地址 |
| 公共属性 | [Birthday](8f3ca62a-f3b1-716a-e34f-18f277479d82.htm) | 生日 |
| 公共属性 | [Id](e27bdaf4-fe99-f947-0698-bde89c3bf8ed.htm) | 身份证号 |
| 公共属性 | [Name](eb2adfaa-795d-7be7-3dda-8828eed3a8b6.htm) | 名字 |
| 公共属性 | [Nation](cbce8863-68df-3356-453d-d1eee371161d.htm) | 民族 |
| 公共属性 | [Organ](f58b4a79-0f09-507e-a9f7-ecfaa11b07ae.htm) | 发证机关 |
| 公共属性 | [Portrait](2d066987-4c10-478f-a23d-324f1ab322c4.htm) | 头像信息 |
| 公共属性 | [Sex](971c6b7b-1089-1ba3-71f6-1fcf4c6ed117.htm) | 性别 |
| 公共属性 | [ValidityEndDate](7792f6f5-d2af-35a3-be13-5214dcc953f1.htm) | 有效期日期的结束日期 |
| 公共属性 | [ValidityStartDate](28128ab6-34b5-67e4-e81f-c169bfcae121.htm) | 有效期日期的起始日期 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[IdentityCard 类](596fe590-4a15-954f-47da-25d8386a8cde.htm)

[HslCommunication.Profinet.IDCard 命名空间](54fb02ce-3bfd-9c1d-96f8-511375556490.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Address 属性 

[原文連結](http://api.hslcommunication.cn/html/41a5131c-5317-fe1a-9043-6d7408869a9b.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.IDCard](../html/54fb02ce-3bfd-9c1d-96f8-511375556490.htm "HslCommunication.Profinet.IDCard")

[IdentityCard 类](../html/596fe590-4a15-954f-47da-25d8386a8cde.htm "IdentityCard 类")

[IdentityCard 属性](../html/a6854f05-9e49-fd16-f424-f9c494cd2366.htm "IdentityCard 属性")

[Address 属性](../html/41a5131c-5317-fe1a-9043-6d7408869a9b.htm "Address 属性 ")

[Birthday 属性](../html/8f3ca62a-f3b1-716a-e34f-18f277479d82.htm "Birthday 属性 ")

[Id 属性](../html/e27bdaf4-fe99-f947-0698-bde89c3bf8ed.htm "Id 属性 ")

[Name 属性](../html/eb2adfaa-795d-7be7-3dda-8828eed3a8b6.htm "Name 属性 ")

[Nation 属性](../html/cbce8863-68df-3356-453d-d1eee371161d.htm "Nation 属性 ")

[Organ 属性](../html/f58b4a79-0f09-507e-a9f7-ecfaa11b07ae.htm "Organ 属性 ")

[Portrait 属性](../html/2d066987-4c10-478f-a23d-324f1ab322c4.htm "Portrait 属性 ")

[Sex 属性](../html/971c6b7b-1089-1ba3-71f6-1fcf4c6ed117.htm "Sex 属性 ")

[ValidityEndDate 属性](../html/7792f6f5-d2af-35a3-be13-5214dcc953f1.htm "ValidityEndDate 属性 ")

[ValidityStartDate 属性](../html/28128ab6-34b5-67e4-e81f-c169bfcae121.htm "ValidityStartDate 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IdentityCardAddress 属性 |

地址

**命名空间：**
 [HslCommunication.Profinet.IDCard](54fb02ce-3bfd-9c1d-96f8-511375556490.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public string Address { get; set; }
```

```
Public Property Address As String
	Get
	Set
```

```
public:
property String^ Address {
	String^ get ();
	void set (String^ value);
}
```

```
member Address : string with get, set
```

#### 属性值

类型：String

![](../icons/SectionExpanded.png)参见

#### 引用

[IdentityCard 类](596fe590-4a15-954f-47da-25d8386a8cde.htm)

[HslCommunication.Profinet.IDCard 命名空间](54fb02ce-3bfd-9c1d-96f8-511375556490.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Birthday 属性 

[原文連結](http://api.hslcommunication.cn/html/8f3ca62a-f3b1-716a-e34f-18f277479d82.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.IDCard](../html/54fb02ce-3bfd-9c1d-96f8-511375556490.htm "HslCommunication.Profinet.IDCard")

[IdentityCard 类](../html/596fe590-4a15-954f-47da-25d8386a8cde.htm "IdentityCard 类")

[IdentityCard 属性](../html/a6854f05-9e49-fd16-f424-f9c494cd2366.htm "IdentityCard 属性")

[Address 属性](../html/41a5131c-5317-fe1a-9043-6d7408869a9b.htm "Address 属性 ")

[Birthday 属性](../html/8f3ca62a-f3b1-716a-e34f-18f277479d82.htm "Birthday 属性 ")

[Id 属性](../html/e27bdaf4-fe99-f947-0698-bde89c3bf8ed.htm "Id 属性 ")

[Name 属性](../html/eb2adfaa-795d-7be7-3dda-8828eed3a8b6.htm "Name 属性 ")

[Nation 属性](../html/cbce8863-68df-3356-453d-d1eee371161d.htm "Nation 属性 ")

[Organ 属性](../html/f58b4a79-0f09-507e-a9f7-ecfaa11b07ae.htm "Organ 属性 ")

[Portrait 属性](../html/2d066987-4c10-478f-a23d-324f1ab322c4.htm "Portrait 属性 ")

[Sex 属性](../html/971c6b7b-1089-1ba3-71f6-1fcf4c6ed117.htm "Sex 属性 ")

[ValidityEndDate 属性](../html/7792f6f5-d2af-35a3-be13-5214dcc953f1.htm "ValidityEndDate 属性 ")

[ValidityStartDate 属性](../html/28128ab6-34b5-67e4-e81f-c169bfcae121.htm "ValidityStartDate 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IdentityCardBirthday 属性 |

生日

**命名空间：**
 [HslCommunication.Profinet.IDCard](54fb02ce-3bfd-9c1d-96f8-511375556490.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public DateTime Birthday { get; set; }
```

```
Public Property Birthday As DateTime
	Get
	Set
```

```
public:
property DateTime Birthday {
	DateTime get ();
	void set (DateTime value);
}
```

```
member Birthday : DateTime with get, set
```

#### 属性值

类型：DateTime

![](../icons/SectionExpanded.png)参见

#### 引用

[IdentityCard 类](596fe590-4a15-954f-47da-25d8386a8cde.htm)

[HslCommunication.Profinet.IDCard 命名空间](54fb02ce-3bfd-9c1d-96f8-511375556490.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Id 属性 

[原文連結](http://api.hslcommunication.cn/html/e27bdaf4-fe99-f947-0698-bde89c3bf8ed.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.IDCard](../html/54fb02ce-3bfd-9c1d-96f8-511375556490.htm "HslCommunication.Profinet.IDCard")

[IdentityCard 类](../html/596fe590-4a15-954f-47da-25d8386a8cde.htm "IdentityCard 类")

[IdentityCard 属性](../html/a6854f05-9e49-fd16-f424-f9c494cd2366.htm "IdentityCard 属性")

[Address 属性](../html/41a5131c-5317-fe1a-9043-6d7408869a9b.htm "Address 属性 ")

[Birthday 属性](../html/8f3ca62a-f3b1-716a-e34f-18f277479d82.htm "Birthday 属性 ")

[Id 属性](../html/e27bdaf4-fe99-f947-0698-bde89c3bf8ed.htm "Id 属性 ")

[Name 属性](../html/eb2adfaa-795d-7be7-3dda-8828eed3a8b6.htm "Name 属性 ")

[Nation 属性](../html/cbce8863-68df-3356-453d-d1eee371161d.htm "Nation 属性 ")

[Organ 属性](../html/f58b4a79-0f09-507e-a9f7-ecfaa11b07ae.htm "Organ 属性 ")

[Portrait 属性](../html/2d066987-4c10-478f-a23d-324f1ab322c4.htm "Portrait 属性 ")

[Sex 属性](../html/971c6b7b-1089-1ba3-71f6-1fcf4c6ed117.htm "Sex 属性 ")

[ValidityEndDate 属性](../html/7792f6f5-d2af-35a3-be13-5214dcc953f1.htm "ValidityEndDate 属性 ")

[ValidityStartDate 属性](../html/28128ab6-34b5-67e4-e81f-c169bfcae121.htm "ValidityStartDate 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IdentityCardId 属性 |

身份证号

**命名空间：**
 [HslCommunication.Profinet.IDCard](54fb02ce-3bfd-9c1d-96f8-511375556490.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public string Id { get; set; }
```

```
Public Property Id As String
	Get
	Set
```

```
public:
property String^ Id {
	String^ get ();
	void set (String^ value);
}
```

```
member Id : string with get, set
```

#### 属性值

类型：String

![](../icons/SectionExpanded.png)参见

#### 引用

[IdentityCard 类](596fe590-4a15-954f-47da-25d8386a8cde.htm)

[HslCommunication.Profinet.IDCard 命名空间](54fb02ce-3bfd-9c1d-96f8-511375556490.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Name 属性 

[原文連結](http://api.hslcommunication.cn/html/eb2adfaa-795d-7be7-3dda-8828eed3a8b6.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.IDCard](../html/54fb02ce-3bfd-9c1d-96f8-511375556490.htm "HslCommunication.Profinet.IDCard")

[IdentityCard 类](../html/596fe590-4a15-954f-47da-25d8386a8cde.htm "IdentityCard 类")

[IdentityCard 属性](../html/a6854f05-9e49-fd16-f424-f9c494cd2366.htm "IdentityCard 属性")

[Address 属性](../html/41a5131c-5317-fe1a-9043-6d7408869a9b.htm "Address 属性 ")

[Birthday 属性](../html/8f3ca62a-f3b1-716a-e34f-18f277479d82.htm "Birthday 属性 ")

[Id 属性](../html/e27bdaf4-fe99-f947-0698-bde89c3bf8ed.htm "Id 属性 ")

[Name 属性](../html/eb2adfaa-795d-7be7-3dda-8828eed3a8b6.htm "Name 属性 ")

[Nation 属性](../html/cbce8863-68df-3356-453d-d1eee371161d.htm "Nation 属性 ")

[Organ 属性](../html/f58b4a79-0f09-507e-a9f7-ecfaa11b07ae.htm "Organ 属性 ")

[Portrait 属性](../html/2d066987-4c10-478f-a23d-324f1ab322c4.htm "Portrait 属性 ")

[Sex 属性](../html/971c6b7b-1089-1ba3-71f6-1fcf4c6ed117.htm "Sex 属性 ")

[ValidityEndDate 属性](../html/7792f6f5-d2af-35a3-be13-5214dcc953f1.htm "ValidityEndDate 属性 ")

[ValidityStartDate 属性](../html/28128ab6-34b5-67e4-e81f-c169bfcae121.htm "ValidityStartDate 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IdentityCardName 属性 |

名字

**命名空间：**
 [HslCommunication.Profinet.IDCard](54fb02ce-3bfd-9c1d-96f8-511375556490.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public string Name { get; set; }
```

```
Public Property Name As String
	Get
	Set
```

```
public:
property String^ Name {
	String^ get ();
	void set (String^ value);
}
```

```
member Name : string with get, set
```

#### 属性值

类型：String

![](../icons/SectionExpanded.png)参见

#### 引用

[IdentityCard 类](596fe590-4a15-954f-47da-25d8386a8cde.htm)

[HslCommunication.Profinet.IDCard 命名空间](54fb02ce-3bfd-9c1d-96f8-511375556490.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Nation 属性 

[原文連結](http://api.hslcommunication.cn/html/cbce8863-68df-3356-453d-d1eee371161d.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.IDCard](../html/54fb02ce-3bfd-9c1d-96f8-511375556490.htm "HslCommunication.Profinet.IDCard")

[IdentityCard 类](../html/596fe590-4a15-954f-47da-25d8386a8cde.htm "IdentityCard 类")

[IdentityCard 属性](../html/a6854f05-9e49-fd16-f424-f9c494cd2366.htm "IdentityCard 属性")

[Address 属性](../html/41a5131c-5317-fe1a-9043-6d7408869a9b.htm "Address 属性 ")

[Birthday 属性](../html/8f3ca62a-f3b1-716a-e34f-18f277479d82.htm "Birthday 属性 ")

[Id 属性](../html/e27bdaf4-fe99-f947-0698-bde89c3bf8ed.htm "Id 属性 ")

[Name 属性](../html/eb2adfaa-795d-7be7-3dda-8828eed3a8b6.htm "Name 属性 ")

[Nation 属性](../html/cbce8863-68df-3356-453d-d1eee371161d.htm "Nation 属性 ")

[Organ 属性](../html/f58b4a79-0f09-507e-a9f7-ecfaa11b07ae.htm "Organ 属性 ")

[Portrait 属性](../html/2d066987-4c10-478f-a23d-324f1ab322c4.htm "Portrait 属性 ")

[Sex 属性](../html/971c6b7b-1089-1ba3-71f6-1fcf4c6ed117.htm "Sex 属性 ")

[ValidityEndDate 属性](../html/7792f6f5-d2af-35a3-be13-5214dcc953f1.htm "ValidityEndDate 属性 ")

[ValidityStartDate 属性](../html/28128ab6-34b5-67e4-e81f-c169bfcae121.htm "ValidityStartDate 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IdentityCardNation 属性 |

民族

**命名空间：**
 [HslCommunication.Profinet.IDCard](54fb02ce-3bfd-9c1d-96f8-511375556490.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public string Nation { get; set; }
```

```
Public Property Nation As String
	Get
	Set
```

```
public:
property String^ Nation {
	String^ get ();
	void set (String^ value);
}
```

```
member Nation : string with get, set
```

#### 属性值

类型：String

![](../icons/SectionExpanded.png)参见

#### 引用

[IdentityCard 类](596fe590-4a15-954f-47da-25d8386a8cde.htm)

[HslCommunication.Profinet.IDCard 命名空间](54fb02ce-3bfd-9c1d-96f8-511375556490.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Organ 属性 

[原文連結](http://api.hslcommunication.cn/html/f58b4a79-0f09-507e-a9f7-ecfaa11b07ae.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.IDCard](../html/54fb02ce-3bfd-9c1d-96f8-511375556490.htm "HslCommunication.Profinet.IDCard")

[IdentityCard 类](../html/596fe590-4a15-954f-47da-25d8386a8cde.htm "IdentityCard 类")

[IdentityCard 属性](../html/a6854f05-9e49-fd16-f424-f9c494cd2366.htm "IdentityCard 属性")

[Address 属性](../html/41a5131c-5317-fe1a-9043-6d7408869a9b.htm "Address 属性 ")

[Birthday 属性](../html/8f3ca62a-f3b1-716a-e34f-18f277479d82.htm "Birthday 属性 ")

[Id 属性](../html/e27bdaf4-fe99-f947-0698-bde89c3bf8ed.htm "Id 属性 ")

[Name 属性](../html/eb2adfaa-795d-7be7-3dda-8828eed3a8b6.htm "Name 属性 ")

[Nation 属性](../html/cbce8863-68df-3356-453d-d1eee371161d.htm "Nation 属性 ")

[Organ 属性](../html/f58b4a79-0f09-507e-a9f7-ecfaa11b07ae.htm "Organ 属性 ")

[Portrait 属性](../html/2d066987-4c10-478f-a23d-324f1ab322c4.htm "Portrait 属性 ")

[Sex 属性](../html/971c6b7b-1089-1ba3-71f6-1fcf4c6ed117.htm "Sex 属性 ")

[ValidityEndDate 属性](../html/7792f6f5-d2af-35a3-be13-5214dcc953f1.htm "ValidityEndDate 属性 ")

[ValidityStartDate 属性](../html/28128ab6-34b5-67e4-e81f-c169bfcae121.htm "ValidityStartDate 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IdentityCardOrgan 属性 |

发证机关

**命名空间：**
 [HslCommunication.Profinet.IDCard](54fb02ce-3bfd-9c1d-96f8-511375556490.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public string Organ { get; set; }
```

```
Public Property Organ As String
	Get
	Set
```

```
public:
property String^ Organ {
	String^ get ();
	void set (String^ value);
}
```

```
member Organ : string with get, set
```

#### 属性值

类型：String

![](../icons/SectionExpanded.png)参见

#### 引用

[IdentityCard 类](596fe590-4a15-954f-47da-25d8386a8cde.htm)

[HslCommunication.Profinet.IDCard 命名空间](54fb02ce-3bfd-9c1d-96f8-511375556490.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Portrait 属性 

[原文連結](http://api.hslcommunication.cn/html/2d066987-4c10-478f-a23d-324f1ab322c4.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.IDCard](../html/54fb02ce-3bfd-9c1d-96f8-511375556490.htm "HslCommunication.Profinet.IDCard")

[IdentityCard 类](../html/596fe590-4a15-954f-47da-25d8386a8cde.htm "IdentityCard 类")

[IdentityCard 属性](../html/a6854f05-9e49-fd16-f424-f9c494cd2366.htm "IdentityCard 属性")

[Address 属性](../html/41a5131c-5317-fe1a-9043-6d7408869a9b.htm "Address 属性 ")

[Birthday 属性](../html/8f3ca62a-f3b1-716a-e34f-18f277479d82.htm "Birthday 属性 ")

[Id 属性](../html/e27bdaf4-fe99-f947-0698-bde89c3bf8ed.htm "Id 属性 ")

[Name 属性](../html/eb2adfaa-795d-7be7-3dda-8828eed3a8b6.htm "Name 属性 ")

[Nation 属性](../html/cbce8863-68df-3356-453d-d1eee371161d.htm "Nation 属性 ")

[Organ 属性](../html/f58b4a79-0f09-507e-a9f7-ecfaa11b07ae.htm "Organ 属性 ")

[Portrait 属性](../html/2d066987-4c10-478f-a23d-324f1ab322c4.htm "Portrait 属性 ")

[Sex 属性](../html/971c6b7b-1089-1ba3-71f6-1fcf4c6ed117.htm "Sex 属性 ")

[ValidityEndDate 属性](../html/7792f6f5-d2af-35a3-be13-5214dcc953f1.htm "ValidityEndDate 属性 ")

[ValidityStartDate 属性](../html/28128ab6-34b5-67e4-e81f-c169bfcae121.htm "ValidityStartDate 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IdentityCardPortrait 属性 |

头像信息

**命名空间：**
 [HslCommunication.Profinet.IDCard](54fb02ce-3bfd-9c1d-96f8-511375556490.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public byte[] Portrait { get; set; }
```

```
Public Property Portrait As Byte()
	Get
	Set
```

```
public:
property array<unsigned char>^ Portrait {
	array<unsigned char>^ get ();
	void set (array<unsigned char>^ value);
}
```

```
member Portrait : byte[] with get, set
```

#### 属性值

类型：Byte

![](../icons/SectionExpanded.png)参见

#### 引用

[IdentityCard 类](596fe590-4a15-954f-47da-25d8386a8cde.htm)

[HslCommunication.Profinet.IDCard 命名空间](54fb02ce-3bfd-9c1d-96f8-511375556490.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Sex 属性 

[原文連結](http://api.hslcommunication.cn/html/971c6b7b-1089-1ba3-71f6-1fcf4c6ed117.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.IDCard](../html/54fb02ce-3bfd-9c1d-96f8-511375556490.htm "HslCommunication.Profinet.IDCard")

[IdentityCard 类](../html/596fe590-4a15-954f-47da-25d8386a8cde.htm "IdentityCard 类")

[IdentityCard 属性](../html/a6854f05-9e49-fd16-f424-f9c494cd2366.htm "IdentityCard 属性")

[Address 属性](../html/41a5131c-5317-fe1a-9043-6d7408869a9b.htm "Address 属性 ")

[Birthday 属性](../html/8f3ca62a-f3b1-716a-e34f-18f277479d82.htm "Birthday 属性 ")

[Id 属性](../html/e27bdaf4-fe99-f947-0698-bde89c3bf8ed.htm "Id 属性 ")

[Name 属性](../html/eb2adfaa-795d-7be7-3dda-8828eed3a8b6.htm "Name 属性 ")

[Nation 属性](../html/cbce8863-68df-3356-453d-d1eee371161d.htm "Nation 属性 ")

[Organ 属性](../html/f58b4a79-0f09-507e-a9f7-ecfaa11b07ae.htm "Organ 属性 ")

[Portrait 属性](../html/2d066987-4c10-478f-a23d-324f1ab322c4.htm "Portrait 属性 ")

[Sex 属性](../html/971c6b7b-1089-1ba3-71f6-1fcf4c6ed117.htm "Sex 属性 ")

[ValidityEndDate 属性](../html/7792f6f5-d2af-35a3-be13-5214dcc953f1.htm "ValidityEndDate 属性 ")

[ValidityStartDate 属性](../html/28128ab6-34b5-67e4-e81f-c169bfcae121.htm "ValidityStartDate 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IdentityCardSex 属性 |

性别

**命名空间：**
 [HslCommunication.Profinet.IDCard](54fb02ce-3bfd-9c1d-96f8-511375556490.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public string Sex { get; set; }
```

```
Public Property Sex As String
	Get
	Set
```

```
public:
property String^ Sex {
	String^ get ();
	void set (String^ value);
}
```

```
member Sex : string with get, set
```

#### 属性值

类型：String

![](../icons/SectionExpanded.png)参见

#### 引用

[IdentityCard 类](596fe590-4a15-954f-47da-25d8386a8cde.htm)

[HslCommunication.Profinet.IDCard 命名空间](54fb02ce-3bfd-9c1d-96f8-511375556490.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ValidityEndDate 属性 

[原文連結](http://api.hslcommunication.cn/html/7792f6f5-d2af-35a3-be13-5214dcc953f1.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.IDCard](../html/54fb02ce-3bfd-9c1d-96f8-511375556490.htm "HslCommunication.Profinet.IDCard")

[IdentityCard 类](../html/596fe590-4a15-954f-47da-25d8386a8cde.htm "IdentityCard 类")

[IdentityCard 属性](../html/a6854f05-9e49-fd16-f424-f9c494cd2366.htm "IdentityCard 属性")

[Address 属性](../html/41a5131c-5317-fe1a-9043-6d7408869a9b.htm "Address 属性 ")

[Birthday 属性](../html/8f3ca62a-f3b1-716a-e34f-18f277479d82.htm "Birthday 属性 ")

[Id 属性](../html/e27bdaf4-fe99-f947-0698-bde89c3bf8ed.htm "Id 属性 ")

[Name 属性](../html/eb2adfaa-795d-7be7-3dda-8828eed3a8b6.htm "Name 属性 ")

[Nation 属性](../html/cbce8863-68df-3356-453d-d1eee371161d.htm "Nation 属性 ")

[Organ 属性](../html/f58b4a79-0f09-507e-a9f7-ecfaa11b07ae.htm "Organ 属性 ")

[Portrait 属性](../html/2d066987-4c10-478f-a23d-324f1ab322c4.htm "Portrait 属性 ")

[Sex 属性](../html/971c6b7b-1089-1ba3-71f6-1fcf4c6ed117.htm "Sex 属性 ")

[ValidityEndDate 属性](../html/7792f6f5-d2af-35a3-be13-5214dcc953f1.htm "ValidityEndDate 属性 ")

[ValidityStartDate 属性](../html/28128ab6-34b5-67e4-e81f-c169bfcae121.htm "ValidityStartDate 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IdentityCardValidityEndDate 属性 |

有效期日期的结束日期

**命名空间：**
 [HslCommunication.Profinet.IDCard](54fb02ce-3bfd-9c1d-96f8-511375556490.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public DateTime ValidityEndDate { get; set; }
```

```
Public Property ValidityEndDate As DateTime
	Get
	Set
```

```
public:
property DateTime ValidityEndDate {
	DateTime get ();
	void set (DateTime value);
}
```

```
member ValidityEndDate : DateTime with get, set
```

#### 属性值

类型：DateTime

![](../icons/SectionExpanded.png)参见

#### 引用

[IdentityCard 类](596fe590-4a15-954f-47da-25d8386a8cde.htm)

[HslCommunication.Profinet.IDCard 命名空间](54fb02ce-3bfd-9c1d-96f8-511375556490.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ValidityStartDate 属性 

[原文連結](http://api.hslcommunication.cn/html/28128ab6-34b5-67e4-e81f-c169bfcae121.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.IDCard](../html/54fb02ce-3bfd-9c1d-96f8-511375556490.htm "HslCommunication.Profinet.IDCard")

[IdentityCard 类](../html/596fe590-4a15-954f-47da-25d8386a8cde.htm "IdentityCard 类")

[IdentityCard 属性](../html/a6854f05-9e49-fd16-f424-f9c494cd2366.htm "IdentityCard 属性")

[Address 属性](../html/41a5131c-5317-fe1a-9043-6d7408869a9b.htm "Address 属性 ")

[Birthday 属性](../html/8f3ca62a-f3b1-716a-e34f-18f277479d82.htm "Birthday 属性 ")

[Id 属性](../html/e27bdaf4-fe99-f947-0698-bde89c3bf8ed.htm "Id 属性 ")

[Name 属性](../html/eb2adfaa-795d-7be7-3dda-8828eed3a8b6.htm "Name 属性 ")

[Nation 属性](../html/cbce8863-68df-3356-453d-d1eee371161d.htm "Nation 属性 ")

[Organ 属性](../html/f58b4a79-0f09-507e-a9f7-ecfaa11b07ae.htm "Organ 属性 ")

[Portrait 属性](../html/2d066987-4c10-478f-a23d-324f1ab322c4.htm "Portrait 属性 ")

[Sex 属性](../html/971c6b7b-1089-1ba3-71f6-1fcf4c6ed117.htm "Sex 属性 ")

[ValidityEndDate 属性](../html/7792f6f5-d2af-35a3-be13-5214dcc953f1.htm "ValidityEndDate 属性 ")

[ValidityStartDate 属性](../html/28128ab6-34b5-67e4-e81f-c169bfcae121.htm "ValidityStartDate 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IdentityCardValidityStartDate 属性 |

有效期日期的起始日期

**命名空间：**
 [HslCommunication.Profinet.IDCard](54fb02ce-3bfd-9c1d-96f8-511375556490.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public DateTime ValidityStartDate { get; set; }
```

```
Public Property ValidityStartDate As DateTime
	Get
	Set
```

```
public:
property DateTime ValidityStartDate {
	DateTime get ();
	void set (DateTime value);
}
```

```
member ValidityStartDate : DateTime with get, set
```

#### 属性值

类型：DateTime

![](../icons/SectionExpanded.png)参见

#### 引用

[IdentityCard 类](596fe590-4a15-954f-47da-25d8386a8cde.htm)

[HslCommunication.Profinet.IDCard 命名空间](54fb02ce-3bfd-9c1d-96f8-511375556490.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## IdentityCard 方法

[原文連結](http://api.hslcommunication.cn/html/3320b06f-433d-42d4-ca92-7fb6247c2668.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.IDCard](../html/54fb02ce-3bfd-9c1d-96f8-511375556490.htm "HslCommunication.Profinet.IDCard")

[IdentityCard 类](../html/596fe590-4a15-954f-47da-25d8386a8cde.htm "IdentityCard 类")

[IdentityCard 方法](../html/3320b06f-433d-42d4-ca92-7fb6247c2668.htm "IdentityCard 方法")

[ToString 方法](../html/4c2ab48a-b5fa-2cb7-1ca5-e5f65404eee5.htm "ToString 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IdentityCard 方法 |

[IdentityCard](596fe590-4a15-954f-47da-25d8386a8cde.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法 | [ToString](4c2ab48a-b5fa-2cb7-1ca5-e5f65404eee5.htm) | 返回表示当前对象的字符串 (重写 ObjectToString.) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[IdentityCard 类](596fe590-4a15-954f-47da-25d8386a8cde.htm)

[HslCommunication.Profinet.IDCard 命名空间](54fb02ce-3bfd-9c1d-96f8-511375556490.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ToString 方法 

[原文連結](http://api.hslcommunication.cn/html/4c2ab48a-b5fa-2cb7-1ca5-e5f65404eee5.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.IDCard](../html/54fb02ce-3bfd-9c1d-96f8-511375556490.htm "HslCommunication.Profinet.IDCard")

[IdentityCard 类](../html/596fe590-4a15-954f-47da-25d8386a8cde.htm "IdentityCard 类")

[IdentityCard 方法](../html/3320b06f-433d-42d4-ca92-7fb6247c2668.htm "IdentityCard 方法")

[ToString 方法](../html/4c2ab48a-b5fa-2cb7-1ca5-e5f65404eee5.htm "ToString 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| IdentityCardToString 方法 |

返回表示当前对象的字符串

**命名空间：**
 [HslCommunication.Profinet.IDCard](54fb02ce-3bfd-9c1d-96f8-511375556490.htm)  
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
字符串

![](../icons/SectionExpanded.png)参见

#### 引用

[IdentityCard 类](596fe590-4a15-954f-47da-25d8386a8cde.htm)

[HslCommunication.Profinet.IDCard 命名空间](54fb02ce-3bfd-9c1d-96f8-511375556490.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## SAMSerial 类

[原文連結](http://api.hslcommunication.cn/html/48f6a433-86f0-a4f8-1751-462e12ad6923.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.IDCard](../html/54fb02ce-3bfd-9c1d-96f8-511375556490.htm "HslCommunication.Profinet.IDCard")

[SAMSerial 类](../html/48f6a433-86f0-a4f8-1751-462e12ad6923.htm "SAMSerial 类")

[SAMSerial 构造函数](../html/60865a16-8442-3d5c-1e81-24534fa4c1da.htm "SAMSerial 构造函数 ")

[SAMSerial 属性](../html/e37df3e7-0468-6303-0143-d7558a731a5d.htm "SAMSerial 属性")

[SAMSerial 方法](../html/3ee3795f-b4e0-f2ff-9050-7961204b2d49.htm "SAMSerial 方法")

[SAMSerial 字段](../html/dd810ad3-9acc-1696-cc05-6a68d6116b8b.htm "SAMSerial 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| SAMSerial 类 |

基于SAM协议的串口通信类，支持读取身份证的数据信息，详细参见API文档  
Network class implemented by Tcp based on the SAM protocol, which supports reading ID card data information,
see API documentation for details

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  [HslCommunication.Core.NetBinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)  
    [HslCommunication.SerialSerialBase](e877e20b-ef65-266c-05a5-5db58e6bad0c.htm)  
      HslCommunication.Profinet.IDCardSAMSerial

**命名空间：**
 [HslCommunication.Profinet.IDCard](54fb02ce-3bfd-9c1d-96f8-511375556490.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class SAMSerial : SerialBase
```

```
Public Class SAMSerial
	Inherits SerialBase
```

```
public ref class SAMSerial : public SerialBase
```

```
type SAMSerial =  
    class
        inherit SerialBase
    end
```

SAMSerial 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [SAMSerial](60865a16-8442-3d5c-1e81-24534fa4c1da.htm) | 实例化一个默认的对象  Instantiate a default object |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [BaudRate](1c33c44a-d2b9-badd-df76-09e155d26cc1.htm) | 当前连接串口信息的波特率  Baud rate of current connection serial port information (继承自 [SerialBase](e877e20b-ef65-266c-05a5-5db58e6bad0c.htm)。) |
| 公共属性 | [CommunicationPipe](cb528f9b-fe2c-5ac7-6912-2c3c5c1fdf2a.htm) | 获取或设置当前的管道信息，管道类型为[CommunicationPipe](ad8e903b-d289-21b0-077e-15ede8005221.htm)的继承类，内置了[PipeTcpNet](26b94bd5-2695-36e7-6378-c3d41846be8f.htm)管道，[PipeUdpNet](5a0b810e-187e-826b-0a13-34fd5576aa0e.htm)管道，[PipeSerialPort](9d924607-3dfc-3916-b532-c5dec48fe347.htm)管道等  Get or set the current pipeline information, the pipeline type is [CommunicationPipe](ad8e903b-d289-21b0-077e-15ede8005221.htm) inheritance class, [PipeTcpNet](26b94bd5-2695-36e7-6378-c3d41846be8f.htm) pipeline, [PipeUdpNet](5a0b810e-187e-826b-0a13-34fd5576aa0e.htm) pipeline, [PipeSerialPort](9d924607-3dfc-3916-b532-c5dec48fe347.htm) pipeline, etc (继承自 [SerialBase](e877e20b-ef65-266c-05a5-5db58e6bad0c.htm)。) |
| 公共属性 | [ConnectionId](ce1efb4b-fbc9-4683-163e-0a8e71e99d6b.htm) | 当前连接的唯一ID号，默认为长度20的guid码加随机数组成，方便列表管理，也可以自己指定  The unique ID number of the current connection. The default is a 20-digit guid code plus a random number. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [IsClearCacheBeforeRead](cf2aaeff-e4eb-849a-4783-d10e5e965515.htm) | 是否在发送数据前清空缓冲数据，默认是false  Whether to empty the buffer before sending data, the default is false (继承自 [SerialBase](e877e20b-ef65-266c-05a5-5db58e6bad0c.htm)。) |
| 公共属性代码示例 | [LogNet](8c56f246-88ac-afb8-8c06-972fafe4f78a.htm) | 组件的日志工具，支持日志记录，只要实例化后，当前网络的基本信息，就以[DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm)等级进行输出  The component's logging tool supports logging. As long as the instantiation of the basic network information, the output will be output at [DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm) (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [PortName](250b8746-0bb7-4caf-2a43-11109f2bda20.htm) | 当前连接串口信息的端口号名称  The port name of the current connection serial port information (继承自 [SerialBase](e877e20b-ef65-266c-05a5-5db58e6bad0c.htm)。) |
| 公共属性 | [ReceiveEmptyDataCount](480b2206-1de4-6a74-6f8e-7c30a2ca05f9.htm) | 获取或设置连续接收空的数据次数，在数据接收完成时有效，每个单位消耗的时间为[SleepTime](1e3a83ef-b6fc-bc8c-8950-ea3f6ae5c19b.htm)。  Obtain or set the number of consecutive times to receive empty data, which is valid when the data is received, and the time consumed by each unit is [SleepTime](1e3a83ef-b6fc-bc8c-8950-ea3f6ae5c19b.htm) (继承自 [SerialBase](e877e20b-ef65-266c-05a5-5db58e6bad0c.htm)。) |
| 公共属性代码示例 | [ReceiveTimeOut](250560a7-332b-ce1b-5eb0-e3abe4ee6012.htm) | 获取或设置接收服务器反馈的时间，如果为负数，则不接收反馈   Gets or sets the time to receive server feedback, and if it is a negative number, does not receive feedback (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [RtsEnable](18c16d98-49b0-5ae2-5c2b-0c6862801c92.htm) | 获取或设置一个值，该值指示在串行通信中是否启用请求发送 (RTS) 信号。  Gets or sets a value indicating whether the request sending (RTS) signal is enabled in serial communication. (继承自 [SerialBase](e877e20b-ef65-266c-05a5-5db58e6bad0c.htm)。) |
| 公共属性 | [SendBeforeHex](72aa1de4-116f-abfe-4adb-814a0be0e560.htm) | 获取或设置在发送通信报文前追加发送的字节信息，HEX格式，通常用于lora组网时，需要携带 00 00 00 02 四个字节的站地址功能。  Obtain or set the byte information sent before sending communication packets, HEX format, usually used for LORA networking, you need to carry 00 00 00 02 four-byte station address function. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [SleepTime](86a5fbb9-3456-6beb-7113-5bbad5eabbc5.htm) | 获取或设置在正式接收对方返回数据前的时候，需要休息的时间，当设置为0的时候，不需要休息。  Get or set the time required to rest before officially receiving the data from the other party. When it is set to 0, no rest is required. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法静态成员 | [BuildReadCommand](26b20941-6ffd-7305-bdae-3ae3b482428e.htm) | 根据SAM的实际的指令，来生成实际的指令信息 |
| 公共方法静态成员 | [CheckADSCommandAndSum](9381caf4-17a8-3d72-9117-9d7ce75409ea.htm) | 检查当前的指令是否是正确的 |
| 公共方法静态成员 | [CheckADSCommandCompletion](2ecc327d-7734-26ee-365c-61175d225a17.htm) | 检查当前的接收数据信息是否一条完整的数据信息 |
| 公共方法 | [CheckSafeModuleStatus](9ac7a019-82d2-b621-758d-e6d2182e3f4c.htm) | 检测安全模块状态  Detecting Security Module Status |
| 公共方法 | [Close](992f28b3-b99b-2196-274e-2f25e9bd3c1d.htm) | 关闭当前的串口连接  Close the current serial connection (继承自 [SerialBase](e877e20b-ef65-266c-05a5-5db58e6bad0c.htm)。) |
| 受保护的方法 | [DecideWhetherQAMessage](acc40cda-e310-8537-6685-59b67eeb16ff.htm) | 决定当前的消息是否是用于问答机制返回的消息，默认直接返回 true, 实际的情况需要根据协议进行重写方法  To determine whether the current message is the message returned by the question answering mechanism, the default is true. In actual cases, the rewriting method needs to be performed according to the protocol (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [Dispose](c83e3a27-3a09-f67c-9c7f-96b7a5a4c8b0.htm) | 释放当前的对象 (继承自 [SerialBase](e877e20b-ef65-266c-05a5-5db58e6bad0c.htm)。) |
| 受保护的方法 | [Dispose(Boolean)](c509106d-35a7-2ff2-7dd8-ca67e16d7950.htm) | 释放当前的对象 (继承自 [SerialBase](e877e20b-ef65-266c-05a5-5db58e6bad0c.htm)。) |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | [ExtraAfterReadFromCoreServer](06a49a15-a3ce-5e39-28c7-2b78851eca2e.htm) | 和服务器交互完成的时候调用的方法，可以根据读写结果进行一些额外的操作，具体的操作需要根据实际的需求来重写实现  The method called when the interaction with the server is completed can perform some additional operations based on the read and write results. The specific operations need to be rewritten according to actual needs. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法静态成员 | [ExtractIdentityCard](36f5add1-fb33-d4e2-f96c-312f56f0cddd.htm) | 从数据中提取出真实的身份证信息 |
| 公共方法静态成员 | [ExtractSafeModuleNumber](5d7947c2-f718-b70c-9841-86350cfa0905.htm) | 提炼安全的模块数据信息 |
| 受保护的方法代码示例 | [ExtraOnDisconnect](de9dbba2-3b2c-1c5e-53c8-4e6f8a44e2c4.htm) | 根据实际的协议选择是否重写本方法，有些协议在断开连接之前，需要发送一些报文来关闭当前的网络通道  Select whether to rewrite this method according to the actual protocol. Some protocols need to send some packets to close the current network channel before disconnecting. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法代码示例 | [ExtraOnDisconnectAsync](210f593e-16d8-77c4-1178-58fbcfe9b27b.htm) | 根据实际的协议选择是否重写本方法，有些协议在断开连接之前，需要发送一些报文来关闭当前的网络通道  Select whether to rewrite this method according to the actual protocol. Some protocols need to send some packets to close the current network channel before disconnecting. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法静态成员 | [GetErrorDescription](f2c1ca18-618b-861b-86e5-a9ebaa6ba521.htm) | 获取错误的文本信息 |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 受保护的方法 | [GetLogTextFromBinary](67f3b2f7-957c-fa0b-8320-6799237157db.htm) | 获取当前的报文进行日志记录的时候，是否使用二进制的格式记录，默认返回 [LogMsgFormatBinary](2ed2ff87-caec-0e44-f08b-e2b015225a58.htm)，重写可以根据session对象分别返回不同记录模式  Whether to log the current packet in binary format, the default return is [LogMsgFormatBinary](2ed2ff87-caec-0e44-f08b-e2b015225a58.htm). If you want to override it, different recording modes can be returned according to session (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法静态成员 | [GetNationEnumerator](25e4d575-c7fe-6aa7-e005-b748dc4ccda6.htm) | 枚举当前的所有的民族信息，共计五十六个民族 |
| 公共方法静态成员 | [GetNationText](4010b2cb-f8d6-aaf9-b18a-5313c00c1842.htm) | 根据民族的代号来获取到民族的文本描述信息 |
| 受保护的方法 | [GetNewNetMessage](55f69034-bd64-b15b-938f-1a381e9cf19f.htm) | 获取一个新的消息对象的方法，需要在继承类里面进行重写  The method to get a new message object needs to be overridden in the inheritance class (重写 [BinaryCommunicationGetNewNetMessage](16636a1b-a6f0-9c0f-9d9c-c94e26c67459.htm).) |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法代码示例 | [InitializationOnConnect](4ab6989f-a5b5-05fe-8b44-bfa45824a2a7.htm) | 根据实际的协议选择是否重写本方法，有些协议在创建连接之后，需要进行一些初始化的信号握手，才能最终建立网络通道。  Whether to rewrite this method is based on the actual protocol. Some protocols require some initial signal handshake to establish a network channel after the connection is created. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法代码示例 | [InitializationOnConnectAsync](7da2af04-d4fc-5045-19a4-4b0d1722eef3.htm) | 根据实际的协议选择是否重写本方法，有些协议在创建连接之后，需要进行一些初始化的信号握手，才能最终建立网络通道。  Whether to rewrite this method is based on the actual protocol. Some protocols require some initial signal handshake to establish a network channel after the connection is created. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [IsOpen](24a0e05a-dd59-2e3c-e4c5-d3bb6ebc8881.htm) | 获取一个值，指示串口是否处于打开状态  Gets a value indicating whether the serial port is open (继承自 [SerialBase](e877e20b-ef65-266c-05a5-5db58e6bad0c.htm)。) |
| 受保护的方法 | [LogRevcMessage(Byte)](dab1657f-1e75-7fe1-7003-ccedfc69855d.htm) | 使用日志记录一个接收的报文信息  Logs are used to record information about a received packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [LogRevcMessage(Byte, PipeSession)](0f51ae45-575e-4228-f978-01b37d6d5df1.htm) | 使用日志记录一个接收的报文信息  Logs are used to record information about a received packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [LogSendMessage(Byte)](4f90eecc-68a2-1ef2-44ef-369f19bed3ec.htm) | 使用日志记录一个发送的报文信息  Logs are used to record information about a send packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [LogSendMessage(Byte, PipeSession)](209caadf-4f0f-40aa-04c2-ebc61f8dab13.htm) | 使用日志记录一个发送的报文信息  Logs are used to record information about a send packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法 | [Open](3501fe5b-65b5-a973-5f25-3de45d0054d6.htm) | 打开当前的管道信息，返回是否成功打开的结果对象，并通过属性 "Content" 指示当前是否为新创建的连接对象，如果是，则该值为 true  Opens the current pipe information, returns whether the result object was successfully opened, and indicates whether the current connection object is a newly created connection object through the property "Content", if so, the value is true (继承自 [SerialBase](e877e20b-ef65-266c-05a5-5db58e6bad0c.htm)。) |
| 公共方法 | [PackCommandWithHeader](d0b2cd76-c306-11fc-38c7-af8ad306cce7.htm) | 对当前的命令进行打包处理，通常是携带命令头内容，标记当前的命令的长度信息，需要进行重写，否则默认不打包  The current command is packaged, usually carrying the content of the command header, marking the length of the current command, and it needs to be rewritten, otherwise it is not packaged by default (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法静态成员 | [PackToSAMCommand](6dc8c357-7f5c-1711-419b-0cb2ab75fc43.htm) | 将指令进行打包成可以发送的数据对象 |
| 公共方法 | [ReadCard](3934f7b5-4172-ca5a-8949-4d1ff811da74.htm) | 读取卡片，如果成功的话，就返回身份证的所有的信息  Read the card, if successful, return all the information of the ID cards |
| 公共方法代码示例 | [ReadFromCoreServer(Byte)](8146dc0c-1edc-3ca7-428a-a3bdbd35e6c5.htm) | 将二进制的数据发送到管道中去，然后从管道里接收二进制的数据回来，并返回是否成功的结果对象。  Send binary data to the pipeline, and then receive binary data back from the pipeline, and return whether the success of the result object (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServer(IEnumerableByte)](044cfdd9-ec76-b89f-2f2b-5f439025006c.htm) | 将多个数据报文按顺序发到设备，并从设备接收返回的数据内容，然后拼接成一个Byte[]信息，需要重写[UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm)方法才能返回正确的结果。  Send multiple data packets to the device in sequence, and receive the returned data content from the device, and then splicing them into a Byte[] message, you need to rewrite [UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm) method to return the correct result. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServer(Byte, Boolean, Boolean)](b32ce8e8-a959-fcc2-6736-14b39fdb9a40.htm) | 将二进制的数据发送到管道中去，然后从管道里接收二进制的数据回来，并返回是否成功的结果对象。  Send binary data to the pipeline, and then receive binary data back from the pipeline, and return whether the success of the result object (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServer(CommunicationPipe, Byte, Boolean, Boolean)](ef36d56b-df09-f727-e8d3-316d37a13ee9.htm) | 使用指定的管道来进行数据通信，发送原始数据到管道，然后从管道接收相关的数据返回，本方法无锁 (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServerAsync(Byte)](a4de5a24-c67a-7b02-db05-33b85eaf6e1e.htm) | 将二进制的数据发送到管道中去，然后从管道里接收二进制的数据回来，并返回是否成功的结果对象。  Send binary data to the pipeline, and then receive binary data back from the pipeline, and return whether the success of the result object (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServerAsync(IEnumerableByte)](bcdb59b7-a70a-d5c0-1fc6-9ececcd39ea2.htm) | 将多个数据报文按顺序发到设备，并从设备接收返回的数据内容，然后拼接成一个Byte[]信息，需要重写[UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm)方法才能返回正确的结果。  Send multiple data packets to the device in sequence, and receive the returned data content from the device, and then splicing them into a Byte[] message, you need to rewrite [UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm) method to return the correct result. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServerAsync(Byte, Boolean, Boolean)](1c79c204-5ab5-9342-e65d-7455c20ea577.htm) | 将二进制的数据发送到管道中去，然后从管道里接收二进制的数据回来，并返回是否成功的结果对象。  Send binary data to the pipeline, and then receive binary data back from the pipeline, and return whether the success of the result object (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServerAsync(CommunicationPipe, Byte, Boolean, Boolean)](a254a19d-3c29-147b-335b-8adccdf4ee2f.htm) | 使用指定的管道来进行数据通信，发送原始数据到管道，然后从管道接收相关的数据返回，本方法无锁 (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadSafeModuleNumber](1ea6f232-1d72-9133-c1e6-9519576160e5.htm) | 读取身份证设备的安全模块号  Read the security module number of the ID device |
| 公共方法 | [SearchCard](4b6a12f3-7bd3-bad9-77ba-40b3a7afacae.htm) | 寻找卡片，并返回是否成功  Find cards and return success |
| 公共方法 | [SelectCard](51354f2d-81b3-490a-6b4a-2e1762216d60.htm) | 选择卡片，并返回是否成功  Select card and return success |
| 公共方法 | [SerialPortInni(ActionSerialPort)](53585409-6243-9359-1846-bf8b4e266708.htm) | 根据自定义初始化方法进行初始化串口信息  Initialize the serial port information according to the custom initialization method (继承自 [SerialBase](e877e20b-ef65-266c-05a5-5db58e6bad0c.htm)。) |
| 公共方法 | [SerialPortInni(String)](7dff5e26-6594-291e-4792-b13db9c5d613.htm) | 初始化串口信息，9600波特率，8位数据位，1位停止位，无奇偶校验  Initial serial port information, 9600 baud rate, 8 data bits, 1 stop bit, no parity (继承自 [SerialBase](e877e20b-ef65-266c-05a5-5db58e6bad0c.htm)。) |
| 公共方法 | [SerialPortInni(String, Int32)](0c7aa163-1d81-54fd-f4f5-103a7e9a341d.htm) | 初始化串口信息，波特率，8位数据位，1位停止位，无奇偶校验  Initializes serial port information, baud rate, 8-bit data bit, 1-bit stop bit, no parity (继承自 [SerialBase](e877e20b-ef65-266c-05a5-5db58e6bad0c.htm)。) |
| 公共方法 | [SerialPortInni(String, Int32, Int32, StopBits, Parity)](420c818b-48e9-fa02-59e2-f4d5897d6802.htm) | 初始化串口信息，波特率，数据位，停止位，奇偶校验需要全部自己来指定  Start serial port information, baud rate, data bit, stop bit, parity all need to be specified (继承自 [SerialBase](e877e20b-ef65-266c-05a5-5db58e6bad0c.htm)。) |
| 公共方法 | [SetDtuPipe](f765d1d7-234a-b229-9457-97e25e024dc9.htm) | 将当前的通信对象设置DTU模式，允许传入现成的管道，并返回初始化结果，如果该设备重写了握手报文，就是返回握手结果  Set the current communication object to DTU mode, allow the existing pipe to be passed in, and return the initialization result, if the device rewrites the handshake packet, the handshake result is returned (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [SetDtuPipeAsync](9d1fd12d-b804-3828-a627-a2c0eef1d23f.htm) | 将当前的通信对象设置DTU模式，允许传入现成的管道，并返回初始化结果，如果该设备重写了握手报文，就是返回握手结果  Set the current communication object to DTU mode, allow the existing pipe to be passed in, and return the initialization result, if the device rewrites the handshake packet, the handshake result is returned (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ToString](6b53a0d9-79d5-4509-63d5-dbcbe689c874.htm) | (重写 [SerialBaseToString](1ac09324-4d19-863d-d683-cbe804ef2124.htm).) |
| 公共方法 | [UnpackResponseContent](208017a7-8da2-33f8-3273-d7e7e850100f.htm) | 根据对方返回的报文命令，对命令进行基本的拆包，例如各种Modbus协议拆包为统一的核心报文，还支持对报文的验证  According to the message command returned by the other party, the command is basically unpacked, for example, various Modbus protocols are unpacked into a unified core message, and the verification of the message is also supported (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)字段

|  | 名称 | 说明 |
| --- | --- | --- |
| 受保护的字段 | [LogMsgFormatBinary](2ed2ff87-caec-0e44-f08b-e2b015225a58.htm) | 设置日志记录报文是否二进制，如果为False，那就使用ASCII码  Set whether the log message is binary, if it is False, then use ASCII code (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)示例

在使用之前需要实例化当前的对象，然后根据实际的情况填写好串口的信息，否则连接不上去。

实例化操作

[复制](# "复制")

```
private SAMSerial sAMSerial = new SAMSerial( );

public void Open( )
{
    // 下面是初始化连接的代码，在读取身份证信息之前需要被调用一次
    // The following is the code to initialize the connection. It needs to be called once before reading the ID information.
    try
    {
        sAMSerial.SerialPortInni( sp =>
        {
            sp.PortName = "COM1";
            sp.BaudRate = 115200;
            sp.DataBits = 0;
            sp.StopBits = System.IO.Ports.StopBits.None;
            sp.Parity = System.IO.Ports.Parity.None;
        } );

        sAMSerial.Open( );

    }
    catch (Exception ex)
    {
        Console.WriteLine( ex.Message );
    }
}
```

在实际的读取，我们一般放在后台进行循环扫描的操作，参见下面的代码

基本的读取操作

[复制](# "复制")

```
// 通常是放到后台进行循环扫描，此处举例开了一个线程的操作，实际上你在开发的时候需要注意GC回收垃圾的事
// Usually it is placed in the background for circular scanning. Here is an example of a thread operation. In fact, you need to pay attention to GC garbage collection when you develop
public void StartRead( )
{
    new Thread( new ThreadStart( ThreadBackgroundReadCard ) ) { IsBackground = true }.Start( );
}

private void ThreadBackgroundReadCard( )
{
    while (true)
    {
        Thread.Sleep( 100 );
        // 首先进行寻卡，成功才进行下一步
        // Find the card first, then proceed to the next step
        OperateResult search = sAMSerial.SearchCard( );
        {
            if (!search.IsSuccess)
            {
                continue;
            }
        }

        Thread.Sleep( 100 );
        // 寻卡成功后开始选卡，选卡成功才进行下一步
        // Select card after successful card search, then proceed to the next step after successful card selection
        if (sAMSerial.SelectCard( ).IsSuccess)
        {
            OperateResult<IdentityCard> read = sAMSerial.ReadCard( );
            if (read.IsSuccess)
            {
                // read.Content，详细的身份证信息，需要查看 IdentityCard 类型的定义，包含身份的名字，
                Console.WriteLine( read.Content.ToString( ) );
            }
            else
            {
                Console.WriteLine( $"读卡失败：{read.Message}" );
            }
        }
    }
}
```

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.Profinet.IDCard 命名空间](54fb02ce-3bfd-9c1d-96f8-511375556490.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## SAMSerial 构造函数 

[原文連結](http://api.hslcommunication.cn/html/60865a16-8442-3d5c-1e81-24534fa4c1da.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.IDCard](../html/54fb02ce-3bfd-9c1d-96f8-511375556490.htm "HslCommunication.Profinet.IDCard")

[SAMSerial 类](../html/48f6a433-86f0-a4f8-1751-462e12ad6923.htm "SAMSerial 类")

[SAMSerial 构造函数](../html/60865a16-8442-3d5c-1e81-24534fa4c1da.htm "SAMSerial 构造函数 ")

[SAMSerial 属性](../html/e37df3e7-0468-6303-0143-d7558a731a5d.htm "SAMSerial 属性")

[SAMSerial 方法](../html/3ee3795f-b4e0-f2ff-9050-7961204b2d49.htm "SAMSerial 方法")

[SAMSerial 字段](../html/dd810ad3-9acc-1696-cc05-6a68d6116b8b.htm "SAMSerial 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| SAMSerial 构造函数 |

实例化一个默认的对象  
Instantiate a default object

**命名空间：**
 [HslCommunication.Profinet.IDCard](54fb02ce-3bfd-9c1d-96f8-511375556490.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public SAMSerial()
```

```
Public Sub New
```

```
public:
SAMSerial()
```

```
new : unit -> SAMSerial
```

![](../icons/SectionExpanded.png)参见

#### 引用

[SAMSerial 类](48f6a433-86f0-a4f8-1751-462e12ad6923.htm)

[HslCommunication.Profinet.IDCard 命名空间](54fb02ce-3bfd-9c1d-96f8-511375556490.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## SAMSerial 属性

[原文連結](http://api.hslcommunication.cn/html/e37df3e7-0468-6303-0143-d7558a731a5d.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.IDCard](../html/54fb02ce-3bfd-9c1d-96f8-511375556490.htm "HslCommunication.Profinet.IDCard")

[SAMSerial 类](../html/48f6a433-86f0-a4f8-1751-462e12ad6923.htm "SAMSerial 类")

[SAMSerial 构造函数](../html/60865a16-8442-3d5c-1e81-24534fa4c1da.htm "SAMSerial 构造函数 ")

[SAMSerial 属性](../html/e37df3e7-0468-6303-0143-d7558a731a5d.htm "SAMSerial 属性")

[SAMSerial 方法](../html/3ee3795f-b4e0-f2ff-9050-7961204b2d49.htm "SAMSerial 方法")

[SAMSerial 字段](../html/dd810ad3-9acc-1696-cc05-6a68d6116b8b.htm "SAMSerial 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| SAMSerial 属性 |

[SAMSerial](48f6a433-86f0-a4f8-1751-462e12ad6923.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [BaudRate](1c33c44a-d2b9-badd-df76-09e155d26cc1.htm) | 当前连接串口信息的波特率  Baud rate of current connection serial port information (继承自 [SerialBase](e877e20b-ef65-266c-05a5-5db58e6bad0c.htm)。) |
| 公共属性 | [CommunicationPipe](cb528f9b-fe2c-5ac7-6912-2c3c5c1fdf2a.htm) | 获取或设置当前的管道信息，管道类型为[CommunicationPipe](ad8e903b-d289-21b0-077e-15ede8005221.htm)的继承类，内置了[PipeTcpNet](26b94bd5-2695-36e7-6378-c3d41846be8f.htm)管道，[PipeUdpNet](5a0b810e-187e-826b-0a13-34fd5576aa0e.htm)管道，[PipeSerialPort](9d924607-3dfc-3916-b532-c5dec48fe347.htm)管道等  Get or set the current pipeline information, the pipeline type is [CommunicationPipe](ad8e903b-d289-21b0-077e-15ede8005221.htm) inheritance class, [PipeTcpNet](26b94bd5-2695-36e7-6378-c3d41846be8f.htm) pipeline, [PipeUdpNet](5a0b810e-187e-826b-0a13-34fd5576aa0e.htm) pipeline, [PipeSerialPort](9d924607-3dfc-3916-b532-c5dec48fe347.htm) pipeline, etc (继承自 [SerialBase](e877e20b-ef65-266c-05a5-5db58e6bad0c.htm)。) |
| 公共属性 | [ConnectionId](ce1efb4b-fbc9-4683-163e-0a8e71e99d6b.htm) | 当前连接的唯一ID号，默认为长度20的guid码加随机数组成，方便列表管理，也可以自己指定  The unique ID number of the current connection. The default is a 20-digit guid code plus a random number. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [IsClearCacheBeforeRead](cf2aaeff-e4eb-849a-4783-d10e5e965515.htm) | 是否在发送数据前清空缓冲数据，默认是false  Whether to empty the buffer before sending data, the default is false (继承自 [SerialBase](e877e20b-ef65-266c-05a5-5db58e6bad0c.htm)。) |
| 公共属性代码示例 | [LogNet](8c56f246-88ac-afb8-8c06-972fafe4f78a.htm) | 组件的日志工具，支持日志记录，只要实例化后，当前网络的基本信息，就以[DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm)等级进行输出  The component's logging tool supports logging. As long as the instantiation of the basic network information, the output will be output at [DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm) (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [PortName](250b8746-0bb7-4caf-2a43-11109f2bda20.htm) | 当前连接串口信息的端口号名称  The port name of the current connection serial port information (继承自 [SerialBase](e877e20b-ef65-266c-05a5-5db58e6bad0c.htm)。) |
| 公共属性 | [ReceiveEmptyDataCount](480b2206-1de4-6a74-6f8e-7c30a2ca05f9.htm) | 获取或设置连续接收空的数据次数，在数据接收完成时有效，每个单位消耗的时间为[SleepTime](1e3a83ef-b6fc-bc8c-8950-ea3f6ae5c19b.htm)。  Obtain or set the number of consecutive times to receive empty data, which is valid when the data is received, and the time consumed by each unit is [SleepTime](1e3a83ef-b6fc-bc8c-8950-ea3f6ae5c19b.htm) (继承自 [SerialBase](e877e20b-ef65-266c-05a5-5db58e6bad0c.htm)。) |
| 公共属性代码示例 | [ReceiveTimeOut](250560a7-332b-ce1b-5eb0-e3abe4ee6012.htm) | 获取或设置接收服务器反馈的时间，如果为负数，则不接收反馈   Gets or sets the time to receive server feedback, and if it is a negative number, does not receive feedback (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [RtsEnable](18c16d98-49b0-5ae2-5c2b-0c6862801c92.htm) | 获取或设置一个值，该值指示在串行通信中是否启用请求发送 (RTS) 信号。  Gets or sets a value indicating whether the request sending (RTS) signal is enabled in serial communication. (继承自 [SerialBase](e877e20b-ef65-266c-05a5-5db58e6bad0c.htm)。) |
| 公共属性 | [SendBeforeHex](72aa1de4-116f-abfe-4adb-814a0be0e560.htm) | 获取或设置在发送通信报文前追加发送的字节信息，HEX格式，通常用于lora组网时，需要携带 00 00 00 02 四个字节的站地址功能。  Obtain or set the byte information sent before sending communication packets, HEX format, usually used for LORA networking, you need to carry 00 00 00 02 four-byte station address function. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [SleepTime](86a5fbb9-3456-6beb-7113-5bbad5eabbc5.htm) | 获取或设置在正式接收对方返回数据前的时候，需要休息的时间，当设置为0的时候，不需要休息。  Get or set the time required to rest before officially receiving the data from the other party. When it is set to 0, no rest is required. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[SAMSerial 类](48f6a433-86f0-a4f8-1751-462e12ad6923.htm)

[HslCommunication.Profinet.IDCard 命名空间](54fb02ce-3bfd-9c1d-96f8-511375556490.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## SAMSerial 方法

[原文連結](http://api.hslcommunication.cn/html/3ee3795f-b4e0-f2ff-9050-7961204b2d49.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.IDCard](../html/54fb02ce-3bfd-9c1d-96f8-511375556490.htm "HslCommunication.Profinet.IDCard")

[SAMSerial 类](../html/48f6a433-86f0-a4f8-1751-462e12ad6923.htm "SAMSerial 类")

[SAMSerial 方法](../html/3ee3795f-b4e0-f2ff-9050-7961204b2d49.htm "SAMSerial 方法")

[BuildReadCommand 方法](../html/26b20941-6ffd-7305-bdae-3ae3b482428e.htm "BuildReadCommand 方法 ")

[CheckADSCommandAndSum 方法](../html/9381caf4-17a8-3d72-9117-9d7ce75409ea.htm "CheckADSCommandAndSum 方法 ")

[CheckADSCommandCompletion 方法](../html/2ecc327d-7734-26ee-365c-61175d225a17.htm "CheckADSCommandCompletion 方法 ")

[CheckSafeModuleStatus 方法](../html/9ac7a019-82d2-b621-758d-e6d2182e3f4c.htm "CheckSafeModuleStatus 方法 ")

[ExtractIdentityCard 方法](../html/36f5add1-fb33-d4e2-f96c-312f56f0cddd.htm "ExtractIdentityCard 方法 ")

[ExtractSafeModuleNumber 方法](../html/5d7947c2-f718-b70c-9841-86350cfa0905.htm "ExtractSafeModuleNumber 方法 ")

[GetErrorDescription 方法](../html/f2c1ca18-618b-861b-86e5-a9ebaa6ba521.htm "GetErrorDescription 方法 ")

[GetNationEnumerator 方法](../html/25e4d575-c7fe-6aa7-e005-b748dc4ccda6.htm "GetNationEnumerator 方法 ")

[GetNationText 方法](../html/4010b2cb-f8d6-aaf9-b18a-5313c00c1842.htm "GetNationText 方法 ")

[GetNewNetMessage 方法](../html/55f69034-bd64-b15b-938f-1a381e9cf19f.htm "GetNewNetMessage 方法 ")

[PackToSAMCommand 方法](../html/6dc8c357-7f5c-1711-419b-0cb2ab75fc43.htm "PackToSAMCommand 方法 ")

[ReadCard 方法](../html/3934f7b5-4172-ca5a-8949-4d1ff811da74.htm "ReadCard 方法 ")

[ReadSafeModuleNumber 方法](../html/1ea6f232-1d72-9133-c1e6-9519576160e5.htm "ReadSafeModuleNumber 方法 ")

[SearchCard 方法](../html/4b6a12f3-7bd3-bad9-77ba-40b3a7afacae.htm "SearchCard 方法 ")

[SelectCard 方法](../html/51354f2d-81b3-490a-6b4a-2e1762216d60.htm "SelectCard 方法 ")

[ToString 方法](../html/6b53a0d9-79d5-4509-63d5-dbcbe689c874.htm "ToString 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| SAMSerial 方法 |

[SAMSerial](48f6a433-86f0-a4f8-1751-462e12ad6923.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法静态成员 | [BuildReadCommand](26b20941-6ffd-7305-bdae-3ae3b482428e.htm) | 根据SAM的实际的指令，来生成实际的指令信息 |
| 公共方法静态成员 | [CheckADSCommandAndSum](9381caf4-17a8-3d72-9117-9d7ce75409ea.htm) | 检查当前的指令是否是正确的 |
| 公共方法静态成员 | [CheckADSCommandCompletion](2ecc327d-7734-26ee-365c-61175d225a17.htm) | 检查当前的接收数据信息是否一条完整的数据信息 |
| 公共方法 | [CheckSafeModuleStatus](9ac7a019-82d2-b621-758d-e6d2182e3f4c.htm) | 检测安全模块状态  Detecting Security Module Status |
| 公共方法 | [Close](992f28b3-b99b-2196-274e-2f25e9bd3c1d.htm) | 关闭当前的串口连接  Close the current serial connection (继承自 [SerialBase](e877e20b-ef65-266c-05a5-5db58e6bad0c.htm)。) |
| 受保护的方法 | [DecideWhetherQAMessage](acc40cda-e310-8537-6685-59b67eeb16ff.htm) | 决定当前的消息是否是用于问答机制返回的消息，默认直接返回 true, 实际的情况需要根据协议进行重写方法  To determine whether the current message is the message returned by the question answering mechanism, the default is true. In actual cases, the rewriting method needs to be performed according to the protocol (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [Dispose](c83e3a27-3a09-f67c-9c7f-96b7a5a4c8b0.htm) | 释放当前的对象 (继承自 [SerialBase](e877e20b-ef65-266c-05a5-5db58e6bad0c.htm)。) |
| 受保护的方法 | [Dispose(Boolean)](c509106d-35a7-2ff2-7dd8-ca67e16d7950.htm) | 释放当前的对象 (继承自 [SerialBase](e877e20b-ef65-266c-05a5-5db58e6bad0c.htm)。) |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | [ExtraAfterReadFromCoreServer](06a49a15-a3ce-5e39-28c7-2b78851eca2e.htm) | 和服务器交互完成的时候调用的方法，可以根据读写结果进行一些额外的操作，具体的操作需要根据实际的需求来重写实现  The method called when the interaction with the server is completed can perform some additional operations based on the read and write results. The specific operations need to be rewritten according to actual needs. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法静态成员 | [ExtractIdentityCard](36f5add1-fb33-d4e2-f96c-312f56f0cddd.htm) | 从数据中提取出真实的身份证信息 |
| 公共方法静态成员 | [ExtractSafeModuleNumber](5d7947c2-f718-b70c-9841-86350cfa0905.htm) | 提炼安全的模块数据信息 |
| 受保护的方法代码示例 | [ExtraOnDisconnect](de9dbba2-3b2c-1c5e-53c8-4e6f8a44e2c4.htm) | 根据实际的协议选择是否重写本方法，有些协议在断开连接之前，需要发送一些报文来关闭当前的网络通道  Select whether to rewrite this method according to the actual protocol. Some protocols need to send some packets to close the current network channel before disconnecting. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法代码示例 | [ExtraOnDisconnectAsync](210f593e-16d8-77c4-1178-58fbcfe9b27b.htm) | 根据实际的协议选择是否重写本方法，有些协议在断开连接之前，需要发送一些报文来关闭当前的网络通道  Select whether to rewrite this method according to the actual protocol. Some protocols need to send some packets to close the current network channel before disconnecting. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法静态成员 | [GetErrorDescription](f2c1ca18-618b-861b-86e5-a9ebaa6ba521.htm) | 获取错误的文本信息 |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 受保护的方法 | [GetLogTextFromBinary](67f3b2f7-957c-fa0b-8320-6799237157db.htm) | 获取当前的报文进行日志记录的时候，是否使用二进制的格式记录，默认返回 [LogMsgFormatBinary](2ed2ff87-caec-0e44-f08b-e2b015225a58.htm)，重写可以根据session对象分别返回不同记录模式  Whether to log the current packet in binary format, the default return is [LogMsgFormatBinary](2ed2ff87-caec-0e44-f08b-e2b015225a58.htm). If you want to override it, different recording modes can be returned according to session (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法静态成员 | [GetNationEnumerator](25e4d575-c7fe-6aa7-e005-b748dc4ccda6.htm) | 枚举当前的所有的民族信息，共计五十六个民族 |
| 公共方法静态成员 | [GetNationText](4010b2cb-f8d6-aaf9-b18a-5313c00c1842.htm) | 根据民族的代号来获取到民族的文本描述信息 |
| 受保护的方法 | [GetNewNetMessage](55f69034-bd64-b15b-938f-1a381e9cf19f.htm) | 获取一个新的消息对象的方法，需要在继承类里面进行重写  The method to get a new message object needs to be overridden in the inheritance class (重写 [BinaryCommunicationGetNewNetMessage](16636a1b-a6f0-9c0f-9d9c-c94e26c67459.htm).) |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法代码示例 | [InitializationOnConnect](4ab6989f-a5b5-05fe-8b44-bfa45824a2a7.htm) | 根据实际的协议选择是否重写本方法，有些协议在创建连接之后，需要进行一些初始化的信号握手，才能最终建立网络通道。  Whether to rewrite this method is based on the actual protocol. Some protocols require some initial signal handshake to establish a network channel after the connection is created. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法代码示例 | [InitializationOnConnectAsync](7da2af04-d4fc-5045-19a4-4b0d1722eef3.htm) | 根据实际的协议选择是否重写本方法，有些协议在创建连接之后，需要进行一些初始化的信号握手，才能最终建立网络通道。  Whether to rewrite this method is based on the actual protocol. Some protocols require some initial signal handshake to establish a network channel after the connection is created. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [IsOpen](24a0e05a-dd59-2e3c-e4c5-d3bb6ebc8881.htm) | 获取一个值，指示串口是否处于打开状态  Gets a value indicating whether the serial port is open (继承自 [SerialBase](e877e20b-ef65-266c-05a5-5db58e6bad0c.htm)。) |
| 受保护的方法 | [LogRevcMessage(Byte)](dab1657f-1e75-7fe1-7003-ccedfc69855d.htm) | 使用日志记录一个接收的报文信息  Logs are used to record information about a received packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [LogRevcMessage(Byte, PipeSession)](0f51ae45-575e-4228-f978-01b37d6d5df1.htm) | 使用日志记录一个接收的报文信息  Logs are used to record information about a received packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [LogSendMessage(Byte)](4f90eecc-68a2-1ef2-44ef-369f19bed3ec.htm) | 使用日志记录一个发送的报文信息  Logs are used to record information about a send packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [LogSendMessage(Byte, PipeSession)](209caadf-4f0f-40aa-04c2-ebc61f8dab13.htm) | 使用日志记录一个发送的报文信息  Logs are used to record information about a send packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法 | [Open](3501fe5b-65b5-a973-5f25-3de45d0054d6.htm) | 打开当前的管道信息，返回是否成功打开的结果对象，并通过属性 "Content" 指示当前是否为新创建的连接对象，如果是，则该值为 true  Opens the current pipe information, returns whether the result object was successfully opened, and indicates whether the current connection object is a newly created connection object through the property "Content", if so, the value is true (继承自 [SerialBase](e877e20b-ef65-266c-05a5-5db58e6bad0c.htm)。) |
| 公共方法 | [PackCommandWithHeader](d0b2cd76-c306-11fc-38c7-af8ad306cce7.htm) | 对当前的命令进行打包处理，通常是携带命令头内容，标记当前的命令的长度信息，需要进行重写，否则默认不打包  The current command is packaged, usually carrying the content of the command header, marking the length of the current command, and it needs to be rewritten, otherwise it is not packaged by default (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法静态成员 | [PackToSAMCommand](6dc8c357-7f5c-1711-419b-0cb2ab75fc43.htm) | 将指令进行打包成可以发送的数据对象 |
| 公共方法 | [ReadCard](3934f7b5-4172-ca5a-8949-4d1ff811da74.htm) | 读取卡片，如果成功的话，就返回身份证的所有的信息  Read the card, if successful, return all the information of the ID cards |
| 公共方法代码示例 | [ReadFromCoreServer(Byte)](8146dc0c-1edc-3ca7-428a-a3bdbd35e6c5.htm) | 将二进制的数据发送到管道中去，然后从管道里接收二进制的数据回来，并返回是否成功的结果对象。  Send binary data to the pipeline, and then receive binary data back from the pipeline, and return whether the success of the result object (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServer(IEnumerableByte)](044cfdd9-ec76-b89f-2f2b-5f439025006c.htm) | 将多个数据报文按顺序发到设备，并从设备接收返回的数据内容，然后拼接成一个Byte[]信息，需要重写[UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm)方法才能返回正确的结果。  Send multiple data packets to the device in sequence, and receive the returned data content from the device, and then splicing them into a Byte[] message, you need to rewrite [UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm) method to return the correct result. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServer(Byte, Boolean, Boolean)](b32ce8e8-a959-fcc2-6736-14b39fdb9a40.htm) | 将二进制的数据发送到管道中去，然后从管道里接收二进制的数据回来，并返回是否成功的结果对象。  Send binary data to the pipeline, and then receive binary data back from the pipeline, and return whether the success of the result object (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServer(CommunicationPipe, Byte, Boolean, Boolean)](ef36d56b-df09-f727-e8d3-316d37a13ee9.htm) | 使用指定的管道来进行数据通信，发送原始数据到管道，然后从管道接收相关的数据返回，本方法无锁 (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServerAsync(Byte)](a4de5a24-c67a-7b02-db05-33b85eaf6e1e.htm) | 将二进制的数据发送到管道中去，然后从管道里接收二进制的数据回来，并返回是否成功的结果对象。  Send binary data to the pipeline, and then receive binary data back from the pipeline, and return whether the success of the result object (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServerAsync(IEnumerableByte)](bcdb59b7-a70a-d5c0-1fc6-9ececcd39ea2.htm) | 将多个数据报文按顺序发到设备，并从设备接收返回的数据内容，然后拼接成一个Byte[]信息，需要重写[UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm)方法才能返回正确的结果。  Send multiple data packets to the device in sequence, and receive the returned data content from the device, and then splicing them into a Byte[] message, you need to rewrite [UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm) method to return the correct result. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServerAsync(Byte, Boolean, Boolean)](1c79c204-5ab5-9342-e65d-7455c20ea577.htm) | 将二进制的数据发送到管道中去，然后从管道里接收二进制的数据回来，并返回是否成功的结果对象。  Send binary data to the pipeline, and then receive binary data back from the pipeline, and return whether the success of the result object (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServerAsync(CommunicationPipe, Byte, Boolean, Boolean)](a254a19d-3c29-147b-335b-8adccdf4ee2f.htm) | 使用指定的管道来进行数据通信，发送原始数据到管道，然后从管道接收相关的数据返回，本方法无锁 (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadSafeModuleNumber](1ea6f232-1d72-9133-c1e6-9519576160e5.htm) | 读取身份证设备的安全模块号  Read the security module number of the ID device |
| 公共方法 | [SearchCard](4b6a12f3-7bd3-bad9-77ba-40b3a7afacae.htm) | 寻找卡片，并返回是否成功  Find cards and return success |
| 公共方法 | [SelectCard](51354f2d-81b3-490a-6b4a-2e1762216d60.htm) | 选择卡片，并返回是否成功  Select card and return success |
| 公共方法 | [SerialPortInni(ActionSerialPort)](53585409-6243-9359-1846-bf8b4e266708.htm) | 根据自定义初始化方法进行初始化串口信息  Initialize the serial port information according to the custom initialization method (继承自 [SerialBase](e877e20b-ef65-266c-05a5-5db58e6bad0c.htm)。) |
| 公共方法 | [SerialPortInni(String)](7dff5e26-6594-291e-4792-b13db9c5d613.htm) | 初始化串口信息，9600波特率，8位数据位，1位停止位，无奇偶校验  Initial serial port information, 9600 baud rate, 8 data bits, 1 stop bit, no parity (继承自 [SerialBase](e877e20b-ef65-266c-05a5-5db58e6bad0c.htm)。) |
| 公共方法 | [SerialPortInni(String, Int32)](0c7aa163-1d81-54fd-f4f5-103a7e9a341d.htm) | 初始化串口信息，波特率，8位数据位，1位停止位，无奇偶校验  Initializes serial port information, baud rate, 8-bit data bit, 1-bit stop bit, no parity (继承自 [SerialBase](e877e20b-ef65-266c-05a5-5db58e6bad0c.htm)。) |
| 公共方法 | [SerialPortInni(String, Int32, Int32, StopBits, Parity)](420c818b-48e9-fa02-59e2-f4d5897d6802.htm) | 初始化串口信息，波特率，数据位，停止位，奇偶校验需要全部自己来指定  Start serial port information, baud rate, data bit, stop bit, parity all need to be specified (继承自 [SerialBase](e877e20b-ef65-266c-05a5-5db58e6bad0c.htm)。) |
| 公共方法 | [SetDtuPipe](f765d1d7-234a-b229-9457-97e25e024dc9.htm) | 将当前的通信对象设置DTU模式，允许传入现成的管道，并返回初始化结果，如果该设备重写了握手报文，就是返回握手结果  Set the current communication object to DTU mode, allow the existing pipe to be passed in, and return the initialization result, if the device rewrites the handshake packet, the handshake result is returned (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [SetDtuPipeAsync](9d1fd12d-b804-3828-a627-a2c0eef1d23f.htm) | 将当前的通信对象设置DTU模式，允许传入现成的管道，并返回初始化结果，如果该设备重写了握手报文，就是返回握手结果  Set the current communication object to DTU mode, allow the existing pipe to be passed in, and return the initialization result, if the device rewrites the handshake packet, the handshake result is returned (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ToString](6b53a0d9-79d5-4509-63d5-dbcbe689c874.htm) | (重写 [SerialBaseToString](1ac09324-4d19-863d-d683-cbe804ef2124.htm).) |
| 公共方法 | [UnpackResponseContent](208017a7-8da2-33f8-3273-d7e7e850100f.htm) | 根据对方返回的报文命令，对命令进行基本的拆包，例如各种Modbus协议拆包为统一的核心报文，还支持对报文的验证  According to the message command returned by the other party, the command is basically unpacked, for example, various Modbus protocols are unpacked into a unified core message, and the verification of the message is also supported (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[SAMSerial 类](48f6a433-86f0-a4f8-1751-462e12ad6923.htm)

[HslCommunication.Profinet.IDCard 命名空间](54fb02ce-3bfd-9c1d-96f8-511375556490.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## BuildReadCommand 方法 

[原文連結](http://api.hslcommunication.cn/html/26b20941-6ffd-7305-bdae-3ae3b482428e.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.IDCard](../html/54fb02ce-3bfd-9c1d-96f8-511375556490.htm "HslCommunication.Profinet.IDCard")

[SAMSerial 类](../html/48f6a433-86f0-a4f8-1751-462e12ad6923.htm "SAMSerial 类")

[SAMSerial 方法](../html/3ee3795f-b4e0-f2ff-9050-7961204b2d49.htm "SAMSerial 方法")

[BuildReadCommand 方法](../html/26b20941-6ffd-7305-bdae-3ae3b482428e.htm "BuildReadCommand 方法 ")

[CheckADSCommandAndSum 方法](../html/9381caf4-17a8-3d72-9117-9d7ce75409ea.htm "CheckADSCommandAndSum 方法 ")

[CheckADSCommandCompletion 方法](../html/2ecc327d-7734-26ee-365c-61175d225a17.htm "CheckADSCommandCompletion 方法 ")

[CheckSafeModuleStatus 方法](../html/9ac7a019-82d2-b621-758d-e6d2182e3f4c.htm "CheckSafeModuleStatus 方法 ")

[ExtractIdentityCard 方法](../html/36f5add1-fb33-d4e2-f96c-312f56f0cddd.htm "ExtractIdentityCard 方法 ")

[ExtractSafeModuleNumber 方法](../html/5d7947c2-f718-b70c-9841-86350cfa0905.htm "ExtractSafeModuleNumber 方法 ")

[GetErrorDescription 方法](../html/f2c1ca18-618b-861b-86e5-a9ebaa6ba521.htm "GetErrorDescription 方法 ")

[GetNationEnumerator 方法](../html/25e4d575-c7fe-6aa7-e005-b748dc4ccda6.htm "GetNationEnumerator 方法 ")

[GetNationText 方法](../html/4010b2cb-f8d6-aaf9-b18a-5313c00c1842.htm "GetNationText 方法 ")

[GetNewNetMessage 方法](../html/55f69034-bd64-b15b-938f-1a381e9cf19f.htm "GetNewNetMessage 方法 ")

[PackToSAMCommand 方法](../html/6dc8c357-7f5c-1711-419b-0cb2ab75fc43.htm "PackToSAMCommand 方法 ")

[ReadCard 方法](../html/3934f7b5-4172-ca5a-8949-4d1ff811da74.htm "ReadCard 方法 ")

[ReadSafeModuleNumber 方法](../html/1ea6f232-1d72-9133-c1e6-9519576160e5.htm "ReadSafeModuleNumber 方法 ")

[SearchCard 方法](../html/4b6a12f3-7bd3-bad9-77ba-40b3a7afacae.htm "SearchCard 方法 ")

[SelectCard 方法](../html/51354f2d-81b3-490a-6b4a-2e1762216d60.htm "SelectCard 方法 ")

[ToString 方法](../html/6b53a0d9-79d5-4509-63d5-dbcbe689c874.htm "ToString 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| SAMSerialBuildReadCommand 方法 |

根据SAM的实际的指令，来生成实际的指令信息

**命名空间：**
 [HslCommunication.Profinet.IDCard](54fb02ce-3bfd-9c1d-96f8-511375556490.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static byte[] BuildReadCommand(
	byte cmd,
	byte para,
	byte[] data
)
```

```
Public Shared Function BuildReadCommand ( 
	cmd As Byte,
	para As Byte,
	data As Byte()
) As Byte()
```

```
public:
static array<unsigned char>^ BuildReadCommand(
	unsigned char cmd, 
	unsigned char para, 
	array<unsigned char>^ data
)
```

```
static member BuildReadCommand : 
        cmd : byte * 
        para : byte * 
        data : byte[] -> byte[] 
```

#### 参数

cmd
:   类型：SystemByte  
    命令码

para
:   类型：SystemByte  
    参数信息

data
:   类型：SystemByte  
    数据内容

#### 返回值

类型：Byte  
字符串的结果信息

![](../icons/SectionExpanded.png)参见

#### 引用

[SAMSerial 类](48f6a433-86f0-a4f8-1751-462e12ad6923.htm)

[HslCommunication.Profinet.IDCard 命名空间](54fb02ce-3bfd-9c1d-96f8-511375556490.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## CheckADSCommandAndSum 方法 

[原文連結](http://api.hslcommunication.cn/html/9381caf4-17a8-3d72-9117-9d7ce75409ea.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.IDCard](../html/54fb02ce-3bfd-9c1d-96f8-511375556490.htm "HslCommunication.Profinet.IDCard")

[SAMSerial 类](../html/48f6a433-86f0-a4f8-1751-462e12ad6923.htm "SAMSerial 类")

[SAMSerial 方法](../html/3ee3795f-b4e0-f2ff-9050-7961204b2d49.htm "SAMSerial 方法")

[BuildReadCommand 方法](../html/26b20941-6ffd-7305-bdae-3ae3b482428e.htm "BuildReadCommand 方法 ")

[CheckADSCommandAndSum 方法](../html/9381caf4-17a8-3d72-9117-9d7ce75409ea.htm "CheckADSCommandAndSum 方法 ")

[CheckADSCommandCompletion 方法](../html/2ecc327d-7734-26ee-365c-61175d225a17.htm "CheckADSCommandCompletion 方法 ")

[CheckSafeModuleStatus 方法](../html/9ac7a019-82d2-b621-758d-e6d2182e3f4c.htm "CheckSafeModuleStatus 方法 ")

[ExtractIdentityCard 方法](../html/36f5add1-fb33-d4e2-f96c-312f56f0cddd.htm "ExtractIdentityCard 方法 ")

[ExtractSafeModuleNumber 方法](../html/5d7947c2-f718-b70c-9841-86350cfa0905.htm "ExtractSafeModuleNumber 方法 ")

[GetErrorDescription 方法](../html/f2c1ca18-618b-861b-86e5-a9ebaa6ba521.htm "GetErrorDescription 方法 ")

[GetNationEnumerator 方法](../html/25e4d575-c7fe-6aa7-e005-b748dc4ccda6.htm "GetNationEnumerator 方法 ")

[GetNationText 方法](../html/4010b2cb-f8d6-aaf9-b18a-5313c00c1842.htm "GetNationText 方法 ")

[GetNewNetMessage 方法](../html/55f69034-bd64-b15b-938f-1a381e9cf19f.htm "GetNewNetMessage 方法 ")

[PackToSAMCommand 方法](../html/6dc8c357-7f5c-1711-419b-0cb2ab75fc43.htm "PackToSAMCommand 方法 ")

[ReadCard 方法](../html/3934f7b5-4172-ca5a-8949-4d1ff811da74.htm "ReadCard 方法 ")

[ReadSafeModuleNumber 方法](../html/1ea6f232-1d72-9133-c1e6-9519576160e5.htm "ReadSafeModuleNumber 方法 ")

[SearchCard 方法](../html/4b6a12f3-7bd3-bad9-77ba-40b3a7afacae.htm "SearchCard 方法 ")

[SelectCard 方法](../html/51354f2d-81b3-490a-6b4a-2e1762216d60.htm "SelectCard 方法 ")

[ToString 方法](../html/6b53a0d9-79d5-4509-63d5-dbcbe689c874.htm "ToString 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| SAMSerialCheckADSCommandAndSum 方法 |

检查当前的指令是否是正确的

**命名空间：**
 [HslCommunication.Profinet.IDCard](54fb02ce-3bfd-9c1d-96f8-511375556490.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult CheckADSCommandAndSum(
	byte[] input
)
```

```
Public Shared Function CheckADSCommandAndSum ( 
	input As Byte()
) As OperateResult
```

```
public:
static OperateResult^ CheckADSCommandAndSum(
	array<unsigned char>^ input
)
```

```
static member CheckADSCommandAndSum : 
        input : byte[] -> OperateResult 
```

#### 参数

input
:   类型：SystemByte  
    输入的指令信息

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否校验成功

![](../icons/SectionExpanded.png)参见

#### 引用

[SAMSerial 类](48f6a433-86f0-a4f8-1751-462e12ad6923.htm)

[HslCommunication.Profinet.IDCard 命名空间](54fb02ce-3bfd-9c1d-96f8-511375556490.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## CheckADSCommandCompletion 方法 

[原文連結](http://api.hslcommunication.cn/html/2ecc327d-7734-26ee-365c-61175d225a17.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.IDCard](../html/54fb02ce-3bfd-9c1d-96f8-511375556490.htm "HslCommunication.Profinet.IDCard")

[SAMSerial 类](../html/48f6a433-86f0-a4f8-1751-462e12ad6923.htm "SAMSerial 类")

[SAMSerial 方法](../html/3ee3795f-b4e0-f2ff-9050-7961204b2d49.htm "SAMSerial 方法")

[BuildReadCommand 方法](../html/26b20941-6ffd-7305-bdae-3ae3b482428e.htm "BuildReadCommand 方法 ")

[CheckADSCommandAndSum 方法](../html/9381caf4-17a8-3d72-9117-9d7ce75409ea.htm "CheckADSCommandAndSum 方法 ")

[CheckADSCommandCompletion 方法](../html/2ecc327d-7734-26ee-365c-61175d225a17.htm "CheckADSCommandCompletion 方法 ")

[CheckSafeModuleStatus 方法](../html/9ac7a019-82d2-b621-758d-e6d2182e3f4c.htm "CheckSafeModuleStatus 方法 ")

[ExtractIdentityCard 方法](../html/36f5add1-fb33-d4e2-f96c-312f56f0cddd.htm "ExtractIdentityCard 方法 ")

[ExtractSafeModuleNumber 方法](../html/5d7947c2-f718-b70c-9841-86350cfa0905.htm "ExtractSafeModuleNumber 方法 ")

[GetErrorDescription 方法](../html/f2c1ca18-618b-861b-86e5-a9ebaa6ba521.htm "GetErrorDescription 方法 ")

[GetNationEnumerator 方法](../html/25e4d575-c7fe-6aa7-e005-b748dc4ccda6.htm "GetNationEnumerator 方法 ")

[GetNationText 方法](../html/4010b2cb-f8d6-aaf9-b18a-5313c00c1842.htm "GetNationText 方法 ")

[GetNewNetMessage 方法](../html/55f69034-bd64-b15b-938f-1a381e9cf19f.htm "GetNewNetMessage 方法 ")

[PackToSAMCommand 方法](../html/6dc8c357-7f5c-1711-419b-0cb2ab75fc43.htm "PackToSAMCommand 方法 ")

[ReadCard 方法](../html/3934f7b5-4172-ca5a-8949-4d1ff811da74.htm "ReadCard 方法 ")

[ReadSafeModuleNumber 方法](../html/1ea6f232-1d72-9133-c1e6-9519576160e5.htm "ReadSafeModuleNumber 方法 ")

[SearchCard 方法](../html/4b6a12f3-7bd3-bad9-77ba-40b3a7afacae.htm "SearchCard 方法 ")

[SelectCard 方法](../html/51354f2d-81b3-490a-6b4a-2e1762216d60.htm "SelectCard 方法 ")

[ToString 方法](../html/6b53a0d9-79d5-4509-63d5-dbcbe689c874.htm "ToString 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| SAMSerialCheckADSCommandCompletion 方法 |

检查当前的接收数据信息是否一条完整的数据信息

**命名空间：**
 [HslCommunication.Profinet.IDCard](54fb02ce-3bfd-9c1d-96f8-511375556490.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static bool CheckADSCommandCompletion(
	List<byte> input
)
```

```
Public Shared Function CheckADSCommandCompletion ( 
	input As List(Of Byte)
) As Boolean
```

```
public:
static bool CheckADSCommandCompletion(
	List<unsigned char>^ input
)
```

```
static member CheckADSCommandCompletion : 
        input : List<byte> -> bool 
```

#### 参数

input
:   类型：System.Collections.GenericListByte  
    输入的信息

#### 返回值

类型：Boolean  
是否接收完成

![](../icons/SectionExpanded.png)参见

#### 引用

[SAMSerial 类](48f6a433-86f0-a4f8-1751-462e12ad6923.htm)

[HslCommunication.Profinet.IDCard 命名空间](54fb02ce-3bfd-9c1d-96f8-511375556490.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## CheckSafeModuleStatus 方法 

[原文連結](http://api.hslcommunication.cn/html/9ac7a019-82d2-b621-758d-e6d2182e3f4c.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.IDCard](../html/54fb02ce-3bfd-9c1d-96f8-511375556490.htm "HslCommunication.Profinet.IDCard")

[SAMSerial 类](../html/48f6a433-86f0-a4f8-1751-462e12ad6923.htm "SAMSerial 类")

[SAMSerial 方法](../html/3ee3795f-b4e0-f2ff-9050-7961204b2d49.htm "SAMSerial 方法")

[BuildReadCommand 方法](../html/26b20941-6ffd-7305-bdae-3ae3b482428e.htm "BuildReadCommand 方法 ")

[CheckADSCommandAndSum 方法](../html/9381caf4-17a8-3d72-9117-9d7ce75409ea.htm "CheckADSCommandAndSum 方法 ")

[CheckADSCommandCompletion 方法](../html/2ecc327d-7734-26ee-365c-61175d225a17.htm "CheckADSCommandCompletion 方法 ")

[CheckSafeModuleStatus 方法](../html/9ac7a019-82d2-b621-758d-e6d2182e3f4c.htm "CheckSafeModuleStatus 方法 ")

[ExtractIdentityCard 方法](../html/36f5add1-fb33-d4e2-f96c-312f56f0cddd.htm "ExtractIdentityCard 方法 ")

[ExtractSafeModuleNumber 方法](../html/5d7947c2-f718-b70c-9841-86350cfa0905.htm "ExtractSafeModuleNumber 方法 ")

[GetErrorDescription 方法](../html/f2c1ca18-618b-861b-86e5-a9ebaa6ba521.htm "GetErrorDescription 方法 ")

[GetNationEnumerator 方法](../html/25e4d575-c7fe-6aa7-e005-b748dc4ccda6.htm "GetNationEnumerator 方法 ")

[GetNationText 方法](../html/4010b2cb-f8d6-aaf9-b18a-5313c00c1842.htm "GetNationText 方法 ")

[GetNewNetMessage 方法](../html/55f69034-bd64-b15b-938f-1a381e9cf19f.htm "GetNewNetMessage 方法 ")

[PackToSAMCommand 方法](../html/6dc8c357-7f5c-1711-419b-0cb2ab75fc43.htm "PackToSAMCommand 方法 ")

[ReadCard 方法](../html/3934f7b5-4172-ca5a-8949-4d1ff811da74.htm "ReadCard 方法 ")

[ReadSafeModuleNumber 方法](../html/1ea6f232-1d72-9133-c1e6-9519576160e5.htm "ReadSafeModuleNumber 方法 ")

[SearchCard 方法](../html/4b6a12f3-7bd3-bad9-77ba-40b3a7afacae.htm "SearchCard 方法 ")

[SelectCard 方法](../html/51354f2d-81b3-490a-6b4a-2e1762216d60.htm "SelectCard 方法 ")

[ToString 方法](../html/6b53a0d9-79d5-4509-63d5-dbcbe689c874.htm "ToString 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| SAMSerialCheckSafeModuleStatus 方法 |

检测安全模块状态  
Detecting Security Module Status

**命名空间：**
 [HslCommunication.Profinet.IDCard](54fb02ce-3bfd-9c1d-96f8-511375556490.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public OperateResult CheckSafeModuleStatus()
```

```
Public Function CheckSafeModuleStatus As OperateResult
```

```
public:
OperateResult^ CheckSafeModuleStatus()
```

```
member CheckSafeModuleStatus : unit -> OperateResult 
```

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
返回是否检测成功

![](../icons/SectionExpanded.png)参见

#### 引用

[SAMSerial 类](48f6a433-86f0-a4f8-1751-462e12ad6923.htm)

[HslCommunication.Profinet.IDCard 命名空间](54fb02ce-3bfd-9c1d-96f8-511375556490.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ExtractIdentityCard 方法 

[原文連結](http://api.hslcommunication.cn/html/36f5add1-fb33-d4e2-f96c-312f56f0cddd.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.IDCard](../html/54fb02ce-3bfd-9c1d-96f8-511375556490.htm "HslCommunication.Profinet.IDCard")

[SAMSerial 类](../html/48f6a433-86f0-a4f8-1751-462e12ad6923.htm "SAMSerial 类")

[SAMSerial 方法](../html/3ee3795f-b4e0-f2ff-9050-7961204b2d49.htm "SAMSerial 方法")

[BuildReadCommand 方法](../html/26b20941-6ffd-7305-bdae-3ae3b482428e.htm "BuildReadCommand 方法 ")

[CheckADSCommandAndSum 方法](../html/9381caf4-17a8-3d72-9117-9d7ce75409ea.htm "CheckADSCommandAndSum 方法 ")

[CheckADSCommandCompletion 方法](../html/2ecc327d-7734-26ee-365c-61175d225a17.htm "CheckADSCommandCompletion 方法 ")

[CheckSafeModuleStatus 方法](../html/9ac7a019-82d2-b621-758d-e6d2182e3f4c.htm "CheckSafeModuleStatus 方法 ")

[ExtractIdentityCard 方法](../html/36f5add1-fb33-d4e2-f96c-312f56f0cddd.htm "ExtractIdentityCard 方法 ")

[ExtractSafeModuleNumber 方法](../html/5d7947c2-f718-b70c-9841-86350cfa0905.htm "ExtractSafeModuleNumber 方法 ")

[GetErrorDescription 方法](../html/f2c1ca18-618b-861b-86e5-a9ebaa6ba521.htm "GetErrorDescription 方法 ")

[GetNationEnumerator 方法](../html/25e4d575-c7fe-6aa7-e005-b748dc4ccda6.htm "GetNationEnumerator 方法 ")

[GetNationText 方法](../html/4010b2cb-f8d6-aaf9-b18a-5313c00c1842.htm "GetNationText 方法 ")

[GetNewNetMessage 方法](../html/55f69034-bd64-b15b-938f-1a381e9cf19f.htm "GetNewNetMessage 方法 ")

[PackToSAMCommand 方法](../html/6dc8c357-7f5c-1711-419b-0cb2ab75fc43.htm "PackToSAMCommand 方法 ")

[ReadCard 方法](../html/3934f7b5-4172-ca5a-8949-4d1ff811da74.htm "ReadCard 方法 ")

[ReadSafeModuleNumber 方法](../html/1ea6f232-1d72-9133-c1e6-9519576160e5.htm "ReadSafeModuleNumber 方法 ")

[SearchCard 方法](../html/4b6a12f3-7bd3-bad9-77ba-40b3a7afacae.htm "SearchCard 方法 ")

[SelectCard 方法](../html/51354f2d-81b3-490a-6b4a-2e1762216d60.htm "SelectCard 方法 ")

[ToString 方法](../html/6b53a0d9-79d5-4509-63d5-dbcbe689c874.htm "ToString 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| SAMSerialExtractIdentityCard 方法 |

从数据中提取出真实的身份证信息

**命名空间：**
 [HslCommunication.Profinet.IDCard](54fb02ce-3bfd-9c1d-96f8-511375556490.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult<IdentityCard> ExtractIdentityCard(
	byte[] data
)
```

```
Public Shared Function ExtractIdentityCard ( 
	data As Byte()
) As OperateResult(Of IdentityCard)
```

```
public:
static OperateResult<IdentityCard^>^ ExtractIdentityCard(
	array<unsigned char>^ data
)
```

```
static member ExtractIdentityCard : 
        data : byte[] -> OperateResult<IdentityCard> 
```

#### 参数

data
:   类型：SystemByte  
    原始数据内容

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)[IdentityCard](596fe590-4a15-954f-47da-25d8386a8cde.htm)  
包含结果对象的身份证数据

![](../icons/SectionExpanded.png)参见

#### 引用

[SAMSerial 类](48f6a433-86f0-a4f8-1751-462e12ad6923.htm)

[HslCommunication.Profinet.IDCard 命名空间](54fb02ce-3bfd-9c1d-96f8-511375556490.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ExtractSafeModuleNumber 方法 

[原文連結](http://api.hslcommunication.cn/html/5d7947c2-f718-b70c-9841-86350cfa0905.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.IDCard](../html/54fb02ce-3bfd-9c1d-96f8-511375556490.htm "HslCommunication.Profinet.IDCard")

[SAMSerial 类](../html/48f6a433-86f0-a4f8-1751-462e12ad6923.htm "SAMSerial 类")

[SAMSerial 方法](../html/3ee3795f-b4e0-f2ff-9050-7961204b2d49.htm "SAMSerial 方法")

[BuildReadCommand 方法](../html/26b20941-6ffd-7305-bdae-3ae3b482428e.htm "BuildReadCommand 方法 ")

[CheckADSCommandAndSum 方法](../html/9381caf4-17a8-3d72-9117-9d7ce75409ea.htm "CheckADSCommandAndSum 方法 ")

[CheckADSCommandCompletion 方法](../html/2ecc327d-7734-26ee-365c-61175d225a17.htm "CheckADSCommandCompletion 方法 ")

[CheckSafeModuleStatus 方法](../html/9ac7a019-82d2-b621-758d-e6d2182e3f4c.htm "CheckSafeModuleStatus 方法 ")

[ExtractIdentityCard 方法](../html/36f5add1-fb33-d4e2-f96c-312f56f0cddd.htm "ExtractIdentityCard 方法 ")

[ExtractSafeModuleNumber 方法](../html/5d7947c2-f718-b70c-9841-86350cfa0905.htm "ExtractSafeModuleNumber 方法 ")

[GetErrorDescription 方法](../html/f2c1ca18-618b-861b-86e5-a9ebaa6ba521.htm "GetErrorDescription 方法 ")

[GetNationEnumerator 方法](../html/25e4d575-c7fe-6aa7-e005-b748dc4ccda6.htm "GetNationEnumerator 方法 ")

[GetNationText 方法](../html/4010b2cb-f8d6-aaf9-b18a-5313c00c1842.htm "GetNationText 方法 ")

[GetNewNetMessage 方法](../html/55f69034-bd64-b15b-938f-1a381e9cf19f.htm "GetNewNetMessage 方法 ")

[PackToSAMCommand 方法](../html/6dc8c357-7f5c-1711-419b-0cb2ab75fc43.htm "PackToSAMCommand 方法 ")

[ReadCard 方法](../html/3934f7b5-4172-ca5a-8949-4d1ff811da74.htm "ReadCard 方法 ")

[ReadSafeModuleNumber 方法](../html/1ea6f232-1d72-9133-c1e6-9519576160e5.htm "ReadSafeModuleNumber 方法 ")

[SearchCard 方法](../html/4b6a12f3-7bd3-bad9-77ba-40b3a7afacae.htm "SearchCard 方法 ")

[SelectCard 方法](../html/51354f2d-81b3-490a-6b4a-2e1762216d60.htm "SelectCard 方法 ")

[ToString 方法](../html/6b53a0d9-79d5-4509-63d5-dbcbe689c874.htm "ToString 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| SAMSerialExtractSafeModuleNumber 方法 |

提炼安全的模块数据信息

**命名空间：**
 [HslCommunication.Profinet.IDCard](54fb02ce-3bfd-9c1d-96f8-511375556490.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult<string> ExtractSafeModuleNumber(
	byte[] data
)
```

```
Public Shared Function ExtractSafeModuleNumber ( 
	data As Byte()
) As OperateResult(Of String)
```

```
public:
static OperateResult<String^>^ ExtractSafeModuleNumber(
	array<unsigned char>^ data
)
```

```
static member ExtractSafeModuleNumber : 
        data : byte[] -> OperateResult<string> 
```

#### 参数

data
:   类型：SystemByte  
    数据

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)String  
结果对象

![](../icons/SectionExpanded.png)参见

#### 引用

[SAMSerial 类](48f6a433-86f0-a4f8-1751-462e12ad6923.htm)

[HslCommunication.Profinet.IDCard 命名空间](54fb02ce-3bfd-9c1d-96f8-511375556490.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GetErrorDescription 方法 

[原文連結](http://api.hslcommunication.cn/html/f2c1ca18-618b-861b-86e5-a9ebaa6ba521.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.IDCard](../html/54fb02ce-3bfd-9c1d-96f8-511375556490.htm "HslCommunication.Profinet.IDCard")

[SAMSerial 类](../html/48f6a433-86f0-a4f8-1751-462e12ad6923.htm "SAMSerial 类")

[SAMSerial 方法](../html/3ee3795f-b4e0-f2ff-9050-7961204b2d49.htm "SAMSerial 方法")

[BuildReadCommand 方法](../html/26b20941-6ffd-7305-bdae-3ae3b482428e.htm "BuildReadCommand 方法 ")

[CheckADSCommandAndSum 方法](../html/9381caf4-17a8-3d72-9117-9d7ce75409ea.htm "CheckADSCommandAndSum 方法 ")

[CheckADSCommandCompletion 方法](../html/2ecc327d-7734-26ee-365c-61175d225a17.htm "CheckADSCommandCompletion 方法 ")

[CheckSafeModuleStatus 方法](../html/9ac7a019-82d2-b621-758d-e6d2182e3f4c.htm "CheckSafeModuleStatus 方法 ")

[ExtractIdentityCard 方法](../html/36f5add1-fb33-d4e2-f96c-312f56f0cddd.htm "ExtractIdentityCard 方法 ")

[ExtractSafeModuleNumber 方法](../html/5d7947c2-f718-b70c-9841-86350cfa0905.htm "ExtractSafeModuleNumber 方法 ")

[GetErrorDescription 方法](../html/f2c1ca18-618b-861b-86e5-a9ebaa6ba521.htm "GetErrorDescription 方法 ")

[GetNationEnumerator 方法](../html/25e4d575-c7fe-6aa7-e005-b748dc4ccda6.htm "GetNationEnumerator 方法 ")

[GetNationText 方法](../html/4010b2cb-f8d6-aaf9-b18a-5313c00c1842.htm "GetNationText 方法 ")

[GetNewNetMessage 方法](../html/55f69034-bd64-b15b-938f-1a381e9cf19f.htm "GetNewNetMessage 方法 ")

[PackToSAMCommand 方法](../html/6dc8c357-7f5c-1711-419b-0cb2ab75fc43.htm "PackToSAMCommand 方法 ")

[ReadCard 方法](../html/3934f7b5-4172-ca5a-8949-4d1ff811da74.htm "ReadCard 方法 ")

[ReadSafeModuleNumber 方法](../html/1ea6f232-1d72-9133-c1e6-9519576160e5.htm "ReadSafeModuleNumber 方法 ")

[SearchCard 方法](../html/4b6a12f3-7bd3-bad9-77ba-40b3a7afacae.htm "SearchCard 方法 ")

[SelectCard 方法](../html/51354f2d-81b3-490a-6b4a-2e1762216d60.htm "SelectCard 方法 ")

[ToString 方法](../html/6b53a0d9-79d5-4509-63d5-dbcbe689c874.htm "ToString 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| SAMSerialGetErrorDescription 方法 |

获取错误的文本信息

**命名空间：**
 [HslCommunication.Profinet.IDCard](54fb02ce-3bfd-9c1d-96f8-511375556490.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static string GetErrorDescription(
	int err
)
```

```
Public Shared Function GetErrorDescription ( 
	err As Integer
) As String
```

```
public:
static String^ GetErrorDescription(
	int err
)
```

```
static member GetErrorDescription : 
        err : int -> string 
```

#### 参数

err
:   类型：SystemInt32  
    错误号

#### 返回值

类型：String  
错误信息

![](../icons/SectionExpanded.png)参见

#### 引用

[SAMSerial 类](48f6a433-86f0-a4f8-1751-462e12ad6923.htm)

[HslCommunication.Profinet.IDCard 命名空间](54fb02ce-3bfd-9c1d-96f8-511375556490.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GetNationEnumerator 方法 

[原文連結](http://api.hslcommunication.cn/html/25e4d575-c7fe-6aa7-e005-b748dc4ccda6.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.IDCard](../html/54fb02ce-3bfd-9c1d-96f8-511375556490.htm "HslCommunication.Profinet.IDCard")

[SAMSerial 类](../html/48f6a433-86f0-a4f8-1751-462e12ad6923.htm "SAMSerial 类")

[SAMSerial 方法](../html/3ee3795f-b4e0-f2ff-9050-7961204b2d49.htm "SAMSerial 方法")

[BuildReadCommand 方法](../html/26b20941-6ffd-7305-bdae-3ae3b482428e.htm "BuildReadCommand 方法 ")

[CheckADSCommandAndSum 方法](../html/9381caf4-17a8-3d72-9117-9d7ce75409ea.htm "CheckADSCommandAndSum 方法 ")

[CheckADSCommandCompletion 方法](../html/2ecc327d-7734-26ee-365c-61175d225a17.htm "CheckADSCommandCompletion 方法 ")

[CheckSafeModuleStatus 方法](../html/9ac7a019-82d2-b621-758d-e6d2182e3f4c.htm "CheckSafeModuleStatus 方法 ")

[ExtractIdentityCard 方法](../html/36f5add1-fb33-d4e2-f96c-312f56f0cddd.htm "ExtractIdentityCard 方法 ")

[ExtractSafeModuleNumber 方法](../html/5d7947c2-f718-b70c-9841-86350cfa0905.htm "ExtractSafeModuleNumber 方法 ")

[GetErrorDescription 方法](../html/f2c1ca18-618b-861b-86e5-a9ebaa6ba521.htm "GetErrorDescription 方法 ")

[GetNationEnumerator 方法](../html/25e4d575-c7fe-6aa7-e005-b748dc4ccda6.htm "GetNationEnumerator 方法 ")

[GetNationText 方法](../html/4010b2cb-f8d6-aaf9-b18a-5313c00c1842.htm "GetNationText 方法 ")

[GetNewNetMessage 方法](../html/55f69034-bd64-b15b-938f-1a381e9cf19f.htm "GetNewNetMessage 方法 ")

[PackToSAMCommand 方法](../html/6dc8c357-7f5c-1711-419b-0cb2ab75fc43.htm "PackToSAMCommand 方法 ")

[ReadCard 方法](../html/3934f7b5-4172-ca5a-8949-4d1ff811da74.htm "ReadCard 方法 ")

[ReadSafeModuleNumber 方法](../html/1ea6f232-1d72-9133-c1e6-9519576160e5.htm "ReadSafeModuleNumber 方法 ")

[SearchCard 方法](../html/4b6a12f3-7bd3-bad9-77ba-40b3a7afacae.htm "SearchCard 方法 ")

[SelectCard 方法](../html/51354f2d-81b3-490a-6b4a-2e1762216d60.htm "SelectCard 方法 ")

[ToString 方法](../html/6b53a0d9-79d5-4509-63d5-dbcbe689c874.htm "ToString 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| SAMSerialGetNationEnumerator 方法 |

枚举当前的所有的民族信息，共计五十六个民族

**命名空间：**
 [HslCommunication.Profinet.IDCard](54fb02ce-3bfd-9c1d-96f8-511375556490.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static IEnumerator<string> GetNationEnumerator()
```

```
Public Shared Function GetNationEnumerator As IEnumerator(Of String)
```

```
public:
static IEnumerator<String^>^ GetNationEnumerator()
```

```
static member GetNationEnumerator : unit -> IEnumerator<string> 
```

#### 返回值

类型：IEnumeratorString  
枚举信息

![](../icons/SectionExpanded.png)参见

#### 引用

[SAMSerial 类](48f6a433-86f0-a4f8-1751-462e12ad6923.htm)

[HslCommunication.Profinet.IDCard 命名空间](54fb02ce-3bfd-9c1d-96f8-511375556490.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GetNationText 方法 

[原文連結](http://api.hslcommunication.cn/html/4010b2cb-f8d6-aaf9-b18a-5313c00c1842.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.IDCard](../html/54fb02ce-3bfd-9c1d-96f8-511375556490.htm "HslCommunication.Profinet.IDCard")

[SAMSerial 类](../html/48f6a433-86f0-a4f8-1751-462e12ad6923.htm "SAMSerial 类")

[SAMSerial 方法](../html/3ee3795f-b4e0-f2ff-9050-7961204b2d49.htm "SAMSerial 方法")

[BuildReadCommand 方法](../html/26b20941-6ffd-7305-bdae-3ae3b482428e.htm "BuildReadCommand 方法 ")

[CheckADSCommandAndSum 方法](../html/9381caf4-17a8-3d72-9117-9d7ce75409ea.htm "CheckADSCommandAndSum 方法 ")

[CheckADSCommandCompletion 方法](../html/2ecc327d-7734-26ee-365c-61175d225a17.htm "CheckADSCommandCompletion 方法 ")

[CheckSafeModuleStatus 方法](../html/9ac7a019-82d2-b621-758d-e6d2182e3f4c.htm "CheckSafeModuleStatus 方法 ")

[ExtractIdentityCard 方法](../html/36f5add1-fb33-d4e2-f96c-312f56f0cddd.htm "ExtractIdentityCard 方法 ")

[ExtractSafeModuleNumber 方法](../html/5d7947c2-f718-b70c-9841-86350cfa0905.htm "ExtractSafeModuleNumber 方法 ")

[GetErrorDescription 方法](../html/f2c1ca18-618b-861b-86e5-a9ebaa6ba521.htm "GetErrorDescription 方法 ")

[GetNationEnumerator 方法](../html/25e4d575-c7fe-6aa7-e005-b748dc4ccda6.htm "GetNationEnumerator 方法 ")

[GetNationText 方法](../html/4010b2cb-f8d6-aaf9-b18a-5313c00c1842.htm "GetNationText 方法 ")

[GetNewNetMessage 方法](../html/55f69034-bd64-b15b-938f-1a381e9cf19f.htm "GetNewNetMessage 方法 ")

[PackToSAMCommand 方法](../html/6dc8c357-7f5c-1711-419b-0cb2ab75fc43.htm "PackToSAMCommand 方法 ")

[ReadCard 方法](../html/3934f7b5-4172-ca5a-8949-4d1ff811da74.htm "ReadCard 方法 ")

[ReadSafeModuleNumber 方法](../html/1ea6f232-1d72-9133-c1e6-9519576160e5.htm "ReadSafeModuleNumber 方法 ")

[SearchCard 方法](../html/4b6a12f3-7bd3-bad9-77ba-40b3a7afacae.htm "SearchCard 方法 ")

[SelectCard 方法](../html/51354f2d-81b3-490a-6b4a-2e1762216d60.htm "SelectCard 方法 ")

[ToString 方法](../html/6b53a0d9-79d5-4509-63d5-dbcbe689c874.htm "ToString 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| SAMSerialGetNationText 方法 |

根据民族的代号来获取到民族的文本描述信息

**命名空间：**
 [HslCommunication.Profinet.IDCard](54fb02ce-3bfd-9c1d-96f8-511375556490.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static string GetNationText(
	int nation
)
```

```
Public Shared Function GetNationText ( 
	nation As Integer
) As String
```

```
public:
static String^ GetNationText(
	int nation
)
```

```
static member GetNationText : 
        nation : int -> string 
```

#### 参数

nation
:   类型：SystemInt32  
    民族代码

#### 返回值

类型：String  
民族的文本信息

![](../icons/SectionExpanded.png)参见

#### 引用

[SAMSerial 类](48f6a433-86f0-a4f8-1751-462e12ad6923.htm)

[HslCommunication.Profinet.IDCard 命名空间](54fb02ce-3bfd-9c1d-96f8-511375556490.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GetNewNetMessage 方法 

[原文連結](http://api.hslcommunication.cn/html/55f69034-bd64-b15b-938f-1a381e9cf19f.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.IDCard](../html/54fb02ce-3bfd-9c1d-96f8-511375556490.htm "HslCommunication.Profinet.IDCard")

[SAMSerial 类](../html/48f6a433-86f0-a4f8-1751-462e12ad6923.htm "SAMSerial 类")

[SAMSerial 方法](../html/3ee3795f-b4e0-f2ff-9050-7961204b2d49.htm "SAMSerial 方法")

[BuildReadCommand 方法](../html/26b20941-6ffd-7305-bdae-3ae3b482428e.htm "BuildReadCommand 方法 ")

[CheckADSCommandAndSum 方法](../html/9381caf4-17a8-3d72-9117-9d7ce75409ea.htm "CheckADSCommandAndSum 方法 ")

[CheckADSCommandCompletion 方法](../html/2ecc327d-7734-26ee-365c-61175d225a17.htm "CheckADSCommandCompletion 方法 ")

[CheckSafeModuleStatus 方法](../html/9ac7a019-82d2-b621-758d-e6d2182e3f4c.htm "CheckSafeModuleStatus 方法 ")

[ExtractIdentityCard 方法](../html/36f5add1-fb33-d4e2-f96c-312f56f0cddd.htm "ExtractIdentityCard 方法 ")

[ExtractSafeModuleNumber 方法](../html/5d7947c2-f718-b70c-9841-86350cfa0905.htm "ExtractSafeModuleNumber 方法 ")

[GetErrorDescription 方法](../html/f2c1ca18-618b-861b-86e5-a9ebaa6ba521.htm "GetErrorDescription 方法 ")

[GetNationEnumerator 方法](../html/25e4d575-c7fe-6aa7-e005-b748dc4ccda6.htm "GetNationEnumerator 方法 ")

[GetNationText 方法](../html/4010b2cb-f8d6-aaf9-b18a-5313c00c1842.htm "GetNationText 方法 ")

[GetNewNetMessage 方法](../html/55f69034-bd64-b15b-938f-1a381e9cf19f.htm "GetNewNetMessage 方法 ")

[PackToSAMCommand 方法](../html/6dc8c357-7f5c-1711-419b-0cb2ab75fc43.htm "PackToSAMCommand 方法 ")

[ReadCard 方法](../html/3934f7b5-4172-ca5a-8949-4d1ff811da74.htm "ReadCard 方法 ")

[ReadSafeModuleNumber 方法](../html/1ea6f232-1d72-9133-c1e6-9519576160e5.htm "ReadSafeModuleNumber 方法 ")

[SearchCard 方法](../html/4b6a12f3-7bd3-bad9-77ba-40b3a7afacae.htm "SearchCard 方法 ")

[SelectCard 方法](../html/51354f2d-81b3-490a-6b4a-2e1762216d60.htm "SelectCard 方法 ")

[ToString 方法](../html/6b53a0d9-79d5-4509-63d5-dbcbe689c874.htm "ToString 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| SAMSerialGetNewNetMessage 方法 |

获取一个新的消息对象的方法，需要在继承类里面进行重写  
The method to get a new message object needs to be overridden in the inheritance class

**命名空间：**
 [HslCommunication.Profinet.IDCard](54fb02ce-3bfd-9c1d-96f8-511375556490.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
protected override INetMessage GetNewNetMessage()
```

```
Protected Overrides Function GetNewNetMessage As INetMessage
```

```
protected:
virtual INetMessage^ GetNewNetMessage() override
```

```
abstract GetNewNetMessage : unit -> INetMessage 
override GetNewNetMessage : unit -> INetMessage
```

#### 返回值

类型：[INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm)  
消息类对象

![](../icons/SectionExpanded.png)参见

#### 引用

[SAMSerial 类](48f6a433-86f0-a4f8-1751-462e12ad6923.htm)

[HslCommunication.Profinet.IDCard 命名空间](54fb02ce-3bfd-9c1d-96f8-511375556490.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)