# HslCommunication - HslCommunication.Profinet.Beckhoff

> 分類頁數: 30



---
## HslCommunication.Profinet.Beckhoff

[原文連結](http://api.hslcommunication.cn/html/fb2af334-3662-b9ca-be3c-b0090bf4ef56.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Beckhoff](../html/fb2af334-3662-b9ca-be3c-b0090bf4ef56.htm "HslCommunication.Profinet.Beckhoff")

[AdsDeviceInfo 类](../html/2e264493-3101-edad-3c48-3a954214b540.htm "AdsDeviceInfo 类")

[AmsTcpHeaderFlags 枚举](../html/f0aebe48-7132-acf5-afca-c6307c90bddb.htm "AmsTcpHeaderFlags 枚举")

[BeckhoffAdsNet 类](../html/7ab15712-8688-063d-6b30-4b24803393fd.htm "BeckhoffAdsNet 类")

[BeckhoffAdsServer 类](../html/e199f7e7-84be-6ff6-8ad7-08e55c5e86da.htm "BeckhoffAdsServer 类")

[BeckhoffCommandId 类](../html/78f1c59b-6b0b-2c97-353f-7fc45ec00a8d.htm "BeckhoffCommandId 类")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslCommunication.Profinet.Beckhoff 命名空间 |

[缺少 "N:HslCommunication.Profinet.Beckhoff" 的 <summary> 文档]

![](../icons/SectionExpanded.png)类

|  | 类 | 说明 |
| --- | --- | --- |
| 公共类 | [AdsDeviceInfo](2e264493-3101-edad-3c48-3a954214b540.htm) | Ads设备的相关信息，主要是版本号，设备名称  Information about Ads devices, primarily version numbers, device names. |
| 公共类代码示例 | [BeckhoffAdsNet](7ab15712-8688-063d-6b30-4b24803393fd.htm) | 倍福的ADS协议，支持读取倍福的地址数据，关于端口号的选择，TwinCAT2，端口号801；TwinCAT3，端口号为851，NETID可以选择手动输入，自动输入方式，具体参考API文档的示例代码  Beckhoff's ADS protocol supports reading Beckhoff address data. Regarding the choice of port number, TwinCAT2, port number is 801; TwinCAT3, port number is 851, NETID can be input manually or automatically. For details, please refer to the example of API documentation code |
| 公共类 | [BeckhoffAdsServer](e199f7e7-84be-6ff6-8ad7-08e55c5e86da.htm) | 倍福Ads协议的虚拟服务器 |
| 公共类 | [BeckhoffCommandId](78f1c59b-6b0b-2c97-353f-7fc45ec00a8d.htm) | 倍福PLC的命令码  Command Id: https://infosys.beckhoff.com/english.php?content=../content/1033/tc3\_ads\_intro/index.html |

![](../icons/SectionExpanded.png)枚举

|  | 枚举 | 说明 |
| --- | --- | --- |
| 公共枚举 | [AmsTcpHeaderFlags](f0aebe48-7132-acf5-afca-c6307c90bddb.htm) | AMS消息的命令号 |

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## AdsDeviceInfo 类

[原文連結](http://api.hslcommunication.cn/html/2e264493-3101-edad-3c48-3a954214b540.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Beckhoff](../html/fb2af334-3662-b9ca-be3c-b0090bf4ef56.htm "HslCommunication.Profinet.Beckhoff")

[AdsDeviceInfo 类](../html/2e264493-3101-edad-3c48-3a954214b540.htm "AdsDeviceInfo 类")

[AdsDeviceInfo 构造函数](../html/201c7a60-3927-086a-3b53-221152e59696.htm "AdsDeviceInfo 构造函数 ")

[AdsDeviceInfo 属性](../html/6b0ae67c-7e16-d9b8-90a0-6c45ee42f395.htm "AdsDeviceInfo 属性")

[AdsDeviceInfo 方法](../html/6c8103c3-256c-e903-dd0d-7ee8119d6343.htm "AdsDeviceInfo 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| AdsDeviceInfo 类 |

Ads设备的相关信息，主要是版本号，设备名称  
Information about Ads devices, primarily version numbers, device names.

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  HslCommunication.Profinet.BeckhoffAdsDeviceInfo

**命名空间：**
 [HslCommunication.Profinet.Beckhoff](fb2af334-3662-b9ca-be3c-b0090bf4ef56.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class AdsDeviceInfo
```

```
Public Class AdsDeviceInfo
```

```
public ref class AdsDeviceInfo
```

```
type AdsDeviceInfo =  class end
```

AdsDeviceInfo 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [AdsDeviceInfo](c0c7a78d-d82a-7e22-59cd-a03c71f4f561.htm) | 实例化一个默认的对象  Instantiate a default object |
| 公共方法 | [AdsDeviceInfo(Byte)](6bbffa05-231e-4a94-c3d5-8b2bf0bb625a.htm) | 根据原始的数据内容来实例化一个对象  Instantiate an object based on the original data content |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [Build](5d5a0608-e89c-983b-8cf6-aee3489c2e9e.htm) | 构建版本号  Build version |
| 公共属性 | [DeviceName](76eb7925-b5c4-83e0-6a94-f76d611ac840.htm) | 设备的名字  Device Name |
| 公共属性 | [Major](20f34da6-e2bd-f944-e49b-17c242db7a05.htm) | 主版本号  Main Version |
| 公共属性 | [Minor](35426b5c-bc98-17e2-b9e7-3926ea7f5bf8.htm) | 次版本号  Minor Version |

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

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.Profinet.Beckhoff 命名空间](fb2af334-3662-b9ca-be3c-b0090bf4ef56.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## AdsDeviceInfo 构造函数 

[原文連結](http://api.hslcommunication.cn/html/201c7a60-3927-086a-3b53-221152e59696.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Beckhoff](../html/fb2af334-3662-b9ca-be3c-b0090bf4ef56.htm "HslCommunication.Profinet.Beckhoff")

[AdsDeviceInfo 类](../html/2e264493-3101-edad-3c48-3a954214b540.htm "AdsDeviceInfo 类")

[AdsDeviceInfo 构造函数](../html/201c7a60-3927-086a-3b53-221152e59696.htm "AdsDeviceInfo 构造函数 ")

[AdsDeviceInfo 构造函数](../html/c0c7a78d-d82a-7e22-59cd-a03c71f4f561.htm "AdsDeviceInfo 构造函数 ")

[AdsDeviceInfo 构造函数 (Byte[])](../html/6bbffa05-231e-4a94-c3d5-8b2bf0bb625a.htm "AdsDeviceInfo 构造函数 (Byte[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| AdsDeviceInfo 构造函数 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [AdsDeviceInfo](c0c7a78d-d82a-7e22-59cd-a03c71f4f561.htm) | 实例化一个默认的对象  Instantiate a default object |
| 公共方法 | [AdsDeviceInfo(Byte)](6bbffa05-231e-4a94-c3d5-8b2bf0bb625a.htm) | 根据原始的数据内容来实例化一个对象  Instantiate an object based on the original data content |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[AdsDeviceInfo 类](2e264493-3101-edad-3c48-3a954214b540.htm)

[HslCommunication.Profinet.Beckhoff 命名空间](fb2af334-3662-b9ca-be3c-b0090bf4ef56.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## AdsDeviceInfo 构造函数 

[原文連結](http://api.hslcommunication.cn/html/c0c7a78d-d82a-7e22-59cd-a03c71f4f561.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Beckhoff](../html/fb2af334-3662-b9ca-be3c-b0090bf4ef56.htm "HslCommunication.Profinet.Beckhoff")

[AdsDeviceInfo 类](../html/2e264493-3101-edad-3c48-3a954214b540.htm "AdsDeviceInfo 类")

[AdsDeviceInfo 构造函数](../html/201c7a60-3927-086a-3b53-221152e59696.htm "AdsDeviceInfo 构造函数 ")

[AdsDeviceInfo 构造函数](../html/c0c7a78d-d82a-7e22-59cd-a03c71f4f561.htm "AdsDeviceInfo 构造函数 ")

[AdsDeviceInfo 构造函数 (Byte[])](../html/6bbffa05-231e-4a94-c3d5-8b2bf0bb625a.htm "AdsDeviceInfo 构造函数 (Byte[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| AdsDeviceInfo 构造函数 |

实例化一个默认的对象  
Instantiate a default object

**命名空间：**
 [HslCommunication.Profinet.Beckhoff](fb2af334-3662-b9ca-be3c-b0090bf4ef56.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public AdsDeviceInfo()
```

```
Public Sub New
```

```
public:
AdsDeviceInfo()
```

```
new : unit -> AdsDeviceInfo
```

![](../icons/SectionExpanded.png)参见

#### 引用

[AdsDeviceInfo 类](2e264493-3101-edad-3c48-3a954214b540.htm)

[AdsDeviceInfo 重载](201c7a60-3927-086a-3b53-221152e59696.htm)

[HslCommunication.Profinet.Beckhoff 命名空间](fb2af334-3662-b9ca-be3c-b0090bf4ef56.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## AdsDeviceInfo 构造函数 (Byte[])

[原文連結](http://api.hslcommunication.cn/html/6bbffa05-231e-4a94-c3d5-8b2bf0bb625a.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Beckhoff](../html/fb2af334-3662-b9ca-be3c-b0090bf4ef56.htm "HslCommunication.Profinet.Beckhoff")

[AdsDeviceInfo 类](../html/2e264493-3101-edad-3c48-3a954214b540.htm "AdsDeviceInfo 类")

[AdsDeviceInfo 构造函数](../html/201c7a60-3927-086a-3b53-221152e59696.htm "AdsDeviceInfo 构造函数 ")

[AdsDeviceInfo 构造函数](../html/c0c7a78d-d82a-7e22-59cd-a03c71f4f561.htm "AdsDeviceInfo 构造函数 ")

[AdsDeviceInfo 构造函数 (Byte[])](../html/6bbffa05-231e-4a94-c3d5-8b2bf0bb625a.htm "AdsDeviceInfo 构造函数 (Byte[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| AdsDeviceInfo 构造函数 (Byte) |

根据原始的数据内容来实例化一个对象  
Instantiate an object based on the original data content

**命名空间：**
 [HslCommunication.Profinet.Beckhoff](fb2af334-3662-b9ca-be3c-b0090bf4ef56.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public AdsDeviceInfo(
	byte[] data
)
```

```
Public Sub New ( 
	data As Byte()
)
```

```
public:
AdsDeviceInfo(
	array<unsigned char>^ data
)
```

```
new : 
        data : byte[] -> AdsDeviceInfo
```

#### 参数

data
:   类型：SystemByte  
    原始的数据内容

![](../icons/SectionExpanded.png)参见

#### 引用

[AdsDeviceInfo 类](2e264493-3101-edad-3c48-3a954214b540.htm)

[AdsDeviceInfo 重载](201c7a60-3927-086a-3b53-221152e59696.htm)

[HslCommunication.Profinet.Beckhoff 命名空间](fb2af334-3662-b9ca-be3c-b0090bf4ef56.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## AdsDeviceInfo 属性

[原文連結](http://api.hslcommunication.cn/html/6b0ae67c-7e16-d9b8-90a0-6c45ee42f395.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Beckhoff](../html/fb2af334-3662-b9ca-be3c-b0090bf4ef56.htm "HslCommunication.Profinet.Beckhoff")

[AdsDeviceInfo 类](../html/2e264493-3101-edad-3c48-3a954214b540.htm "AdsDeviceInfo 类")

[AdsDeviceInfo 属性](../html/6b0ae67c-7e16-d9b8-90a0-6c45ee42f395.htm "AdsDeviceInfo 属性")

[Build 属性](../html/5d5a0608-e89c-983b-8cf6-aee3489c2e9e.htm "Build 属性 ")

[DeviceName 属性](../html/76eb7925-b5c4-83e0-6a94-f76d611ac840.htm "DeviceName 属性 ")

[Major 属性](../html/20f34da6-e2bd-f944-e49b-17c242db7a05.htm "Major 属性 ")

[Minor 属性](../html/35426b5c-bc98-17e2-b9e7-3926ea7f5bf8.htm "Minor 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| AdsDeviceInfo 属性 |

[AdsDeviceInfo](2e264493-3101-edad-3c48-3a954214b540.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [Build](5d5a0608-e89c-983b-8cf6-aee3489c2e9e.htm) | 构建版本号  Build version |
| 公共属性 | [DeviceName](76eb7925-b5c4-83e0-6a94-f76d611ac840.htm) | 设备的名字  Device Name |
| 公共属性 | [Major](20f34da6-e2bd-f944-e49b-17c242db7a05.htm) | 主版本号  Main Version |
| 公共属性 | [Minor](35426b5c-bc98-17e2-b9e7-3926ea7f5bf8.htm) | 次版本号  Minor Version |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[AdsDeviceInfo 类](2e264493-3101-edad-3c48-3a954214b540.htm)

[HslCommunication.Profinet.Beckhoff 命名空间](fb2af334-3662-b9ca-be3c-b0090bf4ef56.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Build 属性 

[原文連結](http://api.hslcommunication.cn/html/5d5a0608-e89c-983b-8cf6-aee3489c2e9e.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Beckhoff](../html/fb2af334-3662-b9ca-be3c-b0090bf4ef56.htm "HslCommunication.Profinet.Beckhoff")

[AdsDeviceInfo 类](../html/2e264493-3101-edad-3c48-3a954214b540.htm "AdsDeviceInfo 类")

[AdsDeviceInfo 属性](../html/6b0ae67c-7e16-d9b8-90a0-6c45ee42f395.htm "AdsDeviceInfo 属性")

[Build 属性](../html/5d5a0608-e89c-983b-8cf6-aee3489c2e9e.htm "Build 属性 ")

[DeviceName 属性](../html/76eb7925-b5c4-83e0-6a94-f76d611ac840.htm "DeviceName 属性 ")

[Major 属性](../html/20f34da6-e2bd-f944-e49b-17c242db7a05.htm "Major 属性 ")

[Minor 属性](../html/35426b5c-bc98-17e2-b9e7-3926ea7f5bf8.htm "Minor 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| AdsDeviceInfoBuild 属性 |

构建版本号  
Build version

**命名空间：**
 [HslCommunication.Profinet.Beckhoff](fb2af334-3662-b9ca-be3c-b0090bf4ef56.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public ushort Build { get; set; }
```

```
Public Property Build As UShort
	Get
	Set
```

```
public:
property unsigned short Build {
	unsigned short get ();
	void set (unsigned short value);
}
```

```
member Build : uint16 with get, set
```

#### 属性值

类型：UInt16

![](../icons/SectionExpanded.png)参见

#### 引用

[AdsDeviceInfo 类](2e264493-3101-edad-3c48-3a954214b540.htm)

[HslCommunication.Profinet.Beckhoff 命名空间](fb2af334-3662-b9ca-be3c-b0090bf4ef56.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## DeviceName 属性 

[原文連結](http://api.hslcommunication.cn/html/76eb7925-b5c4-83e0-6a94-f76d611ac840.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Beckhoff](../html/fb2af334-3662-b9ca-be3c-b0090bf4ef56.htm "HslCommunication.Profinet.Beckhoff")

[AdsDeviceInfo 类](../html/2e264493-3101-edad-3c48-3a954214b540.htm "AdsDeviceInfo 类")

[AdsDeviceInfo 属性](../html/6b0ae67c-7e16-d9b8-90a0-6c45ee42f395.htm "AdsDeviceInfo 属性")

[Build 属性](../html/5d5a0608-e89c-983b-8cf6-aee3489c2e9e.htm "Build 属性 ")

[DeviceName 属性](../html/76eb7925-b5c4-83e0-6a94-f76d611ac840.htm "DeviceName 属性 ")

[Major 属性](../html/20f34da6-e2bd-f944-e49b-17c242db7a05.htm "Major 属性 ")

[Minor 属性](../html/35426b5c-bc98-17e2-b9e7-3926ea7f5bf8.htm "Minor 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| AdsDeviceInfoDeviceName 属性 |

设备的名字  
Device Name

**命名空间：**
 [HslCommunication.Profinet.Beckhoff](fb2af334-3662-b9ca-be3c-b0090bf4ef56.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public string DeviceName { get; set; }
```

```
Public Property DeviceName As String
	Get
	Set
```

```
public:
property String^ DeviceName {
	String^ get ();
	void set (String^ value);
}
```

```
member DeviceName : string with get, set
```

#### 属性值

类型：String

![](../icons/SectionExpanded.png)参见

#### 引用

[AdsDeviceInfo 类](2e264493-3101-edad-3c48-3a954214b540.htm)

[HslCommunication.Profinet.Beckhoff 命名空间](fb2af334-3662-b9ca-be3c-b0090bf4ef56.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Major 属性 

[原文連結](http://api.hslcommunication.cn/html/20f34da6-e2bd-f944-e49b-17c242db7a05.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Beckhoff](../html/fb2af334-3662-b9ca-be3c-b0090bf4ef56.htm "HslCommunication.Profinet.Beckhoff")

[AdsDeviceInfo 类](../html/2e264493-3101-edad-3c48-3a954214b540.htm "AdsDeviceInfo 类")

[AdsDeviceInfo 属性](../html/6b0ae67c-7e16-d9b8-90a0-6c45ee42f395.htm "AdsDeviceInfo 属性")

[Build 属性](../html/5d5a0608-e89c-983b-8cf6-aee3489c2e9e.htm "Build 属性 ")

[DeviceName 属性](../html/76eb7925-b5c4-83e0-6a94-f76d611ac840.htm "DeviceName 属性 ")

[Major 属性](../html/20f34da6-e2bd-f944-e49b-17c242db7a05.htm "Major 属性 ")

[Minor 属性](../html/35426b5c-bc98-17e2-b9e7-3926ea7f5bf8.htm "Minor 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| AdsDeviceInfoMajor 属性 |

主版本号  
Main Version

**命名空间：**
 [HslCommunication.Profinet.Beckhoff](fb2af334-3662-b9ca-be3c-b0090bf4ef56.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public byte Major { get; set; }
```

```
Public Property Major As Byte
	Get
	Set
```

```
public:
property unsigned char Major {
	unsigned char get ();
	void set (unsigned char value);
}
```

```
member Major : byte with get, set
```

#### 属性值

类型：Byte

![](../icons/SectionExpanded.png)参见

#### 引用

[AdsDeviceInfo 类](2e264493-3101-edad-3c48-3a954214b540.htm)

[HslCommunication.Profinet.Beckhoff 命名空间](fb2af334-3662-b9ca-be3c-b0090bf4ef56.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Minor 属性 

[原文連結](http://api.hslcommunication.cn/html/35426b5c-bc98-17e2-b9e7-3926ea7f5bf8.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Beckhoff](../html/fb2af334-3662-b9ca-be3c-b0090bf4ef56.htm "HslCommunication.Profinet.Beckhoff")

[AdsDeviceInfo 类](../html/2e264493-3101-edad-3c48-3a954214b540.htm "AdsDeviceInfo 类")

[AdsDeviceInfo 属性](../html/6b0ae67c-7e16-d9b8-90a0-6c45ee42f395.htm "AdsDeviceInfo 属性")

[Build 属性](../html/5d5a0608-e89c-983b-8cf6-aee3489c2e9e.htm "Build 属性 ")

[DeviceName 属性](../html/76eb7925-b5c4-83e0-6a94-f76d611ac840.htm "DeviceName 属性 ")

[Major 属性](../html/20f34da6-e2bd-f944-e49b-17c242db7a05.htm "Major 属性 ")

[Minor 属性](../html/35426b5c-bc98-17e2-b9e7-3926ea7f5bf8.htm "Minor 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| AdsDeviceInfoMinor 属性 |

次版本号  
Minor Version

**命名空间：**
 [HslCommunication.Profinet.Beckhoff](fb2af334-3662-b9ca-be3c-b0090bf4ef56.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public byte Minor { get; set; }
```

```
Public Property Minor As Byte
	Get
	Set
```

```
public:
property unsigned char Minor {
	unsigned char get ();
	void set (unsigned char value);
}
```

```
member Minor : byte with get, set
```

#### 属性值

类型：Byte

![](../icons/SectionExpanded.png)参见

#### 引用

[AdsDeviceInfo 类](2e264493-3101-edad-3c48-3a954214b540.htm)

[HslCommunication.Profinet.Beckhoff 命名空间](fb2af334-3662-b9ca-be3c-b0090bf4ef56.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## AdsDeviceInfo 方法

[原文連結](http://api.hslcommunication.cn/html/6c8103c3-256c-e903-dd0d-7ee8119d6343.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Beckhoff](../html/fb2af334-3662-b9ca-be3c-b0090bf4ef56.htm "HslCommunication.Profinet.Beckhoff")

[AdsDeviceInfo 类](../html/2e264493-3101-edad-3c48-3a954214b540.htm "AdsDeviceInfo 类")

[AdsDeviceInfo 构造函数](../html/201c7a60-3927-086a-3b53-221152e59696.htm "AdsDeviceInfo 构造函数 ")

[AdsDeviceInfo 属性](../html/6b0ae67c-7e16-d9b8-90a0-6c45ee42f395.htm "AdsDeviceInfo 属性")

[AdsDeviceInfo 方法](../html/6c8103c3-256c-e903-dd0d-7ee8119d6343.htm "AdsDeviceInfo 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| AdsDeviceInfo 方法 |

[AdsDeviceInfo](2e264493-3101-edad-3c48-3a954214b540.htm) 类型公开以下成员。

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

[AdsDeviceInfo 类](2e264493-3101-edad-3c48-3a954214b540.htm)

[HslCommunication.Profinet.Beckhoff 命名空间](fb2af334-3662-b9ca-be3c-b0090bf4ef56.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## AmsTcpHeaderFlags 枚举

[原文連結](http://api.hslcommunication.cn/html/f0aebe48-7132-acf5-afca-c6307c90bddb.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Beckhoff](../html/fb2af334-3662-b9ca-be3c-b0090bf4ef56.htm "HslCommunication.Profinet.Beckhoff")

[AdsDeviceInfo 类](../html/2e264493-3101-edad-3c48-3a954214b540.htm "AdsDeviceInfo 类")

[AmsTcpHeaderFlags 枚举](../html/f0aebe48-7132-acf5-afca-c6307c90bddb.htm "AmsTcpHeaderFlags 枚举")

[BeckhoffAdsNet 类](../html/7ab15712-8688-063d-6b30-4b24803393fd.htm "BeckhoffAdsNet 类")

[BeckhoffAdsServer 类](../html/e199f7e7-84be-6ff6-8ad7-08e55c5e86da.htm "BeckhoffAdsServer 类")

[BeckhoffCommandId 类](../html/78f1c59b-6b0b-2c97-353f-7fc45ec00a8d.htm "BeckhoffCommandId 类")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| AmsTcpHeaderFlags 枚举 |

AMS消息的命令号

**命名空间：**
 [HslCommunication.Profinet.Beckhoff](fb2af334-3662-b9ca-be3c-b0090bf4ef56.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public enum AmsTcpHeaderFlags
```

```
Public Enumeration AmsTcpHeaderFlags
```

```
public enum class AmsTcpHeaderFlags
```

```
type AmsTcpHeaderFlags
```

![](../icons/SectionExpanded.png)成员

|  | 成员名称 | 值 | 说明 |
| --- | --- | --- | --- |
|  | Command | 0 | AmsCommand (AMS\_TCP\_PORT\_AMS\_CMD, 0x0000) |
|  | PortClose | 1 | Port Close command (AMS\_TCP\_PORT\_CLOSE, 0x0001) |
|  | PortConnect | 4096 | Port connect command (AMS\_TCP\_PORT\_CONNECT, 0x1000) |
|  | RouterNotification | 4097 | Router Notification (AMS\_TCP\_PORT\_ROUTER\_NOTE, 0x1001) |
|  | GetLocalNetId | 4098 | Get LocalNetId header |

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.Profinet.Beckhoff 命名空间](fb2af334-3662-b9ca-be3c-b0090bf4ef56.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## BeckhoffAdsNet 类

[原文連結](http://api.hslcommunication.cn/html/7ab15712-8688-063d-6b30-4b24803393fd.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Beckhoff](../html/fb2af334-3662-b9ca-be3c-b0090bf4ef56.htm "HslCommunication.Profinet.Beckhoff")

[BeckhoffAdsNet 类](../html/7ab15712-8688-063d-6b30-4b24803393fd.htm "BeckhoffAdsNet 类")

[BeckhoffAdsNet 构造函数](../html/eb96f2a7-8a9b-62bb-52e2-172e65d2d879.htm "BeckhoffAdsNet 构造函数 ")

[BeckhoffAdsNet 属性](../html/7ce17828-db34-c0f5-016c-ec693f153a7f.htm "BeckhoffAdsNet 属性")

[BeckhoffAdsNet 方法](../html/e8bad1d3-d4ae-37ef-6ba7-e199dac65f0c.htm "BeckhoffAdsNet 方法")

[BeckhoffAdsNet 字段](../html/c861230e-e3b4-d06d-303d-1c75f710bcfc.htm "BeckhoffAdsNet 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| BeckhoffAdsNet 类 |

倍福的ADS协议，支持读取倍福的地址数据，关于端口号的选择，TwinCAT2，端口号801；TwinCAT3，端口号为851，NETID可以选择手动输入，自动输入方式，具体参考API文档的示例代码  
Beckhoff's ADS protocol supports reading Beckhoff address data. Regarding the choice of port number, TwinCAT2, port number is 801; TwinCAT3, port number is 851, NETID can be input manually or automatically.
For details, please refer to the example of API documentation code

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  [HslCommunication.Core.NetBinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)  
    [HslCommunication.Core.DeviceDeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)  
      [HslCommunication.Core.DeviceDeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)  
        HslCommunication.Profinet.BeckhoffBeckhoffAdsNet

**命名空间：**
 [HslCommunication.Profinet.Beckhoff](fb2af334-3662-b9ca-be3c-b0090bf4ef56.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class BeckhoffAdsNet : DeviceTcpNet
```

```
Public Class BeckhoffAdsNet
	Inherits DeviceTcpNet
```

```
public ref class BeckhoffAdsNet : public DeviceTcpNet
```

```
type BeckhoffAdsNet =  
    class
        inherit DeviceTcpNet
    end
```

BeckhoffAdsNet 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [BeckhoffAdsNet](c6da3e11-3f78-2433-01d0-0d19145efc44.htm) | 实例化一个默认的对象  Instantiate a default object |
| 公共方法 | [BeckhoffAdsNet(String, Int32)](72cd48b9-6409-df6b-b93e-03220e25bd24.htm) | 通过指定的ip地址以及端口号实例化一个默认的对象  Instantiate a default object with the specified IP address and port number |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [AmsPort](69677d41-296b-f1e9-4c91-1317977145c8.htm) | 获取或设置Ams的端口号信息，TwinCAT2，端口号801,811,821,831；TwinCAT3，端口号为851,852,853  Get or set the port number information of Ams, TwinCAT2, the port number is 801, 811, 821, 831; TwinCAT3, the port number is 851, 852, 853 |
| 公共属性代码示例 | [ByteTransform](f35118c9-3691-f752-ae6f-2c5549714faa.htm) | 当前的数据变换机制，当你需要从字节数据转换类型数据的时候需要。  The current data transformation mechanism is required when you need to convert type data from byte data. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共属性 | [CommunicationPipe](ad8e903b-d289-21b0-077e-15ede8005221.htm) | 获取或设置当前的管道信息，管道类型为[CommunicationPipe](ad8e903b-d289-21b0-077e-15ede8005221.htm)的继承类，内置了[PipeTcpNet](26b94bd5-2695-36e7-6378-c3d41846be8f.htm)管道，[PipeUdpNet](5a0b810e-187e-826b-0a13-34fd5576aa0e.htm)管道，[PipeSerialPort](9d924607-3dfc-3916-b532-c5dec48fe347.htm)管道等  Get or set the current pipeline information, the pipeline type is [CommunicationPipe](ad8e903b-d289-21b0-077e-15ede8005221.htm) inheritance class, [PipeTcpNet](26b94bd5-2695-36e7-6378-c3d41846be8f.htm) pipeline, [PipeUdpNet](5a0b810e-187e-826b-0a13-34fd5576aa0e.htm) pipeline, [PipeSerialPort](9d924607-3dfc-3916-b532-c5dec48fe347.htm) pipeline, etc (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [ConnectionId](ce1efb4b-fbc9-4683-163e-0a8e71e99d6b.htm) | 当前连接的唯一ID号，默认为长度20的guid码加随机数组成，方便列表管理，也可以自己指定  The unique ID number of the current connection. The default is a 20-digit guid code plus a random number. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性代码示例 | [ConnectTimeOut](71a651e5-737d-6e88-75a2-fb891b9af6cd.htm) | 获取或设置接收服务器反馈的时间，如果为负数，则不接收反馈   Gets or sets the time to receive server feedback, and if it is a negative number, does not receive feedback (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共属性代码示例 | [IpAddress](dafc8bcd-2bb3-a0a1-6191-dcadd9e01293.htm) | 获取或是设置远程服务器的IP地址，如果是本机测试，那么需要设置为127.0.0.1   Get or set the IP address of the remote server. If it is a local test, then it needs to be set to 127.0.0.1 (重写 [DeviceTcpNetIpAddress](c378f4e8-3fdf-1932-7239-ed44bc5af584.htm).) |
| 公共属性 | [LocalBinding](622b418f-5969-3726-5dcc-04456c4b97a7.htm) | 获取或设置绑定的本地的IP地址和端口号信息，如果端口设置为0，代表任何可用的端口  Get or set the bound local IP address and port number information, if the port is set to 0, it means any available port (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共属性代码示例 | [LogNet](8c56f246-88ac-afb8-8c06-972fafe4f78a.htm) | 组件的日志工具，支持日志记录，只要实例化后，当前网络的基本信息，就以[DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm)等级进行输出  The component's logging tool supports logging. As long as the instantiation of the basic network information, the output will be output at [DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm) (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性代码示例 | [Port](6f1ad24a-c916-5b67-aa17-a8c7dc2acb2e.htm) | 获取或设置服务器的端口号，具体的值需要取决于对方的配置  Gets or sets the port number of the server. The specific value depends on the configuration of the other party. (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共属性代码示例 | [ReceiveTimeOut](250560a7-332b-ce1b-5eb0-e3abe4ee6012.htm) | 获取或设置接收服务器反馈的时间，如果为负数，则不接收反馈   Gets or sets the time to receive server feedback, and if it is a negative number, does not receive feedback (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [SendBeforeHex](72aa1de4-116f-abfe-4adb-814a0be0e560.htm) | 获取或设置在发送通信报文前追加发送的字节信息，HEX格式，通常用于lora组网时，需要携带 00 00 00 02 四个字节的站地址功能。  Obtain or set the byte information sent before sending communication packets, HEX format, usually used for LORA networking, you need to carry 00 00 00 02 four-byte station address function. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [SleepTime](86a5fbb9-3456-6beb-7113-5bbad5eabbc5.htm) | 获取或设置在正式接收对方返回数据前的时候，需要休息的时间，当设置为0的时候，不需要休息。  Get or set the time required to rest before officially receiving the data from the other party. When it is set to 0, no rest is required. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [SocketKeepAliveTime](96ab7703-8c4d-b4c2-ddd9-0dc93ebcfc8d.htm) | 获取或设置客户端的Socket的心跳时间信息，这个是Socket底层自动实现的心跳包，不基于协议层实现。默认小于0，不开启心跳检测，如果需要开启，设置 60\_000 比较合适，单位毫秒  Get or set the heartbeat time information of the Socket of the client. This is the heartbeat packet automatically implemented by the bottom layer of the Socket, not based on the protocol layer. The default value is less than 0, and heartbeat detection is not enabled. If you need to enable it, it is more appropriate to set 60\_000, in milliseconds. (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共属性 | [UseAutoAmsNetID](d4280c64-cf22-b428-bc3f-3ec08522754c.htm) | 是否使用服务器自动的NETID信息，默认手动设置  Whether to use the server's automatic NETID information, manually set by default |
| 公共属性 | [UseTagCache](94a53279-b7cc-a5f5-ce4c-93e224e74b29.htm) | 是否使用标签的名称缓存功能，默认为 False  Whether to use tag name caching. The default is False |
| 受保护的属性 | [WordLength](f1a505cc-e85c-6c8b-b189-58a1ab5753f9.htm) | 一个字单位的数据表示的地址长度，西门子为2，三菱，欧姆龙，modbusTcp就为1，AB PLC无效  The address length represented by one word of data, Siemens is 2, Mitsubishi, Omron, modbusTcp is 1, AB PLC is invalid (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [ClearTagsCache](3ca05c73-2ab5-b7b8-2d46-03da8a9c67ad.htm) | 清除标签缓存信息  Clear the tag cache information |
| 公共方法代码示例 | [ConnectClose](569a9c0e-7628-b8d5-e124-0853452747b8.htm) | 手动断开与远程服务器的连接，如果当前是长连接模式，那么就会切换到短连接模式  Manually disconnect from the remote server, if it is currently in long connection mode, it will switch to short connection mode (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共方法代码示例 | [ConnectCloseAsync](76be8150-ff22-95b9-16b1-7c3440c54ae8.htm) | 手动断开与远程服务器的连接，如果当前是长连接模式，那么就会切换到短连接模式  Manually disconnect from the remote server, if it is currently in long connection mode, it will switch to short connection mode (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共方法代码示例 | [ConnectServer](6299cc1f-2fab-b754-4597-1eddc2b21308.htm) | 尝试连接远程的服务器，如果连接成功，就切换短连接模式到长连接模式，后面的每次请求都共享一个通道，使得通讯速度更快速  Try to connect to a remote server. If the connection is successful, switch the short connection mode to the long connection mode. Each subsequent request will share a channel, making the communication speed faster. (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共方法代码示例 | [ConnectServerAsync](e7eff5c0-8e2a-3b6c-5027-16f4d8bcf5ae.htm) | 尝试连接远程的服务器，如果连接成功，就切换短连接模式到长连接模式，后面的每次请求都共享一个通道，使得通讯速度更快速  Try to connect to a remote server. If the connection is successful, switch the short connection mode to the long connection mode. Each subsequent request will share a channel, making the communication speed faster. (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 受保护的方法 | [DecideWhetherQAMessage](771586fe-9184-5312-e3da-bae61fe91d9a.htm) | 决定当前的消息是否是用于问答机制返回的消息，默认直接返回 true, 实际的情况需要根据协议进行重写方法  To determine whether the current message is the message returned by the question answering mechanism, the default is true. In actual cases, the rewriting method needs to be performed according to the protocol (重写 [BinaryCommunicationDecideWhetherQAMessage(CommunicationPipe, OperateResultByte)](acc40cda-e310-8537-6685-59b67eeb16ff.htm).) |
| 公共方法 | [Dispose](b9dda6bf-d342-254e-689a-93fa0bb265d1.htm) | (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 受保护的方法 | [Dispose(Boolean)](0bbb6ed3-111e-069c-e9f5-8cbc3f7f7b77.htm) | (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | [ExtraAfterReadFromCoreServer](3d26d83d-9287-54f5-3d35-c6d074ce6263.htm) | 和服务器交互完成的时候调用的方法，可以根据读写结果进行一些额外的操作，具体的操作需要根据实际的需求来重写实现  The method called when the interaction with the server is completed can perform some additional operations based on the read and write results. The specific operations need to be rewritten according to actual needs. (重写 [BinaryCommunicationExtraAfterReadFromCoreServer(OperateResult)](06a49a15-a3ce-5e39-28c7-2b78851eca2e.htm).) |
| 受保护的方法代码示例 | [ExtraOnDisconnect](de9dbba2-3b2c-1c5e-53c8-4e6f8a44e2c4.htm) | 根据实际的协议选择是否重写本方法，有些协议在断开连接之前，需要发送一些报文来关闭当前的网络通道  Select whether to rewrite this method according to the actual protocol. Some protocols need to send some packets to close the current network channel before disconnecting. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法代码示例 | [ExtraOnDisconnectAsync](210f593e-16d8-77c4-1178-58fbcfe9b27b.htm) | 根据实际的协议选择是否重写本方法，有些协议在断开连接之前，需要发送一些报文来关闭当前的网络通道  Select whether to rewrite this method according to the actual protocol. Some protocols need to send some packets to close the current network channel before disconnecting. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 受保护的方法 | [GetLogTextFromBinary](67f3b2f7-957c-fa0b-8320-6799237157db.htm) | 获取当前的报文进行日志记录的时候，是否使用二进制的格式记录，默认返回 [LogMsgFormatBinary](2ed2ff87-caec-0e44-f08b-e2b015225a58.htm)，重写可以根据session对象分别返回不同记录模式  Whether to log the current packet in binary format, the default return is [LogMsgFormatBinary](2ed2ff87-caec-0e44-f08b-e2b015225a58.htm). If you want to override it, different recording modes can be returned according to session (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [GetNewNetMessage](215ca52b-b151-9cb6-0443-00b9be103188.htm) | 获取一个新的消息对象的方法，需要在继承类里面进行重写  The method to get a new message object needs to be overridden in the inheritance class (重写 [BinaryCommunicationGetNewNetMessage](16636a1b-a6f0-9c0f-9d9c-c94e26c67459.htm).) |
| 公共方法 | [GetSenderAMSNetId](f743d65e-7896-7d75-ba57-6be4991ab9af.htm) | 获取当前发送的AMS的网络ID信息 |
| 公共方法 | [GetTargetAMSNetId](dcc58722-6399-8dd4-affa-2e531a098ae0.htm) | 获取当前目标的AMS网络的ID信息 |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法 | [GetWordLength](c27821ea-1739-8eb9-a848-23a00828fec8.htm) | 一个字单位的数据表示的地址长度，西门子为2，三菱，欧姆龙，modbusTcp就为1，AB PLC无效  The address length represented by one word of data, Siemens is 2, Mitsubishi, Omron, modbusTcp is 1, AB PLC is invalid (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 受保护的方法代码示例 | [InitializationOnConnect](ea40c0b0-f579-3a64-7b45-573259ef706d.htm) | 根据实际的协议选择是否重写本方法，有些协议在创建连接之后，需要进行一些初始化的信号握手，才能最终建立网络通道。  Whether to rewrite this method is based on the actual protocol. Some protocols require some initial signal handshake to establish a network channel after the connection is created. (重写 [BinaryCommunicationInitializationOnConnect](4ab6989f-a5b5-05fe-8b44-bfa45824a2a7.htm).) |
| 受保护的方法代码示例 | [InitializationOnConnectAsync](c2a3f0f3-3a78-e102-d07f-d15053af2633.htm) | 根据实际的协议选择是否重写本方法，有些协议在创建连接之后，需要进行一些初始化的信号握手，才能最终建立网络通道。  Whether to rewrite this method is based on the actual protocol. Some protocols require some initial signal handshake to establish a network channel after the connection is created. (重写 [BinaryCommunicationInitializationOnConnectAsync](7da2af04-d4fc-5045-19a4-4b0d1722eef3.htm).) |
| 公共方法 | [IpAddressPing](0ac126b4-a830-a822-4ef5-855474942cd7.htm) | 对当前设备的IP地址进行PING的操作，返回PING的结果，正常来说，返回Success  PING the IP address of the current device and return the PING result. Normally, it returns Success (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 受保护的方法 | [LogRevcMessage(Byte)](dab1657f-1e75-7fe1-7003-ccedfc69855d.htm) | 使用日志记录一个接收的报文信息  Logs are used to record information about a received packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [LogRevcMessage(Byte, PipeSession)](0f51ae45-575e-4228-f978-01b37d6d5df1.htm) | 使用日志记录一个接收的报文信息  Logs are used to record information about a received packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [LogSendMessage(Byte)](4f90eecc-68a2-1ef2-44ef-369f19bed3ec.htm) | 使用日志记录一个发送的报文信息  Logs are used to record information about a send packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [LogSendMessage(Byte, PipeSession)](209caadf-4f0f-40aa-04c2-ebc61f8dab13.htm) | 使用日志记录一个发送的报文信息  Logs are used to record information about a send packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法 | [PackCommandWithHeader](3895f6bc-bf81-5dbe-cb27-c12c6602126d.htm) | 对当前的命令进行打包处理，通常是携带命令头内容，标记当前的命令的长度信息，需要进行重写，否则默认不打包  The current command is packaged, usually carrying the content of the command header, marking the length of the current command, and it needs to be rewritten, otherwise it is not packaged by default (重写 [BinaryCommunicationPackCommandWithHeader(Byte)](d0b2cd76-c306-11fc-38c7-af8ad306cce7.htm).) |
| 公共方法 | [Read(String, UInt16)](ac58c335-e208-6bf5-6d30-85dc57694090.htm) | 读取PLC的数据，地址共有三种格式，一：I,Q,M数据信息，举例M0,M100；二：内存地址，i=100000；三：标签地址，s=A  Read PLC data, there are three formats of address, one: I, Q, M data information, such as M0, M100; two: memory address, i = 100000; three: tag address, s = A (重写 [DeviceCommunicationRead(String, UInt16)](9dd09452-c200-bc65-75ff-69aa34b366d4.htm).) |
| 公共方法 | [Read(String, UInt16)](24dae660-f973-3403-a4a8-ed39d114483b.htm) | 批量读取PLC的数据，需要传入地址数组，以及读取的长度数组信息，长度单位为字节单位，如果是读取bool变量的，则以bool为单位，统一返回一串字节数据信息，需要进行二次解析的操作。  To read PLC data in batches, you need to pass in the address array and the read length array information. The unit of length is in bytes. If you read a bool variable, it will return a string of byte data information in units of bool. , which requires a secondary parsing operation. |
| 公共方法代码示例 | [ReadT](6af520c0-a616-6709-ea80-14e26306e2be.htm) | 读取支持Hsl特性的数据内容，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明，或是官网的详细文档  Read the data content of the Hsl attribute. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (重写 [DeviceCommunicationReadT](95abe064-5594-5ae7-a849-98c53f88fee0.htm).) |
| 公共方法 | [ReadAdsDeviceInfo](1b58e4ea-7a7c-3b42-1f11-b077010b9a39.htm) | 读取Ads设备的设备信息。主要是版本号，设备名称  Read the device information of the Ads device. Mainly version number, device name |
| 公共方法 | [ReadAdsDeviceInfoAsync](18ef1d19-2639-b4b9-99c6-3e7caed8ef5a.htm) | 读取Ads设备的设备信息。主要是版本号，设备名称  Read the device information of the Ads device. Mainly version number, device name |
| 公共方法 | [ReadAdsState](4645af4b-15b0-49a7-a3d3-bc07ef913ad7.htm) | 读取Ads设备的状态信息，其中[Content1](fba87d7e-fbf8-0667-b0f2-7f669d3179fe.htm)是Ads State，[Content2](312b5df6-df3c-07d6-83c2-391380a3fb33.htm)是Device State  Read the status information of the Ads device, where [Content1](fba87d7e-fbf8-0667-b0f2-7f669d3179fe.htm) is the Ads State, and [Content2](312b5df6-df3c-07d6-83c2-391380a3fb33.htm) is the Device State |
| 公共方法 | [ReadAdsStateAsync](831eac3f-b313-b2bd-3131-71bf3948b652.htm) | 读取Ads设备的状态信息，其中[Content1](fba87d7e-fbf8-0667-b0f2-7f669d3179fe.htm)是Ads State，[Content2](312b5df6-df3c-07d6-83c2-391380a3fb33.htm)是Device State  Read the status information of the Ads device, where [Content1](fba87d7e-fbf8-0667-b0f2-7f669d3179fe.htm) is the Ads State, and [Content2](312b5df6-df3c-07d6-83c2-391380a3fb33.htm) is the Device State |
| 公共方法 | [ReadAsync(String, UInt16)](82d4ecb0-9450-53ba-101a-4252209345f5.htm) | 读取PLC的数据，地址共有三种格式，一：I,Q,M数据信息，举例M0,M100；二：内存地址，i=100000；三：标签地址，s=A  Read PLC data, there are three formats of address, one: I, Q, M data information, such as M0, M100; two: memory address, i = 100000; three: tag address, s = A (重写 [DeviceCommunicationReadAsync(String, UInt16)](d8c14d37-56dc-2a31-1a19-3e2b5b65197d.htm).) |
| 公共方法 | [ReadAsync(String, UInt16)](2b8c8a10-8a40-2e2b-6917-ff16cc6a0c2d.htm) | 批量读取PLC的数据，需要传入地址数组，以及读取的长度数组信息，长度单位为字节单位，如果是读取bool变量的，则以bool为单位，统一返回一串字节数据信息，需要进行二次解析的操作。  To read PLC data in batches, you need to pass in the address array and the read length array information. The unit of length is in bytes. If you read a bool variable, it will return a string of byte data information in units of bool. , which requires a secondary parsing operation. |
| 公共方法代码示例 | [ReadAsyncT](c5887b82-93a5-0799-23a9-be999466b175.htm) | 异步读取支持Hsl特性的数据内容，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Asynchronously read the data content of the Hsl attribute. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (重写 [DeviceCommunicationReadAsyncT](0522e740-4bf0-2e31-c4c5-e41384321b37.htm).) |
| 公共方法 | [ReadBool(String)](9b2faa2e-eaa9-7ace-c49f-7a0318f1df1a.htm) | 读取单个的Boolean数据信息  Read a single Boolean data message (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadBool(String, UInt16)](92020f06-15c1-4b5f-ce8d-36a4d6182a9c.htm) | 读取PLC的数据，地址共有三种格式，一：I,Q,M数据信息，举例M0,M100；二：内存地址，i=100000；三：标签地址，s=A  Read PLC data, there are three formats of address, one: I, Q, M data information, such as M0, M100; two: memory address, i = 100000; three: tag address, s = A (重写 [DeviceCommunicationReadBool(String, UInt16)](69d0d24f-34a7-5b41-5c1b-619edc7e6b54.htm).) |
| 公共方法 | [ReadBoolAsync(String)](cc37da1f-d673-fe67-78ee-4e465d29a2bd.htm) | 异步读取单个的Boolean数据信息  Asynchronously read a single Boolean data message (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadBoolAsync(String, UInt16)](d43cef4a-36ed-4bd2-b999-ecf789b95a28.htm) | 读取PLC的数据，地址共有三种格式，一：I,Q,M数据信息，举例M0,M100；二：内存地址，i=100000；三：标签地址，s=A  Read PLC data, there are three formats of address, one: I, Q, M data information, such as M0, M100; two: memory address, i = 100000; three: tag address, s = A (重写 [DeviceCommunicationReadBoolAsync(String, UInt16)](a45d383d-0095-3ea7-32f8-e75af88d60f8.htm).) |
| 公共方法 | [ReadByte](dadb0178-8d9c-2ce1-5495-e53d72771b82.htm) | 读取PLC的数据，地址共有三种格式，一：I,Q,M数据信息，举例M0,M100；二：内存地址，i=100000；三：标签地址，s=A  Read PLC data, there are three formats of address, one: I, Q, M data information, such as M0, M100; two: memory address, i = 100000; three: tag address, s = A |
| 公共方法 | [ReadByteAsync](d61db681-dca5-f888-2dfd-c35faa269a4f.htm) | 读取PLC的数据，地址共有三种格式，一：I,Q,M数据信息，举例M0,M100；二：内存地址，i=100000；三：标签地址，s=A  Read PLC data, there are three formats of address, one: I, Q, M data information, such as M0, M100; two: memory address, i = 100000; three: tag address, s = A |
| 公共方法代码示例 | [ReadCustomerT(String)](ca069688-d8f0-7dbd-e962-86a10dce83cf.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，返回一个新的类型的实例对象。  To read a custom data type, you need to inherit from the IDataTransfer interface and return an instance object of a new type. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadCustomerT(String, T)](ca7dc1a1-c97d-689f-9819-4dae4af67206.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，传入一个实例，对这个实例进行赋值，并返回该实例的对象。  To read a custom data type, you need to inherit from the IDataTransfer interface, pass in an instance, assign a value to this instance, and return the object of the instance. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadCustomerAsyncT(String)](f3fcda2b-255d-aadd-567e-c7b3df3636d9.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，返回一个新的类型的实例对象。  To read a custom data type, you need to inherit from the IDataTransfer interface and return an instance object of a new type. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadCustomerAsyncT(String, T)](527c9fae-4743-6ebd-7b9a-5ef3a8b7c455.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，传入一个实例，对这个实例进行赋值，并返回该实例的对象。  To read a custom data type, you need to inherit from the IDataTransfer interface, pass in an instance, assign a value to this instance, and return the object of the instance. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadDouble(String)](e7e31c43-71b4-e692-6bdc-ba9533e1d83d.htm) | 读取双浮点的数据  Read double floating point data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadDouble(String, UInt16)](4eee0f05-f861-4b48-3f06-99737fb24aee.htm) | 读取双浮点数据的数组  Read double floating point data array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadDoubleAsync(String)](2cf895fa-1033-1def-f048-0e8caefaa71d.htm) | 异步读取双浮点的数据  Asynchronously read double floating point data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadDoubleAsync(String, UInt16)](2256d406-d4d5-da54-210b-ccef58a5e914.htm) | 异步读取双浮点数据的数组  Asynchronously read double floating point data array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadFloat(String)](bc2801ae-6600-e845-7c6c-522c9156eb0e.htm) | 读取单浮点数据  Read single floating point data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadFloat(String, UInt16)](6a9fb098-b13c-313b-3b27-20d3c787945d.htm) | 读取单浮点精度的数组  Read single floating point array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadFloatAsync(String)](89f743f8-4b3a-7085-1e49-a1d7adb10052.htm) | 异步读取单浮点数据  Asynchronously read single floating point data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadFloatAsync(String, UInt16)](cb6d5c5d-9b5c-2467-3e6b-5efa9d20cd72.htm) | 异步读取单浮点精度的数组  Asynchronously read single floating point array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServer(Byte)](8146dc0c-1edc-3ca7-428a-a3bdbd35e6c5.htm) | 将二进制的数据发送到管道中去，然后从管道里接收二进制的数据回来，并返回是否成功的结果对象。  Send binary data to the pipeline, and then receive binary data back from the pipeline, and return whether the success of the result object (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServer(IEnumerableByte)](044cfdd9-ec76-b89f-2f2b-5f439025006c.htm) | 将多个数据报文按顺序发到设备，并从设备接收返回的数据内容，然后拼接成一个Byte[]信息，需要重写[UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm)方法才能返回正确的结果。  Send multiple data packets to the device in sequence, and receive the returned data content from the device, and then splicing them into a Byte[] message, you need to rewrite [UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm) method to return the correct result. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServer(Byte, Boolean, Boolean)](b32ce8e8-a959-fcc2-6736-14b39fdb9a40.htm) | 将二进制的数据发送到管道中去，然后从管道里接收二进制的数据回来，并返回是否成功的结果对象。  Send binary data to the pipeline, and then receive binary data back from the pipeline, and return whether the success of the result object (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServer(CommunicationPipe, Byte, Boolean, Boolean)](ef36d56b-df09-f727-e8d3-316d37a13ee9.htm) | 使用指定的管道来进行数据通信，发送原始数据到管道，然后从管道接收相关的数据返回，本方法无锁 (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServerAsync(Byte)](a4de5a24-c67a-7b02-db05-33b85eaf6e1e.htm) | 将二进制的数据发送到管道中去，然后从管道里接收二进制的数据回来，并返回是否成功的结果对象。  Send binary data to the pipeline, and then receive binary data back from the pipeline, and return whether the success of the result object (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServerAsync(IEnumerableByte)](bcdb59b7-a70a-d5c0-1fc6-9ececcd39ea2.htm) | 将多个数据报文按顺序发到设备，并从设备接收返回的数据内容，然后拼接成一个Byte[]信息，需要重写[UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm)方法才能返回正确的结果。  Send multiple data packets to the device in sequence, and receive the returned data content from the device, and then splicing them into a Byte[] message, you need to rewrite [UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm) method to return the correct result. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServerAsync(Byte, Boolean, Boolean)](1c79c204-5ab5-9342-e65d-7455c20ea577.htm) | 将二进制的数据发送到管道中去，然后从管道里接收二进制的数据回来，并返回是否成功的结果对象。  Send binary data to the pipeline, and then receive binary data back from the pipeline, and return whether the success of the result object (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServerAsync(CommunicationPipe, Byte, Boolean, Boolean)](a254a19d-3c29-147b-335b-8adccdf4ee2f.htm) | 使用指定的管道来进行数据通信，发送原始数据到管道，然后从管道接收相关的数据返回，本方法无锁 (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法代码示例 | [ReadInt16(String)](57a16003-2fd5-b2bd-09bf-14366be899b8.htm) | 读取16位的有符号的整型数据  Read 16-bit signed integer data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt16(String, UInt16)](b2bde158-f51c-93fc-66ca-5e4cf6ea1fb7.htm) | 读取16位的有符号整型数组  Read 16-bit signed integer array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt16Async(String)](ac021793-890d-ceab-eae4-47ff763390c1.htm) | 异步读取16位的有符号的整型数据  Asynchronously read 16-bit signed integer data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt16Async(String, UInt16)](f30e02fc-8078-f1b3-c226-5991e3e23caf.htm) | 异步读取16位的有符号整型数组  Asynchronously read 16-bit signed integer array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt32(String)](bfbea035-9a83-40e3-084c-28122041b37a.htm) | 读取32位的有符号整型  Read 32-bit signed integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt32(String, UInt16)](095240a3-bc0f-8647-4867-88e40cc089e3.htm) | 读取32位有符号整型数组  Read 32-bit signed integer array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt32Async(String)](96726ea0-82aa-1229-21f6-43cd62abca43.htm) | 异步读取32位的有符号整型  Asynchronously read 32-bit signed integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt32Async(String, UInt16)](d6356095-a289-18eb-136c-b7e6ae099616.htm) | 异步读取32位有符号整型数组  Asynchronously read 32-bit signed integer array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt64(String)](0d0c2284-2da1-6f9a-b353-537561cd5e95.htm) | 读取64位的有符号整型  Read 64-bit signed integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt64(String, UInt16)](35a55019-efb1-bf2d-be1b-01c8eed99243.htm) | 读取64位的有符号整型数组  Read 64-bit signed integer array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt64Async(String)](26d46c8a-3775-3abd-d3ec-5cd1758a317e.htm) | 异步读取64位的有符号整型  Asynchronously read 64-bit signed integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt64Async(String, UInt16)](8dc97c4a-8904-9a34-9ddc-76c1a9477ee5.htm) | 异步读取64位的有符号整型数组  Asynchronously read 64-bit signed integer array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadString(String, UInt16)](07812d76-415a-6400-c8a2-1c21da51e7a3.htm) | 读取字符串数据，默认为最常见的ASCII编码  Read string data, default is the most common ASCII encoding (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadString(String, UInt16, Encoding)](b34c71a6-58f3-6bf4-0b93-471be4f2ed31.htm) | 使用指定的编码，读取字符串数据  Reads string data using the specified encoding (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadStringAsync(String, UInt16)](1835da81-b811-d151-0b53-907d8d87d0a8.htm) | 异步读取字符串数据，默认为最常见的ASCII编码  Asynchronously read string data, default is the most common ASCII encoding (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadStringAsync(String, UInt16, Encoding)](091dfc67-3d25-f33e-3a55-a0228863815e.htm) | 异步使用指定的编码，读取字符串数据  Asynchronously reads string data using the specified encoding (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadStructT(String)](b3f311a4-68d8-87b4-a972-ad8ec4b3f0df.htm) | 读取结构体的信息，传入结构体的类型，以及结构体的起始地址  Read the information of the structure, the type of the incoming structure, and the start address of the structure |
| 公共方法代码示例 | [ReadStructT(String, UInt16)](49853082-fbad-52e5-d756-0615cedb4b83.htm) | 读取结构体类型的数据，根据结构体自身的定义，读取原始字节数组，然后解析出实际的结构体数据，结构体需要实现[HslStructAttribute](87e427a8-0c7b-43a8-c421-c27264eda90f.htm)特性 (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadStructAsyncT](36461e61-d0ed-62cf-79ce-dcc7f71a3fa2.htm) | 读取结构体类型的数据，根据结构体自身的定义，读取原始字节数组，然后解析出实际的结构体数据，结构体需要实现[HslStructAttribute](87e427a8-0c7b-43a8-c421-c27264eda90f.htm)特性 (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt16(String)](736379d1-9560-5052-658a-13c126d56d78.htm) | 读取16位的无符号整型  Read 16-bit unsigned integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt16(String, UInt16)](e284c4bd-d3f8-ff30-eb42-0aa035cfb4f6.htm) | 读取16位的无符号整型数组  Read 16-bit unsigned integer array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt16Async(String)](361acc93-43b5-a9a5-6fa7-fa8728b7b96c.htm) | 异步读取16位的无符号整型  Asynchronously read 16-bit unsigned integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt16Async(String, UInt16)](73dd24b8-b501-df53-1de6-ccb2708f656d.htm) | 异步读取16位的无符号整型数组  Asynchronously read 16-bit unsigned integer array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt32(String)](ca73c958-39c7-9b34-aa44-3db2d9b0cad3.htm) | 读取32位的无符号整型  Read 32-bit unsigned integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt32(String, UInt16)](878b4709-0278-fcdf-3b65-afbf0f496c1e.htm) | 读取32位的无符号整型数组  Read 32-bit unsigned integer array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt32Async(String)](2343953f-6350-331e-6ba7-18a4449f0575.htm) | 异步读取32位的无符号整型  Asynchronously read 32-bit unsigned integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt32Async(String, UInt16)](e0eba6c1-56a2-5190-ba03-66116bcafe3d.htm) | 异步读取32位的无符号整型数组  Asynchronously read 32-bit unsigned integer array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt64(String)](947487dd-d4a3-72d8-d9fb-2224aad70e12.htm) | 读取64位的无符号整型  Read 64-bit unsigned integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt64(String, UInt16)](858766a9-9d62-87a8-f727-a00838f3d2b6.htm) | 读取64位的无符号整型的数组  Read 64-bit unsigned integer array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt64Async(String)](3411ff3a-772f-dbe1-8931-75e00a731579.htm) | 异步读取64位的无符号整型  Asynchronously read 64-bit unsigned integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt64Async(String, UInt16)](c2095a6e-ff89-dd98-8a21-fe72f82dbbc6.htm) | 异步读取64位的无符号整型的数组  Asynchronously read 64-bit unsigned integer array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadValueHandle](7a0c126f-5f3e-f077-f90b-34ca2191c5b2.htm) | 根据当前标签的地址获取到内存偏移地址  Get the memory offset address based on the address of the current label |
| 公共方法 | [ReadValueHandleAsync](0dadc9a4-90cc-1f29-627f-f8112aaa7df2.htm) | 根据当前标签的地址获取到内存偏移地址  Get the memory offset address based on the address of the current label |
| 公共方法 | [ReleaseSystemHandle](e2fa36f0-6d1f-2e02-e3de-bc2344e8d553.htm) | 释放当前的系统句柄，该句柄是通过[ReadValueHandle(String)](7a0c126f-5f3e-f077-f90b-34ca2191c5b2.htm)获取的 |
| 公共方法 | [ReleaseSystemHandleAsync](cb4896ea-bbd4-52fa-4587-db7c137ff2d6.htm) | 释放当前的系统句柄，该句柄是通过[ReadValueHandle(String)](7a0c126f-5f3e-f077-f90b-34ca2191c5b2.htm)获取的 |
| 公共方法 | [SetDtuPipe](f765d1d7-234a-b229-9457-97e25e024dc9.htm) | 将当前的通信对象设置DTU模式，允许传入现成的管道，并返回初始化结果，如果该设备重写了握手报文，就是返回握手结果  Set the current communication object to DTU mode, allow the existing pipe to be passed in, and return the initialization result, if the device rewrites the handshake packet, the handshake result is returned (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [SetDtuPipeAsync](9d1fd12d-b804-3828-a627-a2c0eef1d23f.htm) | 将当前的通信对象设置DTU模式，允许传入现成的管道，并返回初始化结果，如果该设备重写了握手报文，就是返回握手结果  Set the current communication object to DTU mode, allow the existing pipe to be passed in, and return the initialization result, if the device rewrites the handshake packet, the handshake result is returned (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [SetPersistentConnection](f903dae0-4a4e-c7db-8350-50e36b257ec5.htm) | **已过时。** V11版本及之前设置长连接的方法，在V12版本以上中没有任何效果，默认长连接，删除调用即可，此处保留方法是为了部分用户保持兼容性升级。  The method of setting the long connection in V11 and before, has no effect in V12 and above. this method can be deleted. The method is retained here to maintain compatibility upgrades for some users. (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共方法 | [SetSenderAMSNetId](6c085271-688a-909e-53c9-dec86b2a8e8a.htm) | 设置原目标地址 举例 192.168.0.100.1.1；也可以是带端口号 192.168.0.100.1.1:34567  Set the original destination address Example: 192.168.0.100.1.1; it can also be the port number 192.168.0.100.1.1: 34567 |
| 公共方法 | [SetTargetAMSNetId](af74dcd7-fb14-527d-5b79-881112e412ad.htm) | 目标的地址，举例 192.168.0.1.1.1；也可以是带端口号 192.168.0.1.1.1:801  The address of the destination, for example 192.168.0.1.1.1; it can also be the port number 192.168.0.1.1.1: 801 |
| 公共方法 | [ToString](5cf941e3-75eb-cfcb-ec45-0065be0058c5.htm) | (重写 [DeviceTcpNetToString](209a196b-90ea-2b73-e915-ce4b11f1263d.htm).) |
| 公共方法 | [TransValueHandle](f7d19175-e124-9a8d-d36e-638b14d9c919.htm) | 将字符串的地址转换为内存的地址，其他地址则不操作  Converts the address of a string to the address of a memory, other addresses do not operate |
| 公共方法 | [TransValueHandleAsync](efa735b1-7770-c061-0e6f-0486b84ff6b8.htm) | 将字符串的地址转换为内存的地址，其他地址则不操作  Converts the address of a string to the address of a memory, other addresses do not operate |
| 公共方法 | [UnpackResponseContent](9578724a-11cf-6cec-c183-18fb69db6412.htm) | 根据对方返回的报文命令，对命令进行基本的拆包，例如各种Modbus协议拆包为统一的核心报文，还支持对报文的验证  According to the message command returned by the other party, the command is basically unpacked, for example, various Modbus protocols are unpacked into a unified core message, and the verification of the message is also supported (重写 [BinaryCommunicationUnpackResponseContent(Byte, Byte)](208017a7-8da2-33f8-3273-d7e7e850100f.htm).) |
| 公共方法 | [Wait(String, Boolean, Int32, Int32)](0e4db26e-1177-2b18-9492-ba05afb6fcc4.htm) | 等待指定地址的Boolean值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the Boolean value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [Wait(String, Int16, Int32, Int32)](e06fe098-8c8b-83f9-73ca-2a7cce19d43d.htm) | 等待指定地址的Int16值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the Int16 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [Wait(String, Int32, Int32, Int32)](1bce0aae-27e9-f1d6-e8f3-25ae6b47554e.htm) | 等待指定地址的Int32值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the Int32 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [Wait(String, Int64, Int32, Int32)](073354e9-2508-2158-1b01-881c3a34844b.htm) | 等待指定地址的Int64值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the Int64 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [Wait(String, UInt16, Int32, Int32)](d97f8d26-1c1a-df69-93f0-f60209597ced.htm) | 等待指定地址的UInt16值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the UInt16 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [Wait(String, UInt32, Int32, Int32)](1e61c395-b07f-e66c-5f9c-2708093150d9.htm) | 等待指定地址的UInt32值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the UInt32 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [Wait(String, UInt64, Int32, Int32)](d659729b-5766-c497-d418-856c34ec2b24.htm) | 等待指定地址的UInt64值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the UInt64 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WaitAsync(String, Boolean, Int32, Int32)](65399e87-9d38-f819-ccd3-005f68c2c3b3.htm) | 等待指定地址的Boolean值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the Boolean value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WaitAsync(String, Int16, Int32, Int32)](01277e39-afb1-86a4-3aa3-560160e66fb9.htm) | 等待指定地址的Int16值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the Int16 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WaitAsync(String, Int32, Int32, Int32)](4a6d380f-6b2a-1270-07a7-39b88a2712fd.htm) | 等待指定地址的Int32值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the Int32 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WaitAsync(String, Int64, Int32, Int32)](d7eede9b-4d44-479b-fb88-3ff2e1bb7568.htm) | 等待指定地址的Int64值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the Int64 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WaitAsync(String, UInt16, Int32, Int32)](0a8395c7-dc7f-4d5f-4316-98e18fe70f4c.htm) | 等待指定地址的UInt16值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the UInt16 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WaitAsync(String, UInt32, Int32, Int32)](19e9074e-99a2-99fd-91ef-a5314a115397.htm) | 等待指定地址的UInt32值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the UInt32 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WaitAsync(String, UInt64, Int32, Int32)](5109dac0-0ba2-cbb9-e383-1a869c6dc22e.htm) | 等待指定地址的UInt64值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the UInt64 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [Write(String, Boolean)](06d0b9af-7f28-0b9c-098f-4da28e69504a.htm) | 写入单个的Boolean数据，返回是否成功  Write a single Boolean data, and return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Double)](64c44576-7e95-cbcb-a8c9-e6b7ff7e3c50.htm) | 写入double数据，返回是否成功  Write double data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Double)](df56e6b1-0688-ece7-3bf1-9db32fec9f6e.htm) | 写入double数组，返回是否成功  Write double array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Int16)](0e23506e-b604-e67c-b5eb-06271b6c6335.htm) | 写入short数据，返回是否成功  Write short data, returns whether success (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Int16)](c7cfd67d-ca31-3087-c684-ee4c10780ab1.htm) | 写入short数组，返回是否成功  Write short array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Int32)](ba80d1a5-c8d8-7e0b-5c6b-9435cf20cadb.htm) | 写入int数据，返回是否成功  Write int data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Int32)](57dac21c-8b28-de7c-c213-77cda465cccb.htm) | 写入int[]数组，返回是否成功  Write int array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Int64)](ba246b33-5a12-a22c-6d40-0969a64b508b.htm) | 写入long数据，返回是否成功  Write long data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Int64)](5b910f9c-6e6b-bb9d-6842-1b496640c1f4.htm) | 写入long数组，返回是否成功  Write long array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Single)](a2265fa7-5b56-a5aa-9d21-a0dc5458cbf9.htm) | 写入float数据，返回是否成功  Write float data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Single)](bb76227e-3e95-8ca2-5dcf-c087d90b8a01.htm) | 写入float数组，返回是否成功  Write float array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, String)](4e0ec363-ba64-7ca7-b8d3-2e9a8087bd55.htm) | 写入字符串信息，编码为ASCII  Write string information, encoded as ASCII (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, UInt16)](e5bc9466-65b4-fc64-e0d0-1867b1113130.htm) | 写入ushort数据，返回是否成功  Write ushort data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, UInt16)](65544962-82ac-bb47-d72f-5dd0c3d4ac6c.htm) | 写入ushort数组，返回是否成功  Write ushort array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, UInt32)](c99d8246-f95e-4f79-738b-16e622072eb4.htm) | 写入uint数据，返回是否成功  Write uint data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, UInt32)](b80b3099-b34d-fde8-6a65-837da3bdc5ab.htm) | 写入uint[]数组，返回是否成功  Write uint array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, UInt64)](4924086e-9b36-d446-4ae4-6682134a6ad1.htm) | 写入ulong数据，返回是否成功  Write ulong data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, UInt64)](b1131bdd-8b51-ec6d-dd0e-fec1f6cb51b6.htm) | 写入ulong数组，返回是否成功  Write ulong array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [Write(String, Boolean)](0a254c9a-8a5a-3bea-e82f-4ba1e13ceb9a.htm) | 写入PLC的数据，地址共有三种格式，一：I,Q,M数据信息，举例M0,M100；二：内存地址，i=100000；三：标签地址，s=A  There are three formats for the data written into the PLC. One: I, Q, M data information, such as M0, M100; two: memory address, i = 100000; three: tag address, s = A (重写 [DeviceCommunicationWrite(String, Boolean)](5b046c82-3d8e-550b-fdf3-41ca5558169b.htm).) |
| 公共方法 | [Write(String, Byte)](40483d77-6945-118f-6b87-b018f6e6bc3c.htm) | 写入PLC的数据，地址共有三种格式，一：I,Q,M数据信息，举例M0,M100；二：内存地址，i=100000；三：标签地址，s=A  There are three formats for the data written into the PLC. One: I, Q, M data information, such as M0, M100; two: memory address, i = 100000; three: tag address, s = A |
| 公共方法 | [Write(String, Byte)](39cd682e-4c5f-2ff7-acb2-c4771ab24081.htm) | 写入PLC的数据，地址共有三种格式，一：I,Q,M数据信息，举例M0,M100；二：内存地址，i=100000；三：标签地址，s=A  There are three formats for the data written into the PLC. One: I, Q, M data information, such as M0, M100; two: memory address, i = 100000; three: tag address, s = A (重写 [DeviceCommunicationWrite(String, Byte)](b5536dd6-ccde-6883-ac27-bff9e4b37806.htm).) |
| 公共方法代码示例 | [Write(String, String, Int32)](02fe2ca9-49a3-ede4-d1af-8c6da1f644da.htm) | 写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为ASCII  Write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, it is filled with 0 and the encoding is ASCII. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, String, Encoding)](7b602661-902d-2b70-1fdb-ba87ecd307bf.htm) | 写入字符串信息，需要指定的编码信息  Write string information, need to specify the encoding information (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, String, Int32, Encoding)](e8cb30d9-4c5e-6c68-b9af-1b1b6c2dc714.htm) | 写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为指定的编码信息  Write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, then the operation is complemented with 0 , you should specified the encoding information (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteT(T)](9f48853f-fbeb-c37d-25a2-1c76d1dee407.htm) | 写入支持Hsl特性的数据，返回是否写入成功，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Write data that supports the Hsl attribute, and return whether the write was successful. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WriteAdsState](c05364e4-eeb6-6b26-e817-0413145e4e8b.htm) | 写入Ads的状态，可以携带数据信息，数据可以为空  Write the status of Ads, can carry data information, and the data can be empty |
| 公共方法 | [WriteAdsStateAsync](7a6c4cf6-7837-0873-3d9d-e6b99acec800.htm) | 写入Ads的状态，可以携带数据信息，数据可以为空  Write the status of Ads, can carry data information, and the data can be empty |
| 公共方法 | [WriteAsync(String, Boolean)](10fe669c-805e-1e5a-f016-f319ad363e90.htm) | 异步批量写入Boolean数组数据，返回是否成功  Asynchronously batch write Boolean array data, return success (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Double)](068adaa4-d1de-dffc-8376-6f483e2c1d96.htm) | 异步写入double数据，返回是否成功  Asynchronously write double data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Double)](36199b83-c4ea-b02f-3675-fb6ac48d0c59.htm) | 异步写入double数组，返回是否成功  Asynchronously write double array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int16)](ccd95ccc-10a3-af67-62e8-c2312b4871b0.htm) | 异步写入short数据，返回是否成功  Asynchronously write short data, returns whether success (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int16)](2625e868-7096-e526-ce25-b29afc639b41.htm) | 异步写入short数组，返回是否成功  Asynchronously write short array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int32)](df91d629-2080-328d-67cd-fbefb9a45404.htm) | 异步写入int数据，返回是否成功  Asynchronously write int data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int32)](99962a08-9624-cea6-6009-5e6f207236b7.htm) | 异步写入int[]数组，返回是否成功  Asynchronously write int array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int64)](ec6c1e24-a0d8-f1f5-56af-109b0c546fc8.htm) | 异步写入long数据，返回是否成功  Asynchronously write long data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int64)](4b2ab94c-2891-e0a9-f10b-b67c63bd1cf9.htm) | 异步写入long数组，返回是否成功  Asynchronously write long array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Single)](910189df-13f9-b7d0-a626-b5b25ce5d32a.htm) | 异步写入float数据，返回是否成功  Asynchronously write float data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Single)](dbb2359a-fb54-cf93-4413-20d5a4ef9d36.htm) | 异步写入float数组，返回是否成功  Asynchronously write float array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, String)](fe839c9a-46f3-e317-0155-b10a8a83a9b1.htm) | 异步写入字符串信息，编码为ASCII  Asynchronously write string information, encoded as ASCII (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt16)](a04f3aa2-a51e-26cd-974e-bf70a2db359d.htm) | 异步写入ushort数据，返回是否成功  Asynchronously write ushort data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt16)](34a1a321-972e-06a8-76fb-794dd1229005.htm) | 异步写入ushort数组，返回是否成功  Asynchronously write ushort array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt32)](abad215c-bc99-ac90-b9f3-7adb3a7cca8f.htm) | 异步写入uint数据，返回是否成功  Asynchronously write uint data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt32)](5b655bf3-5ba7-ad31-6897-89c59051ac60.htm) | 异步写入uint[]数组，返回是否成功  Asynchronously write uint array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt64)](e396074d-6cf1-b45d-5e8e-01026f6e391b.htm) | 异步写入ulong数据，返回是否成功  Asynchronously write ulong data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt64)](ae3df8ec-29c0-42ff-3d6a-68f496997b3d.htm) | 异步写入ulong数组，返回是否成功  Asynchronously write ulong array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WriteAsync(String, Boolean)](c97f4f59-2c83-022b-9612-d004221cfc2a.htm) | 写入PLC的数据，地址共有三种格式，一：I,Q,M数据信息，举例M0,M100；二：内存地址，i=100000；三：标签地址，s=A  There are three formats for the data written into the PLC. One: I, Q, M data information, such as M0, M100; two: memory address, i = 100000; three: tag address, s = A (重写 [DeviceCommunicationWriteAsync(String, Boolean)](fe872dae-7ef9-12c1-d883-e92ba8676416.htm).) |
| 公共方法 | [WriteAsync(String, Byte)](632f4ccf-893d-5312-15c1-b4653e077f3f.htm) | 写入PLC的数据，地址共有三种格式，一：I,Q,M数据信息，举例M0,M100；二：内存地址，i=100000；三：标签地址，s=A  There are three formats for the data written into the PLC. One: I, Q, M data information, such as M0, M100; two: memory address, i = 100000; three: tag address, s = A |
| 公共方法 | [WriteAsync(String, Byte)](aa8d4d98-cb38-abdf-2c6a-c1e0527d5c15.htm) | 写入PLC的数据，地址共有三种格式，一：I,Q,M数据信息，举例M0,M100；二：内存地址，i=100000；三：标签地址，s=A  There are three formats for the data written into the PLC. One: I, Q, M data information, such as M0, M100; two: memory address, i = 100000; three: tag address, s = A (重写 [DeviceCommunicationWriteAsync(String, Byte)](8f1cec30-8b86-6749-c8f3-7d39168d59df.htm).) |
| 公共方法 | [WriteAsync(String, String, Int32)](b739b983-b03c-8ecf-2713-75f19233ae2f.htm) | 异步写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为ASCII  Asynchronously write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, it is filled with 0 and the encoding is ASCII. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, String, Encoding)](876e9800-0406-6407-bdf8-8ab1cf28df38.htm) | 异步写入字符串信息，需要指定的编码信息  Asynchronously write string information, need to specify the encoding information (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WriteAsync(String, String, Int32, Encoding)](ae9fb2a3-fddb-48eb-e9ad-6cc0ebbac1f6.htm) | 异步写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为指定的编码信息  Asynchronously write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, then the operation is complemented with 0 , you should specified the encoding information (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsyncT(T)](25d96e77-d1b2-918c-33f9-e5e4cc57c9cd.htm) | 异步写入支持Hsl特性的数据，返回是否写入成功，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Asynchronously write data that supports the Hsl attribute, and return whether the write was successful. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteCustomerT](84a37f82-3f67-a0f8-fd3f-0630772565e4.htm) | 写入自定义类型的数据，该类型必须继承自IDataTransfer接口  Write data of a custom type, which must inherit from the IDataTransfer interface (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteCustomerAsyncT](04628397-cfa1-d785-7790-4c02cb377ab6.htm) | 写入自定义类型的数据，该类型必须继承自IDataTransfer接口  Write data of a custom type, which must inherit from the IDataTransfer interface (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WriteStructT](a7aeee39-ab11-376a-71b5-02a066630636.htm) | 将一个结构体写入到指定的地址中去，需要指定写入的起始地址  To write a structure to a specified address, you need to specify the start address of the write |

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

![](../icons/SectionExpanded.png)备注

支持的地址格式分四种，第一种是绝对的地址表示，比如M100，I100，Q100；第二种是字符串地址，采用s=aaaa;的表示方式；第三种是绝对内存地址采用i=1000000;的表示方式，第四种是自定义的index group, IG=0xF020;0 的地址  
There are four supported address formats, the first is absolute address representation, such as M100, I100, Q100; the second is string address, using s=aaaa; representation;
the third is absolute memory address using i =1000000; representation, the fourth is the custom index group, the address of IG=0xF020;0
  

| 重要事项 **重要事项** |
| --- |
| 在实际的测试中，由于打开了VS软件对倍福PLC进行编程操作，会导致HslCommunicationDemo读取PLC发生间歇性读写失败的问题，此时需要关闭Visual Studio软件对倍福的连接，之后HslCommunicationDemo就会读写成功，感谢QQ：1813782515 提供的解决思路。 |

![](../icons/SectionExpanded.png)示例

地址既支持 M100, I100，Q100 ，读取bool时，支持输入 M100.0, 也支持符号的地址，s=MAIN.a ,也支持绝对地址的形式， i=1235467;  
下面是实例化的例子，可选两种方式

实例化

[复制](# "复制")

```
private BeckhoffAdsNet beckhoff;

private void Sample1( )
{
    // 主要分为两种实例化方式，手动的NETID赋值，和自动的方式，先说明手动的方式
    beckhoff = new BeckhoffAdsNet( "127.0.0.1", 48898 );    // IP地址是TWINCAT 软件所在的电脑的IP，端口固定
    beckhoff.SetTargetAMSNetId( "192.168.64.1.1.1:851" );
    beckhoff.SetSenderAMSNetId( "192.168.64.1.1.1:35086" );
    beckhoff.UseTagCache = true;          // 启用标签符号缓存

    // 我们再来看看自动的方式是什么
    beckhoff = new BeckhoffAdsNet( "127.0.0.1", 48898 );   // IP地址是TWINCAT 软件所在的电脑的IP，端口固定
    beckhoff.AmsPort = 851;
    beckhoff.UseAutoAmsNetID = true;      // 自动AMS
    beckhoff.UseTagCache = true;          // 启用标签符号缓存
}
```

实例化之后，就可以连接操作了

连接

[复制](# "复制")

```
private void Connect( )
{
    OperateResult connect = beckhoff.ConnectServer( );
    if (connect.IsSuccess)
    {
        // 连接成功
    }
    else
    {
        // 连接失败
    }

}
```

连接成功之后，就可以进行读写操作了

读写示例

[复制](# "复制")

```
private void Read( )
{
    // 演示简单的单个读取操作
    OperateResult<short> read = beckhoff.ReadInt16( "M100" );
    if (read.IsSuccess)
    {
        Console.WriteLine( "M00: " + read.Content );
    }
    else
    {
        Console.WriteLine( "Read failed: " + read.ToMessageShowString( ) );
    }

    // 如果是读取变量信息，比如 MIAN.dd  如果是全局变量，则不需要输入MAIN.
    OperateResult<float> read_dd = beckhoff.ReadFloat( "s=MAIN.dd" );
    if (read_dd.IsSuccess)
    {
        Console.WriteLine( "MAIN.dd: " + read_dd.Content );
    }
    else
    {
        Console.WriteLine( "Read failed: " + read_dd.ToMessageShowString( ) );
    }

    // 写入也是同理的。
    beckhoff.Write( "M100", (short)123 );
    beckhoff.Write( "s=MAIN.dd", 1.23f );

    // 其他类型都是同理，我们看下数据的情况，如果在PLC里，有个变量 MIAN.ee 是 bool数组，ee: ARRAY[0..10] OF BOOL;
    OperateResult<bool[]> read_bool = beckhoff.ReadBool( "s=MAIN.ee", 11 );
    if (read_bool.IsSuccess)
    {
        // 一次性全部读回来，写入同理
    }

    // 如果要读取数组中的某个值
    OperateResult<bool> read_bool2 = beckhoff.ReadBool( "s=MAIN.ee[2]" );
    if (read_bool2.IsSuccess)
    {
        // 读取了索引2的值，注意，读取某个值的时候，只能读取1个数据
    }
    // 写入同理
    beckhoff.Write( "s=MAIN.ee[7]", true );
}
```

也可以高级的批量读取，需要自己手动解析下数据

批量读取

[复制](# "复制")

```
public void Sample4( )
{
    // 此处演示批量读取，读取数组可以一定程度的提高效率，如果还想要提供效率，可以使用 批量读取。
    // 我们假设要读取4个数据
    // 1. M100           的short数据
    // 2. s=MAIN.dd      的flaot数据  PLC上显示REAL类型
    // 3. s=MAIN.a       的short数据  PLC上显示INT类型
    // 4. s=MAIN.ee      的bool数组类型，读5个bool

    OperateResult<byte[]> read = beckhoff.Read(
        new string[] {
            "M100",
            "s=MAIN.dd",
            "s=MAIN.a",
            "s=MAIN.ee"
        },
        new ushort[]
        {
            2,
            4,
            2,
            5
        } );
    if (read.IsSuccess)
    {
        // 读取成功，开始提取数据信息，总计回复 2 + 4 + 2 + 5 = 13个字节数据
        short m100 = beckhoff.ByteTransform.TransInt16( read.Content, 0 );
        float dd   = beckhoff.ByteTransform.TransSingle( read.Content, 2 );
        short a    = beckhoff.ByteTransform.TransInt16( read.Content, 6 );
        bool[] ee  = read.Content.SelectMiddle( 8, 5 ).Select( m => m != 0x00 ).ToArray( );  // 示例  01 00 01 00 01  为00 就是false，01就是true，其他不可识别
    }

    // 特别说明，如果存在符号地址，第一次读取会慢一点，因为还需要读取符号地址对应的内存地址，后续的读取就会很快
    // In particular, if there is a symbolic address, the first read will be slower, because the memory address corresponding to the symbolic address needs to be read, and subsequent reads will be very fast
}
```

当然，还可以进一步，既实现了批量的高性能读取，又自动解析。

类型读取

[复制](# "复制")

```
// 此处演示类型读取，读取类型也可以一定程度的提高效率，还方便的支持自动的解析。
// 我们假设要读取5个数据
// The type reading is demonstrated here. Reading the type can also improve the efficiency to a certain extent, and it is also convenient to support automatic parsing.

// 1. s=MAIN.a       的short数据  PLC上显示INT类型                 INT type displayed on PLC
// 2. s=MAIN.cc      的int数据    PLC上显示DINT类型                Display DINT type on PLC
// 3. s=MAIN.dd      的flaot数据  PLC上显示REAL类型                REAL type displayed on PLC
// 4. s=MAIN.ee      的bool数组类型，PLC显示11个长度的bool数组      The bool array type, PLC displays a bool array of 11 lengths
// 5. s=MAIN.ff      的int数组， PLC显示3个长度的INT数组            The int array, the PLC displays an int array of 3 lengths

// 定义类对象，类名及属性名可以自定义                               Define class objects, class names and attribute names can be customized
public class MyData
{
    [HslCommunication.Reflection.HslDeviceAddress("s=MAIN.a")]
    public short A { get; set; }

    [HslCommunication.Reflection.HslDeviceAddress( "s=MAIN.cc" )]
    public int CC { get; set; }

    [HslCommunication.Reflection.HslDeviceAddress( "s=MAIN.dd" )]
    public float DD { get; set; }

    [HslCommunication.Reflection.HslDeviceAddress( "s=MAIN.ee", 11 )]
    public bool[] EE { get; set; }

    [HslCommunication.Reflection.HslDeviceAddress( "s=MAIN.ff", 3 )]
    public short[] FF { get; set; }
}

// 然后下面方法演示读取示例
public void Sample5( )
{
    OperateResult<MyData> read = beckhoff.Read<MyData>( );
    if (read.IsSuccess)
    {
        // 读取成功，打印出来看看                                         Read successfully, print it out to see
        Console.WriteLine( "A: "  + read.Content.A );
        Console.WriteLine( "CC: " + read.Content.CC );
        Console.WriteLine( "DD: " + read.Content.DD );
        Console.WriteLine( "EE: " + read.Content.EE.ToArrayString( ) );
        Console.WriteLine( "FF: " + read.Content.FF.ToArrayString( ) );
    }

    // 特别说明，如果存在符号地址，第一次读取会慢一点，因为还需要读取符号地址对应的内存地址，后续的读取就会很快
    // In particular, if there is a symbolic address, the first read will be slower, because the memory address corresponding to the symbolic address needs to be read, and subsequent reads will be very fast
}
```

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.Profinet.Beckhoff 命名空间](fb2af334-3662-b9ca-be3c-b0090bf4ef56.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## BeckhoffAdsNet 构造函数 

[原文連結](http://api.hslcommunication.cn/html/eb96f2a7-8a9b-62bb-52e2-172e65d2d879.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Beckhoff](../html/fb2af334-3662-b9ca-be3c-b0090bf4ef56.htm "HslCommunication.Profinet.Beckhoff")

[BeckhoffAdsNet 类](../html/7ab15712-8688-063d-6b30-4b24803393fd.htm "BeckhoffAdsNet 类")

[BeckhoffAdsNet 构造函数](../html/eb96f2a7-8a9b-62bb-52e2-172e65d2d879.htm "BeckhoffAdsNet 构造函数 ")

[BeckhoffAdsNet 构造函数](../html/c6da3e11-3f78-2433-01d0-0d19145efc44.htm "BeckhoffAdsNet 构造函数 ")

[BeckhoffAdsNet 构造函数 (String, Int32)](../html/72cd48b9-6409-df6b-b93e-03220e25bd24.htm "BeckhoffAdsNet 构造函数 (String, Int32)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| BeckhoffAdsNet 构造函数 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [BeckhoffAdsNet](c6da3e11-3f78-2433-01d0-0d19145efc44.htm) | 实例化一个默认的对象  Instantiate a default object |
| 公共方法 | [BeckhoffAdsNet(String, Int32)](72cd48b9-6409-df6b-b93e-03220e25bd24.htm) | 通过指定的ip地址以及端口号实例化一个默认的对象  Instantiate a default object with the specified IP address and port number |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[BeckhoffAdsNet 类](7ab15712-8688-063d-6b30-4b24803393fd.htm)

[HslCommunication.Profinet.Beckhoff 命名空间](fb2af334-3662-b9ca-be3c-b0090bf4ef56.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## BeckhoffAdsNet 构造函数 

[原文連結](http://api.hslcommunication.cn/html/c6da3e11-3f78-2433-01d0-0d19145efc44.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Beckhoff](../html/fb2af334-3662-b9ca-be3c-b0090bf4ef56.htm "HslCommunication.Profinet.Beckhoff")

[BeckhoffAdsNet 类](../html/7ab15712-8688-063d-6b30-4b24803393fd.htm "BeckhoffAdsNet 类")

[BeckhoffAdsNet 构造函数](../html/eb96f2a7-8a9b-62bb-52e2-172e65d2d879.htm "BeckhoffAdsNet 构造函数 ")

[BeckhoffAdsNet 构造函数](../html/c6da3e11-3f78-2433-01d0-0d19145efc44.htm "BeckhoffAdsNet 构造函数 ")

[BeckhoffAdsNet 构造函数 (String, Int32)](../html/72cd48b9-6409-df6b-b93e-03220e25bd24.htm "BeckhoffAdsNet 构造函数 (String, Int32)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| BeckhoffAdsNet 构造函数 |

实例化一个默认的对象  
Instantiate a default object

**命名空间：**
 [HslCommunication.Profinet.Beckhoff](fb2af334-3662-b9ca-be3c-b0090bf4ef56.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public BeckhoffAdsNet()
```

```
Public Sub New
```

```
public:
BeckhoffAdsNet()
```

```
new : unit -> BeckhoffAdsNet
```

![](../icons/SectionExpanded.png)参见

#### 引用

[BeckhoffAdsNet 类](7ab15712-8688-063d-6b30-4b24803393fd.htm)

[BeckhoffAdsNet 重载](eb96f2a7-8a9b-62bb-52e2-172e65d2d879.htm)

[HslCommunication.Profinet.Beckhoff 命名空间](fb2af334-3662-b9ca-be3c-b0090bf4ef56.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## BeckhoffAdsNet 构造函数 (String, Int32)

[原文連結](http://api.hslcommunication.cn/html/72cd48b9-6409-df6b-b93e-03220e25bd24.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Beckhoff](../html/fb2af334-3662-b9ca-be3c-b0090bf4ef56.htm "HslCommunication.Profinet.Beckhoff")

[BeckhoffAdsNet 类](../html/7ab15712-8688-063d-6b30-4b24803393fd.htm "BeckhoffAdsNet 类")

[BeckhoffAdsNet 构造函数](../html/eb96f2a7-8a9b-62bb-52e2-172e65d2d879.htm "BeckhoffAdsNet 构造函数 ")

[BeckhoffAdsNet 构造函数](../html/c6da3e11-3f78-2433-01d0-0d19145efc44.htm "BeckhoffAdsNet 构造函数 ")

[BeckhoffAdsNet 构造函数 (String, Int32)](../html/72cd48b9-6409-df6b-b93e-03220e25bd24.htm "BeckhoffAdsNet 构造函数 (String, Int32)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| BeckhoffAdsNet 构造函数 (String, Int32) |

通过指定的ip地址以及端口号实例化一个默认的对象  
Instantiate a default object with the specified IP address and port number

**命名空间：**
 [HslCommunication.Profinet.Beckhoff](fb2af334-3662-b9ca-be3c-b0090bf4ef56.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public BeckhoffAdsNet(
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
BeckhoffAdsNet(
	String^ ipAddress, 
	int port
)
```

```
new : 
        ipAddress : string * 
        port : int -> BeckhoffAdsNet
```

#### 参数

ipAddress
:   类型：SystemString  
    IP地址信息

port
:   类型：SystemInt32  
    端口号

![](../icons/SectionExpanded.png)参见

#### 引用

[BeckhoffAdsNet 类](7ab15712-8688-063d-6b30-4b24803393fd.htm)

[BeckhoffAdsNet 重载](eb96f2a7-8a9b-62bb-52e2-172e65d2d879.htm)

[HslCommunication.Profinet.Beckhoff 命名空间](fb2af334-3662-b9ca-be3c-b0090bf4ef56.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## BeckhoffAdsNet 属性

[原文連結](http://api.hslcommunication.cn/html/7ce17828-db34-c0f5-016c-ec693f153a7f.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Beckhoff](../html/fb2af334-3662-b9ca-be3c-b0090bf4ef56.htm "HslCommunication.Profinet.Beckhoff")

[BeckhoffAdsNet 类](../html/7ab15712-8688-063d-6b30-4b24803393fd.htm "BeckhoffAdsNet 类")

[BeckhoffAdsNet 属性](../html/7ce17828-db34-c0f5-016c-ec693f153a7f.htm "BeckhoffAdsNet 属性")

[AmsPort 属性](../html/69677d41-296b-f1e9-4c91-1317977145c8.htm "AmsPort 属性 ")

[IpAddress 属性](../html/dafc8bcd-2bb3-a0a1-6191-dcadd9e01293.htm "IpAddress 属性 ")

[UseAutoAmsNetID 属性](../html/d4280c64-cf22-b428-bc3f-3ec08522754c.htm "UseAutoAmsNetID 属性 ")

[UseTagCache 属性](../html/94a53279-b7cc-a5f5-ce4c-93e224e74b29.htm "UseTagCache 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| BeckhoffAdsNet 属性 |

[BeckhoffAdsNet](7ab15712-8688-063d-6b30-4b24803393fd.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [AmsPort](69677d41-296b-f1e9-4c91-1317977145c8.htm) | 获取或设置Ams的端口号信息，TwinCAT2，端口号801,811,821,831；TwinCAT3，端口号为851,852,853  Get or set the port number information of Ams, TwinCAT2, the port number is 801, 811, 821, 831; TwinCAT3, the port number is 851, 852, 853 |
| 公共属性代码示例 | [ByteTransform](f35118c9-3691-f752-ae6f-2c5549714faa.htm) | 当前的数据变换机制，当你需要从字节数据转换类型数据的时候需要。  The current data transformation mechanism is required when you need to convert type data from byte data. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共属性 | [CommunicationPipe](ad8e903b-d289-21b0-077e-15ede8005221.htm) | 获取或设置当前的管道信息，管道类型为[CommunicationPipe](ad8e903b-d289-21b0-077e-15ede8005221.htm)的继承类，内置了[PipeTcpNet](26b94bd5-2695-36e7-6378-c3d41846be8f.htm)管道，[PipeUdpNet](5a0b810e-187e-826b-0a13-34fd5576aa0e.htm)管道，[PipeSerialPort](9d924607-3dfc-3916-b532-c5dec48fe347.htm)管道等  Get or set the current pipeline information, the pipeline type is [CommunicationPipe](ad8e903b-d289-21b0-077e-15ede8005221.htm) inheritance class, [PipeTcpNet](26b94bd5-2695-36e7-6378-c3d41846be8f.htm) pipeline, [PipeUdpNet](5a0b810e-187e-826b-0a13-34fd5576aa0e.htm) pipeline, [PipeSerialPort](9d924607-3dfc-3916-b532-c5dec48fe347.htm) pipeline, etc (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [ConnectionId](ce1efb4b-fbc9-4683-163e-0a8e71e99d6b.htm) | 当前连接的唯一ID号，默认为长度20的guid码加随机数组成，方便列表管理，也可以自己指定  The unique ID number of the current connection. The default is a 20-digit guid code plus a random number. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性代码示例 | [ConnectTimeOut](71a651e5-737d-6e88-75a2-fb891b9af6cd.htm) | 获取或设置接收服务器反馈的时间，如果为负数，则不接收反馈   Gets or sets the time to receive server feedback, and if it is a negative number, does not receive feedback (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共属性代码示例 | [IpAddress](dafc8bcd-2bb3-a0a1-6191-dcadd9e01293.htm) | 获取或是设置远程服务器的IP地址，如果是本机测试，那么需要设置为127.0.0.1   Get or set the IP address of the remote server. If it is a local test, then it needs to be set to 127.0.0.1 (重写 [DeviceTcpNetIpAddress](c378f4e8-3fdf-1932-7239-ed44bc5af584.htm).) |
| 公共属性 | [LocalBinding](622b418f-5969-3726-5dcc-04456c4b97a7.htm) | 获取或设置绑定的本地的IP地址和端口号信息，如果端口设置为0，代表任何可用的端口  Get or set the bound local IP address and port number information, if the port is set to 0, it means any available port (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共属性代码示例 | [LogNet](8c56f246-88ac-afb8-8c06-972fafe4f78a.htm) | 组件的日志工具，支持日志记录，只要实例化后，当前网络的基本信息，就以[DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm)等级进行输出  The component's logging tool supports logging. As long as the instantiation of the basic network information, the output will be output at [DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm) (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性代码示例 | [Port](6f1ad24a-c916-5b67-aa17-a8c7dc2acb2e.htm) | 获取或设置服务器的端口号，具体的值需要取决于对方的配置  Gets or sets the port number of the server. The specific value depends on the configuration of the other party. (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共属性代码示例 | [ReceiveTimeOut](250560a7-332b-ce1b-5eb0-e3abe4ee6012.htm) | 获取或设置接收服务器反馈的时间，如果为负数，则不接收反馈   Gets or sets the time to receive server feedback, and if it is a negative number, does not receive feedback (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [SendBeforeHex](72aa1de4-116f-abfe-4adb-814a0be0e560.htm) | 获取或设置在发送通信报文前追加发送的字节信息，HEX格式，通常用于lora组网时，需要携带 00 00 00 02 四个字节的站地址功能。  Obtain or set the byte information sent before sending communication packets, HEX format, usually used for LORA networking, you need to carry 00 00 00 02 four-byte station address function. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [SleepTime](86a5fbb9-3456-6beb-7113-5bbad5eabbc5.htm) | 获取或设置在正式接收对方返回数据前的时候，需要休息的时间，当设置为0的时候，不需要休息。  Get or set the time required to rest before officially receiving the data from the other party. When it is set to 0, no rest is required. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [SocketKeepAliveTime](96ab7703-8c4d-b4c2-ddd9-0dc93ebcfc8d.htm) | 获取或设置客户端的Socket的心跳时间信息，这个是Socket底层自动实现的心跳包，不基于协议层实现。默认小于0，不开启心跳检测，如果需要开启，设置 60\_000 比较合适，单位毫秒  Get or set the heartbeat time information of the Socket of the client. This is the heartbeat packet automatically implemented by the bottom layer of the Socket, not based on the protocol layer. The default value is less than 0, and heartbeat detection is not enabled. If you need to enable it, it is more appropriate to set 60\_000, in milliseconds. (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共属性 | [UseAutoAmsNetID](d4280c64-cf22-b428-bc3f-3ec08522754c.htm) | 是否使用服务器自动的NETID信息，默认手动设置  Whether to use the server's automatic NETID information, manually set by default |
| 公共属性 | [UseTagCache](94a53279-b7cc-a5f5-ce4c-93e224e74b29.htm) | 是否使用标签的名称缓存功能，默认为 False  Whether to use tag name caching. The default is False |
| 受保护的属性 | [WordLength](f1a505cc-e85c-6c8b-b189-58a1ab5753f9.htm) | 一个字单位的数据表示的地址长度，西门子为2，三菱，欧姆龙，modbusTcp就为1，AB PLC无效  The address length represented by one word of data, Siemens is 2, Mitsubishi, Omron, modbusTcp is 1, AB PLC is invalid (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[BeckhoffAdsNet 类](7ab15712-8688-063d-6b30-4b24803393fd.htm)

[HslCommunication.Profinet.Beckhoff 命名空间](fb2af334-3662-b9ca-be3c-b0090bf4ef56.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## AmsPort 属性 

[原文連結](http://api.hslcommunication.cn/html/69677d41-296b-f1e9-4c91-1317977145c8.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Beckhoff](../html/fb2af334-3662-b9ca-be3c-b0090bf4ef56.htm "HslCommunication.Profinet.Beckhoff")

[BeckhoffAdsNet 类](../html/7ab15712-8688-063d-6b30-4b24803393fd.htm "BeckhoffAdsNet 类")

[BeckhoffAdsNet 属性](../html/7ce17828-db34-c0f5-016c-ec693f153a7f.htm "BeckhoffAdsNet 属性")

[AmsPort 属性](../html/69677d41-296b-f1e9-4c91-1317977145c8.htm "AmsPort 属性 ")

[IpAddress 属性](../html/dafc8bcd-2bb3-a0a1-6191-dcadd9e01293.htm "IpAddress 属性 ")

[UseAutoAmsNetID 属性](../html/d4280c64-cf22-b428-bc3f-3ec08522754c.htm "UseAutoAmsNetID 属性 ")

[UseTagCache 属性](../html/94a53279-b7cc-a5f5-ce4c-93e224e74b29.htm "UseTagCache 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| BeckhoffAdsNetAmsPort 属性 |

获取或设置Ams的端口号信息，TwinCAT2，端口号801,811,821,831；TwinCAT3，端口号为851,852,853  
Get or set the port number information of Ams, TwinCAT2, the port number is 801, 811, 821, 831; TwinCAT3, the port number is 851, 852, 853

**命名空间：**
 [HslCommunication.Profinet.Beckhoff](fb2af334-3662-b9ca-be3c-b0090bf4ef56.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public int AmsPort { get; set; }
```

```
Public Property AmsPort As Integer
	Get
	Set
```

```
public:
property int AmsPort {
	int get ();
	void set (int value);
}
```

```
member AmsPort : int with get, set
```

#### 属性值

类型：Int32

![](../icons/SectionExpanded.png)参见

#### 引用

[BeckhoffAdsNet 类](7ab15712-8688-063d-6b30-4b24803393fd.htm)

[HslCommunication.Profinet.Beckhoff 命名空间](fb2af334-3662-b9ca-be3c-b0090bf4ef56.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## IpAddress 属性 

[原文連結](http://api.hslcommunication.cn/html/dafc8bcd-2bb3-a0a1-6191-dcadd9e01293.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Beckhoff](../html/fb2af334-3662-b9ca-be3c-b0090bf4ef56.htm "HslCommunication.Profinet.Beckhoff")

[BeckhoffAdsNet 类](../html/7ab15712-8688-063d-6b30-4b24803393fd.htm "BeckhoffAdsNet 类")

[BeckhoffAdsNet 属性](../html/7ce17828-db34-c0f5-016c-ec693f153a7f.htm "BeckhoffAdsNet 属性")

[AmsPort 属性](../html/69677d41-296b-f1e9-4c91-1317977145c8.htm "AmsPort 属性 ")

[IpAddress 属性](../html/dafc8bcd-2bb3-a0a1-6191-dcadd9e01293.htm "IpAddress 属性 ")

[UseAutoAmsNetID 属性](../html/d4280c64-cf22-b428-bc3f-3ec08522754c.htm "UseAutoAmsNetID 属性 ")

[UseTagCache 属性](../html/94a53279-b7cc-a5f5-ce4c-93e224e74b29.htm "UseTagCache 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| BeckhoffAdsNetIpAddress 属性 |

获取或是设置远程服务器的IP地址，如果是本机测试，那么需要设置为127.0.0.1   
Get or set the IP address of the remote server. If it is a local test, then it needs to be set to 127.0.0.1

**命名空间：**
 [HslCommunication.Profinet.Beckhoff](fb2af334-3662-b9ca-be3c-b0090bf4ef56.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public override string IpAddress { get; set; }
```

```
Public Overrides Property IpAddress As String
	Get
	Set
```

```
public:
virtual property String^ IpAddress {
	String^ get () override;
	void set (String^ value) override;
}
```

```
abstract IpAddress : string with get, set
override IpAddress : string with get, set
```

#### 属性值

类型：String

![](../icons/SectionExpanded.png)备注

最好实在初始化的时候进行指定，当使用短连接的时候，支持动态更改，切换；当使用长连接后，无法动态更改  
支持使用域名的网址方式，例如：www.hslcommunication.cn

![](../icons/SectionExpanded.png)示例

以下举例modbus-tcp的短连接及动态更改ip地址的示例

IpAddress示例

[复制](# "复制")

```
ModbusTcpNet modbus = new ModbusTcpNet( "192.168.0.100" );
// 读取线圈100的值
bool coil_ip100 = modbus.ReadCoil( "100" ).Content;
// 切换ip地址
modbus.IpAddress = "192.168.0.101";
bool coil_ip101 = modbus.ReadCoil( "100" ).Content;
```

![](../icons/SectionExpanded.png)参见

#### 引用

[BeckhoffAdsNet 类](7ab15712-8688-063d-6b30-4b24803393fd.htm)

[HslCommunication.Profinet.Beckhoff 命名空间](fb2af334-3662-b9ca-be3c-b0090bf4ef56.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## UseAutoAmsNetID 属性 

[原文連結](http://api.hslcommunication.cn/html/d4280c64-cf22-b428-bc3f-3ec08522754c.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Beckhoff](../html/fb2af334-3662-b9ca-be3c-b0090bf4ef56.htm "HslCommunication.Profinet.Beckhoff")

[BeckhoffAdsNet 类](../html/7ab15712-8688-063d-6b30-4b24803393fd.htm "BeckhoffAdsNet 类")

[BeckhoffAdsNet 属性](../html/7ce17828-db34-c0f5-016c-ec693f153a7f.htm "BeckhoffAdsNet 属性")

[AmsPort 属性](../html/69677d41-296b-f1e9-4c91-1317977145c8.htm "AmsPort 属性 ")

[IpAddress 属性](../html/dafc8bcd-2bb3-a0a1-6191-dcadd9e01293.htm "IpAddress 属性 ")

[UseAutoAmsNetID 属性](../html/d4280c64-cf22-b428-bc3f-3ec08522754c.htm "UseAutoAmsNetID 属性 ")

[UseTagCache 属性](../html/94a53279-b7cc-a5f5-ce4c-93e224e74b29.htm "UseTagCache 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| BeckhoffAdsNetUseAutoAmsNetID 属性 |

是否使用服务器自动的NETID信息，默认手动设置  
Whether to use the server's automatic NETID information, manually set by default

**命名空间：**
 [HslCommunication.Profinet.Beckhoff](fb2af334-3662-b9ca-be3c-b0090bf4ef56.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public bool UseAutoAmsNetID { get; set; }
```

```
Public Property UseAutoAmsNetID As Boolean
	Get
	Set
```

```
public:
property bool UseAutoAmsNetID {
	bool get ();
	void set (bool value);
}
```

```
member UseAutoAmsNetID : bool with get, set
```

#### 属性值

类型：Boolean

![](../icons/SectionExpanded.png)参见

#### 引用

[BeckhoffAdsNet 类](7ab15712-8688-063d-6b30-4b24803393fd.htm)

[HslCommunication.Profinet.Beckhoff 命名空间](fb2af334-3662-b9ca-be3c-b0090bf4ef56.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## UseTagCache 属性 

[原文連結](http://api.hslcommunication.cn/html/94a53279-b7cc-a5f5-ce4c-93e224e74b29.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Beckhoff](../html/fb2af334-3662-b9ca-be3c-b0090bf4ef56.htm "HslCommunication.Profinet.Beckhoff")

[BeckhoffAdsNet 类](../html/7ab15712-8688-063d-6b30-4b24803393fd.htm "BeckhoffAdsNet 类")

[BeckhoffAdsNet 属性](../html/7ce17828-db34-c0f5-016c-ec693f153a7f.htm "BeckhoffAdsNet 属性")

[AmsPort 属性](../html/69677d41-296b-f1e9-4c91-1317977145c8.htm "AmsPort 属性 ")

[IpAddress 属性](../html/dafc8bcd-2bb3-a0a1-6191-dcadd9e01293.htm "IpAddress 属性 ")

[UseAutoAmsNetID 属性](../html/d4280c64-cf22-b428-bc3f-3ec08522754c.htm "UseAutoAmsNetID 属性 ")

[UseTagCache 属性](../html/94a53279-b7cc-a5f5-ce4c-93e224e74b29.htm "UseTagCache 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| BeckhoffAdsNetUseTagCache 属性 |

是否使用标签的名称缓存功能，默认为 False  
Whether to use tag name caching. The default is False

**命名空间：**
 [HslCommunication.Profinet.Beckhoff](fb2af334-3662-b9ca-be3c-b0090bf4ef56.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public bool UseTagCache { get; set; }
```

```
Public Property UseTagCache As Boolean
	Get
	Set
```

```
public:
property bool UseTagCache {
	bool get ();
	void set (bool value);
}
```

```
member UseTagCache : bool with get, set
```

#### 属性值

类型：Boolean

![](../icons/SectionExpanded.png)参见

#### 引用

[BeckhoffAdsNet 类](7ab15712-8688-063d-6b30-4b24803393fd.htm)

[HslCommunication.Profinet.Beckhoff 命名空间](fb2af334-3662-b9ca-be3c-b0090bf4ef56.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## BeckhoffAdsNet 方法

[原文連結](http://api.hslcommunication.cn/html/e8bad1d3-d4ae-37ef-6ba7-e199dac65f0c.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Beckhoff](../html/fb2af334-3662-b9ca-be3c-b0090bf4ef56.htm "HslCommunication.Profinet.Beckhoff")

[BeckhoffAdsNet 类](../html/7ab15712-8688-063d-6b30-4b24803393fd.htm "BeckhoffAdsNet 类")

[BeckhoffAdsNet 方法](../html/e8bad1d3-d4ae-37ef-6ba7-e199dac65f0c.htm "BeckhoffAdsNet 方法")

[ClearTagsCache 方法](../html/3ca05c73-2ab5-b7b8-2d46-03da8a9c67ad.htm "ClearTagsCache 方法 ")

[DecideWhetherQAMessage 方法](../html/771586fe-9184-5312-e3da-bae61fe91d9a.htm "DecideWhetherQAMessage 方法 ")

[ExtraAfterReadFromCoreServer 方法](../html/3d26d83d-9287-54f5-3d35-c6d074ce6263.htm "ExtraAfterReadFromCoreServer 方法 ")

[GetNewNetMessage 方法](../html/215ca52b-b151-9cb6-0443-00b9be103188.htm "GetNewNetMessage 方法 ")

[GetSenderAMSNetId 方法](../html/f743d65e-7896-7d75-ba57-6be4991ab9af.htm "GetSenderAMSNetId 方法 ")

[GetTargetAMSNetId 方法](../html/dcc58722-6399-8dd4-affa-2e531a098ae0.htm "GetTargetAMSNetId 方法 ")

[InitializationOnConnect 方法](../html/ea40c0b0-f579-3a64-7b45-573259ef706d.htm "InitializationOnConnect 方法 ")

[InitializationOnConnectAsync 方法](../html/c2a3f0f3-3a78-e102-d07f-d15053af2633.htm "InitializationOnConnectAsync 方法 ")

[PackCommandWithHeader 方法](../html/3895f6bc-bf81-5dbe-cb27-c12c6602126d.htm "PackCommandWithHeader 方法 ")

[Read 方法](../html/504e0a30-6e4b-ac79-e7e4-ce3aecedecb9.htm "Read 方法 ")

[ReadAdsDeviceInfo 方法](../html/1b58e4ea-7a7c-3b42-1f11-b077010b9a39.htm "ReadAdsDeviceInfo 方法 ")

[ReadAdsDeviceInfoAsync 方法](../html/18ef1d19-2639-b4b9-99c6-3e7caed8ef5a.htm "ReadAdsDeviceInfoAsync 方法 ")

[ReadAdsState 方法](../html/4645af4b-15b0-49a7-a3d3-bc07ef913ad7.htm "ReadAdsState 方法 ")

[ReadAdsStateAsync 方法](../html/831eac3f-b313-b2bd-3131-71bf3948b652.htm "ReadAdsStateAsync 方法 ")

[ReadAsync 方法](../html/ffaa1904-b9ab-79ff-efde-2543c610aa26.htm "ReadAsync 方法 ")

[ReadBool 方法](../html/b4ab6787-c834-7bf7-e069-1e8f7d5a4cad.htm "ReadBool 方法 ")

[ReadBoolAsync 方法](../html/1114ae0d-0359-668d-a6b2-6be76e7da28b.htm "ReadBoolAsync 方法 ")

[ReadByte 方法](../html/dadb0178-8d9c-2ce1-5495-e53d72771b82.htm "ReadByte 方法 ")

[ReadByteAsync 方法](../html/d61db681-dca5-f888-2dfd-c35faa269a4f.htm "ReadByteAsync 方法 ")

[ReadStruct 方法](../html/8e09ddd7-faa8-67c9-e099-eee24d59aa8d.htm "ReadStruct 方法 ")

[ReadValueHandle 方法](../html/7a0c126f-5f3e-f077-f90b-34ca2191c5b2.htm "ReadValueHandle 方法 ")

[ReadValueHandleAsync 方法](../html/0dadc9a4-90cc-1f29-627f-f8112aaa7df2.htm "ReadValueHandleAsync 方法 ")

[ReleaseSystemHandle 方法](../html/e2fa36f0-6d1f-2e02-e3de-bc2344e8d553.htm "ReleaseSystemHandle 方法 ")

[ReleaseSystemHandleAsync 方法](../html/cb4896ea-bbd4-52fa-4587-db7c137ff2d6.htm "ReleaseSystemHandleAsync 方法 ")

[SetSenderAMSNetId 方法](../html/6c085271-688a-909e-53c9-dec86b2a8e8a.htm "SetSenderAMSNetId 方法 ")

[SetTargetAMSNetId 方法](../html/af74dcd7-fb14-527d-5b79-881112e412ad.htm "SetTargetAMSNetId 方法 ")

[ToString 方法](../html/5cf941e3-75eb-cfcb-ec45-0065be0058c5.htm "ToString 方法 ")

[TransValueHandle 方法](../html/f7d19175-e124-9a8d-d36e-638b14d9c919.htm "TransValueHandle 方法 ")

[TransValueHandleAsync 方法](../html/efa735b1-7770-c061-0e6f-0486b84ff6b8.htm "TransValueHandleAsync 方法 ")

[UnpackResponseContent 方法](../html/9578724a-11cf-6cec-c183-18fb69db6412.htm "UnpackResponseContent 方法 ")

[Write 方法](../html/d6d8c193-854e-1069-a9f0-d426a76813cc.htm "Write 方法 ")

[WriteAdsState 方法](../html/c05364e4-eeb6-6b26-e817-0413145e4e8b.htm "WriteAdsState 方法 ")

[WriteAdsStateAsync 方法](../html/7a6c4cf6-7837-0873-3d9d-e6b99acec800.htm "WriteAdsStateAsync 方法 ")

[WriteAsync 方法](../html/7f760de5-65c9-46e0-b53e-fa719a52dafc.htm "WriteAsync 方法 ")

[WriteStruct(T) 方法](../html/a7aeee39-ab11-376a-71b5-02a066630636.htm "WriteStruct(T) 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| BeckhoffAdsNet 方法 |

[BeckhoffAdsNet](7ab15712-8688-063d-6b30-4b24803393fd.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [ClearTagsCache](3ca05c73-2ab5-b7b8-2d46-03da8a9c67ad.htm) | 清除标签缓存信息  Clear the tag cache information |
| 公共方法代码示例 | [ConnectClose](569a9c0e-7628-b8d5-e124-0853452747b8.htm) | 手动断开与远程服务器的连接，如果当前是长连接模式，那么就会切换到短连接模式  Manually disconnect from the remote server, if it is currently in long connection mode, it will switch to short connection mode (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共方法代码示例 | [ConnectCloseAsync](76be8150-ff22-95b9-16b1-7c3440c54ae8.htm) | 手动断开与远程服务器的连接，如果当前是长连接模式，那么就会切换到短连接模式  Manually disconnect from the remote server, if it is currently in long connection mode, it will switch to short connection mode (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共方法代码示例 | [ConnectServer](6299cc1f-2fab-b754-4597-1eddc2b21308.htm) | 尝试连接远程的服务器，如果连接成功，就切换短连接模式到长连接模式，后面的每次请求都共享一个通道，使得通讯速度更快速  Try to connect to a remote server. If the connection is successful, switch the short connection mode to the long connection mode. Each subsequent request will share a channel, making the communication speed faster. (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共方法代码示例 | [ConnectServerAsync](e7eff5c0-8e2a-3b6c-5027-16f4d8bcf5ae.htm) | 尝试连接远程的服务器，如果连接成功，就切换短连接模式到长连接模式，后面的每次请求都共享一个通道，使得通讯速度更快速  Try to connect to a remote server. If the connection is successful, switch the short connection mode to the long connection mode. Each subsequent request will share a channel, making the communication speed faster. (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 受保护的方法 | [DecideWhetherQAMessage](771586fe-9184-5312-e3da-bae61fe91d9a.htm) | 决定当前的消息是否是用于问答机制返回的消息，默认直接返回 true, 实际的情况需要根据协议进行重写方法  To determine whether the current message is the message returned by the question answering mechanism, the default is true. In actual cases, the rewriting method needs to be performed according to the protocol (重写 [BinaryCommunicationDecideWhetherQAMessage(CommunicationPipe, OperateResultByte)](acc40cda-e310-8537-6685-59b67eeb16ff.htm).) |
| 公共方法 | [Dispose](b9dda6bf-d342-254e-689a-93fa0bb265d1.htm) | (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 受保护的方法 | [Dispose(Boolean)](0bbb6ed3-111e-069c-e9f5-8cbc3f7f7b77.htm) | (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | [ExtraAfterReadFromCoreServer](3d26d83d-9287-54f5-3d35-c6d074ce6263.htm) | 和服务器交互完成的时候调用的方法，可以根据读写结果进行一些额外的操作，具体的操作需要根据实际的需求来重写实现  The method called when the interaction with the server is completed can perform some additional operations based on the read and write results. The specific operations need to be rewritten according to actual needs. (重写 [BinaryCommunicationExtraAfterReadFromCoreServer(OperateResult)](06a49a15-a3ce-5e39-28c7-2b78851eca2e.htm).) |
| 受保护的方法代码示例 | [ExtraOnDisconnect](de9dbba2-3b2c-1c5e-53c8-4e6f8a44e2c4.htm) | 根据实际的协议选择是否重写本方法，有些协议在断开连接之前，需要发送一些报文来关闭当前的网络通道  Select whether to rewrite this method according to the actual protocol. Some protocols need to send some packets to close the current network channel before disconnecting. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法代码示例 | [ExtraOnDisconnectAsync](210f593e-16d8-77c4-1178-58fbcfe9b27b.htm) | 根据实际的协议选择是否重写本方法，有些协议在断开连接之前，需要发送一些报文来关闭当前的网络通道  Select whether to rewrite this method according to the actual protocol. Some protocols need to send some packets to close the current network channel before disconnecting. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 受保护的方法 | [GetLogTextFromBinary](67f3b2f7-957c-fa0b-8320-6799237157db.htm) | 获取当前的报文进行日志记录的时候，是否使用二进制的格式记录，默认返回 [LogMsgFormatBinary](2ed2ff87-caec-0e44-f08b-e2b015225a58.htm)，重写可以根据session对象分别返回不同记录模式  Whether to log the current packet in binary format, the default return is [LogMsgFormatBinary](2ed2ff87-caec-0e44-f08b-e2b015225a58.htm). If you want to override it, different recording modes can be returned according to session (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [GetNewNetMessage](215ca52b-b151-9cb6-0443-00b9be103188.htm) | 获取一个新的消息对象的方法，需要在继承类里面进行重写  The method to get a new message object needs to be overridden in the inheritance class (重写 [BinaryCommunicationGetNewNetMessage](16636a1b-a6f0-9c0f-9d9c-c94e26c67459.htm).) |
| 公共方法 | [GetSenderAMSNetId](f743d65e-7896-7d75-ba57-6be4991ab9af.htm) | 获取当前发送的AMS的网络ID信息 |
| 公共方法 | [GetTargetAMSNetId](dcc58722-6399-8dd4-affa-2e531a098ae0.htm) | 获取当前目标的AMS网络的ID信息 |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法 | [GetWordLength](c27821ea-1739-8eb9-a848-23a00828fec8.htm) | 一个字单位的数据表示的地址长度，西门子为2，三菱，欧姆龙，modbusTcp就为1，AB PLC无效  The address length represented by one word of data, Siemens is 2, Mitsubishi, Omron, modbusTcp is 1, AB PLC is invalid (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 受保护的方法代码示例 | [InitializationOnConnect](ea40c0b0-f579-3a64-7b45-573259ef706d.htm) | 根据实际的协议选择是否重写本方法，有些协议在创建连接之后，需要进行一些初始化的信号握手，才能最终建立网络通道。  Whether to rewrite this method is based on the actual protocol. Some protocols require some initial signal handshake to establish a network channel after the connection is created. (重写 [BinaryCommunicationInitializationOnConnect](4ab6989f-a5b5-05fe-8b44-bfa45824a2a7.htm).) |
| 受保护的方法代码示例 | [InitializationOnConnectAsync](c2a3f0f3-3a78-e102-d07f-d15053af2633.htm) | 根据实际的协议选择是否重写本方法，有些协议在创建连接之后，需要进行一些初始化的信号握手，才能最终建立网络通道。  Whether to rewrite this method is based on the actual protocol. Some protocols require some initial signal handshake to establish a network channel after the connection is created. (重写 [BinaryCommunicationInitializationOnConnectAsync](7da2af04-d4fc-5045-19a4-4b0d1722eef3.htm).) |
| 公共方法 | [IpAddressPing](0ac126b4-a830-a822-4ef5-855474942cd7.htm) | 对当前设备的IP地址进行PING的操作，返回PING的结果，正常来说，返回Success  PING the IP address of the current device and return the PING result. Normally, it returns Success (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 受保护的方法 | [LogRevcMessage(Byte)](dab1657f-1e75-7fe1-7003-ccedfc69855d.htm) | 使用日志记录一个接收的报文信息  Logs are used to record information about a received packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [LogRevcMessage(Byte, PipeSession)](0f51ae45-575e-4228-f978-01b37d6d5df1.htm) | 使用日志记录一个接收的报文信息  Logs are used to record information about a received packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [LogSendMessage(Byte)](4f90eecc-68a2-1ef2-44ef-369f19bed3ec.htm) | 使用日志记录一个发送的报文信息  Logs are used to record information about a send packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [LogSendMessage(Byte, PipeSession)](209caadf-4f0f-40aa-04c2-ebc61f8dab13.htm) | 使用日志记录一个发送的报文信息  Logs are used to record information about a send packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法 | [PackCommandWithHeader](3895f6bc-bf81-5dbe-cb27-c12c6602126d.htm) | 对当前的命令进行打包处理，通常是携带命令头内容，标记当前的命令的长度信息，需要进行重写，否则默认不打包  The current command is packaged, usually carrying the content of the command header, marking the length of the current command, and it needs to be rewritten, otherwise it is not packaged by default (重写 [BinaryCommunicationPackCommandWithHeader(Byte)](d0b2cd76-c306-11fc-38c7-af8ad306cce7.htm).) |
| 公共方法 | [Read(String, UInt16)](ac58c335-e208-6bf5-6d30-85dc57694090.htm) | 读取PLC的数据，地址共有三种格式，一：I,Q,M数据信息，举例M0,M100；二：内存地址，i=100000；三：标签地址，s=A  Read PLC data, there are three formats of address, one: I, Q, M data information, such as M0, M100; two: memory address, i = 100000; three: tag address, s = A (重写 [DeviceCommunicationRead(String, UInt16)](9dd09452-c200-bc65-75ff-69aa34b366d4.htm).) |
| 公共方法 | [Read(String, UInt16)](24dae660-f973-3403-a4a8-ed39d114483b.htm) | 批量读取PLC的数据，需要传入地址数组，以及读取的长度数组信息，长度单位为字节单位，如果是读取bool变量的，则以bool为单位，统一返回一串字节数据信息，需要进行二次解析的操作。  To read PLC data in batches, you need to pass in the address array and the read length array information. The unit of length is in bytes. If you read a bool variable, it will return a string of byte data information in units of bool. , which requires a secondary parsing operation. |
| 公共方法代码示例 | [ReadT](6af520c0-a616-6709-ea80-14e26306e2be.htm) | 读取支持Hsl特性的数据内容，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明，或是官网的详细文档  Read the data content of the Hsl attribute. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (重写 [DeviceCommunicationReadT](95abe064-5594-5ae7-a849-98c53f88fee0.htm).) |
| 公共方法 | [ReadAdsDeviceInfo](1b58e4ea-7a7c-3b42-1f11-b077010b9a39.htm) | 读取Ads设备的设备信息。主要是版本号，设备名称  Read the device information of the Ads device. Mainly version number, device name |
| 公共方法 | [ReadAdsDeviceInfoAsync](18ef1d19-2639-b4b9-99c6-3e7caed8ef5a.htm) | 读取Ads设备的设备信息。主要是版本号，设备名称  Read the device information of the Ads device. Mainly version number, device name |
| 公共方法 | [ReadAdsState](4645af4b-15b0-49a7-a3d3-bc07ef913ad7.htm) | 读取Ads设备的状态信息，其中[Content1](fba87d7e-fbf8-0667-b0f2-7f669d3179fe.htm)是Ads State，[Content2](312b5df6-df3c-07d6-83c2-391380a3fb33.htm)是Device State  Read the status information of the Ads device, where [Content1](fba87d7e-fbf8-0667-b0f2-7f669d3179fe.htm) is the Ads State, and [Content2](312b5df6-df3c-07d6-83c2-391380a3fb33.htm) is the Device State |
| 公共方法 | [ReadAdsStateAsync](831eac3f-b313-b2bd-3131-71bf3948b652.htm) | 读取Ads设备的状态信息，其中[Content1](fba87d7e-fbf8-0667-b0f2-7f669d3179fe.htm)是Ads State，[Content2](312b5df6-df3c-07d6-83c2-391380a3fb33.htm)是Device State  Read the status information of the Ads device, where [Content1](fba87d7e-fbf8-0667-b0f2-7f669d3179fe.htm) is the Ads State, and [Content2](312b5df6-df3c-07d6-83c2-391380a3fb33.htm) is the Device State |
| 公共方法 | [ReadAsync(String, UInt16)](82d4ecb0-9450-53ba-101a-4252209345f5.htm) | 读取PLC的数据，地址共有三种格式，一：I,Q,M数据信息，举例M0,M100；二：内存地址，i=100000；三：标签地址，s=A  Read PLC data, there are three formats of address, one: I, Q, M data information, such as M0, M100; two: memory address, i = 100000; three: tag address, s = A (重写 [DeviceCommunicationReadAsync(String, UInt16)](d8c14d37-56dc-2a31-1a19-3e2b5b65197d.htm).) |
| 公共方法 | [ReadAsync(String, UInt16)](2b8c8a10-8a40-2e2b-6917-ff16cc6a0c2d.htm) | 批量读取PLC的数据，需要传入地址数组，以及读取的长度数组信息，长度单位为字节单位，如果是读取bool变量的，则以bool为单位，统一返回一串字节数据信息，需要进行二次解析的操作。  To read PLC data in batches, you need to pass in the address array and the read length array information. The unit of length is in bytes. If you read a bool variable, it will return a string of byte data information in units of bool. , which requires a secondary parsing operation. |
| 公共方法代码示例 | [ReadAsyncT](c5887b82-93a5-0799-23a9-be999466b175.htm) | 异步读取支持Hsl特性的数据内容，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Asynchronously read the data content of the Hsl attribute. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (重写 [DeviceCommunicationReadAsyncT](0522e740-4bf0-2e31-c4c5-e41384321b37.htm).) |
| 公共方法 | [ReadBool(String)](9b2faa2e-eaa9-7ace-c49f-7a0318f1df1a.htm) | 读取单个的Boolean数据信息  Read a single Boolean data message (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadBool(String, UInt16)](92020f06-15c1-4b5f-ce8d-36a4d6182a9c.htm) | 读取PLC的数据，地址共有三种格式，一：I,Q,M数据信息，举例M0,M100；二：内存地址，i=100000；三：标签地址，s=A  Read PLC data, there are three formats of address, one: I, Q, M data information, such as M0, M100; two: memory address, i = 100000; three: tag address, s = A (重写 [DeviceCommunicationReadBool(String, UInt16)](69d0d24f-34a7-5b41-5c1b-619edc7e6b54.htm).) |
| 公共方法 | [ReadBoolAsync(String)](cc37da1f-d673-fe67-78ee-4e465d29a2bd.htm) | 异步读取单个的Boolean数据信息  Asynchronously read a single Boolean data message (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadBoolAsync(String, UInt16)](d43cef4a-36ed-4bd2-b999-ecf789b95a28.htm) | 读取PLC的数据，地址共有三种格式，一：I,Q,M数据信息，举例M0,M100；二：内存地址，i=100000；三：标签地址，s=A  Read PLC data, there are three formats of address, one: I, Q, M data information, such as M0, M100; two: memory address, i = 100000; three: tag address, s = A (重写 [DeviceCommunicationReadBoolAsync(String, UInt16)](a45d383d-0095-3ea7-32f8-e75af88d60f8.htm).) |
| 公共方法 | [ReadByte](dadb0178-8d9c-2ce1-5495-e53d72771b82.htm) | 读取PLC的数据，地址共有三种格式，一：I,Q,M数据信息，举例M0,M100；二：内存地址，i=100000；三：标签地址，s=A  Read PLC data, there are three formats of address, one: I, Q, M data information, such as M0, M100; two: memory address, i = 100000; three: tag address, s = A |
| 公共方法 | [ReadByteAsync](d61db681-dca5-f888-2dfd-c35faa269a4f.htm) | 读取PLC的数据，地址共有三种格式，一：I,Q,M数据信息，举例M0,M100；二：内存地址，i=100000；三：标签地址，s=A  Read PLC data, there are three formats of address, one: I, Q, M data information, such as M0, M100; two: memory address, i = 100000; three: tag address, s = A |
| 公共方法代码示例 | [ReadCustomerT(String)](ca069688-d8f0-7dbd-e962-86a10dce83cf.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，返回一个新的类型的实例对象。  To read a custom data type, you need to inherit from the IDataTransfer interface and return an instance object of a new type. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadCustomerT(String, T)](ca7dc1a1-c97d-689f-9819-4dae4af67206.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，传入一个实例，对这个实例进行赋值，并返回该实例的对象。  To read a custom data type, you need to inherit from the IDataTransfer interface, pass in an instance, assign a value to this instance, and return the object of the instance. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadCustomerAsyncT(String)](f3fcda2b-255d-aadd-567e-c7b3df3636d9.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，返回一个新的类型的实例对象。  To read a custom data type, you need to inherit from the IDataTransfer interface and return an instance object of a new type. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadCustomerAsyncT(String, T)](527c9fae-4743-6ebd-7b9a-5ef3a8b7c455.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，传入一个实例，对这个实例进行赋值，并返回该实例的对象。  To read a custom data type, you need to inherit from the IDataTransfer interface, pass in an instance, assign a value to this instance, and return the object of the instance. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadDouble(String)](e7e31c43-71b4-e692-6bdc-ba9533e1d83d.htm) | 读取双浮点的数据  Read double floating point data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadDouble(String, UInt16)](4eee0f05-f861-4b48-3f06-99737fb24aee.htm) | 读取双浮点数据的数组  Read double floating point data array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadDoubleAsync(String)](2cf895fa-1033-1def-f048-0e8caefaa71d.htm) | 异步读取双浮点的数据  Asynchronously read double floating point data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadDoubleAsync(String, UInt16)](2256d406-d4d5-da54-210b-ccef58a5e914.htm) | 异步读取双浮点数据的数组  Asynchronously read double floating point data array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadFloat(String)](bc2801ae-6600-e845-7c6c-522c9156eb0e.htm) | 读取单浮点数据  Read single floating point data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadFloat(String, UInt16)](6a9fb098-b13c-313b-3b27-20d3c787945d.htm) | 读取单浮点精度的数组  Read single floating point array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadFloatAsync(String)](89f743f8-4b3a-7085-1e49-a1d7adb10052.htm) | 异步读取单浮点数据  Asynchronously read single floating point data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadFloatAsync(String, UInt16)](cb6d5c5d-9b5c-2467-3e6b-5efa9d20cd72.htm) | 异步读取单浮点精度的数组  Asynchronously read single floating point array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServer(Byte)](8146dc0c-1edc-3ca7-428a-a3bdbd35e6c5.htm) | 将二进制的数据发送到管道中去，然后从管道里接收二进制的数据回来，并返回是否成功的结果对象。  Send binary data to the pipeline, and then receive binary data back from the pipeline, and return whether the success of the result object (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServer(IEnumerableByte)](044cfdd9-ec76-b89f-2f2b-5f439025006c.htm) | 将多个数据报文按顺序发到设备，并从设备接收返回的数据内容，然后拼接成一个Byte[]信息，需要重写[UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm)方法才能返回正确的结果。  Send multiple data packets to the device in sequence, and receive the returned data content from the device, and then splicing them into a Byte[] message, you need to rewrite [UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm) method to return the correct result. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServer(Byte, Boolean, Boolean)](b32ce8e8-a959-fcc2-6736-14b39fdb9a40.htm) | 将二进制的数据发送到管道中去，然后从管道里接收二进制的数据回来，并返回是否成功的结果对象。  Send binary data to the pipeline, and then receive binary data back from the pipeline, and return whether the success of the result object (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServer(CommunicationPipe, Byte, Boolean, Boolean)](ef36d56b-df09-f727-e8d3-316d37a13ee9.htm) | 使用指定的管道来进行数据通信，发送原始数据到管道，然后从管道接收相关的数据返回，本方法无锁 (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServerAsync(Byte)](a4de5a24-c67a-7b02-db05-33b85eaf6e1e.htm) | 将二进制的数据发送到管道中去，然后从管道里接收二进制的数据回来，并返回是否成功的结果对象。  Send binary data to the pipeline, and then receive binary data back from the pipeline, and return whether the success of the result object (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServerAsync(IEnumerableByte)](bcdb59b7-a70a-d5c0-1fc6-9ececcd39ea2.htm) | 将多个数据报文按顺序发到设备，并从设备接收返回的数据内容，然后拼接成一个Byte[]信息，需要重写[UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm)方法才能返回正确的结果。  Send multiple data packets to the device in sequence, and receive the returned data content from the device, and then splicing them into a Byte[] message, you need to rewrite [UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm) method to return the correct result. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServerAsync(Byte, Boolean, Boolean)](1c79c204-5ab5-9342-e65d-7455c20ea577.htm) | 将二进制的数据发送到管道中去，然后从管道里接收二进制的数据回来，并返回是否成功的结果对象。  Send binary data to the pipeline, and then receive binary data back from the pipeline, and return whether the success of the result object (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServerAsync(CommunicationPipe, Byte, Boolean, Boolean)](a254a19d-3c29-147b-335b-8adccdf4ee2f.htm) | 使用指定的管道来进行数据通信，发送原始数据到管道，然后从管道接收相关的数据返回，本方法无锁 (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法代码示例 | [ReadInt16(String)](57a16003-2fd5-b2bd-09bf-14366be899b8.htm) | 读取16位的有符号的整型数据  Read 16-bit signed integer data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt16(String, UInt16)](b2bde158-f51c-93fc-66ca-5e4cf6ea1fb7.htm) | 读取16位的有符号整型数组  Read 16-bit signed integer array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt16Async(String)](ac021793-890d-ceab-eae4-47ff763390c1.htm) | 异步读取16位的有符号的整型数据  Asynchronously read 16-bit signed integer data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt16Async(String, UInt16)](f30e02fc-8078-f1b3-c226-5991e3e23caf.htm) | 异步读取16位的有符号整型数组  Asynchronously read 16-bit signed integer array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt32(String)](bfbea035-9a83-40e3-084c-28122041b37a.htm) | 读取32位的有符号整型  Read 32-bit signed integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt32(String, UInt16)](095240a3-bc0f-8647-4867-88e40cc089e3.htm) | 读取32位有符号整型数组  Read 32-bit signed integer array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt32Async(String)](96726ea0-82aa-1229-21f6-43cd62abca43.htm) | 异步读取32位的有符号整型  Asynchronously read 32-bit signed integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt32Async(String, UInt16)](d6356095-a289-18eb-136c-b7e6ae099616.htm) | 异步读取32位有符号整型数组  Asynchronously read 32-bit signed integer array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt64(String)](0d0c2284-2da1-6f9a-b353-537561cd5e95.htm) | 读取64位的有符号整型  Read 64-bit signed integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt64(String, UInt16)](35a55019-efb1-bf2d-be1b-01c8eed99243.htm) | 读取64位的有符号整型数组  Read 64-bit signed integer array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt64Async(String)](26d46c8a-3775-3abd-d3ec-5cd1758a317e.htm) | 异步读取64位的有符号整型  Asynchronously read 64-bit signed integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt64Async(String, UInt16)](8dc97c4a-8904-9a34-9ddc-76c1a9477ee5.htm) | 异步读取64位的有符号整型数组  Asynchronously read 64-bit signed integer array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadString(String, UInt16)](07812d76-415a-6400-c8a2-1c21da51e7a3.htm) | 读取字符串数据，默认为最常见的ASCII编码  Read string data, default is the most common ASCII encoding (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadString(String, UInt16, Encoding)](b34c71a6-58f3-6bf4-0b93-471be4f2ed31.htm) | 使用指定的编码，读取字符串数据  Reads string data using the specified encoding (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadStringAsync(String, UInt16)](1835da81-b811-d151-0b53-907d8d87d0a8.htm) | 异步读取字符串数据，默认为最常见的ASCII编码  Asynchronously read string data, default is the most common ASCII encoding (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadStringAsync(String, UInt16, Encoding)](091dfc67-3d25-f33e-3a55-a0228863815e.htm) | 异步使用指定的编码，读取字符串数据  Asynchronously reads string data using the specified encoding (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadStructT(String)](b3f311a4-68d8-87b4-a972-ad8ec4b3f0df.htm) | 读取结构体的信息，传入结构体的类型，以及结构体的起始地址  Read the information of the structure, the type of the incoming structure, and the start address of the structure |
| 公共方法代码示例 | [ReadStructT(String, UInt16)](49853082-fbad-52e5-d756-0615cedb4b83.htm) | 读取结构体类型的数据，根据结构体自身的定义，读取原始字节数组，然后解析出实际的结构体数据，结构体需要实现[HslStructAttribute](87e427a8-0c7b-43a8-c421-c27264eda90f.htm)特性 (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadStructAsyncT](36461e61-d0ed-62cf-79ce-dcc7f71a3fa2.htm) | 读取结构体类型的数据，根据结构体自身的定义，读取原始字节数组，然后解析出实际的结构体数据，结构体需要实现[HslStructAttribute](87e427a8-0c7b-43a8-c421-c27264eda90f.htm)特性 (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt16(String)](736379d1-9560-5052-658a-13c126d56d78.htm) | 读取16位的无符号整型  Read 16-bit unsigned integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt16(String, UInt16)](e284c4bd-d3f8-ff30-eb42-0aa035cfb4f6.htm) | 读取16位的无符号整型数组  Read 16-bit unsigned integer array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt16Async(String)](361acc93-43b5-a9a5-6fa7-fa8728b7b96c.htm) | 异步读取16位的无符号整型  Asynchronously read 16-bit unsigned integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt16Async(String, UInt16)](73dd24b8-b501-df53-1de6-ccb2708f656d.htm) | 异步读取16位的无符号整型数组  Asynchronously read 16-bit unsigned integer array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt32(String)](ca73c958-39c7-9b34-aa44-3db2d9b0cad3.htm) | 读取32位的无符号整型  Read 32-bit unsigned integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt32(String, UInt16)](878b4709-0278-fcdf-3b65-afbf0f496c1e.htm) | 读取32位的无符号整型数组  Read 32-bit unsigned integer array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt32Async(String)](2343953f-6350-331e-6ba7-18a4449f0575.htm) | 异步读取32位的无符号整型  Asynchronously read 32-bit unsigned integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt32Async(String, UInt16)](e0eba6c1-56a2-5190-ba03-66116bcafe3d.htm) | 异步读取32位的无符号整型数组  Asynchronously read 32-bit unsigned integer array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt64(String)](947487dd-d4a3-72d8-d9fb-2224aad70e12.htm) | 读取64位的无符号整型  Read 64-bit unsigned integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt64(String, UInt16)](858766a9-9d62-87a8-f727-a00838f3d2b6.htm) | 读取64位的无符号整型的数组  Read 64-bit unsigned integer array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt64Async(String)](3411ff3a-772f-dbe1-8931-75e00a731579.htm) | 异步读取64位的无符号整型  Asynchronously read 64-bit unsigned integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt64Async(String, UInt16)](c2095a6e-ff89-dd98-8a21-fe72f82dbbc6.htm) | 异步读取64位的无符号整型的数组  Asynchronously read 64-bit unsigned integer array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadValueHandle](7a0c126f-5f3e-f077-f90b-34ca2191c5b2.htm) | 根据当前标签的地址获取到内存偏移地址  Get the memory offset address based on the address of the current label |
| 公共方法 | [ReadValueHandleAsync](0dadc9a4-90cc-1f29-627f-f8112aaa7df2.htm) | 根据当前标签的地址获取到内存偏移地址  Get the memory offset address based on the address of the current label |
| 公共方法 | [ReleaseSystemHandle](e2fa36f0-6d1f-2e02-e3de-bc2344e8d553.htm) | 释放当前的系统句柄，该句柄是通过[ReadValueHandle(String)](7a0c126f-5f3e-f077-f90b-34ca2191c5b2.htm)获取的 |
| 公共方法 | [ReleaseSystemHandleAsync](cb4896ea-bbd4-52fa-4587-db7c137ff2d6.htm) | 释放当前的系统句柄，该句柄是通过[ReadValueHandle(String)](7a0c126f-5f3e-f077-f90b-34ca2191c5b2.htm)获取的 |
| 公共方法 | [SetDtuPipe](f765d1d7-234a-b229-9457-97e25e024dc9.htm) | 将当前的通信对象设置DTU模式，允许传入现成的管道，并返回初始化结果，如果该设备重写了握手报文，就是返回握手结果  Set the current communication object to DTU mode, allow the existing pipe to be passed in, and return the initialization result, if the device rewrites the handshake packet, the handshake result is returned (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [SetDtuPipeAsync](9d1fd12d-b804-3828-a627-a2c0eef1d23f.htm) | 将当前的通信对象设置DTU模式，允许传入现成的管道，并返回初始化结果，如果该设备重写了握手报文，就是返回握手结果  Set the current communication object to DTU mode, allow the existing pipe to be passed in, and return the initialization result, if the device rewrites the handshake packet, the handshake result is returned (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [SetPersistentConnection](f903dae0-4a4e-c7db-8350-50e36b257ec5.htm) | **已过时。** V11版本及之前设置长连接的方法，在V12版本以上中没有任何效果，默认长连接，删除调用即可，此处保留方法是为了部分用户保持兼容性升级。  The method of setting the long connection in V11 and before, has no effect in V12 and above. this method can be deleted. The method is retained here to maintain compatibility upgrades for some users. (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共方法 | [SetSenderAMSNetId](6c085271-688a-909e-53c9-dec86b2a8e8a.htm) | 设置原目标地址 举例 192.168.0.100.1.1；也可以是带端口号 192.168.0.100.1.1:34567  Set the original destination address Example: 192.168.0.100.1.1; it can also be the port number 192.168.0.100.1.1: 34567 |
| 公共方法 | [SetTargetAMSNetId](af74dcd7-fb14-527d-5b79-881112e412ad.htm) | 目标的地址，举例 192.168.0.1.1.1；也可以是带端口号 192.168.0.1.1.1:801  The address of the destination, for example 192.168.0.1.1.1; it can also be the port number 192.168.0.1.1.1: 801 |
| 公共方法 | [ToString](5cf941e3-75eb-cfcb-ec45-0065be0058c5.htm) | (重写 [DeviceTcpNetToString](209a196b-90ea-2b73-e915-ce4b11f1263d.htm).) |
| 公共方法 | [TransValueHandle](f7d19175-e124-9a8d-d36e-638b14d9c919.htm) | 将字符串的地址转换为内存的地址，其他地址则不操作  Converts the address of a string to the address of a memory, other addresses do not operate |
| 公共方法 | [TransValueHandleAsync](efa735b1-7770-c061-0e6f-0486b84ff6b8.htm) | 将字符串的地址转换为内存的地址，其他地址则不操作  Converts the address of a string to the address of a memory, other addresses do not operate |
| 公共方法 | [UnpackResponseContent](9578724a-11cf-6cec-c183-18fb69db6412.htm) | 根据对方返回的报文命令，对命令进行基本的拆包，例如各种Modbus协议拆包为统一的核心报文，还支持对报文的验证  According to the message command returned by the other party, the command is basically unpacked, for example, various Modbus protocols are unpacked into a unified core message, and the verification of the message is also supported (重写 [BinaryCommunicationUnpackResponseContent(Byte, Byte)](208017a7-8da2-33f8-3273-d7e7e850100f.htm).) |
| 公共方法 | [Wait(String, Boolean, Int32, Int32)](0e4db26e-1177-2b18-9492-ba05afb6fcc4.htm) | 等待指定地址的Boolean值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the Boolean value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [Wait(String, Int16, Int32, Int32)](e06fe098-8c8b-83f9-73ca-2a7cce19d43d.htm) | 等待指定地址的Int16值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the Int16 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [Wait(String, Int32, Int32, Int32)](1bce0aae-27e9-f1d6-e8f3-25ae6b47554e.htm) | 等待指定地址的Int32值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the Int32 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [Wait(String, Int64, Int32, Int32)](073354e9-2508-2158-1b01-881c3a34844b.htm) | 等待指定地址的Int64值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the Int64 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [Wait(String, UInt16, Int32, Int32)](d97f8d26-1c1a-df69-93f0-f60209597ced.htm) | 等待指定地址的UInt16值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the UInt16 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [Wait(String, UInt32, Int32, Int32)](1e61c395-b07f-e66c-5f9c-2708093150d9.htm) | 等待指定地址的UInt32值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the UInt32 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [Wait(String, UInt64, Int32, Int32)](d659729b-5766-c497-d418-856c34ec2b24.htm) | 等待指定地址的UInt64值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the UInt64 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WaitAsync(String, Boolean, Int32, Int32)](65399e87-9d38-f819-ccd3-005f68c2c3b3.htm) | 等待指定地址的Boolean值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the Boolean value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WaitAsync(String, Int16, Int32, Int32)](01277e39-afb1-86a4-3aa3-560160e66fb9.htm) | 等待指定地址的Int16值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the Int16 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WaitAsync(String, Int32, Int32, Int32)](4a6d380f-6b2a-1270-07a7-39b88a2712fd.htm) | 等待指定地址的Int32值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the Int32 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WaitAsync(String, Int64, Int32, Int32)](d7eede9b-4d44-479b-fb88-3ff2e1bb7568.htm) | 等待指定地址的Int64值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the Int64 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WaitAsync(String, UInt16, Int32, Int32)](0a8395c7-dc7f-4d5f-4316-98e18fe70f4c.htm) | 等待指定地址的UInt16值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the UInt16 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WaitAsync(String, UInt32, Int32, Int32)](19e9074e-99a2-99fd-91ef-a5314a115397.htm) | 等待指定地址的UInt32值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the UInt32 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WaitAsync(String, UInt64, Int32, Int32)](5109dac0-0ba2-cbb9-e383-1a869c6dc22e.htm) | 等待指定地址的UInt64值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the UInt64 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [Write(String, Boolean)](06d0b9af-7f28-0b9c-098f-4da28e69504a.htm) | 写入单个的Boolean数据，返回是否成功  Write a single Boolean data, and return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Double)](64c44576-7e95-cbcb-a8c9-e6b7ff7e3c50.htm) | 写入double数据，返回是否成功  Write double data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Double)](df56e6b1-0688-ece7-3bf1-9db32fec9f6e.htm) | 写入double数组，返回是否成功  Write double array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Int16)](0e23506e-b604-e67c-b5eb-06271b6c6335.htm) | 写入short数据，返回是否成功  Write short data, returns whether success (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Int16)](c7cfd67d-ca31-3087-c684-ee4c10780ab1.htm) | 写入short数组，返回是否成功  Write short array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Int32)](ba80d1a5-c8d8-7e0b-5c6b-9435cf20cadb.htm) | 写入int数据，返回是否成功  Write int data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Int32)](57dac21c-8b28-de7c-c213-77cda465cccb.htm) | 写入int[]数组，返回是否成功  Write int array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Int64)](ba246b33-5a12-a22c-6d40-0969a64b508b.htm) | 写入long数据，返回是否成功  Write long data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Int64)](5b910f9c-6e6b-bb9d-6842-1b496640c1f4.htm) | 写入long数组，返回是否成功  Write long array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Single)](a2265fa7-5b56-a5aa-9d21-a0dc5458cbf9.htm) | 写入float数据，返回是否成功  Write float data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Single)](bb76227e-3e95-8ca2-5dcf-c087d90b8a01.htm) | 写入float数组，返回是否成功  Write float array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, String)](4e0ec363-ba64-7ca7-b8d3-2e9a8087bd55.htm) | 写入字符串信息，编码为ASCII  Write string information, encoded as ASCII (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, UInt16)](e5bc9466-65b4-fc64-e0d0-1867b1113130.htm) | 写入ushort数据，返回是否成功  Write ushort data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, UInt16)](65544962-82ac-bb47-d72f-5dd0c3d4ac6c.htm) | 写入ushort数组，返回是否成功  Write ushort array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, UInt32)](c99d8246-f95e-4f79-738b-16e622072eb4.htm) | 写入uint数据，返回是否成功  Write uint data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, UInt32)](b80b3099-b34d-fde8-6a65-837da3bdc5ab.htm) | 写入uint[]数组，返回是否成功  Write uint array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, UInt64)](4924086e-9b36-d446-4ae4-6682134a6ad1.htm) | 写入ulong数据，返回是否成功  Write ulong data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, UInt64)](b1131bdd-8b51-ec6d-dd0e-fec1f6cb51b6.htm) | 写入ulong数组，返回是否成功  Write ulong array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [Write(String, Boolean)](0a254c9a-8a5a-3bea-e82f-4ba1e13ceb9a.htm) | 写入PLC的数据，地址共有三种格式，一：I,Q,M数据信息，举例M0,M100；二：内存地址，i=100000；三：标签地址，s=A  There are three formats for the data written into the PLC. One: I, Q, M data information, such as M0, M100; two: memory address, i = 100000; three: tag address, s = A (重写 [DeviceCommunicationWrite(String, Boolean)](5b046c82-3d8e-550b-fdf3-41ca5558169b.htm).) |
| 公共方法 | [Write(String, Byte)](40483d77-6945-118f-6b87-b018f6e6bc3c.htm) | 写入PLC的数据，地址共有三种格式，一：I,Q,M数据信息，举例M0,M100；二：内存地址，i=100000；三：标签地址，s=A  There are three formats for the data written into the PLC. One: I, Q, M data information, such as M0, M100; two: memory address, i = 100000; three: tag address, s = A |
| 公共方法 | [Write(String, Byte)](39cd682e-4c5f-2ff7-acb2-c4771ab24081.htm) | 写入PLC的数据，地址共有三种格式，一：I,Q,M数据信息，举例M0,M100；二：内存地址，i=100000；三：标签地址，s=A  There are three formats for the data written into the PLC. One: I, Q, M data information, such as M0, M100; two: memory address, i = 100000; three: tag address, s = A (重写 [DeviceCommunicationWrite(String, Byte)](b5536dd6-ccde-6883-ac27-bff9e4b37806.htm).) |
| 公共方法代码示例 | [Write(String, String, Int32)](02fe2ca9-49a3-ede4-d1af-8c6da1f644da.htm) | 写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为ASCII  Write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, it is filled with 0 and the encoding is ASCII. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, String, Encoding)](7b602661-902d-2b70-1fdb-ba87ecd307bf.htm) | 写入字符串信息，需要指定的编码信息  Write string information, need to specify the encoding information (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, String, Int32, Encoding)](e8cb30d9-4c5e-6c68-b9af-1b1b6c2dc714.htm) | 写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为指定的编码信息  Write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, then the operation is complemented with 0 , you should specified the encoding information (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteT(T)](9f48853f-fbeb-c37d-25a2-1c76d1dee407.htm) | 写入支持Hsl特性的数据，返回是否写入成功，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Write data that supports the Hsl attribute, and return whether the write was successful. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WriteAdsState](c05364e4-eeb6-6b26-e817-0413145e4e8b.htm) | 写入Ads的状态，可以携带数据信息，数据可以为空  Write the status of Ads, can carry data information, and the data can be empty |
| 公共方法 | [WriteAdsStateAsync](7a6c4cf6-7837-0873-3d9d-e6b99acec800.htm) | 写入Ads的状态，可以携带数据信息，数据可以为空  Write the status of Ads, can carry data information, and the data can be empty |
| 公共方法 | [WriteAsync(String, Boolean)](10fe669c-805e-1e5a-f016-f319ad363e90.htm) | 异步批量写入Boolean数组数据，返回是否成功  Asynchronously batch write Boolean array data, return success (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Double)](068adaa4-d1de-dffc-8376-6f483e2c1d96.htm) | 异步写入double数据，返回是否成功  Asynchronously write double data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Double)](36199b83-c4ea-b02f-3675-fb6ac48d0c59.htm) | 异步写入double数组，返回是否成功  Asynchronously write double array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int16)](ccd95ccc-10a3-af67-62e8-c2312b4871b0.htm) | 异步写入short数据，返回是否成功  Asynchronously write short data, returns whether success (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int16)](2625e868-7096-e526-ce25-b29afc639b41.htm) | 异步写入short数组，返回是否成功  Asynchronously write short array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int32)](df91d629-2080-328d-67cd-fbefb9a45404.htm) | 异步写入int数据，返回是否成功  Asynchronously write int data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int32)](99962a08-9624-cea6-6009-5e6f207236b7.htm) | 异步写入int[]数组，返回是否成功  Asynchronously write int array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int64)](ec6c1e24-a0d8-f1f5-56af-109b0c546fc8.htm) | 异步写入long数据，返回是否成功  Asynchronously write long data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int64)](4b2ab94c-2891-e0a9-f10b-b67c63bd1cf9.htm) | 异步写入long数组，返回是否成功  Asynchronously write long array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Single)](910189df-13f9-b7d0-a626-b5b25ce5d32a.htm) | 异步写入float数据，返回是否成功  Asynchronously write float data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Single)](dbb2359a-fb54-cf93-4413-20d5a4ef9d36.htm) | 异步写入float数组，返回是否成功  Asynchronously write float array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, String)](fe839c9a-46f3-e317-0155-b10a8a83a9b1.htm) | 异步写入字符串信息，编码为ASCII  Asynchronously write string information, encoded as ASCII (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt16)](a04f3aa2-a51e-26cd-974e-bf70a2db359d.htm) | 异步写入ushort数据，返回是否成功  Asynchronously write ushort data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt16)](34a1a321-972e-06a8-76fb-794dd1229005.htm) | 异步写入ushort数组，返回是否成功  Asynchronously write ushort array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt32)](abad215c-bc99-ac90-b9f3-7adb3a7cca8f.htm) | 异步写入uint数据，返回是否成功  Asynchronously write uint data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt32)](5b655bf3-5ba7-ad31-6897-89c59051ac60.htm) | 异步写入uint[]数组，返回是否成功  Asynchronously write uint array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt64)](e396074d-6cf1-b45d-5e8e-01026f6e391b.htm) | 异步写入ulong数据，返回是否成功  Asynchronously write ulong data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt64)](ae3df8ec-29c0-42ff-3d6a-68f496997b3d.htm) | 异步写入ulong数组，返回是否成功  Asynchronously write ulong array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WriteAsync(String, Boolean)](c97f4f59-2c83-022b-9612-d004221cfc2a.htm) | 写入PLC的数据，地址共有三种格式，一：I,Q,M数据信息，举例M0,M100；二：内存地址，i=100000；三：标签地址，s=A  There are three formats for the data written into the PLC. One: I, Q, M data information, such as M0, M100; two: memory address, i = 100000; three: tag address, s = A (重写 [DeviceCommunicationWriteAsync(String, Boolean)](fe872dae-7ef9-12c1-d883-e92ba8676416.htm).) |
| 公共方法 | [WriteAsync(String, Byte)](632f4ccf-893d-5312-15c1-b4653e077f3f.htm) | 写入PLC的数据，地址共有三种格式，一：I,Q,M数据信息，举例M0,M100；二：内存地址，i=100000；三：标签地址，s=A  There are three formats for the data written into the PLC. One: I, Q, M data information, such as M0, M100; two: memory address, i = 100000; three: tag address, s = A |
| 公共方法 | [WriteAsync(String, Byte)](aa8d4d98-cb38-abdf-2c6a-c1e0527d5c15.htm) | 写入PLC的数据，地址共有三种格式，一：I,Q,M数据信息，举例M0,M100；二：内存地址，i=100000；三：标签地址，s=A  There are three formats for the data written into the PLC. One: I, Q, M data information, such as M0, M100; two: memory address, i = 100000; three: tag address, s = A (重写 [DeviceCommunicationWriteAsync(String, Byte)](8f1cec30-8b86-6749-c8f3-7d39168d59df.htm).) |
| 公共方法 | [WriteAsync(String, String, Int32)](b739b983-b03c-8ecf-2713-75f19233ae2f.htm) | 异步写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为ASCII  Asynchronously write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, it is filled with 0 and the encoding is ASCII. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, String, Encoding)](876e9800-0406-6407-bdf8-8ab1cf28df38.htm) | 异步写入字符串信息，需要指定的编码信息  Asynchronously write string information, need to specify the encoding information (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WriteAsync(String, String, Int32, Encoding)](ae9fb2a3-fddb-48eb-e9ad-6cc0ebbac1f6.htm) | 异步写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为指定的编码信息  Asynchronously write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, then the operation is complemented with 0 , you should specified the encoding information (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsyncT(T)](25d96e77-d1b2-918c-33f9-e5e4cc57c9cd.htm) | 异步写入支持Hsl特性的数据，返回是否写入成功，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Asynchronously write data that supports the Hsl attribute, and return whether the write was successful. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteCustomerT](84a37f82-3f67-a0f8-fd3f-0630772565e4.htm) | 写入自定义类型的数据，该类型必须继承自IDataTransfer接口  Write data of a custom type, which must inherit from the IDataTransfer interface (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteCustomerAsyncT](04628397-cfa1-d785-7790-4c02cb377ab6.htm) | 写入自定义类型的数据，该类型必须继承自IDataTransfer接口  Write data of a custom type, which must inherit from the IDataTransfer interface (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WriteStructT](a7aeee39-ab11-376a-71b5-02a066630636.htm) | 将一个结构体写入到指定的地址中去，需要指定写入的起始地址  To write a structure to a specified address, you need to specify the start address of the write |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[BeckhoffAdsNet 类](7ab15712-8688-063d-6b30-4b24803393fd.htm)

[HslCommunication.Profinet.Beckhoff 命名空间](fb2af334-3662-b9ca-be3c-b0090bf4ef56.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ClearTagsCache 方法 

[原文連結](http://api.hslcommunication.cn/html/3ca05c73-2ab5-b7b8-2d46-03da8a9c67ad.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Beckhoff](../html/fb2af334-3662-b9ca-be3c-b0090bf4ef56.htm "HslCommunication.Profinet.Beckhoff")

[BeckhoffAdsNet 类](../html/7ab15712-8688-063d-6b30-4b24803393fd.htm "BeckhoffAdsNet 类")

[BeckhoffAdsNet 方法](../html/e8bad1d3-d4ae-37ef-6ba7-e199dac65f0c.htm "BeckhoffAdsNet 方法")

[ClearTagsCache 方法](../html/3ca05c73-2ab5-b7b8-2d46-03da8a9c67ad.htm "ClearTagsCache 方法 ")

[DecideWhetherQAMessage 方法](../html/771586fe-9184-5312-e3da-bae61fe91d9a.htm "DecideWhetherQAMessage 方法 ")

[ExtraAfterReadFromCoreServer 方法](../html/3d26d83d-9287-54f5-3d35-c6d074ce6263.htm "ExtraAfterReadFromCoreServer 方法 ")

[GetNewNetMessage 方法](../html/215ca52b-b151-9cb6-0443-00b9be103188.htm "GetNewNetMessage 方法 ")

[GetSenderAMSNetId 方法](../html/f743d65e-7896-7d75-ba57-6be4991ab9af.htm "GetSenderAMSNetId 方法 ")

[GetTargetAMSNetId 方法](../html/dcc58722-6399-8dd4-affa-2e531a098ae0.htm "GetTargetAMSNetId 方法 ")

[InitializationOnConnect 方法](../html/ea40c0b0-f579-3a64-7b45-573259ef706d.htm "InitializationOnConnect 方法 ")

[InitializationOnConnectAsync 方法](../html/c2a3f0f3-3a78-e102-d07f-d15053af2633.htm "InitializationOnConnectAsync 方法 ")

[PackCommandWithHeader 方法](../html/3895f6bc-bf81-5dbe-cb27-c12c6602126d.htm "PackCommandWithHeader 方法 ")

[Read 方法](../html/504e0a30-6e4b-ac79-e7e4-ce3aecedecb9.htm "Read 方法 ")

[ReadAdsDeviceInfo 方法](../html/1b58e4ea-7a7c-3b42-1f11-b077010b9a39.htm "ReadAdsDeviceInfo 方法 ")

[ReadAdsDeviceInfoAsync 方法](../html/18ef1d19-2639-b4b9-99c6-3e7caed8ef5a.htm "ReadAdsDeviceInfoAsync 方法 ")

[ReadAdsState 方法](../html/4645af4b-15b0-49a7-a3d3-bc07ef913ad7.htm "ReadAdsState 方法 ")

[ReadAdsStateAsync 方法](../html/831eac3f-b313-b2bd-3131-71bf3948b652.htm "ReadAdsStateAsync 方法 ")

[ReadAsync 方法](../html/ffaa1904-b9ab-79ff-efde-2543c610aa26.htm "ReadAsync 方法 ")

[ReadBool 方法](../html/b4ab6787-c834-7bf7-e069-1e8f7d5a4cad.htm "ReadBool 方法 ")

[ReadBoolAsync 方法](../html/1114ae0d-0359-668d-a6b2-6be76e7da28b.htm "ReadBoolAsync 方法 ")

[ReadByte 方法](../html/dadb0178-8d9c-2ce1-5495-e53d72771b82.htm "ReadByte 方法 ")

[ReadByteAsync 方法](../html/d61db681-dca5-f888-2dfd-c35faa269a4f.htm "ReadByteAsync 方法 ")

[ReadStruct 方法](../html/8e09ddd7-faa8-67c9-e099-eee24d59aa8d.htm "ReadStruct 方法 ")

[ReadValueHandle 方法](../html/7a0c126f-5f3e-f077-f90b-34ca2191c5b2.htm "ReadValueHandle 方法 ")

[ReadValueHandleAsync 方法](../html/0dadc9a4-90cc-1f29-627f-f8112aaa7df2.htm "ReadValueHandleAsync 方法 ")

[ReleaseSystemHandle 方法](../html/e2fa36f0-6d1f-2e02-e3de-bc2344e8d553.htm "ReleaseSystemHandle 方法 ")

[ReleaseSystemHandleAsync 方法](../html/cb4896ea-bbd4-52fa-4587-db7c137ff2d6.htm "ReleaseSystemHandleAsync 方法 ")

[SetSenderAMSNetId 方法](../html/6c085271-688a-909e-53c9-dec86b2a8e8a.htm "SetSenderAMSNetId 方法 ")

[SetTargetAMSNetId 方法](../html/af74dcd7-fb14-527d-5b79-881112e412ad.htm "SetTargetAMSNetId 方法 ")

[ToString 方法](../html/5cf941e3-75eb-cfcb-ec45-0065be0058c5.htm "ToString 方法 ")

[TransValueHandle 方法](../html/f7d19175-e124-9a8d-d36e-638b14d9c919.htm "TransValueHandle 方法 ")

[TransValueHandleAsync 方法](../html/efa735b1-7770-c061-0e6f-0486b84ff6b8.htm "TransValueHandleAsync 方法 ")

[UnpackResponseContent 方法](../html/9578724a-11cf-6cec-c183-18fb69db6412.htm "UnpackResponseContent 方法 ")

[Write 方法](../html/d6d8c193-854e-1069-a9f0-d426a76813cc.htm "Write 方法 ")

[WriteAdsState 方法](../html/c05364e4-eeb6-6b26-e817-0413145e4e8b.htm "WriteAdsState 方法 ")

[WriteAdsStateAsync 方法](../html/7a6c4cf6-7837-0873-3d9d-e6b99acec800.htm "WriteAdsStateAsync 方法 ")

[WriteAsync 方法](../html/7f760de5-65c9-46e0-b53e-fa719a52dafc.htm "WriteAsync 方法 ")

[WriteStruct(T) 方法](../html/a7aeee39-ab11-376a-71b5-02a066630636.htm "WriteStruct(T) 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| BeckhoffAdsNetClearTagsCache 方法 |

清除标签缓存信息  
Clear the tag cache information

**命名空间：**
 [HslCommunication.Profinet.Beckhoff](fb2af334-3662-b9ca-be3c-b0090bf4ef56.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public void ClearTagsCache()
```

```
Public Sub ClearTagsCache
```

```
public:
void ClearTagsCache()
```

```
member ClearTagsCache : unit -> unit 
```

![](../icons/SectionExpanded.png)参见

#### 引用

[BeckhoffAdsNet 类](7ab15712-8688-063d-6b30-4b24803393fd.htm)

[HslCommunication.Profinet.Beckhoff 命名空间](fb2af334-3662-b9ca-be3c-b0090bf4ef56.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## DecideWhetherQAMessage 方法 

[原文連結](http://api.hslcommunication.cn/html/771586fe-9184-5312-e3da-bae61fe91d9a.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Beckhoff](../html/fb2af334-3662-b9ca-be3c-b0090bf4ef56.htm "HslCommunication.Profinet.Beckhoff")

[BeckhoffAdsNet 类](../html/7ab15712-8688-063d-6b30-4b24803393fd.htm "BeckhoffAdsNet 类")

[BeckhoffAdsNet 方法](../html/e8bad1d3-d4ae-37ef-6ba7-e199dac65f0c.htm "BeckhoffAdsNet 方法")

[ClearTagsCache 方法](../html/3ca05c73-2ab5-b7b8-2d46-03da8a9c67ad.htm "ClearTagsCache 方法 ")

[DecideWhetherQAMessage 方法](../html/771586fe-9184-5312-e3da-bae61fe91d9a.htm "DecideWhetherQAMessage 方法 ")

[ExtraAfterReadFromCoreServer 方法](../html/3d26d83d-9287-54f5-3d35-c6d074ce6263.htm "ExtraAfterReadFromCoreServer 方法 ")

[GetNewNetMessage 方法](../html/215ca52b-b151-9cb6-0443-00b9be103188.htm "GetNewNetMessage 方法 ")

[GetSenderAMSNetId 方法](../html/f743d65e-7896-7d75-ba57-6be4991ab9af.htm "GetSenderAMSNetId 方法 ")

[GetTargetAMSNetId 方法](../html/dcc58722-6399-8dd4-affa-2e531a098ae0.htm "GetTargetAMSNetId 方法 ")

[InitializationOnConnect 方法](../html/ea40c0b0-f579-3a64-7b45-573259ef706d.htm "InitializationOnConnect 方法 ")

[InitializationOnConnectAsync 方法](../html/c2a3f0f3-3a78-e102-d07f-d15053af2633.htm "InitializationOnConnectAsync 方法 ")

[PackCommandWithHeader 方法](../html/3895f6bc-bf81-5dbe-cb27-c12c6602126d.htm "PackCommandWithHeader 方法 ")

[Read 方法](../html/504e0a30-6e4b-ac79-e7e4-ce3aecedecb9.htm "Read 方法 ")

[ReadAdsDeviceInfo 方法](../html/1b58e4ea-7a7c-3b42-1f11-b077010b9a39.htm "ReadAdsDeviceInfo 方法 ")

[ReadAdsDeviceInfoAsync 方法](../html/18ef1d19-2639-b4b9-99c6-3e7caed8ef5a.htm "ReadAdsDeviceInfoAsync 方法 ")

[ReadAdsState 方法](../html/4645af4b-15b0-49a7-a3d3-bc07ef913ad7.htm "ReadAdsState 方法 ")

[ReadAdsStateAsync 方法](../html/831eac3f-b313-b2bd-3131-71bf3948b652.htm "ReadAdsStateAsync 方法 ")

[ReadAsync 方法](../html/ffaa1904-b9ab-79ff-efde-2543c610aa26.htm "ReadAsync 方法 ")

[ReadBool 方法](../html/b4ab6787-c834-7bf7-e069-1e8f7d5a4cad.htm "ReadBool 方法 ")

[ReadBoolAsync 方法](../html/1114ae0d-0359-668d-a6b2-6be76e7da28b.htm "ReadBoolAsync 方法 ")

[ReadByte 方法](../html/dadb0178-8d9c-2ce1-5495-e53d72771b82.htm "ReadByte 方法 ")

[ReadByteAsync 方法](../html/d61db681-dca5-f888-2dfd-c35faa269a4f.htm "ReadByteAsync 方法 ")

[ReadStruct 方法](../html/8e09ddd7-faa8-67c9-e099-eee24d59aa8d.htm "ReadStruct 方法 ")

[ReadValueHandle 方法](../html/7a0c126f-5f3e-f077-f90b-34ca2191c5b2.htm "ReadValueHandle 方法 ")

[ReadValueHandleAsync 方法](../html/0dadc9a4-90cc-1f29-627f-f8112aaa7df2.htm "ReadValueHandleAsync 方法 ")

[ReleaseSystemHandle 方法](../html/e2fa36f0-6d1f-2e02-e3de-bc2344e8d553.htm "ReleaseSystemHandle 方法 ")

[ReleaseSystemHandleAsync 方法](../html/cb4896ea-bbd4-52fa-4587-db7c137ff2d6.htm "ReleaseSystemHandleAsync 方法 ")

[SetSenderAMSNetId 方法](../html/6c085271-688a-909e-53c9-dec86b2a8e8a.htm "SetSenderAMSNetId 方法 ")

[SetTargetAMSNetId 方法](../html/af74dcd7-fb14-527d-5b79-881112e412ad.htm "SetTargetAMSNetId 方法 ")

[ToString 方法](../html/5cf941e3-75eb-cfcb-ec45-0065be0058c5.htm "ToString 方法 ")

[TransValueHandle 方法](../html/f7d19175-e124-9a8d-d36e-638b14d9c919.htm "TransValueHandle 方法 ")

[TransValueHandleAsync 方法](../html/efa735b1-7770-c061-0e6f-0486b84ff6b8.htm "TransValueHandleAsync 方法 ")

[UnpackResponseContent 方法](../html/9578724a-11cf-6cec-c183-18fb69db6412.htm "UnpackResponseContent 方法 ")

[Write 方法](../html/d6d8c193-854e-1069-a9f0-d426a76813cc.htm "Write 方法 ")

[WriteAdsState 方法](../html/c05364e4-eeb6-6b26-e817-0413145e4e8b.htm "WriteAdsState 方法 ")

[WriteAdsStateAsync 方法](../html/7a6c4cf6-7837-0873-3d9d-e6b99acec800.htm "WriteAdsStateAsync 方法 ")

[WriteAsync 方法](../html/7f760de5-65c9-46e0-b53e-fa719a52dafc.htm "WriteAsync 方法 ")

[WriteStruct(T) 方法](../html/a7aeee39-ab11-376a-71b5-02a066630636.htm "WriteStruct(T) 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| BeckhoffAdsNetDecideWhetherQAMessage 方法 |

决定当前的消息是否是用于问答机制返回的消息，默认直接返回 true, 实际的情况需要根据协议进行重写方法  
To determine whether the current message is the message returned by the question answering mechanism,
the default is true. In actual cases, the rewriting method needs to be performed according to the protocol

**命名空间：**
 [HslCommunication.Profinet.Beckhoff](fb2af334-3662-b9ca-be3c-b0090bf4ef56.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
protected override bool DecideWhetherQAMessage(
	CommunicationPipe pipe,
	OperateResult<byte[]> receive
)
```

```
Protected Overrides Function DecideWhetherQAMessage ( 
	pipe As CommunicationPipe,
	receive As OperateResult(Of Byte())
) As Boolean
```

```
protected:
virtual bool DecideWhetherQAMessage(
	CommunicationPipe^ pipe, 
	OperateResult<array<unsigned char>^>^ receive
) override
```

```
abstract DecideWhetherQAMessage : 
        pipe : CommunicationPipe * 
        receive : OperateResult<byte[]> -> bool 
override DecideWhetherQAMessage : 
        pipe : CommunicationPipe * 
        receive : OperateResult<byte[]> -> bool
```

#### 参数

pipe
:   类型：[HslCommunication.Core.PipeCommunicationPipe](ab3f9c5c-144b-4209-4999-b382a1d94450.htm)  
    管道信息

receive
:   类型：[HslCommunicationOperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
    接收的数据信息

#### 返回值

类型：Boolean  
是否是问答的数据

![](../icons/SectionExpanded.png)参见

#### 引用

[BeckhoffAdsNet 类](7ab15712-8688-063d-6b30-4b24803393fd.htm)

[HslCommunication.Profinet.Beckhoff 命名空间](fb2af334-3662-b9ca-be3c-b0090bf4ef56.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ExtraAfterReadFromCoreServer 方法 

[原文連結](http://api.hslcommunication.cn/html/3d26d83d-9287-54f5-3d35-c6d074ce6263.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Beckhoff](../html/fb2af334-3662-b9ca-be3c-b0090bf4ef56.htm "HslCommunication.Profinet.Beckhoff")

[BeckhoffAdsNet 类](../html/7ab15712-8688-063d-6b30-4b24803393fd.htm "BeckhoffAdsNet 类")

[BeckhoffAdsNet 方法](../html/e8bad1d3-d4ae-37ef-6ba7-e199dac65f0c.htm "BeckhoffAdsNet 方法")

[ClearTagsCache 方法](../html/3ca05c73-2ab5-b7b8-2d46-03da8a9c67ad.htm "ClearTagsCache 方法 ")

[DecideWhetherQAMessage 方法](../html/771586fe-9184-5312-e3da-bae61fe91d9a.htm "DecideWhetherQAMessage 方法 ")

[ExtraAfterReadFromCoreServer 方法](../html/3d26d83d-9287-54f5-3d35-c6d074ce6263.htm "ExtraAfterReadFromCoreServer 方法 ")

[GetNewNetMessage 方法](../html/215ca52b-b151-9cb6-0443-00b9be103188.htm "GetNewNetMessage 方法 ")

[GetSenderAMSNetId 方法](../html/f743d65e-7896-7d75-ba57-6be4991ab9af.htm "GetSenderAMSNetId 方法 ")

[GetTargetAMSNetId 方法](../html/dcc58722-6399-8dd4-affa-2e531a098ae0.htm "GetTargetAMSNetId 方法 ")

[InitializationOnConnect 方法](../html/ea40c0b0-f579-3a64-7b45-573259ef706d.htm "InitializationOnConnect 方法 ")

[InitializationOnConnectAsync 方法](../html/c2a3f0f3-3a78-e102-d07f-d15053af2633.htm "InitializationOnConnectAsync 方法 ")

[PackCommandWithHeader 方法](../html/3895f6bc-bf81-5dbe-cb27-c12c6602126d.htm "PackCommandWithHeader 方法 ")

[Read 方法](../html/504e0a30-6e4b-ac79-e7e4-ce3aecedecb9.htm "Read 方法 ")

[ReadAdsDeviceInfo 方法](../html/1b58e4ea-7a7c-3b42-1f11-b077010b9a39.htm "ReadAdsDeviceInfo 方法 ")

[ReadAdsDeviceInfoAsync 方法](../html/18ef1d19-2639-b4b9-99c6-3e7caed8ef5a.htm "ReadAdsDeviceInfoAsync 方法 ")

[ReadAdsState 方法](../html/4645af4b-15b0-49a7-a3d3-bc07ef913ad7.htm "ReadAdsState 方法 ")

[ReadAdsStateAsync 方法](../html/831eac3f-b313-b2bd-3131-71bf3948b652.htm "ReadAdsStateAsync 方法 ")

[ReadAsync 方法](../html/ffaa1904-b9ab-79ff-efde-2543c610aa26.htm "ReadAsync 方法 ")

[ReadBool 方法](../html/b4ab6787-c834-7bf7-e069-1e8f7d5a4cad.htm "ReadBool 方法 ")

[ReadBoolAsync 方法](../html/1114ae0d-0359-668d-a6b2-6be76e7da28b.htm "ReadBoolAsync 方法 ")

[ReadByte 方法](../html/dadb0178-8d9c-2ce1-5495-e53d72771b82.htm "ReadByte 方法 ")

[ReadByteAsync 方法](../html/d61db681-dca5-f888-2dfd-c35faa269a4f.htm "ReadByteAsync 方法 ")

[ReadStruct 方法](../html/8e09ddd7-faa8-67c9-e099-eee24d59aa8d.htm "ReadStruct 方法 ")

[ReadValueHandle 方法](../html/7a0c126f-5f3e-f077-f90b-34ca2191c5b2.htm "ReadValueHandle 方法 ")

[ReadValueHandleAsync 方法](../html/0dadc9a4-90cc-1f29-627f-f8112aaa7df2.htm "ReadValueHandleAsync 方法 ")

[ReleaseSystemHandle 方法](../html/e2fa36f0-6d1f-2e02-e3de-bc2344e8d553.htm "ReleaseSystemHandle 方法 ")

[ReleaseSystemHandleAsync 方法](../html/cb4896ea-bbd4-52fa-4587-db7c137ff2d6.htm "ReleaseSystemHandleAsync 方法 ")

[SetSenderAMSNetId 方法](../html/6c085271-688a-909e-53c9-dec86b2a8e8a.htm "SetSenderAMSNetId 方法 ")

[SetTargetAMSNetId 方法](../html/af74dcd7-fb14-527d-5b79-881112e412ad.htm "SetTargetAMSNetId 方法 ")

[ToString 方法](../html/5cf941e3-75eb-cfcb-ec45-0065be0058c5.htm "ToString 方法 ")

[TransValueHandle 方法](../html/f7d19175-e124-9a8d-d36e-638b14d9c919.htm "TransValueHandle 方法 ")

[TransValueHandleAsync 方法](../html/efa735b1-7770-c061-0e6f-0486b84ff6b8.htm "TransValueHandleAsync 方法 ")

[UnpackResponseContent 方法](../html/9578724a-11cf-6cec-c183-18fb69db6412.htm "UnpackResponseContent 方法 ")

[Write 方法](../html/d6d8c193-854e-1069-a9f0-d426a76813cc.htm "Write 方法 ")

[WriteAdsState 方法](../html/c05364e4-eeb6-6b26-e817-0413145e4e8b.htm "WriteAdsState 方法 ")

[WriteAdsStateAsync 方法](../html/7a6c4cf6-7837-0873-3d9d-e6b99acec800.htm "WriteAdsStateAsync 方法 ")

[WriteAsync 方法](../html/7f760de5-65c9-46e0-b53e-fa719a52dafc.htm "WriteAsync 方法 ")

[WriteStruct(T) 方法](../html/a7aeee39-ab11-376a-71b5-02a066630636.htm "WriteStruct(T) 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| BeckhoffAdsNetExtraAfterReadFromCoreServer 方法 |

和服务器交互完成的时候调用的方法，可以根据读写结果进行一些额外的操作，具体的操作需要根据实际的需求来重写实现  
The method called when the interaction with the server is completed can perform some additional operations based on the read and write results.
The specific operations need to be rewritten according to actual needs.

**命名空间：**
 [HslCommunication.Profinet.Beckhoff](fb2af334-3662-b9ca-be3c-b0090bf4ef56.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
protected override void ExtraAfterReadFromCoreServer(
	OperateResult read
)
```

```
Protected Overrides Sub ExtraAfterReadFromCoreServer ( 
	read As OperateResult
)
```

```
protected:
virtual void ExtraAfterReadFromCoreServer(
	OperateResult^ read
) override
```

```
abstract ExtraAfterReadFromCoreServer : 
        read : OperateResult -> unit 
override ExtraAfterReadFromCoreServer : 
        read : OperateResult -> unit
```

#### 参数

read
:   类型：[HslCommunicationOperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
    读取结果

![](../icons/SectionExpanded.png)参见

#### 引用

[BeckhoffAdsNet 类](7ab15712-8688-063d-6b30-4b24803393fd.htm)

[HslCommunication.Profinet.Beckhoff 命名空间](fb2af334-3662-b9ca-be3c-b0090bf4ef56.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GetNewNetMessage 方法 

[原文連結](http://api.hslcommunication.cn/html/215ca52b-b151-9cb6-0443-00b9be103188.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Beckhoff](../html/fb2af334-3662-b9ca-be3c-b0090bf4ef56.htm "HslCommunication.Profinet.Beckhoff")

[BeckhoffAdsNet 类](../html/7ab15712-8688-063d-6b30-4b24803393fd.htm "BeckhoffAdsNet 类")

[BeckhoffAdsNet 方法](../html/e8bad1d3-d4ae-37ef-6ba7-e199dac65f0c.htm "BeckhoffAdsNet 方法")

[ClearTagsCache 方法](../html/3ca05c73-2ab5-b7b8-2d46-03da8a9c67ad.htm "ClearTagsCache 方法 ")

[DecideWhetherQAMessage 方法](../html/771586fe-9184-5312-e3da-bae61fe91d9a.htm "DecideWhetherQAMessage 方法 ")

[ExtraAfterReadFromCoreServer 方法](../html/3d26d83d-9287-54f5-3d35-c6d074ce6263.htm "ExtraAfterReadFromCoreServer 方法 ")

[GetNewNetMessage 方法](../html/215ca52b-b151-9cb6-0443-00b9be103188.htm "GetNewNetMessage 方法 ")

[GetSenderAMSNetId 方法](../html/f743d65e-7896-7d75-ba57-6be4991ab9af.htm "GetSenderAMSNetId 方法 ")

[GetTargetAMSNetId 方法](../html/dcc58722-6399-8dd4-affa-2e531a098ae0.htm "GetTargetAMSNetId 方法 ")

[InitializationOnConnect 方法](../html/ea40c0b0-f579-3a64-7b45-573259ef706d.htm "InitializationOnConnect 方法 ")

[InitializationOnConnectAsync 方法](../html/c2a3f0f3-3a78-e102-d07f-d15053af2633.htm "InitializationOnConnectAsync 方法 ")

[PackCommandWithHeader 方法](../html/3895f6bc-bf81-5dbe-cb27-c12c6602126d.htm "PackCommandWithHeader 方法 ")

[Read 方法](../html/504e0a30-6e4b-ac79-e7e4-ce3aecedecb9.htm "Read 方法 ")

[ReadAdsDeviceInfo 方法](../html/1b58e4ea-7a7c-3b42-1f11-b077010b9a39.htm "ReadAdsDeviceInfo 方法 ")

[ReadAdsDeviceInfoAsync 方法](../html/18ef1d19-2639-b4b9-99c6-3e7caed8ef5a.htm "ReadAdsDeviceInfoAsync 方法 ")

[ReadAdsState 方法](../html/4645af4b-15b0-49a7-a3d3-bc07ef913ad7.htm "ReadAdsState 方法 ")

[ReadAdsStateAsync 方法](../html/831eac3f-b313-b2bd-3131-71bf3948b652.htm "ReadAdsStateAsync 方法 ")

[ReadAsync 方法](../html/ffaa1904-b9ab-79ff-efde-2543c610aa26.htm "ReadAsync 方法 ")

[ReadBool 方法](../html/b4ab6787-c834-7bf7-e069-1e8f7d5a4cad.htm "ReadBool 方法 ")

[ReadBoolAsync 方法](../html/1114ae0d-0359-668d-a6b2-6be76e7da28b.htm "ReadBoolAsync 方法 ")

[ReadByte 方法](../html/dadb0178-8d9c-2ce1-5495-e53d72771b82.htm "ReadByte 方法 ")

[ReadByteAsync 方法](../html/d61db681-dca5-f888-2dfd-c35faa269a4f.htm "ReadByteAsync 方法 ")

[ReadStruct 方法](../html/8e09ddd7-faa8-67c9-e099-eee24d59aa8d.htm "ReadStruct 方法 ")

[ReadValueHandle 方法](../html/7a0c126f-5f3e-f077-f90b-34ca2191c5b2.htm "ReadValueHandle 方法 ")

[ReadValueHandleAsync 方法](../html/0dadc9a4-90cc-1f29-627f-f8112aaa7df2.htm "ReadValueHandleAsync 方法 ")

[ReleaseSystemHandle 方法](../html/e2fa36f0-6d1f-2e02-e3de-bc2344e8d553.htm "ReleaseSystemHandle 方法 ")

[ReleaseSystemHandleAsync 方法](../html/cb4896ea-bbd4-52fa-4587-db7c137ff2d6.htm "ReleaseSystemHandleAsync 方法 ")

[SetSenderAMSNetId 方法](../html/6c085271-688a-909e-53c9-dec86b2a8e8a.htm "SetSenderAMSNetId 方法 ")

[SetTargetAMSNetId 方法](../html/af74dcd7-fb14-527d-5b79-881112e412ad.htm "SetTargetAMSNetId 方法 ")

[ToString 方法](../html/5cf941e3-75eb-cfcb-ec45-0065be0058c5.htm "ToString 方法 ")

[TransValueHandle 方法](../html/f7d19175-e124-9a8d-d36e-638b14d9c919.htm "TransValueHandle 方法 ")

[TransValueHandleAsync 方法](../html/efa735b1-7770-c061-0e6f-0486b84ff6b8.htm "TransValueHandleAsync 方法 ")

[UnpackResponseContent 方法](../html/9578724a-11cf-6cec-c183-18fb69db6412.htm "UnpackResponseContent 方法 ")

[Write 方法](../html/d6d8c193-854e-1069-a9f0-d426a76813cc.htm "Write 方法 ")

[WriteAdsState 方法](../html/c05364e4-eeb6-6b26-e817-0413145e4e8b.htm "WriteAdsState 方法 ")

[WriteAdsStateAsync 方法](../html/7a6c4cf6-7837-0873-3d9d-e6b99acec800.htm "WriteAdsStateAsync 方法 ")

[WriteAsync 方法](../html/7f760de5-65c9-46e0-b53e-fa719a52dafc.htm "WriteAsync 方法 ")

[WriteStruct(T) 方法](../html/a7aeee39-ab11-376a-71b5-02a066630636.htm "WriteStruct(T) 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| BeckhoffAdsNetGetNewNetMessage 方法 |

获取一个新的消息对象的方法，需要在继承类里面进行重写  
The method to get a new message object needs to be overridden in the inheritance class

**命名空间：**
 [HslCommunication.Profinet.Beckhoff](fb2af334-3662-b9ca-be3c-b0090bf4ef56.htm)  
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

[BeckhoffAdsNet 类](7ab15712-8688-063d-6b30-4b24803393fd.htm)

[HslCommunication.Profinet.Beckhoff 命名空间](fb2af334-3662-b9ca-be3c-b0090bf4ef56.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GetSenderAMSNetId 方法 

[原文連結](http://api.hslcommunication.cn/html/f743d65e-7896-7d75-ba57-6be4991ab9af.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Beckhoff](../html/fb2af334-3662-b9ca-be3c-b0090bf4ef56.htm "HslCommunication.Profinet.Beckhoff")

[BeckhoffAdsNet 类](../html/7ab15712-8688-063d-6b30-4b24803393fd.htm "BeckhoffAdsNet 类")

[BeckhoffAdsNet 方法](../html/e8bad1d3-d4ae-37ef-6ba7-e199dac65f0c.htm "BeckhoffAdsNet 方法")

[ClearTagsCache 方法](../html/3ca05c73-2ab5-b7b8-2d46-03da8a9c67ad.htm "ClearTagsCache 方法 ")

[DecideWhetherQAMessage 方法](../html/771586fe-9184-5312-e3da-bae61fe91d9a.htm "DecideWhetherQAMessage 方法 ")

[ExtraAfterReadFromCoreServer 方法](../html/3d26d83d-9287-54f5-3d35-c6d074ce6263.htm "ExtraAfterReadFromCoreServer 方法 ")

[GetNewNetMessage 方法](../html/215ca52b-b151-9cb6-0443-00b9be103188.htm "GetNewNetMessage 方法 ")

[GetSenderAMSNetId 方法](../html/f743d65e-7896-7d75-ba57-6be4991ab9af.htm "GetSenderAMSNetId 方法 ")

[GetTargetAMSNetId 方法](../html/dcc58722-6399-8dd4-affa-2e531a098ae0.htm "GetTargetAMSNetId 方法 ")

[InitializationOnConnect 方法](../html/ea40c0b0-f579-3a64-7b45-573259ef706d.htm "InitializationOnConnect 方法 ")

[InitializationOnConnectAsync 方法](../html/c2a3f0f3-3a78-e102-d07f-d15053af2633.htm "InitializationOnConnectAsync 方法 ")

[PackCommandWithHeader 方法](../html/3895f6bc-bf81-5dbe-cb27-c12c6602126d.htm "PackCommandWithHeader 方法 ")

[Read 方法](../html/504e0a30-6e4b-ac79-e7e4-ce3aecedecb9.htm "Read 方法 ")

[ReadAdsDeviceInfo 方法](../html/1b58e4ea-7a7c-3b42-1f11-b077010b9a39.htm "ReadAdsDeviceInfo 方法 ")

[ReadAdsDeviceInfoAsync 方法](../html/18ef1d19-2639-b4b9-99c6-3e7caed8ef5a.htm "ReadAdsDeviceInfoAsync 方法 ")

[ReadAdsState 方法](../html/4645af4b-15b0-49a7-a3d3-bc07ef913ad7.htm "ReadAdsState 方法 ")

[ReadAdsStateAsync 方法](../html/831eac3f-b313-b2bd-3131-71bf3948b652.htm "ReadAdsStateAsync 方法 ")

[ReadAsync 方法](../html/ffaa1904-b9ab-79ff-efde-2543c610aa26.htm "ReadAsync 方法 ")

[ReadBool 方法](../html/b4ab6787-c834-7bf7-e069-1e8f7d5a4cad.htm "ReadBool 方法 ")

[ReadBoolAsync 方法](../html/1114ae0d-0359-668d-a6b2-6be76e7da28b.htm "ReadBoolAsync 方法 ")

[ReadByte 方法](../html/dadb0178-8d9c-2ce1-5495-e53d72771b82.htm "ReadByte 方法 ")

[ReadByteAsync 方法](../html/d61db681-dca5-f888-2dfd-c35faa269a4f.htm "ReadByteAsync 方法 ")

[ReadStruct 方法](../html/8e09ddd7-faa8-67c9-e099-eee24d59aa8d.htm "ReadStruct 方法 ")

[ReadValueHandle 方法](../html/7a0c126f-5f3e-f077-f90b-34ca2191c5b2.htm "ReadValueHandle 方法 ")

[ReadValueHandleAsync 方法](../html/0dadc9a4-90cc-1f29-627f-f8112aaa7df2.htm "ReadValueHandleAsync 方法 ")

[ReleaseSystemHandle 方法](../html/e2fa36f0-6d1f-2e02-e3de-bc2344e8d553.htm "ReleaseSystemHandle 方法 ")

[ReleaseSystemHandleAsync 方法](../html/cb4896ea-bbd4-52fa-4587-db7c137ff2d6.htm "ReleaseSystemHandleAsync 方法 ")

[SetSenderAMSNetId 方法](../html/6c085271-688a-909e-53c9-dec86b2a8e8a.htm "SetSenderAMSNetId 方法 ")

[SetTargetAMSNetId 方法](../html/af74dcd7-fb14-527d-5b79-881112e412ad.htm "SetTargetAMSNetId 方法 ")

[ToString 方法](../html/5cf941e3-75eb-cfcb-ec45-0065be0058c5.htm "ToString 方法 ")

[TransValueHandle 方法](../html/f7d19175-e124-9a8d-d36e-638b14d9c919.htm "TransValueHandle 方法 ")

[TransValueHandleAsync 方法](../html/efa735b1-7770-c061-0e6f-0486b84ff6b8.htm "TransValueHandleAsync 方法 ")

[UnpackResponseContent 方法](../html/9578724a-11cf-6cec-c183-18fb69db6412.htm "UnpackResponseContent 方法 ")

[Write 方法](../html/d6d8c193-854e-1069-a9f0-d426a76813cc.htm "Write 方法 ")

[WriteAdsState 方法](../html/c05364e4-eeb6-6b26-e817-0413145e4e8b.htm "WriteAdsState 方法 ")

[WriteAdsStateAsync 方法](../html/7a6c4cf6-7837-0873-3d9d-e6b99acec800.htm "WriteAdsStateAsync 方法 ")

[WriteAsync 方法](../html/7f760de5-65c9-46e0-b53e-fa719a52dafc.htm "WriteAsync 方法 ")

[WriteStruct(T) 方法](../html/a7aeee39-ab11-376a-71b5-02a066630636.htm "WriteStruct(T) 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| BeckhoffAdsNetGetSenderAMSNetId 方法 |

获取当前发送的AMS的网络ID信息

**命名空间：**
 [HslCommunication.Profinet.Beckhoff](fb2af334-3662-b9ca-be3c-b0090bf4ef56.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public string GetSenderAMSNetId()
```

```
Public Function GetSenderAMSNetId As String
```

```
public:
String^ GetSenderAMSNetId()
```

```
member GetSenderAMSNetId : unit -> string 
```

#### 返回值

类型：String  
AMS发送信息

![](../icons/SectionExpanded.png)参见

#### 引用

[BeckhoffAdsNet 类](7ab15712-8688-063d-6b30-4b24803393fd.htm)

[HslCommunication.Profinet.Beckhoff 命名空间](fb2af334-3662-b9ca-be3c-b0090bf4ef56.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GetTargetAMSNetId 方法 

[原文連結](http://api.hslcommunication.cn/html/dcc58722-6399-8dd4-affa-2e531a098ae0.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Beckhoff](../html/fb2af334-3662-b9ca-be3c-b0090bf4ef56.htm "HslCommunication.Profinet.Beckhoff")

[BeckhoffAdsNet 类](../html/7ab15712-8688-063d-6b30-4b24803393fd.htm "BeckhoffAdsNet 类")

[BeckhoffAdsNet 方法](../html/e8bad1d3-d4ae-37ef-6ba7-e199dac65f0c.htm "BeckhoffAdsNet 方法")

[ClearTagsCache 方法](../html/3ca05c73-2ab5-b7b8-2d46-03da8a9c67ad.htm "ClearTagsCache 方法 ")

[DecideWhetherQAMessage 方法](../html/771586fe-9184-5312-e3da-bae61fe91d9a.htm "DecideWhetherQAMessage 方法 ")

[ExtraAfterReadFromCoreServer 方法](../html/3d26d83d-9287-54f5-3d35-c6d074ce6263.htm "ExtraAfterReadFromCoreServer 方法 ")

[GetNewNetMessage 方法](../html/215ca52b-b151-9cb6-0443-00b9be103188.htm "GetNewNetMessage 方法 ")

[GetSenderAMSNetId 方法](../html/f743d65e-7896-7d75-ba57-6be4991ab9af.htm "GetSenderAMSNetId 方法 ")

[GetTargetAMSNetId 方法](../html/dcc58722-6399-8dd4-affa-2e531a098ae0.htm "GetTargetAMSNetId 方法 ")

[InitializationOnConnect 方法](../html/ea40c0b0-f579-3a64-7b45-573259ef706d.htm "InitializationOnConnect 方法 ")

[InitializationOnConnectAsync 方法](../html/c2a3f0f3-3a78-e102-d07f-d15053af2633.htm "InitializationOnConnectAsync 方法 ")

[PackCommandWithHeader 方法](../html/3895f6bc-bf81-5dbe-cb27-c12c6602126d.htm "PackCommandWithHeader 方法 ")

[Read 方法](../html/504e0a30-6e4b-ac79-e7e4-ce3aecedecb9.htm "Read 方法 ")

[ReadAdsDeviceInfo 方法](../html/1b58e4ea-7a7c-3b42-1f11-b077010b9a39.htm "ReadAdsDeviceInfo 方法 ")

[ReadAdsDeviceInfoAsync 方法](../html/18ef1d19-2639-b4b9-99c6-3e7caed8ef5a.htm "ReadAdsDeviceInfoAsync 方法 ")

[ReadAdsState 方法](../html/4645af4b-15b0-49a7-a3d3-bc07ef913ad7.htm "ReadAdsState 方法 ")

[ReadAdsStateAsync 方法](../html/831eac3f-b313-b2bd-3131-71bf3948b652.htm "ReadAdsStateAsync 方法 ")

[ReadAsync 方法](../html/ffaa1904-b9ab-79ff-efde-2543c610aa26.htm "ReadAsync 方法 ")

[ReadBool 方法](../html/b4ab6787-c834-7bf7-e069-1e8f7d5a4cad.htm "ReadBool 方法 ")

[ReadBoolAsync 方法](../html/1114ae0d-0359-668d-a6b2-6be76e7da28b.htm "ReadBoolAsync 方法 ")

[ReadByte 方法](../html/dadb0178-8d9c-2ce1-5495-e53d72771b82.htm "ReadByte 方法 ")

[ReadByteAsync 方法](../html/d61db681-dca5-f888-2dfd-c35faa269a4f.htm "ReadByteAsync 方法 ")

[ReadStruct 方法](../html/8e09ddd7-faa8-67c9-e099-eee24d59aa8d.htm "ReadStruct 方法 ")

[ReadValueHandle 方法](../html/7a0c126f-5f3e-f077-f90b-34ca2191c5b2.htm "ReadValueHandle 方法 ")

[ReadValueHandleAsync 方法](../html/0dadc9a4-90cc-1f29-627f-f8112aaa7df2.htm "ReadValueHandleAsync 方法 ")

[ReleaseSystemHandle 方法](../html/e2fa36f0-6d1f-2e02-e3de-bc2344e8d553.htm "ReleaseSystemHandle 方法 ")

[ReleaseSystemHandleAsync 方法](../html/cb4896ea-bbd4-52fa-4587-db7c137ff2d6.htm "ReleaseSystemHandleAsync 方法 ")

[SetSenderAMSNetId 方法](../html/6c085271-688a-909e-53c9-dec86b2a8e8a.htm "SetSenderAMSNetId 方法 ")

[SetTargetAMSNetId 方法](../html/af74dcd7-fb14-527d-5b79-881112e412ad.htm "SetTargetAMSNetId 方法 ")

[ToString 方法](../html/5cf941e3-75eb-cfcb-ec45-0065be0058c5.htm "ToString 方法 ")

[TransValueHandle 方法](../html/f7d19175-e124-9a8d-d36e-638b14d9c919.htm "TransValueHandle 方法 ")

[TransValueHandleAsync 方法](../html/efa735b1-7770-c061-0e6f-0486b84ff6b8.htm "TransValueHandleAsync 方法 ")

[UnpackResponseContent 方法](../html/9578724a-11cf-6cec-c183-18fb69db6412.htm "UnpackResponseContent 方法 ")

[Write 方法](../html/d6d8c193-854e-1069-a9f0-d426a76813cc.htm "Write 方法 ")

[WriteAdsState 方法](../html/c05364e4-eeb6-6b26-e817-0413145e4e8b.htm "WriteAdsState 方法 ")

[WriteAdsStateAsync 方法](../html/7a6c4cf6-7837-0873-3d9d-e6b99acec800.htm "WriteAdsStateAsync 方法 ")

[WriteAsync 方法](../html/7f760de5-65c9-46e0-b53e-fa719a52dafc.htm "WriteAsync 方法 ")

[WriteStruct(T) 方法](../html/a7aeee39-ab11-376a-71b5-02a066630636.htm "WriteStruct(T) 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| BeckhoffAdsNetGetTargetAMSNetId 方法 |

获取当前目标的AMS网络的ID信息

**命名空间：**
 [HslCommunication.Profinet.Beckhoff](fb2af334-3662-b9ca-be3c-b0090bf4ef56.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public string GetTargetAMSNetId()
```

```
Public Function GetTargetAMSNetId As String
```

```
public:
String^ GetTargetAMSNetId()
```

```
member GetTargetAMSNetId : unit -> string 
```

#### 返回值

类型：String  
AMS目标信息

![](../icons/SectionExpanded.png)参见

#### 引用

[BeckhoffAdsNet 类](7ab15712-8688-063d-6b30-4b24803393fd.htm)

[HslCommunication.Profinet.Beckhoff 命名空间](fb2af334-3662-b9ca-be3c-b0090bf4ef56.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## InitializationOnConnect 方法 

[原文連結](http://api.hslcommunication.cn/html/ea40c0b0-f579-3a64-7b45-573259ef706d.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Beckhoff](../html/fb2af334-3662-b9ca-be3c-b0090bf4ef56.htm "HslCommunication.Profinet.Beckhoff")

[BeckhoffAdsNet 类](../html/7ab15712-8688-063d-6b30-4b24803393fd.htm "BeckhoffAdsNet 类")

[BeckhoffAdsNet 方法](../html/e8bad1d3-d4ae-37ef-6ba7-e199dac65f0c.htm "BeckhoffAdsNet 方法")

[ClearTagsCache 方法](../html/3ca05c73-2ab5-b7b8-2d46-03da8a9c67ad.htm "ClearTagsCache 方法 ")

[DecideWhetherQAMessage 方法](../html/771586fe-9184-5312-e3da-bae61fe91d9a.htm "DecideWhetherQAMessage 方法 ")

[ExtraAfterReadFromCoreServer 方法](../html/3d26d83d-9287-54f5-3d35-c6d074ce6263.htm "ExtraAfterReadFromCoreServer 方法 ")

[GetNewNetMessage 方法](../html/215ca52b-b151-9cb6-0443-00b9be103188.htm "GetNewNetMessage 方法 ")

[GetSenderAMSNetId 方法](../html/f743d65e-7896-7d75-ba57-6be4991ab9af.htm "GetSenderAMSNetId 方法 ")

[GetTargetAMSNetId 方法](../html/dcc58722-6399-8dd4-affa-2e531a098ae0.htm "GetTargetAMSNetId 方法 ")

[InitializationOnConnect 方法](../html/ea40c0b0-f579-3a64-7b45-573259ef706d.htm "InitializationOnConnect 方法 ")

[InitializationOnConnectAsync 方法](../html/c2a3f0f3-3a78-e102-d07f-d15053af2633.htm "InitializationOnConnectAsync 方法 ")

[PackCommandWithHeader 方法](../html/3895f6bc-bf81-5dbe-cb27-c12c6602126d.htm "PackCommandWithHeader 方法 ")

[Read 方法](../html/504e0a30-6e4b-ac79-e7e4-ce3aecedecb9.htm "Read 方法 ")

[ReadAdsDeviceInfo 方法](../html/1b58e4ea-7a7c-3b42-1f11-b077010b9a39.htm "ReadAdsDeviceInfo 方法 ")

[ReadAdsDeviceInfoAsync 方法](../html/18ef1d19-2639-b4b9-99c6-3e7caed8ef5a.htm "ReadAdsDeviceInfoAsync 方法 ")

[ReadAdsState 方法](../html/4645af4b-15b0-49a7-a3d3-bc07ef913ad7.htm "ReadAdsState 方法 ")

[ReadAdsStateAsync 方法](../html/831eac3f-b313-b2bd-3131-71bf3948b652.htm "ReadAdsStateAsync 方法 ")

[ReadAsync 方法](../html/ffaa1904-b9ab-79ff-efde-2543c610aa26.htm "ReadAsync 方法 ")

[ReadBool 方法](../html/b4ab6787-c834-7bf7-e069-1e8f7d5a4cad.htm "ReadBool 方法 ")

[ReadBoolAsync 方法](../html/1114ae0d-0359-668d-a6b2-6be76e7da28b.htm "ReadBoolAsync 方法 ")

[ReadByte 方法](../html/dadb0178-8d9c-2ce1-5495-e53d72771b82.htm "ReadByte 方法 ")

[ReadByteAsync 方法](../html/d61db681-dca5-f888-2dfd-c35faa269a4f.htm "ReadByteAsync 方法 ")

[ReadStruct 方法](../html/8e09ddd7-faa8-67c9-e099-eee24d59aa8d.htm "ReadStruct 方法 ")

[ReadValueHandle 方法](../html/7a0c126f-5f3e-f077-f90b-34ca2191c5b2.htm "ReadValueHandle 方法 ")

[ReadValueHandleAsync 方法](../html/0dadc9a4-90cc-1f29-627f-f8112aaa7df2.htm "ReadValueHandleAsync 方法 ")

[ReleaseSystemHandle 方法](../html/e2fa36f0-6d1f-2e02-e3de-bc2344e8d553.htm "ReleaseSystemHandle 方法 ")

[ReleaseSystemHandleAsync 方法](../html/cb4896ea-bbd4-52fa-4587-db7c137ff2d6.htm "ReleaseSystemHandleAsync 方法 ")

[SetSenderAMSNetId 方法](../html/6c085271-688a-909e-53c9-dec86b2a8e8a.htm "SetSenderAMSNetId 方法 ")

[SetTargetAMSNetId 方法](../html/af74dcd7-fb14-527d-5b79-881112e412ad.htm "SetTargetAMSNetId 方法 ")

[ToString 方法](../html/5cf941e3-75eb-cfcb-ec45-0065be0058c5.htm "ToString 方法 ")

[TransValueHandle 方法](../html/f7d19175-e124-9a8d-d36e-638b14d9c919.htm "TransValueHandle 方法 ")

[TransValueHandleAsync 方法](../html/efa735b1-7770-c061-0e6f-0486b84ff6b8.htm "TransValueHandleAsync 方法 ")

[UnpackResponseContent 方法](../html/9578724a-11cf-6cec-c183-18fb69db6412.htm "UnpackResponseContent 方法 ")

[Write 方法](../html/d6d8c193-854e-1069-a9f0-d426a76813cc.htm "Write 方法 ")

[WriteAdsState 方法](../html/c05364e4-eeb6-6b26-e817-0413145e4e8b.htm "WriteAdsState 方法 ")

[WriteAdsStateAsync 方法](../html/7a6c4cf6-7837-0873-3d9d-e6b99acec800.htm "WriteAdsStateAsync 方法 ")

[WriteAsync 方法](../html/7f760de5-65c9-46e0-b53e-fa719a52dafc.htm "WriteAsync 方法 ")

[WriteStruct(T) 方法](../html/a7aeee39-ab11-376a-71b5-02a066630636.htm "WriteStruct(T) 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| BeckhoffAdsNetInitializationOnConnect 方法 |

根据实际的协议选择是否重写本方法，有些协议在创建连接之后，需要进行一些初始化的信号握手，才能最终建立网络通道。  
Whether to rewrite this method is based on the actual protocol. Some protocols require some initial signal handshake to establish a network channel after the connection is created.

**命名空间：**
 [HslCommunication.Profinet.Beckhoff](fb2af334-3662-b9ca-be3c-b0090bf4ef56.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
protected override OperateResult InitializationOnConnect()
```

```
Protected Overrides Function InitializationOnConnect As OperateResult
```

```
protected:
virtual OperateResult^ InitializationOnConnect() override
```

```
abstract InitializationOnConnect : unit -> OperateResult 
override InitializationOnConnect : unit -> OperateResult
```

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否初始化成功，依据具体的协议进行重写

![](../icons/SectionExpanded.png)示例

有些协议不需要握手信号，比如三菱的MC协议，Modbus协议，西门子和欧姆龙就存在握手信息，此处的例子是继承本类后重写的西门子的协议示例

西门子重连示例

[复制](# "复制")

```
        /// <inheritdoc/>
        public override OperateResult<byte[]> ReadFromCoreServer( CommunicationPipe pipe, byte[] send, bool hasResponseData, bool usePackAndUnpack )
        {
            while (true)
            {
                OperateResult<byte[]> read = base.ReadFromCoreServer( pipe, send, hasResponseData, usePackAndUnpack );
                if (!read.IsSuccess) return read;

                if (read.Content?.Length >= 4)
                {
                    if ((read.Content[2] * 256 + read.Content[3]) != 0x07) return read;
                }
            }
        }

        /// <inheritdoc/>
        protected override OperateResult InitializationOnConnect( )
        {
            // 第一次握手 -> First handshake
            OperateResult<byte[]> read_first = ReadFromCoreServer( this.CommunicationPipe, plcHead1, hasResponseData: true, usePackAndUnpack: true );
            if (!read_first.IsSuccess) return read_first;

            // 第二次握手 -> Second handshake
            OperateResult<byte[]> read_second = ReadFromCoreServer( this.CommunicationPipe, plcHead2, hasResponseData: true, usePackAndUnpack: true );
            if (!read_second.IsSuccess) return read_second;

            // 调整单次接收的pdu长度信息
            pdu_length = ByteTransform.TransUInt16( read_second.Content.SelectLast( 2 ), 0 ) - 28;
            if (pdu_length < 200) pdu_length = 200;

            incrementCount = new SoftIncrementCount( ushort.MaxValue, 1 );
            // 返回成功的信号 -> Return a successful signal
            return OperateResult.CreateSuccessResult( );
        }

#if !NET35 && !NET20
        /// <inheritdoc/>
        public async override Task<OperateResult<byte[]>> ReadFromCoreServerAsync( CommunicationPipe pipe, byte[] send, bool hasResponseData, bool usePackAndUnpack )
        {
            while (true)
            {
                OperateResult<byte[]> read = await base.ReadFromCoreServerAsync( pipe, send, hasResponseData, usePackAndUnpack );
                if (!read.IsSuccess) return read;

                if (read.Content?.Length >= 4)
                {
                    if ((read.Content[2] * 256 + read.Content[3]) != 0x07) return read;
                }
            }
        }

        /// <inheritdoc/>
        protected override async Task<OperateResult> InitializationOnConnectAsync( )
        {
            // 第一次握手 -> First handshake
            OperateResult<byte[]> read_first = await ReadFromCoreServerAsync( this.CommunicationPipe, plcHead1, hasResponseData: true, usePackAndUnpack: true );
            if (!read_first.IsSuccess) return read_first;

            // 第二次握手 -> Second handshake
            OperateResult<byte[]> read_second = await ReadFromCoreServerAsync( this.CommunicationPipe, plcHead2, hasResponseData: true, usePackAndUnpack: true );
            if (!read_second.IsSuccess) return read_second;

            // 调整单次接收的pdu长度信息
            pdu_length = ByteTransform.TransUInt16( read_second.Content.SelectLast( 2 ), 0 ) - 28;
            if (pdu_length < 200) pdu_length = 200;

            incrementCount = new SoftIncrementCount( ushort.MaxValue, 1 );
            // 返回成功的信号 -> Return a successful signal
            return OperateResult.CreateSuccessResult( );
        }
#endif
```

![](../icons/SectionExpanded.png)参见

#### 引用

[BeckhoffAdsNet 类](7ab15712-8688-063d-6b30-4b24803393fd.htm)

[HslCommunication.Profinet.Beckhoff 命名空间](fb2af334-3662-b9ca-be3c-b0090bf4ef56.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## InitializationOnConnectAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/c2a3f0f3-3a78-e102-d07f-d15053af2633.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Beckhoff](../html/fb2af334-3662-b9ca-be3c-b0090bf4ef56.htm "HslCommunication.Profinet.Beckhoff")

[BeckhoffAdsNet 类](../html/7ab15712-8688-063d-6b30-4b24803393fd.htm "BeckhoffAdsNet 类")

[BeckhoffAdsNet 方法](../html/e8bad1d3-d4ae-37ef-6ba7-e199dac65f0c.htm "BeckhoffAdsNet 方法")

[ClearTagsCache 方法](../html/3ca05c73-2ab5-b7b8-2d46-03da8a9c67ad.htm "ClearTagsCache 方法 ")

[DecideWhetherQAMessage 方法](../html/771586fe-9184-5312-e3da-bae61fe91d9a.htm "DecideWhetherQAMessage 方法 ")

[ExtraAfterReadFromCoreServer 方法](../html/3d26d83d-9287-54f5-3d35-c6d074ce6263.htm "ExtraAfterReadFromCoreServer 方法 ")

[GetNewNetMessage 方法](../html/215ca52b-b151-9cb6-0443-00b9be103188.htm "GetNewNetMessage 方法 ")

[GetSenderAMSNetId 方法](../html/f743d65e-7896-7d75-ba57-6be4991ab9af.htm "GetSenderAMSNetId 方法 ")

[GetTargetAMSNetId 方法](../html/dcc58722-6399-8dd4-affa-2e531a098ae0.htm "GetTargetAMSNetId 方法 ")

[InitializationOnConnect 方法](../html/ea40c0b0-f579-3a64-7b45-573259ef706d.htm "InitializationOnConnect 方法 ")

[InitializationOnConnectAsync 方法](../html/c2a3f0f3-3a78-e102-d07f-d15053af2633.htm "InitializationOnConnectAsync 方法 ")

[PackCommandWithHeader 方法](../html/3895f6bc-bf81-5dbe-cb27-c12c6602126d.htm "PackCommandWithHeader 方法 ")

[Read 方法](../html/504e0a30-6e4b-ac79-e7e4-ce3aecedecb9.htm "Read 方法 ")

[ReadAdsDeviceInfo 方法](../html/1b58e4ea-7a7c-3b42-1f11-b077010b9a39.htm "ReadAdsDeviceInfo 方法 ")

[ReadAdsDeviceInfoAsync 方法](../html/18ef1d19-2639-b4b9-99c6-3e7caed8ef5a.htm "ReadAdsDeviceInfoAsync 方法 ")

[ReadAdsState 方法](../html/4645af4b-15b0-49a7-a3d3-bc07ef913ad7.htm "ReadAdsState 方法 ")

[ReadAdsStateAsync 方法](../html/831eac3f-b313-b2bd-3131-71bf3948b652.htm "ReadAdsStateAsync 方法 ")

[ReadAsync 方法](../html/ffaa1904-b9ab-79ff-efde-2543c610aa26.htm "ReadAsync 方法 ")

[ReadBool 方法](../html/b4ab6787-c834-7bf7-e069-1e8f7d5a4cad.htm "ReadBool 方法 ")

[ReadBoolAsync 方法](../html/1114ae0d-0359-668d-a6b2-6be76e7da28b.htm "ReadBoolAsync 方法 ")

[ReadByte 方法](../html/dadb0178-8d9c-2ce1-5495-e53d72771b82.htm "ReadByte 方法 ")

[ReadByteAsync 方法](../html/d61db681-dca5-f888-2dfd-c35faa269a4f.htm "ReadByteAsync 方法 ")

[ReadStruct 方法](../html/8e09ddd7-faa8-67c9-e099-eee24d59aa8d.htm "ReadStruct 方法 ")

[ReadValueHandle 方法](../html/7a0c126f-5f3e-f077-f90b-34ca2191c5b2.htm "ReadValueHandle 方法 ")

[ReadValueHandleAsync 方法](../html/0dadc9a4-90cc-1f29-627f-f8112aaa7df2.htm "ReadValueHandleAsync 方法 ")

[ReleaseSystemHandle 方法](../html/e2fa36f0-6d1f-2e02-e3de-bc2344e8d553.htm "ReleaseSystemHandle 方法 ")

[ReleaseSystemHandleAsync 方法](../html/cb4896ea-bbd4-52fa-4587-db7c137ff2d6.htm "ReleaseSystemHandleAsync 方法 ")

[SetSenderAMSNetId 方法](../html/6c085271-688a-909e-53c9-dec86b2a8e8a.htm "SetSenderAMSNetId 方法 ")

[SetTargetAMSNetId 方法](../html/af74dcd7-fb14-527d-5b79-881112e412ad.htm "SetTargetAMSNetId 方法 ")

[ToString 方法](../html/5cf941e3-75eb-cfcb-ec45-0065be0058c5.htm "ToString 方法 ")

[TransValueHandle 方法](../html/f7d19175-e124-9a8d-d36e-638b14d9c919.htm "TransValueHandle 方法 ")

[TransValueHandleAsync 方法](../html/efa735b1-7770-c061-0e6f-0486b84ff6b8.htm "TransValueHandleAsync 方法 ")

[UnpackResponseContent 方法](../html/9578724a-11cf-6cec-c183-18fb69db6412.htm "UnpackResponseContent 方法 ")

[Write 方法](../html/d6d8c193-854e-1069-a9f0-d426a76813cc.htm "Write 方法 ")

[WriteAdsState 方法](../html/c05364e4-eeb6-6b26-e817-0413145e4e8b.htm "WriteAdsState 方法 ")

[WriteAdsStateAsync 方法](../html/7a6c4cf6-7837-0873-3d9d-e6b99acec800.htm "WriteAdsStateAsync 方法 ")

[WriteAsync 方法](../html/7f760de5-65c9-46e0-b53e-fa719a52dafc.htm "WriteAsync 方法 ")

[WriteStruct(T) 方法](../html/a7aeee39-ab11-376a-71b5-02a066630636.htm "WriteStruct(T) 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| BeckhoffAdsNetInitializationOnConnectAsync 方法 |

根据实际的协议选择是否重写本方法，有些协议在创建连接之后，需要进行一些初始化的信号握手，才能最终建立网络通道。  
Whether to rewrite this method is based on the actual protocol. Some protocols require some initial signal handshake to establish a network channel after the connection is created.

**命名空间：**
 [HslCommunication.Profinet.Beckhoff](fb2af334-3662-b9ca-be3c-b0090bf4ef56.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
protected override Task<OperateResult> InitializationOnConnectAsync()
```

```
Protected Overrides Function InitializationOnConnectAsync As Task(Of OperateResult)
```

```
protected:
virtual Task<OperateResult^>^ InitializationOnConnectAsync() override
```

```
abstract InitializationOnConnectAsync : unit -> Task<OperateResult> 
override InitializationOnConnectAsync : unit -> Task<OperateResult>
```

#### 返回值

类型：Task[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否初始化成功，依据具体的协议进行重写

![](../icons/SectionExpanded.png)示例

有些协议不需要握手信号，比如三菱的MC协议，Modbus协议，西门子和欧姆龙就存在握手信息，此处的例子是继承本类后重写的西门子的协议示例

西门子重连示例

[复制](# "复制")

```
        /// <inheritdoc/>
        public override OperateResult<byte[]> ReadFromCoreServer( CommunicationPipe pipe, byte[] send, bool hasResponseData, bool usePackAndUnpack )
        {
            while (true)
            {
                OperateResult<byte[]> read = base.ReadFromCoreServer( pipe, send, hasResponseData, usePackAndUnpack );
                if (!read.IsSuccess) return read;

                if (read.Content?.Length >= 4)
                {
                    if ((read.Content[2] * 256 + read.Content[3]) != 0x07) return read;
                }
            }
        }

        /// <inheritdoc/>
        protected override OperateResult InitializationOnConnect( )
        {
            // 第一次握手 -> First handshake
            OperateResult<byte[]> read_first = ReadFromCoreServer( this.CommunicationPipe, plcHead1, hasResponseData: true, usePackAndUnpack: true );
            if (!read_first.IsSuccess) return read_first;

            // 第二次握手 -> Second handshake
            OperateResult<byte[]> read_second = ReadFromCoreServer( this.CommunicationPipe, plcHead2, hasResponseData: true, usePackAndUnpack: true );
            if (!read_second.IsSuccess) return read_second;

            // 调整单次接收的pdu长度信息
            pdu_length = ByteTransform.TransUInt16( read_second.Content.SelectLast( 2 ), 0 ) - 28;
            if (pdu_length < 200) pdu_length = 200;

            incrementCount = new SoftIncrementCount( ushort.MaxValue, 1 );
            // 返回成功的信号 -> Return a successful signal
            return OperateResult.CreateSuccessResult( );
        }

#if !NET35 && !NET20
        /// <inheritdoc/>
        public async override Task<OperateResult<byte[]>> ReadFromCoreServerAsync( CommunicationPipe pipe, byte[] send, bool hasResponseData, bool usePackAndUnpack )
        {
            while (true)
            {
                OperateResult<byte[]> read = await base.ReadFromCoreServerAsync( pipe, send, hasResponseData, usePackAndUnpack );
                if (!read.IsSuccess) return read;

                if (read.Content?.Length >= 4)
                {
                    if ((read.Content[2] * 256 + read.Content[3]) != 0x07) return read;
                }
            }
        }

        /// <inheritdoc/>
        protected override async Task<OperateResult> InitializationOnConnectAsync( )
        {
            // 第一次握手 -> First handshake
            OperateResult<byte[]> read_first = await ReadFromCoreServerAsync( this.CommunicationPipe, plcHead1, hasResponseData: true, usePackAndUnpack: true );
            if (!read_first.IsSuccess) return read_first;

            // 第二次握手 -> Second handshake
            OperateResult<byte[]> read_second = await ReadFromCoreServerAsync( this.CommunicationPipe, plcHead2, hasResponseData: true, usePackAndUnpack: true );
            if (!read_second.IsSuccess) return read_second;

            // 调整单次接收的pdu长度信息
            pdu_length = ByteTransform.TransUInt16( read_second.Content.SelectLast( 2 ), 0 ) - 28;
            if (pdu_length < 200) pdu_length = 200;

            incrementCount = new SoftIncrementCount( ushort.MaxValue, 1 );
            // 返回成功的信号 -> Return a successful signal
            return OperateResult.CreateSuccessResult( );
        }
#endif
```

![](../icons/SectionExpanded.png)参见

#### 引用

[BeckhoffAdsNet 类](7ab15712-8688-063d-6b30-4b24803393fd.htm)

[HslCommunication.Profinet.Beckhoff 命名空间](fb2af334-3662-b9ca-be3c-b0090bf4ef56.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)