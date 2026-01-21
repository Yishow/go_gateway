# HslCommunication - HslCommunication.Robot.Hyundai

> 分類頁數: 30



---
## HslCommunication.Robot.Hyundai

[原文連結](http://api.hslcommunication.cn/html/caddf3fe-6919-18da-f3ee-fc10fac8c1bf.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.Hyundai](../html/caddf3fe-6919-18da-f3ee-fc10fac8c1bf.htm "HslCommunication.Robot.Hyundai")

[HyundaiData 类](../html/86650e61-38f2-c3ee-2fcd-c45d8d120f9f.htm "HyundaiData 类")

[HyundaiUdpNet 类](../html/63e7fe67-d9ec-efd5-a7d8-a8cde49538e2.htm "HyundaiUdpNet 类")

[HyundaiUdpNet.OnHyundaiMessageReceiveDelegate 委托](../html/5f6a905d-0c0e-42cb-c1d2-51e9812169db.htm "HyundaiUdpNet.OnHyundaiMessageReceiveDelegate 委托")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslCommunication.Robot.Hyundai 命名空间 |

[缺少 "N:HslCommunication.Robot.Hyundai" 的 <summary> 文档]

![](../icons/SectionExpanded.png)类

|  | 类 | 说明 |
| --- | --- | --- |
| 公共类 | [HyundaiData](86650e61-38f2-c3ee-2fcd-c45d8d120f9f.htm) | Hyundai的数据类对象 |
| 公共类 | [HyundaiUdpNet](63e7fe67-d9ec-efd5-a7d8-a8cde49538e2.htm) | 现代机器人的UDP通讯类，注意本类是服务器，需要等待机器人先配置好ip地址及端口，然后连接到本服务器才能正确的进行操作。详细参见api文档注释  The UDP communication class of modern robots. Note that this class is a server. You need to wait for the robot to configure the IP address and port first, and then connect to this server to operate correctly. See api documentation for details |

![](../icons/SectionExpanded.png)委托

|  | 委托 | 说明 |
| --- | --- | --- |
| 公共委托 | [HyundaiUdpNetOnHyundaiMessageReceiveDelegate](5f6a905d-0c0e-42cb-c1d2-51e9812169db.htm) | 收到机器人消息的事件委托 |

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## HyundaiData 类

[原文連結](http://api.hslcommunication.cn/html/86650e61-38f2-c3ee-2fcd-c45d8d120f9f.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.Hyundai](../html/caddf3fe-6919-18da-f3ee-fc10fac8c1bf.htm "HslCommunication.Robot.Hyundai")

[HyundaiData 类](../html/86650e61-38f2-c3ee-2fcd-c45d8d120f9f.htm "HyundaiData 类")

[HyundaiData 构造函数](../html/0b9a7201-f20f-38a1-e8cb-cb2d84214806.htm "HyundaiData 构造函数 ")

[HyundaiData 属性](../html/50ed49dd-87ac-204e-6164-898e5ef2fb31.htm "HyundaiData 属性")

[HyundaiData 方法](../html/c8612406-272d-33dc-b404-b0008b3b2767.htm "HyundaiData 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HyundaiData 类 |

Hyundai的数据类对象

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  HslCommunication.Robot.HyundaiHyundaiData

**命名空间：**
 [HslCommunication.Robot.Hyundai](caddf3fe-6919-18da-f3ee-fc10fac8c1bf.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class HyundaiData
```

```
Public Class HyundaiData
```

```
public ref class HyundaiData
```

```
type HyundaiData =  class end
```

HyundaiData 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [HyundaiData](fea8a726-d649-ca0f-cd42-dcbd1f6fe18f.htm) | 实例化一个默认的对象 |
| 公共方法 | [HyundaiData(Byte)](25014cc4-a3dc-7712-2f78-d70e3a0b9db6.htm) | 通过缓存对象实例化一个 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [CharDummy](e5796267-5048-c787-dae9-f0a6ee45c0fa.htm) | 虚标记 |
| 公共属性 | [Command](aa7fd03a-39c0-44e7-1282-7bf86bc1793b.htm) | 命令码，从控制器发数据到PC和PC到控制器，两者的命令不一样 |
| 公共属性 | [Count](3a60c870-22e5-444a-70a0-60e32e2e8ae6.htm) | 标记数据，从PLC发送给机器人的数据，原封不动的返回 |
| 公共属性 | [Data](aa511f8a-78ce-7d57-ec81-e0ca1425d3d5.htm) | 关节坐标数据，包含X,Y,Z,W,P,R，三个位置数据，三个角度数据。 |
| 公共属性 | [IntDummy](b42a4d50-2f24-6098-5def-89a801d38e6e.htm) | 虚标记 |
| 公共属性 | [State](657c8050-8cce-cc62-40ce-2c7ef9df1bef.htm) | 状态码 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法 | GetType | (继承自 Object。) |
| 公共方法 | [LoadBy](34da99e5-6c2c-0365-410c-34846c0bc90e.htm) | 从字节数组的指定索引开始加载现在机器人的数据 |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法 | [ToBytes](3cc5be93-38ea-4f3b-a338-620464accd37.htm) | 将现代机器人的数据转换为字节数组 |
| 公共方法 | [ToString](8905a222-5ea4-ed24-547e-33f2ca530baa.htm) | (重写 ObjectToString.) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.Robot.Hyundai 命名空间](caddf3fe-6919-18da-f3ee-fc10fac8c1bf.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## HyundaiData 构造函数 

[原文連結](http://api.hslcommunication.cn/html/0b9a7201-f20f-38a1-e8cb-cb2d84214806.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.Hyundai](../html/caddf3fe-6919-18da-f3ee-fc10fac8c1bf.htm "HslCommunication.Robot.Hyundai")

[HyundaiData 类](../html/86650e61-38f2-c3ee-2fcd-c45d8d120f9f.htm "HyundaiData 类")

[HyundaiData 构造函数](../html/0b9a7201-f20f-38a1-e8cb-cb2d84214806.htm "HyundaiData 构造函数 ")

[HyundaiData 构造函数](../html/fea8a726-d649-ca0f-cd42-dcbd1f6fe18f.htm "HyundaiData 构造函数 ")

[HyundaiData 构造函数 (Byte[])](../html/25014cc4-a3dc-7712-2f78-d70e3a0b9db6.htm "HyundaiData 构造函数 (Byte[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HyundaiData 构造函数 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [HyundaiData](fea8a726-d649-ca0f-cd42-dcbd1f6fe18f.htm) | 实例化一个默认的对象 |
| 公共方法 | [HyundaiData(Byte)](25014cc4-a3dc-7712-2f78-d70e3a0b9db6.htm) | 通过缓存对象实例化一个 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[HyundaiData 类](86650e61-38f2-c3ee-2fcd-c45d8d120f9f.htm)

[HslCommunication.Robot.Hyundai 命名空间](caddf3fe-6919-18da-f3ee-fc10fac8c1bf.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## HyundaiData 构造函数 

[原文連結](http://api.hslcommunication.cn/html/fea8a726-d649-ca0f-cd42-dcbd1f6fe18f.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.Hyundai](../html/caddf3fe-6919-18da-f3ee-fc10fac8c1bf.htm "HslCommunication.Robot.Hyundai")

[HyundaiData 类](../html/86650e61-38f2-c3ee-2fcd-c45d8d120f9f.htm "HyundaiData 类")

[HyundaiData 构造函数](../html/0b9a7201-f20f-38a1-e8cb-cb2d84214806.htm "HyundaiData 构造函数 ")

[HyundaiData 构造函数](../html/fea8a726-d649-ca0f-cd42-dcbd1f6fe18f.htm "HyundaiData 构造函数 ")

[HyundaiData 构造函数 (Byte[])](../html/25014cc4-a3dc-7712-2f78-d70e3a0b9db6.htm "HyundaiData 构造函数 (Byte[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HyundaiData 构造函数 |

实例化一个默认的对象

**命名空间：**
 [HslCommunication.Robot.Hyundai](caddf3fe-6919-18da-f3ee-fc10fac8c1bf.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public HyundaiData()
```

```
Public Sub New
```

```
public:
HyundaiData()
```

```
new : unit -> HyundaiData
```

![](../icons/SectionExpanded.png)参见

#### 引用

[HyundaiData 类](86650e61-38f2-c3ee-2fcd-c45d8d120f9f.htm)

[HyundaiData 重载](0b9a7201-f20f-38a1-e8cb-cb2d84214806.htm)

[HslCommunication.Robot.Hyundai 命名空间](caddf3fe-6919-18da-f3ee-fc10fac8c1bf.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## HyundaiData 构造函数 (Byte[])

[原文連結](http://api.hslcommunication.cn/html/25014cc4-a3dc-7712-2f78-d70e3a0b9db6.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.Hyundai](../html/caddf3fe-6919-18da-f3ee-fc10fac8c1bf.htm "HslCommunication.Robot.Hyundai")

[HyundaiData 类](../html/86650e61-38f2-c3ee-2fcd-c45d8d120f9f.htm "HyundaiData 类")

[HyundaiData 构造函数](../html/0b9a7201-f20f-38a1-e8cb-cb2d84214806.htm "HyundaiData 构造函数 ")

[HyundaiData 构造函数](../html/fea8a726-d649-ca0f-cd42-dcbd1f6fe18f.htm "HyundaiData 构造函数 ")

[HyundaiData 构造函数 (Byte[])](../html/25014cc4-a3dc-7712-2f78-d70e3a0b9db6.htm "HyundaiData 构造函数 (Byte[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HyundaiData 构造函数 (Byte) |

通过缓存对象实例化一个

**命名空间：**
 [HslCommunication.Robot.Hyundai](caddf3fe-6919-18da-f3ee-fc10fac8c1bf.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public HyundaiData(
	byte[] buffer
)
```

```
Public Sub New ( 
	buffer As Byte()
)
```

```
public:
HyundaiData(
	array<unsigned char>^ buffer
)
```

```
new : 
        buffer : byte[] -> HyundaiData
```

#### 参数

buffer
:   类型：SystemByte  

    [缺少 "M:HslCommunication.Robot.Hyundai.HyundaiData.#ctor(System.Byte[])" 的 <param name="buffer"/> 文档]

![](../icons/SectionExpanded.png)参见

#### 引用

[HyundaiData 类](86650e61-38f2-c3ee-2fcd-c45d8d120f9f.htm)

[HyundaiData 重载](0b9a7201-f20f-38a1-e8cb-cb2d84214806.htm)

[HslCommunication.Robot.Hyundai 命名空间](caddf3fe-6919-18da-f3ee-fc10fac8c1bf.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## HyundaiData 属性

[原文連結](http://api.hslcommunication.cn/html/50ed49dd-87ac-204e-6164-898e5ef2fb31.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.Hyundai](../html/caddf3fe-6919-18da-f3ee-fc10fac8c1bf.htm "HslCommunication.Robot.Hyundai")

[HyundaiData 类](../html/86650e61-38f2-c3ee-2fcd-c45d8d120f9f.htm "HyundaiData 类")

[HyundaiData 属性](../html/50ed49dd-87ac-204e-6164-898e5ef2fb31.htm "HyundaiData 属性")

[CharDummy 属性](../html/e5796267-5048-c787-dae9-f0a6ee45c0fa.htm "CharDummy 属性 ")

[Command 属性](../html/aa7fd03a-39c0-44e7-1282-7bf86bc1793b.htm "Command 属性 ")

[Count 属性](../html/3a60c870-22e5-444a-70a0-60e32e2e8ae6.htm "Count 属性 ")

[Data 属性](../html/aa511f8a-78ce-7d57-ec81-e0ca1425d3d5.htm "Data 属性 ")

[IntDummy 属性](../html/b42a4d50-2f24-6098-5def-89a801d38e6e.htm "IntDummy 属性 ")

[State 属性](../html/657c8050-8cce-cc62-40ce-2c7ef9df1bef.htm "State 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HyundaiData 属性 |

[HyundaiData](86650e61-38f2-c3ee-2fcd-c45d8d120f9f.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [CharDummy](e5796267-5048-c787-dae9-f0a6ee45c0fa.htm) | 虚标记 |
| 公共属性 | [Command](aa7fd03a-39c0-44e7-1282-7bf86bc1793b.htm) | 命令码，从控制器发数据到PC和PC到控制器，两者的命令不一样 |
| 公共属性 | [Count](3a60c870-22e5-444a-70a0-60e32e2e8ae6.htm) | 标记数据，从PLC发送给机器人的数据，原封不动的返回 |
| 公共属性 | [Data](aa511f8a-78ce-7d57-ec81-e0ca1425d3d5.htm) | 关节坐标数据，包含X,Y,Z,W,P,R，三个位置数据，三个角度数据。 |
| 公共属性 | [IntDummy](b42a4d50-2f24-6098-5def-89a801d38e6e.htm) | 虚标记 |
| 公共属性 | [State](657c8050-8cce-cc62-40ce-2c7ef9df1bef.htm) | 状态码 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[HyundaiData 类](86650e61-38f2-c3ee-2fcd-c45d8d120f9f.htm)

[HslCommunication.Robot.Hyundai 命名空间](caddf3fe-6919-18da-f3ee-fc10fac8c1bf.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## CharDummy 属性 

[原文連結](http://api.hslcommunication.cn/html/e5796267-5048-c787-dae9-f0a6ee45c0fa.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.Hyundai](../html/caddf3fe-6919-18da-f3ee-fc10fac8c1bf.htm "HslCommunication.Robot.Hyundai")

[HyundaiData 类](../html/86650e61-38f2-c3ee-2fcd-c45d8d120f9f.htm "HyundaiData 类")

[HyundaiData 属性](../html/50ed49dd-87ac-204e-6164-898e5ef2fb31.htm "HyundaiData 属性")

[CharDummy 属性](../html/e5796267-5048-c787-dae9-f0a6ee45c0fa.htm "CharDummy 属性 ")

[Command 属性](../html/aa7fd03a-39c0-44e7-1282-7bf86bc1793b.htm "Command 属性 ")

[Count 属性](../html/3a60c870-22e5-444a-70a0-60e32e2e8ae6.htm "Count 属性 ")

[Data 属性](../html/aa511f8a-78ce-7d57-ec81-e0ca1425d3d5.htm "Data 属性 ")

[IntDummy 属性](../html/b42a4d50-2f24-6098-5def-89a801d38e6e.htm "IntDummy 属性 ")

[State 属性](../html/657c8050-8cce-cc62-40ce-2c7ef9df1bef.htm "State 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HyundaiDataCharDummy 属性 |

虚标记

**命名空间：**
 [HslCommunication.Robot.Hyundai](caddf3fe-6919-18da-f3ee-fc10fac8c1bf.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public string CharDummy { get; set; }
```

```
Public Property CharDummy As String
	Get
	Set
```

```
public:
property String^ CharDummy {
	String^ get ();
	void set (String^ value);
}
```

```
member CharDummy : string with get, set
```

#### 属性值

类型：String

![](../icons/SectionExpanded.png)参见

#### 引用

[HyundaiData 类](86650e61-38f2-c3ee-2fcd-c45d8d120f9f.htm)

[HslCommunication.Robot.Hyundai 命名空间](caddf3fe-6919-18da-f3ee-fc10fac8c1bf.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Command 属性 

[原文連結](http://api.hslcommunication.cn/html/aa7fd03a-39c0-44e7-1282-7bf86bc1793b.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.Hyundai](../html/caddf3fe-6919-18da-f3ee-fc10fac8c1bf.htm "HslCommunication.Robot.Hyundai")

[HyundaiData 类](../html/86650e61-38f2-c3ee-2fcd-c45d8d120f9f.htm "HyundaiData 类")

[HyundaiData 属性](../html/50ed49dd-87ac-204e-6164-898e5ef2fb31.htm "HyundaiData 属性")

[CharDummy 属性](../html/e5796267-5048-c787-dae9-f0a6ee45c0fa.htm "CharDummy 属性 ")

[Command 属性](../html/aa7fd03a-39c0-44e7-1282-7bf86bc1793b.htm "Command 属性 ")

[Count 属性](../html/3a60c870-22e5-444a-70a0-60e32e2e8ae6.htm "Count 属性 ")

[Data 属性](../html/aa511f8a-78ce-7d57-ec81-e0ca1425d3d5.htm "Data 属性 ")

[IntDummy 属性](../html/b42a4d50-2f24-6098-5def-89a801d38e6e.htm "IntDummy 属性 ")

[State 属性](../html/657c8050-8cce-cc62-40ce-2c7ef9df1bef.htm "State 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HyundaiDataCommand 属性 |

命令码，从控制器发数据到PC和PC到控制器，两者的命令不一样

**命名空间：**
 [HslCommunication.Robot.Hyundai](caddf3fe-6919-18da-f3ee-fc10fac8c1bf.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public char Command { get; set; }
```

```
Public Property Command As Char
	Get
	Set
```

```
public:
property wchar_t Command {
	wchar_t get ();
	void set (wchar_t value);
}
```

```
member Command : char with get, set
```

#### 属性值

类型：Char

![](../icons/SectionExpanded.png)参见

#### 引用

[HyundaiData 类](86650e61-38f2-c3ee-2fcd-c45d8d120f9f.htm)

[HslCommunication.Robot.Hyundai 命名空间](caddf3fe-6919-18da-f3ee-fc10fac8c1bf.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Count 属性 

[原文連結](http://api.hslcommunication.cn/html/3a60c870-22e5-444a-70a0-60e32e2e8ae6.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.Hyundai](../html/caddf3fe-6919-18da-f3ee-fc10fac8c1bf.htm "HslCommunication.Robot.Hyundai")

[HyundaiData 类](../html/86650e61-38f2-c3ee-2fcd-c45d8d120f9f.htm "HyundaiData 类")

[HyundaiData 属性](../html/50ed49dd-87ac-204e-6164-898e5ef2fb31.htm "HyundaiData 属性")

[CharDummy 属性](../html/e5796267-5048-c787-dae9-f0a6ee45c0fa.htm "CharDummy 属性 ")

[Command 属性](../html/aa7fd03a-39c0-44e7-1282-7bf86bc1793b.htm "Command 属性 ")

[Count 属性](../html/3a60c870-22e5-444a-70a0-60e32e2e8ae6.htm "Count 属性 ")

[Data 属性](../html/aa511f8a-78ce-7d57-ec81-e0ca1425d3d5.htm "Data 属性 ")

[IntDummy 属性](../html/b42a4d50-2f24-6098-5def-89a801d38e6e.htm "IntDummy 属性 ")

[State 属性](../html/657c8050-8cce-cc62-40ce-2c7ef9df1bef.htm "State 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HyundaiDataCount 属性 |

标记数据，从PLC发送给机器人的数据，原封不动的返回

**命名空间：**
 [HslCommunication.Robot.Hyundai](caddf3fe-6919-18da-f3ee-fc10fac8c1bf.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public int Count { get; set; }
```

```
Public Property Count As Integer
	Get
	Set
```

```
public:
property int Count {
	int get ();
	void set (int value);
}
```

```
member Count : int with get, set
```

#### 属性值

类型：Int32

![](../icons/SectionExpanded.png)参见

#### 引用

[HyundaiData 类](86650e61-38f2-c3ee-2fcd-c45d8d120f9f.htm)

[HslCommunication.Robot.Hyundai 命名空间](caddf3fe-6919-18da-f3ee-fc10fac8c1bf.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Data 属性 

[原文連結](http://api.hslcommunication.cn/html/aa511f8a-78ce-7d57-ec81-e0ca1425d3d5.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.Hyundai](../html/caddf3fe-6919-18da-f3ee-fc10fac8c1bf.htm "HslCommunication.Robot.Hyundai")

[HyundaiData 类](../html/86650e61-38f2-c3ee-2fcd-c45d8d120f9f.htm "HyundaiData 类")

[HyundaiData 属性](../html/50ed49dd-87ac-204e-6164-898e5ef2fb31.htm "HyundaiData 属性")

[CharDummy 属性](../html/e5796267-5048-c787-dae9-f0a6ee45c0fa.htm "CharDummy 属性 ")

[Command 属性](../html/aa7fd03a-39c0-44e7-1282-7bf86bc1793b.htm "Command 属性 ")

[Count 属性](../html/3a60c870-22e5-444a-70a0-60e32e2e8ae6.htm "Count 属性 ")

[Data 属性](../html/aa511f8a-78ce-7d57-ec81-e0ca1425d3d5.htm "Data 属性 ")

[IntDummy 属性](../html/b42a4d50-2f24-6098-5def-89a801d38e6e.htm "IntDummy 属性 ")

[State 属性](../html/657c8050-8cce-cc62-40ce-2c7ef9df1bef.htm "State 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HyundaiDataData 属性 |

关节坐标数据，包含X,Y,Z,W,P,R，三个位置数据，三个角度数据。

**命名空间：**
 [HslCommunication.Robot.Hyundai](caddf3fe-6919-18da-f3ee-fc10fac8c1bf.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public double[] Data { get; set; }
```

```
Public Property Data As Double()
	Get
	Set
```

```
public:
property array<double>^ Data {
	array<double>^ get ();
	void set (array<double>^ value);
}
```

```
member Data : float[] with get, set
```

#### 属性值

类型：Double

![](../icons/SectionExpanded.png)参见

#### 引用

[HyundaiData 类](86650e61-38f2-c3ee-2fcd-c45d8d120f9f.htm)

[HslCommunication.Robot.Hyundai 命名空间](caddf3fe-6919-18da-f3ee-fc10fac8c1bf.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## IntDummy 属性 

[原文連結](http://api.hslcommunication.cn/html/b42a4d50-2f24-6098-5def-89a801d38e6e.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.Hyundai](../html/caddf3fe-6919-18da-f3ee-fc10fac8c1bf.htm "HslCommunication.Robot.Hyundai")

[HyundaiData 类](../html/86650e61-38f2-c3ee-2fcd-c45d8d120f9f.htm "HyundaiData 类")

[HyundaiData 属性](../html/50ed49dd-87ac-204e-6164-898e5ef2fb31.htm "HyundaiData 属性")

[CharDummy 属性](../html/e5796267-5048-c787-dae9-f0a6ee45c0fa.htm "CharDummy 属性 ")

[Command 属性](../html/aa7fd03a-39c0-44e7-1282-7bf86bc1793b.htm "Command 属性 ")

[Count 属性](../html/3a60c870-22e5-444a-70a0-60e32e2e8ae6.htm "Count 属性 ")

[Data 属性](../html/aa511f8a-78ce-7d57-ec81-e0ca1425d3d5.htm "Data 属性 ")

[IntDummy 属性](../html/b42a4d50-2f24-6098-5def-89a801d38e6e.htm "IntDummy 属性 ")

[State 属性](../html/657c8050-8cce-cc62-40ce-2c7ef9df1bef.htm "State 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HyundaiDataIntDummy 属性 |

虚标记

**命名空间：**
 [HslCommunication.Robot.Hyundai](caddf3fe-6919-18da-f3ee-fc10fac8c1bf.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public int IntDummy { get; set; }
```

```
Public Property IntDummy As Integer
	Get
	Set
```

```
public:
property int IntDummy {
	int get ();
	void set (int value);
}
```

```
member IntDummy : int with get, set
```

#### 属性值

类型：Int32

![](../icons/SectionExpanded.png)参见

#### 引用

[HyundaiData 类](86650e61-38f2-c3ee-2fcd-c45d8d120f9f.htm)

[HslCommunication.Robot.Hyundai 命名空间](caddf3fe-6919-18da-f3ee-fc10fac8c1bf.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## State 属性 

[原文連結](http://api.hslcommunication.cn/html/657c8050-8cce-cc62-40ce-2c7ef9df1bef.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.Hyundai](../html/caddf3fe-6919-18da-f3ee-fc10fac8c1bf.htm "HslCommunication.Robot.Hyundai")

[HyundaiData 类](../html/86650e61-38f2-c3ee-2fcd-c45d8d120f9f.htm "HyundaiData 类")

[HyundaiData 属性](../html/50ed49dd-87ac-204e-6164-898e5ef2fb31.htm "HyundaiData 属性")

[CharDummy 属性](../html/e5796267-5048-c787-dae9-f0a6ee45c0fa.htm "CharDummy 属性 ")

[Command 属性](../html/aa7fd03a-39c0-44e7-1282-7bf86bc1793b.htm "Command 属性 ")

[Count 属性](../html/3a60c870-22e5-444a-70a0-60e32e2e8ae6.htm "Count 属性 ")

[Data 属性](../html/aa511f8a-78ce-7d57-ec81-e0ca1425d3d5.htm "Data 属性 ")

[IntDummy 属性](../html/b42a4d50-2f24-6098-5def-89a801d38e6e.htm "IntDummy 属性 ")

[State 属性](../html/657c8050-8cce-cc62-40ce-2c7ef9df1bef.htm "State 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HyundaiDataState 属性 |

状态码

**命名空间：**
 [HslCommunication.Robot.Hyundai](caddf3fe-6919-18da-f3ee-fc10fac8c1bf.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public int State { get; set; }
```

```
Public Property State As Integer
	Get
	Set
```

```
public:
property int State {
	int get ();
	void set (int value);
}
```

```
member State : int with get, set
```

#### 属性值

类型：Int32

![](../icons/SectionExpanded.png)参见

#### 引用

[HyundaiData 类](86650e61-38f2-c3ee-2fcd-c45d8d120f9f.htm)

[HslCommunication.Robot.Hyundai 命名空间](caddf3fe-6919-18da-f3ee-fc10fac8c1bf.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## HyundaiData 方法

[原文連結](http://api.hslcommunication.cn/html/c8612406-272d-33dc-b404-b0008b3b2767.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.Hyundai](../html/caddf3fe-6919-18da-f3ee-fc10fac8c1bf.htm "HslCommunication.Robot.Hyundai")

[HyundaiData 类](../html/86650e61-38f2-c3ee-2fcd-c45d8d120f9f.htm "HyundaiData 类")

[HyundaiData 方法](../html/c8612406-272d-33dc-b404-b0008b3b2767.htm "HyundaiData 方法")

[LoadBy 方法](../html/34da99e5-6c2c-0365-410c-34846c0bc90e.htm "LoadBy 方法 ")

[ToBytes 方法](../html/3cc5be93-38ea-4f3b-a338-620464accd37.htm "ToBytes 方法 ")

[ToString 方法](../html/8905a222-5ea4-ed24-547e-33f2ca530baa.htm "ToString 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HyundaiData 方法 |

[HyundaiData](86650e61-38f2-c3ee-2fcd-c45d8d120f9f.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法 | GetType | (继承自 Object。) |
| 公共方法 | [LoadBy](34da99e5-6c2c-0365-410c-34846c0bc90e.htm) | 从字节数组的指定索引开始加载现在机器人的数据 |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法 | [ToBytes](3cc5be93-38ea-4f3b-a338-620464accd37.htm) | 将现代机器人的数据转换为字节数组 |
| 公共方法 | [ToString](8905a222-5ea4-ed24-547e-33f2ca530baa.htm) | (重写 ObjectToString.) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[HyundaiData 类](86650e61-38f2-c3ee-2fcd-c45d8d120f9f.htm)

[HslCommunication.Robot.Hyundai 命名空间](caddf3fe-6919-18da-f3ee-fc10fac8c1bf.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## LoadBy 方法 

[原文連結](http://api.hslcommunication.cn/html/34da99e5-6c2c-0365-410c-34846c0bc90e.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.Hyundai](../html/caddf3fe-6919-18da-f3ee-fc10fac8c1bf.htm "HslCommunication.Robot.Hyundai")

[HyundaiData 类](../html/86650e61-38f2-c3ee-2fcd-c45d8d120f9f.htm "HyundaiData 类")

[HyundaiData 方法](../html/c8612406-272d-33dc-b404-b0008b3b2767.htm "HyundaiData 方法")

[LoadBy 方法](../html/34da99e5-6c2c-0365-410c-34846c0bc90e.htm "LoadBy 方法 ")

[ToBytes 方法](../html/3cc5be93-38ea-4f3b-a338-620464accd37.htm "ToBytes 方法 ")

[ToString 方法](../html/8905a222-5ea4-ed24-547e-33f2ca530baa.htm "ToString 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HyundaiDataLoadBy 方法 |

从字节数组的指定索引开始加载现在机器人的数据

**命名空间：**
 [HslCommunication.Robot.Hyundai](caddf3fe-6919-18da-f3ee-fc10fac8c1bf.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public void LoadBy(
	byte[] buffer,
	int index = 0
)
```

```
Public Sub LoadBy ( 
	buffer As Byte(),
	Optional index As Integer = 0
)
```

```
public:
void LoadBy(
	array<unsigned char>^ buffer, 
	int index = 0
)
```

```
member LoadBy : 
        buffer : byte[] * 
        ?index : int 
(* Defaults:
        let _index = defaultArg index 0
*)
-> unit 
```

#### 参数

buffer
:   类型：SystemByte  
    原始的字节数据

index (Optional)
:   类型：SystemInt32  
    起始的索引信息

![](../icons/SectionExpanded.png)参见

#### 引用

[HyundaiData 类](86650e61-38f2-c3ee-2fcd-c45d8d120f9f.htm)

[HslCommunication.Robot.Hyundai 命名空间](caddf3fe-6919-18da-f3ee-fc10fac8c1bf.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ToBytes 方法 

[原文連結](http://api.hslcommunication.cn/html/3cc5be93-38ea-4f3b-a338-620464accd37.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.Hyundai](../html/caddf3fe-6919-18da-f3ee-fc10fac8c1bf.htm "HslCommunication.Robot.Hyundai")

[HyundaiData 类](../html/86650e61-38f2-c3ee-2fcd-c45d8d120f9f.htm "HyundaiData 类")

[HyundaiData 方法](../html/c8612406-272d-33dc-b404-b0008b3b2767.htm "HyundaiData 方法")

[LoadBy 方法](../html/34da99e5-6c2c-0365-410c-34846c0bc90e.htm "LoadBy 方法 ")

[ToBytes 方法](../html/3cc5be93-38ea-4f3b-a338-620464accd37.htm "ToBytes 方法 ")

[ToString 方法](../html/8905a222-5ea4-ed24-547e-33f2ca530baa.htm "ToString 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HyundaiDataToBytes 方法 |

将现代机器人的数据转换为字节数组

**命名空间：**
 [HslCommunication.Robot.Hyundai](caddf3fe-6919-18da-f3ee-fc10fac8c1bf.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public byte[] ToBytes()
```

```
Public Function ToBytes As Byte()
```

```
public:
array<unsigned char>^ ToBytes()
```

```
member ToBytes : unit -> byte[] 
```

#### 返回值

类型：Byte  
字节数组

![](../icons/SectionExpanded.png)参见

#### 引用

[HyundaiData 类](86650e61-38f2-c3ee-2fcd-c45d8d120f9f.htm)

[HslCommunication.Robot.Hyundai 命名空间](caddf3fe-6919-18da-f3ee-fc10fac8c1bf.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ToString 方法 

[原文連結](http://api.hslcommunication.cn/html/8905a222-5ea4-ed24-547e-33f2ca530baa.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.Hyundai](../html/caddf3fe-6919-18da-f3ee-fc10fac8c1bf.htm "HslCommunication.Robot.Hyundai")

[HyundaiData 类](../html/86650e61-38f2-c3ee-2fcd-c45d8d120f9f.htm "HyundaiData 类")

[HyundaiData 方法](../html/c8612406-272d-33dc-b404-b0008b3b2767.htm "HyundaiData 方法")

[LoadBy 方法](../html/34da99e5-6c2c-0365-410c-34846c0bc90e.htm "LoadBy 方法 ")

[ToBytes 方法](../html/3cc5be93-38ea-4f3b-a338-620464accd37.htm "ToBytes 方法 ")

[ToString 方法](../html/8905a222-5ea4-ed24-547e-33f2ca530baa.htm "ToString 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HyundaiDataToString 方法 |

[缺少 "M:HslCommunication.Robot.Hyundai.HyundaiData.ToString" 的 <summary> 文档]

**命名空间：**
 [HslCommunication.Robot.Hyundai](caddf3fe-6919-18da-f3ee-fc10fac8c1bf.htm)  
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

[缺少 "M:HslCommunication.Robot.Hyundai.HyundaiData.ToString" 的 <returns> 文档]

![](../icons/SectionExpanded.png)参见

#### 引用

[HyundaiData 类](86650e61-38f2-c3ee-2fcd-c45d8d120f9f.htm)

[HslCommunication.Robot.Hyundai 命名空间](caddf3fe-6919-18da-f3ee-fc10fac8c1bf.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## HyundaiUdpNet 构造函数 

[原文連結](http://api.hslcommunication.cn/html/a760615f-a05d-ac04-364f-c35bd88dae6f.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.Hyundai](../html/caddf3fe-6919-18da-f3ee-fc10fac8c1bf.htm "HslCommunication.Robot.Hyundai")

[HyundaiUdpNet 类](../html/63e7fe67-d9ec-efd5-a7d8-a8cde49538e2.htm "HyundaiUdpNet 类")

[HyundaiUdpNet 构造函数](../html/a760615f-a05d-ac04-364f-c35bd88dae6f.htm "HyundaiUdpNet 构造函数 ")

[HyundaiUdpNet 属性](../html/ce20a3c5-6c99-f1de-8e83-58c03ac5f17a.htm "HyundaiUdpNet 属性")

[HyundaiUdpNet 方法](../html/8d5c4acb-4433-86d1-b295-511cb669b68f.htm "HyundaiUdpNet 方法")

[HyundaiUdpNet 事件](../html/18aacb19-c509-c0bc-371c-8606a5dfab8e.htm "HyundaiUdpNet 事件")

[HyundaiUdpNet 字段](../html/05a0d747-ce32-be60-bb8c-6a2f81043aef.htm "HyundaiUdpNet 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HyundaiUdpNet 构造函数 |

实例化一个默认的对象  
Instantiate a default object

**命名空间：**
 [HslCommunication.Robot.Hyundai](caddf3fe-6919-18da-f3ee-fc10fac8c1bf.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public HyundaiUdpNet()
```

```
Public Sub New
```

```
public:
HyundaiUdpNet()
```

```
new : unit -> HyundaiUdpNet
```

![](../icons/SectionExpanded.png)参见

#### 引用

[HyundaiUdpNet 类](63e7fe67-d9ec-efd5-a7d8-a8cde49538e2.htm)

[HslCommunication.Robot.Hyundai 命名空间](caddf3fe-6919-18da-f3ee-fc10fac8c1bf.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## HyundaiUdpNet 属性

[原文連結](http://api.hslcommunication.cn/html/ce20a3c5-6c99-f1de-8e83-58c03ac5f17a.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.Hyundai](../html/caddf3fe-6919-18da-f3ee-fc10fac8c1bf.htm "HslCommunication.Robot.Hyundai")

[HyundaiUdpNet 类](../html/63e7fe67-d9ec-efd5-a7d8-a8cde49538e2.htm "HyundaiUdpNet 类")

[HyundaiUdpNet 构造函数](../html/a760615f-a05d-ac04-364f-c35bd88dae6f.htm "HyundaiUdpNet 构造函数 ")

[HyundaiUdpNet 属性](../html/ce20a3c5-6c99-f1de-8e83-58c03ac5f17a.htm "HyundaiUdpNet 属性")

[HyundaiUdpNet 方法](../html/8d5c4acb-4433-86d1-b295-511cb669b68f.htm "HyundaiUdpNet 方法")

[HyundaiUdpNet 事件](../html/18aacb19-c509-c0bc-371c-8606a5dfab8e.htm "HyundaiUdpNet 事件")

[HyundaiUdpNet 字段](../html/05a0d747-ce32-be60-bb8c-6a2f81043aef.htm "HyundaiUdpNet 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HyundaiUdpNet 属性 |

[HyundaiUdpNet](63e7fe67-d9ec-efd5-a7d8-a8cde49538e2.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [CheckSerialDataComplete](bef5c4d1-8fac-1f94-42ab-155ad17f76d5.htm) | 检查串口接收到的数据是否完整 (继承自 [CommunicationServer](5a3cb700-c0f6-e35d-6807-dfe8bfecf747.htm)。) |
| 公共属性 | [CreateNewMessage](19bf00bd-dfa5-8afb-6716-a985f22da99a.htm) | 获取一个新的消息对象的方法，需要在继承类里面进行重写  The method to get a new message object needs to be overridden in the inheritance class (继承自 [CommunicationServer](5a3cb700-c0f6-e35d-6807-dfe8bfecf747.htm)。) |
| 公共属性 | [CreatePipeSession](5bebf8f4-912a-0e90-f877-7ae0fe9f02ed.htm) | 创建会话状态的委托对象，也就可以自己指定创建自定义的会话 (继承自 [CommunicationServer](5a3cb700-c0f6-e35d-6807-dfe8bfecf747.htm)。) |
| 公共属性 | [EnableIPv6](f15cbcb3-5dcf-a658-0a93-3f0921383c32.htm) | 获取或设置服务器是否支持IPv6的地址协议信息  Get or set whether the server supports IPv6 address protocol information (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |
| 公共属性 | [ForceSerialReceiveOnce](4e1ee30d-f0c6-72fd-161a-83cc6c2e8e5e.htm) | 获取或设置当前的服务器接收串口数据时候，是否强制只接收一次数据，默认为false，适合点对点通信，如果你总线形式的连接，则需要设置 True  Get or set whether to force the data to be received only once when the current server receives serial port data. The default value is false, which is suitable for point-to-point communication. If you have a bus connection, you need to set True (继承自 [CommunicationServer](5a3cb700-c0f6-e35d-6807-dfe8bfecf747.htm)。) |
| 公共属性 | [IsStarted](8d6786a8-d416-4f45-158c-de6be0634aca.htm) | 服务器引擎是否启动  Whether the server engine is started (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |
| 公共属性 | [IsUseSSL](f4daa1e8-f433-0ee2-8284-16c4754eba25.htm) | 获取当前的服务器是否使用了SSL证书功能 (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |
| 公共属性 | [LocalAddress](c0effb49-cfe6-3267-9aef-afdab27a41c3.htm) | 获取或设置服务器绑定的本地IP地址，默认为空，使用本地所有可用的ip地址  Obtain or set the local IP address bound to the server, which is empty by default, and use all available IP addresses locally (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |
| 公共属性 | [LogDebugMessage](8a6db355-1130-2af8-cc6a-aeccc1ac0b29.htm) | 记录一些调试日志的委托，将会进行输出调试文本。  The delegate that records some debug logs will output debug text. (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |
| 公共属性代码示例 | [LogNet](746b148d-1ee3-1ee7-ee5d-3938eb72b2df.htm) | 组件的日志工具，支持日志记录，只要实例化后，当前网络的基本信息，就以[DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm)等级进行输出  The component's logging tool supports logging. As long as the instantiation of the basic network information, the output will be output at [DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm) (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |
| 公共属性 | [Port](6f700f7a-bd7e-c4b3-1275-a30491469e8f.htm) | 获取或设置服务器的端口号，如果是设置，需要在服务器启动前设置完成，才能生效。  Gets or sets the port number of the server. If it is set, it needs to be set before the server starts to take effect. (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |
| 公共属性 | [SerialReceiveAtleastTime](5802707e-1dd6-f7f9-e293-dd387ae817fd.htm) | 获取或设置串口模式下，接收一条数据最短的时间要求，当设备发送的数据非常慢的时候，或是分割发送数据的时候，就需要将本值设置的大一点，默认为20ms  Get or set the shortest time required to receive a piece of data in serial port mode. When the data sent by the device is very slow, or when the data is divided and sent, you need to set this value to a larger value, the default is 20ms (继承自 [CommunicationServer](5a3cb700-c0f6-e35d-6807-dfe8bfecf747.htm)。) |
| 公共属性 | [SessionsMax](c65bac2d-a238-f93d-02f2-a7ed1297cf7a.htm) | 获取或设置当前允许登录的最大客户端数量，默认为 uint.MaxValue = 4294967295  Gets or sets the maximum number of clients that are currently allowed to log in, which defaults to uint.MaxValue = 4294967295 (继承自 [CommunicationServer](5a3cb700-c0f6-e35d-6807-dfe8bfecf747.htm)。) |
| 公共属性 | [SocketKeepAliveTime](7d77e3df-9bbe-1a6e-a051-c87d9516c3a2.htm) | 获取或设置客户端的Socket的心跳时间信息，这个是Socket底层自动实现的心跳包，不基于协议层实现。默认小于0，不开启心跳检测，如果需要开启，设置 60\_000 比较合适，单位毫秒  Get or set the heartbeat time information of the Socket of the client. This is the heartbeat packet automatically implemented by the bottom layer of the Socket, not based on the protocol layer. The default value is less than 0, and heartbeat detection is not enabled. If you need to enable it, it is more appropriate to set 60\_000, in milliseconds. (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |
| 公共属性 | [ThreadPoolLoginAfterClientCheck](e926b0f9-3849-07d2-ed1d-0a249fad29ae.htm) | 当线程检查后，进行登录之前的检查，通常用于自定义的握手包校验操作。仅对TCP通信的时候有效。 (继承自 [CommunicationServer](5a3cb700-c0f6-e35d-6807-dfe8bfecf747.htm)。) |
| 公共属性 | [UdpBufferSize](762072d5-f18f-8422-fc5f-19fb0eae7eff.htm) | 获取或设置一次接收时的数据长度，默认2KB数据长度，特殊情况的时候需要调整  Gets or sets the length of data received at a time. The default length is 2KB (继承自 [CommunicationServer](5a3cb700-c0f6-e35d-6807-dfe8bfecf747.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[HyundaiUdpNet 类](63e7fe67-d9ec-efd5-a7d8-a8cde49538e2.htm)

[HslCommunication.Robot.Hyundai 命名空间](caddf3fe-6919-18da-f3ee-fc10fac8c1bf.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## HyundaiUdpNet 方法

[原文連結](http://api.hslcommunication.cn/html/8d5c4acb-4433-86d1-b295-511cb669b68f.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.Hyundai](../html/caddf3fe-6919-18da-f3ee-fc10fac8c1bf.htm "HslCommunication.Robot.Hyundai")

[HyundaiUdpNet 类](../html/63e7fe67-d9ec-efd5-a7d8-a8cde49538e2.htm "HyundaiUdpNet 类")

[HyundaiUdpNet 方法](../html/8d5c4acb-4433-86d1-b295-511cb669b68f.htm "HyundaiUdpNet 方法")

[MoveX 方法](../html/e8a255d1-b98d-8aa7-452c-e6e57fb3514b.htm "MoveX 方法 ")

[MoveY 方法](../html/7bd11b89-d4d2-b19f-88da-dc3218e83a83.htm "MoveY 方法 ")

[MoveZ 方法](../html/41b575e3-242b-7144-4289-0fc59a1e5c2d.htm "MoveZ 方法 ")

[RotateX 方法](../html/f154ee9b-d5d5-4ed7-a52e-93885323fb8e.htm "RotateX 方法 ")

[RotateY 方法](../html/2dbb6aa6-1040-6aff-d5cd-eeb5932e1ad5.htm "RotateY 方法 ")

[RotateZ 方法](../html/f96cafec-8f12-63d6-0903-4e20a6282476.htm "RotateZ 方法 ")

[ToString 方法](../html/8010794a-c51c-3618-8667-024a32e2cadd.htm "ToString 方法 ")

[Write 方法](../html/061ebd46-cffc-5e79-1de2-0044dca27b70.htm "Write 方法 ")

[WriteIncrementPos 方法](../html/50870834-55a9-e75b-e524-ffb8a32046f3.htm "WriteIncrementPos 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HyundaiUdpNet 方法 |

[HyundaiUdpNet](63e7fe67-d9ec-efd5-a7d8-a8cde49538e2.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [AddSession](22a49542-4dd7-d360-ed5d-85266665b102.htm) | 新增加一个管道会话信息  A new pipeline session information has been added (继承自 [CommunicationServer](5a3cb700-c0f6-e35d-6807-dfe8bfecf747.htm)。) |
| 受保护的方法 | [CheckSerialReceiveDataComplete](a66e7f2a-c605-d949-b1ec-9271b75df1b6.htm) | 检查串口接收的数据是否完成的方法，如果接收完成，则返回True (继承自 [CommunicationServer](5a3cb700-c0f6-e35d-6807-dfe8bfecf747.htm)。) |
| 公共方法 | [CloseSerialSlave](f5a8b6f0-0481-b9b3-c31b-90d9a945b416.htm) | 关闭提供从机服务的串口对象  Close the serial port object that provides slave services (继承自 [CommunicationServer](5a3cb700-c0f6-e35d-6807-dfe8bfecf747.htm)。) |
| 公共方法 | [ConnectHslAlientClient](b98a142a-9a96-4cc8-4fdc-704dbfc7b438.htm) | 创建一个指定的异形客户端连接，使用Hsl协议来发送注册包  Create a specified profiled client connection and use the Hsl protocol to send registration packets (继承自 [CommunicationServer](5a3cb700-c0f6-e35d-6807-dfe8bfecf747.htm)。) |
| 公共方法 | [ConnectRemoteServer](5946c84a-5591-14d2-2c75-42b95914959e.htm) | 新增一个主动连接的请求，将不会收到是否连接成功的信息，当网络中断及奔溃之后，会自动重新连接。  A new active connection request will not receive a message whether the connection is successful. When the network is interrupted and crashed, it will automatically reconnect. (继承自 [CommunicationServer](5a3cb700-c0f6-e35d-6807-dfe8bfecf747.htm)。) |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | [ExtraOnClose](ede241ce-62f0-59aa-45b0-1cc5622a6acf.htm) | 关闭的时候额外执行的功能代码 (继承自 [CommunicationServer](5a3cb700-c0f6-e35d-6807-dfe8bfecf747.htm)。) |
| 受保护的方法 | [ExtraOnStart](2397069b-cc49-1ba5-41c6-ffe84e25219c.htm) | 服务器启动的时候额外执行的功能代码 (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 受保护的方法 | [GetNewNetMessage](89004772-77de-0bf5-785b-46d2a4ce3a3d.htm) | 获取一个新的消息对象的方法，需要在继承类里面进行重写  The method to get a new message object needs to be overridden in the inheritance class (继承自 [CommunicationServer](5a3cb700-c0f6-e35d-6807-dfe8bfecf747.htm)。) |
| 公共方法 | [GetPipeSessions](fbbc274f-38c2-ad5a-dadf-93ac5a7eacff.htm) | 获取管道会话的列表  Get a list of pipeline sessions (继承自 [CommunicationServer](5a3cb700-c0f6-e35d-6807-dfe8bfecf747.htm)。) |
| 公共方法 | [GetTrustedClients](1a468d3b-7c50-7a36-f1f9-76a1d1488a8a.htm) | 获取受信任的客户端列表  Get a list of trusted clients (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法 | [LogDebugMsg](3dc4b43f-2fdc-1d8e-ef40-66319b9f7055.htm) | 记录当前的日志信息 (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法 | [MoveX](e8a255d1-b98d-8aa7-452c-e6e57fb3514b.htm) | 机器人在X轴上移动一小段距离，单位毫米  The robot moves a short distance on the X axis, in millimeters |
| 公共方法 | [MoveY](7bd11b89-d4d2-b19f-88da-dc3218e83a83.htm) | 机器人在Y轴上移动一小段距离，单位毫米  The robot moves a short distance on the Y axis, in millimeters |
| 公共方法 | [MoveZ](41b575e3-242b-7144-4289-0fc59a1e5c2d.htm) | 机器人在Z轴上移动一小段距离，单位毫米  The robot moves a short distance on the Z axis, in millimeters |
| 公共方法 | [RemoveSession(TimeSpan)](863eb441-97a5-7c9a-2545-0665cd0e46b8.htm) | 指定超时时间移除当前的会话列表，只有是TCP的管道（[PipeTcpNet](26b94bd5-2695-36e7-6378-c3d41846be8f.htm)）才需要被移除。  Specify a timeout to remove the current session list. Only TCP pipe ([PipeTcpNet](26b94bd5-2695-36e7-6378-c3d41846be8f.htm)) need to be removed. (继承自 [CommunicationServer](5a3cb700-c0f6-e35d-6807-dfe8bfecf747.htm)。) |
| 公共方法 | [RemoveSession(PipeSession, String)](92e9b2f6-96c3-b586-b554-f8ecd73df580.htm) | 移除一个管道会话  Remove a pipeline session (继承自 [CommunicationServer](5a3cb700-c0f6-e35d-6807-dfe8bfecf747.htm)。) |
| 公共方法 | [RotateX](f154ee9b-d5d5-4ed7-a52e-93885323fb8e.htm) | 机器人在X轴方向上旋转指定角度，单位角度  The robot rotates the specified angle in the X axis direction, the unit angle |
| 公共方法 | [RotateY](2dbb6aa6-1040-6aff-d5cd-eeb5932e1ad5.htm) | 机器人在Y轴方向上旋转指定角度，单位角度  The robot rotates the specified angle in the Y axis direction, the unit angle |
| 公共方法 | [RotateZ](f96cafec-8f12-63d6-0903-4e20a6282476.htm) | 机器人在Z轴方向上旋转指定角度，单位角度  The robot rotates the specified angle in the Z axis direction, the unit angle |
| 公共方法 | [ServerClose](04a70d3a-d427-ed8d-53cf-3d3be56bdbd3.htm) | 关闭服务器的引擎  Shut down the server's engine (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |
| 公共方法 | [ServerStart](6aa1f94d-768c-4711-80e5-a529a5040f98.htm) | 指定端口号来启动服务器的引擎  Specify the port number to start the server's engine (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |
| 公共方法 | [ServerStart(Int32)](06e14b74-3deb-8e46-81c6-688f078375d9.htm) | 指定端口号来启动服务器的引擎  Specify the port number to start the server's engine (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |
| 公共方法 | [ServerStart(Int32, Boolean)](af75f437-3c8e-ac5c-96ab-96ea7bb4cdb4.htm) | 指定端口号来启动服务器的引擎  Specify the port number to start the server's engine (继承自 [CommunicationServer](5a3cb700-c0f6-e35d-6807-dfe8bfecf747.htm)。) |
| 公共方法 | [ServerStart(Int32, Int32)](1feb9368-f6bf-85e6-122a-d0c0d5a67bed.htm) | 指定一个TCP端口及UDP端口，同时启动两种模式的服务器  Specify a TCP port and a UDP port to start the server in both modes at the same time (继承自 [CommunicationServer](5a3cb700-c0f6-e35d-6807-dfe8bfecf747.htm)。) |
| 公共方法 | [SetNetMessage](836e8de1-0461-348c-b56b-2765e37c0ec5.htm) | 设置当前的服务器接收的消息信息 (继承自 [CommunicationServer](5a3cb700-c0f6-e35d-6807-dfe8bfecf747.htm)。) |
| 公共方法 | [SetSslPipeAction](5ecd4347-5dae-126d-f43c-76d4c7390705.htm) | 设置一个SSL的管道操作对象，在管道实例化之后，可以进行一些初始化的属性设置，例如自定义 SslProtocols 枚举 (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |
| 公共方法 | [SetTrustedIpAddress](6204b3e7-af37-7c20-98db-08114431b5ac.htm) | 设置并启动受信任的客户端登录并读写，如果为null，将关闭对客户端的ip验证  Set and start the trusted client login and read and write, if it is null, the client's IP verification will be turned off (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |
| 受保护的方法 | [SocketAcceptExtraCheck](3ed24041-5b60-57a7-7629-f7c49baf8bf9.htm) | 当客户端的socket登录的时候额外检查的操作，并返回操作的结果信息。  The operation is additionally checked when the client's socket logs in, and the result information of the operation is returned. (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |
| 公共方法 | [StartSerialSlave(ActionSerialPort)](dde390ea-45f6-bc3f-be69-3902f6a5641b.htm) | 启动串口的从机服务，使用自定义的初始化方法初始化串口的参数  Start the slave service of serial and initialize the parameters of the serial port using a custom initialization method (继承自 [CommunicationServer](5a3cb700-c0f6-e35d-6807-dfe8bfecf747.htm)。) |
| 公共方法 | [StartSerialSlave(String)](cfed6d77-50f0-eb64-37d0-7c466cc2d7bf.htm) | 启动串口的从机服务，使用默认的参数进行初始化串口，9600波特率，8位数据位，无奇偶校验，1位停止位  Start the slave service of serial, initialize the serial port with default parameters, 9600 baud rate, 8 data bits, no parity, 1 stop bit (继承自 [CommunicationServer](5a3cb700-c0f6-e35d-6807-dfe8bfecf747.htm)。) |
| 公共方法 | [StartSerialSlave(String, Int32)](155bfb40-a6bc-8ca0-3ebc-14307221ba59.htm) | 启动串口的从机服务，使用默认的参数进行初始化串口，8位数据位，无奇偶校验，1位停止位  Start the slave service of serial, initialize the serial port with default parameters, 8 data bits, no parity, 1 stop bit (继承自 [CommunicationServer](5a3cb700-c0f6-e35d-6807-dfe8bfecf747.htm)。) |
| 公共方法 | [StartSerialSlave(String, Int32, Int32, Parity, StopBits)](7045bee6-9c7f-65e0-191a-205c360dd42c.htm) | 启动串口的从机服务，使用指定的参数进行初始化串口，指定数据位，指定奇偶校验，指定停止位 (继承自 [CommunicationServer](5a3cb700-c0f6-e35d-6807-dfe8bfecf747.htm)。) |
| 受保护的方法 | [ThreadPoolLogin](7ebe2e21-f51e-3a08-a499-624462a0fede.htm) | 当客户端连接到服务器，并听过额外的检查后，进行回调的方法  Callback method when the client connects to the server and has heard additional checks (继承自 [CommunicationServer](5a3cb700-c0f6-e35d-6807-dfe8bfecf747.htm)。) |
| 公共方法 | [ToString](8010794a-c51c-3618-8667-024a32e2cadd.htm) | (重写 [CommunicationServerToString](955ea032-e1fd-11ff-9ec8-72cb0572d36a.htm).) |
| 公共方法 | [UseSSL(X509Certificate)](3577f66a-4337-6d6c-0f73-1f904f6bed38.htm) | 使用SSL通信，传递一个证书的对象 (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |
| 公共方法 | [UseSSL(String, String)](e6f37bac-c743-e05a-729a-0e4a64daf486.htm) | 使用SSL通信，传递一个证书的路径，以及证书的密码 (继承自 [CommunicationTcpServer](8752ca5b-81d1-9516-3e67-4bf2531c6fda.htm)。) |
| 公共方法 | [Write](061ebd46-cffc-5e79-1de2-0044dca27b70.htm) | 将指定的命令写入机器人，该命令是完全自定义的，需要遵循机器人的通讯协议，在写入之前，需要调用[ServerStart(Int32)](06e14b74-3deb-8e46-81c6-688f078375d9.htm) 方法  Write the specified command to the robot. The command is completely customized and needs to follow the robot's communication protocol. Before writing, you need to call the [ServerStart(Int32)](06e14b74-3deb-8e46-81c6-688f078375d9.htm) |
| 公共方法 | [WriteIncrementPos(Double)](fad0645a-8622-ef21-05f5-2a5c93bc0907.htm) | 将指定的增量写入机器人，需要指定6个参数，位置和角度信息，其中位置单位为mm，角度单位为°  To write the specified increment to the robot, you need to specify 6 parameters, position and angle information, where the position unit is mm and the angle unit is ° |
| 公共方法 | [WriteIncrementPos(Double, Double, Double, Double, Double, Double)](b4e11f16-cdbd-1bbb-941e-d6f36bfba447.htm) | 将指定的增量写入机器人，需要指定6个参数，位置和角度信息，其中位置单位为mm，角度单位为°  To write the specified increment to the robot, you need to specify 6 parameters, position and angle information, where the position unit is mm and the angle unit is ° |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[HyundaiUdpNet 类](63e7fe67-d9ec-efd5-a7d8-a8cde49538e2.htm)

[HslCommunication.Robot.Hyundai 命名空间](caddf3fe-6919-18da-f3ee-fc10fac8c1bf.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## MoveX 方法 

[原文連結](http://api.hslcommunication.cn/html/e8a255d1-b98d-8aa7-452c-e6e57fb3514b.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.Hyundai](../html/caddf3fe-6919-18da-f3ee-fc10fac8c1bf.htm "HslCommunication.Robot.Hyundai")

[HyundaiUdpNet 类](../html/63e7fe67-d9ec-efd5-a7d8-a8cde49538e2.htm "HyundaiUdpNet 类")

[HyundaiUdpNet 方法](../html/8d5c4acb-4433-86d1-b295-511cb669b68f.htm "HyundaiUdpNet 方法")

[MoveX 方法](../html/e8a255d1-b98d-8aa7-452c-e6e57fb3514b.htm "MoveX 方法 ")

[MoveY 方法](../html/7bd11b89-d4d2-b19f-88da-dc3218e83a83.htm "MoveY 方法 ")

[MoveZ 方法](../html/41b575e3-242b-7144-4289-0fc59a1e5c2d.htm "MoveZ 方法 ")

[RotateX 方法](../html/f154ee9b-d5d5-4ed7-a52e-93885323fb8e.htm "RotateX 方法 ")

[RotateY 方法](../html/2dbb6aa6-1040-6aff-d5cd-eeb5932e1ad5.htm "RotateY 方法 ")

[RotateZ 方法](../html/f96cafec-8f12-63d6-0903-4e20a6282476.htm "RotateZ 方法 ")

[ToString 方法](../html/8010794a-c51c-3618-8667-024a32e2cadd.htm "ToString 方法 ")

[Write 方法](../html/061ebd46-cffc-5e79-1de2-0044dca27b70.htm "Write 方法 ")

[WriteIncrementPos 方法](../html/50870834-55a9-e75b-e524-ffb8a32046f3.htm "WriteIncrementPos 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HyundaiUdpNetMoveX 方法 |

机器人在X轴上移动一小段距离，单位毫米  
The robot moves a short distance on the X axis, in millimeters

**命名空间：**
 [HslCommunication.Robot.Hyundai](caddf3fe-6919-18da-f3ee-fc10fac8c1bf.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public OperateResult MoveX(
	double value
)
```

```
Public Function MoveX ( 
	value As Double
) As OperateResult
```

```
public:
OperateResult^ MoveX(
	double value
)
```

```
member MoveX : 
        value : float -> OperateResult 
```

#### 参数

value
:   类型：SystemDouble  
    移动距离，单位毫米

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否写入成功

![](../icons/SectionExpanded.png)参见

#### 引用

[HyundaiUdpNet 类](63e7fe67-d9ec-efd5-a7d8-a8cde49538e2.htm)

[HslCommunication.Robot.Hyundai 命名空间](caddf3fe-6919-18da-f3ee-fc10fac8c1bf.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## MoveY 方法 

[原文連結](http://api.hslcommunication.cn/html/7bd11b89-d4d2-b19f-88da-dc3218e83a83.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.Hyundai](../html/caddf3fe-6919-18da-f3ee-fc10fac8c1bf.htm "HslCommunication.Robot.Hyundai")

[HyundaiUdpNet 类](../html/63e7fe67-d9ec-efd5-a7d8-a8cde49538e2.htm "HyundaiUdpNet 类")

[HyundaiUdpNet 方法](../html/8d5c4acb-4433-86d1-b295-511cb669b68f.htm "HyundaiUdpNet 方法")

[MoveX 方法](../html/e8a255d1-b98d-8aa7-452c-e6e57fb3514b.htm "MoveX 方法 ")

[MoveY 方法](../html/7bd11b89-d4d2-b19f-88da-dc3218e83a83.htm "MoveY 方法 ")

[MoveZ 方法](../html/41b575e3-242b-7144-4289-0fc59a1e5c2d.htm "MoveZ 方法 ")

[RotateX 方法](../html/f154ee9b-d5d5-4ed7-a52e-93885323fb8e.htm "RotateX 方法 ")

[RotateY 方法](../html/2dbb6aa6-1040-6aff-d5cd-eeb5932e1ad5.htm "RotateY 方法 ")

[RotateZ 方法](../html/f96cafec-8f12-63d6-0903-4e20a6282476.htm "RotateZ 方法 ")

[ToString 方法](../html/8010794a-c51c-3618-8667-024a32e2cadd.htm "ToString 方法 ")

[Write 方法](../html/061ebd46-cffc-5e79-1de2-0044dca27b70.htm "Write 方法 ")

[WriteIncrementPos 方法](../html/50870834-55a9-e75b-e524-ffb8a32046f3.htm "WriteIncrementPos 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HyundaiUdpNetMoveY 方法 |

机器人在Y轴上移动一小段距离，单位毫米  
The robot moves a short distance on the Y axis, in millimeters

**命名空间：**
 [HslCommunication.Robot.Hyundai](caddf3fe-6919-18da-f3ee-fc10fac8c1bf.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public OperateResult MoveY(
	double value
)
```

```
Public Function MoveY ( 
	value As Double
) As OperateResult
```

```
public:
OperateResult^ MoveY(
	double value
)
```

```
member MoveY : 
        value : float -> OperateResult 
```

#### 参数

value
:   类型：SystemDouble  
    移动距离，单位毫米

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否写入成功

![](../icons/SectionExpanded.png)参见

#### 引用

[HyundaiUdpNet 类](63e7fe67-d9ec-efd5-a7d8-a8cde49538e2.htm)

[HslCommunication.Robot.Hyundai 命名空间](caddf3fe-6919-18da-f3ee-fc10fac8c1bf.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## MoveZ 方法 

[原文連結](http://api.hslcommunication.cn/html/41b575e3-242b-7144-4289-0fc59a1e5c2d.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.Hyundai](../html/caddf3fe-6919-18da-f3ee-fc10fac8c1bf.htm "HslCommunication.Robot.Hyundai")

[HyundaiUdpNet 类](../html/63e7fe67-d9ec-efd5-a7d8-a8cde49538e2.htm "HyundaiUdpNet 类")

[HyundaiUdpNet 方法](../html/8d5c4acb-4433-86d1-b295-511cb669b68f.htm "HyundaiUdpNet 方法")

[MoveX 方法](../html/e8a255d1-b98d-8aa7-452c-e6e57fb3514b.htm "MoveX 方法 ")

[MoveY 方法](../html/7bd11b89-d4d2-b19f-88da-dc3218e83a83.htm "MoveY 方法 ")

[MoveZ 方法](../html/41b575e3-242b-7144-4289-0fc59a1e5c2d.htm "MoveZ 方法 ")

[RotateX 方法](../html/f154ee9b-d5d5-4ed7-a52e-93885323fb8e.htm "RotateX 方法 ")

[RotateY 方法](../html/2dbb6aa6-1040-6aff-d5cd-eeb5932e1ad5.htm "RotateY 方法 ")

[RotateZ 方法](../html/f96cafec-8f12-63d6-0903-4e20a6282476.htm "RotateZ 方法 ")

[ToString 方法](../html/8010794a-c51c-3618-8667-024a32e2cadd.htm "ToString 方法 ")

[Write 方法](../html/061ebd46-cffc-5e79-1de2-0044dca27b70.htm "Write 方法 ")

[WriteIncrementPos 方法](../html/50870834-55a9-e75b-e524-ffb8a32046f3.htm "WriteIncrementPos 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HyundaiUdpNetMoveZ 方法 |

机器人在Z轴上移动一小段距离，单位毫米  
The robot moves a short distance on the Z axis, in millimeters

**命名空间：**
 [HslCommunication.Robot.Hyundai](caddf3fe-6919-18da-f3ee-fc10fac8c1bf.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public OperateResult MoveZ(
	double value
)
```

```
Public Function MoveZ ( 
	value As Double
) As OperateResult
```

```
public:
OperateResult^ MoveZ(
	double value
)
```

```
member MoveZ : 
        value : float -> OperateResult 
```

#### 参数

value
:   类型：SystemDouble  
    移动距离，单位毫米

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否写入成功

![](../icons/SectionExpanded.png)参见

#### 引用

[HyundaiUdpNet 类](63e7fe67-d9ec-efd5-a7d8-a8cde49538e2.htm)

[HslCommunication.Robot.Hyundai 命名空间](caddf3fe-6919-18da-f3ee-fc10fac8c1bf.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## RotateX 方法 

[原文連結](http://api.hslcommunication.cn/html/f154ee9b-d5d5-4ed7-a52e-93885323fb8e.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.Hyundai](../html/caddf3fe-6919-18da-f3ee-fc10fac8c1bf.htm "HslCommunication.Robot.Hyundai")

[HyundaiUdpNet 类](../html/63e7fe67-d9ec-efd5-a7d8-a8cde49538e2.htm "HyundaiUdpNet 类")

[HyundaiUdpNet 方法](../html/8d5c4acb-4433-86d1-b295-511cb669b68f.htm "HyundaiUdpNet 方法")

[MoveX 方法](../html/e8a255d1-b98d-8aa7-452c-e6e57fb3514b.htm "MoveX 方法 ")

[MoveY 方法](../html/7bd11b89-d4d2-b19f-88da-dc3218e83a83.htm "MoveY 方法 ")

[MoveZ 方法](../html/41b575e3-242b-7144-4289-0fc59a1e5c2d.htm "MoveZ 方法 ")

[RotateX 方法](../html/f154ee9b-d5d5-4ed7-a52e-93885323fb8e.htm "RotateX 方法 ")

[RotateY 方法](../html/2dbb6aa6-1040-6aff-d5cd-eeb5932e1ad5.htm "RotateY 方法 ")

[RotateZ 方法](../html/f96cafec-8f12-63d6-0903-4e20a6282476.htm "RotateZ 方法 ")

[ToString 方法](../html/8010794a-c51c-3618-8667-024a32e2cadd.htm "ToString 方法 ")

[Write 方法](../html/061ebd46-cffc-5e79-1de2-0044dca27b70.htm "Write 方法 ")

[WriteIncrementPos 方法](../html/50870834-55a9-e75b-e524-ffb8a32046f3.htm "WriteIncrementPos 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HyundaiUdpNetRotateX 方法 |

机器人在X轴方向上旋转指定角度，单位角度  
The robot rotates the specified angle in the X axis direction, the unit angle

**命名空间：**
 [HslCommunication.Robot.Hyundai](caddf3fe-6919-18da-f3ee-fc10fac8c1bf.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public OperateResult RotateX(
	double value
)
```

```
Public Function RotateX ( 
	value As Double
) As OperateResult
```

```
public:
OperateResult^ RotateX(
	double value
)
```

```
member RotateX : 
        value : float -> OperateResult 
```

#### 参数

value
:   类型：SystemDouble  
    旋转角度，单位角度

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否写入成功

![](../icons/SectionExpanded.png)参见

#### 引用

[HyundaiUdpNet 类](63e7fe67-d9ec-efd5-a7d8-a8cde49538e2.htm)

[HslCommunication.Robot.Hyundai 命名空间](caddf3fe-6919-18da-f3ee-fc10fac8c1bf.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## RotateY 方法 

[原文連結](http://api.hslcommunication.cn/html/2dbb6aa6-1040-6aff-d5cd-eeb5932e1ad5.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.Hyundai](../html/caddf3fe-6919-18da-f3ee-fc10fac8c1bf.htm "HslCommunication.Robot.Hyundai")

[HyundaiUdpNet 类](../html/63e7fe67-d9ec-efd5-a7d8-a8cde49538e2.htm "HyundaiUdpNet 类")

[HyundaiUdpNet 方法](../html/8d5c4acb-4433-86d1-b295-511cb669b68f.htm "HyundaiUdpNet 方法")

[MoveX 方法](../html/e8a255d1-b98d-8aa7-452c-e6e57fb3514b.htm "MoveX 方法 ")

[MoveY 方法](../html/7bd11b89-d4d2-b19f-88da-dc3218e83a83.htm "MoveY 方法 ")

[MoveZ 方法](../html/41b575e3-242b-7144-4289-0fc59a1e5c2d.htm "MoveZ 方法 ")

[RotateX 方法](../html/f154ee9b-d5d5-4ed7-a52e-93885323fb8e.htm "RotateX 方法 ")

[RotateY 方法](../html/2dbb6aa6-1040-6aff-d5cd-eeb5932e1ad5.htm "RotateY 方法 ")

[RotateZ 方法](../html/f96cafec-8f12-63d6-0903-4e20a6282476.htm "RotateZ 方法 ")

[ToString 方法](../html/8010794a-c51c-3618-8667-024a32e2cadd.htm "ToString 方法 ")

[Write 方法](../html/061ebd46-cffc-5e79-1de2-0044dca27b70.htm "Write 方法 ")

[WriteIncrementPos 方法](../html/50870834-55a9-e75b-e524-ffb8a32046f3.htm "WriteIncrementPos 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HyundaiUdpNetRotateY 方法 |

机器人在Y轴方向上旋转指定角度，单位角度  
The robot rotates the specified angle in the Y axis direction, the unit angle

**命名空间：**
 [HslCommunication.Robot.Hyundai](caddf3fe-6919-18da-f3ee-fc10fac8c1bf.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public OperateResult RotateY(
	double value
)
```

```
Public Function RotateY ( 
	value As Double
) As OperateResult
```

```
public:
OperateResult^ RotateY(
	double value
)
```

```
member RotateY : 
        value : float -> OperateResult 
```

#### 参数

value
:   类型：SystemDouble  
    旋转角度，单位角度

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否写入成功

![](../icons/SectionExpanded.png)参见

#### 引用

[HyundaiUdpNet 类](63e7fe67-d9ec-efd5-a7d8-a8cde49538e2.htm)

[HslCommunication.Robot.Hyundai 命名空间](caddf3fe-6919-18da-f3ee-fc10fac8c1bf.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## RotateZ 方法 

[原文連結](http://api.hslcommunication.cn/html/f96cafec-8f12-63d6-0903-4e20a6282476.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.Hyundai](../html/caddf3fe-6919-18da-f3ee-fc10fac8c1bf.htm "HslCommunication.Robot.Hyundai")

[HyundaiUdpNet 类](../html/63e7fe67-d9ec-efd5-a7d8-a8cde49538e2.htm "HyundaiUdpNet 类")

[HyundaiUdpNet 方法](../html/8d5c4acb-4433-86d1-b295-511cb669b68f.htm "HyundaiUdpNet 方法")

[MoveX 方法](../html/e8a255d1-b98d-8aa7-452c-e6e57fb3514b.htm "MoveX 方法 ")

[MoveY 方法](../html/7bd11b89-d4d2-b19f-88da-dc3218e83a83.htm "MoveY 方法 ")

[MoveZ 方法](../html/41b575e3-242b-7144-4289-0fc59a1e5c2d.htm "MoveZ 方法 ")

[RotateX 方法](../html/f154ee9b-d5d5-4ed7-a52e-93885323fb8e.htm "RotateX 方法 ")

[RotateY 方法](../html/2dbb6aa6-1040-6aff-d5cd-eeb5932e1ad5.htm "RotateY 方法 ")

[RotateZ 方法](../html/f96cafec-8f12-63d6-0903-4e20a6282476.htm "RotateZ 方法 ")

[ToString 方法](../html/8010794a-c51c-3618-8667-024a32e2cadd.htm "ToString 方法 ")

[Write 方法](../html/061ebd46-cffc-5e79-1de2-0044dca27b70.htm "Write 方法 ")

[WriteIncrementPos 方法](../html/50870834-55a9-e75b-e524-ffb8a32046f3.htm "WriteIncrementPos 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HyundaiUdpNetRotateZ 方法 |

机器人在Z轴方向上旋转指定角度，单位角度  
The robot rotates the specified angle in the Z axis direction, the unit angle

**命名空间：**
 [HslCommunication.Robot.Hyundai](caddf3fe-6919-18da-f3ee-fc10fac8c1bf.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public OperateResult RotateZ(
	double value
)
```

```
Public Function RotateZ ( 
	value As Double
) As OperateResult
```

```
public:
OperateResult^ RotateZ(
	double value
)
```

```
member RotateZ : 
        value : float -> OperateResult 
```

#### 参数

value
:   类型：SystemDouble  
    旋转角度，单位角度

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否写入成功

![](../icons/SectionExpanded.png)参见

#### 引用

[HyundaiUdpNet 类](63e7fe67-d9ec-efd5-a7d8-a8cde49538e2.htm)

[HslCommunication.Robot.Hyundai 命名空间](caddf3fe-6919-18da-f3ee-fc10fac8c1bf.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ToString 方法 

[原文連結](http://api.hslcommunication.cn/html/8010794a-c51c-3618-8667-024a32e2cadd.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.Hyundai](../html/caddf3fe-6919-18da-f3ee-fc10fac8c1bf.htm "HslCommunication.Robot.Hyundai")

[HyundaiUdpNet 类](../html/63e7fe67-d9ec-efd5-a7d8-a8cde49538e2.htm "HyundaiUdpNet 类")

[HyundaiUdpNet 方法](../html/8d5c4acb-4433-86d1-b295-511cb669b68f.htm "HyundaiUdpNet 方法")

[MoveX 方法](../html/e8a255d1-b98d-8aa7-452c-e6e57fb3514b.htm "MoveX 方法 ")

[MoveY 方法](../html/7bd11b89-d4d2-b19f-88da-dc3218e83a83.htm "MoveY 方法 ")

[MoveZ 方法](../html/41b575e3-242b-7144-4289-0fc59a1e5c2d.htm "MoveZ 方法 ")

[RotateX 方法](../html/f154ee9b-d5d5-4ed7-a52e-93885323fb8e.htm "RotateX 方法 ")

[RotateY 方法](../html/2dbb6aa6-1040-6aff-d5cd-eeb5932e1ad5.htm "RotateY 方法 ")

[RotateZ 方法](../html/f96cafec-8f12-63d6-0903-4e20a6282476.htm "RotateZ 方法 ")

[ToString 方法](../html/8010794a-c51c-3618-8667-024a32e2cadd.htm "ToString 方法 ")

[Write 方法](../html/061ebd46-cffc-5e79-1de2-0044dca27b70.htm "Write 方法 ")

[WriteIncrementPos 方法](../html/50870834-55a9-e75b-e524-ffb8a32046f3.htm "WriteIncrementPos 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HyundaiUdpNetToString 方法 |

[缺少 "M:HslCommunication.Robot.Hyundai.HyundaiUdpNet.ToString" 的 <summary> 文档]

**命名空间：**
 [HslCommunication.Robot.Hyundai](caddf3fe-6919-18da-f3ee-fc10fac8c1bf.htm)  
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

[缺少 "M:HslCommunication.Robot.Hyundai.HyundaiUdpNet.ToString" 的 <returns> 文档]

![](../icons/SectionExpanded.png)参见

#### 引用

[HyundaiUdpNet 类](63e7fe67-d9ec-efd5-a7d8-a8cde49538e2.htm)

[HslCommunication.Robot.Hyundai 命名空间](caddf3fe-6919-18da-f3ee-fc10fac8c1bf.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Write 方法 

[原文連結](http://api.hslcommunication.cn/html/061ebd46-cffc-5e79-1de2-0044dca27b70.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.Hyundai](../html/caddf3fe-6919-18da-f3ee-fc10fac8c1bf.htm "HslCommunication.Robot.Hyundai")

[HyundaiUdpNet 类](../html/63e7fe67-d9ec-efd5-a7d8-a8cde49538e2.htm "HyundaiUdpNet 类")

[HyundaiUdpNet 方法](../html/8d5c4acb-4433-86d1-b295-511cb669b68f.htm "HyundaiUdpNet 方法")

[MoveX 方法](../html/e8a255d1-b98d-8aa7-452c-e6e57fb3514b.htm "MoveX 方法 ")

[MoveY 方法](../html/7bd11b89-d4d2-b19f-88da-dc3218e83a83.htm "MoveY 方法 ")

[MoveZ 方法](../html/41b575e3-242b-7144-4289-0fc59a1e5c2d.htm "MoveZ 方法 ")

[RotateX 方法](../html/f154ee9b-d5d5-4ed7-a52e-93885323fb8e.htm "RotateX 方法 ")

[RotateY 方法](../html/2dbb6aa6-1040-6aff-d5cd-eeb5932e1ad5.htm "RotateY 方法 ")

[RotateZ 方法](../html/f96cafec-8f12-63d6-0903-4e20a6282476.htm "RotateZ 方法 ")

[ToString 方法](../html/8010794a-c51c-3618-8667-024a32e2cadd.htm "ToString 方法 ")

[Write 方法](../html/061ebd46-cffc-5e79-1de2-0044dca27b70.htm "Write 方法 ")

[WriteIncrementPos 方法](../html/50870834-55a9-e75b-e524-ffb8a32046f3.htm "WriteIncrementPos 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HyundaiUdpNetWrite 方法 |

将指定的命令写入机器人，该命令是完全自定义的，需要遵循机器人的通讯协议，在写入之前，需要调用[ServerStart(Int32)](06e14b74-3deb-8e46-81c6-688f078375d9.htm) 方法  
Write the specified command to the robot. The command is completely customized and needs to follow the robot's communication protocol.
Before writing, you need to call the [ServerStart(Int32)](06e14b74-3deb-8e46-81c6-688f078375d9.htm)

**命名空间：**
 [HslCommunication.Robot.Hyundai](caddf3fe-6919-18da-f3ee-fc10fac8c1bf.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public OperateResult Write(
	HyundaiData data
)
```

```
Public Function Write ( 
	data As HyundaiData
) As OperateResult
```

```
public:
OperateResult^ Write(
	HyundaiData^ data
)
```

```
member Write : 
        data : HyundaiData -> OperateResult 
```

#### 参数

data
:   类型：[HslCommunication.Robot.HyundaiHyundaiData](86650e61-38f2-c3ee-2fcd-c45d8d120f9f.htm)  
    机器人数据

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否写入成功

![](../icons/SectionExpanded.png)参见

#### 引用

[HyundaiUdpNet 类](63e7fe67-d9ec-efd5-a7d8-a8cde49538e2.htm)

[HslCommunication.Robot.Hyundai 命名空间](caddf3fe-6919-18da-f3ee-fc10fac8c1bf.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## WriteIncrementPos 方法 

[原文連結](http://api.hslcommunication.cn/html/50870834-55a9-e75b-e524-ffb8a32046f3.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.Hyundai](../html/caddf3fe-6919-18da-f3ee-fc10fac8c1bf.htm "HslCommunication.Robot.Hyundai")

[HyundaiUdpNet 类](../html/63e7fe67-d9ec-efd5-a7d8-a8cde49538e2.htm "HyundaiUdpNet 类")

[HyundaiUdpNet 方法](../html/8d5c4acb-4433-86d1-b295-511cb669b68f.htm "HyundaiUdpNet 方法")

[WriteIncrementPos 方法](../html/50870834-55a9-e75b-e524-ffb8a32046f3.htm "WriteIncrementPos 方法 ")

[WriteIncrementPos 方法 (Double[])](../html/fad0645a-8622-ef21-05f5-2a5c93bc0907.htm "WriteIncrementPos 方法 (Double[])")

[WriteIncrementPos 方法 (Double, Double, Double, Double, Double, Double)](../html/b4e11f16-cdbd-1bbb-941e-d6f36bfba447.htm "WriteIncrementPos 方法 (Double, Double, Double, Double, Double, Double)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HyundaiUdpNetWriteIncrementPos 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [WriteIncrementPos(Double)](fad0645a-8622-ef21-05f5-2a5c93bc0907.htm) | 将指定的增量写入机器人，需要指定6个参数，位置和角度信息，其中位置单位为mm，角度单位为°  To write the specified increment to the robot, you need to specify 6 parameters, position and angle information, where the position unit is mm and the angle unit is ° |
| 公共方法 | [WriteIncrementPos(Double, Double, Double, Double, Double, Double)](b4e11f16-cdbd-1bbb-941e-d6f36bfba447.htm) | 将指定的增量写入机器人，需要指定6个参数，位置和角度信息，其中位置单位为mm，角度单位为°  To write the specified increment to the robot, you need to specify 6 parameters, position and angle information, where the position unit is mm and the angle unit is ° |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[HyundaiUdpNet 类](63e7fe67-d9ec-efd5-a7d8-a8cde49538e2.htm)

[HslCommunication.Robot.Hyundai 命名空间](caddf3fe-6919-18da-f3ee-fc10fac8c1bf.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## WriteIncrementPos 方法 (Double[])

[原文連結](http://api.hslcommunication.cn/html/fad0645a-8622-ef21-05f5-2a5c93bc0907.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Robot.Hyundai](../html/caddf3fe-6919-18da-f3ee-fc10fac8c1bf.htm "HslCommunication.Robot.Hyundai")

[HyundaiUdpNet 类](../html/63e7fe67-d9ec-efd5-a7d8-a8cde49538e2.htm "HyundaiUdpNet 类")

[HyundaiUdpNet 方法](../html/8d5c4acb-4433-86d1-b295-511cb669b68f.htm "HyundaiUdpNet 方法")

[WriteIncrementPos 方法](../html/50870834-55a9-e75b-e524-ffb8a32046f3.htm "WriteIncrementPos 方法 ")

[WriteIncrementPos 方法 (Double[])](../html/fad0645a-8622-ef21-05f5-2a5c93bc0907.htm "WriteIncrementPos 方法 (Double[])")

[WriteIncrementPos 方法 (Double, Double, Double, Double, Double, Double)](../html/b4e11f16-cdbd-1bbb-941e-d6f36bfba447.htm "WriteIncrementPos 方法 (Double, Double, Double, Double, Double, Double)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HyundaiUdpNetWriteIncrementPos 方法 (Double) |

将指定的增量写入机器人，需要指定6个参数，位置和角度信息，其中位置单位为mm，角度单位为°  
To write the specified increment to the robot, you need to specify 6 parameters, position and angle information, where the position unit is mm and the angle unit is °

**命名空间：**
 [HslCommunication.Robot.Hyundai](caddf3fe-6919-18da-f3ee-fc10fac8c1bf.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public OperateResult WriteIncrementPos(
	double[] pos
)
```

```
Public Function WriteIncrementPos ( 
	pos As Double()
) As OperateResult
```

```
public:
OperateResult^ WriteIncrementPos(
	array<double>^ pos
)
```

```
member WriteIncrementPos : 
        pos : float[] -> OperateResult 
```

#### 参数

pos
:   类型：SystemDouble  
    增量的数组信息

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否写入机器人成功

![](../icons/SectionExpanded.png)参见

#### 引用

[HyundaiUdpNet 类](63e7fe67-d9ec-efd5-a7d8-a8cde49538e2.htm)

[WriteIncrementPos 重载](50870834-55a9-e75b-e524-ffb8a32046f3.htm)

[HslCommunication.Robot.Hyundai 命名空间](caddf3fe-6919-18da-f3ee-fc10fac8c1bf.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)