# HslCommunication - HslCommunication.Profinet.Delta.Helper

> 分類頁數: 22



---
## HslCommunication.Profinet.Delta.Helper

[原文連結](http://api.hslcommunication.cn/html/cdc3609c-da0b-e755-d4f9-5d1bc3fa17e1.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Delta.Helper](../html/cdc3609c-da0b-e755-d4f9-5d1bc3fa17e1.htm "HslCommunication.Profinet.Delta.Helper")

[DeltaASHelper 类](../html/08716b14-6e57-3a77-f6b9-a1e0fd2372b8.htm "DeltaASHelper 类")

[DeltaDvpHelper 类](../html/29156f81-12f6-a557-f020-e351396152fd.htm "DeltaDvpHelper 类")

[DeltaHelper 类](../html/ffea0bf8-6cae-1863-da2c-e772ad28c54e.htm "DeltaHelper 类")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslCommunication.Profinet.Delta.Helper 命名空间 |

[缺少 "N:HslCommunication.Profinet.Delta.Helper" 的 <summary> 文档]

![](../icons/SectionExpanded.png)类

|  | 类 | 说明 |
| --- | --- | --- |
| 公共类 | [DeltaASHelper](08716b14-6e57-3a77-f6b9-a1e0fd2372b8.htm) | 台达AS300的辅助帮助类信息 |
| 公共类 | [DeltaDvpHelper](29156f81-12f6-a557-f020-e351396152fd.htm) | 台达PLC的相关的帮助类，公共的地址解析的方法。  Delta PLC related help classes, public address resolution methods. |
| 公共类 | [DeltaHelper](ffea0bf8-6cae-1863-da2c-e772ad28c54e.htm) | 台达的想关的辅助类 |

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## DeltaASHelper 类

[原文連結](http://api.hslcommunication.cn/html/08716b14-6e57-3a77-f6b9-a1e0fd2372b8.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Delta.Helper](../html/cdc3609c-da0b-e755-d4f9-5d1bc3fa17e1.htm "HslCommunication.Profinet.Delta.Helper")

[DeltaASHelper 类](../html/08716b14-6e57-3a77-f6b9-a1e0fd2372b8.htm "DeltaASHelper 类")

[DeltaASHelper 构造函数](../html/3559f448-2d8d-5e67-d337-6187582725e1.htm "DeltaASHelper 构造函数 ")

[DeltaASHelper 方法](../html/f9305465-80c6-fc2d-a799-39c0d9bf11b8.htm "DeltaASHelper 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DeltaASHelper 类 |

台达AS300的辅助帮助类信息

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  HslCommunication.Profinet.Delta.HelperDeltaASHelper

**命名空间：**
 [HslCommunication.Profinet.Delta.Helper](cdc3609c-da0b-e755-d4f9-5d1bc3fa17e1.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class DeltaASHelper
```

```
Public Class DeltaASHelper
```

```
public ref class DeltaASHelper
```

```
type DeltaASHelper =  class end
```

DeltaASHelper 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [DeltaASHelper](3559f448-2d8d-5e67-d337-6187582725e1.htm) | 初始化 DeltaASHelper 类的一个新实例 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法静态成员 | [ParseDeltaASAddress](39ebaa3d-c414-852c-9896-afb25fc25e5e.htm) | 根据台达AS300的PLC的地址，解析出转换后的modbus协议信息，适用AS300系列，当前的地址仍然支持站号指定，例如s=2;D100  According to the PLC address of Delta AS300, the converted modbus protocol information is parsed, and it is applicable to AS300 series. The current address still supports station number designation, for example, s=2;D100 |
| 公共方法 | ToString | (继承自 Object。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.Profinet.Delta.Helper 命名空间](cdc3609c-da0b-e755-d4f9-5d1bc3fa17e1.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## DeltaASHelper 构造函数 

[原文連結](http://api.hslcommunication.cn/html/3559f448-2d8d-5e67-d337-6187582725e1.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Delta.Helper](../html/cdc3609c-da0b-e755-d4f9-5d1bc3fa17e1.htm "HslCommunication.Profinet.Delta.Helper")

[DeltaASHelper 类](../html/08716b14-6e57-3a77-f6b9-a1e0fd2372b8.htm "DeltaASHelper 类")

[DeltaASHelper 构造函数](../html/3559f448-2d8d-5e67-d337-6187582725e1.htm "DeltaASHelper 构造函数 ")

[DeltaASHelper 方法](../html/f9305465-80c6-fc2d-a799-39c0d9bf11b8.htm "DeltaASHelper 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DeltaASHelper 构造函数 |

初始化 [DeltaASHelper](08716b14-6e57-3a77-f6b9-a1e0fd2372b8.htm) 类的一个新实例

**命名空间：**
 [HslCommunication.Profinet.Delta.Helper](cdc3609c-da0b-e755-d4f9-5d1bc3fa17e1.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public DeltaASHelper()
```

```
Public Sub New
```

```
public:
DeltaASHelper()
```

```
new : unit -> DeltaASHelper
```

![](../icons/SectionExpanded.png)参见

#### 引用

[DeltaASHelper 类](08716b14-6e57-3a77-f6b9-a1e0fd2372b8.htm)

[HslCommunication.Profinet.Delta.Helper 命名空间](cdc3609c-da0b-e755-d4f9-5d1bc3fa17e1.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## DeltaASHelper 方法

[原文連結](http://api.hslcommunication.cn/html/f9305465-80c6-fc2d-a799-39c0d9bf11b8.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Delta.Helper](../html/cdc3609c-da0b-e755-d4f9-5d1bc3fa17e1.htm "HslCommunication.Profinet.Delta.Helper")

[DeltaASHelper 类](../html/08716b14-6e57-3a77-f6b9-a1e0fd2372b8.htm "DeltaASHelper 类")

[DeltaASHelper 方法](../html/f9305465-80c6-fc2d-a799-39c0d9bf11b8.htm "DeltaASHelper 方法")

[ParseDeltaASAddress 方法](../html/39ebaa3d-c414-852c-9896-afb25fc25e5e.htm "ParseDeltaASAddress 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DeltaASHelper 方法 |

[DeltaASHelper](08716b14-6e57-3a77-f6b9-a1e0fd2372b8.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法静态成员 | [ParseDeltaASAddress](39ebaa3d-c414-852c-9896-afb25fc25e5e.htm) | 根据台达AS300的PLC的地址，解析出转换后的modbus协议信息，适用AS300系列，当前的地址仍然支持站号指定，例如s=2;D100  According to the PLC address of Delta AS300, the converted modbus protocol information is parsed, and it is applicable to AS300 series. The current address still supports station number designation, for example, s=2;D100 |
| 公共方法 | ToString | (继承自 Object。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[DeltaASHelper 类](08716b14-6e57-3a77-f6b9-a1e0fd2372b8.htm)

[HslCommunication.Profinet.Delta.Helper 命名空间](cdc3609c-da0b-e755-d4f9-5d1bc3fa17e1.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ParseDeltaASAddress 方法 

[原文連結](http://api.hslcommunication.cn/html/39ebaa3d-c414-852c-9896-afb25fc25e5e.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Delta.Helper](../html/cdc3609c-da0b-e755-d4f9-5d1bc3fa17e1.htm "HslCommunication.Profinet.Delta.Helper")

[DeltaASHelper 类](../html/08716b14-6e57-3a77-f6b9-a1e0fd2372b8.htm "DeltaASHelper 类")

[DeltaASHelper 方法](../html/f9305465-80c6-fc2d-a799-39c0d9bf11b8.htm "DeltaASHelper 方法")

[ParseDeltaASAddress 方法](../html/39ebaa3d-c414-852c-9896-afb25fc25e5e.htm "ParseDeltaASAddress 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DeltaASHelperParseDeltaASAddress 方法 |

根据台达AS300的PLC的地址，解析出转换后的modbus协议信息，适用AS300系列，当前的地址仍然支持站号指定，例如s=2;D100  
According to the PLC address of Delta AS300, the converted modbus protocol information is parsed,
and it is applicable to AS300 series. The current address still supports station number designation, for example, s=2;D100

**命名空间：**
 [HslCommunication.Profinet.Delta.Helper](cdc3609c-da0b-e755-d4f9-5d1bc3fa17e1.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult<string> ParseDeltaASAddress(
	string address,
	byte modbusCode
)
```

```
Public Shared Function ParseDeltaASAddress ( 
	address As String,
	modbusCode As Byte
) As OperateResult(Of String)
```

```
public:
static OperateResult<String^>^ ParseDeltaASAddress(
	String^ address, 
	unsigned char modbusCode
)
```

```
static member ParseDeltaASAddress : 
        address : string * 
        modbusCode : byte -> OperateResult<string> 
```

#### 参数

address
:   类型：SystemString  
    台达plc的地址信息

modbusCode
:   类型：SystemByte  
    原始的对应的modbus信息

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)String  
还原后的modbus地址

![](../icons/SectionExpanded.png)参见

#### 引用

[DeltaASHelper 类](08716b14-6e57-3a77-f6b9-a1e0fd2372b8.htm)

[HslCommunication.Profinet.Delta.Helper 命名空间](cdc3609c-da0b-e755-d4f9-5d1bc3fa17e1.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## DeltaDvpHelper 类

[原文連結](http://api.hslcommunication.cn/html/29156f81-12f6-a557-f020-e351396152fd.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Delta.Helper](../html/cdc3609c-da0b-e755-d4f9-5d1bc3fa17e1.htm "HslCommunication.Profinet.Delta.Helper")

[DeltaDvpHelper 类](../html/29156f81-12f6-a557-f020-e351396152fd.htm "DeltaDvpHelper 类")

[DeltaDvpHelper 构造函数](../html/d7ed5d89-301c-d835-4c12-f84c2187c649.htm "DeltaDvpHelper 构造函数 ")

[DeltaDvpHelper 方法](../html/811ac077-29f1-f44b-7308-481786c712f6.htm "DeltaDvpHelper 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DeltaDvpHelper 类 |

台达PLC的相关的帮助类，公共的地址解析的方法。  
Delta PLC related help classes, public address resolution methods.

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  HslCommunication.Profinet.Delta.HelperDeltaDvpHelper

**命名空间：**
 [HslCommunication.Profinet.Delta.Helper](cdc3609c-da0b-e755-d4f9-5d1bc3fa17e1.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class DeltaDvpHelper
```

```
Public Class DeltaDvpHelper
```

```
public ref class DeltaDvpHelper
```

```
type DeltaDvpHelper =  class end
```

DeltaDvpHelper 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [DeltaDvpHelper](d7ed5d89-301c-d835-4c12-f84c2187c649.htm) | 初始化 DeltaDvpHelper 类的一个新实例 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法静态成员 | [ParseDeltaDvpAddress](f89cf60a-cc25-7867-0f05-b21c819470e1.htm) | 根据台达PLC的地址，解析出转换后的modbus协议信息，适用DVP系列，当前的地址仍然支持站号指定，例如s=2;D100  According to the address of Delta PLC, the converted modbus protocol information is parsed out, applicable to DVP series, the current address still supports station number designation, such as s=2;D100 |
| 公共方法静态成员 | [Read](b2b0c6a2-eff4-3a00-6d9d-4c91f8427c91.htm) | 读取台达PLC的原始字节变量，重写了读D地址时，跨区域读4096地址时，将会分割多次读取 |
| 公共方法静态成员 | [ReadAsync](8dba9dbc-4874-0d75-9b17-7443d45092f9.htm) | 读取台达PLC的原始字节变量，重写了读D地址时，跨区域读4096地址时，将会分割多次读取 |
| 公共方法静态成员 | [ReadBool](253ddbd6-6903-24ac-4f04-a7323687dc96.htm) | 读取台达PLC的bool变量，重写了读M地址时，跨区域读1536地址时，将会分割多次读取 |
| 公共方法静态成员 | [ReadBoolAsync](e1185968-3493-27d4-eefc-a8f8a70c2b5f.htm) | 读取台达PLC的bool变量，重写了读M地址时，跨区域读1536地址时，将会分割多次读取 |
| 公共方法 | ToString | (继承自 Object。) |
| 公共方法静态成员 | [Write(FuncString, Boolean, OperateResult, String, Boolean)](db83d633-4946-68f7-19a0-105612a919e9.htm) | 写入台达PLC的bool数据，当发现是M类型的数据，并且地址出现跨1536时，进行切割写入操作 |
| 公共方法静态成员 | [Write(FuncString, Byte, OperateResult, String, Byte)](ff119d13-d37c-183d-015f-54ec55816e55.htm) | 写入台达PLC的原始字节数据，当发现是D类型的数据，并且地址出现跨4096时，进行切割写入操作 |
| 公共方法静态成员 | [WriteAsync(FuncString, Boolean, TaskOperateResult, String, Boolean)](a303aec2-4f6c-6d34-d5bf-6a04328dc9bb.htm) | 写入台达PLC的bool数据，当发现是M类型的数据，并且地址出现跨1536时，进行切割写入操作 |
| 公共方法静态成员 | [WriteAsync(FuncString, Byte, TaskOperateResult, String, Byte)](88c953ac-9a5c-4fd9-2ca0-f20e5fbb00a5.htm) | 写入台达PLC的原始字节数据，当发现是D类型的数据，并且地址出现跨4096时，进行切割写入操作 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.Profinet.Delta.Helper 命名空间](cdc3609c-da0b-e755-d4f9-5d1bc3fa17e1.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## DeltaDvpHelper 构造函数 

[原文連結](http://api.hslcommunication.cn/html/d7ed5d89-301c-d835-4c12-f84c2187c649.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Delta.Helper](../html/cdc3609c-da0b-e755-d4f9-5d1bc3fa17e1.htm "HslCommunication.Profinet.Delta.Helper")

[DeltaDvpHelper 类](../html/29156f81-12f6-a557-f020-e351396152fd.htm "DeltaDvpHelper 类")

[DeltaDvpHelper 构造函数](../html/d7ed5d89-301c-d835-4c12-f84c2187c649.htm "DeltaDvpHelper 构造函数 ")

[DeltaDvpHelper 方法](../html/811ac077-29f1-f44b-7308-481786c712f6.htm "DeltaDvpHelper 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DeltaDvpHelper 构造函数 |

初始化 [DeltaDvpHelper](29156f81-12f6-a557-f020-e351396152fd.htm) 类的一个新实例

**命名空间：**
 [HslCommunication.Profinet.Delta.Helper](cdc3609c-da0b-e755-d4f9-5d1bc3fa17e1.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public DeltaDvpHelper()
```

```
Public Sub New
```

```
public:
DeltaDvpHelper()
```

```
new : unit -> DeltaDvpHelper
```

![](../icons/SectionExpanded.png)参见

#### 引用

[DeltaDvpHelper 类](29156f81-12f6-a557-f020-e351396152fd.htm)

[HslCommunication.Profinet.Delta.Helper 命名空间](cdc3609c-da0b-e755-d4f9-5d1bc3fa17e1.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## DeltaDvpHelper 方法

[原文連結](http://api.hslcommunication.cn/html/811ac077-29f1-f44b-7308-481786c712f6.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Delta.Helper](../html/cdc3609c-da0b-e755-d4f9-5d1bc3fa17e1.htm "HslCommunication.Profinet.Delta.Helper")

[DeltaDvpHelper 类](../html/29156f81-12f6-a557-f020-e351396152fd.htm "DeltaDvpHelper 类")

[DeltaDvpHelper 方法](../html/811ac077-29f1-f44b-7308-481786c712f6.htm "DeltaDvpHelper 方法")

[ParseDeltaDvpAddress 方法](../html/f89cf60a-cc25-7867-0f05-b21c819470e1.htm "ParseDeltaDvpAddress 方法 ")

[Read 方法](../html/b2b0c6a2-eff4-3a00-6d9d-4c91f8427c91.htm "Read 方法 ")

[ReadAsync 方法](../html/8dba9dbc-4874-0d75-9b17-7443d45092f9.htm "ReadAsync 方法 ")

[ReadBool 方法](../html/253ddbd6-6903-24ac-4f04-a7323687dc96.htm "ReadBool 方法 ")

[ReadBoolAsync 方法](../html/e1185968-3493-27d4-eefc-a8f8a70c2b5f.htm "ReadBoolAsync 方法 ")

[Write 方法](../html/61e345cf-1c92-632e-6140-3bd719327c74.htm "Write 方法 ")

[WriteAsync 方法](../html/1477c5e0-8acc-ae40-1935-2d7efe9f1bfb.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DeltaDvpHelper 方法 |

[DeltaDvpHelper](29156f81-12f6-a557-f020-e351396152fd.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法 | GetType | (继承自 Object。) |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法静态成员 | [ParseDeltaDvpAddress](f89cf60a-cc25-7867-0f05-b21c819470e1.htm) | 根据台达PLC的地址，解析出转换后的modbus协议信息，适用DVP系列，当前的地址仍然支持站号指定，例如s=2;D100  According to the address of Delta PLC, the converted modbus protocol information is parsed out, applicable to DVP series, the current address still supports station number designation, such as s=2;D100 |
| 公共方法静态成员 | [Read](b2b0c6a2-eff4-3a00-6d9d-4c91f8427c91.htm) | 读取台达PLC的原始字节变量，重写了读D地址时，跨区域读4096地址时，将会分割多次读取 |
| 公共方法静态成员 | [ReadAsync](8dba9dbc-4874-0d75-9b17-7443d45092f9.htm) | 读取台达PLC的原始字节变量，重写了读D地址时，跨区域读4096地址时，将会分割多次读取 |
| 公共方法静态成员 | [ReadBool](253ddbd6-6903-24ac-4f04-a7323687dc96.htm) | 读取台达PLC的bool变量，重写了读M地址时，跨区域读1536地址时，将会分割多次读取 |
| 公共方法静态成员 | [ReadBoolAsync](e1185968-3493-27d4-eefc-a8f8a70c2b5f.htm) | 读取台达PLC的bool变量，重写了读M地址时，跨区域读1536地址时，将会分割多次读取 |
| 公共方法 | ToString | (继承自 Object。) |
| 公共方法静态成员 | [Write(FuncString, Boolean, OperateResult, String, Boolean)](db83d633-4946-68f7-19a0-105612a919e9.htm) | 写入台达PLC的bool数据，当发现是M类型的数据，并且地址出现跨1536时，进行切割写入操作 |
| 公共方法静态成员 | [Write(FuncString, Byte, OperateResult, String, Byte)](ff119d13-d37c-183d-015f-54ec55816e55.htm) | 写入台达PLC的原始字节数据，当发现是D类型的数据，并且地址出现跨4096时，进行切割写入操作 |
| 公共方法静态成员 | [WriteAsync(FuncString, Boolean, TaskOperateResult, String, Boolean)](a303aec2-4f6c-6d34-d5bf-6a04328dc9bb.htm) | 写入台达PLC的bool数据，当发现是M类型的数据，并且地址出现跨1536时，进行切割写入操作 |
| 公共方法静态成员 | [WriteAsync(FuncString, Byte, TaskOperateResult, String, Byte)](88c953ac-9a5c-4fd9-2ca0-f20e5fbb00a5.htm) | 写入台达PLC的原始字节数据，当发现是D类型的数据，并且地址出现跨4096时，进行切割写入操作 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[DeltaDvpHelper 类](29156f81-12f6-a557-f020-e351396152fd.htm)

[HslCommunication.Profinet.Delta.Helper 命名空间](cdc3609c-da0b-e755-d4f9-5d1bc3fa17e1.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ParseDeltaDvpAddress 方法 

[原文連結](http://api.hslcommunication.cn/html/f89cf60a-cc25-7867-0f05-b21c819470e1.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Delta.Helper](../html/cdc3609c-da0b-e755-d4f9-5d1bc3fa17e1.htm "HslCommunication.Profinet.Delta.Helper")

[DeltaDvpHelper 类](../html/29156f81-12f6-a557-f020-e351396152fd.htm "DeltaDvpHelper 类")

[DeltaDvpHelper 方法](../html/811ac077-29f1-f44b-7308-481786c712f6.htm "DeltaDvpHelper 方法")

[ParseDeltaDvpAddress 方法](../html/f89cf60a-cc25-7867-0f05-b21c819470e1.htm "ParseDeltaDvpAddress 方法 ")

[Read 方法](../html/b2b0c6a2-eff4-3a00-6d9d-4c91f8427c91.htm "Read 方法 ")

[ReadAsync 方法](../html/8dba9dbc-4874-0d75-9b17-7443d45092f9.htm "ReadAsync 方法 ")

[ReadBool 方法](../html/253ddbd6-6903-24ac-4f04-a7323687dc96.htm "ReadBool 方法 ")

[ReadBoolAsync 方法](../html/e1185968-3493-27d4-eefc-a8f8a70c2b5f.htm "ReadBoolAsync 方法 ")

[Write 方法](../html/61e345cf-1c92-632e-6140-3bd719327c74.htm "Write 方法 ")

[WriteAsync 方法](../html/1477c5e0-8acc-ae40-1935-2d7efe9f1bfb.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DeltaDvpHelperParseDeltaDvpAddress 方法 |

根据台达PLC的地址，解析出转换后的modbus协议信息，适用DVP系列，当前的地址仍然支持站号指定，例如s=2;D100  
According to the address of Delta PLC, the converted modbus protocol information is parsed out, applicable to DVP series,
the current address still supports station number designation, such as s=2;D100

**命名空间：**
 [HslCommunication.Profinet.Delta.Helper](cdc3609c-da0b-e755-d4f9-5d1bc3fa17e1.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult<string> ParseDeltaDvpAddress(
	string address,
	byte modbusCode
)
```

```
Public Shared Function ParseDeltaDvpAddress ( 
	address As String,
	modbusCode As Byte
) As OperateResult(Of String)
```

```
public:
static OperateResult<String^>^ ParseDeltaDvpAddress(
	String^ address, 
	unsigned char modbusCode
)
```

```
static member ParseDeltaDvpAddress : 
        address : string * 
        modbusCode : byte -> OperateResult<string> 
```

#### 参数

address
:   类型：SystemString  
    台达plc的地址信息

modbusCode
:   类型：SystemByte  
    原始的对应的modbus信息

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)String  
还原后的modbus地址

![](../icons/SectionExpanded.png)参见

#### 引用

[DeltaDvpHelper 类](29156f81-12f6-a557-f020-e351396152fd.htm)

[HslCommunication.Profinet.Delta.Helper 命名空间](cdc3609c-da0b-e755-d4f9-5d1bc3fa17e1.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Read 方法 

[原文連結](http://api.hslcommunication.cn/html/b2b0c6a2-eff4-3a00-6d9d-4c91f8427c91.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Delta.Helper](../html/cdc3609c-da0b-e755-d4f9-5d1bc3fa17e1.htm "HslCommunication.Profinet.Delta.Helper")

[DeltaDvpHelper 类](../html/29156f81-12f6-a557-f020-e351396152fd.htm "DeltaDvpHelper 类")

[DeltaDvpHelper 方法](../html/811ac077-29f1-f44b-7308-481786c712f6.htm "DeltaDvpHelper 方法")

[ParseDeltaDvpAddress 方法](../html/f89cf60a-cc25-7867-0f05-b21c819470e1.htm "ParseDeltaDvpAddress 方法 ")

[Read 方法](../html/b2b0c6a2-eff4-3a00-6d9d-4c91f8427c91.htm "Read 方法 ")

[ReadAsync 方法](../html/8dba9dbc-4874-0d75-9b17-7443d45092f9.htm "ReadAsync 方法 ")

[ReadBool 方法](../html/253ddbd6-6903-24ac-4f04-a7323687dc96.htm "ReadBool 方法 ")

[ReadBoolAsync 方法](../html/e1185968-3493-27d4-eefc-a8f8a70c2b5f.htm "ReadBoolAsync 方法 ")

[Write 方法](../html/61e345cf-1c92-632e-6140-3bd719327c74.htm "Write 方法 ")

[WriteAsync 方法](../html/1477c5e0-8acc-ae40-1935-2d7efe9f1bfb.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DeltaDvpHelperRead 方法 |

读取台达PLC的原始字节变量，重写了读D地址时，跨区域读4096地址时，将会分割多次读取

**命名空间：**
 [HslCommunication.Profinet.Delta.Helper](cdc3609c-da0b-e755-d4f9-5d1bc3fa17e1.htm)  
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
	Func<string, ushort, OperateResult<byte[]>> readFunc,
	string address,
	ushort length
)
```

```
Public Shared Function Read ( 
	readFunc As Func(Of String, UShort, OperateResult(Of Byte())),
	address As String,
	length As UShort
) As OperateResult(Of Byte())
```

```
public:
static OperateResult<array<unsigned char>^>^ Read(
	Func<String^, unsigned short, OperateResult<array<unsigned char>^>^>^ readFunc, 
	String^ address, 
	unsigned short length
)
```

```
static member Read : 
        readFunc : Func<string, uint16, OperateResult<byte[]>> * 
        address : string * 
        length : uint16 -> OperateResult<byte[]> 
```

#### 参数

readFunc
:   类型：SystemFuncString, UInt16, [OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
    底层基础的读取方法

address
:   类型：SystemString  
    PLC的地址信息

length
:   类型：SystemUInt16  
    读取的长度信息

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
读取的结果

![](../icons/SectionExpanded.png)参见

#### 引用

[DeltaDvpHelper 类](29156f81-12f6-a557-f020-e351396152fd.htm)

[HslCommunication.Profinet.Delta.Helper 命名空间](cdc3609c-da0b-e755-d4f9-5d1bc3fa17e1.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/8dba9dbc-4874-0d75-9b17-7443d45092f9.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Delta.Helper](../html/cdc3609c-da0b-e755-d4f9-5d1bc3fa17e1.htm "HslCommunication.Profinet.Delta.Helper")

[DeltaDvpHelper 类](../html/29156f81-12f6-a557-f020-e351396152fd.htm "DeltaDvpHelper 类")

[DeltaDvpHelper 方法](../html/811ac077-29f1-f44b-7308-481786c712f6.htm "DeltaDvpHelper 方法")

[ParseDeltaDvpAddress 方法](../html/f89cf60a-cc25-7867-0f05-b21c819470e1.htm "ParseDeltaDvpAddress 方法 ")

[Read 方法](../html/b2b0c6a2-eff4-3a00-6d9d-4c91f8427c91.htm "Read 方法 ")

[ReadAsync 方法](../html/8dba9dbc-4874-0d75-9b17-7443d45092f9.htm "ReadAsync 方法 ")

[ReadBool 方法](../html/253ddbd6-6903-24ac-4f04-a7323687dc96.htm "ReadBool 方法 ")

[ReadBoolAsync 方法](../html/e1185968-3493-27d4-eefc-a8f8a70c2b5f.htm "ReadBoolAsync 方法 ")

[Write 方法](../html/61e345cf-1c92-632e-6140-3bd719327c74.htm "Write 方法 ")

[WriteAsync 方法](../html/1477c5e0-8acc-ae40-1935-2d7efe9f1bfb.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DeltaDvpHelperReadAsync 方法 |

读取台达PLC的原始字节变量，重写了读D地址时，跨区域读4096地址时，将会分割多次读取

**命名空间：**
 [HslCommunication.Profinet.Delta.Helper](cdc3609c-da0b-e755-d4f9-5d1bc3fa17e1.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static Task<OperateResult<byte[]>> ReadAsync(
	Func<string, ushort, Task<OperateResult<byte[]>>> readFunc,
	string address,
	ushort length
)
```

```
Public Shared Function ReadAsync ( 
	readFunc As Func(Of String, UShort, Task(Of OperateResult(Of Byte()))),
	address As String,
	length As UShort
) As Task(Of OperateResult(Of Byte()))
```

```
public:
static Task<OperateResult<array<unsigned char>^>^>^ ReadAsync(
	Func<String^, unsigned short, Task<OperateResult<array<unsigned char>^>^>^>^ readFunc, 
	String^ address, 
	unsigned short length
)
```

```
static member ReadAsync : 
        readFunc : Func<string, uint16, Task<OperateResult<byte[]>>> * 
        address : string * 
        length : uint16 -> Task<OperateResult<byte[]>> 
```

#### 参数

readFunc
:   类型：SystemFuncString, UInt16, Task[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
    底层基础的读取方法

address
:   类型：SystemString  
    PLC的地址信息

length
:   类型：SystemUInt16  
    读取的长度信息

#### 返回值

类型：Task[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Byte  
读取的结果

![](../icons/SectionExpanded.png)参见

#### 引用

[DeltaDvpHelper 类](29156f81-12f6-a557-f020-e351396152fd.htm)

[HslCommunication.Profinet.Delta.Helper 命名空间](cdc3609c-da0b-e755-d4f9-5d1bc3fa17e1.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadBool 方法 

[原文連結](http://api.hslcommunication.cn/html/253ddbd6-6903-24ac-4f04-a7323687dc96.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Delta.Helper](../html/cdc3609c-da0b-e755-d4f9-5d1bc3fa17e1.htm "HslCommunication.Profinet.Delta.Helper")

[DeltaDvpHelper 类](../html/29156f81-12f6-a557-f020-e351396152fd.htm "DeltaDvpHelper 类")

[DeltaDvpHelper 方法](../html/811ac077-29f1-f44b-7308-481786c712f6.htm "DeltaDvpHelper 方法")

[ParseDeltaDvpAddress 方法](../html/f89cf60a-cc25-7867-0f05-b21c819470e1.htm "ParseDeltaDvpAddress 方法 ")

[Read 方法](../html/b2b0c6a2-eff4-3a00-6d9d-4c91f8427c91.htm "Read 方法 ")

[ReadAsync 方法](../html/8dba9dbc-4874-0d75-9b17-7443d45092f9.htm "ReadAsync 方法 ")

[ReadBool 方法](../html/253ddbd6-6903-24ac-4f04-a7323687dc96.htm "ReadBool 方法 ")

[ReadBoolAsync 方法](../html/e1185968-3493-27d4-eefc-a8f8a70c2b5f.htm "ReadBoolAsync 方法 ")

[Write 方法](../html/61e345cf-1c92-632e-6140-3bd719327c74.htm "Write 方法 ")

[WriteAsync 方法](../html/1477c5e0-8acc-ae40-1935-2d7efe9f1bfb.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DeltaDvpHelperReadBool 方法 |

读取台达PLC的bool变量，重写了读M地址时，跨区域读1536地址时，将会分割多次读取

**命名空间：**
 [HslCommunication.Profinet.Delta.Helper](cdc3609c-da0b-e755-d4f9-5d1bc3fa17e1.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static OperateResult<bool[]> ReadBool(
	Func<string, ushort, OperateResult<bool[]>> readBoolFunc,
	string address,
	ushort length
)
```

```
Public Shared Function ReadBool ( 
	readBoolFunc As Func(Of String, UShort, OperateResult(Of Boolean())),
	address As String,
	length As UShort
) As OperateResult(Of Boolean())
```

```
public:
static OperateResult<array<bool>^>^ ReadBool(
	Func<String^, unsigned short, OperateResult<array<bool>^>^>^ readBoolFunc, 
	String^ address, 
	unsigned short length
)
```

```
static member ReadBool : 
        readBoolFunc : Func<string, uint16, OperateResult<bool[]>> * 
        address : string * 
        length : uint16 -> OperateResult<bool[]> 
```

#### 参数

readBoolFunc
:   类型：SystemFuncString, UInt16, [OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Boolean  
    底层基础的读取方法

address
:   类型：SystemString  
    PLC的地址信息

length
:   类型：SystemUInt16  
    读取的长度信息

#### 返回值

类型：[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Boolean  
读取的结果

![](../icons/SectionExpanded.png)参见

#### 引用

[DeltaDvpHelper 类](29156f81-12f6-a557-f020-e351396152fd.htm)

[HslCommunication.Profinet.Delta.Helper 命名空间](cdc3609c-da0b-e755-d4f9-5d1bc3fa17e1.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReadBoolAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/e1185968-3493-27d4-eefc-a8f8a70c2b5f.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Delta.Helper](../html/cdc3609c-da0b-e755-d4f9-5d1bc3fa17e1.htm "HslCommunication.Profinet.Delta.Helper")

[DeltaDvpHelper 类](../html/29156f81-12f6-a557-f020-e351396152fd.htm "DeltaDvpHelper 类")

[DeltaDvpHelper 方法](../html/811ac077-29f1-f44b-7308-481786c712f6.htm "DeltaDvpHelper 方法")

[ParseDeltaDvpAddress 方法](../html/f89cf60a-cc25-7867-0f05-b21c819470e1.htm "ParseDeltaDvpAddress 方法 ")

[Read 方法](../html/b2b0c6a2-eff4-3a00-6d9d-4c91f8427c91.htm "Read 方法 ")

[ReadAsync 方法](../html/8dba9dbc-4874-0d75-9b17-7443d45092f9.htm "ReadAsync 方法 ")

[ReadBool 方法](../html/253ddbd6-6903-24ac-4f04-a7323687dc96.htm "ReadBool 方法 ")

[ReadBoolAsync 方法](../html/e1185968-3493-27d4-eefc-a8f8a70c2b5f.htm "ReadBoolAsync 方法 ")

[Write 方法](../html/61e345cf-1c92-632e-6140-3bd719327c74.htm "Write 方法 ")

[WriteAsync 方法](../html/1477c5e0-8acc-ae40-1935-2d7efe9f1bfb.htm "WriteAsync 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DeltaDvpHelperReadBoolAsync 方法 |

读取台达PLC的bool变量，重写了读M地址时，跨区域读1536地址时，将会分割多次读取

**命名空间：**
 [HslCommunication.Profinet.Delta.Helper](cdc3609c-da0b-e755-d4f9-5d1bc3fa17e1.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static Task<OperateResult<bool[]>> ReadBoolAsync(
	Func<string, ushort, Task<OperateResult<bool[]>>> readBoolFunc,
	string address,
	ushort length
)
```

```
Public Shared Function ReadBoolAsync ( 
	readBoolFunc As Func(Of String, UShort, Task(Of OperateResult(Of Boolean()))),
	address As String,
	length As UShort
) As Task(Of OperateResult(Of Boolean()))
```

```
public:
static Task<OperateResult<array<bool>^>^>^ ReadBoolAsync(
	Func<String^, unsigned short, Task<OperateResult<array<bool>^>^>^>^ readBoolFunc, 
	String^ address, 
	unsigned short length
)
```

```
static member ReadBoolAsync : 
        readBoolFunc : Func<string, uint16, Task<OperateResult<bool[]>>> * 
        address : string * 
        length : uint16 -> Task<OperateResult<bool[]>> 
```

#### 参数

readBoolFunc
:   类型：SystemFuncString, UInt16, Task[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Boolean  
    底层基础的读取方法

address
:   类型：SystemString  
    PLC的地址信息

length
:   类型：SystemUInt16  
    读取的长度信息

#### 返回值

类型：Task[OperateResult](cf3af065-3e0e-e0ad-7063-8d7eab691b48.htm)Boolean  
读取的结果

![](../icons/SectionExpanded.png)参见

#### 引用

[DeltaDvpHelper 类](29156f81-12f6-a557-f020-e351396152fd.htm)

[HslCommunication.Profinet.Delta.Helper 命名空间](cdc3609c-da0b-e755-d4f9-5d1bc3fa17e1.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Write 方法 

[原文連結](http://api.hslcommunication.cn/html/61e345cf-1c92-632e-6140-3bd719327c74.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Delta.Helper](../html/cdc3609c-da0b-e755-d4f9-5d1bc3fa17e1.htm "HslCommunication.Profinet.Delta.Helper")

[DeltaDvpHelper 类](../html/29156f81-12f6-a557-f020-e351396152fd.htm "DeltaDvpHelper 类")

[DeltaDvpHelper 方法](../html/811ac077-29f1-f44b-7308-481786c712f6.htm "DeltaDvpHelper 方法")

[Write 方法](../html/61e345cf-1c92-632e-6140-3bd719327c74.htm "Write 方法 ")

[Write 方法 (Func(String, Boolean[], OperateResult), String, Boolean[])](../html/db83d633-4946-68f7-19a0-105612a919e9.htm "Write 方法 (Func(String, Boolean[], OperateResult), String, Boolean[])")

[Write 方法 (Func(String, Byte[], OperateResult), String, Byte[])](../html/ff119d13-d37c-183d-015f-54ec55816e55.htm "Write 方法 (Func(String, Byte[], OperateResult), String, Byte[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DeltaDvpHelperWrite 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法静态成员 | [Write(FuncString, Boolean, OperateResult, String, Boolean)](db83d633-4946-68f7-19a0-105612a919e9.htm) | 写入台达PLC的bool数据，当发现是M类型的数据，并且地址出现跨1536时，进行切割写入操作 |
| 公共方法静态成员 | [Write(FuncString, Byte, OperateResult, String, Byte)](ff119d13-d37c-183d-015f-54ec55816e55.htm) | 写入台达PLC的原始字节数据，当发现是D类型的数据，并且地址出现跨4096时，进行切割写入操作 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[DeltaDvpHelper 类](29156f81-12f6-a557-f020-e351396152fd.htm)

[HslCommunication.Profinet.Delta.Helper 命名空间](cdc3609c-da0b-e755-d4f9-5d1bc3fa17e1.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Write 方法 (Func(String, Boolean[], OperateResult), String, Boolean[])

[原文連結](http://api.hslcommunication.cn/html/db83d633-4946-68f7-19a0-105612a919e9.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Delta.Helper](../html/cdc3609c-da0b-e755-d4f9-5d1bc3fa17e1.htm "HslCommunication.Profinet.Delta.Helper")

[DeltaDvpHelper 类](../html/29156f81-12f6-a557-f020-e351396152fd.htm "DeltaDvpHelper 类")

[DeltaDvpHelper 方法](../html/811ac077-29f1-f44b-7308-481786c712f6.htm "DeltaDvpHelper 方法")

[Write 方法](../html/61e345cf-1c92-632e-6140-3bd719327c74.htm "Write 方法 ")

[Write 方法 (Func(String, Boolean[], OperateResult), String, Boolean[])](../html/db83d633-4946-68f7-19a0-105612a919e9.htm "Write 方法 (Func(String, Boolean[], OperateResult), String, Boolean[])")

[Write 方法 (Func(String, Byte[], OperateResult), String, Byte[])](../html/ff119d13-d37c-183d-015f-54ec55816e55.htm "Write 方法 (Func(String, Byte[], OperateResult), String, Byte[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DeltaDvpHelperWrite 方法 (FuncString, Boolean, OperateResult, String, Boolean) |

写入台达PLC的bool数据，当发现是M类型的数据，并且地址出现跨1536时，进行切割写入操作

**命名空间：**
 [HslCommunication.Profinet.Delta.Helper](cdc3609c-da0b-e755-d4f9-5d1bc3fa17e1.htm)  
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
	Func<string, bool[], OperateResult> writeBoolFunc,
	string address,
	bool[] value
)
```

```
Public Shared Function Write ( 
	writeBoolFunc As Func(Of String, Boolean(), OperateResult),
	address As String,
	value As Boolean()
) As OperateResult
```

```
public:
static OperateResult^ Write(
	Func<String^, array<bool>^, OperateResult^>^ writeBoolFunc, 
	String^ address, 
	array<bool>^ value
)
```

```
static member Write : 
        writeBoolFunc : Func<string, bool[], OperateResult> * 
        address : string * 
        value : bool[] -> OperateResult 
```

#### 参数

writeBoolFunc
:   类型：SystemFuncString, Boolean, [OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
    底层的写入操作方法

address
:   类型：SystemString  
    PLC的起始地址信息

value
:   类型：SystemBoolean  
    等待写入的数据信息

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否写入成功

![](../icons/SectionExpanded.png)参见

#### 引用

[DeltaDvpHelper 类](29156f81-12f6-a557-f020-e351396152fd.htm)

[Write 重载](61e345cf-1c92-632e-6140-3bd719327c74.htm)

[HslCommunication.Profinet.Delta.Helper 命名空间](cdc3609c-da0b-e755-d4f9-5d1bc3fa17e1.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Write 方法 (Func(String, Byte[], OperateResult), String, Byte[])

[原文連結](http://api.hslcommunication.cn/html/ff119d13-d37c-183d-015f-54ec55816e55.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Delta.Helper](../html/cdc3609c-da0b-e755-d4f9-5d1bc3fa17e1.htm "HslCommunication.Profinet.Delta.Helper")

[DeltaDvpHelper 类](../html/29156f81-12f6-a557-f020-e351396152fd.htm "DeltaDvpHelper 类")

[DeltaDvpHelper 方法](../html/811ac077-29f1-f44b-7308-481786c712f6.htm "DeltaDvpHelper 方法")

[Write 方法](../html/61e345cf-1c92-632e-6140-3bd719327c74.htm "Write 方法 ")

[Write 方法 (Func(String, Boolean[], OperateResult), String, Boolean[])](../html/db83d633-4946-68f7-19a0-105612a919e9.htm "Write 方法 (Func(String, Boolean[], OperateResult), String, Boolean[])")

[Write 方法 (Func(String, Byte[], OperateResult), String, Byte[])](../html/ff119d13-d37c-183d-015f-54ec55816e55.htm "Write 方法 (Func(String, Byte[], OperateResult), String, Byte[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DeltaDvpHelperWrite 方法 (FuncString, Byte, OperateResult, String, Byte) |

写入台达PLC的原始字节数据，当发现是D类型的数据，并且地址出现跨4096时，进行切割写入操作

**命名空间：**
 [HslCommunication.Profinet.Delta.Helper](cdc3609c-da0b-e755-d4f9-5d1bc3fa17e1.htm)  
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
	Func<string, byte[], OperateResult> writeFunc,
	string address,
	byte[] value
)
```

```
Public Shared Function Write ( 
	writeFunc As Func(Of String, Byte(), OperateResult),
	address As String,
	value As Byte()
) As OperateResult
```

```
public:
static OperateResult^ Write(
	Func<String^, array<unsigned char>^, OperateResult^>^ writeFunc, 
	String^ address, 
	array<unsigned char>^ value
)
```

```
static member Write : 
        writeFunc : Func<string, byte[], OperateResult> * 
        address : string * 
        value : byte[] -> OperateResult 
```

#### 参数

writeFunc
:   类型：SystemFuncString, Byte, [OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
    底层的写入操作方法

address
:   类型：SystemString  
    PLC的起始地址信息

value
:   类型：SystemByte  
    等待写入的数据信息

#### 返回值

类型：[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否写入成功

![](../icons/SectionExpanded.png)参见

#### 引用

[DeltaDvpHelper 类](29156f81-12f6-a557-f020-e351396152fd.htm)

[Write 重载](61e345cf-1c92-632e-6140-3bd719327c74.htm)

[HslCommunication.Profinet.Delta.Helper 命名空间](cdc3609c-da0b-e755-d4f9-5d1bc3fa17e1.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## WriteAsync 方法 

[原文連結](http://api.hslcommunication.cn/html/1477c5e0-8acc-ae40-1935-2d7efe9f1bfb.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Delta.Helper](../html/cdc3609c-da0b-e755-d4f9-5d1bc3fa17e1.htm "HslCommunication.Profinet.Delta.Helper")

[DeltaDvpHelper 类](../html/29156f81-12f6-a557-f020-e351396152fd.htm "DeltaDvpHelper 类")

[DeltaDvpHelper 方法](../html/811ac077-29f1-f44b-7308-481786c712f6.htm "DeltaDvpHelper 方法")

[WriteAsync 方法](../html/1477c5e0-8acc-ae40-1935-2d7efe9f1bfb.htm "WriteAsync 方法 ")

[WriteAsync 方法 (Func(String, Boolean[], Task(OperateResult)), String, Boolean[])](../html/a303aec2-4f6c-6d34-d5bf-6a04328dc9bb.htm "WriteAsync 方法 (Func(String, Boolean[], Task(OperateResult)), String, Boolean[])")

[WriteAsync 方法 (Func(String, Byte[], Task(OperateResult)), String, Byte[])](../html/88c953ac-9a5c-4fd9-2ca0-f20e5fbb00a5.htm "WriteAsync 方法 (Func(String, Byte[], Task(OperateResult)), String, Byte[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DeltaDvpHelperWriteAsync 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法静态成员 | [WriteAsync(FuncString, Boolean, TaskOperateResult, String, Boolean)](a303aec2-4f6c-6d34-d5bf-6a04328dc9bb.htm) | 写入台达PLC的bool数据，当发现是M类型的数据，并且地址出现跨1536时，进行切割写入操作 |
| 公共方法静态成员 | [WriteAsync(FuncString, Byte, TaskOperateResult, String, Byte)](88c953ac-9a5c-4fd9-2ca0-f20e5fbb00a5.htm) | 写入台达PLC的原始字节数据，当发现是D类型的数据，并且地址出现跨4096时，进行切割写入操作 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[DeltaDvpHelper 类](29156f81-12f6-a557-f020-e351396152fd.htm)

[HslCommunication.Profinet.Delta.Helper 命名空间](cdc3609c-da0b-e755-d4f9-5d1bc3fa17e1.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## WriteAsync 方法 (Func(String, Boolean[], Task(OperateResult)), String, Boolean[])

[原文連結](http://api.hslcommunication.cn/html/a303aec2-4f6c-6d34-d5bf-6a04328dc9bb.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Delta.Helper](../html/cdc3609c-da0b-e755-d4f9-5d1bc3fa17e1.htm "HslCommunication.Profinet.Delta.Helper")

[DeltaDvpHelper 类](../html/29156f81-12f6-a557-f020-e351396152fd.htm "DeltaDvpHelper 类")

[DeltaDvpHelper 方法](../html/811ac077-29f1-f44b-7308-481786c712f6.htm "DeltaDvpHelper 方法")

[WriteAsync 方法](../html/1477c5e0-8acc-ae40-1935-2d7efe9f1bfb.htm "WriteAsync 方法 ")

[WriteAsync 方法 (Func(String, Boolean[], Task(OperateResult)), String, Boolean[])](../html/a303aec2-4f6c-6d34-d5bf-6a04328dc9bb.htm "WriteAsync 方法 (Func(String, Boolean[], Task(OperateResult)), String, Boolean[])")

[WriteAsync 方法 (Func(String, Byte[], Task(OperateResult)), String, Byte[])](../html/88c953ac-9a5c-4fd9-2ca0-f20e5fbb00a5.htm "WriteAsync 方法 (Func(String, Byte[], Task(OperateResult)), String, Byte[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DeltaDvpHelperWriteAsync 方法 (FuncString, Boolean, TaskOperateResult, String, Boolean) |

写入台达PLC的bool数据，当发现是M类型的数据，并且地址出现跨1536时，进行切割写入操作

**命名空间：**
 [HslCommunication.Profinet.Delta.Helper](cdc3609c-da0b-e755-d4f9-5d1bc3fa17e1.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static Task<OperateResult> WriteAsync(
	Func<string, bool[], Task<OperateResult>> writeBoolFunc,
	string address,
	bool[] value
)
```

```
Public Shared Function WriteAsync ( 
	writeBoolFunc As Func(Of String, Boolean(), Task(Of OperateResult)),
	address As String,
	value As Boolean()
) As Task(Of OperateResult)
```

```
public:
static Task<OperateResult^>^ WriteAsync(
	Func<String^, array<bool>^, Task<OperateResult^>^>^ writeBoolFunc, 
	String^ address, 
	array<bool>^ value
)
```

```
static member WriteAsync : 
        writeBoolFunc : Func<string, bool[], Task<OperateResult>> * 
        address : string * 
        value : bool[] -> Task<OperateResult> 
```

#### 参数

writeBoolFunc
:   类型：SystemFuncString, Boolean, Task[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
    底层的写入操作方法

address
:   类型：SystemString  
    PLC的起始地址信息

value
:   类型：SystemBoolean  
    等待写入的数据信息

#### 返回值

类型：Task[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否写入成功

![](../icons/SectionExpanded.png)参见

#### 引用

[DeltaDvpHelper 类](29156f81-12f6-a557-f020-e351396152fd.htm)

[WriteAsync 重载](1477c5e0-8acc-ae40-1935-2d7efe9f1bfb.htm)

[HslCommunication.Profinet.Delta.Helper 命名空间](cdc3609c-da0b-e755-d4f9-5d1bc3fa17e1.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## WriteAsync 方法 (Func(String, Byte[], Task(OperateResult)), String, Byte[])

[原文連結](http://api.hslcommunication.cn/html/88c953ac-9a5c-4fd9-2ca0-f20e5fbb00a5.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Delta.Helper](../html/cdc3609c-da0b-e755-d4f9-5d1bc3fa17e1.htm "HslCommunication.Profinet.Delta.Helper")

[DeltaDvpHelper 类](../html/29156f81-12f6-a557-f020-e351396152fd.htm "DeltaDvpHelper 类")

[DeltaDvpHelper 方法](../html/811ac077-29f1-f44b-7308-481786c712f6.htm "DeltaDvpHelper 方法")

[WriteAsync 方法](../html/1477c5e0-8acc-ae40-1935-2d7efe9f1bfb.htm "WriteAsync 方法 ")

[WriteAsync 方法 (Func(String, Boolean[], Task(OperateResult)), String, Boolean[])](../html/a303aec2-4f6c-6d34-d5bf-6a04328dc9bb.htm "WriteAsync 方法 (Func(String, Boolean[], Task(OperateResult)), String, Boolean[])")

[WriteAsync 方法 (Func(String, Byte[], Task(OperateResult)), String, Byte[])](../html/88c953ac-9a5c-4fd9-2ca0-f20e5fbb00a5.htm "WriteAsync 方法 (Func(String, Byte[], Task(OperateResult)), String, Byte[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DeltaDvpHelperWriteAsync 方法 (FuncString, Byte, TaskOperateResult, String, Byte) |

写入台达PLC的原始字节数据，当发现是D类型的数据，并且地址出现跨4096时，进行切割写入操作

**命名空间：**
 [HslCommunication.Profinet.Delta.Helper](cdc3609c-da0b-e755-d4f9-5d1bc3fa17e1.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static Task<OperateResult> WriteAsync(
	Func<string, byte[], Task<OperateResult>> writeFunc,
	string address,
	byte[] value
)
```

```
Public Shared Function WriteAsync ( 
	writeFunc As Func(Of String, Byte(), Task(Of OperateResult)),
	address As String,
	value As Byte()
) As Task(Of OperateResult)
```

```
public:
static Task<OperateResult^>^ WriteAsync(
	Func<String^, array<unsigned char>^, Task<OperateResult^>^>^ writeFunc, 
	String^ address, 
	array<unsigned char>^ value
)
```

```
static member WriteAsync : 
        writeFunc : Func<string, byte[], Task<OperateResult>> * 
        address : string * 
        value : byte[] -> Task<OperateResult> 
```

#### 参数

writeFunc
:   类型：SystemFuncString, Byte, Task[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
    底层的写入操作方法

address
:   类型：SystemString  
    PLC的起始地址信息

value
:   类型：SystemByte  
    等待写入的数据信息

#### 返回值

类型：Task[OperateResult](a8ab4838-9351-2ebe-b761-3b0fddbcda6b.htm)  
是否写入成功

![](../icons/SectionExpanded.png)参见

#### 引用

[DeltaDvpHelper 类](29156f81-12f6-a557-f020-e351396152fd.htm)

[WriteAsync 重载](1477c5e0-8acc-ae40-1935-2d7efe9f1bfb.htm)

[HslCommunication.Profinet.Delta.Helper 命名空间](cdc3609c-da0b-e755-d4f9-5d1bc3fa17e1.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## DeltaHelper 类

[原文連結](http://api.hslcommunication.cn/html/ffea0bf8-6cae-1863-da2c-e772ad28c54e.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Delta.Helper](../html/cdc3609c-da0b-e755-d4f9-5d1bc3fa17e1.htm "HslCommunication.Profinet.Delta.Helper")

[DeltaHelper 类](../html/ffea0bf8-6cae-1863-da2c-e772ad28c54e.htm "DeltaHelper 类")

[DeltaHelper 构造函数](../html/960b806d-825a-87d0-5b73-89bfaaf6ef1a.htm "DeltaHelper 构造函数 ")

[DeltaHelper 方法](../html/ebca4d5a-a922-3ebf-180a-08b1c1769423.htm "DeltaHelper 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DeltaHelper 类 |

台达的想关的辅助类

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  HslCommunication.Profinet.Delta.HelperDeltaHelper

**命名空间：**
 [HslCommunication.Profinet.Delta.Helper](cdc3609c-da0b-e755-d4f9-5d1bc3fa17e1.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class DeltaHelper
```

```
Public Class DeltaHelper
```

```
public ref class DeltaHelper
```

```
type DeltaHelper =  class end
```

DeltaHelper 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [DeltaHelper](960b806d-825a-87d0-5b73-89bfaaf6ef1a.htm) | 初始化 DeltaHelper 类的一个新实例 |

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

[HslCommunication.Profinet.Delta.Helper 命名空间](cdc3609c-da0b-e755-d4f9-5d1bc3fa17e1.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## DeltaHelper 构造函数 

[原文連結](http://api.hslcommunication.cn/html/960b806d-825a-87d0-5b73-89bfaaf6ef1a.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Delta.Helper](../html/cdc3609c-da0b-e755-d4f9-5d1bc3fa17e1.htm "HslCommunication.Profinet.Delta.Helper")

[DeltaHelper 类](../html/ffea0bf8-6cae-1863-da2c-e772ad28c54e.htm "DeltaHelper 类")

[DeltaHelper 构造函数](../html/960b806d-825a-87d0-5b73-89bfaaf6ef1a.htm "DeltaHelper 构造函数 ")

[DeltaHelper 方法](../html/ebca4d5a-a922-3ebf-180a-08b1c1769423.htm "DeltaHelper 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DeltaHelper 构造函数 |

初始化 [DeltaHelper](ffea0bf8-6cae-1863-da2c-e772ad28c54e.htm) 类的一个新实例

**命名空间：**
 [HslCommunication.Profinet.Delta.Helper](cdc3609c-da0b-e755-d4f9-5d1bc3fa17e1.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public DeltaHelper()
```

```
Public Sub New
```

```
public:
DeltaHelper()
```

```
new : unit -> DeltaHelper
```

![](../icons/SectionExpanded.png)参见

#### 引用

[DeltaHelper 类](ffea0bf8-6cae-1863-da2c-e772ad28c54e.htm)

[HslCommunication.Profinet.Delta.Helper 命名空间](cdc3609c-da0b-e755-d4f9-5d1bc3fa17e1.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## DeltaHelper 方法

[原文連結](http://api.hslcommunication.cn/html/ebca4d5a-a922-3ebf-180a-08b1c1769423.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.Delta.Helper](../html/cdc3609c-da0b-e755-d4f9-5d1bc3fa17e1.htm "HslCommunication.Profinet.Delta.Helper")

[DeltaHelper 类](../html/ffea0bf8-6cae-1863-da2c-e772ad28c54e.htm "DeltaHelper 类")

[DeltaHelper 构造函数](../html/960b806d-825a-87d0-5b73-89bfaaf6ef1a.htm "DeltaHelper 构造函数 ")

[DeltaHelper 方法](../html/ebca4d5a-a922-3ebf-180a-08b1c1769423.htm "DeltaHelper 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| DeltaHelper 方法 |

[DeltaHelper](ffea0bf8-6cae-1863-da2c-e772ad28c54e.htm) 类型公开以下成员。

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

[DeltaHelper 类](ffea0bf8-6cae-1863-da2c-e772ad28c54e.htm)

[HslCommunication.Profinet.Delta.Helper 命名空间](cdc3609c-da0b-e755-d4f9-5d1bc3fa17e1.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)