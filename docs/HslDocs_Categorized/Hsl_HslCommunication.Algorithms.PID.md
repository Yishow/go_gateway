# HslCommunication - HslCommunication.Algorithms.PID

> 分類頁數: 14



---
## HslCommunication.Algorithms.PID

[原文連結](http://api.hslcommunication.cn/html/e119de70-6089-f087-5e46-755795fa33c0.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Algorithms.PID](../html/e119de70-6089-f087-5e46-755795fa33c0.htm "HslCommunication.Algorithms.PID")

[PIDHelper 类](../html/eaa06577-fe3f-57e1-0cf7-e906f4587e04.htm "PIDHelper 类")

[PidMode 枚举](../html/88867bf9-a36f-5029-1afc-91646e44b13a.htm "PidMode 枚举")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslCommunication.Algorithms.PID 命名空间 |

[缺少 "N:HslCommunication.Algorithms.PID" 的 <summary> 文档]

![](../icons/SectionExpanded.png)类

|  | 类 | 说明 |
| --- | --- | --- |
| 公共类 | [PIDHelper](eaa06577-fe3f-57e1-0cf7-e906f4587e04.htm) | 一个PID的辅助类，可以设置 P,I,D 三者的值，用来模拟信号波动的时候，信号的收敛情况 |

![](../icons/SectionExpanded.png)枚举

|  | 枚举 | 说明 |
| --- | --- | --- |
| 公共枚举 | [PidMode](88867bf9-a36f-5029-1afc-91646e44b13a.htm) | Pid的模式选择 |

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## PIDHelper 类

[原文連結](http://api.hslcommunication.cn/html/eaa06577-fe3f-57e1-0cf7-e906f4587e04.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Algorithms.PID](../html/e119de70-6089-f087-5e46-755795fa33c0.htm "HslCommunication.Algorithms.PID")

[PIDHelper 类](../html/eaa06577-fe3f-57e1-0cf7-e906f4587e04.htm "PIDHelper 类")

[PIDHelper 构造函数](../html/e73d2f2a-209d-896d-f108-33af9958bb5e.htm "PIDHelper 构造函数 ")

[PIDHelper 属性](../html/5ea37cfa-1260-b972-3a8f-7401785d1974.htm "PIDHelper 属性")

[PIDHelper 方法](../html/4e903b1d-3732-4465-2e98-35718af9acab.htm "PIDHelper 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| PIDHelper 类 |

一个PID的辅助类，可以设置 P,I,D 三者的值，用来模拟信号波动的时候，信号的收敛情况

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  HslCommunication.Algorithms.PIDPIDHelper

**命名空间：**
 [HslCommunication.Algorithms.PID](e119de70-6089-f087-5e46-755795fa33c0.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class PIDHelper
```

```
Public Class PIDHelper
```

```
public ref class PIDHelper
```

```
type PIDHelper =  class end
```

PIDHelper 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [PIDHelper](e73d2f2a-209d-896d-f108-33af9958bb5e.htm) | 实例化一个默认的对象 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [DeadBand](7b622a95-f3df-5f68-fcb4-1e16d435524c.htm) | 获取或设置死区的值 |
| 公共属性 | [Kd](e35edbee-1a85-8fe7-a2e4-c2e277dec003.htm) | 微分的参数信息 |
| 公共属性 | [Ki](2aeb7f72-333a-c560-5820-23d55a9d5273.htm) | 积分的参数信息 |
| 公共属性 | [Kp](ec4edeae-f643-e04e-7d0c-96cccf1c6dc6.htm) | -rando 比例的参数信息 |
| 公共属性 | [MaxLimit](4e794846-ffe3-ca9e-c9a1-a4c34f2672ce.htm) | 获取或设置输出的上限，默认为没有设置 |
| 公共属性 | [MinLimit](5363d5cc-63eb-803f-b18a-e5100b52aa03.htm) | 获取或设置输出的下限，默认为没有设置 |
| 公共属性 | [SetValue](fbfe17cc-dc49-58eb-b5e4-bfeb2104c873.htm) | 获取或设置当前设置的值 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法 | [PidCalculate](cf87f2b3-ce97-1b6c-f0b9-bbea07ca7d2e.htm) | 计算Pid数据的值 |
| 公共方法 | ToString | (继承自 Object。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.Algorithms.PID 命名空间](e119de70-6089-f087-5e46-755795fa33c0.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## PIDHelper 构造函数 

[原文連結](http://api.hslcommunication.cn/html/e73d2f2a-209d-896d-f108-33af9958bb5e.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Algorithms.PID](../html/e119de70-6089-f087-5e46-755795fa33c0.htm "HslCommunication.Algorithms.PID")

[PIDHelper 类](../html/eaa06577-fe3f-57e1-0cf7-e906f4587e04.htm "PIDHelper 类")

[PIDHelper 构造函数](../html/e73d2f2a-209d-896d-f108-33af9958bb5e.htm "PIDHelper 构造函数 ")

[PIDHelper 属性](../html/5ea37cfa-1260-b972-3a8f-7401785d1974.htm "PIDHelper 属性")

[PIDHelper 方法](../html/4e903b1d-3732-4465-2e98-35718af9acab.htm "PIDHelper 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| PIDHelper 构造函数 |

实例化一个默认的对象

**命名空间：**
 [HslCommunication.Algorithms.PID](e119de70-6089-f087-5e46-755795fa33c0.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public PIDHelper()
```

```
Public Sub New
```

```
public:
PIDHelper()
```

```
new : unit -> PIDHelper
```

![](../icons/SectionExpanded.png)参见

#### 引用

[PIDHelper 类](eaa06577-fe3f-57e1-0cf7-e906f4587e04.htm)

[HslCommunication.Algorithms.PID 命名空间](e119de70-6089-f087-5e46-755795fa33c0.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## PIDHelper 属性

[原文連結](http://api.hslcommunication.cn/html/5ea37cfa-1260-b972-3a8f-7401785d1974.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Algorithms.PID](../html/e119de70-6089-f087-5e46-755795fa33c0.htm "HslCommunication.Algorithms.PID")

[PIDHelper 类](../html/eaa06577-fe3f-57e1-0cf7-e906f4587e04.htm "PIDHelper 类")

[PIDHelper 属性](../html/5ea37cfa-1260-b972-3a8f-7401785d1974.htm "PIDHelper 属性")

[DeadBand 属性](../html/7b622a95-f3df-5f68-fcb4-1e16d435524c.htm "DeadBand 属性 ")

[Kd 属性](../html/e35edbee-1a85-8fe7-a2e4-c2e277dec003.htm "Kd 属性 ")

[Ki 属性](../html/2aeb7f72-333a-c560-5820-23d55a9d5273.htm "Ki 属性 ")

[Kp 属性](../html/ec4edeae-f643-e04e-7d0c-96cccf1c6dc6.htm "Kp 属性 ")

[MaxLimit 属性](../html/4e794846-ffe3-ca9e-c9a1-a4c34f2672ce.htm "MaxLimit 属性 ")

[MinLimit 属性](../html/5363d5cc-63eb-803f-b18a-e5100b52aa03.htm "MinLimit 属性 ")

[SetValue 属性](../html/fbfe17cc-dc49-58eb-b5e4-bfeb2104c873.htm "SetValue 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| PIDHelper 属性 |

[PIDHelper](eaa06577-fe3f-57e1-0cf7-e906f4587e04.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [DeadBand](7b622a95-f3df-5f68-fcb4-1e16d435524c.htm) | 获取或设置死区的值 |
| 公共属性 | [Kd](e35edbee-1a85-8fe7-a2e4-c2e277dec003.htm) | 微分的参数信息 |
| 公共属性 | [Ki](2aeb7f72-333a-c560-5820-23d55a9d5273.htm) | 积分的参数信息 |
| 公共属性 | [Kp](ec4edeae-f643-e04e-7d0c-96cccf1c6dc6.htm) | -rando 比例的参数信息 |
| 公共属性 | [MaxLimit](4e794846-ffe3-ca9e-c9a1-a4c34f2672ce.htm) | 获取或设置输出的上限，默认为没有设置 |
| 公共属性 | [MinLimit](5363d5cc-63eb-803f-b18a-e5100b52aa03.htm) | 获取或设置输出的下限，默认为没有设置 |
| 公共属性 | [SetValue](fbfe17cc-dc49-58eb-b5e4-bfeb2104c873.htm) | 获取或设置当前设置的值 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[PIDHelper 类](eaa06577-fe3f-57e1-0cf7-e906f4587e04.htm)

[HslCommunication.Algorithms.PID 命名空间](e119de70-6089-f087-5e46-755795fa33c0.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## DeadBand 属性 

[原文連結](http://api.hslcommunication.cn/html/7b622a95-f3df-5f68-fcb4-1e16d435524c.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Algorithms.PID](../html/e119de70-6089-f087-5e46-755795fa33c0.htm "HslCommunication.Algorithms.PID")

[PIDHelper 类](../html/eaa06577-fe3f-57e1-0cf7-e906f4587e04.htm "PIDHelper 类")

[PIDHelper 属性](../html/5ea37cfa-1260-b972-3a8f-7401785d1974.htm "PIDHelper 属性")

[DeadBand 属性](../html/7b622a95-f3df-5f68-fcb4-1e16d435524c.htm "DeadBand 属性 ")

[Kd 属性](../html/e35edbee-1a85-8fe7-a2e4-c2e277dec003.htm "Kd 属性 ")

[Ki 属性](../html/2aeb7f72-333a-c560-5820-23d55a9d5273.htm "Ki 属性 ")

[Kp 属性](../html/ec4edeae-f643-e04e-7d0c-96cccf1c6dc6.htm "Kp 属性 ")

[MaxLimit 属性](../html/4e794846-ffe3-ca9e-c9a1-a4c34f2672ce.htm "MaxLimit 属性 ")

[MinLimit 属性](../html/5363d5cc-63eb-803f-b18a-e5100b52aa03.htm "MinLimit 属性 ")

[SetValue 属性](../html/fbfe17cc-dc49-58eb-b5e4-bfeb2104c873.htm "SetValue 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| PIDHelperDeadBand 属性 |

获取或设置死区的值

**命名空间：**
 [HslCommunication.Algorithms.PID](e119de70-6089-f087-5e46-755795fa33c0.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public double DeadBand { get; set; }
```

```
Public Property DeadBand As Double
	Get
	Set
```

```
public:
property double DeadBand {
	double get ();
	void set (double value);
}
```

```
member DeadBand : float with get, set
```

#### 属性值

类型：Double

![](../icons/SectionExpanded.png)参见

#### 引用

[PIDHelper 类](eaa06577-fe3f-57e1-0cf7-e906f4587e04.htm)

[HslCommunication.Algorithms.PID 命名空间](e119de70-6089-f087-5e46-755795fa33c0.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Kd 属性 

[原文連結](http://api.hslcommunication.cn/html/e35edbee-1a85-8fe7-a2e4-c2e277dec003.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Algorithms.PID](../html/e119de70-6089-f087-5e46-755795fa33c0.htm "HslCommunication.Algorithms.PID")

[PIDHelper 类](../html/eaa06577-fe3f-57e1-0cf7-e906f4587e04.htm "PIDHelper 类")

[PIDHelper 属性](../html/5ea37cfa-1260-b972-3a8f-7401785d1974.htm "PIDHelper 属性")

[DeadBand 属性](../html/7b622a95-f3df-5f68-fcb4-1e16d435524c.htm "DeadBand 属性 ")

[Kd 属性](../html/e35edbee-1a85-8fe7-a2e4-c2e277dec003.htm "Kd 属性 ")

[Ki 属性](../html/2aeb7f72-333a-c560-5820-23d55a9d5273.htm "Ki 属性 ")

[Kp 属性](../html/ec4edeae-f643-e04e-7d0c-96cccf1c6dc6.htm "Kp 属性 ")

[MaxLimit 属性](../html/4e794846-ffe3-ca9e-c9a1-a4c34f2672ce.htm "MaxLimit 属性 ")

[MinLimit 属性](../html/5363d5cc-63eb-803f-b18a-e5100b52aa03.htm "MinLimit 属性 ")

[SetValue 属性](../html/fbfe17cc-dc49-58eb-b5e4-bfeb2104c873.htm "SetValue 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| PIDHelperKd 属性 |

微分的参数信息

**命名空间：**
 [HslCommunication.Algorithms.PID](e119de70-6089-f087-5e46-755795fa33c0.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public double Kd { get; set; }
```

```
Public Property Kd As Double
	Get
	Set
```

```
public:
property double Kd {
	double get ();
	void set (double value);
}
```

```
member Kd : float with get, set
```

#### 属性值

类型：Double

![](../icons/SectionExpanded.png)参见

#### 引用

[PIDHelper 类](eaa06577-fe3f-57e1-0cf7-e906f4587e04.htm)

[HslCommunication.Algorithms.PID 命名空间](e119de70-6089-f087-5e46-755795fa33c0.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Ki 属性 

[原文連結](http://api.hslcommunication.cn/html/2aeb7f72-333a-c560-5820-23d55a9d5273.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Algorithms.PID](../html/e119de70-6089-f087-5e46-755795fa33c0.htm "HslCommunication.Algorithms.PID")

[PIDHelper 类](../html/eaa06577-fe3f-57e1-0cf7-e906f4587e04.htm "PIDHelper 类")

[PIDHelper 属性](../html/5ea37cfa-1260-b972-3a8f-7401785d1974.htm "PIDHelper 属性")

[DeadBand 属性](../html/7b622a95-f3df-5f68-fcb4-1e16d435524c.htm "DeadBand 属性 ")

[Kd 属性](../html/e35edbee-1a85-8fe7-a2e4-c2e277dec003.htm "Kd 属性 ")

[Ki 属性](../html/2aeb7f72-333a-c560-5820-23d55a9d5273.htm "Ki 属性 ")

[Kp 属性](../html/ec4edeae-f643-e04e-7d0c-96cccf1c6dc6.htm "Kp 属性 ")

[MaxLimit 属性](../html/4e794846-ffe3-ca9e-c9a1-a4c34f2672ce.htm "MaxLimit 属性 ")

[MinLimit 属性](../html/5363d5cc-63eb-803f-b18a-e5100b52aa03.htm "MinLimit 属性 ")

[SetValue 属性](../html/fbfe17cc-dc49-58eb-b5e4-bfeb2104c873.htm "SetValue 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| PIDHelperKi 属性 |

积分的参数信息

**命名空间：**
 [HslCommunication.Algorithms.PID](e119de70-6089-f087-5e46-755795fa33c0.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public double Ki { get; set; }
```

```
Public Property Ki As Double
	Get
	Set
```

```
public:
property double Ki {
	double get ();
	void set (double value);
}
```

```
member Ki : float with get, set
```

#### 属性值

类型：Double

![](../icons/SectionExpanded.png)参见

#### 引用

[PIDHelper 类](eaa06577-fe3f-57e1-0cf7-e906f4587e04.htm)

[HslCommunication.Algorithms.PID 命名空间](e119de70-6089-f087-5e46-755795fa33c0.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Kp 属性 

[原文連結](http://api.hslcommunication.cn/html/ec4edeae-f643-e04e-7d0c-96cccf1c6dc6.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Algorithms.PID](../html/e119de70-6089-f087-5e46-755795fa33c0.htm "HslCommunication.Algorithms.PID")

[PIDHelper 类](../html/eaa06577-fe3f-57e1-0cf7-e906f4587e04.htm "PIDHelper 类")

[PIDHelper 属性](../html/5ea37cfa-1260-b972-3a8f-7401785d1974.htm "PIDHelper 属性")

[DeadBand 属性](../html/7b622a95-f3df-5f68-fcb4-1e16d435524c.htm "DeadBand 属性 ")

[Kd 属性](../html/e35edbee-1a85-8fe7-a2e4-c2e277dec003.htm "Kd 属性 ")

[Ki 属性](../html/2aeb7f72-333a-c560-5820-23d55a9d5273.htm "Ki 属性 ")

[Kp 属性](../html/ec4edeae-f643-e04e-7d0c-96cccf1c6dc6.htm "Kp 属性 ")

[MaxLimit 属性](../html/4e794846-ffe3-ca9e-c9a1-a4c34f2672ce.htm "MaxLimit 属性 ")

[MinLimit 属性](../html/5363d5cc-63eb-803f-b18a-e5100b52aa03.htm "MinLimit 属性 ")

[SetValue 属性](../html/fbfe17cc-dc49-58eb-b5e4-bfeb2104c873.htm "SetValue 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| PIDHelperKp 属性 |

-rando
比例的参数信息

**命名空间：**
 [HslCommunication.Algorithms.PID](e119de70-6089-f087-5e46-755795fa33c0.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public double Kp { get; set; }
```

```
Public Property Kp As Double
	Get
	Set
```

```
public:
property double Kp {
	double get ();
	void set (double value);
}
```

```
member Kp : float with get, set
```

#### 属性值

类型：Double

![](../icons/SectionExpanded.png)参见

#### 引用

[PIDHelper 类](eaa06577-fe3f-57e1-0cf7-e906f4587e04.htm)

[HslCommunication.Algorithms.PID 命名空间](e119de70-6089-f087-5e46-755795fa33c0.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## MaxLimit 属性 

[原文連結](http://api.hslcommunication.cn/html/4e794846-ffe3-ca9e-c9a1-a4c34f2672ce.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Algorithms.PID](../html/e119de70-6089-f087-5e46-755795fa33c0.htm "HslCommunication.Algorithms.PID")

[PIDHelper 类](../html/eaa06577-fe3f-57e1-0cf7-e906f4587e04.htm "PIDHelper 类")

[PIDHelper 属性](../html/5ea37cfa-1260-b972-3a8f-7401785d1974.htm "PIDHelper 属性")

[DeadBand 属性](../html/7b622a95-f3df-5f68-fcb4-1e16d435524c.htm "DeadBand 属性 ")

[Kd 属性](../html/e35edbee-1a85-8fe7-a2e4-c2e277dec003.htm "Kd 属性 ")

[Ki 属性](../html/2aeb7f72-333a-c560-5820-23d55a9d5273.htm "Ki 属性 ")

[Kp 属性](../html/ec4edeae-f643-e04e-7d0c-96cccf1c6dc6.htm "Kp 属性 ")

[MaxLimit 属性](../html/4e794846-ffe3-ca9e-c9a1-a4c34f2672ce.htm "MaxLimit 属性 ")

[MinLimit 属性](../html/5363d5cc-63eb-803f-b18a-e5100b52aa03.htm "MinLimit 属性 ")

[SetValue 属性](../html/fbfe17cc-dc49-58eb-b5e4-bfeb2104c873.htm "SetValue 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| PIDHelperMaxLimit 属性 |

获取或设置输出的上限，默认为没有设置

**命名空间：**
 [HslCommunication.Algorithms.PID](e119de70-6089-f087-5e46-755795fa33c0.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public double MaxLimit { get; set; }
```

```
Public Property MaxLimit As Double
	Get
	Set
```

```
public:
property double MaxLimit {
	double get ();
	void set (double value);
}
```

```
member MaxLimit : float with get, set
```

#### 属性值

类型：Double

![](../icons/SectionExpanded.png)参见

#### 引用

[PIDHelper 类](eaa06577-fe3f-57e1-0cf7-e906f4587e04.htm)

[HslCommunication.Algorithms.PID 命名空间](e119de70-6089-f087-5e46-755795fa33c0.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## MinLimit 属性 

[原文連結](http://api.hslcommunication.cn/html/5363d5cc-63eb-803f-b18a-e5100b52aa03.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Algorithms.PID](../html/e119de70-6089-f087-5e46-755795fa33c0.htm "HslCommunication.Algorithms.PID")

[PIDHelper 类](../html/eaa06577-fe3f-57e1-0cf7-e906f4587e04.htm "PIDHelper 类")

[PIDHelper 属性](../html/5ea37cfa-1260-b972-3a8f-7401785d1974.htm "PIDHelper 属性")

[DeadBand 属性](../html/7b622a95-f3df-5f68-fcb4-1e16d435524c.htm "DeadBand 属性 ")

[Kd 属性](../html/e35edbee-1a85-8fe7-a2e4-c2e277dec003.htm "Kd 属性 ")

[Ki 属性](../html/2aeb7f72-333a-c560-5820-23d55a9d5273.htm "Ki 属性 ")

[Kp 属性](../html/ec4edeae-f643-e04e-7d0c-96cccf1c6dc6.htm "Kp 属性 ")

[MaxLimit 属性](../html/4e794846-ffe3-ca9e-c9a1-a4c34f2672ce.htm "MaxLimit 属性 ")

[MinLimit 属性](../html/5363d5cc-63eb-803f-b18a-e5100b52aa03.htm "MinLimit 属性 ")

[SetValue 属性](../html/fbfe17cc-dc49-58eb-b5e4-bfeb2104c873.htm "SetValue 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| PIDHelperMinLimit 属性 |

获取或设置输出的下限，默认为没有设置

**命名空间：**
 [HslCommunication.Algorithms.PID](e119de70-6089-f087-5e46-755795fa33c0.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public double MinLimit { get; set; }
```

```
Public Property MinLimit As Double
	Get
	Set
```

```
public:
property double MinLimit {
	double get ();
	void set (double value);
}
```

```
member MinLimit : float with get, set
```

#### 属性值

类型：Double

![](../icons/SectionExpanded.png)参见

#### 引用

[PIDHelper 类](eaa06577-fe3f-57e1-0cf7-e906f4587e04.htm)

[HslCommunication.Algorithms.PID 命名空间](e119de70-6089-f087-5e46-755795fa33c0.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## SetValue 属性 

[原文連結](http://api.hslcommunication.cn/html/fbfe17cc-dc49-58eb-b5e4-bfeb2104c873.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Algorithms.PID](../html/e119de70-6089-f087-5e46-755795fa33c0.htm "HslCommunication.Algorithms.PID")

[PIDHelper 类](../html/eaa06577-fe3f-57e1-0cf7-e906f4587e04.htm "PIDHelper 类")

[PIDHelper 属性](../html/5ea37cfa-1260-b972-3a8f-7401785d1974.htm "PIDHelper 属性")

[DeadBand 属性](../html/7b622a95-f3df-5f68-fcb4-1e16d435524c.htm "DeadBand 属性 ")

[Kd 属性](../html/e35edbee-1a85-8fe7-a2e4-c2e277dec003.htm "Kd 属性 ")

[Ki 属性](../html/2aeb7f72-333a-c560-5820-23d55a9d5273.htm "Ki 属性 ")

[Kp 属性](../html/ec4edeae-f643-e04e-7d0c-96cccf1c6dc6.htm "Kp 属性 ")

[MaxLimit 属性](../html/4e794846-ffe3-ca9e-c9a1-a4c34f2672ce.htm "MaxLimit 属性 ")

[MinLimit 属性](../html/5363d5cc-63eb-803f-b18a-e5100b52aa03.htm "MinLimit 属性 ")

[SetValue 属性](../html/fbfe17cc-dc49-58eb-b5e4-bfeb2104c873.htm "SetValue 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| PIDHelperSetValue 属性 |

获取或设置当前设置的值

**命名空间：**
 [HslCommunication.Algorithms.PID](e119de70-6089-f087-5e46-755795fa33c0.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public double SetValue { get; set; }
```

```
Public Property SetValue As Double
	Get
	Set
```

```
public:
property double SetValue {
	double get ();
	void set (double value);
}
```

```
member SetValue : float with get, set
```

#### 属性值

类型：Double

![](../icons/SectionExpanded.png)参见

#### 引用

[PIDHelper 类](eaa06577-fe3f-57e1-0cf7-e906f4587e04.htm)

[HslCommunication.Algorithms.PID 命名空间](e119de70-6089-f087-5e46-755795fa33c0.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## PIDHelper 方法

[原文連結](http://api.hslcommunication.cn/html/4e903b1d-3732-4465-2e98-35718af9acab.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Algorithms.PID](../html/e119de70-6089-f087-5e46-755795fa33c0.htm "HslCommunication.Algorithms.PID")

[PIDHelper 类](../html/eaa06577-fe3f-57e1-0cf7-e906f4587e04.htm "PIDHelper 类")

[PIDHelper 方法](../html/4e903b1d-3732-4465-2e98-35718af9acab.htm "PIDHelper 方法")

[PidCalculate 方法](../html/cf87f2b3-ce97-1b6c-f0b9-bbea07ca7d2e.htm "PidCalculate 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| PIDHelper 方法 |

[PIDHelper](eaa06577-fe3f-57e1-0cf7-e906f4587e04.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法 | [PidCalculate](cf87f2b3-ce97-1b6c-f0b9-bbea07ca7d2e.htm) | 计算Pid数据的值 |
| 公共方法 | ToString | (继承自 Object。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[PIDHelper 类](eaa06577-fe3f-57e1-0cf7-e906f4587e04.htm)

[HslCommunication.Algorithms.PID 命名空间](e119de70-6089-f087-5e46-755795fa33c0.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## PidCalculate 方法 

[原文連結](http://api.hslcommunication.cn/html/cf87f2b3-ce97-1b6c-f0b9-bbea07ca7d2e.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Algorithms.PID](../html/e119de70-6089-f087-5e46-755795fa33c0.htm "HslCommunication.Algorithms.PID")

[PIDHelper 类](../html/eaa06577-fe3f-57e1-0cf7-e906f4587e04.htm "PIDHelper 类")

[PIDHelper 方法](../html/4e903b1d-3732-4465-2e98-35718af9acab.htm "PIDHelper 方法")

[PidCalculate 方法](../html/cf87f2b3-ce97-1b6c-f0b9-bbea07ca7d2e.htm "PidCalculate 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| PIDHelperPidCalculate 方法 |

计算Pid数据的值

**命名空间：**
 [HslCommunication.Algorithms.PID](e119de70-6089-f087-5e46-755795fa33c0.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public double PidCalculate()
```

```
Public Function PidCalculate As Double
```

```
public:
double PidCalculate()
```

```
member PidCalculate : unit -> float 
```

#### 返回值

类型：Double  
计算值

![](../icons/SectionExpanded.png)参见

#### 引用

[PIDHelper 类](eaa06577-fe3f-57e1-0cf7-e906f4587e04.htm)

[HslCommunication.Algorithms.PID 命名空间](e119de70-6089-f087-5e46-755795fa33c0.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## PidMode 枚举

[原文連結](http://api.hslcommunication.cn/html/88867bf9-a36f-5029-1afc-91646e44b13a.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Algorithms.PID](../html/e119de70-6089-f087-5e46-755795fa33c0.htm "HslCommunication.Algorithms.PID")

[PIDHelper 类](../html/eaa06577-fe3f-57e1-0cf7-e906f4587e04.htm "PIDHelper 类")

[PidMode 枚举](../html/88867bf9-a36f-5029-1afc-91646e44b13a.htm "PidMode 枚举")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| PidMode 枚举 |

Pid的模式选择

**命名空间：**
 [HslCommunication.Algorithms.PID](e119de70-6089-f087-5e46-755795fa33c0.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public enum PidMode
```

```
Public Enumeration PidMode
```

```
public enum class PidMode
```

```
type PidMode
```

![](../icons/SectionExpanded.png)成员

|  | 成员名称 | 值 | 说明 |
| --- | --- | --- | --- |
|  | Increment | 1 | 增量模式 |

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.Algorithms.PID 命名空间](e119de70-6089-f087-5e46-755795fa33c0.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)