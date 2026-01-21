# HslCommunication - HslCommunication.Profinet.AllenBradley

> 分類頁數: 30



---
## HslCommunication.Profinet.AllenBradley

[原文連結](http://api.hslcommunication.cn/html/fe226e46-d0b0-8ed5-5051-c925f25529d4.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.AllenBradley](../html/fe226e46-d0b0-8ed5-5051-c925f25529d4.htm "HslCommunication.Profinet.AllenBradley")

[AbStructHandle 类](../html/7edf9ded-9d79-9b6d-9cce-63bc65b21a72.htm "AbStructHandle 类")

[AbTagItem 类](../html/13d4e75c-bbdc-08cf-426f-799bb4099deb.htm "AbTagItem 类")

[AbTagValue 类](../html/53107c68-a30b-eb91-77de-a01b4880ba67.htm "AbTagValue 类")

[AllenBradleyConnectedCipNet 类](../html/b212ddf4-4eb6-e300-1d3e-cee8cdb4b184.htm "AllenBradleyConnectedCipNet 类")

[AllenBradleyDF1Serial 类](../html/2ac292da-a6e9-5680-0792-894dd2fc5cb9.htm "AllenBradleyDF1Serial 类")

[AllenBradleyHelper 类](../html/64e310fb-5dbb-1711-0cf1-5f64c4598205.htm "AllenBradleyHelper 类")

[AllenBradleyItemValue 类](../html/98562de3-0556-e3e3-50fe-ee4e8b4b165f.htm "AllenBradleyItemValue 类")

[AllenBradleyMicroCip 类](../html/17bfe40c-d284-cf2b-c183-5f3fb2332106.htm "AllenBradleyMicroCip 类")

[AllenBradleyNet 类](../html/149ba939-1936-de99-4987-20ec87ae9c17.htm "AllenBradleyNet 类")

[AllenBradleyPcccNet 类](../html/032dd112-7a1e-d816-ad09-1a4a4e7283df.htm "AllenBradleyPcccNet 类")

[AllenBradleyPcccServer 类](../html/1e0bd557-70ae-1aa4-0969-d8011bed3fc5.htm "AllenBradleyPcccServer 类")

[AllenBradleyServer 类](../html/c228b7ec-84a5-c3fa-a214-4ba6a3cb7110.htm "AllenBradleyServer 类")

[AllenBradleySLCNet 类](../html/fc3dfc16-849f-5245-eac5-4aaec29998ff.htm "AllenBradleySLCNet 类")

[IReadWriteCip 接口](../html/2cc21fd0-71e1-03d9-f451-60616f4f083c.htm "IReadWriteCip 接口")

[MessageRouter 类](../html/e711ce6f-4660-c0c5-92e5-8c01c024b9a6.htm "MessageRouter 类")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| HslCommunication.Profinet.AllenBradley 命名空间 |

[缺少 "N:HslCommunication.Profinet.AllenBradley" 的 <summary> 文档]

![](../icons/SectionExpanded.png)类

|  | 类 | 说明 |
| --- | --- | --- |
| 公共类 | [AbStructHandle](7edf9ded-9d79-9b6d-9cce-63bc65b21a72.htm) | 结构体的句柄信息 |
| 公共类 | [AbTagItem](13d4e75c-bbdc-08cf-426f-799bb4099deb.htm) | AB PLC的数据标签实体类  Data tag entity class of AB PLC |
| 公共类 | [AbTagValue](53107c68-a30b-eb91-77de-a01b4880ba67.htm) | AB标签的数据信息 |
| 公共类代码示例 | [AllenBradleyConnectedCipNet](b212ddf4-4eb6-e300-1d3e-cee8cdb4b184.htm) | 基于连接的对象访问的CIP协议的实现，用于对罗克韦尔 PLC进行标签的数据读写，对数组，多维数组进行读写操作，支持的数据类型请参照API文档手册。  The implementation of the CIP protocol based on connected object access is used to read and write tag data to Rockwell PLC, and read and write arrays and multidimensional arrays. For the supported data types, please refer to the API documentation manual. |
| 公共类 | [AllenBradleyDF1Serial](2ac292da-a6e9-5680-0792-894dd2fc5cb9.htm) | AB-PLC的DF1通信协议，基于串口实现，通信机制为半双工，目前适用于 Micro-Logix1000,SLC500,SLC 5/03,SLC 5/04，地址示例：N7:1 |
| 公共类 | [AllenBradleyHelper](64e310fb-5dbb-1711-0cf1-5f64c4598205.htm) | AB PLC的辅助类，用来辅助生成基本的指令信息 |
| 公共类 | [AllenBradleyItemValue](98562de3-0556-e3e3-50fe-ee4e8b4b165f.htm) | AB PLC的标签节点数据信息 |
| 公共类 | [AllenBradleyMicroCip](17bfe40c-d284-cf2b-c183-5f3fb2332106.htm) | AB PLC的cip通信实现类，适用Micro800系列控制系统  AB PLC's cip communication implementation class, suitable for Micro800 series control system |
| 公共类 | [AllenBradleyNet](149ba939-1936-de99-4987-20ec87ae9c17.htm) | AB PLC的数据通信类，使用CIP协议实现，适用1756，1769等型号，支持使用标签的形式进行读写操作，支持标量数据，一维数组，二维数组，三维数组等等。如果是局部变量，那么使用 Program:MainProgram.[变量名]。  The data communication class of AB PLC is implemented using the CIP protocol. It is suitable for 1756, 1769 and other models. It supports reading and writing in the form of tags, scalar data, one-dimensional array, two-dimensional array, three-dimensional array, and so on. If it is a local variable, use the Program:MainProgram.[Variable name]. |
| 公共类代码示例 | [AllenBradleyPcccNet](032dd112-7a1e-d816-ad09-1a4a4e7283df.htm) | 在CIP协议里，使用PCCC命令进行访问设备的原始数据文件的通信方法， |
| 公共类 | [AllenBradleyPcccServer](1e0bd557-70ae-1aa4-0969-d8011bed3fc5.htm) | 虚拟的PCCC服务器，模拟的AB 1400通信 |
| 公共类 | [AllenBradleyServer](c228b7ec-84a5-c3fa-a214-4ba6a3cb7110.htm) | AB PLC的虚拟服务器，仅支持和HSL组件的完美通信，可以手动添加一些节点。  AB PLC's virtual server only supports perfect communication with HSL components. You can manually add some nodes. |
| 公共类代码示例 | [AllenBradleySLCNet](fc3dfc16-849f-5245-eac5-4aaec29998ff.htm) | AllenBradley品牌的PLC，针对SLC系列的通信的实现，测试PLC为1747。  AllenBradley brand PLC, for the realization of SLC series communication, the test PLC is 1747. |
| 公共类 | [MessageRouter](e711ce6f-4660-c0c5-92e5-8c01c024b9a6.htm) | 自定义的消息路由类，可以实现CIP协议自定义的路由消息  A custom message routing class that can implement custom routing messages of the CIP protocol |

![](../icons/SectionExpanded.png)接口

|  | 接口 | 说明 |
| --- | --- | --- |
| 公共接口 | [IReadWriteCip](2cc21fd0-71e1-03d9-f451-60616f4f083c.htm) | CIP协议的基础接口信息 |

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## AbStructHandle 类

[原文連結](http://api.hslcommunication.cn/html/7edf9ded-9d79-9b6d-9cce-63bc65b21a72.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.AllenBradley](../html/fe226e46-d0b0-8ed5-5051-c925f25529d4.htm "HslCommunication.Profinet.AllenBradley")

[AbStructHandle 类](../html/7edf9ded-9d79-9b6d-9cce-63bc65b21a72.htm "AbStructHandle 类")

[AbStructHandle 构造函数](../html/89182ea3-2ecb-e00c-b23a-502edc9b06cf.htm "AbStructHandle 构造函数 ")

[AbStructHandle 属性](../html/920cd83c-31b1-5cda-3afa-0c9caa05930d.htm "AbStructHandle 属性")

[AbStructHandle 方法](../html/2265386f-1c3c-527d-5d32-df0680221835.htm "AbStructHandle 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| AbStructHandle 类 |

结构体的句柄信息

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  HslCommunication.Profinet.AllenBradleyAbStructHandle

**命名空间：**
 [HslCommunication.Profinet.AllenBradley](fe226e46-d0b0-8ed5-5051-c925f25529d4.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class AbStructHandle
```

```
Public Class AbStructHandle
```

```
public ref class AbStructHandle
```

```
type AbStructHandle =  class end
```

AbStructHandle 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [AbStructHandle](96d0e6a0-d87c-6904-620f-61f12e2a360c.htm) | 实例化一个默认的对象  instantiate a default object |
| 公共方法 | [AbStructHandle(Byte, Int32)](722ae95a-0a89-6493-d622-a99a692b0041.htm) | 使用原始字节的数据，索引信息来实例化一个对象  Instantiate an object with raw bytes of data, index information |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [MemberCount](0fde9476-79d7-7e33-8d98-a489e0f996c0.htm) | 成员数量 |
| 公共属性 | [ReturnCount](b3af0944-0bbf-f834-f0f8-60f0ca67ea5f.htm) | 返回项数 |
| 公共属性 | [StructureHandle](28a42c12-3a6f-84a7-cee4-23878dd52c80.htm) | 结构体的handle |
| 公共属性 | [TemplateObjectDefinitionSize](e81a9baf-1b1a-4fa1-ae0e-040235eb5f02.htm) | 结构体定义大小 |
| 公共属性 | [TemplateStructureSize](7c1de037-60f8-c0ae-73b5-3c72bb3dc970.htm) | 使用读取标记服务读取结构时在线路上传输的字节数 |

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

[HslCommunication.Profinet.AllenBradley 命名空间](fe226e46-d0b0-8ed5-5051-c925f25529d4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## AbStructHandle 构造函数 

[原文連結](http://api.hslcommunication.cn/html/89182ea3-2ecb-e00c-b23a-502edc9b06cf.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.AllenBradley](../html/fe226e46-d0b0-8ed5-5051-c925f25529d4.htm "HslCommunication.Profinet.AllenBradley")

[AbStructHandle 类](../html/7edf9ded-9d79-9b6d-9cce-63bc65b21a72.htm "AbStructHandle 类")

[AbStructHandle 构造函数](../html/89182ea3-2ecb-e00c-b23a-502edc9b06cf.htm "AbStructHandle 构造函数 ")

[AbStructHandle 构造函数](../html/96d0e6a0-d87c-6904-620f-61f12e2a360c.htm "AbStructHandle 构造函数 ")

[AbStructHandle 构造函数 (Byte[], Int32)](../html/722ae95a-0a89-6493-d622-a99a692b0041.htm "AbStructHandle 构造函数 (Byte[], Int32)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| AbStructHandle 构造函数 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [AbStructHandle](96d0e6a0-d87c-6904-620f-61f12e2a360c.htm) | 实例化一个默认的对象  instantiate a default object |
| 公共方法 | [AbStructHandle(Byte, Int32)](722ae95a-0a89-6493-d622-a99a692b0041.htm) | 使用原始字节的数据，索引信息来实例化一个对象  Instantiate an object with raw bytes of data, index information |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[AbStructHandle 类](7edf9ded-9d79-9b6d-9cce-63bc65b21a72.htm)

[HslCommunication.Profinet.AllenBradley 命名空间](fe226e46-d0b0-8ed5-5051-c925f25529d4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## AbStructHandle 构造函数 

[原文連結](http://api.hslcommunication.cn/html/96d0e6a0-d87c-6904-620f-61f12e2a360c.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.AllenBradley](../html/fe226e46-d0b0-8ed5-5051-c925f25529d4.htm "HslCommunication.Profinet.AllenBradley")

[AbStructHandle 类](../html/7edf9ded-9d79-9b6d-9cce-63bc65b21a72.htm "AbStructHandle 类")

[AbStructHandle 构造函数](../html/89182ea3-2ecb-e00c-b23a-502edc9b06cf.htm "AbStructHandle 构造函数 ")

[AbStructHandle 构造函数](../html/96d0e6a0-d87c-6904-620f-61f12e2a360c.htm "AbStructHandle 构造函数 ")

[AbStructHandle 构造函数 (Byte[], Int32)](../html/722ae95a-0a89-6493-d622-a99a692b0041.htm "AbStructHandle 构造函数 (Byte[], Int32)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| AbStructHandle 构造函数 |

实例化一个默认的对象  
instantiate a default object

**命名空间：**
 [HslCommunication.Profinet.AllenBradley](fe226e46-d0b0-8ed5-5051-c925f25529d4.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public AbStructHandle()
```

```
Public Sub New
```

```
public:
AbStructHandle()
```

```
new : unit -> AbStructHandle
```

![](../icons/SectionExpanded.png)参见

#### 引用

[AbStructHandle 类](7edf9ded-9d79-9b6d-9cce-63bc65b21a72.htm)

[AbStructHandle 重载](89182ea3-2ecb-e00c-b23a-502edc9b06cf.htm)

[HslCommunication.Profinet.AllenBradley 命名空间](fe226e46-d0b0-8ed5-5051-c925f25529d4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## AbStructHandle 构造函数 (Byte[], Int32)

[原文連結](http://api.hslcommunication.cn/html/722ae95a-0a89-6493-d622-a99a692b0041.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.AllenBradley](../html/fe226e46-d0b0-8ed5-5051-c925f25529d4.htm "HslCommunication.Profinet.AllenBradley")

[AbStructHandle 类](../html/7edf9ded-9d79-9b6d-9cce-63bc65b21a72.htm "AbStructHandle 类")

[AbStructHandle 构造函数](../html/89182ea3-2ecb-e00c-b23a-502edc9b06cf.htm "AbStructHandle 构造函数 ")

[AbStructHandle 构造函数](../html/96d0e6a0-d87c-6904-620f-61f12e2a360c.htm "AbStructHandle 构造函数 ")

[AbStructHandle 构造函数 (Byte[], Int32)](../html/722ae95a-0a89-6493-d622-a99a692b0041.htm "AbStructHandle 构造函数 (Byte[], Int32)")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| AbStructHandle 构造函数 (Byte, Int32) |

使用原始字节的数据，索引信息来实例化一个对象  
Instantiate an object with raw bytes of data, index information

**命名空间：**
 [HslCommunication.Profinet.AllenBradley](fe226e46-d0b0-8ed5-5051-c925f25529d4.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public AbStructHandle(
	byte[] source,
	int index
)
```

```
Public Sub New ( 
	source As Byte(),
	index As Integer
)
```

```
public:
AbStructHandle(
	array<unsigned char>^ source, 
	int index
)
```

```
new : 
        source : byte[] * 
        index : int -> AbStructHandle
```

#### 参数

source
:   类型：SystemByte  
    原始字节数据

index
:   类型：SystemInt32  
    起始的偏移索引

![](../icons/SectionExpanded.png)参见

#### 引用

[AbStructHandle 类](7edf9ded-9d79-9b6d-9cce-63bc65b21a72.htm)

[AbStructHandle 重载](89182ea3-2ecb-e00c-b23a-502edc9b06cf.htm)

[HslCommunication.Profinet.AllenBradley 命名空间](fe226e46-d0b0-8ed5-5051-c925f25529d4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## AbStructHandle 属性

[原文連結](http://api.hslcommunication.cn/html/920cd83c-31b1-5cda-3afa-0c9caa05930d.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.AllenBradley](../html/fe226e46-d0b0-8ed5-5051-c925f25529d4.htm "HslCommunication.Profinet.AllenBradley")

[AbStructHandle 类](../html/7edf9ded-9d79-9b6d-9cce-63bc65b21a72.htm "AbStructHandle 类")

[AbStructHandle 属性](../html/920cd83c-31b1-5cda-3afa-0c9caa05930d.htm "AbStructHandle 属性")

[MemberCount 属性](../html/0fde9476-79d7-7e33-8d98-a489e0f996c0.htm "MemberCount 属性 ")

[ReturnCount 属性](../html/b3af0944-0bbf-f834-f0f8-60f0ca67ea5f.htm "ReturnCount 属性 ")

[StructureHandle 属性](../html/28a42c12-3a6f-84a7-cee4-23878dd52c80.htm "StructureHandle 属性 ")

[TemplateObjectDefinitionSize 属性](../html/e81a9baf-1b1a-4fa1-ae0e-040235eb5f02.htm "TemplateObjectDefinitionSize 属性 ")

[TemplateStructureSize 属性](../html/7c1de037-60f8-c0ae-73b5-3c72bb3dc970.htm "TemplateStructureSize 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| AbStructHandle 属性 |

[AbStructHandle](7edf9ded-9d79-9b6d-9cce-63bc65b21a72.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [MemberCount](0fde9476-79d7-7e33-8d98-a489e0f996c0.htm) | 成员数量 |
| 公共属性 | [ReturnCount](b3af0944-0bbf-f834-f0f8-60f0ca67ea5f.htm) | 返回项数 |
| 公共属性 | [StructureHandle](28a42c12-3a6f-84a7-cee4-23878dd52c80.htm) | 结构体的handle |
| 公共属性 | [TemplateObjectDefinitionSize](e81a9baf-1b1a-4fa1-ae0e-040235eb5f02.htm) | 结构体定义大小 |
| 公共属性 | [TemplateStructureSize](7c1de037-60f8-c0ae-73b5-3c72bb3dc970.htm) | 使用读取标记服务读取结构时在线路上传输的字节数 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[AbStructHandle 类](7edf9ded-9d79-9b6d-9cce-63bc65b21a72.htm)

[HslCommunication.Profinet.AllenBradley 命名空间](fe226e46-d0b0-8ed5-5051-c925f25529d4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## MemberCount 属性 

[原文連結](http://api.hslcommunication.cn/html/0fde9476-79d7-7e33-8d98-a489e0f996c0.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.AllenBradley](../html/fe226e46-d0b0-8ed5-5051-c925f25529d4.htm "HslCommunication.Profinet.AllenBradley")

[AbStructHandle 类](../html/7edf9ded-9d79-9b6d-9cce-63bc65b21a72.htm "AbStructHandle 类")

[AbStructHandle 属性](../html/920cd83c-31b1-5cda-3afa-0c9caa05930d.htm "AbStructHandle 属性")

[MemberCount 属性](../html/0fde9476-79d7-7e33-8d98-a489e0f996c0.htm "MemberCount 属性 ")

[ReturnCount 属性](../html/b3af0944-0bbf-f834-f0f8-60f0ca67ea5f.htm "ReturnCount 属性 ")

[StructureHandle 属性](../html/28a42c12-3a6f-84a7-cee4-23878dd52c80.htm "StructureHandle 属性 ")

[TemplateObjectDefinitionSize 属性](../html/e81a9baf-1b1a-4fa1-ae0e-040235eb5f02.htm "TemplateObjectDefinitionSize 属性 ")

[TemplateStructureSize 属性](../html/7c1de037-60f8-c0ae-73b5-3c72bb3dc970.htm "TemplateStructureSize 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| AbStructHandleMemberCount 属性 |

成员数量

**命名空间：**
 [HslCommunication.Profinet.AllenBradley](fe226e46-d0b0-8ed5-5051-c925f25529d4.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public ushort MemberCount { get; set; }
```

```
Public Property MemberCount As UShort
	Get
	Set
```

```
public:
property unsigned short MemberCount {
	unsigned short get ();
	void set (unsigned short value);
}
```

```
member MemberCount : uint16 with get, set
```

#### 属性值

类型：UInt16

![](../icons/SectionExpanded.png)备注

This is the number of structure members

![](../icons/SectionExpanded.png)参见

#### 引用

[AbStructHandle 类](7edf9ded-9d79-9b6d-9cce-63bc65b21a72.htm)

[HslCommunication.Profinet.AllenBradley 命名空间](fe226e46-d0b0-8ed5-5051-c925f25529d4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ReturnCount 属性 

[原文連結](http://api.hslcommunication.cn/html/b3af0944-0bbf-f834-f0f8-60f0ca67ea5f.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.AllenBradley](../html/fe226e46-d0b0-8ed5-5051-c925f25529d4.htm "HslCommunication.Profinet.AllenBradley")

[AbStructHandle 类](../html/7edf9ded-9d79-9b6d-9cce-63bc65b21a72.htm "AbStructHandle 类")

[AbStructHandle 属性](../html/920cd83c-31b1-5cda-3afa-0c9caa05930d.htm "AbStructHandle 属性")

[MemberCount 属性](../html/0fde9476-79d7-7e33-8d98-a489e0f996c0.htm "MemberCount 属性 ")

[ReturnCount 属性](../html/b3af0944-0bbf-f834-f0f8-60f0ca67ea5f.htm "ReturnCount 属性 ")

[StructureHandle 属性](../html/28a42c12-3a6f-84a7-cee4-23878dd52c80.htm "StructureHandle 属性 ")

[TemplateObjectDefinitionSize 属性](../html/e81a9baf-1b1a-4fa1-ae0e-040235eb5f02.htm "TemplateObjectDefinitionSize 属性 ")

[TemplateStructureSize 属性](../html/7c1de037-60f8-c0ae-73b5-3c72bb3dc970.htm "TemplateStructureSize 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| AbStructHandleReturnCount 属性 |

返回项数

**命名空间：**
 [HslCommunication.Profinet.AllenBradley](fe226e46-d0b0-8ed5-5051-c925f25529d4.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public ushort ReturnCount { get; set; }
```

```
Public Property ReturnCount As UShort
	Get
	Set
```

```
public:
property unsigned short ReturnCount {
	unsigned short get ();
	void set (unsigned short value);
}
```

```
member ReturnCount : uint16 with get, set
```

#### 属性值

类型：UInt16

![](../icons/SectionExpanded.png)备注

Count of Items returned

![](../icons/SectionExpanded.png)参见

#### 引用

[AbStructHandle 类](7edf9ded-9d79-9b6d-9cce-63bc65b21a72.htm)

[HslCommunication.Profinet.AllenBradley 命名空间](fe226e46-d0b0-8ed5-5051-c925f25529d4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## StructureHandle 属性 

[原文連結](http://api.hslcommunication.cn/html/28a42c12-3a6f-84a7-cee4-23878dd52c80.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.AllenBradley](../html/fe226e46-d0b0-8ed5-5051-c925f25529d4.htm "HslCommunication.Profinet.AllenBradley")

[AbStructHandle 类](../html/7edf9ded-9d79-9b6d-9cce-63bc65b21a72.htm "AbStructHandle 类")

[AbStructHandle 属性](../html/920cd83c-31b1-5cda-3afa-0c9caa05930d.htm "AbStructHandle 属性")

[MemberCount 属性](../html/0fde9476-79d7-7e33-8d98-a489e0f996c0.htm "MemberCount 属性 ")

[ReturnCount 属性](../html/b3af0944-0bbf-f834-f0f8-60f0ca67ea5f.htm "ReturnCount 属性 ")

[StructureHandle 属性](../html/28a42c12-3a6f-84a7-cee4-23878dd52c80.htm "StructureHandle 属性 ")

[TemplateObjectDefinitionSize 属性](../html/e81a9baf-1b1a-4fa1-ae0e-040235eb5f02.htm "TemplateObjectDefinitionSize 属性 ")

[TemplateStructureSize 属性](../html/7c1de037-60f8-c0ae-73b5-3c72bb3dc970.htm "TemplateStructureSize 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| AbStructHandleStructureHandle 属性 |

结构体的handle

**命名空间：**
 [HslCommunication.Profinet.AllenBradley](fe226e46-d0b0-8ed5-5051-c925f25529d4.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public ushort StructureHandle { get; set; }
```

```
Public Property StructureHandle As UShort
	Get
	Set
```

```
public:
property unsigned short StructureHandle {
	unsigned short get ();
	void set (unsigned short value);
}
```

```
member StructureHandle : uint16 with get, set
```

#### 属性值

类型：UInt16

![](../icons/SectionExpanded.png)备注

This is the Tag Type Parameter used in Read/Write Tag service

![](../icons/SectionExpanded.png)参见

#### 引用

[AbStructHandle 类](7edf9ded-9d79-9b6d-9cce-63bc65b21a72.htm)

[HslCommunication.Profinet.AllenBradley 命名空间](fe226e46-d0b0-8ed5-5051-c925f25529d4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## TemplateObjectDefinitionSize 属性 

[原文連結](http://api.hslcommunication.cn/html/e81a9baf-1b1a-4fa1-ae0e-040235eb5f02.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.AllenBradley](../html/fe226e46-d0b0-8ed5-5051-c925f25529d4.htm "HslCommunication.Profinet.AllenBradley")

[AbStructHandle 类](../html/7edf9ded-9d79-9b6d-9cce-63bc65b21a72.htm "AbStructHandle 类")

[AbStructHandle 属性](../html/920cd83c-31b1-5cda-3afa-0c9caa05930d.htm "AbStructHandle 属性")

[MemberCount 属性](../html/0fde9476-79d7-7e33-8d98-a489e0f996c0.htm "MemberCount 属性 ")

[ReturnCount 属性](../html/b3af0944-0bbf-f834-f0f8-60f0ca67ea5f.htm "ReturnCount 属性 ")

[StructureHandle 属性](../html/28a42c12-3a6f-84a7-cee4-23878dd52c80.htm "StructureHandle 属性 ")

[TemplateObjectDefinitionSize 属性](../html/e81a9baf-1b1a-4fa1-ae0e-040235eb5f02.htm "TemplateObjectDefinitionSize 属性 ")

[TemplateStructureSize 属性](../html/7c1de037-60f8-c0ae-73b5-3c72bb3dc970.htm "TemplateStructureSize 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| AbStructHandleTemplateObjectDefinitionSize 属性 |

结构体定义大小

**命名空间：**
 [HslCommunication.Profinet.AllenBradley](fe226e46-d0b0-8ed5-5051-c925f25529d4.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public uint TemplateObjectDefinitionSize { get; set; }
```

```
Public Property TemplateObjectDefinitionSize As UInteger
	Get
	Set
```

```
public:
property unsigned int TemplateObjectDefinitionSize {
	unsigned int get ();
	void set (unsigned int value);
}
```

```
member TemplateObjectDefinitionSize : uint32 with get, set
```

#### 属性值

类型：UInt32

![](../icons/SectionExpanded.png)备注

This is the number of structure members

![](../icons/SectionExpanded.png)参见

#### 引用

[AbStructHandle 类](7edf9ded-9d79-9b6d-9cce-63bc65b21a72.htm)

[HslCommunication.Profinet.AllenBradley 命名空间](fe226e46-d0b0-8ed5-5051-c925f25529d4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## TemplateStructureSize 属性 

[原文連結](http://api.hslcommunication.cn/html/7c1de037-60f8-c0ae-73b5-3c72bb3dc970.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.AllenBradley](../html/fe226e46-d0b0-8ed5-5051-c925f25529d4.htm "HslCommunication.Profinet.AllenBradley")

[AbStructHandle 类](../html/7edf9ded-9d79-9b6d-9cce-63bc65b21a72.htm "AbStructHandle 类")

[AbStructHandle 属性](../html/920cd83c-31b1-5cda-3afa-0c9caa05930d.htm "AbStructHandle 属性")

[MemberCount 属性](../html/0fde9476-79d7-7e33-8d98-a489e0f996c0.htm "MemberCount 属性 ")

[ReturnCount 属性](../html/b3af0944-0bbf-f834-f0f8-60f0ca67ea5f.htm "ReturnCount 属性 ")

[StructureHandle 属性](../html/28a42c12-3a6f-84a7-cee4-23878dd52c80.htm "StructureHandle 属性 ")

[TemplateObjectDefinitionSize 属性](../html/e81a9baf-1b1a-4fa1-ae0e-040235eb5f02.htm "TemplateObjectDefinitionSize 属性 ")

[TemplateStructureSize 属性](../html/7c1de037-60f8-c0ae-73b5-3c72bb3dc970.htm "TemplateStructureSize 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| AbStructHandleTemplateStructureSize 属性 |

使用读取标记服务读取结构时在线路上传输的字节数

**命名空间：**
 [HslCommunication.Profinet.AllenBradley](fe226e46-d0b0-8ed5-5051-c925f25529d4.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public uint TemplateStructureSize { get; set; }
```

```
Public Property TemplateStructureSize As UInteger
	Get
	Set
```

```
public:
property unsigned int TemplateStructureSize {
	unsigned int get ();
	void set (unsigned int value);
}
```

```
member TemplateStructureSize : uint32 with get, set
```

#### 属性值

类型：UInt32

![](../icons/SectionExpanded.png)备注

This is the number of bytes of the structure data

![](../icons/SectionExpanded.png)参见

#### 引用

[AbStructHandle 类](7edf9ded-9d79-9b6d-9cce-63bc65b21a72.htm)

[HslCommunication.Profinet.AllenBradley 命名空间](fe226e46-d0b0-8ed5-5051-c925f25529d4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## AbStructHandle 方法

[原文連結](http://api.hslcommunication.cn/html/2265386f-1c3c-527d-5d32-df0680221835.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.AllenBradley](../html/fe226e46-d0b0-8ed5-5051-c925f25529d4.htm "HslCommunication.Profinet.AllenBradley")

[AbStructHandle 类](../html/7edf9ded-9d79-9b6d-9cce-63bc65b21a72.htm "AbStructHandle 类")

[AbStructHandle 构造函数](../html/89182ea3-2ecb-e00c-b23a-502edc9b06cf.htm "AbStructHandle 构造函数 ")

[AbStructHandle 属性](../html/920cd83c-31b1-5cda-3afa-0c9caa05930d.htm "AbStructHandle 属性")

[AbStructHandle 方法](../html/2265386f-1c3c-527d-5d32-df0680221835.htm "AbStructHandle 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| AbStructHandle 方法 |

[AbStructHandle](7edf9ded-9d79-9b6d-9cce-63bc65b21a72.htm) 类型公开以下成员。

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

[AbStructHandle 类](7edf9ded-9d79-9b6d-9cce-63bc65b21a72.htm)

[HslCommunication.Profinet.AllenBradley 命名空间](fe226e46-d0b0-8ed5-5051-c925f25529d4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## AbTagItem 类

[原文連結](http://api.hslcommunication.cn/html/13d4e75c-bbdc-08cf-426f-799bb4099deb.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.AllenBradley](../html/fe226e46-d0b0-8ed5-5051-c925f25529d4.htm "HslCommunication.Profinet.AllenBradley")

[AbTagItem 类](../html/13d4e75c-bbdc-08cf-426f-799bb4099deb.htm "AbTagItem 类")

[AbTagItem 构造函数](../html/9c81e2cf-f3c8-70a9-095e-2687a7361b87.htm "AbTagItem 构造函数 ")

[AbTagItem 属性](../html/e883451b-7e80-5671-01d2-62bb8ca55996.htm "AbTagItem 属性")

[AbTagItem 方法](../html/e7d1a604-ae05-f304-2ec5-6b12a0db5d9f.htm "AbTagItem 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| AbTagItem 类 |

AB PLC的数据标签实体类  
Data tag entity class of AB PLC

![](../icons/SectionExpanded.png)继承层次

SystemObject  
  HslCommunication.Profinet.AllenBradleyAbTagItem

**命名空间：**
 [HslCommunication.Profinet.AllenBradley](fe226e46-d0b0-8ed5-5051-c925f25529d4.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public class AbTagItem
```

```
Public Class AbTagItem
```

```
public ref class AbTagItem
```

```
type AbTagItem =  class end
```

AbTagItem 类型公开以下成员。

![](../icons/SectionExpanded.png)构造函数

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法 | [AbTagItem](9c81e2cf-f3c8-70a9-095e-2687a7361b87.htm) | 实例化一个默认的对象  instantiate a default object |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [ArrayDimension](200ce643-e41a-b74f-468c-4c637a39f693.htm) | 数据的维度信息，默认是0，标量数据，1表示一维数组，2表示二维数组  The dimension information of the data, the default is 0, scalar data, 1 means a one-dimensional array, 2 means a two-dimensional array |
| 公共属性 | [ArrayLength](b232b555-a1e2-f3c0-2888-23b070007c73.htm) | 当前如果是数组，表示数组的长度，仅在读取结构体的变量信息时有效，为-1则是无效。 |
| 公共属性 | [ByteOffset](a21fa1e8-bb43-fbb4-37c5-bec95127a9ea.htm) | 获取或设置本属性实际数据在结构体中的偏移位置信息  Get or set the offset position information of the actual data of this property in the structure |
| 公共属性 | [InstanceID](b9457f34-dbf5-cf80-1682-d18373e64613.htm) | 实例ID  instance ID |
| 公共属性 | [IsStruct](681d2a85-6f7a-4d2c-ede8-72fe3e2763b7.htm) | 当前的标签是否结构体数据  Whether the current label is structured data |
| 公共属性 | [Members](07a4b5f5-5196-e873-3eba-aaa523d14e28.htm) | 如果当前的标签是结构体的标签，则表示为结构体的成员信息 |
| 公共属性 | [Name](62135725-d0d3-80e3-ed31-1aa3676b7d4f.htm) | 当前标签的名字  the name of the current label |
| 公共属性 | [SymbolType](74ea6511-d581-bcb1-7393-0cd594cce306.htm) | 当前标签的类型代号，例如 0x0C1 表示bool类型，如果当前的标签的[IsStruct](681d2a85-6f7a-4d2c-ede8-72fe3e2763b7.htm)为 True，那么本属性表示结构体的实例ID  The type code of the current tag, for example 0x0C1 means bool type, if the current tag's [IsStruct](681d2a85-6f7a-4d2c-ede8-72fe3e2763b7.htm) is True, then this attribute indicates the instance ID of the structure |
| 公共属性 | [Tag](04e96cbb-f325-2efa-ff3e-00c85d60b0f2.htm) | 用户自定义的额外的对象 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法静态成员 | [CloneBy(AbTagItem)](ed9932cb-f946-4ee4-2d8e-497eeb8b5d22.htm) | 克隆单个的标签数据信息 |
| 公共方法静态成员 | [CloneBy(AbTagItem)](8bf31d04-d7e3-fc21-cf5b-35ec3d29434f.htm) | 克隆整个的标签数组信息 |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法 | GetType | (继承自 Object。) |
| 公共方法 | [GetTypeText](8ebfb7ca-65b0-e615-9415-52cbbfa694e7.htm) | 获取类型的文本描述信息 |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法静态成员 | [PraseAbTagItem](8a285c81-fdbe-e2a1-049d-6d8e32a866e0.htm) | 从指定的原始字节的数据中，解析出实际的节点信息 |
| 公共方法静态成员 | [PraseAbTagItems](3de4c46b-50d7-f67c-c13f-3981a0459cda.htm) | 从指定的原始字节的数据中，解析出实际的标签数组，如果是系统保留的数组，或是\_\_开头的，则自动忽略。 |
| 公共方法静态成员 | [PraseAbTagItemsFromStruct](415ecdff-2640-5f19-5877-70d69b1c2862.htm) | 从结构体的数据中解析出实际的子标签信息 |
| 公共方法 | [ToString](a0378f14-1381-c49b-0a4e-77e2ff6b63c1.htm) | (重写 ObjectToString.) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[HslCommunication.Profinet.AllenBradley 命名空间](fe226e46-d0b0-8ed5-5051-c925f25529d4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## AbTagItem 构造函数 

[原文連結](http://api.hslcommunication.cn/html/9c81e2cf-f3c8-70a9-095e-2687a7361b87.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.AllenBradley](../html/fe226e46-d0b0-8ed5-5051-c925f25529d4.htm "HslCommunication.Profinet.AllenBradley")

[AbTagItem 类](../html/13d4e75c-bbdc-08cf-426f-799bb4099deb.htm "AbTagItem 类")

[AbTagItem 构造函数](../html/9c81e2cf-f3c8-70a9-095e-2687a7361b87.htm "AbTagItem 构造函数 ")

[AbTagItem 属性](../html/e883451b-7e80-5671-01d2-62bb8ca55996.htm "AbTagItem 属性")

[AbTagItem 方法](../html/e7d1a604-ae05-f304-2ec5-6b12a0db5d9f.htm "AbTagItem 方法")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| AbTagItem 构造函数 |

实例化一个默认的对象  
instantiate a default object

**命名空间：**
 [HslCommunication.Profinet.AllenBradley](fe226e46-d0b0-8ed5-5051-c925f25529d4.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public AbTagItem()
```

```
Public Sub New
```

```
public:
AbTagItem()
```

```
new : unit -> AbTagItem
```

![](../icons/SectionExpanded.png)参见

#### 引用

[AbTagItem 类](13d4e75c-bbdc-08cf-426f-799bb4099deb.htm)

[HslCommunication.Profinet.AllenBradley 命名空间](fe226e46-d0b0-8ed5-5051-c925f25529d4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## AbTagItem 属性

[原文連結](http://api.hslcommunication.cn/html/e883451b-7e80-5671-01d2-62bb8ca55996.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.AllenBradley](../html/fe226e46-d0b0-8ed5-5051-c925f25529d4.htm "HslCommunication.Profinet.AllenBradley")

[AbTagItem 类](../html/13d4e75c-bbdc-08cf-426f-799bb4099deb.htm "AbTagItem 类")

[AbTagItem 属性](../html/e883451b-7e80-5671-01d2-62bb8ca55996.htm "AbTagItem 属性")

[ArrayDimension 属性](../html/200ce643-e41a-b74f-468c-4c637a39f693.htm "ArrayDimension 属性 ")

[ArrayLength 属性](../html/b232b555-a1e2-f3c0-2888-23b070007c73.htm "ArrayLength 属性 ")

[ByteOffset 属性](../html/a21fa1e8-bb43-fbb4-37c5-bec95127a9ea.htm "ByteOffset 属性 ")

[InstanceID 属性](../html/b9457f34-dbf5-cf80-1682-d18373e64613.htm "InstanceID 属性 ")

[IsStruct 属性](../html/681d2a85-6f7a-4d2c-ede8-72fe3e2763b7.htm "IsStruct 属性 ")

[Members 属性](../html/07a4b5f5-5196-e873-3eba-aaa523d14e28.htm "Members 属性 ")

[Name 属性](../html/62135725-d0d3-80e3-ed31-1aa3676b7d4f.htm "Name 属性 ")

[SymbolType 属性](../html/74ea6511-d581-bcb1-7393-0cd594cce306.htm "SymbolType 属性 ")

[Tag 属性](../html/04e96cbb-f325-2efa-ff3e-00c85d60b0f2.htm "Tag 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| AbTagItem 属性 |

[AbTagItem](13d4e75c-bbdc-08cf-426f-799bb4099deb.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)属性

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共属性 | [ArrayDimension](200ce643-e41a-b74f-468c-4c637a39f693.htm) | 数据的维度信息，默认是0，标量数据，1表示一维数组，2表示二维数组  The dimension information of the data, the default is 0, scalar data, 1 means a one-dimensional array, 2 means a two-dimensional array |
| 公共属性 | [ArrayLength](b232b555-a1e2-f3c0-2888-23b070007c73.htm) | 当前如果是数组，表示数组的长度，仅在读取结构体的变量信息时有效，为-1则是无效。 |
| 公共属性 | [ByteOffset](a21fa1e8-bb43-fbb4-37c5-bec95127a9ea.htm) | 获取或设置本属性实际数据在结构体中的偏移位置信息  Get or set the offset position information of the actual data of this property in the structure |
| 公共属性 | [InstanceID](b9457f34-dbf5-cf80-1682-d18373e64613.htm) | 实例ID  instance ID |
| 公共属性 | [IsStruct](681d2a85-6f7a-4d2c-ede8-72fe3e2763b7.htm) | 当前的标签是否结构体数据  Whether the current label is structured data |
| 公共属性 | [Members](07a4b5f5-5196-e873-3eba-aaa523d14e28.htm) | 如果当前的标签是结构体的标签，则表示为结构体的成员信息 |
| 公共属性 | [Name](62135725-d0d3-80e3-ed31-1aa3676b7d4f.htm) | 当前标签的名字  the name of the current label |
| 公共属性 | [SymbolType](74ea6511-d581-bcb1-7393-0cd594cce306.htm) | 当前标签的类型代号，例如 0x0C1 表示bool类型，如果当前的标签的[IsStruct](681d2a85-6f7a-4d2c-ede8-72fe3e2763b7.htm)为 True，那么本属性表示结构体的实例ID  The type code of the current tag, for example 0x0C1 means bool type, if the current tag's [IsStruct](681d2a85-6f7a-4d2c-ede8-72fe3e2763b7.htm) is True, then this attribute indicates the instance ID of the structure |
| 公共属性 | [Tag](04e96cbb-f325-2efa-ff3e-00c85d60b0f2.htm) | 用户自定义的额外的对象 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[AbTagItem 类](13d4e75c-bbdc-08cf-426f-799bb4099deb.htm)

[HslCommunication.Profinet.AllenBradley 命名空间](fe226e46-d0b0-8ed5-5051-c925f25529d4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ArrayDimension 属性 

[原文連結](http://api.hslcommunication.cn/html/200ce643-e41a-b74f-468c-4c637a39f693.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.AllenBradley](../html/fe226e46-d0b0-8ed5-5051-c925f25529d4.htm "HslCommunication.Profinet.AllenBradley")

[AbTagItem 类](../html/13d4e75c-bbdc-08cf-426f-799bb4099deb.htm "AbTagItem 类")

[AbTagItem 属性](../html/e883451b-7e80-5671-01d2-62bb8ca55996.htm "AbTagItem 属性")

[ArrayDimension 属性](../html/200ce643-e41a-b74f-468c-4c637a39f693.htm "ArrayDimension 属性 ")

[ArrayLength 属性](../html/b232b555-a1e2-f3c0-2888-23b070007c73.htm "ArrayLength 属性 ")

[ByteOffset 属性](../html/a21fa1e8-bb43-fbb4-37c5-bec95127a9ea.htm "ByteOffset 属性 ")

[InstanceID 属性](../html/b9457f34-dbf5-cf80-1682-d18373e64613.htm "InstanceID 属性 ")

[IsStruct 属性](../html/681d2a85-6f7a-4d2c-ede8-72fe3e2763b7.htm "IsStruct 属性 ")

[Members 属性](../html/07a4b5f5-5196-e873-3eba-aaa523d14e28.htm "Members 属性 ")

[Name 属性](../html/62135725-d0d3-80e3-ed31-1aa3676b7d4f.htm "Name 属性 ")

[SymbolType 属性](../html/74ea6511-d581-bcb1-7393-0cd594cce306.htm "SymbolType 属性 ")

[Tag 属性](../html/04e96cbb-f325-2efa-ff3e-00c85d60b0f2.htm "Tag 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| AbTagItemArrayDimension 属性 |

数据的维度信息，默认是0，标量数据，1表示一维数组，2表示二维数组  
The dimension information of the data, the default is 0, scalar data, 1 means a one-dimensional array, 2 means a two-dimensional array

**命名空间：**
 [HslCommunication.Profinet.AllenBradley](fe226e46-d0b0-8ed5-5051-c925f25529d4.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public int ArrayDimension { get; set; }
```

```
Public Property ArrayDimension As Integer
	Get
	Set
```

```
public:
property int ArrayDimension {
	int get ();
	void set (int value);
}
```

```
member ArrayDimension : int with get, set
```

#### 属性值

类型：Int32

![](../icons/SectionExpanded.png)参见

#### 引用

[AbTagItem 类](13d4e75c-bbdc-08cf-426f-799bb4099deb.htm)

[HslCommunication.Profinet.AllenBradley 命名空间](fe226e46-d0b0-8ed5-5051-c925f25529d4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ArrayLength 属性 

[原文連結](http://api.hslcommunication.cn/html/b232b555-a1e2-f3c0-2888-23b070007c73.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.AllenBradley](../html/fe226e46-d0b0-8ed5-5051-c925f25529d4.htm "HslCommunication.Profinet.AllenBradley")

[AbTagItem 类](../html/13d4e75c-bbdc-08cf-426f-799bb4099deb.htm "AbTagItem 类")

[AbTagItem 属性](../html/e883451b-7e80-5671-01d2-62bb8ca55996.htm "AbTagItem 属性")

[ArrayDimension 属性](../html/200ce643-e41a-b74f-468c-4c637a39f693.htm "ArrayDimension 属性 ")

[ArrayLength 属性](../html/b232b555-a1e2-f3c0-2888-23b070007c73.htm "ArrayLength 属性 ")

[ByteOffset 属性](../html/a21fa1e8-bb43-fbb4-37c5-bec95127a9ea.htm "ByteOffset 属性 ")

[InstanceID 属性](../html/b9457f34-dbf5-cf80-1682-d18373e64613.htm "InstanceID 属性 ")

[IsStruct 属性](../html/681d2a85-6f7a-4d2c-ede8-72fe3e2763b7.htm "IsStruct 属性 ")

[Members 属性](../html/07a4b5f5-5196-e873-3eba-aaa523d14e28.htm "Members 属性 ")

[Name 属性](../html/62135725-d0d3-80e3-ed31-1aa3676b7d4f.htm "Name 属性 ")

[SymbolType 属性](../html/74ea6511-d581-bcb1-7393-0cd594cce306.htm "SymbolType 属性 ")

[Tag 属性](../html/04e96cbb-f325-2efa-ff3e-00c85d60b0f2.htm "Tag 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| AbTagItemArrayLength 属性 |

当前如果是数组，表示数组的长度，仅在读取结构体的变量信息时有效，为-1则是无效。

**命名空间：**
 [HslCommunication.Profinet.AllenBradley](fe226e46-d0b0-8ed5-5051-c925f25529d4.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public int[] ArrayLength { get; set; }
```

```
Public Property ArrayLength As Integer()
	Get
	Set
```

```
public:
property array<int>^ ArrayLength {
	array<int>^ get ();
	void set (array<int>^ value);
}
```

```
member ArrayLength : int[] with get, set
```

#### 属性值

类型：Int32

![](../icons/SectionExpanded.png)参见

#### 引用

[AbTagItem 类](13d4e75c-bbdc-08cf-426f-799bb4099deb.htm)

[HslCommunication.Profinet.AllenBradley 命名空间](fe226e46-d0b0-8ed5-5051-c925f25529d4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## ByteOffset 属性 

[原文連結](http://api.hslcommunication.cn/html/a21fa1e8-bb43-fbb4-37c5-bec95127a9ea.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.AllenBradley](../html/fe226e46-d0b0-8ed5-5051-c925f25529d4.htm "HslCommunication.Profinet.AllenBradley")

[AbTagItem 类](../html/13d4e75c-bbdc-08cf-426f-799bb4099deb.htm "AbTagItem 类")

[AbTagItem 属性](../html/e883451b-7e80-5671-01d2-62bb8ca55996.htm "AbTagItem 属性")

[ArrayDimension 属性](../html/200ce643-e41a-b74f-468c-4c637a39f693.htm "ArrayDimension 属性 ")

[ArrayLength 属性](../html/b232b555-a1e2-f3c0-2888-23b070007c73.htm "ArrayLength 属性 ")

[ByteOffset 属性](../html/a21fa1e8-bb43-fbb4-37c5-bec95127a9ea.htm "ByteOffset 属性 ")

[InstanceID 属性](../html/b9457f34-dbf5-cf80-1682-d18373e64613.htm "InstanceID 属性 ")

[IsStruct 属性](../html/681d2a85-6f7a-4d2c-ede8-72fe3e2763b7.htm "IsStruct 属性 ")

[Members 属性](../html/07a4b5f5-5196-e873-3eba-aaa523d14e28.htm "Members 属性 ")

[Name 属性](../html/62135725-d0d3-80e3-ed31-1aa3676b7d4f.htm "Name 属性 ")

[SymbolType 属性](../html/74ea6511-d581-bcb1-7393-0cd594cce306.htm "SymbolType 属性 ")

[Tag 属性](../html/04e96cbb-f325-2efa-ff3e-00c85d60b0f2.htm "Tag 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| AbTagItemByteOffset 属性 |

获取或设置本属性实际数据在结构体中的偏移位置信息  
Get or set the offset position information of the actual data of this property in the structure

**命名空间：**
 [HslCommunication.Profinet.AllenBradley](fe226e46-d0b0-8ed5-5051-c925f25529d4.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public int ByteOffset { get; set; }
```

```
Public Property ByteOffset As Integer
	Get
	Set
```

```
public:
property int ByteOffset {
	int get ();
	void set (int value);
}
```

```
member ByteOffset : int with get, set
```

#### 属性值

类型：Int32

![](../icons/SectionExpanded.png)参见

#### 引用

[AbTagItem 类](13d4e75c-bbdc-08cf-426f-799bb4099deb.htm)

[HslCommunication.Profinet.AllenBradley 命名空间](fe226e46-d0b0-8ed5-5051-c925f25529d4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## InstanceID 属性 

[原文連結](http://api.hslcommunication.cn/html/b9457f34-dbf5-cf80-1682-d18373e64613.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.AllenBradley](../html/fe226e46-d0b0-8ed5-5051-c925f25529d4.htm "HslCommunication.Profinet.AllenBradley")

[AbTagItem 类](../html/13d4e75c-bbdc-08cf-426f-799bb4099deb.htm "AbTagItem 类")

[AbTagItem 属性](../html/e883451b-7e80-5671-01d2-62bb8ca55996.htm "AbTagItem 属性")

[ArrayDimension 属性](../html/200ce643-e41a-b74f-468c-4c637a39f693.htm "ArrayDimension 属性 ")

[ArrayLength 属性](../html/b232b555-a1e2-f3c0-2888-23b070007c73.htm "ArrayLength 属性 ")

[ByteOffset 属性](../html/a21fa1e8-bb43-fbb4-37c5-bec95127a9ea.htm "ByteOffset 属性 ")

[InstanceID 属性](../html/b9457f34-dbf5-cf80-1682-d18373e64613.htm "InstanceID 属性 ")

[IsStruct 属性](../html/681d2a85-6f7a-4d2c-ede8-72fe3e2763b7.htm "IsStruct 属性 ")

[Members 属性](../html/07a4b5f5-5196-e873-3eba-aaa523d14e28.htm "Members 属性 ")

[Name 属性](../html/62135725-d0d3-80e3-ed31-1aa3676b7d4f.htm "Name 属性 ")

[SymbolType 属性](../html/74ea6511-d581-bcb1-7393-0cd594cce306.htm "SymbolType 属性 ")

[Tag 属性](../html/04e96cbb-f325-2efa-ff3e-00c85d60b0f2.htm "Tag 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| AbTagItemInstanceID 属性 |

实例ID  
instance ID

**命名空间：**
 [HslCommunication.Profinet.AllenBradley](fe226e46-d0b0-8ed5-5051-c925f25529d4.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public uint InstanceID { get; set; }
```

```
Public Property InstanceID As UInteger
	Get
	Set
```

```
public:
property unsigned int InstanceID {
	unsigned int get ();
	void set (unsigned int value);
}
```

```
member InstanceID : uint32 with get, set
```

#### 属性值

类型：UInt32

![](../icons/SectionExpanded.png)参见

#### 引用

[AbTagItem 类](13d4e75c-bbdc-08cf-426f-799bb4099deb.htm)

[HslCommunication.Profinet.AllenBradley 命名空间](fe226e46-d0b0-8ed5-5051-c925f25529d4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## IsStruct 属性 

[原文連結](http://api.hslcommunication.cn/html/681d2a85-6f7a-4d2c-ede8-72fe3e2763b7.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.AllenBradley](../html/fe226e46-d0b0-8ed5-5051-c925f25529d4.htm "HslCommunication.Profinet.AllenBradley")

[AbTagItem 类](../html/13d4e75c-bbdc-08cf-426f-799bb4099deb.htm "AbTagItem 类")

[AbTagItem 属性](../html/e883451b-7e80-5671-01d2-62bb8ca55996.htm "AbTagItem 属性")

[ArrayDimension 属性](../html/200ce643-e41a-b74f-468c-4c637a39f693.htm "ArrayDimension 属性 ")

[ArrayLength 属性](../html/b232b555-a1e2-f3c0-2888-23b070007c73.htm "ArrayLength 属性 ")

[ByteOffset 属性](../html/a21fa1e8-bb43-fbb4-37c5-bec95127a9ea.htm "ByteOffset 属性 ")

[InstanceID 属性](../html/b9457f34-dbf5-cf80-1682-d18373e64613.htm "InstanceID 属性 ")

[IsStruct 属性](../html/681d2a85-6f7a-4d2c-ede8-72fe3e2763b7.htm "IsStruct 属性 ")

[Members 属性](../html/07a4b5f5-5196-e873-3eba-aaa523d14e28.htm "Members 属性 ")

[Name 属性](../html/62135725-d0d3-80e3-ed31-1aa3676b7d4f.htm "Name 属性 ")

[SymbolType 属性](../html/74ea6511-d581-bcb1-7393-0cd594cce306.htm "SymbolType 属性 ")

[Tag 属性](../html/04e96cbb-f325-2efa-ff3e-00c85d60b0f2.htm "Tag 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| AbTagItemIsStruct 属性 |

当前的标签是否结构体数据  
Whether the current label is structured data

**命名空间：**
 [HslCommunication.Profinet.AllenBradley](fe226e46-d0b0-8ed5-5051-c925f25529d4.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public bool IsStruct { get; set; }
```

```
Public Property IsStruct As Boolean
	Get
	Set
```

```
public:
property bool IsStruct {
	bool get ();
	void set (bool value);
}
```

```
member IsStruct : bool with get, set
```

#### 属性值

类型：Boolean

![](../icons/SectionExpanded.png)参见

#### 引用

[AbTagItem 类](13d4e75c-bbdc-08cf-426f-799bb4099deb.htm)

[HslCommunication.Profinet.AllenBradley 命名空间](fe226e46-d0b0-8ed5-5051-c925f25529d4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Members 属性 

[原文連結](http://api.hslcommunication.cn/html/07a4b5f5-5196-e873-3eba-aaa523d14e28.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.AllenBradley](../html/fe226e46-d0b0-8ed5-5051-c925f25529d4.htm "HslCommunication.Profinet.AllenBradley")

[AbTagItem 类](../html/13d4e75c-bbdc-08cf-426f-799bb4099deb.htm "AbTagItem 类")

[AbTagItem 属性](../html/e883451b-7e80-5671-01d2-62bb8ca55996.htm "AbTagItem 属性")

[ArrayDimension 属性](../html/200ce643-e41a-b74f-468c-4c637a39f693.htm "ArrayDimension 属性 ")

[ArrayLength 属性](../html/b232b555-a1e2-f3c0-2888-23b070007c73.htm "ArrayLength 属性 ")

[ByteOffset 属性](../html/a21fa1e8-bb43-fbb4-37c5-bec95127a9ea.htm "ByteOffset 属性 ")

[InstanceID 属性](../html/b9457f34-dbf5-cf80-1682-d18373e64613.htm "InstanceID 属性 ")

[IsStruct 属性](../html/681d2a85-6f7a-4d2c-ede8-72fe3e2763b7.htm "IsStruct 属性 ")

[Members 属性](../html/07a4b5f5-5196-e873-3eba-aaa523d14e28.htm "Members 属性 ")

[Name 属性](../html/62135725-d0d3-80e3-ed31-1aa3676b7d4f.htm "Name 属性 ")

[SymbolType 属性](../html/74ea6511-d581-bcb1-7393-0cd594cce306.htm "SymbolType 属性 ")

[Tag 属性](../html/04e96cbb-f325-2efa-ff3e-00c85d60b0f2.htm "Tag 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| AbTagItemMembers 属性 |

如果当前的标签是结构体的标签，则表示为结构体的成员信息

**命名空间：**
 [HslCommunication.Profinet.AllenBradley](fe226e46-d0b0-8ed5-5051-c925f25529d4.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public AbTagItem[] Members { get; set; }
```

```
Public Property Members As AbTagItem()
	Get
	Set
```

```
public:
property array<AbTagItem^>^ Members {
	array<AbTagItem^>^ get ();
	void set (array<AbTagItem^>^ value);
}
```

```
member Members : AbTagItem[] with get, set
```

#### 属性值

类型：[AbTagItem](13d4e75c-bbdc-08cf-426f-799bb4099deb.htm)

![](../icons/SectionExpanded.png)参见

#### 引用

[AbTagItem 类](13d4e75c-bbdc-08cf-426f-799bb4099deb.htm)

[HslCommunication.Profinet.AllenBradley 命名空间](fe226e46-d0b0-8ed5-5051-c925f25529d4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Name 属性 

[原文連結](http://api.hslcommunication.cn/html/62135725-d0d3-80e3-ed31-1aa3676b7d4f.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.AllenBradley](../html/fe226e46-d0b0-8ed5-5051-c925f25529d4.htm "HslCommunication.Profinet.AllenBradley")

[AbTagItem 类](../html/13d4e75c-bbdc-08cf-426f-799bb4099deb.htm "AbTagItem 类")

[AbTagItem 属性](../html/e883451b-7e80-5671-01d2-62bb8ca55996.htm "AbTagItem 属性")

[ArrayDimension 属性](../html/200ce643-e41a-b74f-468c-4c637a39f693.htm "ArrayDimension 属性 ")

[ArrayLength 属性](../html/b232b555-a1e2-f3c0-2888-23b070007c73.htm "ArrayLength 属性 ")

[ByteOffset 属性](../html/a21fa1e8-bb43-fbb4-37c5-bec95127a9ea.htm "ByteOffset 属性 ")

[InstanceID 属性](../html/b9457f34-dbf5-cf80-1682-d18373e64613.htm "InstanceID 属性 ")

[IsStruct 属性](../html/681d2a85-6f7a-4d2c-ede8-72fe3e2763b7.htm "IsStruct 属性 ")

[Members 属性](../html/07a4b5f5-5196-e873-3eba-aaa523d14e28.htm "Members 属性 ")

[Name 属性](../html/62135725-d0d3-80e3-ed31-1aa3676b7d4f.htm "Name 属性 ")

[SymbolType 属性](../html/74ea6511-d581-bcb1-7393-0cd594cce306.htm "SymbolType 属性 ")

[Tag 属性](../html/04e96cbb-f325-2efa-ff3e-00c85d60b0f2.htm "Tag 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| AbTagItemName 属性 |

当前标签的名字  
the name of the current label

**命名空间：**
 [HslCommunication.Profinet.AllenBradley](fe226e46-d0b0-8ed5-5051-c925f25529d4.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public string Name { get; set; }
```

```
Public Property Name As String
	Get
	Set
```

```
public:
property String^ Name {
	String^ get ();
	void set (String^ value);
}
```

```
member Name : string with get, set
```

#### 属性值

类型：String

![](../icons/SectionExpanded.png)参见

#### 引用

[AbTagItem 类](13d4e75c-bbdc-08cf-426f-799bb4099deb.htm)

[HslCommunication.Profinet.AllenBradley 命名空间](fe226e46-d0b0-8ed5-5051-c925f25529d4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## SymbolType 属性 

[原文連結](http://api.hslcommunication.cn/html/74ea6511-d581-bcb1-7393-0cd594cce306.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.AllenBradley](../html/fe226e46-d0b0-8ed5-5051-c925f25529d4.htm "HslCommunication.Profinet.AllenBradley")

[AbTagItem 类](../html/13d4e75c-bbdc-08cf-426f-799bb4099deb.htm "AbTagItem 类")

[AbTagItem 属性](../html/e883451b-7e80-5671-01d2-62bb8ca55996.htm "AbTagItem 属性")

[ArrayDimension 属性](../html/200ce643-e41a-b74f-468c-4c637a39f693.htm "ArrayDimension 属性 ")

[ArrayLength 属性](../html/b232b555-a1e2-f3c0-2888-23b070007c73.htm "ArrayLength 属性 ")

[ByteOffset 属性](../html/a21fa1e8-bb43-fbb4-37c5-bec95127a9ea.htm "ByteOffset 属性 ")

[InstanceID 属性](../html/b9457f34-dbf5-cf80-1682-d18373e64613.htm "InstanceID 属性 ")

[IsStruct 属性](../html/681d2a85-6f7a-4d2c-ede8-72fe3e2763b7.htm "IsStruct 属性 ")

[Members 属性](../html/07a4b5f5-5196-e873-3eba-aaa523d14e28.htm "Members 属性 ")

[Name 属性](../html/62135725-d0d3-80e3-ed31-1aa3676b7d4f.htm "Name 属性 ")

[SymbolType 属性](../html/74ea6511-d581-bcb1-7393-0cd594cce306.htm "SymbolType 属性 ")

[Tag 属性](../html/04e96cbb-f325-2efa-ff3e-00c85d60b0f2.htm "Tag 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| AbTagItemSymbolType 属性 |

当前标签的类型代号，例如 0x0C1 表示bool类型，如果当前的标签的[IsStruct](681d2a85-6f7a-4d2c-ede8-72fe3e2763b7.htm)为 True，那么本属性表示结构体的实例ID  
The type code of the current tag, for example 0x0C1 means bool type, if the current tag's [IsStruct](681d2a85-6f7a-4d2c-ede8-72fe3e2763b7.htm) is True,
then this attribute indicates the instance ID of the structure

**命名空间：**
 [HslCommunication.Profinet.AllenBradley](fe226e46-d0b0-8ed5-5051-c925f25529d4.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public ushort SymbolType { get; set; }
```

```
Public Property SymbolType As UShort
	Get
	Set
```

```
public:
property unsigned short SymbolType {
	unsigned short get ();
	void set (unsigned short value);
}
```

```
member SymbolType : uint16 with get, set
```

#### 属性值

类型：UInt16

![](../icons/SectionExpanded.png)参见

#### 引用

[AbTagItem 类](13d4e75c-bbdc-08cf-426f-799bb4099deb.htm)

[HslCommunication.Profinet.AllenBradley 命名空间](fe226e46-d0b0-8ed5-5051-c925f25529d4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## Tag 属性 

[原文連結](http://api.hslcommunication.cn/html/04e96cbb-f325-2efa-ff3e-00c85d60b0f2.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.AllenBradley](../html/fe226e46-d0b0-8ed5-5051-c925f25529d4.htm "HslCommunication.Profinet.AllenBradley")

[AbTagItem 类](../html/13d4e75c-bbdc-08cf-426f-799bb4099deb.htm "AbTagItem 类")

[AbTagItem 属性](../html/e883451b-7e80-5671-01d2-62bb8ca55996.htm "AbTagItem 属性")

[ArrayDimension 属性](../html/200ce643-e41a-b74f-468c-4c637a39f693.htm "ArrayDimension 属性 ")

[ArrayLength 属性](../html/b232b555-a1e2-f3c0-2888-23b070007c73.htm "ArrayLength 属性 ")

[ByteOffset 属性](../html/a21fa1e8-bb43-fbb4-37c5-bec95127a9ea.htm "ByteOffset 属性 ")

[InstanceID 属性](../html/b9457f34-dbf5-cf80-1682-d18373e64613.htm "InstanceID 属性 ")

[IsStruct 属性](../html/681d2a85-6f7a-4d2c-ede8-72fe3e2763b7.htm "IsStruct 属性 ")

[Members 属性](../html/07a4b5f5-5196-e873-3eba-aaa523d14e28.htm "Members 属性 ")

[Name 属性](../html/62135725-d0d3-80e3-ed31-1aa3676b7d4f.htm "Name 属性 ")

[SymbolType 属性](../html/74ea6511-d581-bcb1-7393-0cd594cce306.htm "SymbolType 属性 ")

[Tag 属性](../html/04e96cbb-f325-2efa-ff3e-00c85d60b0f2.htm "Tag 属性 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| AbTagItemTag 属性 |

用户自定义的额外的对象

**命名空间：**
 [HslCommunication.Profinet.AllenBradley](fe226e46-d0b0-8ed5-5051-c925f25529d4.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public Object Tag { get; set; }
```

```
Public Property Tag As Object
	Get
	Set
```

```
public:
property Object^ Tag {
	Object^ get ();
	void set (Object^ value);
}
```

```
member Tag : Object with get, set
```

#### 属性值

类型：Object

![](../icons/SectionExpanded.png)参见

#### 引用

[AbTagItem 类](13d4e75c-bbdc-08cf-426f-799bb4099deb.htm)

[HslCommunication.Profinet.AllenBradley 命名空间](fe226e46-d0b0-8ed5-5051-c925f25529d4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## AbTagItem 方法

[原文連結](http://api.hslcommunication.cn/html/e7d1a604-ae05-f304-2ec5-6b12a0db5d9f.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.AllenBradley](../html/fe226e46-d0b0-8ed5-5051-c925f25529d4.htm "HslCommunication.Profinet.AllenBradley")

[AbTagItem 类](../html/13d4e75c-bbdc-08cf-426f-799bb4099deb.htm "AbTagItem 类")

[AbTagItem 方法](../html/e7d1a604-ae05-f304-2ec5-6b12a0db5d9f.htm "AbTagItem 方法")

[CloneBy 方法](../html/b26b10b1-5a73-c7cb-64ce-847611d0dba4.htm "CloneBy 方法 ")

[GetTypeText 方法](../html/8ebfb7ca-65b0-e615-9415-52cbbfa694e7.htm "GetTypeText 方法 ")

[PraseAbTagItem 方法](../html/8a285c81-fdbe-e2a1-049d-6d8e32a866e0.htm "PraseAbTagItem 方法 ")

[PraseAbTagItems 方法](../html/3de4c46b-50d7-f67c-c13f-3981a0459cda.htm "PraseAbTagItems 方法 ")

[PraseAbTagItemsFromStruct 方法](../html/415ecdff-2640-5f19-5877-70d69b1c2862.htm "PraseAbTagItemsFromStruct 方法 ")

[ToString 方法](../html/a0378f14-1381-c49b-0a4e-77e2ff6b63c1.htm "ToString 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| AbTagItem 方法 |

[AbTagItem](13d4e75c-bbdc-08cf-426f-799bb4099deb.htm) 类型公开以下成员。

![](../icons/SectionExpanded.png)方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法静态成员 | [CloneBy(AbTagItem)](ed9932cb-f946-4ee4-2d8e-497eeb8b5d22.htm) | 克隆单个的标签数据信息 |
| 公共方法静态成员 | [CloneBy(AbTagItem)](8bf31d04-d7e3-fc21-cf5b-35ec3d29434f.htm) | 克隆整个的标签数组信息 |
| 公共方法 | Equals | (继承自 Object。) |
| 受保护的方法 | Finalize | (继承自 Object。) |
| 公共方法 | GetHashCode | (继承自 Object。) |
| 公共方法 | GetType | (继承自 Object。) |
| 公共方法 | [GetTypeText](8ebfb7ca-65b0-e615-9415-52cbbfa694e7.htm) | 获取类型的文本描述信息 |
| 受保护的方法 | MemberwiseClone | (继承自 Object。) |
| 公共方法静态成员 | [PraseAbTagItem](8a285c81-fdbe-e2a1-049d-6d8e32a866e0.htm) | 从指定的原始字节的数据中，解析出实际的节点信息 |
| 公共方法静态成员 | [PraseAbTagItems](3de4c46b-50d7-f67c-c13f-3981a0459cda.htm) | 从指定的原始字节的数据中，解析出实际的标签数组，如果是系统保留的数组，或是\_\_开头的，则自动忽略。 |
| 公共方法静态成员 | [PraseAbTagItemsFromStruct](415ecdff-2640-5f19-5877-70d69b1c2862.htm) | 从结构体的数据中解析出实际的子标签信息 |
| 公共方法 | [ToString](a0378f14-1381-c49b-0a4e-77e2ff6b63c1.htm) | (重写 ObjectToString.) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)扩展方法

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共扩展器方法 | [ToJsonString](828b0b70-5381-936e-e606-e6c0a905f688.htm) | 获取当前对象的JSON格式表示的字符串。  Gets the string represented by the JSON format of the current object. (由 [HslExtension](05bbd254-20d0-2658-04fb-05c400c448dc.htm) 定义。) |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[AbTagItem 类](13d4e75c-bbdc-08cf-426f-799bb4099deb.htm)

[HslCommunication.Profinet.AllenBradley 命名空间](fe226e46-d0b0-8ed5-5051-c925f25529d4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## CloneBy 方法 

[原文連結](http://api.hslcommunication.cn/html/b26b10b1-5a73-c7cb-64ce-847611d0dba4.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.AllenBradley](../html/fe226e46-d0b0-8ed5-5051-c925f25529d4.htm "HslCommunication.Profinet.AllenBradley")

[AbTagItem 类](../html/13d4e75c-bbdc-08cf-426f-799bb4099deb.htm "AbTagItem 类")

[AbTagItem 方法](../html/e7d1a604-ae05-f304-2ec5-6b12a0db5d9f.htm "AbTagItem 方法")

[CloneBy 方法](../html/b26b10b1-5a73-c7cb-64ce-847611d0dba4.htm "CloneBy 方法 ")

[CloneBy 方法 (AbTagItem)](../html/ed9932cb-f946-4ee4-2d8e-497eeb8b5d22.htm "CloneBy 方法 (AbTagItem)")

[CloneBy 方法 (AbTagItem[])](../html/8bf31d04-d7e3-fc21-cf5b-35ec3d29434f.htm "CloneBy 方法 (AbTagItem[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| AbTagItemCloneBy 方法 |

![](../icons/SectionExpanded.png)重载列表

|  | 名称 | 说明 |
| --- | --- | --- |
| 公共方法静态成员 | [CloneBy(AbTagItem)](ed9932cb-f946-4ee4-2d8e-497eeb8b5d22.htm) | 克隆单个的标签数据信息 |
| 公共方法静态成员 | [CloneBy(AbTagItem)](8bf31d04-d7e3-fc21-cf5b-35ec3d29434f.htm) | 克隆整个的标签数组信息 |

[Top](#PageHeader)

![](../icons/SectionExpanded.png)参见

#### 引用

[AbTagItem 类](13d4e75c-bbdc-08cf-426f-799bb4099deb.htm)

[HslCommunication.Profinet.AllenBradley 命名空间](fe226e46-d0b0-8ed5-5051-c925f25529d4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## CloneBy 方法 (AbTagItem)

[原文連結](http://api.hslcommunication.cn/html/ed9932cb-f946-4ee4-2d8e-497eeb8b5d22.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.AllenBradley](../html/fe226e46-d0b0-8ed5-5051-c925f25529d4.htm "HslCommunication.Profinet.AllenBradley")

[AbTagItem 类](../html/13d4e75c-bbdc-08cf-426f-799bb4099deb.htm "AbTagItem 类")

[AbTagItem 方法](../html/e7d1a604-ae05-f304-2ec5-6b12a0db5d9f.htm "AbTagItem 方法")

[CloneBy 方法](../html/b26b10b1-5a73-c7cb-64ce-847611d0dba4.htm "CloneBy 方法 ")

[CloneBy 方法 (AbTagItem)](../html/ed9932cb-f946-4ee4-2d8e-497eeb8b5d22.htm "CloneBy 方法 (AbTagItem)")

[CloneBy 方法 (AbTagItem[])](../html/8bf31d04-d7e3-fc21-cf5b-35ec3d29434f.htm "CloneBy 方法 (AbTagItem[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| AbTagItemCloneBy 方法 (AbTagItem) |

克隆单个的标签数据信息

**命名空间：**
 [HslCommunication.Profinet.AllenBradley](fe226e46-d0b0-8ed5-5051-c925f25529d4.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static AbTagItem CloneBy(
	AbTagItem abTagItem
)
```

```
Public Shared Function CloneBy ( 
	abTagItem As AbTagItem
) As AbTagItem
```

```
public:
static AbTagItem^ CloneBy(
	AbTagItem^ abTagItem
)
```

```
static member CloneBy : 
        abTagItem : AbTagItem -> AbTagItem 
```

#### 参数

abTagItem
:   类型：[HslCommunication.Profinet.AllenBradleyAbTagItem](13d4e75c-bbdc-08cf-426f-799bb4099deb.htm)  
    标签信息

#### 返回值

类型：[AbTagItem](13d4e75c-bbdc-08cf-426f-799bb4099deb.htm)  
新的实例的标签

![](../icons/SectionExpanded.png)参见

#### 引用

[AbTagItem 类](13d4e75c-bbdc-08cf-426f-799bb4099deb.htm)

[CloneBy 重载](b26b10b1-5a73-c7cb-64ce-847611d0dba4.htm)

[HslCommunication.Profinet.AllenBradley 命名空间](fe226e46-d0b0-8ed5-5051-c925f25529d4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## CloneBy 方法 (AbTagItem[])

[原文連結](http://api.hslcommunication.cn/html/8bf31d04-d7e3-fc21-cf5b-35ec3d29434f.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.AllenBradley](../html/fe226e46-d0b0-8ed5-5051-c925f25529d4.htm "HslCommunication.Profinet.AllenBradley")

[AbTagItem 类](../html/13d4e75c-bbdc-08cf-426f-799bb4099deb.htm "AbTagItem 类")

[AbTagItem 方法](../html/e7d1a604-ae05-f304-2ec5-6b12a0db5d9f.htm "AbTagItem 方法")

[CloneBy 方法](../html/b26b10b1-5a73-c7cb-64ce-847611d0dba4.htm "CloneBy 方法 ")

[CloneBy 方法 (AbTagItem)](../html/ed9932cb-f946-4ee4-2d8e-497eeb8b5d22.htm "CloneBy 方法 (AbTagItem)")

[CloneBy 方法 (AbTagItem[])](../html/8bf31d04-d7e3-fc21-cf5b-35ec3d29434f.htm "CloneBy 方法 (AbTagItem[])")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| AbTagItemCloneBy 方法 (AbTagItem) |

克隆整个的标签数组信息

**命名空间：**
 [HslCommunication.Profinet.AllenBradley](fe226e46-d0b0-8ed5-5051-c925f25529d4.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static AbTagItem[] CloneBy(
	AbTagItem[] abTagItems
)
```

```
Public Shared Function CloneBy ( 
	abTagItems As AbTagItem()
) As AbTagItem()
```

```
public:
static array<AbTagItem^>^ CloneBy(
	array<AbTagItem^>^ abTagItems
)
```

```
static member CloneBy : 
        abTagItems : AbTagItem[] -> AbTagItem[] 
```

#### 参数

abTagItems
:   类型：[HslCommunication.Profinet.AllenBradleyAbTagItem](13d4e75c-bbdc-08cf-426f-799bb4099deb.htm)  
    标签数组信息

#### 返回值

类型：[AbTagItem](13d4e75c-bbdc-08cf-426f-799bb4099deb.htm)  
标签数组

![](../icons/SectionExpanded.png)参见

#### 引用

[AbTagItem 类](13d4e75c-bbdc-08cf-426f-799bb4099deb.htm)

[CloneBy 重载](b26b10b1-5a73-c7cb-64ce-847611d0dba4.htm)

[HslCommunication.Profinet.AllenBradley 命名空间](fe226e46-d0b0-8ed5-5051-c925f25529d4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## GetTypeText 方法 

[原文連結](http://api.hslcommunication.cn/html/8ebfb7ca-65b0-e615-9415-52cbbfa694e7.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.AllenBradley](../html/fe226e46-d0b0-8ed5-5051-c925f25529d4.htm "HslCommunication.Profinet.AllenBradley")

[AbTagItem 类](../html/13d4e75c-bbdc-08cf-426f-799bb4099deb.htm "AbTagItem 类")

[AbTagItem 方法](../html/e7d1a604-ae05-f304-2ec5-6b12a0db5d9f.htm "AbTagItem 方法")

[CloneBy 方法](../html/b26b10b1-5a73-c7cb-64ce-847611d0dba4.htm "CloneBy 方法 ")

[GetTypeText 方法](../html/8ebfb7ca-65b0-e615-9415-52cbbfa694e7.htm "GetTypeText 方法 ")

[PraseAbTagItem 方法](../html/8a285c81-fdbe-e2a1-049d-6d8e32a866e0.htm "PraseAbTagItem 方法 ")

[PraseAbTagItems 方法](../html/3de4c46b-50d7-f67c-c13f-3981a0459cda.htm "PraseAbTagItems 方法 ")

[PraseAbTagItemsFromStruct 方法](../html/415ecdff-2640-5f19-5877-70d69b1c2862.htm "PraseAbTagItemsFromStruct 方法 ")

[ToString 方法](../html/a0378f14-1381-c49b-0a4e-77e2ff6b63c1.htm "ToString 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| AbTagItemGetTypeText 方法 |

获取类型的文本描述信息

**命名空间：**
 [HslCommunication.Profinet.AllenBradley](fe226e46-d0b0-8ed5-5051-c925f25529d4.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public string GetTypeText()
```

```
Public Function GetTypeText As String
```

```
public:
String^ GetTypeText()
```

```
member GetTypeText : unit -> string 
```

#### 返回值

类型：String  
文本信息

![](../icons/SectionExpanded.png)参见

#### 引用

[AbTagItem 类](13d4e75c-bbdc-08cf-426f-799bb4099deb.htm)

[HslCommunication.Profinet.AllenBradley 命名空间](fe226e46-d0b0-8ed5-5051-c925f25529d4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)

---
## PraseAbTagItem 方法 

[原文連結](http://api.hslcommunication.cn/html/8a285c81-fdbe-e2a1-049d-6d8e32a866e0.htm)

HslCommunication Library Document

[HslCommunication Library Document](../html\c136d3de-eab7-9b0f-4bdf-d891297c8018.htm "HslCommunication Library Document")

[HslCommunication.Profinet.AllenBradley](../html/fe226e46-d0b0-8ed5-5051-c925f25529d4.htm "HslCommunication.Profinet.AllenBradley")

[AbTagItem 类](../html/13d4e75c-bbdc-08cf-426f-799bb4099deb.htm "AbTagItem 类")

[AbTagItem 方法](../html/e7d1a604-ae05-f304-2ec5-6b12a0db5d9f.htm "AbTagItem 方法")

[CloneBy 方法](../html/b26b10b1-5a73-c7cb-64ce-847611d0dba4.htm "CloneBy 方法 ")

[GetTypeText 方法](../html/8ebfb7ca-65b0-e615-9415-52cbbfa694e7.htm "GetTypeText 方法 ")

[PraseAbTagItem 方法](../html/8a285c81-fdbe-e2a1-049d-6d8e32a866e0.htm "PraseAbTagItem 方法 ")

[PraseAbTagItems 方法](../html/3de4c46b-50d7-f67c-c13f-3981a0459cda.htm "PraseAbTagItems 方法 ")

[PraseAbTagItemsFromStruct 方法](../html/415ecdff-2640-5f19-5877-70d69b1c2862.htm "PraseAbTagItemsFromStruct 方法 ")

[ToString 方法](../html/a0378f14-1381-c49b-0a4e-77e2ff6b63c1.htm "ToString 方法 ")

![点击或拖拽改变大小](../icons/TocOpen.gif "点击或拖拽改变大小")![点击或拖拽改变大小](../icons/TocClose.gif "点击或拖拽改变大小")

|  |
| --- |
| AbTagItemPraseAbTagItem 方法 |

从指定的原始字节的数据中，解析出实际的节点信息

**命名空间：**
 [HslCommunication.Profinet.AllenBradley](fe226e46-d0b0-8ed5-5051-c925f25529d4.htm)  
**程序集：**
 HslCommunication (在 HslCommunication.dll 中) 版本：12.5.3.0 (12.5.3.0)

![](../icons/SectionExpanded.png)语法

[C#](#)

[VB](#)

[C++](#)

[F#](#)

[复制](# "复制")

```
public static AbTagItem PraseAbTagItem(
	byte[] source,
	ref int index
)
```

```
Public Shared Function PraseAbTagItem ( 
	source As Byte(),
	ByRef index As Integer
) As AbTagItem
```

```
public:
static AbTagItem^ PraseAbTagItem(
	array<unsigned char>^ source, 
	int% index
)
```

```
static member PraseAbTagItem : 
        source : byte[] * 
        index : int byref -> AbTagItem 
```

#### 参数

source
:   类型：SystemByte  
    原始字节数据

index
:   类型：SystemInt32  
    起始的索引

#### 返回值

类型：[AbTagItem](13d4e75c-bbdc-08cf-426f-799bb4099deb.htm)  
标签信息

![](../icons/SectionExpanded.png)参见

#### 引用

[AbTagItem 类](13d4e75c-bbdc-08cf-426f-799bb4099deb.htm)

[HslCommunication.Profinet.AllenBradley 命名空间](fe226e46-d0b0-8ed5-5051-c925f25529d4.htm)

(C) 2017-2023 杭州胡工物联科技有限公司, All Rights Reserved

有关这个主题的评论请发邮件到
[hsl200909@163.com](mailto:hsl200909%40163.com?Subject=HslCommunication%20Library%20Document)