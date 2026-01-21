# HslCommunication - HslCommunication.Profinet.Inovance

> 分類頁數: 30



---
## HslCommunication.Profinet.Inovance

[原文連結](http://api.hslcommunication.cn/html/aa61d757-9239-26ce-525f-b7e52699e229.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Inovance](../html/aa61d757-9239-26ce-525f-b7e52699e229.htm "HslCommunication.Profinet.Inovance")

[InovanceHelper 类](../html/944b43ff-5625-c33b-3105-1c1b8db259d9.htm "InovanceHelper 类")

[InovanceSerial 类](../html/434659f2-7d84-27a2-72b5-aebd489024da.htm "InovanceSerial 类")

[InovanceSerialOverTcp 类](../html/f82b954d-1a19-0b0e-f705-677d531e17f6.htm "InovanceSerialOverTcp 类")

[InovanceSeries 枚举](../html/624b6c23-3fff-4ce9-aa3c-72b3897aec74.htm "InovanceSeries 枚举")

[InovanceTcpNet 类](../html/a395ef3c-2bbb-0d7e-733d-9db9475ca1e3.htm "InovanceTcpNet 类")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslCommunication.Profinet.Inovance 命名空间 |

[缺少 "N:HslCommunication.Profinet.Inovance" 的 <summary> 文档]

![](../icons/SectionExpanded.png)类

|  | 类 | 说明 |
| --- | --- | --- |
| 公共类 | [InovanceHelper](944b43ff-5625-c33b-3105-1c1b8db259d9.htm) | 汇川PLC的辅助类，提供一些地址解析的方法  Auxiliary class of Yaskawa robot, providing some methods of address resolution |
| 公共类 | [InovanceSerial](434659f2-7d84-27a2-72b5-aebd489024da.htm) | 汇川的串口通信协议，A适用于AM400、 AM400\_800、 AC800、H3U, XP, H5U 等系列底层走的是MODBUS-RTU协议，地址说明参见标记  Huichuan's serial communication protocol is applicable to AM400, AM400\_800, AC800 and other series. The bottom layer is MODBUS-RTU protocol. For the address description, please refer to the mark |
| 公共类 | [InovanceSerialOverTcp](f82b954d-1a19-0b0e-f705-677d531e17f6.htm) | 汇川的串口转网口通信协议，A适用于AM400、 AM400\_800、 AC800、H3U, XP, H5U 等系列底层走的是MODBUS-RTU协议，地址说明参见标记  Huichuan's serial communication protocol is applicable to AM400, AM400\_800, AC800 and other series. The bottom layer is MODBUS-RTU protocol. For the address description, please refer to the mark |
| 公共类 | [InovanceTcpNet](a395ef3c-2bbb-0d7e-733d-9db9475ca1e3.htm) | 汇川的网口通信协议，A适用于AM400、 AM400\_800、 AC800、H3U, XP, H5U 等系列底层走的是MODBUS-TCP协议，地址说明参见标记  Huichuan's serial communication protocol is applicable to AM400, AM400\_800, AC800 and other series. The bottom layer is MODBUS-TCP protocol. For the address description, please refer to the mark |

![](../icons/SectionExpanded.png)枚举

|  | 枚举 | 说明 |
| --- | --- | --- |
| 公共枚举 | [InovanceSeries](624b6c23-3fff-4ce9-aa3c-72b3897aec74.htm) | 汇川的系列枚举信息 |

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## InovanceHelper 类

[原文連結](http://api.hslcommunication.cn/html/944b43ff-5625-c33b-3105-1c1b8db259d9.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Inovance](../html/aa61d757-9239-26ce-525f-b7e52699e229.htm "HslCommunication.Profinet.Inovance")

[InovanceHelper 类](../html/944b43ff-5625-c33b-3105-1c1b8db259d9.htm "InovanceHelper 类")

[InovanceHelper 构造函数](../html/834fe363-1951-ef53-bf2d-b9578f4f3c5f.htm "InovanceHelper 构造函数 ")

[InovanceHelper 方法](../html/05a274bc-a7e4-39bb-7f78-0383d043cf99.htm "InovanceHelper 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| InovanceHelper 类 |

汇川PLC的辅助类，提供一些地址解析的方法  
Auxiliary class of Yaskawa robot, providing some methods of address resolution

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  HslCommunication.Profinet.InovanceInovanceHelper

**命名空间：**
 [HslCommunication.Profinet.Inovance](aa61d757-9239-26ce-525f-b7e52699e229.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class InovanceHelper
```

```
Public Class InovanceHelper
```

```
public ref class InovanceHelper
```

```
type InovanceHelper =  class end
```

InovanceHelper 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [InovanceHelper](834fe363-1951-ef53-bf2d-b9578f4f3c5f.htm) | 初始化 InovanceHelper 类的一个新实例 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法静态成员 | [PraseInovanceAddress](f3d46fa3-fc1f-bb51-0efb-c0331acb938a.htm) | 根据汇川PLC的地址，解析出转换后的modbus协议信息，适用AM,H3U,H5U系列的PLC  According to the address of Inovance PLC, analyze the converted modbus protocol information, which is suitable for AM, H3U, H5U series PLC |
| 公共方法静态成员 | [PraseInovanceAMAddress](9c7e153c-c315-4461-62dc-7ba419a0b6ac.htm) | 根据汇川PLC的地址，解析出转换后的modbus协议信息，适用AM,H3U,H5U系列的PLC  According to the address of Inovance PLC, analyze the converted modbus protocol information, which is suitable for AM, H3U, H5U series PLC |
| 公共方法静态成员 | [PraseInovanceH3UAddress](2a040cb9-febb-6a01-473e-e68edb5897fe.htm) | 根据汇川PLC的地址，解析出转换后的modbus协议信息，适用AM,H3U,H5U系列的PLC  According to the address of Inovance PLC, analyze the converted modbus protocol information, which is suitable for AM, H3U, H5U series PLC |
| 公共方法静态成员 | [PraseInovanceH5UAddress](594cc79d-ad03-e922-4644-36674190c3d9.htm) | 根据汇川PLC的地址，解析出转换后的modbus协议信息，适用AM,H3U,H5U系列的PLC  According to the address of Inovance PLC, analyze the converted modbus protocol information, which is suitable for AM, H3U, H5U series PLC |
| 公共方法静态成员 | [ReadByte](e20dbf83-6657-b7db-c5fb-4f8226e7dee7.htm) | 按照字节读取汇川M地址的数据，地址示例： MB100，MB101，需要注意的是，MB100 及 MB101 的地址是 MW50 的数据。  Read the data of Inovance M address according to the byte, address example: MB100, MB101, it should be noted that the addresses of MB100 and MB101 are the data of MW50. |
| 公共方法静态成员 | [ReadByteAsync](3afd1bcd-229d-1cd0-402b-26ff82e0c480.htm) | 按照字节读取汇川M地址的数据，地址示例： MB100，MB101，需要注意的是，MB100 及 MB101 的地址是 MW50 的数据。  Read the data of Inovance M address according to the byte, address example: MB100, MB101, it should be noted that the addresses of MB100 and MB101 are the data of MW50. |
| 公共方法 | ToString | (继承自 Object。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.Profinet.Inovance 命名空间](aa61d757-9239-26ce-525f-b7e52699e229.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## InovanceHelper 构造函数 

[原文連結](http://api.hslcommunication.cn/html/834fe363-1951-ef53-bf2d-b9578f4f3c5f.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Inovance](../html/aa61d757-9239-26ce-525f-b7e52699e229.htm "HslCommunication.Profinet.Inovance")

[InovanceHelper 类](../html/944b43ff-5625-c33b-3105-1c1b8db259d9.htm "InovanceHelper 类")

[InovanceHelper 构造函数](../html/834fe363-1951-ef53-bf2d-b9578f4f3c5f.htm "InovanceHelper 构造函数 ")

[InovanceHelper 方法](../html/05a274bc-a7e4-39bb-7f78-0383d043cf99.htm "InovanceHelper 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| InovanceHelper 构造函数 |

初始化 [InovanceHelper](944b43ff-5625-c33b-3105-1c1b8db259d9.htm) 类的一个新实例

**命名空间：**
 [HslCommunication.Profinet.Inovance](aa61d757-9239-26ce-525f-b7e52699e229.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public InovanceHelper()
```

```
Public Sub New
```

```
public:
InovanceHelper()
```

```
new : unit -> InovanceHelper
```

![](../icons/SectionExpanded.png)参见

#### 引用

[InovanceHelper 类](944b43ff-5625-c33b-3105-1c1b8db259d9.htm)

[HslCommunication.Profinet.Inovance 命名空间](aa61d757-9239-26ce-525f-b7e52699e229.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## InovanceHelper 方法

[原文連結](http://api.hslcommunication.cn/html/05a274bc-a7e4-39bb-7f78-0383d043cf99.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Inovance](../html/aa61d757-9239-26ce-525f-b7e52699e229.htm "HslCommunication.Profinet.Inovance")

[InovanceHelper 类](../html/944b43ff-5625-c33b-3105-1c1b8db259d9.htm "InovanceHelper 类")

[InovanceHelper 方法](../html/05a274bc-a7e4-39bb-7f78-0383d043cf99.htm "InovanceHelper 方法")

[PraseInovanceAddress 方法](../html/f3d46fa3-fc1f-bb51-0efb-c0331acb938a.htm "PraseInovanceAddress 方法 ")

[PraseInovanceAMAddress 方法](../html/9c7e153c-c315-4461-62dc-7ba419a0b6ac.htm "PraseInovanceAMAddress 方法 ")

[PraseInovanceH3UAddress 方法](../html/2a040cb9-febb-6a01-473e-e68edb5897fe.htm "PraseInovanceH3UAddress 方法 ")

[PraseInovanceH5UAddress 方法](../html/594cc79d-ad03-e922-4644-36674190c3d9.htm "PraseInovanceH5UAddress 方法 ")

[ReadByte 方法](../html/e20dbf83-6657-b7db-c5fb-4f8226e7dee7.htm "ReadByte 方法 ")

[ReadByteAsync 方法](../html/3afd1bcd-229d-1cd0-402b-26ff82e0c480.htm "ReadByteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| InovanceHelper 方法 |

[InovanceHelper](944b43ff-5625-c33b-3105-1c1b8db259d9.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法静态成员 | [PraseInovanceAddress](f3d46fa3-fc1f-bb51-0efb-c0331acb938a.htm) | 根据汇川PLC的地址，解析出转换后的modbus协议信息，适用AM,H3U,H5U系列的PLC  According to the address of Inovance PLC, analyze the converted modbus protocol information, which is suitable for AM, H3U, H5U series PLC |
| 公共方法静态成员 | [PraseInovanceAMAddress](9c7e153c-c315-4461-62dc-7ba419a0b6ac.htm) | 根据汇川PLC的地址，解析出转换后的modbus协议信息，适用AM,H3U,H5U系列的PLC  According to the address of Inovance PLC, analyze the converted modbus protocol information, which is suitable for AM, H3U, H5U series PLC |
| 公共方法静态成员 | [PraseInovanceH3UAddress](2a040cb9-febb-6a01-473e-e68edb5897fe.htm) | 根据汇川PLC的地址，解析出转换后的modbus协议信息，适用AM,H3U,H5U系列的PLC  According to the address of Inovance PLC, analyze the converted modbus protocol information, which is suitable for AM, H3U, H5U series PLC |
| 公共方法静态成员 | [PraseInovanceH5UAddress](594cc79d-ad03-e922-4644-36674190c3d9.htm) | 根据汇川PLC的地址，解析出转换后的modbus协议信息，适用AM,H3U,H5U系列的PLC  According to the address of Inovance PLC, analyze the converted modbus protocol information, which is suitable for AM, H3U, H5U series PLC |
| 公共方法静态成员 | [ReadByte](e20dbf83-6657-b7db-c5fb-4f8226e7dee7.htm) | 按照字节读取汇川M地址的数据，地址示例： MB100，MB101，需要注意的是，MB100 及 MB101 的地址是 MW50 的数据。  Read the data of Inovance M address according to the byte, address example: MB100, MB101, it should be noted that the addresses of MB100 and MB101 are the data of MW50. |
| 公共方法静态成员 | [ReadByteAsync](3afd1bcd-229d-1cd0-402b-26ff82e0c480.htm) | 按照字节读取汇川M地址的数据，地址示例： MB100，MB101，需要注意的是，MB100 及 MB101 的地址是 MW50 的数据。  Read the data of Inovance M address according to the byte, address example: MB100, MB101, it should be noted that the addresses of MB100 and MB101 are the data of MW50. |
| 公共方法 | ToString | (继承自 Object。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[InovanceHelper 类](944b43ff-5625-c33b-3105-1c1b8db259d9.htm)

[HslCommunication.Profinet.Inovance 命名空间](aa61d757-9239-26ce-525f-b7e52699e229.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## PraseInovanceAddress 方法 

[原文連結](http://api.hslcommunication.cn/html/f3d46fa3-fc1f-bb51-0efb-c0331acb938a.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Inovance](../html/aa61d757-9239-26ce-525f-b7e52699e229.htm "HslCommunication.Profinet.Inovance")

[InovanceHelper 类](../html/944b43ff-5625-c33b-3105-1c1b8db259d9.htm "InovanceHelper 类")

[InovanceHelper 方法](../html/05a274bc-a7e4-39bb-7f78-0383d043cf99.htm "InovanceHelper 方法")

[PraseInovanceAddress 方法](../html/f3d46fa3-fc1f-bb51-0efb-c0331acb938a.htm "PraseInovanceAddress 方法 ")

[PraseInovanceAMAddress 方法](../html/9c7e153c-c315-4461-62dc-7ba419a0b6ac.htm "PraseInovanceAMAddress 方法 ")

[PraseInovanceH3UAddress 方法](../html/2a040cb9-febb-6a01-473e-e68edb5897fe.htm "PraseInovanceH3UAddress 方法 ")

[PraseInovanceH5UAddress 方法](../html/594cc79d-ad03-e922-4644-36674190c3d9.htm "PraseInovanceH5UAddress 方法 ")

[ReadByte 方法](../html/e20dbf83-6657-b7db-c5fb-4f8226e7dee7.htm "ReadByte 方法 ")

[ReadByteAsync 方法](../html/3afd1bcd-229d-1cd0-402b-26ff82e0c480.htm "ReadByteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| InovanceHelperPraseInovanceAddress 方法 |

根据汇川PLC的地址，解析出转换后的modbus协议信息，适用AM,H3U,H5U系列的PLC  
According to the address of Inovance PLC, analyze the converted modbus protocol information, which is suitable for AM, H3U, H5U series PLC

**命名空间：**
 [HslCommunication.Profinet.Inovance](aa61d757-9239-26ce-525f-b7e52699e229.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult<string> PraseInovanceAddress(
	InovanceSeries series,
	string address,
	byte modbusCode
)
```

```
Public Shared Function PraseInovanceAddress ( 
	series As InovanceSeries,
	address As String,
	modbusCode As Byte
) As OperateResult(Of String)
```

```
public:
static OperateResult<String^>^ PraseInovanceAddress(
	InovanceSeries series, 
	String^ address, 
	unsigned char modbusCode
)
```

```
static member PraseInovanceAddress : 
        series : InovanceSeries * 
        address : string * 
        modbusCode : byte -> OperateResult<string> 
```

#### 参数

series
:   类型：[HslCommunication.Profinet.InovanceInovanceSeries](624b6c23-3fff-4ce9-aa3c-72b3897aec74.htm)  
    PLC的系列

address
:   类型：SystemString  
    汇川plc的地址信息

modbusCode
:   类型：SystemByte  
    原始的对应的modbus信息

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)String  
Modbus格式的地址

![](../icons/SectionExpanded.png)参见

#### 引用

[InovanceHelper 类](944b43ff-5625-c33b-3105-1c1b8db259d9.htm)

[HslCommunication.Profinet.Inovance 命名空间](aa61d757-9239-26ce-525f-b7e52699e229.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## PraseInovanceAMAddress 方法 

[原文連結](http://api.hslcommunication.cn/html/9c7e153c-c315-4461-62dc-7ba419a0b6ac.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Inovance](../html/aa61d757-9239-26ce-525f-b7e52699e229.htm "HslCommunication.Profinet.Inovance")

[InovanceHelper 类](../html/944b43ff-5625-c33b-3105-1c1b8db259d9.htm "InovanceHelper 类")

[InovanceHelper 方法](../html/05a274bc-a7e4-39bb-7f78-0383d043cf99.htm "InovanceHelper 方法")

[PraseInovanceAddress 方法](../html/f3d46fa3-fc1f-bb51-0efb-c0331acb938a.htm "PraseInovanceAddress 方法 ")

[PraseInovanceAMAddress 方法](../html/9c7e153c-c315-4461-62dc-7ba419a0b6ac.htm "PraseInovanceAMAddress 方法 ")

[PraseInovanceH3UAddress 方法](../html/2a040cb9-febb-6a01-473e-e68edb5897fe.htm "PraseInovanceH3UAddress 方法 ")

[PraseInovanceH5UAddress 方法](../html/594cc79d-ad03-e922-4644-36674190c3d9.htm "PraseInovanceH5UAddress 方法 ")

[ReadByte 方法](../html/e20dbf83-6657-b7db-c5fb-4f8226e7dee7.htm "ReadByte 方法 ")

[ReadByteAsync 方法](../html/3afd1bcd-229d-1cd0-402b-26ff82e0c480.htm "ReadByteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| InovanceHelperPraseInovanceAMAddress 方法 |

根据汇川PLC的地址，解析出转换后的modbus协议信息，适用AM,H3U,H5U系列的PLC  
According to the address of Inovance PLC, analyze the converted modbus protocol information, which is suitable for AM, H3U, H5U series PLC

**命名空间：**
 [HslCommunication.Profinet.Inovance](aa61d757-9239-26ce-525f-b7e52699e229.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult<string> PraseInovanceAMAddress(
	string address,
	byte modbusCode
)
```

```
Public Shared Function PraseInovanceAMAddress ( 
	address As String,
	modbusCode As Byte
) As OperateResult(Of String)
```

```
public:
static OperateResult<String^>^ PraseInovanceAMAddress(
	String^ address, 
	unsigned char modbusCode
)
```

```
static member PraseInovanceAMAddress : 
        address : string * 
        modbusCode : byte -> OperateResult<string> 
```

#### 参数

address
:   类型：SystemString  
    汇川plc的地址信息

modbusCode
:   类型：SystemByte  
    原始的对应的modbus信息

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)String  
Modbus格式的地址

![](../icons/SectionExpanded.png)参见

#### 引用

[InovanceHelper 类](944b43ff-5625-c33b-3105-1c1b8db259d9.htm)

[HslCommunication.Profinet.Inovance 命名空间](aa61d757-9239-26ce-525f-b7e52699e229.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## PraseInovanceH3UAddress 方法 

[原文連結](http://api.hslcommunication.cn/html/2a040cb9-febb-6a01-473e-e68edb5897fe.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Inovance](../html/aa61d757-9239-26ce-525f-b7e52699e229.htm "HslCommunication.Profinet.Inovance")

[InovanceHelper 类](../html/944b43ff-5625-c33b-3105-1c1b8db259d9.htm "InovanceHelper 类")

[InovanceHelper 方法](../html/05a274bc-a7e4-39bb-7f78-0383d043cf99.htm "InovanceHelper 方法")

[PraseInovanceAddress 方法](../html/f3d46fa3-fc1f-bb51-0efb-c0331acb938a.htm "PraseInovanceAddress 方法 ")

[PraseInovanceAMAddress 方法](../html/9c7e153c-c315-4461-62dc-7ba419a0b6ac.htm "PraseInovanceAMAddress 方法 ")

[PraseInovanceH3UAddress 方法](../html/2a040cb9-febb-6a01-473e-e68edb5897fe.htm "PraseInovanceH3UAddress 方法 ")

[PraseInovanceH5UAddress 方法](../html/594cc79d-ad03-e922-4644-36674190c3d9.htm "PraseInovanceH5UAddress 方法 ")

[ReadByte 方法](../html/e20dbf83-6657-b7db-c5fb-4f8226e7dee7.htm "ReadByte 方法 ")

[ReadByteAsync 方法](../html/3afd1bcd-229d-1cd0-402b-26ff82e0c480.htm "ReadByteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| InovanceHelperPraseInovanceH3UAddress 方法 |

根据汇川PLC的地址，解析出转换后的modbus协议信息，适用AM,H3U,H5U系列的PLC  
According to the address of Inovance PLC, analyze the converted modbus protocol information, which is suitable for AM, H3U, H5U series PLC

**命名空间：**
 [HslCommunication.Profinet.Inovance](aa61d757-9239-26ce-525f-b7e52699e229.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult<string> PraseInovanceH3UAddress(
	string address,
	byte modbusCode
)
```

```
Public Shared Function PraseInovanceH3UAddress ( 
	address As String,
	modbusCode As Byte
) As OperateResult(Of String)
```

```
public:
static OperateResult<String^>^ PraseInovanceH3UAddress(
	String^ address, 
	unsigned char modbusCode
)
```

```
static member PraseInovanceH3UAddress : 
        address : string * 
        modbusCode : byte -> OperateResult<string> 
```

#### 参数

address
:   类型：SystemString  
    汇川plc的地址信息

modbusCode
:   类型：SystemByte  
    原始的对应的modbus信息

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)String  
Modbus格式的地址

![](../icons/SectionExpanded.png)参见

#### 引用

[InovanceHelper 类](944b43ff-5625-c33b-3105-1c1b8db259d9.htm)

[HslCommunication.Profinet.Inovance 命名空间](aa61d757-9239-26ce-525f-b7e52699e229.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## PraseInovanceH5UAddress 方法 

[原文連結](http://api.hslcommunication.cn/html/594cc79d-ad03-e922-4644-36674190c3d9.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Inovance](../html/aa61d757-9239-26ce-525f-b7e52699e229.htm "HslCommunication.Profinet.Inovance")

[InovanceHelper 类](../html/944b43ff-5625-c33b-3105-1c1b8db259d9.htm "InovanceHelper 类")

[InovanceHelper 方法](../html/05a274bc-a7e4-39bb-7f78-0383d043cf99.htm "InovanceHelper 方法")

[PraseInovanceAddress 方法](../html/f3d46fa3-fc1f-bb51-0efb-c0331acb938a.htm "PraseInovanceAddress 方法 ")

[PraseInovanceAMAddress 方法](../html/9c7e153c-c315-4461-62dc-7ba419a0b6ac.htm "PraseInovanceAMAddress 方法 ")

[PraseInovanceH3UAddress 方法](../html/2a040cb9-febb-6a01-473e-e68edb5897fe.htm "PraseInovanceH3UAddress 方法 ")

[PraseInovanceH5UAddress 方法](../html/594cc79d-ad03-e922-4644-36674190c3d9.htm "PraseInovanceH5UAddress 方法 ")

[ReadByte 方法](../html/e20dbf83-6657-b7db-c5fb-4f8226e7dee7.htm "ReadByte 方法 ")

[ReadByteAsync 方法](../html/3afd1bcd-229d-1cd0-402b-26ff82e0c480.htm "ReadByteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| InovanceHelperPraseInovanceH5UAddress 方法 |

根据汇川PLC的地址，解析出转换后的modbus协议信息，适用AM,H3U,H5U系列的PLC  
According to the address of Inovance PLC, analyze the converted modbus protocol information, which is suitable for AM, H3U, H5U series PLC

**命名空间：**
 [HslCommunication.Profinet.Inovance](aa61d757-9239-26ce-525f-b7e52699e229.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult<string> PraseInovanceH5UAddress(
	string address,
	byte modbusCode
)
```

```
Public Shared Function PraseInovanceH5UAddress ( 
	address As String,
	modbusCode As Byte
) As OperateResult(Of String)
```

```
public:
static OperateResult<String^>^ PraseInovanceH5UAddress(
	String^ address, 
	unsigned char modbusCode
)
```

```
static member PraseInovanceH5UAddress : 
        address : string * 
        modbusCode : byte -> OperateResult<string> 
```

#### 参数

address
:   类型：SystemString  
    汇川plc的地址信息

modbusCode
:   类型：SystemByte  
    原始的对应的modbus信息

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)String  
Modbus格式的地址

![](../icons/SectionExpanded.png)参见

#### 引用

[InovanceHelper 类](944b43ff-5625-c33b-3105-1c1b8db259d9.htm)

[HslCommunication.Profinet.Inovance 命名空间](aa61d757-9239-26ce-525f-b7e52699e229.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadByte 方法 

[原文連結](http://api.hslcommunication.cn/html/e20dbf83-6657-b7db-c5fb-4f8226e7dee7.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Inovance](../html/aa61d757-9239-26ce-525f-b7e52699e229.htm "HslCommunication.Profinet.Inovance")

[InovanceHelper 类](../html/944b43ff-5625-c33b-3105-1c1b8db259d9.htm "InovanceHelper 类")

[InovanceHelper 方法](../html/05a274bc-a7e4-39bb-7f78-0383d043cf99.htm "InovanceHelper 方法")

[PraseInovanceAddress 方法](../html/f3d46fa3-fc1f-bb51-0efb-c0331acb938a.htm "PraseInovanceAddress 方法 ")

[PraseInovanceAMAddress 方法](../html/9c7e153c-c315-4461-62dc-7ba419a0b6ac.htm "PraseInovanceAMAddress 方法 ")

[PraseInovanceH3UAddress 方法](../html/2a040cb9-febb-6a01-473e-e68edb5897fe.htm "PraseInovanceH3UAddress 方法 ")

[PraseInovanceH5UAddress 方法](../html/594cc79d-ad03-e922-4644-36674190c3d9.htm "PraseInovanceH5UAddress 方法 ")

[ReadByte 方法](../html/e20dbf83-6657-b7db-c5fb-4f8226e7dee7.htm "ReadByte 方法 ")

[ReadByteAsync 方法](../html/3afd1bcd-229d-1cd0-402b-26ff82e0c480.htm "ReadByteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| InovanceHelperReadByte 方法 |

按照字节读取汇川M地址的数据，地址示例： MB100，MB101，需要注意的是，MB100 及 MB101 的地址是 MW50 的数据。  
Read the data of Inovance M address according to the byte, address example: MB100, MB101, it should be noted that the addresses of MB100 and MB101 are the data of MW50.

**命名空间：**
 [HslCommunication.Profinet.Inovance](aa61d757-9239-26ce-525f-b7e52699e229.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult<byte> ReadByte(
	IModbus modbus,
	string address
)
```

```
Public Shared Function ReadByte ( 
	modbus As IModbus,
	address As String
) As OperateResult(Of Byte)
```

```
public:
static OperateResult<unsigned char>^ ReadByte(
	IModbus^ modbus, 
	String^ address
)
```

```
static member ReadByte : 
        modbus : IModbus * 
        address : string -> OperateResult<byte> 
```

#### 参数

modbus
:   类型：[HslCommunication.ModBusIModbus](1d3ec4c8-6dd5-984d-e7ac-c107336985a5.htm)  
    汇川的PLC对象

address
:   类型：SystemString  
    地址信息

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
读取的结果数据

![](../icons/SectionExpanded.png)参见

#### 引用

[InovanceHelper 类](944b43ff-5625-c33b-3105-1c1b8db259d9.htm)

[HslCommunication.Profinet.Inovance 命名空间](aa61d757-9239-26ce-525f-b7e52699e229.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadByteAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/3afd1bcd-229d-1cd0-402b-26ff82e0c480.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Inovance](../html/aa61d757-9239-26ce-525f-b7e52699e229.htm "HslCommunication.Profinet.Inovance")

[InovanceHelper 类](../html/944b43ff-5625-c33b-3105-1c1b8db259d9.htm "InovanceHelper 类")

[InovanceHelper 方法](../html/05a274bc-a7e4-39bb-7f78-0383d043cf99.htm "InovanceHelper 方法")

[PraseInovanceAddress 方法](../html/f3d46fa3-fc1f-bb51-0efb-c0331acb938a.htm "PraseInovanceAddress 方法 ")

[PraseInovanceAMAddress 方法](../html/9c7e153c-c315-4461-62dc-7ba419a0b6ac.htm "PraseInovanceAMAddress 方法 ")

[PraseInovanceH3UAddress 方法](../html/2a040cb9-febb-6a01-473e-e68edb5897fe.htm "PraseInovanceH3UAddress 方法 ")

[PraseInovanceH5UAddress 方法](../html/594cc79d-ad03-e922-4644-36674190c3d9.htm "PraseInovanceH5UAddress 方法 ")

[ReadByte 方法](../html/e20dbf83-6657-b7db-c5fb-4f8226e7dee7.htm "ReadByte 方法 ")

[ReadByteAsync 方法](../html/3afd1bcd-229d-1cd0-402b-26ff82e0c480.htm "ReadByteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| InovanceHelperReadByteAsync 方法 |

按照字节读取汇川M地址的数据，地址示例： MB100，MB101，需要注意的是，MB100 及 MB101 的地址是 MW50 的数据。  
Read the data of Inovance M address according to the byte, address example: MB100, MB101, it should be noted that the addresses of MB100 and MB101 are the data of MW50.

**命名空间：**
 [HslCommunication.Profinet.Inovance](aa61d757-9239-26ce-525f-b7e52699e229.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static Task<OperateResult<byte>> ReadByteAsync(
	IModbus modbus,
	string address
)
```

```
Public Shared Function ReadByteAsync ( 
	modbus As IModbus,
	address As String
) As Task(Of OperateResult(Of Byte))
```

```
public:
static Task<OperateResult<unsigned char>^>^ ReadByteAsync(
	IModbus^ modbus, 
	String^ address
)
```

```
static member ReadByteAsync : 
        modbus : IModbus * 
        address : string -> Task<OperateResult<byte>> 
```

#### 参数

modbus
:   类型：[HslCommunication.ModBusIModbus](1d3ec4c8-6dd5-984d-e7ac-c107336985a5.htm)  
    汇川的PLC对象

address
:   类型：SystemString  
    地址信息

#### 返回值

类型：Task[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
读取的结果数据

![](../icons/SectionExpanded.png)参见

#### 引用

[InovanceHelper 类](944b43ff-5625-c33b-3105-1c1b8db259d9.htm)

[HslCommunication.Profinet.Inovance 命名空间](aa61d757-9239-26ce-525f-b7e52699e229.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## InovanceSerial 类

[原文連結](http://api.hslcommunication.cn/html/434659f2-7d84-27a2-72b5-aebd489024da.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Inovance](../html/aa61d757-9239-26ce-525f-b7e52699e229.htm "HslCommunication.Profinet.Inovance")

[InovanceSerial 类](../html/434659f2-7d84-27a2-72b5-aebd489024da.htm "InovanceSerial 类")

[InovanceSerial 构造函数](../html/9fd1afb8-916b-d937-2a93-6357b2a446b5.htm "InovanceSerial 构造函数 ")

[InovanceSerial 属性](../html/df3ac411-58e0-d3f4-5d96-4fc524001008.htm "InovanceSerial 属性")

[InovanceSerial 方法](../html/a202dad8-3af6-5794-c23f-1ddd2151e1a6.htm "InovanceSerial 方法")

[InovanceSerial 字段](../html/c8f8f731-c2b4-9995-4aaa-ba397ee9696c.htm "InovanceSerial 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| InovanceSerial 类 |

汇川的串口通信协议，A适用于AM400、 AM400\_800、 AC800、H3U, XP, H5U 等系列底层走的是MODBUS-RTU协议，地址说明参见标记  
Huichuan's serial communication protocol is applicable to AM400, AM400\_800, AC800 and other series. The bottom layer is MODBUS-RTU protocol. For the address description, please refer to the mark

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  [HslCommunication.Core.NetBinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)  
    [HslCommunication.Core.DeviceDeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)  
      [HslCommunication.Core.DeviceDeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)  
        [HslCommunication.ModBusModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)  
          HslCommunication.Profinet.InovanceInovanceSerial

**命名空间：**
 [HslCommunication.Profinet.Inovance](aa61d757-9239-26ce-525f-b7e52699e229.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class InovanceSerial : ModbusRtu
```

```
Public Class InovanceSerial
	Inherits ModbusRtu
```

```
public ref class InovanceSerial : public ModbusRtu
```

```
type InovanceSerial =  
    class
        inherit ModbusRtu
    end
```

InovanceSerial 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [InovanceSerial](c2f782e9-a32a-8fcc-3cb7-b504e792d9d4.htm) | 实例化一个默认的对象 |
| 公共方法 | [InovanceSerial(Byte)](487ba2f7-0980-8c26-ee8a-3041efecec52.htm) | 指定服务器地址，端口号，客户端自己的站号来初始化  Specify the server address, port number, and client's own station number to initialize |
| 公共方法 | [InovanceSerial(InovanceSeries, Byte)](3ab7f8c6-bd6a-1803-2f00-5ec1d60b4f3f.htm) | 指定服务器地址，端口号，客户端自己的站号来初始化  Specify the server address, port number, and client's own station number to initialize |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [AddressStartWithZero](4d17309a-509d-62c1-9119-3a5735c8315d.htm) | 获取或设置起始的地址是否从0开始，默认为True  Gets or sets whether the starting address starts from 0. The default is True (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共属性 | [BaudRate](5554f563-5e29-699c-435b-c8e7ad5fa0e4.htm) | 当前连接串口信息的波特率  Baud rate of current connection serial port information (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 公共属性 | [BroadcastStation](27196ecf-f9c7-30ac-d439-bf70d26dc328.htm) | 获取或是设置当前广播模式对应的站号，广播模式意味着不接收设备方的数据返回操作，默认为 -1，表示不使用广播模式。  Gets or sets the station number corresponding to the current broadcast mode. Broadcast mode means that the data return operation of the device is not received. The default value is -1, indicating that broadcast mode is not used. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共属性代码示例 | [ByteTransform](f35118c9-3691-f752-ae6f-2c5549714faa.htm) | 当前的数据变换机制，当你需要从字节数据转换类型数据的时候需要。  The current data transformation mechanism is required when you need to convert type data from byte data. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共属性 | [CommunicationPipe](a4fcb477-64d2-8da7-c712-687f0261755b.htm) | 获取或设置当前的管道信息，管道类型为[CommunicationPipe](ad8e903b-d289-21b0-077e-15ede8005221.htm)的继承类，内置了[PipeTcpNet](26b94bd5-2695-36e7-6378-c3d41846be8f.htm)管道，[PipeUdpNet](5a0b810e-187e-826b-0a13-34fd5576aa0e.htm)管道，[PipeSerialPort](9d924607-3dfc-3916-b532-c5dec48fe347.htm)管道等  Get or set the current pipeline information, the pipeline type is [CommunicationPipe](ad8e903b-d289-21b0-077e-15ede8005221.htm) inheritance class, [PipeTcpNet](26b94bd5-2695-36e7-6378-c3d41846be8f.htm) pipeline, [PipeUdpNet](5a0b810e-187e-826b-0a13-34fd5576aa0e.htm) pipeline, [PipeSerialPort](9d924607-3dfc-3916-b532-c5dec48fe347.htm) pipeline, etc (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 公共属性 | [ConnectionId](ce1efb4b-fbc9-4683-163e-0a8e71e99d6b.htm) | 当前连接的唯一ID号，默认为长度20的guid码加随机数组成，方便列表管理，也可以自己指定  The unique ID number of the current connection. The default is a 20-digit guid code plus a random number. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [Crc16CheckEnable](8b8b3de6-eaf7-0484-604f-657dc210ab61.htm) | 获取或设置是否启用CRC16校验码的检查功能，默认启用，如果需要忽略检查CRC16，则设置为 false 即可。  Gets or sets whether to enable the check function of CRC16 check code. It is enabled by default. If you need to ignore the check of CRC16, you can set it to false. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共属性 | [DataFormat](eeb434e5-c320-09f1-6cf4-ff4798c2505a.htm) | 获取或设置数据解析的格式，可选ABCD, BADC，CDAB，DCBA格式，对int,uint,float,double,long,ulong类型有作用  Get or set the format of the data analysis, optional ABCD, BADC, CDAB, DCBA format, effective for int, uint, float, double, long, ulong type (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共属性 | [DisableFunctionCode06](096f15a9-323c-69db-65a1-be3723d914d1.htm) | 获取或设置是否禁用功能码05的写入操作，设置为true时，表示不允许使用功能码 06 进行写入操作。写入单个的字数据也使用 0x10 功能码  Get or set whether to disable the write operation of function code 06. When set to true, it indicates that the write operation of function code 05 is not allowed. Writing individual character data also uses the 0x10 function code (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共属性 | [EnableWriteMaskCode](85372ceb-e693-1067-0285-581b163f8da4.htm) | 获取或设置当前掩码写入的功能码是否激活状态，设置为 false 时，再执行写入位时，会通过读字，修改位，写字的方式来间接实现。  When the function code is set to false, and then the write bit is executed, it will be indirectly implemented by reading, modifying the bit, and writing the word. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共属性 | [IsClearCacheBeforeRead](81dc971f-1f86-8535-ae9a-586b6590a584.htm) | 是否在发送数据前清空缓冲数据，默认是false  Whether to empty the buffer before sending data, the default is false (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 公共属性 | [IsStringReverse](c063a2e6-4c63-cb3f-8956-9164f3c853f0.htm) | 字符串数据是否按照字来反转，默认为False  Whether the string data is reversed according to words. The default is False. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共属性代码示例 | [LogNet](8c56f246-88ac-afb8-8c06-972fafe4f78a.htm) | 组件的日志工具，支持日志记录，只要实例化后，当前网络的基本信息，就以[DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm)等级进行输出  The component's logging tool supports logging. As long as the instantiation of the basic network information, the output will be output at [DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm) (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [PortName](f2b8e4e3-7079-5240-4ff1-14c6c9b376e6.htm) | 当前连接串口信息的端口号名称  The port name of the current connection serial port information (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 公共属性 | [ReceiveEmptyDataCount](f800fd56-7d49-87c4-985b-103a38330078.htm) | 获取或设置连续接收空的数据次数，在数据接收完成时有效，每个单位消耗的时间为[SleepTime](1e3a83ef-b6fc-bc8c-8950-ea3f6ae5c19b.htm)。  Obtain or set the number of consecutive times to receive empty data, which is valid when the data is received, and the time consumed by each unit is [SleepTime](1e3a83ef-b6fc-bc8c-8950-ea3f6ae5c19b.htm) (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 公共属性代码示例 | [ReceiveTimeOut](250560a7-332b-ce1b-5eb0-e3abe4ee6012.htm) | 获取或设置接收服务器反馈的时间，如果为负数，则不接收反馈   Gets or sets the time to receive server feedback, and if it is a negative number, does not receive feedback (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [RtsEnable](dd4057ef-88fe-f7d0-d2c4-f6e15488fe78.htm) | 获取或设置一个值，该值指示在串行通信中是否启用请求发送 (RTS) 信号。  Gets or sets a value indicating whether the request sending (RTS) signal is enabled in serial communication. (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 公共属性 | [SendBeforeHex](72aa1de4-116f-abfe-4adb-814a0be0e560.htm) | 获取或设置在发送通信报文前追加发送的字节信息，HEX格式，通常用于lora组网时，需要携带 00 00 00 02 四个字节的站地址功能。  Obtain or set the byte information sent before sending communication packets, HEX format, usually used for LORA networking, you need to carry 00 00 00 02 four-byte station address function. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [Series](e8a85bd7-58bd-bea4-db17-519d9f494168.htm) | 获取或设置汇川的系列，默认为AM系列 |
| 公共属性 | [SleepTime](86a5fbb9-3456-6beb-7113-5bbad5eabbc5.htm) | 获取或设置在正式接收对方返回数据前的时候，需要休息的时间，当设置为0的时候，不需要休息。  Get or set the time required to rest before officially receiving the data from the other party. When it is set to 0, no rest is required. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [Station](063d55ef-6207-33e0-9073-a725af469563.htm) | 获取或者重新修改服务器的默认站号信息，当然，你可以再读写的时候动态指定，例如地址 "s=2;100", 更详细的例子可以参考DEMO界面上的地址示例。  Get or re-modify the server's default station number information, of course, you can dynamically specify when reading and writing, such as the address "s=2; 100", a more detailed example can refer to the DEMO interface address example. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共属性 | [StationCheckMatch](03c10d53-c7a2-9256-df12-d780b8efd32c.htm) | 获取或设置当前是否启用站号检查的功能，默认启用，读写数据时将进行站号确认操作，如果需要忽略站号，则设置为 false 即可。  Gets or sets whether the station ID check function is enabled. It is enabled by default. When reading or writing data, the station ID is confirmed. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 受保护的属性 | [WordLength](f1a505cc-e85c-6c8b-b189-58a1ab5753f9.htm) | 一个字单位的数据表示的地址长度，西门子为2，三菱，欧姆龙，modbusTcp就为1，AB PLC无效  The address length represented by one word of data, Siemens is 2, Mitsubishi, Omron, modbusTcp is 1, AB PLC is invalid (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共属性 | [WordReadBatchLength](75150ff2-55eb-9881-6d35-d916a737a46d.htm) | 当进行字读取的时候，商业授权用户支持分批读取，本属性表示分批的长度信息，默认为 120，可以根据实际情况来设置其他值。  When reading of word type, commercial authorization users support batch reading. This attribute represents the length information of the batch, with a default of 120. Other values can be set according to the actual situation. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [Close](d3a2ebfb-4e6b-eeda-8bc8-c0972e21b5ff.htm) | 关闭当前的串口连接  Close the current serial connection (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 受保护的方法 | [DecideWhetherQAMessage](acc40cda-e310-8537-6685-59b67eeb16ff.htm) | 决定当前的消息是否是用于问答机制返回的消息，默认直接返回 true, 实际的情况需要根据协议进行重写方法  To determine whether the current message is the message returned by the question answering mechanism, the default is true. In actual cases, the rewriting method needs to be performed according to the protocol (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [Dispose](b9dda6bf-d342-254e-689a-93fa0bb265d1.htm) | (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 受保护的方法 | [Dispose(Boolean)](0bbb6ed3-111e-069c-e9f5-8cbc3f7f7b77.htm) | (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | [ExtraAfterReadFromCoreServer](06a49a15-a3ce-5e39-28c7-2b78851eca2e.htm) | 和服务器交互完成的时候调用的方法，可以根据读写结果进行一些额外的操作，具体的操作需要根据实际的需求来重写实现  The method called when the interaction with the server is completed can perform some additional operations based on the read and write results. The specific operations need to be rewritten according to actual needs. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法代码示例 | [ExtraOnDisconnect](de9dbba2-3b2c-1c5e-53c8-4e6f8a44e2c4.htm) | 根据实际的协议选择是否重写本方法，有些协议在断开连接之前，需要发送一些报文来关闭当前的网络通道  Select whether to rewrite this method according to the actual protocol. Some protocols need to send some packets to close the current network channel before disconnecting. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法代码示例 | [ExtraOnDisconnectAsync](210f593e-16d8-77c4-1178-58fbcfe9b27b.htm) | 根据实际的协议选择是否重写本方法，有些协议在断开连接之前，需要发送一些报文来关闭当前的网络通道  Select whether to rewrite this method according to the actual protocol. Some protocols need to send some packets to close the current network channel before disconnecting. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 受保护的方法 | [GetLogTextFromBinary](67f3b2f7-957c-fa0b-8320-6799237157db.htm) | 获取当前的报文进行日志记录的时候，是否使用二进制的格式记录，默认返回 [LogMsgFormatBinary](2ed2ff87-caec-0e44-f08b-e2b015225a58.htm)，重写可以根据session对象分别返回不同记录模式  Whether to log the current packet in binary format, the default return is [LogMsgFormatBinary](2ed2ff87-caec-0e44-f08b-e2b015225a58.htm). If you want to override it, different recording modes can be returned according to session (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [GetNewNetMessage](11f57be7-edd6-53dd-f91f-74c38a04578c.htm) | 获取一个新的消息对象的方法，需要在继承类里面进行重写  The method to get a new message object needs to be overridden in the inheritance class (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法 | [GetWordLength](c27821ea-1739-8eb9-a848-23a00828fec8.htm) | 一个字单位的数据表示的地址长度，西门子为2，三菱，欧姆龙，modbusTcp就为1，AB PLC无效  The address length represented by one word of data, Siemens is 2, Mitsubishi, Omron, modbusTcp is 1, AB PLC is invalid (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 受保护的方法代码示例 | [InitializationOnConnect](4ab6989f-a5b5-05fe-8b44-bfa45824a2a7.htm) | 根据实际的协议选择是否重写本方法，有些协议在创建连接之后，需要进行一些初始化的信号握手，才能最终建立网络通道。  Whether to rewrite this method is based on the actual protocol. Some protocols require some initial signal handshake to establish a network channel after the connection is created. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法代码示例 | [InitializationOnConnectAsync](7da2af04-d4fc-5045-19a4-4b0d1722eef3.htm) | 根据实际的协议选择是否重写本方法，有些协议在创建连接之后，需要进行一些初始化的信号握手，才能最终建立网络通道。  Whether to rewrite this method is based on the actual protocol. Some protocols require some initial signal handshake to establish a network channel after the connection is created. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [IsOpen](57cc9454-4b5a-37de-df7a-7184463ad5ef.htm) | 获取一个值，指示串口是否处于打开状态  Gets a value indicating whether the serial port is open (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 受保护的方法 | [LogRevcMessage(Byte)](dab1657f-1e75-7fe1-7003-ccedfc69855d.htm) | 使用日志记录一个接收的报文信息  Logs are used to record information about a received packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [LogRevcMessage(Byte, PipeSession)](0f51ae45-575e-4228-f978-01b37d6d5df1.htm) | 使用日志记录一个接收的报文信息  Logs are used to record information about a received packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [LogSendMessage(Byte)](4f90eecc-68a2-1ef2-44ef-369f19bed3ec.htm) | 使用日志记录一个发送的报文信息  Logs are used to record information about a send packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [LogSendMessage(Byte, PipeSession)](209caadf-4f0f-40aa-04c2-ebc61f8dab13.htm) | 使用日志记录一个发送的报文信息  Logs are used to record information about a send packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法 | [Open](e6840d02-f668-0c1a-d187-4d379de79908.htm) | 打开一个新的串行端口连接  Open a new serial port connection (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 公共方法 | [PackCommandWithHeader](ad370105-3cc4-77f9-797d-0ca5edafd761.htm) | 对当前的命令进行打包处理，通常是携带命令头内容，标记当前的命令的长度信息，需要进行重写，否则默认不打包  The current command is packaged, usually carrying the content of the command header, marking the length of the current command, and it needs to be rewritten, otherwise it is not packaged by default (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [Read(String, UInt16)](6571f657-29f0-1869-7e0c-4d4add79fc2b.htm) | 从Modbus服务器批量读取寄存器的信息，需要指定起始地址，读取长度，如果富文本地址不指定，默认使用的功能码是 0x03，如果需要使用04功能码，那么地址就写成 x=4;100  To read the register information from the Modbus server in batches, you need to specify the start address and read length. If the rich text address is not specified, the default function code is 0x03. If you need to use the 04 function code, the address is written as x = 4; 100 (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [ReadT](95abe064-5594-5ae7-a849-98c53f88fee0.htm) | 读取支持Hsl特性的数据内容，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明，或是官网的详细文档  Read the data content of the Hsl attribute. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadAsync(String, UInt16)](d8c14d37-56dc-2a31-1a19-3e2b5b65197d.htm) | 异步批量读取字节数组信息，需要指定地址和长度，返回原始的字节数组  Asynchronous batch read byte array information, need to specify the address and length, return the original byte array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadAsyncT](0522e740-4bf0-2e31-c4c5-e41384321b37.htm) | 异步读取支持Hsl特性的数据内容，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Asynchronously read the data content of the Hsl attribute. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadBool(String)](9b2faa2e-eaa9-7ace-c49f-7a0318f1df1a.htm) | 读取单个的Boolean数据信息  Read a single Boolean data message (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadBool(String, UInt16)](2d125640-18a5-f9a9-29b9-b121915c77ce.htm) | 批量读取线圈或是离散的数据信息，需要指定地址和长度，具体的结果取决于实现，如果富文本地址不指定，默认使用的功能码是 0x01  To read coils or discrete data in batches, you need to specify the address and length. The specific result depends on the implementation. If the rich text address is not specified, the default function code is 0x01. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [ReadBoolAsync(String)](cc37da1f-d673-fe67-78ee-4e465d29a2bd.htm) | 异步读取单个的Boolean数据信息  Asynchronously read a single Boolean data message (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadBoolAsync(String, UInt16)](a45d383d-0095-3ea7-32f8-e75af88d60f8.htm) | 异步批量读取Boolean数组信息，需要指定地址和长度，返回Boolean 数组  Asynchronously batch read Boolean array information, need to specify the address and length, return Boolean array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadByte](f2698a84-b0c3-6069-888e-4fcbb045db72.htm) | 按照字节读取汇川M地址的数据，地址示例： MB100，MB101，需要注意的是，MB100 及 MB101 的地址是 MW50 的数据。  Read the data of Inovance M address according to the byte, address example: MB100, MB101, it should be noted that the addresses of MB100 and MB101 are the data of MW50. |
| 公共方法 | [ReadByteAsync](70ab77e7-0972-e072-3e26-206a0fdf51e4.htm) | 按照字节读取汇川M地址的数据，地址示例： MB100，MB101，需要注意的是，MB100 及 MB101 的地址是 MW50 的数据。  Read the data of Inovance M address according to the byte, address example: MB100, MB101, it should be noted that the addresses of MB100 and MB101 are the data of MW50. |
| 公共方法 | [ReadCoil(String)](2d2adf84-ab2c-3b95-269f-3a8ddc659018.htm) | 读取线圈，需要指定起始地址，如果富文本地址不指定，默认使用的功能码是 0x01  To read the coil, you need to specify the start address. If the rich text address is not specified, the default function code is 0x01. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [ReadCoil(String, UInt16)](2869f250-557a-4888-5ebc-ea5e98028aed.htm) | 批量的读取线圈，需要指定起始地址，读取长度，如果富文本地址不指定，默认使用的功能码是 0x01  For batch reading coils, you need to specify the start address and read length. If the rich text address is not specified, the default function code is 0x01. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [ReadCoilAsync(String)](d5239965-a5cf-8e54-b4be-e4cb86b38cc0.htm) | 读取线圈，需要指定起始地址，如果富文本地址不指定，默认使用的功能码是 0x01  To read the coil, you need to specify the start address. If the rich text address is not specified, the default function code is 0x01. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [ReadCoilAsync(String, UInt16)](dfc142d3-3390-5eec-1128-6e19cac1ea12.htm) | 批量的读取线圈，需要指定起始地址，读取长度，如果富文本地址不指定，默认使用的功能码是 0x01  For batch reading coils, you need to specify the start address and read length. If the rich text address is not specified, the default function code is 0x01. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [ReadCustomerT(String)](ca069688-d8f0-7dbd-e962-86a10dce83cf.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，返回一个新的类型的实例对象。  To read a custom data type, you need to inherit from the IDataTransfer interface and return an instance object of a new type. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadCustomerT(String, T)](ca7dc1a1-c97d-689f-9819-4dae4af67206.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，传入一个实例，对这个实例进行赋值，并返回该实例的对象。  To read a custom data type, you need to inherit from the IDataTransfer interface, pass in an instance, assign a value to this instance, and return the object of the instance. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadCustomerAsyncT(String)](f3fcda2b-255d-aadd-567e-c7b3df3636d9.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，返回一个新的类型的实例对象。  To read a custom data type, you need to inherit from the IDataTransfer interface and return an instance object of a new type. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadCustomerAsyncT(String, T)](527c9fae-4743-6ebd-7b9a-5ef3a8b7c455.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，传入一个实例，对这个实例进行赋值，并返回该实例的对象。  To read a custom data type, you need to inherit from the IDataTransfer interface, pass in an instance, assign a value to this instance, and return the object of the instance. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadDiscrete(String)](cffe0847-c797-bc7f-71f1-f93092565ec6.htm) | 读取输入线圈，需要指定起始地址，如果富文本地址不指定，默认使用的功能码是 0x02  To read the input coil, you need to specify the start address. If the rich text address is not specified, the default function code is 0x02. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [ReadDiscrete(String, UInt16)](181cbf55-3a34-a988-9aa5-18432b9da97a.htm) | 批量的读取输入点，需要指定起始地址，读取长度，如果富文本地址不指定，默认使用的功能码是 0x02  To read input points in batches, you need to specify the start address and read length. If the rich text address is not specified, the default function code is 0x02 (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [ReadDiscreteAsync(String)](9d11737c-b5d9-5c24-09a1-0b70a87bdfc3.htm) | 读取输入线圈，需要指定起始地址，如果富文本地址不指定，默认使用的功能码是 0x02  To read the input coil, you need to specify the start address. If the rich text address is not specified, the default function code is 0x02. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [ReadDiscreteAsync(String, UInt16)](9c34359d-3d18-fdbf-d6c3-35134eec3b9d.htm) | 批量的读取输入点，需要指定起始地址，读取长度，如果富文本地址不指定，默认使用的功能码是 0x02  To read input points in batches, you need to specify the start address and read length. If the rich text address is not specified, the default function code is 0x02 (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [ReadDouble(String)](e7e31c43-71b4-e692-6bdc-ba9533e1d83d.htm) | 读取双浮点的数据  Read double floating point data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadDouble(String, UInt16)](9df438ac-8869-5002-7ace-bfb7c8f7aa6b.htm) | 读取双浮点数据的数组  Read double floating point data array (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [ReadDoubleAsync(String)](2cf895fa-1033-1def-f048-0e8caefaa71d.htm) | 异步读取双浮点的数据  Asynchronously read double floating point data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadDoubleAsync(String, UInt16)](e5c602ad-d3ce-a677-908b-03d72a4388cd.htm) | 异步读取双浮点数据的数组  Asynchronously read double floating point data array (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [ReadFile](702ddeff-3ba5-952d-abce-99a988a1d3cc.htm) | 使用0x14功能码读取文件的数据内容，需要指定文件编号，地址，长度，返回读取的结果数据，读取的长度为任意长度，内部自动根据实际情况来切割操作。 (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [ReadFloat(String)](bc2801ae-6600-e845-7c6c-522c9156eb0e.htm) | 读取单浮点数据  Read single floating point data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadFloat(String, UInt16)](909486dc-cc84-c34f-2de5-ec17b3d04182.htm) | 读取单浮点精度的数组  Read single floating point array (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [ReadFloatAsync(String)](89f743f8-4b3a-7085-1e49-a1d7adb10052.htm) | 异步读取单浮点数据  Asynchronously read single floating point data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadFloatAsync(String, UInt16)](eaeb9795-d78c-f884-fa54-c698c5e37404.htm) | 异步读取单浮点精度的数组  Asynchronously read single floating point array (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [ReadFromCoreServer(IEnumerableByte)](044cfdd9-ec76-b89f-2f2b-5f439025006c.htm) | 将多个数据报文按顺序发到设备，并从设备接收返回的数据内容，然后拼接成一个Byte[]信息，需要重写[UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm)方法才能返回正确的结果。  Send multiple data packets to the device in sequence, and receive the returned data content from the device, and then splicing them into a Byte[] message, you need to rewrite [UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm) method to return the correct result. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServer(Byte)](ff061b13-df60-6617-08c2-40b81cde38ea.htm) | 将Modbus报文数据发送到当前的通道中，并从通道中接收Modbus的报文，通道将根据当前连接自动获取，本方法是线程安全的。  Send Modbus message data to the current channel, and receive Modbus messages from the channel. The channel will automatically obtain it according to the current connection. This method is thread-safe. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServer(Byte, Boolean, Boolean)](b32ce8e8-a959-fcc2-6736-14b39fdb9a40.htm) | 将二进制的数据发送到管道中去，然后从管道里接收二进制的数据回来，并返回是否成功的结果对象。  Send binary data to the pipeline, and then receive binary data back from the pipeline, and return whether the success of the result object (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServer(CommunicationPipe, Byte, Boolean, Boolean)](ef36d56b-df09-f727-e8d3-316d37a13ee9.htm) | 使用指定的管道来进行数据通信，发送原始数据到管道，然后从管道接收相关的数据返回，本方法无锁 (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServerAsync(IEnumerableByte)](bcdb59b7-a70a-d5c0-1fc6-9ececcd39ea2.htm) | 将多个数据报文按顺序发到设备，并从设备接收返回的数据内容，然后拼接成一个Byte[]信息，需要重写[UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm)方法才能返回正确的结果。  Send multiple data packets to the device in sequence, and receive the returned data content from the device, and then splicing them into a Byte[] message, you need to rewrite [UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm) method to return the correct result. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServerAsync(Byte)](61427407-a616-84ae-62c4-abea8bc7a451.htm) | 将Modbus报文数据发送到当前的通道中，并从通道中接收Modbus的报文，通道将根据当前连接自动获取，本方法是线程安全的。  Send Modbus message data to the current channel, and receive Modbus messages from the channel. The channel will automatically obtain it according to the current connection. This method is thread-safe. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServerAsync(Byte, Boolean, Boolean)](1c79c204-5ab5-9342-e65d-7455c20ea577.htm) | 将二进制的数据发送到管道中去，然后从管道里接收二进制的数据回来，并返回是否成功的结果对象。  Send binary data to the pipeline, and then receive binary data back from the pipeline, and return whether the success of the result object (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServerAsync(CommunicationPipe, Byte, Boolean, Boolean)](a254a19d-3c29-147b-335b-8adccdf4ee2f.htm) | 使用指定的管道来进行数据通信，发送原始数据到管道，然后从管道接收相关的数据返回，本方法无锁 (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法代码示例 | [ReadInt16(String)](57a16003-2fd5-b2bd-09bf-14366be899b8.htm) | 读取16位的有符号的整型数据  Read 16-bit signed integer data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt16(String, UInt16)](b2bde158-f51c-93fc-66ca-5e4cf6ea1fb7.htm) | 读取16位的有符号整型数组  Read 16-bit signed integer array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt16Async(String)](ac021793-890d-ceab-eae4-47ff763390c1.htm) | 异步读取16位的有符号的整型数据  Asynchronously read 16-bit signed integer data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt16Async(String, UInt16)](f30e02fc-8078-f1b3-c226-5991e3e23caf.htm) | 异步读取16位的有符号整型数组  Asynchronously read 16-bit signed integer array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt32(String)](bfbea035-9a83-40e3-084c-28122041b37a.htm) | 读取32位的有符号整型  Read 32-bit signed integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt32(String, UInt16)](ca0b8024-b31c-4e45-ad3e-916e188efebe.htm) | 读取32位有符号整型数组  Read 32-bit signed integer array (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [ReadInt32Async(String)](96726ea0-82aa-1229-21f6-43cd62abca43.htm) | 异步读取32位的有符号整型  Asynchronously read 32-bit signed integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt32Async(String, UInt16)](4d8ec145-5fcb-0e5e-6065-957de3434321.htm) | 异步读取32位有符号整型数组  Asynchronously read 32-bit signed integer array (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [ReadInt64(String)](0d0c2284-2da1-6f9a-b353-537561cd5e95.htm) | 读取64位的有符号整型  Read 64-bit signed integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt64(String, UInt16)](662b75d9-9f79-a959-bd6f-a797741d84c8.htm) | 读取64位的有符号整型数组  Read 64-bit signed integer array (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [ReadInt64Async(String)](26d46c8a-3775-3abd-d3ec-5cd1758a317e.htm) | 异步读取64位的有符号整型  Asynchronously read 64-bit signed integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt64Async(String, UInt16)](c3fa5d8b-cf3e-2424-4541-99bb11f7b8d7.htm) | 异步读取64位的有符号整型数组  Asynchronously read 64-bit signed integer array (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [ReadString(String, UInt16)](07812d76-415a-6400-c8a2-1c21da51e7a3.htm) | 读取字符串数据，默认为最常见的ASCII编码  Read string data, default is the most common ASCII encoding (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadString(String, UInt16, Encoding)](62ed8198-3bc2-4374-b362-02b71138ab98.htm) | 使用指定的编码，读取字符串数据  Reads string data using the specified encoding (重写 [DeviceCommunicationReadString(String, UInt16, Encoding)](b34c71a6-58f3-6bf4-0b93-471be4f2ed31.htm).) |
| 公共方法代码示例 | [ReadStringAsync(String, UInt16)](1835da81-b811-d151-0b53-907d8d87d0a8.htm) | 异步读取字符串数据，默认为最常见的ASCII编码  Asynchronously read string data, default is the most common ASCII encoding (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadStringAsync(String, UInt16, Encoding)](8d084430-6e0d-a700-923a-ca7493cda482.htm) | 异步使用指定的编码，读取字符串数据  Asynchronously reads string data using the specified encoding (重写 [DeviceCommunicationReadStringAsync(String, UInt16, Encoding)](091dfc67-3d25-f33e-3a55-a0228863815e.htm).) |
| 公共方法代码示例 | [ReadStructT](49853082-fbad-52e5-d756-0615cedb4b83.htm) | 读取结构体类型的数据，根据结构体自身的定义，读取原始字节数组，然后解析出实际的结构体数据，结构体需要实现[HslStructAttribute](87e427a8-0c7b-43a8-c421-c27264eda90f.htm)特性 (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadStructAsyncT](36461e61-d0ed-62cf-79ce-dcc7f71a3fa2.htm) | 读取结构体类型的数据，根据结构体自身的定义，读取原始字节数组，然后解析出实际的结构体数据，结构体需要实现[HslStructAttribute](87e427a8-0c7b-43a8-c421-c27264eda90f.htm)特性 (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt16(String)](736379d1-9560-5052-658a-13c126d56d78.htm) | 读取16位的无符号整型  Read 16-bit unsigned integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt16(String, UInt16)](e284c4bd-d3f8-ff30-eb42-0aa035cfb4f6.htm) | 读取16位的无符号整型数组  Read 16-bit unsigned integer array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt16Async(String)](361acc93-43b5-a9a5-6fa7-fa8728b7b96c.htm) | 异步读取16位的无符号整型  Asynchronously read 16-bit unsigned integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt16Async(String, UInt16)](73dd24b8-b501-df53-1de6-ccb2708f656d.htm) | 异步读取16位的无符号整型数组  Asynchronously read 16-bit unsigned integer array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt32(String)](ca73c958-39c7-9b34-aa44-3db2d9b0cad3.htm) | 读取32位的无符号整型  Read 32-bit unsigned integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt32(String, UInt16)](7ef202c8-7f3a-596b-12c0-bf5e5a5abd92.htm) | 读取32位的无符号整型数组  Read 32-bit unsigned integer array (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [ReadUInt32Async(String)](2343953f-6350-331e-6ba7-18a4449f0575.htm) | 异步读取32位的无符号整型  Asynchronously read 32-bit unsigned integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt32Async(String, UInt16)](ca94295c-48ce-dbdf-01eb-4d01701a9f50.htm) | 异步读取32位的无符号整型数组  Asynchronously read 32-bit unsigned integer array (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [ReadUInt64(String)](947487dd-d4a3-72d8-d9fb-2224aad70e12.htm) | 读取64位的无符号整型  Read 64-bit unsigned integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt64(String, UInt16)](130d0d35-93a1-d72d-c2f4-53b525c11f02.htm) | 读取64位的无符号整型的数组  Read 64-bit unsigned integer array (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [ReadUInt64Async(String)](3411ff3a-772f-dbe1-8931-75e00a731579.htm) | 异步读取64位的无符号整型  Asynchronously read 64-bit unsigned integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt64Async(String, UInt16)](6f22e2db-a0eb-b4b7-f43b-030555667ade.htm) | 异步读取64位的无符号整型的数组  Asynchronously read 64-bit unsigned integer array (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [ReadWrite](455d5de4-9bc6-c0c0-db5c-b6a334800126.htm) | 使用0x17功能码来实现同时写入并读取数据的操作，使用一条报文来实现，需要指定读取的地址，长度，写入的地址，写入的数据信息，返回读取的结果数据。  Use 0x17 function code to write and read data at the same time, and use a message to implement it, you need to specify the read address, length, written address, written data information, and return the read result data. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [RegisteredAddressMapping](a00ba5f5-373c-6e60-845a-3d8ad093568e.htm) | 注册一个新的地址映射关系，注册地址映射关系后，就可以使用新的地址来读写Modbus数据了，通常用于其他的支持Modbus协议的PLC。  After registering a new address mapping, you can use the new address to read and write Modbus data, which is usually used for other PLCs that support the Modbus protocol. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [SerialPortInni(ActionSerialPort)](3728a0f2-4553-0bc8-5c92-253c13e84f05.htm) | 根据自定义初始化方法进行初始化串口信息  Initialize the serial port information according to the custom initialization method (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 公共方法 | [SerialPortInni(String)](f82e0ba4-af2f-5165-dcde-e047c3351e74.htm) | 初始化串口信息，9600波特率，8位数据位，1位停止位，无奇偶校验  Initial serial port information, 9600 baud rate, 8 data bits, 1 stop bit, no parity (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 公共方法 | [SerialPortInni(String, Int32)](433523a8-1606-f029-50db-68ef68034046.htm) | 初始化串口信息，波特率，8位数据位，1位停止位，无奇偶校验  Initializes serial port information, baud rate, 8-bit data bit, 1-bit stop bit, no parity (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 公共方法 | [SerialPortInni(String, Int32, Int32, StopBits, Parity)](cefa7cd2-9c1f-b79d-9ea3-d4d27d7ddeae.htm) | 初始化串口信息，波特率，数据位，停止位，奇偶校验需要全部自己来指定  Start serial port information, baud rate, data bit, stop bit, parity all need to be specified (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 公共方法 | [SetDtuPipe](f765d1d7-234a-b229-9457-97e25e024dc9.htm) | 将当前的通信对象设置DTU模式，允许传入现成的管道，并返回初始化结果，如果该设备重写了握手报文，就是返回握手结果  Set the current communication object to DTU mode, allow the existing pipe to be passed in, and return the initialization result, if the device rewrites the handshake packet, the handshake result is returned (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [SetDtuPipeAsync](9d1fd12d-b804-3828-a627-a2c0eef1d23f.htm) | 将当前的通信对象设置DTU模式，允许传入现成的管道，并返回初始化结果，如果该设备重写了握手报文，就是返回握手结果  Set the current communication object to DTU mode, allow the existing pipe to be passed in, and return the initialization result, if the device rewrites the handshake packet, the handshake result is returned (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ToString](a8ec9534-c499-74d1-b92e-bcb0edc12a82.htm) | (重写 [ModbusRtuToString](11f9cb41-d3e4-d45d-abd1-36e7bb2577e6.htm).) |
| 公共方法 | [TranslateToModbusAddress](e554c6c9-1c35-e7b1-5c1e-83bc6d4f4287.htm) | 将当前的地址信息转换成Modbus格式的地址，如果转换失败，返回失败的消息。默认不进行任何的转换。  Convert the current address information into a Modbus format address. If the conversion fails, a failure message will be returned. No conversion is performed by default. (重写 [ModbusRtuTranslateToModbusAddress(String, Byte)](43c76b4c-2f7e-c60a-8b90-ee2c3898ce42.htm).) |
| 公共方法 | [UnpackResponseContent](2af5ade0-337b-6bff-777e-99f624a94606.htm) | 根据对方返回的报文命令，对命令进行基本的拆包，例如各种Modbus协议拆包为统一的核心报文，还支持对报文的验证  According to the message command returned by the other party, the command is basically unpacked, for example, various Modbus protocols are unpacked into a unified core message, and the verification of the message is also supported (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
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
| 公共方法代码示例 | [Write(String, Double)](64c44576-7e95-cbcb-a8c9-e6b7ff7e3c50.htm) | 写入double数据，返回是否成功  Write double data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Int16)](c7cfd67d-ca31-3087-c684-ee4c10780ab1.htm) | 写入short数组，返回是否成功  Write short array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Int32)](ba80d1a5-c8d8-7e0b-5c6b-9435cf20cadb.htm) | 写入int数据，返回是否成功  Write int data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Int64)](ba246b33-5a12-a22c-6d40-0969a64b508b.htm) | 写入long数据，返回是否成功  Write long data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Single)](a2265fa7-5b56-a5aa-9d21-a0dc5458cbf9.htm) | 写入float数据，返回是否成功  Write float data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, String)](4e0ec363-ba64-7ca7-b8d3-2e9a8087bd55.htm) | 写入字符串信息，编码为ASCII  Write string information, encoded as ASCII (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, UInt16)](65544962-82ac-bb47-d72f-5dd0c3d4ac6c.htm) | 写入ushort数组，返回是否成功  Write ushort array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, UInt32)](c99d8246-f95e-4f79-738b-16e622072eb4.htm) | 写入uint数据，返回是否成功  Write uint data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, UInt64)](4924086e-9b36-d446-4ae4-6682134a6ad1.htm) | 写入ulong数据，返回是否成功  Write ulong data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [Write(String, Boolean)](a40c7e8d-5770-9ad2-e263-66a92c45f61c.htm) | 向线圈中写入bool数值，返回是否写入成功，如果富文本地址不指定，默认使用的功能码是 0x05， 如果你的地址为字地址，例如100.2，那么将使用0x16的功能码，通过掩码的方式来修改寄存器的某一位，需要Modbus服务器支持，否则自动切换为读取字数据，修改位，在写入字的方式。  Write bool value to the coil and return whether the writing is successful. If the rich text address is not specified, the default function code is 0x05. If your address is a word address, such as 100.2, then you will use the function code of 0x16 to modify a bit of the register through a mask. It needs Modbus server support, Otherwise, it automatically switches to read the word data, modifies the bits, and writes the word in the way. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [Write(String, Boolean)](fc664b2b-f304-683e-573c-239379b83eb0.htm) | 向线圈中写入bool数组，返回是否写入成功，如果富文本地址不指定，默认使用的功能码是 0x0F  Write the bool array to the coil, and return whether the writing is successful. If the rich text address is not specified, the default function code is 0x0F. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [Write(String, Byte)](ec533f87-fd4a-c833-0256-fdd83da3609f.htm) | 将数据写入到Modbus的寄存器上去，需要指定起始地址和数据内容，如果富文本地址不指定，默认使用的功能码是 0x10  To write data to Modbus registers, you need to specify the start address and data content. If the rich text address is not specified, the default function code is 0x10 (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [Write(String, Double)](e4621d28-0353-fd60-d600-f457a46a661e.htm) | 写入double数组，返回是否成功  Write double array, return whether the write was successful (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [Write(String, Int16)](87345acc-1644-e27d-f5b3-51a01e0d5b02.htm) | 将数据写入到Modbus的单个寄存器上去，需要指定起始地址和数据值，如果富文本地址不指定，默认使用的功能码是 0x06  To write data to a single register of Modbus, you need to specify the start address and data value. If the rich text address is not specified, the default function code is 0x06. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [Write(String, Int32)](cb0f528c-c39e-a439-6093-65c7c8e8ca10.htm) | 写入int[]数组，返回是否成功  Write int array, return whether the write was successful (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [Write(String, Int64)](d3322d92-64ca-0080-27c6-b8693206ad96.htm) | 写入long数组，返回是否成功  Write long array, return whether the write was successful (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [Write(String, Single)](08a2fa24-4d18-9e7d-1fc9-873e8870ecc2.htm) | 写入float数组，返回是否成功  Write float array, return whether the write was successful (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [Write(String, UInt16)](4c43fcf9-d5c0-e4e4-70c2-c2d1aa8bfe4b.htm) | 将数据写入到Modbus的单个寄存器上去，需要指定起始地址和数据值，如果富文本地址不指定，默认使用的功能码是 0x06  To write data to a single register of Modbus, you need to specify the start address and data value. If the rich text address is not specified, the default function code is 0x06. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [Write(String, UInt32)](3113766d-8db6-7f55-c9c9-0580b9b1b756.htm) | 写入uint[]数组，返回是否成功  Write uint array, return whether the write was successful (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [Write(String, UInt64)](414c37ec-7955-1077-addb-8a7c0a861215.htm) | 写入ulong数组，返回是否成功  Write ulong array, return whether the write was successful (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [Write(String, String, Int32)](02fe2ca9-49a3-ede4-d1af-8c6da1f644da.htm) | 写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为ASCII  Write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, it is filled with 0 and the encoding is ASCII. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, String, Encoding)](7b602661-902d-2b70-1fdb-ba87ecd307bf.htm) | 写入字符串信息，需要指定的编码信息  Write string information, need to specify the encoding information (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, String, Int32, Encoding)](e8cb30d9-4c5e-6c68-b9af-1b1b6c2dc714.htm) | 写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为指定的编码信息  Write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, then the operation is complemented with 0 , you should specified the encoding information (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteT(T)](9f48853f-fbeb-c37d-25a2-1c76d1dee407.htm) | 写入支持Hsl特性的数据，返回是否写入成功，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Write data that supports the Hsl attribute, and return whether the write was successful. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WriteAsync(String, Boolean)](fe872dae-7ef9-12c1-d883-e92ba8676416.htm) | 异步批量写入Boolean数组数据，返回是否成功  Asynchronously batch write Boolean array data, return success (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Byte)](8f1cec30-8b86-6749-c8f3-7d39168d59df.htm) | 异步写入原始的byte数组数据到指定的地址，返回是否写入成功  Asynchronously writes the original byte array data to the specified address, and returns whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Double)](068adaa4-d1de-dffc-8376-6f483e2c1d96.htm) | 异步写入double数据，返回是否成功  Asynchronously write double data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int16)](2625e868-7096-e526-ce25-b29afc639b41.htm) | 异步写入short数组，返回是否成功  Asynchronously write short array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int32)](df91d629-2080-328d-67cd-fbefb9a45404.htm) | 异步写入int数据，返回是否成功  Asynchronously write int data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int64)](ec6c1e24-a0d8-f1f5-56af-109b0c546fc8.htm) | 异步写入long数据，返回是否成功  Asynchronously write long data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Single)](910189df-13f9-b7d0-a626-b5b25ce5d32a.htm) | 异步写入float数据，返回是否成功  Asynchronously write float data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, String)](fe839c9a-46f3-e317-0155-b10a8a83a9b1.htm) | 异步写入字符串信息，编码为ASCII  Asynchronously write string information, encoded as ASCII (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt16)](34a1a321-972e-06a8-76fb-794dd1229005.htm) | 异步写入ushort数组，返回是否成功  Asynchronously write ushort array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt32)](abad215c-bc99-ac90-b9f3-7adb3a7cca8f.htm) | 异步写入uint数据，返回是否成功  Asynchronously write uint data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt64)](e396074d-6cf1-b45d-5e8e-01026f6e391b.htm) | 异步写入ulong数据，返回是否成功  Asynchronously write ulong data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WriteAsync(String, Boolean)](90b9711b-1891-383c-0e56-19dd224b7aab.htm) | 向线圈中写入bool数值，返回是否写入成功，如果富文本地址不指定，默认使用的功能码是 0x05， 如果你的地址为字地址，例如100.2，那么将使用0x16的功能码，通过掩码的方式来修改寄存器的某一位，需要Modbus服务器支持，否则自动切换为读取字数据，修改位，在写入字的方式。  Write bool value to the coil and return whether the writing is successful. If the rich text address is not specified, the default function code is 0x05. If your address is a word address, such as 100.2, then you will use the function code of 0x16 to modify a bit of the register through a mask. It needs Modbus server support, Otherwise, it automatically switches to read the word data, modifies the bits, and writes the word in the way. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Double)](d67e053b-fae0-fbbb-d688-b514f417bbd7.htm) | 异步写入double数组，返回是否成功  Asynchronously write double array, return whether the write was successful (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [WriteAsync(String, Int16)](7e819772-f315-0f56-8c1e-8e2bd3aa2396.htm) | 将数据写入到Modbus的单个寄存器上去，需要指定起始地址和数据值，如果富文本地址不指定，默认使用的功能码是 0x06  To write data to a single register of Modbus, you need to specify the start address and data value. If the rich text address is not specified, the default function code is 0x06. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int32)](9e2c12ba-21f7-1c91-71ae-1eb2f2f853c3.htm) | 异步写入int[]数组，返回是否成功  Asynchronously write int array, return whether the write was successful (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int64)](3ef3dd46-318a-f55d-7296-57d303aafc93.htm) | 异步写入long数组，返回是否成功  Asynchronously write long array, return whether the write was successful (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Single)](ffa2fa0f-8321-2c1f-3916-c2cd90882f3c.htm) | 异步写入float数组，返回是否成功  Asynchronously write float array, return whether the write was successful (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [WriteAsync(String, UInt16)](4bb55e75-5cb4-4b22-c8b7-2983eae63e31.htm) | 将数据写入到Modbus的单个寄存器上去，需要指定起始地址和数据值，如果富文本地址不指定，默认使用的功能码是 0x06  To write data to a single register of Modbus, you need to specify the start address and data value. If the rich text address is not specified, the default function code is 0x06. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt32)](e108f604-e11a-2599-b88c-d69e630212b5.htm) | 异步写入uint[]数组，返回是否成功  Asynchronously write uint array, return whether the write was successful (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt64)](d304bc30-22bb-24bc-dd75-f767cfba2c35.htm) | 异步写入ulong数组，返回是否成功  Asynchronously write ulong array, return whether the write was successful (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [WriteAsync(String, String, Int32)](b739b983-b03c-8ecf-2713-75f19233ae2f.htm) | 异步写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为ASCII  Asynchronously write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, it is filled with 0 and the encoding is ASCII. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, String, Encoding)](876e9800-0406-6407-bdf8-8ab1cf28df38.htm) | 异步写入字符串信息，需要指定的编码信息  Asynchronously write string information, need to specify the encoding information (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WriteAsync(String, String, Int32, Encoding)](ae9fb2a3-fddb-48eb-e9ad-6cc0ebbac1f6.htm) | 异步写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为指定的编码信息  Asynchronously write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, then the operation is complemented with 0 , you should specified the encoding information (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsyncT(T)](25d96e77-d1b2-918c-33f9-e5e4cc57c9cd.htm) | 异步写入支持Hsl特性的数据，返回是否写入成功，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Asynchronously write data that supports the Hsl attribute, and return whether the write was successful. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteCustomerT](84a37f82-3f67-a0f8-fd3f-0630772565e4.htm) | 写入自定义类型的数据，该类型必须继承自IDataTransfer接口  Write data of a custom type, which must inherit from the IDataTransfer interface (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteCustomerAsyncT](04628397-cfa1-d785-7790-4c02cb377ab6.htm) | 写入自定义类型的数据，该类型必须继承自IDataTransfer接口  Write data of a custom type, which must inherit from the IDataTransfer interface (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WriteFile](8b9201b8-f363-7860-6b09-e63811574239.htm) | 使用0x15功能码写入文件的数据内容，需要指定文件编号，地址，数据，返回写入的结果数据。写入的数据字节长度不能大于248个字节，如果想写入更多的字节，就需要切片写入操作 (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [WriteMask](6e5e9e07-cc0e-f60c-07be-eb4a63249863.htm) | 向设备写入掩码数据，使用0x16功能码，需要确认对方是否支持相关的操作，掩码数据的操作主要针对寄存器。  To write mask data to the server, using the 0x16 function code, you need to confirm whether the other party supports related operations. The operation of mask data is mainly directed to the register. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [WriteMaskAsync](fd2cff63-8ffc-bf53-a7b6-24253b445aff.htm) | 向设备写入掩码数据，使用0x16功能码，需要确认对方是否支持相关的操作，掩码数据的操作主要针对寄存器。  To write mask data to the server, using the 0x16 function code, you need to confirm whether the other party supports related operations. The operation of mask data is mainly directed to the register. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [WriteOneRegister(String, Int16)](ed003503-3828-ed79-d4b4-08cfc417c8ab.htm) | 将数据写入到Modbus的单个寄存器上去，需要指定起始地址和数据值，如果富文本地址不指定，默认使用的功能码是 0x06  To write data to a single register of Modbus, you need to specify the start address and data value. If the rich text address is not specified, the default function code is 0x06. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [WriteOneRegister(String, UInt16)](9603651d-3202-1560-9141-b2c33949e06c.htm) | 将数据写入到Modbus的单个寄存器上去，需要指定起始地址和数据值，如果富文本地址不指定，默认使用的功能码是 0x06  To write data to a single register of Modbus, you need to specify the start address and data value. If the rich text address is not specified, the default function code is 0x06. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [WriteOneRegisterAsync(String, Int16)](b5f960df-3437-910f-aba4-f4dcbbfedad8.htm) | 将数据写入到Modbus的单个寄存器上去，需要指定起始地址和数据值，如果富文本地址不指定，默认使用的功能码是 0x06  To write data to a single register of Modbus, you need to specify the start address and data value. If the rich text address is not specified, the default function code is 0x06. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [WriteOneRegisterAsync(String, UInt16)](66dad72a-c201-d696-fdf2-bcdac070fc57.htm) | 将数据写入到Modbus的单个寄存器上去，需要指定起始地址和数据值，如果富文本地址不指定，默认使用的功能码是 0x06  To write data to a single register of Modbus, you need to specify the start address and data value. If the rich text address is not specified, the default function code is 0x06. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |

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

AM400\_800 的元件有 Q 区，I 区，M 区这三种，分别都可以按位，按字节，按字和按双字进行访问，在本组件的条件下，仅支持按照位，字访问。  
位地址支持 Q, I, M 地址类型，字地址支持 SM, SD，支持对字地址的位访问，例如 ReadBool("SD0.5");
H3U 系列控制器支持 M/SM/S/T/C/X/Y 等 bit 型变量（也称线圈） 的访问、 D/SD/R/T/C 等 word 型变量的访问；  
H5U 系列控制器支持 M/B/S/X/Y 等 bit 型变量（也称线圈） 的访问、 D/R 等 word 型变量的访问；内部 W 元件，不支持通信访问。

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.Profinet.Inovance 命名空间](aa61d757-9239-26ce-525f-b7e52699e229.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## InovanceSerial 构造函数 

[原文連結](http://api.hslcommunication.cn/html/9fd1afb8-916b-d937-2a93-6357b2a446b5.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Inovance](../html/aa61d757-9239-26ce-525f-b7e52699e229.htm "HslCommunication.Profinet.Inovance")

[InovanceSerial 类](../html/434659f2-7d84-27a2-72b5-aebd489024da.htm "InovanceSerial 类")

[InovanceSerial 构造函数](../html/9fd1afb8-916b-d937-2a93-6357b2a446b5.htm "InovanceSerial 构造函数 ")

[InovanceSerial 构造函数](../html/c2f782e9-a32a-8fcc-3cb7-b504e792d9d4.htm "InovanceSerial 构造函数 ")

[InovanceSerial 构造函数 (Byte)](../html/487ba2f7-0980-8c26-ee8a-3041efecec52.htm "InovanceSerial 构造函数 (Byte)")

[InovanceSerial 构造函数 (InovanceSeries, Byte)](../html/3ab7f8c6-bd6a-1803-2f00-5ec1d60b4f3f.htm "InovanceSerial 构造函数 (InovanceSeries, Byte)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| InovanceSerial 构造函数 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [InovanceSerial](c2f782e9-a32a-8fcc-3cb7-b504e792d9d4.htm) | 实例化一个默认的对象 |
| 公共方法 | [InovanceSerial(Byte)](487ba2f7-0980-8c26-ee8a-3041efecec52.htm) | 指定服务器地址，端口号，客户端自己的站号来初始化  Specify the server address, port number, and client's own station number to initialize |
| 公共方法 | [InovanceSerial(InovanceSeries, Byte)](3ab7f8c6-bd6a-1803-2f00-5ec1d60b4f3f.htm) | 指定服务器地址，端口号，客户端自己的站号来初始化  Specify the server address, port number, and client's own station number to initialize |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[InovanceSerial 类](434659f2-7d84-27a2-72b5-aebd489024da.htm)

[HslCommunication.Profinet.Inovance 命名空间](aa61d757-9239-26ce-525f-b7e52699e229.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## InovanceSerial 构造函数 

[原文連結](http://api.hslcommunication.cn/html/c2f782e9-a32a-8fcc-3cb7-b504e792d9d4.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Inovance](../html/aa61d757-9239-26ce-525f-b7e52699e229.htm "HslCommunication.Profinet.Inovance")

[InovanceSerial 类](../html/434659f2-7d84-27a2-72b5-aebd489024da.htm "InovanceSerial 类")

[InovanceSerial 构造函数](../html/9fd1afb8-916b-d937-2a93-6357b2a446b5.htm "InovanceSerial 构造函数 ")

[InovanceSerial 构造函数](../html/c2f782e9-a32a-8fcc-3cb7-b504e792d9d4.htm "InovanceSerial 构造函数 ")

[InovanceSerial 构造函数 (Byte)](../html/487ba2f7-0980-8c26-ee8a-3041efecec52.htm "InovanceSerial 构造函数 (Byte)")

[InovanceSerial 构造函数 (InovanceSeries, Byte)](../html/3ab7f8c6-bd6a-1803-2f00-5ec1d60b4f3f.htm "InovanceSerial 构造函数 (InovanceSeries, Byte)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| InovanceSerial 构造函数 |

实例化一个默认的对象

**命名空间：**
 [HslCommunication.Profinet.Inovance](aa61d757-9239-26ce-525f-b7e52699e229.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public InovanceSerial()
```

```
Public Sub New
```

```
public:
InovanceSerial()
```

```
new : unit -> InovanceSerial
```

![](../icons/SectionExpanded.png)参见

#### 引用

[InovanceSerial 类](434659f2-7d84-27a2-72b5-aebd489024da.htm)

[InovanceSerial 重载](9fd1afb8-916b-d937-2a93-6357b2a446b5.htm)

[HslCommunication.Profinet.Inovance 命名空间](aa61d757-9239-26ce-525f-b7e52699e229.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## InovanceSerial 构造函数 (Byte)

[原文連結](http://api.hslcommunication.cn/html/487ba2f7-0980-8c26-ee8a-3041efecec52.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Inovance](../html/aa61d757-9239-26ce-525f-b7e52699e229.htm "HslCommunication.Profinet.Inovance")

[InovanceSerial 类](../html/434659f2-7d84-27a2-72b5-aebd489024da.htm "InovanceSerial 类")

[InovanceSerial 构造函数](../html/9fd1afb8-916b-d937-2a93-6357b2a446b5.htm "InovanceSerial 构造函数 ")

[InovanceSerial 构造函数](../html/c2f782e9-a32a-8fcc-3cb7-b504e792d9d4.htm "InovanceSerial 构造函数 ")

[InovanceSerial 构造函数 (Byte)](../html/487ba2f7-0980-8c26-ee8a-3041efecec52.htm "InovanceSerial 构造函数 (Byte)")

[InovanceSerial 构造函数 (InovanceSeries, Byte)](../html/3ab7f8c6-bd6a-1803-2f00-5ec1d60b4f3f.htm "InovanceSerial 构造函数 (InovanceSeries, Byte)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| InovanceSerial 构造函数 (Byte) |

指定服务器地址，端口号，客户端自己的站号来初始化  
Specify the server address, port number, and client's own station number to initialize

**命名空间：**
 [HslCommunication.Profinet.Inovance](aa61d757-9239-26ce-525f-b7e52699e229.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public InovanceSerial(
	byte station = 1
)
```

```
Public Sub New ( 
	Optional station As Byte = 1
)
```

```
public:
InovanceSerial(
	unsigned char station = 1
)
```

```
new : 
        ?station : byte 
(* Defaults:
        let _station = defaultArg station 1
*)
-> InovanceSerial
```

#### 参数

station (Optional)
:   类型：SystemByte  
    客户端自身的站号

![](../icons/SectionExpanded.png)参见

#### 引用

[InovanceSerial 类](434659f2-7d84-27a2-72b5-aebd489024da.htm)

[InovanceSerial 重载](9fd1afb8-916b-d937-2a93-6357b2a446b5.htm)

[HslCommunication.Profinet.Inovance 命名空间](aa61d757-9239-26ce-525f-b7e52699e229.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## InovanceSerial 构造函数 (InovanceSeries, Byte)

[原文連結](http://api.hslcommunication.cn/html/3ab7f8c6-bd6a-1803-2f00-5ec1d60b4f3f.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Inovance](../html/aa61d757-9239-26ce-525f-b7e52699e229.htm "HslCommunication.Profinet.Inovance")

[InovanceSerial 类](../html/434659f2-7d84-27a2-72b5-aebd489024da.htm "InovanceSerial 类")

[InovanceSerial 构造函数](../html/9fd1afb8-916b-d937-2a93-6357b2a446b5.htm "InovanceSerial 构造函数 ")

[InovanceSerial 构造函数](../html/c2f782e9-a32a-8fcc-3cb7-b504e792d9d4.htm "InovanceSerial 构造函数 ")

[InovanceSerial 构造函数 (Byte)](../html/487ba2f7-0980-8c26-ee8a-3041efecec52.htm "InovanceSerial 构造函数 (Byte)")

[InovanceSerial 构造函数 (InovanceSeries, Byte)](../html/3ab7f8c6-bd6a-1803-2f00-5ec1d60b4f3f.htm "InovanceSerial 构造函数 (InovanceSeries, Byte)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| InovanceSerial 构造函数 (InovanceSeries, Byte) |

指定服务器地址，端口号，客户端自己的站号来初始化  
Specify the server address, port number, and client's own station number to initialize

**命名空间：**
 [HslCommunication.Profinet.Inovance](aa61d757-9239-26ce-525f-b7e52699e229.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public InovanceSerial(
	InovanceSeries series,
	byte station = 1
)
```

```
Public Sub New ( 
	series As InovanceSeries,
	Optional station As Byte = 1
)
```

```
public:
InovanceSerial(
	InovanceSeries series, 
	unsigned char station = 1
)
```

```
new : 
        series : InovanceSeries * 
        ?station : byte 
(* Defaults:
        let _station = defaultArg station 1
*)
-> InovanceSerial
```

#### 参数

series
:   类型：[HslCommunication.Profinet.InovanceInovanceSeries](624b6c23-3fff-4ce9-aa3c-72b3897aec74.htm)  
    PLC的系列选择

station (Optional)
:   类型：SystemByte  
    客户端自身的站号

![](../icons/SectionExpanded.png)参见

#### 引用

[InovanceSerial 类](434659f2-7d84-27a2-72b5-aebd489024da.htm)

[InovanceSerial 重载](9fd1afb8-916b-d937-2a93-6357b2a446b5.htm)

[HslCommunication.Profinet.Inovance 命名空间](aa61d757-9239-26ce-525f-b7e52699e229.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## InovanceSerial 属性

[原文連結](http://api.hslcommunication.cn/html/df3ac411-58e0-d3f4-5d96-4fc524001008.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Inovance](../html/aa61d757-9239-26ce-525f-b7e52699e229.htm "HslCommunication.Profinet.Inovance")

[InovanceSerial 类](../html/434659f2-7d84-27a2-72b5-aebd489024da.htm "InovanceSerial 类")

[InovanceSerial 属性](../html/df3ac411-58e0-d3f4-5d96-4fc524001008.htm "InovanceSerial 属性")

[Series 属性](../html/e8a85bd7-58bd-bea4-db17-519d9f494168.htm "Series 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| InovanceSerial 属性 |

[InovanceSerial](434659f2-7d84-27a2-72b5-aebd489024da.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [AddressStartWithZero](4d17309a-509d-62c1-9119-3a5735c8315d.htm) | 获取或设置起始的地址是否从0开始，默认为True  Gets or sets whether the starting address starts from 0. The default is True (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共属性 | [BaudRate](5554f563-5e29-699c-435b-c8e7ad5fa0e4.htm) | 当前连接串口信息的波特率  Baud rate of current connection serial port information (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 公共属性 | [BroadcastStation](27196ecf-f9c7-30ac-d439-bf70d26dc328.htm) | 获取或是设置当前广播模式对应的站号，广播模式意味着不接收设备方的数据返回操作，默认为 -1，表示不使用广播模式。  Gets or sets the station number corresponding to the current broadcast mode. Broadcast mode means that the data return operation of the device is not received. The default value is -1, indicating that broadcast mode is not used. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共属性代码示例 | [ByteTransform](f35118c9-3691-f752-ae6f-2c5549714faa.htm) | 当前的数据变换机制，当你需要从字节数据转换类型数据的时候需要。  The current data transformation mechanism is required when you need to convert type data from byte data. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共属性 | [CommunicationPipe](a4fcb477-64d2-8da7-c712-687f0261755b.htm) | 获取或设置当前的管道信息，管道类型为[CommunicationPipe](ad8e903b-d289-21b0-077e-15ede8005221.htm)的继承类，内置了[PipeTcpNet](26b94bd5-2695-36e7-6378-c3d41846be8f.htm)管道，[PipeUdpNet](5a0b810e-187e-826b-0a13-34fd5576aa0e.htm)管道，[PipeSerialPort](9d924607-3dfc-3916-b532-c5dec48fe347.htm)管道等  Get or set the current pipeline information, the pipeline type is [CommunicationPipe](ad8e903b-d289-21b0-077e-15ede8005221.htm) inheritance class, [PipeTcpNet](26b94bd5-2695-36e7-6378-c3d41846be8f.htm) pipeline, [PipeUdpNet](5a0b810e-187e-826b-0a13-34fd5576aa0e.htm) pipeline, [PipeSerialPort](9d924607-3dfc-3916-b532-c5dec48fe347.htm) pipeline, etc (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 公共属性 | [ConnectionId](ce1efb4b-fbc9-4683-163e-0a8e71e99d6b.htm) | 当前连接的唯一ID号，默认为长度20的guid码加随机数组成，方便列表管理，也可以自己指定  The unique ID number of the current connection. The default is a 20-digit guid code plus a random number. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [Crc16CheckEnable](8b8b3de6-eaf7-0484-604f-657dc210ab61.htm) | 获取或设置是否启用CRC16校验码的检查功能，默认启用，如果需要忽略检查CRC16，则设置为 false 即可。  Gets or sets whether to enable the check function of CRC16 check code. It is enabled by default. If you need to ignore the check of CRC16, you can set it to false. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共属性 | [DataFormat](eeb434e5-c320-09f1-6cf4-ff4798c2505a.htm) | 获取或设置数据解析的格式，可选ABCD, BADC，CDAB，DCBA格式，对int,uint,float,double,long,ulong类型有作用  Get or set the format of the data analysis, optional ABCD, BADC, CDAB, DCBA format, effective for int, uint, float, double, long, ulong type (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共属性 | [DisableFunctionCode06](096f15a9-323c-69db-65a1-be3723d914d1.htm) | 获取或设置是否禁用功能码05的写入操作，设置为true时，表示不允许使用功能码 06 进行写入操作。写入单个的字数据也使用 0x10 功能码  Get or set whether to disable the write operation of function code 06. When set to true, it indicates that the write operation of function code 05 is not allowed. Writing individual character data also uses the 0x10 function code (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共属性 | [EnableWriteMaskCode](85372ceb-e693-1067-0285-581b163f8da4.htm) | 获取或设置当前掩码写入的功能码是否激活状态，设置为 false 时，再执行写入位时，会通过读字，修改位，写字的方式来间接实现。  When the function code is set to false, and then the write bit is executed, it will be indirectly implemented by reading, modifying the bit, and writing the word. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共属性 | [IsClearCacheBeforeRead](81dc971f-1f86-8535-ae9a-586b6590a584.htm) | 是否在发送数据前清空缓冲数据，默认是false  Whether to empty the buffer before sending data, the default is false (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 公共属性 | [IsStringReverse](c063a2e6-4c63-cb3f-8956-9164f3c853f0.htm) | 字符串数据是否按照字来反转，默认为False  Whether the string data is reversed according to words. The default is False. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共属性代码示例 | [LogNet](8c56f246-88ac-afb8-8c06-972fafe4f78a.htm) | 组件的日志工具，支持日志记录，只要实例化后，当前网络的基本信息，就以[DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm)等级进行输出  The component's logging tool supports logging. As long as the instantiation of the basic network information, the output will be output at [DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm) (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [PortName](f2b8e4e3-7079-5240-4ff1-14c6c9b376e6.htm) | 当前连接串口信息的端口号名称  The port name of the current connection serial port information (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 公共属性 | [ReceiveEmptyDataCount](f800fd56-7d49-87c4-985b-103a38330078.htm) | 获取或设置连续接收空的数据次数，在数据接收完成时有效，每个单位消耗的时间为[SleepTime](1e3a83ef-b6fc-bc8c-8950-ea3f6ae5c19b.htm)。  Obtain or set the number of consecutive times to receive empty data, which is valid when the data is received, and the time consumed by each unit is [SleepTime](1e3a83ef-b6fc-bc8c-8950-ea3f6ae5c19b.htm) (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 公共属性代码示例 | [ReceiveTimeOut](250560a7-332b-ce1b-5eb0-e3abe4ee6012.htm) | 获取或设置接收服务器反馈的时间，如果为负数，则不接收反馈   Gets or sets the time to receive server feedback, and if it is a negative number, does not receive feedback (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [RtsEnable](dd4057ef-88fe-f7d0-d2c4-f6e15488fe78.htm) | 获取或设置一个值，该值指示在串行通信中是否启用请求发送 (RTS) 信号。  Gets or sets a value indicating whether the request sending (RTS) signal is enabled in serial communication. (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 公共属性 | [SendBeforeHex](72aa1de4-116f-abfe-4adb-814a0be0e560.htm) | 获取或设置在发送通信报文前追加发送的字节信息，HEX格式，通常用于lora组网时，需要携带 00 00 00 02 四个字节的站地址功能。  Obtain or set the byte information sent before sending communication packets, HEX format, usually used for LORA networking, you need to carry 00 00 00 02 four-byte station address function. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [Series](e8a85bd7-58bd-bea4-db17-519d9f494168.htm) | 获取或设置汇川的系列，默认为AM系列 |
| 公共属性 | [SleepTime](86a5fbb9-3456-6beb-7113-5bbad5eabbc5.htm) | 获取或设置在正式接收对方返回数据前的时候，需要休息的时间，当设置为0的时候，不需要休息。  Get or set the time required to rest before officially receiving the data from the other party. When it is set to 0, no rest is required. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [Station](063d55ef-6207-33e0-9073-a725af469563.htm) | 获取或者重新修改服务器的默认站号信息，当然，你可以再读写的时候动态指定，例如地址 "s=2;100", 更详细的例子可以参考DEMO界面上的地址示例。  Get or re-modify the server's default station number information, of course, you can dynamically specify when reading and writing, such as the address "s=2; 100", a more detailed example can refer to the DEMO interface address example. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共属性 | [StationCheckMatch](03c10d53-c7a2-9256-df12-d780b8efd32c.htm) | 获取或设置当前是否启用站号检查的功能，默认启用，读写数据时将进行站号确认操作，如果需要忽略站号，则设置为 false 即可。  Gets or sets whether the station ID check function is enabled. It is enabled by default. When reading or writing data, the station ID is confirmed. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 受保护的属性 | [WordLength](f1a505cc-e85c-6c8b-b189-58a1ab5753f9.htm) | 一个字单位的数据表示的地址长度，西门子为2，三菱，欧姆龙，modbusTcp就为1，AB PLC无效  The address length represented by one word of data, Siemens is 2, Mitsubishi, Omron, modbusTcp is 1, AB PLC is invalid (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共属性 | [WordReadBatchLength](75150ff2-55eb-9881-6d35-d916a737a46d.htm) | 当进行字读取的时候，商业授权用户支持分批读取，本属性表示分批的长度信息，默认为 120，可以根据实际情况来设置其他值。  When reading of word type, commercial authorization users support batch reading. This attribute represents the length information of the batch, with a default of 120. Other values can be set according to the actual situation. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[InovanceSerial 类](434659f2-7d84-27a2-72b5-aebd489024da.htm)

[HslCommunication.Profinet.Inovance 命名空间](aa61d757-9239-26ce-525f-b7e52699e229.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Series 属性 

[原文連結](http://api.hslcommunication.cn/html/e8a85bd7-58bd-bea4-db17-519d9f494168.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Inovance](../html/aa61d757-9239-26ce-525f-b7e52699e229.htm "HslCommunication.Profinet.Inovance")

[InovanceSerial 类](../html/434659f2-7d84-27a2-72b5-aebd489024da.htm "InovanceSerial 类")

[InovanceSerial 属性](../html/df3ac411-58e0-d3f4-5d96-4fc524001008.htm "InovanceSerial 属性")

[Series 属性](../html/e8a85bd7-58bd-bea4-db17-519d9f494168.htm "Series 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| InovanceSerialSeries 属性 |

获取或设置汇川的系列，默认为AM系列

**命名空间：**
 [HslCommunication.Profinet.Inovance](aa61d757-9239-26ce-525f-b7e52699e229.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public InovanceSeries Series { get; set; }
```

```
Public Property Series As InovanceSeries
	Get
	Set
```

```
public:
property InovanceSeries Series {
	InovanceSeries get ();
	void set (InovanceSeries value);
}
```

```
member Series : InovanceSeries with get, set
```

#### 属性值

类型：[InovanceSeries](624b6c23-3fff-4ce9-aa3c-72b3897aec74.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[InovanceSerial 类](434659f2-7d84-27a2-72b5-aebd489024da.htm)

[HslCommunication.Profinet.Inovance 命名空间](aa61d757-9239-26ce-525f-b7e52699e229.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## InovanceSerial 方法

[原文連結](http://api.hslcommunication.cn/html/a202dad8-3af6-5794-c23f-1ddd2151e1a6.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Inovance](../html/aa61d757-9239-26ce-525f-b7e52699e229.htm "HslCommunication.Profinet.Inovance")

[InovanceSerial 类](../html/434659f2-7d84-27a2-72b5-aebd489024da.htm "InovanceSerial 类")

[InovanceSerial 方法](../html/a202dad8-3af6-5794-c23f-1ddd2151e1a6.htm "InovanceSerial 方法")

[ReadByte 方法](../html/f2698a84-b0c3-6069-888e-4fcbb045db72.htm "ReadByte 方法 ")

[ReadByteAsync 方法](../html/70ab77e7-0972-e072-3e26-206a0fdf51e4.htm "ReadByteAsync 方法 ")

[ReadString 方法](../html/47ff96ed-6d9a-2110-e28c-8a4c529ef384.htm "ReadString 方法 ")

[ReadStringAsync 方法](../html/e62e86b6-223b-20b8-0c2c-b28af24cf3fe.htm "ReadStringAsync 方法 ")

[ToString 方法](../html/a8ec9534-c499-74d1-b92e-bcb0edc12a82.htm "ToString 方法 ")

[TranslateToModbusAddress 方法](../html/e554c6c9-1c35-e7b1-5c1e-83bc6d4f4287.htm "TranslateToModbusAddress 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| InovanceSerial 方法 |

[InovanceSerial](434659f2-7d84-27a2-72b5-aebd489024da.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [Close](d3a2ebfb-4e6b-eeda-8bc8-c0972e21b5ff.htm) | 关闭当前的串口连接  Close the current serial connection (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 受保护的方法 | [DecideWhetherQAMessage](acc40cda-e310-8537-6685-59b67eeb16ff.htm) | 决定当前的消息是否是用于问答机制返回的消息，默认直接返回 true, 实际的情况需要根据协议进行重写方法  To determine whether the current message is the message returned by the question answering mechanism, the default is true. In actual cases, the rewriting method needs to be performed according to the protocol (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [Dispose](b9dda6bf-d342-254e-689a-93fa0bb265d1.htm) | (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 受保护的方法 | [Dispose(Boolean)](0bbb6ed3-111e-069c-e9f5-8cbc3f7f7b77.htm) | (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | [ExtraAfterReadFromCoreServer](06a49a15-a3ce-5e39-28c7-2b78851eca2e.htm) | 和服务器交互完成的时候调用的方法，可以根据读写结果进行一些额外的操作，具体的操作需要根据实际的需求来重写实现  The method called when the interaction with the server is completed can perform some additional operations based on the read and write results. The specific operations need to be rewritten according to actual needs. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法代码示例 | [ExtraOnDisconnect](de9dbba2-3b2c-1c5e-53c8-4e6f8a44e2c4.htm) | 根据实际的协议选择是否重写本方法，有些协议在断开连接之前，需要发送一些报文来关闭当前的网络通道  Select whether to rewrite this method according to the actual protocol. Some protocols need to send some packets to close the current network channel before disconnecting. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法代码示例 | [ExtraOnDisconnectAsync](210f593e-16d8-77c4-1178-58fbcfe9b27b.htm) | 根据实际的协议选择是否重写本方法，有些协议在断开连接之前，需要发送一些报文来关闭当前的网络通道  Select whether to rewrite this method according to the actual protocol. Some protocols need to send some packets to close the current network channel before disconnecting. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 受保护的方法 | [GetLogTextFromBinary](67f3b2f7-957c-fa0b-8320-6799237157db.htm) | 获取当前的报文进行日志记录的时候，是否使用二进制的格式记录，默认返回 [LogMsgFormatBinary](2ed2ff87-caec-0e44-f08b-e2b015225a58.htm)，重写可以根据session对象分别返回不同记录模式  Whether to log the current packet in binary format, the default return is [LogMsgFormatBinary](2ed2ff87-caec-0e44-f08b-e2b015225a58.htm). If you want to override it, different recording modes can be returned according to session (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [GetNewNetMessage](11f57be7-edd6-53dd-f91f-74c38a04578c.htm) | 获取一个新的消息对象的方法，需要在继承类里面进行重写  The method to get a new message object needs to be overridden in the inheritance class (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法 | [GetWordLength](c27821ea-1739-8eb9-a848-23a00828fec8.htm) | 一个字单位的数据表示的地址长度，西门子为2，三菱，欧姆龙，modbusTcp就为1，AB PLC无效  The address length represented by one word of data, Siemens is 2, Mitsubishi, Omron, modbusTcp is 1, AB PLC is invalid (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 受保护的方法代码示例 | [InitializationOnConnect](4ab6989f-a5b5-05fe-8b44-bfa45824a2a7.htm) | 根据实际的协议选择是否重写本方法，有些协议在创建连接之后，需要进行一些初始化的信号握手，才能最终建立网络通道。  Whether to rewrite this method is based on the actual protocol. Some protocols require some initial signal handshake to establish a network channel after the connection is created. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法代码示例 | [InitializationOnConnectAsync](7da2af04-d4fc-5045-19a4-4b0d1722eef3.htm) | 根据实际的协议选择是否重写本方法，有些协议在创建连接之后，需要进行一些初始化的信号握手，才能最终建立网络通道。  Whether to rewrite this method is based on the actual protocol. Some protocols require some initial signal handshake to establish a network channel after the connection is created. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [IsOpen](57cc9454-4b5a-37de-df7a-7184463ad5ef.htm) | 获取一个值，指示串口是否处于打开状态  Gets a value indicating whether the serial port is open (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 受保护的方法 | [LogRevcMessage(Byte)](dab1657f-1e75-7fe1-7003-ccedfc69855d.htm) | 使用日志记录一个接收的报文信息  Logs are used to record information about a received packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [LogRevcMessage(Byte, PipeSession)](0f51ae45-575e-4228-f978-01b37d6d5df1.htm) | 使用日志记录一个接收的报文信息  Logs are used to record information about a received packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [LogSendMessage(Byte)](4f90eecc-68a2-1ef2-44ef-369f19bed3ec.htm) | 使用日志记录一个发送的报文信息  Logs are used to record information about a send packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [LogSendMessage(Byte, PipeSession)](209caadf-4f0f-40aa-04c2-ebc61f8dab13.htm) | 使用日志记录一个发送的报文信息  Logs are used to record information about a send packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法 | [Open](e6840d02-f668-0c1a-d187-4d379de79908.htm) | 打开一个新的串行端口连接  Open a new serial port connection (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 公共方法 | [PackCommandWithHeader](ad370105-3cc4-77f9-797d-0ca5edafd761.htm) | 对当前的命令进行打包处理，通常是携带命令头内容，标记当前的命令的长度信息，需要进行重写，否则默认不打包  The current command is packaged, usually carrying the content of the command header, marking the length of the current command, and it needs to be rewritten, otherwise it is not packaged by default (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [Read(String, UInt16)](6571f657-29f0-1869-7e0c-4d4add79fc2b.htm) | 从Modbus服务器批量读取寄存器的信息，需要指定起始地址，读取长度，如果富文本地址不指定，默认使用的功能码是 0x03，如果需要使用04功能码，那么地址就写成 x=4;100  To read the register information from the Modbus server in batches, you need to specify the start address and read length. If the rich text address is not specified, the default function code is 0x03. If you need to use the 04 function code, the address is written as x = 4; 100 (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [ReadT](95abe064-5594-5ae7-a849-98c53f88fee0.htm) | 读取支持Hsl特性的数据内容，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明，或是官网的详细文档  Read the data content of the Hsl attribute. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadAsync(String, UInt16)](d8c14d37-56dc-2a31-1a19-3e2b5b65197d.htm) | 异步批量读取字节数组信息，需要指定地址和长度，返回原始的字节数组  Asynchronous batch read byte array information, need to specify the address and length, return the original byte array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadAsyncT](0522e740-4bf0-2e31-c4c5-e41384321b37.htm) | 异步读取支持Hsl特性的数据内容，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Asynchronously read the data content of the Hsl attribute. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadBool(String)](9b2faa2e-eaa9-7ace-c49f-7a0318f1df1a.htm) | 读取单个的Boolean数据信息  Read a single Boolean data message (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadBool(String, UInt16)](2d125640-18a5-f9a9-29b9-b121915c77ce.htm) | 批量读取线圈或是离散的数据信息，需要指定地址和长度，具体的结果取决于实现，如果富文本地址不指定，默认使用的功能码是 0x01  To read coils or discrete data in batches, you need to specify the address and length. The specific result depends on the implementation. If the rich text address is not specified, the default function code is 0x01. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [ReadBoolAsync(String)](cc37da1f-d673-fe67-78ee-4e465d29a2bd.htm) | 异步读取单个的Boolean数据信息  Asynchronously read a single Boolean data message (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadBoolAsync(String, UInt16)](a45d383d-0095-3ea7-32f8-e75af88d60f8.htm) | 异步批量读取Boolean数组信息，需要指定地址和长度，返回Boolean 数组  Asynchronously batch read Boolean array information, need to specify the address and length, return Boolean array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadByte](f2698a84-b0c3-6069-888e-4fcbb045db72.htm) | 按照字节读取汇川M地址的数据，地址示例： MB100，MB101，需要注意的是，MB100 及 MB101 的地址是 MW50 的数据。  Read the data of Inovance M address according to the byte, address example: MB100, MB101, it should be noted that the addresses of MB100 and MB101 are the data of MW50. |
| 公共方法 | [ReadByteAsync](70ab77e7-0972-e072-3e26-206a0fdf51e4.htm) | 按照字节读取汇川M地址的数据，地址示例： MB100，MB101，需要注意的是，MB100 及 MB101 的地址是 MW50 的数据。  Read the data of Inovance M address according to the byte, address example: MB100, MB101, it should be noted that the addresses of MB100 and MB101 are the data of MW50. |
| 公共方法 | [ReadCoil(String)](2d2adf84-ab2c-3b95-269f-3a8ddc659018.htm) | 读取线圈，需要指定起始地址，如果富文本地址不指定，默认使用的功能码是 0x01  To read the coil, you need to specify the start address. If the rich text address is not specified, the default function code is 0x01. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [ReadCoil(String, UInt16)](2869f250-557a-4888-5ebc-ea5e98028aed.htm) | 批量的读取线圈，需要指定起始地址，读取长度，如果富文本地址不指定，默认使用的功能码是 0x01  For batch reading coils, you need to specify the start address and read length. If the rich text address is not specified, the default function code is 0x01. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [ReadCoilAsync(String)](d5239965-a5cf-8e54-b4be-e4cb86b38cc0.htm) | 读取线圈，需要指定起始地址，如果富文本地址不指定，默认使用的功能码是 0x01  To read the coil, you need to specify the start address. If the rich text address is not specified, the default function code is 0x01. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [ReadCoilAsync(String, UInt16)](dfc142d3-3390-5eec-1128-6e19cac1ea12.htm) | 批量的读取线圈，需要指定起始地址，读取长度，如果富文本地址不指定，默认使用的功能码是 0x01  For batch reading coils, you need to specify the start address and read length. If the rich text address is not specified, the default function code is 0x01. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [ReadCustomerT(String)](ca069688-d8f0-7dbd-e962-86a10dce83cf.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，返回一个新的类型的实例对象。  To read a custom data type, you need to inherit from the IDataTransfer interface and return an instance object of a new type. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadCustomerT(String, T)](ca7dc1a1-c97d-689f-9819-4dae4af67206.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，传入一个实例，对这个实例进行赋值，并返回该实例的对象。  To read a custom data type, you need to inherit from the IDataTransfer interface, pass in an instance, assign a value to this instance, and return the object of the instance. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadCustomerAsyncT(String)](f3fcda2b-255d-aadd-567e-c7b3df3636d9.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，返回一个新的类型的实例对象。  To read a custom data type, you need to inherit from the IDataTransfer interface and return an instance object of a new type. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadCustomerAsyncT(String, T)](527c9fae-4743-6ebd-7b9a-5ef3a8b7c455.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，传入一个实例，对这个实例进行赋值，并返回该实例的对象。  To read a custom data type, you need to inherit from the IDataTransfer interface, pass in an instance, assign a value to this instance, and return the object of the instance. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadDiscrete(String)](cffe0847-c797-bc7f-71f1-f93092565ec6.htm) | 读取输入线圈，需要指定起始地址，如果富文本地址不指定，默认使用的功能码是 0x02  To read the input coil, you need to specify the start address. If the rich text address is not specified, the default function code is 0x02. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [ReadDiscrete(String, UInt16)](181cbf55-3a34-a988-9aa5-18432b9da97a.htm) | 批量的读取输入点，需要指定起始地址，读取长度，如果富文本地址不指定，默认使用的功能码是 0x02  To read input points in batches, you need to specify the start address and read length. If the rich text address is not specified, the default function code is 0x02 (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [ReadDiscreteAsync(String)](9d11737c-b5d9-5c24-09a1-0b70a87bdfc3.htm) | 读取输入线圈，需要指定起始地址，如果富文本地址不指定，默认使用的功能码是 0x02  To read the input coil, you need to specify the start address. If the rich text address is not specified, the default function code is 0x02. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [ReadDiscreteAsync(String, UInt16)](9c34359d-3d18-fdbf-d6c3-35134eec3b9d.htm) | 批量的读取输入点，需要指定起始地址，读取长度，如果富文本地址不指定，默认使用的功能码是 0x02  To read input points in batches, you need to specify the start address and read length. If the rich text address is not specified, the default function code is 0x02 (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [ReadDouble(String)](e7e31c43-71b4-e692-6bdc-ba9533e1d83d.htm) | 读取双浮点的数据  Read double floating point data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadDouble(String, UInt16)](9df438ac-8869-5002-7ace-bfb7c8f7aa6b.htm) | 读取双浮点数据的数组  Read double floating point data array (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [ReadDoubleAsync(String)](2cf895fa-1033-1def-f048-0e8caefaa71d.htm) | 异步读取双浮点的数据  Asynchronously read double floating point data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadDoubleAsync(String, UInt16)](e5c602ad-d3ce-a677-908b-03d72a4388cd.htm) | 异步读取双浮点数据的数组  Asynchronously read double floating point data array (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [ReadFile](702ddeff-3ba5-952d-abce-99a988a1d3cc.htm) | 使用0x14功能码读取文件的数据内容，需要指定文件编号，地址，长度，返回读取的结果数据，读取的长度为任意长度，内部自动根据实际情况来切割操作。 (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [ReadFloat(String)](bc2801ae-6600-e845-7c6c-522c9156eb0e.htm) | 读取单浮点数据  Read single floating point data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadFloat(String, UInt16)](909486dc-cc84-c34f-2de5-ec17b3d04182.htm) | 读取单浮点精度的数组  Read single floating point array (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [ReadFloatAsync(String)](89f743f8-4b3a-7085-1e49-a1d7adb10052.htm) | 异步读取单浮点数据  Asynchronously read single floating point data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadFloatAsync(String, UInt16)](eaeb9795-d78c-f884-fa54-c698c5e37404.htm) | 异步读取单浮点精度的数组  Asynchronously read single floating point array (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [ReadFromCoreServer(IEnumerableByte)](044cfdd9-ec76-b89f-2f2b-5f439025006c.htm) | 将多个数据报文按顺序发到设备，并从设备接收返回的数据内容，然后拼接成一个Byte[]信息，需要重写[UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm)方法才能返回正确的结果。  Send multiple data packets to the device in sequence, and receive the returned data content from the device, and then splicing them into a Byte[] message, you need to rewrite [UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm) method to return the correct result. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServer(Byte)](ff061b13-df60-6617-08c2-40b81cde38ea.htm) | 将Modbus报文数据发送到当前的通道中，并从通道中接收Modbus的报文，通道将根据当前连接自动获取，本方法是线程安全的。  Send Modbus message data to the current channel, and receive Modbus messages from the channel. The channel will automatically obtain it according to the current connection. This method is thread-safe. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServer(Byte, Boolean, Boolean)](b32ce8e8-a959-fcc2-6736-14b39fdb9a40.htm) | 将二进制的数据发送到管道中去，然后从管道里接收二进制的数据回来，并返回是否成功的结果对象。  Send binary data to the pipeline, and then receive binary data back from the pipeline, and return whether the success of the result object (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServer(CommunicationPipe, Byte, Boolean, Boolean)](ef36d56b-df09-f727-e8d3-316d37a13ee9.htm) | 使用指定的管道来进行数据通信，发送原始数据到管道，然后从管道接收相关的数据返回，本方法无锁 (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServerAsync(IEnumerableByte)](bcdb59b7-a70a-d5c0-1fc6-9ececcd39ea2.htm) | 将多个数据报文按顺序发到设备，并从设备接收返回的数据内容，然后拼接成一个Byte[]信息，需要重写[UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm)方法才能返回正确的结果。  Send multiple data packets to the device in sequence, and receive the returned data content from the device, and then splicing them into a Byte[] message, you need to rewrite [UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm) method to return the correct result. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServerAsync(Byte)](61427407-a616-84ae-62c4-abea8bc7a451.htm) | 将Modbus报文数据发送到当前的通道中，并从通道中接收Modbus的报文，通道将根据当前连接自动获取，本方法是线程安全的。  Send Modbus message data to the current channel, and receive Modbus messages from the channel. The channel will automatically obtain it according to the current connection. This method is thread-safe. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServerAsync(Byte, Boolean, Boolean)](1c79c204-5ab5-9342-e65d-7455c20ea577.htm) | 将二进制的数据发送到管道中去，然后从管道里接收二进制的数据回来，并返回是否成功的结果对象。  Send binary data to the pipeline, and then receive binary data back from the pipeline, and return whether the success of the result object (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServerAsync(CommunicationPipe, Byte, Boolean, Boolean)](a254a19d-3c29-147b-335b-8adccdf4ee2f.htm) | 使用指定的管道来进行数据通信，发送原始数据到管道，然后从管道接收相关的数据返回，本方法无锁 (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法代码示例 | [ReadInt16(String)](57a16003-2fd5-b2bd-09bf-14366be899b8.htm) | 读取16位的有符号的整型数据  Read 16-bit signed integer data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt16(String, UInt16)](b2bde158-f51c-93fc-66ca-5e4cf6ea1fb7.htm) | 读取16位的有符号整型数组  Read 16-bit signed integer array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt16Async(String)](ac021793-890d-ceab-eae4-47ff763390c1.htm) | 异步读取16位的有符号的整型数据  Asynchronously read 16-bit signed integer data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt16Async(String, UInt16)](f30e02fc-8078-f1b3-c226-5991e3e23caf.htm) | 异步读取16位的有符号整型数组  Asynchronously read 16-bit signed integer array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt32(String)](bfbea035-9a83-40e3-084c-28122041b37a.htm) | 读取32位的有符号整型  Read 32-bit signed integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt32(String, UInt16)](ca0b8024-b31c-4e45-ad3e-916e188efebe.htm) | 读取32位有符号整型数组  Read 32-bit signed integer array (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [ReadInt32Async(String)](96726ea0-82aa-1229-21f6-43cd62abca43.htm) | 异步读取32位的有符号整型  Asynchronously read 32-bit signed integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt32Async(String, UInt16)](4d8ec145-5fcb-0e5e-6065-957de3434321.htm) | 异步读取32位有符号整型数组  Asynchronously read 32-bit signed integer array (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [ReadInt64(String)](0d0c2284-2da1-6f9a-b353-537561cd5e95.htm) | 读取64位的有符号整型  Read 64-bit signed integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt64(String, UInt16)](662b75d9-9f79-a959-bd6f-a797741d84c8.htm) | 读取64位的有符号整型数组  Read 64-bit signed integer array (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [ReadInt64Async(String)](26d46c8a-3775-3abd-d3ec-5cd1758a317e.htm) | 异步读取64位的有符号整型  Asynchronously read 64-bit signed integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt64Async(String, UInt16)](c3fa5d8b-cf3e-2424-4541-99bb11f7b8d7.htm) | 异步读取64位的有符号整型数组  Asynchronously read 64-bit signed integer array (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [ReadString(String, UInt16)](07812d76-415a-6400-c8a2-1c21da51e7a3.htm) | 读取字符串数据，默认为最常见的ASCII编码  Read string data, default is the most common ASCII encoding (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadString(String, UInt16, Encoding)](62ed8198-3bc2-4374-b362-02b71138ab98.htm) | 使用指定的编码，读取字符串数据  Reads string data using the specified encoding (重写 [DeviceCommunicationReadString(String, UInt16, Encoding)](b34c71a6-58f3-6bf4-0b93-471be4f2ed31.htm).) |
| 公共方法代码示例 | [ReadStringAsync(String, UInt16)](1835da81-b811-d151-0b53-907d8d87d0a8.htm) | 异步读取字符串数据，默认为最常见的ASCII编码  Asynchronously read string data, default is the most common ASCII encoding (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadStringAsync(String, UInt16, Encoding)](8d084430-6e0d-a700-923a-ca7493cda482.htm) | 异步使用指定的编码，读取字符串数据  Asynchronously reads string data using the specified encoding (重写 [DeviceCommunicationReadStringAsync(String, UInt16, Encoding)](091dfc67-3d25-f33e-3a55-a0228863815e.htm).) |
| 公共方法代码示例 | [ReadStructT](49853082-fbad-52e5-d756-0615cedb4b83.htm) | 读取结构体类型的数据，根据结构体自身的定义，读取原始字节数组，然后解析出实际的结构体数据，结构体需要实现[HslStructAttribute](87e427a8-0c7b-43a8-c421-c27264eda90f.htm)特性 (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadStructAsyncT](36461e61-d0ed-62cf-79ce-dcc7f71a3fa2.htm) | 读取结构体类型的数据，根据结构体自身的定义，读取原始字节数组，然后解析出实际的结构体数据，结构体需要实现[HslStructAttribute](87e427a8-0c7b-43a8-c421-c27264eda90f.htm)特性 (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt16(String)](736379d1-9560-5052-658a-13c126d56d78.htm) | 读取16位的无符号整型  Read 16-bit unsigned integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt16(String, UInt16)](e284c4bd-d3f8-ff30-eb42-0aa035cfb4f6.htm) | 读取16位的无符号整型数组  Read 16-bit unsigned integer array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt16Async(String)](361acc93-43b5-a9a5-6fa7-fa8728b7b96c.htm) | 异步读取16位的无符号整型  Asynchronously read 16-bit unsigned integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt16Async(String, UInt16)](73dd24b8-b501-df53-1de6-ccb2708f656d.htm) | 异步读取16位的无符号整型数组  Asynchronously read 16-bit unsigned integer array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt32(String)](ca73c958-39c7-9b34-aa44-3db2d9b0cad3.htm) | 读取32位的无符号整型  Read 32-bit unsigned integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt32(String, UInt16)](7ef202c8-7f3a-596b-12c0-bf5e5a5abd92.htm) | 读取32位的无符号整型数组  Read 32-bit unsigned integer array (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [ReadUInt32Async(String)](2343953f-6350-331e-6ba7-18a4449f0575.htm) | 异步读取32位的无符号整型  Asynchronously read 32-bit unsigned integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt32Async(String, UInt16)](ca94295c-48ce-dbdf-01eb-4d01701a9f50.htm) | 异步读取32位的无符号整型数组  Asynchronously read 32-bit unsigned integer array (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [ReadUInt64(String)](947487dd-d4a3-72d8-d9fb-2224aad70e12.htm) | 读取64位的无符号整型  Read 64-bit unsigned integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt64(String, UInt16)](130d0d35-93a1-d72d-c2f4-53b525c11f02.htm) | 读取64位的无符号整型的数组  Read 64-bit unsigned integer array (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [ReadUInt64Async(String)](3411ff3a-772f-dbe1-8931-75e00a731579.htm) | 异步读取64位的无符号整型  Asynchronously read 64-bit unsigned integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt64Async(String, UInt16)](6f22e2db-a0eb-b4b7-f43b-030555667ade.htm) | 异步读取64位的无符号整型的数组  Asynchronously read 64-bit unsigned integer array (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [ReadWrite](455d5de4-9bc6-c0c0-db5c-b6a334800126.htm) | 使用0x17功能码来实现同时写入并读取数据的操作，使用一条报文来实现，需要指定读取的地址，长度，写入的地址，写入的数据信息，返回读取的结果数据。  Use 0x17 function code to write and read data at the same time, and use a message to implement it, you need to specify the read address, length, written address, written data information, and return the read result data. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [RegisteredAddressMapping](a00ba5f5-373c-6e60-845a-3d8ad093568e.htm) | 注册一个新的地址映射关系，注册地址映射关系后，就可以使用新的地址来读写Modbus数据了，通常用于其他的支持Modbus协议的PLC。  After registering a new address mapping, you can use the new address to read and write Modbus data, which is usually used for other PLCs that support the Modbus protocol. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [SerialPortInni(ActionSerialPort)](3728a0f2-4553-0bc8-5c92-253c13e84f05.htm) | 根据自定义初始化方法进行初始化串口信息  Initialize the serial port information according to the custom initialization method (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 公共方法 | [SerialPortInni(String)](f82e0ba4-af2f-5165-dcde-e047c3351e74.htm) | 初始化串口信息，9600波特率，8位数据位，1位停止位，无奇偶校验  Initial serial port information, 9600 baud rate, 8 data bits, 1 stop bit, no parity (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 公共方法 | [SerialPortInni(String, Int32)](433523a8-1606-f029-50db-68ef68034046.htm) | 初始化串口信息，波特率，8位数据位，1位停止位，无奇偶校验  Initializes serial port information, baud rate, 8-bit data bit, 1-bit stop bit, no parity (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 公共方法 | [SerialPortInni(String, Int32, Int32, StopBits, Parity)](cefa7cd2-9c1f-b79d-9ea3-d4d27d7ddeae.htm) | 初始化串口信息，波特率，数据位，停止位，奇偶校验需要全部自己来指定  Start serial port information, baud rate, data bit, stop bit, parity all need to be specified (继承自 [DeviceSerialPort](591fa2a4-f95b-44f5-5afb-e9238bf499e4.htm)。) |
| 公共方法 | [SetDtuPipe](f765d1d7-234a-b229-9457-97e25e024dc9.htm) | 将当前的通信对象设置DTU模式，允许传入现成的管道，并返回初始化结果，如果该设备重写了握手报文，就是返回握手结果  Set the current communication object to DTU mode, allow the existing pipe to be passed in, and return the initialization result, if the device rewrites the handshake packet, the handshake result is returned (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [SetDtuPipeAsync](9d1fd12d-b804-3828-a627-a2c0eef1d23f.htm) | 将当前的通信对象设置DTU模式，允许传入现成的管道，并返回初始化结果，如果该设备重写了握手报文，就是返回握手结果  Set the current communication object to DTU mode, allow the existing pipe to be passed in, and return the initialization result, if the device rewrites the handshake packet, the handshake result is returned (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ToString](a8ec9534-c499-74d1-b92e-bcb0edc12a82.htm) | (重写 [ModbusRtuToString](11f9cb41-d3e4-d45d-abd1-36e7bb2577e6.htm).) |
| 公共方法 | [TranslateToModbusAddress](e554c6c9-1c35-e7b1-5c1e-83bc6d4f4287.htm) | 将当前的地址信息转换成Modbus格式的地址，如果转换失败，返回失败的消息。默认不进行任何的转换。  Convert the current address information into a Modbus format address. If the conversion fails, a failure message will be returned. No conversion is performed by default. (重写 [ModbusRtuTranslateToModbusAddress(String, Byte)](43c76b4c-2f7e-c60a-8b90-ee2c3898ce42.htm).) |
| 公共方法 | [UnpackResponseContent](2af5ade0-337b-6bff-777e-99f624a94606.htm) | 根据对方返回的报文命令，对命令进行基本的拆包，例如各种Modbus协议拆包为统一的核心报文，还支持对报文的验证  According to the message command returned by the other party, the command is basically unpacked, for example, various Modbus protocols are unpacked into a unified core message, and the verification of the message is also supported (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
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
| 公共方法代码示例 | [Write(String, Double)](64c44576-7e95-cbcb-a8c9-e6b7ff7e3c50.htm) | 写入double数据，返回是否成功  Write double data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Int16)](c7cfd67d-ca31-3087-c684-ee4c10780ab1.htm) | 写入short数组，返回是否成功  Write short array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Int32)](ba80d1a5-c8d8-7e0b-5c6b-9435cf20cadb.htm) | 写入int数据，返回是否成功  Write int data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Int64)](ba246b33-5a12-a22c-6d40-0969a64b508b.htm) | 写入long数据，返回是否成功  Write long data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Single)](a2265fa7-5b56-a5aa-9d21-a0dc5458cbf9.htm) | 写入float数据，返回是否成功  Write float data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, String)](4e0ec363-ba64-7ca7-b8d3-2e9a8087bd55.htm) | 写入字符串信息，编码为ASCII  Write string information, encoded as ASCII (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, UInt16)](65544962-82ac-bb47-d72f-5dd0c3d4ac6c.htm) | 写入ushort数组，返回是否成功  Write ushort array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, UInt32)](c99d8246-f95e-4f79-738b-16e622072eb4.htm) | 写入uint数据，返回是否成功  Write uint data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, UInt64)](4924086e-9b36-d446-4ae4-6682134a6ad1.htm) | 写入ulong数据，返回是否成功  Write ulong data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [Write(String, Boolean)](a40c7e8d-5770-9ad2-e263-66a92c45f61c.htm) | 向线圈中写入bool数值，返回是否写入成功，如果富文本地址不指定，默认使用的功能码是 0x05， 如果你的地址为字地址，例如100.2，那么将使用0x16的功能码，通过掩码的方式来修改寄存器的某一位，需要Modbus服务器支持，否则自动切换为读取字数据，修改位，在写入字的方式。  Write bool value to the coil and return whether the writing is successful. If the rich text address is not specified, the default function code is 0x05. If your address is a word address, such as 100.2, then you will use the function code of 0x16 to modify a bit of the register through a mask. It needs Modbus server support, Otherwise, it automatically switches to read the word data, modifies the bits, and writes the word in the way. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [Write(String, Boolean)](fc664b2b-f304-683e-573c-239379b83eb0.htm) | 向线圈中写入bool数组，返回是否写入成功，如果富文本地址不指定，默认使用的功能码是 0x0F  Write the bool array to the coil, and return whether the writing is successful. If the rich text address is not specified, the default function code is 0x0F. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [Write(String, Byte)](ec533f87-fd4a-c833-0256-fdd83da3609f.htm) | 将数据写入到Modbus的寄存器上去，需要指定起始地址和数据内容，如果富文本地址不指定，默认使用的功能码是 0x10  To write data to Modbus registers, you need to specify the start address and data content. If the rich text address is not specified, the default function code is 0x10 (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [Write(String, Double)](e4621d28-0353-fd60-d600-f457a46a661e.htm) | 写入double数组，返回是否成功  Write double array, return whether the write was successful (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [Write(String, Int16)](87345acc-1644-e27d-f5b3-51a01e0d5b02.htm) | 将数据写入到Modbus的单个寄存器上去，需要指定起始地址和数据值，如果富文本地址不指定，默认使用的功能码是 0x06  To write data to a single register of Modbus, you need to specify the start address and data value. If the rich text address is not specified, the default function code is 0x06. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [Write(String, Int32)](cb0f528c-c39e-a439-6093-65c7c8e8ca10.htm) | 写入int[]数组，返回是否成功  Write int array, return whether the write was successful (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [Write(String, Int64)](d3322d92-64ca-0080-27c6-b8693206ad96.htm) | 写入long数组，返回是否成功  Write long array, return whether the write was successful (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [Write(String, Single)](08a2fa24-4d18-9e7d-1fc9-873e8870ecc2.htm) | 写入float数组，返回是否成功  Write float array, return whether the write was successful (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [Write(String, UInt16)](4c43fcf9-d5c0-e4e4-70c2-c2d1aa8bfe4b.htm) | 将数据写入到Modbus的单个寄存器上去，需要指定起始地址和数据值，如果富文本地址不指定，默认使用的功能码是 0x06  To write data to a single register of Modbus, you need to specify the start address and data value. If the rich text address is not specified, the default function code is 0x06. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [Write(String, UInt32)](3113766d-8db6-7f55-c9c9-0580b9b1b756.htm) | 写入uint[]数组，返回是否成功  Write uint array, return whether the write was successful (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [Write(String, UInt64)](414c37ec-7955-1077-addb-8a7c0a861215.htm) | 写入ulong数组，返回是否成功  Write ulong array, return whether the write was successful (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [Write(String, String, Int32)](02fe2ca9-49a3-ede4-d1af-8c6da1f644da.htm) | 写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为ASCII  Write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, it is filled with 0 and the encoding is ASCII. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, String, Encoding)](7b602661-902d-2b70-1fdb-ba87ecd307bf.htm) | 写入字符串信息，需要指定的编码信息  Write string information, need to specify the encoding information (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, String, Int32, Encoding)](e8cb30d9-4c5e-6c68-b9af-1b1b6c2dc714.htm) | 写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为指定的编码信息  Write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, then the operation is complemented with 0 , you should specified the encoding information (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteT(T)](9f48853f-fbeb-c37d-25a2-1c76d1dee407.htm) | 写入支持Hsl特性的数据，返回是否写入成功，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Write data that supports the Hsl attribute, and return whether the write was successful. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WriteAsync(String, Boolean)](fe872dae-7ef9-12c1-d883-e92ba8676416.htm) | 异步批量写入Boolean数组数据，返回是否成功  Asynchronously batch write Boolean array data, return success (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Byte)](8f1cec30-8b86-6749-c8f3-7d39168d59df.htm) | 异步写入原始的byte数组数据到指定的地址，返回是否写入成功  Asynchronously writes the original byte array data to the specified address, and returns whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Double)](068adaa4-d1de-dffc-8376-6f483e2c1d96.htm) | 异步写入double数据，返回是否成功  Asynchronously write double data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int16)](2625e868-7096-e526-ce25-b29afc639b41.htm) | 异步写入short数组，返回是否成功  Asynchronously write short array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int32)](df91d629-2080-328d-67cd-fbefb9a45404.htm) | 异步写入int数据，返回是否成功  Asynchronously write int data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int64)](ec6c1e24-a0d8-f1f5-56af-109b0c546fc8.htm) | 异步写入long数据，返回是否成功  Asynchronously write long data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Single)](910189df-13f9-b7d0-a626-b5b25ce5d32a.htm) | 异步写入float数据，返回是否成功  Asynchronously write float data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, String)](fe839c9a-46f3-e317-0155-b10a8a83a9b1.htm) | 异步写入字符串信息，编码为ASCII  Asynchronously write string information, encoded as ASCII (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt16)](34a1a321-972e-06a8-76fb-794dd1229005.htm) | 异步写入ushort数组，返回是否成功  Asynchronously write ushort array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt32)](abad215c-bc99-ac90-b9f3-7adb3a7cca8f.htm) | 异步写入uint数据，返回是否成功  Asynchronously write uint data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt64)](e396074d-6cf1-b45d-5e8e-01026f6e391b.htm) | 异步写入ulong数据，返回是否成功  Asynchronously write ulong data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WriteAsync(String, Boolean)](90b9711b-1891-383c-0e56-19dd224b7aab.htm) | 向线圈中写入bool数值，返回是否写入成功，如果富文本地址不指定，默认使用的功能码是 0x05， 如果你的地址为字地址，例如100.2，那么将使用0x16的功能码，通过掩码的方式来修改寄存器的某一位，需要Modbus服务器支持，否则自动切换为读取字数据，修改位，在写入字的方式。  Write bool value to the coil and return whether the writing is successful. If the rich text address is not specified, the default function code is 0x05. If your address is a word address, such as 100.2, then you will use the function code of 0x16 to modify a bit of the register through a mask. It needs Modbus server support, Otherwise, it automatically switches to read the word data, modifies the bits, and writes the word in the way. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Double)](d67e053b-fae0-fbbb-d688-b514f417bbd7.htm) | 异步写入double数组，返回是否成功  Asynchronously write double array, return whether the write was successful (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [WriteAsync(String, Int16)](7e819772-f315-0f56-8c1e-8e2bd3aa2396.htm) | 将数据写入到Modbus的单个寄存器上去，需要指定起始地址和数据值，如果富文本地址不指定，默认使用的功能码是 0x06  To write data to a single register of Modbus, you need to specify the start address and data value. If the rich text address is not specified, the default function code is 0x06. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int32)](9e2c12ba-21f7-1c91-71ae-1eb2f2f853c3.htm) | 异步写入int[]数组，返回是否成功  Asynchronously write int array, return whether the write was successful (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int64)](3ef3dd46-318a-f55d-7296-57d303aafc93.htm) | 异步写入long数组，返回是否成功  Asynchronously write long array, return whether the write was successful (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Single)](ffa2fa0f-8321-2c1f-3916-c2cd90882f3c.htm) | 异步写入float数组，返回是否成功  Asynchronously write float array, return whether the write was successful (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [WriteAsync(String, UInt16)](4bb55e75-5cb4-4b22-c8b7-2983eae63e31.htm) | 将数据写入到Modbus的单个寄存器上去，需要指定起始地址和数据值，如果富文本地址不指定，默认使用的功能码是 0x06  To write data to a single register of Modbus, you need to specify the start address and data value. If the rich text address is not specified, the default function code is 0x06. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt32)](e108f604-e11a-2599-b88c-d69e630212b5.htm) | 异步写入uint[]数组，返回是否成功  Asynchronously write uint array, return whether the write was successful (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt64)](d304bc30-22bb-24bc-dd75-f767cfba2c35.htm) | 异步写入ulong数组，返回是否成功  Asynchronously write ulong array, return whether the write was successful (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [WriteAsync(String, String, Int32)](b739b983-b03c-8ecf-2713-75f19233ae2f.htm) | 异步写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为ASCII  Asynchronously write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, it is filled with 0 and the encoding is ASCII. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, String, Encoding)](876e9800-0406-6407-bdf8-8ab1cf28df38.htm) | 异步写入字符串信息，需要指定的编码信息  Asynchronously write string information, need to specify the encoding information (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WriteAsync(String, String, Int32, Encoding)](ae9fb2a3-fddb-48eb-e9ad-6cc0ebbac1f6.htm) | 异步写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为指定的编码信息  Asynchronously write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, then the operation is complemented with 0 , you should specified the encoding information (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsyncT(T)](25d96e77-d1b2-918c-33f9-e5e4cc57c9cd.htm) | 异步写入支持Hsl特性的数据，返回是否写入成功，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Asynchronously write data that supports the Hsl attribute, and return whether the write was successful. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteCustomerT](84a37f82-3f67-a0f8-fd3f-0630772565e4.htm) | 写入自定义类型的数据，该类型必须继承自IDataTransfer接口  Write data of a custom type, which must inherit from the IDataTransfer interface (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteCustomerAsyncT](04628397-cfa1-d785-7790-4c02cb377ab6.htm) | 写入自定义类型的数据，该类型必须继承自IDataTransfer接口  Write data of a custom type, which must inherit from the IDataTransfer interface (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WriteFile](8b9201b8-f363-7860-6b09-e63811574239.htm) | 使用0x15功能码写入文件的数据内容，需要指定文件编号，地址，数据，返回写入的结果数据。写入的数据字节长度不能大于248个字节，如果想写入更多的字节，就需要切片写入操作 (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [WriteMask](6e5e9e07-cc0e-f60c-07be-eb4a63249863.htm) | 向设备写入掩码数据，使用0x16功能码，需要确认对方是否支持相关的操作，掩码数据的操作主要针对寄存器。  To write mask data to the server, using the 0x16 function code, you need to confirm whether the other party supports related operations. The operation of mask data is mainly directed to the register. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [WriteMaskAsync](fd2cff63-8ffc-bf53-a7b6-24253b445aff.htm) | 向设备写入掩码数据，使用0x16功能码，需要确认对方是否支持相关的操作，掩码数据的操作主要针对寄存器。  To write mask data to the server, using the 0x16 function code, you need to confirm whether the other party supports related operations. The operation of mask data is mainly directed to the register. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [WriteOneRegister(String, Int16)](ed003503-3828-ed79-d4b4-08cfc417c8ab.htm) | 将数据写入到Modbus的单个寄存器上去，需要指定起始地址和数据值，如果富文本地址不指定，默认使用的功能码是 0x06  To write data to a single register of Modbus, you need to specify the start address and data value. If the rich text address is not specified, the default function code is 0x06. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [WriteOneRegister(String, UInt16)](9603651d-3202-1560-9141-b2c33949e06c.htm) | 将数据写入到Modbus的单个寄存器上去，需要指定起始地址和数据值，如果富文本地址不指定，默认使用的功能码是 0x06  To write data to a single register of Modbus, you need to specify the start address and data value. If the rich text address is not specified, the default function code is 0x06. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [WriteOneRegisterAsync(String, Int16)](b5f960df-3437-910f-aba4-f4dcbbfedad8.htm) | 将数据写入到Modbus的单个寄存器上去，需要指定起始地址和数据值，如果富文本地址不指定，默认使用的功能码是 0x06  To write data to a single register of Modbus, you need to specify the start address and data value. If the rich text address is not specified, the default function code is 0x06. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |
| 公共方法 | [WriteOneRegisterAsync(String, UInt16)](66dad72a-c201-d696-fdf2-bcdac070fc57.htm) | 将数据写入到Modbus的单个寄存器上去，需要指定起始地址和数据值，如果富文本地址不指定，默认使用的功能码是 0x06  To write data to a single register of Modbus, you need to specify the start address and data value. If the rich text address is not specified, the default function code is 0x06. (继承自 [ModbusRtu](a0944e49-821b-30ff-a5a0-0d8c71c64280.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[InovanceSerial 类](434659f2-7d84-27a2-72b5-aebd489024da.htm)

[HslCommunication.Profinet.Inovance 命名空间](aa61d757-9239-26ce-525f-b7e52699e229.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadByte 方法 

[原文連結](http://api.hslcommunication.cn/html/f2698a84-b0c3-6069-888e-4fcbb045db72.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Inovance](../html/aa61d757-9239-26ce-525f-b7e52699e229.htm "HslCommunication.Profinet.Inovance")

[InovanceSerial 类](../html/434659f2-7d84-27a2-72b5-aebd489024da.htm "InovanceSerial 类")

[InovanceSerial 方法](../html/a202dad8-3af6-5794-c23f-1ddd2151e1a6.htm "InovanceSerial 方法")

[ReadByte 方法](../html/f2698a84-b0c3-6069-888e-4fcbb045db72.htm "ReadByte 方法 ")

[ReadByteAsync 方法](../html/70ab77e7-0972-e072-3e26-206a0fdf51e4.htm "ReadByteAsync 方法 ")

[ReadString 方法](../html/47ff96ed-6d9a-2110-e28c-8a4c529ef384.htm "ReadString 方法 ")

[ReadStringAsync 方法](../html/e62e86b6-223b-20b8-0c2c-b28af24cf3fe.htm "ReadStringAsync 方法 ")

[ToString 方法](../html/a8ec9534-c499-74d1-b92e-bcb0edc12a82.htm "ToString 方法 ")

[TranslateToModbusAddress 方法](../html/e554c6c9-1c35-e7b1-5c1e-83bc6d4f4287.htm "TranslateToModbusAddress 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| InovanceSerialReadByte 方法 |

按照字节读取汇川M地址的数据，地址示例： MB100，MB101，需要注意的是，MB100 及 MB101 的地址是 MW50 的数据。  
Read the data of Inovance M address according to the byte, address example: MB100, MB101, it should be noted that the addresses of MB100 and MB101 are the data of MW50.

**命名空间：**
 [HslCommunication.Profinet.Inovance](aa61d757-9239-26ce-525f-b7e52699e229.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public OperateResult<byte> ReadByte(
	string address
)
```

```
Public Function ReadByte ( 
	address As String
) As OperateResult(Of Byte)
```

```
public:
OperateResult<unsigned char>^ ReadByte(
	String^ address
)
```

```
member ReadByte : 
        address : string -> OperateResult<byte> 
```

#### 参数

address
:   类型：SystemString  
    地址信息

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
读取的结果数据

![](../icons/SectionExpanded.png)参见

#### 引用

[InovanceSerial 类](434659f2-7d84-27a2-72b5-aebd489024da.htm)

[HslCommunication.Profinet.Inovance 命名空间](aa61d757-9239-26ce-525f-b7e52699e229.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadByteAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/70ab77e7-0972-e072-3e26-206a0fdf51e4.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Inovance](../html/aa61d757-9239-26ce-525f-b7e52699e229.htm "HslCommunication.Profinet.Inovance")

[InovanceSerial 类](../html/434659f2-7d84-27a2-72b5-aebd489024da.htm "InovanceSerial 类")

[InovanceSerial 方法](../html/a202dad8-3af6-5794-c23f-1ddd2151e1a6.htm "InovanceSerial 方法")

[ReadByte 方法](../html/f2698a84-b0c3-6069-888e-4fcbb045db72.htm "ReadByte 方法 ")

[ReadByteAsync 方法](../html/70ab77e7-0972-e072-3e26-206a0fdf51e4.htm "ReadByteAsync 方法 ")

[ReadString 方法](../html/47ff96ed-6d9a-2110-e28c-8a4c529ef384.htm "ReadString 方法 ")

[ReadStringAsync 方法](../html/e62e86b6-223b-20b8-0c2c-b28af24cf3fe.htm "ReadStringAsync 方法 ")

[ToString 方法](../html/a8ec9534-c499-74d1-b92e-bcb0edc12a82.htm "ToString 方法 ")

[TranslateToModbusAddress 方法](../html/e554c6c9-1c35-e7b1-5c1e-83bc6d4f4287.htm "TranslateToModbusAddress 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| InovanceSerialReadByteAsync 方法 |

按照字节读取汇川M地址的数据，地址示例： MB100，MB101，需要注意的是，MB100 及 MB101 的地址是 MW50 的数据。  
Read the data of Inovance M address according to the byte, address example: MB100, MB101, it should be noted that the addresses of MB100 and MB101 are the data of MW50.

**命名空间：**
 [HslCommunication.Profinet.Inovance](aa61d757-9239-26ce-525f-b7e52699e229.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public Task<OperateResult<byte>> ReadByteAsync(
	string address
)
```

```
Public Function ReadByteAsync ( 
	address As String
) As Task(Of OperateResult(Of Byte))
```

```
public:
Task<OperateResult<unsigned char>^>^ ReadByteAsync(
	String^ address
)
```

```
member ReadByteAsync : 
        address : string -> Task<OperateResult<byte>> 
```

#### 参数

address
:   类型：SystemString  
    地址信息

#### 返回值

类型：Task[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
读取的结果数据

![](../icons/SectionExpanded.png)参见

#### 引用

[InovanceSerial 类](434659f2-7d84-27a2-72b5-aebd489024da.htm)

[HslCommunication.Profinet.Inovance 命名空间](aa61d757-9239-26ce-525f-b7e52699e229.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadString 方法 

[原文連結](http://api.hslcommunication.cn/html/47ff96ed-6d9a-2110-e28c-8a4c529ef384.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Inovance](../html/aa61d757-9239-26ce-525f-b7e52699e229.htm "HslCommunication.Profinet.Inovance")

[InovanceSerial 类](../html/434659f2-7d84-27a2-72b5-aebd489024da.htm "InovanceSerial 类")

[InovanceSerial 方法](../html/a202dad8-3af6-5794-c23f-1ddd2151e1a6.htm "InovanceSerial 方法")

[ReadString 方法](../html/47ff96ed-6d9a-2110-e28c-8a4c529ef384.htm "ReadString 方法 ")

[ReadString 方法 (String, UInt16, Encoding)](../html/62ed8198-3bc2-4374-b362-02b71138ab98.htm "ReadString 方法 (String, UInt16, Encoding)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| InovanceSerialReadString 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法代码示例 | [ReadString(String, UInt16)](07812d76-415a-6400-c8a2-1c21da51e7a3.htm) | 读取字符串数据，默认为最常见的ASCII编码  Read string data, default is the most common ASCII encoding (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadString(String, UInt16, Encoding)](62ed8198-3bc2-4374-b362-02b71138ab98.htm) | 使用指定的编码，读取字符串数据  Reads string data using the specified encoding (重写 [DeviceCommunicationReadString(String, UInt16, Encoding)](b34c71a6-58f3-6bf4-0b93-471be4f2ed31.htm).) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[InovanceSerial 类](434659f2-7d84-27a2-72b5-aebd489024da.htm)

[HslCommunication.Profinet.Inovance 命名空间](aa61d757-9239-26ce-525f-b7e52699e229.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadString 方法 (String, UInt16, Encoding)

[原文連結](http://api.hslcommunication.cn/html/62ed8198-3bc2-4374-b362-02b71138ab98.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Inovance](../html/aa61d757-9239-26ce-525f-b7e52699e229.htm "HslCommunication.Profinet.Inovance")

[InovanceSerial 类](../html/434659f2-7d84-27a2-72b5-aebd489024da.htm "InovanceSerial 类")

[InovanceSerial 方法](../html/a202dad8-3af6-5794-c23f-1ddd2151e1a6.htm "InovanceSerial 方法")

[ReadString 方法](../html/47ff96ed-6d9a-2110-e28c-8a4c529ef384.htm "ReadString 方法 ")

[ReadString 方法 (String, UInt16, Encoding)](../html/62ed8198-3bc2-4374-b362-02b71138ab98.htm "ReadString 方法 (String, UInt16, Encoding)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| InovanceSerialReadString 方法 (String, UInt16, Encoding) |

使用指定的编码，读取字符串数据  
Reads string data using the specified encoding

**命名空间：**
 [HslCommunication.Profinet.Inovance](aa61d757-9239-26ce-525f-b7e52699e229.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public override OperateResult<string> ReadString(
	string address,
	ushort length,
	Encoding encoding
)
```

```
Public Overrides Function ReadString ( 
	address As String,
	length As UShort,
	encoding As Encoding
) As OperateResult(Of String)
```

```
public:
virtual OperateResult<String^>^ ReadString(
	String^ address, 
	unsigned short length, 
	Encoding^ encoding
) override
```

```
abstract ReadString : 
        address : string * 
        length : uint16 * 
        encoding : Encoding -> OperateResult<string> 
override ReadString : 
        address : string * 
        length : uint16 * 
        encoding : Encoding -> OperateResult<string>
```

#### 参数

address
:   类型：SystemString  
    起始地址

length
:   类型：SystemUInt16  
    数据长度

encoding
:   类型：System.TextEncoding  
    指定的自定义的编码

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)String  
带有成功标识的string数据

#### 实现

[IReadWriteNetReadString(String, UInt16, Encoding)](fc289def-e53b-b47d-22fe-7e02b5e354e0.htm)

![](../icons/SectionExpanded.png)示例

以下为三菱的连接对象示例，其他的设备读写情况参照下面的代码：

String类型示例

[复制](# "复制")

```
MelsecMcNet melsec_net = new MelsecMcNet( "192.168.0.100", 6000 );

// 以下是简单的读取，没有仔细校验的方式
string d100_value = melsec_net.ReadString( "D100", 5 ).Content;

// 如果需要判断是否读取成功，使用 Unicode 编码即可读取中文，如果还是乱码，就需要自己指定编码来实现
OperateResult<string> R_d100_value = melsec_net.ReadString( "D100", 5, Encoding.Unicode );
if (R_d100_value.IsSuccess)
{
    // success
    string value = R_d100_value.Content;
}
else
{
    // failed
}
```

![](../icons/SectionExpanded.png)参见

#### 引用

[InovanceSerial 类](434659f2-7d84-27a2-72b5-aebd489024da.htm)

[ReadString 重载](47ff96ed-6d9a-2110-e28c-8a4c529ef384.htm)

[HslCommunication.Profinet.Inovance 命名空间](aa61d757-9239-26ce-525f-b7e52699e229.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadStringAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/e62e86b6-223b-20b8-0c2c-b28af24cf3fe.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Inovance](../html/aa61d757-9239-26ce-525f-b7e52699e229.htm "HslCommunication.Profinet.Inovance")

[InovanceSerial 类](../html/434659f2-7d84-27a2-72b5-aebd489024da.htm "InovanceSerial 类")

[InovanceSerial 方法](../html/a202dad8-3af6-5794-c23f-1ddd2151e1a6.htm "InovanceSerial 方法")

[ReadStringAsync 方法](../html/e62e86b6-223b-20b8-0c2c-b28af24cf3fe.htm "ReadStringAsync 方法 ")

[ReadStringAsync 方法 (String, UInt16, Encoding)](../html/8d084430-6e0d-a700-923a-ca7493cda482.htm "ReadStringAsync 方法 (String, UInt16, Encoding)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| InovanceSerialReadStringAsync 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法代码示例 | [ReadStringAsync(String, UInt16)](1835da81-b811-d151-0b53-907d8d87d0a8.htm) | 异步读取字符串数据，默认为最常见的ASCII编码  Asynchronously read string data, default is the most common ASCII encoding (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadStringAsync(String, UInt16, Encoding)](8d084430-6e0d-a700-923a-ca7493cda482.htm) | 异步使用指定的编码，读取字符串数据  Asynchronously reads string data using the specified encoding (重写 [DeviceCommunicationReadStringAsync(String, UInt16, Encoding)](091dfc67-3d25-f33e-3a55-a0228863815e.htm).) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[InovanceSerial 类](434659f2-7d84-27a2-72b5-aebd489024da.htm)

[HslCommunication.Profinet.Inovance 命名空间](aa61d757-9239-26ce-525f-b7e52699e229.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadStringAsync 方法 (String, UInt16, Encoding)

[原文連結](http://api.hslcommunication.cn/html/8d084430-6e0d-a700-923a-ca7493cda482.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Inovance](../html/aa61d757-9239-26ce-525f-b7e52699e229.htm "HslCommunication.Profinet.Inovance")

[InovanceSerial 类](../html/434659f2-7d84-27a2-72b5-aebd489024da.htm "InovanceSerial 类")

[InovanceSerial 方法](../html/a202dad8-3af6-5794-c23f-1ddd2151e1a6.htm "InovanceSerial 方法")

[ReadStringAsync 方法](../html/e62e86b6-223b-20b8-0c2c-b28af24cf3fe.htm "ReadStringAsync 方法 ")

[ReadStringAsync 方法 (String, UInt16, Encoding)](../html/8d084430-6e0d-a700-923a-ca7493cda482.htm "ReadStringAsync 方法 (String, UInt16, Encoding)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| InovanceSerialReadStringAsync 方法 (String, UInt16, Encoding) |

异步使用指定的编码，读取字符串数据  
Asynchronously reads string data using the specified encoding

**命名空间：**
 [HslCommunication.Profinet.Inovance](aa61d757-9239-26ce-525f-b7e52699e229.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public override Task<OperateResult<string>> ReadStringAsync(
	string address,
	ushort length,
	Encoding encoding
)
```

```
Public Overrides Function ReadStringAsync ( 
	address As String,
	length As UShort,
	encoding As Encoding
) As Task(Of OperateResult(Of String))
```

```
public:
virtual Task<OperateResult<String^>^>^ ReadStringAsync(
	String^ address, 
	unsigned short length, 
	Encoding^ encoding
) override
```

```
abstract ReadStringAsync : 
        address : string * 
        length : uint16 * 
        encoding : Encoding -> Task<OperateResult<string>> 
override ReadStringAsync : 
        address : string * 
        length : uint16 * 
        encoding : Encoding -> Task<OperateResult<string>>
```

#### 参数

address
:   类型：SystemString  
    起始地址

length
:   类型：SystemUInt16  
    数据长度

encoding
:   类型：System.TextEncoding  
    指定的自定义的编码

#### 返回值

类型：Task[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)String  
带有成功标识的string数据

#### 实现

[IReadWriteNetReadStringAsync(String, UInt16, Encoding)](db06baf4-e674-22de-e89e-cb26646413f8.htm)

![](../icons/SectionExpanded.png)示例

以下为三菱的连接对象示例，其他的设备读写情况参照下面的代码：

String类型示例

[复制](# "复制")

```
MelsecMcNet melsec_net = new MelsecMcNet( "192.168.0.100", 6000 );

// 以下是简单的读取，没有仔细校验的方式
string d100_value = melsec_net.ReadString( "D100", 5 ).Content;

// 如果需要判断是否读取成功，使用 Unicode 编码即可读取中文，如果还是乱码，就需要自己指定编码来实现
OperateResult<string> R_d100_value = melsec_net.ReadString( "D100", 5, Encoding.Unicode );
if (R_d100_value.IsSuccess)
{
    // success
    string value = R_d100_value.Content;
}
else
{
    // failed
}
```

![](../icons/SectionExpanded.png)参见

#### 引用

[InovanceSerial 类](434659f2-7d84-27a2-72b5-aebd489024da.htm)

[ReadStringAsync 重载](e62e86b6-223b-20b8-0c2c-b28af24cf3fe.htm)

[HslCommunication.Profinet.Inovance 命名空间](aa61d757-9239-26ce-525f-b7e52699e229.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ToString 方法 

[原文連結](http://api.hslcommunication.cn/html/a8ec9534-c499-74d1-b92e-bcb0edc12a82.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Inovance](../html/aa61d757-9239-26ce-525f-b7e52699e229.htm "HslCommunication.Profinet.Inovance")

[InovanceSerial 类](../html/434659f2-7d84-27a2-72b5-aebd489024da.htm "InovanceSerial 类")

[InovanceSerial 方法](../html/a202dad8-3af6-5794-c23f-1ddd2151e1a6.htm "InovanceSerial 方法")

[ReadByte 方法](../html/f2698a84-b0c3-6069-888e-4fcbb045db72.htm "ReadByte 方法 ")

[ReadByteAsync 方法](../html/70ab77e7-0972-e072-3e26-206a0fdf51e4.htm "ReadByteAsync 方法 ")

[ReadString 方法](../html/47ff96ed-6d9a-2110-e28c-8a4c529ef384.htm "ReadString 方法 ")

[ReadStringAsync 方法](../html/e62e86b6-223b-20b8-0c2c-b28af24cf3fe.htm "ReadStringAsync 方法 ")

[ToString 方法](../html/a8ec9534-c499-74d1-b92e-bcb0edc12a82.htm "ToString 方法 ")

[TranslateToModbusAddress 方法](../html/e554c6c9-1c35-e7b1-5c1e-83bc6d4f4287.htm "TranslateToModbusAddress 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| InovanceSerialToString 方法 |

[缺少 "M:HslCommunication.Profinet.Inovance.InovanceSerial.ToString" 的 <summary> 文档]

**命名空间：**
 [HslCommunication.Profinet.Inovance](aa61d757-9239-26ce-525f-b7e52699e229.htm)  
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

[缺少 "M:HslCommunication.Profinet.Inovance.InovanceSerial.ToString" 的 <returns> 文档]

![](../icons/SectionExpanded.png)参见

#### 引用

[InovanceSerial 类](434659f2-7d84-27a2-72b5-aebd489024da.htm)

[HslCommunication.Profinet.Inovance 命名空间](aa61d757-9239-26ce-525f-b7e52699e229.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## TranslateToModbusAddress 方法 

[原文連結](http://api.hslcommunication.cn/html/e554c6c9-1c35-e7b1-5c1e-83bc6d4f4287.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Inovance](../html/aa61d757-9239-26ce-525f-b7e52699e229.htm "HslCommunication.Profinet.Inovance")

[InovanceSerial 类](../html/434659f2-7d84-27a2-72b5-aebd489024da.htm "InovanceSerial 类")

[InovanceSerial 方法](../html/a202dad8-3af6-5794-c23f-1ddd2151e1a6.htm "InovanceSerial 方法")

[ReadByte 方法](../html/f2698a84-b0c3-6069-888e-4fcbb045db72.htm "ReadByte 方法 ")

[ReadByteAsync 方法](../html/70ab77e7-0972-e072-3e26-206a0fdf51e4.htm "ReadByteAsync 方法 ")

[ReadString 方法](../html/47ff96ed-6d9a-2110-e28c-8a4c529ef384.htm "ReadString 方法 ")

[ReadStringAsync 方法](../html/e62e86b6-223b-20b8-0c2c-b28af24cf3fe.htm "ReadStringAsync 方法 ")

[ToString 方法](../html/a8ec9534-c499-74d1-b92e-bcb0edc12a82.htm "ToString 方法 ")

[TranslateToModbusAddress 方法](../html/e554c6c9-1c35-e7b1-5c1e-83bc6d4f4287.htm "TranslateToModbusAddress 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| InovanceSerialTranslateToModbusAddress 方法 |

将当前的地址信息转换成Modbus格式的地址，如果转换失败，返回失败的消息。默认不进行任何的转换。  
Convert the current address information into a Modbus format address. If the conversion fails, a failure message will be returned. No conversion is performed by default.

**命名空间：**
 [HslCommunication.Profinet.Inovance](aa61d757-9239-26ce-525f-b7e52699e229.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public override OperateResult<string> TranslateToModbusAddress(
	string address,
	byte modbusCode
)
```

```
Public Overrides Function TranslateToModbusAddress ( 
	address As String,
	modbusCode As Byte
) As OperateResult(Of String)
```

```
public:
virtual OperateResult<String^>^ TranslateToModbusAddress(
	String^ address, 
	unsigned char modbusCode
) override
```

```
abstract TranslateToModbusAddress : 
        address : string * 
        modbusCode : byte -> OperateResult<string> 
override TranslateToModbusAddress : 
        address : string * 
        modbusCode : byte -> OperateResult<string>
```

#### 参数

address
:   类型：SystemString  
    传入的地址

modbusCode
:   类型：SystemByte  
    Modbus的功能码

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)String  
转换之后Modbus的地址

#### 实现

[IModbusTranslateToModbusAddress(String, Byte)](68338057-dac7-e74e-f841-de013cea8da3.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[InovanceSerial 类](434659f2-7d84-27a2-72b5-aebd489024da.htm)

[HslCommunication.Profinet.Inovance 命名空间](aa61d757-9239-26ce-525f-b7e52699e229.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## InovanceSerial 字段

[原文連結](http://api.hslcommunication.cn/html/c8f8f731-c2b4-9995-4aaa-ba397ee9696c.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Inovance](../html/aa61d757-9239-26ce-525f-b7e52699e229.htm "HslCommunication.Profinet.Inovance")

[InovanceSerial 类](../html/434659f2-7d84-27a2-72b5-aebd489024da.htm "InovanceSerial 类")

[InovanceSerial 构造函数](../html/9fd1afb8-916b-d937-2a93-6357b2a446b5.htm "InovanceSerial 构造函数 ")

[InovanceSerial 属性](../html/df3ac411-58e0-d3f4-5d96-4fc524001008.htm "InovanceSerial 属性")

[InovanceSerial 方法](../html/a202dad8-3af6-5794-c23f-1ddd2151e1a6.htm "InovanceSerial 方法")

[InovanceSerial 字段](../html/c8f8f731-c2b4-9995-4aaa-ba397ee9696c.htm "InovanceSerial 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| InovanceSerial 字段 |

[InovanceSerial](434659f2-7d84-27a2-72b5-aebd489024da.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)字段

|  | 名称 | 说明 |
| --- | --- | --- |
| 受保护的字段 | [LogMsgFormatBinary](2ed2ff87-caec-0e44-f08b-e2b015225a58.htm) | 设置日志记录报文是否二进制，如果为False，那就使用ASCII码  Set whether the log message is binary, if it is False, then use ASCII code (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[InovanceSerial 类](434659f2-7d84-27a2-72b5-aebd489024da.htm)

[HslCommunication.Profinet.Inovance 命名空间](aa61d757-9239-26ce-525f-b7e52699e229.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## InovanceSerialOverTcp 类

[原文連結](http://api.hslcommunication.cn/html/f82b954d-1a19-0b0e-f705-677d531e17f6.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Inovance](../html/aa61d757-9239-26ce-525f-b7e52699e229.htm "HslCommunication.Profinet.Inovance")

[InovanceSerialOverTcp 类](../html/f82b954d-1a19-0b0e-f705-677d531e17f6.htm "InovanceSerialOverTcp 类")

[InovanceSerialOverTcp 构造函数](../html/c346256c-782a-1381-f646-c6221355f534.htm "InovanceSerialOverTcp 构造函数 ")

[InovanceSerialOverTcp 属性](../html/6b2e9e95-ec0b-89f3-1c88-3abe92c1ab7a.htm "InovanceSerialOverTcp 属性")

[InovanceSerialOverTcp 方法](../html/ca7557a8-9346-30eb-bf40-1b757ad91f09.htm "InovanceSerialOverTcp 方法")

[InovanceSerialOverTcp 字段](../html/06aefc47-ebfd-6a8b-090b-c7539a496e22.htm "InovanceSerialOverTcp 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| InovanceSerialOverTcp 类 |

汇川的串口转网口通信协议，A适用于AM400、 AM400\_800、 AC800、H3U, XP, H5U 等系列底层走的是MODBUS-RTU协议，地址说明参见标记  
Huichuan's serial communication protocol is applicable to AM400, AM400\_800, AC800 and other series. The bottom layer is MODBUS-RTU protocol. For the address description, please refer to the mark

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  [HslCommunication.Core.NetBinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)  
    [HslCommunication.Core.DeviceDeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)  
      [HslCommunication.Core.DeviceDeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)  
        [HslCommunication.ModBusModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)  
          HslCommunication.Profinet.InovanceInovanceSerialOverTcp

**命名空间：**
 [HslCommunication.Profinet.Inovance](aa61d757-9239-26ce-525f-b7e52699e229.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class InovanceSerialOverTcp : ModbusRtuOverTcp
```

```
Public Class InovanceSerialOverTcp
	Inherits ModbusRtuOverTcp
```

```
public ref class InovanceSerialOverTcp : public ModbusRtuOverTcp
```

```
type InovanceSerialOverTcp =  
    class
        inherit ModbusRtuOverTcp
    end
```

InovanceSerialOverTcp 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [InovanceSerialOverTcp](523bbb48-d091-524a-3b2a-7e06104a0664.htm) | 实例化一个默认的对象 |
| 公共方法 | [InovanceSerialOverTcp(String, Int32, Byte)](160107db-161f-966a-1751-bf7fcd70572a.htm) | 通过指定站号，ip地址，端口号来实例化一个新的对象 |
| 公共方法 | [InovanceSerialOverTcp(InovanceSeries, String, Int32, Byte)](b65bbcec-6ea5-2f5e-ae3c-570cc11ef158.htm) | 通过指定站号，IP地址，端口以及PLC的系列来实例化一个新的对象  Instantiate a new object by specifying the station number and PLC series |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [AddressStartWithZero](f4122661-5369-da8f-7c8c-84c85d70a480.htm) | 获取或设置起始的地址是否从0开始，默认为True  Gets or sets whether the starting address starts from 0. The default is True (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共属性 | [BroadcastStation](f307b9fd-b01f-cb79-bb97-0d1d5e8c7d29.htm) | 获取或是设置当前广播模式对应的站号，广播模式意味着不接收设备方的数据返回操作，默认为 -1，表示不使用广播模式。  Gets or sets the station number corresponding to the current broadcast mode. Broadcast mode means that the data return operation of the device is not received. The default value is -1, indicating that broadcast mode is not used. (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共属性代码示例 | [ByteTransform](f35118c9-3691-f752-ae6f-2c5549714faa.htm) | 当前的数据变换机制，当你需要从字节数据转换类型数据的时候需要。  The current data transformation mechanism is required when you need to convert type data from byte data. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共属性 | [CommunicationPipe](ad8e903b-d289-21b0-077e-15ede8005221.htm) | 获取或设置当前的管道信息，管道类型为[CommunicationPipe](ad8e903b-d289-21b0-077e-15ede8005221.htm)的继承类，内置了[PipeTcpNet](26b94bd5-2695-36e7-6378-c3d41846be8f.htm)管道，[PipeUdpNet](5a0b810e-187e-826b-0a13-34fd5576aa0e.htm)管道，[PipeSerialPort](9d924607-3dfc-3916-b532-c5dec48fe347.htm)管道等  Get or set the current pipeline information, the pipeline type is [CommunicationPipe](ad8e903b-d289-21b0-077e-15ede8005221.htm) inheritance class, [PipeTcpNet](26b94bd5-2695-36e7-6378-c3d41846be8f.htm) pipeline, [PipeUdpNet](5a0b810e-187e-826b-0a13-34fd5576aa0e.htm) pipeline, [PipeSerialPort](9d924607-3dfc-3916-b532-c5dec48fe347.htm) pipeline, etc (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [ConnectionId](ce1efb4b-fbc9-4683-163e-0a8e71e99d6b.htm) | 当前连接的唯一ID号，默认为长度20的guid码加随机数组成，方便列表管理，也可以自己指定  The unique ID number of the current connection. The default is a 20-digit guid code plus a random number. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性代码示例 | [ConnectTimeOut](71a651e5-737d-6e88-75a2-fb891b9af6cd.htm) | 获取或设置接收服务器反馈的时间，如果为负数，则不接收反馈   Gets or sets the time to receive server feedback, and if it is a negative number, does not receive feedback (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共属性 | [Crc16CheckEnable](fe1f096a-5f06-b612-7dab-4e56be636a76.htm) | 获取或设置是否启用CRC16校验码的检查功能，默认启用，如果需要忽略检查CRC16，则设置为 false 即可。  Gets or sets whether to enable the check function of CRC16 check code. It is enabled by default. If you need to ignore the check of CRC16, you can set it to false. (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共属性 | [DataFormat](ee89a24c-d5b0-82e0-9c40-9cc4ab57f2b7.htm) | 获取或设置数据解析的格式，可选ABCD, BADC，CDAB，DCBA格式，对int,uint,float,double,long,ulong类型有作用  Get or set the format of the data analysis, optional ABCD, BADC, CDAB, DCBA format, effective for int, uint, float, double, long, ulong type (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共属性 | [DisableFunctionCode06](a1be5198-fdfa-3e32-32ed-0edbdaf577ce.htm) | 获取或设置是否禁用功能码05的写入操作，设置为true时，表示不允许使用功能码 06 进行写入操作。写入单个的字数据也使用 0x10 功能码  Get or set whether to disable the write operation of function code 06. When set to true, it indicates that the write operation of function code 05 is not allowed. Writing individual character data also uses the 0x10 function code (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共属性 | [EnableWriteMaskCode](358ac1fd-3000-898e-1c63-6213649d851a.htm) | 获取或设置当前掩码写入的功能码是否激活状态，设置为 false 时，再执行写入位时，会通过读字，修改位，写字的方式来间接实现。  When the function code is set to false, and then the write bit is executed, it will be indirectly implemented by reading, modifying the bit, and writing the word. (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共属性代码示例 | [IpAddress](c378f4e8-3fdf-1932-7239-ed44bc5af584.htm) | 获取或是设置远程服务器的IP地址，如果是本机测试，那么需要设置为127.0.0.1   Get or set the IP address of the remote server. If it is a local test, then it needs to be set to 127.0.0.1 (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共属性 | [IsStringReverse](2b20320d-43f0-d660-a603-580f923e9e43.htm) | 字符串数据是否按照字来反转，默认为False  Whether the string data is reversed according to words. The default is False. (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共属性 | [LocalBinding](622b418f-5969-3726-5dcc-04456c4b97a7.htm) | 获取或设置绑定的本地的IP地址和端口号信息，如果端口设置为0，代表任何可用的端口  Get or set the bound local IP address and port number information, if the port is set to 0, it means any available port (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共属性代码示例 | [LogNet](8c56f246-88ac-afb8-8c06-972fafe4f78a.htm) | 组件的日志工具，支持日志记录，只要实例化后，当前网络的基本信息，就以[DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm)等级进行输出  The component's logging tool supports logging. As long as the instantiation of the basic network information, the output will be output at [DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm) (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性代码示例 | [Port](6f1ad24a-c916-5b67-aa17-a8c7dc2acb2e.htm) | 获取或设置服务器的端口号，具体的值需要取决于对方的配置  Gets or sets the port number of the server. The specific value depends on the configuration of the other party. (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共属性代码示例 | [ReceiveTimeOut](250560a7-332b-ce1b-5eb0-e3abe4ee6012.htm) | 获取或设置接收服务器反馈的时间，如果为负数，则不接收反馈   Gets or sets the time to receive server feedback, and if it is a negative number, does not receive feedback (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [SendBeforeHex](72aa1de4-116f-abfe-4adb-814a0be0e560.htm) | 获取或设置在发送通信报文前追加发送的字节信息，HEX格式，通常用于lora组网时，需要携带 00 00 00 02 四个字节的站地址功能。  Obtain or set the byte information sent before sending communication packets, HEX format, usually used for LORA networking, you need to carry 00 00 00 02 four-byte station address function. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [Series](565f07e7-017e-9072-ab09-c1f21bc4f29a.htm) | 获取或设置汇川的系列，默认为AM系列 |
| 公共属性 | [SleepTime](86a5fbb9-3456-6beb-7113-5bbad5eabbc5.htm) | 获取或设置在正式接收对方返回数据前的时候，需要休息的时间，当设置为0的时候，不需要休息。  Get or set the time required to rest before officially receiving the data from the other party. When it is set to 0, no rest is required. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [SocketKeepAliveTime](96ab7703-8c4d-b4c2-ddd9-0dc93ebcfc8d.htm) | 获取或设置客户端的Socket的心跳时间信息，这个是Socket底层自动实现的心跳包，不基于协议层实现。默认小于0，不开启心跳检测，如果需要开启，设置 60\_000 比较合适，单位毫秒  Get or set the heartbeat time information of the Socket of the client. This is the heartbeat packet automatically implemented by the bottom layer of the Socket, not based on the protocol layer. The default value is less than 0, and heartbeat detection is not enabled. If you need to enable it, it is more appropriate to set 60\_000, in milliseconds. (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共属性 | [Station](9e7022b4-45cc-58e8-8606-f9866a9e4e66.htm) | 获取或者重新修改服务器的默认站号信息，当然，你可以再读写的时候动态指定，例如地址 "s=2;100", 更详细的例子可以参考DEMO界面上的地址示例。  Get or re-modify the server's default station number information, of course, you can dynamically specify when reading and writing, such as the address "s=2; 100", a more detailed example can refer to the DEMO interface address example. (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共属性 | [StationCheckMacth](861ed947-785a-3866-afd5-ec8042e7bebc.htm) | 获取或设置当前是否启用站号检查的功能，默认启用，读写数据时将进行站号确认操作，如果需要忽略站号，则设置为 false 即可。  Gets or sets whether the station ID check function is enabled. It is enabled by default. When reading or writing data, the station ID is confirmed. (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 受保护的属性 | [WordLength](f1a505cc-e85c-6c8b-b189-58a1ab5753f9.htm) | 一个字单位的数据表示的地址长度，西门子为2，三菱，欧姆龙，modbusTcp就为1，AB PLC无效  The address length represented by one word of data, Siemens is 2, Mitsubishi, Omron, modbusTcp is 1, AB PLC is invalid (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共属性 | [WordReadBatchLength](1d4b2647-54fd-22b6-dddf-e88356c270a5.htm) | 当进行字读取的时候，商业授权用户支持分批读取，本属性表示分批的长度信息，默认为 120，可以根据实际情况来设置其他值。  When reading of word type, commercial authorization users support batch reading. This attribute represents the length information of the batch, with a default of 120. Other values can be set according to the actual situation. (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法代码示例 | [ConnectClose](569a9c0e-7628-b8d5-e124-0853452747b8.htm) | 手动断开与远程服务器的连接，如果当前是长连接模式，那么就会切换到短连接模式  Manually disconnect from the remote server, if it is currently in long connection mode, it will switch to short connection mode (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共方法代码示例 | [ConnectCloseAsync](76be8150-ff22-95b9-16b1-7c3440c54ae8.htm) | 手动断开与远程服务器的连接，如果当前是长连接模式，那么就会切换到短连接模式  Manually disconnect from the remote server, if it is currently in long connection mode, it will switch to short connection mode (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共方法代码示例 | [ConnectServer](6299cc1f-2fab-b754-4597-1eddc2b21308.htm) | 尝试连接远程的服务器，如果连接成功，就切换短连接模式到长连接模式，后面的每次请求都共享一个通道，使得通讯速度更快速  Try to connect to a remote server. If the connection is successful, switch the short connection mode to the long connection mode. Each subsequent request will share a channel, making the communication speed faster. (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共方法代码示例 | [ConnectServerAsync](e7eff5c0-8e2a-3b6c-5027-16f4d8bcf5ae.htm) | 尝试连接远程的服务器，如果连接成功，就切换短连接模式到长连接模式，后面的每次请求都共享一个通道，使得通讯速度更快速  Try to connect to a remote server. If the connection is successful, switch the short connection mode to the long connection mode. Each subsequent request will share a channel, making the communication speed faster. (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 受保护的方法 | [DecideWhetherQAMessage](acc40cda-e310-8537-6685-59b67eeb16ff.htm) | 决定当前的消息是否是用于问答机制返回的消息，默认直接返回 true, 实际的情况需要根据协议进行重写方法  To determine whether the current message is the message returned by the question answering mechanism, the default is true. In actual cases, the rewriting method needs to be performed according to the protocol (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [Dispose](b9dda6bf-d342-254e-689a-93fa0bb265d1.htm) | (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 受保护的方法 | [Dispose(Boolean)](0bbb6ed3-111e-069c-e9f5-8cbc3f7f7b77.htm) | (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | [ExtraAfterReadFromCoreServer](06a49a15-a3ce-5e39-28c7-2b78851eca2e.htm) | 和服务器交互完成的时候调用的方法，可以根据读写结果进行一些额外的操作，具体的操作需要根据实际的需求来重写实现  The method called when the interaction with the server is completed can perform some additional operations based on the read and write results. The specific operations need to be rewritten according to actual needs. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法代码示例 | [ExtraOnDisconnect](de9dbba2-3b2c-1c5e-53c8-4e6f8a44e2c4.htm) | 根据实际的协议选择是否重写本方法，有些协议在断开连接之前，需要发送一些报文来关闭当前的网络通道  Select whether to rewrite this method according to the actual protocol. Some protocols need to send some packets to close the current network channel before disconnecting. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法代码示例 | [ExtraOnDisconnectAsync](210f593e-16d8-77c4-1178-58fbcfe9b27b.htm) | 根据实际的协议选择是否重写本方法，有些协议在断开连接之前，需要发送一些报文来关闭当前的网络通道  Select whether to rewrite this method according to the actual protocol. Some protocols need to send some packets to close the current network channel before disconnecting. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 受保护的方法 | [GetLogTextFromBinary](67f3b2f7-957c-fa0b-8320-6799237157db.htm) | 获取当前的报文进行日志记录的时候，是否使用二进制的格式记录，默认返回 [LogMsgFormatBinary](2ed2ff87-caec-0e44-f08b-e2b015225a58.htm)，重写可以根据session对象分别返回不同记录模式  Whether to log the current packet in binary format, the default return is [LogMsgFormatBinary](2ed2ff87-caec-0e44-f08b-e2b015225a58.htm). If you want to override it, different recording modes can be returned according to session (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [GetNewNetMessage](c08212c0-06b8-a8c2-2fa9-b1b36ef1212d.htm) | 获取一个新的消息对象的方法，需要在继承类里面进行重写  The method to get a new message object needs to be overridden in the inheritance class (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法 | [GetWordLength](c27821ea-1739-8eb9-a848-23a00828fec8.htm) | 一个字单位的数据表示的地址长度，西门子为2，三菱，欧姆龙，modbusTcp就为1，AB PLC无效  The address length represented by one word of data, Siemens is 2, Mitsubishi, Omron, modbusTcp is 1, AB PLC is invalid (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 受保护的方法代码示例 | [InitializationOnConnect](4ab6989f-a5b5-05fe-8b44-bfa45824a2a7.htm) | 根据实际的协议选择是否重写本方法，有些协议在创建连接之后，需要进行一些初始化的信号握手，才能最终建立网络通道。  Whether to rewrite this method is based on the actual protocol. Some protocols require some initial signal handshake to establish a network channel after the connection is created. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法代码示例 | [InitializationOnConnectAsync](7da2af04-d4fc-5045-19a4-4b0d1722eef3.htm) | 根据实际的协议选择是否重写本方法，有些协议在创建连接之后，需要进行一些初始化的信号握手，才能最终建立网络通道。  Whether to rewrite this method is based on the actual protocol. Some protocols require some initial signal handshake to establish a network channel after the connection is created. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [IpAddressPing](0ac126b4-a830-a822-4ef5-855474942cd7.htm) | 对当前设备的IP地址进行PING的操作，返回PING的结果，正常来说，返回Success  PING the IP address of the current device and return the PING result. Normally, it returns Success (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 受保护的方法 | [LogRevcMessage(Byte)](dab1657f-1e75-7fe1-7003-ccedfc69855d.htm) | 使用日志记录一个接收的报文信息  Logs are used to record information about a received packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [LogRevcMessage(Byte, PipeSession)](0f51ae45-575e-4228-f978-01b37d6d5df1.htm) | 使用日志记录一个接收的报文信息  Logs are used to record information about a received packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [LogSendMessage(Byte)](4f90eecc-68a2-1ef2-44ef-369f19bed3ec.htm) | 使用日志记录一个发送的报文信息  Logs are used to record information about a send packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [LogSendMessage(Byte, PipeSession)](209caadf-4f0f-40aa-04c2-ebc61f8dab13.htm) | 使用日志记录一个发送的报文信息  Logs are used to record information about a send packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法 | [PackCommandWithHeader](ca0c7547-4769-2e5c-b3e5-97e68f1bb00b.htm) | 对当前的命令进行打包处理，通常是携带命令头内容，标记当前的命令的长度信息，需要进行重写，否则默认不打包  The current command is packaged, usually carrying the content of the command header, marking the length of the current command, and it needs to be rewritten, otherwise it is not packaged by default (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法代码示例 | [Read(String, UInt16)](a061c449-5af4-7b69-e6ed-76f1740e178f.htm) | 从Modbus服务器批量读取寄存器的信息，需要指定起始地址，读取长度，如果富文本地址不指定，默认使用的功能码是 0x03，如果需要使用04功能码，那么地址就写成 x=4;100  To read the register information from the Modbus server in batches, you need to specify the start address and read length. If the rich text address is not specified, the default function code is 0x03. If you need to use the 04 function code, the address is written as x = 4; 100 (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法代码示例 | [ReadT](95abe064-5594-5ae7-a849-98c53f88fee0.htm) | 读取支持Hsl特性的数据内容，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明，或是官网的详细文档  Read the data content of the Hsl attribute. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadAsync(String, UInt16)](fb1da34a-99cc-17bc-6e2a-e5155a6b9015.htm) | 从Modbus服务器批量读取寄存器的信息，需要指定起始地址，读取长度，如果富文本地址不指定，默认使用的功能码是 0x03，如果需要使用04功能码，那么地址就写成 x=4;100  To read the register information from the Modbus server in batches, you need to specify the start address and read length. If the rich text address is not specified, the default function code is 0x03. If you need to use the 04 function code, the address is written as x = 4; 100 (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法代码示例 | [ReadAsyncT](0522e740-4bf0-2e31-c4c5-e41384321b37.htm) | 异步读取支持Hsl特性的数据内容，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Asynchronously read the data content of the Hsl attribute. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadBool(String)](9b2faa2e-eaa9-7ace-c49f-7a0318f1df1a.htm) | 读取单个的Boolean数据信息  Read a single Boolean data message (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadBool(String, UInt16)](b9d9664e-4a71-41e9-dd8a-bc6167dcac6c.htm) | 批量读取线圈或是离散的数据信息，需要指定地址和长度，具体的结果取决于实现，如果富文本地址不指定，默认使用的功能码是 0x01  To read coils or discrete data in batches, you need to specify the address and length. The specific result depends on the implementation. If the rich text address is not specified, the default function code is 0x01. (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法 | [ReadBoolAsync(String)](cc37da1f-d673-fe67-78ee-4e465d29a2bd.htm) | 异步读取单个的Boolean数据信息  Asynchronously read a single Boolean data message (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadBoolAsync(String, UInt16)](4ab1379f-3369-4060-6837-413f93bc424a.htm) | 批量读取线圈或是离散的数据信息，需要指定地址和长度，具体的结果取决于实现，如果富文本地址不指定，默认使用的功能码是 0x01  To read coils or discrete data in batches, you need to specify the address and length. The specific result depends on the implementation. If the rich text address is not specified, the default function code is 0x01. (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法 | [ReadByte](910f2447-d09b-877b-9012-9033c8c9acaf.htm) | 按照字节读取汇川M地址的数据，地址示例： MB100，MB101，需要注意的是，MB100 及 MB101 的地址是 MW50 的数据。  Read the data of Inovance M address according to the byte, address example: MB100, MB101, it should be noted that the addresses of MB100 and MB101 are the data of MW50. |
| 公共方法 | [ReadByteAsync](fdd9be24-a813-017b-bc04-25e1356306b4.htm) | 按照字节读取汇川M地址的数据，地址示例： MB100，MB101，需要注意的是，MB100 及 MB101 的地址是 MW50 的数据。  Read the data of Inovance M address according to the byte, address example: MB100, MB101, it should be noted that the addresses of MB100 and MB101 are the data of MW50. |
| 公共方法 | [ReadCoil(String)](9b750e8c-d4e3-2200-8536-7fba9a12b819.htm) | 读取线圈，需要指定起始地址，如果富文本地址不指定，默认使用的功能码是 0x01  To read the coil, you need to specify the start address. If the rich text address is not specified, the default function code is 0x01. (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法 | [ReadCoil(String, UInt16)](4633c27f-4e6f-7ed2-88df-cf6b88256496.htm) | 批量的读取线圈，需要指定起始地址，读取长度，如果富文本地址不指定，默认使用的功能码是 0x01  For batch reading coils, you need to specify the start address and read length. If the rich text address is not specified, the default function code is 0x01. (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法 | [ReadCoilAsync(String)](f744db67-b882-392c-1824-71070c6d19a4.htm) | 读取线圈，需要指定起始地址，如果富文本地址不指定，默认使用的功能码是 0x01  To read the coil, you need to specify the start address. If the rich text address is not specified, the default function code is 0x01. (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法 | [ReadCoilAsync(String, UInt16)](2e8b3480-8188-ed91-cca1-7e53cc641ca2.htm) | 批量的读取线圈，需要指定起始地址，读取长度，如果富文本地址不指定，默认使用的功能码是 0x01  For batch reading coils, you need to specify the start address and read length. If the rich text address is not specified, the default function code is 0x01. (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法代码示例 | [ReadCustomerT(String)](ca069688-d8f0-7dbd-e962-86a10dce83cf.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，返回一个新的类型的实例对象。  To read a custom data type, you need to inherit from the IDataTransfer interface and return an instance object of a new type. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadCustomerT(String, T)](ca7dc1a1-c97d-689f-9819-4dae4af67206.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，传入一个实例，对这个实例进行赋值，并返回该实例的对象。  To read a custom data type, you need to inherit from the IDataTransfer interface, pass in an instance, assign a value to this instance, and return the object of the instance. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadCustomerAsyncT(String)](f3fcda2b-255d-aadd-567e-c7b3df3636d9.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，返回一个新的类型的实例对象。  To read a custom data type, you need to inherit from the IDataTransfer interface and return an instance object of a new type. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadCustomerAsyncT(String, T)](527c9fae-4743-6ebd-7b9a-5ef3a8b7c455.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，传入一个实例，对这个实例进行赋值，并返回该实例的对象。  To read a custom data type, you need to inherit from the IDataTransfer interface, pass in an instance, assign a value to this instance, and return the object of the instance. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadDiscrete(String)](4886013f-4945-b8b5-7a85-8fd979367a5b.htm) | 读取输入线圈，需要指定起始地址，如果富文本地址不指定，默认使用的功能码是 0x02  To read the input coil, you need to specify the start address. If the rich text address is not specified, the default function code is 0x02. (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法 | [ReadDiscrete(String, UInt16)](f907fa38-a034-8af3-5b1d-99ccd6062667.htm) | 批量的读取输入点，需要指定起始地址，读取长度，如果富文本地址不指定，默认使用的功能码是 0x02  To read input points in batches, you need to specify the start address and read length. If the rich text address is not specified, the default function code is 0x02 (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法 | [ReadDiscreteAsync(String)](abc82d59-8329-5160-b68d-e19588fe4a20.htm) | 读取输入线圈，需要指定起始地址，如果富文本地址不指定，默认使用的功能码是 0x02  To read the input coil, you need to specify the start address. If the rich text address is not specified, the default function code is 0x02. (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法 | [ReadDiscreteAsync(String, UInt16)](70c766bd-2a88-287f-7c06-c2da7870f010.htm) | 批量的读取输入点，需要指定起始地址，读取长度，如果富文本地址不指定，默认使用的功能码是 0x02  To read input points in batches, you need to specify the start address and read length. If the rich text address is not specified, the default function code is 0x02 (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法代码示例 | [ReadDouble(String)](e7e31c43-71b4-e692-6bdc-ba9533e1d83d.htm) | 读取双浮点的数据  Read double floating point data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadDouble(String, UInt16)](98be3669-7b4b-fa15-efd7-12cdda2a1ac2.htm) | 读取双浮点数据的数组  Read double floating point data array (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法代码示例 | [ReadDoubleAsync(String)](2cf895fa-1033-1def-f048-0e8caefaa71d.htm) | 异步读取双浮点的数据  Asynchronously read double floating point data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadDoubleAsync(String, UInt16)](82fc6208-2830-7463-9003-fd7a33ea6dab.htm) | 异步读取双浮点数据的数组  Asynchronously read double floating point data array (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法 | [ReadFile](6ff35712-0972-5094-2aee-05f62228db24.htm) | 使用0x14功能码读取文件的数据内容，需要指定文件编号，地址，长度，返回读取的结果数据，读取的长度为任意长度，内部自动根据实际情况来切割操作。 (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法代码示例 | [ReadFloat(String)](bc2801ae-6600-e845-7c6c-522c9156eb0e.htm) | 读取单浮点数据  Read single floating point data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadFloat(String, UInt16)](de69a6c8-ee84-464e-d735-aa479a2fbe4f.htm) | 读取单浮点精度的数组  Read single floating point array (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法代码示例 | [ReadFloatAsync(String)](89f743f8-4b3a-7085-1e49-a1d7adb10052.htm) | 异步读取单浮点数据  Asynchronously read single floating point data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadFloatAsync(String, UInt16)](ffac5e66-8ba9-5e95-b6f9-21ac11013093.htm) | 异步读取单浮点精度的数组  Asynchronously read single floating point array (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法 | [ReadFromCoreServer(IEnumerableByte)](044cfdd9-ec76-b89f-2f2b-5f439025006c.htm) | 将多个数据报文按顺序发到设备，并从设备接收返回的数据内容，然后拼接成一个Byte[]信息，需要重写[UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm)方法才能返回正确的结果。  Send multiple data packets to the device in sequence, and receive the returned data content from the device, and then splicing them into a Byte[] message, you need to rewrite [UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm) method to return the correct result. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServer(Byte)](b516dde4-8446-eaaa-f60c-35500a2c9a25.htm) | 将二进制的数据发送到管道中去，然后从管道里接收二进制的数据回来，并返回是否成功的结果对象。  Send binary data to the pipeline, and then receive binary data back from the pipeline, and return whether the success of the result object (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServer(Byte, Boolean, Boolean)](b32ce8e8-a959-fcc2-6736-14b39fdb9a40.htm) | 将二进制的数据发送到管道中去，然后从管道里接收二进制的数据回来，并返回是否成功的结果对象。  Send binary data to the pipeline, and then receive binary data back from the pipeline, and return whether the success of the result object (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServer(CommunicationPipe, Byte, Boolean, Boolean)](ef36d56b-df09-f727-e8d3-316d37a13ee9.htm) | 使用指定的管道来进行数据通信，发送原始数据到管道，然后从管道接收相关的数据返回，本方法无锁 (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServerAsync(IEnumerableByte)](bcdb59b7-a70a-d5c0-1fc6-9ececcd39ea2.htm) | 将多个数据报文按顺序发到设备，并从设备接收返回的数据内容，然后拼接成一个Byte[]信息，需要重写[UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm)方法才能返回正确的结果。  Send multiple data packets to the device in sequence, and receive the returned data content from the device, and then splicing them into a Byte[] message, you need to rewrite [UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm) method to return the correct result. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServerAsync(Byte)](90b7ccb6-1139-1b6f-edd1-d5d0b845eb6c.htm) | 将二进制的数据发送到管道中去，然后从管道里接收二进制的数据回来，并返回是否成功的结果对象。  Send binary data to the pipeline, and then receive binary data back from the pipeline, and return whether the success of the result object (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法代码示例 | [ReadFromCoreServerAsync(Byte, Boolean, Boolean)](1c79c204-5ab5-9342-e65d-7455c20ea577.htm) | 将二进制的数据发送到管道中去，然后从管道里接收二进制的数据回来，并返回是否成功的结果对象。  Send binary data to the pipeline, and then receive binary data back from the pipeline, and return whether the success of the result object (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServerAsync(CommunicationPipe, Byte, Boolean, Boolean)](a254a19d-3c29-147b-335b-8adccdf4ee2f.htm) | 使用指定的管道来进行数据通信，发送原始数据到管道，然后从管道接收相关的数据返回，本方法无锁 (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法代码示例 | [ReadInt16(String)](57a16003-2fd5-b2bd-09bf-14366be899b8.htm) | 读取16位的有符号的整型数据  Read 16-bit signed integer data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt16(String, UInt16)](b2bde158-f51c-93fc-66ca-5e4cf6ea1fb7.htm) | 读取16位的有符号整型数组  Read 16-bit signed integer array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt16Async(String)](ac021793-890d-ceab-eae4-47ff763390c1.htm) | 异步读取16位的有符号的整型数据  Asynchronously read 16-bit signed integer data (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt16Async(String, UInt16)](f30e02fc-8078-f1b3-c226-5991e3e23caf.htm) | 异步读取16位的有符号整型数组  Asynchronously read 16-bit signed integer array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt32(String)](bfbea035-9a83-40e3-084c-28122041b37a.htm) | 读取32位的有符号整型  Read 32-bit signed integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt32(String, UInt16)](44156b71-97d1-917a-782a-5e716ffd286c.htm) | 读取32位有符号整型数组  Read 32-bit signed integer array (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法代码示例 | [ReadInt32Async(String)](96726ea0-82aa-1229-21f6-43cd62abca43.htm) | 异步读取32位的有符号整型  Asynchronously read 32-bit signed integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt32Async(String, UInt16)](4957686b-6205-802f-177e-b5322e1eabda.htm) | 异步读取32位有符号整型数组  Asynchronously read 32-bit signed integer array (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法代码示例 | [ReadInt64(String)](0d0c2284-2da1-6f9a-b353-537561cd5e95.htm) | 读取64位的有符号整型  Read 64-bit signed integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt64(String, UInt16)](2a1dd46e-c6be-563e-6ae7-86d7432b64f6.htm) | 读取64位的有符号整型数组  Read 64-bit signed integer array (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法代码示例 | [ReadInt64Async(String)](26d46c8a-3775-3abd-d3ec-5cd1758a317e.htm) | 异步读取64位的有符号整型  Asynchronously read 64-bit signed integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadInt64Async(String, UInt16)](5a7ebc54-70a1-50a1-6c3a-9c4b3c4f9af2.htm) | 异步读取64位的有符号整型数组  Asynchronously read 64-bit signed integer array (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法代码示例 | [ReadString(String, UInt16)](07812d76-415a-6400-c8a2-1c21da51e7a3.htm) | 读取字符串数据，默认为最常见的ASCII编码  Read string data, default is the most common ASCII encoding (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadString(String, UInt16, Encoding)](32a0448f-89be-a39d-4053-f7c4f9309bf4.htm) | 使用指定的编码，读取字符串数据  Reads string data using the specified encoding (重写 [DeviceCommunicationReadString(String, UInt16, Encoding)](b34c71a6-58f3-6bf4-0b93-471be4f2ed31.htm).) |
| 公共方法代码示例 | [ReadStringAsync(String, UInt16)](1835da81-b811-d151-0b53-907d8d87d0a8.htm) | 异步读取字符串数据，默认为最常见的ASCII编码  Asynchronously read string data, default is the most common ASCII encoding (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadStringAsync(String, UInt16, Encoding)](ba19ecc8-7d44-2476-c633-c83e8cf5223b.htm) | 异步使用指定的编码，读取字符串数据  Asynchronously reads string data using the specified encoding (重写 [DeviceCommunicationReadStringAsync(String, UInt16, Encoding)](091dfc67-3d25-f33e-3a55-a0228863815e.htm).) |
| 公共方法代码示例 | [ReadStructT](49853082-fbad-52e5-d756-0615cedb4b83.htm) | 读取结构体类型的数据，根据结构体自身的定义，读取原始字节数组，然后解析出实际的结构体数据，结构体需要实现[HslStructAttribute](87e427a8-0c7b-43a8-c421-c27264eda90f.htm)特性 (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadStructAsyncT](36461e61-d0ed-62cf-79ce-dcc7f71a3fa2.htm) | 读取结构体类型的数据，根据结构体自身的定义，读取原始字节数组，然后解析出实际的结构体数据，结构体需要实现[HslStructAttribute](87e427a8-0c7b-43a8-c421-c27264eda90f.htm)特性 (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt16(String)](736379d1-9560-5052-658a-13c126d56d78.htm) | 读取16位的无符号整型  Read 16-bit unsigned integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt16(String, UInt16)](e284c4bd-d3f8-ff30-eb42-0aa035cfb4f6.htm) | 读取16位的无符号整型数组  Read 16-bit unsigned integer array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt16Async(String)](361acc93-43b5-a9a5-6fa7-fa8728b7b96c.htm) | 异步读取16位的无符号整型  Asynchronously read 16-bit unsigned integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt16Async(String, UInt16)](73dd24b8-b501-df53-1de6-ccb2708f656d.htm) | 异步读取16位的无符号整型数组  Asynchronously read 16-bit unsigned integer array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt32(String)](ca73c958-39c7-9b34-aa44-3db2d9b0cad3.htm) | 读取32位的无符号整型  Read 32-bit unsigned integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt32(String, UInt16)](c5543613-3abd-2cb0-dc1d-eb89b3da3a1e.htm) | 读取32位的无符号整型数组  Read 32-bit unsigned integer array (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法代码示例 | [ReadUInt32Async(String)](2343953f-6350-331e-6ba7-18a4449f0575.htm) | 异步读取32位的无符号整型  Asynchronously read 32-bit unsigned integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt32Async(String, UInt16)](272c5014-a400-5e74-69f5-1de1237e4ae0.htm) | 异步读取32位的无符号整型数组  Asynchronously read 32-bit unsigned integer array (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法代码示例 | [ReadUInt64(String)](947487dd-d4a3-72d8-d9fb-2224aad70e12.htm) | 读取64位的无符号整型  Read 64-bit unsigned integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt64(String, UInt16)](f3a9d180-d7bc-aacc-b07d-0edc4608547c.htm) | 读取64位的无符号整型的数组  Read 64-bit unsigned integer array (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法代码示例 | [ReadUInt64Async(String)](3411ff3a-772f-dbe1-8931-75e00a731579.htm) | 异步读取64位的无符号整型  Asynchronously read 64-bit unsigned integer (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadUInt64Async(String, UInt16)](20fab556-ee24-9f7b-9a2c-dfdbd2b2dca4.htm) | 异步读取64位的无符号整型的数组  Asynchronously read 64-bit unsigned integer array (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法 | [ReadWrite](4d80e7e8-42a3-5005-e8f0-5dd028431ccd.htm) | 使用0x17功能码来实现同时写入并读取数据的操作，使用一条报文来实现，需要指定读取的地址，长度，写入的地址，写入的数据信息，返回读取的结果数据。  Use 0x17 function code to write and read data at the same time, and use a message to implement it, you need to specify the read address, length, written address, written data information, and return the read result data. (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法 | [RegisteredAddressMapping](5686f821-9dcb-c44c-5a99-663d0b75748a.htm) | 注册一个新的地址映射关系，注册地址映射关系后，就可以使用新的地址来读写Modbus数据了，通常用于其他的支持Modbus协议的PLC。  After registering a new address mapping, you can use the new address to read and write Modbus data, which is usually used for other PLCs that support the Modbus protocol. (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法 | [SetDtuPipe](f765d1d7-234a-b229-9457-97e25e024dc9.htm) | 将当前的通信对象设置DTU模式，允许传入现成的管道，并返回初始化结果，如果该设备重写了握手报文，就是返回握手结果  Set the current communication object to DTU mode, allow the existing pipe to be passed in, and return the initialization result, if the device rewrites the handshake packet, the handshake result is returned (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [SetDtuPipeAsync](9d1fd12d-b804-3828-a627-a2c0eef1d23f.htm) | 将当前的通信对象设置DTU模式，允许传入现成的管道，并返回初始化结果，如果该设备重写了握手报文，就是返回握手结果  Set the current communication object to DTU mode, allow the existing pipe to be passed in, and return the initialization result, if the device rewrites the handshake packet, the handshake result is returned (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [SetPersistentConnection](f903dae0-4a4e-c7db-8350-50e36b257ec5.htm) | **已过时。** V11版本及之前设置长连接的方法，在V12版本以上中没有任何效果，默认长连接，删除调用即可，此处保留方法是为了部分用户保持兼容性升级。  The method of setting the long connection in V11 and before, has no effect in V12 and above. this method can be deleted. The method is retained here to maintain compatibility upgrades for some users. (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共方法 | [ToString](70ff3d60-984c-329c-201b-5de7d81d4077.htm) | (重写 [ModbusRtuOverTcpToString](c409cc5d-0b8d-7cab-dc72-b323757b8449.htm).) |
| 公共方法 | [TranslateToModbusAddress](b42ef779-c21b-3501-d4ec-55535980a43c.htm) | 将当前的地址信息转换成Modbus格式的地址，如果转换失败，返回失败的消息。默认不进行任何的转换。  Convert the current address information into a Modbus format address. If the conversion fails, a failure message will be returned. No conversion is performed by default. (重写 [ModbusRtuOverTcpTranslateToModbusAddress(String, Byte)](fd4a9aee-b09c-08e3-dff3-51a3dbdd6087.htm).) |
| 公共方法 | [UnpackResponseContent](4dc1fa2e-6465-742e-bf0f-f133a6686523.htm) | 根据对方返回的报文命令，对命令进行基本的拆包，例如各种Modbus协议拆包为统一的核心报文，还支持对报文的验证  According to the message command returned by the other party, the command is basically unpacked, for example, various Modbus protocols are unpacked into a unified core message, and the verification of the message is also supported (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
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
| 公共方法代码示例 | [Write(String, Double)](64c44576-7e95-cbcb-a8c9-e6b7ff7e3c50.htm) | 写入double数据，返回是否成功  Write double data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Int16)](c7cfd67d-ca31-3087-c684-ee4c10780ab1.htm) | 写入short数组，返回是否成功  Write short array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Int32)](ba80d1a5-c8d8-7e0b-5c6b-9435cf20cadb.htm) | 写入int数据，返回是否成功  Write int data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Int64)](ba246b33-5a12-a22c-6d40-0969a64b508b.htm) | 写入long数据，返回是否成功  Write long data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, Single)](a2265fa7-5b56-a5aa-9d21-a0dc5458cbf9.htm) | 写入float数据，返回是否成功  Write float data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, String)](4e0ec363-ba64-7ca7-b8d3-2e9a8087bd55.htm) | 写入字符串信息，编码为ASCII  Write string information, encoded as ASCII (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, UInt16)](65544962-82ac-bb47-d72f-5dd0c3d4ac6c.htm) | 写入ushort数组，返回是否成功  Write ushort array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, UInt32)](c99d8246-f95e-4f79-738b-16e622072eb4.htm) | 写入uint数据，返回是否成功  Write uint data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, UInt64)](4924086e-9b36-d446-4ae4-6682134a6ad1.htm) | 写入ulong数据，返回是否成功  Write ulong data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [Write(String, Boolean)](0be65147-ace7-ba53-504d-d4b923a30333.htm) | 向线圈中写入bool数值，返回是否写入成功，如果富文本地址不指定，默认使用的功能码是 0x05， 如果你的地址为字地址，例如100.2，那么将使用0x16的功能码，通过掩码的方式来修改寄存器的某一位，需要Modbus服务器支持，否则自动切换为读取字数据，修改位，在写入字的方式。  Write bool value to the coil and return whether the writing is successful. If the rich text address is not specified, the default function code is 0x05. If your address is a word address, such as 100.2, then you will use the function code of 0x16 to modify a bit of the register through a mask. It needs Modbus server support, Otherwise, it automatically switches to read the word data, modifies the bits, and writes the word in the way. (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法 | [Write(String, Boolean)](4fe07f04-9b89-7291-d2d7-1f474293f4de.htm) | 向线圈中写入bool数组，返回是否写入成功，如果富文本地址不指定，默认使用的功能码是 0x0F  Write the bool array to the coil, and return whether the writing is successful. If the rich text address is not specified, the default function code is 0x0F. (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法代码示例 | [Write(String, Byte)](fdfef281-9b77-397a-d6df-078cd26afd17.htm) | 将数据写入到Modbus的寄存器上去，需要指定起始地址和数据内容，如果富文本地址不指定，默认使用的功能码是 0x10  To write data to Modbus registers, you need to specify the start address and data content. If the rich text address is not specified, the default function code is 0x10 (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法代码示例 | [Write(String, Double)](4d6f1238-0d16-f592-924e-c35cbd755e83.htm) | 写入double数组，返回是否成功  Write double array, return whether the write was successful (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法 | [Write(String, Int16)](77d3ad67-4385-bef7-4223-cc2340afbb7e.htm) | 将数据写入到Modbus的单个寄存器上去，需要指定起始地址和数据值，如果富文本地址不指定，默认使用的功能码是 0x06  To write data to a single register of Modbus, you need to specify the start address and data value. If the rich text address is not specified, the default function code is 0x06. (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法代码示例 | [Write(String, Int32)](ed9597f2-87de-c1ef-4b83-4267527a38c9.htm) | 写入int[]数组，返回是否成功  Write int array, return whether the write was successful (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法代码示例 | [Write(String, Int64)](72aef2fe-26dd-dc40-55df-9fa4aa2ae6bf.htm) | 写入long数组，返回是否成功  Write long array, return whether the write was successful (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法代码示例 | [Write(String, Single)](9e904462-199f-488f-2e9f-3c972ca9955e.htm) | 写入float数组，返回是否成功  Write float array, return whether the write was successful (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法 | [Write(String, UInt16)](a3316efc-b1cc-2b7f-ece1-f945b3b0182d.htm) | 将数据写入到Modbus的单个寄存器上去，需要指定起始地址和数据值，如果富文本地址不指定，默认使用的功能码是 0x06  To write data to a single register of Modbus, you need to specify the start address and data value. If the rich text address is not specified, the default function code is 0x06. (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法代码示例 | [Write(String, UInt32)](3f968111-808b-5681-7126-74765225c918.htm) | 写入uint[]数组，返回是否成功  Write uint array, return whether the write was successful (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法代码示例 | [Write(String, UInt64)](a0df0554-5053-ebed-bfab-5f8a3c331a48.htm) | 写入ulong数组，返回是否成功  Write ulong array, return whether the write was successful (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法代码示例 | [Write(String, String, Int32)](02fe2ca9-49a3-ede4-d1af-8c6da1f644da.htm) | 写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为ASCII  Write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, it is filled with 0 and the encoding is ASCII. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, String, Encoding)](7b602661-902d-2b70-1fdb-ba87ecd307bf.htm) | 写入字符串信息，需要指定的编码信息  Write string information, need to specify the encoding information (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, String, Int32, Encoding)](e8cb30d9-4c5e-6c68-b9af-1b1b6c2dc714.htm) | 写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为指定的编码信息  Write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, then the operation is complemented with 0 , you should specified the encoding information (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteT(T)](9f48853f-fbeb-c37d-25a2-1c76d1dee407.htm) | 写入支持Hsl特性的数据，返回是否写入成功，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Write data that supports the Hsl attribute, and return whether the write was successful. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Double)](068adaa4-d1de-dffc-8376-6f483e2c1d96.htm) | 异步写入double数据，返回是否成功  Asynchronously write double data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int16)](2625e868-7096-e526-ce25-b29afc639b41.htm) | 异步写入short数组，返回是否成功  Asynchronously write short array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int32)](df91d629-2080-328d-67cd-fbefb9a45404.htm) | 异步写入int数据，返回是否成功  Asynchronously write int data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int64)](ec6c1e24-a0d8-f1f5-56af-109b0c546fc8.htm) | 异步写入long数据，返回是否成功  Asynchronously write long data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Single)](910189df-13f9-b7d0-a626-b5b25ce5d32a.htm) | 异步写入float数据，返回是否成功  Asynchronously write float data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, String)](fe839c9a-46f3-e317-0155-b10a8a83a9b1.htm) | 异步写入字符串信息，编码为ASCII  Asynchronously write string information, encoded as ASCII (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt16)](34a1a321-972e-06a8-76fb-794dd1229005.htm) | 异步写入ushort数组，返回是否成功  Asynchronously write ushort array, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt32)](abad215c-bc99-ac90-b9f3-7adb3a7cca8f.htm) | 异步写入uint数据，返回是否成功  Asynchronously write uint data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt64)](e396074d-6cf1-b45d-5e8e-01026f6e391b.htm) | 异步写入ulong数据，返回是否成功  Asynchronously write ulong data, return whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WriteAsync(String, Boolean)](62f1700b-e043-2ec6-f2aa-47aa41548cec.htm) | 向线圈中写入bool数值，返回是否写入成功，如果富文本地址不指定，默认使用的功能码是 0x05， 如果你的地址为字地址，例如100.2，那么将使用0x16的功能码，通过掩码的方式来修改寄存器的某一位，需要Modbus服务器支持，否则自动切换为读取字数据，修改位，在写入字的方式。  Write bool value to the coil and return whether the writing is successful. If the rich text address is not specified, the default function code is 0x05. If your address is a word address, such as 100.2, then you will use the function code of 0x16 to modify a bit of the register through a mask. It needs Modbus server support, Otherwise, it automatically switches to read the word data, modifies the bits, and writes the word in the way. (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法 | [WriteAsync(String, Boolean)](dbc7b234-d200-d7e8-e308-fed865d7aa95.htm) | 向线圈中写入bool数组，返回是否写入成功，如果富文本地址不指定，默认使用的功能码是 0x0F  Write the bool array to the coil, and return whether the writing is successful. If the rich text address is not specified, the default function code is 0x0F. (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Byte)](104f3bf6-6822-f73c-3419-e8c99d5a3c15.htm) | 将数据写入到Modbus的寄存器上去，需要指定起始地址和数据内容，如果富文本地址不指定，默认使用的功能码是 0x10  To write data to Modbus registers, you need to specify the start address and data content. If the rich text address is not specified, the default function code is 0x10 (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Double)](910251ea-28ae-052b-b95b-c6caf8fe6483.htm) | 异步写入double数组，返回是否成功  Asynchronously write double array, return whether the write was successful (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法 | [WriteAsync(String, Int16)](17151d4f-bde5-cab5-e749-0992ad67de8b.htm) | 将数据写入到Modbus的单个寄存器上去，需要指定起始地址和数据值，如果富文本地址不指定，默认使用的功能码是 0x06  To write data to a single register of Modbus, you need to specify the start address and data value. If the rich text address is not specified, the default function code is 0x06. (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int32)](7bdf041e-af53-261c-ec07-0b1b85c1c735.htm) | 异步写入int[]数组，返回是否成功  Asynchronously write int array, return whether the write was successful (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int64)](d8069128-a5ea-a0c1-2221-3b151b626c53.htm) | 异步写入long数组，返回是否成功  Asynchronously write long array, return whether the write was successful (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Single)](2af0147b-de42-f977-b99f-e91f585e0c9a.htm) | 异步写入float数组，返回是否成功  Asynchronously write float array, return whether the write was successful (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法 | [WriteAsync(String, UInt16)](d08d6403-4da8-13fe-3762-86932dab62de.htm) | 将数据写入到Modbus的单个寄存器上去，需要指定起始地址和数据值，如果富文本地址不指定，默认使用的功能码是 0x06  To write data to a single register of Modbus, you need to specify the start address and data value. If the rich text address is not specified, the default function code is 0x06. (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt32)](eb068ca1-b343-5288-05a9-f8481d323bdf.htm) | 异步写入uint[]数组，返回是否成功  Asynchronously write uint array, return whether the write was successful (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt64)](e19ff622-f678-4c52-1662-2af481317995.htm) | 异步写入ulong数组，返回是否成功  Asynchronously write ulong array, return whether the write was successful (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法 | [WriteAsync(String, String, Int32)](b739b983-b03c-8ecf-2713-75f19233ae2f.htm) | 异步写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为ASCII  Asynchronously write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, it is filled with 0 and the encoding is ASCII. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, String, Encoding)](876e9800-0406-6407-bdf8-8ab1cf28df38.htm) | 异步写入字符串信息，需要指定的编码信息  Asynchronously write string information, need to specify the encoding information (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WriteAsync(String, String, Int32, Encoding)](ae9fb2a3-fddb-48eb-e9ad-6cc0ebbac1f6.htm) | 异步写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为指定的编码信息  Asynchronously write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, then the operation is complemented with 0 , you should specified the encoding information (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsyncT(T)](25d96e77-d1b2-918c-33f9-e5e4cc57c9cd.htm) | 异步写入支持Hsl特性的数据，返回是否写入成功，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Asynchronously write data that supports the Hsl attribute, and return whether the write was successful. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteCustomerT](84a37f82-3f67-a0f8-fd3f-0630772565e4.htm) | 写入自定义类型的数据，该类型必须继承自IDataTransfer接口  Write data of a custom type, which must inherit from the IDataTransfer interface (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteCustomerAsyncT](04628397-cfa1-d785-7790-4c02cb377ab6.htm) | 写入自定义类型的数据，该类型必须继承自IDataTransfer接口  Write data of a custom type, which must inherit from the IDataTransfer interface (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WriteFile](c888070e-6866-6228-2d68-533c2b238fd8.htm) | 使用0x15功能码写入文件的数据内容，需要指定文件编号，地址，数据，返回写入的结果数据。写入的数据字节长度不能大于248个字节，如果想写入更多的字节，就需要切片写入操作 (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法 | [WriteMask](5cce8a41-41b1-90d8-b44e-07c85b38041a.htm) | 向设备写入掩码数据，使用0x16功能码，需要确认对方是否支持相关的操作，掩码数据的操作主要针对寄存器。  To write mask data to the server, using the 0x16 function code, you need to confirm whether the other party supports related operations. The operation of mask data is mainly directed to the register. (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法 | [WriteMaskAsync](24048bc4-090d-04c2-b762-85ac5c836c8d.htm) | 向设备写入掩码数据，使用0x16功能码，需要确认对方是否支持相关的操作，掩码数据的操作主要针对寄存器。  To write mask data to the server, using the 0x16 function code, you need to confirm whether the other party supports related operations. The operation of mask data is mainly directed to the register. (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法 | [WriteOneRegister(String, Int16)](685792aa-f82c-7400-8319-d62d56af0aba.htm) | 将数据写入到Modbus的单个寄存器上去，需要指定起始地址和数据值，如果富文本地址不指定，默认使用的功能码是 0x06  To write data to a single register of Modbus, you need to specify the start address and data value. If the rich text address is not specified, the default function code is 0x06. (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法 | [WriteOneRegister(String, UInt16)](38dca771-988c-7bbb-1963-8c354830cc3f.htm) | 将数据写入到Modbus的单个寄存器上去，需要指定起始地址和数据值，如果富文本地址不指定，默认使用的功能码是 0x06  To write data to a single register of Modbus, you need to specify the start address and data value. If the rich text address is not specified, the default function code is 0x06. (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法 | [WriteOneRegisterAsync(String, Int16)](7c9fb70d-655c-0544-fa34-086ffe8ad357.htm) | 将数据写入到Modbus的单个寄存器上去，需要指定起始地址和数据值，如果富文本地址不指定，默认使用的功能码是 0x06  To write data to a single register of Modbus, you need to specify the start address and data value. If the rich text address is not specified, the default function code is 0x06. (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |
| 公共方法 | [WriteOneRegisterAsync(String, UInt16)](5d3d62d9-e9fc-a807-36d6-5771b9bac661.htm) | 将数据写入到Modbus的单个寄存器上去，需要指定起始地址和数据值，如果富文本地址不指定，默认使用的功能码是 0x06  To write data to a single register of Modbus, you need to specify the start address and data value. If the rich text address is not specified, the default function code is 0x06. (继承自 [ModbusRtuOverTcp](122ebebd-f8a8-4214-b165-280215a9543e.htm)。) |

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

AM400\_800 的元件有 Q 区，I 区，M 区这三种，分别都可以按位，按字节，按字和按双字进行访问，在本组件的条件下，仅支持按照位，字访问。  
位地址支持 Q, I, M 地址类型，字地址支持 SM, SD，支持对字地址的位访问，例如 ReadBool("SD0.5");  
H3U 系列控制器支持 M/SM/S/T/C/X/Y 等 bit 型变量（也称线圈） 的访问、 D/SD/R/T/C 等 word 型变量的访问；  
H5U 系列控制器支持 M/B/S/X/Y 等 bit 型变量（也称线圈） 的访问、 D/R 等 word 型变量的访问；内部 W 元件，不支持通信访问。

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.Profinet.Inovance 命名空间](aa61d757-9239-26ce-525f-b7e52699e229.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## InovanceSerialOverTcp 构造函数 

[原文連結](http://api.hslcommunication.cn/html/c346256c-782a-1381-f646-c6221355f534.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Inovance](../html/aa61d757-9239-26ce-525f-b7e52699e229.htm "HslCommunication.Profinet.Inovance")

[InovanceSerialOverTcp 类](../html/f82b954d-1a19-0b0e-f705-677d531e17f6.htm "InovanceSerialOverTcp 类")

[InovanceSerialOverTcp 构造函数](../html/c346256c-782a-1381-f646-c6221355f534.htm "InovanceSerialOverTcp 构造函数 ")

[InovanceSerialOverTcp 构造函数](../html/523bbb48-d091-524a-3b2a-7e06104a0664.htm "InovanceSerialOverTcp 构造函数 ")

[InovanceSerialOverTcp 构造函数 (String, Int32, Byte)](../html/160107db-161f-966a-1751-bf7fcd70572a.htm "InovanceSerialOverTcp 构造函数 (String, Int32, Byte)")

[InovanceSerialOverTcp 构造函数 (InovanceSeries, String, Int32, Byte)](../html/b65bbcec-6ea5-2f5e-ae3c-570cc11ef158.htm "InovanceSerialOverTcp 构造函数 (InovanceSeries, String, Int32, Byte)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| InovanceSerialOverTcp 构造函数 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [InovanceSerialOverTcp](523bbb48-d091-524a-3b2a-7e06104a0664.htm) | 实例化一个默认的对象 |
| 公共方法 | [InovanceSerialOverTcp(String, Int32, Byte)](160107db-161f-966a-1751-bf7fcd70572a.htm) | 通过指定站号，ip地址，端口号来实例化一个新的对象 |
| 公共方法 | [InovanceSerialOverTcp(InovanceSeries, String, Int32, Byte)](b65bbcec-6ea5-2f5e-ae3c-570cc11ef158.htm) | 通过指定站号，IP地址，端口以及PLC的系列来实例化一个新的对象  Instantiate a new object by specifying the station number and PLC series |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[InovanceSerialOverTcp 类](f82b954d-1a19-0b0e-f705-677d531e17f6.htm)

[HslCommunication.Profinet.Inovance 命名空间](aa61d757-9239-26ce-525f-b7e52699e229.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## InovanceSerialOverTcp 构造函数 

[原文連結](http://api.hslcommunication.cn/html/523bbb48-d091-524a-3b2a-7e06104a0664.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Inovance](../html/aa61d757-9239-26ce-525f-b7e52699e229.htm "HslCommunication.Profinet.Inovance")

[InovanceSerialOverTcp 类](../html/f82b954d-1a19-0b0e-f705-677d531e17f6.htm "InovanceSerialOverTcp 类")

[InovanceSerialOverTcp 构造函数](../html/c346256c-782a-1381-f646-c6221355f534.htm "InovanceSerialOverTcp 构造函数 ")

[InovanceSerialOverTcp 构造函数](../html/523bbb48-d091-524a-3b2a-7e06104a0664.htm "InovanceSerialOverTcp 构造函数 ")

[InovanceSerialOverTcp 构造函数 (String, Int32, Byte)](../html/160107db-161f-966a-1751-bf7fcd70572a.htm "InovanceSerialOverTcp 构造函数 (String, Int32, Byte)")

[InovanceSerialOverTcp 构造函数 (InovanceSeries, String, Int32, Byte)](../html/b65bbcec-6ea5-2f5e-ae3c-570cc11ef158.htm "InovanceSerialOverTcp 构造函数 (InovanceSeries, String, Int32, Byte)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| InovanceSerialOverTcp 构造函数 |

实例化一个默认的对象

**命名空间：**
 [HslCommunication.Profinet.Inovance](aa61d757-9239-26ce-525f-b7e52699e229.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public InovanceSerialOverTcp()
```

```
Public Sub New
```

```
public:
InovanceSerialOverTcp()
```

```
new : unit -> InovanceSerialOverTcp
```

![](../icons/SectionExpanded.png)参见

#### 引用

[InovanceSerialOverTcp 类](f82b954d-1a19-0b0e-f705-677d531e17f6.htm)

[InovanceSerialOverTcp 重载](c346256c-782a-1381-f646-c6221355f534.htm)

[HslCommunication.Profinet.Inovance 命名空间](aa61d757-9239-26ce-525f-b7e52699e229.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)