# HslCommunication - HslCommunication.Profinet.Yokogawa

> 分類頁數: 30



---
## HslCommunication.Profinet.Yokogawa

[原文連結](http://api.hslcommunication.cn/html/e985e513-76fd-7ba4-62ba-2b2b3d61b929.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Yokogawa](../html/e985e513-76fd-7ba4-62ba-2b2b3d61b929.htm "HslCommunication.Profinet.Yokogawa")

[YokogawaLinkHelper 类](../html/6097616d-a102-a94a-1bd2-b7a00819a5fa.htm "YokogawaLinkHelper 类")

[YokogawaLinkServer 类](../html/ab5927f5-b6a1-d76a-7c06-10cbcca7ac37.htm "YokogawaLinkServer 类")

[YokogawaLinkTcp 类](../html/31845cf4-eb9e-f801-f6b8-c2759f9b9a5b.htm "YokogawaLinkTcp 类")

[YokogawaSystemInfo 类](../html/b2464ddd-7f4e-705a-a5e3-561a10b040fe.htm "YokogawaSystemInfo 类")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslCommunication.Profinet.Yokogawa 命名空间 |

[缺少 "N:HslCommunication.Profinet.Yokogawa" 的 <summary> 文档]

![](../icons/SectionExpanded.png)类

|  | 类 | 说明 |
| --- | --- | --- |
| 公共类 | [YokogawaLinkHelper](6097616d-a102-a94a-1bd2-b7a00819a5fa.htm) | 横河PLC的通信辅助类。 |
| 公共类代码示例 | [YokogawaLinkServer](ab5927f5-b6a1-d76a-7c06-10cbcca7ac37.htm) | 横河PLC的虚拟服务器，支持X,Y,I,E,M,T,C,L继电器类型的数据读写，支持D,B,F,R,V,Z,W,TN,CN寄存器类型的数据读写，可以用来测试横河PLC的二进制通信类型  Yokogawa PLC's virtual server, supports X, Y, I, E, M, T, C, L relay type data read and write, supports D, B, F, R, V, Z, W, TN, CN register types The data read and write can be used to test the binary communication type of Yokogawa PLC |
| 公共类代码示例 | [YokogawaLinkTcp](31845cf4-eb9e-f801-f6b8-c2759f9b9a5b.htm) | 横河PLC的二进制通信类，支持X,Y,I,E,M,T,C,L继电器类型的数据读写，支持D,B,F,R,V,Z,W,TN,CN寄存器类型的数据读写，还支持一些高级的信息读写接口，详细参考Demo界面。  Yokogawa PLC's binary communication type, supports X, Y, I, E, M, T, C, L relay type data read and write, supports D, B, F, R, V, Z, W, TN, CN registers Types of data reading and writing, and some advanced information reading and writing interfaces are also supported. Please refer to the Demo UI. |
| 公共类 | [YokogawaSystemInfo](b2464ddd-7f4e-705a-a5e3-561a10b040fe.htm) | 横河PLC的系统基本信息  Basic system information of Yokogawa PLC |

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## YokogawaLinkHelper 类

[原文連結](http://api.hslcommunication.cn/html/6097616d-a102-a94a-1bd2-b7a00819a5fa.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Yokogawa](../html/e985e513-76fd-7ba4-62ba-2b2b3d61b929.htm "HslCommunication.Profinet.Yokogawa")

[YokogawaLinkHelper 类](../html/6097616d-a102-a94a-1bd2-b7a00819a5fa.htm "YokogawaLinkHelper 类")

[YokogawaLinkHelper 构造函数](../html/25dbdb25-83b0-3dae-677c-adaa17f49cb8.htm "YokogawaLinkHelper 构造函数 ")

[YokogawaLinkHelper 方法](../html/63a4ca2d-bd5d-a701-ae7d-0e3af16bab4a.htm "YokogawaLinkHelper 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| YokogawaLinkHelper 类 |

横河PLC的通信辅助类。

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  HslCommunication.Profinet.YokogawaYokogawaLinkHelper

**命名空间：**
 [HslCommunication.Profinet.Yokogawa](e985e513-76fd-7ba4-62ba-2b2b3d61b929.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class YokogawaLinkHelper
```

```
Public Class YokogawaLinkHelper
```

```
public ref class YokogawaLinkHelper
```

```
type YokogawaLinkHelper =  class end
```

YokogawaLinkHelper 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [YokogawaLinkHelper](25dbdb25-83b0-3dae-677c-adaa17f49cb8.htm) | 初始化 YokogawaLinkHelper 类的一个新实例 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法静态成员 | [GetErrorMsg](11b56442-2f76-97d3-66fb-2fc474155576.htm) | 获取横河PLC的错误的具体描述信息 |
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

[HslCommunication.Profinet.Yokogawa 命名空间](e985e513-76fd-7ba4-62ba-2b2b3d61b929.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## YokogawaLinkHelper 构造函数 

[原文連結](http://api.hslcommunication.cn/html/25dbdb25-83b0-3dae-677c-adaa17f49cb8.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Yokogawa](../html/e985e513-76fd-7ba4-62ba-2b2b3d61b929.htm "HslCommunication.Profinet.Yokogawa")

[YokogawaLinkHelper 类](../html/6097616d-a102-a94a-1bd2-b7a00819a5fa.htm "YokogawaLinkHelper 类")

[YokogawaLinkHelper 构造函数](../html/25dbdb25-83b0-3dae-677c-adaa17f49cb8.htm "YokogawaLinkHelper 构造函数 ")

[YokogawaLinkHelper 方法](../html/63a4ca2d-bd5d-a701-ae7d-0e3af16bab4a.htm "YokogawaLinkHelper 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| YokogawaLinkHelper 构造函数 |

初始化 [YokogawaLinkHelper](6097616d-a102-a94a-1bd2-b7a00819a5fa.htm) 类的一个新实例

**命名空间：**
 [HslCommunication.Profinet.Yokogawa](e985e513-76fd-7ba4-62ba-2b2b3d61b929.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public YokogawaLinkHelper()
```

```
Public Sub New
```

```
public:
YokogawaLinkHelper()
```

```
new : unit -> YokogawaLinkHelper
```

![](../icons/SectionExpanded.png)参见

#### 引用

[YokogawaLinkHelper 类](6097616d-a102-a94a-1bd2-b7a00819a5fa.htm)

[HslCommunication.Profinet.Yokogawa 命名空间](e985e513-76fd-7ba4-62ba-2b2b3d61b929.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## YokogawaLinkHelper 方法

[原文連結](http://api.hslcommunication.cn/html/63a4ca2d-bd5d-a701-ae7d-0e3af16bab4a.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Yokogawa](../html/e985e513-76fd-7ba4-62ba-2b2b3d61b929.htm "HslCommunication.Profinet.Yokogawa")

[YokogawaLinkHelper 类](../html/6097616d-a102-a94a-1bd2-b7a00819a5fa.htm "YokogawaLinkHelper 类")

[YokogawaLinkHelper 方法](../html/63a4ca2d-bd5d-a701-ae7d-0e3af16bab4a.htm "YokogawaLinkHelper 方法")

[GetErrorMsg 方法](../html/11b56442-2f76-97d3-66fb-2fc474155576.htm "GetErrorMsg 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| YokogawaLinkHelper 方法 |

[YokogawaLinkHelper](6097616d-a102-a94a-1bd2-b7a00819a5fa.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法静态成员 | [GetErrorMsg](11b56442-2f76-97d3-66fb-2fc474155576.htm) | 获取横河PLC的错误的具体描述信息 |
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

[YokogawaLinkHelper 类](6097616d-a102-a94a-1bd2-b7a00819a5fa.htm)

[HslCommunication.Profinet.Yokogawa 命名空间](e985e513-76fd-7ba4-62ba-2b2b3d61b929.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GetErrorMsg 方法 

[原文連結](http://api.hslcommunication.cn/html/11b56442-2f76-97d3-66fb-2fc474155576.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Yokogawa](../html/e985e513-76fd-7ba4-62ba-2b2b3d61b929.htm "HslCommunication.Profinet.Yokogawa")

[YokogawaLinkHelper 类](../html/6097616d-a102-a94a-1bd2-b7a00819a5fa.htm "YokogawaLinkHelper 类")

[YokogawaLinkHelper 方法](../html/63a4ca2d-bd5d-a701-ae7d-0e3af16bab4a.htm "YokogawaLinkHelper 方法")

[GetErrorMsg 方法](../html/11b56442-2f76-97d3-66fb-2fc474155576.htm "GetErrorMsg 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| YokogawaLinkHelperGetErrorMsg 方法 |

获取横河PLC的错误的具体描述信息

**命名空间：**
 [HslCommunication.Profinet.Yokogawa](e985e513-76fd-7ba4-62ba-2b2b3d61b929.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static string GetErrorMsg(
	byte code
)
```

```
Public Shared Function GetErrorMsg ( 
	code As Byte
) As String
```

```
public:
static String^ GetErrorMsg(
	unsigned char code
)
```

```
static member GetErrorMsg : 
        code : byte -> string 
```

#### 参数

code
:   类型：SystemByte  
    错误码

#### 返回值

类型：String  
错误的描述信息

![](../icons/SectionExpanded.png)参见

#### 引用

[YokogawaLinkHelper 类](6097616d-a102-a94a-1bd2-b7a00819a5fa.htm)

[HslCommunication.Profinet.Yokogawa 命名空间](e985e513-76fd-7ba4-62ba-2b2b3d61b929.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## YokogawaLinkServer 类

[原文連結](http://api.hslcommunication.cn/html/ab5927f5-b6a1-d76a-7c06-10cbcca7ac37.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Yokogawa](../html/e985e513-76fd-7ba4-62ba-2b2b3d61b929.htm "HslCommunication.Profinet.Yokogawa")

[YokogawaLinkServer 类](../html/ab5927f5-b6a1-d76a-7c06-10cbcca7ac37.htm "YokogawaLinkServer 类")

[YokogawaLinkServer 构造函数](../html/30b1499f-a533-626f-81ae-4a3051d664b2.htm "YokogawaLinkServer 构造函数 ")

[YokogawaLinkServer 属性](../html/c56bc036-ab9b-8156-d3d2-430f68508349.htm "YokogawaLinkServer 属性")

[YokogawaLinkServer 方法](../html/710635c8-b6db-c84f-fafd-1a8fb1d3b10f.htm "YokogawaLinkServer 方法")

[YokogawaLinkServer 事件](../html/459eec4e-5336-2966-0b53-80009ba8ea32.htm "YokogawaLinkServer 事件")

[YokogawaLinkServer 字段](../html/ae5eb6e5-ad9c-dafe-cb5d-1b5e6c7d05bd.htm "YokogawaLinkServer 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| YokogawaLinkServer 类 |

横河PLC的虚拟服务器，支持X,Y,I,E,M,T,C,L继电器类型的数据读写，支持D,B,F,R,V,Z,W,TN,CN寄存器类型的数据读写，可以用来测试横河PLC的二进制通信类型  
Yokogawa PLC's virtual server, supports X, Y, I, E, M, T, C, L relay type data read and write,
supports D, B, F, R, V, Z, W, TN, CN register types The data read and write can be used to test the binary communication type of Yokogawa PLC

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  [HslCommunication.Core.NetBinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)  
    [HslCommunication.Core.DeviceDeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)  
      [HslCommunication.Core.DeviceDeviceServer](bad891c2-0527-acc1-3bcb-0eecb1a393a2.htm)  
        HslCommunication.Profinet.YokogawaYokogawaLinkServer

**命名空间：**
 [HslCommunication.Profinet.Yokogawa](e985e513-76fd-7ba4-62ba-2b2b3d61b929.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class YokogawaLinkServer : DeviceServer
```

```
Public Class YokogawaLinkServer
	Inherits DeviceServer
```

```
public ref class YokogawaLinkServer : public DeviceServer
```

```
type YokogawaLinkServer =  
    class
        inherit DeviceServer
    end
```

YokogawaLinkServer 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [YokogawaLinkServer](30b1499f-a533-626f-81ae-4a3051d664b2.htm) | 实例化一个横河PLC的服务器，支持X,Y,I,E,M,T,C,L继电器类型的数据读写，支持D,B,F,R,V,Z,W,TN,CN寄存器类型的数据读写  Instantiate a Yokogawa PLC server, support X, Y, I, E, M, T, C, L relay type data read and write, support D, B, F, R, V, Z, W, TN, CN Register type data reading and writing |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [ActiveTimeSpan](1e32f805-ec96-2746-dc55-44e253ee8514.htm) | 获取或设置两次数据交互时的最小时间间隔，默认为24小时。如果超过该设定的时间不进行数据交互，服务器就会强制断开当前的连接操作。  Get or set the minimum time interval between two data interactions, the default is 24 hours. If the data exchange is not performed for more than the set time, the server will forcibly disconnect the current connection operation. (继承自 [DeviceServer](bad891c2-0527-acc1-3bcb-0eecb1a393a2.htm)。) |
| 公共属性 | [BothModeUdpPort](a2b83e02-9594-1480-2698-f5f0d23e53ca.htm) | 当服务器同时启动TCP及UDP服务的时候，获取当前的UDP服务的端口号  When the server starts TCP and UDP services at the same time, it obtains the port number of the current UDP service (继承自 [DeviceServer](bad891c2-0527-acc1-3bcb-0eecb1a393a2.htm)。) |
| 公共属性代码示例 | [ByteTransform](f35118c9-3691-f752-ae6f-2c5549714faa.htm) | 当前的数据变换机制，当你需要从字节数据转换类型数据的时候需要。  The current data transformation mechanism is required when you need to convert type data from byte data. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共属性 | [CommunicationPipe](ad8e903b-d289-21b0-077e-15ede8005221.htm) | 获取或设置当前的管道信息，管道类型为[CommunicationPipe](ad8e903b-d289-21b0-077e-15ede8005221.htm)的继承类，内置了[PipeTcpNet](26b94bd5-2695-36e7-6378-c3d41846be8f.htm)管道，[PipeUdpNet](5a0b810e-187e-826b-0a13-34fd5576aa0e.htm)管道，[PipeSerialPort](9d924607-3dfc-3916-b532-c5dec48fe347.htm)管道等  Get or set the current pipeline information, the pipeline type is [CommunicationPipe](ad8e903b-d289-21b0-077e-15ede8005221.htm) inheritance class, [PipeTcpNet](26b94bd5-2695-36e7-6378-c3d41846be8f.htm) pipeline, [PipeUdpNet](5a0b810e-187e-826b-0a13-34fd5576aa0e.htm) pipeline, [PipeSerialPort](9d924607-3dfc-3916-b532-c5dec48fe347.htm) pipeline, etc (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [ConnectionId](ce1efb4b-fbc9-4683-163e-0a8e71e99d6b.htm) | 当前连接的唯一ID号，默认为长度20的guid码加随机数组成，方便列表管理，也可以自己指定  The unique ID number of the current connection. The default is a 20-digit guid code plus a random number. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [EnableIPv6](12e23164-745f-c405-453f-26b5177ac5d4.htm) | 获取或设置服务器是否支持IPv6的地址协议信息  Get or set whether the server supports IPv6 address protocol information (继承自 [DeviceServer](bad891c2-0527-acc1-3bcb-0eecb1a393a2.htm)。) |
| 公共属性 | [EnableWrite](e0eeb367-12e3-cf66-19ee-e5938dffbc3e.htm) | 获取或设置当前的服务器是否允许远程客户端进行写入数据操作，默认为True  Gets or sets whether the current server allows remote clients to write data, the default is True (继承自 [DeviceServer](bad891c2-0527-acc1-3bcb-0eecb1a393a2.htm)。) |
| 公共属性 | [ForceSerialReceiveOnce](54cca9f0-f9cd-4743-6bfe-ffc98f5b3163.htm) | 获取或设置当前的服务器接收串口数据时候，是否强制只接收一次数据，默认为false，适合点对点通信，如果你总线形式的连接，则需要设置 True  Get or set whether to force the data to be received only once when the current server receives serial port data. The default value is false, which is suitable for point-to-point communication. If you have a bus connection, you need to set True (继承自 [DeviceServer](bad891c2-0527-acc1-3bcb-0eecb1a393a2.htm)。) |
| 公共属性 | [IsStarted](fbe63af3-720e-11d0-e808-51ff135541d3.htm) | 服务器引擎是否启动  Whether the server engine is started (继承自 [DeviceServer](bad891c2-0527-acc1-3bcb-0eecb1a393a2.htm)。) |
| 公共属性 | [LocalAddress](2d19d9e8-e035-28b9-8aba-3649ef66cc9a.htm) | 获取或设置服务器绑定的本地IP地址，默认为空，使用本地所有可用的ip地址  Obtain or set the local IP address bound to the server, which is empty by default, and use all available IP addresses locally (继承自 [DeviceServer](bad891c2-0527-acc1-3bcb-0eecb1a393a2.htm)。) |
| 公共属性代码示例 | [LogNet](8c56f246-88ac-afb8-8c06-972fafe4f78a.htm) | 组件的日志工具，支持日志记录，只要实例化后，当前网络的基本信息，就以[DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm)等级进行输出  The component's logging tool supports logging. As long as the instantiation of the basic network information, the output will be output at [DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm) (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [OnlineCount](0c146ed0-b958-3374-5952-8069aaee70d2.htm) | 获取在线的客户端的数量  Get the number of clients online (继承自 [DeviceServer](bad891c2-0527-acc1-3bcb-0eecb1a393a2.htm)。) |
| 公共属性 | [Port](858c6445-e24b-21e4-7c2f-3dff7ccb0cb2.htm) | 获取或设置服务器的端口号，如果是设置，需要在服务器启动前设置完成，才能生效。  Gets or sets the port number of the server. If it is set, it needs to be set before the server starts to take effect. (继承自 [DeviceServer](bad891c2-0527-acc1-3bcb-0eecb1a393a2.htm)。) |
| 公共属性代码示例 | [ReceiveTimeOut](250560a7-332b-ce1b-5eb0-e3abe4ee6012.htm) | 获取或设置接收服务器反馈的时间，如果为负数，则不接收反馈   Gets or sets the time to receive server feedback, and if it is a negative number, does not receive feedback (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [SendBeforeHex](72aa1de4-116f-abfe-4adb-814a0be0e560.htm) | 获取或设置在发送通信报文前追加发送的字节信息，HEX格式，通常用于lora组网时，需要携带 00 00 00 02 四个字节的站地址功能。  Obtain or set the byte information sent before sending communication packets, HEX format, usually used for LORA networking, you need to carry 00 00 00 02 four-byte station address function. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [SerialReceiveAtleastTime](75be7ae9-07a9-46a6-0d31-a978821490c8.htm) | 获取或设置串口模式下，接收一条数据最短的时间要求，当设备发送的数据非常慢的时候，或是分割发送数据的时候，就需要将本值设置的大一点，默认为20ms  Get or set the shortest time required to receive a piece of data in serial port mode. When the data sent by the device is very slow, or when the data is divided and sent, you need to set this value to a larger value, the default is 20ms (继承自 [DeviceServer](bad891c2-0527-acc1-3bcb-0eecb1a393a2.htm)。) |
| 公共属性 | [ServerMode](2a437b8e-94ea-1b29-44b9-a1d22fa6aa03.htm) | 当前服务器的模式，0：TCP服务器，1：UDP服务器，2：TCP及UDP服务器  Gets whether the current server is a TCP server or a UDP server (继承自 [DeviceServer](bad891c2-0527-acc1-3bcb-0eecb1a393a2.htm)。) |
| 公共属性 | [SleepTime](86a5fbb9-3456-6beb-7113-5bbad5eabbc5.htm) | 获取或设置在正式接收对方返回数据前的时候，需要休息的时间，当设置为0的时候，不需要休息。  Get or set the time required to rest before officially receiving the data from the other party. When it is set to 0, no rest is required. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [SocketKeepAliveTime](6ba850ca-d152-1045-fbd9-39bc60410ab6.htm) | 获取或设置客户端的Socket的心跳时间信息，这个是Socket底层自动实现的心跳包，不基于协议层实现。默认小于0，不开启心跳检测，如果需要开启，设置 60\_000 比较合适，单位毫秒  Get or set the heartbeat time information of the Socket of the client. This is the heartbeat packet automatically implemented by the bottom layer of the Socket, not based on the protocol layer. The default value is less than 0, and heartbeat detection is not enabled. If you need to enable it, it is more appropriate to set 60\_000, in milliseconds. (继承自 [DeviceServer](bad891c2-0527-acc1-3bcb-0eecb1a393a2.htm)。) |
| 受保护的属性 | [WordLength](f1a505cc-e85c-6c8b-b189-58a1ab5753f9.htm) | 一个字单位的数据表示的地址长度，西门子为2，三菱，欧姆龙，modbusTcp就为1，AB PLC无效  The address length represented by one word of data, Siemens is 2, Mitsubishi, Omron, modbusTcp is 1, AB PLC is invalid (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 受保护的方法 | [CheckSerialReceiveDataComplete](f55d7a1e-3d6b-6a34-b0c3-f9e615218479.htm) | 检查串口接收的数据是否完成的方法，如果接收完成，则返回True (继承自 [DeviceServer](bad891c2-0527-acc1-3bcb-0eecb1a393a2.htm)。) |
| 受保护的方法 | [CloseAction](37b80be6-d1d8-dd5c-0ab7-c0f191302c20.htm) | 服务器关闭的时候需要做的事情  Things to do when the server is down (继承自 [DeviceServer](bad891c2-0527-acc1-3bcb-0eecb1a393a2.htm)。) |
| 公共方法 | [CloseSerialSlave](634feaea-15dd-5051-f41d-9b82de27af7b.htm) | 关闭提供从机服务的串口对象  Close the serial port object that provides slave services (继承自 [DeviceServer](bad891c2-0527-acc1-3bcb-0eecb1a393a2.htm)。) |
| 受保护的方法 | [DecideWhetherQAMessage](acc40cda-e310-8537-6685-59b67eeb16ff.htm) | 决定当前的消息是否是用于问答机制返回的消息，默认直接返回 true, 实际的情况需要根据协议进行重写方法  To determine whether the current message is the message returned by the question answering mechanism, the default is true. In actual cases, the rewriting method needs to be performed according to the protocol (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [Dispose](b9dda6bf-d342-254e-689a-93fa0bb265d1.htm) | (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 受保护的方法 | [Dispose(Boolean)](09988646-d3b6-cd16-4fcc-d93ebd932096.htm) | 释放被 YokogawaLinkServer 使用的非托管资源，并且是否托管资源（可选） (重写 [DeviceServerDispose(Boolean)](219f18c3-64d3-a044-b798-de6e2bf866ef.htm).) |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | [ExtraAfterReadFromCoreServer](06a49a15-a3ce-5e39-28c7-2b78851eca2e.htm) | 和服务器交互完成的时候调用的方法，可以根据读写结果进行一些额外的操作，具体的操作需要根据实际的需求来重写实现  The method called when the interaction with the server is completed can perform some additional operations based on the read and write results. The specific operations need to be rewritten according to actual needs. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法代码示例 | [ExtraOnDisconnect](de9dbba2-3b2c-1c5e-53c8-4e6f8a44e2c4.htm) | 根据实际的协议选择是否重写本方法，有些协议在断开连接之前，需要发送一些报文来关闭当前的网络通道  Select whether to rewrite this method according to the actual protocol. Some protocols need to send some packets to close the current network channel before disconnecting. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法代码示例 | [ExtraOnDisconnectAsync](210f593e-16d8-77c4-1178-58fbcfe9b27b.htm) | 根据实际的协议选择是否重写本方法，有些协议在断开连接之前，需要发送一些报文来关闭当前的网络通道  Select whether to rewrite this method according to the actual protocol. Some protocols need to send some packets to close the current network channel before disconnecting. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | [GetCommunicationServer](160e612f-d0c3-0bc4-7a2c-7b60d66835cd.htm) | 获取当前的核心服务器信息 (继承自 [DeviceServer](bad891c2-0527-acc1-3bcb-0eecb1a393a2.htm)。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 受保护的方法 | [GetLogTextFromBinary](67f3b2f7-957c-fa0b-8320-6799237157db.htm) | 获取当前的报文进行日志记录的时候，是否使用二进制的格式记录，默认返回 [LogMsgFormatBinary](2ed2ff87-caec-0e44-f08b-e2b015225a58.htm)，重写可以根据session对象分别返回不同记录模式  Whether to log the current packet in binary format, the default return is [LogMsgFormatBinary](2ed2ff87-caec-0e44-f08b-e2b015225a58.htm). If you want to override it, different recording modes can be returned according to session (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [GetNewNetMessage](c8e4092d-6573-61f0-1035-4a9a062719c9.htm) | 获取一个新的消息对象的方法，需要在继承类里面进行重写  The method to get a new message object needs to be overridden in the inheritance class (重写 [BinaryCommunicationGetNewNetMessage](16636a1b-a6f0-9c0f-9d9c-c94e26c67459.htm).) |
| 公共方法 | [GetTrustedClients](0531cf85-32d8-c438-ccce-b544fc954393.htm) | 获取受信任的客户端列表  Get a list of trusted clients (继承自 [DeviceServer](bad891c2-0527-acc1-3bcb-0eecb1a393a2.htm)。) |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法 | [GetWordLength](c27821ea-1739-8eb9-a848-23a00828fec8.htm) | 一个字单位的数据表示的地址长度，西门子为2，三菱，欧姆龙，modbusTcp就为1，AB PLC无效  The address length represented by one word of data, Siemens is 2, Mitsubishi, Omron, modbusTcp is 1, AB PLC is invalid (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 受保护的方法代码示例 | [InitializationOnConnect](4ab6989f-a5b5-05fe-8b44-bfa45824a2a7.htm) | 根据实际的协议选择是否重写本方法，有些协议在创建连接之后，需要进行一些初始化的信号握手，才能最终建立网络通道。  Whether to rewrite this method is based on the actual protocol. Some protocols require some initial signal handshake to establish a network channel after the connection is created. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法代码示例 | [InitializationOnConnectAsync](7da2af04-d4fc-5045-19a4-4b0d1722eef3.htm) | 根据实际的协议选择是否重写本方法，有些协议在创建连接之后，需要进行一些初始化的信号握手，才能最终建立网络通道。  Whether to rewrite this method is based on the actual protocol. Some protocols require some initial signal handshake to establish a network channel after the connection is created. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [LoadDataPool](346b1333-ac23-cdfc-8edd-3a5ed66ab4bc.htm) | 从文件加载数据池信息  Load datapool information from a file (继承自 [DeviceServer](bad891c2-0527-acc1-3bcb-0eecb1a393a2.htm)。) |
| 受保护的方法 | [LoadFromBytes](b6683fe3-9938-32fe-5511-9dc08ef6f0eb.htm) | 从字节数据加载数据信息，需要进行重写方法  Loading data information from byte data requires rewriting method (重写 [DeviceServerLoadFromBytes(Byte)](4a16c815-c1b7-9620-c705-21db6518272b.htm).) |
| 受保护的方法 | [LogRevcMessage(Byte)](dab1657f-1e75-7fe1-7003-ccedfc69855d.htm) | 使用日志记录一个接收的报文信息  Logs are used to record information about a received packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [LogRevcMessage(Byte, PipeSession)](0f51ae45-575e-4228-f978-01b37d6d5df1.htm) | 使用日志记录一个接收的报文信息  Logs are used to record information about a received packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [LogSendMessage(Byte)](4f90eecc-68a2-1ef2-44ef-369f19bed3ec.htm) | 使用日志记录一个发送的报文信息  Logs are used to record information about a send packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [LogSendMessage(Byte, PipeSession)](209caadf-4f0f-40aa-04c2-ebc61f8dab13.htm) | 使用日志记录一个发送的报文信息  Logs are used to record information about a send packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法 | [PackCommandWithHeader](d0b2cd76-c306-11fc-38c7-af8ad306cce7.htm) | 对当前的命令进行打包处理，通常是携带命令头内容，标记当前的命令的长度信息，需要进行重写，否则默认不打包  The current command is packaged, usually carrying the content of the command header, marking the length of the current command, and it needs to be rewritten, otherwise it is not packaged by default (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [RaiseDataReceived](18d5bad7-d0b9-b744-a88c-ec268dd951a0.htm) | 触发一个数据接收的事件信息  Event information that triggers a data reception (继承自 [DeviceServer](bad891c2-0527-acc1-3bcb-0eecb1a393a2.htm)。) |
| 受保护的方法 | [RaiseDataSend](f3aafc2a-1633-37d7-6931-80f791b8e82e.htm) | 触发一个数据发送的事件信息  Event information that triggers a data transmission (继承自 [DeviceServer](bad891c2-0527-acc1-3bcb-0eecb1a393a2.htm)。) |
| 公共方法 | [Read(String, UInt16)](89ea561c-d11c-40f1-0ab3-0654d01d6852.htm) | 批量读取字节数组信息，需要指定地址和长度，返回原始的字节数组  Batch read byte array information, need to specify the address and length, return the original byte array (重写 [DeviceCommunicationRead(String, UInt16)](9dd09452-c200-bc65-75ff-69aa34b366d4.htm).) |
| 公共方法代码示例 | [ReadT](95abe064-5594-5ae7-a849-98c53f88fee0.htm) | 读取支持Hsl特性的数据内容，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明，或是官网的详细文档  Read the data content of the Hsl attribute. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadAsync(String, UInt16)](d8c14d37-56dc-2a31-1a19-3e2b5b65197d.htm) | 异步批量读取字节数组信息，需要指定地址和长度，返回原始的字节数组  Asynchronous batch read byte array information, need to specify the address and length, return the original byte array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadAsyncT](0522e740-4bf0-2e31-c4c5-e41384321b37.htm) | 异步读取支持Hsl特性的数据内容，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Asynchronously read the data content of the Hsl attribute. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadBool(String)](9b2faa2e-eaa9-7ace-c49f-7a0318f1df1a.htm) | 读取单个的Boolean数据信息  Read a single Boolean data message (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadBool(String, UInt16)](3ea0f3cf-91d1-b748-194e-fcddc80b2430.htm) | 批量读取Boolean数组信息，需要指定地址和长度，返回Boolean 数组  Batch read Boolean array information, need to specify the address and length, return Boolean array (重写 [DeviceCommunicationReadBool(String, UInt16)](69d0d24f-34a7-5b41-5c1b-619edc7e6b54.htm).) |
| 公共方法 | [ReadBoolAsync(String)](cc37da1f-d673-fe67-78ee-4e465d29a2bd.htm) | 异步读取单个的Boolean数据信息  Asynchronously read a single Boolean data message (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadBoolAsync(String, UInt16)](a45d383d-0095-3ea7-32f8-e75af88d60f8.htm) | 异步批量读取Boolean数组信息，需要指定地址和长度，返回Boolean 数组  Asynchronously batch read Boolean array information, need to specify the address and length, return Boolean array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
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
| 受保护的方法代码示例 | [ReadFromCoreServer(PipeSession, Byte)](3950b0e0-214c-75f3-71a8-87421282ea1f.htm) | 将数据发送到当前的网络通道中，并从网络通道中接收一个[INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm)指定的完整的报文，网络通道将根据[GetAvailableSocket](e86a6fc1-eb91-8de9-5ff8-799b0ea3786b.htm)方法自动获取，本方法是线程安全的。  Send data to the current network channel and receive a complete message specified by [INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm) from the network channel. The network channel will be automatically obtained according to the [GetAvailableSocket](e86a6fc1-eb91-8de9-5ff8-799b0ea3786b.htm) method This method is thread-safe. (重写 [DeviceServerReadFromCoreServer(PipeSession, Byte)](bad4d1eb-5428-086b-8767-3440239b6538.htm).) |
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
| 公共方法代码示例 | [ReadStructT](49853082-fbad-52e5-d756-0615cedb4b83.htm) | 读取结构体类型的数据，根据结构体自身的定义，读取原始字节数组，然后解析出实际的结构体数据，结构体需要实现[HslStructAttribute](87e427a8-0c7b-43a8-c421-c27264eda90f.htm)特性 (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
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
| 公共方法 | [SaveDataPool](f3d0c3b1-3383-d593-4b20-0bfe8b1f56bf.htm) | 将本系统的数据池数据存储到指定的文件  Store the data pool data of this system to the specified file (继承自 [DeviceServer](bad891c2-0527-acc1-3bcb-0eecb1a393a2.htm)。) |
| 受保护的方法 | [SaveToBytes](bbb5fdfe-745d-1fc9-62bb-a401c143934f.htm) | 将数据信息存储到字节数组去，需要进行重写方法  To store data information into a byte array, a rewrite method is required (重写 [DeviceServerSaveToBytes](871f8a64-6b6d-53ea-143c-479af077c5a3.htm).) |
| 公共方法 | [ServerClose](b919d91b-5ed0-73c5-dff4-1a1bd2f03bae.htm) | 关闭服务器的引擎  Shut down the server's engine (继承自 [DeviceServer](bad891c2-0527-acc1-3bcb-0eecb1a393a2.htm)。) |
| 公共方法 | [ServerStart](9e137e74-16b0-fcf5-97f3-b1369f55edf0.htm) | 使用已经配置好的端口启动服务器的引擎，并且使用TCP模式  Use the configured port to start the server's engine (继承自 [DeviceServer](bad891c2-0527-acc1-3bcb-0eecb1a393a2.htm)。) |
| 公共方法 | [ServerStart(Int32, Boolean)](bcda1e54-84c1-903a-a284-c2a08870e75e.htm) | 指定端口号来启动服务器的引擎  Specify the port number to start the server's engine (继承自 [DeviceServer](bad891c2-0527-acc1-3bcb-0eecb1a393a2.htm)。) |
| 公共方法 | [ServerStart(Int32, Int32)](f8a8fc7f-c557-9fae-f763-7d8e4b6df3ed.htm) | 指定一个TCP端口及UDP端口，同时启动两种模式的服务器  Specify a TCP port and a UDP port to start the server in both modes at the same time (继承自 [DeviceServer](bad891c2-0527-acc1-3bcb-0eecb1a393a2.htm)。) |
| 公共方法 | [SetDtuPipe](f765d1d7-234a-b229-9457-97e25e024dc9.htm) | 将当前的通信对象设置DTU模式，允许传入现成的管道，并返回初始化结果，如果该设备重写了握手报文，就是返回握手结果  Set the current communication object to DTU mode, allow the existing pipe to be passed in, and return the initialization result, if the device rewrites the handshake packet, the handshake result is returned (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [SetDtuPipeAsync](9d1fd12d-b804-3828-a627-a2c0eef1d23f.htm) | 将当前的通信对象设置DTU模式，允许传入现成的管道，并返回初始化结果，如果该设备重写了握手报文，就是返回握手结果  Set the current communication object to DTU mode, allow the existing pipe to be passed in, and return the initialization result, if the device rewrites the handshake packet, the handshake result is returned (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [SetTrustedIpAddress](7c62dce6-ef4b-bc97-11c9-2ec592d36915.htm) | 设置并启动受信任的客户端登录并读写，如果为null，将关闭对客户端的ip验证  Set and start the trusted client login and read and write, if it is null, the client's IP verification will be turned off (继承自 [DeviceServer](bad891c2-0527-acc1-3bcb-0eecb1a393a2.htm)。) |
| 受保护的方法 | [SocketAcceptExtraCheck](3e3d1f6e-95f9-2f55-ca94-f1802aac9fc2.htm) | 当客户端的socket登录的时候额外检查的操作，并返回操作的结果信息。  The operation is additionally checked when the client's socket logs in, and the result information of the operation is returned. (继承自 [DeviceServer](bad891c2-0527-acc1-3bcb-0eecb1a393a2.htm)。) |
| 受保护的方法 | [StartInitialization](5c580365-04f4-736f-0acc-3146cfb94e10.htm) | 服务器启动时额外的初始化信息，可以用于启动一些额外的服务的操作。  The extra initialization information when the server starts can be used to start some additional service operations. (继承自 [DeviceServer](bad891c2-0527-acc1-3bcb-0eecb1a393a2.htm)。) |
| 公共方法 | [StartProgram](3dc02561-e8a1-57cf-6ef6-bb613eb4e6a0.htm) | 如果未执行程序，则开始执行程序  Starts executing a program if it is not being executed |
| 公共方法 | [StartSerialSlave(ActionSerialPort)](311866fc-5e1c-201d-a1e8-e1a5fafe47f7.htm) | 启动串口的从机服务，使用自定义的初始化方法初始化串口的参数  Start the slave service of serial and initialize the parameters of the serial port using a custom initialization method (继承自 [DeviceServer](bad891c2-0527-acc1-3bcb-0eecb1a393a2.htm)。) |
| 公共方法 | [StartSerialSlave(String)](c719d5ed-8fef-bf84-adc3-0e6af138235d.htm) | 启动串口的从机服务，使用默认的参数进行初始化串口，9600波特率，8位数据位，无奇偶校验，1位停止位  Start the slave service of serial, initialize the serial port with default parameters, 9600 baud rate, 8 data bits, no parity, 1 stop bit (继承自 [DeviceServer](bad891c2-0527-acc1-3bcb-0eecb1a393a2.htm)。) |
| 公共方法 | [StartSerialSlave(String, Int32)](1b6bc90f-70f4-4e2d-01bf-9cd3d6cde88e.htm) | 启动串口的从机服务，使用默认的参数进行初始化串口，8位数据位，无奇偶校验，1位停止位  Start the slave service of serial, initialize the serial port with default parameters, 8 data bits, no parity, 1 stop bit (继承自 [DeviceServer](bad891c2-0527-acc1-3bcb-0eecb1a393a2.htm)。) |
| 公共方法 | [StartSerialSlave(String, Int32, Int32, Parity, StopBits)](353ae919-c1df-6dd7-7b08-1a9011fbd40e.htm) | 启动串口的从机服务，使用指定的参数进行初始化串口，指定数据位，指定奇偶校验，指定停止位 (继承自 [DeviceServer](bad891c2-0527-acc1-3bcb-0eecb1a393a2.htm)。) |
| 公共方法 | [StopProgram](0a212b25-16a8-da15-a0c9-6c8130126998.htm) | 停止当前正在执行程序  Stops the executing program. |
| 受保护的方法 | [ThreadPoolLoginAfterClientCheck](097524c4-f7f3-4529-af08-3773195a5794.htm) | 当客户端登录后，在Ip信息的过滤后，然后触发本方法，进行后续的数据接收，处理，并返回相关的数据信息  When the client logs in, after filtering the IP information, this method is then triggered to perform subsequent data reception, processing, and return related data information (继承自 [DeviceServer](bad891c2-0527-acc1-3bcb-0eecb1a393a2.htm)。) |
| 公共方法 | [ToString](572f031c-6c28-1c80-4cd0-e0b99d3144de.htm) | (重写 [DeviceCommunicationToString](5667665a-da16-d227-f350-5b3f0ab9ab04.htm).) |
| 公共方法 | [UnpackResponseContent](208017a7-8da2-33f8-3273-d7e7e850100f.htm) | 根据对方返回的报文命令，对命令进行基本的拆包，例如各种Modbus协议拆包为统一的核心报文，还支持对报文的验证  According to the message command returned by the other party, the command is basically unpacked, for example, various Modbus protocols are unpacked into a unified core message, and the verification of the message is also supported (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [UseSSL(X509Certificate)](3dd75f50-9dd7-58e2-28c3-d8870a6b2b1b.htm) | 使用SSL通信，传递一个证书的对象 (继承自 [DeviceServer](bad891c2-0527-acc1-3bcb-0eecb1a393a2.htm)。) |
| 公共方法 | [UseSSL(String, String)](08755f4f-11f2-e957-58f2-54abf3706213.htm) | 使用SSL通信，传递一个证书的路径，以及证书的密码 (继承自 [DeviceServer](bad891c2-0527-acc1-3bcb-0eecb1a393a2.htm)。) |
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
| 公共方法 | [Write(String, Boolean)](f1d3e757-a2d4-320d-11cc-9853ffc9b2ab.htm) | 批量写入Boolean数组数据，返回是否成功  Batch write Boolean array data, return whether the write was successful (重写 [DeviceCommunicationWrite(String, Boolean)](5b046c82-3d8e-550b-fdf3-41ca5558169b.htm).) |
| 公共方法 | [Write(String, Byte)](6e9f0b54-4ae7-f483-a1e8-70952e632b24.htm) | 写入原始的byte数组数据到指定的地址，返回是否写入成功  Write the original byte array data to the specified address, and return whether the write was successful (重写 [DeviceCommunicationWrite(String, Byte)](b5536dd6-ccde-6883-ac27-bff9e4b37806.htm).) |
| 公共方法代码示例 | [Write(String, String, Int32)](02fe2ca9-49a3-ede4-d1af-8c6da1f644da.htm) | 写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为ASCII  Write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, it is filled with 0 and the encoding is ASCII. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, String, Encoding)](7b602661-902d-2b70-1fdb-ba87ecd307bf.htm) | 写入字符串信息，需要指定的编码信息  Write string information, need to specify the encoding information (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, String, Int32, Encoding)](e8cb30d9-4c5e-6c68-b9af-1b1b6c2dc714.htm) | 写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为指定的编码信息  Write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, then the operation is complemented with 0 , you should specified the encoding information (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteT(T)](9f48853f-fbeb-c37d-25a2-1c76d1dee407.htm) | 写入支持Hsl特性的数据，返回是否写入成功，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Write data that supports the Hsl attribute, and return whether the write was successful. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WriteAsync(String, Boolean)](10fe669c-805e-1e5a-f016-f319ad363e90.htm) | 异步批量写入Boolean数组数据，返回是否成功  Asynchronously batch write Boolean array data, return success (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WriteAsync(String, Boolean)](fe872dae-7ef9-12c1-d883-e92ba8676416.htm) | 异步批量写入Boolean数组数据，返回是否成功  Asynchronously batch write Boolean array data, return success (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Byte)](8f1cec30-8b86-6749-c8f3-7d39168d59df.htm) | 异步写入原始的byte数组数据到指定的地址，返回是否写入成功  Asynchronously writes the original byte array data to the specified address, and returns whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
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
| 公共方法 | [WriteAsync(String, String, Int32)](b739b983-b03c-8ecf-2713-75f19233ae2f.htm) | 异步写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为ASCII  Asynchronously write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, it is filled with 0 and the encoding is ASCII. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, String, Encoding)](876e9800-0406-6407-bdf8-8ab1cf28df38.htm) | 异步写入字符串信息，需要指定的编码信息  Asynchronously write string information, need to specify the encoding information (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WriteAsync(String, String, Int32, Encoding)](ae9fb2a3-fddb-48eb-e9ad-6cc0ebbac1f6.htm) | 异步写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为指定的编码信息  Asynchronously write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, then the operation is complemented with 0 , you should specified the encoding information (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsyncT(T)](25d96e77-d1b2-918c-33f9-e5e4cc57c9cd.htm) | 异步写入支持Hsl特性的数据，返回是否写入成功，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Asynchronously write data that supports the Hsl attribute, and return whether the write was successful. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteCustomerT](84a37f82-3f67-a0f8-fd3f-0630772565e4.htm) | 写入自定义类型的数据，该类型必须继承自IDataTransfer接口  Write data of a custom type, which must inherit from the IDataTransfer interface (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteCustomerAsyncT](04628397-cfa1-d785-7790-4c02cb377ab6.htm) | 写入自定义类型的数据，该类型必须继承自IDataTransfer接口  Write data of a custom type, which must inherit from the IDataTransfer interface (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)事件

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共事件代码示例 | [OnDataReceived](821fdec2-5ea2-1d1c-7f5c-08f594d7d1ba.htm) | 接收到数据的时候就触发的事件，示例详细参考API文档信息  An event that is triggered when data is received (继承自 [DeviceServer](bad891c2-0527-acc1-3bcb-0eecb1a393a2.htm)。) |
| 公共事件 | [OnDataSend](deabf344-282f-a8c2-8d71-8f1512a66e38.htm) | 数据发送的时候就触发的事件  Events that are triggered when data is sent (继承自 [DeviceServer](bad891c2-0527-acc1-3bcb-0eecb1a393a2.htm)。) |

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

其中的X继电器可以在服务器进行读写操作，但是远程的PLC只能进行读取，所有的数据读写的最大的范围按照协议进行了限制。

![](../icons/SectionExpanded.png)示例

你可以很快速并且简单的创建一个虚拟的横河服务器

简单的创建服务器

[复制](# "复制")

```
private YokogawaLinkServer yokogawaLinkServer;
public void S7Start( )
{
    try
    {
        yokogawaLinkServer = new YokogawaLinkServer( );
        yokogawaLinkServer.ServerStart( 12289 );
    }
    catch (Exception ex)
    {
        Console.Write( "Failed:" + HslCommunication.BasicFramework.SoftBasic.GetExceptionMessage( ex ) );
    }
}
```

当然如果需要高级的服务器，指定日志，限制客户端的IP地址，获取客户端发送的信息，在服务器初始化的时候就要参照下面的代码：

定制服务器

[复制](# "复制")

```
public void S7Start4( )
{
    try
    {
        yokogawaLinkServer = new YokogawaLinkServer( );
        yokogawaLinkServer.LogNet = new HslCommunication.LogNet.LogNetSingle( "logs.txt" );        // 配置日志信息
        yokogawaLinkServer.SetTrustedIpAddress( new List<string>( ) { "127.0.0.1" } );             // 仅仅限制本机客户端读写
        yokogawaLinkServer.OnDataReceived += S7NetServer_OnDataReceived;
        yokogawaLinkServer.ServerStart( 12289 );
    }
    catch (Exception ex)
    {
        Console.Write( "Failed:" + HslCommunication.BasicFramework.SoftBasic.GetExceptionMessage( ex ) );
    }
}

private void S7NetServer_OnDataReceived( object sender, object source, byte[] data )
{
    Console.WriteLine( HslCommunication.BasicFramework.SoftBasic.ByteToHexString( data, ' ' ) ); // 打印客户端发送的数据
}
```

服务器创建好之后，我们就可以对服务器进行一些读写的操作了，下面的代码是基础的BCL类型的读写操作。

基础的读写示例

[复制](# "复制")

```
private void ReadExample( )
{
    // 此处以D100寄存器作为示例
    bool bool_M100_0 = yokogawaLinkServer.ReadBool( "M100" ).Content;
    short short_D100 = yokogawaLinkServer.ReadInt16( "D100" ).Content;        // 读取D100-M101组成的字
    ushort ushort_D100 = yokogawaLinkServer.ReadUInt16( "D100" ).Content;       // 读取D100-M101组成的无符号的值
    int int_D100 = yokogawaLinkServer.ReadInt32( "D100" ).Content;        // 读取D100-M103组成的有符号的数据
    uint uint_D100 = yokogawaLinkServer.ReadUInt32( "D100" ).Content;       // 读取D100-M103组成的无符号的值
    float float_D100 = yokogawaLinkServer.ReadFloat( "D100" ).Content;        // 读取D100-M103组成的单精度值
    long long_D100 = yokogawaLinkServer.ReadInt64( "D100" ).Content;        // 读取D100-M107组成的大数据值
    ulong ulong_D100 = yokogawaLinkServer.ReadUInt64( "D100" ).Content;       // 读取D100-M107组成的无符号大数据
    double double_D100 = yokogawaLinkServer.ReadDouble( "D100" ).Content;       // 读取D100-M107组成的双精度值
    string string_D100 = yokogawaLinkServer.ReadString( "D100", 10 ).Content;   // 读取D100-M109组成的ASCII字符串数据

    // 读取数组
    short[] short_D100_array = yokogawaLinkServer.ReadInt16( "D100", 10 ).Content;      // 读取D100-M101组成的字
    ushort[] ushort_D100_array = yokogawaLinkServer.ReadUInt16( "D100", 10 ).Content;     // 读取D100-M101组成的无符号的值
    int[] int_D100_array = yokogawaLinkServer.ReadInt32( "D100", 10 ).Content;      // 读取D100-M103组成的有符号的数据
    uint[] uint_D100_array = yokogawaLinkServer.ReadUInt32( "D100", 10 ).Content;     // 读取D100-M103组成的无符号的值
    float[] float_D100_array = yokogawaLinkServer.ReadFloat( "D100", 10 ).Content;      // 读取D100-M103组成的单精度值
    long[] long_D100_array = yokogawaLinkServer.ReadInt64( "D100", 10 ).Content;      // 读取D100-M107组成的大数据值
    ulong[] ulong_D100_array = yokogawaLinkServer.ReadUInt64( "D100", 10 ).Content;     // 读取D100-M107组成的无符号大数据
    double[] double_D100_array = yokogawaLinkServer.ReadDouble( "D100", 10 ).Content;     // 读取D100-M107组成的双精度值
}

private void WriteExample( )
{

    // 此处以D100寄存器作为示例
    yokogawaLinkServer.Write( "M100", true );                       // 写入D100  bool值
    yokogawaLinkServer.Write( "D100", (byte)123 );                  // 写入D100  byte值
    yokogawaLinkServer.Write( "D100", (short)1234 );                // 写入D100  short值
    yokogawaLinkServer.Write( "D100", (ushort)45678 );              // 写入D100  ushort值
    yokogawaLinkServer.Write( "D100", 1234566 );                    // 写入D100  int值
    yokogawaLinkServer.Write( "D100", (uint)1234566 );              // 写入D100  uint值
    yokogawaLinkServer.Write( "D100", 123.456f );                   // 写入D100  float值
    yokogawaLinkServer.Write( "D100", 123.456d );                   // 写入D100  double值
    yokogawaLinkServer.Write( "D100", 123456661235123534L );        // 写入D100  long值
    yokogawaLinkServer.Write( "D100", 523456661235123534UL );       // 写入D100  ulong值
    yokogawaLinkServer.Write( "D100", "K123456789" );               // 写入D100  string值

    // 读取数组
    yokogawaLinkServer.Write( "D100", new short[] { 123, 3566, -123 } );                                      // 写入D100  short值  ,W3C0,R3C0 效果是一样的
    yokogawaLinkServer.Write( "D100", new ushort[] { 12242, 42321, 12323 } );                                 // 写入D100  ushort值
    yokogawaLinkServer.Write( "D100", new int[] { 1234312312, 12312312, -1237213 } );                         // 写入D100  int值
    yokogawaLinkServer.Write( "D100", new uint[] { 523123212, 213, 13123 } );                                 // 写入D100  uint值
    yokogawaLinkServer.Write( "D100", new float[] { 123.456f, 35.3f, -675.2f } );                             // 写入D100  float值
    yokogawaLinkServer.Write( "D100", new double[] { 12343.542312d, 213123.123d, -231232.53432d } );          // 写入D100  double值
    yokogawaLinkServer.Write( "D100", new long[] { 1231231242312, 34312312323214, -1283862312631823 } );      // 写入D100  long值
    yokogawaLinkServer.Write( "D100", new ulong[] { 1231231242312, 34312312323214, 9731283862312631823 } );   // 写入D100  ulong值

}
```

高级的对于byte数组类型的数据进行批量化的读写操作如下：

字节的读写示例

[复制](# "复制")

```
public void ReadExample2( )
{

    OperateResult<byte[]> read = yokogawaLinkServer.Read( "D100", 8 );
    if (read.IsSuccess)
    {
        float temp = yokogawaLinkServer.ByteTransform.TransInt16( read.Content, 0 ) / 10f;
        float press = yokogawaLinkServer.ByteTransform.TransInt16( read.Content, 2 ) / 100f;
        int count = yokogawaLinkServer.ByteTransform.TransInt32( read.Content, 2 );

        // do something
    }
    else
    {
        // failed
    }
}

public void WriteExample2( )
{
    // 拼凑数据，这样的话，一次通讯就完成数据的全部写入
    byte[] buffer = new byte[8];
    yokogawaLinkServer.ByteTransform.TransByte( (short)1234 ).CopyTo( buffer, 0 );
    yokogawaLinkServer.ByteTransform.TransByte( (short)2100 ).CopyTo( buffer, 2 );
    yokogawaLinkServer.ByteTransform.TransByte( 12353423 ).CopyTo( buffer, 4 );

    OperateResult write = yokogawaLinkServer.Write( "D100", buffer );
    if (write.IsSuccess)
    {
        // success
    }
    else
    {
        // failed
    }

    // 上面的功能等同于三个数据分别写入，下面的方式性能稍微差一点点，几乎看不出来
    // yokogawaLinkServer.Write( "D100", (short)1234 );
    // yokogawaLinkServer.Write( "D100", (short)2100 );
    // yokogawaLinkServer.Write( "D100", 12353423 );

}
```

更高级操作请参见源代码。

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.Profinet.Yokogawa 命名空间](e985e513-76fd-7ba4-62ba-2b2b3d61b929.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## YokogawaLinkServer 构造函数 

[原文連結](http://api.hslcommunication.cn/html/30b1499f-a533-626f-81ae-4a3051d664b2.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Yokogawa](../html/e985e513-76fd-7ba4-62ba-2b2b3d61b929.htm "HslCommunication.Profinet.Yokogawa")

[YokogawaLinkServer 类](../html/ab5927f5-b6a1-d76a-7c06-10cbcca7ac37.htm "YokogawaLinkServer 类")

[YokogawaLinkServer 构造函数](../html/30b1499f-a533-626f-81ae-4a3051d664b2.htm "YokogawaLinkServer 构造函数 ")

[YokogawaLinkServer 属性](../html/c56bc036-ab9b-8156-d3d2-430f68508349.htm "YokogawaLinkServer 属性")

[YokogawaLinkServer 方法](../html/710635c8-b6db-c84f-fafd-1a8fb1d3b10f.htm "YokogawaLinkServer 方法")

[YokogawaLinkServer 事件](../html/459eec4e-5336-2966-0b53-80009ba8ea32.htm "YokogawaLinkServer 事件")

[YokogawaLinkServer 字段](../html/ae5eb6e5-ad9c-dafe-cb5d-1b5e6c7d05bd.htm "YokogawaLinkServer 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| YokogawaLinkServer 构造函数 |

实例化一个横河PLC的服务器，支持X,Y,I,E,M,T,C,L继电器类型的数据读写，支持D,B,F,R,V,Z,W,TN,CN寄存器类型的数据读写  
Instantiate a Yokogawa PLC server, support X, Y, I, E, M, T, C, L relay type data read and write,
support D, B, F, R, V, Z, W, TN, CN Register type data reading and writing

**命名空间：**
 [HslCommunication.Profinet.Yokogawa](e985e513-76fd-7ba4-62ba-2b2b3d61b929.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public YokogawaLinkServer()
```

```
Public Sub New
```

```
public:
YokogawaLinkServer()
```

```
new : unit -> YokogawaLinkServer
```

![](../icons/SectionExpanded.png)参见

#### 引用

[YokogawaLinkServer 类](ab5927f5-b6a1-d76a-7c06-10cbcca7ac37.htm)

[HslCommunication.Profinet.Yokogawa 命名空间](e985e513-76fd-7ba4-62ba-2b2b3d61b929.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## YokogawaLinkServer 属性

[原文連結](http://api.hslcommunication.cn/html/c56bc036-ab9b-8156-d3d2-430f68508349.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Yokogawa](../html/e985e513-76fd-7ba4-62ba-2b2b3d61b929.htm "HslCommunication.Profinet.Yokogawa")

[YokogawaLinkServer 类](../html/ab5927f5-b6a1-d76a-7c06-10cbcca7ac37.htm "YokogawaLinkServer 类")

[YokogawaLinkServer 构造函数](../html/30b1499f-a533-626f-81ae-4a3051d664b2.htm "YokogawaLinkServer 构造函数 ")

[YokogawaLinkServer 属性](../html/c56bc036-ab9b-8156-d3d2-430f68508349.htm "YokogawaLinkServer 属性")

[YokogawaLinkServer 方法](../html/710635c8-b6db-c84f-fafd-1a8fb1d3b10f.htm "YokogawaLinkServer 方法")

[YokogawaLinkServer 事件](../html/459eec4e-5336-2966-0b53-80009ba8ea32.htm "YokogawaLinkServer 事件")

[YokogawaLinkServer 字段](../html/ae5eb6e5-ad9c-dafe-cb5d-1b5e6c7d05bd.htm "YokogawaLinkServer 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| YokogawaLinkServer 属性 |

[YokogawaLinkServer](ab5927f5-b6a1-d76a-7c06-10cbcca7ac37.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [ActiveTimeSpan](1e32f805-ec96-2746-dc55-44e253ee8514.htm) | 获取或设置两次数据交互时的最小时间间隔，默认为24小时。如果超过该设定的时间不进行数据交互，服务器就会强制断开当前的连接操作。  Get or set the minimum time interval between two data interactions, the default is 24 hours. If the data exchange is not performed for more than the set time, the server will forcibly disconnect the current connection operation. (继承自 [DeviceServer](bad891c2-0527-acc1-3bcb-0eecb1a393a2.htm)。) |
| 公共属性 | [BothModeUdpPort](a2b83e02-9594-1480-2698-f5f0d23e53ca.htm) | 当服务器同时启动TCP及UDP服务的时候，获取当前的UDP服务的端口号  When the server starts TCP and UDP services at the same time, it obtains the port number of the current UDP service (继承自 [DeviceServer](bad891c2-0527-acc1-3bcb-0eecb1a393a2.htm)。) |
| 公共属性代码示例 | [ByteTransform](f35118c9-3691-f752-ae6f-2c5549714faa.htm) | 当前的数据变换机制，当你需要从字节数据转换类型数据的时候需要。  The current data transformation mechanism is required when you need to convert type data from byte data. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共属性 | [CommunicationPipe](ad8e903b-d289-21b0-077e-15ede8005221.htm) | 获取或设置当前的管道信息，管道类型为[CommunicationPipe](ad8e903b-d289-21b0-077e-15ede8005221.htm)的继承类，内置了[PipeTcpNet](26b94bd5-2695-36e7-6378-c3d41846be8f.htm)管道，[PipeUdpNet](5a0b810e-187e-826b-0a13-34fd5576aa0e.htm)管道，[PipeSerialPort](9d924607-3dfc-3916-b532-c5dec48fe347.htm)管道等  Get or set the current pipeline information, the pipeline type is [CommunicationPipe](ad8e903b-d289-21b0-077e-15ede8005221.htm) inheritance class, [PipeTcpNet](26b94bd5-2695-36e7-6378-c3d41846be8f.htm) pipeline, [PipeUdpNet](5a0b810e-187e-826b-0a13-34fd5576aa0e.htm) pipeline, [PipeSerialPort](9d924607-3dfc-3916-b532-c5dec48fe347.htm) pipeline, etc (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [ConnectionId](ce1efb4b-fbc9-4683-163e-0a8e71e99d6b.htm) | 当前连接的唯一ID号，默认为长度20的guid码加随机数组成，方便列表管理，也可以自己指定  The unique ID number of the current connection. The default is a 20-digit guid code plus a random number. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [EnableIPv6](12e23164-745f-c405-453f-26b5177ac5d4.htm) | 获取或设置服务器是否支持IPv6的地址协议信息  Get or set whether the server supports IPv6 address protocol information (继承自 [DeviceServer](bad891c2-0527-acc1-3bcb-0eecb1a393a2.htm)。) |
| 公共属性 | [EnableWrite](e0eeb367-12e3-cf66-19ee-e5938dffbc3e.htm) | 获取或设置当前的服务器是否允许远程客户端进行写入数据操作，默认为True  Gets or sets whether the current server allows remote clients to write data, the default is True (继承自 [DeviceServer](bad891c2-0527-acc1-3bcb-0eecb1a393a2.htm)。) |
| 公共属性 | [ForceSerialReceiveOnce](54cca9f0-f9cd-4743-6bfe-ffc98f5b3163.htm) | 获取或设置当前的服务器接收串口数据时候，是否强制只接收一次数据，默认为false，适合点对点通信，如果你总线形式的连接，则需要设置 True  Get or set whether to force the data to be received only once when the current server receives serial port data. The default value is false, which is suitable for point-to-point communication. If you have a bus connection, you need to set True (继承自 [DeviceServer](bad891c2-0527-acc1-3bcb-0eecb1a393a2.htm)。) |
| 公共属性 | [IsStarted](fbe63af3-720e-11d0-e808-51ff135541d3.htm) | 服务器引擎是否启动  Whether the server engine is started (继承自 [DeviceServer](bad891c2-0527-acc1-3bcb-0eecb1a393a2.htm)。) |
| 公共属性 | [LocalAddress](2d19d9e8-e035-28b9-8aba-3649ef66cc9a.htm) | 获取或设置服务器绑定的本地IP地址，默认为空，使用本地所有可用的ip地址  Obtain or set the local IP address bound to the server, which is empty by default, and use all available IP addresses locally (继承自 [DeviceServer](bad891c2-0527-acc1-3bcb-0eecb1a393a2.htm)。) |
| 公共属性代码示例 | [LogNet](8c56f246-88ac-afb8-8c06-972fafe4f78a.htm) | 组件的日志工具，支持日志记录，只要实例化后，当前网络的基本信息，就以[DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm)等级进行输出  The component's logging tool supports logging. As long as the instantiation of the basic network information, the output will be output at [DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm) (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [OnlineCount](0c146ed0-b958-3374-5952-8069aaee70d2.htm) | 获取在线的客户端的数量  Get the number of clients online (继承自 [DeviceServer](bad891c2-0527-acc1-3bcb-0eecb1a393a2.htm)。) |
| 公共属性 | [Port](858c6445-e24b-21e4-7c2f-3dff7ccb0cb2.htm) | 获取或设置服务器的端口号，如果是设置，需要在服务器启动前设置完成，才能生效。  Gets or sets the port number of the server. If it is set, it needs to be set before the server starts to take effect. (继承自 [DeviceServer](bad891c2-0527-acc1-3bcb-0eecb1a393a2.htm)。) |
| 公共属性代码示例 | [ReceiveTimeOut](250560a7-332b-ce1b-5eb0-e3abe4ee6012.htm) | 获取或设置接收服务器反馈的时间，如果为负数，则不接收反馈   Gets or sets the time to receive server feedback, and if it is a negative number, does not receive feedback (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [SendBeforeHex](72aa1de4-116f-abfe-4adb-814a0be0e560.htm) | 获取或设置在发送通信报文前追加发送的字节信息，HEX格式，通常用于lora组网时，需要携带 00 00 00 02 四个字节的站地址功能。  Obtain or set the byte information sent before sending communication packets, HEX format, usually used for LORA networking, you need to carry 00 00 00 02 four-byte station address function. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [SerialReceiveAtleastTime](75be7ae9-07a9-46a6-0d31-a978821490c8.htm) | 获取或设置串口模式下，接收一条数据最短的时间要求，当设备发送的数据非常慢的时候，或是分割发送数据的时候，就需要将本值设置的大一点，默认为20ms  Get or set the shortest time required to receive a piece of data in serial port mode. When the data sent by the device is very slow, or when the data is divided and sent, you need to set this value to a larger value, the default is 20ms (继承自 [DeviceServer](bad891c2-0527-acc1-3bcb-0eecb1a393a2.htm)。) |
| 公共属性 | [ServerMode](2a437b8e-94ea-1b29-44b9-a1d22fa6aa03.htm) | 当前服务器的模式，0：TCP服务器，1：UDP服务器，2：TCP及UDP服务器  Gets whether the current server is a TCP server or a UDP server (继承自 [DeviceServer](bad891c2-0527-acc1-3bcb-0eecb1a393a2.htm)。) |
| 公共属性 | [SleepTime](86a5fbb9-3456-6beb-7113-5bbad5eabbc5.htm) | 获取或设置在正式接收对方返回数据前的时候，需要休息的时间，当设置为0的时候，不需要休息。  Get or set the time required to rest before officially receiving the data from the other party. When it is set to 0, no rest is required. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [SocketKeepAliveTime](6ba850ca-d152-1045-fbd9-39bc60410ab6.htm) | 获取或设置客户端的Socket的心跳时间信息，这个是Socket底层自动实现的心跳包，不基于协议层实现。默认小于0，不开启心跳检测，如果需要开启，设置 60\_000 比较合适，单位毫秒  Get or set the heartbeat time information of the Socket of the client. This is the heartbeat packet automatically implemented by the bottom layer of the Socket, not based on the protocol layer. The default value is less than 0, and heartbeat detection is not enabled. If you need to enable it, it is more appropriate to set 60\_000, in milliseconds. (继承自 [DeviceServer](bad891c2-0527-acc1-3bcb-0eecb1a393a2.htm)。) |
| 受保护的属性 | [WordLength](f1a505cc-e85c-6c8b-b189-58a1ab5753f9.htm) | 一个字单位的数据表示的地址长度，西门子为2，三菱，欧姆龙，modbusTcp就为1，AB PLC无效  The address length represented by one word of data, Siemens is 2, Mitsubishi, Omron, modbusTcp is 1, AB PLC is invalid (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[YokogawaLinkServer 类](ab5927f5-b6a1-d76a-7c06-10cbcca7ac37.htm)

[HslCommunication.Profinet.Yokogawa 命名空间](e985e513-76fd-7ba4-62ba-2b2b3d61b929.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## YokogawaLinkServer 方法

[原文連結](http://api.hslcommunication.cn/html/710635c8-b6db-c84f-fafd-1a8fb1d3b10f.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Yokogawa](../html/e985e513-76fd-7ba4-62ba-2b2b3d61b929.htm "HslCommunication.Profinet.Yokogawa")

[YokogawaLinkServer 类](../html/ab5927f5-b6a1-d76a-7c06-10cbcca7ac37.htm "YokogawaLinkServer 类")

[YokogawaLinkServer 方法](../html/710635c8-b6db-c84f-fafd-1a8fb1d3b10f.htm "YokogawaLinkServer 方法")

[Dispose 方法](../html/b6a495c4-aed3-7e36-99a8-50eaa2d2f3a2.htm "Dispose 方法 ")

[GetNewNetMessage 方法](../html/c8e4092d-6573-61f0-1035-4a9a062719c9.htm "GetNewNetMessage 方法 ")

[LoadFromBytes 方法](../html/b6683fe3-9938-32fe-5511-9dc08ef6f0eb.htm "LoadFromBytes 方法 ")

[Read 方法](../html/e2cb5a1c-d288-f950-13c8-76bd934f7b55.htm "Read 方法 ")

[ReadBool 方法](../html/b6569f0d-0fc5-aa2f-f127-440c80cef640.htm "ReadBool 方法 ")

[ReadFromCoreServer 方法](../html/94ce8884-86d2-0122-7176-8f070a76ea28.htm "ReadFromCoreServer 方法 ")

[SaveToBytes 方法](../html/bbb5fdfe-745d-1fc9-62bb-a401c143934f.htm "SaveToBytes 方法 ")

[StartProgram 方法](../html/3dc02561-e8a1-57cf-6ef6-bb613eb4e6a0.htm "StartProgram 方法 ")

[StopProgram 方法](../html/0a212b25-16a8-da15-a0c9-6c8130126998.htm "StopProgram 方法 ")

[ToString 方法](../html/572f031c-6c28-1c80-4cd0-e0b99d3144de.htm "ToString 方法 ")

[Write 方法](../html/9f7f0851-0088-71ce-da94-f7932fba2f92.htm "Write 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| YokogawaLinkServer 方法 |

[YokogawaLinkServer](ab5927f5-b6a1-d76a-7c06-10cbcca7ac37.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 受保护的方法 | [CheckSerialReceiveDataComplete](f55d7a1e-3d6b-6a34-b0c3-f9e615218479.htm) | 检查串口接收的数据是否完成的方法，如果接收完成，则返回True (继承自 [DeviceServer](bad891c2-0527-acc1-3bcb-0eecb1a393a2.htm)。) |
| 受保护的方法 | [CloseAction](37b80be6-d1d8-dd5c-0ab7-c0f191302c20.htm) | 服务器关闭的时候需要做的事情  Things to do when the server is down (继承自 [DeviceServer](bad891c2-0527-acc1-3bcb-0eecb1a393a2.htm)。) |
| 公共方法 | [CloseSerialSlave](634feaea-15dd-5051-f41d-9b82de27af7b.htm) | 关闭提供从机服务的串口对象  Close the serial port object that provides slave services (继承自 [DeviceServer](bad891c2-0527-acc1-3bcb-0eecb1a393a2.htm)。) |
| 受保护的方法 | [DecideWhetherQAMessage](acc40cda-e310-8537-6685-59b67eeb16ff.htm) | 决定当前的消息是否是用于问答机制返回的消息，默认直接返回 true, 实际的情况需要根据协议进行重写方法  To determine whether the current message is the message returned by the question answering mechanism, the default is true. In actual cases, the rewriting method needs to be performed according to the protocol (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [Dispose](b9dda6bf-d342-254e-689a-93fa0bb265d1.htm) | (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 受保护的方法 | [Dispose(Boolean)](09988646-d3b6-cd16-4fcc-d93ebd932096.htm) | (重写 [DeviceServerDispose(Boolean)](219f18c3-64d3-a044-b798-de6e2bf866ef.htm).) |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | [ExtraAfterReadFromCoreServer](06a49a15-a3ce-5e39-28c7-2b78851eca2e.htm) | 和服务器交互完成的时候调用的方法，可以根据读写结果进行一些额外的操作，具体的操作需要根据实际的需求来重写实现  The method called when the interaction with the server is completed can perform some additional operations based on the read and write results. The specific operations need to be rewritten according to actual needs. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法代码示例 | [ExtraOnDisconnect](de9dbba2-3b2c-1c5e-53c8-4e6f8a44e2c4.htm) | 根据实际的协议选择是否重写本方法，有些协议在断开连接之前，需要发送一些报文来关闭当前的网络通道  Select whether to rewrite this method according to the actual protocol. Some protocols need to send some packets to close the current network channel before disconnecting. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法代码示例 | [ExtraOnDisconnectAsync](210f593e-16d8-77c4-1178-58fbcfe9b27b.htm) | 根据实际的协议选择是否重写本方法，有些协议在断开连接之前，需要发送一些报文来关闭当前的网络通道  Select whether to rewrite this method according to the actual protocol. Some protocols need to send some packets to close the current network channel before disconnecting. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | [GetCommunicationServer](160e612f-d0c3-0bc4-7a2c-7b60d66835cd.htm) | 获取当前的核心服务器信息 (继承自 [DeviceServer](bad891c2-0527-acc1-3bcb-0eecb1a393a2.htm)。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 受保护的方法 | [GetLogTextFromBinary](67f3b2f7-957c-fa0b-8320-6799237157db.htm) | 获取当前的报文进行日志记录的时候，是否使用二进制的格式记录，默认返回 [LogMsgFormatBinary](2ed2ff87-caec-0e44-f08b-e2b015225a58.htm)，重写可以根据session对象分别返回不同记录模式  Whether to log the current packet in binary format, the default return is [LogMsgFormatBinary](2ed2ff87-caec-0e44-f08b-e2b015225a58.htm). If you want to override it, different recording modes can be returned according to session (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [GetNewNetMessage](c8e4092d-6573-61f0-1035-4a9a062719c9.htm) | 获取一个新的消息对象的方法，需要在继承类里面进行重写  The method to get a new message object needs to be overridden in the inheritance class (重写 [BinaryCommunicationGetNewNetMessage](16636a1b-a6f0-9c0f-9d9c-c94e26c67459.htm).) |
| 公共方法 | [GetTrustedClients](0531cf85-32d8-c438-ccce-b544fc954393.htm) | 获取受信任的客户端列表  Get a list of trusted clients (继承自 [DeviceServer](bad891c2-0527-acc1-3bcb-0eecb1a393a2.htm)。) |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法 | [GetWordLength](c27821ea-1739-8eb9-a848-23a00828fec8.htm) | 一个字单位的数据表示的地址长度，西门子为2，三菱，欧姆龙，modbusTcp就为1，AB PLC无效  The address length represented by one word of data, Siemens is 2, Mitsubishi, Omron, modbusTcp is 1, AB PLC is invalid (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 受保护的方法代码示例 | [InitializationOnConnect](4ab6989f-a5b5-05fe-8b44-bfa45824a2a7.htm) | 根据实际的协议选择是否重写本方法，有些协议在创建连接之后，需要进行一些初始化的信号握手，才能最终建立网络通道。  Whether to rewrite this method is based on the actual protocol. Some protocols require some initial signal handshake to establish a network channel after the connection is created. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法代码示例 | [InitializationOnConnectAsync](7da2af04-d4fc-5045-19a4-4b0d1722eef3.htm) | 根据实际的协议选择是否重写本方法，有些协议在创建连接之后，需要进行一些初始化的信号握手，才能最终建立网络通道。  Whether to rewrite this method is based on the actual protocol. Some protocols require some initial signal handshake to establish a network channel after the connection is created. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [LoadDataPool](346b1333-ac23-cdfc-8edd-3a5ed66ab4bc.htm) | 从文件加载数据池信息  Load datapool information from a file (继承自 [DeviceServer](bad891c2-0527-acc1-3bcb-0eecb1a393a2.htm)。) |
| 受保护的方法 | [LoadFromBytes](b6683fe3-9938-32fe-5511-9dc08ef6f0eb.htm) | 从字节数据加载数据信息，需要进行重写方法  Loading data information from byte data requires rewriting method (重写 [DeviceServerLoadFromBytes(Byte)](4a16c815-c1b7-9620-c705-21db6518272b.htm).) |
| 受保护的方法 | [LogRevcMessage(Byte)](dab1657f-1e75-7fe1-7003-ccedfc69855d.htm) | 使用日志记录一个接收的报文信息  Logs are used to record information about a received packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [LogRevcMessage(Byte, PipeSession)](0f51ae45-575e-4228-f978-01b37d6d5df1.htm) | 使用日志记录一个接收的报文信息  Logs are used to record information about a received packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [LogSendMessage(Byte)](4f90eecc-68a2-1ef2-44ef-369f19bed3ec.htm) | 使用日志记录一个发送的报文信息  Logs are used to record information about a send packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [LogSendMessage(Byte, PipeSession)](209caadf-4f0f-40aa-04c2-ebc61f8dab13.htm) | 使用日志记录一个发送的报文信息  Logs are used to record information about a send packet (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法 | [PackCommandWithHeader](d0b2cd76-c306-11fc-38c7-af8ad306cce7.htm) | 对当前的命令进行打包处理，通常是携带命令头内容，标记当前的命令的长度信息，需要进行重写，否则默认不打包  The current command is packaged, usually carrying the content of the command header, marking the length of the current command, and it needs to be rewritten, otherwise it is not packaged by default (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法 | [RaiseDataReceived](18d5bad7-d0b9-b744-a88c-ec268dd951a0.htm) | 触发一个数据接收的事件信息  Event information that triggers a data reception (继承自 [DeviceServer](bad891c2-0527-acc1-3bcb-0eecb1a393a2.htm)。) |
| 受保护的方法 | [RaiseDataSend](f3aafc2a-1633-37d7-6931-80f791b8e82e.htm) | 触发一个数据发送的事件信息  Event information that triggers a data transmission (继承自 [DeviceServer](bad891c2-0527-acc1-3bcb-0eecb1a393a2.htm)。) |
| 公共方法 | [Read(String, UInt16)](89ea561c-d11c-40f1-0ab3-0654d01d6852.htm) | 批量读取字节数组信息，需要指定地址和长度，返回原始的字节数组  Batch read byte array information, need to specify the address and length, return the original byte array (重写 [DeviceCommunicationRead(String, UInt16)](9dd09452-c200-bc65-75ff-69aa34b366d4.htm).) |
| 公共方法代码示例 | [ReadT](95abe064-5594-5ae7-a849-98c53f88fee0.htm) | 读取支持Hsl特性的数据内容，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明，或是官网的详细文档  Read the data content of the Hsl attribute. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadAsync(String, UInt16)](d8c14d37-56dc-2a31-1a19-3e2b5b65197d.htm) | 异步批量读取字节数组信息，需要指定地址和长度，返回原始的字节数组  Asynchronous batch read byte array information, need to specify the address and length, return the original byte array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadAsyncT](0522e740-4bf0-2e31-c4c5-e41384321b37.htm) | 异步读取支持Hsl特性的数据内容，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Asynchronously read the data content of the Hsl attribute. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadBool(String)](9b2faa2e-eaa9-7ace-c49f-7a0318f1df1a.htm) | 读取单个的Boolean数据信息  Read a single Boolean data message (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadBool(String, UInt16)](3ea0f3cf-91d1-b748-194e-fcddc80b2430.htm) | 批量读取Boolean数组信息，需要指定地址和长度，返回Boolean 数组  Batch read Boolean array information, need to specify the address and length, return Boolean array (重写 [DeviceCommunicationReadBool(String, UInt16)](69d0d24f-34a7-5b41-5c1b-619edc7e6b54.htm).) |
| 公共方法 | [ReadBoolAsync(String)](cc37da1f-d673-fe67-78ee-4e465d29a2bd.htm) | 异步读取单个的Boolean数据信息  Asynchronously read a single Boolean data message (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadBoolAsync(String, UInt16)](a45d383d-0095-3ea7-32f8-e75af88d60f8.htm) | 异步批量读取Boolean数组信息，需要指定地址和长度，返回Boolean 数组  Asynchronously batch read Boolean array information, need to specify the address and length, return Boolean array (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
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
| 受保护的方法代码示例 | [ReadFromCoreServer(PipeSession, Byte)](3950b0e0-214c-75f3-71a8-87421282ea1f.htm) | 将数据发送到当前的网络通道中，并从网络通道中接收一个[INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm)指定的完整的报文，网络通道将根据[GetAvailableSocket](e86a6fc1-eb91-8de9-5ff8-799b0ea3786b.htm)方法自动获取，本方法是线程安全的。  Send data to the current network channel and receive a complete message specified by [INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm) from the network channel. The network channel will be automatically obtained according to the [GetAvailableSocket](e86a6fc1-eb91-8de9-5ff8-799b0ea3786b.htm) method This method is thread-safe. (重写 [DeviceServerReadFromCoreServer(PipeSession, Byte)](bad4d1eb-5428-086b-8767-3440239b6538.htm).) |
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
| 公共方法代码示例 | [ReadStructT](49853082-fbad-52e5-d756-0615cedb4b83.htm) | 读取结构体类型的数据，根据结构体自身的定义，读取原始字节数组，然后解析出实际的结构体数据，结构体需要实现[HslStructAttribute](87e427a8-0c7b-43a8-c421-c27264eda90f.htm)特性 (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
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
| 公共方法 | [SaveDataPool](f3d0c3b1-3383-d593-4b20-0bfe8b1f56bf.htm) | 将本系统的数据池数据存储到指定的文件  Store the data pool data of this system to the specified file (继承自 [DeviceServer](bad891c2-0527-acc1-3bcb-0eecb1a393a2.htm)。) |
| 受保护的方法 | [SaveToBytes](bbb5fdfe-745d-1fc9-62bb-a401c143934f.htm) | 将数据信息存储到字节数组去，需要进行重写方法  To store data information into a byte array, a rewrite method is required (重写 [DeviceServerSaveToBytes](871f8a64-6b6d-53ea-143c-479af077c5a3.htm).) |
| 公共方法 | [ServerClose](b919d91b-5ed0-73c5-dff4-1a1bd2f03bae.htm) | 关闭服务器的引擎  Shut down the server's engine (继承自 [DeviceServer](bad891c2-0527-acc1-3bcb-0eecb1a393a2.htm)。) |
| 公共方法 | [ServerStart](9e137e74-16b0-fcf5-97f3-b1369f55edf0.htm) | 使用已经配置好的端口启动服务器的引擎，并且使用TCP模式  Use the configured port to start the server's engine (继承自 [DeviceServer](bad891c2-0527-acc1-3bcb-0eecb1a393a2.htm)。) |
| 公共方法 | [ServerStart(Int32, Boolean)](bcda1e54-84c1-903a-a284-c2a08870e75e.htm) | 指定端口号来启动服务器的引擎  Specify the port number to start the server's engine (继承自 [DeviceServer](bad891c2-0527-acc1-3bcb-0eecb1a393a2.htm)。) |
| 公共方法 | [ServerStart(Int32, Int32)](f8a8fc7f-c557-9fae-f763-7d8e4b6df3ed.htm) | 指定一个TCP端口及UDP端口，同时启动两种模式的服务器  Specify a TCP port and a UDP port to start the server in both modes at the same time (继承自 [DeviceServer](bad891c2-0527-acc1-3bcb-0eecb1a393a2.htm)。) |
| 公共方法 | [SetDtuPipe](f765d1d7-234a-b229-9457-97e25e024dc9.htm) | 将当前的通信对象设置DTU模式，允许传入现成的管道，并返回初始化结果，如果该设备重写了握手报文，就是返回握手结果  Set the current communication object to DTU mode, allow the existing pipe to be passed in, and return the initialization result, if the device rewrites the handshake packet, the handshake result is returned (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [SetDtuPipeAsync](9d1fd12d-b804-3828-a627-a2c0eef1d23f.htm) | 将当前的通信对象设置DTU模式，允许传入现成的管道，并返回初始化结果，如果该设备重写了握手报文，就是返回握手结果  Set the current communication object to DTU mode, allow the existing pipe to be passed in, and return the initialization result, if the device rewrites the handshake packet, the handshake result is returned (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [SetTrustedIpAddress](7c62dce6-ef4b-bc97-11c9-2ec592d36915.htm) | 设置并启动受信任的客户端登录并读写，如果为null，将关闭对客户端的ip验证  Set and start the trusted client login and read and write, if it is null, the client's IP verification will be turned off (继承自 [DeviceServer](bad891c2-0527-acc1-3bcb-0eecb1a393a2.htm)。) |
| 受保护的方法 | [SocketAcceptExtraCheck](3e3d1f6e-95f9-2f55-ca94-f1802aac9fc2.htm) | 当客户端的socket登录的时候额外检查的操作，并返回操作的结果信息。  The operation is additionally checked when the client's socket logs in, and the result information of the operation is returned. (继承自 [DeviceServer](bad891c2-0527-acc1-3bcb-0eecb1a393a2.htm)。) |
| 受保护的方法 | [StartInitialization](5c580365-04f4-736f-0acc-3146cfb94e10.htm) | 服务器启动时额外的初始化信息，可以用于启动一些额外的服务的操作。  The extra initialization information when the server starts can be used to start some additional service operations. (继承自 [DeviceServer](bad891c2-0527-acc1-3bcb-0eecb1a393a2.htm)。) |
| 公共方法 | [StartProgram](3dc02561-e8a1-57cf-6ef6-bb613eb4e6a0.htm) | 如果未执行程序，则开始执行程序  Starts executing a program if it is not being executed |
| 公共方法 | [StartSerialSlave(ActionSerialPort)](311866fc-5e1c-201d-a1e8-e1a5fafe47f7.htm) | 启动串口的从机服务，使用自定义的初始化方法初始化串口的参数  Start the slave service of serial and initialize the parameters of the serial port using a custom initialization method (继承自 [DeviceServer](bad891c2-0527-acc1-3bcb-0eecb1a393a2.htm)。) |
| 公共方法 | [StartSerialSlave(String)](c719d5ed-8fef-bf84-adc3-0e6af138235d.htm) | 启动串口的从机服务，使用默认的参数进行初始化串口，9600波特率，8位数据位，无奇偶校验，1位停止位  Start the slave service of serial, initialize the serial port with default parameters, 9600 baud rate, 8 data bits, no parity, 1 stop bit (继承自 [DeviceServer](bad891c2-0527-acc1-3bcb-0eecb1a393a2.htm)。) |
| 公共方法 | [StartSerialSlave(String, Int32)](1b6bc90f-70f4-4e2d-01bf-9cd3d6cde88e.htm) | 启动串口的从机服务，使用默认的参数进行初始化串口，8位数据位，无奇偶校验，1位停止位  Start the slave service of serial, initialize the serial port with default parameters, 8 data bits, no parity, 1 stop bit (继承自 [DeviceServer](bad891c2-0527-acc1-3bcb-0eecb1a393a2.htm)。) |
| 公共方法 | [StartSerialSlave(String, Int32, Int32, Parity, StopBits)](353ae919-c1df-6dd7-7b08-1a9011fbd40e.htm) | 启动串口的从机服务，使用指定的参数进行初始化串口，指定数据位，指定奇偶校验，指定停止位 (继承自 [DeviceServer](bad891c2-0527-acc1-3bcb-0eecb1a393a2.htm)。) |
| 公共方法 | [StopProgram](0a212b25-16a8-da15-a0c9-6c8130126998.htm) | 停止当前正在执行程序  Stops the executing program. |
| 受保护的方法 | [ThreadPoolLoginAfterClientCheck](097524c4-f7f3-4529-af08-3773195a5794.htm) | 当客户端登录后，在Ip信息的过滤后，然后触发本方法，进行后续的数据接收，处理，并返回相关的数据信息  When the client logs in, after filtering the IP information, this method is then triggered to perform subsequent data reception, processing, and return related data information (继承自 [DeviceServer](bad891c2-0527-acc1-3bcb-0eecb1a393a2.htm)。) |
| 公共方法 | [ToString](572f031c-6c28-1c80-4cd0-e0b99d3144de.htm) | (重写 [DeviceCommunicationToString](5667665a-da16-d227-f350-5b3f0ab9ab04.htm).) |
| 公共方法 | [UnpackResponseContent](208017a7-8da2-33f8-3273-d7e7e850100f.htm) | 根据对方返回的报文命令，对命令进行基本的拆包，例如各种Modbus协议拆包为统一的核心报文，还支持对报文的验证  According to the message command returned by the other party, the command is basically unpacked, for example, various Modbus protocols are unpacked into a unified core message, and the verification of the message is also supported (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [UseSSL(X509Certificate)](3dd75f50-9dd7-58e2-28c3-d8870a6b2b1b.htm) | 使用SSL通信，传递一个证书的对象 (继承自 [DeviceServer](bad891c2-0527-acc1-3bcb-0eecb1a393a2.htm)。) |
| 公共方法 | [UseSSL(String, String)](08755f4f-11f2-e957-58f2-54abf3706213.htm) | 使用SSL通信，传递一个证书的路径，以及证书的密码 (继承自 [DeviceServer](bad891c2-0527-acc1-3bcb-0eecb1a393a2.htm)。) |
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
| 公共方法 | [Write(String, Boolean)](f1d3e757-a2d4-320d-11cc-9853ffc9b2ab.htm) | 批量写入Boolean数组数据，返回是否成功  Batch write Boolean array data, return whether the write was successful (重写 [DeviceCommunicationWrite(String, Boolean)](5b046c82-3d8e-550b-fdf3-41ca5558169b.htm).) |
| 公共方法 | [Write(String, Byte)](6e9f0b54-4ae7-f483-a1e8-70952e632b24.htm) | 写入原始的byte数组数据到指定的地址，返回是否写入成功  Write the original byte array data to the specified address, and return whether the write was successful (重写 [DeviceCommunicationWrite(String, Byte)](b5536dd6-ccde-6883-ac27-bff9e4b37806.htm).) |
| 公共方法代码示例 | [Write(String, String, Int32)](02fe2ca9-49a3-ede4-d1af-8c6da1f644da.htm) | 写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为ASCII  Write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, it is filled with 0 and the encoding is ASCII. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, String, Encoding)](7b602661-902d-2b70-1fdb-ba87ecd307bf.htm) | 写入字符串信息，需要指定的编码信息  Write string information, need to specify the encoding information (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, String, Int32, Encoding)](e8cb30d9-4c5e-6c68-b9af-1b1b6c2dc714.htm) | 写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为指定的编码信息  Write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, then the operation is complemented with 0 , you should specified the encoding information (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteT(T)](9f48853f-fbeb-c37d-25a2-1c76d1dee407.htm) | 写入支持Hsl特性的数据，返回是否写入成功，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Write data that supports the Hsl attribute, and return whether the write was successful. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WriteAsync(String, Boolean)](10fe669c-805e-1e5a-f016-f319ad363e90.htm) | 异步批量写入Boolean数组数据，返回是否成功  Asynchronously batch write Boolean array data, return success (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WriteAsync(String, Boolean)](fe872dae-7ef9-12c1-d883-e92ba8676416.htm) | 异步批量写入Boolean数组数据，返回是否成功  Asynchronously batch write Boolean array data, return success (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Byte)](8f1cec30-8b86-6749-c8f3-7d39168d59df.htm) | 异步写入原始的byte数组数据到指定的地址，返回是否写入成功  Asynchronously writes the original byte array data to the specified address, and returns whether the write was successful (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
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
| 公共方法 | [WriteAsync(String, String, Int32)](b739b983-b03c-8ecf-2713-75f19233ae2f.htm) | 异步写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为ASCII  Asynchronously write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, it is filled with 0 and the encoding is ASCII. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, String, Encoding)](876e9800-0406-6407-bdf8-8ab1cf28df38.htm) | 异步写入字符串信息，需要指定的编码信息  Asynchronously write string information, need to specify the encoding information (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WriteAsync(String, String, Int32, Encoding)](ae9fb2a3-fddb-48eb-e9ad-6cc0ebbac1f6.htm) | 异步写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为指定的编码信息  Asynchronously write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, then the operation is complemented with 0 , you should specified the encoding information (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsyncT(T)](25d96e77-d1b2-918c-33f9-e5e4cc57c9cd.htm) | 异步写入支持Hsl特性的数据，返回是否写入成功，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Asynchronously write data that supports the Hsl attribute, and return whether the write was successful. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteCustomerT](84a37f82-3f67-a0f8-fd3f-0630772565e4.htm) | 写入自定义类型的数据，该类型必须继承自IDataTransfer接口  Write data of a custom type, which must inherit from the IDataTransfer interface (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteCustomerAsyncT](04628397-cfa1-d785-7790-4c02cb377ab6.htm) | 写入自定义类型的数据，该类型必须继承自IDataTransfer接口  Write data of a custom type, which must inherit from the IDataTransfer interface (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[YokogawaLinkServer 类](ab5927f5-b6a1-d76a-7c06-10cbcca7ac37.htm)

[HslCommunication.Profinet.Yokogawa 命名空间](e985e513-76fd-7ba4-62ba-2b2b3d61b929.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Dispose 方法 

[原文連結](http://api.hslcommunication.cn/html/b6a495c4-aed3-7e36-99a8-50eaa2d2f3a2.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Yokogawa](../html/e985e513-76fd-7ba4-62ba-2b2b3d61b929.htm "HslCommunication.Profinet.Yokogawa")

[YokogawaLinkServer 类](../html/ab5927f5-b6a1-d76a-7c06-10cbcca7ac37.htm "YokogawaLinkServer 类")

[YokogawaLinkServer 方法](../html/710635c8-b6db-c84f-fafd-1a8fb1d3b10f.htm "YokogawaLinkServer 方法")

[Dispose 方法](../html/b6a495c4-aed3-7e36-99a8-50eaa2d2f3a2.htm "Dispose 方法 ")

[Dispose 方法 (Boolean)](../html/09988646-d3b6-cd16-4fcc-d93ebd932096.htm "Dispose 方法 (Boolean)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| YokogawaLinkServerDispose 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [Dispose](b9dda6bf-d342-254e-689a-93fa0bb265d1.htm) | (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 受保护的方法 | [Dispose(Boolean)](09988646-d3b6-cd16-4fcc-d93ebd932096.htm) | 释放被 [YokogawaLinkServer](ab5927f5-b6a1-d76a-7c06-10cbcca7ac37.htm) 使用的非托管资源，并且是否托管资源（可选） (重写 [DeviceServerDispose(Boolean)](219f18c3-64d3-a044-b798-de6e2bf866ef.htm).) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[YokogawaLinkServer 类](ab5927f5-b6a1-d76a-7c06-10cbcca7ac37.htm)

[HslCommunication.Profinet.Yokogawa 命名空间](e985e513-76fd-7ba4-62ba-2b2b3d61b929.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Dispose 方法 (Boolean)

[原文連結](http://api.hslcommunication.cn/html/09988646-d3b6-cd16-4fcc-d93ebd932096.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Yokogawa](../html/e985e513-76fd-7ba4-62ba-2b2b3d61b929.htm "HslCommunication.Profinet.Yokogawa")

[YokogawaLinkServer 类](../html/ab5927f5-b6a1-d76a-7c06-10cbcca7ac37.htm "YokogawaLinkServer 类")

[YokogawaLinkServer 方法](../html/710635c8-b6db-c84f-fafd-1a8fb1d3b10f.htm "YokogawaLinkServer 方法")

[Dispose 方法](../html/b6a495c4-aed3-7e36-99a8-50eaa2d2f3a2.htm "Dispose 方法 ")

[Dispose 方法 (Boolean)](../html/09988646-d3b6-cd16-4fcc-d93ebd932096.htm "Dispose 方法 (Boolean)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| YokogawaLinkServerDispose 方法 (Boolean) |

释放被 [YokogawaLinkServer](ab5927f5-b6a1-d76a-7c06-10cbcca7ac37.htm) 使用的非托管资源，并且是否托管资源（可选）

**命名空间：**
 [HslCommunication.Profinet.Yokogawa](e985e513-76fd-7ba4-62ba-2b2b3d61b929.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
protected override void Dispose(
	bool disposing
)
```

```
Protected Overrides Sub Dispose ( 
	disposing As Boolean
)
```

```
protected:
virtual void Dispose(
	bool disposing
) override
```

```
abstract Dispose : 
        disposing : bool -> unit 
override Dispose : 
        disposing : bool -> unit
```

#### 参数

disposing
:   类型：SystemBoolean  
    为 true 则同时释放托管资源和非托管资源；为 false 则只释放非托管资源

![](../icons/SectionExpanded.png)参见

#### 引用

[YokogawaLinkServer 类](ab5927f5-b6a1-d76a-7c06-10cbcca7ac37.htm)

[Dispose 重载](b6a495c4-aed3-7e36-99a8-50eaa2d2f3a2.htm)

[HslCommunication.Profinet.Yokogawa 命名空间](e985e513-76fd-7ba4-62ba-2b2b3d61b929.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GetNewNetMessage 方法 

[原文連結](http://api.hslcommunication.cn/html/c8e4092d-6573-61f0-1035-4a9a062719c9.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Yokogawa](../html/e985e513-76fd-7ba4-62ba-2b2b3d61b929.htm "HslCommunication.Profinet.Yokogawa")

[YokogawaLinkServer 类](../html/ab5927f5-b6a1-d76a-7c06-10cbcca7ac37.htm "YokogawaLinkServer 类")

[YokogawaLinkServer 方法](../html/710635c8-b6db-c84f-fafd-1a8fb1d3b10f.htm "YokogawaLinkServer 方法")

[Dispose 方法](../html/b6a495c4-aed3-7e36-99a8-50eaa2d2f3a2.htm "Dispose 方法 ")

[GetNewNetMessage 方法](../html/c8e4092d-6573-61f0-1035-4a9a062719c9.htm "GetNewNetMessage 方法 ")

[LoadFromBytes 方法](../html/b6683fe3-9938-32fe-5511-9dc08ef6f0eb.htm "LoadFromBytes 方法 ")

[Read 方法](../html/e2cb5a1c-d288-f950-13c8-76bd934f7b55.htm "Read 方法 ")

[ReadBool 方法](../html/b6569f0d-0fc5-aa2f-f127-440c80cef640.htm "ReadBool 方法 ")

[ReadFromCoreServer 方法](../html/94ce8884-86d2-0122-7176-8f070a76ea28.htm "ReadFromCoreServer 方法 ")

[SaveToBytes 方法](../html/bbb5fdfe-745d-1fc9-62bb-a401c143934f.htm "SaveToBytes 方法 ")

[StartProgram 方法](../html/3dc02561-e8a1-57cf-6ef6-bb613eb4e6a0.htm "StartProgram 方法 ")

[StopProgram 方法](../html/0a212b25-16a8-da15-a0c9-6c8130126998.htm "StopProgram 方法 ")

[ToString 方法](../html/572f031c-6c28-1c80-4cd0-e0b99d3144de.htm "ToString 方法 ")

[Write 方法](../html/9f7f0851-0088-71ce-da94-f7932fba2f92.htm "Write 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| YokogawaLinkServerGetNewNetMessage 方法 |

获取一个新的消息对象的方法，需要在继承类里面进行重写  
The method to get a new message object needs to be overridden in the inheritance class

**命名空间：**
 [HslCommunication.Profinet.Yokogawa](e985e513-76fd-7ba4-62ba-2b2b3d61b929.htm)  
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

[YokogawaLinkServer 类](ab5927f5-b6a1-d76a-7c06-10cbcca7ac37.htm)

[HslCommunication.Profinet.Yokogawa 命名空间](e985e513-76fd-7ba4-62ba-2b2b3d61b929.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## LoadFromBytes 方法 

[原文連結](http://api.hslcommunication.cn/html/b6683fe3-9938-32fe-5511-9dc08ef6f0eb.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Yokogawa](../html/e985e513-76fd-7ba4-62ba-2b2b3d61b929.htm "HslCommunication.Profinet.Yokogawa")

[YokogawaLinkServer 类](../html/ab5927f5-b6a1-d76a-7c06-10cbcca7ac37.htm "YokogawaLinkServer 类")

[YokogawaLinkServer 方法](../html/710635c8-b6db-c84f-fafd-1a8fb1d3b10f.htm "YokogawaLinkServer 方法")

[Dispose 方法](../html/b6a495c4-aed3-7e36-99a8-50eaa2d2f3a2.htm "Dispose 方法 ")

[GetNewNetMessage 方法](../html/c8e4092d-6573-61f0-1035-4a9a062719c9.htm "GetNewNetMessage 方法 ")

[LoadFromBytes 方法](../html/b6683fe3-9938-32fe-5511-9dc08ef6f0eb.htm "LoadFromBytes 方法 ")

[Read 方法](../html/e2cb5a1c-d288-f950-13c8-76bd934f7b55.htm "Read 方法 ")

[ReadBool 方法](../html/b6569f0d-0fc5-aa2f-f127-440c80cef640.htm "ReadBool 方法 ")

[ReadFromCoreServer 方法](../html/94ce8884-86d2-0122-7176-8f070a76ea28.htm "ReadFromCoreServer 方法 ")

[SaveToBytes 方法](../html/bbb5fdfe-745d-1fc9-62bb-a401c143934f.htm "SaveToBytes 方法 ")

[StartProgram 方法](../html/3dc02561-e8a1-57cf-6ef6-bb613eb4e6a0.htm "StartProgram 方法 ")

[StopProgram 方法](../html/0a212b25-16a8-da15-a0c9-6c8130126998.htm "StopProgram 方法 ")

[ToString 方法](../html/572f031c-6c28-1c80-4cd0-e0b99d3144de.htm "ToString 方法 ")

[Write 方法](../html/9f7f0851-0088-71ce-da94-f7932fba2f92.htm "Write 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| YokogawaLinkServerLoadFromBytes 方法 |

从字节数据加载数据信息，需要进行重写方法  
Loading data information from byte data requires rewriting method

**命名空间：**
 [HslCommunication.Profinet.Yokogawa](e985e513-76fd-7ba4-62ba-2b2b3d61b929.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
protected override void LoadFromBytes(
	byte[] content
)
```

```
Protected Overrides Sub LoadFromBytes ( 
	content As Byte()
)
```

```
protected:
virtual void LoadFromBytes(
	array<unsigned char>^ content
) override
```

```
abstract LoadFromBytes : 
        content : byte[] -> unit 
override LoadFromBytes : 
        content : byte[] -> unit
```

#### 参数

content
:   类型：SystemByte  
    字节数据

![](../icons/SectionExpanded.png)参见

#### 引用

[YokogawaLinkServer 类](ab5927f5-b6a1-d76a-7c06-10cbcca7ac37.htm)

[HslCommunication.Profinet.Yokogawa 命名空间](e985e513-76fd-7ba4-62ba-2b2b3d61b929.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Read 方法 

[原文連結](http://api.hslcommunication.cn/html/e2cb5a1c-d288-f950-13c8-76bd934f7b55.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Yokogawa](../html/e985e513-76fd-7ba4-62ba-2b2b3d61b929.htm "HslCommunication.Profinet.Yokogawa")

[YokogawaLinkServer 类](../html/ab5927f5-b6a1-d76a-7c06-10cbcca7ac37.htm "YokogawaLinkServer 类")

[YokogawaLinkServer 方法](../html/710635c8-b6db-c84f-fafd-1a8fb1d3b10f.htm "YokogawaLinkServer 方法")

[Read 方法](../html/e2cb5a1c-d288-f950-13c8-76bd934f7b55.htm "Read 方法 ")

[Read 方法 (String, UInt16)](../html/89ea561c-d11c-40f1-0ab3-0654d01d6852.htm "Read 方法 (String, UInt16)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| YokogawaLinkServerRead 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法代码示例 | [ReadT](95abe064-5594-5ae7-a849-98c53f88fee0.htm) | 读取支持Hsl特性的数据内容，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明，或是官网的详细文档  Read the data content of the Hsl attribute. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [Read(String, UInt16)](89ea561c-d11c-40f1-0ab3-0654d01d6852.htm) | 批量读取字节数组信息，需要指定地址和长度，返回原始的字节数组  Batch read byte array information, need to specify the address and length, return the original byte array (重写 [DeviceCommunicationRead(String, UInt16)](9dd09452-c200-bc65-75ff-69aa34b366d4.htm).) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[YokogawaLinkServer 类](ab5927f5-b6a1-d76a-7c06-10cbcca7ac37.htm)

[HslCommunication.Profinet.Yokogawa 命名空间](e985e513-76fd-7ba4-62ba-2b2b3d61b929.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Read 方法 (String, UInt16)

[原文連結](http://api.hslcommunication.cn/html/89ea561c-d11c-40f1-0ab3-0654d01d6852.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Yokogawa](../html/e985e513-76fd-7ba4-62ba-2b2b3d61b929.htm "HslCommunication.Profinet.Yokogawa")

[YokogawaLinkServer 类](../html/ab5927f5-b6a1-d76a-7c06-10cbcca7ac37.htm "YokogawaLinkServer 类")

[YokogawaLinkServer 方法](../html/710635c8-b6db-c84f-fafd-1a8fb1d3b10f.htm "YokogawaLinkServer 方法")

[Read 方法](../html/e2cb5a1c-d288-f950-13c8-76bd934f7b55.htm "Read 方法 ")

[Read 方法 (String, UInt16)](../html/89ea561c-d11c-40f1-0ab3-0654d01d6852.htm "Read 方法 (String, UInt16)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| YokogawaLinkServerRead 方法 (String, UInt16) |

批量读取字节数组信息，需要指定地址和长度，返回原始的字节数组  
Batch read byte array information, need to specify the address and length, return the original byte array

**命名空间：**
 [HslCommunication.Profinet.Yokogawa](e985e513-76fd-7ba4-62ba-2b2b3d61b929.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public override OperateResult<byte[]> Read(
	string address,
	ushort length
)
```

```
Public Overrides Function Read ( 
	address As String,
	length As UShort
) As OperateResult(Of Byte())
```

```
public:
virtual OperateResult<array<unsigned char>^>^ Read(
	String^ address, 
	unsigned short length
) override
```

```
abstract Read : 
        address : string * 
        length : uint16 -> OperateResult<byte[]> 
override Read : 
        address : string * 
        length : uint16 -> OperateResult<byte[]>
```

#### 参数

address
:   类型：SystemString  
    数据地址

length
:   类型：SystemUInt16  
    数据长度

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
带有成功标识的byte[]数组

#### 实现

[IReadWriteNetRead(String, UInt16)](b8d4187f-ab9a-79e9-c7ca-8a01ec0df233.htm)

![](../icons/SectionExpanded.png)备注

读取的线圈地址支持X,Y,I,E,M,T,C,L，寄存器地址支持D,B,F,R,V,Z,W,TN,CN，举例：D100；也可以携带CPU进行访问，举例：cpu=2;D100  
**[商业授权]** 如果想要读取特殊模块的数据，需要使用 **Special:** 开头标记，举例：Special:unit=0;slot=1;100  
The read coil address supports X, Y, I, E, M, T, C, L, and the register address supports D, B, F, R, V, Z, W, TN, CN, for example: D100;
it can also be carried CPU access, for example: cpu=2;D100.   
**[Authorization]** If you want to read the data of a special module, you need to use the **Special:** beginning tag, for example: Special:unit=0;slot=1;100

![](../icons/SectionExpanded.png)参见

#### 引用

[YokogawaLinkServer 类](ab5927f5-b6a1-d76a-7c06-10cbcca7ac37.htm)

[Read 重载](e2cb5a1c-d288-f950-13c8-76bd934f7b55.htm)

[HslCommunication.Profinet.Yokogawa 命名空间](e985e513-76fd-7ba4-62ba-2b2b3d61b929.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadBool 方法 

[原文連結](http://api.hslcommunication.cn/html/b6569f0d-0fc5-aa2f-f127-440c80cef640.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Yokogawa](../html/e985e513-76fd-7ba4-62ba-2b2b3d61b929.htm "HslCommunication.Profinet.Yokogawa")

[YokogawaLinkServer 类](../html/ab5927f5-b6a1-d76a-7c06-10cbcca7ac37.htm "YokogawaLinkServer 类")

[YokogawaLinkServer 方法](../html/710635c8-b6db-c84f-fafd-1a8fb1d3b10f.htm "YokogawaLinkServer 方法")

[ReadBool 方法](../html/b6569f0d-0fc5-aa2f-f127-440c80cef640.htm "ReadBool 方法 ")

[ReadBool 方法 (String, UInt16)](../html/3ea0f3cf-91d1-b748-194e-fcddc80b2430.htm "ReadBool 方法 (String, UInt16)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| YokogawaLinkServerReadBool 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [ReadBool(String)](9b2faa2e-eaa9-7ace-c49f-7a0318f1df1a.htm) | 读取单个的Boolean数据信息  Read a single Boolean data message (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadBool(String, UInt16)](3ea0f3cf-91d1-b748-194e-fcddc80b2430.htm) | 批量读取Boolean数组信息，需要指定地址和长度，返回Boolean 数组  Batch read Boolean array information, need to specify the address and length, return Boolean array (重写 [DeviceCommunicationReadBool(String, UInt16)](69d0d24f-34a7-5b41-5c1b-619edc7e6b54.htm).) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[YokogawaLinkServer 类](ab5927f5-b6a1-d76a-7c06-10cbcca7ac37.htm)

[HslCommunication.Profinet.Yokogawa 命名空间](e985e513-76fd-7ba4-62ba-2b2b3d61b929.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadBool 方法 (String, UInt16)

[原文連結](http://api.hslcommunication.cn/html/3ea0f3cf-91d1-b748-194e-fcddc80b2430.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Yokogawa](../html/e985e513-76fd-7ba4-62ba-2b2b3d61b929.htm "HslCommunication.Profinet.Yokogawa")

[YokogawaLinkServer 类](../html/ab5927f5-b6a1-d76a-7c06-10cbcca7ac37.htm "YokogawaLinkServer 类")

[YokogawaLinkServer 方法](../html/710635c8-b6db-c84f-fafd-1a8fb1d3b10f.htm "YokogawaLinkServer 方法")

[ReadBool 方法](../html/b6569f0d-0fc5-aa2f-f127-440c80cef640.htm "ReadBool 方法 ")

[ReadBool 方法 (String, UInt16)](../html/3ea0f3cf-91d1-b748-194e-fcddc80b2430.htm "ReadBool 方法 (String, UInt16)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| YokogawaLinkServerReadBool 方法 (String, UInt16) |

批量读取Boolean数组信息，需要指定地址和长度，返回Boolean 数组  
Batch read Boolean array information, need to specify the address and length, return Boolean array

**命名空间：**
 [HslCommunication.Profinet.Yokogawa](e985e513-76fd-7ba4-62ba-2b2b3d61b929.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public override OperateResult<bool[]> ReadBool(
	string address,
	ushort length
)
```

```
Public Overrides Function ReadBool ( 
	address As String,
	length As UShort
) As OperateResult(Of Boolean())
```

```
public:
virtual OperateResult<array<bool>^>^ ReadBool(
	String^ address, 
	unsigned short length
) override
```

```
abstract ReadBool : 
        address : string * 
        length : uint16 -> OperateResult<bool[]> 
override ReadBool : 
        address : string * 
        length : uint16 -> OperateResult<bool[]>
```

#### 参数

address
:   类型：SystemString  
    数据地址

length
:   类型：SystemUInt16  
    数据长度

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Boolean  
带有成功标识的 bool[] 数组

#### 实现

[IReadWriteNetReadBool(String, UInt16)](16e00e01-28d8-bfb7-e3bb-9d8b8778d7b9.htm)

![](../icons/SectionExpanded.png)备注

读取的线圈地址支持X,Y,I,E,M,T,C,L，举例：Y100；也可以携带CPU进行访问，举例：cpu=2;Y100  
The read coil address supports X, Y, I, E, M, T, C, L, for example: Y100; you can also carry the CPU for access, for example: cpu=2;Y100

![](../icons/SectionExpanded.png)参见

#### 引用

[YokogawaLinkServer 类](ab5927f5-b6a1-d76a-7c06-10cbcca7ac37.htm)

[ReadBool 重载](b6569f0d-0fc5-aa2f-f127-440c80cef640.htm)

[HslCommunication.Profinet.Yokogawa 命名空间](e985e513-76fd-7ba4-62ba-2b2b3d61b929.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadFromCoreServer 方法 

[原文連結](http://api.hslcommunication.cn/html/94ce8884-86d2-0122-7176-8f070a76ea28.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Yokogawa](../html/e985e513-76fd-7ba4-62ba-2b2b3d61b929.htm "HslCommunication.Profinet.Yokogawa")

[YokogawaLinkServer 类](../html/ab5927f5-b6a1-d76a-7c06-10cbcca7ac37.htm "YokogawaLinkServer 类")

[YokogawaLinkServer 方法](../html/710635c8-b6db-c84f-fafd-1a8fb1d3b10f.htm "YokogawaLinkServer 方法")

[ReadFromCoreServer 方法](../html/94ce8884-86d2-0122-7176-8f070a76ea28.htm "ReadFromCoreServer 方法 ")

[ReadFromCoreServer 方法 (PipeSession, Byte[])](../html/3950b0e0-214c-75f3-71a8-87421282ea1f.htm "ReadFromCoreServer 方法 (PipeSession, Byte[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| YokogawaLinkServerReadFromCoreServer 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法代码示例 | [ReadFromCoreServer(Byte)](8146dc0c-1edc-3ca7-428a-a3bdbd35e6c5.htm) | 将二进制的数据发送到管道中去，然后从管道里接收二进制的数据回来，并返回是否成功的结果对象。  Send binary data to the pipeline, and then receive binary data back from the pipeline, and return whether the success of the result object (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServer(IEnumerableByte)](044cfdd9-ec76-b89f-2f2b-5f439025006c.htm) | 将多个数据报文按顺序发到设备，并从设备接收返回的数据内容，然后拼接成一个Byte[]信息，需要重写[UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm)方法才能返回正确的结果。  Send multiple data packets to the device in sequence, and receive the returned data content from the device, and then splicing them into a Byte[] message, you need to rewrite [UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm) method to return the correct result. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 受保护的方法代码示例 | [ReadFromCoreServer(PipeSession, Byte)](3950b0e0-214c-75f3-71a8-87421282ea1f.htm) | 将数据发送到当前的网络通道中，并从网络通道中接收一个[INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm)指定的完整的报文，网络通道将根据[GetAvailableSocket](e86a6fc1-eb91-8de9-5ff8-799b0ea3786b.htm)方法自动获取，本方法是线程安全的。  Send data to the current network channel and receive a complete message specified by [INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm) from the network channel. The network channel will be automatically obtained according to the [GetAvailableSocket](e86a6fc1-eb91-8de9-5ff8-799b0ea3786b.htm) method This method is thread-safe. (重写 [DeviceServerReadFromCoreServer(PipeSession, Byte)](bad4d1eb-5428-086b-8767-3440239b6538.htm).) |
| 公共方法代码示例 | [ReadFromCoreServer(Byte, Boolean, Boolean)](b32ce8e8-a959-fcc2-6736-14b39fdb9a40.htm) | 将二进制的数据发送到管道中去，然后从管道里接收二进制的数据回来，并返回是否成功的结果对象。  Send binary data to the pipeline, and then receive binary data back from the pipeline, and return whether the success of the result object (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [ReadFromCoreServer(CommunicationPipe, Byte, Boolean, Boolean)](ef36d56b-df09-f727-e8d3-316d37a13ee9.htm) | 使用指定的管道来进行数据通信，发送原始数据到管道，然后从管道接收相关的数据返回，本方法无锁 (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[YokogawaLinkServer 类](ab5927f5-b6a1-d76a-7c06-10cbcca7ac37.htm)

[HslCommunication.Profinet.Yokogawa 命名空间](e985e513-76fd-7ba4-62ba-2b2b3d61b929.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadFromCoreServer 方法 (PipeSession, Byte[])

[原文連結](http://api.hslcommunication.cn/html/3950b0e0-214c-75f3-71a8-87421282ea1f.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Yokogawa](../html/e985e513-76fd-7ba4-62ba-2b2b3d61b929.htm "HslCommunication.Profinet.Yokogawa")

[YokogawaLinkServer 类](../html/ab5927f5-b6a1-d76a-7c06-10cbcca7ac37.htm "YokogawaLinkServer 类")

[YokogawaLinkServer 方法](../html/710635c8-b6db-c84f-fafd-1a8fb1d3b10f.htm "YokogawaLinkServer 方法")

[ReadFromCoreServer 方法](../html/94ce8884-86d2-0122-7176-8f070a76ea28.htm "ReadFromCoreServer 方法 ")

[ReadFromCoreServer 方法 (PipeSession, Byte[])](../html/3950b0e0-214c-75f3-71a8-87421282ea1f.htm "ReadFromCoreServer 方法 (PipeSession, Byte[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| YokogawaLinkServerReadFromCoreServer 方法 (PipeSession, Byte) |

将数据发送到当前的网络通道中，并从网络通道中接收一个[INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm)指定的完整的报文，网络通道将根据[GetAvailableSocket](e86a6fc1-eb91-8de9-5ff8-799b0ea3786b.htm)方法自动获取，本方法是线程安全的。  
Send data to the current network channel and receive a complete message specified by [INetMessage](a37971ba-f7b5-5ce1-ae41-f1334c2cba3c.htm) from the network channel.
The network channel will be automatically obtained according to the [GetAvailableSocket](e86a6fc1-eb91-8de9-5ff8-799b0ea3786b.htm) method This method is thread-safe.

**命名空间：**
 [HslCommunication.Profinet.Yokogawa](e985e513-76fd-7ba4-62ba-2b2b3d61b929.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
protected override OperateResult<byte[]> ReadFromCoreServer(
	PipeSession session,
	byte[] receive
)
```

```
Protected Overrides Function ReadFromCoreServer ( 
	session As PipeSession,
	receive As Byte()
) As OperateResult(Of Byte())
```

```
protected:
virtual OperateResult<array<unsigned char>^>^ ReadFromCoreServer(
	PipeSession^ session, 
	array<unsigned char>^ receive
) override
```

```
abstract ReadFromCoreServer : 
        session : PipeSession * 
        receive : byte[] -> OperateResult<byte[]> 
override ReadFromCoreServer : 
        session : PipeSession * 
        receive : byte[] -> OperateResult<byte[]>
```

#### 参数

session
:   类型：[HslCommunication.Core.NetPipeSession](bb565655-08a0-c9eb-9c75-26ad6a263256.htm)  

    [缺少 "M:HslCommunication.Profinet.Yokogawa.YokogawaLinkServer.ReadFromCoreServer(HslCommunication.Core.Net.PipeSession,System.Byte[])" 的 <param name="session"/> 文档]

receive
:   类型：SystemByte  

    [缺少 "M:HslCommunication.Profinet.Yokogawa.YokogawaLinkServer.ReadFromCoreServer(HslCommunication.Core.Net.PipeSession,System.Byte[])" 的 <param name="receive"/> 文档]

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
接收的完整的报文信息

![](../icons/SectionExpanded.png)备注

本方法用于实现本组件还未实现的一些报文功能，例如有些modbus服务器会有一些特殊的功能码支持，需要收发特殊的报文，详细请看示例

![](../icons/SectionExpanded.png)示例

此处举例有个modbus服务器，有个特殊的功能码0x09，后面携带子数据0x01即可，发送字节为 0x00 0x00 0x00 0x00 0x00 0x03 0x01 0x09 0x01

ReadFromCoreServer示例

[复制](# "复制")

```
ModbusTcpNet modbus = new ModbusTcpNet( "192.168.0.100" );

// 此处举例实现特殊的modbus功能码
OperateResult<byte[]> read = modbus.ReadFromCoreServer( SoftBasic.HexStringToBytes( "0x00 0x00 0x00 0x00 0x00 0x03 0x01 0x09 0x01" ) );
if (read.IsSuccess)
{
    // 成功，开始解析从服务器返回的数据，是一条完整的报文信息
    Console.WriteLine( SoftBasic.ByteToHexString( read.Content, ' ' ) );
}
else
{
    // 失败
}
```

![](../icons/SectionExpanded.png)参见

#### 引用

[YokogawaLinkServer 类](ab5927f5-b6a1-d76a-7c06-10cbcca7ac37.htm)

[ReadFromCoreServer 重载](94ce8884-86d2-0122-7176-8f070a76ea28.htm)

[HslCommunication.Profinet.Yokogawa 命名空间](e985e513-76fd-7ba4-62ba-2b2b3d61b929.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## SaveToBytes 方法 

[原文連結](http://api.hslcommunication.cn/html/bbb5fdfe-745d-1fc9-62bb-a401c143934f.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Yokogawa](../html/e985e513-76fd-7ba4-62ba-2b2b3d61b929.htm "HslCommunication.Profinet.Yokogawa")

[YokogawaLinkServer 类](../html/ab5927f5-b6a1-d76a-7c06-10cbcca7ac37.htm "YokogawaLinkServer 类")

[YokogawaLinkServer 方法](../html/710635c8-b6db-c84f-fafd-1a8fb1d3b10f.htm "YokogawaLinkServer 方法")

[Dispose 方法](../html/b6a495c4-aed3-7e36-99a8-50eaa2d2f3a2.htm "Dispose 方法 ")

[GetNewNetMessage 方法](../html/c8e4092d-6573-61f0-1035-4a9a062719c9.htm "GetNewNetMessage 方法 ")

[LoadFromBytes 方法](../html/b6683fe3-9938-32fe-5511-9dc08ef6f0eb.htm "LoadFromBytes 方法 ")

[Read 方法](../html/e2cb5a1c-d288-f950-13c8-76bd934f7b55.htm "Read 方法 ")

[ReadBool 方法](../html/b6569f0d-0fc5-aa2f-f127-440c80cef640.htm "ReadBool 方法 ")

[ReadFromCoreServer 方法](../html/94ce8884-86d2-0122-7176-8f070a76ea28.htm "ReadFromCoreServer 方法 ")

[SaveToBytes 方法](../html/bbb5fdfe-745d-1fc9-62bb-a401c143934f.htm "SaveToBytes 方法 ")

[StartProgram 方法](../html/3dc02561-e8a1-57cf-6ef6-bb613eb4e6a0.htm "StartProgram 方法 ")

[StopProgram 方法](../html/0a212b25-16a8-da15-a0c9-6c8130126998.htm "StopProgram 方法 ")

[ToString 方法](../html/572f031c-6c28-1c80-4cd0-e0b99d3144de.htm "ToString 方法 ")

[Write 方法](../html/9f7f0851-0088-71ce-da94-f7932fba2f92.htm "Write 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| YokogawaLinkServerSaveToBytes 方法 |

将数据信息存储到字节数组去，需要进行重写方法  
To store data information into a byte array, a rewrite method is required

**命名空间：**
 [HslCommunication.Profinet.Yokogawa](e985e513-76fd-7ba4-62ba-2b2b3d61b929.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
protected override byte[] SaveToBytes()
```

```
Protected Overrides Function SaveToBytes As Byte()
```

```
protected:
virtual array<unsigned char>^ SaveToBytes() override
```

```
abstract SaveToBytes : unit -> byte[] 
override SaveToBytes : unit -> byte[]
```

#### 返回值

类型：Byte  
所有的内容

![](../icons/SectionExpanded.png)参见

#### 引用

[YokogawaLinkServer 类](ab5927f5-b6a1-d76a-7c06-10cbcca7ac37.htm)

[HslCommunication.Profinet.Yokogawa 命名空间](e985e513-76fd-7ba4-62ba-2b2b3d61b929.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## StartProgram 方法 

[原文連結](http://api.hslcommunication.cn/html/3dc02561-e8a1-57cf-6ef6-bb613eb4e6a0.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Yokogawa](../html/e985e513-76fd-7ba4-62ba-2b2b3d61b929.htm "HslCommunication.Profinet.Yokogawa")

[YokogawaLinkServer 类](../html/ab5927f5-b6a1-d76a-7c06-10cbcca7ac37.htm "YokogawaLinkServer 类")

[YokogawaLinkServer 方法](../html/710635c8-b6db-c84f-fafd-1a8fb1d3b10f.htm "YokogawaLinkServer 方法")

[Dispose 方法](../html/b6a495c4-aed3-7e36-99a8-50eaa2d2f3a2.htm "Dispose 方法 ")

[GetNewNetMessage 方法](../html/c8e4092d-6573-61f0-1035-4a9a062719c9.htm "GetNewNetMessage 方法 ")

[LoadFromBytes 方法](../html/b6683fe3-9938-32fe-5511-9dc08ef6f0eb.htm "LoadFromBytes 方法 ")

[Read 方法](../html/e2cb5a1c-d288-f950-13c8-76bd934f7b55.htm "Read 方法 ")

[ReadBool 方法](../html/b6569f0d-0fc5-aa2f-f127-440c80cef640.htm "ReadBool 方法 ")

[ReadFromCoreServer 方法](../html/94ce8884-86d2-0122-7176-8f070a76ea28.htm "ReadFromCoreServer 方法 ")

[SaveToBytes 方法](../html/bbb5fdfe-745d-1fc9-62bb-a401c143934f.htm "SaveToBytes 方法 ")

[StartProgram 方法](../html/3dc02561-e8a1-57cf-6ef6-bb613eb4e6a0.htm "StartProgram 方法 ")

[StopProgram 方法](../html/0a212b25-16a8-da15-a0c9-6c8130126998.htm "StopProgram 方法 ")

[ToString 方法](../html/572f031c-6c28-1c80-4cd0-e0b99d3144de.htm "ToString 方法 ")

[Write 方法](../html/9f7f0851-0088-71ce-da94-f7932fba2f92.htm "Write 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| YokogawaLinkServerStartProgram 方法 |

如果未执行程序，则开始执行程序  
Starts executing a program if it is not being executed

**命名空间：**
 [HslCommunication.Profinet.Yokogawa](e985e513-76fd-7ba4-62ba-2b2b3d61b929.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public void StartProgram()
```

```
Public Sub StartProgram
```

```
public:
void StartProgram()
```

```
member StartProgram : unit -> unit 
```

![](../icons/SectionExpanded.png)参见

#### 引用

[YokogawaLinkServer 类](ab5927f5-b6a1-d76a-7c06-10cbcca7ac37.htm)

[HslCommunication.Profinet.Yokogawa 命名空间](e985e513-76fd-7ba4-62ba-2b2b3d61b929.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## StopProgram 方法 

[原文連結](http://api.hslcommunication.cn/html/0a212b25-16a8-da15-a0c9-6c8130126998.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Yokogawa](../html/e985e513-76fd-7ba4-62ba-2b2b3d61b929.htm "HslCommunication.Profinet.Yokogawa")

[YokogawaLinkServer 类](../html/ab5927f5-b6a1-d76a-7c06-10cbcca7ac37.htm "YokogawaLinkServer 类")

[YokogawaLinkServer 方法](../html/710635c8-b6db-c84f-fafd-1a8fb1d3b10f.htm "YokogawaLinkServer 方法")

[Dispose 方法](../html/b6a495c4-aed3-7e36-99a8-50eaa2d2f3a2.htm "Dispose 方法 ")

[GetNewNetMessage 方法](../html/c8e4092d-6573-61f0-1035-4a9a062719c9.htm "GetNewNetMessage 方法 ")

[LoadFromBytes 方法](../html/b6683fe3-9938-32fe-5511-9dc08ef6f0eb.htm "LoadFromBytes 方法 ")

[Read 方法](../html/e2cb5a1c-d288-f950-13c8-76bd934f7b55.htm "Read 方法 ")

[ReadBool 方法](../html/b6569f0d-0fc5-aa2f-f127-440c80cef640.htm "ReadBool 方法 ")

[ReadFromCoreServer 方法](../html/94ce8884-86d2-0122-7176-8f070a76ea28.htm "ReadFromCoreServer 方法 ")

[SaveToBytes 方法](../html/bbb5fdfe-745d-1fc9-62bb-a401c143934f.htm "SaveToBytes 方法 ")

[StartProgram 方法](../html/3dc02561-e8a1-57cf-6ef6-bb613eb4e6a0.htm "StartProgram 方法 ")

[StopProgram 方法](../html/0a212b25-16a8-da15-a0c9-6c8130126998.htm "StopProgram 方法 ")

[ToString 方法](../html/572f031c-6c28-1c80-4cd0-e0b99d3144de.htm "ToString 方法 ")

[Write 方法](../html/9f7f0851-0088-71ce-da94-f7932fba2f92.htm "Write 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| YokogawaLinkServerStopProgram 方法 |

停止当前正在执行程序  
Stops the executing program.

**命名空间：**
 [HslCommunication.Profinet.Yokogawa](e985e513-76fd-7ba4-62ba-2b2b3d61b929.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public void StopProgram()
```

```
Public Sub StopProgram
```

```
public:
void StopProgram()
```

```
member StopProgram : unit -> unit 
```

![](../icons/SectionExpanded.png)参见

#### 引用

[YokogawaLinkServer 类](ab5927f5-b6a1-d76a-7c06-10cbcca7ac37.htm)

[HslCommunication.Profinet.Yokogawa 命名空间](e985e513-76fd-7ba4-62ba-2b2b3d61b929.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ToString 方法 

[原文連結](http://api.hslcommunication.cn/html/572f031c-6c28-1c80-4cd0-e0b99d3144de.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Yokogawa](../html/e985e513-76fd-7ba4-62ba-2b2b3d61b929.htm "HslCommunication.Profinet.Yokogawa")

[YokogawaLinkServer 类](../html/ab5927f5-b6a1-d76a-7c06-10cbcca7ac37.htm "YokogawaLinkServer 类")

[YokogawaLinkServer 方法](../html/710635c8-b6db-c84f-fafd-1a8fb1d3b10f.htm "YokogawaLinkServer 方法")

[Dispose 方法](../html/b6a495c4-aed3-7e36-99a8-50eaa2d2f3a2.htm "Dispose 方法 ")

[GetNewNetMessage 方法](../html/c8e4092d-6573-61f0-1035-4a9a062719c9.htm "GetNewNetMessage 方法 ")

[LoadFromBytes 方法](../html/b6683fe3-9938-32fe-5511-9dc08ef6f0eb.htm "LoadFromBytes 方法 ")

[Read 方法](../html/e2cb5a1c-d288-f950-13c8-76bd934f7b55.htm "Read 方法 ")

[ReadBool 方法](../html/b6569f0d-0fc5-aa2f-f127-440c80cef640.htm "ReadBool 方法 ")

[ReadFromCoreServer 方法](../html/94ce8884-86d2-0122-7176-8f070a76ea28.htm "ReadFromCoreServer 方法 ")

[SaveToBytes 方法](../html/bbb5fdfe-745d-1fc9-62bb-a401c143934f.htm "SaveToBytes 方法 ")

[StartProgram 方法](../html/3dc02561-e8a1-57cf-6ef6-bb613eb4e6a0.htm "StartProgram 方法 ")

[StopProgram 方法](../html/0a212b25-16a8-da15-a0c9-6c8130126998.htm "StopProgram 方法 ")

[ToString 方法](../html/572f031c-6c28-1c80-4cd0-e0b99d3144de.htm "ToString 方法 ")

[Write 方法](../html/9f7f0851-0088-71ce-da94-f7932fba2f92.htm "Write 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| YokogawaLinkServerToString 方法 |

[缺少 "M:HslCommunication.Profinet.Yokogawa.YokogawaLinkServer.ToString" 的 <summary> 文档]

**命名空间：**
 [HslCommunication.Profinet.Yokogawa](e985e513-76fd-7ba4-62ba-2b2b3d61b929.htm)  
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

[缺少 "M:HslCommunication.Profinet.Yokogawa.YokogawaLinkServer.ToString" 的 <returns> 文档]

![](../icons/SectionExpanded.png)参见

#### 引用

[YokogawaLinkServer 类](ab5927f5-b6a1-d76a-7c06-10cbcca7ac37.htm)

[HslCommunication.Profinet.Yokogawa 命名空间](e985e513-76fd-7ba4-62ba-2b2b3d61b929.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Write 方法 

[原文連結](http://api.hslcommunication.cn/html/9f7f0851-0088-71ce-da94-f7932fba2f92.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Yokogawa](../html/e985e513-76fd-7ba4-62ba-2b2b3d61b929.htm "HslCommunication.Profinet.Yokogawa")

[YokogawaLinkServer 类](../html/ab5927f5-b6a1-d76a-7c06-10cbcca7ac37.htm "YokogawaLinkServer 类")

[YokogawaLinkServer 方法](../html/710635c8-b6db-c84f-fafd-1a8fb1d3b10f.htm "YokogawaLinkServer 方法")

[Write 方法](../html/9f7f0851-0088-71ce-da94-f7932fba2f92.htm "Write 方法 ")

[Write 方法 (String, Boolean[])](../html/f1d3e757-a2d4-320d-11cc-9853ffc9b2ab.htm "Write 方法 (String, Boolean[])")

[Write 方法 (String, Byte[])](../html/6e9f0b54-4ae7-f483-a1e8-70952e632b24.htm "Write 方法 (String, Byte[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| YokogawaLinkServerWrite 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法代码示例 | [WriteT(T)](9f48853f-fbeb-c37d-25a2-1c76d1dee407.htm) | 写入支持Hsl特性的数据，返回是否写入成功，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Write data that supports the Hsl attribute, and return whether the write was successful. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
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
| 公共方法 | [Write(String, Boolean)](f1d3e757-a2d4-320d-11cc-9853ffc9b2ab.htm) | 批量写入Boolean数组数据，返回是否成功  Batch write Boolean array data, return whether the write was successful (重写 [DeviceCommunicationWrite(String, Boolean)](5b046c82-3d8e-550b-fdf3-41ca5558169b.htm).) |
| 公共方法 | [Write(String, Byte)](6e9f0b54-4ae7-f483-a1e8-70952e632b24.htm) | 写入原始的byte数组数据到指定的地址，返回是否写入成功  Write the original byte array data to the specified address, and return whether the write was successful (重写 [DeviceCommunicationWrite(String, Byte)](b5536dd6-ccde-6883-ac27-bff9e4b37806.htm).) |
| 公共方法代码示例 | [Write(String, String, Int32)](02fe2ca9-49a3-ede4-d1af-8c6da1f644da.htm) | 写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为ASCII  Write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, it is filled with 0 and the encoding is ASCII. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, String, Encoding)](7b602661-902d-2b70-1fdb-ba87ecd307bf.htm) | 写入字符串信息，需要指定的编码信息  Write string information, need to specify the encoding information (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, String, Int32, Encoding)](e8cb30d9-4c5e-6c68-b9af-1b1b6c2dc714.htm) | 写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为指定的编码信息  Write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, then the operation is complemented with 0 , you should specified the encoding information (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[YokogawaLinkServer 类](ab5927f5-b6a1-d76a-7c06-10cbcca7ac37.htm)

[HslCommunication.Profinet.Yokogawa 命名空间](e985e513-76fd-7ba4-62ba-2b2b3d61b929.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Write 方法 (String, Boolean[])

[原文連結](http://api.hslcommunication.cn/html/f1d3e757-a2d4-320d-11cc-9853ffc9b2ab.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Yokogawa](../html/e985e513-76fd-7ba4-62ba-2b2b3d61b929.htm "HslCommunication.Profinet.Yokogawa")

[YokogawaLinkServer 类](../html/ab5927f5-b6a1-d76a-7c06-10cbcca7ac37.htm "YokogawaLinkServer 类")

[YokogawaLinkServer 方法](../html/710635c8-b6db-c84f-fafd-1a8fb1d3b10f.htm "YokogawaLinkServer 方法")

[Write 方法](../html/9f7f0851-0088-71ce-da94-f7932fba2f92.htm "Write 方法 ")

[Write 方法 (String, Boolean[])](../html/f1d3e757-a2d4-320d-11cc-9853ffc9b2ab.htm "Write 方法 (String, Boolean[])")

[Write 方法 (String, Byte[])](../html/6e9f0b54-4ae7-f483-a1e8-70952e632b24.htm "Write 方法 (String, Byte[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| YokogawaLinkServerWrite 方法 (String, Boolean) |

批量写入Boolean数组数据，返回是否成功  
Batch write Boolean array data, return whether the write was successful

**命名空间：**
 [HslCommunication.Profinet.Yokogawa](e985e513-76fd-7ba4-62ba-2b2b3d61b929.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public override OperateResult Write(
	string address,
	bool[] value
)
```

```
Public Overrides Function Write ( 
	address As String,
	value As Boolean()
) As OperateResult
```

```
public:
virtual OperateResult^ Write(
	String^ address, 
	array<bool>^ value
) override
```

```
abstract Write : 
        address : string * 
        value : bool[] -> OperateResult 
override Write : 
        address : string * 
        value : bool[] -> OperateResult
```

#### 参数

address
:   类型：SystemString  
    起始地址

value
:   类型：SystemBoolean  
    写入值

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否写入成功

#### 实现

[IReadWriteNetWrite(String, Boolean)](2a99122c-1496-d766-f0aa-44a18918fcd3.htm)

![](../icons/SectionExpanded.png)备注

写入的线圈地址支持Y,I,E,M,T,C,L，举例：Y100；也可以携带CPU进行访问，举例：cpu=2;Y100  
The write coil address supports Y, I, E, M, T, C, L, for example: Y100; you can also carry the CPU for access, for example: cpu=2;Y100

![](../icons/SectionExpanded.png)参见

#### 引用

[YokogawaLinkServer 类](ab5927f5-b6a1-d76a-7c06-10cbcca7ac37.htm)

[Write 重载](9f7f0851-0088-71ce-da94-f7932fba2f92.htm)

[HslCommunication.Profinet.Yokogawa 命名空间](e985e513-76fd-7ba4-62ba-2b2b3d61b929.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Write 方法 (String, Byte[])

[原文連結](http://api.hslcommunication.cn/html/6e9f0b54-4ae7-f483-a1e8-70952e632b24.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Yokogawa](../html/e985e513-76fd-7ba4-62ba-2b2b3d61b929.htm "HslCommunication.Profinet.Yokogawa")

[YokogawaLinkServer 类](../html/ab5927f5-b6a1-d76a-7c06-10cbcca7ac37.htm "YokogawaLinkServer 类")

[YokogawaLinkServer 方法](../html/710635c8-b6db-c84f-fafd-1a8fb1d3b10f.htm "YokogawaLinkServer 方法")

[Write 方法](../html/9f7f0851-0088-71ce-da94-f7932fba2f92.htm "Write 方法 ")

[Write 方法 (String, Boolean[])](../html/f1d3e757-a2d4-320d-11cc-9853ffc9b2ab.htm "Write 方法 (String, Boolean[])")

[Write 方法 (String, Byte[])](../html/6e9f0b54-4ae7-f483-a1e8-70952e632b24.htm "Write 方法 (String, Byte[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| YokogawaLinkServerWrite 方法 (String, Byte) |

写入原始的byte数组数据到指定的地址，返回是否写入成功  
Write the original byte array data to the specified address, and return whether the write was successful

**命名空间：**
 [HslCommunication.Profinet.Yokogawa](e985e513-76fd-7ba4-62ba-2b2b3d61b929.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public override OperateResult Write(
	string address,
	byte[] value
)
```

```
Public Overrides Function Write ( 
	address As String,
	value As Byte()
) As OperateResult
```

```
public:
virtual OperateResult^ Write(
	String^ address, 
	array<unsigned char>^ value
) override
```

```
abstract Write : 
        address : string * 
        value : byte[] -> OperateResult 
override Write : 
        address : string * 
        value : byte[] -> OperateResult
```

#### 参数

address
:   类型：SystemString  
    起始地址

value
:   类型：SystemByte  
    写入值

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否写入成功

#### 实现

[IReadWriteNetWrite(String, Byte)](b8d8d94e-300d-c7b7-2e80-6859c62479fb.htm)

![](../icons/SectionExpanded.png)备注

写入的线圈地址支持Y,I,E,M,T,C,L，寄存器地址支持D,B,F,R,V,Z,W,TN,CN，举例：D100；也可以携带CPU进行访问，举例：cpu=2;D100  
如果想要写入特殊模块的数据，需要使用 **Special:** 开头标记，举例：Special:unit=0;slot=1;100  
The read coil address supports Y, I, E, M, T, C, L, and the register address supports D, B, F, R, V, Z, W, TN, CN, for example: D100;
it can also be carried CPU access, for example: cpu=2;D100.
If you want to read the data of a special module, you need to use the **Special:** beginning tag, for example: Special:unit=0;slot=1;100

![](../icons/SectionExpanded.png)参见

#### 引用

[YokogawaLinkServer 类](ab5927f5-b6a1-d76a-7c06-10cbcca7ac37.htm)

[Write 重载](9f7f0851-0088-71ce-da94-f7932fba2f92.htm)

[HslCommunication.Profinet.Yokogawa 命名空间](e985e513-76fd-7ba4-62ba-2b2b3d61b929.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## YokogawaLinkServer 事件

[原文連結](http://api.hslcommunication.cn/html/459eec4e-5336-2966-0b53-80009ba8ea32.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Yokogawa](../html/e985e513-76fd-7ba4-62ba-2b2b3d61b929.htm "HslCommunication.Profinet.Yokogawa")

[YokogawaLinkServer 类](../html/ab5927f5-b6a1-d76a-7c06-10cbcca7ac37.htm "YokogawaLinkServer 类")

[YokogawaLinkServer 构造函数](../html/30b1499f-a533-626f-81ae-4a3051d664b2.htm "YokogawaLinkServer 构造函数 ")

[YokogawaLinkServer 属性](../html/c56bc036-ab9b-8156-d3d2-430f68508349.htm "YokogawaLinkServer 属性")

[YokogawaLinkServer 方法](../html/710635c8-b6db-c84f-fafd-1a8fb1d3b10f.htm "YokogawaLinkServer 方法")

[YokogawaLinkServer 事件](../html/459eec4e-5336-2966-0b53-80009ba8ea32.htm "YokogawaLinkServer 事件")

[YokogawaLinkServer 字段](../html/ae5eb6e5-ad9c-dafe-cb5d-1b5e6c7d05bd.htm "YokogawaLinkServer 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| YokogawaLinkServer 事件 |

[YokogawaLinkServer](ab5927f5-b6a1-d76a-7c06-10cbcca7ac37.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)事件

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共事件代码示例 | [OnDataReceived](821fdec2-5ea2-1d1c-7f5c-08f594d7d1ba.htm) | 接收到数据的时候就触发的事件，示例详细参考API文档信息  An event that is triggered when data is received (继承自 [DeviceServer](bad891c2-0527-acc1-3bcb-0eecb1a393a2.htm)。) |
| 公共事件 | [OnDataSend](deabf344-282f-a8c2-8d71-8f1512a66e38.htm) | 数据发送的时候就触发的事件  Events that are triggered when data is sent (继承自 [DeviceServer](bad891c2-0527-acc1-3bcb-0eecb1a393a2.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[YokogawaLinkServer 类](ab5927f5-b6a1-d76a-7c06-10cbcca7ac37.htm)

[HslCommunication.Profinet.Yokogawa 命名空间](e985e513-76fd-7ba4-62ba-2b2b3d61b929.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## YokogawaLinkServer 字段

[原文連結](http://api.hslcommunication.cn/html/ae5eb6e5-ad9c-dafe-cb5d-1b5e6c7d05bd.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Yokogawa](../html/e985e513-76fd-7ba4-62ba-2b2b3d61b929.htm "HslCommunication.Profinet.Yokogawa")

[YokogawaLinkServer 类](../html/ab5927f5-b6a1-d76a-7c06-10cbcca7ac37.htm "YokogawaLinkServer 类")

[YokogawaLinkServer 构造函数](../html/30b1499f-a533-626f-81ae-4a3051d664b2.htm "YokogawaLinkServer 构造函数 ")

[YokogawaLinkServer 属性](../html/c56bc036-ab9b-8156-d3d2-430f68508349.htm "YokogawaLinkServer 属性")

[YokogawaLinkServer 方法](../html/710635c8-b6db-c84f-fafd-1a8fb1d3b10f.htm "YokogawaLinkServer 方法")

[YokogawaLinkServer 事件](../html/459eec4e-5336-2966-0b53-80009ba8ea32.htm "YokogawaLinkServer 事件")

[YokogawaLinkServer 字段](../html/ae5eb6e5-ad9c-dafe-cb5d-1b5e6c7d05bd.htm "YokogawaLinkServer 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| YokogawaLinkServer 字段 |

[YokogawaLinkServer](ab5927f5-b6a1-d76a-7c06-10cbcca7ac37.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)字段

|  | 名称 | 说明 |
| --- | --- | --- |
| 受保护的字段 | [LogMsgFormatBinary](2ed2ff87-caec-0e44-f08b-e2b015225a58.htm) | 设置日志记录报文是否二进制，如果为False，那就使用ASCII码  Set whether the log message is binary, if it is False, then use ASCII code (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[YokogawaLinkServer 类](ab5927f5-b6a1-d76a-7c06-10cbcca7ac37.htm)

[HslCommunication.Profinet.Yokogawa 命名空间](e985e513-76fd-7ba4-62ba-2b2b3d61b929.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## YokogawaLinkTcp 类

[原文連結](http://api.hslcommunication.cn/html/31845cf4-eb9e-f801-f6b8-c2759f9b9a5b.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Yokogawa](../html/e985e513-76fd-7ba4-62ba-2b2b3d61b929.htm "HslCommunication.Profinet.Yokogawa")

[YokogawaLinkTcp 类](../html/31845cf4-eb9e-f801-f6b8-c2759f9b9a5b.htm "YokogawaLinkTcp 类")

[YokogawaLinkTcp 构造函数](../html/3b6774ca-ec80-e83e-59ea-4a8e82957659.htm "YokogawaLinkTcp 构造函数 ")

[YokogawaLinkTcp 属性](../html/0db863b4-4e9a-e5c0-5693-40c666a418b4.htm "YokogawaLinkTcp 属性")

[YokogawaLinkTcp 方法](../html/2a4654f3-ad6a-4087-12ab-c51f22cf3153.htm "YokogawaLinkTcp 方法")

[YokogawaLinkTcp 字段](../html/d6a91526-94b8-dd18-6134-94c8c7471549.htm "YokogawaLinkTcp 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| YokogawaLinkTcp 类 |

横河PLC的二进制通信类，支持X,Y,I,E,M,T,C,L继电器类型的数据读写，支持D,B,F,R,V,Z,W,TN,CN寄存器类型的数据读写，还支持一些高级的信息读写接口，详细参考Demo界面。  
Yokogawa PLC's binary communication type, supports X, Y, I, E, M, T, C, L relay type data read and write,
supports D, B, F, R, V, Z, W, TN, CN registers Types of data reading and writing, and some advanced information reading and writing interfaces are also supported.
Please refer to the Demo UI.

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  [HslCommunication.Core.NetBinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)  
    [HslCommunication.Core.DeviceDeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)  
      [HslCommunication.Core.DeviceDeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)  
        HslCommunication.Profinet.YokogawaYokogawaLinkTcp

**命名空间：**
 [HslCommunication.Profinet.Yokogawa](e985e513-76fd-7ba4-62ba-2b2b3d61b929.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class YokogawaLinkTcp : DeviceTcpNet
```

```
Public Class YokogawaLinkTcp
	Inherits DeviceTcpNet
```

```
public ref class YokogawaLinkTcp : public DeviceTcpNet
```

```
type YokogawaLinkTcp =  
    class
        inherit DeviceTcpNet
    end
```

YokogawaLinkTcp 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [YokogawaLinkTcp](581a37aa-2b02-82b6-e47e-71d7becfad11.htm) | 实例化一个默认的对象  Instantiate a default object |
| 公共方法 | [YokogawaLinkTcp(String, Int32)](78e41756-6adc-0621-fa2d-2904748ce0ce.htm) | 指定IP地址和端口号来实例化一个对象  Specify the IP address and port number to instantiate an object |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性代码示例 | [ByteTransform](f35118c9-3691-f752-ae6f-2c5549714faa.htm) | 当前的数据变换机制，当你需要从字节数据转换类型数据的时候需要。  The current data transformation mechanism is required when you need to convert type data from byte data. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共属性 | [CommunicationPipe](ad8e903b-d289-21b0-077e-15ede8005221.htm) | 获取或设置当前的管道信息，管道类型为[CommunicationPipe](ad8e903b-d289-21b0-077e-15ede8005221.htm)的继承类，内置了[PipeTcpNet](26b94bd5-2695-36e7-6378-c3d41846be8f.htm)管道，[PipeUdpNet](5a0b810e-187e-826b-0a13-34fd5576aa0e.htm)管道，[PipeSerialPort](9d924607-3dfc-3916-b532-c5dec48fe347.htm)管道等  Get or set the current pipeline information, the pipeline type is [CommunicationPipe](ad8e903b-d289-21b0-077e-15ede8005221.htm) inheritance class, [PipeTcpNet](26b94bd5-2695-36e7-6378-c3d41846be8f.htm) pipeline, [PipeUdpNet](5a0b810e-187e-826b-0a13-34fd5576aa0e.htm) pipeline, [PipeSerialPort](9d924607-3dfc-3916-b532-c5dec48fe347.htm) pipeline, etc (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [ConnectionId](ce1efb4b-fbc9-4683-163e-0a8e71e99d6b.htm) | 当前连接的唯一ID号，默认为长度20的guid码加随机数组成，方便列表管理，也可以自己指定  The unique ID number of the current connection. The default is a 20-digit guid code plus a random number. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性代码示例 | [ConnectTimeOut](71a651e5-737d-6e88-75a2-fb891b9af6cd.htm) | 获取或设置接收服务器反馈的时间，如果为负数，则不接收反馈   Gets or sets the time to receive server feedback, and if it is a negative number, does not receive feedback (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共属性 | [CpuNumber](55c7c07c-aaa5-f287-4ffb-ec23b733860b.htm) | 获取或设置当前的CPU Number，默认值为1  Get or set the current CPU Number, the default value is 1 |
| 公共属性代码示例 | [IpAddress](c378f4e8-3fdf-1932-7239-ed44bc5af584.htm) | 获取或是设置远程服务器的IP地址，如果是本机测试，那么需要设置为127.0.0.1   Get or set the IP address of the remote server. If it is a local test, then it needs to be set to 127.0.0.1 (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共属性 | [LocalBinding](622b418f-5969-3726-5dcc-04456c4b97a7.htm) | 获取或设置绑定的本地的IP地址和端口号信息，如果端口设置为0，代表任何可用的端口  Get or set the bound local IP address and port number information, if the port is set to 0, it means any available port (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共属性代码示例 | [LogNet](8c56f246-88ac-afb8-8c06-972fafe4f78a.htm) | 组件的日志工具，支持日志记录，只要实例化后，当前网络的基本信息，就以[DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm)等级进行输出  The component's logging tool supports logging. As long as the instantiation of the basic network information, the output will be output at [DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm) (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性代码示例 | [Port](6f1ad24a-c916-5b67-aa17-a8c7dc2acb2e.htm) | 获取或设置服务器的端口号，具体的值需要取决于对方的配置  Gets or sets the port number of the server. The specific value depends on the configuration of the other party. (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共属性代码示例 | [ReceiveTimeOut](250560a7-332b-ce1b-5eb0-e3abe4ee6012.htm) | 获取或设置接收服务器反馈的时间，如果为负数，则不接收反馈   Gets or sets the time to receive server feedback, and if it is a negative number, does not receive feedback (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [SendBeforeHex](72aa1de4-116f-abfe-4adb-814a0be0e560.htm) | 获取或设置在发送通信报文前追加发送的字节信息，HEX格式，通常用于lora组网时，需要携带 00 00 00 02 四个字节的站地址功能。  Obtain or set the byte information sent before sending communication packets, HEX format, usually used for LORA networking, you need to carry 00 00 00 02 four-byte station address function. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [SleepTime](86a5fbb9-3456-6beb-7113-5bbad5eabbc5.htm) | 获取或设置在正式接收对方返回数据前的时候，需要休息的时间，当设置为0的时候，不需要休息。  Get or set the time required to rest before officially receiving the data from the other party. When it is set to 0, no rest is required. (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共属性 | [SocketKeepAliveTime](96ab7703-8c4d-b4c2-ddd9-0dc93ebcfc8d.htm) | 获取或设置客户端的Socket的心跳时间信息，这个是Socket底层自动实现的心跳包，不基于协议层实现。默认小于0，不开启心跳检测，如果需要开启，设置 60\_000 比较合适，单位毫秒  Get or set the heartbeat time information of the Socket of the client. This is the heartbeat packet automatically implemented by the bottom layer of the Socket, not based on the protocol layer. The default value is less than 0, and heartbeat detection is not enabled. If you need to enable it, it is more appropriate to set 60\_000, in milliseconds. (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 受保护的属性 | [WordLength](f1a505cc-e85c-6c8b-b189-58a1ab5753f9.htm) | 一个字单位的数据表示的地址长度，西门子为2，三菱，欧姆龙，modbusTcp就为1，AB PLC无效  The address length represented by one word of data, Siemens is 2, Mitsubishi, Omron, modbusTcp is 1, AB PLC is invalid (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法静态成员 | [BuildReadCommand](76e0eb89-d39d-195d-ac10-257f082dfeb5.htm) | 构建读取命令的原始报文信息 |
| 公共方法静态成员 | [BuildReadRandomCommand](c2020bb2-f1ba-9503-2dde-0ce77a73a0d5.htm) | 构建随机读取的原始报文的初始命令 |
| 公共方法静态成员 | [BuildReadSpecialModule(Byte, String, UInt16)](e5d9d03f-23d1-7c29-1e32-6de041846dbc.htm) | 构建读取特殊模块的命令报文，需要传入高级地址，必须以 **Special:** 开头表示特殊模块地址，示例：Special:cpu=1;unit=0;slot=1;100  To construct a command message to read a special module, the advanced address needs to be passed in. It must start with **Special:** to indicate the address of the special module, for example: Special:cpu=1;unit=0;slot=1;100 |
| 公共方法静态成员 | [BuildReadSpecialModule(Byte, Byte, Byte, UInt16, UInt16)](e5457549-40d3-f10b-2a3e-ac75887eca17.htm) | 构建读取特殊模块的命令报文 |
| 公共方法静态成员 | [BuildStartCommand](ce8b1f0e-dd23-1135-7e4a-952a66181f7e.htm) | 构建启动PLC的命令报文 |
| 公共方法静态成员 | [BuildStopCommand](6ca89884-0633-77a5-1e14-044982625879.htm) | 构建停止PLC的命令报文 |
| 公共方法静态成员 | [BuildWriteBoolCommand](d98af2ff-880d-96ef-da36-39003722ba0a.htm) | 构建批量写入Bool数组的命令，需要指定CPU Number信息和设备地址信息 |
| 公共方法静态成员 | [BuildWriteRandomBoolCommand](bd217357-fe1e-3c93-cde0-991f851782da.htm) | 构建批量随机写入Bool数组的命令，需要指定CPU Number信息和设备地址信息 |
| 公共方法静态成员 | [BuildWriteRandomWordCommand](da1b9b15-441a-6f30-96f7-660362e984d3.htm) | 构建随机写入字的命令的报文 |
| 公共方法静态成员 | [BuildWriteSpecialModule(Byte, String, Byte)](f9997857-24bd-0dd7-d827-8c7b1a655e5a.htm) | 构建写入特殊模块的命令报文，需要传入高级地址，必须以 **Special:** 开头表示特殊模块地址，示例：Special:cpu=1;unit=0;slot=1;100  To construct a command message to write a special module, the advanced address needs to be passed in. It must start with **Special:** to indicate the address of the special module, for example: Special:cpu=1;unit=0;slot=1;100 |
| 公共方法静态成员 | [BuildWriteSpecialModule(Byte, Byte, Byte, UInt16, Byte)](4b066ee0-0e07-3bab-0504-9b184b479ef8.htm) | 构建读取特殊模块的命令报文 |
| 公共方法静态成员 | [BuildWriteWordCommand](6c7b7935-d703-220d-d31b-062d9fb88d1a.htm) | 构建字写入的命令报文信息，需要指定设备地址 |
| 公共方法静态成员 | [CheckContent](95dafa81-721c-8985-862b-fdbb2f4fbbb3.htm) | 检查当前的反馈内容，如果没有发生错误，就解析出实际的数据内容。  Check the current feedback content, if there is no error, parse out the actual data content. |
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
| 受保护的方法 | [GetNewNetMessage](1ab89c30-6d40-a3d1-e115-82aec583e5e3.htm) | 获取一个新的消息对象的方法，需要在继承类里面进行重写  The method to get a new message object needs to be overridden in the inheritance class (重写 [BinaryCommunicationGetNewNetMessage](16636a1b-a6f0-9c0f-9d9c-c94e26c67459.htm).) |
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
| 公共方法 | [ModuleReset](77db8ab8-1773-41d1-5f60-1daf1ca18e18.htm) | **[商业授权]** 重置当前的模块，当前打开的连接被强制关闭。 模块中所做的设置也将被清除。然后当前对象需要重连PLC。 **[Authorization]** When this command is executed via an Ethernet interface module or an Ethernet connection of an F3SP66, F3SP67, F3SP71 or F3SP76 sequence CPU module, the connection which is currently open is forced to close. The settings made in the modules are also cleared. Then the current object needs to reconnect to the PLC. |
| 公共方法 | [ModuleResetAsync](85cb362e-607f-5801-63d9-243b1813957d.htm) | **[商业授权]** 重置当前的模块，当前打开的连接被强制关闭。 模块中所做的设置也将被清除。然后当前对象需要重连PLC。 **[Authorization]** When this command is executed via an Ethernet interface module or an Ethernet connection of an F3SP66, F3SP67, F3SP71 or F3SP76 sequence CPU module, the connection which is currently open is forced to close. The settings made in the modules are also cleared. Then the current object needs to reconnect to the PLC. |
| 公共方法 | [PackCommandWithHeader](d0b2cd76-c306-11fc-38c7-af8ad306cce7.htm) | 对当前的命令进行打包处理，通常是携带命令头内容，标记当前的命令的长度信息，需要进行重写，否则默认不打包  The current command is packaged, usually carrying the content of the command header, marking the length of the current command, and it needs to be rewritten, otherwise it is not packaged by default (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [Read(String, UInt16)](a7f62417-078d-a9bd-2e97-dc6ab4a4ab6b.htm) | 批量读取字节数组信息，需要指定地址和长度，返回原始的字节数组  Batch read byte array information, need to specify the address and length, return the original byte array (重写 [DeviceCommunicationRead(String, UInt16)](9dd09452-c200-bc65-75ff-69aa34b366d4.htm).) |
| 公共方法代码示例 | [ReadT](95abe064-5594-5ae7-a849-98c53f88fee0.htm) | 读取支持Hsl特性的数据内容，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明，或是官网的详细文档  Read the data content of the Hsl attribute. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadAsync(String, UInt16)](192e4ce5-7afc-b1b6-95a5-ee1b7f6a9afa.htm) | 异步批量读取字节数组信息，需要指定地址和长度，返回原始的字节数组  Asynchronous batch read byte array information, need to specify the address and length, return the original byte array (重写 [DeviceCommunicationReadAsync(String, UInt16)](d8c14d37-56dc-2a31-1a19-3e2b5b65197d.htm).) |
| 公共方法代码示例 | [ReadAsyncT](0522e740-4bf0-2e31-c4c5-e41384321b37.htm) | 异步读取支持Hsl特性的数据内容，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Asynchronously read the data content of the Hsl attribute. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadBool(String)](9b2faa2e-eaa9-7ace-c49f-7a0318f1df1a.htm) | 读取单个的Boolean数据信息  Read a single Boolean data message (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadBool(String, UInt16)](284fe2c1-0f65-fffc-b844-2bddc1d063a7.htm) | 批量读取Boolean数组信息，需要指定地址和长度，返回Boolean 数组  Batch read Boolean array information, need to specify the address and length, return Boolean array (重写 [DeviceCommunicationReadBool(String, UInt16)](69d0d24f-34a7-5b41-5c1b-619edc7e6b54.htm).) |
| 公共方法 | [ReadBoolAsync(String)](cc37da1f-d673-fe67-78ee-4e465d29a2bd.htm) | 异步读取单个的Boolean数据信息  Asynchronously read a single Boolean data message (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadBoolAsync(String, UInt16)](7b155086-c613-3f11-c639-c066c66fc0cb.htm) | 异步批量读取Boolean数组信息，需要指定地址和长度，返回Boolean 数组  Asynchronously batch read Boolean array information, need to specify the address and length, return Boolean array (重写 [DeviceCommunicationReadBoolAsync(String, UInt16)](a45d383d-0095-3ea7-32f8-e75af88d60f8.htm).) |
| 公共方法代码示例 | [ReadCustomerT(String)](ca069688-d8f0-7dbd-e962-86a10dce83cf.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，返回一个新的类型的实例对象。  To read a custom data type, you need to inherit from the IDataTransfer interface and return an instance object of a new type. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadCustomerT(String, T)](ca7dc1a1-c97d-689f-9819-4dae4af67206.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，传入一个实例，对这个实例进行赋值，并返回该实例的对象。  To read a custom data type, you need to inherit from the IDataTransfer interface, pass in an instance, assign a value to this instance, and return the object of the instance. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadCustomerAsyncT(String)](f3fcda2b-255d-aadd-567e-c7b3df3636d9.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，返回一个新的类型的实例对象。  To read a custom data type, you need to inherit from the IDataTransfer interface and return an instance object of a new type. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadCustomerAsyncT(String, T)](527c9fae-4743-6ebd-7b9a-5ef3a8b7c455.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，传入一个实例，对这个实例进行赋值，并返回该实例的对象。  To read a custom data type, you need to inherit from the IDataTransfer interface, pass in an instance, assign a value to this instance, and return the object of the instance. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadDateTime](927316f6-c5f9-3717-c5f2-5f9c9fc9c57e.htm) | **[商业授权]** 读取当前PLC的时间信息，包含年月日时分秒 **[Authorization]** Read current PLC time information, including year, month, day, hour, minute, and second |
| 公共方法 | [ReadDateTimeAsync](f9e4c339-0b1a-6d1d-ae62-95d6360f0e02.htm) | **[商业授权]** 读取当前PLC的时间信息，包含年月日时分秒 **[Authorization]** Read current PLC time information, including year, month, day, hour, minute, and second |
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
| 公共方法 | [ReadProgramStatus](34d80693-1574-28b2-b2b1-378fa6cf530b.htm) | **[商业授权]** 读取当前PLC的程序状态，返回1：RUN；2：Stop；3：Debug；255：ROM writer **[Authorization]** Read the program status. return code 1:RUN; 2:Stop; 3:Debug; 255:ROM writer |
| 公共方法 | [ReadProgramStatusAsync](79cc598b-d97c-3b5c-9741-44780f2067b0.htm) | **[商业授权]** 读取当前PLC的程序状态，返回1：RUN；2：Stop；3：Debug；255：ROM writer **[Authorization]** Read the program status. return code 1:RUN; 2:Stop; 3:Debug; 255:ROM writer |
| 公共方法 | [ReadRandom](af00117f-de91-2989-a76d-9e9f15def2fd.htm) | **[商业授权]** 随机读取Byte数组信息，主需要出传入String数组地址信息，就可以返回批量Byte值 **[Authorization]** Random read Byte array information, the master needs to pass in the String array address information, and then the batch can be returned to Byte value |
| 公共方法 | [ReadRandomAsync](6e6974f4-1c96-c671-7105-24a0ff2cbb9b.htm) | **[商业授权]** 随机读取Byte数组信息，主需要出传入String数组地址信息，就可以返回批量Byte值 **[Authorization]** Random read Byte array information, the master needs to pass in the String array address information, and then the batch can be returned to Byte value |
| 公共方法 | [ReadRandomBool](37ef8788-9fc2-b34d-e547-1d465d87077a.htm) | **[商业授权]** 随机读取Boolean数组信息，主需要出传入String数组地址信息，就可以返回批量Boolean值 **[Authorization]** Random read Boolean array information, the master needs to pass in the String array address information, and then the batch can be returned to Boolean value |
| 公共方法 | [ReadRandomBoolAsync](42fcd2fd-8e53-06aa-1931-8f93d831fbf1.htm) | **[商业授权]** 随机读取Boolean数组信息，主需要出传入String数组地址信息，就可以返回批量Boolean值 **[Authorization]** Random read Boolean array information, the master needs to pass in the String array address information, and then the batch can be returned to Boolean value |
| 公共方法 | [ReadRandomInt16](57c28834-515b-fce2-ef82-592416c55596.htm) | **[商业授权]** 随机读取Int16数组信息，主需要出传入String数组地址信息，就可以返回批量Int16值 **[Authorization]** Random read Int16 array information, the master needs to pass in the Int16 array address information, and then the batch can be returned to Int16 value |
| 公共方法 | [ReadRandomInt16Async](6a3a24d2-bd4e-0d86-933a-e01055c62344.htm) | **[商业授权]** 随机读取Int16数组信息，主需要出传入String数组地址信息，就可以返回批量Int16值 **[Authorization]** Random read Int16 array information, the master needs to pass in the Int16 array address information, and then the batch can be returned to Int16 value |
| 公共方法 | [ReadRandomUInt16](2bbf66b5-91fa-b315-4d1e-2b3247e26cee.htm) | **[商业授权]** 随机读取UInt16数组信息，主需要出传入String数组地址信息，就可以返回批量UInt16值 **[Authorization]** Random read UInt16 array information, the master needs to pass in the UInt16 array address information, and then the batch can be returned to UInt16 value |
| 公共方法 | [ReadRandomUInt16Async](cf0bc1fb-f4af-67a8-53bc-16fc2b91c2d4.htm) | **[商业授权]** 随机读取UInt16数组信息，主需要出传入String数组地址信息，就可以返回批量UInt16值 **[Authorization]** Random read UInt16 array information, the master needs to pass in the UInt16 array address information, and then the batch can be returned to UInt16 value |
| 公共方法 | [ReadSpecialModule](c897675e-4faf-0083-1f60-c51536141b29.htm) | **[商业授权]** 读取特殊模块的数据信息，需要指定模块单元号，模块站号，数据地址，长度信息。 **[Authorization]** To read the data information of a special module, you need to specify the module unit number, module slot number, data address, and length information. |
| 公共方法 | [ReadSpecialModuleAsync](e7e7ed87-ec4e-8734-e94d-6c0e168a8f74.htm) | **[商业授权]** 读取特殊模块的数据信息，需要指定模块单元号，模块站号，数据地址，长度信息。 **[Authorization]** To read the data information of a special module, you need to specify the module unit number, module slot number, data address, and length information. |
| 公共方法代码示例 | [ReadString(String, UInt16)](07812d76-415a-6400-c8a2-1c21da51e7a3.htm) | 读取字符串数据，默认为最常见的ASCII编码  Read string data, default is the most common ASCII encoding (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadString(String, UInt16, Encoding)](b34c71a6-58f3-6bf4-0b93-471be4f2ed31.htm) | 使用指定的编码，读取字符串数据  Reads string data using the specified encoding (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadStringAsync(String, UInt16)](1835da81-b811-d151-0b53-907d8d87d0a8.htm) | 异步读取字符串数据，默认为最常见的ASCII编码  Asynchronously read string data, default is the most common ASCII encoding (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadStringAsync(String, UInt16, Encoding)](091dfc67-3d25-f33e-3a55-a0228863815e.htm) | 异步使用指定的编码，读取字符串数据  Asynchronously reads string data using the specified encoding (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadStructT](49853082-fbad-52e5-d756-0615cedb4b83.htm) | 读取结构体类型的数据，根据结构体自身的定义，读取原始字节数组，然后解析出实际的结构体数据，结构体需要实现[HslStructAttribute](87e427a8-0c7b-43a8-c421-c27264eda90f.htm)特性 (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [ReadStructAsyncT](36461e61-d0ed-62cf-79ce-dcc7f71a3fa2.htm) | 读取结构体类型的数据，根据结构体自身的定义，读取原始字节数组，然后解析出实际的结构体数据，结构体需要实现[HslStructAttribute](87e427a8-0c7b-43a8-c421-c27264eda90f.htm)特性 (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [ReadSystemInfo](ac7dfc05-f058-2ca7-13c9-5a1849fdbf07.htm) | **[商业授权]** 读取当前PLC的系统状态，系统的ID，CPU类型，程序大小信息 **[Authorization]** Read current PLC system status, system ID, CPU type, program size information |
| 公共方法 | [ReadSystemInfoAsync](c24a5c48-0be7-7d7a-0ec5-6d7ed984411e.htm) | **[商业授权]** 读取当前PLC的系统状态，系统的ID，CPU类型，程序大小信息 **[Authorization]** Read current PLC system status, system ID, CPU type, program size information |
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
| 公共方法 | [SetDtuPipe](f765d1d7-234a-b229-9457-97e25e024dc9.htm) | 将当前的通信对象设置DTU模式，允许传入现成的管道，并返回初始化结果，如果该设备重写了握手报文，就是返回握手结果  Set the current communication object to DTU mode, allow the existing pipe to be passed in, and return the initialization result, if the device rewrites the handshake packet, the handshake result is returned (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [SetDtuPipeAsync](9d1fd12d-b804-3828-a627-a2c0eef1d23f.htm) | 将当前的通信对象设置DTU模式，允许传入现成的管道，并返回初始化结果，如果该设备重写了握手报文，就是返回握手结果  Set the current communication object to DTU mode, allow the existing pipe to be passed in, and return the initialization result, if the device rewrites the handshake packet, the handshake result is returned (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
| 公共方法 | [SetPersistentConnection](f903dae0-4a4e-c7db-8350-50e36b257ec5.htm) | **已过时。** V11版本及之前设置长连接的方法，在V12版本以上中没有任何效果，默认长连接，删除调用即可，此处保留方法是为了部分用户保持兼容性升级。  The method of setting the long connection in V11 and before, has no effect in V12 and above. this method can be deleted. The method is retained here to maintain compatibility upgrades for some users. (继承自 [DeviceTcpNet](09da6b9b-3727-e84e-5c46-19061382088e.htm)。) |
| 公共方法 | [Start](95c77910-9fcd-2b4d-f8e6-21e63b316edc.htm) | **[商业授权]** 如果未执行程序，则开始执行程序 **[Authorization]** Starts executing a program if it is not being executed |
| 公共方法 | [StartAsync](e252e129-e3aa-49ae-b57f-597f37bca330.htm) | **[商业授权]** 如果未执行程序，则开始执行程序 **[Authorization]** Starts executing a program if it is not being executed |
| 公共方法 | [Stop](9d670a20-07b5-cfad-7ff8-99a3f97b1dae.htm) | **[商业授权]** 停止当前正在执行程序 **[Authorization]** Stops the executing program. |
| 公共方法 | [StopAsync](7e1f03af-aa51-0265-2243-6e87c71d642a.htm) | **[商业授权]** 停止当前正在执行程序 **[Authorization]** Stops the executing program. |
| 公共方法 | [ToString](ac2f962c-3384-9880-a6d9-c06891fe9a3f.htm) | (重写 [DeviceTcpNetToString](209a196b-90ea-2b73-e915-ce4b11f1263d.htm).) |
| 公共方法 | [UnpackResponseContent](208017a7-8da2-33f8-3273-d7e7e850100f.htm) | 根据对方返回的报文命令，对命令进行基本的拆包，例如各种Modbus协议拆包为统一的核心报文，还支持对报文的验证  According to the message command returned by the other party, the command is basically unpacked, for example, various Modbus protocols are unpacked into a unified core message, and the verification of the message is also supported (继承自 [BinaryCommunication](1edb8336-6f5d-2697-7b3b-911625a330f6.htm)。) |
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
| 公共方法 | [Write(String, Boolean)](fa759d55-e9a6-f4e4-f58a-17adb391b61c.htm) | 批量写入Boolean数组数据，返回是否成功  Batch write Boolean array data, return whether the write was successful (重写 [DeviceCommunicationWrite(String, Boolean)](5b046c82-3d8e-550b-fdf3-41ca5558169b.htm).) |
| 公共方法 | [Write(String, Byte)](9f740db5-926d-41aa-7ae2-c575f938a8e4.htm) | 写入原始的byte数组数据到指定的地址，返回是否写入成功  Write the original byte array data to the specified address, and return whether the write was successful (重写 [DeviceCommunicationWrite(String, Byte)](b5536dd6-ccde-6883-ac27-bff9e4b37806.htm).) |
| 公共方法代码示例 | [Write(String, String, Int32)](02fe2ca9-49a3-ede4-d1af-8c6da1f644da.htm) | 写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为ASCII  Write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, it is filled with 0 and the encoding is ASCII. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, String, Encoding)](7b602661-902d-2b70-1fdb-ba87ecd307bf.htm) | 写入字符串信息，需要指定的编码信息  Write string information, need to specify the encoding information (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [Write(String, String, Int32, Encoding)](e8cb30d9-4c5e-6c68-b9af-1b1b6c2dc714.htm) | 写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为指定的编码信息  Write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, then the operation is complemented with 0 , you should specified the encoding information (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteT(T)](9f48853f-fbeb-c37d-25a2-1c76d1dee407.htm) | 写入支持Hsl特性的数据，返回是否写入成功，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Write data that supports the Hsl attribute, and return whether the write was successful. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
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
| 公共方法 | [WriteAsync(String, Boolean)](744b96be-bd79-16f4-44b8-5a9780f6fe1c.htm) | 异步批量写入Boolean数组数据，返回是否成功  Asynchronously batch write Boolean array data, return success (重写 [DeviceCommunicationWriteAsync(String, Boolean)](fe872dae-7ef9-12c1-d883-e92ba8676416.htm).) |
| 公共方法代码示例 | [WriteAsync(String, Byte)](97e64efb-851b-602e-58b0-b465d0be521d.htm) | 异步写入原始的byte数组数据到指定的地址，返回是否写入成功  Asynchronously writes the original byte array data to the specified address, and returns whether the write was successful (重写 [DeviceCommunicationWriteAsync(String, Byte)](8f1cec30-8b86-6749-c8f3-7d39168d59df.htm).) |
| 公共方法 | [WriteAsync(String, String, Int32)](b739b983-b03c-8ecf-2713-75f19233ae2f.htm) | 异步写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为ASCII  Asynchronously write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, it is filled with 0 and the encoding is ASCII. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, String, Encoding)](876e9800-0406-6407-bdf8-8ab1cf28df38.htm) | 异步写入字符串信息，需要指定的编码信息  Asynchronously write string information, need to specify the encoding information (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WriteAsync(String, String, Int32, Encoding)](ae9fb2a3-fddb-48eb-e9ad-6cc0ebbac1f6.htm) | 异步写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为指定的编码信息  Asynchronously write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, then the operation is complemented with 0 , you should specified the encoding information (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteAsyncT(T)](25d96e77-d1b2-918c-33f9-e5e4cc57c9cd.htm) | 异步写入支持Hsl特性的数据，返回是否写入成功，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Asynchronously write data that supports the Hsl attribute, and return whether the write was successful. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteCustomerT](84a37f82-3f67-a0f8-fd3f-0630772565e4.htm) | 写入自定义类型的数据，该类型必须继承自IDataTransfer接口  Write data of a custom type, which must inherit from the IDataTransfer interface (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法代码示例 | [WriteCustomerAsyncT](04628397-cfa1-d785-7790-4c02cb377ab6.htm) | 写入自定义类型的数据，该类型必须继承自IDataTransfer接口  Write data of a custom type, which must inherit from the IDataTransfer interface (继承自 [DeviceCommunication](d345c980-d488-c83d-72ce-f958039275db.htm)。) |
| 公共方法 | [WriteRandom(String, Byte)](fada762f-8559-1ecf-8a91-7500f1553797.htm) | **[商业授权]** 随机写入Byte数组信息，主需要出传入String数组地址信息，以及对应的Byte数组值 **[Authorization]** Randomly write the Byte array information, the main need to pass in the String array address information, and the corresponding Byte array value |
| 公共方法 | [WriteRandom(String, Int16)](2b9cc33a-7e6a-0137-4d95-85df2f4e8d45.htm) | **[商业授权]** 随机写入Int16数组信息，主需要出传入String数组地址信息，以及对应的Int16数组值 **[Authorization]** Randomly write the Int16 array information, the main need to pass in the Int16 array address information, and the corresponding Int16 array value |
| 公共方法 | [WriteRandom(String, UInt16)](3a85aefc-970b-da93-a575-831792dd8ec0.htm) | **[商业授权]** 随机写入UInt16数组信息，主需要出传入String数组地址信息，以及对应的UInt16数组值 **[Authorization]** Randomly write the UInt16 array information, the main need to pass in the UInt16 array address information, and the corresponding UInt16 array value |
| 公共方法 | [WriteRandomAsync(String, Byte)](f73849a0-b8aa-959a-4ffa-2843e6a37987.htm) | **[商业授权]** 随机写入Byte数组信息，主需要出传入String数组地址信息，以及对应的Byte数组值 **[Authorization]** Randomly write the Byte array information, the main need to pass in the String array address information, and the corresponding Byte array value |
| 公共方法 | [WriteRandomAsync(String, Int16)](43fee8a6-b5b7-0cd8-3f7f-98074fb12f62.htm) | **[商业授权]** 随机写入Int16数组信息，主需要出传入String数组地址信息，以及对应的Int16数组值 **[Authorization]** Randomly write the Int16 array information, the main need to pass in the Int16 array address information, and the corresponding Int16 array value |
| 公共方法 | [WriteRandomAsync(String, UInt16)](4e5427ef-9ea9-f642-f1d9-6c4469cf2464.htm) | **[商业授权]** 随机写入UInt16数组信息，主需要出传入String数组地址信息，以及对应的UInt16数组值 **[Authorization]** Randomly write the UInt16 array information, the main need to pass in the UInt16 array address information, and the corresponding UInt16 array value |
| 公共方法 | [WriteRandomBool](46f3ddfc-942a-e974-d416-9be4c4aff586.htm) | **[商业授权]** 随机写入Boolean数组信息，主需要出传入String数组地址信息，以及对应的Boolean数组值 **[Authorization]** Randomly write the Boolean array information, the main need to pass in the String array address information, and the corresponding Boolean array value |
| 公共方法 | [WriteRandomBoolAsync](79a52a28-496c-b49e-c8ac-3e1ac896c08a.htm) | **[商业授权]** 随机写入Boolean数组信息，主需要出传入String数组地址信息，以及对应的Boolean数组值 **[Authorization]** Randomly write the Boolean array information, the main need to pass in the String array address information, and the corresponding Boolean array value |

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

基础的数据读写面向VIP用户开放，高级的读写随机数据，启动停止命令，读取程序状态，
系统信息，PLC时间，读写特殊的模块数据需要商业用户授权，读取的数据长度，读取的随机地址长度，在商业授权下，长度不受限制，可以无限大。

![](../icons/SectionExpanded.png)示例

例如我们正常读取一个D100的数据如下：

Read示例

[复制](# "复制")

```
// 实例化对象，指定PLC的ip地址
YokogawaLinkTcp yokogawa = new YokogawaLinkTcp( " 192.168.0.2", 12289 );
OperateResult connect = yokogawa.ConnectServer( );
if (connect.IsSuccess)
{
    Console.WriteLine( "Connect plc success!" );
}
else
{
    Console.WriteLine( "Connect plc failed, reason:" + connect.Message );
    return;
}

// 举例读取D100的值
OperateResult<short> read = yokogawa.ReadInt16( "D100" );
if (read.IsSuccess)
{
    Console.WriteLine( "Read D100: " + read.Content );
}
else
{
    Console.WriteLine( "Read failed, reason: " + read.Message );
    return;
}

// 写入的原理是一样的
OperateResult write = yokogawa.Write( "D100", (short)100 );
if (write.IsSuccess)
{
    Console.WriteLine( "write plc success!" );
}
else
{
    Console.WriteLine( "write failed:" + write.Message );
    return;
}

// 关闭连接
yokogawa.ConnectClose( );
```

我们在读取的时候可以动态的变更cpu信息，参考下面的代码

Read示例

[复制](# "复制")

```
// 在读取的时候，可以携带CPU信息，例如
OperateResult<short> read = yokogawa.ReadInt16( "cpu=2;D100" );
if (read.IsSuccess)
{
    Console.WriteLine( "Read D100: " + read.Content );
}
else
{
    Console.WriteLine( "Read failed, reason: " + read.Message );
    return;
}
```

关于随机读写的代码示例，可以读写地址分布很散的地址，参考下面的代码

Read示例

[复制](# "复制")

```
// 使用随机读取的方式读取不同地址的数据信息
OperateResult<bool[]> read = yokogawa.ReadRandomBool( new string[] { "X0", "Y100", "M200", "M1000" } );
if (read.IsSuccess)
{
    bool x_0 = read.Content[0];
    bool y_100 = read.Content[1];
    bool m_200 = read.Content[2];
    bool m_1000 = read.Content[3];
}
else
{
    Console.WriteLine( "Read failed, reason: " + read.Message );
}


// 写入也是类似的
OperateResult write = yokogawa.WriteRandomBool( new string[] { "X0", "Y100", "M200", "M1000" },
    new bool[] { true, false, false, true } );
if (write.IsSuccess)
{
    Console.WriteLine( "write plc success!" );
}
else
{
    Console.WriteLine( "write failed:" + write.Message );
    return;
}
```

最后看一下读取特殊模块的数据，可以读取基本的字节数据，也可以使用富文本的地址读取

Read示例

[复制](# "复制")

```
// 读取PLC的特殊模块的地址数据
OperateResult<byte[]> read = yokogawa.ReadSpecialModule( 1, 2, 100, 2 );
if (read.IsSuccess)
{
    // 我们读取到了原始的数据内容
    byte[] content = read.Content;
}
else
{
    Console.WriteLine( "Read failed, reason: " + read.Message );
}


// 我们想要读取特殊模块的short数据，上面的读取方式我们还要转换，比较麻烦，所以用下面的方式，必须 Special:开头
OperateResult<short> readInt16 = yokogawa.ReadInt16( "Special:unit=1;slot=2;100" );
if (readInt16.IsSuccess)
{
    Console.WriteLine( "Read D100: " + readInt16.Content );
}
else
{
    Console.WriteLine( "Read failed, reason: " + readInt16.Message );
    return;
}

// 写入都是类似的
```

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.Profinet.Yokogawa 命名空间](e985e513-76fd-7ba4-62ba-2b2b3d61b929.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## YokogawaLinkTcp 构造函数 

[原文連結](http://api.hslcommunication.cn/html/3b6774ca-ec80-e83e-59ea-4a8e82957659.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Yokogawa](../html/e985e513-76fd-7ba4-62ba-2b2b3d61b929.htm "HslCommunication.Profinet.Yokogawa")

[YokogawaLinkTcp 类](../html/31845cf4-eb9e-f801-f6b8-c2759f9b9a5b.htm "YokogawaLinkTcp 类")

[YokogawaLinkTcp 构造函数](../html/3b6774ca-ec80-e83e-59ea-4a8e82957659.htm "YokogawaLinkTcp 构造函数 ")

[YokogawaLinkTcp 构造函数](../html/581a37aa-2b02-82b6-e47e-71d7becfad11.htm "YokogawaLinkTcp 构造函数 ")

[YokogawaLinkTcp 构造函数 (String, Int32)](../html/78e41756-6adc-0621-fa2d-2904748ce0ce.htm "YokogawaLinkTcp 构造函数 (String, Int32)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| YokogawaLinkTcp 构造函数 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [YokogawaLinkTcp](581a37aa-2b02-82b6-e47e-71d7becfad11.htm) | 实例化一个默认的对象  Instantiate a default object |
| 公共方法 | [YokogawaLinkTcp(String, Int32)](78e41756-6adc-0621-fa2d-2904748ce0ce.htm) | 指定IP地址和端口号来实例化一个对象  Specify the IP address and port number to instantiate an object |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[YokogawaLinkTcp 类](31845cf4-eb9e-f801-f6b8-c2759f9b9a5b.htm)

[HslCommunication.Profinet.Yokogawa 命名空间](e985e513-76fd-7ba4-62ba-2b2b3d61b929.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)