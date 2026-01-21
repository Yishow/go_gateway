# HslCommunication - HslCommunication.Profinet.Toledo

> 分類頁數: 30



---
## HslCommunication.Profinet.Toledo

[原文連結](http://api.hslcommunication.cn/html/3455a9c0-e5c9-7245-1dce-b1fee92e6846.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Toledo](../html/3455a9c0-e5c9-7245-1dce-b1fee92e6846.htm "HslCommunication.Profinet.Toledo")

[ToledoSerial 类](../html/61773df3-22c2-08bf-1635-9cc4d0e62430.htm "ToledoSerial 类")

[ToledoSerial.ToledoStandardDataReceivedDelegate 委托](../html/4d9acb69-0a01-8cba-3658-6977be61668b.htm "ToledoSerial.ToledoStandardDataReceivedDelegate 委托")

[ToledoStandardData 类](../html/13fdb4fb-3312-d193-7e77-a383ddffbf1d.htm "ToledoStandardData 类")

[ToledoTcpServer 类](../html/32226f81-8b8f-228c-4a5f-55dedb6ab8e2.htm "ToledoTcpServer 类")

[ToledoTcpServer.ToledoStandardDataReceivedDelegate 委托](../html/49d9d452-7681-6598-36ad-c6c03053782b.htm "ToledoTcpServer.ToledoStandardDataReceivedDelegate 委托")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslCommunication.Profinet.Toledo 命名空间 |

[缺少 "N:HslCommunication.Profinet.Toledo" 的 <summary> 文档]

![](../icons/SectionExpanded.png)类

|  | 类 | 说明 |
| --- | --- | --- |
| 公共类 | [ToledoSerial](61773df3-22c2-08bf-1635-9cc4d0e62430.htm) | 托利多电子秤的串口服务器对象 |
| 公共类 | [ToledoStandardData](13fdb4fb-3312-d193-7e77-a383ddffbf1d.htm) | 托利多标准格式的数据类对象 |
| 公共类 | [ToledoTcpServer](32226f81-8b8f-228c-4a5f-55dedb6ab8e2.htm) | 托利多电子秤的TCP服务器，启动服务器后，等待电子秤的数据连接。 |

![](../icons/SectionExpanded.png)委托

|  | 委托 | 说明 |
| --- | --- | --- |
| 公共委托 | [ToledoSerialToledoStandardDataReceivedDelegate](4d9acb69-0a01-8cba-3658-6977be61668b.htm) | 托利多数据接收时的委托 |
| 公共委托 | [ToledoTcpServerToledoStandardDataReceivedDelegate](49d9d452-7681-6598-36ad-c6c03053782b.htm) | 托利多数据接收时的委托 |

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ToledoSerial 类

[原文連結](http://api.hslcommunication.cn/html/61773df3-22c2-08bf-1635-9cc4d0e62430.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Toledo](../html/3455a9c0-e5c9-7245-1dce-b1fee92e6846.htm "HslCommunication.Profinet.Toledo")

[ToledoSerial 类](../html/61773df3-22c2-08bf-1635-9cc4d0e62430.htm "ToledoSerial 类")

[ToledoSerial 构造函数](../html/a90f33c0-1d80-82a1-fffe-1e9ef69bd52d.htm "ToledoSerial 构造函数 ")

[ToledoSerial 属性](../html/cb84b93a-70b6-403b-0829-c73f29baef54.htm "ToledoSerial 属性")

[ToledoSerial 方法](../html/244168d4-6935-958b-3a55-e0c8d7b8ff21.htm "ToledoSerial 方法")

[ToledoSerial 事件](../html/5833568b-8df7-5029-ce9c-0daaf2616061.htm "ToledoSerial 事件")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ToledoSerial 类 |

托利多电子秤的串口服务器对象

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  HslCommunication.Profinet.ToledoToledoSerial

**命名空间：**
 [HslCommunication.Profinet.Toledo](3455a9c0-e5c9-7245-1dce-b1fee92e6846.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class ToledoSerial
```

```
Public Class ToledoSerial
```

```
public ref class ToledoSerial
```

```
type ToledoSerial =  class end
```

ToledoSerial 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [ToledoSerial](a90f33c0-1d80-82a1-fffe-1e9ef69bd52d.htm) | 实例化一个默认的对象 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [BaudRate](fa1faa63-acbb-ba2f-ddbe-4f4af4f40377.htm) | 当前连接串口信息的波特率  Baud rate of current connection serial port information |
| 公共属性 | [HasChk](819adc97-8de5-f822-8775-222c7edabe5d.htm) | 获取或设置当前的报文否是含有校验的，默认为含有校验 |
| 公共属性代码示例 | [LogNet](d70c17a2-b265-02ca-45f4-34de2209c31c.htm) | 组件的日志工具，支持日志记录，只要实例化后，当前网络的基本信息，就以[DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm)等级进行输出  The component's logging tool supports logging. As long as the instantiation of the basic network information, the output will be output at [DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm) |
| 公共属性 | [PortName](a465ce32-cb75-a950-fd76-63312a9966e3.htm) | 当前连接串口信息的端口号名称  The port name of the current connection serial port information |
| 公共属性 | [ReceiveTimeout](d0348a07-8906-a6fb-bad0-4515ec117f43.htm) | 接收数据的超时时间，默认5000ms  Timeout for receiving data, default is 5000ms |
| 公共属性 | [RtsEnable](b602ea9e-46ee-216e-0c8e-afd4cdc39d08.htm) | 获取或设置一个值，该值指示在串行通信中是否启用请求发送 (RTS) 信号。  Gets or sets a value indicating whether the request sending (RTS) signal is enabled in serial communication. |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [Close](bfcd36b1-f1a0-9377-912b-ed62a44aca0c.htm) | 关闭当前的串口连接  Close the current serial connection |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法 | GetType | (继承自 Object。) |
| 公共方法 | [IsOpen](2e3d66ef-9397-a8b7-1ef9-67a699af15c0.htm) | 获取一个值，指示串口是否处于打开状态  Gets a value indicating whether the serial port is open |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法 | [Open](b5a78c2c-b159-4634-c9ec-ee64d885ff70.htm) | 打开一个新的串行端口连接  Open a new serial port connection |
| 公共方法 | [SerialPortInni(ActionSerialPort)](82011d60-4795-6e72-51b1-be0bd201d530.htm) | 根据自定义初始化方法进行初始化串口信息  Initialize the serial port information according to the custom initialization method |
| 公共方法 | [SerialPortInni(String)](b09ded66-4339-b33e-8828-4281fbe0b28a.htm) | 初始化串口信息，9600波特率，8位数据位，1位停止位，无奇偶校验  Initial serial port information, 9600 baud rate, 8 data bits, 1 stop bit, no parity |
| 公共方法 | [SerialPortInni(String, Int32)](84aac95f-c526-96bb-95e1-959921d057cd.htm) | 初始化串口信息，波特率，8位数据位，1位停止位，无奇偶校验  Initializes serial port information, baud rate, 8-bit data bit, 1-bit stop bit, no parity |
| 公共方法 | [SerialPortInni(String, Int32, Int32, StopBits, Parity)](26a464c9-69d2-2e04-6217-a7fc8f88035d.htm) | 初始化串口信息，波特率，数据位，停止位，奇偶校验需要全部自己来指定  Start serial port information, baud rate, data bit, stop bit, parity all need to be specified |
| 公共方法 | [ToString](008685e3-0f3c-5219-7ad8-b828a1da9932.htm) | (重写 ObjectToString.) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)事件

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共事件 | [OnToledoStandardDataReceived](de23564d-f97b-79ba-9236-36651a923868.htm) | 当接收到一条新的托利多的数据的时候触发 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.Profinet.Toledo 命名空间](3455a9c0-e5c9-7245-1dce-b1fee92e6846.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ToledoSerial 构造函数 

[原文連結](http://api.hslcommunication.cn/html/a90f33c0-1d80-82a1-fffe-1e9ef69bd52d.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Toledo](../html/3455a9c0-e5c9-7245-1dce-b1fee92e6846.htm "HslCommunication.Profinet.Toledo")

[ToledoSerial 类](../html/61773df3-22c2-08bf-1635-9cc4d0e62430.htm "ToledoSerial 类")

[ToledoSerial 构造函数](../html/a90f33c0-1d80-82a1-fffe-1e9ef69bd52d.htm "ToledoSerial 构造函数 ")

[ToledoSerial 属性](../html/cb84b93a-70b6-403b-0829-c73f29baef54.htm "ToledoSerial 属性")

[ToledoSerial 方法](../html/244168d4-6935-958b-3a55-e0c8d7b8ff21.htm "ToledoSerial 方法")

[ToledoSerial 事件](../html/5833568b-8df7-5029-ce9c-0daaf2616061.htm "ToledoSerial 事件")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ToledoSerial 构造函数 |

实例化一个默认的对象

**命名空间：**
 [HslCommunication.Profinet.Toledo](3455a9c0-e5c9-7245-1dce-b1fee92e6846.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public ToledoSerial()
```

```
Public Sub New
```

```
public:
ToledoSerial()
```

```
new : unit -> ToledoSerial
```

![](../icons/SectionExpanded.png)参见

#### 引用

[ToledoSerial 类](61773df3-22c2-08bf-1635-9cc4d0e62430.htm)

[HslCommunication.Profinet.Toledo 命名空间](3455a9c0-e5c9-7245-1dce-b1fee92e6846.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ToledoSerial 属性

[原文連結](http://api.hslcommunication.cn/html/cb84b93a-70b6-403b-0829-c73f29baef54.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Toledo](../html/3455a9c0-e5c9-7245-1dce-b1fee92e6846.htm "HslCommunication.Profinet.Toledo")

[ToledoSerial 类](../html/61773df3-22c2-08bf-1635-9cc4d0e62430.htm "ToledoSerial 类")

[ToledoSerial 属性](../html/cb84b93a-70b6-403b-0829-c73f29baef54.htm "ToledoSerial 属性")

[BaudRate 属性](../html/fa1faa63-acbb-ba2f-ddbe-4f4af4f40377.htm "BaudRate 属性 ")

[HasChk 属性](../html/819adc97-8de5-f822-8775-222c7edabe5d.htm "HasChk 属性 ")

[LogNet 属性](../html/d70c17a2-b265-02ca-45f4-34de2209c31c.htm "LogNet 属性 ")

[PortName 属性](../html/a465ce32-cb75-a950-fd76-63312a9966e3.htm "PortName 属性 ")

[ReceiveTimeout 属性](../html/d0348a07-8906-a6fb-bad0-4515ec117f43.htm "ReceiveTimeout 属性 ")

[RtsEnable 属性](../html/b602ea9e-46ee-216e-0c8e-afd4cdc39d08.htm "RtsEnable 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ToledoSerial 属性 |

[ToledoSerial](61773df3-22c2-08bf-1635-9cc4d0e62430.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [BaudRate](fa1faa63-acbb-ba2f-ddbe-4f4af4f40377.htm) | 当前连接串口信息的波特率  Baud rate of current connection serial port information |
| 公共属性 | [HasChk](819adc97-8de5-f822-8775-222c7edabe5d.htm) | 获取或设置当前的报文否是含有校验的，默认为含有校验 |
| 公共属性代码示例 | [LogNet](d70c17a2-b265-02ca-45f4-34de2209c31c.htm) | 组件的日志工具，支持日志记录，只要实例化后，当前网络的基本信息，就以[DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm)等级进行输出  The component's logging tool supports logging. As long as the instantiation of the basic network information, the output will be output at [DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm) |
| 公共属性 | [PortName](a465ce32-cb75-a950-fd76-63312a9966e3.htm) | 当前连接串口信息的端口号名称  The port name of the current connection serial port information |
| 公共属性 | [ReceiveTimeout](d0348a07-8906-a6fb-bad0-4515ec117f43.htm) | 接收数据的超时时间，默认5000ms  Timeout for receiving data, default is 5000ms |
| 公共属性 | [RtsEnable](b602ea9e-46ee-216e-0c8e-afd4cdc39d08.htm) | 获取或设置一个值，该值指示在串行通信中是否启用请求发送 (RTS) 信号。  Gets or sets a value indicating whether the request sending (RTS) signal is enabled in serial communication. |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[ToledoSerial 类](61773df3-22c2-08bf-1635-9cc4d0e62430.htm)

[HslCommunication.Profinet.Toledo 命名空间](3455a9c0-e5c9-7245-1dce-b1fee92e6846.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## BaudRate 属性 

[原文連結](http://api.hslcommunication.cn/html/fa1faa63-acbb-ba2f-ddbe-4f4af4f40377.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Toledo](../html/3455a9c0-e5c9-7245-1dce-b1fee92e6846.htm "HslCommunication.Profinet.Toledo")

[ToledoSerial 类](../html/61773df3-22c2-08bf-1635-9cc4d0e62430.htm "ToledoSerial 类")

[ToledoSerial 属性](../html/cb84b93a-70b6-403b-0829-c73f29baef54.htm "ToledoSerial 属性")

[BaudRate 属性](../html/fa1faa63-acbb-ba2f-ddbe-4f4af4f40377.htm "BaudRate 属性 ")

[HasChk 属性](../html/819adc97-8de5-f822-8775-222c7edabe5d.htm "HasChk 属性 ")

[LogNet 属性](../html/d70c17a2-b265-02ca-45f4-34de2209c31c.htm "LogNet 属性 ")

[PortName 属性](../html/a465ce32-cb75-a950-fd76-63312a9966e3.htm "PortName 属性 ")

[ReceiveTimeout 属性](../html/d0348a07-8906-a6fb-bad0-4515ec117f43.htm "ReceiveTimeout 属性 ")

[RtsEnable 属性](../html/b602ea9e-46ee-216e-0c8e-afd4cdc39d08.htm "RtsEnable 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ToledoSerialBaudRate 属性 |

当前连接串口信息的波特率  
Baud rate of current connection serial port information

**命名空间：**
 [HslCommunication.Profinet.Toledo](3455a9c0-e5c9-7245-1dce-b1fee92e6846.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public int BaudRate { get; }
```

```
Public ReadOnly Property BaudRate As Integer
	Get
```

```
public:
property int BaudRate {
	int get ();
}
```

```
member BaudRate : int with get
```

#### 属性值

类型：Int32

![](../icons/SectionExpanded.png)参见

#### 引用

[ToledoSerial 类](61773df3-22c2-08bf-1635-9cc4d0e62430.htm)

[HslCommunication.Profinet.Toledo 命名空间](3455a9c0-e5c9-7245-1dce-b1fee92e6846.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## HasChk 属性 

[原文連結](http://api.hslcommunication.cn/html/819adc97-8de5-f822-8775-222c7edabe5d.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Toledo](../html/3455a9c0-e5c9-7245-1dce-b1fee92e6846.htm "HslCommunication.Profinet.Toledo")

[ToledoSerial 类](../html/61773df3-22c2-08bf-1635-9cc4d0e62430.htm "ToledoSerial 类")

[ToledoSerial 属性](../html/cb84b93a-70b6-403b-0829-c73f29baef54.htm "ToledoSerial 属性")

[BaudRate 属性](../html/fa1faa63-acbb-ba2f-ddbe-4f4af4f40377.htm "BaudRate 属性 ")

[HasChk 属性](../html/819adc97-8de5-f822-8775-222c7edabe5d.htm "HasChk 属性 ")

[LogNet 属性](../html/d70c17a2-b265-02ca-45f4-34de2209c31c.htm "LogNet 属性 ")

[PortName 属性](../html/a465ce32-cb75-a950-fd76-63312a9966e3.htm "PortName 属性 ")

[ReceiveTimeout 属性](../html/d0348a07-8906-a6fb-bad0-4515ec117f43.htm "ReceiveTimeout 属性 ")

[RtsEnable 属性](../html/b602ea9e-46ee-216e-0c8e-afd4cdc39d08.htm "RtsEnable 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ToledoSerialHasChk 属性 |

获取或设置当前的报文否是含有校验的，默认为含有校验

**命名空间：**
 [HslCommunication.Profinet.Toledo](3455a9c0-e5c9-7245-1dce-b1fee92e6846.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public bool HasChk { get; set; }
```

```
Public Property HasChk As Boolean
	Get
	Set
```

```
public:
property bool HasChk {
	bool get ();
	void set (bool value);
}
```

```
member HasChk : bool with get, set
```

#### 属性值

类型：Boolean

![](../icons/SectionExpanded.png)参见

#### 引用

[ToledoSerial 类](61773df3-22c2-08bf-1635-9cc4d0e62430.htm)

[HslCommunication.Profinet.Toledo 命名空间](3455a9c0-e5c9-7245-1dce-b1fee92e6846.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## LogNet 属性 

[原文連結](http://api.hslcommunication.cn/html/d70c17a2-b265-02ca-45f4-34de2209c31c.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Toledo](../html/3455a9c0-e5c9-7245-1dce-b1fee92e6846.htm "HslCommunication.Profinet.Toledo")

[ToledoSerial 类](../html/61773df3-22c2-08bf-1635-9cc4d0e62430.htm "ToledoSerial 类")

[ToledoSerial 属性](../html/cb84b93a-70b6-403b-0829-c73f29baef54.htm "ToledoSerial 属性")

[BaudRate 属性](../html/fa1faa63-acbb-ba2f-ddbe-4f4af4f40377.htm "BaudRate 属性 ")

[HasChk 属性](../html/819adc97-8de5-f822-8775-222c7edabe5d.htm "HasChk 属性 ")

[LogNet 属性](../html/d70c17a2-b265-02ca-45f4-34de2209c31c.htm "LogNet 属性 ")

[PortName 属性](../html/a465ce32-cb75-a950-fd76-63312a9966e3.htm "PortName 属性 ")

[ReceiveTimeout 属性](../html/d0348a07-8906-a6fb-bad0-4515ec117f43.htm "ReceiveTimeout 属性 ")

[RtsEnable 属性](../html/b602ea9e-46ee-216e-0c8e-afd4cdc39d08.htm "RtsEnable 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ToledoSerialLogNet 属性 |

组件的日志工具，支持日志记录，只要实例化后，当前网络的基本信息，就以[DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm)等级进行输出  
The component's logging tool supports logging. As long as the instantiation of the basic network information, the output will be output at [DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm)

**命名空间：**
 [HslCommunication.Profinet.Toledo](3455a9c0-e5c9-7245-1dce-b1fee92e6846.htm)  
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

![](../icons/SectionExpanded.png)备注

只要实例化即可以记录日志，实例化的对象需要实现接口 [ILogNet](d95d4704-db0d-b379-41e6-03879927a543.htm) ，本组件提供了三个日志记录类，你可以实现基于 [ILogNet](d95d4704-db0d-b379-41e6-03879927a543.htm) 的对象。

![](../icons/SectionExpanded.png)示例

如下的实例化适用于所有的Network及其派生类，以下举两个例子，三菱的设备类及服务器类

LogNet示例

[复制](# "复制")

```
// 设备连接对象的日志
MelsecMcNet melsec = new MelsecMcNet( "192.168.0.100", 6000 );

// 举例实现日志文件为单日志文件
melsec.LogNet = new HslCommunication.LogNet.LogNetSingle( "D://123.txt" );
```

LogNet示例

[复制](# "复制")

```
// 一般服务器对象的
NetSimplifyServer simplifyServer = new NetSimplifyServer( );
simplifyServer.LogNet = new HslCommunication.LogNet.LogNetSingle( "D://log.txt" );
simplifyServer.ReceiveStringEvent += ( HslCommunication.Core.Net.AppSession session, HslCommunication.NetHandle handle, string data ) =>
{
    simplifyServer.SendMessage( session, handle, "Back:" + data );
};
simplifyServer.ServerStart( 45678 );
```

![](../icons/SectionExpanded.png)参见

#### 引用

[ToledoSerial 类](61773df3-22c2-08bf-1635-9cc4d0e62430.htm)

[HslCommunication.Profinet.Toledo 命名空间](3455a9c0-e5c9-7245-1dce-b1fee92e6846.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## PortName 属性 

[原文連結](http://api.hslcommunication.cn/html/a465ce32-cb75-a950-fd76-63312a9966e3.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Toledo](../html/3455a9c0-e5c9-7245-1dce-b1fee92e6846.htm "HslCommunication.Profinet.Toledo")

[ToledoSerial 类](../html/61773df3-22c2-08bf-1635-9cc4d0e62430.htm "ToledoSerial 类")

[ToledoSerial 属性](../html/cb84b93a-70b6-403b-0829-c73f29baef54.htm "ToledoSerial 属性")

[BaudRate 属性](../html/fa1faa63-acbb-ba2f-ddbe-4f4af4f40377.htm "BaudRate 属性 ")

[HasChk 属性](../html/819adc97-8de5-f822-8775-222c7edabe5d.htm "HasChk 属性 ")

[LogNet 属性](../html/d70c17a2-b265-02ca-45f4-34de2209c31c.htm "LogNet 属性 ")

[PortName 属性](../html/a465ce32-cb75-a950-fd76-63312a9966e3.htm "PortName 属性 ")

[ReceiveTimeout 属性](../html/d0348a07-8906-a6fb-bad0-4515ec117f43.htm "ReceiveTimeout 属性 ")

[RtsEnable 属性](../html/b602ea9e-46ee-216e-0c8e-afd4cdc39d08.htm "RtsEnable 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ToledoSerialPortName 属性 |

当前连接串口信息的端口号名称  
The port name of the current connection serial port information

**命名空间：**
 [HslCommunication.Profinet.Toledo](3455a9c0-e5c9-7245-1dce-b1fee92e6846.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public string PortName { get; }
```

```
Public ReadOnly Property PortName As String
	Get
```

```
public:
property String^ PortName {
	String^ get ();
}
```

```
member PortName : string with get
```

#### 属性值

类型：String

![](../icons/SectionExpanded.png)参见

#### 引用

[ToledoSerial 类](61773df3-22c2-08bf-1635-9cc4d0e62430.htm)

[HslCommunication.Profinet.Toledo 命名空间](3455a9c0-e5c9-7245-1dce-b1fee92e6846.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReceiveTimeout 属性 

[原文連結](http://api.hslcommunication.cn/html/d0348a07-8906-a6fb-bad0-4515ec117f43.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Toledo](../html/3455a9c0-e5c9-7245-1dce-b1fee92e6846.htm "HslCommunication.Profinet.Toledo")

[ToledoSerial 类](../html/61773df3-22c2-08bf-1635-9cc4d0e62430.htm "ToledoSerial 类")

[ToledoSerial 属性](../html/cb84b93a-70b6-403b-0829-c73f29baef54.htm "ToledoSerial 属性")

[BaudRate 属性](../html/fa1faa63-acbb-ba2f-ddbe-4f4af4f40377.htm "BaudRate 属性 ")

[HasChk 属性](../html/819adc97-8de5-f822-8775-222c7edabe5d.htm "HasChk 属性 ")

[LogNet 属性](../html/d70c17a2-b265-02ca-45f4-34de2209c31c.htm "LogNet 属性 ")

[PortName 属性](../html/a465ce32-cb75-a950-fd76-63312a9966e3.htm "PortName 属性 ")

[ReceiveTimeout 属性](../html/d0348a07-8906-a6fb-bad0-4515ec117f43.htm "ReceiveTimeout 属性 ")

[RtsEnable 属性](../html/b602ea9e-46ee-216e-0c8e-afd4cdc39d08.htm "RtsEnable 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ToledoSerialReceiveTimeout 属性 |

接收数据的超时时间，默认5000ms  
Timeout for receiving data, default is 5000ms

**命名空间：**
 [HslCommunication.Profinet.Toledo](3455a9c0-e5c9-7245-1dce-b1fee92e6846.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public int ReceiveTimeout { get; set; }
```

```
Public Property ReceiveTimeout As Integer
	Get
	Set
```

```
public:
property int ReceiveTimeout {
	int get ();
	void set (int value);
}
```

```
member ReceiveTimeout : int with get, set
```

#### 属性值

类型：Int32

![](../icons/SectionExpanded.png)参见

#### 引用

[ToledoSerial 类](61773df3-22c2-08bf-1635-9cc4d0e62430.htm)

[HslCommunication.Profinet.Toledo 命名空间](3455a9c0-e5c9-7245-1dce-b1fee92e6846.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## RtsEnable 属性 

[原文連結](http://api.hslcommunication.cn/html/b602ea9e-46ee-216e-0c8e-afd4cdc39d08.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Toledo](../html/3455a9c0-e5c9-7245-1dce-b1fee92e6846.htm "HslCommunication.Profinet.Toledo")

[ToledoSerial 类](../html/61773df3-22c2-08bf-1635-9cc4d0e62430.htm "ToledoSerial 类")

[ToledoSerial 属性](../html/cb84b93a-70b6-403b-0829-c73f29baef54.htm "ToledoSerial 属性")

[BaudRate 属性](../html/fa1faa63-acbb-ba2f-ddbe-4f4af4f40377.htm "BaudRate 属性 ")

[HasChk 属性](../html/819adc97-8de5-f822-8775-222c7edabe5d.htm "HasChk 属性 ")

[LogNet 属性](../html/d70c17a2-b265-02ca-45f4-34de2209c31c.htm "LogNet 属性 ")

[PortName 属性](../html/a465ce32-cb75-a950-fd76-63312a9966e3.htm "PortName 属性 ")

[ReceiveTimeout 属性](../html/d0348a07-8906-a6fb-bad0-4515ec117f43.htm "ReceiveTimeout 属性 ")

[RtsEnable 属性](../html/b602ea9e-46ee-216e-0c8e-afd4cdc39d08.htm "RtsEnable 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ToledoSerialRtsEnable 属性 |

获取或设置一个值，该值指示在串行通信中是否启用请求发送 (RTS) 信号。  
Gets or sets a value indicating whether the request sending (RTS) signal is enabled in serial communication.

**命名空间：**
 [HslCommunication.Profinet.Toledo](3455a9c0-e5c9-7245-1dce-b1fee92e6846.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public bool RtsEnable { get; set; }
```

```
Public Property RtsEnable As Boolean
	Get
	Set
```

```
public:
property bool RtsEnable {
	bool get ();
	void set (bool value);
}
```

```
member RtsEnable : bool with get, set
```

#### 属性值

类型：Boolean

![](../icons/SectionExpanded.png)参见

#### 引用

[ToledoSerial 类](61773df3-22c2-08bf-1635-9cc4d0e62430.htm)

[HslCommunication.Profinet.Toledo 命名空间](3455a9c0-e5c9-7245-1dce-b1fee92e6846.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ToledoSerial 方法

[原文連結](http://api.hslcommunication.cn/html/244168d4-6935-958b-3a55-e0c8d7b8ff21.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Toledo](../html/3455a9c0-e5c9-7245-1dce-b1fee92e6846.htm "HslCommunication.Profinet.Toledo")

[ToledoSerial 类](../html/61773df3-22c2-08bf-1635-9cc4d0e62430.htm "ToledoSerial 类")

[ToledoSerial 方法](../html/244168d4-6935-958b-3a55-e0c8d7b8ff21.htm "ToledoSerial 方法")

[Close 方法](../html/bfcd36b1-f1a0-9377-912b-ed62a44aca0c.htm "Close 方法 ")

[IsOpen 方法](../html/2e3d66ef-9397-a8b7-1ef9-67a699af15c0.htm "IsOpen 方法 ")

[Open 方法](../html/b5a78c2c-b159-4634-c9ec-ee64d885ff70.htm "Open 方法 ")

[SerialPortInni 方法](../html/df188f59-2279-df42-047a-89b95dce1c46.htm "SerialPortInni 方法 ")

[ToString 方法](../html/008685e3-0f3c-5219-7ad8-b828a1da9932.htm "ToString 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ToledoSerial 方法 |

[ToledoSerial](61773df3-22c2-08bf-1635-9cc4d0e62430.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [Close](bfcd36b1-f1a0-9377-912b-ed62a44aca0c.htm) | 关闭当前的串口连接  Close the current serial connection |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法 | GetType | (继承自 Object。) |
| 公共方法 | [IsOpen](2e3d66ef-9397-a8b7-1ef9-67a699af15c0.htm) | 获取一个值，指示串口是否处于打开状态  Gets a value indicating whether the serial port is open |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法 | [Open](b5a78c2c-b159-4634-c9ec-ee64d885ff70.htm) | 打开一个新的串行端口连接  Open a new serial port connection |
| 公共方法 | [SerialPortInni(ActionSerialPort)](82011d60-4795-6e72-51b1-be0bd201d530.htm) | 根据自定义初始化方法进行初始化串口信息  Initialize the serial port information according to the custom initialization method |
| 公共方法 | [SerialPortInni(String)](b09ded66-4339-b33e-8828-4281fbe0b28a.htm) | 初始化串口信息，9600波特率，8位数据位，1位停止位，无奇偶校验  Initial serial port information, 9600 baud rate, 8 data bits, 1 stop bit, no parity |
| 公共方法 | [SerialPortInni(String, Int32)](84aac95f-c526-96bb-95e1-959921d057cd.htm) | 初始化串口信息，波特率，8位数据位，1位停止位，无奇偶校验  Initializes serial port information, baud rate, 8-bit data bit, 1-bit stop bit, no parity |
| 公共方法 | [SerialPortInni(String, Int32, Int32, StopBits, Parity)](26a464c9-69d2-2e04-6217-a7fc8f88035d.htm) | 初始化串口信息，波特率，数据位，停止位，奇偶校验需要全部自己来指定  Start serial port information, baud rate, data bit, stop bit, parity all need to be specified |
| 公共方法 | [ToString](008685e3-0f3c-5219-7ad8-b828a1da9932.htm) | (重写 ObjectToString.) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[ToledoSerial 类](61773df3-22c2-08bf-1635-9cc4d0e62430.htm)

[HslCommunication.Profinet.Toledo 命名空间](3455a9c0-e5c9-7245-1dce-b1fee92e6846.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Close 方法 

[原文連結](http://api.hslcommunication.cn/html/bfcd36b1-f1a0-9377-912b-ed62a44aca0c.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Toledo](../html/3455a9c0-e5c9-7245-1dce-b1fee92e6846.htm "HslCommunication.Profinet.Toledo")

[ToledoSerial 类](../html/61773df3-22c2-08bf-1635-9cc4d0e62430.htm "ToledoSerial 类")

[ToledoSerial 方法](../html/244168d4-6935-958b-3a55-e0c8d7b8ff21.htm "ToledoSerial 方法")

[Close 方法](../html/bfcd36b1-f1a0-9377-912b-ed62a44aca0c.htm "Close 方法 ")

[IsOpen 方法](../html/2e3d66ef-9397-a8b7-1ef9-67a699af15c0.htm "IsOpen 方法 ")

[Open 方法](../html/b5a78c2c-b159-4634-c9ec-ee64d885ff70.htm "Open 方法 ")

[SerialPortInni 方法](../html/df188f59-2279-df42-047a-89b95dce1c46.htm "SerialPortInni 方法 ")

[ToString 方法](../html/008685e3-0f3c-5219-7ad8-b828a1da9932.htm "ToString 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ToledoSerialClose 方法 |

关闭当前的串口连接  
Close the current serial connection

**命名空间：**
 [HslCommunication.Profinet.Toledo](3455a9c0-e5c9-7245-1dce-b1fee92e6846.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public void Close()
```

```
Public Sub Close
```

```
public:
void Close()
```

```
member Close : unit -> unit 
```

![](../icons/SectionExpanded.png)参见

#### 引用

[ToledoSerial 类](61773df3-22c2-08bf-1635-9cc4d0e62430.htm)

[HslCommunication.Profinet.Toledo 命名空间](3455a9c0-e5c9-7245-1dce-b1fee92e6846.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## IsOpen 方法 

[原文連結](http://api.hslcommunication.cn/html/2e3d66ef-9397-a8b7-1ef9-67a699af15c0.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Toledo](../html/3455a9c0-e5c9-7245-1dce-b1fee92e6846.htm "HslCommunication.Profinet.Toledo")

[ToledoSerial 类](../html/61773df3-22c2-08bf-1635-9cc4d0e62430.htm "ToledoSerial 类")

[ToledoSerial 方法](../html/244168d4-6935-958b-3a55-e0c8d7b8ff21.htm "ToledoSerial 方法")

[Close 方法](../html/bfcd36b1-f1a0-9377-912b-ed62a44aca0c.htm "Close 方法 ")

[IsOpen 方法](../html/2e3d66ef-9397-a8b7-1ef9-67a699af15c0.htm "IsOpen 方法 ")

[Open 方法](../html/b5a78c2c-b159-4634-c9ec-ee64d885ff70.htm "Open 方法 ")

[SerialPortInni 方法](../html/df188f59-2279-df42-047a-89b95dce1c46.htm "SerialPortInni 方法 ")

[ToString 方法](../html/008685e3-0f3c-5219-7ad8-b828a1da9932.htm "ToString 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ToledoSerialIsOpen 方法 |

获取一个值，指示串口是否处于打开状态  
Gets a value indicating whether the serial port is open

**命名空间：**
 [HslCommunication.Profinet.Toledo](3455a9c0-e5c9-7245-1dce-b1fee92e6846.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public bool IsOpen()
```

```
Public Function IsOpen As Boolean
```

```
public:
bool IsOpen()
```

```
member IsOpen : unit -> bool 
```

#### 返回值

类型：Boolean  
是或否

![](../icons/SectionExpanded.png)参见

#### 引用

[ToledoSerial 类](61773df3-22c2-08bf-1635-9cc4d0e62430.htm)

[HslCommunication.Profinet.Toledo 命名空间](3455a9c0-e5c9-7245-1dce-b1fee92e6846.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Open 方法 

[原文連結](http://api.hslcommunication.cn/html/b5a78c2c-b159-4634-c9ec-ee64d885ff70.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Toledo](../html/3455a9c0-e5c9-7245-1dce-b1fee92e6846.htm "HslCommunication.Profinet.Toledo")

[ToledoSerial 类](../html/61773df3-22c2-08bf-1635-9cc4d0e62430.htm "ToledoSerial 类")

[ToledoSerial 方法](../html/244168d4-6935-958b-3a55-e0c8d7b8ff21.htm "ToledoSerial 方法")

[Close 方法](../html/bfcd36b1-f1a0-9377-912b-ed62a44aca0c.htm "Close 方法 ")

[IsOpen 方法](../html/2e3d66ef-9397-a8b7-1ef9-67a699af15c0.htm "IsOpen 方法 ")

[Open 方法](../html/b5a78c2c-b159-4634-c9ec-ee64d885ff70.htm "Open 方法 ")

[SerialPortInni 方法](../html/df188f59-2279-df42-047a-89b95dce1c46.htm "SerialPortInni 方法 ")

[ToString 方法](../html/008685e3-0f3c-5219-7ad8-b828a1da9932.htm "ToString 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ToledoSerialOpen 方法 |

打开一个新的串行端口连接  
Open a new serial port connection

**命名空间：**
 [HslCommunication.Profinet.Toledo](3455a9c0-e5c9-7245-1dce-b1fee92e6846.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public void Open()
```

```
Public Sub Open
```

```
public:
void Open()
```

```
member Open : unit -> unit 
```

![](../icons/SectionExpanded.png)参见

#### 引用

[ToledoSerial 类](61773df3-22c2-08bf-1635-9cc4d0e62430.htm)

[HslCommunication.Profinet.Toledo 命名空间](3455a9c0-e5c9-7245-1dce-b1fee92e6846.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## SerialPortInni 方法 

[原文連結](http://api.hslcommunication.cn/html/df188f59-2279-df42-047a-89b95dce1c46.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Toledo](../html/3455a9c0-e5c9-7245-1dce-b1fee92e6846.htm "HslCommunication.Profinet.Toledo")

[ToledoSerial 类](../html/61773df3-22c2-08bf-1635-9cc4d0e62430.htm "ToledoSerial 类")

[ToledoSerial 方法](../html/244168d4-6935-958b-3a55-e0c8d7b8ff21.htm "ToledoSerial 方法")

[SerialPortInni 方法](../html/df188f59-2279-df42-047a-89b95dce1c46.htm "SerialPortInni 方法 ")

[SerialPortInni 方法 (Action(SerialPort))](../html/82011d60-4795-6e72-51b1-be0bd201d530.htm "SerialPortInni 方法 (Action(SerialPort))")

[SerialPortInni 方法 (String)](../html/b09ded66-4339-b33e-8828-4281fbe0b28a.htm "SerialPortInni 方法 (String)")

[SerialPortInni 方法 (String, Int32)](../html/84aac95f-c526-96bb-95e1-959921d057cd.htm "SerialPortInni 方法 (String, Int32)")

[SerialPortInni 方法 (String, Int32, Int32, StopBits, Parity)](../html/26a464c9-69d2-2e04-6217-a7fc8f88035d.htm "SerialPortInni 方法 (String, Int32, Int32, StopBits, Parity)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ToledoSerialSerialPortInni 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [SerialPortInni(ActionSerialPort)](82011d60-4795-6e72-51b1-be0bd201d530.htm) | 根据自定义初始化方法进行初始化串口信息  Initialize the serial port information according to the custom initialization method |
| 公共方法 | [SerialPortInni(String)](b09ded66-4339-b33e-8828-4281fbe0b28a.htm) | 初始化串口信息，9600波特率，8位数据位，1位停止位，无奇偶校验  Initial serial port information, 9600 baud rate, 8 data bits, 1 stop bit, no parity |
| 公共方法 | [SerialPortInni(String, Int32)](84aac95f-c526-96bb-95e1-959921d057cd.htm) | 初始化串口信息，波特率，8位数据位，1位停止位，无奇偶校验  Initializes serial port information, baud rate, 8-bit data bit, 1-bit stop bit, no parity |
| 公共方法 | [SerialPortInni(String, Int32, Int32, StopBits, Parity)](26a464c9-69d2-2e04-6217-a7fc8f88035d.htm) | 初始化串口信息，波特率，数据位，停止位，奇偶校验需要全部自己来指定  Start serial port information, baud rate, data bit, stop bit, parity all need to be specified |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[ToledoSerial 类](61773df3-22c2-08bf-1635-9cc4d0e62430.htm)

[HslCommunication.Profinet.Toledo 命名空间](3455a9c0-e5c9-7245-1dce-b1fee92e6846.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## SerialPortInni 方法 (Action(SerialPort))

[原文連結](http://api.hslcommunication.cn/html/82011d60-4795-6e72-51b1-be0bd201d530.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Toledo](../html/3455a9c0-e5c9-7245-1dce-b1fee92e6846.htm "HslCommunication.Profinet.Toledo")

[ToledoSerial 类](../html/61773df3-22c2-08bf-1635-9cc4d0e62430.htm "ToledoSerial 类")

[ToledoSerial 方法](../html/244168d4-6935-958b-3a55-e0c8d7b8ff21.htm "ToledoSerial 方法")

[SerialPortInni 方法](../html/df188f59-2279-df42-047a-89b95dce1c46.htm "SerialPortInni 方法 ")

[SerialPortInni 方法 (Action(SerialPort))](../html/82011d60-4795-6e72-51b1-be0bd201d530.htm "SerialPortInni 方法 (Action(SerialPort))")

[SerialPortInni 方法 (String)](../html/b09ded66-4339-b33e-8828-4281fbe0b28a.htm "SerialPortInni 方法 (String)")

[SerialPortInni 方法 (String, Int32)](../html/84aac95f-c526-96bb-95e1-959921d057cd.htm "SerialPortInni 方法 (String, Int32)")

[SerialPortInni 方法 (String, Int32, Int32, StopBits, Parity)](../html/26a464c9-69d2-2e04-6217-a7fc8f88035d.htm "SerialPortInni 方法 (String, Int32, Int32, StopBits, Parity)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ToledoSerialSerialPortInni 方法 (ActionSerialPort) |

根据自定义初始化方法进行初始化串口信息  
Initialize the serial port information according to the custom initialization method

**命名空间：**
 [HslCommunication.Profinet.Toledo](3455a9c0-e5c9-7245-1dce-b1fee92e6846.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public void SerialPortInni(
	Action<SerialPort> initi
)
```

```
Public Sub SerialPortInni ( 
	initi As Action(Of SerialPort)
)
```

```
public:
void SerialPortInni(
	Action<SerialPort^>^ initi
)
```

```
member SerialPortInni : 
        initi : Action<SerialPort> -> unit 
```

#### 参数

initi
:   类型：SystemActionSerialPort  
    初始化的委托方法

![](../icons/SectionExpanded.png)参见

#### 引用

[ToledoSerial 类](61773df3-22c2-08bf-1635-9cc4d0e62430.htm)

[SerialPortInni 重载](df188f59-2279-df42-047a-89b95dce1c46.htm)

[HslCommunication.Profinet.Toledo 命名空间](3455a9c0-e5c9-7245-1dce-b1fee92e6846.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## SerialPortInni 方法 (String)

[原文連結](http://api.hslcommunication.cn/html/b09ded66-4339-b33e-8828-4281fbe0b28a.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Toledo](../html/3455a9c0-e5c9-7245-1dce-b1fee92e6846.htm "HslCommunication.Profinet.Toledo")

[ToledoSerial 类](../html/61773df3-22c2-08bf-1635-9cc4d0e62430.htm "ToledoSerial 类")

[ToledoSerial 方法](../html/244168d4-6935-958b-3a55-e0c8d7b8ff21.htm "ToledoSerial 方法")

[SerialPortInni 方法](../html/df188f59-2279-df42-047a-89b95dce1c46.htm "SerialPortInni 方法 ")

[SerialPortInni 方法 (Action(SerialPort))](../html/82011d60-4795-6e72-51b1-be0bd201d530.htm "SerialPortInni 方法 (Action(SerialPort))")

[SerialPortInni 方法 (String)](../html/b09ded66-4339-b33e-8828-4281fbe0b28a.htm "SerialPortInni 方法 (String)")

[SerialPortInni 方法 (String, Int32)](../html/84aac95f-c526-96bb-95e1-959921d057cd.htm "SerialPortInni 方法 (String, Int32)")

[SerialPortInni 方法 (String, Int32, Int32, StopBits, Parity)](../html/26a464c9-69d2-2e04-6217-a7fc8f88035d.htm "SerialPortInni 方法 (String, Int32, Int32, StopBits, Parity)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ToledoSerialSerialPortInni 方法 (String) |

初始化串口信息，9600波特率，8位数据位，1位停止位，无奇偶校验  
Initial serial port information, 9600 baud rate, 8 data bits, 1 stop bit, no parity

**命名空间：**
 [HslCommunication.Profinet.Toledo](3455a9c0-e5c9-7245-1dce-b1fee92e6846.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public void SerialPortInni(
	string portName
)
```

```
Public Sub SerialPortInni ( 
	portName As String
)
```

```
public:
void SerialPortInni(
	String^ portName
)
```

```
member SerialPortInni : 
        portName : string -> unit 
```

#### 参数

portName
:   类型：SystemString  
    端口号信息，例如"COM3"

![](../icons/SectionExpanded.png)参见

#### 引用

[ToledoSerial 类](61773df3-22c2-08bf-1635-9cc4d0e62430.htm)

[SerialPortInni 重载](df188f59-2279-df42-047a-89b95dce1c46.htm)

[HslCommunication.Profinet.Toledo 命名空间](3455a9c0-e5c9-7245-1dce-b1fee92e6846.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## SerialPortInni 方法 (String, Int32)

[原文連結](http://api.hslcommunication.cn/html/84aac95f-c526-96bb-95e1-959921d057cd.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Toledo](../html/3455a9c0-e5c9-7245-1dce-b1fee92e6846.htm "HslCommunication.Profinet.Toledo")

[ToledoSerial 类](../html/61773df3-22c2-08bf-1635-9cc4d0e62430.htm "ToledoSerial 类")

[ToledoSerial 方法](../html/244168d4-6935-958b-3a55-e0c8d7b8ff21.htm "ToledoSerial 方法")

[SerialPortInni 方法](../html/df188f59-2279-df42-047a-89b95dce1c46.htm "SerialPortInni 方法 ")

[SerialPortInni 方法 (Action(SerialPort))](../html/82011d60-4795-6e72-51b1-be0bd201d530.htm "SerialPortInni 方法 (Action(SerialPort))")

[SerialPortInni 方法 (String)](../html/b09ded66-4339-b33e-8828-4281fbe0b28a.htm "SerialPortInni 方法 (String)")

[SerialPortInni 方法 (String, Int32)](../html/84aac95f-c526-96bb-95e1-959921d057cd.htm "SerialPortInni 方法 (String, Int32)")

[SerialPortInni 方法 (String, Int32, Int32, StopBits, Parity)](../html/26a464c9-69d2-2e04-6217-a7fc8f88035d.htm "SerialPortInni 方法 (String, Int32, Int32, StopBits, Parity)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ToledoSerialSerialPortInni 方法 (String, Int32) |

初始化串口信息，波特率，8位数据位，1位停止位，无奇偶校验  
Initializes serial port information, baud rate, 8-bit data bit, 1-bit stop bit, no parity

**命名空间：**
 [HslCommunication.Profinet.Toledo](3455a9c0-e5c9-7245-1dce-b1fee92e6846.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public void SerialPortInni(
	string portName,
	int baudRate
)
```

```
Public Sub SerialPortInni ( 
	portName As String,
	baudRate As Integer
)
```

```
public:
void SerialPortInni(
	String^ portName, 
	int baudRate
)
```

```
member SerialPortInni : 
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

[ToledoSerial 类](61773df3-22c2-08bf-1635-9cc4d0e62430.htm)

[SerialPortInni 重载](df188f59-2279-df42-047a-89b95dce1c46.htm)

[HslCommunication.Profinet.Toledo 命名空间](3455a9c0-e5c9-7245-1dce-b1fee92e6846.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## SerialPortInni 方法 (String, Int32, Int32, StopBits, Parity)

[原文連結](http://api.hslcommunication.cn/html/26a464c9-69d2-2e04-6217-a7fc8f88035d.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Toledo](../html/3455a9c0-e5c9-7245-1dce-b1fee92e6846.htm "HslCommunication.Profinet.Toledo")

[ToledoSerial 类](../html/61773df3-22c2-08bf-1635-9cc4d0e62430.htm "ToledoSerial 类")

[ToledoSerial 方法](../html/244168d4-6935-958b-3a55-e0c8d7b8ff21.htm "ToledoSerial 方法")

[SerialPortInni 方法](../html/df188f59-2279-df42-047a-89b95dce1c46.htm "SerialPortInni 方法 ")

[SerialPortInni 方法 (Action(SerialPort))](../html/82011d60-4795-6e72-51b1-be0bd201d530.htm "SerialPortInni 方法 (Action(SerialPort))")

[SerialPortInni 方法 (String)](../html/b09ded66-4339-b33e-8828-4281fbe0b28a.htm "SerialPortInni 方法 (String)")

[SerialPortInni 方法 (String, Int32)](../html/84aac95f-c526-96bb-95e1-959921d057cd.htm "SerialPortInni 方法 (String, Int32)")

[SerialPortInni 方法 (String, Int32, Int32, StopBits, Parity)](../html/26a464c9-69d2-2e04-6217-a7fc8f88035d.htm "SerialPortInni 方法 (String, Int32, Int32, StopBits, Parity)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ToledoSerialSerialPortInni 方法 (String, Int32, Int32, StopBits, Parity) |

初始化串口信息，波特率，数据位，停止位，奇偶校验需要全部自己来指定  
Start serial port information, baud rate, data bit, stop bit, parity all need to be specified

**命名空间：**
 [HslCommunication.Profinet.Toledo](3455a9c0-e5c9-7245-1dce-b1fee92e6846.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public void SerialPortInni(
	string portName,
	int baudRate,
	int dataBits,
	StopBits stopBits,
	Parity parity
)
```

```
Public Sub SerialPortInni ( 
	portName As String,
	baudRate As Integer,
	dataBits As Integer,
	stopBits As StopBits,
	parity As Parity
)
```

```
public:
void SerialPortInni(
	String^ portName, 
	int baudRate, 
	int dataBits, 
	StopBits stopBits, 
	Parity parity
)
```

```
member SerialPortInni : 
        portName : string * 
        baudRate : int * 
        dataBits : int * 
        stopBits : StopBits * 
        parity : Parity -> unit 
```

#### 参数

portName
:   类型：SystemString  
    端口号信息，例如"COM3"

baudRate
:   类型：SystemInt32  
    波特率

dataBits
:   类型：SystemInt32  
    数据位

stopBits
:   类型：System.IO.PortsStopBits  
    停止位

parity
:   类型：System.IO.PortsParity  
    奇偶校验

![](../icons/SectionExpanded.png)参见

#### 引用

[ToledoSerial 类](61773df3-22c2-08bf-1635-9cc4d0e62430.htm)

[SerialPortInni 重载](df188f59-2279-df42-047a-89b95dce1c46.htm)

[HslCommunication.Profinet.Toledo 命名空间](3455a9c0-e5c9-7245-1dce-b1fee92e6846.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ToString 方法 

[原文連結](http://api.hslcommunication.cn/html/008685e3-0f3c-5219-7ad8-b828a1da9932.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Toledo](../html/3455a9c0-e5c9-7245-1dce-b1fee92e6846.htm "HslCommunication.Profinet.Toledo")

[ToledoSerial 类](../html/61773df3-22c2-08bf-1635-9cc4d0e62430.htm "ToledoSerial 类")

[ToledoSerial 方法](../html/244168d4-6935-958b-3a55-e0c8d7b8ff21.htm "ToledoSerial 方法")

[Close 方法](../html/bfcd36b1-f1a0-9377-912b-ed62a44aca0c.htm "Close 方法 ")

[IsOpen 方法](../html/2e3d66ef-9397-a8b7-1ef9-67a699af15c0.htm "IsOpen 方法 ")

[Open 方法](../html/b5a78c2c-b159-4634-c9ec-ee64d885ff70.htm "Open 方法 ")

[SerialPortInni 方法](../html/df188f59-2279-df42-047a-89b95dce1c46.htm "SerialPortInni 方法 ")

[ToString 方法](../html/008685e3-0f3c-5219-7ad8-b828a1da9932.htm "ToString 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ToledoSerialToString 方法 |

[缺少 "M:HslCommunication.Profinet.Toledo.ToledoSerial.ToString" 的 <summary> 文档]

**命名空间：**
 [HslCommunication.Profinet.Toledo](3455a9c0-e5c9-7245-1dce-b1fee92e6846.htm)  
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

[缺少 "M:HslCommunication.Profinet.Toledo.ToledoSerial.ToString" 的 <returns> 文档]

![](../icons/SectionExpanded.png)参见

#### 引用

[ToledoSerial 类](61773df3-22c2-08bf-1635-9cc4d0e62430.htm)

[HslCommunication.Profinet.Toledo 命名空间](3455a9c0-e5c9-7245-1dce-b1fee92e6846.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ToledoSerial 事件

[原文連結](http://api.hslcommunication.cn/html/5833568b-8df7-5029-ce9c-0daaf2616061.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Toledo](../html/3455a9c0-e5c9-7245-1dce-b1fee92e6846.htm "HslCommunication.Profinet.Toledo")

[ToledoSerial 类](../html/61773df3-22c2-08bf-1635-9cc4d0e62430.htm "ToledoSerial 类")

[ToledoSerial 事件](../html/5833568b-8df7-5029-ce9c-0daaf2616061.htm "ToledoSerial 事件")

[OnToledoStandardDataReceived 事件](../html/de23564d-f97b-79ba-9236-36651a923868.htm "OnToledoStandardDataReceived 事件")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ToledoSerial 事件 |

[ToledoSerial](61773df3-22c2-08bf-1635-9cc4d0e62430.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)事件

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共事件 | [OnToledoStandardDataReceived](de23564d-f97b-79ba-9236-36651a923868.htm) | 当接收到一条新的托利多的数据的时候触发 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[ToledoSerial 类](61773df3-22c2-08bf-1635-9cc4d0e62430.htm)

[HslCommunication.Profinet.Toledo 命名空间](3455a9c0-e5c9-7245-1dce-b1fee92e6846.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## OnToledoStandardDataReceived 事件

[原文連結](http://api.hslcommunication.cn/html/de23564d-f97b-79ba-9236-36651a923868.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Toledo](../html/3455a9c0-e5c9-7245-1dce-b1fee92e6846.htm "HslCommunication.Profinet.Toledo")

[ToledoSerial 类](../html/61773df3-22c2-08bf-1635-9cc4d0e62430.htm "ToledoSerial 类")

[ToledoSerial 事件](../html/5833568b-8df7-5029-ce9c-0daaf2616061.htm "ToledoSerial 事件")

[OnToledoStandardDataReceived 事件](../html/de23564d-f97b-79ba-9236-36651a923868.htm "OnToledoStandardDataReceived 事件")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ToledoSerialOnToledoStandardDataReceived 事件 |

当接收到一条新的托利多的数据的时候触发

**命名空间：**
 [HslCommunication.Profinet.Toledo](3455a9c0-e5c9-7245-1dce-b1fee92e6846.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public event ToledoSerialToledoStandardDataReceivedDelegate OnToledoStandardDataReceived
```

```
Public Event OnToledoStandardDataReceived As ToledoSerialToledoStandardDataReceivedDelegate
```

```
public:
 event ToledoSerialToledoStandardDataReceivedDelegate^ OnToledoStandardDataReceived {
	void add (ToledoSerialToledoStandardDataReceivedDelegate^ value);
	void remove (ToledoSerialToledoStandardDataReceivedDelegate^ value);
}
```

```
member OnToledoStandardDataReceived : IEvent<ToledoSerialToledoStandardDataReceivedDelegate,
    ToledoStandardData>
```

#### 值

类型：[HslCommunication.Profinet.ToledoToledoSerialToledoStandardDataReceivedDelegate](4d9acb69-0a01-8cba-3658-6977be61668b.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[ToledoSerial 类](61773df3-22c2-08bf-1635-9cc4d0e62430.htm)

[HslCommunication.Profinet.Toledo 命名空间](3455a9c0-e5c9-7245-1dce-b1fee92e6846.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ToledoSerial.ToledoStandardDataReceivedDelegate 委托

[原文連結](http://api.hslcommunication.cn/html/4d9acb69-0a01-8cba-3658-6977be61668b.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Toledo](../html/3455a9c0-e5c9-7245-1dce-b1fee92e6846.htm "HslCommunication.Profinet.Toledo")

[ToledoSerial 类](../html/61773df3-22c2-08bf-1635-9cc4d0e62430.htm "ToledoSerial 类")

[ToledoSerial.ToledoStandardDataReceivedDelegate 委托](../html/4d9acb69-0a01-8cba-3658-6977be61668b.htm "ToledoSerial.ToledoStandardDataReceivedDelegate 委托")

[ToledoStandardData 类](../html/13fdb4fb-3312-d193-7e77-a383ddffbf1d.htm "ToledoStandardData 类")

[ToledoTcpServer 类](../html/32226f81-8b8f-228c-4a5f-55dedb6ab8e2.htm "ToledoTcpServer 类")

[ToledoTcpServer.ToledoStandardDataReceivedDelegate 委托](../html/49d9d452-7681-6598-36ad-c6c03053782b.htm "ToledoTcpServer.ToledoStandardDataReceivedDelegate 委托")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ToledoSerialToledoStandardDataReceivedDelegate 委托 |

托利多数据接收时的委托

**命名空间：**
 [HslCommunication.Profinet.Toledo](3455a9c0-e5c9-7245-1dce-b1fee92e6846.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public delegate void ToledoStandardDataReceivedDelegate(
	Object sender,
	ToledoStandardData toledoStandardData
)
```

```
Public Delegate Sub ToledoStandardDataReceivedDelegate ( 
	sender As Object,
	toledoStandardData As ToledoStandardData
)
```

```
public delegate void ToledoStandardDataReceivedDelegate(
	Object^ sender, 
	ToledoStandardData^ toledoStandardData
)
```

```
type ToledoStandardDataReceivedDelegate = 
    delegate of 
        sender : Object * 
        toledoStandardData : ToledoStandardData -> unit
```

#### 参数

sender
:   类型：SystemObject  
    数据发送对象

toledoStandardData
:   类型：[HslCommunication.Profinet.ToledoToledoStandardData](13fdb4fb-3312-d193-7e77-a383ddffbf1d.htm)  
    数据对象

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.Profinet.Toledo 命名空间](3455a9c0-e5c9-7245-1dce-b1fee92e6846.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ToledoStandardData 类

[原文連結](http://api.hslcommunication.cn/html/13fdb4fb-3312-d193-7e77-a383ddffbf1d.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Toledo](../html/3455a9c0-e5c9-7245-1dce-b1fee92e6846.htm "HslCommunication.Profinet.Toledo")

[ToledoStandardData 类](../html/13fdb4fb-3312-d193-7e77-a383ddffbf1d.htm "ToledoStandardData 类")

[ToledoStandardData 构造函数](../html/429b9ba6-8efb-64b7-7322-ea9f6fcb181c.htm "ToledoStandardData 构造函数 ")

[ToledoStandardData 属性](../html/3fe81910-7894-7631-de26-d2ac7f140782.htm "ToledoStandardData 属性")

[ToledoStandardData 方法](../html/4edb7686-7e68-aeba-a551-a75a3397729e.htm "ToledoStandardData 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ToledoStandardData 类 |

托利多标准格式的数据类对象

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  HslCommunication.Profinet.ToledoToledoStandardData

**命名空间：**
 [HslCommunication.Profinet.Toledo](3455a9c0-e5c9-7245-1dce-b1fee92e6846.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class ToledoStandardData
```

```
Public Class ToledoStandardData
```

```
public ref class ToledoStandardData
```

```
type ToledoStandardData =  class end
```

ToledoStandardData 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [ToledoStandardData](99d58cda-b94d-381f-5b40-937653dfb4df.htm) | 实例化一个默认的对象 |
| 公共方法 | [ToledoStandardData(Byte)](e1f2a774-85b0-6c74-89de-40bc56ebf1e9.htm) | 从缓存里加载一个标准格式的对象 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [BeyondScope](465a0314-8a63-a00c-b04f-4d65a6af3a1e.htm) | 是否在范围之外 |
| 公共属性 | [DataValid](fd26e9ee-6a81-1372-53c6-4429ddbf94e3.htm) | 数据是否有效 |
| 公共属性 | [DynamicState](77a3a834-bad9-6a3d-c649-beb1b38fc0f9.htm) | 是否为动态，为 True 则是动态，为 False 则为稳态 |
| 公共属性 | [IsExpandOutput](37149140-100a-0f52-5c0f-d6f8c0ec208d.htm) | 是否属于扩展输出模式 |
| 公共属性 | [IsPrint](2122f950-b596-0009-dbce-3a66e1b1f84c.htm) | 是否打印 |
| 公共属性 | [IsTenExtend](b4cebc10-ab06-b5fc-513c-9fd5ef2fd3c4.htm) | 是否10被扩展 |
| 公共属性 | [SourceData](d478b7f6-ca65-6aa0-3156-581a1cd68a33.htm) | 解析数据的原始字节 |
| 公共属性 | [Suttle](e1959b5f-89e8-17fe-d07b-4ceee83e0e7e.htm) | 为 True 则是净重，为 False 则为毛重 |
| 公共属性 | [Symbol](df9ef7d7-6a4e-04a0-6dbd-bfb6a4100ada.htm) | 为 True 则是正，为 False 则为负 |
| 公共属性 | [Tare](3c34323e-f46e-bf12-6780-9e5ceb60384d.htm) | 皮重 |
| 公共属性 | [TareType](9eb382bf-639d-c4a2-d11c-6a50400cd482.htm) | 皮重类型，0: 无皮重; 1: 按键去皮; 2: 预置去皮; 3: 皮重内存，仅在扩展输出下有效 |
| 公共属性 | [Unit](4f500913-4b2e-da80-459d-e4513c0a5a3d.htm) | 单位 |
| 公共属性 | [Weight](089a04ce-cb69-50fa-c9b4-508768d0c9dd.htm) | 重量 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法 | [ToString](d29e8058-63fc-1179-fb87-32408dfe9d2d.htm) | (重写 ObjectToString.) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.Profinet.Toledo 命名空间](3455a9c0-e5c9-7245-1dce-b1fee92e6846.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ToledoStandardData 构造函数 

[原文連結](http://api.hslcommunication.cn/html/429b9ba6-8efb-64b7-7322-ea9f6fcb181c.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Toledo](../html/3455a9c0-e5c9-7245-1dce-b1fee92e6846.htm "HslCommunication.Profinet.Toledo")

[ToledoStandardData 类](../html/13fdb4fb-3312-d193-7e77-a383ddffbf1d.htm "ToledoStandardData 类")

[ToledoStandardData 构造函数](../html/429b9ba6-8efb-64b7-7322-ea9f6fcb181c.htm "ToledoStandardData 构造函数 ")

[ToledoStandardData 构造函数](../html/99d58cda-b94d-381f-5b40-937653dfb4df.htm "ToledoStandardData 构造函数 ")

[ToledoStandardData 构造函数 (Byte[])](../html/e1f2a774-85b0-6c74-89de-40bc56ebf1e9.htm "ToledoStandardData 构造函数 (Byte[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ToledoStandardData 构造函数 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [ToledoStandardData](99d58cda-b94d-381f-5b40-937653dfb4df.htm) | 实例化一个默认的对象 |
| 公共方法 | [ToledoStandardData(Byte)](e1f2a774-85b0-6c74-89de-40bc56ebf1e9.htm) | 从缓存里加载一个标准格式的对象 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[ToledoStandardData 类](13fdb4fb-3312-d193-7e77-a383ddffbf1d.htm)

[HslCommunication.Profinet.Toledo 命名空间](3455a9c0-e5c9-7245-1dce-b1fee92e6846.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ToledoStandardData 构造函数 

[原文連結](http://api.hslcommunication.cn/html/99d58cda-b94d-381f-5b40-937653dfb4df.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Toledo](../html/3455a9c0-e5c9-7245-1dce-b1fee92e6846.htm "HslCommunication.Profinet.Toledo")

[ToledoStandardData 类](../html/13fdb4fb-3312-d193-7e77-a383ddffbf1d.htm "ToledoStandardData 类")

[ToledoStandardData 构造函数](../html/429b9ba6-8efb-64b7-7322-ea9f6fcb181c.htm "ToledoStandardData 构造函数 ")

[ToledoStandardData 构造函数](../html/99d58cda-b94d-381f-5b40-937653dfb4df.htm "ToledoStandardData 构造函数 ")

[ToledoStandardData 构造函数 (Byte[])](../html/e1f2a774-85b0-6c74-89de-40bc56ebf1e9.htm "ToledoStandardData 构造函数 (Byte[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ToledoStandardData 构造函数 |

实例化一个默认的对象

**命名空间：**
 [HslCommunication.Profinet.Toledo](3455a9c0-e5c9-7245-1dce-b1fee92e6846.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public ToledoStandardData()
```

```
Public Sub New
```

```
public:
ToledoStandardData()
```

```
new : unit -> ToledoStandardData
```

![](../icons/SectionExpanded.png)参见

#### 引用

[ToledoStandardData 类](13fdb4fb-3312-d193-7e77-a383ddffbf1d.htm)

[ToledoStandardData 重载](429b9ba6-8efb-64b7-7322-ea9f6fcb181c.htm)

[HslCommunication.Profinet.Toledo 命名空间](3455a9c0-e5c9-7245-1dce-b1fee92e6846.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ToledoStandardData 构造函数 (Byte[])

[原文連結](http://api.hslcommunication.cn/html/e1f2a774-85b0-6c74-89de-40bc56ebf1e9.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Toledo](../html/3455a9c0-e5c9-7245-1dce-b1fee92e6846.htm "HslCommunication.Profinet.Toledo")

[ToledoStandardData 类](../html/13fdb4fb-3312-d193-7e77-a383ddffbf1d.htm "ToledoStandardData 类")

[ToledoStandardData 构造函数](../html/429b9ba6-8efb-64b7-7322-ea9f6fcb181c.htm "ToledoStandardData 构造函数 ")

[ToledoStandardData 构造函数](../html/99d58cda-b94d-381f-5b40-937653dfb4df.htm "ToledoStandardData 构造函数 ")

[ToledoStandardData 构造函数 (Byte[])](../html/e1f2a774-85b0-6c74-89de-40bc56ebf1e9.htm "ToledoStandardData 构造函数 (Byte[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ToledoStandardData 构造函数 (Byte) |

从缓存里加载一个标准格式的对象

**命名空间：**
 [HslCommunication.Profinet.Toledo](3455a9c0-e5c9-7245-1dce-b1fee92e6846.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public ToledoStandardData(
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
ToledoStandardData(
	array<unsigned char>^ buffer
)
```

```
new : 
        buffer : byte[] -> ToledoStandardData
```

#### 参数

buffer
:   类型：SystemByte  
    缓存

![](../icons/SectionExpanded.png)参见

#### 引用

[ToledoStandardData 类](13fdb4fb-3312-d193-7e77-a383ddffbf1d.htm)

[ToledoStandardData 重载](429b9ba6-8efb-64b7-7322-ea9f6fcb181c.htm)

[HslCommunication.Profinet.Toledo 命名空间](3455a9c0-e5c9-7245-1dce-b1fee92e6846.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ToledoStandardData 属性

[原文連結](http://api.hslcommunication.cn/html/3fe81910-7894-7631-de26-d2ac7f140782.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Toledo](../html/3455a9c0-e5c9-7245-1dce-b1fee92e6846.htm "HslCommunication.Profinet.Toledo")

[ToledoStandardData 类](../html/13fdb4fb-3312-d193-7e77-a383ddffbf1d.htm "ToledoStandardData 类")

[ToledoStandardData 属性](../html/3fe81910-7894-7631-de26-d2ac7f140782.htm "ToledoStandardData 属性")

[BeyondScope 属性](../html/465a0314-8a63-a00c-b04f-4d65a6af3a1e.htm "BeyondScope 属性 ")

[DataValid 属性](../html/fd26e9ee-6a81-1372-53c6-4429ddbf94e3.htm "DataValid 属性 ")

[DynamicState 属性](../html/77a3a834-bad9-6a3d-c649-beb1b38fc0f9.htm "DynamicState 属性 ")

[IsExpandOutput 属性](../html/37149140-100a-0f52-5c0f-d6f8c0ec208d.htm "IsExpandOutput 属性 ")

[IsPrint 属性](../html/2122f950-b596-0009-dbce-3a66e1b1f84c.htm "IsPrint 属性 ")

[IsTenExtend 属性](../html/b4cebc10-ab06-b5fc-513c-9fd5ef2fd3c4.htm "IsTenExtend 属性 ")

[SourceData 属性](../html/d478b7f6-ca65-6aa0-3156-581a1cd68a33.htm "SourceData 属性 ")

[Suttle 属性](../html/e1959b5f-89e8-17fe-d07b-4ceee83e0e7e.htm "Suttle 属性 ")

[Symbol 属性](../html/df9ef7d7-6a4e-04a0-6dbd-bfb6a4100ada.htm "Symbol 属性 ")

[Tare 属性](../html/3c34323e-f46e-bf12-6780-9e5ceb60384d.htm "Tare 属性 ")

[TareType 属性](../html/9eb382bf-639d-c4a2-d11c-6a50400cd482.htm "TareType 属性 ")

[Unit 属性](../html/4f500913-4b2e-da80-459d-e4513c0a5a3d.htm "Unit 属性 ")

[Weight 属性](../html/089a04ce-cb69-50fa-c9b4-508768d0c9dd.htm "Weight 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ToledoStandardData 属性 |

[ToledoStandardData](13fdb4fb-3312-d193-7e77-a383ddffbf1d.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [BeyondScope](465a0314-8a63-a00c-b04f-4d65a6af3a1e.htm) | 是否在范围之外 |
| 公共属性 | [DataValid](fd26e9ee-6a81-1372-53c6-4429ddbf94e3.htm) | 数据是否有效 |
| 公共属性 | [DynamicState](77a3a834-bad9-6a3d-c649-beb1b38fc0f9.htm) | 是否为动态，为 True 则是动态，为 False 则为稳态 |
| 公共属性 | [IsExpandOutput](37149140-100a-0f52-5c0f-d6f8c0ec208d.htm) | 是否属于扩展输出模式 |
| 公共属性 | [IsPrint](2122f950-b596-0009-dbce-3a66e1b1f84c.htm) | 是否打印 |
| 公共属性 | [IsTenExtend](b4cebc10-ab06-b5fc-513c-9fd5ef2fd3c4.htm) | 是否10被扩展 |
| 公共属性 | [SourceData](d478b7f6-ca65-6aa0-3156-581a1cd68a33.htm) | 解析数据的原始字节 |
| 公共属性 | [Suttle](e1959b5f-89e8-17fe-d07b-4ceee83e0e7e.htm) | 为 True 则是净重，为 False 则为毛重 |
| 公共属性 | [Symbol](df9ef7d7-6a4e-04a0-6dbd-bfb6a4100ada.htm) | 为 True 则是正，为 False 则为负 |
| 公共属性 | [Tare](3c34323e-f46e-bf12-6780-9e5ceb60384d.htm) | 皮重 |
| 公共属性 | [TareType](9eb382bf-639d-c4a2-d11c-6a50400cd482.htm) | 皮重类型，0: 无皮重; 1: 按键去皮; 2: 预置去皮; 3: 皮重内存，仅在扩展输出下有效 |
| 公共属性 | [Unit](4f500913-4b2e-da80-459d-e4513c0a5a3d.htm) | 单位 |
| 公共属性 | [Weight](089a04ce-cb69-50fa-c9b4-508768d0c9dd.htm) | 重量 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[ToledoStandardData 类](13fdb4fb-3312-d193-7e77-a383ddffbf1d.htm)

[HslCommunication.Profinet.Toledo 命名空间](3455a9c0-e5c9-7245-1dce-b1fee92e6846.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## BeyondScope 属性 

[原文連結](http://api.hslcommunication.cn/html/465a0314-8a63-a00c-b04f-4d65a6af3a1e.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Toledo](../html/3455a9c0-e5c9-7245-1dce-b1fee92e6846.htm "HslCommunication.Profinet.Toledo")

[ToledoStandardData 类](../html/13fdb4fb-3312-d193-7e77-a383ddffbf1d.htm "ToledoStandardData 类")

[ToledoStandardData 属性](../html/3fe81910-7894-7631-de26-d2ac7f140782.htm "ToledoStandardData 属性")

[BeyondScope 属性](../html/465a0314-8a63-a00c-b04f-4d65a6af3a1e.htm "BeyondScope 属性 ")

[DataValid 属性](../html/fd26e9ee-6a81-1372-53c6-4429ddbf94e3.htm "DataValid 属性 ")

[DynamicState 属性](../html/77a3a834-bad9-6a3d-c649-beb1b38fc0f9.htm "DynamicState 属性 ")

[IsExpandOutput 属性](../html/37149140-100a-0f52-5c0f-d6f8c0ec208d.htm "IsExpandOutput 属性 ")

[IsPrint 属性](../html/2122f950-b596-0009-dbce-3a66e1b1f84c.htm "IsPrint 属性 ")

[IsTenExtend 属性](../html/b4cebc10-ab06-b5fc-513c-9fd5ef2fd3c4.htm "IsTenExtend 属性 ")

[SourceData 属性](../html/d478b7f6-ca65-6aa0-3156-581a1cd68a33.htm "SourceData 属性 ")

[Suttle 属性](../html/e1959b5f-89e8-17fe-d07b-4ceee83e0e7e.htm "Suttle 属性 ")

[Symbol 属性](../html/df9ef7d7-6a4e-04a0-6dbd-bfb6a4100ada.htm "Symbol 属性 ")

[Tare 属性](../html/3c34323e-f46e-bf12-6780-9e5ceb60384d.htm "Tare 属性 ")

[TareType 属性](../html/9eb382bf-639d-c4a2-d11c-6a50400cd482.htm "TareType 属性 ")

[Unit 属性](../html/4f500913-4b2e-da80-459d-e4513c0a5a3d.htm "Unit 属性 ")

[Weight 属性](../html/089a04ce-cb69-50fa-c9b4-508768d0c9dd.htm "Weight 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ToledoStandardDataBeyondScope 属性 |

是否在范围之外

**命名空间：**
 [HslCommunication.Profinet.Toledo](3455a9c0-e5c9-7245-1dce-b1fee92e6846.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public bool BeyondScope { get; set; }
```

```
Public Property BeyondScope As Boolean
	Get
	Set
```

```
public:
property bool BeyondScope {
	bool get ();
	void set (bool value);
}
```

```
member BeyondScope : bool with get, set
```

#### 属性值

类型：Boolean

![](../icons/SectionExpanded.png)参见

#### 引用

[ToledoStandardData 类](13fdb4fb-3312-d193-7e77-a383ddffbf1d.htm)

[HslCommunication.Profinet.Toledo 命名空间](3455a9c0-e5c9-7245-1dce-b1fee92e6846.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## DataValid 属性 

[原文連結](http://api.hslcommunication.cn/html/fd26e9ee-6a81-1372-53c6-4429ddbf94e3.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Toledo](../html/3455a9c0-e5c9-7245-1dce-b1fee92e6846.htm "HslCommunication.Profinet.Toledo")

[ToledoStandardData 类](../html/13fdb4fb-3312-d193-7e77-a383ddffbf1d.htm "ToledoStandardData 类")

[ToledoStandardData 属性](../html/3fe81910-7894-7631-de26-d2ac7f140782.htm "ToledoStandardData 属性")

[BeyondScope 属性](../html/465a0314-8a63-a00c-b04f-4d65a6af3a1e.htm "BeyondScope 属性 ")

[DataValid 属性](../html/fd26e9ee-6a81-1372-53c6-4429ddbf94e3.htm "DataValid 属性 ")

[DynamicState 属性](../html/77a3a834-bad9-6a3d-c649-beb1b38fc0f9.htm "DynamicState 属性 ")

[IsExpandOutput 属性](../html/37149140-100a-0f52-5c0f-d6f8c0ec208d.htm "IsExpandOutput 属性 ")

[IsPrint 属性](../html/2122f950-b596-0009-dbce-3a66e1b1f84c.htm "IsPrint 属性 ")

[IsTenExtend 属性](../html/b4cebc10-ab06-b5fc-513c-9fd5ef2fd3c4.htm "IsTenExtend 属性 ")

[SourceData 属性](../html/d478b7f6-ca65-6aa0-3156-581a1cd68a33.htm "SourceData 属性 ")

[Suttle 属性](../html/e1959b5f-89e8-17fe-d07b-4ceee83e0e7e.htm "Suttle 属性 ")

[Symbol 属性](../html/df9ef7d7-6a4e-04a0-6dbd-bfb6a4100ada.htm "Symbol 属性 ")

[Tare 属性](../html/3c34323e-f46e-bf12-6780-9e5ceb60384d.htm "Tare 属性 ")

[TareType 属性](../html/9eb382bf-639d-c4a2-d11c-6a50400cd482.htm "TareType 属性 ")

[Unit 属性](../html/4f500913-4b2e-da80-459d-e4513c0a5a3d.htm "Unit 属性 ")

[Weight 属性](../html/089a04ce-cb69-50fa-c9b4-508768d0c9dd.htm "Weight 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ToledoStandardDataDataValid 属性 |

数据是否有效

**命名空间：**
 [HslCommunication.Profinet.Toledo](3455a9c0-e5c9-7245-1dce-b1fee92e6846.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public bool DataValid { get; set; }
```

```
Public Property DataValid As Boolean
	Get
	Set
```

```
public:
property bool DataValid {
	bool get ();
	void set (bool value);
}
```

```
member DataValid : bool with get, set
```

#### 属性值

类型：Boolean

![](../icons/SectionExpanded.png)参见

#### 引用

[ToledoStandardData 类](13fdb4fb-3312-d193-7e77-a383ddffbf1d.htm)

[HslCommunication.Profinet.Toledo 命名空间](3455a9c0-e5c9-7245-1dce-b1fee92e6846.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)