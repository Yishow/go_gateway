# HslCommunication - HslCommunication.Instrument.Light

> 分類頁數: 30



---
## HslCommunication.Instrument.Light

[原文連結](http://api.hslcommunication.cn/html/00d11148-efb6-2e43-ae78-c995e930d868.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.Light](../html/00d11148-efb6-2e43-ae78-c995e930d868.htm "HslCommunication.Instrument.Light")

[ShineInLightData 类](../html/6981aaa1-3c66-69ba-49da-a707306a9267.htm "ShineInLightData 类")

[ShineInLightSourceController 类](../html/5de4a57b-1202-4bcd-4650-a524ba7166df.htm "ShineInLightSourceController 类")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslCommunication.Instrument.Light 命名空间 |

[缺少 "N:HslCommunication.Instrument.Light" 的 <summary> 文档]

![](../icons/SectionExpanded.png)类

|  | 类 | 说明 |
| --- | --- | --- |
| 公共类 | [ShineInLightData](6981aaa1-3c66-69ba-49da-a707306a9267.htm) | 光源的数据信息 |
| 公共类 | [ShineInLightSourceController](5de4a57b-1202-4bcd-4650-a524ba7166df.htm) | 昱行智造科技（深圳）有限公司的光源控制器，可以控制灯的亮暗，控制灯的颜色，通道等信息。 |

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ShineInLightData 类

[原文連結](http://api.hslcommunication.cn/html/6981aaa1-3c66-69ba-49da-a707306a9267.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.Light](../html/00d11148-efb6-2e43-ae78-c995e930d868.htm "HslCommunication.Instrument.Light")

[ShineInLightData 类](../html/6981aaa1-3c66-69ba-49da-a707306a9267.htm "ShineInLightData 类")

[ShineInLightData 构造函数](../html/9f0ba81c-0a7d-6a67-a292-363dd231230e.htm "ShineInLightData 构造函数 ")

[ShineInLightData 属性](../html/b15b673a-6a71-b189-ef25-c1295ccf5afc.htm "ShineInLightData 属性")

[ShineInLightData 方法](../html/b2756ccc-ffce-a407-b9db-4bf0f09b2d22.htm "ShineInLightData 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ShineInLightData 类 |

光源的数据信息

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  HslCommunication.Instrument.LightShineInLightData

**命名空间：**
 [HslCommunication.Instrument.Light](00d11148-efb6-2e43-ae78-c995e930d868.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class ShineInLightData
```

```
Public Class ShineInLightData
```

```
public ref class ShineInLightData
```

```
type ShineInLightData =  class end
```

ShineInLightData 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [ShineInLightData](a0a4671e-c91b-44d4-3dbd-2aadb9208f32.htm) | 实例化一个默认的对象 |
| 公共方法 | [ShineInLightData(Byte)](9460ea6c-8082-a93c-10de-556996d28c70.htm) | 使用指定的原始数据来获取当前的对象 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [Address](d4aa65dc-1600-be79-5cce-5e97cbfa687b.htm) | 控制器的地址选择位 |
| 公共属性 | [Channel](c4da2fb9-6521-ab18-afb8-8c1f137f928b.htm) | 通道数据，01-08H的值 |
| 公共属性 | [Color](a7a1657b-42e4-08fa-60ae-237772e9eb68.htm) | 光源颜色信息，1:红色 2:绿色 3:蓝色 4:白色(默认) |
| 公共属性 | [Light](352f43b3-73e7-8320-de63-e74fcfff1c2c.htm) | 光源的亮度信息，00-FF，值越大，亮度越大 |
| 公共属性 | [LightDegree](95236d7b-3597-7fe8-7928-7da08ebfda6a.htm) | 光源的亮度等级，1-3 |
| 公共属性 | [PulseWidth](eafb5dee-b971-98ed-1f1b-d46cf349dafb.htm) | 脉冲宽度，01-14H |
| 公共属性 | [WorkMode](e8dfd953-e6d0-5a43-f743-791f969bdf26.htm) | 光源的工作模式，00:延时常亮 01:通道一频闪 02:通道二频闪 03:通道一二频闪 04:普通常亮 05:关闭 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法 | [GetSourceData](7762dc95-6ca5-d81f-2bf1-ea91f4bb1962.htm) | 获取原始的数据信息 |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法 | [ParseFrom](3f246eee-bf80-4e67-6a49-9b9cc8235014.htm) | 从原始的信息解析光源的数据 |
| 公共方法 | [ToString](617f906c-304f-e852-9ec9-52f6d4cb5b9d.htm) | (重写 ObjectToString.) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.Instrument.Light 命名空间](00d11148-efb6-2e43-ae78-c995e930d868.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ShineInLightData 构造函数 

[原文連結](http://api.hslcommunication.cn/html/9f0ba81c-0a7d-6a67-a292-363dd231230e.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.Light](../html/00d11148-efb6-2e43-ae78-c995e930d868.htm "HslCommunication.Instrument.Light")

[ShineInLightData 类](../html/6981aaa1-3c66-69ba-49da-a707306a9267.htm "ShineInLightData 类")

[ShineInLightData 构造函数](../html/9f0ba81c-0a7d-6a67-a292-363dd231230e.htm "ShineInLightData 构造函数 ")

[ShineInLightData 构造函数](../html/a0a4671e-c91b-44d4-3dbd-2aadb9208f32.htm "ShineInLightData 构造函数 ")

[ShineInLightData 构造函数 (Byte[])](../html/9460ea6c-8082-a93c-10de-556996d28c70.htm "ShineInLightData 构造函数 (Byte[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ShineInLightData 构造函数 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [ShineInLightData](a0a4671e-c91b-44d4-3dbd-2aadb9208f32.htm) | 实例化一个默认的对象 |
| 公共方法 | [ShineInLightData(Byte)](9460ea6c-8082-a93c-10de-556996d28c70.htm) | 使用指定的原始数据来获取当前的对象 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[ShineInLightData 类](6981aaa1-3c66-69ba-49da-a707306a9267.htm)

[HslCommunication.Instrument.Light 命名空间](00d11148-efb6-2e43-ae78-c995e930d868.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ShineInLightData 构造函数 

[原文連結](http://api.hslcommunication.cn/html/a0a4671e-c91b-44d4-3dbd-2aadb9208f32.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.Light](../html/00d11148-efb6-2e43-ae78-c995e930d868.htm "HslCommunication.Instrument.Light")

[ShineInLightData 类](../html/6981aaa1-3c66-69ba-49da-a707306a9267.htm "ShineInLightData 类")

[ShineInLightData 构造函数](../html/9f0ba81c-0a7d-6a67-a292-363dd231230e.htm "ShineInLightData 构造函数 ")

[ShineInLightData 构造函数](../html/a0a4671e-c91b-44d4-3dbd-2aadb9208f32.htm "ShineInLightData 构造函数 ")

[ShineInLightData 构造函数 (Byte[])](../html/9460ea6c-8082-a93c-10de-556996d28c70.htm "ShineInLightData 构造函数 (Byte[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ShineInLightData 构造函数 |

实例化一个默认的对象

**命名空间：**
 [HslCommunication.Instrument.Light](00d11148-efb6-2e43-ae78-c995e930d868.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public ShineInLightData()
```

```
Public Sub New
```

```
public:
ShineInLightData()
```

```
new : unit -> ShineInLightData
```

![](../icons/SectionExpanded.png)参见

#### 引用

[ShineInLightData 类](6981aaa1-3c66-69ba-49da-a707306a9267.htm)

[ShineInLightData 重载](9f0ba81c-0a7d-6a67-a292-363dd231230e.htm)

[HslCommunication.Instrument.Light 命名空间](00d11148-efb6-2e43-ae78-c995e930d868.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ShineInLightData 构造函数 (Byte[])

[原文連結](http://api.hslcommunication.cn/html/9460ea6c-8082-a93c-10de-556996d28c70.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.Light](../html/00d11148-efb6-2e43-ae78-c995e930d868.htm "HslCommunication.Instrument.Light")

[ShineInLightData 类](../html/6981aaa1-3c66-69ba-49da-a707306a9267.htm "ShineInLightData 类")

[ShineInLightData 构造函数](../html/9f0ba81c-0a7d-6a67-a292-363dd231230e.htm "ShineInLightData 构造函数 ")

[ShineInLightData 构造函数](../html/a0a4671e-c91b-44d4-3dbd-2aadb9208f32.htm "ShineInLightData 构造函数 ")

[ShineInLightData 构造函数 (Byte[])](../html/9460ea6c-8082-a93c-10de-556996d28c70.htm "ShineInLightData 构造函数 (Byte[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ShineInLightData 构造函数 (Byte) |

使用指定的原始数据来获取当前的对象

**命名空间：**
 [HslCommunication.Instrument.Light](00d11148-efb6-2e43-ae78-c995e930d868.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public ShineInLightData(
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
ShineInLightData(
	array<unsigned char>^ data
)
```

```
new : 
        data : byte[] -> ShineInLightData
```

#### 参数

data
:   类型：SystemByte  
    原始数据

![](../icons/SectionExpanded.png)参见

#### 引用

[ShineInLightData 类](6981aaa1-3c66-69ba-49da-a707306a9267.htm)

[ShineInLightData 重载](9f0ba81c-0a7d-6a67-a292-363dd231230e.htm)

[HslCommunication.Instrument.Light 命名空间](00d11148-efb6-2e43-ae78-c995e930d868.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ShineInLightData 属性

[原文連結](http://api.hslcommunication.cn/html/b15b673a-6a71-b189-ef25-c1295ccf5afc.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.Light](../html/00d11148-efb6-2e43-ae78-c995e930d868.htm "HslCommunication.Instrument.Light")

[ShineInLightData 类](../html/6981aaa1-3c66-69ba-49da-a707306a9267.htm "ShineInLightData 类")

[ShineInLightData 属性](../html/b15b673a-6a71-b189-ef25-c1295ccf5afc.htm "ShineInLightData 属性")

[Address 属性](../html/d4aa65dc-1600-be79-5cce-5e97cbfa687b.htm "Address 属性 ")

[Channel 属性](../html/c4da2fb9-6521-ab18-afb8-8c1f137f928b.htm "Channel 属性 ")

[Color 属性](../html/a7a1657b-42e4-08fa-60ae-237772e9eb68.htm "Color 属性 ")

[Light 属性](../html/352f43b3-73e7-8320-de63-e74fcfff1c2c.htm "Light 属性 ")

[LightDegree 属性](../html/95236d7b-3597-7fe8-7928-7da08ebfda6a.htm "LightDegree 属性 ")

[PulseWidth 属性](../html/eafb5dee-b971-98ed-1f1b-d46cf349dafb.htm "PulseWidth 属性 ")

[WorkMode 属性](../html/e8dfd953-e6d0-5a43-f743-791f969bdf26.htm "WorkMode 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ShineInLightData 属性 |

[ShineInLightData](6981aaa1-3c66-69ba-49da-a707306a9267.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [Address](d4aa65dc-1600-be79-5cce-5e97cbfa687b.htm) | 控制器的地址选择位 |
| 公共属性 | [Channel](c4da2fb9-6521-ab18-afb8-8c1f137f928b.htm) | 通道数据，01-08H的值 |
| 公共属性 | [Color](a7a1657b-42e4-08fa-60ae-237772e9eb68.htm) | 光源颜色信息，1:红色 2:绿色 3:蓝色 4:白色(默认) |
| 公共属性 | [Light](352f43b3-73e7-8320-de63-e74fcfff1c2c.htm) | 光源的亮度信息，00-FF，值越大，亮度越大 |
| 公共属性 | [LightDegree](95236d7b-3597-7fe8-7928-7da08ebfda6a.htm) | 光源的亮度等级，1-3 |
| 公共属性 | [PulseWidth](eafb5dee-b971-98ed-1f1b-d46cf349dafb.htm) | 脉冲宽度，01-14H |
| 公共属性 | [WorkMode](e8dfd953-e6d0-5a43-f743-791f969bdf26.htm) | 光源的工作模式，00:延时常亮 01:通道一频闪 02:通道二频闪 03:通道一二频闪 04:普通常亮 05:关闭 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[ShineInLightData 类](6981aaa1-3c66-69ba-49da-a707306a9267.htm)

[HslCommunication.Instrument.Light 命名空间](00d11148-efb6-2e43-ae78-c995e930d868.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Address 属性 

[原文連結](http://api.hslcommunication.cn/html/d4aa65dc-1600-be79-5cce-5e97cbfa687b.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.Light](../html/00d11148-efb6-2e43-ae78-c995e930d868.htm "HslCommunication.Instrument.Light")

[ShineInLightData 类](../html/6981aaa1-3c66-69ba-49da-a707306a9267.htm "ShineInLightData 类")

[ShineInLightData 属性](../html/b15b673a-6a71-b189-ef25-c1295ccf5afc.htm "ShineInLightData 属性")

[Address 属性](../html/d4aa65dc-1600-be79-5cce-5e97cbfa687b.htm "Address 属性 ")

[Channel 属性](../html/c4da2fb9-6521-ab18-afb8-8c1f137f928b.htm "Channel 属性 ")

[Color 属性](../html/a7a1657b-42e4-08fa-60ae-237772e9eb68.htm "Color 属性 ")

[Light 属性](../html/352f43b3-73e7-8320-de63-e74fcfff1c2c.htm "Light 属性 ")

[LightDegree 属性](../html/95236d7b-3597-7fe8-7928-7da08ebfda6a.htm "LightDegree 属性 ")

[PulseWidth 属性](../html/eafb5dee-b971-98ed-1f1b-d46cf349dafb.htm "PulseWidth 属性 ")

[WorkMode 属性](../html/e8dfd953-e6d0-5a43-f743-791f969bdf26.htm "WorkMode 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ShineInLightDataAddress 属性 |

控制器的地址选择位

**命名空间：**
 [HslCommunication.Instrument.Light](00d11148-efb6-2e43-ae78-c995e930d868.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public byte Address { get; set; }
```

```
Public Property Address As Byte
	Get
	Set
```

```
public:
property unsigned char Address {
	unsigned char get ();
	void set (unsigned char value);
}
```

```
member Address : byte with get, set
```

#### 属性值

类型：Byte

![](../icons/SectionExpanded.png)参见

#### 引用

[ShineInLightData 类](6981aaa1-3c66-69ba-49da-a707306a9267.htm)

[HslCommunication.Instrument.Light 命名空间](00d11148-efb6-2e43-ae78-c995e930d868.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Channel 属性 

[原文連結](http://api.hslcommunication.cn/html/c4da2fb9-6521-ab18-afb8-8c1f137f928b.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.Light](../html/00d11148-efb6-2e43-ae78-c995e930d868.htm "HslCommunication.Instrument.Light")

[ShineInLightData 类](../html/6981aaa1-3c66-69ba-49da-a707306a9267.htm "ShineInLightData 类")

[ShineInLightData 属性](../html/b15b673a-6a71-b189-ef25-c1295ccf5afc.htm "ShineInLightData 属性")

[Address 属性](../html/d4aa65dc-1600-be79-5cce-5e97cbfa687b.htm "Address 属性 ")

[Channel 属性](../html/c4da2fb9-6521-ab18-afb8-8c1f137f928b.htm "Channel 属性 ")

[Color 属性](../html/a7a1657b-42e4-08fa-60ae-237772e9eb68.htm "Color 属性 ")

[Light 属性](../html/352f43b3-73e7-8320-de63-e74fcfff1c2c.htm "Light 属性 ")

[LightDegree 属性](../html/95236d7b-3597-7fe8-7928-7da08ebfda6a.htm "LightDegree 属性 ")

[PulseWidth 属性](../html/eafb5dee-b971-98ed-1f1b-d46cf349dafb.htm "PulseWidth 属性 ")

[WorkMode 属性](../html/e8dfd953-e6d0-5a43-f743-791f969bdf26.htm "WorkMode 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ShineInLightDataChannel 属性 |

通道数据，01-08H的值

**命名空间：**
 [HslCommunication.Instrument.Light](00d11148-efb6-2e43-ae78-c995e930d868.htm)  
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

[ShineInLightData 类](6981aaa1-3c66-69ba-49da-a707306a9267.htm)

[HslCommunication.Instrument.Light 命名空间](00d11148-efb6-2e43-ae78-c995e930d868.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Color 属性 

[原文連結](http://api.hslcommunication.cn/html/a7a1657b-42e4-08fa-60ae-237772e9eb68.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.Light](../html/00d11148-efb6-2e43-ae78-c995e930d868.htm "HslCommunication.Instrument.Light")

[ShineInLightData 类](../html/6981aaa1-3c66-69ba-49da-a707306a9267.htm "ShineInLightData 类")

[ShineInLightData 属性](../html/b15b673a-6a71-b189-ef25-c1295ccf5afc.htm "ShineInLightData 属性")

[Address 属性](../html/d4aa65dc-1600-be79-5cce-5e97cbfa687b.htm "Address 属性 ")

[Channel 属性](../html/c4da2fb9-6521-ab18-afb8-8c1f137f928b.htm "Channel 属性 ")

[Color 属性](../html/a7a1657b-42e4-08fa-60ae-237772e9eb68.htm "Color 属性 ")

[Light 属性](../html/352f43b3-73e7-8320-de63-e74fcfff1c2c.htm "Light 属性 ")

[LightDegree 属性](../html/95236d7b-3597-7fe8-7928-7da08ebfda6a.htm "LightDegree 属性 ")

[PulseWidth 属性](../html/eafb5dee-b971-98ed-1f1b-d46cf349dafb.htm "PulseWidth 属性 ")

[WorkMode 属性](../html/e8dfd953-e6d0-5a43-f743-791f969bdf26.htm "WorkMode 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ShineInLightDataColor 属性 |

光源颜色信息，1:红色 2:绿色 3:蓝色 4:白色(默认)

**命名空间：**
 [HslCommunication.Instrument.Light](00d11148-efb6-2e43-ae78-c995e930d868.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public byte Color { get; set; }
```

```
Public Property Color As Byte
	Get
	Set
```

```
public:
property unsigned char Color {
	unsigned char get ();
	void set (unsigned char value);
}
```

```
member Color : byte with get, set
```

#### 属性值

类型：Byte

![](../icons/SectionExpanded.png)参见

#### 引用

[ShineInLightData 类](6981aaa1-3c66-69ba-49da-a707306a9267.htm)

[HslCommunication.Instrument.Light 命名空间](00d11148-efb6-2e43-ae78-c995e930d868.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Light 属性 

[原文連結](http://api.hslcommunication.cn/html/352f43b3-73e7-8320-de63-e74fcfff1c2c.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.Light](../html/00d11148-efb6-2e43-ae78-c995e930d868.htm "HslCommunication.Instrument.Light")

[ShineInLightData 类](../html/6981aaa1-3c66-69ba-49da-a707306a9267.htm "ShineInLightData 类")

[ShineInLightData 属性](../html/b15b673a-6a71-b189-ef25-c1295ccf5afc.htm "ShineInLightData 属性")

[Address 属性](../html/d4aa65dc-1600-be79-5cce-5e97cbfa687b.htm "Address 属性 ")

[Channel 属性](../html/c4da2fb9-6521-ab18-afb8-8c1f137f928b.htm "Channel 属性 ")

[Color 属性](../html/a7a1657b-42e4-08fa-60ae-237772e9eb68.htm "Color 属性 ")

[Light 属性](../html/352f43b3-73e7-8320-de63-e74fcfff1c2c.htm "Light 属性 ")

[LightDegree 属性](../html/95236d7b-3597-7fe8-7928-7da08ebfda6a.htm "LightDegree 属性 ")

[PulseWidth 属性](../html/eafb5dee-b971-98ed-1f1b-d46cf349dafb.htm "PulseWidth 属性 ")

[WorkMode 属性](../html/e8dfd953-e6d0-5a43-f743-791f969bdf26.htm "WorkMode 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ShineInLightDataLight 属性 |

光源的亮度信息，00-FF，值越大，亮度越大

**命名空间：**
 [HslCommunication.Instrument.Light](00d11148-efb6-2e43-ae78-c995e930d868.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public byte Light { get; set; }
```

```
Public Property Light As Byte
	Get
	Set
```

```
public:
property unsigned char Light {
	unsigned char get ();
	void set (unsigned char value);
}
```

```
member Light : byte with get, set
```

#### 属性值

类型：Byte

![](../icons/SectionExpanded.png)参见

#### 引用

[ShineInLightData 类](6981aaa1-3c66-69ba-49da-a707306a9267.htm)

[HslCommunication.Instrument.Light 命名空间](00d11148-efb6-2e43-ae78-c995e930d868.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## LightDegree 属性 

[原文連結](http://api.hslcommunication.cn/html/95236d7b-3597-7fe8-7928-7da08ebfda6a.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.Light](../html/00d11148-efb6-2e43-ae78-c995e930d868.htm "HslCommunication.Instrument.Light")

[ShineInLightData 类](../html/6981aaa1-3c66-69ba-49da-a707306a9267.htm "ShineInLightData 类")

[ShineInLightData 属性](../html/b15b673a-6a71-b189-ef25-c1295ccf5afc.htm "ShineInLightData 属性")

[Address 属性](../html/d4aa65dc-1600-be79-5cce-5e97cbfa687b.htm "Address 属性 ")

[Channel 属性](../html/c4da2fb9-6521-ab18-afb8-8c1f137f928b.htm "Channel 属性 ")

[Color 属性](../html/a7a1657b-42e4-08fa-60ae-237772e9eb68.htm "Color 属性 ")

[Light 属性](../html/352f43b3-73e7-8320-de63-e74fcfff1c2c.htm "Light 属性 ")

[LightDegree 属性](../html/95236d7b-3597-7fe8-7928-7da08ebfda6a.htm "LightDegree 属性 ")

[PulseWidth 属性](../html/eafb5dee-b971-98ed-1f1b-d46cf349dafb.htm "PulseWidth 属性 ")

[WorkMode 属性](../html/e8dfd953-e6d0-5a43-f743-791f969bdf26.htm "WorkMode 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ShineInLightDataLightDegree 属性 |

光源的亮度等级，1-3

**命名空间：**
 [HslCommunication.Instrument.Light](00d11148-efb6-2e43-ae78-c995e930d868.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public byte LightDegree { get; set; }
```

```
Public Property LightDegree As Byte
	Get
	Set
```

```
public:
property unsigned char LightDegree {
	unsigned char get ();
	void set (unsigned char value);
}
```

```
member LightDegree : byte with get, set
```

#### 属性值

类型：Byte

![](../icons/SectionExpanded.png)参见

#### 引用

[ShineInLightData 类](6981aaa1-3c66-69ba-49da-a707306a9267.htm)

[HslCommunication.Instrument.Light 命名空间](00d11148-efb6-2e43-ae78-c995e930d868.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## PulseWidth 属性 

[原文連結](http://api.hslcommunication.cn/html/eafb5dee-b971-98ed-1f1b-d46cf349dafb.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.Light](../html/00d11148-efb6-2e43-ae78-c995e930d868.htm "HslCommunication.Instrument.Light")

[ShineInLightData 类](../html/6981aaa1-3c66-69ba-49da-a707306a9267.htm "ShineInLightData 类")

[ShineInLightData 属性](../html/b15b673a-6a71-b189-ef25-c1295ccf5afc.htm "ShineInLightData 属性")

[Address 属性](../html/d4aa65dc-1600-be79-5cce-5e97cbfa687b.htm "Address 属性 ")

[Channel 属性](../html/c4da2fb9-6521-ab18-afb8-8c1f137f928b.htm "Channel 属性 ")

[Color 属性](../html/a7a1657b-42e4-08fa-60ae-237772e9eb68.htm "Color 属性 ")

[Light 属性](../html/352f43b3-73e7-8320-de63-e74fcfff1c2c.htm "Light 属性 ")

[LightDegree 属性](../html/95236d7b-3597-7fe8-7928-7da08ebfda6a.htm "LightDegree 属性 ")

[PulseWidth 属性](../html/eafb5dee-b971-98ed-1f1b-d46cf349dafb.htm "PulseWidth 属性 ")

[WorkMode 属性](../html/e8dfd953-e6d0-5a43-f743-791f969bdf26.htm "WorkMode 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ShineInLightDataPulseWidth 属性 |

脉冲宽度，01-14H

**命名空间：**
 [HslCommunication.Instrument.Light](00d11148-efb6-2e43-ae78-c995e930d868.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public byte PulseWidth { get; set; }
```

```
Public Property PulseWidth As Byte
	Get
	Set
```

```
public:
property unsigned char PulseWidth {
	unsigned char get ();
	void set (unsigned char value);
}
```

```
member PulseWidth : byte with get, set
```

#### 属性值

类型：Byte

![](../icons/SectionExpanded.png)参见

#### 引用

[ShineInLightData 类](6981aaa1-3c66-69ba-49da-a707306a9267.htm)

[HslCommunication.Instrument.Light 命名空间](00d11148-efb6-2e43-ae78-c995e930d868.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## WorkMode 属性 

[原文連結](http://api.hslcommunication.cn/html/e8dfd953-e6d0-5a43-f743-791f969bdf26.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.Light](../html/00d11148-efb6-2e43-ae78-c995e930d868.htm "HslCommunication.Instrument.Light")

[ShineInLightData 类](../html/6981aaa1-3c66-69ba-49da-a707306a9267.htm "ShineInLightData 类")

[ShineInLightData 属性](../html/b15b673a-6a71-b189-ef25-c1295ccf5afc.htm "ShineInLightData 属性")

[Address 属性](../html/d4aa65dc-1600-be79-5cce-5e97cbfa687b.htm "Address 属性 ")

[Channel 属性](../html/c4da2fb9-6521-ab18-afb8-8c1f137f928b.htm "Channel 属性 ")

[Color 属性](../html/a7a1657b-42e4-08fa-60ae-237772e9eb68.htm "Color 属性 ")

[Light 属性](../html/352f43b3-73e7-8320-de63-e74fcfff1c2c.htm "Light 属性 ")

[LightDegree 属性](../html/95236d7b-3597-7fe8-7928-7da08ebfda6a.htm "LightDegree 属性 ")

[PulseWidth 属性](../html/eafb5dee-b971-98ed-1f1b-d46cf349dafb.htm "PulseWidth 属性 ")

[WorkMode 属性](../html/e8dfd953-e6d0-5a43-f743-791f969bdf26.htm "WorkMode 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ShineInLightDataWorkMode 属性 |

光源的工作模式，00:延时常亮 01:通道一频闪 02:通道二频闪 03:通道一二频闪 04:普通常亮 05:关闭

**命名空间：**
 [HslCommunication.Instrument.Light](00d11148-efb6-2e43-ae78-c995e930d868.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public byte WorkMode { get; set; }
```

```
Public Property WorkMode As Byte
	Get
	Set
```

```
public:
property unsigned char WorkMode {
	unsigned char get ();
	void set (unsigned char value);
}
```

```
member WorkMode : byte with get, set
```

#### 属性值

类型：Byte

![](../icons/SectionExpanded.png)参见

#### 引用

[ShineInLightData 类](6981aaa1-3c66-69ba-49da-a707306a9267.htm)

[HslCommunication.Instrument.Light 命名空间](00d11148-efb6-2e43-ae78-c995e930d868.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ShineInLightData 方法

[原文連結](http://api.hslcommunication.cn/html/b2756ccc-ffce-a407-b9db-4bf0f09b2d22.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.Light](../html/00d11148-efb6-2e43-ae78-c995e930d868.htm "HslCommunication.Instrument.Light")

[ShineInLightData 类](../html/6981aaa1-3c66-69ba-49da-a707306a9267.htm "ShineInLightData 类")

[ShineInLightData 方法](../html/b2756ccc-ffce-a407-b9db-4bf0f09b2d22.htm "ShineInLightData 方法")

[GetSourceData 方法](../html/7762dc95-6ca5-d81f-2bf1-ea91f4bb1962.htm "GetSourceData 方法 ")

[ParseFrom 方法](../html/3f246eee-bf80-4e67-6a49-9b9cc8235014.htm "ParseFrom 方法 ")

[ToString 方法](../html/617f906c-304f-e852-9ec9-52f6d4cb5b9d.htm "ToString 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ShineInLightData 方法 |

[ShineInLightData](6981aaa1-3c66-69ba-49da-a707306a9267.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法 | [GetSourceData](7762dc95-6ca5-d81f-2bf1-ea91f4bb1962.htm) | 获取原始的数据信息 |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法 | [ParseFrom](3f246eee-bf80-4e67-6a49-9b9cc8235014.htm) | 从原始的信息解析光源的数据 |
| 公共方法 | [ToString](617f906c-304f-e852-9ec9-52f6d4cb5b9d.htm) | (重写 ObjectToString.) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[ShineInLightData 类](6981aaa1-3c66-69ba-49da-a707306a9267.htm)

[HslCommunication.Instrument.Light 命名空间](00d11148-efb6-2e43-ae78-c995e930d868.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GetSourceData 方法 

[原文連結](http://api.hslcommunication.cn/html/7762dc95-6ca5-d81f-2bf1-ea91f4bb1962.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.Light](../html/00d11148-efb6-2e43-ae78-c995e930d868.htm "HslCommunication.Instrument.Light")

[ShineInLightData 类](../html/6981aaa1-3c66-69ba-49da-a707306a9267.htm "ShineInLightData 类")

[ShineInLightData 方法](../html/b2756ccc-ffce-a407-b9db-4bf0f09b2d22.htm "ShineInLightData 方法")

[GetSourceData 方法](../html/7762dc95-6ca5-d81f-2bf1-ea91f4bb1962.htm "GetSourceData 方法 ")

[ParseFrom 方法](../html/3f246eee-bf80-4e67-6a49-9b9cc8235014.htm "ParseFrom 方法 ")

[ToString 方法](../html/617f906c-304f-e852-9ec9-52f6d4cb5b9d.htm "ToString 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ShineInLightDataGetSourceData 方法 |

获取原始的数据信息

**命名空间：**
 [HslCommunication.Instrument.Light](00d11148-efb6-2e43-ae78-c995e930d868.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public byte[] GetSourceData()
```

```
Public Function GetSourceData As Byte()
```

```
public:
array<unsigned char>^ GetSourceData()
```

```
member GetSourceData : unit -> byte[] 
```

#### 返回值

类型：Byte  
原始的字节信息

![](../icons/SectionExpanded.png)参见

#### 引用

[ShineInLightData 类](6981aaa1-3c66-69ba-49da-a707306a9267.htm)

[HslCommunication.Instrument.Light 命名空间](00d11148-efb6-2e43-ae78-c995e930d868.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ParseFrom 方法 

[原文連結](http://api.hslcommunication.cn/html/3f246eee-bf80-4e67-6a49-9b9cc8235014.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.Light](../html/00d11148-efb6-2e43-ae78-c995e930d868.htm "HslCommunication.Instrument.Light")

[ShineInLightData 类](../html/6981aaa1-3c66-69ba-49da-a707306a9267.htm "ShineInLightData 类")

[ShineInLightData 方法](../html/b2756ccc-ffce-a407-b9db-4bf0f09b2d22.htm "ShineInLightData 方法")

[GetSourceData 方法](../html/7762dc95-6ca5-d81f-2bf1-ea91f4bb1962.htm "GetSourceData 方法 ")

[ParseFrom 方法](../html/3f246eee-bf80-4e67-6a49-9b9cc8235014.htm "ParseFrom 方法 ")

[ToString 方法](../html/617f906c-304f-e852-9ec9-52f6d4cb5b9d.htm "ToString 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ShineInLightDataParseFrom 方法 |

从原始的信息解析光源的数据

**命名空间：**
 [HslCommunication.Instrument.Light](00d11148-efb6-2e43-ae78-c995e930d868.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public void ParseFrom(
	byte[] data
)
```

```
Public Sub ParseFrom ( 
	data As Byte()
)
```

```
public:
void ParseFrom(
	array<unsigned char>^ data
)
```

```
member ParseFrom : 
        data : byte[] -> unit 
```

#### 参数

data
:   类型：SystemByte  
    原始的数据信息

![](../icons/SectionExpanded.png)参见

#### 引用

[ShineInLightData 类](6981aaa1-3c66-69ba-49da-a707306a9267.htm)

[HslCommunication.Instrument.Light 命名空间](00d11148-efb6-2e43-ae78-c995e930d868.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ToString 方法 

[原文連結](http://api.hslcommunication.cn/html/617f906c-304f-e852-9ec9-52f6d4cb5b9d.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.Light](../html/00d11148-efb6-2e43-ae78-c995e930d868.htm "HslCommunication.Instrument.Light")

[ShineInLightData 类](../html/6981aaa1-3c66-69ba-49da-a707306a9267.htm "ShineInLightData 类")

[ShineInLightData 方法](../html/b2756ccc-ffce-a407-b9db-4bf0f09b2d22.htm "ShineInLightData 方法")

[GetSourceData 方法](../html/7762dc95-6ca5-d81f-2bf1-ea91f4bb1962.htm "GetSourceData 方法 ")

[ParseFrom 方法](../html/3f246eee-bf80-4e67-6a49-9b9cc8235014.htm "ParseFrom 方法 ")

[ToString 方法](../html/617f906c-304f-e852-9ec9-52f6d4cb5b9d.htm "ToString 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ShineInLightDataToString 方法 |

[缺少 "M:HslCommunication.Instrument.Light.ShineInLightData.ToString" 的 <summary> 文档]

**命名空间：**
 [HslCommunication.Instrument.Light](00d11148-efb6-2e43-ae78-c995e930d868.htm)  
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

[缺少 "M:HslCommunication.Instrument.Light.ShineInLightData.ToString" 的 <returns> 文档]

![](../icons/SectionExpanded.png)参见

#### 引用

[ShineInLightData 类](6981aaa1-3c66-69ba-49da-a707306a9267.htm)

[HslCommunication.Instrument.Light 命名空间](00d11148-efb6-2e43-ae78-c995e930d868.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ShineInLightSourceController 类

[原文連結](http://api.hslcommunication.cn/html/5de4a57b-1202-4bcd-4650-a524ba7166df.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.Light](../html/00d11148-efb6-2e43-ae78-c995e930d868.htm "HslCommunication.Instrument.Light")

[ShineInLightSourceController 类](../html/5de4a57b-1202-4bcd-4650-a524ba7166df.htm "ShineInLightSourceController 类")

[ShineInLightSourceController 构造函数](../html/bcec4279-07cc-5bf8-615a-5c0be053134b.htm "ShineInLightSourceController 构造函数 ")

[ShineInLightSourceController 属性](../html/0b153ce1-fdbd-236f-fa54-0155c7ad19ab.htm "ShineInLightSourceController 属性")

[ShineInLightSourceController 方法](../html/cbbaa7f0-2b3e-2557-2808-9aa825c888d7.htm "ShineInLightSourceController 方法")

[ShineInLightSourceController 字段](../html/fbd7a8ae-835d-0ec4-a250-2764569c1f83.htm "ShineInLightSourceController 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ShineInLightSourceController 类 |

昱行智造科技（深圳）有限公司的光源控制器，可以控制灯的亮暗，控制灯的颜色，通道等信息。

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  [HslCommunication.Core.NetBinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)  
    [HslCommunication.SerialSerialBase](e877e20b-ef65-266c-05a5-5db58e6bad0c.htm)  
      HslCommunication.Instrument.LightShineInLightSourceController

**命名空间：**
 [HslCommunication.Instrument.Light](00d11148-efb6-2e43-ae78-c995e930d868.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class ShineInLightSourceController : SerialBase
```

```
Public Class ShineInLightSourceController
	Inherits SerialBase
```

```
public ref class ShineInLightSourceController : public SerialBase
```

```
type ShineInLightSourceController =  
    class
        inherit SerialBase
    end
```

ShineInLightSourceController 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [ShineInLightSourceController](bcec4279-07cc-5bf8-615a-5c0be053134b.htm) | 实例化一个默认的对象 |

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
| 公共方法静态成员 | [BuildReadCommand](148fde6a-f94f-92b0-8f1d-80f57f35ce05.htm) | 构建读取数据的报文命令 |
| 公共方法静态成员 | [BuildWriteCommand](59521eec-5735-ef60-b7c4-f4707ecf9dec.htm) | 构建写入数据的报文命令 |
| 公共方法 | [Close](992f28b3-b99b-2196-274e-2f25e9bd3c1d.htm) | 关闭当前的串口连接  Close the current serial connection (继承自 [SerialBase](e877e20b-ef65-266c-05a5-5db58e6bad0c.htm)。) |
| 受保护的方法 | [DecideWhetherQAMessage](acc40cda-e310-8537-6685-59b67eeb16ff.htm) | 决定当前的消息是否是用于问答机制返回的消息，默认直接返回 true, 实际的情况需要根据协议进行重写方法  To determine whether the current message is the message returned by the question answering mechanism, the default is true. In actual cases, the rewriting method needs to be performed according to the protocol (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [Dispose](c83e3a27-3a09-f67c-9c7f-96b7a5a4c8b0.htm) | 释放当前的对象 (继承自 [SerialBase](e877e20b-ef65-266c-05a5-5db58e6bad0c.htm)。) |
| 受保护的方法 | [Dispose(Boolean)](c509106d-35a7-2ff2-7dd8-ca67e16d7950.htm) | 释放当前的对象 (继承自 [SerialBase](e877e20b-ef65-266c-05a5-5db58e6bad0c.htm)。) |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | [ExtraAfterReadFromCoreServer](06a49a15-a3ce-5e39-28c7-2b78851eca2e.htm) | 和服务器交互完成的时候调用的方法，可以根据读写结果进行一些额外的操作，具体的操作需要根据实际的需求来重写实现  The method called when the interaction with the server is completed can perform some additional operations based on the read and write results. The specific operations need to be rewritten according to actual needs. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法静态成员 | [ExtractActualData](9937e217-ca0b-0ea1-5e75-28be1a0af552.htm) | 把服务器反馈的数据解析成实际的命令 |
| 受保护的方法代码示例 | [ExtraOnDisconnect](de9dbba2-3b2c-1c5e-53c8-4e6f8a44e2c4.htm) | 根据实际的协议选择是否重写本方法，有些协议在断开连接之前，需要发送一些报文来关闭当前的网络通道  Select whether to rewrite this method according to the actual protocol. Some protocols need to send some packets to close the current network channel before disconnecting. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法代码示例 | [ExtraOnDisconnectAsync](210f593e-16d8-77c4-1178-58fbcfe9b27b.htm) | 根据实际的协议选择是否重写本方法，有些协议在断开连接之前，需要发送一些报文来关闭当前的网络通道  Select whether to rewrite this method according to the actual protocol. Some protocols need to send some packets to close the current network channel before disconnecting. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 受保护的方法 | [GetLogTextFromBinary](67f3b2f7-957c-fa0b-8320-6799237157db.htm) | 获取当前的报文进行日志记录的时候，是否使用二进制的格式记录，默认返回 [LogMsgFormatBinary](2ed2ff87-caec-0e44-f08b-e2b015225a58.htm)，重写可以根据session对象分别返回不同记录模式  Whether to log the current packet in binary format, the default return is [LogMsgFormatBinary](2ed2ff87-caec-0e44-f08b-e2b015225a58.htm). If you want to override it, different recording modes can be returned according to session (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [GetNewNetMessage](16636a1b-a6f0-9c0f-9d9c-c94e26c67459.htm) | 获取一个新的消息对象的方法，需要在继承类里面进行重写  The method to get a new message object needs to be overridden in the inheritance class (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
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
| 公共方法静态成员 | [PackCommand](1094f5da-f719-64c7-4ffc-c42c5f8004c7.htm) | 将命令和数据打包成用于发送的报文 |
| 公共方法 | [PackCommandWithHeader](d0b2cd76-c306-11fc-38c7-af8ad306cce7.htm) | 对当前的命令进行打包处理，通常是携带命令头内容，标记当前的命令的长度信息，需要进行重写，否则默认不打包  The current command is packaged, usually carrying the content of the command header, marking the length of the current command, and it needs to be rewritten, otherwise it is not packaged by default (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [Read](d4363a0b-2a6e-d4e2-4e40-7abb240b0382.htm) | 读取光源控制器的参数信息，需要传入通道号信息，读取到详细的内容参照[ShineInLightData](6981aaa1-3c66-69ba-49da-a707306a9267.htm)的值 |
| 公共方法代码示例 | [ReadFromCoreServer(Byte)](8146dc0c-1edc-3ca7-428a-a3bdbd35e6c5.htm) | 将二进制的数据发送到管道中去，然后从管道里接收二进制的数据回来，并返回是否成功的结果对象。  Send binary data to the pipeline, and then receive binary data back from the pipeline, and return whether the success of the result object (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServer(IEnumerableByte)](044cfdd9-ec76-b89f-2f2b-5f439025006c.htm) | 将多个数据报文按顺序发到设备，并从设备接收返回的数据内容，然后拼接成一个Byte[]信息，需要重写[UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm)方法才能返回正确的结果。  Send multiple data packets to the device in sequence, and receive the returned data content from the device, and then splicing them into a Byte[] message, you need to rewrite [UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm) method to return the correct result. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServer(Byte, Boolean, Boolean)](b32ce8e8-a959-fcc2-6736-14b39fdb9a40.htm) | 将二进制的数据发送到管道中去，然后从管道里接收二进制的数据回来，并返回是否成功的结果对象。  Send binary data to the pipeline, and then receive binary data back from the pipeline, and return whether the success of the result object (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServer(CommunicationPipe, Byte, Boolean, Boolean)](ef36d56b-df09-f727-e8d3-316d37a13ee9.htm) | 使用指定的管道来进行数据通信，发送原始数据到管道，然后从管道接收相关的数据返回，本方法无锁 (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServerAsync(Byte)](a4de5a24-c67a-7b02-db05-33b85eaf6e1e.htm) | 将二进制的数据发送到管道中去，然后从管道里接收二进制的数据回来，并返回是否成功的结果对象。  Send binary data to the pipeline, and then receive binary data back from the pipeline, and return whether the success of the result object (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServerAsync(IEnumerableByte)](bcdb59b7-a70a-d5c0-1fc6-9ececcd39ea2.htm) | 将多个数据报文按顺序发到设备，并从设备接收返回的数据内容，然后拼接成一个Byte[]信息，需要重写[UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm)方法才能返回正确的结果。  Send multiple data packets to the device in sequence, and receive the returned data content from the device, and then splicing them into a Byte[] message, you need to rewrite [UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm) method to return the correct result. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServerAsync(Byte, Boolean, Boolean)](1c79c204-5ab5-9342-e65d-7455c20ea577.htm) | 将二进制的数据发送到管道中去，然后从管道里接收二进制的数据回来，并返回是否成功的结果对象。  Send binary data to the pipeline, and then receive binary data back from the pipeline, and return whether the success of the result object (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServerAsync(CommunicationPipe, Byte, Boolean, Boolean)](a254a19d-3c29-147b-335b-8adccdf4ee2f.htm) | 使用指定的管道来进行数据通信，发送原始数据到管道，然后从管道接收相关的数据返回，本方法无锁 (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [SerialPortInni(String)](76fee18c-6275-b5cc-2f03-e9095178cae5.htm) | 初始化串口信息，57600波特率，8位数据位，1位停止位，偶校验  Initial serial port information, 57600 baud rate, 8 data bits, 1 stop bit, even parity (重写 [SerialBaseSerialPortInni(String)](7dff5e26-6594-291e-4792-b13db9c5d613.htm).) |
| 公共方法 | [SerialPortInni(ActionSerialPort)](53585409-6243-9359-1846-bf8b4e266708.htm) | 根据自定义初始化方法进行初始化串口信息  Initialize the serial port information according to the custom initialization method (继承自 [SerialBase](e877e20b-ef65-266c-05a5-5db58e6bad0c.htm)。) |
| 公共方法 | [SerialPortInni(String, Int32)](601606e3-f3c0-4bd7-49f3-cf420c123fb9.htm) | 初始化串口信息，波特率，8位数据位，1位停止位，偶校验  Initializes serial port information, baud rate, 8-bit data bit, 1-bit stop bit, even parity (重写 [SerialBaseSerialPortInni(String, Int32)](0c7aa163-1d81-54fd-f4f5-103a7e9a341d.htm).) |
| 公共方法 | [SerialPortInni(String, Int32, Int32, StopBits, Parity)](420c818b-48e9-fa02-59e2-f4d5897d6802.htm) | 初始化串口信息，波特率，数据位，停止位，奇偶校验需要全部自己来指定  Start serial port information, baud rate, data bit, stop bit, parity all need to be specified (继承自 [SerialBase](e877e20b-ef65-266c-05a5-5db58e6bad0c.htm)。) |
| 公共方法 | [SetDtuPipe](f765d1d7-234a-b229-9457-97e25e024dc9.htm) | 将当前的通信对象设置DTU模式，允许传入现成的管道，并返回初始化结果，如果该设备重写了握手报文，就是返回握手结果  Set the current communication object to DTU mode, allow the existing pipe to be passed in, and return the initialization result, if the device rewrites the handshake packet, the handshake result is returned (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [SetDtuPipeAsync](9d1fd12d-b804-3828-a627-a2c0eef1d23f.htm) | 将当前的通信对象设置DTU模式，允许传入现成的管道，并返回初始化结果，如果该设备重写了握手报文，就是返回握手结果  Set the current communication object to DTU mode, allow the existing pipe to be passed in, and return the initialization result, if the device rewrites the handshake packet, the handshake result is returned (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ToString](1154f133-182e-ecfb-ebc5-347911fc7506.htm) | (重写 [SerialBaseToString](1ac09324-4d19-863d-d683-cbe804ef2124.htm).) |
| 公共方法 | [UnpackResponseContent](208017a7-8da2-33f8-3273-d7e7e850100f.htm) | 根据对方返回的报文命令，对命令进行基本的拆包，例如各种Modbus协议拆包为统一的核心报文，还支持对报文的验证  According to the message command returned by the other party, the command is basically unpacked, for example, various Modbus protocols are unpacked into a unified core message, and the verification of the message is also supported (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [Write](9d0bc031-6b4d-20cd-62e9-f52cb3f4b0b6.htm) | 将光源控制器的数据写入到设备，返回是否写入成功 |

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

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.Instrument.Light 命名空间](00d11148-efb6-2e43-ae78-c995e930d868.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ShineInLightSourceController 构造函数 

[原文連結](http://api.hslcommunication.cn/html/bcec4279-07cc-5bf8-615a-5c0be053134b.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.Light](../html/00d11148-efb6-2e43-ae78-c995e930d868.htm "HslCommunication.Instrument.Light")

[ShineInLightSourceController 类](../html/5de4a57b-1202-4bcd-4650-a524ba7166df.htm "ShineInLightSourceController 类")

[ShineInLightSourceController 构造函数](../html/bcec4279-07cc-5bf8-615a-5c0be053134b.htm "ShineInLightSourceController 构造函数 ")

[ShineInLightSourceController 属性](../html/0b153ce1-fdbd-236f-fa54-0155c7ad19ab.htm "ShineInLightSourceController 属性")

[ShineInLightSourceController 方法](../html/cbbaa7f0-2b3e-2557-2808-9aa825c888d7.htm "ShineInLightSourceController 方法")

[ShineInLightSourceController 字段](../html/fbd7a8ae-835d-0ec4-a250-2764569c1f83.htm "ShineInLightSourceController 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ShineInLightSourceController 构造函数 |

实例化一个默认的对象

**命名空间：**
 [HslCommunication.Instrument.Light](00d11148-efb6-2e43-ae78-c995e930d868.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public ShineInLightSourceController()
```

```
Public Sub New
```

```
public:
ShineInLightSourceController()
```

```
new : unit -> ShineInLightSourceController
```

![](../icons/SectionExpanded.png)参见

#### 引用

[ShineInLightSourceController 类](5de4a57b-1202-4bcd-4650-a524ba7166df.htm)

[HslCommunication.Instrument.Light 命名空间](00d11148-efb6-2e43-ae78-c995e930d868.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ShineInLightSourceController 属性

[原文連結](http://api.hslcommunication.cn/html/0b153ce1-fdbd-236f-fa54-0155c7ad19ab.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.Light](../html/00d11148-efb6-2e43-ae78-c995e930d868.htm "HslCommunication.Instrument.Light")

[ShineInLightSourceController 类](../html/5de4a57b-1202-4bcd-4650-a524ba7166df.htm "ShineInLightSourceController 类")

[ShineInLightSourceController 构造函数](../html/bcec4279-07cc-5bf8-615a-5c0be053134b.htm "ShineInLightSourceController 构造函数 ")

[ShineInLightSourceController 属性](../html/0b153ce1-fdbd-236f-fa54-0155c7ad19ab.htm "ShineInLightSourceController 属性")

[ShineInLightSourceController 方法](../html/cbbaa7f0-2b3e-2557-2808-9aa825c888d7.htm "ShineInLightSourceController 方法")

[ShineInLightSourceController 字段](../html/fbd7a8ae-835d-0ec4-a250-2764569c1f83.htm "ShineInLightSourceController 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ShineInLightSourceController 属性 |

[ShineInLightSourceController](5de4a57b-1202-4bcd-4650-a524ba7166df.htm) 类型公开以下成员。

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

[ShineInLightSourceController 类](5de4a57b-1202-4bcd-4650-a524ba7166df.htm)

[HslCommunication.Instrument.Light 命名空间](00d11148-efb6-2e43-ae78-c995e930d868.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ShineInLightSourceController 方法

[原文連結](http://api.hslcommunication.cn/html/cbbaa7f0-2b3e-2557-2808-9aa825c888d7.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.Light](../html/00d11148-efb6-2e43-ae78-c995e930d868.htm "HslCommunication.Instrument.Light")

[ShineInLightSourceController 类](../html/5de4a57b-1202-4bcd-4650-a524ba7166df.htm "ShineInLightSourceController 类")

[ShineInLightSourceController 方法](../html/cbbaa7f0-2b3e-2557-2808-9aa825c888d7.htm "ShineInLightSourceController 方法")

[BuildReadCommand 方法](../html/148fde6a-f94f-92b0-8f1d-80f57f35ce05.htm "BuildReadCommand 方法 ")

[BuildWriteCommand 方法](../html/59521eec-5735-ef60-b7c4-f4707ecf9dec.htm "BuildWriteCommand 方法 ")

[ExtractActualData 方法](../html/9937e217-ca0b-0ea1-5e75-28be1a0af552.htm "ExtractActualData 方法 ")

[PackCommand 方法](../html/1094f5da-f719-64c7-4ffc-c42c5f8004c7.htm "PackCommand 方法 ")

[Read 方法](../html/d4363a0b-2a6e-d4e2-4e40-7abb240b0382.htm "Read 方法 ")

[SerialPortInni 方法](../html/b06cbc7a-09af-2f7e-78e8-a6015e79cd15.htm "SerialPortInni 方法 ")

[ToString 方法](../html/1154f133-182e-ecfb-ebc5-347911fc7506.htm "ToString 方法 ")

[Write 方法](../html/9d0bc031-6b4d-20cd-62e9-f52cb3f4b0b6.htm "Write 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ShineInLightSourceController 方法 |

[ShineInLightSourceController](5de4a57b-1202-4bcd-4650-a524ba7166df.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法静态成员 | [BuildReadCommand](148fde6a-f94f-92b0-8f1d-80f57f35ce05.htm) | 构建读取数据的报文命令 |
| 公共方法静态成员 | [BuildWriteCommand](59521eec-5735-ef60-b7c4-f4707ecf9dec.htm) | 构建写入数据的报文命令 |
| 公共方法 | [Close](992f28b3-b99b-2196-274e-2f25e9bd3c1d.htm) | 关闭当前的串口连接  Close the current serial connection (继承自 [SerialBase](e877e20b-ef65-266c-05a5-5db58e6bad0c.htm)。) |
| 受保护的方法 | [DecideWhetherQAMessage](acc40cda-e310-8537-6685-59b67eeb16ff.htm) | 决定当前的消息是否是用于问答机制返回的消息，默认直接返回 true, 实际的情况需要根据协议进行重写方法  To determine whether the current message is the message returned by the question answering mechanism, the default is true. In actual cases, the rewriting method needs to be performed according to the protocol (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [Dispose](c83e3a27-3a09-f67c-9c7f-96b7a5a4c8b0.htm) | 释放当前的对象 (继承自 [SerialBase](e877e20b-ef65-266c-05a5-5db58e6bad0c.htm)。) |
| 受保护的方法 | [Dispose(Boolean)](c509106d-35a7-2ff2-7dd8-ca67e16d7950.htm) | 释放当前的对象 (继承自 [SerialBase](e877e20b-ef65-266c-05a5-5db58e6bad0c.htm)。) |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | [ExtraAfterReadFromCoreServer](06a49a15-a3ce-5e39-28c7-2b78851eca2e.htm) | 和服务器交互完成的时候调用的方法，可以根据读写结果进行一些额外的操作，具体的操作需要根据实际的需求来重写实现  The method called when the interaction with the server is completed can perform some additional operations based on the read and write results. The specific operations need to be rewritten according to actual needs. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法静态成员 | [ExtractActualData](9937e217-ca0b-0ea1-5e75-28be1a0af552.htm) | 把服务器反馈的数据解析成实际的命令 |
| 受保护的方法代码示例 | [ExtraOnDisconnect](de9dbba2-3b2c-1c5e-53c8-4e6f8a44e2c4.htm) | 根据实际的协议选择是否重写本方法，有些协议在断开连接之前，需要发送一些报文来关闭当前的网络通道  Select whether to rewrite this method according to the actual protocol. Some protocols need to send some packets to close the current network channel before disconnecting. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法代码示例 | [ExtraOnDisconnectAsync](210f593e-16d8-77c4-1178-58fbcfe9b27b.htm) | 根据实际的协议选择是否重写本方法，有些协议在断开连接之前，需要发送一些报文来关闭当前的网络通道  Select whether to rewrite this method according to the actual protocol. Some protocols need to send some packets to close the current network channel before disconnecting. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 受保护的方法 | [GetLogTextFromBinary](67f3b2f7-957c-fa0b-8320-6799237157db.htm) | 获取当前的报文进行日志记录的时候，是否使用二进制的格式记录，默认返回 [LogMsgFormatBinary](2ed2ff87-caec-0e44-f08b-e2b015225a58.htm)，重写可以根据session对象分别返回不同记录模式  Whether to log the current packet in binary format, the default return is [LogMsgFormatBinary](2ed2ff87-caec-0e44-f08b-e2b015225a58.htm). If you want to override it, different recording modes can be returned according to session (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [GetNewNetMessage](16636a1b-a6f0-9c0f-9d9c-c94e26c67459.htm) | 获取一个新的消息对象的方法，需要在继承类里面进行重写  The method to get a new message object needs to be overridden in the inheritance class (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
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
| 公共方法静态成员 | [PackCommand](1094f5da-f719-64c7-4ffc-c42c5f8004c7.htm) | 将命令和数据打包成用于发送的报文 |
| 公共方法 | [PackCommandWithHeader](d0b2cd76-c306-11fc-38c7-af8ad306cce7.htm) | 对当前的命令进行打包处理，通常是携带命令头内容，标记当前的命令的长度信息，需要进行重写，否则默认不打包  The current command is packaged, usually carrying the content of the command header, marking the length of the current command, and it needs to be rewritten, otherwise it is not packaged by default (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [Read](d4363a0b-2a6e-d4e2-4e40-7abb240b0382.htm) | 读取光源控制器的参数信息，需要传入通道号信息，读取到详细的内容参照[ShineInLightData](6981aaa1-3c66-69ba-49da-a707306a9267.htm)的值 |
| 公共方法代码示例 | [ReadFromCoreServer(Byte)](8146dc0c-1edc-3ca7-428a-a3bdbd35e6c5.htm) | 将二进制的数据发送到管道中去，然后从管道里接收二进制的数据回来，并返回是否成功的结果对象。  Send binary data to the pipeline, and then receive binary data back from the pipeline, and return whether the success of the result object (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServer(IEnumerableByte)](044cfdd9-ec76-b89f-2f2b-5f439025006c.htm) | 将多个数据报文按顺序发到设备，并从设备接收返回的数据内容，然后拼接成一个Byte[]信息，需要重写[UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm)方法才能返回正确的结果。  Send multiple data packets to the device in sequence, and receive the returned data content from the device, and then splicing them into a Byte[] message, you need to rewrite [UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm) method to return the correct result. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServer(Byte, Boolean, Boolean)](b32ce8e8-a959-fcc2-6736-14b39fdb9a40.htm) | 将二进制的数据发送到管道中去，然后从管道里接收二进制的数据回来，并返回是否成功的结果对象。  Send binary data to the pipeline, and then receive binary data back from the pipeline, and return whether the success of the result object (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServer(CommunicationPipe, Byte, Boolean, Boolean)](ef36d56b-df09-f727-e8d3-316d37a13ee9.htm) | 使用指定的管道来进行数据通信，发送原始数据到管道，然后从管道接收相关的数据返回，本方法无锁 (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServerAsync(Byte)](a4de5a24-c67a-7b02-db05-33b85eaf6e1e.htm) | 将二进制的数据发送到管道中去，然后从管道里接收二进制的数据回来，并返回是否成功的结果对象。  Send binary data to the pipeline, and then receive binary data back from the pipeline, and return whether the success of the result object (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServerAsync(IEnumerableByte)](bcdb59b7-a70a-d5c0-1fc6-9ececcd39ea2.htm) | 将多个数据报文按顺序发到设备，并从设备接收返回的数据内容，然后拼接成一个Byte[]信息，需要重写[UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm)方法才能返回正确的结果。  Send multiple data packets to the device in sequence, and receive the returned data content from the device, and then splicing them into a Byte[] message, you need to rewrite [UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm) method to return the correct result. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServerAsync(Byte, Boolean, Boolean)](1c79c204-5ab5-9342-e65d-7455c20ea577.htm) | 将二进制的数据发送到管道中去，然后从管道里接收二进制的数据回来，并返回是否成功的结果对象。  Send binary data to the pipeline, and then receive binary data back from the pipeline, and return whether the success of the result object (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServerAsync(CommunicationPipe, Byte, Boolean, Boolean)](a254a19d-3c29-147b-335b-8adccdf4ee2f.htm) | 使用指定的管道来进行数据通信，发送原始数据到管道，然后从管道接收相关的数据返回，本方法无锁 (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [SerialPortInni(String)](76fee18c-6275-b5cc-2f03-e9095178cae5.htm) | 初始化串口信息，57600波特率，8位数据位，1位停止位，偶校验  Initial serial port information, 57600 baud rate, 8 data bits, 1 stop bit, even parity (重写 [SerialBaseSerialPortInni(String)](7dff5e26-6594-291e-4792-b13db9c5d613.htm).) |
| 公共方法 | [SerialPortInni(ActionSerialPort)](53585409-6243-9359-1846-bf8b4e266708.htm) | 根据自定义初始化方法进行初始化串口信息  Initialize the serial port information according to the custom initialization method (继承自 [SerialBase](e877e20b-ef65-266c-05a5-5db58e6bad0c.htm)。) |
| 公共方法 | [SerialPortInni(String, Int32)](601606e3-f3c0-4bd7-49f3-cf420c123fb9.htm) | 初始化串口信息，波特率，8位数据位，1位停止位，偶校验  Initializes serial port information, baud rate, 8-bit data bit, 1-bit stop bit, even parity (重写 [SerialBaseSerialPortInni(String, Int32)](0c7aa163-1d81-54fd-f4f5-103a7e9a341d.htm).) |
| 公共方法 | [SerialPortInni(String, Int32, Int32, StopBits, Parity)](420c818b-48e9-fa02-59e2-f4d5897d6802.htm) | 初始化串口信息，波特率，数据位，停止位，奇偶校验需要全部自己来指定  Start serial port information, baud rate, data bit, stop bit, parity all need to be specified (继承自 [SerialBase](e877e20b-ef65-266c-05a5-5db58e6bad0c.htm)。) |
| 公共方法 | [SetDtuPipe](f765d1d7-234a-b229-9457-97e25e024dc9.htm) | 将当前的通信对象设置DTU模式，允许传入现成的管道，并返回初始化结果，如果该设备重写了握手报文，就是返回握手结果  Set the current communication object to DTU mode, allow the existing pipe to be passed in, and return the initialization result, if the device rewrites the handshake packet, the handshake result is returned (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [SetDtuPipeAsync](9d1fd12d-b804-3828-a627-a2c0eef1d23f.htm) | 将当前的通信对象设置DTU模式，允许传入现成的管道，并返回初始化结果，如果该设备重写了握手报文，就是返回握手结果  Set the current communication object to DTU mode, allow the existing pipe to be passed in, and return the initialization result, if the device rewrites the handshake packet, the handshake result is returned (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ToString](1154f133-182e-ecfb-ebc5-347911fc7506.htm) | (重写 [SerialBaseToString](1ac09324-4d19-863d-d683-cbe804ef2124.htm).) |
| 公共方法 | [UnpackResponseContent](208017a7-8da2-33f8-3273-d7e7e850100f.htm) | 根据对方返回的报文命令，对命令进行基本的拆包，例如各种Modbus协议拆包为统一的核心报文，还支持对报文的验证  According to the message command returned by the other party, the command is basically unpacked, for example, various Modbus protocols are unpacked into a unified core message, and the verification of the message is also supported (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [Write](9d0bc031-6b4d-20cd-62e9-f52cb3f4b0b6.htm) | 将光源控制器的数据写入到设备，返回是否写入成功 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[ShineInLightSourceController 类](5de4a57b-1202-4bcd-4650-a524ba7166df.htm)

[HslCommunication.Instrument.Light 命名空间](00d11148-efb6-2e43-ae78-c995e930d868.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## BuildReadCommand 方法 

[原文連結](http://api.hslcommunication.cn/html/148fde6a-f94f-92b0-8f1d-80f57f35ce05.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.Light](../html/00d11148-efb6-2e43-ae78-c995e930d868.htm "HslCommunication.Instrument.Light")

[ShineInLightSourceController 类](../html/5de4a57b-1202-4bcd-4650-a524ba7166df.htm "ShineInLightSourceController 类")

[ShineInLightSourceController 方法](../html/cbbaa7f0-2b3e-2557-2808-9aa825c888d7.htm "ShineInLightSourceController 方法")

[BuildReadCommand 方法](../html/148fde6a-f94f-92b0-8f1d-80f57f35ce05.htm "BuildReadCommand 方法 ")

[BuildWriteCommand 方法](../html/59521eec-5735-ef60-b7c4-f4707ecf9dec.htm "BuildWriteCommand 方法 ")

[ExtractActualData 方法](../html/9937e217-ca0b-0ea1-5e75-28be1a0af552.htm "ExtractActualData 方法 ")

[PackCommand 方法](../html/1094f5da-f719-64c7-4ffc-c42c5f8004c7.htm "PackCommand 方法 ")

[Read 方法](../html/d4363a0b-2a6e-d4e2-4e40-7abb240b0382.htm "Read 方法 ")

[SerialPortInni 方法](../html/b06cbc7a-09af-2f7e-78e8-a6015e79cd15.htm "SerialPortInni 方法 ")

[ToString 方法](../html/1154f133-182e-ecfb-ebc5-347911fc7506.htm "ToString 方法 ")

[Write 方法](../html/9d0bc031-6b4d-20cd-62e9-f52cb3f4b0b6.htm "Write 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ShineInLightSourceControllerBuildReadCommand 方法 |

构建读取数据的报文命令

**命名空间：**
 [HslCommunication.Instrument.Light](00d11148-efb6-2e43-ae78-c995e930d868.htm)  
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
	byte channel
)
```

```
Public Shared Function BuildReadCommand ( 
	channel As Byte
) As Byte()
```

```
public:
static array<unsigned char>^ BuildReadCommand(
	unsigned char channel
)
```

```
static member BuildReadCommand : 
        channel : byte -> byte[] 
```

#### 参数

channel
:   类型：SystemByte  
    通道信息

#### 返回值

类型：Byte  
构建读取的命令

![](../icons/SectionExpanded.png)参见

#### 引用

[ShineInLightSourceController 类](5de4a57b-1202-4bcd-4650-a524ba7166df.htm)

[HslCommunication.Instrument.Light 命名空间](00d11148-efb6-2e43-ae78-c995e930d868.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## BuildWriteCommand 方法 

[原文連結](http://api.hslcommunication.cn/html/59521eec-5735-ef60-b7c4-f4707ecf9dec.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.Light](../html/00d11148-efb6-2e43-ae78-c995e930d868.htm "HslCommunication.Instrument.Light")

[ShineInLightSourceController 类](../html/5de4a57b-1202-4bcd-4650-a524ba7166df.htm "ShineInLightSourceController 类")

[ShineInLightSourceController 方法](../html/cbbaa7f0-2b3e-2557-2808-9aa825c888d7.htm "ShineInLightSourceController 方法")

[BuildReadCommand 方法](../html/148fde6a-f94f-92b0-8f1d-80f57f35ce05.htm "BuildReadCommand 方法 ")

[BuildWriteCommand 方法](../html/59521eec-5735-ef60-b7c4-f4707ecf9dec.htm "BuildWriteCommand 方法 ")

[ExtractActualData 方法](../html/9937e217-ca0b-0ea1-5e75-28be1a0af552.htm "ExtractActualData 方法 ")

[PackCommand 方法](../html/1094f5da-f719-64c7-4ffc-c42c5f8004c7.htm "PackCommand 方法 ")

[Read 方法](../html/d4363a0b-2a6e-d4e2-4e40-7abb240b0382.htm "Read 方法 ")

[SerialPortInni 方法](../html/b06cbc7a-09af-2f7e-78e8-a6015e79cd15.htm "SerialPortInni 方法 ")

[ToString 方法](../html/1154f133-182e-ecfb-ebc5-347911fc7506.htm "ToString 方法 ")

[Write 方法](../html/9d0bc031-6b4d-20cd-62e9-f52cb3f4b0b6.htm "Write 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ShineInLightSourceControllerBuildWriteCommand 方法 |

构建写入数据的报文命令

**命名空间：**
 [HslCommunication.Instrument.Light](00d11148-efb6-2e43-ae78-c995e930d868.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static byte[] BuildWriteCommand(
	ShineInLightData shineInLightData
)
```

```
Public Shared Function BuildWriteCommand ( 
	shineInLightData As ShineInLightData
) As Byte()
```

```
public:
static array<unsigned char>^ BuildWriteCommand(
	ShineInLightData^ shineInLightData
)
```

```
static member BuildWriteCommand : 
        shineInLightData : ShineInLightData -> byte[] 
```

#### 参数

shineInLightData
:   类型：[HslCommunication.Instrument.LightShineInLightData](6981aaa1-3c66-69ba-49da-a707306a9267.htm)  
    准备写入的数据

#### 返回值

类型：Byte  
报文命令

![](../icons/SectionExpanded.png)参见

#### 引用

[ShineInLightSourceController 类](5de4a57b-1202-4bcd-4650-a524ba7166df.htm)

[HslCommunication.Instrument.Light 命名空间](00d11148-efb6-2e43-ae78-c995e930d868.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ExtractActualData 方法 

[原文連結](http://api.hslcommunication.cn/html/9937e217-ca0b-0ea1-5e75-28be1a0af552.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.Light](../html/00d11148-efb6-2e43-ae78-c995e930d868.htm "HslCommunication.Instrument.Light")

[ShineInLightSourceController 类](../html/5de4a57b-1202-4bcd-4650-a524ba7166df.htm "ShineInLightSourceController 类")

[ShineInLightSourceController 方法](../html/cbbaa7f0-2b3e-2557-2808-9aa825c888d7.htm "ShineInLightSourceController 方法")

[BuildReadCommand 方法](../html/148fde6a-f94f-92b0-8f1d-80f57f35ce05.htm "BuildReadCommand 方法 ")

[BuildWriteCommand 方法](../html/59521eec-5735-ef60-b7c4-f4707ecf9dec.htm "BuildWriteCommand 方法 ")

[ExtractActualData 方法](../html/9937e217-ca0b-0ea1-5e75-28be1a0af552.htm "ExtractActualData 方法 ")

[PackCommand 方法](../html/1094f5da-f719-64c7-4ffc-c42c5f8004c7.htm "PackCommand 方法 ")

[Read 方法](../html/d4363a0b-2a6e-d4e2-4e40-7abb240b0382.htm "Read 方法 ")

[SerialPortInni 方法](../html/b06cbc7a-09af-2f7e-78e8-a6015e79cd15.htm "SerialPortInni 方法 ")

[ToString 方法](../html/1154f133-182e-ecfb-ebc5-347911fc7506.htm "ToString 方法 ")

[Write 方法](../html/9d0bc031-6b4d-20cd-62e9-f52cb3f4b0b6.htm "Write 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ShineInLightSourceControllerExtractActualData 方法 |

把服务器反馈的数据解析成实际的命令

**命名空间：**
 [HslCommunication.Instrument.Light](00d11148-efb6-2e43-ae78-c995e930d868.htm)  
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
    反馈的数据

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
结果内容

![](../icons/SectionExpanded.png)参见

#### 引用

[ShineInLightSourceController 类](5de4a57b-1202-4bcd-4650-a524ba7166df.htm)

[HslCommunication.Instrument.Light 命名空间](00d11148-efb6-2e43-ae78-c995e930d868.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## PackCommand 方法 

[原文連結](http://api.hslcommunication.cn/html/1094f5da-f719-64c7-4ffc-c42c5f8004c7.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.Light](../html/00d11148-efb6-2e43-ae78-c995e930d868.htm "HslCommunication.Instrument.Light")

[ShineInLightSourceController 类](../html/5de4a57b-1202-4bcd-4650-a524ba7166df.htm "ShineInLightSourceController 类")

[ShineInLightSourceController 方法](../html/cbbaa7f0-2b3e-2557-2808-9aa825c888d7.htm "ShineInLightSourceController 方法")

[BuildReadCommand 方法](../html/148fde6a-f94f-92b0-8f1d-80f57f35ce05.htm "BuildReadCommand 方法 ")

[BuildWriteCommand 方法](../html/59521eec-5735-ef60-b7c4-f4707ecf9dec.htm "BuildWriteCommand 方法 ")

[ExtractActualData 方法](../html/9937e217-ca0b-0ea1-5e75-28be1a0af552.htm "ExtractActualData 方法 ")

[PackCommand 方法](../html/1094f5da-f719-64c7-4ffc-c42c5f8004c7.htm "PackCommand 方法 ")

[Read 方法](../html/d4363a0b-2a6e-d4e2-4e40-7abb240b0382.htm "Read 方法 ")

[SerialPortInni 方法](../html/b06cbc7a-09af-2f7e-78e8-a6015e79cd15.htm "SerialPortInni 方法 ")

[ToString 方法](../html/1154f133-182e-ecfb-ebc5-347911fc7506.htm "ToString 方法 ")

[Write 方法](../html/9d0bc031-6b4d-20cd-62e9-f52cb3f4b0b6.htm "Write 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ShineInLightSourceControllerPackCommand 方法 |

将命令和数据打包成用于发送的报文

**命名空间：**
 [HslCommunication.Instrument.Light](00d11148-efb6-2e43-ae78-c995e930d868.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static byte[] PackCommand(
	byte cmd,
	byte[] data
)
```

```
Public Shared Function PackCommand ( 
	cmd As Byte,
	data As Byte()
) As Byte()
```

```
public:
static array<unsigned char>^ PackCommand(
	unsigned char cmd, 
	array<unsigned char>^ data
)
```

```
static member PackCommand : 
        cmd : byte * 
        data : byte[] -> byte[] 
```

#### 参数

cmd
:   类型：SystemByte  
    命令

data
:   类型：SystemByte  
    命令数据

#### 返回值

类型：Byte  
可用于发送的报文

![](../icons/SectionExpanded.png)参见

#### 引用

[ShineInLightSourceController 类](5de4a57b-1202-4bcd-4650-a524ba7166df.htm)

[HslCommunication.Instrument.Light 命名空间](00d11148-efb6-2e43-ae78-c995e930d868.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Read 方法 

[原文連結](http://api.hslcommunication.cn/html/d4363a0b-2a6e-d4e2-4e40-7abb240b0382.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.Light](../html/00d11148-efb6-2e43-ae78-c995e930d868.htm "HslCommunication.Instrument.Light")

[ShineInLightSourceController 类](../html/5de4a57b-1202-4bcd-4650-a524ba7166df.htm "ShineInLightSourceController 类")

[ShineInLightSourceController 方法](../html/cbbaa7f0-2b3e-2557-2808-9aa825c888d7.htm "ShineInLightSourceController 方法")

[BuildReadCommand 方法](../html/148fde6a-f94f-92b0-8f1d-80f57f35ce05.htm "BuildReadCommand 方法 ")

[BuildWriteCommand 方法](../html/59521eec-5735-ef60-b7c4-f4707ecf9dec.htm "BuildWriteCommand 方法 ")

[ExtractActualData 方法](../html/9937e217-ca0b-0ea1-5e75-28be1a0af552.htm "ExtractActualData 方法 ")

[PackCommand 方法](../html/1094f5da-f719-64c7-4ffc-c42c5f8004c7.htm "PackCommand 方法 ")

[Read 方法](../html/d4363a0b-2a6e-d4e2-4e40-7abb240b0382.htm "Read 方法 ")

[SerialPortInni 方法](../html/b06cbc7a-09af-2f7e-78e8-a6015e79cd15.htm "SerialPortInni 方法 ")

[ToString 方法](../html/1154f133-182e-ecfb-ebc5-347911fc7506.htm "ToString 方法 ")

[Write 方法](../html/9d0bc031-6b4d-20cd-62e9-f52cb3f4b0b6.htm "Write 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ShineInLightSourceControllerRead 方法 |

读取光源控制器的参数信息，需要传入通道号信息，读取到详细的内容参照[ShineInLightData](6981aaa1-3c66-69ba-49da-a707306a9267.htm)的值

**命名空间：**
 [HslCommunication.Instrument.Light](00d11148-efb6-2e43-ae78-c995e930d868.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public OperateResult<ShineInLightData> Read(
	byte channel
)
```

```
Public Function Read ( 
	channel As Byte
) As OperateResult(Of ShineInLightData)
```

```
public:
OperateResult<ShineInLightData^>^ Read(
	unsigned char channel
)
```

```
member Read : 
        channel : byte -> OperateResult<ShineInLightData> 
```

#### 参数

channel
:   类型：SystemByte  
    读取的通道信息

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)[ShineInLightData](6981aaa1-3c66-69ba-49da-a707306a9267.htm)  
读取的参数值

![](../icons/SectionExpanded.png)参见

#### 引用

[ShineInLightSourceController 类](5de4a57b-1202-4bcd-4650-a524ba7166df.htm)

[HslCommunication.Instrument.Light 命名空间](00d11148-efb6-2e43-ae78-c995e930d868.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## SerialPortInni 方法 

[原文連結](http://api.hslcommunication.cn/html/b06cbc7a-09af-2f7e-78e8-a6015e79cd15.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.Light](../html/00d11148-efb6-2e43-ae78-c995e930d868.htm "HslCommunication.Instrument.Light")

[ShineInLightSourceController 类](../html/5de4a57b-1202-4bcd-4650-a524ba7166df.htm "ShineInLightSourceController 类")

[ShineInLightSourceController 方法](../html/cbbaa7f0-2b3e-2557-2808-9aa825c888d7.htm "ShineInLightSourceController 方法")

[SerialPortInni 方法](../html/b06cbc7a-09af-2f7e-78e8-a6015e79cd15.htm "SerialPortInni 方法 ")

[SerialPortInni 方法 (String)](../html/76fee18c-6275-b5cc-2f03-e9095178cae5.htm "SerialPortInni 方法 (String)")

[SerialPortInni 方法 (String, Int32)](../html/601606e3-f3c0-4bd7-49f3-cf420c123fb9.htm "SerialPortInni 方法 (String, Int32)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ShineInLightSourceControllerSerialPortInni 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [SerialPortInni(String)](76fee18c-6275-b5cc-2f03-e9095178cae5.htm) | 初始化串口信息，57600波特率，8位数据位，1位停止位，偶校验  Initial serial port information, 57600 baud rate, 8 data bits, 1 stop bit, even parity (重写 [SerialBaseSerialPortInni(String)](7dff5e26-6594-291e-4792-b13db9c5d613.htm).) |
| 公共方法 | [SerialPortInni(ActionSerialPort)](53585409-6243-9359-1846-bf8b4e266708.htm) | 根据自定义初始化方法进行初始化串口信息  Initialize the serial port information according to the custom initialization method (继承自 [SerialBase](e877e20b-ef65-266c-05a5-5db58e6bad0c.htm)。) |
| 公共方法 | [SerialPortInni(String, Int32)](601606e3-f3c0-4bd7-49f3-cf420c123fb9.htm) | 初始化串口信息，波特率，8位数据位，1位停止位，偶校验  Initializes serial port information, baud rate, 8-bit data bit, 1-bit stop bit, even parity (重写 [SerialBaseSerialPortInni(String, Int32)](0c7aa163-1d81-54fd-f4f5-103a7e9a341d.htm).) |
| 公共方法 | [SerialPortInni(String, Int32, Int32, StopBits, Parity)](420c818b-48e9-fa02-59e2-f4d5897d6802.htm) | 初始化串口信息，波特率，数据位，停止位，奇偶校验需要全部自己来指定  Start serial port information, baud rate, data bit, stop bit, parity all need to be specified (继承自 [SerialBase](e877e20b-ef65-266c-05a5-5db58e6bad0c.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[ShineInLightSourceController 类](5de4a57b-1202-4bcd-4650-a524ba7166df.htm)

[HslCommunication.Instrument.Light 命名空间](00d11148-efb6-2e43-ae78-c995e930d868.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## SerialPortInni 方法 (String)

[原文連結](http://api.hslcommunication.cn/html/76fee18c-6275-b5cc-2f03-e9095178cae5.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.Light](../html/00d11148-efb6-2e43-ae78-c995e930d868.htm "HslCommunication.Instrument.Light")

[ShineInLightSourceController 类](../html/5de4a57b-1202-4bcd-4650-a524ba7166df.htm "ShineInLightSourceController 类")

[ShineInLightSourceController 方法](../html/cbbaa7f0-2b3e-2557-2808-9aa825c888d7.htm "ShineInLightSourceController 方法")

[SerialPortInni 方法](../html/b06cbc7a-09af-2f7e-78e8-a6015e79cd15.htm "SerialPortInni 方法 ")

[SerialPortInni 方法 (String)](../html/76fee18c-6275-b5cc-2f03-e9095178cae5.htm "SerialPortInni 方法 (String)")

[SerialPortInni 方法 (String, Int32)](../html/601606e3-f3c0-4bd7-49f3-cf420c123fb9.htm "SerialPortInni 方法 (String, Int32)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ShineInLightSourceControllerSerialPortInni 方法 (String) |

初始化串口信息，57600波特率，8位数据位，1位停止位，偶校验  
Initial serial port information, 57600 baud rate, 8 data bits, 1 stop bit, even parity

**命名空间：**
 [HslCommunication.Instrument.Light](00d11148-efb6-2e43-ae78-c995e930d868.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public override void SerialPortInni(
	string portName
)
```

```
Public Overrides Sub SerialPortInni ( 
	portName As String
)
```

```
public:
virtual void SerialPortInni(
	String^ portName
) override
```

```
abstract SerialPortInni : 
        portName : string -> unit 
override SerialPortInni : 
        portName : string -> unit
```

#### 参数

portName
:   类型：SystemString  
    端口号信息，例如"COM3"

![](../icons/SectionExpanded.png)参见

#### 引用

[ShineInLightSourceController 类](5de4a57b-1202-4bcd-4650-a524ba7166df.htm)

[SerialPortInni 重载](b06cbc7a-09af-2f7e-78e8-a6015e79cd15.htm)

[HslCommunication.Instrument.Light 命名空间](00d11148-efb6-2e43-ae78-c995e930d868.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## SerialPortInni 方法 (String, Int32)

[原文連結](http://api.hslcommunication.cn/html/601606e3-f3c0-4bd7-49f3-cf420c123fb9.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.Light](../html/00d11148-efb6-2e43-ae78-c995e930d868.htm "HslCommunication.Instrument.Light")

[ShineInLightSourceController 类](../html/5de4a57b-1202-4bcd-4650-a524ba7166df.htm "ShineInLightSourceController 类")

[ShineInLightSourceController 方法](../html/cbbaa7f0-2b3e-2557-2808-9aa825c888d7.htm "ShineInLightSourceController 方法")

[SerialPortInni 方法](../html/b06cbc7a-09af-2f7e-78e8-a6015e79cd15.htm "SerialPortInni 方法 ")

[SerialPortInni 方法 (String)](../html/76fee18c-6275-b5cc-2f03-e9095178cae5.htm "SerialPortInni 方法 (String)")

[SerialPortInni 方法 (String, Int32)](../html/601606e3-f3c0-4bd7-49f3-cf420c123fb9.htm "SerialPortInni 方法 (String, Int32)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ShineInLightSourceControllerSerialPortInni 方法 (String, Int32) |

初始化串口信息，波特率，8位数据位，1位停止位，偶校验  
Initializes serial port information, baud rate, 8-bit data bit, 1-bit stop bit, even parity

**命名空间：**
 [HslCommunication.Instrument.Light](00d11148-efb6-2e43-ae78-c995e930d868.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public override void SerialPortInni(
	string portName,
	int baudRate
)
```

```
Public Overrides Sub SerialPortInni ( 
	portName As String,
	baudRate As Integer
)
```

```
public:
virtual void SerialPortInni(
	String^ portName, 
	int baudRate
) override
```

```
abstract SerialPortInni : 
        portName : string * 
        baudRate : int -> unit 
override SerialPortInni : 
        portName : string * 
        baudRate : int -> unit
```

#### 参数

portName
:   类型：SystemString  
    端口号信息，例如"COM3"

baudRate
:   类型：SystemInt32  
    波特率

![](../icons/SectionExpanded.png)参见

#### 引用

[ShineInLightSourceController 类](5de4a57b-1202-4bcd-4650-a524ba7166df.htm)

[SerialPortInni 重载](b06cbc7a-09af-2f7e-78e8-a6015e79cd15.htm)

[HslCommunication.Instrument.Light 命名空间](00d11148-efb6-2e43-ae78-c995e930d868.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ToString 方法 

[原文連結](http://api.hslcommunication.cn/html/1154f133-182e-ecfb-ebc5-347911fc7506.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.Light](../html/00d11148-efb6-2e43-ae78-c995e930d868.htm "HslCommunication.Instrument.Light")

[ShineInLightSourceController 类](../html/5de4a57b-1202-4bcd-4650-a524ba7166df.htm "ShineInLightSourceController 类")

[ShineInLightSourceController 方法](../html/cbbaa7f0-2b3e-2557-2808-9aa825c888d7.htm "ShineInLightSourceController 方法")

[BuildReadCommand 方法](../html/148fde6a-f94f-92b0-8f1d-80f57f35ce05.htm "BuildReadCommand 方法 ")

[BuildWriteCommand 方法](../html/59521eec-5735-ef60-b7c4-f4707ecf9dec.htm "BuildWriteCommand 方法 ")

[ExtractActualData 方法](../html/9937e217-ca0b-0ea1-5e75-28be1a0af552.htm "ExtractActualData 方法 ")

[PackCommand 方法](../html/1094f5da-f719-64c7-4ffc-c42c5f8004c7.htm "PackCommand 方法 ")

[Read 方法](../html/d4363a0b-2a6e-d4e2-4e40-7abb240b0382.htm "Read 方法 ")

[SerialPortInni 方法](../html/b06cbc7a-09af-2f7e-78e8-a6015e79cd15.htm "SerialPortInni 方法 ")

[ToString 方法](../html/1154f133-182e-ecfb-ebc5-347911fc7506.htm "ToString 方法 ")

[Write 方法](../html/9d0bc031-6b4d-20cd-62e9-f52cb3f4b0b6.htm "Write 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ShineInLightSourceControllerToString 方法 |

[缺少 "M:HslCommunication.Instrument.Light.ShineInLightSourceController.ToString" 的 <summary> 文档]

**命名空间：**
 [HslCommunication.Instrument.Light](00d11148-efb6-2e43-ae78-c995e930d868.htm)  
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

[缺少 "M:HslCommunication.Instrument.Light.ShineInLightSourceController.ToString" 的 <returns> 文档]

![](../icons/SectionExpanded.png)参见

#### 引用

[ShineInLightSourceController 类](5de4a57b-1202-4bcd-4650-a524ba7166df.htm)

[HslCommunication.Instrument.Light 命名空间](00d11148-efb6-2e43-ae78-c995e930d868.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)