# HslCommunication - HslCommunication.Secs.Types

> 分類頁數: 30



---
## HslCommunication.Secs.Types

[原文連結](http://api.hslcommunication.cn/html/a75c6382-21bf-e81e-9704-41aecc3cdc4b.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Secs.Types](../html/a75c6382-21bf-e81e-9704-41aecc3cdc4b.htm "HslCommunication.Secs.Types")

[ISecs 接口](../html/08134ad0-26a2-2a6b-52a7-0a3f930b5731.htm "ISecs 接口")

[OnlineData 类](../html/94067ba2-5c27-cc1a-c397-236b70c3fecb.htm "OnlineData 类")

[SecsItemType 枚举](../html/25a238c4-3ec8-ecf1-c580-83e39b5f9595.htm "SecsItemType 枚举")

[SecsMessage 类](../html/3594b6c1-13b3-f4d8-6932-cf74c11ea232.htm "SecsMessage 类")

[SecsMessageExtension 类](../html/a5e48ed5-8b1e-79b4-40b8-568c622983d4.htm "SecsMessageExtension 类")

[SecsValue 类](../html/06a37040-2c85-e52a-86bb-2679fa81cf40.htm "SecsValue 类")

[VariableName 类](../html/66b8d742-dcdd-e507-ab9a-87998843e5f1.htm "VariableName 类")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslCommunication.Secs.Types 命名空间 |

[缺少 "N:HslCommunication.Secs.Types" 的 <summary> 文档]

![](../icons/SectionExpanded.png)类

|  | 类 | 说明 |
| --- | --- | --- |
| 公共类 | [OnlineData](94067ba2-5c27-cc1a-c397-236b70c3fecb.htm) | 在线数据信息 |
| 公共类 | [SecsMessage](3594b6c1-13b3-f4d8-6932-cf74c11ea232.htm) | Secs的消息类对象 |
| 公共类 | [SecsMessageExtension](a5e48ed5-8b1e-79b4-40b8-568c622983d4.htm) | 扩展类 |
| 公共类代码示例 | [SecsValue](06a37040-2c85-e52a-86bb-2679fa81cf40.htm) | SECS数据的对象信息，可以用来表示层级及嵌套的数据内容，如果需要显示，只需要[ToString](bccbab8a-839a-4952-2c7d-d536313e648a.htm) 方法即可， 如果需要发送SECS设备，只需要 [ToSourceBytes](b6b6fdc2-085d-27f1-a16d-caa3a8ed88cf.htm)，并支持反序列化操作 [ParseFromSource(Byte, Encoding)](a6a154e2-8c1a-2a2e-c509-2f656d6fc77b.htm)，无论是XML元素还是byte[]类型。  SECS data object information, can be used to represent the hierarchy and nested data content, if you need to display, just need to [ToString](bccbab8a-839a-4952-2c7d-d536313e648a.htm) method can be. If you need to send SECS equipment, only need [ToSourceBytes](b6b6fdc2-085d-27f1-a16d-caa3a8ed88cf.htm), and support the deserialization operation [ParseFromSource(Byte, Encoding)](a6a154e2-8c1a-2a2e-c509-2f656d6fc77b.htm). Whether it's an XML element or byte[] type. |
| 公共类 | [VariableName](66b8d742-dcdd-e507-ab9a-87998843e5f1.htm) | 变量名称类对象 |

![](../icons/SectionExpanded.png)接口

|  | 接口 | 说明 |
| --- | --- | --- |
| 公共接口 | [ISecs](08134ad0-26a2-2a6b-52a7-0a3f930b5731.htm) | ISecs的接口信息，支持了将数据发送到对方，或是使用问答机制从设备获取数据  The interface information of ISecs supports sending data to the other party, or using the question and answer mechanism to obtain data from the device |

![](../icons/SectionExpanded.png)枚举

|  | 枚举 | 说明 |
| --- | --- | --- |
| 公共枚举 | [SecsItemType](25a238c4-3ec8-ecf1-c580-83e39b5f9595.htm) | 数据类型的定义 |

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ISecs 接口

[原文連結](http://api.hslcommunication.cn/html/08134ad0-26a2-2a6b-52a7-0a3f930b5731.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Secs.Types](../html/a75c6382-21bf-e81e-9704-41aecc3cdc4b.htm "HslCommunication.Secs.Types")

[ISecs 接口](../html/08134ad0-26a2-2a6b-52a7-0a3f930b5731.htm "ISecs 接口")

[ISecs 方法](../html/01acd2e8-e0a2-3b5c-01f9-72a7ffac2d75.htm "ISecs 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ISecs 接口 |

ISecs的接口信息，支持了将数据发送到对方，或是使用问答机制从设备获取数据  
The interface information of ISecs supports sending data to the other party, or using the question and answer mechanism to obtain data from the device

**命名空间：**
 [HslCommunication.Secs.Types](a75c6382-21bf-e81e-9704-41aecc3cdc4b.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public interface ISecs
```

```
Public Interface ISecs
```

```
public interface class ISecs
```

```
type ISecs =  interface end
```

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [ReadSecsMessage(Byte, Byte, SecsValue, Boolean)](8e9fd7f5-29ee-056b-b2d2-422519fc4889.htm) | 根据指定的功能码将数据报文发送给设备，并且等待从SECS设备返回Secs消息，本访问机制是问答模式的。  Send the data message to the device according to the specified function code, and wait for the Secs message to be returned from the SECS device. This access mechanism is in question-and-answer mode. |
| 公共方法 | [ReadSecsMessage(Byte, Byte, Byte, Boolean)](e2137630-3690-7c5e-5c5a-635104598421.htm) | 根据指定的功能码将数据报文发送给设备，并且等待从SECS设备返回Secs消息，本访问机制是问答模式的。  Send the data message to the device according to the specified function code, and wait for the Secs message to be returned from the SECS device. This access mechanism is in question-and-answer mode. |
| 公共方法 | [SendByCommand(Byte, Byte, SecsValue, Boolean)](7b9351e2-c311-d33b-6a3f-6abb5ac6b0dc.htm) | 将数据发送到设备方去，只是单纯的发送数据过去，并不等待设备的数据返回，返回是否发送成功，以及当前发送时候的 SystemBytes 的值。  Sending data to the device side simply sends the data to the past, and does not wait for the data from the device to return, and returns whether the transmission is successful. |
| 公共方法 | [SendByCommand(Byte, Byte, Byte, Boolean)](1806e328-bc88-85e6-1ed5-8380da8d0b8b.htm) | 将数据发送到设备方去，只是单纯的发送数据过去，并不等待设备的数据返回，返回是否发送成功，以及当前发送时候的 SystemBytes 的值。  Sending data to the device side simply sends the data to the past, and does not wait for the data from the device to return, and returns whether the transmission is successful. |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.Secs.Types 命名空间](a75c6382-21bf-e81e-9704-41aecc3cdc4b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ISecs 方法

[原文連結](http://api.hslcommunication.cn/html/01acd2e8-e0a2-3b5c-01f9-72a7ffac2d75.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Secs.Types](../html/a75c6382-21bf-e81e-9704-41aecc3cdc4b.htm "HslCommunication.Secs.Types")

[ISecs 接口](../html/08134ad0-26a2-2a6b-52a7-0a3f930b5731.htm "ISecs 接口")

[ISecs 方法](../html/01acd2e8-e0a2-3b5c-01f9-72a7ffac2d75.htm "ISecs 方法")

[ReadSecsMessage 方法](../html/ce300dcc-de34-3933-6b95-7647f2297b8d.htm "ReadSecsMessage 方法 ")

[SendByCommand 方法](../html/7c95127a-d935-5b7f-450c-84fe3ffff93c.htm "SendByCommand 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ISecs 方法 |

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [ReadSecsMessage(Byte, Byte, SecsValue, Boolean)](8e9fd7f5-29ee-056b-b2d2-422519fc4889.htm) | 根据指定的功能码将数据报文发送给设备，并且等待从SECS设备返回Secs消息，本访问机制是问答模式的。  Send the data message to the device according to the specified function code, and wait for the Secs message to be returned from the SECS device. This access mechanism is in question-and-answer mode. |
| 公共方法 | [ReadSecsMessage(Byte, Byte, Byte, Boolean)](e2137630-3690-7c5e-5c5a-635104598421.htm) | 根据指定的功能码将数据报文发送给设备，并且等待从SECS设备返回Secs消息，本访问机制是问答模式的。  Send the data message to the device according to the specified function code, and wait for the Secs message to be returned from the SECS device. This access mechanism is in question-and-answer mode. |
| 公共方法 | [SendByCommand(Byte, Byte, SecsValue, Boolean)](7b9351e2-c311-d33b-6a3f-6abb5ac6b0dc.htm) | 将数据发送到设备方去，只是单纯的发送数据过去，并不等待设备的数据返回，返回是否发送成功，以及当前发送时候的 SystemBytes 的值。  Sending data to the device side simply sends the data to the past, and does not wait for the data from the device to return, and returns whether the transmission is successful. |
| 公共方法 | [SendByCommand(Byte, Byte, Byte, Boolean)](1806e328-bc88-85e6-1ed5-8380da8d0b8b.htm) | 将数据发送到设备方去，只是单纯的发送数据过去，并不等待设备的数据返回，返回是否发送成功，以及当前发送时候的 SystemBytes 的值。  Sending data to the device side simply sends the data to the past, and does not wait for the data from the device to return, and returns whether the transmission is successful. |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[ISecs 接口](08134ad0-26a2-2a6b-52a7-0a3f930b5731.htm)

[HslCommunication.Secs.Types 命名空间](a75c6382-21bf-e81e-9704-41aecc3cdc4b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadSecsMessage 方法 

[原文連結](http://api.hslcommunication.cn/html/ce300dcc-de34-3933-6b95-7647f2297b8d.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Secs.Types](../html/a75c6382-21bf-e81e-9704-41aecc3cdc4b.htm "HslCommunication.Secs.Types")

[ISecs 接口](../html/08134ad0-26a2-2a6b-52a7-0a3f930b5731.htm "ISecs 接口")

[ISecs 方法](../html/01acd2e8-e0a2-3b5c-01f9-72a7ffac2d75.htm "ISecs 方法")

[ReadSecsMessage 方法](../html/ce300dcc-de34-3933-6b95-7647f2297b8d.htm "ReadSecsMessage 方法 ")

[ReadSecsMessage 方法 (Byte, Byte, SecsValue, Boolean)](../html/8e9fd7f5-29ee-056b-b2d2-422519fc4889.htm "ReadSecsMessage 方法 (Byte, Byte, SecsValue, Boolean)")

[ReadSecsMessage 方法 (Byte, Byte, Byte[], Boolean)](../html/e2137630-3690-7c5e-5c5a-635104598421.htm "ReadSecsMessage 方法 (Byte, Byte, Byte[], Boolean)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ISecsReadSecsMessage 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [ReadSecsMessage(Byte, Byte, SecsValue, Boolean)](8e9fd7f5-29ee-056b-b2d2-422519fc4889.htm) | 根据指定的功能码将数据报文发送给设备，并且等待从SECS设备返回Secs消息，本访问机制是问答模式的。  Send the data message to the device according to the specified function code, and wait for the Secs message to be returned from the SECS device. This access mechanism is in question-and-answer mode. |
| 公共方法 | [ReadSecsMessage(Byte, Byte, Byte, Boolean)](e2137630-3690-7c5e-5c5a-635104598421.htm) | 根据指定的功能码将数据报文发送给设备，并且等待从SECS设备返回Secs消息，本访问机制是问答模式的。  Send the data message to the device according to the specified function code, and wait for the Secs message to be returned from the SECS device. This access mechanism is in question-and-answer mode. |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[ISecs 接口](08134ad0-26a2-2a6b-52a7-0a3f930b5731.htm)

[HslCommunication.Secs.Types 命名空间](a75c6382-21bf-e81e-9704-41aecc3cdc4b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadSecsMessage 方法 (Byte, Byte, SecsValue, Boolean)

[原文連結](http://api.hslcommunication.cn/html/8e9fd7f5-29ee-056b-b2d2-422519fc4889.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Secs.Types](../html/a75c6382-21bf-e81e-9704-41aecc3cdc4b.htm "HslCommunication.Secs.Types")

[ISecs 接口](../html/08134ad0-26a2-2a6b-52a7-0a3f930b5731.htm "ISecs 接口")

[ISecs 方法](../html/01acd2e8-e0a2-3b5c-01f9-72a7ffac2d75.htm "ISecs 方法")

[ReadSecsMessage 方法](../html/ce300dcc-de34-3933-6b95-7647f2297b8d.htm "ReadSecsMessage 方法 ")

[ReadSecsMessage 方法 (Byte, Byte, SecsValue, Boolean)](../html/8e9fd7f5-29ee-056b-b2d2-422519fc4889.htm "ReadSecsMessage 方法 (Byte, Byte, SecsValue, Boolean)")

[ReadSecsMessage 方法 (Byte, Byte, Byte[], Boolean)](../html/e2137630-3690-7c5e-5c5a-635104598421.htm "ReadSecsMessage 方法 (Byte, Byte, Byte[], Boolean)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ISecsReadSecsMessage 方法 (Byte, Byte, SecsValue, Boolean) |

根据指定的功能码将数据报文发送给设备，并且等待从SECS设备返回Secs消息，本访问机制是问答模式的。  
Send the data message to the device according to the specified function code, and wait for the Secs message to be returned from the SECS device. This access mechanism is in question-and-answer mode.

**命名空间：**
 [HslCommunication.Secs.Types](a75c6382-21bf-e81e-9704-41aecc3cdc4b.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
OperateResult<SecsMessage> ReadSecsMessage(
	byte stream,
	byte function,
	SecsValue data,
	bool back
)
```

```
Function ReadSecsMessage ( 
	stream As Byte,
	function As Byte,
	data As SecsValue,
	back As Boolean
) As OperateResult(Of SecsMessage)
```

```
OperateResult<SecsMessage^>^ ReadSecsMessage(
	unsigned char stream, 
	unsigned char function, 
	SecsValue^ data, 
	bool back
)
```

```
abstract ReadSecsMessage : 
        stream : byte * 
        function : byte * 
        data : SecsValue * 
        back : bool -> OperateResult<SecsMessage> 
```

#### 参数

stream
:   类型：SystemByte  
    功能码1

function
:   类型：SystemByte  
    功能码2

data
:   类型：[HslCommunication.Secs.TypesSecsValue](06a37040-2c85-e52a-86bb-2679fa81cf40.htm)  
    Secs格式的对象信息

back
:   类型：SystemBoolean  
    是否必须返回，此标记仅仅是secs报文的是否返回标记，不表示问答模式

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)[SecsMessage](3594b6c1-13b3-f4d8-6932-cf74c11ea232.htm)  
返回SECS消息结果对象

![](../icons/SectionExpanded.png)参见

#### 引用

[ISecs 接口](08134ad0-26a2-2a6b-52a7-0a3f930b5731.htm)

[ReadSecsMessage 重载](ce300dcc-de34-3933-6b95-7647f2297b8d.htm)

[HslCommunication.Secs.Types 命名空间](a75c6382-21bf-e81e-9704-41aecc3cdc4b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadSecsMessage 方法 (Byte, Byte, Byte[], Boolean)

[原文連結](http://api.hslcommunication.cn/html/e2137630-3690-7c5e-5c5a-635104598421.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Secs.Types](../html/a75c6382-21bf-e81e-9704-41aecc3cdc4b.htm "HslCommunication.Secs.Types")

[ISecs 接口](../html/08134ad0-26a2-2a6b-52a7-0a3f930b5731.htm "ISecs 接口")

[ISecs 方法](../html/01acd2e8-e0a2-3b5c-01f9-72a7ffac2d75.htm "ISecs 方法")

[ReadSecsMessage 方法](../html/ce300dcc-de34-3933-6b95-7647f2297b8d.htm "ReadSecsMessage 方法 ")

[ReadSecsMessage 方法 (Byte, Byte, SecsValue, Boolean)](../html/8e9fd7f5-29ee-056b-b2d2-422519fc4889.htm "ReadSecsMessage 方法 (Byte, Byte, SecsValue, Boolean)")

[ReadSecsMessage 方法 (Byte, Byte, Byte[], Boolean)](../html/e2137630-3690-7c5e-5c5a-635104598421.htm "ReadSecsMessage 方法 (Byte, Byte, Byte[], Boolean)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ISecsReadSecsMessage 方法 (Byte, Byte, Byte, Boolean) |

根据指定的功能码将数据报文发送给设备，并且等待从SECS设备返回Secs消息，本访问机制是问答模式的。  
Send the data message to the device according to the specified function code, and wait for the Secs message to be returned from the SECS device. This access mechanism is in question-and-answer mode.

**命名空间：**
 [HslCommunication.Secs.Types](a75c6382-21bf-e81e-9704-41aecc3cdc4b.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
OperateResult<SecsMessage> ReadSecsMessage(
	byte stream,
	byte function,
	byte[] data,
	bool back
)
```

```
Function ReadSecsMessage ( 
	stream As Byte,
	function As Byte,
	data As Byte(),
	back As Boolean
) As OperateResult(Of SecsMessage)
```

```
OperateResult<SecsMessage^>^ ReadSecsMessage(
	unsigned char stream, 
	unsigned char function, 
	array<unsigned char>^ data, 
	bool back
)
```

```
abstract ReadSecsMessage : 
        stream : byte * 
        function : byte * 
        data : byte[] * 
        back : bool -> OperateResult<SecsMessage> 
```

#### 参数

stream
:   类型：SystemByte  
    功能码1

function
:   类型：SystemByte  
    功能码2

data
:   类型：SystemByte  
    原始的字节数据

back
:   类型：SystemBoolean  
    是否必须返回，此标记仅仅是secs报文的是否返回标记，不表示问答模式

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)[SecsMessage](3594b6c1-13b3-f4d8-6932-cf74c11ea232.htm)  
返回SECS消息结果对象

![](../icons/SectionExpanded.png)参见

#### 引用

[ISecs 接口](08134ad0-26a2-2a6b-52a7-0a3f930b5731.htm)

[ReadSecsMessage 重载](ce300dcc-de34-3933-6b95-7647f2297b8d.htm)

[HslCommunication.Secs.Types 命名空间](a75c6382-21bf-e81e-9704-41aecc3cdc4b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## SendByCommand 方法 

[原文連結](http://api.hslcommunication.cn/html/7c95127a-d935-5b7f-450c-84fe3ffff93c.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Secs.Types](../html/a75c6382-21bf-e81e-9704-41aecc3cdc4b.htm "HslCommunication.Secs.Types")

[ISecs 接口](../html/08134ad0-26a2-2a6b-52a7-0a3f930b5731.htm "ISecs 接口")

[ISecs 方法](../html/01acd2e8-e0a2-3b5c-01f9-72a7ffac2d75.htm "ISecs 方法")

[SendByCommand 方法](../html/7c95127a-d935-5b7f-450c-84fe3ffff93c.htm "SendByCommand 方法 ")

[SendByCommand 方法 (Byte, Byte, SecsValue, Boolean)](../html/7b9351e2-c311-d33b-6a3f-6abb5ac6b0dc.htm "SendByCommand 方法 (Byte, Byte, SecsValue, Boolean)")

[SendByCommand 方法 (Byte, Byte, Byte[], Boolean)](../html/1806e328-bc88-85e6-1ed5-8380da8d0b8b.htm "SendByCommand 方法 (Byte, Byte, Byte[], Boolean)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ISecsSendByCommand 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [SendByCommand(Byte, Byte, SecsValue, Boolean)](7b9351e2-c311-d33b-6a3f-6abb5ac6b0dc.htm) | 将数据发送到设备方去，只是单纯的发送数据过去，并不等待设备的数据返回，返回是否发送成功，以及当前发送时候的 SystemBytes 的值。  Sending data to the device side simply sends the data to the past, and does not wait for the data from the device to return, and returns whether the transmission is successful. |
| 公共方法 | [SendByCommand(Byte, Byte, Byte, Boolean)](1806e328-bc88-85e6-1ed5-8380da8d0b8b.htm) | 将数据发送到设备方去，只是单纯的发送数据过去，并不等待设备的数据返回，返回是否发送成功，以及当前发送时候的 SystemBytes 的值。  Sending data to the device side simply sends the data to the past, and does not wait for the data from the device to return, and returns whether the transmission is successful. |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[ISecs 接口](08134ad0-26a2-2a6b-52a7-0a3f930b5731.htm)

[HslCommunication.Secs.Types 命名空间](a75c6382-21bf-e81e-9704-41aecc3cdc4b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## SendByCommand 方法 (Byte, Byte, SecsValue, Boolean)

[原文連結](http://api.hslcommunication.cn/html/7b9351e2-c311-d33b-6a3f-6abb5ac6b0dc.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Secs.Types](../html/a75c6382-21bf-e81e-9704-41aecc3cdc4b.htm "HslCommunication.Secs.Types")

[ISecs 接口](../html/08134ad0-26a2-2a6b-52a7-0a3f930b5731.htm "ISecs 接口")

[ISecs 方法](../html/01acd2e8-e0a2-3b5c-01f9-72a7ffac2d75.htm "ISecs 方法")

[SendByCommand 方法](../html/7c95127a-d935-5b7f-450c-84fe3ffff93c.htm "SendByCommand 方法 ")

[SendByCommand 方法 (Byte, Byte, SecsValue, Boolean)](../html/7b9351e2-c311-d33b-6a3f-6abb5ac6b0dc.htm "SendByCommand 方法 (Byte, Byte, SecsValue, Boolean)")

[SendByCommand 方法 (Byte, Byte, Byte[], Boolean)](../html/1806e328-bc88-85e6-1ed5-8380da8d0b8b.htm "SendByCommand 方法 (Byte, Byte, Byte[], Boolean)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ISecsSendByCommand 方法 (Byte, Byte, SecsValue, Boolean) |

将数据发送到设备方去，只是单纯的发送数据过去，并不等待设备的数据返回，返回是否发送成功，以及当前发送时候的 SystemBytes 的值。  
Sending data to the device side simply sends the data to the past, and does not wait for the data from the device to return, and returns whether the transmission is successful.

**命名空间：**
 [HslCommunication.Secs.Types](a75c6382-21bf-e81e-9704-41aecc3cdc4b.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
OperateResult<uint> SendByCommand(
	byte stream,
	byte function,
	SecsValue data,
	bool back
)
```

```
Function SendByCommand ( 
	stream As Byte,
	function As Byte,
	data As SecsValue,
	back As Boolean
) As OperateResult(Of UInteger)
```

```
OperateResult<unsigned int>^ SendByCommand(
	unsigned char stream, 
	unsigned char function, 
	SecsValue^ data, 
	bool back
)
```

```
abstract SendByCommand : 
        stream : byte * 
        function : byte * 
        data : SecsValue * 
        back : bool -> OperateResult<uint32> 
```

#### 参数

stream
:   类型：SystemByte  
    功能码1

function
:   类型：SystemByte  
    功能码2

data
:   类型：[HslCommunication.Secs.TypesSecsValue](06a37040-2c85-e52a-86bb-2679fa81cf40.htm)  
    Secs格式的对象信息

back
:   类型：SystemBoolean  
    是否必须返回，此标记仅仅是secs报文的是否返回标记，不表示问答模式

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)UInt32  
是否发送成功的结果对象

![](../icons/SectionExpanded.png)参见

#### 引用

[ISecs 接口](08134ad0-26a2-2a6b-52a7-0a3f930b5731.htm)

[SendByCommand 重载](7c95127a-d935-5b7f-450c-84fe3ffff93c.htm)

[HslCommunication.Secs.Types 命名空间](a75c6382-21bf-e81e-9704-41aecc3cdc4b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## SendByCommand 方法 (Byte, Byte, Byte[], Boolean)

[原文連結](http://api.hslcommunication.cn/html/1806e328-bc88-85e6-1ed5-8380da8d0b8b.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Secs.Types](../html/a75c6382-21bf-e81e-9704-41aecc3cdc4b.htm "HslCommunication.Secs.Types")

[ISecs 接口](../html/08134ad0-26a2-2a6b-52a7-0a3f930b5731.htm "ISecs 接口")

[ISecs 方法](../html/01acd2e8-e0a2-3b5c-01f9-72a7ffac2d75.htm "ISecs 方法")

[SendByCommand 方法](../html/7c95127a-d935-5b7f-450c-84fe3ffff93c.htm "SendByCommand 方法 ")

[SendByCommand 方法 (Byte, Byte, SecsValue, Boolean)](../html/7b9351e2-c311-d33b-6a3f-6abb5ac6b0dc.htm "SendByCommand 方法 (Byte, Byte, SecsValue, Boolean)")

[SendByCommand 方法 (Byte, Byte, Byte[], Boolean)](../html/1806e328-bc88-85e6-1ed5-8380da8d0b8b.htm "SendByCommand 方法 (Byte, Byte, Byte[], Boolean)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ISecsSendByCommand 方法 (Byte, Byte, Byte, Boolean) |

将数据发送到设备方去，只是单纯的发送数据过去，并不等待设备的数据返回，返回是否发送成功，以及当前发送时候的 SystemBytes 的值。  
Sending data to the device side simply sends the data to the past, and does not wait for the data from the device to return, and returns whether the transmission is successful.

**命名空间：**
 [HslCommunication.Secs.Types](a75c6382-21bf-e81e-9704-41aecc3cdc4b.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
OperateResult<uint> SendByCommand(
	byte stream,
	byte function,
	byte[] data,
	bool back
)
```

```
Function SendByCommand ( 
	stream As Byte,
	function As Byte,
	data As Byte(),
	back As Boolean
) As OperateResult(Of UInteger)
```

```
OperateResult<unsigned int>^ SendByCommand(
	unsigned char stream, 
	unsigned char function, 
	array<unsigned char>^ data, 
	bool back
)
```

```
abstract SendByCommand : 
        stream : byte * 
        function : byte * 
        data : byte[] * 
        back : bool -> OperateResult<uint32> 
```

#### 参数

stream
:   类型：SystemByte  
    功能码1

function
:   类型：SystemByte  
    功能码2

data
:   类型：SystemByte  
    原始的字节数据

back
:   类型：SystemBoolean  
    是否必须返回，此标记仅仅是secs报文的是否返回标记，不表示问答模式

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)UInt32  
是否发送成功的结果对象

![](../icons/SectionExpanded.png)参见

#### 引用

[ISecs 接口](08134ad0-26a2-2a6b-52a7-0a3f930b5731.htm)

[SendByCommand 重载](7c95127a-d935-5b7f-450c-84fe3ffff93c.htm)

[HslCommunication.Secs.Types 命名空间](a75c6382-21bf-e81e-9704-41aecc3cdc4b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## OnlineData 类

[原文連結](http://api.hslcommunication.cn/html/94067ba2-5c27-cc1a-c397-236b70c3fecb.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Secs.Types](../html/a75c6382-21bf-e81e-9704-41aecc3cdc4b.htm "HslCommunication.Secs.Types")

[OnlineData 类](../html/94067ba2-5c27-cc1a-c397-236b70c3fecb.htm "OnlineData 类")

[OnlineData 构造函数](../html/1768000b-0209-5825-44f3-c1cf1986d959.htm "OnlineData 构造函数 ")

[OnlineData 属性](../html/cd6cd38a-a353-f42c-1ce2-f947a700e0ff.htm "OnlineData 属性")

[OnlineData 方法](../html/85799c14-8a95-aea9-a659-4c29b63b394b.htm "OnlineData 方法")

[OnlineData 类型转换](../html/016a1000-c05c-5e5b-27d3-0b4d5663b191.htm "OnlineData 类型转换")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| OnlineData 类 |

在线数据信息

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  HslCommunication.Secs.TypesOnlineData

**命名空间：**
 [HslCommunication.Secs.Types](a75c6382-21bf-e81e-9704-41aecc3cdc4b.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class OnlineData
```

```
Public Class OnlineData
```

```
public ref class OnlineData
```

```
type OnlineData =  class end
```

OnlineData 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [OnlineData](536e9996-b1ff-ecc4-d339-3c3779b0f0ff.htm) | 实例化一个默认的对象 |
| 公共方法 | [OnlineData(String, String)](8bbb0c59-5f91-a1db-80cc-b34d30b8dc33.htm) | 指定类型及其版本号来实例化一个对象 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [ModelType](eed8f154-6b4c-a4f1-bb67-b4ad29035915.htm) | equipment model type |
| 公共属性 | [SoftVersion](f8089fb4-e76f-577b-7a1b-154d8c22671f.htm) | software revision |

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

![](../icons/SectionExpanded.png)运算符

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共运算符静态成员 | [(OnlineData to SecsValue)](49c97311-9a30-4b67-8ce1-8b13aa08a21f.htm) | 也可以赋值给[SecsValue](06a37040-2c85-e52a-86bb-2679fa81cf40.htm) 数据 |
| 公共运算符静态成员 | [(SecsValue to OnlineData)](e5288bf0-6ea7-879b-60cc-b7dfeba2e5e4.htm) | 赋值操作，可以直接赋值 OnlineData 数据 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.Secs.Types 命名空间](a75c6382-21bf-e81e-9704-41aecc3cdc4b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## OnlineData 构造函数 

[原文連結](http://api.hslcommunication.cn/html/1768000b-0209-5825-44f3-c1cf1986d959.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Secs.Types](../html/a75c6382-21bf-e81e-9704-41aecc3cdc4b.htm "HslCommunication.Secs.Types")

[OnlineData 类](../html/94067ba2-5c27-cc1a-c397-236b70c3fecb.htm "OnlineData 类")

[OnlineData 构造函数](../html/1768000b-0209-5825-44f3-c1cf1986d959.htm "OnlineData 构造函数 ")

[OnlineData 构造函数](../html/536e9996-b1ff-ecc4-d339-3c3779b0f0ff.htm "OnlineData 构造函数 ")

[OnlineData 构造函数 (String, String)](../html/8bbb0c59-5f91-a1db-80cc-b34d30b8dc33.htm "OnlineData 构造函数 (String, String)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| OnlineData 构造函数 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [OnlineData](536e9996-b1ff-ecc4-d339-3c3779b0f0ff.htm) | 实例化一个默认的对象 |
| 公共方法 | [OnlineData(String, String)](8bbb0c59-5f91-a1db-80cc-b34d30b8dc33.htm) | 指定类型及其版本号来实例化一个对象 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[OnlineData 类](94067ba2-5c27-cc1a-c397-236b70c3fecb.htm)

[HslCommunication.Secs.Types 命名空间](a75c6382-21bf-e81e-9704-41aecc3cdc4b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## OnlineData 构造函数 

[原文連結](http://api.hslcommunication.cn/html/536e9996-b1ff-ecc4-d339-3c3779b0f0ff.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Secs.Types](../html/a75c6382-21bf-e81e-9704-41aecc3cdc4b.htm "HslCommunication.Secs.Types")

[OnlineData 类](../html/94067ba2-5c27-cc1a-c397-236b70c3fecb.htm "OnlineData 类")

[OnlineData 构造函数](../html/1768000b-0209-5825-44f3-c1cf1986d959.htm "OnlineData 构造函数 ")

[OnlineData 构造函数](../html/536e9996-b1ff-ecc4-d339-3c3779b0f0ff.htm "OnlineData 构造函数 ")

[OnlineData 构造函数 (String, String)](../html/8bbb0c59-5f91-a1db-80cc-b34d30b8dc33.htm "OnlineData 构造函数 (String, String)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| OnlineData 构造函数 |

实例化一个默认的对象

**命名空间：**
 [HslCommunication.Secs.Types](a75c6382-21bf-e81e-9704-41aecc3cdc4b.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public OnlineData()
```

```
Public Sub New
```

```
public:
OnlineData()
```

```
new : unit -> OnlineData
```

![](../icons/SectionExpanded.png)参见

#### 引用

[OnlineData 类](94067ba2-5c27-cc1a-c397-236b70c3fecb.htm)

[OnlineData 重载](1768000b-0209-5825-44f3-c1cf1986d959.htm)

[HslCommunication.Secs.Types 命名空间](a75c6382-21bf-e81e-9704-41aecc3cdc4b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## OnlineData 构造函数 (String, String)

[原文連結](http://api.hslcommunication.cn/html/8bbb0c59-5f91-a1db-80cc-b34d30b8dc33.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Secs.Types](../html/a75c6382-21bf-e81e-9704-41aecc3cdc4b.htm "HslCommunication.Secs.Types")

[OnlineData 类](../html/94067ba2-5c27-cc1a-c397-236b70c3fecb.htm "OnlineData 类")

[OnlineData 构造函数](../html/1768000b-0209-5825-44f3-c1cf1986d959.htm "OnlineData 构造函数 ")

[OnlineData 构造函数](../html/536e9996-b1ff-ecc4-d339-3c3779b0f0ff.htm "OnlineData 构造函数 ")

[OnlineData 构造函数 (String, String)](../html/8bbb0c59-5f91-a1db-80cc-b34d30b8dc33.htm "OnlineData 构造函数 (String, String)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| OnlineData 构造函数 (String, String) |

指定类型及其版本号来实例化一个对象

**命名空间：**
 [HslCommunication.Secs.Types](a75c6382-21bf-e81e-9704-41aecc3cdc4b.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public OnlineData(
	string model,
	string version
)
```

```
Public Sub New ( 
	model As String,
	version As String
)
```

```
public:
OnlineData(
	String^ model, 
	String^ version
)
```

```
new : 
        model : string * 
        version : string -> OnlineData
```

#### 参数

model
:   类型：SystemString  
    类型信息

version
:   类型：SystemString  
    版本号

![](../icons/SectionExpanded.png)参见

#### 引用

[OnlineData 类](94067ba2-5c27-cc1a-c397-236b70c3fecb.htm)

[OnlineData 重载](1768000b-0209-5825-44f3-c1cf1986d959.htm)

[HslCommunication.Secs.Types 命名空间](a75c6382-21bf-e81e-9704-41aecc3cdc4b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## OnlineData 属性

[原文連結](http://api.hslcommunication.cn/html/cd6cd38a-a353-f42c-1ce2-f947a700e0ff.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Secs.Types](../html/a75c6382-21bf-e81e-9704-41aecc3cdc4b.htm "HslCommunication.Secs.Types")

[OnlineData 类](../html/94067ba2-5c27-cc1a-c397-236b70c3fecb.htm "OnlineData 类")

[OnlineData 属性](../html/cd6cd38a-a353-f42c-1ce2-f947a700e0ff.htm "OnlineData 属性")

[ModelType 属性](../html/eed8f154-6b4c-a4f1-bb67-b4ad29035915.htm "ModelType 属性 ")

[SoftVersion 属性](../html/f8089fb4-e76f-577b-7a1b-154d8c22671f.htm "SoftVersion 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| OnlineData 属性 |

[OnlineData](94067ba2-5c27-cc1a-c397-236b70c3fecb.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [ModelType](eed8f154-6b4c-a4f1-bb67-b4ad29035915.htm) | equipment model type |
| 公共属性 | [SoftVersion](f8089fb4-e76f-577b-7a1b-154d8c22671f.htm) | software revision |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[OnlineData 类](94067ba2-5c27-cc1a-c397-236b70c3fecb.htm)

[HslCommunication.Secs.Types 命名空间](a75c6382-21bf-e81e-9704-41aecc3cdc4b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ModelType 属性 

[原文連結](http://api.hslcommunication.cn/html/eed8f154-6b4c-a4f1-bb67-b4ad29035915.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Secs.Types](../html/a75c6382-21bf-e81e-9704-41aecc3cdc4b.htm "HslCommunication.Secs.Types")

[OnlineData 类](../html/94067ba2-5c27-cc1a-c397-236b70c3fecb.htm "OnlineData 类")

[OnlineData 属性](../html/cd6cd38a-a353-f42c-1ce2-f947a700e0ff.htm "OnlineData 属性")

[ModelType 属性](../html/eed8f154-6b4c-a4f1-bb67-b4ad29035915.htm "ModelType 属性 ")

[SoftVersion 属性](../html/f8089fb4-e76f-577b-7a1b-154d8c22671f.htm "SoftVersion 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| OnlineDataModelType 属性 |

equipment model type

**命名空间：**
 [HslCommunication.Secs.Types](a75c6382-21bf-e81e-9704-41aecc3cdc4b.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public string ModelType { get; set; }
```

```
Public Property ModelType As String
	Get
	Set
```

```
public:
property String^ ModelType {
	String^ get ();
	void set (String^ value);
}
```

```
member ModelType : string with get, set
```

#### 属性值

类型：String

![](../icons/SectionExpanded.png)参见

#### 引用

[OnlineData 类](94067ba2-5c27-cc1a-c397-236b70c3fecb.htm)

[HslCommunication.Secs.Types 命名空间](a75c6382-21bf-e81e-9704-41aecc3cdc4b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## SoftVersion 属性 

[原文連結](http://api.hslcommunication.cn/html/f8089fb4-e76f-577b-7a1b-154d8c22671f.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Secs.Types](../html/a75c6382-21bf-e81e-9704-41aecc3cdc4b.htm "HslCommunication.Secs.Types")

[OnlineData 类](../html/94067ba2-5c27-cc1a-c397-236b70c3fecb.htm "OnlineData 类")

[OnlineData 属性](../html/cd6cd38a-a353-f42c-1ce2-f947a700e0ff.htm "OnlineData 属性")

[ModelType 属性](../html/eed8f154-6b4c-a4f1-bb67-b4ad29035915.htm "ModelType 属性 ")

[SoftVersion 属性](../html/f8089fb4-e76f-577b-7a1b-154d8c22671f.htm "SoftVersion 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| OnlineDataSoftVersion 属性 |

software revision

**命名空间：**
 [HslCommunication.Secs.Types](a75c6382-21bf-e81e-9704-41aecc3cdc4b.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public string SoftVersion { get; set; }
```

```
Public Property SoftVersion As String
	Get
	Set
```

```
public:
property String^ SoftVersion {
	String^ get ();
	void set (String^ value);
}
```

```
member SoftVersion : string with get, set
```

#### 属性值

类型：String

![](../icons/SectionExpanded.png)参见

#### 引用

[OnlineData 类](94067ba2-5c27-cc1a-c397-236b70c3fecb.htm)

[HslCommunication.Secs.Types 命名空间](a75c6382-21bf-e81e-9704-41aecc3cdc4b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## OnlineData 方法

[原文連結](http://api.hslcommunication.cn/html/85799c14-8a95-aea9-a659-4c29b63b394b.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Secs.Types](../html/a75c6382-21bf-e81e-9704-41aecc3cdc4b.htm "HslCommunication.Secs.Types")

[OnlineData 类](../html/94067ba2-5c27-cc1a-c397-236b70c3fecb.htm "OnlineData 类")

[OnlineData 构造函数](../html/1768000b-0209-5825-44f3-c1cf1986d959.htm "OnlineData 构造函数 ")

[OnlineData 属性](../html/cd6cd38a-a353-f42c-1ce2-f947a700e0ff.htm "OnlineData 属性")

[OnlineData 方法](../html/85799c14-8a95-aea9-a659-4c29b63b394b.htm "OnlineData 方法")

[OnlineData 类型转换](../html/016a1000-c05c-5e5b-27d3-0b4d5663b191.htm "OnlineData 类型转换")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| OnlineData 方法 |

[OnlineData](94067ba2-5c27-cc1a-c397-236b70c3fecb.htm) 类型公开以下成员。

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

[OnlineData 类](94067ba2-5c27-cc1a-c397-236b70c3fecb.htm)

[HslCommunication.Secs.Types 命名空间](a75c6382-21bf-e81e-9704-41aecc3cdc4b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## OnlineData 类型转换

[原文連結](http://api.hslcommunication.cn/html/016a1000-c05c-5e5b-27d3-0b4d5663b191.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Secs.Types](../html/a75c6382-21bf-e81e-9704-41aecc3cdc4b.htm "HslCommunication.Secs.Types")

[OnlineData 类](../html/94067ba2-5c27-cc1a-c397-236b70c3fecb.htm "OnlineData 类")

[OnlineData 类型转换](../html/016a1000-c05c-5e5b-27d3-0b4d5663b191.htm "OnlineData 类型转换")

[Implicit 转换运算符](../html/422092b1-d478-8393-60d1-0d47b4d40a44.htm "Implicit 转换运算符")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| OnlineData 类型转换 |

![](../icons/SectionExpanded.png)运算符

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共运算符静态成员 | [(OnlineData to SecsValue)](49c97311-9a30-4b67-8ce1-8b13aa08a21f.htm) | 也可以赋值给[SecsValue](06a37040-2c85-e52a-86bb-2679fa81cf40.htm) 数据 |
| 公共运算符静态成员 | [(SecsValue to OnlineData)](e5288bf0-6ea7-879b-60cc-b7dfeba2e5e4.htm) | 赋值操作，可以直接赋值 [OnlineData](94067ba2-5c27-cc1a-c397-236b70c3fecb.htm) 数据 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[OnlineData 类](94067ba2-5c27-cc1a-c397-236b70c3fecb.htm)

[HslCommunication.Secs.Types 命名空间](a75c6382-21bf-e81e-9704-41aecc3cdc4b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Implicit 转换运算符

[原文連結](http://api.hslcommunication.cn/html/422092b1-d478-8393-60d1-0d47b4d40a44.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Secs.Types](../html/a75c6382-21bf-e81e-9704-41aecc3cdc4b.htm "HslCommunication.Secs.Types")

[OnlineData 类](../html/94067ba2-5c27-cc1a-c397-236b70c3fecb.htm "OnlineData 类")

[OnlineData 类型转换](../html/016a1000-c05c-5e5b-27d3-0b4d5663b191.htm "OnlineData 类型转换")

[Implicit 转换运算符](../html/422092b1-d478-8393-60d1-0d47b4d40a44.htm "Implicit 转换运算符")

[Implicit 转换 (OnlineData to SecsValue)](../html/49c97311-9a30-4b67-8ce1-8b13aa08a21f.htm "Implicit 转换 (OnlineData to SecsValue)")

[Implicit 转换 (SecsValue to OnlineData)](../html/e5288bf0-6ea7-879b-60cc-b7dfeba2e5e4.htm "Implicit 转换 (SecsValue to OnlineData)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| OnlineData  转换运算符 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共运算符静态成员 | [(OnlineData to SecsValue)](49c97311-9a30-4b67-8ce1-8b13aa08a21f.htm) | 也可以赋值给[SecsValue](06a37040-2c85-e52a-86bb-2679fa81cf40.htm) 数据 |
| 公共运算符静态成员 | [(SecsValue to OnlineData)](e5288bf0-6ea7-879b-60cc-b7dfeba2e5e4.htm) | 赋值操作，可以直接赋值 [OnlineData](94067ba2-5c27-cc1a-c397-236b70c3fecb.htm) 数据 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[OnlineData 类](94067ba2-5c27-cc1a-c397-236b70c3fecb.htm)

[HslCommunication.Secs.Types 命名空间](a75c6382-21bf-e81e-9704-41aecc3cdc4b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Implicit 转换 (OnlineData to SecsValue)

[原文連結](http://api.hslcommunication.cn/html/49c97311-9a30-4b67-8ce1-8b13aa08a21f.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Secs.Types](../html/a75c6382-21bf-e81e-9704-41aecc3cdc4b.htm "HslCommunication.Secs.Types")

[OnlineData 类](../html/94067ba2-5c27-cc1a-c397-236b70c3fecb.htm "OnlineData 类")

[OnlineData 类型转换](../html/016a1000-c05c-5e5b-27d3-0b4d5663b191.htm "OnlineData 类型转换")

[Implicit 转换运算符](../html/422092b1-d478-8393-60d1-0d47b4d40a44.htm "Implicit 转换运算符")

[Implicit 转换 (OnlineData to SecsValue)](../html/49c97311-9a30-4b67-8ce1-8b13aa08a21f.htm "Implicit 转换 (OnlineData to SecsValue)")

[Implicit 转换 (SecsValue to OnlineData)](../html/e5288bf0-6ea7-879b-60cc-b7dfeba2e5e4.htm "Implicit 转换 (SecsValue to OnlineData)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| OnlineData  转换 (OnlineData to SecsValue) |

也可以赋值给[SecsValue](06a37040-2c85-e52a-86bb-2679fa81cf40.htm) 数据

**命名空间：**
 [HslCommunication.Secs.Types](a75c6382-21bf-e81e-9704-41aecc3cdc4b.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static implicit operator SecsValue (
	OnlineData value
)
```

```
Public Shared Widening Operator CType ( 
	value As OnlineData
) As SecsValue
```

```
static implicit operator SecsValue^ (
	OnlineData^ value
)
```

```
F# does not support the declaration of new casting operators.
```

#### 参数

value
:   类型：[HslCommunication.Secs.TypesOnlineData](94067ba2-5c27-cc1a-c397-236b70c3fecb.htm)  
    [SecsValue](06a37040-2c85-e52a-86bb-2679fa81cf40.htm) 对象

#### 返回值

类型：[SecsValue](06a37040-2c85-e52a-86bb-2679fa81cf40.htm)  
等值的消息对象

![](../icons/SectionExpanded.png)参见

#### 引用

[OnlineData 类](94067ba2-5c27-cc1a-c397-236b70c3fecb.htm)

[重载](422092b1-d478-8393-60d1-0d47b4d40a44.htm)

[HslCommunication.Secs.Types 命名空间](a75c6382-21bf-e81e-9704-41aecc3cdc4b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Implicit 转换 (SecsValue to OnlineData)

[原文連結](http://api.hslcommunication.cn/html/e5288bf0-6ea7-879b-60cc-b7dfeba2e5e4.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Secs.Types](../html/a75c6382-21bf-e81e-9704-41aecc3cdc4b.htm "HslCommunication.Secs.Types")

[OnlineData 类](../html/94067ba2-5c27-cc1a-c397-236b70c3fecb.htm "OnlineData 类")

[OnlineData 类型转换](../html/016a1000-c05c-5e5b-27d3-0b4d5663b191.htm "OnlineData 类型转换")

[Implicit 转换运算符](../html/422092b1-d478-8393-60d1-0d47b4d40a44.htm "Implicit 转换运算符")

[Implicit 转换 (OnlineData to SecsValue)](../html/49c97311-9a30-4b67-8ce1-8b13aa08a21f.htm "Implicit 转换 (OnlineData to SecsValue)")

[Implicit 转换 (SecsValue to OnlineData)](../html/e5288bf0-6ea7-879b-60cc-b7dfeba2e5e4.htm "Implicit 转换 (SecsValue to OnlineData)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| OnlineData  转换 (SecsValue to OnlineData) |

赋值操作，可以直接赋值 [OnlineData](94067ba2-5c27-cc1a-c397-236b70c3fecb.htm) 数据

**命名空间：**
 [HslCommunication.Secs.Types](a75c6382-21bf-e81e-9704-41aecc3cdc4b.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static implicit operator OnlineData (
	SecsValue value
)
```

```
Public Shared Widening Operator CType ( 
	value As SecsValue
) As OnlineData
```

```
static implicit operator OnlineData^ (
	SecsValue^ value
)
```

```
F# does not support the declaration of new casting operators.
```

#### 参数

value
:   类型：[HslCommunication.Secs.TypesSecsValue](06a37040-2c85-e52a-86bb-2679fa81cf40.htm)  
    [SecsValue](06a37040-2c85-e52a-86bb-2679fa81cf40.htm) 数值

#### 返回值

类型：[OnlineData](94067ba2-5c27-cc1a-c397-236b70c3fecb.htm)  
等值的消息对象

![](../icons/SectionExpanded.png)参见

#### 引用

[OnlineData 类](94067ba2-5c27-cc1a-c397-236b70c3fecb.htm)

[重载](422092b1-d478-8393-60d1-0d47b4d40a44.htm)

[HslCommunication.Secs.Types 命名空间](a75c6382-21bf-e81e-9704-41aecc3cdc4b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## SecsItemType 枚举

[原文連結](http://api.hslcommunication.cn/html/25a238c4-3ec8-ecf1-c580-83e39b5f9595.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Secs.Types](../html/a75c6382-21bf-e81e-9704-41aecc3cdc4b.htm "HslCommunication.Secs.Types")

[ISecs 接口](../html/08134ad0-26a2-2a6b-52a7-0a3f930b5731.htm "ISecs 接口")

[OnlineData 类](../html/94067ba2-5c27-cc1a-c397-236b70c3fecb.htm "OnlineData 类")

[SecsItemType 枚举](../html/25a238c4-3ec8-ecf1-c580-83e39b5f9595.htm "SecsItemType 枚举")

[SecsMessage 类](../html/3594b6c1-13b3-f4d8-6932-cf74c11ea232.htm "SecsMessage 类")

[SecsMessageExtension 类](../html/a5e48ed5-8b1e-79b4-40b8-568c622983d4.htm "SecsMessageExtension 类")

[SecsValue 类](../html/06a37040-2c85-e52a-86bb-2679fa81cf40.htm "SecsValue 类")

[VariableName 类](../html/66b8d742-dcdd-e507-ab9a-87998843e5f1.htm "VariableName 类")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| SecsItemType 枚举 |

数据类型的定义

**命名空间：**
 [HslCommunication.Secs.Types](a75c6382-21bf-e81e-9704-41aecc3cdc4b.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public enum SecsItemType
```

```
Public Enumeration SecsItemType
```

```
public enum class SecsItemType
```

```
type SecsItemType
```

![](../icons/SectionExpanded.png)成员

|  | 成员名称 | 值 | 说明 |
| --- | --- | --- | --- |
|  | List | 0 | 列表数据类型，代号：L |
|  | Bool | 1 | Bool值类型，代号：BOOLEAN |
|  | Binary | 2 | 二进制数据，代号：B |
|  | ASCII | 3 | ASCII编码的字符串，代号：A |
|  | JIS8 | 4 | JIS类型，代号：J |
|  | SByte | 5 | 一个字节的有符号数据，代号：I1 |
|  | Byte | 6 | 一个字节的无符号数据，代号：U1 |
|  | Int16 | 7 | 两个字节的有符号数据，代号：I2 |
|  | UInt16 | 8 | 两个字节的无符号数据，代号：U2 |
|  | Int32 | 9 | 四个字节的有符号数据，代号：I4 |
|  | UInt32 | 10 | 四个字节的无符号数据，代号：U4 |
|  | Int64 | 11 | 八个字节的有符号数据，代号：I8 |
|  | UInt64 | 12 | 八个字节的无符号数据，代号：U8 |
|  | Single | 13 | 四个字节的浮点数，代号：F4 |
|  | Double | 14 | 八个字节的浮点数，代号：F8 |
|  | None | 15 | 这是一个空的类型信息 |

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.Secs.Types 命名空间](a75c6382-21bf-e81e-9704-41aecc3cdc4b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## SecsMessage 类

[原文連結](http://api.hslcommunication.cn/html/3594b6c1-13b3-f4d8-6932-cf74c11ea232.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Secs.Types](../html/a75c6382-21bf-e81e-9704-41aecc3cdc4b.htm "HslCommunication.Secs.Types")

[SecsMessage 类](../html/3594b6c1-13b3-f4d8-6932-cf74c11ea232.htm "SecsMessage 类")

[SecsMessage 构造函数](../html/81cb805f-29be-cd54-2744-0395498d0d86.htm "SecsMessage 构造函数 ")

[SecsMessage 属性](../html/d47b093f-be16-cd01-cc0d-a02f8407bab2.htm "SecsMessage 属性")

[SecsMessage 方法](../html/f5980fc1-ec74-1c53-e787-fbc71915671a.htm "SecsMessage 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| SecsMessage 类 |

Secs的消息类对象

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  HslCommunication.Secs.TypesSecsMessage

**命名空间：**
 [HslCommunication.Secs.Types](a75c6382-21bf-e81e-9704-41aecc3cdc4b.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class SecsMessage
```

```
Public Class SecsMessage
```

```
public ref class SecsMessage
```

```
type SecsMessage =  class end
```

SecsMessage 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [SecsMessage](1f17c47e-1b6e-e3e4-f55d-151a21a2ac4a.htm) | 实例化一个默认的对象 |
| 公共方法 | [SecsMessage(Byte)](40239390-cce2-2f2f-3b8f-f5204106f5dc.htm) | 通过原始的报文信息来实例化一个默认的对象 |
| 公共方法 | [SecsMessage(Byte, Int32)](7352b0d7-a929-a72f-01a6-a3c4bb6612fd.htm) | 通过原始的报文信息来实例化一个默认的对象 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [BlockNo](04e3e901-eb3c-a92c-a1a6-3f8a66cca5eb.htm) | 获取或设置区块号信息 |
| 公共属性 | [Data](a05f39d0-2f9e-9016-81cc-6dbcac9374c8.htm) | 消息数据对象 |
| 公共属性 | [DeviceID](0bbd5052-abc6-bc10-b8bf-385b6b7cb9e2.htm) | 设备的ID信息 |
| 公共属性 | [E](bc85307f-740f-2e88-e340-40e48b3810e5.htm) | E=false, 尚有Block; E=true, 此为最后一个Block |
| 公共属性 | [FunctionNo](b7767047-d202-14d3-4058-5605bb30ef3e.htm) | Function功能码 |
| 公共属性 | [MessageID](fc758947-88bd-8876-07c4-a2f38036d62f.htm) | 获取或设置消息ID信息 |
| 公共属性 | [R](86125e46-b3ee-6b20-b806-e7893f12dfd9.htm) | R=false, Host → Equipment; R=true, Host ← Equipment |
| 公共属性 | [StreamNo](9fe85c68-f851-c04c-398a-03225c7fe237.htm) | Stream功能码 |
| 公共属性 | [StringEncoding](21f37a98-862b-a836-ca5e-a01f75229d66.htm) | 获取或设置用于字符串解析的编码信息 |
| 公共属性 | [W](d0f4bbf6-6afd-01f3-50e4-22665ba3af31.htm) | W=false, 不必回复讯息；W=true, 必须回复讯息 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法 | [GetItemValues](3b6adeef-44b1-7c7f-a914-e807adfae6b9.htm) | 获取当前消息的所有对象信息 |
| 公共方法 | [GetItemValues(Encoding)](9ad29937-2f33-4afe-7e31-cfa91c03f9fc.htm) | 使用指定的编码获取当前消息的所有对象信息 |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法 | [ToString](5f7b801d-5c15-584f-ff02-d8f0594c28f8.htm) | (重写 ObjectToString.) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.Secs.Types 命名空间](a75c6382-21bf-e81e-9704-41aecc3cdc4b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## SecsMessage 构造函数 

[原文連結](http://api.hslcommunication.cn/html/81cb805f-29be-cd54-2744-0395498d0d86.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Secs.Types](../html/a75c6382-21bf-e81e-9704-41aecc3cdc4b.htm "HslCommunication.Secs.Types")

[SecsMessage 类](../html/3594b6c1-13b3-f4d8-6932-cf74c11ea232.htm "SecsMessage 类")

[SecsMessage 构造函数](../html/81cb805f-29be-cd54-2744-0395498d0d86.htm "SecsMessage 构造函数 ")

[SecsMessage 构造函数](../html/1f17c47e-1b6e-e3e4-f55d-151a21a2ac4a.htm "SecsMessage 构造函数 ")

[SecsMessage 构造函数 (Byte[])](../html/40239390-cce2-2f2f-3b8f-f5204106f5dc.htm "SecsMessage 构造函数 (Byte[])")

[SecsMessage 构造函数 (Byte[], Int32)](../html/7352b0d7-a929-a72f-01a6-a3c4bb6612fd.htm "SecsMessage 构造函数 (Byte[], Int32)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| SecsMessage 构造函数 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [SecsMessage](1f17c47e-1b6e-e3e4-f55d-151a21a2ac4a.htm) | 实例化一个默认的对象 |
| 公共方法 | [SecsMessage(Byte)](40239390-cce2-2f2f-3b8f-f5204106f5dc.htm) | 通过原始的报文信息来实例化一个默认的对象 |
| 公共方法 | [SecsMessage(Byte, Int32)](7352b0d7-a929-a72f-01a6-a3c4bb6612fd.htm) | 通过原始的报文信息来实例化一个默认的对象 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[SecsMessage 类](3594b6c1-13b3-f4d8-6932-cf74c11ea232.htm)

[HslCommunication.Secs.Types 命名空间](a75c6382-21bf-e81e-9704-41aecc3cdc4b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## SecsMessage 构造函数 

[原文連結](http://api.hslcommunication.cn/html/1f17c47e-1b6e-e3e4-f55d-151a21a2ac4a.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Secs.Types](../html/a75c6382-21bf-e81e-9704-41aecc3cdc4b.htm "HslCommunication.Secs.Types")

[SecsMessage 类](../html/3594b6c1-13b3-f4d8-6932-cf74c11ea232.htm "SecsMessage 类")

[SecsMessage 构造函数](../html/81cb805f-29be-cd54-2744-0395498d0d86.htm "SecsMessage 构造函数 ")

[SecsMessage 构造函数](../html/1f17c47e-1b6e-e3e4-f55d-151a21a2ac4a.htm "SecsMessage 构造函数 ")

[SecsMessage 构造函数 (Byte[])](../html/40239390-cce2-2f2f-3b8f-f5204106f5dc.htm "SecsMessage 构造函数 (Byte[])")

[SecsMessage 构造函数 (Byte[], Int32)](../html/7352b0d7-a929-a72f-01a6-a3c4bb6612fd.htm "SecsMessage 构造函数 (Byte[], Int32)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| SecsMessage 构造函数 |

实例化一个默认的对象

**命名空间：**
 [HslCommunication.Secs.Types](a75c6382-21bf-e81e-9704-41aecc3cdc4b.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public SecsMessage()
```

```
Public Sub New
```

```
public:
SecsMessage()
```

```
new : unit -> SecsMessage
```

![](../icons/SectionExpanded.png)参见

#### 引用

[SecsMessage 类](3594b6c1-13b3-f4d8-6932-cf74c11ea232.htm)

[SecsMessage 重载](81cb805f-29be-cd54-2744-0395498d0d86.htm)

[HslCommunication.Secs.Types 命名空间](a75c6382-21bf-e81e-9704-41aecc3cdc4b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## SecsMessage 构造函数 (Byte[])

[原文連結](http://api.hslcommunication.cn/html/40239390-cce2-2f2f-3b8f-f5204106f5dc.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Secs.Types](../html/a75c6382-21bf-e81e-9704-41aecc3cdc4b.htm "HslCommunication.Secs.Types")

[SecsMessage 类](../html/3594b6c1-13b3-f4d8-6932-cf74c11ea232.htm "SecsMessage 类")

[SecsMessage 构造函数](../html/81cb805f-29be-cd54-2744-0395498d0d86.htm "SecsMessage 构造函数 ")

[SecsMessage 构造函数](../html/1f17c47e-1b6e-e3e4-f55d-151a21a2ac4a.htm "SecsMessage 构造函数 ")

[SecsMessage 构造函数 (Byte[])](../html/40239390-cce2-2f2f-3b8f-f5204106f5dc.htm "SecsMessage 构造函数 (Byte[])")

[SecsMessage 构造函数 (Byte[], Int32)](../html/7352b0d7-a929-a72f-01a6-a3c4bb6612fd.htm "SecsMessage 构造函数 (Byte[], Int32)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| SecsMessage 构造函数 (Byte) |

通过原始的报文信息来实例化一个默认的对象

**命名空间：**
 [HslCommunication.Secs.Types](a75c6382-21bf-e81e-9704-41aecc3cdc4b.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public SecsMessage(
	byte[] message
)
```

```
Public Sub New ( 
	message As Byte()
)
```

```
public:
SecsMessage(
	array<unsigned char>^ message
)
```

```
new : 
        message : byte[] -> SecsMessage
```

#### 参数

message
:   类型：SystemByte  
    原始的字节信息

![](../icons/SectionExpanded.png)参见

#### 引用

[SecsMessage 类](3594b6c1-13b3-f4d8-6932-cf74c11ea232.htm)

[SecsMessage 重载](81cb805f-29be-cd54-2744-0395498d0d86.htm)

[HslCommunication.Secs.Types 命名空间](a75c6382-21bf-e81e-9704-41aecc3cdc4b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## SecsMessage 构造函数 (Byte[], Int32)

[原文連結](http://api.hslcommunication.cn/html/7352b0d7-a929-a72f-01a6-a3c4bb6612fd.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Secs.Types](../html/a75c6382-21bf-e81e-9704-41aecc3cdc4b.htm "HslCommunication.Secs.Types")

[SecsMessage 类](../html/3594b6c1-13b3-f4d8-6932-cf74c11ea232.htm "SecsMessage 类")

[SecsMessage 构造函数](../html/81cb805f-29be-cd54-2744-0395498d0d86.htm "SecsMessage 构造函数 ")

[SecsMessage 构造函数](../html/1f17c47e-1b6e-e3e4-f55d-151a21a2ac4a.htm "SecsMessage 构造函数 ")

[SecsMessage 构造函数 (Byte[])](../html/40239390-cce2-2f2f-3b8f-f5204106f5dc.htm "SecsMessage 构造函数 (Byte[])")

[SecsMessage 构造函数 (Byte[], Int32)](../html/7352b0d7-a929-a72f-01a6-a3c4bb6612fd.htm "SecsMessage 构造函数 (Byte[], Int32)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| SecsMessage 构造函数 (Byte, Int32) |

通过原始的报文信息来实例化一个默认的对象

**命名空间：**
 [HslCommunication.Secs.Types](a75c6382-21bf-e81e-9704-41aecc3cdc4b.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public SecsMessage(
	byte[] message,
	int startIndex
)
```

```
Public Sub New ( 
	message As Byte(),
	startIndex As Integer
)
```

```
public:
SecsMessage(
	array<unsigned char>^ message, 
	int startIndex
)
```

```
new : 
        message : byte[] * 
        startIndex : int -> SecsMessage
```

#### 参数

message
:   类型：SystemByte  
    原始的字节信息

startIndex
:   类型：SystemInt32  
    起始的偏移地址

![](../icons/SectionExpanded.png)参见

#### 引用

[SecsMessage 类](3594b6c1-13b3-f4d8-6932-cf74c11ea232.htm)

[SecsMessage 重载](81cb805f-29be-cd54-2744-0395498d0d86.htm)

[HslCommunication.Secs.Types 命名空间](a75c6382-21bf-e81e-9704-41aecc3cdc4b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## SecsMessage 属性

[原文連結](http://api.hslcommunication.cn/html/d47b093f-be16-cd01-cc0d-a02f8407bab2.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Secs.Types](../html/a75c6382-21bf-e81e-9704-41aecc3cdc4b.htm "HslCommunication.Secs.Types")

[SecsMessage 类](../html/3594b6c1-13b3-f4d8-6932-cf74c11ea232.htm "SecsMessage 类")

[SecsMessage 属性](../html/d47b093f-be16-cd01-cc0d-a02f8407bab2.htm "SecsMessage 属性")

[BlockNo 属性](../html/04e3e901-eb3c-a92c-a1a6-3f8a66cca5eb.htm "BlockNo 属性 ")

[Data 属性](../html/a05f39d0-2f9e-9016-81cc-6dbcac9374c8.htm "Data 属性 ")

[DeviceID 属性](../html/0bbd5052-abc6-bc10-b8bf-385b6b7cb9e2.htm "DeviceID 属性 ")

[E 属性](../html/bc85307f-740f-2e88-e340-40e48b3810e5.htm "E 属性 ")

[FunctionNo 属性](../html/b7767047-d202-14d3-4058-5605bb30ef3e.htm "FunctionNo 属性 ")

[MessageID 属性](../html/fc758947-88bd-8876-07c4-a2f38036d62f.htm "MessageID 属性 ")

[R 属性](../html/86125e46-b3ee-6b20-b806-e7893f12dfd9.htm "R 属性 ")

[StreamNo 属性](../html/9fe85c68-f851-c04c-398a-03225c7fe237.htm "StreamNo 属性 ")

[StringEncoding 属性](../html/21f37a98-862b-a836-ca5e-a01f75229d66.htm "StringEncoding 属性 ")

[W 属性](../html/d0f4bbf6-6afd-01f3-50e4-22665ba3af31.htm "W 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| SecsMessage 属性 |

[SecsMessage](3594b6c1-13b3-f4d8-6932-cf74c11ea232.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [BlockNo](04e3e901-eb3c-a92c-a1a6-3f8a66cca5eb.htm) | 获取或设置区块号信息 |
| 公共属性 | [Data](a05f39d0-2f9e-9016-81cc-6dbcac9374c8.htm) | 消息数据对象 |
| 公共属性 | [DeviceID](0bbd5052-abc6-bc10-b8bf-385b6b7cb9e2.htm) | 设备的ID信息 |
| 公共属性 | [E](bc85307f-740f-2e88-e340-40e48b3810e5.htm) | E=false, 尚有Block; E=true, 此为最后一个Block |
| 公共属性 | [FunctionNo](b7767047-d202-14d3-4058-5605bb30ef3e.htm) | Function功能码 |
| 公共属性 | [MessageID](fc758947-88bd-8876-07c4-a2f38036d62f.htm) | 获取或设置消息ID信息 |
| 公共属性 | [R](86125e46-b3ee-6b20-b806-e7893f12dfd9.htm) | R=false, Host → Equipment; R=true, Host ← Equipment |
| 公共属性 | [StreamNo](9fe85c68-f851-c04c-398a-03225c7fe237.htm) | Stream功能码 |
| 公共属性 | [StringEncoding](21f37a98-862b-a836-ca5e-a01f75229d66.htm) | 获取或设置用于字符串解析的编码信息 |
| 公共属性 | [W](d0f4bbf6-6afd-01f3-50e4-22665ba3af31.htm) | W=false, 不必回复讯息；W=true, 必须回复讯息 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[SecsMessage 类](3594b6c1-13b3-f4d8-6932-cf74c11ea232.htm)

[HslCommunication.Secs.Types 命名空间](a75c6382-21bf-e81e-9704-41aecc3cdc4b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## BlockNo 属性 

[原文連結](http://api.hslcommunication.cn/html/04e3e901-eb3c-a92c-a1a6-3f8a66cca5eb.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Secs.Types](../html/a75c6382-21bf-e81e-9704-41aecc3cdc4b.htm "HslCommunication.Secs.Types")

[SecsMessage 类](../html/3594b6c1-13b3-f4d8-6932-cf74c11ea232.htm "SecsMessage 类")

[SecsMessage 属性](../html/d47b093f-be16-cd01-cc0d-a02f8407bab2.htm "SecsMessage 属性")

[BlockNo 属性](../html/04e3e901-eb3c-a92c-a1a6-3f8a66cca5eb.htm "BlockNo 属性 ")

[Data 属性](../html/a05f39d0-2f9e-9016-81cc-6dbcac9374c8.htm "Data 属性 ")

[DeviceID 属性](../html/0bbd5052-abc6-bc10-b8bf-385b6b7cb9e2.htm "DeviceID 属性 ")

[E 属性](../html/bc85307f-740f-2e88-e340-40e48b3810e5.htm "E 属性 ")

[FunctionNo 属性](../html/b7767047-d202-14d3-4058-5605bb30ef3e.htm "FunctionNo 属性 ")

[MessageID 属性](../html/fc758947-88bd-8876-07c4-a2f38036d62f.htm "MessageID 属性 ")

[R 属性](../html/86125e46-b3ee-6b20-b806-e7893f12dfd9.htm "R 属性 ")

[StreamNo 属性](../html/9fe85c68-f851-c04c-398a-03225c7fe237.htm "StreamNo 属性 ")

[StringEncoding 属性](../html/21f37a98-862b-a836-ca5e-a01f75229d66.htm "StringEncoding 属性 ")

[W 属性](../html/d0f4bbf6-6afd-01f3-50e4-22665ba3af31.htm "W 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| SecsMessageBlockNo 属性 |

获取或设置区块号信息

**命名空间：**
 [HslCommunication.Secs.Types](a75c6382-21bf-e81e-9704-41aecc3cdc4b.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public int BlockNo { get; set; }
```

```
Public Property BlockNo As Integer
	Get
	Set
```

```
public:
property int BlockNo {
	int get ();
	void set (int value);
}
```

```
member BlockNo : int with get, set
```

#### 属性值

类型：Int32

![](../icons/SectionExpanded.png)参见

#### 引用

[SecsMessage 类](3594b6c1-13b3-f4d8-6932-cf74c11ea232.htm)

[HslCommunication.Secs.Types 命名空间](a75c6382-21bf-e81e-9704-41aecc3cdc4b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Data 属性 

[原文連結](http://api.hslcommunication.cn/html/a05f39d0-2f9e-9016-81cc-6dbcac9374c8.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Secs.Types](../html/a75c6382-21bf-e81e-9704-41aecc3cdc4b.htm "HslCommunication.Secs.Types")

[SecsMessage 类](../html/3594b6c1-13b3-f4d8-6932-cf74c11ea232.htm "SecsMessage 类")

[SecsMessage 属性](../html/d47b093f-be16-cd01-cc0d-a02f8407bab2.htm "SecsMessage 属性")

[BlockNo 属性](../html/04e3e901-eb3c-a92c-a1a6-3f8a66cca5eb.htm "BlockNo 属性 ")

[Data 属性](../html/a05f39d0-2f9e-9016-81cc-6dbcac9374c8.htm "Data 属性 ")

[DeviceID 属性](../html/0bbd5052-abc6-bc10-b8bf-385b6b7cb9e2.htm "DeviceID 属性 ")

[E 属性](../html/bc85307f-740f-2e88-e340-40e48b3810e5.htm "E 属性 ")

[FunctionNo 属性](../html/b7767047-d202-14d3-4058-5605bb30ef3e.htm "FunctionNo 属性 ")

[MessageID 属性](../html/fc758947-88bd-8876-07c4-a2f38036d62f.htm "MessageID 属性 ")

[R 属性](../html/86125e46-b3ee-6b20-b806-e7893f12dfd9.htm "R 属性 ")

[StreamNo 属性](../html/9fe85c68-f851-c04c-398a-03225c7fe237.htm "StreamNo 属性 ")

[StringEncoding 属性](../html/21f37a98-862b-a836-ca5e-a01f75229d66.htm "StringEncoding 属性 ")

[W 属性](../html/d0f4bbf6-6afd-01f3-50e4-22665ba3af31.htm "W 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| SecsMessageData 属性 |

消息数据对象

**命名空间：**
 [HslCommunication.Secs.Types](a75c6382-21bf-e81e-9704-41aecc3cdc4b.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public byte[] Data { get; set; }
```

```
Public Property Data As Byte()
	Get
	Set
```

```
public:
property array<unsigned char>^ Data {
	array<unsigned char>^ get ();
	void set (array<unsigned char>^ value);
}
```

```
member Data : byte[] with get, set
```

#### 属性值

类型：Byte

![](../icons/SectionExpanded.png)参见

#### 引用

[SecsMessage 类](3594b6c1-13b3-f4d8-6932-cf74c11ea232.htm)

[HslCommunication.Secs.Types 命名空间](a75c6382-21bf-e81e-9704-41aecc3cdc4b.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)