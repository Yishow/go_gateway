# HslCommunication - HslCommunication.Profinet.Geniitek

> 分類頁數: 30



---
## HslCommunication.Profinet.Geniitek

[原文連結](http://api.hslcommunication.cn/html/307ec5f8-718a-eb03-ee2d-5de6962e72cd.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Geniitek](../html/307ec5f8-718a-eb03-ee2d-5de6962e72cd.htm "HslCommunication.Profinet.Geniitek")

[VibrationSensorActualValue 结构](../html/b464dceb-49ff-b456-418f-0db478690aaa.htm "VibrationSensorActualValue 结构")

[VibrationSensorClient 类](../html/f4100855-734d-f593-67b3-891e63f61804.htm "VibrationSensorClient 类")

[VibrationSensorClient.OnActualValueReceiveDelegate 委托](../html/25739141-a12a-ca22-7443-62162c2c19d9.htm "VibrationSensorClient.OnActualValueReceiveDelegate 委托")

[VibrationSensorClient.OnClientConnectedDelegate 委托](../html/ff07e8ff-aac9-0fc1-08a1-d6878b6f40ea.htm "VibrationSensorClient.OnClientConnectedDelegate 委托")

[VibrationSensorClient.OnPeekValueReceiveDelegate 委托](../html/91ea5b2f-d211-3f51-f32c-d73a03c83ca6.htm "VibrationSensorClient.OnPeekValueReceiveDelegate 委托")

[VibrationSensorLongMessage 类](../html/900e9fa4-2bdf-fb0f-10b3-b1da76626266.htm "VibrationSensorLongMessage 类")

[VibrationSensorPeekValue 类](../html/00a9c380-142d-e06d-054a-b7e276725a22.htm "VibrationSensorPeekValue 类")

[VibrationSensorShortMessage 类](../html/9f1c87b5-9a6e-401e-1243-f37a10ad1fad.htm "VibrationSensorShortMessage 类")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslCommunication.Profinet.Geniitek 命名空间 |

[缺少 "N:HslCommunication.Profinet.Geniitek" 的 <summary> 文档]

![](../icons/SectionExpanded.png)类

|  | 类 | 说明 |
| --- | --- | --- |
| 公共类 | [VibrationSensorClient](f4100855-734d-f593-67b3-891e63f61804.htm) | Geniitek-VB31 型号的智能无线振动传感器，来自苏州捷杰传感器技术有限公司 |
| 公共类 | [VibrationSensorLongMessage](900e9fa4-2bdf-fb0f-10b3-b1da76626266.htm) | 完整的数据报文信息 |
| 公共类 | [VibrationSensorPeekValue](00a9c380-142d-e06d-054a-b7e276725a22.htm) | 振动传感器的峰值数据类 |
| 公共类 | [VibrationSensorShortMessage](9f1c87b5-9a6e-401e-1243-f37a10ad1fad.htm) | 短消息的报文内容 |

![](../icons/SectionExpanded.png)结构

|  | 结构 | 说明 |
| --- | --- | --- |
| 公共结构 | [VibrationSensorActualValue](b464dceb-49ff-b456-418f-0db478690aaa.htm) | 振动传感器的加速度值 |

![](../icons/SectionExpanded.png)委托

|  | 委托 | 说明 |
| --- | --- | --- |
| 公共委托 | [VibrationSensorClientOnActualValueReceiveDelegate](25739141-a12a-ca22-7443-62162c2c19d9.htm) | 震动传感器实时数据事件委托  Vibration sensor real-time data event delegation |
| 公共委托 | [VibrationSensorClientOnClientConnectedDelegate](ff07e8ff-aac9-0fc1-08a1-d6878b6f40ea.htm) | 连接服务器成功的委托  Connection server successfully delegated |
| 公共委托 | [VibrationSensorClientOnPeekValueReceiveDelegate](91ea5b2f-d211-3f51-f32c-d73a03c83ca6.htm) | 震动传感器峰值数据事件委托  Shock sensor peak data event delegation |

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## VibrationSensorActualValue 结构

[原文連結](http://api.hslcommunication.cn/html/b464dceb-49ff-b456-418f-0db478690aaa.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Geniitek](../html/307ec5f8-718a-eb03-ee2d-5de6962e72cd.htm "HslCommunication.Profinet.Geniitek")

[VibrationSensorActualValue 结构](../html/b464dceb-49ff-b456-418f-0db478690aaa.htm "VibrationSensorActualValue 结构")

[VibrationSensorActualValue 属性](../html/f269b234-e75e-c60e-2afb-838d117b0541.htm "VibrationSensorActualValue 属性")

[VibrationSensorActualValue 方法](../html/f7fbe8da-2ab8-d302-3de6-82f3b0dbfb80.htm "VibrationSensorActualValue 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| VibrationSensorActualValue 结构 |

振动传感器的加速度值

**命名空间：**
 [HslCommunication.Profinet.Geniitek](307ec5f8-718a-eb03-ee2d-5de6962e72cd.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public struct VibrationSensorActualValue
```

```
Public Structure VibrationSensorActualValue
```

```
public value class VibrationSensorActualValue
```

```
[<SealedAttribute>]
type VibrationSensorActualValue =  struct end
```

VibrationSensorActualValue 类型公开以下成员。

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [AcceleratedSpeedX](7778c23d-0ded-4547-47c9-7e00bbccfa74.htm) | X轴的实时加速度 |
| 公共属性 | [AcceleratedSpeedY](ccf10509-7726-02c1-beda-4ede919bea7e.htm) | Y轴的实时加速度 |
| 公共属性 | [AcceleratedSpeedZ](515094af-0df8-e6c3-8e49-b06076107e7c.htm) | Z轴的实时加速度 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | Equals | (继承自 ValueType。) |
| 公共方法 | GetHashCode | (继承自 ValueType。) |
| 公共方法 | GetType | (继承自 Object。) |
| 公共方法 | [ToString](a463c26d-c789-d364-6c31-6657c8378994.htm) | (重写 ValueTypeToString.) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.Profinet.Geniitek 命名空间](307ec5f8-718a-eb03-ee2d-5de6962e72cd.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## VibrationSensorActualValue 属性

[原文連結](http://api.hslcommunication.cn/html/f269b234-e75e-c60e-2afb-838d117b0541.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Geniitek](../html/307ec5f8-718a-eb03-ee2d-5de6962e72cd.htm "HslCommunication.Profinet.Geniitek")

[VibrationSensorActualValue 结构](../html/b464dceb-49ff-b456-418f-0db478690aaa.htm "VibrationSensorActualValue 结构")

[VibrationSensorActualValue 属性](../html/f269b234-e75e-c60e-2afb-838d117b0541.htm "VibrationSensorActualValue 属性")

[AcceleratedSpeedX 属性](../html/7778c23d-0ded-4547-47c9-7e00bbccfa74.htm "AcceleratedSpeedX 属性 ")

[AcceleratedSpeedY 属性](../html/ccf10509-7726-02c1-beda-4ede919bea7e.htm "AcceleratedSpeedY 属性 ")

[AcceleratedSpeedZ 属性](../html/515094af-0df8-e6c3-8e49-b06076107e7c.htm "AcceleratedSpeedZ 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| VibrationSensorActualValue 属性 |

[VibrationSensorActualValue](b464dceb-49ff-b456-418f-0db478690aaa.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [AcceleratedSpeedX](7778c23d-0ded-4547-47c9-7e00bbccfa74.htm) | X轴的实时加速度 |
| 公共属性 | [AcceleratedSpeedY](ccf10509-7726-02c1-beda-4ede919bea7e.htm) | Y轴的实时加速度 |
| 公共属性 | [AcceleratedSpeedZ](515094af-0df8-e6c3-8e49-b06076107e7c.htm) | Z轴的实时加速度 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[VibrationSensorActualValue 结构](b464dceb-49ff-b456-418f-0db478690aaa.htm)

[HslCommunication.Profinet.Geniitek 命名空间](307ec5f8-718a-eb03-ee2d-5de6962e72cd.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## AcceleratedSpeedX 属性 

[原文連結](http://api.hslcommunication.cn/html/7778c23d-0ded-4547-47c9-7e00bbccfa74.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Geniitek](../html/307ec5f8-718a-eb03-ee2d-5de6962e72cd.htm "HslCommunication.Profinet.Geniitek")

[VibrationSensorActualValue 结构](../html/b464dceb-49ff-b456-418f-0db478690aaa.htm "VibrationSensorActualValue 结构")

[VibrationSensorActualValue 属性](../html/f269b234-e75e-c60e-2afb-838d117b0541.htm "VibrationSensorActualValue 属性")

[AcceleratedSpeedX 属性](../html/7778c23d-0ded-4547-47c9-7e00bbccfa74.htm "AcceleratedSpeedX 属性 ")

[AcceleratedSpeedY 属性](../html/ccf10509-7726-02c1-beda-4ede919bea7e.htm "AcceleratedSpeedY 属性 ")

[AcceleratedSpeedZ 属性](../html/515094af-0df8-e6c3-8e49-b06076107e7c.htm "AcceleratedSpeedZ 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| VibrationSensorActualValueAcceleratedSpeedX 属性 |

X轴的实时加速度

**命名空间：**
 [HslCommunication.Profinet.Geniitek](307ec5f8-718a-eb03-ee2d-5de6962e72cd.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public float AcceleratedSpeedX { get; set; }
```

```
Public Property AcceleratedSpeedX As Single
	Get
	Set
```

```
public:
property float AcceleratedSpeedX {
	float get ();
	void set (float value);
}
```

```
member AcceleratedSpeedX : float32 with get, set
```

#### 属性值

类型：Single

![](../icons/SectionExpanded.png)参见

#### 引用

[VibrationSensorActualValue 结构](b464dceb-49ff-b456-418f-0db478690aaa.htm)

[HslCommunication.Profinet.Geniitek 命名空间](307ec5f8-718a-eb03-ee2d-5de6962e72cd.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## AcceleratedSpeedY 属性 

[原文連結](http://api.hslcommunication.cn/html/ccf10509-7726-02c1-beda-4ede919bea7e.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Geniitek](../html/307ec5f8-718a-eb03-ee2d-5de6962e72cd.htm "HslCommunication.Profinet.Geniitek")

[VibrationSensorActualValue 结构](../html/b464dceb-49ff-b456-418f-0db478690aaa.htm "VibrationSensorActualValue 结构")

[VibrationSensorActualValue 属性](../html/f269b234-e75e-c60e-2afb-838d117b0541.htm "VibrationSensorActualValue 属性")

[AcceleratedSpeedX 属性](../html/7778c23d-0ded-4547-47c9-7e00bbccfa74.htm "AcceleratedSpeedX 属性 ")

[AcceleratedSpeedY 属性](../html/ccf10509-7726-02c1-beda-4ede919bea7e.htm "AcceleratedSpeedY 属性 ")

[AcceleratedSpeedZ 属性](../html/515094af-0df8-e6c3-8e49-b06076107e7c.htm "AcceleratedSpeedZ 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| VibrationSensorActualValueAcceleratedSpeedY 属性 |

Y轴的实时加速度

**命名空间：**
 [HslCommunication.Profinet.Geniitek](307ec5f8-718a-eb03-ee2d-5de6962e72cd.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public float AcceleratedSpeedY { get; set; }
```

```
Public Property AcceleratedSpeedY As Single
	Get
	Set
```

```
public:
property float AcceleratedSpeedY {
	float get ();
	void set (float value);
}
```

```
member AcceleratedSpeedY : float32 with get, set
```

#### 属性值

类型：Single

![](../icons/SectionExpanded.png)参见

#### 引用

[VibrationSensorActualValue 结构](b464dceb-49ff-b456-418f-0db478690aaa.htm)

[HslCommunication.Profinet.Geniitek 命名空间](307ec5f8-718a-eb03-ee2d-5de6962e72cd.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## AcceleratedSpeedZ 属性 

[原文連結](http://api.hslcommunication.cn/html/515094af-0df8-e6c3-8e49-b06076107e7c.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Geniitek](../html/307ec5f8-718a-eb03-ee2d-5de6962e72cd.htm "HslCommunication.Profinet.Geniitek")

[VibrationSensorActualValue 结构](../html/b464dceb-49ff-b456-418f-0db478690aaa.htm "VibrationSensorActualValue 结构")

[VibrationSensorActualValue 属性](../html/f269b234-e75e-c60e-2afb-838d117b0541.htm "VibrationSensorActualValue 属性")

[AcceleratedSpeedX 属性](../html/7778c23d-0ded-4547-47c9-7e00bbccfa74.htm "AcceleratedSpeedX 属性 ")

[AcceleratedSpeedY 属性](../html/ccf10509-7726-02c1-beda-4ede919bea7e.htm "AcceleratedSpeedY 属性 ")

[AcceleratedSpeedZ 属性](../html/515094af-0df8-e6c3-8e49-b06076107e7c.htm "AcceleratedSpeedZ 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| VibrationSensorActualValueAcceleratedSpeedZ 属性 |

Z轴的实时加速度

**命名空间：**
 [HslCommunication.Profinet.Geniitek](307ec5f8-718a-eb03-ee2d-5de6962e72cd.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public float AcceleratedSpeedZ { get; set; }
```

```
Public Property AcceleratedSpeedZ As Single
	Get
	Set
```

```
public:
property float AcceleratedSpeedZ {
	float get ();
	void set (float value);
}
```

```
member AcceleratedSpeedZ : float32 with get, set
```

#### 属性值

类型：Single

![](../icons/SectionExpanded.png)参见

#### 引用

[VibrationSensorActualValue 结构](b464dceb-49ff-b456-418f-0db478690aaa.htm)

[HslCommunication.Profinet.Geniitek 命名空间](307ec5f8-718a-eb03-ee2d-5de6962e72cd.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## VibrationSensorActualValue 方法

[原文連結](http://api.hslcommunication.cn/html/f7fbe8da-2ab8-d302-3de6-82f3b0dbfb80.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Geniitek](../html/307ec5f8-718a-eb03-ee2d-5de6962e72cd.htm "HslCommunication.Profinet.Geniitek")

[VibrationSensorActualValue 结构](../html/b464dceb-49ff-b456-418f-0db478690aaa.htm "VibrationSensorActualValue 结构")

[VibrationSensorActualValue 方法](../html/f7fbe8da-2ab8-d302-3de6-82f3b0dbfb80.htm "VibrationSensorActualValue 方法")

[ToString 方法](../html/a463c26d-c789-d364-6c31-6657c8378994.htm "ToString 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| VibrationSensorActualValue 方法 |

[VibrationSensorActualValue](b464dceb-49ff-b456-418f-0db478690aaa.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | Equals | (继承自 ValueType。) |
| 公共方法 | GetHashCode | (继承自 ValueType。) |
| 公共方法 | GetType | (继承自 Object。) |
| 公共方法 | [ToString](a463c26d-c789-d364-6c31-6657c8378994.htm) | (重写 ValueTypeToString.) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[VibrationSensorActualValue 结构](b464dceb-49ff-b456-418f-0db478690aaa.htm)

[HslCommunication.Profinet.Geniitek 命名空间](307ec5f8-718a-eb03-ee2d-5de6962e72cd.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ToString 方法 

[原文連結](http://api.hslcommunication.cn/html/a463c26d-c789-d364-6c31-6657c8378994.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Geniitek](../html/307ec5f8-718a-eb03-ee2d-5de6962e72cd.htm "HslCommunication.Profinet.Geniitek")

[VibrationSensorActualValue 结构](../html/b464dceb-49ff-b456-418f-0db478690aaa.htm "VibrationSensorActualValue 结构")

[VibrationSensorActualValue 方法](../html/f7fbe8da-2ab8-d302-3de6-82f3b0dbfb80.htm "VibrationSensorActualValue 方法")

[ToString 方法](../html/a463c26d-c789-d364-6c31-6657c8378994.htm "ToString 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| VibrationSensorActualValueToString 方法 |

[缺少 "M:HslCommunication.Profinet.Geniitek.VibrationSensorActualValue.ToString" 的 <summary> 文档]

**命名空间：**
 [HslCommunication.Profinet.Geniitek](307ec5f8-718a-eb03-ee2d-5de6962e72cd.htm)  
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

[缺少 "M:HslCommunication.Profinet.Geniitek.VibrationSensorActualValue.ToString" 的 <returns> 文档]

![](../icons/SectionExpanded.png)参见

#### 引用

[VibrationSensorActualValue 结构](b464dceb-49ff-b456-418f-0db478690aaa.htm)

[HslCommunication.Profinet.Geniitek 命名空间](307ec5f8-718a-eb03-ee2d-5de6962e72cd.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## VibrationSensorClient 类

[原文連結](http://api.hslcommunication.cn/html/f4100855-734d-f593-67b3-891e63f61804.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Geniitek](../html/307ec5f8-718a-eb03-ee2d-5de6962e72cd.htm "HslCommunication.Profinet.Geniitek")

[VibrationSensorClient 类](../html/f4100855-734d-f593-67b3-891e63f61804.htm "VibrationSensorClient 类")

[VibrationSensorClient 构造函数](../html/d8565314-1894-17a8-a16a-9776a924cb8a.htm "VibrationSensorClient 构造函数 ")

[VibrationSensorClient 属性](../html/9da4a7a8-fc48-2eb3-3478-45cdda9f0e50.htm "VibrationSensorClient 属性")

[VibrationSensorClient 方法](../html/a7ca1954-ca0b-784f-4a99-2bb3ef651141.htm "VibrationSensorClient 方法")

[VibrationSensorClient 事件](../html/edc6c96c-eab1-7e9e-e1c9-a79c676754a5.htm "VibrationSensorClient 事件")

[VibrationSensorClient 字段](../html/50704ab2-06f9-69f8-6024-2715d539a965.htm "VibrationSensorClient 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| VibrationSensorClient 类 |

Geniitek-VB31 型号的智能无线振动传感器，来自苏州捷杰传感器技术有限公司

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  [HslCommunication.Core.NetNetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)  
    [HslCommunication.Core.NetNetworkXBase](ab528a6a-e713-5202-59e0-90dd6a66de5f.htm)  
      HslCommunication.Profinet.GeniitekVibrationSensorClient

**命名空间：**
 [HslCommunication.Profinet.Geniitek](307ec5f8-718a-eb03-ee2d-5de6962e72cd.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class VibrationSensorClient : NetworkXBase
```

```
Public Class VibrationSensorClient
	Inherits NetworkXBase
```

```
public ref class VibrationSensorClient : public NetworkXBase
```

```
type VibrationSensorClient =  
    class
        inherit NetworkXBase
    end
```

VibrationSensorClient 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [VibrationSensorClient](d8565314-1894-17a8-a16a-9776a924cb8a.htm) | 使用指定的ip，端口来实例化一个默认的对象 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [Address](08f34568-0a4b-e548-17b8-73158f4a705f.htm) | 当前设备的地址信息 |
| 公共属性 | [CheckSeconds](518e4c3b-021d-6a00-2ab4-8808fa91f5f9.htm) | 获取或设置当前的客户端假死超时检查时间，单位为秒，默认60秒，60秒内没有接收到传感器的数据，则强制重连。 |
| 公共属性 | [ConnectTimeOut](07b3c8d2-d1e3-da43-4742-235b02299de5.htm) | 获取或设置当前客户端的连接超时时间，默认10,000毫秒，单位ms  Gets or sets the connection timeout of the current client. The default is 10,000 milliseconds. The unit is ms. |
| 公共属性代码示例 | [LogNet](46dfe698-46c6-73a2-5458-8fd568d264ce.htm) | 组件的日志工具，支持日志记录，只要实例化后，当前网络的基本信息，就以[DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm)等级进行输出  The component's logging tool supports logging. As long as the instantiation of the basic network information, the output will be output at [DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm) (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 公共属性代码示例 | [Token](80db00a7-aaa5-164f-e87c-af5400dc80fa.htm) | 网络类的身份令牌，在hsl协议的模式下会有效，在和设备进行通信的时候是无效的  Network-type identity tokens will be valid in the hsl protocol mode and will not be valid when communicating with the device (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法静态成员 | [BulidLongMessage](96de5e3f-4a4a-df2e-262e-ad19dc2c5b2f.htm) | 根据地址，命令，数据，创建向传感器发送的数据信息 |
| 受保护的方法 | [CheckRemoteToken](1d231984-101b-0fff-0daf-1d18581a4241.htm) | 检查当前的头子节信息的令牌是否是正确的，仅用于某些特殊的协议实现  Check whether the token of the current header subsection information is correct, only for some special protocol implementations (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 公共方法静态成员 | [CheckXor](a6ea980c-719c-22a0-7d04-516f7b4dbb1b.htm) | 检查当前的数据是否XOR校验成功 |
| 公共方法 | [ConnectClose](beb87b8e-1ec5-c550-7471-642cc9c798a2.htm) | 关闭Mqtt服务器的连接。  Close the connection to the Mqtt server. |
| 公共方法 | [ConnectServer](ab05e736-dd37-a9b3-5a4b-0f267052ef77.htm) | 连接服务器，实例化客户端之后，至少要调用成功一次，如果返回失败，那些请过一段时间后重新调用本方法连接。  After connecting to the server, the client must be called at least once after instantiating the client. If the return fails, please call this method to connect again after a period of time. |
| 公共方法 | [ConnectServerAsync](7a6ef23e-9bde-da7f-d8fa-2d8e2a87e6a2.htm) | 连接服务器，实例化客户端之后，至少要调用成功一次，如果返回失败，那些请过一段时间后重新调用本方法连接。  After connecting to the server, the client must be called at least once after instantiating the client. If the return fails, please call this method to connect again after a period of time. |
| 受保护的方法代码示例 | [CreateSocketAndConnect(String, Int32)](a153ae27-92e3-4fe7-b2e7-c8d207bebbb0.htm) | 创建一个新的socket对象并连接到远程的地址，默认超时时间为10秒钟，需要指定ip地址以及端口号信息  Create a new socket object and connect to the remote address. The default timeout is 10 seconds. You need to specify the IP address and port number. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法代码示例 | [CreateSocketAndConnect(IPEndPoint, Int32, IPEndPoint)](7ce6288f-c6db-e968-d676-f851d2b05c77.htm) | 创建一个新的socket对象并连接到远程的地址，需要指定远程终结点，超时时间（单位是毫秒），如果需要绑定本地的IP或是端口，传入 local对象  To create a new socket object and connect to the remote address, you need to specify the remote endpoint, the timeout period (in milliseconds), if you need to bind the local IP or port, pass in the local object (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法代码示例 | [CreateSocketAndConnect(String, Int32, Int32)](d0e14a5f-1e97-e59d-078f-4fbf0412b865.htm) | 创建一个新的socket对象并连接到远程的地址，需要指定ip地址以及端口号信息，还有超时时间，单位是毫秒  To create a new socket object and connect to a remote address, you need to specify the IP address and port number information, and the timeout period in milliseconds (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法代码示例 | [CreateSocketAndConnectAsync(String, Int32)](0227bc60-1c7f-3c1d-0cd6-e49bb904d6df.htm) | 创建一个新的socket对象并连接到远程的地址，默认超时时间为10秒钟，需要指定ip地址以及端口号信息  Create a new socket object and connect to the remote address. The default timeout is 10 seconds. You need to specify the IP address and port number. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法代码示例 | [CreateSocketAndConnectAsync(IPEndPoint, Int32, IPEndPoint)](ccec07bf-cfa4-b302-210e-441821e095a4.htm) | 创建一个新的socket对象并连接到远程的地址，需要指定远程终结点，超时时间（单位是毫秒），如果需要绑定本地的IP或是端口，传入 local对象  To create a new socket object and connect to the remote address, you need to specify the remote endpoint, the timeout period (in milliseconds), if you need to bind the local IP or port, pass in the local object (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法代码示例 | [CreateSocketAndConnectAsync(String, Int32, Int32)](7ec25883-60e2-c7ac-2f90-0ea1ae0c4945.htm) | 创建一个新的socket对象并连接到远程的地址，需要指定ip地址以及端口号信息，还有超时时间，单位是毫秒  To create a new socket object and connect to a remote address, you need to specify the IP address and port number information, and the timeout period in milliseconds (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [DeleteFileByName](59afa4a4-3f3c-9831-26ee-c7c13a213bf7.htm) | 删除一个指定的文件，如果文件不存在，直接返回 True，如果文件存在则直接删除，删除成功返回 True，如果发生了异常，返回False  Delete a specified file, if the file does not exist, return True directly, if the file exists, delete it directly, if the deletion is successful, return True, if an exception occurs, return  False (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
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
| 受保护的方法 | [ReceiveFileFromSocket(Socket, Stream, ActionInt64, Int64)](a227f0b6-3050-e5d1-79c8-88ea9719bd06.htm) | [自校验] 从网络中接收一个文件，写入数据流，如果结果异常，则结束通讯，参数顺序文件名，文件大小，文件标识，上传人  [Self-checking] Receive a file from the network. If the result is abnormal, the communication ends. (继承自 [NetworkXBase](ab528a6a-e713-5202-59e0-90dd6a66de5f.htm)。) |
| 受保护的方法 | [ReceiveFileFromSocket(Socket, String, ActionInt64, Int64)](6bcd5a10-d79c-649e-fb65-cd78fecf99b7.htm) | [自校验] 从网络中接收一个文件，如果结果异常，则结束通讯  [Self-checking] Receive a file from the network. If the result is abnormal, the communication ends. (继承自 [NetworkXBase](ab528a6a-e713-5202-59e0-90dd6a66de5f.htm)。) |
| 受保护的方法 | [ReceiveFileFromSocketAsync(Socket, Stream, ActionInt64, Int64)](e3b415d2-58a9-7ab6-0097-0bdc736ca0dc.htm) | [自校验] 从网络中接收一个文件，写入数据流，如果结果异常，则结束通讯，参数顺序文件名，文件大小，文件标识，上传人  [Self-checking] Receive a file from the network. If the result is abnormal, the communication ends. (继承自 [NetworkXBase](ab528a6a-e713-5202-59e0-90dd6a66de5f.htm)。) |
| 受保护的方法 | [ReceiveFileFromSocketAsync(Socket, String, ActionInt64, Int64)](ab9b6255-580b-464b-447a-5de92b6f3dd9.htm) | [自校验] 从网络中接收一个文件，如果结果异常，则结束通讯  [Self-checking] Receive a file from the network. If the result is abnormal, the communication ends. (继承自 [NetworkXBase](ab528a6a-e713-5202-59e0-90dd6a66de5f.htm)。) |
| 受保护的方法 | [ReceiveFileHeadFromSocket](37b7d2d2-faeb-2263-627d-6a74bcbce47d.htm) | [自校验] 从套接字中接收文件头信息  [Self-checking] Receive file header information from socket (继承自 [NetworkXBase](ab528a6a-e713-5202-59e0-90dd6a66de5f.htm)。) |
| 受保护的方法 | [ReceiveFileHeadFromSocketAsync](9f10bf85-76d0-c721-60f1-1b42425a52af.htm) | [自校验] 从套接字中接收文件头信息  [Self-checking] Receive file header information from socket (继承自 [NetworkXBase](ab528a6a-e713-5202-59e0-90dd6a66de5f.htm)。) |
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
| 受保护的方法 | [SendFileAndCheckReceive(Socket, Stream, String, String, String, ActionInt64, Int64)](602977f6-83fa-486f-7b4e-7913594bb847.htm) | [自校验] 将流数据发送至套接字，具体发送细节将在继承类中实现，如果结果异常，则结束通讯  [Self-checking] Send stream data to the socket. The specific sending details will be implemented in the inherited class. If the result is abnormal, the communication will be terminated (继承自 [NetworkXBase](ab528a6a-e713-5202-59e0-90dd6a66de5f.htm)。) |
| 受保护的方法 | [SendFileAndCheckReceive(Socket, String, String, String, String, ActionInt64, Int64)](cad4c115-2f7b-5301-36e4-66b6da8d864f.htm) | [自校验] 将文件数据发送至套接字，具体发送细节将在继承类中实现，如果结果异常，则结束通讯  [Self-checking] Send the file data to the socket. The specific sending details will be implemented in the inherited class. If the result is abnormal, the communication will end (继承自 [NetworkXBase](ab528a6a-e713-5202-59e0-90dd6a66de5f.htm)。) |
| 受保护的方法 | [SendFileAndCheckReceiveAsync(Socket, Stream, String, String, String, ActionInt64, Int64)](8ad77d7e-8710-5bb5-e39a-9e613957cf72.htm) | [自校验] 将流数据发送至套接字，具体发送细节将在继承类中实现，如果结果异常，则结束通讯  [Self-checking] Send stream data to the socket. The specific sending details will be implemented in the inherited class. If the result is abnormal, the communication will be terminated (继承自 [NetworkXBase](ab528a6a-e713-5202-59e0-90dd6a66de5f.htm)。) |
| 受保护的方法 | [SendFileAndCheckReceiveAsync(Socket, String, String, String, String, ActionInt64, Int64)](9eff5d17-5f46-6900-b36b-f06e76f3605e.htm) | [自校验] 将文件数据发送至套接字，具体发送细节将在继承类中实现，如果结果异常，则结束通讯  [Self-checking] Send the file data to the socket. The specific sending details will be implemented in the inherited class. If the result is abnormal, the communication will end (继承自 [NetworkXBase](ab528a6a-e713-5202-59e0-90dd6a66de5f.htm)。) |
| 受保护的方法 | [SendFileStreamToSocket](73df4306-17c8-3649-8214-b0704cd009af.htm) | [自校验] 将文件数据发送至套接字，如果结果异常，则结束通讯  [Self-check] Send the file data to the socket. If the result is abnormal, the communication is ended. (继承自 [NetworkXBase](ab528a6a-e713-5202-59e0-90dd6a66de5f.htm)。) |
| 受保护的方法 | [SendFileStreamToSocketAsync](41b6ebdb-22c3-fe61-113b-a60f1674750c.htm) | [自校验] 将文件数据发送至套接字，如果结果异常，则结束通讯  [Self-check] Send the file data to the socket. If the result is abnormal, the communication is ended. (继承自 [NetworkXBase](ab528a6a-e713-5202-59e0-90dd6a66de5f.htm)。) |
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
| 公共方法 | [SetReadActual](6bbeb3f4-690a-b65f-dbce-bdec86750011.htm) | 设置读取震动传感器的实时加速度  Set the real-time acceleration of the vibration sensor |
| 公共方法 | [SetReadStatus](3e411f8c-ebb8-189f-432d-417dbda9b682.htm) | 设置读取震动传感器的状态数据  Set to read the status data of the shock sensor |
| 公共方法 | [SetReadStatusInterval](958c6cda-4e99-28d9-a088-2e30f07c134a.htm) | 设置当前的震动传感器的数据发送间隔为指定的时间，单位为秒  Set the current vibration sensor data transmission interval to the specified time in seconds |
| 公共方法 | [ToString](77cd262f-cbad-bdd6-ab4e-9ffedd3940a5.htm) | (重写 [NetworkXBaseToString](08bf006b-af6e-7c58-6935-8e8c387e96ab.htm).) |
| 受保护的方法 | [WriteStreamFromSocket](97bf95ea-51fe-8f2e-5da6-1c83766aaa78.htm) | 从套接字中接收所有的数据然后写入到指定的流当中去，需要指定数据的长度，支持按照百分比进行进度报告  Receives all data from the socket and writes it to the specified stream. The length of the data needs to be specified, and progress reporting is supported in percentage. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [WriteStreamFromSocketAsync](b4cf7f6d-1e0c-a266-4c49-683fdb14a922.htm) | 从套接字中接收所有的数据然后写入到指定的流当中去，需要指定数据的长度，支持按照百分比进行进度报告  Receives all data from the socket and writes it to the specified stream. The length of the data needs to be specified, and progress reporting is supported in percentage. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)事件

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共事件 | [OnActualValueReceive](b88d88ef-cb1e-80db-732c-2cd1f794d45b.htm) | 接收到震动传感器实时数据时触发  Triggered when real-time data from shock sensor is received |
| 公共事件 | [OnClientConnected](e371eef0-ea60-a452-ed75-437958da9036.htm) | 当客户端连接成功触发事件，就算是重新连接服务器后，也是会触发的  The event is triggered when the client is connected successfully, even after reconnecting to the server. |
| 公共事件 | [OnNetworkError](51179b61-d3a9-44dc-cf21-e33c6beb87fe.htm) | 当网络发生异常的时候触发的事件，用户应该在事件里进行重连服务器 |
| 公共事件 | [OnPeekValueReceive](9ce2d21b-57a6-103e-3800-95d0af9f300b.htm) | 接收到震动传感器峰值数据时触发  Triggered when peak data of vibration sensor is received |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)字段

|  | 名称 | 说明 |
| --- | --- | --- |
| 受保护的字段 | [CoreSocket](6831da9f-21e0-8967-f17c-641b330e0d80.htm) | 对客户端而言是的通讯用的套接字，对服务器来说是用于侦听的套接字  A communication socket for the client, or a listening socket for the server (继承自 [NetworkXBase](ab528a6a-e713-5202-59e0-90dd6a66de5f.htm)。) |
| 受保护的字段 | [fileCacheSize](2328fb89-cfda-7fa9-88d0-98f1ae3c7c54.htm) | 文件传输的时候的缓存大小，直接影响传输的速度，值越大，传输速度越快，越占内存，默认为100K大小  The size of the cache during file transfer directly affects the speed of the transfer. The larger the value, the faster the transfer speed and the more memory it takes. The default size is 100K. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.Profinet.Geniitek 命名空间](307ec5f8-718a-eb03-ee2d-5de6962e72cd.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## VibrationSensorClient 构造函数 

[原文連結](http://api.hslcommunication.cn/html/d8565314-1894-17a8-a16a-9776a924cb8a.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Geniitek](../html/307ec5f8-718a-eb03-ee2d-5de6962e72cd.htm "HslCommunication.Profinet.Geniitek")

[VibrationSensorClient 类](../html/f4100855-734d-f593-67b3-891e63f61804.htm "VibrationSensorClient 类")

[VibrationSensorClient 构造函数](../html/d8565314-1894-17a8-a16a-9776a924cb8a.htm "VibrationSensorClient 构造函数 ")

[VibrationSensorClient 属性](../html/9da4a7a8-fc48-2eb3-3478-45cdda9f0e50.htm "VibrationSensorClient 属性")

[VibrationSensorClient 方法](../html/a7ca1954-ca0b-784f-4a99-2bb3ef651141.htm "VibrationSensorClient 方法")

[VibrationSensorClient 事件](../html/edc6c96c-eab1-7e9e-e1c9-a79c676754a5.htm "VibrationSensorClient 事件")

[VibrationSensorClient 字段](../html/50704ab2-06f9-69f8-6024-2715d539a965.htm "VibrationSensorClient 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| VibrationSensorClient 构造函数 |

使用指定的ip，端口来实例化一个默认的对象

**命名空间：**
 [HslCommunication.Profinet.Geniitek](307ec5f8-718a-eb03-ee2d-5de6962e72cd.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public VibrationSensorClient(
	string ipAddress = "192.168.1.1",
	int port = 3001
)
```

```
Public Sub New ( 
	Optional ipAddress As String = "192.168.1.1",
	Optional port As Integer = 3001
)
```

```
public:
VibrationSensorClient(
	String^ ipAddress = L"192.168.1.1", 
	int port = 3001
)
```

```
new : 
        ?ipAddress : string * 
        ?port : int 
(* Defaults:
        let _ipAddress = defaultArg ipAddress "192.168.1.1"
        let _port = defaultArg port 3001
*)
-> VibrationSensorClient
```

#### 参数

ipAddress (Optional)
:   类型：SystemString  
    Ip地址信息

port (Optional)
:   类型：SystemInt32  
    端口号信息

![](../icons/SectionExpanded.png)参见

#### 引用

[VibrationSensorClient 类](f4100855-734d-f593-67b3-891e63f61804.htm)

[HslCommunication.Profinet.Geniitek 命名空间](307ec5f8-718a-eb03-ee2d-5de6962e72cd.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## VibrationSensorClient 属性

[原文連結](http://api.hslcommunication.cn/html/9da4a7a8-fc48-2eb3-3478-45cdda9f0e50.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Geniitek](../html/307ec5f8-718a-eb03-ee2d-5de6962e72cd.htm "HslCommunication.Profinet.Geniitek")

[VibrationSensorClient 类](../html/f4100855-734d-f593-67b3-891e63f61804.htm "VibrationSensorClient 类")

[VibrationSensorClient 属性](../html/9da4a7a8-fc48-2eb3-3478-45cdda9f0e50.htm "VibrationSensorClient 属性")

[Address 属性](../html/08f34568-0a4b-e548-17b8-73158f4a705f.htm "Address 属性 ")

[CheckSeconds 属性](../html/518e4c3b-021d-6a00-2ab4-8808fa91f5f9.htm "CheckSeconds 属性 ")

[ConnectTimeOut 属性](../html/07b3c8d2-d1e3-da43-4742-235b02299de5.htm "ConnectTimeOut 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| VibrationSensorClient 属性 |

[VibrationSensorClient](f4100855-734d-f593-67b3-891e63f61804.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [Address](08f34568-0a4b-e548-17b8-73158f4a705f.htm) | 当前设备的地址信息 |
| 公共属性 | [CheckSeconds](518e4c3b-021d-6a00-2ab4-8808fa91f5f9.htm) | 获取或设置当前的客户端假死超时检查时间，单位为秒，默认60秒，60秒内没有接收到传感器的数据，则强制重连。 |
| 公共属性 | [ConnectTimeOut](07b3c8d2-d1e3-da43-4742-235b02299de5.htm) | 获取或设置当前客户端的连接超时时间，默认10,000毫秒，单位ms  Gets or sets the connection timeout of the current client. The default is 10,000 milliseconds. The unit is ms. |
| 公共属性代码示例 | [LogNet](46dfe698-46c6-73a2-5458-8fd568d264ce.htm) | 组件的日志工具，支持日志记录，只要实例化后，当前网络的基本信息，就以[DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm)等级进行输出  The component's logging tool supports logging. As long as the instantiation of the basic network information, the output will be output at [DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm) (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 公共属性代码示例 | [Token](80db00a7-aaa5-164f-e87c-af5400dc80fa.htm) | 网络类的身份令牌，在hsl协议的模式下会有效，在和设备进行通信的时候是无效的  Network-type identity tokens will be valid in the hsl protocol mode and will not be valid when communicating with the device (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[VibrationSensorClient 类](f4100855-734d-f593-67b3-891e63f61804.htm)

[HslCommunication.Profinet.Geniitek 命名空间](307ec5f8-718a-eb03-ee2d-5de6962e72cd.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Address 属性 

[原文連結](http://api.hslcommunication.cn/html/08f34568-0a4b-e548-17b8-73158f4a705f.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Geniitek](../html/307ec5f8-718a-eb03-ee2d-5de6962e72cd.htm "HslCommunication.Profinet.Geniitek")

[VibrationSensorClient 类](../html/f4100855-734d-f593-67b3-891e63f61804.htm "VibrationSensorClient 类")

[VibrationSensorClient 属性](../html/9da4a7a8-fc48-2eb3-3478-45cdda9f0e50.htm "VibrationSensorClient 属性")

[Address 属性](../html/08f34568-0a4b-e548-17b8-73158f4a705f.htm "Address 属性 ")

[CheckSeconds 属性](../html/518e4c3b-021d-6a00-2ab4-8808fa91f5f9.htm "CheckSeconds 属性 ")

[ConnectTimeOut 属性](../html/07b3c8d2-d1e3-da43-4742-235b02299de5.htm "ConnectTimeOut 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| VibrationSensorClientAddress 属性 |

当前设备的地址信息

**命名空间：**
 [HslCommunication.Profinet.Geniitek](307ec5f8-718a-eb03-ee2d-5de6962e72cd.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public ushort Address { get; set; }
```

```
Public Property Address As UShort
	Get
	Set
```

```
public:
property unsigned short Address {
	unsigned short get ();
	void set (unsigned short value);
}
```

```
member Address : uint16 with get, set
```

#### 属性值

类型：UInt16

![](../icons/SectionExpanded.png)参见

#### 引用

[VibrationSensorClient 类](f4100855-734d-f593-67b3-891e63f61804.htm)

[HslCommunication.Profinet.Geniitek 命名空间](307ec5f8-718a-eb03-ee2d-5de6962e72cd.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## CheckSeconds 属性 

[原文連結](http://api.hslcommunication.cn/html/518e4c3b-021d-6a00-2ab4-8808fa91f5f9.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Geniitek](../html/307ec5f8-718a-eb03-ee2d-5de6962e72cd.htm "HslCommunication.Profinet.Geniitek")

[VibrationSensorClient 类](../html/f4100855-734d-f593-67b3-891e63f61804.htm "VibrationSensorClient 类")

[VibrationSensorClient 属性](../html/9da4a7a8-fc48-2eb3-3478-45cdda9f0e50.htm "VibrationSensorClient 属性")

[Address 属性](../html/08f34568-0a4b-e548-17b8-73158f4a705f.htm "Address 属性 ")

[CheckSeconds 属性](../html/518e4c3b-021d-6a00-2ab4-8808fa91f5f9.htm "CheckSeconds 属性 ")

[ConnectTimeOut 属性](../html/07b3c8d2-d1e3-da43-4742-235b02299de5.htm "ConnectTimeOut 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| VibrationSensorClientCheckSeconds 属性 |

获取或设置当前的客户端假死超时检查时间，单位为秒，默认60秒，60秒内没有接收到传感器的数据，则强制重连。

**命名空间：**
 [HslCommunication.Profinet.Geniitek](307ec5f8-718a-eb03-ee2d-5de6962e72cd.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public int CheckSeconds { get; set; }
```

```
Public Property CheckSeconds As Integer
	Get
	Set
```

```
public:
property int CheckSeconds {
	int get ();
	void set (int value);
}
```

```
member CheckSeconds : int with get, set
```

#### 属性值

类型：Int32

![](../icons/SectionExpanded.png)参见

#### 引用

[VibrationSensorClient 类](f4100855-734d-f593-67b3-891e63f61804.htm)

[HslCommunication.Profinet.Geniitek 命名空间](307ec5f8-718a-eb03-ee2d-5de6962e72cd.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ConnectTimeOut 属性 

[原文連結](http://api.hslcommunication.cn/html/07b3c8d2-d1e3-da43-4742-235b02299de5.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Geniitek](../html/307ec5f8-718a-eb03-ee2d-5de6962e72cd.htm "HslCommunication.Profinet.Geniitek")

[VibrationSensorClient 类](../html/f4100855-734d-f593-67b3-891e63f61804.htm "VibrationSensorClient 类")

[VibrationSensorClient 属性](../html/9da4a7a8-fc48-2eb3-3478-45cdda9f0e50.htm "VibrationSensorClient 属性")

[Address 属性](../html/08f34568-0a4b-e548-17b8-73158f4a705f.htm "Address 属性 ")

[CheckSeconds 属性](../html/518e4c3b-021d-6a00-2ab4-8808fa91f5f9.htm "CheckSeconds 属性 ")

[ConnectTimeOut 属性](../html/07b3c8d2-d1e3-da43-4742-235b02299de5.htm "ConnectTimeOut 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| VibrationSensorClientConnectTimeOut 属性 |

获取或设置当前客户端的连接超时时间，默认10,000毫秒，单位ms  
Gets or sets the connection timeout of the current client. The default is 10,000 milliseconds. The unit is ms.

**命名空间：**
 [HslCommunication.Profinet.Geniitek](307ec5f8-718a-eb03-ee2d-5de6962e72cd.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public int ConnectTimeOut { get; set; }
```

```
Public Property ConnectTimeOut As Integer
	Get
	Set
```

```
public:
property int ConnectTimeOut {
	int get ();
	void set (int value);
}
```

```
member ConnectTimeOut : int with get, set
```

#### 属性值

类型：Int32

![](../icons/SectionExpanded.png)参见

#### 引用

[VibrationSensorClient 类](f4100855-734d-f593-67b3-891e63f61804.htm)

[HslCommunication.Profinet.Geniitek 命名空间](307ec5f8-718a-eb03-ee2d-5de6962e72cd.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## VibrationSensorClient 方法

[原文連結](http://api.hslcommunication.cn/html/a7ca1954-ca0b-784f-4a99-2bb3ef651141.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Geniitek](../html/307ec5f8-718a-eb03-ee2d-5de6962e72cd.htm "HslCommunication.Profinet.Geniitek")

[VibrationSensorClient 类](../html/f4100855-734d-f593-67b3-891e63f61804.htm "VibrationSensorClient 类")

[VibrationSensorClient 方法](../html/a7ca1954-ca0b-784f-4a99-2bb3ef651141.htm "VibrationSensorClient 方法")

[BulidLongMessage 方法](../html/96de5e3f-4a4a-df2e-262e-ad19dc2c5b2f.htm "BulidLongMessage 方法 ")

[CheckXor 方法](../html/a6ea980c-719c-22a0-7d04-516f7b4dbb1b.htm "CheckXor 方法 ")

[ConnectClose 方法](../html/beb87b8e-1ec5-c550-7471-642cc9c798a2.htm "ConnectClose 方法 ")

[ConnectServer 方法](../html/ab05e736-dd37-a9b3-5a4b-0f267052ef77.htm "ConnectServer 方法 ")

[ConnectServerAsync 方法](../html/7a6ef23e-9bde-da7f-d8fa-2d8e2a87e6a2.htm "ConnectServerAsync 方法 ")

[SetReadActual 方法](../html/6bbeb3f4-690a-b65f-dbce-bdec86750011.htm "SetReadActual 方法 ")

[SetReadStatus 方法](../html/3e411f8c-ebb8-189f-432d-417dbda9b682.htm "SetReadStatus 方法 ")

[SetReadStatusInterval 方法](../html/958c6cda-4e99-28d9-a088-2e30f07c134a.htm "SetReadStatusInterval 方法 ")

[ToString 方法](../html/77cd262f-cbad-bdd6-ab4e-9ffedd3940a5.htm "ToString 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| VibrationSensorClient 方法 |

[VibrationSensorClient](f4100855-734d-f593-67b3-891e63f61804.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法静态成员 | [BulidLongMessage](96de5e3f-4a4a-df2e-262e-ad19dc2c5b2f.htm) | 根据地址，命令，数据，创建向传感器发送的数据信息 |
| 受保护的方法 | [CheckRemoteToken](1d231984-101b-0fff-0daf-1d18581a4241.htm) | 检查当前的头子节信息的令牌是否是正确的，仅用于某些特殊的协议实现  Check whether the token of the current header subsection information is correct, only for some special protocol implementations (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 公共方法静态成员 | [CheckXor](a6ea980c-719c-22a0-7d04-516f7b4dbb1b.htm) | 检查当前的数据是否XOR校验成功 |
| 公共方法 | [ConnectClose](beb87b8e-1ec5-c550-7471-642cc9c798a2.htm) | 关闭Mqtt服务器的连接。  Close the connection to the Mqtt server. |
| 公共方法 | [ConnectServer](ab05e736-dd37-a9b3-5a4b-0f267052ef77.htm) | 连接服务器，实例化客户端之后，至少要调用成功一次，如果返回失败，那些请过一段时间后重新调用本方法连接。  After connecting to the server, the client must be called at least once after instantiating the client. If the return fails, please call this method to connect again after a period of time. |
| 公共方法 | [ConnectServerAsync](7a6ef23e-9bde-da7f-d8fa-2d8e2a87e6a2.htm) | 连接服务器，实例化客户端之后，至少要调用成功一次，如果返回失败，那些请过一段时间后重新调用本方法连接。  After connecting to the server, the client must be called at least once after instantiating the client. If the return fails, please call this method to connect again after a period of time. |
| 受保护的方法代码示例 | [CreateSocketAndConnect(String, Int32)](a153ae27-92e3-4fe7-b2e7-c8d207bebbb0.htm) | 创建一个新的socket对象并连接到远程的地址，默认超时时间为10秒钟，需要指定ip地址以及端口号信息  Create a new socket object and connect to the remote address. The default timeout is 10 seconds. You need to specify the IP address and port number. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法代码示例 | [CreateSocketAndConnect(IPEndPoint, Int32, IPEndPoint)](7ce6288f-c6db-e968-d676-f851d2b05c77.htm) | 创建一个新的socket对象并连接到远程的地址，需要指定远程终结点，超时时间（单位是毫秒），如果需要绑定本地的IP或是端口，传入 local对象  To create a new socket object and connect to the remote address, you need to specify the remote endpoint, the timeout period (in milliseconds), if you need to bind the local IP or port, pass in the local object (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法代码示例 | [CreateSocketAndConnect(String, Int32, Int32)](d0e14a5f-1e97-e59d-078f-4fbf0412b865.htm) | 创建一个新的socket对象并连接到远程的地址，需要指定ip地址以及端口号信息，还有超时时间，单位是毫秒  To create a new socket object and connect to a remote address, you need to specify the IP address and port number information, and the timeout period in milliseconds (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法代码示例 | [CreateSocketAndConnectAsync(String, Int32)](0227bc60-1c7f-3c1d-0cd6-e49bb904d6df.htm) | 创建一个新的socket对象并连接到远程的地址，默认超时时间为10秒钟，需要指定ip地址以及端口号信息  Create a new socket object and connect to the remote address. The default timeout is 10 seconds. You need to specify the IP address and port number. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法代码示例 | [CreateSocketAndConnectAsync(IPEndPoint, Int32, IPEndPoint)](ccec07bf-cfa4-b302-210e-441821e095a4.htm) | 创建一个新的socket对象并连接到远程的地址，需要指定远程终结点，超时时间（单位是毫秒），如果需要绑定本地的IP或是端口，传入 local对象  To create a new socket object and connect to the remote address, you need to specify the remote endpoint, the timeout period (in milliseconds), if you need to bind the local IP or port, pass in the local object (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法代码示例 | [CreateSocketAndConnectAsync(String, Int32, Int32)](7ec25883-60e2-c7ac-2f90-0ea1ae0c4945.htm) | 创建一个新的socket对象并连接到远程的地址，需要指定ip地址以及端口号信息，还有超时时间，单位是毫秒  To create a new socket object and connect to a remote address, you need to specify the IP address and port number information, and the timeout period in milliseconds (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [DeleteFileByName](59afa4a4-3f3c-9831-26ee-c7c13a213bf7.htm) | 删除一个指定的文件，如果文件不存在，直接返回 True，如果文件存在则直接删除，删除成功返回 True，如果发生了异常，返回False  Delete a specified file, if the file does not exist, return True directly, if the file exists, delete it directly, if the deletion is successful, return True, if an exception occurs, return  False (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
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
| 受保护的方法 | [ReceiveFileFromSocket(Socket, Stream, ActionInt64, Int64)](a227f0b6-3050-e5d1-79c8-88ea9719bd06.htm) | [自校验] 从网络中接收一个文件，写入数据流，如果结果异常，则结束通讯，参数顺序文件名，文件大小，文件标识，上传人  [Self-checking] Receive a file from the network. If the result is abnormal, the communication ends. (继承自 [NetworkXBase](ab528a6a-e713-5202-59e0-90dd6a66de5f.htm)。) |
| 受保护的方法 | [ReceiveFileFromSocket(Socket, String, ActionInt64, Int64)](6bcd5a10-d79c-649e-fb65-cd78fecf99b7.htm) | [自校验] 从网络中接收一个文件，如果结果异常，则结束通讯  [Self-checking] Receive a file from the network. If the result is abnormal, the communication ends. (继承自 [NetworkXBase](ab528a6a-e713-5202-59e0-90dd6a66de5f.htm)。) |
| 受保护的方法 | [ReceiveFileFromSocketAsync(Socket, Stream, ActionInt64, Int64)](e3b415d2-58a9-7ab6-0097-0bdc736ca0dc.htm) | [自校验] 从网络中接收一个文件，写入数据流，如果结果异常，则结束通讯，参数顺序文件名，文件大小，文件标识，上传人  [Self-checking] Receive a file from the network. If the result is abnormal, the communication ends. (继承自 [NetworkXBase](ab528a6a-e713-5202-59e0-90dd6a66de5f.htm)。) |
| 受保护的方法 | [ReceiveFileFromSocketAsync(Socket, String, ActionInt64, Int64)](ab9b6255-580b-464b-447a-5de92b6f3dd9.htm) | [自校验] 从网络中接收一个文件，如果结果异常，则结束通讯  [Self-checking] Receive a file from the network. If the result is abnormal, the communication ends. (继承自 [NetworkXBase](ab528a6a-e713-5202-59e0-90dd6a66de5f.htm)。) |
| 受保护的方法 | [ReceiveFileHeadFromSocket](37b7d2d2-faeb-2263-627d-6a74bcbce47d.htm) | [自校验] 从套接字中接收文件头信息  [Self-checking] Receive file header information from socket (继承自 [NetworkXBase](ab528a6a-e713-5202-59e0-90dd6a66de5f.htm)。) |
| 受保护的方法 | [ReceiveFileHeadFromSocketAsync](9f10bf85-76d0-c721-60f1-1b42425a52af.htm) | [自校验] 从套接字中接收文件头信息  [Self-checking] Receive file header information from socket (继承自 [NetworkXBase](ab528a6a-e713-5202-59e0-90dd6a66de5f.htm)。) |
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
| 受保护的方法 | [SendFileAndCheckReceive(Socket, Stream, String, String, String, ActionInt64, Int64)](602977f6-83fa-486f-7b4e-7913594bb847.htm) | [自校验] 将流数据发送至套接字，具体发送细节将在继承类中实现，如果结果异常，则结束通讯  [Self-checking] Send stream data to the socket. The specific sending details will be implemented in the inherited class. If the result is abnormal, the communication will be terminated (继承自 [NetworkXBase](ab528a6a-e713-5202-59e0-90dd6a66de5f.htm)。) |
| 受保护的方法 | [SendFileAndCheckReceive(Socket, String, String, String, String, ActionInt64, Int64)](cad4c115-2f7b-5301-36e4-66b6da8d864f.htm) | [自校验] 将文件数据发送至套接字，具体发送细节将在继承类中实现，如果结果异常，则结束通讯  [Self-checking] Send the file data to the socket. The specific sending details will be implemented in the inherited class. If the result is abnormal, the communication will end (继承自 [NetworkXBase](ab528a6a-e713-5202-59e0-90dd6a66de5f.htm)。) |
| 受保护的方法 | [SendFileAndCheckReceiveAsync(Socket, Stream, String, String, String, ActionInt64, Int64)](8ad77d7e-8710-5bb5-e39a-9e613957cf72.htm) | [自校验] 将流数据发送至套接字，具体发送细节将在继承类中实现，如果结果异常，则结束通讯  [Self-checking] Send stream data to the socket. The specific sending details will be implemented in the inherited class. If the result is abnormal, the communication will be terminated (继承自 [NetworkXBase](ab528a6a-e713-5202-59e0-90dd6a66de5f.htm)。) |
| 受保护的方法 | [SendFileAndCheckReceiveAsync(Socket, String, String, String, String, ActionInt64, Int64)](9eff5d17-5f46-6900-b36b-f06e76f3605e.htm) | [自校验] 将文件数据发送至套接字，具体发送细节将在继承类中实现，如果结果异常，则结束通讯  [Self-checking] Send the file data to the socket. The specific sending details will be implemented in the inherited class. If the result is abnormal, the communication will end (继承自 [NetworkXBase](ab528a6a-e713-5202-59e0-90dd6a66de5f.htm)。) |
| 受保护的方法 | [SendFileStreamToSocket](73df4306-17c8-3649-8214-b0704cd009af.htm) | [自校验] 将文件数据发送至套接字，如果结果异常，则结束通讯  [Self-check] Send the file data to the socket. If the result is abnormal, the communication is ended. (继承自 [NetworkXBase](ab528a6a-e713-5202-59e0-90dd6a66de5f.htm)。) |
| 受保护的方法 | [SendFileStreamToSocketAsync](41b6ebdb-22c3-fe61-113b-a60f1674750c.htm) | [自校验] 将文件数据发送至套接字，如果结果异常，则结束通讯  [Self-check] Send the file data to the socket. If the result is abnormal, the communication is ended. (继承自 [NetworkXBase](ab528a6a-e713-5202-59e0-90dd6a66de5f.htm)。) |
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
| 公共方法 | [SetReadActual](6bbeb3f4-690a-b65f-dbce-bdec86750011.htm) | 设置读取震动传感器的实时加速度  Set the real-time acceleration of the vibration sensor |
| 公共方法 | [SetReadStatus](3e411f8c-ebb8-189f-432d-417dbda9b682.htm) | 设置读取震动传感器的状态数据  Set to read the status data of the shock sensor |
| 公共方法 | [SetReadStatusInterval](958c6cda-4e99-28d9-a088-2e30f07c134a.htm) | 设置当前的震动传感器的数据发送间隔为指定的时间，单位为秒  Set the current vibration sensor data transmission interval to the specified time in seconds |
| 公共方法 | [ToString](77cd262f-cbad-bdd6-ab4e-9ffedd3940a5.htm) | (重写 [NetworkXBaseToString](08bf006b-af6e-7c58-6935-8e8c387e96ab.htm).) |
| 受保护的方法 | [WriteStreamFromSocket](97bf95ea-51fe-8f2e-5da6-1c83766aaa78.htm) | 从套接字中接收所有的数据然后写入到指定的流当中去，需要指定数据的长度，支持按照百分比进行进度报告  Receives all data from the socket and writes it to the specified stream. The length of the data needs to be specified, and progress reporting is supported in percentage. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |
| 受保护的方法 | [WriteStreamFromSocketAsync](b4cf7f6d-1e0c-a266-4c49-683fdb14a922.htm) | 从套接字中接收所有的数据然后写入到指定的流当中去，需要指定数据的长度，支持按照百分比进行进度报告  Receives all data from the socket and writes it to the specified stream. The length of the data needs to be specified, and progress reporting is supported in percentage. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[VibrationSensorClient 类](f4100855-734d-f593-67b3-891e63f61804.htm)

[HslCommunication.Profinet.Geniitek 命名空间](307ec5f8-718a-eb03-ee2d-5de6962e72cd.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## BulidLongMessage 方法 

[原文連結](http://api.hslcommunication.cn/html/96de5e3f-4a4a-df2e-262e-ad19dc2c5b2f.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Geniitek](../html/307ec5f8-718a-eb03-ee2d-5de6962e72cd.htm "HslCommunication.Profinet.Geniitek")

[VibrationSensorClient 类](../html/f4100855-734d-f593-67b3-891e63f61804.htm "VibrationSensorClient 类")

[VibrationSensorClient 方法](../html/a7ca1954-ca0b-784f-4a99-2bb3ef651141.htm "VibrationSensorClient 方法")

[BulidLongMessage 方法](../html/96de5e3f-4a4a-df2e-262e-ad19dc2c5b2f.htm "BulidLongMessage 方法 ")

[CheckXor 方法](../html/a6ea980c-719c-22a0-7d04-516f7b4dbb1b.htm "CheckXor 方法 ")

[ConnectClose 方法](../html/beb87b8e-1ec5-c550-7471-642cc9c798a2.htm "ConnectClose 方法 ")

[ConnectServer 方法](../html/ab05e736-dd37-a9b3-5a4b-0f267052ef77.htm "ConnectServer 方法 ")

[ConnectServerAsync 方法](../html/7a6ef23e-9bde-da7f-d8fa-2d8e2a87e6a2.htm "ConnectServerAsync 方法 ")

[SetReadActual 方法](../html/6bbeb3f4-690a-b65f-dbce-bdec86750011.htm "SetReadActual 方法 ")

[SetReadStatus 方法](../html/3e411f8c-ebb8-189f-432d-417dbda9b682.htm "SetReadStatus 方法 ")

[SetReadStatusInterval 方法](../html/958c6cda-4e99-28d9-a088-2e30f07c134a.htm "SetReadStatusInterval 方法 ")

[ToString 方法](../html/77cd262f-cbad-bdd6-ab4e-9ffedd3940a5.htm "ToString 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| VibrationSensorClientBulidLongMessage 方法 |

根据地址，命令，数据，创建向传感器发送的数据信息

**命名空间：**
 [HslCommunication.Profinet.Geniitek](307ec5f8-718a-eb03-ee2d-5de6962e72cd.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static byte[] BulidLongMessage(
	ushort address,
	byte cmd,
	byte[] data
)
```

```
Public Shared Function BulidLongMessage ( 
	address As UShort,
	cmd As Byte,
	data As Byte()
) As Byte()
```

```
public:
static array<unsigned char>^ BulidLongMessage(
	unsigned short address, 
	unsigned char cmd, 
	array<unsigned char>^ data
)
```

```
static member BulidLongMessage : 
        address : uint16 * 
        cmd : byte * 
        data : byte[] -> byte[] 
```

#### 参数

address
:   类型：SystemUInt16  
    设备地址

cmd
:   类型：SystemByte  
    命令

data
:   类型：SystemByte  
    数据信息

#### 返回值

类型：Byte  
原始的数据内容

![](../icons/SectionExpanded.png)参见

#### 引用

[VibrationSensorClient 类](f4100855-734d-f593-67b3-891e63f61804.htm)

[HslCommunication.Profinet.Geniitek 命名空间](307ec5f8-718a-eb03-ee2d-5de6962e72cd.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## CheckXor 方法 

[原文連結](http://api.hslcommunication.cn/html/a6ea980c-719c-22a0-7d04-516f7b4dbb1b.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Geniitek](../html/307ec5f8-718a-eb03-ee2d-5de6962e72cd.htm "HslCommunication.Profinet.Geniitek")

[VibrationSensorClient 类](../html/f4100855-734d-f593-67b3-891e63f61804.htm "VibrationSensorClient 类")

[VibrationSensorClient 方法](../html/a7ca1954-ca0b-784f-4a99-2bb3ef651141.htm "VibrationSensorClient 方法")

[BulidLongMessage 方法](../html/96de5e3f-4a4a-df2e-262e-ad19dc2c5b2f.htm "BulidLongMessage 方法 ")

[CheckXor 方法](../html/a6ea980c-719c-22a0-7d04-516f7b4dbb1b.htm "CheckXor 方法 ")

[ConnectClose 方法](../html/beb87b8e-1ec5-c550-7471-642cc9c798a2.htm "ConnectClose 方法 ")

[ConnectServer 方法](../html/ab05e736-dd37-a9b3-5a4b-0f267052ef77.htm "ConnectServer 方法 ")

[ConnectServerAsync 方法](../html/7a6ef23e-9bde-da7f-d8fa-2d8e2a87e6a2.htm "ConnectServerAsync 方法 ")

[SetReadActual 方法](../html/6bbeb3f4-690a-b65f-dbce-bdec86750011.htm "SetReadActual 方法 ")

[SetReadStatus 方法](../html/3e411f8c-ebb8-189f-432d-417dbda9b682.htm "SetReadStatus 方法 ")

[SetReadStatusInterval 方法](../html/958c6cda-4e99-28d9-a088-2e30f07c134a.htm "SetReadStatusInterval 方法 ")

[ToString 方法](../html/77cd262f-cbad-bdd6-ab4e-9ffedd3940a5.htm "ToString 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| VibrationSensorClientCheckXor 方法 |

检查当前的数据是否XOR校验成功

**命名空间：**
 [HslCommunication.Profinet.Geniitek](307ec5f8-718a-eb03-ee2d-5de6962e72cd.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static bool CheckXor(
	byte[] data
)
```

```
Public Shared Function CheckXor ( 
	data As Byte()
) As Boolean
```

```
public:
static bool CheckXor(
	array<unsigned char>^ data
)
```

```
static member CheckXor : 
        data : byte[] -> bool 
```

#### 参数

data
:   类型：SystemByte  
    数据信息

#### 返回值

类型：Boolean  
校验结果

![](../icons/SectionExpanded.png)参见

#### 引用

[VibrationSensorClient 类](f4100855-734d-f593-67b3-891e63f61804.htm)

[HslCommunication.Profinet.Geniitek 命名空间](307ec5f8-718a-eb03-ee2d-5de6962e72cd.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ConnectClose 方法 

[原文連結](http://api.hslcommunication.cn/html/beb87b8e-1ec5-c550-7471-642cc9c798a2.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Geniitek](../html/307ec5f8-718a-eb03-ee2d-5de6962e72cd.htm "HslCommunication.Profinet.Geniitek")

[VibrationSensorClient 类](../html/f4100855-734d-f593-67b3-891e63f61804.htm "VibrationSensorClient 类")

[VibrationSensorClient 方法](../html/a7ca1954-ca0b-784f-4a99-2bb3ef651141.htm "VibrationSensorClient 方法")

[BulidLongMessage 方法](../html/96de5e3f-4a4a-df2e-262e-ad19dc2c5b2f.htm "BulidLongMessage 方法 ")

[CheckXor 方法](../html/a6ea980c-719c-22a0-7d04-516f7b4dbb1b.htm "CheckXor 方法 ")

[ConnectClose 方法](../html/beb87b8e-1ec5-c550-7471-642cc9c798a2.htm "ConnectClose 方法 ")

[ConnectServer 方法](../html/ab05e736-dd37-a9b3-5a4b-0f267052ef77.htm "ConnectServer 方法 ")

[ConnectServerAsync 方法](../html/7a6ef23e-9bde-da7f-d8fa-2d8e2a87e6a2.htm "ConnectServerAsync 方法 ")

[SetReadActual 方法](../html/6bbeb3f4-690a-b65f-dbce-bdec86750011.htm "SetReadActual 方法 ")

[SetReadStatus 方法](../html/3e411f8c-ebb8-189f-432d-417dbda9b682.htm "SetReadStatus 方法 ")

[SetReadStatusInterval 方法](../html/958c6cda-4e99-28d9-a088-2e30f07c134a.htm "SetReadStatusInterval 方法 ")

[ToString 方法](../html/77cd262f-cbad-bdd6-ab4e-9ffedd3940a5.htm "ToString 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| VibrationSensorClientConnectClose 方法 |

关闭Mqtt服务器的连接。  
Close the connection to the Mqtt server.

**命名空间：**
 [HslCommunication.Profinet.Geniitek](307ec5f8-718a-eb03-ee2d-5de6962e72cd.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public void ConnectClose()
```

```
Public Sub ConnectClose
```

```
public:
void ConnectClose()
```

```
member ConnectClose : unit -> unit 
```

![](../icons/SectionExpanded.png)参见

#### 引用

[VibrationSensorClient 类](f4100855-734d-f593-67b3-891e63f61804.htm)

[HslCommunication.Profinet.Geniitek 命名空间](307ec5f8-718a-eb03-ee2d-5de6962e72cd.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ConnectServer 方法 

[原文連結](http://api.hslcommunication.cn/html/ab05e736-dd37-a9b3-5a4b-0f267052ef77.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Geniitek](../html/307ec5f8-718a-eb03-ee2d-5de6962e72cd.htm "HslCommunication.Profinet.Geniitek")

[VibrationSensorClient 类](../html/f4100855-734d-f593-67b3-891e63f61804.htm "VibrationSensorClient 类")

[VibrationSensorClient 方法](../html/a7ca1954-ca0b-784f-4a99-2bb3ef651141.htm "VibrationSensorClient 方法")

[BulidLongMessage 方法](../html/96de5e3f-4a4a-df2e-262e-ad19dc2c5b2f.htm "BulidLongMessage 方法 ")

[CheckXor 方法](../html/a6ea980c-719c-22a0-7d04-516f7b4dbb1b.htm "CheckXor 方法 ")

[ConnectClose 方法](../html/beb87b8e-1ec5-c550-7471-642cc9c798a2.htm "ConnectClose 方法 ")

[ConnectServer 方法](../html/ab05e736-dd37-a9b3-5a4b-0f267052ef77.htm "ConnectServer 方法 ")

[ConnectServerAsync 方法](../html/7a6ef23e-9bde-da7f-d8fa-2d8e2a87e6a2.htm "ConnectServerAsync 方法 ")

[SetReadActual 方法](../html/6bbeb3f4-690a-b65f-dbce-bdec86750011.htm "SetReadActual 方法 ")

[SetReadStatus 方法](../html/3e411f8c-ebb8-189f-432d-417dbda9b682.htm "SetReadStatus 方法 ")

[SetReadStatusInterval 方法](../html/958c6cda-4e99-28d9-a088-2e30f07c134a.htm "SetReadStatusInterval 方法 ")

[ToString 方法](../html/77cd262f-cbad-bdd6-ab4e-9ffedd3940a5.htm "ToString 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| VibrationSensorClientConnectServer 方法 |

连接服务器，实例化客户端之后，至少要调用成功一次，如果返回失败，那些请过一段时间后重新调用本方法连接。  
After connecting to the server, the client must be called at least once after instantiating the client.
If the return fails, please call this method to connect again after a period of time.

**命名空间：**
 [HslCommunication.Profinet.Geniitek](307ec5f8-718a-eb03-ee2d-5de6962e72cd.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public OperateResult ConnectServer()
```

```
Public Function ConnectServer As OperateResult
```

```
public:
OperateResult^ ConnectServer()
```

```
member ConnectServer : unit -> OperateResult 
```

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
连接是否成功

![](../icons/SectionExpanded.png)参见

#### 引用

[VibrationSensorClient 类](f4100855-734d-f593-67b3-891e63f61804.htm)

[HslCommunication.Profinet.Geniitek 命名空间](307ec5f8-718a-eb03-ee2d-5de6962e72cd.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ConnectServerAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/7a6ef23e-9bde-da7f-d8fa-2d8e2a87e6a2.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Geniitek](../html/307ec5f8-718a-eb03-ee2d-5de6962e72cd.htm "HslCommunication.Profinet.Geniitek")

[VibrationSensorClient 类](../html/f4100855-734d-f593-67b3-891e63f61804.htm "VibrationSensorClient 类")

[VibrationSensorClient 方法](../html/a7ca1954-ca0b-784f-4a99-2bb3ef651141.htm "VibrationSensorClient 方法")

[BulidLongMessage 方法](../html/96de5e3f-4a4a-df2e-262e-ad19dc2c5b2f.htm "BulidLongMessage 方法 ")

[CheckXor 方法](../html/a6ea980c-719c-22a0-7d04-516f7b4dbb1b.htm "CheckXor 方法 ")

[ConnectClose 方法](../html/beb87b8e-1ec5-c550-7471-642cc9c798a2.htm "ConnectClose 方法 ")

[ConnectServer 方法](../html/ab05e736-dd37-a9b3-5a4b-0f267052ef77.htm "ConnectServer 方法 ")

[ConnectServerAsync 方法](../html/7a6ef23e-9bde-da7f-d8fa-2d8e2a87e6a2.htm "ConnectServerAsync 方法 ")

[SetReadActual 方法](../html/6bbeb3f4-690a-b65f-dbce-bdec86750011.htm "SetReadActual 方法 ")

[SetReadStatus 方法](../html/3e411f8c-ebb8-189f-432d-417dbda9b682.htm "SetReadStatus 方法 ")

[SetReadStatusInterval 方法](../html/958c6cda-4e99-28d9-a088-2e30f07c134a.htm "SetReadStatusInterval 方法 ")

[ToString 方法](../html/77cd262f-cbad-bdd6-ab4e-9ffedd3940a5.htm "ToString 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| VibrationSensorClientConnectServerAsync 方法 |

连接服务器，实例化客户端之后，至少要调用成功一次，如果返回失败，那些请过一段时间后重新调用本方法连接。  
After connecting to the server, the client must be called at least once after instantiating the client.
If the return fails, please call this method to connect again after a period of time.

**命名空间：**
 [HslCommunication.Profinet.Geniitek](307ec5f8-718a-eb03-ee2d-5de6962e72cd.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public Task<OperateResult> ConnectServerAsync()
```

```
Public Function ConnectServerAsync As Task(Of OperateResult)
```

```
public:
Task<OperateResult^>^ ConnectServerAsync()
```

```
member ConnectServerAsync : unit -> Task<OperateResult> 
```

#### 返回值

类型：Task[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
连接是否成功

![](../icons/SectionExpanded.png)参见

#### 引用

[VibrationSensorClient 类](f4100855-734d-f593-67b3-891e63f61804.htm)

[HslCommunication.Profinet.Geniitek 命名空间](307ec5f8-718a-eb03-ee2d-5de6962e72cd.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## SetReadActual 方法 

[原文連結](http://api.hslcommunication.cn/html/6bbeb3f4-690a-b65f-dbce-bdec86750011.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Geniitek](../html/307ec5f8-718a-eb03-ee2d-5de6962e72cd.htm "HslCommunication.Profinet.Geniitek")

[VibrationSensorClient 类](../html/f4100855-734d-f593-67b3-891e63f61804.htm "VibrationSensorClient 类")

[VibrationSensorClient 方法](../html/a7ca1954-ca0b-784f-4a99-2bb3ef651141.htm "VibrationSensorClient 方法")

[BulidLongMessage 方法](../html/96de5e3f-4a4a-df2e-262e-ad19dc2c5b2f.htm "BulidLongMessage 方法 ")

[CheckXor 方法](../html/a6ea980c-719c-22a0-7d04-516f7b4dbb1b.htm "CheckXor 方法 ")

[ConnectClose 方法](../html/beb87b8e-1ec5-c550-7471-642cc9c798a2.htm "ConnectClose 方法 ")

[ConnectServer 方法](../html/ab05e736-dd37-a9b3-5a4b-0f267052ef77.htm "ConnectServer 方法 ")

[ConnectServerAsync 方法](../html/7a6ef23e-9bde-da7f-d8fa-2d8e2a87e6a2.htm "ConnectServerAsync 方法 ")

[SetReadActual 方法](../html/6bbeb3f4-690a-b65f-dbce-bdec86750011.htm "SetReadActual 方法 ")

[SetReadStatus 方法](../html/3e411f8c-ebb8-189f-432d-417dbda9b682.htm "SetReadStatus 方法 ")

[SetReadStatusInterval 方法](../html/958c6cda-4e99-28d9-a088-2e30f07c134a.htm "SetReadStatusInterval 方法 ")

[ToString 方法](../html/77cd262f-cbad-bdd6-ab4e-9ffedd3940a5.htm "ToString 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| VibrationSensorClientSetReadActual 方法 |

设置读取震动传感器的实时加速度  
Set the real-time acceleration of the vibration sensor

**命名空间：**
 [HslCommunication.Profinet.Geniitek](307ec5f8-718a-eb03-ee2d-5de6962e72cd.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public OperateResult SetReadActual()
```

```
Public Function SetReadActual As OperateResult
```

```
public:
OperateResult^ SetReadActual()
```

```
member SetReadActual : unit -> OperateResult 
```

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否发送成功

![](../icons/SectionExpanded.png)参见

#### 引用

[VibrationSensorClient 类](f4100855-734d-f593-67b3-891e63f61804.htm)

[HslCommunication.Profinet.Geniitek 命名空间](307ec5f8-718a-eb03-ee2d-5de6962e72cd.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## SetReadStatus 方法 

[原文連結](http://api.hslcommunication.cn/html/3e411f8c-ebb8-189f-432d-417dbda9b682.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Geniitek](../html/307ec5f8-718a-eb03-ee2d-5de6962e72cd.htm "HslCommunication.Profinet.Geniitek")

[VibrationSensorClient 类](../html/f4100855-734d-f593-67b3-891e63f61804.htm "VibrationSensorClient 类")

[VibrationSensorClient 方法](../html/a7ca1954-ca0b-784f-4a99-2bb3ef651141.htm "VibrationSensorClient 方法")

[BulidLongMessage 方法](../html/96de5e3f-4a4a-df2e-262e-ad19dc2c5b2f.htm "BulidLongMessage 方法 ")

[CheckXor 方法](../html/a6ea980c-719c-22a0-7d04-516f7b4dbb1b.htm "CheckXor 方法 ")

[ConnectClose 方法](../html/beb87b8e-1ec5-c550-7471-642cc9c798a2.htm "ConnectClose 方法 ")

[ConnectServer 方法](../html/ab05e736-dd37-a9b3-5a4b-0f267052ef77.htm "ConnectServer 方法 ")

[ConnectServerAsync 方法](../html/7a6ef23e-9bde-da7f-d8fa-2d8e2a87e6a2.htm "ConnectServerAsync 方法 ")

[SetReadActual 方法](../html/6bbeb3f4-690a-b65f-dbce-bdec86750011.htm "SetReadActual 方法 ")

[SetReadStatus 方法](../html/3e411f8c-ebb8-189f-432d-417dbda9b682.htm "SetReadStatus 方法 ")

[SetReadStatusInterval 方法](../html/958c6cda-4e99-28d9-a088-2e30f07c134a.htm "SetReadStatusInterval 方法 ")

[ToString 方法](../html/77cd262f-cbad-bdd6-ab4e-9ffedd3940a5.htm "ToString 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| VibrationSensorClientSetReadStatus 方法 |

设置读取震动传感器的状态数据  
Set to read the status data of the shock sensor

**命名空间：**
 [HslCommunication.Profinet.Geniitek](307ec5f8-718a-eb03-ee2d-5de6962e72cd.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public OperateResult SetReadStatus()
```

```
Public Function SetReadStatus As OperateResult
```

```
public:
OperateResult^ SetReadStatus()
```

```
member SetReadStatus : unit -> OperateResult 
```

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否发送成功

![](../icons/SectionExpanded.png)参见

#### 引用

[VibrationSensorClient 类](f4100855-734d-f593-67b3-891e63f61804.htm)

[HslCommunication.Profinet.Geniitek 命名空间](307ec5f8-718a-eb03-ee2d-5de6962e72cd.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## SetReadStatusInterval 方法 

[原文連結](http://api.hslcommunication.cn/html/958c6cda-4e99-28d9-a088-2e30f07c134a.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Geniitek](../html/307ec5f8-718a-eb03-ee2d-5de6962e72cd.htm "HslCommunication.Profinet.Geniitek")

[VibrationSensorClient 类](../html/f4100855-734d-f593-67b3-891e63f61804.htm "VibrationSensorClient 类")

[VibrationSensorClient 方法](../html/a7ca1954-ca0b-784f-4a99-2bb3ef651141.htm "VibrationSensorClient 方法")

[BulidLongMessage 方法](../html/96de5e3f-4a4a-df2e-262e-ad19dc2c5b2f.htm "BulidLongMessage 方法 ")

[CheckXor 方法](../html/a6ea980c-719c-22a0-7d04-516f7b4dbb1b.htm "CheckXor 方法 ")

[ConnectClose 方法](../html/beb87b8e-1ec5-c550-7471-642cc9c798a2.htm "ConnectClose 方法 ")

[ConnectServer 方法](../html/ab05e736-dd37-a9b3-5a4b-0f267052ef77.htm "ConnectServer 方法 ")

[ConnectServerAsync 方法](../html/7a6ef23e-9bde-da7f-d8fa-2d8e2a87e6a2.htm "ConnectServerAsync 方法 ")

[SetReadActual 方法](../html/6bbeb3f4-690a-b65f-dbce-bdec86750011.htm "SetReadActual 方法 ")

[SetReadStatus 方法](../html/3e411f8c-ebb8-189f-432d-417dbda9b682.htm "SetReadStatus 方法 ")

[SetReadStatusInterval 方法](../html/958c6cda-4e99-28d9-a088-2e30f07c134a.htm "SetReadStatusInterval 方法 ")

[ToString 方法](../html/77cd262f-cbad-bdd6-ab4e-9ffedd3940a5.htm "ToString 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| VibrationSensorClientSetReadStatusInterval 方法 |

设置当前的震动传感器的数据发送间隔为指定的时间，单位为秒  
Set the current vibration sensor data transmission interval to the specified time in seconds

**命名空间：**
 [HslCommunication.Profinet.Geniitek](307ec5f8-718a-eb03-ee2d-5de6962e72cd.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public OperateResult SetReadStatusInterval(
	int seconds
)
```

```
Public Function SetReadStatusInterval ( 
	seconds As Integer
) As OperateResult
```

```
public:
OperateResult^ SetReadStatusInterval(
	int seconds
)
```

```
member SetReadStatusInterval : 
        seconds : int -> OperateResult 
```

#### 参数

seconds
:   类型：SystemInt32  
    时间信息，单位为秒

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否发送成功

![](../icons/SectionExpanded.png)参见

#### 引用

[VibrationSensorClient 类](f4100855-734d-f593-67b3-891e63f61804.htm)

[HslCommunication.Profinet.Geniitek 命名空间](307ec5f8-718a-eb03-ee2d-5de6962e72cd.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ToString 方法 

[原文連結](http://api.hslcommunication.cn/html/77cd262f-cbad-bdd6-ab4e-9ffedd3940a5.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Geniitek](../html/307ec5f8-718a-eb03-ee2d-5de6962e72cd.htm "HslCommunication.Profinet.Geniitek")

[VibrationSensorClient 类](../html/f4100855-734d-f593-67b3-891e63f61804.htm "VibrationSensorClient 类")

[VibrationSensorClient 方法](../html/a7ca1954-ca0b-784f-4a99-2bb3ef651141.htm "VibrationSensorClient 方法")

[BulidLongMessage 方法](../html/96de5e3f-4a4a-df2e-262e-ad19dc2c5b2f.htm "BulidLongMessage 方法 ")

[CheckXor 方法](../html/a6ea980c-719c-22a0-7d04-516f7b4dbb1b.htm "CheckXor 方法 ")

[ConnectClose 方法](../html/beb87b8e-1ec5-c550-7471-642cc9c798a2.htm "ConnectClose 方法 ")

[ConnectServer 方法](../html/ab05e736-dd37-a9b3-5a4b-0f267052ef77.htm "ConnectServer 方法 ")

[ConnectServerAsync 方法](../html/7a6ef23e-9bde-da7f-d8fa-2d8e2a87e6a2.htm "ConnectServerAsync 方法 ")

[SetReadActual 方法](../html/6bbeb3f4-690a-b65f-dbce-bdec86750011.htm "SetReadActual 方法 ")

[SetReadStatus 方法](../html/3e411f8c-ebb8-189f-432d-417dbda9b682.htm "SetReadStatus 方法 ")

[SetReadStatusInterval 方法](../html/958c6cda-4e99-28d9-a088-2e30f07c134a.htm "SetReadStatusInterval 方法 ")

[ToString 方法](../html/77cd262f-cbad-bdd6-ab4e-9ffedd3940a5.htm "ToString 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| VibrationSensorClientToString 方法 |

[缺少 "M:HslCommunication.Profinet.Geniitek.VibrationSensorClient.ToString" 的 <summary> 文档]

**命名空间：**
 [HslCommunication.Profinet.Geniitek](307ec5f8-718a-eb03-ee2d-5de6962e72cd.htm)  
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

[缺少 "M:HslCommunication.Profinet.Geniitek.VibrationSensorClient.ToString" 的 <returns> 文档]

![](../icons/SectionExpanded.png)参见

#### 引用

[VibrationSensorClient 类](f4100855-734d-f593-67b3-891e63f61804.htm)

[HslCommunication.Profinet.Geniitek 命名空间](307ec5f8-718a-eb03-ee2d-5de6962e72cd.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## VibrationSensorClient 事件

[原文連結](http://api.hslcommunication.cn/html/edc6c96c-eab1-7e9e-e1c9-a79c676754a5.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Geniitek](../html/307ec5f8-718a-eb03-ee2d-5de6962e72cd.htm "HslCommunication.Profinet.Geniitek")

[VibrationSensorClient 类](../html/f4100855-734d-f593-67b3-891e63f61804.htm "VibrationSensorClient 类")

[VibrationSensorClient 事件](../html/edc6c96c-eab1-7e9e-e1c9-a79c676754a5.htm "VibrationSensorClient 事件")

[OnActualValueReceive 事件](../html/b88d88ef-cb1e-80db-732c-2cd1f794d45b.htm "OnActualValueReceive 事件")

[OnClientConnected 事件](../html/e371eef0-ea60-a452-ed75-437958da9036.htm "OnClientConnected 事件")

[OnNetworkError 事件](../html/51179b61-d3a9-44dc-cf21-e33c6beb87fe.htm "OnNetworkError 事件")

[OnPeekValueReceive 事件](../html/9ce2d21b-57a6-103e-3800-95d0af9f300b.htm "OnPeekValueReceive 事件")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| VibrationSensorClient 事件 |

[VibrationSensorClient](f4100855-734d-f593-67b3-891e63f61804.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)事件

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共事件 | [OnActualValueReceive](b88d88ef-cb1e-80db-732c-2cd1f794d45b.htm) | 接收到震动传感器实时数据时触发  Triggered when real-time data from shock sensor is received |
| 公共事件 | [OnClientConnected](e371eef0-ea60-a452-ed75-437958da9036.htm) | 当客户端连接成功触发事件，就算是重新连接服务器后，也是会触发的  The event is triggered when the client is connected successfully, even after reconnecting to the server. |
| 公共事件 | [OnNetworkError](51179b61-d3a9-44dc-cf21-e33c6beb87fe.htm) | 当网络发生异常的时候触发的事件，用户应该在事件里进行重连服务器 |
| 公共事件 | [OnPeekValueReceive](9ce2d21b-57a6-103e-3800-95d0af9f300b.htm) | 接收到震动传感器峰值数据时触发  Triggered when peak data of vibration sensor is received |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[VibrationSensorClient 类](f4100855-734d-f593-67b3-891e63f61804.htm)

[HslCommunication.Profinet.Geniitek 命名空间](307ec5f8-718a-eb03-ee2d-5de6962e72cd.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## OnActualValueReceive 事件

[原文連結](http://api.hslcommunication.cn/html/b88d88ef-cb1e-80db-732c-2cd1f794d45b.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Geniitek](../html/307ec5f8-718a-eb03-ee2d-5de6962e72cd.htm "HslCommunication.Profinet.Geniitek")

[VibrationSensorClient 类](../html/f4100855-734d-f593-67b3-891e63f61804.htm "VibrationSensorClient 类")

[VibrationSensorClient 事件](../html/edc6c96c-eab1-7e9e-e1c9-a79c676754a5.htm "VibrationSensorClient 事件")

[OnActualValueReceive 事件](../html/b88d88ef-cb1e-80db-732c-2cd1f794d45b.htm "OnActualValueReceive 事件")

[OnClientConnected 事件](../html/e371eef0-ea60-a452-ed75-437958da9036.htm "OnClientConnected 事件")

[OnNetworkError 事件](../html/51179b61-d3a9-44dc-cf21-e33c6beb87fe.htm "OnNetworkError 事件")

[OnPeekValueReceive 事件](../html/9ce2d21b-57a6-103e-3800-95d0af9f300b.htm "OnPeekValueReceive 事件")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| VibrationSensorClientOnActualValueReceive 事件 |

接收到震动传感器实时数据时触发  
Triggered when real-time data from shock sensor is received

**命名空间：**
 [HslCommunication.Profinet.Geniitek](307ec5f8-718a-eb03-ee2d-5de6962e72cd.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public event VibrationSensorClientOnActualValueReceiveDelegate OnActualValueReceive
```

```
Public Event OnActualValueReceive As VibrationSensorClientOnActualValueReceiveDelegate
```

```
public:
 event VibrationSensorClientOnActualValueReceiveDelegate^ OnActualValueReceive {
	void add (VibrationSensorClientOnActualValueReceiveDelegate^ value);
	void remove (VibrationSensorClientOnActualValueReceiveDelegate^ value);
}
```

```
member OnActualValueReceive : IEvent<VibrationSensorClientOnActualValueReceiveDelegate,
    EventArgs>
```

#### 值

类型：[HslCommunication.Profinet.GeniitekVibrationSensorClientOnActualValueReceiveDelegate](25739141-a12a-ca22-7443-62162c2c19d9.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[VibrationSensorClient 类](f4100855-734d-f593-67b3-891e63f61804.htm)

[HslCommunication.Profinet.Geniitek 命名空间](307ec5f8-718a-eb03-ee2d-5de6962e72cd.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## OnClientConnected 事件

[原文連結](http://api.hslcommunication.cn/html/e371eef0-ea60-a452-ed75-437958da9036.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Geniitek](../html/307ec5f8-718a-eb03-ee2d-5de6962e72cd.htm "HslCommunication.Profinet.Geniitek")

[VibrationSensorClient 类](../html/f4100855-734d-f593-67b3-891e63f61804.htm "VibrationSensorClient 类")

[VibrationSensorClient 事件](../html/edc6c96c-eab1-7e9e-e1c9-a79c676754a5.htm "VibrationSensorClient 事件")

[OnActualValueReceive 事件](../html/b88d88ef-cb1e-80db-732c-2cd1f794d45b.htm "OnActualValueReceive 事件")

[OnClientConnected 事件](../html/e371eef0-ea60-a452-ed75-437958da9036.htm "OnClientConnected 事件")

[OnNetworkError 事件](../html/51179b61-d3a9-44dc-cf21-e33c6beb87fe.htm "OnNetworkError 事件")

[OnPeekValueReceive 事件](../html/9ce2d21b-57a6-103e-3800-95d0af9f300b.htm "OnPeekValueReceive 事件")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| VibrationSensorClientOnClientConnected 事件 |

当客户端连接成功触发事件，就算是重新连接服务器后，也是会触发的  
The event is triggered when the client is connected successfully, even after reconnecting to the server.

**命名空间：**
 [HslCommunication.Profinet.Geniitek](307ec5f8-718a-eb03-ee2d-5de6962e72cd.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public event VibrationSensorClientOnClientConnectedDelegate OnClientConnected
```

```
Public Event OnClientConnected As VibrationSensorClientOnClientConnectedDelegate
```

```
public:
 event VibrationSensorClientOnClientConnectedDelegate^ OnClientConnected {
	void add (VibrationSensorClientOnClientConnectedDelegate^ value);
	void remove (VibrationSensorClientOnClientConnectedDelegate^ value);
}
```

```
member OnClientConnected : IEvent<VibrationSensorClientOnClientConnectedDelegate,
    EventArgs>
```

#### 值

类型：[HslCommunication.Profinet.GeniitekVibrationSensorClientOnClientConnectedDelegate](ff07e8ff-aac9-0fc1-08a1-d6878b6f40ea.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[VibrationSensorClient 类](f4100855-734d-f593-67b3-891e63f61804.htm)

[HslCommunication.Profinet.Geniitek 命名空间](307ec5f8-718a-eb03-ee2d-5de6962e72cd.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## OnNetworkError 事件

[原文連結](http://api.hslcommunication.cn/html/51179b61-d3a9-44dc-cf21-e33c6beb87fe.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Geniitek](../html/307ec5f8-718a-eb03-ee2d-5de6962e72cd.htm "HslCommunication.Profinet.Geniitek")

[VibrationSensorClient 类](../html/f4100855-734d-f593-67b3-891e63f61804.htm "VibrationSensorClient 类")

[VibrationSensorClient 事件](../html/edc6c96c-eab1-7e9e-e1c9-a79c676754a5.htm "VibrationSensorClient 事件")

[OnActualValueReceive 事件](../html/b88d88ef-cb1e-80db-732c-2cd1f794d45b.htm "OnActualValueReceive 事件")

[OnClientConnected 事件](../html/e371eef0-ea60-a452-ed75-437958da9036.htm "OnClientConnected 事件")

[OnNetworkError 事件](../html/51179b61-d3a9-44dc-cf21-e33c6beb87fe.htm "OnNetworkError 事件")

[OnPeekValueReceive 事件](../html/9ce2d21b-57a6-103e-3800-95d0af9f300b.htm "OnPeekValueReceive 事件")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| VibrationSensorClientOnNetworkError 事件 |

当网络发生异常的时候触发的事件，用户应该在事件里进行重连服务器

**命名空间：**
 [HslCommunication.Profinet.Geniitek](307ec5f8-718a-eb03-ee2d-5de6962e72cd.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public event EventHandler OnNetworkError
```

```
Public Event OnNetworkError As EventHandler
```

```
public:
 event EventHandler^ OnNetworkError {
	void add (EventHandler^ value);
	void remove (EventHandler^ value);
}
```

```
member OnNetworkError : IEvent<EventHandler,
    EventArgs>
```

#### 值

类型：SystemEventHandler

![](../icons/SectionExpanded.png)参见

#### 引用

[VibrationSensorClient 类](f4100855-734d-f593-67b3-891e63f61804.htm)

[HslCommunication.Profinet.Geniitek 命名空间](307ec5f8-718a-eb03-ee2d-5de6962e72cd.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## OnPeekValueReceive 事件

[原文連結](http://api.hslcommunication.cn/html/9ce2d21b-57a6-103e-3800-95d0af9f300b.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Geniitek](../html/307ec5f8-718a-eb03-ee2d-5de6962e72cd.htm "HslCommunication.Profinet.Geniitek")

[VibrationSensorClient 类](../html/f4100855-734d-f593-67b3-891e63f61804.htm "VibrationSensorClient 类")

[VibrationSensorClient 事件](../html/edc6c96c-eab1-7e9e-e1c9-a79c676754a5.htm "VibrationSensorClient 事件")

[OnActualValueReceive 事件](../html/b88d88ef-cb1e-80db-732c-2cd1f794d45b.htm "OnActualValueReceive 事件")

[OnClientConnected 事件](../html/e371eef0-ea60-a452-ed75-437958da9036.htm "OnClientConnected 事件")

[OnNetworkError 事件](../html/51179b61-d3a9-44dc-cf21-e33c6beb87fe.htm "OnNetworkError 事件")

[OnPeekValueReceive 事件](../html/9ce2d21b-57a6-103e-3800-95d0af9f300b.htm "OnPeekValueReceive 事件")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| VibrationSensorClientOnPeekValueReceive 事件 |

接收到震动传感器峰值数据时触发  
Triggered when peak data of vibration sensor is received

**命名空间：**
 [HslCommunication.Profinet.Geniitek](307ec5f8-718a-eb03-ee2d-5de6962e72cd.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public event VibrationSensorClientOnPeekValueReceiveDelegate OnPeekValueReceive
```

```
Public Event OnPeekValueReceive As VibrationSensorClientOnPeekValueReceiveDelegate
```

```
public:
 event VibrationSensorClientOnPeekValueReceiveDelegate^ OnPeekValueReceive {
	void add (VibrationSensorClientOnPeekValueReceiveDelegate^ value);
	void remove (VibrationSensorClientOnPeekValueReceiveDelegate^ value);
}
```

```
member OnPeekValueReceive : IEvent<VibrationSensorClientOnPeekValueReceiveDelegate,
    EventArgs>
```

#### 值

类型：[HslCommunication.Profinet.GeniitekVibrationSensorClientOnPeekValueReceiveDelegate](91ea5b2f-d211-3f51-f32c-d73a03c83ca6.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[VibrationSensorClient 类](f4100855-734d-f593-67b3-891e63f61804.htm)

[HslCommunication.Profinet.Geniitek 命名空间](307ec5f8-718a-eb03-ee2d-5de6962e72cd.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## VibrationSensorClient 字段

[原文連結](http://api.hslcommunication.cn/html/50704ab2-06f9-69f8-6024-2715d539a965.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Geniitek](../html/307ec5f8-718a-eb03-ee2d-5de6962e72cd.htm "HslCommunication.Profinet.Geniitek")

[VibrationSensorClient 类](../html/f4100855-734d-f593-67b3-891e63f61804.htm "VibrationSensorClient 类")

[VibrationSensorClient 构造函数](../html/d8565314-1894-17a8-a16a-9776a924cb8a.htm "VibrationSensorClient 构造函数 ")

[VibrationSensorClient 属性](../html/9da4a7a8-fc48-2eb3-3478-45cdda9f0e50.htm "VibrationSensorClient 属性")

[VibrationSensorClient 方法](../html/a7ca1954-ca0b-784f-4a99-2bb3ef651141.htm "VibrationSensorClient 方法")

[VibrationSensorClient 事件](../html/edc6c96c-eab1-7e9e-e1c9-a79c676754a5.htm "VibrationSensorClient 事件")

[VibrationSensorClient 字段](../html/50704ab2-06f9-69f8-6024-2715d539a965.htm "VibrationSensorClient 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| VibrationSensorClient 字段 |

[VibrationSensorClient](f4100855-734d-f593-67b3-891e63f61804.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)字段

|  | 名称 | 说明 |
| --- | --- | --- |
| 受保护的字段 | [CoreSocket](6831da9f-21e0-8967-f17c-641b330e0d80.htm) | 对客户端而言是的通讯用的套接字，对服务器来说是用于侦听的套接字  A communication socket for the client, or a listening socket for the server (继承自 [NetworkXBase](ab528a6a-e713-5202-59e0-90dd6a66de5f.htm)。) |
| 受保护的字段 | [fileCacheSize](2328fb89-cfda-7fa9-88d0-98f1ae3c7c54.htm) | 文件传输的时候的缓存大小，直接影响传输的速度，值越大，传输速度越快，越占内存，默认为100K大小  The size of the cache during file transfer directly affects the speed of the transfer. The larger the value, the faster the transfer speed and the more memory it takes. The default size is 100K. (继承自 [NetworkBase](783ebe14-9f50-5b88-f0eb-5ea5ef6b2c56.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[VibrationSensorClient 类](f4100855-734d-f593-67b3-891e63f61804.htm)

[HslCommunication.Profinet.Geniitek 命名空间](307ec5f8-718a-eb03-ee2d-5de6962e72cd.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)