# HslCommunication - HslCommunication.Instrument.CJT.Helper

> 分類頁數: 30



---
## HslCommunication.Instrument.CJT.Helper

[原文連結](http://api.hslcommunication.cn/html/e7617fc1-440e-5d36-1b24-a85e8d4dffb0.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.CJT.Helper](../html/e7617fc1-440e-5d36-1b24-a85e8d4dffb0.htm "HslCommunication.Instrument.CJT.Helper")

[CJT188Helper 类](../html/2e426e36-c009-46b5-43dd-b42e3a2b11c5.htm "CJT188Helper 类")

[CJTControl 类](../html/70e1d569-367b-1da4-fd16-2592bb57b9bf.htm "CJTControl 类")

[ICjt188 接口](../html/3fabd38b-b4c3-f1b2-0469-ad7704214895.htm "ICjt188 接口")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslCommunication.Instrument.CJT.Helper 命名空间 |

[缺少 "N:HslCommunication.Instrument.CJT.Helper" 的 <summary> 文档]

![](../icons/SectionExpanded.png)类

|  | 类 | 说明 |
| --- | --- | --- |
| 公共类 | [CJT188Helper](2e426e36-c009-46b5-43dd-b42e3a2b11c5.htm) | CJT188辅助方法 |
| 公共类 | [CJTControl](70e1d569-367b-1da4-fd16-2592bb57b9bf.htm) | 控制码信息 |

![](../icons/SectionExpanded.png)接口

|  | 接口 | 说明 |
| --- | --- | --- |
| 公共接口 | [ICjt188](3fabd38b-b4c3-f1b2-0469-ad7704214895.htm) | CJT188设备的接口 |

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## CJT188Helper 类

[原文連結](http://api.hslcommunication.cn/html/2e426e36-c009-46b5-43dd-b42e3a2b11c5.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.CJT.Helper](../html/e7617fc1-440e-5d36-1b24-a85e8d4dffb0.htm "HslCommunication.Instrument.CJT.Helper")

[CJT188Helper 类](../html/2e426e36-c009-46b5-43dd-b42e3a2b11c5.htm "CJT188Helper 类")

[CJT188Helper 构造函数](../html/a885502b-9fdb-9bd3-984e-98558fd98f3b.htm "CJT188Helper 构造函数 ")

[CJT188Helper 方法](../html/b6d542ae-0467-f6cb-eec2-8e1313c84576.htm "CJT188Helper 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CJT188Helper 类 |

CJT188辅助方法

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  HslCommunication.Instrument.CJT.HelperCJT188Helper

**命名空间：**
 [HslCommunication.Instrument.CJT.Helper](e7617fc1-440e-5d36-1b24-a85e8d4dffb0.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class CJT188Helper
```

```
Public Class CJT188Helper
```

```
public ref class CJT188Helper
```

```
type CJT188Helper =  class end
```

CJT188Helper 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [CJT188Helper](a885502b-9fdb-9bd3-984e-98558fd98f3b.htm) | 初始化 CJT188Helper 类的一个新实例 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法静态成员 | [AnalysisBytesAddress](19f104be-2349-b906-d8f3-bb6c6a9e5738.htm) | 从用户输入的地址信息中解析出真实的地址及数据标识 |
| 公共方法静态成员 | [Build188EntireCommand](240f6ac6-815e-f3cf-0b0f-f928e65765a3.htm) | 将指定的地址信息，控制码信息，数据域信息打包成完整的报文命令 |
| 公共方法静态成员 | [CheckResponse](c0df1d6d-197f-6039-ce25-8afb6366feea.htm) | 检查当前的反馈数据信息是否正确 |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法静态成员 | [GetAddressByteFromString](88fb6eb9-c67b-28b8-6302-52517c8cc4bc.htm) | 将地址解析成BCD码的地址，并且扩充到14位，不够的补0操作 |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法静态成员 | [Read](51e94d04-015b-a8eb-4522-f6fa925ea30c.htm) | 根据指定的数据标识来读取相关的原始数据信息，地址标识根据手册来，从高位到地位，例如 91-1F，分割符可以任意特殊字符或是没有分隔符。  Read the relevant original data information according to the specified data identifier. The address identifier is based on the manual, from high to position, such as 91-1F. The separator can be any special character or no separator. |
| 公共方法静态成员 | [ReadAddress](a81b1674-c8a5-1e05-00e5-e69fd0a3b48b.htm) | 读取设备的通信地址，仅支持点对点通讯的情况，返回地址域数据，例如：14910000729012  Read the communication address of the device, only support point-to-point communication, and return the address field data, for example: 149100007290 |
| 公共方法静态成员 | [ReadStringArray](bc790e20-12da-1989-1a45-ed11fcc08909.htm) | 读取指定地址的所有的字符串数据信息，一般来说，一个地址只有一个数据 |
| 公共方法静态成员 | [ReadValueT](06ba322b-7249-60c5-bd64-749fbbb3f6cb.htm) | 读取数据的数组信息，需要指定如何从字符串转换的功能方法 |
| 公共方法 | ToString | (继承自 Object。) |
| 公共方法静态成员 | [Write](1d27c65a-d3db-8521-04af-d708904e71d8.htm) | 根据指定的数据标识来写入相关的原始数据信息，地址标识根据手册来，从高位到地位，例如 90-1F，分割符可以任意特殊字符或是没有分隔符。  Read the relevant original data information according to the specified data identifier. The address identifier is based on the manual, from high to position, such as 90-1F. The separator can be any special character or no separator. |
| 公共方法静态成员 | [WriteAddress](f54f5a32-f807-edcf-f43b-9b07c29b840c.htm) | 写入设备的地址域信息，仅支持点对点通讯的情况，需要指定地址域信息，例如：14910000729011  Write the address domain information of the device, only support point-to-point communication, you need to specify the address domain information, for example: 14910000729011 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.Instrument.CJT.Helper 命名空间](e7617fc1-440e-5d36-1b24-a85e8d4dffb0.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## CJT188Helper 构造函数 

[原文連結](http://api.hslcommunication.cn/html/a885502b-9fdb-9bd3-984e-98558fd98f3b.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.CJT.Helper](../html/e7617fc1-440e-5d36-1b24-a85e8d4dffb0.htm "HslCommunication.Instrument.CJT.Helper")

[CJT188Helper 类](../html/2e426e36-c009-46b5-43dd-b42e3a2b11c5.htm "CJT188Helper 类")

[CJT188Helper 构造函数](../html/a885502b-9fdb-9bd3-984e-98558fd98f3b.htm "CJT188Helper 构造函数 ")

[CJT188Helper 方法](../html/b6d542ae-0467-f6cb-eec2-8e1313c84576.htm "CJT188Helper 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CJT188Helper 构造函数 |

初始化 [CJT188Helper](2e426e36-c009-46b5-43dd-b42e3a2b11c5.htm) 类的一个新实例

**命名空间：**
 [HslCommunication.Instrument.CJT.Helper](e7617fc1-440e-5d36-1b24-a85e8d4dffb0.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public CJT188Helper()
```

```
Public Sub New
```

```
public:
CJT188Helper()
```

```
new : unit -> CJT188Helper
```

![](../icons/SectionExpanded.png)参见

#### 引用

[CJT188Helper 类](2e426e36-c009-46b5-43dd-b42e3a2b11c5.htm)

[HslCommunication.Instrument.CJT.Helper 命名空间](e7617fc1-440e-5d36-1b24-a85e8d4dffb0.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## CJT188Helper 方法

[原文連結](http://api.hslcommunication.cn/html/b6d542ae-0467-f6cb-eec2-8e1313c84576.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.CJT.Helper](../html/e7617fc1-440e-5d36-1b24-a85e8d4dffb0.htm "HslCommunication.Instrument.CJT.Helper")

[CJT188Helper 类](../html/2e426e36-c009-46b5-43dd-b42e3a2b11c5.htm "CJT188Helper 类")

[CJT188Helper 方法](../html/b6d542ae-0467-f6cb-eec2-8e1313c84576.htm "CJT188Helper 方法")

[AnalysisBytesAddress 方法](../html/19f104be-2349-b906-d8f3-bb6c6a9e5738.htm "AnalysisBytesAddress 方法 ")

[Build188EntireCommand 方法](../html/240f6ac6-815e-f3cf-0b0f-f928e65765a3.htm "Build188EntireCommand 方法 ")

[CheckResponse 方法](../html/c0df1d6d-197f-6039-ce25-8afb6366feea.htm "CheckResponse 方法 ")

[GetAddressByteFromString 方法](../html/88fb6eb9-c67b-28b8-6302-52517c8cc4bc.htm "GetAddressByteFromString 方法 ")

[Read 方法](../html/51e94d04-015b-a8eb-4522-f6fa925ea30c.htm "Read 方法 ")

[ReadAddress 方法](../html/a81b1674-c8a5-1e05-00e5-e69fd0a3b48b.htm "ReadAddress 方法 ")

[ReadStringArray 方法](../html/bc790e20-12da-1989-1a45-ed11fcc08909.htm "ReadStringArray 方法 ")

[ReadValue(T) 方法](../html/06ba322b-7249-60c5-bd64-749fbbb3f6cb.htm "ReadValue(T) 方法 ")

[Write 方法](../html/1d27c65a-d3db-8521-04af-d708904e71d8.htm "Write 方法 ")

[WriteAddress 方法](../html/f54f5a32-f807-edcf-f43b-9b07c29b840c.htm "WriteAddress 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CJT188Helper 方法 |

[CJT188Helper](2e426e36-c009-46b5-43dd-b42e3a2b11c5.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法静态成员 | [AnalysisBytesAddress](19f104be-2349-b906-d8f3-bb6c6a9e5738.htm) | 从用户输入的地址信息中解析出真实的地址及数据标识 |
| 公共方法静态成员 | [Build188EntireCommand](240f6ac6-815e-f3cf-0b0f-f928e65765a3.htm) | 将指定的地址信息，控制码信息，数据域信息打包成完整的报文命令 |
| 公共方法静态成员 | [CheckResponse](c0df1d6d-197f-6039-ce25-8afb6366feea.htm) | 检查当前的反馈数据信息是否正确 |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法静态成员 | [GetAddressByteFromString](88fb6eb9-c67b-28b8-6302-52517c8cc4bc.htm) | 将地址解析成BCD码的地址，并且扩充到14位，不够的补0操作 |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法静态成员 | [Read](51e94d04-015b-a8eb-4522-f6fa925ea30c.htm) | 根据指定的数据标识来读取相关的原始数据信息，地址标识根据手册来，从高位到地位，例如 91-1F，分割符可以任意特殊字符或是没有分隔符。  Read the relevant original data information according to the specified data identifier. The address identifier is based on the manual, from high to position, such as 91-1F. The separator can be any special character or no separator. |
| 公共方法静态成员 | [ReadAddress](a81b1674-c8a5-1e05-00e5-e69fd0a3b48b.htm) | 读取设备的通信地址，仅支持点对点通讯的情况，返回地址域数据，例如：14910000729012  Read the communication address of the device, only support point-to-point communication, and return the address field data, for example: 149100007290 |
| 公共方法静态成员 | [ReadStringArray](bc790e20-12da-1989-1a45-ed11fcc08909.htm) | 读取指定地址的所有的字符串数据信息，一般来说，一个地址只有一个数据 |
| 公共方法静态成员 | [ReadValueT](06ba322b-7249-60c5-bd64-749fbbb3f6cb.htm) | 读取数据的数组信息，需要指定如何从字符串转换的功能方法 |
| 公共方法 | ToString | (继承自 Object。) |
| 公共方法静态成员 | [Write](1d27c65a-d3db-8521-04af-d708904e71d8.htm) | 根据指定的数据标识来写入相关的原始数据信息，地址标识根据手册来，从高位到地位，例如 90-1F，分割符可以任意特殊字符或是没有分隔符。  Read the relevant original data information according to the specified data identifier. The address identifier is based on the manual, from high to position, such as 90-1F. The separator can be any special character or no separator. |
| 公共方法静态成员 | [WriteAddress](f54f5a32-f807-edcf-f43b-9b07c29b840c.htm) | 写入设备的地址域信息，仅支持点对点通讯的情况，需要指定地址域信息，例如：14910000729011  Write the address domain information of the device, only support point-to-point communication, you need to specify the address domain information, for example: 14910000729011 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[CJT188Helper 类](2e426e36-c009-46b5-43dd-b42e3a2b11c5.htm)

[HslCommunication.Instrument.CJT.Helper 命名空间](e7617fc1-440e-5d36-1b24-a85e8d4dffb0.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## AnalysisBytesAddress 方法 

[原文連結](http://api.hslcommunication.cn/html/19f104be-2349-b906-d8f3-bb6c6a9e5738.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.CJT.Helper](../html/e7617fc1-440e-5d36-1b24-a85e8d4dffb0.htm "HslCommunication.Instrument.CJT.Helper")

[CJT188Helper 类](../html/2e426e36-c009-46b5-43dd-b42e3a2b11c5.htm "CJT188Helper 类")

[CJT188Helper 方法](../html/b6d542ae-0467-f6cb-eec2-8e1313c84576.htm "CJT188Helper 方法")

[AnalysisBytesAddress 方法](../html/19f104be-2349-b906-d8f3-bb6c6a9e5738.htm "AnalysisBytesAddress 方法 ")

[Build188EntireCommand 方法](../html/240f6ac6-815e-f3cf-0b0f-f928e65765a3.htm "Build188EntireCommand 方法 ")

[CheckResponse 方法](../html/c0df1d6d-197f-6039-ce25-8afb6366feea.htm "CheckResponse 方法 ")

[GetAddressByteFromString 方法](../html/88fb6eb9-c67b-28b8-6302-52517c8cc4bc.htm "GetAddressByteFromString 方法 ")

[Read 方法](../html/51e94d04-015b-a8eb-4522-f6fa925ea30c.htm "Read 方法 ")

[ReadAddress 方法](../html/a81b1674-c8a5-1e05-00e5-e69fd0a3b48b.htm "ReadAddress 方法 ")

[ReadStringArray 方法](../html/bc790e20-12da-1989-1a45-ed11fcc08909.htm "ReadStringArray 方法 ")

[ReadValue(T) 方法](../html/06ba322b-7249-60c5-bd64-749fbbb3f6cb.htm "ReadValue(T) 方法 ")

[Write 方法](../html/1d27c65a-d3db-8521-04af-d708904e71d8.htm "Write 方法 ")

[WriteAddress 方法](../html/f54f5a32-f807-edcf-f43b-9b07c29b840c.htm "WriteAddress 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CJT188HelperAnalysisBytesAddress 方法 |

从用户输入的地址信息中解析出真实的地址及数据标识

**命名空间：**
 [HslCommunication.Instrument.CJT.Helper](e7617fc1-440e-5d36-1b24-a85e8d4dffb0.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult<string, byte[]> AnalysisBytesAddress(
	string address,
	string defaultStation
)
```

```
Public Shared Function AnalysisBytesAddress ( 
	address As String,
	defaultStation As String
) As OperateResult(Of String, Byte())
```

```
public:
static OperateResult<String^, array<unsigned char>^>^ AnalysisBytesAddress(
	String^ address, 
	String^ defaultStation
)
```

```
static member AnalysisBytesAddress : 
        address : string * 
        defaultStation : string -> OperateResult<string, byte[]> 
```

#### 参数

address
:   类型：SystemString  
    用户输入的地址信息

defaultStation
:   类型：SystemString  
    默认的地址域

#### 返回值

类型：[OperateResult](f52f888f-5e8d-b0c4-2302-81e70c230a31.htm)String, Byte  
解析结果信息

![](../icons/SectionExpanded.png)参见

#### 引用

[CJT188Helper 类](2e426e36-c009-46b5-43dd-b42e3a2b11c5.htm)

[HslCommunication.Instrument.CJT.Helper 命名空间](e7617fc1-440e-5d36-1b24-a85e8d4dffb0.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Build188EntireCommand 方法 

[原文連結](http://api.hslcommunication.cn/html/240f6ac6-815e-f3cf-0b0f-f928e65765a3.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.CJT.Helper](../html/e7617fc1-440e-5d36-1b24-a85e8d4dffb0.htm "HslCommunication.Instrument.CJT.Helper")

[CJT188Helper 类](../html/2e426e36-c009-46b5-43dd-b42e3a2b11c5.htm "CJT188Helper 类")

[CJT188Helper 方法](../html/b6d542ae-0467-f6cb-eec2-8e1313c84576.htm "CJT188Helper 方法")

[AnalysisBytesAddress 方法](../html/19f104be-2349-b906-d8f3-bb6c6a9e5738.htm "AnalysisBytesAddress 方法 ")

[Build188EntireCommand 方法](../html/240f6ac6-815e-f3cf-0b0f-f928e65765a3.htm "Build188EntireCommand 方法 ")

[CheckResponse 方法](../html/c0df1d6d-197f-6039-ce25-8afb6366feea.htm "CheckResponse 方法 ")

[GetAddressByteFromString 方法](../html/88fb6eb9-c67b-28b8-6302-52517c8cc4bc.htm "GetAddressByteFromString 方法 ")

[Read 方法](../html/51e94d04-015b-a8eb-4522-f6fa925ea30c.htm "Read 方法 ")

[ReadAddress 方法](../html/a81b1674-c8a5-1e05-00e5-e69fd0a3b48b.htm "ReadAddress 方法 ")

[ReadStringArray 方法](../html/bc790e20-12da-1989-1a45-ed11fcc08909.htm "ReadStringArray 方法 ")

[ReadValue(T) 方法](../html/06ba322b-7249-60c5-bd64-749fbbb3f6cb.htm "ReadValue(T) 方法 ")

[Write 方法](../html/1d27c65a-d3db-8521-04af-d708904e71d8.htm "Write 方法 ")

[WriteAddress 方法](../html/f54f5a32-f807-edcf-f43b-9b07c29b840c.htm "WriteAddress 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CJT188HelperBuild188EntireCommand 方法 |

将指定的地址信息，控制码信息，数据域信息打包成完整的报文命令

**命名空间：**
 [HslCommunication.Instrument.CJT.Helper](e7617fc1-440e-5d36-1b24-a85e8d4dffb0.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult<byte[]> Build188EntireCommand(
	string address,
	byte type,
	byte control,
	byte[] dataArea
)
```

```
Public Shared Function Build188EntireCommand ( 
	address As String,
	type As Byte,
	control As Byte,
	dataArea As Byte()
) As OperateResult(Of Byte())
```

```
public:
static OperateResult<array<unsigned char>^>^ Build188EntireCommand(
	String^ address, 
	unsigned char type, 
	unsigned char control, 
	array<unsigned char>^ dataArea
)
```

```
static member Build188EntireCommand : 
        address : string * 
        type : byte * 
        control : byte * 
        dataArea : byte[] -> OperateResult<byte[]> 
```

#### 参数

address
:   类型：SystemString  
    地址域信息，地址域由7个字节构成，每字节2位BCD码，地址长度可达14位十进制数。地址域支持锁位寻址，即从若干低位起，剩余高位补AAH作为通配符进行读表操作

type
:   类型：SystemByte  
    仪表类型

control
:   类型：SystemByte  
    控制码信息

dataArea
:   类型：SystemByte  
    数据域的内容

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
返回是否报文创建成功

![](../icons/SectionExpanded.png)参见

#### 引用

[CJT188Helper 类](2e426e36-c009-46b5-43dd-b42e3a2b11c5.htm)

[HslCommunication.Instrument.CJT.Helper 命名空间](e7617fc1-440e-5d36-1b24-a85e8d4dffb0.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## CheckResponse 方法 

[原文連結](http://api.hslcommunication.cn/html/c0df1d6d-197f-6039-ce25-8afb6366feea.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.CJT.Helper](../html/e7617fc1-440e-5d36-1b24-a85e8d4dffb0.htm "HslCommunication.Instrument.CJT.Helper")

[CJT188Helper 类](../html/2e426e36-c009-46b5-43dd-b42e3a2b11c5.htm "CJT188Helper 类")

[CJT188Helper 方法](../html/b6d542ae-0467-f6cb-eec2-8e1313c84576.htm "CJT188Helper 方法")

[AnalysisBytesAddress 方法](../html/19f104be-2349-b906-d8f3-bb6c6a9e5738.htm "AnalysisBytesAddress 方法 ")

[Build188EntireCommand 方法](../html/240f6ac6-815e-f3cf-0b0f-f928e65765a3.htm "Build188EntireCommand 方法 ")

[CheckResponse 方法](../html/c0df1d6d-197f-6039-ce25-8afb6366feea.htm "CheckResponse 方法 ")

[GetAddressByteFromString 方法](../html/88fb6eb9-c67b-28b8-6302-52517c8cc4bc.htm "GetAddressByteFromString 方法 ")

[Read 方法](../html/51e94d04-015b-a8eb-4522-f6fa925ea30c.htm "Read 方法 ")

[ReadAddress 方法](../html/a81b1674-c8a5-1e05-00e5-e69fd0a3b48b.htm "ReadAddress 方法 ")

[ReadStringArray 方法](../html/bc790e20-12da-1989-1a45-ed11fcc08909.htm "ReadStringArray 方法 ")

[ReadValue(T) 方法](../html/06ba322b-7249-60c5-bd64-749fbbb3f6cb.htm "ReadValue(T) 方法 ")

[Write 方法](../html/1d27c65a-d3db-8521-04af-d708904e71d8.htm "Write 方法 ")

[WriteAddress 方法](../html/f54f5a32-f807-edcf-f43b-9b07c29b840c.htm "WriteAddress 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CJT188HelperCheckResponse 方法 |

检查当前的反馈数据信息是否正确

**命名空间：**
 [HslCommunication.Instrument.CJT.Helper](e7617fc1-440e-5d36-1b24-a85e8d4dffb0.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult CheckResponse(
	ICjt188 cjt,
	byte[] response
)
```

```
Public Shared Function CheckResponse ( 
	cjt As ICjt188,
	response As Byte()
) As OperateResult
```

```
public:
static OperateResult^ CheckResponse(
	ICjt188^ cjt, 
	array<unsigned char>^ response
)
```

```
static member CheckResponse : 
        cjt : ICjt188 * 
        response : byte[] -> OperateResult 
```

#### 参数

cjt
:   类型：[HslCommunication.Instrument.CJT.HelperICjt188](3fabd38b-b4c3-f1b2-0469-ad7704214895.htm)  
    CJT188的通信对象

response
:   类型：SystemByte  
    从仪表反馈的数据信息

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否校验成功

![](../icons/SectionExpanded.png)参见

#### 引用

[CJT188Helper 类](2e426e36-c009-46b5-43dd-b42e3a2b11c5.htm)

[HslCommunication.Instrument.CJT.Helper 命名空间](e7617fc1-440e-5d36-1b24-a85e8d4dffb0.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GetAddressByteFromString 方法 

[原文連結](http://api.hslcommunication.cn/html/88fb6eb9-c67b-28b8-6302-52517c8cc4bc.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.CJT.Helper](../html/e7617fc1-440e-5d36-1b24-a85e8d4dffb0.htm "HslCommunication.Instrument.CJT.Helper")

[CJT188Helper 类](../html/2e426e36-c009-46b5-43dd-b42e3a2b11c5.htm "CJT188Helper 类")

[CJT188Helper 方法](../html/b6d542ae-0467-f6cb-eec2-8e1313c84576.htm "CJT188Helper 方法")

[AnalysisBytesAddress 方法](../html/19f104be-2349-b906-d8f3-bb6c6a9e5738.htm "AnalysisBytesAddress 方法 ")

[Build188EntireCommand 方法](../html/240f6ac6-815e-f3cf-0b0f-f928e65765a3.htm "Build188EntireCommand 方法 ")

[CheckResponse 方法](../html/c0df1d6d-197f-6039-ce25-8afb6366feea.htm "CheckResponse 方法 ")

[GetAddressByteFromString 方法](../html/88fb6eb9-c67b-28b8-6302-52517c8cc4bc.htm "GetAddressByteFromString 方法 ")

[Read 方法](../html/51e94d04-015b-a8eb-4522-f6fa925ea30c.htm "Read 方法 ")

[ReadAddress 方法](../html/a81b1674-c8a5-1e05-00e5-e69fd0a3b48b.htm "ReadAddress 方法 ")

[ReadStringArray 方法](../html/bc790e20-12da-1989-1a45-ed11fcc08909.htm "ReadStringArray 方法 ")

[ReadValue(T) 方法](../html/06ba322b-7249-60c5-bd64-749fbbb3f6cb.htm "ReadValue(T) 方法 ")

[Write 方法](../html/1d27c65a-d3db-8521-04af-d708904e71d8.htm "Write 方法 ")

[WriteAddress 方法](../html/f54f5a32-f807-edcf-f43b-9b07c29b840c.htm "WriteAddress 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CJT188HelperGetAddressByteFromString 方法 |

将地址解析成BCD码的地址，并且扩充到14位，不够的补0操作

**命名空间：**
 [HslCommunication.Instrument.CJT.Helper](e7617fc1-440e-5d36-1b24-a85e8d4dffb0.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult<byte[]> GetAddressByteFromString(
	string address
)
```

```
Public Shared Function GetAddressByteFromString ( 
	address As String
) As OperateResult(Of Byte())
```

```
public:
static OperateResult<array<unsigned char>^>^ GetAddressByteFromString(
	String^ address
)
```

```
static member GetAddressByteFromString : 
        address : string -> OperateResult<byte[]> 
```

#### 参数

address
:   类型：SystemString  
    地址域信息

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
实际的结果

![](../icons/SectionExpanded.png)参见

#### 引用

[CJT188Helper 类](2e426e36-c009-46b5-43dd-b42e3a2b11c5.htm)

[HslCommunication.Instrument.CJT.Helper 命名空间](e7617fc1-440e-5d36-1b24-a85e8d4dffb0.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Read 方法 

[原文連結](http://api.hslcommunication.cn/html/51e94d04-015b-a8eb-4522-f6fa925ea30c.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.CJT.Helper](../html/e7617fc1-440e-5d36-1b24-a85e8d4dffb0.htm "HslCommunication.Instrument.CJT.Helper")

[CJT188Helper 类](../html/2e426e36-c009-46b5-43dd-b42e3a2b11c5.htm "CJT188Helper 类")

[CJT188Helper 方法](../html/b6d542ae-0467-f6cb-eec2-8e1313c84576.htm "CJT188Helper 方法")

[AnalysisBytesAddress 方法](../html/19f104be-2349-b906-d8f3-bb6c6a9e5738.htm "AnalysisBytesAddress 方法 ")

[Build188EntireCommand 方法](../html/240f6ac6-815e-f3cf-0b0f-f928e65765a3.htm "Build188EntireCommand 方法 ")

[CheckResponse 方法](../html/c0df1d6d-197f-6039-ce25-8afb6366feea.htm "CheckResponse 方法 ")

[GetAddressByteFromString 方法](../html/88fb6eb9-c67b-28b8-6302-52517c8cc4bc.htm "GetAddressByteFromString 方法 ")

[Read 方法](../html/51e94d04-015b-a8eb-4522-f6fa925ea30c.htm "Read 方法 ")

[ReadAddress 方法](../html/a81b1674-c8a5-1e05-00e5-e69fd0a3b48b.htm "ReadAddress 方法 ")

[ReadStringArray 方法](../html/bc790e20-12da-1989-1a45-ed11fcc08909.htm "ReadStringArray 方法 ")

[ReadValue(T) 方法](../html/06ba322b-7249-60c5-bd64-749fbbb3f6cb.htm "ReadValue(T) 方法 ")

[Write 方法](../html/1d27c65a-d3db-8521-04af-d708904e71d8.htm "Write 方法 ")

[WriteAddress 方法](../html/f54f5a32-f807-edcf-f43b-9b07c29b840c.htm "WriteAddress 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CJT188HelperRead 方法 |

根据指定的数据标识来读取相关的原始数据信息，地址标识根据手册来，从高位到地位，例如 91-1F，分割符可以任意特殊字符或是没有分隔符。  
Read the relevant original data information according to the specified data identifier. The address identifier is based on the manual,
from high to position, such as 91-1F. The separator can be any special character or no separator.

**命名空间：**
 [HslCommunication.Instrument.CJT.Helper](e7617fc1-440e-5d36-1b24-a85e8d4dffb0.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult<byte[]> Read(
	ICjt188 cjt,
	string address,
	int length
)
```

```
Public Shared Function Read ( 
	cjt As ICjt188,
	address As String,
	length As Integer
) As OperateResult(Of Byte())
```

```
public:
static OperateResult<array<unsigned char>^>^ Read(
	ICjt188^ cjt, 
	String^ address, 
	int length
)
```

```
static member Read : 
        cjt : ICjt188 * 
        address : string * 
        length : int -> OperateResult<byte[]> 
```

#### 参数

cjt
:   类型：[HslCommunication.Instrument.CJT.HelperICjt188](3fabd38b-b4c3-f1b2-0469-ad7704214895.htm)  
    CJT188的通信对象

address
:   类型：SystemString  
    数据标识，具体需要查找手册来对应

length
:   类型：SystemInt32  
    数据长度信息

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
结果信息

![](../icons/SectionExpanded.png)备注

地址可以携带地址域信息，例如 "s=2;90-1F" 或是 "s=100000;90-1F"，关于数据域信息，需要查找手册，例如:D1-20 表示： 上一月结算日累积流量

![](../icons/SectionExpanded.png)参见

#### 引用

[CJT188Helper 类](2e426e36-c009-46b5-43dd-b42e3a2b11c5.htm)

[HslCommunication.Instrument.CJT.Helper 命名空间](e7617fc1-440e-5d36-1b24-a85e8d4dffb0.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadAddress 方法 

[原文連結](http://api.hslcommunication.cn/html/a81b1674-c8a5-1e05-00e5-e69fd0a3b48b.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.CJT.Helper](../html/e7617fc1-440e-5d36-1b24-a85e8d4dffb0.htm "HslCommunication.Instrument.CJT.Helper")

[CJT188Helper 类](../html/2e426e36-c009-46b5-43dd-b42e3a2b11c5.htm "CJT188Helper 类")

[CJT188Helper 方法](../html/b6d542ae-0467-f6cb-eec2-8e1313c84576.htm "CJT188Helper 方法")

[AnalysisBytesAddress 方法](../html/19f104be-2349-b906-d8f3-bb6c6a9e5738.htm "AnalysisBytesAddress 方法 ")

[Build188EntireCommand 方法](../html/240f6ac6-815e-f3cf-0b0f-f928e65765a3.htm "Build188EntireCommand 方法 ")

[CheckResponse 方法](../html/c0df1d6d-197f-6039-ce25-8afb6366feea.htm "CheckResponse 方法 ")

[GetAddressByteFromString 方法](../html/88fb6eb9-c67b-28b8-6302-52517c8cc4bc.htm "GetAddressByteFromString 方法 ")

[Read 方法](../html/51e94d04-015b-a8eb-4522-f6fa925ea30c.htm "Read 方法 ")

[ReadAddress 方法](../html/a81b1674-c8a5-1e05-00e5-e69fd0a3b48b.htm "ReadAddress 方法 ")

[ReadStringArray 方法](../html/bc790e20-12da-1989-1a45-ed11fcc08909.htm "ReadStringArray 方法 ")

[ReadValue(T) 方法](../html/06ba322b-7249-60c5-bd64-749fbbb3f6cb.htm "ReadValue(T) 方法 ")

[Write 方法](../html/1d27c65a-d3db-8521-04af-d708904e71d8.htm "Write 方法 ")

[WriteAddress 方法](../html/f54f5a32-f807-edcf-f43b-9b07c29b840c.htm "WriteAddress 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CJT188HelperReadAddress 方法 |

读取设备的通信地址，仅支持点对点通讯的情况，返回地址域数据，例如：14910000729012  
Read the communication address of the device, only support point-to-point communication, and return the address field data, for example: 149100007290

**命名空间：**
 [HslCommunication.Instrument.CJT.Helper](e7617fc1-440e-5d36-1b24-a85e8d4dffb0.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult<string> ReadAddress(
	ICjt188 cjt
)
```

```
Public Shared Function ReadAddress ( 
	cjt As ICjt188
) As OperateResult(Of String)
```

```
public:
static OperateResult<String^>^ ReadAddress(
	ICjt188^ cjt
)
```

```
static member ReadAddress : 
        cjt : ICjt188 -> OperateResult<string> 
```

#### 参数

cjt
:   类型：[HslCommunication.Instrument.CJT.HelperICjt188](3fabd38b-b4c3-f1b2-0469-ad7704214895.htm)  
    CJT188的通信对象

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)String  
设备的通信地址

![](../icons/SectionExpanded.png)参见

#### 引用

[CJT188Helper 类](2e426e36-c009-46b5-43dd-b42e3a2b11c5.htm)

[HslCommunication.Instrument.CJT.Helper 命名空间](e7617fc1-440e-5d36-1b24-a85e8d4dffb0.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadStringArray 方法 

[原文連結](http://api.hslcommunication.cn/html/bc790e20-12da-1989-1a45-ed11fcc08909.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.CJT.Helper](../html/e7617fc1-440e-5d36-1b24-a85e8d4dffb0.htm "HslCommunication.Instrument.CJT.Helper")

[CJT188Helper 类](../html/2e426e36-c009-46b5-43dd-b42e3a2b11c5.htm "CJT188Helper 类")

[CJT188Helper 方法](../html/b6d542ae-0467-f6cb-eec2-8e1313c84576.htm "CJT188Helper 方法")

[AnalysisBytesAddress 方法](../html/19f104be-2349-b906-d8f3-bb6c6a9e5738.htm "AnalysisBytesAddress 方法 ")

[Build188EntireCommand 方法](../html/240f6ac6-815e-f3cf-0b0f-f928e65765a3.htm "Build188EntireCommand 方法 ")

[CheckResponse 方法](../html/c0df1d6d-197f-6039-ce25-8afb6366feea.htm "CheckResponse 方法 ")

[GetAddressByteFromString 方法](../html/88fb6eb9-c67b-28b8-6302-52517c8cc4bc.htm "GetAddressByteFromString 方法 ")

[Read 方法](../html/51e94d04-015b-a8eb-4522-f6fa925ea30c.htm "Read 方法 ")

[ReadAddress 方法](../html/a81b1674-c8a5-1e05-00e5-e69fd0a3b48b.htm "ReadAddress 方法 ")

[ReadStringArray 方法](../html/bc790e20-12da-1989-1a45-ed11fcc08909.htm "ReadStringArray 方法 ")

[ReadValue(T) 方法](../html/06ba322b-7249-60c5-bd64-749fbbb3f6cb.htm "ReadValue(T) 方法 ")

[Write 方法](../html/1d27c65a-d3db-8521-04af-d708904e71d8.htm "Write 方法 ")

[WriteAddress 方法](../html/f54f5a32-f807-edcf-f43b-9b07c29b840c.htm "WriteAddress 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CJT188HelperReadStringArray 方法 |

读取指定地址的所有的字符串数据信息，一般来说，一个地址只有一个数据

**命名空间：**
 [HslCommunication.Instrument.CJT.Helper](e7617fc1-440e-5d36-1b24-a85e8d4dffb0.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult<string[]> ReadStringArray(
	ICjt188 cjt,
	string address
)
```

```
Public Shared Function ReadStringArray ( 
	cjt As ICjt188,
	address As String
) As OperateResult(Of String())
```

```
public:
static OperateResult<array<String^>^>^ ReadStringArray(
	ICjt188^ cjt, 
	String^ address
)
```

```
static member ReadStringArray : 
        cjt : ICjt188 * 
        address : string -> OperateResult<string[]> 
```

#### 参数

cjt
:   类型：[HslCommunication.Instrument.CJT.HelperICjt188](3fabd38b-b4c3-f1b2-0469-ad7704214895.htm)  
    CJT通信对象

address
:   类型：SystemString  
    数据标识，具体需要查找手册来对应

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)String  
字符串数组信息

![](../icons/SectionExpanded.png)备注

地址可以携带地址域信息，例如 "s=2;90-1F" 或是 "s=100000;90-1F"

![](../icons/SectionExpanded.png)参见

#### 引用

[CJT188Helper 类](2e426e36-c009-46b5-43dd-b42e3a2b11c5.htm)

[HslCommunication.Instrument.CJT.Helper 命名空间](e7617fc1-440e-5d36-1b24-a85e8d4dffb0.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadValue(T) 方法 

[原文連結](http://api.hslcommunication.cn/html/06ba322b-7249-60c5-bd64-749fbbb3f6cb.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.CJT.Helper](../html/e7617fc1-440e-5d36-1b24-a85e8d4dffb0.htm "HslCommunication.Instrument.CJT.Helper")

[CJT188Helper 类](../html/2e426e36-c009-46b5-43dd-b42e3a2b11c5.htm "CJT188Helper 类")

[CJT188Helper 方法](../html/b6d542ae-0467-f6cb-eec2-8e1313c84576.htm "CJT188Helper 方法")

[AnalysisBytesAddress 方法](../html/19f104be-2349-b906-d8f3-bb6c6a9e5738.htm "AnalysisBytesAddress 方法 ")

[Build188EntireCommand 方法](../html/240f6ac6-815e-f3cf-0b0f-f928e65765a3.htm "Build188EntireCommand 方法 ")

[CheckResponse 方法](../html/c0df1d6d-197f-6039-ce25-8afb6366feea.htm "CheckResponse 方法 ")

[GetAddressByteFromString 方法](../html/88fb6eb9-c67b-28b8-6302-52517c8cc4bc.htm "GetAddressByteFromString 方法 ")

[Read 方法](../html/51e94d04-015b-a8eb-4522-f6fa925ea30c.htm "Read 方法 ")

[ReadAddress 方法](../html/a81b1674-c8a5-1e05-00e5-e69fd0a3b48b.htm "ReadAddress 方法 ")

[ReadStringArray 方法](../html/bc790e20-12da-1989-1a45-ed11fcc08909.htm "ReadStringArray 方法 ")

[ReadValue(T) 方法](../html/06ba322b-7249-60c5-bd64-749fbbb3f6cb.htm "ReadValue(T) 方法 ")

[Write 方法](../html/1d27c65a-d3db-8521-04af-d708904e71d8.htm "Write 方法 ")

[WriteAddress 方法](../html/f54f5a32-f807-edcf-f43b-9b07c29b840c.htm "WriteAddress 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CJT188HelperReadValueT 方法 |

读取数据的数组信息，需要指定如何从字符串转换的功能方法

**命名空间：**
 [HslCommunication.Instrument.CJT.Helper](e7617fc1-440e-5d36-1b24-a85e8d4dffb0.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult<T[]> ReadValue<T>(
	ICjt188 cjt,
	string address,
	ushort length,
	Func<string, T> trans
)
```

```
Public Shared Function ReadValue(Of T) ( 
	cjt As ICjt188,
	address As String,
	length As UShort,
	trans As Func(Of String, T)
) As OperateResult(Of T())
```

```
public:
generic<typename T>
static OperateResult<array<T>^>^ ReadValue(
	ICjt188^ cjt, 
	String^ address, 
	unsigned short length, 
	Func<String^, T>^ trans
)
```

```
static member ReadValue : 
        cjt : ICjt188 * 
        address : string * 
        length : uint16 * 
        trans : Func<string, 'T> -> OperateResult<'T[]> 
```

#### 参数

cjt
:   类型：[HslCommunication.Instrument.CJT.HelperICjt188](3fabd38b-b4c3-f1b2-0469-ad7704214895.htm)  
    CJT188的通信对象

address
:   类型：SystemString  
    地址信息

length
:   类型：SystemUInt16  
    数据长度

trans
:   类型：SystemFuncString, T  
    转换方法

#### 类型参数

T
:   类型信息

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)T  
包含泛型数组的结果对象

![](../icons/SectionExpanded.png)参见

#### 引用

[CJT188Helper 类](2e426e36-c009-46b5-43dd-b42e3a2b11c5.htm)

[HslCommunication.Instrument.CJT.Helper 命名空间](e7617fc1-440e-5d36-1b24-a85e8d4dffb0.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Write 方法 

[原文連結](http://api.hslcommunication.cn/html/1d27c65a-d3db-8521-04af-d708904e71d8.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.CJT.Helper](../html/e7617fc1-440e-5d36-1b24-a85e8d4dffb0.htm "HslCommunication.Instrument.CJT.Helper")

[CJT188Helper 类](../html/2e426e36-c009-46b5-43dd-b42e3a2b11c5.htm "CJT188Helper 类")

[CJT188Helper 方法](../html/b6d542ae-0467-f6cb-eec2-8e1313c84576.htm "CJT188Helper 方法")

[AnalysisBytesAddress 方法](../html/19f104be-2349-b906-d8f3-bb6c6a9e5738.htm "AnalysisBytesAddress 方法 ")

[Build188EntireCommand 方法](../html/240f6ac6-815e-f3cf-0b0f-f928e65765a3.htm "Build188EntireCommand 方法 ")

[CheckResponse 方法](../html/c0df1d6d-197f-6039-ce25-8afb6366feea.htm "CheckResponse 方法 ")

[GetAddressByteFromString 方法](../html/88fb6eb9-c67b-28b8-6302-52517c8cc4bc.htm "GetAddressByteFromString 方法 ")

[Read 方法](../html/51e94d04-015b-a8eb-4522-f6fa925ea30c.htm "Read 方法 ")

[ReadAddress 方法](../html/a81b1674-c8a5-1e05-00e5-e69fd0a3b48b.htm "ReadAddress 方法 ")

[ReadStringArray 方法](../html/bc790e20-12da-1989-1a45-ed11fcc08909.htm "ReadStringArray 方法 ")

[ReadValue(T) 方法](../html/06ba322b-7249-60c5-bd64-749fbbb3f6cb.htm "ReadValue(T) 方法 ")

[Write 方法](../html/1d27c65a-d3db-8521-04af-d708904e71d8.htm "Write 方法 ")

[WriteAddress 方法](../html/f54f5a32-f807-edcf-f43b-9b07c29b840c.htm "WriteAddress 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CJT188HelperWrite 方法 |

根据指定的数据标识来写入相关的原始数据信息，地址标识根据手册来，从高位到地位，例如 90-1F，分割符可以任意特殊字符或是没有分隔符。  
Read the relevant original data information according to the specified data identifier. The address identifier is based on the manual,
from high to position, such as 90-1F. The separator can be any special character or no separator.

**命名空间：**
 [HslCommunication.Instrument.CJT.Helper](e7617fc1-440e-5d36-1b24-a85e8d4dffb0.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult Write(
	ICjt188 cjt,
	string address,
	byte[] value
)
```

```
Public Shared Function Write ( 
	cjt As ICjt188,
	address As String,
	value As Byte()
) As OperateResult
```

```
public:
static OperateResult^ Write(
	ICjt188^ cjt, 
	String^ address, 
	array<unsigned char>^ value
)
```

```
static member Write : 
        cjt : ICjt188 * 
        address : string * 
        value : byte[] -> OperateResult 
```

#### 参数

cjt
:   类型：[HslCommunication.Instrument.CJT.HelperICjt188](3fabd38b-b4c3-f1b2-0469-ad7704214895.htm)  
    CJT188的通信对象

address
:   类型：SystemString  
    地址信息

value
:   类型：SystemByte  
    写入的数据值

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否写入成功

![](../icons/SectionExpanded.png)备注

地址可以携带地址域信息，例如 "s=2;90-1F" 或是 "s=100000;90-1F"，关于数据域信息，需要查找手册，例如:00-01-00-00 表示： (当前)正向有功总电能  
注意：本命令必须与编程键配合使用

![](../icons/SectionExpanded.png)参见

#### 引用

[CJT188Helper 类](2e426e36-c009-46b5-43dd-b42e3a2b11c5.htm)

[HslCommunication.Instrument.CJT.Helper 命名空间](e7617fc1-440e-5d36-1b24-a85e8d4dffb0.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## WriteAddress 方法 

[原文連結](http://api.hslcommunication.cn/html/f54f5a32-f807-edcf-f43b-9b07c29b840c.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.CJT.Helper](../html/e7617fc1-440e-5d36-1b24-a85e8d4dffb0.htm "HslCommunication.Instrument.CJT.Helper")

[CJT188Helper 类](../html/2e426e36-c009-46b5-43dd-b42e3a2b11c5.htm "CJT188Helper 类")

[CJT188Helper 方法](../html/b6d542ae-0467-f6cb-eec2-8e1313c84576.htm "CJT188Helper 方法")

[AnalysisBytesAddress 方法](../html/19f104be-2349-b906-d8f3-bb6c6a9e5738.htm "AnalysisBytesAddress 方法 ")

[Build188EntireCommand 方法](../html/240f6ac6-815e-f3cf-0b0f-f928e65765a3.htm "Build188EntireCommand 方法 ")

[CheckResponse 方法](../html/c0df1d6d-197f-6039-ce25-8afb6366feea.htm "CheckResponse 方法 ")

[GetAddressByteFromString 方法](../html/88fb6eb9-c67b-28b8-6302-52517c8cc4bc.htm "GetAddressByteFromString 方法 ")

[Read 方法](../html/51e94d04-015b-a8eb-4522-f6fa925ea30c.htm "Read 方法 ")

[ReadAddress 方法](../html/a81b1674-c8a5-1e05-00e5-e69fd0a3b48b.htm "ReadAddress 方法 ")

[ReadStringArray 方法](../html/bc790e20-12da-1989-1a45-ed11fcc08909.htm "ReadStringArray 方法 ")

[ReadValue(T) 方法](../html/06ba322b-7249-60c5-bd64-749fbbb3f6cb.htm "ReadValue(T) 方法 ")

[Write 方法](../html/1d27c65a-d3db-8521-04af-d708904e71d8.htm "Write 方法 ")

[WriteAddress 方法](../html/f54f5a32-f807-edcf-f43b-9b07c29b840c.htm "WriteAddress 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CJT188HelperWriteAddress 方法 |

写入设备的地址域信息，仅支持点对点通讯的情况，需要指定地址域信息，例如：14910000729011  
Write the address domain information of the device, only support point-to-point communication,
you need to specify the address domain information, for example: 14910000729011

**命名空间：**
 [HslCommunication.Instrument.CJT.Helper](e7617fc1-440e-5d36-1b24-a85e8d4dffb0.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult WriteAddress(
	ICjt188 cjt,
	string address
)
```

```
Public Shared Function WriteAddress ( 
	cjt As ICjt188,
	address As String
) As OperateResult
```

```
public:
static OperateResult^ WriteAddress(
	ICjt188^ cjt, 
	String^ address
)
```

```
static member WriteAddress : 
        cjt : ICjt188 * 
        address : string -> OperateResult 
```

#### 参数

cjt
:   类型：[HslCommunication.Instrument.CJT.HelperICjt188](3fabd38b-b4c3-f1b2-0469-ad7704214895.htm)  
    CJT188的通信对象

address
:   类型：SystemString  
    等待写入的地址域

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否写入成功

![](../icons/SectionExpanded.png)参见

#### 引用

[CJT188Helper 类](2e426e36-c009-46b5-43dd-b42e3a2b11c5.htm)

[HslCommunication.Instrument.CJT.Helper 命名空间](e7617fc1-440e-5d36-1b24-a85e8d4dffb0.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## CJTControl 类

[原文連結](http://api.hslcommunication.cn/html/70e1d569-367b-1da4-fd16-2592bb57b9bf.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.CJT.Helper](../html/e7617fc1-440e-5d36-1b24-a85e8d4dffb0.htm "HslCommunication.Instrument.CJT.Helper")

[CJTControl 类](../html/70e1d569-367b-1da4-fd16-2592bb57b9bf.htm "CJTControl 类")

[CJTControl 构造函数](../html/ba354ca2-389f-80d8-2e15-f15a094ecd1b.htm "CJTControl 构造函数 ")

[CJTControl 方法](../html/62458b7a-1a1a-80bf-f2f6-cb4426794229.htm "CJTControl 方法")

[CJTControl 字段](../html/7b9ef95f-5630-1f27-54e3-09ee9eb9636e.htm "CJTControl 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CJTControl 类 |

控制码信息

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  HslCommunication.Instrument.CJT.HelperCJTControl

**命名空间：**
 [HslCommunication.Instrument.CJT.Helper](e7617fc1-440e-5d36-1b24-a85e8d4dffb0.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class CJTControl
```

```
Public Class CJTControl
```

```
public ref class CJTControl
```

```
type CJTControl =  class end
```

CJTControl 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [CJTControl](ba354ca2-389f-80d8-2e15-f15a094ecd1b.htm) | 初始化 CJTControl 类的一个新实例 |

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

![](../icons/SectionExpanded.png)字段

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共字段静态成员 | [ReadAddress](731d91c5-ac26-8dbc-be26-285e512291f8.htm) | 读地址 |
| 公共字段静态成员 | [ReadData](d89e83f1-d822-47c9-09c5-4ec396804da5.htm) | 读数据 |
| 公共字段静态成员 | [WriteAddress](b4bb3e84-c4b4-9afc-1943-e1b3e7cbd138.htm) | 写地址 |
| 公共字段静态成员 | [WriteData](42d8303c-f84c-289d-71b8-c2659f427763.htm) | 写数据 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.Instrument.CJT.Helper 命名空间](e7617fc1-440e-5d36-1b24-a85e8d4dffb0.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## CJTControl 构造函数 

[原文連結](http://api.hslcommunication.cn/html/ba354ca2-389f-80d8-2e15-f15a094ecd1b.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.CJT.Helper](../html/e7617fc1-440e-5d36-1b24-a85e8d4dffb0.htm "HslCommunication.Instrument.CJT.Helper")

[CJTControl 类](../html/70e1d569-367b-1da4-fd16-2592bb57b9bf.htm "CJTControl 类")

[CJTControl 构造函数](../html/ba354ca2-389f-80d8-2e15-f15a094ecd1b.htm "CJTControl 构造函数 ")

[CJTControl 方法](../html/62458b7a-1a1a-80bf-f2f6-cb4426794229.htm "CJTControl 方法")

[CJTControl 字段](../html/7b9ef95f-5630-1f27-54e3-09ee9eb9636e.htm "CJTControl 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CJTControl 构造函数 |

初始化 [CJTControl](70e1d569-367b-1da4-fd16-2592bb57b9bf.htm) 类的一个新实例

**命名空间：**
 [HslCommunication.Instrument.CJT.Helper](e7617fc1-440e-5d36-1b24-a85e8d4dffb0.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public CJTControl()
```

```
Public Sub New
```

```
public:
CJTControl()
```

```
new : unit -> CJTControl
```

![](../icons/SectionExpanded.png)参见

#### 引用

[CJTControl 类](70e1d569-367b-1da4-fd16-2592bb57b9bf.htm)

[HslCommunication.Instrument.CJT.Helper 命名空间](e7617fc1-440e-5d36-1b24-a85e8d4dffb0.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## CJTControl 方法

[原文連結](http://api.hslcommunication.cn/html/62458b7a-1a1a-80bf-f2f6-cb4426794229.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.CJT.Helper](../html/e7617fc1-440e-5d36-1b24-a85e8d4dffb0.htm "HslCommunication.Instrument.CJT.Helper")

[CJTControl 类](../html/70e1d569-367b-1da4-fd16-2592bb57b9bf.htm "CJTControl 类")

[CJTControl 构造函数](../html/ba354ca2-389f-80d8-2e15-f15a094ecd1b.htm "CJTControl 构造函数 ")

[CJTControl 方法](../html/62458b7a-1a1a-80bf-f2f6-cb4426794229.htm "CJTControl 方法")

[CJTControl 字段](../html/7b9ef95f-5630-1f27-54e3-09ee9eb9636e.htm "CJTControl 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CJTControl 方法 |

[CJTControl](70e1d569-367b-1da4-fd16-2592bb57b9bf.htm) 类型公开以下成员。

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

[CJTControl 类](70e1d569-367b-1da4-fd16-2592bb57b9bf.htm)

[HslCommunication.Instrument.CJT.Helper 命名空间](e7617fc1-440e-5d36-1b24-a85e8d4dffb0.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## CJTControl 字段

[原文連結](http://api.hslcommunication.cn/html/7b9ef95f-5630-1f27-54e3-09ee9eb9636e.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.CJT.Helper](../html/e7617fc1-440e-5d36-1b24-a85e8d4dffb0.htm "HslCommunication.Instrument.CJT.Helper")

[CJTControl 类](../html/70e1d569-367b-1da4-fd16-2592bb57b9bf.htm "CJTControl 类")

[CJTControl 字段](../html/7b9ef95f-5630-1f27-54e3-09ee9eb9636e.htm "CJTControl 字段")

[ReadAddress 字段](../html/731d91c5-ac26-8dbc-be26-285e512291f8.htm "ReadAddress 字段")

[ReadData 字段](../html/d89e83f1-d822-47c9-09c5-4ec396804da5.htm "ReadData 字段")

[WriteAddress 字段](../html/b4bb3e84-c4b4-9afc-1943-e1b3e7cbd138.htm "WriteAddress 字段")

[WriteData 字段](../html/42d8303c-f84c-289d-71b8-c2659f427763.htm "WriteData 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CJTControl 字段 |

[CJTControl](70e1d569-367b-1da4-fd16-2592bb57b9bf.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)字段

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共字段静态成员 | [ReadAddress](731d91c5-ac26-8dbc-be26-285e512291f8.htm) | 读地址 |
| 公共字段静态成员 | [ReadData](d89e83f1-d822-47c9-09c5-4ec396804da5.htm) | 读数据 |
| 公共字段静态成员 | [WriteAddress](b4bb3e84-c4b4-9afc-1943-e1b3e7cbd138.htm) | 写地址 |
| 公共字段静态成员 | [WriteData](42d8303c-f84c-289d-71b8-c2659f427763.htm) | 写数据 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[CJTControl 类](70e1d569-367b-1da4-fd16-2592bb57b9bf.htm)

[HslCommunication.Instrument.CJT.Helper 命名空间](e7617fc1-440e-5d36-1b24-a85e8d4dffb0.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadAddress 字段

[原文連結](http://api.hslcommunication.cn/html/731d91c5-ac26-8dbc-be26-285e512291f8.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.CJT.Helper](../html/e7617fc1-440e-5d36-1b24-a85e8d4dffb0.htm "HslCommunication.Instrument.CJT.Helper")

[CJTControl 类](../html/70e1d569-367b-1da4-fd16-2592bb57b9bf.htm "CJTControl 类")

[CJTControl 字段](../html/7b9ef95f-5630-1f27-54e3-09ee9eb9636e.htm "CJTControl 字段")

[ReadAddress 字段](../html/731d91c5-ac26-8dbc-be26-285e512291f8.htm "ReadAddress 字段")

[ReadData 字段](../html/d89e83f1-d822-47c9-09c5-4ec396804da5.htm "ReadData 字段")

[WriteAddress 字段](../html/b4bb3e84-c4b4-9afc-1943-e1b3e7cbd138.htm "WriteAddress 字段")

[WriteData 字段](../html/42d8303c-f84c-289d-71b8-c2659f427763.htm "WriteData 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CJTControlReadAddress 字段 |

读地址

**命名空间：**
 [HslCommunication.Instrument.CJT.Helper](e7617fc1-440e-5d36-1b24-a85e8d4dffb0.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public const byte ReadAddress = 3
```

```
Public Const ReadAddress As Byte = 3
```

```
public:
literal unsigned char ReadAddress = 3
```

```
static val mutable ReadAddress: byte
```

#### 字段值

类型：Byte

![](../icons/SectionExpanded.png)参见

#### 引用

[CJTControl 类](70e1d569-367b-1da4-fd16-2592bb57b9bf.htm)

[HslCommunication.Instrument.CJT.Helper 命名空间](e7617fc1-440e-5d36-1b24-a85e8d4dffb0.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadData 字段

[原文連結](http://api.hslcommunication.cn/html/d89e83f1-d822-47c9-09c5-4ec396804da5.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.CJT.Helper](../html/e7617fc1-440e-5d36-1b24-a85e8d4dffb0.htm "HslCommunication.Instrument.CJT.Helper")

[CJTControl 类](../html/70e1d569-367b-1da4-fd16-2592bb57b9bf.htm "CJTControl 类")

[CJTControl 字段](../html/7b9ef95f-5630-1f27-54e3-09ee9eb9636e.htm "CJTControl 字段")

[ReadAddress 字段](../html/731d91c5-ac26-8dbc-be26-285e512291f8.htm "ReadAddress 字段")

[ReadData 字段](../html/d89e83f1-d822-47c9-09c5-4ec396804da5.htm "ReadData 字段")

[WriteAddress 字段](../html/b4bb3e84-c4b4-9afc-1943-e1b3e7cbd138.htm "WriteAddress 字段")

[WriteData 字段](../html/42d8303c-f84c-289d-71b8-c2659f427763.htm "WriteData 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CJTControlReadData 字段 |

读数据

**命名空间：**
 [HslCommunication.Instrument.CJT.Helper](e7617fc1-440e-5d36-1b24-a85e8d4dffb0.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public const byte ReadData = 1
```

```
Public Const ReadData As Byte = 1
```

```
public:
literal unsigned char ReadData = 1
```

```
static val mutable ReadData: byte
```

#### 字段值

类型：Byte

![](../icons/SectionExpanded.png)参见

#### 引用

[CJTControl 类](70e1d569-367b-1da4-fd16-2592bb57b9bf.htm)

[HslCommunication.Instrument.CJT.Helper 命名空间](e7617fc1-440e-5d36-1b24-a85e8d4dffb0.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## WriteAddress 字段

[原文連結](http://api.hslcommunication.cn/html/b4bb3e84-c4b4-9afc-1943-e1b3e7cbd138.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.CJT.Helper](../html/e7617fc1-440e-5d36-1b24-a85e8d4dffb0.htm "HslCommunication.Instrument.CJT.Helper")

[CJTControl 类](../html/70e1d569-367b-1da4-fd16-2592bb57b9bf.htm "CJTControl 类")

[CJTControl 字段](../html/7b9ef95f-5630-1f27-54e3-09ee9eb9636e.htm "CJTControl 字段")

[ReadAddress 字段](../html/731d91c5-ac26-8dbc-be26-285e512291f8.htm "ReadAddress 字段")

[ReadData 字段](../html/d89e83f1-d822-47c9-09c5-4ec396804da5.htm "ReadData 字段")

[WriteAddress 字段](../html/b4bb3e84-c4b4-9afc-1943-e1b3e7cbd138.htm "WriteAddress 字段")

[WriteData 字段](../html/42d8303c-f84c-289d-71b8-c2659f427763.htm "WriteData 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CJTControlWriteAddress 字段 |

写地址

**命名空间：**
 [HslCommunication.Instrument.CJT.Helper](e7617fc1-440e-5d36-1b24-a85e8d4dffb0.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public const byte WriteAddress = 21
```

```
Public Const WriteAddress As Byte = 21
```

```
public:
literal unsigned char WriteAddress = 21
```

```
static val mutable WriteAddress: byte
```

#### 字段值

类型：Byte

![](../icons/SectionExpanded.png)参见

#### 引用

[CJTControl 类](70e1d569-367b-1da4-fd16-2592bb57b9bf.htm)

[HslCommunication.Instrument.CJT.Helper 命名空间](e7617fc1-440e-5d36-1b24-a85e8d4dffb0.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## WriteData 字段

[原文連結](http://api.hslcommunication.cn/html/42d8303c-f84c-289d-71b8-c2659f427763.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.CJT.Helper](../html/e7617fc1-440e-5d36-1b24-a85e8d4dffb0.htm "HslCommunication.Instrument.CJT.Helper")

[CJTControl 类](../html/70e1d569-367b-1da4-fd16-2592bb57b9bf.htm "CJTControl 类")

[CJTControl 字段](../html/7b9ef95f-5630-1f27-54e3-09ee9eb9636e.htm "CJTControl 字段")

[ReadAddress 字段](../html/731d91c5-ac26-8dbc-be26-285e512291f8.htm "ReadAddress 字段")

[ReadData 字段](../html/d89e83f1-d822-47c9-09c5-4ec396804da5.htm "ReadData 字段")

[WriteAddress 字段](../html/b4bb3e84-c4b4-9afc-1943-e1b3e7cbd138.htm "WriteAddress 字段")

[WriteData 字段](../html/42d8303c-f84c-289d-71b8-c2659f427763.htm "WriteData 字段")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| CJTControlWriteData 字段 |

写数据

**命名空间：**
 [HslCommunication.Instrument.CJT.Helper](e7617fc1-440e-5d36-1b24-a85e8d4dffb0.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public const byte WriteData = 4
```

```
Public Const WriteData As Byte = 4
```

```
public:
literal unsigned char WriteData = 4
```

```
static val mutable WriteData: byte
```

#### 字段值

类型：Byte

![](../icons/SectionExpanded.png)参见

#### 引用

[CJTControl 类](70e1d569-367b-1da4-fd16-2592bb57b9bf.htm)

[HslCommunication.Instrument.CJT.Helper 命名空间](e7617fc1-440e-5d36-1b24-a85e8d4dffb0.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ICjt188 接口

[原文連結](http://api.hslcommunication.cn/html/3fabd38b-b4c3-f1b2-0469-ad7704214895.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.CJT.Helper](../html/e7617fc1-440e-5d36-1b24-a85e8d4dffb0.htm "HslCommunication.Instrument.CJT.Helper")

[ICjt188 接口](../html/3fabd38b-b4c3-f1b2-0469-ad7704214895.htm "ICjt188 接口")

[ICjt188 属性](../html/7dbdd1e5-a59c-0253-ecc4-dea40ae4633a.htm "ICjt188 属性")

[ICjt188 方法](../html/44dc2683-e490-896d-05a0-cee4a0920cea.htm "ICjt188 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ICjt188 接口 |

CJT188设备的接口

**命名空间：**
 [HslCommunication.Instrument.CJT.Helper](e7617fc1-440e-5d36-1b24-a85e8d4dffb0.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public interface ICjt188 : IReadWriteDevice, 
	IReadWriteNet
```

```
Public Interface ICjt188
	Inherits IReadWriteDevice, IReadWriteNet
```

```
public interface class ICjt188 : IReadWriteDevice, 
	IReadWriteNet
```

```
type ICjt188 =  
    interface
        interface IReadWriteDevice
        interface IReadWriteNet
    end
```

ICjt188 类型公开以下成员。

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性代码示例 | [ByteTransform](0836ec33-d464-ed39-37ff-b2ada0d286fc.htm) | 当前的数据变换机制，当你需要从字节数据转换类型数据的时候需要。  The current data transformation mechanism is required when you need to convert type data from byte data. (继承自 [IReadWriteDevice](d987afcc-f99c-6aed-3847-5b1890504635.htm)。) |
| 公共属性 | [ConnectionId](20787592-e3ec-e0d0-27c2-abfa24654bd2.htm) | 当前连接的唯一ID号，默认为长度20的guid码加随机数组成，方便列表管理，也可以自己指定  The unique ID number of the current connection. The default is a 20-digit guid code plus a random number. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共属性 | [EnableCodeFE](2420a2be-c28d-4ddf-63e7-05e2eb050af8.htm) | 获取或设置是否在每一次的报文通信时，增加"FE FE"的命令头  Get or set whether to add the command header of "FE FE" in each message communication |
| 公共属性 | [InstrumentType](54d205f6-517b-1b53-bb21-afa097006156.htm) | 仪表的类型 |
| 公共属性代码示例 | [LogNet](38760c1f-22e3-f8fe-7512-15849a4b6eb0.htm) | 组件的日志工具，支持日志记录，只要实例化后，当前网络的基本信息，就以[DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm)等级进行输出  The component's logging tool supports logging. As long as the instantiation of the basic network information, the output will be output at [DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm) (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共属性 | [Station](36f1a758-e626-a565-2d13-3a0cdea61613.htm) | 获取或设置当前的地址域信息，是一个12个字符的BCD码，例如：149100007290  Get or set the current address domain information, which is a 12-character BCD code, for example: 149100007290 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [ActiveDeveice](f00c3a89-859f-c8ba-5c0a-b58cd4162eaa.htm) | 激活设备的命令，只发送数据到设备，不等待设备数据返回  The command to activate the device, only send data to the device, do not wait for the device data to return |
| 公共方法 | [Read(String, UInt16)](b8d4187f-ab9a-79e9-c7ca-8a01ec0df233.htm) | 批量读取字节数组信息，需要指定地址和长度，返回原始的字节数组  Batch read byte array information, need to specify the address and length, return the original byte array (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadT](1be4d687-47e0-f01f-2969-570e5b057c4a.htm) | 读取支持Hsl特性的数据内容，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明，或是官网的详细文档  Read the data content of the Hsl attribute. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [ReadAddress](a669de5c-80f0-74df-235c-2eccb0ff9c38.htm) | 读取设备的通信地址，仅支持点对点通讯的情况，返回地址域数据，例如：14910000729012  Read the communication address of the device, only support point-to-point communication, and return the address field data, for example: 149100007290 |
| 公共方法 | [ReadAsync(String, UInt16)](59839c27-5a76-2f53-730c-c9d2bb41c2c4.htm) | 异步批量读取字节数组信息，需要指定地址和长度，返回原始的字节数组  Asynchronous batch read byte array information, need to specify the address and length, return the original byte array (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadAsyncT](2c18ca3b-cee6-e5ac-8962-48d1f7a99ff4.htm) | 异步读取支持Hsl特性的数据内容，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Asynchronously read the data content of the Hsl attribute. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [ReadBool(String)](67c3a6fd-ad81-62bf-a72d-74b99522a93e.htm) | 读取单个的Boolean数据信息  Read a single Boolean data message (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [ReadBool(String, UInt16)](16e00e01-28d8-bfb7-e3bb-9d8b8778d7b9.htm) | 批量读取Boolean数组信息，需要指定地址和长度，返回Boolean 数组  Batch read Boolean array information, need to specify the address and length, return Boolean array (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [ReadBoolAsync(String)](f99fec07-3c4c-82a3-4427-4cc26281c871.htm) | 异步读取单个的Boolean数据信息  Asynchronously read a single Boolean data message (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [ReadBoolAsync(String, UInt16)](74e40b23-198a-944b-6ed7-3f58ca51da5d.htm) | 异步批量读取Boolean数组信息，需要指定地址和长度，返回Boolean 数组  Asynchronously batch read Boolean array information, need to specify the address and length, return Boolean array (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadCustomerT(String)](04bd0f87-9d33-c4d1-9420-7ca995687eb3.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，返回一个新的类型的实例对象。  To read a custom data type, you need to inherit from the IDataTransfer interface and return an instance object of a new type. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadCustomerT(String, T)](72f8e1a8-fd49-efc7-3485-fb2684f1d948.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，传入一个实例，对这个实例进行赋值，并返回该实例的对象。  To read a custom data type, you need to inherit from the IDataTransfer interface, pass in an instance, assign a value to this instance, and return the object of the instance. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadCustomerAsyncT(String)](482ef33f-8bb1-bea7-2bad-5c23bafcd173.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，返回一个新的类型的实例对象。  To read a custom data type, you need to inherit from the IDataTransfer interface and return an instance object of a new type. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadCustomerAsyncT(String, T)](0405df60-0745-1be6-0101-59edaf4e4d38.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，传入一个实例，对这个实例进行赋值，并返回该实例的对象。  To read a custom data type, you need to inherit from the IDataTransfer interface, pass in an instance, assign a value to this instance, and return the object of the instance. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadDouble(String)](18936f56-2b3c-59cc-788b-63490168e4e2.htm) | 读取双浮点的数据  Read double floating point data (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadDouble(String, UInt16)](a8468555-105b-337a-196f-b0a667856edd.htm) | 读取双浮点数据的数组  Read double floating point data array (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadDoubleAsync(String)](d902471c-124a-67a2-be35-e2feebbb1298.htm) | 异步读取双浮点的数据  Asynchronously read double floating point data (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadDoubleAsync(String, UInt16)](3a436322-c7f5-bb1d-e97b-dc6c0d6990d2.htm) | 异步读取双浮点数据的数组  Asynchronously read double floating point data array (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadFloat(String)](95ca7390-d065-ef28-cfdb-e8c161437ecd.htm) | 读取单浮点数据  Read single floating point data (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadFloat(String, UInt16)](a8ae302d-6269-8baa-cacc-f35df9fcd131.htm) | 读取单浮点精度的数组  Read single floating point array (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadFloatAsync(String)](de191b71-b9e7-6f9c-e7a6-9a117e13fb05.htm) | 异步读取单浮点数据  Asynchronously read single floating point data (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadFloatAsync(String, UInt16)](26a13b36-422b-0520-3cb0-66a235e3d7a8.htm) | 异步读取单浮点精度的数组  Asynchronously read single floating point array (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [ReadFromCoreServer(Byte)](413a68ef-2f4b-638f-8bad-973673920de4.htm) | 将当前的数据报文发送到设备去，具体使用什么通信方式取决于设备信息，然后从设备接收数据回来，并返回给调用者。  Send the current data message to the device, the specific communication method used depends on the device information, and then receive the data back from the device and return it to the caller. (继承自 [IReadWriteDevice](d987afcc-f99c-6aed-3847-5b1890504635.htm)。) |
| 公共方法 | [ReadFromCoreServer(IEnumerableByte)](236bfe99-f977-a092-d283-c313ab8aaa7c.htm) | 将多个数据报文按顺序发到设备，并从设备接收返回的数据内容，然后拼接成一个Byte[]信息，需要重写[UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm)方法才能返回正确的结果。  Send multiple data packets to the device in sequence, and receive the returned data content from the device, and then splicing them into a Byte[] message, you need to rewrite [UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm) method to return the correct result. (继承自 [IReadWriteDevice](d987afcc-f99c-6aed-3847-5b1890504635.htm)。) |
| 公共方法 | [ReadFromCoreServerAsync(Byte)](646de5ac-9ca6-1fec-f9de-a318dbc34dc1.htm) | 将当前的数据报文发送到设备去，具体使用什么通信方式取决于设备信息，然后从设备接收数据回来，并返回给调用者。  Send the current data message to the device, the specific communication method used depends on the device information, and then receive the data back from the device and return it to the caller. (继承自 [IReadWriteDevice](d987afcc-f99c-6aed-3847-5b1890504635.htm)。) |
| 公共方法 | [ReadFromCoreServerAsync(IEnumerableByte)](90bb695c-5fc3-eb72-39f6-195819857394.htm) | 将多个数据报文按顺序发到设备，并从设备接收返回的数据内容，然后拼接成一个Byte[]信息，需要重写[UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm)方法才能返回正确的结果。  Send multiple data packets to the device in sequence, and receive the returned data content from the device, and then splicing them into a Byte[] message, you need to rewrite [UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm) method to return the correct result. (继承自 [IReadWriteDevice](d987afcc-f99c-6aed-3847-5b1890504635.htm)。) |
| 公共方法代码示例 | [ReadInt16(String)](c5308422-4397-2d19-f2de-6d4213394e3d.htm) | 读取16位的有符号的整型数据  Read 16-bit signed integer data (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadInt16(String, UInt16)](3173b3fb-e0ac-1b58-ffff-3741aa992458.htm) | 读取16位的有符号整型数组  Read 16-bit signed integer array (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadInt16Async(String)](bef4e722-646e-3196-0fc7-bf12cfa7eb57.htm) | 异步读取16位的有符号的整型数据  Asynchronously read 16-bit signed integer data (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadInt16Async(String, UInt16)](1d194744-a756-5f40-90a8-ec59633a3970.htm) | 异步读取16位的有符号整型数组  Asynchronously read 16-bit signed integer array (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadInt32(String)](51ba7f65-d19b-20dd-e81e-f3913226b5cc.htm) | 读取32位的有符号整型  Read 32-bit signed integer (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadInt32(String, UInt16)](39024553-e8a3-ceb6-946c-0812c211c218.htm) | 读取32位有符号整型数组  Read 32-bit signed integer array (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadInt32Async(String)](fcf027ad-23d6-a6a4-c715-b0c1dcf18b31.htm) | 异步读取32位的有符号整型  Asynchronously read 32-bit signed integer (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadInt32Async(String, UInt16)](f873bfc0-32e5-3303-ec0a-b0c5515754c6.htm) | 异步读取32位有符号整型数组  Asynchronously read 32-bit signed integer array (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadInt64(String)](a3bea048-fd83-2158-4580-231419212efe.htm) | 读取64位的有符号整型  Read 64-bit signed integer (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadInt64(String, UInt16)](32190a15-c77a-7831-9d59-f2419a7f8bb2.htm) | 读取64位的有符号整型数组  Read 64-bit signed integer array (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadInt64Async(String)](d2220dae-6131-a6aa-8533-a5cce335e2c4.htm) | 异步读取64位的有符号整型  Asynchronously read 64-bit signed integer (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadInt64Async(String, UInt16)](b511a119-ee98-1abc-b578-edcac930fcd4.htm) | 异步读取64位的有符号整型数组  Asynchronously read 64-bit signed integer array (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadString(String, UInt16)](a8ca8b30-6612-307e-c046-d8ff01b30df1.htm) | 读取字符串数据，默认为最常见的ASCII编码  Read string data, default is the most common ASCII encoding (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadString(String, UInt16, Encoding)](fc289def-e53b-b47d-22fe-7e02b5e354e0.htm) | 使用指定的编码，读取字符串数据  Reads string data using the specified encoding (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [ReadStringArray](226fd33b-1900-9236-7bb7-ef5b195d6101.htm) | 读取指定地址的所有的字符串数据信息，一般来说，一个地址只有一个数据 |
| 公共方法代码示例 | [ReadStringAsync(String, UInt16)](960c5f9b-c87b-d7e9-a4e0-048e4f2f5981.htm) | 异步读取字符串数据，默认为最常见的ASCII编码  Asynchronously read string data, default is the most common ASCII encoding (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadStringAsync(String, UInt16, Encoding)](db06baf4-e674-22de-e89e-cb26646413f8.htm) | 异步使用指定的编码，读取字符串数据  Asynchronously reads string data using the specified encoding (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadStructT](50df0eb8-8f9c-ad9c-4dc9-bf75b4438d46.htm) | 读取结构体类型的数据，根据结构体自身的定义，读取原始字节数组，然后解析出实际的结构体数据，结构体需要实现[HslStructAttribute](87e427a8-0c7b-43a8-c421-c27264eda90f.htm)特性 (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadStructAsyncT](af626b7e-2adb-ecd1-f2a0-82aaa6e8b68d.htm) | 读取结构体类型的数据，根据结构体自身的定义，读取原始字节数组，然后解析出实际的结构体数据，结构体需要实现[HslStructAttribute](87e427a8-0c7b-43a8-c421-c27264eda90f.htm)特性 (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadUInt16(String)](66238795-7948-2cf9-c639-d1bb3da126dd.htm) | 读取16位的无符号整型  Read 16-bit unsigned integer (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadUInt16(String, UInt16)](28888661-1fcc-f941-5746-9124eb539fe8.htm) | 读取16位的无符号整型数组  Read 16-bit unsigned integer array (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadUInt16Async(String)](1744a25b-c440-7d1e-dee5-606b0ff485bc.htm) | 异步读取16位的无符号整型  Asynchronously read 16-bit unsigned integer (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadUInt16Async(String, UInt16)](0748f071-d6a7-581b-ed7a-d73317b16762.htm) | 异步读取16位的无符号整型数组  Asynchronously read 16-bit unsigned integer array (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadUInt32(String)](b252ad01-4035-bc3e-e297-7fd97c3ac7af.htm) | 读取32位的无符号整型  Read 32-bit unsigned integer (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadUInt32(String, UInt16)](64d8381b-05b6-a990-c03e-2c8a1b7cfb7f.htm) | 读取32位的无符号整型数组  Read 32-bit unsigned integer array (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadUInt32Async(String)](a5b848e0-82bd-2bb5-2dd9-6634063e7b48.htm) | 异步读取32位的无符号整型  Asynchronously read 32-bit unsigned integer (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadUInt32Async(String, UInt16)](59ba0cbe-08f8-5a30-6029-0a9f718c7d76.htm) | 异步读取32位的无符号整型数组  Asynchronously read 32-bit unsigned integer array (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadUInt64(String)](529c82ba-1d5e-f33b-2eb9-c862799f952f.htm) | 读取64位的无符号整型  Read 64-bit unsigned integer (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadUInt64(String, UInt16)](2af999cc-c8d4-e25a-68f1-d621e5d84af7.htm) | 读取64位的无符号整型的数组  Read 64-bit unsigned integer array (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadUInt64Async(String)](0c521a11-248c-ba7a-fc27-0e5a645611b5.htm) | 异步读取64位的无符号整型  Asynchronously read 64-bit unsigned integer (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadUInt64Async(String, UInt16)](b510070f-beb5-959b-7c6b-20f212a68dd1.htm) | 异步读取64位的无符号整型的数组  Asynchronously read 64-bit unsigned integer array (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [Wait(String, Boolean, Int32, Int32)](c8c01cf3-bf99-761c-23d4-82dd54081952.htm) | 等待指定地址的Boolean值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the Boolean value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [Wait(String, Int16, Int32, Int32)](f8a820b0-cf90-668c-6fba-524bd0659042.htm) | 等待指定地址的Int16值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the Int16 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [Wait(String, Int32, Int32, Int32)](96ad6909-b660-2894-9c17-a7f4f4d4a2bd.htm) | 等待指定地址的Int32值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the Int32 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [Wait(String, Int64, Int32, Int32)](28ba4678-2aeb-7502-95e8-d6b855bbe971.htm) | 等待指定地址的Int64值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the Int64 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [Wait(String, UInt16, Int32, Int32)](44af3173-f7f4-04f5-753b-7855581a171a.htm) | 等待指定地址的UInt16值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the UInt16 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [Wait(String, UInt32, Int32, Int32)](0e877e55-9b76-f662-6841-9102a9e9b9b8.htm) | 等待指定地址的UInt32值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the UInt32 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [Wait(String, UInt64, Int32, Int32)](6801831d-c349-800d-e3d3-34d2fde09e0f.htm) | 等待指定地址的UInt64值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the UInt64 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [WaitAsync(String, Boolean, Int32, Int32)](b51daff5-9c09-c614-9ea3-f14e907a36b8.htm) | 等待指定地址的Boolean值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the Boolean value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [WaitAsync(String, Int16, Int32, Int32)](cfd8af62-371b-84e9-5695-87539fbb204b.htm) | 等待指定地址的Int16值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the Int16 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [WaitAsync(String, Int32, Int32, Int32)](639f2d1b-9a41-0cf2-0335-7f9641fbb0c8.htm) | 等待指定地址的Int32值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the Int32 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [WaitAsync(String, Int64, Int32, Int32)](4355bcb4-e923-9226-ac82-55dd1325f34f.htm) | 等待指定地址的Int64值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the Int64 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [WaitAsync(String, UInt16, Int32, Int32)](291b6f29-9cbf-6fc3-4828-226f412e9cc6.htm) | 等待指定地址的UInt16值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the UInt16 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [WaitAsync(String, UInt32, Int32, Int32)](e1009c58-9f02-660e-225a-b6861eb8781f.htm) | 等待指定地址的UInt32值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the UInt32 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [WaitAsync(String, UInt64, Int32, Int32)](4236aa40-6c9d-2a0c-8470-c55dab497074.htm) | 等待指定地址的UInt64值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the UInt64 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [Write(String, Boolean)](4d742738-6864-465c-5382-bd49032eae37.htm) | 写入单个的Boolean数据，返回是否成功  Write a single Boolean data, and return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [Write(String, Boolean)](2a99122c-1496-d766-f0aa-44a18918fcd3.htm) | 批量写入Boolean数组数据，返回是否成功  Batch write Boolean array data, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [Write(String, Byte)](b8d8d94e-300d-c7b7-2e80-6859c62479fb.htm) | 写入原始的byte数组数据到指定的地址，返回是否写入成功  Write the original byte array data to the specified address, and return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [Write(String, Double)](34e1c967-7214-bb8e-3daf-fdb93fc1db59.htm) | 写入double数据，返回是否成功  Write double data, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [Write(String, Double)](f2a8cc81-ee85-4d3a-1eaf-af33b15de5e6.htm) | 写入double数组，返回是否成功  Write double array, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [Write(String, Int16)](53b27553-763c-2dd9-b8a1-df2ebf68ab2c.htm) | 写入short数据，返回是否成功  Write short data, returns whether success (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [Write(String, Int16)](0b8fff2d-0da9-d73a-bc11-866a501e6112.htm) | 写入short数组，返回是否成功  Write short array, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [Write(String, Int32)](281ce6f6-5c53-da62-9a09-9d439390a972.htm) | 写入int数据，返回是否成功  Write int data, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [Write(String, Int32)](496e02c2-62f2-4f72-aff1-1d691cc87726.htm) | 写入int[]数组，返回是否成功  Write int array, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [Write(String, Int64)](7ec919a0-1d6b-b7d7-1f2f-eb1644beb912.htm) | 写入long数据，返回是否成功  Write long data, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [Write(String, Int64)](026100a5-d52e-b640-f239-a140c903201d.htm) | 写入long数组，返回是否成功  Write long array, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [Write(String, Single)](7c72c919-367f-9bc6-567c-550813459297.htm) | 写入float数据，返回是否成功  Write float data, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [Write(String, Single)](195c858e-1e2d-863a-8a2c-dc49cb2063d3.htm) | 写入float数组，返回是否成功  Write float array, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [Write(String, String)](fc988c8f-b2d6-18f9-3274-9b57d4a04e1c.htm) | 写入字符串信息，编码为ASCII  Write string information, encoded as ASCII (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [Write(String, UInt16)](7084bc37-6bda-e506-4118-74aff979efff.htm) | 写入ushort数据，返回是否成功  Write ushort data, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [Write(String, UInt16)](9c21745f-7c5f-cae3-7f2a-2bcf260e54a9.htm) | 写入ushort数组，返回是否成功  Write ushort array, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [Write(String, UInt32)](ed79fea9-bc6e-3ef3-5294-134b872141cd.htm) | 写入uint数据，返回是否成功  Write uint data, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [Write(String, UInt32)](3ae38bac-66cc-5cdf-be68-483362b17ac3.htm) | 写入uint[]数组，返回是否成功  Write uint array, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [Write(String, UInt64)](3927095f-0705-e823-bb24-74c9363fab36.htm) | 写入ulong数据，返回是否成功  Write ulong data, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [Write(String, UInt64)](2eb4a3ce-c350-9949-405b-ccb6b5d9786e.htm) | 写入ulong数组，返回是否成功  Write ulong array, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [Write(String, String, Int32)](1ed07021-8f14-743a-071c-c348025634f5.htm) | 写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为ASCII  Write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, it is filled with 0 and the encoding is ASCII. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [Write(String, String, Encoding)](c66a9dfb-840a-3ac4-96ef-48e23fb85ec1.htm) | 写入字符串信息，需要指定的编码信息  Write string information, need to specify the encoding information (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [Write(String, String, Int32, Encoding)](61396f95-43cd-34fb-9055-37676c4e436d.htm) | 写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为指定的编码信息  Write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, then the operation is complemented with 0 , you should specified the encoding information (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteT(T)](ec80557a-6d74-6030-7e09-57c0c6a63cc4.htm) | 写入支持Hsl特性的数据，返回是否写入成功，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Write data that supports the Hsl attribute, and return whether the write was successful. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [WriteAddress](dfdff82a-e574-bf9a-381a-471bfb3c0300.htm) | 写入设备的地址域信息，仅支持点对点通讯的情况，需要指定地址域信息，例如：14910000729011  Write the address domain information of the device, only support point-to-point communication, you need to specify the address domain information, for example: 14910000729011 |
| 公共方法 | [WriteAsync(String, Boolean)](559924da-b298-9c81-794f-83d53aafc602.htm) | 异步批量写入Boolean数组数据，返回是否成功  Asynchronously batch write Boolean array data, return success (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [WriteAsync(String, Boolean)](cf8f8e56-df3d-607b-1c52-c1a76378db39.htm) | 异步批量写入Boolean数组数据，返回是否成功  Asynchronously batch write Boolean array data, return success (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Byte)](1a2e83c8-7782-847a-5d89-08d5556ac4ad.htm) | 异步写入原始的byte数组数据到指定的地址，返回是否写入成功  Asynchronously writes the original byte array data to the specified address, and returns whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Double)](ca502ab0-a05c-9dc7-c8b0-9afd3a24cb63.htm) | 异步写入double数据，返回是否成功  Asynchronously write double data, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Double)](33c58466-53fe-e221-1d73-6650c434e2b1.htm) | 异步写入double数组，返回是否成功  Asynchronously write double array, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int16)](3e923476-059a-0de9-4fc5-6df5dc7a93a5.htm) | 异步写入short数据，返回是否成功  Asynchronously write short data, returns whether success (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int16)](8955c5b3-7709-57ce-51fa-a58369f1147f.htm) | 异步写入short数组，返回是否成功  Asynchronously write short array, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int32)](448b905c-dfee-7f97-4dd9-8caf3539eef4.htm) | 异步写入int数据，返回是否成功  Asynchronously write int data, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int32)](cd5f66ad-6ffc-5394-4554-43abb07d977a.htm) | 异步写入int[]数组，返回是否成功  Asynchronously write int array, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int64)](972d6a25-bb1a-3c92-c12a-4b6942fc026f.htm) | 异步写入long数据，返回是否成功  Asynchronously write long data, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int64)](89e51562-f954-65b5-0c20-d1602d8e78a8.htm) | 异步写入long数组，返回是否成功  Asynchronously write long array, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Single)](213794a0-4dc8-36b4-3b3f-d1ee74e35155.htm) | 异步写入float数据，返回是否成功  Asynchronously write float data, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Single)](dca9a471-a3cb-3242-b3fe-edeb04d4f786.htm) | 异步写入float数组，返回是否成功  Asynchronously write float array, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, String)](265dfc0f-8630-cc0a-4a25-cd31d50f6056.htm) | 异步写入字符串信息，编码为ASCII  Asynchronously write string information, encoded as ASCII (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt16)](ee35584e-d8c8-720e-7a05-24a2d3011021.htm) | 异步写入ushort数据，返回是否成功  Asynchronously write ushort data, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt16)](9df3d402-f0c7-5888-65d0-f527a6ada1b2.htm) | 异步写入ushort数组，返回是否成功  Asynchronously write ushort array, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt32)](76ac4f61-3a88-384c-fae1-857066a74ac0.htm) | 异步写入uint数据，返回是否成功  Asynchronously write uint data, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt32)](aae8db84-1d63-31a3-e84c-f03a5f784ff0.htm) | 异步写入uint[]数组，返回是否成功  Asynchronously write uint array, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt64)](81bbb9e0-b222-067d-7dd9-a879277859c4.htm) | 异步写入ulong数据，返回是否成功  Asynchronously write ulong data, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt64)](f56c270e-dc8d-9adc-ea60-36d47e2c6e50.htm) | 异步写入ulong数组，返回是否成功  Asynchronously write ulong array, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [WriteAsync(String, String, Int32)](c9efea95-68b3-ab09-5176-d388576844e8.htm) | 异步写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为ASCII  Asynchronously write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, it is filled with 0 and the encoding is ASCII. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, String, Encoding)](cda78ec9-ab0c-b6f1-db38-0e51e6aa06d0.htm) | 异步写入字符串信息，需要指定的编码信息  Asynchronously write string information, need to specify the encoding information (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [WriteAsync(String, String, Int32, Encoding)](c077e3ad-d170-aea3-0058-2888f2b27126.htm) | 异步写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为指定的编码信息  Asynchronously write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, then the operation is complemented with 0 , you should specified the encoding information (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteAsyncT(T)](da0f692b-4d73-0520-6e17-7890b6ef51cb.htm) | 异步写入支持Hsl特性的数据，返回是否写入成功，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Asynchronously write data that supports the Hsl attribute, and return whether the write was successful. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteCustomerT](466c49b4-95ac-148b-c18d-a419c5f6dab3.htm) | 写入自定义类型的数据，该类型必须继承自IDataTransfer接口  Write data of a custom type, which must inherit from the IDataTransfer interface (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteCustomerAsyncT](82aa8f35-3c06-d4f6-e8c5-affb6ab6ce41.htm) | 写入自定义类型的数据，该类型必须继承自IDataTransfer接口  Write data of a custom type, which must inherit from the IDataTransfer interface (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.Instrument.CJT.Helper 命名空间](e7617fc1-440e-5d36-1b24-a85e8d4dffb0.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ICjt188 属性

[原文連結](http://api.hslcommunication.cn/html/7dbdd1e5-a59c-0253-ecc4-dea40ae4633a.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.CJT.Helper](../html/e7617fc1-440e-5d36-1b24-a85e8d4dffb0.htm "HslCommunication.Instrument.CJT.Helper")

[ICjt188 接口](../html/3fabd38b-b4c3-f1b2-0469-ad7704214895.htm "ICjt188 接口")

[ICjt188 属性](../html/7dbdd1e5-a59c-0253-ecc4-dea40ae4633a.htm "ICjt188 属性")

[EnableCodeFE 属性](../html/2420a2be-c28d-4ddf-63e7-05e2eb050af8.htm "EnableCodeFE 属性 ")

[InstrumentType 属性](../html/54d205f6-517b-1b53-bb21-afa097006156.htm "InstrumentType 属性 ")

[Station 属性](../html/36f1a758-e626-a565-2d13-3a0cdea61613.htm "Station 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ICjt188 属性 |

[ICjt188](3fabd38b-b4c3-f1b2-0469-ad7704214895.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性代码示例 | [ByteTransform](0836ec33-d464-ed39-37ff-b2ada0d286fc.htm) | 当前的数据变换机制，当你需要从字节数据转换类型数据的时候需要。  The current data transformation mechanism is required when you need to convert type data from byte data. (继承自 [IReadWriteDevice](d987afcc-f99c-6aed-3847-5b1890504635.htm)。) |
| 公共属性 | [ConnectionId](20787592-e3ec-e0d0-27c2-abfa24654bd2.htm) | 当前连接的唯一ID号，默认为长度20的guid码加随机数组成，方便列表管理，也可以自己指定  The unique ID number of the current connection. The default is a 20-digit guid code plus a random number. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共属性 | [EnableCodeFE](2420a2be-c28d-4ddf-63e7-05e2eb050af8.htm) | 获取或设置是否在每一次的报文通信时，增加"FE FE"的命令头  Get or set whether to add the command header of "FE FE" in each message communication |
| 公共属性 | [InstrumentType](54d205f6-517b-1b53-bb21-afa097006156.htm) | 仪表的类型 |
| 公共属性代码示例 | [LogNet](38760c1f-22e3-f8fe-7512-15849a4b6eb0.htm) | 组件的日志工具，支持日志记录，只要实例化后，当前网络的基本信息，就以[DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm)等级进行输出  The component's logging tool supports logging. As long as the instantiation of the basic network information, the output will be output at [DEBUG](3ce37fee-7678-e278-6a0a-7656126003b6.htm) (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共属性 | [Station](36f1a758-e626-a565-2d13-3a0cdea61613.htm) | 获取或设置当前的地址域信息，是一个12个字符的BCD码，例如：149100007290  Get or set the current address domain information, which is a 12-character BCD code, for example: 149100007290 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[ICjt188 接口](3fabd38b-b4c3-f1b2-0469-ad7704214895.htm)

[HslCommunication.Instrument.CJT.Helper 命名空间](e7617fc1-440e-5d36-1b24-a85e8d4dffb0.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## EnableCodeFE 属性 

[原文連結](http://api.hslcommunication.cn/html/2420a2be-c28d-4ddf-63e7-05e2eb050af8.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.CJT.Helper](../html/e7617fc1-440e-5d36-1b24-a85e8d4dffb0.htm "HslCommunication.Instrument.CJT.Helper")

[ICjt188 接口](../html/3fabd38b-b4c3-f1b2-0469-ad7704214895.htm "ICjt188 接口")

[ICjt188 属性](../html/7dbdd1e5-a59c-0253-ecc4-dea40ae4633a.htm "ICjt188 属性")

[EnableCodeFE 属性](../html/2420a2be-c28d-4ddf-63e7-05e2eb050af8.htm "EnableCodeFE 属性 ")

[InstrumentType 属性](../html/54d205f6-517b-1b53-bb21-afa097006156.htm "InstrumentType 属性 ")

[Station 属性](../html/36f1a758-e626-a565-2d13-3a0cdea61613.htm "Station 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ICjt188EnableCodeFE 属性 |

获取或设置是否在每一次的报文通信时，增加"FE FE"的命令头  
Get or set whether to add the command header of "FE FE" in each message communication

**命名空间：**
 [HslCommunication.Instrument.CJT.Helper](e7617fc1-440e-5d36-1b24-a85e8d4dffb0.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
bool EnableCodeFE { get; set; }
```

```
Property EnableCodeFE As Boolean
	Get
	Set
```

```
property bool EnableCodeFE {
	bool get ();
	void set (bool value);
}
```

```
abstract EnableCodeFE : bool with get, set
```

#### 属性值

类型：Boolean

![](../icons/SectionExpanded.png)参见

#### 引用

[ICjt188 接口](3fabd38b-b4c3-f1b2-0469-ad7704214895.htm)

[HslCommunication.Instrument.CJT.Helper 命名空间](e7617fc1-440e-5d36-1b24-a85e8d4dffb0.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## InstrumentType 属性 

[原文連結](http://api.hslcommunication.cn/html/54d205f6-517b-1b53-bb21-afa097006156.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.CJT.Helper](../html/e7617fc1-440e-5d36-1b24-a85e8d4dffb0.htm "HslCommunication.Instrument.CJT.Helper")

[ICjt188 接口](../html/3fabd38b-b4c3-f1b2-0469-ad7704214895.htm "ICjt188 接口")

[ICjt188 属性](../html/7dbdd1e5-a59c-0253-ecc4-dea40ae4633a.htm "ICjt188 属性")

[EnableCodeFE 属性](../html/2420a2be-c28d-4ddf-63e7-05e2eb050af8.htm "EnableCodeFE 属性 ")

[InstrumentType 属性](../html/54d205f6-517b-1b53-bb21-afa097006156.htm "InstrumentType 属性 ")

[Station 属性](../html/36f1a758-e626-a565-2d13-3a0cdea61613.htm "Station 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ICjt188InstrumentType 属性 |

仪表的类型

**命名空间：**
 [HslCommunication.Instrument.CJT.Helper](e7617fc1-440e-5d36-1b24-a85e8d4dffb0.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
byte InstrumentType { get; set; }
```

```
Property InstrumentType As Byte
	Get
	Set
```

```
property unsigned char InstrumentType {
	unsigned char get ();
	void set (unsigned char value);
}
```

```
abstract InstrumentType : byte with get, set
```

#### 属性值

类型：Byte

![](../icons/SectionExpanded.png)参见

#### 引用

[ICjt188 接口](3fabd38b-b4c3-f1b2-0469-ad7704214895.htm)

[HslCommunication.Instrument.CJT.Helper 命名空间](e7617fc1-440e-5d36-1b24-a85e8d4dffb0.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Station 属性 

[原文連結](http://api.hslcommunication.cn/html/36f1a758-e626-a565-2d13-3a0cdea61613.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.CJT.Helper](../html/e7617fc1-440e-5d36-1b24-a85e8d4dffb0.htm "HslCommunication.Instrument.CJT.Helper")

[ICjt188 接口](../html/3fabd38b-b4c3-f1b2-0469-ad7704214895.htm "ICjt188 接口")

[ICjt188 属性](../html/7dbdd1e5-a59c-0253-ecc4-dea40ae4633a.htm "ICjt188 属性")

[EnableCodeFE 属性](../html/2420a2be-c28d-4ddf-63e7-05e2eb050af8.htm "EnableCodeFE 属性 ")

[InstrumentType 属性](../html/54d205f6-517b-1b53-bb21-afa097006156.htm "InstrumentType 属性 ")

[Station 属性](../html/36f1a758-e626-a565-2d13-3a0cdea61613.htm "Station 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ICjt188Station 属性 |

获取或设置当前的地址域信息，是一个12个字符的BCD码，例如：149100007290  
Get or set the current address domain information, which is a 12-character BCD code, for example: 149100007290

**命名空间：**
 [HslCommunication.Instrument.CJT.Helper](e7617fc1-440e-5d36-1b24-a85e8d4dffb0.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
string Station { get; set; }
```

```
Property Station As String
	Get
	Set
```

```
property String^ Station {
	String^ get ();
	void set (String^ value);
}
```

```
abstract Station : string with get, set
```

#### 属性值

类型：String

![](../icons/SectionExpanded.png)参见

#### 引用

[ICjt188 接口](3fabd38b-b4c3-f1b2-0469-ad7704214895.htm)

[HslCommunication.Instrument.CJT.Helper 命名空间](e7617fc1-440e-5d36-1b24-a85e8d4dffb0.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ICjt188 方法

[原文連結](http://api.hslcommunication.cn/html/44dc2683-e490-896d-05a0-cee4a0920cea.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.CJT.Helper](../html/e7617fc1-440e-5d36-1b24-a85e8d4dffb0.htm "HslCommunication.Instrument.CJT.Helper")

[ICjt188 接口](../html/3fabd38b-b4c3-f1b2-0469-ad7704214895.htm "ICjt188 接口")

[ICjt188 方法](../html/44dc2683-e490-896d-05a0-cee4a0920cea.htm "ICjt188 方法")

[ActiveDeveice 方法](../html/f00c3a89-859f-c8ba-5c0a-b58cd4162eaa.htm "ActiveDeveice 方法 ")

[ReadAddress 方法](../html/a669de5c-80f0-74df-235c-2eccb0ff9c38.htm "ReadAddress 方法 ")

[ReadStringArray 方法](../html/226fd33b-1900-9236-7bb7-ef5b195d6101.htm "ReadStringArray 方法 ")

[WriteAddress 方法](../html/dfdff82a-e574-bf9a-381a-471bfb3c0300.htm "WriteAddress 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ICjt188 方法 |

[ICjt188](3fabd38b-b4c3-f1b2-0469-ad7704214895.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [ActiveDeveice](f00c3a89-859f-c8ba-5c0a-b58cd4162eaa.htm) | 激活设备的命令，只发送数据到设备，不等待设备数据返回  The command to activate the device, only send data to the device, do not wait for the device data to return |
| 公共方法 | [Read(String, UInt16)](b8d4187f-ab9a-79e9-c7ca-8a01ec0df233.htm) | 批量读取字节数组信息，需要指定地址和长度，返回原始的字节数组  Batch read byte array information, need to specify the address and length, return the original byte array (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadT](1be4d687-47e0-f01f-2969-570e5b057c4a.htm) | 读取支持Hsl特性的数据内容，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明，或是官网的详细文档  Read the data content of the Hsl attribute. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [ReadAddress](a669de5c-80f0-74df-235c-2eccb0ff9c38.htm) | 读取设备的通信地址，仅支持点对点通讯的情况，返回地址域数据，例如：14910000729012  Read the communication address of the device, only support point-to-point communication, and return the address field data, for example: 149100007290 |
| 公共方法 | [ReadAsync(String, UInt16)](59839c27-5a76-2f53-730c-c9d2bb41c2c4.htm) | 异步批量读取字节数组信息，需要指定地址和长度，返回原始的字节数组  Asynchronous batch read byte array information, need to specify the address and length, return the original byte array (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadAsyncT](2c18ca3b-cee6-e5ac-8962-48d1f7a99ff4.htm) | 异步读取支持Hsl特性的数据内容，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Asynchronously read the data content of the Hsl attribute. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [ReadBool(String)](67c3a6fd-ad81-62bf-a72d-74b99522a93e.htm) | 读取单个的Boolean数据信息  Read a single Boolean data message (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [ReadBool(String, UInt16)](16e00e01-28d8-bfb7-e3bb-9d8b8778d7b9.htm) | 批量读取Boolean数组信息，需要指定地址和长度，返回Boolean 数组  Batch read Boolean array information, need to specify the address and length, return Boolean array (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [ReadBoolAsync(String)](f99fec07-3c4c-82a3-4427-4cc26281c871.htm) | 异步读取单个的Boolean数据信息  Asynchronously read a single Boolean data message (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [ReadBoolAsync(String, UInt16)](74e40b23-198a-944b-6ed7-3f58ca51da5d.htm) | 异步批量读取Boolean数组信息，需要指定地址和长度，返回Boolean 数组  Asynchronously batch read Boolean array information, need to specify the address and length, return Boolean array (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadCustomerT(String)](04bd0f87-9d33-c4d1-9420-7ca995687eb3.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，返回一个新的类型的实例对象。  To read a custom data type, you need to inherit from the IDataTransfer interface and return an instance object of a new type. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadCustomerT(String, T)](72f8e1a8-fd49-efc7-3485-fb2684f1d948.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，传入一个实例，对这个实例进行赋值，并返回该实例的对象。  To read a custom data type, you need to inherit from the IDataTransfer interface, pass in an instance, assign a value to this instance, and return the object of the instance. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadCustomerAsyncT(String)](482ef33f-8bb1-bea7-2bad-5c23bafcd173.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，返回一个新的类型的实例对象。  To read a custom data type, you need to inherit from the IDataTransfer interface and return an instance object of a new type. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadCustomerAsyncT(String, T)](0405df60-0745-1be6-0101-59edaf4e4d38.htm) | 读取自定义的数据类型，需要继承自IDataTransfer接口，传入一个实例，对这个实例进行赋值，并返回该实例的对象。  To read a custom data type, you need to inherit from the IDataTransfer interface, pass in an instance, assign a value to this instance, and return the object of the instance. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadDouble(String)](18936f56-2b3c-59cc-788b-63490168e4e2.htm) | 读取双浮点的数据  Read double floating point data (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadDouble(String, UInt16)](a8468555-105b-337a-196f-b0a667856edd.htm) | 读取双浮点数据的数组  Read double floating point data array (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadDoubleAsync(String)](d902471c-124a-67a2-be35-e2feebbb1298.htm) | 异步读取双浮点的数据  Asynchronously read double floating point data (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadDoubleAsync(String, UInt16)](3a436322-c7f5-bb1d-e97b-dc6c0d6990d2.htm) | 异步读取双浮点数据的数组  Asynchronously read double floating point data array (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadFloat(String)](95ca7390-d065-ef28-cfdb-e8c161437ecd.htm) | 读取单浮点数据  Read single floating point data (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadFloat(String, UInt16)](a8ae302d-6269-8baa-cacc-f35df9fcd131.htm) | 读取单浮点精度的数组  Read single floating point array (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadFloatAsync(String)](de191b71-b9e7-6f9c-e7a6-9a117e13fb05.htm) | 异步读取单浮点数据  Asynchronously read single floating point data (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadFloatAsync(String, UInt16)](26a13b36-422b-0520-3cb0-66a235e3d7a8.htm) | 异步读取单浮点精度的数组  Asynchronously read single floating point array (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [ReadFromCoreServer(Byte)](413a68ef-2f4b-638f-8bad-973673920de4.htm) | 将当前的数据报文发送到设备去，具体使用什么通信方式取决于设备信息，然后从设备接收数据回来，并返回给调用者。  Send the current data message to the device, the specific communication method used depends on the device information, and then receive the data back from the device and return it to the caller. (继承自 [IReadWriteDevice](d987afcc-f99c-6aed-3847-5b1890504635.htm)。) |
| 公共方法 | [ReadFromCoreServer(IEnumerableByte)](236bfe99-f977-a092-d283-c313ab8aaa7c.htm) | 将多个数据报文按顺序发到设备，并从设备接收返回的数据内容，然后拼接成一个Byte[]信息，需要重写[UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm)方法才能返回正确的结果。  Send multiple data packets to the device in sequence, and receive the returned data content from the device, and then splicing them into a Byte[] message, you need to rewrite [UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm) method to return the correct result. (继承自 [IReadWriteDevice](d987afcc-f99c-6aed-3847-5b1890504635.htm)。) |
| 公共方法 | [ReadFromCoreServerAsync(Byte)](646de5ac-9ca6-1fec-f9de-a318dbc34dc1.htm) | 将当前的数据报文发送到设备去，具体使用什么通信方式取决于设备信息，然后从设备接收数据回来，并返回给调用者。  Send the current data message to the device, the specific communication method used depends on the device information, and then receive the data back from the device and return it to the caller. (继承自 [IReadWriteDevice](d987afcc-f99c-6aed-3847-5b1890504635.htm)。) |
| 公共方法 | [ReadFromCoreServerAsync(IEnumerableByte)](90bb695c-5fc3-eb72-39f6-195819857394.htm) | 将多个数据报文按顺序发到设备，并从设备接收返回的数据内容，然后拼接成一个Byte[]信息，需要重写[UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm)方法才能返回正确的结果。  Send multiple data packets to the device in sequence, and receive the returned data content from the device, and then splicing them into a Byte[] message, you need to rewrite [UnpackResponseContent(Byte, Byte)](32bad9a0-b59d-c171-ce85-74055edfc597.htm) method to return the correct result. (继承自 [IReadWriteDevice](d987afcc-f99c-6aed-3847-5b1890504635.htm)。) |
| 公共方法代码示例 | [ReadInt16(String)](c5308422-4397-2d19-f2de-6d4213394e3d.htm) | 读取16位的有符号的整型数据  Read 16-bit signed integer data (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadInt16(String, UInt16)](3173b3fb-e0ac-1b58-ffff-3741aa992458.htm) | 读取16位的有符号整型数组  Read 16-bit signed integer array (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadInt16Async(String)](bef4e722-646e-3196-0fc7-bf12cfa7eb57.htm) | 异步读取16位的有符号的整型数据  Asynchronously read 16-bit signed integer data (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadInt16Async(String, UInt16)](1d194744-a756-5f40-90a8-ec59633a3970.htm) | 异步读取16位的有符号整型数组  Asynchronously read 16-bit signed integer array (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadInt32(String)](51ba7f65-d19b-20dd-e81e-f3913226b5cc.htm) | 读取32位的有符号整型  Read 32-bit signed integer (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadInt32(String, UInt16)](39024553-e8a3-ceb6-946c-0812c211c218.htm) | 读取32位有符号整型数组  Read 32-bit signed integer array (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadInt32Async(String)](fcf027ad-23d6-a6a4-c715-b0c1dcf18b31.htm) | 异步读取32位的有符号整型  Asynchronously read 32-bit signed integer (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadInt32Async(String, UInt16)](f873bfc0-32e5-3303-ec0a-b0c5515754c6.htm) | 异步读取32位有符号整型数组  Asynchronously read 32-bit signed integer array (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadInt64(String)](a3bea048-fd83-2158-4580-231419212efe.htm) | 读取64位的有符号整型  Read 64-bit signed integer (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadInt64(String, UInt16)](32190a15-c77a-7831-9d59-f2419a7f8bb2.htm) | 读取64位的有符号整型数组  Read 64-bit signed integer array (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadInt64Async(String)](d2220dae-6131-a6aa-8533-a5cce335e2c4.htm) | 异步读取64位的有符号整型  Asynchronously read 64-bit signed integer (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadInt64Async(String, UInt16)](b511a119-ee98-1abc-b578-edcac930fcd4.htm) | 异步读取64位的有符号整型数组  Asynchronously read 64-bit signed integer array (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadString(String, UInt16)](a8ca8b30-6612-307e-c046-d8ff01b30df1.htm) | 读取字符串数据，默认为最常见的ASCII编码  Read string data, default is the most common ASCII encoding (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadString(String, UInt16, Encoding)](fc289def-e53b-b47d-22fe-7e02b5e354e0.htm) | 使用指定的编码，读取字符串数据  Reads string data using the specified encoding (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [ReadStringArray](226fd33b-1900-9236-7bb7-ef5b195d6101.htm) | 读取指定地址的所有的字符串数据信息，一般来说，一个地址只有一个数据 |
| 公共方法代码示例 | [ReadStringAsync(String, UInt16)](960c5f9b-c87b-d7e9-a4e0-048e4f2f5981.htm) | 异步读取字符串数据，默认为最常见的ASCII编码  Asynchronously read string data, default is the most common ASCII encoding (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadStringAsync(String, UInt16, Encoding)](db06baf4-e674-22de-e89e-cb26646413f8.htm) | 异步使用指定的编码，读取字符串数据  Asynchronously reads string data using the specified encoding (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadStructT](50df0eb8-8f9c-ad9c-4dc9-bf75b4438d46.htm) | 读取结构体类型的数据，根据结构体自身的定义，读取原始字节数组，然后解析出实际的结构体数据，结构体需要实现[HslStructAttribute](87e427a8-0c7b-43a8-c421-c27264eda90f.htm)特性 (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadStructAsyncT](af626b7e-2adb-ecd1-f2a0-82aaa6e8b68d.htm) | 读取结构体类型的数据，根据结构体自身的定义，读取原始字节数组，然后解析出实际的结构体数据，结构体需要实现[HslStructAttribute](87e427a8-0c7b-43a8-c421-c27264eda90f.htm)特性 (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadUInt16(String)](66238795-7948-2cf9-c639-d1bb3da126dd.htm) | 读取16位的无符号整型  Read 16-bit unsigned integer (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadUInt16(String, UInt16)](28888661-1fcc-f941-5746-9124eb539fe8.htm) | 读取16位的无符号整型数组  Read 16-bit unsigned integer array (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadUInt16Async(String)](1744a25b-c440-7d1e-dee5-606b0ff485bc.htm) | 异步读取16位的无符号整型  Asynchronously read 16-bit unsigned integer (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadUInt16Async(String, UInt16)](0748f071-d6a7-581b-ed7a-d73317b16762.htm) | 异步读取16位的无符号整型数组  Asynchronously read 16-bit unsigned integer array (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadUInt32(String)](b252ad01-4035-bc3e-e297-7fd97c3ac7af.htm) | 读取32位的无符号整型  Read 32-bit unsigned integer (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadUInt32(String, UInt16)](64d8381b-05b6-a990-c03e-2c8a1b7cfb7f.htm) | 读取32位的无符号整型数组  Read 32-bit unsigned integer array (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadUInt32Async(String)](a5b848e0-82bd-2bb5-2dd9-6634063e7b48.htm) | 异步读取32位的无符号整型  Asynchronously read 32-bit unsigned integer (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadUInt32Async(String, UInt16)](59ba0cbe-08f8-5a30-6029-0a9f718c7d76.htm) | 异步读取32位的无符号整型数组  Asynchronously read 32-bit unsigned integer array (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadUInt64(String)](529c82ba-1d5e-f33b-2eb9-c862799f952f.htm) | 读取64位的无符号整型  Read 64-bit unsigned integer (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadUInt64(String, UInt16)](2af999cc-c8d4-e25a-68f1-d621e5d84af7.htm) | 读取64位的无符号整型的数组  Read 64-bit unsigned integer array (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadUInt64Async(String)](0c521a11-248c-ba7a-fc27-0e5a645611b5.htm) | 异步读取64位的无符号整型  Asynchronously read 64-bit unsigned integer (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [ReadUInt64Async(String, UInt16)](b510070f-beb5-959b-7c6b-20f212a68dd1.htm) | 异步读取64位的无符号整型的数组  Asynchronously read 64-bit unsigned integer array (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [Wait(String, Boolean, Int32, Int32)](c8c01cf3-bf99-761c-23d4-82dd54081952.htm) | 等待指定地址的Boolean值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the Boolean value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [Wait(String, Int16, Int32, Int32)](f8a820b0-cf90-668c-6fba-524bd0659042.htm) | 等待指定地址的Int16值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the Int16 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [Wait(String, Int32, Int32, Int32)](96ad6909-b660-2894-9c17-a7f4f4d4a2bd.htm) | 等待指定地址的Int32值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the Int32 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [Wait(String, Int64, Int32, Int32)](28ba4678-2aeb-7502-95e8-d6b855bbe971.htm) | 等待指定地址的Int64值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the Int64 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [Wait(String, UInt16, Int32, Int32)](44af3173-f7f4-04f5-753b-7855581a171a.htm) | 等待指定地址的UInt16值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the UInt16 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [Wait(String, UInt32, Int32, Int32)](0e877e55-9b76-f662-6841-9102a9e9b9b8.htm) | 等待指定地址的UInt32值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the UInt32 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [Wait(String, UInt64, Int32, Int32)](6801831d-c349-800d-e3d3-34d2fde09e0f.htm) | 等待指定地址的UInt64值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the UInt64 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [WaitAsync(String, Boolean, Int32, Int32)](b51daff5-9c09-c614-9ea3-f14e907a36b8.htm) | 等待指定地址的Boolean值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the Boolean value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [WaitAsync(String, Int16, Int32, Int32)](cfd8af62-371b-84e9-5695-87539fbb204b.htm) | 等待指定地址的Int16值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the Int16 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [WaitAsync(String, Int32, Int32, Int32)](639f2d1b-9a41-0cf2-0335-7f9641fbb0c8.htm) | 等待指定地址的Int32值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the Int32 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [WaitAsync(String, Int64, Int32, Int32)](4355bcb4-e923-9226-ac82-55dd1325f34f.htm) | 等待指定地址的Int64值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the Int64 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [WaitAsync(String, UInt16, Int32, Int32)](291b6f29-9cbf-6fc3-4828-226f412e9cc6.htm) | 等待指定地址的UInt16值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the UInt16 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [WaitAsync(String, UInt32, Int32, Int32)](e1009c58-9f02-660e-225a-b6861eb8781f.htm) | 等待指定地址的UInt32值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the UInt32 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [WaitAsync(String, UInt64, Int32, Int32)](4236aa40-6c9d-2a0c-8470-c55dab497074.htm) | 等待指定地址的UInt64值为指定的值，可以指定刷新数据的频率，等待的超时时间，如果超时时间为-1的话，则是无期限等待。  Waiting for the UInt64 value of the specified address to be the specified value, you can specify the frequency of refreshing the data, and the timeout time to wait. If the timeout time is -1, it is an indefinite wait. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [Write(String, Boolean)](4d742738-6864-465c-5382-bd49032eae37.htm) | 写入单个的Boolean数据，返回是否成功  Write a single Boolean data, and return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [Write(String, Boolean)](2a99122c-1496-d766-f0aa-44a18918fcd3.htm) | 批量写入Boolean数组数据，返回是否成功  Batch write Boolean array data, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [Write(String, Byte)](b8d8d94e-300d-c7b7-2e80-6859c62479fb.htm) | 写入原始的byte数组数据到指定的地址，返回是否写入成功  Write the original byte array data to the specified address, and return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [Write(String, Double)](34e1c967-7214-bb8e-3daf-fdb93fc1db59.htm) | 写入double数据，返回是否成功  Write double data, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [Write(String, Double)](f2a8cc81-ee85-4d3a-1eaf-af33b15de5e6.htm) | 写入double数组，返回是否成功  Write double array, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [Write(String, Int16)](53b27553-763c-2dd9-b8a1-df2ebf68ab2c.htm) | 写入short数据，返回是否成功  Write short data, returns whether success (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [Write(String, Int16)](0b8fff2d-0da9-d73a-bc11-866a501e6112.htm) | 写入short数组，返回是否成功  Write short array, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [Write(String, Int32)](281ce6f6-5c53-da62-9a09-9d439390a972.htm) | 写入int数据，返回是否成功  Write int data, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [Write(String, Int32)](496e02c2-62f2-4f72-aff1-1d691cc87726.htm) | 写入int[]数组，返回是否成功  Write int array, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [Write(String, Int64)](7ec919a0-1d6b-b7d7-1f2f-eb1644beb912.htm) | 写入long数据，返回是否成功  Write long data, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [Write(String, Int64)](026100a5-d52e-b640-f239-a140c903201d.htm) | 写入long数组，返回是否成功  Write long array, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [Write(String, Single)](7c72c919-367f-9bc6-567c-550813459297.htm) | 写入float数据，返回是否成功  Write float data, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [Write(String, Single)](195c858e-1e2d-863a-8a2c-dc49cb2063d3.htm) | 写入float数组，返回是否成功  Write float array, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [Write(String, String)](fc988c8f-b2d6-18f9-3274-9b57d4a04e1c.htm) | 写入字符串信息，编码为ASCII  Write string information, encoded as ASCII (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [Write(String, UInt16)](7084bc37-6bda-e506-4118-74aff979efff.htm) | 写入ushort数据，返回是否成功  Write ushort data, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [Write(String, UInt16)](9c21745f-7c5f-cae3-7f2a-2bcf260e54a9.htm) | 写入ushort数组，返回是否成功  Write ushort array, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [Write(String, UInt32)](ed79fea9-bc6e-3ef3-5294-134b872141cd.htm) | 写入uint数据，返回是否成功  Write uint data, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [Write(String, UInt32)](3ae38bac-66cc-5cdf-be68-483362b17ac3.htm) | 写入uint[]数组，返回是否成功  Write uint array, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [Write(String, UInt64)](3927095f-0705-e823-bb24-74c9363fab36.htm) | 写入ulong数据，返回是否成功  Write ulong data, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [Write(String, UInt64)](2eb4a3ce-c350-9949-405b-ccb6b5d9786e.htm) | 写入ulong数组，返回是否成功  Write ulong array, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [Write(String, String, Int32)](1ed07021-8f14-743a-071c-c348025634f5.htm) | 写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为ASCII  Write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, it is filled with 0 and the encoding is ASCII. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [Write(String, String, Encoding)](c66a9dfb-840a-3ac4-96ef-48e23fb85ec1.htm) | 写入字符串信息，需要指定的编码信息  Write string information, need to specify the encoding information (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [Write(String, String, Int32, Encoding)](61396f95-43cd-34fb-9055-37676c4e436d.htm) | 写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为指定的编码信息  Write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, then the operation is complemented with 0 , you should specified the encoding information (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteT(T)](ec80557a-6d74-6030-7e09-57c0c6a63cc4.htm) | 写入支持Hsl特性的数据，返回是否写入成功，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Write data that supports the Hsl attribute, and return whether the write was successful. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [WriteAddress](dfdff82a-e574-bf9a-381a-471bfb3c0300.htm) | 写入设备的地址域信息，仅支持点对点通讯的情况，需要指定地址域信息，例如：14910000729011  Write the address domain information of the device, only support point-to-point communication, you need to specify the address domain information, for example: 14910000729011 |
| 公共方法 | [WriteAsync(String, Boolean)](559924da-b298-9c81-794f-83d53aafc602.htm) | 异步批量写入Boolean数组数据，返回是否成功  Asynchronously batch write Boolean array data, return success (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [WriteAsync(String, Boolean)](cf8f8e56-df3d-607b-1c52-c1a76378db39.htm) | 异步批量写入Boolean数组数据，返回是否成功  Asynchronously batch write Boolean array data, return success (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Byte)](1a2e83c8-7782-847a-5d89-08d5556ac4ad.htm) | 异步写入原始的byte数组数据到指定的地址，返回是否写入成功  Asynchronously writes the original byte array data to the specified address, and returns whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Double)](ca502ab0-a05c-9dc7-c8b0-9afd3a24cb63.htm) | 异步写入double数据，返回是否成功  Asynchronously write double data, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Double)](33c58466-53fe-e221-1d73-6650c434e2b1.htm) | 异步写入double数组，返回是否成功  Asynchronously write double array, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int16)](3e923476-059a-0de9-4fc5-6df5dc7a93a5.htm) | 异步写入short数据，返回是否成功  Asynchronously write short data, returns whether success (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int16)](8955c5b3-7709-57ce-51fa-a58369f1147f.htm) | 异步写入short数组，返回是否成功  Asynchronously write short array, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int32)](448b905c-dfee-7f97-4dd9-8caf3539eef4.htm) | 异步写入int数据，返回是否成功  Asynchronously write int data, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int32)](cd5f66ad-6ffc-5394-4554-43abb07d977a.htm) | 异步写入int[]数组，返回是否成功  Asynchronously write int array, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int64)](972d6a25-bb1a-3c92-c12a-4b6942fc026f.htm) | 异步写入long数据，返回是否成功  Asynchronously write long data, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Int64)](89e51562-f954-65b5-0c20-d1602d8e78a8.htm) | 异步写入long数组，返回是否成功  Asynchronously write long array, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Single)](213794a0-4dc8-36b4-3b3f-d1ee74e35155.htm) | 异步写入float数据，返回是否成功  Asynchronously write float data, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, Single)](dca9a471-a3cb-3242-b3fe-edeb04d4f786.htm) | 异步写入float数组，返回是否成功  Asynchronously write float array, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, String)](265dfc0f-8630-cc0a-4a25-cd31d50f6056.htm) | 异步写入字符串信息，编码为ASCII  Asynchronously write string information, encoded as ASCII (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt16)](ee35584e-d8c8-720e-7a05-24a2d3011021.htm) | 异步写入ushort数据，返回是否成功  Asynchronously write ushort data, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt16)](9df3d402-f0c7-5888-65d0-f527a6ada1b2.htm) | 异步写入ushort数组，返回是否成功  Asynchronously write ushort array, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt32)](76ac4f61-3a88-384c-fae1-857066a74ac0.htm) | 异步写入uint数据，返回是否成功  Asynchronously write uint data, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt32)](aae8db84-1d63-31a3-e84c-f03a5f784ff0.htm) | 异步写入uint[]数组，返回是否成功  Asynchronously write uint array, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt64)](81bbb9e0-b222-067d-7dd9-a879277859c4.htm) | 异步写入ulong数据，返回是否成功  Asynchronously write ulong data, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, UInt64)](f56c270e-dc8d-9adc-ea60-36d47e2c6e50.htm) | 异步写入ulong数组，返回是否成功  Asynchronously write ulong array, return whether the write was successful (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [WriteAsync(String, String, Int32)](c9efea95-68b3-ab09-5176-d388576844e8.htm) | 异步写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为ASCII  Asynchronously write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, it is filled with 0 and the encoding is ASCII. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteAsync(String, String, Encoding)](cda78ec9-ab0c-b6f1-db38-0e51e6aa06d0.htm) | 异步写入字符串信息，需要指定的编码信息  Asynchronously write string information, need to specify the encoding information (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法 | [WriteAsync(String, String, Int32, Encoding)](c077e3ad-d170-aea3-0058-2888f2b27126.htm) | 异步写入指定长度的字符串信息，如果超出，就截断字符串，如果长度不足，那就补0操作，编码为指定的编码信息  Asynchronously write string information of the specified length. If it exceeds the value, the string is truncated. If the length is not enough, then the operation is complemented with 0 , you should specified the encoding information (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteAsyncT(T)](da0f692b-4d73-0520-6e17-7890b6ef51cb.htm) | 异步写入支持Hsl特性的数据，返回是否写入成功，该特性为[HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm)，详细参考api文档说明  Asynchronously write data that supports the Hsl attribute, and return whether the write was successful. The attribute is [HslDeviceAddressAttribute](e39ddd01-95cf-ea40-1f29-20801dddfbe0.htm), please refer to the api documentation for details. (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteCustomerT](466c49b4-95ac-148b-c18d-a419c5f6dab3.htm) | 写入自定义类型的数据，该类型必须继承自IDataTransfer接口  Write data of a custom type, which must inherit from the IDataTransfer interface (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |
| 公共方法代码示例 | [WriteCustomerAsyncT](82aa8f35-3c06-d4f6-e8c5-affb6ab6ce41.htm) | 写入自定义类型的数据，该类型必须继承自IDataTransfer接口  Write data of a custom type, which must inherit from the IDataTransfer interface (继承自 [IReadWriteNet](d90cbd6f-aef5-0e38-b28f-56bb2248b12f.htm)。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[ICjt188 接口](3fabd38b-b4c3-f1b2-0469-ad7704214895.htm)

[HslCommunication.Instrument.CJT.Helper 命名空间](e7617fc1-440e-5d36-1b24-a85e8d4dffb0.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ActiveDeveice 方法 

[原文連結](http://api.hslcommunication.cn/html/f00c3a89-859f-c8ba-5c0a-b58cd4162eaa.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.CJT.Helper](../html/e7617fc1-440e-5d36-1b24-a85e8d4dffb0.htm "HslCommunication.Instrument.CJT.Helper")

[ICjt188 接口](../html/3fabd38b-b4c3-f1b2-0469-ad7704214895.htm "ICjt188 接口")

[ICjt188 方法](../html/44dc2683-e490-896d-05a0-cee4a0920cea.htm "ICjt188 方法")

[ActiveDeveice 方法](../html/f00c3a89-859f-c8ba-5c0a-b58cd4162eaa.htm "ActiveDeveice 方法 ")

[ReadAddress 方法](../html/a669de5c-80f0-74df-235c-2eccb0ff9c38.htm "ReadAddress 方法 ")

[ReadStringArray 方法](../html/226fd33b-1900-9236-7bb7-ef5b195d6101.htm "ReadStringArray 方法 ")

[WriteAddress 方法](../html/dfdff82a-e574-bf9a-381a-471bfb3c0300.htm "WriteAddress 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ICjt188ActiveDeveice 方法 |

激活设备的命令，只发送数据到设备，不等待设备数据返回  
The command to activate the device, only send data to the device, do not wait for the device data to return

**命名空间：**
 [HslCommunication.Instrument.CJT.Helper](e7617fc1-440e-5d36-1b24-a85e8d4dffb0.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
OperateResult ActiveDeveice()
```

```
Function ActiveDeveice As OperateResult
```

```
OperateResult^ ActiveDeveice()
```

```
abstract ActiveDeveice : unit -> OperateResult 
```

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否发送成功

![](../icons/SectionExpanded.png)参见

#### 引用

[ICjt188 接口](3fabd38b-b4c3-f1b2-0469-ad7704214895.htm)

[HslCommunication.Instrument.CJT.Helper 命名空间](e7617fc1-440e-5d36-1b24-a85e8d4dffb0.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadAddress 方法 

[原文連結](http://api.hslcommunication.cn/html/a669de5c-80f0-74df-235c-2eccb0ff9c38.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Instrument.CJT.Helper](../html/e7617fc1-440e-5d36-1b24-a85e8d4dffb0.htm "HslCommunication.Instrument.CJT.Helper")

[ICjt188 接口](../html/3fabd38b-b4c3-f1b2-0469-ad7704214895.htm "ICjt188 接口")

[ICjt188 方法](../html/44dc2683-e490-896d-05a0-cee4a0920cea.htm "ICjt188 方法")

[ActiveDeveice 方法](../html/f00c3a89-859f-c8ba-5c0a-b58cd4162eaa.htm "ActiveDeveice 方法 ")

[ReadAddress 方法](../html/a669de5c-80f0-74df-235c-2eccb0ff9c38.htm "ReadAddress 方法 ")

[ReadStringArray 方法](../html/226fd33b-1900-9236-7bb7-ef5b195d6101.htm "ReadStringArray 方法 ")

[WriteAddress 方法](../html/dfdff82a-e574-bf9a-381a-471bfb3c0300.htm "WriteAddress 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| ICjt188ReadAddress 方法 |

读取设备的通信地址，仅支持点对点通讯的情况，返回地址域数据，例如：14910000729012  
Read the communication address of the device, only support point-to-point communication, and return the address field data, for example: 149100007290

**命名空间：**
 [HslCommunication.Instrument.CJT.Helper](e7617fc1-440e-5d36-1b24-a85e8d4dffb0.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
OperateResult<string> ReadAddress()
```

```
Function ReadAddress As OperateResult(Of String)
```

```
OperateResult<String^>^ ReadAddress()
```

```
abstract ReadAddress : unit -> OperateResult<string> 
```

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)String  
设备的通信地址

![](../icons/SectionExpanded.png)参见

#### 引用

[ICjt188 接口](3fabd38b-b4c3-f1b2-0469-ad7704214895.htm)

[HslCommunication.Instrument.CJT.Helper 命名空间](e7617fc1-440e-5d36-1b24-a85e8d4dffb0.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)